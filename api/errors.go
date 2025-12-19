// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import "fmt"

// Error représente une erreur de l'API
type Error struct {
	Type    ErrorType
	Message string
	Cause   error
}

// ErrorType représente le type d'erreur
type ErrorType string

const (
	ErrorTypeParse      ErrorType = "parse"
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeExecution  ErrorType = "execution"
	ErrorTypeConfig     ErrorType = "config"
	ErrorTypeIO         ErrorType = "io"
	ErrorTypeInternal   ErrorType = "internal"
)

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// ParseError représente une erreur de parsing avec position
type ParseError struct {
	Filename string
	Line     int
	Column   int
	Message  string
	Cause    error
}

func (e *ParseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s:%d:%d: %s: %v", e.Filename, e.Line, e.Column, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s:%d:%d: %s", e.Filename, e.Line, e.Column, e.Message)
}

func (e *ParseError) Unwrap() error {
	return e.Cause
}

// ConfigError représente une erreur de configuration
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration invalide pour '%s': %s", e.Field, e.Message)
}

// XupleSpaceError représente une erreur liée aux xuple-spaces
type XupleSpaceError struct {
	SpaceName string
	Operation string
	Message   string
	Cause     error
}

func (e *XupleSpaceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("xuple-space '%s': %s: %s: %v", e.SpaceName, e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("xuple-space '%s': %s: %s", e.SpaceName, e.Operation, e.Message)
}

func (e *XupleSpaceError) Unwrap() error {
	return e.Cause
}
