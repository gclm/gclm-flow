---
name: executing-plans
description: Use when you already have a written implementation plan and need to execute it step by step, either in one session or in controlled batches with review checkpoints.
---

# 按计划执行

按计划落实，不在执行中临时发明需求。这个 skill 统一覆盖原来的“批次执行”和“当前会话逐任务推进”两种模式。

## 核心规则

- 没有现成计划时，先回到 `writing-plans`
- 先识别计划缺口，再执行，不盲冲
- 每个任务或批次完成后，都要做验证和 review
- 遇到阻塞立即停下，不猜、不硬推

## 两种执行模式

### 模式 A：当前会话逐任务推进

适合：
- 任务之间相对独立
- 希望每一步都 review
- 想减少上下文污染

做法：
1. 读完整计划并拆成任务。
2. 每次只推进一个任务。
3. 每个任务结束后做两层检查：
   - 是否符合需求/计划
   - 是否存在代码质量问题
4. 修完再进下一个任务。

### 模式 B：分批执行

适合：
- 任务较多
- 想按阶段汇报
- 需要用户在批次边界确认

做法：
1. 读完整计划，指出缺口与风险。
2. 建立任务清单，明确当前批次。
3. 批量执行当前批次。
4. 批次结束后汇报完成项、验证证据、剩余风险。

## 选择规则

- 如果任务高度耦合，优先模式 A
- 如果任务规模较大、用户想阶段性确认，优先模式 B
- 如果你无法明确边界，先小批次推进，不要一次铺满

## 联动技能

- `writing-plans`
- `code-review`
- `verification-before-completion`
- `finishing-a-development-branch`
