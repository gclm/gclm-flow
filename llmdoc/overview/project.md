# 项目介绍

## 项目目标

**gclm-flow** 的目标是为 Claude Code 提供一个智能分流工作流系统，自动判断任务类型并选择最优开发流程。

### 核心价值

1. **自动化决策**: 减少用户选择工作流的认知负担
2. **最佳实践**: 内置 SpecDD、TDD、Document-First 等开发方法论
3. **并行效率**: 多 Agent 并行执行，提升开发效率
4. **文档优先**: llmdoc 策略确保代码上下文理解

---

## 项目范围

### 包含功能

- **智能分流**: 自动检测任务类型 (DOCUMENT / CODE_SIMPLE / CODE_COMPLEX)
- **工作流引擎**: 9 阶段工作流，支持跳过和定制
- **Agent 管理**: 6 个自定义 Agent + 2 个官方插件 Agent
- **代码搜索**: auggie -> llmdoc -> Grep 分层回退策略
- **状态管理**: 工作流状态持久化和恢复
- **Hooks**: 通知和停止钩子

### 主要组件

| 组件 | 描述 |
|:---|:---|
| `agents/` | 自定义 Agent 定义文件 |
| `commands/` | Claude Code 命令定义 |
| `skills/` | 核心工作流 Skill 逻辑 |
| `rules/` | 工作流规则文档 |
| `hooks/` | 生命周期 Hooks |

---

## 非目标

以下内容**不在项目范围内**：

- Claude Code 核心功能的修改
- 第三方 Agent 的实现（仅使用官方插件）
- 独立的测试框架（依赖项目自带的测试工具）
- 持续集成/部署配置（CI/CD）

---

## 技术栈

| 类别 | 技术 |
|:---|:---|
| **脚本语言** | Bash/Zsh |
| **配置格式** | Markdown, YAML, JSON |
| **插件系统** | Claude Code Plugin API |
| **通信协议** | MCP (Model Context Protocol) |
| **AI 模型** | Claude Opus 4.5, Sonnet 4.5, Haiku 4.5 |
| **外部工具** | auggie (可选，推荐) |
