# Module RETE - Architecture Refactorisée

## 📁 Structure du Projet

```
rete/
├── pkg/                    # Packages publics
│   ├── domain/            # Types et interfaces du domaine
│   │   ├── facts.go       # Types Fact, Token, WorkingMemory
│   │   ├── interfaces.go  # Interfaces Node, Storage, Logger
│   │   └── errors.go      # Types d'erreurs structurés
│   └── nodes/             # Implémentations des nœuds
│       └── base.go        # Nœud de base commun
├── internal/              # Packages internes
│   ├── config/           # Configuration structurée
│   └── logger/           # Système de logging
├── test/                  # Tests organisés
│   ├── unit/             # Tests unitaires
│   ├── integration/      # Tests d'intégration
│   └── coverage/         # Rapports de couverture
├── docs/                  # Documentation
├── scripts/              # Scripts utilitaires
├── Makefile              # Commandes de build
└── *.go                  # Sources principales à la racine
```

## 🚀 Utilisation Rapide

### Via Makefile
```bash
make help        # Afficher l'aide
make build       # Compiler
make test        # Tests complets avec couverture
make validate    # Valider la structure
make clean       # Nettoyer les artefacts
```

### Via Scripts
```bash
./scripts/run_tests.sh        # Tests avec couverture détaillée
./scripts/validate_structure.sh  # Validation complète
./scripts/clean.sh               # Nettoyage
```

## 🏗️ Architecture

### Packages Publics (`pkg/`)
- **`pkg/domain/`** : Types fondamentaux et interfaces du système RETE
- **`pkg/nodes/`** : Implémentations concrètes des nœuds

### Packages Internes (`internal/`)
- **`internal/config/`** : Gestion de configuration
- **`internal/logger/`** : Système de logging structuré

## 📊 Tests et Couverture

- **Tests unitaires** : `test/unit/`
- **Tests d'intégration** : `test/integration/`
- **Rapports de couverture** : `test/coverage/reports/`

Couverture actuelle : **89%**

## 🔧 Développement

### Prérequis
- Go 1.19+
- Make (optionnel)

### Configuration
```go
import "github.com/treivax/tsd/rete/pkg/domain"

// Utilisation des nouveaux types
fact := domain.NewFact("temperature", 25.5)
wm := domain.NewWorkingMemory()
wm.AddFact(fact)
```

### Ajout de Nouveaux Nœuds
```go
import (
    "github.com/treivax/tsd/rete/pkg/domain"
    "github.com/treivax/tsd/rete/pkg/nodes"
)

type MyNode struct {
    *nodes.BaseNode
    // Champs spécifiques
}

func (n *MyNode) Process(token domain.Token) error {
    // Implémentation spécifique
    return n.ProcessChildren(token)
}
```

## 📝 Bonnes Pratiques Implémentées

### Principes SOLID
- **SRP** : Chaque package a une responsabilité unique
- **OCP** : Interfaces extensibles sans modification
- **LSP** : Substitution des implémentations
- **ISP** : Interfaces ségrégées par rôle
- **DIP** : Dépendances vers les abstractions

### Organisation Go Standard
- `pkg/` pour les API publiques
- `internal/` pour le code privé
- `test/` pour l'organisation des tests
- `docs/` pour la documentation
- `scripts/` pour les utilitaires

### Gestion d'Erreurs
```go
import "github.com/treivax/tsd/rete/pkg/domain"

// Erreurs structurées avec contexte
err := domain.NewValidationError("invalid fact", "field", "temperature")
if domain.IsValidationError(err) {
    // Gestion spécifique
}
```

## 🔄 Migration depuis l'Ancienne Structure

1. **Imports à modifier** :
   ```go
   // Ancien
   import "github.com/treivax/tsd/rete"
   
   // Nouveau
   import "github.com/treivax/tsd/rete/pkg/domain"
   import "github.com/treivax/tsd/rete/pkg/nodes"
   ```

2. **Types déplacés** :
   - `Fact` → `domain.Fact`
   - `Token` → `domain.Token`
   - `WorkingMemory` → `domain.WorkingMemory`

3. **Tests organisés** : Déplacés vers `test/unit/`

## 📖 Documentation Complète

- **[Guide de Refactoring](docs/REFACTORING_RECOMMENDATIONS.md)** : Analyse détaillée
- **[Tests Alpha](docs/ALPHA_TESTS_DOCUMENTATION.md)** : Documentation des tests
- **[Résumé des Tests](docs/TESTS_SUMMARY.md)** : État de la couverture

## 🎯 Prochaines Étapes

1. Migration progressive des fichiers racine vers `pkg/`
2. Implémentation des tests d'intégration
3. Documentation API complète
4. Optimisations de performance

---

> **Note** : Cette refactorisation maintient la compatibilité avec l'API existante tout en introduisant une structure plus maintenable et extensible.