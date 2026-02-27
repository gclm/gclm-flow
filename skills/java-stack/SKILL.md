---
name: java-stack
description: |
  Java/Spring Boot 技术栈完整开发指南。当检测到 Java 项目（pom.xml、build.gradle）
  或用户明确要求 Java/Spring Boot/Quarkus 开发时自动触发。包含：
  (1) 项目结构规范 (2) Spring Boot 最佳实践 (3) 测试模式 (4) 安全规范 (5) 性能优化
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - java
    - spring-boot
    - quarkus
---

# Java 技术栈开发指南

## 框架检测

检测项目使用的框架并加载对应参考文档：

- 存在 `spring-boot-starter` 依赖 → Spring Boot，详见 [springboot.md](references/springboot.md)
- 存在 `quarkus-bom` 依赖 → Quarkus，详见 [quarkus.md](references/quarkus.md)
- 测试相关 → 详见 [testing.md](references/testing.md)

## 标准项目结构

```
src/main/java/com/example/
├── config/           # 配置类
│   ├── SecurityConfig.java
│   ├── WebConfig.java
│   └── RedisConfig.java
├── controller/       # REST 控制器
├── service/          # 业务逻辑
├── repository/       # 数据访问
├── entity/           # JPA 实体
├── dto/              # 数据传输对象
│   ├── request/
│   └── response/
├── exception/        # 异常处理
│   ├── GlobalExceptionHandler.java
│   └── BusinessException.java
└── util/             # 工具类
```

## 核心规范

### 依赖注入

```java
// 推荐：构造函数注入 + Lombok
@Service
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;
    private final EmailService emailService;
}
```

### 统一响应格式

```java
@Data
@AllArgsConstructor
public class ApiResponse<T> {
    private boolean success;
    private T data;
    private String error;
    private LocalDateTime timestamp;

    public static <T> ApiResponse<T> success(T data) {
        return new ApiResponse<>(true, data, null, LocalDateTime.now());
    }

    public static <T> ApiResponse<T> error(String error) {
        return new ApiResponse<>(false, null, error, LocalDateTime.now());
    }
}
```

### REST 控制器

```java
@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping
    public ApiResponse<Page<UserDTO>> list(
        @RequestParam(defaultValue = "0") int page,
        @RequestParam(defaultValue = "20") int size
    ) {
        return ApiResponse.success(userService.findAll(page, size));
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public ApiResponse<UserDTO> create(@Valid @RequestBody CreateUserRequest request) {
        return ApiResponse.success(userService.create(request));
    }
}
```

### 异常处理

```java
@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(EntityNotFoundException.class)
    public ResponseEntity<ApiResponse<?>> handleNotFound(EntityNotFoundException e) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND)
            .body(ApiResponse.error(e.getMessage()));
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ApiResponse<?>> handleValidation(MethodArgumentNotValidException e) {
        String message = e.getBindingResult().getFieldErrors().stream()
            .map(FieldError::getDefaultMessage)
            .collect(Collectors.joining(", "));
        return ResponseEntity.badRequest().body(ApiResponse.error(message));
    }
}
```

## 测试规范

详见 [testing.md](references/testing.md)

- 使用 JUnit 5 + Mockito
- TDD：RED-GREEN-REFACTOR 循环
- 目标覆盖率：80%+

## 最佳实践

详见 [springboot.md](references/springboot.md)

- 配置管理：使用 `@ConfigurationProperties`
- 事件驱动：使用 `ApplicationEvent`
- 缓存：使用 `@Cacheable`
- 异步：使用 `@Async`

## 相关技能

- `code-review` - Java 代码审查
- `testing` - Java 测试模式
- `database` - JPA/Repository 模式
