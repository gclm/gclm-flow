---
name: subagent-driven-development
description: Use when executing a multi-step implementation in the current session and the work can be split into clear task-sized chunks with review checkpoints.
---

# 子 Agent 驱动开发

把实现拆成小任务推进，每个任务完成后马上 review，而不是最后一起爆雷。

## 工作流

1. 读取计划并拆成任务。
2. 每次只推进一个任务。
3. 每个任务结束后做两层检查：
   - 是否符合需求/计划
   - 是否存在代码质量问题
4. 修复后再进入下一个任务。
5. 全部完成后做一次总体回顾。

## 适合

- 任务之间相对独立
- 需要持续 review 闭环
- 希望减少上下文污染
