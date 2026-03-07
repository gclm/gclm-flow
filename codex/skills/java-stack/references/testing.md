# Java 测试模式

## JUnit 5 基础

```java
import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;

class UserServiceTest {

    @BeforeEach
    void setUp() {
        // 每个测试前的准备工作
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

## TDD 工作流

```
RED     → 先写失败的测试
GREEN   → 写最少代码通过测试
REFACTOR → 重构代码，保持测试通过
REPEAT  → 继续下一个需求
```

## Mockito

```java
@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    @Mock
    private OrderRepository orderRepository;

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

## Web 层测试

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
}
```

## 数据层测试

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

## 参数化测试

```java
@ParameterizedTest
@CsvSource({
    "1, 2, 3",
    "10, 20, 30",
    "-5, 5, 0"
})
void shouldAddNumbers(int a, int b, int expected) {
    assertEquals(expected, calculator.add(a, b));
}
```

## 覆盖率目标

| 代码类型 | 目标 |
|---------|------|
| 核心业务逻辑 | 100% |
| 公共 API | 90%+ |
| 一般代码 | 80%+ |
