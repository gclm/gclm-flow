---
name: gclm
description: "智能分流工作流引擎 - 基于任务类型自动选择最优工作流，支持动态扩展工作流定义"
allowed-tools: [
  "Bash(gclm-engine *)",
  "Read(*)", "Write(*)", "Edit(*)",
  "Glob(*)", "Grep(*)", "Task(*)"
]
version: "3.0"
engine: "gclm-engine Go Engine"
---

## 核心设计

**动态工作流引擎** - 通过 Go 引擎 + YAML 配置实现可扩展的工作流系统：

```
用户请求 → Go 引擎智能分类 → 加载对应 YAML 工作流 → 执行阶段序列
```

### 关键特性

- **动态扩展**: 在 `workflows/` 添加 YAML 文件即可新增工作流
- **智能分流**: 引擎自动检测任务类型并选择工作流
- **状态持久化**: SQLite 数据库管理任务和阶段状态
- **Agent 编排**: 支持 Agent 定义、并行执行、依赖管理

---

## 工作流机制

### 1. 工作流定义（YAML）

工作流通过 `workflows/*.yaml` 文件定义：

```yaml
name: my_workflow
workflow_type: "MY_TYPE"  # 用于智能分类
display_name: "我的工作流"
description: "工作流描述"

nodes:
  - ref: discovery
    display_name: "需求发现"
    agent: investigator
    model: haiku
    timeout: 60
    required: true

  - ref: clarification
    depends_on: [discovery]
    display_name: "澄清确认"
    agent: investigator
    model: haiku

  # ... 更多节点
```

### 2. 智能分类

Go 引擎根据用户输入自动选择工作流：

| 机制 | 说明 |
|:---|:---|
| **关键词匹配** | `service/classifier.go` 中的关键词评分系统 |
| **手动指定** | 使用 `--workflow-type` 参数强制指定 |
| **默认回退** | 无匹配时使用默认工作流 |

**查看可用工作流**:
```bash
~/.gclm-flow/gclm-engine workflow list
```

### 3. 阶段依赖

通过 `depends_on` 定义阶段依赖关系，引擎自动计算执行顺序：

```yaml
nodes:
  - ref: phase_a
    # 无依赖，可立即执行

  - ref: phase_b
    depends_on: [phase_a]  # 等待 phase_a 完成

  - ref: phase_c
    depends_on: [phase_a, phase_b]  # 等待两者完成
```

**并行执行**: 无相互依赖的阶段会自动并行执行。

### 4. Agent 体系

Agent 在 `agents/*.md` 中定义，包含：
- 职责描述
- 推荐模型
- 输入输出格式
- 使用场景

**常用 Agent**:
- `investigator` - 探索分析 (Haiku)
- `architect` - 架构设计 (Opus)
- `worker` - 任务执行 (Sonnet)
- `tdd-guide` - TDD 指导 (Sonnet)
- `code-reviewer` - 代码审查 (Sonnet)

---

## 执行流程

### 初始化

```bash
# 创建任务，引擎自动选择工作流
~/.gclm-flow/gclm-engine workflow start "<任务描述>" --json

# 返回当前阶段信息
{
  "task_id": "task-xxx",
  "current_phase": { "phase_name": "discovery", ... },
  "total_phases": 5
}
```

### 阶段循环

```bash
# 1. 获取当前阶段
~/.gclm-flow/gclm-engine task current <task-id> --json

# 2. 执行阶段（调用相应 Agent 或 Task）

# 3. 完成阶段
~/.gclm-flow/gclm-engine task complete <task-id> <phase-id> --output "<输出>" --json

# 4. 重复步骤 1-3 直到所有阶段完成
```

### 状态查询

```bash
# 查看执行计划
~/.gclm-flow/gclm-engine task plan <task-id>

# 查看事件日志
~/.gclm-flow/gclm-engine task events <task-id>

# 列出所有任务
~/.gclm-flow/gclm-engine task list
```

---

## 状态管理

**重要**: 所有状态由 Go 引擎在 SQLite 数据库中维护。

- **数据库位置**: `~/.gclm-flow/gclm-engine.db`
- **无需手动状态文件**: 引擎自动管理任务、阶段、事件表

---

## 添加新工作流

1. 在 `workflows/` 创建 YAML 文件
2. 定义 `workflow_type` 用于智能分类
3. 添加阶段节点和依赖关系
4. 更新 `service/classifier.go` 中的关键词（可选）

示例：
```yaml
# workflows/review.yaml
name: review
workflow_type: "REVIEW"
display_name: "代码审查"
nodes:
  - ref: analyze
    agent: investigator
    model: haiku
  - ref: review
    depends_on: [analyze]
    agent: code-reviewer
    model: sonnet
```

---

## 硬约束

1. **Phase 0 强制**: 任何操作前先读取 llmdoc
2. **代码搜索分层回退**: auggie → llmdoc → Grep
3. **依赖优先**: 阶段必须满足依赖条件才能执行
4. **状态持久化**: 每个阶段后调用引擎更新状态

---

## 代码搜索

### 分层回退策略

| 方法 | 优势 | 劣势 | 状态 |
|:---|:---|:---|:---:|
| **auggie** | 语义搜索、自然语言查询 | 需要外部服务 | 推荐 |
| **llmdoc** | 结构化文档、本地 | 覆盖范围有限 | 默认 |
| **Grep** | 完整代码搜索 | 速度较慢 | 备选 |

### 安装 auggie（可选）

```bash
npm install -g @augmentcode/auggie@prerelease
```
