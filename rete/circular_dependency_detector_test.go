// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestNewCircularDependencyDetector teste la création du détecteur
func TestNewCircularDependencyDetector(t *testing.T) {
	detector := NewCircularDependencyDetector()
	if detector == nil {
		t.Fatal("Expected detector to be created, got nil")
	}
	if detector.graph == nil {
		t.Error("Expected graph to be initialized")
	}
	if detector.colors == nil {
		t.Error("Expected colors to be initialized")
	}
	if detector.metadata == nil {
		t.Error("Expected metadata to be initialized")
	}
}

// TestCircularDependency_NoCycle teste un graphe sans cycle
func TestCircularDependency_NoCycle(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Graphe linéaire: A -> B -> C
	detector.AddNode("A", []string{"B"})
	detector.AddNode("B", []string{"C"})
	detector.AddNode("C", []string{})
	hasCycle := detector.DetectCycles()
	if hasCycle {
		t.Errorf("Expected no cycle, but cycle detected: %v", detector.GetCyclePath())
	}
}

// TestCircularDependency_SimpleCycle teste un cycle simple
func TestCircularDependency_SimpleCycle(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Cycle simple: A -> B -> A
	detector.AddNode("A", []string{"B"})
	detector.AddNode("B", []string{"A"})
	hasCycle := detector.DetectCycles()
	if !hasCycle {
		t.Error("Expected cycle to be detected")
	}
	cyclePath := detector.GetCyclePath()
	if len(cyclePath) == 0 {
		t.Error("Expected cycle path to be non-empty")
	}
	t.Logf("Cycle path: %s", detector.FormatCyclePath())
	// Vérifier que le cycle contient A et B
	hasA := false
	hasB := false
	for _, node := range cyclePath {
		if node == "A" {
			hasA = true
		}
		if node == "B" {
			hasB = true
		}
	}
	if !hasA || !hasB {
		t.Errorf("Expected cycle to contain A and B, got %v", cyclePath)
	}
}

// TestCircularDependency_SelfCycle teste un cycle sur soi-même
func TestCircularDependency_SelfCycle(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// A dépend de lui-même
	detector.AddNode("A", []string{"A"})
	hasCycle := detector.DetectCycles()
	if !hasCycle {
		t.Error("Expected self-cycle to be detected")
	}
	cyclePath := detector.GetCyclePath()
	t.Logf("Self-cycle path: %s", detector.FormatCyclePath())
	if len(cyclePath) < 2 {
		t.Errorf("Expected cycle path with at least 2 nodes (start and end), got %d", len(cyclePath))
	}
}

// TestCircularDependency_ComplexCycle teste un cycle complexe
func TestCircularDependency_ComplexCycle(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Cycle: A -> B -> C -> D -> B
	detector.AddNode("A", []string{"B"})
	detector.AddNode("B", []string{"C"})
	detector.AddNode("C", []string{"D"})
	detector.AddNode("D", []string{"B"})
	hasCycle := detector.DetectCycles()
	if !hasCycle {
		t.Error("Expected cycle to be detected")
	}
	cyclePath := detector.GetCyclePath()
	t.Logf("Complex cycle path: %s", detector.FormatCyclePath())
	if len(cyclePath) < 3 {
		t.Errorf("Expected cycle path with at least 3 nodes, got %d", len(cyclePath))
	}
}

// TestCircularDependency_MultiplePaths teste un graphe avec plusieurs chemins
func TestCircularDependency_MultiplePaths(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Graphe en diamant sans cycle:
	//     A
	//    / \
	//   B   C
	//    \ /
	//     D
	detector.AddNode("A", []string{"B", "C"})
	detector.AddNode("B", []string{"D"})
	detector.AddNode("C", []string{"D"})
	detector.AddNode("D", []string{})
	hasCycle := detector.DetectCycles()
	if hasCycle {
		t.Errorf("Expected no cycle in diamond graph, but cycle detected: %v", detector.GetCyclePath())
	}
}

// TestCircularDependency_DisconnectedGraph teste un graphe déconnecté
func TestCircularDependency_DisconnectedGraph(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Deux composantes déconnectées:
	// Composante 1: A -> B
	// Composante 2: C -> D -> C (avec cycle)
	detector.AddNode("A", []string{"B"})
	detector.AddNode("B", []string{})
	detector.AddNode("C", []string{"D"})
	detector.AddNode("D", []string{"C"})
	hasCycle := detector.DetectCycles()
	if !hasCycle {
		t.Error("Expected cycle to be detected in disconnected graph")
	}
	t.Logf("Cycle in disconnected graph: %s", detector.FormatCyclePath())
}

// TestCircularDependency_EmptyGraph teste un graphe vide
func TestCircularDependency_EmptyGraph(t *testing.T) {
	detector := NewCircularDependencyDetector()
	hasCycle := detector.DetectCycles()
	if hasCycle {
		t.Error("Expected no cycle in empty graph")
	}
}

// TestCircularDependency_Validate teste la validation complète
func TestCircularDependency_Validate(t *testing.T) {
	tests := []struct {
		name          string
		nodes         map[string][]string
		expectValid   bool
		expectCycle   bool
		expectWarning bool
	}{
		{
			name: "valid_linear",
			nodes: map[string][]string{
				"temp_1": {},
				"temp_2": {"temp_1"},
				"temp_3": {"temp_2"},
			},
			expectValid: true,
			expectCycle: false,
		},
		{
			name: "cycle_detected",
			nodes: map[string][]string{
				"temp_1": {"temp_2"},
				"temp_2": {"temp_1"},
			},
			expectValid: false,
			expectCycle: true,
		},
		{
			name: "missing_dependency",
			nodes: map[string][]string{
				"temp_1": {"temp_missing"},
			},
			expectValid: false,
			expectCycle: false,
		},
		{
			name: "isolated_node",
			nodes: map[string][]string{
				"temp_1": {},
				"temp_2": {"temp_3"},
				"temp_3": {},
			},
			expectValid:   true,
			expectWarning: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := NewCircularDependencyDetector()
			for node, deps := range tt.nodes {
				detector.AddNode(node, deps)
			}
			report := detector.Validate()
			if report.Valid != tt.expectValid {
				t.Errorf("Expected Valid=%v, got %v. Error: %s",
					tt.expectValid, report.Valid, report.ErrorMessage)
			}
			if report.HasCircularDeps != tt.expectCycle {
				t.Errorf("Expected HasCircularDeps=%v, got %v",
					tt.expectCycle, report.HasCircularDeps)
			}
			if tt.expectWarning && len(report.Warnings) == 0 {
				t.Error("Expected warnings, got none")
			}
			t.Logf("Report for %s: Valid=%v, Cycle=%v, MaxDepth=%d, Warnings=%d",
				tt.name, report.Valid, report.HasCircularDeps, report.MaxDepth, len(report.Warnings))
		})
	}
}

// TestCircularDependency_MaxDepth teste le calcul de la profondeur maximale
func TestCircularDependency_MaxDepth(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Chaîne de profondeur 5:
	// temp_5 -> temp_4 -> temp_3 -> temp_2 -> temp_1 -> (rien)
	detector.AddNode("temp_1", []string{})
	detector.AddNode("temp_2", []string{"temp_1"})
	detector.AddNode("temp_3", []string{"temp_2"})
	detector.AddNode("temp_4", []string{"temp_3"})
	detector.AddNode("temp_5", []string{"temp_4"})
	report := detector.Validate()
	if !report.Valid {
		t.Errorf("Expected valid graph, got error: %s", report.ErrorMessage)
	}
	expectedDepth := 4 // Profondeur = nombre d'arêtes
	if report.MaxDepth != expectedDepth {
		t.Errorf("Expected MaxDepth=%d, got %d", expectedDepth, report.MaxDepth)
	}
	t.Logf("Max depth for chain of 5 nodes: %d", report.MaxDepth)
}

// TestCircularDependency_TopologicalSort teste le tri topologique
func TestCircularDependency_TopologicalSort(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// DAG: temp_3 -> temp_2 -> temp_1
	detector.AddNode("temp_1", []string{})
	detector.AddNode("temp_2", []string{"temp_1"})
	detector.AddNode("temp_3", []string{"temp_2"})
	sorted, err := detector.GetTopologicalSort()
	if err != nil {
		t.Fatalf("Expected successful topological sort, got error: %v", err)
	}
	if len(sorted) != 3 {
		t.Errorf("Expected 3 nodes in sorted order, got %d", len(sorted))
	}
	t.Logf("Topological sort: %v", sorted)
	// Vérifier que temp_1 vient avant temp_2
	idx1 := -1
	idx2 := -1
	idx3 := -1
	for i, node := range sorted {
		if node == "temp_1" {
			idx1 = i
		}
		if node == "temp_2" {
			idx2 = i
		}
		if node == "temp_3" {
			idx3 = i
		}
	}
	if idx1 >= idx2 {
		t.Errorf("Expected temp_1 before temp_2, got indices %d, %d", idx1, idx2)
	}
	if idx2 >= idx3 {
		t.Errorf("Expected temp_2 before temp_3, got indices %d, %d", idx2, idx3)
	}
}

// TestCircularDependency_TopologicalSort_WithCycle teste que le tri échoue avec un cycle
func TestCircularDependency_TopologicalSort_WithCycle(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Cycle: temp_1 -> temp_2 -> temp_1
	detector.AddNode("temp_1", []string{"temp_2"})
	detector.AddNode("temp_2", []string{"temp_1"})
	_, err := detector.GetTopologicalSort()
	if err == nil {
		t.Error("Expected error for topological sort with cycle, got nil")
	}
	t.Logf("Expected error: %v", err)
}

// TestCircularDependency_GetDependencyChain teste la récupération de la chaîne complète
func TestCircularDependency_GetDependencyChain(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Graphe:
	// temp_3 -> temp_2 -> temp_1
	//        \-> temp_4
	detector.AddNode("temp_1", []string{})
	detector.AddNode("temp_2", []string{"temp_1"})
	detector.AddNode("temp_3", []string{"temp_2", "temp_4"})
	detector.AddNode("temp_4", []string{})
	chain := detector.GetDependencyChain("temp_3")
	if len(chain) != 4 {
		t.Errorf("Expected chain of length 4, got %d: %v", len(chain), chain)
	}
	t.Logf("Dependency chain for temp_3: %v", chain)
	// Vérifier que tous les nœuds sont présents
	hasTemp3 := false
	hasTemp2 := false
	hasTemp1 := false
	hasTemp4 := false
	for _, node := range chain {
		if node == "temp_3" {
			hasTemp3 = true
		}
		if node == "temp_2" {
			hasTemp2 = true
		}
		if node == "temp_1" {
			hasTemp1 = true
		}
		if node == "temp_4" {
			hasTemp4 = true
		}
	}
	if !hasTemp3 || !hasTemp2 || !hasTemp1 || !hasTemp4 {
		t.Errorf("Expected all nodes in chain, got %v", chain)
	}
}

// TestCircularDependency_GetStatistics teste les statistiques
func TestCircularDependency_GetStatistics(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Graphe simple
	detector.AddNode("temp_1", []string{})
	detector.AddNode("temp_2", []string{"temp_1"})
	detector.AddNode("temp_3", []string{"temp_1", "temp_2"})
	stats := detector.GetStatistics()
	if stats["total_nodes"].(int) != 3 {
		t.Errorf("Expected 3 nodes, got %d", stats["total_nodes"])
	}
	if stats["total_edges"].(int) != 3 {
		t.Errorf("Expected 3 edges, got %d", stats["total_edges"])
	}
	if stats["has_cycles"].(bool) {
		t.Error("Expected no cycles")
	}
	t.Logf("Graph statistics: %+v", stats)
}

// TestCircularDependency_Clear teste la réinitialisation
func TestCircularDependency_Clear(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Ajouter des nœuds
	detector.AddNode("temp_1", []string{})
	detector.AddNode("temp_2", []string{"temp_1"})
	// Vérifier qu'ils sont présents
	if len(detector.graph) != 2 {
		t.Errorf("Expected 2 nodes before clear, got %d", len(detector.graph))
	}
	// Clear
	detector.Clear()
	// Vérifier qu'ils sont supprimés
	if len(detector.graph) != 0 {
		t.Errorf("Expected 0 nodes after clear, got %d", len(detector.graph))
	}
	if len(detector.colors) != 0 {
		t.Errorf("Expected 0 colors after clear, got %d", len(detector.colors))
	}
}

// TestCircularDependency_ValidateAlphaChain teste la validation d'une chaîne d'alpha nodes
func TestCircularDependency_ValidateAlphaChain(t *testing.T) {
	detector := NewCircularDependencyDetector()
	storage := NewMemoryStorage()
	// Créer une chaîne valide
	nodes := []*AlphaNode{
		{
			BaseNode: BaseNode{
				ID:      "alpha_1",
				Type:    "alpha",
				Storage: storage,
			},
			ResultName:   "temp_1",
			IsAtomic:     true,
			Dependencies: []string{},
		},
		{
			BaseNode: BaseNode{
				ID:      "alpha_2",
				Type:    "alpha",
				Storage: storage,
			},
			ResultName:   "temp_2",
			IsAtomic:     true,
			Dependencies: []string{"temp_1"},
		},
		{
			BaseNode: BaseNode{
				ID:      "alpha_3",
				Type:    "alpha",
				Storage: storage,
			},
			ResultName:   "temp_3",
			IsAtomic:     true,
			Dependencies: []string{"temp_2"},
		},
	}
	report := detector.ValidateAlphaChain(nodes)
	if !report.Valid {
		t.Errorf("Expected valid chain, got error: %s", report.ErrorMessage)
	}
	if report.HasCircularDeps {
		t.Error("Expected no circular dependencies")
	}
	if report.MaxDepth != 2 {
		t.Errorf("Expected max depth 2, got %d", report.MaxDepth)
	}
	t.Logf("Alpha chain validation report: Valid=%v, MaxDepth=%d, TotalNodes=%d",
		report.Valid, report.MaxDepth, report.TotalNodes)
}

// TestCircularDependency_ValidateAlphaChain_WithCycle teste la détection de cycle dans une chaîne
func TestCircularDependency_ValidateAlphaChain_WithCycle(t *testing.T) {
	detector := NewCircularDependencyDetector()
	storage := NewMemoryStorage()
	// Créer une chaîne avec cycle
	nodes := []*AlphaNode{
		{
			BaseNode: BaseNode{
				ID:      "alpha_1",
				Type:    "alpha",
				Storage: storage,
			},
			ResultName:   "temp_1",
			IsAtomic:     true,
			Dependencies: []string{"temp_2"}, // Dépend de temp_2
		},
		{
			BaseNode: BaseNode{
				ID:      "alpha_2",
				Type:    "alpha",
				Storage: storage,
			},
			ResultName:   "temp_2",
			IsAtomic:     true,
			Dependencies: []string{"temp_1"}, // Dépend de temp_1 -> cycle!
		},
	}
	report := detector.ValidateAlphaChain(nodes)
	if report.Valid {
		t.Error("Expected invalid chain due to cycle")
	}
	if !report.HasCircularDeps {
		t.Error("Expected circular dependencies to be detected")
	}
	t.Logf("Cycle detection report: Valid=%v, Error=%s, Cycle=%v",
		report.Valid, report.ErrorMessage, report.CyclePath)
}

// TestCircularDependency_AddNodeWithMetadata teste l'ajout de métadonnées
func TestCircularDependency_AddNodeWithMetadata(t *testing.T) {
	detector := NewCircularDependencyDetector()
	metadata := &NodeMetadata{
		ResultName:  "temp_1",
		Depth:       0,
		Expression:  "a * b",
		IsAtomic:    true,
		Descendants: []string{},
	}
	detector.AddNodeWithMetadata("temp_1", []string{}, metadata)
	if detector.metadata["temp_1"] == nil {
		t.Error("Expected metadata to be stored")
	}
	if detector.metadata["temp_1"].Expression != "a * b" {
		t.Errorf("Expected expression 'a * b', got '%s'", detector.metadata["temp_1"].Expression)
	}
}

// TestCircularDependency_RealWorldScenario teste un scénario réel d'expression arithmétique
func TestCircularDependency_RealWorldScenario(t *testing.T) {
	detector := NewCircularDependencyDetector()
	// Expression: (c.qte * 23 - 10 + c.remise * 43) > 0
	// Décomposition:
	// temp_1 = c.qte * 23
	// temp_2 = temp_1 - 10
	// temp_3 = c.remise * 43
	// temp_4 = temp_2 + temp_3
	// temp_5 = temp_4 > 0
	detector.AddNode("temp_1", []string{})                   // c.qte * 23
	detector.AddNode("temp_2", []string{"temp_1"})           // temp_1 - 10
	detector.AddNode("temp_3", []string{})                   // c.remise * 43
	detector.AddNode("temp_4", []string{"temp_2", "temp_3"}) // temp_2 + temp_3
	detector.AddNode("temp_5", []string{"temp_4"})           // temp_4 > 0
	report := detector.Validate()
	if !report.Valid {
		t.Errorf("Expected valid decomposition, got error: %s", report.ErrorMessage)
	}
	if report.HasCircularDeps {
		t.Error("Expected no circular dependencies")
	}
	t.Logf("Real-world scenario validation: Valid=%v, MaxDepth=%d, TotalNodes=%d",
		report.Valid, report.MaxDepth, report.TotalNodes)
	// Vérifier le tri topologique
	sorted, err := detector.GetTopologicalSort()
	if err != nil {
		t.Fatalf("Expected successful topological sort, got error: %v", err)
	}
	t.Logf("Execution order (topological sort): %v", sorted)
}

// TestCircularDependency_String teste la représentation textuelle
func TestCircularDependency_String(t *testing.T) {
	detector := NewCircularDependencyDetector()
	detector.AddNode("temp_1", []string{})
	detector.AddNode("temp_2", []string{"temp_1"})
	str := detector.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	t.Logf("String representation:\n%s", str)
}

// BenchmarkCircularDependency_DetectCycles benchmark la détection de cycles
func BenchmarkCircularDependency_DetectCycles(b *testing.B) {
	detector := NewCircularDependencyDetector()
	// Créer une chaîne de 100 nœuds
	for i := 0; i < 100; i++ {
		if i == 0 {
			detector.AddNode("temp_0", []string{})
		} else {
			detector.AddNode("temp_"+string(rune('0'+i)), []string{"temp_" + string(rune('0'+i-1))})
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.DetectCycles()
	}
}

// BenchmarkCircularDependency_Validate benchmark la validation complète
func BenchmarkCircularDependency_Validate(b *testing.B) {
	detector := NewCircularDependencyDetector()
	// Créer un graphe complexe
	for i := 0; i < 50; i++ {
		if i == 0 {
			detector.AddNode("temp_0", []string{})
		} else if i%2 == 0 {
			detector.AddNode("temp_"+string(rune('0'+i)), []string{"temp_" + string(rune('0'+i-1))})
		} else {
			detector.AddNode("temp_"+string(rune('0'+i)), []string{"temp_" + string(rune('0'+i-2))})
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector.Validate()
	}
}
