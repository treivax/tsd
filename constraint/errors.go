package constraint

import (
	"fmt"
	"strings"
)

// ValidationError represents a non-blocking validation error that occurs during parsing.
// It captures the context of the error including the file, type of element (fact/rule),
// the error message, and optionally the line number where the error occurred.
type ValidationError struct {
	File    string // Source file where the error occurred
	Type    string // Type of element: "fact", "rule", or "type"
	Message string // Descriptive error message
	Line    int    // Line number in the source file (0 if unknown)
}

// Error implements the error interface for ValidationError.
func (ve ValidationError) Error() string {
	if ve.Line > 0 {
		return fmt.Sprintf("%s:%d: %s in %s", ve.File, ve.Line, ve.Message, ve.Type)
	}
	return fmt.Sprintf("%s: %s in %s", ve.File, ve.Message, ve.Type)
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

// Error implements the error interface for ValidationErrors.
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "no validation errors"
	}
	if len(ve) == 1 {
		return ve[0].Error()
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d validation errors:\n", len(ve))
	for i, err := range ve {
		fmt.Fprintf(&sb, "  %d. %s\n", i+1, err.Error())
	}
	return sb.String()
}

// HasErrors returns true if there are any validation errors.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Count returns the number of validation errors.
func (ve ValidationErrors) Count() int {
	return len(ve)
}

// Error type constants for validation errors
const (
	ErrorTypeFact = "fact"
	ErrorTypeRule = "rule"
	ErrorTypeType = "type"
)
