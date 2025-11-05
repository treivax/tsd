# Parser de Contraintes Go avec PEG

Ce projet implémente un parser pour des expressions de contraintes avec des types de données, utilisant PEG (Parsing Expression Grammar) et Go.

## Structure des fichiers

- `constraint.peg` - Grammaire PEG principale
- `constraint_types.go` - Structures de données Go pour l'AST
- `constraint_utils.go` - Fonctions utilitaires et validation
- `constraint_main.go` - Programme principal
- `build.sh` - Script de construction
- `test_input.txt` - Exemple d'entrée
- `go.mod` - Module Go avec dépendances

## Installation et Construction

### 1. Installer pigeon (générateur PEG pour Go)
```bash
go install github.com/mna/pigeon@latest
```

### 2. Générer le parser
```bash
chmod +x build.sh
./build.sh
```

Ou manuellement :
```bash
pigeon -o parser.go constraint.peg
```

### 3. Compilation
```bash
go mod tidy
go build -o constraint-parser parser.go constraint_types.go constraint_utils.go constraint_main.go
```

## Format du langage

### Définition de types
```
type NomType : < champ1: type1, champ2: type2, ... >
```

Types atomiques supportés :
- `string` - Chaînes de caractères  
- `number` - Nombres (entiers et décimaux)
- `bool` - Booléens (true/false)

### Expressions de contraintes (une ou plusieurs)
```
{ variable1: Type1, variable2: Type2 } / contraintes1

{ variable3: Type3, variable4: Type4 } / contraintes2

...
```

### Exemple complet avec plusieurs expressions
```
type Personne : < nom: string, age: number, adulte: bool >
type Animal : < nom: string, age: number, domestique: bool >
type Produit : < prix: number, stock: number, categorie: string >

{ p1: Personne, p2: Personne } / p1.age > p2.age AND p1.adulte = true

{ a: Animal } / a.age > 2 AND a.domestique = true

{ prod: Produit, client: Personne } / prod.prix < 100 AND client.adulte = true
```

## Opérateurs supportés

### Arithmétiques
- `+`, `-`, `*`, `/`

### Comparaisons  
- `=`, `!=`, `<`, `>`, `<=`, `>=`, `==`

### Logiques
- `AND`, `OR`, `&&`, `||`, `&`, `|`

### Accès aux champs
- `variable.champ`

## Usage

### Exécuter le parser
```bash
go run parser.go constraint_types.go constraint_utils.go constraint_main.go test_input.txt
```

### Sortie JSON
Le parser génère un AST au format JSON :
```json
{
  "types": [
    {
      "type": "typeDefinition",
      "name": "Personne", 
      "fields": [
        {"name": "nom", "type": "string"},
        {"name": "age", "type": "number"}
      ]
    }
  ],
  "expression": {
    "type": "expression",
    "set": {
      "type": "set",
      "variables": [
        {"type": "typedVariable", "name": "p", "dataType": "Personne"}
      ]
    },
    "constraints": { ... }
  }
}
```

## Validation

Le parser inclut des fonctions de validation :
- Vérification que tous les types référencés sont définis
- Validation de l'accès aux champs
- Vérification de cohérence des types

## Exemples d'entrées valides

### E-commerce
```
type Produit : < prix: number, stock: number, promotion: bool >
type Client : < age: number, vip: bool, budget: number >

{ p: Produit, c: Client } / p.prix <= c.budget AND (c.vip = true OR p.promotion = true)
```

### Gestion d'utilisateurs
```
type Utilisateur : < score: number, niveau: string, actif: bool >

{ u1: Utilisateur, u2: Utilisateur } / u1.score > u2.score AND u1.actif = true
```

### Contraintes complexes
```
type Compte : < solde: number, type: string, verrouille: bool >

{ c1: Compte, c2: Compte } / c1.solde + c2.solde > 1000 AND c1.type = "premium" AND c2.verrouille = false
```

## Extension

Pour ajouter de nouvelles fonctionnalités :

1. Modifiez `constraint.peg` pour la grammaire
2. Ajoutez les structures correspondantes dans `constraint_parser.go`
3. Régénérez avec `pigeon -o parser.go constraint.peg`
4. Ajoutez la validation dans les fonctions utilitaires

## Dépendances

- Go 1.21+
- github.com/mna/pigeon v1.1.0