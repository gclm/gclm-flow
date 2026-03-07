---
name: database
description: Use when working on SQL, schema changes, migrations, indexes, query performance, transaction behavior, or persistence-layer design across PostgreSQL, MySQL, MongoDB, or Redis.
---

# 数据库

这个 skill 负责数据库层的设计判断、变更风险、查询性能和迁移检查，不替代具体语言栈 skill。

## 核心规则

- 先确认变更属于哪一类：查询、schema、migration、索引、事务、缓存
- 改 schema 或 migration 时，优先考虑兼容性、数据安全、回滚路径
- 优化查询前先明确瓶颈，不凭感觉加索引或上缓存
- 任何会影响线上数据正确性的改动，都默认按高风险处理

## 常见场景

### 查询与性能

关注：
- 是否有 N+1
- 是否全表扫描、未分页、`SELECT *`
- 是否有不必要的排序、聚合、重复查询
- 是否应该用索引、批处理、缓存，而不是单点修补

### Schema 与 Migration

关注：
- 是否需要兼容窗口
- 默认值、非空约束、唯一约束是否会阻塞线上数据
- 是否需要双写、回填、分阶段迁移
- 回滚时数据是否仍可读、可写、可恢复

### 事务与一致性

关注：
- 事务边界是否过大或过小
- 是否存在部分提交、重试副作用、并发写冲突
- 缓存更新和数据库写入是否可能失一致

## 使用顺序

1. 先识别数据库类型和具体变更面。
2. 明确风险：正确性、性能、兼容性、回滚。
3. 只做最小但正确的设计或修复。
4. 通过 query plan、测试、迁移验证来证明结论。

## 参考资料

- [query-review-checklist.md](references/query-review-checklist.md)
- [migration-safety-checklist.md](references/migration-safety-checklist.md)

## 联动技能

- `code-review`
- `testing`
- `devops`
- 对应语言栈 skill
