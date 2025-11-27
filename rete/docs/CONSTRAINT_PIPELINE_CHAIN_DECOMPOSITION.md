# Constraint Pipeline Chain Decomposition

## Vue d'ensemble

Le système de décomposition en chaînes du Constraint Pipeline intègre l'analyseur d'expressions RETE pour décomposer automatiquement les expressions logiques complexes (AND) en chaînes d'AlphaNodes optimisées et partageables.

## Fonctionnalités

### 1. Analyse Automatique des Expressions

Le pipeline analyse chaque condition de règle avec `AnalyzeExpression()` pour déterminer son type :
- **Simple** : Une condition unique (ex: `p.age > 18`)
- **AND** : Expressions conjonctives (ex: `p.age > 18 AND p.salary >= 50000`)
- **OR** : Expressions disjonctives (ex: `p.age < 18 OR p.age > 65`)
- **NOT** : Négations
- **Arithmetic** : Expressions arithmétiques (ex: `p.salary * 1.1 > 60000`)

### 2. Décomposition en Chaînes

Pour les expressions **AND** décomposables :

1. **Extraction** : `ExtractConditions()` décompose l'expression en conditions atomiques
2. **Normalisation** : `NormalizeConditions()` réordonne les conditions pour maximiser le partage
3. **Construction** : `BuildChain()` crée une chaîne d'AlphaNodes
4. **Partage** : Les nœuds identiques sont automatiquement réutilisés entre règles

### 3. Comportements Spéciaux

#### Expressions OR
Créent un seul AlphaNode normalisé (pas de décomposition) car les conditions OR ne peuvent pas être évaluées séquentiellement.

#### Expressions Simples
Utilisent le comportement classique sans décomposition.

#### Erreurs d'Analyse
Fallback automatique vers le comportement simple pour assurer la robustesse.

## Architecture

### Flux de Traitement

```
Condition → AnalyzeExpression() → CanDecompose() ?
                                          ↓
                                     Type = AND ?
                                          ↓
                                   ExtractConditions()
                                          ↓
                                   NormalizeConditions()
                                          ↓
                                      BuildChain()
                                          ↓
                                   Attach Terminal
```

### Fonctions Clés

#### `createAlphaNodeWithTerminal()`
Fonction principale qui orchestre le processus :
- Analyse l'expression
- Décide de la stratégie (chaîne ou simple)
- Construit le réseau
- Attache le TerminalNode

#### `createSimpleAlphaNodeWithTerminal()`
Fonction de fallback qui implémente le comportement original :
- Crée un seul AlphaNode
- Partage le nœud si possible
- Attache directement le TerminalNode

## Exemples

### Exemple 1 : Condition Simple

```go
// Expression: p.age > 18
condition := map[string]interface{}{
    "type": "binaryOperation",
    "left": constraint.FieldAccess{Object: "p", Field: "age"},
    "operator": ">",
    "right": constraint.NumberLiteral{Value: 18},
}

// Résultat: 1 AlphaNode créé (pas de chaîne)
// ✨ Nouveau AlphaNode créé: alpha_xxx
```

### Exemple 2 : Expression AND (2 conditions)

```go
// Expression: p.age > 18 AND p.salary >= 50000
condition := constraint.LogicalExpression{
    Left: constraint.BinaryOperation{...}, // p.age > 18
    Operations: []constraint.LogicalOperation{
        {Op: "AND", Right: constraint.BinaryOperation{...}}, // p.salary >= 50000
    },
}

// Résultat: Chaîne de 2 AlphaNodes
// 🔍 Expression de type ExprTypeAND détectée
// 🔗 Décomposition en chaîne: 2 conditions détectées
// ✨ Nouveau AlphaNode créé: alpha_aaa (condition 1)
// ✨ Nouveau AlphaNode créé: alpha_bbb (condition 2)
// ✓ TerminalNode attaché au nœud final
```

### Exemple 3 : Expression AND (3 conditions)

```go
// Expression: p.age > 18 AND p.salary >= 50000 AND p.experience > 5
condition := constraint.LogicalExpression{
    Left: constraint.BinaryOperation{...},
    Operations: []constraint.LogicalOperation{
        {Op: "AND", Right: ...},
        {Op: "AND", Right: ...},
    },
}

// Résultat: Chaîne de 3 AlphaNodes
// 🔗 Décomposition en chaîne: 3 conditions détectées
// ✅ Chaîne construite: 3 nœud(s), 0 partagé(s)
```

### Exemple 4 : Deux Règles avec Partage

```go
// Règle 1: p.age > 18 AND p.salary >= 50000
createAlphaNodeWithTerminal(network, "rule1", condition1, ...)
// Résultat: 2 nouveaux AlphaNodes créés

// Règle 2: p.age > 18 AND p.salary >= 50000 (même conditions)
createAlphaNodeWithTerminal(network, "rule2", condition2, ...)
// Résultat: 0 nouveaux AlphaNodes (réutilisation complète)
// ♻️  AlphaNode partagé réutilisé: alpha_aaa
// ♻️  AlphaNode partagé réutilisé: alpha_bbb
// ✅ Chaîne construite: 2 nœud(s), 2 partagé(s)
```

### Exemple 5 : Expression OR

```go
// Expression: p.age < 18 OR p.age > 65
condition := constraint.LogicalExpression{
    Left: constraint.BinaryOperation{...},
    Operations: []constraint.LogicalOperation{
        {Op: "OR", Right: ...},
    },
}

// Résultat: 1 AlphaNode normalisé (pas de chaîne)
// ℹ️  Expression OR détectée, création d'un nœud alpha normalisé unique
```

## Logging

Le système fournit un logging détaillé avec des emojis pour faciliter le débogage :

### Messages de Décomposition
- `🔍 Expression de type ExprTypeAND détectée, tentative de décomposition...`
- `🔗 Décomposition en chaîne: X conditions détectées (opérateur: AND)`
- `📋 Conditions normalisées: X condition(s)`

### Messages de Construction
- `✨ Nouveau AlphaNode créé: [hash]` - Nœud nouvellement créé
- `♻️  AlphaNode partagé réutilisé: [hash]` - Nœud existant réutilisé
- `✅ Chaîne construite: X nœud(s), Y partagé(s)` - Résumé de la construction

### Messages de Fallback
- `ℹ️  Expression de type X non décomposable, utilisation du nœud simple`
- `⚠️  Erreur analyse expression: ..., fallback vers comportement simple`

### Messages de Terminal
- `✓ TerminalNode [id] attaché au nœud final [id] de la chaîne`

## Avantages

### 1. Performance
- **Partage maximal** : Les conditions communes entre règles sont partagées automatiquement
- **Évaluation séquentielle** : Les chaînes AND permettent un court-circuit dès qu'une condition échoue
- **Réduction mémoire** : Moins de nœuds dupliqués dans le réseau

### 2. Maintenabilité
- **Backward compatible** : Les conditions simples fonctionnent exactement comme avant
- **Fallback robuste** : Les erreurs ne cassent pas le système
- **Logging transparent** : Facile de voir ce qui se passe

### 3. Évolutivité
- **Architecture modulaire** : Facile d'ajouter de nouveaux types d'optimisations
- **Réutilisation du code** : Utilise les composants existants (analyzer, extractor, builder)

## Cas d'Usage

### Scénario 1 : Règles RH

```constraint
// Règle 1: Employés éligibles aux bonus
WHEN Employee e WHERE e.age >= 25 AND e.salary < 80000 AND e.performance > 8.0
THEN bonus(e)

// Règle 2: Employés éligibles à la promotion
WHEN Employee e WHERE e.age >= 25 AND e.salary < 80000 AND e.years_service > 5
THEN promote(e)

// Résultat: Les 2 premières conditions (age et salary) sont partagées
// Économie: 2 AlphaNodes partagés au lieu de 4 nœuds distincts
```

### Scénario 2 : Système de Tarification

```constraint
// Tarif Premium
WHEN Customer c WHERE c.age > 18 AND c.credit_score > 700 AND c.income > 50000
THEN premium_rate(c)

// Tarif Standard
WHEN Customer c WHERE c.age > 18 AND c.credit_score > 600
THEN standard_rate(c)

// Résultat: Les 2 conditions age et credit_score partagées partiellement
```

### Scénario 3 : Détection de Fraude

```constraint
// Alerte niveau 1
WHEN Transaction t WHERE t.amount > 1000 AND t.country == "foreign" AND t.time == "night"
THEN alert_level_1(t)

// Alerte niveau 2
WHEN Transaction t WHERE t.amount > 1000 AND t.country == "foreign" AND t.velocity > 5
THEN alert_level_2(t)

// Résultat: Conditions amount et country partagées
```

## Compatibilité

### Rétrocompatibilité
✅ **100% compatible** avec le code existant :
- Les règles existantes fonctionnent sans modification
- Les conditions simples utilisent le même code path qu'avant
- Pas de breaking changes dans l'API

### Dépendances
Nécessite les modules suivants :
- `expression_analyzer.go` - Analyse de types d'expressions
- `alpha_chain_extractor.go` - Extraction de conditions
- `alpha_chain_builder.go` - Construction de chaînes
- `alpha_sharing_manager.go` - Gestion du partage

## Tests

Le système inclut des tests complets :

### Tests Unitaires
- `TestPipeline_SimpleCondition_NoChange` - Conditions simples inchangées
- `TestPipeline_AND_CreatesChain` - Décomposition AND
- `TestPipeline_OR_SingleNode` - OR crée un seul nœud
- `TestPipeline_TwoRules_ShareChain` - Partage entre règles
- `TestPipeline_ErrorHandling_FallbackToSimple` - Gestion d'erreurs
- `TestPipeline_ComplexAND_ThreeConditions` - Chaînes complexes
- `TestPipeline_Arithmetic_NoChain` - Expressions arithmétiques

### Lancer les Tests

```bash
# Tous les tests pipeline
go test ./rete -v -run "TestPipeline_"

# Tests spécifiques de chaîne
go test ./rete -v -run "TestPipeline_.*Chain"

# Test d'une fonctionnalité spécifique
go test ./rete -v -run "TestPipeline_AND_CreatesChain"
```

## Limitations Connues

### 1. Expressions OR
Les expressions OR ne sont pas décomposées car elles nécessitent une évaluation complète de toutes les branches, pas une évaluation séquentielle.

### 2. Expressions Mixtes
Les expressions avec AND et OR mélangés (Mixed) ne sont pas décomposées pour préserver la sémantique d'évaluation.

### 3. Conditions Arithmétiques Complexes
Les expressions arithmétiques complexes ne sont pas décomposées pour éviter des évaluations partielles incorrectes.

## Roadmap Future

### Court Terme
- [ ] Support de la décomposition des expressions NOT avec De Morgan
- [ ] Métriques Prometheus pour le partage de nœuds
- [ ] Dashboard de visualisation des chaînes

### Moyen Terme
- [ ] Optimisation basée sur la sélectivité (réordonnancement)
- [ ] Support des expressions Mixed avec décomposition partielle
- [ ] Cache de décomposition pour éviter les re-analyses

### Long Terme
- [ ] Optimiseur basé sur les coûts
- [ ] Décomposition adaptative selon les statistiques d'exécution
- [ ] Support de la décomposition des expressions OR avec branches multiples

## Références

### Fichiers Associés
- `tsd/rete/constraint_pipeline_helpers.go` - Implémentation principale
- `tsd/rete/constraint_pipeline_chain_test.go` - Tests d'intégration
- `tsd/rete/expression_analyzer.go` - Analyseur d'expressions
- `tsd/rete/alpha_chain_extractor.go` - Extracteur de conditions
- `tsd/rete/alpha_chain_builder.go` - Constructeur de chaînes

### Documentation Connexe
- `EXPRESSION_ANALYZER_V1.3.0_FEATURES.md` - Fonctionnalités de l'analyseur
- `ALPHA_CHAIN_BUILDER.md` - Documentation du builder
- `ALPHA_SHARING_MANAGER.md` - Gestion du partage

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License