package semver

import (
	"errors"
	"strconv"
	"strings"
)

type VersionType string

const (
	PATCH VersionType = "patch"
	MINOR VersionType = "minor"
	MAJOR VersionType = "major"
)

type Version interface {
	Get(versionType string) (string, error)
	SetNext(versionType string) (string, error)
}

type version struct {
	current string
}

func NewVersion(v string) (Version, error) {
	if !IsValidVersion(v) {
		return nil, errors.New("version invalid - please provide version number in the following format: <major>.<minor>.<patch>")
	}
	return &version{v}, nil
}

func (v *version) Get(vt string) (string, error) {
	versionType := VersionType(vt)

	numbers := strings.Split(v.current, ".")

	for i, n := range numbers {
		numbers[i] = strings.TrimSuffix(n, "\n")
	}

	var err error
	switch versionType {
	case MAJOR:
		major, e := strconv.Atoi(numbers[0])
		err = e
		numbers[0] = strconv.Itoa(major + 1)
		numbers[1] = "0"
		numbers[2] = "0"
	case MINOR:
		minor, e := strconv.Atoi(numbers[1])
		err = e
		numbers[1] = strconv.Itoa(minor + 1)
		numbers[2] = "0"
	case PATCH:
		patch, e := strconv.Atoi(numbers[2])
		err = e
		numbers[2] = strconv.Itoa(patch + 1)
	}

	return strings.Join(numbers, "."), err
}

func (v *version) SetNext(vt string) (string, error) {
	versionType := VersionType(vt)
	if !IsValidNextVersionType(versionType) {
		return "", errors.New("invalid next version type")
	}

	next, _ := v.Get(vt)
	v.current = next

	return v.current, nil
}
