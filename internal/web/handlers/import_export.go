package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/exporter"
	"github.com/timerzz/cc-venv/internal/importer"
)

// ExportResponse 导出响应
type ExportResponse struct {
	DownloadURL string `json:"downloadUrl"`
}

// ImportResponse 导入响应
type ImportResponse struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// ExportEnv 导出环境
func ExportEnv(c *gin.Context) {
	name := c.Param("name")

	// 使用受控导出目录
	exportDir := filepath.Join(os.TempDir(), "ccv-exports")
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		InternalError(c, fmt.Sprintf("create export directory: %v", err))
		return
	}

	// 生成归档文件名
	timestamp := time.Now().Format("20060102-150405")
	archiveName := fmt.Sprintf("%s-%s.tar.gz", name, timestamp)
	archivePath := filepath.Join(exportDir, archiveName)

	// 调用 exporter 模块
	if _, err := exporter.Export(name, exporter.Options{
		OutputPath: archivePath,
	}); err != nil {
		Error(c, err.Error())
		return
	}

	// 返回下载URL
	Success(c, ExportResponse{
		DownloadURL: fmt.Sprintf("/api/downloads/%s", archiveName),
	})
}

// ImportEnv 导入环境
func ImportEnv(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		BadRequest(c, "no file uploaded")
		return
	}
	defer file.Close()

	// 检查文件扩展名
	if !strings.HasSuffix(header.Filename, ".tar.gz") && !strings.HasSuffix(header.Filename, ".tgz") {
		BadRequest(c, "file must be a .tar.gz archive")
		return
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "import-*.tar.gz")
	if err != nil {
		InternalError(c, fmt.Sprintf("create temp file: %v", err))
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// 保存上传的文件
	if _, err := tmpFile.ReadFrom(file); err != nil {
		InternalError(c, fmt.Sprintf("save uploaded file: %v", err))
		return
	}

	// 检查是否强制覆盖
	force := c.PostForm("force") == "true"

	// 调用 importer 模块
	result, err := importer.ImportArchive(tmpFile.Name(), importer.Options{
		Force: force,
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			BadRequest(c, err.Error())
			return
		}
		Error(c, err.Error())
		return
	}

	Success(c, ImportResponse{
		Name: result.EnvName,
		Path: result.Path,
	})
}

// DownloadFile 下载文件
func DownloadFile(c *gin.Context) {
	filename := c.Param("filename")

	// 安全检查：防止路径遍历
	if strings.Contains(filename, "..") || strings.ContainsAny(filename, "/\\") {
		BadRequest(c, "invalid filename")
		return
	}

	// 只允许下载 ccv-exports 目录下的文件
	exportDir := filepath.Join(os.TempDir(), "ccv-exports")
	filePath := filepath.Join(exportDir, filename)

	// 确保文件路径在导出目录内
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		InternalError(c, "resolve file path")
		return
	}

	absExportDir, err := filepath.Abs(exportDir)
	if err != nil {
		InternalError(c, "resolve export directory")
		return
	}

	if !strings.HasPrefix(absPath, absExportDir) {
		BadRequest(c, "invalid filename")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		NotFound(c, "file not found")
		return
	}

	c.FileAttachment(filePath, filename)
}
