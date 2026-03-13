---
paths:
  - "**/*.java"
---
# Java Security

> This file extends [common/security.md](../common/security.md) with Java specific content.

## SQL 注入

- 禁止拼接 SQL 字符串，使用 JPA/JPQL 参数绑定或 `PreparedStatement`
- MyBatis 使用 `#{}` 而非 `${}`

## 认证与授权

- Spring Security 统一管理认证，不手动实现
- 方法级权限用 `@PreAuthorize`
- JWT 存 HttpOnly Cookie，不存 localStorage

## 敏感数据

- 密码用 `BCryptPasswordEncoder`，禁止 MD5/SHA1
- 日志中屏蔽密码、token、手机号等字段
- 配置文件中的 secrets 通过环境变量或 Vault 注入

## Reference

See skill: `security-review` for complete security checklist.
