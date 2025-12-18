# Prompt 09 - Finalisation de la documentation du système xuples

## 🎯 Objectif

Finaliser et compléter toute la documentation du système xuples pour garantir que :
- Les utilisateurs peuvent comprendre et utiliser le système
- Les développeurs peuvent maintenir et étendre le code
- L'architecture et les décisions sont documentées
- Des exemples complets sont fournis
- Les guides de migration sont disponibles

## 📋 Tâches

### 1. Créer le guide utilisateur complet

**Objectif** : Document principal pour les utilisateurs du système xuples.

**Fichier à créer** : `tsd/docs/xuples/user-guide/complete-guide.md`

**Contenu attendu** :

```markdown
# Guide Complet du Système Xuples

## Table des matières

1. Introduction
2. Concepts fondamentaux
3. Déclaration de xuple-spaces
4. Utilisation de l'action Xuple
5. Politiques en détail
6. Cas d'usage
7. Bonnes pratiques
8. Dépannage

## 1. Introduction

### Qu'est-ce qu'un xuple ?

Un **xuple** (tuple étendu) est une structure de données qui contient :
- Un **fait principal** : le fait passé à l'action Xuple
- Des **faits déclencheurs** : tous les faits qui ont activé la règle

Les xuples permettent de créer des espaces de coordination entre le moteur de règles RETE et des agents externes.

### Architecture

[Diagramme architecture RETE ↔ xuples]

## 2. Concepts fondamentaux

### Xuple-space

Un xuple-space est un espace nommé qui :
- Stocke des xuples
- Applique des politiques de sélection, consommation et rétention
- Permet aux agents d'accéder aux xuples

### Politiques

Trois types de politiques configurent un xuple-space :

1. **Sélection** : comment choisir un xuple parmi plusieurs
2. **Consommation** : comment les xuples peuvent être consommés
3. **Rétention** : combien de temps les xuples sont conservés

## 3. Déclaration de xuple-spaces

### Syntaxe

\```tsd
xuple-space <nom> {
    selection: <random|fifo|lifo>
    consumption: <once|per-agent|limited(n)>
    retention: <unlimited|duration(temps)>
}
\```

### Exemples

#### Xuple-space basique

\```tsd
xuple-space tasks {
    selection: fifo
    consumption: once
    retention: unlimited
}
\```

#### Avec expiration temporelle

\```tsd
xuple-space notifications {
    selection: random
    consumption: per-agent
    retention: duration(5m)
}
\```

#### Avec consommation limitée

\```tsd
xuple-space shared-data {
    selection: lifo
    consumption: limited(3)
    retention: duration(1h)
}
\```

## 4. Utilisation de l'action Xuple

### Syntaxe

\```tsd
Xuple("<nom-xuple-space>", <fait>)
\```

### Exemple complet

\```tsd
xuple-space alerts {
    selection: fifo
    consumption: once
    retention: unlimited
}

fact Alert(level: string, message: string, timestamp: int)

rule "critical-alert" {
    when {
        alert: Alert(level == "critical")
    }
    then {
        Print("Critical alert: " + alert.message)
        Log("Alert logged at " + alert.timestamp)
        Xuple("alerts", alert)
    }
}
\```

### Faits déclencheurs

Le xuple créé contient automatiquement tous les faits qui ont déclenché la règle :

\```tsd
xuple-space assignments {
    selection: fifo
    consumption: once
    retention: unlimited
}

fact Person(name: string, age: int)
fact Department(name: string)
fact Assignment(person: string, dept: string)

rule "valid-assignment" {
    when {
        p: Person(age >= 18)
        d: Department(name == "Engineering")
        a: Assignment(person == p.name, dept == d.name)
    }
    then {
        Xuple("assignments", a)  // Le xuple contiendra p, d, et a
    }
}
\```

## 5. Politiques en détail

### Politiques de sélection

#### random
Sélectionne un xuple aléatoirement parmi ceux disponibles.

**Cas d'usage** : Distribution équitable de charge, éviter les biais d'ordre.

#### fifo (First In, First Out)
Sélectionne le xuple le plus ancien.

**Cas d'usage** : Files d'attente, traitement séquentiel, ordre chronologique important.

#### lifo (Last In, First Out)
Sélectionne le xuple le plus récent.

**Cas d'usage** : Traitement en pile, priorisation des événements récents.

### Politiques de consommation

#### once
Un xuple ne peut être consommé qu'une seule fois au total.

**Cas d'usage** : Tâches uniques, événements one-shot.

\```tsd
xuple-space unique-tasks {
    selection: fifo
    consumption: once
    retention: unlimited
}
\```

#### per-agent
Chaque agent peut consommer le xuple une fois.

**Cas d'usage** : Broadcasting, notification à plusieurs agents.

\```tsd
xuple-space broadcasts {
    selection: random
    consumption: per-agent
    retention: unlimited
}
\```

#### limited(n)
Le xuple peut être consommé maximum n fois.

**Cas d'usage** : Ressources partagées limitées, quotas.

\```tsd
xuple-space shared-resources {
    selection: fifo
    consumption: limited(5)
    retention: unlimited
}
\```

### Politiques de rétention

#### unlimited
Les xuples ne expirent jamais.

**Cas d'usage** : Données persistantes, historique.

#### duration(temps)
Les xuples expirent après la durée spécifiée.

**Format** : `<nombre><unité>` où unité = s (secondes), m (minutes), h (heures), d (jours)

**Exemples** :
- `duration(30s)` : 30 secondes
- `duration(5m)` : 5 minutes
- `duration(2h)` : 2 heures
- `duration(7d)` : 7 jours

**Cas d'usage** : Données temporaires, caches, événements éphémères.

\```tsd
xuple-space temporary-cache {
    selection: fifo
    consumption: per-agent
    retention: duration(10m)
}
\```

## 6. Cas d'usage

### Workflow orchestration

\```tsd
xuple-space workflow-tasks {
    selection: fifo
    consumption: once
    retention: unlimited
}

fact WorkflowStep(id: string, workflow: string, order: int)
fact WorkflowContext(workflow: string, status: string)

rule "queue-workflow-step" {
    when {
        ctx: WorkflowContext(status == "active")
        step: WorkflowStep(workflow == ctx.workflow)
    }
    then {
        Xuple("workflow-tasks", step)
    }
}
\```

### Event broadcasting

\```tsd
xuple-space system-events {
    selection: random
    consumption: per-agent
    retention: duration(1h)
}

fact SystemEvent(type: string, data: string)

rule "broadcast-event" {
    when {
        event: SystemEvent(type == "config-changed")
    }
    then {
        Xuple("system-events", event)
    }
}
\```

### Resource allocation

\```tsd
xuple-space available-slots {
    selection: fifo
    consumption: limited(10)
    retention: unlimited
}

fact ResourceSlot(id: string, capacity: int)
fact Request(requester: string, amount: int)

rule "allocate-slot" {
    when {
        slot: ResourceSlot(capacity > 0)
        req: Request(amount <= slot.capacity)
    }
    then {
        Xuple("available-slots", slot)
    }
}
\```

## 7. Bonnes pratiques

### Nommage des xuple-spaces

- Utilisez des noms descriptifs : `user-notifications`, `pending-tasks`
- Évitez les noms génériques : `data`, `items`
- Soyez cohérent avec votre domaine métier

### Choix des politiques

1. **Commencez simple** : `fifo` / `once` / `unlimited` pour la plupart des cas
2. **Ajoutez de la complexité si nécessaire** :
   - `per-agent` si plusieurs consommateurs doivent voir le même xuple
   - `duration` si les données deviennent obsolètes
   - `limited` si vous avez des quotas

### Performance

- Les xuples avec rétention illimitée s'accumulent : nettoyez-les ou utilisez `duration`
- Pour beaucoup de xuples, préférez `fifo` ou `lifo` (plus rapides que `random`)
- Appelez périodiquement `Cleanup()` sur les xuple-spaces avec rétention temporelle

### Debugging

Activez le logging pour voir les actions :

\```tsd
rule "debug-rule" {
    when {
        event: Event()
    }
    then {
        Log("Creating xuple for event: " + event.id)
        Xuple("myspace", event)
        Log("Xuple created successfully")
    }
}
\```

## 8. Dépannage

### Erreur : "xuple-space not found"

**Cause** : Le xuple-space n'a pas été déclaré.

**Solution** : Ajoutez une déclaration `xuple-space` avant de l'utiliser dans une règle.

### Erreur : "cannot redefine default action 'Xuple'"

**Cause** : Tentative de redéfinir l'action Xuple.

**Solution** : Supprimez la déclaration `action Xuple(...)`.

### Aucun xuple n'est créé

**Causes possibles** :
1. La règle ne se déclenche pas (vérifiez les conditions)
2. Le xuple-space n'existe pas
3. Erreur dans l'action (vérifiez les logs)

**Solution** : Ajoutez du logging avant et après l'action Xuple.

### Les xuples ne sont pas disponibles

**Causes possibles** :
1. Tous consommés (politique `once`)
2. Expirés (politique `duration`)
3. Limite atteinte (politique `limited`)

**Solution** : Vérifiez les politiques et l'état du xuple-space.
```

**Livrables** :
- [ ] Guide utilisateur complet créé
- [ ] Tous les concepts expliqués
- [ ] Exemples pour chaque cas d'usage
- [ ] Section dépannage complète

### 2. Créer la documentation d'architecture

**Objectif** : Documenter l'architecture technique pour les développeurs.

**Fichier à créer** : `tsd/docs/xuples/architecture/overview.md`

**Contenu attendu** :

```markdown
# Architecture du Système Xuples

## Vue d'ensemble

Le système xuples est composé de plusieurs modules indépendants :

1. **Parser** : Parsing de la commande `xuple-space`
2. **Compiler** : Validation et instanciation des xuple-spaces
3. **Module xuples** : Gestion des xuples et politiques
4. **RETE Actions** : Intégration avec le moteur de règles
5. **Default Actions** : Système d'actions par défaut

## Diagramme de composants

[Diagramme UML des composants]

## Séparation RETE ↔ Xuples

Le système est conçu pour maintenir un découplage fort entre :

- **RETE** : Moteur de règles, évaluation de conditions, propagation
- **Xuples** : Système de coordination, stockage, politiques

### Points d'intégration

1. **Action Xuple** : Interface entre RETE et xuples
2. **XupleManager** : Injecté dans BuiltinActionExecutor
3. **Token** : Extraction des faits déclencheurs

## Flux de données

### 1. Compilation

\```
Programme TSD
    ↓
Parser (grammar.peg)
    ↓
AST (XupleSpaceDeclaration)
    ↓
Compiler (validation)
    ↓
XupleManager (instanciation)
    ↓
XupleSpace (avec politiques)
\```

### 2. Exécution

\```
Fait inséré dans RETE
    ↓
Règle activée
    ↓
Action Xuple invoquée
    ↓
BuiltinActionExecutor.executeXuple()
    ↓
Extraction faits déclencheurs (Token)
    ↓
XupleManager.CreateXuple()
    ↓
XupleSpace.Insert()
    ↓
Application politiques
    ↓
Xuple stocké
\```

### 3. Consommation (futur)

\```
Agent demande un xuple
    ↓
XupleSpace.Retrieve(agentID)
    ↓
SelectionPolicy.Select()
    ↓
ConsumptionPolicy.CanConsume()
    ↓
Xuple retourné
    ↓
XupleSpace.MarkConsumed()
    ↓
ConsumptionPolicy.OnConsumed()
    ↓
RetentionPolicy.ShouldRetain()
\```

## Modules en détail

### Module xuples

**Responsabilités** :
- Gestion des xuples et métadonnées
- Implémentation des xuple-spaces
- Implémentation des politiques
- Thread-safety (sync.RWMutex)

**Exports publics** :
- `Xuple`, `XupleMetadata`, `XupleState`
- `XupleManager`, `XupleSpace` (interfaces)
- `XupleSpaceConfig`
- `SelectionPolicy`, `ConsumptionPolicy`, `RetentionPolicy` (interfaces)
- Implémentations de politiques (New*Policy)
- Erreurs (`Err*`)

**Interne** :
- `DefaultXupleManager`, `DefaultXupleSpace` (implémentations)

### Module rete/actions

**Responsabilités** :
- Exécution des actions par défaut
- Intégration avec XupleManager
- Extraction des faits déclencheurs

**Exports publics** :
- `BuiltinActionExecutor`

### Module internal/defaultactions

**Responsabilités** :
- Chargement des définitions d'actions par défaut
- Fichier defaults.tsd embarqué

**Exports publics** :
- `LoadDefaultActions()`
- `IsDefaultAction()`
- `DefaultActionNames`

## Décisions architecturales

### 1. Découplage RETE ↔ Xuples

**Décision** : Les modules sont totalement indépendants.

**Raison** :
- Maintenabilité : modifications isolées
- Testabilité : tests unitaires indépendants
- Réutilisabilité : xuples peut être utilisé ailleurs

**Alternative rejetée** : Xuples intégré dans RETE (couplage fort).

### 2. Injection de dépendances

**Décision** : XupleManager injecté dans BuiltinActionExecutor.

**Raison** :
- Testabilité : mocks faciles
- Flexibilité : implémentations alternatives
- Pas de dépendances globales

**Alternative rejetée** : XupleManager global/singleton (non testable).

### 3. Actions par défaut via fichier

**Décision** : Actions définies dans defaults.tsd, parsé à l'init.

**Raison** :
- Pas de hardcoding
- Cohérence avec le langage TSD
- Facilité de modification
- Vérification au compile-time (du fichier defaults.tsd)

**Alternative rejetée** : Actions hardcodées dans le compilateur.

### 4. Politiques en interfaces

**Décision** : Politiques définies par des interfaces.

**Raison** :
- Extensibilité : nouvelles politiques faciles
- Strategy pattern : comportements interchangeables
- Testabilité : mocks de politiques

**Alternative rejetée** : Enum avec switch (non extensible).

### 5. Thread-safety

**Décision** : sync.RWMutex dans XupleSpace et XupleManager.

**Raison** :
- Accès concurrent sûr
- Performance acceptable (read lock partagé)
- Simplicité d'implémentation

**Alternative rejetée** : Channels (complexité inutile pour ce cas).

## Patterns utilisés

- **Strategy Pattern** : Politiques
- **Factory Pattern** : New*Policy(), NewXupleSpace()
- **Dependency Injection** : XupleManager → BuiltinActionExecutor
- **Observer Pattern** : ActionObserver
- **Embedded Resources** : go:embed defaults.tsd

## Évolutions futures

1. **API REST** pour agents externes
2. **Persistance** des xuples (optionnelle)
3. **Politiques personnalisées** via plugins
4. **Métriques** et observabilité avancée
5. **Clustering** pour distribution
```

**Livrables** :
- [ ] Documentation d'architecture complète
- [ ] Diagrammes de composants
- [ ] Flux de données documentés
- [ ] Décisions architecturales justifiées

### 3. Créer le guide de contribution

**Objectif** : Aider les contributeurs à étendre le système.

**Fichier à créer** : `tsd/docs/xuples/contributing/extending-xuples.md`

**Contenu attendu** :

```markdown
# Guide de Contribution - Système Xuples

## Ajouter une nouvelle politique de sélection

1. Créer l'implémentation dans `tsd/xuples/policy_selection.go` :

\```go
type CustomSelectionPolicy struct {
    // Vos champs
}

func NewCustomSelectionPolicy() *CustomSelectionPolicy {
    return &CustomSelectionPolicy{}
}

func (p *CustomSelectionPolicy) Select(xuples []*Xuple) *Xuple {
    // Votre logique
}

func (p *CustomSelectionPolicy) Name() string {
    return "custom"
}
\```

2. Ajouter les tests dans `tsd/xuples/policies_test.go`

3. Étendre le parser pour supporter la nouvelle politique

4. Mettre à jour la documentation utilisateur

## Ajouter une nouvelle politique de consommation

Même processus dans `policy_consumption.go`.

## Ajouter une nouvelle politique de rétention

Même processus dans `policy_retention.go`.

## Ajouter une action par défaut

1. Ajouter la signature dans `tsd/internal/defaultactions/defaults.tsd`

2. Implémenter dans `tsd/rete/actions/builtin.go` :

\```go
func (e *BuiltinActionExecutor) executeMyAction(args []interface{}, token *rete.Token) error {
    // Validation des arguments
    // Implémentation
}
\```

3. Ajouter dans le switch de Execute()

4. Tester dans `builtin_test.go`

5. Documenter dans le guide utilisateur

## Standards de code

Voir `.github/prompts/common.md` pour tous les standards.

### Spécifiques au module xuples

- Thread-safety obligatoire (sync.RWMutex)
- Tests de concurrence requis
- Documentation GoDoc complète
- Pas de dépendances externes (sauf rete pour Fact)
```

**Livrables** :
- [ ] Guide de contribution créé
- [ ] Instructions claires pour extensions
- [ ] Standards documentés

### 4. Créer des exemples avancés

**Objectif** : Fournir des exemples réalistes et complets.

**Fichiers à créer** :
- `tsd/examples/xuples/workflow-orchestration.tsd`
- `tsd/examples/xuples/event-broadcasting.tsd`
- `tsd/examples/xuples/resource-allocation.tsd`

**Exemple workflow-orchestration.tsd** :

```tsd
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License

// ============================================================================
// EXEMPLE : Orchestration de Workflow
// ============================================================================
//
// Cet exemple montre comment utiliser les xuples pour orchestrer un workflow
// de traitement de commandes avec plusieurs étapes séquentielles.
//
// Workflow :
//   1. Validation de la commande
//   2. Vérification de stock
//   3. Paiement
//   4. Livraison
//
// Chaque étape crée un xuple pour la suivante.
//

// Xuple-space pour les étapes en attente
xuple-space workflow-steps {
    selection: fifo
    consumption: once
    retention: duration(1h)
}

// Types de faits
fact Order(id: string, customer: string, amount: int)
fact OrderValidated(orderId: string)
fact StockChecked(orderId: string, available: bool)
fact PaymentProcessed(orderId: string, success: bool)
fact ReadyForDelivery(orderId: string)

// Règle 1: Nouvelle commande → Validation
rule "start-workflow" {
    when {
        order: Order()
    }
    then {
        Print("Starting workflow for order: " + order.id)
        
        // Créer le fait de validation
        Insert(OrderValidated(order.id))
        
        // L'agent de validation récupérera ce xuple
        Xuple("workflow-steps", order)
    }
}

// Règle 2: Validation → Vérification de stock
rule "after-validation" {
    when {
        order: Order()
        validated: OrderValidated(orderId == order.id)
    }
    then {
        Print("Order validated: " + order.id)
        
        // Simuler la vérification de stock
        Insert(StockChecked(order.id, true))
        
        Xuple("workflow-steps", validated)
    }
}

// Règle 3: Stock OK → Paiement
rule "after-stock-check" {
    when {
        order: Order()
        stock: StockChecked(orderId == order.id, available == true)
    }
    then {
        Print("Stock available for order: " + order.id)
        
        // Simuler le paiement
        Insert(PaymentProcessed(order.id, true))
        
        Xuple("workflow-steps", stock)
    }
}

// Règle 4: Paiement OK → Livraison
rule "after-payment" {
    when {
        order: Order()
        payment: PaymentProcessed(orderId == order.id, success == true)
    }
    then {
        Print("Payment successful for order: " + order.id)
        
        // Créer le fait de livraison
        Insert(ReadyForDelivery(order.id))
        
        Xuple("workflow-steps", payment)
    }
}

// Règle 5: Livraison prête
rule "ready-for-delivery" {
    when {
        order: Order()
        ready: ReadyForDelivery(orderId == order.id)
    }
    then {
        Print("Order ready for delivery: " + order.id)
        Log("Workflow completed for order " + order.id)
    }
}
```

**Livrables** :
- [ ] Exemples avancés créés
- [ ] Commentaires explicatifs complets
- [ ] Cas d'usage réalistes
- [ ] Exemples testés et fonctionnels

### 5. Créer un INDEX de documentation

**Objectif** : Faciliter la navigation dans toute la documentation.

**Fichier à créer** : `tsd/docs/xuples/README.md`

**Contenu attendu** :

```markdown
# Documentation du Système Xuples

Bienvenue dans la documentation complète du système xuples de TSD.

## 📚 Pour les utilisateurs

- **[Guide Complet](user-guide/complete-guide.md)** - Documentation principale
- **[Commande xuple-space](user-guide/xuplespace-command.md)** - Référence syntaxe
- **[Utilisation des xuples](user-guide/using-xuples.md)** - Guide pratique

### Exemples

- **[Exemples de base](../../examples/xuples/)** - Exemples simples
- **[Workflow orchestration](../../examples/xuples/workflow-orchestration.tsd)** - Exemple avancé
- **[Event broadcasting](../../examples/xuples/event-broadcasting.tsd)** - Broadcasting
- **[Resource allocation](../../examples/xuples/resource-allocation.tsd)** - Allocation

## 🏗️ Pour les développeurs

### Architecture

- **[Vue d'ensemble](architecture/overview.md)** - Architecture générale
- **[Décisions](architecture/decisions.md)** - ADR (Architecture Decision Records)

### Implémentation

- **[Analyse existant](analysis/)** - Analyse du code existant
- **[Conception](design/)** - Spécifications de conception
- **[Implémentation](implementation/)** - Notes d'implémentation

### Tests

- **[Stratégie de tests](testing/test-strategy.md)** - Approche de tests
- **[Rapport de tests](testing/test-report.md)** - Résultats de tests

### Contribution

- **[Guide de contribution](contributing/extending-xuples.md)** - Étendre le système
- **[Standards de code](../../.github/prompts/common.md)** - Standards du projet

## 🚀 Démarrage rapide

1. **Déclarer un xuple-space** :

\```tsd
xuple-space tasks {
    selection: fifo
    consumption: once
    retention: unlimited
}
\```

2. **Créer un xuple dans une règle** :

\```tsd
rule "create-task" {
    when {
        task: Task(priority > 5)
    }
    then {
        Xuple("tasks", task)
    }
}
\```

3. **Récupérer un xuple** (depuis un agent externe - futur) :

\```go
xuple, err := xupleManager.Retrieve("tasks", "agent1")
\```

## 📖 Concepts clés

- **Xuple** : Tuple étendu (fait + faits déclencheurs)
- **Xuple-space** : Espace de stockage avec politiques
- **Politiques** : Règles de sélection, consommation, rétention
- **Agent** : Programme externe consommant des xuples

## 🔗 Liens utiles

- [Guide complet](user-guide/complete-guide.md)
- [Architecture](architecture/overview.md)
- [Exemples](../../examples/xuples/)
- [API Documentation](https://pkg.go.dev/tsd/xuples)

## 📝 Historique

- **v1.0** - Implémentation initiale du système xuples
  - Parsing `xuple-space`
  - Actions par défaut (Print, Log, Update, Insert, Retract, Xuple)
  - Exécution immédiate des actions
  - Module xuples avec politiques
  - Intégration RETE ↔ xuples

## 🆘 Support

- Consultez le [Guide de dépannage](user-guide/complete-guide.md#8-dépannage)
- Lisez les [exemples](../../examples/xuples/)
- Vérifiez les [tests](../../tests/)
```

**Livrables** :
- [ ] INDEX de documentation créé
- [ ] Navigation facilitée
- [ ] Liens vers toutes les ressources
- [ ] Démarrage rapide inclus

### 6. Créer un document de migration

**Objectif** : Aider à migrer du système tuple-space vers xuples.

**Fichier à créer** : `tsd/docs/xuples/migration/from-tuple-space.md`

**Contenu attendu** :

```markdown
# Migration de Tuple-Space vers Xuples

## Vue d'ensemble

Le système xuples remplace l'ancien système tuple-space avec :

- ✅ Exécution immédiate des actions (pas de stockage)
- ✅ Module xuples découplé de RETE
- ✅ Politiques configurables
- ✅ Actions par défaut (Print, Log, Update, Insert, Retract, Xuple)
- ✅ Parsing de la commande xuple-space

## Différences principales

| Aspect | Ancien (tuple-space) | Nouveau (xuples) |
|--------|---------------------|------------------|
| Actions | Stockées dans terminal nodes | Exécutées immédiatement |
| Récupération | Via collectActivations | Via XupleManager |
| Configuration | Hardcodée | Politiques déclaratives |
| Découplage | Intégré dans RETE | Module indépendant |

## Changements requis

### 1. Déclaration de xuple-spaces

**Avant** : Pas de déclaration explicite

**Après** : Déclaration obligatoire

\```tsd
xuple-space myspace {
    selection: fifo
    consumption: once
    retention: unlimited
}
\```

### 2. Utilisation de l'action Xuple

**Avant** : Pas d'action dédiée

**Après** : Utiliser l'action Xuple

\```tsd
rule "my-rule" {
    when {
        fact: Fact()
    }
    then {
        Xuple("myspace", fact)
    }
}
\```

### 3. Récupération des activations

**Avant** :
\```go
activations := collectActivations(network)
\```

**Après** :
\```go
xuplespace, _ := xupleManager.GetXupleSpace("myspace")
xuple, _ := xuplespace.Retrieve("agent1")
\```

### 4. Tests

**Avant** : Vérification de terminal.Memory.Tokens

**Après** : Utilisation d'observer ou GetExecutionCount

\```go
// Avant
if len(terminal.Memory.Tokens) != 1 {
    t.Error("Expected 1 activation")
}

// Après
if terminal.GetExecutionCount() != 1 {
    t.Error("Expected 1 execution")
}
\```

## Checklist de migration

- [ ] Déclarer tous les xuple-spaces nécessaires
- [ ] Remplacer les références aux activations par Xuple
- [ ] Migrer les tests (observer, GetExecutionCount)
- [ ] Supprimer les appels à collectActivations
- [ ] Vérifier que toutes les actions s'exécutent
- [ ] Tester le système complet
```

**Livrables** :
- [ ] Guide de migration créé
- [ ] Différences documentées
- [ ] Checklist de migration fournie

### 7. Générer la documentation GoDoc

**Objectif** : Générer et vérifier la documentation API.

**Tâches** :

```bash
# Générer la documentation GoDoc
go doc -all tsd/xuples > tsd/docs/xuples/api/godoc.txt

# Vérifier que toutes les fonctions exportées sont documentées
for file in tsd/xuples/*.go; do
    echo "Checking $file for GoDoc..."
    grep -E "^func [A-Z]" "$file" | while read line; do
        # Vérifier qu'il y a un commentaire au-dessus
    done
done

# Générer la documentation HTML (optionnel)
godoc -http=:6060 &
# Visiter http://localhost:6060/pkg/tsd/xuples/
```

**Livrables** :
- [ ] Toutes les fonctions exportées documentées
- [ ] Exemples d'utilisation en commentaires
- [ ] Documentation GoDoc générée

### 8. Créer un CHANGELOG

**Objectif** : Documenter l'historique des changements.

**Fichier à créer/modifier** : `tsd/CHANGELOG.md`

**Ajouts attendus** :

```markdown
## [Unreleased]

### Added - Système Xuples

#### Nouvelles fonctionnalités

- **Xuple-spaces** : Espaces de coordination configurables
  - Politiques de sélection (random, fifo, lifo)
  - Politiques de consommation (once, per-agent, limited)
  - Politiques de rétention (unlimited, duration)
  
- **Commande xuple-space** : Déclaration de xuple-spaces dans le langage TSD
  
- **Actions par défaut** : Système d'actions prédéfinies
  - Print(string) : Affichage sur stdout
  - Log(string) : Génération de traces
  - Update(fact) : Modification de fait
  - Insert(fact) : Création de fait
  - Retract(id) : Suppression de fait
  - Xuple(xuplespace, fact) : Création de xuple

- **Module xuples** : Nouveau package pour la gestion des xuples
  - Thread-safe (sync.RWMutex)
  - Politiques extensibles (interfaces)
  - Cycle de vie complet des xuples

#### Changements

- **Terminal nodes** : Exécution immédiate au lieu de stockage
  - Les actions sont exécutées dès l'activation
  - Suppression du stockage des tokens
  - Observer pattern pour l'observabilité

- **Compilateur** : Support du nouveau système
  - Chargement automatique des actions par défaut
  - Validation des xuple-spaces
  - Instanciation des politiques

#### Deprecated

- `collectActivations()` : Remplacé par le système d'observer

#### Documentation

- Guide utilisateur complet
- Documentation d'architecture
- Exemples avancés (workflow, broadcasting, allocation)
- Guide de migration tuple-space → xuples
- API documentation (GoDoc)

#### Tests

- Tests unitaires (>80% couverture)
- Tests d'intégration
- Tests E2E
- Tests de performance (benchmarks)
- Tests de concurrence (race detector)
```

**Livrables** :
- [ ] CHANGELOG mis à jour
- [ ] Toutes les nouveautés documentées
- [ ] Format respect de Keep a Changelog

## 📁 Structure finale de la documentation

```
tsd/docs/xuples/
├── README.md                           # INDEX principal
├── user-guide/
│   ├── complete-guide.md               # Guide complet utilisateur
│   ├── xuplespace-command.md           # Référence commande
│   └── using-xuples.md                 # Guide pratique
├── architecture/
│   ├── overview.md                     # Vue d'ensemble
│   └── decisions.md                    # ADR
├── analysis/
│   ├── 00-INDEX.md
│   ├── 01-current-action-parsing.md
│   ├── 02-terminal-nodes.md
│   ├── 03-token-fact-structures.md
│   ├── 04-action-executor.md
│   └── 05-existing-tests.md
├── design/
│   ├── 00-INDEX.md
│   ├── 01-data-structures.md
│   ├── 02-interfaces.md
│   ├── 03-policies.md
│   ├── 04-rete-integration.md
│   ├── 05-lifecycle.md
│   ├── 06-agent-interface.md
│   └── 07-package-structure.md
├── implementation/
│   ├── 01-parser-analysis.md
│   ├── 02-xuplespace-syntax.md
│   ├── 03-current-action-system.md
│   ├── 04-default-actions-design.md
│   ├── 05-terminal-node-current-behavior.md
│   ├── 06-immediate-execution-design.md
│   └── 07-rete-xuples-integration.md
├── testing/
│   ├── test-strategy.md
│   └── test-report.md
├── contributing/
│   └── extending-xuples.md
├── migration/
│   └── from-tuple-space.md
└── api/
    └── godoc.txt
```

## ✅ Critères de succès

- [ ] Guide utilisateur complet et clair
- [ ] Documentation d'architecture détaillée
- [ ] Guide de contribution pour extensions
- [ ] Exemples avancés réalistes
- [ ] INDEX de documentation créé
- [ ] Guide de migration fourni
- [ ] Documentation GoDoc complète
- [ ] CHANGELOG mis à jour
- [ ] Toutes les références croisées correctes
- [ ] Documentation relue et validée
- [ ] Aucune référence cassée
- [ ] Format Markdown valide

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- Keep a Changelog - https://keepachangelog.com/
- GoDoc Best Practices - https://go.dev/blog/godoc
- Arc42 - Template de documentation d'architecture

## 🎯 Prochaine étape

Une fois la documentation finalisée, passer au prompt **10-final-validation.md** pour la validation finale complète du système xuples avant intégration.