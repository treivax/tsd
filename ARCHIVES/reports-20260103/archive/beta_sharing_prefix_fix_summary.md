# Correction du Bug de Partage de Préfixes Beta

**Date**: 2025-12-13  
**Statut**: ✅ Résolu  
**Impact**: Critique - Provoquait des erreurs d'exécution sur des règles multi-variables

---

## 🐛 Problème Identifié

### Symptômes

Trois tests E2E échouaient systématiquement avec l'erreur suivante :

```
Variable 'p' not found in context
Variables available: [u o]
```

**Tests affectés** :
- `tests/fixtures/beta/beta_join_complex.tsd` (r2)
- `tests/fixtures/beta/join_multi_variable_complex.tsd` (r2)
- `tests/fixtures/integration/beta_exhaustive_coverage.tsd` (r24)

### Analyse de la Cause Racine

Le problème se situait dans le **mécanisme de partage de préfixes** (`prefix sharing`) du `BetaChainBuilder`, pas dans le système de bindings immuables (`BindingChain`) qui avait été récemment introduit.

#### Scénario du Bug

Pour deux règles avec les mêmes variables mais des conditions différentes :

```tsd
rule r1 : {u: User, o: Order, p: Product} / <conditions_r1> ==> action1(...)
rule r2 : {u: User, o: Order, p: Product} / <conditions_r2> ==> action2(...)
```

**Comportement attendu** :
- r1 crée 2 JoinNodes : `[u] ⋈ [o]` (level 0) et `[u,o] ⋈ [p]` (level 1)
- r2 crée 2 JoinNodes : `[u] ⋈ [o]` (level 0) et `[u,o] ⋈ [p]` (level 1)
- Les JoinNodes peuvent être partagés si les conditions sont identiques

**Comportement bugué** :
- r1 crée 2 JoinNodes correctement
- r2 **réutilise le préfixe de r1** (premier JoinNode) au lieu de créer le sien
- r2 ne crée qu'1 seul JoinNode : `[u,o] ⋈ [p]` (level 1)
- Résultat : r2 reçoit des tokens incomplets `[u, o]` sans `p`

#### Origine Technique

Le bug se trouvait dans la fonction `computePrefixKey()` du fichier `rete/beta_chain_optimizer.go` :

```go
// ❌ Code bugué (AVANT)
func (bcb *BetaChainBuilder) computePrefixKey(patterns []JoinPattern) string {
    key := ""
    for _, pattern := range patterns {
        key += fmt.Sprintf("%v|%v|", pattern.LeftVars, pattern.RightVars)
    }
    return key
}
```

**Problème** : La clé de cache ne contenait **pas le `ruleID`**, permettant à différentes règles de partager leurs préfixes même si leurs contextes étaient différents.

---

## ✅ Solution Implémentée

### Modifications Principales

#### 1. Ajout du `ruleID` dans la clé de préfixe

**Fichier** : `rete/beta_chain_optimizer.go`

```go
// ✅ Code corrigé (APRÈS)
func (bcb *BetaChainBuilder) computePrefixKey(patterns []JoinPattern, ruleID string) string {
    // Inclure ruleID pour éviter le partage de préfixes entre règles différentes
    key := ruleID + "::"
    for _, pattern := range patterns {
        key += fmt.Sprintf("%v|%v|", pattern.LeftVars, pattern.RightVars)
    }
    return key
}
```

**Format de clé** :
- Avant : `[u]|[o]|` (identique pour r1 et r2)
- Après : `r1::[u]|[o]|` et `r2::[u]|[o]|` (distinctes)

#### 2. Ajout du `cascadeLevel` dans la signature des JoinNodes

**Fichier** : `rete/beta_sharing_interface.go`

Ajout du champ `CascadeLevel` à `JoinNodeSignature` et `CanonicalJoinSignature` :

```go
type JoinNodeSignature struct {
    Condition    interface{}
    LeftVars     []string
    RightVars    []string
    AllVars      []string
    VarTypes     map[string]string
    CascadeLevel int  // ✅ NOUVEAU : Évite le partage entre niveaux différents
}
```

**Objectif** : Empêcher qu'un JoinNode de niveau 0 soit partagé avec un JoinNode de niveau 1, même s'ils ont les mêmes variables.

#### 3. Mise à jour de l'interface `GetOrCreateJoinNode`

```go
// Signature mise à jour
GetOrCreateJoinNode(
    condition interface{},
    leftVars []string,
    rightVars []string,
    allVars []string,
    varTypes map[string]string,
    storage Storage,
    cascadeLevel int,  // ✅ NOUVEAU paramètre
) (*JoinNode, string, bool, error)
```

### Fichiers Modifiés

| Fichier | Type | Description |
|---------|------|-------------|
| `rete/beta_chain_optimizer.go` | **FIX** | Ajout `ruleID` dans `computePrefixKey()` |
| `rete/beta_sharing_interface.go` | **FIX** | Ajout `CascadeLevel` aux signatures |
| `rete/beta_sharing.go` | **FIX** | Mise à jour `GetOrCreateJoinNode()` |
| `rete/beta_sharing_hash.go` | **FIX** | Normalisation avec `cascadeLevel` |
| `rete/beta_chain_builder_orchestration.go` | **FIX** | Passage de `cascadeLevel` et `ruleID` |
| `rete/builder_join_rules_binary_orchestration.go` | **FIX** | Passage de `cascadeLevel = 0` |
| `rete/beta_sharing_helpers.go` | **DOC** | Mise à jour commentaires |
| Tous les fichiers `*_test.go` | **TEST** | Mise à jour appels `GetOrCreateJoinNode()` |
| `rete/beta_sharing_prefix_regression_test.go` | **NEW** | Tests de régression |

---

## 🧪 Tests de Régression

Deux nouveaux tests ont été créés dans `rete/beta_sharing_prefix_regression_test.go` :

### Test 1 : `TestPrefixSharingDoesNotCrossRules`

Vérifie que le prefix sharing ne partage **pas** les préfixes entre règles différentes.

**Assertions clés** :
- ✅ Chaque règle crée 2 JoinNodes (cascade complète)
- ✅ Les clés de cache incluent le `ruleID` (`r1::[u]|[o]|` vs `r2::[u]|[o]|`)
- ✅ r1 et r2 ont des JoinNodes distincts malgré des variables identiques

### Test 2 : `TestCascadeLevelInSignature`

Vérifie que le `cascadeLevel` empêche le partage entre niveaux différents.

**Assertions clés** :
- ✅ Même condition + mêmes variables + `cascadeLevel` différent = JoinNodes différents
- ✅ Hash différent pour chaque niveau de cascade

---

## 📊 Résultats

### Tests E2E

| Test | Avant | Après |
|------|-------|-------|
| `beta_join_complex.tsd` | ❌ FAIL | ✅ PASS |
| `join_multi_variable_complex.tsd` | ❌ FAIL | ✅ PASS |
| `beta_exhaustive_coverage.tsd` | ❌ FAIL | ✅ PASS |

### Partage des JoinNodes

Le partage **légitime** des JoinNodes fonctionne toujours :

**Exemple** : Test `TestArithmeticExpressionsE2E`
- Règle 1 et Règle 3 ont les mêmes conditions → **partagent 1 JoinNode**
- Règle 2 a des conditions différentes → **JoinNode séparé**
- **Résultat** : 2 JoinNodes pour 3 règles (partage efficace ✅)

### Performance

- ✅ Tous les tests unitaires passent
- ✅ Tous les tests d'intégration passent
- ✅ Tous les tests E2E passent
- ✅ Aucune régression détectée

---

## 🎯 Principes de la Correction

### 1. Séparation des Responsabilités

- **Prefix Sharing** : Utilise `ruleID` pour éviter le partage cross-rule
- **JoinNode Sharing** : Utilise `cascadeLevel + conditions` pour partager intelligemment

### 2. Partage Conservé Quand Approprié

La correction **ne désactive pas** le partage des JoinNodes. Elle le rend plus précis :

```
✅ PARTAGE AUTORISÉ :
   - Même cascadeLevel
   - Mêmes LeftVars, RightVars
   - Même condition
   → Peut être n'importe quelle règle

❌ PARTAGE INTERDIT :
   - cascadeLevel différent
   - Variables différentes
   - Condition différente
```

### 3. Isolation des Règles via Prefix Cache

Chaque règle construit sa cascade **complète** indépendamment, empêchant :
- Les préfixes incomplets
- Les connexions incorrectes
- Les tokens avec bindings manquants

---

## 📝 Notes Techniques

### Structure du Réseau RETE

Après la correction, pour `beta_join_complex.tsd` :

```
Rule r1: {u: User, o: Order, p: Product} / <cond1> ==> action1(...)
  TypeNode(User)  ──→ PassthroughAlpha ──→ JoinNode_1 (level 0: [u] ⋈ [o])
  TypeNode(Order) ──→ PassthroughAlpha ──┘              │
                                                        ↓
  TypeNode(Product) ─→ PassthroughAlpha ──→ JoinNode_2 (level 1: [u,o] ⋈ [p])
                                                        │
                                                        ↓
                                                  r1_terminal

Rule r2: {u: User, o: Order, p: Product} / <cond2> ==> action2(...)
  TypeNode(User)  ──→ PassthroughAlpha ──→ JoinNode_3 (level 0: [u] ⋈ [o])
  TypeNode(Order) ──→ PassthroughAlpha ──┘              │
                                                        ↓
  TypeNode(Product) ─→ PassthroughAlpha ──→ JoinNode_4 (level 1: [u,o] ⋈ [p])
                                                        │
                                                        ↓
                                                  r2_terminal
```

**Note** : Si `<cond1> == <cond2>`, alors `JoinNode_1 == JoinNode_3` et `JoinNode_2 == JoinNode_4` (partage légitime).

---

## 🔍 Leçons Apprises

### 1. Importance du Contexte dans le Partage

Le partage de ressources (JoinNodes, préfixes) doit **toujours** tenir compte du contexte :
- **RuleID** pour les préfixes (évite confusion entre règles)
- **CascadeLevel** pour les joins (évite confusion entre niveaux)

### 2. Tests de Régression Cruciaux

Les tests E2E ont révélé un bug subtil qui n'apparaissait pas dans les tests unitaires. Les tests de régression ajoutés empêcheront sa réapparition.

### 3. Debugging Systématique

La méthode de debugging utilisée :
1. Instrumentation ciblée (logs stderr)
2. Test isolé reproduisant le problème
3. Dumps de structure réseau
4. Analyse de la construction du réseau

Cette approche a permis d'identifier rapidement la cause racine.

---

## ✅ Validation Finale

- [x] Tous les tests E2E passent
- [x] Tous les tests unitaires passent
- [x] Tests de régression créés
- [x] Partage des JoinNodes conservé et fonctionnel
- [x] Aucune régression de performance
- [x] Documentation mise à jour
- [x] Code respecte les standards (`.github/prompts/common.md`)

---

## 🔬 Validation du Partage Incrémental

### Test des Conditions Supplémentaires

Des tests approfondis ont été ajoutés pour valider que le partage fonctionne **correctement** même quand une règle a des conditions supplémentaires par rapport à une autre.

#### Scénario de Test

```tsd
// Règle 1: Condition de base uniquement
rule r1 : {u: User, o: Order} / u.id == o.user_id ==> basic_order(...)

// Règle 2: Condition de base + filtre alpha
rule r2 : {u: User, o: Order} / u.id == o.user_id AND o.amount > 100 ==> high_value_order(...)

// Règle 3: Condition de base + 2 filtres alpha
rule r3 : {u: User, o: Order} / u.id == o.user_id AND u.status == "vip" AND o.amount > 100 ==> vip_high_value_order(...)
```

#### Résultats Validés

✅ **Partage Optimal Confirmé**
- Les 3 règles **partagent le même JoinNode** pour la jointure `u.id == o.user_id`
- Les conditions supplémentaires sont gérées par des **AlphaNodes séparés**
- **Efficacité** : 1 JoinNode au lieu de 3 (67% d'économie)

#### Métriques de Partage

```
Total requests: 3
Shared reuses: 2
Unique creations: 1
Sharing ratio: 66.67%
```

#### Tests Créés

| Test | Objectif | Résultat |
|------|----------|----------|
| `TestJoinNodeSharingWithIncrementalConditions` | Partage avec conditions alpha | ✅ PASS |
| `TestJoinNodeSharingWithDifferentAdditionalConditions` | 3 règles, 1 JoinNode partagé | ✅ PASS |
| `TestNoSharingWhenJoinConditionsDiffer` | Pas de partage si conditions différentes | ✅ PASS |
| `join_incremental_conditions.tsd` (E2E) | Fixture complète avec 7 actions | ✅ PASS |

### Conclusion du Partage

Le partage des JoinNodes fonctionne **parfaitement** dans tous les cas :

1. ✅ **Partage activé** : Mêmes conditions de jointure → JoinNode partagé
2. ✅ **Conditions alpha** : Gérées par AlphaNodes séparés (pas de conflit)
3. ✅ **Isolation des règles** : Prefix sharing ne croise pas les règles
4. ✅ **Efficacité maximale** : Jusqu'à 67% d'économie de JoinNodes

---

**Conclusion** : Le bug de partage de préfixes beta a été résolu avec succès en incluant le `ruleID` dans la clé de cache des préfixes et le `cascadeLevel` dans la signature des JoinNodes. Le partage légitime des JoinNodes est **conservé et optimisé**, permettant à des règles avec conditions incrémentales de partager efficacement les JoinNodes communs tout en gardant leurs filtres alpha séparés. Cela améliore l'efficacité du réseau RETE tout en garantissant la correction sémantique.