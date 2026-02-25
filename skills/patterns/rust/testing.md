# Rust 测试模式

包括单元测试、集成测试、文档测试、基准测试和覆盖率的 Rust 测试模式。

## 何时激活

- 编写新的 Rust 函数或模块
- 为现有代码添加测试覆盖
- 遵循 Rust 项目 TDD 工作流
- 搭建测试基础设施

## TDD 工作流

```rust
// 步骤 1: 编写失败的测试 (RED)
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add() {
        assert_eq!(add(2, 3), 5);
    }
}

// 步骤 2: 实现最少代码 (GREEN)
pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

// 步骤 3: 按需重构
```

## 单元测试

### 基本测试结构

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_basic() {
        assert!(true);
    }

    #[test]
    fn test_equality() {
        assert_eq!(2 + 2, 4);
    }

    #[test]
    fn test_inequality() {
        assert_ne!(2 + 2, 5);
    }
}
```

### 测试 Result

```rust
#[test]
fn test_result() -> Result<(), Box<dyn std::error::Error>> {
    let result = parse_number("42")?;
    assert_eq!(result, 42);
    Ok(())
}
```

### 测试 Panic

```rust
#[test]
#[should_panic(expected = "division by zero")]
fn test_divide_by_zero() {
    divide(10, 0);
}
```

### 忽略测试

```rust
#[test]
#[ignore]
fn expensive_test() {
    // 耗时测试
}

// 运行: cargo test -- --ignored
```

## 断言

```rust
// 基本断言
assert!(condition);
assert_eq!(actual, expected);
assert_ne!(actual, unexpected);

// 带自定义消息
assert_eq!(result, 42, "结果应该是 42");

// 调试断言（仅在调试构建中）
debug_assert!(condition);
debug_assert_eq!(a, b);
```

## 集成测试

### 目录结构

```
project/
├── src/
│   └── lib.rs
└── tests/
    ├── common/
    │   └── mod.rs        # 共享测试工具
    └── integration_test.rs
```

### 集成测试示例

```rust
// tests/integration_test.rs
use myproject::*;

#[test]
fn test_integration() {
    let result = public_function();
    assert!(result.is_ok());
}
```

## 文档测试

```rust
/// 两数相加。
///
/// # Examples
///
/// ```
/// use myproject::add;
///
/// assert_eq!(add(2, 3), 5);
/// ```
pub fn add(a: i32, b: i32) -> i32 {
    a + b
}
```

## 异步测试

```rust
#[tokio::test]
async fn test_async() {
    let result = async_function().await;
    assert_eq!(result, 42);
}

#[tokio::test]
async fn test_async_result() -> Result<(), Box<dyn std::error::Error>> {
    let result = fetch_data().await?;
    assert!(!result.is_empty());
    Ok(())
}
```

## 使用 mockall 进行 Mocking

```rust
use mockall::*;

#[automock]
trait Database {
    fn get_user(&self, id: u64) -> Option<User>;
}

#[test]
fn test_with_mock() {
    let mut mock_db = MockDatabase::new();
    mock_db
        .expect_get_user()
        .with(eq(1))
        .times(1)
        .returning(|_| Some(User { id: 1, name: "Alice".into() }));

    let service = UserService::new(Box::new(mock_db));
    let user = service.get_user(1);

    assert!(user.is_some());
}
```

## 基准测试

```rust
// Cargo.toml
// [dev-dependencies]
// criterion = "0.5"

use criterion::{black_box, criterion_group, criterion_main, Criterion};

fn bench_add(c: &mut Criterion) {
    c.bench_function("add", |b| {
        b.iter(|| add(black_box(2), black_box(3)))
    });
}

criterion_group!(benches, bench_add);
criterion_main!(benches);

// 运行: cargo bench
```

## 测试命令

```bash
# 运行所有测试
cargo test

# 运行特定测试
cargo test test_name

# 运行带输出的测试
cargo test -- --nocapture

# 运行被忽略的测试
cargo test -- --ignored

# 发布模式运行测试
cargo test --release

# 仅运行文档测试
cargo test --doc

# 并行测试执行
cargo test -- --test-threads=4
```

## 覆盖率

```bash
# 使用 cargo-tarpaulin
cargo install cargo-tarpaulin
cargo tarpaulin --out Html

# 使用 cargo-llvm-cov
cargo install cargo-llvm-cov
cargo llvm-cov --html
```

| 代码类型 | 目标 |
|---------|------|
| 核心逻辑 | 100% |
| 公共 API | 90%+ |
| 一般代码 | 80%+ |

## 最佳实践

**应该：**
- 先写测试（TDD）
- 单元测试使用 `#[test]`
- 集成测试使用 `tests/` 目录
- 文档示例使用文档测试
- 测试错误路径

**不应该：**
- 直接测试私有函数
- 忽略失败的测试
- 在测试中使用 `unwrap()`（使用 `?` 或 `expect`）
