# Index de la Documentation TSD

Guide de navigation dans la documentation du projet TSD (Type System Development).

## 📋 Documentation Principale

### À la Racine
- [**README.md**](README.md) - Vue d'ensemble du projet, installation, quick start
- [**CHANGELOG.md**](CHANGELOG.md) - Historique des versions et changements
- [**LOGGING_GUIDE.md**](LOGGING_GUIDE.md) - Guide du système de logging
- [**THIRD_PARTY_LICENSES.md**](THIRD_PARTY_LICENSES.md) - Licences des dépendances

### Documentation Générale (`docs/`)
- [**README**](docs/README.md) - Index de la documentation
- [**TUTORIAL**](docs/TUTORIAL.md) - Tutoriel complet de zéro à héros
- [**FEATURES**](docs/FEATURES.md) - Toutes les fonctionnalités du projet
- [**OPTIMIZATIONS**](docs/OPTIMIZATIONS.md) - Guide complet des optimisations
- [**API_REFERENCE**](docs/API_REFERENCE.md) - Référence complète de l'API
- [**EXAMPLES**](docs/EXAMPLES.md) - Exemples d'utilisation
- [**GRAMMAR_GUIDE**](docs/GRAMMAR_GUIDE.md) - Guide de la grammaire TSD

### Guides Avancés (`docs/`)
- [**STRONG_MODE_TUNING_GUIDE**](docs/STRONG_MODE_TUNING_GUIDE.md) - Configuration du Strong Mode
- [**PROMETHEUS_INTEGRATION**](docs/PROMETHEUS_INTEGRATION.md) - Intégration Prometheus
- [**TRANSACTION_README**](docs/TRANSACTION_README.md) - Système de transactions
- [**TRANSACTION_ARCHITECTURE**](docs/TRANSACTION_ARCHITECTURE.md) - Architecture des transactions
- [**development_guidelines**](docs/development_guidelines.md) - Guidelines de développement

---

## 🔧 Module Constraint (`constraint/`)

Parser et validateur pour le langage TSD.

### Documentation
- [**README**](constraint/README.md) - Vue d'ensemble du module
- [**docs/README**](constraint/docs/README.md) - Index de la documentation
- [**GUIDE_CONTRAINTES**](constraint/docs/GUIDE_CONTRAINTES.md) - Guide complet des contraintes
- [**TUTORIEL_ACTIONS**](constraint/docs/TUTORIEL_ACTIONS.md) - Tutoriel des actions
- [**TUTORIEL_CONTRAINTES**](constraint/docs/TUTORIEL_CONTRAINTES.md) - Tutoriel des contraintes
- [**GRAMMAR_COMPLETE**](constraint/docs/GRAMMAR_COMPLETE.md) - Grammaire PEG complète
- [**TYPE_VALIDATION**](constraint/docs/TYPE_VALIDATION.md) - Validation des types

### Tests
- [**test/validation/README**](constraint/test/validation/README.md) - Tests de validation
- [**test/integration/NEGATION_RESULTS_COMPLETE**](constraint/test/integration/NEGATION_RESULTS_COMPLETE.md) - Résultats tests de négation

---

## ⚙️ Module RETE (`rete/`)

Moteur d'inférence basé sur l'algorithme RETE.

### Documentation Principale
- [**README**](rete/README.md) - Vue d'ensemble du module RETE

### Guides Utilisateur (`rete/docs/`)

#### Fonctionnalités de Base
- [**ACTIONS**](rete/docs/ACTIONS.md) - Guide des actions
- [**ARITHMETIC**](rete/docs/ARITHMETIC.md) - Expressions arithmétiques
- [**NESTED_OR**](rete/docs/NESTED_OR.md) - Expressions OR imbriquées

#### Optimisations
- [**ALPHA_CHAINS**](rete/docs/ALPHA_CHAINS.md) - Chaînes d'alpha nodes
- [**BETA_SHARING**](rete/docs/BETA_SHARING.md) - Partage des jointures
- [**BETA_CHAINS**](rete/docs/BETA_CHAINS.md) - Chaînes de jointures
- [**OPTIMIZATIONS**](rete/docs/OPTIMIZATIONS.md) - Guide des optimisations

#### Fonctionnalités Avancées
- [**MULTI_SOURCE_AGGREGATION**](rete/docs/MULTI_SOURCE_AGGREGATION.md) - Agrégations multi-sources
- [**NODE_LIFECYCLE**](rete/docs/NODE_LIFECYCLE.md) - Cycle de vie des nœuds
- [**NORMALIZATION**](rete/docs/NORMALIZATION.md) - Normalisation des expressions

### Guides Techniques (`rete/docs/`)

#### Architecture
- [**ADVANCED_NODES_IMPLEMENTATION**](rete/docs/ADVANCED_NODES_IMPLEMENTATION.md) - Implémentation nœuds avancés
- [**ADVANCED_NODES_USAGE_GUIDE**](rete/docs/ADVANCED_NODES_USAGE_GUIDE.md) - Utilisation nœuds avancés
- [**TUPLE_SPACE_IMPLEMENTATION**](rete/docs/TUPLE_SPACE_IMPLEMENTATION.md) - Espace de tuples

#### Beta Nodes
- [**BETA_NODES_GUIDE**](rete/docs/BETA_NODES_GUIDE.md) - Guide des nœuds beta
- [**BETA_NODES_ANALYSIS**](rete/docs/BETA_NODES_ANALYSIS.md) - Analyse des nœuds beta
- [**BETA_NODES_ARCHITECTURE_DIAGRAMS**](rete/docs/BETA_NODES_ARCHITECTURE_DIAGRAMS.md) - Diagrammes d'architecture

#### Conception
- [**ALPHA_BETA_CHAINS_COMPARISON**](rete/docs/ALPHA_BETA_CHAINS_COMPARISON.md) - Comparaison alpha/beta chains
- [**BETA_CHAINS_DESIGN**](rete/docs/BETA_CHAINS_DESIGN.md) - Conception beta chains
- [**BETA_SHARING_DESIGN**](rete/docs/BETA_SHARING_DESIGN.md) - Conception beta sharing

#### Exemples
- [**BETA_CHAINS_EXAMPLES**](rete/docs/BETA_CHAINS_EXAMPLES.md) - Exemples beta chains
- [**BETA_SHARING_EXAMPLES**](rete/docs/BETA_SHARING_EXAMPLES.md) - Exemples beta sharing

#### Outils
- [**EXPRESSION_ANALYZER_README**](rete/docs/EXPRESSION_ANALYZER_README.md) - Analyseur d'expressions
- [**TESTING**](rete/docs/TESTING.md) - Guide des tests

#### Features
- [**FEATURE_ARITHMETIC_ALPHA_NODES**](rete/docs/FEATURE_ARITHMETIC_ALPHA_NODES.md) - Alpha nodes arithmétiques
- [**FEATURE_PASSTHROUGH_PER_RULE**](rete/docs/FEATURE_PASSTHROUGH_PER_RULE.md) - Optimisation passthrough

---

## 📚 Exemples (`examples/`)

- [**examples/README**](examples/README.md) - Index des exemples
- [**beta_chains/README**](examples/beta_chains/README.md) - Exemple beta chains
- [**lru_cache/README**](examples/lru_cache/README.md) - Exemple cache LRU
- [**multi_source_aggregations/README**](examples/multi_source_aggregations/README.md) - Exemple agrégations

---

## 🧪 Tests (`test/` et `tests/`)

- [**test/README**](test/README.md) - Guide des tests (nouveau système)
- [**tests/README**](tests/README.md) - Guide des tests (ancien système)

---

## 🚀 Quick Start par Use Case

### Je veux apprendre TSD
1. [README.md](README.md) - Vue d'ensemble
2. [docs/TUTORIAL.md](docs/TUTORIAL.md) - Tutoriel complet
3. [docs/EXAMPLES.md](docs/EXAMPLES.md) - Exemples pratiques

### Je veux optimiser les performances
1. [docs/OPTIMIZATIONS.md](docs/OPTIMIZATIONS.md) - Guide des optimisations
2. [rete/docs/ALPHA_CHAINS.md](rete/docs/ALPHA_CHAINS.md) - Alpha chains
3. [rete/docs/BETA_SHARING.md](rete/docs/BETA_SHARING.md) - Beta sharing
4. [docs/STRONG_MODE_TUNING_GUIDE.md](docs/STRONG_MODE_TUNING_GUIDE.md) - Strong mode

### Je veux comprendre le moteur RETE
1. [rete/README.md](rete/README.md) - Vue d'ensemble RETE
2. [rete/docs/ADVANCED_NODES_IMPLEMENTATION.md](rete/docs/ADVANCED_NODES_IMPLEMENTATION.md) - Nœuds avancés
3. [rete/docs/BETA_NODES_GUIDE.md](rete/docs/BETA_NODES_GUIDE.md) - Nœuds beta

### Je veux écrire des règles complexes
1. [constraint/docs/GUIDE_CONTRAINTES.md](constraint/docs/GUIDE_CONTRAINTES.md) - Guide contraintes
2. [constraint/docs/TUTORIEL_ACTIONS.md](constraint/docs/TUTORIEL_ACTIONS.md) - Actions
3. [rete/docs/ARITHMETIC.md](rete/docs/ARITHMETIC.md) - Arithmétique
4. [rete/docs/NESTED_OR.md](rete/docs/NESTED_OR.md) - Expressions OR

### Je veux contribuer au projet
1. [docs/development_guidelines.md](docs/development_guidelines.md) - Guidelines
2. [rete/docs/TESTING.md](rete/docs/TESTING.md) - Tests
3. [CHANGELOG.md](CHANGELOG.md) - Historique des changements

---

## 📞 Support et Ressources

- **Documentation complète** : Voir [docs/README.md](docs/README.md)
- **API Reference** : [docs/API_REFERENCE.md](docs/API_REFERENCE.md)
- **Issues GitHub** : Reporter bugs et demander features
- **Examples** : Exemples concrets dans `examples/`

---

**Version** : 1.0  
**Dernière mise à jour** : Janvier 2025  
**Total de fichiers de documentation** : ~60 fichiers organisés