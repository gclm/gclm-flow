---
name: writing-skills
description: Use when creating, editing, or refactoring agent skills so they are discoverable, concise, reusable, and cheap to load.
---

# 编写 Skills

写 skill 的目标不是“解释自己会做什么”，而是让模型在正确时机稳定触发，并在尽量少的上下文里执行正确流程。

## 核心原则

- frontmatter 只保留 `name` 和 `description`
- `description` 只写何时触发，不写流程摘要
- 一个 skill 只解决一个高频问题，不把所有场景揉成大杂烩
- `SKILL.md` 保持精简，把长清单、模板、案例拆到 `references/`
- 能脚本化的重复动作放 `scripts/`，不要把长命令全堆在正文里
- 为搜索优化：写清触发词、症状、同义词、边界条件
- 控制 token 成本：默认先读主文档，按需展开 `references/`

## 目录设计

详细结构见 [packaging-guide.md](references/packaging-guide.md)

- `SKILL.md`：触发条件、边界、最短执行路径、输出要求
- `references/`：长清单、案例、对照表、专项指南
- `scripts/`：可复用脚本或模板生成器
- `assets/`：确实需要复用的模板或静态素材

## 工作流

1. 先识别这个 skill 解决的重复问题和失败模式。
2. 写触发描述：让模型知道“什么时候该读它”。
3. 写最短主流程：步骤、边界、输出格式、常见误用。
4. 把超过约 30 行的专项清单/案例拆到 `references/`。
5. 检查是否存在更小的拆分方式，避免 skill 继续膨胀。
6. 修改后验证：frontmatter、目录、引用路径、同步结果、触发是否清晰。

## 反模式

- description 里把完整流程讲完，导致模型只看 description 不读正文
- skill 同时承担多个不相干职责
- 把项目私货、绝对路径、凭证、一次性经验写进通用 skill
- 主文档太长，没有分层引用
- 明明适合脚本化，却让模型每次手写一大段重复内容

## 验证清单

- 名称是否清晰、可搜索
- description 是否只描述触发条件
- 主流程是否 1 分钟内可扫描完成
- `references/` 链接是否真实存在
- 同步到 `~/.codex` 后是否无漂移
