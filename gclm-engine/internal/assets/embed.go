package assets

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var assetsFS embed.FS

// Init initializes the embed package with the filesystem
func Init(fs embed.FS) {
	assetsFS = fs
}

// GetFS returns the embedded filesystem
func GetFS() *embed.FS {
	return &assetsFS
}

// MigrationsFS returns the migrations sub-filesystem
func MigrationsFS() fs.FS {
	sub, err := fs.Sub(assetsFS, "migrations")
	if err != nil {
		return nil
	}
	return sub
}

// WorkflowsFS returns the workflows sub-filesystem
func WorkflowsFS() fs.FS {
	sub, err := fs.Sub(assetsFS, "workflows")
	if err != nil {
		return nil
	}
	return sub
}

// WebFS returns the web sub-filesystem for HTTP serving
// Returns an fs.FS compatible with io/fs interface
func WebFS() fs.FS {
	sub, err := fs.Sub(assetsFS, "web")
	if err != nil {
		return nil
	}
	return sub
}

// GetDefaultConfig returns the default config file content
func GetDefaultConfig() ([]byte, error) {
	return fs.ReadFile(assetsFS, "gclm_engine_config.yaml")
}

// ExportDefaultConfig exports the default config to target directory
func ExportDefaultConfig(targetDir string, force bool) (bool, error) {
	configPath := filepath.Join(targetDir, "gclm_engine_config.yaml")

	// 检查文件是否已存在
	if !force {
		if _, err := os.Stat(configPath); err == nil {
			return false, nil // 文件存在，跳过
		}
	}

	// 确保目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return false, err
	}

	// 读取嵌入文件
	data, err := GetDefaultConfig()
	if err != nil {
		return false, err
	}

	// 写入文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return false, err
	}

	return true, nil
}

// ExportBuiltinWorkflows exports builtin workflows to target directory
func ExportBuiltinWorkflows(targetDir string, force bool) ([]string, error) {
	exported := []string{}

	// 确保目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, err
	}

	// 遍历根文件系统，只处理 workflows/ 目录下的 .yaml 文件
	err := fs.WalkDir(assetsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 只处理 workflows/ 目录下的 .yaml 文件
		if !strings.HasPrefix(path, "workflows/") {
			return nil
		}
		if filepath.Ext(path) != ".yaml" {
			return nil
		}

		// 提取文件名（去掉 workflows/ 前缀）
		filename := filepath.Base(path)
		dstPath := filepath.Join(targetDir, filename)

		// 检查文件是否已存在
		if !force {
			if _, err := os.Stat(dstPath); err == nil {
				return nil // 文件存在，跳过
			}
		}

		// 读取嵌入文件
		data, err := fs.ReadFile(assetsFS, path)
		if err != nil {
			return err
		}

		// 写入文件
		if err := os.WriteFile(dstPath, data, 0644); err != nil {
			return err
		}

		exported = append(exported, filename)
		return nil
	})

	return exported, err
}
