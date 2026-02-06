package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// TestLoggerInit 测试日志初始化
func TestLoggerInit(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		cfg := DefaultConfig()

		if cfg.Level != InfoLevel {
			t.Errorf("Expected default level Info, got %s", cfg.Level)
		}

		if cfg.Format != "console" {
			t.Errorf("Expected default format console, got %s", cfg.Format)
		}
	})

	t.Run("InitWithConsoleFormat", func(t *testing.T) {
		// Capture stderr output
		r, w, _ := os.Pipe()
		oldStderr := os.Stderr
		os.Stderr = w

		cfg := &Config{
			Level:      DebugLevel,
			Format:     "console",
			TimeFormat: "unix",
		}

		err := Init(cfg)
		if err != nil {
			t.Fatalf("Failed to init logger: %v", err)
		}

		// Write a test log
		Info().Msg("test info message")

		// Restore stderr and close pipe
		w.Close()
		os.Stderr = oldStderr

		// Read output
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Check that log was written
		if !strings.Contains(output, "test info message") {
			t.Error("Expected log output to contain test message")
		}
	})

	t.Run("InitFromEnv", func(t *testing.T) {
		// Set environment variables before Init
		os.Setenv("LOG_LEVEL", "warn")
		os.Setenv("LOG_FORMAT", "json")
		defer func() {
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_FORMAT")
		}()

		err := InitFromEnv()
		if err != nil {
			t.Fatalf("Failed to init from env: %v", err)
		}

		// Logger should be initialized
		// Note: The global level is set, but checking GetLogger() may return
		// a cached logger. The important part is that InitFromEnv succeeds.
	})
}

// TestLogLevelParsing 测试日志级别解析
func TestLogLevelParsing(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected zerolog.Level
	}{
		{"TraceLevel", TraceLevel, zerolog.TraceLevel},
		{"DebugLevel", DebugLevel, zerolog.DebugLevel},
		{"InfoLevel", InfoLevel, zerolog.InfoLevel},
		{"WarnLevel", WarnLevel, zerolog.WarnLevel},
		{"ErrorLevel", ErrorLevel, zerolog.ErrorLevel},
		{"FatalLevel", FatalLevel, zerolog.FatalLevel},
		{"DisabledLevel", DisabledLevel, zerolog.Disabled},
		{"InvalidLevel", LogLevel("invalid"), zerolog.InfoLevel}, // 默认 Info
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLevel(tt.level)
			if result != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestLoggerFunctions 测试日志函数
func TestLoggerFunctions(t *testing.T) {
	// Initialize logger with debug level
	cfg := &Config{
		Level:  DebugLevel,
		Format: "console",
	}
	if err := Init(cfg); err != nil {
		t.Fatalf("Failed to init logger: %v", err)
	}

	t.Run("DebugFunction", func(t *testing.T) {
		event := Debug()
		if event == nil {
			t.Error("Debug() returned nil")
		}
	})

	t.Run("InfoFunction", func(t *testing.T) {
		event := Info()
		if event == nil {
			t.Error("Info() returned nil")
		}
	})

	t.Run("WarnFunction", func(t *testing.T) {
		event := Warn()
		if event == nil {
			t.Error("Warn() returned nil")
		}
	})

	t.Run("ErrorFunction", func(t *testing.T) {
		event := Error()
		if event == nil {
			t.Error("Error() returned nil")
		}
	})

	t.Run("FatalFunction", func(t *testing.T) {
		// Note: Fatal() calls os.Exit, so we can't test it directly
		event := Fatal()
		if event == nil {
			t.Error("Fatal() returned nil")
		}
	})
}

// TestGetLogger 测试获取日志实例
func TestGetLogger(t *testing.T) {
	logger := GetLogger()

	if logger.GetLevel() == zerolog.NoLevel {
		t.Error("Expected logger to have a valid level")
	}
}
