package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/internal/errors"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"github.com/spf13/cobra"
)

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
		Long:  "Create a new task using the specified workflow.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCreate,
	}
	createCmd.Flags().String("workflow", "", "Workflow name (analyze, docs, feat, fix)")
	createCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task get
	getCmd := &cobra.Command{
		Use:   "get <task-id>",
		Short: "Get task details",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskGet,
	}
	getCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE:  c.runTaskList,
	}
	listCmd.Flags().String("status", "", "Filter by status")
	listCmd.Flags().Int("limit", 20, "Maximum number of tasks to show")
	listCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task current - 获取当前应该执行的阶段
	currentCmd := &cobra.Command{
		Use:   "current <task-id>",
		Short: "Get current phase to execute",
		Long:  "Get the next pending phase that should be executed. Used by skills to determine what to do next.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCurrent,
	}
	currentCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task plan - 获取完整执行计划
	planCmd := &cobra.Command{
		Use:   "plan <task-id>",
		Short: "Get execution plan",
		Long:  "Get the complete execution plan with all phases and dependencies. Used by skills to understand the workflow.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskPlan,
	}
	planCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

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
	updateCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task complete - 完成阶段
	completeCmd := &cobra.Command{
		Use:   "complete <task-id> <phase-id>",
		Short: "Mark phase as completed",
		Long:  "Mark a phase as completed with output. Shortcut for 'task update ... completed'",
		Args:  cobra.ExactArgs(2),
		RunE:  c.runTaskComplete,
	}
	completeCmd.Flags().String("output", "", "Phase output")
	completeCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task fail - 标记阶段失败
	failCmd := &cobra.Command{
		Use:   "fail <task-id> <phase-id>",
		Short: "Mark phase as failed",
		Long:  "Mark a phase as failed with error message. Shortcut for 'task update ... failed'",
		Args:  cobra.ExactArgs(2),
		RunE:  c.runTaskFail,
	}
	failCmd.Flags().String("error", "", "Error message")
	failCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task phases
	phasesCmd := &cobra.Command{
		Use:   "phases <task-id>",
		Short: "Show task phases",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskPhases,
	}
	phasesCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	// task events
	eventsCmd := &cobra.Command{
		Use:   "events <task-id>",
		Short: "Show task events",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskEvents,
	}
	eventsCmd.Flags().Int("limit", 50, "Maximum number of events to show")
	eventsCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

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
	pauseCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	resumeCmd := &cobra.Command{
		Use:   "resume <task-id>",
		Short: "Resume a paused task",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskResume,
	}
	resumeCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	cancelCmd := &cobra.Command{
		Use:   "cancel <task-id>",
		Short: "Cancel a task",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runTaskCancel,
	}
	cancelCmd.Flags().BoolP("json", "j", false, "Output in JSON format")

	cmd.AddCommand(createCmd, getCmd, listCmd, currentCmd, planCmd, updateCmd,
		completeCmd, failCmd, phasesCmd, eventsCmd, exportCmd,
		pauseCmd, resumeCmd, cancelCmd)

	return cmd
}

// ============================================================================
// Task command implementations
// ============================================================================

func (c *CLI) runTaskCreate(cmd *cobra.Command, args []string) error {
	prompt := args[0]

	// Get flags
	workflow, _ := cmd.Flags().GetString("workflow")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()

	// workflow is required
	if workflow == "" {
		return fmt.Errorf("--workflow is required (use: analyze, docs, feat, fix)")
	}

	// 使用 TaskService 创建任务
	task, err := c.taskSvc.CreateTask(ctx, prompt, workflow)
	if err != nil {
		c.printFriendlyError(err)
		return err
	}

	// 输出结果
	output := map[string]interface{}{
		"task_id":       task.ID,
		"status":        task.Status,
		"workflow_type": task.WorkflowType,
		"workflow":      task.WorkflowID,
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

	ctx := context.Background()
	var status *types.TaskStatus
	if statusStr != "" {
		s := types.TaskStatus(statusStr)
		status = &s
	}

	// 使用 taskRepo 通过 adapter 接口获取任务
	tasks, err := c.taskRepo.ListTasks(ctx, status, limit)
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
	// 注意: 这里需要使用 taskRepo 获取任务，然后用 workflowLoader 加载流水线
	task, err := c.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	wf, err := c.workflowLoader.Load(ctx, task.WorkflowID)
	if err != nil {
		return fmt.Errorf("load workflow: %w", err)
	}
	node := findWorkflowNode(wf, phase.PhaseName)

	output := map[string]interface{}{
		"task_id":      taskID,
		"phase_id":     phase.ID,
		"phase_name":   phase.PhaseName,
		"display_name": phase.DisplayName,
		"agent":        phase.AgentName,
		"model":        phase.ModelName,
		"timeout":      node.Timeout,
		"required":     node.Required,
	}

	if len(node.DependsOn) > 0 {
		output["dependencies"] = node.DependsOn
	}

	c.printOutput(output, jsonOutput)
	return nil
}

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

	ctx := context.Background()
	phases, err := c.taskRepo.GetPhasesByTask(ctx, taskID)
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

	ctx := context.Background()
	events, err := c.taskRepo.GetEventsByTask(ctx, taskID, limit)
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

// 辅助函数
func errMsg(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func findWorkflowNode(wf *types.Workflow, ref string) *types.WorkflowNode {
	for i := range wf.Nodes {
		if wf.Nodes[i].Ref == ref {
			return &wf.Nodes[i]
		}
	}
	return nil
}
