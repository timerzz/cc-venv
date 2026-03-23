package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/timerzz/cc-venv/internal/env"
)

func TestRootCommandContainsExpectedCommands(t *testing.T) {
	t.Parallel()

	cmd := newRootCmd()
	got := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		got[sub.Name()] = true
	}

	for _, want := range []string{"create", "list", "active", "remove", "run", "web", "export", "import"} {
		if !got[want] {
			t.Fatalf("missing subcommand %q", want)
		}
	}
}

func TestListCommandRuns(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"list"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(list) error = %v", err)
	}
}

func TestCreateAndRemoveCommandsRun(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"create", "demo"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(create) error = %v", err)
	}

	if _, err := os.Stat(home + "/.ccv/envs/demo"); err != nil {
		t.Fatalf("expected created environment to exist: %v", err)
	}

	cmd = newRootCmd()
	cmd.SetArgs([]string{"remove", "--force", "demo"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(remove) error = %v", err)
	}
}

func TestExportCommandRequiresExistingEnv(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cmd := newRootCmd()
	cmd.SetArgs([]string{"export", "nonexistent-env"})
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("Execute(export) should fail for nonexistent env")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("Execute(export) error = %v, want substring 'not found'", err)
	}
}

func TestImportCommandRequiresExistingArchive(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cmd := newRootCmd()
	cmd.SetArgs([]string{"import", "nonexistent-archive.tar.gz"})
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("Execute(import) should fail for nonexistent archive")
	}
}

func TestRunCommandPassesThroughClaudeArgs(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	e, err := env.Create("demo")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	toolDir := t.TempDir()
	outDir := t.TempDir()
	claudePath := toolDir + "/claude"
	script := "#!/bin/sh\nprintf '%s\n' \"$@\" > \"$CCV_TEST_OUTPUT_DIR/args.txt\"\n"
	if err := os.WriteFile(claudePath, []byte(script), 0o755); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	t.Setenv("PATH", toolDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("CCV_TEST_OUTPUT_DIR", outDir)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"run", "demo", "--model", "claude-opus", "-p", "hello"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(run) error = %v", err)
	}

	got, err := os.ReadFile(outDir + "/args.txt")
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}

	want := "--append-system-prompt-file\n" + e.ClaudeConfigDir + "/CLAUDE.md\n" +
		"--mcp-config\n" + e.McpConfigPath + "\n" +
		"--model\nclaude-opus\n-p\nhello\n"
	if string(got) != want {
		t.Fatalf("args = %q, want %q", string(got), want)
	}
}
