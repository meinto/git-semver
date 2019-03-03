package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/pkg/errors"
)

// WriteJSONVersionFile writes given jsonContent back to json file
func WriteJSONVersionFile(jsonContent map[string]interface{}, versionFile string) error {
	newJSONContent, _ := json.MarshalIndent(jsonContent, "", "  ")
	err := ioutil.WriteFile(versionFile, newJSONContent, 0644)
	return errors.Wrap(err, fmt.Sprintf("error writing %s", versionFile))
}

// WriteRAWVersionFile writes given jsonContent back to json file
func WriteRAWVersionFile(version, versionFile string) error {
	err := ioutil.WriteFile(versionFile, []byte(version), 0644)
	return errors.Wrap(err, fmt.Sprintf("error writing %s", versionFile))
}

// CheckIfSSHFileExists validate the existence of the given ssh file path
func CheckIfSSHFileExists(sshFilePath string) error {
	_, err := os.Stat(sshFilePath)
	return errors.Wrap(err, "ssh file not found")
}

// GetDefaultSSHFilePath returns the absolute path to ~/.ssh/id_rsa
func GetDefaultSSHFilePath() (string, error) {
	currentUser, err := user.Current()
	defaultSSHFilePath := currentUser.HomeDir + "/.ssh/id_rsa"
	return defaultSSHFilePath, errors.Wrap(err, "error getting default ssh file path")
}
