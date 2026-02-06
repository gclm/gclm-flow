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

---

## 全局命令：`/gclm`

当用户调用 `/gclm` 命令时，必须严格遵循以下流程：

### 第一步：获取可用工作流

```bash
~/.gclm-flow/gclm-engine workflow list --json
```

返回示例：
```json
{
  "workflows": [
    {
      "name": "code_simple",
      "displayName": "CODE_SIMPLE 工作流",
      "workflowType": "fix",
      "description": "Bug修复、小修改、单文件变更的标准流程"
    },
    {
      "name": "code_complex",
      "displayName": "复杂功能开发工作流",
      "workflowType": "feat",
      "description": "用于新功能开发、模块开发、跨文件重构等复杂任务"
    }
  ]
}
```

### 第二步：LLM 语义匹配

根据用户任务描述和工作流类型进行匹配：

| 用户输入关键词 | 匹配类型 | 选择工作流 |
|:---|:---|:---|
| "修复 bug"、"修复登录"、"解决报错" | fix | code_simple |
| "添加功能"、"新功能"、"实现" | feat | code_complex |
| "分析代码"、"查看问题"、"性能分析" | analyze | analyze |
| "写文档"、"编写说明"、"API文档" | docs | document |

### 第三步：向用户展示分析

```
📋 任务分析报告

任务："<用户原始任务>"

我为您选择：code_simple (CODE_SIMPLE 工作流)

选择理由：
- 任务类型：Bug修复
- 工作流描述：Bug修复、小修改、单文件变更
- 预计阶段：3个
- 预计耗时：5-15分钟

工作流节点：
1. 🔍 investigator - 需求探索 (haiku-4.5)
2. 🔧 worker - 代码实现 (sonnet-4.5)
3. 👀 code-reviewer - 代码审查 (sonnet-4.5)
```

### 第四步：用户确认

使用 `AskUserQuestion` 让用户选择：
- ✅ 使用推荐工作流
- 🔄 手动选择其他工作流
- ❌ 取消任务

### 第五步：启动任务

```bash
~/.gclm-flow/gclm-engine task create "<任务描述>" --workflow <workflow_name> --json
```

返回：
```json
{
  "task": {
    "id": "task-123",
    "prompt": "<任务描述>",
    "workflowType": "fix",
    "status": "running",
    "currentPhase": 0,
    "totalPhases": 3
  }
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
| `code_simple` | fix | Bug修复、小修改 | 3阶段，快速迭代 |
| `code_complex` | feat | 功能开发、重构 | 6+阶段，完整流程 |
| `analyze` | analyze | 代码分析、问题诊断 | 探索+分析，轻量级 |
| `document` | docs | 文档编写、技术方案 | 规范+文档，注重质量 |

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
Body: { "prompt": "任务描述", "workflowType": "fix" }

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
