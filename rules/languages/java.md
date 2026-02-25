# Java/Spring Boot 规则

Java 和 Spring Boot 项目的编码规范和最佳实践。

## Java 编码规范

### 1. 命名规范
```java
// 类名：PascalCase
public class UserService {}

// 接口名：PascalCase，可以用 I 前缀或 able 后缀
public interface UserRepository {}
public interface Serializable {}

// 方法名：camelCase，动词开头
public User findById(Long id) {}
public boolean isValid() {}

// 常量：UPPER_SNAKE_CASE
public static final String DEFAULT_CHARSET = "UTF-8";

// 包名：全小写
package com.example.user.service;
```

### 2. 代码结构
```java
public class ExampleClass {
    // 1. 静态常量
    private static final Logger log = LoggerFactory.getLogger(ExampleClass.class);

    // 2. 实例变量
    private final UserRepository userRepository;

    // 3. 构造函数
    public ExampleClass(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    // 4. 公共方法
    public User getUser(Long id) {
        // ...
    }

    // 5. 私有方法
    private void validateId(Long id) {
        // ...
    }

    // 6. 内部类
    private static class Helper {}
}
```

### 3. 空值处理
```java
// 使用 Optional
public Optional<User> findById(Long id) {
    return userRepository.findById(id);
}

// 使用 Objects.requireNonNull
public void setName(String name) {
    this.name = Objects.requireNonNull(name, "name must not be null");
}

// 使用 @Nullable 和 @NonNull 注解
public void process(@NonNull String input, @Nullable String optional) {}
```

## Spring Boot 最佳实践

### 1. 项目结构
```
src/main/java/com/example/
├── config/          # 配置类
├── controller/      # REST 控制器
├── service/         # 业务逻辑
├── repository/      # 数据访问
├── entity/          # 实体类
├── dto/             # 数据传输对象
├── exception/       # 异常处理
└── util/            # 工具类
```

### 2. 依赖注入
```java
// 推荐：构造函数注入
@Service
@RequiredArgsConstructor // Lombok
public class UserService {
    private final UserRepository userRepository;
    private final EmailService emailService;
}

// 不推荐：字段注入
@Service
public class UserService {
    @Autowired
    private UserRepository userRepository; // 不推荐
}
```

### 3. REST API 设计
```java
@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponse<UserDTO>> getUser(@PathVariable Long id) {
        return userService.findById(id)
            .map(user -> ResponseEntity.ok(ApiResponse.success(user)))
            .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public ApiResponse<UserDTO> createUser(@Valid @RequestBody CreateUserRequest request) {
        return ApiResponse.success(userService.create(request));
    }
}
```

### 4. 统一响应格式
```java
@Data
@AllArgsConstructor
@NoArgsConstructor
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
}
```

### 6. 配置管理
```java
// 使用 @ConfigurationProperties
@ConfigurationProperties(prefix = "app")
@Data
public class AppProperties {
    private String name;
    private int maxConnections;
    private Duration timeout;
}

// application.yml
app:
  name: my-app
  max-connections: 100
  timeout: 30s
```

### 7. 数据访问
```java
// 使用 Spring Data JPA
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);

    @Query("SELECT u FROM User u WHERE u.status = :status")
    List<User> findByStatus(@Param("status") UserStatus status);
}

// Service 层
@Service
@Transactional(readOnly = true)
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;

    @Transactional
    public User create(User user) {
        return userRepository.save(user);
    }
}
```

## 测试规范

### 1. 单元测试
```java
@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserRepository userRepository;

    @InjectMocks
    private UserService userService;

    @Test
    void findById_shouldReturnUser_whenUserExists() {
        // Arrange
        User user = new User(1L, "test@example.com");
        when(userRepository.findById(1L)).thenReturn(Optional.of(user));

        // Act
        Optional<User> result = userService.findById(1L);

        // Assert
        assertThat(result).isPresent();
        assertThat(result.get().getEmail()).isEqualTo("test@example.com");
    }
}
```

### 2. 集成测试
```java
@SpringBootTest
@AutoConfigureMockMvc
class UserControllerIntegrationTest {

    @Autowired
    private MockMvc mockMvc;

    @Test
    void getUser_shouldReturn200_whenUserExists() throws Exception {
        mockMvc.perform(get("/api/v1/users/1"))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.success").value(true))
            .andExpect(jsonPath("$.data.id").value(1));
    }
}
```

## 性能优化

### 1. 数据库优化
- 使用 @EntityGraph 解决 N+1 问题
- 合理使用 @BatchSize
- 使用投影减少数据传输
- 分页查询

### 2. 缓存
```java
@Service
public class UserService {

    @Cacheable(value = "users", key = "#id")
    public User findById(Long id) {
        return userRepository.findById(id).orElseThrow();
    }

    @CacheEvict(value = "users", key = "#user.id")
    public User update(User user) {
        return userRepository.save(user);
    }
}
```

### 3. 异步处理
```java
@Service
public class NotificationService {

    @Async
    public CompletableFuture<Void> sendEmailAsync(String to, String subject, String body) {
        // 异步发送邮件
        return CompletableFuture.completedFuture(null);
    }
}
```
