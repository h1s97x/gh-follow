package errors

import (
	"fmt"
)

// ErrorCode represents an error code
type ErrorCode string

const (
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeRateLimit       ErrorCode = "RATE_LIMIT"
	ErrCodeNetwork         ErrorCode = "NETWORK_ERROR"
	ErrCodeInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrCodeAlreadyExists   ErrorCode = "ALREADY_EXISTS"
	ErrCodeSyncFailed      ErrorCode = "SYNC_FAILED"
	ErrCodeConfigError     ErrorCode = "CONFIG_ERROR"
	ErrCodeUserNotFound    ErrorCode = "USER_NOT_FOUND"
	ErrCodeAPILimitExceeded ErrorCode = "API_LIMIT_EXCEEDED"
)

// Sentinel errors for common cases
var (
	ErrEmptyUsername      = NewError(ErrCodeInvalidInput, "username cannot be empty", nil)
	ErrNoToken            = NewError(ErrCodeUnauthorized, "GitHub token not found. Please run 'gh auth login'", nil)
	ErrUserAlreadyFollowed = NewError(ErrCodeAlreadyExists, "user is already in the follow list", nil)
	ErrUserNotFound       = NewError(ErrCodeNotFound, "user not found in follow list", nil)
	ErrInvalidFormat      = NewError(ErrCodeInvalidInput, "invalid export format. Use 'json' or 'csv'", nil)
	ErrNetworkError       = NewError(ErrCodeNetwork, "network error occurred", nil)
	ErrUnauthorized       = NewError(ErrCodeUnauthorized, "unauthorized. Please check your GitHub token", nil)
	ErrAPILimitExceeded   = NewError(ErrCodeRateLimit, "GitHub API rate limit exceeded", nil)
)

// FollowError represents a structured error
type FollowError struct {
	Code     ErrorCode
	Message  string
	Cause    error
	Username string
}

// Error implements the error interface
func (e *FollowError) Error() string {
	if e.Cause != nil {
		if e.Username != "" {
			return fmt.Sprintf("%s: %s (user: %s, caused by: %v)", e.Code, e.Message, e.Username, e.Cause)
		}
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	if e.Username != "" {
		return fmt.Sprintf("%s: %s (user: %s)", e.Code, e.Message, e.Username)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *FollowError) Unwrap() error {
	return e.Cause
}

// NewError creates a new FollowError
func NewError(code ErrorCode, message string, cause error) *FollowError {
	return &FollowError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewFollowError creates a new FollowError with operation context
func NewFollowError(operation string, username string, cause error) *FollowError {
	message := fmt.Sprintf("operation '%s' failed", operation)
	return &FollowError{
		Code:     ErrCodeNetwork,
		Message:  message,
		Cause:    cause,
		Username: username,
	}
}

// NotFound creates a not found error
func NotFound(message string, cause error) *FollowError {
	return NewError(ErrCodeNotFound, message, cause)
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string, cause error) *FollowError {
	return NewError(ErrCodeUnauthorized, message, cause)
}

// RateLimit creates a rate limit error
func RateLimit(message string, cause error) *FollowError {
	return NewError(ErrCodeRateLimit, message, cause)
}

// Network creates a network error
func Network(message string, cause error) *FollowError {
	return NewError(ErrCodeNetwork, message, cause)
}

// InvalidInput creates an invalid input error
func InvalidInput(message string, cause error) *FollowError {
	return NewError(ErrCodeInvalidInput, message, cause)
}

// AlreadyExists creates an already exists error
func AlreadyExists(message string, cause error) *FollowError {
	return NewError(ErrCodeAlreadyExists, message, cause)
}

// SyncFailed creates a sync failed error
func SyncFailed(message string, cause error) *FollowError {
	return NewError(ErrCodeSyncFailed, message, cause)
}

// ConfigError creates a config error
func ConfigError(message string, cause error) *FollowError {
	return NewError(ErrCodeConfigError, message, cause)
}

// IsNotFoundError checks if error is a not found error
func IsNotFoundError(err error) bool {
	if e, ok := err.(*FollowError); ok {
		return e.Code == ErrCodeNotFound
	}
	return false
}

// IsUnauthorizedError checks if error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	if e, ok := err.(*FollowError); ok {
		return e.Code == ErrCodeUnauthorized
	}
	return false
}

// IsRateLimitError checks if error is a rate limit error
func IsRateLimitError(err error) bool {
	if e, ok := err.(*FollowError); ok {
		return e.Code == ErrCodeRateLimit
	}
	return false
}

// IsNetworkError checks if error is a network error
func IsNetworkError(err error) bool {
	if e, ok := err.(*FollowError); ok {
		return e.Code == ErrCodeNetwork
	}
	return false
}
