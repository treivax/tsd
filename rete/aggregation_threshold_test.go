// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
package rete
import (
	"os"
	"path/filepath"
	"testing"
)
// TestAggregationThreshold_GreaterThan tests aggregation with > threshold
func TestAggregationThreshold_GreaterThan(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "threshold_gt.tsd")
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
action print(message: string)
rule high_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id AND avg_sal > 50000 ==> print("High average salary")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	// Submit department
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "Engineering",
		},
	}
	// Test case 1: Average below threshold (should NOT fire)
	emp1 := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"deptId": "d1",
			"salary": 40000, // avg = 40000, should NOT fire (40000 > 50000 is false)
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount != 0 {
		t.Errorf("Expected 0 activations when avg (40000) <= threshold (50000), got %d", activatedCount)
	}
	// Test case 2: Average above threshold (should fire)
	emp2 := Fact{
		ID:   "E2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e2",
			"deptId": "d1",
			"salary": 70000, // avg = (40000 + 70000) / 2 = 55000, should fire
		},
	}
	network.SubmitFact(&emp2)
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation when avg (55000) > threshold (50000), got %d", activatedCount)
	}
	t.Logf("✅ Threshold > correctly filters aggregation results")
}
// TestAggregationThreshold_GreaterThanOrEqual tests aggregation with >= threshold
func TestAggregationThreshold_GreaterThanOrEqual(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "threshold_gte.tsd")
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
action print(message: string)
rule decent_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id AND avg_sal >= 50000 ==> print("Decent average")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "Sales",
		},
	}
	// Average exactly at threshold (should fire with >=)
	emp1 := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"deptId": "d1",
			"salary": 50000, // avg = 50000, should fire (50000 >= 50000 is true)
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation when avg (50000) >= threshold (50000), got %d", activatedCount)
	}
	t.Logf("✅ Threshold >= correctly handles equal values")
}
// TestAggregationThreshold_LessThan tests aggregation with < threshold
func TestAggregationThreshold_LessThan(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "threshold_lt.tsd")
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
action print(message: string)
rule low_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id AND avg_sal < 40000 ==> print("Low average salary")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "Support",
		},
	}
	// Low average (should fire)
	emp1 := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"deptId": "d1",
			"salary": 35000, // avg = 35000, should fire (35000 < 40000)
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation when avg (35000) < threshold (40000), got %d", activatedCount)
	}
	t.Logf("✅ Threshold < correctly filters aggregation results")
}
// TestAggregationThreshold_MultipleConditions tests aggregation with multiple threshold conditions
func TestAggregationThreshold_MultipleConditions(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "threshold_multi.tsd")
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
action print(message: string)
rule mid_range_avg : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id AND avg_sal > 40000 AND avg_sal < 60000 ==> print("Mid-range average")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "Marketing",
		},
	}
	// Test case 1: Below range (should NOT fire)
	emp1 := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"deptId": "d1",
			"salary": 30000,
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount != 0 {
		t.Errorf("Expected 0 activations when avg (30000) is below range, got %d", activatedCount)
	}
	// Test case 2: Within range (should fire)
	emp2 := Fact{
		ID:   "E2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e2",
			"deptId": "d1",
			"salary": 60000, // avg = (30000 + 60000) / 2 = 45000, should fire
		},
	}
	network.SubmitFact(&emp2)
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation when avg (45000) is in range (40000-60000), got %d", activatedCount)
	}
	t.Logf("✅ Multiple threshold conditions work correctly")
}
// TestAggregationThreshold_COUNT tests COUNT aggregation with threshold
func TestAggregationThreshold_COUNT(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "threshold_count.tsd")
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
action print(message: string)
rule large_dept : {d: Department, emp_count: COUNT(e.id)} / {e: Employee} / e.deptId == d.id AND emp_count >= 3 ==> print("Large department")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "HR",
		},
	}
	network.SubmitFact(&dept)
	// Add 2 employees (should NOT fire, count < 3)
	for i := 1; i <= 2; i++ {
		emp := Fact{
			ID:   "E" + string(rune('0'+i)),
			Type: "Employee",
			Fields: map[string]interface{}{
				"id":     "e" + string(rune('0'+i)),
				"deptId": "d1",
				"salary": 50000,
			},
		}
		network.SubmitFact(&emp)
	}
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount != 0 {
		t.Errorf("Expected 0 activations when count (2) < threshold (3), got %d", activatedCount)
	}
	// Add 3rd employee (should fire, count >= 3)
	emp3 := Fact{
		ID:   "E3",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e3",
			"deptId": "d1",
			"salary": 50000,
		},
	}
	network.SubmitFact(&emp3)
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation when count (3) >= threshold (3), got %d", activatedCount)
	}
	t.Logf("✅ COUNT aggregation with threshold works correctly")
}
// TestAggregationThreshold_NoThreshold tests that aggregation without threshold fires always
func TestAggregationThreshold_NoThreshold(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "no_threshold.tsd")
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
action print(message: string)
rule any_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Any average")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "IT",
		},
	}
	emp := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"deptId": "d1",
			"salary": 1000, // Very low salary, but should still fire with no threshold
		},
	}
	network.SubmitFact(&dept)
	network.SubmitFact(&emp)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activatedCount += len(terminalNode.GetMemory().Tokens)
	}
	if activatedCount < 1 {
		t.Errorf("Expected at least 1 activation with no threshold condition, got %d", activatedCount)
	}
	t.Logf("✅ Aggregation without threshold fires unconditionally")
}