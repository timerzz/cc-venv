package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/env"
)

type resourceListResponse struct {
	Items []string `json:"items"`
}

type resourceContentResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type resourceUpsertRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content"`
}

func ListResourceFiles(c *gin.Context) {
	e, kind, baseDir, recursive, ok := loadResourceContext(c)
	if !ok {
		return
	}

	items, err := listResourceFiles(baseDir, recursive)
	if err != nil {
		InternalError(c, fmt.Sprintf("list %s: %v", kind, err))
		return
	}

	Success(c, resourceListResponse{Items: items})
	_ = e
}

func GetResourceFile(c *gin.Context) {
	_, kind, baseDir, recursive, ok := loadResourceContext(c)
	if !ok {
		return
	}

	name := c.Query("name")
	if name == "" {
		BadRequest(c, "name is required")
		return
	}

	path, err := resourceFilePath(baseDir, kind, name, recursive)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			NotFound(c, fmt.Sprintf("%s %q not found", kind, name))
			return
		}
		InternalError(c, fmt.Sprintf("read %s: %v", kind, err))
		return
	}

	Success(c, resourceContentResponse{
		Name:    name,
		Content: string(data),
	})
}

func UpsertResourceFile(c *gin.Context) {
	_, kind, baseDir, recursive, ok := loadResourceContext(c)
	if !ok {
		return
	}

	var req resourceUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "invalid request body")
		return
	}

	path, err := resourceFilePath(baseDir, kind, req.Name, recursive)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		InternalError(c, fmt.Sprintf("create %s directory: %v", kind, err))
		return
	}

	if err := os.WriteFile(path, []byte(req.Content), 0o644); err != nil {
		InternalError(c, fmt.Sprintf("write %s: %v", kind, err))
		return
	}

	Success(c, resourceContentResponse{
		Name:    req.Name,
		Content: req.Content,
	})
}

func DeleteResourceFile(c *gin.Context) {
	_, kind, baseDir, recursive, ok := loadResourceContext(c)
	if !ok {
		return
	}

	name := c.Query("name")
	if name == "" {
		BadRequest(c, "name is required")
		return
	}

	path, err := resourceFilePath(baseDir, kind, name, recursive)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			NotFound(c, fmt.Sprintf("%s %q not found", kind, name))
			return
		}
		InternalError(c, fmt.Sprintf("delete %s: %v", kind, err))
		return
	}

	Success(c, gin.H{"name": name})
}

func loadResourceContext(c *gin.Context) (env.Environment, string, string, bool, bool) {
	name := c.Param("name")
	kind := c.Param("kind")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return env.Environment{}, "", "", false, false
	}

	baseDir, recursive, err := resourceBaseDir(e.EnvDir, kind)
	if err != nil {
		BadRequest(c, err.Error())
		return env.Environment{}, "", "", false, false
	}

	return e, kind, baseDir, recursive, true
}

func resourceBaseDir(envDir, kind string) (string, bool, error) {
	switch kind {
	case "agents":
		return filepath.Join(envDir, ".claude", "agents"), false, nil
	case "commands":
		return filepath.Join(envDir, ".claude", "commands"), false, nil
	case "rules":
		return filepath.Join(envDir, ".claude", "rules"), true, nil
	default:
		return "", false, fmt.Errorf("unsupported resource kind %q", kind)
	}
}

func resourceFilePath(baseDir, kind, name string, recursive bool) (string, error) {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return "", fmt.Errorf("name is required")
	}

	cleanName = strings.TrimSuffix(cleanName, ".md")
	cleanName = filepath.Clean(cleanName)
	if cleanName == "." || cleanName == "" {
		return "", fmt.Errorf("invalid %s name", kind)
	}
	if strings.HasPrefix(cleanName, "..") || filepath.IsAbs(cleanName) {
		return "", fmt.Errorf("invalid %s name", kind)
	}
	if !recursive && strings.Contains(cleanName, string(filepath.Separator)) {
		return "", fmt.Errorf("%s names cannot contain subdirectories", kind)
	}

	path := filepath.Join(baseDir, cleanName+".md")
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("resolve %s directory: %w", kind, err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve %s path: %w", kind, err)
	}
	if absPath != absBaseDir+".md" && !strings.HasPrefix(absPath, absBaseDir+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid %s path", kind)
	}

	return path, nil
}

func listResourceFiles(baseDir string, recursive bool) ([]string, error) {
	if _, err := os.Stat(baseDir); err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	items := []string{}
	if recursive {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
				return nil
			}
			rel, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}
			items = append(items, strings.TrimSuffix(filepath.ToSlash(rel), ".md"))
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			items = append(items, strings.TrimSuffix(entry.Name(), ".md"))
		}
	}

	sort.Strings(items)
	return items, nil
}
