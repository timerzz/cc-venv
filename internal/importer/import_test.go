package importer

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/timerzz/cc-venv/internal/env"
	"github.com/timerzz/cc-venv/internal/exporter"
)

func setHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)
}

func TestImportRequiresPath(t *testing.T) {
	_, err := ImportArchive("", Options{})
	if err == nil || !strings.Contains(err.Error(), "archive path is required") {
		t.Fatalf("ImportArchive(\"\") error = %v", err)
	 }
}

func TestImportFileNotFound(t *testing.T) {
	_, err := ImportArchive("nonexistent-archive-for-test.tar.gz", Options{})
	if err == nil {
		t.Fatalf("ImportArchive() should fail for nonexistent file")
	 }
}

func TestImportRoundTrip(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-import")
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

	// 导出
	exportPath := filepath.Join(t.TempDir(), "test-import.tar.gz")
	_, err = exporter.Export("test-import", exporter.Options{OutputPath: exportPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 删除原环境
	if err := env.Remove("test-import", true); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// 导入
	importResult, err := ImportArchive(exportPath, Options{})
	if err != nil {
		t.Fatalf("ImportArchive() error = %v", err)
	}

	// 验证结果
	if importResult.EnvName != "test-import" {
		t.Errorf("EnvName = %v, want test-import", importResult.EnvName)
	}
	if !importResult.Created {
		t.Error("Created = false, want true")
	}

	// 验证环境存在
	loaded, err := env.Load("test-import")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Name != "test-import" {
		t.Errorf("Loaded.Name = %v, want test-import", loaded.Name)
	}

	// 验证 skill 孁在在
	importedSkillDir := filepath.Join(loaded.EnvDir, ".claude", "skills", "test-skill")
	if _, err := os.Stat(importedSkillDir); os.IsNotExist(err) {
		t.Fatal("Imported skill directory not found")
	}
}

func TestImportForceOverwrite(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-force")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 添加原始资源
	skillDir := filepath.Join(e.EnvDir, ".claude", "skills", "original-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// 导出
	exportPath := filepath.Join(t.TempDir(), "test-force.tar.gz")
	_, err = exporter.Export("test-force", exporter.Options{OutputPath: exportPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 修改环境（添加新资源）
	newSkillDir := filepath.Join(e.EnvDir, ".claude", "skills", "new-skill")
	if err := os.MkdirAll(newSkillDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// 不使用 force 导入应该失败
	_, err = ImportArchive(exportPath, Options{})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("ImportArchive() should fail with 'already exists' error, got: %v", err)
	}

	// 使用 force 导入应该成功
	_, err = ImportArchive(exportPath, Options{Force: true})
	if err != nil {
		t.Fatalf("ImportArchive() with Force error = %v", err)
	}

	// 验证 new-skill 被覆盖（不存在）
	if _, err := os.Stat(newSkillDir); !os.IsNotExist(err) {
		t.Fatal("new-skill should be removed after force import")
	}
}

func TestImportInvalidArchive(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建无效的归档文件
	invalidArchive := filepath.Join(t.TempDir(), "invalid.tar.gz")
	if err := os.WriteFile(invalidArchive, []byte("not a valid tar.gz"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	_, err := ImportArchive(invalidArchive, Options{})
	if err == nil {
		t.Fatal("ImportArchive() should fail with invalid archive")
	}
}

func TestImportMissingManifest(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建没有 manifest.json 的归档
	archivePath := filepath.Join(t.TempDir(), "no-manifest.tar.gz")
	if err := createTarWithoutManifest(archivePath); err != nil {
		t.Fatalf("createTarWithoutManifest error = %v", err)
	}

	_, err := ImportArchive(archivePath, Options{})
	if err == nil || !strings.Contains(err.Error(), "manifest.json not found") {
		t.Fatalf("ImportArchive() should fail with 'manifest.json not found' error, got: %v", err)
	}
}

func TestImportInvalidManifest(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	// 创建 manifest.json 无效的归档
	archivePath := filepath.Join(t.TempDir(), "invalid-manifest.tar.gz")
	if err := createTarWithInvalidManifest(archivePath); err != nil {
		t.Fatalf("createTarWithInvalidManifest error = %v", err)
	}

	_, err := ImportArchive(archivePath, Options{})
	if err == nil {
		t.Fatal("ImportArchive() should fail with invalid manifest")
	}
}

// createTarWithoutManifest 创建没有 manifest.json 的 tar.gz
func createTarWithoutManifest(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 添加一个空目录
	hdr := &tar.Header{Name: "test/", Typeflag: tar.TypeDir, Mode: 0o755}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	return nil
}

// createTarWithInvalidManifest 创建 manifest.json 无效的 tar.gz
func createTarWithInvalidManifest(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 添加无效的 manifest.json
	hdr := &tar.Header{Name: "manifest.json", Size: int64(len("invalid")), Mode: 0o644}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := tw.Write([]byte("invalid")); err != nil {
		return err
	}
	return nil
}

func TestVerifyChecksums(t *testing.T) {
	// 创建带 checksums.txt 的归档并验证
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境并导出
	e, err := env.Create("test-checksums")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	exportPath := filepath.Join(t.TempDir(), "test-checksums.tar.gz")
	_, err = exporter.Export("test-checksums", exporter.Options{OutputPath: exportPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 删除原环境
	if err := env.Remove("test-checksums", true); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// 导入（会验证校验和）
	_, err = ImportArchive(exportPath, Options{})
	if err != nil {
		t.Fatalf("ImportArchive() error = %v", err)
	}
}

func TestReadEnvNameFromManifest(t *testing.T) {
	// 测试从 manifest.json 读取环境名
	home := t.TempDir()
	setHome(t, home)

	// 创建测试环境
	e, err := env.Create("test-manifest-name")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	defer os.RemoveAll(e.RootPath)

	// 导出
	exportPath := filepath.Join(t.TempDir(), "test-manifest-name.tar.gz")
	_, err = exporter.Export("test-manifest-name", exporter.Options{OutputPath: exportPath})
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// 删除原环境
	if err := env.Remove("test-manifest-name", true); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// 导入
	result, err := ImportArchive(exportPath, Options{})
	if err != nil {
		t.Fatalf("ImportArchive() error = %v", err)
	}

	// 验证环境名
	if result.EnvName != "test-manifest-name" {
		t.Errorf("EnvName = %v, want test-manifest-name", result.EnvName)
	}
}
