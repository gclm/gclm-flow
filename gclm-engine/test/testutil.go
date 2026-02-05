package test

import (
	"os"
	"path/filepath"
	"testing"
)

// getConfigPath returns the path to the workflows directory
func getConfigPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// If we're in test directory, go up two levels to project root
	if filepath.Base(wd) == "test" {
		return filepath.Join(filepath.Dir(filepath.Dir(wd)), "workflows")
	}

	// If we're in gclm-engine, go up one level
	if filepath.Base(wd) == "gclm-engine" {
		return filepath.Join(filepath.Dir(wd), "workflows")
	}

	return filepath.Join(wd, "workflows")
}
