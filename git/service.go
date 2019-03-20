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
	PushTag(name string) error
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
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	return strings.TrimSuffix(stdout.String(), "\n"), errors.Wrap(err, fmt.Sprintf("pkg(git) GitRepoPath(): %s", stderr.String()))
}

func (s service) IsRepoClean() (bool, error) {
	cmd := exec.Command(s.gitPath, "status", "-s")
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	return stdout.String() == "", errors.Wrap(err, fmt.Sprintf("pkg(git) IsRepoClean(): %s", stderr.String()))
}

func (s service) CreateTag(version string) error {
	cmd := exec.Command(s.gitPath, "tag", "-a", "v"+version, "-m", fmt.Sprintf("create new tag v%s", version))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) CreateTag(): %s", stderr.String()))
}

func (s service) Push() error {
	cmd := exec.Command(s.gitPath, "push", "--follow-tags")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) Push(): %s", stderr.String()))
}

func (s service) PushTag(version string) error {
	cmd := exec.Command(s.gitPath, "push", "origin", "v"+version)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) PushTag(): %s", stderr.String()))
}

func (s service) AddVersionChanges(filename string) error {
	repoPath, err := s.GitRepoPath()
	if err != nil {
		return err
	}
	filePath := repoPath + "/" + filename
	cmd := exec.Command(s.gitPath, "add", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) AddVersionChanges(): %s", stderr.String()))
}

func (s service) CommitVersionChanges(version string) error {
	cmd := exec.Command(s.gitPath, "commit", "-m", fmt.Sprintf("add changes for version %s", version))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) CommitVersionChanges(): %s", stderr.String()))
}
