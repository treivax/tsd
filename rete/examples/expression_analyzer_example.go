//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("=== Expression Analyzer Example ===\n")

	// Example 1: Simple condition
	fmt.Println("Example 1: Simple Condition")
	simpleExpr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
		Operator: ">",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
	}
	analyzeAndPrint("p.age > 18", simpleExpr)

	// Example 2: AND expression
	fmt.Println("\nExample 2: AND Expression")
	andExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
				},
			},
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
		},
	}
	analyzeAndPrint("p.age > 18 AND p.salary >= 50000 AND p.active == true", andExpr)

	// Example 3: OR expression
	fmt.Println("\nExample 3: OR Expression")
	orExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
			Operator: "==",
			Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "active"},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
					Operator: "==",
					Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "pending"},
				},
			},
		},
	}
	analyzeAndPrint("p.status == 'active' OR p.status == 'pending'", orExpr)

	// Example 4: Mixed expression
	fmt.Println("\nExample 4: Mixed Expression (AND + OR)")
	mixedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
			Operator: ">",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
					Operator: ">=",
					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
				},
			},
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "vip"},
					Operator: "==",
					Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
				},
			},
		},
	}
	analyzeAndPrint("(p.age > 18 AND p.salary >= 50000) OR p.vip == true", mixedExpr)

	// Example 5: Arithmetic expression
	fmt.Println("\nExample 5: Arithmetic Expression")
	arithmeticExpr := constraint.BinaryOperation{
		Type:     "binaryOperation",
		Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "price"},
		Operator: "*",
		Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 1.2},
	}
	analyzeAndPrint("p.price * 1.2", arithmeticExpr)

	// Example 6: NOT expression
	fmt.Println("\nExample 6: NOT Expression")
	notExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
			Operator: "==",
			Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
		},
	}
	analyzeAndPrint("NOT p.active == true", notExpr)

	// Example 7: Complex NOT expression (NOT with AND inside)
	fmt.Println("\nExample 7: NOT with Complex Expression")
	complexNotExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
						Operator: "<",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
					},
				},
			},
		},
	}
	analyzeAndPrint("NOT (p.age > 18 AND p.salary < 50000)", complexNotExpr)

	// Example 8: Parenthesized expression
	fmt.Println("\nExample 8: Parenthesized Expression")
	parenthesizedExpr := map[string]interface{}{
		"type": "parenthesized",
		"expression": map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type":     "binaryOperation",
				"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
				"operator": ">",
				"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
			},
			"operations": []interface{}{
				map[string]interface{}{
					"op": "AND",
					"right": map[string]interface{}{
						"type":     "binaryOperation",
						"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "active"},
						"operator": "==",
						"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
					},
				},
			},
		},
	}
	analyzeAndPrint("(p.age > 18 AND p.active == true)", parenthesizedExpr)

	// Example 9: NOT with parenthesized expression
	fmt.Println("\nExample 9: NOT with Parenthesized Expression")
	notParenthesizedExpr := map[string]interface{}{
		"type": "not",
		"expr": map[string]interface{}{
			"type": "parenthesized",
			"expression": map[string]interface{}{
				"type": "logicalExpr",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "age"},
					"operator": ">",
					"right":    map[string]interface{}{"type": "numberLiteral", "value": 18},
				},
				"operations": []interface{}{
					map[string]interface{}{
						"op": "OR",
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"left":     map[string]interface{}{"type": "fieldAccess", "object": "p", "field": "vip"},
							"operator": "==",
							"right":    map[string]interface{}{"type": "booleanLiteral", "value": true},
						},
					},
				},
			},
		},
	}
	analyzeAndPrint("NOT (p.age > 18 OR p.vip == true)", notParenthesizedExpr)

	// Example 10: Inner expression analysis
	fmt.Println("\nExample 10: Inner Expression Analysis for NOT")
	info, err := rete.GetExpressionInfo(complexNotExpr)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Expression: NOT (p.age > 18 AND p.salary < 50000)\n")
		fmt.Printf("  Outer type: %s\n", info.Type)
		fmt.Printf("  Outer complexity: %d\n", info.Complexity)
		if info.InnerInfo != nil {
			fmt.Printf("  Inner type: %s\n", info.InnerInfo.Type)
			fmt.Printf("  Inner complexity: %d\n", info.InnerInfo.Complexity)
			fmt.Printf("  Inner can decompose: %v\n", info.InnerInfo.CanDecompose)
			fmt.Printf("  Inner requires beta: %v\n", info.InnerInfo.RequiresBeta)
			fmt.Println("  ✓ Inner expression analyzed recursively")
		}
	}

	// Example 11: Processing decision based on analysis
	fmt.Println("\n=== Processing Decision Example ===")
	expressions := []struct {
		name string
		expr interface{}
	}{
		{"Simple", simpleExpr},
		{"AND", andExpr},
		{"OR", orExpr},
		{"Mixed", mixedExpr},
		{"Arithmetic", arithmeticExpr},
		{"NOT", notExpr},
		{"ComplexNOT", complexNotExpr},
		{"Parenthesized", parenthesizedExpr},
		{"NOTParenthesized", notParenthesizedExpr},
	}

	for _, e := range expressions {
		fmt.Printf("\n%s Expression:\n", e.name)
		info, err := rete.GetExpressionInfo(e.expr)
		if err != nil {
			log.Printf("  Error: %v\n", err)
			continue
		}

		// Decide processing strategy
		if info.CanDecompose {
			if info.Type == rete.ExprTypeAND {
				fmt.Println("  → Strategy: Build alpha chain (optimal for AND)")
			} else if info.Type == rete.ExprTypeArithmetic {
				fmt.Println("  → Strategy: Build arithmetic evaluation chain")
			} else if info.Type == rete.ExprTypeNOT {
				fmt.Println("  → Strategy: Create alpha node with negation flag")
			} else {
				fmt.Println("  → Strategy: Create single alpha node")
			}
		} else {
			if info.ShouldNormalize {
				fmt.Println("  → Strategy: Normalize first, then process")
			}
			if info.RequiresBeta {
				fmt.Println("  → Strategy: Use beta nodes for branching")
			}
		}

		fmt.Printf("  Estimated complexity: %d\n", info.Complexity)
	}

	// Example 12: De Morgan Transformation - NOT(A OR B)
	fmt.Println("\n=== De Morgan Transformation Examples ===")
	fmt.Println("\nExample 12: De Morgan - NOT(A OR B) -> (NOT A) AND (NOT B)")
	notOrExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
						Operator: "<",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
					},
				},
			},
		},
	}
	fmt.Println("Original: NOT (p.age > 18 OR p.salary < 50000)")
	originalType, _ := rete.AnalyzeExpression(notOrExpr)
	fmt.Printf("  Original type: %s\n", originalType)

	transformed, applied := rete.ApplyDeMorganTransformation(notOrExpr)
	if applied {
		fmt.Println("  ✓ De Morgan transformation applied!")
		transformedType, _ := rete.AnalyzeExpression(transformed)
		fmt.Printf("  Transformed type: %s\n", transformedType)
		fmt.Println("  Result: (NOT p.age > 18) AND (NOT p.salary < 50000)")
	}

	// Example 13: De Morgan Transformation - NOT(A AND B)
	fmt.Println("\nExample 13: De Morgan - NOT(A AND B) -> (NOT A) OR (NOT B)")
	notAndExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
				Operator: ">",
				Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "AND",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
						Operator: ">=",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
					},
				},
			},
		},
	}
	fmt.Println("Original: NOT (p.age > 18 AND p.salary >= 50000)")
	originalType2, _ := rete.AnalyzeExpression(notAndExpr)
	fmt.Printf("  Original type: %s\n", originalType2)

	transformed2, applied2 := rete.ApplyDeMorganTransformation(notAndExpr)
	if applied2 {
		fmt.Println("  ✓ De Morgan transformation applied!")
		transformedType2, _ := rete.AnalyzeExpression(transformed2)
		fmt.Printf("  Transformed type: %s\n", transformedType2)
		fmt.Println("  Result: (NOT p.age > 18) OR (NOT p.salary >= 50000)")
	}

	// Example 14: Should Apply De Morgan Decision
	fmt.Println("\nExample 14: De Morgan Decision Logic")
	testExprs := []struct {
		name string
		expr interface{}
	}{
		{"NOT(A OR B)", notOrExpr},
		{"NOT(A AND B)", notAndExpr},
		{"NOT(simple)", notExpr},
		{"Simple", simpleExpr},
	}

	for _, test := range testExprs {
		shouldApply := rete.ShouldApplyDeMorgan(test.expr)
		fmt.Printf("  %s: should apply = %v\n", test.name, shouldApply)
	}

	// Example 15: Optimization Hints
	fmt.Println("\n=== Optimization Hints Examples ===")
	fmt.Println("\nExample 15: Optimization Hints for Various Expressions")

	hintExprs := []struct {
		name string
		expr interface{}
	}{
		{"NOT(A OR B)", notOrExpr},
		{"NOT(A AND B)", notAndExpr},
		{"Mixed Expression", mixedExpr},
		{"Complex AND", andExpr},
		{"Simple OR", orExpr},
		{"Arithmetic", arithmeticExpr},
	}

	for _, test := range hintExprs {
		info, err := rete.GetExpressionInfo(test.expr)
		if err != nil {
			log.Printf("  Error: %v\n", err)
			continue
		}

		fmt.Printf("\n%s:\n", test.name)
		fmt.Printf("  Type: %s, Complexity: %d\n", info.Type, info.Complexity)
		if len(info.OptimizationHints) > 0 {
			fmt.Println("  Optimization Hints:")
			for _, hint := range info.OptimizationHints {
				fmt.Printf("    - %s\n", hint)
			}
		} else {
			fmt.Println("  No optimization hints (already optimal)")
		}
	}

	// Example 16: Complete Workflow with Hints
	fmt.Println("\n=== Complete Optimization Workflow ===")
	fmt.Println("\nExample 16: Analyze -> Check Hints -> Apply Optimizations")

	workflowExpr := constraint.NotConstraint{
		Type: "notConstraint",
		Expression: constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
				Operator: "==",
				Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "active"},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "pending"},
					},
				},
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "review"},
					},
				},
			},
		},
	}

	fmt.Println("Step 1: Analyze expression")
	fmt.Println("  Expression: NOT (status='active' OR status='pending' OR status='review')")

	workflowInfo, err := rete.GetExpressionInfo(workflowExpr)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("  Type: %s\n", workflowInfo.Type)
	fmt.Printf("  Complexity: %d\n", workflowInfo.Complexity)
	fmt.Printf("  Can Decompose: %v\n", workflowInfo.CanDecompose)
	fmt.Printf("  Requires Beta: %v\n", workflowInfo.RequiresBeta)

	fmt.Println("\nStep 2: Check optimization hints")
	if len(workflowInfo.OptimizationHints) > 0 {
		for _, hint := range workflowInfo.OptimizationHints {
			fmt.Printf("  Hint: %s\n", hint)
		}
	}

	fmt.Println("\nStep 3: Apply recommended optimizations")
	if rete.ShouldApplyDeMorgan(workflowExpr) {
		fmt.Println("  → Applying De Morgan transformation...")
		optimized, applied := rete.ApplyDeMorganTransformation(workflowExpr)
		if applied {
			optimizedInfo, err := rete.GetExpressionInfo(optimized)
			if err != nil {
				log.Printf("Error getting optimized info: %v\n", err)
				return
			}
			fmt.Println("  ✓ Transformation successful!")
			fmt.Printf("  New type: %s\n", optimizedInfo.Type)
			fmt.Printf("  New complexity: %d\n", optimizedInfo.Complexity)
			fmt.Printf("  Can decompose now: %v\n", optimizedInfo.CanDecompose)
			fmt.Println("  Result: (NOT status='active') AND (NOT status='pending') AND (NOT status='review')")
			fmt.Println("  → This can now be processed as an alpha chain!")
		}
	}

	fmt.Println("\n=== Example Complete ===")
}

func analyzeAndPrint(description string, expr interface{}) {
	fmt.Printf("Expression: %s\n", description)

	// Analyze the expression
	exprType, err := rete.AnalyzeExpression(expr)
	if err != nil {
		log.Printf("Error analyzing expression: %v\n", err)
		return
	}

	// Get comprehensive info
	info, err := rete.GetExpressionInfo(expr)
	if err != nil {
		log.Printf("Error getting expression info: %v\n", err)
		return
	}

	// Print results
	fmt.Printf("  Type: %s\n", exprType)
	fmt.Printf("  Can decompose into chain: %v\n", info.CanDecompose)
	fmt.Printf("  Should normalize: %v\n", info.ShouldNormalize)
	fmt.Printf("  Requires beta node: %v\n", info.RequiresBeta)
	fmt.Printf("  Complexity level: %d\n", info.Complexity)

	// Print inner expression info if available
	if info.InnerInfo != nil {
		fmt.Printf("  Inner expression type: %s\n", info.InnerInfo.Type)
		fmt.Printf("  Inner complexity: %d\n", info.InnerInfo.Complexity)
	}

	// Print optimization hints if available
	if len(info.OptimizationHints) > 0 {
		fmt.Printf("  Optimization hints: %v\n", info.OptimizationHints)
	}

	// Processing recommendation
	if rete.CanDecompose(exprType) {
		fmt.Println("  ✓ This expression can be decomposed into an alpha chain")
	} else {
		fmt.Println("  ✗ This expression needs special handling (beta nodes or normalization)")
	}
}
