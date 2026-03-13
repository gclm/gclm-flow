---
name: security-review
description: Use when adding authentication, handling user input, working with secrets, creating API endpoints, implementing payment or sensitive features, or before any production deployment.
---

# 安全审查

## 核心检查项

### 1. Secrets 管理

- 禁止硬编码 API key、密码、token
- 所有 secrets 通过环境变量注入
- `.env.local` 在 `.gitignore` 中
- 验证 secrets 存在，缺失时 fail fast

### 2. 输入验证

- 所有用户输入用 schema 验证（zod、pydantic、joi 等）
- 文件上传校验：大小、类型、扩展名
- 使用白名单而非黑名单
- 错误信息不暴露内部细节

### 3. SQL 注入防护

- 禁止字符串拼接 SQL
- 始终使用参数化查询或 ORM

### 4. 认证与授权

- Token 存 httpOnly cookie，不存 localStorage
- 每个敏感操作前校验权限
- 实现 RBAC

### 5. XSS 防护

- 用户提供的 HTML 需 sanitize
- 配置 CSP headers

### 6. 敏感数据暴露

- 日志中不记录密码、token、卡号
- 生产环境错误信息对用户通用化，详细错误只写服务端日志
- 不向用户暴露 stack trace

### 7. 依赖安全

- `npm audit` / `pip audit` 定期运行
- lock 文件提交到 git
- CI 中启用依赖扫描

## 上线前检查清单

- [ ] 无硬编码 secrets
- [ ] 所有用户输入已验证
- [ ] SQL 查询已参数化
- [ ] 认证 token 存储安全
- [ ] 权限检查到位
- [ ] 错误信息已脱敏
- [ ] 日志无敏感数据
- [ ] 依赖无已知漏洞
- [ ] HTTPS 强制开启

## 联动技能

- `code-review`
- `verification-before-completion`
- `testing`
