package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var versionCmdOptions struct {
	RepoPath    string
	VersionFile string
	DryRun      bool
	CreateTag   bool
	Push        bool
	Author      string
	Email       string
	SSHFilePath string
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
	versionCmd.Flags().StringVarP(&versionCmdOptions.Author, "author", "a", "semver", "name of the author")
	versionCmd.Flags().StringVarP(&versionCmdOptions.Email, "email", "e", "semver@no-reply.git", "email of the author")
	versionCmd.Flags().StringVarP(&versionCmdOptions.VersionFile, "outfile", "o", "semver.json", "name of version file")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.DryRun, "dryrun", "d", false, "only log how version number would change")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.CreateTag, "tag", "t", false, "create a git tag")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.Push, "push", "P", false, "push git tags and version changes")

	currentUser, err := user.Current()
	var defaultSSHFilePath string
	if err != nil {
		log.Println("cannot set default ssh file path")
		defaultSSHFilePath = ""
	} else {
		defaultSSHFilePath = currentUser.HomeDir + "/.ssh/id_rsa"
	}
	versionCmd.Flags().StringVar(&versionCmdOptions.SSHFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "create new version for repository",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalln("please provide the next version type (major, minor, patch).")
		}

		nextVersionType := args[0]
		if nextVersionType != "major" && nextVersionType != "minor" && nextVersionType != "patch" {
			log.Fatalln("please choose one of these values: major, minor, patch")
		}

		gitRepoPath, err := filepath.Abs(versionCmdOptions.RepoPath)
		if err != nil {
			log.Fatalln("cannot resolve repo path: ", err)
		}

		var jsonContent map[string]interface{}
		pathToVersionFile := gitRepoPath + "/" + versionCmdOptions.VersionFile
		if _, err := os.Stat(pathToVersionFile); os.IsNotExist(err) {
			log.Printf("%s doesn't exist. creating one...", versionCmdOptions.VersionFile)
			jsonContent = make(map[string]interface{})
			jsonContent["version"] = "1.0.0"
		} else {
			versionFile, err := os.Open(pathToVersionFile)
			if err != nil {
				log.Fatalf("cannot read %s: %s", versionCmdOptions.VersionFile, err.Error())
			}
			defer versionFile.Close()

			byteValue, _ := ioutil.ReadAll(versionFile)
			json.Unmarshal(byteValue, &jsonContent)

			currentVersion, ok := jsonContent["version"]
			if !ok {
				log.Fatalln("current version not set")
			}
			nextVersion, err := makeVersion(currentVersion.(string), nextVersionType)
			if err != nil {
				log.Fatalln(err)
			}

			jsonContent["version"] = nextVersion
		}

		nextVersion := jsonContent["version"].(string)
		log.Println("new version: ", nextVersion)
		if versionCmdOptions.DryRun {
			log.Println("dry run finished...")
			os.Exit(1)
		}

		if versionCmdOptions.Push {
			if err = checkIfRepoIsClean(versionCmdOptions.RepoPath); err != nil {
				log.Fatal(err)
			}
			if err = checkIfSSHFileExists(versionCmdOptions.SSHFilePath); err != nil {
				log.Fatal(err)
			}
		}

		writeVersionFile(jsonContent)
		if err = addVersionChanges(versionCmdOptions.RepoPath, versionCmdOptions.VersionFile, nextVersion); err != nil {
			log.Fatal(err)
		}

		var createGitTagError error
		if versionCmdOptions.CreateTag {
			createGitTagError = makeGitTag(versionCmdOptions.RepoPath, nextVersion)
		}

		if versionCmdOptions.Push && createGitTagError == nil {
			if err = push(versionCmdOptions.RepoPath); err != nil {
				log.Fatalf("cannot push tag: %s", err.Error())
			}
		}
	},
}

func makeVersion(currentVersion, nextVersionType string) (string, error) {
	numbers := strings.Split(currentVersion, ".")
	if len(numbers) != 3 {
		return "", errors.New("please provide version number in the following format: <major>.<minor>.<patch>")
	}

	switch nextVersionType {
	case "major":
		major, _ := strconv.Atoi(numbers[0])
		numbers[0] = strconv.Itoa(major + 1)
		numbers[1] = "0"
		numbers[2] = "0"
	case "minor":
		minor, _ := strconv.Atoi(numbers[1])
		numbers[1] = strconv.Itoa(minor + 1)
		numbers[2] = "0"
	case "patch":
		patch, _ := strconv.Atoi(numbers[2])
		numbers[2] = strconv.Itoa(patch + 1)
	}

	return strings.Join(numbers, "."), nil
}

func writeVersionFile(jsonContent map[string]interface{}) {
	newJSONContent, _ := json.MarshalIndent(jsonContent, "", "  ")
	err := ioutil.WriteFile(versionCmdOptions.VersionFile, newJSONContent, 0644)
	if err != nil {
		log.Fatalf("error writing %s: %s", versionCmdOptions.VersionFile, err.Error())
	}
}

func makeGitTag(repoPath, version string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
	} else {
		headRef, err := r.Head()
		if err != nil {
			return err
		}

		tag := fmt.Sprintf("refs/tags/v%s", version)
		ref := plumbing.NewHashReference(plumbing.ReferenceName(tag), headRef.Hash())

		err = r.Storer.SetReference(ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkIfRepoIsClean(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}
	status, err := w.Status()
	if err != nil {
		return err
	}
	if !status.IsClean() {
		return errors.New("please commit all files before versioning")
	}
	return nil
}

func addVersionChanges(repoPath, configFile, version string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(configFile)
	if err != nil {
		return err
	}
	_, err = w.Commit("new version: "+version, &git.CommitOptions{
		Author: &object.Signature{
			Name:  versionCmdOptions.Author,
			Email: versionCmdOptions.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func push(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", versionCmdOptions.SSHFilePath, "")
	if err != nil {
		return err
	}

	tags := config.RefSpec("refs/tags/*:refs/tags/*")
	heads := config.RefSpec("refs/heads/*:refs/heads/*")
	err = r.Push(&git.PushOptions{
		Auth:     sshAuth,
		RefSpecs: []config.RefSpec{tags, heads},
	})
	if err != nil {
		return err
	}
	return nil
}

func checkIfSSHFileExists(sshFilePath string) error {
	if _, err := os.Stat(sshFilePath); os.IsNotExist(err) {
		return fmt.Errorf("ssh file not found: %s", err.Error())
	}
	return nil
}
