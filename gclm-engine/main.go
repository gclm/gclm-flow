package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gclm/gclm-flow/gclm-engine/internal/cli"
)

func main() {
	// Determine workflows directory
	workflowsDir := os.Getenv("GCLM_ENGINE_WORKFLOWS_DIR")
	if workflowsDir == "" {
		// Default to ~/.gclm-flow/workflows
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		workflowsDir = filepath.Join(homeDir, ".gclm-flow", "workflows")
	}

	// Initialize CLI
	c, err := cli.New(workflowsDir)
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
