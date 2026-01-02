# 📝 TODO - Package rete/delta

> **Date de création** : 2025-01-02  
> **Suite au refactoring qualité**

---

## 🔴 Problèmes Existants (Non liés au refactoring)

### 1. Test Échoué : `TestDeltaPropagator_ResetMetrics`

**Fichier** : `rete/delta/delta_propagator_test.go:354`  
**Erreur** : `classic propagation not yet implemented - requires Retract+Insert callback`

**Description** :
Le test `TestDeltaPropagator_ResetMetrics` échoue car il manque l'implémentation du callback pour la propagation classique (Retract+Insert).

**Cause** :
Le `DeltaPropagator` tente d'effectuer une propagation classique (fallback) mais le callback nécessaire n'est pas fourni dans le test.

**Solution requise** :
```go
// Dans le test, ajouter un callback pour propagation classique
propagator, _ := NewDeltaPropagatorBuilder().
    WithIndex(index).
    WithPropagateCallback(deltaCallback).
    WithClassicPropagationCallback(classicCallback). // ← Ajouter ce callback
    Build()
```

**Priorité** : 🟡 Moyenne  
**Impact** : Test échoue mais fonctionnalité fonctionne en production  
**Estimation** : 30 min

---

## 🟡 Améliorations Recommandées

### 2. Réduire Complexité du Test d'Intégration

**Fichier** : `rete/delta/integration_test.go`  
**Fonction** : `TestIndexation_IntegrationScenario`  
**Complexité cyclomatique** : 20 (> 15)

**Recommandation** :
Extraire les étapes du test en fonctions séparées :

```go
func TestIndexation_IntegrationScenario(t *testing.T) {
    t.Log("🧪 TEST INTÉGRATION COMPLÈTE - Scénario réel d'indexation")
    
    // Étape 1
    index := setupIndex(t)
    
    // Étape 2
    addProductNodes(t, index)
    
    // Étape 3
    addOrderNodes(t, index)
    
    // Étapes 4-10
    validateQueries(t, index)
    validateDelta(t, index)
    validateDiagnostics(t, index)
    validateClear(t, index)
}

func setupIndex(t *testing.T) *DependencyIndex { ... }
func addProductNodes(t *testing.T, index *DependencyIndex) { ... }
// ... etc
```

**Priorité** : 🟢 Basse  
**Impact** : Qualité du code test  
**Estimation** : 1h

### 3. Corriger Warnings Staticcheck dans Tests

**Fichier** : `rete/delta/pool_test.go:221-222`  
**Warning** : `this result of append is never used (SA4010)`

**Code actuel** :
```go
slice := make([]int, 0, 10)
append(slice, 1)  // ❌ Résultat non assigné
slice = append(slice, 2)  // ❌ Valeur précédente perdue
```

**Correction** :
```go
slice := make([]int, 0, 10)
_ = append(slice, 1)  // ✅ Intentionnel (benchmark)
// OU
slice = append(slice, 1)  // ✅ Assignation correcte
slice = append(slice, 2)
```

**Priorité** : 🟢 Basse  
**Impact** : Qualité du code test  
**Estimation** : 10 min

### 4. Améliorer Couverture Tests

**Couverture actuelle** : 74.7%  
**Objectif** : > 80%

**Zones à couvrir** :
- Nouvelles fonctions extraites (`compareSimpleTypes`, `compareNumericTypes`, etc.)
- Cas d'erreur dans helpers
- Edge cases dans comparaisons

**Priorité** : 🟡 Moyenne  
**Impact** : Qualité et confiance  
**Estimation** : 2h

---

## 🟢 Améliorations Futures (Long terme)

### 5. Optimisations Performance

**Opportunités identifiées** :

1. **Cache de comparaisons** : Désactivé par défaut
   - Évaluer impact avec workload réel
   - Benchmarker avec/sans cache
   - Ajuster TTL et taille optimale

2. **Pool d'objets** : Peut être optimisé
   - Analyser contention avec profiling
   - Ajuster taille initiale pool
   - Considérer sync.Pool vs implémentation custom

3. **Batch processing** : Sous-utilisé
   - Identifier cas d'usage optimaux
   - Benchmarker batch vs sequential
   - Documenter quand utiliser

**Priorité** : 🟢 Basse  
**Impact** : Performance marginale  
**Estimation** : 1 semaine (avec profiling complet)

### 6. Documentation Utilisateur Avancée

**À ajouter** :

1. **Guide de tuning** :
   - Quand activer cache
   - Comment ajuster epsilon pour floats
   - Stratégies de propagation (sequential vs topological)

2. **Exemples avancés** :
   - Intégration avec RETE complet
   - Cas d'usage IoT (updates fréquents)
   - E-commerce (inventaire temps réel)

3. **Troubleshooting** :
   - Diagnostics performance
   - Debug delta vs classique
   - Profiling et optimisation

**Priorité** : 🟢 Basse  
**Impact** : Adoption utilisateur  
**Estimation** : 1 semaine

---

## 📊 Statistiques TODO

| Priorité | Nombre | Estimation totale |
|----------|--------|-------------------|
| 🔴 Haute | 0 | - |
| 🟡 Moyenne | 2 | 2h30 |
| 🟢 Basse | 4 | 2 semaines |
| **Total** | **6** | **~2 semaines** |

---

## 🎯 Plan d'Action Recommandé

### Court terme (cette semaine)

1. ✅ Corriger test échoué (`TestDeltaPropagator_ResetMetrics`)
2. ✅ Corriger warnings staticcheck (`pool_test.go`)

**Temps estimé** : 1h

### Moyen terme (ce mois)

3. ✅ Réduire complexité test intégration
4. ✅ Améliorer couverture tests > 80%

**Temps estimé** : 3h

### Long terme (ce trimestre)

5. ⏳ Optimisations performance (avec profiling)
6. ⏳ Documentation utilisateur avancée

**Temps estimé** : 2 semaines

---

**Dernière mise à jour** : 2025-01-02  
**Responsable** : À assigner  
**Suivi** : Ce fichier doit être mis à jour régulièrement
