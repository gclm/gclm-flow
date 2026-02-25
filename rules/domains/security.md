# 安全规则

所有项目通用的安全编码规范。

## 认证与授权

### 认证
- 使用成熟的认证框架（JWT、OAuth2、Session）
- 密码必须使用强加密算法存储（bcrypt、argon2）
- 实现 Token 过期和刷新机制
- 支持 MFA（多因素认证）

### 授权
- 实施最小权限原则
- 使用 RBAC（基于角色的访问控制）
- 验证每个请求的权限
- 不依赖客户端进行授权检查

## 输入验证

### 原则
- **永远不要信任用户输入**
- 在边界层验证所有输入
- 使用白名单验证
- 转义特殊字符

### 验证规则
```python
# 使用验证库
from pydantic import BaseModel, EmailStr, Field

class UserInput(BaseModel):
    email: EmailStr
    age: int = Field(ge=0, le=150)
    name: str = Field(min_length=1, max_length=100)
```

## 常见漏洞防护

### SQL 注入
```python
# ✅ 正确：参数化查询
cursor.execute("SELECT * FROM users WHERE id = %s", (user_id,))

# ❌ 错误：字符串拼接
cursor.execute(f"SELECT * FROM users WHERE id = {user_id}")
```

### XSS（跨站脚本）
```python
# 输出时转义
from markupsafe import escape
safe_output = escape(user_input)
```

### CSRF（跨站请求伪造）
- 使用 CSRF Token
- 验证 Referer 头
- SameSite Cookie 属性

### 路径遍历
```python
# 验证文件路径
import os
safe_path = os.path.join(base_dir, user_path)
if not os.path.abspath(safe_path).startswith(os.path.abspath(base_dir)):
    raise SecurityError("Invalid path")
```

## 敏感数据处理

### 存储
- 敏感数据加密存储
- 使用环境变量存储密钥
- 不在代码中硬编码密钥
- 使用密钥管理服务（KMS）

### 传输
- 强制使用 HTTPS
- 使用 TLS 1.2+
- 验证 SSL 证书

### 日志
- 不记录敏感信息（密码、Token、PII）
- 日志脱敏处理
- 安全存储日志

## API 安全

### 请求限制
- 实现速率限制
- 防止暴力破解
- 使用 CAPTCHA

### 响应
- 不暴露内部错误信息
- 不泄露技术栈信息
- 使用统一的错误响应格式

## 依赖安全

- 定期更新依赖
- 使用依赖扫描工具
- 审查第三方库的安全性
- 锁定依赖版本
