package util

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func VersionFilePath(gitRepoPath, versionFile string) string {
	return gitRepoPath + "/" + versionFile
}

func GetVersion(versionFileType string, fileContent []byte) string {
	switch versionFileType {
	default:
		log.Fatal("unknown file type")
		return ""
	case "json":
		var jsonContent = make(map[string]interface{})
		err := json.Unmarshal(fileContent, &jsonContent)
		LogFatalOnErr(errors.Wrap(err, "cannot read json"))
		var version interface{}
		version, ok := jsonContent["version"]
		LogFatalIfNotOk(ok, "key version is not defined in map")
		return version.(string)
	case "raw":
		return string(fileContent)
	}
}

func ValidateNextVersionType(nextVersionType string) {
	if nextVersionType != "major" && nextVersionType != "minor" && nextVersionType != "patch" {
		log.Fatal("please choose one of these values: major, minor, patch")
	}
}

func PrintNextVersion(nextVersionType, nextVersion string, raw bool) {
	printVersion(nextVersion, fmt.Sprintf("Next %s version: ", nextVersionType), raw)
}

func PrintCurrentVersion(currentVersion string, raw bool) {
	printVersion(currentVersion, "Current version: ", raw)
}

func printVersion(nextVersion, message string, raw bool) {
	if !raw {
		fmt.Print(message)
	}
	fmt.Println(nextVersion)
}
