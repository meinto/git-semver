package semver

import "strings"

func IsValidNextVersionType(nvt VersionType) bool {
	if nvt == MAJOR || nvt == MINOR || nvt == PATCH {
		return true
	}
	return false
}

func IsValidVersion(v string) bool {
	numbers := strings.Split(v, ".")
	if len(numbers) != 3 {
		return false
	}
	return true
}
