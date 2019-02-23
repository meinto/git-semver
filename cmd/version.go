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

		var packageJSON map[string]interface{}
		pathToPackageJSON := gitRepoPath + "/package.json"
		if _, err := os.Stat(pathToPackageJSON); os.IsNotExist(err) {
			log.Println("package.json doesn't exist. creating one...")
			packageJSON = make(map[string]interface{})
			packageJSON["version"] = "1.0.0"
		} else {
			packageJSONFile, err := os.Open(pathToPackageJSON)
			if err != nil {
				log.Fatalln("cannot read package.json: ", err)
			}
			defer packageJSONFile.Close()

			byteValue, _ := ioutil.ReadAll(packageJSONFile)
			json.Unmarshal(byteValue, &packageJSON)

			currentVersion, ok := packageJSON["version"]
			if !ok {
				log.Fatalln("current version not set")
			}
			nextVersion, err := makeVersion(currentVersion.(string), nextVersionType)
			if err != nil {
				log.Fatalln(err)
			}

			packageJSON["version"] = nextVersion
		}

		newPackageJSON, _ := json.MarshalIndent(packageJSON, "", "  ")
		err = ioutil.WriteFile("package.json", newPackageJSON, 0644)
		if err != nil {
			log.Fatalln("error writing package.json: ", err)
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
