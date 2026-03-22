package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Success(c, gin.H{"message": "ok"})

	if w.Code != http.StatusOK {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp Response[map[string]any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if resp.Code != CodeSuccess {
		t.Errorf("Code = %v, want %v", resp.Code, CodeSuccess)
	}

	if resp.Data["message"] != "ok" {
		t.Errorf("message = %v, want ok", resp.Data["message"])
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	BadRequest(c, "invalid request")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var resp Response[any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if resp.Code != CodeError {
		t.Errorf("Code = %v, want %v", resp.Code, CodeError)
	}

	if resp.Msg != "invalid request" {
		t.Errorf("Msg = %v, want invalid request", resp.Msg)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	NotFound(c, "resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusNotFound)
	}

	var resp Response[any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if resp.Code != CodeError {
		t.Errorf("Code = %v, want %v", resp.Code, CodeError)
	}
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	InternalError(c, "internal error")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	var resp Response[any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if resp.Code != CodeError {
		t.Errorf("Code = %v, want %v", resp.Code, CodeError)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/", nil)

	Error(c, "custom error")

	// Error 函数返回 HTTP 200，但 code = CodeError
	if w.Code != http.StatusOK {
		t.Errorf("Status = %v, want %v", w.Code, http.StatusOK)
	}

	var resp Response[any]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if resp.Code != CodeError {
		t.Errorf("Code = %v, want %v", resp.Code, CodeError)
	}

	if resp.Msg != "custom error" {
		t.Errorf("Msg = %v, want custom error", resp.Msg)
	}
}

func TestIsValidEnvName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid simple", "demo", true},
		{"valid with hyphen", "my-env", true},
		{"valid with underscore", "my_env", true},
		{"valid with numbers", "env123", true},
		{"empty", "", false},
		{"too long", "a12345678901234567890123456789012345678901234567890123456789012345", false},
		{"with space", "my env", false},
		{"with dot", "my.env", false},
		{"with slash", "my/env", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidEnvName(tt.input)
			if got != tt.want {
				t.Errorf("isValidEnvName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"normal", "demo", "demo"},
		{"with dots", "demo..test", "demotest"},
		{"with slash", "demo/test", "demotest"},
		{"with backslash", "demo\\test", "demotest"},
		{"with spaces", "  demo  ", "demo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeName(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeName(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
