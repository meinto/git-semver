package util

import "fmt"

func noValidGitRepo(err error) error {
	return fmt.Errorf("this is no valid git repository: %s", err.Error())
}
