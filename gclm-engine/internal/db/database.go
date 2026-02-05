package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gclm/gclm-flow/gclm-engine/pkg/types"
	"github.com/pressly/goose/v3"
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the SQLite database connection
type Database struct {
	conn *sql.DB
	dsn  string
}

// Config holds database configuration
type Config struct {
	Path string // Path to SQLite database file
}

// DefaultConfig returns default database configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Path: filepath.Join(homeDir, ".gclm-flow", "gclm-engine.db"),
	}
}

// New creates a new database instance
func New(cfg *Config) (*Database, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Ensure directory exists
	dbDir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", cfg.Path)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite works best with single writer
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	database := &Database{
		conn: db,
		dsn:  dsn,
	}

	// Initialize schema
	if err := database.init(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return database, nil
}

// init initializes the database schema using goose migrations
func (d *Database) init() error {
	// Configure goose
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Determine migrations directory
	// For development: use ./migrations relative to working directory
	// For production: use ./migrations relative to executable
	migrationsDir := "migrations"

	// Check if migrations directory exists in working directory
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try relative to executable
		execDir, err := os.Executable()
		if err == nil {
			migrationsDir = filepath.Join(filepath.Dir(execDir), "migrations")
		}
	}

	// Run migrations
	if err := goose.Up(d.conn, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.conn.Close()
}

// GetDB returns the underlying sql.DB instance
func (d *Database) GetDB() *sql.DB {
	return d.conn
}

// InitWorkflows initializes builtin workflows from the workflows directory
func (d *Database) InitWorkflows(workflowsDir string) error {
	wfRepo := NewWorkflowRepository(d)
	return wfRepo.InitializeBuiltinWorkflows(workflowsDir)
}

// BeginTx starts a new transaction
func (d *Database) BeginTx() (*sql.Tx, error) {
	return d.conn.Begin()
}

// Repository provides database operations
type Repository struct {
	db *Database
}

// NewRepository creates a new repository
func NewRepository(db *Database) *Repository {
	return &Repository{db: db}
}

// Task operations

// CreateTask creates a new task
func (r *Repository) CreateTask(task *types.Task) error {
	query := `
		INSERT INTO tasks (
			id, workflow_id, prompt, workflow_type, status,
			current_phase, total_phases, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.conn.Exec(query,
		task.ID,
		task.WorkflowID,
		task.Prompt,
		task.WorkflowType,
		task.Status,
		task.CurrentPhase,
		task.TotalPhases,
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

// GetTask retrieves a task by ID
func (r *Repository) GetTask(id string) (*types.Task, error) {
	query := `
		SELECT id, workflow_id, prompt, workflow_type, status,
		       current_phase, total_phases, result, error_message,
		       created_at, started_at, completed_at, updated_at
		FROM tasks WHERE id = ?
	`

	task := &types.Task{}
	var startedAt, completedAt, result, errMsg, createdAt, updatedAt sql.NullString

	err := r.db.conn.QueryRow(query, id).Scan(
		&task.ID,
		&task.WorkflowID,
		&task.Prompt,
		&task.WorkflowType,
		&task.Status,
		&task.CurrentPhase,
		&task.TotalPhases,
		&result,
		&errMsg,
		&createdAt,
		&startedAt,
		&completedAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	// Parse time fields
	if createdAt.Valid {
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
	}
	if updatedAt.Valid {
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt.String)
	}
	if result.Valid {
		task.Result = result.String
	}
	if errMsg.Valid {
		task.Error = errMsg.String
	}
	if startedAt.Valid {
		t, _ := time.Parse(time.RFC3339, startedAt.String)
		task.StartedAt = &t
	}
	if completedAt.Valid {
		t, _ := time.Parse(time.RFC3339, completedAt.String)
		task.CompletedAt = &t
	}

	return task, nil
}

// ListTasks retrieves all tasks with optional filtering
func (r *Repository) ListTasks(status *types.TaskStatus, limit int) ([]*types.Task, error) {
	query := `
		SELECT id, workflow_id, prompt, workflow_type, status,
		       current_phase, total_phases, result, error_message,
		       created_at, started_at, completed_at, updated_at
		FROM tasks
	`
	args := []any{}

	if status != nil {
		query += " WHERE status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := r.db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*types.Task
	for rows.Next() {
		task := &types.Task{}
		var startedAt, completedAt, result, errMsg, createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.WorkflowID,
			&task.Prompt,
			&task.WorkflowType,
			&task.Status,
			&task.CurrentPhase,
			&task.TotalPhases,
			&result,
			&errMsg,
			&createdAt,
			&startedAt,
			&completedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse time fields
		if createdAt.Valid {
			task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
		}
		if updatedAt.Valid {
			task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt.String)
		}
		if result.Valid {
			task.Result = result.String
		}
		if errMsg.Valid {
			task.Error = errMsg.String
		}
		if startedAt.Valid {
			t, _ := time.Parse(time.RFC3339, startedAt.String)
			task.StartedAt = &t
		}
		if completedAt.Valid {
			t, _ := time.Parse(time.RFC3339, completedAt.String)
			task.CompletedAt = &t
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

// UpdateTaskStatus updates the status of a task
func (r *Repository) UpdateTaskStatus(id string, status types.TaskStatus) error {
	query := `UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.conn.Exec(query, status, time.Now().Format(time.RFC3339), id)
	return err
}

// UpdateTaskProgress updates the current phase of a task
func (r *Repository) UpdateTaskProgress(id string, currentPhase int) error {
	query := `UPDATE tasks SET current_phase = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.conn.Exec(query, currentPhase, time.Now().Format(time.RFC3339), id)
	return err
}

// CompleteTask marks a task as completed
func (r *Repository) CompleteTask(id string, result string) error {
	query := `
		UPDATE tasks
		SET status = ?, result = ?, completed_at = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now().Format(time.RFC3339)
	_, err := r.db.conn.Exec(query, types.TaskStatusCompleted, result, now, now, id)
	return err
}

// FailTask marks a task as failed
func (r *Repository) FailTask(id string, errMsg string) error {
	query := `
		UPDATE tasks
		SET status = ?, error_message = ?, completed_at = ?, updated_at = ?
		WHERE id = ?
	`
	now := time.Now().Format(time.RFC3339)
	_, err := r.db.conn.Exec(query, types.TaskStatusFailed, errMsg, now, now, id)
	return err
}

// TaskPhase operations

// CreatePhase creates a new task phase
func (r *Repository) CreatePhase(phase *types.TaskPhase) error {
	query := `
		INSERT INTO task_phases (
			id, task_id, phase_name, display_name, sequence,
			agent_name, model_name, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.conn.Exec(query,
		phase.ID,
		phase.TaskID,
		phase.PhaseName,
		phase.DisplayName,
		phase.Sequence,
		phase.AgentName,
		phase.ModelName,
		phase.Status,
		phase.CreatedAt.Format(time.RFC3339),
		phase.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

// GetPhasesByTask retrieves all phases for a task
func (r *Repository) GetPhasesByTask(taskID string) ([]*types.TaskPhase, error) {
	query := `
		SELECT id, task_id, phase_name, display_name, sequence,
		       agent_name, model_name, status, input_prompt, output_text,
		       error_message, started_at, completed_at, duration_ms,
		       created_at, updated_at
		FROM task_phases
		WHERE task_id = ?
		ORDER BY sequence ASC
	`

	rows, err := r.db.conn.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phases []*types.TaskPhase
	for rows.Next() {
		phase := &types.TaskPhase{}
		var startedAt, completedAt, createdAt, updatedAt sql.NullString
		var inputPrompt, outputText, errMsg sql.NullString
		var durationMs sql.NullInt64

		err := rows.Scan(
			&phase.ID,
			&phase.TaskID,
			&phase.PhaseName,
			&phase.DisplayName,
			&phase.Sequence,
			&phase.AgentName,
			&phase.ModelName,
			&phase.Status,
			&inputPrompt,
			&outputText,
			&errMsg,
			&startedAt,
			&completedAt,
			&durationMs,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse time fields
		if createdAt.Valid {
			phase.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
		}
		if updatedAt.Valid {
			phase.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt.String)
		}
		if durationMs.Valid {
			phase.DurationMs = int(durationMs.Int64)
		}
		if inputPrompt.Valid {
			phase.InputPrompt = inputPrompt.String
		}
		if outputText.Valid {
			phase.OutputText = outputText.String
		}
		if errMsg.Valid {
			phase.Error = errMsg.String
		}
		if startedAt.Valid {
			t, _ := time.Parse(time.RFC3339, startedAt.String)
			phase.StartedAt = &t
		}
		if completedAt.Valid {
			t, _ := time.Parse(time.RFC3339, completedAt.String)
			phase.CompletedAt = &t
		}

		phases = append(phases, phase)
	}

	return phases, rows.Err()
}

// GetPhase retrieves a single phase by ID
func (r *Repository) GetPhase(id string) (*types.TaskPhase, error) {
	query := `
		SELECT id, task_id, phase_name, display_name, sequence,
		       agent_name, model_name, status, input_prompt, output_text,
		       error_message, started_at, completed_at, duration_ms,
		       created_at, updated_at
		FROM task_phases
		WHERE id = ?
	`

	phase := &types.TaskPhase{}
	var startedAt, completedAt, createdAt, updatedAt sql.NullString
	var inputPrompt, outputText, errMsg sql.NullString
	var durationMs sql.NullInt64

	err := r.db.conn.QueryRow(query, id).Scan(
		&phase.ID,
		&phase.TaskID,
		&phase.PhaseName,
		&phase.DisplayName,
		&phase.Sequence,
		&phase.AgentName,
		&phase.ModelName,
		&phase.Status,
		&inputPrompt,
		&outputText,
		&errMsg,
		&startedAt,
		&completedAt,
		&durationMs,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse time fields
	if createdAt.Valid {
		phase.CreatedAt, _ = time.Parse(time.RFC3339, createdAt.String)
	}
	if updatedAt.Valid {
		phase.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt.String)
	}
	if durationMs.Valid {
		phase.DurationMs = int(durationMs.Int64)
	}
	if inputPrompt.Valid {
		phase.InputPrompt = inputPrompt.String
	}
	if outputText.Valid {
		phase.OutputText = outputText.String
	}
	if errMsg.Valid {
		phase.Error = errMsg.String
	}
	if startedAt.Valid {
		t, _ := time.Parse(time.RFC3339, startedAt.String)
		phase.StartedAt = &t
	}
	if completedAt.Valid {
		t, _ := time.Parse(time.RFC3339, completedAt.String)
		phase.CompletedAt = &t
	}

	return phase, nil
}

// UpdatePhaseStatus updates the status of a phase
func (r *Repository) UpdatePhaseStatus(id string, status types.PhaseStatus) error {
	query := `UPDATE task_phases SET status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.conn.Exec(query, status, time.Now().Format(time.RFC3339), id)
	return err
}

// UpdatePhaseOutput updates the output of a completed phase
func (r *Repository) UpdatePhaseOutput(id string, output string) error {
	query := `UPDATE task_phases SET output_text = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.conn.Exec(query, output, time.Now().Format(time.RFC3339), id)
	return err
}

// Event operations

// CreateEvent creates a new event
func (r *Repository) CreateEvent(event *types.Event) error {
	query := `
		INSERT INTO events (id, task_id, phase_id, event_type, event_level, data, occurred_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.conn.Exec(query,
		event.ID,
		event.TaskID,
		event.PhaseID,
		event.EventType,
		event.EventLevel,
		event.Data,
		event.OccurredAt.Format(time.RFC3339),
	)

	return err
}

// GetEventsByTask retrieves all events for a task
func (r *Repository) GetEventsByTask(taskID string, limit int) ([]*types.Event, error) {
	query := `
		SELECT id, task_id, phase_id, event_type, event_level, data, occurred_at
		FROM events
		WHERE task_id = ?
		ORDER BY occurred_at DESC
	`

	if limit > 0 {
		query += " LIMIT ?"
	}

	rows, err := r.db.conn.Query(query, taskID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*types.Event
	for rows.Next() {
		event := &types.Event{}
		var phaseID, data, occurredAt sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.TaskID,
			&phaseID,
			&event.EventType,
			&event.EventLevel,
			&data,
			&occurredAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse time field
		if occurredAt.Valid {
			event.OccurredAt, _ = time.Parse(time.RFC3339, occurredAt.String)
		}
		if phaseID.Valid {
			event.PhaseID = phaseID.String
		}
		if data.Valid {
			event.Data = data.String
		}

		events = append(events, event)
	}

	return events, rows.Err()
}
