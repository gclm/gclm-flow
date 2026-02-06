package types

import "time"

// TaskStatus represents the current state of a task
type TaskStatus string

const (
	TaskStatusCreated   TaskStatus = "created"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusPaused    TaskStatus = "paused"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// PhaseStatus represents the current state of a phase
type PhaseStatus string

const (
	PhaseStatusPending   PhaseStatus = "pending"
	PhaseStatusRunning   PhaseStatus = "running"
	PhaseStatusCompleted PhaseStatus = "completed"
	PhaseStatusFailed    PhaseStatus = "failed"
	PhaseStatusSkipped   PhaseStatus = "skipped"
)

// EventType represents the type of event
type EventType string

const (
	EventTypeTaskCreated       EventType = "task_created"
	EventTypeTaskStarted       EventType = "task_started"
	EventTypeTaskCompleted     EventType = "task_completed"
	EventTypeTaskFailed        EventType = "task_failed"
	EventTypeTaskCancelled     EventType = "task_cancelled"
	EventTypeTaskPaused        EventType = "task_paused"
	EventTypeTaskResumed       EventType = "task_resumed"
	EventTypePhaseStarted      EventType = "phase_started"
	EventTypePhaseCompleted    EventType = "phase_completed"
	EventTypePhaseFailed       EventType = "phase_failed"
	EventTypePhaseSkipped      EventType = "phase_skipped"
	EventTypeAgentInvoked      EventType = "agent_invoked"
	EventTypeAgentCompleted    EventType = "agent_completed"
	EventTypeAgentFailed       EventType = "agent_failed"
)

// EventLevel represents the severity level of an event
type EventLevel string

const (
	EventLevelDebug EventLevel = "debug"
	EventLevelInfo  EventLevel = "info"
	EventLevelWarn  EventLevel = "warn"
	EventLevelError EventLevel = "error"
)

// Task represents a user task
type Task struct {
	ID           string        `json:"id" db:"id"`
	WorkflowID   string        `json:"workflowId" db:"workflow_id"`
	Prompt       string        `json:"prompt" db:"prompt"`
	WorkflowType string        `json:"workflowType" db:"workflow_type"`
	Status       TaskStatus    `json:"status" db:"status"`
	CurrentPhase int           `json:"currentPhase" db:"current_phase"`
	TotalPhases  int           `json:"totalPhases" db:"total_phases"`
	Result       string        `json:"result,omitempty" db:"result"`
	Error        string        `json:"error,omitempty" db:"error_message"`
	CreatedAt    time.Time     `json:"createdAt" db:"created_at"`
	StartedAt    *time.Time    `json:"startedAt,omitempty" db:"started_at"`
	CompletedAt  *time.Time    `json:"completedAt,omitempty" db:"completed_at"`
	UpdatedAt    time.Time     `json:"updatedAt" db:"updated_at"`
}

// TaskPhase represents a single phase in a task
type TaskPhase struct {
	ID          string      `json:"id" db:"id"`
	TaskID      string      `json:"taskId" db:"task_id"`
	PhaseName   string      `json:"phaseName" db:"phase_name"`
	DisplayName string      `json:"displayName" db:"display_name"`
	Sequence    int         `json:"sequence" db:"sequence"`
	AgentName   string      `json:"agentName,omitempty" db:"agent_name"`
	ModelName   string      `json:"modelName,omitempty" db:"model_name"`
	Status      PhaseStatus `json:"status" db:"status"`
	InputPrompt string      `json:"inputPrompt,omitempty" db:"input_prompt"`
	OutputText  string      `json:"outputText,omitempty" db:"output_text"`
	Error       string      `json:"error,omitempty" db:"error_message"`
	StartedAt   *time.Time  `json:"startedAt,omitempty" db:"started_at"`
	CompletedAt *time.Time  `json:"completedAt,omitempty" db:"completed_at"`
	DurationMs  int         `json:"durationMs,omitempty" db:"duration_ms"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at"`
}

// Event represents an audit log event
type Event struct {
	ID         string     `json:"id" db:"id"`
	TaskID     string     `json:"taskId,omitempty" db:"task_id"`
	PhaseID    string     `json:"phaseId,omitempty" db:"phase_id"`
	EventType  EventType  `json:"eventType" db:"event_type"`
	EventLevel EventLevel `json:"eventLevel" db:"event_level"`
	Data       string     `json:"data,omitempty" db:"data"` // JSON string
	OccurredAt time.Time  `json:"occurredAt" db:"occurred_at"`
}

// CreateTaskRequest represents a request to create a new task
type CreateTaskRequest struct {
	Prompt       string `json:"prompt" validate:"required"`
	WorkflowType string `json:"workflowType,omitempty"`
	WorkflowID   string `json:"workflowId,omitempty"`
}

// TaskResponse represents a task response with additional metadata
type TaskResponse struct {
	Task   *Task       `json:"task"`
	Phases []*TaskPhase `json:"phases,omitempty"`
	Events []*Event    `json:"events,omitempty"`
}
