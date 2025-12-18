# 🎯 Plan Global - Finalisation Xuples E2E

> **Date**: 2025-12-18  
> **Objectif**: Éliminer toutes les limitations actuelles et rendre l'intégration des xuples 100% automatique

---

## 📋 Vue d'Ensemble

Ce plan d'action vise à finaliser l'exécution correcte des pipelines TSD intégrant la gestion des xuples, en éliminant toutes les limitations identifiées dans `XUPLES_E2E_AUTOMATIC.md`.

### Principes Directeurs

1. **Xuples fait partie intégrante de TSD** - Ce n'est PAS un module optionnel
2. **Zéro configuration manuelle** - Les tests soumettent uniquement le fichier TSD
3. **Architecture simplifiée** - Pipeline unique et cohérent
4. **Support parser complet** - Faits inline dans les actions avec références aux variables
5. **Automatisation totale** - Création des xuple-spaces, xuples, et exécution des actions

---

## 🏗️ Architecture Cible

### Approche Retenue: Package API Pipeline Complet

**Avant** (architecture actuelle avec limitations):
```
┌─────────────────────────────────────────┐
│ Test E2E                                │
│                                         │
│ 1. Créer réseau RETE                   │
│ 2. Configurer factory manuellement      │ ← MANUEL
│ 3. Appeler IngestFile()                 │
│ 4. Créer xuples manuellement            │ ← MANUEL (parser limité)
│ 5. Vérifier résultats                   │
└─────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ rete.ConstraintPipeline                 │
│                                         │
│ - Parse TSD                             │
│ - Construit réseau RETE                 │
│ - Appelle factory (si configurée)       │ ← OPTIONNEL
│ - Ne connaît pas xuples (cycle import)  │
└─────────────────────────────────────────┘
```

**Après** (architecture cible simplifiée):
```
┌─────────────────────────────────────────┐
│ Test E2E                                │
│                                         │
│ 1. pipeline.IngestFile("fichier.tsd")  │ ← UN SEUL APPEL
│ 2. Vérifier résultats                   │
└─────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ api.Pipeline (NOUVEAU)                  │
│                                         │
│ - Importe rete + xuples + constraint    │
│ - Configuration automatique             │
│ - Point d'entrée unique et simple       │
│ - Gestion complète du cycle de vie      │
└─────────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────┐
│ rete.ConstraintPipeline                 │
│                                         │
│ - Parse TSD (support faits inline)      │
│ - Construit réseau RETE                 │
│ - Crée xuple-spaces dès parsing         │
│ - Enregistre actions Xuple auto         │
└─────────────────────────────────────────┘
```

---

## 📦 Structure des Prompts

Chaque prompt est conçu pour être exécuté dans une session unique (contexte 128k max).

### Prompt 01 - Parser TSD (Support Faits Inline)
**Fichier**: `01-parser-faits-inline.md`  
**Objectif**: Étendre le parser TSD pour supporter complètement les faits inline dans les actions

**Livrables**:
- Support syntaxe `Xuple("space", Alert(level: "CRITICAL", message: msg))`
- Support multi-ligne dans les actions
- Références aux champs des faits déclencheurs (ex: `s.sensorId`, `s.temperature`)
- Tests parser pour valider toutes les variantes syntaxiques

**Impact**:
- `constraint/parser.go` (ou fichiers PEG)
- `constraint/ast.go` (si besoin nouveaux nœuds AST)
- Tests unitaires parser
- Documentation syntaxe TSD

---

### Prompt 02 - Package API Pipeline
**Fichier**: `02-package-api-pipeline.md`  
**Objectif**: Créer le package `api` centralisant le pipeline complet

**Livrables**:
- Nouveau package `tsd/api/`
- `api.Pipeline` avec méthode `IngestFile(filename string) (*Result, error)`
- Configuration automatique rete + xuples
- Gestion automatique des xuple-spaces
- Documentation GoDoc complète

**Impact**:
- `api/pipeline.go` (nouveau)
- `api/result.go` (nouveau)
- `api/doc.go` (nouveau)
- Tests unitaires API
- README.md mise à jour

---

### Prompt 03 - Création Automatique Xuple-Spaces
**Fichier**: `03-creation-auto-xuple-spaces.md`  
**Objectif**: Créer les xuple-spaces dès l'étape de parsing (réaction aux commandes de définition)

**Livrables**:
- Hook dans le pipeline pour créer xuple-spaces pendant parsing
- Élimination de la factory (intégration directe via package API)
- Validation des configurations xuple-spaces
- Tests de création automatique

**Impact**:
- `rete/constraint_pipeline.go` (modification flow)
- `api/pipeline.go` (intégration directe xuples)
- Suppression de `XupleSpaceFactoryFunc` (obsolète)
- Tests pipeline

---

### Prompt 04 - Actions Xuple Automatiques
**Fichier**: `04-actions-xuple-automatiques.md`  
**Objectif**: Automatiser complètement l'exécution des actions Xuple dans les règles

**Livrables**:
- Enregistrement automatique de l'action `Xuple` au démarrage
- Exécution automatique lors du déclenchement des règles
- Support complet des faits inline (utilisant le parser amélioré)
- Gestion des références aux faits déclencheurs (triggeringFacts)
- Tests validant la création automatique des xuples par les règles

**Impact**:
- `rete/actions.go` (si modifications nécessaires)
- `rete/terminal_node.go` (exécution actions)
- `api/pipeline.go` (enregistrement automatique)
- Tests E2E avec règles créant des xuples
- Documentation actions TSD

---

### Prompt 05 - Migration Tests E2E
**Fichier**: `05-migration-tests-e2e.md`  
**Objectif**: Migrer tous les tests E2E pour utiliser le nouveau package API

**Livrables**:
- Migration `tests/e2e/xuples_e2e_test.go` vers `api.Pipeline`
- Suppression de toute configuration manuelle
- Suppression de toute création manuelle de xuples/xuple-spaces
- Validation que les xuples sont créés par les règles
- Tests additionnels pour cas limites

**Impact**:
- `tests/e2e/xuples_e2e_test.go` (simplification majeure)
- Autres tests E2E si applicable
- Suppression du code de workaround
- Rapport E2E mis à jour

---

### Prompt 06 - Refactoring et Nettoyage
**Fichier**: `06-refactoring-nettoyage.md`  
**Objectif**: Nettoyer le code obsolète et refactorer pour cohérence

**Livrables**:
- Suppression du pattern factory (obsolète avec package API)
- Suppression des méthodes `SetXupleSpaceFactory`, `GetXupleSpaceFactory`
- Nettoyage des imports et dépendances
- Vérification qu'aucun code mort ne subsiste
- Mise à jour de la documentation

**Impact**:
- `rete/network.go` (suppression factory)
- `rete/constraint_pipeline.go` (simplification)
- Documentation architecture
- XUPLES_E2E_AUTOMATIC.md (mise à jour ou archivage)

---

### Prompt 07 - Tests Complets et Documentation
**Fichier**: `07-tests-documentation.md`  
**Objectif**: Compléter la couverture de tests et finaliser la documentation

**Livrables**:
- Tests unitaires pour toutes les nouvelles fonctionnalités
- Tests d'intégration validant le flow complet
- Couverture > 80% pour les nouveaux packages
- Documentation utilisateur (README, guides)
- Documentation développeur (architecture, GoDoc)
- Exemples TSD complets

**Impact**:
- Tests unitaires (`api/`, `rete/`, `constraint/`)
- Tests d'intégration
- `README.md`
- `docs/xuples-guide.md` (nouveau)
- `docs/api-pipeline.md` (nouveau)
- `examples/xuples/` (fichiers .tsd exemples)

---

## 🔄 Ordre d'Exécution

Les prompts doivent être exécutés dans l'ordre suivant (dépendances):

```
01 Parser ──────────┐
                    ├──> 04 Actions ──┐
02 API Pipeline ────┤                 ├──> 05 Migration Tests
                    ├──> 03 Spaces ───┤
                    │                 │
                    └─────────────────┴──> 06 Refactoring ──> 07 Tests & Doc
```

**Explications**:
- **Prompt 01** (Parser) est indépendant et peut être fait en premier
- **Prompt 02** (API Pipeline) est le fondement de la nouvelle architecture
- **Prompt 03** (Xuple-Spaces) dépend de l'API Pipeline
- **Prompt 04** (Actions) dépend du Parser et de l'API Pipeline
- **Prompt 05** (Migration) dépend de tous les précédents
- **Prompt 06** (Refactoring) nettoie après la migration
- **Prompt 07** (Tests & Doc) finalise tout

---

## ✅ Critères de Succès

### Fonctionnels

- [ ] Un test E2E se résume à: `pipeline.IngestFile("test.tsd")` + vérification résultats
- [ ] Aucune configuration manuelle requise
- [ ] Aucune création manuelle de xuples ou xuple-spaces
- [ ] Le parser supporte `Xuple("space", Fact(field: var.subfield))`
- [ ] Les règles créent automatiquement les xuples lors de leur déclenchement
- [ ] Les xuple-spaces sont créés dès le parsing de leur définition

### Techniques

- [ ] Aucun cycle d'importation
- [ ] Couverture de tests > 80%
- [ ] Tous les tests passent (`make test`)
- [ ] Validation complète passe (`make validate`)
- [ ] Documentation à jour et complète
- [ ] Code conforme aux standards TSD (voir `.github/prompts/common.md`)

### Architecture

- [ ] Package `api` centralise le pipeline complet
- [ ] `rete` reste indépendant de `xuples` (mais utilisable via `api`)
- [ ] Xuples intégré de manière transparente (pas optionnel)
- [ ] Pattern factory supprimé (obsolète)
- [ ] Code obsolète nettoyé

---

## 📊 Impact Estimé

### Fichiers Nouveaux

```
api/
├── pipeline.go      (pipeline complet, point d'entrée)
├── result.go        (résultats d'ingestion)
├── config.go        (configuration optionnelle)
├── doc.go           (documentation package)
└── pipeline_test.go (tests unitaires)

docs/
├── xuples-guide.md     (guide utilisateur xuples)
├── api-pipeline.md     (doc API pipeline)
└── migration-guide.md  (migration anciens codes)

examples/xuples/
├── monitoring.tsd      (exemple monitoring capteurs)
├── workflow.tsd        (exemple workflow avec xuples)
└── README.md           (explication exemples)
```

### Fichiers Modifiés

```
constraint/
├── parser.go (ou fichiers PEG) - Support faits inline
├── ast.go                      - Nouveaux nœuds AST si besoin
└── parser_test.go              - Tests nouvelles syntaxes

rete/
├── constraint_pipeline.go      - Intégration xuple-spaces pendant parsing
├── network.go                  - Suppression factory, simplification
├── actions.go                  - Actions Xuple (si modifs)
└── terminal_node.go            - Exécution actions (si modifs)

tests/e2e/
└── xuples_e2e_test.go          - Simplification majeure

README.md                        - Ajout section xuples et API
XUPLES_E2E_AUTOMATIC.md         - Mise à jour ou archivage
```

### Fichiers Supprimés (Code Obsolète)

- Méthodes factory dans `rete/network.go`
- Code de workaround dans tests
- Documentation obsolète

---

## 🎯 Résultat Final

### Avant (Workflow Actuel - Complexe)

```go
// Test E2E - 9 étapes dont 7 manuelles
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)
pipeline := rete.NewConstraintPipeline()

// Configurer factory (complexe, répétitif)
network.SetXupleSpaceFactory(func(net *rete.ReteNetwork, defs []interface{}) error {
    xupleManager := xuples.NewXupleManager()
    // ... 50 lignes de configuration manuelle ...
    return nil
})

// Ingérer fichier
network, metrics, err := pipeline.IngestFile("test.tsd", network, storage)

// Créer xuples manuellement (workaround parser)
xupleManager := network.GetXupleManager().(xuples.XupleManager)
xupleManager.CreateXuple("space", alert, triggeringFacts)
// ... répéter pour chaque xuple ...

// Vérifier
space, _ := xupleManager.GetXupleSpace("alerts")
xuples := space.ListAll()
assert.Equal(t, 6, len(xuples))
```

### Après (Workflow Cible - Simple)

```go
// Test E2E - 3 étapes, 100% automatique
import "github.com/treivax/tsd/api"

// 1. Créer pipeline (une ligne)
pipeline := api.NewPipeline()

// 2. Ingérer fichier (tout est automatique)
result, err := pipeline.IngestFile("test.tsd")
require.NoError(t, err)

// 3. Vérifier résultats (API simple)
xuples := result.GetXuples("critical_alerts")
assert.Equal(t, 2, len(xuples))

commands := result.GetXuples("command_queue")
assert.Equal(t, 3, len(commands))
```

**Réduction de complexité**: 
- De 9 étapes à 3 étapes
- De 7 lignes manuelles à 0 ligne manuelle
- De ~100 lignes de code test à ~10 lignes

---

## 📝 Notes Importantes

### Respect des Standards TSD

Tous les prompts doivent respecter `.github/prompts/common.md`:
- En-tête copyright obligatoire sur tous les nouveaux fichiers
- Aucun hardcoding (valeurs, chemins, configs)
- Tout privé par défaut, exports minimaux
- Tests avec couverture > 80%
- GoDoc complet pour exports
- `make validate` doit passer

### Contexte des Sessions

Chaque prompt inclura:
- **Contexte minimal**: fichiers strictement nécessaires
- **Objectif précis**: une seule responsabilité par prompt
- **Livrables clairs**: fichiers à créer/modifier
- **Tests**: validation de l'objectif
- **Checklist**: étapes de validation

### Gestion des Dépendances

Si un prompt dépend d'un autre:
- Lire les fichiers créés par le prompt précédent
- Ne pas dupliquer le code
- Référencer les interfaces/types définis précédemment

---

## 🚀 Démarrage

Pour commencer, exécuter les prompts dans l'ordre:

```bash
# Session 1
cat scripts/maj-xuples/01-parser-faits-inline.md

# Session 2
cat scripts/maj-xuples/02-package-api-pipeline.md

# Session 3
cat scripts/maj-xuples/03-creation-auto-xuple-spaces.md

# ... et ainsi de suite
```

Chaque session est indépendante et peut être validée avant de passer à la suivante.

---

**Status**: ✅ Plan global défini  
**Prochaine étape**: Créer les prompts détaillés 01 à 07