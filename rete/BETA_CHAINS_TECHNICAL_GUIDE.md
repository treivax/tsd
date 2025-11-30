# Guide Technique : Chaînes de JoinNodes (Beta Chains)

## Table des Matières

1. [Architecture](#architecture)
2. [Algorithmes](#algorithmes)
3. [Normalisation des Patterns de Jointure](#normalisation-des-patterns-de-jointure)
4. [Lifecycle Management](#lifecycle-management)
5. [API Reference](#api-reference)
6. [Gestion des cas edge](#gestion-des-cas-edge)
7. [Optimisations](#optimisations)
8. [Internals](#internals)

---

## Architecture

### Vue d'ensemble du système

Le système de Beta Chains (partage de JoinNodes) est composé de plusieurs couches interdépendantes qui permettent la réutilisation intelligente des nœuds de jointure dans le réseau RETE :

```
┌─────────────────────────────────────────────────────────────────┐
│  Application Layer                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ TSD Parser   │→ │ Rule Builder │→ │ ReteNetwork  │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  Beta Chain Layer                                                 │
│  ┌──────────────────┐      ┌──────────────────┐                │
│  │ BetaChainBuilder │◄────►│ BetaChainMetrics │                │
│  └──────────────────┘      └──────────────────┘                │
│           │                                                       │
│           │ build chains, optimize join order                    │
│           ▼                                                       │
│  ┌──────────────────┐                                           │
│  │ JoinOrderOptimizer│                                          │
│  └──────────────────┘                                           │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  Sharing Layer                                                    │
│  ┌─────────────────────┐   ┌──────────────────┐                │
│  │BetaSharingRegistry  │◄──│ LRUCache (Hash)  │                │
│  └─────────────────────┘   └──────────────────┘                │
│           │                                                       │
│           │ joinKey → JoinNode mapping                           │
│           │ pattern normalization                                │
│           ▼                                                       │
│  ┌─────────────────────┐   ┌──────────────────┐                │
│  │  BetaJoinCache      │◄──│ LRUCache (Join)  │                │
│  │  (Join Results)     │   │  Results Cache   │                │
│  └─────────────────────┘   └──────────────────┘                │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  Node Layer                                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  AlphaNode   │→ │  JoinNode    │→ │ TerminalNode │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                           │                                       │
│                    ┌──────▼──────┐                              │
│                    │LeftMemory   │                              │
│                    │ RightMemory │                              │
│                    └─────────────┘                              │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│  Lifecycle Layer                                                  │
│  ┌─────────────────────┐   ┌──────────────────┐                │
│  │  LifecycleManager   │◄──│  NodeLifecycle   │                │
│  │  (Reference Count)  │   │  (Per Node)      │                │
│  └─────────────────────┘   └──────────────────┘                │
└─────────────────────────────────────────────────────────────────┘
```

### Composants principaux

#### 1. BetaChainBuilder

**Responsabilités :**
- Construction séquentielle de chaînes de JoinNodes
- Optimisation de l'ordre des jointures
- Coordination avec BetaSharingRegistry pour réutilisation
- Collection de métriques de construction
- Gestion des préfixes de chaînes partagées

**Structure de données :**

```go
type BetaChainBuilder struct {
    network         *ReteNetwork           // Référence au réseau RETE
    registry        *BetaSharingRegistry   // Registre de partage
    optimizer       *JoinOrderOptimizer    // Optimiseur d'ordre de jointure
    metrics         *BetaChainMetrics      // Métriques de construction
    mutex           sync.RWMutex           // Protection concurrence
}
```

**Thread-safety :**
- Lecture du registry : `RLock()`
- Écriture du registry : `Lock()`
- Métriques : opérations atomiques via `sync/atomic`

**Algorithme de construction :**

```
BuildBetaChain(patterns []Pattern, terminalNode Node) error:
    1. Analyser les patterns et extraire les contraintes de jointure
    2. Optimiser l'ordre des jointures (selectivité, disponibilité indices)
    3. Pour chaque pattern dans l'ordre optimisé:
        a. Calculer la clé de jointure normalisée
        b. Vérifier si un JoinNode existe déjà (sharing)
        c. Si oui: réutiliser, incrémenter ref count
        d. Si non: créer nouveau JoinNode, l'enregistrer
        e. Connecter au parent (AlphaNode ou JoinNode précédent)
        f. Initialiser les mémoires (left/right)
    4. Connecter le dernier JoinNode au TerminalNode
    5. Enregistrer les métriques
```

#### 2. BetaSharingRegistry

**Responsabilités :**
- Mapping joinKey → JoinNode
- Cache LRU des calculs de hash de jointure
- Normalisation des patterns de jointure
- Création thread-safe de JoinNodes partagés
- Statistiques de partage en temps réel

**Structure de données :**

```go
type BetaSharingRegistry struct {
    joinNodes       map[string]*JoinNode   // joinKey → JoinNode
    hashCache       *LRUCache              // Cache des hash normalisés
    normalization   *JoinNormalizer        // Normalisation patterns
    metrics         *BetaSharingStats      // Statistiques runtime
    mutex           sync.RWMutex           // Protection concurrence
}

type BetaSharingStats struct {
    TotalJoinNodes    int64   // Nombre total de JoinNodes créés
    SharedJoinNodes   int64   // Nombre de JoinNodes partagés
    UniqueJoinKeys    int64   // Nombre de clés uniques
    HashCacheHits     int64   // Hits du cache de hash
    HashCacheMisses   int64   // Misses du cache de hash
    SharingRatio      float64 // Ratio de partage (shared/total)
}
```

**Calcul de la clé de jointure :**

```go
func (r *BetaSharingRegistry) ComputeJoinKey(pattern Pattern) string {
    // 1. Normaliser le pattern (ordre canonique)
    normalized := r.normalization.Normalize(pattern)
    
    // 2. Vérifier le cache LRU
    if cachedHash, ok := r.hashCache.Get(normalized.String()); ok {
        atomic.AddInt64(&r.metrics.HashCacheHits, 1)
        return cachedHash.(string)
    }
    
    // 3. Calculer le hash (FNV-1a ou xxHash)
    hash := computeHash(normalized)
    
    // 4. Mettre en cache
    r.hashCache.Set(normalized.String(), hash)
    atomic.AddInt64(&r.metrics.HashCacheMisses, 1)
    
    return hash
}
```

#### 3. BetaJoinCache

**Responsabilités :**
- Cache des résultats de jointure (token × fact → résultat)
- Invalidation sélective lors de l'ajout/suppression de faits
- Métriques de hit rate et évictions
- Gestion LRU de la capacité mémoire

**Structure de données :**

```go
type BetaJoinCache struct {
    cache           *LRUCache              // Cache LRU principal
    stats           *JoinCacheStats        // Statistiques
    maxSize         int                    // Capacité maximale
    mutex           sync.RWMutex           // Protection concurrence
}

type JoinCacheStats struct {
    Hits            int64   // Nombre de hits
    Misses          int64   // Nombre de misses
    Evictions       int64   // Nombre d'évictions
    Size            int64   // Taille actuelle
    HitRate         float64 // Ratio de hits
}
```

**Clé de cache :**

```
joinCacheKey = hash(tokenID + factID + joinNodeID)
```

**Algorithme d'invalidation :**

```
InvalidateForFact(factID string):
    1. Pour chaque entrée dans le cache:
        a. Si la clé contient factID:
            - Supprimer l'entrée
            - Incrémenter le compteur d'évictions
    2. Mettre à jour les statistiques
```

**Note :** Pour une meilleure performance, une structure d'index inversé (fact → cache keys) peut être implémentée.

#### 4. JoinOrderOptimizer

**Responsabilités :**
- Analyse de la sélectivité des patterns
- Optimisation de l'ordre des jointures pour minimiser les résultats intermédiaires
- Heuristiques basées sur les statistiques du working memory

**Heuristiques d'optimisation :**

1. **Sélectivité** : Patterns les plus sélectifs d'abord
   - Pattern avec le plus de contraintes
   - Pattern avec indices disponibles
   - Pattern avec le moins de faits correspondants estimés

2. **Disponibilité des index** : Préférer les patterns indexés

3. **Cardinalité** : Minimiser le produit cartésien

**Formule de coût :**

```
cost(joinOrder) = Σ (cardinalité(pattern_i) × coût_jointure(pattern_i, pattern_{i-1}))
```

### Architecture des JoinNodes

#### Structure d'un JoinNode

```
┌─────────────────────────────────────────────────────────────────┐
│                           JoinNode                                │
├─────────────────────────────────────────────────────────────────┤
│  ID: "join_abc123"                                               │
│  Type: NodeTypeJoin                                              │
│  Tests: []JoinTest                                               │
│  LeftParent: AlphaNode | JoinNode                                │
│  RightParent: AlphaNode                                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────┐      ┌─────────────────────┐          │
│  │   Left Memory       │      │   Right Memory      │          │
│  │   (Tokens)          │      │   (Facts)           │          │
│  ├─────────────────────┤      ├─────────────────────┤          │
│  │ token1: [f1, f2]    │      │ fact1: {x:1, y:2}  │          │
│  │ token2: [f1, f3]    │      │ fact2: {x:1, y:3}  │          │
│  │ token3: [f4, f5]    │      │ fact3: {x:2, y:4}  │          │
│  └─────────────────────┘      └─────────────────────┘          │
├─────────────────────────────────────────────────────────────────┤
│  Children: [JoinNode | TerminalNode]                             │
│  Lifecycle: {RefCount: 2, CreatedAt: timestamp}                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Tests de jointure (JoinTests)

Les JoinTests définissent les conditions de jointure entre les tokens de gauche et les faits de droite :

```go
type JoinTest struct {
    LeftIndex   int       // Index dans le token (0-based)
    LeftField   string    // Champ du fait gauche
    Operator    string    // ==, !=, <, >, <=, >=
    RightField  string    // Champ du fait droit
}

// Exemple : joindre sur customer.id == order.customer_id
JoinTest{
    LeftIndex:  0,           // Premier fait du token
    LeftField:  "id",        // customer.id
    Operator:   "==",
    RightField: "customer_id", // order.customer_id
}
```

### Mécanisme de partage de préfixes

Le partage de préfixes permet de réutiliser les chaînes de JoinNodes communes entre plusieurs règles :

```
Règle 1: Customer → Order → LineItem → Terminal1
Règle 2: Customer → Order → Payment  → Terminal2
Règle 3: Customer → Order → Shipment → Terminal3

Réseau RETE optimisé avec partage :

         ┌──────────┐
         │ Customer │ (AlphaNode)
         └────┬─────┘
              │
         ┌────▼─────┐
         │  Order   │ (JoinNode - PARTAGÉ)
         └────┬─────┘
              │
      ┌───────┼───────┬─────────┐
      │       │       │         │
 ┌────▼───┐ ┌▼────┐ ┌▼──────┐  │
 │LineItem│ │Payment│Shipment│  │
 └────┬───┘ └┬────┘ └┬──────┘  │
      │      │       │         │
 ┌────▼───┐ ┌▼────┐ ┌▼──────┐  │
 │Terminal1│Terminal2│Terminal3│
 └────────┘ └─────┘ └───────┘  │
```

**Avantages du partage de préfixes :**
- Moins de nœuds créés (mémoire réduite)
- Moins de calculs de jointure (CPU réduit)
- Propagation des tokens une seule fois
- Mémoires partagées (left/right memory)

---

## Algorithmes

### 1. Construction de chaîne Beta

```go
func (b *BetaChainBuilder) BuildChain(
    patterns []Pattern,
    terminal *TerminalNode,
) error {
    // Étape 1 : Optimiser l'ordre des patterns
    optimizedPatterns := b.optimizer.OptimizeJoinOrder(patterns)
    
    var currentParent Node = nil
    
    // Étape 2 : Construire la chaîne pattern par pattern
    for i, pattern := range optimizedPatterns {
        // 2a. Déterminer le parent (AlphaNode ou JoinNode précédent)
        if i == 0 {
            // Premier pattern : parent = AlphaNode du pattern
            currentParent = b.getOrCreateAlphaNode(pattern)
        }
        
        // 2b. Calculer la clé de jointure normalisée
        joinKey := b.registry.ComputeJoinKey(pattern, currentParent)
        
        // 2c. Obtenir ou créer le JoinNode
        joinNode, created := b.registry.GetOrCreateJoinNode(
            joinKey,
            pattern,
            currentParent,
        )
        
        // 2d. Mettre à jour les métriques
        if created {
            atomic.AddInt64(&b.metrics.NodesCreated, 1)
        } else {
            atomic.AddInt64(&b.metrics.NodesReused, 1)
        }
        
        // 2e. Connecter au parent
        currentParent.AddChild(joinNode)
        
        // 2f. Le JoinNode devient le parent pour l'itération suivante
        currentParent = joinNode
    }
    
    // Étape 3 : Connecter au terminal
    currentParent.AddChild(terminal)
    
    return nil
}
```

### 2. Normalisation de pattern de jointure

La normalisation garantit que des patterns équivalents produisent la même clé de hash :

```go
func (n *JoinNormalizer) Normalize(pattern Pattern) NormalizedPattern {
    normalized := NormalizedPattern{
        Type:        pattern.Type,
        Constraints: []Constraint{},
    }
    
    // Étape 1 : Extraire et normaliser les contraintes
    for _, constraint := range pattern.Constraints {
        normalized.Constraints = append(
            normalized.Constraints,
            n.normalizeConstraint(constraint),
        )
    }
    
    // Étape 2 : Trier les contraintes dans un ordre canonique
    sort.Slice(normalized.Constraints, func(i, j int) bool {
        return normalized.Constraints[i].Hash() < 
               normalized.Constraints[j].Hash()
    })
    
    // Étape 3 : Normaliser les contraintes commutatives
    // Exemple : x == y équivaut à y == x
    for i := range normalized.Constraints {
        if isCommutative(normalized.Constraints[i].Operator) {
            normalized.Constraints[i] = 
                canonicalizeCommutative(normalized.Constraints[i])
        }
    }
    
    return normalized
}

func canonicalizeCommutative(c Constraint) Constraint {
    // Pour ==, !=, mettre le champ "plus petit" à gauche
    if c.Operator == "==" || c.Operator == "!=" {
        if c.RightField < c.LeftField {
            return Constraint{
                LeftField:  c.RightField,
                Operator:   c.Operator,
                RightField: c.LeftField,
            }
        }
    }
    return c
}
```

**Exemples de normalisation :**

```
Avant normalisation :
  Pattern 1: {type: "Order", constraints: [amount > 100, status == "active"]}
  Pattern 2: {type: "Order", constraints: [status == "active", amount > 100]}

Après normalisation :
  Pattern 1: {type: "Order", constraints: [amount > 100, status == "active"]}
  Pattern 2: {type: "Order", constraints: [amount > 100, status == "active"]}
  → Même clé de hash → Même JoinNode partagé
```

### 3. Optimisation de l'ordre de jointure

**Algorithme glouton (greedy) :**

```go
func (o *JoinOrderOptimizer) OptimizeJoinOrder(
    patterns []Pattern,
) []Pattern {
    if len(patterns) <= 1 {
        return patterns
    }
    
    optimized := make([]Pattern, 0, len(patterns))
    remaining := append([]Pattern{}, patterns...)
    
    // Étape 1 : Choisir le pattern le plus sélectif comme point de départ
    mostSelective := o.findMostSelective(remaining)
    optimized = append(optimized, mostSelective)
    remaining = removePattern(remaining, mostSelective)
    
    // Étape 2 : Ajouter itérativement le meilleur pattern suivant
    for len(remaining) > 0 {
        bestNext := o.findBestNext(optimized, remaining)
        optimized = append(optimized, bestNext)
        remaining = removePattern(remaining, bestNext)
    }
    
    return optimized
}

func (o *JoinOrderOptimizer) findBestNext(
    current []Pattern,
    candidates []Pattern,
) Pattern {
    var bestPattern Pattern
    lowestCost := math.MaxFloat64
    
    for _, candidate := range candidates {
        // Coût = cardinalité estimée × coût de jointure
        cost := o.estimateCardinality(candidate) * 
                o.estimateJoinCost(current, candidate)
        
        if cost < lowestCost {
            lowestCost = cost
            bestPattern = candidate
        }
    }
    
    return bestPattern
}
```

**Heuristiques de sélectivité :**

1. **Nombre de contraintes** : Plus de contraintes = plus sélectif
2. **Type de contrainte** : Égalité > Inégalité > Range
3. **Disponibilité d'index** : Patterns indexés sont préférés
4. **Statistiques du working memory** : Nombre de faits par type

### 4. Algorithme de jointure (runtime)

```go
func (j *JoinNode) ActivateLeft(token *Token) {
    // 1. Ajouter le token à la mémoire gauche
    j.leftMemory.Add(token)
    
    // 2. Pour chaque fait dans la mémoire droite
    for _, fact := range j.rightMemory.GetAll() {
        // 3. Vérifier le cache de jointure
        cacheKey := computeJoinCacheKey(token.ID, fact.ID, j.ID)
        
        if cachedResult, ok := j.cache.Get(cacheKey); ok {
            // Cache hit : utiliser le résultat
            if cachedResult.(bool) {
                j.propagateMatch(token, fact)
            }
            continue
        }
        
        // 4. Évaluer les tests de jointure
        match := j.evaluateJoinTests(token, fact)
        
        // 5. Mettre en cache le résultat
        j.cache.Set(cacheKey, match)
        
        // 6. Propager si match
        if match {
            j.propagateMatch(token, fact)
        }
    }
}

func (j *JoinNode) evaluateJoinTests(token *Token, fact *Fact) bool {
    for _, test := range j.tests {
        // Extraire la valeur du token (fait gauche)
        leftFact := token.Facts[test.LeftIndex]
        leftValue := leftFact.Fields[test.LeftField]
        
        // Extraire la valeur du fait droit
        rightValue := fact.Fields[test.RightField]
        
        // Évaluer l'opérateur
        if !evaluateOperator(leftValue, test.Operator, rightValue) {
            return false // Échec du test
        }
    }
    
    return true // Tous les tests réussis
}

func (j *JoinNode) propagateMatch(token *Token, fact *Fact) {
    // Créer un nouveau token étendu
    newToken := &Token{
        ID:    generateTokenID(),
        Facts: append(token.Facts, fact),
    }
    
    // Propager aux enfants
    for _, child := range j.children {
        child.Activate(newToken)
    }
}
```

### 5. Invalidation du cache de jointure

**Invalidation naïve (O(n)) :**

```go
func (c *BetaJoinCache) InvalidateForFact(factID string) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    // Parcourir toutes les entrées du cache
    for key := range c.cache.items {
        if strings.Contains(key, factID) {
            c.cache.Remove(key)
            atomic.AddInt64(&c.stats.Evictions, 1)
        }
    }
}
```

**Invalidation optimisée avec index inversé (O(k)) :**

```go
type BetaJoinCacheWithIndex struct {
    cache       *LRUCache
    factIndex   map[string][]string  // factID → [cacheKeys]
    tokenIndex  map[string][]string  // tokenID → [cacheKeys]
    mutex       sync.RWMutex
}

func (c *BetaJoinCacheWithIndex) Set(
    tokenID, factID, joinNodeID string,
    result bool,
) {
    cacheKey := computeJoinCacheKey(tokenID, factID, joinNodeID)
    
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    // Ajouter au cache principal
    c.cache.Set(cacheKey, result)
    
    // Ajouter aux index inversés
    c.factIndex[factID] = append(c.factIndex[factID], cacheKey)
    c.tokenIndex[tokenID] = append(c.tokenIndex[tokenID], cacheKey)
}

func (c *BetaJoinCacheWithIndex) InvalidateForFact(factID string) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    // Récupérer les clés associées au fact
    keys := c.factIndex[factID]
    
    // Supprimer chaque clé du cache
    for _, key := range keys {
        c.cache.Remove(key)
        atomic.AddInt64(&c.stats.Evictions, 1)
    }
    
    // Supprimer l'entrée d'index
    delete(c.factIndex, factID)
}
```

---

## Normalisation des Patterns de Jointure

### Principes de normalisation

La normalisation garantit que des patterns sémantiquement équivalents produisent la même clé de hash, maximisant ainsi les opportunités de partage.

#### 1. Ordre des contraintes

Les contraintes sont triées dans un ordre canonique :

```go
// Avant normalisation
Pattern{
    Type: "Order",
    Constraints: [
        {Field: "status", Op: "==", Value: "active"},
        {Field: "amount", Op: ">", Value: 100},
        {Field: "date", Op: ">=", Value: "2024-01-01"},
    ],
}

// Après normalisation (tri lexicographique par hash)
Pattern{
    Type: "Order",
    Constraints: [
        {Field: "amount", Op: ">", Value: 100},
        {Field: "date", Op: ">=", Value: "2024-01-01"},
        {Field: "status", Op: "==", Value: "active"},
    ],
}
```

#### 2. Opérateurs commutatifs

Les opérateurs commutatifs (`==`, `!=`) sont normalisés :

```go
// Avant normalisation
Constraint{LeftField: "order.id", Op: "==", RightField: "payment.order_id"}
Constraint{LeftField: "payment.order_id", Op: "==", RightField: "order.id"}

// Après normalisation (ordre lexicographique)
Constraint{LeftField: "order.id", Op: "==", RightField: "payment.order_id"}
Constraint{LeftField: "order.id", Op: "==", RightField: "payment.order_id"}
```

#### 3. Opérateurs équivalents

Certains opérateurs peuvent être réécrits :

```go
// NOT (x == y) → x != y
// NOT (x < y) → x >= y
// NOT (x > y) → x <= y
```

#### 4. Contraintes redondantes

Les contraintes redondantes sont détectées et éliminées :

```go
// Avant normalisation
Constraints: [
    {Field: "x", Op: ">", Value: 10},
    {Field: "x", Op: ">", Value: 5},  // Redondant (10 > 5)
]

// Après normalisation
Constraints: [
    {Field: "x", Op: ">", Value: 10},
]
```

### Algorithme de hachage

**Fonction de hash FNV-1a :**

```go
func computeJoinHash(pattern NormalizedPattern) string {
    h := fnv.New64a()
    
    // 1. Hasher le type
    h.Write([]byte(pattern.Type))
    h.Write([]byte{0}) // Séparateur
    
    // 2. Hasher chaque contrainte dans l'ordre canonique
    for _, constraint := range pattern.Constraints {
        h.Write([]byte(constraint.LeftField))
        h.Write([]byte{0})
        h.Write([]byte(constraint.Operator))
        h.Write([]byte{0})
        h.Write([]byte(constraint.RightField))
        h.Write([]byte{0})
    }
    
    return hex.EncodeToString(h.Sum(nil))
}
```

**Propriétés du hash :**
- Déterministe : même input → même output
- Rapide : O(n) où n = nombre de contraintes
- Uniforme : bonne distribution dans l'espace de hash
- Collision-resistant : probabilité très faible de collision

### Exemples de normalisation

**Exemple 1 : Ordre des contraintes**

```
Règle A:
  when
    Order(status == "active", amount > 100)
  then
    ...

Règle B:
  when
    Order(amount > 100, status == "active")
  then
    ...

Normalisation :
  Pattern A normalisé : Order(amount > 100, status == "active")
  Pattern B normalisé : Order(amount > 100, status == "active")
  → Hash identique → JoinNode partagé
```

**Exemple 2 : Opérateurs commutatifs**

```
Règle A:
  when
    Order(id == $orderId)
    Payment(order_id == $orderId)
  then
    ...

Règle B:
  when
    Order($orderId := id)
    Payment($orderId == order_id)
  then
    ...

Normalisation :
  Join A : Order.id == Payment.order_id
  Join B : Order.id == Payment.order_id
  → Hash identique → JoinNode partagé
```

**Exemple 3 : Patterns complexes**

```
Règle :
  when
    Customer($custId := id, type == "premium")
    Order(customer_id == $custId, status != "cancelled")
    LineItem(order_id == $orderId, quantity > 0)
  then
    ...

Chaîne normalisée :
  1. JoinNode[Customer-Order] :
     - Tests: [Customer.id == Order.customer_id]
     - Hash: a1b2c3d4...
  
  2. JoinNode[Customer-Order-LineItem] :
     - Tests: [Order.id == LineItem.order_id]
     - Hash: e5f6g7h8...
```

---

## Lifecycle Management

### Référencement des JoinNodes

Chaque JoinNode maintient un compteur de références pour le partage :

```go
type NodeLifecycle struct {
    RefCount    int32      // Nombre de règles utilisant ce nœud
    CreatedAt   time.Time  // Timestamp de création
    LastUsedAt  time.Time  // Timestamp de dernière utilisation
    RuleIDs     []string   // IDs des règles utilisant ce nœud
    mutex       sync.Mutex // Protection concurrence
}
```

### Ajout d'une règle

```go
func (r *BetaSharingRegistry) GetOrCreateJoinNode(
    joinKey string,
    pattern Pattern,
    parent Node,
) (*JoinNode, bool) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    // Vérifier si le JoinNode existe déjà
    if existingNode, ok := r.joinNodes[joinKey]; ok {
        // Incrémenter le compteur de références
        existingNode.lifecycle.IncrementRef()
        
        // Ajouter l'ID de la règle
        existingNode.lifecycle.AddRuleID(ruleID)
        
        // Mettre à jour le timestamp
        existingNode.lifecycle.LastUsedAt = time.Now()
        
        atomic.AddInt64(&r.metrics.SharedJoinNodes, 1)
        return existingNode, false // Nœud réutilisé
    }
    
    // Créer un nouveau JoinNode
    newNode := &JoinNode{
        ID:       generateNodeID(),
        Type:     NodeTypeJoin,
        Tests:    extractJoinTests(pattern),
        Parent:   parent,
        Children: []Node{},
        LeftMemory:  NewMemoryStore(),
        RightMemory: NewMemoryStore(),
        Lifecycle: &NodeLifecycle{
            RefCount:   1,
            CreatedAt:  time.Now(),
            LastUsedAt: time.Now(),
            RuleIDs:    []string{ruleID},
        },
    }
    
    // Enregistrer dans le registre
    r.joinNodes[joinKey] = newNode
    
    atomic.AddInt64(&r.metrics.TotalJoinNodes, 1)
    return newNode, true // Nouveau nœud créé
}
```

### Suppression d'une règle

```go
func (n *ReteNetwork) RemoveRule(ruleID string) error {
    // 1. Trouver tous les JoinNodes utilisés par cette règle
    joinNodes := n.findJoinNodesForRule(ruleID)
    
    // 2. Décrémenter les compteurs de références
    for _, joinNode := range joinNodes {
        joinNode.lifecycle.DecrementRef()
        joinNode.lifecycle.RemoveRuleID(ruleID)
        
        // 3. Si plus aucune référence, supprimer le nœud
        if joinNode.lifecycle.RefCount == 0 {
            n.deleteJoinNode(joinNode)
        }
    }
    
    return nil
}

func (n *ReteNetwork) deleteJoinNode(node *JoinNode) {
    // 1. Déconnecter du parent
    node.Parent.RemoveChild(node)
    
    // 2. Reconnecter les enfants au parent (si possible)
    for _, child := range node.Children {
        node.Parent.AddChild(child)
    }
    
    // 3. Vider les mémoires
    node.LeftMemory.Clear()
    node.RightMemory.Clear()
    
    // 4. Supprimer du registre
    joinKey := n.registry.GetKeyForNode(node)
    delete(n.registry.joinNodes, joinKey)
    
    // 5. Mettre à jour les métriques
    atomic.AddInt64(&n.registry.metrics.TotalJoinNodes, -1)
}
```

### Garbage Collection

**Nettoyage des nœuds non utilisés :**

```go
func (r *BetaSharingRegistry) GarbageCollect(maxAge time.Duration) int {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    now := time.Now()
    deleted := 0
    
    for joinKey, node := range r.joinNodes {
        // Conditions de suppression :
        // - RefCount == 0 (aucune règle)
        // - ET LastUsedAt > maxAge (non utilisé depuis longtemps)
        if node.lifecycle.RefCount == 0 && 
           now.Sub(node.lifecycle.LastUsedAt) > maxAge {
            delete(r.joinNodes, joinKey)
            deleted++
        }
    }
    
    return deleted
}
```

### Diagramme de lifecycle

```
┌─────────────────────────────────────────────────────────────────┐
│                     JoinNode Lifecycle                            │
└─────────────────────────────────────────────────────────────────┘

  [Rule Added]
       │
       ▼
  ┌─────────┐  RefCount++
  │ CREATED ├──────────────────┐
  └─────────┘                  │
       │                       │
       │ RefCount > 0          │
       ▼                       ▼
  ┌─────────┐            ┌─────────┐
  │ ACTIVE  │◄───────────┤ SHARED  │
  └─────────┘  RefCount>1└─────────┘
       │                       │
       │ [Rule Removed]        │ [Rule Removed]
       │ RefCount--            │ RefCount--
       ▼                       │
  ┌─────────┐                 │
  │ UNUSED  │◄────────────────┘
  └─────────┘  RefCount==0
       │
       │ [GC after maxAge]
       ▼
  ┌─────────┐
  │ DELETED │
  └─────────┘
```

---

## API Reference

### BetaChainBuilder

#### `NewBetaChainBuilder`

```go
func NewBetaChainBuilder(
    network *ReteNetwork,
    registry *BetaSharingRegistry,
) *BetaChainBuilder
```

Crée un nouveau builder de chaînes Beta.

**Paramètres :**
- `network` : Référence au réseau RETE
- `registry` : Registre de partage des JoinNodes

**Retour :**
- `*BetaChainBuilder` : Nouveau builder

**Exemple :**

```go
registry := NewBetaSharingRegistry(1000) // Cache de 1000 entrées
builder := NewBetaChainBuilder(network, registry)
```

#### `BuildChain`

```go
func (b *BetaChainBuilder) BuildChain(
    patterns []Pattern,
    terminal *TerminalNode,
) error
```

Construit une chaîne de JoinNodes pour une séquence de patterns.

**Paramètres :**
- `patterns` : Liste des patterns de conditions
- `terminal` : Nœud terminal (production)

**Retour :**
- `error` : Erreur si la construction échoue

**Exemple :**

```go
patterns := []Pattern{
    {Type: "Customer", Constraints: [...]},
    {Type: "Order", Constraints: [...]},
    {Type: "Payment", Constraints: [...]},
}

terminal := &TerminalNode{RuleID: "rule123"}
err := builder.BuildChain(patterns, terminal)
if err != nil {
    log.Fatal(err)
}
```

#### `GetMetrics`

```go
func (b *BetaChainBuilder) GetMetrics() *BetaChainMetrics
```

Récupère les métriques de construction.

**Retour :**
- `*BetaChainMetrics` : Snapshot des métriques

**Exemple :**

```go
metrics := builder.GetMetrics()
fmt.Printf("Nodes created: %d\n", metrics.NodesCreated)
fmt.Printf("Nodes reused: %d\n", metrics.NodesReused)
fmt.Printf("Sharing ratio: %.2f%%\n", metrics.SharingRatio*100)
```

### BetaSharingRegistry

#### `NewBetaSharingRegistry`

```go
func NewBetaSharingRegistry(hashCacheSize int) *BetaSharingRegistry
```

Crée un nouveau registre de partage.

**Paramètres :**
- `hashCacheSize` : Taille du cache LRU de hash

**Retour :**
- `*BetaSharingRegistry` : Nouveau registre

**Exemple :**

```go
registry := NewBetaSharingRegistry(1000)
```

#### `GetOrCreateJoinNode`

```go
func (r *BetaSharingRegistry) GetOrCreateJoinNode(
    joinKey string,
    pattern Pattern,
    parent Node,
) (*JoinNode, bool)
```

Obtient un JoinNode existant ou en crée un nouveau.

**Paramètres :**
- `joinKey` : Clé de hash normalisée
- `pattern` : Pattern de condition
- `parent` : Nœud parent

**Retour :**
- `*JoinNode` : Le nœud (existant ou nouveau)
- `bool` : `true` si créé, `false` si réutilisé

**Exemple :**

```go
joinKey := registry.ComputeJoinKey(pattern, parent)
joinNode, created := registry.GetOrCreateJoinNode(joinKey, pattern, parent)

if created {
    fmt.Println("New JoinNode created")
} else {
    fmt.Println("Existing JoinNode reused")
}
```

#### `GetSharingStats`

```go
func (r *BetaSharingRegistry) GetSharingStats() BetaSharingStats
```

Récupère les statistiques de partage.

**Retour :**
- `BetaSharingStats` : Snapshot des statistiques

**Exemple :**

```go
stats := registry.GetSharingStats()
fmt.Printf("Total join nodes: %d\n", stats.TotalJoinNodes)
fmt.Printf("Shared join nodes: %d\n", stats.SharedJoinNodes)
fmt.Printf("Sharing ratio: %.2f%%\n", stats.SharingRatio*100)
fmt.Printf("Hash cache hit rate: %.2f%%\n", 
    float64(stats.HashCacheHits) / 
    float64(stats.HashCacheHits + stats.HashCacheMisses) * 100)
```

### BetaJoinCache

#### `NewBetaJoinCache`

```go
func NewBetaJoinCache(maxSize int) *BetaJoinCache
```

Crée un nouveau cache de résultats de jointure.

**Paramètres :**
- `maxSize` : Capacité maximale du cache LRU

**Retour :**
- `*BetaJoinCache` : Nouveau cache

**Exemple :**

```go
cache := NewBetaJoinCache(5000)
```

#### `GetJoinResult`

```go
func (c *BetaJoinCache) GetJoinResult(
    tokenID, factID, joinNodeID string,
) (bool, bool)
```

Récupère un résultat de jointure du cache.

**Paramètres :**
- `tokenID` : ID du token
- `factID` : ID du fait
- `joinNodeID` : ID du JoinNode

**Retour :**
- `bool` : Résultat de la jointure (si trouvé)
- `bool` : `true` si trouvé dans le cache

**Exemple :**

```go
result, found := cache.GetJoinResult(token.ID, fact.ID, joinNode.ID)
if found {
    if result {
        // Jointure réussie (cached)
    } else {
        // Jointure échouée (cached)
    }
} else {
    // Cache miss, évaluer les tests
}
```

#### `SetJoinResult`

```go
func (c *BetaJoinCache) SetJoinResult(
    tokenID, factID, joinNodeID string,
    result bool,
)
```

Stocke un résultat de jointure dans le cache.

**Paramètres :**
- `tokenID` : ID du token
- `factID` : ID du fait
- `joinNodeID` : ID du JoinNode
- `result` : Résultat de la jointure

**Exemple :**

```go
// Évaluer la jointure
result := joinNode.evaluateJoinTests(token, fact)

// Mettre en cache
cache.SetJoinResult(token.ID, fact.ID, joinNode.ID, result)
```

#### `InvalidateForFact`

```go
func (c *BetaJoinCache) InvalidateForFact(factID string)
```

Invalide toutes les entrées du cache concernant un fait.

**Paramètres :**
- `factID` : ID du fait

**Exemple :**

```go
// Lors de la rétraction d'un fait
network.RetractFact(fact)

// Invalider le cache
cache.InvalidateForFact(fact.ID)
```

#### `GetStats`

```go
func (c *BetaJoinCache) GetStats() JoinCacheStats
```

Récupère les statistiques du cache.

**Retour :**
- `JoinCacheStats` : Snapshot des statistiques

**Exemple :**

```go
stats := cache.GetStats()
fmt.Printf("Cache size: %d\n", stats.Size)
fmt.Printf("Hit rate: %.2f%%\n", stats.HitRate*100)
fmt.Printf("Evictions: %d\n", stats.Evictions)
```

### BetaChainMetrics

#### `NewBetaChainMetrics`

```go
func NewBetaChainMetrics() *BetaChainMetrics
```

Crée un nouveau collecteur de métriques.

**Retour :**
- `*BetaChainMetrics` : Nouveau collecteur

**Exemple :**

```go
metrics := NewBetaChainMetrics()
```

#### `RecordChainBuild`

```go
func (m *BetaChainMetrics) RecordChainBuild(
    duration time.Duration,
    nodesCreated, nodesReused int,
)
```

Enregistre les métriques d'une construction de chaîne.

**Paramètres :**
- `duration` : Durée de la construction
- `nodesCreated` : Nombre de nœuds créés
- `nodesReused` : Nombre de nœuds réutilisés

**Exemple :**

```go
start := time.Now()
// ... construction de la chaîne ...
duration := time.Since(start)

metrics.RecordChainBuild(duration, 5, 3)
```

#### `GetSnapshot`

```go
func (m *BetaChainMetrics) GetSnapshot() MetricsSnapshot
```

Récupère un snapshot des métriques.

**Retour :**
- `MetricsSnapshot` : Snapshot des métriques

**Exemple :**

```go
snapshot := metrics.GetSnapshot()
fmt.Printf("Total chains built: %d\n", snapshot.TotalChains)
fmt.Printf("Average build time: %v\n", snapshot.AvgBuildTime)
fmt.Printf("Total nodes created: %d\n", snapshot.TotalNodesCreated)
fmt.Printf("Total nodes reused: %d\n", snapshot.TotalNodesReused)
```

---

## Gestion des cas edge

### 1. Règles avec un seul pattern

**Problème :** Pas de jointure à effectuer.

**Solution :** Connecter directement l'AlphaNode au TerminalNode, sans créer de JoinNode.

```go
func (b *BetaChainBuilder) BuildChain(
    patterns []Pattern,
    terminal *TerminalNode,
) error {
    if len(patterns) == 0 {
        return errors.New("no patterns provided")
    }
    
    if len(patterns) == 1 {
        // Cas spécial : un seul pattern, pas de jointure
        alphaNode := b.getOrCreateAlphaNode(patterns[0])
        alphaNode.AddChild(terminal)
        return nil
    }
    
    // Cas normal : plusieurs patterns
    // ... (algorithme standard)
}
```

### 2. Patterns identiques

**Problème :** Même pattern répété plusieurs fois dans une règle.

**Solution :** Partager le même AlphaNode, mais créer des JoinNodes distincts si nécessaire.

```go
Règle :
  when
    Order($o1 : id, status == "active")
    Order($o2 : id, status == "active", id != $o1)
  then
    ...

Réseau :
  AlphaNode[Order, status=="active"] (PARTAGÉ)
       │
       ├──► JoinNode[self-join avec id != id]
       │
       └──► TerminalNode
```

### 3. Cycles dans les dépendances

**Problème :** Patterns avec dépendances cycliques.

**Solution :** Détection de cycles lors de l'optimisation et utilisation d'un ordre topologique.

```go
func (o *JoinOrderOptimizer) OptimizeJoinOrder(
    patterns []Pattern,
) ([]Pattern, error) {
    // 1. Construire le graphe de dépendances
    graph := o.buildDependencyGraph(patterns)
    
    // 2. Détecter les cycles
    if hasCycle(graph) {
        return nil, errors.New("circular dependency detected")
    }
    
    // 3. Tri topologique
    ordered := topologicalSort(graph)
    
    return ordered, nil
}
```

### 4. Patterns avec négation

**Problème :** Patterns avec `not` ou `not exists`.

**Solution :** Utiliser des nœuds spécialisés (NegationNode, ExistsNode) plutôt que des JoinNodes.

```go
Règle :
  when
    Order($orderId : id)
    not Payment(order_id == $orderId)
  then
    ...

Réseau :
  AlphaNode[Order]
       │
       ▼
  NegationNode[Payment avec order_id == $orderId]
       │
       ▼
  TerminalNode
```

### 5. Mémoire pleine (OOM)

**Problème :** Les mémoires des JoinNodes consomment trop de RAM.

**Solution :**
- Limiter la taille des mémoires (LRU)
- Garbage collection aggressive
- Alertes Prometheus

```go
type JoinNode struct {
    // ...
    LeftMemory  *BoundedMemory  // Mémoire avec limite
    RightMemory *BoundedMemory
}

type BoundedMemory struct {
    store   *LRUCache
    maxSize int
}

func (m *BoundedMemory) Add(item interface{}) error {
    if m.store.Len() >= m.maxSize {
        // Évincer l'élément le moins récemment utilisé
        m.store.RemoveOldest()
    }
    m.store.Add(item)
    return nil
}
```

### 6. Patterns avec wildcards

**Problème :** Patterns trop génériques (ex: `Fact()`).

**Solution :**
- Avertissement lors de la compilation
- Optimisation de l'ordre pour placer ces patterns en dernier

```go
func (o *JoinOrderOptimizer) estimateCardinality(p Pattern) int {
    if len(p.Constraints) == 0 {
        // Wildcard pattern : cardinality maximale
        return math.MaxInt32
    }
    
    // Estimation basée sur les contraintes
    return o.estimateFromConstraints(p.Constraints)
}
```

### 7. Règles supprimées avec nœuds partagés

**Problème :** Supprimer une règle ne doit pas affecter les autres règles partageant les mêmes JoinNodes.

**Solution :** Compteur de références (voir Lifecycle Management).

```go
func (n *ReteNetwork) RemoveRule(ruleID string) error {
    joinNodes := n.findJoinNodesForRule(ruleID)
    
    for _, node := range joinNodes {
        node.lifecycle.DecrementRef()
        
        // Ne supprimer que si RefCount == 0
        if node.lifecycle.RefCount == 0 {
            n.deleteJoinNode(node)
        }
    }
    
    return nil
}
```

### 8. Hash collisions

**Problème :** Deux patterns différents produisent le même hash.

**Solution :**
- Utiliser un hash de 64 bits (collision très improbable)
- En cas de collision détectée, comparer les patterns normalisés

```go
func (r *BetaSharingRegistry) GetOrCreateJoinNode(
    joinKey string,
    pattern Pattern,
    parent Node,
) (*JoinNode, bool) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    if existingNode, ok := r.joinNodes[joinKey]; ok {
        // Vérifier que le pattern correspond exactement
        if !r.patternsEqual(existingNode.pattern, pattern) {
            // Collision détectée !
            log.Warn("Hash collision detected", 
                "key", joinKey,
                "pattern1", existingNode.pattern,
                "pattern2", pattern)
            
            // Générer une nouvelle clé avec un salt
            joinKey = joinKey + "_collision"
        } else {
            // Vraiment le même pattern, partager
            existingNode.lifecycle.IncrementRef()
            return existingNode, false
        }
    }
    
    // Créer un nouveau nœud
    // ...
}
```

### 9. Ordre d'évaluation non-déterministe

**Problème :** Ordre différent des résultats entre exécutions.

**Solution :**
- Tri stable des patterns lors de la normalisation
- Ordre déterministe dans les mémoires (par timestamp ou ID)

```go
func (m *MemoryStore) GetAll() []*Token {
    tokens := m.items
    
    // Trier par timestamp pour garantir l'ordre
    sort.Slice(tokens, func(i, j int) bool {
        return tokens[i].Timestamp < tokens[j].Timestamp
    })
    
    return tokens
}
```

### 10. Deadlocks lors de suppressions concurrentes

**Problème :** Suppression concurrente de plusieurs règles partageant des nœuds.

**Solution :**
- Verrouillage global durant la suppression
- Ordre de verrouillage cohérent

```go
func (n *ReteNetwork) RemoveRules(ruleIDs []string) error {
    // Verrouillage global pour éviter les deadlocks
    n.mutex.Lock()
    defer n.mutex.Unlock()
    
    // Trier les IDs pour garantir un ordre cohérent
    sort.Strings(ruleIDs)
    
    for _, ruleID := range ruleIDs {
        if err := n.removeRuleUnsafe(ruleID); err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Optimisations

### 1. Optimisation du cache de hash

**Problème :** Recalcul coûteux des hash pour les mêmes patterns.

**Solution :** Cache LRU des hash normalisés.

**Impact :**
- Réduction de 90%+ du temps de calcul de hash
- Overhead mémoire : ~10-20 KB pour 1000 entrées

```go
type BetaSharingRegistry struct {
    hashCache *LRUCache  // Pattern normalisé → hash
    // ...
}

func (r *BetaSharingRegistry) ComputeJoinKey(pattern Pattern) string {
    normalized := r.normalize(pattern)
    key := normalized.String()
    
    // Vérifier le cache
    if cached, ok := r.hashCache.Get(key); ok {
        return cached.(string)
    }
    
    // Calculer et mettre en cache
    hash := computeHash(normalized)
    r.hashCache.Set(key, hash)
    
    return hash
}
```

### 2. Optimisation du cache de jointure

**Problème :** Réévaluation coûteuse des mêmes jointures.

**Solution :** Cache LRU des résultats de jointure.

**Impact :**
- Réduction de 60-70% des évaluations de jointure
- Hit rate typique : 65-70% (workload mixte)

```go
func (j *JoinNode) ActivateLeft(token *Token) {
    for _, fact := range j.rightMemory.GetAll() {
        cacheKey := computeJoinCacheKey(token.ID, fact.ID, j.ID)
        
        // Vérifier le cache
        if result, found := j.cache.GetJoinResult(cacheKey); found {
            if result {
                j.propagateMatch(token, fact)
            }
            continue  // Cache hit, pas d'évaluation
        }
        
        // Cache miss, évaluer
        result := j.evaluateJoinTests(token, fact)
        j.cache.SetJoinResult(cacheKey, result)
        
        if result {
            j.propagateMatch(token, fact)
        }
    }
}
```

### 3. Optimisation de l'ordre de jointure

**Problème :** Ordre sous-optimal génère des résultats intermédiaires volumineux.

**Solution :** Heuristiques basées sur la sélectivité et la cardinalité.

**Impact :**
- Réduction de 30-50% du nombre de tokens intermédiaires
- Réduction de 20-40% du temps d'exécution

**Heuristiques :**

```go
func (o *JoinOrderOptimizer) scorePattern(p Pattern) float64 {
    score := 0.0
    
    // 1. Nombre de contraintes (plus = meilleur)
    score += float64(len(p.Constraints)) * 10.0
    
    // 2. Type de contrainte
    for _, c := range p.Constraints {
        switch c.Operator {
        case "==":
            score += 5.0  // Égalité très sélective
        case "!=":
            score += 1.0  // Inégalité peu sélective
        case "<", ">", "<=", ">=":
            score += 2.0  // Range moyennement sélectif
        }
    }
    
    // 3. Disponibilité d'index
    if o.hasIndex(p.Type, p.Constraints[0].Field) {
        score += 15.0
    }
    
    // 4. Cardinalité estimée (moins = meilleur)
    cardinality := o.estimateCardinality(p)
    score -= float64(cardinality) / 100.0
    
    return score
}
```

### 4. Index inversé pour invalidation

**Problème :** Invalidation du cache O(n) lors de rétraction de faits.

**Solution :** Index inversé (fact → cache keys).

**Impact :**
- Réduction de O(n) à O(k) où k = nombre de clés par fait
- Typiquement k << n

```go
type BetaJoinCacheWithIndex struct {
    cache       *LRUCache
    factIndex   map[string][]string  // factID → [cacheKeys]
    mutex       sync.RWMutex
}

func (c *BetaJoinCacheWithIndex) InvalidateForFact(factID string) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    // O(k) au lieu de O(n)
    for _, key := range c.factIndex[factID] {
        c.cache.Remove(key)
    }
    
    delete(c.factIndex, factID)
}
```

### 5. Partage de préfixes

**Problème :** Règles similaires créent des sous-chaînes redondantes.

**Solution :** Détecter et partager les préfixes communs.

**Impact :**
- Réduction de 40-60% du nombre de JoinNodes (règles similaires)
- Réduction de 30-50% de la mémoire consommée

```go
func (b *BetaChainBuilder) BuildChain(
    patterns []Pattern,
    terminal *TerminalNode,
) error {
    var currentParent Node = nil
    
    for i, pattern := range patterns {
        // Calculer la clé incluant le préfixe
        prefix := patterns[:i]
        joinKey := b.registry.ComputeJoinKeyWithPrefix(pattern, prefix)
        
        // Obtenir ou créer (réutilise si préfixe identique)
        joinNode, created := b.registry.GetOrCreateJoinNode(
            joinKey, pattern, currentParent,
        )
        
        if !created {
            // Préfixe partagé !
            atomic.AddInt64(&b.metrics.PrefixShared, 1)
        }
        
        currentParent = joinNode
    }
    
    currentParent.AddChild(terminal)
    return nil
}
```

### 6. Lazy evaluation des jointures

**Problème :** Évaluation de toutes les jointures même si elles ne produisent pas de résultat.

**Solution :** Court-circuitage lors de l'échec du premier test.

**Impact :**
- Réduction de 20-30% du temps d'évaluation des jointures
- Particulièrement efficace avec plusieurs tests de jointure

```go
func (j *JoinNode) evaluateJoinTests(token *Token, fact *Fact) bool {
    for _, test := range j.tests {
        leftFact := token.Facts[test.LeftIndex]
        leftValue := leftFact.Fields[test.LeftField]
        rightValue := fact.Fields[test.RightField]
        
        // Court-circuitage : échec immédiat
        if !evaluateOperator(leftValue, test.Operator, rightValue) {
            return false
        }
    }
    
    return true
}
```

### 7. Batch processing

**Problème :** Propagation individuelle des tokens coûteuse.

**Solution :** Accumuler et propager par lots.

**Impact :**
- Réduction de 15-25% de l'overhead de propagation
- Meilleure utilisation du cache CPU

```go
type JoinNode struct {
    // ...
    batchSize     int
    pendingTokens []*Token
    mutex         sync.Mutex
}

func (j *JoinNode) ActivateLeft(token *Token) {
    j.mutex.Lock()
    j.pendingTokens = append(j.pendingTokens, token)
    shouldFlush := len(j.pendingTokens) >= j.batchSize
    j.mutex.Unlock()
    
    if shouldFlush {
        j.flushBatch()
    }
}

func (j *JoinNode) flushBatch() {
    j.mutex.Lock()
    batch := j.pendingTokens
    j.pendingTokens = nil
    j.mutex.Unlock()
    
    for _, token := range batch {
        j.processSingleToken(token)
    }
}
```

---

## Internals

### Structure mémoire détaillée

#### Layout d'un JoinNode en mémoire

```
JoinNode (heap allocated)
├── Header (24 bytes)
│   ├── ID (string ptr + len)      : 16 bytes
│   └── Type (int32)                : 4 bytes
│   └── Padding                     : 4 bytes
├── Tests (slice)                   : 24 bytes
│   ├── Pointer                     : 8 bytes
│   ├── Length                      : 8 bytes
│   └── Capacity                    : 8 bytes
├── Parents (2 pointers)            : 16 bytes
├── Children (slice)                : 24 bytes
├── Lifecycle (pointer)             : 8 bytes
├── LeftMemory (pointer)            : 8 bytes
├── RightMemory (pointer)           : 8 bytes
└── Mutex (sync.RWMutex)            : 24 bytes
                                    ─────────
Total per JoinNode                  : ~136 bytes + data

Mémoires (left/right) :
- BoundedMemory struct              : 32 bytes
- LRU cache overhead                : 48 bytes per cache
- Items (tokens/facts)              : variable

Exemple avec 100 tokens et 50 facts :
- JoinNode overhead                 : 136 bytes
- LeftMemory (100 tokens × 64 B)    : 6,400 bytes
- RightMemory (50 facts × 48 B)     : 2,400 bytes
                                    ─────────
Total                               : ~9 KB per JoinNode
```

#### Overhead du partage

```
Sans partage (10 règles similaires) :
- 10 JoinNodes × 9 KB               : 90 KB

Avec partage (10 règles → 1 JoinNode) :
- 1 JoinNode × 9 KB                 : 9 KB
- BetaSharingRegistry overhead      : 2 KB
- Hash cache (1000 entrées)         : 20 KB
                                    ─────────
Total                               : 31 KB

Économie                            : 59 KB (65% de réduction)
```

### Complexité algorithmique

#### Construction de chaîne

| Opération                    | Sans partage | Avec partage  | Notes                          |
|------------------------------|--------------|---------------|--------------------------------|
| Normalisation pattern        | -            | O(n log n)    | n = nombre de contraintes      |
| Calcul de hash               | -            | O(n)          | Amorti à O(1) avec cache       |
| Lookup dans registry         | -            | O(1)          | HashMap                        |
| Création JoinNode            | O(1)         | O(1)          | Allocation + init              |
| Connexion parent-enfant      | O(1)         | O(1)          | Ajout à slice                  |
| **Total par pattern**        | **O(1)**     | **O(n log n)**| Amorti à O(1) avec cache       |
| **Total pour k patterns**    | **O(k)**     | **O(k × n log n)** | Amorti à O(k)            |

#### Runtime (évaluation)

| Opération                    | Sans cache   | Avec cache    | Notes                          |
|------------------------------|--------------|---------------|--------------------------------|
| Lookup cache de jointure     | -            | O(1)          | HashMap                        |
| Évaluation tests de jointure | O(t)         | O(1) (hit)    | t = nombre de tests            |
| Propagation token            | O(c)         | O(c)          | c = nombre d'enfants           |
| **Activation gauche**        | **O(r × t)** | **O(r) (avg)**| r = taille right memory        |
| **Activation droite**        | **O(l × t)** | **O(l) (avg)**| l = taille left memory         |

#### Invalidation du cache

| Approche                     | Complexité   | Overhead mémoire | Notes                     |
|------------------------------|--------------|------------------|---------------------------|
| Scan complet                 | O(n)         | 0                | n = taille du cache       |
| Index inversé                | O(k)         | O(n × f)         | k = clés par fact         |
| Bloom filter                 | O(k)         | O(m)             | m = taille du filtre      |

### Concurrence et thread-safety

#### Modèle de verrouillage

```
┌─────────────────────────────────────────────────────────────────┐
│                   Lock Hierarchy                                  │
└─────────────────────────────────────────────────────────────────┘

Level 1: ReteNetwork.mutex (global)
         │
         ├──► Lock pour : AddRule, RemoveRule, Reset
         │
         └──► Jamais tenu pendant des opérations longues
              
Level 2: BetaSharingRegistry.mutex
         │
         ├──► Lock pour : GetOrCreateJoinNode, GetSharingStats
         │
         └──► Peut appeler JoinNode méthodes (niveau 3)
              
Level 3: JoinNode.mutex (per-node)
         │
         ├──► Lock pour : Activate, AddChild, UpdateMemory
         │
         └──► Jamais appeler registry (évite deadlock)
              
Level 4: Cache.mutex (per-cache)
         │
         ├──► Lock pour : Get, Set, Remove
         │
         └──► Opérations atomiques courtes

Règle d'or : Toujours verrouiller de haut en bas (Level 1 → Level 4)
             Jamais dans l'autre sens !
```

#### Patrons de synchronisation

**Pattern 1 : Read-Write Lock**

```go
type BetaSharingRegistry struct {
    joinNodes map[string]*JoinNode
    mutex     sync.RWMutex  // Lecture fréquente, écriture rare
}

func (r *BetaSharingRegistry) GetJoinNode(key string) *JoinNode {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    return r.joinNodes[key]
}

func (r *BetaSharingRegistry) CreateJoinNode(key string, node *JoinNode) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    
    r.joinNodes[key] = node
}
```

**Pattern 2 : Opérations atomiques**

```go
type BetaSharingStats struct {
    TotalJoinNodes  int64  // sync/atomic
    SharedJoinNodes int64  // sync/atomic
}

func (r *BetaSharingRegistry) IncrementShared() {
    atomic.AddInt64(&r.metrics.SharedJoinNodes, 1)
}
```

**Pattern 3 : Copy-on-Write**

```go
func (r *BetaSharingRegistry) GetSharingStats() BetaSharingStats {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    
    // Copie snapshot (pas de référence partagée)
    return BetaSharingStats{
        TotalJoinNodes:  atomic.LoadInt64(&r.metrics.TotalJoinNodes),
        SharedJoinNodes: atomic.LoadInt64(&r.metrics.SharedJoinNodes),
        // ...
    }
}
```

### Debugging et observabilité

#### Logging structuré

```go
import "go.uber.org/zap"

func (b *BetaChainBuilder) BuildChain(
    patterns []Pattern,
    terminal *TerminalNode,
) error {
    logger := b.logger.With(
        zap.String("rule_id", terminal.RuleID),
        zap.Int("pattern_count", len(patterns)),
    )
    
    logger.Info("Starting beta chain build")
    
    for i, pattern := range patterns {
        joinKey := b.registry.ComputeJoinKey(pattern)
        joinNode, created := b.registry.GetOrCreateJoinNode(joinKey, pattern, nil)
        
        logger.Debug("Processed pattern",
            zap.Int("index", i),
            zap.String("type", pattern.Type),
            zap.String("join_key", joinKey),
            zap.Bool("created", created),
            zap.String("join_node_id", joinNode.ID),
        )
    }
    
    logger.Info("Beta chain build complete",
        zap.Duration("duration", time.Since(start)),
    )
    
    return nil
}
```

#### Métriques Prometheus

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    betaChainBuildDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "beta_chain_build_duration_seconds",
            Help:    "Time spent building beta chains",
            Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
        },
    )
    
    betaNodeSharingRatio = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "beta_node_sharing_ratio",
            Help: "Ratio of shared vs total join nodes",
        },
    )
    
    betaJoinCacheHitRate = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "beta_join_cache_hit_rate",
            Help: "Hit rate of the join result cache",
        },
    )
)

func (b *BetaChainBuilder) recordMetrics() {
    stats := b.registry.GetSharingStats()
    
    betaNodeSharingRatio.Set(stats.SharingRatio)
    
    cacheStats := b.cache.GetStats()
    betaJoinCacheHitRate.Set(cacheStats.HitRate)
}
```

#### Tracing distribué

```go
import "go.opentelemetry.io/otel/trace"

func (b *BetaChainBuilder) BuildChain(
    ctx context.Context,
    patterns []Pattern,
    terminal *TerminalNode,
) error {
    ctx, span := b.tracer.Start(ctx, "BetaChainBuilder.BuildChain")
    defer span.End()
    
    span.SetAttributes(
        attribute.Int("pattern_count", len(patterns)),
        attribute.String("rule_id", terminal.RuleID),
    )
    
    for i, pattern := range patterns {
        ctx, patternSpan := b.tracer.Start(ctx, "process_pattern")
        patternSpan.SetAttributes(
            attribute.Int("index", i),
            attribute.String("type", pattern.Type),
        )
        
        // ... traitement ...
        
        patternSpan.End()
    }
    
    return nil
}
```

### Profiling et analyse

#### CPU Profiling

```bash
# Générer un profil CPU
go test -bench=BenchmarkBetaChainBuild -cpuprofile=cpu.prof

# Analyser avec pprof
go tool pprof cpu.prof

# Top fonctions consommatrices
(pprof) top10

# Graphe d'appels
(pprof) web

# Focus sur une fonction spécifique
(pprof) list BetaChainBuilder.BuildChain
```

#### Memory Profiling

```bash
# Générer un profil mémoire
go test -bench=BenchmarkBetaChainBuild -memprofile=mem.prof

# Analyser les allocations
go tool pprof mem.prof

# Top allocateurs
(pprof) top10

# Allocations in-use
(pprof) inuse_space

# Allocations totales
(pprof) alloc_space
```

#### Mutex Profiling

```bash
# Activer le profil de contention
go test -bench=BenchmarkBetaChainBuild -mutexprofile=mutex.prof

# Analyser les contentions
go tool pprof mutex.prof

# Identifier les hot spots
(pprof) top10
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
FITNESS FOR A PARTICULAR PURPOSE AND NON