# Database Skills

数据库相关技能，包括 PostgreSQL、MySQL、Redis 等。

## 何时使用

- 数据库设计
- 查询优化
- 迁移管理
- 性能调优

## 包含内容

### 1. PostgreSQL 模式

#### 连接池配置
```python
# SQLAlchemy
engine = create_engine(
    DATABASE_URL,
    pool_size=10,
    max_overflow=20,
    pool_pre_ping=True
)
```

```java
// Spring Boot HikariCP
spring:
  datasource:
    hikari:
      maximum-pool-size: 20
      minimum-idle: 5
```

#### 查询优化
- 使用 `EXPLAIN ANALYZE` 分析查询
- 创建合适的索引
- 避免 `SELECT *`
- 使用连接池

### 2. 数据库迁移

#### Flyway (Java)
```
db/migration/
├── V1__create_users_table.sql
├── V2__add_email_index.sql
└── V3__create_orders_table.sql
```

#### Alembic (Python)
```bash
alembic revision --autogenerate -m "add user table"
alembic upgrade head
```

#### golang-migrate
```bash
migrate create -ext sql -dir migrations -seq create_users
migrate -database $DATABASE_URL -path migrations up
```

### 3. Redis 缓存

#### 缓存模式
```python
# Cache-Aside
def get_user(user_id):
    cached = redis.get(f"user:{user_id}")
    if cached:
        return cached

    user = db.query(User).get(user_id)
    redis.setex(f"user:{user_id}", 3600, user)
    return user
```

#### 常用模式
- 缓存穿透：使用空值缓存
- 缓存雪崩：设置随机过期时间
- 缓存击穿：使用分布式锁

### 4. 性能优化

| 问题 | 解决方案 |
|------|----------|
| N+1 查询 | 使用 JOIN 或批量查询 |
| 缺少索引 | 分析慢查询，创建索引 |
| 大表扫描 | 分区表、添加索引 |
| 连接数过多 | 使用连接池 |

## 相关命令

- `/gclm:review --scope performance` - 性能审查
