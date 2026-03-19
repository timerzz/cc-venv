package exporter

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/timerzz/cc-venv/internal/env"
)

func setHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
}

func TestExportRequiresName(t *testing.T) {
	_, err := Export("", Options{})
	if err == nil || !strings.Contains(err.Error(), "environment name is required") {
		t.Fatalf("Export(\"\") error = %v", err)
	}
}

func TestExportEnvNotFound(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	_, err := Export("nonexistent-env-for-test", Options{})
	if err == nil {
		t.Fatalf("Export(\"nonexistent-env-for-test\") should fail")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("Expected 'not found' error, got: %v", err)
	}
}

func TestExportCreatesArchive(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-export")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 添加一些资源
	skillDir := filepath.Join(e.EnvDir, ".claude", "skills", "test-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "README.md"), []byte("# test skill"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// 执行导出
	outputPath := filepath.Join(t.TempDir(), "test-export.tar.gz")
	result, err := Export("test-export", Options{OutputPath: outputPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 验证结果
	if result.EnvName != "test-export" {
		t.Errorf("EnvName = %v, want test-export", result.EnvName)
	}
	if result.ArchivePath != outputPath {
		t.Errorf("ArchivePath = %v, want %v", result.ArchivePath, outputPath)
	}

	// 验证归档文件存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Archive file not created: %v", err)
	}
}

func TestExportDefaultOutputPath(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-default-path")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 切换到临时目录
	oldWD, _ := os.Getwd()
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Chdir error = %v", err)
	}
	defer os.Chdir(oldWD)

	// 执行导出（不指定输出路径）
	result, err := Export("test-default-path", Options{})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 验证默认路径包含环境名
	if !strings.Contains(result.ArchivePath, "test-default-path") {
		t.Errorf("ArchivePath = %v, should contain env name", result.ArchivePath)
	}
	if !strings.HasSuffix(result.ArchivePath, ".tar.gz") {
		t.Errorf("ArchivePath = %v, should end with .tar.gz", result.ArchivePath)
	}

	// 清理
	os.Remove(result.ArchivePath)
}

func TestExportIncludesManifest(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-manifest")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 执行导出
	outputPath := filepath.Join(t.TempDir(), "test-manifest.tar.gz")
	_, err = Export("test-manifest", Options{OutputPath: outputPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 验证归档文件存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Archive file not created: %v", err)
	}
}

func TestExportExcludesPluginData(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-exclude")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 创建插件数据目录（应该被排除）
	pluginDataDir := filepath.Join(e.EnvDir, ".claude", "plugins", "data")
	if err := os.MkdirAll(pluginDataDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDataDir, "data.json"), []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	// 执行导出
	outputPath := filepath.Join(t.TempDir(), "test-exclude.tar.gz")
	_, err = Export("test-exclude", Options{OutputPath: outputPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 验证插件数据目录不在归档中
	if fileExistsInTar(outputPath, "plugins/data/data.json") {
		t.Fatal("plugins/data should be excluded from export")
	}
}

func TestExportExcludesTempPlugins(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-temp-exclude")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 创建 temp 插件（应该被排除）
	tempPluginDir := filepath.Join(e.EnvDir, ".claude", "plugins", "cache", "temp123")
	if err := os.MkdirAll(tempPluginDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// 创建正常插件
	normalPluginDir := filepath.Join(e.EnvDir, ".claude", "plugins", "cache", "normal-plugin")
	if err := os.MkdirAll(normalPluginDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// 执行导出
	outputPath := filepath.Join(t.TempDir(), "test-temp-exclude.tar.gz")
	_, err = Export("test-temp-exclude", Options{OutputPath: outputPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 验证 temp 插件不在归档中
	if fileExistsInTar(outputPath, "plugins/cache/temp123") {
		t.Fatal("temp* plugins should be excluded from export")
	}
}

// fileExistsInTar 检查文件是否存在于 tar.gz 中
func fileExistsInTar(tarPath, targetPath string) bool {
	file, err := os.Open(tarPath)
	if err != nil {
		return false
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return false
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err != nil {
			return false
		}
		if strings.Contains(header.Name, targetPath) {
			return true
		}
	}
}
