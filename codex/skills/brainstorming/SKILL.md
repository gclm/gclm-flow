---
name: brainstorming
description: Use when defining a new feature, workflow, behavior change, or non-trivial configuration change before implementation. Trigger on 设计方案, 方案比较, 怎么做更合适, architecture choice.
---

# 方案设计

先收敛目标与边界，再进入实现。

## 何时使用

- 新功能、行为变更、工作流编排、较大配置改造
- 需要在 2 个以上方案之间做取舍
- 需求不完整、成功标准不清晰、风险较高

## 工作流

1. 只读探索：代码、文档、约束、现状。
2. 澄清目标：范围、非目标、验收标准、风险。
3. 给出 2-3 个方案，写清 trade-off 和推荐项。
4. 输出设计结论：架构、数据流、错误处理、验证方式。
5. 等用户确认后再进入实现或计划。

## 产出要求

- 结论先行，不写散乱脑暴
- 明确推荐方案和不选其他方案的原因
- 如果要改代码或配置，必须给出受影响文件/模块范围
