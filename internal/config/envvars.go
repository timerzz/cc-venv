package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type EnvVars map[string]string

func WriteEnvJSON(path string, vars EnvVars) error {
	if vars == nil {
		vars = EnvVars{}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create env config directory: %w", err)
	}

	data, err := json.MarshalIndent(vars, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal env.json: %w", err)
	}

	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write env.json: %w", err)
	}

	return nil
}

func ReadEnvJSON(path string) (EnvVars, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return EnvVars{}, nil
		}
		return nil, fmt.Errorf("read env.json: %w", err)
	}

	var vars EnvVars
	if err := json.Unmarshal(data, &vars); err != nil {
		return nil, fmt.Errorf("parse env.json: %w", err)
	}

	if vars == nil {
		vars = EnvVars{}
	}

	return vars, nil
}

// SettingsJSON Claude settings.json 结构
type SettingsJSON struct {
	Env map[string]string `json:"env,omitempty"`
}

// ReadSettingsJSONEnv 从 settings.json 中读取 env 字段
func ReadSettingsJSONEnv(path string) (EnvVars, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return EnvVars{}, nil
		}
		return nil, fmt.Errorf("read settings.json: %w", err)
	}

	var cfg SettingsJSON
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse settings.json: %w", err)
	}

	if cfg.Env == nil {
		cfg.Env = EnvVars{}
	}

	return cfg.Env, nil
}

// WriteSettingsJSONEnv 写入环境变量到 settings.json
func WriteSettingsJSONEnv(path string, envVars EnvVars) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	// 读取现有 settings.json
	data, err := os.ReadFile(path)
	var settings map[string]any
	if err != nil {
		if os.IsNotExist(err) {
			settings = make(map[string]any)
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal(data, &settings); err != nil {
			return err
		}
	}

	// 更新 env 字段
	settings["env"] = envVars

	// 写回文件
	newData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	newData = append(newData, '\n')
	return os.WriteFile(path, newData, 0o644)
}

func MergeEnv(base []string, overlay EnvVars) []string {
	merged := make(map[string]string, len(base)+len(overlay))

	for _, item := range base {
		for i := 0; i < len(item); i++ {
			if item[i] == '=' {
				merged[item[:i]] = item[i+1:]
				break
			}
		}
	}

	for k, v := range overlay {
		merged[k] = v
	}

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+merged[k])
	}

	return out
}
