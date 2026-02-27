---
name: go-stack
description: |
  Go/Gin 技术栈完整开发指南。当检测到 Go 项目（go.mod）
  或用户明确要求 Go/Gin/Echo 开发时自动触发。包含：
  (1) 项目结构规范 (2) Gin 最佳实践 (3) 测试模式 (4) 错误处理
metadata:
  author: gclm-flow
  version: "2.0.0"
  platforms:
    - claude-code
    - codex-cli
  tags:
    - go
    - gin
    - echo
---

# Go 技术栈开发指南

## 框架检测

- 存在 `github.com/gin-gonic/gin` → Gin，详见 [gin.md](references/gin.md)
- 存在 `github.com/labstack/echo` → Echo
- 测试相关 → 详见 [testing.md](references/testing.md)

## 标准项目结构

```
myapp/
├── cmd/
│   └── server/
│       └── main.go          # 应用入口
├── internal/
│   ├── handler/             # HTTP 处理器
│   ├── service/             # 业务逻辑
│   ├── repository/          # 数据访问
│   ├── model/               # 数据模型
│   ├── middleware/          # 中间件
│   └── config/              # 配置
├── pkg/                     # 公共库
│   ├── response/
│   └── utils/
├── go.mod
└── go.sum
```

## 核心规范

### 统一响应格式

```go
package response

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func Success(data interface{}) Response {
    return Response{Success: true, Data: data}
}

func Error(msg string) Response {
    return Response{Success: false, Error: msg}
}
```

### 处理器结构

```go
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
        users.DELETE("/:id", h.Delete)
    }
}
```

### 错误处理

```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

func handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, ErrUserNotFound):
        c.JSON(http.StatusNotFound, response.Error(err.Error()))
    case errors.Is(err, ErrInvalidInput):
        c.JSON(http.StatusBadRequest, response.Error(err.Error()))
    default:
        c.JSON(http.StatusInternalServerError, response.Error("internal error"))
    }
}
```

### 中间件

```go
// 认证中间件
func Auth(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
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
        c.Next()
        log.Printf("[%s] %s %d %v",
            c.Request.Method, c.Request.URL.Path,
            c.Writer.Status(), time.Since(start))
    }
}
```

### 仓储模式

```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return &user, err
}
```

## 测试规范

详见 [testing.md](references/testing.md)

- 使用内置 testing 包
- 目标覆盖率：80%+

## 相关技能

- `code-review` - Go 代码审查
- `testing` - Go 测试模式
- `database` - GORM 模式
