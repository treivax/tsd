// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
)

// This file demonstrates transaction configuration and monitoring patterns.
// TSD uses in-memory storage with strong consistency guarantees.

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   TSD Transactions - Configuration Examples               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	demonstrateConfigurations()
	demonstrateMonitoring()
	demonstrateTuningProcess()

	fmt.Println("\n✅ Examples complete!")
	fmt.Println("\nFor more information:")
	fmt.Println("  📖 User Guide: docs/USER_GUIDE.md")
	fmt.Println("  📊 Architecture: docs/ARCHITECTURE.md")
	fmt.Println("  🔧 API Reference: rete/coherence_mode.go")
}

func demonstrateConfigurations() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("1. Transaction Configuration Patterns")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Default Configuration (optimized for in-memory)
	fmt.Println("🔹 Default Configuration (In-Memory Storage)")
	defaultOpts := rete.DefaultTransactionOptions()
	fmt.Printf("   SubmissionTimeout:  %v\n", defaultOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", defaultOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", defaultOpts.MaxVerifyRetries)
	fmt.Printf("   VerifyOnCommit:     %v\n", defaultOpts.VerifyOnCommit)
	fmt.Printf("   Use case: Single-node in-memory storage\n")
	fmt.Printf("   Performance: ~10,000-50,000 facts/sec\n")
	fmt.Println()

	// Low Latency Configuration
	fmt.Println("🔹 Low Latency Configuration")
	lowLatencyOpts := &rete.TransactionOptions{
		SubmissionTimeout: 5 * time.Second,
		VerifyRetryDelay:  5 * time.Millisecond,
		MaxVerifyRetries:  3,
		VerifyOnCommit:    true,
	}
	fmt.Printf("   SubmissionTimeout:  %v\n", lowLatencyOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", lowLatencyOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", lowLatencyOpts.MaxVerifyRetries)
	fmt.Printf("   Use case: Low-latency requirements, fast in-memory operations\n")
	fmt.Printf("   Performance: ~20,000-50,000 facts/sec\n")
	fmt.Println()

	// Future: Network Replication Configuration
	fmt.Println("🔹 Network Replication Configuration (Future)")
	replicationOpts := &rete.TransactionOptions{
		SubmissionTimeout: 30 * time.Second,
		VerifyRetryDelay:  50 * time.Millisecond,
		MaxVerifyRetries:  10,
		VerifyOnCommit:    true,
	}
	fmt.Printf("   SubmissionTimeout:  %v\n", replicationOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", replicationOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", replicationOpts.MaxVerifyRetries)
	fmt.Printf("   Use case: Multi-node replication via Raft\n")
	fmt.Printf("   Performance: ~1,000-10,000 facts/sec (depends on network)\n")
	fmt.Println()

	// Code examples
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Code Examples")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("// Using default configuration")
	fmt.Println("storage := rete.NewMemoryStorage()")
	fmt.Println("network := rete.NewReteNetwork(storage)")
	fmt.Println("tx := network.BeginTransaction()")
	fmt.Println()

	fmt.Println("// Using custom low-latency configuration")
	fmt.Println("opts := &rete.TransactionOptions{")
	fmt.Println("    SubmissionTimeout: 5 * time.Second,")
	fmt.Println("    VerifyRetryDelay:  5 * time.Millisecond,")
	fmt.Println("    MaxVerifyRetries:  3,")
	fmt.Println("    VerifyOnCommit:    true,")
	fmt.Println("}")
	fmt.Println("tx := network.BeginTransactionWithOptions(opts)")
	fmt.Println()
}

func demonstrateMonitoring() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("2. Performance Monitoring Pattern")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("// Initialize performance metrics collector")
	fmt.Println("perfMetrics := rete.NewStrongModePerformanceMetrics()")
	fmt.Println()

	fmt.Println("// For each transaction:")
	fmt.Println("start := time.Now()")
	fmt.Println("tx := network.BeginTransaction()")
	fmt.Println("// ... add facts ...")
	fmt.Println("err := tx.Commit()")
	fmt.Println("duration := time.Since(start)")
	fmt.Println()

	fmt.Println("// Record metrics")
	fmt.Println("coherenceMetrics := tx.GetCoherenceMetrics()")
	fmt.Println("perfMetrics.RecordTransaction(")
	fmt.Println("    duration,")
	fmt.Println("    factCount,")
	fmt.Println("    err == nil,")
	fmt.Println("    coherenceMetrics,")
	fmt.Println(")")
	fmt.Println()

	fmt.Println("// Generate performance report")
	fmt.Println("fmt.Println(perfMetrics.GetReport())")
	fmt.Println()

	// Simulate some metrics
	perfMetrics := rete.NewStrongModePerformanceMetrics()

	// Simulate 10 successful transactions
	for i := 0; i < 10; i++ {
		perfMetrics.RecordTransaction(
			time.Duration(10+i)*time.Millisecond,
			100,
			true,
			&rete.CoherenceMetrics{
				FactsSubmitted:      100,
				FactsPersisted:      100,
				TotalVerifyAttempts: 1,
				TotalWaitTime:       5 * time.Millisecond,
				TotalSyncTime:       5 * time.Millisecond,
			},
		)
	}

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Sample Performance Report")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println(perfMetrics.GetReport())
	fmt.Println()

	fmt.Println("// Check system health")
	fmt.Println("if !perfMetrics.IsHealthy {")
	fmt.Println("    log.Warn(\"Transaction performance degraded\")")
	fmt.Println("    for _, rec := range perfMetrics.Recommendations {")
	fmt.Println("        log.Info(rec)")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println()
}

func demonstrateTuningProcess() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("3. Performance Tuning Process")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("Step 1: Start with Default Configuration")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("opts := rete.DefaultTransactionOptions()")
	fmt.Println()

	fmt.Println("Step 2: Collect Baseline Metrics (100+ transactions)")
	fmt.Println("────────────────────────────────────────")
	perfMetrics := rete.NewStrongModePerformanceMetrics()

	// Simulate fast storage scenario
	for i := 0; i < 100; i++ {
		perfMetrics.RecordTransaction(
			time.Duration(8+i%5)*time.Millisecond,
			50,
			true,
			&rete.CoherenceMetrics{
				FactsSubmitted:      50,
				FactsPersisted:      50,
				TotalVerifyAttempts: 1,
				TotalWaitTime:       3 * time.Millisecond,
				TotalSyncTime:       4 * time.Millisecond,
			},
		)
	}

	fmt.Println("Baseline collected:")
	fmt.Printf("  Average Transaction Time: %.2fms\n", perfMetrics.AvgTransactionTime.Seconds()*1000)
	fmt.Printf("  Transaction Count: %d\n", perfMetrics.TransactionCount)
	if perfMetrics.TransactionCount > 0 {
		successRate := float64(perfMetrics.SuccessfulTransactions) / float64(perfMetrics.TransactionCount) * 100
		fmt.Printf("  Success Rate: %.1f%%\n", successRate)
	}
	fmt.Println()

	fmt.Println("Step 3: Analyze and Tune")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("// For fast in-memory operations:")
	fmt.Println("opts := &rete.TransactionOptions{")
	fmt.Println("    SubmissionTimeout: 5 * time.Second,  // Reduced from 30s")
	fmt.Println("    VerifyRetryDelay:  5 * time.Millisecond,  // Reduced from 50ms")
	fmt.Println("    MaxVerifyRetries:  3,  // Reduced from 10")
	fmt.Println("    VerifyOnCommit:    true,")
	fmt.Println("}")
	fmt.Println()

	fmt.Println("Step 4: Validate Tuned Configuration")
	fmt.Println("────────────────────────────────────────")
	perfMetrics2 := rete.NewStrongModePerformanceMetrics()

	// Simulate with tuned config (should be faster)
	for i := 0; i < 100; i++ {
		perfMetrics2.RecordTransaction(
			time.Duration(6+i%3)*time.Millisecond,
			50,
			true,
			&rete.CoherenceMetrics{
				FactsSubmitted:      50,
				FactsPersisted:      50,
				TotalVerifyAttempts: 1,
				TotalWaitTime:       2 * time.Millisecond,
				TotalSyncTime:       3 * time.Millisecond,
			},
		)
	}

	fmt.Printf("Improvement: %.1f%% faster\n",
		(1-(perfMetrics2.AvgTransactionTime.Seconds()/perfMetrics.AvgTransactionTime.Seconds()))*100)
	fmt.Println()

	fmt.Println("Step 5: Monitor Continuously")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("// Set up periodic health checks")
	fmt.Println("ticker := time.NewTicker(1 * time.Minute)")
	fmt.Println("go func() {")
	fmt.Println("    for range ticker.C {")
	fmt.Println("        if !perfMetrics.IsHealthy {")
	fmt.Println("            logger.Warn(\"Transaction health degraded\")")
	fmt.Println("            // Review recommendations and adjust")
	fmt.Println("        }")
	fmt.Println("    }")
	fmt.Println("}()")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("Key Tuning Guidelines")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("1. SubmissionTimeout:")
	fmt.Println("   - In-memory: 5-10 seconds")
	fmt.Println("   - Replicated: 20-30 seconds")
	fmt.Println()
	fmt.Println("2. VerifyRetryDelay:")
	fmt.Println("   - In-memory: 5-10ms")
	fmt.Println("   - Replicated: 50-100ms")
	fmt.Println()
	fmt.Println("3. MaxVerifyRetries:")
	fmt.Println("   - In-memory: 3-5 retries")
	fmt.Println("   - Replicated: 10-15 retries")
	fmt.Println()
	fmt.Println("4. VerifyOnCommit:")
	fmt.Println("   - Always true for strong consistency")
	fmt.Println()
}
