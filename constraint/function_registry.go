// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"sync"
)

// FunctionSignature defines the signature of a built-in function
type FunctionSignature struct {
	Name       string   // Function name
	ReturnType string   // Return type (string, number, bool)
	ParamTypes []string // Parameter types (optional, for future validation)
}

// FunctionRegistry manages built-in function signatures for type inference
type FunctionRegistry struct {
	mu        sync.RWMutex
	functions map[string]*FunctionSignature
}

// NewFunctionRegistry creates a new function registry with default built-in functions
func NewFunctionRegistry() *FunctionRegistry {
	fr := &FunctionRegistry{
		functions: make(map[string]*FunctionSignature),
	}

	// String functions
	fr.register("LENGTH", ValueTypeNumber, nil)
	fr.register("SUBSTRING", ValueTypeString, []string{ValueTypeString, ValueTypeNumber, ValueTypeNumber})
	fr.register("UPPER", ValueTypeString, []string{ValueTypeString})
	fr.register("LOWER", ValueTypeString, []string{ValueTypeString})
	fr.register("TRIM", ValueTypeString, []string{ValueTypeString})

	// Math functions
	fr.register("ABS", ValueTypeNumber, []string{ValueTypeNumber})
	fr.register("ROUND", ValueTypeNumber, []string{ValueTypeNumber})
	fr.register("FLOOR", ValueTypeNumber, []string{ValueTypeNumber})
	fr.register("CEIL", ValueTypeNumber, []string{ValueTypeNumber})

	return fr
}

// register adds a function signature to the registry (internal use)
func (fr *FunctionRegistry) register(name, returnType string, paramTypes []string) {
	fr.functions[strings.ToUpper(name)] = &FunctionSignature{
		Name:       name,
		ReturnType: returnType,
		ParamTypes: paramTypes,
	}
}

// RegisterFunction adds or updates a function signature in the registry
func (fr *FunctionRegistry) RegisterFunction(name, returnType string, paramTypes []string) {
	fr.mu.Lock()
	defer fr.mu.Unlock()
	fr.register(name, returnType, paramTypes)
}

// GetReturnType returns the return type of a function, or the default type if not found
func (fr *FunctionRegistry) GetReturnType(funcName string, defaultType string) string {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	sig, exists := fr.functions[strings.ToUpper(funcName)]
	if !exists {
		return defaultType
	}
	return sig.ReturnType
}

// GetSignature returns the full function signature if it exists
func (fr *FunctionRegistry) GetSignature(funcName string) (*FunctionSignature, bool) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	sig, exists := fr.functions[strings.ToUpper(funcName)]
	return sig, exists
}

// HasFunction checks if a function is registered
func (fr *FunctionRegistry) HasFunction(funcName string) bool {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	_, exists := fr.functions[strings.ToUpper(funcName)]
	return exists
}

// DefaultFunctionRegistry is the global default registry
var DefaultFunctionRegistry = NewFunctionRegistry()
