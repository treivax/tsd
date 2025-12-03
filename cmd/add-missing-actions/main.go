// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/treivax/tsd/constraint"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file.tsd> [<file2.tsd> ...]")
		fmt.Println("   or: go run main.go test/coverage/alpha/*.tsd")
		os.Exit(1)
	}

	for _, filePath := range os.Args[1:] {
		if err := processFile(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filePath, err)
		}
	}
}

func processFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)

	// Parse the file to get AST
	result, err := constraint.ParseConstraintFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	// Convert to Program structure
	program, err := constraint.ConvertResultToProgram(result)
	if err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	// Build a map of type definitions to know field types
	typeFields := make(map[string]map[string]string) // typeName -> fieldName -> fieldType
	for _, typeDecl := range program.Types {
		fields := make(map[string]string)
		for _, field := range typeDecl.Fields {
			fields[field.Name] = field.Type
		}
		typeFields[typeDecl.Name] = fields
	}

	// Find existing action definitions
	definedActions := make(map[string]bool)
	for _, action := range program.Actions {
		definedActions[action.Name] = true
	}

	// Find all action calls in rules
	// Use a custom parser instead of regex to handle nested parentheses
	matches := extractActionCalls(contentStr)

	if len(matches) == 0 {
		fmt.Printf("✓ %s: no actions to add\n", filepath.Base(filePath))
		return nil
	}

	// Analyze missing actions and their parameter types
	missingActions := make(map[string][]ActionParam)

	for _, match := range matches {
		actionName := match.Name
		if definedActions[actionName] {
			continue
		}

		argsStr := strings.TrimSpace(match.Args)
		if argsStr == "" {
			if _, exists := missingActions[actionName]; !exists {
				missingActions[actionName] = []ActionParam{}
			}
			continue
		}

		// Parse arguments to infer types
		args := strings.Split(argsStr, ",")
		params := make([]ActionParam, len(args))

		for i, arg := range args {
			arg = strings.TrimSpace(arg)
			paramType := inferArgumentType(arg, program, typeFields)
			params[i] = ActionParam{
				Name: fmt.Sprintf("arg%d", i+1),
				Type: paramType,
			}
		}

		// Store the signature with the most parameters
		if existing, exists := missingActions[actionName]; !exists || len(params) > len(existing) {
			missingActions[actionName] = params
		}
	}

	if len(missingActions) == 0 {
		fmt.Printf("✓ %s: all actions already defined\n", filepath.Base(filePath))
		return nil
	}

	// Generate action definitions
	var actionDefs strings.Builder
	actionDefs.WriteString("\n// Action definitions\n")

	for actionName, params := range missingActions {
		actionDefs.WriteString("action ")
		actionDefs.WriteString(actionName)
		actionDefs.WriteString("(")

		for i, param := range params {
			if i > 0 {
				actionDefs.WriteString(", ")
			}
			actionDefs.WriteString(param.Name)
			actionDefs.WriteString(": ")
			actionDefs.WriteString(param.Type)
		}

		actionDefs.WriteString(")\n")
	}
	actionDefs.WriteString("\n")

	// Insert action definitions after type definitions, before rules
	ruleRegex := regexp.MustCompile(`(?m)^rule\s+`)
	ruleLoc := ruleRegex.FindStringIndex(contentStr)

	var newContent string
	if ruleLoc != nil {
		insertPos := ruleLoc[0]
		newContent = contentStr[:insertPos] + actionDefs.String() + contentStr[insertPos:]
	} else {
		// No rules found, append at end
		newContent = contentStr + "\n" + actionDefs.String()
	}

	// Write back to file
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return err
	}

	fmt.Printf("✓ %s: added %d action(s)\n", filepath.Base(filePath), len(missingActions))
	for actionName, params := range missingActions {
		var paramStrs []string
		for _, p := range params {
			paramStrs = append(paramStrs, fmt.Sprintf("%s: %s", p.Name, p.Type))
		}
		fmt.Printf("  - %s(%s)\n", actionName, strings.Join(paramStrs, ", "))
	}

	return nil
}

type ActionParam struct {
	Name string
	Type string
}

// inferArgumentType tries to infer the type of an argument based on its syntax
func inferArgumentType(arg string, program *constraint.Program, typeFields map[string]map[string]string) string {
	arg = strings.TrimSpace(arg)

	// Check for arithmetic expressions (contains +, -, *, /)
	// These always return number type
	if containsArithmeticOperator(arg) {
		return "number"
	}

	// Check for field access (e.g., b.amount, p.price)
	if strings.Contains(arg, ".") {
		parts := strings.Split(arg, ".")
		if len(parts) == 2 {
			// Try to find the binding and its type
			varName := parts[0]
			fieldName := parts[1]

			// Look for the variable in rule patterns
			for _, expr := range program.Expressions {
				// Check in Set (backward compatibility)
				if len(expr.Set.Variables) > 0 {
					for _, v := range expr.Set.Variables {
						if v.Name == varName {
							typeName := v.DataType
							if fields, ok := typeFields[typeName]; ok {
								if fieldType, ok := fields[fieldName]; ok {
									return fieldType
								}
							}
						}
					}
				}

				// Check in Patterns (new format)
				if expr.Patterns != nil {
					for _, pattern := range expr.Patterns {
						for _, v := range pattern.Variables {
							if v.Name == varName {
								typeName := v.DataType
								if fields, ok := typeFields[typeName]; ok {
									if fieldType, ok := fields[fieldName]; ok {
										return fieldType
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Check for literal values
	if strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`) {
		return "string"
	}

	if arg == "true" || arg == "false" {
		return "bool"
	}

	// Check if it's a number
	if matched, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, arg); matched {
		return "number"
	}

	// Check for function calls (they typically return specific types)
	if strings.Contains(arg, "(") {
		funcName := strings.Split(arg, "(")[0]
		funcName = strings.TrimSpace(funcName)

		// Common functions and their return types
		switch strings.ToUpper(funcName) {
		case "ABS", "LEN", "LENGTH", "FLOOR", "CEIL", "ROUND", "MIN", "MAX":
			return "number"
		case "UPPER", "LOWER", "TRIM", "SUBSTRING":
			return "string"
		case "CONTAINS", "MATCHES", "LIKE":
			return "bool"
		}

		// If function operates on a field, infer from the inner expression
		innerStart := strings.Index(arg, "(")
		innerEnd := strings.LastIndex(arg, ")")
		if innerStart >= 0 && innerEnd > innerStart {
			innerArg := arg[innerStart+1 : innerEnd]
			return inferArgumentType(innerArg, program, typeFields)
		}
	}

	// Default to string if we can't determine the type
	return "string"
}

// containsArithmeticOperator checks if an expression contains arithmetic operators
// outside of string literals
func containsArithmeticOperator(expr string) bool {
	inString := false

	for i := 0; i < len(expr); i++ {
		ch := expr[i]

		// Track string literals
		if ch == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		// Check for arithmetic operators anywhere in the expression (including inside parentheses)
		// Parentheses in expressions like (a + b) or (x * y) still indicate arithmetic
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			// Make sure it's not a negative sign at the start
			if ch == '-' && i == 0 {
				continue
			}
			// Check if it's a binary operator (has non-whitespace on both sides)
			if i > 0 && i < len(expr)-1 {
				// Look back for non-whitespace
				hasLeft := false
				for j := i - 1; j >= 0; j-- {
					if expr[j] != ' ' && expr[j] != '\t' {
						hasLeft = true
						break
					}
				}
				// Look ahead for non-whitespace
				hasRight := false
				for j := i + 1; j < len(expr); j++ {
					if expr[j] != ' ' && expr[j] != '\t' {
						hasRight = true
						break
					}
				}
				if hasLeft && hasRight {
					return true
				}
			}
		}
	}

	return false
}

// ActionCall represents a parsed action call
type ActionCall struct {
	Name string
	Args string
}

// extractActionCalls parses action calls from content, handling nested parentheses
func extractActionCalls(content string) []ActionCall {
	var calls []ActionCall

	// Find all occurrences of "==> actionName("
	i := 0
	for i < len(content) {
		// Look for "==>"
		idx := strings.Index(content[i:], "==>")
		if idx == -1 {
			break
		}
		idx += i

		// Skip whitespace after "==>"
		j := idx + 3
		for j < len(content) && (content[j] == ' ' || content[j] == '\t') {
			j++
		}

		// Extract action name
		nameStart := j
		for j < len(content) && (isAlphaNum(content[j]) || content[j] == '_') {
			j++
		}
		if j == nameStart {
			i = j
			continue
		}
		actionName := content[nameStart:j]

		// Skip whitespace
		for j < len(content) && (content[j] == ' ' || content[j] == '\t') {
			j++
		}

		// Check for opening parenthesis
		if j >= len(content) || content[j] != '(' {
			i = j
			continue
		}
		j++ // skip '('

		// Extract arguments with proper parentheses matching
		argsStart := j
		parenDepth := 1
		inString := false

		for j < len(content) && parenDepth > 0 {
			ch := content[j]

			if ch == '"' && (j == 0 || content[j-1] != '\\') {
				inString = !inString
			} else if !inString {
				if ch == '(' {
					parenDepth++
				} else if ch == ')' {
					parenDepth--
				}
			}

			if parenDepth > 0 {
				j++
			}
		}

		if parenDepth == 0 {
			args := content[argsStart:j]
			calls = append(calls, ActionCall{
				Name: actionName,
				Args: args,
			})
		}

		i = j
	}

	return calls
}

// isAlphaNum checks if a character is alphanumeric
func isAlphaNum(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}
