package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/timerzz/cc-venv/internal/config"
)

type ExecSpec struct {
	Dir string
	Env []string
}

func RunClaude(e Environment) error {
	spec, err := prepareExecSpec(e, WorkingDirCurrent)
	if err != nil {
		return err
	}

	args := []string{
		"--append-system-prompt-file",
		filepath.Join(e.ClaudeConfigDir, "CLAUDE.md"),
		"--mcp-config",
		e.McpConfigPath,
	}

	cmd := exec.Command("claude", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = spec.Dir
	cmd.Env = spec.Env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run claude in environment %q: %w", e.Name, err)
	}

	return nil
}

type WorkingDirMode int

const (
	WorkingDirCurrent WorkingDirMode = iota
	WorkingDirEnvRoot
)

func prepareExecSpec(e Environment, dirMode WorkingDirMode) (ExecSpec, error) {
	envVars, err := buildProcessEnv(e)
	if err != nil {
		return ExecSpec{}, err
	}

	dir, err := resolveWorkingDir(e, dirMode)
	if err != nil {
		return ExecSpec{}, err
	}

	return ExecSpec{
		Dir: dir,
		Env: envVars,
	}, nil
}

func buildProcessEnv(e Environment) ([]string, error) {
	userVars, err := config.ReadSettingsJSONEnv(e.SettingsPath)
	if err != nil {
		return nil, err
	}

	userVars["CLAUDE_CONFIG_DIR"] = e.ClaudeConfigDir
	userVars["CCV_ENV_NAME"] = e.Name
	userVars["CCV_ENV_ROOT"] = e.RootPath
	userVars["CCV_ACTIVE"] = "1"

	return config.MergeEnv(os.Environ(), userVars), nil
}

func resolveWorkingDir(e Environment, mode WorkingDirMode) (string, error) {
	switch mode {
	case WorkingDirCurrent:
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get current working directory: %w", err)
		}
		return wd, nil
	case WorkingDirEnvRoot:
		return e.RootPath, nil
	default:
		return "", fmt.Errorf("unsupported working directory mode")
	}
}
