# 系统架构

## 概述

gclm-flow 采用 Go 引擎 + YAML 工作流 + 多 Agent 的架构设计。

```
┌─────────────────────────────────────────────────────────────┐
│                        用户请求                              │
│                  (/gclm <任务描述>)                          │
└────────────────────────────┬────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                      Skills 层                               │
│                   skills/gclm/SKILL.md                       │
│              - 接收用户请求                                  │
│              - 调用 gclm-engine 命令                         │
└────────────────────────────┬────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                     Go 引擎层                                │
│                  gclm-engine/                                │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ CLI (internal/cli/)                                   │  │
│  │  - workflow start <prompt>                           │  │
│  │  - task current <task-id>                            │  │
│  │  - task complete <task-id> <phase-id> --output "..."  │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 任务服务 (internal/service/)                           │  │
│  │  - 智能分流 (关键词评分)                                │  │
│  │  - 阶段管理 (pending → running → completed)              │  │
│  │  - 工作流加载 (YAML 解析)                              │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 数据库 (internal/db/)                                  │  │
│  │  - SQLite (WAL 模式)                                  │  │
│  │  - tasks / task_phases / events                       │  │
│  └───────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                    工作流配置层                               │
│                     workflows/                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │  document   │  │code_simple  │  │code_complex │         │
│  │   .yaml     │  │   .yaml     │  │   .yaml     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└────────────────────────────┬────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                     Agent 执行层                              │
│                      agents/                                │
│  ┌───────────┐  ┌──────────┐  ┌─────────┐  ┌──────────┐     │
│  │investigator│  │architect │  │ tdd-guide│  │  worker  │ ... │
│  └───────────┘  └──────────┘  └─────────┘  └──────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## 核心流程

### 1. 工作流启动

```
用户请求 → Skills → gclm-engine workflow start
                              ↓
                        智能分流 (关键词评分)
                              ↓
                    ┌─────────┼─────────┐
                    ↓         ↓         ↓
                DOCUMENT  CODE_SIMPLE  CODE_COMPLEX
                    ↓         ↓         ↓
                加载对应 YAML 工作流配置
                              ↓
                        创建 Task 和 Phase 记录
                              ↓
                        返回第一阶段信息
```

### 2. 阶段执行

```
Skills → gclm-engine task current <task-id>
           ↓
        查询下一个待执行 Phase (pending)
           ↓
        返回 Phase 信息 (agent, model, prompt)
           ↓
Skills → 调用对应 Agent
           ↓
Agent → 执行任务
           ↓
Skills → gclm-engine task complete <task-id> <phase-id> --output "..."
           ↓
        更新 Phase 状态 (running → completed)
           ↓
        检查下一阶段 (并行组处理)
```

---

## 关键组件

### Go 引擎 (gclm-engine/)

| 组件 | 文件 | 职责 |
|:---|:---|:---|
| **CLI** | `internal/cli/commands.go` | 命令接口，JSON 输出 |
| **任务服务** | `internal/service/task.go` | 智能分流、阶段管理 |
| **流水线解析** | `internal/pipeline/parser.go` | YAML 解析、依赖检查 |
| **数据库** | `internal/db/database.go` | SQLite 操作 |
| **类型定义** | `pkg/types/` | 共享数据结构 |

### 工作流配置 (workflows/)

| 文件 | workflow_type | 阶段数 |
|:---|:---|:---:|
| `document.yaml` | `DOCUMENT` | 7 |
| `code_simple.yaml` | `CODE_SIMPLE` | 6 |
| `code_complex.yaml` | `CODE_COMPLEX` | 9 |

### Skills 集成 (skills/gclm/)

| 命令 | 功能 |
|:---|:---|
| `workflow start <prompt>` | 创建任务，返回第一阶段 |
| `task current <task-id>` | 获取下一阶段 |
| `task complete ... --output "..."` | 完成阶段 |

---

## 数据流

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户请求   │────▶│  Skills 层   │────▶│  Go 引擎层   │
└──────────────┘     └──────────────┘     └──────────────┘
                                                  ↓
                                            ┌──────────────┐
                                            │   SQLite DB   │
                                            │ ~/.gclm-flow/ │
                                            │ gclm-engine.db│
                                            └──────────────┘
                                                  ↑
                                            ┌──────────────┐
                                            │ workflows/    │
                                            │ *.yaml        │
                                            └──────────────┘
```

---

## 并行执行

工作流支持通过 `parallel_group` 实现阶段并行执行：

```yaml
nodes:
  - ref: review_1
    parallel_group: review  # 与同组节点并行
  - ref: review_2
    parallel_group: review  # 与 review_1 并行
  - ref: review_3
    parallel_group: review  # 与 review_1, review_2 并行
```

Go 引擎通过拓扑排序计算执行顺序，识别并行组。
