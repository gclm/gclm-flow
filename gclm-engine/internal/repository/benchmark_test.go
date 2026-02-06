package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/internal/workflow"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// getBenchmarkDB 创建基准测试数据库
func getBenchmarkDB() *db.Database {
	cfg := &db.Config{
		Path: ":memory:",
	}
	database, err := db.New(cfg)
	if err != nil {
		panic(err)
	}
	return database
}

// BenchmarkTaskRepository_CreateTask 基准测试：创建任务
func BenchmarkTaskRepository_CreateTask(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &types.Task{
			ID:           "bench-task-" + string(rune('a'+i%26)),
			WorkflowID:   "code_simple",
			Prompt:       "Benchmark test task",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}
		_ = taskRepo.CreateTask(ctx, task)
	}
}

// BenchmarkTaskRepository_GetTask 基准测试：获取任务
func BenchmarkTaskRepository_GetTask(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	// 创建测试数据
	task := &types.Task{
		ID:           "bench-get-task",
		WorkflowID:   "code_simple",
		Prompt:       "Benchmark test",
		WorkflowType: "fix",
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  5,
	}
	taskRepo.CreateTask(ctx, task)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = taskRepo.GetTask(ctx, "bench-get-task")
	}
}

// BenchmarkTaskRepository_ListTasks 基准测试：列出任务
func BenchmarkTaskRepository_ListTasks(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	// 创建测试数据
	for i := 0; i < 100; i++ {
		task := &types.Task{
			ID:           "list-task-" + string(rune('a'+i%26)),
			WorkflowID:   "code_simple",
			Prompt:       "List benchmark",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}
		taskRepo.CreateTask(ctx, task)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = taskRepo.ListTasks(ctx, nil, 0)
	}
}

// BenchmarkTaskRepository_UpdateTaskStatus 基准测试：更新任务状态
func BenchmarkTaskRepository_UpdateTaskStatus(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	task := &types.Task{
		ID:           "bench-update-task",
		WorkflowID:   "code_simple",
		Prompt:       "Update benchmark",
		WorkflowType: "fix",
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  5,
	}
	taskRepo.CreateTask(ctx, task)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = taskRepo.UpdateTaskStatus(ctx, "bench-update-task", types.TaskStatusRunning)
	}
}

// BenchmarkTaskRepository_CreatePhase 基准测试：创建阶段
func BenchmarkTaskRepository_CreatePhase(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	// 创建任务
	task := &types.Task{
		ID:           "bench-phase-task",
		WorkflowID:   "code_simple",
		Prompt:       "Phase benchmark",
		WorkflowType: "fix",
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  5,
	}
	taskRepo.CreateTask(ctx, task)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		phase := &types.TaskPhase{
			ID:          "bench-phase-" + string(rune('a'+i%26)),
			TaskID:      "bench-phase-task",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    i % 5,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}
		_ = taskRepo.CreatePhase(ctx, phase)
	}
}

// BenchmarkTaskRepository_GetPhasesByTask 基准测试：获取任务的所有阶段
func BenchmarkTaskRepository_GetPhasesByTask(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	// 创建任务和阶段
	task := &types.Task{
		ID:           "bench-get-phases",
		WorkflowID:   "code_simple",
		Prompt:       "Get phases benchmark",
		WorkflowType: "fix",
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  5,
	}
	taskRepo.CreateTask(ctx, task)

	for i := 0; i < 5; i++ {
		phase := &types.TaskPhase{
			ID:          "phase-" + string(rune('a'+i)),
			TaskID:      "bench-get-phases",
			PhaseName:   "phase",
			DisplayName: "Phase",
			Sequence:    i,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}
		taskRepo.CreatePhase(ctx, phase)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = taskRepo.GetPhasesByTask(ctx, "bench-get-phases")
	}
}

// BenchmarkTaskRepository_CreateEvent 基准测试：创建事件
func BenchmarkTaskRepository_CreateEvent(b *testing.B) {
	database := getBenchmarkDB()
	defer database.Close()

	repo := db.NewRepository(database)
	taskRepo := NewTaskRepository(repo)
	ctx := context.Background()

	// 创建任务
	task := &types.Task{
		ID:           "bench-event-task",
		WorkflowID:   "code_simple",
		Prompt:       "Event benchmark",
		WorkflowType: "fix",
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  1,
	}
	taskRepo.CreateTask(ctx, task)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := &types.Event{
			ID:         "bench-event-" + string(rune('a'+i%26)),
			TaskID:     "bench-event-task",
			EventType:  types.EventTypeTaskCreated,
			EventLevel: types.EventLevelInfo,
		}
		_ = taskRepo.CreateEvent(ctx, event)
	}
}

// BenchmarkWorkflowLoader_Load 基准测试：加载工作流（含缓存）
func BenchmarkWorkflowLoader_Load_Cached(b *testing.B) {
	loader := setupTestLoader()
	ctx := context.Background()

	// 预热缓存
	loader.Load(ctx, "code_simple")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.Load(ctx, "code_simple")
	}
}

// BenchmarkWorkflowLoader_LoadAll 基准测试：加载所有工作流
func BenchmarkWorkflowLoader_LoadAll(b *testing.B) {
	loader := setupTestLoader()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.LoadAll(ctx)
	}
}

// BenchmarkWorkflowLoader_GetExecutionOrder 基准测试：计算执行顺序
func BenchmarkWorkflowLoader_GetExecutionOrder(b *testing.B) {
	loader := setupTestLoader()
	ctx := context.Background()

	// 预加载工作流
	workflow, _ := loader.Load(ctx, "code_simple")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loader.GetExecutionOrder(ctx, workflow)
	}
}

// setupTestLoader 创建测试用的工作流加载器
func setupTestLoader() domain.WorkflowLoader {
	workflowsDir := getWorkflowsDirForBench()
	parser := workflow.NewParser(workflowsDir)
	return NewWorkflowLoader(parser)
}

// getWorkflowsDirForBench 找到工作流目录
func getWorkflowsDirForBench() string {
	wd, _ := os.Getwd()
	for i := 0; i < 4; i++ {
		testPath := filepath.Join(wd, "workflows")
		if info, err := os.Stat(testPath); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(testPath)
			if len(entries) > 0 {
				return testPath
			}
		}
		wd = filepath.Dir(wd)
		if filepath.Base(wd) == "gclm-engine" {
			return filepath.Join(wd, "workflows")
		}
	}
	return "workflows"
}
