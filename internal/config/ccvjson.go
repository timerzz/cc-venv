package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CcvConfig struct {
	SchemaVersion int             `json:"schemaVersion"`
	Name          string          `json:"name"`
	EnvType       string          `json:"envType"`
	Claude        ClaudeConfig    `json:"claude"`
	Resources     ResourceSummary `json:"resources,omitempty"`
}

type ClaudeConfig struct {
	ConfigDirMode string `json:"configDirMode"`
}

type ResourceSummary struct {
	Skills    []string `json:"skills,omitempty"`
	Agents    []string `json:"agents,omitempty"`
	Plugins   []string `json:"plugins,omitempty"`
	Commands  []string `json:"commands,omitempty"`
	Rules     []string `json:"rules,omitempty"`
	Hooks     []string `json:"hooks,omitempty"`
	MCPServer []string `json:"mcpServers,omitempty"`
}

func WriteCcvJSON(path string, cfg CcvConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create manifest directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal ccv.json: %w", err)
	}

	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write ccv.json: %w", err)
	}

	return nil
}
