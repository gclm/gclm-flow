package test

import (
	"os"
	"testing"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/internal/db"
	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
)

// getTestDB returns a test database instance
func getTestDB(t *testing.T) (*db.Database, func()) {
	// Use in-memory database for testing
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

// TestDatabaseCreation 测试数据库创建
func TestDatabaseCreation(t *testing.T) {
	database, cleanup := getTestDB(t)
	defer cleanup()

	// Check if we can query the database
	rows, err := database.GetDB().Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}
		tables = append(tables, name)
	}

	// Check that required tables exist
	requiredTables := []string{"tasks", "task_phases", "events"}
	for _, required := range requiredTables {
		found := false
		for _, table := range tables {
			if table == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Required table %s not found", required)
		}
	}
}

// TestTaskCRUD 测试任务 CRUD 操作
func TestTaskCRUD(t *testing.T) {
	database, cleanup := getTestDB(t)
	defer cleanup()

	repo := db.NewRepository(database)

	t.Run("CreateTask", func(t *testing.T) {
		task := &types.Task{
			ID:           "test-task-1",
			WorkflowID:   "code_simple",
			Prompt:       "Test prompt",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}

		err := repo.CreateTask(task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		// Verify task was created
		retrieved, err := repo.GetTask("test-task-1")
		if err != nil {
			t.Fatalf("Failed to retrieve task: %v", err)
		}

		if retrieved.Prompt != "Test prompt" {
			t.Errorf("Expected prompt 'Test prompt', got '%s'", retrieved.Prompt)
		}
	})

	t.Run("GetTask", func(t *testing.T) {
		task := &types.Task{
			ID:           "test-task-2",
			WorkflowID:   "code_simple",
			Prompt:       "Get test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}

		repo.CreateTask(task)

		retrieved, err := repo.GetTask("test-task-2")
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if retrieved.ID != "test-task-2" {
			t.Errorf("Expected ID 'test-task-2', got '%s'", retrieved.ID)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		// Create multiple tasks
		for i := 0; i < 3; i++ {
			task := &types.Task{
				ID:           "test-list-" + string(rune('a'+i)),
				WorkflowID:   "code_simple",
				Prompt:       "List test",
				WorkflowType: "CODE_SIMPLE",
				Status:       types.TaskStatusCreated,
				CurrentPhase: 0,
				TotalPhases:  5,
			}
			repo.CreateTask(task)
		}

		tasks, err := repo.ListTasks(nil, 0)
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) < 3 {
			t.Errorf("Expected at least 3 tasks, got %d", len(tasks))
		}
	})

	t.Run("UpdateTaskStatus", func(t *testing.T) {
		task := &types.Task{
			ID:           "test-status-1",
			WorkflowID:   "code_simple",
			Prompt:       "Status test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  5,
		}

		repo.CreateTask(task)

		err := repo.UpdateTaskStatus("test-status-1", types.TaskStatusRunning)
		if err != nil {
			t.Fatalf("Failed to update task status: %v", err)
		}

		retrieved, _ := repo.GetTask("test-status-1")
		if retrieved.Status != types.TaskStatusRunning {
			t.Errorf("Expected status 'running', got '%s'", retrieved.Status)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		task := &types.Task{
			ID:           "test-complete-1",
			WorkflowID:   "code_simple",
			Prompt:       "Complete test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusRunning,
			CurrentPhase: 5,
			TotalPhases:  5,
		}

		repo.CreateTask(task)

		result := "Task completed successfully"
		err := repo.CompleteTask("test-complete-1", result)
		if err != nil {
			t.Fatalf("Failed to complete task: %v", err)
		}

		retrieved, _ := repo.GetTask("test-complete-1")
		if retrieved.Status != types.TaskStatusCompleted {
			t.Errorf("Expected status 'completed', got '%s'", retrieved.Status)
		}
		if retrieved.Result != result {
			t.Errorf("Expected result '%s', got '%s'", result, retrieved.Result)
		}
	})
}

// TestTaskPhases 测试任务阶段操作
func TestTaskPhases(t *testing.T) {
	database, cleanup := getTestDB(t)
	defer cleanup()

	repo := db.NewRepository(database)

	t.Run("CreatePhase", func(t *testing.T) {
		// Create a task first
		task := &types.Task{
			ID:           "test-phase-task-1",
			WorkflowID:   "code_simple",
			Prompt:       "Phase test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		repo.CreateTask(task)

		// Create a phase
		phase := &types.TaskPhase{
			ID:          "test-phase-1",
			TaskID:      "test-phase-task-1",
			PhaseName:   "discovery",
			DisplayName: "Discovery / 需求发现",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}

		err := repo.CreatePhase(phase)
		if err != nil {
			t.Fatalf("Failed to create phase: %v", err)
		}

		// Verify phase was created
		phases, err := repo.GetPhasesByTask("test-phase-task-1")
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

	t.Run("UpdatePhaseStatus", func(t *testing.T) {
		// Create task and phase
		task := &types.Task{
			ID:           "test-phase-task-2",
			WorkflowID:   "code_simple",
			Prompt:       "Phase status test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		repo.CreateTask(task)

		phase := &types.TaskPhase{
			ID:          "test-phase-2",
			TaskID:      "test-phase-task-2",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusPending,
		}
		repo.CreatePhase(phase)

		// Update status
		err := repo.UpdatePhaseStatus("test-phase-2", types.PhaseStatusRunning)
		if err != nil {
			t.Fatalf("Failed to update phase status: %v", err)
		}

		// Verify update
		phases, _ := repo.GetPhasesByTask("test-phase-task-2")
		if phases[0].Status != types.PhaseStatusRunning {
			t.Errorf("Expected status 'running', got '%s'", phases[0].Status)
		}
	})

	t.Run("UpdatePhaseOutput", func(t *testing.T) {
		// Create task and phase
		task := &types.Task{
			ID:           "test-phase-task-3",
			WorkflowID:   "code_simple",
			Prompt:       "Phase output test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
		}
		repo.CreateTask(task)

		phase := &types.TaskPhase{
			ID:          "test-phase-3",
			TaskID:      "test-phase-task-3",
			PhaseName:   "discovery",
			DisplayName: "Discovery",
			Sequence:    0,
			AgentName:   "investigator",
			ModelName:   "haiku",
			Status:      types.PhaseStatusRunning,
		}
		repo.CreatePhase(phase)

		// Update output
		output := "Phase completed successfully"
		err := repo.UpdatePhaseOutput("test-phase-3", output)
		if err != nil {
			t.Fatalf("Failed to update phase output: %v", err)
		}

		// Verify update
		phases, _ := repo.GetPhasesByTask("test-phase-task-3")
		if phases[0].OutputText != output {
			t.Errorf("Expected output '%s', got '%s'", output, phases[0].OutputText)
		}
	})
}

// TestEvents 测试事件操作
func TestEvents(t *testing.T) {
	database, cleanup := getTestDB(t)
	defer cleanup()

	repo := db.NewRepository(database)

	t.Run("CreateEvent", func(t *testing.T) {
		// Create a task first
		now := time.Now()
		task := &types.Task{
			ID:           "test-event-task-1",
			WorkflowID:   "code_simple",
			Prompt:       "Event test",
			WorkflowType: "CODE_SIMPLE",
			Status:       types.TaskStatusCreated,
			CurrentPhase: 0,
			TotalPhases:  1,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		repo.CreateTask(task)

		// Create an event
		event := &types.Event{
			ID:         "test-event-1",
			TaskID:     "test-event-task-1",
			EventType:  types.EventTypeTaskCreated,
			EventLevel: types.EventLevelInfo,
		}

		err := repo.CreateEvent(event)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}

		// Verify event was created
		events, err := repo.GetEventsByTask("test-event-task-1", 0)
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
}

// TestPersistence 测试文件数据库持久化
func TestPersistence(t *testing.T) {
	// Create temp file
	tmpFile := "/tmp/gclm-test-" + os.TempDir() + ".db"
	defer os.Remove(tmpFile)

	cfg := &db.Config{Path: tmpFile}
	database, err := db.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	repo := db.NewRepository(database)

	// Create a task
	task := &types.Task{
		ID:           "test-persist-1",
		WorkflowID:   "code_simple",
		Prompt:       "Persistence test",
		WorkflowType: "CODE_SIMPLE",
		Status:       types.TaskStatusCreated,
		CurrentPhase: 0,
		TotalPhases:  1,
	}
	repo.CreateTask(task)

	// Close database
	database.Close()

	// Reopen database
	database, err = db.New(cfg)
	if err != nil {
		t.Fatalf("Failed to reopen database: %v", err)
	}
	defer database.Close()

	repo = db.NewRepository(database)

	// Verify task still exists
	retrieved, err := repo.GetTask("test-persist-1")
	if err != nil {
		t.Fatalf("Failed to retrieve task after reopen: %v", err)
	}

	if retrieved.Prompt != "Persistence test" {
		t.Errorf("Expected prompt 'Persistence test', got '%s'", retrieved.Prompt)
	}
}
