package file

import (
	"os"
	"os/user"

	"github.com/pkg/errors"
)

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
