package cli

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/internal/assets"
	"github.com/gclm/gclm-flow/gclm-engine/internal/config"
	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/errors"
	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/internal/service"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CLI represents the command-line interface
type CLI struct {
	rootCmd   *cobra.Command
	db        *db.Database
	parser    *workflow.Parser
	repo      *db.Repository
	taskSvc   *service.TaskService
	configDir string
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

	// Initialize builtin workflows from YAML files
	workflowsDir := filepath.Join(configDir, "workflows")
	if err := database.InitWorkflows(workflowsDir); err != nil {
		return nil, fmt.Errorf("failed to initialize workflows: %w", err)
	}

	// Initialize pipeline parser (still needed for YAML loading operations)
	parser := workflow.NewParser(configDir)

	// Initialize repository
	repo := db.NewRepository(database)

	// Initialize task service
	taskSvc := service.NewTaskService(repo, parser)

	cli := &CLI{
		db:        database,
		parser:    parser,
		repo:      repo,
		taskSvc:   taskSvc,
		configDir: configDir,
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
	// Export default config
	if _, err := assets.ExportDefaultConfig(configDir, false); err != nil {
		return err
	}

	// Export builtin workflows
	workflowsDir := filepath.Join(configDir, "workflows")
	if _, err := assets.ExportBuiltinWorkflows(workflowsDir, false); err != nil {
		return err
	}

	// Note: Database initialization happens in db.New() via migrations

	return nil
}

// createRootCommand creates the root command
func (c *CLI) createRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "gclm-engine",
		Short: "gclm-flow workflow engine",
		Long:  "Go-based workflow engine for AI-driven development tasks",
	}

	// 全局 JSON 输出选项
	root.PersistentFlags().Bool("json", false, "Output in JSON format")
	root.PersistentFlags().Bool("pretty", true, "Pretty print JSON output")

	// Add subcommands
	root.AddCommand(c.createInitCommand())
	root.AddCommand(c.createTaskCommand())
	root.AddCommand(c.createPipelineCommand())
	root.AddCommand(c.createWorkflowCommand())
	root.AddCommand(c.createVersionCommand())

	return root
}

// createInitCommand creates the init command (top-level)
func (c *CLI) createInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize gclm-engine configuration",
		Long:  "Create default configuration and workflow files in ~/.gclm-flow/\n\n" +
			"If --force is specified, existing files will be overwritten.\n" +
			"If --silent is specified, no output will be printed (useful for automatic init).",
		RunE: c.runWorkflowInit,
	}
	cmd.Flags().Bool("force", false, "Overwrite existing files")
	cmd.Flags().Bool("silent", false, "Suppress output (for automatic init)")
	return cmd
}

// createTaskCommand creates task management commands
func (c *CLI) createTaskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task management commands",
	}

	// task create
	createCmd := &cobra.Command{
		Use:   "create <prompt>",
		Short: "Create a new task",
		Long:  "Create a new task. Workflow type is auto-detected from prompt keywords.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCreate,
	}
	createCmd.Flags().String("workflow-type", "", "Workflow type (CODE_SIMPLE, CODE_COMPLEX, DOCUMENT). Auto-detected if not specified.")
	createCmd.Flags().String("pipeline", "", "Pipeline name (overrides workflow-type)")

	// task get
	getCmd := &cobra.Command{
		Use:   "get <task-id>",
		Short: "Get task details",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskGet,
	}

	// task list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE:  c.runTaskList,
	}
	listCmd.Flags().String("status", "", "Filter by status")
	listCmd.Flags().Int("limit", 20, "Maximum number of tasks to show")

	// task current - 获取当前应该执行的阶段
	currentCmd := &cobra.Command{
		Use:   "current <task-id>",
		Short: "Get current phase to execute",
		Long:  "Get the next pending phase that should be executed. Used by skills to determine what to do next.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCurrent,
	}

	// task plan - 获取完整执行计划
	planCmd := &cobra.Command{
		Use:   "plan <task-id>",
		Short: "Get execution plan",
		Long:  "Get the complete execution plan with all phases and dependencies. Used by skills to understand the workflow.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskPlan,
	}

	// task update - 更新阶段状态
	updateCmd := &cobra.Command{
		Use:   "update <task-id> <phase-id> <status>",
		Short: "Update phase status",
		Long:  "Update phase status. Used by skills to report completion or failure.",
		Args:  cobra.ExactArgs(3),
		RunE:  c.runTaskUpdate,
	}
	updateCmd.Flags().String("output", "", "Phase output (for completed status)")
	updateCmd.Flags().String("error", "", "Error message (for failed status)")

	// task complete - 完成阶段
	completeCmd := &cobra.Command{
		Use:   "complete <task-id> <phase-id>",
		Short: "Mark phase as completed",
		Long:  "Mark a phase as completed with output. Shortcut for 'task update ... completed'",
		Args:  cobra.ExactArgs(2),
		RunE:  c.runTaskComplete,
	}
	completeCmd.Flags().String("output", "", "Phase output")

	// task fail - 标记阶段失败
	failCmd := &cobra.Command{
		Use:   "fail <task-id> <phase-id>",
		Short: "Mark phase as failed",
		Long:  "Mark a phase as failed with error message. Shortcut for 'task update ... failed'",
		Args:  cobra.ExactArgs(2),
		RunE:  c.runTaskFail,
	}
	failCmd.Flags().String("error", "", "Error message")

	// task phases
	phasesCmd := &cobra.Command{
		Use:   "phases <task-id>",
		Short: "Show task phases",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskPhases,
	}

	// task events
	eventsCmd := &cobra.Command{
		Use:   "events <task-id>",
		Short: "Show task events",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskEvents,
	}
	eventsCmd.Flags().Int("limit", 50, "Maximum number of events to show")

	// task export - 导出状态文件
	exportCmd := &cobra.Command{
		Use:   "export <task-id> <output-file>",
		Short: "Export task state to file",
		Long:  "Export task state to a markdown file with YAML frontmatter. Compatible with existing skills.",
		Args:  cobra.ExactArgs(2),
		RunE:  c.runTaskExport,
	}

	// task pause/resume/cancel
	pauseCmd := &cobra.Command{
		Use:   "pause <task-id>",
		Short: "Pause a task",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskPause,
	}

	resumeCmd := &cobra.Command{
		Use:   "resume <task-id>",
		Short: "Resume a paused task",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskResume,
	}

	cancelCmd := &cobra.Command{
		Use:   "cancel <task-id>",
		Short: "Cancel a task",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCancel,
	}

	cmd.AddCommand(createCmd, getCmd, listCmd, currentCmd, planCmd, updateCmd,
		completeCmd, failCmd, phasesCmd, eventsCmd, exportCmd,
		pauseCmd, resumeCmd, cancelCmd)

	return cmd
}

// createPipelineCommand creates pipeline management commands
func (c *CLI) createPipelineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Pipeline management commands",
	}

	// pipeline list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available pipelines",
		RunE:  c.runPipelineList,
	}

	// pipeline get
	getCmd := &cobra.Command{
		Use:   "get <pipeline-name>",
		Short: "Get pipeline details",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runPipelineGet,
	}

	// pipeline recommend
	recommendCmd := &cobra.Command{
		Use:   "recommend <prompt>",
		Short: "Recommend pipeline based on prompt",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runPipelineRecommend,
	}

	cmd.AddCommand(listCmd, getCmd, recommendCmd)

	return cmd
}

// createWorkflowCommand creates workflow commands
func (c *CLI) createWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow management commands",
	}

	// workflow start - 一键开始工作流
	startCmd := &cobra.Command{
		Use:   "start <prompt>",
		Short: "Start a new workflow",
		Long:  "Create a new task and return the first phase to execute.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runWorkflowStart,
	}
	startCmd.Flags().String("workflow", "", "Use specific workflow (auto-detected if not specified)")

	// workflow next - 获取下一步
	nextCmd := &cobra.Command{
		Use:   "next <task-id>",
		Short: "Get next phase to execute",
		Long:  "Get the next pending phase.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCurrent,
	}

	// workflow validate - 验证工作流配置
	validateCmd := &cobra.Command{
		Use:   "validate <yaml-file>",
		Short: "Validate a workflow configuration",
		Long:  "Validate a workflow YAML file. Works with any file path.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runWorkflowValidate,
	}

	// workflow install - 安装工作流
	installCmd := &cobra.Command{
		Use:   "install <yaml-file>",
		Short: "Install a workflow configuration",
		Long:  "Install a workflow YAML file to gclm-engine.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runWorkflowInstall,
	}
	installCmd.Flags().String("name", "", "Custom workflow name")

	// workflow uninstall - 卸载工作流
	uninstallCmd := &cobra.Command{
		Use:   "uninstall <workflow-name>",
		Short: "Uninstall a workflow",
		Long:  "Remove a workflow from gclm-engine.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runWorkflowUninstall,
	}

	// workflow list - 列出工作流
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all workflows",
		RunE: c.runWorkflowList,
	}

	// workflow export - 导出工作流
	exportCmd := &cobra.Command{
		Use:   "export <workflow-name> [output-file]",
		Short: "Export a workflow",
		Long:  "Export a workflow to a YAML file.",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  c.runWorkflowExport,
	}

	// workflow info - 显示工作流详情
	infoCmd := &cobra.Command{
		Use:   "info <workflow-name>",
		Short: "Show workflow details",
		Long:  "Display detailed information about a workflow.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runWorkflowInfo,
	}

	// workflow sync - 同步工作流到数据库
	syncCmd := &cobra.Command{
		Use:   "sync [yaml-file]",
		Short: "Sync workflow YAML files to database",
		Long: "Sync workflow YAML files (draft) to database (production).\n\n" +
			"Arguments:\n" +
			"  yaml-file    Path to YAML file (relative or absolute)\n" +
			"               If omitted, sync all workflows from workflows/ directory\n\n" +
			"Examples:\n" +
			"  gclm-engine workflow sync                           # Sync all\n" +
			"  gclm-engine workflow sync workflows/feat.yaml      # Sync specific file\n" +
			"  gclm-engine workflow sync ../custom/my_workflow.yaml # Sync from custom path\n\n" +
			"YAML files are treated as draft data, database stores production data.\n" +
			"Modifying YAML files does not affect the running system until sync is executed.",
		RunE: c.runWorkflowSync,
	}
	syncCmd.Flags().Bool("force", false, "Force sync even if validation fails")

	cmd.AddCommand(startCmd, nextCmd, validateCmd, installCmd, uninstallCmd, listCmd, exportCmd, infoCmd, syncCmd)

	return cmd
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

// ============================================================================
// Task commands
// ============================================================================

func (c *CLI) runTaskCreate(cmd *cobra.Command, args []string) error {
	prompt := args[0]

	// Get flags
	workflowType, _ := cmd.Flags().GetString("workflow-type")
	pipelineName, _ := cmd.Flags().GetString("pipeline")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()

	// workflow-type 或 pipeline 必需指定其一
	if workflowType == "" && pipelineName == "" {
		return fmt.Errorf("either --workflow-type or --pipeline is required")
	}

	var task *types.Task
	var err error

	// 使用 TaskService 创建任务
	if pipelineName != "" {
		// 直接使用 pipeline 名称
		task, err = c.taskSvc.CreateTask(ctx, prompt, pipelineName)
	} else {
		// 使用 workflow_type
		task, err = c.taskSvc.CreateTask(ctx, prompt, workflowType)
	}

	if err != nil {
		c.printFriendlyError(err)
		return err
	}

	// 输出结果
	output := map[string]interface{}{
		"task_id":       task.ID,
		"status":        task.Status,
		"workflow_type": task.WorkflowType,
		"pipeline":      task.WorkflowID,
		"total_phases":  task.TotalPhases,
		"current_phase": task.CurrentPhase,
		"message":       "Task created successfully",
	}

	c.printOutput(output, jsonOutput)
	return nil
}

func (c *CLI) runTaskGet(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	status, err := c.taskSvc.GetTaskStatus(ctx, taskID)
	if err != nil {
		c.printFriendlyError(errors.TaskNotFound(taskID))
		return err
	}

	c.printOutput(status, jsonOutput)
	return nil
}

func (c *CLI) runTaskList(cmd *cobra.Command, args []string) error {
	statusStr, _ := cmd.Flags().GetString("status")
	limit, _ := cmd.Flags().GetInt("limit")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	var status *types.TaskStatus
	if statusStr != "" {
		s := types.TaskStatus(statusStr)
		status = &s
	}

	tasks, err := c.repo.ListTasks(status, limit)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	if jsonOutput {
		c.printOutput(tasks, true)
		return nil
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	// 文本格式输出
	for _, task := range tasks {
		statusIcon := " "
		switch task.Status {
		case types.TaskStatusCompleted:
			statusIcon = "✓"
		case types.TaskStatusRunning:
			statusIcon = "→"
		case types.TaskStatusFailed:
			statusIcon = "✗"
		case types.TaskStatusPaused:
			statusIcon = "⏸"
		}

		fmt.Printf("[%s] %s | %s | %s | Phase %d/%d\n",
			statusIcon,
			task.ID[:8],
			task.WorkflowType,
			task.Status,
			task.CurrentPhase,
			task.TotalPhases)

		if len(task.Prompt) > 0 {
			fmt.Printf("    %s\n", truncate(task.Prompt, 60))
		}
	}

	return nil
}

// runTaskCurrent 获取当前应该执行的阶段（skills 使用）
func (c *CLI) runTaskCurrent(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	phase, err := c.taskSvc.GetCurrentPhase(ctx, taskID)
	if err != nil {
		// 如果没有待执行的阶段，返回空
		if strings.Contains(err.Error(), "no pending phase") {
			c.printOutput(map[string]interface{}{
				"task_id": taskID,
				"phase":   nil,
				"message": "All phases completed or no pending phase",
			}, jsonOutput)
			return nil
		}
		return fmt.Errorf("failed to get current phase: %w", err)
	}

	// 加载流水线获取阶段详情
	task, _ := c.repo.GetTask(taskID)
	wf, _ := c.parser.LoadWorkflow(task.WorkflowID)
	node := findNode(wf, phase.PhaseName)

	output := map[string]interface{}{
		"task_id":      taskID,
		"phase_id":     phase.ID,
		"phase_name":   phase.PhaseName,
		"display_name": phase.DisplayName,
		"sequence":     phase.Sequence,
		"agent":        phase.AgentName,
		"model":        phase.ModelName,
		"status":       phase.Status,
	}

	if node != nil {
		output["required"] = node.Required
		output["timeout"] = node.Timeout
		output["dependencies"] = node.DependsOn
	}

	c.printOutput(output, jsonOutput)
	return nil
}

// runTaskPlan 获取完整执行计划
func (c *CLI) runTaskPlan(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	plan, err := c.taskSvc.GetExecutionPlan(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get execution plan: %w", err)
	}

	c.printOutput(plan, jsonOutput)
	return nil
}

// runTaskUpdate 更新阶段状态
func (c *CLI) runTaskUpdate(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	phaseID := args[1]
	status := args[2]

	output, _ := cmd.Flags().GetString("output")
	errorMsg, _ := cmd.Flags().GetString("error")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()

	switch types.PhaseStatus(status) {
	case types.PhaseStatusCompleted:
		if output == "" {
			return fmt.Errorf("output is required for completed status")
		}
		err := c.taskSvc.ReportPhaseOutput(ctx, taskID, phaseID, output)
		if err != nil {
			return fmt.Errorf("failed to report phase output: %w", err)
		}

	case types.PhaseStatusFailed:
		if errorMsg == "" {
			return fmt.Errorf("error message is required for failed status")
		}
		err := c.taskSvc.ReportPhaseError(ctx, taskID, phaseID, errorMsg)
		if err != nil {
			return fmt.Errorf("failed to report phase error: %w", err)
		}

	default:
		return fmt.Errorf("unsupported status: %s (use 'completed' or 'failed')", status)
	}

	result := map[string]interface{}{
		"task_id":  taskID,
		"phase_id": phaseID,
		"status":   status,
		"message":  "Phase updated successfully",
	}

	c.printOutput(result, jsonOutput)
	return nil
}

// runTaskComplete 完成阶段（快捷命令）
func (c *CLI) runTaskComplete(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	phaseID := args[1]
	output, _ := cmd.Flags().GetString("output")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if output == "" {
		return fmt.Errorf("output is required")
	}

	ctx := context.Background()
	err := c.taskSvc.ReportPhaseOutput(ctx, taskID, phaseID, output)
	if err != nil {
		return fmt.Errorf("failed to complete phase: %w", err)
	}

	result := map[string]interface{}{
		"task_id":  taskID,
		"phase_id": phaseID,
		"status":   "completed",
		"message":  "Phase completed successfully",
	}

	c.printOutput(result, jsonOutput)
	return nil
}

// runTaskFail 标记阶段失败（快捷命令）
func (c *CLI) runTaskFail(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	phaseID := args[1]
	errorMsg, _ := cmd.Flags().GetString("error")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if errorMsg == "" {
		errorMsg = "Phase failed"
	}

	ctx := context.Background()
	err := c.taskSvc.ReportPhaseError(ctx, taskID, phaseID, errorMsg)
	if err != nil {
		// 如果是必需阶段失败，这是预期行为
		result := map[string]interface{}{
			"task_id":  taskID,
			"phase_id": phaseID,
			"status":   "failed",
			"message":  errMsg(err),
		}
		c.printOutput(result, jsonOutput)
		return nil
	}

	result := map[string]interface{}{
		"task_id":  taskID,
		"phase_id": phaseID,
		"status":   "failed",
		"message":  "Phase marked as failed",
	}

	c.printOutput(result, jsonOutput)
	return nil
}

func (c *CLI) runTaskPhases(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	phases, err := c.repo.GetPhasesByTask(taskID)
	if err != nil {
		return fmt.Errorf("failed to get phases: %w", err)
	}

	if jsonOutput {
		c.printOutput(phases, true)
		return nil
	}

	// 文本格式输出
	for _, phase := range phases {
		statusIcon := " "
		switch phase.Status {
		case types.PhaseStatusCompleted:
			statusIcon = "✓"
		case types.PhaseStatusRunning:
			statusIcon = "→"
		case types.PhaseStatusFailed:
			statusIcon = "✗"
		case types.PhaseStatusPending:
			statusIcon = "○"
		}

		fmt.Printf("[%s] %d. %s (%s/%s)\n",
			statusIcon,
			phase.Sequence,
			phase.DisplayName,
			phase.AgentName,
			phase.ModelName)

		if phase.OutputText != "" && len(phase.OutputText) > 0 {
			fmt.Printf("    Output: %s\n", truncate(phase.OutputText, 100))
		}
		if phase.Error != "" && len(phase.Error) > 0 {
			fmt.Printf("    Error: %s\n", truncate(phase.Error, 100))
		}
	}

	return nil
}

func (c *CLI) runTaskEvents(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	limit, _ := cmd.Flags().GetInt("limit")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	events, err := c.repo.GetEventsByTask(taskID, limit)
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	if jsonOutput {
		c.printOutput(events, true)
		return nil
	}

	for _, event := range events {
		fmt.Printf("[%s] %s: %s\n",
			event.OccurredAt.Format("15:04:05"),
			event.EventLevel,
			event.EventType)

		if event.Data != "" {
			fmt.Printf("    %s\n", event.Data)
		}
	}

	return nil
}

// runTaskExport 导出任务状态到文件（兼容旧版 skills）
func (c *CLI) runTaskExport(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	outputFile := args[1]

	ctx := context.Background()

	// 获取任务状态
	status, err := c.taskSvc.GetTaskStatus(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get task status: %w", err)
	}

	// 获取执行计划
	plan, err := c.taskSvc.GetExecutionPlan(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to get execution plan: %w", err)
	}

	// 生成 Markdown 文件（YAML frontmatter 格式）
	content := c.generateStateFileMarkdown(status, plan)

	// 写入文件
	if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	fmt.Printf("State file exported to: %s\n", outputFile)
	return nil
}

func (c *CLI) runTaskPause(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	if err := c.taskSvc.PauseTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to pause task: %w", err)
	}

	result := map[string]interface{}{
		"task_id": taskID,
		"status":  "paused",
		"message": "Task paused successfully",
	}

	c.printOutput(result, jsonOutput)
	return nil
}

func (c *CLI) runTaskResume(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	if err := c.taskSvc.ResumeTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to resume task: %w", err)
	}

	result := map[string]interface{}{
		"task_id": taskID,
		"status":  "resumed",
		"message": "Task resumed successfully",
	}

	c.printOutput(result, jsonOutput)
	return nil
}

func (c *CLI) runTaskCancel(cmd *cobra.Command, args []string) error {
	taskID := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	if err := c.taskSvc.CancelTask(ctx, taskID); err != nil {
		return fmt.Errorf("failed to cancel task: %w", err)
	}

	result := map[string]interface{}{
		"task_id": taskID,
		"status":  "cancelled",
		"message": "Task cancelled successfully",
	}

	c.printOutput(result, jsonOutput)
	return nil
}

// ============================================================================
// Workflow commands (快捷命令)
// ============================================================================

// runWorkflowStart 一键开始工作流
func (c *CLI) runWorkflowStart(cmd *cobra.Command, args []string) error {
	prompt := args[0]
	workflow, _ := cmd.Flags().GetString("workflow")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// workflow 参数必需（不再自动检测）
	if workflow == "" {
		return fmt.Errorf("workflow is required. Use --workflow <name> or let LLM select from list")
	}

	ctx := context.Background()

	// 创建任务（workflow 参数作为工作流名称）
	task, err := c.taskSvc.CreateTask(ctx, prompt, workflow)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// 获取当前阶段
	phase, err := c.taskSvc.GetCurrentPhase(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("failed to get current phase: %w", err)
	}

	// 加载流水线获取阶段详情
	wf, _ := c.parser.LoadWorkflow(task.WorkflowID)
	node := findNode(wf, phase.PhaseName)

	output := map[string]interface{}{
		"task_id":       task.ID,
		"workflow":      task.WorkflowID,
		"workflow_type": task.WorkflowType,
		"total_phases":  task.TotalPhases,
		"current_phase": map[string]interface{}{
			"phase_id":     phase.ID,
			"phase_name":   phase.PhaseName,
			"display_name": phase.DisplayName,
			"agent":        phase.AgentName,
			"model":        phase.ModelName,
			"sequence":     phase.Sequence,
			"required":     node.Required,
			"timeout":      node.Timeout,
		},
		"message": "Workflow started successfully",
	}

	c.printOutput(output, jsonOutput)
	return nil
}

// ============================================================================
// Workflow commands
// ============================================================================

func (c *CLI) runPipelineList(cmd *cobra.Command, args []string) error {
	workflows, err := c.parser.ListWorkflows()
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	c.printOutput(workflows, jsonOutput)
	return nil
}

func (c *CLI) runPipelineGet(cmd *cobra.Command, args []string) error {
	name := args[0]

	workflow, err := c.parser.LoadWorkflow(name)
	if err != nil {
		return fmt.Errorf("failed to load workflow: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	c.printOutput(workflow, jsonOutput)
	return nil
}

func (c *CLI) runPipelineRecommend(cmd *cobra.Command, args []string) error {
	prompt := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// 不再自动检测，提示用户使用 workflow list 查看
	output := map[string]interface{}{
		"message": "Auto-detection removed. Please use 'gclm-engine workflow list --json' to see available workflows and let LLM select based on prompt.",
		"prompt":  prompt,
		"hint":    "Skills should call 'workflow list --json' and semantically match the best workflow",
	}

	c.printOutput(output, jsonOutput)
	return nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// generateStateFileMarkdown 生成状态文件 Markdown 内容
func (c *CLI) generateStateFileMarkdown(status *service.TaskStatusResponse, plan *service.ExecutionPlan) string {
	var sb strings.Builder

	// YAML frontmatter
	sb.WriteString("---\n")
	active := status.Status == types.TaskStatusRunning
	sb.WriteString(fmt.Sprintf("active: %v\n", active))
	sb.WriteString(fmt.Sprintf("current_phase: %d\n", status.CurrentPhase))
	sb.WriteString(fmt.Sprintf("workflow_type: %s\n", status.WorkflowType))
	sb.WriteString(fmt.Sprintf("total_phases: %d\n", status.TotalPhases))
	sb.WriteString("phases:\n")

	for _, step := range plan.Steps {
		sb.WriteString(fmt.Sprintf("  - sequence: %d\n", step.Sequence))
		sb.WriteString(fmt.Sprintf("    name: %s\n", step.PhaseName))
		sb.WriteString(fmt.Sprintf("    display_name: %s\n", step.DisplayName))
		sb.WriteString(fmt.Sprintf("    agent: %s\n", step.Agent))
		sb.WriteString(fmt.Sprintf("    model: %s\n", step.Model))
		sb.WriteString(fmt.Sprintf("    status: %s\n", step.Status))
		if len(step.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("    depends_on: %v\n", step.Dependencies))
		}
	}

	sb.WriteString("---\n\n")

	// Markdown 内容
	sb.WriteString(fmt.Sprintf("# Task: %s\n\n", status.TaskID))
	sb.WriteString(fmt.Sprintf("**Status**: %s\n", status.Status))
	sb.WriteString(fmt.Sprintf("**Workflow**: %s\n", status.WorkflowType))
	sb.WriteString(fmt.Sprintf("**Progress**: Phase %d/%d\n\n", status.CurrentPhase, status.TotalPhases))

	sb.WriteString("## Phases\n\n")

	for _, phaseStatus := range status.Phases {
		statusIcon := "○"
		switch phaseStatus.Status {
		case types.PhaseStatusCompleted:
			statusIcon = "✓"
		case types.PhaseStatusRunning:
			statusIcon = "→"
		case types.PhaseStatusFailed:
			statusIcon = "✗"
		}

		sb.WriteString(fmt.Sprintf("%s **%d. %s** (%s/%s)\n",
			statusIcon,
			phaseStatus.Sequence,
			phaseStatus.DisplayName,
			phaseStatus.Agent,
			phaseStatus.Model))
	}

	return sb.String()
}

// printOutput 输出结果（JSON 或文本）
func (c *CLI) printOutput(data interface{}, jsonOutput bool) {
	if jsonOutput {
		pretty, _ := c.rootCmd.PersistentFlags().GetBool("pretty")
		if pretty {
			jsonBytes, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			jsonBytes, _ := json.Marshal(data)
			fmt.Println(string(jsonBytes))
		}
	} else {
		// 默认也用 JSON 输出，方便解析
		jsonBytes, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonBytes))
	}
}

// Run executes the CLI
func (c *CLI) Run() error {
	return c.rootCmd.Execute()
}

// Close closes the CLI resources
func (c *CLI) Close() error {
	return c.db.Close()
}

// Init initializes the CLI for use in main.go
func Init(configDir string) (*CLI, error) {
	return New(configDir)
}

// ============================================================================
// Helper functions
// ============================================================================

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// runWorkflowValidate 验证工作流配置文件
func (c *CLI) runWorkflowValidate(cmd *cobra.Command, args []string) error {
	yamlFile := args[0]

	if _, err := os.Stat(yamlFile); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", yamlFile)
	}

	wf, err := c.parser.LoadYAMLFile(yamlFile)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if wf.Name == "" {
		return fmt.Errorf("'name' field is required")
	}

	if wf.WorkflowType == "" {
		return fmt.Errorf("'workflow_type' field is required")
	}

	if len(wf.Nodes) == 0 {
		return fmt.Errorf("workflow must have at least one node")
	}

	fmt.Printf("OK: Workflow validation successful\n")
	fmt.Printf("  Name: %s\n", wf.Name)
	fmt.Printf("  Type: %s\n", wf.WorkflowType)
	fmt.Printf("  Nodes: %d\n", len(wf.Nodes))
	fmt.Printf("\nTo install: gclm-engine workflow install %s\n", yamlFile)

	return nil
}

// runWorkflowInstall 安装工作流
func (c *CLI) runWorkflowInstall(cmd *cobra.Command, args []string) error {
	yamlFile := args[0]
	customName, _ := cmd.Flags().GetString("name")

	// Read YAML file
	input, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}

	// Validate YAML structure
	if _, err := c.parser.LoadYAMLFile(yamlFile); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Determine workflow name
	workflowName := customName
	if workflowName == "" {
		baseName := filepath.Base(yamlFile)
		workflowName = strings.TrimSuffix(baseName, filepath.Ext(baseName))
	}

	workflowName = strings.ToLower(workflowName)
	workflowName = strings.ReplaceAll(workflowName, "-", "_")
	workflowName = strings.ReplaceAll(workflowName, " ", "_")

	// Install: copy YAML file to workflows directory
	workflowsDir := filepath.Join(c.configDir, "workflows")
	destPath := filepath.Join(workflowsDir, workflowName+".yaml")

	if err := os.WriteFile(destPath, input, 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	fmt.Printf("OK: Workflow '%s' installed\n", workflowName)
	fmt.Printf("  File: %s\n", destPath)
	fmt.Printf("  Note: Run 'gclm-engine workflow sync' to publish to database\n")

	return nil
}

// runWorkflowUninstall 卸载工作流
func (c *CLI) runWorkflowUninstall(cmd *cobra.Command, args []string) error {
	workflowName := args[0]

	// Check if it's a builtin workflow
	wfRepo := db.NewWorkflowRepository(c.db)
	record, err := wfRepo.GetWorkflow(workflowName)
	if err != nil {
		return err
	}

	if record.IsBuiltin {
		return fmt.Errorf("cannot uninstall builtin workflow '%s'", workflowName)
	}

	// Uninstall: delete YAML file
	workflowsDir := filepath.Join(c.configDir, "workflows")
	yamlPath := filepath.Join(workflowsDir, workflowName+".yaml")

	if err := os.Remove(yamlPath); err != nil {
		return fmt.Errorf("failed to remove workflow file: %w", err)
	}

	// Remove from database
	_, err = c.db.GetDB().Exec("DELETE FROM workflows WHERE name = ?", workflowName)
	if err != nil {
		return fmt.Errorf("failed to remove from database: %w", err)
	}

	fmt.Printf("OK: Workflow '%s' uninstalled\n", workflowName)
	return nil
}

// runWorkflowList 列出所有工作流
func (c *CLI) runWorkflowList(cmd *cobra.Command, args []string) error {
	wfRepo := db.NewWorkflowRepository(c.db)
	workflows, err := wfRepo.ListWorkflows()
	if err != nil {
		return fmt.Errorf("failed to list: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")

	if jsonOutput {
		// JSON 输出：返回工作流列表供 LLM 匹配
		result := make([]map[string]interface{}, 0, len(workflows))
		for _, w := range workflows {
			result = append(result, map[string]interface{}{
				"name":         w.Name,
				"display_name": w.DisplayName,
				"description":  w.Description,
				"workflow_type": w.WorkflowType,
				"version":      w.Version,
				"is_builtin":   w.IsBuiltin,
			})
		}
		c.printOutput(result, true)
		return nil
	}

	// 文本格式输出
	if len(workflows) == 0 {
		fmt.Println("No workflows found")
		return nil
	}

	fmt.Println("\nBuilt-in workflows:")
	hasBuiltin := false
	for _, w := range workflows {
		if w.IsBuiltin {
			fmt.Printf("  %-20s %-15s %s\n", w.Name, w.WorkflowType, w.DisplayName)
			hasBuiltin = true
		}
	}
	if !hasBuiltin {
		fmt.Println("  (none)")
	}

	fmt.Println("\nCustom workflows:")
	hasCustom := false
	for _, w := range workflows {
		if !w.IsBuiltin {
			fmt.Printf("  %-20s %-15s %s\n", w.Name, w.WorkflowType, w.DisplayName)
			hasCustom = true
		}
	}
	if !hasCustom {
		fmt.Println("  (none)")
	}

	return nil
}

// runWorkflowExport 导出工作流
func (c *CLI) runWorkflowExport(cmd *cobra.Command, args []string) error {
	workflowName := args[0]

	// Get workflow from database
	wfRepo := db.NewWorkflowRepository(c.db)
	record, err := wfRepo.GetWorkflow(workflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	outputFile := workflowName + ".yaml"
	if len(args) > 1 {
		outputFile = args[1]
	}

	if err := os.WriteFile(outputFile, []byte(record.ConfigYAML), 0644); err != nil {
		return fmt.Errorf("failed to export: %w", err)
	}

	fmt.Printf("OK: Exported to %s\n", outputFile)
	return nil
}

// runWorkflowInfo 显示工作流详细信息
func (c *CLI) runWorkflowInfo(cmd *cobra.Command, args []string) error {
	workflowName := args[0]

	// Get workflow from database
	wfRepo := db.NewWorkflowRepository(c.db)
	record, err := wfRepo.GetWorkflow(workflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	// Parse YAML to get workflow details
	var wfDef types.Workflow
	if err := yaml.Unmarshal([]byte(record.ConfigYAML), &wfDef); err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	fmt.Printf("Workflow: %s\n", wfDef.Name)
	fmt.Printf("  Display: %s\n", wfDef.DisplayName)
	fmt.Printf("  Type: %s\n", wfDef.WorkflowType)
	fmt.Printf("  Version: %s\n\n", wfDef.Version)

	fmt.Printf("Nodes (%d):\n", len(wfDef.Nodes))
	for i, node := range wfDef.Nodes {
		deps := ""
		if len(node.DependsOn) > 0 {
			deps = " (after: " + strings.Join(node.DependsOn, ", ") + ")"
		}

		fmt.Printf("  %d. %s\n", i+1, node.DisplayName)
		fmt.Printf("     Ref: %s, %s/%s%s\n", node.Ref, node.Agent, node.Model, deps)
	}

	return nil
}

// runWorkflowInit 初始化 gclm-engine 配置
func (c *CLI) runWorkflowInit(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	silent, _ := cmd.Flags().GetBool("silent")

	if !silent {
		fmt.Println("Initializing gclm-engine...")
	}

	// 1. Export default config
	if !silent {
		fmt.Println("\n[1/3] Setting up configuration...")
	}

	created, err := assets.ExportDefaultConfig(c.configDir, force)
	if err != nil {
		return fmt.Errorf("failed to export default config: %w", err)
	}
	if created && !silent {
		fmt.Printf("  ✓ Created: %s\n", filepath.Join(c.configDir, "gclm_engine_config.yaml"))
	} else if !silent {
		fmt.Printf("  − Config exists (use --force to overwrite)\n")
	}

	// 2. Export builtin workflows
	if !silent {
		fmt.Println("\n[2/3] Setting up workflow definitions...")
	}

	workflowsDir := filepath.Join(c.configDir, "workflows")
	exported, err := assets.ExportBuiltinWorkflows(workflowsDir, force)
	if err != nil {
		return fmt.Errorf("failed to export builtin workflows: %w", err)
	}
	if len(exported) > 0 && !silent {
		for _, name := range exported {
			fmt.Printf("  ✓ Created: %s\n", filepath.Join(workflowsDir, name))
		}
	} else if !silent {
		fmt.Printf("  − Workflows exist (use --force to overwrite)\n")
	}

	// 3. Initialize database and load workflows
	if !silent {
		fmt.Println("\n[3/3] Initializing database...")
	}

	// Initialize builtin workflows from the exported YAML files
	wfRepo := db.NewWorkflowRepository(c.db)
	if err := wfRepo.InitializeBuiltinWorkflows(workflowsDir); err != nil {
		return fmt.Errorf("failed to initialize workflows: %w", err)
	}
	if !silent {
		fmt.Printf("  ✓ Database initialized\n")
	}

	if !silent {
		fmt.Println("\n✓ Initialization complete!")
		fmt.Printf("  Config: %s\n", filepath.Join(c.configDir, "gclm_engine_config.yaml"))
		fmt.Printf("  Workflows: %s\n", workflowsDir)
		fmt.Printf("  Database: %s\n", filepath.Join(c.configDir, "gclm-engine.db"))
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Edit workflow YAML files in the workflows directory (optional)")
		fmt.Println("  2. Run 'gclm-engine workflow sync' to publish changes to database")
		fmt.Println("  3. Run 'gclm-engine workflow list' to see available workflows")
	}

	return nil
}

// runWorkflowSync 同步工作流 YAML 文件到数据库
func (c *CLI) runWorkflowSync(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")
	workflowsDir := filepath.Join(c.configDir, "workflows")
	wfRepo := db.NewWorkflowRepository(c.db)

	if len(args) == 0 {
		// Sync all workflows
		fmt.Printf("Syncing workflows from %s...\n\n", workflowsDir)

		entries, err := os.ReadDir(workflowsDir)
		if err != nil {
			return fmt.Errorf("failed to read workflows directory: %w", err)
		}

		successCount := 0
		skipCount := 0
		errorCount := 0

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
				continue
			}

			workflowName := strings.TrimSuffix(entry.Name(), ".yaml")
			yamlPath := filepath.Join(workflowsDir, entry.Name())

			result, err := c.syncOneWorkflow(wfRepo, workflowName, yamlPath, force)
			if err != nil {
				fmt.Fprintf(os.Stderr, "✗ %s: %v\n", workflowName, err)
				errorCount++
			} else if result == "skipped" {
				fmt.Printf("− %s: unchanged\n", workflowName)
				skipCount++
			} else {
				fmt.Printf("✓ %s: synced\n", workflowName)
				successCount++
			}
		}

		fmt.Printf("\nSync complete: %d synced, %d unchanged, %d errors\n",
			successCount, skipCount, errorCount)

		if errorCount > 0 && !force {
			return fmt.Errorf("%d workflows failed to sync", errorCount)
		}

		return nil
	}

	// Sync single workflow by file path
	yamlPath := args[0]

	// Expand path if relative
	if !filepath.IsAbs(yamlPath) {
		// Try workflows directory first
		workflowsDirPath := filepath.Join(workflowsDir, yamlPath)
		if _, err := os.Stat(workflowsDirPath); err == nil {
			yamlPath = workflowsDirPath
		} else {
			// Use current directory
			absPath, err := filepath.Abs(yamlPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path: %w", err)
			}
			yamlPath = absPath
		}
	}

	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		return fmt.Errorf("workflow file not found: %s", yamlPath)
	}

	// Extract workflow name from filename
	workflowName := strings.TrimSuffix(filepath.Base(yamlPath), ".yaml")

	result, err := c.syncOneWorkflow(wfRepo, workflowName, yamlPath, force)
	if err != nil {
		return err
	}

	if result == "skipped" {
		fmt.Printf("%s: unchanged (no sync needed)\n", yamlPath)
	} else {
		fmt.Printf("%s: synced successfully\n", yamlPath)
	}

	return nil
}

// syncResult 表示同步结果
type syncResult string

const (
	syncResultSuccess syncResult = "synced"
	syncResultSkipped syncResult = "skipped"
)

// syncOneWorkflow 同步单个工作流
func (c *CLI) syncOneWorkflow(wfRepo *db.WorkflowRepository, workflowName, yamlPath string, force bool) (syncResult, error) {
	// Read YAML file
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Check if workflow exists and compare
	var existingYAML string
	var exists bool
	err = c.db.GetDB().QueryRow("SELECT config_yaml FROM workflows WHERE name = ?", workflowName).
		Scan(&existingYAML)

	if err == sql.ErrNoRows {
		exists = false
	} else if err != nil {
		return "", fmt.Errorf("failed to check existing workflow: %w", err)
	} else {
		exists = true
	}

	// Check if content changed
	if exists && string(yamlData) == existingYAML {
		return syncResultSkipped, nil
	}

	// Validate YAML before syncing (unless force)
	if !force {
		var wf types.Workflow
		if err := yaml.Unmarshal(yamlData, &wf); err != nil {
			return "", fmt.Errorf("YAML validation failed: %w", err)
		}

		if wf.Name == "" {
			return "", fmt.Errorf("workflow name is required")
		}

		if wf.WorkflowType == "" {
			return "", fmt.Errorf("workflow_type is required")
		}

		// Validate against config
		cfg, err := config.Load()
		if err != nil {
			return "", fmt.Errorf("failed to load config: %w", err)
		}

		if err := cfg.ValidateWorkflowType(wf.WorkflowType); err != nil {
			return "", fmt.Errorf("invalid workflow_type: %w", err)
		}
	}

	// Insert or update
	if exists {
		_, err = c.db.GetDB().Exec(`
			UPDATE workflows
			SET config_yaml = ?, display_name = ?, description = ?,
			    workflow_type = ?, version = ?, updated_at = ?
			WHERE name = ?
		`, string(yamlData), extractDisplayName(yamlData), extractDescription(yamlData),
			extractWorkflowType(yamlData), extractVersion(yamlData),
			time.Now().Format(time.RFC3339), workflowName)
	} else {
		_, err = c.db.GetDB().Exec(`
			INSERT INTO workflows (name, display_name, description, workflow_type, version, is_builtin, config_yaml)
			VALUES (?, ?, ?, ?, ?, 0, ?)
		`, workflowName, extractDisplayName(yamlData), extractDescription(yamlData),
			extractWorkflowType(yamlData), extractVersion(yamlData), string(yamlData))
	}

	if err != nil {
		return "", fmt.Errorf("failed to save to database: %w", err)
	}

	return syncResultSuccess, nil
}

// Helper functions to extract metadata from YAML without full parsing
func extractDisplayName(yamlData []byte) string {
	if match := findYAMLField(yamlData, "display_name"); match != "" {
		return match
	}
	return ""
}

func extractDescription(yamlData []byte) string {
	if match := findYAMLField(yamlData, "description"); match != "" {
		return match
	}
	return ""
}

func extractWorkflowType(yamlData []byte) string {
	if match := findYAMLField(yamlData, "workflow_type"); match != "" {
		return match
	}
	return ""
}

func extractVersion(yamlData []byte) string {
	if match := findYAMLField(yamlData, "version"); match != "" {
		return match
	}
	return "1.0.0"
}

func findYAMLField(data []byte, field string) string {
	// Simple regex-like search for field: value
	prefix := []byte(field + ":")
	idx := bytes.Index(data, prefix)
	if idx == -1 {
		return ""
	}

	// Find value after prefix
	rest := data[idx+len(prefix):]
	rest = bytes.TrimLeft(rest, " \t")

	if len(rest) == 0 {
		return ""
	}

	// Extract quoted string or unquoted value
	if rest[0] == '"' {
		end := bytes.Index(rest[1:], []byte("\""))
		if end == -1 {
			return ""
		}
		return string(rest[1 : end+1])
	}

	// Extract until newline or comment
	end := bytes.IndexAny(rest, "\n#")
	if end == -1 {
		end = len(rest)
	}
	return strings.TrimSpace(string(rest[:end]))
}

func findNode(wf *types.Workflow, ref string) *types.WorkflowNode {
	for i := range wf.Nodes {
		if wf.Nodes[i].Ref == ref {
			return &wf.Nodes[i]
		}
	}
	return nil
}

func errMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// printFriendlyError 打印友好的错误信息
func (c *CLI) printFriendlyError(err error) {
	if err == nil {
		return
	}

	// 检查是否为友好错误
	if friendly, ok := err.(*errors.FriendlyError); ok {
		fmt.Fprintln(os.Stderr, friendly.FormatUser())
		return
	}

	// 对于非友好错误，尝试包装
	errStr := err.Error()

	// 根据错误内容提供友好提示
	switch {
	case strings.Contains(errStr, "not found"):
		if strings.Contains(errStr, "task") {
			fmt.Fprintln(os.Stderr, errors.TaskNotFound(extractID(errStr, "task")).FormatUser())
		} else if strings.Contains(errStr, "workflow") || strings.Contains(errStr, "pipeline") {
			fmt.Fprintln(os.Stderr, errors.WorkflowNotFound(extractID(errStr, "")).FormatUser())
		}
	case strings.Contains(errStr, "failed to load pipeline"):
		fmt.Fprintln(os.Stderr, errors.PipelineLoadError(extractID(errStr, ""), err).FormatUser())
	case strings.Contains(errStr, "no such file"):
		fmt.Fprintln(os.Stderr, errors.ConfigDirectoryNotFound(extractPath(errStr)).FormatUser())
	case strings.Contains(errStr, "yaml") || strings.Contains(errStr, "unmarshal"):
		fmt.Fprintln(os.Stderr, errors.InvalidYAMLFormat(extractPath(errStr), err).FormatUser())
	default:
		// 默认输出原始错误
		fmt.Fprintf(os.Stderr, "❌ 错误: %s\n", errStr)
		fmt.Fprintln(os.Stderr, "\n💡 如需帮助，请运行 `gclm-engine --help`")
	}
}

// extractID 从错误消息中提取 ID
func extractID(errStr, prefix string) string {
	// 简单实现：查找引号或单引号中的内容
	start := strings.Index(errStr, "'")
	if start == -1 {
		start = strings.Index(errStr, "\"")
	}
	if start == -1 {
		return "unknown"
	}
	end := strings.Index(errStr[start+1:], "'")
	if end == -1 {
		end = strings.Index(errStr[start+1:], "\"")
	}
	if end == -1 {
		return "unknown"
	}
	return errStr[start+1 : start+1+end]
}

// extractPath 从错误消息中提取路径
func extractPath(errStr string) string {
	// 简单实现：查找 .yaml 或 .yml 文件
	yamlIdx := strings.Index(errStr, ".yaml")
	if yamlIdx == -1 {
		yamlIdx = strings.Index(errStr, ".yml")
	}
	if yamlIdx == -1 {
		return "unknown"
	}
	// 向前查找文件名开始
	start := strings.LastIndex(errStr[:yamlIdx], "/")
	if start == -1 {
		start = strings.LastIndex(errStr[:yamlIdx], " ")
	}
	if start == -1 {
		return "unknown"
	}
	return errStr[start+1 : yamlIdx+5]
}
