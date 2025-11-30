# Prompt 1: Analyse de l'Existant des BetaNodes

**Objectif:** Comprendre l'implémentation actuelle des JoinNodes et identifier les opportunités d'optimisation.

**Fichier prompt:** `.github/prompts/beta-analyze-existing.md`

---

## Contexte

Tu es un expert en moteurs RETE et optimisation de code.

Le projet TSD dispose d'un moteur RETE mature où les AlphaNodes bénéficient déjà d'un système de partage très efficace (70-85% de réutilisation). Les BetaNodes (JoinNodes) n'ont actuellement aucun mécanisme de partage, ce qui représente une opportunité d'optimisation majeure.

---

## Mission

Analyse l'implémentation actuelle des BetaNodes (JoinNodes) dans le projet TSD pour identifier les opportunités d'optimisation via le partage de nœuds.

---

## Axes d'Analyse

### 1. Architecture actuelle
   - Structure des JoinNodes
   - Comment sont-ils créés et connectés
   - Gestion de la mémoire et des tokens
   - Algorithme de jointure actuel

### 2. Identifier les patterns
   - Patterns de jointure courants
   - Conditions de jointure dupliquées
   - Opportunités de partage

### 3. Comparaison avec AlphaNodes
   - Similitudes et différences
   - Ce qui a bien fonctionné pour les alpha
   - Adaptations nécessaires pour les beta

### 4. Points d'amélioration
   - Où le partage peut être appliqué
   - Quelles optimisations sont possibles
   - Risques et contraintes

---

## Livrables Attendus

### 1. Rapport d'analyse
**Fichier:** `rete/docs/BETA_NODES_ANALYSIS.md`

Contenu attendu:
- État actuel du code
- Architecture détaillée des JoinNodes
- Patterns de jointure identifiés
- Comparaison avec AlphaNodes
- Liste des opportunités d'optimisation
- Plan technique d'implémentation
- Analyse des risques et contraintes
- Métriques et critères de succès

### 2. Diagrammes d'architecture
**Fichier:** `rete/docs/BETA_NODES_ARCHITECTURE_DIAGRAMS.md`

Contenu attendu:
- Diagrammes de l'architecture actuelle (sans partage)
- Diagrammes de l'architecture proposée (avec partage)
- Flux de données
- Comparaisons visuelles avant/après

### 3. Liste des opportunités d'optimisation
**Fichier:** `rete/docs/BETA_OPTIMIZATION_OPPORTUNITIES.md`

Contenu attendu:
- Liste priorisée des optimisations
- Impact estimé de chaque optimisation
- Effort d'implémentation
- Dépendances entre optimisations
- Roadmap recommandée

---

## Fichiers à Analyser

### Code Principal
- `rete/node_join.go` - Implémentation des JoinNodes
- `rete/network.go` - Intégration réseau
- `rete/constraint_pipeline_builder.go` - Construction du réseau
- `rete/alpha_sharing.go` - Référence pour le partage (AlphaNodes)
- `rete/node_lifecycle.go` - Gestion du cycle de vie

### Tests
- `rete/node_join_cascade_test.go` - Tests des cascades
- `rete/alpha_sharing_test.go` - Tests de partage Alpha (référence)
- `rete/alpha_sharing_integration_test.go` - Tests d'intégration

### Documentation Existante
- `rete/ALPHA_NODE_SHARING.md` - Guide du partage Alpha (référence)
- `rete/NODE_LIFECYCLE_README.md` - Cycle de vie des nœuds

---

## Format du Rapport

Le rapport principal doit inclure:

1. **Executive Summary**
   - Vue d'ensemble
   - Problème principal
   - Impact attendu

2. **Architecture Actuelle**
   - Structure détaillée
   - Algorithmes
   - Intégration réseau

3. **Patterns Identifiés**
   - Patterns de jointure courants
   - Cas d'usage réels
   - Statistiques (si possible)

4. **Comparaison Alpha vs Beta**
   - Similitudes exploitables
   - Différences à gérer
   - Leçons apprises

5. **Opportunités d'Optimisation**
   - Liste détaillée et priorisée
   - Impact quantifié
   - Faisabilité technique

6. **Plan Technique**
   - Phases d'implémentation
   - Timeline estimée
   - Prérequis et dépendances

7. **Risques et Contraintes**
   - Risques techniques identifiés
   - Stratégies de mitigation
   - Contraintes à respecter

8. **Métriques et Validation**
   - Critères de succès
   - Métriques à mesurer
   - Benchmarks attendus

---

## Critères de Qualité

- ✅ Analyse technique approfondie et rigoureuse
- ✅ Recommandations concrètes et actionnables
- ✅ Impact quantifié (chiffres, estimations)
- ✅ Diagrammes clairs et informatifs
- ✅ Plan d'implémentation réaliste
- ✅ Documentation claire et structurée

---

## Notes

- Inspire-toi de l'implémentation réussie du partage des AlphaNodes
- Identifie les patterns réutilisables
- Propose une approche incrémentale (phases)
- Quantifie l'impact attendu (mémoire, performance)
- Considère la compatibilité ascendante
- Documente les décisions techniques

---

**Prompt créé par:** Équipe TSD
**Date:** 2025-01-27
**Statut:** Prêt pour exécution