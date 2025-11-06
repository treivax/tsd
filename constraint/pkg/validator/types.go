package validator

import (
	"fmt"
	"sync"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

// TypeRegistry implémente l'interface domain.TypeRegistry
type TypeRegistry struct {
	types map[string]domain.TypeDefinition
	mutex sync.RWMutex
}

// NewTypeRegistry crée un nouveau registre de types
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		types: make(map[string]domain.TypeDefinition),
		mutex: sync.RWMutex{},
	}
}

// RegisterType enregistre un nouveau type
func (tr *TypeRegistry) RegisterType(typeDef domain.TypeDefinition) error {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	if _, exists := tr.types[typeDef.Name]; exists {
		return domain.NewValidationError(
			fmt.Sprintf("type '%s' already exists", typeDef.Name),
			domain.Context{Type: typeDef.Name},
		)
	}

	tr.types[typeDef.Name] = typeDef
	return nil
}

// GetType récupère un type par son nom
func (tr *TypeRegistry) GetType(name string) (*domain.TypeDefinition, error) {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	if typeDef, exists := tr.types[name]; exists {
		return &typeDef, nil
	}

	return nil, domain.NewUnknownTypeError(
		name,
		domain.Context{Type: name},
	)
}

// HasType vérifie si un type existe
func (tr *TypeRegistry) HasType(name string) bool {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	_, exists := tr.types[name]
	return exists
}

// ListTypes retourne tous les types enregistrés
func (tr *TypeRegistry) ListTypes() []domain.TypeDefinition {
	tr.mutex.RLock()
	defer tr.mutex.RUnlock()

	types := make([]domain.TypeDefinition, 0, len(tr.types))
	for _, typeDef := range tr.types {
		types = append(types, typeDef)
	}
	return types
}

// GetTypeFields retourne les champs d'un type
func (tr *TypeRegistry) GetTypeFields(typeName string) (map[string]string, error) {
	typeDef, err := tr.GetType(typeName)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	for _, field := range typeDef.Fields {
		fields[field.Name] = field.Type
	}

	return fields, nil
}

// Clear supprime tous les types (utile pour les tests)
func (tr *TypeRegistry) Clear() {
	tr.mutex.Lock()
	defer tr.mutex.Unlock()

	tr.types = make(map[string]domain.TypeDefinition)
}

// TypeChecker implémente l'interface domain.TypeChecker
type TypeChecker struct {
	registry domain.TypeRegistry
}

// NewTypeChecker crée un nouveau vérificateur de types
func NewTypeChecker(registry domain.TypeRegistry) *TypeChecker {
	return &TypeChecker{
		registry: registry,
	}
}

// GetFieldType retourne le type d'un champ
func (tc *TypeChecker) GetFieldType(fieldAccess interface{}, variables []domain.TypedVariable, types []domain.TypeDefinition) (string, error) {
	// Cast vers FieldAccess
	fa, ok := fieldAccess.(*domain.FieldAccess)
	if !ok {
		// Essayer avec une map (format JSON)
		if faMap, ok := fieldAccess.(map[string]interface{}); ok {
			objectName, _ := faMap["object"].(string)
			fieldName, _ := faMap["field"].(string)

			fa = &domain.FieldAccess{
				Object: objectName,
				Field:  fieldName,
			}
		} else {
			return "", domain.NewValidationError(
				"invalid field access format",
				domain.Context{Value: fieldAccess},
			)
		}
	}

	// Trouver la variable
	var variableType string
	for _, variable := range variables {
		if variable.Name == fa.Object {
			variableType = variable.DataType
			break
		}
	}

	if variableType == "" {
		return "", domain.NewValidationError(
			fmt.Sprintf("variable '%s' not found", fa.Object),
			domain.Context{Variable: fa.Object},
		)
	}

	// Trouver le type de la variable
	typeDef, err := tc.registry.GetType(variableType)
	if err != nil {
		return "", err
	}

	// Trouver le champ dans le type
	field := typeDef.GetFieldByName(fa.Field)
	if field == nil {
		return "", domain.NewFieldNotFoundError(
			fa.Field,
			variableType,
			domain.Context{
				Field:    fa.Field,
				Type:     variableType,
				Variable: fa.Object,
			},
		)
	}

	return field.Type, nil
}

// GetValueType retourne le type d'une valeur
func (tc *TypeChecker) GetValueType(value interface{}) string {
	switch v := value.(type) {
	case bool:
		return "bool"
	case int, int8, int16, int32, int64:
		return "integer"
	case float32, float64:
		return "number"
	case string:
		return "string"
	case map[string]interface{}:
		// Format JSON, essayer de déterminer le type
		if valueType, ok := v["type"].(string); ok {
			switch valueType {
			case "booleanLiteral":
				return "bool"
			case "integerLiteral":
				return "integer"
			case "numberLiteral":
				return "number"
			case "stringLiteral":
				return "string"
			}
		}
		// Si on a une valeur directe
		if val, ok := v["value"]; ok {
			return tc.GetValueType(val)
		}
		return "unknown"
	default:
		return "unknown"
	}
}

// ValidateTypeCompatibility vérifie la compatibilité entre types
func (tc *TypeChecker) ValidateTypeCompatibility(leftType, rightType, operator string) error {
	// Opérateurs de comparaison
	comparisonOps := map[string]bool{
		"==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
	}

	// Opérateurs logiques
	logicalOps := map[string]bool{
		"AND": true, "OR": true, "NOT": true,
	}

	// Opérateurs arithmétiques
	arithmeticOps := map[string]bool{
		"+": true, "-": true, "*": true, "/": true, "%": true,
	}

	// Vérification basique des opérateurs
	if !domain.IsValidOperator(operator) {
		return domain.NewValidationError(
			fmt.Sprintf("invalid operator: %s", operator),
			domain.Context{
				Expected: "valid operator (==, !=, <, >, <=, >=, AND, OR, NOT, +, -, *, /, %)",
				Actual:   operator,
			},
		)
	}

	// Règles spécifiques par opérateur
	if comparisonOps[operator] {
		// Pour les comparaisons, les types doivent être compatibles
		return tc.validateComparisonTypes(leftType, rightType, operator)
	}

	if logicalOps[operator] {
		// Pour les opérations logiques, les opérandes doivent être booléens
		return tc.validateLogicalTypes(leftType, rightType, operator)
	}

	if arithmeticOps[operator] {
		// Pour l'arithmétique, les types doivent être numériques
		return tc.validateArithmeticTypes(leftType, rightType, operator)
	}

	return nil
}

// validateComparisonTypes valide les types pour les opérateurs de comparaison
func (tc *TypeChecker) validateComparisonTypes(leftType, rightType, operator string) error {
	// Égalité/inégalité : tous types compatibles si identiques
	if operator == "==" || operator == "!=" {
		if leftType != rightType {
			return domain.NewTypeMismatchError(
				leftType,
				rightType,
				domain.Context{},
			)
		}
		return nil
	}

	// Comparaisons ordinales : seulement pour les types numériques et strings
	orderableTypes := map[string]bool{
		"number": true, "integer": true, "string": true,
	}

	if !orderableTypes[leftType] || !orderableTypes[rightType] {
		return domain.NewValidationError(
			fmt.Sprintf("operator '%s' not supported for types '%s' and '%s'",
				operator, leftType, rightType),
			domain.Context{
				Expected: "number, integer, or string",
				Actual:   fmt.Sprintf("%s, %s", leftType, rightType),
			},
		)
	}

	if leftType != rightType {
		return domain.NewTypeMismatchError(
			leftType,
			rightType,
			domain.Context{},
		)
	}

	return nil
}

// validateLogicalTypes valide les types pour les opérateurs logiques
func (tc *TypeChecker) validateLogicalTypes(leftType, rightType, operator string) error {
	if operator == "NOT" {
		// NOT est unaire, vérifier seulement le type de gauche
		if leftType != "bool" {
			return domain.NewTypeMismatchError(
				"bool",
				leftType,
				domain.Context{},
			)
		}
		return nil
	}

	// AND/OR : les deux opérandes doivent être booléens
	if leftType != "bool" || rightType != "bool" {
		return domain.NewValidationError(
			fmt.Sprintf("logical operator '%s' requires boolean operands", operator),
			domain.Context{
				Expected: "bool, bool",
				Actual:   fmt.Sprintf("%s, %s", leftType, rightType),
			},
		)
	}

	return nil
}

// validateArithmeticTypes valide les types pour les opérateurs arithmétiques
func (tc *TypeChecker) validateArithmeticTypes(leftType, rightType, operator string) error {
	numericTypes := map[string]bool{
		"number": true, "integer": true,
	}

	if !numericTypes[leftType] || !numericTypes[rightType] {
		return domain.NewValidationError(
			fmt.Sprintf("arithmetic operator '%s' requires numeric operands", operator),
			domain.Context{
				Expected: "number or integer",
				Actual:   fmt.Sprintf("%s, %s", leftType, rightType),
			},
		)
	}

	return nil
}
