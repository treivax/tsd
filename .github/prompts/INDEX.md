# 📇 INDEX - Prompts Réutilisables TSD

> Navigation rapide vers tous les prompts disponibles

---

## 🎯 Par Besoin

### Je veux vérifier la conformité
- **Vérifier la conformité de licence** → [`verify-license-compliance.md`](verify-license-compliance.md)

### Je veux tester
- **Lancer tous les tests** → [`run-tests.md`](run-tests.md)
- **Ajouter des tests** → [`add-test.md`](add-test.md)
- **Débugger un test qui échoue** → [`debug-test.md`](debug-test.md)
- **Valider un réseau RETE** → [`validate-network.md`](validate-network.md)

### Je veux développer
- **Ajouter une fonctionnalité** → [`add-feature.md`](add-feature.md)
- **Modifier un comportement existant** → [`modify-behavior.md`](modify-behavior.md)
- **Corriger un bug** → [`fix-bug.md`](fix-bug.md)
- **Refactoriser du code** → [`refactor.md`](refactor.md)
- **Nettoyage approfondi du code** → [`deep-clean.md`](deep-clean.md)

### Je veux comprendre
- **Expliquer du code** → [`explain-code.md`](explain-code.md)
- **Analyser une erreur** → [`analyze-error.md`](analyze-error.md)
- **Investiguer un comportement étrange** → [`investigate.md`](investigate.md)

### Je veux valider
- **Code review** → [`code-review.md`](code-review.md)
- **Valider un réseau RETE** → [`validate-network.md`](validate-network.md)

### Je veux optimiser
- **Optimiser les performances** → [`optimize-performance.md`](optimize-performance.md)

### Je veux documenter
- **Mettre à jour la documentation** → [`update-docs.md`](update-docs.md)
- **Générer des exemples RETE** → [`generate-examples.md`](generate-examples.md)

### Je veux migrer
- **Migrer version Go ou dépendances** → [`migrate.md`](migrate.md)

### Je veux analyser
- **Statistiques du code** → [`stats-code.md`](stats-code.md)

### Je veux vérifier
- **Conformité de licence** → [`verify-license-compliance.md`](verify-license-compliance.md)

---

## 📚 Par Catégorie

### 🧪 Tests
| Prompt | Description | Taille |
|--------|-------------|--------|
| [run-tests.md](run-tests.md) | Lancer l'ensemble des tests (unitaires + runner universel) | 4K |
| [add-test.md](add-test.md) | Ajouter des tests (unitaires, intégration, RETE) | 12K |
| [debug-test.md](debug-test.md) | Débugger un test qui échoue avec analyse complète | 8K |

### 🔧 Développement
| Prompt | Description | Taille |
|--------|-------------|--------|
| [add-feature.md](add-feature.md) | Ajouter une nouvelle fonctionnalité au projet | 8K |
| [modify-behavior.md](modify-behavior.md) | Modifier un comportement ou fonctionnalité existante | 16K |
| [fix-bug.md](fix-bug.md) | Corriger un bug identifié avec méthodologie complète | 16K |
| [refactor.md](refactor.md) | Refactoriser du code sans changer le comportement | 12K |
| [deep-clean.md](deep-clean.md) | Nettoyage approfondi (fichiers, code mort, refactoring, tests, doc) | 20K |

### 🐛 Debug & Diagnostique
| Prompt | Description | Taille |
|--------|-------------|--------|
| [analyze-error.md](analyze-error.md) | Analyser une erreur avec diagnostic complet | 8K |
| [investigate.md](investigate.md) | Investiguer un comportement étrange sans erreur explicite | 12K |

### ⚡ Performance
| Prompt | Description | Taille |
|--------|-------------|--------|
| [optimize-performance.md](optimize-performance.md) | Optimiser les performances avec profiling et benchmarks | 16K |

### 👀 Revue & Qualité
| Prompt | Description | Taille |
|--------|-------------|--------|
| [code-review.md](code-review.md) | Code review complète et structurée | 12K |

### 📖 Documentation
| Prompt | Description | Taille |
|--------|-------------|--------|
| [explain-code.md](explain-code.md) | Expliquer une partie du code en détail | 12K |
| [update-docs.md](update-docs.md) | Mettre à jour la documentation (README, docs/, GoDoc, CHANGELOG) | 16K |
| [generate-examples.md](generate-examples.md) | Générer des exemples RETE (.constraint, .facts, documentation) | 16K |

### ✓ Validation RETE
| Prompt | Description | Taille |
|--------|-------------|--------|
| [validate-network.md](validate-network.md) | Valider un réseau RETE (structure, propagation, résultats) | 16K |

### 🔄 Migration & Maintenance
| Prompt | Description | Taille |
|--------|-------------|--------|
| [migrate.md](migrate.md) | Migrer version Go, dépendances, ou adapter à changements d'API | 14K |

### 📊 Analyse & Statistiques
| Prompt | Description | Taille |
|--------|-------------|--------|
| [stats-code.md](stats-code.md) | Générer statistiques du code (lignes, complexité, fichiers volumineux) | 10K |

### 📄 Conformité & Licence
| Prompt | Description | Taille |
|--------|-------------|--------|
| [verify-license-compliance.md](verify-license-compliance.md) | Vérifier conformité complète de licence (en-têtes, dépendances, documentation) | 18K |

---

## 📖 Documentation Générale

| Fichier | Description | Taille |
|---------|-------------|--------|
| [README.md](README.md) | Documentation complète du système de prompts | 8K |
| [QUICK_REFERENCE.md](QUICK_REFERENCE.md) | Référence rapide avec exemples | 8K |
| [INDEX.md](INDEX.md) | Ce fichier - navigation globale | 4K |

---

## 🚀 Utilisation Rapide

### Format Standard
```
Utilise le prompt "[nom]"
```

### Exemples Concrets

```
# Tests
Relance moi l'ensemble des tests, dont le runner universel

# Debug
Utilise le prompt "debug-test" pour TestIncrementalPropagation

# Feature
Utilise le prompt "add-feature" pour ajouter support opérateurs chaînes

# Modification
Utilise le prompt "modify-behavior" pour modifier evaluateJoinConditions

# Nettoyage
Utilise le prompt "deep-clean" pour nettoyer le projet

# Review
Utilise le prompt "code-review" pour rete/node_join.go

# Explication
Utilise le prompt "explain-code" pour evaluateJoinConditions (niveau intermédiaire)

# Validation RETE
Utilise le prompt "validate-network" pour beta_coverage_tests/join_complex.constraint
```

---

## 🔍 Recherche par Mot-Clé

| Mot-clé | Prompts Associés |
**Mot-clé** | **Prompts Associés** |
|---------|------------------|
| **test** | run-tests, add-test, debug-test, validate-network |
| **erreur** | analyze-error, debug-test, fix-bug |
| **bug** | fix-bug, debug-test, analyze-error |
| **code** | explain-code, code-review, add-feature, modify-behavior, refactor, deep-clean |
| **RETE** | validate-network, explain-code, generate-examples |
| **debug** | debug-test, analyze-error, investigate |
| **feature** | add-feature, modify-behavior |
| **modification** | modify-behavior, refactor |
| **nettoyage** | deep-clean, refactor |
| **review** | code-review |
| **comprendre** | explain-code, investigate |
| **performance** | optimize-performance, investigate |
| **documentation** | update-docs, generate-examples, explain-code |
| **migration** | migrate |
| **refactoring** | refactor, deep-clean |
| **statistiques** | stats-code |
| **analyse** | stats-code, investigate, analyze-error |
| **métriques** | stats-code, optimize-performance |

---

## 📊 Statistiques

- **Total prompts** : 18
- **Taille totale** : ~260 KB
- **Catégories** : 10
- **Documentation** : 3 fichiers
- **Niveau de détail** : Débutant à Expert

---

## 🎓 Parcours Recommandés

### 👶 Nouveau sur le Projet
1. [explain-code.md](explain-code.md) - Comprendre l'architecture
2. [run-tests.md](run-tests.md) - Valider l'environnement
3. [validate-network.md](validate-network.md) - Explorer RETE
4. [generate-examples.md](generate-examples.md) - Créer exemples pour apprendre

### 👨‍💻 Développeur
1. [add-feature.md](add-feature.md) - Implémenter nouvelle feature
2. [modify-behavior.md](modify-behavior.md) - Modifier feature existante
3. [fix-bug.md](fix-bug.md) - Corriger un bug
4. [refactor.md](refactor.md) - Améliorer la structure du code
5. [add-test.md](add-test.md) - Ajouter des tests
6. [run-tests.md](run-tests.md) - Tester
7. [code-review.md](code-review.md) - Valider
8. [deep-clean.md](deep-clean.md) - Nettoyer régulièrement

### 🐛 Debugger
1. [analyze-error.md](analyze-error.md) - Comprendre l'erreur
2. [debug-test.md](debug-test.md) - Corriger un test
3. [fix-bug.md](fix-bug.md) - Corriger un bug
4. [investigate.md](investigate.md) - Investiguer comportement étrange
5. [run-tests.md](run-tests.md) - Valider

### ⚡ Performance Engineer
1. [optimize-performance.md](optimize-performance.md) - Optimiser
2. [analyze-error.md](analyze-error.md) - Diagnostiquer
3. [investigate.md](investigate.md) - Profiler et analyser

### 🏗️ Architecte RETE
1. [validate-network.md](validate-network.md) - Tester réseau
2. [explain-code.md](explain-code.md) - Comprendre implémentation
3. [generate-examples.md](generate-examples.md) - Créer cas de test
4. [add-feature.md](add-feature.md) - Étendre

### 📝 Documentation Writer
1. [update-docs.md](update-docs.md) - Mettre à jour documentation
2. [generate-examples.md](generate-examples.md) - Créer exemples
3. [explain-code.md](explain-code.md) - Documenter le code

### 🔧 Mainteneur
1. [migrate.md](migrate.md) - Migrer versions/dépendances
2. [deep-clean.md](deep-clean.md) - Nettoyer le projet
3. [refactor.md](refactor.md) - Améliorer la qualité
4. [update-docs.md](update-docs.md) - Maintenir documentation

### 📊 Analyste Qualité
1. [stats-code.md](stats-code.md) - Analyser statistiques du code
2. [code-review.md](code-review.md) - Reviewer la qualité
3. [optimize-performance.md](optimize-performance.md) - Optimiser
4. [refactor.md](refactor.md) - Refactoriser

---

## 🔗 Liens Utiles

- [Makefile du projet](../../Makefile)
- [Tests](../../test/)
- [Documentation RETE](../../docs/)
- [Code source](../../rete/)

---

## 📝 Notes

- Les prompts sont conçus pour être réutilisables et adaptables
- Chaque prompt suit une structure standard : Contexte → Objectif → Instructions → Critères
- N'hésitez pas à combiner plusieurs prompts pour des tâches complexes
- Vous pouvez adapter les prompts à vos besoins spécifiques

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Mainteneur** : Équipe TSD