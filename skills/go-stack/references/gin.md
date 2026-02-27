# Gin 最佳实践

## 路由组织

```go
func main() {
    r := gin.Default()

    r.GET("/health", healthCheck)

    v1 := r.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("", listUsers)
            users.GET("/:id", getUser)
            users.POST("", createUser)
            users.DELETE("/:id", deleteUser)
        }
    }

    r.Run(":8080")
}
```

## 依赖注入

```go
type App struct {
    db      *gorm.DB
    router  *gin.Engine
    handler *handler.UserHandler
}

func NewApp(cfg *config.Config) (*App, error) {
    db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    router := gin.Default()
    userHandler.RegisterRoutes(router.Group("/api/v1"))

    return &App{db: db, router: router, handler: userHandler}, nil
}
```

## 配置管理

```go
import "github.com/spf13/viper"

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    var config Config
    viper.Unmarshal(&config)
    return &config, nil
}
```
