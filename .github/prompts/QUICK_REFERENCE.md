# ⚡ Référence Rapide - Prompts TSD

## 🎯 Utilisation Express

**Format** : `Utilise le prompt "[nom]"` ou `Applique [nom].md`

## ⚠️ RÈGLES STRICTES APPLIQUÉES

### 🚫 Code Golang
- ❌ **AUCUN HARDCODING** (valeurs, magic numbers, configs)
- ✅ **CODE GÉNÉRIQUE** (paramètres, constantes, interfaces)
- ✅ **BONNES PRATIQUES GO** (Effective Go, go vet, golangci-lint)

### 🚫 Tests RETE
- ❌ **AUCUNE SIMULATION** de résultats
- ✅ **EXTRACTION OBLIGATOIRE** depuis le réseau RETE réel
- ✅ **VALIDATION RÉELLE** (TerminalNodes, mémoires)

**Prompts concernés** :
- Code Go : `add-feature`, `debug-test`, `analyze-error`
- Tests RETE : `debug-test`, `validate-network`, `analyze-error`

---

## 🧪 TESTS

### 🚀 Lancer Tous les Tests
```
Utilise le prompt "run-tests"
```
**Quand** : Validation complète du système  
**Résultat** : Tests unitaires + runner universel + rapport

### 🧪 Ajouter des Tests
```
Utilise le prompt "add-test"
```
**Quand** : Créer nouveaux tests (unitaires, intégration, RETE)  
**Résultat** : Tests ajoutés + couverture améliorée + validation

### 🐛 Débugger un Test
```
Utilise le prompt "debug-test" pour [TestNom]
```
**Quand** : Un test échoue  
**Résultat** : Analyse + correction + validation

---

## 🔧 DÉVELOPPEMENT

### ✨ Ajouter une Fonctionnalité
```
Utilise le prompt "add-feature" pour [description]
```
**Quand** : Nouvelle feature à implémenter  
**Résultat** : Code + tests + documentation

### 🔄 Modifier un Comportement
```
Utilise le prompt "modify-behavior"
```
**Quand** : Modifier une fonctionnalité existante  
**Résultat** : Comportement modifié + tests à jour + doc à jour + prompts mis à jour

### 🐛 Corriger un Bug
```
Utilise le prompt "fix-bug"
```
**Quand** : Bug identifié à corriger  
**Résultat** : Bug corrigé + tests de régression + documentation

### 🔄 Refactoriser du Code
```
Utilise le prompt "refactor"
```
**Quand** : Améliorer structure sans changer comportement  
**Résultat** : Code plus lisible + complexité réduite + tests passent

### 🧹 Nettoyage Approfondi
```
Utilise le prompt "deep-clean"
```
**Quand** : Nettoyage de printemps du code  
**Résultat** : Code propre + fichiers nettoyés + tests améliorés + doc à jour

---

## 🐛 DEBUG

### 🔍 Analyser une Erreur
```
Utilise le prompt "analyze-error"
[Copier l'erreur complète]
```
**Quand** : Erreur incomprise  
**Résultat** : Diagnostic + cause racine + solution

### 🔍 Investiguer un Comportement
```
Utilise le prompt "investigate"
```
**Quand** : Comportement étrange sans erreur explicite  
**Résultat** : Investigation + cause identifiée + rapport

---

## ⚡ PERFORMANCE

### ⚡ Optimiser les Performances
```
Utilise le prompt "optimize-performance"
```
**Quand** : Améliorer performance d'une fonction/module  
**Résultat** : Code optimisé + benchmarks + gains documentés

---

## 👀 REVIEW

### ✅ Code Review
```
Utilise le prompt "code-review"
Fichiers : [liste]
```
**Quand** : Avant merge ou validation  
**Résultat** : Analyse qualité + suggestions + verdict

---

## 📖 DOCUMENTATION

### 📚 Expliquer du Code
```
Utilise le prompt "explain-code"
Fichier : [chemin]
Fonction : [nom]
Niveau : [débutant/intermédiaire/expert]
```
**Quand** : Code incompris  
**Résultat** : Explication détaillée + exemples + contexte

### 📝 Mettre à Jour Documentation
```
Utilise le prompt "update-docs"
```
**Quand** : Documentation obsolète ou incomplète  
**Résultat** : README + CHANGELOG + docs/ + GoDoc à jour

### 🎯 Générer des Exemples RETE
```
Utilise le prompt "generate-examples"
```
**Quand** : Créer exemples .constraint et .facts  
**Résultat** : Fichiers exemples + documentation + cas de test

---

## ✓ VALIDATION

### 🎯 Valider un Réseau RETE
```
Utilise le prompt "validate-network"
Contrainte : [fichier.constraint]
Faits : [fichier.facts]
```
**Quand** : Nouveau réseau ou modification RETE  
**Résultat** : Validation structure + propagation + résultats

---

## 🔄 MIGRATION

### 🔄 Migrer Version / Dépendances
```
Utilise le prompt "migrate"
```
**Quand** : Mettre à jour Go, dépendances, ou adapter API  
**Résultat** : Migration complète + tests + documentation + guide

---

## 📊 ANALYSE

### 📊 Statistiques du Code
```
Utilise le prompt "stats-code"
```
**Quand** : Analyser volume, complexité, identifier refactoring  
**Résultat** : Rapport complet (lignes/module, fichiers volumineux, fonctions longues)

---

## 📋 COMMANDES MAKE UTILES

```bash
make test              # Tests unitaires
make test-coverage     # Tests avec couverture
make test-integration  # Tests d'intégration
make rete-unified      # Runner universel (58 tests)
make validate          # Validation complète
make lint              # Analyse statique
make format            # Formatage code
```

---

## 💡 EXEMPLES CONCRETS

### Scénario 1 : Test Échoue
```
Le test TestIncrementalPropagation échoue avec :
"variable non liée: p (variables disponibles: [u o])"

Utilise le prompt "debug-test" pour corriger ce problème.
```

### Scénario 2 : Ajouter Opérateur
```
Je veux ajouter le support de l'opérateur "contains" pour les chaînes :
{p: Person} / p.name contains "Alice" ==> action(p)

Utilise le prompt "add-feature".
```

### Scénario 3 : Comprendre Fonction
```
Je ne comprends pas evaluateJoinConditions dans rete/node_join.go.

Utilise le prompt "explain-code" niveau intermédiaire.
```

### Scénario 4 : Valider Nouveau Réseau
```
Nouveau réseau : beta_coverage_tests/join_complex.constraint
Attendu : 3 TypeNodes, jointure 3-way, 2 tokens terminaux

Utilise le prompt "validate-network".
```

### Scénario 6 : Modifier un Comportement
```
Je veux modifier le comportement de evaluateJoinConditions pour gérer
intelligemment les variables manquantes au lieu de générer une erreur.

Utilise le prompt "modify-behavior".
```

### Scénario 7 : Nettoyage Approfondi
```
Le projet a accumulé du code au fil du temps : fichiers inutilisés,
code mort, duplication, tests obsolètes.

Utilise le prompt "deep-clean" pour faire un grand nettoyage.
```

### Scénario 5 : Review Avant Merge
```
Fichiers modifiés :
- rete/node_join.go (+150 -30)
- rete/node_join_test.go (+50 -0)

Utilise le prompt "code-review".
```

---

## 🔥 RACCOURCIS

| Besoin | Prompt | Format Court |
|--------|--------|--------------|
| Tester tout | `run-tests` | `Relance les tests` |
| Ajouter tests | `add-test` | `Ajoute tests pour X` |
| Fix test | `debug-test` | `Debug TestNom` |
| Nouvelle feature | `add-feature` | `Ajoute support de X` |
| Modifier feature | `modify-behavior` | `Modifie comportement de X` |
| Corriger bug | `fix-bug` | `Corrige le bug X` |
| Refactoriser | `refactor` | `Refactorise la fonction X` |
| Nettoyage | `deep-clean` | `Nettoie le projet` |
| Erreur | `analyze-error` | `Analyse cette erreur: ...` |
| Investigation | `investigate` | `Investigue ce comportement: ...` |
| Performance | `optimize-performance` | `Optimise X` |
| Review | `code-review` | `Review de X.go` |
| Expliquer | `explain-code` | `Explique la fonction X` |
| Mettre à jour docs | `update-docs` | `Mets à jour la doc` |
| Générer exemples | `generate-examples` | `Crée exemple pour X` |
| Valider RETE | `validate-network` | `Valide le réseau X.constraint` |
| Migrer | `migrate` | `Migre vers Go 1.21` |
| Stats code | `stats-code` | `Statistiques du code` |

---

## 🎓 WORKFLOW RECOMMANDÉ

### Nouveau Développeur
1. `explain-code` → Comprendre le projet
2. `run-tests` → Valider environnement
3. `validate-network` → Explorer RETE
4. `generate-examples` → Créer exemples pour apprendre

### Feature Development
1. `add-feature` → Implémenter nouvelle feature
2. `add-test` → Ajouter tests
3. `run-tests` → Tester
4. `code-review` → Valider
5. `update-docs` → Documenter
6. Merge ✅

### Bug Fixing
1. `analyze-error` → Comprendre (si erreur explicite)
2. `investigate` → Investiguer (si comportement étrange)
3. `fix-bug` → Corriger avec méthodologie
4. `debug-test` → Corriger tests
5. `run-tests` → Valider
6. `code-review` → Review

### RETE Development
1. `generate-examples` → Créer `.constraint` et `.facts`
2. `validate-network` → Tester réseau
3. `run-tests` → Intégrer
4. `code-review` → Valider
5. `update-docs` → Documenter

### Maintenance Régulière
1. `refactor` → Améliorer structure code
2. `deep-clean` → Nettoyer code
3. `run-tests` → Valider
4. `update-docs` → Maintenir documentation
5. `code-review` → Vérifier qualité
6. Commit ✅

### Performance Optimization
1. `optimize-performance` → Profiler et optimiser
2. `investigate` → Analyser comportement
3. `run-tests` → Valider (benchmarks)
4. `code-review` → Review

### Migration / Upgrade
1. `migrate` → Planifier et exécuter migration
2. `run-tests` → Valider tous les tests
3. `update-docs` → Documenter changements
4. `code-review` → Review complète

### Analyse Qualité
1. `stats-code` → Analyser statistiques du code
2. `code-review` → Reviewer qualité
3. `refactor` → Améliorer code identifié
4. `optimize-performance` → Optimiser si nécessaire

---

## 📞 AIDE

**Documentation complète** : `.github/prompts/README.md`

**Structure des prompts** :
- Contexte : Informations projet
- Objectif : Ce qu'on veut accomplir
- Instructions : Étapes détaillées
- Critères de succès : Validation
- Exemple : Utilisation concrète

**Créer un nouveau prompt** : Voir `.github/prompts/README.md#ajouter-un-nouveau-prompt`

---

**Version** : 2.0 (18 prompts disponibles)  
**Dernière mise à jour** : Novembre 2025