package test

import (
	"os"
	"path/filepath"
	"testing"
)

// getConfigPath returns the path to the gclm-engine directory (configDir)
// The gclm-engine directory contains:
//   - gclm_engine_config.yaml
//   - workflows/*.yaml
//   - migrations/*.sql
func getConfigPath(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// If we're in test directory, go up one level to gclm-engine
	if filepath.Base(wd) == "test" {
		return filepath.Dir(wd)
	}

	// If we're in gclm-engine, use current directory
	if filepath.Base(wd) == "gclm-engine" {
		return wd
	}

	// Default: assume current directory is gclm-engine
	return wd
}

// getWorkflowsPath returns the path to the workflows directory
func getWorkflowsPath(t *testing.T) string {
	return filepath.Join(getConfigPath(t), "workflows")
}
