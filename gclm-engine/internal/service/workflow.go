package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"gopkg.in/yaml.v3"
)

// WorkflowService 工作流管理服务
type WorkflowService struct {
	workflowRepo domain.WorkflowRepository
	loader       domain.WorkflowLoader
}

// NewWorkflowService 创建工作流服务
func NewWorkflowService(workflowRepo domain.WorkflowRepository, loader domain.WorkflowLoader) *WorkflowService {
	return &WorkflowService{
		workflowRepo: workflowRepo,
		loader:       loader,
	}
}

// ListWorkflows 列出所有工作流
func (s *WorkflowService) ListWorkflows(ctx context.Context) ([]*domain.WorkflowInfo, error) {
	records, err := s.workflowRepo.ListWorkflows(ctx)
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}

	result := make([]*domain.WorkflowInfo, len(records))
	for i, r := range records {
		result[i] = &domain.WorkflowInfo{
			Name:         r.Name,
			DisplayName:  r.DisplayName,
			Description:  r.Description,
			Version:      r.Version,
			WorkflowType: r.WorkflowType,
		}
	}
	return result, nil
}

// GetWorkflow 获取工作流详情
func (s *WorkflowService) GetWorkflow(ctx context.Context, name string) (*domain.WorkflowDetail, error) {
	record, err := s.workflowRepo.GetWorkflow(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}

	// 解析 YAML 获取节点信息
	var workflow types.Workflow
	if err := yaml.Unmarshal([]byte(record.ConfigYAML), &workflow); err != nil {
		return nil, fmt.Errorf("parse workflow YAML: %w", err)
	}

	nodes := make([]*domain.NodeDetail, len(workflow.Nodes))
	for i, node := range workflow.Nodes {
		nodes[i] = &domain.NodeDetail{
			Ref:          node.Ref,
			DisplayName:  node.DisplayName,
			Agent:        node.Agent,
			Model:        node.Model,
			Timeout:      node.Timeout,
			Required:     node.Required,
			Dependencies: node.DependsOn,
		}
	}

	return &domain.WorkflowDetail{
		Name:         record.Name,
		DisplayName:  record.DisplayName,
		Description:  record.Description,
		Version:      record.Version,
		WorkflowType: record.WorkflowType,
		IsBuiltin:    record.IsBuiltin,
		Nodes:        nodes,
		ConfigYAML:   record.ConfigYAML,
	}, nil
}

// GetWorkflowByType 按类型获取工作流详情
func (s *WorkflowService) GetWorkflowByType(ctx context.Context, workflowType string) (*domain.WorkflowDetail, error) {
	// List all workflows and find the one matching the type
	records, err := s.workflowRepo.ListWorkflows(ctx)
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}

	// Find the workflow with matching type
	var record *domain.WorkflowRecord
	for _, r := range records {
		if r.WorkflowType == workflowType {
			record = r
			break
		}
	}

	if record == nil {
		return nil, fmt.Errorf("workflow with type '%s' not found", workflowType)
	}

	// Parse YAML to get node information
	var workflow types.Workflow
	if err := yaml.Unmarshal([]byte(record.ConfigYAML), &workflow); err != nil {
		return nil, fmt.Errorf("parse workflow YAML: %w", err)
	}

	nodes := make([]*domain.NodeDetail, len(workflow.Nodes))
	for i, node := range workflow.Nodes {
		nodes[i] = &domain.NodeDetail{
			Ref:          node.Ref,
			DisplayName:  node.DisplayName,
			Agent:        node.Agent,
			Model:        node.Model,
			Timeout:      node.Timeout,
			Required:     node.Required,
			Dependencies: node.DependsOn,
		}
	}

	return &domain.WorkflowDetail{
		Name:         record.Name,
		DisplayName:  record.DisplayName,
		Description:  record.Description,
		Version:      record.Version,
		WorkflowType: record.WorkflowType,
		IsBuiltin:    record.IsBuiltin,
		Nodes:        nodes,
		ConfigYAML:   record.ConfigYAML,
	}, nil
}

// ValidateWorkflow 验证工作流配置
func (s *WorkflowService) ValidateWorkflow(ctx context.Context, yamlFile string) (*types.Workflow, error) {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("read workflow file: %w", err)
	}

	var workflow types.Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return nil, fmt.Errorf("parse workflow YAML: %w", err)
	}

	if err := s.loader.Validate(ctx, &workflow); err != nil {
		return nil, fmt.Errorf("validate workflow: %w", err)
	}

	return &workflow, nil
}

// InstallWorkflow 安装工作流
func (s *WorkflowService) InstallWorkflow(ctx context.Context, name string, yamlFile string) error {
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("read workflow file: %w", err)
	}

	// 验证工作流
	var workflow types.Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return fmt.Errorf("parse workflow YAML: %w", err)
	}

	if err := s.loader.Validate(ctx, &workflow); err != nil {
		return fmt.Errorf("validate workflow: %w", err)
	}

	// 安装
	if err := s.workflowRepo.InstallWorkflow(ctx, name, data); err != nil {
		return fmt.Errorf("install workflow: %w", err)
	}

	return nil
}

// UninstallWorkflow 卸载工作流
func (s *WorkflowService) UninstallWorkflow(ctx context.Context, name string) error {
	return s.workflowRepo.UninstallWorkflow(ctx, name)
}

// ExportWorkflow 导出工作流配置
func (s *WorkflowService) ExportWorkflow(ctx context.Context, name string) ([]byte, error) {
	record, err := s.workflowRepo.GetWorkflow(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get workflow: %w", err)
	}

	return []byte(record.ConfigYAML), nil
}

// SyncWorkflows 同步工作流
func (s *WorkflowService) SyncWorkflows(ctx context.Context, yamlDir string) ([]string, error) {
	entries, err := os.ReadDir(yamlDir)
	if err != nil {
		return nil, fmt.Errorf("read workflows directory: %w", err)
	}

	var synced []string

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".yaml")
		yamlPath := filepath.Join(yamlDir, entry.Name())

		data, err := os.ReadFile(yamlPath)
		if err != nil {
			return nil, fmt.Errorf("read workflow file %s: %w", yamlPath, err)
		}

		// 验证工作流
		var workflow types.Workflow
		if err := yaml.Unmarshal(data, &workflow); err != nil {
			return nil, fmt.Errorf("parse workflow YAML %s: %w", yamlPath, err)
		}

		if err := s.loader.Validate(ctx, &workflow); err != nil {
			return nil, fmt.Errorf("validate workflow %s: %w", name, err)
		}

		// 安装或更新
		if err := s.workflowRepo.InstallWorkflow(ctx, name, data); err != nil {
			// 如果已存在，不算错误（跳过）
			if !strings.Contains(err.Error(), "already exists") {
				return nil, fmt.Errorf("sync workflow %s: %w", name, err)
			}
		}

		synced = append(synced, name)
	}

	return synced, nil
}
