package env

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/timerzz/cc-venv/internal/config"
)

func TestCreateLoadListAndRemove(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	for _, rel := range []string{
		".claude/commands",
		".claude/agents",
		".claude/skills",
		".claude/rules",
		".claude/plugins/cache",
		".claude/plugins/data",
		".claude/CLAUDE.md",
		".claude/settings.json",
		".claude.json",
		"ccv.json",
	} {
		if _, err := os.Stat(filepath.Join(e.RootPath, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	loaded, err := Load("demo")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.RootPath != e.RootPath {
		t.Fatalf("Load().RootPath = %q, want %q", loaded.RootPath, e.RootPath)
	}
	if loaded.ClaudeConfigDir != filepath.Join(e.RootPath, ".claude") {
		t.Fatalf("Load().ClaudeConfigDir = %q, want %q", loaded.ClaudeConfigDir, filepath.Join(e.RootPath, ".claude"))
	}

	envs, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(envs) != 1 || envs[0].Name != "demo" {
		t.Fatalf("List() = %+v, want one env named demo", envs)
	}

	if err := Remove("demo", true); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if _, err := os.Stat(e.RootPath); !os.IsNotExist(err) {
		t.Fatalf("expected removed path to not exist, stat err = %v", err)
	}
}

func TestLoadMissingEnvironment(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	_, err := Load("missing")
	if err == nil || !strings.Contains(err.Error(), `environment "missing" not found`) {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestListReturnsEmptyWhenEnvHomeMissing(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	envs, err := List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(envs) != 0 {
		t.Fatalf("List() len = %d, want 0", len(envs))
	}
}

func TestCreateRejectsEmptyName(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	_, err := Create("")
	if err == nil || !strings.Contains(err.Error(), "environment name is required") {
		t.Fatalf("Create(\"\") error = %v", err)
	}
}

func TestCreateCopiesDefaultClaudeMDWhenPresent(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	defaultClaudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(defaultClaudeDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	want := "# Shared Claude Memory\n\nUse shared defaults.\n"
	if err := os.WriteFile(filepath.Join(defaultClaudeDir, "CLAUDE.md"), []byte(want), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(e.RootPath, ".claude", "CLAUDE.md"))
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	if string(got) != want {
		t.Fatalf("CLAUDE.md = %q, want %q", string(got), want)
	}
}

func TestCreateFallsBackToDefaultClaudeMDTemplateWhenMissing(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(e.RootPath, ".claude", "CLAUDE.md"))
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	want := "# demo\n\n"
	if string(got) != want {
		t.Fatalf("CLAUDE.md = %q, want %q", string(got), want)
	}
}

func TestCreateCopiesDefaultSettingsEnvWhenPresent(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	defaultClaudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(defaultClaudeDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	want := config.EnvVars{
		"ANTHROPIC_API_KEY": "secret",
		"OPENAI_BASE_URL":   "https://example.test/v1",
	}
	if err := config.WriteSettingsJSONEnv(filepath.Join(defaultClaudeDir, "settings.json"), want); err != nil {
		t.Fatalf("WriteSettingsJSONEnv() error = %v", err)
	}

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := config.ReadSettingsJSONEnv(filepath.Join(e.RootPath, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("ReadSettingsJSONEnv() error = %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("len(env) = %d, want %d", len(got), len(want))
	}
	for key, value := range want {
		if got[key] != value {
			t.Fatalf("env[%q] = %q, want %q", key, got[key], value)
		}
	}
}

func TestCreateFallsBackToEmptySettingsEnvWhenMissing(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := config.ReadSettingsJSONEnv(filepath.Join(e.RootPath, ".claude", "settings.json"))
	if err != nil {
		t.Fatalf("ReadSettingsJSONEnv() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("env = %#v, want empty", got)
	}
}

func TestCreateCopiesDefaultClaudeJSONWhenPresent(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	want := "{\n  \"mcpServers\": {\n    \"demo\": {\n      \"command\": \"uvx\"\n    }\n  }\n}\n"
	if err := os.WriteFile(filepath.Join(home, ".claude.json"), []byte(want), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(e.RootPath, ".claude.json"))
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	if string(got) != want {
		t.Fatalf(".claude.json = %q, want %q", string(got), want)
	}
}

func TestCreateFallsBackToEmptyClaudeJSONWhenMissing(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(e.RootPath, ".claude.json"))
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	want := "{}\n"
	if string(got) != want {
		t.Fatalf(".claude.json = %q, want %q", string(got), want)
	}
}

func TestBuildProcessEnvReadsEnvJSONAndOverridesReservedVars(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	err = writeSettingsEnv(e.SettingsPath, config.EnvVars{
		"ANTHROPIC_API_KEY": "secret",
		"CLAUDE_CONFIG_DIR": "wrong",
		"CCV_ENV_NAME":      "wrong",
	})
	if err != nil {
		t.Fatalf("writeSettingsEnv() error = %v", err)
	}

	got, err := buildProcessEnv(e)
	if err != nil {
		t.Fatalf("buildProcessEnv() error = %v", err)
	}

	envMap := envSliceToMap(got)
	if envMap["ANTHROPIC_API_KEY"] != "secret" {
		t.Fatalf("ANTHROPIC_API_KEY = %q, want secret", envMap["ANTHROPIC_API_KEY"])
	}
	if envMap["CLAUDE_CONFIG_DIR"] != e.ClaudeConfigDir {
		t.Fatalf("CLAUDE_CONFIG_DIR = %q, want %q", envMap["CLAUDE_CONFIG_DIR"], e.ClaudeConfigDir)
	}
	if envMap["CCV_ENV_NAME"] != e.Name {
		t.Fatalf("CCV_ENV_NAME = %q, want %q", envMap["CCV_ENV_NAME"], e.Name)
	}
	if envMap["CCV_ENV_ROOT"] != e.RootPath {
		t.Fatalf("CCV_ENV_ROOT = %q, want %q", envMap["CCV_ENV_ROOT"], e.RootPath)
	}
	if envMap["CCV_ACTIVE"] != "1" {
		t.Fatalf("CCV_ACTIVE = %q, want 1", envMap["CCV_ACTIVE"])
	}
}

func TestBuildProcessEnvInvalidEnvJSON(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := os.WriteFile(e.SettingsPath, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	_, err = buildProcessEnv(e)
	if err == nil || !strings.Contains(err.Error(), "parse settings.json") {
		t.Fatalf("buildProcessEnv() error = %v, want parse settings.json error", err)
	}
}

func TestActivateUsesCurrentWorkingDirectoryAndInjectsVars(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := writeSettingsEnv(e.SettingsPath, config.EnvVars{
		"TEST_USER_VAR": "hello",
	}); err != nil {
		t.Fatalf("writeSettingsEnv() error = %v", err)
	}

	workDir := filepath.Join(t.TempDir(), "workspace")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	scriptPath, outDir := writeShellRecorder(t)
	t.Setenv("SHELL", scriptPath)

	oldWD, _ := os.Getwd()
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	t.Setenv("CCV_TEST_OUTPUT_DIR", outDir)
	if err := Activate(e); err != nil {
		t.Fatalf("Activate() error = %v", err)
	}

	_ = w.Close()
	var output bytes.Buffer
	if _, err := output.ReadFrom(r); err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	if !strings.Contains(output.String(), "[ccv] active environment: demo") {
		t.Fatalf("stdout = %q, want active message", output.String())
	}

	assertFileEquals(t, filepath.Join(outDir, "pwd.txt"), workDir)
	assertFileEquals(t, filepath.Join(outDir, "claude_config_dir.txt"), e.ClaudeConfigDir)
	assertFileEquals(t, filepath.Join(outDir, "ccv_env_name.txt"), e.Name)
	assertFileEquals(t, filepath.Join(outDir, "ccv_env_root.txt"), e.RootPath)
	assertFileEquals(t, filepath.Join(outDir, "ccv_active.txt"), "1")
	assertFileEquals(t, filepath.Join(outDir, "test_user_var.txt"), "hello")
}

func TestRunClaudeUsesEnvironmentAndRootPath(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := writeSettingsEnv(e.SettingsPath, config.EnvVars{
		"TEST_USER_VAR": "from-run",
	}); err != nil {
		t.Fatalf("writeSettingsEnv() error = %v", err)
	}

	toolDir := t.TempDir()
	outDir := t.TempDir()
	claudePath := filepath.Join(toolDir, "claude")
	script := "#!/bin/sh\n" +
		"printf '%s' \"$PWD\" > \"$CCV_TEST_OUTPUT_DIR/pwd.txt\"\n" +
		"printf '%s' \"$CLAUDE_CONFIG_DIR\" > \"$CCV_TEST_OUTPUT_DIR/claude_config_dir.txt\"\n" +
		"printf '%s' \"$CCV_ENV_NAME\" > \"$CCV_TEST_OUTPUT_DIR/ccv_env_name.txt\"\n" +
		"printf '%s' \"$CCV_ENV_ROOT\" > \"$CCV_TEST_OUTPUT_DIR/ccv_env_root.txt\"\n" +
		"printf '%s' \"$CCV_ACTIVE\" > \"$CCV_TEST_OUTPUT_DIR/ccv_active.txt\"\n" +
		"printf '%s' \"$TEST_USER_VAR\" > \"$CCV_TEST_OUTPUT_DIR/test_user_var.txt\"\n" +
		"printf '%s\n' \"$@\" > \"$CCV_TEST_OUTPUT_DIR/args.txt\"\n"
	if err := os.WriteFile(claudePath, []byte(script), 0o755); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	t.Setenv("PATH", toolDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CCV_TEST_OUTPUT_DIR", outDir)

	workDir := filepath.Join(t.TempDir(), "project")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}
	oldWD, _ := os.Getwd()
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	if err := RunClaude(e, nil); err != nil {
		t.Fatalf("RunClaude() error = %v", err)
	}

	assertFileEquals(t, filepath.Join(outDir, "pwd.txt"), workDir)
	assertFileEquals(t, filepath.Join(outDir, "claude_config_dir.txt"), e.ClaudeConfigDir)
	assertFileEquals(t, filepath.Join(outDir, "ccv_env_name.txt"), e.Name)
	assertFileEquals(t, filepath.Join(outDir, "ccv_env_root.txt"), e.RootPath)
	assertFileEquals(t, filepath.Join(outDir, "ccv_active.txt"), "1")
	assertFileEquals(t, filepath.Join(outDir, "test_user_var.txt"), "from-run")
	assertFileEquals(
		t,
		filepath.Join(outDir, "args.txt"),
		"--append-system-prompt-file\n"+filepath.Join(e.ClaudeConfigDir, "CLAUDE.md")+"\n"+
			"--mcp-config\n"+e.McpConfigPath+"\n",
	)
}

func TestRunClaudeAppendsPassthroughArgs(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	toolDir := t.TempDir()
	outDir := t.TempDir()
	claudePath := filepath.Join(toolDir, "claude")
	script := "#!/bin/sh\n" +
		"printf '%s\n' \"$@\" > \"$CCV_TEST_OUTPUT_DIR/args.txt\"\n"
	if err := os.WriteFile(claudePath, []byte(script), 0o755); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	t.Setenv("PATH", toolDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CCV_TEST_OUTPUT_DIR", outDir)

	if err := RunClaude(e, []string{"--model", "claude-opus", "-p", "hello"}); err != nil {
		t.Fatalf("RunClaude() error = %v", err)
	}

	assertFileEquals(
		t,
		filepath.Join(outDir, "args.txt"),
		"--append-system-prompt-file\n"+filepath.Join(e.ClaudeConfigDir, "CLAUDE.md")+"\n"+
			"--mcp-config\n"+e.McpConfigPath+"\n"+
			"--model\nclaude-opus\n-p\nhello\n",
	)
}

func TestRemoveRequiresConfirmationWhenNotForced(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	restoreStdin := replaceStdinWithTempFile(t, "n\n")
	defer restoreStdin()

	err = Remove("demo", false)
	if err == nil || !strings.Contains(err.Error(), "remove cancelled") {
		t.Fatalf("Remove() error = %v, want remove cancelled", err)
	}

	if _, err := os.Stat(e.RootPath); err != nil {
		t.Fatalf("expected env to still exist, stat err = %v", err)
	}
}

func TestRemoveDeletesWhenConfirmed(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	restoreStdin := replaceStdinWithTempFile(t, "yes\n")
	defer restoreStdin()

	if err := Remove("demo", false); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	if _, err := os.Stat(e.RootPath); !os.IsNotExist(err) {
		t.Fatalf("expected env to be removed, stat err = %v", err)
	}
}

func TestResolveWorkingDirModes(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	e, err := Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	workDir := filepath.Join(t.TempDir(), "cwd")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	oldWD, _ := os.Getwd()
	if err := os.Chdir(workDir); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	gotCurrent, err := resolveWorkingDir(e, WorkingDirCurrent)
	if err != nil {
		t.Fatalf("resolveWorkingDir(current) error = %v", err)
	}
	if gotCurrent != workDir {
		t.Fatalf("resolveWorkingDir(current) = %q, want %q", gotCurrent, workDir)
	}

	gotEnvRoot, err := resolveWorkingDir(e, WorkingDirEnvRoot)
	if err != nil {
		t.Fatalf("resolveWorkingDir(envroot) error = %v", err)
	}
	if gotEnvRoot != e.RootPath {
		t.Fatalf("resolveWorkingDir(envroot) = %q, want %q", gotEnvRoot, e.RootPath)
	}
}

func setHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
}

// writeSettingsEnv 写入环境变量到 settings.json（测试辅助函数）
func writeSettingsEnv(settingsPath string, vars config.EnvVars) error {
	// 读取现有 settings.json
	data, err := os.ReadFile(settingsPath)
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
	settings["env"] = vars

	// 写回文件
	newData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	newData = append(newData, '\n')
	return os.WriteFile(settingsPath, newData, 0o644)
}

func envSliceToMap(values []string) map[string]string {
	out := make(map[string]string, len(values))
	for _, v := range values {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) == 2 {
			out[parts[0]] = parts[1]
		}
	}
	return out
}

func writeShellRecorder(t *testing.T) (string, string) {
	t.Helper()

	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("os.MkdirAll() error = %v", err)
	}

	scriptPath := filepath.Join(dir, "fake-shell.sh")
	script := "#!/bin/sh\n" +
		"printf '%s' \"$PWD\" > \"$CCV_TEST_OUTPUT_DIR/pwd.txt\"\n" +
		"printf '%s' \"$CLAUDE_CONFIG_DIR\" > \"$CCV_TEST_OUTPUT_DIR/claude_config_dir.txt\"\n" +
		"printf '%s' \"$CCV_ENV_NAME\" > \"$CCV_TEST_OUTPUT_DIR/ccv_env_name.txt\"\n" +
		"printf '%s' \"$CCV_ENV_ROOT\" > \"$CCV_TEST_OUTPUT_DIR/ccv_env_root.txt\"\n" +
		"printf '%s' \"$CCV_ACTIVE\" > \"$CCV_TEST_OUTPUT_DIR/ccv_active.txt\"\n" +
		"printf '%s' \"$TEST_USER_VAR\" > \"$CCV_TEST_OUTPUT_DIR/test_user_var.txt\"\n"
	if err := os.WriteFile(scriptPath, []byte(script), 0o755); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	return scriptPath, outDir
}

func assertFileEquals(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", path, err)
	}
	got := string(data)
	if got != want {
		t.Fatalf("%s = %q, want %q", path, got, want)
	}
}

func replaceStdinWithTempFile(t *testing.T, content string) func() {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "stdin-*")
	if err != nil {
		t.Fatalf("os.CreateTemp() error = %v", err)
	}
	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("WriteString() error = %v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("Seek() error = %v", err)
	}

	old := os.Stdin
	os.Stdin = file

	return func() {
		os.Stdin = old
		_ = file.Close()
	}
}
