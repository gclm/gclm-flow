---
paths:
  - "**/*.rs"
---
# Rust Patterns

> This file extends [common/patterns.md](../common/patterns.md) with Rust specific content.

## 项目结构（Axum 服务）

```
src/
├── main.rs          # 启动入口
├── lib.rs           # 公共 API 导出
├── config.rs        # 配置加载
├── error.rs         # 统一错误类型
├── routes/          # HTTP 路由
├── handlers/        # 请求处理
├── services/        # 业务逻辑
├── models/          # 数据模型
└── db/              # 数据库访问
```

## 常用模式

- **Newtype**：用 newtype 包装基础类型增加类型安全（`struct UserId(Uuid)`）
- **Builder**：复杂结构体用 `derive_builder` 或手写 builder
- **State 注入**：Axum 用 `State<AppState>` 传递依赖，不用全局变量
- **Error 统一**：`error.rs` 定义 `AppError`，实现 `IntoResponse`

## 禁止

- 禁止全局可变状态（`static mut`）
- 禁止 `clone()` 性能热路径（用借用替代）
- 禁止在 `async` 上下文中使用阻塞 IO（用 `tokio::fs`、`spawn_blocking`）

## Reference

See skill: `rust-stack` for Axum implementation patterns.
