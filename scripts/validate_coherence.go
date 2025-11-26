// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// TypeDefinition represents a constraint type
type TypeDefinition struct {
	Name   string
	Fields []string
}

// RuleReference represents a field reference in a rule
type RuleReference struct {
	Variable string
	Type     string
	Field    string
	Line     int
	Rule     string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run validate_coherence.go <constraint_file>")
	}

	filename := os.Args[1]
	fmt.Printf("🔍 Validation cohérence type/champ: %s\n", filename)

	types, rules, err := parseConstraintFile(filename)
	if err != nil {
		log.Fatalf("❌ Erreur parsing: %v", err)
	}

	fmt.Printf("📋 Types trouvés: %d\n", len(types))
	fmt.Printf("📋 Règles trouvées: %d\n", len(rules))

	errors := validateReferences(types, rules)

	if len(errors) == 0 {
		fmt.Println("✅ Validation réussie: Aucune incohérence détectée")
		return
	}

	fmt.Printf("❌ %d erreur(s) détectée(s):\n", len(errors))
	for _, err := range errors {
		fmt.Printf("   Ligne %d: %s\n", err.Line, err)
	}

	// Proposer corrections automatiques
	fmt.Println("\n🔧 Suggestions de correction:")
	suggestFixes(errors, types)

	os.Exit(1)
}

func parseConstraintFile(filename string) (map[string]TypeDefinition, []RuleReference, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	types := make(map[string]TypeDefinition)
	var rules []RuleReference

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Regex patterns
	typePattern := regexp.MustCompile(`type\s+(\w+)\s*:\s*<([^>]+)>`)
	rulePattern := regexp.MustCompile(`\{([^}]+)\}\s*/.*`)
	fieldRefPattern := regexp.MustCompile(`(\w+)\.(\w+)`)

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "//") || line == "" {
			continue
		}

		// Parse type definitions
		if match := typePattern.FindStringSubmatch(line); match != nil {
			typeName := match[1]
			fieldsStr := match[2]

			fields := parseFields(fieldsStr)
			types[typeName] = TypeDefinition{
				Name:   typeName,
				Fields: fields,
			}
		}

		// Parse rule references
		if match := rulePattern.FindStringSubmatch(line); match != nil {
			variables := parseVariables(match[1])

			// Find field references in the entire line
			fieldRefs := fieldRefPattern.FindAllStringSubmatch(line, -1)
			for _, ref := range fieldRefs {
				varName := ref[1]
				fieldName := ref[2]

				typeName := ""
				if varType, exists := variables[varName]; exists {
					typeName = varType
				}

				rules = append(rules, RuleReference{
					Variable: varName,
					Type:     typeName,
					Field:    fieldName,
					Line:     lineNum,
					Rule:     line,
				})
			}
		}
	}

	return types, rules, scanner.Err()
}

func parseFields(fieldsStr string) []string {
	var fields []string
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if colonPos := strings.Index(part, ":"); colonPos != -1 {
			fieldName := strings.TrimSpace(part[:colonPos])
			fields = append(fields, fieldName)
		}
	}

	return fields
}

func parseVariables(varsStr string) map[string]string {
	variables := make(map[string]string)

	parts := strings.Split(varsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if colonPos := strings.Index(part, ":"); colonPos != -1 {
			varName := strings.TrimSpace(part[:colonPos])
			typeName := strings.TrimSpace(part[colonPos+1:])
			variables[varName] = typeName
		}
	}

	return variables
}

func validateReferences(types map[string]TypeDefinition, rules []RuleReference) []ValidationError {
	var errors []ValidationError

	for _, rule := range rules {
		if rule.Type == "" {
			continue // Skip if type couldn't be determined
		}

		typeDesc, exists := types[rule.Type]
		if !exists {
			errors = append(errors, ValidationError{
				Line:    rule.Line,
				Message: fmt.Sprintf("Type inexistant: %s", rule.Type),
				Rule:    rule.Rule,
			})
			continue
		}

		if !contains(typeDesc.Fields, rule.Field) {
			errors = append(errors, ValidationError{
				Line:    rule.Line,
				Message: fmt.Sprintf("Champ inexistant: %s.%s (type %s n'a que: %v)", rule.Variable, rule.Field, rule.Type, typeDesc.Fields),
				Rule:    rule.Rule,
			})
		}
	}

	return errors
}

func suggestFixes(errors []ValidationError, types map[string]TypeDefinition) {
	for _, err := range errors {
		if strings.Contains(err.Message, "Champ inexistant") {
			parts := strings.Split(err.Message, " ")
			if len(parts) >= 3 {
				fieldRef := parts[2] // e.g., "p.available"
				if dotPos := strings.Index(fieldRef, "."); dotPos != -1 {
					varName := fieldRef[:dotPos]
					wrongField := fieldRef[dotPos+1:]

					// Find similar fields
					for typeName, typeDef := range types {
						for _, field := range typeDef.Fields {
							if similarField(wrongField, field) {
								fmt.Printf("   → Ligne %d: Remplacer '%s.%s' par '%s.%s' (type %s)\n",
									err.Line, varName, wrongField, varName, field, typeName)
							}
						}
					}
				}
			}
		}
	}
}

type ValidationError struct {
	Line    int
	Message string
	Rule    string
}

func (e ValidationError) Error() string {
	return e.Message
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func similarField(field1, field2 string) bool {
	return levenshteinDistance(field1, field2) <= 2 && len(field1) > 2
}

func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min(matrix[i-1][j], matrix[i][j-1], matrix[i-1][j-1]) + 1
			}
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
