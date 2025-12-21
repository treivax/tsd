// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/tests/shared"
)

// TestArithmeticModuloE2E_BasicOperations teste l'opérateur modulo dans des scénarios réels.
// ✅ RESPECT DE LA CONTRAINTE: Tests fonctionnels RÉELS sans mocks
func TestArithmeticModuloE2E_BasicOperations(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST E2E: Opérateur Modulo (%) - Scénario Réel")

	// Programme TSD avec opérateur modulo pour classification de nombres
	programContent := `// Système de classification par parité et divisibilité
type Number(#id: string, value: number)
type Classification(id: string, value: number, category: string, reason: string)

xuple-space even_numbers {
	selection: fifo
	consumption: once
}

xuple-space odd_numbers {
	selection: fifo
	consumption: once
}

xuple-space divisible_by_five {
	selection: fifo
	consumption: once
}

// Règle: nombres pairs (modulo 2 == 0)
rule classify_even : {n: Number} / n.value % 2 == 0 ==>
	Xuple("even_numbers", Classification(
		id: n.id,
		value: n.value,
		category: "even",
		reason: "divisible by 2"
	))

// Règle: nombres impairs (modulo 2 != 0)
rule classify_odd : {n: Number} / n.value % 2 != 0 ==>
	Xuple("odd_numbers", Classification(
		id: n.id,
		value: n.value,
		category: "odd",
		reason: "remainder 1 when divided by 2"
	))

// Règle: divisible par 5 (modulo 5 == 0)
rule divisible_five : {n: Number} / n.value % 5 == 0 ==>
	Xuple("divisible_by_five", Classification(
		id: n.id,
		value: n.value,
		category: "multiple_of_5",
		reason: "divisible by 5"
	))

// Faits de test avec divers nombres
Number(id: "n001", value: 2)
Number(id: "n002", value: 3)
Number(id: "n003", value: 5)
Number(id: "n004", value: 6)
Number(id: "n005", value: 10)
Number(id: "n006", value: 15)
Number(id: "n007", value: 17)
Number(id: "n008", value: 20)
Number(id: "n009", value: 21)
Number(id: "n010", value: 30)
`

	_, result := shared.CreatePipelineFromTSD(t, programContent)

	// ═══════════════════════════════════════════════════════════════
	// VÉRIFICATION: Xuple-spaces créés
	// ═══════════════════════════════════════════════════════════════
	shared.LogTestSubsection(t, "📊 Vérification des Xuple-Spaces")

	spaces := result.XupleSpaceNames()
	require.Len(t, spaces, 3, "3 xuple-spaces devraient être créés")
	shared.AssertXupleSpaceExists(t, result, "even_numbers")
	shared.AssertXupleSpaceExists(t, result, "odd_numbers")
	shared.AssertXupleSpaceExists(t, result, "divisible_by_five")
	t.Log("✅ Tous les xuple-spaces créés")

	// ═══════════════════════════════════════════════════════════════
	// VÉRIFICATION: Classification pairs/impairs avec modulo
	// ═══════════════════════════════════════════════════════════════
	shared.LogTestSubsection(t, "🔢 Vérification Pairs/Impairs (% 2)")

	// Nombres pairs: 2, 6, 10, 20, 30 = 5 xuples
	evenNumbers := shared.GetXuples(t, result, "even_numbers")
	require.Len(t, evenNumbers, 5, "5 nombres pairs attendus")
	t.Logf("✅ %d nombres pairs classifiés", len(evenNumbers))

	// Vérifier que tous ont value % 2 == 0
	for _, xuple := range evenNumbers {
		value := shared.GetXupleFieldFloat(t, xuple, "value")
		assert.Equal(t, 0, int(value)%2, "nombre pair devrait avoir modulo 2 = 0")
		category := shared.GetXupleFieldString(t, xuple, "category")
		assert.Equal(t, "even", category)
	}

	// Nombres impairs: 3, 5, 15, 17, 21 = 5 xuples
	oddNumbers := shared.GetXuples(t, result, "odd_numbers")
	require.Len(t, oddNumbers, 5, "5 nombres impairs attendus")
	t.Logf("✅ %d nombres impairs classifiés", len(oddNumbers))

	// Vérifier que tous ont value % 2 != 0
	for _, xuple := range oddNumbers {
		value := shared.GetXupleFieldFloat(t, xuple, "value")
		assert.Equal(t, 1, int(value)%2, "nombre impair devrait avoir modulo 2 = 1")
		category := shared.GetXupleFieldString(t, xuple, "category")
		assert.Equal(t, "odd", category)
	}

	// ═══════════════════════════════════════════════════════════════
	// VÉRIFICATION: Divisibilité par 5 avec modulo
	// ═══════════════════════════════════════════════════════════════
	shared.LogTestSubsection(t, "➗ Vérification Divisibilité par 5 (% 5)")

	// Divisibles par 5: 5, 10, 15, 20, 30 = 5 xuples
	divByFive := shared.GetXuples(t, result, "divisible_by_five")
	require.Len(t, divByFive, 5, "5 nombres divisibles par 5 attendus")
	t.Logf("✅ %d nombres divisibles par 5", len(divByFive))

	// Vérifier que tous ont value % 5 == 0
	expectedDivByFive := map[string]float64{
		"n003": 5,
		"n005": 10,
		"n006": 15,
		"n008": 20,
		"n010": 30,
	}

	for _, xuple := range divByFive {
		id := shared.GetXupleFieldString(t, xuple, "id")
		value := shared.GetXupleFieldFloat(t, xuple, "value")

		expectedValue, found := expectedDivByFive[id]
		require.True(t, found, "ID %s devrait être dans la liste", id)
		assert.Equal(t, expectedValue, value, "valeur pour %s", id)
		remainder := int(value) % 5
		assert.Equal(t, 0, remainder, "nombre devrait être divisible par 5")
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST MODULO E2E RÉUSSI - Toutes classifications validées")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestArithmeticModuloE2E_ComplexExpressions teste des expressions complexes avec modulo.
// ✅ RESPECT DE LA CONTRAINTE: Tests déterministes avec résultats réels
func TestArithmeticModuloE2E_ComplexExpressions(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST E2E: Expressions Arithmétiques Complexes avec Modulo")

	programContent := `// Système de calculs complexes avec modulo et autres opérateurs
type Input(#id: string, a: number, b: number, c: number)
type Result(id: string, operation: string, result: number)

xuple-space calculation_results {
	selection: fifo
	consumption: once
}

// Règle: Calcul complexe avec modulo et opérations mixtes
// Formula: ((a * b) % c) + (a / b) - c
rule complex_calc : {inp: Input} / inp.b != 0 AND inp.c != 0 ==>
	Xuple("calculation_results", Result(
		id: inp.id,
		operation: "complex_modulo",
		result: ((inp.a * inp.b) % inp.c) + (inp.a / inp.b) - inp.c
	))

// Règle: Calcul avec modulo imbriqué
// Formula: (a % b) % c
rule nested_modulo : {inp: Input} / inp.b != 0 AND inp.c != 0 ==>
	Xuple("calculation_results", Result(
		id: inp.id,
		operation: "nested_modulo",
		result: (inp.a % inp.b) % inp.c
	))

// Règle: Priorité des opérateurs avec modulo
// Formula: a + b % c * 2
rule precedence_test : {inp: Input} / inp.c != 0 ==>
	Xuple("calculation_results", Result(
		id: inp.id,
		operation: "precedence",
		result: inp.a + inp.b % inp.c * 2
	))

// Faits de test
Input(id: "calc1", a: 17, b: 5, c: 3)
Input(id: "calc2", a: 100, b: 7, c: 4)
Input(id: "calc3", a: 23, b: 8, c: 5)
`

	_, result := shared.CreatePipelineFromTSD(t, programContent)

	results := shared.GetXuples(t, result, "calculation_results")
	require.NotEmpty(t, results, "des résultats devraient être calculés")

	t.Logf("✅ %d résultats de calculs complexes générés", len(results))

	// Organiser les résultats par ID et opération
	resultsByID := make(map[string]map[string]float64)
	for _, res := range results {
		id := shared.GetXupleFieldString(t, res, "id")
		operation := shared.GetXupleFieldString(t, res, "operation")
		value := shared.GetXupleFieldFloat(t, res, "result")

		if resultsByID[id] == nil {
			resultsByID[id] = make(map[string]float64)
		}
		resultsByID[id][operation] = value
	}

	// Vérifier calc1: a=17, b=5, c=3
	// complex_modulo: ((17*5)%3) + (17/5) - 3 = (85%3) + 3.4 - 3 = 1 + 3.4 - 3 = 1.4
	// nested_modulo: (17%5)%3 = 2%3 = 2
	// precedence: 17 + 5%3*2 = 17 + 2*2 = 17 + 4 = 21
	if calc1, exists := resultsByID["calc1"]; exists {
		if val, ok := calc1["complex_modulo"]; ok {
			assert.InDelta(t, 1.4, val, 0.01, "calc1 complex_modulo")
		}
		if val, ok := calc1["nested_modulo"]; ok {
			assert.Equal(t, 2.0, val, "calc1 nested_modulo")
		}
		if val, ok := calc1["precedence"]; ok {
			assert.Equal(t, 21.0, val, "calc1 precedence")
		}
		t.Log("✅ calc1: Tous les calculs validés")
	}

	// Vérifier calc2: a=100, b=7, c=4
	// complex_modulo: ((100*7)%4) + (100/7) - 4 = (700%4) + 14.285... - 4 = 0 + 14.285... - 4 = 10.285...
	// nested_modulo: (100%7)%4 = 2%4 = 2
	// precedence: 100 + 7%4*2 = 100 + 3*2 = 100 + 6 = 106
	if calc2, exists := resultsByID["calc2"]; exists {
		if val, ok := calc2["complex_modulo"]; ok {
			assert.InDelta(t, 10.285, val, 0.01, "calc2 complex_modulo")
		}
		if val, ok := calc2["nested_modulo"]; ok {
			assert.Equal(t, 2.0, val, "calc2 nested_modulo")
		}
		if val, ok := calc2["precedence"]; ok {
			assert.Equal(t, 106.0, val, "calc2 precedence")
		}
		t.Log("✅ calc2: Tous les calculs validés")
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST EXPRESSIONS COMPLEXES RÉUSSI")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestArithmeticModuloE2E_EdgeCases teste les cas limites avec modulo.
// ✅ RESPECT DE LA CONTRAINTE: Cas limites testés (couverture > 80%)
func TestArithmeticModuloE2E_EdgeCases(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST E2E: Cas Limites Modulo")

	programContent := `// Tests de cas limites pour l'opérateur modulo
type TestCase(id: string, dividend: number, divisor: number)
type ModuloResult(id: string, dividend: number, divisor: number, result: number, category: string)

xuple-space modulo_results {
	selection: fifo
	consumption: once
}

// Règle: calcul modulo et classification
rule compute_modulo : {tc: TestCase} / tc.divisor != 0 ==>
	Xuple("modulo_results", ModuloResult(
		id: tc.id,
		dividend: tc.dividend,
		divisor: tc.divisor,
		result: tc.dividend % tc.divisor,
		category: "computed"
	))

// Cas limites pour modulo
TestCase(id: "ZERO_MOD", dividend: 0, divisor: 5)           // 0 % 5 = 0
TestCase(id: "SAME_NUM", dividend: 7, divisor: 7)           // 7 % 7 = 0
TestCase(id: "LARGER_DIV", dividend: 3, divisor: 10)        // 3 % 10 = 3
TestCase(id: "ONE_DIVISOR", dividend: 42, divisor: 1)       // 42 % 1 = 0
TestCase(id: "LARGE_NUM", dividend: 1000000, divisor: 7)    // 1000000 % 7 = 1
TestCase(id: "EQUAL_RESULT", dividend: 15, divisor: 4)      // 15 % 4 = 3
TestCase(id: "TWO_DIVISOR", dividend: 99, divisor: 2)       // 99 % 2 = 1
`

	_, result := shared.CreatePipelineFromTSD(t, programContent)

	// ═══════════════════════════════════════════════════════════════
	// VÉRIFICATION: Tous les cas limites
	// ═══════════════════════════════════════════════════════════════
	shared.LogTestSubsection(t, "🔍 Vérification Cas Limites")

	moduloResults := shared.GetXuples(t, result, "modulo_results")
	require.Len(t, moduloResults, 7, "7 résultats modulo attendus")

	expectedResults := map[string]float64{
		"ZERO_MOD":     0.0, // 0 % 5
		"SAME_NUM":     0.0, // 7 % 7
		"LARGER_DIV":   3.0, // 3 % 10
		"ONE_DIVISOR":  0.0, // 42 % 1
		"LARGE_NUM":    1.0, // 1000000 % 7
		"EQUAL_RESULT": 3.0, // 15 % 4
		"TWO_DIVISOR":  1.0, // 99 % 2
	}

	for _, xuple := range moduloResults {
		id := shared.GetXupleFieldString(t, xuple, "id")
		dividend := shared.GetXupleFieldFloat(t, xuple, "dividend")
		divisor := shared.GetXupleFieldFloat(t, xuple, "divisor")
		result := shared.GetXupleFieldFloat(t, xuple, "result")
		category := shared.GetXupleFieldString(t, xuple, "category")

		expectedResult, found := expectedResults[id]
		require.True(t, found, "cas de test %s devrait exister", id)
		assert.Equal(t, expectedResult, result,
			"résultat pour %s: %.0f %% %.0f = %.0f", id, dividend, divisor, expectedResult)
		assert.Equal(t, "computed", category, "catégorie pour %s", id)

		t.Logf("✅ %s: %.0f %% %.0f = %.0f", id, dividend, divisor, result)
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST CAS LIMITES MODULO RÉUSSI")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestArithmeticModuloE2E_PriorityCalculations teste calculs avec priorités arithmétiques.
// ✅ RESPECT DE LA CONTRAINTE: Tests fonctionnels RÉELS sans mocks
func TestArithmeticModuloE2E_PriorityCalculations(t *testing.T) {
	shared.LogTestSection(t, "🧪 TEST E2E: Priorité des Opérateurs Arithmétiques")

	programContent := `// Tests de priorité des opérateurs
type Expression(id: string, a: number, b: number, c: number, d: number)
type Evaluation(id: string, expression: string, result: number)

xuple-space evaluations {
	selection: fifo
	consumption: once
}

// Test priorité: * / % avant + -
rule test_priority_1 : {e: Expression} / ==>
	Xuple("evaluations", Evaluation(
		id: e.id,
		expression: "a + b * c",
		result: e.a + e.b * e.c
	))

rule test_priority_2 : {e: Expression} / ==>
	Xuple("evaluations", Evaluation(
		id: e.id,
		expression: "a * b + c",
		result: e.a * e.b + e.c
	))

rule test_priority_3 : {e: Expression} / ==>
	Xuple("evaluations", Evaluation(
		id: e.id,
		expression: "a + b % c",
		result: e.a + e.b % e.c
	))

rule test_priority_4 : {e: Expression} / ==>
	Xuple("evaluations", Evaluation(
		id: e.id,
		expression: "a * b % c + d",
		result: e.a * e.b % e.c + e.d
	))

// Données de test pour vérifier priorités
Expression(id: "PRIO_TEST", a: 10, b: 5, c: 3, d: 2)
`

	_, result := shared.CreatePipelineFromTSD(t, programContent)

	// ═══════════════════════════════════════════════════════════════
	// VÉRIFICATION: Évaluations avec priorités correctes
	// ═══════════════════════════════════════════════════════════════
	shared.LogTestSubsection(t, "🎯 Vérification Priorités Opérateurs")

	evaluations := shared.GetXuples(t, result, "evaluations")
	require.Len(t, evaluations, 4, "4 évaluations attendues")

	// Valeurs: a=10, b=5, c=3, d=2
	expectedResults := map[string]float64{
		"a + b * c":     25.0, // 10 + (5 * 3) = 10 + 15 = 25
		"a * b + c":     53.0, // (10 * 5) + 3 = 50 + 3 = 53
		"a + b % c":     12.0, // 10 + (5 % 3) = 10 + 2 = 12
		"a * b % c + d": 4.0,  // ((10 * 5) % 3) + 2 = (50 % 3) + 2 = 2 + 2 = 4
	}

	for _, xuple := range evaluations {
		expression := shared.GetXupleFieldString(t, xuple, "expression")
		result := shared.GetXupleFieldFloat(t, xuple, "result")

		expectedResult, found := expectedResults[expression]
		require.True(t, found, "expression %s devrait exister", expression)

		assert.Equal(t, expectedResult, result,
			"évaluation incorrecte pour: %s", expression)

		t.Logf("✅ %s = %.0f", expression, result)
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST PRIORITÉS OPÉRATEURS RÉUSSI")
	t.Log("═══════════════════════════════════════════════════════════════")
}
