# Module Constraint

Ce module contient la grammaire, le parser et les utilitaires pour le système de contraintes et règles métier.

## Structure du module

```
constraint/
├── go.mod                  # Configuration du module Go
├── build.sh               # Script de construction
├── README.md              # Cette documentation
├── api.go                 # API publique du module
├── constraint_types.go    # Définitions des types d'AST
├── constraint_utils.go    # Fonctions utilitaires de validation
├── parser.go             # Parser généré par pigeon (ne pas modifier)
├── cmd/
│   ├── go.mod            # Module pour l'exécutable
│   └── main.go           # Point d'entrée principal
├── grammar/
│   ├── constraint.peg    # Grammaire PEG source
│   └── SetConstraint.g4  # Grammaire ANTLR (alternative)
├── tests/
│   ├── test_input.txt           # Tests de base
│   ├── test_type_valid.txt      # Tests de validation de types
│   ├── test_actions.txt         # Tests d'actions
│   ├── test_multi_expressions.txt
│   ├── test_multiple_actions.txt
│   └── ... (autres fichiers de test)
└── docs/
    ├── GUIDE_CONTRAINTES.md     # Guide de référence complet
    ├── TUTORIEL_CONTRAINTES.md  # Tutoriel d'apprentissage
    └── PARSER_README.md         # Documentation du parser
```

## Installation et construction

### Prérequis

1. **Go 1.21 ou supérieur**
2. **Pigeon** (générateur de parser PEG pour Go)
   ```bash
   go install github.com/mna/pigeon@latest
   ```

### Construction

```bash
# Option 1: Utiliser le script de build (recommandé)
./build.sh

# Option 2: Build manuel
cd constraint
pigeon -o parser.go grammar/constraint.peg
cd cmd
go build -o constraint-parser main.go
```

## Utilisation

### Ligne de commande

```bash
# Depuis le dossier constraint/cmd
./constraint-parser ../tests/test_input.txt

# Ou directement avec go run
go run main.go ../tests/test_input.txt
```

### En tant que bibliothèque Go

```go
package main

import (
    "fmt"
    "log"
    "os"
    "constraint"
)

func main() {
    // Lire un fichier de contraintes
    input, err := os.ReadFile("example.txt")
    if err != nil {
        log.Fatal(err)
    }

    // Parser le contenu
    result, err := constraint.ParseConstraint("example.txt", input)
    if err != nil {
        log.Fatal("Erreur de parsing:", err)
    }

    // Valider la structure
    err = constraint.ValidateConstraintProgram(result)
    if err != nil {
        log.Fatal("Erreur de validation:", err)
    }

    fmt.Println("Parsing et validation réussis!")
}
```

## Grammaire

La grammaire est définie dans `grammar/constraint.peg` et supporte :

- **Définition de types** avec champs typés
- **Variables typées** dans des ensembles
- **Contraintes logiques** avec opérateurs de comparaison
- **Actions** déclenchées par les règles
- **Expressions arithmétiques** 
- **Accès aux champs** avec notation pointée

### Exemple de syntaxe

```
type Personne : < nom: string, age: number, adulte: bool >
type Animal : < nom: string, age: number, domestique: bool >

{ p: Personne, a: Animal } / p.age > 18 AND a.domestique = true ==> adoption(p, a)
```

## Tests

Les fichiers de test dans `tests/` couvrent différents aspects :

- `test_input.txt` : Exemple de base
- `test_type_valid.txt` : Validation de types
- `test_actions.txt` : Tests d'actions  
- `test_field_*.txt` : Tests d'accès aux champs
- `test_multi_expressions.txt` : Règles multiples

## Documentation

- **[Guide de référence](docs/GUIDE_CONTRAINTES.md)** : Documentation complète de la syntaxe
- **[Tutoriel pratique](docs/TUTORIEL_CONTRAINTES.md)** : Apprentissage par l'exemple
- **[Parser README](docs/PARSER_README.md)** : Documentation technique du parser

## Développement

### Modifier la grammaire

1. Éditez `grammar/constraint.peg`
2. Régénérez le parser : `pigeon -o parser.go grammar/constraint.peg`
3. Testez avec `./build.sh`

### Ajouter de nouveaux tests

1. Créez un fichier `tests/test_nouveau_cas.txt`
2. Testez avec : `./constraint-parser tests/test_nouveau_cas.txt`

### Types d'erreurs communes

- **Erreur de syntaxe** : Vérifiez la grammaire PEG
- **Type non défini** : Assurez-vous que tous les types référencés sont déclarés
- **Champ inexistant** : Vérifiez les noms de champs dans les accès
- **Erreur de compilation Go** : Vérifiez que le parser est régénéré

## API du module

### Fonctions publiques

```go
// ParseConstraint analyse un fichier de contraintes
func ParseConstraint(filename string, input []byte) (interface{}, error)

// ValidateConstraintProgram valide un programme d'AST 
func ValidateConstraintProgram(result interface{}) error

// ParseConstraintFile analyse depuis le système de fichiers
func ParseConstraintFile(filename string) (interface{}, error)
```

### Types principaux

```go
type Program struct {
    Types       []TypeDefinition `json:"types"`
    Expressions []Expression     `json:"expressions"`
}

type TypeDefinition struct {
    Type   string  `json:"type"`
    Name   string  `json:"name"`
    Fields []Field `json:"fields"`
}

type Expression struct {
    Type        string      `json:"type"`
    Set         Set         `json:"set"`
    Constraints interface{} `json:"constraints"`
    Action      *Action     `json:"action,omitempty"`
}
```