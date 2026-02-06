package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// createWorkflowCommand creates workflow commands
func (c *CLI) createWorkflowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Workflow management commands",
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
		RunE:  c.runWorkflowList,
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
		Use: "sync [yaml-file]",
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

	cmd.AddCommand(validateCmd, installCmd, uninstallCmd, listCmd, exportCmd, infoCmd, syncCmd)

	return cmd
}

// ============================================================================
// Workflow command implementations
// ============================================================================

func (c *CLI) runWorkflowValidate(cmd *cobra.Command, args []string) error {
	yamlFile := args[0]

	ctx := context.Background()
	workflow, err := c.workflowSvc.ValidateWorkflow(ctx, yamlFile)
	if err != nil {
		return fmt.Errorf("workflow validation failed: %w", err)
	}

	fmt.Printf("Workflow '%s' is valid\n", workflow.Name)
	fmt.Printf("  Type: %s\n", workflow.WorkflowType)
	fmt.Printf("  Nodes: %d\n", len(workflow.Nodes))
	return nil
}

func (c *CLI) runWorkflowInstall(cmd *cobra.Command, args []string) error {
	yamlFile := args[0]
	name, _ := cmd.Flags().GetString("name")

	// 如果没有指定名称，从文件名提取
	if name == "" {
		name = strings.TrimSuffix(filepath.Base(yamlFile), ".yaml")
	}

	ctx := context.Background()
	if err := c.workflowSvc.InstallWorkflow(ctx, name, yamlFile); err != nil {
		return fmt.Errorf("failed to install workflow: %w", err)
	}

	fmt.Printf("Workflow '%s' installed successfully\n", name)
	return nil
}

func (c *CLI) runWorkflowUninstall(cmd *cobra.Command, args []string) error {
	name := args[0]

	ctx := context.Background()
	if err := c.workflowSvc.UninstallWorkflow(ctx, name); err != nil {
		return fmt.Errorf("failed to uninstall workflow: %w", err)
	}

	fmt.Printf("Workflow '%s' uninstalled successfully\n", name)
	return nil
}

func (c *CLI) runWorkflowList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	workflows, err := c.workflowSvc.ListWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	if len(workflows) == 0 {
		fmt.Println("No workflows found")
		return nil
	}

	for _, wf := range workflows {
		builtin := ""
		if strings.HasPrefix(wf.Name, "builtin:") {
			builtin = " [builtin]"
		}
		fmt.Printf("- %s%s\n", wf.Name, builtin)
		fmt.Printf("  Display: %s\n", wf.DisplayName)
		fmt.Printf("  Type: %s\n", wf.WorkflowType)
		fmt.Printf("  Version: %s\n", wf.Version)
		if wf.Description != "" {
			fmt.Printf("  Description: %s\n", wf.Description)
		}
		fmt.Println()
	}

	return nil
}

func (c *CLI) runWorkflowExport(cmd *cobra.Command, args []string) error {
	name := args[0]

	outputFile := name + ".yaml"
	if len(args) > 1 {
		outputFile = args[1]
	}

	ctx := context.Background()
	data, err := c.workflowSvc.ExportWorkflow(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to export workflow: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Workflow exported to: %s\n", outputFile)
	return nil
}

func (c *CLI) runWorkflowInfo(cmd *cobra.Command, args []string) error {
	name := args[0]
	jsonOutput, _ := cmd.Flags().GetBool("json")

	ctx := context.Background()
	detail, err := c.workflowSvc.GetWorkflow(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get workflow info: %w", err)
	}

	c.printOutput(detail, jsonOutput)
	return nil
}

func (c *CLI) runWorkflowSync(cmd *cobra.Command, args []string) error {
	var yamlDir string

	if len(args) > 0 {
		// 同步单个文件
		yamlFile := args[0]
		name := strings.TrimSuffix(filepath.Base(yamlFile), ".yaml")

		ctx := context.Background()
		data, err := os.ReadFile(yamlFile)
		if err != nil {
			return fmt.Errorf("failed to read workflow file: %w", err)
		}

		if err := c.workflowRepo.InstallWorkflow(ctx, name, data); err != nil {
			// 如果已存在，尝试更新
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("failed to sync workflow: %w", err)
			}
			fmt.Printf("Workflow '%s' already exists (use 'workflow uninstall' first to replace)\n", name)
			return nil
		}

		fmt.Printf("Workflow '%s' synced successfully\n", name)
		return nil
	}

	// 同步所有工作流
	yamlDir = c.workflowsDir

	ctx := context.Background()
	synced, err := c.workflowSvc.SyncWorkflows(ctx, yamlDir)
	if err != nil {
		return fmt.Errorf("failed to sync workflows: %w", err)
	}

	if len(synced) == 0 {
		fmt.Println("No new workflows to sync")
		return nil
	}

	fmt.Printf("Synced %d workflows:\n", len(synced))
	for _, name := range synced {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}
