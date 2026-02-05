-- +goose Up
-- Rename pipeline_id to workflow_id across the database
-- This aligns the schema with the codebase refactoring (pipeline -> workflow)

-- Note: SQLite doesn't support ALTER TABLE RENAME COLUMN directly
-- We need to recreate tables

-- Step 1: Drop views that depend on tasks table
DROP VIEW IF EXISTS active_tasks;
DROP VIEW IF EXISTS task_phases_summary;

-- Step 2: Drop triggers that depend on tasks table
DROP TRIGGER IF EXISTS update_tasks_timestamp;
DROP TRIGGER IF EXISTS calculate_phase_duration;

-- Step 3: Create new versions of tables with workflow_id
CREATE TABLE IF NOT EXISTS tasks_new (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL,
    prompt TEXT NOT NULL,
    workflow_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'created',
    current_phase INTEGER DEFAULT 0,
    total_phases INTEGER NOT NULL,
    result TEXT,
    error_message TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    started_at TEXT,
    completed_at TEXT,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_tasks_new_status ON tasks_new(status);
CREATE INDEX IF NOT EXISTS idx_tasks_new_workflow_type ON tasks_new(workflow_type);
CREATE INDEX IF NOT EXISTS idx_tasks_new_created_at ON tasks_new(created_at DESC);

-- Step 4: Migrate data from old table to new table
INSERT INTO tasks_new (
    id, workflow_id, prompt, workflow_type, status,
    current_phase, total_phases, result, error_message,
    created_at, started_at, completed_at, updated_at
)
SELECT
    id, pipeline_id, prompt, workflow_type, status,
    current_phase, total_phases, result, error_message,
    created_at, started_at, completed_at, updated_at
FROM tasks;

-- Step 5: Drop old table and rename new table
DROP TABLE tasks;
ALTER TABLE tasks_new RENAME TO tasks;

-- Step 6: Recreate triggers
CREATE TRIGGER IF NOT EXISTS update_tasks_timestamp AFTER UPDATE ON tasks BEGIN UPDATE tasks SET updated_at = datetime('now') WHERE id = NEW.id; END;

-- Step 7: Recreate views with workflow_id
CREATE VIEW IF NOT EXISTS active_tasks AS SELECT id, workflow_id, workflow_type, status, current_phase, total_phases, prompt, created_at, updated_at FROM tasks WHERE status IN ('created', 'running');

-- +goose Down
-- Revert: Rename workflow_id back to pipeline_id

-- Step 1: Drop views that depend on tasks table
DROP VIEW IF EXISTS active_tasks;
DROP VIEW IF EXISTS task_phases_summary;

-- Step 2: Drop triggers
DROP TRIGGER IF EXISTS update_tasks_timestamp;
DROP TRIGGER IF EXISTS calculate_phase_duration;

-- Step 3: Create old table with pipeline_id
CREATE TABLE IF NOT EXISTS tasks_old (
    id TEXT PRIMARY KEY,
    pipeline_id TEXT NOT NULL,
    prompt TEXT NOT NULL,
    workflow_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'created',
    current_phase INTEGER DEFAULT 0,
    total_phases INTEGER NOT NULL,
    result TEXT,
    error_message TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    started_at TEXT,
    completed_at TEXT,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_tasks_old_status ON tasks_old(status);
CREATE INDEX IF NOT EXISTS idx_tasks_old_workflow_type ON tasks_old(workflow_type);
CREATE INDEX IF NOT EXISTS idx_tasks_old_created_at ON tasks_old(created_at DESC);

-- Step 4: Migrate data back
INSERT INTO tasks_old (
    id, pipeline_id, prompt, workflow_type, status,
    current_phase, total_phases, result, error_message,
    created_at, started_at, completed_at, updated_at
)
SELECT
    id, workflow_id, prompt, workflow_type, status,
    current_phase, total_phases, result, error_message,
    created_at, started_at, completed_at, updated_at
FROM tasks;

-- Step 5: Replace table
DROP TABLE tasks;
ALTER TABLE tasks_old RENAME TO tasks;

-- Step 6: Recreate triggers
CREATE TRIGGER IF NOT EXISTS update_tasks_timestamp AFTER UPDATE ON tasks BEGIN UPDATE tasks SET updated_at = datetime('now') WHERE id = NEW.id; END;

-- Step 7: Recreate views with pipeline_id
CREATE VIEW IF NOT EXISTS active_tasks AS SELECT id, pipeline_id, workflow_type, status, current_phase, total_phases, prompt, created_at, updated_at FROM tasks WHERE status IN ('created', 'running');
