# 📊 RAPPORT STATISTIQUES CODE - TSD

**Date** : 26 novembre 2025  
**Commit** : `b4e9916` (2025-11-25 19:49:44)  
**Branche** : main  
**Scope** : Code manuel uniquement (hors tests, hors généré)

---

## 📈 RÉSUMÉ EXÉCUTIF

### Vue d'Ensemble
- **Lignes de code manuel** : 11,040 lignes (67.9% du projet)
- **Lignes de code généré** : 5,230 lignes (32.1% du projet)
- **Lignes de tests** : 5,241 lignes (ratio 47.5% - excellent ✅)
- **Fichiers Go fonctionnels** : 48 fichiers
- **Fonctions/Méthodes** : 463 fonctions (344 méthodes + 119 fonctions libres)

### Indicateurs Qualité

| Indicateur | Valeur | Cible | État |
|------------|--------|-------|------|
| **Lignes/Fichier** | 230 | < 400 | ✅ Excellent |
| **Lignes/Fonction** | 24 | < 50 | ✅ Excellent |
| **Complexité Max** | 37 | < 15 | 🔴 Problématique |
| **Ratio Commentaires** | 17.2% | > 15% | ✅ Bon |
| **Coverage Tests** | 32.9% | > 70% | 🔴 Insuffisant |
| **Fichiers > 800 lignes** | 2 | 0 | ⚠️ À corriger |
| **Fonctions > 100 lignes** | 6 | 0 | ⚠️ À corriger |
| **Complexité > 15** | 10 | 0 | 🔴 Urgent |

### 🎯 Priorités Urgentes

1. 🔴 **CRITIQUE** : Réduire complexité de `createSingleRule()` (37) et `main()` (36)
2. 🔴 **CRITIQUE** : Augmenter coverage de 32.9% à 70% (focus: rete, validator, nodes)
3. 🔴 **URGENT** : Refactoriser 2 fichiers > 800 lignes (constraint_pipeline, evaluator)
4. ⚠️ **IMPORTANT** : Simplifier 10 fonctions avec complexité > 15

---

## 🔥 HOTSPOTS CRITIQUES

### Top 2 Fichiers à Refactoriser URGEMMENT

1. **`rete/constraint_pipeline.go`** (1,039 lignes)
   - 🔥 Modifié 16 fois en 6 mois (HOTSPOT #1)
   - Complexité max : 37
   - 4 fonctions > 100 lignes
   - **Action** : Découper en 4 modules (parser, validator, builder, helpers)
   - **Estimation** : 6-8h

2. **`rete/evaluator.go`** (1,011 lignes)
   - 🔥 Modifié 15 fois en 6 mois (HOTSPOT #2)
   - Complexité max : 28
   - **Action** : Découper en 4 modules (alpha, map, functions, operators)
   - **Estimation** : 5-7h

---

## 📊 STATISTIQUES DÉTAILLÉES

### Répartition par Module

| Module | Lignes | Fichiers | % Total | Fonctions |
|--------|--------|----------|---------|-----------|
| **rete/** | 6,811 | 28 | 61.7% | 294 |
| **constraint/** | 3,073 | 14 | 27.8% | 133 |
| **cmd/** | 387 | 2 | 3.5% | 3 |
| **internal/** | 490 | 3 | 4.4% | 22 |
| **scripts/** | 279 | 1 | 2.5% | 11 |

### Top 10 Fonctions par Complexité

| # | Complexité | Fonction | Fichier |
|---|------------|----------|---------|
| 1 | **37** 🔴 | `createSingleRule()` | `constraint_pipeline.go` |
| 2 | **36** 🔴 | `main()` | `cmd/tsd/main.go` |
| 3 | **33** 🔴 | `evaluateJoinConditions()` | `node_join.go` |
| 4 | **30** 🔴 | `createExistsRule()` | `constraint_pipeline.go` |
| 5 | **28** 🔴 | `evaluateValueFromMap()` | `evaluator.go` |
| 6 | **24** 🔴 | `calculateAggregateForFacts()` | `node_accumulate.go` |
| 7 | **24** 🔴 | `ValidateConstraintFieldAccess()` | `constraint_utils.go` |
| 8 | **23** 🔴 | `extractAggregationInfo()` | `constraint_pipeline.go` |
| 9 | **22** 🔴 | `computeMin()` | `advanced_beta.go` |
| 10 | **22** 🔴 | `computeMax()` | `advanced_beta.go` |

**⚠️ ALERTE** : 10 fonctions avec complexité > 15 (seuil max recommandé)

---

## 🧪 COUVERTURE DE TESTS

### Coverage Globale : 32.9% 🔴

| Package | Coverage | État |
|---------|----------|------|
| `constraint` | 59.2% | ⚠️ Moyen |
| `rete` | 34.3% | 🔴 Insuffisant |
| `test/integration` | 29.4% | 🔴 Insuffisant |
| **cmd/tsd** | **0.0%** | 🔴 **CRITIQUE** |
| **cmd/universal-rete-runner** | **0.0%** | 🔴 **CRITIQUE** |
| **constraint/pkg/validator** | **0.0%** | 🔴 **CRITIQUE** |
| **rete/pkg/nodes** | **0.0%** | 🔴 **CRITIQUE** |
| **rete/pkg/domain** | **0.0%** | 🔴 **CRITIQUE** |

### Packages Sans Tests (0%) - CRITIQUE

1. **cmd/\*** - CLI non testés (impact critique utilisateurs)
2. **constraint/pkg/validator** - Validation non testée (impact sécurité)
3. **rete/pkg/nodes** - Cœur moteur RETE non testé (impact fonctionnel)

**Action urgente** : Ajouter tests pour atteindre min 60% coverage global

---

## 📈 MÉTRIQUES QUALITÉ

### Ratio Commentaires

| Module | Ratio | Évaluation |
|--------|-------|------------|
| **constraint/** | 23.2% | ✅ Excellent |
| **rete/** | 16.5% | ✅ Bon |
| **GLOBAL** | **17.2%** | ✅ Bon |
| **cmd/** | 7.1% | 🔴 Insuffisant |

### Longueur Fonctions

- **Moyenne** : 24 lignes ✅
- **< 25 lignes** : 67.4% des fonctions ✅
- **> 100 lignes** : 6 fonctions 🔴

---

## 🎯 PLAN D'ACTION

### Sprint 1 (Semaine 1-2) - CRITIQUE

| # | Action | Effort | Impact |
|---|--------|--------|--------|
| 1 | Refactoriser `constraint_pipeline.go` | 6-8h | 🔥🔥🔥 |
| 2 | Refactoriser `evaluator.go` | 5-7h | 🔥🔥🔥 |
| 3 | Ajouter tests `cmd/*` | 6-8h | 🔥🔥🔥 |
| 4 | Ajouter tests `validator` | 4-6h | 🔥🔥🔥 |
| 5 | Ajouter tests `rete/pkg/nodes` | 8-10h | 🔥🔥🔥 |

**Total Sprint 1** : 29-39h (~1-2 semaines pour 2 devs)

### Sprint 2 (Semaine 3-4) - IMPORTANT

| # | Action | Effort |
|---|--------|--------|
| 6 | Simplifier top 4 fonctions complexes | 9-13h |
| 7 | Augmenter coverage rete/ (34% → 55%) | 8-10h |
| 8 | Augmenter coverage integration/ (29% → 50%) | 4-5h |
| 9 | Améliorer doc CMD | 2-3h |

**Total Sprint 2** : 23-31h

---

## ✅ CONCLUSION

### Score Global : 6.5/10 ⚠️

**Points Forts** :
- ✅ Bonne taille globale (11K lignes)
- ✅ Fonctions courtes (moy 24 lignes)
- ✅ Bonne documentation (17.2%)
- ✅ Excellent ratio tests/code (47.5%)

**Points Faibles** :
- 🔴 2 hotspots critiques (1,000+ lignes)
- 🔴 10 fonctions complexité > 15 (max 37)
- 🔴 Coverage 32.9% (cible 70%)
- 🔴 5 packages critiques à 0% coverage

**Actions Immédiates** :
1. Refactoriser hotspots (constraint_pipeline + evaluator)
2. Ajouter tests pour packages critiques (cmd, validator, nodes)
3. Simplifier fonctions trop complexes

---

**Prochaine analyse** : Dans 1 mois (après refactoring hotspots)  
**Généré avec** : prompt `stats-code.md` v2.0
