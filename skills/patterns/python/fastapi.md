# Python Patterns

Python/Flask/FastAPI 开发专用模式和最佳实践。

## 技能描述

这个技能包含 Python 后端开发的专用模式，适用于 Flask 和 FastAPI 框架。

## 包含的模式

### 1. 项目结构

```
app/
├── main.py              # 应用入口
├── config.py            # 配置管理
├── dependencies.py      # 依赖注入
├── routers/             # 路由模块
│   ├── __init__.py
│   ├── users.py
│   └── auth.py
├── services/            # 业务逻辑
│   ├── __init__.py
│   └── user_service.py
├── models/              # 数据模型
│   ├── __init__.py
│   ├── user.py
│   └── schemas.py
├── repositories/        # 数据访问
│   ├── __init__.py
│   └── user_repo.py
├── utils/               # 工具函数
│   ├── __init__.py
│   └── security.py
└── exceptions.py        # 异常定义
```

### 2. FastAPI 依赖注入

```python
from fastapi import Depends
from sqlalchemy.orm import Session

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

# 使用
@router.get("/me")
async def get_me(current_user: User = Depends(get_current_user)):
    return current_user
```

### 3. Pydantic 模型

```python
from pydantic import BaseModel, EmailStr, Field
from datetime import datetime
from typing import Optional

# 请求模型
class UserCreate(BaseModel):
    email: EmailStr
    name: str = Field(..., min_length=2, max_length=100)
    password: str = Field(..., min_length=8)

# 响应模型
class UserResponse(BaseModel):
    id: int
    email: str
    name: str
    created_at: datetime

    class Config:
        from_attributes = True

# 更新模型
class UserUpdate(BaseModel):
    name: Optional[str] = Field(None, min_length=2, max_length=100)
    email: Optional[EmailStr] = None
```

### 4. 统一响应格式

```python
from typing import Generic, TypeVar, Optional
from pydantic import BaseModel

T = TypeVar("T")

class ApiResponse(BaseModel, Generic[T]):
    success: bool = True
    data: Optional[T] = None
    error: Optional[str] = None
    message: Optional[str] = None

    @classmethod
    def ok(cls, data: T, message: str = None) -> "ApiResponse[T]":
        return cls(success=True, data=data, message=message)

    @classmethod
    def error(cls, error: str) -> "ApiResponse[T]":
        return cls(success=False, error=error)
```

### 5. 异常处理

```python
from fastapi import Request, HTTPException
from fastapi.responses import JSONResponse

class AppException(Exception):
    def __init__(self, status_code: int, detail: str):
        self.status_code = status_code
        self.detail = detail

# 全局异常处理器
@app.exception_handler(AppException)
async def app_exception_handler(request: Request, exc: AppException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"success": False, "error": exc.detail}
    )

@app.exception_handler(HTTPException)
async def http_exception_handler(request: Request, exc: HTTPException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"success": False, "error": exc.detail}
    )
```

### 6. 仓储模式

```python
from typing import Optional, List

class UserRepository:
    def __init__(self, db: Session):
        self.db = db

    def find_by_id(self, id: int) -> Optional[User]:
        return self.db.query(User).filter(User.id == id).first()

    def find_by_email(self, email: str) -> Optional[User]:
        return self.db.query(User).filter(User.email == email).first()

    def find_all(self, skip: int = 0, limit: int = 100) -> List[User]:
        return self.db.query(User).offset(skip).limit(limit).all()

    def create(self, user: User) -> User:
        self.db.add(user)
        self.db.commit()
        self.db.refresh(user)
        return user

    def delete(self, id: int) -> bool:
        user = self.find_by_id(id)
        if user:
            self.db.delete(user)
            self.db.commit()
            return True
        return False
```

### 7. 配置管理

```python
from pydantic_settings import BaseSettings
from functools import lru_cache

class Settings(BaseSettings):
    app_name: str = "My API"
    debug: bool = False
    database_url: str
    secret_key: str
    jwt_expiration_hours: int = 24

    class Config:
        env_file = ".env"

@lru_cache()
def get_settings() -> Settings:
    return Settings()
```

### 8. 后台任务

```python
from fastapi import BackgroundTasks

def send_email(to: str, subject: str, body: str):
    # 发送邮件逻辑
    pass

@router.post("/users/")
async def create_user(
    user_in: UserCreate,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db)
):
    user = create_user_in_db(db, user_in)
    background_tasks.add_task(
        send_email,
        to=user.email,
        subject="Welcome",
        body="Welcome to our service"
    )
    return ApiResponse.ok(UserResponse.model_validate(user))
```

## 使用场景

- 创建新的 FastAPI 项目
- REST API 开发
- 依赖注入设计
- 异常处理
