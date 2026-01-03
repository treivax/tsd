# 🔗 Prompt 06 - Intégration avec Action Update

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Intégrer le système de propagation delta dans l'action `Update` existante. Modifier le flux d'exécution pour utiliser la propagation sélective au lieu de Retract+Insert classique.

Cette intégration est le point où toutes les composantes développées (détection delta, index, propagation) sont connectées au moteur RETE existant.

**⚠️ IMPORTANT** : Ce prompt génère du code. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompts 01-05 validés** : Tous les composants delta implémentés
- [x] **Tests passent** : `go test ./rete/delta/... -v` (100% success)
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `REPORTS/sequence_update_actuel.md`
  - `rete/action_executor_evaluation.go` - Évaluation Update actuelle
  - `rete/action_executor_facts.go` - Gestion faits
  - `rete/network_manager.go` - Insertion/retract
  - `rete/delta/delta_propagator.go` - Propagateur delta

---

## 📂 Fichiers à Modifier/Créer

```
rete/
├── network.go                      # Ajouter DeltaPropagator
├── network_manager.go              # Ajouter UpdateFact avec delta
├── action_executor_facts.go       # Intégrer propagation delta
└── delta/
    ├── integration.go              # Helper intégration RETE
    ├── integration_test.go         # Tests intégration
    └── network_callbacks.go        # Callbacks vers réseau RETE
```

---

## 🔧 Tâche 1 : Extension ReteNetwork

### Fichier : `rete/network.go`

**Modifications** :

```go
// Ajouter à la structure ReteNetwork :

type ReteNetwork struct {
    // ... champs existants ...
    
    // Propagation delta (nouveau)
    DeltaPropagator       *delta.DeltaPropagator `json:"-"`
    DependencyIndex       *delta.DependencyIndex `json:"-"`
    EnableDeltaPropagation bool                  `json:"-"`
    
    // ... reste des champs ...
}
```

**Ajouter méthode d'initialisation** :

```go
// InitializeDeltaPropagation initialise le système de propagation delta.
//
// Cette méthode construit l'index de dépendances depuis le réseau existant
// et crée le DeltaPropagator configuré.
//
// Doit être appelée après la construction complète du réseau RETE.
func (rn *ReteNetwork) InitializeDeltaPropagation() error {
    // 1. Construire l'index de dépendances
    indexBuilder := delta.NewIndexBuilder()
    indexBuilder.EnableDiagnostics()
    
    // 2. Parcourir les nœuds alpha
    for nodeID, alphaNode := range rn.AlphaNodes {
        factType := alphaNode.FactType
        condition := alphaNode.Condition
        
        err := indexBuilder.BuildFromAlphaNode(
            rn.DependencyIndex,
            nodeID,
            factType,
            condition,
        )
        if err != nil {
            return fmt.Errorf("failed to index alpha node %s: %w", nodeID, err)
        }
    }
    
    // 3. Parcourir les nœuds beta
    for nodeID, betaNode := range rn.BetaNodes {
        // Extraction type et condition de jointure
        // (implémentation dépend de la structure BetaNode)
        err := indexBuilder.BuildFromBetaNode(
            rn.DependencyIndex,
            nodeID,
            factType,
            joinCondition,
        )
        if err != nil {
            return fmt.Errorf("failed to index beta node %s: %w", nodeID, err)
        }
    }
    
    // 4. Parcourir les nœuds terminaux
    for nodeID, terminalNode := range rn.TerminalNodes {
        factType := terminalNode.FactType
        actions := terminalNode.Actions
        
        err := indexBuilder.BuildFromTerminalNode(
            rn.DependencyIndex,
            nodeID,
            factType,
            actions,
        )
        if err != nil {
            return fmt.Errorf("failed to index terminal node %s: %w", nodeID, err)
        }
    }
    
    // 5. Construire le DeltaPropagator
    propagator, err := delta.NewDeltaPropagatorBuilder().
        WithIndex(rn.DependencyIndex).
        WithDetector(delta.NewDeltaDetector()).
        WithStrategy(&delta.SequentialStrategy{}).
        WithConfig(delta.DefaultPropagationConfig()).
        WithPropagateCallback(rn.propagateDeltaToNode).
        Build()
    
    if err != nil {
        return fmt.Errorf("failed to build delta propagator: %w", err)
    }
    
    rn.DeltaPropagator = propagator
    
    return nil
}

// propagateDeltaToNode est le callback pour propager un delta vers un nœud.
func (rn *ReteNetwork) propagateDeltaToNode(nodeID string, delta *delta.FactDelta) error {
    // Trouver le nœud
    if alphaNode, exists := rn.AlphaNodes[nodeID]; exists {
        return rn.propagateDeltaToAlpha(alphaNode, delta)
    }
    
    if betaNode, exists := rn.BetaNodes[nodeID]; exists {
        return rn.propagateDeltaToBeta(betaNode, delta)
    }
    
    if terminalNode, exists := rn.TerminalNodes[nodeID]; exists {
        return rn.propagateDeltaToTerminal(terminalNode, delta)
    }
    
    return fmt.Errorf("node not found: %s", nodeID)
}

// propagateDeltaToAlpha propage un delta vers un nœud alpha.
func (rn *ReteNetwork) propagateDeltaToAlpha(node *AlphaNode, delta *delta.FactDelta) error {
    // Évaluer la condition avec le fait modifié
    // Si condition satisfaite, propager vers successeurs
    // (implémentation dépend de la structure AlphaNode)
    return nil
}

// propagateDeltaToBeta propage un delta vers un nœud beta.
func (rn *ReteNetwork) propagateDeltaToBeta(node interface{}, delta *delta.FactDelta) error {
    // Ré-évaluer les jointures concernées
    // (implémentation dépend de la structure BetaNode)
    return nil
}

// propagateDeltaToTerminal propage un delta vers un nœud terminal.
func (rn *ReteNetwork) propagateDeltaToTerminal(node *TerminalNode, delta *delta.FactDelta) error {
    // Activer la règle si nécessaire
    // (implémentation dépend de la structure TerminalNode)
    return nil
}
```

---

## 🔧 Tâche 2 : Modification NetworkManager

### Fichier : `rete/network_manager.go`

**Ajouter nouvelle méthode UpdateFact** :

```go
// UpdateFact met à jour un fait dans le réseau avec propagation delta.
//
// Cette méthode est le nouveau point d'entrée pour les mises à jour de faits.
// Elle utilise la propagation delta si activée et applicable, sinon fallback
// sur Retract+Insert classique.
//
// Paramètres :
//   - oldFact : fait avant modification
//   - newFact : fait après modification
//   - factID : identifiant interne du fait
//   - factType : type du fait
//
// Retourne une erreur si la mise à jour échoue.
func (nm *NetworkManager) UpdateFact(
    oldFact, newFact map[string]interface{},
    factID, factType string,
) error {
    // Vérifier si propagation delta activée
    if !nm.network.EnableDeltaPropagation || nm.network.DeltaPropagator == nil {
        // Fallback mode classique
        return nm.updateFactClassic(oldFact, newFact, factID, factType)
    }
    
    // Tenter propagation delta
    err := nm.network.DeltaPropagator.PropagateUpdate(
        oldFact, newFact,
        factID, factType,
    )
    
    if err != nil {
        // Si échec et retry activé, fallback classique
        if nm.network.DeltaPropagator.GetConfig().RetryOnError {
            return nm.updateFactClassic(oldFact, newFact, factID, factType)
        }
        return err
    }
    
    // Mettre à jour le storage
    if nm.network.Storage != nil {
        if err := nm.network.Storage.UpdateFact(factID, newFact); err != nil {
            return fmt.Errorf("failed to update fact in storage: %w", err)
        }
    }
    
    return nil
}

// updateFactClassic effectue une mise à jour classique (Retract+Insert).
func (nm *NetworkManager) updateFactClassic(
    oldFact, newFact map[string]interface{},
    factID, factType string,
) error {
    // Retract ancien fait
    if err := nm.RetractFact(factID, factType); err != nil {
        return fmt.Errorf("retract failed: %w", err)
    }
    
    // Insert nouveau fait
    if err := nm.InsertFact(newFact, factID, factType); err != nil {
        return fmt.Errorf("insert failed: %w", err)
    }
    
    return nil
}
```

---

## 🔧 Tâche 3 : Intégration dans ActionExecutor

### Fichier : `rete/action_executor_facts.go`

**Modifier la méthode executeUpdateWithModifications** :

```go
// executeUpdateWithModifications exécute une action Update avec modifications.
//
// Ancien comportement : Retract + Insert
// Nouveau comportement : Propagation delta si applicable
func (ae *ActionExecutor) executeUpdateWithModifications(
    action map[string]interface{},
    ctx *ExecutionContext,
) error {
    // 1. Évaluer la variable (fait à modifier)
    varName, ok := action["variable"].(string)
    if !ok {
        return fmt.Errorf("invalid variable in Update action")
    }
    
    oldFact := ctx.GetVariable(varName)
    if oldFact == nil {
        return fmt.Errorf("variable %s not found", varName)
    }
    
    // 2. Extraire l'ID et le type du fait
    factID := ae.getFactID(oldFact)
    factType := ae.getFactType(oldFact)
    
    // 3. Appliquer les modifications
    modifications, ok := action["modifications"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("invalid modifications in Update action")
    }
    
    newFact := ae.applyModifications(oldFact, modifications, ctx)
    
    // 4. Vérification no-op (déjà implémentée)
    if ae.areFactsEqual(oldFact, newFact) {
        // Aucun changement, skip
        return nil
    }
    
    // 5. Utiliser UpdateFact avec support delta
    networkManager := ae.getNetworkManager() // Accès au NetworkManager
    
    return networkManager.UpdateFact(oldFact, newFact, factID, factType)
}

// applyModifications applique les modifications à un fait.
func (ae *ActionExecutor) applyModifications(
    oldFact map[string]interface{},
    modifications map[string]interface{},
    ctx *ExecutionContext,
) map[string]interface{} {
    // Créer une copie du fait
    newFact := make(map[string]interface{})
    for k, v := range oldFact {
        newFact[k] = v
    }
    
    // Appliquer chaque modification
    for field, value := range modifications {
        // Évaluer la valeur (peut être une expression)
        evaluatedValue, err := ae.evaluateArgument(value, ctx)
        if err != nil {
            // Log erreur et skip cette modification
            continue
        }
        
        newFact[field] = evaluatedValue
    }
    
    return newFact
}

// getFactID extrait l'ID interne d'un fait.
func (ae *ActionExecutor) getFactID(fact map[string]interface{}) string {
    // L'ID interne est stocké dans un champ spécial (ex: "__internal_id")
    // ou construit depuis les clés primaires
    if id, ok := fact["__internal_id"].(string); ok {
        return id
    }
    
    // Fallback : construire depuis le type et les PKs
    // (implémentation dépend de la gestion des IDs)
    return ""
}

// getFactType extrait le type d'un fait.
func (ae *ActionExecutor) getFactType(fact map[string]interface{}) string {
    if factType, ok := fact["__type"].(string); ok {
        return factType
    }
    return "Unknown"
}
```

---

## 🔧 Tâche 4 : Callbacks Réseau RETE

### Fichier : `rete/delta/network_callbacks.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// NetworkCallbacks définit les callbacks pour interagir avec le réseau RETE.
//
// Cette interface découple le package delta du package rete principal,
// évitant ainsi les dépendances circulaires.
type NetworkCallbacks interface {
    // PropagateToAlpha propage un delta vers un nœud alpha
    PropagateToAlpha(nodeID string, delta *FactDelta) error
    
    // PropagateToBeta propage un delta vers un nœud beta
    PropagateToBeta(nodeID string, delta *FactDelta) error
    
    // PropagateToTerminal propage un delta vers un nœud terminal
    PropagateToTerminal(nodeID string, delta *FactDelta) error
    
    // GetNode récupère un nœud par son ID
    GetNode(nodeID string) (interface{}, error)
    
    // UpdateStorage met à jour le storage avec le fait modifié
    UpdateStorage(factID string, newFact map[string]interface{}) error
}

// DefaultNetworkCallbacks est une implémentation par défaut (no-op).
type DefaultNetworkCallbacks struct{}

func (dnc *DefaultNetworkCallbacks) PropagateToAlpha(nodeID string, delta *FactDelta) error {
    return nil
}

func (dnc *DefaultNetworkCallbacks) PropagateToBeta(nodeID string, delta *FactDelta) error {
    return nil
}

func (dnc *DefaultNetworkCallbacks) PropagateToTerminal(nodeID string, delta *FactDelta) error {
    return nil
}

func (dnc *DefaultNetworkCallbacks) GetNode(nodeID string) (interface{}, error) {
    return nil, nil
}

func (dnc *DefaultNetworkCallbacks) UpdateStorage(factID string, newFact map[string]interface{}) error {
    return nil
}
```

---

## 🔧 Tâche 5 : Helper Intégration

### Fichier : `rete/delta/integration.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
    "fmt"
)

// IntegrationHelper facilite l'intégration du système delta avec RETE.
type IntegrationHelper struct {
    propagator *DeltaPropagator
    index      *DependencyIndex
    callbacks  NetworkCallbacks
}

// NewIntegrationHelper crée un nouveau helper d'intégration.
func NewIntegrationHelper(
    propagator *DeltaPropagator,
    index *DependencyIndex,
    callbacks NetworkCallbacks,
) *IntegrationHelper {
    return &IntegrationHelper{
        propagator: propagator,
        index:      index,
        callbacks:  callbacks,
    }
}

// ProcessUpdate traite une mise à jour de fait de bout en bout.
//
// Cette méthode coordonne :
// 1. Détection du delta
// 2. Recherche des nœuds affectés
// 3. Décision delta vs classique
// 4. Propagation
// 5. Mise à jour storage
//
// Paramètres :
//   - oldFact, newFact : faits avant/après
//   - factID, factType : identifiant et type
//
// Retourne une erreur si le traitement échoue.
func (ih *IntegrationHelper) ProcessUpdate(
    oldFact, newFact map[string]interface{},
    factID, factType string,
) error {
    // Déléguer à la propagation delta
    err := ih.propagator.PropagateUpdate(oldFact, newFact, factID, factType)
    if err != nil {
        return fmt.Errorf("delta propagation failed: %w", err)
    }
    
    // Mettre à jour le storage
    if ih.callbacks != nil {
        if err := ih.callbacks.UpdateStorage(factID, newFact); err != nil {
            return fmt.Errorf("storage update failed: %w", err)
        }
    }
    
    return nil
}

// RebuildIndex reconstruit l'index de dépendances depuis le réseau.
//
// Cette méthode doit être appelée si le réseau RETE est modifié
// (ajout/suppression de règles).
func (ih *IntegrationHelper) RebuildIndex(networkNodes interface{}) error {
    // Clear index existant
    ih.index.Clear()
    
    // Reconstruire depuis les nœuds
    // (implémentation dépend de la structure du réseau)
    
    return nil
}

// GetMetrics retourne les métriques du propagateur.
func (ih *IntegrationHelper) GetMetrics() PropagationMetrics {
    return ih.propagator.GetMetrics()
}
```

---

## ✅ Validation

Après implémentation, exécuter :

```bash
# 1. Formattage
go fmt ./rete/...
goimports -w ./rete/

# 2. Validation statique
go vet ./rete/...
staticcheck ./rete/...

# 3. Tests unitaires
go test ./rete/delta/... -v
go test ./rete/... -v -run TestUpdate

# 4. Tests d'intégration
go test ./tests/integration/... -v -run Delta

# 5. Tests de régression
go test ./tests/... -v

# 6. Race detector
go test ./rete/... -race

# 7. Validation complète
make validate
make test
```

**Critères de succès** :
- [ ] Tous les tests passent (100%)
- [ ] Aucune régression sur tests existants
- [ ] Action Update utilise propagation delta quand applicable
- [ ] Fallback classique fonctionne si delta désactivé
- [ ] Métriques collectées correctement
- [ ] Aucune race condition

---

## 🧪 Tests d'Intégration

### Fichier : `rete/delta/integration_test.go`

**Scénarios de test** :

```go
// 1. Test Update avec delta activé
func TestIntegration_UpdateWithDelta(t *testing.T) {
    // Setup réseau + index + propagateur
    // Insérer un fait
    // Modifier 1 champ
    // Vérifier propagation delta utilisée
    // Vérifier métriques
}

// 2. Test Update avec delta désactivé (fallback)
func TestIntegration_UpdateWithoutDelta(t *testing.T) {
    // Setup avec EnableDeltaPropagation = false
    // Modifier fait
    // Vérifier Retract+Insert classique utilisé
}

// 3. Test Update avec changement PK (fallback)
func TestIntegration_UpdatePrimaryKey(t *testing.T) {
    // Modifier clé primaire
    // Vérifier fallback classique
}

// 4. Test Update avec ratio élevé (fallback)
func TestIntegration_UpdateHighRatio(t *testing.T) {
    // Modifier 80% des champs
    // Vérifier fallback classique
}

// 5. Test Update concurrent
func TestIntegration_ConcurrentUpdates(t *testing.T) {
    // Lancer plusieurs updates en parallèle
    // Vérifier cohérence
}

// 6. Test reconstruction index
func TestIntegration_RebuildIndex(t *testing.T) {
    // Ajouter une règle au réseau
    // Reconstruire index
    // Vérifier nouveaux nœuds indexés
}
```

---

## 📊 Livrables

À la fin de ce prompt :

1. **Code** :
   - ✅ `rete/network.go` - Extension avec DeltaPropagator
   - ✅ `rete/network_manager.go` - UpdateFact avec delta
   - ✅ `rete/action_executor_facts.go` - Intégration Update
   - ✅ `rete/delta/network_callbacks.go` - Interface callbacks
   - ✅ `rete/delta/integration.go` - Helper intégration

2. **Tests** :
   - ✅ `rete/delta/integration_test.go`
   - ✅ Tests de régression Update

3. **Validation** :
   - ✅ Tous tests passent
   - ✅ Aucune régression
   - ✅ Documentation inline complète

---

## 🚀 Commit

Une fois validé :

```bash
git add rete/
git commit -m "feat(rete): [Prompt 06] Intégration propagation delta dans Update

- Extension ReteNetwork avec DeltaPropagator
- NetworkManager.UpdateFact avec support delta
- ActionExecutor intégré avec propagation sélective
- Callbacks pour découplage delta/rete
- Helper IntegrationHelper pour coordination
- Fallback automatique vers Retract+Insert classique
- Tests d'intégration complets
- Métriques de propagation collectées
- Aucune régression sur tests existants"
```

---

## 🚦 Prochaine Étape

Passer au **Prompt 07 - Tests Unitaires**

---

**Durée estimée** : 2-3 heures  
**Difficulté** : Élevée (intégration système)  
**Prérequis** : Prompts 01-05 validés  
**Couverture cible** : > 85%