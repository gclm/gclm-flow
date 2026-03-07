---
name: python-stack
description: |
  Python/FastAPI/Flask 技术栈完整开发指南。当检测到 Python 项目（requirements.txt、pyproject.toml）
  或用户明确要求 Python/FastAPI/Flask 开发时自动触发。包含：
  (1) 项目结构规范 (2) FastAPI 最佳实践 (3) 测试模式 (4) 安全规范
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - python
    - fastapi
    - flask
---

# Python 技术栈开发指南

## 框架检测

- 存在 `fastapi` 依赖 → FastAPI，详见 [fastapi.md](references/fastapi.md)
- 存在 `flask` 依赖 → Flask
- 测试相关 → 详见 [testing.md](references/testing.md)

## 标准项目结构

```
app/
├── main.py              # 应用入口
├── config.py            # 配置管理
├── dependencies.py      # 依赖注入
├── routers/             # 路由模块
│   ├── users.py
│   └── auth.py
├── services/            # 业务逻辑
├── models/              # 数据模型
│   ├── user.py
│   └── schemas.py
├── repositories/        # 数据访问
├── utils/               # 工具函数
└── exceptions.py        # 异常定义
```

## 核心规范

### 依赖注入 (FastAPI)

```python
from fastapi import Depends

# 数据库依赖
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# 认证依赖
async def get_current_user(
    token: str = Depends(oauth2_scheme),
    db: Session = Depends(get_db)
) -> User:
    payload = decode_token(token)
    user = db.query(User).filter(User.id == payload["sub"]).first()
    if not user:
        raise HTTPException(status_code=401, detail="Invalid token")
    return user
```

### Pydantic 模型

```python
from pydantic import BaseModel, EmailStr, Field

class UserCreate(BaseModel):
    email: EmailStr
    name: str = Field(..., min_length=2, max_length=100)
    password: str = Field(..., min_length=8)

class UserResponse(BaseModel):
    id: int
    email: str
    name: str

    class Config:
        from_attributes = True
```

### 统一响应格式

```python
from typing import Generic, TypeVar, Optional
from pydantic import BaseModel

T = TypeVar("T")

class ApiResponse(BaseModel, Generic[T]):
    success: bool = True
    data: Optional[T] = None
    error: Optional[str] = None

    @classmethod
    def ok(cls, data: T) -> "ApiResponse[T]":
        return cls(success=True, data=data)

    @classmethod
    def error(cls, error: str) -> "ApiResponse[T]":
        return cls(success=False, error=error)
```

### 异常处理

```python
from fastapi import Request
from fastapi.responses import JSONResponse

class AppException(Exception):
    def __init__(self, status_code: int, detail: str):
        self.status_code = status_code
        self.detail = detail

@app.exception_handler(AppException)
async def app_exception_handler(request: Request, exc: AppException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"success": False, "error": exc.detail}
    )
```

### 配置管理

```python
from pydantic_settings import BaseSettings
from functools import lru_cache

class Settings(BaseSettings):
    app_name: str = "My API"
    database_url: str
    secret_key: str

    class Config:
        env_file = ".env"

@lru_cache()
def get_settings() -> Settings:
    return Settings()
```

## 测试规范

详见 [testing.md](references/testing.md)

- 使用 pytest
- 目标覆盖率：80%+

## 相关技能

- `code-review` - Python 代码审查
- `testing` - pytest 测试模式
- `database` - SQLAlchemy 模式
