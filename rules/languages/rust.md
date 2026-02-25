# Rust/Axum/Actix 规则

Rust、Axum、Actix 项目的编码规范和最佳实践。

## Rust 编码规范

### 1. 命名规范
```rust
// 类型/结构体/枚举：PascalCase
struct UserService {}
enum UserStatus {}

// 函数/方法/变量：snake_case
fn get_user_by_id(user_id: i64) -> Option<User> {}
let user_count = 0;

// 常量：SCREAMING_SNAKE_CASE
const MAX_RETRY_COUNT: i32 = 3;
const DEFAULT_TIMEOUT_SECS: u64 = 30;

// 模块：snake_case
mod user_service;

// 生命周期：单个小写字母
struct User<'a> {
    name: &'a str,
}

// 泛型类型：PascalCase，通常单个大写字母
struct Container<T> {
    value: T,
}
```

### 2. 项目结构
```
myapp/
├── Cargo.toml
├── src/
│   ├── main.rs           # 二进制入口
│   ├── lib.rs            # 库入口
│   ├── config.rs         # 配置
│   ├── error.rs          # 错误类型
│   ├── handlers/         # HTTP 处理器
│   │   ├── mod.rs
│   │   └── users.rs
│   ├── services/         # 业务逻辑
│   │   ├── mod.rs
│   │   └── user_service.rs
│   ├── models/           # 数据模型
│   │   ├── mod.rs
│   │   └── user.rs
│   ├── repositories/     # 数据访问
│   │   ├── mod.rs
│   │   └── user_repo.rs
│   └── utils/            # 工具函数
│       └── mod.rs
└── tests/                # 集成测试
    └── integration_test.rs
```

### 3. 错误处理
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

    #[error("Internal server error")]
    InternalError,
}

// 实现 IntoResponse for Axum
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            AppError::UserNotFound(_) => (StatusCode::NOT_FOUND, self.to_string()),
            AppError::InvalidInput(_) => (StatusCode::BAD_REQUEST, self.to_string()),
            AppError::DatabaseError(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Database error".to_string()),
            AppError::InternalError => (StatusCode::INTERNAL_SERVER_ERROR, "Internal error".to_string()),
        };

        (status, Json(ApiResponse::<()>::error(&message))).into_response()
    }
}

// 结果类型别名
pub type AppResult<T> = Result<T, AppError>;
```

### 4. Option 和 Result 处理
```rust
// 使用 ? 运算符
fn get_user_email(id: i64) -> AppResult<String> {
    let user = repository::find_by_id(id)?
        .ok_or(AppError::UserNotFound(id))?;
    Ok(user.email)
}

// 使用 map/and_then
fn get_user_name(id: i64) -> Option<String> {
    repository::find_by_id(id)
        .map(|user| user.name)
}

// 模式匹配
fn process_user(user: Option<User>) -> String {
    match user {
        Some(u) => u.name,
        None => "Unknown".to_string(),
    }
}

// if let
fn process_if_exists(user: Option<User>) {
    if let Some(u) = user {
        println!("User: {}", u.name);
    }
}
```

## Axum 最佳实践

### 1. 路由组织
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

### 2. 处理器
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
    Ok((StatusCode::CREATED, Json(ApiResponse::success(UserResponse::from(user)))))
}
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

### 4. 中间件
```rust
use axum::{
    middleware::{self, Next},
    response::Response,
};

// 认证中间件
async fn auth_middleware(
    State(state): State<AppState>,
    mut req: Request,
    next: Next,
) -> AppResult<Response> {
    let auth_header = req
        .headers()
        .get("Authorization")
        .and_then(|h| h.to_str().ok())
        .ok_or(AppError::Unauthorized)?;

    let claims = validate_token(auth_header, &state.jwt_secret)?;
    req.extensions_mut().insert(claims);

    Ok(next.run(req).await)
}

// 日志中间件
async fn logging_middleware(req: Request, next: Next) -> Response {
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

### 5. 状态管理
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

## Actix 最佳实践

### 1. 路由和处理器
```rust
use actix_web::{web, App, HttpServer, HttpResponse, Responder};

async fn get_user(path: web::Path<i64>, state: web::Data<AppState>) -> impl Responder {
    match state.user_service.get_by_id(*path).await {
        Ok(user) => HttpResponse::Ok().json(ApiResponse::success(user)),
        Err(e) => HttpResponse::NotFound().json(ApiResponse::<()>::error(&e.to_string())),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let state = web::Data::new(AppState::new().await);

    HttpServer::new(move || {
        App::new()
            .app_data(state.clone())
            .service(
                web::scope("/api/v1/users")
                    .route("", web::get().to(list_users))
                    .route("", web::post().to(create_user))
                    .route("/{id}", web::get().to(get_user))
            )
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
```

## 测试规范

### 1. 单元测试
```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_user_creation() {
        let user = User {
            id: 1,
            email: "test@example.com".to_string(),
            name: "Test User".to_string(),
        };

        assert_eq!(user.email, "test@example.com");
        assert_eq!(user.name, "Test User");
    }

    #[tokio::test]
    async fn test_get_user() {
        let state = create_test_state().await;
        let service = &state.user_service;

        let result = service.get_by_id(1).await;

        assert!(result.is_ok());
        let user = result.unwrap();
        assert_eq!(user.id, 1);
    }
}
```

### 2. 集成测试
```rust
#[tokio::test]
async fn test_create_user_api() {
    let app = create_test_app().await;

    let response = app
        .oneshot(
            Request::builder()
                .method("POST")
                .uri("/api/v1/users")
                .header("Content-Type", "application/json")
                .body(Body::from(r#"{"email":"test@example.com","name":"Test"}"#))
                .unwrap(),
        )
        .await
        .unwrap();

    assert_eq!(response.status(), StatusCode::CREATED);
}
```

## 性能优化

### 1. 异步和并发
```rust
use tokio::try_join;

async fn get_user_with_orders(id: i64) -> AppResult<(User, Vec<Order>)> {
    let user_fut = user_repo.find_by_id(id);
    let orders_fut = order_repo.find_by_user_id(id);

    let (user, orders) = try_join!(user_fut, orders_fut)?;
    let user = user.ok_or(AppError::UserNotFound(id))?;

    Ok((user, orders))
}
```

### 2. 连接池
```rust
use sqlx::postgres::PgPoolOptions;

async fn create_pool(database_url: &str) -> Result<PgPool, sqlx::Error> {
    PgPoolOptions::new()
        .max_connections(20)
        .min_connections(5)
        .acquire_timeout(Duration::from_secs(3))
        .connect(database_url)
        .await
}
```
