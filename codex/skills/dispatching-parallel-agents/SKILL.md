---
name: dispatching-parallel-agents
description: Use when there are 2 or more independent investigations or implementation tasks that can run in parallel without conflicting files or shared state.
---

# 并行 Agent 调度

仅在任务彼此独立时并行，避免“并行制造冲突”。

## 使用条件

- 至少 2 个子任务彼此独立
- 每个子任务都有独立交付物
- 不会频繁编辑同一文件或依赖同一中间结果

## 操作步骤

1. 先拆出子任务清单。
2. 为每个子任务定义范围、输入、输出和禁止事项。
3. 并行执行后统一汇总：改动文件、验证结果、冲突情况。
4. 若发现共享状态或顺序依赖，立即改回串行。

## 不要这样用

- 同一个模块的连续重构
- 先后强依赖的任务
- 你还没搞清楚问题边界时
