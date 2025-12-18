# Documentation Xuples - Analyse et Refactoring

Ce répertoire contient l'analyse complète du système d'actions et terminal nodes de TSD, ainsi que la planification de la refonte vers l'architecture xuples.

## 📚 Documents d'Analyse

### [00-INDEX.md](./00-INDEX.md) - **COMMENCER ICI**

Synthèse générale de l'analyse et stratégie de migration complète.

**Contenu** :
- Vue d'ensemble du projet
- Cartographie du code
- Points d'intervention identifiés
- Stratégie de migration en 3 phases
- Risques et mitigations

**👉 À lire en premier pour comprendre l'ensemble du projet**

---

### [01-current-action-parsing.md](./01-current-action-parsing.md)

Analyse du parsing des actions dans le langage TSD.

**Contenu** :
- Grammaire PEG complète
- Structures AST (ActionDefinition, Action, JobCall)
- Flux de parsing (TSD → AST → Go)
- Validations existantes
- Points d'intervention

**Utile pour** : Comprendre comment les actions sont définies et parsées

---

### [02-terminal-nodes.md](./02-terminal-nodes.md)

Architecture des Terminal Nodes et stockage des tokens.

**Contenu** :
- Structure TerminalNode et BaseNode
- Cycle de vie d'un token
- Exécution des actions
- Interface avec ActionExecutor
- Diagramme de séquence

**Utile pour** : Comprendre comment les tokens sont stockés et les actions exécutées

---

### [03-token-fact-structures.md](./03-token-fact-structures.md)

Structures de données Token et Fact.

**Contenu** :
- Structure Fact complète
- Structure Token complète
- TokenMetadata (traçage)
- BindingChain immuable
- WorkingMemory et indexation
- Implications pour xuples

**Utile pour** : Comprendre les structures de données fondamentales

---

### [04-action-executor.md](./04-action-executor.md)

ActionExecutor et son interface.

**Contenu** :
- Interface ActionHandler
- ActionRegistry (thread-safe)
- ActionExecutor (exécution, évaluation)
- ExecutionContext
- Évaluation des arguments
- Propositions d'actions par défaut

**Utile pour** : Comprendre comment les actions sont exécutées

---

### [05-existing-tests.md](./05-existing-tests.md)

Recensement et analyse des tests existants.

**Contenu** :
- Statistiques (222 fichiers de test)
- Patterns de test utilisés
- Couverture par fonctionnalité
- Points faibles identifiés
- Recommandations pour nouveaux tests

**Utile pour** : Comprendre la couverture de tests actuelle et ce qui manque

---

## 📊 Rapports et Planification

Ces documents se trouvent dans `/REPORTS/` et à la racine du projet :

### [SYNTHESE-EXECUTIVE-2025-12-17.md](/REPORTS/SYNTHESE-EXECUTIVE-2025-12-17.md)

Synthèse exécutive de toute l'analyse et du refactoring Phase 1.

**Contenu** :
- Objectifs de la mission
- Résultats clés
- Principales découvertes
- État des tests
- Recommandations
- Métriques de qualité

**👉 À lire pour avoir une vue d'ensemble rapide des résultats**

---

### [refactoring-phase1-2025-12-17.md](/REPORTS/refactoring-phase1-2025-12-17.md)

Rapport détaillé du refactoring Phase 1.

**Contenu** :
- Modifications effectuées (code)
- Fix thread-safety `generateTokenID()`
- Refactoring `executeAction()`
- Validation (tests)
- Impact et bénéfices

**Utile pour** : Comprendre exactement quelles modifications ont été faites au code

---

### [TODO-XUPLES.md](/TODO-XUPLES.md)

Plan d'action détaillé pour les Phases 2 et 3.

**Contenu** :
- Phase 1 : ✅ Terminée
- Phase 2 : Tests manquants, ActionRegistry, actions par défaut
- Phase 3 : Architecture xuples complète
- Priorités et estimations
- Quick start pour développeurs

**👉 À consulter pour savoir quoi faire ensuite**

---

## 🚀 Quick Start

### Pour comprendre le système actuel

1. Lire [00-INDEX.md](./00-INDEX.md) pour la vue d'ensemble
2. Lire [02-terminal-nodes.md](./02-terminal-nodes.md) pour comprendre les terminal nodes
3. Lire [04-action-executor.md](./04-action-executor.md) pour comprendre l'exécution

### Pour contribuer au refactoring

1. Lire [TODO-XUPLES.md](/TODO-XUPLES.md) pour voir les tâches à faire
2. Lire [refactoring-phase1-2025-12-17.md](/REPORTS/refactoring-phase1-2025-12-17.md) pour comprendre ce qui a été fait
3. Choisir une tâche dans Phase 2 et commencer !

### Pour ajouter une action

1. Lire [04-action-executor.md](./04-action-executor.md) section "Propositions pour actions par défaut"
2. Créer un nouveau fichier `rete/action_<nom>.go`
3. Implémenter l'interface `ActionHandler`
4. Enregistrer dans `ActionExecutor.RegisterDefaultActions()`
5. Ajouter tests

---

## 📈 Progression

| Phase | Statut | Documents | Code | Tests |
|-------|--------|-----------|------|-------|
| **Phase 1** | ✅ Terminée | 6 docs (3650+ lignes) | 2 fichiers (64 lignes) | ✅ 100% passent |
| **Phase 2** | 📋 Planifiée | - | À faire | À ajouter |
| **Phase 3** | 📅 À planifier | - | À concevoir | À créer |

---

## 🎯 Objectif Final

Créer une architecture **xuples** qui :
- Sépare clairement RETE (pattern matching) et Xuples (tuple-space)
- Améliore la gestion mémoire (stratégies de rétention)
- Enrichit les fonctionnalités (actions par défaut, métriques, callbacks)
- Facilite l'extensibilité future

---

## 📝 Standards et Références

- [common.md](../../.github/prompts/common.md) - Standards du projet
- [review.md](../../.github/prompts/review.md) - Process de revue
- [Effective Go](https://go.dev/doc/effective_go) - Bonnes pratiques Go

---

**Date de création** : 2025-12-17  
**Dernière mise à jour** : 2025-12-17  
**Auteur** : Analyse automatique pour refonte xuples
