# TSD Examples

Ce répertoire contient des exemples d'utilisation du système TSD (Type System with Dependencies) et du moteur de règles RETE.

## Exemples disponibles

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

### Arithmétiques (dans les arguments)
- `+` : Addition
- `-` : Soustraction
- `*` : Multiplication
- `/` : Division

## Fonctions d'agrégation

| Fonction | Description | Exemple |
|----------|-------------|---------|
| `AVG(expr)` | Moyenne | `avg_sal: AVG(e.salary)` |
| `SUM(expr)` | Somme | `total: SUM(e.amount)` |
| `COUNT(expr)` | Comptage | `num: COUNT(e.id)` |
| `MIN(expr)` | Minimum | `min_age: MIN(p.age)` |
| `MAX(expr)` | Maximum | `max_sal: MAX(e.salary)` |

## Ressources

- [Documentation complète](../docs/multiple_actions.md)
- [Tests d'intégration](../test/integration/)
- [Grammaire PEG](../constraint/grammar/constraint.peg)
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