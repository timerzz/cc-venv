package webui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/config"
	"github.com/timerzz/cc-venv/internal/env"
	"github.com/timerzz/cc-venv/internal/scan"
)

// EnvListItem 环境列表项
type EnvListItem struct {
	Name      string         `json:"name"`
	Path      string         `json:"path"`
	Resources ResourceCounts `json:"resources"`
}

// ResourceCounts 资源计数
type ResourceCounts struct {
	Skills     int `json:"skills"`
	Agents     int `json:"agents"`
	Commands   int `json:"commands"`
	Rules      int `json:"rules"`
	Hooks      int `json:"hooks"`
	MCPServers int `json:"mcpServers"`
}

// EnvDetail 环境详情
type EnvDetail struct {
	Name        string                 `json:"name"`
	Path        string                 `json:"path"`
	ClaudeMd    string                 `json:"claudeMd,omitempty"`
	Settings    map[string]any         `json:"settings,omitempty"`
	EnvVars     config.EnvVars         `json:"envVars,omitempty"`
	MCPServers  map[string]MCPServer   `json:"mcpServers,omitempty"`
	Resources   ResourceLists          `json:"resources"`
}

// ResourceLists 资源列表
type ResourceLists struct {
	Skills   []string `json:"skills"`
	Agents   []string `json:"agents"`
	Commands []string `json:"commands"`
	Rules    []string `json:"rules"`
	Hooks    []string `json:"hooks"`
	Plugins  []string `json:"plugins"`
}

// MCPServer MCP服务器配置
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// CreateEnvRequest 创建环境请求
type CreateEnvRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateEnvRequest 更新环境请求
type UpdateEnvRequest struct {
	Name     string           `json:"name,omitempty"`
	ClaudeMd string           `json:"claudeMd,omitempty"`
	EnvVars  config.EnvVars   `json:"envVars,omitempty"`
}

// ListEnvs 列出所有环境
func ListEnvs(c *gin.Context) {
	envs, err := env.List()
	if err != nil {
		InternalError(c, fmt.Sprintf("list environments: %v", err))
		return
	}

	items := make([]EnvListItem, 0, len(envs))
	for _, e := range envs {
		item := EnvListItem{
			Name: e.Name,
			Path: e.RootPath,
		}

		// 读取 ccv.json 获取资源摘要
		cfg, err := readCcvJSON(e.ManifestPath)
		if err == nil {
			item.Resources = ResourceCounts{
				Skills:     len(cfg.Resources.Skills),
				Agents:     len(cfg.Resources.Agents),
				Commands:   len(cfg.Resources.Commands),
				Rules:      len(cfg.Resources.Rules),
				Hooks:      len(cfg.Resources.Hooks),
				MCPServers: len(cfg.Resources.MCPServer),
			}
		}

		items = append(items, item)
	}

	Success(c, gin.H{"envs": items})
}

// GetEnv 获取环境详情
func GetEnv(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	detail := EnvDetail{
		Name: e.Name,
		Path: e.RootPath,
	}

	// 读取 CLAUDE.md
	claudeMdPath := filepath.Join(e.EnvDir, ".claude", "CLAUDE.md")
	if data, err := os.ReadFile(claudeMdPath); err == nil {
		detail.ClaudeMd = string(data)
	}

	// 读取 settings.json 中的 env
	envVars, err := config.ReadSettingsJSONEnv(e.SettingsPath)
	if err == nil {
		detail.EnvVars = envVars
	}

	// 读取完整的 settings.json
	settingsPath := e.SettingsPath // 这是 .claude/settings.json
	if data, err := os.ReadFile(settingsPath); err == nil {
		var settings map[string]any
		if json.Unmarshal(data, &settings) == nil {
			detail.Settings = settings
		}
	}

	// 使用统一的 scanner 模块扫描资源
	resources, err := scan.ScanEnvironment(e.EnvDir, scan.Options{ScanMCP: true})
	if err != nil {
		InternalError(c, fmt.Sprintf("scan environment resources: %v", err))
		return
	}

	detail.Resources = ResourceLists{
		Skills:   resources.Skills,
		Agents:   resources.Agents,
		Commands: resources.Commands,
		Rules:    resources.Rules,
		Hooks:    resources.Hooks,
		Plugins:  resources.Plugins,
	}

	// 转换 MCP servers 到响应格式
	if len(resources.MCPServers) > 0 {
		detail.MCPServers = make(map[string]MCPServer)
		for name, cfg := range resources.MCPServers {
			detail.MCPServers[name] = MCPServer{
				Command: cfg.Command,
				Args:    cfg.Args,
				Env:     cfg.Env,
			}
		}
	}

	Success(c, detail)
}

// CreateEnv 创建环境
func CreateEnv(c *gin.Context) {
	var req CreateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "name is required")
		return
	}

	e, err := env.Create(req.Name)
	if err != nil {
		Error(c, fmt.Sprintf("create environment: %v", err))
		return
	}

	Success(c, gin.H{
		"name": e.Name,
		"path": e.RootPath,
	})
}

// UpdateEnv 更新环境
func UpdateEnv(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	var req UpdateEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	// 更新 CLAUDE.md
	if req.ClaudeMd != "" {
		claudeMdPath := filepath.Join(e.EnvDir, ".claude", "CLAUDE.md")
		if err := os.WriteFile(claudeMdPath, []byte(req.ClaudeMd), 0o644); err != nil {
			InternalError(c, fmt.Sprintf("write CLAUDE.md: %v", err))
			return
		}
	}

	// 更新环境变量到 settings.json
	if req.EnvVars != nil {
		if err := config.WriteSettingsJSONEnv(e.SettingsPath, req.EnvVars); err != nil {
			InternalError(c, fmt.Sprintf("write settings.json: %v", err))
			return
		}
	}

	// 重命名环境
	if req.Name != "" && req.Name != name {
		// 校验新名称
		if !isValidEnvName(req.Name) {
			BadRequest(c, "invalid environment name: only letters, numbers, underscore and hyphen are allowed")
			return
		}

		// 检查新名称是否已存在
		_, err := env.Load(req.Name)
		if err == nil {
			BadRequest(c, fmt.Sprintf("environment %q already exists", req.Name))
			return
		}

		// 执行重命名
		oldPath := e.RootPath
		newPath := filepath.Join(filepath.Dir(oldPath), req.Name)
		if err := os.Rename(oldPath, newPath); err != nil {
			InternalError(c, fmt.Sprintf("rename environment: %v", err))
			return
		}

		// 更新 ccv.json 中的 name（这是 rename 的正式步骤）
		newManifestPath := filepath.Join(newPath, "ccv.json")
		cfg, err := readCcvJSON(newManifestPath)
		if err != nil {
			InternalError(c, fmt.Sprintf("read ccv.json after rename: %v", err))
			return
		}

		cfg.Name = req.Name
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			InternalError(c, fmt.Sprintf("marshal ccv.json: %v", err))
			return
		}
		data = append(data, '\n')
		if err := os.WriteFile(newManifestPath, data, 0o644); err != nil {
			InternalError(c, fmt.Sprintf("write ccv.json: %v", err))
			return
		}

		// 返回新名称
		Success(c, gin.H{"name": req.Name, "renamed": true})
		return
	}

	Success(c, gin.H{"name": e.Name})
}

// DeleteEnv 删除环境
func DeleteEnv(c *gin.Context) {
	name := c.Param("name")

	// force=true 跳过确认
	if err := env.Remove(name, true); err != nil {
		Error(c, fmt.Sprintf("delete environment: %v", err))
		return
	}

	Success(c, gin.H{"name": name})
}

// readCcvJSON 读取 ccv.json
func readCcvJSON(path string) (*config.CcvConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg config.CcvConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// isValidEnvName 检查环境名称是否合法
func isValidEnvName(name string) bool {
	if name == "" || len(name) > 64 {
		return false
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
}
