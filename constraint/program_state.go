// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"encoding/json"
	"fmt"
	"github.com/treivax/tsd/tsdio"
	"sort"
)

// ProgramState maintains the cumulative state of parsed types, rules, and facts
// across multiple file parsing operations.
//
// Thread-Safety: ProgramState is NOT thread-safe. Concurrent access must be
// synchronized externally using sync.Mutex or similar mechanisms.
// All methods that modify state (ParseAndMerge, ParseAndMergeContent, Reset)
// must not be called concurrently.
type ProgramState struct {
	types       map[string]*TypeDefinition
	rules       []*Expression
	facts       []*Fact
	filesParsed []string
	errors      []ValidationError
	ruleIDs     map[string]bool // Track used rule IDs (reset clears this)
}

// NewProgramState creates a new empty program state
func NewProgramState() *ProgramState {
	return &ProgramState{
		types:       make(map[string]*TypeDefinition),
		rules:       make([]*Expression, 0),
		facts:       make([]*Fact, 0),
		filesParsed: make([]string, 0),
		errors:      make([]ValidationError, 0),
		ruleIDs:     make(map[string]bool),
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

	return ps.mergeParseResult(result, filename)
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

	return ps.mergeParseResult(result, filename)
}

// mergeParseResult is the common implementation for parsing and merging
func (ps *ProgramState) mergeParseResult(result interface{}, filename string) error {
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
	ps.filesParsed = append(ps.filesParsed, filename)

	return nil
}

// mergeTypes merges new type definitions into the program state
func (ps *ProgramState) mergeTypes(newTypes []TypeDefinition, filename string) error {
	if ps == nil {
		return fmt.Errorf("ProgramState is nil")
	}
	if ps.types == nil {
		ps.types = make(map[string]*TypeDefinition)
	}
	if newTypes == nil {
		return nil
	}

	for _, typeDef := range newTypes {
		// Create a copy with filename information
		newType := &TypeDefinition{
			Type:   typeDef.Type,
			Name:   typeDef.Name,
			Fields: typeDef.Fields,
		}

		// Check if type already exists
		if existingType, exists := ps.types[typeDef.Name]; exists {
			// Validate compatibility
			if !ps.areTypesCompatible(existingType, newType) {
				return fmt.Errorf("type %s redefined incompatibly in %s",
					typeDef.Name, filename)
			}
			// Use the more detailed definition
			if len(newType.Fields) > len(existingType.Fields) {
				ps.types[typeDef.Name] = newType
			}
		} else {
			ps.types[typeDef.Name] = newType
		}
	}

	return nil
}

// mergeRules merges new rules and validates them against existing types
func (ps *ProgramState) mergeRules(newRules []Expression, filename string) error {
	if ps == nil {
		return fmt.Errorf("ProgramState is nil")
	}

	// Initialize ruleIDs map if nil
	if ps.ruleIDs == nil {
		ps.ruleIDs = make(map[string]bool)
	}

	for _, rule := range newRules {
		// Check for duplicate rule ID
		if rule.RuleId != "" {
			if ps.ruleIDs[rule.RuleId] {
				// Non-blocking error: record and skip this rule
				errMsg := fmt.Sprintf("rule ID '%s' already used, ignoring duplicate rule", rule.RuleId)
				ps.recordAndSkipError(filename, ErrorTypeRule, errMsg)
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
		if err := ps.validateRule(newRule, filename); err != nil {
			ps.recordAndSkipError(filename, ErrorTypeRule, err.Error())
			continue
		}

		// Mark this rule ID as used
		if newRule.RuleId != "" {
			ps.ruleIDs[newRule.RuleId] = true
		}

		ps.rules = append(ps.rules, newRule)
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
		if err := ps.validateFact(newFact, filename); err != nil {
			ps.recordAndSkipError(filename, ErrorTypeFact, err.Error())
			continue
		}

		ps.facts = append(ps.facts, newFact)
	}

	return nil
}

// recordAndSkipError records a validation error and logs it
func (ps *ProgramState) recordAndSkipError(filename, errorType, message string) {
	ps.errors = append(ps.errors, ValidationError{
		File:    filename,
		Type:    errorType,
		Message: message,
		Line:    0,
	})
	tsdio.Printf("⚠️  Skipping invalid %s in %s: %s\n", errorType, filename, message)
}

// validateRule validates a rule against existing type definitions
func (ps *ProgramState) validateRule(rule *Expression, filename string) error {
	// Extract variables from the rule
	variables := ps.extractRuleVariables(rule)

	// Validate each variable type exists
	for varName, varType := range variables {
		if _, exists := ps.types[varType]; !exists {
			return fmt.Errorf("variable %s references undefined type %s in %s", varName, varType, filename)
		}
	}

	// Validate field accesses in constraints
	if err := ps.validateRuleConstraints(rule, variables); err != nil {
		return err
	}

	// Validate field accesses in action
	if err := ps.validateRuleAction(rule, variables); err != nil {
		return err
	}

	return nil
}

// extractRuleVariables extracts all variables and their types from a rule
func (ps *ProgramState) extractRuleVariables(rule *Expression) map[string]string {
	variables := make(map[string]string)

	// New multi-pattern syntax (aggregation with joins)
	if len(rule.Patterns) > 0 {
		for _, pattern := range rule.Patterns {
			for _, variable := range pattern.Variables {
				// Skip aggregation variables (they don't have a dataType in the traditional sense)
				if variable.DataType != "" {
					variables[variable.Name] = variable.DataType
				}
			}
		}
	} else {
		// Old single-pattern syntax (backward compatibility)
		for _, variable := range rule.Set.Variables {
			variables[variable.Name] = variable.DataType
		}
	}

	return variables
}

// validateRuleConstraints validates field accesses in rule constraints
func (ps *ProgramState) validateRuleConstraints(rule *Expression, variables map[string]string) error {
	err := ps.validateFieldAccesses(rule.Constraints, variables)
	if err != nil {
		return fmt.Errorf("constraint validation failed: %w", err)
	}
	return nil
}

// validateRuleAction validates field accesses in rule action
func (ps *ProgramState) validateRuleAction(rule *Expression, variables map[string]string) error {
	if rule.Action == nil {
		return nil
	}

	// Convert Action struct to map[string]interface{} for validation
	actionBytes, err := json.Marshal(rule.Action)
	if err != nil {
		return fmt.Errorf("failed to serialize action: %w", err)
	}

	var actionMap map[string]interface{}
	if err := json.Unmarshal(actionBytes, &actionMap); err != nil {
		return fmt.Errorf("failed to deserialize action: %w", err)
	}

	err = ps.validateFieldAccesses(actionMap, variables)
	if err != nil {
		return fmt.Errorf("action validation failed: %w", err)
	}

	return nil
}

// validateFact validates a fact against existing type definitions
func (ps *ProgramState) validateFact(fact *Fact, filename string) error {
	// Check if type exists
	typeDef, exists := ps.types[fact.TypeName]
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
// Returns types sorted by name for deterministic output
func (ps *ProgramState) ToProgram() *Program {
	// Convert types and sort by name
	var types []TypeDefinition
	for _, typeDef := range ps.types {
		types = append(types, *typeDef)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

	// Convert rules
	var expressions []Expression
	for _, rule := range ps.rules {
		expressions = append(expressions, *rule)
	}

	// Convert facts
	var facts []Fact
	for _, fact := range ps.facts {
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
		if objType, hasType := v[JSONKeyType].(string); hasType && objType == ConstraintTypeFieldAccess {
			object, _ := extractMapStringValue(v, JSONKeyObject)
			field, _ := extractMapStringValue(v, JSONKeyField)

			// Validate the field access
			if varType, exists := variables[object]; exists {
				if typeDef, typeExists := ps.types[varType]; typeExists {
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
	case ValueTypeString:
		if value.Type != ValueTypeString && value.Type != ValueTypeIdentifier {
			return fmt.Errorf("expected string, got %s", value.Type)
		}
	case ValueTypeNumber:
		if value.Type != ValueTypeNumber {
			return fmt.Errorf("expected number, got %s", value.Type)
		}
	case ValueTypeBool:
		if value.Type != ValueTypeBoolean {
			return fmt.Errorf("expected boolean, got %s", value.Type)
		}
	}

	return nil
}

func extractMapStringValue(m map[string]interface{}, key string) (string, bool) {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}
