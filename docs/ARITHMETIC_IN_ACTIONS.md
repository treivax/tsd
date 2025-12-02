# Expressions Arithmétiques dans les Actions

## Vue d'ensemble

Le système TSD permet maintenant d'utiliser des expressions arithmétiques directement dans les actions pour calculer dynamiquement des valeurs lors de la création ou modification de faits. Cette fonctionnalité permet d'effectuer des calculs en utilisant les variables liées par la règle.

## Opérateurs Supportés

Les opérateurs arithmétiques suivants sont disponibles :

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| `+` | Addition | `a.age + 5` |
| `-` | Soustraction | `a.age - e.age` |
| `*` | Multiplication | `p.price * p.quantity` |
| `/` | Division | `total / count` |
| `%` | Modulo | `value % 10` |

## Utilisation

### 1. Création de fait avec calcul arithmétique

Vous pouvez utiliser des expressions arithmétiques pour calculer la valeur d'un attribut lors de la création d'un nouveau fait :

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

Dans cet exemple :
- `a.age - e.age` calcule la différence d'âge entre l'adulte et l'enfant
- Le résultat est assigné au champ `ageParentALaNaissance` du nouveau fait `Naissance`

### 2. Modification de fait avec calcul arithmétique

Les expressions arithmétiques peuvent également être utilisées pour modifier un champ existant :

```tsd
{ p: Person } / p.age < 30 
==> setFact(p[bonus] = p.salary * 0.1)
```

Ici, le champ `bonus` est calculé comme 10% du salaire.

### 3. Expressions arithmétiques imbriquées

Les expressions peuvent être combinées et imbriquées pour des calculs complexes :

```tsd
{ prod: Product } / prod.available = true
==> setFact(
    Invoice(
        productId: prod.id,
        subtotal: prod.price * prod.quantity,
        tax: (prod.price * prod.quantity) * 0.20,
        total: (prod.price * prod.quantity) * 1.20
    )
)
```

### 4. Combinaison avec des valeurs littérales

Vous pouvez mixer variables et valeurs littérales :

```tsd
{ s: Score } / s.points > 100
==> setFact(s[adjustedScore] = s.points + 50)
```

## Format Interne

Les expressions arithmétiques sont représentées en JSON avec la structure suivante :

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

### Types supportés dans les opérations

- **fieldAccess** : Accès à un champ d'une variable (`a.age`)
- **number** : Valeur numérique littérale (`42`, `3.14`)
- **binaryOperation** : Expression arithmétique imbriquée

## Exemple Complet

Voici un exemple complet utilisant plusieurs types d'expressions arithmétiques :

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

Dans cet exemple :
1. La règle se déclenche quand un adulte et un enfant sont liés (père/enfant)
2. Un nouveau fait `Naissance` est créé avec l'âge du parent à la naissance calculé
3. Le fait `Enfant` existant est modifié pour ajouter la différence d'âge

## Gestion des Erreurs

Le système gère automatiquement les cas d'erreur suivants :

### Division par zéro

```tsd
{ n: Numbers } / n.divisor = 0
==> setFact(n[result] = n.value / n.divisor)  // Erreur : "division par zéro"
```

### Modulo par zéro

```tsd
{ n: Numbers } / n.divisor = 0
==> setFact(n[result] = n.value % n.divisor)  // Erreur : "modulo par zéro"
```

### Type incompatible

Les opérations arithmétiques nécessitent des valeurs numériques. Une erreur est levée si les opérandes ne sont pas des nombres :

```tsd
{ p: Person } / p.name != ""
==> setFact(p[invalid] = p.name + 10)  // Erreur : types incompatibles
```

## Validation de Type

Le système valide automatiquement que :
1. Les résultats des calculs correspondent au type attendu du champ
2. Les champs modifiés existent dans la définition du type
3. Les champs requis sont présents lors de la création de faits

## Précédence des Opérateurs

La précédence standard des opérateurs mathématiques est respectée :

1. `*`, `/`, `%` (multiplication, division, modulo) - priorité haute
2. `+`, `-` (addition, soustraction) - priorité basse

Utilisez des parenthèses pour forcer un ordre d'évaluation spécifique :

```tsd
// Sans parenthèses : (10 * 5) + 2 = 52
result: 10 * 5 + 2

// Avec parenthèses : 10 * (5 + 2) = 70
result: 10 * (5 + 2)
```

## Performances

Les expressions arithmétiques sont évaluées de manière récursive lors de l'exécution de l'action. Pour des performances optimales :

- Évitez des imbrications excessivement profondes
- Pré-calculez les valeurs constantes si possible
- Les opérations sont effectuées en `float64` pour la précision

## Limitations

1. Les opérations arithmétiques ne supportent que les types numériques
2. Les comparaisons dans les actions (`<`, `>`, etc.) retournent des booléens
3. Les expressions arithmétiques dans les conditions de règles utilisent la même syntaxe mais sont évaluées différemment

## Tests

Des tests complets sont disponibles dans `rete/action_arithmetic_test.go`, incluant :

- Création de faits avec expressions arithmétiques
- Modification de faits avec calculs
- Expressions imbriquées
- Gestion des erreurs (division par zéro, types incompatibles)
- Opérations avec valeurs littérales et variables

## Voir Aussi

- [Système d'Actions](ACTIONS_SYSTEM.md)
- [Guide des Actions](ACTIONS_GUIDE.md)
- [Action Print](ACTIONS_PRINT.md)
- [Documentation des Types](TYPES.md)