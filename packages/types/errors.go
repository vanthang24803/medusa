package types

import (
	"errors"
	"fmt"
)

// Sentinel errors — use errors.Is to check in handler layer
var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation failed")
)

// AppError wraps error with code + message to return to client
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:    "not_found",
		Message: fmt.Sprintf("%s not found", resource),
		Err:     ErrNotFound,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{Code: "conflict", Message: message, Err: ErrConflict}
}

func NewValidation(message string) *AppError {
	return &AppError{Code: "validation_error", Message: message, Err: ErrValidation}
}

func NewForbidden(message string) *AppError {
	return &AppError{Code: "forbidden", Message: message, Err: ErrForbidden}
}
