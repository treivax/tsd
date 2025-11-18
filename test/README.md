# 🧪 Tests TSD

Ce dossier contient tous les tests du projet TSD organisés selon les bonnes pratiques.

## 📁 Structure

```
test/
├── unit/           # Tests unitaires des modules individuels
├── integration/    # Tests d'intégration système
├── benchmark/      # Tests de performance
└── coverage/       # Tests de couverture fonctionnelle
    └── alpha/      # Tests de couverture des Alpha nodes
```

## 🚀 Exécution des Tests

### Tests de Couverture Alpha
```bash
# Exécuter tous les tests Alpha
go run test/coverage/alpha_coverage_runner.go

# Tests individuels disponibles dans test/coverage/alpha/
ls test/coverage/alpha/*.constraint
```

### Tests Unitaires
```bash
# Exécuter tous les tests unitaires
go test ./test/unit/...

# Test avec couverture
go test -cover ./test/unit/...
```

### Tests d'Intégration
```bash
# Tests d'intégration complets
go test ./test/integration/...
```

### Benchmarks
```bash
# Tests de performance
go test -bench=. ./test/benchmark/...
```

## 📋 Tests de Couverture Alpha

Les tests Alpha valident le fonctionnement de tous les opérateurs:

- **Booléens**: `==`, `!=` avec `true`/`false`
- **Comparaisons**: `>`, `<`, `>=`, `<=`
- **Chaînes**: Égalité et patterns
- **Fonctions**: `LENGTH()`, `ABS()`, `UPPER()`
- **Patterns**: `CONTAINS`, `LIKE`, `MATCHES`, `IN`
- **Négations**: `NOT()` avec tous les opérateurs

## ✅ Statut des Tests

- **26 tests Alpha** : 100% conformes
- **Tous opérateurs** : Fonctionnels
- **Négations complexes** : Supportées

## 🔧 Ajout de Nouveaux Tests

1. Créer un fichier `.constraint` avec la règle
2. Créer un fichier `.facts` avec les données test
3. Placer dans le dossier approprié
4. Exécuter le runner pour valider
