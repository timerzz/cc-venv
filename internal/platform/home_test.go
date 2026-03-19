package platform

import (
	"path/filepath"
	"testing"
)

func TestEnvPaths(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	home, err := HomeDir()
	if err != nil {
		t.Fatalf("HomeDir() error = %v", err)
	}

	envHome, err := EnvHome()
	if err != nil {
		t.Fatalf("EnvHome() error = %v", err)
	}
	if envHome != filepath.Join(home, ".ccv", "envs") {
		t.Fatalf("EnvHome() = %q", envHome)
	}

	root, err := EnvRoot("demo")
	if err != nil {
		t.Fatalf("EnvRoot() error = %v", err)
	}
	if root != filepath.Join(home, ".ccv", "envs", "demo") {
		t.Fatalf("EnvRoot() = %q", root)
	}

	globalRoot, err := GlobalEnvRoot("demo")
	if err != nil {
		t.Fatalf("GlobalEnvRoot() error = %v", err)
	}
	if globalRoot != root {
		t.Fatalf("GlobalEnvRoot() = %q, want %q", globalRoot, root)
	}
}
