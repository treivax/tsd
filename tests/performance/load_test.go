// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

//go:build performance

package performance

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/treivax/tsd/tests/shared/testutil"
)

const (
	// DefaultLoadTestTimeout timeout par défaut pour les tests de charge
	DefaultLoadTestTimeout = 30 * time.Second

	// ExtendedLoadTestTimeout timeout pour les tests de charge étendus
	ExtendedLoadTestTimeout = 60 * time.Second

	// LongLoadTestTimeout timeout pour les tests de charge longs
	LongLoadTestTimeout = 120 * time.Second

	// MinExpectedActivations nombre minimum d'activations attendues
	MinExpectedActivations = 1
)

func TestLoad_100Facts(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	// Generate rule with 100 facts
	rule := generateRuleWithFacts(100)

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFile(t, tempFile)

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 100)

	t.Logf("Successfully processed 100 facts: %d activations", result.Activations)
}

func TestLoad_1000Facts(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	rule := generateRuleWithFacts(1000)

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        DefaultLoadTestTimeout,
		MaxActivations: -1,
	})

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 1000)

	t.Logf("Successfully processed 1000 facts: %d activations", result.Activations)
}

func TestLoad_5000Facts(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	rule := generateRuleWithFacts(5000)

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        ExtendedLoadTestTimeout,
		MaxActivations: -1,
	})

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 5000)

	t.Logf("Successfully processed 5000 facts: %d activations", result.Activations)
}

func TestLoad_10000Facts(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	rule := generateRuleWithFacts(10000)

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        LongLoadTestTimeout,
		MaxActivations: -1,
	})

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 10000)

	t.Logf("Successfully processed 10000 facts: %d activations", result.Activations)
}

func TestLoad_MultipleRulesWithFacts(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	// Multiple rules with 500 facts each
	rule := `type Person(name: string, age: number, salary: number, active: bool)

action print(arg: string)

rule r1 : {p: Person} / p.age > 18 ==> print("adult")
rule r2 : {p: Person} / p.salary > 50000 ==> print("high_earner")
rule r3 : {p: Person} / p.active == true ==> print("active")
rule r4 : {p: Person} / p.age > 30 and p.salary > 60000 ==> print("senior_high_earner")

`

	for i := 0; i < 500; i++ {
		active := "true"
		if i%3 == 0 {
			active = "false"
		}
		rule += fmt.Sprintf(`Person(name:"Person%d", age:%d, salary:%d, active:%s)
`, i, 20+(i%50), 30000+(i*100), active)
	}

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        ExtendedLoadTestTimeout,
		MaxActivations: -1,
	})

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 500)

	t.Logf("Successfully processed 500 facts with 4 rules: %d activations", result.Activations)
}

func TestLoad_ComplexConstraints(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	// Complex constraint with multiple conditions
	rule := `type Transaction(#id: number, amount: number, category: string, verified: bool, timestamp: number)

action print(arg: string)

rule fraud_detection : {tx: Transaction} /
    tx.amount > 10000 and
    tx.verified == false and
    tx.timestamp > 1000000000
    ==> print("potential_fraud")

rule high_value : {tx: Transaction} /
    tx.amount > 50000 and
    tx.verified == true
    ==> print("high_value_verified")

`

	for i := 0; i < 1000; i++ {
		verified := "true"
		if i%5 == 0 {
			verified = "false"
		}
		categories := []string{"retail", "online", "transfer", "withdrawal"}
		category := categories[i%len(categories)]

		rule += fmt.Sprintf(`Transaction(id:%d, amount:%d, category:"%s", verified:%s, timestamp:%d)
`, i, 5000+(i*100), category, verified, 1000000000+i)
	}

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        ExtendedLoadTestTimeout,
		MaxActivations: -1,
	})

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 1000)

	t.Logf("Complex constraints with 1000 facts: %d activations", result.Activations)
}

func TestLoad_JoinHeavy(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	// NOTE: Les jointures à 3+ variables avec BindingChain fonctionnent correctement
	// (voir join_multi_variable_complex.tsd qui utilise des IDs string)
	// Important: Utiliser string pour les IDs de jointure, pas number
	// (les comparaisons number==number dans les jointures ont un bug de conversion int/float)

	// Create scenario with joins between multiple types
	rule := `type Employee(#id: string, name: string, dept_id: string)
type Department(#id: string, name: string, budget: number)
type Project(#id: string, dept_id: string, name: string)

action print(arg: string)

rule emp_dept_project : {e: Employee, d: Department, p: Project} /
    e.dept_id == d.id and
    p.dept_id == d.id and
    d.budget > 100000
    ==> print("employee_on_funded_project")

`

	// Generate 100 employees, 10 departments, 50 projects
	for i := 0; i < 100; i++ {
		deptID := (i % 10) + 1
		rule += fmt.Sprintf(`Employee(id:"e%d", name:"Employee%d", dept_id:"d%d")
`, i, i, deptID)
	}

	for i := 1; i <= 10; i++ {
		budget := 100000 + (i * 10000)
		rule += fmt.Sprintf(`Department(id:"d%d", name:"Dept%d", budget:%d)
`, i, i, budget)
	}

	for i := 0; i < 50; i++ {
		deptID := (i % 10) + 1
		rule += fmt.Sprintf(`Project(id:"p%d", dept_id:"d%d", name:"Project%d")
`, i, deptID, i)
	}

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        ExtendedLoadTestTimeout,
		MaxActivations: -1,
		CaptureOutput:  true,
	})

	testutil.AssertNoError(t, result)

	// Count activations from output since they may not persist in terminal memory
	// Each action execution produces "🎯 ACTION EXÉCUTÉE" in the output
	actionCount := 0
	if result.Output != "" {
		lines := strings.Split(result.Output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "🎯 ACTION EXÉCUTÉE") || strings.Contains(line, "ACTION EXÉCUTÉE") {
				actionCount++
			}
		}
	}

	// Use action count if we found any, otherwise use terminal memory count
	if actionCount > 0 {
		result.Activations = actionCount
	}

	testutil.AssertMinActivations(t, result, MinExpectedActivations)

	t.Logf("Join-heavy test (100 employees, 10 depts, 50 projects): %d activations", result.Activations)
}

func TestLoad_IncrementalFactAddition(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	// Test incremental addition performance
	batchSizes := []int{10, 50, 100, 200, 500}

	for _, batchSize := range batchSizes {
		t.Run(fmt.Sprintf("Batch_%d", batchSize), func(t *testing.T) {
			rule := generateRuleWithFacts(batchSize)

			tempFile := testutil.CreateTempTSDFile(t, rule)
			defer testutil.CleanupTempFiles(t, tempFile)

			start := time.Now()
			result := testutil.ExecuteTSDFile(t, tempFile)
			duration := time.Since(start)

			testutil.AssertNoError(t, result)
			testutil.AssertFactCount(t, result, batchSize)

			t.Logf("Batch size %d: %v, %d activations", batchSize, duration, result.Activations)
		})
	}
}

func TestLoad_MemoryStress(t *testing.T) {
	testutil.SkipIfShort(t, "performance tests skipped in short mode")

	// Large facts with multiple fields to stress memory
	rule := `type LargeRecord(
    id: number,
    field1: string,
    field2: string,
    field3: number,
    field4: number,
    field5: bool,
    field6: string,
    field7: number,
    field8: bool
)

action print(arg: string)

rule r1 : {lr: LargeRecord} / lr.id > 0 ==> print("processed")

`

	for i := 0; i < 2000; i++ {
		rule += fmt.Sprintf(`LargeRecord(id:%d, field1:"value%d", field2:"data%d", field3:%d, field4:%d, field5:%s, field6:"extra%d", field7:%d, field8:%s)
`,
			i,
			i, i,
			i*10, i*20,
			fmt.Sprintf("%t", i%2 == 0),
			i,
			i*30,
			fmt.Sprintf("%t", i%3 == 0))
	}

	tempFile := testutil.CreateTempTSDFile(t, rule)
	defer testutil.CleanupTempFiles(t, tempFile)

	result := testutil.ExecuteTSDFileWithOptions(t, tempFile, &testutil.ExecutionOptions{
		Timeout:        LongLoadTestTimeout,
		MaxActivations: -1,
	})

	testutil.AssertNoError(t, result)
	testutil.AssertFactCount(t, result, 2000)

	t.Logf("Memory stress test with 2000 large records: %d activations", result.Activations)
}

// generateRuleWithFacts creates a TSD rule with the specified number of facts
func generateRuleWithFacts(count int) string {
	rule := `type Item(#id: number, value: string, score: number)

action print(arg: string)

rule r1 : {i: Item} / i.score > 0 ==> print("positive_score")

`

	for i := 0; i < count; i++ {
		rule += fmt.Sprintf(`Item(id:%d, value:"item_%d", score:%d)
`, i, i, i%100)
	}

	return rule
}
