# TODO - Debug Tests E2E Bindings (3 Variables)

**Priorité** : 🔴 CRITIQUE - 3 tests E2E échouent  
**Créé** : 2025-12-12  
**Contexte** : Suite à l'investigation approfondie (voir SESSION_DEBUG_BINDINGS_REPORT.md)

---

## 🎯 Objectif

Résoudre le bug où la variable 'p' (Product) est perdue dans les cascades de jointures à 3 variables.

**Tests échouants** :
- `tests/fixtures/beta/beta_join_complex.tsd` - Règle r2
- `tests/fixtures/beta/join_multi_variable_complex.tsd` - Règle r2
- `tests/fixtures/integration/beta_exhaustive_coverage.tsd` - Règle r24

**Erreur** : `Variable 'p' non trouvée dans le contexte. Variables disponibles: [u o]`

---

## ✅ Validation Préalable

**IMPORTANT** : L'investigation a démontré que :
- ✅ Le système `BindingChain` fonctionne correctement
- ✅ Le code de jointure (`performJoinWithTokens`) est correct
- ✅ La propagation des tokens est correcte
- ✅ Le test manuel (`rete/node_join_debug_test.go`) PASSE avec succès

**Conclusion** : Le problème n'est PAS dans l'architecture immuable, mais dans la configuration ou l'exécution du système réel.

---

## 📋 Actions Prioritaires

### 1. Activer le Logging pour Tests E2E (URGENT)

**Problème** : Les `fmt.Printf` ne s'affichent pas dans les tests E2E car stdout est capturé.

**Solution** : Utiliser stderr ou un fichier de log.

**Fichiers à modifier** :
- `rete/node_join.go`
- `rete/node_alpha.go`

**Code à ajouter** :

```go
// Dans node_join.go - performJoinWithTokens
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
    // DEBUG E2E - Écrire sur stderr
    fmt.Fprintf(os.Stderr, "\n🔍 [JOIN_%s] performJoinWithTokens\n", jn.ID)
    fmt.Fprintf(os.Stderr, "   Token1: %v, Token2: %v\n", token1.GetVariables(), token2.GetVariables())
    
    // ... code existant ...
    
    newBindings = token1.Bindings.Merge(token2.Bindings)
    fmt.Fprintf(os.Stderr, "   After merge: %v\n", newBindings.Variables())
    
    // ... suite du code ...
    
    if !jn.evaluateJoinConditions(newBindings) {
        fmt.Fprintf(os.Stderr, "   ❌ Join conditions FAILED\n")
        return nil
    }
    
    fmt.Fprintf(os.Stderr, "   ✅ Join conditions PASSED\n")
    fmt.Fprintf(os.Stderr, "   Final token: %v\n", newBindings.Variables())
    
    return joinedToken
}
```

**Commande de test** :
```bash
go test -v -tags=e2e ./tests/e2e -run "TestBetaFixtures/beta_join_complex" 2>&1 | tee debug_e2e.log
```

**Chercher dans les logs** :
- Est-ce que `performJoinWithTokens` est appelé pour les 2 JoinNodes ?
- Quels bindings sont présents à chaque étape ?
- Est-ce que "Join conditions FAILED" apparaît ?

### 2. Vérifier evaluateJoinConditions (URGENT)

**Hypothèse** : La fonction `evaluateJoinConditions()` pourrait retourner `false` pour le second JoinNode.

**Fichier** : `rete/node_join.go`

**Debug à ajouter** :

```go
func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool {
    fmt.Fprintf(os.Stderr, "\n🔍 [JOIN_%s] evaluateJoinConditions\n", jn.ID)
    fmt.Fprintf(os.Stderr, "   Bindings disponibles: %v\n", bindings.Variables())
    fmt.Fprintf(os.Stderr, "   JoinConditions à évaluer: %d\n", len(jn.JoinConditions))
    
    // Pour chaque condition
    for i, cond := range jn.JoinConditions {
        fmt.Fprintf(os.Stderr, "   Condition %d: %s.%s %s %s.%s\n", 
            i, cond.LeftVar, cond.LeftField, cond.Operator, cond.RightVar, cond.RightField)
        
        leftFact := bindings.Get(cond.LeftVar)
        rightFact := bindings.Get(cond.RightVar)
        
        if leftFact == nil {
            fmt.Fprintf(os.Stderr, "      ❌ LeftVar '%s' NOT FOUND in bindings\n", cond.LeftVar)
        }
        if rightFact == nil {
            fmt.Fprintf(os.Stderr, "      ❌ RightVar '%s' NOT FOUND in bindings\n", cond.RightVar)
        }
    }
    
    // ... code existant d'évaluation ...
    
    if result {
        fmt.Fprintf(os.Stderr, "   ✅ Toutes les conditions passent\n")
    } else {
        fmt.Fprintf(os.Stderr, "   ❌ Au moins une condition échoue\n")
    }
    
    return result
}
```

**Scénario de bug potentiel** :
- Si `evaluateJoinConditions()` échoue pour JoinNode2
- Aucun token `[u, o, p]` n'est créé
- Mais le token `[u, o]` du JoinNode1 pourrait quand même arriver au TerminalNode (comment ?)

### 3. Inspecter l'État du Réseau Après Construction (HAUTE)

**Créer un utilitaire de dump** :

**Fichier** : `rete/debug_utils.go` (nouveau fichier)

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "fmt"
    "os"
)

// DumpNetworkStructure affiche la structure complète du réseau pour debug
func DumpNetworkStructure(network *ReteNetwork, outputPath string) error {
    f, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer f.Close()
    
    fmt.Fprintf(f, "RETE Network Structure\n")
    fmt.Fprintf(f, "======================\n\n")
    
    // TypeNodes
    fmt.Fprintf(f, "TypeNodes: %d\n", len(network.TypeNodes))
    for name, node := range network.TypeNodes {
        fmt.Fprintf(f, "  %s (ID: %s)\n", name, node.ID)
        fmt.Fprintf(f, "    TypeName: %s\n", node.TypeName)
        fmt.Fprintf(f, "    Children: %d\n", len(node.Children))
        for i, child := range node.Children {
            fmt.Fprintf(f, "      %d: %s (type: %s)\n", i, child.GetID(), child.GetType())
        }
    }
    
    // BetaNodes (JoinNodes)
    fmt.Fprintf(f, "\nBetaNodes (JoinNodes): %d\n", len(network.BetaNodes))
    for id, node := range network.BetaNodes {
        if jn, ok := node.(*JoinNode); ok {
            fmt.Fprintf(f, "  %s\n", id)
            fmt.Fprintf(f, "    LeftVars: %v\n", jn.LeftVariables)
            fmt.Fprintf(f, "    RightVars: %v\n", jn.RightVariables)
            fmt.Fprintf(f, "    AllVars: %v\n", jn.AllVariables)
            fmt.Fprintf(f, "    VariableTypes: %v\n", jn.VariableTypes)
            fmt.Fprintf(f, "    JoinConditions: %d\n", len(jn.JoinConditions))
            for i, cond := range jn.JoinConditions {
                fmt.Fprintf(f, "      %d: %s.%s %s %s.%s\n",
                    i, cond.LeftVar, cond.LeftField, cond.Operator, cond.RightVar, cond.RightField)
            }
            fmt.Fprintf(f, "    Children: %d\n", len(jn.Children))
            for i, child := range jn.Children {
                fmt.Fprintf(f, "      %d: %s (type: %s)\n", i, child.GetID(), child.GetType())
            }
        }
    }
    
    // PassthroughAlpha
    fmt.Fprintf(f, "\nPassthroughAlpha: %d\n", len(network.PassthroughRegistry))
    for key, node := range network.PassthroughRegistry {
        fmt.Fprintf(f, "  %s\n", key)
        fmt.Fprintf(f, "    ID: %s\n", node.ID)
        fmt.Fprintf(f, "    VariableName: %s\n", node.VariableName)
        if condMap, ok := node.Condition.(map[string]interface{}); ok {
            if side, exists := condMap["side"]; exists {
                fmt.Fprintf(f, "    Side: %v\n", side)
            }
        }
        fmt.Fprintf(f, "    Children: %d\n", len(node.Children))
        for i, child := range node.Children {
            fmt.Fprintf(f, "      %d: %s (type: %s)\n", i, child.GetID(), child.GetType())
        }
    }
    
    return nil
}
```

**Utilisation dans le test** :

```go
// Dans tests/e2e ou après construction du réseau
DumpNetworkStructure(network, "/tmp/rete_network_structure.txt")
```

**Vérifier** :
- Est-ce que les PassthroughAlpha sont bien connectés aux bons JoinNodes ?
- Est-ce que JoinNode1 est bien connecté à JoinNode2 ?
- Est-ce que les `AllVariables` sont corrects pour chaque JoinNode ?

### 4. Tracer l'Ordre de Soumission des Faits (MOYENNE)

**Observation** : L'erreur se produit lors de la soumission du fait `ORD001` (Order).

**Hypothèse** : L'ordre de soumission pourrait causer un problème.

**Debug à ajouter** :

```go
// Dans rete/rete_network.go - SubmitFact
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
    fmt.Fprintf(os.Stderr, "\n📥 [NETWORK] SubmitFact: %s (Type: %s)\n", fact.ID, fact.Type)
    
    // ... code existant ...
    
    return err
}
```

**Inspecter les mémoires après chaque soumission** :

```go
// Après chaque soumission dans le test
for id, node := range network.BetaNodes {
    if jn, ok := node.(*JoinNode); ok {
        fmt.Fprintf(os.Stderr, "  JoinNode %s: Left=%d, Right=%d, Result=%d\n",
            id, len(jn.LeftMemory.Tokens), len(jn.RightMemory.Tokens), len(jn.ResultMemory.Tokens))
    }
}
```

---

## 🔬 Scénarios de Debug

### Scénario A : evaluateJoinConditions Échoue

**Si les logs montrent** :
```
🔍 [JOIN_xxx] performJoinWithTokens
   After merge: [u o p]
   ❌ Join conditions FAILED
```

**Action** : Examiner pourquoi les conditions échouent. Possibilités :
- Les valeurs des champs ne correspondent pas (ex: `u.id != o.user_id`)
- Un binding est nil alors qu'il ne devrait pas l'être
- La logique d'évaluation a un bug

### Scénario B : performJoinWithTokens Jamais Appelé pour JoinNode2

**Si les logs montrent** :
```
🔍 [JOIN_1] performJoinWithTokens  ← JoinNode1
   After merge: [u o]
   ✅ Join conditions PASSED

(pas de log pour JoinNode2)
```

**Action** : Vérifier que :
- Product est bien dans la RightMemory de JoinNode2
- Le token [u, o] arrive bien dans la LeftMemory de JoinNode2
- JoinNode2.ActivateLeft et ActivateRight sont bien appelés

### Scénario C : Mauvaise Connexion des PassthroughAlpha

**Si `network_structure.txt` montre** :
```
PassthroughAlpha: passthrough_r2_p_Product_right
    Children: 0   ← ❌ PAS D'ENFANTS !
```

**Action** : Le PassthroughAlpha n'est pas connecté au JoinNode2. Vérifier :
- `builder_join_rules_cascade.go` - ligne 274 : `passthroughAlpha.AddChild(joinNode)`
- Ou `builder_utils.go` - `ConnectTypeNodeToBetaNode`

---

## 📝 Checklist de Debug

- [ ] Ajouter logging stderr dans `performJoinWithTokens`
- [ ] Ajouter logging stderr dans `evaluateJoinConditions`
- [ ] Ajouter logging stderr dans `ActivateLeft` et `ActivateRight`
- [ ] Créer `DumpNetworkStructure` utilitaire
- [ ] Exécuter test E2E avec logs : `beta_join_complex`
- [ ] Examiner `debug_e2e.log` pour identifier le problème
- [ ] Examiner `/tmp/rete_network_structure.txt` pour vérifier les connexions
- [ ] Identifier la cause exacte du bug
- [ ] Implémenter la correction
- [ ] Vérifier que les 3 tests passent
- [ ] Retirer tous les logs de debug
- [ ] Lancer `make test-complete` → 83/83 tests ✅
- [ ] Mettre à jour la documentation

---

## 🎯 Critère de Succès

**Les 3 tests E2E doivent passer** :
```bash
go test -v -tags=e2e ./tests/e2e -run "TestBetaFixtures/beta_join_complex"
go test -v -tags=e2e ./tests/e2e -run "TestBetaFixtures/join_multi_variable_complex"
go test -v -tags=e2e ./tests/e2e -run "TestBetaFixtures/beta_exhaustive_coverage"
```

**Résultat attendu** : 83/83 tests E2E passent (100%)

---

## 📚 Références

- `SESSION_DEBUG_BINDINGS_REPORT.md` - Rapport détaillé de l'investigation
- `docs/architecture/BINDINGS_DESIGN.md` - Spécification du système immuable
- `rete/node_join_debug_test.go` - Test manuel qui FONCTIONNE

---

**Créé** : 2025-12-12  
**Priorité** : 🔴 CRITIQUE  
**Estimation** : 2-4 heures pour identifier et corriger le bug