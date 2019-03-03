package util

import (
	"encoding/json"
	"fmt"
	"log"

	semverUtil "github.com/meinto/git-semver/util"
	"github.com/pkg/errors"
)

func VersionFilePath(gitRepoPath, versionFile string) string {
	return gitRepoPath + "/" + versionFile
}

func GetJsonContent(fileContent []byte) map[string]interface{} {
	var jsonContent = make(map[string]interface{})
	err := json.Unmarshal(fileContent, &jsonContent)
	LogFatalOnErr(errors.Wrap(err, "cannot read json"))
	return jsonContent
}

func GetVersion(versionFileType string, fileContent []byte) string {
	switch versionFileType {
	case "json":
		jsonContent := GetJsonContent(fileContent)
		var version interface{}
		version, ok := jsonContent["version"]
		LogFatalIfNotOk(ok, "key version is not defined in map")
		return version.(string)
	case "raw":
		return string(fileContent)
	default:
		log.Fatal("unknown file type")
		return ""
	}
}

func WriteVersion(versionFileType, versionFileName, version string, fileContent []byte) error {
	var err error
	switch versionFileType {
	case "json":
		jsonContent := GetJsonContent(fileContent)
		jsonContent["version"] = version
		err = semverUtil.WriteJSONVersionFile(jsonContent, versionFileName)
		break
	case "raw":
		err = semverUtil.WriteRAWVersionFile(version, versionFileName)
	default:
		err = errors.New("unknown file type")
	}
	return err
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
