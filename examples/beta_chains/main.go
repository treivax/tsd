// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/treivax/tsd/rete"
)

// Configuration flags
var (
	scenario      = flag.String("scenario", "", "Scenario to run: simple, complex, advanced")
	compare       = flag.Bool("compare", false, "Compare all scenarios")
	export        = flag.String("export", "", "Export results to file (json or csv)")
	configType    = flag.String("config", "default", "Config type: default, high-performance, memory-optimized")
	hashCacheSize = flag.Int("hash-cache", 0, "Hash cache size (0 = use config default)")
	noMetrics     = flag.Bool("no-metrics", false, "Disable metrics collection")
)

// Results stores benchmark results
type Results struct {
	Scenario           string                 `json:"scenario"`
	Config             map[string]interface{} `json:"config"`
	Metrics            MetricsData            `json:"metrics"`
	Comparison         *ComparisonData        `json:"comparison,omitempty"`
	ChainVisualization string                 `json:"chain_visualization,omitempty"`
}

type MetricsData struct {
	NodesCreated    int     `json:"nodes_created"`
	NodesReused     int     `json:"nodes_reused"`
	SharingRatio    float64 `json:"sharing_ratio"`
	BuildTimeUs     int64   `json:"build_time_us"`
	MemorySavedKB   int     `json:"memory_saved_kb"`
	CacheEfficiency float64 `json:"cache_efficiency"`
	TotalRules      int     `json:"total_rules"`
	AvgBuildTimeUs  int64   `json:"avg_build_time_us"`
}

type ComparisonData struct {
	WithoutSharing WithoutSharingMetrics `json:"without_sharing"`
	Gains          GainsData             `json:"gains"`
}

type WithoutSharingMetrics struct {
	NodesCreated int   `json:"nodes_created"`
	BuildTimeUs  int64 `json:"build_time_us"`
}

type GainsData struct {
	NodesSaved     int     `json:"nodes_saved"`
	TimeSavedPct   float64 `json:"time_saved_pct"`
	MemorySavedPct float64 `json:"memory_saved_pct"`
}

func main() {
	flag.Parse()

	printHeader()

	if *compare {
		runComparison()
		return
	}

	if *scenario == "" {
		runInteractive()
	} else {
		runScenario(*scenario)
	}
}

func printHeader() {
	fmt.Println("🚀 Beta Chains Interactive Example")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
}

func runInteractive() {
	fmt.Println("Select a scenario:")
	fmt.Println("1. Simple (2-3 joins, high sharing)")
	fmt.Println("2. Complex (5+ joins, mixed sharing)")
	fmt.Println("3. Advanced (real-world monitoring)")
	fmt.Println("4. Compare all scenarios")
	fmt.Println()
	fmt.Print("Choice [1-4]: ")

	var choice int
	fmt.Scanf("%d", &choice)
	fmt.Println()

	switch choice {
	case 1:
		runScenario("simple")
	case 2:
		runScenario("complex")
	case 3:
		runScenario("advanced")
	case 4:
		runComparison()
	default:
		fmt.Println("❌ Invalid choice")
		os.Exit(1)
	}
}

func runScenario(name string) {
	fmt.Printf("🏃 Running scenario: %s\n", name)
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println()

	// Get configuration
	config := getConfig()

	// Print configuration
	printConfig(config, name)

	// Run with Beta Sharing (now always enabled)
	fmt.Println("🏗️  Building network with Beta Sharing (always enabled)...")
	resultsWith := runWithConfig(name, config)

	// Display results
	fmt.Println()
	displayResults(resultsWith, nil)

	// Export if requested
	if *export != "" {
		exportResults(resultsWith, nil)
	}
}

func getConfig() *rete.ChainPerformanceConfig {
	var config *rete.ChainPerformanceConfig

	// Apply preset configurations
	switch *configType {
	case "high-performance":
		config = rete.HighPerformanceConfig()
	case "memory-optimized":
		config = rete.LowMemoryConfig()
	default:
		config = rete.DefaultChainPerformanceConfig()
	}

	// Override with CLI flags
	if *hashCacheSize > 0 {
		config.BetaHashCacheMaxSize = *hashCacheSize
	}
	if *noMetrics {
		config.MetricsEnabled = false
	} else {
		config.MetricsEnabled = true
	}

	// Beta Sharing is now always enabled (no longer configurable)

	return config
}

func printConfig(config *rete.ChainPerformanceConfig, scenario string) {
	fmt.Println("📊 Configuration:")
	fmt.Printf("  Scenario:         %s\n", scenario)
	fmt.Println("  Beta Sharing:     ✅ Always Enabled")
	fmt.Printf("  Hash Cache Size:  %d\n", config.BetaHashCacheMaxSize)
	fmt.Printf("  Join Cache Size:  %d\n", config.BetaJoinResultCacheMaxSize)
	if config.MetricsEnabled {
		fmt.Println("  Metrics:          ✅ Enabled")
	} else {
		fmt.Println("  Metrics:          ❌ Disabled")
	}
	fmt.Println()
}

func runWithConfig(scenario string, config *rete.ChainPerformanceConfig) *Results {
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetworkWithConfig(storage, config)

	// Capture memory before
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// Build network
	startTime := time.Now()

	var ruleCount int
	switch scenario {
	case "simple":
		ruleCount = runSimpleScenario(network)
	case "complex":
		ruleCount = runComplexScenario(network)
	case "advanced":
		ruleCount = runAdvancedScenario(network)
	default:
		fmt.Printf("❌ Unknown scenario: %s\n", scenario)
		os.Exit(1)
	}

	buildTime := time.Since(startTime)

	// Capture memory after
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)
	memoryUsedKB := int((memAfter.Alloc - memBefore.Alloc) / 1024)

	// Get metrics
	var metricsData MetricsData
	if config.MetricsEnabled {
		metrics := network.GetBetaChainMetrics()
		snapshot := metrics.GetSnapshot()

		metricsData = MetricsData{
			NodesCreated:    snapshot.TotalNodesCreated,
			NodesReused:     snapshot.TotalNodesReused,
			SharingRatio:    snapshot.SharingRatio,
			BuildTimeUs:     buildTime.Microseconds(),
			MemorySavedKB:   memoryUsedKB,
			CacheEfficiency: metrics.GetJoinCacheEfficiency(),
			TotalRules:      ruleCount,
			AvgBuildTimeUs:  buildTime.Microseconds() / int64(ruleCount),
		}
	} else {
		metricsData = MetricsData{
			BuildTimeUs:    buildTime.Microseconds(),
			MemorySavedKB:  memoryUsedKB,
			TotalRules:     ruleCount,
			AvgBuildTimeUs: buildTime.Microseconds() / int64(ruleCount),
		}
	}

	results := &Results{
		Scenario: scenario,
		Config: map[string]interface{}{
			"beta_sharing":    "always_enabled",
			"join_cache_size": config.BetaJoinResultCacheMaxSize,
			"hash_cache_size": config.BetaHashCacheMaxSize,
			"enable_metrics":  config.MetricsEnabled,
		},
		Metrics: metricsData,
	}

	return results
}

func displayResults(with *Results, without *Results) {
	fmt.Println("📈 Beta Sharing Metrics:")
	fmt.Println(strings.Repeat("-", 60))

	// Node metrics
	if with.Metrics.NodesCreated > 0 || with.Metrics.NodesReused > 0 {
		fmt.Println()
		fmt.Println("┌─────────────────────────────────────────────┐")
		fmt.Println("│ Nodes                                       │")
		fmt.Println("├─────────────────────────────────────────────┤")
		fmt.Printf("│ Created:              %-22d│\n", with.Metrics.NodesCreated)
		fmt.Printf("│ Reused:               %-22d│\n", with.Metrics.NodesReused)
		fmt.Printf("│ Sharing Ratio:        %-20.1f%% │\n", with.Metrics.SharingRatio*100)
		fmt.Println("└─────────────────────────────────────────────┘")
	}

	// Performance metrics
	fmt.Println()
	fmt.Println("┌─────────────────────────────────────────────┐")
	fmt.Println("│ Performance                                 │")
	fmt.Println("├─────────────────────────────────────────────┤")
	fmt.Printf("│ Build Time:           %-20dµs │\n", with.Metrics.BuildTimeUs)
	fmt.Printf("│ Avg Time/Rule:        %-20dµs │\n", with.Metrics.AvgBuildTimeUs)
	fmt.Printf("│ Total Rules:          %-22d│\n", with.Metrics.TotalRules)
	fmt.Printf("│ Memory Used:          ~%-20dKB │\n", with.Metrics.MemorySavedKB)
	fmt.Println("└─────────────────────────────────────────────┘")

	// Cache metrics
	if with.Metrics.CacheEfficiency > 0 {
		fmt.Println()
		fmt.Println("┌─────────────────────────────────────────────┐")
		fmt.Println("│ Cache                                       │")
		fmt.Println("├─────────────────────────────────────────────┤")
		fmt.Printf("│ Cache Efficiency:     %-20.1f%% │\n", with.Metrics.CacheEfficiency*100)
		fmt.Println("└─────────────────────────────────────────────┘")
	}

	// Comparison if available
	if without != nil {
		fmt.Println()
		fmt.Println("📊 Comparison: WITH vs WITHOUT Beta Sharing")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println()
		fmt.Printf("%-20s %15s %15s %10s\n", "", "WITH Sharing", "WITHOUT", "GAIN")
		fmt.Println(strings.Repeat("─", 60))

		nodeGain := 0.0
		if without.Metrics.NodesCreated > 0 {
			nodeGain = float64(without.Metrics.NodesCreated-with.Metrics.NodesCreated) / float64(without.Metrics.NodesCreated) * 100
		}
		fmt.Printf("%-20s %15d %15d %9.0f%%\n", "Nodes Created",
			with.Metrics.NodesCreated, without.Metrics.NodesCreated, nodeGain)

		timeGain := 0.0
		if without.Metrics.BuildTimeUs > 0 {
			timeGain = float64(without.Metrics.BuildTimeUs-with.Metrics.BuildTimeUs) / float64(without.Metrics.BuildTimeUs) * 100
		}
		fmt.Printf("%-20s %13dµs %13dµs %9.0f%%\n", "Build Time",
			with.Metrics.BuildTimeUs, without.Metrics.BuildTimeUs, timeGain)

		memGain := 0.0
		if without.Metrics.MemorySavedKB > 0 {
			memGain = float64(without.Metrics.MemorySavedKB-with.Metrics.MemorySavedKB) / float64(without.Metrics.MemorySavedKB) * 100
		}
		fmt.Printf("%-20s %14dKB %14dKB %9.0f%%\n", "Memory Used",
			with.Metrics.MemorySavedKB, without.Metrics.MemorySavedKB, memGain)

		fmt.Println()
		fmt.Println("💰 Total Savings:")
		if without.Metrics.NodesCreated > with.Metrics.NodesCreated {
			fmt.Printf("  - %d JoinNodes saved\n", without.Metrics.NodesCreated-with.Metrics.NodesCreated)
		}
		if without.Metrics.BuildTimeUs > with.Metrics.BuildTimeUs {
			fmt.Printf("  - %dµs faster (%.0f%% improvement)\n",
				without.Metrics.BuildTimeUs-with.Metrics.BuildTimeUs, timeGain)
		}
		if without.Metrics.MemorySavedKB > with.Metrics.MemorySavedKB {
			fmt.Printf("  - %dKB memory saved (%.0f%% reduction)\n",
				without.Metrics.MemorySavedKB-with.Metrics.MemorySavedKB, memGain)
		}
	}

	fmt.Println()
}

func runComparison() {
	scenarios := []string{"simple", "complex", "advanced"}
	config := getConfig()

	fmt.Println("📊 Running Beta Sharing across all scenarios")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	for _, scenario := range scenarios {
		fmt.Printf("📍 Scenario: %s\n", scenario)
		fmt.Println(strings.Repeat("-", 60))

		// Beta Sharing is always enabled
		results := runWithConfig(scenario, config)

		// Display results
		displayResults(results, nil)
		fmt.Println()
	}
}

func exportResults(with *Results, without *Results) {
	// Add comparison if available
	if without != nil {
		nodesSaved := without.Metrics.NodesCreated - with.Metrics.NodesCreated
		timeSavedPct := 0.0
		if without.Metrics.BuildTimeUs > 0 {
			timeSavedPct = float64(without.Metrics.BuildTimeUs-with.Metrics.BuildTimeUs) / float64(without.Metrics.BuildTimeUs) * 100
		}
		memSavedPct := 0.0
		if without.Metrics.MemorySavedKB > 0 {
			memSavedPct = float64(without.Metrics.MemorySavedKB-with.Metrics.MemorySavedKB) / float64(without.Metrics.MemorySavedKB) * 100
		}

		with.Comparison = &ComparisonData{
			WithoutSharing: WithoutSharingMetrics{
				NodesCreated: without.Metrics.NodesCreated,
				BuildTimeUs:  without.Metrics.BuildTimeUs,
			},
			Gains: GainsData{
				NodesSaved:     nodesSaved,
				TimeSavedPct:   timeSavedPct,
				MemorySavedPct: memSavedPct,
			},
		}
	}

	// Export based on file extension
	if strings.HasSuffix(*export, ".json") {
		exportJSON(with)
	} else if strings.HasSuffix(*export, ".csv") {
		exportCSV(with)
	} else {
		fmt.Printf("❌ Unknown export format: %s (use .json or .csv)\n", *export)
		os.Exit(1)
	}

	fmt.Printf("✅ Results exported to: %s\n", *export)
}

func exportJSON(results *Results) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Printf("❌ Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*export, data, 0644)
	if err != nil {
		fmt.Printf("❌ Error writing file: %v\n", err)
		os.Exit(1)
	}
}

func exportCSV(results *Results) {
	file, err := os.Create(*export)
	if err != nil {
		fmt.Printf("❌ Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	header := []string{"Metric", "Value"}
	writer.Write(header)

	// Data rows
	rows := [][]string{
		{"Scenario", results.Scenario},
		{"Nodes Created", fmt.Sprintf("%d", results.Metrics.NodesCreated)},
		{"Nodes Reused", fmt.Sprintf("%d", results.Metrics.NodesReused)},
		{"Sharing Ratio", fmt.Sprintf("%.2f%%", results.Metrics.SharingRatio*100)},
		{"Build Time (µs)", fmt.Sprintf("%d", results.Metrics.BuildTimeUs)},
		{"Memory Used (KB)", fmt.Sprintf("%d", results.Metrics.MemorySavedKB)},
		{"Cache Efficiency", fmt.Sprintf("%.2f%%", results.Metrics.CacheEfficiency*100)},
		{"Total Rules", fmt.Sprintf("%d", results.Metrics.TotalRules)},
	}

	if results.Comparison != nil {
		rows = append(rows, [][]string{
			{"Nodes Saved", fmt.Sprintf("%d", results.Comparison.Gains.NodesSaved)},
			{"Time Saved", fmt.Sprintf("%.1f%%", results.Comparison.Gains.TimeSavedPct)},
			{"Memory Saved", fmt.Sprintf("%.1f%%", results.Comparison.Gains.MemorySavedPct)},
		}...)
	}

	for _, row := range rows {
		writer.Write(row)
	}
}

// Placeholder scenario functions - to be implemented in separate files
func runSimpleScenario(network *rete.ReteNetwork) int {
	fmt.Println("  ℹ️  Simple scenario: Creating 5 rules with shared joins...")
	fmt.Println("  ℹ️  (Placeholder - actual implementation in scenarios/simple.go)")
	return 5
}

func runComplexScenario(network *rete.ReteNetwork) int {
	fmt.Println("  ℹ️  Complex scenario: Creating 10 rules with cascading joins...")
	fmt.Println("  ℹ️  (Placeholder - actual implementation in scenarios/complex.go)")
	return 10
}

func runAdvancedScenario(network *rete.ReteNetwork) int {
	fmt.Println("  ℹ️  Advanced scenario: Creating 20 rules for monitoring...")
	fmt.Println("  ℹ️  (Placeholder - actual implementation in scenarios/advanced.go)")
	return 20
}
