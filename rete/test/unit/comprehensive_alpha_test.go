package rete

import (
	"math"
	"testing"
)

// TestAlphaConditionEvaluator_ComprehensiveCoverage teste tous les types d'expressions possibles
func TestAlphaConditionEvaluator_ComprehensiveCoverage(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	// Créer différents types de faits pour les tests
	integerFact := &Fact{
		ID:   "integer_fact",
		Type: "IntegerEvent",
		Fields: map[string]interface{}{
			"id":       42,
			"count":    100,
			"negative": -15,
			"zero":     0,
		},
	}

	floatFact := &Fact{
		ID:   "float_fact",
		Type: "FloatEvent",
		Fields: map[string]interface{}{
			"score":       95.5,
			"percentage":  0.85,
			"temperature": -2.5,
			"pi":          3.14159,
		},
	}

	stringFact := &Fact{
		ID:   "string_fact",
		Type: "StringEvent",
		Fields: map[string]interface{}{
			"name":        "TestEvent",
			"status":      "active",
			"category":    "urgent",
			"description": "",
		},
	}

	booleanFact := &Fact{
		ID:   "boolean_fact",
		Type: "BooleanEvent",
		Fields: map[string]interface{}{
			"active":    true,
			"validated": false,
			"confirmed": true,
		},
	}

	mixedFact := &Fact{
		ID:   "mixed_fact",
		Type: "MixedEvent",
		Fields: map[string]interface{}{
			"id":     123,
			"score":  87.3,
			"name":   "MixedTest",
			"active": true,
		},
	}

	tests := []struct {
		name        string
		fact        *Fact
		condition   interface{}
		variable    string
		expected    bool
		expectError bool
		description string
	}{
		// Tests des littéraux booléens
		{
			name:        "BooleanLiteral_True",
			fact:        integerFact,
			condition:   builder.True(),
			variable:    "event",
			expected:    true,
			description: "Littéral booléen toujours vrai",
		},
		{
			name:        "BooleanLiteral_False",
			fact:        integerFact,
			condition:   builder.False(),
			variable:    "event",
			expected:    false,
			description: "Littéral booléen toujours faux",
		},

		// Tests d'égalité avec différents types
		{
			name:        "IntegerEquality_True",
			fact:        integerFact,
			condition:   builder.FieldEquals("event", "id", 42),
			variable:    "event",
			expected:    true,
			description: "Égalité entier positive",
		},
		{
			name:        "IntegerEquality_False",
			fact:        integerFact,
			condition:   builder.FieldEquals("event", "id", 43),
			variable:    "event",
			expected:    false,
			description: "Égalité entier négative",
		},
		{
			name:        "FloatEquality_True",
			fact:        floatFact,
			condition:   builder.FieldEquals("event", "score", 95.5),
			variable:    "event",
			expected:    true,
			description: "Égalité float positive",
		},
		{
			name:        "FloatEquality_False",
			fact:        floatFact,
			condition:   builder.FieldEquals("event", "score", 95.6),
			variable:    "event",
			expected:    false,
			description: "Égalité float négative",
		},
		{
			name:        "StringEquality_True",
			fact:        stringFact,
			condition:   builder.FieldEquals("event", "name", "TestEvent"),
			variable:    "event",
			expected:    true,
			description: "Égalité string positive",
		},
		{
			name:        "StringEquality_False",
			fact:        stringFact,
			condition:   builder.FieldEquals("event", "name", "OtherEvent"),
			variable:    "event",
			expected:    false,
			description: "Égalité string négative",
		},
		{
			name:        "BooleanEquality_True",
			fact:        booleanFact,
			condition:   builder.FieldEquals("event", "active", true),
			variable:    "event",
			expected:    true,
			description: "Égalité booléenne positive",
		},
		{
			name:        "BooleanEquality_False",
			fact:        booleanFact,
			condition:   builder.FieldEquals("event", "active", false),
			variable:    "event",
			expected:    false,
			description: "Égalité booléenne négative",
		},

		// Tests d'inégalité
		{
			name:        "IntegerNotEquals_True",
			fact:        integerFact,
			condition:   builder.FieldNotEquals("event", "id", 43),
			variable:    "event",
			expected:    true,
			description: "Inégalité entier positive",
		},
		{
			name:        "IntegerNotEquals_False",
			fact:        integerFact,
			condition:   builder.FieldNotEquals("event", "id", 42),
			variable:    "event",
			expected:    false,
			description: "Inégalité entier négative",
		},
		{
			name:        "StringNotEquals_True",
			fact:        stringFact,
			condition:   builder.FieldNotEquals("event", "status", "inactive"),
			variable:    "event",
			expected:    true,
			description: "Inégalité string positive",
		},

		// Tests de comparaison numérique
		{
			name:        "IntegerLessThan_True",
			fact:        integerFact,
			condition:   builder.FieldLessThan("event", "id", 50),
			variable:    "event",
			expected:    true,
			description: "Comparaison entier < positive",
		},
		{
			name:        "IntegerLessThan_False",
			fact:        integerFact,
			condition:   builder.FieldLessThan("event", "id", 40),
			variable:    "event",
			expected:    false,
			description: "Comparaison entier < négative",
		},
		{
			name:        "FloatLessThan_True",
			fact:        floatFact,
			condition:   builder.FieldLessThan("event", "score", 96.0),
			variable:    "event",
			expected:    true,
			description: "Comparaison float < positive",
		},
		{
			name:        "IntegerLessOrEqual_True_Less",
			fact:        integerFact,
			condition:   builder.FieldLessOrEqual("event", "id", 50),
			variable:    "event",
			expected:    true,
			description: "Comparaison entier <= positive (moins)",
		},
		{
			name:        "IntegerLessOrEqual_True_Equal",
			fact:        integerFact,
			condition:   builder.FieldLessOrEqual("event", "id", 42),
			variable:    "event",
			expected:    true,
			description: "Comparaison entier <= positive (égal)",
		},
		{
			name:        "IntegerLessOrEqual_False",
			fact:        integerFact,
			condition:   builder.FieldLessOrEqual("event", "id", 40),
			variable:    "event",
			expected:    false,
			description: "Comparaison entier <= négative",
		},
		{
			name:        "IntegerGreaterThan_True",
			fact:        integerFact,
			condition:   builder.FieldGreaterThan("event", "count", 50),
			variable:    "event",
			expected:    true,
			description: "Comparaison entier > positive",
		},
		{
			name:        "IntegerGreaterThan_False",
			fact:        integerFact,
			condition:   builder.FieldGreaterThan("event", "count", 150),
			variable:    "event",
			expected:    false,
			description: "Comparaison entier > négative",
		},
		{
			name:        "FloatGreaterOrEqual_True_Greater",
			fact:        floatFact,
			condition:   builder.FieldGreaterOrEqual("event", "score", 90.0),
			variable:    "event",
			expected:    true,
			description: "Comparaison float >= positive (plus grand)",
		},
		{
			name:        "FloatGreaterOrEqual_True_Equal",
			fact:        floatFact,
			condition:   builder.FieldGreaterOrEqual("event", "score", 95.5),
			variable:    "event",
			expected:    true,
			description: "Comparaison float >= positive (égal)",
		},

		// Tests de comparaison de chaînes
		{
			name:        "StringLessThan_True",
			fact:        stringFact,
			condition:   builder.FieldLessThan("event", "status", "inactive"),
			variable:    "event",
			expected:    true,
			description: "Comparaison string < positive (active < inactive)",
		},
		{
			name:        "StringGreaterThan_True",
			fact:        stringFact,
			condition:   builder.FieldGreaterThan("event", "category", "normal"),
			variable:    "event",
			expected:    true,
			description: "Comparaison string > positive (urgent > normal)",
		},

		// Tests avec valeurs négatives et zéro
		{
			name:        "NegativeInteger_LessThan",
			fact:        integerFact,
			condition:   builder.FieldLessThan("event", "negative", 0),
			variable:    "event",
			expected:    true,
			description: "Comparaison entier négatif < 0",
		},
		{
			name:        "ZeroComparison_Equal",
			fact:        integerFact,
			condition:   builder.FieldEquals("event", "zero", 0),
			variable:    "event",
			expected:    true,
			description: "Comparaison avec zéro",
		},
		{
			name:        "NegativeFloat_GreaterThan",
			fact:        floatFact,
			condition:   builder.FieldGreaterThan("event", "temperature", -5.0),
			variable:    "event",
			expected:    true,
			description: "Comparaison float négatif > -5",
		},

		// Tests d'expressions logiques AND
		{
			name: "LogicalAnd_True_True",
			fact: mixedFact,
			condition: builder.And(
				builder.FieldEquals("event", "active", true),
				builder.FieldGreaterThan("event", "score", 80.0),
			),
			variable:    "event",
			expected:    true,
			description: "Expression logique AND (true && true)",
		},
		{
			name: "LogicalAnd_True_False",
			fact: mixedFact,
			condition: builder.And(
				builder.FieldEquals("event", "active", true),
				builder.FieldLessThan("event", "score", 50.0),
			),
			variable:    "event",
			expected:    false,
			description: "Expression logique AND (true && false)",
		},
		{
			name: "LogicalAnd_False_True",
			fact: mixedFact,
			condition: builder.And(
				builder.FieldEquals("event", "active", false),
				builder.FieldGreaterThan("event", "score", 80.0),
			),
			variable:    "event",
			expected:    false,
			description: "Expression logique AND (false && true)",
		},
		{
			name: "LogicalAnd_False_False",
			fact: mixedFact,
			condition: builder.And(
				builder.FieldEquals("event", "active", false),
				builder.FieldLessThan("event", "score", 50.0),
			),
			variable:    "event",
			expected:    false,
			description: "Expression logique AND (false && false)",
		},

		// Tests d'expressions logiques OR
		{
			name: "LogicalOr_True_True",
			fact: mixedFact,
			condition: builder.Or(
				builder.FieldEquals("event", "active", true),
				builder.FieldGreaterThan("event", "score", 80.0),
			),
			variable:    "event",
			expected:    true,
			description: "Expression logique OR (true || true)",
		},
		{
			name: "LogicalOr_True_False",
			fact: mixedFact,
			condition: builder.Or(
				builder.FieldEquals("event", "active", true),
				builder.FieldLessThan("event", "score", 50.0),
			),
			variable:    "event",
			expected:    true,
			description: "Expression logique OR (true || false)",
		},
		{
			name: "LogicalOr_False_True",
			fact: mixedFact,
			condition: builder.Or(
				builder.FieldEquals("event", "active", false),
				builder.FieldGreaterThan("event", "score", 80.0),
			),
			variable:    "event",
			expected:    true,
			description: "Expression logique OR (false || true)",
		},
		{
			name: "LogicalOr_False_False",
			fact: mixedFact,
			condition: builder.Or(
				builder.FieldEquals("event", "active", false),
				builder.FieldLessThan("event", "score", 50.0),
			),
			variable:    "event",
			expected:    false,
			description: "Expression logique OR (false || false)",
		},

		// Tests d'expressions imbriquées complexes
		{
			name: "Complex_Nested_And_Or",
			fact: mixedFact,
			condition: builder.And(
				builder.Or(
					builder.FieldEquals("event", "name", "MixedTest"),
					builder.FieldEquals("event", "name", "OtherTest"),
				),
				builder.FieldGreaterThan("event", "id", 100),
			),
			variable:    "event",
			expected:    true,
			description: "Expression complexe imbriquée ((name == 'MixedTest' || name == 'OtherTest') && id > 100)",
		},
		{
			name: "Complex_Multiple_Conditions",
			fact: mixedFact,
			condition: builder.Or(
				builder.And(
					builder.FieldEquals("event", "active", true),
					builder.FieldGreaterThan("event", "score", 85.0),
				),
				builder.And(
					builder.FieldEquals("event", "active", false),
					builder.FieldLessThan("event", "id", 50),
				),
			),
			variable:    "event",
			expected:    true,
			description: "Expression complexe avec multiples conditions",
		},

		// Tests avec des chaînes vides et valeurs spéciales
		{
			name:        "EmptyString_Equality",
			fact:        stringFact,
			condition:   builder.FieldEquals("event", "description", ""),
			variable:    "event",
			expected:    true,
			description: "Égalité avec chaîne vide",
		},
		{
			name:        "EmptyString_NotEquals",
			fact:        stringFact,
			condition:   builder.FieldNotEquals("event", "name", ""),
			variable:    "event",
			expected:    true,
			description: "Inégalité avec chaîne vide",
		},

		// Tests de conversion de types (int vers float)
		{
			name:        "TypeConversion_Int_To_Float",
			fact:        mixedFact,
			condition:   builder.FieldGreaterThan("event", "id", 120.5),
			variable:    "event",
			expected:    true,
			description: "Conversion automatique int vers float dans comparaison",
		},
	}

	// Exécuter tous les tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Nettoyer les liaisons avant chaque test
			evaluator.ClearBindings()

			result, err := evaluator.EvaluateCondition(tt.condition, tt.fact, tt.variable)

			if tt.expectError {
				if err == nil {
					t.Errorf("Attendait une erreur mais n'en a pas reçu")
				}
				return
			}

			if err != nil {
				t.Errorf("Erreur inattendue: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Résultat incorrect. Attendu: %v, Reçu: %v", tt.expected, result)
				t.Logf("Description: %s", tt.description)
				t.Logf("Condition: %+v", tt.condition)
				t.Logf("Fait: %+v", tt.fact)
			}
		})
	}
}

// TestAlphaConditionEvaluator_ExtendedErrorCases teste les cas d'erreur étendus
func TestAlphaConditionEvaluator_ExtendedErrorCases(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	// Fait de test avec des valeurs incompatibles pour certaines comparaisons
	testFact := &Fact{
		ID:   "error_test_fact",
		Type: "ErrorTestEvent",
		Fields: map[string]interface{}{
			"string_field": "test",
			"int_field":    42,
			"bool_field":   true,
		},
	}

	errorTests := []struct {
		name        string
		condition   interface{}
		fact        *Fact
		variable    string
		description string
	}{
		{
			name:        "NonExistent_Field",
			condition:   builder.FieldEquals("event", "nonexistent_field", "value"),
			fact:        testFact,
			variable:    "event",
			description: "Accès à un champ qui n'existe pas",
		},
		{
			name:        "Invalid_Expression_Type",
			condition:   "invalid_expression",
			fact:        testFact,
			variable:    "event",
			description: "Type d'expression non supporté",
		},
		{
			name:        "Incompatible_Type_Comparison_String_Int",
			condition:   builder.FieldLessThan("event", "string_field", 42),
			fact:        testFact,
			variable:    "event",
			description: "Comparaison incompatible entre string et int",
		},
		{
			name:        "Incompatible_Type_Comparison_Bool_Int",
			condition:   builder.FieldGreaterThan("event", "bool_field", 10),
			fact:        testFact,
			variable:    "event",
			description: "Comparaison incompatible entre bool et int",
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator.ClearBindings()

			_, err := evaluator.EvaluateCondition(tt.condition, tt.fact, tt.variable)

			if err == nil {
				t.Errorf("Attendait une erreur pour %s mais n'en a pas reçu", tt.description)
			} else {
				t.Logf("Erreur attendue reçue: %v", err)
			}
		})
	}
}

// TestAlphaConditionEvaluator_EdgeCases teste les cas limites
func TestAlphaConditionEvaluator_EdgeCases(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	// Tests avec des valeurs limites
	edgeFact := &Fact{
		ID:   "edge_case_fact",
		Type: "EdgeCaseEvent",
		Fields: map[string]interface{}{
			"max_int":    math.MaxInt64,
			"min_int":    math.MinInt64,
			"max_float":  math.MaxFloat64,
			"min_float":  -math.MaxFloat64,
			"inf_float":  math.Inf(1),
			"neg_inf":    math.Inf(-1),
			"nan_float":  math.NaN(),
			"zero_int":   0,
			"zero_float": 0.0,
		},
	}

	edgeTests := []struct {
		name        string
		condition   interface{}
		expected    bool
		expectError bool
		description string
	}{
		{
			name:        "MaxInt64_Comparison",
			condition:   builder.FieldGreaterThan("event", "max_int", 1000000),
			expected:    true,
			description: "Comparaison avec MaxInt64",
		},
		{
			name:        "MinInt64_Comparison",
			condition:   builder.FieldLessThan("event", "min_int", -1000000),
			expected:    true,
			description: "Comparaison avec MinInt64",
		},
		{
			name:        "MaxFloat64_Comparison",
			condition:   builder.FieldGreaterThan("event", "max_float", 1e308),
			expected:    true,
			description: "Comparaison avec MaxFloat64",
		},
		{
			name:        "Zero_Int_Float_Equality",
			condition:   builder.FieldEquals("event", "zero_int", 0.0),
			expected:    true,
			description: "Égalité entre zéro entier et zéro float",
		},
		{
			name:        "Infinity_Comparison",
			condition:   builder.FieldGreaterThan("event", "inf_float", 1e308),
			expected:    true,
			description: "Comparaison avec l'infini positif",
		},
		{
			name:        "Negative_Infinity_Comparison",
			condition:   builder.FieldLessThan("event", "neg_inf", -1e308),
			expected:    true,
			description: "Comparaison avec l'infini négatif",
		},
	}

	for _, tt := range edgeTests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator.ClearBindings()

			result, err := evaluator.EvaluateCondition(tt.condition, edgeFact, "event")

			if tt.expectError {
				if err == nil {
					t.Errorf("Attendait une erreur mais n'en a pas reçu")
				}
				return
			}

			if err != nil {
				t.Errorf("Erreur inattendue: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Résultat incorrect. Attendu: %v, Reçu: %v", tt.expected, result)
				t.Logf("Description: %s", tt.description)
			}
		})
	}
}

// TestAlphaConditionBuilder_AllMethods teste toutes les méthodes du builder
func TestAlphaConditionBuilder_AllMethods(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	// Test que toutes les méthodes du builder produisent des structures valides
	methodTests := []struct {
		name      string
		condition interface{}
	}{
		{"True", builder.True()},
		{"False", builder.False()},
		{"FieldEquals", builder.FieldEquals("var", "field", "value")},
		{"FieldNotEquals", builder.FieldNotEquals("var", "field", "value")},
		{"FieldLessThan", builder.FieldLessThan("var", "field", 10)},
		{"FieldLessOrEqual", builder.FieldLessOrEqual("var", "field", 10)},
		{"FieldGreaterThan", builder.FieldGreaterThan("var", "field", 10)},
		{"FieldGreaterOrEqual", builder.FieldGreaterOrEqual("var", "field", 10)},
		{"FieldRange", builder.FieldRange("var", "field", 5, 15)},
		{"And", builder.And(builder.True(), builder.False())},
		{"Or", builder.Or(builder.True(), builder.False())},
		{"AndMultiple", builder.AndMultiple(builder.True(), builder.False(), builder.True())},
		{"OrMultiple", builder.OrMultiple(builder.True(), builder.False(), builder.True())},
	}

	for _, tt := range methodTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.condition == nil {
				t.Errorf("La méthode %s a retourné nil", tt.name)
			}

			// Vérifier que la condition peut être sérialisée (structure valide)
			switch cond := tt.condition.(type) {
			case map[string]interface{}:
				if cond["type"] == nil {
					t.Errorf("La condition %s n'a pas de type défini", tt.name)
				}
			default:
				// Pour les types de base (bool, etc.), c'est acceptable
			}
		})
	}
}

// TestAlphaConditionEvaluator_VariableBindings teste la gestion des liaisons de variables
func TestAlphaConditionEvaluator_VariableBindings(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	fact1 := &Fact{ID: "fact1", Type: "Event1", Fields: map[string]interface{}{"value": 10}}
	fact2 := &Fact{ID: "fact2", Type: "Event2", Fields: map[string]interface{}{"value": 20}}

	// Test liaison de variable
	condition := builder.FieldEquals("event", "value", 10)
	result, err := evaluator.EvaluateCondition(condition, fact1, "event")

	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if !result {
		t.Errorf("La condition aurait dû être vraie")
	}

	// Vérifier que la variable est liée
	bindings := evaluator.GetBindings()
	if len(bindings) != 1 {
		t.Errorf("Attendait 1 liaison, reçu %d", len(bindings))
	}

	if bindings["event"] != fact1 {
		t.Errorf("La variable 'event' devrait être liée à fact1")
	}

	// Test changement de liaison
	result, err = evaluator.EvaluateCondition(condition, fact2, "event")
	if err != nil {
		t.Fatalf("Erreur inattendue: %v", err)
	}

	if result {
		t.Errorf("La condition aurait dû être fausse pour fact2")
	}

	// Vérifier que la liaison a changé
	bindings = evaluator.GetBindings()
	if bindings["event"] != fact2 {
		t.Errorf("La variable 'event' devrait maintenant être liée à fact2")
	}

	// Test nettoyage des liaisons
	evaluator.ClearBindings()
	bindings = evaluator.GetBindings()
	if len(bindings) != 0 {
		t.Errorf("Toutes les liaisons auraient dû être effacées")
	}
}

// TestAlphaConditionEvaluator_Integration teste l'intégration avec les nœuds Alpha
func TestAlphaConditionEvaluator_Integration(t *testing.T) {
	builder := NewAlphaConditionBuilder()

	// Créer un mock storage pour le test
	mockStorage := &MockStorage{}

	// Créer un nœud Alpha avec une condition
	condition := builder.And(
		builder.FieldEquals("event", "type", "Alert"),
		builder.FieldGreaterThan("event", "priority", 5),
	)

	alphaNode := NewAlphaNode("AlertNode", condition, "event", mockStorage)

	// Créer des faits de test
	alertFact := &Fact{
		ID:   "alert1",
		Type: "Event",
		Fields: map[string]interface{}{
			"type":     "Alert",
			"priority": 8,
			"message":  "High priority alert",
		},
	}

	normalFact := &Fact{
		ID:   "normal1",
		Type: "Event",
		Fields: map[string]interface{}{
			"type":     "Info",
			"priority": 3,
			"message":  "Normal message",
		},
	}

	lowPriorityAlert := &Fact{
		ID:   "alert2",
		Type: "Event",
		Fields: map[string]interface{}{
			"type":     "Alert",
			"priority": 3,
			"message":  "Low priority alert",
		},
	}

	// Compter les activations avec un mock successor
	activationCount := 0
	activateLeftCount := 0

	mockSuccessor := &MockBetaNode{
		id: "mock_successor",
		activateFunc: func(fact *Fact, token *Token) {
			activationCount++
		},
		activateLeftFunc: func(token *Token) {
			activateLeftCount++
		},
	}
	alphaNode.AddChild(mockSuccessor)

	// Tester les activations
	err := alphaNode.ActivateRight(alertFact)
	if err != nil {
		t.Errorf("Erreur lors de l'activation avec alertFact: %v", err)
	}

	err = alphaNode.ActivateRight(normalFact)
	if err != nil {
		t.Errorf("Erreur lors de l'activation avec normalFact: %v", err)
	}

	err = alphaNode.ActivateRight(lowPriorityAlert)
	if err != nil {
		t.Errorf("Erreur lors de l'activation avec lowPriorityAlert: %v", err)
	}

	// Les nœuds Alpha propagent via ActivateLeft avec des tokens
	// Vérifier que seul alertFact (type=Alert ET priority>5) a activé le nœud
	if activateLeftCount != 1 {
		t.Errorf("Attendait 1 activation ActivateLeft, reçu %d", activateLeftCount)
	}

	// Vérifier aussi que la condition a bien fonctionné
	if alphaNode.Memory.Facts == nil || len(alphaNode.Memory.Facts) != 1 {
		t.Errorf("Attendait 1 fait en mémoire, reçu %d", len(alphaNode.Memory.Facts))
	}
}

// MockStorage est un storage de test
type MockStorage struct{}

func (m *MockStorage) SaveMemory(nodeID string, memory *WorkingMemory) error {
	return nil
}

func (m *MockStorage) LoadMemory(nodeID string) (*WorkingMemory, error) {
	return &WorkingMemory{NodeID: nodeID}, nil
}

func (m *MockStorage) DeleteMemory(nodeID string) error {
	return nil
}

func (m *MockStorage) ListNodes() ([]string, error) {
	return []string{}, nil
}

// MockBetaNode est un nœud bêta de test pour l'intégration
type MockBetaNode struct {
	id               string
	activateFunc     func(*Fact, *Token)
	activateLeftFunc func(*Token)
}

func (m *MockBetaNode) GetID() string {
	return m.id
}

func (m *MockBetaNode) GetType() string {
	return "MockBeta"
}

func (m *MockBetaNode) GetMemory() *WorkingMemory {
	return &WorkingMemory{NodeID: m.id}
}

func (m *MockBetaNode) ActivateLeft(token *Token) error {
	if m.activateLeftFunc != nil {
		m.activateLeftFunc(token)
	}
	return nil
}

func (m *MockBetaNode) ActivateRight(fact *Fact) error {
	if m.activateFunc != nil {
		m.activateFunc(fact, nil)
	}
	return nil
}

func (m *MockBetaNode) AddChild(child Node) {
	// Not implemented for mock
}

func (m *MockBetaNode) GetChildren() []Node {
	return []Node{}
}

// Benchmark pour mesurer les performances
func BenchmarkAlphaConditionEvaluator(b *testing.B) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	fact := &Fact{
		ID:   "benchmark_fact",
		Type: "BenchmarkEvent",
		Fields: map[string]interface{}{
			"id":       42,
			"score":    95.5,
			"active":   true,
			"category": "test",
		},
	}

	// Condition complexe pour le benchmark
	condition := builder.And(
		builder.Or(
			builder.FieldEquals("event", "active", true),
			builder.FieldGreaterThan("event", "score", 90.0),
		),
		builder.FieldRange("event", "id", 0, 100),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evaluator.ClearBindings()
		_, err := evaluator.EvaluateCondition(condition, fact, "event")
		if err != nil {
			b.Fatalf("Erreur dans le benchmark: %v", err)
		}
	}
}
