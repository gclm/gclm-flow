---
name: code-review
description: Use when reviewing code changes, pull requests, regressions, security concerns, implementation quality, audit feedback, or when preparing a focused review of recent changes.
---

# 代码审查

像资深工程师一样做结构化审查，先给出问题，再给总结。审查、发起复审、处理审查反馈统一走这个 skill。

## 核心规则

- 默认从 `git diff`、PR 变更、最近修改文件开始，不做无边界全文审查。
- 输出以 `findings` 为主，按严重级别排序；没有问题时也要明确说明“未发现问题”。
- 关注行为回归、边界条件、安全、性能、可维护性、测试缺口，不以风格意见淹没关键问题。
- 发现高风险结论时，必须指出触发条件、影响范围、复现线索和最小修复方向。
- 收到 reviewer 反馈时，先验证反馈是否成立，再决定修复、澄清或反驳。

## 使用场景

### 1. 审查代码变更

1. 明确变更范围：`git status`、`git diff --stat`、关键文件 diff。
2. 理解意图：需求、计划、设计文档、现有约定、相关 tests。
3. 逐文件审查：优先看新增逻辑、危险路径、删除/重构区域。
4. 交叉验证：检查实现是否和测试、文档、配置、迁移脚本一致。
5. 输出报告：先列 findings，再补充风险、假设、测试缺口和变更概览。

当 diff 里出现删除旧代码、移除配置、废弃接口、迁移脚本时，额外执行 `removal-plan` 检查。

### 2. 发起复审

在任务完成、合并前、或高风险重构后，先准备这几项：
- review 范围：`git diff --stat`、关键文件、测试变更
- 改动目标和非目标
- 希望 reviewer 重点关注的风险面：安全、回归、并发、性能、兼容性
- 对比基线：从哪个 commit/分支改到现在

### 3. 处理审查反馈

1. 先完整读完反馈，不表演式认同。
2. 回到代码、需求、测试里确认反馈是否成立。
3. 成立就修复；不明确就提澄清问题；不成立就给出技术理由。
4. 多条反馈先按阻塞级别排序，逐条处理，不要半懂半改。

## 重点检查项

- 正确性：逻辑错误、状态遗漏、空值/并发/顺序问题、兼容性回归
- 安全性：认证授权、注入、XSS/SSRF/路径遍历、敏感信息泄露、权限边界
- 性能：N+1、重复 IO、无界循环、缓存缺失、锁竞争、过大对象复制
- 可维护性：命名、职责分离、重复逻辑、隐式耦合、错误处理、可观测性
- 测试：是否覆盖新行为、失败路径、回归场景；测试是否真验证了需求
- 删除/弃用：大规模删除、废弃开关、配置迁移、数据迁移是否有依赖核查和回滚路径
- 反馈处理：评论是否成立、是否与真实代码/需求一致、是否存在更小修复面

详细清单见：
- [severity-guide.md](references/severity-guide.md)
- [security-checklist.md](references/security-checklist.md)
- [performance-checklist.md](references/performance-checklist.md)
- [removal-plan.md](references/removal-plan.md)
- [handling-review-feedback.md](references/handling-review-feedback.md)
- [review-output-template.md](references/review-output-template.md)

## 输出格式

### 有问题时

```markdown
1. [严重级别] [文件:行号] 问题结论
   影响：会导致什么实际后果
   依据：为什么这是问题
   建议：最小修复方向

## Open Questions
- 需要确认的前提或设计意图

## Residual Risks
- 未覆盖但值得关注的风险

## Change Summary
- 简述本次变更做了什么
```

### 无问题时

```markdown
未发现需要阻塞的代码问题。

Residual Risks
- 仍未验证的边界或测试空白

Change Summary
- 简述审查范围与主要改动
```

## 严重级别

- `P0`：会导致数据破坏、权限绕过、严重安全事故、生产不可用
- `P1`：明显功能错误、稳定性问题、较高概率回归
- `P2`：可维护性/性能/测试缺口，建议在合并前处理
- `P3`：低风险建议，非阻塞

## 审查边界

- 不把“我更喜欢另一种写法”当作问题。
- 不在证据不足时下结论；需要时明确写出假设。
- 用户要求 review 时，主回答必须先给 findings，而不是先写摘要。

## 联动技能

- `testing`
- `verification-before-completion`
- `systematic-debugging`
