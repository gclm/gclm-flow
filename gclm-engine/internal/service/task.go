package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/pipeline"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// TaskService 任务服务 - 大管家核心
// 职责: 工作流配置和状态管理，不执行 Agent
type TaskService struct {
	repo   *db.Repository
	parser *pipeline.Parser
}

// NewTaskService 创建任务服务
func NewTaskService(repo *db.Repository, parser *pipeline.Parser) *TaskService {
	return &TaskService{
		repo:   repo,
		parser: parser,
	}
}

// CreateTask 创建新任务（替代 setup-gclm.sh 的任务创建）
func (s *TaskService) CreateTask(ctx context.Context, prompt string, workflowType string) (*types.Task, error) {
	// 如果未指定工作流类型，自动检测
	if workflowType == "" {
		workflowType = s.detectWorkflowType(prompt)
	}

	// 加载流水线配置
	pipeline, err := s.parser.GetPipelineByWorkflowType(workflowType)
	if err != nil {
		return nil, fmt.Errorf("failed to load pipeline: %w", err)
	}

	// 创建任务
	now := time.Now()
	task := &types.Task{
		ID:           generateID("task"),
		PipelineID:   pipeline.Name,
		Prompt:       prompt,
		WorkflowType: types.WorkflowType(workflowType),
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  len(pipeline.Nodes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}

	// 创建阶段
	if err := s.createPhases(task, pipeline); err != nil {
		return nil, err
	}

	// 记录事件
	s.recordEvent(task.ID, types.EventTypeTaskCreated, "Task created")

	// 标记为运行中
	s.repo.UpdateTaskStatus(task.ID, types.TaskStatusRunning)

	return task, nil
}

// GetExecutionPlan 获取执行计划（供 Skills 查询）
// 这是 setup-gclm.sh 状态文件的 Go 版本
func (s *TaskService) GetExecutionPlan(ctx context.Context, taskID string) (*ExecutionPlan, error) {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return nil, err
	}

	// 加载流水线
	pipe, err := s.parser.LoadPipeline(task.PipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to load pipeline: %w", err)
	}

	// 获取阶段
	phases, err := s.repo.GetPhasesByTask(taskID)
	if err != nil {
		return nil, err
	}

	// 构建执行计划
	plan := &ExecutionPlan{
		TaskID:      taskID,
		PipelineID:  task.PipelineID,
		TotalSteps:  len(phases),
		WorkflowType: string(task.WorkflowType),
		Steps:       make([]*ExecutionStep, 0),
	}

	// 按顺序组织步骤
	for _, phase := range phases {
		step := &ExecutionStep{
			Sequence:    phase.Sequence,
			PhaseID:     phase.ID,
			PhaseName:   phase.PhaseName,
			DisplayName: phase.DisplayName,
			Agent:       phase.AgentName,
			Model:       phase.ModelName,
			Status:      phase.Status,
			Dependencies: make([]string, 0),
		}

		// 找出依赖的阶段
		node := findNode(pipe, phase.PhaseName)
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
	phases, err := s.repo.GetPhasesByTask(taskID)
	if err != nil {
		return nil, err
	}

	// 找到下一个待执行的阶段
	for _, phase := range phases {
		if phase.Status == types.PhaseStatusPending {
			// 检查依赖是否满足
			if s.areDependenciesSatisfied(phase, phases) {
				return phase, nil
			}
		}
	}

	return nil, fmt.Errorf("no pending phase with satisfied dependencies")
}

// ReportPhaseOutput Skills 完成阶段后报告输出
func (s *TaskService) ReportPhaseOutput(ctx context.Context, taskID, phaseID string, output string) error {
	// 更新阶段输出
	if err := s.repo.UpdatePhaseOutput(phaseID, output); err != nil {
		return err
	}

	// 更新阶段状态为完成
	if err := s.repo.UpdatePhaseStatus(phaseID, types.PhaseStatusCompleted); err != nil {
		return err
	}

	// 更新任务进度
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return err
	}

	phases, err := s.repo.GetPhasesByTask(taskID)
	if err != nil {
		return err
	}

	completedCount := 0
	for _, p := range phases {
		if p.Status == types.PhaseStatusCompleted {
			completedCount++
		}

		// 更新当前阶段索引
		if p.Status == types.PhaseStatusCompleted && p.Sequence+1 > task.CurrentPhase {
			s.repo.UpdateTaskProgress(taskID, p.Sequence+1)
		}
	}

	// 检查是否全部完成
	if completedCount >= task.TotalPhases {
		s.repo.CompleteTask(taskID, "All phases completed successfully")
	}

	// 记录事件
	s.recordEvent(taskID, types.EventTypePhaseCompleted, fmt.Sprintf("Phase %s completed", phaseID))

	return nil
}

// ReportPhaseError Skills 报告阶段错误
func (s *TaskService) ReportPhaseError(ctx context.Context, taskID, phaseID string, errMsg string) error {
	// 更新阶段输出
	s.repo.UpdatePhaseOutput(phaseID, fmt.Sprintf("Error: %s", errMsg))

	// 更新阶段状态为失败
	s.repo.UpdatePhaseStatus(phaseID, types.PhaseStatusFailed)

	// 检查是否为必需阶段
	task, _ := s.repo.GetTask(taskID)
	phases, _ := s.repo.GetPhasesByTask(taskID)

	var currentPhase *types.TaskPhase
	for _, p := range phases {
		if p.ID == phaseID {
			currentPhase = p
			break
		}
	}

	if currentPhase != nil {
		pipe, _ := s.parser.LoadPipeline(task.PipelineID)
		node := findNode(pipe, currentPhase.PhaseName)

		if node != nil && node.Required {
			// 必需阶段失败，任务失败
			s.repo.FailTask(taskID, fmt.Sprintf("Required phase %s failed: %s", currentPhase.PhaseName, errMsg))
			return fmt.Errorf("required phase failed")
		}
		// 非必需阶段失败，继续
	}

	return nil
}

// PauseTask 暂停任务
func (s *TaskService) PauseTask(ctx context.Context, taskID string) error {
	return s.repo.UpdateTaskStatus(taskID, types.TaskStatusPaused)
}

// ResumeTask 恢复任务
func (s *TaskService) ResumeTask(ctx context.Context, taskID string) error {
	return s.repo.UpdateTaskStatus(taskID, types.TaskStatusRunning)
}

// CancelTask 取消任务
func (s *TaskService) CancelTask(ctx context.Context, taskID string) error {
	return s.repo.UpdateTaskStatus(taskID, types.TaskStatusCancelled)
}

// GetTaskStatus 获取任务状态（替代 setup-gclm.sh 状态文件读取）
func (s *TaskService) GetTaskStatus(ctx context.Context, taskID string) (*TaskStatusResponse, error) {
	task, err := s.repo.GetTask(taskID)
	if err != nil {
		return nil, err
	}

	phases, err := s.repo.GetPhasesByTask(taskID)
	if err != nil {
		return nil, err
	}

	response := &TaskStatusResponse{
		TaskID:       task.ID,
		Status:       task.Status,
		CurrentPhase: task.CurrentPhase,
		TotalPhases:  task.TotalPhases,
		WorkflowType: string(task.WorkflowType),
		Phases:       make([]*PhaseStatus, 0),
	}

	for _, phase := range phases {
		response.Phases = append(response.Phases, &PhaseStatus{
			PhaseName:   phase.PhaseName,
			DisplayName: phase.DisplayName,
			Status:      phase.Status,
			Agent:       phase.AgentName,
			Model:       phase.ModelName,
			Sequence:    phase.Sequence,
		})
	}

	return response, nil
}

// 内部方法

// detectWorkflowType 检测工作流类型（使用统一分类器）
func (s *TaskService) detectWorkflowType(prompt string) string {
	return DetectWorkflowType(prompt)
}

func (s *TaskService) createPhases(task *types.Task, pipe *types.Pipeline) error {
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

		if err := s.repo.CreatePhase(phase); err != nil {
			return err
		}
	}

	return nil
}

func (s *TaskService) areDependenciesSatisfied(phase *types.TaskPhase, allPhases []*types.TaskPhase) bool {
	task, _ := s.repo.GetTask(phase.TaskID)
	pipe, _ := s.parser.LoadPipeline(task.PipelineID)

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

func (s *TaskService) recordEvent(taskID string, eventType types.EventType, data string) {
	event := &types.Event{
		ID:         generateID("event"),
		TaskID:     taskID,
		EventType:  eventType,
		EventLevel: types.EventLevelInfo,
		Data:       data,
		OccurredAt: time.Now(),
	}
	s.repo.CreateEvent(event)
}

// 辅助函数

func findNode(pipeline *types.Pipeline, ref string) *types.PipelineNode {
	for i := range pipeline.Nodes {
		if pipeline.Nodes[i].Ref == ref {
			return &pipeline.Nodes[i]
		}
	}
	return nil
}

func generateID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// ExecutionPlan 执行计划（替代 setup-gclm.sh 的状态文件）
type ExecutionPlan struct {
	TaskID       string
	PipelineID   string
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

// TaskStatusResponse 任务状态响应
type TaskStatusResponse struct {
	TaskID       string
	Status       types.TaskStatus
	CurrentPhase int
	TotalPhases  int
	WorkflowType string
	Phases       []*PhaseStatus
}

// PhaseStatus 阶段状态
type PhaseStatus struct {
	PhaseName   string
	DisplayName string
	Status      types.PhaseStatus
	Agent       string
	Model       string
	Sequence    int
}
