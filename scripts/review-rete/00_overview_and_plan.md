# 🔍 Revue de Code Module RETE - Vue d'Ensemble et Planification

**Objectif:** Revue systématique et complète du module `rete` selon les standards définis dans `.github/prompts/review.md`

---

## 📋 Contexte

Le module **rete** est le cœur du projet TSD, représentant 68% du code total (101,372 lignes sur 149,075).

### Statistiques Actuelles

```
Fichiers source:     152 fichiers .go (hors tests)
Fichiers tests:      123 fichiers *_test.go
Lignes de code:      52,671 lignes (source)
Lignes de tests:     48,701 lignes (tests rete)
Packages:            2 (rete, rete/internal/config)
Couverture tests:    80.8%
Complexité > 15:     98 fonctions (dont plusieurs dans rete)
```

### Problématiques Identifiées

1. **Complexité élevée** - Plusieurs fonctions >15 complexité cyclomatique
2. **Taille du module** - 68% du projet, considérer sous-modules
3. **Duplication potentielle** - À analyser
4. **Documentation** - Vérifier exhaustivité GoDoc

---

## 🎯 Objectifs de la Revue

### Primaires
- ✅ Garantir conformité aux standards Go et projet
- ✅ Réduire complexité cyclomatique (<15 partout)
- ✅ Éliminer duplication de code
- ✅ Améliorer lisibilité et maintenabilité
- ✅ Valider encapsulation et architecture

### Secondaires
- ✅ Améliorer documentation (GoDoc)
- ✅ Identifier opportunités optimisation
- ✅ Vérifier gestion erreurs
- ✅ Valider tests (couverture >80%)

---

## 📊 Organisation de la Revue

La revue est organisée en **10 prompts numérotés**, chacun couvrant un sous-ensemble cohérent du module, adapté aux contraintes de contexte (128k tokens).

### Découpage par Domaine Fonctionnel

#### Prompt 01 - Core RETE (Nœuds Fondamentaux)
**Fichiers:** 8 fichiers (~2,000 lignes)
- `network.go` - Réseau RETE principal
- `node_*.go` - Nœuds (alpha, beta, join, terminal, etc.)
- `memory.go` - Gestion mémoire
- `token.go` - Tokens et bindings

**Focus:** Architecture de base, performance, encapsulation

#### Prompt 02 - Bindings et Chaînes Immuables
**Fichiers:** 6 fichiers (~1,500 lignes)
- `binding_chain.go` - Système immuable
- `beta_chain.go` - Chaînes beta
- `chain_*.go` - Configuration et métriques
- `token_metadata.go` - Métadonnées

**Focus:** Immuabilité, thread-safety, performance

#### Prompt 03 - Alpha Network (Construction et Partage)
**Fichiers:** 10 fichiers (~2,500 lignes)
- `alpha_builder.go` - Construction
- `alpha_chain_*.go` - Chaînes et extraction
- `alpha_sharing*.go` - Mécanismes de partage
- `alpha_condition*.go` - Évaluation conditions

**Focus:** Partage, normalisation, optimisation

#### Prompt 04 - Beta Network (Jointures et Partage)
**Fichiers:** 8 fichiers (~2,200 lignes)
- `beta_chain_builder*.go` - Construction chaînes
- `beta_sharing*.go` - Partage JoinNodes
- `beta_chain_optimizer.go` - Optimisations
- `node_join*.go` - Logique jointures

**Focus:** Partage JoinNode, cascade, performance

#### Prompt 05 - Expressions Arithmétiques
**Fichiers:** 8 fichiers (~2,800 lignes)
- `arithmetic_*.go` - Décomposition, cache, évaluation
- `expression_*.go` - Analyse et normalisation
- `nested_or_normalizer*.go` - Normalisation OR

**Focus:** Décomposition, cache, optimisation

#### Prompt 06 - Builders et Construction du Réseau
**Fichiers:** 12 fichiers (~3,000 lignes)
- `builder_*.go` - Construction règles (alpha, join, exists, accumulator)
- `builder_orchestration.go` - Orchestration
- `builder_types.go` - Types et structures

**Focus:** Architecture builders, séparation responsabilités

#### Prompt 07 - Actions et Exécution
**Fichiers:** 8 fichiers (~1,800 lignes)
- `action_*.go` - Exécuteur, handlers, contexte
- `command_*.go` - Commandes (assertions, rétractions)
- `rule_activation.go` - Activation règles

**Focus:** Gestion actions, thread-safety, erreurs

#### Prompt 08 - Pipeline et Validation
**Fichiers:** 6 fichiers (~2,000 lignes)
- `constraint_pipeline*.go` - Pipeline principal
- `constraint_validator*.go` - Validation contraintes
- `type_checker.go` - Vérification types
- `coherence_*.go` - Cohérence

**Focus:** Validation, gestion erreurs, robustesse

#### Prompt 09 - Métriques et Diagnostics
**Fichiers:** 10 fichiers (~2,500 lignes)
- `*_metrics.go` - Collecte métriques
- `*_stats.go` - Statistiques
- `debug_*.go` - Debug et logging
- `print_network*.go` - Visualisation

**Focus:** Observabilité, debug, documentation

#### Prompt 10 - Utilitaires et Helpers
**Fichiers:** Restants (~1,500 lignes)
- `utils.go` - Utilitaires généraux
- `circular_dependency_detector.go`
- `evaluator.go`
- `rule_router.go`
- Autres helpers

**Focus:** Généricité, réutilisabilité, simplicité

---

## 📝 Format de Chaque Prompt

Chaque prompt suivra cette structure standardisée:

```markdown
# 🔍 Revue RETE - [Numéro]: [Domaine]

## 📋 Périmètre
- Fichiers couverts (liste)
- Lignes de code
- Complexité estimée

## 🎯 Objectifs Spécifiques
- Objectif 1
- Objectif 2
- ...

## 📖 Instructions

1. **Analyser** les fichiers du périmètre
2. **Identifier** problèmes selon checklist review.md
3. **Proposer** corrections/améliorations
4. **Implémenter** les changements
5. **Valider** (tests, métriques)
6. **Documenter** les modifications

## ✅ Checklist de Revue

[Checklist spécifique au domaine]

## 📊 Métriques Attendues

- Complexité: < 15 partout
- Couverture: > 80%
- Duplication: Minimale
- GoDoc: 100% exports

## 🎯 Livrables

- [ ] Code refactoré (si nécessaire)
- [ ] Tests validés
- [ ] Documentation mise à jour
- [ ] Rapport de revue
```

---

## 🔄 Workflow d'Exécution

### Pour Chaque Prompt (01-10)

1. **Lire le prompt** dans l'ordre numérique
2. **Charger les fichiers** du périmètre dans le contexte
3. **Appliquer la checklist** de review.md
4. **Identifier les problèmes** (complexité, duplication, etc.)
5. **Proposer les corrections**
6. **Implémenter les changements** de manière incrémentale
7. **Valider avec tests** (`go test ./rete/...`)
8. **Générer rapport** pour ce prompt
9. **Passer au prompt suivant**

### Validation Entre Prompts

Après chaque prompt:
```bash
# Tests
go test ./rete/... -v

# Complexité
gocyclo -over 15 rete/

# Formatage
go fmt ./rete/...

# Vérifications
go vet ./rete/...
```

---

## 📈 Critères de Succès Globaux

### Métriques Cibles (Fin de Revue)

| Métrique | Actuel | Cible | Critique |
|----------|--------|-------|----------|
| Complexité max | 48 | <20 | ⚠️ Oui |
| Fonctions >15 | 98 | <30 | ⚠️ Oui |
| Couverture | 80.8% | >85% | Non |
| Duplication | ? | <5% | Oui |
| GoDoc exports | ~90% | 100% | Non |
| Warnings vet | 0 | 0 | ✅ OK |

### Qualité Globale

- ✅ Tous les tests passent (100%)
- ✅ Aucune régression fonctionnelle
- ✅ Performance préservée ou améliorée
- ✅ Architecture SOLID respectée
- ✅ Code auto-documenté
- ✅ Encapsulation rigoureuse

---

## 🎓 Principes Directeurs

### Refactoring
1. **Incrémental** - Petites étapes validées
2. **Testé** - Tests passent après chaque changement
3. **Documenté** - Changements expliqués
4. **Réversible** - Commits atomiques

### Qualité
1. **Simplicité** - La solution la plus simple
2. **Lisibilité** - Code auto-documenté
3. **Maintenabilité** - Facile à modifier
4. **Performance** - Pas de dégradation

### Standards
1. **Go idiomatique** - Conventions Go respectées
2. **Projet cohérent** - Standards TSD appliqués
3. **Documentation** - GoDoc exhaustif
4. **Tests** - Couverture >80%

---

## 📚 Références

- `.github/prompts/review.md` - Standards de revue
- `.github/prompts/common.md` - Standards projet
- `REPORTS/MAINTENANCE_REPORT.md` - État actuel
- `docs/architecture/` - Documentation architecture

---

## 🚀 Prochaines Étapes

1. **Lire ce prompt** (00) - Comprendre la structure
2. **Exécuter prompt 01** - Core RETE
3. **Continuer séquentiellement** jusqu'au prompt 10
4. **Générer rapport final** - Synthèse complète

---

## 📊 Suivi de Progression

| Prompt | Domaine | Fichiers | Status | Durée |
|--------|---------|----------|--------|-------|
| 00 | Overview | - | ✅ | - |
| 01 | Core RETE | 8 | ⏳ Pending | - |
| 02 | Bindings | 6 | ⏳ Pending | - |
| 03 | Alpha Network | 10 | ⏳ Pending | - |
| 04 | Beta Network | 8 | ⏳ Pending | - |
| 05 | Arithmétique | 8 | ⏳ Pending | - |
| 06 | Builders | 12 | ⏳ Pending | - |
| 07 | Actions | 8 | ⏳ Pending | - |
| 08 | Pipeline | 6 | ⏳ Pending | - |
| 09 | Métriques | 10 | ⏳ Pending | - |
| 10 | Utilitaires | ~10 | ⏳ Pending | - |

---

**Date création:** 2024-12-15  
**Version:** 1.0  
**Auteur:** Équipe TSD  
**Statut:** 📋 Plan validé - Prêt pour exécution