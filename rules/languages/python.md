# Python/Flask/FastAPI 规则

Python、Flask、FastAPI 项目的编码规范和最佳实践。

## Python 编码规范

### 1. 命名规范
```python
# 类名：PascalCase
class UserService:
    pass

# 函数/方法名：snake_case
def get_user_by_id(user_id: int) -> User:
    pass

# 常量：UPPER_SNAKE_CASE
MAX_RETRY_COUNT = 3
DEFAULT_TIMEOUT = 30

# 私有属性：_前缀
class Example:
    def __init__(self):
        self._internal_value = None

# 模块级私有：__前缀
def __internal_helper():
    pass
```

### 2. 类型注解
```python
from typing import Optional, List, Dict, Any

def get_user(user_id: int) -> Optional[User]:
    """获取用户信息"""
    pass

def process_items(items: List[str]) -> Dict[str, Any]:
    """处理项目列表"""
    pass

# Python 3.10+ 联合类型
def get_value(key: str) -> str | None:
    pass
```

### 3. 文档字符串
```python
def create_user(email: str, name: str) -> User:
    """创建新用户。

    Args:
        email: 用户邮箱地址
        name: 用户名称

    Returns:
        创建的用户对象

    Raises:
        ValueError: 如果邮箱格式无效
        DuplicateEmailError: 如果邮箱已存在
    """
    pass
```

### 4. 上下文管理器
```python
# 使用 with 语句管理资源
with open('file.txt', 'r') as f:
    content = f.read()

# 自定义上下文管理器
from contextlib import contextmanager

@contextmanager
def get_db_session():
    session = Session()
    try:
        yield session
        session.commit()
    except Exception:
        session.rollback()
        raise
    finally:
        session.close()
```

## FastAPI 最佳实践

### 1. 项目结构
```
app/
├── main.py              # 应用入口
├── config.py            # 配置管理
├── dependencies.py      # 依赖注入
├── routers/             # 路由模块
│   ├── users.py
│   └── auth.py
├── services/            # 业务逻辑
│   └── user_service.py
├── models/              # 数据模型
│   ├── user.py
│   └── schemas.py
├── repositories/        # 数据访问
│   └── user_repo.py
└── utils/               # 工具函数
    └── security.py
```

### 2. 依赖注入
```python
from fastapi import Depends

# 数据库依赖
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

# 当前用户依赖
async def get_current_user(
    token: str = Depends(oauth2_scheme),
    db: Session = Depends(get_db)
) -> User:
    payload = decode_token(token)
    user = db.query(User).filter(User.id == payload["sub"]).first()
    if not user:
        raise HTTPException(status_code=401, detail="Invalid token")
    return user

# 在路由中使用
@router.get("/me")
async def get_me(current_user: User = Depends(get_current_user)):
    return current_user
```

### 3. 路由定义
```python
from fastapi import APIRouter, HTTPException, status
from pydantic import BaseModel

router = APIRouter(prefix="/users", tags=["users"])

# 请求/响应模型
class UserCreate(BaseModel):
    email: str
    name: str

class UserResponse(BaseModel):
    id: int
    email: str
    name: str

    class Config:
        from_attributes = True

@router.post("/", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(user_in: UserCreate, db: Session = Depends(get_db)):
    """创建新用户"""
    if db.query(User).filter(User.email == user_in.email).first():
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="Email already registered"
        )
    user = User(**user_in.model_dump())
    db.add(user)
    db.commit()
    return user
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

# 使用
@router.get("/{user_id}", response_model=ApiResponse[UserResponse])
async def get_user(user_id: int, db: Session = Depends(get_db)):
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        return ApiResponse.error("User not found")
    return ApiResponse.ok(UserResponse.model_validate(user))
```

### 5. 异常处理
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
        content={
            "success": False,
            "error": exc.detail,
        }
    )

# 使用
@router.get("/{user_id}")
async def get_user(user_id: int, db: Session = Depends(get_db)):
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise AppException(status_code=404, detail="User not found")
    return user
```

### 6. 配置管理
```python
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    app_name: str = "My API"
    debug: bool = False
    database_url: str
    secret_key: str

    class Config:
        env_file = ".env"

settings = Settings()
```

## Flask 最佳实践

### 1. 应用工厂
```python
from flask import Flask

def create_app(config_name="default"):
    app = Flask(__name__)
    app.config.from_object(config[config_name])

    # 注册蓝图
    from .api import api_bp
    app.register_blueprint(api_bp, url_prefix="/api")

    # 初始化扩展
    db.init_app(app)

    return app
```

### 2. 蓝图组织
```python
from flask import Blueprint, jsonify, request

users_bp = Blueprint("users", __name__)

@users_bp.route("/", methods=["GET"])
def list_users():
    users = User.query.all()
    return jsonify([u.to_dict() for u in users])

@users_bp.route("/", methods=["POST"])
def create_user():
    data = request.get_json()
    user = User(**data)
    db.session.add(user)
    db.session.commit()
    return jsonify(user.to_dict()), 201
```

## 测试规范

### 1. pytest 单元测试
```python
import pytest
from unittest.mock import Mock, patch

@pytest.fixture
def mock_db():
    return Mock()

class TestUserService:
    def test_get_user_by_id(self, mock_db):
        # Arrange
        mock_db.query.return_value.filter.return_value.first.return_value = User(id=1)

        # Act
        service = UserService(mock_db)
        result = service.get_by_id(1)

        # Assert
        assert result.id == 1

    def test_get_user_not_found(self, mock_db):
        mock_db.query.return_value.filter.return_value.first.return_value = None

        service = UserService(mock_db)

        with pytest.raises(NotFoundError):
            service.get_by_id(999)
```

### 2. FastAPI 测试
```python
from fastapi.testclient import TestClient

def test_create_user(client: TestClient):
    response = client.post(
        "/users/",
        json={"email": "test@example.com", "name": "Test User"}
    )
    assert response.status_code == 201
    data = response.json()
    assert data["email"] == "test@example.com"
```

## 性能优化

### 1. 异步处理
```python
# FastAPI 异步
@router.get("/users/{user_id}")
async def get_user(user_id: int, db: Session = Depends(get_db)):
    # 使用 async/await
    user = await user_service.get_by_id(user_id)
    return user

# 后台任务
from fastapi import BackgroundTasks

def send_email(to: str, subject: str, body: str):
    # 发送邮件
    pass

@router.post("/users/")
async def create_user(
    user_in: UserCreate,
    background_tasks: BackgroundTasks
):
    user = await create_user_in_db(user_in)
    background_tasks.add_task(
        send_email,
        to=user.email,
        subject="Welcome",
        body="Welcome to our service"
    )
    return user
```

### 2. 缓存
```python
from functools import lru_cache
from cachetools import cached, TTLCache

# 简单缓存
@lru_cache(maxsize=128)
def get_config(key: str):
    return db.query(Config).filter(Config.key == key).first()

# TTL 缓存
cache = TTLCache(maxsize=100, ttl=300)

@cached(cache)
def get_user_permissions(user_id: int):
    return db.query(Permission).filter(Permission.user_id == user_id).all()
```
