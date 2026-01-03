# 🚀 Plan d'Implémentation - Propagation Optimale (RETE-II/TREAT)

## 📋 Vue d'ensemble

Ce répertoire contient le plan d'action détaillé pour implémenter la **Solution B : Incremental Update avec Propagation Delta (RETE-II / TREAT)**.

Cette solution optimise les mises à jour de faits via l'action `Update` en ne propageant que les changements (delta) au lieu de faire un Retract+Insert complet du fait.

## 🎯 Objectif

Implémenter un mécanisme de propagation différentielle qui :
- ✅ Détecte quels champs ont changé lors d'un Update
- ✅ Indexe les nœuds RETE par champs auxquels ils sont sensibles
- ✅ Propage uniquement vers les nœuds impactés par les changements
- ✅ Optimise drastiquement les performances sur les mises à jour partielles
- ✅ Maintient la compatibilité avec l'architecture RETE existante

## 📚 Contexte Théorique

### RETE Classique (Actuel)
```
Update(fact, { field1: newValue })
  → Retract(fact)           // Retire le fait du réseau
  → Insert(fact')           // Réinsère le fait modifié
  → Propagation COMPLÈTE    // Tous les nœuds sont réévalués
```

**Problème** : Si seul `field1` change, on réévalue inutilement tous les nœuds qui dépendent d'autres champs.

### RETE-II / TREAT (Cible)
```
Update(fact, { field1: newValue })
  → DetectDelta(fact, fact')           // Identifie field1 modifié
  → FindAffectedNodes(field1)          // Trouve nœuds sensibles à field1
  → PropagateSelective(delta, nodes)   // Propage uniquement le delta
```

**Avantage** : Propagation O(nodes sensibles) au lieu de O(tous les nœuds).

## 📂 Structure du Plan

Le plan est divisé en **10 prompts séquentiels**, chacun dans son propre fichier :

| Prompt | Fichier | Description | Estimation |
|--------|---------|-------------|------------|
| 01 | `01_analyse_architecture.md` | Analyse RETE actuel + conception delta | 2-3h |
| 02 | `02_modele_donnees.md` | Structures FieldDelta, DependencyIndex | 2-3h |
| 03 | `03_indexation_dependances.md` | Index nœuds par champs | 3-4h |
| 04 | `04_detection_delta.md` | Détection des champs modifiés | 2-3h |
| 05 | `05_propagation_selective.md` | Moteur de propagation delta | 4-5h |
| 06 | `06_integration_update.md` | Intégration dans action Update | 2-3h |
| 07 | `07_tests_unitaires.md` | Tests des composants delta | 3-4h |
| 08 | `08_tests_integration.md` | Tests système complet | 3-4h |
| 09 | `09_optimisations.md` | Profiling et optimisations | 2-3h |
| 10 | `10_documentation.md` | Documentation complète | 2-3h |

**Durée totale estimée** : 25-35 heures de développement

## 🔄 Workflow d'Exécution

### Prérequis
- ✅ Branche `main` à jour
- ✅ Tous les tests passent (`make test`)
- ✅ Action `Update` actuelle fonctionnelle
- ✅ Validation complète (`make validate`)

### Exécution Séquentielle

Chaque prompt **DOIT** être exécuté dans l'ordre :

```bash
# 1. Créer une branche feature
git checkout -b feature/propagation-delta-rete-ii

# 2. Exécuter chaque prompt dans l'ordre
# Ouvrir 01_analyse_architecture.md dans Zed
# Implémenter selon les instructions
# Valider : make validate && make test

# 3. Commit après chaque prompt réussi
git add .
git commit -m "feat(rete): [Prompt 01] Analyse architecture et conception delta"

# 4. Répéter pour les prompts 02 à 10

# 5. Merge final
git checkout main
git merge feature/propagation-delta-rete-ii
```

### Règles Strictes
- ⚠️ **UN SEUL prompt à la fois** - Ne pas sauter d'étapes
- ⚠️ **Validation obligatoire** après chaque prompt : `make validate && make test`
- ⚠️ **Contexte 128k max** - Chaque prompt est auto-suffisant
- ⚠️ **Commit systématique** après succès d'un prompt
- ⚠️ **Pas de code avant prompt 02** - Prompt 01 = analyse uniquement

## 📊 Métriques de Succès

### Performance
- [ ] Update partiel (1 champ) : **> 80% plus rapide** que Retract+Insert
- [ ] Update de n champs : propagation **O(nœuds sensibles)** vs O(tous nœuds)
- [ ] Overhead index : **< 5%** mémoire supplémentaire
- [ ] Benchmark : **> 10x speedup** sur cas réalistes (grande base, updates fréquents)

### Qualité
- [ ] Couverture tests : **> 90%** pour nouveaux modules
- [ ] Tous tests existants : **100% passing**
- [ ] Zero régression fonctionnelle
- [ ] Documentation complète (GoDoc + guides)

### Compatibilité
- [ ] API publique inchangée (sauf nouvelles options configurables)
- [ ] Backward compatible avec règles existantes
- [ ] Migration transparente (activation opt-in si nécessaire)

## 🔍 Points d'Attention

### Complexité
- **Alpha Nodes** : Indexer par champs testés (conditions)
- **Beta Nodes** : Indexer par champs de jointure
- **Terminal Nodes** : Pas d'index nécessaire (propagation finale)

### Edge Cases
- Champs calculés / dérivés : propagation transitive
- Clés primaires modifiées : gérer changement d'ID interne
- Concurrence : index thread-safe (sync.RWMutex)
- Transactions : delta tracking par transaction

### Performance
- Index : utiliser structures optimisées (map + slices pré-allouées)
- Cache : éviter re-calcul de deltas identiques
- Lazy evaluation : construire index à la demande si réseau petit

## 📖 Références

### Standards Projet
- `.github/prompts/common.md` - Standards de code
- `.github/prompts/develop.md` - Workflow développement
- `docs/architecture/rete.md` - Architecture RETE actuelle

### Littérature
- Forgy, C. (1982). "Rete: A Fast Algorithm for the Many Pattern/Many Object Pattern Match Problem"
- Miranker, D. (1990). "TREAT: A New and Efficient Match Algorithm for AI Production Systems"
- Doorenbos, R. (1995). "Production Matching for Large Learning Systems" (PhD thesis)

### Code Existant
- `rete/network.go` - Structure ReteNetwork
- `rete/action_executor_evaluation.go` - Évaluation Update actuelle
- `rete/alpha_node.go` - Nœuds alpha (conditions)
- `rete/beta_node.go` - Nœuds beta (jointures)

## 🚦 Statut d'Exécution

Tracker la progression :

- [ ] **Prompt 01** - Analyse architecture ✏️
- [ ] **Prompt 02** - Modèle données
- [ ] **Prompt 03** - Indexation dépendances
- [ ] **Prompt 04** - Détection delta
- [ ] **Prompt 05** - Propagation sélective
- [ ] **Prompt 06** - Intégration Update
- [ ] **Prompt 07** - Tests unitaires
- [ ] **Prompt 08** - Tests intégration
- [ ] **Prompt 09** - Optimisations
- [ ] **Prompt 10** - Documentation

---

**Date de création** : 2025-01-02  
**Version** : 1.0  
**Auteur** : TSD Contributors  
**Licence** : MIT