package domain

import (
	"fmt"
)

// ErrorType définit les types d'erreurs du module constraint
type ErrorType string

const (
	// Erreurs de parsing
	ParseError ErrorType = "PARSE_ERROR"

	// Erreurs de validation
	ValidationError    ErrorType = "VALIDATION_ERROR"
	TypeMismatchError  ErrorType = "TYPE_MISMATCH_ERROR"
	FieldNotFoundError ErrorType = "FIELD_NOT_FOUND_ERROR"
	UnknownTypeError   ErrorType = "UNKNOWN_TYPE_ERROR"

	// Erreurs de contraintes
	ConstraintValidationError ErrorType = "CONSTRAINT_ERROR"
	ActionError               ErrorType = "ACTION_ERROR"
)

// Error représente une erreur structurée du module constraint
type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Context Context   `json:"context"`
	Cause   error     `json:"cause,omitempty"`
}

// Context fournit des informations contextuelles sur l'erreur
type Context struct {
	File     string      `json:"file,omitempty"`
	Line     int         `json:"line,omitempty"`
	Column   int         `json:"column,omitempty"`
	Field    string      `json:"field,omitempty"`
	Type     string      `json:"type,omitempty"`
	Variable string      `json:"variable,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	Expected string      `json:"expected,omitempty"`
	Actual   string      `json:"actual,omitempty"`
}

// Error implémente l'interface error
func (ce *Error) Error() string {
	if ce.Context.File != "" {
		return fmt.Sprintf("%s at %s:%d:%d: %s",
			ce.Type, ce.Context.File, ce.Context.Line, ce.Context.Column, ce.Message)
	}
	return fmt.Sprintf("%s: %s", ce.Type, ce.Message)
}

// Unwrap permet le unwrapping des erreurs
func (ce *Error) Unwrap() error {
	return ce.Cause
}

// Is permet la comparaison des types d'erreurs
func (ce *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return ce.Type == t.Type
	}
	return false
}

// Constructeurs d'erreurs

// NewParseError crée une erreur de parsing
func NewParseError(message, file string, line, column int) *Error {
	return &Error{
		Type:    ParseError,
		Message: message,
		Context: Context{
			File:   file,
			Line:   line,
			Column: column,
		},
	}
}

// NewValidationError crée une erreur de validation
func NewValidationError(message string, ctx Context) *Error {
	return &Error{
		Type:    ValidationError,
		Message: message,
		Context: ctx,
	}
}

// NewTypeMismatchError crée une erreur de type incompatible
func NewTypeMismatchError(expected, actual string, ctx Context) *Error {
	ctx.Expected = expected
	ctx.Actual = actual
	return &Error{
		Type:    TypeMismatchError,
		Message: fmt.Sprintf("type mismatch: expected %s, got %s", expected, actual),
		Context: ctx,
	}
}

// NewFieldNotFoundError crée une erreur de champ introuvable
func NewFieldNotFoundError(field, typeName string, ctx Context) *Error {
	ctx.Field = field
	ctx.Type = typeName
	return &Error{
		Type:    FieldNotFoundError,
		Message: fmt.Sprintf("field '%s' not found in type '%s'", field, typeName),
		Context: ctx,
	}
}

// NewUnknownTypeError crée une erreur de type inconnu
func NewUnknownTypeError(typeName string, ctx Context) *Error {
	ctx.Type = typeName
	return &Error{
		Type:    UnknownTypeError,
		Message: fmt.Sprintf("unknown type '%s'", typeName),
		Context: ctx,
	}
}

// NewConstraintError crée une erreur de contrainte
func NewConstraintError(message string, ctx Context) *Error {
	return &Error{
		Type:    ConstraintValidationError,
		Message: message,
		Context: ctx,
	}
}

// NewActionError crée une erreur d'action
func NewActionError(message string, ctx Context) *Error {
	return &Error{
		Type:    ActionError,
		Message: message,
		Context: ctx,
	}
}

// Helpers pour vérifier les types d'erreurs

// IsParseError vérifie si l'erreur est une erreur de parsing
func IsParseError(err error) bool {
	if ce, ok := err.(*Error); ok {
		return ce.Type == ParseError
	}
	return false
}

// IsValidationError vérifie si l'erreur est une erreur de validation
func IsValidationError(err error) bool {
	if ce, ok := err.(*Error); ok {
		return ce.Type == ValidationError ||
			ce.Type == TypeMismatchError ||
			ce.Type == FieldNotFoundError ||
			ce.Type == UnknownTypeError
	}
	return false
}

// IsTypeMismatchError vérifie si l'erreur est une erreur de type incompatible
func IsTypeMismatchError(err error) bool {
	if ce, ok := err.(*Error); ok {
		return ce.Type == TypeMismatchError
	}
	return false
}

// IsFieldNotFoundError vérifie si l'erreur est une erreur de champ introuvable
func IsFieldNotFoundError(err error) bool {
	if ce, ok := err.(*Error); ok {
		return ce.Type == FieldNotFoundError
	}
	return false
}

// IsUnknownTypeError vérifie si l'erreur est une erreur de type inconnu
func IsUnknownTypeError(err error) bool {
	if ce, ok := err.(*Error); ok {
		return ce.Type == UnknownTypeError
	}
	return false
}

// ErrorCollection permet de collecter plusieurs erreurs
type ErrorCollection struct {
	Errors []*Error `json:"errors"`
}

// NewErrorCollection crée une nouvelle collection d'erreurs
func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{
		Errors: make([]*Error, 0),
	}
}

// Add ajoute une erreur à la collection
func (ec *ErrorCollection) Add(err *Error) {
	ec.Errors = append(ec.Errors, err)
}

// HasErrors vérifie s'il y a des erreurs
func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.Errors) > 0
}

// Error implémente l'interface error pour la collection
func (ec *ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}
	if len(ec.Errors) == 1 {
		return ec.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors: %s (and %d more)",
		len(ec.Errors), ec.Errors[0].Error(), len(ec.Errors)-1)
}

// First retourne la première erreur ou nil
func (ec *ErrorCollection) First() *Error {
	if len(ec.Errors) > 0 {
		return ec.Errors[0]
	}
	return nil
}
