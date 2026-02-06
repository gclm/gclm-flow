package repository

import (
	"context"
	"testing"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// getTestDB creates an in-memory database for testing
func getTestDB(t *testing.T) (*db.Database, func()) {
	cfg := &db.Config{
		Path: ":memory:",
	}

	database, err := db.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	cleanup := func() {
		database.Close()
	}

	return database, cleanup
}

// TestTaskRepository 测试 TaskRepository 适配器
func TestTaskRepository(t *testing.T) {
	t.Run("CreateAndGetTask", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-task-1",
			WorkflowID:   "code_simple",
			Prompt:       "Test prompt",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}

		err := taskRepo.CreateTask(ctx, task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		// Get task
		retrieved, err := taskRepo.GetTask(ctx, task.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if retrieved.Prompt != "Test prompt" {
			t.Errorf("Expected prompt 'Test prompt', got '%s'", retrieved.Prompt)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		// Create test tasks
		for i := 0; i < 3; i++ {
			task := &types.Task{
				ID:           "test-list-" + string(rune('a'+i)),
				WorkflowID:   "code_simple",
				Prompt:       "List test",
				WorkflowType: "fix",
				Status:       types.TaskStatusCreated,
				CurrentPhase: 0,
				TotalPhases:  5,
			}
			taskRepo.CreateTask(ctx, task)
		}

		// List tasks
		tasks, err := taskRepo.ListTasks(ctx, nil, 0)
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) < 3 {
			t.Errorf("Expected at least 3 tasks, got %d", len(tasks))
		}
	})

	t.Run("UpdateTaskStatus", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-status-1",
			WorkflowID:   "code_simple",
			Prompt:       "Status test",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}

		taskRepo.CreateTask(ctx, task)

		err := taskRepo.UpdateTaskStatus(ctx, task.ID, types.TaskStatusRunning)
		if err != nil {
			t.Fatalf("Failed to update task status: %v", err)
		}

		retrieved, _ := taskRepo.GetTask(ctx, task.ID)
		if retrieved.Status != types.TaskStatusRunning {
			t.Errorf("Expected status 'running', got '%s'", retrieved.Status)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-complete-1",
			WorkflowID:   "code_simple",
			Prompt:       "Complete test",
			WorkflowType: "fix",
			Status:       types.TaskStatusRunning,
			CurrentPhase: 5,
			TotalPhases:  5,
		}

		taskRepo.CreateTask(ctx, task)

		result := "Task completed successfully"
		err := taskRepo.CompleteTask(ctx, task.ID, result)
		if err != nil {
			t.Fatalf("Failed to complete task: %v", err)
		}

		retrieved, _ := taskRepo.GetTask(ctx, task.ID)
		if retrieved.Status != types.TaskStatusCompleted {
			t.Errorf("Expected status 'completed', got '%s'", retrieved.Status)
		}
		if retrieved.Result != result {
			t.Errorf("Expected result '%s', got '%s'", result, retrieved.Result)
		}
	})

	t.Run("FailTask", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-fail-1",
			WorkflowID:   "code_simple",
			Prompt:       "Fail test",
			WorkflowType: "fix",
			Status:       types.TaskStatusRunning,
			CurrentPhase: 2,
			TotalPhases:  5,
		}

		taskRepo.CreateTask(ctx, task)

		errMsg := "Task failed due to error"
		err := taskRepo.FailTask(ctx, task.ID, errMsg)
		if err != nil {
			t.Fatalf("Failed to fail task: %v", err)
		}

		retrieved, _ := taskRepo.GetTask(ctx, task.ID)
		if retrieved.Status != types.TaskStatusFailed {
			t.Errorf("Expected status 'failed', got '%s'", retrieved.Status)
		}
		if retrieved.Error != errMsg {
			t.Errorf("Expected error '%s', got '%s'", errMsg, retrieved.Error)
		}
	})

	t.Run("CreatePhase", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		// Create task first
		task := &types.Task{
			ID:           "test-phase-task-1",
			WorkflowID:   "code_simple",
			Prompt:       "Phase test",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		taskRepo.CreateTask(ctx, task)

		// Create phase
		phase := &types.TaskPhase{
			ID:          "test-phase-1",
			TaskID:      "test-phase-task-1",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}

		err := taskRepo.CreatePhase(ctx, phase)
		if err != nil {
			t.Fatalf("Failed to create phase: %v", err)
		}

		// Verify phase was created
		phases, err := taskRepo.GetPhasesByTask(ctx, "test-phase-task-1")
		if err != nil {
			t.Fatalf("Failed to get phases: %v", err)
		}

		if len(phases) != 1 {
			t.Errorf("Expected 1 phase, got %d", len(phases))
		}

		if phases[0].PhaseName != "discovery" {
			t.Errorf("Expected phase name 'discovery', got '%s'", phases[0].PhaseName)
		}
	})

	t.Run("GetPhase", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-get-phase-1",
			WorkflowID:   "code_simple",
			Prompt:       "Get phase test",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		taskRepo.CreateTask(ctx, task)

		phase := &types.TaskPhase{
			ID:          "test-phase-get-1",
			TaskID:      "test-get-phase-1",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}
		taskRepo.CreatePhase(ctx, phase)

		// Get phase
		retrieved, err := taskRepo.GetPhase(ctx, phase.ID)
		if err != nil {
			t.Fatalf("Failed to get phase: %v", err)
		}

		if retrieved.ID != "test-phase-get-1" {
			t.Errorf("Expected ID 'test-phase-get-1', got '%s'", retrieved.ID)
		}
	})

	t.Run("UpdatePhaseStatus", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-phase-status-1",
			WorkflowID:   "code_simple",
			Prompt:       "Phase status test",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		taskRepo.CreateTask(ctx, task)

		phase := &types.TaskPhase{
			ID:          "test-phase-status-1",
			TaskID:      "test-phase-status-1",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}
		taskRepo.CreatePhase(ctx, phase)

		// Update status
		err := taskRepo.UpdatePhaseStatus(ctx, phase.ID, types.PhaseStatusRunning)
		if err != nil {
			t.Fatalf("Failed to update phase status: %v", err)
		}

		// Verify update
		phases, _ := taskRepo.GetPhasesByTask(ctx, "test-phase-status-1")
		if phases[0].Status != types.PhaseStatusRunning {
			t.Errorf("Expected status 'running', got '%s'", phases[0].Status)
		}
	})

	t.Run("UpdatePhaseOutput", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-phase-output-1",
			WorkflowID:   "code_simple",
			Prompt:       "Phase output test",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		taskRepo.CreateTask(ctx, task)

		phase := &types.TaskPhase{
			ID:          "test-phase-output-1",
			TaskID:      "test-phase-output-1",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusRunning,
		}
		taskRepo.CreatePhase(ctx, phase)

		// Update output
		output := "Phase completed successfully"
		err := taskRepo.UpdatePhaseOutput(ctx, phase.ID, output)
		if err != nil {
			t.Fatalf("Failed to update phase output: %v", err)
		}

		// Verify update
		phases, _ := taskRepo.GetPhasesByTask(ctx, "test-phase-output-1")
		if phases[0].OutputText != output {
			t.Errorf("Expected output '%s', got '%s'", output, phases[0].OutputText)
		}
	})

	t.Run("CreateEvent", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-event-task-1",
			WorkflowID:   "code_simple",
			Prompt:       "Event test",
			WorkflowType: "fix",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		taskRepo.CreateTask(ctx, task)

		event := &types.Event{
			ID:         "test-event-1",
			TaskID:     "test-event-task-1",
			EventType:  types.EventTypeTaskCreated,
			EventLevel: types.EventLevelInfo,
		}

		err := taskRepo.CreateEvent(ctx, event)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}

		// Verify event was created
		events, err := taskRepo.GetEventsByTask(ctx, "test-event-task-1", 0)
		if err != nil {
			t.Fatalf("Failed to get events: %v", err)
		}

		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}

		if events[0].EventType != types.EventTypeTaskCreated {
			t.Errorf("Expected event type 'task_created', got '%s'", events[0].EventType)
		}
	})

	t.Run("UpdateTaskProgress", func(t *testing.T) {
		database, cleanup := getTestDB(t)
		defer cleanup()

		repo := db.NewRepository(database)
		taskRepo := NewTaskRepository(repo)

		ctx := context.Background()

		task := &types.Task{
			ID:           "test-progress-1",
			WorkflowID:   "code_simple",
			Prompt:       "Progress test",
			WorkflowType: "fix",
			Status:       types.TaskStatusRunning,
			CurrentPhase: 0,
			TotalPhases:  5,
		}

		taskRepo.CreateTask(ctx, task)

		err := taskRepo.UpdateTaskProgress(ctx, task.ID, 2)
		if err != nil {
			t.Fatalf("Failed to update task progress: %v", err)
		}

		retrieved, _ := taskRepo.GetTask(ctx, task.ID)
		if retrieved.CurrentPhase != 2 {
			t.Errorf("Expected current phase 2, got %d", retrieved.CurrentPhase)
		}
	})
}
