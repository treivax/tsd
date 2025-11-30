# Beta Node Sharing : Concepts et Mécanismes

## Table des Matières

1. [Concepts de base](#concepts-de-base)
2. [Différence avec Alpha Sharing](#différence-avec-alpha-sharing)
3. [Quand les JoinNodes sont partagés](#quand-les-joinnodes-sont-partagés)
4. [Diagrammes explicatifs](#diagrammes-explicatifs)
5. [Exemples visuels](#exemples-visuels)
6. [Mécanismes internes](#mécanismes-internes)

---

## Concepts de base

### Qu'est-ce qu'un Beta Node ?

Un **Beta Node** (ou JoinNode) est un nœud du réseau RETE qui effectue une **jointure** entre deux sources de données :

- **Mémoire gauche (left memory)** : Contient des tokens (séquences de faits)
- **Mémoire droite (right memory)** : Contient des faits individuels
- **Tests de jointure** : Conditions pour joindre les tokens et les faits

```
┌─────────────────────────────────────────────────────────┐
│                      JoinNode                           │
│                                                         │
│  ┌──────────────────┐         ┌──────────────────┐    │
│  │  Left Memory     │         │  Right Memory    │    │
│  │  (Tokens)        │         │  (Facts)         │    │
│  ├──────────────────┤         ├──────────────────┤    │
│  │ [Customer, ]     │         │ Order{id:1}      │    │
│  │ [Customer, ]     │   ⚡    │ Order{id:2}      │    │
│  │ [Customer, ]     │  Join   │ Order{id:3}      │    │
│  └──────────────────┘  Tests  └──────────────────┘    │
│                                                         │
│  Tests: Customer.id == Order.customerId                │
│                                                         │
│  Output: [Customer, Order] tokens                      │
└─────────────────────────────────────────────────────────┘
```

### Qu'est-ce que le Beta Node Sharing ?

Le **Beta Node Sharing** est une technique d'optimisation qui permet de **réutiliser le même JoinNode** pour plusieurs règles ayant des conditions de jointure identiques.

**Principe :**

```
Sans partage :
  Règle A → JoinNode A → Terminal A
  Règle B → JoinNode B → Terminal B
  (JoinNode A et B effectuent la même jointure)

Avec partage :
  Règle A ──┐
           JoinNode (partagé) → Terminal A
  Règle B ──┘                 → Terminal B
  (Une seule évaluation de la jointure)
```

### Pourquoi partager les Beta Nodes ?

**Avantages :**

1. **Moins de nœuds créés** → Réduction de 40-70% de la mémoire
2. **Calculs partagés** → Réduction de 30-50% du temps d'exécution
3. **Mémoires partagées** → Meilleure utilisation du cache CPU
4. **Propagation unique** → Moins d'overhead de communication

**Exemple concret :**

```
10 règles similaires sans partage :
- 10 JoinNodes × 10 KB = 100 KB
- 10 évaluations de jointure par activation

10 règles similaires avec partage :
- 1 JoinNode × 10 KB = 10 KB (90% d'économie)
- 1 évaluation de jointure par activation (90% de gain)
```

---

## Différence avec Alpha Sharing

### Alpha Sharing vs Beta Sharing

| Aspect | Alpha Sharing | Beta Sharing |
|--------|---------------|--------------|
| **Nœuds concernés** | AlphaNodes | JoinNodes |
| **Type de test** | Test sur 1 fait | Jointure entre 2 faits |
| **Mémoire** | Pas de mémoire | Left + Right memory |
| **Complexité** | O(1) par test | O(n×m) avec n,m tailles mémoires |
| **Impact mémoire** | Faible (~1KB/node) | Élevé (~10KB/node) |
| **Gain typique** | 20-30% | 40-70% |

### Alpha Node (tests sur un fait)

```
AlphaNode : Teste les propriétés d'un seul fait

┌────────────────────────────┐
│      AlphaNode             │
│                            │
│  Type: Order               │
│  Tests:                    │
│    - status == "active"    │
│    - amount > 100          │
│                            │
│  Input: Order fact         │
│  Output: Order (si match)  │
└────────────────────────────┘

Exemple de règle :
  when
    Order(status == "active", amount > 100)
  then
    ...
```

**Alpha Sharing :** Réutilise l'AlphaNode pour toutes les règles testant `Order(status == "active", amount > 100)`.

### Beta Node (jointure entre faits)

```
JoinNode : Joint deux faits selon des conditions

┌────────────────────────────────────────────┐
│           JoinNode                         │
│                                            │
│  Left: Customer fact                       │
│  Right: Order fact                         │
│  Tests:                                    │
│    - Customer.id == Order.customerId       │
│                                            │
│  Input: (Customer token, Order fact)       │
│  Output: [Customer, Order] token          │
└────────────────────────────────────────────┘

Exemple de règle :
  when
    Customer($custId : id)
    Order(customerId == $custId)
  then
    ...
```

**Beta Sharing :** Réutilise le JoinNode pour toutes les règles joignant Customer et Order sur `id == customerId`.

### Visualisation de la différence

```
Réseau RETE avec Alpha et Beta Sharing :

                   [RootNode]
                       │
              ┌────────┴────────┐
              │                 │
        [TypeNode]        [TypeNode]
        Customer           Order
              │                 │
              │                 │
        [AlphaNode]       [AlphaNode]        ← ALPHA SHARING
        tier=="premium"   status=="active"
              │                 │
              └────────┬────────┘
                       │
                  [JoinNode]                  ← BETA SHARING
              Customer.id==Order.customerId
                       │
              ┌────────┴────────┬────────┐
              │                 │        │
         [Terminal]        [Terminal] [Terminal]
          Rule1             Rule2     Rule3
```

### Différence dans le code

**Alpha Sharing (simple) :**

```go
type AlphaNode struct {
    ID         string
    Type       string
    Conditions []Condition
    Children   []Node
}

func (n *AlphaNode) Activate(fact *Fact) {
    // Évaluer les conditions
    if n.evaluateConditions(fact) {
        // Propager aux enfants
        for _, child := range n.Children {
            child.Activate(fact)
        }
    }
}
```

**Beta Sharing (complexe) :**

```go
type JoinNode struct {
    ID          string
    Tests       []JoinTest
    LeftParent  Node
    RightParent Node
    Children    []Node
    LeftMemory  *Memory    // Stock des tokens
    RightMemory *Memory    // Stock des faits
    Lifecycle   *NodeLifecycle  // Compteur de références
}

func (n *JoinNode) ActivateLeft(token *Token) {
    // 1. Stocker le token en mémoire gauche
    n.LeftMemory.Add(token)
    
    // 2. Joindre avec tous les faits de la mémoire droite
    for _, fact := range n.RightMemory.GetAll() {
        if n.evaluateJoinTests(token, fact) {
            newToken := append(token.Facts, fact)
            // 3. Propager aux enfants
            for _, child := range n.Children {
                child.Activate(newToken)
            }
        }
    }
}
```

---

## Quand les JoinNodes sont partagés

### Critères de partage

Un JoinNode est partagé entre deux règles si et seulement si :

1. **Même séquence de patterns** jusqu'à ce point
2. **Mêmes tests de jointure** (après normalisation)
3. **Même structure de token** en entrée

### Exemple 1 : Partage complet

```tsd
rule "Rule1"
when
    Customer($custId : id, type == "premium")
    Order(customerId == $custId, status == "active")
then
    action1();
end

rule "Rule2"
when
    Customer($custId : id, type == "premium")
    Order(customerId == $custId, status == "active")
then
    action2();
end
```

**Résultat :** Les deux règles partagent **tous** les nœuds (AlphaNodes ET JoinNode).

```
[Customer AlphaNode] ─┐
                      │
                      ├─► [JoinNode] ─┬─► [Terminal Rule1]
                      │              │
[Order AlphaNode] ────┘              └─► [Terminal Rule2]

Sharing ratio: 100%
```

### Exemple 2 : Partage partiel

```tsd
rule "Rule1"
when
    Customer($custId : id, type == "premium")
    Order(customerId == $custId, status == "active", amount > 100)
then
    action1();
end

rule "Rule2"
when
    Customer($custId : id, type == "premium")
    Order(customerId == $custId, status == "active", amount > 500)
then
    action2();
end
```

**Résultat :** Les règles partagent le JoinNode Customer-Order, mais ont des AlphaNodes Order différents.

```
[Customer AlphaNode]
         │
         ├─────────────┬─────────────┐
         │             │             │
[Order AlphaNode]  [Order AlphaNode]  (différents: amount)
 amount > 100      amount > 500
         │             │
    [JoinNode]    [JoinNode]      (PARTAGÉ si normalisé identiquement)
         │             │
   [Terminal]     [Terminal]
     Rule1          Rule2
```

### Exemple 3 : Pas de partage

```tsd
rule "Rule1"
when
    Customer($custId : id, type == "premium")
    Order(customerId == $custId)
then
    action1();
end

rule "Rule2"
when
    Supplier($suppId : id, region == "EU")
    Product(supplierId == $suppId)
then
    action2();
end
```

**Résultat :** Aucun partage (patterns complètement différents).

```
[Customer AlphaNode] ──► [JoinNode 1] ──► [Terminal Rule1]

[Supplier AlphaNode] ──► [JoinNode 2] ──► [Terminal Rule2]

Sharing ratio: 0%
```

### Exemple 4 : Partage de préfixe

```tsd
rule "Rule1"
when
    Customer($custId : id)
    Order(customerId == $custId)
    LineItem(orderId == $orderId, quantity > 1)
then
    action1();
end

rule "Rule2"
when
    Customer($custId : id)
    Order(customerId == $custId)
    Payment(orderId == $orderId, method == "credit")
then
    action2();
end
```

**Résultat :** Les règles partagent le préfixe Customer-Order.

```
[Customer AlphaNode]
         │
         └──► [JoinNode Customer-Order] (PARTAGÉ)
                       │
              ┌────────┴────────┐
              │                 │
        [JoinNode]         [JoinNode]
       Order-LineItem    Order-Payment
              │                 │
         [Terminal]        [Terminal]
           Rule1             Rule2
```

### Normalisation et partage

La **normalisation** des patterns garantit que des conditions équivalentes produisent la même clé de hash :

**Exemple : Ordre des contraintes**

```tsd
// Pattern A
Order(status == "active", amount > 100)

// Pattern B (ordre différent)
Order(amount > 100, status == "active")

// Après normalisation : IDENTIQUES
Order(amount > 100, status == "active")
```

**Exemple : Opérateurs commutatifs**

```tsd
// Pattern A
Customer.id == Order.customerId

// Pattern B (inversé)
Order.customerId == Customer.id

// Après normalisation : IDENTIQUES
Customer.id == Order.customerId
```

---

## Diagrammes explicatifs

### Diagramme 1 : Anatomie d'un JoinNode

```
┌───────────────────────────────────────────────────────────────────┐
│                         JoinNode                                  │
├───────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Métadonnées:                                                     │
│    ID: "join_abc123"                                              │
│    Type: NodeTypeJoin                                             │
│    Lifecycle: {RefCount: 3, RuleIDs: [r1, r2, r3]}               │
│                                                                   │
├─────────────────────────┬─────────────────────────────────────────┤
│   Left Memory           │         Right Memory                    │
│   (Tokens)              │         (Facts)                         │
├─────────────────────────┼─────────────────────────────────────────┤
│                         │                                         │
│  Token 1:               │    Fact 1:                              │
│    [Customer{id:c1}]    │      Order{id:o1, custId:c1}           │
│                         │                                         │
│  Token 2:               │    Fact 2:                              │
│    [Customer{id:c2}]    │      Order{id:o2, custId:c2}           │
│                         │                                         │
│  Token 3:               │    Fact 3:                              │
│    [Customer{id:c3}]    │      Order{id:o3, custId:c1}           │
│                         │                                         │
└─────────────────────────┴─────────────────────────────────────────┘
              │                            │
              └──────────┬─────────────────┘
                         │
                    [Join Tests]
              Customer.id == Order.customerId
                         │
                         ▼
               ┌─────────────────────┐
               │  Output Tokens:     │
               ├─────────────────────┤
               │  [Customer{c1},     │
               │   Order{o1}]        │
               │                     │
               │  [Customer{c2},     │
               │   Order{o2}]        │
               │                     │
               │  [Customer{c1},     │
               │   Order{o3}]        │
               └─────────────────────┘
```

### Diagramme 2 : Flux d'activation

```
Activation d'un JoinNode :

1. Activation Gauche (nouveau token)
   ┌─────────────┐
   │ New Token:  │
   │ [Customer]  │
   └──────┬──────┘
          │
          ▼
   ┌─────────────────┐
   │ Add to Left     │
   │ Memory          │
   └──────┬──────────┘
          │
          ▼
   ┌─────────────────────────────────┐
   │ For each fact in Right Memory:  │
   │   - Evaluate join tests         │
   │   - If match: create new token  │
   │   - Propagate to children       │
   └─────────────────────────────────┘

2. Activation Droite (nouveau fait)
   ┌─────────────┐
   │ New Fact:   │
   │ Order       │
   └──────┬──────┘
          │
          ▼
   ┌─────────────────┐
   │ Add to Right    │
   │ Memory          │
   └──────┬──────────┘
          │
          ▼
   ┌─────────────────────────────────┐
   │ For each token in Left Memory:  │
   │   - Evaluate join tests         │
   │   - If match: create new token  │
   │   - Propagate to children       │
   └─────────────────────────────────┘
```

### Diagramme 3 : Cycle de vie d'un JoinNode partagé

```
Lifecycle d'un JoinNode partagé :

Temps T0 : Création
═══════════════════
Rule1 added → JoinNode created (RefCount: 1)
┌─────────────┐
│  JoinNode   │
│ RefCount: 1 │
│ Rules: [R1] │
└─────────────┘

Temps T1 : Partage
═══════════════════
Rule2 added → JoinNode reused (RefCount: 2)
┌─────────────┐
│  JoinNode   │
│ RefCount: 2 │
│ Rules: [R1, │
│        R2]  │
└─────────────┘

Temps T2 : Plus de partage
═══════════════════
Rule3 added → JoinNode reused (RefCount: 3)
┌─────────────┐
│  JoinNode   │
│ RefCount: 3 │
│ Rules: [R1, │
│   R2, R3]   │
└─────────────┘

Temps T3 : Suppression partielle
═══════════════════
Rule2 removed → JoinNode kept (RefCount: 2)
┌─────────────┐
│  JoinNode   │
│ RefCount: 2 │
│ Rules: [R1, │
│        R3]  │
└─────────────┘

Temps T4 : Suppression complète
═══════════════════
Rule1 removed → JoinNode kept (RefCount: 1)
Rule3 removed → JoinNode deleted (RefCount: 0)
┌─────────────┐
│  JoinNode   │
│  DELETED    │
└─────────────┘
```

### Diagramme 4 : Cache de jointure

```
Join Cache : Évite la réévaluation des jointures

┌─────────────────────────────────────────────────────────────┐
│                    Join Cache (LRU)                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Key: hash(tokenID + factID + joinNodeID)                  │
│  Value: bool (match result)                                │
│                                                             │
│  ┌───────────────────────────────────────────────┐         │
│  │ "t1_f1_j1" → true   (HIT)                     │         │
│  │ "t1_f2_j1" → false  (HIT)                     │         │
│  │ "t2_f1_j1" → true   (HIT)                     │         │
│  │ "t2_f2_j1" → false  (MISS - evicted)          │         │
│  │ "t3_f1_j1" → true   (MISS - new entry)        │         │
│  └───────────────────────────────────────────────┘         │
│                                                             │
│  Stats:                                                     │
│    Hits: 3                                                  │
│    Misses: 2                                                │
│    Hit Rate: 60%                                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘

Workflow avec cache :

┌────────────────┐
│ Token + Fact   │
└───────┬────────┘
        │
        ▼
┌────────────────────┐      YES     ┌──────────────┐
│ Check cache        │─────────────►│ Use cached   │
│ (hash lookup)      │              │ result       │
└────────┬───────────┘              └──────┬───────┘
         │ NO                               │
         ▼                                  │
┌────────────────────┐                     │
│ Evaluate join      │                     │
│ tests              │                     │
└────────┬───────────┘                     │
         │                                  │
         ▼                                  │
┌────────────────────┐                     │
│ Store in cache     │                     │
└────────┬───────────┘                     │
         │                                  │
         └──────────────┬───────────────────┘
                        ▼
                ┌───────────────┐
                │ Propagate     │
                └───────────────┘
```

---

## Exemples visuels

### Exemple 1 : Réseau sans Beta Sharing

```
Règle 1: Customer → Order → Terminal1
Règle 2: Customer → Order → Terminal2
Règle 3: Customer → Order → Terminal3

Sans Beta Sharing :
════════════════════

     [Customer]    [Customer]    [Customer]
     AlphaNode     AlphaNode     AlphaNode
          │             │             │
     [Order]       [Order]       [Order]
     AlphaNode     AlphaNode     AlphaNode
          │             │             │
     [JoinNode]    [JoinNode]    [JoinNode]   ← 3 nœuds distincts
          │             │             │
     [Terminal1]   [Terminal2]   [Terminal3]

Mémoire : 3 × 10 KB = 30 KB
Calculs : 3 × évaluation jointure
```

### Exemple 2 : Réseau avec Beta Sharing

```
Règle 1: Customer → Order → Terminal1
Règle 2: Customer → Order → Terminal2
Règle 3: Customer → Order → Terminal3

Avec Beta Sharing :
═══════════════════

          [Customer]
          AlphaNode (PARTAGÉ)
               │
          [Order]
          AlphaNode (PARTAGÉ)
               │
          [JoinNode]          ← 1 seul nœud partagé
          (RefCount: 3)
               │
     ┌─────────┼─────────┐
     │         │         │
[Terminal1][Terminal2][Terminal3]

Mémoire : 1 × 10 KB = 10 KB (67% de réduction)
Calculs : 1 × évaluation jointure (67% de gain)
```

### Exemple 3 : Chaîne complexe avec préfixes partagés

```
Règle 1: Customer → Order → LineItem → Terminal1
Règle 2: Customer → Order → Payment → Terminal2
Règle 3: Customer → Order → Shipment → Terminal3

Avec Beta Sharing :
═══════════════════

                    [Customer]
                    AlphaNode
                        │
                    [Order]
                    AlphaNode
                        │
                   [JoinNode]
                 Customer-Order
                  (PARTAGÉ)
                        │
           ┌────────────┼────────────┐
           │            │            │
      [LineItem]   [Payment]   [Shipment]
      AlphaNode    AlphaNode   AlphaNode
           │            │            │
      [JoinNode]   [JoinNode]   [JoinNode]
      Order-       Order-       Order-
      LineItem     Payment      Shipment
           │            │            │
      [Terminal1]  [Terminal2]  [Terminal3]

Préfixe partagé : Customer → Order
Partage ratio : 33% (1 nœud partagé sur 3 créés)
```

### Exemple 4 : Patterns avec variations

```
Règle 1: Customer(premium) → Order(active) → Terminal1
Règle 2: Customer(premium) → Order(active, amount>100) → Terminal2
Règle 3: Customer(gold) → Order(active) → Terminal3

Avec Beta Sharing :
═══════════════════

        [Customer]              [Customer]
        premium                 gold
            │                       │
            ├───────────┐           │
            │           │           │
        [Order]     [Order]     [Order]
        active      active      active
                    amount>100
            │           │           │
       [JoinNode]  [JoinNode]  [JoinNode]
            │           │           │
       [Terminal1] [Terminal2] [Terminal3]

Partage ratio : 0% (patterns trop différents)
Amélioration possible : Refactoriser pour plus de partage
```

### Exemple 5 : Avant/Après optimisation

```
AVANT (règles non optimisées) :
════════════════════════════════

rule "R1"
when
    Order(status=="active", amount>100)
    Customer(id==$orderId.customerId)
then ...

rule "R2"
when
    Customer($custId : id)
    Order(customerId==$custId, status=="active", amount>100)
then ...

Réseau :
    [Order]              [Customer]
       │                     │
    [Customer]           [Order]
       │                     │
   [JoinNode]           [JoinNode]    ← 2 nœuds distincts
       │                     │
   [Terminal R1]        [Terminal R2]


APRÈS (règles optimisées) :
════════════════════════════

rule "R1"
when
    Customer($custId : id)
    Order(customerId==$custId, status=="active", amount>100)
then ...

rule "R2"
when
    Customer($custId : id)
    Order(customerId==$custId, status=="active", amount>100)
then ...

Réseau :
        [Customer]
            │
        [Order]
            │
       [JoinNode]          ← 1 seul nœud partagé
            │
      ┌─────┴─────┐
      │           │
  [Terminal]  [Terminal]
     R1          R2

Gain : 50% de nœuds en moins, 2× plus rapide
```

---

## Mécanismes internes

### 1. Calcul de la clé de partage

```
Algorithme de génération de clé :

Input: Pattern de jointure
Output: Clé de hash unique

Étapes :
1. Extraire les contraintes de jointure
2. Normaliser les contraintes (ordre canonique)
3. Calculer le hash (FNV-1a ou xxHash)
4. Retourner la clé

Exemple :

Pattern original :
  Customer($custId : id, type == "premium")
  Order(customerId == $custId, status == "active")

Extraction :
  - Test: Customer.id == Order.customerId
  - Pattern Customer: type == "premium"
  - Pattern Order: status == "active"

Normalisation :
  - Join: Customer.id == Order.customerId (ordre canonique)
  - Customer constraints: [type == "premium"] (ordre alpha)
  - Order constraints: [status == "active"] (ordre alpha)

Hash :
  FNV-1a(
    "Customer" + 
    "type==premium" + 
    "Order" + 
    "status==active" + 
    "Customer.id==Order.customerId"
  )
  = "a1b2c3d4e5f67890"

Clé finale : "join_a1b2c3d4e5f67890"
```

### 2. Registre de partage

```
BetaSharingRegistry :

┌──────────────────────────────────────────────────────────┐
│                  BetaSharingRegistry                     │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  joinNodes: map[string]*JoinNode                         │
│  ┌────────────────────────────────────────────────┐     │
│  │ "join_a1b2..." → JoinNode{                     │     │
│  │                    ID: "node_123",              │     │
│  │                    RefCount: 3,                 │     │
│  │                    RuleIDs: ["r1","r2","r3"]    │     │
│  │                  }                              │     │
│  │                                                 │     │
│  │ "join_c3d4..." → JoinNode{                     │     │
│  │                    ID: "node_456",              │     │
│  │                    RefCount: 1,                 │     │
│  │                    RuleIDs: ["r4"]              │     │
│  │                  }                              │     │
│  └────────────────────────────────────────────────┘     │
│                                                          │
│  hashCache: LRUCache (normalizedPattern → hash)         │
│  ┌────────────────────────────────────────────────┐     │
│  │ "Customer[...]|Order[...]" → "join_a1b2..."    │     │
│  │ "Account[...]|Transaction[...]" → "join_c3d4..."│    │
│  └────────────────────────────────────────────────┘     │
│                                                          │
│  Méthodes:                                               │
│    - GetOrCreateJoinNode(key, pattern, parent)           │
│    - ComputeJoinKey(pattern) → key                       │
│    - GetSharingStats() → stats                           │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

### 3. Séquence de construction

```
Construction d'une règle avec Beta Sharing :

┌───────────────────┐
│ Parse rule        │
│ (TSD → AST)       │
└────────┬──────────┘
         │
         ▼
┌───────────────────┐
│ Extract patterns  │
│ [Customer, Order] │
└────────┬──────────┘
         │
         ▼
┌───────────────────────────┐
│ For each pattern pair:    │
│   1. Normalize pattern    │
│   2. Compute join key     │
│   3. Check registry       │
└────────┬──────────────────┘
         │
         ▼
┌────────────────────┐   YES   ┌──────────────────┐
│ JoinNode exists?   │────────►│ Reuse node       │
└────────┬───────────┘         │ RefCount++       │
         │ NO                  └────────┬─────────┘
         ▼                              │
┌────────────────────┐                 │
│ Create new         │                 │
│ JoinNode           │                 │
│ RefCount = 1       │                 │
└────────┬───────────┘                 │
         │                              │
         └──────────────┬───────────────┘
                        │
                        ▼
                ┌───────────────┐
                │ Register in   │
                │ registry      │
                └───────┬───────┘
                        │
                        ▼
                ┌───────────────┐
                │ Connect to    │
                │ terminal      │
                └───────────────┘
```

### 4. Gestion de la concurrence

```
Thread-safety du Beta Sharing :

Scénario : Deux règles ajoutées simultanément

Thread 1                          Thread 2
────────                          ────────
AddRule("Rule1")                  AddRule("Rule2")
    │                                 │
    ▼                                 ▼
ComputeJoinKey()                  ComputeJoinKey()
    │                                 │
    ▼                                 ▼
registry.Lock()                   registry.Lock() ⏳ (wait)
    │                                 │
    ▼                                 │
Check joinNodes map                  │
    │                                 │
    ▼                                 │
Create new JoinNode                  │
    │                                 │
    ▼                                 │
Store in map                         │
    │                                 │
    ▼                                 │
registry.Unlock()                    │
                                     ▼
                              registry.Lock() ✓ (acquired)
                                     │
                                     ▼
                              Check joinNodes map
                                     │
                                     ▼
                              Find existing node
                                     │
                                     ▼
                              RefCount++ (atomic)
                                     │
                                     ▼
                              registry.Unlock()

Garanties :
✓ Pas de race conditions
✓ Pas de créations dupliquées
✓ RefCount toujours cohérent
```

---

## Licence

Copyright (c) 2024 TSD Project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINF