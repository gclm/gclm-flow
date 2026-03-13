# Commit Type Guide

用于选择 commit type 和 scope。

## 何时查看

- 已经明确要提交，但 type 或 scope 还不稳定
- 需要快速判断这次改动属于 feat、fix、refactor 还是 docs

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

## 注意事项

- commit type 服务于后续检索和回滚，不要为了“显得正式”而选错类型
