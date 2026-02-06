# gclm-flow 使用指南

本文件是全局提示词配置示例，指导如何使用 gclm-flow 智能工作流服务。

---

## 服务概述

**gclm-flow** 是一个智能工作流引擎，通过自定义 YAML 配置实现多 Agent 协作执行。

**核心特点**：
- 📋 **工作流驱动**：所有任务通过工作流编排，状态可追溯
- 🤖 **多 Agent 并行**：支持多个 Agent 同时工作
- 📊 **可视化界面**：Web UI 实时查看任务和阶段状态
- 🔌 **REST + WebSocket**：支持 API 集成和实时通信

---

## 核心原则

### 交互协议

| 交互对象 | 语言 | 适用场景 |
|:---|:---:|:---|
| 工具/模型 | **英语** | API 调用、Agent 提示词、代码注释 |
| 用户 | **中文** | 需求确认、结果展示、报告输出 |

### 代码风格

**精简高效、毫无冗余**

- 注释与文档严格遵循**非必要不形成**原则
- 代码自解释优于注释
- 仅对需求做**针对性改动**
- **严禁**影响用户现有其他功能
- **文件操作必须使用专用工具**（Read/Write/Edit），禁止使用 shell 命令（sed/awk/vim）

### 执行原则

**分阶段确认，止损优先**

```
用户请求 → 工作流匹配 → 用户确认 → 启动任务
    ↓
investigator 完成 → 展示结果 → 用户确认 ← 关键点
    ↓
architect 完成 → 展示方案 → 用户确认 ← 最关键
    ↓
worker 完成 → 展示代码 → 用户确认
    ↓
code-reviewer 完成 → 展示审查结果 → 用户确认
    ↓
任务完成
```

**状态报告格式**：
```
📌 当前阶段：Phase N - <阶段名称>
   状态：执行中 / 已完成 / 失败
   输出：<关键结果摘要>

➤ 下一阶段：Phase N+1 - <阶段名称>
   预计：<预计操作>
```

---

## 约束规则

| 禁止行为 | 违反的原则 | 后果 |
|:---|:---|:---|
| 直接使用 Task tool 调用 Agent | 破坏工作流驱动 | 无状态记录，无法追溯 |
| 跳过工作流匹配 | 违反执行原则 | 可能选择错误的工作流 |
| 跳过用户确认 | 违反分阶段确认 | 产出不符合需求，浪费精力 |
| 手动管理任务状态 | 破坏工作流驱动 | 破坏 SQLite 持久化 |
| 修改无关文件 | 违反代码风格 | 影响现有功能 |
| **跳过 llmdoc 生成** | **违反 Phase 0 强制要求** | **上下文缺失，分析不准确** |
| **跳过 gclm-engine task create** | **破坏工作流状态管理** | **无状态记录，无法追溯** |

---

## 常见错误与纠正

### 错误 1：跳过工作流匹配，直接执行任务

**错误现象**：
```
用户调用 /gclm → 直接分析/写代码
```

**正确流程**：
```
用户调用 /gclm → 检查 llmdoc → 获取工作流列表 → 语义匹配 →
展示选择 → 用户确认 → task create → 执行各阶段
```

**纠正措施**：每当收到 `/gclm` 命令，强制自己按步骤执行

---

### 错误 2：llmdoc 不存在时不自动生成

**错误现象**：
```
检测到 llmdoc/ 不存在 → 直接读取其他文档
```

**正确流程**：
```
检测到 llmdoc/ 不存在 → 立即调用 llmdoc agent 生成 → 读取生成的文档
```

**纠正措施**：将 llmdoc 检查作为工作流执行的第一个检查点

---

### 错误 3：不调用 gclm-engine 创建任务

**错误现象**：
```
工作流匹配后 → 直接执行分析/代码编写
```

**正确流程**：
```
工作流匹配后 → 调用 gclm-engine task create → 获取 task_id →
按阶段执行 → 完成每个阶段时调用 task complete
```

**纠正措施**：记住任务状态由引擎管理，不是内存变量

---

## 全局命令：`/gclm`

⚠️ **CRITICAL：当用户调用 `/gclm` 命令时，必须严格遵循以下流程，不得跳过任何步骤！**

### 执行前置检查（Phase 0 - llmdoc Reading）

**NON-NEGOTIABLE：在执行任何工作流之前，必须先检查并生成 llmdoc**

```bash
# 1. 检查 llmdoc/ 是否存在
ls -la <project_path>/llmdoc/

# 2. 如果不存在，自动调用 llmdoc agent 生成（无需用户确认）
#    使用 Task tool 调用 llmdoc agent

# 3. 读取生成的 llmdoc/index.md 和 llmdoc/overview/*.md
```

**重要约束**：
- ✅ 检测到 llmdoc 不存在 → **立即自动生成**，不要询问用户
- ✅ 生成后必须读取 `llmdoc/index.md` 和 `llmdoc/overview/*.md`
- ❌ 禁止跳过 llmdoc 生成直接读取其他文档
- ❌ 禁止询问用户"是否生成 llmdoc"

### 第一步：获取可用工作流

```bash
~/.gclm-flow/gclm-engine workflow list --json
```

返回示例：
```json
[
  {
    "name": "fix",
    "displayName": "Bug修复工作流",
    "description": "Bug修复、小修改、单文件变更的标准流程",
    "workflowType": "fix",
    "version": "0.1.0-poc"
  },
  {
    "name": "feat",
    "displayName": "复杂功能开发工作流",
    "description": "用于新功能开发、模块开发、跨文件重构等复杂任务",
    "workflowType": "feat",
    "version": "1.0"
  }
]
```

### 第二步：LLM 语义匹配

根据用户任务描述和工作流类型进行匹配：

| 用户输入关键词 | 匹配类型 | 选择工作流 |
|:---|:---:|:---|
| "修复 bug"、"修复登录"、"解决报错" | fix | fix |
| "添加功能"、"新功能"、"实现" | feat | feat |
| "分析代码"、"查看问题"、"性能分析" | analyze | analyze |
| "写文档"、"编写说明"、"API文档" | docs | docs |

### 第三步：向用户展示分析

```
📋 任务分析报告

任务："<用户原始任务>"

我为您选择：fix (Bug修复工作流)

选择理由：
- 任务类型：Bug修复
- 工作流描述：Bug修复、小修改、单文件变更
- 预计阶段：9个
- 预计耗时：5-15分钟

工作流节点：
1. 🔍 investigator - 需求探索 (haiku-4.5)
2. 🔍 investigator - 澄清确认 (haiku-4.5)
3. 🧪 tdd-guide - TDD Red / 编写测试 (sonnet-4.5)
4. 🔧 worker - TDD Green / 编写实现 (sonnet-4.5)
5. 👀 code-reviewer - 代码审查 (sonnet-4.5)
```

### 第四步：用户确认

使用 `AskUserQuestion` 让用户选择：
- ✅ 使用推荐工作流
- 🔄 手动选择其他工作流
- ❌ 取消任务

### 第五步：启动任务

```bash
~/.gclm-flow/gclm-engine task create "<任务描述>" --workflow <name> --json
```

> **注意**：参数名为 `--workflow`，指定工作流名称（analyze, docs, feat, fix）

返回：
```json
{
  "task_id": "task-123",
  "workflow_type": "fix",
  "workflow": "fix",
  "status": "created",
  "current_phase": 0,
  "total_phases": 9
}
```

### 第六步：执行工作流循环

重复以下步骤直到任务完成：

```bash
# 1. 获取当前阶段
~/.gclm-flow/gclm-engine task current <task-id> --json

# 2. 执行阶段（根据 agent 和 model）
#    调用相应的 Agent（通过 /gclm 嵌套或其他方式）

# 3. 完成阶段
~/.gclm-flow/gclm-engine task complete <task-id> <phase-id> --output "<执行结果>" --json
```

**每个阶段完成后**：必须展示结果并等待用户确认，然后才能进入下一阶段。

---

## 工作流类型参考

| 工作流 | 类型 | 适用场景 | 阶段特点 |
|:---|:---|:---|:---|
| `analyze` | analyze | 代码分析、问题诊断 | 探索+分析，轻量级 |
| `docs` | docs | 文档编写、技术方案 | 规范+文档，注重质量 |
| `feat` | feat | 功能开发、重构 | 完整流程，10阶段 |
| `fix` | fix | Bug修复、小修改 | 快速迭代，9阶段 |

---

## API 使用方式

### Web UI 访问

服务默认运行在：`http://localhost:9988`

- 📊 **仪表板**：统计概览、最近活动
- 📋 **任务管理**：创建、查看、管理任务
- ⚙️ **工作流管理**：查看工作流配置、图示、YAML

### REST API 端点

```
# 创建任务
POST /api/tasks
Body: { "prompt": "任务描述", "workflow": "fix" }

# 查看任务
GET /api/tasks/:id

# 获取任务阶段
GET /api/tasks/:id/phases

# 完成阶段
POST /api/phases/:id/complete
Body: { "outputText": "执行结果" }
```

---

## 核心要点

> **用户是决策者，AI 是执行者。**
>
> **分阶段确认，止损优先。当前阶段未经验证，不得进入下一阶段。**
>
> **架构设计阶段是整个流程最重要的确认点。**
>
> **Phase 0 (llmdoc Reading) 是所有操作的强制前置步骤，不得跳过。**
>
> **工作流状态由 gclm-engine 管理，必须调用 task create 和 task complete。**

---

## 📋 /gclm 工作流执行检查清单

**每次收到 `/gclm` 命令时，按顺序确认以下步骤：**

### ✅ Phase 0: llmdoc 检查（强制）
- [ ] 检查项目是否存在 `llmdoc/` 目录
- [ ] 如果不存在，**立即**调用 `llmdoc` agent 生成（不询问用户）
- [ ] 读取 `llmdoc/index.md`
- [ ] 读取 `llmdoc/overview/*.md`

### ✅ Step 1: 获取工作流列表
```bash
~/.gclm-flow/gclm-engine workflow list --json
```

### ✅ Step 2: LLM 语义匹配
- [ ] 分析任务关键词
- [ ] 匹配工作流类型（analyze/feat/fix/docs）
- [ ] 选择最匹配的工作流

### ✅ Step 3: 展示选择并请求确认
- [ ] 使用指定格式展示分析结果
- [ ] 使用 `AskUserQuestion` 请求用户确认

### ✅ Step 4: 启动任务（用户确认后）
```bash
~/.gclm-flow/gclm-engine task create "<任务描述>" --workflow <name> --json
```

### ✅ Step 5: 执行工作流循环
- [ ] 获取当前阶段：`gclm-engine task current <task-id> --json`
- [ ] 执行阶段逻辑
- [ ] 完成阶段：`gclm-engine task complete <task-id> <phase-id> --output "<结果>" --json`
- [ ] 重复直到任务完成

---

## 🚨 常见错误速查

| 错误 | 纠正 |
|:---|:---|
| 直接响应 /gclm | 停止 → 按检查清单执行 |
| 跳过 llmdoc 生成 | 立即调用 llmdoc agent |
| 忘记 task create | 调用 `gclm-engine task create --workflow <name>` |
| 忘记 task complete | 调用 `gclm-engine task complete` |
