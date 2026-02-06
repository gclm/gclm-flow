package api

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gclm/gclm-flow/gclm-engine/internal/api/websocket"
	"github.com/gclm/gclm-flow/gclm-engine/internal/domain"
	"github.com/gclm/gclm-flow/gclm-engine/internal/logger"
)

// Server HTTP API 服务器
type Server struct {
	router          *gin.Engine
	addr            string
	taskSvc         domain.TaskService
	workflowSvc     domain.WorkflowService
	wsHub           *websocket.Hub
	webFS           fs.FS
	shutdownTimeout int
	mu              sync.RWMutex
}

// NewServer 创建新的 HTTP 服务器
func NewServer(
	addr string,
	taskSvc domain.TaskService,
	workflowSvc domain.WorkflowService,
	wsHub *websocket.Hub,
	webFS fs.FS,
) *Server {
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		addr:            addr,
		taskSvc:         taskSvc,
		workflowSvc:     workflowSvc,
		wsHub:           wsHub,
		webFS:           webFS,
		shutdownTimeout: 10,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware 设置中间件
func (s *Server) setupMiddleware() {
	s.router = gin.New()

	// 恢复中间件
	s.router.Use(gin.Recovery())

	// CORS 中间件
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// 日志中间件
	s.router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 日志记录
		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Msg("HTTP request")
	})
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// API 路由
	api := s.router.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", s.listTasks)
			tasks.POST("", s.createTask)
			tasks.GET("/:id", s.getTask)
			tasks.GET("/:id/phases", s.getTaskPhases)
			tasks.GET("/:id/events", s.getTaskEvents)
			tasks.POST("/:id/pause", s.pauseTask)
			tasks.POST("/:id/resume", s.resumeTask)
			tasks.POST("/:id/cancel", s.cancelTask)
		}

		phases := api.Group("/phases")
		{
			phases.POST("/:id/complete", s.completePhase)
			phases.POST("/:id/fail", s.failPhase)
		}

		workflows := api.Group("/workflows")
		{
			workflows.GET("", s.listWorkflows)
			workflows.GET("/:name", s.getWorkflow)
			workflows.GET("/type/:type", s.getWorkflowByType)
			workflows.GET("/:name/yaml", s.getWorkflowYAML)
		}
	}

	// WebSocket 路由
	s.router.GET("/ws/tasks/:id", s.wsHub.HandleWebSocket)

	// 静态文件 - 使用 embed 文件系统
	if s.webFS != nil {
		// 获取 static 子文件系统
		staticFS, err := fs.Sub(s.webFS, "static")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create static filesystem")
		} else {
			// 使用 StaticFS 服务静态文件，需要转换为 http.FileSystem
			s.router.StaticFS("/static", http.FS(staticFS))
		}
	}
	s.router.GET("/", s.indexHandler)
}

// Start 启动 HTTP 服务器
func (s *Server) Start(ctx context.Context) error {
	logger.Info().Str("addr", s.addr).Msg("Starting HTTP server")

	// 启动 WebSocket Hub
	if s.wsHub != nil {
		go s.wsHub.Run(ctx)
	}

	// 在 goroutine 中启动服务器
	serverErr := make(chan error, 1)
	go func() {
		logger.Info().Str("addr", s.addr).Msg("HTTP server listening")
		if err := s.router.Run(s.addr); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// 等待上下文取消或服务器错误
	select {
	case <-ctx.Done():
		logger.Info().Msg("Shutting down HTTP server")
		return s.shutdown()
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}
}

// shutdown 优雅关闭服务器
func (s *Server) shutdown() error {
	logger.Info().Msg("HTTP server shutdown complete")
	return nil
}

// indexHandler 首页处理器
func (s *Server) indexHandler(c *gin.Context) {
	if s.webFS != nil {
		// 从 embed 文件系统读取 index.html
		file, err := s.webFS.Open("index.html")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to open index.html from embed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load index.html"})
			return
		}
		defer file.Close()

		// 读取文件内容
		content, err := io.ReadAll(file)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to read index.html")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read index.html"})
			return
		}

		// 设置正确的 Content-Type
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
		return
	}

	// Fallback: 返回错误
	c.JSON(http.StatusNotFound, gin.H{"error": "Web UI not available (web files not embedded)"})
}

// ============================================================================
// API Handlers
// ============================================================================
// All handlers are implemented in handlers.go
