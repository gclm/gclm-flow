# 性能规则

所有项目通用的性能优化规范。

## 通用原则

### 性能优先级
1. **正确性**：先保证功能正确
2. **可读性**：代码要易于理解
3. **性能**：在需要时优化

### 过早优化
> "Premature optimization is the root of all evil" — Donald Knuth

- 不要过早优化
- 基于测量进行优化
- 优化热点代码

## 数据库性能

### 查询优化
```sql
-- ✅ 使用索引
CREATE INDEX idx_users_email ON users(email);

-- ✅ 只查询需要的列
SELECT id, name FROM users WHERE active = true;

-- ❌ 避免 SELECT *
SELECT * FROM users;
```

### N+1 问题
```python
# ❌ N+1 查询
for user in users:
    orders = db.query(Order).filter(Order.user_id == user.id).all()

# ✅ 批量查询
user_ids = [u.id for u in users]
orders = db.query(Order).filter(Order.user_id.in_(user_ids)).all()
```

### 连接池
```python
# 配置连接池
engine = create_engine(
    DATABASE_URL,
    pool_size=10,
    max_overflow=20,
    pool_pre_ping=True
)
```

## 缓存策略

### 缓存层次
```
请求 → 本地缓存 → 分布式缓存 → 数据库
```

### 缓存模式

#### Cache-Aside
```python
def get_user(user_id):
    # 先查缓存
    cached = redis.get(f"user:{user_id}")
    if cached:
        return cached

    # 缓存未命中，查数据库
    user = db.query(User).get(user_id)

    # 写入缓存
    redis.setex(f"user:{user_id}", 3600, user)
    return user
```

#### 缓存问题
| 问题 | 解决方案 |
|------|----------|
| 缓存穿透 | 缓存空值、布隆过滤器 |
| 缓存雪崩 | 随机过期时间、多级缓存 |
| 缓存击穿 | 分布式锁、永不过期 |

## 代码优化

### 算法复杂度
- 了解常见算法的时间复杂度
- 避免在循环中进行数据库查询
- 使用合适的数据结构

### 内存管理
```python
# ✅ 使用生成器处理大数据
def process_large_file():
    with open('large_file.txt') as f:
        for line in f:
            yield process_line(line)

# ❌ 一次性加载所有数据
def process_large_file_bad():
    with open('large_file.txt') as f:
        lines = f.readlines()  # 占用大量内存
```

### 懒加载
```python
# 懒加载属性
@property
def expensive_data(self):
    if self._expensive_data is None:
        self._expensive_data = self._compute_expensive_data()
    return self._expensive_data
```

## 并发处理

### 异步操作
```python
# FastAPI 异步
@app.get("/users/{user_id}")
async def get_user(user_id: int):
    user = await user_service.get_by_id(user_id)
    return user
```

### 后台任务
```python
# 非阻塞发送邮件
@app.post("/users/")
async def create_user(
    user: UserCreate,
    background_tasks: BackgroundTasks
):
    background_tasks.add_task(send_welcome_email, user.email)
    return user
```

## API 性能

### 分页
```python
# 使用游标分页（大数据集）
def get_users(cursor=None, limit=20):
    query = db.query(User)
    if cursor:
        query = query.filter(User.id > cursor)
    return query.limit(limit).all()
```

### 压缩
- 启用 Gzip 压缩
- 压缩响应体
- 使用 CDN

### 批量操作
```python
# ✅ 批量插入
db.bulk_insert_mappings(User, users_data)

# ❌ 逐条插入
for user_data in users_data:
    db.add(User(**user_data))
```

## 监控

### 关键指标
| 指标 | 描述 |
|------|------|
| 响应时间 | P50, P95, P99 |
| 吞吐量 | QPS/RPS |
| 错误率 | 失败请求比例 |
| 资源使用 | CPU、内存、I/O |

### 性能分析
- 使用 APM 工具
- 记录慢查询
- 分析热点函数
