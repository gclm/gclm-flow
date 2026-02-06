package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

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

// generateStateFileMarkdown 生成状态文件 Markdown 内容
func (c *CLI) generateStateFileMarkdown(status *domain.TaskStatusResponse, plan *domain.ExecutionPlan) string {
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
			phaseStatus.AgentName,
			phaseStatus.ModelName))
	}

	return sb.String()
}

// printFriendlyError 打印友好的错误信息
func (c *CLI) printFriendlyError(err error) {
	if err == nil {
		return
	}

	// 对于非友好错误，尝试包装
	errStr := err.Error()

	// 根据错误内容提供友好提示
	switch {
	case strings.Contains(errStr, "not found"):
		if strings.Contains(errStr, "task") {
			fmt.Fprintf(os.Stderr, "Error: Task not found\n")
		} else if strings.Contains(errStr, "workflow") {
			fmt.Fprintf(os.Stderr, "Error: Workflow not found\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", errStr)
		}
	case strings.Contains(errStr, "required"):
		fmt.Fprintf(os.Stderr, "Error: Required phase failed: %s\n", errStr)
	default:
		fmt.Fprintf(os.Stderr, "Error: %s\n", errStr)
	}
}
