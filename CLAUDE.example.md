# gclm-flow 项目配置

本文件为 Claude Code (claude.ai/code) 在此代码库中工作提供指导。

---

## 项目概述

**gclm-flow** 是一个基于 Go 引擎的智能工作流系统，支持自定义工作流 YAML 配置和多 Agent 并行执行。

```
用户请求 → gclm-engine (Go 引擎) → 工作流编排 → Agent 执行
    ↓              ↓                    ↓
 自然语言    SQLite 状态管理      多 Agent 并行
```

---

## 强制执行流程

**当用户调用 `/gclm` 命令时，必须严格遵循以下流程：**

### 1. 获取工作流列表 (NON-NEGOTIABLE)

```bash
~/.gclm-flow/gclm-engine workflow list --json
```

### 2. LLM 语义匹配

根据用户任务描述和工作流描述进行匹配：

| 用户输入示例 | 关键词 | 匹配工作流 |
|:---|:---|:---|
| "分析项目不足" | 分析、不足 | `analyze` |
| "修复登录 bug" | 修复、bug | `code_simple` |
| "添加新功能" | 添加、功能 | `code_complex` |
| "编写文档" | 编写、文档 | `document` |

### 3. 向用户展示选择

```
📋 工作流选择分析

根据您的任务: "<用户任务>"

我为您选择: <workflow_name> (<display_name>)

选择理由:
- 任务关键词: <关键词>
- 工作流描述: <描述>
- 类型匹配: <workflow_type>
- 阶段数: <n>
```

### 4. 用户确认

使用 `AskUserQuestion` 让用户确认：
- 使用推荐工作流
- 手动选择其他工作流
- 取消

### 5. 启动工作流

```bash
~/.gclm-flow/gclm-engine workflow start "<任务描述>" --workflow <name> --json
```

### 6. 阶段循环

重复以下步骤直到完成：
1. 获取当前阶段：`task current <task-id> --json`
2. 执行阶段（根据 agent 和 model 调用相应 Agent）
3. 完成阶段：`task complete <task-id> <phase-id> --output "..." --json`

---

## 禁止行为

| 禁止行为 | 原因 |
|:---|:---|
| **直接使用 Task tool 调用 Agent** | 绕过工作流引擎，无状态记录 |
| **跳过工作流列表获取** | 无法进行语义匹配 |
| **跳过用户确认** | 用户失去控制权 |
| **手动管理任务状态** | 破坏 SQLite 持久化 |

---

## 工作流执行注意事项

### 关键阶段确认点

**在以下阶段完成后，必须向用户展示结果并确认：**

| 阶段 | 确认内容 | 原因 |
|:---|:---|:---|
| **Discovery** 后 | 需求理解、分析范围 | 确保理解一致 |
| **Exploration** 后 | 探索发现、关键问题 | 让用户了解初步发现 |
| **Architecture** 后 | 改进方案、架构设计 | **最关键**：方案决定后续方向 |
| **Spec** 后 | 规范文档 | 确认细节正确 |

### 常见错误

| 错误 | 表现 | 后果 | 解决方法 |
|:---|:---|:---|:---|
| **机械执行** | 按顺序执行阶段，无中间沟通 | 产出不符合用户需求 | 每个关键阶段后展示结果 |
| **跳过确认** | Architecture 后直接写 Spec | 方案理解不一致，浪费精力 | **Architecture 后必须确认** |
| **忽略反馈** | 用户提出疑问但继续执行 | 用户体验差 | 暂停并调整方向 |

### 正确执行示例

```
Discovery → 展示需求理解 → 用户确认
    ↓
Exploration → 展示发现 → 用户确认
    ↓
Architecture → 展示方案 → 用户确认 ← **关键点**
    ↓
Spec → 编写规范 → 用户确认
    ↓
实现...
```

### 错误执行示例

```
Discovery → Exploration → Architecture → Spec ❌
              (无中间展示和确认)
    ↓
结果：Spec 文档基于错误理解，完全浪费
```

---

## 记住

> **用户是决策者，AI 是执行者。每个关键决策点都需要用户确认。**

Architecture 阶段是**最重要的确认点**，因为它决定了后续所有工作的方向。
