# 数据库设计

## 概述

gclm-engine 使用 SQLite (WAL 模式) 存储工作流状态，数据库位置：`~/.gclm-flow/gclm-engine.db`

**迁移系统**: 使用 Goose 进行数据库版本管理，迁移文件嵌入到二进制中。

---

## 数据库配置

```go
// db/database.go
dsn := "file:~/.gclm-flow/gclm-engine.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"
```

**配置说明**：
- `foreign_keys(1)`: 启用外键约束
- `journal_mode(WAL)`: Write-Ahead Logging，提升并发性能
- `SetMaxOpenConns(1)`: SQLite 单写入者限制

---

## 数据表结构

### 1. workflows - 工作流定义表

```sql
CREATE TABLE workflows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,              -- 工作流唯一标识
    display_name TEXT NOT NULL,             -- 显示名称
    description TEXT,                       -- 描述
    workflow_type TEXT NOT NULL,           -- DOCUMENT/CODE_SIMPLE/CODE_COMPLEX
    version TEXT NOT NULL,                 -- 版本号
    is_builtin BOOLEAN NOT NULL DEFAULT 0, -- 是否内置
    config_yaml TEXT NOT NULL,             -- YAML 配置内容
    created_at TEXT NOT NULL,              -- 创建时间
    updated_at TEXT NOT NULL               -- 更新时间
);
```

**索引**：
- `name` - 唯一索引
- `workflow_type` - 类型查询索引
- `is_builtin` - 内置/自定义区分

---

### 2. tasks - 任务表

```sql
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,                   -- 任务 UUID
    workflow_id TEXT NOT NULL,            -- 工作流名称 (原 pipeline_id)
    prompt TEXT NOT NULL,                  -- 用户输入
    workflow_type TEXT NOT NULL,          -- DOCUMENT/CODE_SIMPLE/CODE_COMPLEX
    status TEXT NOT NULL DEFAULT 'created', -- pending/running/completed/failed
    current_phase INTEGER NOT NULL,       -- 当前阶段序号
    total_phases INTEGER NOT NULL,        -- 总阶段数
    result TEXT,                          -- 任务结果
    error_message TEXT,                    -- 错误信息
    created_at TEXT NOT NULL,             -- 创建时间
    started_at TEXT,                      -- 开始时间
    completed_at TEXT,                    -- 完成时间
    updated_at TEXT NOT NULL              -- 更新时间
);
```

**索引**：
- `id` - 主键
- `status` - 状态查询
- `created_at` - 时间排序

---

### 3. task_phases - 任务阶段表

```sql
CREATE TABLE task_phases (
    id TEXT PRIMARY KEY,                   -- 阶段 UUID
    task_id TEXT NOT NULL,                -- 关联任务 ID
    phase_name TEXT NOT NULL,             -- 阶段名称 (如 discovery)
    display_name TEXT NOT NULL,           -- 显示名称
    sequence INTEGER NOT NULL,            -- 执行顺序
    agent_name TEXT NOT NULL,             -- Agent 名称
    model_name TEXT NOT NULL,             -- 模型名称 (haiku/sonnet/opus)
    status TEXT NOT NULL,                 -- pending/running/completed/failed/skipped
    input_prompt TEXT,                    -- 输入提示
    output_text TEXT,                     -- 输出结果
    error_message TEXT,                   -- 错误信息
    started_at TEXT,                      -- 开始时间
    completed_at TEXT,                    -- 完成时间
    duration_ms INTEGER,                  -- 执行时长（毫秒）
    created_at TEXT NOT NULL,             -- 创建时间
    updated_at TEXT NOT NULL,             -- 更新时间
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
```

**索引**：
- `id` - 主键
- `task_id` - 关联查询
- `status` - 状态查询
- `sequence` - 顺序查询

---

### 4. events - 事件日志表

```sql
CREATE TABLE events (
    id TEXT PRIMARY KEY,                   -- 事件 UUID
    task_id TEXT NOT NULL,                -- 关联任务 ID
    phase_id TEXT,                        -- 关联阶段 ID
    event_type TEXT NOT NULL,             -- 事件类型
    event_level TEXT NOT NULL,            -- 事件级别 (info/warning/error)
    data TEXT,                            -- 事件数据 (JSON)
    occurred_at TEXT NOT NULL,            -- 发生时间
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
```

**索引**：
- `id` - 主键
- `task_id` - 关联查询
- `occurred_at` - 时间排序

---

## 关系图

```
┌─────────────┐
│  workflows  │
│             │
│ - name      │
│ - config    │
│ - is_builtin│
└─────────────┘
       │
       │ 1:N
       ↓
┌─────────────┐     ┌─────────────┐
│   tasks     │────▶│ task_phases │
│             │     │             │
│ - id        │     │ - phase     │
│ - workflow  │     │ - agent     │
│ - status    │     │ - output    │
└─────────────┘     └─────────────┘
       │                     │
       │ N:1                 │ N:1
       ↓                     ↓
┌─────────────┐     ┌─────────────┐
│   events    │     │   events    │
│             │     │             │
│ - task_id   │     │ - phase_id  │
│ - type      │     │ - data      │
└─────────────┘     └─────────────┘
```

---

## 数据访问层

### Database (db/database.go)

```go
type Database struct {
    conn *sql.DB
    dsn  string
}

// 主要方法
func New(cfg *Config) (*Database, error)
func (d *Database) Close() error
func (d *Database) BeginTx() (*sql.Tx, error)
func (d *Database) InitWorkflows(workflowsDir string) error
func (d *Database) GetDB() *sql.DB
```

**迁移系统**:
- 使用 Goose (`github.com/pressly/goose/v3`)
- 迁移文件嵌入到二进制 (`embed.FS`)
- 自动检测并应用未执行的迁移

### WorkflowRepository (db/workflow.go)

```go
type WorkflowRepository struct {
    db *Database
}

// 主要方法
func (r *WorkflowRepository) InitializeBuiltinWorkflows(workflowsDir string) error
func (r *WorkflowRepository) GetWorkflow(name string) (*WorkflowRecord, error)
func (r *WorkflowRepository) GetWorkflowByType(workflowType string) (*WorkflowRecord, error)
func (r *WorkflowRepository) ListWorkflows() ([]WorkflowRecord, error)
func (r *WorkflowRepository) InstallWorkflow(name string, yamlData []byte) error
func (r *WorkflowRepository) UninstallWorkflow(name string) error
```

**WorkflowRecord**:
```go
type WorkflowRecord struct {
    Name        string
    DisplayName string
    Description string
    WorkflowType string
    Version     string
    IsBuiltin   bool
    ConfigYAML  string
}
```

### Repository (db/database.go)

```go
type Repository struct {
    db *Database
}

// Task 操作
func (r *Repository) CreateTask(task *types.Task) error
func (r *Repository) GetTask(id string) (*types.Task, error)
func (r *Repository) UpdateTaskStatus(id string, status types.TaskStatus) error
func (r *Repository) CompleteTask(id string, result string) error

// Phase 操作
func (r *Repository) CreatePhase(phase *types.TaskPhase) error
func (r *Repository) GetPhasesByTask(taskID string) ([]*types.TaskPhase, error)
func (r *Repository) UpdatePhaseOutput(id string, output string) error

// Event 操作
func (r *Repository) CreateEvent(event *types.Event) error
func (r *Repository) GetEventsByTask(taskID string, limit int) ([]*types.Event, error)
```

---

## 状态转换

### Task 状态

```
pending → running → completed
                  ↘ failed
```

### Phase 状态

```
pending → running → completed
                  ↘ failed
                  ↘ skipped (非必需阶段)
```

---

## 并发控制

由于 SQLite 的单写入者限制：

```go
db.SetMaxOpenConns(1)  // 单连接
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(time.Hour)
```

**影响**：
- 同时只能有一个写操作
- WAL 模式允许读写并发
- 读操作可以并发进行

---

## 性能优化

1. **WAL 模式**: 写操作不阻塞读操作
2. **索引**: 常用查询字段建立索引
3. **连接池**: 限制连接数，避免锁竞争
4. **批量操作**: 阶段完成时批量更新
5. **自动时间戳**: 使用触发器自动更新 `updated_at`
6. **视图**: 提供常用查询的预定义视图
   - `active_tasks`: 活跃任务视图
   - `task_phases_summary`: 任务阶段汇总视图

---

## 迁移历史

| 版本 | 文件 | 说明 |
|:---|:---|:---|
| 00001 | `baseline.sql` | 初始数据库结构 |
| 00002 | `rename_pipeline_to_workflow.sql` | 重命名 pipeline → workflow |
| 00003 | `optimize_indexes.sql` | 优化索引和添加视图 |
