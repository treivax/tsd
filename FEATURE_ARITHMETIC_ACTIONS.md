# Feature: Expressions Arithmétiques dans les Actions

## 📋 Résumé

Cette fonctionnalité permet d'utiliser des expressions arithmétiques directement dans les actions pour calculer dynamiquement des valeurs lors de la création ou modification de faits, en utilisant les variables liées par les règles.

## 🎯 Objectif

Permettre aux utilisateurs de TSD d'effectuer des calculs arithmétiques dans les actions sans avoir à écrire de code supplémentaire ou créer des actions personnalisées pour chaque type de calcul.

## ✨ Fonctionnalités

### Opérateurs Supportés

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| `+` | Addition | `a.age + 5` |
| `-` | Soustraction | `a.age - e.age` |
| `*` | Multiplication | `p.price * p.quantity` |
| `/` | Division | `total / count` |
| `%` | Modulo | `value % 10` |

### Cas d'Usage

#### 1. Création de fait avec calcul

```tsd
{ a: Adulte, e: Enfant } / a.age > e.age AND e.pere = a.ID 
==> setFact(
    Naissance(
        id: e.ID,
        parent: a.ID,
        ageParentALaNaissance: a.age - e.age
    )
)
```

#### 2. Modification de fait avec calcul

```tsd
{ p: Person } / p.age < 30 
==> setFact(p[bonus] = p.salary * 0.1)
```

#### 3. Expressions complexes imbriquées

```tsd
{ prod: Product } / prod.available = true
==> setFact(
    Invoice(
        subtotal: prod.price * prod.quantity,
        tax: (prod.price * prod.quantity) * 0.20,
        total: (prod.price * prod.quantity) * 1.20
    )
)
```

## 🏗️ Architecture

### Composants Modifiés

- **ActionExecutor** (`rete/action_executor.go`)
  - `evaluateBinaryOperation()` : Gère les types `binaryOperation`, `binaryOp`, `binary_operation`
  - `evaluateArithmeticOperation()` : Exécute les calculs arithmétiques
  - `evaluateComparison()` : Gère les comparaisons
  - `areEqual()` : Compare deux valeurs

### Nouveaux Fichiers

```
tsd/
├── rete/
│   ├── action_arithmetic_test.go           # Tests complets
│   └── examples/
│       └── arithmetic_actions_example.go   # Exemples pratiques
└── docs/
    └── ARITHMETIC_IN_ACTIONS.md            # Documentation complète
```

## 🧪 Tests

### Couverture des Tests

- ✅ 7 tests dans `TestArithmeticExpressionsInActions`
- ✅ 5 tests dans `TestArithmeticExpressionEvaluation`
- ✅ 1 test d'intégration complet

### Scénarios Testés

1. Création de fait avec soustraction
2. Modification de fait avec addition
3. Multiplication complexe
4. Expressions imbriquées (ex: `(a * b) + c`)
5. Division par zéro (validation d'erreur)
6. Opération modulo
7. Arithmétique avec valeurs littérales

### Exécuter les Tests

```bash
cd rete
go test -v -run TestArithmetic
```

## 📖 Documentation

### Guides Disponibles

1. **Guide Complet** : [`docs/ARITHMETIC_IN_ACTIONS.md`](docs/ARITHMETIC_IN_ACTIONS.md)
   - Vue d'ensemble détaillée
   - Tous les opérateurs supportés
   - Exemples d'utilisation
   - Gestion des erreurs
   - Format interne

2. **Système d'Actions** : [`rete/ACTIONS_SYSTEM.md`](rete/ACTIONS_SYSTEM.md)
   - Section "Expressions Arithmétiques"
   - Intégration avec le système existant

3. **Changelog** : [`CHANGELOG_ARITHMETIC_ACTIONS.md`](CHANGELOG_ARITHMETIC_ACTIONS.md)
   - Détails techniques complets
   - Historique des modifications

### Exemples Pratiques

Exécuter les exemples :

```bash
cd rete/examples
go run arithmetic_actions_example.go
```

Les exemples incluent :
- Calcul d'âge parent/enfant
- Calcul de facture avec TVA
- Calcul de bonus salarial
- Expressions arithmétiques complexes

## 🛡️ Gestion des Erreurs

Le système gère automatiquement :

### Division/Modulo par Zéro

```tsd
{ n: Numbers } / n.divisor = 0
==> setFact(n[result] = n.value / n.divisor)
// Erreur : "division par zéro"
```

### Types Incompatibles

```tsd
{ p: Person } / p.name != ""
==> setFact(p[invalid] = p.name + 10)
// Erreur : "opération arithmétique nécessite des nombres"
```

### Validation de Type

Le résultat d'un calcul doit correspondre au type attendu du champ dans la définition du type.

## 🚀 Performance

- **Évaluation récursive** : Les expressions sont évaluées à la volée
- **Type de calcul** : Toutes les opérations utilisent `float64`
- **Impact** : Aucun impact sur les actions existantes
- **Optimisation** : Les valeurs constantes peuvent être pré-calculées

## 🔄 Compatibilité

- ✅ **Rétrocompatible** : Aucune modification des actions existantes
- ✅ **Pas de breaking changes** : Ajout de fonctionnalités uniquement
- ✅ **Tests existants** : Tous les tests passent

## 📊 Exemple Complet

### Définition des Types

```tsd
type Adulte {
    ID: string
    age: number
}

type Enfant {
    ID: string
    pere: string
    age: number
    differenceAgeParent: number
}

type Naissance {
    id: string
    parent: string
    ageParentALaNaissance: number
}
```

### Règle avec Calculs

```tsd
{ a: Adulte, e: Enfant } / a.age > e.age AND e.pere = a.ID 
==> setFact(
        Naissance(
            id: e.ID,
            parent: a.ID,
            ageParentALaNaissance: a.age - e.age
        )
    ),
    setFact(e[differenceAgeParent] = a.age - e.age)
```

### Résultat

Si `a.age = 45` et `e.age = 18` :
- Un nouveau fait `Naissance` est créé avec `ageParentALaNaissance = 27`
- Le fait `Enfant` est modifié avec `differenceAgeParent = 27`

## 🔮 Évolutions Futures

Cette fonctionnalité ouvre la voie à :

1. **Fonctions mathématiques** : `sqrt()`, `pow()`, `abs()`, `round()`
2. **Fonctions d'agrégation** : `sum()`, `avg()`, `min()`, `max()`
3. **Expressions conditionnelles** : Opérateur ternaire `? :`
4. **Optimisations** : Pré-calcul des sous-expressions constantes

## 🎓 Précédence des Opérateurs

Les règles mathématiques standard s'appliquent :

1. **Priorité haute** : `*`, `/`, `%`
2. **Priorité basse** : `+`, `-`
3. **Parenthèses** : Force l'ordre d'évaluation

### Exemples

```
10 + 5 * 2    = 10 + (5 * 2)   = 20
(10 + 5) * 2  = 15 * 2          = 30
10 * 5 + 2    = (10 * 5) + 2    = 52
```

## 📝 Notes d'Implémentation

### Format Interne

Les expressions utilisent le type `"binaryOperation"` :

```json
{
    "type": "binaryOperation",
    "operator": "-",
    "left": {"type": "fieldAccess", "object": "a", "field": "age"},
    "right": {"type": "fieldAccess", "object": "e", "field": "age"}
}
```

### Types Supportés dans les Expressions

- `fieldAccess` : Accès à un champ (`a.age`)
- `number` : Valeur numérique littérale (`42`, `3.14`)
- `binaryOperation` : Expression imbriquée

## 👥 Contribution

Pour contribuer à cette fonctionnalité :

1. Lire la documentation dans `docs/ARITHMETIC_IN_ACTIONS.md`
2. Examiner les tests dans `rete/action_arithmetic_test.go`
3. Voir les exemples dans `rete/examples/arithmetic_actions_example.go`

## 📞 Support

Pour toute question ou problème :

1. Consulter la documentation complète
2. Vérifier les exemples pratiques
3. Examiner les tests unitaires
4. Créer une issue avec un exemple reproductible

---

**Version** : 1.0.0  
**Date** : 2025-12-01  
**Statut** : ✅ Production Ready