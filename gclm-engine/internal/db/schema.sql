-- gclm-engine PoC Database Schema
-- SQLite version 3.x
-- WAL mode enabled for concurrency

-- Enable WAL mode for better concurrency
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;

-- ============================================================================
-- Tasks Table
-- ============================================================================
-- Represents a user task that needs to be processed through a pipeline
-- ============================================================================
CREATE TABLE IF NOT EXISTS tasks (
    -- Primary key: UUID v4
    id TEXT PRIMARY KEY,

    -- Pipeline reference
    pipeline_id TEXT NOT NULL,

    -- User input
    prompt TEXT NOT NULL,

    -- Workflow classification
    workflow_type TEXT NOT NULL,  -- 'DOCUMENT', 'CODE_SIMPLE', 'CODE_COMPLEX'

    -- Status tracking
    status TEXT NOT NULL DEFAULT 'created',  -- created, running, paused, completed, failed, cancelled

    -- Phase tracking
    current_phase INTEGER DEFAULT 0,
    total_phases INTEGER NOT NULL,

    -- Results
    result TEXT,                    -- Final output summary
    error_message TEXT,             -- Error details if failed

    -- Timestamps
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    started_at TEXT,
    completed_at TEXT,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Indexes for tasks
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_workflow_type ON tasks(workflow_type);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);


-- ============================================================================
-- Task Phases Table
-- ============================================================================
-- Tracks execution of individual phases within a task
-- ============================================================================
CREATE TABLE IF NOT EXISTS task_phases (
    -- Primary key
    id TEXT PRIMARY KEY,

    -- Foreign key to tasks
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,

    -- Phase information
    phase_name TEXT NOT NULL,       -- e.g., 'discovery', 'clarification'
    display_name TEXT NOT NULL,     -- e.g., 'Discovery / 需求发现'

    -- Execution order
    sequence INTEGER NOT NULL,

    -- Agent configuration
    agent_name TEXT,                -- e.g., 'investigator', 'tdd-guide'
    model_name TEXT,                -- e.g., 'haiku', 'sonnet', 'opus'

    -- Status
    status TEXT NOT NULL DEFAULT 'pending',  -- pending, running, completed, failed, skipped

    -- Execution details
    input_prompt TEXT,              -- Input given to the agent
    output_text TEXT,               -- Output from the agent
    error_message TEXT,             -- Error details if failed

    -- Timing
    started_at TEXT,
    completed_at TEXT,
    duration_ms INTEGER,            -- Execution duration in milliseconds

    -- Timestamps
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Indexes for task_phases
CREATE INDEX IF NOT EXISTS idx_task_phases_task_id ON task_phases(task_id);
CREATE INDEX IF NOT EXISTS idx_task_phases_sequence ON task_phases(task_id, sequence);
CREATE INDEX IF NOT EXISTS idx_task_phases_status ON task_phases(status);


-- ============================================================================
-- Events Table
-- ============================================================================
-- Audit log and event stream for all state changes
-- ============================================================================
CREATE TABLE IF NOT EXISTS events (
    -- Primary key
    id TEXT PRIMARY KEY,

    -- Event source
    task_id TEXT REFERENCES tasks(id) ON DELETE CASCADE,
    phase_id TEXT REFERENCES task_phases(id) ON DELETE SET NULL,

    -- Event details
    event_type TEXT NOT NULL,       -- task_created, phase_started, phase_completed, etc.
    event_level TEXT NOT NULL DEFAULT 'info',  -- debug, info, warn, error

    -- Event data (JSON)
    data TEXT,                      -- Additional event context as JSON

    -- Timestamp
    occurred_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Indexes for events
CREATE INDEX IF NOT EXISTS idx_events_task_id ON events(task_id);
CREATE INDEX IF NOT EXISTS idx_events_occurred_at ON events(occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);


-- ============================================================================
-- Workflows Table
-- ============================================================================
-- Stores workflow definitions (previously in YAML files)
-- ============================================================================
CREATE TABLE IF NOT EXISTS workflows (
    -- Primary key: workflow name (unique identifier)
    name TEXT PRIMARY KEY,

    -- Display information
    display_name TEXT NOT NULL,
    description TEXT,

    -- Workflow classification
    workflow_type TEXT NOT NULL,  -- 'document', 'code_simple', 'code_complex'

    -- Versioning
    version TEXT NOT NULL DEFAULT '1.0.0',

    -- Builtin flag (true for default workflows that cannot be deleted)
    is_builtin INTEGER NOT NULL DEFAULT 0,

    -- Workflow configuration (YAML as text)
    config_yaml TEXT NOT NULL,

    -- Timestamps
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Indexes for workflows
CREATE INDEX IF NOT EXISTS idx_workflows_type ON workflows(workflow_type);
CREATE INDEX IF NOT EXISTS idx_workflows_builtin ON workflows(is_builtin);

-- Trigger to update workflows.updated_at
CREATE TRIGGER IF NOT EXISTS update_workflows_timestamp
AFTER UPDATE ON workflows
BEGIN
    UPDATE workflows SET updated_at = datetime('now') WHERE name = NEW.name;
END;


-- ============================================================================
-- Triggers for automatic timestamp updates
-- ============================================================================

-- Update tasks.updated_at on any modification
CREATE TRIGGER IF NOT EXISTS update_tasks_timestamp
AFTER UPDATE ON tasks
BEGIN
    UPDATE tasks SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Update task_phases.updated_at on any modification
CREATE TRIGGER IF NOT EXISTS update_task_phases_timestamp
AFTER UPDATE ON task_phases
BEGIN
    UPDATE task_phases SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Calculate duration_ms when phase completes
CREATE TRIGGER IF NOT EXISTS calculate_phase_duration
AFTER UPDATE OF status ON task_phases
WHEN NEW.status IN ('completed', 'failed') AND OLD.status NOT IN ('completed', 'failed')
BEGIN
    UPDATE task_phases
    SET duration_ms = CAST((julianday(NEW.completed_at) - julianday(NEW.started_at)) * 86400000 AS INTEGER)
    WHERE id = NEW.id;
END;


-- ============================================================================
-- Views for common queries
-- ============================================================================

-- Active tasks view
CREATE VIEW IF NOT EXISTS active_tasks AS
SELECT id, pipeline_id, workflow_type, status, current_phase, total_phases,
       prompt, created_at, updated_at
FROM tasks
WHERE status IN ('created', 'running');


-- Task phases summary view
CREATE VIEW IF NOT EXISTS task_phases_summary AS
SELECT
    t.id AS task_id,
    t.status AS task_status,
    COUNT(tp.id) AS total_phases,
    SUM(CASE WHEN tp.status = 'completed' THEN 1 ELSE 0 END) AS completed_phases,
    SUM(CASE WHEN tp.status = 'pending' THEN 1 ELSE 0 END) AS pending_phases,
    SUM(CASE WHEN tp.status = 'running' THEN 1 ELSE 0 END) AS running_phases,
    SUM(CASE WHEN tp.status = 'failed' THEN 1 ELSE 0 END) AS failed_phases
FROM tasks t
LEFT JOIN task_phases tp ON t.id = tp.task_id
GROUP BY t.id;


-- ============================================================================
-- Sample data for PoC testing
-- ============================================================================

-- Insert a sample CODE_SIMPLE pipeline (will be replaced by YAML config)
-- This is just for initial testing
