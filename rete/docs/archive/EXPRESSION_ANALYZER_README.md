# Expression Analyzer

## Overview

The Expression Analyzer is a component of the RETE network implementation that analyzes expressions to determine their type and how they should be processed in the network. It provides decision-making logic for expression decomposition, normalization requirements, and processing strategies.

## Purpose

When building RETE networks, different types of expressions require different handling strategies:

- **Simple conditions** can be directly converted to alpha nodes
- **AND expressions** can be decomposed into chains of alpha nodes
- **OR expressions** require special handling with branches or duplication
- **Mixed expressions** need normalization before processing
- **Arithmetic chains** can be decomposed into sequential operations

The Expression Analyzer identifies which category an expression belongs to and provides guidance on how to process it. It also supports parenthesized expressions and recursively analyzes inner expressions within NOT constraints.

## Expression Types

### ExprTypeSimple

Represents an atomic condition that can be directly evaluated.

**Examples:**
```go
p.age > 18
p.status == "active"
p.salary >= 50000
```

**Characteristics:**
- Can decompose: ✓
- Requires normalization: ✗
- Requires beta node: ✗
- Complexity: 1

### ExprTypeAND

Represents an expression with only AND operators connecting conditions.

**Examples:**
```go
p.age > 18 AND p.salary >= 50000
p.active == true AND p.verified == true AND p.score > 100
```

**Characteristics:**
- Can decompose: ✓ (into alpha node chain)
- Requires normalization: ✗
- Requires beta node: ✗
- Complexity: 2

### ExprTypeOR

Represents an expression with only OR operators connecting conditions.

**Examples:**
```go
p.status == "active" OR p.status == "pending"
p.type == "A" OR p.type == "B" OR p.type == "C"
```

**Characteristics:**
- Can decompose: ✗
- Requires normalization: ✓
- Requires beta node: ✓
- Complexity: 3

### ExprTypeMixed

Represents an expression mixing AND and OR operators.

**Examples:**
```go
(p.age > 18 AND p.salary >= 50000) OR p.vip == true
p.admin == true OR (p.age > 21 AND p.verified == true)
```

**Characteristics:**
- Can decompose: ✗
- Requires normalization: ✓ (must convert to DNF or CNF)
- Requires beta node: ✓
- Complexity: 4

### ExprTypeArithmetic

Represents arithmetic operations that can be evaluated in sequence.

**Examples:**
```go
p.price * 1.2
p.total + 5
p.value / 2
```

**Characteristics:**
- Can decompose: ✓ (into operation chain)
- Requires normalization: ✗
- Requires beta node: ✗
- Complexity: 2

### ExprTypeNOT

Represents a negation of an expression using the NOT operator.

**Examples:**
```go
NOT p.active
NOT (p.age > 18 AND p.salary < 50000)
NOT (p.status == "active" OR p.status == "pending")
```

**Characteristics:**
- Can decompose: ✓ (into alpha node with negation flag)
- Requires normalization: ✗ (but inner expression may need it)
- Requires beta node: ✗
- Complexity: 2

**Note:** As of version 1.2.0, the analyzer automatically analyzes inner expressions within NOT constraints. The `ExpressionInfo` structure includes an `InnerInfo` field that contains the analysis of the expression being negated. This allows for recursive analysis and better understanding of complex NOT expressions.

## API Reference

### Core Functions

#### AnalyzeExpression

```go
func AnalyzeExpression(expr interface{}) (ExpressionType, error)
```

Analyzes an expression and determines its type.

**Parameters:**
- `expr`: Expression to analyze (can be `constraint.BinaryOperation`, `constraint.LogicalExpression`, map, etc.)

**Returns:**
- `ExpressionType`: The identified expression type
- `error`: Error if analysis fails

**Example:**
```go
expr := constraint.BinaryOperation{
    Type:     "binaryOperation",
    Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
    Operator: ">",
    Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
}

exprType, err := AnalyzeExpression(expr)
if err != nil {
    log.Fatal(err)
}
fmt.Println(exprType) // Output: ExprTypeSimple
```

#### CanDecompose

```go
func CanDecompose(exprType ExpressionType) bool
```

Determines if an expression can be decomposed into an alpha node chain.

**Parameters:**
- `exprType`: The expression type to check

**Returns:**
- `true` if the expression can be decomposed (ExprTypeSimple, ExprTypeAND, ExprTypeArithmetic)
- `false` otherwise (ExprTypeOR, ExprTypeMixed)

**Example:**
```go
if CanDecompose(exprType) {
    // Build alpha chain
    chain, err := builder.BuildChain(conditions, variableName, parentNode, ruleID)
} else {
    // Use alternative strategy (beta nodes, normalization, etc.)
}
```

#### ShouldNormalize

```go
func ShouldNormalize(exprType ExpressionType) bool
```

Determines if an expression requires normalization before processing.

**Parameters:**
- `exprType`: The expression type to check

**Returns:**
- `true` if normalization is needed (ExprTypeOR, ExprTypeMixed)
- `false` otherwise

**Example:**
```go
if ShouldNormalize(exprType) {
    // Convert to DNF or CNF first
    normalized := NormalizeExpression(expr)
    // Then process normalized expression
}
```

### Helper Functions

#### GetExpressionComplexity

```go
func GetExpressionComplexity(exprType ExpressionType) int
```

Returns a complexity estimate for an expression type (1-4 scale).

#### RequiresBetaNode

```go
func RequiresBetaNode(exprType ExpressionType) bool
```

Determines if an expression requires beta nodes for evaluation.

#### GetExpressionInfo

```go
func GetExpressionInfo(expr interface{}) (*ExpressionInfo, error)
```

Returns comprehensive information about an expression in a single call.

**Returns:**
```go
type ExpressionInfo struct {
    Type            ExpressionType
    CanDecompose    bool
    ShouldNormalize bool
    Complexity      int
    RequiresBeta    bool
}
```

**Example:**
```go
info, err := GetExpressionInfo(expr)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Type: %s\n", info.Type)
fmt.Printf("Can decompose: %v\n", info.CanDecompose)
fmt.Printf("Complexity: %d\n", info.Complexity)
```

## Usage Examples

### Example 1: Simple Decision Making

```go
expr := /* your expression */

exprType, err := AnalyzeExpression(expr)
if err != nil {
    return err
}

if CanDecompose(exprType) {
    // Use alpha chain builder
    conditions, _, err := ExtractConditions(expr)
    if err != nil {
        return err
    }
    
    builder := NewAlphaChainBuilder(network, storage)
    chain, err := builder.BuildChain(conditions, "p", rootNode, "rule1")
    if err != nil {
        return err
    }
    
    fmt.Printf("Built chain with %d nodes\n", len(chain.Nodes))
} else {
    // Use alternative strategy
    fmt.Println("Expression requires special handling")
}
```

### Example 2: Comprehensive Analysis

```go
// Analyze multiple expressions
expressions := []interface{}{
    simpleBinaryOp,
    andExpression,
    orExpression,
    mixedExpression,
}

for i, expr := range expressions {
    info, err := GetExpressionInfo(expr)
    if err != nil {
        fmt.Printf("Expression %d: error - %v\n", i, err)
        continue
    }
    
    fmt.Printf("Expression %d:\n", i)
    fmt.Printf("  Type: %s\n", info.Type)
    fmt.Printf("  Complexity: %d\n", info.Complexity)
    fmt.Printf("  Can decompose: %v\n", info.CanDecompose)
    fmt.Printf("  Should normalize: %v\n", info.ShouldNormalize)
    fmt.Printf("  Requires beta: %v\n", info.RequiresBeta)
}
```

### Example 3: Processing Pipeline

```go
func ProcessExpression(expr interface{}, network *ReteNetwork) error {
    // Analyze expression
    exprType, err := AnalyzeExpression(expr)
    if err != nil {
        return fmt.Errorf("analysis failed: %w", err)
    }
    
    // Normalize if needed
    if ShouldNormalize(exprType) {
        expr = NormalizeExpression(expr)
        // Re-analyze after normalization
        exprType, err = AnalyzeExpression(expr)
        if err != nil {
            return fmt.Errorf("re-analysis failed: %w", err)
        }
    }
    
    // Process based on type
    if CanDecompose(exprType) {
        if exprType == ExprTypeNOT {
            return buildAlphaNodeWithNegation(expr, network)
        }
        return buildAlphaChain(expr, network)
    } else if RequiresBetaNode(exprType) {
        return buildBetaNetwork(expr, network)
    }
    
    return nil
}
```

### Example 4: Handling NOT Expressions

```go
// Analyze a NOT expression
notExpr := constraint.NotConstraint{
    Type: "notConstraint",
    Expression: constraint.BinaryOperation{
        Type:     "binaryOperation",
        Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "active"},
        Operator: "==",
        Right:    constraint.BooleanLiteral{Type: "booleanLiteral", Value: true},
    },
}

exprType, err := AnalyzeExpression(notExpr)
if err != nil {
    log.Fatal(err)
}

if exprType == ExprTypeNOT {
    fmt.Println("This is a NOT expression")
    
    // You can still decompose it into an alpha node
    if CanDecompose(exprType) {
        // Build alpha node with negation flag set to true
        node := buildAlphaNodeWithNegation(notExpr.Expression)
    }
}
```

## Integration with RETE Network

## Parenthesized Expressions (v1.2.0+)

The analyzer transparently handles parenthesized expressions by analyzing the inner content:

```go
parenthesizedExpr := map[string]interface{}{
    "type": "parenthesized",
    "expression": map[string]interface{}{
        "type": "logicalExpr",
        "left": ...,
        "operations": [...],
    },
}

exprType, err := rete.AnalyzeExpression(parenthesizedExpr)
// Returns the type of the inner expression (e.g., ExprTypeAND)
```

Supported parenthesized formats:
- `"type": "parenthesized"` with `"expression"` field
- `"type": "parenthesizedExpression"` with `"expr"` field
- `"type": "group"` with `"inner"` field

## Recursive Inner Expression Analysis (v1.2.0+)

For NOT constraints and parenthesized expressions, the analyzer can recursively analyze inner expressions:

```go
// NOT with complex inner expression
notExpr := constraint.NotConstraint{
    Type: "notConstraint",
    Expression: constraint.LogicalExpression{
        Type: "logicalExpr",
        Left: ...,
        Operations: [...], // AND operations
    },
}

info, err := rete.GetExpressionInfo(notExpr)
// info.Type == ExprTypeNOT
// info.InnerInfo != nil
// info.InnerInfo.Type == ExprTypeAND
// info.Complexity == 2 + info.InnerInfo.Complexity
```

The `ExpressionInfo` structure includes:
- `InnerInfo *ExpressionInfo` - Analysis of the inner expression (for NOT and parenthesized expressions)
- Automatically adjusted complexity based on inner expression

Use `AnalyzeInnerExpression()` to directly analyze the inner expression:

```go
innerType, err := rete.AnalyzeInnerExpression(notExpr)
// Returns the type of the inner expression
```

## Integration with RETE Components

The Expression Analyzer integrates seamlessly with other RETE components:

1. **Alpha Chain Builder**: Uses `CanDecompose()` to determine if a chain can be built
2. **Expression Normalizer**: Uses `ShouldNormalize()` to decide when normalization is needed
3. **Network Builder**: Uses `RequiresBetaNode()` to decide on node types
4. **Optimizer**: Uses `GetExpressionComplexity()` for optimization decisions

## Design Decisions

### Why Separate Analysis from Processing?

The analyzer provides information about expressions without actually modifying them. This separation of concerns allows:

- **Flexibility**: Different processing strategies can be chosen based on context
- **Testability**: Analysis logic can be tested independently
- **Performance**: Analysis is lightweight and can be done early in the pipeline
- **Maintainability**: Changes to processing don't affect analysis logic

### Expression Type Granularity

The six expression types provide the right balance between:

- **Specificity**: Enough detail to make informed processing decisions
- **Simplicity**: Not so many types that the API becomes unwieldy
- **Extensibility**: New types can be added without breaking existing code

### NOT Expression Handling

The `ExprTypeNOT` is treated as a decomposable type because:

1. **Alpha Node Compatibility**: NOT can be implemented as an alpha node with a negation flag
2. **Single Operand**: Unlike OR which requires branching, NOT has a single operand
3. **Performance**: Negation can be efficiently evaluated at the alpha node level
4. **Composability**: NOT expressions can contain any other expression type

However, note that the complexity of the inner expression may affect the overall processing strategy. For example, `NOT (A OR B)` can be transformed using De Morgan's law to `(NOT A) AND (NOT B)`, which may be more efficient in some cases.

## Performance Considerations

- **Fast Analysis**: Expression analysis is O(1) for simple expressions and O(n) for logical expressions where n is the number of operations
- **No Mutations**: The analyzer never modifies input expressions
- **Minimal Allocations**: Most functions return scalar values or pre-allocated constants

## Testing

## Version History

### v1.2.0 (Current)
- **Parenthesized Expression Support**: Automatic handling of parenthesized expressions
- **Recursive Inner Analysis**: Automatic analysis of expressions within NOT constraints
- **Enhanced ExpressionInfo**: Added `InnerInfo` field for nested expression analysis
- **Adjusted Complexity Calculation**: Complexity now accounts for inner expressions
- **New Function**: `AnalyzeInnerExpression()` for direct inner expression analysis

### v1.1.0
- **NOT Operator Support**: Added `ExprTypeNOT` for negation expressions
- Support for multiple NOT formats (`notConstraint`, `not`, `negation`)
- Handling of nested NOT expressions (double negation)
- NOT expression characteristics and processing guidelines

### v1.0.0
- Initial implementation with 5 core expression types
- Basic analysis functions (`AnalyzeExpression`, `CanDecompose`, etc.)
- Comprehensive test suite

## Tests

The expression analyzer includes comprehensive tests:

- `TestAnalyzeExpression_Simple`: Tests simple condition detection
- `TestAnalyzeExpression_AND`: Tests AND expression detection
- `TestAnalyzeExpression_OR`: Tests OR expression detection
- `TestAnalyzeExpression_Mixed_AND_OR`: Tests mixed expression detection
- `TestAnalyzeExpression_Arithmetic`: Tests arithmetic operation detection
- `TestAnalyzeExpression_NOT`: Tests NOT expression detection (5 test cases)
- `TestAnalyzeExpression_NOT_Nested`: Tests nested and complex NOT expressions (3 test cases)
- `TestGetExpressionInfo_NOT`: Tests comprehensive info for NOT expressions
- `TestCanDecompose_AllTypes`: Tests decomposition logic for all types
- `TestShouldNormalize_AllTypes`: Tests normalization requirements
- `TestAnalyzeExpression_EdgeCases`: Tests edge cases and error handling

Run tests with:
```bash
cd rete
go test -v -run TestAnalyzeExpression
go test -v -run TestCanDecompose
go test -v -run TestShouldNormalize
go test -v -run TestAnalyzeExpression_NOT
```

All tests pass with 100% success rate.

## License

This implementation is part of TSD and is licensed under the MIT License. See the LICENSE file in the project root for full license text.

## See Also

- `alpha_chain_builder.go`: Uses the analyzer to build alpha node chains
- `alpha_chain_extractor.go`: Extracts conditions from complex expressions
- `normalizer.go`: Normalizes expressions to standard forms
- `ALPHA_CHAIN_BUILDER_README.md`: Documentation for the alpha chain builder