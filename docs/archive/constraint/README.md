# Documentation du Module Constraint

Parser et validateur pour le langage de règles TSD.

## 📖 Vue d'ensemble

Le module `constraint` fournit :

- **Parser PEG** : Analyse syntaxique des fichiers `.tsd`
- **Validation de types** : Vérification de la cohérence des types
- **AST Construction** : Génération d'arbres syntaxiques abstraits
- **API Programmatique** : Interface pour construire des programmes TSD

## 📚 Documentation

### Guides Utilisateur

- [**Guide des Contraintes**](GUIDE_CONTRAINTES.md) - Guide complet des contraintes
- [**Tutoriel Actions**](TUTORIEL_ACTIONS.md) - Comment définir et utiliser les actions
- [**Tutoriel Contraintes**](TUTORIEL_CONTRAINTES.md) - Introduction aux contraintes

### Références Techniques

- [**Grammar Complete**](GRAMMAR_COMPLETE.md) - Grammaire PEG complète du langage
- [**Type Validation**](TYPE_VALIDATION.md) - Système de validation des types

## 🚀 Utilisation Rapide

### Parser un Fichier TSD

```go
package main

import (
    "github.com/yourusername/tsd/constraint"
)

func main() {
    // Parser un fichier
    program, err := constraint.ParseFile("rules.tsd")
    if err != nil {
        panic(err)
    }
    
    // Accéder aux types
    for _, typeDef := range program.Types {
        fmt.Printf("Type: %s\n", typeDef.Name)
    }
    
    // Accéder aux règles
    for _, rule := range program.Rules {
        fmt.Printf("Rule: %s\n", rule.Name)
    }
}
```

### Validation des Types

```go
// Valider un programme
validator := constraint.NewTypeValidator()
err := validator.Validate(program)
if err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```

### Construction Programmatique

```go
// Créer un programme par code
program := &constraint.Program{
    Types: []*constraint.TypeDefinition{
        {
            Name: "Person",
            Fields: map[string]string{
                "name": "string",
                "age":  "number",
            },
        },
    },
    Facts: []*constraint.Fact{
        {
            Type: "Person",
            Fields: map[string]interface{}{
                "name": "Alice",
                "age":  30,
            },
        },
    },
}
```

## 📝 Format du Langage TSD

### Types

```tsd
type Person : <name: string, age: number, city: string>
type Company : <name: string, employees: number>
```

### Faits

```tsd
Person(name="Alice", age=30, city="Paris")
Company(name="TechCorp", employees=100)
```

### Règles

```tsd
rule "adult_check" {
    when {
        p: Person(age >= 18)
    }
    then {
        print("Adult: " + p.name)
    }
}
```

### Actions

```tsd
action AdultStatus : <name: string, status: string>

rule "classify_adult" {
    when {
        p: Person(age >= 18)
    }
    then {
        assert(AdultStatus(name=p.name, status="adult"))
    }
}
```

## 🔧 API Principale

### Parsing

```go
// Parser depuis un fichier
program, err := constraint.ParseFile(filename)

// Parser depuis une chaîne
program, err := constraint.ParseString(content)

// Parser depuis un Reader
program, err := constraint.Parse(reader)
```

### Validation

```go
// Créer un validateur
validator := constraint.NewTypeValidator()

// Valider un programme
err := validator.Validate(program)

// Valider un type spécifique
err := validator.ValidateType(typeDef)

// Valider une règle
err := validator.ValidateRule(rule)
```

### Manipulation de l'AST

```go
// Visiter les nœuds
visitor := &MyVisitor{}
constraint.Walk(program, visitor)

// Transformer l'AST
transformer := &MyTransformer{}
newProgram := constraint.Transform(program, transformer)
```

## 🎯 Types Supportés

| Type | Description | Exemple |
|------|-------------|---------|
| `string` | Chaîne de caractères | `"hello"` |
| `number` | Nombre (int ou float) | `42`, `3.14` |
| `boolean` | Booléen | `true`, `false` |
| `list` | Liste de valeurs | `[1, 2, 3]` |
| `map` | Dictionnaire | `{"key": "value"}` |

## ⚠️ Erreurs Courantes

### Type Non Défini

```tsd
# ❌ Erreur: Company n'est pas défini
rule "check" {
    when { c: Company(employees > 10) }
    then { print("Big company") }
}
```

**Solution** : Définir le type avant de l'utiliser

```tsd
# ✅ Correct
type Company : <name: string, employees: number>

rule "check" {
    when { c: Company(employees > 10) }
    then { print("Big company") }
}
```

### Champ Inexistant

```tsd
type Person : <name: string, age: number>

# ❌ Erreur: salary n'existe pas sur Person
rule "check" {
    when { p: Person(salary > 50000) }
    then { print("High salary") }
}
```

**Solution** : Ajouter le champ au type

```tsd
type Person : <name: string, age: number, salary: number>

# ✅ Correct
rule "check" {
    when { p: Person(salary > 50000) }
    then { print("High salary") }
}
```

### Type Incompatible

```tsd
type Person : <name: string, age: number>

# ❌ Erreur: age est un number, pas un string
Person(name="Alice", age="thirty")
```

**Solution** : Utiliser le bon type

```tsd
# ✅ Correct
Person(name="Alice", age=30)
```

## 🔗 Liens Utiles

### Documentation Générale
- [README Principal](../../README.md)
- [Documentation Complète](../../docs/README.md)
- [Tutorial](../../docs/TUTORIAL.md)

### Module RETE
- [RETE Engine](../../rete/README.md)
- [RETE Documentation](../../rete/docs/)

## 🧪 Tests

```bash
# Tous les tests
go test ./constraint/...

# Tests de parsing
go test ./constraint/... -run Parser

# Tests de validation
go test ./constraint/... -run Validation

# Tests avec coverage
go test ./constraint/... -cover
```

## 🤝 Contribution

Pour contribuer au module constraint :

1. Lire [Development Guidelines](../../docs/development_guidelines.md)
2. Modifier la grammaire PEG si nécessaire
3. Ajouter des tests de parsing et validation
4. Soumettre une Pull Request

## 📄 License

Voir [LICENSE](../../LICENSE) à la racine du projet.

---

**Version** : 1.0  
**Dernière mise à jour** : Janvier 2025