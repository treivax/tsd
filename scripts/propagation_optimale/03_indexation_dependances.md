# 🗂️ Prompt 03 - Indexation des Dépendances

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Implémenter le système d'indexation des dépendances : `DependencyIndex` qui permet de déterminer rapidement quels nœuds RETE sont sensibles à chaque champ de chaque type de fait.

Cette indexation est le cœur du système de propagation delta : elle permet de ne propager que vers les nœuds réellement impactés par un changement.

**⚠️ IMPORTANT** : Ce prompt génère du code. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompt 01 validé** : Conception disponible
- [x] **Prompt 02 validé** : Modèle de données delta implémenté
- [x] **Tests passent** : `go test ./rete/delta/... -v` (100% success)
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `REPORTS/metadata_noeuds.md`
  - `REPORTS/ast_conditions_mapping.md`

---

## 📂 Fichiers à Créer

Ajouter au package `rete/delta` :

```
rete/delta/
├── dependency_index.go           # Structure DependencyIndex
├── dependency_index_test.go      # Tests unitaires
├── index_builder.go              # Construction de l'index
├── index_builder_test.go         # Tests construction
├── field_extractor.go            # Extraction champs depuis AST
├── field_extractor_test.go       # Tests extraction
└── index_metrics.go              # Métriques d'indexation
```

---

## 🔧 Tâche 1 : Structure DependencyIndex

### Fichier : `rete/delta/dependency_index.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"sync"
	"time"
)

// NodeReference représente une référence à un nœud RETE.
//
// Cette structure stocke les informations nécessaires pour identifier
// et accéder à un nœud sans stocker directement le pointeur (évite cycles).
type NodeReference struct {
	// NodeID est l'identifiant unique du nœud
	NodeID string
	
	// NodeType indique le type de nœud ("alpha", "beta", "terminal")
	NodeType string
	
	// FactType est le type de fait concerné (ex: "Product")
	FactType string
	
	// Fields est la liste des champs utilisés par ce nœud
	Fields []string
}

// String retourne une représentation string de la référence
func (nr NodeReference) String() string {
	return fmt.Sprintf("%s[%s](%s)", nr.NodeType, nr.NodeID, nr.FactType)
}

// DependencyIndex est un index inversé permettant de trouver rapidement
// tous les nœuds RETE sensibles à un champ donné d'un type donné.
//
// Structure : factType → field → [nodeRefs]
//
// Exemple :
//   index.Get("Product", "price") → [alpha1, alpha2, beta3, terminal5]
//
// Thread-safety : toutes les opérations sont protégées par mutex.
type DependencyIndex struct {
	// alphaIndex indexe les nœuds alpha par (factType, field)
	// Structure : factType → field → [nodeIDs]
	alphaIndex map[string]map[string][]string
	
	// betaIndex indexe les nœuds beta par (factType, field)
	betaIndex map[string]map[string][]string
	
	// terminalIndex indexe les nœuds terminaux par (factType, field)
	terminalIndex map[string]map[string][]string
	
	// nodeReferences stocke les détails de chaque nœud indexé
	// Structure : nodeID → NodeReference
	nodeReferences map[string]NodeReference
	
	// metadata stocke des informations sur l'index
	builtAt    time.Time
	nodeCount  int
	fieldCount int
	
	// mutex protège l'accès concurrent
	mutex sync.RWMutex
}

// NewDependencyIndex crée un nouvel index de dépendances vide.
func NewDependencyIndex() *DependencyIndex {
	return &DependencyIndex{
		alphaIndex:     make(map[string]map[string][]string),
		betaIndex:      make(map[string]map[string][]string),
		terminalIndex:  make(map[string]map[string][]string),
		nodeReferences: make(map[string]NodeReference),
		builtAt:        time.Now(),
	}
}

// AddAlphaNode enregistre un nœud alpha pour un ensemble de champs.
//
// Paramètres :
//   - nodeID : identifiant unique du nœud
//   - factType : type de fait concerné
//   - fields : liste des champs testés par ce nœud
func (idx *DependencyIndex) AddAlphaNode(nodeID, factType string, fields []string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()
	
	idx.addNodeToIndex(idx.alphaIndex, nodeID, factType, fields)
	idx.nodeReferences[nodeID] = NodeReference{
		NodeID:   nodeID,
		NodeType: "alpha",
		FactType: factType,
		Fields:   fields,
	}
	idx.nodeCount++
}

// AddBetaNode enregistre un nœud beta pour un ensemble de champs.
//
// Paramètres :
//   - nodeID : identifiant unique du nœud
//   - factType : type de fait concerné
//   - fields : liste des champs utilisés dans les jointures
func (idx *DependencyIndex) AddBetaNode(nodeID, factType string, fields []string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()
	
	idx.addNodeToIndex(idx.betaIndex, nodeID, factType, fields)
	idx.nodeReferences[nodeID] = NodeReference{
		NodeID:   nodeID,
		NodeType: "beta",
		FactType: factType,
		Fields:   fields,
	}
	idx.nodeCount++
}

// AddTerminalNode enregistre un nœud terminal pour un ensemble de champs.
//
// Paramètres :
//   - nodeID : identifiant unique du nœud
//   - factType : type de fait concerné
//   - fields : liste des champs utilisés dans les actions
func (idx *DependencyIndex) AddTerminalNode(nodeID, factType string, fields []string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()
	
	idx.addNodeToIndex(idx.terminalIndex, nodeID, factType, fields)
	idx.nodeReferences[nodeID] = NodeReference{
		NodeID:   nodeID,
		NodeType: "terminal",
		FactType: factType,
		Fields:   fields,
	}
	idx.nodeCount++
}

// addNodeToIndex est une fonction helper privée pour ajouter un nœud à un index.
// ATTENTION : doit être appelée avec mutex déjà acquis.
func (idx *DependencyIndex) addNodeToIndex(
	index map[string]map[string][]string,
	nodeID, factType string,
	fields []string,
) {
	// Initialiser la map pour ce factType si nécessaire
	if index[factType] == nil {
		index[factType] = make(map[string][]string)
	}
	
	// Ajouter le nœud pour chaque champ
	for _, field := range fields {
		// Vérifier si le nœud n'est pas déjà indexé pour ce champ
		nodes := index[factType][field]
		alreadyIndexed := false
		for _, existingNodeID := range nodes {
			if existingNodeID == nodeID {
				alreadyIndexed = true
				break
			}
		}
		
		if !alreadyIndexed {
			index[factType][field] = append(index[factType][field], nodeID)
			idx.fieldCount++
		}
	}
}

// GetAffectedNodes retourne tous les nœuds affectés par un changement de champ.
//
// Paramètres :
//   - factType : type de fait
//   - field : nom du champ modifié
//
// Retourne la liste des NodeReferences des nœuds sensibles à ce champ.
func (idx *DependencyIndex) GetAffectedNodes(factType, field string) []NodeReference {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()
	
	affectedNodeIDs := make(map[string]bool)
	
	// Collecter les nœuds alpha
	if alphaFields := idx.alphaIndex[factType]; alphaFields != nil {
		if nodes := alphaFields[field]; nodes != nil {
			for _, nodeID := range nodes {
				affectedNodeIDs[nodeID] = true
			}
		}
	}
	
	// Collecter les nœuds beta
	if betaFields := idx.betaIndex[factType]; betaFields != nil {
		if nodes := betaFields[field]; nodes != nil {
			for _, nodeID := range nodes {
				affectedNodeIDs[nodeID] = true
			}
		}
	}
	
	// Collecter les nœuds terminaux
	if terminalFields := idx.terminalIndex[factType]; terminalFields != nil {
		if nodes := terminalFields[field]; nodes != nil {
			for _, nodeID := range nodes {
				affectedNodeIDs[nodeID] = true
			}
		}
	}
	
	// Convertir en NodeReferences
	result := make([]NodeReference, 0, len(affectedNodeIDs))
	for nodeID := range affectedNodeIDs {
		if ref, exists := idx.nodeReferences[nodeID]; exists {
			result = append(result, ref)
		}
	}
	
	return result
}

// GetAffectedNodesForDelta retourne tous les nœuds affectés par un FactDelta.
//
// Paramètres :
//   - delta : le FactDelta contenant les champs modifiés
//
// Retourne la liste des NodeReferences des nœuds affectés (dédupliqués).
func (idx *DependencyIndex) GetAffectedNodesForDelta(delta *FactDelta) []NodeReference {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()
	
	affectedNodeIDs := make(map[string]bool)
	
	// Pour chaque champ modifié
	for fieldName := range delta.Fields {
		// Alpha nodes
		if alphaFields := idx.alphaIndex[delta.FactType]; alphaFields != nil {
			if nodes := alphaFields[fieldName]; nodes != nil {
				for _, nodeID := range nodes {
					affectedNodeIDs[nodeID] = true
				}
			}
		}
		
		// Beta nodes
		if betaFields := idx.betaIndex[delta.FactType]; betaFields != nil {
			if nodes := betaFields[fieldName]; nodes != nil {
				for _, nodeID := range nodes {
					affectedNodeIDs[nodeID] = true
				}
			}
		}
		
		// Terminal nodes
		if terminalFields := idx.terminalIndex[delta.FactType]; terminalFields != nil {
			if nodes := terminalFields[fieldName]; nodes != nil {
				for _, nodeID := range nodes {
					affectedNodeIDs[nodeID] = true
				}
			}
		}
	}
	
	// Convertir en NodeReferences
	result := make([]NodeReference, 0, len(affectedNodeIDs))
	for nodeID := range affectedNodeIDs {
		if ref, exists := idx.nodeReferences[nodeID]; exists {
			result = append(result, ref)
		}
	}
	
	return result
}

// Clear vide complètement l'index.
func (idx *DependencyIndex) Clear() {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()
	
	idx.alphaIndex = make(map[string]map[string][]string)
	idx.betaIndex = make(map[string]map[string][]string)
	idx.terminalIndex = make(map[string]map[string][]string)
	idx.nodeReferences = make(map[string]NodeReference)
	idx.nodeCount = 0
	idx.fieldCount = 0
	idx.builtAt = time.Now()
}

// Stats retourne des statistiques sur l'index.
type IndexStats struct {
	NodeCount       int
	FieldCount      int
	AlphaNodeCount  int
	BetaNodeCount   int
	TerminalCount   int
	FactTypes       []string
	BuiltAt         time.Time
	MemoryEstimate  int64 // Estimation mémoire en bytes
}

// GetStats retourne les statistiques de l'index.
func (idx *DependencyIndex) GetStats() IndexStats {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()
	
	stats := IndexStats{
		NodeCount:  idx.nodeCount,
		FieldCount: idx.fieldCount,
		BuiltAt:    idx.builtAt,
	}
	
	// Compter par type de nœud
	for _, ref := range idx.nodeReferences {
		switch ref.NodeType {
		case "alpha":
			stats.AlphaNodeCount++
		case "beta":
			stats.BetaNodeCount++
		case "terminal":
			stats.TerminalCount++
		}
	}
	
	// Collecter les types de faits
	factTypeSet := make(map[string]bool)
	for factType := range idx.alphaIndex {
		factTypeSet[factType] = true
	}
	for factType := range idx.betaIndex {
		factTypeSet[factType] = true
	}
	for factType := range idx.terminalIndex {
		factTypeSet[factType] = true
	}
	
	stats.FactTypes = make([]string, 0, len(factTypeSet))
	for factType := range factTypeSet {
		stats.FactTypes = append(stats.FactTypes, factType)
	}
	
	// Estimation mémoire approximative
	// nodeID (string) ~ 50 bytes, ref ~ 100 bytes, slice overhead ~ 24 bytes
	stats.MemoryEstimate = int64(idx.nodeCount * 150)
	stats.MemoryEstimate += int64(idx.fieldCount * 74) // index entries
	
	return stats
}

// String retourne une représentation string de l'index.
func (idx *DependencyIndex) String() string {
	stats := idx.GetStats()
	return fmt.Sprintf(
		"DependencyIndex[nodes=%d, fields=%d, alpha=%d, beta=%d, terminal=%d, types=%d]",
		stats.NodeCount, stats.FieldCount,
		stats.AlphaNodeCount, stats.BetaNodeCount, stats.TerminalCount,
		len(stats.FactTypes),
	)
}
```

### Tests : `rete/delta/dependency_index_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestNewDependencyIndex(t *testing.T) {
	idx := NewDependencyIndex()
	
	if idx == nil {
		t.Fatal("NewDependencyIndex returned nil")
	}
	
	if idx.alphaIndex == nil || idx.betaIndex == nil || idx.terminalIndex == nil {
		t.Error("Indexes not initialized")
	}
	
	if idx.nodeReferences == nil {
		t.Error("nodeReferences not initialized")
	}
	
	// Vérifier timestamp
	if time.Since(idx.builtAt) > time.Second {
		t.Error("builtAt timestamp incorrect")
	}
}

func TestDependencyIndex_AddAlphaNode(t *testing.T) {
	idx := NewDependencyIndex()
	
	idx.AddAlphaNode("alpha1", "Product", []string{"price", "status"})
	
	// Vérifier que le nœud est indexé
	nodes := idx.GetAffectedNodes("Product", "price")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}
	
	if nodes[0].NodeID != "alpha1" {
		t.Errorf("Expected alpha1, got %s", nodes[0].NodeID)
	}
	
	if nodes[0].NodeType != "alpha" {
		t.Errorf("Expected alpha type, got %s", nodes[0].NodeType)
	}
	
	// Vérifier pour l'autre champ
	nodes = idx.GetAffectedNodes("Product", "status")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node for status, got %d", len(nodes))
	}
}

func TestDependencyIndex_AddBetaNode(t *testing.T) {
	idx := NewDependencyIndex()
	
	idx.AddBetaNode("beta1", "Order", []string{"customer_id"})
	
	nodes := idx.GetAffectedNodes("Order", "customer_id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}
	
	if nodes[0].NodeType != "beta" {
		t.Errorf("Expected beta type, got %s", nodes[0].NodeType)
	}
}

func TestDependencyIndex_AddTerminalNode(t *testing.T) {
	idx := NewDependencyIndex()
	
	idx.AddTerminalNode("term1", "Alert", []string{"severity", "message"})
	
	nodes := idx.GetAffectedNodes("Alert", "severity")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(nodes))
	}
	
	if nodes[0].NodeType != "terminal" {
		t.Errorf("Expected terminal type, got %s", nodes[0].NodeType)
	}
}

func TestDependencyIndex_MultipleNodesPerField(t *testing.T) {
	idx := NewDependencyIndex()
	
	// Plusieurs nœuds sensibles au même champ
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha2", "Product", []string{"price", "category"})
	idx.AddBetaNode("beta1", "Product", []string{"price"})
	
	nodes := idx.GetAffectedNodes("Product", "price")
	
	if len(nodes) != 3 {
		t.Fatalf("Expected 3 nodes for price, got %d", len(nodes))
	}
	
	// Vérifier que tous les nœuds sont présents
	nodeIDs := make(map[string]bool)
	for _, node := range nodes {
		nodeIDs[node.NodeID] = true
	}
	
	if !nodeIDs["alpha1"] || !nodeIDs["alpha2"] || !nodeIDs["beta1"] {
		t.Errorf("Missing expected nodes: %v", nodeIDs)
	}
}

func TestDependencyIndex_NoduplicateIndexing(t *testing.T) {
	idx := NewDependencyIndex()
	
	// Ajouter le même nœud deux fois
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	
	nodes := idx.GetAffectedNodes("Product", "price")
	
	// Le nœud ne devrait être indexé qu'une fois
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node (no duplicate), got %d", len(nodes))
	}
}

func TestDependencyIndex_GetAffectedNodesForDelta(t *testing.T) {
	idx := NewDependencyIndex()
	
	// Setup : indexer plusieurs nœuds
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddAlphaNode("alpha2", "Product", []string{"status"})
	idx.AddBetaNode("beta1", "Product", []string{"price", "category"})
	idx.AddTerminalNode("term1", "Product", []string{"status"})
	
	// Créer un delta modifiant price et status
	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100.0, 150.0)
	delta.AddFieldChange("status", "active", "inactive")
	
	nodes := idx.GetAffectedNodesForDelta(delta)
	
	// Devrait retourner : alpha1 (price), alpha2 (status), beta1 (price), term1 (status)
	// = 4 nœuds uniques
	if len(nodes) != 4 {
		t.Fatalf("Expected 4 affected nodes, got %d", len(nodes))
	}
	
	// Vérifier présence de tous
	nodeIDs := make(map[string]bool)
	for _, node := range nodes {
		nodeIDs[node.NodeID] = true
	}
	
	expected := []string{"alpha1", "alpha2", "beta1", "term1"}
	for _, expectedID := range expected {
		if !nodeIDs[expectedID] {
			t.Errorf("Missing expected node: %s", expectedID)
		}
	}
}

func TestDependencyIndex_DifferentFactTypes(t *testing.T) {
	idx := NewDependencyIndex()
	
	// Même nom de champ, types différents
	idx.AddAlphaNode("alpha_product", "Product", []string{"id"})
	idx.AddAlphaNode("alpha_order", "Order", []string{"id"})
	
	// Requête pour Product.id
	nodes := idx.GetAffectedNodes("Product", "id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node for Product.id, got %d", len(nodes))
	}
	if nodes[0].NodeID != "alpha_product" {
		t.Errorf("Expected alpha_product, got %s", nodes[0].NodeID)
	}
	
	// Requête pour Order.id
	nodes = idx.GetAffectedNodes("Order", "id")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node for Order.id, got %d", len(nodes))
	}
	if nodes[0].NodeID != "alpha_order" {
		t.Errorf("Expected alpha_order, got %s", nodes[0].NodeID)
	}
}

func TestDependencyIndex_Clear(t *testing.T) {
	idx := NewDependencyIndex()
	
	// Ajouter quelques nœuds
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	idx.AddBetaNode("beta1", "Order", []string{"total"})
	
	// Vérifier qu'ils sont là
	if len(idx.GetAffectedNodes("Product", "price")) == 0 {
		t.Fatal("Nodes should exist before Clear")
	}
	
	// Clear
	idx.Clear()
	
	// Vérifier que tout est vide
	if len(idx.GetAffectedNodes("Product", "price")) != 0 {
		t.Error("Index should be empty after Clear")
	}
	
	stats := idx.GetStats()
	if stats.NodeCount != 0 || stats.FieldCount != 0 {
		t.Errorf("Stats should be zero after Clear: %+v", stats)
	}
}

func TestDependencyIndex_GetStats(t *testing.T) {
	idx := NewDependencyIndex()
	
	idx.AddAlphaNode("alpha1", "Product", []string{"price", "status"})
	idx.AddAlphaNode("alpha2", "Product", []string{"category"})
	idx.AddBetaNode("beta1", "Order", []string{"customer_id"})
	idx.AddTerminalNode("term1", "Alert", []string{"severity"})
	
	stats := idx.GetStats()
	
	if stats.NodeCount != 4 {
		t.Errorf("Expected 4 nodes, got %d", stats.NodeCount)
	}
	
	if stats.AlphaNodeCount != 2 {
		t.Errorf("Expected 2 alpha nodes, got %d", stats.AlphaNodeCount)
	}
	
	if stats.BetaNodeCount != 1 {
		t.Errorf("Expected 1 beta node, got %d", stats.BetaNodeCount)
	}
	
	if stats.TerminalCount != 1 {
		t.Errorf("Expected 1 terminal node, got %d", stats.TerminalCount)
	}
	
	if len(stats.FactTypes) != 3 {
		t.Errorf("Expected 3 fact types, got %d", len(stats.FactTypes))
	}
	
	if stats.MemoryEstimate <= 0 {
		t.Error("MemoryEstimate should be positive")
	}
}

func TestDependencyIndex_String(t *testing.T) {
	idx := NewDependencyIndex()
	idx.AddAlphaNode("alpha1", "Product", []string{"price"})
	
	str := idx.String()
	if str == "" {
		t.Error("String() should not be empty")
	}
	
	// Vérifier que contient des informations clés
	if len(str) < 20 {
		t.Errorf("String() too short: %s", str)
	}
}

func TestNodeReference_String(t *testing.T) {
	ref := NodeReference{
		NodeID:   "alpha1",
		NodeType: "alpha",
		FactType: "Product",
		Fields:   []string{"price"},
	}
	
	str := ref.String()
	expectedMatch := "alpha[alpha1](Product)"
	if str != expectedMatch {
		t.Errorf("String() = %s, want %s", str, expectedMatch)
	}
}

func TestDependencyIndex_ConcurrentAccess(t *testing.T) {
	idx := NewDependencyIndex()
	
	// Test accès concurrent (détection race conditions)
	done := make(chan bool, 2)
	
	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			idx.AddAlphaNode("alpha1", "Product", []string{"price"})
		}
		done <- true
	}()
	
	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			_ = idx.GetAffectedNodes("Product", "price")
		}
		done <- true
	}()
	
	// Attendre fin
	<-done
	<-done
	
	// Vérifier état final cohérent
	nodes := idx.GetAffectedNodes("Product", "price")
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node after concurrent access, got %d", len(nodes))
	}
}

func BenchmarkDependencyIndex_AddAlphaNode(b *testing.B) {
	idx := NewDependencyIndex()
	fields := []string{"field1", "field2", "field3"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nodeID := "alpha" + string(rune(i))
		idx.AddAlphaNode(nodeID, "TestType", fields)
	}
}

func BenchmarkDependencyIndex_GetAffectedNodes(b *testing.B) {
	idx := NewDependencyIndex()
	
	// Setup : 100 nœuds indexés
	for i := 0; i < 100; i++ {
		nodeID := "alpha" + string(rune(i))
		idx.AddAlphaNode(nodeID, "Product", []string{"price"})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.GetAffectedNodes("Product", "price")
	}
}

func BenchmarkDependencyIndex_GetAffectedNodesForDelta(b *testing.B) {
	idx := NewDependencyIndex()
	
	// Setup
	for i := 0; i < 50; i++ {
		idx.AddAlphaNode("alpha"+string(rune(i)), "Product", []string{"price", "status"})
	}
	
	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100, 150)
	delta.AddFieldChange("status", "active", "inactive")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.GetAffectedNodesForDelta(delta)
	}
}
```

---

## 🔧 Tâche 2 : Extraction de Champs depuis AST

### Fichier : `rete/delta/field_extractor.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
)

// FieldExtractor extrait les noms de champs depuis diverses structures AST.
//
// Cette interface abstraite permet de supporter différents formats de conditions
// (alpha, beta, terminal) de manière uniforme.
type FieldExtractor interface {
	// ExtractFields extrait les champs depuis une condition/expression
	ExtractFields(condition interface{}) ([]string, error)
}

// AlphaConditionExtractor extrait les champs depuis les conditions de nœuds alpha.
type AlphaConditionExtractor struct{}

// ExtractFields extrait les champs depuis une condition alpha.
//
// Les conditions alpha sont typiquement des expressions binaires ou comparaisons
// sur des champs d'un fait.
//
// Exemples :
//   - "product.price > 100" → ["price"]
//   - "order.status == 'active' && order.total > 50" → ["status", "total"]
//
// Paramètres :
//   - condition : structure de condition (map[string]interface{} depuis parser)
//
// Retourne la liste des noms de champs (dédupliqués).
func (ace *AlphaConditionExtractor) ExtractFields(condition interface{}) ([]string, error) {
	fields := make(map[string]bool)
	
	if err := extractFieldsRecursive(condition, fields); err != nil {
		return nil, err
	}
	
	// Convertir map en slice
	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}
	
	return result, nil
}

// BetaConditionExtractor extrait les champs depuis les conditions de jointure beta.
type BetaConditionExtractor struct{}

// ExtractFields extrait les champs depuis une condition de jointure beta.
//
// Les conditions beta sont des comparaisons entre champs de différents faits.
//
// Exemple :
//   - "order.customer_id == customer.id" → {"order": ["customer_id"], "customer": ["id"]}
//
// Note : Cette version retourne une liste plate de tous les champs.
// Une version future pourrait retourner une map factType → fields.
func (bce *BetaConditionExtractor) ExtractFields(condition interface{}) ([]string, error) {
	fields := make(map[string]bool)
	
	if err := extractFieldsRecursive(condition, fields); err != nil {
		return nil, err
	}
	
	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}
	
	return result, nil
}

// ActionFieldExtractor extrait les champs depuis les actions des nœuds terminaux.
type ActionFieldExtractor struct{}

// ExtractFields extrait les champs depuis une action.
//
// Exemples :
//   - Update(product, {price: product.price * 1.1}) → ["price"]
//   - Insert(Alert, {message: order.id}) → ["id"] (du contexte)
//
// Cette fonction examine les arguments et modifications pour extraire
// tous les champs référencés.
func (afe *ActionFieldExtractor) ExtractFields(action interface{}) ([]string, error) {
	fields := make(map[string]bool)
	
	if err := extractFieldsRecursive(action, fields); err != nil {
		return nil, err
	}
	
	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}
	
	return result, nil
}

// extractFieldsRecursive est une fonction récursive privée qui parcourt
// une structure AST et collecte tous les champs référencés.
//
// Cette fonction est générique et fonctionne pour alpha, beta et terminal nodes.
func extractFieldsRecursive(node interface{}, fields map[string]bool) error {
	if node == nil {
		return nil
	}
	
	// Cas 1: Map (structure du parser)
	if nodeMap, ok := node.(map[string]interface{}); ok {
		nodeType, hasType := nodeMap["type"].(string)
		
		if hasType {
			switch nodeType {
			case "fieldAccess":
				// Accès à un champ : variable.field
				if field, ok := nodeMap["field"].(string); ok {
					fields[field] = true
				}
				
			case "binaryOp":
				// Opération binaire : récurser sur left et right
				if err := extractFieldsRecursive(nodeMap["left"], fields); err != nil {
					return err
				}
				if err := extractFieldsRecursive(nodeMap["right"], fields); err != nil {
					return err
				}
				
			case "comparison":
				// Comparaison : extraire des deux côtés
				if err := extractFieldsRecursive(nodeMap["left"], fields); err != nil {
					return err
				}
				if err := extractFieldsRecursive(nodeMap["right"], fields); err != nil {
					return err
				}
				
			case "updateWithModifications":
				// Action Update : extraire champs modifiés
				if modifications, ok := nodeMap["modifications"].(map[string]interface{}); ok {
					for fieldName := range modifications {
						fields[fieldName] = true
					}
				}
				
			case "factCreation":
				// Insert : extraire champs du fait créé
				if factFields, ok := nodeMap["fields"].(map[string]interface{}); ok {
					for fieldName := range factFields {
						fields[fieldName] = true
					}
				}
			}
		}
		
		// Récurser sur toutes les valeurs de la map
		for _, value := range nodeMap {
			if err := extractFieldsRecursive(value, fields); err != nil {
				return err
			}
		}
	}
	
	// Cas 2: Slice (liste d'expressions)
	if nodeSlice, ok := node.([]interface{}); ok {
		for _, item := range nodeSlice {
			if err := extractFieldsRecursive(item, fields); err != nil {
				return err
			}
		}
	}
	
	// Cas 3: Types primitifs → ignorer (pas de champs)
	
	return nil
}

// ExtractFieldsFromAlphaCondition est une fonction helper pour extraire
// les champs depuis une condition alpha.
func ExtractFieldsFromAlphaCondition(condition interface{}) ([]string, error) {
	extractor := &AlphaConditionExtractor{}
	return extractor.ExtractFields(condition)
}

// ExtractFieldsFromBetaCondition est une fonction helper pour extraire
// les champs depuis une condition beta.
func ExtractFieldsFromBetaCondition(condition interface{}) ([]string, error) {
	extractor := &BetaConditionExtractor{}
	return extractor.ExtractFields(condition)
}

// ExtractFieldsFromAction est une fonction helper pour extraire
// les champs depuis une action.
func ExtractFieldsFromAction(action interface{}) ([]string, error) {
	extractor := &ActionFieldExtractor{}
	return extractor.ExtractFields(action)
}
```

### Tests : `rete/delta/field_extractor_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"sort"
	"testing"
)

func TestAlphaConditionExtractor_SimpleFieldAccess(t *testing.T) {
	// Condition : product.price
	condition := map[string]interface{}{
		"type":  "fieldAccess",
		"field": "price",
	}
	
	extractor := &AlphaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}
	
	if fields[0] != "price" {
		t.Errorf("Expected 'price', got '%s'", fields[0])
	}
}

func TestAlphaConditionExtractor_Comparison(t *testing.T) {
	// Condition : product.price > 100
	condition := map[string]interface{}{
		"type": "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100,
		},
	}
	
	extractor := &AlphaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(fields))
	}
	
	if fields[0] != "price" {
		t.Errorf("Expected 'price', got '%s'", fields[0])
	}
}

func TestAlphaConditionExtractor_BinaryOp(t *testing.T) {
	// Condition : product.price > 100 && product.status == "active"
	condition := map[string]interface{}{
		"type": "binaryOp",
		"operator": "&&",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "price",
			},
			"right": 100,
		},
		"right": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "status",
			},
			"right": "active",
		},
	}
	
	extractor := &AlphaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
	
	// Trier pour comparaison stable
	sort.Strings(fields)
	
	if fields[0] != "price" || fields[1] != "status" {
		t.Errorf("Expected ['price', 'status'], got %v", fields)
	}
}

func TestBetaConditionExtractor_JoinCondition(t *testing.T) {
	// Condition : order.customer_id == customer.id
	condition := map[string]interface{}{
		"type": "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":     "fieldAccess",
			"variable": "order",
			"field":    "customer_id",
		},
		"right": map[string]interface{}{
			"type":     "fieldAccess",
			"variable": "customer",
			"field":    "id",
		},
	}
	
	extractor := &BetaConditionExtractor{}
	fields, err := extractor.ExtractFields(condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
	
	// Vérifier que customer_id et id sont présents
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}
	
	if !fieldMap["customer_id"] || !fieldMap["id"] {
		t.Errorf("Expected ['customer_id', 'id'], got %v", fields)
	}
}

func TestActionFieldExtractor_UpdateAction(t *testing.T) {
	// Action : Update(product, {price: newValue})
	action := map[string]interface{}{
		"type": "updateWithModifications",
		"variable": "product",
		"modifications": map[string]interface{}{
			"price":  150,
			"status": "updated",
		},
	}
	
	extractor := &ActionFieldExtractor{}
	fields, err := extractor.ExtractFields(action)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
	
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}
	
	if !fieldMap["price"] || !fieldMap["status"] {
		t.Errorf("Expected ['price', 'status'], got %v", fields)
	}
}

func TestActionFieldExtractor_InsertAction(t *testing.T) {
	// Action : Insert(Alert, {severity: "high", message: "test"})
	action := map[string]interface{}{
		"type": "factCreation",
		"factType": "Alert",
		"fields": map[string]interface{}{
			"severity": "high",
			"message":  "test",
		},
	}
	
	extractor := &ActionFieldExtractor{}
	fields, err := extractor.ExtractFields(action)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
	
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}
	
	if !fieldMap["severity"] || !fieldMap["message"] {
		t.Errorf("Expected ['severity', 'message'], got %v", fields)
	}
}

func TestExtractFieldsRecursive_EmptyNode(t *testing.T) {
	fields := make(map[string]bool)
	err := extractFieldsRecursive(nil, fields)
	
	if err != nil {
		t.Errorf("Expected no error for nil node, got %v", err)
	}
	
	if len(fields) != 0 {
		t.Error("Expected no fields for nil node")
	}
}

func TestExtractFieldsRecursive_NestedStructure(t *testing.T) {
	// Structure complexe imbriquée
	node := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type": "binaryOp",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "field1",
			},
			"right": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "field2",
			},
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "field3",
		},
	}
	
	fields := make(map[string]bool)
	err := extractFieldsRecursive(node, fields)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 3 {
		t.Fatalf("Expected 3 fields, got %d", len(fields))
	}
	
	for i := 1; i <= 3; i++ {
		fieldName := "field" + string(rune('0'+i))
		if !fields[fieldName] {
			t.Errorf("Expected field '%s' to be extracted", fieldName)
		}
	}
}

func TestExtractFieldsRecursive_SliceInput(t *testing.T) {
	// Liste d'accès à des champs
	nodes := []interface{}{
		map[string]interface{}{
			"type":  "fieldAccess",
			"field": "field1",
		},
		map[string]interface{}{
			"type":  "fieldAccess",
			"field": "field2",
		},
	}
	
	fields := make(map[string]bool)
	err := extractFieldsRecursive(nodes, fields)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
}

func TestExtractFieldsFromAlphaCondition_HelperFunction(t *testing.T) {
	condition := map[string]interface{}{
		"type":  "fieldAccess",
		"field": "testField",
	}
	
	fields, err := ExtractFieldsFromAlphaCondition(condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 1 || fields[0] != "testField" {
		t.Errorf("Expected ['testField'], got %v", fields)
	}
}

func TestExtractFieldsFromBetaCondition_HelperFunction(t *testing.T) {
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "leftField",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "rightField",
		},
	}
	
	fields, err := ExtractFieldsFromBetaCondition(condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(fields))
	}
}

func TestExtractFieldsFromAction_HelperFunction(t *testing.T) {
	action := map[string]interface{}{
		"type": "updateWithModifications",
		"modifications": map[string]interface{}{
			"field1": "value1",
		},
	}
	
	fields, err := ExtractFieldsFromAction(action)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if len(fields) != 1 || fields[0] != "field1" {
		t.Errorf("Expected ['field1'], got %v", fields)
	}
}

func BenchmarkExtractFieldsFromAlphaCondition(b *testing.B) {
	condition := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "status",
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ExtractFieldsFromAlphaCondition(condition)
	}
}
```

---

## 🔧 Tâche 3 : Construction de l'Index depuis le Réseau RETE

### Fichier : `rete/delta/index_builder.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
)

// IndexBuilder construit un DependencyIndex depuis un réseau RETE complet.
//
// Cette structure orchestre l'extraction de champs depuis tous les nœuds
// et la construction de l'index de dépendances.
type IndexBuilder struct {
	alphaExtractor    *AlphaConditionExtractor
	betaExtractor     *BetaConditionExtractor
	actionExtractor   *ActionFieldExtractor
	enableDiagnostics bool
	diagnostics       *BuildDiagnostics
}

// BuildDiagnostics contient des informations de diagnostic sur la construction.
type BuildDiagnostics struct {
	NodesProcessed    int
	NodesSkipped      int
	FieldsExtracted   int
	Errors            []string
	Warnings          []string
}

// NewIndexBuilder crée un nouveau constructeur d'index.
func NewIndexBuilder() *IndexBuilder {
	return &IndexBuilder{
		alphaExtractor:    &AlphaConditionExtractor{},
		betaExtractor:     &BetaConditionExtractor{},
		actionExtractor:   &ActionFieldExtractor{},
		enableDiagnostics: false,
		diagnostics:       &BuildDiagnostics{},
	}
}

// EnableDiagnostics active la collecte d'informations de diagnostic.
func (ib *IndexBuilder) EnableDiagnostics() {
	ib.enableDiagnostics = true
}

// GetDiagnostics retourne les diagnostics de la dernière construction.
func (ib *IndexBuilder) GetDiagnostics() *BuildDiagnostics {
	return ib.diagnostics
}

// BuildFromNetwork construit un index depuis un réseau RETE.
//
// Cette méthode parcourt tous les nœuds du réseau, extrait les champs
// utilisés, et construit l'index de dépendances.
//
// Paramètres :
//   - network : interface{} représentant le ReteNetwork
//     (typiquement *rete.ReteNetwork, mais on utilise interface{} pour
//     éviter les dépendances circulaires)
//
// Retourne un DependencyIndex complet et une erreur si échec.
//
// Note : Cette fonction utilise la reflection pour accéder aux champs
// du ReteNetwork afin d'éviter les dépendances circulaires.
func (ib *IndexBuilder) BuildFromNetwork(network interface{}) (*DependencyIndex, error) {
	idx := NewDependencyIndex()
	
	// Reset diagnostics
	ib.diagnostics = &BuildDiagnostics{}
	
	// Extraire les structures du réseau via reflection/type assertion
	// Note : Dans l'implémentation réelle, nous utiliserons les types concrets
	// Cette version utilise des interfaces pour la conception
	
	// TODO: Implémenter extraction depuis ReteNetwork
	// Pour l'instant, retourner index vide
	// L'implémentation complète sera faite dans le prompt 06 lors de l'intégration
	
	return idx, nil
}

// BuildFromAlphaNode traite un nœud alpha et l'ajoute à l'index.
//
// Paramètres :
//   - idx : index de dépendances à remplir
//   - nodeID : identifiant du nœud
//   - factType : type de fait
//   - condition : condition du nœud (AST)
//
// Retourne une erreur si l'extraction de champs échoue.
func (ib *IndexBuilder) BuildFromAlphaNode(
	idx *DependencyIndex,
	nodeID, factType string,
	condition interface{},
) error {
	fields, err := ib.alphaExtractor.ExtractFields(condition)
	if err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Alpha node %s: %v", nodeID, err),
			)
			ib.diagnostics.NodesSkipped++
		}
		return err
	}
	
	if len(fields) == 0 {
		if ib.enableDiagnostics {
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Alpha node %s: no fields extracted", nodeID),
			)
		}
	}
	
	idx.AddAlphaNode(nodeID, factType, fields)
	
	if ib.enableDiagnostics {
		ib.diagnostics.NodesProcessed++
		ib.diagnostics.FieldsExtracted += len(fields)
	}
	
	return nil
}

// BuildFromBetaNode traite un nœud beta et l'ajoute à l'index.
//
// Paramètres :
//   - idx : index de dépendances
//   - nodeID : identifiant du nœud
//   - factType : type de fait
//   - joinCondition : condition de jointure
//
// Retourne une erreur si l'extraction échoue.
func (ib *IndexBuilder) BuildFromBetaNode(
	idx *DependencyIndex,
	nodeID, factType string,
	joinCondition interface{},
) error {
	fields, err := ib.betaExtractor.ExtractFields(joinCondition)
	if err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Beta node %s: %v", nodeID, err),
			)
			ib.diagnostics.NodesSkipped++
		}
		return err
	}
	
	if len(fields) == 0 {
		if ib.enableDiagnostics {
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Beta node %s: no fields extracted", nodeID),
			)
		}
	}
	
	idx.AddBetaNode(nodeID, factType, fields)
	
	if ib.enableDiagnostics {
		ib.diagnostics.NodesProcessed++
		ib.diagnostics.FieldsExtracted += len(fields)
	}
	
	return nil
}

// BuildFromTerminalNode traite un nœud terminal et l'ajoute à l'index.
//
// Paramètres :
//   - idx : index de dépendances
//   - nodeID : identifiant du nœud
//   - factType : type de fait
//   - actions : liste des actions (AST)
//
// Retourne une erreur si l'extraction échoue.
func (ib *IndexBuilder) BuildFromTerminalNode(
	idx *DependencyIndex,
	nodeID, factType string,
	actions []interface{},
) error {
	allFields := make(map[string]bool)
	
	// Extraire champs depuis toutes les actions
	for _, action := range actions {
		fields, err := ib.actionExtractor.ExtractFields(action)
		if err != nil {
			if ib.enableDiagnostics {
				ib.diagnostics.Errors = append(
					ib.diagnostics.Errors,
					fmt.Sprintf("Terminal node %s action: %v", nodeID, err),
				)
			}
			// Continue avec les autres actions
			continue
		}
		
		for _, field := range fields {
			allFields[field] = true
		}
	}
	
	// Convertir map en slice
	fieldSlice := make([]string, 0, len(allFields))
	for field := range allFields {
		fieldSlice = append(fieldSlice, field)
	}
	
	if len(fieldSlice) == 0 {
		if ib.enableDiagnostics {
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Terminal node %s: no fields extracted", nodeID),
			)
		}
	}
	
	idx.AddTerminalNode(nodeID, factType, fieldSlice)
	
	if ib.enableDiagnostics {
		ib.diagnostics.NodesProcessed++
		ib.diagnostics.FieldsExtracted += len(fieldSlice)
	}
	
	return nil
}
```

### Tests : `rete/delta/index_builder_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

func TestNewIndexBuilder(t *testing.T) {
	builder := NewIndexBuilder()
	
	if builder == nil {
		t.Fatal("NewIndexBuilder returned nil")
	}
	
	if builder.alphaExtractor == nil {
		t.Error("alphaExtractor not initialized")
	}
	
	if builder.betaExtractor == nil {
		t.Error("betaExtractor not initialized")
	}
	
	if builder.actionExtractor == nil {
		t.Error("actionExtractor not initialized")
	}
	
	if builder.diagnostics == nil {
		t.Error("diagnostics not initialized")
	}
}

func TestIndexBuilder_BuildFromAlphaNode(t *testing.T) {
	builder := NewIndexBuilder()
	builder.EnableDiagnostics()
	idx := NewDependencyIndex()
	
	// Condition : product.price > 100
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":  "fieldAccess",
			"field": "price",
		},
		"right": 100,
	}
	
	err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Vérifier que le nœud est indexé
	nodes := idx.GetAffectedNodes("Product", "price")
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 node indexed, got %d", len(nodes))
	}
	
	if nodes[0].NodeID != "alpha1" {
		t.Errorf("Expected alpha1, got %s", nodes[0].NodeID)
	}
	
	// Vérifier diagnostics
	diag := builder.GetDiagnostics()
	if diag.NodesProcessed != 1 {
		t.Errorf("Expected 1 node processed, got %d", diag.NodesProcessed)
	}
	if diag.FieldsExtracted != 1 {
		t.Errorf("Expected 1 field extracted, got %d", diag.FieldsExtracted)
	}
}

func TestIndexBuilder_BuildFromAlphaNode_MultipleFields(t *testing.T) {
	builder := NewIndexBuilder()
	idx := NewDependencyIndex()
	
	// Condition : product.price > 100 && product.status == "active"
	condition := map[string]interface{}{
		"type": "binaryOp",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":  "fieldAccess",
				"field": "price",
			},
			"right": 100,
		},
		"right": map[string]interface{}{