package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	cmdUtil "github.com/meinto/git-semver/cmd/internal/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install semver",
	Run: func(cmd *cobra.Command, args []string) {

		flist, err := fileList(".")
		cmdUtil.LogFatalOnErr(errors.Wrap(err, "file listing failed"))

		index, _, err := cmdUtil.PromptSelect(
			"Select your downloaded semver file",
			flist,
		)
		cmdUtil.LogFatalOnErr(err)

		filePath, err := filepath.Abs(flist[index])
		cmdUtil.LogFatalOnErr(errors.Wrap(err, "error getting path to semver file"))

		index, _, err = cmdUtil.PromptSelect(
			"How do you want to use semver",
			[]string{"global", "git plugin"},
		)
		cmdUtil.LogFatalOnErr(err)

		var newFileName string
		switch index {
		case 0:
			newFileName = "/usr/local/bin/semver"
		case 1:
			newFileName = "/usr/local/bin/git-semver"
		}

		if _, err := os.Stat(newFileName); !os.IsNotExist(err) {
			err := replaceFile(newFileName)
			cmdUtil.LogFatalOnErr(err)
		}

		err = os.Rename(filePath, newFileName)
		cmdUtil.LogFatalOnErr(err)

		fmt.Println("successfully moved semver")
	},
}

func replaceFile(filePath string) error {
	prompt := promptui.Select{
		Label: "File exists. Do you want to replace it",
		Items: []string{"yes", "no"},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return err
	}

	if index == 0 {
		return nil
	}
	return errors.New("file not replaced")
}

func fileList(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !strings.HasPrefix(path, "semver_") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, errors.Wrap(err, "error creating file list")
	}

	return files, nil
}
