// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete implements a high-performance RETE algorithm for rule-based systems.
//
// The RETE algorithm is a pattern matching algorithm for implementing rule-based systems.
// It was originally developed by Charles Forgy at Carnegie Mellon University (1974-1979)
// and is now in the public domain. This is an original implementation specifically
// designed for the TSD constraint language.
//
// This package provides a complete implementation with support for:
//   - Alpha nodes (single-fact pattern matching)
//   - Beta nodes (multi-fact joins)
//   - Terminal nodes (rule actions)
//   - Memory management and token propagation
//   - Advanced features (NOT, EXISTS, aggregation)
//
// # Architecture
//
// The RETE network consists of several types of nodes:
//   - RootNode: Entry point for all facts
//   - TypeNode: Filters facts by type
//   - AlphaNode: Single-fact pattern matching and filtering
//   - JoinNode/BetaNode: Multi-fact joins and complex conditions
//   - TerminalNode: Executes actions when patterns match
//   - Advanced nodes: NotNode, ExistsNode, AccumulateNode
//
// # Example Usage
//
//	// Create a RETE network
//	storage := rete.NewMemoryStorage()
//	network := rete.NewReteNetwork(storage)
//
//	// Load rules from constraint AST
//	err := network.LoadFromAST(program)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Submit facts to the network
//	fact := &rete.Fact{
//		ID:   "person_001",
//		Type: "Person",
//		Fields: map[string]interface{}{
//			"name": "Alice",
//			"age":  25,
//		},
//	}
//	network.SubmitFact(fact)
//
// # Integration
//
// This package integrates with the constraint package to automatically build
// RETE networks from constraint program ASTs. It provides seamless conversion
// from high-level constraint rules to efficient pattern matching networks.
//
// # Performance
//
// The implementation is optimized for high-throughput fact processing:
//   - Efficient token propagation
//   - Optimized memory management
//   - Hash-based joins for O(1) lookups
//   - Minimal object allocation during matching
//
// # Storage Backends
//
// The package supports multiple storage backends:
//   - MemoryStorage: In-memory storage for development and testing
//   - Future: Persistent storage backends for production use
package rete
