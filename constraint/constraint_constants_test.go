// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import "testing"

func TestIsBinaryOperationType(t *testing.T) {
	t.Log("🧪 TEST isBinaryOperationType")
	t.Log("==============================")

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Cas nominaux - formats valides
		{
			name:     "format primaire binaryOp",
			input:    ArgTypeBinaryOp,
			expected: true,
		},
		{
			name:     "format legacy binaryOperation",
			input:    ArgTypeBinaryOp2,
			expected: true,
		},
		{
			name:     "format legacy binary_operation",
			input:    ArgTypeBinaryOp3,
			expected: true,
		},

		// Cas négatifs - types non binaires
		{
			name:     "type string literal",
			input:    ArgTypeStringLiteral,
			expected: false,
		},
		{
			name:     "type number literal",
			input:    ArgTypeNumberLiteral,
			expected: false,
		},
		{
			name:     "type function call",
			input:    ArgTypeFunctionCall,
			expected: false,
		},

		// Cas limites
		{
			name:     "chaîne vide",
			input:    "",
			expected: false,
		},
		{
			name:     "type invalide",
			input:    "invalidType",
			expected: false,
		},
		{
			name:     "casse différente",
			input:    "BinaryOp",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBinaryOperationType(tt.input)

			if result != tt.expected {
				t.Errorf("❌ isBinaryOperationType(%q) = %v, attendu %v",
					tt.input, result, tt.expected)
				return
			}

			t.Logf("✅ Test réussi: isBinaryOperationType(%q) = %v", tt.input, result)
		})
	}
}

func TestIsValidOperator(t *testing.T) {
	t.Log("🧪 TEST IsValidOperator")
	t.Log("=======================")

	tests := []struct {
		name     string
		operator string
		expected bool
	}{
		// Opérateurs arithmétiques
		{name: "addition", operator: OpAdd, expected: true},
		{name: "soustraction", operator: OpSub, expected: true},
		{name: "multiplication", operator: OpMul, expected: true},
		{name: "division", operator: OpDiv, expected: true},
		{name: "modulo", operator: OpMod, expected: true},

		// Opérateurs de comparaison
		{name: "égalité", operator: OpEq, expected: true},
		{name: "inégalité", operator: OpNeq, expected: true},
		{name: "inférieur", operator: OpLt, expected: true},
		{name: "supérieur", operator: OpGt, expected: true},
		{name: "inférieur ou égal", operator: OpLte, expected: true},
		{name: "supérieur ou égal", operator: OpGte, expected: true},

		// Opérateurs logiques
		{name: "AND logique", operator: OpAnd, expected: true},
		{name: "OR logique", operator: OpOr, expected: true},
		{name: "NOT logique", operator: OpNot, expected: true},

		// Cas négatifs
		{name: "opérateur invalide", operator: "???", expected: false},
		{name: "chaîne vide", operator: "", expected: false},
		{name: "opérateur inconnu", operator: "XOR", expected: false},
		{name: "casse incorrecte and", operator: "and", expected: false},
		{name: "casse incorrecte or", operator: "or", expected: false},
		{name: "opérateur bitwise", operator: "&", expected: false},
		{name: "opérateur bitwise OR", operator: "|", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidOperator(tt.operator)

			if result != tt.expected {
				t.Errorf("❌ IsValidOperator(%q) = %v, attendu %v",
					tt.operator, result, tt.expected)
				return
			}

			t.Logf("✅ Test réussi: IsValidOperator(%q) = %v", tt.operator, result)
		})
	}
}

func TestIsValidPrimitiveType(t *testing.T) {
	t.Log("🧪 TEST IsValidPrimitiveType")
	t.Log("============================")

	tests := []struct {
		name     string
		typeName string
		expected bool
	}{
		// Types primitifs valides
		{name: "type string", typeName: ValueTypeString, expected: true},
		{name: "type number", typeName: ValueTypeNumber, expected: true},
		{name: "type bool", typeName: ValueTypeBool, expected: true},
		{name: "type boolean", typeName: ValueTypeBoolean, expected: true},
		{name: "type integer", typeName: "integer", expected: true},

		// Cas négatifs
		{name: "type invalide", typeName: "invalid", expected: false},
		{name: "chaîne vide", typeName: "", expected: false},
		{name: "type object", typeName: "object", expected: false},
		{name: "type array", typeName: "array", expected: false},
		{name: "type null", typeName: "null", expected: false},
		{name: "type undefined", typeName: "undefined", expected: false},
		{name: "type float", typeName: "float", expected: false},
		{name: "type double", typeName: "double", expected: false},

		// Cas limites - casse
		{name: "String avec majuscule", typeName: "String", expected: false},
		{name: "NUMBER en majuscules", typeName: "NUMBER", expected: false},
		{name: "Boolean mixte", typeName: "Boolean", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPrimitiveType(tt.typeName)

			if result != tt.expected {
				t.Errorf("❌ IsValidPrimitiveType(%q) = %v, attendu %v",
					tt.typeName, result, tt.expected)
				return
			}

			t.Logf("✅ Test réussi: IsValidPrimitiveType(%q) = %v", tt.typeName, result)
		})
	}
}

func TestGetValidOperators(t *testing.T) {
	t.Log("🧪 TEST getValidOperators")
	t.Log("=========================")

	operators := getValidOperators()

	// Vérifier que tous les opérateurs attendus sont présents
	expectedOps := []string{
		OpAdd, OpSub, OpMul, OpDiv, OpMod,
		OpEq, OpNeq, OpLt, OpGt, OpLte, OpGte,
		OpAnd, OpOr, OpNot,
	}

	for _, op := range expectedOps {
		if !operators[op] {
			t.Errorf("❌ Opérateur attendu manquant: %q", op)
		}
	}

	t.Logf("✅ Tous les opérateurs attendus sont présents (%d opérateurs)", len(expectedOps))

	// Vérifier l'immutabilité (nouvelle map à chaque appel)
	operators2 := getValidOperators()
	if &operators == &operators2 {
		t.Error("❌ getValidOperators() devrait retourner une nouvelle map à chaque appel")
	} else {
		t.Log("✅ Immutabilité vérifiée: nouvelle map à chaque appel")
	}
}

func TestGetValidPrimitiveTypes(t *testing.T) {
	t.Log("🧪 TEST getValidPrimitiveTypes")
	t.Log("==============================")

	types := getValidPrimitiveTypes()

	// Vérifier que tous les types attendus sont présents
	expectedTypes := []string{
		ValueTypeString,
		ValueTypeNumber,
		ValueTypeBool,
		ValueTypeBoolean,
		"integer",
	}

	for _, typ := range expectedTypes {
		if !types[typ] {
			t.Errorf("❌ Type primitif attendu manquant: %q", typ)
		}
	}

	t.Logf("✅ Tous les types primitifs attendus sont présents (%d types)", len(expectedTypes))

	// Vérifier l'immutabilité (nouvelle map à chaque appel)
	types2 := getValidPrimitiveTypes()
	if &types == &types2 {
		t.Error("❌ getValidPrimitiveTypes() devrait retourner une nouvelle map à chaque appel")
	} else {
		t.Log("✅ Immutabilité vérifiée: nouvelle map à chaque appel")
	}
}

func TestBackwardCompatibilityConstants(t *testing.T) {
	t.Log("🧪 TEST Rétrocompatibilité des Constantes")
	t.Log("=========================================")

	t.Run("ValidOperators deprecated var", func(t *testing.T) {
		// Vérifier que la variable dépréciée existe et contient les bons opérateurs
		if ValidOperators == nil {
			t.Fatal("❌ ValidOperators ne devrait pas être nil")
		}

		if !ValidOperators[OpAdd] {
			t.Error("❌ ValidOperators devrait contenir OpAdd")
		}

		if !ValidOperators[OpEq] {
			t.Error("❌ ValidOperators devrait contenir OpEq")
		}

		t.Log("✅ ValidOperators (deprecated) fonctionne correctement")
	})

	t.Run("ValidPrimitiveTypes deprecated var", func(t *testing.T) {
		// Vérifier que la variable dépréciée existe et contient les bons types
		if ValidPrimitiveTypes == nil {
			t.Fatal("❌ ValidPrimitiveTypes ne devrait pas être nil")
		}

		if !ValidPrimitiveTypes[ValueTypeString] {
			t.Error("❌ ValidPrimitiveTypes devrait contenir ValueTypeString")
		}

		if !ValidPrimitiveTypes[ValueTypeNumber] {
			t.Error("❌ ValidPrimitiveTypes devrait contenir ValueTypeNumber")
		}

		t.Log("✅ ValidPrimitiveTypes (deprecated) fonctionne correctement")
	})
}

func TestConstantValues(t *testing.T) {
	t.Log("🧪 TEST Valeurs des Constantes")
	t.Log("==============================")

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		// Types de contraintes
		{"ConstraintTypeFieldAccess", ConstraintTypeFieldAccess, "fieldAccess"},
		{"ConstraintTypeComparison", ConstraintTypeComparison, "comparison"},
		{"ConstraintTypeLogicalExpr", ConstraintTypeLogicalExpr, "logicalExpr"},
		{"ConstraintTypeBinaryOp", ConstraintTypeBinaryOp, "binaryOp"},

		// Types de valeurs
		{"ValueTypeString", ValueTypeString, "string"},
		{"ValueTypeNumber", ValueTypeNumber, "number"},
		{"ValueTypeBoolean", ValueTypeBoolean, "boolean"},
		{"ValueTypeBool", ValueTypeBool, "bool"},
		{"ValueTypeIdentifier", ValueTypeIdentifier, "identifier"},
		{"ValueTypeVariable", ValueTypeVariable, "variable"},
		{"ValueTypeUnknown", ValueTypeUnknown, "unknown"},

		// Noms de champs spéciaux
		{"FieldNameID", FieldNameID, "id"},
		{"FieldNameReteType", FieldNameReteType, "reteType"},

		// Clés JSON
		{"JSONKeyType", JSONKeyType, "type"},
		{"JSONKeyFieldAccess", JSONKeyFieldAccess, "fieldAccess"},
		{"JSONKeyObject", JSONKeyObject, "object"},
		{"JSONKeyField", JSONKeyField, "field"},
		{"JSONKeyTypes", JSONKeyTypes, "types"},
		{"JSONKeyActions", JSONKeyActions, "actions"},
		{"JSONKeyExpressions", JSONKeyExpressions, "expressions"},
		{"JSONKeyRuleRemovals", JSONKeyRuleRemovals, "ruleRemovals"},

		// Types d'arguments
		{"ArgTypeStringLiteral", ArgTypeStringLiteral, "stringLiteral"},
		{"ArgTypeNumberLiteral", ArgTypeNumberLiteral, "numberLiteral"},
		{"ArgTypeBoolLiteral", ArgTypeBoolLiteral, "booleanLiteral"},
		{"ArgTypeFunctionCall", ArgTypeFunctionCall, "functionCall"},
		{"ArgTypeBinaryOp", ArgTypeBinaryOp, "binaryOp"},
		{"ArgTypeBinaryOp2", ArgTypeBinaryOp2, "binaryOperation"},
		{"ArgTypeBinaryOp3", ArgTypeBinaryOp3, "binary_operation"},

		// Opérateurs
		{"OpAdd", OpAdd, "+"},
		{"OpSub", OpSub, "-"},
		{"OpMul", OpMul, "*"},
		{"OpDiv", OpDiv, "/"},
		{"OpMod", OpMod, "%"},
		{"OpEq", OpEq, "=="},
		{"OpNeq", OpNeq, "!="},
		{"OpLt", OpLt, "<"},
		{"OpGt", OpGt, ">"},
		{"OpLte", OpLte, "<="},
		{"OpGte", OpGte, ">="},
		{"OpAnd", OpAnd, "AND"},
		{"OpOr", OpOr, "OR"},
		{"OpNot", OpNot, "NOT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("❌ %s = %q, attendu %q", tt.name, tt.actual, tt.expected)
				return
			}
			t.Logf("✅ %s = %q", tt.name, tt.actual)
		})
	}
}

func TestValidationLimits(t *testing.T) {
	t.Log("🧪 TEST Limites de Validation")
	t.Log("=============================")

	t.Run("MaxValidationDepth", func(t *testing.T) {
		if MaxValidationDepth != 100 {
			t.Errorf("❌ MaxValidationDepth = %d, attendu 100", MaxValidationDepth)
			return
		}
		t.Logf("✅ MaxValidationDepth = %d", MaxValidationDepth)
	})

	t.Run("MaxBase64DecodeSize", func(t *testing.T) {
		expectedSize := 1024 * 1024 // 1MB
		if MaxBase64DecodeSize != expectedSize {
			t.Errorf("❌ MaxBase64DecodeSize = %d, attendu %d", MaxBase64DecodeSize, expectedSize)
			return
		}
		t.Logf("✅ MaxBase64DecodeSize = %d (1MB)", MaxBase64DecodeSize)
	})
}
