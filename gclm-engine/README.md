# gclm-engine PoC

> **Proof of Concept** - gclm-flow Go 引擎原型验证

## 项目概述

gclm-engine 是 gclm-flow 工作流的 Go 语言实现，旨在替代现有的 Shell 脚本方案，提供更健壮的状态管理、更好的并发支持和可扩展的流水线配置。

**PoC 目标**: 验证核心架构可行性，包括：
- SQLite 状态管理
- YAML 流水线配置
- 节点依赖解析和拓扑排序
- CLI 基础命令

## 项目结构

```
gclm-engine/
├── main.go                 # 入口文件
├── go.mod                  # Go 模块定义
├── internal/
│   ├── cli/               # CLI 命令实现
│   ├── db/                # 数据库层 (SQLite)
│   ├── pipeline/          # 流水线解析器
│   └── ...
├── pkg/types/             # 共享类型定义
├── configs/pipelines/     # 流水线 YAML 配置
│   └── code_simple.yaml   # CODE_SIMPLE 流水线
└── test/                  # 测试文件
```

## 快速开始

### 安装依赖

```bash
cd gclm-engine
go mod download
```

### 构建项目

```bash
go build -o gclm-engine
```

### 运行命令

```bash
# 查看版本
./gclm-engine version

# 列出所有流水线
./gclm-engine pipeline list

# 查看流水线详情
./gclm-engine pipeline get code_simple

# 推荐流水线
./gclm-engine pipeline recommend "修复登录页面样式问题"

# 创建任务
./gclm-engine task create "修复登录按钮颜色错误" --workflow-type CODE_SIMPLE

# 查看任务详情
./gclm-engine task get <task-id>

# 列出任务
./gclm-engine task list

# 查看任务阶段
./gclm-engine task phases <task-id>

# 查看任务事件
./gclm-engine task events <task-id>
```

## 数据库

PoC 使用 SQLite 数据库，默认位置：

```
~/.gclm-flow/gclm-engine.db
```

### 数据库 Schema

```sql
-- 任务表
tasks (id, pipeline_id, prompt, workflow_type, status, ...)

-- 阶段表
task_phases (id, task_id, phase_name, agent_name, status, ...)

-- 事件表
events (id, task_id, event_type, event_level, data, ...)
```

## 流水线配置

流水线使用 YAML 配置，示例：

```yaml
# configs/pipelines/code_simple.yaml
name: code_simple
display_name: "CODE_SIMPLE 工作流"
workflow_type: "CODE_SIMPLE"
version: "0.1.0-poc"

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

  # ... more nodes
```

### 节点配置项

| 字段 | 类型 | 说明 |
|:---|:---|:---|
| `ref` | string | 节点唯一标识符 |
| `display_name` | string | 显示名称 |
| `agent` | string | Agent 名称 |
| `model` | string | 模型名称 (haiku/sonnet/opus) |
| `timeout` | int | 超时时间（秒） |
| `required` | bool | 是否必需 |
| `depends_on` | []string | 依赖节点列表 |
| `parallel_group` | string | 并行组标识 |

## PoC 范围

### 包含的功能

| 功能 | 状态 |
|:---|:---:|
| SQLite 数据库 | ✅ |
| 任务 CRUD | ✅ |
| 阶段追踪 | ✅ |
| 事件日志 | ✅ |
| YAML 解析 | ✅ |
| 依赖图构建 | ✅ |
| 拓扑排序 | ✅ |
| CLI 命令 | ✅ |
| 工作流检测 | ✅ |

### 不包含的功能

| 功能 | 计划阶段 |
|:---|:---:|
| 流水线执行器 | Phase 1 |
| Agent 调用器 | Phase 1 |
| 并行执行 | Phase 2 |
| Web UI | Phase 4 |
| WebSocket | Phase 4 |
| REST API | Phase 3 |

## 开发状态

- **版本**: 0.1.0-poc
- **状态**: 开发中
- **最后更新**: 2026-02-04

## 下一步

1. ✅ 项目结构
2. ✅ 数据库 Schema
3. ✅ 类型定义
4. ✅ CLI 命令
5. ⏳ 单元测试
6. ⏳ E2E 测试
7. ⏳ 流水线执行器
8. ⏳ Agent 集成

## License

MIT
