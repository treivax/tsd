// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// This file serves as the main entry point for nested OR normalization.
// The implementation has been split into focused modules:
//
// - nested_or_normalizer_analysis.go: Complexity analysis and type definitions
// - nested_or_normalizer_flattening.go: OR flattening operations
// - nested_or_normalizer_dnf.go: DNF transformation logic
// - nested_or_normalizer_helpers.go: Main normalization orchestration
//
// Public API:
// - AnalyzeNestedOR: Analyze expression complexity
// - FlattenNestedOR: Flatten nested OR expressions
// - TransformToDNF: Transform to Disjunctive Normal Form
// - NormalizeNestedOR: Complete normalization pipeline
