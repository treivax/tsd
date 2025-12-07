// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License

package rete

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// =============================================================================
// Benchmark: Multi-Source Aggregation Performance
// =============================================================================

// BenchmarkMultiSourceAggregation_TwoSources_SmallScale benchmarks 2-source aggregation with small dataset
func BenchmarkMultiSourceAggregation_TwoSources_SmallScale(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    10,
		NumSourceFacts1: 50,
		NumSourceFacts2: 50,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_TwoSources_MediumScale benchmarks 2-source aggregation with medium dataset
func BenchmarkMultiSourceAggregation_TwoSources_MediumScale(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    100,
		NumSourceFacts1: 500,
		NumSourceFacts2: 500,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_TwoSources_LargeScale benchmarks 2-source aggregation with large dataset
func BenchmarkMultiSourceAggregation_TwoSources_LargeScale(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    1000,
		NumSourceFacts1: 5000,
		NumSourceFacts2: 5000,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_ThreeSources_SmallScale benchmarks 3-source aggregation with small dataset
func BenchmarkMultiSourceAggregation_ThreeSources_SmallScale(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    10,
		NumSourceFacts1: 50,
		NumSourceFacts2: 50,
		NumSourceFacts3: 50,
		NumAggVars:      3,
		NumSources:      3,
	})
}

// BenchmarkMultiSourceAggregation_ThreeSources_MediumScale benchmarks 3-source aggregation with medium dataset
func BenchmarkMultiSourceAggregation_ThreeSources_MediumScale(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    100,
		NumSourceFacts1: 500,
		NumSourceFacts2: 500,
		NumSourceFacts3: 500,
		NumAggVars:      3,
		NumSources:      3,
	})
}

// BenchmarkMultiSourceAggregation_ThreeSources_LargeScale benchmarks 3-source aggregation with large dataset
func BenchmarkMultiSourceAggregation_ThreeSources_LargeScale(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    500,
		NumSourceFacts1: 2500,
		NumSourceFacts2: 2500,
		NumSourceFacts3: 2500,
		NumAggVars:      3,
		NumSources:      3,
	})
}

// BenchmarkMultiSourceAggregation_ManyAggregates benchmarks with many aggregation functions
func BenchmarkMultiSourceAggregation_ManyAggregates(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    100,
		NumSourceFacts1: 500,
		NumSourceFacts2: 500,
		NumAggVars:      5, // AVG, SUM, COUNT, MIN, MAX
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_HighFanout benchmarks high fanout scenario (many source facts per main fact)
func BenchmarkMultiSourceAggregation_HighFanout(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    10,
		NumSourceFacts1: 1000, // 100 source facts per main fact
		NumSourceFacts2: 1000,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_LowFanout benchmarks low fanout scenario (few source facts per main fact)
func BenchmarkMultiSourceAggregation_LowFanout(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    100,
		NumSourceFacts1: 150, // 1.5 source facts per main fact
		NumSourceFacts2: 150,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_WithThresholds benchmarks aggregation with threshold evaluation
func BenchmarkMultiSourceAggregation_WithThresholds(b *testing.B) {
	benchmarkMultiSourceAggregationWithThresholds(b, BenchmarkConfig{
		NumMainFacts:    100,
		NumSourceFacts1: 500,
		NumSourceFacts2: 500,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_Retraction benchmarks fact retraction and recomputation
func BenchmarkMultiSourceAggregation_Retraction(b *testing.B) {
	benchmarkMultiSourceAggregationRetraction(b, BenchmarkConfig{
		NumMainFacts:    100,
		NumSourceFacts1: 500,
		NumSourceFacts2: 500,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_IncrementalUpdate benchmarks incremental fact addition
func BenchmarkMultiSourceAggregation_IncrementalUpdate(b *testing.B) {
	benchmarkMultiSourceAggregationIncremental(b, BenchmarkConfig{
		NumMainFacts:    50,
		NumSourceFacts1: 250,
		NumSourceFacts2: 250,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_Memory_SmallScale measures memory usage for small scale
func BenchmarkMultiSourceAggregation_Memory_SmallScale(b *testing.B) {
	benchmarkMultiSourceAggregationMemory(b, BenchmarkConfig{
		NumMainFacts:    10,
		NumSourceFacts1: 50,
		NumSourceFacts2: 50,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_Memory_LargeScale measures memory usage for large scale
func BenchmarkMultiSourceAggregation_Memory_LargeScale(b *testing.B) {
	benchmarkMultiSourceAggregationMemory(b, BenchmarkConfig{
		NumMainFacts:    1000,
		NumSourceFacts1: 5000,
		NumSourceFacts2: 5000,
		NumAggVars:      2,
		NumSources:      2,
	})
}

// =============================================================================
// Benchmark Configuration
// =============================================================================

type BenchmarkConfig struct {
	NumMainFacts    int
	NumSourceFacts1 int
	NumSourceFacts2 int
	NumSourceFacts3 int
	NumAggVars      int
	NumSources      int
}

// =============================================================================
// Core Benchmark Functions
// =============================================================================

func benchmarkMultiSourceAggregation(b *testing.B, config BenchmarkConfig) {
	// Setup
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "benchmark.tsd")

	content := generateTSDContent(config, false)
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()

	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}

	// Pre-generate facts
	mainFacts := generateMainFacts(config.NumMainFacts)
	sourceFacts1 := generateSourceFacts("Employee", "emp", "deptId", config.NumMainFacts, config.NumSourceFacts1)
	sourceFacts2 := generateSourceFacts("Performance", "perf", "employeeId", config.NumSourceFacts1, config.NumSourceFacts2)
	var sourceFacts3 []*Fact
	if config.NumSources >= 3 {
		sourceFacts3 = generateSourceFacts("Training", "train", "employeeId", config.NumSourceFacts1, config.NumSourceFacts3)
	}

	b.ResetTimer()
	b.ReportAllocs()

	startTime := time.Now()

	for i := 0; i < b.N; i++ {
		// Submit facts
		for _, fact := range mainFacts {
			network.SubmitFact(fact)
		}
		for _, fact := range sourceFacts1 {
			network.SubmitFact(fact)
		}
		for _, fact := range sourceFacts2 {
			network.SubmitFact(fact)
		}
		if config.NumSources >= 3 {
			for _, fact := range sourceFacts3 {
				network.SubmitFact(fact)
			}
		}

		// Reset for next iteration (if multiple iterations)
		if i < b.N-1 {
			network.Reset()
		}
	}

	b.StopTimer()

	// Calculate throughput
	elapsed := time.Since(startTime)
	totalFacts := config.NumMainFacts + config.NumSourceFacts1 + config.NumSourceFacts2
	if config.NumSources >= 3 {
		totalFacts += config.NumSourceFacts3
	}
	factsPerSec := float64(totalFacts*b.N) / elapsed.Seconds()

	// Report metrics
	b.ReportMetric(factsPerSec, "facts/sec")
	b.ReportMetric(float64(totalFacts), "total_facts")
	b.ReportMetric(float64(config.NumMainFacts), "main_facts")
	b.ReportMetric(float64(config.NumAggVars), "agg_vars")
	b.ReportMetric(elapsed.Seconds()/float64(b.N), "sec/op")

	// Count activations
	activationCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activationCount += len(terminalNode.GetMemory().Tokens)
	}
	b.ReportMetric(float64(activationCount), "activations")
}

func benchmarkMultiSourceAggregationWithThresholds(b *testing.B, config BenchmarkConfig) {
	// Setup
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "benchmark_threshold.tsd")

	content := generateTSDContent(config, true)
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()

	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}

	// Pre-generate facts
	mainFacts := generateMainFacts(config.NumMainFacts)
	sourceFacts1 := generateSourceFacts("Employee", "emp", "deptId", config.NumMainFacts, config.NumSourceFacts1)
	sourceFacts2 := generateSourceFacts("Performance", "perf", "employeeId", config.NumSourceFacts1, config.NumSourceFacts2)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Submit facts
		for _, fact := range mainFacts {
			network.SubmitFact(fact)
		}
		for _, fact := range sourceFacts1 {
			network.SubmitFact(fact)
		}
		for _, fact := range sourceFacts2 {
			network.SubmitFact(fact)
		}

		// Reset for next iteration
		if i < b.N-1 {
			network.Reset()
		}
	}

	b.StopTimer()

	// Count activations (should be fewer due to thresholds)
	activationCount := 0
	for _, terminalNode := range network.TerminalNodes {
		activationCount += len(terminalNode.GetMemory().Tokens)
	}
	b.ReportMetric(float64(activationCount), "activations")
}

func benchmarkMultiSourceAggregationRetraction(b *testing.B, config BenchmarkConfig) {
	// Setup
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "benchmark_retraction.tsd")

	content := generateTSDContent(config, false)
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()

	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}

	// Pre-generate facts
	mainFacts := generateMainFacts(config.NumMainFacts)
	sourceFacts1 := generateSourceFacts("Employee", "emp", "deptId", config.NumMainFacts, config.NumSourceFacts1)
	sourceFacts2 := generateSourceFacts("Performance", "perf", "employeeId", config.NumSourceFacts1, config.NumSourceFacts2)

	// Submit all facts initially
	for _, fact := range mainFacts {
		network.SubmitFact(fact)
	}
	for _, fact := range sourceFacts1 {
		network.SubmitFact(fact)
	}
	for _, fact := range sourceFacts2 {
		network.SubmitFact(fact)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark retraction performance
	for i := 0; i < b.N; i++ {
		// Retract 10% of source facts
		retractionCount := config.NumSourceFacts1 / 10
		for j := 0; j < retractionCount; j++ {
			idx := (i*retractionCount + j) % config.NumSourceFacts1
			network.RetractFact(sourceFacts1[idx].ID)
		}

		// Re-add them for next iteration
		if i < b.N-1 {
			for j := 0; j < retractionCount; j++ {
				idx := (i*retractionCount + j) % config.NumSourceFacts1
				network.SubmitFact(sourceFacts1[idx])
			}
		}
	}

	b.StopTimer()
	b.ReportMetric(float64(config.NumSourceFacts1/10), "retractions/op")
}

func benchmarkMultiSourceAggregationIncremental(b *testing.B, config BenchmarkConfig) {
	// Setup
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "benchmark_incremental.tsd")

	content := generateTSDContent(config, false)
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()

	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}

	// Pre-generate facts
	mainFacts := generateMainFacts(config.NumMainFacts)
	sourceFacts1 := generateSourceFacts("Employee", "emp", "deptId", config.NumMainFacts, config.NumSourceFacts1)
	sourceFacts2 := generateSourceFacts("Performance", "perf", "employeeId", config.NumSourceFacts1, config.NumSourceFacts2)

	// Submit main facts once
	for _, fact := range mainFacts {
		network.SubmitFact(fact)
	}

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark incremental fact addition
	for i := 0; i < b.N; i++ {
		// Add one source fact from each source
		idx := i % len(sourceFacts1)
		network.SubmitFact(sourceFacts1[idx])
		network.SubmitFact(sourceFacts2[idx])

		// Clean up for next iteration
		if i < b.N-1 && i%100 == 99 {
			network.Reset()
			for _, fact := range mainFacts {
				network.SubmitFact(fact)
			}
		}
	}

	b.StopTimer()
	b.ReportMetric(2.0, "facts_added/op")
}

func benchmarkMultiSourceAggregationMemory(b *testing.B, config BenchmarkConfig) {
	// Setup
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "benchmark_memory.tsd")

	content := generateTSDContent(config, false)
	err := os.WriteFile(tsdFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()

	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}

	// Pre-generate facts
	mainFacts := generateMainFacts(config.NumMainFacts)
	sourceFacts1 := generateSourceFacts("Employee", "emp", "deptId", config.NumMainFacts, config.NumSourceFacts1)
	sourceFacts2 := generateSourceFacts("Performance", "perf", "employeeId", config.NumSourceFacts1, config.NumSourceFacts2)

	// Force GC before measurement
	runtime.GC()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Submit facts
		for _, fact := range mainFacts {
			network.SubmitFact(fact)
		}
		for _, fact := range sourceFacts1 {
			network.SubmitFact(fact)
		}
		for _, fact := range sourceFacts2 {
			network.SubmitFact(fact)
		}

		// Reset for next iteration
		if i < b.N-1 {
			network.Reset()
		}
	}

	b.StopTimer()

	// Measure memory after
	runtime.GC()
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Report memory metrics
	allocatedMB := float64(memAfter.Alloc-memBefore.Alloc) / 1024 / 1024
	b.ReportMetric(allocatedMB, "MB_allocated")
	b.ReportMetric(float64(memAfter.Mallocs-memBefore.Mallocs), "mallocs")
	b.ReportMetric(float64(memAfter.TotalAlloc-memBefore.TotalAlloc)/1024/1024, "MB_total_alloc")
}

// =============================================================================
// Helper Functions
// =============================================================================

func generateTSDContent(config BenchmarkConfig, withThreshold bool) string {
	content := `type Department(id: string, name:string)
type Employee(id: string, deptId: string, salary:number)
type Performance(id: string, employeeId: string, score:number)
`

	if config.NumSources >= 3 {
		content += `type Training(id: string, employeeId: string, hours:number)
`
	}

	content += "\nrule benchmark_rule : {d: Department"

	// Add aggregation variables
	if config.NumAggVars >= 1 {
		content += ", avg_sal: AVG(e.salary)"
	}
	if config.NumAggVars >= 2 {
		content += ", avg_score: AVG(p.score)"
	}
	if config.NumAggVars >= 3 {
		content += ", total_sal: SUM(e.salary)"
	}
	if config.NumAggVars >= 4 {
		content += ", emp_count: COUNT(e.id)"
	}
	if config.NumAggVars >= 5 {
		content += ", max_score: MAX(p.score)"
	}

	content += "} / {e: Employee} / {p: Performance}"

	if config.NumSources >= 3 {
		content += " / {t: Training}"
	}

	content += " / e.deptId == d.id AND p.employeeId == e.id"

	if config.NumSources >= 3 {
		content += " AND t.employeeId == e.id"
	}

	// Add thresholds if requested
	if withThreshold {
		content += " AND avg_sal > 50000 AND avg_score > 75"
	}

	content += " ==> print(\"Benchmark rule fired\")\n"

	return content
}

func generateMainFacts(count int) []*Fact {
	facts := make([]*Fact, count)
	for i := 0; i < count; i++ {
		facts[i] = &Fact{
			ID:   fmt.Sprintf("dept%d", i),
			Type: "Department",
			Fields: map[string]interface{}{
				"id":   fmt.Sprintf("dept%d", i),
				"name": fmt.Sprintf("Department %d", i),
			},
		}
	}
	return facts
}

func generateSourceFacts(factType, prefix, foreignKey string, numParents, count int) []*Fact {
	facts := make([]*Fact, count)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		parentIdx := i % numParents
		foreignKeyValue := ""

		if factType == "Employee" {
			foreignKeyValue = fmt.Sprintf("dept%d", parentIdx)
		} else {
			foreignKeyValue = fmt.Sprintf("emp%d", parentIdx)
		}

		fields := map[string]interface{}{
			"id":       fmt.Sprintf("%s%d", prefix, i),
			foreignKey: foreignKeyValue,
		}

		// Add type-specific fields
		switch factType {
		case "Employee":
			fields["salary"] = 40000 + rand.Intn(80000) // 40k-120k
		case "Performance":
			fields["score"] = 60 + rand.Intn(40) // 60-100
		case "Training":
			fields["hours"] = 10 + rand.Intn(90) // 10-100
		}

		facts[i] = &Fact{
			ID:     fmt.Sprintf("%s%d", prefix, i),
			Type:   factType,
			Fields: fields,
		}
	}
	return facts
}

// =============================================================================
// Profiling Helpers
// =============================================================================

// BenchmarkMultiSourceAggregation_Profile runs a profiling-friendly benchmark
// Usage: go test -bench=BenchmarkMultiSourceAggregation_Profile -cpuprofile=cpu.prof -memprofile=mem.prof
func BenchmarkMultiSourceAggregation_Profile(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    500,
		NumSourceFacts1: 2500,
		NumSourceFacts2: 2500,
		NumAggVars:      3,
		NumSources:      2,
	})
}

// BenchmarkMultiSourceAggregation_ThreeSources_Profile runs a 3-source profiling benchmark
func BenchmarkMultiSourceAggregation_ThreeSources_Profile(b *testing.B) {
	benchmarkMultiSourceAggregation(b, BenchmarkConfig{
		NumMainFacts:    500,
		NumSourceFacts1: 2500,
		NumSourceFacts2: 2500,
		NumSourceFacts3: 2500,
		NumAggVars:      3,
		NumSources:      3,
	})
}
