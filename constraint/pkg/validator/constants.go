// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package validator

// Operator categories and type rules
var (
	// ComparisonOperators définit les opérateurs de comparaison valides
	ComparisonOperators = map[string]bool{
		"==": true,
		"!=": true,
		"<":  true,
		">":  true,
		"<=": true,
		">=": true,
	}

	// LogicalOperators définit les opérateurs logiques valides
	LogicalOperators = map[string]bool{
		"AND": true,
		"OR":  true,
		"NOT": true,
	}

	// ArithmeticOperators définit les opérateurs arithmétiques valides
	ArithmeticOperators = map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
		"%": true,
	}

	// OrderableTypes définit les types qui peuvent être comparés avec <, >, <=, >=
	OrderableTypes = map[string]bool{
		"number":  true,
		"integer": true,
		"string":  true,
	}

	// NumericTypes définit les types numériques pour les opérations arithmétiques
	NumericTypes = map[string]bool{
		"number":  true,
		"integer": true,
	}
)
