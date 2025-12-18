# Analyse du Parser TSD - Base pour xuple-space

## 🎯 Objectif

Analyser la structure actuelle du parser TSD pour comprendre comment ajouter la commande `xuple-space`.

## 📁 Localisation

### Grammaire PEG
- **Fichier** : `/home/resinsec/dev/tsd/constraint/grammar/constraint.peg`
- **Générateur** : `pigeon` (PEG parser generator pour Go)
- **Commande** : `pigeon -o parser.go constraint.peg`
- **Parser généré** : `/home/resinsec/dev/tsd/constraint/parser.go`

### Structures AST
- **Fichier** : `/home/resinsec/dev/tsd/constraint/constraint_types.go`
- **Package** : `constraint`
- **Structures principales** :
  - `Program` - Racine de l'AST
  - `TypeDefinition` - Définitions de types
  - `ActionDefinition` - Définitions d'actions
  - `Expression` - Règles RETE
  - `Fact` - Faits
  - `Reset` - Commandes reset
  - `RuleRemoval` - Suppressions de règles

## 📊 Structure Actuelle

### Point d'Entrée : Règle `Start`

```peg
Start <- _ statements:StatementList _ EOF {
    // Séparer types, actions, expressions, faits, retractions, ruleRemovals et reset
    types := []interface{}{}
    actions := []interface{}{}
    expressions := []interface{}{}
    facts := []interface{}{}
    retractions := []interface{}{}
    ruleRemovals := []interface{}{}
    resets := []interface{}{}
    
    // Classification des statements
    if statements != nil {
        for _, stmt := range statements.([]interface{}) {
            if stmtMap, ok := stmt.(map[string]interface{}); ok {
                if stmtMap["type"] == "typeDefinition" {
                    types = append(types, stmt)
                } else if stmtMap["type"] == "actionDefinition" {
                    actions = append(actions, stmt)
                } else if stmtMap["type"] == "expression" {
                    expressions = append(expressions, stmt)
                } else if stmtMap["type"] == "fact" {
                    facts = append(facts, stmt)
                } else if stmtMap["type"] == "retraction" {
                    retractions = append(retractions, stmt)
                } else if stmtMap["type"] == "ruleRemoval" {
                    ruleRemovals = append(ruleRemovals, stmt)
                } else if stmtMap["type"] == "reset" {
                    resets = append(resets, stmt)
                }
            }
        }
    }
    
    return map[string]interface{}{
        "types": types,
        "actions": actions,
        "expressions": expressions,
        "facts": facts,
        "retractions": retractions,
        "ruleRemovals": ruleRemovals,
        "resets": resets,
    }, nil
}
```

### Règle `Statement`

```peg
Statement <- TypeDefinition / ActionDefinition / Expression / RemoveRule / RemoveFact / Fact / Reset
```

**Pattern observé** : 
- Chaque type de statement est une alternative dans la règle `Statement`
- Chaque statement retourne une map avec un champ `"type"` unique
- La règle `Start` dispatche selon le type

## 🏗️ Pattern des Commandes Existantes

### 1. TypeDefinition

**Syntaxe** :
```tsd
type Person(#id: string, name: string, age: number)
```

**Règle PEG** :
```peg
TypeDefinition <- "type" _ name:IdentName _ "(" _ fields:FieldList _ ")" {
    return map[string]interface{}{
        "type": "typeDefinition",
        "name": name,
        "fields": fields,
    }, nil
}
```

**Structure AST (Go)** :
```go
type TypeDefinition struct {
    Type   string  `json:"type"`   // "typeDefinition"
    Name   string  `json:"name"`
    Fields []Field `json:"fields"`
}
```

### 2. ActionDefinition

**Syntaxe** :
```tsd
action notify(recipient: string, message: string, priority: number = 1)
```

**Règle PEG** :
```peg
ActionDefinition <- "action" _ name:IdentName _ "(" _ params:ParameterList? _ ")" {
    if params == nil {
        params = []interface{}{}
    }
    return map[string]interface{}{
        "type": "actionDefinition",
        "name": name,
        "parameters": params,
    }, nil
}
```

**Structure AST (Go)** :
```go
type ActionDefinition struct {
    Type       string      `json:"type"`       // "actionDefinition"
    Name       string      `json:"name"`
    Parameters []Parameter `json:"parameters"`
}
```

### 3. Fact

**Syntaxe** :
```tsd
Person(id: "123", name: "John", age: 30)
```

**Règle PEG** :
```peg
Fact <- typeName:IdentName "(" _ fields:FactFieldList _ ")" {
    return map[string]interface{}{
        "type": "fact",
        "typeName": typeName,
        "fields": fields,
    }, nil
}
```

**Structure AST (Go)** :
```go
type Fact struct {
    Type     string                 `json:"type"`     // "fact"
    TypeName string                 `json:"typeName"`
    Fields   map[string]interface{} `json:"fields"`
}
```

### 4. Reset

**Syntaxe** :
```tsd
reset
```

**Règle PEG** :
```peg
Reset <- "reset" {
    return map[string]interface{}{
        "type": "reset",
    }, nil
}
```

**Structure AST (Go)** :
```go
type Reset struct {
    Type string `json:"type"` // "reset"
}
```

## 🔍 Observations Clés

### Pattern Uniforme

1. **Mot-clé initial** : Chaque commande commence par un mot-clé (`type`, `action`, `rule`, `reset`)
2. **Identifiant nommé** : Les commandes avec nom utilisent `IdentName`
3. **Paramètres entre parenthèses** : Lorsqu'il y a des paramètres `(` ... `)`
4. **Corps entre accolades** : Pour structures complexes `{` ... `}`
5. **Retour structuré** : Map avec champ `"type"` discriminant
6. **Séparation espaces** : Règle `_` pour whitespace/commentaires

### Gestion des Listes

**Pattern récurrent** :
```peg
List <- first:Element rest:(_ "," _ Element)* {
    elements := []interface{}{first}
    if rest != nil {
        for _, item := range rest.([]interface{}) {
            elements = append(elements, item.([]interface{})[3])
        }
    }
    return elements, nil
}
```

**Explication** :
- `first` : Premier élément obligatoire
- `rest` : Zéro ou plusieurs éléments additionnels avec séparateur `,`
- Index `[3]` : 4ème élément de la séquence `(_ "," _ Element)`

### Valeurs Optionnelles

**Pattern** :
```peg
optionalValue:(_ "=" _ Value)?
```

**Test nil** :
```go
if optionalValue != nil {
    result["key"] = optionalValue.([]interface{})[3]
}
```

## 📝 Implications pour xuple-space

### Syntaxe Cible

```tsd
xuple-space agents-commands {
    selection: fifo
    consumption: once
    retention: unlimited
}
```

### Structure à Implémenter

1. **Règle principale** : `XupleSpaceDeclaration`
2. **Corps** : `XupleSpaceBody` avec propriétés
3. **Propriétés** : `selection`, `consumption`, `retention`
4. **Valeurs** : Énumérations + paramètres optionnels

### Ajouts Nécessaires

#### 1. Modifier `Statement`
```peg
Statement <- TypeDefinition / ActionDefinition / XupleSpaceDeclaration / Expression / RemoveRule / RemoveFact / Fact / Reset
```

#### 2. Modifier `Start`
```go
xupleSpaces := []interface{}{}
// ... dans la boucle ...
} else if stmtMap["type"] == "xupleSpaceDeclaration" {
    xupleSpaces = append(xupleSpaces, stmt)
}
// ... dans le return ...
"xupleSpaces": xupleSpaces,
```

#### 3. Ajouter structure `Program`
```go
type Program struct {
    Types        []TypeDefinition         `json:"types"`
    Actions      []ActionDefinition       `json:"actions"`
    XupleSpaces  []XupleSpaceDeclaration  `json:"xupleSpaces"`  // NOUVEAU
    Expressions  []Expression             `json:"expressions"`
    Facts        []Fact                   `json:"facts"`
    Resets       []Reset                  `json:"resets"`
    RuleRemovals []RuleRemoval            `json:"ruleRemovals"`
}
```

## 🔧 Outils de Build

### Génération du Parser

```bash
cd /home/resinsec/dev/tsd/constraint/grammar
pigeon -o ../parser.go constraint.peg
```

**Important** : 
- Le fichier `parser.go` est **généré** et ne doit **jamais** être modifié manuellement
- Toutes les modifications doivent se faire dans `constraint.peg`
- Le fichier généré contient l'en-tête : `// Code generated by pigeon; DO NOT EDIT.`

### Validation

```bash
cd /home/resinsec/dev/tsd
make validate  # Validation complète
make test      # Tests unitaires
go fmt ./...   # Formatage
```

## 📚 Références

### Documentation PEG
- [PEG (Parsing Expression Grammar)](https://en.wikipedia.org/wiki/Parsing_expression_grammar)
- [Pigeon - PEG Parser Generator](https://github.com/mna/pigeon)

### Fichiers Projet
- Grammaire : `constraint/grammar/constraint.peg`
- Types AST : `constraint/constraint_types.go`
- Parser généré : `constraint/parser.go` (NE PAS MODIFIER)
- API : `constraint/api.go`

### Standards Projet
- `.github/prompts/common.md` - Standards généraux
- `.github/prompts/review.md` - Checklist revue de code

## ✅ Conclusion

Le parser TSD suit un pattern clair et cohérent :
1. Grammaire PEG déclarative
2. Retour de maps structurées depuis les règles
3. Conversion vers structures Go typées
4. Classification dans `Program` selon le type

Pour ajouter `xuple-space`, il suffit de :
1. ✅ Créer les règles PEG suivant le pattern existant
2. ✅ Définir la structure AST Go
3. ✅ Ajouter dans `Statement` et `Start`
4. ✅ Étendre `Program` avec le nouveau type
5. ✅ Régénérer le parser avec `pigeon`
6. ✅ Créer les tests
