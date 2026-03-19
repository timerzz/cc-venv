package webui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/env"
)

// MCPListResponse MCP列表响应
type MCPListResponse struct {
	Servers map[string]MCPServer `json:"servers"`
}

// AddMCPRequest 添加MCP服务器请求
type AddMCPRequest struct {
	Name   string    `json:"name" binding:"required"`
	Config MCPServer `json:"config" binding:"required"`
}

// UpdateMCPRequest 更新MCP服务器请求
type UpdateMCPRequest struct {
	Config MCPServer `json:"config" binding:"required"`
}

// claudeJSON .claude.json 结构
type claudeJSON struct {
	MCPServers map[string]MCPServer `json:"mcpServers,omitempty"`
}

// readClaudeJSON 读取 .claude.json 获取 MCP 配置
func readClaudeJSON(envDir string) (map[string]MCPServer, error) {
	path := filepath.Join(envDir, ".claude.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cfg claudeJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg.MCPServers, nil
}

// ListMCP 列出MCP服务器
func ListMCP(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	servers, err := readClaudeJSON(e.EnvDir)
	if err != nil {
		InternalError(c, fmt.Sprintf("read .claude.json: %v", err))
		return
	}

	if servers == nil {
		servers = make(map[string]MCPServer)
	}

	Success(c, MCPListResponse{Servers: servers})
}

// AddMCP 添加MCP服务器
func AddMCP(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	var req AddMCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	if req.Name == "" {
		BadRequest(c, "server name is required")
		return
	}

	// 读取现有配置
	servers, err := readClaudeJSON(e.EnvDir)
	if err != nil {
		InternalError(c, fmt.Sprintf("read .claude.json: %v", err))
		return
	}

	if servers == nil {
		servers = make(map[string]MCPServer)
	}

	// 检查是否已存在
	if _, exists := servers[req.Name]; exists {
		Error(c, fmt.Sprintf("MCP server %q already exists", req.Name))
		return
	}

	// 添加新服务器
	servers[req.Name] = req.Config

	// 写入配置
	if err := writeClaudeJSON(e.EnvDir, servers); err != nil {
		InternalError(c, fmt.Sprintf("write .claude.json: %v", err))
		return
	}

	Success(c, gin.H{
		"name":   req.Name,
		"config": req.Config,
	})
}

// UpdateMCP 更新MCP服务器
func UpdateMCP(c *gin.Context) {
	envName := c.Param("name")
	serverName := c.Param("server")

	e, err := env.Load(envName)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", envName))
		return
	}

	var req UpdateMCPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	// 读取现有配置
	servers, err := readClaudeJSON(e.EnvDir)
	if err != nil {
		InternalError(c, fmt.Sprintf("read .claude.json: %v", err))
		return
	}

	if servers == nil {
		servers = make(map[string]MCPServer)
	}

	// 检查是否存在
	if _, exists := servers[serverName]; !exists {
		NotFound(c, fmt.Sprintf("MCP server %q not found", serverName))
		return
	}

	// 更新服务器配置
	servers[serverName] = req.Config

	// 写入配置
	if err := writeClaudeJSON(e.EnvDir, servers); err != nil {
		InternalError(c, fmt.Sprintf("write .claude.json: %v", err))
		return
	}

	Success(c, gin.H{
		"name":   serverName,
		"config": req.Config,
	})
}

// DeleteMCP 删除MCP服务器
func DeleteMCP(c *gin.Context) {
	envName := c.Param("name")
	serverName := c.Param("server")

	e, err := env.Load(envName)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", envName))
		return
	}

	// 读取现有配置
	servers, err := readClaudeJSON(e.EnvDir)
	if err != nil {
		InternalError(c, fmt.Sprintf("read .claude.json: %v", err))
		return
	}

	if servers == nil {
		NotFound(c, fmt.Sprintf("MCP server %q not found", serverName))
		return
	}

	// 检查是否存在
	if _, exists := servers[serverName]; !exists {
		NotFound(c, fmt.Sprintf("MCP server %q not found", serverName))
		return
	}

	// 删除服务器
	delete(servers, serverName)

	// 写入配置
	if err := writeClaudeJSON(e.EnvDir, servers); err != nil {
		InternalError(c, fmt.Sprintf("write .claude.json: %v", err))
		return
	}

	Success(c, gin.H{"name": serverName})
}

// writeClaudeJSON 写入 .claude.json
func writeClaudeJSON(envDir string, servers map[string]MCPServer) error {
	path := filepath.Join(envDir, ".claude.json")

	cfg := claudeJSON{
		MCPServers: servers,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
