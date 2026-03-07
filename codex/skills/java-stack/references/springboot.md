# Spring Boot 最佳实践

## 配置管理

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

## JPA 实体

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

    @CreatedDate
    @Column(updatable = false)
    private LocalDateTime createdAt;

    @LastModifiedDate
    private LocalDateTime updatedAt;

    @Version
    private Long version;
}
```

## Repository 模式

```java
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);

    @Query("SELECT u FROM User u WHERE u.status = :status")
    Page<User> findByStatus(@Param("status") UserStatus status, Pageable pageable);
}
```

## 事件驱动

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

## 缓存配置

```java
@Configuration
@EnableCaching
public class CacheConfig {

    @Bean
    public CacheManager cacheManager(RedisConnectionFactory factory) {
        RedisCacheConfiguration config = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(30))
            .serializeValuesWith(RedisSerializationContext.SerializationPair
                .fromSerializer(new GenericJackson2JsonRedisSerializer()));

        return RedisCacheManager.builder(factory)
            .cacheDefaults(config)
            .build();
    }
}
```

## 异步处理

```java
@Configuration
@EnableAsync
public class AsyncConfig {

    @Bean
    public Executor taskExecutor() {
        ThreadPoolTaskExecutor executor = new ThreadPoolTaskExecutor();
        executor.setCorePoolSize(5);
        executor.setMaxPoolSize(10);
        executor.setQueueCapacity(100);
        executor.setThreadNamePrefix("async-");
        executor.initialize();
        return executor;
    }
}
```
