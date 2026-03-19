package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCcvJSON(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "nested", "ccv.json")
	cfg := CcvConfig{
		SchemaVersion: 1,
		Name:          "demo",
		EnvType:       "named",
		Claude: ClaudeConfig{
			ConfigDirMode: "isolated",
		},
		Resources: ResourceSummary{
			Skills: []string{"frontend-design"},
			Agents: []string{"reviewer-agent"},
		},
	}

	if err := WriteCcvJSON(path, cfg); err != nil {
		t.Fatalf("WriteCcvJSON() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	content := string(data)
	for _, part := range []string{`"name": "demo"`, `"envType": "named"`, `"agents": [`} {
		if !strings.Contains(content, part) {
			t.Fatalf("ccv.json missing %q in %q", part, content)
		}
	}
}
