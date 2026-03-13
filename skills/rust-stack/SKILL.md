---
name: rust-stack
description: Use when working on Rust services, Axum or Actix applications, ownership-heavy refactors, async/concurrency issues, or Rust-specific architecture and runtime concerns.
---

# Rust 技术栈

这个 skill 负责 Rust 项目的入口判断和工程经验索引。主文档只保留适用范围、关键关注点和 references 入口。

## 核心规则

- 先判断问题属于所有权/借用、错误处理、async、框架层还是模块边界
- 通用测试策略看 `testing`；Rust 特有测试细节看 `references/testing.md`
- 详细实践优先写入 `references/`

## 重点关注

- Axum / Actix 路由与状态管理
- 错误类型和边界
- async 任务、资源共享、并发模型
- crate / module 边界与类型设计

## 参考资料

- [axum.md](references/axum.md)
- [testing.md](references/testing.md)
