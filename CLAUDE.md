# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作提供指导。

---

## 项目概述

**gclm-flow** 是一个 Claude Code 的智能工作流路由插件。它自动检测任务类型 (DOCUMENT / CODE_SIMPLE / CODE_COMPLEX) 并路由到合适的开发流程，遵循 SpecDD、TDD 和文档优先原则。

项目组成：
1. **插件 skills/agents** (`skills/`, `agents/`) - Claude Code 集成
2. **Go 引擎** (`gclm-engine/`) - 通过 SQLite 进行工作流编排和状态管理
3. **工作流定义** (`gclm-engine/workflows/`) - 基于 YAML 的流水线配置

---

## 快速开始命令

### Go 引擎 (gclm-engine)

```bash
# 构建和安装
cd gclm-engine
make build
make install          # 安装到 ~/.gclm-flow/

# 运行命令
~/.gclm-flow/gclm-engine version
~/.gclm-flow/gclm-engine workflow start "实现用户登录" --json
~/.gclm-flow/gclm-engine task current <task-id> --json
~/.gclm-flow/gclm-engine task complete <task-id> <phase-id> --output "结果"
```

### 开发

```bash
# 测试
make test              # 运行所有测试
cd gclm-engine && go test ./... -v

# 快速开发构建
make dev              # 直接构建到 ~/.gclm-flow/
```

---

## 架构

### 高层流程

```
用户请求 (/gclm 任务)
    ↓
Go 引擎: workflow start → 在 SQLite 中创建任务和阶段
    ↓
Skills 系统: 通过 task current 读取当前阶段
    ↓
Agent 执行: 运行相应的 agent (investigator/architect/worker 等)
    ↓
状态更新: task complete → 更新 SQLite 阶段状态
    ↓
下一阶段: 重复直到所有阶段完成
```

### 核心组件

| 组件 | 位置 | 用途 |
|:---|:---|:---|
| **Go 引擎 CLI** | `gclm-engine/internal/cli/` | 命令接口，为 skills 提供 JSON 输出 |
| **任务服务** | `gclm-engine/internal/service/task.go` | 核心工作流逻辑，阶段转换 |
| **数据库层** | `gclm-engine/internal/db/` | 任务/阶段/事件的 SQLite 持久化 |
| **流水线解析器** | `gclm-engine/internal/pipeline/` | YAML 工作流解析，依赖解析 |
| **工作流 YAML** | `gclm-engine/workflows/` | 定义 DOCUMENT、CODE_SIMPLE、CODE_COMPLEX 流程 |
| **Skills** | `skills/gclm/SKILL.md` | 编排工作流的主 skill |
| **Agents** | `agents/*.md` | Agent 定义 (investigator、architect、tdd-guide 等) |

### 数据库结构

位于 `~/.gclm-flow/gclm-engine.db`：

- **tasks**: id, pipeline_id, prompt, workflow_type, status, current_phase, total_phases
- **task_phases**: id, task_id, phase_name, agent, model, status, output_text
- **events**: id, task_id, phase_id, event_type, data (审计日志)

---

## 工作流配置

工作流在 `gclm-engine/workflows/` 中通过 YAML 定义：

```yaml
name: code_simple
workflow_type: "CODE_SIMPLE"
nodes:
  - ref: discovery
    display_name: "Discovery / 需求发现"
    agent: investigator
    model: haiku
    timeout: 60
    required: true
  - ref: clarification
    depends_on: [discovery]
    # ... 更多节点
```

### 工作流类型

| 类型 | 触发关键词 | 阶段 |
|:---|:---|:---|
| **DOCUMENT** | "文档", "方案", "设计", "需求" | Draft → Refine → Review |
| **CODE_SIMPLE** | "bug", "修复", "fix error" | TDD Red → TDD Green |
| **CODE_COMPLEX** | "功能", "模块", "开发" | 完整 SpecDD + 架构阶段 |

### 添加新工作流

1. 在 `gclm-engine/workflows/` 中创建 YAML 文件或使用 `workflow install <path>`
2. 使用 `depends_on` 定义节点依赖
3. 使用 `parallel_group` 实现并行执行
4. 用 `required: true` 标记关键节点

---

## Agent 体系

| Agent | 模型 | 阶段 | 用途 |
|:---|:---|:---|:---|
| `investigator` | Haiku | 1, 2, 9 | 快速代码库调查 |
| `architect` | Opus | 4 | 架构设计决策 |
| `spec-guide` | Opus | 5 | SpecDD 文档编写 |
| `tdd-guide` | Sonnet | 6 | TDD 测试编写指导 |
| `worker` | Sonnet | 7 | 代码实现 |
| `code-simplifier` | Sonnet | 8 (并行) | 代码重构 |
| `security-guidance` | Sonnet | 8 (并行) | 安全审查 |
| `code-reviewer` | Sonnet | 8 (并行) | 代码审查 |

### 自定义 Agent

在 `agents/<name>.md` 中定义，使用 YAML frontmatter：

```yaml
---
name: my-agent
description: "Agent 用途"
tools: ["Read", "Write", "Grep"]
model: sonnet
color: blue
permission: auto
---
```

---

## Skills 集成

主 skill: `skills/gclm/SKILL.md`

**关键集成点：**
- `workflow start <prompt>` → 创建任务，返回第一阶段
- `task current <task-id>` → 获取下一个待执行阶段
- `task complete <task-id> <phase-id> --output "..."` → 标记阶段完成
- `task export <task-id> <file>` → 导出状态到 markdown (兼容性)

---

## 约定规范

### 工作流类型检测

关键词评分系统 (位于 `service/task.go`)：
- 文档短语 (+5): "编写文档", "方案设计", "架构设计"
- 文档单词 (+3): "文档", "方案", "需求"
- Bug 修复短语 (-5): "修复bug", "fix bug"
- Bug 修复单词 (-3): "bug", "修复", "debug"
- 功能开发单词 (-1): "功能", "模块", "开发"

阈值：score >= 3 → DOCUMENT, score <= -3 → CODE_SIMPLE, 其他 → CODE_COMPLEX

### 阶段状态

`pending` → `running` → `completed` / `failed` / `skipped`

### 错误处理

- 必需阶段失败 → 任务失败
- 非必需阶段失败 → 继续下一阶段
- 通过工作流 YAML 中的 `required: true/false` 配置

---

## 目录结构

```
gclm-flow/
├── agents/                    # Agent 定义
├── skills/                    # Skill 定义
├── rules/                     # 工作流规则 (phases, tdd, spec)
├── gclm-engine/
│   ├── main.go               # 入口文件
│   ├── internal/
│   │   ├── cli/              # CLI 命令 (cobra)
│   │   ├── db/               # SQLite 操作
│   │   ├── pipeline/         # YAML 解析器
│   │   └── service/          # 任务服务 (工作流逻辑)
│   ├── pkg/types/            # 共享类型
│   ├── workflows/            # 内置工作流 YAML
│   └── Makefile
└── workflows/examples/        # 自定义工作流示例
```

---

## 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
cd gclm-engine && go test ./internal/cli -v
cd gclm-engine && go test ./internal/service -v

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 发布流程

1. 更新 `gclm-engine/internal/cli/commands.go` 中的版本 (createVersionCommand)
2. 创建 git 标签: `git tag v0.x.x`
3. 推送标签: `git push origin v0.x.x`
4. GitHub Actions 构建 darwin/linux amd64/arm64 二进制文件
5. 发布包含：二进制文件、workflows.tar.gz、install.sh、checksums.txt

---

## 重要约束

1. **SQLite 单写入者**: 数据库使用 `SetMaxOpenConns(1)` - SQLite 限制
2. **WAL 模式**: 启用以提升并发性 (`_pragma=journal_mode(WAL)`)
3. **工作流状态**: 存储在 `~/.gclm-flow/gclm-engine.db`
4. **JSON 输出**: 所有引擎命令支持 `--json` 标志以便 skill 集成
5. **阶段依赖**: 必须形成 DAG - 加载时检查循环依赖

---

## 依赖项

- `github.com/spf13/cobra` - CLI 框架
- `github.com/mattn/go-sqlite3` - SQLite 驱动 (需要 CGO)
- `gopkg.in/yaml.v3` - YAML 解析

**构建注意**: SQLite 必须启用 CGO (GitHub Actions 中 `CGO_ENABLED=1`)
