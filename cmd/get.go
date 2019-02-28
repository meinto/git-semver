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

var getCmdOptions struct {
	RepoPath string
	PrintRaw bool
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&getCmdOptions.RepoPath, "path", "p", ".", "path to git repository")
	getCmd.Flags().BoolVarP(&getCmdOptions.PrintRaw, "raw", "r", false, "print only the plain version number")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get version number",
	Run: func(cmd *cobra.Command, args []string) {
		gitRepoPath, err := filepath.Abs(getCmdOptions.RepoPath)
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

			if len(args) > 0 {
				nextVersionType := args[0]
				if nextVersionType != "major" && nextVersionType != "minor" && nextVersionType != "patch" {
					log.Fatalln("please choose one of these values: major, minor, patch")
				}

				nextVersion, err := semverUtil.NextVersion(currentVersion.(string), nextVersionType)
				if err != nil {
					log.Fatalln(err)
				}

				if !getCmdOptions.PrintRaw {
					fmt.Printf("Next %s version: ", nextVersionType)
				}
				fmt.Println(nextVersion)
			} else {
				if !getCmdOptions.PrintRaw {
					fmt.Print("Current version: ")
				}
				fmt.Println(currentVersion)
			}
		}
	},
}
