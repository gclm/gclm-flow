# gclm-engine

> **gclm-flow 工作流引擎** - Go 语言实现的智能工作流编排系统

## 项目概述

gclm-engine 是 gclm-flow 的核心引擎，负责：
- **工作流配置管理** - 基于 YAML 的可扩展工作流定义
- **任务状态管理** - SQLite 持久化，支持暂停/恢复
- **阶段调度** - 依赖解析、拓扑排序、并行执行
- **CLI 接口** - JSON 输出，与 Claude Code skills 集成
- **Web API** - RESTful API + WebSocket 实时推送

## 项目结构

```
gclm-engine/
├── main.go                 # 入口文件
├── go.mod                  # Go 模块定义
├── Makefile                # 构建脚本
├── internal/
│   ├── api/                # Web API 层
│   │   ├── server.go       # HTTP 服务器 (Gin)
│   │   ├── handlers.go     # REST handlers
│   │   ├── utils.go        # 工具函数
│   │   └── websocket/
│   │       └── hub.go       # WebSocket Hub
│   ├── cli/                # CLI 命令 (cobra)
│   │   ├── root.go         # 根命令和初始化
│   │   ├── task_commands.go # 任务管理命令
│   │   ├── workflow_commands.go # 工作流命令
│   │   ├── init_commands.go # 初始化命令
│   │   ├── output.go       # 输出格式化
│   │   └── helpers.go      # 辅助函数
│   ├── db/                 # 数据库层
│   │   ├── database.go     # 数据库初始化
│   │   ├── repository.go   # 数据访问实现
│   │   └── workflow.go     # 工作流存储
│   ├── domain/             # 领域接口层
│   │   ├── errors.go       # 错误类型定义
│   │   ├── loader.go       # WorkflowLoader 接口
│   │   ├── repository.go   # 仓库接口
│   │   └── service.go      # 服务接口
│   ├── repository/         # 适配器层
│   │   ├── task_repository.go    # TaskRepository 适配器
│   │   ├── workflow_repository.go # WorkflowRepository 适配器
│   │   └── workflow_loader.go     # WorkflowLoader 适配器
│   ├── service/            # 服务层
│   │   ├── task.go         # 任务服务实现
│   │   └── workflow.go     # 工作流服务实现
│   ├── workflow/           # 工作流解析
│   │   └── parser.go       # YAML 解析、依赖图
│   ├── logger/             # 统一日志 (zerolog)
│   ├── assets/             # 嵌入资源
│   ├── config/             # 配置管理
│   └── errors/             # 错误处理
├── pkg/types/              # 共享类型定义
│   ├── types.go           # Task, Phase, Event
│   └── workflow.go        # Workflow, WorkflowNode
├── web/                    # Web 前端
│   ├── index.html         # 主页面
│   └── static/
│       ├── css/style.css  # 样式
│       └── js/app.js      # 前端逻辑
├── workflows/             # 工作流 YAML (运行时)
└── test/                  # 测试文件
```

## 快速开始

### 构建

```bash
cd gclm-engine
go build -o gclm-engine .
```

### 初始化

首次运行会自动初始化配置：

```bash
./gclm-engine version
# 自动创建 ~/.gclm-flow/ 目录并导出内置工作流
```

或手动初始化：

```bash
./gclm-engine init
```

### CLI 使用

```bash
# 创建任务
./gclm-engine task create "实现用户登录功能"

# 查看任务详情
./gclm-engine task get <task-id>

# 列出所有任务
./gclm-engine task list
```

### Web API 使用

```bash
# 启动 HTTP 服务器
./gclm-engine serve

# API 端点:
# - GET  /api/tasks - 列出任务
# - POST /api/tasks - 创建任务
# - GET  /api/tasks/:id - 获取任务详情
# - WS   /ws/tasks/:id - WebSocket 实时更新
```

## 命令参考

### 工作流管理

```bash
# 列出所有工作流
gclm-engine workflow list

# 查看工作流详情
gclm-engine workflow info <workflow-name>

# 验证工作流配置
gclm-engine workflow validate <yaml-file>

# 安装自定义工作流
gclm-engine workflow install <yaml-file>

# 卸载工作流
gclm-engine workflow uninstall <workflow-name>

# 导出工作流
gclm-engine workflow export <workflow-name>

# 同步工作流
gclm-engine workflow sync <workflows-dir>
```

### 任务管理

```bash
# 创建任务（自动检测工作流类型）
gclm-engine task create "实现用户登录功能"

# 指定工作流类型
gclm-engine task create "修复 bug" --workflow fix

# 查看任务详情
gclm-engine task get <task-id>

# 列出所有任务
gclm-engine task list
gclm-engine task list --status running

# 获取当前待执行阶段
gclm-engine task current <task-id>

# 获取完整执行计划
gclm-engine task plan <task-id>

# 完成阶段
gclm-engine task complete <task-id> <phase-id> --output "阶段输出"

# 标记阶段失败
gclm-engine task fail <task-id> <phase-id> --error "错误描述"

# 查看任务阶段
gclm-engine task phases <task-id>

# 查看任务事件
gclm-engine task events <task-id>

# 导出任务状态
gclm-engine task export <task-id> <output-file>

# 任务控制
gclm-engine task pause <task-id>
gclm-engine task resume <task-id>
gclm-engine task cancel <task-id>
```

### 工作流管理

```bash
# 列出所有工作流
gclm-engine workflow list

# 查看工作流详情
gclm-engine workflow info <workflow-name>

# 同步工作流配置
gclm-engine workflow sync
```

## 架构设计

### 分层架构

```
┌─────────────────────────────────────────────┐
│              CLI / Web API                   │
└─────────────────────────────────────────────┘
                      │
┌─────────────────────────────────────────────┐
│             Service Layer                    │
│  (TaskService, WorkflowService)             │
└─────────────────────────────────────────────┘
                      │
┌─────────────────────────────────────────────┐
│          Repository Adapters                 │
│  (TaskRepository, WorkflowLoader)           │
└─────────────────────────────────────────────┘
                      │
┌─────────────────────────────────────────────┐
│            Database Layer                    │
│         (SQLite + goose migrations)          │
└─────────────────────────────────────────────┘
```

### 依赖注入

- `domain` 包定义接口
- `repository` 包提供适配器实现
- `service` 包依赖接口而非具体实现
- 支持测试时使用 mock 实现

## Web API

### REST 端点

| 端点 | 方法 | 描述 |
|:---|:---|:---|
| `/api/tasks` | GET | 列出任务 |
| `/api/tasks` | POST | 创建任务 |
| `/api/tasks/:id` | GET | 获取任务详情 |
| `/api/tasks/:id/phases` | GET | 获取任务阶段 |
| `/api/tasks/:id/events` | GET | 获取任务事件 |
| `/api/tasks/:id/pause` | POST | 暂停任务 |
| `/api/tasks/:id/resume` | POST | 恢复任务 |
| `/api/tasks/:id/cancel` | POST | 取消任务 |
| `/api/phases/:id/complete` | POST | 完成阶段 |
| `/api/phases/:id/fail` | POST | 标记阶段失败 |
| `/api/workflows` | GET | 列出工作流 |
| `/api/workflows/:name` | GET | 获取工作流详情 |

### WebSocket 事件

| 事件类型 | 描述 |
|:---|:---|
| `task_created` | 任务创建 |
| `task_started` | 任务开始 |
| `task_completed` | 任务完成 |
| `task_failed` | 任务失败 |
| `task_paused` | 任务暂停 |
| `task_resumed` | 任务恢复 |
| `task_cancelled` | 任务取消 |
| `phase_started` | 阶段开始 |
| `phase_completed` | 阶段完成 |
| `phase_failed` | 阶段失败 |

## 测试

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/repository -v
go test ./internal/service -v

# 测试覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 当前覆盖率

- `internal/logger`: 85.0%
- `test`: 58.8%
- `internal/repository`: 52.3%

### 基准测试

```bash
# 运行基准测试
go test ./internal/workflow/... -bench=. -benchmem
go test ./internal/repository/... -bench=. -benchmem

# CPU profiling
go test ./internal/... -bench=. -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## 数据库

### 位置

```
~/.gclm-flow/gclm-engine.db
```

### Schema

```sql
-- 任务表
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    workflow_id TEXT,
    prompt TEXT,
    workflow_type TEXT,
    status TEXT,
    current_phase INTEGER,
    total_phases INTEGER,
    result TEXT,
    error_message TEXT,
    created_at TEXT,
    started_at TEXT,
    completed_at TEXT,
    updated_at TEXT
);

-- 阶段表
CREATE TABLE task_phases (
    id TEXT PRIMARY KEY,
    task_id TEXT,
    phase_name TEXT,
    display_name TEXT,
    sequence INTEGER,
    agent_name TEXT,
    model_name TEXT,
    status TEXT,
    input_prompt TEXT,
    output_text TEXT,
    error_message TEXT,
    started_at TEXT,
    completed_at TEXT,
    duration_ms INTEGER,
    created_at TEXT,
    updated_at TEXT,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- 事件表
CREATE TABLE events (
    id TEXT PRIMARY KEY,
    task_id TEXT,
    phase_id TEXT,
    event_type TEXT,
    event_level TEXT,
    data TEXT,
    occurred_at TEXT,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

-- 工作流表
CREATE TABLE workflows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE,
    workflow_type TEXT,
    display_name TEXT,
    description TEXT,
    version TEXT,
    config_yaml TEXT,
    is_builtin INTEGER,
    created_at TEXT,
    updated_at TEXT
);
```

## 工作流配置

### 内置工作流

- `analyze` - 代码分析、问题诊断、性能评估、安全审计
- `docs` - 文档编写、设计方案、需求分析
- `feat` - 新功能开发、模块开发、跨文件重构
- `fix` - Bug 修复、小修改、单文件变更

### 节点配置

| 字段 | 类型 | 必需 | 说明 |
|:---|:---|:---:|:---|
| `ref` | string | ✅ | 节点唯一标识符 |
| `display_name` | string | ✅ | 显示名称 |
| `agent` | string | ✅ | Agent 名称 |
| `model` | string | ✅ | 模型 (haiku/sonnet/opus) |
| `timeout` | int | ✅ | 超时时间（秒） |
| `required` | bool | ✅ | 是否必需 |
| `depends_on` | []string | ❌ | 依赖节点列表 |
| `parallel_group` | string | ❌ | 并行组标识 |

## 依赖

- `github.com/spf13/cobra` - CLI 框架
- `github.com/gin-gonic/gin` - HTTP 框架
- `github.com/gorilla/websocket` - WebSocket
- `github.com/mattn/go-sqlite3` - SQLite 驱动
- `github.com/pressly/goose/v3` - 数据库迁移
- `github.com/rs/zerolog` - 结构化日志
- `gopkg.in/yaml.v3` - YAML 解析

## License

MIT
