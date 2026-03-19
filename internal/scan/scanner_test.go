package scan

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestScanEnvironment(t *testing.T) {
	t.Parallel()

	envDir := t.TempDir()
	claudeDir := filepath.Join(envDir, ".claude")

	// 创建目录结构
	dirs := []string{
		filepath.Join(claudeDir, "skills", "skill1"),
		filepath.Join(claudeDir, "skills", "skill2"),
		filepath.Join(claudeDir, "agents"),
		filepath.Join(claudeDir, "commands"),
		filepath.Join(claudeDir, "rules", "subdir"),
		filepath.Join(claudeDir, "hooks"),
		filepath.Join(claudeDir, "plugins", "cache"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", dir, err)
		}
	}

	// 创建 agents 文件
	if err := os.WriteFile(filepath.Join(claudeDir, "agents", "agent1.md"), []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile agent error = %v", err)
	}

	// 创建 commands 文件
	if err := os.WriteFile(filepath.Join(claudeDir, "commands", "cmd1.md"), []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile command error = %v", err)
	}

	// 创建 rules 文件（包括子目录）
	if err := os.WriteFile(filepath.Join(claudeDir, "rules", "rule1.md"), []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile rule error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "rules", "subdir", "rule2.md"), []byte(""), 0o644); err != nil {
		t.Fatalf("WriteFile subrule error = %v", err)
	}

	// 创建 hooks.json
	hooksCfg := struct {
		Hooks map[string]any `json:"hooks"`
	}{
		Hooks: map[string]any{
			"hook1": map[string]any{"event": "PreToolUse"},
		},
	}
	hooksData, _ := json.Marshal(hooksCfg)
	if err := os.WriteFile(filepath.Join(claudeDir, "hooks", "hooks.json"), hooksData, 0o644); err != nil {
		t.Fatalf("WriteFile hooks.json error = %v", err)
	}

	// 创建 plugins（包括 temp 目录）
	if err := os.MkdirAll(filepath.Join(claudeDir, "plugins", "cache", "plugin1"), 0o755); err != nil {
		t.Fatalf("MkdirAll plugin error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(claudeDir, "plugins", "cache", "temp123"), 0o755); err != nil {
		t.Fatalf("MkdirAll temp plugin error = %v", err)
	}

	// 创建 .claude.json
	mcpCfg := struct {
		MCPServers map[string]MCPServerConfig `json:"mcpServers"`
	}{
		MCPServers: map[string]MCPServerConfig{
			"server1": {Command: "node", Args: []string{"server.js"}},
		},
	}
	mcpData, _ := json.Marshal(mcpCfg)
	if err := os.WriteFile(filepath.Join(envDir, ".claude.json"), mcpData, 0o644); err != nil {
		t.Fatalf("WriteFile .claude.json error = %v", err)
	}

	// 执行扫描
	resources, err := ScanEnvironment(envDir, Options{ScanMCP: true})
	if err != nil {
		t.Fatalf("ScanEnvironment() error = %v", err)
	}

	// 验证结果
	if len(resources.Skills) != 2 {
		t.Errorf("Skills = %v, want 2 items", resources.Skills)
	}
	if len(resources.Agents) != 1 || resources.Agents[0] != "agent1" {
		t.Errorf("Agents = %v, want [agent1]", resources.Agents)
	}
	if len(resources.Commands) != 1 || resources.Commands[0] != "cmd1" {
		t.Errorf("Commands = %v, want [cmd1]", resources.Commands)
	}
	if len(resources.Rules) != 2 {
		t.Errorf("Rules = %v, want 2 items", resources.Rules)
	}
	if len(resources.Hooks) != 1 || resources.Hooks[0] != "hook1" {
		t.Errorf("Hooks = %v, want [hook1]", resources.Hooks)
	}
	if len(resources.Plugins) != 1 || resources.Plugins[0] != "plugin1" {
		t.Errorf("Plugins = %v, want [plugin1] (temp* should be ignored)", resources.Plugins)
	}
	if len(resources.MCPServers) != 1 {
		t.Errorf("MCPServers = %v, want 1 item", resources.MCPServers)
	}
}

func TestScanEnvironmentEmpty(t *testing.T) {
	t.Parallel()

	envDir := t.TempDir()

	resources, err := ScanEnvironment(envDir, Options{ScanMCP: true})
	if err != nil {
		t.Fatalf("ScanEnvironment() error = %v", err)
	}

	if len(resources.Skills) != 0 {
		t.Errorf("Skills = %v, want empty", resources.Skills)
	}
	if len(resources.Agents) != 0 {
		t.Errorf("Agents = %v, want empty", resources.Agents)
	}
	if len(resources.Commands) != 0 {
		t.Errorf("Commands = %v, want empty", resources.Commands)
	}
	if len(resources.Rules) != 0 {
		t.Errorf("Rules = %v, want empty", resources.Rules)
	}
	if len(resources.Hooks) != 0 {
		t.Errorf("Hooks = %v, want empty", resources.Hooks)
	}
	if len(resources.Plugins) != 0 {
		t.Errorf("Plugins = %v, want empty", resources.Plugins)
	}
	if len(resources.MCPServers) != 0 {
		t.Errorf("MCPServers = %v, want empty", resources.MCPServers)
	}
}

func TestScanHooksJSONParseError(t *testing.T) {
	t.Parallel()

	envDir := t.TempDir()
	claudeDir := filepath.Join(envDir, ".claude", "hooks")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// 写入无效 JSON
	if err := os.WriteFile(filepath.Join(claudeDir, "hooks.json"), []byte("invalid json"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	_, err := ScanEnvironment(envDir, Options{})
	if err == nil {
		t.Fatal("ScanEnvironment() should fail with invalid hooks.json")
	}
}

func TestScanMCPServersParseError(t *testing.T) {
	t.Parallel()

	envDir := t.TempDir()

	// 写入无效 JSON
	if err := os.WriteFile(filepath.Join(envDir, ".claude.json"), []byte("invalid json"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	_, err := ScanEnvironment(envDir, Options{ScanMCP: true})
	if err == nil {
		t.Fatal("ScanEnvironment() should fail with invalid .claude.json")
	}
}

func TestScanEnvironmentSorted(t *testing.T) {
	t.Parallel()

	envDir := t.TempDir()
	claudeDir := filepath.Join(envDir, ".claude")

	// 创建 skills 目录（名称无序）
	for _, name := range []string{"z-skill", "a-skill", "m-skill"} {
		if err := os.MkdirAll(filepath.Join(claudeDir, "skills", name), 0o755); err != nil {
			t.Fatalf("MkdirAll error = %v", err)
		}
	}

	resources, err := ScanEnvironment(envDir, Options{})
	if err != nil {
		t.Fatalf("ScanEnvironment() error = %v", err)
	}

	// 验证已排序
	expected := []string{"a-skill", "m-skill", "z-skill"}
	for i, v := range expected {
		if i >= len(resources.Skills) || resources.Skills[i] != v {
			t.Errorf("Skills[%d] = %v, want %v", i, resources.Skills[i], v)
		}
	}
}

func TestScanEnvironmentNoMCP(t *testing.T) {
	t.Parallel()

	envDir := t.TempDir()
	claudeDir := filepath.Join(envDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// ScanMCP = false，不应该扫描 .claude.json
	resources, err := ScanEnvironment(envDir, Options{ScanMCP: false})
	if err != nil {
		t.Fatalf("ScanEnvironment() error = %v", err)
	}

	if resources.MCPServers == nil {
		t.Error("MCPServers should be initialized even when ScanMCP is false")
	}
}
