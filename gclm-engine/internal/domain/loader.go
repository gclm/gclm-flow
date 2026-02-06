package domain

import (
	"context"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// ============================================================================
// WorkflowLoader 工作流加载器接口
// ============================================================================

// WorkflowLoader 定义工作流配置加载操作
type WorkflowLoader interface {
	// Load 加载工作流配置
	// name 是工作流名称（如 "feat", "fix", "docs"）
	Load(ctx context.Context, name string) (*types.Workflow, error)

	// LoadAll 加载所有可用工作流
	LoadAll(ctx context.Context) (map[string]*types.Workflow, error)

	// Validate 验证工作流配置是否有效
	Validate(ctx context.Context, workflow *types.Workflow) error

	// GetExecutionOrder 计算工作流的执行顺序
	// 返回按依赖关系排序的节点列表
	GetExecutionOrder(ctx context.Context, workflow *types.Workflow) ([]*types.NodeExecutionOrder, error)
}
