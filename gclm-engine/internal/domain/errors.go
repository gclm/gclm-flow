package domain

import (
	"errors"
	"fmt"
)

// ============================================================================
// 基础错误类型
// ============================================================================

var (
	// ErrNotFound 资源未找到
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput 无效输入
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict 冲突错误
	ErrConflict = errors.New("conflict")

	// ErrInternal 内部错误
	ErrInternal = errors.New("internal error")

	// ErrUnauthorized 未授权
	ErrUnauthorized = errors.New("unauthorized")
)

// ============================================================================
// Repository 错误 (带基础错误包装)
// ============================================================================

// 定义函数返回包装错误，而非全局变量
func ErrTaskNotFound() error {
	return fmt.Errorf("task not found: %w", ErrNotFound)
}

func ErrPhaseNotFound() error {
	return fmt.Errorf("phase not found: %w", ErrNotFound)
}

func ErrWorkflowNotFound() error {
	return fmt.Errorf("workflow not found: %w", ErrNotFound)
}

// ============================================================================
// Service 错误 (带基础错误包装)
// ============================================================================

func ErrInvalidWorkflowType() error {
	return fmt.Errorf("invalid workflow type: %w", ErrInvalidInput)
}

func ErrPhaseDependencyNotSatisfied() error {
	return fmt.Errorf("phase dependency not satisfied: %w", ErrConflict)
}

func ErrTaskAlreadyCompleted() error {
	return fmt.Errorf("task already completed: %w", ErrConflict)
}

func ErrTaskAlreadyFailed() error {
	return fmt.Errorf("task already failed: %w", ErrConflict)
}

// ErrRequiredPhaseFailed 必需阶段失败
var ErrRequiredPhaseFailed = errors.New("required phase failed")

// ============================================================================
// Workflow 错误 (带基础错误包装)
// ============================================================================

func ErrWorkflowValidationFailed() error {
	return fmt.Errorf("workflow validation failed: %w", ErrInvalidInput)
}

func ErrCircularDependencyDetected() error {
	return fmt.Errorf("circular dependency detected: %w", ErrInvalidInput)
}
