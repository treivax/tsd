// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"strings"

	"github.com/treivax/tsd/rete"
)

type TestResult struct {
	Name        string
	Success     bool
	Error       error
	Description string
	Facts       int
	Actions     int
	ExpectedErr bool
}

func main() {
	fmt.Println("🧪 COMPREHENSIVE ACTION ARGUMENTS TEST SUITE")
	fmt.Println("============================================")

	testSuite := []struct {
		name        string
		constraint  string
		facts       string
		description string
		expectError bool
	}{
		{
			name:        "Variable Action Test",
			description: "Test de variables et d'actions",
			constraint:  "../../constraint/test/integration/variable_action_test.constraint",
			facts:       "../../constraint/test/integration/variable_action_test.facts",
		},
		{
			name:        "Comprehensive Args Test",
			description: "Test complet des arguments",
			constraint:  "../../constraint/test/integration/comprehensive_args_test.constraint",
			facts:       "../../constraint/test/integration/comprehensive_args_test.facts",
		},
		{
			name:        "Error Args Test",
			description: "Test des erreurs d'arguments",
			constraint:  "../../constraint/test/integration/error_args_test.constraint",
			facts:       "../../constraint/test/integration/error_args_test.facts",
		},
	}

	results := make([]TestResult, 0, len(testSuite))

	fmt.Printf("🚀 Running %d test cases...\n\n", len(testSuite))

	for i, test := range testSuite {
		fmt.Printf("📝 Test %d/%d: %s\n", i+1, len(testSuite), test.name)
		fmt.Printf("   %s\n", test.description)
		fmt.Printf("   Constraint: %s\n", test.constraint)
		fmt.Printf("   Facts: %s\n", test.facts)

		result := runSingleTest(test.name, test.constraint, test.facts, test.description, test.expectError)
		results = append(results, result)

		if result.Success {
			fmt.Printf("   ✅ PASSED\n")
		} else {
			fmt.Printf("   ❌ FAILED: %v\n", result.Error)
		}
		fmt.Println()
	}

	// Print summary
	printTestSummary(results)
}

func runSingleTest(name, constraintFile, factsFile, description string, expectError bool) TestResult {
	result := TestResult{
		Name:        name,
		Description: description,
		ExpectedErr: expectError,
	}

	// Create pipeline and storage
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	// Try to build network
	network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts(
		constraintFile,
		factsFile,
		storage,
	)

	// Handle expected errors
	if expectError {
		if err != nil {
			result.Success = true
			result.Error = nil
			fmt.Printf("   ✓ Expected error detected: %v\n", err)
			return result
		} else {
			result.Success = false
			result.Error = fmt.Errorf("expected error but test succeeded")
			return result
		}
	}

	// Handle unexpected errors
	if err != nil {
		result.Success = false
		result.Error = err
		return result
	}

	// Test succeeded - validate results
	result.Facts = len(facts)
	result.Actions = len(network.TerminalNodes)

	// Additional validation for successful cases
	err = validateTestResults(network, facts)
	if err != nil {
		result.Success = false
		result.Error = fmt.Errorf("validation failed: %w", err)
		return result
	}

	result.Success = true
	fmt.Printf("   ✓ Network built: %d facts, %d terminal nodes\n", result.Facts, result.Actions)

	// Display some actions triggered for verification
	displayTriggeredActions(network)

	return result
}

func validateTestResults(network *rete.ReteNetwork, facts []*rete.Fact) error {
	// Validate network structure
	if len(network.TypeNodes) == 0 {
		return fmt.Errorf("no type nodes created")
	}

	if len(network.TerminalNodes) == 0 {
		return fmt.Errorf("no terminal nodes created")
	}

	// Validate terminal nodes have proper actions
	for id, terminal := range network.TerminalNodes {
		if terminal.Action == nil {
			return fmt.Errorf("terminal node %s has no action", id)
		}

		if terminal.Action.Job.Name == "" {
			return fmt.Errorf("terminal node %s has empty action name", id)
		}
	}

	// Validate facts were processed
	if len(facts) == 0 {
		return fmt.Errorf("no facts were parsed")
	}

	return nil
}

func displayTriggeredActions(network *rete.ReteNetwork) {
	fmt.Printf("   📋 Terminal nodes and their actions:\n")
	for id, terminal := range network.TerminalNodes {
		if terminal.Action != nil {
			fmt.Printf("      - %s: %s(%v)\n",
				id,
				terminal.Action.Job.Name,
				formatArgs(terminal.Action.Job.Args))
		}
	}
}

func formatArgs(args []interface{}) string {
	formatted := make([]string, len(args))
	for i, arg := range args {
		if argMap, ok := arg.(map[string]interface{}); ok {
			if argType, hasType := argMap["type"]; hasType {
				switch argType {
				case "variable":
					if name, hasName := argMap["name"]; hasName {
						formatted[i] = fmt.Sprintf("var:%s", name)
					}
				case "fieldAccess":
					if obj, hasObj := argMap["object"]; hasObj {
						if field, hasField := argMap["field"]; hasField {
							formatted[i] = fmt.Sprintf("field:%s.%s", obj, field)
						}
					}
				}
			}
		} else {
			formatted[i] = fmt.Sprintf("%v", arg)
		}
	}
	return strings.Join(formatted, ", ")
}

func printTestSummary(results []TestResult) {
	fmt.Println("📊 TEST SUMMARY")
	fmt.Println("===============")

	passed := 0
	failed := 0

	for _, result := range results {
		if result.Success {
			passed++
			fmt.Printf("✅ %s\n", result.Name)
			if !result.ExpectedErr {
				fmt.Printf("   📈 %d facts processed, %d actions available\n", result.Facts, result.Actions)
			} else {
				fmt.Printf("   🛡️ Error correctly detected and handled\n")
			}
		} else {
			failed++
			fmt.Printf("❌ %s\n", result.Name)
			fmt.Printf("   💥 %v\n", result.Error)
		}
	}

	fmt.Printf("\n🎯 FINAL RESULTS: %d passed, %d failed out of %d total tests\n",
		passed, failed, len(results))

	if failed == 0 {
		fmt.Println("🎉 ALL TESTS PASSED! Action arguments system is working correctly.")
	} else {
		fmt.Println("⚠️ Some tests failed. Please review the errors above.")
	}
}
