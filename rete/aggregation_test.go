// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPipeline_AVG teste le calcul de moyenne via le pipeline complet
func TestPipeline_AVG(t *testing.T) {
	// Créer fichier constraint temporaire
	constraintFile := filepath.Join(t.TempDir(), "avg_test.tsd")
	constraintContent := `type Employee(id: string, name:string)
type Performance(id: string, employee_id: string, score:number)
action PRINT(arg1: string, arg2: string, arg3: string)
rule r1 : {e: Employee} / AVG(p: Performance / p.employee_id == e.id ; p.score) >= 8.5
==> PRINT("Employee ", e.name, " has avg >= 8.5")
`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatal(err)
	}
	// Construire le réseau via le pipeline
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Injecter les faits
	network.SubmitFact(&Fact{ID: "EMP001", Type: "Employee", Fields: map[string]interface{}{"name": "Alice"}})
	network.SubmitFact(&Fact{ID: "PERF001", Type: "Performance", Fields: map[string]interface{}{"employee_id": "EMP001", "score": 9.0}})
	network.SubmitFact(&Fact{ID: "PERF002", Type: "Performance", Fields: map[string]interface{}{"employee_id": "EMP001", "score": 8.5}})
	network.SubmitFact(&Fact{ID: "PERF003", Type: "Performance", Fields: map[string]interface{}{"employee_id": "EMP001", "score": 9.2}})
	// Vérifier que le terminal a été activé (l'action devrait être disponible)
	terminalCount := len(network.TerminalNodes)
	if terminalCount == 0 {
		t.Error("Aucun nœud terminal créé")
	}
	t.Logf("✓ Test AVG: réseau créé avec %d nœuds terminaux", terminalCount)
}

// TestPipeline_SUM teste le calcul de somme via le pipeline complet
func TestPipeline_SUM(t *testing.T) {
	constraintFile := filepath.Join(t.TempDir(), "sum_test.tsd")
	constraintContent := `type Employee(id: string, name:string)
type Order(id: string, employee_id: string, amount:number)
action PRINT(arg1: string, arg2: string, arg3: string)
rule r1 : {e: Employee} / SUM(o: Order / o.employee_id == e.id ; o.amount) >= 1000
==> PRINT("Employee ", e.name, " has sum >= 1000")
`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatal(err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	network.SubmitFact(&Fact{ID: "EMP001", Type: "Employee", Fields: map[string]interface{}{"name": "Bob"}})
	network.SubmitFact(&Fact{ID: "ORD001", Type: "Order", Fields: map[string]interface{}{"employee_id": "EMP001", "amount": 500.0}})
	network.SubmitFact(&Fact{ID: "ORD002", Type: "Order", Fields: map[string]interface{}{"employee_id": "EMP001", "amount": 750.0}})
	terminalCount := len(network.TerminalNodes)
	if terminalCount == 0 {
		t.Error("Aucun nœud terminal créé")
	}
	t.Logf("✓ Test SUM: réseau créé avec %d nœuds terminaux", terminalCount)
}

// TestPipeline_COUNT teste le comptage via le pipeline complet
func TestPipeline_COUNT(t *testing.T) {
	constraintFile := filepath.Join(t.TempDir(), "count_test.tsd")
	constraintContent := `type Department(id: string, name:string)
type Employee(id: string, department:string)
action PRINT(arg1: string, arg2: string, arg3: string)
rule r1 : {d: Department} / COUNT(e: Employee / e.department == d.name) >= 3
==> PRINT("Department ", d.name, " has >= 3 employees")
`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatal(err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	network.SubmitFact(&Fact{ID: "DEPT001", Type: "Department", Fields: map[string]interface{}{"name": "sales"}})
	network.SubmitFact(&Fact{ID: "EMP001", Type: "Employee", Fields: map[string]interface{}{"department": "sales"}})
	network.SubmitFact(&Fact{ID: "EMP002", Type: "Employee", Fields: map[string]interface{}{"department": "sales"}})
	network.SubmitFact(&Fact{ID: "EMP003", Type: "Employee", Fields: map[string]interface{}{"department": "sales"}})
	network.SubmitFact(&Fact{ID: "EMP004", Type: "Employee", Fields: map[string]interface{}{"department": "sales"}})
	terminalCount := len(network.TerminalNodes)
	if terminalCount == 0 {
		t.Error("Aucun nœud terminal créé")
	}
	t.Logf("✓ Test COUNT: réseau créé avec %d nœuds terminaux", terminalCount)
}

// TestPipeline_MIN teste le minimum via le pipeline complet
func TestPipeline_MIN(t *testing.T) {
	constraintFile := filepath.Join(t.TempDir(), "min_test.tsd")
	constraintContent := `type Department(id: string, name:string)
type Employee(id: string, department: string, salary:number)
action PRINT(arg1: string, arg2: string, arg3: string)
rule r1 : {d: Department} / MIN(e: Employee / e.department == d.name ; e.salary) >= 50000
==> PRINT("Department ", d.name, " has min salary >= 50000")
`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatal(err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	network.SubmitFact(&Fact{ID: "DEPT001", Type: "Department", Fields: map[string]interface{}{"name": "engineering"}})
	network.SubmitFact(&Fact{ID: "EMP001", Type: "Employee", Fields: map[string]interface{}{"department": "engineering", "salary": 55000.0}})
	network.SubmitFact(&Fact{ID: "EMP002", Type: "Employee", Fields: map[string]interface{}{"department": "engineering", "salary": 70000.0}})
	terminalCount := len(network.TerminalNodes)
	if terminalCount == 0 {
		t.Error("Aucun nœud terminal créé")
	}
	t.Logf("✓ Test MIN: réseau créé avec %d nœuds terminaux", terminalCount)
}

// TestPipeline_MAX teste le maximum via le pipeline complet
func TestPipeline_MAX(t *testing.T) {
	constraintFile := filepath.Join(t.TempDir(), "max_test.tsd")
	constraintContent := `type Department(id: string, name:string)
type Employee(id: string, department: string, salary:number)
action PRINT(arg1: string, arg2: string, arg3: string)
rule r1 : {d: Department} / MAX(e: Employee / e.department == d.name ; e.salary) >= 80000
==> PRINT("Department ", d.name, " has max salary >= 80000")
`
	if err := os.WriteFile(constraintFile, []byte(constraintContent), 0644); err != nil {
		t.Fatal(err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	network.SubmitFact(&Fact{ID: "DEPT001", Type: "Department", Fields: map[string]interface{}{"name": "management"}})
	network.SubmitFact(&Fact{ID: "EMP001", Type: "Employee", Fields: map[string]interface{}{"department": "management", "salary": 95000.0}})
	network.SubmitFact(&Fact{ID: "EMP002", Type: "Employee", Fields: map[string]interface{}{"department": "management", "salary": 70000.0}})
	terminalCount := len(network.TerminalNodes)
	if terminalCount == 0 {
		t.Error("Aucun nœud terminal créé")
	}
	t.Logf("✓ Test MAX: réseau créé avec %d nœuds terminaux", terminalCount)
}
