package platform

import (
	"fmt"
	"os"
	"path/filepath"
)

func HomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}

	return home, nil
}

func GlobalEnvRoot(name string) (string, error) {
	home, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ccv", "envs", name), nil
}

func EnvHome() (string, error) {
	home, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ccv", "envs"), nil
}

func EnvRoot(name string) (string, error) {
	return GlobalEnvRoot(name)
}
