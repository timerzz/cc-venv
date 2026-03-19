package webui

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/config"
	"github.com/timerzz/cc-venv/internal/env"
)

// LLMConfig LLM配置
type LLMConfig struct {
	APIKey   string          `json:"apiKey"`
	BaseURL  string          `json:"baseUrl"`
	Models   LLMModelsConfig `json:"models"`
}

// LLMModelsConfig 模型配置
type LLMModelsConfig struct {
	Default string `json:"default"`
	Sonnet  string `json:"sonnet,omitempty"`
	Opus    string `json:"opus,omitempty"`
	Haiku   string `json:"haiku,omitempty"`
}

// UpdateLLMConfigRequest 更新LLM配置请求
type UpdateLLMConfigRequest struct {
	APIKey   string          `json:"apiKey"`
	BaseURL  string          `json:"baseUrl"`
	Models   LLMModelsConfig `json:"models"`
}

// LLMProvider LLM供应商
type LLMProvider struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	BaseURL string   `json:"baseUrl"`
	Models  []string `json:"models"`
}

// 预定义的LLM供应商列表
// 注意：当前实现使用 Anthropic-compatible 配置格式
var llmProviders = []LLMProvider{
	{
		ID:      "anthropic",
		Name:    "Anthropic",
		BaseURL: "https://api.anthropic.com",
		Models: []string{
			"claude-sonnet-4-6",
			"claude-opus-4-6",
			"claude-haiku-4-5",
		},
	},
}

// GetLLMConfig 获取LLM配置
func GetLLMConfig(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	cfg := LLMConfig{
		Models: LLMModelsConfig{},
	}

	// 从 settings.json 中的 env 字段读取配置
	// 使用 Claude Code 原生支持的环境变量
	envVars, err := config.ReadSettingsJSONEnv(e.SettingsPath)
	if err == nil {
		// API Key 脱敏显示
		if apiKey := envVars["ANTHROPIC_API_KEY"]; apiKey != "" {
			cfg.APIKey = maskAPIKey(apiKey)
		}

		cfg.BaseURL = envVars["ANTHROPIC_BASE_URL"]
		cfg.Models.Default = envVars["ANTHROPIC_SMALL_FAST_MODEL"]
		cfg.Models.Sonnet = envVars["ANTHROPIC_DEFAULT_SONNET_MODEL"]
		cfg.Models.Opus = envVars["ANTHROPIC_DEFAULT_OPUS_MODEL"]
		cfg.Models.Haiku = envVars["ANTHROPIC_DEFAULT_HAIKU_MODEL"]
	}

	Success(c, cfg)
}

// UpdateLLMConfig 更新LLM配置
func UpdateLLMConfig(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	var req UpdateLLMConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	// 读取现有环境变量
	envVars, err := config.ReadSettingsJSONEnv(e.SettingsPath)
	if err != nil {
		envVars = make(map[string]string)
	}

	// 更新配置
	// 使用 Claude Code 原生支持的环境变量
	if req.APIKey != "" {
		envVars["ANTHROPIC_API_KEY"] = req.APIKey
	}
	if req.BaseURL != "" {
		envVars["ANTHROPIC_BASE_URL"] = req.BaseURL
	}
	if req.Models.Default != "" {
		envVars["ANTHROPIC_SMALL_FAST_MODEL"] = req.Models.Default
	}
	if req.Models.Sonnet != "" {
		envVars["ANTHROPIC_DEFAULT_SONNET_MODEL"] = req.Models.Sonnet
	}
	if req.Models.Opus != "" {
		envVars["ANTHROPIC_DEFAULT_OPUS_MODEL"] = req.Models.Opus
	}
	if req.Models.Haiku != "" {
		envVars["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = req.Models.Haiku
	}

	// 写入 settings.json
	if err := config.WriteSettingsJSONEnv(e.SettingsPath, envVars); err != nil {
		InternalError(c, fmt.Sprintf("write settings.json: %v", err))
		return
	}

	Success(c, gin.H{"name": name})
}

// ListLLMProviders 列出支持的LLM供应商
func ListLLMProviders(c *gin.Context) {
	Success(c, gin.H{"providers": llmProviders})
}

// maskAPIKey 脱敏显示 API Key
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}
