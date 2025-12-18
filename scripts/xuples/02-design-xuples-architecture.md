# Prompt 02 - Conception de l'architecture du module xuples

## 🎯 Objectif

Concevoir l'architecture complète du module xuples en définissant :
- Les interfaces et structures de données
- Les responsabilités de chaque composant
- Les interactions entre RETE et xuples
- Le cycle de vie des xuples
- Les politiques de xuple-space
- L'architecture permettant d'éviter tout hardcoding

Cette conception doit être cohérente avec l'analyse précédente et permettre une implémentation progressive.

## 📋 Tâches

### 1. Définir les structures de données core du module xuples

**Objectif** : Concevoir les structures fondamentales sans hardcoding.

- [ ] Définir la structure `Xuple` (fait principal + faits déclencheurs + métadonnées)
- [ ] Définir la structure `XupleSpace` (nom, politique, stockage)
- [ ] Définir les structures de politique (sélection, consommation, rétention)
- [ ] Définir les métadonnées nécessaires (timestamp, consommations, etc.)
- [ ] Définir les interfaces pour extensibilité

**Livrables** :
- Créer `tsd/docs/xuples/design/01-data-structures.md` contenant :
  - Déclaration Go complète de toutes les structures
  - Justification de chaque champ
  - Relations entre structures (diagramme)
  - Exemples d'instanciation
  - Considérations mémoire et performance

### 2. Concevoir l'interface publique du module xuples

**Objectif** : Définir une API claire et découplée pour le module xuples.

- [ ] Définir l'interface `XupleManager` (création, gestion des xuple-spaces)
- [ ] Définir l'interface `XupleSpace` (insertion, récupération, politiques)
- [ ] Définir l'interface `SelectionPolicy` (stratégie de sélection)
- [ ] Définir l'interface `ConsumptionPolicy` (stratégie de consommation)
- [ ] Définir l'interface `RetentionPolicy` (stratégie de rétention)
- [ ] Définir les types d'erreurs spécifiques

**Livrables** :
- Créer `tsd/docs/xuples/design/02-interfaces.md` contenant :
  - Déclaration complète de toutes les interfaces Go
  - Contrat de chaque méthode (params, retours, erreurs)
  - Principe de responsabilité unique respecté
  - Diagramme d'interfaces
  - Exemples d'utilisation

### 3. Concevoir le système de politiques configurable

**Objectif** : Permettre la configuration des politiques sans hardcoding.

- [ ] Définir les types de politiques de sélection (random, FIFO, LIFO)
- [ ] Définir les types de politiques de consommation (once, per-agent, limited)
- [ ] Définir les types de politiques de rétention (unlimited, duration-based)
- [ ] Concevoir un système de configuration des politiques
- [ ] Prévoir l'extensibilité (nouvelles politiques personnalisées)

**Livrables** :
- Créer `tsd/docs/xuples/design/03-policies.md` contenant :
  - Catalogue complet des politiques par défaut
  - Structure de configuration des politiques
  - Mécanisme d'enregistrement de nouvelles politiques
  - Implémentation de chaque politique (algorithme)
  - Exemples de configuration en TSD
  - Diagramme de stratégie (design pattern)

### 4. Concevoir l'intégration RETE ↔ xuples

**Objectif** : Définir comment RETE et xuples communiquent sans couplage fort.

- [ ] Concevoir l'interface de callback pour l'action Xuple
- [ ] Définir comment extraire les faits déclencheurs d'un token
- [ ] Concevoir le passage de données entre RETE et xuples
- [ ] Définir la gestion d'erreurs entre les deux modules
- [ ] Prévoir la testabilité (injection de dépendances)

**Livrables** :
- Créer `tsd/docs/xuples/design/04-rete-integration.md` contenant :
  - Interface de pont entre RETE et xuples
  - Diagramme de séquence de l'action Xuple
  - Extraction des faits déclencheurs du token
  - Gestion d'erreurs et propagation
  - Injection de dépendances pour tests
  - Exemple de code complet

### 5. Concevoir le cycle de vie des xuples

**Objectif** : Définir précisément comment les xuples naissent, vivent et meurent.

- [ ] Définir les états d'un xuple (créé, disponible, consommé, expiré)
- [ ] Concevoir la création de xuple (action Xuple)
- [ ] Concevoir la consommation par les agents
- [ ] Concevoir l'expiration basée sur le temps
- [ ] Concevoir le nettoyage (garbage collection)
- [ ] Définir les événements et notifications possibles

**Livrables** :
- Créer `tsd/docs/xuples/design/05-lifecycle.md` contenant :
  - Machine à états des xuples
  - Diagramme de cycle de vie complet
  - Algorithmes de gestion d'état
  - Stratégie de nettoyage mémoire
  - Gestion de la concurrence (si applicable)
  - Métriques et observabilité

### 6. Concevoir l'interface agent (future)

**Objectif** : Préparer l'interface pour les agents externes (MVP pour cette phase).

- [ ] Définir le concept d'agent (identité, authentification future)
- [ ] Concevoir l'API minimale pour récupérer un xuple
- [ ] Concevoir l'API pour marquer un xuple consommé
- [ ] Définir les formats d'échange (structure du xuple retourné)
- [ ] Prévoir l'extensibilité pour versions futures

**Livrables** :
- Créer `tsd/docs/xuples/design/06-agent-interface.md` contenant :
  - Définition du concept d'agent
  - API minimale (méthodes Go pour cette phase)
  - Format de sérialisation des xuples
  - Considérations futures (REST API, etc.)
  - Gestion de session agent (simplifié pour MVP)
  - Exemples d'utilisation

### 7. Concevoir l'organisation du package xuples

**Objectif** : Structurer le code du module de manière maintenable.

- [ ] Définir l'arborescence du package `tsd/xuples/`
- [ ] Répartir les responsabilités entre sous-packages
- [ ] Définir les exports publics vs privés
- [ ] Concevoir la stratégie de tests (unitaires, intégration)
- [ ] Prévoir les fichiers de documentation

**Livrables** :
- Créer `tsd/docs/xuples/design/07-package-structure.md` contenant :
  - Arborescence complète du package
  - Responsabilité de chaque fichier
  - Exports publics (API du module)
  - Dépendances entre fichiers
  - Stratégie de tests par composant
  - Convention de nommage

Exemple attendu :
```
tsd/xuples/
├── xuples.go              # Types publics et XupleManager
├── xuplespace.go          # Implémentation XupleSpace
├── policies.go            # Politiques de base
├── policy_selection.go    # Implémentations SelectionPolicy
├── policy_consumption.go  # Implémentations ConsumptionPolicy
├── policy_retention.go    # Implémentations RetentionPolicy
├── lifecycle.go           # Gestion cycle de vie
├── errors.go              # Erreurs spécifiques
├── xuples_test.go
├── xuplespace_test.go
└── testdata/
```

### 8. Créer le document de conception complet

**Objectif** : Synthétiser toute la conception en un document maître.

- [ ] Vue d'ensemble de l'architecture
- [ ] Diagramme de composants
- [ ] Diagramme de classes
- [ ] Flux de données complets
- [ ] Décisions architecturales et justifications
- [ ] Plan d'implémentation recommandé

**Livrables** :
- Créer `tsd/docs/xuples/design/00-INDEX.md` contenant :
  - Architecture complète du module xuples
  - Diagrammes de haut niveau (composants, classes)
  - Principes de conception (SOLID, découplage)
  - Décisions architecturales et alternatives considérées
  - Limitations connues et évolutions futures
  - Roadmap d'implémentation
  - Matrice de traçabilité (exigences → conception)

## 📁 Structure de documentation attendue

```
tsd/docs/xuples/
├── analysis/                            # (créé au prompt 01)
└── design/
    ├── 00-INDEX.md                      # Vue d'ensemble conception
    ├── 01-data-structures.md            # Structures de données
    ├── 02-interfaces.md                 # Interfaces publiques
    ├── 03-policies.md                   # Système de politiques
    ├── 04-rete-integration.md           # Intégration RETE
    ├── 05-lifecycle.md                  # Cycle de vie xuples
    ├── 06-agent-interface.md            # Interface agents
    └── 07-package-structure.md          # Organisation code
```

## ✅ Critères de succès

- [ ] Architecture complète et cohérente
- [ ] Tous les composants clairement définis
- [ ] Interfaces respectant SOLID
- [ ] Aucun hardcoding dans la conception
- [ ] Extensibilité garantie (nouvelles politiques, etc.)
- [ ] Découplage fort entre RETE et xuples
- [ ] Testabilité garantie (injection de dépendances)
- [ ] Documentation complète et claire
- [ ] Diagrammes UML/architecture fournis
- [ ] Prêt pour implémentation

## 🎨 Principes de conception à respecter

### Obligatoires (selon common.md)

- **Single Responsibility** - Chaque composant une seule responsabilité
- **Open/Closed** - Extensible sans modification (nouvelles politiques)
- **Dependency Injection** - Pas de dépendances globales hardcodées
- **Composition over Inheritance** - Interfaces et embedding
- **Interfaces** - Petites, focalisées, cohérentes
- **Découplage fort** - RETE et xuples indépendants

### Spécifiques au module xuples

- **Configuration over Code** - Politiques configurables, pas hardcodées
- **Policy Pattern** - Stratégies interchangeables
- **Factory Pattern** - Création de xuples et xuple-spaces
- **Observer Pattern** - Notifications d'événements (optionnel)
- **Strategy Pattern** - Implémentation des politiques

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/analysis/` - Analyse de l'existant (prompt 01)
- Effective Go - https://go.dev/doc/effective_go
- Go Design Patterns
- SOLID Principles

## 🎯 Prochaine étape

Une fois cette conception terminée et validée, passer au prompt **03-extend-parser-xuplespace.md** pour ajouter le parsing de la commande `xuple-space` dans le langage TSD.