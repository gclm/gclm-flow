# Java/Spring Boot 测试模式

使用 JUnit 5、Mockito 和 TDD 方法论的 Java/Spring Boot 应用测试策略。

## 何时激活

- 编写新的 Java/Spring Boot 代码
- 设计 Java 项目测试套件
- 审查 Java 测试覆盖率
- 搭建测试基础设施

## TDD 工作流

### RED-GREEN-REFACTOR 循环

```
RED     → 先写失败的测试
GREEN   → 写最少代码通过测试
REFACTOR → 重构代码，保持测试通过
REPEAT  → 继续下一个需求
```

### TDD 步骤示例

```java
// 步骤 1: 编写失败的测试 (RED)
@Test
void shouldAddTwoNumbers() {
    Calculator calculator = new Calculator();
    int result = calculator.add(2, 3);
    assertEquals(5, result);
}

// 步骤 2: 实现最少代码 (GREEN)
public class Calculator {
    public int add(int a, int b) {
        return a + b;
    }
}

// 步骤 3: 按需重构
```

## JUnit 5 基础

### 基本测试结构

```java
import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;

class UserServiceTest {

    @BeforeEach
    void setUp() {
        // 每个测试前的准备工作
    }

    @AfterEach
    void tearDown() {
        // 每个测试后的清理工作
    }

    @Test
    @DisplayName("应该使用有效数据创建用户")
    void shouldCreateUser() {
        User user = new User("Alice", "alice@example.com");

        assertEquals("Alice", user.getName());
        assertEquals("alice@example.com", user.getEmail());
    }

    @Test
    void multipleAssertions() {
        User user = new User("Bob", "bob@example.com");

        assertAll("user",
            () -> assertEquals("Bob", user.getName()),
            () -> assertEquals("bob@example.com", user.getEmail()),
            () -> assertNotNull(user.getId())
        );
    }
}
```

### 断言

```java
// 相等性
assertEquals(expected, actual);
assertEquals(expected, actual, "自定义消息");
assertNotEquals(unexpected, actual);

// 布尔值
assertTrue(condition);
assertFalse(condition);

// 空值检查
assertNull(value);
assertNotNull(value);

// 异常
assertThrows(IllegalArgumentException.class, () -> {
    service.doSomethingInvalid();
});

Exception ex = assertThrows(RuntimeException.class, () -> {
    throw new RuntimeException("错误消息");
});
assertEquals("错误消息", ex.getMessage());

// 超时
assertTimeout(Duration.ofSeconds(2), () -> {
    slowOperation();
});
```

## 参数化测试

```java
import org.junit.jupiter.params.*;
import org.junit.jupiter.params.provider.*;

@ParameterizedTest
@ValueSource(strings = {"hello", "world", "test"})
void shouldValidateNonEmptyStrings(String input) {
    assertTrue(Validator.isNotEmpty(input));
}

@ParameterizedTest
@CsvSource({
    "1, 2, 3",
    "10, 20, 30",
    "-5, 5, 0"
})
void shouldAddNumbers(int a, int b, int expected) {
    assertEquals(expected, calculator.add(a, b));
}

@ParameterizedTest
@CsvFileSource(resources = "/test-data.csv", numLinesToSkip = 1)
void shouldProcessCsvData(String input, String expected) {
    assertEquals(expected, processor.process(input));
}
```

## Spring Boot 测试

### 测试配置

```java
@SpringBootTest
@TestPropertySource(locations = "classpath:application-test.properties")
class ApplicationIntegrationTest {

    @Autowired
    private UserRepository userRepository;

    @Test
    void shouldLoadContext() {
        assertNotNull(userRepository);
    }
}
```

### Web 层测试

```java
@WebMvcTest(UserController.class)
class UserControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private UserService userService;

    @Test
    void shouldReturnUserById() throws Exception {
        User user = new User(1L, "Alice");
        when(userService.findById(1L)).thenReturn(user);

        mockMvc.perform(get("/api/users/1"))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.name").value("Alice"));
    }

    @Test
    void shouldCreateUser() throws Exception {
        User saved = new User(1L, "Bob");
        when(userService.save(any())).thenReturn(saved);

        mockMvc.perform(post("/api/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content("{\"name\":\"Bob\"}"))
            .andExpect(status().isCreated())
            .andExpect(jsonPath("$.id").value(1));
    }
}
```

### 数据层测试

```java
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
class UserRepositoryTest {

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private TestEntityManager entityManager;

    @Test
    void shouldFindUserByEmail() {
        User user = new User("Alice", "alice@test.com");
        entityManager.persist(user);

        Optional<User> found = userRepository.findByEmail("alice@test.com");

        assertTrue(found.isPresent());
        assertEquals("Alice", found.get().getName());
    }
}
```

## Mockito

### 基本模拟

```java
@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock
    private OrderRepository orderRepository;

    @Mock
    private PaymentService paymentService;

    @InjectMocks
    private OrderService orderService;

    @Test
    void shouldCreateOrder() {
        Order order = new Order("ORDER-1", 100.0);
        when(orderRepository.save(any())).thenReturn(order);

        Order result = orderService.createOrder(100.0);

        assertNotNull(result);
        verify(orderRepository).save(any());
    }
}
```

### 验证

```java
// 验证方法被调用
verify(repository).save(any(User.class));

// 验证方法被特定参数调用
verify(repository).findById(1L);

// 验证调用次数
verify(repository, times(2)).save(any());
verify(repository, never()).delete(any());
verify(repository, atLeastOnce()).findAll();

// 验证调用顺序
InOrder inOrder = inOrder(service, repository);
inOrder.verify(service).validate(any());
inOrder.verify(repository).save(any());

// 捕获参数
ArgumentCaptor<User> captor = ArgumentCaptor.forClass(User.class);
verify(repository).save(captor.capture());
assertEquals("Alice", captor.getValue().getName());
```

## 测试覆盖率

### Maven 配置

```xml
<plugin>
    <groupId>org.jacoco</groupId>
    <artifactId>jacoco-maven-plugin</artifactId>
    <version>0.8.11</version>
</plugin>
```

### 覆盖率目标

| 代码类型 | 目标 |
|---------|------|
| 核心业务逻辑 | 100% |
| 公共 API | 90%+ |
| 一般代码 | 80%+ |

## 最佳实践

**应该：**
- 先写测试（TDD）
- 使用描述性测试名称
- 每个测试只测一件事
- 使用 `@DisplayName` 提高可读性
- 模拟外部依赖
- 目标 80%+ 覆盖率

**不应该：**
- 直接测试私有方法
- 在测试中使用复杂条件
- 忽略测试失败
- 在测试间共享可变状态
