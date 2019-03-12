package internal

import (
	"github.com/pkg/errors"

	"github.com/meinto/git-semver/util"
)

func ValidateReadyForPushingChanges(repoPath, sshFilePath string, shouldPush bool) {
	if shouldPush {
		err := util.CheckIfRepoIsClean(repoPath)
		LogFatalOnErr(errors.Wrap(err, "repo not clean"))

		err = util.CheckIfSSHFileExists(sshFilePath)
		LogFatalOnErr(errors.Wrap(err, "ssh file doesn't exists"))
	}
}
