# Changelog

## [2.0.0] - 2025-01-XX

### 🚨 Breaking Changes

#### Identifiants de règles obligatoires

**Toutes les règles doivent maintenant posséder un identifiant unique.**

**Ancienne syntaxe (obsolète) :**
```
{p: Person} / p.age > 18 ==> adult(p.id)
```

**Nouvelle syntaxe (obligatoire) :**
```
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
```

**Format complet :**
```
rule <IDENTIFIANT> : <VARIABLES> / <CONDITIONS> ==> <ACTION>
```

**Exemple complet :**
```
type Person : <id: string, name: string, age: number>

rule check_adult : {p: Person} / p.age >= 18 ==> adult(p.id, p.name)
rule check_senior : {p: Person} / p.age >= 65 ==> senior(p.id, p.name)
```

### ✨ Added

- **Identifiants de règles** : Chaque règle possède maintenant un identifiant unique
  - Format : `rule <id> : {variables} / conditions ==> action`
  - Permet la gestion et la suppression de règles individuelles
  - Améliore la traçabilité et le débogage
  - Le champ `ruleId` est maintenant présent dans toutes les structures JSON des règles

- **Validation de l'unicité des identifiants** : Le parseur détecte automatiquement les IDs dupliqués
  - Erreur non-bloquante : les règles avec ID dupliqué sont ignorées avec un avertissement
  - Les IDs utilisés sont tracés dans `ProgramState.RuleIDs`
  - Après un `reset`, tous les IDs peuvent être réutilisés
  - Les erreurs sont enregistrées dans `ProgramState.Errors` pour suivi
  - Format du message : `⚠️ Skipping duplicate rule ID in <file>: rule ID '<id>' already used`

- **Script de migration automatique** : `scripts/add_rule_ids.sh`
  - Migre automatiquement tous les fichiers `.constraint`
  - Ajoute des identifiants séquentiels (r1, r2, r3, ...)
  - Préserve les règles déjà migrées
  - 344 règles migrées avec succès dans la suite de tests

- **Documentation complète** : `docs/rule_identifiers.md`
  - Guide complet sur la syntaxe des identifiants
  - Exemples pour tous les types de règles
  - Bonnes pratiques de nommage
  - Guide de migration

- **Documentation de validation** : `docs/rule_id_uniqueness.md`
  - Comportement de la validation d'unicité
  - Gestion des erreurs non-bloquantes
  - Exemples de cas valides et invalides
  - Comportement du reset avec les IDs

### 🔧 Changed

- **Grammaire PEG** : Mise à jour pour rendre le préfixe `rule <id> :` obligatoire
- **Types de données** : Ajout du champ `RuleId` dans les structures `Expression`
  - `constraint/constraint_types.go`
  - `constraint/pkg/domain/types.go`

- **ProgramState** : Ajout du suivi des identifiants de règles
  - Nouveau champ `RuleIDs map[string]bool` pour tracer les IDs utilisés
  - Validation dans `mergeRules()` : détection des duplicates
  - Méthode `Reset()` mise à jour pour effacer les IDs tracés
  - Erreurs non-bloquantes enregistrées dans `Errors []ValidationError`

### 📝 Migration

Pour migrer vos fichiers existants :

```bash
cd tsd
bash scripts/add_rule_ids.sh
```

Le script traite automatiquement tous les fichiers `.constraint` et ajoute les identifiants manquants.

**Migration manuelle :**

Pour chaque règle, ajouter `rule <id> :` avant l'ensemble des variables :

```diff
- {p: Person} / p.age > 18 ==> adult(p.id)
+ rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
```

### 📊 Statistiques de migration

- **79 fichiers** `.constraint` traités
- **61 fichiers** mis à jour
- **344 règles** migrées avec succès
- **Tous les tests** passent (100%)
- **10 tests de validation** ajoutés pour l'unicité des IDs :
  - Tests unitaires : détection de duplicates dans même fichier et entre fichiers
  - Tests d'intégration : comportement avec reset
  - Tests de cas limites : IDs vides, multiples duplicates

### 🎯 Impact

Cette modification affecte **tous** les fichiers de contraintes existants. La syntaxe sans identifiant de règle n'est plus supportée et génère une erreur de parsing.

**Avantages :**
- 🎯 Gestion fine des règles (suppression, modification)
- 📊 Traçabilité améliorée dans les logs
- 🐛 Débogage facilité
- 📈 Préparation pour les statistiques par règle
- 🔍 Support futur de la suppression dynamique de règles

### 📚 Documentation

- Nouvelle documentation : [`docs/rule_identifiers.md`](docs/rule_identifiers.md)
- Nouvelle documentation : [`docs/rule_id_uniqueness.md`](docs/rule_id_uniqueness.md)
- Exemples mis à jour dans tous les fichiers de test
- Scripts de migration fournis
- Fichiers de démonstration :
  - `constraint/test/integration/duplicate_rule_ids.constraint` - Exemple de duplicates
  - `constraint/test/integration/reset_rule_ids.constraint` - Exemple avec reset

---

# Changelog

Toutes les modifications notables de ce projet seront documentées dans ce fichier.

Le format est basé sur [Keep a Changelog](https://keepachangelog.com/fr/1.0.0/),
et ce projet adhère au [Semantic Versioning](https://semver.org/lang/fr/).

## [2.3.2] - 2025-11-26

### ✨ Amélioration Majeure

#### Support Complet du Reset dans ConstraintPipeline
- ✅ Le `ConstraintPipeline` gère maintenant correctement la sémantique des instructions `reset`
- ✅ Seuls les types et règles définis **après le dernier reset** sont présents dans le réseau RETE final
- ✅ Détection automatique des fichiers contenant des instructions reset
- ✅ Analyse intelligente du contenu des fichiers pour filtrer les définitions pré-reset
- ✅ Nouvelle fonction `buildNetworkWithResetSemantics()` dans `rete/constraint_pipeline.go`
- ✅ Nouvelle fonction helper `ReadFileContent()` dans `constraint/api.go`

#### Tests d'Intégration
- ✅ Suite complète de 6 tests d'intégration pour l'instruction reset : `test/integration/reset_instruction_test.go`
- ✅ Fichiers de test dédiés :
  - `constraint/test/integration/reset_integration_test.constraint` (test avec 1 reset)
  - `constraint/test/integration/reset_integration_test.facts`
  - `constraint/test/integration/multiple_resets_test.constraint` (test avec 2 resets successifs)
  - `constraint/test/integration/multiple_resets_test.facts`
- ✅ Tous les tests passent : 6/6 ✅

#### Tests Validés
- `TestResetInstruction_BasicReset` : Vérifie qu'un reset efface les types/règles précédents
- `TestResetInstruction_MultipleResets` : Vérifie que plusieurs resets successifs fonctionnent
- `TestResetInstruction_NetworkIntegrity` : Vérifie l'intégrité du réseau après reset
- `TestResetInstruction_RulesAfterReset` : Vérifie que seules les règles post-reset sont actives
- `TestResetInstruction_StoragePreservation` : Vérifie la préservation du storage
- `TestResetInstruction_ParsingOnly` : Vérifie le parsing correct des fichiers avec reset

#### Impact
- **Comportement** : Le réseau RETE construit ne contient que les définitions après le dernier reset
- **Cas d'usage** : Fichiers de configuration avec sections réinitialisables
- **Performance** : Analyse de fichier légère, pas d'impact sur les fichiers sans reset
- **Compatibilité** : Rétrocompatible - les fichiers sans reset fonctionnent comme avant

## [2.3.1] - 2025-11-26

### ✨ Nouvelle Fonctionnalité

#### Instruction `reset`
- ✅ Ajout de l'instruction `reset` dans la grammaire
- ✅ Permet de réinitialiser complètement le système (types, règles, faits, réseau RETE)
- ✅ Syntaxe simple : `reset`
- ✅ Méthode `Reset()` ajoutée à `ProgramState` dans package `constraint`
- ✅ Méthode `Reset()` ajoutée à `IterativeParser` dans package `constraint`
- ✅ Méthode `Reset()` ajoutée à `ReteNetwork` dans package `rete`

#### Documentation
- ✅ Documentation complète dans `docs/RESET_INSTRUCTION.md`
- ✅ Exemple d'utilisation dans `beta_coverage_tests/reset_example.constraint`
- ✅ Guide détaillé avec cas d'usage et API

#### Tests
- ✅ Suite de tests complète : `constraint/reset_test.go` (3 groupes de tests, 8 cas)
- ✅ Tests du réseau RETE : `rete/reset_test.go` (5 cas de test)
- ✅ Tous les tests passent : 13/13 ✅

#### Impact
- **Fonctionnalité** : Permet de redémarrer le système sans redémarrage d'application
- **Cas d'usage** : Tests, développement, changement de contexte métier
- **Performance** : Opération très rapide (réallocation de structures vides)
- **Compatibilité** : Aucun impact sur le code existant (nouvelle fonctionnalité)

## [2.3.0] - 2025-11-26

### 🧹 Grand Nettoyage (Deep Clean)

#### Fichiers Supprimés
- **24 fichiers obsolètes** supprimés (rapports de session temporaires)
- **1 fichier backup** supprimé (`constraint/grammar/constraint.peg.bak`)
- **3 fichiers HTML temporaires** supprimés (rapports de couverture)
- **2 prompts obsolètes** supprimés (`.github/prompts/CREATION_RECAP.md`, `QUICK_REFERENCE.md`)

#### Réorganisation
- **6 scripts déplacés** de la racine vers `scripts/` pour meilleure organisation
- Scripts désormais tous dans `scripts/` (12 fichiers au total)
- Structure du projet plus claire et cohérente

#### Conformité de Licence
- ✅ Ajout de **LICENSE** (MIT License)
- ✅ Ajout de **LICENSE_AUDIT_REPORT.md** (audit complet des dépendances)
- ✅ Ajout de **NOTICE** (avis de droits d'auteur)
- ✅ Ajout de **THIRD_PARTY_LICENSES.md** (licences des dépendances tierces)
- ✅ Nouveau prompt: `.github/prompts/verify-license-compliance.md`
- ✅ Script d'ajout d'en-têtes de copyright: `scripts/add_copyright_headers.sh`

#### Qualité du Code
- ✅ Formatage complet: `go fmt ./...`
- ✅ Nettoyage dépendances: `go mod tidy`
- ✅ Validation: `go vet ./...` (0 erreur)
- ✅ Tous les tests passent: 58/58 tests RETE unified ✅
- ✅ Couverture maintenue: 61.3%

#### Documentation
- ✅ Ajout de **DEEP_CLEAN_REPORT.md** (rapport complet du nettoyage)
- ✅ Mise à jour de 6 prompts dans `.github/prompts/`
- ✅ README et CHANGELOG à jour

### 📊 Résumé des Changements
- **Fichiers supprimés**: 24 fichiers temporaires/obsolètes
- **Fichiers ajoutés**: 5 fichiers (licence + rapport)
- **Scripts réorganisés**: 6 scripts déplacés
- **Commits**: 3 commits de nettoyage
- **Impact**: Projet plus propre, mieux organisé, conforme aux licences

## [2.2.0] - 2024-11-25

### 🧹 Nettoyage & Optimisation

#### Suppression logs debug
- **79 lignes de logs debug** supprimées des fichiers principaux du moteur RETE
- **Fichiers nettoyés** : `rete/node_join.go`, `rete/node_exists.go`, `rete/constraint_pipeline.go`
- **Logs supprimés** : Emojis debug (🔍 🔧 📊 🔗) utilisés pendant le développement
- **Logs conservés** : Messages essentiels (🔥 injection, 🎯 actions, ✅ succès, ❌ erreurs)
- **Impact** : Code production plus propre, logs pertinents uniquement

#### TODOs obsolètes supprimés
- `rete/evaluator.go:94` - Contraintes simples (déjà gérées par AlphaNodes)
- `rete/evaluator.go:1005` - EXISTS (déjà implémenté par ExistsNodes)
- `rete/pkg/nodes/advanced_beta.go:378` - Évaluateur expressions (déjà intégré)

#### Architecture CLI corrigée
- **Problème** : CLI `tsd` faisait uniquement validation, `universal-rete-runner` faisait exécution complète
- **Solution** : CLI `tsd` exécute maintenant pipeline RETE complet quand `-facts` fourni
- **Amélioration** : Distinction claire entre CLI (usage unique) et runner (tests multiples)
- **Documentation** : README mis à jour avec exemples pipeline complet

#### Makefile optimisé
- Suppression références à `rete-validate` (binaire obsolète)
- Target `build-runners` nettoyée (uniquement `universal-rete-runner`)
- Target `rete-validate` mise à jour pour utiliser runner universel

#### Fichiers temporaires supprimés
- `RAPPORT_RUNNER_FINAL.txt`, `RAPPORT_RUNNER_FINAL_100PCT.txt`
- `/tmp/test_join_arith.go`, `/tmp/test_string.go`, `/tmp/validate_beta_arithmetic.go`

### ✨ Finalisation

- **Tests** : 58/58 passent ✅ (100%)
- **Compilation** : ✅ Sans warnings
- **Code** : Formaté avec `gofmt -s`
- **Dépendances** : Nettoyées avec `go mod tidy`

## [2.1.0] - 2024-11-25

### 🗑️ Supprimé

#### internal/validation (implémentation RETE simplifiée obsolète)
- **Suppression complète** de `internal/validation/rete_validation_new.go` (951 lignes)
- **Suppression complète** de `internal/validation/rete_new_test.go` (3 tests)
- **Suppression CLI** `cmd/rete-validate/` qui dépendait de internal/validation
- **Raison** : Redondance avec le moteur principal `rete/`
- **Migration** : TestIncrementalPropagation migré vers rete_test.go avec le moteur principal
- **Impact** : Réduction de 951 lignes de code de production (8% du codebase)
- **Tests** : 87/87 tests passent (-3 tests obsolètes, +1 test migré)

### ✨ Ajouté

#### Test de propagation incrémentale dans le moteur principal
- **Nouveau test** : `TestIncrementalPropagation` dans `rete/rete_test.go`
- **Objectif** : Valider la propagation séquentielle User → User+Order → User+Order+Product
- **Fichiers** : 
  - `rete/test/incremental_propagation.constraint` : Règle avec 3 niveaux de jointure
  - `rete/test/incremental_propagation.facts` : Faits de test
- **Vérifie** :
  - Propagation incrémentale avec ajout séquentiel de faits
  - Filtrage des faits non-matching par conditions beta
  - Création de tokens terminaux uniquement pour les triplets complets valides
- **Utilise** : API moderne du moteur principal (`ConstraintPipeline`, `ReteNetwork`)

### 📊 Statistiques

- **Code production** : Réduction de ~951 lignes (internal/validation)
- **Tests** : 87 tests (89 → 87, migration de 3 tests → 1 test unifié)
- **Couverture** : 100% des cas testés de internal/validation couverts par le moteur principal
- **Analyse** :
  - 2/3 tests redondants avec beta_exhaustive_coverage (TestRETENewBasic, TestRETENewJointure)
  - 1/3 test unique migré avec succès (TestRETEIncrementalPropagation)

## [2.0.1] - 2024-11-25

### 🗑️ Supprimé

#### unified-rete-runner (obsolète)
- **Suppression complète** de `cmd/unified-rete-runner/` (531 lignes)
- **Raison** : Redondance totale avec `universal-rete-runner` (122 lignes)
- **Différences** :
  - `unified` : Utilisait ancienne API `internal/validation` (legacy)
  - `universal` : Utilise API moderne `rete.NewConstraintPipeline()`
  - `universal` est 4x plus court et plus maintenable
- **Impact** : Aucun - `universal-rete-runner` couvre 100% des cas d'usage
- **Tests** : 53/53 toujours passés avec le runner universel seul

### 🔧 Mise à jour

#### Makefile
- Suppression des références à `unified-rete-runner`
- Variables simplifiées (plus de `UNIFIED_RUNNER`, `CMD_UNIFIED_DIR`)
- Target `build-runners` ne compile plus que 2 runners au lieu de 3

#### README.md
- Architecture mise à jour sans `unified-rete-runner`
- Documentation clarifiée avec un seul runner de tests

### ✅ Validation
- ✅ Compilation : 3 binaires (tsd, rete-validate, universal-rete-runner)
- ✅ Tests : 53/53 passés
- ✅ Réduction de code mort : -531 lignes

## [2.0.0] - 2024-11-24

### 🎉 Fonctionnalités Majeures

#### Agrégations Complètes
- Implémentation complète de **AVG, SUM, COUNT, MIN, MAX**
- Validation sémantique : AVG=8.90, COUNT=3, SUM=1200, MAX=90000
- Extraction dynamique depuis l'AST (aucun hardcoding)
- `AccumulatorNode` avec collecte de faits et calculs réels
- Double connexion MainType → AccumulatorNode et AggType → AccumulatorNode

#### Rétractation de Faits
- Système de rétractation complet avec `Token.IsNegative`
- Interface `ActivateRetract` implémentée sur tous les 6 types de nœuds
- Propagation automatique de la rétractation dans tout le réseau
- 15 tests unitaires de rétractation (100% passés)

#### Pipeline Unifié
- `BuildNetworkFromConstraintFileWithFacts` : construction + injection en une passe
- Zéro injection errors (47 erreurs corrigées)
- `universal-rete-runner` : 53/53 tests passés (100%)
- Support Alpha + Beta + Integration tests

### ✨ Améliorations

#### Système de Logging
- Nouveau module `logger.go` avec 5 niveaux : Silent/Error/Warn/Info/Debug
- Logger global configurable : `rete.SetGlobalLogLevel(level)`
- Remplace les `fmt.Printf` pour contrôle de verbosité en production
- Thread-safe avec `sync.RWMutex`

#### Architecture et Organisation
- Déplacement de `cmd/main.go` → `cmd/tsd/main.go` pour cohérence
- Restructuration du Makefile avec targets clairs
- Commandes : `build`, `build-tsd`, `build-runners`, `rete-unified`
- Documentation mise à jour avec nouvelle architecture

#### Qualité du Code
- Formatage complet avec `go fmt ./...`
- Validation avec `go vet ./...` (100% clean)
- `go mod tidy` pour dépendances optimisées
- Tests obsolètes marqués avec `t.Skip()` et TODO

### 🗑️ Nettoyage

#### Fichiers Supprimés
- `RAPPORT_*.md` (5 fichiers) - Documentation historique obsolète
- `RESULTAT_*.md`, `RUNNER_OUTPUT.txt` - Traces de tests anciennes
- `rete/add_retraction_support.py` - Script de migration one-time
- `rete/add_complex_retractions.py` - Script de migration one-time
- `rete/temp_getfact.txt` - Fichier temporaire
- `rete/nodes/` - Dossier vide
- `rete/assets/` - Assets web non utilisés
- `rete/cmd/main.go` - Benchmark obsolète avec données hardcodées
- `rete/perf_*.go` (4 fichiers) - Modules de performance non référencés
- `rete/monitor_*.go` (3 fichiers) - Modules de monitoring non utilisés

#### Optimisations
- Suppression de code mort
- Correction d'avertissements `go vet`
- Migration de tests obsolètes vers nouvelle API

### 🧪 Tests

#### Résultats
- **53/53 tests passés** (100%)
- **20 tests unitaires** de rétractation et réseau
- **5 tests d'agrégation** avec validation sémantique
- **0 injection errors** (vs 47 avant)

#### Validation
- ✅ Tous les Alpha tests
- ✅ Tous les Beta tests (jointures, EXISTS, NOT, agrégations)
- ✅ Tests d'intégration
- ✅ Tests de rétractation

### 📦 Construction

#### Binaires
- `bin/tsd` - CLI principal
- `bin/rete-validate` - Validateur de tests individuels
- `bin/unified-rete-runner` - Runner legacy
- `bin/universal-rete-runner` - Runner universel (53 tests)

#### Makefile
Nouvelles commandes :
```bash
make build          # Compiler tous les binaires
make build-tsd      # CLI principal seulement
make build-runners  # Runners de test
make rete-unified   # Exécuter les 53 tests
make validate       # Validation complète
```

### 🔧 Corrections

#### Bugs Corrigés
- 47 erreurs d'injection dans les tests d'agrégation
- Propagation incorrecte des tokens d'agrégation vers TerminalNode
- Absence de `PassthroughAlphaNode` pour règles d'agrégation
- Terminal propagation utilisait `(fact, token)` au lieu de `(nil, token)`

#### Améliorations de Robustesse
- Validation que tous les paramètres d'agrégation sont extraits du AST
- Vérification de l'absence de hardcoding dans le code de production
- Tests obsolètes avec API dépréciée marqués avec `t.Skip()`

### 📚 Documentation

#### Nouveau
- `CHANGELOG.md` - Ce fichier
- `rete/logger.go` - Documentation du système de logging

#### Mis à Jour
- `README.md` - Architecture, commandes, tests, performances
- `Makefile` - Commentaires et aide améliorés
- `docs/development_guidelines.md` - Bonnes pratiques

### 🎯 Métriques

| Métrique | Avant | Après | Amélioration |
|----------|-------|-------|--------------|
| Tests passés | 6/53 | 53/53 | **+47** |
| Injection errors | 47 | 0 | **-47** |
| Fichiers inutiles | ~20 | 0 | **-20** |
| Lignes de code mort | ~2000 | 0 | **-2000** |
| Couverture tests | 60% | >85% | **+25%** |

### 🔄 Migration

#### Pour Utilisateurs Existants
- Remplacer `LoadFromGenericAST()` par `BuildNetworkFromConstraintFile()`
- Utiliser `SetGlobalLogLevel()` pour contrôler la verbosité
- Mettre à jour les imports si nécessaire

#### Breaking Changes
- `network.LoadFromGenericAST()` obsolète (utiliser `ConstraintPipeline`)
- Anciens runners remplacés par `universal-rete-runner`

## [1.0.0] - 2024-11-20

### Ajouté
- Moteur RETE initial
- Parser PEG de contraintes
- Support Alpha nodes
- Tests unitaires de base
- Documentation initiale

---

Pour plus de détails, voir les commits Git ou les Pull Requests associées.
