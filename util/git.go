package util

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// CheckIfRepoIsClean checks for uncomitted files
func CheckIfRepoIsClean(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return noValidGitRepo(err)
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}
	status, err := w.Status()
	if err != nil {
		return err
	}
	if !status.IsClean() {
		return errors.New("please commit all files before versioning")
	}
	return nil
}

// MakeGitTag tags the repository with given version number
// format of git tag: "v<version-number>"
func MakeGitTag(repoPath, version string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return noValidGitRepo(err)
	}

	headRef, err := r.Head()
	if err != nil {
		return err
	}

	tag := fmt.Sprintf("refs/tags/v%s", version)
	ref := plumbing.NewHashReference(plumbing.ReferenceName(tag), headRef.Hash())

	err = r.Storer.SetReference(ref)
	if err != nil {
		return err
	}
	return nil
}

// Push pushes all changes made by semver to defined bare repository
// The push includes file changes as well as git tags
func Push(repoPath, sshFilePath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return noValidGitRepo(err)
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshFilePath, "")
	if err != nil {
		return err
	}

	tags := config.RefSpec("refs/tags/*:refs/tags/*")
	heads := config.RefSpec("refs/heads/*:refs/heads/*")
	err = r.Push(&git.PushOptions{
		Auth:     sshAuth,
		RefSpecs: []config.RefSpec{tags, heads},
	})
	if err != nil {
		return err
	}
	return nil
}

// AddVersionChanges adds all changes made by semver during the versioning process
func AddVersionChanges(repoPath, configFile, version, author, email string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return noValidGitRepo(err)
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(configFile)
	if err != nil {
		return err
	}
	_, err = w.Commit("new version: "+version, &git.CommitOptions{
		Author: &object.Signature{
			Name:  author,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
