# Prompt 01 - Analyse de l'implémentation actuelle des actions et terminal nodes

## 🎯 Objectif

Analyser en profondeur l'implémentation actuelle du système d'actions et des terminal nodes dans TSD pour comprendre :
- Comment les actions sont actuellement définies et parsées
- Comment les terminal nodes stockent les tokens (comportement tuple-space actuel)
- Comment les actions sont (ou seraient) exécutées
- Les structures de données impliquées (Token, TerminalNode, ActionExecutor, etc.)
- Les points d'intégration nécessaires pour la refonte

Cette analyse servira de base pour la conception de la nouvelle architecture séparant RETE et xuples.

## 📋 Tâches

### 1. Analyser la définition et le parsing des actions

**Objectif** : Comprendre comment les actions sont actuellement définies dans le langage TSD et parsées.

- [ ] Localiser le code de parsing des actions (grammaire, parser)
- [ ] Identifier la structure AST des actions
- [ ] Documenter les types d'actions supportés actuellement
- [ ] Identifier où sont stockées les définitions d'actions parsées
- [ ] Vérifier les validations existantes (détection de doublons, etc.)

**Livrables** :
- Créer `tsd/docs/xuples/analysis/01-current-action-parsing.md` documentant :
  - La grammaire complète des actions
  - Les structures de données utilisées
  - Le flux de parsing (fichier TSD → AST → structures internes)
  - Les validations existantes
  - Extraits de code pertinents avec références de fichiers et lignes

### 2. Analyser l'implémentation des Terminal Nodes

**Objectif** : Comprendre le comportement actuel des terminal nodes qui stockent les tokens.

- [ ] Examiner `tsd/rete/node_terminal.go` en détail
- [ ] Analyser la méthode `ActivateLeft` et le stockage des tokens
- [ ] Analyser la méthode `executeAction` et son comportement actuel
- [ ] Identifier la structure `WorkingMemory` et `Token`
- [ ] Comprendre comment `collectActivations` récupère les tokens

**Livrables** :
- Créer `tsd/docs/xuples/analysis/02-terminal-nodes.md` documentant :
  - Architecture actuelle des terminal nodes
  - Cycle de vie d'un token (création → stockage → récupération)
  - Structure de données `Token` et `TokenMetadata`
  - Interface `ActionExecutor` actuelle
  - Diagramme de séquence du flux actuel
  - Points d'intervention pour la refonte

### 3. Analyser les structures de données Token et Fact

**Objectif** : Comprendre en détail les structures qui portent les données dans RETE.

- [ ] Examiner `tsd/rete/fact_token.go`
- [ ] Documenter la structure `Token` complète
- [ ] Documenter la structure `Fact` complète
- [ ] Identifier comment les faits déclencheurs sont liés aux tokens
- [ ] Comprendre les métadonnées associées

**Livrables** :
- Créer `tsd/docs/xuples/analysis/03-token-fact-structures.md` documentant :
  - Structures complètes avec tous les champs
  - Relations entre Token, Fact, et WorkingMemory
  - Comment un token "combiné" contient plusieurs faits déclencheurs
  - Schéma des structures de données
  - Implications pour la création de xuples

### 4. Analyser l'ActionExecutor et son interface

**Objectif** : Comprendre comment les actions sont censées être exécutées.

- [ ] Localiser l'interface `ActionExecutor`
- [ ] Identifier les implémentations existantes (si présentes)
- [ ] Comprendre le contrat d'exécution des actions
- [ ] Analyser comment le réseau RETE délègue l'exécution

**Livrables** :
- Créer `tsd/docs/xuples/analysis/04-action-executor.md` documentant :
  - Interface complète avec toutes les méthodes
  - Implémentations actuelles
  - Contrat d'exécution (paramètres, retours, gestion d'erreurs)
  - Intégration avec le réseau RETE
  - Proposition d'évolution pour les actions par défaut

### 5. Analyser les tests existants

**Objectif** : Comprendre comment les actions et terminal nodes sont testés actuellement.

- [ ] Identifier tous les tests liés aux actions
- [ ] Identifier tous les tests liés aux terminal nodes
- [ ] Analyser les patterns de test utilisés
- [ ] Comprendre comment les tests vérifient les activations

**Livrables** :
- Créer `tsd/docs/xuples/analysis/05-existing-tests.md` documentant :
  - Liste exhaustive des fichiers de test pertinents
  - Patterns de test utilisés
  - Couverture actuelle des fonctionnalités
  - Points faibles de la couverture de test
  - Recommendations pour les nouveaux tests

### 6. Créer un document de synthèse

**Objectif** : Consolider toutes les analyses en une vue d'ensemble.

- [ ] Synthétiser les découvertes
- [ ] Identifier les points de modification nécessaires
- [ ] Proposer une stratégie de migration progressive
- [ ] Lister les risques et contraintes

**Livrables** :
- Créer `tsd/docs/xuples/analysis/00-INDEX.md` contenant :
  - Vue d'ensemble de l'architecture actuelle
  - Cartographie complète du code concerné
  - Points d'intervention identifiés pour la refonte
  - Dépendances entre composants
  - Stratégie de migration recommandée
  - Risques identifiés et mitigations proposées

## 📁 Structure de documentation attendue

```
tsd/docs/xuples/
└── analysis/
    ├── 00-INDEX.md                      # Synthèse et vue d'ensemble
    ├── 01-current-action-parsing.md     # Parsing des actions
    ├── 02-terminal-nodes.md             # Terminal nodes
    ├── 03-token-fact-structures.md      # Structures de données
    ├── 04-action-executor.md            # ActionExecutor
    └── 05-existing-tests.md             # Tests existants
```

## ✅ Critères de succès

- [ ] Tous les documents d'analyse créés et complets
- [ ] Architecture actuelle complètement comprise et documentée
- [ ] Points d'intervention clairement identifiés
- [ ] Aucune ambiguïté sur le fonctionnement actuel
- [ ] Stratégie de migration définie
- [ ] Base solide pour la conception de la nouvelle architecture

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/rete/docs/TUPLE_SPACE_IMPLEMENTATION.md` - Documentation actuelle tuple-space
- `tsd/rete/node_terminal.go` - Implémentation terminal nodes
- `tsd/rete/fact_token.go` - Structures Token et Fact
- `tsd/internal/servercmd/servercmd.go` - Utilisation de collectActivations

## 🎯 Prochaine étape

Une fois cette analyse terminée, passer au prompt **02-design-xuples-architecture.md** pour concevoir la nouvelle architecture du module xuples.