package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmdOptions struct {
	RepoPath string
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().StringVarP(&versionCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
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

		var versionFileJSON map[string]interface{}
		pathToVersionFile := gitRepoPath + "/semver.json"
		if _, err := os.Stat(pathToVersionFile); os.IsNotExist(err) {
			log.Println("semver.json doesn't exist. creating one...")
			versionFileJSON = make(map[string]interface{})
			versionFileJSON["version"] = "1.0.0"
		} else {
			versionFile, err := os.Open(pathToVersionFile)
			if err != nil {
				log.Fatalln("cannot read semver.json: ", err)
			}
			defer versionFile.Close()

			byteValue, _ := ioutil.ReadAll(versionFile)
			json.Unmarshal(byteValue, &versionFileJSON)

			currentVersion, ok := versionFileJSON["version"]
			if !ok {
				log.Fatalln("current version not set")
			}
			nextVersion, err := makeVersion(currentVersion.(string), nextVersionType)
			if err != nil {
				log.Fatalln(err)
			}

			versionFileJSON["version"] = nextVersion
		}

		newVersionFileJSON, _ := json.MarshalIndent(versionFileJSON, "", "  ")
		err = ioutil.WriteFile("semver.json", newVersionFileJSON, 0644)
		if err != nil {
			log.Fatalln("error writing semver.json: ", err)
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
