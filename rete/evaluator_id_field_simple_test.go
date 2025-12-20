// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestEvaluator_IDFieldAccess_BasicComparisons teste les comparaisons de faits
// en utilisant des liaisons fact-to-fact (object binding) au lieu de l'ancien champ 'id'.
//
// Nouvelle sémantique:
// - Le champ 'id' n'est plus accessible dans les expressions
// - Les liaisons entre faits utilisent des affectations directes (ex: user: u)
// - Les comparaisons utilisent l'égalité d'objets (ex: m.user == u)
func TestEvaluator_IDFieldAccess_BasicComparisons(t *testing.T) {
	t.Log("🧪 TEST: Evaluator - Comparaisons avec liaisons fact-to-fact")
	t.Log("=============================================================")

	tests := []struct {
		name       string
		userFact   *Fact
		mailFact   *Fact
		expression map[string]interface{}
		wantResult bool
		wantErr    bool
	}{
		{
			name: "Égalité fact-to-fact: m.user == u (liaison directe)",
			userFact: &Fact{
				ID:   "User~Alice",
				Type: "User",
				Fields: map[string]interface{}{
					"nom":    "Alice",
					"prenom": "Dupont",
				},
			},
			mailFact: &Fact{
				ID:   "Mail~1",
				Type: "Mail",
				Fields: map[string]interface{}{
					"adresse": "alice@example.com",
					"user": &Fact{
						ID:   "User~Alice",
						Type: "User",
						Fields: map[string]interface{}{
							"nom":    "Alice",
							"prenom": "Dupont",
						},
					},
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "==",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "m",
					"field":  "user",
				},
				"right": map[string]interface{}{
					"type": "variable",
					"name": "u",
				},
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Inégalité fact-to-fact: m.user != autre_user",
			userFact: &Fact{
				ID:   "User~Bob",
				Type: "User",
				Fields: map[string]interface{}{
					"nom":    "Bob",
					"prenom": "Martin",
				},
			},
			mailFact: &Fact{
				ID:   "Mail~1",
				Type: "Mail",
				Fields: map[string]interface{}{
					"adresse": "alice@example.com",
					"user": &Fact{
						ID:   "User~Alice",
						Type: "User",
						Fields: map[string]interface{}{
							"nom":    "Alice",
							"prenom": "Dupont",
						},
					},
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "!=",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "m",
					"field":  "user",
				},
				"right": map[string]interface{}{
					"type": "variable",
					"name": "u",
				},
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Comparaison champ primaire: u.nom == 'Alice'",
			userFact: &Fact{
				ID:   "User~Alice",
				Type: "User",
				Fields: map[string]interface{}{
					"nom":    "Alice",
					"prenom": "Dupont",
				},
			},
			mailFact: nil, // Pas utilisé dans ce test
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "==",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "u",
					"field":  "nom",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Alice",
				},
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "CONTAINS sur champ texte: u.nom CONTAINS 'Ali'",
			userFact: &Fact{
				ID:   "User~Alice",
				Type: "User",
				Fields: map[string]interface{}{
					"nom":    "Alice",
					"prenom": "Dupont",
				},
			},
			mailFact: nil,
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "CONTAINS",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "u",
					"field":  "nom",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Ali",
				},
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Comparaison champ composite: u.prenom == 'Dupont'",
			userFact: &Fact{
				ID:   "User~Alice",
				Type: "User",
				Fields: map[string]interface{}{
					"nom":    "Alice",
					"prenom": "Dupont",
				},
			},
			mailFact: nil,
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "==",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "u",
					"field":  "prenom",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Dupont",
				},
			},
			wantResult: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := NewAlphaConditionEvaluator()

			// Lier la variable 'u' au fait utilisateur
			if tt.userFact != nil {
				eval.variableBindings["u"] = tt.userFact
			}

			// Lier la variable 'm' au fait mail (si présent)
			if tt.mailFact != nil {
				eval.variableBindings["m"] = tt.mailFact
			}

			// Évaluer l'expression
			result, err := eval.evaluateExpression(tt.expression)

			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ evaluateExpression() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				if result != tt.wantResult {
					t.Errorf("❌ evaluateExpression() = %v, attendu %v", result, tt.wantResult)
				} else {
					t.Logf("✅ Test réussi: résultat = %v", result)
				}
			}
		})
	}

	t.Log("✅ Test complet: Comparaisons avec liaisons fact-to-fact")
}
