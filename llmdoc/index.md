# gclm-flow 文档索引

## 项目概览

**gclm-flow** 是一个基于 Go 引擎的智能工作流系统，支持自定义工作流 YAML 配置和多 Agent 并行执行。

核心特性：
- **Go 引擎**: 工作流编排和状态管理 (SQLite)
- **YAML 工作流**: 可配置的工作流定义
- **智能分流**: 自动判断任务类型 (DOCUMENT / CODE_SIMPLE / CODE_COMPLEX)
- **多 Agent 并行**: 6 个自定义 Agent + 2 个官方插件 Agent
- **代码搜索分层回退**: auggie (语义搜索) → llmdoc (结构化) → Grep (模式匹配)

---

## 快速导航

### 项目概览
- [项目介绍](overview/project.md) - 项目目标、范围
- [技术栈](overview/tech-stack.md) - 技术栈清单
- [目录结构](overview/structure.md) - 文件组织说明

### 架构设计
- [系统架构](architecture/system.md) - Go 引擎 + 工作流 + Agents
- [工作流配置](architecture/workflows.md) - YAML 工作流定义
- [Agent 体系](architecture/agents.md) - 自定义 Agent 和官方插件
- [代码搜索](architecture/code-search.md) - 分层回退搜索策略
- [数据库设计](architecture/database.md) - SQLite 数据库结构

### 使用指南
- [安装指南](guides/installation.md) - 安装和配置步骤
- [快速开始](guides/quickstart.md) - 基本使用方法
- [工作流开发](guides/workflow-development.md) - 自定义工作流开发

### 参考文档
- [命令参考](reference/commands.md) - gclm-engine 命令列表
- [工作流参考](reference/workflows.md) - 内置工作流说明
- [配置参考](reference/configuration.md) - 配置选项说明

---

## 三种工作流类型

| 类型 | workflow_type | 适用场景 |
|:---|:---|:---|
| 📝 **DOCUMENT** | `DOCUMENT` | 文档编写、架构设计、需求分析 |
| 🔧 **CODE_SIMPLE** | `CODE_SIMPLE` | Bug 修复、小修改、单文件变更 |
| 🚀 **CODE_COMPLEX** | `CODE_COMPLEX` | 新功能、模块开发、跨文件变更 |

---

## Agent 体系

| Agent | 职责 | 模型 | 典型阶段 |
|:---|:---|:---|:---|
| `investigator` | 代码库调查、分析 | Haiku | 需求发现、探索、总结 |
| `architect` | 架构设计、方案权衡 | Opus | 架构设计 |
| `spec-guide` | SpecDD 规范文档编写 | Opus | 规范文档 |
| `tdd-guide` | TDD 流程指导 | Sonnet | 测试编写 |
| `worker` | 执行明确定义的任务 | Sonnet | 代码实现 |
| `code-reviewer` | 代码审查 | Sonnet | 代码审查 |
| `recorder` | 文档记录 | Sonnet | 文档更新 |

**官方插件 Agents**:
- `code-simplifier@claude-plugins-official` - 代码简化重构
- `security-guidance@claude-plugins-official` - 安全审查

---

## 核心组件

| 组件 | 位置 | 用途 |
|:---|:---|:---|
| **Go 引擎** | `gclm-engine/` | 工作流编排、状态管理 |
| **工作流定义** | `workflows/` | YAML 配置的工作流 |
| **Skills** | `skills/gclm/` | Claude Code 集成入口 |
| **Agents** | `agents/` | Agent 定义 |
| **Rules** | `rules/` | 工作流规则 |

---

## 快速开始

```bash
# 安装
./install.sh

# 使用
gclm-engine workflow list
gclm-engine workflow start "修复登录页面 bug"
```

---

## 数据流

```
用户请求 → gclm-engine (Go 引擎) → 工作流编排 → Agent 执行
    ↓              ↓                    ↓
 自然语言    SQLite 状态管理      多 Agent 并行
```
