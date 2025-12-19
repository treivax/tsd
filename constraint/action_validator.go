// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ActionValidator validates action calls against action definitions and type definitions.
type ActionValidator struct {
	actions          map[string]*ActionDefinition
	types            map[string]*TypeDefinition
	functionRegistry *FunctionRegistry
}

// NewActionValidator creates a new ActionValidator with the given action and type definitions.
func NewActionValidator(actions []ActionDefinition, types []TypeDefinition) *ActionValidator {
	av := &ActionValidator{
		actions:          make(map[string]*ActionDefinition),
		types:            make(map[string]*TypeDefinition),
		functionRegistry: DefaultFunctionRegistry,
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
	// Validate inputs
	if err := validateInputNotNil(map[string]interface{}{
		"jobCall":       jobCall,
		"ruleVariables": ruleVariables,
	}); err != nil {
		return err
	}

	// Get action definition
	actionDef, exists := av.actions[jobCall.Name]
	if !exists {
		return fmt.Errorf("action '%s' is not defined", sanitizeForLog(jobCall.Name, 100))
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
			sanitizeForLog(jobCall.Name, 100), requiredCount, argCount)
	}
	if argCount > totalParams {
		return fmt.Errorf("action '%s' accepts at most %d arguments, got %d",
			sanitizeForLog(jobCall.Name, 100), totalParams, argCount)
	}

	// Validate each argument against its parameter type
	for i, arg := range jobCall.Args {
		param := actionDef.Parameters[i]

		// Get the type of the argument
		argType, err := av.inferArgumentType(arg, ruleVariables, 0)
		if err != nil {
			return fmt.Errorf("error inferring type of argument %d for action '%s': %v",
				i+1, sanitizeForLog(jobCall.Name, 100), err)
		}

		// Check type compatibility
		if !av.isTypeCompatible(argType, param.Type) {
			return fmt.Errorf("type mismatch for parameter '%s' in action '%s': expected '%s', got '%s'",
				sanitizeForLog(param.Name, 50), sanitizeForLog(jobCall.Name, 100),
				sanitizeForLog(param.Type, 50), sanitizeForLog(argType, 50))
		}
	}

	return nil
}

// inferArgumentType infers the type of an argument expression with recursion depth tracking.
func (av *ActionValidator) inferArgumentType(arg interface{}, ruleVariables map[string]string, depth int) (string, error) {
	if depth > MaxValidationDepth {
		return "", fmt.Errorf("maximum validation depth exceeded (%d)", MaxValidationDepth)
	}

	argMap, ok := arg.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("unexpected argument structure")
	}

	argType, ok := argMap["type"].(string)
	if !ok {
		return "", fmt.Errorf("argument missing type field")
	}

	return av.inferComplexType(argType, argMap, ruleVariables)
}

// inferComplexType infers the type based on argument type and map
func (av *ActionValidator) inferComplexType(argType string, argMap map[string]interface{}, ruleVariables map[string]string) (string, error) {
	switch argType {
	case ValueTypeString, ArgTypeStringLiteral:
		return ValueTypeString, nil
	case ValueTypeNumber, ArgTypeNumberLiteral:
		return ValueTypeNumber, nil
	case ValueTypeBoolean, ArgTypeBoolLiteral, ValueTypeBool:
		return ValueTypeBool, nil
	case ValueTypeVariable:
		return av.inferVariableType(argMap, ruleVariables)
	case ConstraintTypeFieldAccess:
		return av.inferFieldAccessType(argMap, ruleVariables)
	case ArgTypeBinaryOp, ArgTypeBinaryOp2, ArgTypeBinaryOp3:
		return av.inferBinaryOpType(argMap)
	case ArgTypeFunctionCall:
		return av.inferFunctionCallType(argMap)
	case "inlineFact":
		return av.inferInlineFactType(argMap)
	default:
		return "", fmt.Errorf("unknown argument type: %s", sanitizeForLog(argType, 50))
	}
}

// inferVariableType infers the type of a variable reference
func (av *ActionValidator) inferVariableType(argMap map[string]interface{}, ruleVariables map[string]string) (string, error) {
	varName, ok := argMap["name"].(string)
	if !ok {
		return "", fmt.Errorf("variable missing name")
	}
	varType, exists := ruleVariables[varName]
	if !exists {
		return "", fmt.Errorf("variable '%s' not found in rule", sanitizeForLog(varName, 50))
	}
	return varType, nil
}

// inferFieldAccessType infers the type of a field access expression
func (av *ActionValidator) inferFieldAccessType(argMap map[string]interface{}, ruleVariables map[string]string) (string, error) {
	objName, ok := argMap["object"].(string)
	if !ok {
		return "", fmt.Errorf("fieldAccess missing object name")
	}
	fieldName, ok := argMap["field"].(string)
	if !ok {
		return "", fmt.Errorf("fieldAccess missing field name")
	}

	objType, exists := ruleVariables[objName]
	if !exists {
		return "", fmt.Errorf("object '%s' not found in rule", sanitizeForLog(objName, 50))
	}

	// Le champ 'id' est un champ spécial généré automatiquement, toujours de type string
	if fieldName == FieldNameID {
		return "string", nil
	}

	typeDef, exists := av.types[objType]
	if !exists {
		return "", fmt.Errorf("type '%s' not found", sanitizeForLog(objType, 50))
	}

	for _, field := range typeDef.Fields {
		if field.Name == fieldName {
			return field.Type, nil
		}
	}

	return "", fmt.Errorf("field '%s' not found in type '%s'",
		sanitizeForLog(fieldName, 50), sanitizeForLog(objType, 50))
}

// inferBinaryOpType infers the type of a binary operation
func (av *ActionValidator) inferBinaryOpType(argMap map[string]interface{}) (string, error) {
	op, ok := argMap["operator"].(string)
	if !ok {
		return "", fmt.Errorf("binaryOp missing operator")
	}

	if decoded, err := safeBase64Decode(op); err == nil {
		op = decoded
	}

	if isArithmeticOperator(op) {
		return ValueTypeNumber, nil
	}
	if isComparisonOperator(op) {
		return ValueTypeBool, nil
	}
	return "", fmt.Errorf("unknown operator '%s'", sanitizeForLog(op, 20))
}

// inferFunctionCallType infers the type of a function call
func (av *ActionValidator) inferFunctionCallType(argMap map[string]interface{}) (string, error) {
	funcName, ok := argMap["name"].(string)
	if !ok {
		return "", fmt.Errorf("functionCall missing name")
	}
	return av.functionRegistry.GetReturnType(funcName, ValueTypeString), nil
}

// inferInlineFactType infers the type of an inline fact creation.
// An inline fact has the format: TypeName(field: value, ...)
// The type returned is the TypeName, which should be a user-defined type.
func (av *ActionValidator) inferInlineFactType(argMap map[string]interface{}) (string, error) {
	typeName, ok := argMap["typeName"].(string)
	if !ok {
		return "", fmt.Errorf("inlineFact missing typeName")
	}

	// Vérifier que le type existe
	if _, exists := av.types[typeName]; !exists {
		return "", fmt.Errorf("type '%s' not defined", sanitizeForLog(typeName, 50))
	}

	// Le type d'un fait inline est le nom du type du fait
	return typeName, nil
}

// inferFunctionReturnType returns the return type of a function using the function registry.
// Deprecated: Use functionRegistry.GetReturnType() directly
func (av *ActionValidator) inferFunctionReturnType(funcName string) string {
	return av.functionRegistry.GetReturnType(funcName, ValueTypeString)
}

// isTypeCompatible checks if an argument type is compatible with a parameter type.
func (av *ActionValidator) isTypeCompatible(argType, paramType string) bool {
	// Special case: 'any' parameter type accepts any argument type
	if paramType == "any" {
		return true
	}

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
	// Check primitive types (including special type 'any')
	if paramType == "string" || paramType == "number" || paramType == "bool" || paramType == "any" {
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

// AddAction ajoute une nouvelle action au validator.
// Retourne une erreur si l'action existe déjà, avec un message spécifique
// si c'est une action par défaut qui ne peut pas être redéfinie.
func (av *ActionValidator) AddAction(action ActionDefinition) error {
	// Vérifier si l'action existe déjà
	if existing, exists := av.actions[action.Name]; exists {
		if existing.IsDefault {
			return fmt.Errorf("cannot redefine default action '%s' (default actions cannot be overridden)",
				sanitizeForLog(action.Name, 100))
		}
		return fmt.Errorf("action '%s' is already defined",
			sanitizeForLog(action.Name, 100))
	}

	// Ajouter l'action
	av.actions[action.Name] = &action
	return nil
}

// ValidateNonRedefinition vérifie qu'aucune action par défaut n'est redéfinie.
// Cette fonction est utile pour valider un lot d'actions avant de les ajouter.
func (av *ActionValidator) ValidateNonRedefinition(newActions []ActionDefinition) error {
	for _, newAction := range newActions {
		if existing, exists := av.actions[newAction.Name]; exists {
			if existing.IsDefault {
				return fmt.Errorf("cannot redefine default action '%s' (default actions cannot be overridden)",
					sanitizeForLog(newAction.Name, 100))
			}
			return fmt.Errorf("action '%s' is already defined",
				sanitizeForLog(newAction.Name, 100))
		}
	}
	return nil
}
