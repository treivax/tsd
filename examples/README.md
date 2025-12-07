# TSD Examples

Ce répertoire contient des exemples d'utilisation du système TSD (Type System with Dependencies) et du moteur de règles RETE.

## Exemples disponibles

### type-casting.tsd

**Description:** Démontre l'utilisation des opérateurs de casting de types en TSD.

**Fonctionnalités démontrées:**
- Casting `number → string` pour concaténation et comparaison
- Casting `string → number` pour calculs arithmétiques
- Casting `bool → string` pour formatage de messages
- Casting `string → bool` pour évaluation conditionnelle
- Casting `number → bool` (0 = false, non-zero = true)
- Casting `bool → number` (false = 0, true = 1)
- Expressions complexes avec multiples casts
- Exemples complets : e-commerce, configuration système, transformation de données

**Utilisation:**
```bash
# Parser et valider
./tsd examples/type-casting.tsd

# Exécuter avec verbose
./tsd examples/type-casting.tsd -verbose
```

**Concepts clés:**
- Cast explicite requis pour conversions de types
- Syntaxe : `cast(expression as type)` ou `(type)expression`
- Politique de typage stricte (pas de conversion implicite)
- Combinaison de casts avec opérateurs arithmétiques et logiques

**Cas d'usage:**
- Formatage de prix et quantités pour affichage
- Conversion de données utilisateur (strings) en nombres pour calcul
- Validation de configuration avec booléens en string
- Comptage de flags booléens via conversion en nombre

---

### string-operations.tsd

**Description:** Démontre les opérations sur chaînes de caractères en TSD.

**Fonctionnalités démontrées:**
- Concaténation de chaînes avec l'opérateur `+`
- Concaténation avec casting de types (number, bool → string)
- Pattern matching avec `LIKE` (wildcards `%`)
- Recherche de sous-chaînes avec `CONTAINS`
- Expressions régulières avec `MATCHES`
- Appartenance à un ensemble avec `IN`
- Exemples complets : e-commerce, traitement de logs, validation de données

**Utilisation:**
```bash
# Parser et valider
./tsd examples/string-operations.tsd

# Exécuter avec verbose
./tsd examples/string-operations.tsd -verbose
```

**Concepts clés:**
- Concaténation stricte : `string + string` uniquement (pas de conversion implicite)
- Pattern matching : `email LIKE "%@gmail.com"`
- Containment : `content CONTAINS "urgent"`
- Regex : `timestamp MATCHES "[0-9]{4}-[0-9]{2}-[0-9]{2}"`
- Set membership : `status IN ["pending", "processing"]`

**Cas d'usage:**
- Formatage de messages et rapports
- Validation d'emails, téléphones, et autres entrées utilisateur
- Filtrage de logs par patterns
- Détection de contenu sensible ou urgent
- Routage basé sur patterns (régions, statuts, etc.)

---

### multiple_actions_example.tsd

**Description:** Démontre l'utilisation des actions multiples dans les règles RETE.

**Fonctionnalités démontrées:**
- Actions multiples séparées par des virgules
- Règles avec deux actions
- Règles avec trois actions ou plus
- Actions multiples avec agrégation (AVG)
- Utilisation de champs et variables dans les arguments
- Combinaison de types Person, Department et Employee

**Utilisation:**
```bash
# Parser et valider
cd constraint
go run cmd/main.go ../examples/multiple_actions_example.tsd

# Exécuter avec le pipeline complet
cd ..
./tsd -file examples/multiple_actions_example.tsd
```

**Concepts clés:**
- Règle `adult_check` : Deux actions simples (marking + logging)
- Règle `high_earner` : Trois actions pour traiter les hauts salaires
- Règle `dept_avg` : Actions multiples avec agrégation AVG
- Règle `employee_bonus` : Actions multiples avec jointure entre types

**Résultat attendu:**
Le programme parse correctement 4 règles et 9 faits, puis affiche la structure JSON complète avec toutes les actions définies dans le champ `jobs`.

## Comment utiliser les exemples

### 1. Parser uniquement
Pour valider la syntaxe et voir la structure JSON générée :
```bash
cd constraint
go run cmd/main.go ../examples/<nom_fichier>.tsd
```

### 2. Exécution complète
Pour construire le réseau RETE et soumettre les faits :
```bash
./tsd -file examples/<nom_fichier>.tsd
```

### 3. Exécution avec verbose
Pour voir tous les détails de l'exécution :
```bash
./tsd -file examples/<nom_fichier>.tsd -verbose
```

## Structure d'un fichier .tsd

```
// Commentaires (optionnel)

// 1. Définitions de types
type TypeName : <field1: type1, field2: type2, ...>

// 2. Règles
rule rule_name : {patterns} / constraints ==> action1(...), action2(...)

// 3. Faits (optionnel)
TypeName(field1:value1, field2:value2, ...)
```

## Syntaxe des actions multiples

### Format de base
```
rule name : {patterns} / constraints ==> action1(args), action2(args), action3(args)
```

### Exemples

**Deux actions:**
```
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id), log("Adult")
```

**Trois actions:**
```
rule r2 : {p: Person} / p.salary > 50000 ==> 
    flag(p.id), 
    update(p.salary), 
    notify("manager")
```

**Avec agrégation:**
```
rule r3 : {d: Dept, avg: AVG(e.salary)} / {e: Emp} / e.deptId == d.id ==> 
    print(d.name, avg), 
    update_dashboard(d.id, avg), 
    alert(d.id)
```

## Types de données supportés

| Type | Description | Exemple |
|------|-------------|---------|
| `string` | Chaîne de caractères | `"Alice"`, `"test"` |
| `number` | Nombre (entier ou décimal) | `42`, `3.14`, `-10` |
| `bool` | Booléen | `true`, `false` |

## Opérateurs supportés

### Comparaison
- `==` : Égalité
- `!=` : Différence
- `>` : Supérieur
- `<` : Inférieur
- `>=` : Supérieur ou égal
- `<=` : Inférieur ou égal

### Logiques
- `AND` : Et logique
- `OR` : Ou logique
- `NOT` : Négation

### Arithmétiques
- `+` : Addition (nombres) ou Concaténation (chaînes)
- `-` : Soustraction
- `*` : Multiplication
- `/` : Division

### Chaînes de caractères
- `LIKE` : Pattern matching avec wildcards (`%`)
- `CONTAINS` : Recherche de sous-chaîne
- `MATCHES` : Expression régulière
- `IN` : Appartenance à un ensemble

### Casting
- `cast(expr as type)` : Conversion explicite de type
- Types supportés : `number`, `string`, `bool`
- Syntaxe alternative : `(type)expr`

## Fonctions d'agrégation

| Fonction | Description | Exemple |
|----------|-------------|---------|
| `AVG(expr)` | Moyenne | `avg_sal: AVG(e.salary)` |
| `SUM(expr)` | Somme | `total: SUM(e.amount)` |
| `COUNT(expr)` | Comptage | `num: COUNT(e.id)` |
| `MIN(expr)` | Minimum | `min_age: MIN(p.age)` |
| `MAX(expr)` | Maximum | `max_sal: MAX(e.salary)` |

## Ressources

- [Guide utilisateur complet](../docs/USER_GUIDE.md)
- [Guide de démarrage rapide](../docs/QUICK_START.md)
- [Architecture technique](../docs/ARCHITECTURE.md)
- [Guide de contribution](../docs/CONTRIBUTING.md)
- [Référence API](../docs/API_REFERENCE.md)
- [Guide d'authentification](../docs/AUTHENTICATION.md)
- [Tests d'intégration](../tests/e2e/)
- [CHANGELOG](../CHANGELOG.md)

## Contribution

Pour ajouter de nouveaux exemples :

1. Créer un fichier `.tsd` dans ce répertoire
2. Documenter l'exemple dans ce README
3. Ajouter des tests si nécessaire
4. Vérifier que l'exemple parse et s'exécute correctement

## Copyright

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License