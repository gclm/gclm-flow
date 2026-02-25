# Spring Boot Patterns

Spring Boot 框架专用模式和最佳实践。

## 技能描述

这个技能包含 Spring Boot 开发的专用模式，适用于 Java 后端开发。

## 包含的模式

### 1. 项目结构

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

### 2. 依赖注入最佳实践

```java
// 推荐：构造函数注入 + Lombok
@Service
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;
    private final EmailService emailService;
    private final ApplicationEventPublisher eventPublisher;
}
```

### 3. 统一响应格式

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

### 4. REST 控制器

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

    @GetMapping("/{id}")
    public ApiResponse<UserDTO> get(@PathVariable Long id) {
        return ApiResponse.success(userService.findById(id));
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public ApiResponse<UserDTO> create(@Valid @RequestBody CreateUserRequest request) {
        return ApiResponse.success(userService.create(request));
    }
}
```

### 5. 异常处理

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
        return ResponseEntity.badRequest()
            .body(ApiResponse.error(message));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiResponse<?>> handleGeneric(Exception e) {
        log.error("Unexpected error", e);
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
            .body(ApiResponse.error("Internal server error"));
    }
}
```

### 6. JPA 最佳实践

```java
@Entity
@Table(name = "users")
@EntityListeners(AuditingEntityListener.class)
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, unique = true)
    private String email;

    @Column(nullable = false)
    private String name;

    @CreatedDate
    @Column(updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    private LocalDateTime updatedAt;

    @Version
    private Long version;
}

// Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);

    @Query("SELECT u FROM User u WHERE u.status = :status")
    Page<User> findByStatus(@Param("status") UserStatus status, Pageable pageable);
}
```

### 7. 配置管理

```java
@Data
@ConfigurationProperties(prefix = "app")
public class AppProperties {
    private String name;
    private Security security = new Security();
    private Jwt jwt = new Jwt();

    @Data
    public static class Security {
        private boolean enabled = true;
    }

    @Data
    public static class Jwt {
        private String secret;
        private Duration expiration = Duration.ofHours(24);
    }
}
```

### 8. 事件驱动

```java
// 定义事件
public class UserCreatedEvent extends ApplicationEvent {
    private final User user;

    public UserCreatedEvent(Object source, User user) {
        super(source);
        this.user = user;
    }
}

// 发布事件
@Service
@RequiredArgsConstructor
public class UserService {
    private final ApplicationEventPublisher eventPublisher;

    public User create(CreateUserRequest request) {
        User user = // ... 创建用户
        eventPublisher.publishEvent(new UserCreatedEvent(this, user));
        return user;
    }
}

// 监听事件
@Component
public class UserEventListener {

    @EventListener
    @Async
    public void handleUserCreated(UserCreatedEvent event) {
        // 发送欢迎邮件等
    }
}
```

## 使用场景

- 创建新的 Spring Boot 项目
- REST API 开发
- 异常处理设计
- 配置管理
