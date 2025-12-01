// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

// TestMultipleActionsBasic tests basic multiple actions parsing
func TestMultipleActionsBasic(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - CAS DE BASE")
	t.Log("========================================")

	input := `type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id), log("Adult detected")
`

	t.Log("📝 Parsing du fichier avec actions multiples...")
	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Expressions) != 1 {
		t.Fatalf("❌ Attendu 1 expression, reçu %d", len(program.Expressions))
	}

	expr := program.Expressions[0]
	if expr.Action == nil {
		t.Fatal("❌ Action est nil")
	}

	jobs := expr.Action.GetJobs()
	if len(jobs) != 2 {
		t.Fatalf("❌ Attendu 2 jobs, reçu %d", len(jobs))
	}

	// Vérifier le premier job
	if jobs[0].Name != "adult" {
		t.Errorf("❌ Premier job: attendu 'adult', reçu '%s'", jobs[0].Name)
	}

	// Vérifier le second job
	if jobs[1].Name != "log" {
		t.Errorf("❌ Second job: attendu 'log', reçu '%s'", jobs[1].Name)
	}

	t.Log("✅ Test réussi - Actions multiples parsées correctement")
}

// TestMultipleActionsThreeJobs tests three actions in a rule
func TestMultipleActionsThreeJobs(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - TROIS ACTIONS")
	t.Log("==========================================")

	input := `type Person : <id: string, name: string, age: number>
type Order : <id: string, customer_id: string, amount: number>

rule r1 : {p: Person, o: Order} / p.id == o.customer_id ==> customer_order(p.id, o.id), update_stats(o.amount), notify("admin")
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Expressions) != 1 {
		t.Fatalf("❌ Attendu 1 expression, reçu %d", len(program.Expressions))
	}

	expr := program.Expressions[0]
	jobs := expr.Action.GetJobs()

	if len(jobs) != 3 {
		t.Fatalf("❌ Attendu 3 jobs, reçu %d", len(jobs))
	}

	expectedNames := []string{"customer_order", "update_stats", "notify"}
	for i, expectedName := range expectedNames {
		if jobs[i].Name != expectedName {
			t.Errorf("❌ Job %d: attendu '%s', reçu '%s'", i, expectedName, jobs[i].Name)
		}
	}

	t.Log("✅ Test réussi - Trois actions parsées correctement")
}

// TestMultipleActionsWithAggregation tests multiple actions with aggregation
func TestMultipleActionsWithAggregation(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - AVEC AGRÉGATION")
	t.Log("============================================")

	input := `type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_stats : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Stats"), update_dashboard(d.id, avg_sal), alert_hr(d.id)
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Expressions) != 1 {
		t.Fatalf("❌ Attendu 1 expression, reçu %d", len(program.Expressions))
	}

	expr := program.Expressions[0]
	jobs := expr.Action.GetJobs()

	if len(jobs) != 3 {
		t.Fatalf("❌ Attendu 3 jobs, reçu %d", len(jobs))
	}

	expectedNames := []string{"print", "update_dashboard", "alert_hr"}
	for i, expectedName := range expectedNames {
		if jobs[i].Name != expectedName {
			t.Errorf("❌ Job %d: attendu '%s', reçu '%s'", i, expectedName, jobs[i].Name)
		}
	}

	t.Log("✅ Test réussi - Actions multiples avec agrégation")
}

// TestBackwardCompatibilitySingleAction tests backward compatibility with single action
func TestBackwardCompatibilitySingleAction(t *testing.T) {
	t.Log("🧪 TEST RÉTROCOMPATIBILITÉ - ACTION UNIQUE")
	t.Log("==========================================")

	input := `type Person : <id: string, age: number>

rule adults : {p: Person} / p.age >= 18 ==> print("Adult")
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(program.Expressions) != 1 {
		t.Fatalf("❌ Attendu 1 expression, reçu %d", len(program.Expressions))
	}

	expr := program.Expressions[0]
	jobs := expr.Action.GetJobs()

	if len(jobs) != 1 {
		t.Fatalf("❌ Attendu 1 job, reçu %d", len(jobs))
	}

	if jobs[0].Name != "print" {
		t.Errorf("❌ Job: attendu 'print', reçu '%s'", jobs[0].Name)
	}

	t.Log("✅ Test réussi - Rétrocompatibilité maintenue")
}

// TestMultipleActionsWithComplexArguments tests multiple actions with complex arguments
func TestMultipleActionsWithComplexArguments(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - ARGUMENTS COMPLEXES")
	t.Log("================================================")

	input := `type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id), log(p.name, p.age), notify("admin", p.id)
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	expr := program.Expressions[0]
	jobs := expr.Action.GetJobs()

	if len(jobs) != 3 {
		t.Fatalf("❌ Attendu 3 jobs, reçu %d", len(jobs))
	}

	// Vérifier que chaque job a le bon nombre d'arguments
	expectedArgCounts := []int{1, 2, 2}
	for i, expectedCount := range expectedArgCounts {
		if len(jobs[i].Args) != expectedCount {
			t.Errorf("❌ Job %d: attendu %d arguments, reçu %d", i, expectedCount, len(jobs[i].Args))
		}
	}

	t.Log("✅ Test réussi - Arguments complexes gérés correctement")
}

// TestMultipleActionsWithArithmeticExpressions tests actions with arithmetic in arguments
func TestMultipleActionsWithArithmeticExpressions(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - EXPRESSIONS ARITHMÉTIQUES")
	t.Log("======================================================")

	input := `type Person : <id: string, age: number, salary: number>

rule r1 : {p: Person} / p.age > 18 ==> update_bonus(p.id, p.salary * 1.1), log("Bonus calculated"), notify_payroll(p.id)
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	expr := program.Expressions[0]
	jobs := expr.Action.GetJobs()

	if len(jobs) != 3 {
		t.Fatalf("❌ Attendu 3 jobs, reçu %d", len(jobs))
	}

	expectedNames := []string{"update_bonus", "log", "notify_payroll"}
	for i, expectedName := range expectedNames {
		if jobs[i].Name != expectedName {
			t.Errorf("❌ Job %d: attendu '%s', reçu '%s'", i, expectedName, jobs[i].Name)
		}
	}

	t.Log("✅ Test réussi - Expressions arithmétiques dans les arguments")
}

// TestMultipleActionsParseError tests that parse errors are properly caught
func TestMultipleActionsParseError(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - ERREURS DE PARSING")
	t.Log("===============================================")

	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "Virgule manquante entre actions",
			input: `type Person : <id: string>
rule r1 : {p: Person} / p.id == "1" ==> action1(p.id) action2(p.id)`,
		},
		{
			name: "Parenthèse fermante manquante",
			input: `type Person : <id: string>
rule r1 : {p: Person} / p.id == "1" ==> action1(p.id, action2(p.id)`,
		},
		{
			name: "Virgule en trop à la fin",
			input: `type Person : <id: string>
rule r1 : {p: Person} / p.id == "1" ==> action1(p.id), action2(p.id),`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tc.input))
			if err == nil {
				t.Errorf("❌ Attendu une erreur pour '%s', mais parsing réussi", tc.name)
			} else {
				t.Logf("✅ Erreur correctement détectée: %v", err)
			}
		})
	}

	t.Log("✅ Test réussi - Erreurs de parsing détectées")
}

// TestMultipleActionsStructure tests the structure of multiple actions
func TestMultipleActionsStructure(t *testing.T) {
	t.Log("🧪 TEST ACTIONS MULTIPLES - STRUCTURE")
	t.Log("======================================")

	input := `type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id), log("Adult")
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	expr := program.Expressions[0]

	// Vérifier que Jobs est peuplé (nouveau format)
	if len(expr.Action.Jobs) != 2 {
		t.Errorf("❌ Attendu 2 jobs dans le champ Jobs, reçu %d", len(expr.Action.Jobs))
	}

	// Vérifier que Job (ancien format) est nil
	if expr.Action.Job != nil {
		t.Error("❌ Le champ Job devrait être nil pour le nouveau format")
	}

	t.Log("✅ Test réussi - Structure correcte avec nouveau format")
}

// TestGetJobsMethodNilAction tests GetJobs with nil action
func TestGetJobsMethodNilAction(t *testing.T) {
	t.Log("🧪 TEST MÉTHODE GetJobs - ACTION NIL")
	t.Log("=====================================")

	action := &Action{Type: "action"}
	jobs := action.GetJobs()

	if len(jobs) != 0 {
		t.Errorf("❌ Attendu 0 jobs pour action vide, reçu %d", len(jobs))
	}

	t.Log("✅ Test réussi - GetJobs gère correctement l'absence de jobs")
}

// TestGetJobsMethodOldFormat tests GetJobs with old single job format
func TestGetJobsMethodOldFormat(t *testing.T) {
	t.Log("🧪 TEST MÉTHODE GetJobs - ANCIEN FORMAT")
	t.Log("========================================")

	job := &JobCall{
		Type: "jobCall",
		Name: "testJob",
		Args: []interface{}{"arg1"},
	}

	action := &Action{
		Type: "action",
		Job:  job,
	}

	jobs := action.GetJobs()

	if len(jobs) != 1 {
		t.Fatalf("❌ Attendu 1 job, reçu %d", len(jobs))
	}

	if jobs[0].Name != "testJob" {
		t.Errorf("❌ Attendu 'testJob', reçu '%s'", jobs[0].Name)
	}

	t.Log("✅ Test réussi - GetJobs gère l'ancien format correctement")
}

// TestGetJobsMethodNewFormat tests GetJobs with new multiple jobs format
func TestGetJobsMethodNewFormat(t *testing.T) {
	t.Log("🧪 TEST MÉTHODE GetJobs - NOUVEAU FORMAT")
	t.Log("=========================================")

	jobs := []JobCall{
		{Type: "jobCall", Name: "job1", Args: []interface{}{}},
		{Type: "jobCall", Name: "job2", Args: []interface{}{}},
		{Type: "jobCall", Name: "job3", Args: []interface{}{}},
	}

	action := &Action{
		Type: "action",
		Jobs: jobs,
	}

	retrievedJobs := action.GetJobs()

	if len(retrievedJobs) != 3 {
		t.Fatalf("❌ Attendu 3 jobs, reçu %d", len(retrievedJobs))
	}

	for i, job := range retrievedJobs {
		expectedName := jobs[i].Name
		if job.Name != expectedName {
			t.Errorf("❌ Job %d: attendu '%s', reçu '%s'", i, expectedName, job.Name)
		}
	}

	t.Log("✅ Test réussi - GetJobs gère le nouveau format correctement")
}
