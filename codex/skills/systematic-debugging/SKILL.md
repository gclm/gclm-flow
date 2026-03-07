---
name: systematic-debugging
description: Use when facing a bug, failing test, flaky behavior, runtime error, or unexpected output before proposing a fix.
---

# 系统化调试

先找根因，再修。拍脑袋试错最多三轮。

## 四步法

1. 复现：拿到稳定复现方式、日志、输入、环境。
2. 缩小范围：确认是最近变更、依赖、配置还是数据问题。
3. 找根因：用证据解释“为什么会发生”，不要只描述表象。
4. 再修复：优先最小修复，并补验证或回归测试。

## 规则

- 同一问题最多尝试 3 次，每次必须改变假设。
- 没有根因证据前，不给“应该这样改”的结论。
- 如果问题涉及安全、权限、并发、数据破坏，默认按高风险处理。
