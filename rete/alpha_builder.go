// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// AlphaConditionBuilder aide à construire des conditions Alpha
type AlphaConditionBuilder struct{}

// NewAlphaConditionBuilder crée un nouveau constructeur de conditions
func NewAlphaConditionBuilder() *AlphaConditionBuilder {
	return &AlphaConditionBuilder{}
}

// FieldEquals crée une condition d'égalité sur un champ
func (acb *AlphaConditionBuilder) FieldEquals(variable, field string, value interface{}) interface{} {
	return map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": variable,
			"field":  field,
		},
		"right": acb.createLiteral(value),
	}
}

// FieldNotEquals crée une condition d'inégalité sur un champ
func (acb *AlphaConditionBuilder) FieldNotEquals(variable, field string, value interface{}) interface{} {
	return map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "!=",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": variable,
			"field":  field,
		},
		"right": acb.createLiteral(value),
	}
}

// FieldLessThan crée une condition de comparaison inférieure sur un champ
func (acb *AlphaConditionBuilder) FieldLessThan(variable, field string, value interface{}) interface{} {
	return map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "<",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": variable,
			"field":  field,
		},
		"right": acb.createLiteral(value),
	}
}

// FieldLessOrEqual crée une condition de comparaison inférieure ou égale sur un champ
func (acb *AlphaConditionBuilder) FieldLessOrEqual(variable, field string, value interface{}) interface{} {
	return map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "<=",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": variable,
			"field":  field,
		},
		"right": acb.createLiteral(value),
	}
}

// FieldGreaterThan crée une condition de comparaison supérieure sur un champ
func (acb *AlphaConditionBuilder) FieldGreaterThan(variable, field string, value interface{}) interface{} {
	return map[string]interface{}{
		"type":     "binaryOperation",
		"operator": ">",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": variable,
			"field":  field,
		},
		"right": acb.createLiteral(value),
	}
}

// FieldGreaterOrEqual crée une condition de comparaison supérieure ou égale sur un champ
func (acb *AlphaConditionBuilder) FieldGreaterOrEqual(variable, field string, value interface{}) interface{} {
	return map[string]interface{}{
		"type":     "binaryOperation",
		"operator": ">=",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": variable,
			"field":  field,
		},
		"right": acb.createLiteral(value),
	}
}

// And crée une condition logique AND
func (acb *AlphaConditionBuilder) And(left, right interface{}) interface{} {
	return map[string]interface{}{
		"type": "logicalExpression",
		"left": left,
		"operations": []interface{}{
			map[string]interface{}{
				"op":    "AND",
				"right": right,
			},
		},
	}
}

// Or crée une condition logique OR
func (acb *AlphaConditionBuilder) Or(left, right interface{}) interface{} {
	return map[string]interface{}{
		"type": "logicalExpression",
		"left": left,
		"operations": []interface{}{
			map[string]interface{}{
				"op":    "OR",
				"right": right,
			},
		},
	}
}

// AndMultiple crée une condition logique AND avec plusieurs conditions
func (acb *AlphaConditionBuilder) AndMultiple(conditions ...interface{}) interface{} {
	if len(conditions) == 0 {
		return acb.True()
	}
	if len(conditions) == 1 {
		return conditions[0]
	}

	operations := make([]interface{}, 0, len(conditions)-1)
	for _, condition := range conditions[1:] {
		operations = append(operations, map[string]interface{}{
			"op":    "AND",
			"right": condition,
		})
	}

	return map[string]interface{}{
		"type":       "logicalExpression",
		"left":       conditions[0],
		"operations": operations,
	}
}

// OrMultiple crée une condition logique OR avec plusieurs conditions
func (acb *AlphaConditionBuilder) OrMultiple(conditions ...interface{}) interface{} {
	if len(conditions) == 0 {
		return acb.False()
	}
	if len(conditions) == 1 {
		return conditions[0]
	}

	operations := make([]interface{}, 0, len(conditions)-1)
	for _, condition := range conditions[1:] {
		operations = append(operations, map[string]interface{}{
			"op":    "OR",
			"right": condition,
		})
	}

	return map[string]interface{}{
		"type":       "logicalExpression",
		"left":       conditions[0],
		"operations": operations,
	}
}

// True crée une condition toujours vraie
func (acb *AlphaConditionBuilder) True() interface{} {
	return map[string]interface{}{
		"type":  "booleanLiteral",
		"value": true,
	}
}

// False crée une condition toujours fausse
func (acb *AlphaConditionBuilder) False() interface{} {
	return map[string]interface{}{
		"type":  "booleanLiteral",
		"value": false,
	}
}

// FieldRange crée une condition de plage pour un champ (min <= field <= max)
func (acb *AlphaConditionBuilder) FieldRange(variable, field string, min, max interface{}) interface{} {
	minCondition := acb.FieldGreaterOrEqual(variable, field, min)
	maxCondition := acb.FieldLessOrEqual(variable, field, max)
	return acb.And(minCondition, maxCondition)
}

// FieldIn crée une condition de présence dans une liste de valeurs
func (acb *AlphaConditionBuilder) FieldIn(variable, field string, values ...interface{}) interface{} {
	if len(values) == 0 {
		return acb.False()
	}

	conditions := make([]interface{}, len(values))
	for i, value := range values {
		conditions[i] = acb.FieldEquals(variable, field, value)
	}

	return acb.OrMultiple(conditions...)
}

// FieldNotIn crée une condition d'absence dans une liste de valeurs
func (acb *AlphaConditionBuilder) FieldNotIn(variable, field string, values ...interface{}) interface{} {
	if len(values) == 0 {
		return acb.True()
	}

	conditions := make([]interface{}, len(values))
	for i, value := range values {
		conditions[i] = acb.FieldNotEquals(variable, field, value)
	}

	return acb.AndMultiple(conditions...)
}

// createLiteral crée un littéral typé
func (acb *AlphaConditionBuilder) createLiteral(value interface{}) map[string]interface{} {
	switch v := value.(type) {
	case string:
		return map[string]interface{}{
			"type":  "stringLiteral",
			"value": v,
		}
	case int:
		return map[string]interface{}{
			"type":  "numberLiteral",
			"value": float64(v),
		}
	case int32:
		return map[string]interface{}{
			"type":  "numberLiteral",
			"value": float64(v),
		}
	case int64:
		return map[string]interface{}{
			"type":  "numberLiteral",
			"value": float64(v),
		}
	case float32:
		return map[string]interface{}{
			"type":  "numberLiteral",
			"value": float64(v),
		}
	case float64:
		return map[string]interface{}{
			"type":  "numberLiteral",
			"value": v,
		}
	case bool:
		return map[string]interface{}{
			"type":  "booleanLiteral",
			"value": v,
		}
	default:
		// Fallback vers string
		return map[string]interface{}{
			"type":  "stringLiteral",
			"value": fmt.Sprintf("%v", v),
		}
	}
}

// CreateConstraintFromAST crée une condition à partir d'un AST de contrainte
func (acb *AlphaConditionBuilder) CreateConstraintFromAST(constraint interface{}) interface{} {
	// Si c'est déjà une map, la retourner telle quelle
	if constraintMap, ok := constraint.(map[string]interface{}); ok {
		return constraintMap
	}

	// Sinon, retourner tel quel (sera traité par l'évaluateur)
	return constraint
}
