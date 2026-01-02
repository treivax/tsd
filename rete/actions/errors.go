// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
	"fmt"
)

// ValidationError représente une erreur de validation d'arguments d'action.
type ValidationError struct {
	ActionName string
	Expected   int
	Got        int
}

// Error implémente l'interface error.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("action %s expects %d argument(s), got %d", e.ActionName, e.Expected, e.Got)
}

// NewValidationError crée une nouvelle erreur de validation.
func NewValidationError(actionName string, expected, got int) error {
	return &ValidationError{
		ActionName: actionName,
		Expected:   expected,
		Got:        got,
	}
}

// TypeError représente une erreur de type d'argument d'action.
type TypeError struct {
	ActionName    string
	ArgumentIndex int
	ExpectedType  string
	GotValue      interface{}
}

// Error implémente l'interface error.
func (e *TypeError) Error() string {
	return fmt.Sprintf("action %s expects %s as argument %d, got %T",
		e.ActionName, e.ExpectedType, e.ArgumentIndex, e.GotValue)
}

// NewTypeError crée une nouvelle erreur de type.
func NewTypeError(actionName string, index int, expectedType string, gotValue interface{}) error {
	return &TypeError{
		ActionName:    actionName,
		ArgumentIndex: index,
		ExpectedType:  expectedType,
		GotValue:      gotValue,
	}
}
