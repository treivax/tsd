# Expression Analyzer - Summary

## Quick Overview

The Expression Analyzer is a decision-making component that analyzes expressions in the RETE network to determine their type and processing requirements. It supports parenthesized expressions and recursively analyzes inner expressions within NOT constraints (v1.2.0+).

## What It Does

- **Identifies expression types**: Simple, AND, OR, Mixed, Arithmetic
- **Determines decomposability**: Can the expression be converted to an alpha chain?
- **Evaluates normalization needs**: Does the expression need to be normalized?
- **Assesses complexity**: Provides a complexity score (1-4)
- **Guides processing strategy**: Helps decide between alpha chains, beta nodes, or normalization

## Key Types

### ExpressionType Constants

```go
ExprTypeSimple      // Single condition: p.age > 18
ExprTypeAND         // AND only: p.age > 18 AND p.salary >= 50000
ExprTypeOR          // OR only: p.status == "A" OR p.status == "B"
ExprTypeMixed       // Mixed: (A AND B) OR C
ExprTypeArithmetic  // Arithmetic: p.price * 1.2
ExprTypeNOT         // Negation: NOT p.active, NOT (A AND B)
```

## Main Functions

### AnalyzeExpression(expr) → (ExpressionType, error)
Analyzes an expression and returns its type.

### CanDecompose(exprType) → bool
Returns true if the expression can be decomposed into an alpha chain.
- ✓ ExprTypeSimple, ExprTypeAND, ExprTypeArithmetic, ExprTypeNOT
- ✗ ExprTypeOR, ExprTypeMixed

### ShouldNormalize(exprType) → bool
Returns true if the expression needs normalization.
- ✓ ExprTypeOR, ExprTypeMixed
- ✗ ExprTypeSimple, ExprTypeAND, ExprTypeArithmetic, ExprTypeNOT

### GetExpressionInfo(expr) → (*ExpressionInfo, error)
Returns comprehensive analysis in one call.

## Quick Example

```go
// Analyze an expression
exprType, err := rete.AnalyzeExpression(expr)
if err != nil {
    return err
}

// Make processing decision
if rete.CanDecompose(exprType) {
    // Build alpha chain
    chain, err := builder.BuildChain(conditions, "p", parent, ruleID)
} else if rete.ShouldNormalize(exprType) {
    // Normalize first
    normalized := normalize(expr)
    // Then process
}
```

## Expression Type Characteristics

| Type          | Decompose | Normalize | Beta Node | Complexity |
|---------------|-----------|-----------|-----------|------------|
| Simple        | ✓         | ✗         | ✗         | 1          |
| AND           | ✓         | ✗         | ✗         | 2          |
| Arithmetic    | ✓         | ✗         | ✗         | 2          |
| NOT           | ✓         | ✗         | ✗         | 2          |
| OR            | ✗         | ✓         | ✓         | 3          |
| Mixed         | ✗         | ✓         | ✓         | 4          |

## Use Cases

### 1. Chain Building Decision
```go
if CanDecompose(exprType) {
    // Use AlphaChainBuilder
} else {
    // Use alternative strategy
}
```

### 2. Normalization Pipeline
```go
if ShouldNormalize(exprType) {
    expr = NormalizeExpression(expr)
}
```

### 3. Optimization Strategy
```go
complexity := GetExpressionComplexity(exprType)
if complexity > 3 {
    // Apply advanced optimization
}
```

## Integration Points

- **AlphaChainBuilder**: Uses `CanDecompose()` to validate expressions
- **Normalizer**: Uses `ShouldNormalize()` to decide when to normalize
- **Network Builder**: Uses `RequiresBetaNode()` for node type selection
- **Optimizer**: Uses `GetExpressionComplexity()` for cost estimation

## Testing

All 15+ test functions pass:
- TestAnalyzeExpression_Simple
- TestAnalyzeExpression_AND
- TestAnalyzeExpression_OR
- TestAnalyzeExpression_Mixed_AND_OR
- TestAnalyzeExpression_Arithmetic
- TestAnalyzeExpression_NOT (5 test cases)
- TestAnalyzeExpression_NOT_Nested (3 test cases)
- TestGetExpressionInfo_NOT
- TestCanDecompose_AllTypes
- TestShouldNormalize_AllTypes
- TestExpressionType_String
- TestGetExpressionComplexity
- TestRequiresBetaNode
- TestGetExpressionInfo
- TestAnalyzeExpression_EdgeCases

Run tests:
```bash
cd rete
go test -v -run TestAnalyzeExpression
go test -v -run TestCanDecompose
go test -v -run TestShouldNormalize
go test -v -run TestAnalyzeExpression_NOT
```

## Files

- `expression_analyzer.go` - Main implementation (500+ lines, v1.2.0)
- `expression_analyzer_test.go` - Comprehensive tests (1950+ lines, v1.2.0)
- `examples/expression_analyzer_example.go` - Usage example with all features
- `EXPRESSION_ANALYZER_README.md` - Full documentation (updated for v1.2.0)
- `EXPRESSION_ANALYZER_CHANGELOG.md` - Version history

## Version History

### v1.2.0 (Current)
- Parenthesized expression support
- Recursive inner expression analysis
- Enhanced `ExpressionInfo` with `InnerInfo` field
- `AnalyzeInnerExpression()` function
- Adjusted complexity calculation for nested expressions

### v1.1.0
- NOT operator support (`ExprTypeNOT`)
- Multiple NOT format support
- Nested NOT handling

### v1.0.0
- Initial implementation
- 5 core expression types
- Basic analysis functions

## License

MIT License - Copyright (c) 2025 TSD Contributors