// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ProgramValidator valide un programme TSD complet.
type ProgramValidator struct {
	typeSystem    *TypeSystem
	factValidator *FactValidator
}

// NewProgramValidator crée un nouveau validateur de programme.
func NewProgramValidator() *ProgramValidator {
	return &ProgramValidator{}
}

// Validate valide un programme complet.
func (pv *ProgramValidator) Validate(program Program) error {
	if err := pv.validateTypeDefinitions(program.Types); err != nil {
		return fmt.Errorf("validation des types: %v", err)
	}

	pv.typeSystem = NewTypeSystem(program.Types)
	pv.factValidator = NewFactValidator(pv.typeSystem)

	if err := pv.typeSystem.ValidateCircularReferences(); err != nil {
		return fmt.Errorf("validation des types: %v", err)
	}

	if err := pv.validateFactAssignments(program.FactAssignments); err != nil {
		return fmt.Errorf("validation des affectations: %v", err)
	}

	if err := pv.validateFacts(program.Facts); err != nil {
		return fmt.Errorf("validation des faits: %v", err)
	}

	if err := pv.validateExpressions(program.Expressions); err != nil {
		return fmt.Errorf("validation des expressions: %v", err)
	}

	return nil
}

// validateTypeDefinitions valide toutes les définitions de types.
func (pv *ProgramValidator) validateTypeDefinitions(types []TypeDefinition) error {
	for i, typeDef := range types {
		if err := ValidateTypeDefinition(typeDef); err != nil {
			return fmt.Errorf("type %d ('%s'): %v", i+1, typeDef.Name, err)
		}
	}

	ts := NewTypeSystem(types)
	for i, typeDef := range types {
		for _, field := range typeDef.Fields {
			if !ts.IsPrimitiveType(field.Type) && !ts.IsUserDefinedType(field.Type) {
				return fmt.Errorf("type %d ('%s'): champ '%s' utilise un type non défini '%s'", i+1, typeDef.Name, field.Name, field.Type)
			}
		}
	}

	return nil
}

// validateFactAssignments valide les affectations de variables.
func (pv *ProgramValidator) validateFactAssignments(assignments []FactAssignment) error {
	for i, assignment := range assignments {
		if err := pv.factValidator.ValidateFact(assignment.Fact); err != nil {
			return fmt.Errorf("affectation %d (variable '%s'): %v", i+1, assignment.Variable, err)
		}

		if err := pv.typeSystem.RegisterVariable(assignment.Variable, assignment.Fact.TypeName); err != nil {
			return fmt.Errorf("affectation %d: %v", i+1, err)
		}
	}

	return nil
}

// validateFacts valide tous les faits.
func (pv *ProgramValidator) validateFacts(facts []Fact) error {
	for i, fact := range facts {
		if err := pv.factValidator.ValidateFact(fact); err != nil {
			return fmt.Errorf("fait %d: %v", i+1, err)
		}
	}

	return nil
}

// validateExpressions valide toutes les expressions/règles.
func (pv *ProgramValidator) validateExpressions(expressions []Expression) error {
	for i, expr := range expressions {
		if err := pv.validateExpression(expr); err != nil {
			return fmt.Errorf("expression %d (règle '%s'): %v", i+1, expr.RuleId, err)
		}
	}

	return nil
}

// validateExpression valide une expression.
func (pv *ProgramValidator) validateExpression(expr Expression) error {
	varTypes := make(map[string]string)

	patterns := expr.Patterns
	if len(patterns) == 0 && expr.Set.Type == "set" {
		patterns = []Set{expr.Set}
	}

	for _, pattern := range patterns {
		for _, variable := range pattern.Variables {
			if variable.Type == "typedVariable" {
				varTypes[variable.Name] = variable.DataType
			}
		}
	}

	if expr.Constraints != nil {
		if err := pv.validateConstraints(expr.Constraints, varTypes); err != nil {
			return err
		}
	}

	return nil
}

// validateConstraints valide les contraintes d'une expression.
func (pv *ProgramValidator) validateConstraints(constraints interface{}, varTypes map[string]string) error {
	switch c := constraints.(type) {
	case Constraint:
		return pv.validateConstraint(c, varTypes)

	case BinaryOperation:
		return pv.validateBinaryOperation(c, varTypes)

	case LogicalExpression:
		return pv.validateLogicalExpression(c, varTypes)

	case map[string]interface{}:
		return pv.validateConstraintMap(c, varTypes)

	default:
		return nil
	}
}

// validateConstraintMap valide une contrainte sous forme de map.
func (pv *ProgramValidator) validateConstraintMap(c map[string]interface{}, varTypes map[string]string) error {
	constraintType, ok := c["type"].(string)
	if !ok {
		return nil
	}

	switch constraintType {
	case ConstraintTypeComparison, ConstraintTypeBinaryOp:
		if left := c["left"]; left != nil {
			if err := pv.validateConstraints(left, varTypes); err != nil {
				return err
			}
		}
		if right := c["right"]; right != nil {
			if err := pv.validateConstraints(right, varTypes); err != nil {
				return err
			}
		}

		if operator, ok := c["operator"].(string); ok {
			return pv.validateComparisonFromMap(c, operator, varTypes)
		}

	case ConstraintTypeLogicalExpr:
		if left := c["left"]; left != nil {
			if err := pv.validateConstraints(left, varTypes); err != nil {
				return err
			}
		}
		if operations, ok := c["operations"].([]interface{}); ok {
			for _, op := range operations {
				if opMap, ok := op.(map[string]interface{}); ok {
					if right := opMap["right"]; right != nil {
						if err := pv.validateConstraints(right, varTypes); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// validateComparisonFromMap valide une comparaison à partir d'une map.
func (pv *ProgramValidator) validateComparisonFromMap(c map[string]interface{}, operator string, varTypes map[string]string) error {
	if left := c["left"]; left != nil {
		if right := c["right"]; right != nil {
			return pv.validateComparison(left, right, operator, varTypes)
		}
	}
	return nil
}

// validateConstraint valide une contrainte.
func (pv *ProgramValidator) validateConstraint(constraint Constraint, varTypes map[string]string) error {
	if constraint.Operator != "" {
		return pv.validateComparison(constraint.Left, constraint.Right, constraint.Operator, varTypes)
	}

	return nil
}

// validateBinaryOperation valide une opération binaire.
func (pv *ProgramValidator) validateBinaryOperation(op BinaryOperation, varTypes map[string]string) error {
	if op.Type == ConstraintTypeComparison {
		return pv.validateComparison(op.Left, op.Right, op.Operator, varTypes)
	}

	return nil
}

// validateLogicalExpression valide une expression logique.
func (pv *ProgramValidator) validateLogicalExpression(expr LogicalExpression, varTypes map[string]string) error {
	if err := pv.validateConstraints(expr.Left, varTypes); err != nil {
		return err
	}

	for _, op := range expr.Operations {
		if err := pv.validateConstraints(op.Right, varTypes); err != nil {
			return err
		}
	}

	return nil
}

// validateComparison valide une comparaison.
func (pv *ProgramValidator) validateComparison(left, right interface{}, operator string, varTypes map[string]string) error {
	leftType, err := pv.inferExpressionType(left, varTypes)
	if err != nil {
		return fmt.Errorf("expression gauche: %v", err)
	}

	rightType, err := pv.inferExpressionType(right, varTypes)
	if err != nil {
		return fmt.Errorf("expression droite: %v", err)
	}

	if !pv.typeSystem.AreTypesCompatible(leftType, rightType, operator) {
		return fmt.Errorf(
			"types incompatibles pour comparaison %s: '%s' et '%s'",
			operator,
			leftType,
			rightType,
		)
	}

	return nil
}

// inferExpressionType infère le type d'une expression.
func (pv *ProgramValidator) inferExpressionType(expr interface{}, varTypes map[string]string) (string, error) {
	switch e := expr.(type) {
	case FieldAccess:
		return pv.inferFieldAccessType(e, varTypes)

	case Variable:
		return pv.inferVariableType(e, varTypes)

	case StringLiteral:
		return ValueTypeString, nil

	case NumberLiteral:
		return ValueTypeNumber, nil

	case BooleanLiteral:
		return ValueTypeBool, nil

	case map[string]interface{}:
		return pv.inferMapExpressionType(e, varTypes)

	default:
		return "", fmt.Errorf("type d'expression non supporté: %T", expr)
	}
}

// inferMapExpressionType infère le type d'une expression sous forme de map.
func (pv *ProgramValidator) inferMapExpressionType(e map[string]interface{}, varTypes map[string]string) (string, error) {
	exprType, ok := e["type"].(string)
	if !ok {
		return "", fmt.Errorf("expression map sans type")
	}

	switch exprType {
	case ConstraintTypeFieldAccess:
		object, _ := e["object"].(string)
		field, _ := e["field"].(string)
		return pv.inferFieldAccessType(FieldAccess{Object: object, Field: field}, varTypes)

	case ValueTypeVariable:
		name, _ := e["name"].(string)
		return pv.inferVariableType(Variable{Name: name}, varTypes)

	case ValueTypeString, "stringLiteral":
		return ValueTypeString, nil

	case ValueTypeNumber, "numberLiteral":
		return ValueTypeNumber, nil

	case ValueTypeBoolean, "booleanLiteral", ValueTypeBool:
		return ValueTypeBool, nil

	default:
		return "", fmt.Errorf("type d'expression map non supporté: %s", exprType)
	}
}

// inferFieldAccessType infère le type d'un accès de champ.
func (pv *ProgramValidator) inferFieldAccessType(fa FieldAccess, varTypes map[string]string) (string, error) {
	varType, exists := varTypes[fa.Object]
	if !exists {
		return "", fmt.Errorf("variable '%s' non définie dans cette règle", fa.Object)
	}

	fieldType, err := pv.typeSystem.GetFieldType(varType, fa.Field)
	if err != nil {
		return "", err
	}

	return fieldType, nil
}

// inferVariableType infère le type d'une variable.
func (pv *ProgramValidator) inferVariableType(v Variable, varTypes map[string]string) (string, error) {
	varType, exists := varTypes[v.Name]
	if !exists {
		return "", fmt.Errorf("variable '%s' non définie dans cette règle", v.Name)
	}

	return varType, nil
}
