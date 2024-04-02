package testhelper

import (
	"os"
	"testing"
)

func LoadJsonGoldenFile(t *testing.T, path string) []byte {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Read test data file
	jsonData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to get json test data %q: %v", path, err)
	}

	return jsonData
}
