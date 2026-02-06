package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gclm/gclm-flow/gclm-engine/internal/assets"
	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/internal/repository"
	"github.com/gclm/gclm-flow/gclm-engine/internal/service"
	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/spf13/cobra"
)

// CLI represents the command-line interface
type CLI struct {
	rootCmd        *cobra.Command
	db             *db.Database
	parser         *workflow.Parser
	taskRepo       domain.TaskRepository
	workflowRepo   domain.WorkflowRepository
	workflowLoader domain.WorkflowLoader
	taskSvc        *service.TaskService
	workflowSvc    *service.WorkflowService
	configDir      string
	workflowsDir   string
}

// New creates a new CLI instance
func New(configDir string) (*CLI, error) {
	// Check if initialization is needed
	needsInit := checkNeedsInit(configDir)
	if needsInit {
		// Silent auto-init for first-time setup
		if err := autoInitialize(configDir); err != nil {
			return nil, fmt.Errorf("auto-initialization failed: %w", err)
		}
	}

	// Initialize database
	database, err := db.New(db.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Workflows directory
	workflowsDir := filepath.Join(configDir, "workflows")

	// Initialize builtin workflows from YAML files
	if err := database.InitWorkflows(workflowsDir); err != nil {
		return nil, fmt.Errorf("failed to initialize workflows: %w", err)
	}

	// Initialize workflow parser (still needed for YAML loading operations)
	parser := workflow.NewParser(workflowsDir)

	// Initialize repositories
	repo := db.NewRepository(database)
	workflowRepoDB := db.NewWorkflowRepository(database)

	// Create domain adapters
	taskRepo := repository.NewTaskRepository(repo)
	workflowRepo := repository.NewWorkflowRepository(workflowRepoDB)
	workflowLoader := repository.NewWorkflowLoader(parser)

	// Initialize services
	taskSvc := service.NewTaskService(taskRepo, workflowLoader)
	workflowSvc := service.NewWorkflowService(workflowRepo, workflowLoader)

	cli := &CLI{
		db:             database,
		parser:         parser,
		taskRepo:       taskRepo,
		workflowRepo:   workflowRepo,
		workflowLoader: workflowLoader,
		taskSvc:        taskSvc,
		workflowSvc:    workflowSvc,
		configDir:      configDir,
		workflowsDir:   workflowsDir,
	}

	cli.rootCmd = cli.createRootCommand()

	return cli, nil
}

// checkNeedsInit 检查是否需要初始化
func checkNeedsInit(configDir string) bool {
	configFile := filepath.Join(configDir, "gclm_engine_config.yaml")
	workflowsDir := filepath.Join(configDir, "workflows")

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return true
	}

	// Check if workflows directory exists and has files
	entries, err := os.ReadDir(workflowsDir)
	if err != nil || len(entries) == 0 {
		return true
	}

	return false
}

// autoInitialize 静默自动初始化
func autoInitialize(configDir string) error {
	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	// Create workflows directory
	workflowsDir := filepath.Join(configDir, "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		return fmt.Errorf("create workflows directory: %w", err)
	}

	// Extract embedded workflows
	if _, err := assets.ExportBuiltinWorkflows(workflowsDir, true); err != nil {
		return fmt.Errorf("extract workflows: %w", err)
	}

	// Create default config
	if _, err := assets.ExportDefaultConfig(configDir, true); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

// Run executes the CLI
func (c *CLI) Run() error {
	return c.rootCmd.Execute()
}

// Close closes the CLI and releases resources
func (c *CLI) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Init creates and initializes a new CLI instance
func Init(configDir string) (*CLI, error) {
	return New(configDir)
}

// createRootCommand creates the root command
func (c *CLI) createRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gclm-engine",
		Short: "GCLM Flow Engine - 智能工作流引擎",
		Long: `GCLM Flow Engine 是一个基于 Go 的智能工作流引擎，
支持通过 YAML 配置自定义工作流，并协调多个 Agent 协作完成任务。`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	// Add subcommands
	rootCmd.AddCommand(c.createInitCommand())
	rootCmd.AddCommand(c.createTaskCommand())
	rootCmd.AddCommand(c.createWorkflowCommand())
	rootCmd.AddCommand(c.createServeCommand())
	rootCmd.AddCommand(c.createVersionCommand())

	return rootCmd
}
