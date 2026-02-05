package service

import (
	"context"
	"strings"
	"time"
)

// ErrorType 错误类型
type ErrorType int

const (
	ErrorTypeTemporary ErrorType = iota
	ErrorTypePermanent
	ErrorTypeValidation
	ErrorTypeTimeout
	ErrorTypeCancellation
)

// ClassifiedError 分类错误
type ClassifiedError struct {
	Err      error
	Type     ErrorType
	TaskID   string
	PhaseID  string
	Time     time.Time
}

// ErrorClassifier 错误分类器
type ErrorClassifier struct{}

// NewErrorClassifier 创建错误分类器
func NewErrorClassifier() *ErrorClassifier {
	return &ErrorClassifier{}
}

// Classify 分类错误
func (e *ErrorClassifier) Classify(err error) ErrorType {
	if err == nil {
		return ErrorTypeTemporary
	}

	errStr := err.Error()

	// 上下文取消
	if err == context.Canceled || err == context.DeadlineExceeded {
		return ErrorTypeCancellation
	}

	// 超时错误
	if contains(errStr, "timeout") || contains(errStr, "deadline") {
		return ErrorTypeTimeout
	}

	// 验证错误
	if contains(errStr, "validation") || contains(errStr, "invalid") {
		return ErrorTypeValidation
	}

	// 临时错误
	temporaryPatterns := []string{
		"connection refused",
		"temporary",
		"rate limit",
		"unavailable",
	}

	for _, pattern := range temporaryPatterns {
		if contains(errStr, pattern) {
			return ErrorTypeTemporary
		}
	}

	// 默认为永久错误
	return ErrorTypePermanent
}

// IsRecoverable 判断错误是否可恢复
func (e *ErrorClassifier) IsRecoverable(err error) bool {
	errorType := e.Classify(err)

	switch errorType {
	case ErrorTypeTemporary, ErrorTypeTimeout:
		return true
	case ErrorTypeValidation, ErrorTypePermanent, ErrorTypeCancellation:
		return false
	default:
		return false
	}
}

// contains 字符串包含检查（不区分大小写）
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
