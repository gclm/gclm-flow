# Security Checklist

审查以下问题是否存在：

- 输入是否校验、规范化、限制长度
- 是否存在 SQL/NoSQL/command/template 注入
- 前端输出是否存在 XSS 或 HTML 注入
- URL 拉取、代理、回调是否有 SSRF 风险
- 文件读写是否有路径遍历、任意上传、覆盖风险
- 鉴权、租户隔离、资源归属校验是否完整
- 日志、错误、配置中是否泄露 token、cookie、密钥、PII
- Webhook、hook、外部命令是否有权限放大问题
