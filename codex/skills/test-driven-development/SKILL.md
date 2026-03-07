---
name: test-driven-development
description: Use when implementing a feature, bugfix, or behavior change and the project has a practical automated test path.
---

# TDD

先写失败测试，再写最小实现。

## 最小闭环

1. 写一个能表达需求的失败测试。
2. 运行它，确认失败原因正确。
3. 写最少代码让它通过。
4. 再跑测试，确认通过且无回归。
5. 最后再重构。

## 适用判断

优先用于：
- bug 修复
- 新行为
- 回归风险高的改动

不强推用于：
- 纯文档
- 纯配置且没有测试入口
- 生成产物

## 红线

- 不要先写实现再补测试冒充 TDD
- 不要跳过“先看失败”这一步
