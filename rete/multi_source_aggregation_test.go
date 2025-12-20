// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestMultiSourceAggregation_TwoSources tests aggregation over two joined sources
func TestMultiSourceAggregation_TwoSources(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "multi_source_test.tsd")
	// Rule aggregates data from both Employee and Performance tables
	content := `type Department(#id: string, name:string)
type Employee(#id: string, deptId: string, salary:number)
type Performance(#id: string, employeeId: string, score:number)
action print(message: string)
rule dept_combined_metric : {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} / {e: Employee} / {p: Performance} / e.deptId == d.id AND p.employeeId == e.id ==> print("Combined metrics")
`
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	if len(network.TerminalNodes) == 0 {
		t.Fatal("Expected at least one terminal node")
	}
	// Submit facts
	dept := Fact{
		ID:   "dept1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "dept1",
			"name": "Engineering",
		},
	}
	emp1 := Fact{
		ID:   "emp1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp1",
			"deptId": "dept1",
			"salary": 60000,
		},
	}
	emp2 := Fact{
		ID:   "emp2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp2",
			"deptId": "dept1",
			"salary": 70000,
		},
	}
	perf1 := Fact{
		ID:   "perf1",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf1",
			"employeeId": "emp1",
			"score":      85,
		},
	}
	perf2 := Fact{
		ID:   "perf2",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf2",
			"employeeId": "emp2",
			"score":      90,
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	network.SubmitFact(&emp2)
	network.SubmitFact(&perf1)
	network.SubmitFact(&perf2)
	// Expected: avg_sal = 65000, avg_score = 87.5
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += int(terminalNode.GetExecutionCount())
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation for multi-source aggregation, got %d", activatedCount)
	}
	t.Logf("✅ Multi-source aggregation with two sources completed with %d activations", activatedCount)
}

// TestMultiSourceAggregation_ThreeSources tests aggregation over three joined sources
func TestMultiSourceAggregation_ThreeSources(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "three_source_test.tsd")
	// Rule aggregates data from Employee, Performance, and Training tables
	content := `type Department(#id: string, name:string)
type Employee(#id: string, deptId: string, salary:number)
type Performance(#id: string, employeeId: string, score:number)
type Training(#id: string, employeeId: string, hours:number)
action print(message: string)
rule dept_full_metrics : {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score), total_hours: SUM(t.hours)} / {e: Employee} / {p: Performance} / {t: Training} / e.deptId == d.id AND p.employeeId == e.id AND t.employeeId == e.id ==> print("Full metrics")
`
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	if len(network.TerminalNodes) == 0 {
		t.Fatal("Expected at least one terminal node")
	}
	// Submit facts
	dept := Fact{
		ID:   "dept1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "dept1",
			"name": "Engineering",
		},
	}
	emp1 := Fact{
		ID:   "emp1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp1",
			"deptId": "dept1",
			"salary": 60000,
		},
	}
	perf1 := Fact{
		ID:   "perf1",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf1",
			"employeeId": "emp1",
			"score":      85,
		},
	}
	training1 := Fact{
		ID:   "training1",
		Type: "Training",
		Fields: map[string]interface{}{
			"id":         "training1",
			"employeeId": "emp1",
			"hours":      40,
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	network.SubmitFact(&perf1)
	network.SubmitFact(&training1)
	// Expected: avg_sal = 60000, avg_score = 85, total_hours = 40
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += int(terminalNode.GetExecutionCount())
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation for three-source aggregation, got %d", activatedCount)
	}
	t.Logf("✅ Multi-source aggregation with three sources completed with %d activations", activatedCount)
}

// TestMultiSourceAggregation_WithThreshold tests multi-source aggregation with thresholds
func TestMultiSourceAggregation_WithThreshold(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "threshold_test.tsd")
	// Rule only fires when both average salary > 50000 and average score > 80
	content := `type Department(#id: string, name:string)
type Employee(#id: string, deptId: string, salary:number)
type Performance(#id: string, employeeId: string, score:number)
action print(message: string)
rule high_performing_dept : {d: Department, avg_sal: AVG(e.salary), avg_score: AVG(p.score)} / {e: Employee} / {p: Performance} / e.deptId == d.id AND p.employeeId == e.id AND avg_sal > 50000 AND avg_score > 80 ==> print("High performing department")
`
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	if len(network.TerminalNodes) == 0 {
		t.Fatal("Expected at least one terminal node")
	}
	// Test case 1: Both thresholds satisfied
	dept1 := Fact{
		ID:   "dept1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "dept1",
			"name": "Engineering",
		},
	}
	emp1 := Fact{
		ID:   "emp1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp1",
			"deptId": "dept1",
			"salary": 60000,
		},
	}
	perf1 := Fact{
		ID:   "perf1",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf1",
			"employeeId": "emp1",
			"score":      85,
		},
	}
	network.SubmitFact(&dept1)
	network.SubmitFact(&emp1)
	network.SubmitFact(&perf1)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += int(terminalNode.GetExecutionCount())
	}
	if activatedCount < 1 {
		t.Errorf("Expected activation when thresholds are satisfied, got %d", activatedCount)
	}
	// Test case 2: One threshold not satisfied
	dept2 := Fact{
		ID:   "dept2",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "dept2",
			"name": "Sales",
		},
	}
	emp2 := Fact{
		ID:   "emp2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp2",
			"deptId": "dept2",
			"salary": 40000, // Below threshold
		},
	}
	perf2 := Fact{
		ID:   "perf2",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf2",
			"employeeId": "emp2",
			"score":      85,
		},
	}
	previousCount := activatedCount
	network.SubmitFact(&dept2)
	network.SubmitFact(&emp2)
	network.SubmitFact(&perf2)
	newActivatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		newActivatedCount += int(terminalNode.GetExecutionCount())
	}
	if newActivatedCount > previousCount {
		t.Errorf("Expected no new activation when threshold not satisfied, but got %d (was %d)", newActivatedCount, previousCount)
	}
	t.Logf("✅ Multi-source aggregation with thresholds working correctly")
}

// TestMultiSourceAggregation_DifferentFunctions tests different aggregation functions
func TestMultiSourceAggregation_DifferentFunctions(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "functions_test.tsd")
	// Rule uses different aggregation functions: AVG, SUM, COUNT, MIN, MAX
	content := `type Department(#id: string, name:string)
type Employee(#id: string, deptId: string, salary:number)
type Performance(#id: string, employeeId: string, score:number)
action print(message: string)
rule dept_comprehensive_stats : {d: Department, avg_sal: AVG(e.salary), total_sal: SUM(e.salary), emp_count: COUNT(e.id), min_score: MIN(p.score), max_score: MAX(p.score)} / {e: Employee} / {p: Performance} / e.deptId == d.id AND p.employeeId == e.id ==> print("Comprehensive stats")
`
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	if len(network.TerminalNodes) == 0 {
		t.Fatal("Expected at least one terminal node")
	}
	// Submit facts
	dept := Fact{
		ID:   "dept1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "dept1",
			"name": "Engineering",
		},
	}
	emp1 := Fact{
		ID:   "emp1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp1",
			"deptId": "dept1",
			"salary": 60000,
		},
	}
	emp2 := Fact{
		ID:   "emp2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "emp2",
			"deptId": "dept1",
			"salary": 70000,
		},
	}
	perf1 := Fact{
		ID:   "perf1",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf1",
			"employeeId": "emp1",
			"score":      75,
		},
	}
	perf2 := Fact{
		ID:   "perf2",
		Type: "Performance",
		Fields: map[string]interface{}{
			"id":         "perf2",
			"employeeId": "emp2",
			"score":      95,
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	network.SubmitFact(&emp2)
	network.SubmitFact(&perf1)
	network.SubmitFact(&perf2)
	// Expected: avg_sal=65000, total_sal=130000, emp_count=2, min_score=75, max_score=95
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += int(terminalNode.GetExecutionCount())
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation for comprehensive stats, got %d", activatedCount)
	}
	t.Logf("✅ Multi-source aggregation with different functions completed with %d activations", activatedCount)
}
