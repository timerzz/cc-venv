package webui

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timerzz/cc-venv/internal/env"
)

// SkillInfo Skill信息
type SkillInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// SkillListResponse Skill列表响应
type SkillListResponse struct {
	Skills []SkillInfo `json:"skills"`
}

// AddSkillRequest 添加Skill请求（URL方式）
type AddSkillRequest struct {
	URL string `json:"url"`
}

// ListSkills 列出Skills
func ListSkills(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	skillsDir := filepath.Join(e.EnvDir, ".claude", "skills")
	skills, err := listSkillsDir(skillsDir)
	if err != nil && !os.IsNotExist(err) {
		InternalError(c, fmt.Sprintf("list skills: %v", err))
		return
	}

	Success(c, SkillListResponse{Skills: skills})
}

// listSkillsDir 列出skills目录
func listSkillsDir(dir string) ([]SkillInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var skills []SkillInfo
	for _, entry := range entries {
		if entry.IsDir() {
			skills = append(skills, SkillInfo{
				Name: entry.Name(),
				Path: filepath.Join(".claude", "skills", entry.Name()),
			})
		}
	}

	// 按名称排序
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return skills, nil
}

// AddSkill 添加Skill（支持URL下载或zip上传）
func AddSkill(c *gin.Context) {
	name := c.Param("name")

	e, err := env.Load(name)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", name))
		return
	}

	skillsDir := filepath.Join(e.EnvDir, ".claude", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		InternalError(c, fmt.Sprintf("create skills directory: %v", err))
		return
	}

	contentType := c.GetHeader("Content-Type")

	var skillName string

	// 判断是URL下载还是文件上传
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// 文件上传方式
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			BadRequest(c, "no file uploaded")
			return
		}
		defer file.Close()

		skillName, err = extractZipFromReader(file, skillsDir)
		if err != nil {
			Error(c, fmt.Sprintf("extract zip: %v", err))
			return
		}
	} else {
		// JSON方式（URL下载）
		var req AddSkillRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			BadRequest(c, "invalid request body")
			return
		}

		if req.URL == "" {
			BadRequest(c, "url is required")
			return
		}

		// 验证 URL 协议必须是 https
		if !strings.HasPrefix(req.URL, "https://") {
			BadRequest(c, "only https URLs are allowed")
			return
		}

		skillName, err = downloadAndExtractSkill(req.URL, skillsDir)
		if err != nil {
			Error(c, err.Error())
			return
		}
	}

	Success(c, SkillInfo{
		Name: skillName,
		Path: filepath.Join(".claude", "skills", skillName),
	})
}

// downloadAndExtractSkill 从 URL 下载并解压 skill（流式落盘）
func downloadAndExtractSkill(url, skillsDir string) (string, error) {
	// 使用带 timeout 的 http.Client
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download from url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "skill-*.zip")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// 限制响应体大小 (100MB) 并流式写入临时文件
	maxSize := int64(100 * 1024 * 1024)
	limitedReader := io.LimitReader(resp.Body, maxSize+1)

	written, err := io.Copy(tmpFile, limitedReader)
	tmpFile.Close()
	if err != nil {
		return "", fmt.Errorf("write temp file: %w", err)
	}
	if written > maxSize {
		return "", fmt.Errorf("file too large (max 100MB)")
	}

	// 检查 zip magic bytes (PK\x03\x04)
	if err := validateZipFile(tmpPath); err != nil {
		return "", err
	}

	// 直接从临时文件路径解压
	return extractZipFromPath(tmpPath, skillsDir)
}

// validateZipFile 检查文件是否为有效的 zip 文件
func validateZipFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file for validation: %w", err)
	}
	defer file.Close()

	// ZIP 文件的 magic bytes: 0x50 0x4B 0x03 0x04 (PK..)
	magic := make([]byte, 4)
	n, err := file.Read(magic)
	if err != nil {
		return fmt.Errorf("read file magic: %w", err)
	}
	if n < 4 {
		return fmt.Errorf("file too small to be a valid zip")
	}
	if magic[0] != 0x50 || magic[1] != 0x4B || magic[2] != 0x03 || magic[3] != 0x04 {
		return fmt.Errorf("invalid file format: expected zip file")
	}
	return nil
}

// extractZipFromReader 从Reader解压zip到skills目录
func extractZipFromReader(reader io.Reader, skillsDir string) (string, error) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "skill-*.zip")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// 复制zip内容到临时文件
	if _, err := io.Copy(tmpFile, reader); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	// 直接从临时文件路径解压
	return extractZipFromPath(tmpPath, skillsDir)
}

// extractZipFromPath 从文件路径解压zip到skills目录
func extractZipFromPath(zipPath, skillsDir string) (string, error) {
	// 打开zip文件
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("open zip: %w", err)
	}
	defer zipReader.Close()

	// 确定skill名称（使用zip文件中的根目录名）
	var skillName string
	var rootDir string

	if len(zipReader.File) > 0 {
		first := zipReader.File[0]
		if first.FileInfo().IsDir() {
			rootDir = first.Name
			skillName = strings.TrimSuffix(rootDir, "/")
			// 处理 nested archive 目录结构
			if idx := strings.Index(skillName, "/"); idx > 0 {
				skillName = skillName[:idx]
			}
		} else {
			// 如果第一个不是目录，使用文件所在目录
			parts := strings.Split(first.Name, "/")
			if len(parts) > 1 {
				skillName = parts[0]
			} else {
				skillName = "skill"
			}
		}
	}

	if skillName == "" {
		skillName = "skill"
	}

	// 清理skill名称
	skillName = sanitizeName(skillName)

	// 检查是否已存在
	skillPath := filepath.Join(skillsDir, skillName)
	if _, err := os.Stat(skillPath); err == nil {
		return "", fmt.Errorf("skill %q already exists", skillName)
	}

	// 获取skillPath的绝对路径用于后续校验
	absSkillPath, err := filepath.Abs(skillPath)
	if err != nil {
		return "", fmt.Errorf("get absolute path of skillPath: %w", err)
	}

	// 解压文件
	for _, file := range zipReader.File {
		// 安全检查
		cleanName := filepath.Clean(file.Name)

		// 拒绝绝对路径（以 / 开头或包含 Windows 盘符如 C:）
		if filepath.IsAbs(cleanName) {
			continue
		}
		// 检查 Windows 盘符路径（如 C:\ 或 C:/）
		if len(cleanName) >= 2 && cleanName[1] == ':' && ((cleanName[0] >= 'A' && cleanName[0] <= 'Z') || (cleanName[0] >= 'a' && cleanName[0] <= 'z')) {
			continue
		}

		// 拒绝包含 .. 的路径（路径遍历攻击）
		if strings.Contains(cleanName, "..") {
			continue
		}

		// 拒绝符号链接类型的条目
		if file.Mode()&os.ModeSymlink != 0 {
			continue
		}

		// 跳过根目录本身
		if file.Name == rootDir {
			continue
		}

		// 计算目标路径
		relPath := strings.TrimPrefix(file.Name, rootDir)
		if relPath == "" {
			continue
		}
		relPath = strings.TrimPrefix(relPath, "/")

		targetPath := filepath.Join(skillPath, relPath)

		// 最终校验 targetPath 必须位于 skillPath 目录内
		absTarget, err := filepath.Abs(targetPath)
		if err != nil {
			continue
		}
		if !strings.HasPrefix(absTarget, absSkillPath+string(filepath.Separator)) {
			continue
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, 0o755)
			continue
		}

		// 确保父目录存在
		os.MkdirAll(filepath.Dir(targetPath), 0o755)

		// 解压文件
		src, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("open file in zip: %w", err)
		}

		dst, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			src.Close()
			return "", fmt.Errorf("create file: %w", err)
		}

		_, err = io.Copy(dst, src)
		dst.Close()
		src.Close()

		if err != nil {
			return "", fmt.Errorf("write file: %w", err)
		}
	}

	return skillName, nil
}

// DeleteSkill 删除Skill
func DeleteSkill(c *gin.Context) {
	envName := c.Param("name")
	skillName := c.Param("skill")

	e, err := env.Load(envName)
	if err != nil {
		NotFound(c, fmt.Sprintf("environment %q not found", envName))
		return
	}

	skillPath := filepath.Join(e.EnvDir, ".claude", "skills", skillName)

	// 检查是否存在
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		NotFound(c, fmt.Sprintf("skill %q not found", skillName))
		return
	}

	// 删除skill目录
	if err := os.RemoveAll(skillPath); err != nil {
		InternalError(c, fmt.Sprintf("delete skill: %v", err))
		return
	}

	Success(c, gin.H{"name": skillName})
}

// sanitizeName 清理名称
func sanitizeName(name string) string {
	// 移除不安全字符
	name = strings.ReplaceAll(name, "..", "")
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")
	name = strings.TrimSpace(name)
	return name
}
