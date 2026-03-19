package env

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/timerzz/cc-venv/internal/config"
	"github.com/timerzz/cc-venv/internal/platform"
)

func Create(name string) (Environment, error) {
	e, err := newEnvironment(name)
	if err != nil {
		return Environment{}, err
	}

	if err := createLayout(e); err != nil {
		return Environment{}, err
	}

	return e, nil
}

func createLayout(e Environment) error {
	if err := os.MkdirAll(e.EnvDir, 0o755); err != nil {
		return fmt.Errorf("create env dir: %w", err)
	}

	for _, dir := range managedDirs {
		if err := os.MkdirAll(filepath.Join(e.EnvDir, dir), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}

	cfg := config.CcvConfig{
		SchemaVersion: 1,
		Name:          e.Name,
		EnvType:       "named",
		Claude: config.ClaudeConfig{
			ConfigDirMode: "isolated",
		},
	}

	if err := config.WriteCcvJSON(e.ManifestPath, cfg); err != nil {
		return err
	}

	// Create CLAUDE.md in .claude/
	claudeMdPath := filepath.Join(e.EnvDir, ".claude", "CLAUDE.md")
	if err := os.WriteFile(claudeMdPath, []byte("# "+e.Name+"\n\n"), 0o644); err != nil {
		return fmt.Errorf("create CLAUDE.md: %w", err)
	}

	// Create empty .claude.json for MCP configuration
	claudeJsonPath := filepath.Join(e.EnvDir, ".claude.json")
	if err := os.WriteFile(claudeJsonPath, []byte("{}\n"), 0o644); err != nil {
		return fmt.Errorf("create .claude.json: %w", err)
	}

	// Create settings.json in .claude/
	settingsJsonPath := filepath.Join(e.EnvDir, ".claude", "settings.json")
	settingsContent := `{
  "env": {}
}
`
	if err := os.WriteFile(settingsJsonPath, []byte(settingsContent), 0o644); err != nil {
		return fmt.Errorf("create settings.json: %w", err)
	}

	return nil
}

func Load(name string) (Environment, error) {
	e, err := newEnvironment(name)
	if err != nil {
		return Environment{}, err
	}

	info, statErr := os.Stat(e.EnvDir)
	if statErr != nil || !info.IsDir() {
		return Environment{}, fmt.Errorf("environment %q not found", name)
	}

	return e, nil
}

func List() ([]Environment, error) {
	root, err := platform.EnvHome()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return []Environment{}, nil
		}
		return nil, fmt.Errorf("list environments: %w", err)
	}

	envs := make([]Environment, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		e, loadErr := Load(entry.Name())
		if loadErr == nil {
			envs = append(envs, e)
		}
	}

	return envs, nil
}

func Activate(e Environment) error {
	spec, err := prepareExecSpec(e, WorkingDirCurrent)
	if err != nil {
		return err
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	fmt.Printf("[ccv] active environment: %s\n", e.Name)

	cmd := exec.Command(shell)
	cmd.Dir = spec.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = spec.Env

	return cmd.Run()
}

func Remove(name string, force bool) error {
	e, err := Load(name)
	if err != nil {
		return err
	}

	if !force {
		ok, err := confirm(fmt.Sprintf("remove environment %q at %s?", e.Name, e.RootPath))
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("remove cancelled")
		}
	}

	if err := os.RemoveAll(e.RootPath); err != nil {
		return fmt.Errorf("remove environment %q: %w", e.Name, err)
	}

	return nil
}

func newEnvironment(name string) (Environment, error) {
	if name == "" {
		return Environment{}, fmt.Errorf("environment name is required")
	}

	rootPath, err := platform.EnvRoot(name)
	if err != nil {
		return Environment{}, err
	}

	claudeDir := filepath.Join(rootPath, ".claude")

	return Environment{
		Name:            name,
		RootPath:        rootPath,
		EnvDir:          rootPath,
		ManifestPath:    filepath.Join(rootPath, "ccv.json"),
		ClaudeConfigDir: rootPath,
		SettingsPath:    filepath.Join(claudeDir, "settings.json"),
		McpConfigPath:   filepath.Join(rootPath, ".claude.json"),
	}, nil
}

func confirm(prompt string) (bool, error) {
	fmt.Printf("%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
}
