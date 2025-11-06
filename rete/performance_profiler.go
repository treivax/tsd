package rete

import (
	"fmt"
	"runtime"
	"time"
)

// PerformanceReport représente un rapport de performances détaillé
type PerformanceReport struct {
	ComponentName    string                 `json:"component_name"`
	TestName        string                 `json:"test_name"`
	Duration        time.Duration          `json:"duration"`
	OperationsCount int64                  `json:"operations_count"`
	ThroughputOps   float64               `json:"throughput_ops_per_sec"`
	MemoryUsage     MemoryStats           `json:"memory_usage"`
	CacheStats      map[string]interface{} `json:"cache_stats,omitempty"`
	ErrorCount      int                   `json:"error_count"`
	Timestamp       time.Time             `json:"timestamp"`
}

// MemoryStats représente les statistiques mémoire
type MemoryStats struct {
	AllocBytes      uint64 `json:"alloc_bytes"`
	TotalAllocBytes uint64 `json:"total_alloc_bytes"`
	SysBytes        uint64 `json:"sys_bytes"`
	NumGC           uint32 `json:"num_gc"`
	HeapObjects     uint64 `json:"heap_objects"`
}

// PerformanceProfiler gère le profilage des performances
type PerformanceProfiler struct {
	startTime   time.Time
	startMemory runtime.MemStats
	reports     []PerformanceReport
}

// NewPerformanceProfiler crée un nouveau profileur de performances
func NewPerformanceProfiler() *PerformanceProfiler {
	return &PerformanceProfiler{
		reports: make([]PerformanceReport, 0),
	}
}

// StartProfiling démarre le profilage d'un composant
func (p *PerformanceProfiler) StartProfiling(componentName, testName string) {
	p.startTime = time.Now()
	runtime.GC() // Force GC pour des mesures précises
	runtime.ReadMemStats(&p.startMemory)
}

// EndProfiling termine le profilage et génère un rapport
func (p *PerformanceProfiler) EndProfiling(componentName, testName string, operationsCount int64) PerformanceReport {
	endTime := time.Now()
	var endMemory runtime.MemStats
	runtime.ReadMemStats(&endMemory)
	
	duration := endTime.Sub(p.startTime)
	throughput := float64(operationsCount) / duration.Seconds()
	
	report := PerformanceReport{
		ComponentName:   componentName,
		TestName:       testName,
		Duration:       duration,
		OperationsCount: operationsCount,
		ThroughputOps:  throughput,
		MemoryUsage: MemoryStats{
			AllocBytes:      endMemory.Alloc - p.startMemory.Alloc,
			TotalAllocBytes: endMemory.TotalAlloc - p.startMemory.TotalAlloc,
			SysBytes:        endMemory.Sys,
			NumGC:          endMemory.NumGC - p.startMemory.NumGC,
			HeapObjects:    endMemory.HeapObjects,
		},
		Timestamp: endTime,
	}
	
	p.reports = append(p.reports, report)
	return report
}

// AddCacheStats ajoute des statistiques de cache au dernier rapport
func (p *PerformanceProfiler) AddCacheStats(stats map[string]interface{}) {
	if len(p.reports) > 0 {
		p.reports[len(p.reports)-1].CacheStats = stats
	}
}

// GetReports retourne tous les rapports générés
func (p *PerformanceProfiler) GetReports() []PerformanceReport {
	return p.reports
}

// PrintSummary affiche un résumé des performances
func (p *PerformanceProfiler) PrintSummary() {
	fmt.Println("\n=== RAPPORT DE PERFORMANCES ===")
	fmt.Printf("Total tests exécutés: %d\n\n", len(p.reports))
	
	for _, report := range p.reports {
		fmt.Printf("Component: %s | Test: %s\n", report.ComponentName, report.TestName)
		fmt.Printf("  Durée: %v\n", report.Duration)
		fmt.Printf("  Opérations: %d\n", report.OperationsCount)
		fmt.Printf("  Débit: %.2f ops/sec\n", report.ThroughputOps)
		fmt.Printf("  Mémoire allouée: %d bytes\n", report.MemoryUsage.AllocBytes)
		
		if len(report.CacheStats) > 0 {
			fmt.Printf("  Stats cache:\n")
			for key, value := range report.CacheStats {
				fmt.Printf("    %s: %v\n", key, value)
			}
		}
		
		fmt.Printf("  Erreurs: %d\n\n", report.ErrorCount)
	}
}

// ComparisonResult représente le résultat d'une comparaison de performances
type ComparisonResult struct {
	BaselineName    string        `json:"baseline_name"`
	OptimizedName   string        `json:"optimized_name"`
	SpeedupFactor   float64       `json:"speedup_factor"`
	MemoryReduction float64       `json:"memory_reduction_percent"`
	ThroughputGain  float64       `json:"throughput_gain_percent"`
	Improvement     string        `json:"improvement_description"`
}

// ComparePerformance compare deux rapports de performances
func ComparePerformance(baseline, optimized PerformanceReport) ComparisonResult {
	speedup := float64(baseline.Duration) / float64(optimized.Duration)
	
	memoryReduction := 0.0
	if baseline.MemoryUsage.AllocBytes > 0 {
		memoryReduction = (1.0 - float64(optimized.MemoryUsage.AllocBytes)/float64(baseline.MemoryUsage.AllocBytes)) * 100
	}
	
	throughputGain := ((optimized.ThroughputOps - baseline.ThroughputOps) / baseline.ThroughputOps) * 100
	
	var improvement string
	if speedup > 1.5 && memoryReduction > 0 {
		improvement = "Excellente optimisation"
	} else if speedup > 1.2 {
		improvement = "Bonne optimisation"
	} else if speedup > 1.0 {
		improvement = "Optimisation légère"
	} else {
		improvement = "Pas d'amélioration significative"
	}
	
	return ComparisonResult{
		BaselineName:    baseline.TestName,
		OptimizedName:   optimized.TestName,
		SpeedupFactor:   speedup,
		MemoryReduction: memoryReduction,
		ThroughputGain:  throughputGain,
		Improvement:     improvement,
	}
}

// LoadTestConfig représente la configuration d'un test de charge
type LoadTestConfig struct {
	NumWorkers      int           `json:"num_workers"`
	TestDuration    time.Duration `json:"test_duration"`
	RampUpDuration  time.Duration `json:"ramp_up_duration"`
	TargetRPS       float64       `json:"target_requests_per_second"`
	MaxConcurrency  int           `json:"max_concurrency"`
}

// LoadTestResult représente le résultat d'un test de charge
type LoadTestResult struct {
	Config           LoadTestConfig `json:"config"`
	TotalRequests    int64         `json:"total_requests"`
	SuccessfulReqs   int64         `json:"successful_requests"`
	FailedRequests   int64         `json:"failed_requests"`
	AverageLatency   time.Duration `json:"average_latency"`
	P95Latency       time.Duration `json:"p95_latency"`
	P99Latency       time.Duration `json:"p99_latency"`
	MaxLatency       time.Duration `json:"max_latency"`
	ActualRPS        float64       `json:"actual_rps"`
	ErrorRate        float64       `json:"error_rate_percent"`
	ThroughputMBps   float64       `json:"throughput_mbps"`
}

// PerformanceThresholds définit les seuils de performance acceptables
type PerformanceThresholds struct {
	MaxLatency        time.Duration `json:"max_latency"`
	MinThroughput     float64       `json:"min_throughput_ops_per_sec"`
	MaxMemoryUsage    uint64        `json:"max_memory_usage_bytes"`
	MaxErrorRate      float64       `json:"max_error_rate_percent"`
	RequiredSpeedup   float64       `json:"required_speedup_factor"`
}

// ValidatePerformance valide qu'un rapport respecte les seuils définis
func ValidatePerformance(report PerformanceReport, thresholds PerformanceThresholds) []string {
	var violations []string
	
	avgLatency := time.Duration(int64(report.Duration) / report.OperationsCount)
	if avgLatency > thresholds.MaxLatency {
		violations = append(violations, 
			fmt.Sprintf("Latence moyenne (%v) dépasse le seuil (%v)", avgLatency, thresholds.MaxLatency))
	}
	
	if report.ThroughputOps < thresholds.MinThroughput {
		violations = append(violations, 
			fmt.Sprintf("Débit (%.2f ops/sec) en dessous du seuil (%.2f ops/sec)", 
				report.ThroughputOps, thresholds.MinThroughput))
	}
	
	if report.MemoryUsage.AllocBytes > thresholds.MaxMemoryUsage {
		violations = append(violations, 
			fmt.Sprintf("Usage mémoire (%d bytes) dépasse le seuil (%d bytes)", 
				report.MemoryUsage.AllocBytes, thresholds.MaxMemoryUsage))
	}
	
	errorRate := float64(report.ErrorCount) / float64(report.OperationsCount) * 100
	if errorRate > thresholds.MaxErrorRate {
		violations = append(violations, 
			fmt.Sprintf("Taux d'erreur (%.2f%%) dépasse le seuil (%.2f%%)", 
				errorRate, thresholds.MaxErrorRate))
	}
	
	return violations
}

// OptimizationSuggestion représente une suggestion d'optimisation
type OptimizationSuggestion struct {
	Component   string  `json:"component"`
	Issue       string  `json:"issue"`
	Suggestion  string  `json:"suggestion"`
	Priority    string  `json:"priority"` // "HIGH", "MEDIUM", "LOW"
	Impact      string  `json:"estimated_impact"`
}

// AnalyzePerformance analyse les performances et suggère des optimisations
func AnalyzePerformance(reports []PerformanceReport) []OptimizationSuggestion {
	var suggestions []OptimizationSuggestion
	
	for _, report := range reports {
		// Analyse de la latence
		avgLatency := time.Duration(int64(report.Duration) / report.OperationsCount)
		if avgLatency > 10*time.Millisecond {
			suggestions = append(suggestions, OptimizationSuggestion{
				Component:  report.ComponentName,
				Issue:      "Latence élevée détectée",
				Suggestion: "Considérer l'ajout de caches ou l'optimisation des algorithmes",
				Priority:   "HIGH",
				Impact:     "Amélioration potentielle de 2-5x",
			})
		}
		
		// Analyse de l'usage mémoire
		if report.MemoryUsage.AllocBytes > 100*1024*1024 { // > 100MB
			suggestions = append(suggestions, OptimizationSuggestion{
				Component:  report.ComponentName,
				Issue:      "Usage mémoire élevé",
				Suggestion: "Implémenter le pooling d'objets et optimiser les structures de données",
				Priority:   "MEDIUM",
				Impact:     "Réduction mémoire de 20-50%",
			})
		}
		
		// Analyse du débit
		if report.ThroughputOps < 1000 {
			suggestions = append(suggestions, OptimizationSuggestion{
				Component:  report.ComponentName,
				Issue:      "Débit faible",
				Suggestion: "Parallelisation et optimisation des goulots d'étranglement",
				Priority:   "HIGH",
				Impact:     "Amélioration du débit de 3-10x",
			})
		}
		
		// Analyse des GC
		if report.MemoryUsage.NumGC > 10 {
			suggestions = append(suggestions, OptimizationSuggestion{
				Component:  report.ComponentName,
				Issue:      "Trop de collections de garbage",
				Suggestion: "Réduire les allocations et utiliser des pools d'objets",
				Priority:   "MEDIUM",
				Impact:     "Réduction des pauses GC de 50-80%",
			})
		}
	}
	
	return suggestions
}

// PrintOptimizationReport affiche un rapport d'optimisation
func PrintOptimizationReport(suggestions []OptimizationSuggestion) {
	fmt.Println("\n=== RAPPORT D'OPTIMISATION ===")
	
	if len(suggestions) == 0 {
		fmt.Println("Aucune optimisation critique détectée.")
		return
	}
	
	high := 0
	medium := 0
	low := 0
	
	for _, suggestion := range suggestions {
		switch suggestion.Priority {
		case "HIGH":
			high++
		case "MEDIUM":
			medium++
		case "LOW":
			low++
		}
		
		fmt.Printf("\n[%s] %s - %s\n", suggestion.Priority, suggestion.Component, suggestion.Issue)
		fmt.Printf("Suggestion: %s\n", suggestion.Suggestion)
		fmt.Printf("Impact estimé: %s\n", suggestion.Impact)
	}
	
	fmt.Printf("\nRésumé des priorités:\n")
	fmt.Printf("  - Haute priorité: %d\n", high)
	fmt.Printf("  - Priorité moyenne: %d\n", medium)
	fmt.Printf("  - Basse priorité: %d\n", low)
}