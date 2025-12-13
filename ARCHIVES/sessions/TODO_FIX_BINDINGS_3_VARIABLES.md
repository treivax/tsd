# TODO - Correction Bug Bindings (3 Variables)

**Priorité** : 🔴 CRITIQUE - Bloquant pour production  
**Créé** : 2025-12-12  
**Mis à jour** : 2025-12-12 (après session de debug approfondie)  
**Contexte** : Refactoring système de bindings immuable (Session 12/12)

**⚠️ MISE À JOUR IMPORTANTE** : L'investigation approfondie a démontré que **le système de bindings immuables fonctionne correctement**. Le bug n'est PAS dans l'architecture BindingChain. Voir `SESSION_DEBUG_BINDINGS_REPORT.md` et `TODO_DEBUG_E2E_BINDINGS.md` pour les détails.

---

## 🐛 Problème

**Tests échouants** : 3/83 tests E2E (96% de réussite, cible : 100%)

1. `tests/fixtures/beta/beta_join_complex.tsd` - Règle r2
2. `tests/fixtures/beta/join_multi_variable_complex.tsd` - Règle r2  
3. `tests/fixtures/integration/beta_exhaustive_coverage.tsd` - Règle r24

**Erreur observée** :
```
❌ Erreur d'exécution d'action:
Variable 'u' non trouvée dans le contexte
Variables disponibles: [p o]
```

**Comportement attendu** : Variables disponibles = [u, o, p]  
**Comportement réel** : Variables disponibles = [u, o] (le binding 'p' est perdu)

**Note** : L'erreur dans la documentation originale indiquait que 'u' était perdu avec variables=[p,o], mais les tests réels montrent que c'est 'p' qui est perdu avec variables=[u,o].

---

## 🔍 Analyse

### Règle Testée (beta_join_complex.tsd - r2)

```tsd
rule r2 : {u: User, o: Order, p: Product} / 
    u.status == "vip" AND 
    o.user_id == u.id AND 
    p.id == o.product_id AND 
    p.category == "luxury" 
    ==> vip_luxury_purchase(u.id, p.name)
```

**Action** : `vip_luxury_purchase(u.id, p.name)` → utilise 'u' et 'p'  
**Problème** : Le binding 'u' n'est pas disponible au moment de l'exécution

### Cascade de Jointures Attendue

```
TypeNode(User) ──┬──→ JoinNode1 ──┐
                 │    [u ⋈ o]     │
TypeNode(Order) ─┘                ├──→ JoinNode2 ──→ TerminalNode
                                  │    [u,o ⋈ p]
TypeNode(Product) ────────────────┘

Tokens attendus:
- JoinNode1 output : Token{Bindings: [u, o]}
- JoinNode2 output : Token{Bindings: [u, o, p]}
- TerminalNode     : Tous les bindings disponibles ✅
```

### Cascade Réelle (Bug)

```
??? ──→ JoinNode1 ──→ JoinNode2 ──→ TerminalNode
                     Token{Bindings: [p, o]}  ❌ 'u' manquant
```

**Hypothèse** : Le token qui arrive à `JoinNode2.ActivateLeft()` ne contient que [o] au lieu de [u, o]

---

## 🎯 Stratégie de Debug

### Étape 1 : Activer le Mode Debug

**Fichier** : Test ou code de construction du réseau

**Modification** :
```go
// Activer le debug sur tous les JoinNodes d'une règle 3 variables
for _, node := range network.BetaNodes {
    if jn, ok := node.(*JoinNode); ok {
        jn.Debug = true
    }
}
```

**Ou dans le test E2E** :
```go
// tests/e2e/tsd_fixtures_test.go
// Avant de soumettre les faits
network.SetDebugMode(true)
```

### Étape 2 : Tracer le Flux Complet

**Commande** :
```bash
cd /home/resinsec/dev/tsd
go test -v -tags=e2e ./tests/e2e/... -run "beta_join_complex" 2>&1 | tee debug_cascade.log
```

**Analyse du log** :
1. Chercher les lignes `🔍 [JOIN_xxx] ActivateLeft CALLED`
2. Vérifier `Token Bindings` pour chaque activation
3. Identifier à quel moment 'u' disparaît

**Pattern attendu** :
```
🔍 [JOIN_<hash1>] ActivateRight CALLED
   Fact Type: User
   RightVariables: [u]
   → Created right token with bindings: [u]

🔍 [JOIN_<hash1>] ActivateRight CALLED  
   Fact Type: Order
   RightVariables: [o]
   → Created right token with bindings: [o]

🔍 [JOIN_<hash1>] performJoinWithTokens
   Token1 Bindings: [u]
   Token2 Bindings: [o]
   → Created token: Bindings=[u, o]  ✅

🔍 [JOIN_<hash2>] ActivateLeft CALLED
   Token Bindings: [u, o]  ← VÉRIFIER ICI
   LeftVariables: [u, o]

🔍 [JOIN_<hash2>] ActivateRight CALLED
   Fact Type: Product
   RightVariables: [p]
   → Created right token with bindings: [p]

🔍 [JOIN_<hash2>] performJoinWithTokens
   Token1 Bindings: [u, o]  ← SI ICI c'est [o] → BUG TROUVÉ
   Token2 Bindings: [p]
   → Created token: Bindings=[u, o, p]  ou [o, p] si bug
```

### Étape 3 : Vérifier la Construction du Réseau

**Fichier** : `rete/builder_join_rules_cascade.go`  
**Fonction** : `connectChainToNetworkWithAlpha()`

**Ajouter des logs** :
```go
// Après chaque connexion TypeNode → JoinNode
fmt.Printf("📍 [NETWORK BUILD] Connected %s → %s.%s (pattern %d, vars: %v)\n",
    sourceNode.GetID(),
    joinNode.GetID(),
    connectionSide, // "Left" ou "Right"
    i,
    joinNode.AllVariables)
```

**Vérifications** :

Pour une règle `{u: User, o: Order, p: Product}` :

**Pattern 0** (JoinNode1) :
```
✓ TypeNode(User) → JoinNode1.Left (vars: [u])
✓ TypeNode(Order) → JoinNode1.Right (vars: [o])
✓ JoinNode1.AllVariables = [u, o]
```

**Pattern 1** (JoinNode2) :
```
✓ JoinNode1 → JoinNode2.Left (vars: [u, o])
✓ TypeNode(Product) → JoinNode2.Right (vars: [p])
✓ JoinNode2.AllVariables = [u, o, p]
```

---

## 🔧 Corrections Potentielles

### Scénario A : Connexion Incorrecte

**Si le log montre** : TypeNode connecté au mauvais côté (Left/Right)

**Fichier** : `rete/builder_join_rules_cascade.go`  
**Ligne** : ~230-260

**Correction** :
```go
// Connecter le premier JoinNode
if i == 0 {
    // Premier pattern : 2 variables
    leftType := variableTypes[0]   // u: User
    rightType := variableTypes[1]  // o: Order
    
    leftTypeNode := network.TypeNodes[leftType]
    rightTypeNode := network.TypeNodes[rightType]
    
    // IMPORTANT: Vérifier l'ordre Left/Right
    leftTypeNode.AddChild(...) → joinNode.ActivateLeft
    rightTypeNode.AddChild(...) → joinNode.ActivateRight  
} else {
    // Patterns suivants : connecter le JoinNode précédent à gauche
    previousJoinNode := chain.Nodes[i-1]
    currentRightType := variableTypes[i+1]
    
    previousJoinNode.AddChild(...) → currentJoinNode.ActivateLeft
    rightTypeNode.AddChild(...) → currentJoinNode.ActivateRight
}
```

### Scénario B : Token Mal Propagé

**Si le log montre** : Token reçu dans ActivateLeft a des bindings incomplets

**Fichier** : `rete/node_join.go`  
**Fonction** : `ActivateLeft()`

**Vérification** :
```go
func (jn *JoinNode) ActivateLeft(token *Token) error {
    // VÉRIFIER: token.Bindings contient-il TOUS les bindings attendus ?
    expectedVars := jn.LeftVariables
    actualVars := token.GetVariables()
    
    // Debug temporaire
    fmt.Printf("🔍 ActivateLeft: expected %v, got %v\n", expectedVars, actualVars)
    
    // Si actualVars ne contient pas toutes les expectedVars → BUG AMONT
    // ...
}
```

### Scénario C : performJoinWithTokens Défaillant

**Si le log montre** : Merge() ne fusionne pas correctement les bindings

**Fichier** : `rete/node_join.go`  
**Fonction** : `performJoinWithTokens()`

**Vérification** :
```go
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
    // ...
    newBindings := token1.Bindings.Merge(token2.Bindings)
    
    // Debug: vérifier que le merge est complet
    vars := newBindings.Variables()
    expectedVars := jn.AllVariables
    
    if len(vars) != len(expectedVars) {
        fmt.Printf("❌ MERGE INCOMPLETE: expected %v, got %v\n", expectedVars, vars)
    }
    // ...
}
```

---

## ✅ Checklist de Validation

Après correction :

- [ ] Activer debug sur JoinNodes
- [ ] Tracer le flux pour beta_join_complex.tsd r2
- [ ] Identifier la cause exacte du problème
- [ ] Implémenter la correction
- [ ] Vérifier que les 3 tests passent
- [ ] Lancer `make test-e2e` → 83/83 tests ✅
- [ ] Lancer `make test-complete` → PASS ✅
- [ ] Lancer `make validate` → PASS ✅
- [ ] Vérifier qu'aucune régression (autres tests)
- [ ] Documenter la cause et la solution
- [ ] Mettre à jour BINDINGS_DESIGN.md avec status ✅
- [ ] Mettre à jour CHANGELOG.md (EN COURS → COMPLETED)
- [ ] Créer commit avec message détaillé

---

## 📝 Notes

### Commandes Utiles

```bash
# Test d'un seul fichier
go test -v -tags=e2e ./tests/e2e/... -run "beta_join_complex"

# Tous les tests E2E
make test-e2e

# Validation complète
make validate

# Coverage
go test -v -cover ./rete/...

# Benchmarks
go test -bench=. -benchmem ./rete/
```

### Fichiers Clés

**Implementation** :
- `rete/binding_chain.go` - Structure immuable
- `rete/node_join.go` - Logique de jointure
- `rete/builder_join_rules_cascade.go` - Construction des cascades

**Tests** :
- `tests/fixtures/beta/beta_join_complex.tsd` - Test échouant
- `rete/node_join_cascade_test.go` - Tests unitaires cascades

**Documentation** :
- `docs/architecture/BINDINGS_DESIGN.md` - Spécification
- `docs/architecture/BINDINGS_STATUS_REPORT.md` - État actuel

---

## 🎯 Objectif Final

**Critère de succès** : 83/83 tests E2E passent (100%)

**État actuel** : 77/80 tests E2E passent (96%)

**Effort estimé** : 2-4 heures
- 1h : Debug avec traces
- 1-2h : Correction
- 1h : Validation et tests

---

## ✅ Résultats de l'Investigation (2025-12-12)

### Ce qui a été validé

1. ✅ **Architecture BindingChain** - Fonctionne parfaitement
2. ✅ **performJoinWithTokens** - Le merge est correct
3. ✅ **Propagation des tokens** - Correcte dans le test manuel
4. ✅ **Test de debug créé** - `rete/node_join_debug_test.go` PASSE

### Prochaines étapes

Le problème est ailleurs dans le système réel. Voir les fichiers suivants :
- **`SESSION_DEBUG_BINDINGS_REPORT.md`** - Rapport détaillé de l'investigation (3h de debug)
- **`TODO_DEBUG_E2E_BINDINGS.md`** - Actions prioritaires pour résoudre le bug réel

### Hypothèses principales

1. **evaluateJoinConditions** pourrait échouer pour le second JoinNode
2. **PassthroughAlpha** pourraient être mal connectés dans le système réel
3. **Ordre de soumission** des faits pourrait causer un problème
4. **Logging insuffisant** dans les tests E2E (stdout capturé)

---

**Créé par** : Session 12/12 - Documentation et Cleanup  
**Investigation approfondie** : 2025-12-12 (~3h)  
**Dernière mise à jour** : 2025-12-12  
**Priorité** : 🔴 CRITIQUE

**Statut actuel** : Architecture validée ✅, bug réel non encore localisé 🔍
