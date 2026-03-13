---
paths:
  - "**/*.rs"
  - "**/Cargo.toml"
  - "**/Cargo.lock"
---
# Rust Coding Style

> This file extends [common/coding-style.md](../common/coding-style.md) with Rust specific content.

## 所有权与不可变性

- 默认不可变（`let`），只在必要时用 `let mut`
- 优先借用（`&T`），避免不必要的 clone
- 使用 `Arc<Mutex<T>>` 管理共享状态，不用裸指针

## 错误处理

- 使用 `Result<T, E>` 和 `?` 操作符传播错误
- 自定义错误类型用 `thiserror`，不用 `Box<dyn Error>`
- 禁止 `unwrap()`/`expect()` 在生产代码中（测试代码可以）
- 使用 `anyhow` 处理应用层错误，`thiserror` 处理库层错误

## 文件组织

- 模块按功能域划分，`mod.rs` 只做 re-export
- 单文件不超过 500 行
- 公共 API 用 `pub use` 在 `lib.rs` 统一导出

## 代码质量

```bash
cargo clippy -- -D warnings
cargo fmt --check
```

完成前检查：
- [ ] `cargo clippy` 无警告
- [ ] `cargo fmt` 已格式化
- [ ] 无 `unwrap()`/`expect()` 在非测试代码中
- [ ] 无 `#[allow(dead_code)]`（清理未使用代码）
- [ ] 所有 `pub` API 有文档注释
