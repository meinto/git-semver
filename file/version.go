package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type VersionFileService interface {
	WriteVersionFile(version string) error
	WriteVersionJSONFile(version string) error
	ReadVersionFromFile() (string, error)
	ReadVersionFromJSONFile() (string, error)
	readJSONFile() (map[string]interface{}, error)
}

type versionFileService struct {
	filepath string
}

func NewVersionFileService(filepath string) VersionFileService {
	return &versionFileService{filepath}
}

func (s *versionFileService) WriteVersionFile(version string) error {
	err := ioutil.WriteFile(s.filepath, []byte(version), 0644)
	return errors.Wrap(err, fmt.Sprintf("error writing %s", s.filepath))
}

func (s *versionFileService) WriteVersionJSONFile(version string) error {
	jsonContent, err := s.readJSONFile()
	if err != nil {
		return err
	}
	jsonContent["version"] = version
	newJSONContent, _ := json.MarshalIndent(jsonContent, "", "  ")
	err = ioutil.WriteFile(s.filepath, newJSONContent, 0644)
	return errors.Wrap(err, fmt.Sprintf("error writing %s", s.filepath))
}

func (s *versionFileService) ReadVersionFromFile() (string, error) {
	versionFile, err := os.Open(s.filepath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("cannot read %s", s.filepath))
	}
	defer versionFile.Close()
	byteValue, err := ioutil.ReadAll(versionFile)
	return string(byteValue), errors.Wrap(err, "cannot read json")
}

func (s *versionFileService) ReadVersionFromJSONFile() (string, error) {
	jsonContent, err := s.readJSONFile()
	if err != nil {
		return "", err
	}
	version, ok := jsonContent["version"]
	if !ok {
		return "", errors.New("version property not set")
	}
	return version.(string), nil
}

func (s *versionFileService) readJSONFile() (map[string]interface{}, error) {
	versionFile, err := os.Open(s.filepath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("cannot read %s", s.filepath))
	}
	defer versionFile.Close()
	byteValue, err := ioutil.ReadAll(versionFile)
	var jsonContent = make(map[string]interface{})
	err = json.Unmarshal(byteValue, &jsonContent)
	return jsonContent, errors.Wrap(err, "cannot read json")
}
