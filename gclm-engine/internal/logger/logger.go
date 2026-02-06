package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Logger 全局日志实例
	Logger zerolog.Logger
)

// LogLevel 日志级别类型
type LogLevel string

const (
	// TraceLevel 最详细的日志级别
	TraceLevel LogLevel = "trace"
	// DebugLevel 调试日志级别
	DebugLevel LogLevel = "debug"
	// InfoLevel 信息日志级别（默认）
	InfoLevel LogLevel = "info"
	// WarnLevel 警告日志级别
	WarnLevel LogLevel = "warn"
	// ErrorLevel 错误日志级别
	ErrorLevel LogLevel = "error"
	// FatalLevel 致命错误日志级别
	FatalLevel LogLevel = "fatal"
	// DisabledLevel 禁用日志
	DisabledLevel LogLevel = "disabled"
)

// Config 日志配置
type Config struct {
	Level      LogLevel `json:"level" yaml:"level"`
	OutputFile string   `json:"output_file" yaml:"output_file"`
	Format     string   `json:"format" yaml:"format"`     // json, console
	TimeFormat string   `json:"time_format" yaml:"time_format"` // unix, rfc3339
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      InfoLevel,
		OutputFile: "", // 输出到 stderr
		Format:     "console",
		TimeFormat: "rfc3339",
	}
}

// Init 初始化日志系统
func Init(cfg *Config) error {
	// 设置日志级别
	level := parseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	// 配置时间格式
	switch cfg.TimeFormat {
	case "unix":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	default:
		// 默认使用 RFC3339 格式
	}

	// 配置输出
	var outputWriter io.Writer = os.Stderr
	if cfg.OutputFile != "" {
		fileWriter, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		outputWriter = fileWriter
	}

	// 配置格式
	if cfg.Format == "console" {
		// 控制台格式：彩色、易读
		Logger = zerolog.New(outputWriter).With().Timestamp().Logger()
	} else {
		// JSON 格式：结构化，适合日志收集
		Logger = zerolog.New(outputWriter).With().Timestamp().Logger()
	}

	// 设置全局日志
	log.Logger = Logger

	return nil
}

// InitFromEnv 从环境变量初始化日志
func InitFromEnv() error {
	cfg := DefaultConfig()

	// 从环境变量读取配置
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		cfg.Level = LogLevel(level)
	}
	if outputFile := os.Getenv("LOG_FILE"); outputFile != "" {
		cfg.OutputFile = outputFile
	}
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		cfg.Format = format
	}

	return Init(cfg)
}

// parseLevel 解析日志级别
func parseLevel(level LogLevel) zerolog.Level {
	switch level {
	case TraceLevel:
		return zerolog.TraceLevel
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case DisabledLevel:
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

// GetLogger 获取全局日志实例
func GetLogger() zerolog.Logger {
	return Logger
}

// WithFields 创建带字段的日志上下文
func WithFields(fields map[string]interface{}) zerolog.Context {
	return Logger.With().Fields(fields)
}

// Debug 创建 debug 级别日志
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info 创建 info 级别日志
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn 创建 warn 级别日志
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error 创建 error 级别日志
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal 创建 fatal 级别日志
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
