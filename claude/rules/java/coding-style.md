---
paths:
  - "**/*.java"
  - "**/*.kt"
  - "**/pom.xml"
  - "**/build.gradle"
  - "**/build.gradle.kts"
---
# Java Coding Style

> This file extends [common/coding-style.md](../common/coding-style.md) with Java specific content.

## Immutability

- 优先使用 `final` 字段和不可变对象
- 使用 `record`（Java 16+）替代简单的 POJO
- 集合使用 `List.of()`、`Map.of()` 等不可变工厂方法

## 文件组织

- 单文件单 public class
- 包结构按功能域划分，不按类型划分（避免 `controllers/`、`services/` 平铺）
- 单个类不超过 400 行，超过时拆分

## 错误处理

- 业务异常使用自定义 checked/unchecked exception，明确区分
- 不吞异常：catch 后必须 log 或 rethrow
- Spring 项目统一用 `@ControllerAdvice` 处理异常，不在 controller 层散落 try-catch

## 代码质量检查

完成前检查：
- [ ] 无 `System.out.println`（用日志框架）
- [ ] 无裸 `catch (Exception e) {}`
- [ ] 无未使用的 import
- [ ] 字段/方法访问权限最小化
- [ ] 使用 `Optional` 替代返回 null
