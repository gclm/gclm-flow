package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gclm/gclm-flow/gclm-engine/internal/cli"
	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
)

func main() {
	// Set embedded migrations for database initialization
	db.SetMigrationsFS(MigrationsFS())

	// Check command to decide whether to show migration logs
	// Only enable for "init" and "serve" commands for diagnostic purposes
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "init" || cmd == "serve" {
			os.Setenv("GCLM_MIGRATION_VERBOSE", "1")
		}
	}

	// Determine config directory
	configDir := os.Getenv("GCLM_ENGINE_CONFIG_DIR")
	if configDir == "" {
		// Default to ~/.gclm-flow
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		configDir = filepath.Join(homeDir, ".gclm-flow")
	}

	// Initialize CLI
	c, err := cli.New(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing CLI: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()

	// Run CLI
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
