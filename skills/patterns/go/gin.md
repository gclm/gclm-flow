# Golang Patterns

Go/Gin 开发专用模式和最佳实践。

## 技能描述

这个技能包含 Go 后端开发的专用模式，适用于 Gin 框架。

## 包含的模式

### 1. 项目结构

```
myapp/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── handler/             # HTTP 处理器
│   │   ├── handler.go
│   │   └── user_handler.go
│   ├── service/             # 业务逻辑
│   │   └── user_service.go
│   ├── repository/          # 数据访问
│   │   └── user_repo.go
│   ├── model/               # 数据模型
│   │   └── user.go
│   ├── middleware/          # 中间件
│   │   ├── auth.go
│   │   └── logger.go
│   └── config/              # 配置
│       └── config.go
├── pkg/                     # 公共库
│   ├── response/
│   └── utils/
├── api/
│   └── openapi.yaml
├── go.mod
└── go.sum
```

### 2. 统一响应格式

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
```

### 3. 处理器结构

```go
package handler

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
        c.JSON(http.StatusBadRequest, response.Error("invalid id"))
        return
    }

    user, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, response.Success(user))
}
```

### 4. 错误处理

```go
package errors

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
)

// 包装错误
func Wrap(err error, message string) error {
    return fmt.Errorf("%s: %w", message, err)
}

// 处理器中的错误处理
func handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, ErrUserNotFound):
        c.JSON(http.StatusNotFound, response.Error(err.Error()))
    case errors.Is(err, ErrInvalidInput):
        c.JSON(http.StatusBadRequest, response.Error(err.Error()))
    case errors.Is(err, ErrUnauthorized):
        c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
    default:
        c.JSON(http.StatusInternalServerError, response.Error("internal server error"))
    }
}
```

### 5. 中间件

```go
package middleware

// 认证中间件
func Auth(jwtSecret string) gin.HandlerFunc {
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
func Logger() gin.HandlerFunc {
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

// CORS 中间件
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

### 6. 仓储模式

```go
package repository

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}
```

### 7. 配置管理

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

type ServerConfig struct {
    Port int
    Mode string
}

type DatabaseConfig struct {
    URL string
}

type JWTConfig struct {
    Secret     string
    Expiration time.Duration
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

### 8. 依赖注入

```go
package main

type App struct {
    db      *gorm.DB
    router  *gin.Engine
    handler *handler.UserHandler
}

func NewApp(cfg *config.Config) (*App, error) {
    // 初始化数据库
    db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // 创建依赖链
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 创建路由
    router := gin.Default()
    userHandler.RegisterRoutes(router.Group("/api/v1"))

    return &App{
        db:      db,
        router:  router,
        handler: userHandler,
    }, nil
}
```

## 使用场景

- 创建新的 Go 项目
- Gin 框架开发
- 中间件设计
- 错误处理
