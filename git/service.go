package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Service describes all actions which can performed with git
type Service interface {
	GitRepoPath() (string, error)
	IsRepoClean() (bool, error)
	CreateTag(version string) error
	Push() error
	AddVersionChanges(filename string) error
	CommitVersionChanges(version string) error
}

type service struct {
	gitPath string
	Service
}

func NewGitService(gitPath string) Service {
	return service{gitPath: gitPath}
}

func (s service) GitRepoPath() (string, error) {
	cmd := exec.Command(s.gitPath, "rev-parse", "--show-toplevel")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), err
}

func (s service) IsRepoClean() (bool, error) {
	cmd := exec.Command(s.gitPath, "status", "-s")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String() == "", err
}

func (s service) CreateTag(version string) error {
	cmd := exec.Command(s.gitPath, "tag", "-a", "v"+version, "-m", fmt.Sprintf("create new tag v%s", version))
	err := cmd.Run()
	return errors.Wrap(err, "error creating git tag")
}

func (s service) Push() error {
	cmd := exec.Command(s.gitPath, "push", "--follow-tags")
	err := cmd.Run()
	return errors.Wrap(err, "error creating git tag")
}

func (s service) AddVersionChanges(filename string) error {
	repoPath, err := s.GitRepoPath()
	if err != nil {
		return err
	}
	filePath := repoPath + "/" + filename
	cmd := exec.Command(s.gitPath, "add", filePath)
	err = cmd.Run()
	return errors.Wrap(err, "error adding version changes")
}

func (s service) CommitVersionChanges(version string) error {
	cmd := exec.Command(s.gitPath, "commit", "-m", fmt.Sprintf("add changes for version %s"))
	err := cmd.Run()
	return errors.Wrap(err, "error committing added changes")
}
