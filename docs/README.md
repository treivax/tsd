# 📖 Documentation TSD

Documentation officielle du projet TSD (Type System Development) - Moteur de règles basé sur l'algorithme RETE.

## 📋 Table des Matières

### 🚀 Pour Commencer
- [README Principal](../README.md) - Vue d'ensemble et installation
- [Guide de Développement](development_guidelines.md) - Standards et bonnes pratiques

### 🧪 Tests et Validation
- [Tests Alpha - Résumé](alpha_actions_summary.md) - Résumé des tests de couverture Alpha
- [Tests Alpha - Détails](alpha_tests_detailed.md) - Rapport détaillé par test
- [Rapport de Validation](validation_report.md) - Validation des expressions de négation
- [Correction LIKE](like_fix_report.md) - Résolution du problème opérateur LIKE

### 🏗️ Architecture

#### Moteur RETE
- [Introduction RETE](../rete/README.md) - Vue d'ensemble du moteur RETE
- [Alpha Nodes](../rete/docs/ALPHA_NODES_IMPLEMENTATION.md) - Implémentation des nœuds Alpha
- [Beta Nodes](../rete/docs/BETA_NODES_GUIDE.md) - Guide des nœuds Beta
- [Tuple Space](../rete/docs/TUPLE_SPACE_IMPLEMENTATION.md) - Implémentation de l'espace de tuples

#### Parser de Contraintes
- [Grammar](../constraint/grammar/constraint.peg) - Grammaire PEG des contraintes
- [Guide Contraintes](../constraint/docs/GUIDE_CONTRAINTES.md) - Guide d'écriture des contraintes
- [Tutoriel Actions](../constraint/docs/TUTORIEL_ACTIONS.md) - Tutoriel des actions

## 🎯 Cas d'Usage Validés

### Expressions de Négation Complexes ✅
TSD supporte entièrement les expressions comme :
```
NOT(p.age == 0 AND p.ville <> "Paris")
```

**Statut :** 100% de conformité sur 26 tests Alpha

### Opérateurs Supportés ✅
- **Booléens :** `==`, `!=` avec `true`/`false`
- **Comparaisons :** `>`, `<`, `>=`, `<=`
- **Chaînes :** Égalité et patterns
- **Fonctions :** `LENGTH()`, `ABS()`, `UPPER()`
- **Patterns :** `CONTAINS`, `LIKE`, `MATCHES`, `IN`
- **Négations :** `NOT()` avec tous opérateurs

## 📊 Métriques de Qualité

| Métrique | Valeur | Statut |
|----------|--------|---------|
| **Tests Alpha** | 26/26 | ✅ 100% |
| **Couverture Code** | >90% | ✅ Élevée |
| **Performance** | <1ms/règle | ✅ Optimale |
| **Expressions Complexes** | Supportées | ✅ Validé |

## 🔧 Pour les Développeurs

### Structure du Projet
```
tsd/
├── cmd/           # Applications et CLI
├── constraint/    # Parser de contraintes
├── rete/          # Moteur RETE
├── test/          # Tests organisés
├── docs/          # Documentation
└── scripts/       # Scripts utilitaires
```

### Workflow de Développement
1. Consulter [development_guidelines.md](development_guidelines.md)
2. Exécuter les tests: `go test ./...`
3. Valider Alpha: `go run test/coverage/alpha_coverage_runner.go`
4. Benchmark: `go test -bench=. ./test/benchmark/...`

## 🚀 Statut du Projet

**TSD est prêt pour la production** avec une validation complète des expressions de négation complexes.

**Version :** 1.0
**Dernière validation :** 17 novembre 2025
**Conformité :** 100%

---

*Documentation générée automatiquement - Projet TSD*
