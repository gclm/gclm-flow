# 系统架构

## 概述

gclm-flow 采用 Go 引擎 + 草稿/正式分离模型 + 多 Agent 的架构设计。

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
│  │  - 工作流加载 (从数据库加载)                            │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 资源嵌入 (internal/assets/)                            │  │
│  │  - migrations/*.sql                                    │  │
│  │  - workflows/*.yaml                                    │  │
│  │  - gclm_engine_config.yaml                             │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ 数据库 (internal/db/)                                  │  │
│  │  - SQLite (WAL 模式)                                  │  │
│  │  - tasks / task_phases / events                       │  │
│  └───────────────────────────────────────────────────────┘  │
└────────────────────────────┬────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                    工作流配置层 (草稿/正式分离)                │
└────────────────────────────┬────────────────────────────────┘
                             │
         ┌───────────────────┴───────────────────┐
         ↓                                       ↓
┌─────────────────────┐              ┌─────────────────────┐
│   草稿 (YAML)       │              │   正式 (数据库)     │
│  workflows/         │  ← sync →   │  workflows 表      │
│  *.yaml 文件        │              │  (config_yaml)      │
└─────────────────────┘              └─────────────────────┘
         │                                       │
         └───────────────────┬───────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                    内置工作流 (嵌入二进制)                    │
│                  gclm-engine/workflows/                     │
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
| **工作流解析** | `internal/workflow/parser.go` | YAML 解析、依赖检查 |
| **数据库** | `internal/db/database.go` | SQLite + Goose 迁移 |
| **资源嵌入** | `internal/assets/embed.go` | 资源嵌入/导出 |
| **类型定义** | `pkg/types/` | 共享数据结构 |

### 工作流配置 (草稿/正式分离)

**工作流来源**:
1. **内置**: `gclm-engine/workflows/*.yaml` (嵌入二进制)
2. **草稿**: `~/.gclm-flow/workflows/*.yaml` (用户可编辑)
3. **正式**: `workflows` 表 (数据库，通过 `workflow sync` 发布)

**同步命令**:
```bash
# 初始化 (从内置导出到草稿目录)
gclm-engine init

# 同步 (草稿 → 正式)
gclm-engine workflow sync                    # 同步所有
gclm-engine workflow sync workflows/feat.yaml # 同步单个
```

**内置工作流**:

| 文件 | workflow_type | 阶段数 |
|:---|:---|:---:|
| `document.yaml` | `DOCUMENT` | 7 |
| `code_simple.yaml` | `CODE_SIMPLE` | 6 |
| `code_complex.yaml` | `CODE_COMPLEX` | 9 |
| `analyze.yaml` | `ANALYZE` | 5 |

### Skills 集成 (skills/gclm/)

| 命令 | 功能 |
|:---|:---|
| `workflow list` | 列出所有工作流 (从数据库)
| `workflow start <prompt> --workflow <name>` | 创建任务，返回第一阶段 |
| `task current <task-id>` | 获取下一阶段 |
| `task complete ... --output "..."` | 完成阶段 |
| `task export <task-id> <file>` | 导出状态文件 (兼容旧版) |

**新增命令**:
| 命令 | 功能 |
|:---|:---|
| `workflow sync [yaml-file]` | 同步工作流 YAML 到数据库 |
| `workflow validate <yaml-file>` | 验证工作流配置 |
| `workflow install <yaml-file>` | 安装自定义工作流 |
| `workflow uninstall <name>` | 卸载自定义工作流 |
| `workflow info <name>` | 显示工作流详细信息 |
| `workflow export <name> [file]` | 导出工作流到 YAML |

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
                                            │              │
                                            │ - workflows  │ ← 草稿/正式分离
                                            │ - tasks      │
                                            │ - task_phases│
                                            │ - events     │
                                            └──────────────┘
                                                  ↑
                                            ┌──────────────┐
                                            │ workflows/    │ ← 草稿目录
                                            │ *.yaml        │
                                            └──────────────┘
                                                  ↑
                                            ┌──────────────┐
                                            │ 内置工作流    │ ← 嵌入二进制
                                            │ (embed.FS)   │
                                            └──────────────┘
```

---

## 草稿/正式分离模型

gclm-flow 引入了草稿/正式数据分离模型，确保工作流变更的可控性：

```
开发流程:
1. 编辑 ~/.gclm-flow/workflows/*.yaml (草稿)
2. 测试工作流: gclm-engine workflow validate <file>
3. 发布到正式: gclm-engine workflow sync
4. 正式环境从数据库加载工作流
```

**优势**:
- 草稿修改不影响运行中的系统
- 支持版本控制 (git)
- 可回滚 (保留多个版本)
- 零依赖部署 (内置工作流嵌入二进制)

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
