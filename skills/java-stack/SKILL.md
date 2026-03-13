---
name: java-stack
description: Use when working on Java backend systems, Spring Boot, Quarkus, dependency injection, transaction boundaries, JVM service structure, or Java-specific framework concerns.
---

# Java 技术栈

这个 skill 负责 Java 服务端项目的入口判断和框架经验索引。主文档只保留适用范围、关键关注点和 references 入口。

## 核心规则

- 先判断问题属于框架装配、事务、配置、Web 层、数据访问还是运行时行为
- 通用测试策略看 `testing`；Java 特有测试细节看 `references/testing.md`
- 详细实践优先写入 `references/`，不要在主文档重复模板代码

## 重点关注

- Spring Boot / Quarkus 配置与装配
- 依赖注入、事务边界、异常处理
- Controller / Service / Repository 分层
- JVM 服务运行与配置管理

## 参考资料

- [springboot.md](references/springboot.md)
- [testing.md](references/testing.md)
