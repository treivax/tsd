# Normalisation des Conditions Alpha

## Vue d'ensemble

La normalisation des conditions permet d'ordonner les conditions de manière **canonique et déterministe** pour améliorer le partage des nœuds Alpha dans le réseau RETE. Deux expressions logiquement équivalentes mais écrites dans un ordre différent seront normalisées au même ordre canonique.

## Motivation

### Problème

Dans un réseau RETE, les expressions suivantes sont logiquement équivalentes mais créaient des nœuds Alpha différents :

```go
// Règle 1
p.age > 18 AND p.salary >= 50000

// Règle 2
p.salary >= 50000 AND p.age > 18
```

Sans normalisation, ces deux règles créent des chaînes Alpha différentes, ce qui réduit l'efficacité du partage de nœuds.

### Solution

La normalisation réordonne les conditions selon un ordre canonique **quand l'opérateur est commutatif** (AND, OR, +, *, ==, !=). Ainsi :

- `A AND B` et `B AND A` → même ordre canonique
- `A OR B` et `B OR A` → même ordre canonique
- `A - B` et `B - A` → ordre préservé (non-commutatif)

## API

### IsCommutative

```go
func IsCommutative(operator string) bool
```

Détermine si un opérateur est commutatif.

**Opérateurs commutatifs** (peuvent être réordonnés) :
- `AND`, `OR`, `&&`, `||`
- `+`, `*`
- `==`, `!=`, `<>`

**Opérateurs non-commutatifs** (ordre préservé) :
- `-`, `/`
- `<`, `>`, `<=`, `>=`
- Opérateurs séquentiels

**Exemple :**

```go
IsCommutative("AND")  // true
IsCommutative("OR")   // true
IsCommutative("-")    // false
IsCommutative("<")    // false
```

### NormalizeConditions

```go
func NormalizeConditions(conditions []SimpleCondition, operator string) []SimpleCondition
```

Trie les conditions dans un ordre canonique déterministe.

**Comportement :**
- Si l'opérateur est **commutatif** → trie les conditions par ordre lexicographique de leur représentation canonique
- Si l'opérateur est **non-commutatif** → préserve l'ordre original
- Si 0 ou 1 condition → retourne tel quel

**Exemple :**

```go
// Créer deux conditions
condAge := NewSimpleCondition(
    "binaryOperation",
    FieldAccess{Object: "p", Field: "age"},
    ">",
    NumberLiteral{Value: 18},
)

condSalary := NewSimpleCondition(
    "binaryOperation",
    FieldAccess{Object: "p", Field: "salary"},
    ">=",
    NumberLiteral{Value: 50000},
)

// Les deux ordres produisent le même résultat
normalized1 := NormalizeConditions([]SimpleCondition{condAge, condSalary}, "AND")
normalized2 := NormalizeConditions([]SimpleCondition{condSalary, condAge}, "AND")

// normalized1 == normalized2 ✅
```

### NormalizeExpression

```go
func NormalizeExpression(expr interface{}) (interface{}, error)
```

Point d'entrée principal pour normaliser une expression complète.

**Types supportés :**
- `constraint.LogicalExpression`
- `constraint.BinaryOperation`
- `constraint.Constraint`
- `map[string]interface{}`
- Littéraux et accès de champs (retournés inchangés)

**Exemple :**

```go
// Expression originale : salary >= 50000 AND age > 18
expr := constraint.LogicalExpression{
    Type: "logicalExpr",
    Left: BinaryOperation{
        Left: FieldAccess{Object: "p", Field: "salary"},
        Operator: ">=",
        Right: NumberLiteral{Value: 50000},
    },
    Operations: []LogicalOperation{
        {
            Op: "AND",
            Right: BinaryOperation{
                Left: FieldAccess{Object: "p", Field: "age"},
                Operator: ">",
                Right: NumberLiteral{Value: 18},
            },
        },
    },
}

// Normaliser avec reconstruction automatique
normalized, err := NormalizeExpression(expr)
if err != nil {
    log.Fatal(err)
}

// normalized contient maintenant : age > 18 AND salary >= 50000
// L'arbre d'expression a été complètement reconstruit en ordre canonique
```

## Algorithme

### 1. Extraction des Conditions

```
ExtractConditions(expression) → (conditions[], operatorType, error)
```

Extrait toutes les conditions atomiques d'une expression complexe.

### 2. Détermination de la Commutativité

```
IsCommutative(operator) → bool
```

Vérifie si l'opérateur permet le réordonnancement.

### 3. Tri Canonique

```
Si commutatif:
    Trier conditions par CanonicalString(condition)
Sinon:
    Préserver l'ordre original
```

### 4. Représentation Canonique

Chaque condition est convertie en une chaîne unique :

```
binaryOperation(fieldAccess(p,age),>,literal(18))
```

Le tri lexicographique de ces chaînes produit un ordre déterministe.

## Cas d'Usage

### 1. Partage de Nœuds Alpha

```go
// Règle 1: age > 18 AND salary >= 50000
// Règle 2: salary >= 50000 AND age > 18

conditions1, op1, _ := ExtractConditions(rule1.Constraint)
normalized1 := NormalizeConditions(conditions1, op1)

conditions2, op2, _ := ExtractConditions(rule2.Constraint)
normalized2 := NormalizeConditions(conditions2, op2)

// normalized1 == normalized2
// → Même chaîne Alpha → Partage de nœuds ✅
```

### 2. Déduplication de Règles

```go
// Détecter les règles sémantiquement identiques
rules := []Rule{rule1, rule2, rule3}

seen := make(map[string]bool)
for _, rule := range rules {
    conditions, op, _ := ExtractConditions(rule.Constraint)
    normalized := NormalizeConditions(conditions, op)
    
    key := computeKey(normalized)
    if seen[key] {
        log.Printf("Règle dupliquée détectée: %s", rule.Name)
    }
    seen[key] = true
}
```

### 3. Optimisation de Requêtes

```go
// Normaliser avant la création du réseau RETE
for _, rule := range rules {
    conditions, op, _ := ExtractConditions(rule.Constraint)
    rule.NormalizedConditions = NormalizeConditions(conditions, op)
}

// Construire le réseau avec les conditions normalisées
network := BuildReteNetwork(rules)
```

## Exemples Complets

### Exemple 1 : AND Normalization

```go
// Définir deux règles avec AND dans des ordres différents
rule1 := "p.age > 18 AND p.salary >= 50000"
rule2 := "p.salary >= 50000 AND p.age > 18"

// Extraire et normaliser
conds1, op1, _ := ExtractConditions(parseRule(rule1))
normalized1 := NormalizeConditions(conds1, op1)

conds2, op2, _ := ExtractConditions(parseRule(rule2))
normalized2 := NormalizeConditions(conds2, op2)

// Vérifier l'égalité
for i := range normalized1 {
    if !CompareConditions(normalized1[i], normalized2[i]) {
        panic("Normalisation échouée!")
    }
}
// ✅ Succès : même ordre canonique
```

### Exemple 2 : OR Normalization

```go
// OR est également commutatif
rule1 := "status == 'active' OR verified == true"
rule2 := "verified == true OR status == 'active'"

conds1, _ := ExtractConditions(parseRule(rule1))
conds2, _ := ExtractConditions(parseRule(rule2))

normalized1 := NormalizeConditions(conds1, "OR")
normalized2 := NormalizeConditions(conds2, "OR")

// normalized1 == normalized2 ✅
```

### Exemple 3 : Non-Commutative Preservation

```go
// Les opérations non-commutatives préservent l'ordre
expr1 := "x - 10"
expr2 := "10 - x"

conds1, _ := ExtractConditions(parseExpr(expr1))
conds2, _ := ExtractConditions(parseExpr(expr2))

// L'opérateur '-' n'est PAS commutatif
normalized1 := NormalizeConditions(conds1, "-")
normalized2 := NormalizeConditions(conds2, "-")

// normalized1 != normalized2 ✅ (ordre préservé)
```

## Exécuter la Démonstration

Un exemple complet est fourni dans `examples/normalization/main.go` :

```bash
go run ./rete/examples/normalization/main.go
```

**Output attendu (extrait) :**

```
📋 Exemple 1: Normalisation AND (opérateur commutatif)
=============================================================

🔄 Ordre A: age > 18 AND salary >= 50000
   [0] binaryOperation(fieldAccess(p,age),>,literal(18))
   [1] binaryOperation(fieldAccess(p,salary),>=,literal(50000))

🔄 Ordre B: salary >= 50000 AND age > 18
   [0] binaryOperation(fieldAccess(p,age),>,literal(18))
   [1] binaryOperation(fieldAccess(p,salary),>=,literal(50000))

✅ Les deux ordres produisent le MÊME ordre canonique!

...

📋 Exemple 5: Reconstruction d'Expressions Normalisées
=============================================================

🔍 Expression originale (ordre inversé):
   (p.salary >= 50000) AND (p.age > 18)

✨ Normalisation avec RECONSTRUCTION automatique...

📊 Conditions APRÈS normalisation et reconstruction:
   [0] binaryOperation(fieldAccess(p,age),>,literal(18))
   [1] binaryOperation(fieldAccess(p,salary),>=,literal(50000))

🔍 Vérification de l'ordre canonique:
   ✓ Premier élément (Left): p.age > ...
     ✅ Correct ! 'age' vient avant 'salary' en ordre canonique

✅ Résultat:
   🎉 Les deux expressions ont été reconstruites avec le MÊME ordre canonique!
   → Le partage de nœuds Alpha sera maximal
```

## Tests

Tous les tests sont dans `alpha_chain_extractor_normalize_test.go` :

```bash
# Exécuter tous les tests de normalisation
go test -v ./rete -run "TestNormalize|TestIsCommutative"

# Tests spécifiques
go test -v ./rete -run TestNormalizeConditions_AND_OrderIndependent
go test -v ./rete -run TestNormalizeConditions_OR_OrderIndependent
go test -v ./rete -run TestNormalizeConditions_NonCommutative_PreserveOrder
go test -v ./rete -run TestIsCommutative_AllOperators
```

### Couverture des Tests

**Tests de normalisation :**
- ✅ `TestIsCommutative_AllOperators` - Tous les opérateurs (commutatifs et non-commutatifs)
- ✅ `TestNormalizeConditions_AND_OrderIndependent` - AND : A∧B == B∧A
- ✅ `TestNormalizeConditions_OR_OrderIndependent` - OR : A∨B == B∨A
- ✅ `TestNormalizeConditions_NonCommutative_PreserveOrder` - Préservation de l'ordre
- ✅ `TestNormalizeConditions_EmptyAndSingle` - Cas limites (0 et 1 condition)
- ✅ `TestNormalizeConditions_ThreeConditions` - 3+ conditions, toutes permutations
- ✅ `TestNormalizeExpression_ComplexNested` - Expressions imbriquées
- ✅ `TestNormalizeExpression_BinaryOperation` - Opérations binaires simples
- ✅ `TestNormalizeExpression_Map` - Expressions sous forme de map
- ✅ `TestNormalizeExpression_Literals` - Littéraux inchangés
- ✅ `TestNormalizeConditions_DeterministicOrder` - Déterminisme du tri

**Tests de reconstruction :**
- ✅ `TestRebuildLogicalExpression_SingleCondition` - Reconstruction avec 1 condition
- ✅ `TestRebuildLogicalExpression_TwoConditions` - Reconstruction avec 2 conditions
- ✅ `TestRebuildLogicalExpression_ThreeConditions` - Reconstruction avec 3+ conditions
- ✅ `TestRebuildLogicalExpression_Empty` - Cas d'erreur (liste vide)
- ✅ `TestNormalizeExpression_WithReconstruction` - Normalisation complète avec reconstruction
- ✅ `TestNormalizeExpression_PreservesSemantics` - Préservation de la sémantique
- ✅ `TestRebuildLogicalExpressionMap_TwoConditions` - Reconstruction de map
- ✅ `TestNormalizeExpressionMap_WithReconstruction` - Normalisation map avec reconstruction

## Propriétés Garanties

### 1. Idempotence

```go
normalized := NormalizeConditions(conditions, "AND")
normalized2 := NormalizeConditions(normalized, "AND")
// normalized == normalized2 ✅
```

### 2. Déterminisme

```go
// Normaliser 100 fois produit toujours le même résultat
for i := 0; i < 100; i++ {
    result := NormalizeConditions(conditions, "AND")
    assert(result == expected)
}
```

### 3. Préservation Sémantique

```go
// La normalisation ne change PAS la sémantique
original := evaluateConditions(conditions, "AND", fact)
normalized := NormalizeConditions(conditions, "AND")
result := evaluateConditions(normalized, "AND", fact)
// original == result ✅
```

### 4. Commutativité Respectée

```go
// AND et OR : A op B == B op A
if IsCommutative(op) {
    norm1 := NormalizeConditions([A, B], op)
    norm2 := NormalizeConditions([B, A], op)
    assert(norm1 == norm2)
}
```

### 5. Non-Commutativité Respectée

```go
// Soustraction : A - B != B - A
if !IsCommutative(op) {
    norm1 := NormalizeConditions([A, B], op)
    norm2 := NormalizeConditions([B, A], op)
    assert(norm1 == [A, B])  // Ordre préservé
    assert(norm2 == [B, A])  // Ordre préservé
}
```

## Intégration avec le Partage Alpha

La normalisation avec reconstruction est conçue pour maximiser le partage de nœuds Alpha :

```go
// 1. Normaliser l'expression (avec reconstruction automatique)
normalizedExpr, err := NormalizeExpression(expr)
if err != nil {
    log.Fatal(err)
}

// 2. Extraire les conditions de l'expression normalisée
conditions, opType, err := ExtractConditions(normalizedExpr)
if err != nil {
    log.Fatal(err)
}

// 3. Générer les hashes pour le partage
// Les conditions sont maintenant en ordre canonique
for _, cond := range conditions {
    hash := cond.Hash  // Hash unique déjà calculé
    // Rechercher ou créer le nœud Alpha avec ce hash
    alphaNode := getOrCreateAlphaNode(hash, cond)
}

// Résultat : Deux règles avec le même ordre canonique partageront
// exactement les mêmes nœuds Alpha, maximisant l'efficacité
```

## Limitations et Considérations

### 1. ~~Reconstruction d'Expression~~ ✅ **IMPLÉMENTÉ**

La reconstruction complète d'expressions normalisées est maintenant implémentée :

```go
// Reconstruit une LogicalExpression à partir de conditions normalisées
func rebuildLogicalExpression(conditions []SimpleCondition, operator string) (constraint.LogicalExpression, error)

// Reconstruit une expression map à partir de conditions normalisées
func rebuildLogicalExpressionMap(conditions []SimpleCondition, operator string) (map[string]interface{}, error)

// Convertit une SimpleCondition en BinaryOperation
func rebuildConditionAsExpression(cond SimpleCondition) interface{}

// Convertit une SimpleCondition en map
func rebuildConditionAsMap(cond SimpleCondition) map[string]interface{}
```

**Fonctionnement** :
1. Les conditions normalisées sont extraites et triées
2. La première condition devient `Left` de la LogicalExpression
3. Les conditions suivantes deviennent des `Operations`
4. L'arbre d'expression est reconstruit avec la structure correcte

**Exemple** :
```go
// Expression originale : salary >= 50000 AND age > 18
expr := constraint.LogicalExpression{...}

// Normalisation avec reconstruction automatique
normalized, _ := NormalizeExpression(expr)

// Résultat : age > 18 AND salary >= 50000 (ordre canonique)
// Structure complètement reconstruite
```

### 2. Opérateurs Mixtes

Si une expression contient plusieurs types d'opérateurs (ex: `A AND B OR C`), l'opérateur est marqué comme "MIXED" et l'ordre n'est pas modifié.

### 3. Précédence des Opérateurs

La normalisation ne change pas la structure de l'arbre d'expression, seulement l'ordre des conditions au même niveau de précédence.

## Performance

- **Complexité temporelle** : O(n log n) pour n conditions (tri)
- **Complexité spatiale** : O(n) (copie des conditions)
- **Impact** : Négligeable pour des règles typiques (< 10 conditions)

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

## Voir Aussi

- `ALPHA_CHAIN_EXTRACTOR_README.md` - Extraction de conditions
- `ALPHA_NODE_SHARING.md` - Partage de nœuds Alpha
- `alpha_chain_extractor.go` - Implémentation
- `alpha_chain_extractor_normalize_test.go` - Tests complets