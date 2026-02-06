package domain

import (
	"context"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// ============================================================================
// TaskRepository 任务仓库接口
// ============================================================================

// TaskRepository 定义任务数据访问操作
type TaskRepository interface {
	// CreateTask 创建新任务
	CreateTask(ctx context.Context, task *types.Task) error

	// GetTask 根据 ID 获取任务
	GetTask(ctx context.Context, id string) (*types.Task, error)

	// ListTasks 列出任务，支持按状态过滤和限制数量
	ListTasks(ctx context.Context, status *types.TaskStatus, limit int) ([]*types.Task, error)

	// UpdateTaskStatus 更新任务状态
	UpdateTaskStatus(ctx context.Context, id string, status types.TaskStatus) error

	// UpdateTaskProgress 更新任务当前阶段进度
	UpdateTaskProgress(ctx context.Context, id string, currentPhase int) error

	// CompleteTask 标记任务为完成
	CompleteTask(ctx context.Context, id string, result string) error

	// FailTask 标记任务为失败
	FailTask(ctx context.Context, id string, errMsg string) error

	// CreatePhase 创建任务阶段
	CreatePhase(ctx context.Context, phase *types.TaskPhase) error

	// GetPhase 根据 ID 获取阶段
	GetPhase(ctx context.Context, id string) (*types.TaskPhase, error)

	// GetPhasesByTask 获取任务的所有阶段
	GetPhasesByTask(ctx context.Context, taskID string) ([]*types.TaskPhase, error)

	// UpdatePhaseStatus 更新阶段状态
	UpdatePhaseStatus(ctx context.Context, id string, status types.PhaseStatus) error

	// UpdatePhaseOutput 更新阶段输出
	UpdatePhaseOutput(ctx context.Context, id string, output string) error

	// CreateEvent 创建事件记录
	CreateEvent(ctx context.Context, event *types.Event) error

	// GetEventsByTask 获取任务的事件日志
	GetEventsByTask(ctx context.Context, taskID string, limit int) ([]*types.Event, error)
}

// ============================================================================
// WorkflowRepository 工作流仓库接口
// ============================================================================

// WorkflowRecord 工作流记录
type WorkflowRecord struct {
	Name         string
	DisplayName  string
	Description  string
	WorkflowType string
	Version      string
	IsBuiltin    bool
	ConfigYAML   string
}

// WorkflowRepository 定义工作流数据访问操作
type WorkflowRepository interface {
	// GetWorkflow 根据名称获取工作流
	GetWorkflow(ctx context.Context, name string) (*WorkflowRecord, error)

	// GetWorkflowByType 根据类型获取工作流
	GetWorkflowByType(ctx context.Context, workflowType string) (*WorkflowRecord, error)

	// ListWorkflows 列出所有工作流
	ListWorkflows(ctx context.Context) ([]*WorkflowRecord, error)

	// InitializeBuiltinWorkflows 初始化内置工作流
	InitializeBuiltinWorkflows(ctx context.Context, workflowsDir string) error

	// InstallWorkflow 安装工作流
	InstallWorkflow(ctx context.Context, name string, yamlData []byte) error

	// UninstallWorkflow 卸载工作流
	UninstallWorkflow(ctx context.Context, name string) error
}
