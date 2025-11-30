# Prompt : Créer la documentation technique complète pour le Beta Sharing System

## Objectif

Créer la documentation technique des BetaChains (partage de JoinNodes) dans le moteur RETE de TSD.

## Documentation requise

### 1. BETA_CHAINS_TECHNICAL_GUIDE.md (~30-40 pages)

Guide technique détaillé destiné aux développeurs expérimentés et contributeurs.

**Contenu obligatoire :**

- **Architecture détaillée**
  - Vue d'ensemble du système
  - Composants principaux (BetaChainBuilder, BetaSharingRegistry, BetaJoinCache, JoinOrderOptimizer)
  - Structure des JoinNodes (left/right memory, tests de jointure)
  - Mécanisme de partage de préfixes
  - Diagrammes ASCII de l'architecture

- **Algorithmes de partage et construction**
  - Construction de chaîne Beta (étape par étape)
  - Normalisation de pattern de jointure
  - Optimisation de l'ordre de jointure (heuristiques)
  - Algorithme de jointure runtime
  - Invalidation du cache de jointure

- **Normalisation des patterns de jointure**
  - Principes de normalisation
  - Ordre des contraintes
  - Opérateurs commutatifs
  - Opérateurs équivalents
  - Contraintes redondantes
  - Algorithme de hachage (FNV-1a)

- **Lifecycle management**
  - Référencement des JoinNodes (RefCount)
  - Ajout d'une règle
  - Suppression d'une règle
  - Garbage collection
  - Diagramme de lifecycle

- **API de référence**
  - BetaChainBuilder (NewBetaChainBuilder, BuildChain, GetMetrics)
  - BetaSharingRegistry (NewBetaSharingRegistry, GetOrCreateJoinNode, GetSharingStats)
  - BetaJoinCache (NewBetaJoinCache, GetJoinResult, SetJoinResult, InvalidateForFact, GetStats)
  - BetaChainMetrics (NewBetaChainMetrics, RecordChainBuild, GetSnapshot)
  - Signatures complètes avec paramètres et valeurs de retour
  - Exemples de code pour chaque méthode

- **Cas limites (edge cases)**
  - Règles avec un seul pattern
  - Patterns identiques
  - Cycles dans les dépendances
  - Patterns avec négation
  - Mémoire pleine (OOM)
  - Patterns avec wildcards
  - Règles supprimées avec nœuds partagés
  - Hash collisions
  - Ordre d'évaluation non-déterministe
  - Deadlocks lors de suppressions concurrentes

- **Optimisations**
  - Cache de hash (LRU)
  - Cache de jointure (LRU)
  - Optimisation de l'ordre de jointure
  - Index inversé pour invalidation
  - Partage de préfixes
  - Lazy evaluation
  - Batch processing

- **Internals**
  - Structure mémoire détaillée (layout en bytes)
  - Overhead du partage
  - Complexité algorithmique (tableaux Big-O)
  - Concurrence et thread-safety (modèles de verrouillage)
  - Debugging et observabilité (logging, Prometheus, tracing)
  - Profiling (CPU, memory, mutex)

### 2. BETA_CHAINS_USER_GUIDE.md (~20-30 pages)

Guide utilisateur destiné aux développeurs d'applications et intégrateurs.

**Contenu obligatoire :**

- **Introduction au Beta Sharing**
  - Qu'est-ce que le Beta Sharing ?
  - Pourquoi est-ce important ?
  - Exemple simple avec diagrammes avant/après

- **Bénéfices et cas d'usage**
  - Réduction de la consommation mémoire (40-70%)
  - Amélioration des performances (30-50%)
  - Scalabilité améliorée
  - Cache de jointure
  - Cas d'usage détaillés :
    - Règles de validation avec base commune
    - Règles de tarification dynamique
    - Détection de fraude multi-critères
    - Workflows multi-étapes

- **Comment l'activer/configurer**
  - Activation du Beta Sharing (code Go)
  - Configuration du cache de hash
  - Configuration du cache de jointure
  - Configuration avancée (optimisation, GC)

- **Exemples pratiques**
  - Système de recommandation (complet avec code)
  - Validation de commandes (complet avec code)
  - Monitoring et alertes (complet avec code)
  - Affichage des métriques

- **Guide de dépannage**
  - Sharing ratio faible (causes, diagnostic, solutions)
  - Cache hit rate faible (causes, diagnostic, solutions)
  - Consommation mémoire élevée (causes, diagnostic, solutions)
  - Performance dégradée (causes, diagnostic, solutions)
  - Résultats incorrects (causes, diagnostic, solutions)

- **FAQ**
  - 10+ questions fréquentes avec réponses détaillées
  - Compatibilité Alpha/Beta Sharing
  - Overhead du Beta Sharing
  - Désactivation par règle
  - Thread-safety
  - Modification de règles
  - Mesure de l'efficacité
  - Négations
  - Debugging
  - Taille des caches

- **Meilleures pratiques**
  - Nommer les variables de façon cohérente
  - Grouper les contraintes
  - Patterns communs en préfixe
  - Ajuster les caches selon workload
  - Monitorer les métriques
  - Garbage collection périodique
  - Tests avec/sans Beta Sharing
  - Documentation des dépendances

### 3. BETA_NODE_SHARING.md (~15 pages)

Guide conceptuel pour tous niveaux.

**Contenu obligatoire :**

- **Concepts de base**
  - Qu'est-ce qu'un Beta Node / JoinNode ?
  - Qu'est-ce que le Beta Node Sharing ?
  - Pourquoi partager les Beta Nodes ?
  - Exemple concret avec chiffres

- **Différence avec Alpha Sharing**
  - Tableau comparatif (nœuds, tests, mémoire, complexité, impact)
  - Alpha Node vs Beta Node (explications détaillées)
  - Visualisation de la différence (diagramme réseau complet)
  - Différence dans le code (extraits comparatifs)

- **Quand les JoinNodes sont partagés**
  - Critères de partage (3 conditions)
  - Exemple 1 : Partage complet
  - Exemple 2 : Partage partiel
  - Exemple 3 : Pas de partage
  - Exemple 4 : Partage de préfixe
  - Normalisation et partage (exemples concrets)

- **Diagrammes explicatifs**
  - Diagramme 1 : Anatomie d'un JoinNode (ASCII art détaillé)
  - Diagramme 2 : Flux d'activation (gauche/droite)
  - Diagramme 3 : Cycle de vie (états avec transitions)
  - Diagramme 4 : Cache de jointure (workflow)

- **Exemples visuels**
  - Exemple 1 : Réseau sans Beta Sharing (3 règles)
  - Exemple 2 : Réseau avec Beta Sharing (3 règles)
  - Exemple 3 : Chaîne complexe avec préfixes partagés
  - Exemple 4 : Patterns avec variations
  - Exemple 5 : Avant/Après optimisation (side-by-side)

- **Mécanismes internes**
  - Calcul de la clé de partage (algorithme détaillé)
  - Registre de partage (structure de données)
  - Séquence de construction (flowchart)
  - Gestion de la concurrence (threads)

### 4. Index centralisé (BETA_CHAINS_INDEX.md)

**Contenu obligatoire :**

- Vue d'ensemble de toute la documentation
- Liens vers les 3 guides principaux (descriptions)
- Documentation complémentaire (performance, implémentation, migration)
- Quick Start (3 parcours : nouveaux utilisateurs, développeurs, intégrateurs)
- Index par sujet (architecture, algorithmes, configuration, performance, dépannage, API)
- Diagrammes et visualisations (liens)
- Glossaire (termes clés, acronymes)
- Ressources additionnelles (externe, code source, tests)
- Contribution (guidelines)
- Historique des versions
- Support et licence

## Diagrammes et visualisations

**Requis dans la documentation :**

- **ASCII art** :
  - Architecture multi-couches
  - Structure des JoinNodes
  - Flux d'activation
  - Réseaux RETE (avant/après partage)
  - Diagrammes de séquence
  - États du lifecycle

- **Exemples de chaînes** :
  - Avant partage : 3 règles → 3 JoinNodes distincts
  - Après partage : 3 règles → 1 JoinNode partagé
  - Calculs de gain (mémoire, performance)

## Critères de succès

✅ Documentation claire et complète (couvre tous les aspects)
✅ Exemples exécutables (code Go complet et testable)
✅ Diagrammes visuels (ASCII art détaillé)
✅ Style cohérent avec ALPHA_CHAINS_* (même structure, même format)
✅ Licence MIT (en-tête dans chaque document)
✅ Liens entre documents (références croisées)
✅ Pas de code propriétaire ou copié
✅ Tableaux et listes bien formatés (Markdown)

## Structure des fichiers

```
rete/
├── BETA_CHAINS_TECHNICAL_GUIDE.md   (30-40 pages)
├── BETA_CHAINS_USER_GUIDE.md        (20-30 pages)
├── BETA_NODE_SHARING.md             (15 pages)
└── BETA_CHAINS_INDEX.md             (index centralisé)
```

## Style et format

- **Format** : Markdown
- **Langage** : Français
- **Code** : Go (avec syntaxe highlighting)
- **Diagrammes** : ASCII art + descriptions
- **Sections** : Table des matières au début
- **Licence** : MIT (footer dans chaque doc)

## Utilisation du prompt update-docs

Ce prompt doit être utilisé conjointement avec le prompt `update-docs` pour :
- Maintenir la cohérence avec la documentation existante
- Respecter les conventions du projet
- Assurer la qualité et la complétude

## Références

Inspirez-vous des documents Alpha existants pour le style :
- `ALPHA_CHAINS_TECHNICAL_GUIDE.md`
- `ALPHA_CHAINS_USER_GUIDE.md`
- `ALPHA_NODE_SHARING.md`
- `ALPHA_CHAINS_INDEX.md`

## Livrables finaux

1. ✅ BETA_CHAINS_TECHNICAL_GUIDE.md (guide technique complet)
2. ✅ BETA_CHAINS_USER_GUIDE.md (guide utilisateur complet)
3. ✅ BETA_NODE_SHARING.md (concepts et mécanismes)
4. ✅ BETA_CHAINS_INDEX.md (index centralisé)
5. ✅ Diagrammes ASCII dans les documents
6. ✅ Exemples de code exécutables
7. ✅ Liens de navigation entre documents

---

**Note** : Tous les documents doivent être compatibles avec la licence MIT utilisée par TSD.