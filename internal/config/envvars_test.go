package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteAndReadEnvJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config", "env.json")

	want := EnvVars{
		"ANTHROPIC_API_KEY": "key-123",
		"HTTP_PROXY":        "http://127.0.0.1:7890",
	}

	if err := WriteEnvJSON(path, want); err != nil {
		t.Fatalf("WriteEnvJSON() error = %v", err)
	}

	got, err := ReadEnvJSON(path)
	if err != nil {
		t.Fatalf("ReadEnvJSON() error = %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("ReadEnvJSON() len = %d, want %d", len(got), len(want))
	}

	for k, v := range want {
		if got[k] != v {
			t.Fatalf("ReadEnvJSON()[%q] = %q, want %q", k, got[k], v)
		}
	}
}

func TestReadEnvJSONMissingFileReturnsEmptyMap(t *testing.T) {
	t.Parallel()

	got, err := ReadEnvJSON(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("ReadEnvJSON() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("ReadEnvJSON() len = %d, want 0", len(got))
	}
}

func TestReadEnvJSONInvalidJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "env.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	_, err := ReadEnvJSON(path)
	if err == nil || !strings.Contains(err.Error(), "parse env.json") {
		t.Fatalf("ReadEnvJSON() error = %v, want parse env.json error", err)
	}
}

func TestMergeEnvOverlayWins(t *testing.T) {
	t.Parallel()

	base := []string{
		"KEEP=value",
		"OVERRIDE=old",
	}
	overlay := EnvVars{
		"OVERRIDE": "new",
		"ADDED":    "yes",
	}

	got := MergeEnv(base, overlay)
	joined := strings.Join(got, "\n")

	for _, want := range []string{"KEEP=value", "OVERRIDE=new", "ADDED=yes"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("MergeEnv() missing %q in %q", want, joined)
		}
	}
}
