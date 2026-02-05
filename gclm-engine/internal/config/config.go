package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config 引擎配置
type Config struct {
	Version        string                       `yaml:"version"`
	WorkflowTypes  map[string]WorkflowType      `yaml:"workflow_types"`
	Engine         EngineConfig                 `yaml:"engine"`
	Execution      ExecutionConfig              `yaml:"execution"`
}

// WorkflowType 工作流类型定义
type WorkflowType struct {
	DisplayName string `yaml:"display_name"`
	Description string `yaml:"description"`
}

// EngineConfig 引擎配置
type EngineConfig struct {
	DatabasePath string `yaml:"database_path"`
	WorkflowsDir string `yaml:"workflows_dir"`
	LogLevel     string `yaml:"log_level"`
}

// ExecutionConfig 执行配置
type ExecutionConfig struct {
	DefaultTimeout      int `yaml:"default_timeout"`
	MaxRetries          int `yaml:"max_retries"`
	MaxConcurrentTasks  int `yaml:"max_concurrent_tasks"`
}

var (
	globalConfig *Config
	configOnce   sync.Once
	configFile   string
)

// Load 加载配置文件
func Load() (*Config, error) {
	var err error
	configOnce.Do(func() {
		globalConfig, err = loadConfig()
	})
	return globalConfig, err
}

// loadConfig 实际加载配置
func loadConfig() (*Config, error) {
	// 1. 尝试从 CLI 同目录加载用户配置
	cliDir, err := getCLIDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get CLI directory: %w", err)
	}

	userConfigPath := filepath.Join(cliDir, "gclm_engine_config.yaml")
	if _, err := os.Stat(userConfigPath); err == nil {
		return loadFromFile(userConfigPath)
	}

	// 2. 使用内置默认配置
	return getDefaultConfig()
}

// loadFromFile 从文件加载配置
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置配置文件路径（用于引用）
	configFile = path

	return &cfg, nil
}

// getDefaultConfig 获取内置默认配置
func getDefaultConfig() (*Config, error) {
	// 获取可执行文件所在目录
	execDir, err := getCLIDirectory()
	if err != nil {
		return nil, err
	}

	// 尝试从同目录加载默认配置
	defaultConfigPath := filepath.Join(execDir, "workflow_config_default.yaml")
	if _, err := os.Stat(defaultConfigPath); err == nil {
		return loadFromFile(defaultConfigPath)
	}

	// 如果没有默认配置文件，返回硬编码的默认配置
	return &Config{
		Version: "1.0",
		WorkflowTypes: map[string]WorkflowType{
			"analyze": {
				DisplayName: "代码分析",
				Description: "代码分析、问题诊断、性能评估、架构分析",
			},
			"review": {
				DisplayName: "代码审查",
				Description: "代码审查、安全审计、质量检查",
			},
			"feat": {
				DisplayName: "新功能",
				Description: "新功能开发、模块开发、功能实现",
			},
			"fix": {
				DisplayName: "Bug 修复",
				Description: "Bug 修复、错误处理、问题解决",
			},
			"docs": {
				DisplayName: "文档",
				Description: "文档编写、方案设计、需求分析、API 文档",
			},
			"refactor": {
				DisplayName: "重构",
				Description: "代码重构、架构调整、优化改进",
			},
			"test": {
				DisplayName: "测试",
				Description: "测试编写、测试优化、覆盖率提升",
			},
			"chore": {
				DisplayName: "构建/工具",
				Description: "构建配置、工具升级、依赖更新",
			},
			"style": {
				DisplayName: "代码格式",
				Description: "代码格式调整、样式修改（不影响功能）",
			},
			"perf": {
				DisplayName: "性能优化",
				Description: "性能优化、响应时间优化、资源优化",
			},
			"ci": {
				DisplayName: "CI 配置",
				Description: "CI/CD 配置、自动化脚本",
			},
			"deploy": {
				DisplayName: "部署",
				Description: "部署配置、发布流程、环境配置",
			},
		},
		Engine: EngineConfig{
			DatabasePath: "gclm-engine.db",
			WorkflowsDir: "workflows",
			LogLevel:     "info",
		},
		Execution: ExecutionConfig{
			DefaultTimeout:     300,
			MaxRetries:         3,
			MaxConcurrentTasks: 5,
		},
	}, nil
}

// getCLIDirectory 获取 CLI 可执行文件所在目录
func getCLIDirectory() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// 解析符号链接
	realPath, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	return filepath.Dir(realPath), nil
}

// GetConfigPath 获取当前使用的配置文件路径
func GetConfigPath() string {
	return configFile
}

// ValidateWorkflowType 验证工作流类型是否合法
func (c *Config) ValidateWorkflowType(workflowType string) error {
	if workflowType == "" {
		return fmt.Errorf("workflow_type is required")
	}

	if _, exists := c.WorkflowTypes[workflowType]; !exists {
		return fmt.Errorf("invalid workflow_type: %s (not defined in config)", workflowType)
	}

	return nil
}

// GetWorkflowType 获取工作流类型定义
func (c *Config) GetWorkflowType(workflowType string) (WorkflowType, bool) {
	t, exists := c.WorkflowTypes[workflowType]
	return t, exists
}

// ListWorkflowTypes 列出所有工作流类型
func (c *Config) ListWorkflowTypes() map[string]WorkflowType {
	return c.WorkflowTypes
}

// GetDatabasePath 获取数据库路径
func (c *Config) GetDatabasePath() (string, error) {
	if filepath.IsAbs(c.Engine.DatabasePath) {
		return c.Engine.DatabasePath, nil
	}

	// 相对路径，相对于 CLI 目录
	cliDir, err := getCLIDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(cliDir, c.Engine.DatabasePath), nil
}

// GetWorkflowsDir 获取工作流目录路径
func (c *Config) GetWorkflowsDir() (string, error) {
	if filepath.IsAbs(c.Engine.WorkflowsDir) {
		return c.Engine.WorkflowsDir, nil
	}

	// 相对路径，相对于 CLI 目录
	cliDir, err := getCLIDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(cliDir, c.Engine.WorkflowsDir), nil
}
