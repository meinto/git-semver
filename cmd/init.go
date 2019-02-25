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
	Use:   "init",
	Short: "init semver",
	Run: func(cmd *cobra.Command, args []string) {
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
			Label:    "Name of binary",
			Validate: validate,
		}

		fileName, err := getFileName.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		filePath, err := filepath.Abs(fileName)
		if err != nil {
			log.Fatalf("error while creating absolute path: %s", err.Error())
		}

		prompt := promptui.Select{
			Label: "How do you want to use semver?",
			Items: []string{"global", "git plugin"},
		}

		index, _, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		var newFileName string
		switch index {
		case 0:
			newFileName = "/usr/local/bin/semver"
		case 1:
			newFileName = "/usr/local/bin/git-semver"
		}

		err = os.Rename(filePath, newFileName)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("successfully moved semver")
	},
}
