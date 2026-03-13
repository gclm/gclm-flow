---
paths:
  - "**/*.rs"
  - "**/Cargo.toml"
---
# Rust Testing

> This file extends [common/testing.md](../common/testing.md) with Rust specific content.

## Framework

- 单元测试：内置 `#[test]`，放在同文件 `#[cfg(test)]` 模块
- 集成测试：`tests/` 目录
- 异步测试：`#[tokio::test]`

## 运行

```bash
# 所有测试
cargo test

# 带输出
cargo test -- --nocapture

# 特定测试
cargo test test_name
```

## 覆盖率

```bash
cargo llvm-cov --html
```

## 原则

- 单元测试与被测代码同文件，集成测试在 `tests/`
- 使用 `assert_eq!`、`assert_matches!` 精确断言
- 错误路径必须测试（`Result::Err` 分支）
- 异步代码用 `tokio::test` 或 `async-std::test`

## Reference

See skill: `rust-stack` for Axum and async patterns.
