package repository

import (
	"context"
	"fmt"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// taskRepository 实现 domain.TaskRepository 接口
type taskRepository struct {
	repo *db.Repository
}

// NewTaskRepository 创建新的任务仓库
func NewTaskRepository(repo *db.Repository) domain.TaskRepository {
	return &taskRepository{repo: repo}
}

// CreateTask 创建新任务
func (r *taskRepository) CreateTask(ctx context.Context, task *types.Task) error {
	return r.repo.CreateTask(task)
}

// GetTask 根据 ID 获取任务
func (r *taskRepository) GetTask(ctx context.Context, id string) (*types.Task, error) {
	task, err := r.repo.GetTask(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("task not found: %s", id) {
			return nil, fmt.Errorf("get task %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get task %s: %w", id, err)
	}
	return task, nil
}

// ListTasks 列出任务，支持按状态过滤和限制数量
func (r *taskRepository) ListTasks(ctx context.Context, status *types.TaskStatus, limit int) ([]*types.Task, error) {
	return r.repo.ListTasks(status, limit)
}

// UpdateTaskStatus 更新任务状态
func (r *taskRepository) UpdateTaskStatus(ctx context.Context, id string, status types.TaskStatus) error {
	return r.repo.UpdateTaskStatus(id, status)
}

// UpdateTaskProgress 更新任务当前阶段进度
func (r *taskRepository) UpdateTaskProgress(ctx context.Context, id string, currentPhase int) error {
	return r.repo.UpdateTaskProgress(id, currentPhase)
}

// CompleteTask 标记任务为完成
func (r *taskRepository) CompleteTask(ctx context.Context, id string, result string) error {
	return r.repo.CompleteTask(id, result)
}

// FailTask 标记任务为失败
func (r *taskRepository) FailTask(ctx context.Context, id string, errMsg string) error {
	return r.repo.FailTask(id, errMsg)
}

// CreatePhase 创建任务阶段
func (r *taskRepository) CreatePhase(ctx context.Context, phase *types.TaskPhase) error {
	return r.repo.CreatePhase(phase)
}

// GetPhase 根据 ID 获取阶段
func (r *taskRepository) GetPhase(ctx context.Context, id string) (*types.TaskPhase, error) {
	phase, err := r.repo.GetPhase(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("phase not found: %s", id) {
			return nil, fmt.Errorf("get phase %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get phase %s: %w", id, err)
	}
	return phase, nil
}

// GetPhasesByTask 获取任务的所有阶段
func (r *taskRepository) GetPhasesByTask(ctx context.Context, taskID string) ([]*types.TaskPhase, error) {
	return r.repo.GetPhasesByTask(taskID)
}

// UpdatePhaseStatus 更新阶段状态
func (r *taskRepository) UpdatePhaseStatus(ctx context.Context, id string, status types.PhaseStatus) error {
	return r.repo.UpdatePhaseStatus(id, status)
}

// UpdatePhaseOutput 更新阶段输出
func (r *taskRepository) UpdatePhaseOutput(ctx context.Context, id string, output string) error {
	return r.repo.UpdatePhaseOutput(id, output)
}

// CreateEvent 创建事件记录
func (r *taskRepository) CreateEvent(ctx context.Context, event *types.Event) error {
	return r.repo.CreateEvent(event)
}

// GetEventsByTask 获取任务的事件日志
func (r *taskRepository) GetEventsByTask(ctx context.Context, taskID string, limit int) ([]*types.Event, error) {
	return r.repo.GetEventsByTask(taskID, limit)
}
