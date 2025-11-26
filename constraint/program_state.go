// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/json"
	"fmt"
)

// ProgramState maintains the cumulative state of parsed types, rules, and facts
// across multiple file parsing operations
type ProgramState struct {
	Types       map[string]*TypeDefinition `json:"types"`
	Rules       []*Expression              `json:"rules"`
	Facts       []*Fact                    `json:"facts"`
	FilesParsed []string                   `json:"files_parsed"`
	Errors      []ValidationError          `json:"errors"`
	RuleIDs     map[string]bool            `json:"-"` // Track used rule IDs (reset clears this)
}

// NewProgramState creates a new empty program state
func NewProgramState() *ProgramState {
	return &ProgramState{
		Types:       make(map[string]*TypeDefinition),
		Rules:       make([]*Expression, 0),
		Facts:       make([]*Fact, 0),
		FilesParsed: make([]string, 0),
		Errors:      make([]ValidationError, 0),
		RuleIDs:     make(map[string]bool),
	}
}

// ParseAndMerge parses a file and merges the results into the program state
// It validates that new rules and facts are compatible with existing type definitions
func (ps *ProgramState) ParseAndMerge(filename string) error {
	// Parse the file
	result, err := ParseConstraintFile(filename)
	if err != nil {
		return fmt.Errorf("error parsing file %s: %w", filename, err)
	}

	// Convert to program structure
	program, err := ConvertResultToProgram(result)
	if err != nil {
		return fmt.Errorf("error converting result for %s: %w", filename, err)
	}

	// Check for reset instructions first
	if len(program.Resets) > 0 {
		// Apply reset: clear all existing state including rule IDs
		ps.Reset()
		// Note: After reset, we continue to merge the new content from this file
	}

	// Merge types (new types or validate existing ones)
	err = ps.mergeTypes(program.Types, filename)
	if err != nil {
		return fmt.Errorf("error merging types from %s: %w", filename, err)
	}

	// Merge and validate rules
	err = ps.mergeRules(program.Expressions, filename)
	if err != nil {
		return fmt.Errorf("error merging rules from %s: %w", filename, err)
	}

	// Merge and validate facts
	err = ps.mergeFacts(program.Facts, filename)
	if err != nil {
		return fmt.Errorf("error merging facts from %s: %w", filename, err)
	}

	// Record the parsed file
	ps.FilesParsed = append(ps.FilesParsed, filename)

	return nil
}

// ParseAndMergeContent parses content from a string and merges the results into the program state
// It validates that new rules and facts are compatible with existing type definitions
func (ps *ProgramState) ParseAndMergeContent(content, filename string) error {
	if ps == nil {
		return fmt.Errorf("ProgramState is nil")
	}
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	if content == "" {
		return fmt.Errorf("content cannot be empty")
	}

	// Parse the content
	result, err := ParseConstraint(filename, []byte(content))
	if err != nil {
		return fmt.Errorf("error parsing content %s: %w", filename, err)
	}

	// Convert to program structure
	program, err := ConvertResultToProgram(result)
	if err != nil {
		return fmt.Errorf("error converting result for %s: %w", filename, err)
	}

	// Check for reset instructions first
	if len(program.Resets) > 0 {
		// Apply reset: clear all existing state including rule IDs
		ps.Reset()
		// Note: After reset, we continue to merge the new content from this content
	}

	// Merge types (new types or validate existing ones)
	err = ps.mergeTypes(program.Types, filename)
	if err != nil {
		return fmt.Errorf("error merging types from %s: %w", filename, err)
	}

	// Merge and validate rules
	err = ps.mergeRules(program.Expressions, filename)
	if err != nil {
		return fmt.Errorf("error merging rules from %s: %w", filename, err)
	}

	// Merge and validate facts
	err = ps.mergeFacts(program.Facts, filename)
	if err != nil {
		return fmt.Errorf("error merging facts from %s: %w", filename, err)
	}

	// Record the parsed file
	ps.FilesParsed = append(ps.FilesParsed, filename)

	return nil
}

// mergeTypes merges new type definitions into the program state
func (ps *ProgramState) mergeTypes(newTypes []TypeDefinition, filename string) error {
	if ps == nil {
		return fmt.Errorf("ProgramState is nil")
	}
	if ps.Types == nil {
		ps.Types = make(map[string]*TypeDefinition)
	}

	for _, typeDef := range newTypes {
		// Create a copy with filename information
		newType := &TypeDefinition{
			Type:   typeDef.Type,
			Name:   typeDef.Name,
			Fields: typeDef.Fields,
		}

		// Check if type already exists
		if existingType, exists := ps.Types[typeDef.Name]; exists {
			// Validate compatibility
			if !ps.areTypesCompatible(existingType, newType) {
				return fmt.Errorf("type %s redefined incompatibly in %s",
					typeDef.Name, filename)
			}
			// Use the more detailed definition
			if len(newType.Fields) > len(existingType.Fields) {
				ps.Types[typeDef.Name] = newType
			}
		} else {
			ps.Types[typeDef.Name] = newType
		}
	}

	return nil
}

// mergeRules merges new rules and validates them against existing types
func (ps *ProgramState) mergeRules(newRules []Expression, filename string) error {
	if ps == nil {
		return fmt.Errorf("ProgramState is nil")
	}

	// Initialize RuleIDs map if nil
	if ps.RuleIDs == nil {
		ps.RuleIDs = make(map[string]bool)
	}

	for _, rule := range newRules {
		// Check for duplicate rule ID
		if rule.RuleId != "" {
			if ps.RuleIDs[rule.RuleId] {
				// Non-blocking error: record and skip this rule
				errMsg := fmt.Sprintf("rule ID '%s' already used, ignoring duplicate rule", rule.RuleId)
				ps.Errors = append(ps.Errors, ValidationError{
					File:    filename,
					Type:    "rule",
					Message: errMsg,
					Line:    0,
				})
				fmt.Printf("⚠️  Skipping duplicate rule ID in %s: %s\n", filename, errMsg)
				continue
			}
		}

		// Create a copy of the rule
		newRule := &Expression{
			Type:        rule.Type,
			RuleId:      rule.RuleId,
			Set:         rule.Set,
			Constraints: rule.Constraints,
			Action:      rule.Action,
		}

		// Validate rule against existing types
		err := ps.validateRule(newRule, filename)
		if err != nil {
			// Non-blocking error: record and continue
			ps.Errors = append(ps.Errors, ValidationError{
				File:    filename,
				Type:    "rule",
				Message: err.Error(),
				Line:    0,
			})
			fmt.Printf("⚠️  Skipping invalid rule in %s: %v\n", filename, err)
			continue
		}

		// Mark this rule ID as used
		if newRule.RuleId != "" {
			ps.RuleIDs[newRule.RuleId] = true
		}

		ps.Rules = append(ps.Rules, newRule)
	}

	return nil
}

// mergeFacts merges new facts and validates them against existing types
func (ps *ProgramState) mergeFacts(newFacts []Fact, filename string) error {
	if ps == nil {
		return fmt.Errorf("ProgramState is nil")
	}

	for _, fact := range newFacts {
		// Create a copy of the fact
		newFact := &Fact{
			Type:     fact.Type,
			TypeName: fact.TypeName,
			Fields:   fact.Fields,
		}

		// Validate fact against existing types
		err := ps.validateFact(newFact, filename)
		if err != nil {
			// Non-blocking error: record and continue
			ps.Errors = append(ps.Errors, ValidationError{
				File:    filename,
				Type:    "fact",
				Message: err.Error(),
				Line:    0,
			})
			fmt.Printf("⚠️  Skipping invalid fact in %s: %v\n", filename, err)
			continue
		}

		ps.Facts = append(ps.Facts, newFact)
	}

	return nil
}

// validateRule validates a rule against existing type definitions
func (ps *ProgramState) validateRule(rule *Expression, filename string) error {
	// Extract variables from the set
	variables := make(map[string]string)
	for _, variable := range rule.Set.Variables {
		variables[variable.Name] = variable.DataType
	}

	// Validate each variable type exists
	for varName, varType := range variables {
		if _, exists := ps.Types[varType]; !exists {
			return fmt.Errorf("variable %s references undefined type %s in %s", varName, varType, filename)
		}
	}

	// Validate field accesses in constraints and action
	err := ps.validateFieldAccesses(rule.Constraints, variables)
	if err != nil {
		return fmt.Errorf("constraint validation failed: %w", err)
	}

	if rule.Action != nil {
		// Convert Action struct to map[string]interface{} for validation
		actionBytes, err := json.Marshal(rule.Action)
		if err != nil {
			return fmt.Errorf("failed to serialize action: %w", err)
		}
		var actionMap map[string]interface{}
		err = json.Unmarshal(actionBytes, &actionMap)
		if err != nil {
			return fmt.Errorf("failed to deserialize action: %w", err)
		}

		err = ps.validateFieldAccesses(actionMap, variables)
		if err != nil {
			return fmt.Errorf("action validation failed: %w", err)
		}
	}

	return nil
}

// validateFact validates a fact against existing type definitions
func (ps *ProgramState) validateFact(fact *Fact, filename string) error {
	// Check if type exists
	typeDef, exists := ps.Types[fact.TypeName]
	if !exists {
		return fmt.Errorf("fact references undefined type %s in %s", fact.TypeName, filename)
	}

	// Create a map of field definitions for easy lookup
	fieldDefs := make(map[string]Field)
	for _, field := range typeDef.Fields {
		fieldDefs[field.Name] = field
	}

	// Validate each field in the fact
	for _, factField := range fact.Fields {
		// Find field definition
		fieldDef, found := fieldDefs[factField.Name]
		if !found {
			return fmt.Errorf("fact contains undefined field %s for type %s in %s",
				factField.Name, fact.TypeName, filename)
		}

		// Validate field value type
		err := ps.validateFactValue(factField.Value, fieldDef.Type)
		if err != nil {
			return fmt.Errorf("field %s validation failed in %s: %w", factField.Name, filename, err)
		}
	}

	return nil
}

// ToProgram converts the program state back to a Program structure
func (ps *ProgramState) ToProgram() *Program {
	// Convert types
	var types []TypeDefinition
	for _, typeDef := range ps.Types {
		types = append(types, *typeDef)
	}

	// Convert rules
	var expressions []Expression
	for _, rule := range ps.Rules {
		expressions = append(expressions, *rule)
	}

	// Convert facts
	var facts []Fact
	for _, fact := range ps.Facts {
		facts = append(facts, *fact)
	}

	return &Program{
		Types:       types,
		Expressions: expressions,
		Facts:       facts,
	}
}

// Helper functions

func (ps *ProgramState) areTypesCompatible(type1, type2 *TypeDefinition) bool {
	if type1.Name != type2.Name {
		return false
	}

	// Create field maps for comparison
	fields1 := make(map[string]string)
	fields2 := make(map[string]string)

	for _, field := range type1.Fields {
		fields1[field.Name] = field.Type
	}
	for _, field := range type2.Fields {
		fields2[field.Name] = field.Type
	}

	// Check that all common fields have the same type
	for name, type1 := range fields1 {
		if type2, exists := fields2[name]; exists {
			if type1 != type2 {
				return false
			}
		}
	}

	return true
}

func (ps *ProgramState) validateFieldAccesses(data interface{}, variables map[string]string) error {
	// Recursively scan for fieldAccess patterns
	return ps.scanForFieldAccess(data, variables)
}

func (ps *ProgramState) scanForFieldAccess(data interface{}, variables map[string]string) error {
	switch v := data.(type) {
	case map[string]interface{}:
		// Check if this is a fieldAccess
		if objType, hasType := v["type"].(string); hasType && objType == "fieldAccess" {
			object := getStringValue(v, "object")
			field := getStringValue(v, "field")

			// Validate the field access
			if varType, exists := variables[object]; exists {
				if typeDef, typeExists := ps.Types[varType]; typeExists {
					fieldFound := false
					for _, fieldDef := range typeDef.Fields {
						if fieldDef.Name == field {
							fieldFound = true
							break
						}
					}
					if !fieldFound {
						return fmt.Errorf("field %s.%s not found in type %s", object, field, varType)
					}
				}
			}
		}

		// Recursively check all map values
		for _, value := range v {
			err := ps.scanForFieldAccess(value, variables)
			if err != nil {
				return err
			}
		}

	case []interface{}:
		// Recursively check all slice elements
		for _, item := range v {
			err := ps.scanForFieldAccess(item, variables)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ps *ProgramState) validateFactValue(value FactValue, expectedType string) error {
	// Basic type validation
	switch expectedType {
	case "string":
		if value.Type != "string" && value.Type != "identifier" {
			return fmt.Errorf("expected string, got %s", value.Type)
		}
	case "number":
		if value.Type != "number" {
			return fmt.Errorf("expected number, got %s", value.Type)
		}
	case "bool":
		if value.Type != "boolean" {
			return fmt.Errorf("expected boolean, got %s", value.Type)
		}
	}

	return nil
}

func getStringValue(m map[string]interface{}, key string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}
