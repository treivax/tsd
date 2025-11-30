# Prompt 11: Exemples et Migration Guide

**Objectif:** Créer des exemples pratiques et un guide de migration pour le Beta Sharing.

## Documentation requise

### 1. BETA_CHAINS_EXAMPLES.md

Créer `rete/BETA_CHAINS_EXAMPLES.md` avec:

- **10+ exemples concrets** couvrant tous les cas d'usage
- **Règles avec 2, 3, 5 jointures** pour montrer la scalabilité
- **Partage complet, partiel, aucun** (tous les scénarios)
- **Avant/après optimisation** avec visualisations ASCII
- **Métriques de chaque exemple** (temps, mémoire, partage)
- **Visualisations des chaînes** en ASCII art et Mermaid

**Structure attendue:**
```
# Exemples Concrets : Beta Chains

## Table des Matières
1. [Exemples Basiques](#exemples-basiques)
2. [Exemples de Partage](#exemples-de-partage)
3. [Exemples Avancés](#exemples-avancés)
4. [Visualisations](#visualisations)
5. [Métriques de Partage](#métriques-de-partage)
6. [Cas d'Usage Réels](#cas-dusage-réels)

## Exemples Basiques

### Exemple 1: Deux jointures simples
### Exemple 2: Trois jointures en cascade
### Exemple 3: Cinq jointures complexes
[...]

## Exemples de Partage

### Exemple 4: Partage complet (100%)
### Exemple 5: Partage partiel (50%)
### Exemple 6: Aucun partage (0%)
[...]

## Exemples Avancés

### Exemple 7: Optimisation de l'ordre de jointure
### Exemple 8: Réutilisation de préfixes
### Exemple 9: Cache de jointure
### Exemple 10: Monitoring en production
[...]
```

**Contenu de chaque exemple:**
- Code TSD des règles
- Chaîne beta créée (visualisation ASCII)
- Métriques détaillées:
  - Nombre de JoinNodes créés vs réutilisés
  - Ratio de partage (%)
  - Temps de construction
  - Mémoire économisée
  - Hits/miss du cache
- Comparaison avant/après optimisation
- Logs de construction
- Diagrammes Mermaid pour les cas complexes

### 2. BETA_CHAINS_MIGRATION.md

Créer `rete/BETA_CHAINS_MIGRATION.md` avec:

- **Guide de migration pas à pas** pour activer le beta sharing
- **Impact sur le code existant** (breaking changes, API changes)
- **Comment activer le beta sharing** avec exemples de code
- **Configuration recommandée** pour différents cas d'usage
- **Troubleshooting** avec solutions aux problèmes courants
- **Rollback si nécessaire** avec procédure complète

**Structure attendue:**
```
# Guide de Migration : Beta Chains

## Table des Matières
1. [Vue d'ensemble](#vue-densemble)
2. [Impact sur le code existant](#impact-sur-le-code-existant)
3. [Migration pas à pas](#migration-pas-à-pas)
4. [Configuration et tuning](#configuration-et-tuning)
5. [Troubleshooting](#troubleshooting)
6. [Rollback](#rollback)
7. [FAQ Migration](#faq-migration)

## Vue d'ensemble

### Qu'est-ce qui change ?
### Qui est impacté ?
### Compatibilité

## Impact sur le code existant

### Code qui continue de fonctionner
### Nouveau code optionnel
### Breaking changes (s'il y en a)
### Dépendances

## Migration pas à pas

### Étape 1: Prérequis
### Étape 2: Activation basique
### Étape 3: Configuration
### Étape 4: Validation
### Étape 5: Monitoring
### Étape 6: Tuning

## Configuration et tuning

### Configuration par défaut
### Configuration haute performance
### Configuration mémoire optimisée
### Configuration debugging

## Troubleshooting

### Problème 1: Beta sharing ne s'active pas
### Problème 2: Performance dégradée
### Problème 3: Fuite mémoire
### Problème 4: Erreurs de jointure
[...]

## Rollback

### Procédure de rollback
### Vérification post-rollback
### Logs et diagnostics
```

**Contenu détaillé:**
- Code exemples avant/après migration
- Commandes de configuration
- Tests de validation
- Métriques à surveiller
- Checkpoints de migration
- Procédures de rollback testées

### 3. examples/beta_chains/

Créer `examples/beta_chains/` avec:

- **Exemple exécutable en Go** (`main.go`)
- **Configuration avec/sans beta sharing** (`config.go`)
- **Affichage des métriques** (`metrics.go`)
- **Comparaison des performances** (`benchmark.go`)
- **README détaillé** (`README.md`)

**Structure du dossier:**
```
examples/beta_chains/
├── README.md
├── main.go
├── config.go
├── metrics.go
├── benchmark.go
├── scenarios/
│   ├── simple.go
│   ├── complex.go
│   └── advanced.go
└── go.mod (si nécessaire)
```

**Fonctionnalités de l'exemple:**
- Mode interactif pour choisir le scénario
- Comparaison side-by-side avec/sans beta sharing
- Affichage en temps réel des métriques
- Export des résultats en JSON/CSV
- Visualisation ASCII des chaînes construites
- Calcul automatique des gains de performance

**Scénarios à implémenter:**
1. **Simple**: 2-3 jointures, partage évident
2. **Complex**: 5+ jointures, optimisation d'ordre
3. **Advanced**: Cas réels (e-commerce, monitoring, etc.)

### 4. BETA_CHAINS_INDEX.md

Mettre à jour `rete/BETA_CHAINS_INDEX.md` avec:

- **Index centralisé** de toute la documentation
- **Quick start** pour chaque public cible
- **Liens vers tous les guides** avec descriptions
- **FAQ consolidée** (20+ questions/réponses)

**Sections à ajouter:**
```
## Nouveaux Documents

### Exemples Pratiques
- [BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md)
  - 15+ exemples exécutables
  - Visualisations complètes
  - Métriques détaillées

### Guide de Migration
- [BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md)
  - Migration pas à pas
  - Troubleshooting complet
  - Procédures de rollback

### Exemples de Code
- [examples/beta_chains/](../../examples/beta_chains/)
  - Code Go exécutable
  - Comparaisons de performance
  - Scénarios multiples

## Quick Start par Profil

### Pour les débutants
1. Lire BETA_NODE_SHARING.md
2. Exécuter examples/beta_chains/
3. Consulter BETA_CHAINS_EXAMPLES.md (exemples 1-5)

### Pour les développeurs
1. Lire BETA_CHAINS_USER_GUIDE.md
2. Suivre BETA_CHAINS_MIGRATION.md
3. Implémenter avec examples/beta_chains/ comme référence

### Pour les experts
1. Lire BETA_CHAINS_TECHNICAL_GUIDE.md
2. Analyser BETA_CHAINS_EXAMPLES.md (exemples avancés)
3. Optimiser avec les patterns du guide technique

## FAQ Consolidée

### Questions Générales
1. Qu'est-ce que le beta sharing ?
2. Quelle est la différence avec l'alpha sharing ?
3. Quels sont les bénéfices ?
[...]

### Questions Techniques
11. Comment calculer le hash d'un JoinNode ?
12. Comment optimiser l'ordre des jointures ?
13. Comment gérer le cache LRU ?
[...]

### Questions de Migration
21. Est-ce compatible avec mon code existant ?
22. Comment activer le beta sharing ?
23. Que faire si ça ne fonctionne pas ?
[...]
```

## Critères de succès

### Exemples exécutables
- ✅ Tous les exemples du guide peuvent être exécutés avec `go run`
- ✅ Chaque exemple affiche clairement les métriques
- ✅ Les visualisations ASCII sont correctes et lisibles
- ✅ Les comparaisons avant/après sont démonstrables

### Migration sans breaking change
- ✅ Le code existant fonctionne sans modification
- ✅ L'activation est opt-in (pas opt-out)
- ✅ La configuration par défaut est sûre
- ✅ Les tests existants passent tous

### Guide de troubleshooting complet
- ✅ 10+ problèmes courants documentés
- ✅ Solutions testées et vérifiées
- ✅ Procédures de diagnostic claires
- ✅ Logs et commandes de debug

### 15+ exemples au total
- ✅ 3+ exemples basiques (2-3 jointures)
- ✅ 5+ exemples de partage (complet/partiel/aucun)
- ✅ 7+ exemples avancés (optimisations, cache, etc.)
- ✅ Chaque exemple avec métriques complètes

### Documentation liée et cohérente
- ✅ Index centralisé à jour
- ✅ Références croisées entre documents
- ✅ Style uniforme (comme ALPHA_CHAINS_EXAMPLES.md)
- ✅ FAQ consolidée et complète
- ✅ Tous les liens fonctionnent

## Structure finale attendue

```
rete/
├── BETA_CHAINS_EXAMPLES.md        (~25-30 pages, 15+ exemples)
├── BETA_CHAINS_MIGRATION.md       (~20 pages, guide complet)
├── BETA_CHAINS_INDEX.md           (~10 pages, index mis à jour)
└── [autres docs existants]

examples/
├── beta_chains/
│   ├── README.md                  (~5 pages)
│   ├── main.go                    (~200 lignes)
│   ├── config.go                  (~100 lignes)
│   ├── metrics.go                 (~150 lignes)
│   ├── benchmark.go               (~200 lignes)
│   └── scenarios/
│       ├── simple.go              (~100 lignes)
│       ├── complex.go             (~150 lignes)
│       └── advanced.go            (~200 lignes)
└── [autres exemples]
```

## Notes importantes

### Compatibilité License MIT
- ✅ Tout le code produit doit être compatible avec la license MIT de TSD
- ✅ Pas de dépendances avec licenses incompatibles
- ✅ Headers de copyright appropriés dans tous les fichiers
- ✅ Attribution correcte si du code tiers est utilisé

### Style de documentation
- Suivre le style de ALPHA_CHAINS_EXAMPLES.md et ALPHA_CHAINS_MIGRATION.md
- Utiliser des émojis pour les sections (📊, 🔍, ⚡, etc.)
- ASCII art pour les visualisations simples
- Mermaid diagrams pour les cas complexes
- Code blocks avec syntaxe highlighting
- Tableaux pour les comparaisons

### Qualité du code
- Code Go idiomatique et propre
- Tests unitaires si approprié
- Gestion d'erreurs complète
- Documentation inline (godoc)
- Exemples self-contained (pas de dépendances externes non nécessaires)

### Métriques à inclure
Pour chaque exemple:
- Temps de construction (µs/ms)
- Nombre de JoinNodes créés
- Nombre de JoinNodes réutilisés
- Ratio de partage (%)
- Hits/miss du cache de jointure
- Mémoire utilisée (si pertinent)
- Comparaison avec/sans beta sharing

## Utilisation du prompt update-docs

Une fois les fichiers créés, utiliser le prompt `update-docs` pour:
1. Vérifier la cohérence entre tous les documents
2. Mettre à jour les références croisées
3. Compléter l'index BETA_CHAINS_INDEX.md
4. Valider que tous les exemples sont mentionnés
5. S'assurer que la FAQ est exhaustive

## Livrables finaux

1. ✅ `rete/BETA_CHAINS_EXAMPLES.md` - 25-30 pages, 15+ exemples
2. ✅ `rete/BETA_CHAINS_MIGRATION.md` - 20 pages, guide complet
3. ✅ `examples/beta_chains/` - Dossier complet avec code exécutable
4. ✅ `rete/BETA_CHAINS_INDEX.md` - Index mis à jour avec nouveaux liens
5. ✅ Tous les fichiers testés et validés
6. ✅ License MIT respectée partout