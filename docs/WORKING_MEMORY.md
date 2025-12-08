# WorkingMemory - Mémoire de Travail du Réseau RETE

**Module** : `rete`  
**Fichier source** : `rete/fact_token.go`  
**Date** : 8 décembre 2024

---

## Table des Matières

1. [Vue d'Ensemble](#vue-densemble)
2. [Structure de Données](#structure-de-données)
3. [Rôle et Responsabilités](#rôle-et-responsabilités)
4. [Utilisation par Type de Nœud](#utilisation-par-type-de-nœud)
5. [Cycle de Vie](#cycle-de-vie)
6. [Opérations](#opérations)
7. [Gestion de la Mémoire](#gestion-de-la-mémoire)
8. [Patterns d'Usage](#patterns-dusage)
9. [Considérations de Performance](#considérations-de-performance)
10. [Exemples Concrets](#exemples-concrets)

---

## Vue d'Ensemble

La **WorkingMemory** est une structure centrale dans l'implémentation du réseau RETE de TSD. Elle représente la **mémoire de travail** locale de chaque nœud du réseau, stockant les faits et tokens qui ont été traités et validés par ce nœud.

### Définition

```go
type WorkingMemory struct {
    NodeID string            `json:"node_id"`  // Identifiant du nœud propriétaire
    Facts  map[string]*Fact  `json:"facts"`    // Faits stockés (clé = internal_id)
    Tokens map[string]*Token `json:"tokens"`   // Tokens stockés (clé = token_id)
}
```

### Concepts Clés

- **Localité** : Chaque nœud possède sa propre WorkingMemory indépendante
- **État** : Représente l'état actuel des correspondances à ce point du réseau
- **Incrémentalité** : Permet l'évaluation incrémentale de l'algorithme RETE
- **Partage** : Facilite le partage de nœuds entre règles (alpha/beta sharing)

---

## Structure de Données

### Composants

#### 1. NodeID (string)
Identifiant unique du nœud propriétaire de cette mémoire.

```go
NodeID: "alpha_Person_age>18"
NodeID: "join_PersonOrder"
NodeID: "terminal_rule_discount"
```

#### 2. Facts (map[string]*Fact)
Dictionnaire des **faits** validés par ce nœud.

**Clé** : Internal ID au format `Type_ID`
- Ex: `"Person_123"`, `"Order_456"`

**Valeur** : Pointeur vers la structure `Fact`

```go
type Fact struct {
    ID        string                 // "123"
    Type      string                 // "Person"
    Fields    map[string]interface{} // {"name": "Alice", "age": 25}
    Timestamp time.Time
}
```

**Utilisation** :
- Nœuds **TypeNode** : Stockent tous les faits de leur type
- Nœuds **AlphaNode** : Stockent les faits qui satisfont leur condition
- Nœuds **RootNode** : Stockent tous les faits du système

#### 3. Tokens (map[string]*Token)
Dictionnaire des **tokens** (correspondances partielles) créés par ce nœud.

**Clé** : Token ID unique
- Ex: `"token_123"`, `"join_token_456_789"`

**Valeur** : Pointeur vers la structure `Token`

```go
type Token struct {
    ID           string           // Identifiant unique
    Facts        []*Fact          // Liste de faits combinés
    NodeID       string           // Nœud qui a créé ce token
    Parent       *Token           // Token parent (chaînage)
    Bindings     map[string]*Fact // variable_name → fact
    IsJoinResult bool             // Indique un token de jointure
}
```

**Utilisation** :
- Nœuds **JoinNode** : Stockent les tokens résultant des jointures
- Nœuds **AccumulatorNode** : Stockent les tokens avec résultats d'agrégation
- Nœuds **ExistsNode/NotNode** : Stockent les tokens validés

---

## Rôle et Responsabilités

### 1. Conservation de l'État

La WorkingMemory maintient **l'état local** du nœud :
- Quels faits ont été traités et validés ?
- Quelles combinaisons de faits (tokens) ont réussi ?
- Quel est l'historique des correspondances ?

### 2. Support de l'Incrémentalité

RETE est un algorithme **incrémental** : seuls les changements sont recalculés.

La WorkingMemory permet :
- **Ajout de fait** : Comparer avec les faits déjà en mémoire
- **Retrait de fait** : Identifier et supprimer les tokens impactés
- **Jointures** : Croiser les nouvelles données avec l'existant

### 3. Facilitation du Partage

Avec alpha/beta sharing, plusieurs règles peuvent partager un même nœud.

La WorkingMemory centralisée :
- **Évite les recalculs** : Un fait testé une fois pour toutes les règles
- **Réduit la mémoire** : Une seule copie des résultats
- **Améliore la cohérence** : État unique et synchronisé

### 4. Support des Jointures

Pour les **JoinNodes**, la WorkingMemory est cruciale :
- **LeftMemory** : Tokens venant de la branche gauche
- **RightMemory** : Tokens venant de la branche droite
- **ResultMemory** : Tokens résultant des jointures réussies

---

## Utilisation par Type de Nœud

### RootNode
```go
type RootNode struct {
    BaseNode  // Contient Memory *WorkingMemory
}
```

**Usage** :
- `Memory.Facts` : **Tous les faits** du système
- `Memory.Tokens` : Non utilisé (pas de tokens au root)

**Opérations** :
```go
// Ajout d'un fait
rn.Memory.AddFact(fact)

// Propagation aux TypeNodes enfants
for _, child := range rn.Children {
    child.ActivateRight(fact)
}
```

---

### TypeNode
```go
type TypeNode struct {
    BaseNode
    TypeName       string
    TypeDefinition TypeDefinition
}
```

**Usage** :
- `Memory.Facts` : Faits de ce type spécifique
- `Memory.Tokens` : Non utilisé

**Opérations** :
```go
// Validation et stockage
if fact.Type == tn.TypeName {
    tn.Memory.AddFact(fact)
    // Propager aux AlphaNodes
}
```

---

### AlphaNode
```go
type AlphaNode struct {
    BaseNode
    Condition    interface{}
    VariableName string
}
```

**Usage** :
- `Memory.Facts` : Faits satisfaisant la condition alpha
- `Memory.Tokens` : Tokens créés pour ces faits (dans certains cas)

**Opérations** :
```go
// Test et stockage
if evaluator.EvaluateCondition(an.Condition, fact) {
    an.Memory.AddFact(fact)
    
    // Créer token pour propagation
    token := &Token{
        ID:       fmt.Sprintf("alpha_token_%s_%s", an.ID, fact.ID),
        Facts:    []*Fact{fact},
        Bindings: map[string]*Fact{an.VariableName: fact},
    }
    an.PropagateToChildren(nil, token)
}
```

---

### JoinNode (Architecture Triple-Mémoire)
```go
type JoinNode struct {
    BaseNode
    LeftMemory   *WorkingMemory  // Tokens de gauche
    RightMemory  *WorkingMemory  // Tokens de droite
    ResultMemory *WorkingMemory  // Résultats de jointure
}
```

**Architecture RETE classique** :
```
    LeftMemory          RightMemory
        ↓                   ↓
        └─── JoinTest ─────┘
                ↓
          ResultMemory
                ↓
         PropagateToChildren
```

**Usage détaillé** :

#### LeftMemory
- Stocke les **tokens** venant de la branche gauche (généralement des AlphaNodes ou autres JoinNodes)
- Utilisée pour jointures avec nouveaux faits de droite

```go
func (jn *JoinNode) ActivateLeft(token *Token) error {
    jn.LeftMemory.AddToken(token)
    
    // Essayer de joindre avec tous les tokens en RightMemory
    for _, rightToken := range jn.RightMemory.GetTokens() {
        if joinedToken := jn.performJoin(token, rightToken); joinedToken != nil {
            jn.ResultMemory.AddToken(joinedToken)
            jn.PropagateToChildren(nil, joinedToken)
        }
    }
}
```

#### RightMemory
- Stocke les **faits/tokens** venant de la branche droite (nouveaux faits)
- Utilisée pour jointures avec tokens existants en LeftMemory

```go
func (jn *JoinNode) ActivateRight(fact *Fact) error {
    factToken := &Token{Facts: []*Fact{fact}}
    jn.RightMemory.AddToken(factToken)
    
    // Essayer de joindre avec tous les tokens en LeftMemory
    for _, leftToken := range jn.LeftMemory.GetTokens() {
        if joinedToken := jn.performJoin(leftToken, factToken); joinedToken != nil {
            jn.ResultMemory.AddToken(joinedToken)
            jn.PropagateToChildren(nil, joinedToken)
        }
    }
}
```

#### ResultMemory
- Stocke uniquement les **tokens de jointure réussie**
- Permet de tracer les résultats validés
- Utilisée pour retrait/invalidation

```go
// Token de jointure réussi
joinedToken := &Token{
    ID:           fmt.Sprintf("join_%s_%s", leftToken.ID, rightToken.ID),
    Facts:        append(leftToken.Facts, rightToken.Facts...),
    Bindings:     mergedBindings,
    IsJoinResult: true,
}
jn.ResultMemory.AddToken(joinedToken)
```

**Avantage de cette architecture** :
- **Symétrie** : Jointures possibles dans les deux sens (gauche→droite et droite→gauche)
- **Incrémentalité** : Nouveaux faits comparés uniquement avec l'existant
- **Efficacité** : Pas de recalcul complet à chaque ajout

---

### AccumulatorNode
```go
type AccumulatorNode struct {
    BaseNode
    MainFacts    map[string]*Fact  // Faits principaux (indexés séparément)
    AllFacts     map[string]*Fact  // Tous les faits
}
```

**Usage** :
- `Memory.Tokens` : Tokens avec résultats d'agrégation
- `MainFacts` : Faits principaux à agréger
- `AllFacts` : Tous les faits nécessaires aux calculs

**Opérations** :
```go
// Calcul d'agrégation
aggregatedValue := an.calculateAggregate(mainFact, aggregatedFacts)

// Stocker si condition satisfaite
if an.evaluateCondition(aggregatedValue) {
    token := &Token{
        ID:       fmt.Sprintf("accum_%s", mainFact.ID),
        Facts:    []*Fact{mainFact},
        Bindings: map[string]*Fact{an.MainVariable: mainFact},
    }
    an.Memory.AddToken(token)
}
```

---

### ExistsNode
```go
type ExistsNode struct {
    BaseNode
    MainMemory   *WorkingMemory  // Tokens principaux
    ExistsMemory *WorkingMemory  // Faits testés pour existence
    ResultMemory *WorkingMemory  // Tokens validés
}
```

**Usage** :
- `MainMemory` : Tokens en attente de validation
- `ExistsMemory` : Faits qui satisfont le pattern EXISTS
- `ResultMemory` : Tokens validés (au moins un fait existe)

---

### TerminalNode
```go
type TerminalNode struct {
    BaseNode
    RuleID string
    Action Action
}
```

**Usage** :
- `Memory.Tokens` : Tokens qui ont atteint le terminal (règle satisfaite)
- Utilisé pour éviter de réexécuter la même action

---

## Cycle de Vie

### 1. Création
La WorkingMemory est créée lors de l'initialisation du nœud :

```go
node := &AlphaNode{
    BaseNode: BaseNode{
        ID:     "alpha_Person_age>18",
        Memory: &WorkingMemory{
            NodeID: "alpha_Person_age>18",
            Facts:  make(map[string]*Fact),
            Tokens: make(map[string]*Token),
        },
    },
}
```

### 2. Remplissage
Au fur et à mesure des activations :

```go
// Ajout incrémental
node.Memory.AddFact(fact1)
node.Memory.AddFact(fact2)
node.Memory.AddToken(token1)
```

### 3. Requêtage
Pour jointures et évaluations :

```go
// Récupération pour comparaison
existingFacts := node.Memory.GetFacts()
for _, fact := range existingFacts {
    // Comparer avec nouveau fait
}
```

### 4. Nettoyage
Lors de la rétractation :

```go
// Retrait
node.Memory.RemoveFact(factID)
node.Memory.RemoveToken(tokenID)
```

### 5. Persistance (optionnel)
Via le Storage :

```go
// Sauvegarde
storage.SaveMemory(node.ID, node.Memory)

// Restauration
memory, _ := storage.LoadMemory(node.ID)
node.Memory = memory
```

---

## Opérations

### API Complète

#### Gestion des Faits

```go
// AddFact ajoute un fait (clé = internal_id)
func (wm *WorkingMemory) AddFact(fact *Fact) error

// RemoveFact supprime un fait
func (wm *WorkingMemory) RemoveFact(factID string)

// GetFact récupère un fait par internal_id
func (wm *WorkingMemory) GetFact(internalID string) (*Fact, bool)

// GetFactByTypeAndID récupère un fait par type et ID
func (wm *WorkingMemory) GetFactByTypeAndID(factType, factID string) (*Fact, bool)

// GetFacts retourne tous les faits
func (wm *WorkingMemory) GetFacts() []*Fact
```

#### Gestion des Tokens

```go
// AddToken ajoute un token
func (wm *WorkingMemory) AddToken(token *Token)

// RemoveToken supprime un token
func (wm *WorkingMemory) RemoveToken(tokenID string)

// GetTokens retourne tous les tokens
func (wm *WorkingMemory) GetTokens() []*Token
```

#### Utilitaires

```go
// Clone crée une copie profonde
func (wm *WorkingMemory) Clone() *WorkingMemory

// GetFactsByVariable filtre par variable (implémentation future)
func (wm *WorkingMemory) GetFactsByVariable(variables []string) []*Fact

// GetTokensByVariable filtre par variable (implémentation future)
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token
```

---

## Gestion de la Mémoire

### Stratégies de Gestion

#### 1. Croissance Dynamique
Les maps Go grandissent dynamiquement :
```go
Facts:  make(map[string]*Fact)   // Capacité initiale petite
Tokens: make(map[string]*Token)  // Grandit automatiquement
```

#### 2. Éviction (Non Implémenté)
Actuellement : pas d'éviction automatique.

Future : LRU ou Time-based eviction :
```go
type WorkingMemory struct {
    Facts      map[string]*Fact
    FactsLRU   *lru.Cache        // Cache LRU
    MaxSize    int               // Limite de taille
    LastAccess map[string]time.Time
}
```

#### 3. Nettoyage lors de Rétractation
Suppression en cascade lors du retrait de fait :

```go
func (jn *JoinNode) ActivateRetract(factID string) error {
    // Retirer des 3 mémoires
    for tokenID, token := range jn.LeftMemory.Tokens {
        if containsFact(token, factID) {
            jn.LeftMemory.RemoveToken(tokenID)
        }
    }
    // Idem pour RightMemory et ResultMemory
}
```

---

## Patterns d'Usage

### Pattern 1 : Test-and-Store (AlphaNode)

```go
// Tester condition, puis stocker si satisfait
if evaluator.Evaluate(condition, fact) {
    node.Memory.AddFact(fact)
    propagate(fact)
}
```

### Pattern 2 : Join-and-Store (JoinNode)

```go
// Stocker à gauche, puis joindre avec droite
leftMemory.AddToken(leftToken)

for _, rightToken := range rightMemory.GetTokens() {
    if joinTest(leftToken, rightToken) {
        joined := merge(leftToken, rightToken)
        resultMemory.AddToken(joined)
        propagate(joined)
    }
}
```

### Pattern 3 : Accumulate-and-Store (AccumulatorNode)

```go
// Stocker, calculer agrégation, tester condition
allFacts[fact.ID] = fact

if fact.Type == mainType {
    mainFacts[fact.ID] = fact
    aggregatedValue := calculate(fact, relatedFacts)
    
    if condition.Evaluate(aggregatedValue) {
        token := createToken(fact, aggregatedValue)
        memory.AddToken(token)
        propagate(token)
    }
}
```

### Pattern 4 : Exists-and-Validate (ExistsNode)

```go
// Stocker main token, tester existence, valider
mainMemory.AddToken(token)

existingFacts := existsMemory.GetFacts()
if len(existingFacts) > 0 {
    // Au moins un fait existe
    resultMemory.AddToken(token)
    propagate(token)
}
```

---

## Considérations de Performance

### Complexité

#### Temps
- **AddFact/AddToken** : O(1) - insertion dans map
- **RemoveFact/RemoveToken** : O(1) - suppression dans map
- **GetFact/GetToken** : O(1) - lookup dans map
- **GetFacts/GetTokens** : O(n) - itération complète

#### Espace
- **Par nœud** : O(F + T) où F=faits, T=tokens
- **Réseau complet** : O(N × (F + T)) où N=nombre de nœuds

Avec **alpha/beta sharing** :
- **Réduction** : 40-70% moins de mémoire
- **Raison** : Nœuds partagés = WorkingMemory partagée

### Optimisations

#### 1. Indexation (Future)
Pour accélérer les recherches :
```go
type IndexedWorkingMemory struct {
    WorkingMemory
    FactsByType  map[string][]*Fact         // Type → Facts
    FactsByField map[string]map[interface{}][]*Fact // Field → Value → Facts
}
```

#### 2. Compaction (Future)
Nettoyer les tokens obsolètes :
```go
func (wm *WorkingMemory) Compact() {
    // Supprimer tokens sans références
    // Supprimer faits non utilisés
}
```

#### 3. Lazy Loading (Future)
Charger depuis storage à la demande :
```go
func (wm *WorkingMemory) GetFact(id string) (*Fact, bool) {
    if fact, exists := wm.Facts[id]; exists {
        return fact, true
    }
    // Charger depuis storage si pas en cache
    return wm.loadFromStorage(id)
}
```

---

## Exemples Concrets

### Exemple 1 : AlphaNode Simple

```go
// Règle : { p: Person } / p.age > 18 ==> ...

// Création du nœud
alphaNode := &AlphaNode{
    BaseNode: BaseNode{
        ID: "alpha_Person_age>18",
        Memory: &WorkingMemory{
            NodeID: "alpha_Person_age>18",
            Facts:  make(map[string]*Fact),
            Tokens: make(map[string]*Token),
        },
    },
    Condition: map[string]interface{}{
        "type": "comparison",
        "left": "p.age",
        "operator": ">",
        "right": 18,
    },
}

// Arrivée de faits
fact1 := &Fact{ID: "1", Type: "Person", Fields: map[string]interface{}{"age": 25}}
fact2 := &Fact{ID: "2", Type: "Person", Fields: map[string]interface{}{"age": 16}}

// Traitement
alphaNode.ActivateRight(fact1)  // age=25 > 18 ✅
// → alphaNode.Memory.Facts["Person_1"] = fact1
// → Propagation aux enfants

alphaNode.ActivateRight(fact2)  // age=16 > 18 ❌
// → Rejeté, pas stocké

// État final de la mémoire
// alphaNode.Memory.Facts = {"Person_1": fact1}
// alphaNode.Memory.Tokens = {} (vide pour ce type de nœud)
```

### Exemple 2 : JoinNode avec Triple-Mémoire

```go
// Règle : { p: Person, o: Order } / p.id == o.customer_id ==> ...

// Création du nœud
joinNode := &JoinNode{
    BaseNode: BaseNode{
        ID: "join_PersonOrder",
        Memory: &WorkingMemory{...},  // Pour comptage global
    },
    LeftMemory:  &WorkingMemory{...},  // Tokens Person
    RightMemory: &WorkingMemory{...},  // Facts Order
    ResultMemory: &WorkingMemory{...}, // Jointures réussies
    JoinConditions: []JoinCondition{
        {LeftField: "p.id", RightField: "o.customer_id", Operator: "=="},
    },
}

// Scénario
// 1. Arrivée d'un token Person (depuis alpha node)
personToken := &Token{
    ID: "token_p1",
    Facts: []*Fact{{ID: "1", Type: "Person", Fields: map[string]interface{}{"id": "123"}}},
    Bindings: map[string]*Fact{"p": ...},
}
joinNode.ActivateLeft(personToken)
// → joinNode.LeftMemory.Tokens["token_p1"] = personToken
// → Pas encore de RightMemory, donc pas de jointure

// 2. Arrivée d'un fait Order
orderFact := &Fact{
    ID: "100", 
    Type: "Order", 
    Fields: map[string]interface{}{"customer_id": "123"},
}
joinNode.ActivateRight(orderFact)
// → Création d'un token pour Order
// → joinNode.RightMemory.Tokens["token_o100"] = orderToken
// → Test de jointure avec LeftMemory
// → p.id (123) == o.customer_id (123) ✅

// → Création token de jointure
joinedToken := &Token{
    ID: "join_token_p1_o100",
    Facts: []*Fact{personFact, orderFact},
    Bindings: map[string]*Fact{
        "p": personFact,
        "o": orderFact,
    },
    IsJoinResult: true,
}
// → joinNode.ResultMemory.Tokens["join_token_p1_o100"] = joinedToken
// → Propagation aux enfants

// État final des mémoires
// LeftMemory.Tokens = {"token_p1": personToken}
// RightMemory.Tokens = {"token_o100": orderToken}
// ResultMemory.Tokens = {"join_token_p1_o100": joinedToken}
```

### Exemple 3 : AccumulatorNode

```go
// Règle : { e: Employee } / SUM(p.score FROM Performance p WHERE p.employee_id == e.id) > 100 ==> ...

accNode := &AccumulatorNode{
    BaseNode: BaseNode{
        ID: "accum_EmployeePerf",
        Memory: &WorkingMemory{...},
    },
    MainFacts: make(map[string]*Fact),  // Employees
    AllFacts:  make(map[string]*Fact),  // Employees + Performances
    AggregateFunc: "SUM",
    Field: "score",
}

// Arrivée d'un Employee
empFact := &Fact{ID: "e1", Type: "Employee", Fields: map[string]interface{}{"id": "123"}}
accNode.Activate(empFact, nil)
// → accNode.MainFacts["e1"] = empFact
// → accNode.AllFacts["e1"] = empFact
// → Pas encore de Performance, sum=0, condition non satisfaite

// Arrivée d'une Performance
perf1 := &Fact{ID: "p1", Type: "Performance", Fields: map[string]interface{}{
    "employee_id": "123",
    "score": 60,
}}
accNode.Activate(perf1, nil)
// → accNode.AllFacts["p1"] = perf1
// → Recalcul pour Employee e1
// → Performances liées : [perf1]
// → SUM(score) = 60, condition 60 > 100 ❌

// Arrivée d'une autre Performance
perf2 := &Fact{ID: "p2", Type: "Performance", Fields: map[string]interface{}{
    "employee_id": "123",
    "score": 50,
}}
accNode.Activate(perf2, nil)
// → accNode.AllFacts["p2"] = perf2
// → Recalcul pour Employee e1
// → Performances liées : [perf1, perf2]
// → SUM(score) = 110, condition 110 > 100 ✅

// → Création token
token := &Token{
    ID: "accum_e1",
    Facts: []*Fact{empFact},
    Bindings: map[string]*Fact{"e": empFact},
}
// → accNode.Memory.Tokens["accum_e1"] = token
// → Propagation aux enfants

// État final
// MainFacts = {"e1": empFact}
// AllFacts = {"e1": empFact, "p1": perf1, "p2": perf2}
// Memory.Tokens = {"accum_e1": token}
```

---

## Résumé

### Points Clés

1. **WorkingMemory = État Local**
   - Chaque nœud maintient son propre état
   - Faits et tokens validés par ce nœud

2. **Deux Collections Distinctes**
   - `Facts` : Données brutes (faits)
   - `Tokens` : Combinaisons/correspondances (tokens)

3. **Architecture Variable par Nœud**
   - Alpha/Type : Principalement Facts
   - Join : Triple-mémoire (Left/Right/Result)
   - Accumulator : Facts + index spécialisés

4. **Support de l'Incrémentalité**
   - Nouveaux faits comparés avec existant
   - Pas de recalcul complet
   - Performances O(1) pour la plupart des ops

5. **Facilite le Partage**
   - Nœuds partagés = mémoire partagée
   - Calculs une seule fois
   - Réduction 40-70% de la mémoire

### Usage Recommandé

- **Lire la mémoire** avant d'agir (jointures, comparaisons)
- **Écrire après validation** (ne stocker que ce qui passe)
- **Nettoyer lors de rétractations** (cohérence)
- **Cloner si modification** (immutabilité des tokens)

---

**Voir aussi** :
- [fact_token.go](../rete/fact_token.go) - Implémentation
- [node_base.go](../rete/node_base.go) - BaseNode avec Memory
- [node_join.go](../rete/node_join.go) - Triple-mémoire en action
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architecture globale du système