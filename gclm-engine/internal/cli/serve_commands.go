package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gclm/gclm-flow/gclm-engine/internal/api"
	"github.com/gclm/gclm-flow/gclm-engine/internal/api/websocket"
	"github.com/gclm/gclm-flow/gclm-engine/internal/assets"
	"github.com/gclm/gclm-flow/gclm-engine/internal/logger"
	"github.com/spf13/cobra"
)

// createServeCommand creates the serve command
func (c *CLI) createServeCommand() *cobra.Command {
	var port int
	var daemon bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP API server",
		Long: `Start the HTTP API server for gclm-engine.

The server provides:
- REST API for task and workflow management
- WebSocket support for real-time task updates
- Web UI at http://localhost:PORT (with embedded static files)

Examples:
  # Start server in foreground (default port 9988)
  gclm-engine serve

  # Start server on custom port
  gclm-engine serve --port 8888

  # Start server in background (daemon mode)
  gclm-engine serve --daemon

  # Stop background server
  gclm-engine serve --stop`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle --stop flag
			stop, _ := cmd.Flags().GetBool("stop")
			if stop {
				return c.stopDaemon()
			}

			// Handle internal flag (for daemon child process)
			internal, _ := cmd.Flags().GetBool("foreground-internal")
			if internal {
				// This is the child process, run in foreground mode
				return c.runServe(port, true)
			}

			// Handle --daemon flag: fork new process
			if daemon {
				return c.startDaemon(port)
			}

			// Normal foreground mode
			return c.runServe(port, false)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 9988, "Port to listen on (default: 9988)")
	cmd.Flags().BoolVarP(&daemon, "daemon", "d", false, "Run as daemon (background)")
	cmd.Flags().Bool("stop", false, "Stop background server")
	cmd.Flags().Bool("foreground-internal", false, "Internal flag for daemon child process (hidden)")
	_ = cmd.Flags().MarkHidden("foreground-internal")

	return cmd
}

// getPIDFilePath returns the PID file path
func getPIDFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".gclm-flow", "gclm-engine.pid")
}

// stopDaemon stops the background server
func (c *CLI) stopDaemon() error {
	pidFile := getPIDFilePath()

	// Read PID file
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no running server found (PID file does not exist)")
		}
		return fmt.Errorf("failed to read PID file: %w", err)
	}

	var pid int
	_, err = fmt.Sscanf(string(pidData), "%d", &pid)
	if err != nil {
		return fmt.Errorf("invalid PID file: %w", err)
	}

	// Check if process is running
	process, err := os.FindProcess(pid)
	if err != nil {
		os.Remove(pidFile)
		return fmt.Errorf("process not found: %w", err)
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		os.Remove(pidFile)
		return fmt.Errorf("failed to stop process: %w", err)
	}

	// Remove PID file
	os.Remove(pidFile)

	fmt.Printf("Server stopped (PID: %d)\n", pid)
	return nil
}

// startDaemon forks a new process for daemon mode
func (c *CLI) startDaemon(port int) error {
	pidFile := getPIDFilePath()

	// Check if already running
	if _, err := os.Stat(pidFile); err == nil {
		pidData, _ := os.ReadFile(pidFile)
		return fmt.Errorf("server already running (PID: %s). Use --stop to stop it first", string(pidData))
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(pidFile), 0755); err != nil {
		return fmt.Errorf("failed to create PID directory: %w", err)
	}

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Log file path
	logFile := filepath.Join(filepath.Dir(pidFile), "gclm-engine.log")
	logFH, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Create command for child process
	cmd := exec.Command(execPath, "serve", "--port", fmt.Sprintf("%d", port), "--foreground-internal")
	cmd.Dir = filepath.Dir(execPath)

	// Redirect stdout and stderr to log file
	cmd.Stdout = logFH
	cmd.Stderr = logFH

	// Set sysprocattr to detach from terminal (create new session)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Start the daemon process
	if err := cmd.Start(); err != nil {
		logFH.Close()
		return fmt.Errorf("failed to start daemon: %w", err)
	}

	// Write PID file
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		cmd.Process.Kill()
		logFH.Close()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	// Close log file handle (daemon has its own copy)
	logFH.Close()

	// Print startup info and exit (parent process exits here)
	fmt.Printf("Server started in background\n")
	fmt.Printf("  PID: %d\n", cmd.Process.Pid)
	fmt.Printf("  Port: %d\n", port)
	fmt.Printf("  URL: http://localhost:%d\n", port)
	fmt.Printf("  Log: %s\n", logFile)
	fmt.Printf("\nStop with: gclm-engine serve --stop\n")

	return nil
}

// runServe executes the serve command
func (c *CLI) runServe(port int, isDaemon bool) error {
	addr := fmt.Sprintf(":%d", port)

	// For daemon mode, ensure output is redirected to log
	if isDaemon {
		logFile := filepath.Join(filepath.Dir(getPIDFilePath()), "gclm-engine.log")
		logFH, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			os.Stdout = logFH
			os.Stderr = logFH
		}
	}

	// Get embedded web filesystem
	var webFS fs.FS = assets.WebFS()
	if webFS != nil {
		logger.Info().Msg("Using embedded web files")
	} else {
		logger.Warn().Msg("Web files not embedded, web UI will not be available")
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create WebSocket hub
	wsHub := websocket.NewHub(ctx)

	// Create HTTP server with embedded web files
	server := api.NewServer(addr, c.taskSvc, c.workflowSvc, wsHub, webFS)

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		if err := server.Start(ctx); err != nil {
			serverErrChan <- err
		}
	}()

	// Log startup (goes to log file in daemon mode)
	logger.Info().Str("addr", addr).Msg("Server starting")
	logger.Info().Str("url", fmt.Sprintf("http://localhost:%d", port)).Msg("Web UI available")

	// Wait for signal or error
	select {
	case <-sigChan:
		logger.Info().Msg("Received interrupt signal, shutting down...")
		cancel()
	case err := <-serverErrChan:
		// Clean up PID file on error
		if isDaemon {
			os.Remove(getPIDFilePath())
		}
		return err
	}

	// Wait a moment for cleanup
	logger.Info().Msg("Server stopped")

	// Clean up PID file on graceful shutdown
	if isDaemon {
		os.Remove(getPIDFilePath())
	}

	return nil
}
