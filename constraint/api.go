// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadFileContent reads the entire content of a file as a string
func ReadFileContent(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return string(data), nil
}

// ParseConstraint parses constraint input and returns the AST.
// It takes a filename (for error reporting) and the input bytes to parse.
// Returns the parsed AST as an interface{} or an error if parsing fails.
//
// Example:
//
//	content, _ := os.ReadFile("rules.constraint")
//	ast, err := ParseConstraint("rules.constraint", content)
func ParseConstraint(filename string, input []byte) (interface{}, error) {
	return Parse(filename, input)
}

// ValidateConstraintProgram validates a parsed constraint program AST.
// It performs semantic validation including type checking, variable resolution,
// constraint consistency checks, and action signature validation.
//
// Example:
//
//	ast, _ := ParseConstraint("rules.constraint", content)
//	err := ValidateConstraintProgram(ast)
func ValidateConstraintProgram(result interface{}) error {
	// First perform standard validation
	if err := ValidateProgram(result); err != nil {
		return err
	}

	// Convert to Program for action validation
	program, err := ConvertResultToProgram(result)
	if err != nil {
		return fmt.Errorf("failed to convert result to program: %v", err)
	}

	// Validate action definitions and calls
	return ValidateActionCalls(program)
}

// ParseConstraintFile parses a constraint file from the filesystem.
// It reads the file and parses it, returning the AST or an error.
// This is a convenience function that combines file reading and parsing.
//
// Example:
//
//	ast, err := ParseConstraintFile("rules.constraint")
func ParseConstraintFile(filename string) (interface{}, error) {
	return ParseFile(filename)
}

// ParseFactsFile parses a .facts file using the constraint grammar.
// This allows facts to be parsed using the same grammar as constraint files,
// enabling unified processing of constraint and fact definitions.
//
// Example:
//
//	facts, err := ParseFactsFile("data.facts")
func ParseFactsFile(filename string) (interface{}, error) {
	return ParseFile(filename)
}

// ExtractFactsFromProgram extracts facts from a parsed program and converts them to RETE format
func ExtractFactsFromProgram(result interface{}) ([]map[string]interface{}, error) {
	// Convert result to Program structure
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("JSON conversion error: %v", err)
	}

	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	// Convert facts to RETE format
	reteFacts, err := ConvertFactsToReteFormat(program)
	if err != nil {
		return nil, fmt.Errorf("conversion to RETE format failed: %v", err)
	}
	return reteFacts, nil
}

// ValidateActionCalls validates all action calls in a program against their definitions.
func ValidateActionCalls(program *Program) error {
	// Create validator with action and type definitions
	validator := NewActionValidator(program.Actions, program.Types)

	// First, validate action definitions themselves
	if errs := validator.ValidateActionDefinitions(); len(errs) > 0 {
		// Return first error
		return errs[0]
	}

	// Validate action calls in each rule
	for _, expr := range program.Expressions {
		// Build map of rule variables to their types
		ruleVariables := extractRuleVariablesFromExpression(&expr)

		// Validate each action call
		if expr.Action != nil {
			for _, job := range expr.Action.GetJobs() {
				if err := validator.ValidateActionCall(&job, ruleVariables); err != nil {
					return fmt.Errorf("rule '%s': %v", expr.RuleId, err)
				}
			}
		}
	}

	return nil
}

// extractRuleVariablesFromExpression extracts all variables and their types from an expression
func extractRuleVariablesFromExpression(expr *Expression) map[string]string {
	ruleVariables := make(map[string]string)

	// Extract variables from Set (single pattern, backward compatibility)
	if len(expr.Set.Variables) > 0 {
		for _, v := range expr.Set.Variables {
			ruleVariables[v.Name] = v.DataType
		}
	}

	// Extract variables from Patterns (multiple patterns)
	if expr.Patterns != nil {
		for _, pattern := range expr.Patterns {
			for _, v := range pattern.Variables {
				ruleVariables[v.Name] = v.DataType
			}
		}
	}

	return ruleVariables
}

// ConvertResultToProgram converts parser result to Program structure
func ConvertResultToProgram(result interface{}) (*Program, error) {
	// Convert result to Program structure
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("JSON conversion error: %v", err)
	}

	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %v", err)
	}

	return &program, nil
}

// ConvertToReteProgram converts a Program structure from constraints to the format expected by RETE
func ConvertToReteProgram(program *Program) (interface{}, error) {
	// Convert structures to interface arrays using generic helper
	typesInterface, err := convertSliceToInterfaceArray(program.Types)
	if err != nil {
		return nil, fmt.Errorf("error converting types: %w", err)
	}

	actionsInterface, err := convertSliceToInterfaceArray(program.Actions)
	if err != nil {
		return nil, fmt.Errorf("error converting actions: %w", err)
	}

	expressionsInterface, err := convertSliceToInterfaceArray(program.Expressions)
	if err != nil {
		return nil, fmt.Errorf("error converting expressions: %w", err)
	}

	// Create program structure
	reteProgram := map[string]interface{}{
		JSONKeyTypes:       typesInterface,
		JSONKeyActions:     actionsInterface,
		JSONKeyExpressions: expressionsInterface,
	}

	// Add rule removals if present
	if len(program.RuleRemovals) > 0 {
		ruleRemovalsInterface, err := convertSliceToInterfaceArray(program.RuleRemovals)
		if err != nil {
			return nil, fmt.Errorf("error converting rule removals: %w", err)
		}
		reteProgram[JSONKeyRuleRemovals] = ruleRemovalsInterface
	}

	return reteProgram, nil
}

// convertSliceToInterfaceArray converts a slice of any type to []interface{} via JSON marshaling
func convertSliceToInterfaceArray[T any](items []T) ([]interface{}, error) {
	result := make([]interface{}, 0, len(items))

	for _, item := range items {
		// Serialize to JSON then deserialize to map
		jsonData, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal item: %w", err)
		}

		var itemMap map[string]interface{}
		if err := json.Unmarshal(jsonData, &itemMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}

		result = append(result, itemMap)
	}

	return result, nil
}

// NewIterativeParser creates a new iterative parser that can parse multiple files
// and maintain state across parsing operations. This enables parsing type definitions,
// rules, and facts from separate files while ensuring consistency.
//
// Example:
//
//	parser := NewIterativeParser()
//	err := parser.ParseFile("types.constraint")
//	err = parser.ParseFile("rules.constraint")
//	err = parser.ParseFile("facts.constraint")
//	program := parser.GetProgram()
func NewIterativeParser() *IterativeParser {
	return &IterativeParser{
		state: NewProgramState(),
	}
}

// IterativeParser provides iterative parsing capabilities with state management
type IterativeParser struct {
	state *ProgramState
}

// ParseFile parses a file and merges it with the current state.
// The file can contain types, rules, facts, or any combination.
// Validation is performed to ensure consistency with previously parsed content.
func (ip *IterativeParser) ParseFile(filename string) error {
	return ip.state.ParseAndMerge(filename)
}

// ParseContent parses content from a string and merges it with current state.
func (ip *IterativeParser) ParseContent(content, filename string) error {
	return ip.state.ParseAndMergeContent(content, filename)
}

// GetProgram returns the current combined program state
func (ip *IterativeParser) GetProgram() *Program {
	return ip.state.ToProgram()
}

// GetState returns the internal state for advanced usage
func (ip *IterativeParser) GetState() *ProgramState {
	return ip.state
}

// Reset clears all parsed content and resets to empty state
func (ip *IterativeParser) Reset() {
	ip.state = NewProgramState()
}

// GetParsingStatistics returns statistics about the parsing process
func (ip *IterativeParser) GetParsingStatistics() ParsingStatistics {
	return ParsingStatistics{
		TypesCount:       ip.state.GetTypesCount(),
		RulesCount:       ip.state.GetRulesCount(),
		FactsCount:       ip.state.GetFactsCount(),
		FilesParsedCount: len(ip.state.GetFilesParsed()),
		FilesParsed:      ip.state.GetFilesParsed(),
	}
}

// ParsingStatistics contains statistics about the parsing process
type ParsingStatistics struct {
	TypesCount       int      `json:"types_count"`
	RulesCount       int      `json:"rules_count"`
	FactsCount       int      `json:"facts_count"`
	FilesParsedCount int      `json:"files_parsed_count"`
	FilesParsed      []string `json:"files_parsed"`
}
