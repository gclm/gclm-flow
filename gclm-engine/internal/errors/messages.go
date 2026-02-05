package errors

import (
	"fmt"
	"strings"
)

// ErrorSeverity 错误严重程度
type ErrorSeverity int

const (
	SeverityInfo    ErrorSeverity = iota // 信息
	SeverityWarning                       // 警告
	SeverityError                         // 错误
	SeverityFatal                         // 致命错误
)

// FriendlyError 友好的错误信息
type FriendlyError struct {
	// 简短错误描述（用于日志）
	Short string
	// 用户友好的详细描述
	UserMessage string
	// 可操作的解决建议
	Suggestions []string
	// 错误严重程度
	Severity ErrorSeverity
	// 原始错误
	Err error
}

// Error 实现 error 接口
func (e *FriendlyError) Error() string {
	return e.Short
}

// Unwrap 返回原始错误
func (e *FriendlyError) Unwrap() error {
	return e.Err
}

// FormatUser 返回用户友好的格式化消息
func (e *FriendlyError) FormatUser() string {
	var sb strings.Builder

	sb.WriteString("❌ ")
	sb.WriteString(e.UserMessage)
	sb.WriteString("\n")

	if len(e.Suggestions) > 0 {
		sb.WriteString("\n💡 建议操作:\n")
		for i, s := range e.Suggestions {
			sb.WriteString(fmt.Sprintf("   %d. %s\n", i+1, s))
		}
	}

	return sb.String()
}

// New 创建友好的错误
func New(short, userMessage string, suggestions []string, err error) *FriendlyError {
	if err == nil {
		err = fmt.Errorf(short)
	}
	return &FriendlyError{
		Short:       short,
		UserMessage: userMessage,
		Suggestions: suggestions,
		Severity:    SeverityError,
		Err:         err,
	}
}

// 预定义的友好错误

// WorkflowNotFound 工作流未找到
func WorkflowNotFound(name string) *FriendlyError {
	return &FriendlyError{
		Short:       fmt.Sprintf("workflow '%s' not found", name),
		UserMessage: fmt.Sprintf("工作流 '%s' 不存在", name),
		Suggestions: []string{
			"使用 `gclm-engine workflow list` 查看可用的工作流",
			"检查工作流名称拼写是否正确",
		},
		Severity: SeverityError,
	}
}

// TaskNotFound 任务未找到
func TaskNotFound(taskID string) *FriendlyError {
	return &FriendlyError{
		Short:       fmt.Sprintf("task '%s' not found", taskID),
		UserMessage: fmt.Sprintf("任务 '%s' 不存在", taskID),
		Suggestions: []string{
			"使用 `gclm-engine task list` 查看所有任务",
			"检查任务 ID 是否正确",
		},
		Severity: SeverityError,
	}
}

// PipelineLoadError 流水线加载失败
func PipelineLoadError(name string, err error) *FriendlyError {
	return &FriendlyError{
		Short:       fmt.Sprintf("failed to load pipeline '%s'", name),
		UserMessage: fmt.Sprintf("无法加载工作流配置 '%s'", name),
		Suggestions: []string{
			"检查 workflows/ 目录中是否存在对应的 .yaml 文件",
			"验证 YAML 文件格式是否正确",
			"运行 `gclm-engine workflow list` 查看可用工作流",
		},
		Severity: SeverityError,
		Err:      err,
	}
}

// DatabaseInitError 数据库初始化失败
func DatabaseInitError(path string, err error) *FriendlyError {
	return &FriendlyError{
		Short:       "database initialization failed",
		UserMessage: "无法初始化数据库",
		Suggestions: []string{
			fmt.Sprintf("检查目录权限: %s", path),
			"确保 ~/.gclm-flow/ 目录存在",
			"尝试删除数据库文件后重试",
		},
		Severity: SeverityFatal,
		Err:      err,
	}
}

// PhaseDependencyError 阶段依赖错误
func PhaseDependencyError(phaseName string) *FriendlyError {
	return &FriendlyError{
		Short:       fmt.Sprintf("phase '%s' dependencies not satisfied", phaseName),
		UserMessage: fmt.Sprintf("阶段 '%s' 的前置阶段未完成", phaseName),
		Suggestions: []string{
			"使用 `gclm-engine task plan <task-id>` 查看阶段依赖关系",
			"先完成所有前置阶段",
		},
		Severity: SeverityError,
	}
}

// PhaseAlreadyCompleted 阶段已完成
func PhaseAlreadyCompleted(phaseName string) *FriendlyError {
	return &FriendlyError{
		Short:       fmt.Sprintf("phase '%s' already completed", phaseName),
		UserMessage: fmt.Sprintf("阶段 '%s' 已经完成", phaseName),
		Suggestions: []string{
			"使用 `gclm-engine task current <task-id>` 查看下一阶段",
			"如需重新执行，请使用 `gclm-engine task reset <task-id> <phase-id>`",
		},
		Severity: SeverityWarning,
	}
}

// InvalidYAMLFormat YAML 格式错误
func InvalidYAMLFormat(file string, err error) *FriendlyError {
	return &FriendlyError{
		Short:       fmt.Sprintf("invalid YAML format in %s", file),
		UserMessage: fmt.Sprintf("文件 %s 的 YAML 格式不正确", file),
		Suggestions: []string{
			"检查 YAML 缩进是否正确（使用空格而非制表符）",
			"验证 YAML 语法: https://www.yamllint.com/",
			"参考其他工作流 YAML 文件的格式",
		},
		Severity: SeverityError,
		Err:      err,
	}
}

// ConfigDirectoryNotFound 配置目录未找到
func ConfigDirectoryNotFound(dir string) *FriendlyError {
	return &FriendlyError{
		Short:       "config directory not found",
		UserMessage: "配置目录不存在",
		Suggestions: []string{
			fmt.Sprintf("创建目录: mkdir -p %s", dir),
			"运行 install.sh 重新安装",
		},
		Severity: SeverityFatal,
	}
}

// Wrap 将标准错误包装为友好错误
func Wrap(err error, short, userMessage string, suggestions []string) *FriendlyError {
	return &FriendlyError{
		Short:       short,
		UserMessage: userMessage,
		Suggestions: suggestions,
		Severity:    SeverityError,
		Err:         err,
	}
}

// IsFriendlyError 检查是否为友好错误
func IsFriendlyError(err error) bool {
	_, ok := err.(*FriendlyError)
	return ok
}

// GetUserMessage 获取用户友好的错误消息
func GetUserMessage(err error) string {
	if friendly, ok := err.(*FriendlyError); ok {
		return friendly.FormatUser()
	}
	return err.Error()
}
