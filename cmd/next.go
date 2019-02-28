package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	semverUtil "github.com/meinto/git-semver/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var nextCmdOptions struct {
	RepoPath string
}

func init() {
	rootCmd.AddCommand(nextCmd)
	nextCmd.Flags().StringVarP(&versionCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "next version number",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nextVersionType := args[0]
		if nextVersionType != "major" && nextVersionType != "minor" && nextVersionType != "patch" {
			log.Fatalln("please choose one of these values: major, minor, patch")
		}

		gitRepoPath, err := filepath.Abs(nextCmdOptions.RepoPath)
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

			fmt.Println(nextVersion)
		}
	},
}
