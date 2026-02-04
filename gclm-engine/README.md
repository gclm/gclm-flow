# gclm-engine

> **gclm-flow 工作流引擎** - Go 语言实现的智能工作流编排系统

## 项目概述

gclm-engine 是 gclm-flow 的核心引擎，负责：
- **工作流配置管理** - 基于 YAML 的可扩展工作流定义
- **任务状态管理** - SQLite 持久化，支持暂停/恢复
- **阶段调度** - 依赖解析、拓扑排序、并行执行
- **CLI 接口** - JSON 输出，与 Claude Code skills 集成

## 项目结构

```
gclm-engine/
├── main.go                 # 入口文件
├── go.mod                  # Go 模块定义
├── Makefile                # 构建脚本
├── internal/
│   ├── cli/               # CLI 命令实现 (cobra)
│   ├── db/                # 数据库层 (SQLite)
│   │   ├── database.go    # 数据库初始化
│   │   └── workflow.go    # 工作流存储
│   ├── pipeline/          # 流水线解析器
│   │   └── parser.go      # YAML 解析、依赖图
│   └── service/           # 任务服务
│       └── task.go        # 核心工作流逻辑
├── pkg/types/             # 共享类型定义
│   ├── types.go          # Task, Phase, Event
│   └── pipeline.go       # Pipeline, PipelineNode
├── workflows/             # 内置工作流 YAML (已移至项目根目录)
└── test/                  # 测试文件
```

## 快速开始

### 构建

```bash
cd gclm-engine
make build
# 或
go build -o gclm-engine
```

### 安装

```bash
# 安装到 ~/.gclm-flow/
make install

# 或使用项目根目录的安装脚本
cd .. && bash install.sh
```

### 运行

```bash
# 添加到 PATH
export PATH="$PATH:$HOME/.gclm-flow"

# 查看版本
gclm-engine version
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

# 导出工作流
gclm-engine workflow export <workflow-name>
```

### 任务管理

```bash
# 创建任务（自动检测工作流类型）
gclm-engine task create "实现用户登录功能"

# 指定工作流类型
gclm-engine task create "修复 bug" --workflow-type CODE_SIMPLE

# 一键开始工作流
gclm-engine workflow start "添加支付模块"

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
    pipeline_id TEXT,
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

### YAML 格式

工作流通过 YAML 文件定义，位于 `~/.gclm-flow/workflows/`：

```yaml
name: code_simple
display_name: "CODE_SIMPLE 工作流"
description: "Bug 修复、小修改的标准流程"
version: "0.2.0"
author: "gclm-flow"

workflow_type: "CODE_SIMPLE"

nodes:
  - ref: discovery
    display_name: "Discovery / 需求发现"
    agent: investigator
    model: haiku
    timeout: 60
    required: true

  - ref: clarification
    display_name: "Clarification / 澄清确认"
    agent: investigator
    model: haiku
    depends_on: [discovery]
    timeout: 60
    required: true

  - ref: tdd_red
    display_name: "TDD Red / 编写测试"
    agent: tdd-guide
    model: sonnet
    depends_on: [clarification]
    timeout: 120
    required: true

  - ref: tdd_green
    display_name: "TDD Green / 编写实现"
    agent: worker
    model: sonnet
    depends_on: [tdd_red]
    timeout: 180
    required: true

  - ref: code_reviewer
    display_name: "Code Reviewer / 代码审查"
    agent: code-reviewer
    model: sonnet
    depends_on: [tdd_green]
    timeout: 90
    required: true

  - ref: summary
    display_name: "Summary / 完成总结"
    agent: investigator
    model: haiku
    depends_on: [code_reviewer]
    timeout: 60
    required: true
```

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

## 工作流类型检测

引擎使用关键词评分系统自动检测工作流类型：

| 关键词类型 | 示例 | 分数 |
|:---|:---|:---:|
| 文档短语 | "编写文档", "方案设计", "架构设计" | +5 |
| 文档单词 | "文档", "方案", "需求", "分析" | +3 |
| Bug 短语 | "修复bug", "fix bug" | -5 |
| Bug 单词 | "bug", "修复", "debug" | -3 |
| 功能单词 | "功能", "模块", "开发" | -1 |

**阈值**：score ≥ 3 → DOCUMENT, score ≤ -3 → CODE_SIMPLE, 其他 → CODE_COMPLEX

## 开发

### 测试

```bash
# 运行所有测试
make test

# 运行特定包测试
go test ./internal/cli -v
go test ./internal/service -v

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 开发模式

```bash
# 快速构建到 ~/.gclm-flow/
make dev

# 或直接构建
go build -o ~/.gclm-flow/gclm-engine .
```

## 状态说明

### 任务状态

- `created` - 已创建
- `running` - 运行中
- `paused` - 已暂停
- `completed` - 已完成
- `failed` - 已失败
- `cancelled` - 已取消

### 阶段状态

- `pending` - 待执行
- `running` - 执行中
- `completed` - 已完成
- `failed` - 失败
- `skipped` - 已跳过

## 依赖

- `github.com/spf13/cobra` - CLI 框架
- `github.com/mattn/go-sqlite3` - SQLite 驱动 (需要 CGO)
- `gopkg.in/yaml.v3` - YAML 解析

## License

MIT
