# FastAPI 最佳实践

## 路由定义

```python
from fastapi import APIRouter, HTTPException, status
from pydantic import BaseModel

router = APIRouter(prefix="/users", tags=["users"])

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

## 后台任务

```python
from fastapi import BackgroundTasks

def send_email(to: str, subject: str, body: str):
    pass

@router.post("/users/")
async def create_user(
    user_in: UserCreate,
    background_tasks: BackgroundTasks
):
    user = create_user_in_db(user_in)
    background_tasks.add_task(
        send_email,
        to=user.email,
        subject="Welcome"
    )
    return user
```

## 仓储模式

```python
from typing import Optional, List

class UserRepository:
    def __init__(self, db: Session):
        self.db = db

    def find_by_id(self, id: int) -> Optional[User]:
        return self.db.query(User).filter(User.id == id).first()

    def find_by_email(self, email: str) -> Optional[User]:
        return self.db.query(User).filter(User.email == email).first()

    def create(self, user: User) -> User:
        self.db.add(user)
        self.db.commit()
        self.db.refresh(user)
        return user
```
