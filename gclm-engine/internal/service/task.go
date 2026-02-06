package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// TaskService 任务服务 - 大管家核心
// 职责: 工作流配置和状态管理，不执行 Agent
type TaskService struct {
	taskRepo domain.TaskRepository
	loader   domain.WorkflowLoader
}

// NewTaskService 创建任务服务
func NewTaskService(taskRepo domain.TaskRepository, loader domain.WorkflowLoader) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		loader:   loader,
	}
}

// ListTasks 列出所有任务
func (s *TaskService) ListTasks(ctx context.Context, status *types.TaskStatus, limit int) ([]*types.Task, error) {
	tasks, err := s.taskRepo.ListTasks(ctx, status, limit)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}

// CreateTask 创建新任务
// workflow: 必需，指定工作流名称（如 "analyze", "docs", "feat", "fix"）
// prompt: 用户任务描述
func (s *TaskService) CreateTask(ctx context.Context, prompt string, workflow string) (*types.Task, error) {
	// workflow 必需
	if workflow == "" {
		return nil, fmt.Errorf("workflow is required: %w", domain.ErrInvalidInput)
	}

	// 加载工作流配置（按名称）
	workflowDef, err := s.loader.Load(ctx, workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to load workflow '%s': %w", workflow, err)
	}

	// 创建任务
	now := time.Now()
	task := &types.Task{
		ID:           generateID("task"),
		WorkflowID:   workflowDef.Name,
		Prompt:       prompt,
		WorkflowType: workflowDef.WorkflowType,
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  len(workflowDef.Nodes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.taskRepo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	// 创建阶段
	if err := s.createPhases(ctx, task, workflowDef); err != nil {
		return nil, fmt.Errorf("create phases: %w", err)
	}

	// 记录事件
	s.recordEvent(ctx, task.ID, types.EventTypeTaskCreated, "Task created")

	// 标记为运行中
	if err := s.taskRepo.UpdateTaskStatus(ctx, task.ID, types.TaskStatusRunning); err != nil {
		return nil, fmt.Errorf("update task status: %w", err)
	}

	return task, nil
}

// GetExecutionPlan 获取执行计划（供 Skills 查询）
func (s *TaskService) GetExecutionPlan(ctx context.Context, taskID string) (*domain.ExecutionPlan, error) {
	task, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	// 加载工作流
	wf, err := s.loader.Load(ctx, task.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("load workflow: %w", err)
	}

	// 获取阶段
	phases, err := s.taskRepo.GetPhasesByTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get phases: %w", err)
	}

	// 构建执行计划
	plan := &domain.ExecutionPlan{
		TaskID:       taskID,
		WorkflowID:   task.WorkflowID,
		TotalSteps:   len(phases),
		WorkflowType: string(task.WorkflowType),
		Steps:        make([]*domain.ExecutionStep, 0),
	}

	// 按顺序组织步骤
	for _, phase := range phases {
		step := &domain.ExecutionStep{
			Sequence:     phase.Sequence,
			PhaseID:      phase.ID,
			PhaseName:    phase.PhaseName,
			DisplayName:  phase.DisplayName,
			Agent:        phase.AgentName,
			Model:        phase.ModelName,
			Status:       phase.Status,
			Dependencies: make([]string, 0),
		}

		// 找出依赖的阶段
		node := findNode(wf, phase.PhaseName)
		if node != nil {
			step.Dependencies = node.DependsOn
			step.Required = node.Required
			step.Timeout = node.Timeout
		}

		plan.Steps = append(plan.Steps, step)
	}

	return plan, nil
}

// GetCurrentPhase 获取当前应该执行的阶段（供 Skills 查询）
func (s *TaskService) GetCurrentPhase(ctx context.Context, taskID string) (*types.TaskPhase, error) {
	phases, err := s.taskRepo.GetPhasesByTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get phases: %w", err)
	}

	// 找到下一个待执行的阶段
	for _, phase := range phases {
		if phase.Status == types.PhaseStatusPending {
			// 检查依赖是否满足
			if s.areDependenciesSatisfied(ctx, phase, phases) {
				return phase, nil
			}
		}
	}

	return nil, fmt.Errorf("no pending phase with satisfied dependencies")
}

// ReportPhaseOutput Skills 完成阶段后报告输出
func (s *TaskService) ReportPhaseOutput(ctx context.Context, taskID, phaseID, output string) error {
	// 更新阶段输出
	if err := s.taskRepo.UpdatePhaseOutput(ctx, phaseID, output); err != nil {
		return fmt.Errorf("update phase output: %w", err)
	}

	// 更新阶段状态为完成
	if err := s.taskRepo.UpdatePhaseStatus(ctx, phaseID, types.PhaseStatusCompleted); err != nil {
		return fmt.Errorf("update phase status: %w", err)
	}

	// 更新任务进度
	task, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	phases, err := s.taskRepo.GetPhasesByTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get phases: %w", err)
	}

	completedCount := 0
	for _, p := range phases {
		if p.Status == types.PhaseStatusCompleted {
			completedCount++

			// 更新当前阶段索引
			if p.Sequence+1 > task.CurrentPhase {
				if err := s.taskRepo.UpdateTaskProgress(ctx, taskID, p.Sequence+1); err != nil {
					return fmt.Errorf("update task progress: %w", err)
				}
			}
		}
	}

	// 检查是否全部完成
	if completedCount >= task.TotalPhases {
		if err := s.taskRepo.CompleteTask(ctx, taskID, "All phases completed successfully"); err != nil {
			return fmt.Errorf("complete task: %w", err)
		}
	}

	// 记录事件
	s.recordEvent(ctx, taskID, types.EventTypePhaseCompleted, fmt.Sprintf("Phase %s completed", phaseID))

	return nil
}

// ReportPhaseError Skills 报告阶段错误
func (s *TaskService) ReportPhaseError(ctx context.Context, taskID, phaseID, errMsg string) error {
	// 更新阶段输出
	if err := s.taskRepo.UpdatePhaseOutput(ctx, phaseID, fmt.Sprintf("Error: %s", errMsg)); err != nil {
		return fmt.Errorf("update phase output: %w", err)
	}

	// 更新阶段状态为失败
	if err := s.taskRepo.UpdatePhaseStatus(ctx, phaseID, types.PhaseStatusFailed); err != nil {
		return fmt.Errorf("update phase status: %w", err)
	}

	// 检查是否为必需阶段
	task, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	phases, err := s.taskRepo.GetPhasesByTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get phases: %w", err)
	}

	var currentPhase *types.TaskPhase
	for _, p := range phases {
		if p.ID == phaseID {
			currentPhase = p
			break
		}
	}

	if currentPhase != nil {
		pipe, err := s.loader.Load(ctx, task.WorkflowID)
		if err != nil {
			return fmt.Errorf("load workflow: %w", err)
		}

		node := findNode(pipe, currentPhase.PhaseName)
		if node != nil && node.Required {
			// 必需阶段失败，任务失败
			if err := s.taskRepo.FailTask(ctx, taskID, fmt.Sprintf("Required phase %s failed: %s", currentPhase.PhaseName, errMsg)); err != nil {
				return fmt.Errorf("fail task: %w", err)
			}
			return fmt.Errorf("required phase failed: %w", domain.ErrRequiredPhaseFailed)
		}
		// 非必需阶段失败，继续
	}

	return nil
}

// PauseTask 暂停任务
func (s *TaskService) PauseTask(ctx context.Context, taskID string) error {
	return s.taskRepo.UpdateTaskStatus(ctx, taskID, types.TaskStatusPaused)
}

// ResumeTask 恢复任务
func (s *TaskService) ResumeTask(ctx context.Context, taskID string) error {
	return s.taskRepo.UpdateTaskStatus(ctx, taskID, types.TaskStatusRunning)
}

// CancelTask 取消任务
func (s *TaskService) CancelTask(ctx context.Context, taskID string) error {
	return s.taskRepo.UpdateTaskStatus(ctx, taskID, types.TaskStatusCancelled)
}

// GetTaskStatus 获取任务状态（替代 setup-gclm.sh 状态文件读取）
func (s *TaskService) GetTaskStatus(ctx context.Context, taskID string) (*domain.TaskStatusResponse, error) {
	task, err := s.taskRepo.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}

	phases, err := s.taskRepo.GetPhasesByTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get phases: %w", err)
	}

	response := &domain.TaskStatusResponse{
		TaskID:       task.ID,
		Status:       task.Status,
		CurrentPhase: task.CurrentPhase,
		TotalPhases:  task.TotalPhases,
		WorkflowType: string(task.WorkflowType),
		Phases:       make([]*domain.PhaseStatus, 0),
	}

	for _, phase := range phases {
		response.Phases = append(response.Phases, &domain.PhaseStatus{
			PhaseName:   phase.PhaseName,
			DisplayName: phase.DisplayName,
			Status:      phase.Status,
			AgentName:   phase.AgentName,
			ModelName:   phase.ModelName,
			Sequence:    phase.Sequence,
		})
	}

	return response, nil
}

// 内部方法

func (s *TaskService) createPhases(ctx context.Context, task *types.Task, pipe *types.Workflow) error {
	now := time.Now()

	for i, node := range pipe.Nodes {
		phase := &types.TaskPhase{
			ID:          generateID("phase"),
			TaskID:      task.ID,
			PhaseName:   node.Ref,
			DisplayName: node.DisplayName,
			Sequence:    i,
			AgentName:   node.Agent,
			ModelName:   node.Model,
			Status:      types.PhaseStatusPending,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		if err := s.taskRepo.CreatePhase(ctx, phase); err != nil {
			return fmt.Errorf("create phase: %w", err)
		}
	}

	return nil
}

func (s *TaskService) areDependenciesSatisfied(ctx context.Context, phase *types.TaskPhase, allPhases []*types.TaskPhase) bool {
	task, err := s.taskRepo.GetTask(ctx, phase.TaskID)
	if err != nil {
		return false
	}

	pipe, err := s.loader.Load(ctx, task.WorkflowID)
	if err != nil {
		return false
	}

	node := findNode(pipe, phase.PhaseName)
	if node == nil || len(node.DependsOn) == 0 {
		return true
	}

	// 检查所有依赖是否都已完成
	for _, depRef := range node.DependsOn {
		found := false
		for _, p := range allPhases {
			if p.PhaseName == depRef && p.Status == types.PhaseStatusCompleted {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

func (s *TaskService) recordEvent(ctx context.Context, taskID string, eventType types.EventType, data string) {
	event := &types.Event{
		ID:         generateID("event"),
		TaskID:     taskID,
		EventType:  eventType,
		EventLevel: types.EventLevelInfo,
		Data:       data,
		OccurredAt: time.Now(),
	}
	s.taskRepo.CreateEvent(ctx, event)
}

// 辅助函数

func findNode(workflow *types.Workflow, ref string) *types.WorkflowNode {
	for i := range workflow.Nodes {
		if workflow.Nodes[i].Ref == ref {
			return &workflow.Nodes[i]
		}
	}
	return nil
}

func generateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
