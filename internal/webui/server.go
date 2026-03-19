package webui

import (
	"embed"
	"fmt"
	"io/fs"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed all:static
var staticFS embed.FS

// Config Web服务配置
type Config struct {
	Port    int
	DevMode bool
	NoOpen  bool
}

// Server Web服务
type Server struct {
	config Config
	engine *gin.Engine
}

// NewServer 创建Web服务
func NewServer(cfg Config) *Server {
	if cfg.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	// 添加 Recovery 中间件，确保 panic 时能够恢复
	engine.Use(gin.Recovery())
	// 开发模式启用请求日志
	if cfg.DevMode {
		engine.Use(gin.Logger())
	}

	return &Server{
		config: cfg,
		engine: engine,
	}
}

// Start 启动Web服务
func (s *Server) Start() error {
	// 注册路由
	s.setupRoutes()

	// 启动HTTP服务
	addr := fmt.Sprintf(":%d", s.config.Port)
	fmt.Printf("ccv web server running at http://localhost%s\n", addr)

	return s.engine.Run(addr)
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 注册API路由
	api := s.engine.Group("/api")
	registerAPIRoutes(api)

	// 静态资源服务（非开发模式）
	if !s.config.DevMode {
		staticDir, err := fs.Sub(staticFS, "static")
		if err != nil {
			fmt.Printf("warning: failed to load static files: %v\n", err)
			return
		}
		setupStaticRoutes(s.engine, staticDir)
	}
}

// OpenBrowser 打开浏览器
func OpenBrowser(url string) error {
	// 延迟一下等待服务启动
	time.Sleep(500 * time.Millisecond)

	// 使用系统默认浏览器打开URL
	// 这里简化实现，实际可能需要根据平台区分
	fmt.Printf("Please open your browser and visit: %s\n", url)
	return nil
}

// Open 启动Web服务（兼容旧接口）
func Open() error {
	server := NewServer(Config{Port: 3000})
	return server.Start()
}
