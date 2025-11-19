// Package constraint provides parsing, validation and management of constraint programs.
//
// This package implements a complete constraint language parser using PEG (Parsing Expression Grammar)
// and provides facilities for validating constraint programs, extracting facts, and converting
// constraint ASTs for use with RETE networks.
//
// # Core Components
//
// The package provides several key components:
//   - Parser: PEG-based parser for constraint files (.constraint)
//   - Validator: Semantic validation for constraint programs
//   - Type System: Comprehensive type definitions for AST nodes
//   - Fact Extraction: Support for parsing facts within constraint files
//   - API: High-level functions for parsing and processing constraint files
//
// # Example Usage
//
//	// Parse a constraint file
//	program, err := constraint.ParseConstraintFile("rules.constraint")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Validate the program
//	err = constraint.ValidateConstraintProgram(program)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Extract facts from the program
//	facts := constraint.ExtractFactsFromProgram(program)
//
// # Constraint Language
//
// The constraint language supports:
//   - Type definitions: type Person : <id: string, name: string, age: number>
//   - Constraint expressions: p: Person, p.age > 18 ==> approve_adult(p)
//   - Logical operations: AND, OR, NOT
//   - Quantifiers: EXISTS, aggregate functions (SUM, COUNT, etc.)
//   - Facts: Person(id: "P001", name: "Alice", age: 25)
//   - Function calls: LENGTH(p.name) > 5, UPPER(p.status) == "ACTIVE"
//
// # Integration
//
// This package is designed to work seamlessly with the rete package for rule execution.
// Constraint programs can be converted to RETE network representations for efficient
// pattern matching and rule evaluation.
package constraint
