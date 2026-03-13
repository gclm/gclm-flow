---
paths:
  - "**/*.java"
  - "**/pom.xml"
  - "**/build.gradle.kts"
---
# Java Testing

> This file extends [common/testing.md](../common/testing.md) with Java specific content.

## Framework

- 单元测试：JUnit 5 + Mockito
- Spring 集成测试：`@SpringBootTest` + `@MockBean`
- 接口测试：REST Assured 或 MockMvc

## 测试命名

```java
@Test
void should_returnUser_when_idIsValid() { ... }
```

## 覆盖率

```bash
# Maven
mvn test jacoco:report

# Gradle
./gradlew test jacocoTestReport
```

目标：核心业务逻辑 80%+ 覆盖率。

## 原则

- 单元测试不启动 Spring 容器（用纯 Mockito）
- 集成测试用 `@DataJpaTest`、`@WebMvcTest` 等切片测试，不用完整 `@SpringBootTest`
- 测试数据用 Builder 模式或 Test Fixture，不用生产数据

## Reference

See skill: `java-stack` for detailed Spring Boot patterns.
