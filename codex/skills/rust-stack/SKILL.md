---
name: rust-stack
description: |
  Rust/Axum/Actix 技术栈完整开发指南。当检测到 Rust 项目（Cargo.toml）
  或用户明确要求 Rust/Axum/Actix 开发时自动触发。包含：
  (1) 项目结构规范 (2) Axum 最佳实践 (3) 测试模式 (4) 错误处理
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - rust
    - axum
    - actix
---

# Rust 技术栈开发指南

## 框架检测

- 存在 `axum` 依赖 → Axum，详见 [axum.md](references/axum.md)
- 存在 `actix-web` 依赖 → Actix
- 测试相关 → 详见 [testing.md](references/testing.md)

## 标准项目结构

```
myapp/
├── Cargo.toml
├── src/
│   ├── main.rs           # 二进制入口
│   ├── lib.rs            # 库入口
│   ├── handlers/         # HTTP 处理器
│   ├── services/         # 业务逻辑
│   ├── models/           # 数据模型
│   ├── repositories/     # 数据访问
│   ├── error.rs          # 错误类型
│   └── config.rs         # 配置
└── tests/                # 集成测试
```

## 核心规范

### 错误处理

```rust
use thiserror::Error;

#[derive(Error, Debug)]
pub enum AppError {
    #[error("User not found: {0}")]
    UserNotFound(i64),

    #[error("Invalid input: {0}")]
    InvalidInput(String),

    #[error("Database error: {0}")]
    DatabaseError(#[from] sqlx::Error),
}

pub type AppResult<T> = Result<T, AppError>;
```

### 统一响应格式

```rust
use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct ApiResponse<T> {
    pub success: bool,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data: Option<T>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub error: Option<String>,
}

impl<T> ApiResponse<T> {
    pub fn success(data: T) -> Self {
        Self { success: true, data: Some(data), error: None }
    }

    pub fn error(message: &str) -> Self {
        Self { success: false, data: None, error: Some(message.to_string()) }
    }
}
```

### Axum 处理器

```rust
use axum::{
    extract::{Path, State, Json},
    http::StatusCode,
};

pub async fn get_user(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> AppResult<Json<ApiResponse<UserResponse>>> {
    let user = state.user_service.get_by_id(id).await?;
    Ok(Json(ApiResponse::success(UserResponse::from(user))))
}
```

### 状态管理

```rust
use std::sync::Arc;

#[derive(Clone)]
pub struct AppState {
    pub db: PgPool,
    pub user_service: Arc<UserService>,
}

impl AppState {
    pub fn new(db: PgPool) -> Self {
        let user_repo = Arc::new(UserRepository::new(db.clone()));
        let user_service = Arc::new(UserService::new(user_repo));
        Self { db, user_service }
    }
}
```

## 测试规范

```rust
#[tokio::test]
async fn test_get_user() {
    let state = create_test_state().await;
    let service = &state.user_service;

    let result = service.get_by_id(1).await;

    assert!(result.is_ok());
}
```

## 相关技能

- `code-review` - Rust 代码审查
- `testing` - Rust 测试模式
- `database` - SQLx 模式
