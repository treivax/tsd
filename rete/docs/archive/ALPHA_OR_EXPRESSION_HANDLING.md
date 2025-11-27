# Gestion des Expressions OR dans le Moteur RETE

## Vue d'ensemble

Ce document décrit comment le moteur RETE de TSD gère les expressions OR et les expressions mixtes (AND+OR). Contrairement aux expressions AND qui peuvent être décomposées en chaînes d'AlphaNodes, les expressions OR sont traitées comme des nœuds atomiques uniques pour garantir la sémantique correcte de l'évaluation.

## Principes de Conception

### 1. OR n'est PAS décomposé

Les expressions OR ne sont **jamais** décomposées en plusieurs AlphaNodes. Une expression comme `p.status == "VIP" OR p.age > 18` crée un seul AlphaNode qui évalue l'expression complète.

**Raison**: La décomposition d'un OR en chaîne changerait la sémantique. Un fait doit satisfaire AU MOINS UNE des conditions OR pour passer, pas TOUTES.

```
❌ INCORRECT (décomposition en chaîne):
TypeNode → AlphaNode(status="VIP") → AlphaNode(age>18) → Terminal
(nécessiterait les DEUX conditions)

✅ CORRECT (nœud unique):
TypeNode → AlphaNode(status="VIP" OR age>18) → Terminal
(une seule des deux conditions suffit)
```

### 2. OR est normalisé pour le partage

Bien que non décomposé, un OR est **normalisé** pour permettre le partage d'AlphaNodes entre règles différentes avec des termes dans un ordre différent.

**Exemple**:
```
Règle 1: p.status == "VIP" OR p.age > 18
Règle 2: p.age > 18 OR p.status == "VIP"
```

Ces deux règles partagent le même AlphaNode après normalisation car elles sont sémantiquement identiques.

### 3. Expressions mixtes (AND + OR)

Les expressions contenant à la fois AND et OR sont également traitées comme un seul AlphaNode normalisé.

## Architecture

### Flux de Traitement

```
Expression OR
    ↓
AnalyzeExpression() → ExprTypeOR
    ↓
NormalizeORExpression() → Termes triés
    ↓
CreateAlphaNode() → AlphaNode unique
    ↓
ConditionHash() → Hash pour partage
```

### Composants Clés

#### 1. `expression_analyzer.go`

**Fonction**: `AnalyzeExpression(expr interface{}) (ExpressionType, error)`

Détecte le type d'expression:
- `ExprTypeOR`: Expression OR pure
- `ExprTypeMixed`: Expression contenant AND et OR
- `ExprTypeAND`: Expression AND pure

```go
// Exemples de détection
p.age > 18 OR p.status == "VIP"           → ExprTypeOR
(p.age > 18 OR p.status == "VIP") AND ... → ExprTypeMixed
p.age > 18 AND p.status == "VIP"          → ExprTypeAND
```

#### 2. `alpha_chain_extractor.go`

**Fonction**: `NormalizeORExpression(expr interface{}) (interface{}, error)`

Normalise une expression OR:
1. Extrait tous les termes OR
2. Génère une représentation canonique pour chaque terme
3. Trie les termes par ordre alphabétique de leur représentation
4. Reconstruit l'expression avec les termes triés

**Exemple**:
```go
Input:  p.status == "VIP" OR p.age > 18
        ↓
Termes: ["p.status == VIP", "p.age > 18"]
        ↓
Triés:  ["p.age > 18", "p.status == VIP"]
        ↓
Output: p.age > 18 OR p.status == "VIP"
```

**Propriétés garanties**:
- **Idempotence**: Normaliser deux fois donne le même résultat
- **Ordre indépendant**: `A OR B` et `B OR A` donnent le même résultat
- **Déterminisme**: Même entrée → toujours même sortie

#### 3. `constraint_pipeline_helpers.go`

**Fonction**: `createAlphaNodeWithTerminal(...)`

Gère la création d'AlphaNodes avec traitement spécial pour OR:

```go
// Pseudo-code du flux
if exprType == ExprTypeOR {
    normalizedExpr = NormalizeORExpression(actualCondition)
    condition = wrap(normalizedExpr)
    return createSimpleAlphaNode(condition) // UN SEUL nœud
}

if exprType == ExprTypeMixed {
    normalizedExpr = NormalizeORExpression(actualCondition)
    condition = wrap(normalizedExpr)
    return createSimpleAlphaNode(condition) // UN SEUL nœud
}

if exprType == ExprTypeAND {
    // Décomposition en chaîne possible
    return buildChain(extractConditions(actualCondition))
}
```

**Important**: Le traitement OR/Mixed doit être fait AVANT le check `CanDecompose()` car ces types retournent `false`.

#### 4. `evaluator_constraints.go`

**Fonction**: `evaluateConstraintMap(...)`

Évalue les expressions OR au runtime:
- Extrait l'expression wrappée
- Détecte si c'est une `LogicalExpression`
- Route vers `evaluateLogicalExpression()` qui gère AND/OR

```go
// Évaluation d'un OR
left_result = evaluate(left_term)
result = left_result

for each operation in operations:
    right_result = evaluate(operation.right)
    if operation.op == "OR":
        result = result || right_result  // Court-circuit
```

## Tests

### Suite de Tests Complète

Le fichier `alpha_or_expression_test.go` contient 5 tests:

#### 1. `TestOR_SingleNode_NotDecomposed`

Vérifie qu'une expression OR crée un seul AlphaNode.

```go
Expected: 1 AlphaNode
Actual:   ✓ 1 AlphaNode créé
```

#### 2. `TestOR_Normalization_OrderIndependent`

Vérifie que l'ordre des termes n'affecte pas le hash.

```go
Expression 1: p.status == "VIP" OR p.age > 18
Expression 2: p.age > 18 OR p.status == "VIP"

Hash 1: alpha_84ef332f520d58e7
Hash 2: alpha_84ef332f520d58e7 ✓ IDENTIQUE
```

#### 3. `TestMixedAND_OR_SingleNode`

Vérifie qu'une expression mixte crée un seul AlphaNode.

```go
Expression: (p.age > 18 OR p.status == "VIP") AND p.country == "FR"

Expected: 1 AlphaNode
Actual:   ✓ 1 AlphaNode créé
```

#### 4. `TestOR_FactPropagation_Correct`

Vérifie que les faits se propagent correctement à travers un OR.

```go
Fait 1: status="VIP", age=15   → PASSE (satisfait 1ère condition)
Fait 2: status="Regular", age=25 → PASSE (satisfait 2ème condition)
Fait 3: status="VIP", age=30   → PASSE (satisfait les deux)
Fait 4: status="Regular", age=16 → BLOQUÉ (ne satisfait aucune)

Expected: 3 faits propagés
Actual:   ✓ 3 faits propagés
```

#### 5. `TestOR_SharingBetweenRules`

Vérifie le partage d'AlphaNode entre règles avec OR dans ordre différent.

```go
Règle 1: p.status == "VIP" OR p.age > 18
Règle 2: p.age > 18 OR p.status == "VIP"

Expected: 1 AlphaNode partagé, 2 TerminalNodes
Actual:   ✓ 1 AlphaNode, 2 TerminalNodes
```

### Exécution des Tests

```bash
# Tests OR spécifiques
go test -v -run "TestOR_|TestMixedAND_OR" ./rete

# Tous les tests RETE
go test ./rete
```

## Exemples d'Usage

### Règle Simple avec OR

```tsd
rule "VIP_or_Adult" {
    when
        p: Person(p.status == "VIP" OR p.age > 18)
    then
        log("Eligible customer")
}
```

**Résultat**:
- 1 AlphaNode créé: `alpha_84ef332f520d58e7`
- Évalue l'expression OR complète sur chaque fait Person
- Propage uniquement les faits satisfaisant au moins une condition

### Partage entre Règles

```tsd
rule "Rule1" {
    when
        p: Person(p.status == "VIP" OR p.age > 18)
    then
        action1()
}

rule "Rule2" {
    when
        p: Person(p.age > 18 OR p.status == "VIP")  // Ordre différent!
    then
        action2()
}
```

**Résultat**:
- 1 seul AlphaNode partagé (grâce à la normalisation)
- 2 TerminalNodes distincts (un par règle)
- Gain de mémoire: ~50% (1 nœud au lieu de 2)

### Expression Mixte

```tsd
rule "ComplexCondition" {
    when
        p: Person(
            (p.age > 18 OR p.status == "VIP") AND
            p.country == "FR"
        )
    then
        action()
}
```

**Résultat**:
- 1 AlphaNode avec l'expression mixte complète
- Pas de décomposition en chaîne

## Performances

### Métriques

| Scénario | AlphaNodes Créés | Gain Partage |
|----------|------------------|--------------|
| 1 règle OR | 1 | - |
| 2 règles OR identiques (ordre différent) | 1 | 50% |
| Expression mixte | 1 | N/A |

### Optimisations

1. **Normalisation**: Coût O(n log n) où n = nombre de termes OR
2. **Hashing**: Calculé une seule fois à la création
3. **Évaluation**: Court-circuit pour OR (arrêt dès qu'un terme est vrai)

## Limitations et Considérations

### 1. OR Complexes Imbriqués

Les OR imbriqués sont supportés mais traités comme un seul bloc:

```tsd
// Supporté mais non optimisé davantage
p: Person((a OR b) OR (c OR d))
```

### 2. Combinaisons AND/OR

Les expressions mixtes complexes sont normalisées mais le tri peut être limité:

```tsd
// Normalisation partielle
p: Person(a AND (b OR c) AND d)
```

### 3. Performance avec Nombreux Termes

Un OR avec beaucoup de termes reste un seul nœud mais l'évaluation peut être plus lente:

```tsd
// Performance O(n) pour chaque fait
p: Person(cond1 OR cond2 OR ... OR cond50)
```

## Migration et Compatibilité

### Versions Antérieures

Avant cette implémentation, les OR pouvaient être décomposés incorrectement ou mal évalués.

**Migration**: Aucune action requise. Le nouveau comportement est automatique et rétrocompatible.

### Format de Condition

Les conditions OR sont stockées comme:

```json
{
  "type": "constraint",
  "constraint": {
    "type": "logicalExpr",
    "left": { ... },
    "operations": [
      { "op": "OR", "right": { ... } }
    ]
  }
}
```

## Debugging

### Logs de Création

```
ℹ️  Expression OR détectée, normalisation et création d'un nœud alpha unique
✨ Nouveau AlphaNode partageable créé: alpha_84ef332f520d58e7 (hash: alpha_84ef332f520d58e7)
✓ AlphaNode alpha_84ef332f520d58e7 connecté au TypeNode Person
```

### Logs de Partage

```
ℹ️  Expression OR détectée, normalisation et création d'un nœud alpha unique
♻️  AlphaNode partagé réutilisé: alpha_84ef332f520d58e7 (hash: alpha_84ef332f520d58e7)
✓ Règle rule2 attachée à l'AlphaNode partagé alpha_84ef332f520d58e7 via terminal rule2_terminal
```

### Vérification du Hash

```go
// Calculer le hash d'une expression normalisée
normalizedExpr, _ := NormalizeORExpression(expr)
condition := map[string]interface{}{
    "type": "constraint",
    "constraint": normalizedExpr,
}
hash, _ := ConditionHash(condition, "p")
fmt.Printf("Hash: %s\n", hash)
```

## Références

- `expression_analyzer.go`: Détection des types d'expressions
- `alpha_chain_extractor.go`: Normalisation OR
- `constraint_pipeline_helpers.go`: Création d'AlphaNodes
- `evaluator_constraints.go`: Évaluation runtime
- `alpha_or_expression_test.go`: Suite de tests complète
- `ALPHA_NODE_SHARING.md`: Documentation sur le partage d'AlphaNodes

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

---

**Dernière mise à jour**: 2025-01-XX  
**Version**: 1.0.0  
**Auteur**: TSD Team