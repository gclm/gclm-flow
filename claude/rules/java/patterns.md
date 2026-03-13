---
paths:
  - "**/*.java"
---
# Java Patterns

> This file extends [common/patterns.md](../common/patterns.md) with Java specific content.

## Spring Boot 结构

```
src/main/java/com/example/
├── domain/          # 领域模型、业务逻辑
├── application/     # 用例、服务层
├── infrastructure/  # 数据库、外部 API 适配
└── interfaces/      # Controller、DTO
```

## 常用模式

- **Repository Pattern**：数据访问层统一用 Spring Data Repository
- **DTO 分离**：Controller 入参/出参用独立 DTO，不暴露 Entity
- **Builder**：复杂对象构建用 Lombok `@Builder`
- **事务边界**：`@Transactional` 放在 service 层，不放 repository 层

## 禁止

- 禁止在 Entity 上直接加 Jackson 注解（DTO 分离）
- 禁止 service 直接依赖 controller 层
- 禁止跨模块直接访问 repository（通过 service 调用）

## Reference

See skill: `java-stack` for Spring Boot implementation patterns.
