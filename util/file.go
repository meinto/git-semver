package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
