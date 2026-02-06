package repository

import (
	"context"
	"fmt"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
)

// workflowRepository 实现 domain.WorkflowRepository 接口
type workflowRepository struct {
	repo *db.WorkflowRepository
}

// NewWorkflowRepository 创建新的工作流仓库
func NewWorkflowRepository(repo *db.WorkflowRepository) domain.WorkflowRepository {
	return &workflowRepository{repo: repo}
}

// GetWorkflow 根据名称获取工作流
func (r *workflowRepository) GetWorkflow(ctx context.Context, name string) (*domain.WorkflowRecord, error) {
	record, err := r.repo.GetWorkflow(name)
	if err != nil {
		if err.Error() == fmt.Sprintf("workflow '%s' not found", name) {
			return nil, fmt.Errorf("get workflow %s: %w", name, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get workflow %s: %w", name, err)
	}
	return &domain.WorkflowRecord{
		Name:         record.Name,
		DisplayName:  record.DisplayName,
		Description:  record.Description,
		WorkflowType: record.WorkflowType,
		Version:      record.Version,
		IsBuiltin:    record.IsBuiltin,
		ConfigYAML:   record.ConfigYAML,
	}, nil
}

// GetWorkflowByType 根据类型获取工作流
func (r *workflowRepository) GetWorkflowByType(ctx context.Context, workflowType string) (*domain.WorkflowRecord, error) {
	record, err := r.repo.GetWorkflowByType(workflowType)
	if err != nil {
		if err.Error() == fmt.Sprintf("workflow of type '%s' not found", workflowType) {
			return nil, fmt.Errorf("get workflow by type %s: %w", workflowType, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get workflow by type %s: %w", workflowType, err)
	}
	return &domain.WorkflowRecord{
		Name:         record.Name,
		DisplayName:  record.DisplayName,
		Description:  record.Description,
		WorkflowType: record.WorkflowType,
		Version:      record.Version,
		IsBuiltin:    record.IsBuiltin,
		ConfigYAML:   record.ConfigYAML,
	}, nil
}

// ListWorkflows 列出所有工作流
func (r *workflowRepository) ListWorkflows(ctx context.Context) ([]*domain.WorkflowRecord, error) {
	records, err := r.repo.ListWorkflows()
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}

	result := make([]*domain.WorkflowRecord, len(records))
	for i, r := range records {
		result[i] = &domain.WorkflowRecord{
			Name:         r.Name,
			DisplayName:  r.DisplayName,
			Description:  r.Description,
			WorkflowType: r.WorkflowType,
			Version:      r.Version,
			IsBuiltin:    r.IsBuiltin,
			ConfigYAML:   r.ConfigYAML,
		}
	}
	return result, nil
}

// InitializeBuiltinWorkflows 初始化内置工作流
func (r *workflowRepository) InitializeBuiltinWorkflows(ctx context.Context, workflowsDir string) error {
	return r.repo.InitializeBuiltinWorkflows(workflowsDir)
}

// InstallWorkflow 安装工作流
func (r *workflowRepository) InstallWorkflow(ctx context.Context, name string, yamlData []byte) error {
	return r.repo.InstallWorkflow(name, yamlData)
}

// UninstallWorkflow 卸载工作流
func (r *workflowRepository) UninstallWorkflow(ctx context.Context, name string) error {
	return r.repo.UninstallWorkflow(name)
}
