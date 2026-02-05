-- +goose Up
-- Optimize indexes for better query performance
-- Based on actual query patterns in the codebase

-- Step 1: Fix tasks table indexes (remove old _new suffix and add composite index)
DROP INDEX IF EXISTS idx_tasks_new_status;
DROP INDEX IF EXISTS idx_tasks_new_workflow_type;
DROP INDEX IF EXISTS idx_tasks_new_created_at;

-- Composite index for: WHERE status = ? ORDER BY created_at DESC
CREATE INDEX idx_tasks_status_created ON tasks(status, created_at DESC);

-- Keep workflow_type index for workflow filtering
CREATE INDEX IF NOT EXISTS idx_tasks_workflow_type ON tasks(workflow_type);

-- Step 2: Optimize task_phases indexes
-- Current indexes are already good, just verify they exist
-- The composite index (task_id, sequence) is already optimal

-- Step 3: Optimize events indexes (composite index for task query pattern)
DROP INDEX IF EXISTS idx_events_task_id;
DROP INDEX IF EXISTS idx_events_occurred_at;

-- Composite index for: WHERE task_id = ? ORDER BY occurred_at DESC
CREATE INDEX idx_events_task_id_occurred ON events(task_id, occurred_at DESC);

-- Keep event_type index for event type filtering
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);

-- Step 4: Optimize workflows indexes
-- Add composite index for name + is_builtin lookups
CREATE INDEX IF NOT EXISTS idx_workflows_name_builtin ON workflows(name, is_builtin);

-- Keep type and builtin indexes for filtering
CREATE INDEX IF NOT EXISTS idx_workflows_type ON workflows(workflow_type);
CREATE INDEX IF NOT EXISTS idx_workflows_builtin ON workflows(is_builtin);

-- +goose Down
-- Revert index optimizations

-- Step 1: Revert tasks indexes
DROP INDEX IF EXISTS idx_tasks_status_created;
DROP INDEX IF EXISTS idx_tasks_workflow_type;

CREATE INDEX idx_tasks_new_status ON tasks(status);
CREATE INDEX idx_tasks_new_workflow_type ON tasks(workflow_type);
CREATE INDEX idx_tasks_new_created_at ON tasks(created_at DESC);

-- Step 2: Revert events indexes
DROP INDEX IF EXISTS idx_events_task_id_occurred;
DROP INDEX IF EXISTS idx_events_event_type;

CREATE INDEX idx_events_task_id ON events(task_id);
CREATE INDEX idx_events_occurred_at ON events(occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);

-- Step 3: Revert workflows indexes
DROP INDEX IF EXISTS idx_workflows_name_builtin;

-- Keep type and builtin indexes
CREATE INDEX IF NOT EXISTS idx_workflows_type ON workflows(workflow_type);
CREATE INDEX IF NOT EXISTS idx_workflows_builtin ON workflows(is_builtin);
