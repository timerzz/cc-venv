package exporter

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/timerzz/cc-venv/internal/config"
	"github.com/timerzz/cc-venv/internal/env"
	"github.com/timerzz/cc-venv/internal/scan"
)

const (
	FormatVersion = 1
	ToolName      = "ccv"
)

// ExportMeta 导出元数据
type ExportMeta struct {
	ExportedAt string `json:"exportedAt"`
	Version    string `json:"version"`
}

// Options 导出选项
type Options struct {
	OutputPath string // 输出文件路径，如果为空则使用默认路径
}

// Result 导出结果
type Result struct {
	ArchivePath string // 归档文件路径
	EnvName     string // 环境名称
}

// Export 导出环境到 tar.gz 归档
func Export(name string, opts Options) (*Result, error) {
	if name == "" {
		return nil, fmt.Errorf("export: environment name is required")
	}

	// 加载环境
	e, err := env.Load(name)
	if err != nil {
		return nil, fmt.Errorf("load environment: %w", err)
	}

	// 尝试读取现有 ccv.json，如果不存在则使用默认值
	ccvCfg, err := readCcvConfig(e.ManifestPath)
	if err != nil {
		// ccv.json 不可读，使用默认配置
		ccvCfg = &config.CcvConfig{
			SchemaVersion: 1,
			Name:          name,
			EnvType:       "named",
			Claude: config.ClaudeConfig{
				ConfigDirMode: "isolated",
			},
		}
	}

	// 扫描并刷新 ccv.json（这是真正的事实来源）
	ccvCfg, err = scanEnvResources(e)
	if err != nil {
		return nil, fmt.Errorf("scan environment: %w", err)
	}

	// 确保 name 和 envType 正确
	ccvCfg.Name = name
	if ccvCfg.EnvType == "" {
		ccvCfg.EnvType = "named"
	}

	// 写回 ccv.json
	if err := config.WriteCcvJSON(e.ManifestPath, *ccvCfg); err != nil {
		return nil, fmt.Errorf("write ccv.json: %w", err)
	}

	// 创建临时导出目录
	exportDir, err := os.MkdirTemp("", "ccv-export-*")
	if err != nil {
		return nil, fmt.Errorf("create temp directory: %w", err)
	}
	defer os.RemoveAll(exportDir)

	// 构建导出结构
	if err := buildExportStructure(exportDir, e, ccvCfg); err != nil {
		return nil, fmt.Errorf("build export structure: %w", err)
	}

	// 生成校验和
	if err := generateChecksums(exportDir); err != nil {
		return nil, fmt.Errorf("generate checksums: %w", err)
	}

	// 确定输出路径
	outputPath := opts.OutputPath
	if outputPath == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputPath = fmt.Sprintf("%s-%s.tar.gz", name, timestamp)
	}

	// 创建归档
	if err := createTarball(exportDir, outputPath); err != nil {
		return nil, fmt.Errorf("create archive: %w", err)
	}

	return &Result{
		ArchivePath: outputPath,
		EnvName:     name,
	}, nil
}

// buildExportStructure 构建导出目录结构
func buildExportStructure(exportDir string, e env.Environment, ccvCfg *config.CcvConfig) error {
	// 创建目录结构
	// exportDir/
	// ├── manifest.json
	// ├── ccv.json
	// ├── env/
	// └── meta/

	// 创建 manifest.json
	manifest := config.Manifest{
		FormatVersion: FormatVersion,
		Tool:          ToolName,
		EnvName:       e.Name,
		EnvType:       ccvCfg.EnvType,
		Includes:      getIncludesList(),
	}
	if err := writeJSON(filepath.Join(exportDir, "manifest.json"), manifest); err != nil {
		return fmt.Errorf("write manifest.json: %w", err)
	}

	// 复制 ccv.json
	ccvData, err := os.ReadFile(e.ManifestPath)
	if err != nil {
		return fmt.Errorf("read ccv.json: %w", err)
	}
	if err := os.WriteFile(filepath.Join(exportDir, "ccv.json"), ccvData, 0o644); err != nil {
		return fmt.Errorf("write ccv.json: %w", err)
	}

	// 复制环境文件到 env/
	envDir := filepath.Join(exportDir, "env")
	if err := copyEnvDir(e.EnvDir, envDir); err != nil {
		return fmt.Errorf("copy environment: %w", err)
	}

	// 创建 meta/export.json
	metaDir := filepath.Join(exportDir, "meta")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		return fmt.Errorf("create meta directory: %w", err)
	}

	exportMeta := ExportMeta{
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
		Version:    "1.0.0", // TODO: 从 build info 获取版本
	}
	if err := writeJSON(filepath.Join(metaDir, "export.json"), exportMeta); err != nil {
		return fmt.Errorf("write export.json: %w", err)
	}

	return nil
}

// getIncludesList 获取包含的资源列表
func getIncludesList() []string {
	return []string{
		"ccv.json",
		"env/.claude",
		"env/.claude.json",
	}
}

// copyEnvDir 复制环境目录
func copyEnvDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 跳过排除的路径
		if shouldExclude(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		dstPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 跳过符号链接
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// 复制文件
		return copyFile(path, dstPath, info.Mode())
	})
}

// shouldExclude 判断是否应该排除
func shouldExclude(path string) bool {
	// 排除插件数据
	if strings.HasPrefix(path, ".claude/plugins/data") {
		return true
	}

	// 只排除 cache 中的 temp* 文件
	if strings.HasPrefix(path, ".claude/plugins/cache/") {
		// 获取 cache 下的子路径
		subPath := strings.TrimPrefix(path, ".claude/plugins/cache/")
		// 如果是 temp 开头的文件或目录，排除
		if strings.HasPrefix(subPath, "temp") {
			return true
		}
	}

	return false
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

// generateChecksums 生成校验和文件
func generateChecksums(exportDir string) error {
	checksums := make(map[string]string)

	err := filepath.Walk(exportDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 跳过 checksums.txt 本身
		if filepath.Base(path) == "checksums.txt" {
			return nil
		}

		relPath, err := filepath.Rel(exportDir, path)
		if err != nil {
			return err
		}

		// 计算文件校验和
		hash, err := fileChecksum(path)
		if err != nil {
			return err
		}

		checksums[relPath] = hash
		return nil
	})

	if err != nil {
		return err
	}

	// 写入 checksums.txt
	var lines []string
	for path, hash := range checksums {
		lines = append(lines, fmt.Sprintf("%s  %s", hash, path))
	}

	checksumsPath := filepath.Join(exportDir, "checksums.txt")
	return os.WriteFile(checksumsPath, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
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

// createTarball 创建 tar.gz 归档
func createTarball(srcDir, dstPath string) error {
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	baseDir := filepath.Base(srcDir)

	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 创建 tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.Join(baseDir, relPath)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 写入文件内容
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
}

// readCcvConfig 读取 ccv.json
func readCcvConfig(path string) (*config.CcvConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg config.CcvConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// writeJSON 写入 JSON 文件
func writeJSON(path string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	jsonData = append(jsonData, '\n')
	return os.WriteFile(path, jsonData, 0o644)
}

// scanEnvResources 扫描环境目录获取最新资源列表
func scanEnvResources(e env.Environment) (*config.CcvConfig, error) {
	cfg := &config.CcvConfig{
		SchemaVersion: 1,
		Name:          e.Name,
		EnvType:       "named",
		Claude: config.ClaudeConfig{
			ConfigDirMode: "isolated",
		},
	}

	// 使用统一的 scanner 模块扫描资源
	resources, err := scan.ScanEnvironment(e.EnvDir, scan.Options{ScanMCP: true})
	if err != nil {
		return nil, err
	}

	// 转换到 CcvConfig 结构
	cfg.Resources.Skills = resources.Skills
	cfg.Resources.Agents = resources.Agents
	cfg.Resources.Commands = resources.Commands
	cfg.Resources.Rules = resources.Rules
	cfg.Resources.Hooks = resources.Hooks
	cfg.Resources.Plugins = resources.Plugins
	// MCPServers 是 map，需要提取名称
	for name := range resources.MCPServers {
		cfg.Resources.MCPServer = append(cfg.Resources.MCPServer, name)
	}

	return cfg, nil
}
