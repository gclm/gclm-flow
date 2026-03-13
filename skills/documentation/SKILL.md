---
name: documentation
description: Use when initializing project docs, updating README or API docs, restructuring documentation, or checking whether code, config, and user-visible behavior changes need documentation updates.
---

# 文档

这个 skill 统一负责文档初始化、文档维护和文档结构治理。它既覆盖新项目的基础文档骨架，也覆盖后续 README、API、运维和架构文档更新。

## 核心规则

- 先判断是“初始化文档”还是“更新已有文档”
- 文档只更新和当前变更直接相关的部分，不机械大扫除
- 改了接口、配置、命令、用户可见行为后，要检查是否需要同步文档
- 初始化文档时，先给出最小可用骨架，不默认生成过多模板

## 两类场景

### 1. 文档初始化

适用于：
- 新仓库缺少 `README.md`、`AGENTS.md`、`docs/` 或 `llmdoc/` 等基础结构
- 需要建立最小项目说明、开发入口、架构骨架

最小初始化通常包括：
- 项目概览
- 开发/运行入口
- 核心目录说明
- 后续扩展所需的文档目录骨架

详细参考：
- [documentation-bootstrap.md](references/documentation-bootstrap.md)

### 2. 文档更新

适用于：
- API、schema、配置、CLI、hooks、agents、部署流程发生变化
- 需要补 README、运维说明、架构说明或迁移说明

详细参考：
- [documentation-drift-checklist.md](references/documentation-drift-checklist.md)

## 工作顺序

1. 识别变更面：初始化、README、API、配置、运维、架构。
2. 确定最小需要更新的文档集合。
3. 保持文档结构清晰，长清单放 `references/`。
4. 如文档与代码存在漂移，优先修正用户会直接依赖的部分。

## 常见误区

- 新项目一上来生成一大套没人维护的空文档
- 改了接口或配置，却只改代码不改说明
- README 写成营销文案，没有可执行入口
- 把一次性会议记录混进长期文档结构

## 联动技能

- `code-review`
- `testing`
- `verification-before-completion`
