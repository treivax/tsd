# Index de la Documentation Beta Chains

## Vue d'ensemble

Ce document sert d'index centralisé pour toute la documentation relative au système de **Beta Chains** (partage de JoinNodes) dans le moteur RETE de TSD.

---

## Documentation principale

### 1. Exemples Pratiques (BETA_CHAINS_EXAMPLES.md)

**Public cible :** Tous niveaux, développeurs cherchant des exemples concrets

**Contenu :**
- 15+ exemples exécutables et détaillés
- Règles avec 2, 3, 5 jointures
- Partage complet, partiel, aucun
- Métriques et comparaisons avant/après
- Visualisations ASCII et Mermaid
- Cas d'usage réels (e-commerce, monitoring, banking)
- Optimisations de performance

**Niveau :** Débutant à Avancé

**Lien :** [BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md)

---

### 2. Guide de Migration (BETA_CHAINS_MIGRATION.md)

**Public cible :** Développeurs, Ops/DevOps, équipes de migration

**Contenu :**
- Migration pas à pas (0 breaking change)
- Impact sur le code existant
- Configuration et tuning (3 configs prédéfinies)
- Troubleshooting complet (5+ problèmes courants)
- Procédure de rollback
- FAQ migration (20+ questions)
- Métriques à surveiller

**Niveau :** Intermédiaire

**Lien :** [BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md)

---

### 3. Guide Technique (BETA_CHAINS_TECHNICAL_GUIDE.md)

**Public cible :** Développeurs expérimentés, contributeurs au moteur RETE

**Contenu :**
- Architecture détaillée du système de partage
- Algorithmes de construction et normalisation
- Patterns de jointure et optimisation
- Lifecycle management des JoinNodes
- API Reference complète
- Gestion des cas edge
- Optimisations avancées
- Internals et structures mémoire

**Niveau :** Avancé

**Lien :** [BETA_CHAINS_TECHNICAL_GUIDE.md](./BETA_CHAINS_TECHNICAL_GUIDE.md)

---

### 4. Guide Utilisateur (BETA_CHAINS_USER_GUIDE.md)

**Public cible :** Développeurs d'applications utilisant TSD, intégrateurs

**Contenu :**
- Introduction au Beta Sharing
- Bénéfices et cas d'usage
- Configuration et activation
- Exemples pratiques (recommandations, validation, monitoring)
- Guide de dépannage
- FAQ
- Meilleures pratiques

**Niveau :** Intermédiaire

**Lien :** [BETA_CHAINS_USER_GUIDE.md](./BETA_CHAINS_USER_GUIDE.md)

---

### 5. Concepts de Beta Node Sharing (BETA_NODE_SHARING.md)

**Public cible :** Tous niveaux, introduction conceptuelle

**Contenu :**
- Concepts de base du Beta Sharing
- Différence entre Alpha et Beta Sharing
- Critères de partage des JoinNodes
- Diagrammes explicatifs (ASCII art & Mermaid)
- Exemples visuels avant/après
- Mécanismes internes simplifiés

**Niveau :** Débutant à Intermédiaire

**Lien :** [BETA_NODE_SHARING.md](./BETA_NODE_SHARING.md)

---

---

### 6. Code Exécutable (examples/beta_chains/)

**Public cible :** Développeurs voulant tester et expérimenter

**Contenu :**
- Code Go exécutable et interactif
- 3 scénarios (Simple, Complex, Advanced)
- Comparaison avec/sans Beta Sharing
- Métriques en temps réel
- Export JSON/CSV des résultats
- Benchmarks de performance

**Niveau :** Tous niveaux

**Lien :** [examples/beta_chains/](../../examples/beta_chains/)

---

## Documentation complémentaire

### Performance et Benchmarking

- **[BETA_PERFORMANCE_REPORT.md](./docs/BETA_PERFORMANCE_REPORT.md)** : Rapport détaillé des benchmarks de performance
- **[BETA_BENCHMARK_README.md](./BETA_BENCHMARK_README.md)** : Guide d'exécution des benchmarks
- **[BETA_BENCHMARK_STATUS.md](./BETA_BENCHMARK_STATUS.md)** : Statut et résultats des tests

### Implémentation et Intégration

- **[BETA_CHAIN_BUILDER_README.md](./BETA_CHAIN_BUILDER_README.md)** : Documentation du BetaChainBuilder
- **[BETA_CHAIN_BUILDER_SUMMARY.md](./BETA_CHAIN_BUILDER_SUMMARY.md)** : Résumé technique du builder
- **[BETA_JOIN_CACHE_README.md](./BETA_JOIN_CACHE_README.md)** : Cache de résultats de jointure
- **[BETA_JOIN_CACHE_SUMMARY.md](./BETA_JOIN_CACHE_SUMMARY.md)** : Résumé du cache

### Migration et Déploiement

- **[BETA_SHARING_MIGRATION.md](./BETA_SHARING_MIGRATION.md)** : Guide de migration vers Beta Sharing
- **[BETA_SHARING_INTEGRATION_SUMMARY.md](./BETA_SHARING_INTEGRATION_SUMMARY.md)** : Résumé de l'intégration
- **[BETA_SHARING_QUICK_REF.md](./BETA_SHARING_QUICK_REF.md)** : Référence rapide

### Livrables et Suivi

- **[BETA_BENCHMARK_DELIVERABLES.md](./BETA_BENCHMARK_DELIVERABLES.md)** : Liste des livrables de benchmark
- **[BETA_BENCHMARK_INDEX.md](./BETA_BENCHMARK_INDEX.md)** : Index des benchmarks
- **[BETA_BENCHMARK_PR_SUMMARY.md](./BETA_BENCHMARK_PR_SUMMARY.md)** : Résumé pour Pull Request

---

## Documentation Alpha Chains (référence)

Pour comparer avec le système Alpha :

- **[ALPHA_CHAINS_TECHNICAL_GUIDE.md](./ALPHA_CHAINS_TECHNICAL_GUIDE.md)** : Guide technique Alpha
- **[ALPHA_CHAINS_USER_GUIDE.md](./ALPHA_CHAINS_USER_GUIDE.md)** : Guide utilisateur Alpha
- **[ALPHA_NODE_SHARING.md](./ALPHA_NODE_SHARING.md)** : Concepts Alpha Sharing

---

## Quick Start

### Pour les débutants

1. **Découvrir les concepts** : Commencez par [BETA_NODE_SHARING.md](./BETA_NODE_SHARING.md)
2. **Voir des exemples** : Consultez [BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md) (exemples 1-5)
3. **Exécuter du code** : Testez [examples/beta_chains/](../../examples/beta_chains/)
   ```bash
   cd examples/beta_chains
   go run main.go -scenario simple
   ```

### Pour les développeurs

1. **Migration** : [BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md) - Étapes pas à pas
2. **Exemples pratiques** : [BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md) - 15+ exemples
3. **Guide utilisateur** : [BETA_CHAINS_USER_GUIDE.md](./BETA_CHAINS_USER_GUIDE.md)
4. **Code exécutable** : [examples/beta_chains/](../../examples/beta_chains/) - Scenarios testables

### Pour les experts

1. **Architecture** : [BETA_CHAINS_TECHNICAL_GUIDE.md](./BETA_CHAINS_TECHNICAL_GUIDE.md)
2. **Exemples avancés** : [BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md) (exemples 7-12)
3. **Optimisations** : [BETA_CHAINS_TECHNICAL_GUIDE.md](./BETA_CHAINS_TECHNICAL_GUIDE.md) - Section Optimisations
4. **Benchmarks** : [BETA_PERFORMANCE_REPORT.md](./docs/BETA_PERFORMANCE_REPORT.md)

### Pour les Ops/DevOps

1. **Migration** : [BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md)
2. **Configuration** : [BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md) - Section Configuration
3. **Troubleshooting** : [BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md) - Section Troubleshooting
4. **Monitoring** : [BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md) - Exemple 10

---

## Index par sujet

### Exemples et Tutoriels

- [Exemples basiques (2-3 joins)](./BETA_CHAINS_EXAMPLES.md#exemples-basiques)
- [Exemples de partage](./BETA_CHAINS_EXAMPLES.md#exemples-de-partage)
- [Exemples avancés](./BETA_CHAINS_EXAMPLES.md#exemples-avancés)
- [Cas d'usage réels](./BETA_CHAINS_EXAMPLES.md#cas-dusage-réels)
- [Code exécutable](../../examples/beta_chains/)
- [Visualisations](./BETA_CHAINS_EXAMPLES.md#visualisations)

### Migration et Déploiement

- [Guide de migration complet](./BETA_CHAINS_MIGRATION.md)
- [Impact sur code existant](./BETA_CHAINS_MIGRATION.md#impact-sur-le-code-existant)
- [Migration pas à pas](./BETA_CHAINS_MIGRATION.md#migration-pas-à-pas)
- [Configuration recommandée](./BETA_CHAINS_MIGRATION.md#configuration-et-tuning)
- [Troubleshooting](./BETA_CHAINS_MIGRATION.md#troubleshooting)
- [Procédure de rollback](./BETA_CHAINS_MIGRATION.md#rollback)
- [FAQ Migration (20+ questions)](./BETA_CHAINS_MIGRATION.md#faq-migration)

### Architecture et Conception

- [Architecture détaillée](./BETA_CHAINS_TECHNICAL_GUIDE.md#architecture)
- [Structure des JoinNodes](./BETA_NODE_SHARING.md#concepts-de-base)
- [Différence Alpha/Beta](./BETA_NODE_SHARING.md#différence-avec-alpha-sharing)
- [Registre de partage](./BETA_CHAINS_TECHNICAL_GUIDE.md#betasharingregistry)

### Algorithmes

- [Construction de chaîne](./BETA_CHAINS_TECHNICAL_GUIDE.md#algorithmes)
- [Normalisation des patterns](./BETA_CHAINS_TECHNICAL_GUIDE.md#normalisation-des-patterns-de-jointure)
- [Optimisation de l'ordre de jointure](./BETA_CHAINS_TECHNICAL_GUIDE.md#optimisation-de-lordre-de-jointure)
- [Cache de jointure](./BETA_CHAINS_TECHNICAL_GUIDE.md#algorithme-de-jointure-runtime)

### Configuration et Utilisation

- [Activation du Beta Sharing](./BETA_CHAINS_MIGRATION.md#étape-2--activation-basique-opt-in)
- [Configuration des caches](./BETA_CHAINS_MIGRATION.md#configuration-et-tuning)
- [Configuration par défaut](./BETA_CHAINS_MIGRATION.md#configuration-par-défaut)
- [Configuration haute performance](./BETA_CHAINS_MIGRATION.md#configuration-haute-performance)
- [Configuration mémoire optimisée](./BETA_CHAINS_MIGRATION.md#configuration-mémoire-optimisée)
- [Exemples pratiques](./BETA_CHAINS_EXAMPLES.md)
- [Cas d'usage](./BETA_CHAINS_EXAMPLES.md#cas-dusage-réels)

### Performance

- [Métriques de partage](./BETA_CHAINS_EXAMPLES.md#métriques-de-partage)
- [Optimisation ordre de jointure](./BETA_CHAINS_EXAMPLES.md#exemple-7--optimisation-de-lordre-de-jointure)
- [Cache de jointure](./BETA_CHAINS_EXAMPLES.md#exemple-9--cache-de-jointure-betajoincache)
- [Benchmarks](./BETA_PERFORMANCE_REPORT.md)
- [Benchmarks exécutables](../../examples/beta_chains/README.md#benchmarking)
- [Bénéfices du partage](./BETA_CHAINS_USER_GUIDE.md#bénéfices-du-beta-sharing)
- [Optimisations avancées](./BETA_CHAINS_TECHNICAL_GUIDE.md#optimisations)
- [Métriques Prometheus](./BETA_CHAINS_TECHNICAL_GUIDE.md#métriques-prometheus)

### Dépannage et Maintenance

- [Guide de dépannage](./BETA_CHAINS_USER_GUIDE.md#guide-de-dépannage)
- [Cas edge](./BETA_CHAINS_TECHNICAL_GUIDE.md#gestion-des-cas-edge)
- [FAQ](./BETA_CHAINS_USER_GUIDE.md#faq)
- [Debugging](./BETA_CHAINS_TECHNICAL_GUIDE.md#debugging-et-observabilité)

### API et Référence

- [BetaChainBuilder API](./BETA_CHAINS_TECHNICAL_GUIDE.md#betachainbuilder)
- [BetaSharingRegistry API](./BETA_CHAINS_TECHNICAL_GUIDE.md#betasharingregistry)
- [BetaJoinCache API](./BETA_CHAINS_TECHNICAL_GUIDE.md#betajoincache)
- [BetaChainMetrics API](./BETA_CHAINS_TECHNICAL_GUIDE.md#betachainmetrics)

---

## Diagrammes et Visualisations

### Diagrammes d'architecture

- [Vue d'ensemble du système](./BETA_CHAINS_TECHNICAL_GUIDE.md#vue-densemble-du-système)
- [Anatomie d'un JoinNode](./BETA_NODE_SHARING.md#diagramme-1--anatomie-dun-joinnode)
- [Flux d'activation](./BETA_NODE_SHARING.md#diagramme-2--flux-dactivation)
- [Cycle de vie](./BETA_NODE_SHARING.md#diagramme-3--cycle-de-vie-dun-joinnode-partagé)

### Exemples visuels

- [Visualisations ASCII](./BETA_CHAINS_EXAMPLES.md#visualisations)
- [Diagrammes Mermaid](./BETA_CHAINS_EXAMPLES.md#visualisation-mermaid---chaîne-simple)
- [Comparaison avec/sans partage](./BETA_CHAINS_EXAMPLES.md#visualisation-ascii---comparaison-partage)
- [Réseau sans partage](./BETA_NODE_SHARING.md#exemple-1--réseau-sans-beta-sharing)
- [Réseau avec partage](./BETA_NODE_SHARING.md#exemple-2--réseau-avec-beta-sharing)
- [Partage de préfixes](./BETA_NODE_SHARING.md#exemple-3--chaîne-complexe-avec-préfixes-partagés)
- [Avant/Après optimisation](./BETA_NODE_SHARING.md#exemple-5--avantaprès-optimisation)

---

## Glossaire

### Termes clés

- **Beta Node / JoinNode** : Nœud effectuant une jointure entre deux sources de données
- **Beta Sharing** : Réutilisation de JoinNodes entre plusieurs règles
- **Beta Chain** : Séquence ordonnée de JoinNodes en cascade
- **Left Memory** : Mémoire stockant les tokens (séquences de faits) dans un JoinNode
- **Right Memory** : Mémoire stockant les faits individuels dans un JoinNode
- **Join Test** : Condition de jointure entre un token et un fait
- **Token** : Séquence de faits correspondant aux patterns précédents d'une règle
- **Pattern** : Condition de filtrage sur un type de fait
- **Normalisation** : Transformation d'un pattern en forme canonique
- **RefCount** : Compteur de références (nombre de règles utilisant un nœud)
- **Join Cache** : Cache LRU des résultats d'évaluation de jointure
- **Hash Cache** : Cache LRU des hash de patterns normalisés
- **Sharing Ratio** : Pourcentage de JoinNodes réutilisés vs créés
- **Sélectivité** : Estimation du filtrage d'une jointure (0-1, plus bas = plus sélectif)
- **Prefix Sharing** : Réutilisation de sous-chaînes communes entre règles

### Acronymes

- **LRU** : Least Recently Used (cache)
- **GC** : Garbage Collection
- **API** : Application Programming Interface
- **ASCII** : American Standard Code for Information Interchange
- **FNV** : Fowler-Noll-Vo (algorithme de hash)

---

## Ressources additionnelles

### Documentation externe

- [RETE Algorithm Paper (Forgy, 1982)](https://en.wikipedia.org/wiki/Rete_algorithm)
- [Drools Documentation](https://docs.drools.org/) (autre implémentation RETE)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

### Code source

- **[beta_sharing.go](./beta_sharing.go)** : Implémentation du BetaSharingRegistry
- **[beta_chain_builder.go](./beta_chain_builder.go)** : Builder de chaînes Beta
- **[beta_join_cache.go](./beta_join_cache.go)** : Cache de résultats de jointure
- **[beta_chain_metrics.go](./beta_chain_metrics.go)** : Métriques de performance

### Tests et Benchmarks

- **[beta_sharing_test.go](./beta_sharing_test.go)** : Tests unitaires
- **[beta_sharing_integration_test.go](./beta_sharing_integration_test.go)** : Tests d'intégration
- **[beta_chain_performance_test.go](./beta_chain_performance_test.go)** : Benchmarks de performance

### Exemples exécutables

- **[examples/beta_chains/main.go](../../examples/beta_chains/main.go)** : Exemple interactif
- **[examples/beta_chains/scenarios/simple.go](../../examples/beta_chains/scenarios/simple.go)** : Scénario simple
- **[examples/beta_chains/scenarios/complex.go](../../examples/beta_chains/scenarios/complex.go)** : Scénario complexe
- **[examples/beta_chains/scenarios/advanced.go](../../examples/beta_chains/scenarios/advanced.go)** : Scénario avancé
- **[examples/beta_chains/benchmark_test.go](../../examples/beta_chains/benchmark_test.go)** : Benchmarks comparatifs

---

## Contribuer

### Signaler un problème

Si vous trouvez une erreur dans la documentation ou le code :

1. Vérifiez que le problème n'est pas déjà signalé
2. Testez avec/sans Beta Sharing pour confirmer la cause
3. Créez une issue avec un exemple reproductible
4. Incluez les logs et métriques pertinents
5. Utilisez `go run examples/beta_chains/main.go` pour démontrer

### Améliorer la documentation

Pour contribuer à la documentation :

1. Forkez le dépôt
2. Créez une branche (`git checkout -b improve-beta-docs`)
3. Modifiez la documentation
4. Testez les exemples de code
5. Soumettez une Pull Request

### Guidelines de style

- Utilisez des exemples concrets et exécutables
- Incluez des diagrammes ASCII pour les concepts complexes
- Ajoutez des références croisées entre documents
- Respectez la licence MIT

---

## Historique des versions

### Version 1.3.0 (2024-01)
- ✅ Implémentation complète du Beta Sharing
- ✅ Cache de jointure LRU
- ✅ Optimisation de l'ordre de jointure
- ✅ Métriques Prometheus
- ✅ Documentation complète

### Version 1.2.0 (2023-12)
- Alpha Sharing déjà implémenté
- Base pour le Beta Sharing

---

## Support

### Canaux de communication

- **Issues GitHub** : Pour les bugs et feature requests
- **Discussions** : Pour les questions et discussions générales
- **Pull Requests** : Pour les contributions de code

### Ressources de support

- **[BETA_CHAINS_EXAMPLES.md](./BETA_CHAINS_EXAMPLES.md)** : 15+ exemples pour apprendre
- **[BETA_CHAINS_MIGRATION.md](./BETA_CHAINS_MIGRATION.md)** : Guide de migration et troubleshooting
- **[examples/beta_chains/](../../examples/beta_chains/)** : Code exécutable pour tester
- **[FAQ Migration](./BETA_CHAINS_MIGRATION.md#faq-migration)** : 20+ questions/réponses

### Contact

- Email : support@tsd-project.org (exemple)
- Documentation : [https://docs.tsd-project.org](https://docs.tsd-project.org)

---

## Licence

Copyright (c) 2024 TSD Project

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS