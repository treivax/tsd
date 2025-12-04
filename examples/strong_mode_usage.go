// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
)

// This file demonstrates Strong Mode configuration patterns.
// Actual fact insertion uses the Command pattern via RecordAndExecute.

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║   TSD Strong Mode - Configuration Examples                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	demonstrateConfigurations()
	demonstrateMonitoring()
	demonstrateTuningProcess()

	fmt.Println("\n✅ Examples complete!")
	fmt.Println("\nFor more information:")
	fmt.Println("  📖 Tuning Guide: docs/STRONG_MODE_TUNING_GUIDE.md")
	fmt.Println("  📊 Design Doc: docs/PHASE4_COHERENCE_STRONG_MODE.md")
	fmt.Println("  🔧 API Reference: rete/coherence_mode.go")
}

func demonstrateConfigurations() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("1. Configuration Patterns for Different Storage Backends")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Default Configuration (works everywhere)
	fmt.Println("🔹 Default Configuration (Universal)")
	defaultOpts := rete.DefaultTransactionOptions()
	fmt.Printf("   SubmissionTimeout:  %v\n", defaultOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", defaultOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", defaultOpts.MaxVerifyRetries)
	fmt.Printf("   VerifyOnCommit:     %v\n", defaultOpts.VerifyOnCommit)
	fmt.Printf("   Use case: General purpose, safe defaults\n")
	fmt.Printf("   Performance: ~100-1,000 facts/sec\n")
	fmt.Println()

	// PostgreSQL / MySQL (Fast synchronous storage)
	fmt.Println("🔹 PostgreSQL / MySQL Configuration")
	pgOpts := &rete.TransactionOptions{
		SubmissionTimeout: 10 * time.Second,
		VerifyRetryDelay:  10 * time.Millisecond,
		MaxVerifyRetries:  5,
		VerifyOnCommit:    true,
	}
	fmt.Printf("   SubmissionTimeout:  %v\n", pgOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", pgOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", pgOpts.MaxVerifyRetries)
	fmt.Printf("   Use case: Fast, strongly consistent relational databases\n")
	fmt.Printf("   Performance: ~1,000-5,000 facts/sec\n")
	fmt.Println()

	// Redis (Ultra-fast in-memory)
	fmt.Println("🔹 Redis Configuration")
	redisOpts := &rete.TransactionOptions{
		SubmissionTimeout: 5 * time.Second,
		VerifyRetryDelay:  5 * time.Millisecond,
		MaxVerifyRetries:  3,
		VerifyOnCommit:    false, // Optional for synchronous storage
	}
	fmt.Printf("   SubmissionTimeout:  %v\n", redisOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", redisOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", redisOpts.MaxVerifyRetries)
	fmt.Printf("   VerifyOnCommit:     %v\n", redisOpts.VerifyOnCommit)
	fmt.Printf("   Use case: Ultra-fast in-memory storage\n")
	fmt.Printf("   Performance: ~5,000-10,000 facts/sec\n")
	fmt.Println()

	// Cassandra / DynamoDB (Eventually consistent)
	fmt.Println("🔹 Cassandra / DynamoDB Configuration")
	cassandraOpts := &rete.TransactionOptions{
		SubmissionTimeout: 45 * time.Second,
		VerifyRetryDelay:  100 * time.Millisecond,
		MaxVerifyRetries:  12,
		VerifyOnCommit:    true,
	}
	fmt.Printf("   SubmissionTimeout:  %v\n", cassandraOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", cassandraOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", cassandraOpts.MaxVerifyRetries)
	fmt.Printf("   Use case: Eventually consistent distributed databases\n")
	fmt.Printf("   Performance: ~500-2,000 facts/sec\n")
	fmt.Println()

	// High-throughput batch processing
	fmt.Println("🔹 High-Throughput Batch Configuration")
	batchOpts := &rete.TransactionOptions{
		SubmissionTimeout: 5 * time.Second,
		VerifyRetryDelay:  20 * time.Millisecond,
		MaxVerifyRetries:  3,
		VerifyOnCommit:    false, // Trade safety for speed
	}
	fmt.Printf("   SubmissionTimeout:  %v\n", batchOpts.SubmissionTimeout)
	fmt.Printf("   VerifyRetryDelay:   %v\n", batchOpts.VerifyRetryDelay)
	fmt.Printf("   MaxVerifyRetries:   %d\n", batchOpts.MaxVerifyRetries)
	fmt.Printf("   VerifyOnCommit:     %v\n", batchOpts.VerifyOnCommit)
	fmt.Printf("   Use case: Analytics, event streams, non-critical data\n")
	fmt.Printf("   Performance: ~3,000-10,000 facts/sec\n")
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
	fmt.Println("// start := time.Now()")
	fmt.Println("// tx := network.BeginTransaction()")
	fmt.Println("// ... perform operations ...")
	fmt.Println("// err := tx.Commit()")
	fmt.Println("// duration := time.Since(start)")
	fmt.Println()

	fmt.Println("// Record metrics")
	fmt.Println("// coherenceMetrics := collector.Finalize()")
	fmt.Println("// perfMetrics.RecordTransaction(duration, factCount, success, coherenceMetrics)")
	fmt.Println()

	fmt.Println("// Generate performance report")
	fmt.Println("// fmt.Println(perfMetrics.GetReport())")
	fmt.Println()

	// Simulate some metrics
	perfMetrics := rete.NewStrongModePerformanceMetrics()

	// Simulate 10 successful transactions
	for i := 0; i < 10; i++ {
		cm := &rete.CoherenceMetrics{
			FactsSubmitted:      10,
			FactsPersisted:      10,
			TotalRetries:        2,
			TotalVerifyAttempts: 12,
			TotalWaitTime:       100 * time.Millisecond,
		}
		perfMetrics.RecordTransaction(100*time.Millisecond, 10, true, cm)
	}

	fmt.Println("📊 Sample Performance Metrics:")
	fmt.Println(perfMetrics.GetSummary())
	fmt.Println()

	if perfMetrics.IsHealthy {
		fmt.Printf("✅ System Health: %.0f%% (Grade: %s)\n",
			perfMetrics.HealthScore, perfMetrics.PerformanceGrade)
	}

	if len(perfMetrics.Recommendations) > 0 {
		fmt.Println("\n💡 Recommendations:")
		for _, rec := range perfMetrics.Recommendations {
			fmt.Printf("   %s\n", rec)
		}
	}
	fmt.Println()
}

func demonstrateTuningProcess() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("3. Performance Tuning Process")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Println("Step 1: Start with Default Configuration")
	fmt.Println("────────────────────────────────────────")
	opts := rete.DefaultTransactionOptions()
	fmt.Printf("Initial: VerifyRetryDelay=%v, MaxVerifyRetries=%d\n",
		opts.VerifyRetryDelay, opts.MaxVerifyRetries)
	fmt.Println()

	fmt.Println("Step 2: Collect Baseline Metrics (100+ transactions)")
	fmt.Println("────────────────────────────────────────")
	perfMetrics := rete.NewStrongModePerformanceMetrics()

	// Simulate fast storage scenario
	for i := 0; i < 100; i++ {
		cm := &rete.CoherenceMetrics{
			FactsSubmitted:      10,
			FactsPersisted:      10,
			TotalRetries:        1, // Very low retry rate
			TotalVerifyAttempts: 11,
			TotalWaitTime:       50 * time.Millisecond,
		}
		perfMetrics.RecordTransaction(80*time.Millisecond, 10, true, cm)
	}

	fmt.Printf("Baseline Results:\n")
	fmt.Printf("  Avg Retries/Fact: %.2f\n", perfMetrics.AvgRetriesPerFact)
	fmt.Printf("  Health Score:     %.0f%%\n", perfMetrics.HealthScore)
	fmt.Printf("  Timeout Rate:     %.2f%%\n", perfMetrics.TimeoutRate*100)
	fmt.Println()

	fmt.Println("Step 3: Analyze & Tune Configuration")
	fmt.Println("────────────────────────────────────────")

	if perfMetrics.AvgRetriesPerFact < 0.2 && perfMetrics.HealthScore > 90 {
		fmt.Println("✅ Analysis: Low retry rate indicates fast storage")
		fmt.Println("   Recommendation: Optimize for speed")
		fmt.Println()

		// Tune for speed
		oldDelay := opts.VerifyRetryDelay
		opts.VerifyRetryDelay = 20 * time.Millisecond
		opts.MaxVerifyRetries = 5

		perfMetrics.RecordConfigChange("VerifyRetryDelay",
			oldDelay, opts.VerifyRetryDelay,
			"Optimizing for fast storage")

		fmt.Printf("Tuned: VerifyRetryDelay=%v, MaxVerifyRetries=%d\n",
			opts.VerifyRetryDelay, opts.MaxVerifyRetries)
	} else if perfMetrics.AvgRetriesPerFact > 1.0 {
		fmt.Println("⚠️  Analysis: High retry rate indicates slow storage")
		fmt.Println("   Recommendation: Increase delays and retries")
		fmt.Println()

		opts.VerifyRetryDelay = 100 * time.Millisecond
		opts.MaxVerifyRetries = 15

		fmt.Printf("Tuned: VerifyRetryDelay=%v, MaxVerifyRetries=%d\n",
			opts.VerifyRetryDelay, opts.MaxVerifyRetries)
	}
	fmt.Println()

	fmt.Println("Step 4: Validate Tuned Configuration")
	fmt.Println("────────────────────────────────────────")
	perfMetrics2 := rete.NewStrongModePerformanceMetrics()

	// Simulate with tuned config (should be faster)
	for i := 0; i < 100; i++ {
		cm := &rete.CoherenceMetrics{
			FactsSubmitted:      10,
			FactsPersisted:      10,
			TotalRetries:        1,
			TotalVerifyAttempts: 11,
			TotalWaitTime:       40 * time.Millisecond, // Faster with optimized config
		}
		perfMetrics2.RecordTransaction(65*time.Millisecond, 10, true, cm)
	}

	fmt.Printf("After Tuning:\n")
	fmt.Printf("  Avg Transaction Time: %v → %v\n",
		perfMetrics.AvgTransactionTime, perfMetrics2.AvgTransactionTime)

	improvement := float64(perfMetrics.AvgTransactionTime-perfMetrics2.AvgTransactionTime) /
		float64(perfMetrics.AvgTransactionTime) * 100
	fmt.Printf("  Performance Improvement: %.1f%%\n", improvement)
	fmt.Printf("  Health Score: %.0f%% → %.0f%%\n",
		perfMetrics.HealthScore, perfMetrics2.HealthScore)
	fmt.Println()

	fmt.Println("Step 5: Monitor Continuously")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("// Set up periodic health checks")
	fmt.Println("// go func() {")
	fmt.Println("//     ticker := time.NewTicker(5 * time.Minute)")
	fmt.Println("//     for range ticker.C {")
	fmt.Println("//         if !perfMetrics.IsHealthy {")
	fmt.Println("//             logger.Warn(\"Strong mode health degraded\")")
	fmt.Println("//             // Review recommendations and adjust")
	fmt.Println("//         }")
	fmt.Println("//     }")
	fmt.Println("// }()")
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("📈 Key Tuning Guidelines:")
	fmt.Println()
	fmt.Println("1. AvgRetriesPerFact < 0.5:")
	fmt.Println("   ✅ Reduce VerifyRetryDelay and MaxVerifyRetries for speed")
	fmt.Println()
	fmt.Println("2. AvgRetriesPerFact > 1.0:")
	fmt.Println("   ⚠️  Increase VerifyRetryDelay and MaxVerifyRetries")
	fmt.Println()
	fmt.Println("3. TimeoutRate > 5%:")
	fmt.Println("   ⚠️  Increase SubmissionTimeout")
	fmt.Println()
	fmt.Println("4. HealthScore > 95%:")
	fmt.Println("   ✅ Configuration is optimal, consider slight optimizations")
	fmt.Println()
	fmt.Println("5. HealthScore < 80%:")
	fmt.Println("   ❌ Investigate storage performance and tune aggressively")
	fmt.Println()
}
