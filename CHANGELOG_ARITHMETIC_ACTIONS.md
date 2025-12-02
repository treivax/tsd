# Changelog - Expressions Arithmétiques dans les Actions

## [2025-12-01] Ajout du support des expressions arithmétiques dans les actions

### 🎯 Fonctionnalité Ajoutée

Le système d'actions TSD supporte maintenant l'utilisation d'**expressions arithmétiques** pour calculer dynamiquement des valeurs lors de la création ou modification de faits. Cette fonctionnalité permet d'effectuer des calculs en utilisant les variables liées par les règles.

### ✨ Nouveautés

#### Opérateurs Arithmétiques Supportés

- `+` : Addition
- `-` : Soustraction
- `*` : Multiplication
- `/` : Division
- `%` : Modulo

#### Cas d'Utilisation

1. **Création de fait avec calcul arithmétique**
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

2. **Modification de fait avec calcul**
   ```tsd
   { p: Person } / p.age < 30 
   ==> setFact(p[bonus] = p.salary * 0.1)
   ```

3. **Expressions imbriquées complexes**
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

### 🔧 Modifications Techniques

#### Fichiers Modifiés

- `rete/action_executor.go`
  - Ajout de `evaluateBinaryOperation()` pour gérer les types `binaryOperation`, `binaryOp`, `binary_operation`
  - Ajout de `evaluateArithmeticOperation()` pour exécuter les calculs arithmétiques
  - Ajout de `evaluateComparison()` pour les opérations de comparaison
  - Ajout de `areEqual()` pour la comparaison de valeurs
  - Support du modulo (`%`) dans les opérations arithmétiques
  - Gestion des erreurs améliorée (division/modulo par zéro, types incompatibles)

#### Fichiers Ajoutés

- `rete/action_arithmetic_test.go`
  - Tests complets pour les expressions arithmétiques dans les actions
  - 7 scénarios de test principaux
  - Tests de validation des erreurs (division par zéro, types incompatibles)
  - Tests d'expressions imbriquées
  - Test du scénario complet parent/enfant

- `docs/ARITHMETIC_IN_ACTIONS.md`
  - Documentation complète de la fonctionnalité
  - Exemples d'utilisation détaillés
  - Guide de gestion des erreurs
  - Référence des opérateurs supportés

- `rete/examples/arithmetic_actions_example.go`
  - Exemple pratique démontrant 4 scénarios réels
  - Calcul d'âge parent/enfant
  - Calcul de facture avec TVA
  - Calcul de bonus salarial
  - Expressions arithmétiques complexes

#### Documentation Mise à Jour

- `rete/ACTIONS_SYSTEM.md`
  - Ajout d'une section "Expressions Arithmétiques"
  - Exemples d'utilisation intégrés
  - Référence au guide complet

### 📊 Tests

#### Couverture

- **8 tests unitaires** dans `TestArithmeticExpressionsInActions`
  - Création de fait avec soustraction
  - Modification de fait avec addition
  - Multiplication complexe
  - Expressions imbriquées
  - Division par zéro (validation d'erreur)
  - Opération modulo
  - Arithmétique avec valeurs littérales

- **5 tests d'évaluation** dans `TestArithmeticExpressionEvaluation`
  - Addition simple
  - Soustraction avec variables
  - Multiplication
  - Division
  - Modulo

- **1 test d'intégration** dans `TestCompleteScenario_ParentChildAge`
  - Scénario complet avec création et modification de faits
  - Utilisation de multiples expressions arithmétiques

#### Résultats

```
✅ PASS: TestArithmeticExpressionsInActions (7/7 tests)
✅ PASS: TestArithmeticExpressionEvaluation (5/5 tests)
✅ PASS: TestCompleteScenario_ParentChildAge
```

### 🎨 Format Interne

Les expressions arithmétiques utilisent le format `binaryOperation` :

```json
{
    "type": "binaryOperation",
    "operator": "-",
    "left": {
        "type": "fieldAccess",
        "object": "a",
        "field": "age"
    },
    "right": {
        "type": "fieldAccess",
        "object": "e",
        "field": "age"
    }
}
```

### 🛡️ Gestion des Erreurs

Le système gère automatiquement :

- **Division par zéro** : Erreur explicite lors de l'exécution
- **Modulo par zéro** : Erreur explicite lors de l'exécution
- **Types incompatibles** : Validation que les opérandes sont numériques
- **Validation de type** : Les résultats correspondent au type attendu du champ

### 🚀 Performances

- Les expressions sont évaluées récursivement lors de l'exécution
- Toutes les opérations utilisent `float64` pour la précision
- Pas d'impact sur les performances des actions existantes

### 🔄 Compatibilité

- ✅ **Rétrocompatible** : Les actions existantes continuent de fonctionner
- ✅ **Pas de breaking changes** : Ajout de fonctionnalités uniquement
- ✅ **Tests existants** : Tous les tests passent (pas de régression)

### 📚 Documentation

#### Nouveaux Documents

1. `docs/ARITHMETIC_IN_ACTIONS.md` - Guide complet
2. `rete/examples/arithmetic_actions_example.go` - Exemples pratiques

#### Documents Mis à Jour

1. `rete/ACTIONS_SYSTEM.md` - Ajout de la section arithmétique

### 🎯 Exemple d'Utilisation Complet

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

### 🔮 Prochaines Étapes

Cette fonctionnalité pose les bases pour :

- Support de fonctions mathématiques (sqrt, pow, abs, etc.)
- Optimisation des expressions constantes
- Support d'expressions arithmétiques dans plus de contextes

### 👥 Contributeurs

- Implémentation initiale et tests
- Documentation complète
- Exemples pratiques

### 📝 Notes

- La précédence des opérateurs mathématiques est respectée (`*`, `/`, `%` avant `+`, `-`)
- Les parenthèses sont supportées pour forcer l'ordre d'évaluation
- Les expressions peuvent être imbriquées à volonté
- La validation de type est effectuée automatiquement

---

**Date** : 2025-12-01  
**Version** : 1.0.0  
**Statut** : ✅ Complété et testé