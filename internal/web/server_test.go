package web

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	srv := NewServer(Config{
		Port:    3000,
		DevMode: true,
	})
	if srv == nil {
		t.Fatal("NewServer() returned nil")
	}
}
