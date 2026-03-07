---
name: testing
description: Use when designing, writing, running, or reviewing automated tests, or when deciding what level of tests is needed for a code or config change.
---

# 测试

这个 skill 负责测试策略、测试分层、测试执行和覆盖判断。它不替代 `test-driven-development`；TDD 负责“先写失败测试再写实现”，这里负责“该测什么、怎么测、测到什么程度”。

## 核心规则

- 先选最小但足够证明行为的测试层级，不默认上 E2E
- 新行为、bug 修复、危险重构，至少补一条能防回归的自动化验证
- 测试要验证行为，不要只验证实现细节
- 宣称“测试通过”前，必须跑对应命令并看结果

## 选择哪种测试

### 单元测试

优先用于：
- 纯函数、规则判断、边界条件
- 错误处理与分支覆盖
- 快速回归验证

### 集成测试

优先用于：
- 模块交互、数据库访问、外部接口适配层
- 配置装配、协议编解码、文件系统交互

### 端到端测试

只用于：
- 关键用户路径
- 高价值主流程
- 只有全链路才能暴露的问题

## 测试顺序

1. 先定义要证明的行为和失败条件。
2. 选择最小测试层级。
3. 写测试数据和断言，优先验证对外行为。
4. 跑相关测试。
5. 如涉及共享模块或危险路径，再扩大验证范围。

## 覆盖判断

- 核心业务逻辑：优先覆盖正常路径、失败路径、边界条件
- bug 修复：必须有一条能复现原问题的回归测试，或有明确理由说明无法自动化
- 配置/脚本改动：至少提供可重复执行的验证步骤
- UI/交互改动：除自动化测试外，必要时补实际预览或截图验证

## 常见误区

- 只有 happy path，没有失败路径
- 断言太弱，只断言“函数被调用”而不是结果正确
- 用大而慢的 E2E 替代本该写的单元/集成测试
- 改了公共能力，却只跑了局部测试
- 把日志输出当成测试通过证据

## 语言入口

- Java：`java-stack/references/testing.md`
- Python：`python-stack/references/testing.md`
- Go：`go-stack/references/testing.md`
- Rust：`rust-stack/references/testing.md`
- Frontend：`frontend-stack/references/testing.md`

## 联动技能

- `test-driven-development`
- `systematic-debugging`
- `verification-before-completion`
- `code-review`
