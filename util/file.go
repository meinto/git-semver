package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

// WriteVersionJSONFile writes given jsonContent back to json file
func WriteVersionJSONFile(jsonContent map[string]interface{}, versionFile string) error {
	newJSONContent, _ := json.MarshalIndent(jsonContent, "", "  ")
	err := ioutil.WriteFile(versionFile, newJSONContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing %s: %s", versionFile, err.Error())
	}
	return nil
}

// CheckIfSSHFileExists validate the existence of the given ssh file path
func CheckIfSSHFileExists(sshFilePath string) error {
	if _, err := os.Stat(sshFilePath); os.IsNotExist(err) {
		return fmt.Errorf("ssh file not found: %s", err.Error())
	}
	return nil
}

// GetDefaultSSHFilePath returns the absolute path to ~/.ssh/id_rsa
func GetDefaultSSHFilePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	defaultSSHFilePath := currentUser.HomeDir + "/.ssh/id_rsa"
	return defaultSSHFilePath, nil
}
