package rete

import (
	"testing"
)

func TestAlphaConditionEvaluator(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	// Créer un fait de test
	fact := &Fact{
		ID:   "test_event_1",
		Type: "Event",
		Fields: map[string]interface{}{
			"id":       1,
			"priority": 5,
			"active":   true,
			"name":     "test_event",
			"score":    95.5,
		},
	}

	tests := []struct {
		name        string
		condition   interface{}
		expected    bool
		expectError bool
	}{
		{
			name:      "Condition toujours vraie",
			condition: builder.True(),
			expected:  true,
		},
		{
			name:      "Condition toujours fausse",
			condition: builder.False(),
			expected:  false,
		},
		{
			name:      "Égalité entier",
			condition: builder.FieldEquals("event", "id", 1),
			expected:  true,
		},
		{
			name:      "Égalité entier (faux)",
			condition: builder.FieldEquals("event", "id", 2),
			expected:  false,
		},
		{
			name:      "Égalité booléen",
			condition: builder.FieldEquals("event", "active", true),
			expected:  true,
		},
		{
			name:      "Égalité string",
			condition: builder.FieldEquals("event", "name", "test_event"),
			expected:  true,
		},
		{
			name:      "Inégalité",
			condition: builder.FieldNotEquals("event", "id", 2),
			expected:  true,
		},
		{
			name:      "Comparaison inférieure",
			condition: builder.FieldLessThan("event", "priority", 10),
			expected:  true,
		},
		{
			name:      "Comparaison inférieure (faux)",
			condition: builder.FieldLessThan("event", "priority", 3),
			expected:  false,
		},
		{
			name:      "Comparaison supérieure",
			condition: builder.FieldGreaterThan("event", "priority", 3),
			expected:  true,
		},
		{
			name:      "Comparaison supérieure ou égale",
			condition: builder.FieldGreaterOrEqual("event", "priority", 5),
			expected:  true,
		},
		{
			name:      "Comparaison inférieure ou égale",
			condition: builder.FieldLessOrEqual("event", "priority", 5),
			expected:  true,
		},
		{
			name: "Condition AND (vraie)",
			condition: builder.And(
				builder.FieldEquals("event", "active", true),
				builder.FieldGreaterThan("event", "priority", 3),
			),
			expected: true,
		},
		{
			name: "Condition AND (fausse)",
			condition: builder.And(
				builder.FieldEquals("event", "active", true),
				builder.FieldLessThan("event", "priority", 3),
			),
			expected: false,
		},
		{
			name: "Condition OR (vraie)",
			condition: builder.Or(
				builder.FieldEquals("event", "active", false),
				builder.FieldGreaterThan("event", "priority", 3),
			),
			expected: true,
		},
		{
			name: "Condition OR (fausse)",
			condition: builder.Or(
				builder.FieldEquals("event", "active", false),
				builder.FieldLessThan("event", "priority", 3),
			),
			expected: false,
		},
		{
			name:      "Plage de valeurs",
			condition: builder.FieldRange("event", "priority", 1, 10),
			expected:  true,
		},
		{
			name:      "Plage de valeurs (hors limites)",
			condition: builder.FieldRange("event", "priority", 10, 20),
			expected:  false,
		},
		{
			name:      "Valeur dans liste",
			condition: builder.FieldIn("event", "priority", 1, 3, 5, 7),
			expected:  true,
		},
		{
			name:      "Valeur pas dans liste",
			condition: builder.FieldNotIn("event", "priority", 1, 2, 3, 4),
			expected:  true,
		},
		{
			name: "Conditions multiples AND",
			condition: builder.AndMultiple(
				builder.FieldEquals("event", "active", true),
				builder.FieldGreaterThan("event", "priority", 3),
				builder.FieldLessThan("event", "priority", 10),
			),
			expected: true,
		},
		{
			name:      "Comparaison avec float",
			condition: builder.FieldGreaterThan("event", "score", 90.0),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.EvaluateCondition(tt.condition, fact, "event")

			if tt.expectError && err == nil {
				t.Errorf("Attendait une erreur mais n'en a pas eu")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Erreur inattendue: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Résultat attendu %v, obtenu %v", tt.expected, result)
			}
		})
	}
}

func TestAlphaConditionEvaluator_ErrorCases(t *testing.T) {
	evaluator := NewAlphaConditionEvaluator()
	builder := NewAlphaConditionBuilder()

	// Créer un fait de test
	fact := &Fact{
		ID:   "test_event_1",
		Type: "Event",
		Fields: map[string]interface{}{
			"id": 1,
		},
	}

	tests := []struct {
		name      string
		condition interface{}
	}{
		{
			name:      "Champ inexistant",
			condition: builder.FieldEquals("event", "nonexistent", 1),
		},
		{
			name:      "Variable non liée",
			condition: builder.FieldEquals("unknown", "id", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := evaluator.EvaluateCondition(tt.condition, fact, "event")

			if err == nil {
				t.Errorf("Attendait une erreur mais n'en a pas eu")
			}
		})
	}
}

func TestAlphaNode_WithConditions(t *testing.T) {
	storage := NewMemoryStorage()
	builder := NewAlphaConditionBuilder()

	// Créer un nœud alpha avec une condition
	condition := builder.And(
		builder.FieldEquals("event", "active", true),
		builder.FieldGreaterThan("event", "priority", 5),
	)

	alphaNode := NewAlphaNode("alpha_test", condition, "event", storage)

	tests := []struct {
		name       string
		fact       *Fact
		shouldPass bool
	}{
		{
			name: "Fait qui passe la condition",
			fact: &Fact{
				ID:   "event_1",
				Type: "Event",
				Fields: map[string]interface{}{
					"active":   true,
					"priority": 10,
				},
			},
			shouldPass: true,
		},
		{
			name: "Fait qui ne passe pas (active=false)",
			fact: &Fact{
				ID:   "event_2",
				Type: "Event",
				Fields: map[string]interface{}{
					"active":   false,
					"priority": 10,
				},
			},
			shouldPass: false,
		},
		{
			name: "Fait qui ne passe pas (priority trop faible)",
			fact: &Fact{
				ID:   "event_3",
				Type: "Event",
				Fields: map[string]interface{}{
					"active":   true,
					"priority": 3,
				},
			},
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compter les faits avant
			factsBefore := len(alphaNode.GetMemory().Facts)

			// Activer le nœud avec le fait
			err := alphaNode.ActivateRight(tt.fact)
			if err != nil {
				t.Errorf("Erreur inattendue: %v", err)
				return
			}

			// Vérifier si le fait a été ajouté à la mémoire
			factsAfter := len(alphaNode.GetMemory().Facts)

			if tt.shouldPass {
				if factsAfter != factsBefore+1 {
					t.Errorf("Le fait devrait avoir été ajouté à la mémoire")
				}
			} else {
				if factsAfter != factsBefore {
					t.Errorf("Le fait ne devrait pas avoir été ajouté à la mémoire")
				}
			}
		})
	}
}
