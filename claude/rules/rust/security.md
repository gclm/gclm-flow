---
paths:
  - "**/*.rs"
---
# Rust Security

> This file extends [common/security.md](../common/security.md) with Rust specific content.

## 内存安全

- 避免 `unsafe` 块；必须使用时需注释说明理由并 code review
- 不使用 `std::mem::transmute`（类型转换用安全方式）
- 依赖 `cargo audit` 扫描已知漏洞

```bash
cargo audit
```

## 输入验证

- 使用 `validator` 或 `garde` crate 验证用户输入
- 反序列化时用强类型（`serde` + 自定义 Deserialize），拒绝未知字段

## Secrets

- 使用 `secrecy` crate 包装敏感值，防止意外打印
- 配置从环境变量读取（`dotenvy` 或 `config` crate）

## SQL

- 使用 `sqlx` 的参数绑定，禁止字符串拼接 SQL
- `sqlx::query!` 宏在编译时验证 SQL

## Reference

See skill: `security-review` for complete security checklist.
