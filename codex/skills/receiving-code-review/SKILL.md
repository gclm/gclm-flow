---
name: receiving-code-review
description: Use when receiving review feedback before applying suggestions, especially when comments may be ambiguous, incomplete, or technically questionable.
---

# 处理审查反馈

评审意见先验证，再实现。不要表演式认同，也不要盲改。

## 规则

1. 先完整读完反馈。
2. 逐条确认自己是否真正理解。
3. 回到代码和需求里验证反馈是否成立。
4. 成立就修；不成立就给出技术理由或提出澄清问题。

## 红线

- 不要在没验证前直接说“完全正确，我现在就改”
- 不要只改你理解的一半，剩下模糊项留着不问
- 外部 reviewer 的意见不是命令，要以代码库真实约束为准

## 输出方式

- 正确反馈：`已修复 + 位置 + 简要说明`
- 不明确反馈：`我理解了 X，Y 还需要确认`
- 需要反驳：`当前实现之所以这样，是因为 ...；如果目标改为 ...，可以按另一方案处理`
