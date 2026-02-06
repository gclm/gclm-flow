package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"github.com/spf13/cobra"
)

// createPipelineCommand creates pipeline management commands
func (c *CLI) createPipelineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Pipeline management commands",
	}

	// pipeline list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all pipelines",
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
		Short: "Recommend a pipeline for the task",
		Long:  "Analyze the task prompt and recommend the most suitable pipeline.",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runPipelineRecommend,
	}

	cmd.AddCommand(listCmd, getCmd, recommendCmd)

	return cmd
}

// ============================================================================
// Pipeline command implementations
// ============================================================================

func (c *CLI) runPipelineList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// 获取所有工作流（即 pipeline）
	workflows, err := c.workflowRepo.ListWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}

	if len(workflows) == 0 {
		fmt.Println("No pipelines found")
		return nil
	}

	// 按 workflow_type 分组
	byType := make(map[string][]*domain.WorkflowRecord)
	for _, wf := range workflows {
		byType[wf.WorkflowType] = append(byType[wf.WorkflowType], wf)
	}

	// 按标准顺序显示
	types := []string{"document", "code_simple", "code_complex", "fix", "refactor", "test", "docs"}

	for _, t := range types {
		wfs, ok := byType[t]
		if !ok || len(wfs) == 0 {
			continue
		}

		fmt.Printf("\n[%s]\n", strings.ToUpper(t))
		for _, wf := range wfs {
			fmt.Printf("  %s - %s\n", wf.Name, wf.DisplayName)
		}
	}

	return nil
}

func (c *CLI) runPipelineGet(cmd *cobra.Command, args []string) error {
	name := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	detail, err := c.workflowSvc.GetWorkflow(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %w", err)
	}

	c.printOutput(detail, jsonOutput)
	return nil
}

func (c *CLI) runPipelineRecommend(cmd *cobra.Command, args []string) error {
	prompt := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	// 简单的关键词匹配算法
	promptLower := strings.ToLower(prompt)

	// 关键词映射
	typeKeywords := map[string][]string{
		"document":    {"文档", "设计", "方案", "规范", "说明", "编写文档"},
		"code_simple": {"修复", "fix", "bug", "小修改", "调整"},
		"code_complex": {"功能", "模块", "开发", "重构", "实现"},
		"fix":         {"修复", "fix", "bug", "错误"},
		"refactor":    {"重构", "优化", "改进"},
		"test":        {"测试", "test"},
		"docs":        {"文档", "docs"},
	}

	// 计算每个类型的得分
	scores := make(map[string]int)
	for typeName, keywords := range typeKeywords {
		for _, kw := range keywords {
			if strings.Contains(promptLower, kw) {
				scores[typeName]++
			}
		}
	}

	// 找出得分最高的类型
	var bestType string
	maxScore := 0
	for typeName, score := range scores {
		if score > maxScore {
			maxScore = score
			bestType = typeName
		}
	}

	// 如果没有匹配，默认使用 code_complex
	if bestType == "" {
		bestType = "code_complex"
	}

	// 获取该类型的第一个工作流
	ctx := context.Background()
	workflow, err := c.workflowRepo.GetWorkflowByType(ctx, bestType)
	if err != nil {
		return fmt.Errorf("no workflow found for type: %s", bestType)
	}

	result := map[string]interface{}{
		"prompt":        prompt,
		"recommended":   workflow.Name,
		"workflow_type": bestType,
		"display_name":  workflow.DisplayName,
		"confidence":    "high",
		"reasoning":     "Matched keywords: " + fmt.Sprintf("%v", scores[bestType]),
	}

	c.printOutput(result, jsonOutput)
	return nil
}

// 辅助函数：用于旧的代码兼容
type TaskStatus = types.TaskStatus
type PhaseStatus = types.PhaseStatus
