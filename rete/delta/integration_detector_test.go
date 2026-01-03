// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
)

// TestIntegration_DetectorWithIndex vérifie l'intégration complète
// entre DeltaDetector et DependencyIndex
func TestIntegration_DetectorWithIndex(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - DeltaDetector + DependencyIndex")
	t.Log("=======================================================")

	// Étape 1 : Créer l'index de dépendances
	t.Log("\n📋 Étape 1 : Construction de l'index de dépendances")
	idx := NewDependencyIndex()

	// Ajouter des nœuds alpha dépendant de différents champs
	idx.AddAlphaNode("alpha_price", "Product", []string{"price"})
	idx.AddAlphaNode("alpha_stock", "Product", []string{"stock", "status"})
	idx.AddTerminalNode("terminal_update", "Product", []string{"price", "status"})

	stats := idx.GetStats()
	t.Logf("  ✅ Index créé : %d nœuds, %d champs", stats.NodeCount, stats.FieldCount)

	// Étape 2 : Créer le détecteur
	t.Log("\n📋 Étape 2 : Création du DeltaDetector")
	detector := NewDeltaDetector()
	t.Log("  ✅ Détecteur créé")

	// Étape 3 : Simuler une modification de fait
	t.Log("\n📋 Étape 3 : Simulation de modification")
	oldFact := map[string]interface{}{
		"id":     "p123",
		"price":  100.0,
		"stock":  50,
		"status": "active",
	}

	newFact := map[string]interface{}{
		"id":     "p123",
		"price":  150.0, // Modifié
		"stock":  50,
		"status": "active",
	}

	// Étape 4 : Détecter les changements
	t.Log("\n📋 Étape 4 : Détection des changements")
	factDelta, err := detector.DetectDelta(oldFact, newFact, "Product~p123", "Product")
	if err != nil {
		t.Fatalf("Erreur lors de la détection : %v", err)
	}

	t.Logf("  📝 Delta détecté : %d champ(s) modifié(s)", len(factDelta.Fields))
	for fieldName, change := range factDelta.Fields {
		t.Logf("    - %s : %v → %v", fieldName, change.OldValue, change.NewValue)
	}

	// Vérifier que seul "price" a changé
	if len(factDelta.Fields) != 1 {
		t.Errorf("Attendu 1 changement, reçu %d", len(factDelta.Fields))
	}

	if _, exists := factDelta.Fields["price"]; !exists {
		t.Error("Le champ 'price' devrait être dans le delta")
	}

	// Étape 5 : Obtenir les nœuds affectés
	t.Log("\n📋 Étape 5 : Identification des nœuds affectés")
	affectedNodes := idx.GetAffectedNodesForDelta(factDelta)

	t.Logf("  🔍 Nœuds affectés : %d", len(affectedNodes))
	for _, node := range affectedNodes {
		t.Logf("    - %s", node.String())
	}

	// Vérifier que les bons nœuds sont affectés
	// Seuls alpha_price et terminal_update dépendent de "price"
	// alpha_stock ne devrait PAS être affecté (dépend de stock et status)
	expectedAffected := 2 // alpha_price + terminal_update
	if len(affectedNodes) != expectedAffected {
		t.Errorf("Attendu %d nœuds affectés, reçu %d", expectedAffected, len(affectedNodes))
	}

	// Vérifier les IDs des nœuds affectés
	affectedIDs := make(map[string]bool)
	for _, node := range affectedNodes {
		affectedIDs[node.NodeID] = true
	}

	if !affectedIDs["alpha_price"] {
		t.Error("alpha_price devrait être affecté")
	}

	if !affectedIDs["terminal_update"] {
		t.Error("terminal_update devrait être affecté")
	}

	if affectedIDs["alpha_stock"] {
		t.Error("alpha_stock ne devrait PAS être affecté")
	}

	// Étape 6 : Calculer le ratio de changement
	t.Log("\n📋 Étape 6 : Calcul du ratio de changement")
	ratio := factDelta.ChangeRatio()
	t.Logf("  📊 Ratio de changement : %.2f (%d/%d)", ratio, len(factDelta.Fields), factDelta.FieldCount)

	// Étape 7 : Décision de stratégie de propagation
	t.Log("\n📋 Étape 7 : Décision de stratégie")
	const deltaThreshold = 0.3
	if ratio < deltaThreshold {
		t.Logf("  ✅ Propagation DELTA recommandée (ratio=%.2f < %.2f)", ratio, deltaThreshold)
		t.Logf("     → Propager uniquement vers %d nœuds affectés", len(affectedNodes))
	} else {
		t.Logf("  ⚠️  Propagation CLASSIQUE recommandée (ratio=%.2f >= %.2f)", ratio, deltaThreshold)
		t.Log("     → Utiliser retract + insert standard")
	}

	t.Log("\n🎉 TEST INTÉGRATION RÉUSSI")
}

// TestIntegration_MultipleChanges teste le cas avec plusieurs changements
func TestIntegration_MultipleChanges(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - Changements multiples")

	idx := NewDependencyIndex()
	idx.AddAlphaNode("alpha_price", "Product", []string{"price"})
	idx.AddAlphaNode("alpha_stock_status", "Product", []string{"stock", "status"})
	idx.AddTerminalNode("terminal_update", "Product", []string{"price", "status"})

	detector := NewDeltaDetector()

	oldFact := map[string]interface{}{
		"id":     "p123",
		"price":  100.0,
		"stock":  50,
		"status": "active",
	}

	newFact := map[string]interface{}{
		"id":     "p123",
		"price":  150.0, // Modifié
		"stock":  50,
		"status": "sold", // Modifié
	}

	factDelta, err := detector.DetectDelta(oldFact, newFact, "Product~p123", "Product")
	if err != nil {
		t.Fatalf("Erreur : %v", err)
	}

	t.Logf("Delta : %d changements détectés", len(factDelta.Fields))

	affectedNodes := idx.GetAffectedNodesForDelta(factDelta)
	t.Logf("Nœuds affectés : %d", len(affectedNodes))

	// Tous les nœuds devraient être affectés (price OU status)
	expectedAffected := 3 // alpha_price, alpha_stock_status, terminal_update
	if len(affectedNodes) != expectedAffected {
		t.Errorf("Attendu %d nœuds affectés, reçu %d", expectedAffected, len(affectedNodes))
	}

	ratio := factDelta.ChangeRatio()
	t.Logf("Ratio : %.2f", ratio)

	t.Log("✅ Test réussi")
}

// TestIntegration_QuickDetectionNoChanges teste DetectDeltaQuick sans changements
func TestIntegration_QuickDetectionNoChanges(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - DetectDeltaQuick (no-op)")

	detector := NewDeltaDetector()

	fact := map[string]interface{}{
		"id":     "p123",
		"price":  100.0,
		"status": "active",
	}

	// Pas de changements
	factDelta, err := detector.DetectDeltaQuick(fact, fact, "Product~p123", "Product")
	if err != nil {
		t.Fatalf("Erreur : %v", err)
	}

	if factDelta != nil {
		t.Error("Attendu nil (aucun changement), reçu un delta")
	}

	t.Log("✅ Aucun changement détecté (optimisation no-op)")
}

// TestIntegration_WithCache teste l'utilisation du cache
func TestIntegration_WithCache(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - Cache du détecteur")

	config := DefaultDetectorConfig()
	config.CacheComparisons = true
	detector := NewDeltaDetectorWithConfig(config)

	oldFact := map[string]interface{}{"price": 100.0}
	newFact := map[string]interface{}{"price": 150.0}

	// Première détection
	delta1, _ := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")

	// Deuxième détection (devrait venir du cache)
	delta2, _ := detector.DetectDelta(oldFact, newFact, "Product~123", "Product")

	metrics := detector.GetMetrics()
	t.Logf("Métriques : Comparisons=%d, CacheHits=%d, HitRate=%.2f%%",
		metrics.Comparisons, metrics.CacheHits, metrics.HitRate*100)

	if metrics.CacheHits == 0 {
		t.Error("Attendu au moins 1 cache hit")
	}

	if len(delta1.Fields) != len(delta2.Fields) {
		t.Error("Les deux deltas devraient être identiques")
	}

	t.Log("✅ Cache fonctionne correctement")
}
