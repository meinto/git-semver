package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/meinto/git-semver/cmd/internal"
	"github.com/pkg/errors"
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

		versionFileName, err := internal.PromptOptionalText("Name of version file")
		internal.LogFatalOnErr(err)
		config.VersionFileName = versionFileName

		_, versionFileType, err := internal.PromptSelect(
			"File type of version file",
			[]string{"json", "raw"},
		)
		internal.LogFatalOnErr(err)
		config.VersionFileType = versionFileType

		_, shouldBeTagged, err := internal.PromptSelect(
			"Should new version automatically be tagged",
			[]string{"yes", "no"},
		)
		internal.LogFatalOnErr(err)
		if shouldBeTagged == "yes" {
			config.TagVersions = true
		}

		_, changesShouldBePushed, err := internal.PromptSelect(
			"Should changes made by semver automatically be pushed",
			[]string{"yes", "no"},
		)
		internal.LogFatalOnErr(err)

		if changesShouldBePushed == "yes" {
			config.PushChanges = true

			author, err := internal.PromptOptionalText("Author of version commits")
			internal.LogFatalOnErr(err)
			if author != "" {
				config.Author = author
			}

			email, err := internal.PromptOptionalText("Email of version commits")
			internal.LogFatalOnErr(err)
			if email != "" {
				config.Email = email
			}
		}

		jsonContent, _ := json.MarshalIndent(config, "", "  ")
		err = ioutil.WriteFile("semver.config.json", jsonContent, 0644)
		internal.LogFatalOnErr(errors.Wrap(err, "error writing semver.config.json"))
	},
}
