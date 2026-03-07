---
name: writing-skills
description: Use when creating, editing, or refactoring agent skills so they are discoverable, concise, testable, and reusable.
---

# 编写 Skills

写 skill 的目标不是“解释自己会做什么”，而是让模型在正确时机稳定触发并按预期执行。

## 设计原则

- frontmatter 只保留 `name` 和 `description`
- `description` 只写何时触发，不写具体流程
- `SKILL.md` 保持精简，长清单放 `references/`
- 优先写触发词、边界、步骤、输出要求、常见误用

## 维护方式

1. 先识别当前 skill 的触发问题或流程漏洞。
2. 最小化修改 `description` 和结构。
3. 把长案例、清单、模板拆到 `references/`。
4. 修改后做一次本地检查：目录、frontmatter、引用路径、同步结果。

## 常见问题

- 描述过长，模型只看 description 不读正文
- skill 混入项目特例，复用性差
- 没写何时不用，导致误触发
