package repository

import (
	"os"
	"path/filepath"
)

func init() {
	// Change to gclm-engine directory so migrations can be found
	// When running `go test ./internal/repository/...`, the working directory is set to internal/repository/
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// If we're in internal/repository, go up two levels to gclm-engine
	if filepath.Base(wd) == "repository" {
		if err := os.Chdir("../.."); err != nil {
			panic(err)
		}
	}
	// If we're in internal, go up one level
	if filepath.Base(wd) == "internal" {
		if err := os.Chdir(".."); err != nil {
			panic(err)
		}
	}
}
