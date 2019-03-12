package util

import "testing"

func TestNextVersion(t *testing.T) {
	var testCases = []struct {
		currentVersion  string
		nextVersionType string
		nextVersion     string
	}{
		{"1.2.3", "major", "2.0.0"},
		{"1.2.3", "minor", "1.3.0"},
		{"1.2.3", "patch", "1.2.4"},
	}

	for _, testCase := range testCases {
		nextVersion, err := NextVersion(testCase.currentVersion, testCase.nextVersionType)
		if err != nil {
			t.Error(err)
		}
		if nextVersion != testCase.nextVersion {
			t.Errorf("expected: %s, got: %s", testCase.nextVersion, nextVersion)
		}
	}
}
