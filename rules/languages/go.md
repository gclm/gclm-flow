# Go/Gin 规则

Go 和 Gin 框架的编码规范和最佳实践。

## Go 编码规范

### 1. 命名规范
```go
// 包名：小写单词
package userservice

// 导出类型：PascalCase
type UserService struct {}

// 私有类型：camelCase
type userService struct {}

// 接口：通常以 er 结尾
type UserReader interface {
    Read(id int64) (*User, error)
}

// 常量：PascalCase（导出）或 camelCase（私有）
const (
    MaxRetryCount = 3
    defaultTimeout = 30 * time.Second
)

// 枚举常量：类型前缀
type UserStatus int
const (
    UserStatusActive UserStatus = iota
    UserStatusInactive
    UserStatusDeleted
)
```

### 2. 项目结构
```
myapp/
├── cmd/                 # 应用入口
│   └── server/
│       └── main.go
├── internal/            # 私有代码
│   ├── handler/         # HTTP 处理器
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问
│   └── model/           # 数据模型
├── pkg/                 # 公共库
│   └── utils/
├── api/                 # API 定义
│   └── openapi.yaml
├── config/              # 配置文件
├── go.mod
└── go.sum
```

### 3. 错误处理
```go
import "errors"

// 自定义错误
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

// 包装错误
func (s *UserService) GetByID(id int64) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
    }
    return user, nil
}

// 处理错误
func (h *Handler) GetUser(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse("invalid id"))
        return
    }

    user, err := h.service.GetByID(id)
    if errors.Is(err, ErrUserNotFound) {
        c.JSON(http.StatusNotFound, ErrorResponse(err.Error()))
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse(err.Error()))
        return
    }

    c.JSON(http.StatusOK, SuccessResponse(user))
}
```

### 4. 并发模式
```go
// 使用 context 控制超时
func (s *Service) Process(ctx context.Context, items []Item) error {
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()

    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := s.processItem(ctx, item); err != nil {
                return err
            }
        }
    }
    return nil
}

// 使用 errgroup 并发处理
func (s *Service) ProcessAll(ctx context.Context, items []Item) error {
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item // 捕获变量
        g.Go(func() error {
            return s.processItem(ctx, item)
        })
    }

    return g.Wait()
}
```

## Gin 最佳实践

### 1. 路由组织
```go
package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()

    // 健康检查
    r.GET("/health", healthCheck)

    // API v1
    v1 := r.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("", listUsers)
            users.GET("/:id", getUser)
            users.POST("", createUser)
            users.PUT("/:id", updateUser)
            users.DELETE("/:id", deleteUser)
        }
    }

    r.Run(":8080")
}
```

### 2. 处理器结构
```go
package handler

import "github.com/gin-gonic/gin"

type UserHandler struct {
    service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
    return &UserHandler{service: s}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
    users := r.Group("/users")
    {
        users.GET("", h.List)
        users.GET("/:id", h.Get)
        users.POST("", h.Create)
        users.PUT("/:id", h.Update)
        users.DELETE("/:id", h.Delete)
    }
}

func (h *UserHandler) Get(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse("invalid id"))
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, SuccessResponse(user))
}
```

### 3. 统一响应格式
```go
package response

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Message string      `json:"message,omitempty"`
}

func Success(data interface{}) Response {
    return Response{
        Success: true,
        Data:    data,
    }
}

func Error(msg string) Response {
    return Response{
        Success: false,
        Error:   msg,
    }
}

func ErrorWithCode(code int, msg string) Response {
    return Response{
        Success: false,
        Error:   msg,
    }
}
```

### 4. 中间件
```go
// 认证中间件
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
            c.Abort()
            return
        }

        claims, err := validateToken(token, jwtSecret)
        if err != nil {
            c.JSON(http.StatusUnauthorized, response.Error("invalid token"))
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Next()
    }
}

// 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        latency := time.Since(start)
        log.Printf("[%s] %s %d %v",
            c.Request.Method,
            path,
            c.Writer.Status(),
            latency,
        )
    }
}

// 错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                c.JSON(http.StatusInternalServerError, response.Error("internal server error"))
            }
        }()
        c.Next()
    }
}
```

### 5. 请求验证
```go
type CreateUserRequest struct {
    Email string `json:"email" binding:"required,email"`
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
}

func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error(err.Error()))
        return
    }

    user := &model.User{
        Email: req.Email,
        Name:  req.Name,
        Age:   req.Age,
    }

    if err := h.service.Create(c.Request.Context(), user); err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, response.Success(user))
}
```

### 6. 依赖注入
```go
package main

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type App struct {
    db      *gorm.DB
    router  *gin.Engine
    handler *handler.UserHandler
}

func NewApp(db *gorm.DB) *App {
    // 创建依赖链
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    router := gin.Default()
    userHandler.RegisterRoutes(router.Group("/api/v1"))

    return &App{
        db:      db,
        router:  router,
        handler: userHandler,
    }
}
```

## 测试规范

### 1. 单元测试
```go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(id int64) (*model.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.User), args.Error(1)
}

func TestUserService_GetByID(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    expectedUser := &model.User{ID: 1, Email: "test@example.com"}
    mockRepo.On("FindByID", int64(1)).Return(expectedUser, nil)

    user, err := service.GetByID(context.Background(), 1)

    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}
```

### 2. 集成测试
```go
func TestUserAPI(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    app := NewApp(db)
    router := app.router

    // 测试创建用户
    body := `{"email":"test@example.com","name":"Test User"}`
    req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## 性能优化

### 1. 数据库连接池
```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}

sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. 缓存
```go
import "github.com/patrickmn/go-cache"

type CachedUserService struct {
    service *UserService
    cache   *cache.Cache
}

func (s *CachedUserService) GetByID(ctx context.Context, id int64) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    if cached, found := s.cache.Get(cacheKey); found {
        return cached.(*User), nil
    }

    user, err := s.service.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    s.cache.Set(cacheKey, user, cache.DefaultExpiration)
    return user, nil
}
```
