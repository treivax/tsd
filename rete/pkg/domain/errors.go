package domain

import (
	"errors"
	"fmt"
)

// Erreurs de base du domaine RETE
var (
	ErrFactNotFound     = errors.New("fact not found")
	ErrInvalidFactType  = errors.New("invalid fact type")
	ErrInvalidFieldType = errors.New("invalid field type")
	ErrNodeNotFound     = errors.New("node not found")
	ErrStorageError     = errors.New("storage error")
	ErrValidationFailed = errors.New("validation failed")
)

// ValidationError représente une erreur de validation avec contexte.
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s (value: %v)", e.Field, e.Message, e.Value)
}

// NodeError représente une erreur spécifique à un nœud.
type NodeError struct {
	NodeID   string
	NodeType string
	Cause    error
}

func (e *NodeError) Error() string {
	return fmt.Sprintf("error in node %s (%s): %v", e.NodeID, e.NodeType, e.Cause)
}

func (e *NodeError) Unwrap() error {
	return e.Cause
}

// NewValidationError crée une nouvelle erreur de validation.
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// NewNodeError crée une nouvelle erreur de nœud.
func NewNodeError(nodeID, nodeType string, cause error) *NodeError {
	return &NodeError{
		NodeID:   nodeID,
		NodeType: nodeType,
		Cause:    cause,
	}
}
