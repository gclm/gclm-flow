-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
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

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_workflow_type ON tasks(workflow_type);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);

CREATE TABLE IF NOT EXISTS task_phases (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    phase_name TEXT NOT NULL,
    display_name TEXT NOT NULL,
    sequence INTEGER NOT NULL,
    agent_name TEXT,
    model_name TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    input_prompt TEXT,
    output_text TEXT,
    error_message TEXT,
    started_at TEXT,
    completed_at TEXT,
    duration_ms INTEGER,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_task_phases_task_id ON task_phases(task_id);
CREATE INDEX IF NOT EXISTS idx_task_phases_sequence ON task_phases(task_id, sequence);
CREATE INDEX IF NOT EXISTS idx_task_phases_status ON task_phases(status);

CREATE TABLE IF NOT EXISTS events (
    id TEXT PRIMARY KEY,
    task_id TEXT REFERENCES tasks(id) ON DELETE CASCADE,
    phase_id TEXT REFERENCES task_phases(id) ON DELETE SET NULL,
    event_type TEXT NOT NULL,
    event_level TEXT NOT NULL DEFAULT 'info',
    data TEXT,
    occurred_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_events_task_id ON events(task_id);
CREATE INDEX IF NOT EXISTS idx_events_occurred_at ON events(occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);

CREATE TABLE IF NOT EXISTS workflows (
    name TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    description TEXT,
    workflow_type TEXT NOT NULL,
    version TEXT NOT NULL DEFAULT '1.0.0',
    is_builtin INTEGER NOT NULL DEFAULT 0,
    config_yaml TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_workflows_type ON workflows(workflow_type);
CREATE INDEX IF NOT EXISTS idx_workflows_builtin ON workflows(is_builtin);

CREATE TRIGGER IF NOT EXISTS update_tasks_timestamp AFTER UPDATE ON tasks BEGIN UPDATE tasks SET updated_at = datetime('now') WHERE id = NEW.id; END;

CREATE TRIGGER IF NOT EXISTS update_task_phases_timestamp AFTER UPDATE ON task_phases BEGIN UPDATE task_phases SET updated_at = datetime('now') WHERE id = NEW.id; END;

CREATE TRIGGER IF NOT EXISTS calculate_phase_duration AFTER UPDATE OF status ON task_phases WHEN NEW.status IN ('completed', 'failed') AND OLD.status NOT IN ('completed', 'failed') BEGIN UPDATE task_phases SET duration_ms = CAST((julianday(NEW.completed_at) - julianday(NEW.started_at)) * 86400000 AS INTEGER) WHERE id = NEW.id; END;

CREATE VIEW IF NOT EXISTS active_tasks AS SELECT id, pipeline_id, workflow_type, status, current_phase, total_phases, prompt, created_at, updated_at FROM tasks WHERE status IN ('created', 'running');

CREATE VIEW IF NOT EXISTS task_phases_summary AS SELECT t.id AS task_id, t.status AS task_status, COUNT(tp.id) AS total_phases, SUM(CASE WHEN tp.status = 'completed' THEN 1 ELSE 0 END) AS completed_phases, SUM(CASE WHEN tp.status = 'pending' THEN 1 ELSE 0 END) AS pending_phases, SUM(CASE WHEN tp.status = 'running' THEN 1 ELSE 0 END) AS running_phases, SUM(CASE WHEN tp.status = 'failed' THEN 1 ELSE 0 END) AS failed_phases FROM tasks t LEFT JOIN task_phases tp ON t.id = tp.task_id GROUP BY t.id;

-- +goose Down
DROP VIEW IF EXISTS task_phases_summary;
DROP VIEW IF EXISTS active_tasks;
DROP TRIGGER IF EXISTS calculate_phase_duration;
DROP TRIGGER IF EXISTS update_task_phases_timestamp;
DROP TRIGGER IF EXISTS update_tasks_timestamp;
DROP INDEX IF EXISTS idx_workflows_builtin;
DROP INDEX IF EXISTS idx_workflows_type;
DROP TABLE IF EXISTS workflows;
DROP INDEX IF EXISTS idx_events_event_type;
DROP INDEX IF EXISTS idx_events_occurred_at;
DROP INDEX IF EXISTS idx_events_task_id;
DROP TABLE IF EXISTS events;
DROP INDEX IF EXISTS idx_task_phases_status;
DROP INDEX IF EXISTS idx_task_phases_sequence;
DROP INDEX IF EXISTS idx_task_phases_task_id;
DROP TABLE IF EXISTS task_phases;
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_workflow_type;
DROP INDEX IF EXISTS idx_tasks_status;
DROP TABLE IF EXISTS tasks;
