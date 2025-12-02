# Fix: Bug de multiplication des tokens avec partage beta obligatoire

**Date**: 2025-12-02  
**Statut**: ✅ RÉSOLU  
**Impact**: Fonctionnel - Propagation incorrecte des tokens dans JoinNodes partagés

---

## 📋 Résumé

Après avoir rendu le partage beta (JoinNodes) obligatoire et systématique, un bug de multiplication des tokens a été découvert dans le test E2E `TestArithmeticExpressionsE2E`. Les règles recevaient 3× le nombre attendu de tokens en raison d'une propagation incorrecte dans les JoinNodes partagés.

---

## 🐛 Symptômes observés

### Test E2E: `TestArithmeticExpressionsE2E`

**Attendu**:
- Règle 1 (`calcul_facture_base`): 3 tokens
- Règle 2 (`calcul_facture_speciale`): 0 tokens  
- Règle 3 (`calcul_facture_premium`): 3 tokens
- **Total**: 6 tokens

**Obtenu** (avant fix):
- Règle 1: 27 tokens (9× attendu)
- Règle 2: 27 tokens (devrait être 0)
- Règle 3: 27 tokens (9× attendu)
- **Total**: 81 tokens

**Progression du fix**:
1. Après fix partiel 1: 18 tokens par règle (6× attendu)
2. Après fix partiel 2: 9 tokens par règle (3× attendu)
3. Après fix final: 3 tokens pour R1 et R3, 0 pour R2 ✅

---

## 🔍 Analyse de la cause racine

### Problème 1: Hash du JoinNode n'incluait pas les conditions alpha

**Code initial** (`builder_join_rules.go` ligne 143):
```go
// Reconstruct beta-only condition for JoinNode
var joinCondition map[string]interface{}
if len(betaConditions) > 0 {
    joinCondition = splitter.ReconstructBetaCondition(betaConditions)
}

// Create JoinNode with beta conditions only
node, hash, shared, err := network.BetaSharingRegistry.GetOrCreateJoinNode(
    joinCondition,  // ❌ Ne contient que les conditions beta
    leftVars, rightVars, allVars, varTypes, storage,
)
```

**Conséquence**: 
- Les règles 1, 2 et 3 avaient toutes la même condition beta: `c.produit_id == p.id`
- Mais des conditions alpha différentes:
  - R1: `c.qte * 23 - 10 + c.remise * 43 > 0`
  - R2: `c.qte * 23 - 10 + c.remise * 43 < 0` 
  - R3: `c.qte * 23 - 10 + c.remise * 43 > 0` (identique à R1)
- Le hash étant calculé uniquement sur les conditions beta, **les 3 règles partageaient le même JoinNode** alors que seules R1 et R3 auraient dû le partager

**Résultat**: 1 seul JoinNode au lieu de 2 → tous les tokens propagés à toutes les règles

### Problème 2: Reconnexion des inputs pour JoinNodes partagés

**Code initial** (ligne 201-240):
```go
// STEP 5: Connect the network correctly
for i, varName := range variableNames {
    // ... création des passthroughs ...
    passthroughAlpha.AddChild(joinNode)  // ❌ Toujours connecté, même si partagé
}
```

**Conséquence**:
- Quand la règle 3 réutilisait le JoinNode de R1, elle créait ses propres passthroughs
- Ces nouveaux passthroughs étaient **aussi connectés au JoinNode partagé**
- Le JoinNode recevait donc les tokens **2 fois** (une fois de R1, une fois de R3)
- Avec 3 produits LEFT et 3 commandes RIGHT, au lieu de 3 jointures, on avait:
  - 3 produits × 2 (dupliqués) = 6 tokens LEFT
  - 3 commandes × 2 (dupliqués) = 6 tokens RIGHT  
  - Jointures: 3 × 2 × 3 = 18 tokens par règle

### Problème 3: Multiple TerminalNodes connectés au même JoinNode

**Code initial** (ligne 174):
```go
joinNode.AddChild(terminalNode)  // ❌ Toujours ajouté, même si JoinNode partagé
```

**Conséquence**:
- Le JoinNode partagé avait **2 TerminalNodes** comme enfants (R1 et R3)
- Chaque token produit était propagé aux **2 TerminalNodes**
- Avec 3 tokens produits, chaque TerminalNode recevait 3 tokens
- Mais comme le JoinNode recevait des inputs dupliqués (problème 2), il produisait 9 tokens
- Résultat: 9 tokens × 2 TerminalNodes, mais chaque TerminalNode "voyait" 9 tokens

### Problème 4: JoinConditions non extraites de la condition composite

**Code initial** (`node_join.go` ligne 54):
```go
JoinConditions: extractJoinConditions(condition),  // ❌ condition est composite
```

**Conséquence**:
- La condition passée était maintenant `{"beta": {...}, "alpha": {...}}`
- `extractJoinConditions` ne savait pas extraire depuis ce format
- Résultat: `JoinConditions = []` (liste vide)
- Sans JoinConditions, **toutes les combinaisons LEFT × RIGHT étaient acceptées**
- 3 LEFT × 3 RIGHT = 9 jointures au lieu de 3 (seulement celles qui matchent)

---

## ✅ Solutions implémentées

### Fix 1: Inclure les conditions alpha dans le hash du JoinNode

**Fichier**: `rete/builder_join_rules.go` (lignes 146-162)

```go
// STEP 3b: Build composite condition including alpha conditions for proper sharing
// The JoinNode hash must include alpha conditions to prevent incorrect sharing
// between rules with same beta but different alpha conditions
compositeCondition := map[string]interface{}{
    "beta": joinCondition,
}

// Include alpha conditions in the composite to ensure proper hash differentiation
if len(alphaConditions) > 0 {
    alphaCondMap := make(map[string]interface{})
    for _, alphaCond := range alphaConditions {
        // Use variable name as key and condition as value
        varKey := alphaCond.Variable
        alphaCondMap[varKey] = alphaCond.Condition
    }
    compositeCondition["alpha"] = alphaCondMap
}

// Create JoinNode with composite condition (beta + alpha)
node, hash, shared, err := network.BetaSharingRegistry.GetOrCreateJoinNode(
    compositeCondition,  // ✅ Inclut beta ET alpha
    leftVars, rightVars, allVars, varTypes, storage,
)
```

**Résultat**: Les règles avec conditions alpha différentes obtiennent des hashs différents et ne partagent plus le même JoinNode incorrectement.

### Fix 2: Extraire la partie beta avant d'extraire les JoinConditions

**Fichier**: `rete/node_join.go` (lignes 41-48)

```go
// Extract beta condition from composite condition if present
// Composite conditions are in format: {"beta": ..., "alpha": ...}
conditionForExtraction := condition
if betaCond, hasBeta := condition["beta"]; hasBeta {
    if betaMap, ok := betaCond.(map[string]interface{}); ok {
        conditionForExtraction = betaMap
    }
}

return &JoinNode{
    // ...
    JoinConditions: extractJoinConditions(conditionForExtraction),  // ✅ Extrait depuis beta
}
```

**Résultat**: Les JoinConditions sont correctement extraites, permettant le filtrage des jointures (3 au lieu de 9).

### Fix 3: Gérer les conditions composites dans l'évaluation

**Fichier**: `rete/node_join.go` (lignes 268-277)

```go
// Unwrap composite condition (beta + alpha) if present
actualCondition := jn.Condition
if betaCond, isBeta := jn.Condition["beta"]; isBeta {
    // This is a composite condition from beta sharing with alpha conditions
    // Extract only the beta part for join evaluation
    if betaMap, ok := betaCond.(map[string]interface{}); ok {
        actualCondition = betaMap
    }
}
```

**Résultat**: L'évaluation des jointures utilise la partie beta correctement.

### Fix 4: Skip la reconnexion des inputs pour JoinNodes partagés

**Fichier**: `rete/builder_join_rules.go` (lignes 210-252)

```go
// STEP 5: Connect the network correctly
// IMPORTANT: Skip this step if JoinNode was shared - inputs are already connected
if wasShared {
    fmt.Printf("   ⏭️  Skipping input reconnection for shared JoinNode %s\n", joinNode.ID)
}

if !wasShared {
    for i, varName := range variableNames {
        // ... connexion des passthroughs ...
    }
}
```

**Résultat**: Les inputs ne sont connectés qu'une seule fois, évitant la duplication des tokens.

### Fix 5: Utiliser RuleRouterNode pour les TerminalNodes

**Fichier**: `rete/node_rule_router.go` (nouveau fichier)

```go
// RuleRouterNode is an intermediate node between a shared JoinNode and TerminalNodes
type RuleRouterNode struct {
    BaseNode
    RuleID       string
    JoinNodeID   string
    TerminalNode *TerminalNode
}

func (rrn *RuleRouterNode) ActivateLeft(token *Token) error {
    // Route the token to the terminal node
    if rrn.TerminalNode != nil {
        return rrn.TerminalNode.ActivateLeft(token)
    }
    return rrn.PropagateToChildren(nil, token)
}
```

**Fichier**: `rete/builder_join_rules.go` (lignes 194-204)

```go
// STEP 4b: Connect terminal node properly based on sharing status
if wasShared {
    // JoinNode is shared - use RuleRouterNode to avoid token duplication
    router := NewRuleRouterNode(ruleID, joinNode.ID, jrb.utils.storage)
    router.SetTerminalNode(terminalNode)
    joinNode.AddChild(router)
} else {
    // JoinNode is new - connect terminal directly
    joinNode.AddChild(terminalNode)
}
```

**Résultat**: 
- Architecture: `SharedJoinNode -> RuleRouterNode (R1) -> TerminalNode (R1)`
- Architecture: `SharedJoinNode -> TerminalNode (R1)` (première règle, connexion directe)
- Chaque règle route ses tokens indépendamment

**Note**: Le RuleRouterNode est créé uniquement pour les règles qui réutilisent un JoinNode partagé. La première règle qui crée le JoinNode connecte son TerminalNode directement.

### Fix 6: Compatibilité des tests avec le nouveau format de clé

**Fichier**: `rete/builder_join_rules.go` (lignes 209-212)

```go
// Store the JoinNode in the network's BetaNodes
network.BetaNodes[joinNode.ID] = joinNode

// Also store with legacy key format for test compatibility
legacyKey := fmt.Sprintf("%s_join", ruleID)
network.BetaNodes[legacyKey] = joinNode
```

**Fichiers de tests mis à jour**:
- `rete/action_arithmetic_e2e_test.go`: Compter les JoinNodes uniques par ID
- `rete/builder_join_rules_test.go`: Ajuster les attentes pour les entrées dupliquées

---

## 📊 Résultats

### Architecture finale du test E2E

```
TypeNode[Produit] ─┬─> PassthroughAlpha[R1_left] ─┐
                   ├─> PassthroughAlpha[R2_left] ─┼─> JoinNode[R2]
                   └─> PassthroughAlpha[R3_left] ─┘
                   
TypeNode[Commande] ─┬─> AlphaFilter1(> 0) ─> PassthroughAlpha[R1_right] ─┐
                    ├─> AlphaFilter2(< 0) ─> PassthroughAlpha[R2_right] ─┼─> JoinNode[R2]
                    └─> AlphaFilter3(> 0) ─> PassthroughAlpha[R3_right] ─┘

JoinNode[R1_R3_shared] ─┬─> TerminalNode[R1]
                        └─> RuleRouterNode[R3] ─> TerminalNode[R3]

JoinNode[R2] ─> TerminalNode[R2]
```

### Métriques

- **JoinNodes créés**: 2 (au lieu de 3 sans partage, ou 1 avec bug)
- **Tokens générés**: 6 (correct)
  - Règle 1: 3 tokens ✅
  - Règle 2: 0 tokens ✅
  - Règle 3: 3 tokens ✅
- **Partage beta**: ✅ Règles 1 et 3 partagent correctement le JoinNode
- **AlphaNodes partagés**: 5 (décomposition arithmétique réutilisée)

### Tests

Tous les tests passent:
```bash
cd rete && go test
# PASS
# ok  	github.com/treivax/tsd/rete	1.601s
```

---

## 🎯 Leçons apprises

1. **Le hash d'un JoinNode doit inclure TOUTES les conditions** (alpha et beta) qui affectent les tokens qui lui arrivent, pas seulement les conditions beta.

2. **Les inputs d'un nœud partagé ne doivent être connectés qu'une seule fois**, lors de la création initiale. Les règles suivantes réutilisent les mêmes connexions.

3. **Les nœuds intermédiaires (RuleRouterNode) sont nécessaires** pour router correctement les tokens d'un nœud partagé vers les différentes règles.

4. **L'extraction de conditions depuis des formats composites** doit être robuste et gérer les différentes structures possibles.

5. **Les tests doivent vérifier les comportements fonctionnels** (nombre de tokens corrects) et non seulement la structure du réseau.

---

## 📚 Fichiers modifiés

### Core fixes
- `rete/builder_join_rules.go` - Hash composite + skip reconnexion + RuleRouterNode
- `rete/node_join.go` - Extraction beta + évaluation composite
- `rete/node_rule_router.go` - **Nouveau** - Routage des tokens

### Tests
- `rete/action_arithmetic_e2e_test.go` - Compter JoinNodes uniques
- `rete/builder_join_rules_test.go` - Ajuster attentes

### Documentation
- `docs/BETA_SHARING_MANDATORY.md` - Mis à jour avec section bug fix
- `docs/BETA_SHARING_TOKEN_MULTIPLICATION_FIX.md` - **Ce document**

---

## 🔗 Références

- Test de référence: `TestArithmeticExpressionsE2E` (`rete/action_arithmetic_e2e_test.go`)
- Thread de conversation: `RETE Arithmetic Decomposition Metrics`
- Documentation partage beta: `docs/BETA_SHARING_MANDATORY.md`
- Spécification décomposition: `rete/ARITHMETIC_DECOMPOSITION_SPEC.md`

---

**Auteur**: Assistant IA  
**Validé par**: Tests automatisés  
**Statut**: ✅ Production Ready