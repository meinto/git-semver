package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	cmdUtil "github.com/meinto/git-semver/cmd/internal/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "install",
	Short: "install semver",
	Run: func(cmd *cobra.Command, args []string) {

		flist, err := fileList(".")
		if err != nil {
			log.Fatalf("file listing failed: %s", err.Error())
		}

		index, err := cmdUtil.PromptSelect(
			"Select your downloaded semver file",
			flist,
		)
		if err != nil {
			log.Fatal(err)
		}
		filePath, err := filepath.Abs(flist[index])
		if err != nil {
			log.Fatalf("error getting path to semver file: %s", err)
		}

		index, err = cmdUtil.PromptSelect(
			"How do you want to use semver",
			[]string{"global", "git plugin"},
		)
		if err != nil {
			log.Fatal(err)
		}

		var newFileName string
		switch index {
		case 0:
			newFileName = "/usr/local/bin/semver"
		case 1:
			newFileName = "/usr/local/bin/git-semver"
		}

		if _, err := os.Stat(newFileName); !os.IsNotExist(err) {
			replace, err := replaceFile(newFileName)
			if err != nil {
				log.Fatal(err)
			}
			if !replace {
				log.Fatal("file not replaced")
			}
		}

		err = os.Rename(filePath, newFileName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("successfully moved semver")
	},
}

func replaceFile(filePath string) (bool, error) {
	prompt := promptui.Select{
		Label: "File exists. Do you want to replace it",
		Items: []string{"yes", "no"},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return false, err
	}

	if index == 0 {
		return true, nil
	}
	return false, nil
}

func fileList(rootPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if path == "." {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}
