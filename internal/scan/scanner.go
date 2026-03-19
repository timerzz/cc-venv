package scan

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MCPServerConfig MCP 服务器配置
type MCPServerConfig struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// Resources 扫描到的资源
type Resources struct {
	Skills     []string                   `json:"skills"`
	Agents     []string                   `json:"agents"`
	Commands   []string                   `json:"commands"`
	Rules      []string                   `json:"rules"`
	Hooks      []string                   `json:"hooks"`
	Plugins    []string                   `json:"plugins"`
	MCPServers map[string]MCPServerConfig `json:"mcpServers,omitempty"`
}

// Options 扫描选项
type Options struct {
	ScanMCP bool // 是否扫描 MCP servers
}

// ScanEnvironment 扫描环境目录
func ScanEnvironment(envDir string, opts Options) (*Resources, error) {
	resources := &Resources{
		Skills:     []string{},
		Agents:     []string{},
		Commands:   []string{},
		Rules:      []string{},
		Hooks:      []string{},
		Plugins:    []string{},
		MCPServers: map[string]MCPServerConfig{},
	}

	claudeDir := filepath.Join(envDir, ".claude")

	// 扫描 skills（目录）
	scanDir(claudeDir, "skills", func(name string, isDir bool) {
		if isDir {
			resources.Skills = append(resources.Skills, name)
		}
	})

	// 扫描 agents（.md 文件）
	scanDir(claudeDir, "agents", func(name string, isDir bool) {
		if !isDir && strings.HasSuffix(name, ".md") {
			resources.Agents = append(resources.Agents, strings.TrimSuffix(name, ".md"))
		}
	})

	// 扫描 commands（.md 文件）
	scanDir(claudeDir, "commands", func(name string, isDir bool) {
		if !isDir && strings.HasSuffix(name, ".md") {
			resources.Commands = append(resources.Commands, strings.TrimSuffix(name, ".md"))
		}
	})

	// 递归扫描 rules（.md 文件）
	scanRulesRecursive(claudeDir, resources)

	// 解析 hooks.json
	if err := scanHooksJSON(claudeDir, resources); err != nil {
		return nil, err
	}

	// 扫描 plugins/cache（忽略 temp*）
	scanDir(claudeDir, filepath.Join("plugins", "cache"), func(name string, isDir bool) {
		if isDir && !strings.HasPrefix(name, "temp") {
			resources.Plugins = append(resources.Plugins, name)
		}
	})

	// 扫描 MCP servers
	if opts.ScanMCP {
		if err := scanMCPServers(envDir, resources); err != nil {
			return nil, err
		}
	}

	// 统一排序
	sortResources(resources)

	return resources, nil
}

// scanDir 扫描指定目录
func scanDir(claudeDir, subDir string, fn func(name string, isDir bool)) {
	dir := filepath.Join(claudeDir, subDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		fn(entry.Name(), entry.IsDir())
	}
}

// scanRulesRecursive 递归扫描 rules 目录
func scanRulesRecursive(claudeDir string, resources *Resources) {
	rulesDir := filepath.Join(claudeDir, "rules")
	filepath.WalkDir(rulesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			relPath, _ := filepath.Rel(rulesDir, path)
			name := strings.TrimSuffix(relPath, ".md")
			resources.Rules = append(resources.Rules, name)
		}
		return nil
	})
}

// scanHooksJSON 解析 hooks.json 配置
func scanHooksJSON(claudeDir string, resources *Resources) error {
	hooksPath := filepath.Join(claudeDir, "hooks", "hooks.json")
	data, err := os.ReadFile(hooksPath)
	if err != nil {
		// 文件不存在是正常的，静默忽略
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read hooks.json: %w", err)
	}

	var cfg struct {
		Hooks map[string]any `json:"hooks"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse hooks.json: %w", err)
	}

	for name := range cfg.Hooks {
		resources.Hooks = append(resources.Hooks, name)
	}
	return nil
}

// scanMCPServers 从 .claude.json 读取 MCP servers 完整配置
func scanMCPServers(envDir string, resources *Resources) error {
	mcpPath := filepath.Join(envDir, ".claude.json")
	data, err := os.ReadFile(mcpPath)
	if err != nil {
		// 文件不存在是正常的，静默忽略
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read .claude.json: %w", err)
	}

	var cfg struct {
		MCPServers map[string]MCPServerConfig `json:"mcpServers"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse .claude.json: %w", err)
	}

	resources.MCPServers = cfg.MCPServers
	return nil
}

// sortResources 对所有资源列表排序
func sortResources(r *Resources) {
	sort.Strings(r.Skills)
	sort.Strings(r.Agents)
	sort.Strings(r.Commands)
	sort.Strings(r.Rules)
	sort.Strings(r.Hooks)
	sort.Strings(r.Plugins)
	// MCPServers 是 map，不需要排序
}
