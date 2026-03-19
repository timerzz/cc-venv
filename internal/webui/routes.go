package webui

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerAPIRoutes 注册API路由
func registerAPIRoutes(api *gin.RouterGroup) {
	// 环境管理
	api.GET("/envs", ListEnvs)
	api.POST("/envs", CreateEnv)
	api.GET("/envs/:name", GetEnv)
	api.PUT("/envs/:name", UpdateEnv)
	api.DELETE("/envs/:name", DeleteEnv)
	api.POST("/envs/:name/export", ExportEnv)
	api.POST("/envs/import", ImportEnv)

	// LLM配置
	api.GET("/envs/:name/llm", GetLLMConfig)
	api.PUT("/envs/:name/llm", UpdateLLMConfig)
	api.GET("/llm/providers", ListLLMProviders)

	// MCP管理
	api.GET("/envs/:name/mcp", ListMCP)
	api.POST("/envs/:name/mcp", AddMCP)
	api.PUT("/envs/:name/mcp/:server", UpdateMCP)
	api.DELETE("/envs/:name/mcp/:server", DeleteMCP)

	// Skills管理
	api.GET("/envs/:name/skills", ListSkills)
	api.POST("/envs/:name/skills", AddSkill)
	api.DELETE("/envs/:name/skills/:skill", DeleteSkill)

	// 文件下载
	api.GET("/downloads/:filename", DownloadFile)
}

// setupStaticRoutes 设置静态资源路由
func setupStaticRoutes(engine *gin.Engine, staticDir fs.FS) {
	// 静态资源路由
	engine.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(staticDir))
	})

	// SPA fallback - 所有未匹配的路由返回index.html
	engine.NoRoute(func(c *gin.Context) {
		// 如果是API请求但未匹配，返回404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			NotFound(c, "not found")
			return
		}

		// 返回index.html
		c.FileFromFS("index.html", http.FS(staticDir))
	})
}
