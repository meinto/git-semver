package util

import (
	"github.com/pkg/errors"

	semverUtil "github.com/meinto/git-semver/util"
)

func ValidateReadyForPushingChanges(repoPath, sshFilePath string, shouldPush bool) {
	if shouldPush {
		err := semverUtil.CheckIfRepoIsClean(repoPath)
		LogFatalOnErr(errors.Wrap(err, "repo not clean"))

		err = semverUtil.CheckIfSSHFileExists(sshFilePath)
		LogFatalOnErr(errors.Wrap(err, "ssh file doesn't exists"))
	}
}
