---
name: database
description: |
  数据库技能。当用户要求数据库、sql、migration、schema、查询优化时自动触发。
  包含：(1) PostgreSQL (2) MySQL (3) MongoDB (4) Redis
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - database
    - sql
    - migration
---

# 数据库

## 支持的数据库

| 数据库 | ORM/驱动 | 适用场景 |
|--------|----------|----------|
| PostgreSQL | GORM/SQLx/TypeORM | 关系型数据 |
| MySQL | GORM/TypeORM | 关系型数据 |
| MongoDB | Mongoose | 文档存储 |
| Redis | go-redis/redis-py | 缓存 |

## 最佳实践

### 连接池

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 查询优化

- 使用索引
- 避免 SELECT *
- 分页查询
- 避免 N+1 问题

### Migration

```
# 创建迁移
migrate create -ext sql -dir migrations create_users_table

# 执行迁移
migrate -path migrations -database $DB_URL up
```

## 相关技能

- `java-stack` - JPA/Repository 模式
- `python-stack` - SQLAlchemy 模式
- `go-stack` - GORM 模式
