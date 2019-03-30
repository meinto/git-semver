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
	shell      string
	pathToRepo string
	Service
}

func NewGitService(pathToShell string) Service {
	return service{
		shell:      pathToShell,
		pathToRepo: "",
	}
}

func NewRepoPathGitService(pathToShell, pathToRepo string) Service {
	return service{
		shell:      pathToShell,
		pathToRepo: pathToRepo,
	}
}

func (s service) Command(cmd string) (*exec.Cmd, error) {
	repoPath, err := s.GitRepoPath()
	if err != nil {
		return nil, err
	}
	formattedCommand := fmt.Sprintf("cd %s && %s", repoPath, cmd)
	return exec.Command(s.shell, "-c", formattedCommand), nil
}

func (s service) GitRepoPath() (string, error) {
	if s.pathToRepo == "" {
		cmd, err := s.Command("git rev-parse --show-toplevel")
		if err != nil {
			return "", err
		}
		var stdout, stderr bytes.Buffer
		cmd.Stderr = &stderr
		cmd.Stdout = &stdout
		err = cmd.Run()
		return strings.TrimSuffix(stdout.String(), "\n"), errors.Wrap(err, fmt.Sprintf("pkg(git) GitRepoPath(): %s", stderr.String()))
	} else {
		return s.pathToRepo, nil
	}
}

func (s service) IsRepoClean() (bool, error) {
	cmd, err := s.Command("git status -s")
	if err != nil {
		return false, err
	}
	var stdout, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	return stdout.String() == "", errors.Wrap(err, fmt.Sprintf("pkg(git) IsRepoClean(): %s", stderr.String()))
}

func (s service) CreateTag(version string) error {
	cmd, err := s.Command(fmt.Sprintf("git tag -a v%s -m 'create new tag v%s'", version, version))
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) CreateTag(): %s", stderr.String()))
}

func (s service) Push() error {
	cmd, err := s.Command("git push --follow-tags")
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) Push(): %s", stderr.String()))
}

func (s service) PushTag(version string) error {
	cmd, err := s.Command(fmt.Sprintf("git push origin v%s", version))
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) PushTag(): %s", stderr.String()))
}

func (s service) AddVersionChanges(filename string) error {
	repoPath, err := s.GitRepoPath()
	if err != nil {
		return err
	}
	filePath := repoPath + "/" + filename
	cmd, err := s.Command(fmt.Sprintf("git add %s", filePath))
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) AddVersionChanges(): %s", stderr.String()))
}

func (s service) CommitVersionChanges(version string) error {
	cmd, err := s.Command(fmt.Sprintf("git commit -m 'add changes for version %s'", version))
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	return errors.Wrap(err, fmt.Sprintf("pkg(git) CommitVersionChanges(): %s", stderr.String()))
}
