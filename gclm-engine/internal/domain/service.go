package domain

import (
	"context"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// ============================================================================
// TaskService 任务服务接口
// ============================================================================

// TaskService 定义任务管理业务逻辑
type TaskService interface {
	// CreateTask 创建新任务
	// prompt: 用户任务描述
	// workflowType: 工作流类型（如 "feat", "fix", "docs"）
	CreateTask(ctx context.Context, prompt string, workflowType string) (*types.Task, error)

	// ListTasks 列出所有任务
	ListTasks(ctx context.Context, status *types.TaskStatus, limit int) ([]*types.Task, error)

	// GetTaskStatus 获取任务状态
	GetTaskStatus(ctx context.Context, taskID string) (*TaskStatusResponse, error)

	// GetCurrentPhase 获取当前待执行的阶段
	GetCurrentPhase(ctx context.Context, taskID string) (*types.TaskPhase, error)

	// GetExecutionPlan 获取任务执行计划
	GetExecutionPlan(ctx context.Context, taskID string) (*ExecutionPlan, error)

	// ReportPhaseOutput 报告阶段输出
	ReportPhaseOutput(ctx context.Context, taskID, phaseID, output string) error

	// ReportPhaseError 报告阶段错误
	ReportPhaseError(ctx context.Context, taskID, phaseID, errMsg string) error

	// PauseTask 暂停任务
	PauseTask(ctx context.Context, taskID string) error

	// ResumeTask 恢复任务
	ResumeTask(ctx context.Context, taskID string) error

	// CancelTask 取消任务
	CancelTask(ctx context.Context, taskID string) error
}

// ============================================================================
// WorkflowService 工作流服务接口
// ============================================================================

// WorkflowService 定义工作流管理业务逻辑
type WorkflowService interface {
	// ListWorkflows 列出所有工作流
	ListWorkflows(ctx context.Context) ([]*WorkflowInfo, error)

	// GetWorkflow 获取工作流详情
	GetWorkflow(ctx context.Context, name string) (*WorkflowDetail, error)

	// GetWorkflowByType 按类型获取工作流详情
	GetWorkflowByType(ctx context.Context, workflowType string) (*WorkflowDetail, error)

	// ValidateWorkflow 验证工作流配置
	ValidateWorkflow(ctx context.Context, yamlFile string) (*types.Workflow, error)

	// InstallWorkflow 安装工作流
	InstallWorkflow(ctx context.Context, name string, yamlFile string) error

	// UninstallWorkflow 卸载工作流
	UninstallWorkflow(ctx context.Context, name string) error

	// ExportWorkflow 导出工作流配置
	ExportWorkflow(ctx context.Context, name string) ([]byte, error)

	// SyncWorkflows 同步工作流
	SyncWorkflows(ctx context.Context, yamlDir string) ([]string, error)
}

// ============================================================================
// 数据传输对象 (DTO)
// ============================================================================

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	TaskID       string         `json:"taskId"`
	Status       types.TaskStatus `json:"status"`
	CurrentPhase int            `json:"currentPhase"`
	TotalPhases  int            `json:"totalPhases"`
	WorkflowType string         `json:"workflowType"`
	Phases       []*PhaseStatus `json:"phases"`
}

// PhaseStatus 阶段状态
type PhaseStatus struct {
	PhaseName   string          `json:"phaseName"`
	DisplayName string          `json:"displayName"`
	Status      types.PhaseStatus `json:"status"`
	AgentName   string          `json:"agentName"`
	ModelName   string          `json:"modelName"`
	Sequence    int             `json:"sequence"`
}

// ExecutionPlan 执行计划
type ExecutionPlan struct {
	TaskID       string
	WorkflowID   string
	WorkflowType string
	TotalSteps   int
	Steps        []*ExecutionStep
}

// ExecutionStep 执行步骤
type ExecutionStep struct {
	Sequence     int
	PhaseID      string
	PhaseName    string
	DisplayName  string
	Agent        string
	Model        string
	Status       types.PhaseStatus
	Dependencies []string
	Required     bool
	Timeout      int
}

// WorkflowDetail 工作流详情
type WorkflowDetail struct {
	Name         string       `json:"name"`
	DisplayName  string       `json:"displayName"`
	Description  string       `json:"description"`
	Version      string       `json:"version"`
	WorkflowType string       `json:"workflowType"`
	IsBuiltin    bool         `json:"isBuiltin"`
	Nodes        []*NodeDetail `json:"nodes"`
	ConfigYAML   string       `json:"configYaml,omitempty"` // YAML 源文件内容
}

// NodeDetail 节点详情
type NodeDetail struct {
	Ref          string   `json:"ref"`
	DisplayName  string   `json:"displayName"`
	Agent        string   `json:"agent"`
	Model        string   `json:"model"`
	Timeout      int      `json:"timeout"`
	Required     bool     `json:"required"`
	Dependencies []string `json:"dependsOn,omitempty"`
}

// WorkflowInfo 工作流基本信息
type WorkflowInfo struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Description  string `json:"description"`
	WorkflowType string `json:"workflowType"`
	Version      string `json:"version"`
}
