// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import "fmt"

// HasErrors returns true if there are any validation errors
func (ps *ProgramState) HasErrors() bool {
	return len(ps.Errors) > 0
}

// GetErrorCount returns the number of validation errors
func (ps *ProgramState) GetErrorCount() int {
	return len(ps.Errors)
}

// GetErrors returns a copy of all validation errors
func (ps *ProgramState) GetErrors() []ValidationError {
	// Return a copy to prevent external modifications
	errors := make([]ValidationError, len(ps.Errors))
	copy(errors, ps.Errors)
	return errors
}

// PrintErrors prints all validation errors to console
func (ps *ProgramState) PrintErrors() {
	for _, err := range ps.Errors {
		fmt.Printf("⚠️  Error in %s: %s (line %d): %s\n",
			err.File, err.Type, err.Line, err.Message)
	}
}

// AddError adds a validation error to the program state
func (ps *ProgramState) AddError(err ValidationError) {
	ps.Errors = append(ps.Errors, err)
}

// ClearErrors removes all validation errors
func (ps *ProgramState) ClearErrors() {
	ps.Errors = []ValidationError{}
}

// Reset clears all state including types, rules, facts, and errors.
// This completely resets the program state to an empty state.
// Use this when you need to start fresh with a new program.
func (ps *ProgramState) Reset() {
	ps.Types = make(map[string]*TypeDefinition)
	ps.Rules = make([]*Expression, 0)
	ps.Facts = make([]*Fact, 0)
	ps.FilesParsed = make([]string, 0)
	ps.Errors = make([]ValidationError, 0)
}
