package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "install",
	Short: "install semver",
	Run: func(cmd *cobra.Command, args []string) {

		filePath, err := pathToSemverFile()
		if err != nil {
			log.Fatalf("error getting path to semver file: %s", err)
		}

		index, err := usageOptions()
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

func pathToSemverFile() (string, error) {
	validate := func(input string) error {
		filePath, err := filepath.Abs(input)
		if err != nil {
			return fmt.Errorf("error while creating absolute path: %s", err.Error())
		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return errors.New("file does not exist")
		}
		return nil
	}

	getFileName := promptui.Prompt{
		Label:    "Name of binary: ",
		Validate: validate,
	}

	fileName, err := getFileName.Run()
	if err != nil {
		return "", err
	}

	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func usageOptions() (int, error) {
	prompt := promptui.Select{
		Label: "How do you want to use semver?",
		Items: []string{"global", "git plugin"},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return index, nil
}

func replaceFile(filePath string) (bool, error) {
	prompt := promptui.Select{
		Label: "File exists. Do you want to replace it?",
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
