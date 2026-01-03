# 🔗 Prompt 08 - Tests d'Intégration

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Développer une suite de tests d'intégration exhaustive qui valide le fonctionnement du système de propagation delta dans des scénarios réels complets, de bout en bout.

Ces tests garantissent que tous les composants fonctionnent correctement ensemble et que le système se comporte comme attendu dans des conditions réelles d'utilisation.

**⚠️ IMPORTANT** : Ce prompt génère du code de tests. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompts 01-07 validés** : Système complet + tests unitaires
- [x] **Couverture unitaire > 90%** : Tests unitaires passent
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `tests/integration/builtin_actions_test.go` - Exemple intégration existante
  - Tous les fichiers sources du package `rete/delta`

---

## 📂 Fichiers de Tests à Créer

```
tests/integration/
├── delta_propagation_test.go        # Tests propagation delta
├── delta_update_scenarios_test.go   # Scénarios Update complets
├── delta_regression_test.go         # Tests de non-régression
└── delta_performance_test.go        # Tests performance intégration

tests/e2e/
├── delta_e2e_test.go               # Tests end-to-end complets
└── delta_real_scenarios_test.go    # Scénarios réels métier
```

---

## 🔧 Tâche 1 : Tests d'Intégration Propagation Delta

### Fichier : `tests/integration/delta_propagation_test.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package integration

import (
    "testing"
    "github.com/yourusername/tsd/rete"
    "github.com/yourusername/tsd/rete/delta"
)

// TestDeltaPropagation_BasicFlow teste le flux de propagation delta de base.
func TestDeltaPropagation_BasicFlow(t *testing.T) {
    // 1. Setup : créer réseau RETE avec règles
    network := createTestNetwork(t)
    defer network.Shutdown()
    
    // 2. Activer propagation delta
    network.EnableDeltaPropagation = true
    err := network.InitializeDeltaPropagation()
    if err != nil {
        t.Fatalf("Failed to initialize delta propagation: %v", err)
    }
    
    // 3. Insérer un fait initial
    fact := map[string]interface{}{
        "id":     "P001",
        "name":   "Product A",
        "price":  100.0,
        "status": "active",
    }
    
    err = network.InsertFact(fact, "Product~P001", "Product")
    if err != nil {
        t.Fatalf("Failed to insert fact: %v", err)
    }
    
    // 4. Modifier le fait (1 champ)
    updatedFact := map[string]interface{}{
        "id":     "P001",
        "name":   "Product A",
        "price":  150.0, // Modifié
        "status": "active",
    }
    
    err = network.UpdateFact(fact, updatedFact, "Product~P001", "Product")
    if err != nil {
        t.Fatalf("Failed to update fact: %v", err)
    }
    
    // 5. Vérifier que propagation delta a été utilisée
    metrics := network.DeltaPropagator.GetMetrics()
    if metrics.DeltaPropagations != 1 {
        t.Errorf("Expected 1 delta propagation, got %d", metrics.DeltaPropagations)
    }
    
    if metrics.ClassicPropagations != 0 {
        t.Errorf("Expected 0 classic propagations, got %d", metrics.ClassicPropagations)
    }
    
    // 6. Vérifier que les nœuds affectés ont été activés
    // (implémentation dépend de la structure du réseau)
}

// TestDeltaPropagation_FallbackToClassic teste le fallback vers mode classique.
func TestDeltaPropagation_FallbackToClassic(t *testing.T) {
    network := createTestNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // Configuration pour forcer fallback
    config := network.DeltaPropagator.GetConfig()
    config.DeltaThreshold = 0.3 // 30%
    
    // Insérer fait
    fact := map[string]interface{}{
        "field1": "value1",
        "field2": "value2",
        "field3": "value3",
        "field4": "value4",
        "field5": "value5",
    }
    
    network.InsertFact(fact, "Test~1", "Test")
    
    // Modifier 40% des champs (> threshold)
    updatedFact := map[string]interface{}{
        "field1": "modified1",
        "field2": "modified2",
        "field3": "value3",
        "field4": "value4",
        "field5": "value5",
    }
    
    err := network.UpdateFact(fact, updatedFact, "Test~1", "Test")
    if err != nil {
        t.Fatalf("Failed to update: %v", err)
    }
    
    // Vérifier fallback utilisé
    metrics := network.DeltaPropagator.GetMetrics()
    if metrics.ClassicPropagations == 0 {
        t.Error("Expected classic propagation (fallback)")
    }
    
    if metrics.FallbacksDueToRatio == 0 {
        t.Error("Expected fallback due to ratio")
    }
}

// TestDeltaPropagation_ConcurrentUpdates teste les mises à jour concurrentes.
func TestDeltaPropagation_ConcurrentUpdates(t *testing.T) {
    network := createTestNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // Insérer 10 faits
    for i := 0; i < 10; i++ {
        fact := map[string]interface{}{
            "id":    fmt.Sprintf("P%03d", i),
            "price": float64(i * 100),
        }
        network.InsertFact(fact, fmt.Sprintf("Product~P%03d", i), "Product")
    }
    
    // Modifier tous les faits en parallèle
    done := make(chan bool, 10)
    errors := make(chan error, 10)
    
    for i := 0; i < 10; i++ {
        go func(id int) {
            oldFact := map[string]interface{}{
                "id":    fmt.Sprintf("P%03d", id),
                "price": float64(id * 100),
            }
            
            newFact := map[string]interface{}{
                "id":    fmt.Sprintf("P%03d", id),
                "price": float64(id * 100 * 1.1),
            }
            
            err := network.UpdateFact(oldFact, newFact, fmt.Sprintf("Product~P%03d", id), "Product")
            if err != nil {
                errors <- err
            }
            done <- true
        }(i)
    }
    
    // Attendre fin
    for i := 0; i < 10; i++ {
        <-done
    }
    
    close(errors)
    
    // Vérifier aucune erreur
    for err := range errors {
        t.Errorf("Concurrent update error: %v", err)
    }
    
    // Vérifier métriques
    metrics := network.DeltaPropagator.GetMetrics()
    if metrics.TotalPropagations != 10 {
        t.Errorf("Expected 10 propagations, got %d", metrics.TotalPropagations)
    }
}

// TestDeltaPropagation_PrimaryKeyChange teste le changement de clé primaire.
func TestDeltaPropagation_PrimaryKeyChange(t *testing.T) {
    network := createTestNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // Configurer champs PK
    config := network.DeltaPropagator.GetConfig()
    config.PrimaryKeyFields = []string{"id"}
    config.AllowPrimaryKeyChange = false
    
    oldFact := map[string]interface{}{
        "id":    "P001",
        "price": 100.0,
    }
    
    network.InsertFact(oldFact, "Product~P001", "Product")
    
    // Modifier la clé primaire
    newFact := map[string]interface{}{
        "id":    "P002", // PK changée
        "price": 100.0,
    }
    
    err := network.UpdateFact(oldFact, newFact, "Product~P001", "Product")
    if err != nil {
        t.Fatalf("Update failed: %v", err)
    }
    
    // Vérifier fallback classique utilisé
    metrics := network.DeltaPropagator.GetMetrics()
    if metrics.FallbacksDueToPK == 0 {
        t.Error("Expected fallback due to PK change")
    }
}

// Helper : créer un réseau de test
func createTestNetwork(t *testing.T) *rete.ReteNetwork {
    t.Helper()
    
    network := rete.NewReteNetwork()
    
    // Ajouter types
    network.AddType(rete.TypeDefinition{
        Name: "Product",
        Fields: map[string]rete.FieldDefinition{
            "id":     {Type: "string", PrimaryKey: true},
            "name":   {Type: "string"},
            "price":  {Type: "number"},
            "status": {Type: "string"},
        },
    })
    
    // Ajouter règles de test
    network.AddRule(rete.Rule{
        Name: "HighPriceAlert",
        Patterns: []rete.Pattern{
            {Type: "Product", Variable: "p", Conditions: "p.price > 120"},
        },
        Actions: []rete.Action{
            {Type: "Print", Arguments: []interface{}{"High price detected"}},
        },
    })
    
    network.Build()
    
    return network
}
```

---

## 🔧 Tâche 2 : Scénarios Update Complets

### Fichier : `tests/integration/delta_update_scenarios_test.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package integration

import (
    "testing"
)

// TestScenario_OrderProcessing teste un scénario métier complet.
func TestScenario_OrderProcessing(t *testing.T) {
    // Scénario : Traitement de commande avec états
    // 1. Créer commande (pending)
    // 2. Valider commande (pending → confirmed)
    // 3. Préparer commande (confirmed → preparing)
    // 4. Expédier commande (preparing → shipped)
    
    network := setupOrderNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // 1. Créer commande
    order := map[string]interface{}{
        "id":       "ORD001",
        "customer": "CUST001",
        "total":    250.0,
        "status":   "pending",
    }
    
    network.InsertFact(order, "Order~ORD001", "Order")
    
    // 2. Valider commande
    order["status"] = "confirmed"
    err := network.UpdateFact(
        map[string]interface{}{"id": "ORD001", "status": "pending"},
        order,
        "Order~ORD001",
        "Order",
    )
    
    if err != nil {
        t.Fatalf("Failed to confirm order: %v", err)
    }
    
    // Vérifier propagation delta utilisée (1 seul champ modifié)
    metrics := network.DeltaPropagator.GetMetrics()
    if metrics.DeltaPropagations == 0 {
        t.Error("Expected delta propagation for status change")
    }
    
    // 3. Préparer commande
    order["status"] = "preparing"
    network.UpdateFact(
        map[string]interface{}{"id": "ORD001", "status": "confirmed"},
        order,
        "Order~ORD001",
        "Order",
    )
    
    // 4. Expédier commande
    order["status"] = "shipped"
    order["shipped_at"] = "2025-01-02T10:00:00Z"
    
    network.UpdateFact(
        map[string]interface{}{"id": "ORD001", "status": "preparing"},
        order,
        "Order~ORD001",
        "Order",
    )
    
    // Vérifier que toutes les transitions ont utilisé delta
    finalMetrics := network.DeltaPropagator.GetMetrics()
    if finalMetrics.DeltaPropagations != 3 {
        t.Errorf("Expected 3 delta propagations, got %d", finalMetrics.DeltaPropagations)
    }
}

// TestScenario_InventoryManagement teste gestion d'inventaire.
func TestScenario_InventoryManagement(t *testing.T) {
    // Scénario : Gestion stock
    // 1. Créer produit avec stock initial
    // 2. Vendre unités (stock décrémente)
    // 3. Réapprovisionner (stock incrémente)
    // 4. Vérifier règles d'alerte stock bas
    
    network := setupInventoryNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // 1. Produit initial
    product := map[string]interface{}{
        "id":       "PROD001",
        "name":     "Widget",
        "stock":    100,
        "min_stock": 20,
    }
    
    network.InsertFact(product, "Product~PROD001", "Product")
    
    // 2. Vendre 30 unités
    oldProduct := copyMap(product)
    product["stock"] = 70
    
    network.UpdateFact(oldProduct, product, "Product~PROD001", "Product")
    
    // 3. Vendre encore 55 unités (total: 15, alerte stock bas)
    oldProduct = copyMap(product)
    product["stock"] = 15
    
    network.UpdateFact(oldProduct, product, "Product~PROD001", "Product")
    
    // Vérifier que la règle d'alerte a été activée
    // (implémentation dépend du système d'actions)
    
    // 4. Réapprovisionner
    oldProduct = copyMap(product)
    product["stock"] = 115
    
    network.UpdateFact(oldProduct, product, "Product~PROD001", "Product")
    
    // Vérifier métriques
    metrics := network.DeltaPropagator.GetMetrics()
    
    // Toutes les updates ne modifient qu'un champ (stock)
    if metrics.AvgFieldsPerPropagation > 1.5 {
        t.Errorf("Expected ~1 field per propagation, got %.2f", metrics.AvgFieldsPerPropagation)
    }
}

// TestScenario_ComplexRelationships teste relations complexes.
func TestScenario_ComplexRelationships(t *testing.T) {
    // Scénario : Relations Customer-Order-Product
    // 1. Créer customer, products, order
    // 2. Modifier order.total → doit propager vers règles customer
    // 3. Modifier product.price → doit recalculer order.total
    
    network := setupRelationshipNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // Setup données
    customer := map[string]interface{}{
        "id":    "C001",
        "name":  "Alice",
        "tier":  "bronze",
    }
    
    product := map[string]interface{}{
        "id":    "P001",
        "price": 50.0,
    }
    
    order := map[string]interface{}{
        "id":          "O001",
        "customer_id": "C001",
        "product_id":  "P001",
        "quantity":    2,
        "total":       100.0,
    }
    
    network.InsertFact(customer, "Customer~C001", "Customer")
    network.InsertFact(product, "Product~P001", "Product")
    network.InsertFact(order, "Order~O001", "Order")
    
    // Modifier price produit
    oldProduct := copyMap(product)
    product["price"] = 60.0
    
    network.UpdateFact(oldProduct, product, "Product~P001", "Product")
    
    // Le changement de price devrait déclencher recalcul de order.total
    // via une règle (implémentation dépend des règles définies)
    
    // Vérifier que propagation delta a été efficace
    metrics := network.DeltaPropagator.GetMetrics()
    
    efficiencyRatio := metrics.GetEfficiencyRatio()
    if efficiencyRatio < 0.5 {
        t.Errorf("Low efficiency ratio: %.2f (expected > 0.5)", efficiencyRatio)
    }
}

// Helpers
func setupOrderNetwork(t *testing.T) *rete.ReteNetwork {
    t.Helper()
    // Setup network avec types Order et règles de transition
    return createTestNetwork(t)
}

func setupInventoryNetwork(t *testing.T) *rete.ReteNetwork {
    t.Helper()
    // Setup network avec règles de gestion stock
    return createTestNetwork(t)
}

func setupRelationshipNetwork(t *testing.T) *rete.ReteNetwork {
    t.Helper()
    // Setup network avec relations Customer-Order-Product
    return createTestNetwork(t)
}

func copyMap(m map[string]interface{}) map[string]interface{} {
    copy := make(map[string]interface{})
    for k, v := range m {
        copy[k] = v
    }
    return copy
}
```

---

## 🔧 Tâche 3 : Tests de Non-Régression

### Fichier : `tests/integration/delta_regression_test.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package integration

import (
    "testing"
)

// TestRegression_ExistingUpdateBehavior vérifie qu'Update fonctionne comme avant.
func TestRegression_ExistingUpdateBehavior(t *testing.T) {
    // Test : Avec delta DÉSACTIVÉ, comportement identique à avant
    
    network := createTestNetwork(t)
    defer network.Shutdown()
    
    // Delta désactivé
    network.EnableDeltaPropagation = false
    
    fact := map[string]interface{}{
        "id":    "P001",
        "price": 100.0,
    }
    
    network.InsertFact(fact, "Product~P001", "Product")
    
    updatedFact := map[string]interface{}{
        "id":    "P001",
        "price": 150.0,
    }
    
    err := network.UpdateFact(fact, updatedFact, "Product~P001", "Product")
    if err != nil {
        t.Fatalf("Update failed: %v", err)
    }
    
    // Vérifier que le fait a bien été mis à jour dans le storage
    storedFact, err := network.Storage.GetFact("Product~P001")
    if err != nil {
        t.Fatalf("Failed to get fact: %v", err)
    }
    
    if storedFact["price"] != 150.0 {
        t.Errorf("Expected price 150.0, got %v", storedFact["price"])
    }
}

// TestRegression_AllExistingTests rejoue tous les tests Update existants.
func TestRegression_AllExistingTests(t *testing.T) {
    // Importer et exécuter tous les tests d'action Update existants
    // avec delta activé, vérifier résultats identiques
    
    t.Run("with_delta_enabled", func(t *testing.T) {
        // ... exécuter suite de tests existante avec delta ON
    })
    
    t.Run("with_delta_disabled", func(t *testing.T) {
        // ... exécuter suite de tests existante avec delta OFF
    })
    
    // Comparer résultats : doivent être identiques
}

// TestRegression_NoChangeDetection teste détection no-op existante.
func TestRegression_NoChangeDetection(t *testing.T) {
    // Vérifier que détection no-op fonctionne toujours
    
    network := createTestNetwork(t)
    defer network.Shutdown()
    
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    fact := map[string]interface{}{
        "id":    "P001",
        "price": 100.0,
    }
    
    network.InsertFact(fact, "Product~P001", "Product")
    
    // Update avec valeurs identiques
    err := network.UpdateFact(fact, fact, "Product~P001", "Product")
    if err != nil {
        t.Fatalf("No-op update failed: %v", err)
    }
    
    // Vérifier qu'aucune propagation n'a eu lieu
    metrics := network.DeltaPropagator.GetMetrics()
    if metrics.TotalPropagations != 0 {
        t.Errorf("Expected 0 propagations for no-op, got %d", metrics.TotalPropagations)
    }
}
```

---

## 🔧 Tâche 4 : Tests End-to-End

### Fichier : `tests/e2e/delta_e2e_test.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
    "testing"
    "time"
)

// TestE2E_FullWorkflow teste un workflow complet de bout en bout.
func TestE2E_FullWorkflow(t *testing.T) {
    // Test E2E complet : charger fichier TSD, exécuter, vérifier résultats
    
    // 1. Charger fichier TSD avec règles
    tsdContent := `
        type Product {
            id: string primary_key,
            price: number,
            discount: number,
            final_price: number
        }
        
        rule "CalculateFinalPrice" {
            when {
                p: Product(p.price > 0)
            }
            then {
                Update(p, {
                    final_price: p.price * (1 - p.discount)
                })
            }
        }
    `
    
    // 2. Compiler et créer réseau
    network, err := compileTSD(tsdContent)
    if err != nil {
        t.Fatalf("Failed to compile TSD: %v", err)
    }
    defer network.Shutdown()
    
    // 3. Activer propagation delta
    network.EnableDeltaPropagation = true
    network.InitializeDeltaPropagation()
    
    // 4. Insérer faits
    product := map[string]interface{}{
        "id":       "P001",
        "price":    100.0,
        "discount": 0.1,
    }
    
    network.InsertFact(product, "Product~P001", "Product")
    
    // Attendre propagation
    time.Sleep(100 * time.Millisecond)
    
    // 5. Vérifier calcul final_price
    fact, _ := network.Storage.GetFact("Product~P001")
    if fact["final_price"] != 90.0 {
        t.Errorf("Expected final_price 90.0, got %v", fact["final_price"])
    }
    
    // 6. Modifier discount
    oldProduct := copyMap(fact)
    product["discount"] = 0.2
    
    network.UpdateFact(oldProduct, product, "Product~P001", "Product")
    
    time.Sleep(100 * time.Millisecond)
    
    // 7. Vérifier nouveau final_price
    fact, _ = network.Storage.GetFact("Product~P001")
    if fact["final_price"] != 80.0 {
        t.Errorf("Expected final_price 80.0, got %v", fact["final_price"])
    }
    
    // 8. Vérifier métriques
    metrics := network.DeltaPropagator.GetMetrics()
    
    if metrics.DeltaPropagations == 0 {
        t.Error("Expected delta propagations")
    }
    
    t.Logf("E2E metrics: %+v", metrics)
}
```

---

## ✅ Validation

Après implémentation, exécuter :

```bash
# 1. Tests d'intégration
go test ./tests/integration/... -v -run Delta

# 2. Tests E2E
go test ./tests/e2e/... -v -run Delta

# 3. Tests de régression
go test ./tests/integration/... -v -run Regression

# 4. Tous les tests
go test ./tests/... -v

# 5. Race detector
go test ./tests/integration/... -race

# 6. Validation complète
make test-integration
make test-e2e
```

**Critères de succès** :
- [ ] Tous les tests d'intégration passent (100%)
- [ ] Tous les tests E2E passent (100%)
- [ ] Aucune régression détectée
- [ ] Scénarios métier validés
- [ ] Aucune race condition
- [ ] Performance acceptable (voir Prompt 09)

---

## 📊 Livrables

À la fin de ce prompt :

1. **Tests d'intégration** :
   - ✅ `delta_propagation_test.go` - Tests propagation
   - ✅ `delta_update_scenarios_test.go` - Scénarios complets
   - ✅ `delta_regression_test.go` - Non-régression

2. **Tests E2E** :
   - ✅ `delta_e2e_test.go` - Workflow complet

3. **Validation** :
   - ✅ Rapport de tests d'intégration
   - ✅ Scénarios métier validés

---

## 🚀 Commit

Une fois validé :

```bash
git add tests/
git commit -m "test(integration): [Prompt 08] Tests d'intégration et E2E pour propagation delta

- Tests propagation delta (basic flow, fallback, concurrence)
- Scénarios métier complets (order processing, inventory, relationships)
- Tests de non-régression (comportement identique delta ON/OFF)
- Tests E2E workflow complet
- Validation scénarios réels
- Aucune régression détectée"
```

---

## 🚦 Prochaine Étape

Passer au **Prompt 09 - Optimisations et Profiling**

---

**Durée estimée** : 3-4 heures  
**Difficulté** : Moyenne-Élevée  
**Prérequis** : Prompts 01-07 validés  
**Couverture cible** : 100% scénarios métier