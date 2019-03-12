package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func setupTestRepo(t *testing.T) (*git.Repository, string, func()) {
	t.Log("init tmp git repo")

	tmpFolder := "/tmp/github.com/meinto/semver/"
	os.MkdirAll(tmpFolder, os.ModePerm)
	r, err := git.PlainInit(tmpFolder, false)
	if err != nil {
		t.Error(err)
	}

	return r, tmpFolder, func() {
		os.RemoveAll("/tmp/github.com/")
		t.Log("finish")
	}
}

func TestCheckIfRepoIsClean(t *testing.T) {
	err := CheckIfRepoIsClean("/tmp")
	t.Log(err)
	if err == nil {
		t.Error("should throw error because the repo is not initialized")
	}

	repo, repoPath, teardown := setupTestRepo(t)
	defer teardown()

	filename := filepath.Join(repoPath, "example-git-file")
	ioutil.WriteFile(filename, []byte("hello world!"), 0644)

	err = CheckIfRepoIsClean(repoPath)
	t.Log(err)
	if err == nil {
		t.Error("should throw error because there are uncomitted files")
	}

	w, _ := repo.Worktree()
	w.Add(".")
	w.Commit("test commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@email.com",
			When:  time.Now(),
		},
	})

	err = CheckIfRepoIsClean(repoPath)
	if err != nil {
		t.Error(err)
	}
}
