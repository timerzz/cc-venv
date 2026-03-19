package importer

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/timerzz/cc-venv/internal/config"
	"github.com/timerzz/cc-venv/internal/env"
)

// ImportMeta 导入元数据
type ImportMeta struct {
	ExportedAt string `json:"exportedAt"`
	Version    string `json:"version"`
}

// Options 导入选项
type Options struct {
	Force bool // 强制覆盖已存在的环境
}

// Result 导入结果
type Result struct {
	EnvName string // 环境名称
	Path    string // 环境路径
	Created bool   // 是否新创建
}

// ImportArchive 从 tar.gz 归档导入环境
func ImportArchive(archivePath string, opts Options) (*Result, error) {
	if archivePath == "" {
		return nil, fmt.Errorf("import: archive path is required")
	}

	// 解压到临时目录
	extractDir, err := os.MkdirTemp("", "ccv-import-*")
	if err != nil {
		return nil, fmt.Errorf("create temp directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// 解压归档
	if err := extractTarball(archivePath, extractDir); err != nil {
		return nil, fmt.Errorf("extract archive: %w", err)
	}

	// 查找导出根目录（可能是 ccv-export/ 或直接在根目录）
	exportRoot, err := findExportRoot(extractDir)
	if err != nil {
		return nil, fmt.Errorf("find export root: %w", err)
	}

	// 读取并验证 manifest.json
	manifest, err := readAndValidateManifest(exportRoot)
	if err != nil {
		return nil, fmt.Errorf("validate manifest: %w", err)
	}

	// 验证校验和
	if err := verifyChecksums(exportRoot); err != nil {
		return nil, fmt.Errorf("verify checksums: %w", err)
	}

	// 优先从 manifest.json 获取环境名称，fallback 到 ccv.json
	envName := manifest.EnvName
	if envName == "" {
		envName, err = readEnvNameFromCcv(exportRoot)
		if err != nil {
			return nil, fmt.Errorf("read env name: %w", err)
		}
	}

	// 检查环境是否已存在
	e, err := env.Load(envName)
	if err == nil {
		// 环境已存在
		if !opts.Force {
			return nil, fmt.Errorf("environment %q already exists, use force option to overwrite", envName)
		}
		// 删除现有环境
		if err := env.Remove(envName, true); err != nil {
			return nil, fmt.Errorf("remove existing environment: %w", err)
		}
	}

	// 创建新环境
	e, err = env.Create(envName)
	if err != nil {
		return nil, fmt.Errorf("create environment: %w", err)
	}

	// 复制环境文件
	envSrcDir := filepath.Join(exportRoot, "env")
	if err := copyEnvFiles(envSrcDir, e.EnvDir); err != nil {
		return nil, fmt.Errorf("copy environment files: %w", err)
	}

	return &Result{
		EnvName: envName,
		Path:    e.RootPath,
		Created: true,
	}, nil
}

// findExportRoot 查找导出根目录
func findExportRoot(extractDir string) (string, error) {
	// 首先检查 manifest.json 是否在根目录
	manifestPath := filepath.Join(extractDir, "manifest.json")
	if _, err := os.Stat(manifestPath); err == nil {
		return extractDir, nil
	}

	// 查找子目录中的 manifest.json
	entries, err := os.ReadDir(extractDir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subDir := filepath.Join(extractDir, entry.Name())
			manifestPath := filepath.Join(subDir, "manifest.json")
			if _, err := os.Stat(manifestPath); err == nil {
				return subDir, nil
			}
		}
	}

	return "", fmt.Errorf("manifest.json not found in archive")
}

// readAndValidateManifest 读取并验证 manifest
func readAndValidateManifest(exportRoot string) (*config.Manifest, error) {
	manifestPath := filepath.Join(exportRoot, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read manifest.json: %w", err)
	}

	var manifest config.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest.json: %w", err)
	}

	// 验证 tool
	if manifest.Tool != "ccv" {
		return nil, fmt.Errorf("invalid tool: expected 'ccv', got %q", manifest.Tool)
	}

	// 验证格式版本
	if manifest.FormatVersion < 1 || manifest.FormatVersion > 1 {
		return nil, fmt.Errorf("unsupported format version: %d", manifest.FormatVersion)
	}

	return &manifest, nil
}

// readEnvNameFromCcv 从 ccv.json 读取环境名称
func readEnvNameFromCcv(exportRoot string) (string, error) {
	ccvPath := filepath.Join(exportRoot, "ccv.json")
	data, err := os.ReadFile(ccvPath)
	if err != nil {
		return "", fmt.Errorf("read ccv.json: %w", err)
	}

	var cfg config.CcvConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("parse ccv.json: %w", err)
	}

	if cfg.Name == "" {
		return "", fmt.Errorf("environment name not found in ccv.json")
	}

	return cfg.Name, nil
}

// verifyChecksums 验证校验和
func verifyChecksums(exportRoot string) error {
	checksumsPath := filepath.Join(exportRoot, "checksums.txt")

	// 如果没有 checksums.txt，跳过验证
	if _, err := os.Stat(checksumsPath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(checksumsPath)
	if err != nil {
		return fmt.Errorf("open checksums.txt: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// 格式: <hash>  <path>
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			continue
		}

		expectedHash := parts[0]
		relPath := parts[1]

		filePath := filepath.Join(exportRoot, relPath)
		actualHash, err := fileChecksum(filePath)
		if err != nil {
			return fmt.Errorf("checksum %s: %w", relPath, err)
		}

		if actualHash != expectedHash {
			return fmt.Errorf("checksum mismatch for %s: expected %s, got %s", relPath, expectedHash, actualHash)
		}
	}

	return scanner.Err()
}

// fileChecksum 计算文件 SHA256 校验和
func fileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// copyEnvFiles 复制环境文件
func copyEnvFiles(srcDir, dstDir string) error {
	// 清空目标目录（保留目录本身）
	entries, err := os.ReadDir(dstDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, entry := range entries {
		path := filepath.Join(dstDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	// 确保目标目录存在
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	// 复制文件
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 跳过符号链接
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		return copyFile(path, dstPath, info.Mode())
	})
}

// copyFile 复制文件
func copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// extractTarball 解压 tar.gz 归档
func extractTarball(srcPath, dstDir string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 安全检查：防止路径遍历
		cleanName := filepath.Clean(header.Name)
		if strings.Contains(cleanName, "..") {
			continue
		}

		// 拒绝绝对路径
		if filepath.IsAbs(cleanName) {
			continue
		}

		targetPath := filepath.Join(dstDir, cleanName)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}

			f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()

		case tar.TypeSymlink:
			// 拒绝符号链接
			continue

		case tar.TypeLink:
			// 拒绝硬链接
			continue
		}
	}

	return nil
}
