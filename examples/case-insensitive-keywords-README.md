# Mots-Clés Insensibles à la Casse

## 📋 Description

TSD supporte maintenant les mots-clés de la grammaire dans **trois formes de casse** :
- **UPPERCASE** : Style SQL traditionnel (ex: `AND`, `OR`, `NOT`)
- **lowercase** : Style moderne (ex: `and`, `or`, `not`)
- **Capitalized** : Style titre (ex: `And`, `Or`, `Not`)

Cette fonctionnalité améliore la flexibilité et l'expérience utilisateur en permettant d'écrire les règles selon vos préférences de style.

## ✅ Mots-Clés Supportés

### Opérateurs Logiques
- `AND` / `and` / `And`
- `OR` / `or` / `Or`
- `NOT` / `not` / `Not`

### Contraintes Spéciales
- `EXISTS` / `exists` / `Exists`

### Fonctions d'Agrégation
- `AVG` / `avg` / `Avg`
- `COUNT` / `count` / `Count`
- `SUM` / `sum` / `Sum`
- `MIN` / `min` / `Min`
- `MAX` / `max` / `Max`

### Opérateurs de Comparaison
- `IN` / `in` / `In`
- `LIKE` / `like` / `Like`
- `MATCHES` / `matches` / `Matches`
- `CONTAINS` / `contains` / `Contains`

### Fonctions de Manipulation de Chaînes
- `LENGTH` / `length` / `Length`
- `SUBSTRING` / `substring` / `Substring`
- `UPPER` / `upper` / `Upper`
- `LOWER` / `lower` / `Lower`
- `TRIM` / `trim` / `Trim`

### Fonctions Mathématiques
- `ABS` / `abs` / `Abs`
- `ROUND` / `round` / `Round`
- `FLOOR` / `floor` / `Floor`
- `CEIL` / `ceil` / `Ceil`

## 📝 Exemples

### Style UPPERCASE (SQL Traditionnel)

```tsd
rule highSalaryEmployee : {e: Employee} / 
    e.salary > 100000 AND 
    e.active == true AND
    NOT(e.department IN ["Deprecated"]) AND
    LENGTH(e.name) > 5 ==> 
    promote(e.id)
```

### Style lowercase (Moderne)

```tsd
rule highSalaryEmployee : {e: Employee} / 
    e.salary > 100000 and 
    e.active == true and
    not(e.department in ["Deprecated"]) and
    length(e.name) > 5 ==> 
    promote(e.id)
```

### Style Capitalized (Mixte)

```tsd
rule highSalaryEmployee : {e: Employee} / 
    e.salary > 100000 And 
    e.active == true And
    Not(e.department In ["Deprecated"]) And
    Length(e.name) > 5 ==> 
    promote(e.id)
```

### Combinaison de Styles

Vous pouvez même mélanger les styles dans une même règle :

```tsd
rule complexRule : {e: Employee, s: Sale, total: sum(s.amount)} /
    e.active == true And
    exists(s: Sale / s.employeeId == e.id) and
    total > 10000 OR
    length(e.name) < 10 ==>
    processEmployee(e.id)
```

## ❌ Formes Non Supportées

Pour maintenir la lisibilité, seules les trois formes mentionnées ci-dessus sont acceptées. Les formes de casse arbitraires sont **rejetées** :

```tsd
// ❌ INVALIDE - casse mixte arbitraire
rule bad1 : {e: Employee} / e.age > 18 aNd e.age < 65 ==> action()
rule bad2 : {e: Employee} / e.name LiKe "John%" ==> action()
rule bad3 : {e: Employee} / nOt(e.active) ==> action()
rule bad4 : {e: Employee} / eXiStS(s: Sale / s.id == e.id) ==> action()
```

Ces formes produiront une erreur de parsing, ce qui aide à détecter les erreurs de frappe.

## 🎯 Cas d'Usage

### 1. Migration depuis SQL
Si vous êtes habitué au SQL en majuscules, continuez à l'utiliser :
```tsd
rule sqlStyle : {e: Employee} / 
    e.department IN ["Sales", "Marketing"] AND 
    e.salary > 50000 ==> 
    bonus(e.id)
```

### 2. Style Moderne et Épuré
Préférez un style plus moderne avec des minuscules :
```tsd
rule modernStyle : {e: Employee} / 
    e.department in ["Sales", "Marketing"] and 
    e.salary > 50000 ==> 
    bonus(e.id)
```

### 3. Lisibilité avec Capitalisation
Utilisez la capitalisation pour un style intermédiaire :
```tsd
rule readableStyle : {e: Employee} / 
    e.department In ["Sales", "Marketing"] And 
    e.salary > 50000 ==> 
    bonus(e.id)
```

### 4. Agrégations et Fonctions
Toutes les fonctions supportent les trois styles :
```tsd
// UPPERCASE
rule agg1 : {s: Sale, avg: AVG(s.amount)} / avg > 1000 ==> report()

// lowercase
rule agg2 : {s: Sale, avg: avg(s.amount)} / avg > 1000 ==> report()

// Capitalized
rule agg3 : {s:Sale, total:Sum(s.amount)} / total > 1000 ==> alert()
```

### 5. Fonctions d'Accumulation dans AccumulateConstraint
Les fonctions d'accumulation fonctionnent aussi dans les contraintes d'accumulation :
```tsd
// UPPERCASE - Style SQL traditionnel
rule totalSales : {e:Employee} / 
    SUM(s:Sale / s.employeeId == e.id ; s.amount) > 50000 ==> 
    reward(e.id)

// lowercase - Style moderne
rule totalSales : {e:Employee} / 
    sum(s:Sale / s.employeeId == e.id ; s.amount) > 50000 ==> 
    reward(e.id)

// Capitalized - Style titre
rule totalSales : {e:Employee} / 
    Sum(s:Sale / s.employeeId == e.id ; s.amount) > 50000 ==> 
    reward(e.id)
```

Toutes les fonctions d'accumulation supportent les 3 styles :
- `SUM` / `sum` / `Sum`
- `AVG` / `avg` / `Avg`
- `COUNT` / `count` / `Count`
- `MIN` / `min` / `Min`
- `MAX` / `max` / `Max`

## 🔧 Fichier d'Exemple

Le fichier `case-insensitive-keywords.tsd` contient 16 exemples démontrant :
- Tous les mots-clés dans les 3 formes
- Fonctions d'agrégation dans les variables typées
- Fonctions d'accumulation dans AccumulateConstraint
- Combinaisons complexes
- Règles complètes en un seul style
- Règles mixtes combinant plusieurs styles

Pour l'exécuter :
```bash
tsd examples/case-insensitive-keywords.tsd
```

## 📚 Documentation Technique

Pour plus de détails sur l'implémentation :
- Voir `docs/fix-case-insensitive-keywords.md` pour la documentation complète de la correction
- Voir `constraint/parser_case_insensitive_test.go` pour les tests exhaustifs
- Voir `constraint/grammar/constraint.peg` pour la grammaire PEG

## 💡 Recommandations

1. **Cohérence** : Choisissez un style et respectez-le dans tout votre projet
2. **Lisibilité** : Le style lowercase est généralement plus facile à lire
3. **Standards** : Le style UPPERCASE est plus proche du SQL standard
4. **Équipe** : Établissez une convention au sein de votre équipe

## 🐛 Signaler un Problème

Si vous rencontrez des problèmes avec les mots-clés insensibles à la casse :
1. Vérifiez que vous utilisez l'une des trois formes acceptées
2. Consultez les tests dans `constraint/parser_case_insensitive_test.go`
3. Ouvrez une issue sur GitHub avec un exemple de code minimal