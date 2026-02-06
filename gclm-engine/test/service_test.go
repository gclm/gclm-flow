package test

import (
	"context"
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/repository"
	"github.com/gclm/gclm-flow/gclm-engine/internal/service"
	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// TestTaskService 测试任务服务
func TestTaskService(t *testing.T) {
	t.Run("CreateTask", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		// 使用 domain 接口适配器
		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 测试创建任务（指定工作流类型）
		task, err := taskService.CreateTask(ctx, "修复这个bug", "code_simple")
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		if task.WorkflowType != "fix" {
			t.Errorf("Expected workflow type fix, got %s", task.WorkflowType)
		}

		// 验证阶段已创建
		phases, err := taskRepo.GetPhasesByTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to get phases: %v", err)
		}

		if len(phases) == 0 {
			t.Error("Expected at least one phase")
		}
	})

	t.Run("GetExecutionPlan", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 创建任务
		task, err := taskService.CreateTask(ctx, "编写技术文档", "document")
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		// 获取执行计划
		plan, err := taskService.GetExecutionPlan(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to get execution plan: %v", err)
		}

		if plan.WorkflowType != "docs" {
			t.Errorf("Expected workflow type docs, got %s", plan.WorkflowType)
		}

		if len(plan.Steps) == 0 {
			t.Error("Expected at least one step")
		}

		// 验证步骤包含必要信息
		step := plan.Steps[0]
		if step.PhaseID == "" {
			t.Error("Expected phase ID")
		}
		if step.Agent == "" {
			t.Error("Expected agent name")
		}
	})

	t.Run("GetCurrentPhase", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 创建任务
		task, _ := taskService.CreateTask(ctx, "开发新功能", "code_complex")

		// 获取当前阶段
		phase, err := taskService.GetCurrentPhase(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to get current phase: %v", err)
		}

		if phase.Status != types.PhaseStatusPending {
			t.Errorf("Expected status pending, got %s", phase.Status)
		}
	})

	t.Run("ReportPhaseOutput", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 创建任务
		task, _ := taskService.CreateTask(ctx, "简单修复", "code_simple")
		phase, _ := taskService.GetCurrentPhase(ctx, task.ID)

		// 报告输出
		output := "Phase completed successfully"
		err := taskService.ReportPhaseOutput(ctx, task.ID, phase.ID, output)
		if err != nil {
			t.Fatalf("Failed to report phase output: %v", err)
		}

		// 验证阶段状态
		updatedPhase, _ := taskRepo.GetPhase(ctx, phase.ID)
		if updatedPhase.Status != types.PhaseStatusCompleted {
			t.Errorf("Expected status completed, got %s", updatedPhase.Status)
		}

		if updatedPhase.OutputText != output {
			t.Errorf("Expected output '%s', got '%s'", output, updatedPhase.OutputText)
		}
	})

	t.Run("ReportPhaseError", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 创建任务
		task, err := taskService.CreateTask(ctx, "测试任务", "code_simple")
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
		phase, err := taskService.GetCurrentPhase(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to get current phase: %v", err)
		}

		// 报告错误 - 如果是必需阶段，任务会失败并返回错误
		errMsg := "Something went wrong"
		err = taskService.ReportPhaseError(ctx, task.ID, phase.ID, errMsg)

		// 验证阶段状态已更新为失败
		updatedPhase, _ := taskRepo.GetPhase(ctx, phase.ID)
		if updatedPhase.Status != types.PhaseStatusFailed {
			t.Errorf("Expected phase status failed, got %s", updatedPhase.Status)
		}

		// 如果是必需阶段，任务应该失败
		updatedTask, _ := taskRepo.GetTask(ctx, task.ID)
		if updatedTask.Status == types.TaskStatusFailed {
			// 这是预期行为，返回错误是正常的
			if err == nil {
				t.Error("Expected error when required phase fails, got nil")
			}
		} else {
			// 非必需阶段，不应该返回错误
			if err != nil {
				t.Fatalf("Unexpected error for non-required phase: %v", err)
			}
		}
	})

	t.Run("GetTaskStatus", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 创建任务
		task, _ := taskService.CreateTask(ctx, "状态查询测试", "code_simple")

		// 获取任务状态
		status, err := taskService.GetTaskStatus(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to get task status: %v", err)
		}

		if status.TaskID != task.ID {
			t.Errorf("Expected task ID %s, got %s", task.ID, status.TaskID)
		}

		if len(status.Phases) == 0 {
			t.Error("Expected at least one phase")
		}

		if status.WorkflowType != string(task.WorkflowType) {
			t.Errorf("Expected workflow type %s, got %s", task.WorkflowType, status.WorkflowType)
		}
	})

	t.Run("PauseResumeCancelTask", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		parser := workflow.NewParser(getWorkflowsPath(t))

		taskRepo := repository.NewTaskRepository(repo)
		workflowLoader := repository.NewWorkflowLoader(parser)

		taskService := service.NewTaskService(taskRepo, workflowLoader)

		ctx := context.Background()

		// 创建任务
		task, _ := taskService.CreateTask(ctx, "暂停测试", "code_simple")

		// 暂停任务
		err := taskService.PauseTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to pause task: %v", err)
		}

		status, _ := taskService.GetTaskStatus(ctx, task.ID)
		if status.Status != types.TaskStatusPaused {
			t.Errorf("Expected status paused, got %s", status.Status)
		}

		// 恢复任务
		err = taskService.ResumeTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to resume task: %v", err)
		}

		status, _ = taskService.GetTaskStatus(ctx, task.ID)
		if status.Status != types.TaskStatusRunning {
			t.Errorf("Expected status running, got %s", status.Status)
		}

		// 取消任务
		err = taskService.CancelTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to cancel task: %v", err)
		}

		status, _ = taskService.GetTaskStatus(ctx, task.ID)
		if status.Status != types.TaskStatusCancelled {
			t.Errorf("Expected status cancelled, got %s", status.Status)
		}
	})
}

// TestWorkflowTypeDetection 测试工作流类型检测
func TestWorkflowTypeDetection(t *testing.T) {
	database, cleanup := getTestDB(t)
	defer cleanup()

	repo := db.NewRepository(database)
	parser := workflow.NewParser(getWorkflowsPath(t))

	taskRepo := repository.NewTaskRepository(repo)
	workflowLoader := repository.NewWorkflowLoader(parser)

	taskService := service.NewTaskService(taskRepo, workflowLoader)

	ctx := context.Background()

	tests := []struct {
		name             string
		prompt           string
		workflowName     string
		expectedWorkflow string
	}{
		{"文档编写", "编写API文档", "document", "docs"},
		{"方案设计", "设计技术方案", "document", "docs"},
		{"Bug修复", "修复登录bug", "code_simple", "fix"},
		{"小修改", "fix error in auth", "code_simple", "fix"},
		{"新功能", "开发用户管理模块", "code_complex", "feat"},
		{"重构", "重构数据库层", "code_complex", "feat"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := taskService.CreateTask(ctx, tt.prompt, tt.workflowName)
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}

			if task.WorkflowType != tt.expectedWorkflow {
				t.Errorf("Expected workflow type %s, got %s", tt.expectedWorkflow, task.WorkflowType)
			}
		})
	}
}

// TestErrorClassification 测试错误分类
func TestErrorClassification(t *testing.T) {
	classifier := service.NewErrorClassifier()

	tests := []struct {
		name     string
		err      error
		expected service.ErrorType
	}{
		{"ContextCanceled", context.Canceled, service.ErrorTypeCancellation},
		{"ContextDeadline", context.DeadlineExceeded, service.ErrorTypeCancellation},
		{"TimeoutError", &timeoutError{}, service.ErrorTypeTimeout},
		{"ValidationError", &validationError{}, service.ErrorTypeValidation},
		{"TemporaryError", &temporaryError{}, service.ErrorTypeTemporary},
		{"UnknownError", &unknownError{}, service.ErrorTypePermanent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.Classify(tt.err)
			if result != tt.expected {
				t.Errorf("Expected error type %v, got %v", tt.expected, result)
			}
		})
	}

	t.Run("IsRecoverable", func(t *testing.T) {
		if !classifier.IsRecoverable(&temporaryError{}) {
			t.Error("Expected temporary error to be recoverable")
		}

		if classifier.IsRecoverable(&validationError{}) {
			t.Error("Expected validation error to not be recoverable")
		}
	})
}

// 错误类型实现

type timeoutError struct{}

func (e *timeoutError) Error() string { return "operation timeout" }
func (e *timeoutError) Timeout() bool { return true }

type validationError struct{}

func (e *validationError) Error() string { return "validation failed" }
func (e *validationError) Invalid() bool { return true }

type temporaryError struct{}

func (e *temporaryError) Error() string   { return "temporary failure" }
func (e *temporaryError) Temporary() bool { return true }

type unknownError struct{}

func (e *unknownError) Error() string { return "unknown error" }
