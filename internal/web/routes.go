package web

import (
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/web/handlers"
)

// registerAPIRoutes 注册API路由
func registerAPIRoutes(api *gin.RouterGroup) {
	// 环境管理
	api.GET("/envs", handlers.ListEnvs)
	api.POST("/envs", handlers.CreateEnv)
	api.GET("/envs/:name", handlers.GetEnv)
	api.PUT("/envs/:name", handlers.UpdateEnv)
	api.DELETE("/envs/:name", handlers.DeleteEnv)
	api.POST("/envs/:name/export", handlers.ExportEnv)
	api.POST("/envs/import", handlers.ImportEnv)

	// LLM配置
	api.GET("/envs/:name/llm", handlers.GetLLMConfig)
	api.PUT("/envs/:name/llm", handlers.UpdateLLMConfig)
	api.GET("/llm/providers", handlers.ListLLMProviders)

	// MCP管理
	api.GET("/envs/:name/mcp", handlers.ListMCP)
	api.POST("/envs/:name/mcp", handlers.AddMCP)
	api.PUT("/envs/:name/mcp/:server", handlers.UpdateMCP)
	api.DELETE("/envs/:name/mcp/:server", handlers.DeleteMCP)

	// Skills管理
	api.GET("/envs/:name/skills", handlers.ListSkills)
	api.POST("/envs/:name/skills", handlers.AddSkill)
	api.DELETE("/envs/:name/skills/:skill", handlers.DeleteSkill)

	// 文件类资源管理
	api.GET("/envs/:name/resources/:kind", handlers.ListResourceFiles)
	api.GET("/envs/:name/resources/:kind/content", handlers.GetResourceFile)
	api.PUT("/envs/:name/resources/:kind/content", handlers.UpsertResourceFile)
	api.DELETE("/envs/:name/resources/:kind/content", handlers.DeleteResourceFile)

	// 文件下载
	api.GET("/downloads/:filename", handlers.DownloadFile)
}

// setupStaticRoutes 设置静态资源路由
func setupStaticRoutes(engine *gin.Engine, staticDir fs.FS) {
	indexHTML, err := fs.ReadFile(staticDir, "index.html")
	if err != nil {
		fmt.Printf("warning: failed to read embedded index.html: %v\n", err)
		return
	}

	assetsDir, err := fs.Sub(staticDir, "assets")
	if err != nil {
		fmt.Printf("warning: failed to load embedded assets: %v\n", err)
		return
	}

	serveIndex := func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	}

	// SPA 入口
	engine.GET("/", serveIndex)

	// 静态资源路由
	engine.GET("/assets/*filepath", func(c *gin.Context) {
		filepath := strings.TrimPrefix(c.Param("filepath"), "/")
		if filepath == "" {
			handlers.NotFound(c, "not found")
			return
		}

		data, err := fs.ReadFile(assetsDir, filepath)
		if err != nil {
			handlers.NotFound(c, "not found")
			return
		}

		contentType := mime.TypeByExtension(path.Ext(filepath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Data(http.StatusOK, contentType, data)
	})

	// SPA fallback - 所有未匹配的路由返回index.html
	engine.NoRoute(func(c *gin.Context) {
		// 如果是API请求但未匹配，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			handlers.NotFound(c, "not found")
			return
		}

		// 返回index.html
		serveIndex(c)
	})
}
