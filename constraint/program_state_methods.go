// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import "github.com/treivax/tsd/tsdio"

// GetTypes returns an immutable copy of all type definitions
func (ps *ProgramState) GetTypes() map[string]TypeDefinition {
	result := make(map[string]TypeDefinition)
	for k, v := range ps.types {
		result[k] = *v // Defensive copy
	}
	return result
}

// GetRules returns an immutable copy of all rules
func (ps *ProgramState) GetRules() []Expression {
	result := make([]Expression, len(ps.rules))
	for i, rule := range ps.rules {
		result[i] = *rule // Defensive copy
	}
	return result
}

// GetFacts returns an immutable copy of all facts
func (ps *ProgramState) GetFacts() []Fact {
	result := make([]Fact, len(ps.facts))
	for i, fact := range ps.facts {
		result[i] = *fact // Defensive copy
	}
	return result
}

// GetFilesParsed returns an immutable copy of parsed files list
func (ps *ProgramState) GetFilesParsed() []string {
	result := make([]string, len(ps.filesParsed))
	copy(result, ps.filesParsed)
	return result
}

// GetTypesCount returns the number of type definitions
func (ps *ProgramState) GetTypesCount() int {
	return len(ps.types)
}

// GetRulesCount returns the number of rules
func (ps *ProgramState) GetRulesCount() int {
	return len(ps.rules)
}

// GetFactsCount returns the number of facts
func (ps *ProgramState) GetFactsCount() int {
	return len(ps.facts)
}

// HasErrors returns true if there are any validation errors
func (ps *ProgramState) HasErrors() bool {
	return len(ps.errors) > 0
}

// GetErrorCount returns the number of validation errors
func (ps *ProgramState) GetErrorCount() int {
	return len(ps.errors)
}

// GetErrors returns a copy of all validation errors
func (ps *ProgramState) GetErrors() []ValidationError {
	// Return a copy to prevent external modifications
	errors := make([]ValidationError, len(ps.errors))
	copy(errors, ps.errors)
	return errors
}

// PrintErrors prints all validation errors to console
func (ps *ProgramState) PrintErrors() {
	for _, err := range ps.errors {
		tsdio.Printf("⚠️  Error in %s: %s (line %d): %s\n",
			err.File, err.Type, err.Line, err.Message)
	}
}

// AddError adds a validation error to the program state
func (ps *ProgramState) AddError(err ValidationError) {
	ps.errors = append(ps.errors, err)
}

// ClearErrors removes all validation errors
func (ps *ProgramState) ClearErrors() {
	ps.errors = []ValidationError{}
}

// Reset clears all state including types, rules, facts, and errors.
// This completely resets the program state to an empty state.
// Use this when you need to start fresh with a new program.
func (ps *ProgramState) Reset() {
	ps.types = make(map[string]*TypeDefinition)
	ps.rules = make([]*Expression, 0)
	ps.facts = make([]*Fact, 0)
	ps.filesParsed = make([]string, 0)
	ps.errors = make([]ValidationError, 0)
	ps.ruleIDs = make(map[string]bool) // Clear rule IDs on reset
}
