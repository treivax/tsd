package constraint

import (
	"encoding/json"
	"fmt"
)

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
	// Créer une structure compatible avec le format attendu par RETE
	// La structure rete.Program a seulement Types et Expressions, pas de Facts
	reteProgram := map[string]interface{}{
		"types":       program.Types,
		"expressions": program.Expressions,
	}
	return reteProgram
}
