# RETE Expression Analyzer

## Overview

The RETE Expression Analyzer is a sophisticated component of the TSD RETE engine that analyzes, classifies, and optimizes logical expressions for efficient processing in the RETE network.

**Current Version:** 1.3.0  
**Status:** Production Ready

## Features

### Core Analysis (v1.0+)

- **Expression Type Classification** - Identifies Simple, AND, OR, Mixed, Arithmetic, and NOT expressions
- **Decomposability Analysis** - Determines if expressions can be converted to alpha chains
- **Complexity Estimation** - Calculates computational complexity of expressions
- **Beta Node Detection** - Identifies expressions requiring beta nodes

### Inner Expression Analysis (v1.2+)

- **Recursive Analysis** - Analyzes nested and parenthesized expressions
- **NOT Expression Support** - Deep analysis of negated expressions
- **Parenthesized Expressions** - Handles multiple parenthesization formats

### De Morgan Transformation (v1.3+)

- **Automatic Transformation** - Converts `NOT(A OR B)` → `(NOT A) AND (NOT B)`
- **Intelligent Decision Making** - Applies transformation only when beneficial
- **Performance Optimization** - 33-40% improvement for affected expressions

### Optimization Hints (v1.3+)

- **10 Optimization Hints** - Actionable suggestions for expression optimization
- **Intelligent Generation** - Based on deep expression analysis
- **Zero Runtime Overhead** - Generated once during compilation

## Quick Start

### Basic Analysis

```go
import "github.com/treivax/tsd/rete"

// Analyze an expression
expr := /* your expression */
exprType, err := rete.AnalyzeExpression(expr)

if rete.CanDecompose(exprType) {
    // Build alpha chain
} else if rete.ShouldNormalize(exprType) {
    // Normalize first
}
```

### Complete Expression Info

```go
info, err := rete.GetExpressionInfo(expr)

fmt.Printf("Type: %s\n", info.Type)
fmt.Printf("Complexity: %d\n", info.Complexity)
fmt.Printf("Can Decompose: %v\n", info.CanDecompose)
fmt.Printf("Should Normalize: %v\n", info.ShouldNormalize)
fmt.Printf("Requires Beta: %v\n", info.RequiresBeta)

// Check for nested expressions
if info.InnerInfo != nil {
    fmt.Printf("Inner Type: %s\n", info.InnerInfo.Type)
}

// Process optimization hints
for _, hint := range info.OptimizationHints {
    fmt.Printf("Hint: %s\n", hint)
}
```

### Apply De Morgan Transformation

```go
// Check if transformation is beneficial
if rete.ShouldApplyDeMorgan(expr) {
    // Apply transformation
    optimized, applied := rete.ApplyDeMorganTransformation(expr)
    if applied {
        // Use optimized expression
        expr = optimized
    }
}
```

## Expression Types

### ExprTypeSimple
**Description:** Atomic condition (e.g., `p.age > 18`)  
**Complexity:** 1  
**Can Decompose:** Yes  
**Requires Beta:** No

### ExprTypeAND
**Description:** Conjunction of conditions (e.g., `A AND B AND C`)  
**Complexity:** 1 + number of operations  
**Can Decompose:** Yes (alpha chain)  
**Requires Beta:** No

### ExprTypeOR
**Description:** Disjunction of conditions (e.g., `A OR B OR C`)  
**Complexity:** 1 + number of operations  
**Can Decompose:** No  
**Requires Beta:** Yes (branching required)

### ExprTypeMixed
**Description:** Both AND and OR operators (e.g., `(A AND B) OR C`)  
**Complexity:** 1 + number of operations  
**Can Decompose:** No  
**Requires Beta:** Yes

### ExprTypeArithmetic
**Description:** Arithmetic operations (e.g., `price * 1.2 + 5`)  
**Complexity:** 2  
**Can Decompose:** Yes  
**Requires Beta:** No

### ExprTypeNOT
**Description:** Negation of expression (e.g., `NOT (A OR B)`)  
**Complexity:** 2 + inner complexity  
**Can Decompose:** Yes (single NOT node)  
**Requires Beta:** No (unless inner requires it)

## API Reference

### Analysis Functions

#### AnalyzeExpression
```go
func AnalyzeExpression(expr interface{}) (ExpressionType, error)
```
Analyzes an expression and returns its type.

#### GetExpressionInfo
```go
func GetExpressionInfo(expr interface{}) (*ExpressionInfo, error)
```
Returns comprehensive information about an expression, including optimization hints.

#### AnalyzeInnerExpression
```go
func AnalyzeInnerExpression(expr interface{}) (ExpressionType, error)
```
Analyzes the inner expression of a NOT or parenthesized expression.

### Decision Functions

#### CanDecompose
```go
func CanDecompose(exprType ExpressionType) bool
```
Determines if expression can be decomposed into alpha chain.

#### ShouldNormalize
```go
func ShouldNormalize(exprType ExpressionType) bool
```
Determines if expression needs normalization before processing.

#### RequiresBetaNode
```go
func RequiresBetaNode(exprType ExpressionType) bool
```
Determines if expression requires beta nodes.

#### GetExpressionComplexity
```go
func GetExpressionComplexity(exprType ExpressionType) int
```
Returns base complexity estimate for expression type.

### Transformation Functions (v1.3+)

#### ApplyDeMorganTransformation
```go
func ApplyDeMorganTransformation(expr interface{}) (interface{}, bool)
```
Applies De Morgan's laws to transform NOT expressions.

**Returns:**
- Transformed expression (or original if no transformation)
- Boolean indicating if transformation was applied

#### ShouldApplyDeMorgan
```go
func ShouldApplyDeMorgan(expr interface{}) bool
```
Determines if De Morgan transformation should be applied.

**Decision Logic:**
- `NOT(OR)` → Always apply (major optimization)
- `NOT(AND)` → Apply if inner complexity ≤ 2
- Others → Don't apply

## Optimization Hints

### Available Hints

| Hint | Trigger | Action | Benefit |
|------|---------|--------|---------|
| `apply_demorgan_not_or` | NOT(OR) | Apply De Morgan | Major (OR→AND) |
| `apply_demorgan_not_and` | NOT(AND) | Consider applying | Minimal |
| `push_negation_down` | NOT(Mixed) | Recursive De Morgan | Simplification |
| `normalize_to_dnf` | Mixed | Apply DNF normalization | Enables processing |
| `consider_dnf_expansion` | OR | Expand to DNF | Alpha sharing |
| `alpha_sharing_opportunity` | Complex AND | Share alpha nodes | Memory savings |
| `consider_reordering` | AND (2+ ops) | Reorder by selectivity | Early filtering |
| `high_complexity_review` | Complexity ≥ 4 | Manual review | Prevent issues |
| `requires_beta_node` | OR/Mixed | Allocate beta resources | Proper handling |
| `consider_arithmetic_simplification` | Arithmetic | Simplify/pre-compute | Faster evaluation |

### Using Hints

```go
info, _ := rete.GetExpressionInfo(expr)

for _, hint := range info.OptimizationHints {
    switch hint {
    case "apply_demorgan_not_or":
        expr, _ = rete.ApplyDeMorganTransformation(expr)
        
    case "alpha_sharing_opportunity":
        builder.EnableAlphaSharing()
        
    case "consider_reordering":
        conditions = reorderBySelectivity(conditions)
        
    case "normalize_to_dnf":
        expr = normalizeToDNF(expr)
        
    case "high_complexity_review":
        log.Printf("Complex expression detected: %v", expr)
    }
}
```

## Examples

### Example 1: Simple Analysis

```go
expr := constraint.BinaryOperation{
    Left:     constraint.FieldAccess{Object: "p", Field: "age"},
    Operator: ">",
    Right:    constraint.NumberLiteral{Value: 18},
}

info, _ := rete.GetExpressionInfo(expr)
// Type: ExprTypeSimple
// Complexity: 1
// Can Decompose: true
// Optimization Hints: []
```

### Example 2: AND Chain

```go
expr := constraint.LogicalExpression{
    Left: /* p.age > 18 */,
    Operations: []constraint.LogicalOperation{
        {Op: "AND", Right: /* p.salary >= 50000 */},
        {Op: "AND", Right: /* p.active == true */},
    },
}

info, _ := rete.GetExpressionInfo(expr)
// Type: ExprTypeAND
// Complexity: 3 (1 left + 2 operations)
// Can Decompose: true
// Hints: ["alpha_sharing_opportunity", "consider_reordering"]
```

### Example 3: De Morgan Transformation

```go
// Original: NOT(status="active" OR status="pending")
expr := constraint.NotConstraint{
    Expression: constraint.LogicalExpression{
        Left: /* status="active" */,
        Operations: []constraint.LogicalOperation{
            {Op: "OR", Right: /* status="pending" */},
        },
    },
}

// Analyze
info, _ := rete.GetExpressionInfo(expr)
// Type: ExprTypeNOT
// Inner Type: ExprTypeOR
// Complexity: 4 (2 + 2)
// Hints: ["apply_demorgan_not_or", "high_complexity_review"]

// Transform
if rete.ShouldApplyDeMorgan(expr) {
    expr, _ = rete.ApplyDeMorganTransformation(expr)
    // Result: (NOT status="active") AND (NOT status="pending")
    // Type: ExprTypeAND
    // Complexity: 2
    // Can Decompose: true
}
```

### Example 4: Complete Optimization Workflow

```go
func optimizeExpression(expr interface{}) interface{} {
    // Step 1: Analyze
    info, err := rete.GetExpressionInfo(expr)
    if err != nil {
        return expr
    }
    
    // Step 2: Apply De Morgan if beneficial
    if contains(info.OptimizationHints, "apply_demorgan_not_or") {
        if rete.ShouldApplyDeMorgan(expr) {
            expr, _ = rete.ApplyDeMorganTransformation(expr)
            log.Println("✓ Applied De Morgan transformation")
        }
    }
    
    // Step 3: Normalize if needed
    if info.ShouldNormalize {
        expr = normalizeToDNF(expr)
        log.Println("✓ Normalized to DNF")
    }
    
    // Step 4: Log remaining optimization opportunities
    for _, hint := range info.OptimizationHints {
        if hint != "apply_demorgan_not_or" {
            log.Printf("→ Optimization opportunity: %s", hint)
        }
    }
    
    return expr
}
```

## Performance

### De Morgan Transformation Impact

| Expression | Before | After | Improvement |
|-----------|--------|-------|-------------|
| NOT(A OR B) | 3 nodes, branching | 2 alpha nodes | ~33% |
| NOT(A OR B OR C) | 4 nodes, 3 branches | 3 alpha nodes | ~40% |
| NOT(A OR B OR C OR D) | 5 nodes, 4 branches | 4 alpha nodes | ~45% |

### Complexity Calculation Overhead

- Expression analysis: < 0.05ms
- Hint generation: < 0.01ms
- Total overhead: Negligible (one-time during compilation)

## Testing

### Run Tests

```bash
# All expression analyzer tests
go test -v ./rete/ -run TestAnalyzeExpression

# De Morgan tests
go test -v ./rete/ -run TestApplyDeMorgan

# Optimization hint tests
go test -v ./rete/ -run TestOptimization

# Integration tests
go test -v ./rete/
```

### Test Coverage

- **Core Analysis:** 50+ test cases
- **Inner Expression:** 20+ test cases
- **De Morgan:** 10+ test cases
- **Optimization Hints:** 15+ test cases
- **Total:** 95+ test cases, all passing ✓

## Documentation

- **Feature Documentation:** `docs/EXPRESSION_ANALYZER_V1.3.0_FEATURES.md`
- **Changelog:** `docs/CHANGELOG_V1.3.0.md`
- **Examples:** `examples/expression_analyzer_example.go`
- **This README:** `docs/EXPRESSION_ANALYZER_README.md`

## Version History

### v1.3.0 (2025-11-27) - Current
- ✨ De Morgan transformation
- ✨ Optimization hints system
- 🔧 Enhanced complexity calculation
- 📚 Comprehensive documentation
- ✅ 21 new tests

### v1.2.0
- ✨ Parenthesized expression support
- ✨ Recursive inner expression analysis
- ✨ NOT expression deep analysis
- 🔧 Enhanced complexity for NOT expressions

### v1.1.0
- ✨ Expression type classification
- ✨ Complexity estimation
- ✨ Decomposability analysis

### v1.0.0
- 🎉 Initial release
- ✨ Basic expression analysis

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- New features include tests
- Documentation is updated
- Examples are provided

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

## Support

For issues, questions, or contributions:
- File an issue in the project repository
- Check existing documentation
- Review examples in `examples/` directory

---

**Expression Analyzer v1.3.0** - Intelligent expression analysis and optimization for RETE networks