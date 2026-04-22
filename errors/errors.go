package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// custom error type
// carries both human-readable and HTTP status code
type GatewayError struct {
	Code    int    //Http status code
	Message string // human readable message
	Err     error  // the underlying error
}

// Error implements the error interface
// Any type with an Error() string method satisfies the error interface
func (e *GatewayError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work through the chain
func (e *GatewayError) Unwrap() error {
	return e.Err
}

// Sentinel Errors - predefined errors you can compare against
// These are the named errors handlers will return
// they are variables, not string, So it's safe to compare
var (
	ErrInvalidRequest = &GatewayError{
		Code:    http.StatusBadRequest,
		Message: "invalid request",
	}
	ErrPromptRequired = &GatewayError{
		Code:    http.StatusBadRequest,
		Message: "prompt is required",
	}
	ErrMethodNotAllowed = &GatewayError{
		Code:    http.StatusMethodNotAllowed,
		Message: "method not allowed",
	}
	ErrModelUnavailable = &GatewayError{
		Code:    http.StatusServiceUnavailable,
		Message: "model is unavailable",
	}
	ErrRateLimited = &GatewayError{
		Code:    http.StatusTooManyRequests,
		Message: "rate limit exceeded",
	}
)

// Creates a GateewayError wrapping and underlying error
func New(base *GatewayError, err error) *GatewayError {
	return &GatewayError{
		Code:    base.Code,
		Message: base.Message,
		Err:     err,
	}
}

// allows to match GatewayErrors by code + message
func (e *GatewayError) Is(target error) bool {
	var t *GatewayError
	if errors.As(target, &t) {
		return e.Code == t.Code && e.Message == t.Message
	}
	return false
}
