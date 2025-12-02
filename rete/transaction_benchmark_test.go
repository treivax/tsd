// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// Helper: Build a test network with N facts
func buildNetworkWithFacts(size int) *ReteNetwork {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Add a simple type
	typeDef := TypeDefinition{
		Name: "TestType",
		Fields: []Field{
			{Name: "value", Type: "int"},
			{Name: "name", Type: "string"},
		},
	}
	network.Types = append(network.Types, typeDef)

	// Add facts
	for i := 0; i < size; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("fact_%d", i),
			Type: "TestType",
			Fields: map[string]interface{}{
				"value": i,
				"name":  fmt.Sprintf("test_%d", i),
			},
			Timestamp: time.Now(),
		}
		storage.AddFact(fact)
	}

	return network
}

// Helper: Create a test fact
func createTestFact() *Fact {
	return &Fact{
		ID:   fmt.Sprintf("fact_%d", time.Now().UnixNano()),
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
			"name":  "test",
		},
		Timestamp: time.Now(),
	}
}

// Helper: Estimate network size in bytes
func estimateNetworkSize(network *ReteNetwork) int64 {
	var size int64

	// TypeNodes
	size += int64(len(network.TypeNodes) * 1024)

	// AlphaNodes
	size += int64(len(network.AlphaNodes) * 512)

	// BetaNodes
	size += int64(len(network.BetaNodes) * 2048)

	// TerminalNodes
	size += int64(len(network.TerminalNodes) * 1024)

	// Facts
	if network.Storage != nil {
		size += int64(len(network.Storage.GetAllFacts()) * 200)
	}

	return size
}

// Benchmark: BeginTransaction + Commit with different network sizes
func BenchmarkTransaction_BeginCommit(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("NetworkSize_%d", size), func(b *testing.B) {
			network := buildNetworkWithFacts(size)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tx := network.BeginTransaction()
				// Simulate some operation
				_ = createTestFact()
				tx.Commit()
			}
		})
	}
}

// Benchmark: BeginTransaction only (measures snapshot cost)
func BenchmarkTransaction_BeginOnly(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("NetworkSize_%d", size), func(b *testing.B) {
			network := buildNetworkWithFacts(size)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tx := network.BeginTransaction()
				_ = tx
			}
		})
	}
}

// Benchmark: Rollback with different network sizes
func BenchmarkTransaction_Rollback(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("NetworkSize_%d", size), func(b *testing.B) {
			network := buildNetworkWithFacts(size)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tx := network.BeginTransaction()
				// Simulate operation that will be rolled back
				_ = createTestFact()
				tx.Rollback()
			}
		})
	}
}

// Benchmark: Transaction with multiple operations
func BenchmarkTransaction_MultipleOperations(b *testing.B) {
	opCounts := []int{1, 10, 100}

	for _, opCount := range opCounts {
		b.Run(fmt.Sprintf("Operations_%d", opCount), func(b *testing.B) {
			network := buildNetworkWithFacts(1000)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				tx := network.BeginTransaction()

				// Simulate multiple operations
				for j := 0; j < opCount; j++ {
					_ = createTestFact()
				}

				tx.Commit()
			}
		})
	}
}

// Benchmark: Memory overhead measurement
func BenchmarkTransaction_MemoryOverhead(b *testing.B) {
	sizes := []int{1000, 10000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("NetworkSize_%d", size), func(b *testing.B) {
			network := buildNetworkWithFacts(size)

			// Force GC for clean baseline
			runtime.GC()

			var m1 runtime.MemStats
			runtime.ReadMemStats(&m1)
			baselineAlloc := m1.Alloc

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tx := network.BeginTransaction()

				// Measure memory after transaction created
				runtime.GC()
				var m2 runtime.MemStats
				runtime.ReadMemStats(&m2)

				overhead := m2.Alloc - baselineAlloc

				tx.Commit()

				// Report overhead per transaction
				if i == 0 {
					networkSize := estimateNetworkSize(network)
					b.ReportMetric(float64(overhead)/1024, "KB/tx")
					b.ReportMetric(float64(overhead)/float64(networkSize)*100, "%overhead")
				}
			}
		})
	}
}
