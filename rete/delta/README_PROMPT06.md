# 📚 Guide - Intégration Delta Propagation (Prompt 06)

## 🎯 Vue d'ensemble

Ce répertoire contient l'implémentation de l'intégration du système de propagation delta dans le réseau RETE, conformément au **Prompt 06** (`.github/prompts/06_integration_update.md`).

## 📁 Structure

### Fichiers de Code
- `network_callbacks.go` - Interface de découplage delta/rete
- `integration.go` - Helper d'intégration des composants
- `integration_helper_test.go` - Tests d'intégration

### Documentation
- `SYNTHESE_PROMPT06.md` - **À LIRE EN PREMIER** ⭐
- `EXECUTION_SUMMARY_PROMPT06.md` - Rapport d'exécution détaillé
- `CODE_REVIEW_PROMPT06.md` - Revue de code complète

## 🚀 Démarrage Rapide

### 1. Activer la Propagation Delta

```go
// Créer le réseau RETE
network := NewReteNetwork(...)

// Initialiser la propagation delta
err := network.InitializeDeltaPropagation()
if err != nil {
    log.Fatal(err)
}

// Activer la propagation delta
network.EnableDeltaPropagation = true
```

### 2. Utiliser UpdateFact

```go
// L'update utilisera automatiquement la propagation delta si activée
fact := &Fact{
    ID:   "product_1",
    Type: "Product",
    Fields: map[string]interface{}{
        "price": 199.99,  // Modification
        "name":  "Laptop",
    },
}

err := network.UpdateFact(fact)
// La propagation delta est tentée
// Fallback automatique vers Retract+Insert si échec
```

### 3. Consulter les Métriques

```go
// Récupérer les statistiques
helper := network.IntegrationHelper
if helper != nil {
    metrics := helper.GetMetrics()
    fmt.Printf("Delta propagations: %d\n", metrics.DeltaPropagations)
    fmt.Printf("Classic fallbacks: %d\n", metrics.ClassicPropagations)
    
    indexStats := helper.GetIndexMetrics()
    fmt.Printf("Nodes indexed: %d\n", indexStats.NodeCount)
    fmt.Printf("Field dependencies: %d\n", indexStats.FieldCount)
}
```

## ⚠️  Limitations Actuelles

### 🔴 Critiques (À résoudre avant production)

1. **Propagation non opérationnelle**
   - La méthode `propagateDeltaToNode()` ne fait que logger
   - Aucune propagation réelle vers les nœuds alpha/beta/terminal
   - **Solution**: Voir Prompt 07

2. **Nœuds beta non indexés**
   - Les conditions de jointure ne sont pas extraites
   - Les nœuds beta ne sont pas dans l'index
   - **Solution**: Voir Prompt 07

### 🟡 Non-bloquants

3. **Inférence type de fait simpliste**
   - Retourne toujours "Unknown"
   - Impact: Index incomplet

4. **Reconstruction index incomplète**
   - `RebuildIndex()` ne reconstruit pas réellement
   - Impact: Ajout dynamique de règles non supporté

## 📊 État Actuel

### ✅ Fonctionnel
- Architecture d'intégration complète
- Interface `NetworkCallbacks` pour découplage
- `IntegrationHelper` coordinateur
- Stratégie hybride delta + fallback
- Tests d'intégration (84.5% couverture)
- Validation statique (vet, staticcheck)

### ⚠️  En cours (Prompt 07)
- Propagation réelle vers nœuds
- Indexation nœuds beta
- Tests end-to-end complets

## 🧪 Tests

### Exécuter les Tests

```bash
# Tests delta package
go test ./rete/delta/... -v

# Tests avec couverture
go test ./rete/delta/... -cover

# Tests d'intégration spécifiques
go test ./rete/delta/... -run TestIntegrationHelper -v
```

### Tests Disponibles

- `TestIntegrationHelper_New` - Création helper
- `TestIntegrationHelper_ProcessUpdate` - Traitement update
- `TestIntegrationHelper_ErrorCases` - Gestion erreurs
- `TestIntegrationHelper_RebuildIndex` - Reconstruction index
- `TestIntegrationHelper_Metrics` - Métriques
- `TestIntegrationHelper_Diagnostics` - Diagnostics

## 📖 Documentation

### Ordre de Lecture Recommandé

1. **`SYNTHESE_PROMPT06.md`** - Vue d'ensemble et résumé
2. **`EXECUTION_SUMMARY_PROMPT06.md`** - Détails implémentation
3. **`CODE_REVIEW_PROMPT06.md`** - Analyse qualité code
4. Code source avec GoDoc

### Documentation Code

Tous les exports sont documentés avec GoDoc:

```go
// NetworkCallbacks définit les callbacks pour interagir avec le réseau RETE.
//
// Cette interface découple le package delta du package rete principal,
// évitant ainsi les dépendances circulaires.
//
// Les implémentations de cette interface doivent être thread-safe car
// elles peuvent être appelées concurremment par plusieurs goroutines.
type NetworkCallbacks interface {
    // ...
}
```

## 🛠️  Développement

### Standards Appliqués

- `.github/prompts/common.md` - Standards projet
- `.github/prompts/review.md` - Checklist qualité
- Effective Go
- Go Code Review Comments

### Checklist Avant Commit

- [ ] `go fmt` appliqué
- [ ] `go vet` sans erreur
- [ ] `staticcheck` sans erreur
- [ ] Tests passent (couverture > 80%)
- [ ] Documentation à jour
- [ ] TODOs documentés si nécessaire

## 🔗 Références

### Fichiers Projet
- `rete/network.go` - Extension réseau RETE
- `rete/network_manager.go` - UpdateFact optimisé
- `.github/prompts/06_integration_update.md` - Prompt original

### Packages Liés
- `rete/delta/` - Package delta complet
- `rete/` - Réseau RETE

## 🐛 Problèmes Connus

### Issue #1: Propagation Non Implémentée
**Symptôme**: UpdateFact utilise toujours Retract+Insert  
**Cause**: `propagateDeltaToNode()` retourne sans action  
**Workaround**: Désactiver delta (`EnableDeltaPropagation = false`)  
**Fix**: Prompt 07 - Implémenter propagation réelle

### Issue #2: Beta Nodes Non Indexés
**Symptôme**: Changements sur champs de jointure non optimisés  
**Cause**: Extraction conditions jointure non implémentée  
**Workaround**: N/A (fallback classique automatique)  
**Fix**: Prompt 07 - Indexer nœuds beta

## 🚀 Prochaines Étapes

### Prompt 07 (Recommandé)

**Objectif**: Rendre delta propagation pleinement opérationnel

**Tâches**:
1. Implémenter propagation vers alpha nodes
2. Implémenter propagation vers beta nodes
3. Implémenter propagation vers terminal nodes
4. Indexer nœuds beta (extraction conditions)
5. Tests end-to-end
6. Benchmarks performance

**Durée estimée**: 3-4 heures  
**Difficulté**: Moyenne

## 💡 Conseils

### Performance
- La propagation delta est plus efficace avec peu de changements de champs
- Le fallback classique reste optimal pour changements massifs
- Monitorer les métriques pour ajuster la configuration

### Debugging
- Activer diagnostics: `helper.EnableDiagnostics()`
- Consulter métriques: `helper.GetMetrics()`
- Logs debug: Logger du réseau RETE

### Configuration
- Ajuster seuils dans `PropagationConfig`
- Tester différentes stratégies (Sequential, Topological, Optimized)
- Valider avec vrais cas d'usage

## 📞 Support

Questions ou problèmes? Consulter:
1. Documentation dans ce README
2. `SYNTHESE_PROMPT06.md` (FAQ implicite)
3. Code source (commentaires détaillés)
4. TODOs dans le code

---

**Dernière mise à jour**: 2026-01-02  
**Version**: Prompt 06 - Infrastructure complète  
**Statut**: ✅ Validé (propagation à implémenter)
