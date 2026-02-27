---
name: gclm-core
description: |
  Gclm-Flow 智能开发工作流系统。自动分析用户意图并编排工作流：
  (1) 新项目 → 自动调用 gclm-init
  (2) 修复问题 → fix 流程
  (3) 新功能 → plan → do 流程
  (4) 提交代码 → 调用 gclm-commit
  (5) 问答/调研 → ask 流程
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
---

# Gclm-Flow 智能编排器

智能开发工作流系统，自动检测项目状态和用户意图，编排合适的工作流。

## 智能检测

### 项目状态检测

| 检测条件 | 状态 | 动作 |
|----------|------|------|
| 无 llmdoc/AGENTS.md | 新项目 | 调用 `gclm-init` |
| 无 .git | 非 Git 项目 | 建议初始化 Git |
| 有 llmdoc | 已有项目 | 继续分析意图 |

### 语言栈检测

| 检测文件 | 语言栈 | 调用技能 |
|----------|--------|----------|
| `pom.xml` / `build.gradle` | Java | `java-stack` |
| `requirements.txt` / `pyproject.toml` | Python | `python-stack` |
| `go.mod` | Go | `go-stack` |
| `Cargo.toml` | Rust | `rust-stack` |
| `package.json` | 前端 | `frontend-stack` |

### 意图识别

| 关键词 | 意图 | 工作流 |
|--------|------|--------|
| 修复、fix、报错、错误 | 问题修复 | [fix-workflow](references/fix-workflow.md) |
| 实现、开发、新增、功能 | 新功能开发 | [plan](references/plan-workflow.md) → [do](references/do-workflow.md) |
| 怎么、为什么、如何、查询 | 问答/调研 | [ask-workflow](references/ask-workflow.md) |
| 提交、commit | 代码提交 | `gclm-commit` |
| 审查、review、检查 | 代码审查 | `code-review` |
| 测试、test | 运行测试 | `testing` |

## 工作流

### 新功能开发
```
检测语言栈 → 加载对应技能 → plan → do → test → review
```
详见：[plan-workflow.md](references/plan-workflow.md)、[do-workflow.md](references/do-workflow.md)

### 问题修复
```
分析错误 → 定位问题 → 查询历史 → 修复 → 验证 → 记录
```
详见：[fix-workflow.md](references/fix-workflow.md)

### 问答/调研
```
分析问题 → 收集上下文 → 查询历史 → 生成回答
```
详见：[ask-workflow.md](references/ask-workflow.md)

## 代理协作

| 任务 | 代理 | 职责 |
|------|------|------|
| 规划 | planner | 需求分析、任务分解 |
| 实现 | builder | 代码编写、重构 |
| 审查 | reviewer | 质量检查、安全审计 |
| 调研 | investigator | 技术调研、上下文理解 |
| 记录 | recorder | 文档维护、知识管理 |
| 记忆 | remember | 错误记录、模式提取 |

## 可用技能

| 技能 | 调用方式 | 用途 |
|------|----------|------|
| `gclm-init` | `/gclm-init` | 项目初始化 |
| `gclm-commit` | `/gclm-commit` | 智能提交 |
| `java-stack` | `/java-stack` | Java/Spring Boot |
| `python-stack` | `/python-stack` | Python/FastAPI |
| `go-stack` | `/go-stack` | Go/Gin |
| `rust-stack` | `/rust-stack` | Rust/Axum |
| `frontend-stack` | `/frontend-stack` | React/Vue |
| `code-review` | `/code-review` | 代码审查 |
| `testing` | `/testing` | 测试 |
| `documentation` | `/documentation` | 文档 |
| `memory` | `/memory` | 记忆系统 |
| `database` | `/database` | 数据库 |
| `devops` | `/devops` | DevOps |
