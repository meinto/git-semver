package util

import (
	"errors"
	"strconv"
	"strings"
)

// NextVersion returns the next version number corresponding to given next-version-type
// next-version-type can be "major", "minor", "patch"
func NextVersion(currentVersion, nextVersionType string) (string, error) {
	numbers := strings.Split(currentVersion, ".")
	if len(numbers) != 3 {
		return "", errors.New("please provide version number in the following format: <major>.<minor>.<patch>")
	}

	switch nextVersionType {
	case "major":
		major, _ := strconv.Atoi(numbers[0])
		numbers[0] = strconv.Itoa(major + 1)
		numbers[1] = "0"
		numbers[2] = "0"
	case "minor":
		minor, _ := strconv.Atoi(numbers[1])
		numbers[1] = strconv.Itoa(minor + 1)
		numbers[2] = "0"
	case "patch":
		patch, _ := strconv.Atoi(numbers[2])
		numbers[2] = strconv.Itoa(patch + 1)
	}

	return strings.Join(numbers, "."), nil
}
