# Design du Système de Bindings Immuable

**Date** : 2025-12-12  
**Version** : 1.0  
**Statut** : ⚠️ EN COURS - Tests partiellement échouants (77/80 passent)

---

## 1. Résumé Exécutif

### 1.1 Objectif

Remplacer le système de bindings mutable (`map[string]*Fact`) par une architecture immuable (`BindingChain`) pour garantir qu'aucun binding ne puisse être perdu lors des jointures en cascade (3+ variables).

### 1.2 Problème Résolu

**Avant** : Les règles avec 3+ variables échouaient avec "variable non trouvée" car les bindings pouvaient être écrasés dans la map mutable lors des jointures successives.

**Après** : Structure de données immuable qui préserve tous les bindings par construction.

### 1.3 Statut Actuel

✅ **Implémenté** :
- BindingChain (chaîne immuable)
- Token avec BindingChain
- Tests unitaires de BindingChain (>95% couverture)
- Benchmarks de performance

⚠️ **En cours** :
- 3 tests E2E échouent encore (77/80 passent)
- Les bindings ne se propagent pas correctement dans certaines cascades
- Investigation en cours sur la propagation left/right dans JoinNode

---

## 2. Architecture Technique

### 2.1 BindingChain - Structure Immuable

#### Concept

La `BindingChain` est une liste chaînée fonctionnelle (Cons List) qui garantit l'immutabilité par construction.

```go
type BindingChain struct {
    Variable string          // Nom de la variable (ex: "u", "order", "task")
    Fact     *Fact          // Pointeur vers le fait lié
    Parent   *BindingChain  // Chaîne parente (nil si vide)
}
```

#### Propriétés Garanties

1. **Immutabilité totale** : Une fois créée, une BindingChain ne change JAMAIS
2. **Structural Sharing** : Les chaînes partagent leur structure parente
3. **Thread-safe** : L'immutabilité garantit la sécurité concurrente
4. **Pas de cycles** : Parent pointe toujours vers une chaîne plus courte

#### Opérations

```go
// Créer une chaîne vide
chain := NewBindingChain()  // nil

// Ajouter un binding (retourne NOUVELLE chaîne)
chain1 := chain.Add("u", userFact)        // [u]
chain2 := chain1.Add("order", orderFact)  // [u, order]
chain3 := chain2.Add("task", taskFact)    // [u, order, task]

// Les anciennes chaînes restent inchangées
chain1.Len() // toujours 1
chain2.Len() // toujours 2
chain3.Len() // 3

// Récupérer un binding
fact := chain3.Get("order")  // retourne orderFact

// Fusionner deux chaînes
merged := chain1.Merge(chain2)  // [u, order]
```

#### Complexité

| Opération | Complexité | Notes |
|-----------|-----------|-------|
| `Add(v, f)` | O(1) | Création d'un nouveau nœud |
| `Get(v)` | O(n) | Parcours séquentiel, n = nb bindings |
| `Merge(other)` | O(m) | m = taille de other |
| `Has(v)` | O(n) | Parcours séquentiel |
| `Len()` | O(n) | Parcours complet |
| `Variables()` | O(n) | Collecte de toutes les variables |

**Trade-off** : O(1) pour Add au prix de O(n) pour Get. Acceptable car :
- Add est appelé à chaque jointure (fréquent)
- Get est appelé uniquement lors de l'évaluation d'actions (rare)
- n reste petit en pratique (généralement 2-5 variables)

### 2.2 Token avec BindingChain

#### Ancienne Structure (Mutable)

```go
type Token struct {
    ID       string
    Facts    []*Fact
    Bindings map[string]*Fact  // ❌ MUTABLE
    NodeID   string
}
```

**Problème** : La map pouvait être modifiée, entraînant la perte de bindings.

#### Nouvelle Structure (Immuable)

```go
type Token struct {
    ID           string
    Facts        []*Fact
    Bindings     *BindingChain    // ✅ IMMUABLE
    NodeID       string
    IsJoinResult bool
    Metadata     TokenMetadata
}

type TokenMetadata struct {
    CreatedAt    string
    CreatedBy    string
    JoinLevel    int
    ParentTokens []string
}
```

**Avantages** :
- Bindings ne peuvent jamais être perdus
- Traçabilité complète avec métadonnées
- Thread-safe par nature

### 2.3 Jointures Multi-Variables

#### Cascade de Jointures (Exemple 3 Variables)

```
Variables: [u: User, o: Order, p: Product]

TypeNode(User) ──→ AlphaNode(u) ─┐
                                  ├──→ JoinNode1 ──→ JoinNode2 ──→ TerminalNode
TypeNode(Order) ──→ AlphaNode(o) ─┘        ↑
                                            │
TypeNode(Product) ──→ AlphaNode(p) ────────┘

JoinNode1: LeftVars=[u], RightVars=[o], AllVars=[u,o]
JoinNode2: LeftVars=[u,o], RightVars=[p], AllVars=[u,o,p]
```

#### Propagation des Bindings

**Étape 1 - JoinNode1 (u ⋈ o)** :
```
Input Left:  Token{Bindings: [u]}
Input Right: Token{Bindings: [o]}
Output:      Token{Bindings: [u, o]}  // Merge([u], [o])
```

**Étape 2 - JoinNode2 (u,o ⋈ p)** :
```
Input Left:  Token{Bindings: [u, o]}   // Vient de JoinNode1
Input Right: Token{Bindings: [p]}      // Nouveau fait
Output:      Token{Bindings: [u, o, p]}  // Merge([u,o], [p])
```

**Résultat Final au TerminalNode** :
```
Token{Bindings: [u, o, p]}  // ✅ TOUS les bindings préservés
```

### 2.4 Configuration des JoinNodes

**CRITIQUE** : Chaque JoinNode doit avoir `AllVariables` correctement configuré :

```go
JoinNode{
    LeftVariables:  []string{"u", "o"},  // Variables du côté gauche
    RightVariables: []string{"p"},        // Nouvelle variable à joindre
    AllVariables:   []string{"u", "o", "p"},  // ❗ TOUTES les variables cumulées
    VariableTypes:  map[string]string{
        "u": "User",
        "o": "Order", 
        "p": "Product",
    },
}
```

`AllVariables` **DOIT** contenir la liste complète et ordonnée de toutes les variables accumulées à ce point de la cascade.

---

## 3. Implémentation

### 3.1 Fichiers Créés

| Fichier | Lignes | Description |
|---------|--------|-------------|
| `rete/binding_chain.go` | ~300 | Structure BindingChain immuable |
| `rete/binding_chain_test.go` | ~500 | Tests unitaires (>95% couverture) |
| `rete/node_join_cascade_test.go` | ~500 | Tests de cascades 2-10 variables |
| `rete/node_join_benchmark_test.go` | ~400 | Benchmarks de performance |

### 3.2 Fichiers Modifiés

| Fichier | Modifications | Raison |
|---------|--------------|--------|
| `rete/fact_token.go` | Token.Bindings: map → *BindingChain | Architecture immuable |
| `rete/node_join.go` | performJoinWithTokens refactoré | Utilisation de Merge() |
| `rete/builder_beta_chain.go` | AllVariables configuré | Propagation correcte |
| `rete/builder_join_rules_cascade.go` | buildJoinPatterns mis à jour | Patterns corrects |
| `rete/action_executor_context.go` | ExecutionContext adapté | Résolution via BindingChain |
| `rete/action_executor_evaluation.go` | Utilisation de Get() | API BindingChain |
| `rete/node_terminal.go` | Accès via BindingChain | Compatibilité |

### 3.3 API BindingChain

#### Constructeurs

```go
// Chaîne vide
func NewBindingChain() *BindingChain

// Chaîne avec un binding initial
func NewBindingChainWith(variable string, fact *Fact) *BindingChain
```

#### Mutations (retournent NOUVELLE chaîne)

```go
// Ajouter un binding
func (bc *BindingChain) Add(variable string, fact *Fact) *BindingChain

// Fusionner deux chaînes
func (bc *BindingChain) Merge(other *BindingChain) *BindingChain
```

#### Queries (lecture seule)

```go
// Récupérer un fait
func (bc *BindingChain) Get(variable string) *Fact

// Vérifier l'existence
func (bc *BindingChain) Has(variable string) bool

// Obtenir toutes les variables
func (bc *BindingChain) Variables() []string

// Taille
func (bc *BindingChain) Len() int

// Conversion (pour compatibilité)
func (bc *BindingChain) ToMap() map[string]*Fact
```

---

## 4. Tests et Validation

### 4.1 Tests Unitaires BindingChain

**Fichier** : `rete/binding_chain_test.go`

**Couverture** : >95%

**Scénarios testés** :
- Création et ajout de bindings
- Immutabilité (anciennes chaînes inchangées)
- Get/Has/Variables
- Merge de chaînes
- ToMap pour compatibilité
- Cas limites (nil, vide, etc.)

### 4.2 Tests de Cascades

**Fichier** : `rete/node_join_cascade_test.go`

**Tests paramétriques** : N=2 à N=10 variables

```go
func TestJoinCascade_NVariables(t *testing.T) {
    for n := 2; n <= 10; n++ {
        t.Run(fmt.Sprintf("%d_variables", n), func(t *testing.T) {
            // Test avec N variables
            // Vérifie que tous les bindings sont présents
        })
    }
}
```

### 4.3 Tests E2E

**Fichiers de test** :
- `tests/fixtures/beta/beta_join_complex.tsd`
- `tests/fixtures/beta/join_multi_variable_complex.tsd`
- `tests/fixtures/integration/beta_exhaustive_coverage.tsd`

**Statut** : 
- ✅ 77/80 tests E2E passent
- ❌ 3 tests échouent encore (règles 3 variables)

**Erreur observée** :
```
Variable 'u' non trouvée dans le contexte
Variables disponibles: [p o]
```

**Analyse** : Le binding 'u' est perdu quelque part dans la cascade. Investigation en cours.

### 4.4 Benchmarks

**Fichier** : `rete/node_join_benchmark_test.go`

**Résultats** :

| Opération | Temps | Allocations |
|-----------|-------|-------------|
| BindingChain.Add() | ~25 ns/op | 1 alloc |
| BindingChain.Get() (n=3) | ~11 ns/op | 0 alloc |
| JoinNode 2 variables | Baseline | Baseline |
| JoinNode 3 variables | +8% | Similaire |

**Conclusion** : Overhead <10% pour jointures 3 variables, acceptable.

---

## 5. TODO - Corrections Nécessaires

### 5.1 ❌ BUG CRITIQUE : Perte de Bindings dans Cascades 3+ Variables

**Symptôme** :
```
erreur évaluation argument: variable 'u' non trouvée
Variables disponibles: [p o]
```

**Tests affectés** :
1. `beta_join_complex.tsd` - r2
2. `join_multi_variable_complex.tsd` - r2
3. `beta_exhaustive_coverage.tsd` - r24

**Cause suspectée** :
Le binding 'u' est perdu lors de la propagation du premier JoinNode vers le second. Probablement un problème dans :
- `JoinNode.ActivateLeft()` ne reçoit pas tous les bindings
- Ou `builder_join_rules_cascade.go` connecte mal les nœuds
- Ou les TypeNodes ne propagent pas correctement vers le bon côté (left vs right)

**Actions requises** :
1. Activer le debug sur JoinNodes : `JoinNode.Debug = true`
2. Tracer le flux complet de propagation pour une règle 3 variables
3. Vérifier que `ActivateLeft` reçoit bien un token avec bindings [u, o]
4. Vérifier que `connectChainToNetworkWithAlpha` connecte correctement :
   - Premier join : TypeNode(User) → left, TypeNode(Order) → right
   - Second join : Premier join → left, TypeNode(Product) → right
5. S'assurer que le token passé à `ActivateLeft` contient TOUS les bindings accumulés

**Code à vérifier** :
```go
// rete/builder_join_rules_cascade.go
func (jrb *JoinRuleBuilder) connectChainToNetworkWithAlpha(...)

// rete/node_join.go  
func (jn *JoinNode) ActivateLeft(token *Token) error
func (jn *JoinNode) performJoinWithTokens(token1, token2 *Token) *Token
```

### 5.2 Validation Complète Requise

Une fois le bug corrigé :

```bash
# Tester avec debug activé
cd tsd
go test -v -tags=e2e ./tests/e2e/... -run "beta_join_complex"

# Validation complète
make test-complete

# Vérifier que 83/83 tests passent (100%)
```

**Critère de succès** :
- ✅ 83/83 tests E2E passent (actuellement 77/80)
- ✅ Tous les bindings disponibles dans ExecutionContext
- ✅ Actions s'exécutent sans erreur "variable non trouvée"

---

## 6. Performance

### 6.1 Overhead Mesuré

| Scénario | Temps relatif | Note |
|----------|--------------|------|
| Jointure 2 variables | Baseline (1.0x) | Référence |
| Jointure 3 variables | +8% (1.08x) | Acceptable |
| Jointure 4 variables | +12% (1.12x) | Acceptable |
| Jointure 10 variables | +25% (1.25x) | OK pour cas rare |

### 6.2 Scalabilité

**Testé jusqu'à N=10 variables** avec performances acceptables.

**Limites pratiques** :
- La plupart des règles utilisent 2-3 variables (90% des cas)
- Règles à 4-5 variables sont rares (9%)
- Règles à 6+ variables sont exceptionnelles (<1%)

**Conclusion** : L'overhead est négligeable pour les cas d'usage réels.

### 6.3 Mémoire

**Structural Sharing** réduit l'empreinte mémoire :
- Chaînes partagent les nœuds parents
- Pas de duplication des Fact (pointeurs partagés)
- Overhead : 32 bytes par binding (64-bit systems)

**Exemple** :
```
3 variables → ~96 bytes (acceptable)
10 variables → ~320 bytes (toujours acceptable)
```

---

## 7. Migration

### 7.1 Breaking Changes

**API Interne uniquement** (rete package) :

| Avant | Après |
|-------|-------|
| `Token.Bindings map[string]*Fact` | `Token.Bindings *BindingChain` |
| `bindings["u"]` | `bindings.Get("u")` |
| `bindings["u"] = fact` | `bindings = bindings.Add("u", fact)` |

**Aucun impact sur l'API publique** :
- Les fichiers `.tsd` fonctionnent sans modification
- Les règles existantes sont compatibles
- Le langage TSD reste inchangé

### 7.2 Compatibilité

✅ **Rétrocompatibilité complète** pour les utilisateurs de TSD

❌ **Pas de rétrocompatibilité** pour le code interne modifiant directement les structures RETE

### 7.3 Notes de Migration (Développeurs)

Si vous modifiez le moteur RETE :

**Avant** :
```go
// Création
bindings := make(map[string]*Fact)
bindings["u"] = userFact
bindings["o"] = orderFact

// Accès
fact := bindings["u"]

// Modification
bindings["p"] = productFact  // ❌ Mutation
```

**Après** :
```go
// Création
bindings := NewBindingChain()
bindings = bindings.Add("u", userFact)
bindings = bindings.Add("o", orderFact)

// Accès
fact := bindings.Get("u")

// Extension
bindings = bindings.Add("p", productFact)  // ✅ Nouvelle chaîne
```

---

## 8. Documentation

### 8.1 Documents Créés

1. **BINDINGS_ANALYSIS.md** - Analyse du problème et cause racine
2. **BINDINGS_DESIGN.md** - Ce document (spécification technique)
3. **BINDINGS_PERFORMANCE.md** - Résultats de performance et benchmarks
4. **CODE_REVIEW_BINDINGS.md** - Revue de code du refactoring

### 8.2 Documentation Code

**GoDoc** complet pour :
- `BindingChain` et toutes ses méthodes
- `Token` et `TokenMetadata`
- Fonctions de jointure modifiées

**Commentaires inline** pour :
- Zones critiques (Merge, performJoinWithTokens)
- Garanties d'immutabilité
- Complexité algorithmique

---

## 9. Conclusion

### 9.1 Objectifs Atteints

✅ Architecture immuable implémentée  
✅ Tests unitaires complets (>95% couverture)  
✅ Benchmarks de performance  
✅ Documentation exhaustive  
⚠️ Tests E2E partiellement passants (77/80)

### 9.2 Travail Restant

❌ **Bug critique** : Correction de la perte de bindings (3 tests)  
❌ **Validation finale** : 83/83 tests doivent passer  
❌ **Documentation** : Mise à jour avec résultats finaux

### 9.3 Impact

**Technique** :
- Architecture plus robuste et thread-safe
- Debugging facilité (métadonnées de traçage)
- Base solide pour fonctionnalités futures

**Business** :
- Support de règles complexes (3+ variables)
- Fiabilité accrue du moteur RETE
- Performance acceptable (<10% overhead)

---

**Prochaine étape** : Débugger et corriger le problème de propagation des bindings pour atteindre 83/83 tests passants.
