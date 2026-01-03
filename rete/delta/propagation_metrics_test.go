// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestNewPropagationMetrics(t *testing.T) {
	t.Log("🧪 TEST: NewPropagationMetrics - Création instance métriques propagation")

	metrics := NewPropagationMetrics()

	if metrics == nil {
		t.Fatal("❌ NewPropagationMetrics() devrait retourner une instance non-nil")
	}

	// Vérifier initialisation des compteurs
	snapshot := metrics.GetSnapshot()

	if snapshot.DeltaPropagations != 0 {
		t.Errorf("❌ DeltaPropagations initial devrait être 0, got %d", snapshot.DeltaPropagations)
	}

	if snapshot.ClassicPropagations != 0 {
		t.Errorf("❌ ClassicPropagations initial devrait être 0, got %d", snapshot.ClassicPropagations)
	}

	if snapshot.FailedPropagations != 0 {
		t.Errorf("❌ FailedPropagations initial devrait être 0, got %d", snapshot.FailedPropagations)
	}

	if snapshot.TotalNodesEvaluated != 0 {
		t.Errorf("❌ TotalNodesEvaluated initial devrait être 0, got %d", snapshot.TotalNodesEvaluated)
	}

	if snapshot.NodesSkippedByDelta != 0 {
		t.Errorf("❌ NodesSkippedByDelta initial devrait être 0, got %d", snapshot.NodesSkippedByDelta)
	}

	if snapshot.MinPropagationTime != time.Duration(1<<63-1) {
		t.Errorf("❌ MinPropagationTime devrait être max duration, got %v", snapshot.MinPropagationTime)
	}

	t.Log("✅ Métriques de propagation initialisées correctement")
}

func TestPropagationMetrics_RecordDeltaPropagation(t *testing.T) {
	t.Log("🧪 TEST: RecordDeltaPropagation - Enregistrement propagation delta")

	metrics := NewPropagationMetrics()

	// Enregistrer une propagation delta
	metrics.RecordDeltaPropagation(5*time.Millisecond, 3, 2)

	snapshot := metrics.GetSnapshot()

	if snapshot.DeltaPropagations != 1 {
		t.Errorf("❌ DeltaPropagations devrait être 1, got %d", snapshot.DeltaPropagations)
	}

	if snapshot.TotalPropagations != 1 {
		t.Errorf("❌ TotalPropagations devrait être 1, got %d", snapshot.TotalPropagations)
	}

	if snapshot.TotalNodesEvaluated != 3 {
		t.Errorf("❌ TotalNodesEvaluated devrait être 3, got %d", snapshot.TotalNodesEvaluated)
	}

	if snapshot.TotalFieldsChanged != 2 {
		t.Errorf("❌ TotalFieldsChanged devrait être 2, got %d", snapshot.TotalFieldsChanged)
	}

	if snapshot.TotalPropagationTime != 5*time.Millisecond {
		t.Errorf("❌ TotalPropagationTime devrait être 5ms, got %v", snapshot.TotalPropagationTime)
	}

	// Enregistrer une autre propagation
	metrics.RecordDeltaPropagation(10*time.Millisecond, 5, 3)

	snapshot = metrics.GetSnapshot()

	if snapshot.DeltaPropagations != 2 {
		t.Errorf("❌ DeltaPropagations devrait être 2, got %d", snapshot.DeltaPropagations)
	}

	if snapshot.TotalNodesEvaluated != 8 {
		t.Errorf("❌ TotalNodesEvaluated devrait être 8 (3+5), got %d", snapshot.TotalNodesEvaluated)
	}

	if snapshot.TotalFieldsChanged != 5 {
		t.Errorf("❌ TotalFieldsChanged devrait être 5 (2+3), got %d", snapshot.TotalFieldsChanged)
	}

	if snapshot.TotalPropagationTime != 15*time.Millisecond {
		t.Errorf("❌ TotalPropagationTime devrait être 15ms, got %v", snapshot.TotalPropagationTime)
	}

	t.Log("✅ RecordDeltaPropagation fonctionne correctement")
}

func TestPropagationMetrics_RecordClassicPropagation(t *testing.T) {
	t.Log("🧪 TEST: RecordClassicPropagation - Enregistrement propagation classique")

	metrics := NewPropagationMetrics()

	// Enregistrer une propagation classique
	metrics.RecordClassicPropagation(20*time.Millisecond, 10)

	snapshot := metrics.GetSnapshot()

	if snapshot.ClassicPropagations != 1 {
		t.Errorf("❌ ClassicPropagations devrait être 1, got %d", snapshot.ClassicPropagations)
	}

	if snapshot.TotalPropagations != 1 {
		t.Errorf("❌ TotalPropagations devrait être 1, got %d", snapshot.TotalPropagations)
	}

	if snapshot.TotalNodesEvaluated != 10 {
		t.Errorf("❌ TotalNodesEvaluated devrait être 10, got %d", snapshot.TotalNodesEvaluated)
	}

	if snapshot.TotalPropagationTime != 20*time.Millisecond {
		t.Errorf("❌ TotalPropagationTime devrait être 20ms, got %v", snapshot.TotalPropagationTime)
	}

	t.Log("✅ RecordClassicPropagation fonctionne correctement")
}

func TestPropagationMetrics_RecordFailedPropagation(t *testing.T) {
	t.Log("🧪 TEST: RecordFailedPropagation - Enregistrement échecs propagation")

	metrics := NewPropagationMetrics()

	// Enregistrer un échec
	metrics.RecordFailedPropagation()

	snapshot := metrics.GetSnapshot()

	if snapshot.FailedPropagations != 1 {
		t.Errorf("❌ FailedPropagations devrait être 1, got %d", snapshot.FailedPropagations)
	}

	if snapshot.TotalPropagations != 1 {
		t.Errorf("❌ TotalPropagations devrait être 1, got %d", snapshot.TotalPropagations)
	}

	// Enregistrer plusieurs échecs
	metrics.RecordFailedPropagation()
	metrics.RecordFailedPropagation()

	snapshot = metrics.GetSnapshot()

	if snapshot.FailedPropagations != 3 {
		t.Errorf("❌ FailedPropagations devrait être 3, got %d", snapshot.FailedPropagations)
	}

	if snapshot.TotalPropagations != 3 {
		t.Errorf("❌ TotalPropagations devrait être 3, got %d", snapshot.TotalPropagations)
	}

	t.Log("✅ RecordFailedPropagation fonctionne correctement")
}

func TestPropagationMetrics_RecordFallback(t *testing.T) {
	t.Log("🧪 TEST: RecordFallback - Enregistrement raison fallback")

	metrics := NewPropagationMetrics()

	// Enregistrer des raisons de fallback
	metrics.RecordFallback("ratio")
	metrics.RecordFallback("nodes")
	metrics.RecordFallback("pk")
	metrics.RecordFallback("error")
	metrics.RecordFallback("ratio") // Duplicate

	snapshot := metrics.GetSnapshot()

	// Vérifier que les raisons sont comptées
	if snapshot.FallbacksDueToRatio != 2 {
		t.Errorf("❌ FallbacksDueToRatio devrait être 2, got %d", snapshot.FallbacksDueToRatio)
	}

	if snapshot.FallbacksDueToNodes != 1 {
		t.Errorf("❌ FallbacksDueToNodes devrait être 1, got %d", snapshot.FallbacksDueToNodes)
	}

	if snapshot.FallbacksDueToPK != 1 {
		t.Errorf("❌ FallbacksDueToPK devrait être 1, got %d", snapshot.FallbacksDueToPK)
	}

	if snapshot.FallbacksDueToError != 1 {
		t.Errorf("❌ FallbacksDueToError devrait être 1, got %d", snapshot.FallbacksDueToError)
	}

	t.Log("✅ RecordFallback fonctionne correctement")
}

func TestPropagationMetrics_RecordNodesSkipped(t *testing.T) {
	t.Log("🧪 TEST: RecordNodesSkipped - Enregistrement nœuds ignorés")

	metrics := NewPropagationMetrics()

	// Enregistrer des nœuds ignorés
	metrics.RecordNodesSkipped(5)

	snapshot := metrics.GetSnapshot()

	if snapshot.NodesSkippedByDelta != 5 {
		t.Errorf("❌ NodesSkippedByDelta devrait être 5, got %d", snapshot.NodesSkippedByDelta)
	}

	// Enregistrer plus de nœuds ignorés
	metrics.RecordNodesSkipped(3)

	snapshot = metrics.GetSnapshot()

	if snapshot.NodesSkippedByDelta != 8 {
		t.Errorf("❌ NodesSkippedByDelta devrait être 8 (5+3), got %d", snapshot.NodesSkippedByDelta)
	}

	t.Log("✅ RecordNodesSkipped fonctionne correctement")
}

func TestPropagationMetrics_GetEfficiencyRatio(t *testing.T) {
	t.Log("🧪 TEST: GetEfficiencyRatio - Calcul ratio efficacité")

	metrics := NewPropagationMetrics()

	// Sans propagations
	ratio := metrics.GetEfficiencyRatio()
	if ratio != 0.0 {
		t.Errorf("❌ EfficiencyRatio sans propagations devrait être 0.0, got %f", ratio)
	}

	// Avec propagations et nœuds évités
	metrics.RecordDeltaPropagation(5*time.Millisecond, 10, 2)
	metrics.RecordNodesSkipped(90) // Sur 100 nœuds totaux, 90 évités

	ratio = metrics.GetEfficiencyRatio()

	// Ratio = NodesSkipped / (NodesEvaluated + NodesSkipped) = 90 / 100 = 0.9
	expectedRatio := 0.9
	if ratio < expectedRatio-0.01 || ratio > expectedRatio+0.01 {
		t.Errorf("❌ EfficiencyRatio devrait être ~%.2f, got %.2f", expectedRatio, ratio)
	}

	t.Log("✅ GetEfficiencyRatio calcule correctement le ratio")
}

func TestPropagationMetrics_GetDeltaUsageRatio(t *testing.T) {
	t.Log("🧪 TEST: GetDeltaUsageRatio - Calcul ratio utilisation delta")

	metrics := NewPropagationMetrics()

	// Sans propagations
	ratio := metrics.GetDeltaUsageRatio()
	if ratio != 0.0 {
		t.Errorf("❌ DeltaUsageRatio sans propagations devrait être 0.0, got %f", ratio)
	}

	// 3 delta, 1 classic = 75% delta usage
	metrics.RecordDeltaPropagation(1*time.Millisecond, 1, 1)
	metrics.RecordDeltaPropagation(1*time.Millisecond, 1, 1)
	metrics.RecordDeltaPropagation(1*time.Millisecond, 1, 1)
	metrics.RecordClassicPropagation(1*time.Millisecond, 1)

	ratio = metrics.GetDeltaUsageRatio()
	expected := 0.75

	if ratio < expected-0.01 || ratio > expected+0.01 {
		t.Errorf("❌ DeltaUsageRatio devrait être ~%.2f, got %.2f", expected, ratio)
	}

	t.Log("✅ GetDeltaUsageRatio calcule correctement le ratio")
}

func TestPropagationMetrics_Reset(t *testing.T) {
	t.Log("🧪 TEST: Reset - Réinitialisation métriques propagation")

	metrics := NewPropagationMetrics()

	// Remplir avec des données
	metrics.RecordDeltaPropagation(10*time.Millisecond, 5, 2)
	metrics.RecordClassicPropagation(20*time.Millisecond, 10)
	metrics.RecordFailedPropagation()
	metrics.RecordFallback("ratio")
	metrics.RecordNodesSkipped(3)

	// Vérifier que les données sont présentes
	snapshot := metrics.GetSnapshot()
	if snapshot.DeltaPropagations == 0 {
		t.Fatal("❌ DeltaPropagations devrait être > 0 avant reset")
	}

	// Reset
	metrics.Reset()

	// Vérifier réinitialisation
	snapshot = metrics.GetSnapshot()

	if snapshot.DeltaPropagations != 0 {
		t.Errorf("❌ DeltaPropagations après reset devrait être 0, got %d", snapshot.DeltaPropagations)
	}

	if snapshot.ClassicPropagations != 0 {
		t.Errorf("❌ ClassicPropagations après reset devrait être 0, got %d", snapshot.ClassicPropagations)
	}

	if snapshot.FailedPropagations != 0 {
		t.Errorf("❌ FailedPropagations après reset devrait être 0, got %d", snapshot.FailedPropagations)
	}

	if snapshot.TotalNodesEvaluated != 0 {
		t.Errorf("❌ TotalNodesEvaluated après reset devrait être 0, got %d", snapshot.TotalNodesEvaluated)
	}

	if snapshot.NodesSkippedByDelta != 0 {
		t.Errorf("❌ NodesSkippedByDelta après reset devrait être 0, got %d", snapshot.NodesSkippedByDelta)
	}

	if snapshot.TotalPropagationTime != 0 {
		t.Errorf("❌ TotalPropagationTime après reset devrait être 0, got %v", snapshot.TotalPropagationTime)
	}

	if snapshot.TotalFieldsChanged != 0 {
		t.Errorf("❌ TotalFieldsChanged après reset devrait être 0, got %d", snapshot.TotalFieldsChanged)
	}

	if snapshot.FallbacksDueToRatio != 0 {
		t.Errorf("❌ FallbacksDueToRatio après reset devrait être 0, got %d", snapshot.FallbacksDueToRatio)
	}

	t.Log("✅ Reset réinitialise toutes les métriques de propagation")
}

func TestPropagationMetrics_GetSnapshot(t *testing.T) {
	t.Log("🧪 TEST: GetSnapshot - Capture instantané métriques propagation")

	metrics := NewPropagationMetrics()

	// Enregistrer des données
	metrics.RecordDeltaPropagation(10*time.Millisecond, 5, 2)
	metrics.RecordClassicPropagation(30*time.Millisecond, 15)
	metrics.RecordFailedPropagation()
	metrics.RecordFallback("nodes")
	metrics.RecordNodesSkipped(2)

	// Créer snapshot
	snapshot := metrics.GetSnapshot()

	// Vérifier snapshot
	if snapshot.DeltaPropagations != 1 {
		t.Errorf("❌ snapshot.DeltaPropagations devrait être 1, got %d", snapshot.DeltaPropagations)
	}

	if snapshot.ClassicPropagations != 1 {
		t.Errorf("❌ snapshot.ClassicPropagations devrait être 1, got %d", snapshot.ClassicPropagations)
	}

	if snapshot.FailedPropagations != 1 {
		t.Errorf("❌ snapshot.FailedPropagations devrait être 1, got %d", snapshot.FailedPropagations)
	}

	if snapshot.TotalPropagations != 3 {
		t.Errorf("❌ snapshot.TotalPropagations devrait être 3, got %d", snapshot.TotalPropagations)
	}

	if snapshot.TotalNodesEvaluated != 20 { // 5 + 15
		t.Errorf("❌ snapshot.TotalNodesEvaluated devrait être 20, got %d", snapshot.TotalNodesEvaluated)
	}

	if snapshot.NodesSkippedByDelta != 2 {
		t.Errorf("❌ snapshot.NodesSkippedByDelta devrait être 2, got %d", snapshot.NodesSkippedByDelta)
	}

	if snapshot.TotalPropagationTime != 40*time.Millisecond {
		t.Errorf("❌ snapshot.TotalPropagationTime devrait être 40ms, got %v", snapshot.TotalPropagationTime)
	}

	if snapshot.FallbacksDueToNodes != 1 {
		t.Errorf("❌ snapshot.FallbacksDueToNodes devrait être 1, got %d", snapshot.FallbacksDueToNodes)
	}

	// Modifier les métriques originales ne devrait pas affecter le snapshot
	metrics.RecordDeltaPropagation(5*time.Millisecond, 3, 1)

	if snapshot.DeltaPropagations != 1 {
		t.Error("❌ Le snapshot devrait être immutable")
	}

	t.Log("✅ Snapshot capture correctement les métriques de propagation")
}

func TestPropagationMetrics_Averages(t *testing.T) {
	t.Log("🧪 TEST: Averages - Calcul moyennes temps propagation")

	metrics := NewPropagationMetrics()

	// Enregistrer plusieurs propagations
	metrics.RecordDeltaPropagation(10*time.Millisecond, 5, 1)
	metrics.RecordDeltaPropagation(20*time.Millisecond, 10, 2)
	metrics.RecordDeltaPropagation(30*time.Millisecond, 15, 3)

	snapshot := metrics.GetSnapshot()

	// Moyenne = (10+20+30)/3 = 20ms
	expectedAvg := 20 * time.Millisecond
	if snapshot.AvgPropagationTime != expectedAvg {
		t.Errorf("❌ AvgPropagationTime devrait être %v, got %v", expectedAvg, snapshot.AvgPropagationTime)
	}

	// Moyenne nœuds = (5+10+15)/3 = 10
	expectedAvgNodes := 10.0
	if snapshot.AvgNodesPerPropagation < expectedAvgNodes-0.1 || snapshot.AvgNodesPerPropagation > expectedAvgNodes+0.1 {
		t.Errorf("❌ AvgNodesPerPropagation devrait être ~%.1f, got %.1f", expectedAvgNodes, snapshot.AvgNodesPerPropagation)
	}

	// Moyenne fields = (1+2+3)/3 = 2
	expectedAvgFields := 2.0
	if snapshot.AvgFieldsPerPropagation < expectedAvgFields-0.1 || snapshot.AvgFieldsPerPropagation > expectedAvgFields+0.1 {
		t.Errorf("❌ AvgFieldsPerPropagation devrait être ~%.1f, got %.1f", expectedAvgFields, snapshot.AvgFieldsPerPropagation)
	}

	t.Log("✅ Calcul des moyennes correct")
}

func TestPropagationMetrics_MinMaxTiming(t *testing.T) {
	t.Log("🧪 TEST: MinMaxTiming - Suivi min/max temps propagation")

	metrics := NewPropagationMetrics()

	// Enregistrer propagations avec temps variés
	metrics.RecordDeltaPropagation(20*time.Millisecond, 5, 1)
	metrics.RecordDeltaPropagation(5*time.Millisecond, 3, 1)
	metrics.RecordDeltaPropagation(50*time.Millisecond, 10, 2)

	snapshot := metrics.GetSnapshot()

	if snapshot.MinPropagationTime != 5*time.Millisecond {
		t.Errorf("❌ MinPropagationTime devrait être 5ms, got %v", snapshot.MinPropagationTime)
	}

	if snapshot.MaxPropagationTime != 50*time.Millisecond {
		t.Errorf("❌ MaxPropagationTime devrait être 50ms, got %v", snapshot.MaxPropagationTime)
	}

	t.Log("✅ Suivi min/max temps fonctionne correctement")
}

func TestPropagationMetrics_Timestamps(t *testing.T) {
	t.Log("🧪 TEST: Timestamps - Suivi timestamps propagations")

	metrics := NewPropagationMetrics()

	snapshot := metrics.GetSnapshot()
	if !snapshot.FirstPropagation.IsZero() {
		t.Error("❌ FirstPropagation devrait être zero avant toute propagation")
	}

	// Première propagation
	before := time.Now()
	metrics.RecordDeltaPropagation(1*time.Millisecond, 1, 1)
	after := time.Now()

	snapshot = metrics.GetSnapshot()

	if snapshot.FirstPropagation.Before(before) || snapshot.FirstPropagation.After(after) {
		t.Error("❌ FirstPropagation devrait être dans l'intervalle de temps")
	}

	firstTime := snapshot.FirstPropagation

	// Deuxième propagation
	time.Sleep(2 * time.Millisecond)
	metrics.RecordDeltaPropagation(1*time.Millisecond, 1, 1)

	snapshot = metrics.GetSnapshot()

	// FirstPropagation ne devrait pas changer
	if snapshot.FirstPropagation != firstTime {
		t.Error("❌ FirstPropagation ne devrait pas changer après la première propagation")
	}

	// LastPropagation devrait être après FirstPropagation
	if !snapshot.LastPropagation.After(snapshot.FirstPropagation) {
		t.Error("❌ LastPropagation devrait être après FirstPropagation")
	}

	t.Log("✅ Timestamps suivent correctement les propagations")
}

func TestPropagationMetrics_ConcurrentAccess(t *testing.T) {
	t.Log("🧪 TEST: ConcurrentAccess - Accès concurrent métriques propagation")

	metrics := NewPropagationMetrics()

	done := make(chan bool)

	// Goroutine 1: RecordDeltaPropagation
	go func() {
		for i := 0; i < 50; i++ {
			metrics.RecordDeltaPropagation(time.Millisecond, 1, 1)
		}
		done <- true
	}()

	// Goroutine 2: RecordClassicPropagation
	go func() {
		for i := 0; i < 50; i++ {
			metrics.RecordClassicPropagation(2*time.Millisecond, 2)
		}
		done <- true
	}()

	// Goroutine 3: GetSnapshot
	go func() {
		for i := 0; i < 30; i++ {
			_ = metrics.GetSnapshot()
		}
		done <- true
	}()

	// Attendre toutes les goroutines
	<-done
	<-done
	<-done

	// Vérifier cohérence
	snapshot := metrics.GetSnapshot()

	if snapshot.DeltaPropagations != 50 {
		t.Errorf("❌ DeltaPropagations devrait être 50, got %d", snapshot.DeltaPropagations)
	}

	if snapshot.ClassicPropagations != 50 {
		t.Errorf("❌ ClassicPropagations devrait être 50, got %d", snapshot.ClassicPropagations)
	}

	if snapshot.TotalPropagations != 100 {
		t.Errorf("❌ TotalPropagations devrait être 100, got %d", snapshot.TotalPropagations)
	}

	t.Log("✅ Test concurrent terminé avec succès (mutex protège l'accès)")
}
