package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/pipeline"
	"github.com/gclm/gclm-flow/gclm-engine/internal/service"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// CLI represents the command-line interface
type CLI struct {
	rootCmd   *cobra.Command
	db        *db.Database
	parser    *pipeline.Parser
	repo      *db.Repository
	taskSvc   *service.TaskService
	configDir string
}

// New creates a new CLI instance
func New(configDir string) (*CLI, error) {
	// Initialize database
	database, err := db.New(db.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize builtin workflows from YAML files
	if err := database.InitWorkflows(configDir); err != nil {
		return nil, fmt.Errorf("failed to initialize workflows: %w", err)
	}

	// Initialize pipeline parser (still needed for YAML loading operations)
	parser := pipeline.NewParser(configDir)

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
	root.AddCommand(c.createTaskCommand())
	root.AddCommand(c.createPipelineCommand())
	root.AddCommand(c.createWorkflowCommand())
	root.AddCommand(c.createVersionCommand())

	return root
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

	cmd.AddCommand(startCmd, nextCmd, validateCmd, installCmd, uninstallCmd, listCmd, exportCmd, infoCmd)

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

	var task *types.Task
	var err error

	// 使用 TaskService 创建任务
	if pipelineName != "" {
		// 如果指定了 pipeline 名称，需要先获取 workflow_type
		pipe, err := c.parser.LoadPipeline(pipelineName)
		if err != nil {
			return fmt.Errorf("failed to load pipeline: %w", err)
		}
		task, err = c.taskSvc.CreateTask(ctx, prompt, pipe.WorkflowType)
	} else if workflowType != "" {
		task, err = c.taskSvc.CreateTask(ctx, prompt, workflowType)
	} else {
		// 自动检测
		task, err = c.taskSvc.CreateTask(ctx, prompt, "")
	}

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// 输出结果
	output := map[string]interface{}{
		"task_id":       task.ID,
		"status":        task.Status,
		"workflow_type": task.WorkflowType,
		"pipeline":      task.PipelineID,
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
		return fmt.Errorf("failed to get task: %w", err)
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
	pipe, _ := c.parser.LoadPipeline(task.PipelineID)
	node := findNode(pipe, phase.PhaseName)

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
	workflowType, _ := cmd.Flags().GetString("workflow-type")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()

	// 创建任务
	task, err := c.taskSvc.CreateTask(ctx, prompt, workflowType)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	// 获取当前阶段
	phase, err := c.taskSvc.GetCurrentPhase(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("failed to get current phase: %w", err)
	}

	// 加载流水线获取阶段详情
	pipe, _ := c.parser.LoadPipeline(task.PipelineID)
	node := findNode(pipe, phase.PhaseName)

	output := map[string]interface{}{
		"task_id":       task.ID,
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
// Pipeline commands
// ============================================================================

func (c *CLI) runPipelineList(cmd *cobra.Command, args []string) error {
	pipelines, err := c.parser.ListPipelines()
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	c.printOutput(pipelines, jsonOutput)
	return nil
}

func (c *CLI) runPipelineGet(cmd *cobra.Command, args []string) error {
	name := args[0]

	pipeline, err := c.parser.LoadPipeline(name)
	if err != nil {
		return fmt.Errorf("failed to load pipeline: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	c.printOutput(pipeline, jsonOutput)
	return nil
}

func (c *CLI) runPipelineRecommend(cmd *cobra.Command, args []string) error {
	prompt := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// 检测工作流类型
	detectedType := c.detectWorkflowType(prompt)

	pipeline, err := c.parser.GetPipelineByWorkflowType(detectedType)
	if err != nil {
		return fmt.Errorf("failed to find pipeline: %w", err)
	}

	output := map[string]interface{}{
		"workflow_type": detectedType,
		"pipeline":      pipeline.Name,
		"display_name":  pipeline.DisplayName,
		"description":   pipeline.Description,
		"total_nodes":   len(pipeline.Nodes),
	}

	c.printOutput(output, jsonOutput)
	return nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// detectWorkflowType 检测工作流类型（与 service 层保持一致）
func (c *CLI) detectWorkflowType(prompt string) string {
	docPhrases := []string{"编写文档", "文档编写", "方案设计", "设计文档", "需求分析", "技术方案", "架构设计", "api文档", "spec文档"}
	docKeywords := []string{"文档", "方案", "需求", "分析", "架构", "规范", "说明"}

	bugPhrases := []string{"修复bug", "fix bug", "bug修复", "修复错误", "解决bug"}
	bugKeywords := []string{"bug", "修复", "fix error", "error fix", "调试", "debug"}

	featureKeywords := []string{"功能", "模块", "开发", "重构", "实现"}

	score := 0

	// 文档类短语（+5分）
	for _, phrase := range docPhrases {
		if strings.Contains(strings.ToLower(prompt), strings.ToLower(phrase)) {
			score += 5
		}
	}
	// 文档类单词（+3分）
	for _, kw := range docKeywords {
		if contains(prompt, kw) {
			score += 3
		}
	}

	// Bug修复短语（-5分）
	for _, phrase := range bugPhrases {
		if strings.Contains(strings.ToLower(prompt), strings.ToLower(phrase)) {
			score -= 5
		}
	}
	// Bug修复单词（-3分）
	for _, kw := range bugKeywords {
		if contains(prompt, kw) {
			score -= 3
		}
	}

	// 功能开发单词（-1分）
	for _, kw := range featureKeywords {
		if contains(prompt, kw) {
			score -= 1
		}
	}

	// 分类
	if score >= 3 {
		return "DOCUMENT"
	} else if score <= -3 {
		return "CODE_SIMPLE"
	}
	return "CODE_COMPLEX"
}

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

	pipeline, err := c.parser.LoadYAMLFile(yamlFile)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if pipeline.Name == "" {
		return fmt.Errorf("'name' field is required")
	}

	if pipeline.WorkflowType == "" {
		return fmt.Errorf("'workflow_type' field is required")
	}

	if len(pipeline.Nodes) == 0 {
		return fmt.Errorf("workflow must have at least one node")
	}

	fmt.Printf("OK: Workflow validation successful\n")
	fmt.Printf("  Name: %s\n", pipeline.Name)
	fmt.Printf("  Type: %s\n", pipeline.WorkflowType)
	fmt.Printf("  Nodes: %d\n", len(pipeline.Nodes))
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

	// Install to database
	wfRepo := db.NewWorkflowRepository(c.db)
	if err := wfRepo.InstallWorkflow(workflowName, input); err != nil {
		return fmt.Errorf("failed to install: %w", err)
	}

	fmt.Printf("OK: Workflow '%s' installed\n", workflowName)
	fmt.Printf("  Use: gclm-engine workflow start \"<task>\" --workflow %s\n", workflowName)

	return nil
}

// runWorkflowUninstall 卸载工作流
func (c *CLI) runWorkflowUninstall(cmd *cobra.Command, args []string) error {
	workflowName := args[0]

	// Uninstall from database
	wfRepo := db.NewWorkflowRepository(c.db)
	if err := wfRepo.UninstallWorkflow(workflowName); err != nil {
		return err
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
	wf, err := wfRepo.GetWorkflow(workflowName)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	outputFile := workflowName + ".yaml"
	if len(args) > 1 {
		outputFile = args[1]
	}

	if err := os.WriteFile(outputFile, []byte(wf.ConfigYAML), 0644); err != nil {
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
	wf, err := wfRepo.GetWorkflow(workflowName)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	// Parse YAML to get pipeline details
	var pipeline types.Pipeline
	if err := yaml.Unmarshal([]byte(wf.ConfigYAML), &pipeline); err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	fmt.Printf("Workflow: %s\n", pipeline.Name)
	fmt.Printf("  Display: %s\n", pipeline.DisplayName)
	fmt.Printf("  Type: %s\n", pipeline.WorkflowType)
	fmt.Printf("  Version: %s\n\n", pipeline.Version)

	fmt.Printf("Nodes (%d):\n", len(pipeline.Nodes))
	for i, node := range pipeline.Nodes {
		deps := ""
		if len(node.DependsOn) > 0 {
			deps = " (after: " + strings.Join(node.DependsOn, ", ") + ")"
		}

		fmt.Printf("  %d. %s\n", i+1, node.DisplayName)
		fmt.Printf("     Ref: %s, %s/%s%s\n", node.Ref, node.Agent, node.Model, deps)
	}

	return nil
}

func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func findNode(pipeline *types.Pipeline, ref string) *types.PipelineNode {
	for i := range pipeline.Nodes {
		if pipeline.Nodes[i].Ref == ref {
			return &pipeline.Nodes[i]
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
