package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/meinto/cobra-utils"
	"github.com/meinto/git-semver/pkg/cli/cmd/internal"
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

		l := internal.NewLogger(rootCmdFlags.verbose)

		var config struct {
			VersionFile     string `json:"versionFile,omitempty"`
			VersionFileType string `json:"versionFileType,omitempty"`
			TagVersions     bool   `json:"tagVersions,omitempty"`
			PushChanges     bool   `json:"pushChanges,omitempty"`
		}

		versionFile, err := internal.PromptOptionalText("Name of version file")
		l.LogFatalOnError(err)
		config.VersionFile = versionFile

		_, versionFileType, err := cobraUtils.PromptSelect(
			"File type of version file",
			[]string{"json", "raw"},
		)
		l.LogFatalOnError(err)
		config.VersionFileType = versionFileType

		_, shouldBeTagged, err := cobraUtils.PromptSelect(
			"Should new version automatically be tagged",
			[]string{"yes", "no"},
		)
		l.LogFatalOnError(err)
		if shouldBeTagged == "yes" {
			config.TagVersions = true
		}

		_, changesShouldBePushed, err := cobraUtils.PromptSelect(
			"Should changes made by semver automatically be pushed",
			[]string{"yes", "no"},
		)
		l.LogFatalOnError(err)

		if changesShouldBePushed == "yes" {
			config.PushChanges = true
		}

		jsonContent, _ := json.MarshalIndent(config, "", "  ")
		err = ioutil.WriteFile("semver.config.json", jsonContent, 0644)
		l.LogFatalOnError((errors.Wrap(err, "error writing semver.config.json")))
	},
}
