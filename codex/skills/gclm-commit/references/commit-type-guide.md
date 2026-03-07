# Commit Type Guide

用于选择 commit type 和 scope。

## 常用类型

- `feat`: 新功能
- `fix`: Bug 修复
- `refactor`: 重构，不改变外部行为
- `docs`: 文档
- `test`: 测试
- `chore`: 杂项维护

## Scope 选择

优先用：
- 模块名
- 子系统名
- 功能域

避免：
- 过大或过泛的 scope，例如 `project`、`misc`
