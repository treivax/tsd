# Prompt 02 : Design du Système de Bindings Immuable

**Session** : 2/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Avoir complété Prompt 01 et lu `BINDINGS_ANALYSIS.md`

---

## 🎯 Objectif de cette Session

Concevoir l'architecture complète du nouveau système de bindings immuable en :
1. Spécifiant les structures de données (BindingChain, Token)
2. Définissant les interfaces et contrats
3. Planifiant la stratégie de migration (remplacement, pas cohabitation)
4. Établissant le plan de test détaillé

**Livrable final** : `tsd/docs/architecture/BINDINGS_DESIGN.md` (800-1200 lignes)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Analyser les Besoins (30 min)

#### 1.1 Extraire les exigences de BINDINGS_ANALYSIS.md

**Lire** : `tsd/docs/architecture/BINDINGS_ANALYSIS.md` (produit au Prompt 01)

**Questions à répondre** :
1. Quelle est la cause racine exacte du problème ?
2. Quelles sont les contraintes identifiées ?
3. Quels sont les points critiques à ne pas casser ?
4. Quelles opportunités d'amélioration ont été identifiées ?

**Documenter dans** : Section "1. Analyse des Besoins" de BINDINGS_DESIGN.md

---

#### 1.2 Définir les exigences fonctionnelles

**Exigences obligatoires** :
- [ ] **REQ-1** : Un token doit porter TOUS les bindings de sa généalogie
- [ ] **REQ-2** : Les bindings ne peuvent JAMAIS être perdus une fois créés
- [ ] **REQ-3** : Support de N variables (N ≥ 2, sans limite arbitraire)
- [ ] **REQ-4** : Jointure de deux tokens produit un token avec tous les bindings combinés
- [ ] **REQ-5** : Les actions peuvent accéder à toutes les variables déclarées dans la règle

**Exigences non-fonctionnelles** :
- [ ] **NFR-1** : Performance : overhead < 10% pour jointures 2 variables
- [ ] **NFR-2** : Scalabilité : jusqu'à N=10 variables sans dégradation majeure
- [ ] **NFR-3** : Thread-safety : les tokens doivent être thread-safe
- [ ] **NFR-4** : Debuggabilité : traçabilité complète de la chaîne de bindings
- [ ] **NFR-5** : Maintenabilité : code clair et testable

**Documenter dans** : Section "2. Exigences"

---

### Tâche 2 : Concevoir BindingChain (45 min)

#### 2.1 Spécification de la structure

**Structure de base** :
```go
// BindingChain représente une chaîne immuable de bindings variable → fact
// Utilise le pattern "Cons list" pour le partage de structure
type BindingChain struct {
    Variable string          // Nom de la variable (ex: "u", "order", "task")
    Fact     *Fact          // Fait lié à cette variable
    Parent   *BindingChain  // Chaîne parente (nil si racine/vide)
}
```

**Invariants à garantir** :
1. Une BindingChain est **immuable** : une fois créée, elle ne change jamais
2. `Add()` retourne une **nouvelle** chaîne, ne modifie pas l'existante
3. La racine (chaîne vide) est représentée par `Parent == nil`
4. Pas de cycles : `Parent` pointe toujours vers une chaîne plus courte

**Documenter dans** : Section "3. BindingChain - Spécification"

---

#### 2.2 API de BindingChain

**Constructeur** :
```go
// NewBindingChain crée une chaîne vide
func NewBindingChain() *BindingChain

// NewBindingChainWith crée une chaîne avec un binding initial
func NewBindingChainWith(variable string, fact *Fact) *BindingChain
```

**Opérations de lecture** (ne modifient pas la chaîne) :
```go
// Get retourne le fait associé à une variable, ou nil si non trouvé
func (bc *BindingChain) Get(variable string) *Fact

// Has vérifie si une variable existe dans la chaîne
func (bc *BindingChain) Has(variable string) bool

// Len retourne le nombre de bindings dans la chaîne
func (bc *BindingChain) Len() int

// Variables retourne la liste des variables (dans l'ordre d'ajout)
func (bc *BindingChain) Variables() []string

// ToMap convertit la chaîne en map (pour compatibilité/debug)
func (bc *BindingChain) ToMap() map[string]*Fact
```

**Opérations de construction** (retournent une nouvelle chaîne) :
```go
// Add ajoute un binding et retourne une NOUVELLE chaîne
// L'ancienne chaîne reste inchangée (immutabilité)
func (bc *BindingChain) Add(variable string, fact *Fact) *BindingChain

// Merge combine deux chaînes (retourne nouvelle chaîne)
// En cas de conflit, priorité à 'other'
func (bc *BindingChain) Merge(other *BindingChain) *BindingChain
```

**Opérations de debug** :
```go
// String retourne une représentation textuelle pour debug
func (bc *BindingChain) String() string

// Chain retourne la liste des variables depuis la racine (pour traçage)
func (bc *BindingChain) Chain() []string
```

**Documenter dans** : Section "3.2. API de BindingChain"

---

#### 2.3 Complexité algorithmique

**Analyser la complexité** :

| Opération | Complexité Temporelle | Complexité Spatiale | Notes |
|-----------|----------------------|---------------------|-------|
| `Add(v, f)` | O(1) | O(1) | Création d'un nœud |
| `Get(v)` | O(n) | O(1) | n = nombre de bindings |
| `Has(v)` | O(n) | O(1) | Parcours linéaire |
| `Len()` | O(n) | O(1) | Parcours pour compter |
| `Variables()` | O(n) | O(n) | Allocation d'une slice |
| `ToMap()` | O(n) | O(n) | Création d'une map |
| `Merge(other)` | O(m) | O(m) | m = taille de other |

**Optimisations possibles** :
1. **Cache de longueur** : Stocker `length int` dans la structure
2. **Cache de variables** : Stocker `[]string` (calculé à la demande, mis en cache)
3. **Index** : Pour n > seuil (ex: 5), créer un index map en parallèle

**Décision** : Commencer simple (sans cache), optimiser si nécessaire après benchmarks.

**Documenter dans** : Section "3.3. Complexité et Performance"

---

### Tâche 3 : Concevoir Token Immuable (45 min)

#### 3.1 Nouvelle structure Token

**Spécification** :
```go
// Token représente un ensemble de faits liés par des bindings
// Version immuable avec BindingChain
type Token struct {
    ID       string          // Identifiant unique du token
    Facts    []*Fact        // Liste des faits (ordre d'ajout)
    Bindings *BindingChain  // Chaîne immuable de bindings
    NodeID   string          // ID du nœud qui a créé ce token
    Metadata TokenMetadata   // Métadonnées pour debugging
}

// TokenMetadata contient des informations de traçage
type TokenMetadata struct {
    CreatedAt    time.Time   // Date de création
    CreatedBy    string      // ID du nœud créateur
    JoinLevel    int         // Profondeur de cascade (0 = fact initial)
    ParentTokens []string    // IDs des tokens parents (pour traçage)
}
```

**Comparaison avec l'ancienne structure** :
```go
// ANCIEN (à remplacer)
type Token struct {
    ID       string
    Facts    []*Fact
    Bindings map[string]*Fact  // ❌ Mutable, peut perdre des bindings
    NodeID   string
}

// NOUVEAU
type Token struct {
    ID       string
    Facts    []*Fact
    Bindings *BindingChain     // ✅ Immuable, bindings garantis
    NodeID   string
    Metadata TokenMetadata     // ✅ Debugging amélioré
}
```

**Documenter dans** : Section "4. Token - Nouvelle Structure"

---

#### 3.2 API de Token

**Constructeurs** :
```go
// NewToken crée un token vide
func NewToken(nodeID string) *Token

// NewTokenWithFact crée un token avec un fait initial
func NewTokenWithFact(fact *Fact, variable string, nodeID string) *Token

// NewTokenWithBinding crée un token avec des bindings existants
func NewTokenWithBinding(bindings *BindingChain, facts []*Fact, nodeID string) *Token
```

**Opérations** :
```go
// GetBinding retourne le fait lié à une variable
func (t *Token) GetBinding(variable string) *Fact

// HasBinding vérifie si une variable est liée
func (t *Token) HasBinding(variable string) bool

// GetVariables retourne toutes les variables liées
func (t *Token) GetVariables() []string

// Clone crée une copie du token (pour les cas où nécessaire)
func (t *Token) Clone() *Token
```

**Documenter dans** : Section "4.2. API de Token"

---

#### 3.3 Impact sur le code existant

**Fichiers utilisant Token (à modifier)** :
1. `rete/fact_token.go` - Définition de Token (remplacement complet)
2. `rete/node_join.go` - JoinNode utilise Token.Bindings
3. `rete/node_terminal.go` - TerminalNode lit Token.Bindings
4. `rete/action_executor_context.go` - ExecutionContext lit Token.Bindings
5. `rete/node_alpha.go` - AlphaNode peut créer des tokens
6. `rete/network.go` - Network manipule des tokens
7. Tous les tests utilisant Token

**Stratégie de migration** :
1. **Phase 1** : Créer BindingChain (nouveau fichier)
2. **Phase 2** : Modifier Token dans `fact_token.go` (remplacement direct)
3. **Phase 3** : Fixer toutes les erreurs de compilation (migration forcée)
4. **Phase 4** : Adapter les tests

**Documenter dans** : Section "4.3. Stratégie de Migration"

---

### Tâche 4 : Concevoir les Modifications de JoinNode (40 min)

#### 4.1 Nouvelle logique de performJoinWithTokens

**Ancienne implémentation (problématique)** :
```go
func (jn *JoinNode) performJoinWithTokens(token1, token2 *Token) *Token {
    combinedBindings := make(map[string]*Fact)
    
    // Copie token1.Bindings
    for k, v := range token1.Bindings {
        combinedBindings[k] = v
    }
    
    // Copie token2.Bindings
    for k, v := range token2.Bindings {
        combinedBindings[k] = v  // ❌ Peut écraser ? Peut ne pas copier tous ?
    }
    
    // Création du token
    joinedToken := &Token{
        Bindings: combinedBindings,  // ❌ Nouveau map, perte possible
        // ...
    }
    
    return joinedToken
}
```

**Nouvelle implémentation (immuable)** :
```go
func (jn *JoinNode) performJoinWithTokens(token1, token2 *Token) *Token {
    // Composer les chaînes : partir de token1, ajouter token2
    newBindings := token1.Bindings
    
    // Ajouter tous les bindings de token2
    for _, variable := range token2.Bindings.Variables() {
        fact := token2.Bindings.Get(variable)
        newBindings = newBindings.Add(variable, fact)
    }
    
    // Vérifier les conditions de jointure
    if !jn.evaluateJoinConditions(newBindings) {
        return nil
    }
    
    // Créer le token joint
    joinedToken := &Token{
        ID:       generateTokenID(),
        Facts:    append(token1.Facts, token2.Facts...),
        Bindings: newBindings,  // ✅ Chaîne complète garantie
        NodeID:   jn.ID,
        Metadata: TokenMetadata{
            CreatedAt:    time.Now(),
            CreatedBy:    jn.ID,
            JoinLevel:    max(token1.Metadata.JoinLevel, token2.Metadata.JoinLevel) + 1,
            ParentTokens: []string{token1.ID, token2.ID},
        },
    }
    
    return joinedToken
}
```

**Documenter dans** : Section "5. JoinNode - Nouvelle Logique"

---

#### 4.2 Adaptation de evaluateJoinConditions

**Signature actuelle** :
```go
func (jn *JoinNode) evaluateJoinConditions(bindings map[string]*Fact) bool
```

**Nouvelle signature** :
```go
func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool
```

**Changements nécessaires** :
- Remplacer `bindings[variable]` par `bindings.Get(variable)`
- Adapter les accès aux faits dans les conditions

**Documenter dans** : Section "5.2. Évaluation des Conditions"

---

### Tâche 5 : Concevoir les Modifications de ExecutionContext (30 min)

#### 5.1 Nouvelle structure ExecutionContext

**Ancienne structure** :
```go
type ExecutionContext struct {
    varCache map[string]*Fact  // ❌ Référence directe au map
    // ...
}
```

**Nouvelle structure** :
```go
type ExecutionContext struct {
    bindings *BindingChain     // ✅ Référence à la chaîne immuable
    // ...
}
```

**Documenter dans** : Section "6. ExecutionContext"

---

#### 5.2 Résolution de variables améliorée

**Ancienne implémentation** :
```go
func (ctx *ExecutionContext) resolveVariable(name string) (*Fact, error) {
    if fact, ok := ctx.varCache[name]; ok {
        return fact, nil
    }
    return nil, fmt.Errorf("variable '%s' non trouvée", name)
}
```

**Nouvelle implémentation** :
```go
func (ctx *ExecutionContext) resolveVariable(name string) (*Fact, error) {
    if ctx.bindings.Has(name) {
        return ctx.bindings.Get(name), nil
    }
    
    // Message d'erreur amélioré
    available := ctx.bindings.Variables()
    return nil, fmt.Errorf(
        "variable '%s' non trouvée (variables disponibles: %v)", 
        name, available,
    )
}
```

**Documenter dans** : Section "6.2. Résolution de Variables"

---

### Tâche 6 : Plan de Test Détaillé (40 min)

#### 6.1 Tests unitaires pour BindingChain

**Fichier** : `rete/binding_chain_test.go`

**Tests à implémenter** :
```go
// Test de base
func TestBindingChain_CreateEmpty(t *testing.T)
func TestBindingChain_CreateWithBinding(t *testing.T)

// Test d'ajout
func TestBindingChain_Add_Single(t *testing.T)
func TestBindingChain_Add_Multiple(t *testing.T)
func TestBindingChain_Add_Preserves_Parent(t *testing.T)  // Immutabilité

// Test de lecture
func TestBindingChain_Get_Existing(t *testing.T)
func TestBindingChain_Get_NotFound(t *testing.T)
func TestBindingChain_Has(t *testing.T)
func TestBindingChain_Len(t *testing.T)
func TestBindingChain_Variables(t *testing.T)

// Test de conversion
func TestBindingChain_ToMap(t *testing.T)
func TestBindingChain_ToMap_Empty(t *testing.T)

// Test de merge
func TestBindingChain_Merge(t *testing.T)
func TestBindingChain_Merge_Conflicts(t *testing.T)

// Test edge cases
func TestBindingChain_Nil_Operations(t *testing.T)
func TestBindingChain_Long_Chain(t *testing.T)  // 100 bindings

// Test de performance
func BenchmarkBindingChain_Add(b *testing.B)
func BenchmarkBindingChain_Get(b *testing.B)
func BenchmarkBindingChain_Get_DeepChain(b *testing.B)
```

**Couverture cible** : > 95%

**Documenter dans** : Section "7.1. Tests BindingChain"

---

#### 6.2 Tests unitaires pour Token

**Fichier** : `rete/fact_token_test.go` (à adapter)

**Tests à ajouter/modifier** :
```go
func TestToken_CreateWithBindingChain(t *testing.T)
func TestToken_GetBinding(t *testing.T)
func TestToken_HasBinding(t *testing.T)
func TestToken_GetVariables(t *testing.T)
func TestToken_Metadata(t *testing.T)
```

**Documenter dans** : Section "7.2. Tests Token"

---

#### 6.3 Tests d'intégration pour JoinNode

**Fichier** : `rete/node_join_cascade_test.go` (nouveau)

**Tests à implémenter** :
```go
// Régression : 2 variables doivent continuer à fonctionner
func TestJoinNode_2Variables_UserOrder(t *testing.T)

// Cas principal : 3 variables
func TestJoinNode_3Variables_UserOrderProduct(t *testing.T) {
    // Setup : créer 3 TypeNodes, 2 JoinNodes en cascade
    // Test : soumettre User, Order, Product
    // Assert : token final contient [u, o, p]
}

// Variations d'ordre d'arrivée
func TestJoinNode_3Variables_DifferentOrders(t *testing.T) {
    // Tester 6 permutations d'ordre de soumission
}

// Scalabilité : 4 variables
func TestJoinNode_4Variables(t *testing.T)

// Test paramétrique : N variables
func TestJoinNode_NVariables(t *testing.T) {
    for n := 2; n <= 10; n++ {
        t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
            // Créer cascade de N variables
            // Vérifier que token final a N bindings
        })
    }
}
```

**Documenter dans** : Section "7.3. Tests Cascades"

---

#### 6.4 Tests E2E (validation finale)

**Fixtures à tester** :
1. `beta_join_complex.tsd` - Doit passer (actuellement échoue)
2. `join_multi_variable_complex.tsd` - Doit passer (actuellement échoue)
3. Tous les autres tests - Doivent continuer à passer (non-régression)

**Commandes de validation** :
```bash
make test-unit           # Tous les tests unitaires
make test-integration    # Tests d'intégration
make test-e2e           # Tests E2E (83 fixtures)
make test-complete      # Tout
```

**Critère de succès** : 83/83 tests E2E passent (100%)

**Documenter dans** : Section "7.4. Tests E2E"

---

### Tâche 7 : Stratégie de Migration Détaillée (30 min)

#### 7.1 Ordre de migration des fichiers

**Étapes séquentielles (Prompts 03-08)** :

**Prompt 03** : BindingChain (nouveau fichier)
- Créer `rete/binding_chain.go`
- Créer `rete/binding_chain_test.go`
- Valider : tests passent, code compile

**Prompt 04** : Token (remplacement)
- Modifier `rete/fact_token.go` : remplacer Token.Bindings
- Stratégie : renommer Token → TokenOld temporairement
- Créer nouvelle implémentation de Token
- Fixer TOUTES les erreurs de compilation
- Supprimer TokenOld
- Valider : code compile

**Prompt 05** : JoinNode - performJoinWithTokens
- Modifier `rete/node_join.go` : réécrire performJoinWithTokens
- Valider : tests de JoinNode passent

**Prompt 06** : JoinNode - Activation
- Modifier `rete/node_join.go` : réécrire ActivateLeft/ActivateRight
- Valider : tests d'intégration 2 variables passent

**Prompt 07** : BetaChainBuilder
- Modifier `rete/builder_beta_chain.go`
- Modifier `rete/builder_join_rules_cascade.go`
- Valider : construction correcte des cascades

**Prompt 08** : ExecutionContext et Actions
- Modifier `rete/action_executor_context.go`
- Modifier `rete/action_executor_evaluation.go`
- Modifier `rete/node_terminal.go`
- Valider : actions s'exécutent correctement

**Documenter dans** : Section "8. Stratégie de Migration"

---

#### 7.2 Gestion des erreurs de compilation

**Approche** : Migration "big bang" au Prompt 04

**Stratégie** :
1. Identifier TOUS les fichiers utilisant `Token.Bindings`
   ```bash
   grep -r "\.Bindings\[" tsd/rete/
   grep -r "\.Bindings =" tsd/rete/
   ```

2. Pour chaque occurrence, remplacer :
   - `token.Bindings[variable]` → `token.Bindings.Get(variable)`
   - `token.Bindings = map[...]` → `token.Bindings = NewBindingChain()`
   - `for k, v := range token.Bindings` → `for _, k := range token.Bindings.Variables() { v := token.Bindings.Get(k); ... }`

3. Compiler fréquemment :
   ```bash
   go build ./rete/...
   ```

4. Fixer les erreurs une par une

**Documenter dans** : Section "8.2. Migration du Code"

---

### Tâche 8 : Documenter le Design (30 min)

#### 8.1 Créer BINDINGS_DESIGN.md

**Chemin** : `tsd/docs/architecture/BINDINGS_DESIGN.md`

**Structure complète** :

```markdown
# Design du Système de Bindings Immuable

**Date** : [DATE]
**Auteur** : Design (Prompt 02)
**Version** : 1.0
**Status** : Spécification

---

## 1. Analyse des Besoins

### 1.1 Problème à Résoudre
[Résumé de BINDINGS_ANALYSIS.md]

### 1.2 Objectifs
[Ce que doit accomplir ce design]

### 1.3 Contraintes
[Limites et contraintes à respecter]

---

## 2. Exigences

### 2.1 Exigences Fonctionnelles
[REQ-1 à REQ-5]

### 2.2 Exigences Non-Fonctionnelles
[NFR-1 à NFR-5]

---

## 3. BindingChain - Spécification

### 3.1 Structure de Données
[Code Go avec commentaires détaillés]

### 3.2 API Complète
[Toutes les méthodes avec signatures]

### 3.3 Complexité et Performance
[Tableau des complexités]

### 3.4 Invariants
[Propriétés garanties]

### 3.5 Exemples d'Utilisation
[Code exemples concrets]

---

## 4. Token - Nouvelle Structure

### 4.1 Définition
[Code Go de la nouvelle structure]

### 4.2 API
[Méthodes de Token]

### 4.3 Stratégie de Migration
[Comment passer de l'ancien au nouveau]

### 4.4 Impact sur le Code Existant
[Fichiers affectés et changements nécessaires]

---

## 5. JoinNode - Nouvelle Logique

### 5.1 performJoinWithTokens
[Nouvelle implémentation avec explications]

### 5.2 Évaluation des Conditions
[Adaptation de evaluateJoinConditions]

### 5.3 Gestion des Mémoires
[Left/Right memory avec BindingChain]

---

## 6. ExecutionContext

### 6.1 Nouvelle Structure
[Code de ExecutionContext adapté]

### 6.2 Résolution de Variables
[Nouvelle implémentation]

---

## 7. Plan de Test

### 7.1 Tests BindingChain
[Liste des tests unitaires]

### 7.2 Tests Token
[Tests de Token]

### 7.3 Tests Cascades
[Tests d'intégration multi-niveaux]

### 7.4 Tests E2E
[Validation finale]

---

## 8. Stratégie de Migration

### 8.1 Ordre d'Exécution
[Prompts 03 à 08]

### 8.2 Migration du Code
[Stratégie "big bang" vs progressive]

### 8.3 Points de Validation
[Critères de succès à chaque étape]

---

## 9. Risques et Mitigation

### 9.1 Risques Identifiés
[Liste des risques]

### 9.2 Plans de Mitigation
[Comment gérer chaque risque]

---

## 10. Trade-offs et Décisions

### 10.1 Immutabilité vs Performance
**Décision** : Immutabilité prioritaire
**Justification** : [...]

### 10.2 Structural Sharing vs Deep Copy
**Décision** : Structural sharing
**Justification** : [...]

### 10.3 Cache vs Simplicité
**Décision** : Pas de cache initialement
**Justification** : [...]

---

## 11. Alternatives Considérées

### 11.1 Alternative 1 : Map avec Copy-on-Write
**Description** : [...]
**Avantages** : [...]
**Inconvénients** : [...]
**Rejet** : [raison]

### 11.2 Alternative 2 : Persistent Data Structure Library
**Description** : [...]
**Avantages** : [...]
**Inconvénients** : [...]
**Rejet** : [raison]

---

## 12. Validation du Design

### 12.1 Checklist de Validation
- [ ] Répond à toutes les exigences fonctionnelles
- [ ] Respecte les exigences non-fonctionnelles
- [ ] Résout le problème identifié dans BINDINGS_ANALYSIS.md
- [ ] API claire et cohérente
- [ ] Plan de test complet
- [ ] Stratégie de migration définie
- [ ] Risques identifiés et mitigés

### 12.2 Approbation
[À remplir après review]

---

## Annexes

### Annexe A : Diagrammes
[Diagrammes de classes, séquences, etc.]

### Annexe B : Exemples Complets
[Code complet d'utilisation]

### Annexe C : Benchmarks Préliminaires
[Estimations de performance]
```

**Documenter dans** : Créer le fichier complet

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Livrables
- [ ] ✅ Fichier `docs/architecture/BINDINGS_DESIGN.md` complet (800-1200 lignes)
- [ ] ✅ Spécification complète de BindingChain (structure + API + complexité)
- [ ] ✅ Spécification complète de Token (nouvelle structure)
- [ ] ✅ Design de la nouvelle logique JoinNode
- [ ] ✅ Plan de test détaillé (tous les tests à implémenter)
- [ ] ✅ Stratégie de migration séquentielle (Prompts 03-08)

### Qualité du Design
- [ ] Toutes les exigences sont adressées
- [ ] Les trade-offs sont documentés et justifiés
- [ ] Les alternatives sont considérées
- [ ] Les risques sont identifiés et mitigés
- [ ] Le code est lisible et maintenable
- [ ] La performance est acceptable (estimée)

### Validité Technique
- [ ] Les structures de données sont cohérentes
- [ ] Les invariants sont garantissables
- [ ] L'API est ergonomique
- [ ] La migration est faisable
- [ ] Les tests sont exhaustifs

---

## 🎯 Questions Clés - Réponses Attendues

À la fin de cette session, le design doit répondre :

1. **Structure** : Comment BindingChain est-elle organisée ?
   - Réponse : Structure récursive avec Variable, Fact, Parent

2. **Immutabilité** : Comment garantir que les bindings ne sont jamais perdus ?
   - Réponse : Add() retourne nouvelle chaîne, structural sharing

3. **Performance** : Quelle est la complexité de Get() ?
   - Réponse : O(n) avec n < 10 typiquement, acceptable

4. **Migration** : Comment migrer sans cohabitation ancien/nouveau ?
   - Réponse : Remplacement direct de Token.Bindings, fix all compilation errors

5. **Tests** : Comment valider que ça fonctionne ?
   - Réponse : Tests unitaires (BindingChain, Token) + intégration (cascades) + E2E

---

## 🎯 Prochaine Étape

Une fois ce design **complet et validé**, passer au **Prompt 03 - Implémentation de BindingChain**.

Le Prompt 03 implémentera la spécification définie ici.

---

## 💡 Conseils Pratiques

### Pour un Bon Design
1. **Être précis** : Spécifier les signatures exactes, pas juste des idées
2. **Être complet** : Couvrir tous les cas d'usage, y compris edge cases
3. **Être réaliste** : Estimer la complexité et la faisabilité
4. **Être pragmatique** : Privilégier la simplicité, optimiser si nécessaire

### Pour la Documentation
1. Utiliser des exemples de code concrets
2. Expliquer le "pourquoi" des décisions
3. Documenter les alternatives rejetées
4. Inclure des diagrammes quand pertinent

### Pour la Validation
1. Vérifier que chaque exigence est adressée
2. S'assurer que la migration est faisable
3. Valider que les tests sont suffisants
4. Confirmer que les risques sont mitigés

---

**Note** : Ce design est une **spécification**, pas une implémentation. Aucun code de production n'est écrit dans cette session. Le but est de **PLANIFIER**, pas encore de **CODER**.