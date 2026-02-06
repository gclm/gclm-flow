package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createInitCommand creates the init command
func (c *CLI) createInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize gclm-engine configuration",
		Long: `Initialize gclm-engine by creating the configuration directory and
extracting embedded workflow files.

This command is usually run automatically on first use, but can be run
manually to reset the configuration.`,
		RunE: c.runInit,
	}

	cmd.Flags().Bool("silent", false, "Suppress output (for automatic init)")
	return cmd
}

// runInit executes the init command
func (c *CLI) runInit(cmd *cobra.Command, args []string) error {
	silent, _ := cmd.Flags().GetBool("silent")

	// Check if already initialized
	if !checkNeedsInit(c.configDir) {
		if !silent {
			fmt.Println("gclm-engine is already initialized")
			fmt.Printf("Config directory: %s\n", c.configDir)
		}
		return nil
	}

	// Run initialization
	if err := autoInitialize(c.configDir); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	if !silent {
		fmt.Println("gclm-engine initialized successfully")
		fmt.Printf("Config directory: %s\n", c.configDir)
		fmt.Printf("Workflows directory: %s\n", c.workflowsDir)
		fmt.Println("\nYou can now use 'gclm-engine task create' to create a new task")
	}

	return nil
}

// createVersionCommand creates version command
func (c *CLI) createVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gclm-engine v0.2.0")
		},
	}
}
