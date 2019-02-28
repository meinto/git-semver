package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	semverUtil "github.com/meinto/git-semver/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmdOptions struct {
	RepoPath          string
	VersionFile       string
	VersionFileFormat string
	DryRun            bool
	CreateTag         bool
	Push              bool
	Author            string
	Email             string
	SSHFilePath       string
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringVarP(&versionCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
	versionCmd.Flags().StringVarP(&versionCmdOptions.Author, "author", "a", "semver", "name of the author")
	versionCmd.Flags().StringVarP(&versionCmdOptions.Email, "email", "e", "semver@no-reply.git", "email of the author")
	versionCmd.Flags().StringVarP(&versionCmdOptions.VersionFile, "outfile", "o", "semver.json", "name of version file")
	versionCmd.Flags().StringVarP(&versionCmdOptions.VersionFileFormat, "outfileFormat", "f", "json", "format of outfile (json, raw)")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.DryRun, "dryrun", "d", false, "only log how version number would change")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.CreateTag, "tag", "t", false, "create a git tag")
	versionCmd.Flags().BoolVarP(&versionCmdOptions.Push, "push", "P", false, "push git tags and version changes")

	viper.BindPFlag("versionFileName", versionCmd.Flags().Lookup("outfile"))
	viper.BindPFlag("versionFileType", versionCmd.Flags().Lookup("outfileFormat"))
	viper.BindPFlag("tagVersions", versionCmd.Flags().Lookup("tag"))
	viper.BindPFlag("pushChanges", versionCmd.Flags().Lookup("push"))
	viper.BindPFlag("author", versionCmd.Flags().Lookup("author"))
	viper.BindPFlag("email", versionCmd.Flags().Lookup("email"))

	defaultSSHFilePath, err := semverUtil.GetDefaultSSHFilePath()
	if err != nil {
		log.Println(err)
	}
	versionCmd.Flags().StringVar(&versionCmdOptions.SSHFilePath, "sshFilePath", defaultSSHFilePath, "path to your ssh file")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "create new version for repository",
	Args:  cobra.MinimumNArgs(1),
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

		var jsonContent = make(map[string]interface{})
		pathToVersionFile := gitRepoPath + "/" + viper.GetString("versionFileName")
		if _, err := os.Stat(pathToVersionFile); os.IsNotExist(err) {
			log.Printf("%s doesn't exist. creating one...", viper.GetString("versionFileName"))
			jsonContent = make(map[string]interface{})
			jsonContent["version"] = "1.0.0"
		} else {
			versionFile, err := os.Open(pathToVersionFile)
			if err != nil {
				log.Fatalf("cannot read %s: %s", viper.GetString("versionFileName"), err.Error())
			}
			defer versionFile.Close()

			byteValue, _ := ioutil.ReadAll(versionFile)

			switch viper.GetString("versionFileType") {
			case "json":
				json.Unmarshal(byteValue, &jsonContent)
			case "raw":
				jsonContent["version"] = string(byteValue)
			}

			currentVersion, ok := jsonContent["version"]
			if !ok {
				log.Fatalln("current version not set")
			}
			nextVersion, err := semverUtil.NextVersion(currentVersion.(string), nextVersionType)
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

		if viper.GetBool("pushChanges") {
			if err = semverUtil.CheckIfRepoIsClean(versionCmdOptions.RepoPath); err != nil {
				log.Fatal(err)
			}
			if err = semverUtil.CheckIfSSHFileExists(versionCmdOptions.SSHFilePath); err != nil {
				log.Fatal(err)
			}
		}

		switch viper.GetString("versionFileType") {
		case "json":
			err = semverUtil.WriteJSONVersionFile(jsonContent, viper.GetString("versionFileName"))
		case "raw":
			err = semverUtil.WriteRAWVersionFile(jsonContent["version"].(string), viper.GetString("versionFileName"))
		}
		if err != nil {
			log.Fatal(err)
		}
		if err = semverUtil.AddVersionChanges(
			versionCmdOptions.RepoPath,
			viper.GetString("versionFileName"),
			nextVersion,
			viper.GetString("author"),
			viper.GetString("email"),
		); err != nil {
			log.Fatal(err)
		}

		var createGitTagError error
		if viper.GetBool("tagVersions") {
			createGitTagError = semverUtil.MakeGitTag(versionCmdOptions.RepoPath, nextVersion)
		}

		if viper.GetBool("pushChanges") && createGitTagError == nil {
			if err = semverUtil.Push(versionCmdOptions.RepoPath, versionCmdOptions.SSHFilePath); err != nil {
				log.Fatalf("cannot push tag: %s", err.Error())
			}
		}
	},
}
