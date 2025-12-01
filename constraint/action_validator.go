// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
	"strings"
)

// ActionValidator validates action calls against action definitions and type definitions.
type ActionValidator struct {
	actions map[string]*ActionDefinition
	types   map[string]*TypeDefinition
}

// NewActionValidator creates a new ActionValidator with the given action and type definitions.
func NewActionValidator(actions []ActionDefinition, types []TypeDefinition) *ActionValidator {
	av := &ActionValidator{
		actions: make(map[string]*ActionDefinition),
		types:   make(map[string]*TypeDefinition),
	}

	// Index actions by name
	for i := range actions {
		av.actions[actions[i].Name] = &actions[i]
	}

	// Index types by name
	for i := range types {
		av.types[types[i].Name] = &types[i]
	}

	return av
}

// ValidateActionCall validates that a job call matches its action definition.
// It checks:
// - The action exists
// - The number of arguments matches (considering optional parameters)
// - The argument types are compatible with parameter types
func (av *ActionValidator) ValidateActionCall(jobCall *JobCall, ruleVariables map[string]string) error {
	// Get action definition
	actionDef, exists := av.actions[jobCall.Name]
	if !exists {
		return fmt.Errorf("action '%s' is not defined", jobCall.Name)
	}

	// Count required and optional parameters
	requiredCount := 0
	optionalCount := 0
	for _, param := range actionDef.Parameters {
		if param.Optional || param.DefaultValue != nil {
			optionalCount++
		} else {
			requiredCount++
		}
	}

	totalParams := len(actionDef.Parameters)
	argCount := len(jobCall.Args)

	// Check argument count
	if argCount < requiredCount {
		return fmt.Errorf("action '%s' requires at least %d arguments, got %d",
			jobCall.Name, requiredCount, argCount)
	}
	if argCount > totalParams {
		return fmt.Errorf("action '%s' accepts at most %d arguments, got %d",
			jobCall.Name, totalParams, argCount)
	}

	// Validate each argument against its parameter type
	for i, arg := range jobCall.Args {
		param := actionDef.Parameters[i]

		// Get the type of the argument
		argType, err := av.inferArgumentType(arg, ruleVariables)
		if err != nil {
			return fmt.Errorf("error inferring type of argument %d for action '%s': %v",
				i+1, jobCall.Name, err)
		}

		// Check type compatibility
		if !av.isTypeCompatible(argType, param.Type) {
			return fmt.Errorf("type mismatch for parameter '%s' in action '%s': expected '%s', got '%s'",
				param.Name, jobCall.Name, param.Type, argType)
		}
	}

	return nil
}

// inferArgumentType infers the type of an argument expression.
func (av *ActionValidator) inferArgumentType(arg interface{}, ruleVariables map[string]string) (string, error) {
	switch v := arg.(type) {
	case map[string]interface{}:
		argType, ok := v["type"].(string)
		if !ok {
			return "", fmt.Errorf("argument missing type field")
		}

		switch argType {
		case "string", "stringLiteral":
			return "string", nil
		case "number", "numberLiteral":
			return "number", nil
		case "boolean", "booleanLiteral", "bool":
			return "bool", nil
		case "variable":
			// Look up the variable type from rule variables
			varName, ok := v["name"].(string)
			if !ok {
				return "", fmt.Errorf("variable missing name")
			}
			varType, exists := ruleVariables[varName]
			if !exists {
				return "", fmt.Errorf("variable '%s' not found in rule", varName)
			}
			return varType, nil
		case "fieldAccess":
			// For field access, we need to look up the object type and then the field type
			objName, ok := v["object"].(string)
			if !ok {
				return "", fmt.Errorf("fieldAccess missing object name")
			}
			fieldName, ok := v["field"].(string)
			if !ok {
				return "", fmt.Errorf("fieldAccess missing field name")
			}

			// Get the type of the object
			objType, exists := ruleVariables[objName]
			if !exists {
				return "", fmt.Errorf("object '%s' not found in rule", objName)
			}

			// Get the type definition
			typeDef, exists := av.types[objType]
			if !exists {
				return "", fmt.Errorf("type '%s' not found", objType)
			}

			// Find the field in the type
			for _, field := range typeDef.Fields {
				if field.Name == fieldName {
					return field.Type, nil
				}
			}

			return "", fmt.Errorf("field '%s' not found in type '%s'", fieldName, objType)
		case "binaryOp":
			// For binary operations, infer from operands (assume number for arithmetic)
			op, ok := v["operator"].(string)
			if !ok {
				return "", fmt.Errorf("binaryOp missing operator")
			}
			// Arithmetic operations return number
			if op == "+" || op == "-" || op == "*" || op == "/" {
				return "number", nil
			}
			// Comparison operations return bool
			if op == "==" || op == "!=" || op == "<" || op == ">" || op == "<=" || op == ">=" {
				return "bool", nil
			}
			return "", fmt.Errorf("unknown operator '%s'", op)
		case "functionCall":
			// Function calls - would need function signature database
			// For now, return string for string functions, number for numeric functions
			funcName, ok := v["name"].(string)
			if !ok {
				return "", fmt.Errorf("functionCall missing name")
			}
			return av.inferFunctionReturnType(funcName), nil
		default:
			return "", fmt.Errorf("unknown argument type: %s", argType)
		}
	default:
		return "", fmt.Errorf("unexpected argument structure")
	}
}

// inferFunctionReturnType returns the return type of a function.
func (av *ActionValidator) inferFunctionReturnType(funcName string) string {
	switch strings.ToUpper(funcName) {
	case "LENGTH":
		return "number"
	case "SUBSTRING", "UPPER", "LOWER", "TRIM":
		return "string"
	case "ABS", "ROUND", "FLOOR", "CEIL":
		return "number"
	default:
		return "string" // Default to string
	}
}

// isTypeCompatible checks if an argument type is compatible with a parameter type.
func (av *ActionValidator) isTypeCompatible(argType, paramType string) bool {
	// Exact match
	if argType == paramType {
		return true
	}

	// Check if paramType is a user-defined type
	if _, exists := av.types[paramType]; exists {
		// For user-defined types, argument must be exactly that type
		return argType == paramType
	}

	// Primitive type compatibility
	// (for now, require exact match for primitives)
	return false
}

// ValidateActionDefinitions validates action definitions.
// It checks:
// - Parameter types are valid (either primitive or user-defined)
// - Default values match parameter types
func (av *ActionValidator) ValidateActionDefinitions() []error {
	var errors []error

	for _, action := range av.actions {
		for i, param := range action.Parameters {
			// Check if parameter type is valid
			if !av.isValidParameterType(param.Type) {
				errors = append(errors, fmt.Errorf(
					"action '%s' parameter '%s': type '%s' is not defined",
					action.Name, param.Name, param.Type))
			}

			// If there's a default value, check type compatibility
			if param.DefaultValue != nil {
				defaultType := av.inferDefaultValueType(param.DefaultValue)
				if !av.isTypeCompatible(defaultType, param.Type) {
					errors = append(errors, fmt.Errorf(
						"action '%s' parameter '%s': default value type '%s' does not match parameter type '%s'",
						action.Name, param.Name, defaultType, param.Type))
				}
			}

			// Optional parameters with no default value after required parameters is a warning
			// but we'll allow it for flexibility
			if i > 0 && param.Optional && param.DefaultValue == nil {
				prevParam := action.Parameters[i-1]
				if !prevParam.Optional && prevParam.DefaultValue == nil {
					// This is actually fine - optional can come after required
					// Just noting the pattern
				}
			}
		}
	}

	return errors
}

// isValidParameterType checks if a parameter type is valid.
func (av *ActionValidator) isValidParameterType(paramType string) bool {
	// Check primitive types
	if paramType == "string" || paramType == "number" || paramType == "bool" {
		return true
	}

	// Check user-defined types
	_, exists := av.types[paramType]
	return exists
}

// inferDefaultValueType infers the type of a default value.
func (av *ActionValidator) inferDefaultValueType(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		valueType, ok := v["type"].(string)
		if !ok {
			return "unknown"
		}
		switch valueType {
		case "string", "stringLiteral":
			return "string"
		case "number", "numberLiteral":
			return "number"
		case "boolean", "booleanLiteral", "bool":
			return "bool"
		default:
			return "unknown"
		}
	case string:
		return "string"
	case float64, int, int64:
		return "number"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

// GetActionDefinition returns the action definition for a given action name.
func (av *ActionValidator) GetActionDefinition(name string) (*ActionDefinition, bool) {
	action, exists := av.actions[name]
	return action, exists
}

// GetTypeDefinition returns the type definition for a given type name.
func (av *ActionValidator) GetTypeDefinition(name string) (*TypeDefinition, bool) {
	typeDef, exists := av.types[name]
	return typeDef, exists
}
