package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func CheckIfRepoIsClean(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
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

func MakeGitTag(repoPath, version string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
	} else {
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
	}
	return nil
}

func Push(repoPath, sshFilePath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
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

func AddVersionChanges(repoPath, configFile, version, author, email string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Println("this is no valid git repository")
		return err
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
