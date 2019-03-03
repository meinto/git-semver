package util

import (
	"github.com/pkg/errors"
)

func noValidGitRepo(err error) error {
	return errors.Wrap(err, "this is no valid git repository")
}
