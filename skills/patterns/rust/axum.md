# Rust Patterns

Rust/Axum/Actix 开发专用模式和最佳实践。

## 技能描述

这个技能包含 Rust 后端开发的专用模式，适用于 Axum 和 Actix 框架。

## 包含的模式

### 1. 项目结构

```
myapp/
├── Cargo.toml
├── src/
│   ├── main.rs              # 应用入口
│   ├── lib.rs               # 库入口
│   ├── error.rs             # 错误类型
│   ├── config.rs            # 配置
│   ├── handlers/            # HTTP 处理器
│   │   ├── mod.rs
│   │   └── users.rs
│   ├── services/            # 业务逻辑
│   │   ├── mod.rs
│   │   └── user_service.rs
│   ├── models/              # 数据模型
│   │   ├── mod.rs
│   │   └── user.rs
│   ├── repositories/        # 数据访问
│   │   ├── mod.rs
│   │   └── user_repo.rs
│   └── middleware/          # 中间件
│       └── mod.rs
└── tests/
    └── integration_test.rs
```

### 2. 错误处理

```rust
use thiserror::Error;
use axum::{
    http::StatusCode,
    response::{Response, IntoResponse},
    Json,
};

#[derive(Error, Debug)]
pub enum AppError {
    #[error("User not found: {0}")]
    UserNotFound(i64),

    #[error("Invalid input: {0}")]
    InvalidInput(String),

    #[error("Database error: {0}")]
    DatabaseError(#[from] sqlx::Error),

    #[error("Unauthorized")]
    Unauthorized,
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            AppError::UserNotFound(_) => (StatusCode::NOT_FOUND, self.to_string()),
            AppError::InvalidInput(_) => (StatusCode::BAD_REQUEST, self.to_string()),
            AppError::Unauthorized => (StatusCode::UNAUTHORIZED, self.to_string()),
            AppError::DatabaseError(_) => {
                (StatusCode::INTERNAL_SERVER_ERROR, "Database error".to_string())
            }
        };

        (status, Json(ApiResponse::<()>::error(&message))).into_response()
    }
}

pub type AppResult<T> = Result<T, AppError>;
```

### 3. 统一响应格式

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
        Self {
            success: true,
            data: Some(data),
            error: None,
        }
    }

    pub fn error(message: &str) -> Self {
        Self {
            success: false,
            data: None,
            error: Some(message.to_string()),
        }
    }
}
```

### 4. Axum 处理器

```rust
use axum::{
    extract::{Path, State, Json},
    http::StatusCode,
};

#[derive(Debug, serde::Deserialize)]
pub struct CreateUserRequest {
    pub email: String,
    pub name: String,
}

#[derive(Debug, serde::Serialize)]
pub struct UserResponse {
    pub id: i64,
    pub email: String,
    pub name: String,
}

pub async fn get_user(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> AppResult<Json<ApiResponse<UserResponse>>> {
    let user = state.user_service.get_by_id(id).await?;
    Ok(Json(ApiResponse::success(UserResponse::from(user))))
}

pub async fn create_user(
    State(state): State<AppState>,
    Json(req): Json<CreateUserRequest>,
) -> AppResult<(StatusCode, Json<ApiResponse<UserResponse>>)> {
    let user = state.user_service.create(req).await?;
    Ok((
        StatusCode::CREATED,
        Json(ApiResponse::success(UserResponse::from(user)))
    ))
}
```

### 5. 路由组织

```rust
use axum::{
    routing::{get, post, put, delete},
    Router,
};

pub fn create_router(state: AppState) -> Router {
    Router::new()
        .route("/health", get(health_check))
        .nest("/api/v1/users", user_routes())
        .with_state(state)
}

fn user_routes() -> Router<AppState> {
    Router::new()
        .route("/", get(list_users).post(create_user))
        .route("/:id", get(get_user).put(update_user).delete(delete_user))
}
```

### 6. 状态管理

```rust
use std::sync::Arc;

#[derive(Clone)]
pub struct AppState {
    pub db: PgPool,
    pub user_service: Arc<UserService>,
    pub config: Arc<Config>,
}

impl AppState {
    pub fn new(db: PgPool, config: Config) -> Self {
        let user_repo = Arc::new(UserRepository::new(db.clone()));
        let user_service = Arc::new(UserService::new(user_repo));

        Self {
            db,
            user_service,
            config: Arc::new(config),
        }
    }
}
```

### 7. 中间件

```rust
use axum::{
    middleware::{self, Next},
    response::Response,
    body::Body,
    http::Request,
};

// 认证中间件
pub async fn auth_middleware(
    State(state): State<AppState>,
    mut req: Request<Body>,
    next: Next,
) -> AppResult<Response> {
    let auth_header = req
        .headers()
        .get("Authorization")
        .and_then(|h| h.to_str().ok())
        .ok_or(AppError::Unauthorized)?;

    let claims = validate_token(auth_header, &state.config.jwt_secret)?;
    req.extensions_mut().insert(claims);

    Ok(next.run(req).await)
}

// 日志中间件
pub async fn logging_middleware(req: Request<Body>, next: Next) -> Response {
    let start = std::time::Instant::now();
    let method = req.method().clone();
    let path = req.uri().path().to_string();

    let response = next.run(req).await;

    let elapsed = start.elapsed();
    tracing::info!(
        method = %method,
        path = %path,
        status = %response.status().as_u16(),
        elapsed_ms = %elapsed.as_millis(),
        "Request completed"
    );

    response
}
```

### 8. 仓储模式

```rust
use sqlx::PgPool;

pub struct UserRepository {
    db: PgPool,
}

impl UserRepository {
    pub fn new(db: PgPool) -> Self {
        Self { db }
    }

    pub async fn find_by_id(&self, id: i64) -> AppResult<Option<User>> {
        let user = sqlx::query_as::<_, User>(
            "SELECT id, email, name, created_at FROM users WHERE id = $1"
        )
        .bind(id)
        .fetch_optional(&self.db)
        .await?;

        Ok(user)
    }

    pub async fn create(&self, user: &CreateUser) -> AppResult<User> {
        let user = sqlx::query_as::<_, User>(
            "INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, email, name, created_at"
        )
        .bind(&user.email)
        .bind(&user.name)
        .fetch_one(&self.db)
        .await?;

        Ok(user)
    }
}
```

## 使用场景

- 创建新的 Rust 项目
- Axum/Actix 框架开发
- 错误处理设计
- 异步编程
