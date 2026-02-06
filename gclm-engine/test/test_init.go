package test

import (
	"os"
	"path/filepath"
	"testing"
)

func init() {
	// Change to gclm-engine directory so migrations can be found
	// When running `go test ./test/...`, the working directory is set to test/
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// If we're in test directory, go up one level to gclm-engine
	if filepath.Base(wd) == "test" {
		if err := os.Chdir(".."); err != nil {
			panic(err)
		}
	}
}

// TestMain is the main entry point for tests
func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
