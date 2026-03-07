---
name: gclm
description: Use when the user wants a high-level route to the right skill for planning, implementation, review, documentation, testing, or commit work, but the exact entry point is still unclear.
---

# Gclm 路由

这个 skill 不再承担“大总管实现流程”，只负责把任务快速路由到更具体的 skill。

## 路由规则

- 新项目或缺基础文档：`documentation`
- 多步任务、需求不清：`brainstorming` 或 `writing-plans`
- 按计划推进：`executing-plans`
- 代码审查、复审、反馈处理：`code-review`
- 测试策略与测试执行：`testing`
- 提交：`gclm-commit`
- 技术栈专项：`frontend-stack` / `python-stack` / `go-stack` / `java-stack` / `rust-stack`
- 部署与交付：`devops`
- 数据库：`database`

## 使用方式

1. 先识别任务最核心的问题域。
2. 如果已经能明确入口，就直接用对应 skill，不要停留在 `gclm`。
3. 如果任务跨多个域，先选主 skill，再按需联动其他 skill。

## 何时不用

- 已经明确知道要用哪个具体 skill
- 已经在某个领域 skill 的执行上下文里
