# Alpha Chain Extractor - Documentation

## Vue d'ensemble

Le module `alpha_chain_extractor.go` fournit des fonctionnalités pour extraire et analyser les conditions atomiques à partir d'expressions complexes dans le réseau RETE. Il permet de décomposer des expressions logiques imbriquées en conditions simples, de générer des représentations canoniques uniques, et de gérer la déduplication de conditions.

## Fonctionnalités principales

### 1. Extraction de conditions (`ExtractConditions`)

Extrait toutes les conditions simples d'une expression complexe, qu'elle soit sous forme de structures Go typées ou de maps JSON.

**Signature:**
```go
func ExtractConditions(expr interface{}) ([]SimpleCondition, string, error)
```

**Paramètres:**
- `expr`: Expression à analyser (BinaryOperation, LogicalExpression, Constraint, ou map[string]interface{})

**Retours:**
- `[]SimpleCondition`: Liste des conditions atomiques extraites
- `string`: Type d'opérateur principal ("AND", "OR", "MIXED", "SINGLE", "NONE")
- `error`: Erreur si l'expression n'est pas supportée

**Exemples:**

```go
// Comparaison simple: p.age > 18
expr := constraint.BinaryOperation{
    Type:     "binaryOperation",
    Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
    Operator: ">",
    Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
}

conditions, opType, _ := ExtractConditions(expr)
// conditions: [SimpleCondition{Type: "binaryOperation", Operator: ">", ...}]
// opType: "SINGLE"
```

```go
// Expression AND: p.age > 18 AND p.salary >= 50000
expr := constraint.LogicalExpression{
    Type: "logicalExpr",
    Left: /* p.age > 18 */,
    Operations: []constraint.LogicalOperation{
        {Op: "AND", Right: /* p.salary >= 50000 */},
    },
}

conditions, opType, _ := ExtractConditions(expr)
// conditions: [condition1, condition2]
// opType: "AND"
```

```go
// Expressions imbriquées: (p.age > 18 AND p.salary >= 50000) OR p.vip == true
// opType: "MIXED"
```

### 2. Structure SimpleCondition

Représente une condition atomique avec hachage automatique pour comparaison et déduplication.

**Définition:**
```go
type SimpleCondition struct {
    Type     string      // Type: "binaryOperation", "comparison", "constraint", etc.
    Left     interface{} // Opérande gauche
    Operator string      // Opérateur: ">", "==", ">=", etc.
    Right    interface{} // Opérande droite
    Hash     string      // Hash SHA-256 calculé automatiquement
}
```

**Création:**
```go
cond := NewSimpleCondition(
    "binaryOperation",
    constraint.FieldAccess{Object: "p", Field: "age"},
    ">",
    constraint.NumberLiteral{Value: 18},
)
// Le hash est calculé automatiquement
fmt.Println(cond.Hash) // e.g., "a3b5c7d9..."
```

### 3. Représentation canonique (`CanonicalString`)

Génère une représentation textuelle unique et déterministe d'une condition.

**Signature:**
```go
func CanonicalString(condition SimpleCondition) string
```

**Format:** `type(left,operator,right)`

**Exemples:**

| Condition | Représentation canonique |
|-----------|-------------------------|
| `p.age > 18` | `binaryOperation(fieldAccess(p,age),>,literal(18))` |
| `p.salary + 100` | `binaryOperation(fieldAccess(p,salary),+,literal(100))` |
| `p.name == "Alice"` | `binaryOperation(fieldAccess(p,name),==,literal(Alice))` |

**Propriétés:**
- **Déterministe**: Même condition → même string
- **Unique**: Conditions différentes → strings différentes
- **Ordonné**: Les maps sont triées par clés pour garantir la cohérence

```go
cond1 := NewSimpleCondition("binaryOperation", fieldAccess, ">", literal18)
cond2 := NewSimpleCondition("binaryOperation", fieldAccess, ">", literal18)

str1 := CanonicalString(cond1)
str2 := CanonicalString(cond2)
// str1 == str2 (déterministe)

fmt.Println(cond1.Hash == cond2.Hash) // true
```

### 4. Comparaison de conditions (`CompareConditions`)

Compare deux conditions pour l'égalité basée sur leur hash.

**Signature:**
```go
func CompareConditions(c1, c2 SimpleCondition) bool
```

**Exemple:**
```go
if CompareConditions(cond1, cond2) {
    fmt.Println("Les conditions sont identiques")
}
```

### 5. Déduplication (`DeduplicateConditions`)

Supprime les conditions dupliquées d'une liste.

**Signature:**
```go
func DeduplicateConditions(conditions []SimpleCondition) []SimpleCondition
```

**Exemple:**
```go
conditions := []SimpleCondition{cond1, cond2, cond1, cond3}
unique := DeduplicateConditions(conditions)
// unique contient seulement cond1, cond2, cond3 (sans duplication)
```

## Types d'expressions supportés

### Structures Go typées
- `constraint.BinaryOperation`: Opérations binaires et comparaisons
- `constraint.LogicalExpression`: Expressions logiques (AND/OR)
- `constraint.Constraint`: Contraintes simples
- `constraint.FieldAccess`: Accès aux champs d'objets
- `constraint.NumberLiteral`: Littéraux numériques
- `constraint.StringLiteral`: Littéraux de chaînes
- `constraint.BooleanLiteral`: Littéraux booléens

### Format map (JSON)
```go
map[string]interface{}{
    "type": "binaryOperation",
    "operator": ">",
    "left": map[string]interface{}{
        "type": "fieldAccess",
        "object": "p",
        "field": "age",
    },
    "right": map[string]interface{}{
        "type": "numberLiteral",
        "value": 18,
    },
}
```

## Types d'opérateurs retournés

| Type | Description | Exemple |
|------|-------------|---------|
| `SINGLE` | Une seule condition | `p.age > 18` |
| `AND` | Toutes les opérations sont AND | `p.age > 18 AND p.salary >= 50000` |
| `OR` | Toutes les opérations sont OR | `p.age < 25 OR p.age > 65` |
| `MIXED` | Mélange de AND et OR | `p.age > 18 AND p.salary >= 50000 OR p.vip == true` |
| `NONE` | Pas de condition (littéraux, accès champs seuls) | `p.age`, `18` |

## Cas d'usage

### 1. Analyse de règles complexes
```go
// Analyser une règle avec multiples conditions
rule := parseRule("p.age > 18 AND p.salary >= 50000 AND p.active == true")
conditions, opType, _ := ExtractConditions(rule.Constraints)

fmt.Printf("Type: %s, Nombre de conditions: %d\n", opType, len(conditions))
// Output: Type: AND, Nombre de conditions: 3
```

### 2. Construction de chaînes alpha optimisées
```go
// Extraire les conditions pour créer des nœuds alpha
conditions, _, _ := ExtractConditions(expr)
for _, cond := range conditions {
    alphaNode := createAlphaNode(cond)
    network.AddAlphaNode(alphaNode)
}
```

### 3. Détection de conditions dupliquées
```go
// Analyser plusieurs règles et trouver les conditions communes
allConditions := []SimpleCondition{}
for _, rule := range rules {
    conditions, _, _ := ExtractConditions(rule.Constraints)
    allConditions = append(allConditions, conditions...)
}

unique := DeduplicateConditions(allConditions)
duplicates := len(allConditions) - len(unique)
fmt.Printf("%d conditions dupliquées trouvées\n", duplicates)
```

### 4. Cache de conditions
```go
// Utiliser le hash comme clé de cache
cache := make(map[string]*AlphaNode)

conditions, _, _ := ExtractConditions(expr)
for _, cond := range conditions {
    if node, exists := cache[cond.Hash]; exists {
        // Réutiliser le nœud existant
        reuseNode(node)
    } else {
        // Créer un nouveau nœud
        node := createAlphaNode(cond)
        cache[cond.Hash] = node
    }
}
```

## Tests

Le module est couvert par une suite complète de tests unitaires :

- `TestExtractConditions_SimpleComparison`: Comparaisons simples
- `TestExtractConditions_LogicalAND`: Expressions AND
- `TestExtractConditions_LogicalOR`: Expressions OR
- `TestExtractConditions_NestedExpressions`: Expressions imbriquées
- `TestExtractConditions_MixedOperators`: Opérateurs mélangés (AND/OR)
- `TestExtractConditions_ArithmeticOperations`: Opérations arithmétiques
- `TestCanonicalString_Deterministic`: Vérification du déterminisme
- `TestCanonicalString_Uniqueness`: Vérification de l'unicité
- `TestCompareConditions`: Test de comparaison
- `TestDeduplicateConditions`: Test de déduplication

**Exécuter les tests:**
```bash
cd tsd
go test ./rete -run "ExtractConditions|CanonicalString|CompareConditions|DeduplicateConditions" -v
```

## Limitations et notes

1. **Expressions non supportées**: Certaines expressions complexes comme les fonctions personnalisées peuvent nécessiter un traitement spécial.

2. **Performance**: Le calcul de hash SHA-256 est rapide, mais pour des volumes très importants, envisagez un cache.

3. **Format canonique**: Le format est conçu pour être lisible et déterministe, pas pour être analysé à nouveau.

4. **Types de maps**: Les maps JSON doivent avoir un champ "type" pour être correctement analysées.

## Intégration avec le réseau RETE

Ce module est conçu pour faciliter la construction de chaînes alpha optimisées dans le réseau RETE. Les conditions extraites peuvent être utilisées pour :

- Créer des nœuds alpha partagés entre plusieurs règles
- Optimiser le réseau en détectant les conditions communes
- Construire des index de conditions pour un accès rapide
- Analyser la complexité des règles

## Exemples complets

### Exemple 1: Analyse d'une règle métier

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/constraint"
    "github.com/treivax/tsd/rete"
)

func main() {
    // Règle: Les employés de plus de 18 ans avec un salaire >= 50000 sont éligibles
    expr := constraint.LogicalExpression{
        Type: "logicalExpr",
        Left: constraint.BinaryOperation{
            Type:     "binaryOperation",
            Left:     constraint.FieldAccess{Object: "e", Field: "age"},
            Operator: ">",
            Right:    constraint.NumberLiteral{Value: 18},
        },
        Operations: []constraint.LogicalOperation{
            {
                Op: "AND",
                Right: constraint.BinaryOperation{
                    Type:     "binaryOperation",
                    Left:     constraint.FieldAccess{Object: "e", Field: "salary"},
                    Operator: ">=",
                    Right:    constraint.NumberLiteral{Value: 50000},
                },
            },
        },
    }

    conditions, opType, err := rete.ExtractConditions(expr)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Type d'opérateur: %s\n", opType)
    fmt.Printf("Nombre de conditions: %d\n\n", len(conditions))

    for i, cond := range conditions {
        fmt.Printf("Condition %d:\n", i+1)
        fmt.Printf("  Type: %s\n", cond.Type)
        fmt.Printf("  Opérateur: %s\n", cond.Operator)
        fmt.Printf("  Canonique: %s\n", rete.CanonicalString(cond))
        fmt.Printf("  Hash: %s\n\n", cond.Hash)
    }
}
```

### Exemple 2: Détection de partage de conditions entre règles

```go
func analyzeRuleSharing(rules []RuleDefinition) {
    allConditions := make(map[string][]string) // hash -> ruleIDs

    for _, rule := range rules {
        conditions, _, err := rete.ExtractConditions(rule.Constraints)
        if err != nil {
            continue
        }

        for _, cond := range conditions {
            allConditions[cond.Hash] = append(allConditions[cond.Hash], rule.ID)
        }
    }

    // Trouver les conditions partagées
    fmt.Println("Conditions partagées entre règles:")
    for hash, ruleIDs := range allConditions {
        if len(ruleIDs) > 1 {
            fmt.Printf("  Hash %s utilisé par: %v\n", hash[:8], ruleIDs)
        }
    }
}
```

## Références

- Code source: `tsd/rete/alpha_chain_extractor.go`
- Tests: `tsd/rete/alpha_chain_extractor_test.go`
- Documentation RETE: `tsd/rete/README.md`
- Alpha Chains: `tsd/ALPHA_CHAINS_README.md`
