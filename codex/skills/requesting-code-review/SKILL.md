---
name: requesting-code-review
description: Use when finishing a task, before merge, or after a risky refactor to request a focused review against the exact changed scope.
---

# 发起代码复审

让审查针对“刚刚改了什么”，而不是泛泛而谈。

## 何时发起

- 完成一个主要任务后
- 合并前
- 修复杂 bug 或高风险重构后

## 发起前准备

- 明确 review 范围：`git diff --stat`、关键文件、测试变更
- 写出本次改动目标和非目标
- 指定需要 reviewer 重点看的风险面：安全、回归、并发、性能、兼容性

## 建议格式

- 改了什么
- 依据哪份计划/需求
- 对比基线是什么
- 最担心 reviewer 看漏什么
