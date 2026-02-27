---
name: testing
description: |
  测试技能。当用户要求测试、test、unit test、pytest、jest、vitest、go test 时自动触发。
  包含：(1) 单元测试 (2) 集成测试 (3) E2E 测试 (4) 测试覆盖率
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - test
    - tdd
    - coverage
---

# 测试

## 测试原则

### FIRST 原则
- **F**ast：测试要快
- **I**ndependent：测试要独立
- **R**epeatable：测试可重复
- **S**elf-validating：自动验证
- **T**imely：及时编写

### TDD 工作流

```
RED     → 先写失败的测试
GREEN   → 写最少代码通过测试
REFACTOR → 重构代码
REPEAT  → 继续下一个需求
```

## 测试类型

| 类型 | 比例 | 用途 |
|------|------|------|
| 单元测试 | 70% | 测试单个函数/方法 |
| 集成测试 | 20% | 测试模块间交互 |
| E2E 测试 | 10% | 测试关键流程 |

## 覆盖率目标

| 代码类型 | 目标 |
|---------|------|
| 核心业务逻辑 | 90%+ |
| 公共 API | 80%+ |
| 一般代码 | 70%+ |

## 语言特定

| 语言 | 框架 | 参考 |
|------|------|------|
| Java | JUnit 5 + Mockito | `java-stack/references/testing.md` |
| Python | pytest | `python-stack/references/testing.md` |
| Go | testing (内置) | `go-stack/references/testing.md` |
| Rust | cargo test | `rust-stack/references/testing.md` |
| 前端 | Vitest | `frontend-stack/references/testing.md` |
