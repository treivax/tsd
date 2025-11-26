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
// and constraint consistency checks.
//
// Example:
//
//	ast, _ := ParseConstraint("rules.constraint", content)
//	err := ValidateConstraintProgram(ast)
func ValidateConstraintProgram(result interface{}) error {
	return ValidateProgram(result)
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

// ExtractFactsFromProgram extrait les faits d'un programme parsé et les convertit au format RETE
func ExtractFactsFromProgram(result interface{}) ([]map[string]interface{}, error) {
	// Convertir le résultat en structure Program
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("erreur conversion JSON: %v", err)
	}

	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %v", err)
	}

	// Convertir les faits au format RETE
	reteFacts := ConvertFactsToReteFormat(program)
	return reteFacts, nil
}

// ConvertResultToProgram convertit le résultat du parser en structure Program
func ConvertResultToProgram(result interface{}) (*Program, error) {
	// Convertir le résultat en structure Program
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("erreur conversion JSON: %v", err)
	}

	var program Program
	err = json.Unmarshal(jsonData, &program)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing JSON: %v", err)
	}

	return &program, nil
}

// ConvertToReteProgram convertit une structure Program des contraintes vers le format attendu par RETE
func ConvertToReteProgram(program *Program) interface{} {
	// Convertir les types en map[string]interface{} via JSON pour être compatible avec RETE
	typesInterface := make([]interface{}, len(program.Types))
	for i, typeDef := range program.Types {
		// Sérialiser vers JSON puis désérialiser en map
		jsonData, err := json.Marshal(typeDef)
		if err != nil {
			continue
		}
		var typeMap map[string]interface{}
		json.Unmarshal(jsonData, &typeMap)
		typesInterface[i] = typeMap
	}

	// Convertir les expressions en map[string]interface{} via JSON
	expressionsInterface := make([]interface{}, len(program.Expressions))
	for i, expr := range program.Expressions {
		// Sérialiser vers JSON puis désérialiser en map
		jsonData, err := json.Marshal(expr)
		if err != nil {
			continue
		}
		var exprMap map[string]interface{}
		json.Unmarshal(jsonData, &exprMap)
		expressionsInterface[i] = exprMap
	}

	// Créer une structure compatible avec le format attendu par RETE
	reteProgram := map[string]interface{}{
		"types":       typesInterface,
		"expressions": expressionsInterface,
	}
	return reteProgram
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
		TypesCount:       len(ip.state.Types),
		RulesCount:       len(ip.state.Rules),
		FactsCount:       len(ip.state.Facts),
		FilesParsedCount: len(ip.state.FilesParsed),
		FilesParsed:      ip.state.FilesParsed,
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
