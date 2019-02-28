package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	cmdUtil "github.com/meinto/git-semver/cmd/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init semver config",
	Run: func(cmd *cobra.Command, args []string) {

		var config struct {
			VersionFileName string `json:"versionFileName,omitempty"`
			VersionFileType string `json:"versionFileType,omitempty"`
			TagVersions     bool   `json:"tagVersions,omitempty"`
			PushChanges     bool   `json:"pushChanges,omitempty"`
			Author          string `json:"author,omitempty"`
			Email           string `json:"email,omitempty"`
		}

		versionFileName, _ := cmdUtil.PromptOptionalText("Name of version file")
		config.VersionFileName = versionFileName

		_, versionFileType, _ := cmdUtil.PromptSelect(
			"File type of version file",
			[]string{"json", "raw"},
		)
		config.VersionFileType = versionFileType

		_, shouldBeTagged, _ := cmdUtil.PromptSelect(
			"Should new version automatically be tagged",
			[]string{"yes", "no"},
		)
		if shouldBeTagged == "yes" {
			config.TagVersions = true
		}

		_, changesShouldBePushed, _ := cmdUtil.PromptSelect(
			"Should changes made by semver automatically be pushed",
			[]string{"yes", "no"},
		)

		if changesShouldBePushed == "yes" {
			config.PushChanges = true

			author, _ := cmdUtil.PromptOptionalText("Author of version commits")
			if author != "" {
				config.Author = author
			}

			email, _ := cmdUtil.PromptOptionalText("Email of version commits")
			if email != "" {
				config.Email = email
			}
		}

		jsonContent, _ := json.MarshalIndent(config, "", "  ")
		err := ioutil.WriteFile("semver.config.json", jsonContent, 0644)
		if err != nil {
			log.Fatalf("error writing semver.config.json: %s", err.Error())
		}
	},
}
