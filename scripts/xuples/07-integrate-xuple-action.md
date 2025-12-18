# Prompt 07 - Intégration de l'action Xuple avec le module xuples

## 🎯 Objectif

Intégrer l'action `Xuple` avec le module xuples pour permettre la création de xuples depuis les règles RETE.

Cette intégration doit :
- Connecter le BuiltinActionExecutor au XupleManager
- Permettre à l'action Xuple de créer des xuples dans les xuple-spaces
- Extraire correctement les faits déclencheurs des tokens
- Valider que le xuple-space existe
- Gérer les erreurs de manière robuste
- Maintenir le découplage entre RETE et xuples

## 📋 Tâches

### 1. Analyser l'interface entre RETE et xuples

**Objectif** : Identifier précisément comment les deux modules doivent interagir.

- [ ] Examiner l'interface XupleManager définie dans le module xuples
- [ ] Identifier les données nécessaires de RETE (Token, Fact)
- [ ] Définir le contrat d'intégration
- [ ] Identifier les points de configuration

**Livrables** :
- Créer `tsd/docs/xuples/implementation/07-rete-xuples-integration.md` documentant :
  - Diagramme d'architecture montrant RETE et xuples
  - Points d'interaction entre les modules
  - Flux de données complet
  - Contrat d'interface
  - Stratégie d'injection de dépendances
  - Gestion d'erreurs entre modules

### 2. Adapter BuiltinActionExecutor pour accepter XupleManager

**Objectif** : Modifier le BuiltinActionExecutor pour utiliser le XupleManager.

**Fichier à modifier** : `tsd/rete/actions/builtin.go`

**Modifications attendues** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
    "fmt"
    "log"
    
    "tsd/rete"
    "tsd/xuples"
)

// BuiltinActionExecutor exécute les actions par défaut du système
type BuiltinActionExecutor struct {
    network      *rete.Network
    xupleManager xuples.XupleManager
}

// NewBuiltinActionExecutor crée un nouvel exécuteur d'actions natives
func NewBuiltinActionExecutor(network *rete.Network, xupleManager xuples.XupleManager) *BuiltinActionExecutor {
    return &BuiltinActionExecutor{
        network:      network,
        xupleManager: xupleManager,
    }
}

// Execute exécute une action par défaut
func (e *BuiltinActionExecutor) Execute(actionName string, args []interface{}, token *rete.Token) error {
    switch actionName {
    case "Print":
        return e.executePrint(args)
    case "Log":
        return e.executeLog(args)
    case "Update":
        return e.executeUpdate(args, token)
    case "Insert":
        return e.executeInsert(args)
    case "Retract":
        return e.executeRetract(args)
    case "Xuple":
        return e.executeXuple(args, token)
    default:
        return fmt.Errorf("unknown builtin action: %s", actionName)
    }
}

// executeXuple implémente l'action Xuple
func (e *BuiltinActionExecutor) executeXuple(args []interface{}, token *rete.Token) error {
    if len(args) != 2 {
        return fmt.Errorf("Xuple expects 2 arguments (xuplespace, fact), got %d", len(args))
    }
    
    // Valider les arguments
    xuplespace, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Xuple expects string as first argument (xuplespace name), got %T", args[0])
    }
    
    fact, ok := args[1].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Xuple expects fact as second argument, got %T", args[1])
    }
    
    // Vérifier que le xupleManager est configuré
    if e.xupleManager == nil {
        return fmt.Errorf("XupleManager not configured, cannot execute Xuple action")
    }
    
    // Extraire les faits déclencheurs du token
    triggeringFacts := extractTriggeringFacts(token)
    
    log.Printf("🔧 Creating xuple in '%s' with fact %s and %d triggering facts",
        xuplespace, fact.ID, len(triggeringFacts))
    
    // Déléguer au XupleManager
    err := e.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
    if err != nil {
        return fmt.Errorf("failed to create xuple in '%s': %w", xuplespace, err)
    }
    
    log.Printf("✅ Xuple created successfully in '%s'", xuplespace)
    return nil
}

// extractTriggeringFacts extrait tous les faits d'un token combiné
func extractTriggeringFacts(token *rete.Token) []*rete.Fact {
    if token == nil {
        return []*rete.Fact{}
    }
    
    var facts []*rete.Fact
    
    // Parcourir la chaîne de tokens (du plus récent au plus ancien)
    for t := token; t != nil; t = t.Parent {
        if t.Fact != nil {
            facts = append(facts, t.Fact)
        }
    }
    
    // Inverser pour avoir l'ordre chronologique (plus ancien d'abord)
    for i := 0; i < len(facts)/2; i++ {
        facts[i], facts[len(facts)-1-i] = facts[len(facts)-1-i], facts[i]
    }
    
    return facts
}
```

**Livrables** :
- [ ] BuiltinActionExecutor modifié pour accepter XupleManager
- [ ] Implémentation complète de executeXuple
- [ ] Extraction des faits déclencheurs
- [ ] Validation des arguments robuste
- [ ] Gestion d'erreurs avec messages clairs
- [ ] Logging des opérations
- [ ] Documentation GoDoc mise à jour

### 3. Intégrer la création de XupleManager dans le compilateur

**Objectif** : Créer et configurer le XupleManager lors de la compilation.

**Fichier à modifier** : `tsd/compiler/compiler.go`

**Modifications attendues** :

```go
// Dans la structure Compiler
type Compiler struct {
    // ... champs existants ...
    
    xupleManager xuples.XupleManager
}

// Dans NewCompiler
func NewCompiler() (*Compiler, error) {
    c := &Compiler{
        actions:      make(map[string]*ast.ActionDeclaration),
        xupleSpaces:  make(map[string]*ast.XupleSpaceDeclaration),
        xupleManager: xuples.NewXupleManager(), // Créer le manager
        // ... autres champs ...
    }
    
    // ... reste de l'initialisation ...
    
    return c, nil
}

// Méthode pour accéder au XupleManager
func (c *Compiler) GetXupleManager() xuples.XupleManager {
    return c.xupleManager
}
```

**Livrables** :
- [ ] XupleManager créé à l'initialisation du compilateur
- [ ] Méthode d'accès au manager
- [ ] Documentation mise à jour

### 4. Créer les xuple-spaces depuis les déclarations parsées

**Objectif** : Instancier les xuple-spaces déclarés dans le programme TSD.

**Fichier à modifier** : `tsd/compiler/compiler.go`

**Code attendu** :

```go
// InstantiateXupleSpaces crée les xuple-spaces déclarés
func (c *Compiler) InstantiateXupleSpaces() error {
    for name, decl := range c.xupleSpaces {
        // Construire la configuration depuis la déclaration AST
        config, err := c.buildXupleSpaceConfig(decl)
        if err != nil {
            return fmt.Errorf("failed to build config for xuple-space '%s': %w", name, err)
        }
        
        // Créer le xuple-space
        err = c.xupleManager.CreateXupleSpace(name, config)
        if err != nil {
            return fmt.Errorf("failed to create xuple-space '%s': %w", name, err)
        }
        
        log.Printf("✅ Created xuple-space '%s' with policies: selection=%s, consumption=%s, retention=%s",
            name,
            config.SelectionPolicy.Name(),
            config.ConsumptionPolicy.Name(),
            config.RetentionPolicy.Name())
    }
    
    return nil
}

// buildXupleSpaceConfig construit une configuration depuis une déclaration AST
func (c *Compiler) buildXupleSpaceConfig(decl *ast.XupleSpaceDeclaration) (xuples.XupleSpaceConfig, error) {
    // Construire la politique de sélection
    var selectionPolicy xuples.SelectionPolicy
    switch decl.SelectionPolicy {
    case ast.SelectionRandom:
        selectionPolicy = xuples.NewRandomSelectionPolicy()
    case ast.SelectionFIFO:
        selectionPolicy = xuples.NewFIFOSelectionPolicy()
    case ast.SelectionLIFO:
        selectionPolicy = xuples.NewLIFOSelectionPolicy()
    default:
        return xuples.XupleSpaceConfig{}, fmt.Errorf("unknown selection policy: %v", decl.SelectionPolicy)
    }
    
    // Construire la politique de consommation
    var consumptionPolicy xuples.ConsumptionPolicy
    switch decl.ConsumptionPolicy.Type {
    case ast.ConsumptionOnce:
        consumptionPolicy = xuples.NewOnceConsumptionPolicy()
    case ast.ConsumptionPerAgent:
        consumptionPolicy = xuples.NewPerAgentConsumptionPolicy()
    case ast.ConsumptionLimited:
        consumptionPolicy = xuples.NewLimitedConsumptionPolicy(decl.ConsumptionPolicy.Limit)
    default:
        return xuples.XupleSpaceConfig{}, fmt.Errorf("unknown consumption policy: %v", decl.ConsumptionPolicy.Type)
    }
    
    // Construire la politique de rétention
    var retentionPolicy xuples.RetentionPolicy
    switch decl.RetentionPolicy.Type {
    case ast.RetentionUnlimited:
        retentionPolicy = xuples.NewUnlimitedRetentionPolicy()
    case ast.RetentionDuration:
        retentionPolicy = xuples.NewDurationRetentionPolicy(decl.RetentionPolicy.Duration)
    default:
        return xuples.XupleSpaceConfig{}, fmt.Errorf("unknown retention policy: %v", decl.RetentionPolicy.Type)
    }
    
    return xuples.XupleSpaceConfig{
        Name:              decl.Name,
        SelectionPolicy:   selectionPolicy,
        ConsumptionPolicy: consumptionPolicy,
        RetentionPolicy:   retentionPolicy,
    }, nil
}
```

**Livrables** :
- [ ] Méthode InstantiateXupleSpaces implémentée
- [ ] Conversion AST → xuples.XupleSpaceConfig
- [ ] Création des politiques appropriées
- [ ] Validation et gestion d'erreurs
- [ ] Logging de la création
- [ ] Appel de InstantiateXupleSpaces dans le flux de compilation

### 5. Configurer le réseau RETE avec le XupleManager

**Objectif** : Passer le XupleManager au BuiltinActionExecutor lors de la création du réseau.

**Fichier à modifier** : `tsd/internal/servercmd/servercmd.go` (ou équivalent)

**Modifications attendues** :

```go
func executeTSDProgram(source string) (*tsdio.ExecuteResponse, error) {
    // ... parsing ...
    
    // Créer le compilateur (qui crée le XupleManager)
    compiler, err := compiler.NewCompiler()
    if err != nil {
        return nil, fmt.Errorf("failed to create compiler: %w", err)
    }
    
    // ... compilation ...
    
    // Instancier les xuple-spaces déclarés
    err = compiler.InstantiateXupleSpaces()
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate xuple-spaces: %w", err)
    }
    
    // Créer le réseau RETE
    network := compiler.BuildReteNetwork()
    
    // Créer l'executor avec le XupleManager
    executor := actions.NewBuiltinActionExecutor(network, compiler.GetXupleManager())
    
    // Configurer le réseau avec l'executor
    network.SetActionExecutor(executor)
    
    // ... reste de l'exécution ...
}
```

**Livrables** :
- [ ] XupleManager passé au BuiltinActionExecutor
- [ ] Xuple-spaces instanciés avant l'exécution
- [ ] Configuration correcte du réseau
- [ ] Ordre d'initialisation respecté

### 6. Valider que le xuple-space existe lors de l'exécution

**Objectif** : S'assurer qu'une action Xuple ne peut référencer qu'un xuple-space déclaré.

**Validation au compile-time (optionnel mais recommandé)** :

Dans le compilateur, lors de la compilation des règles :

```go
// Dans la validation des règles
func (c *Compiler) validateActionCall(actionCall *ast.ActionCall) error {
    // Si c'est l'action Xuple
    if actionCall.ActionName == "Xuple" {
        if len(actionCall.Arguments) < 1 {
            return fmt.Errorf("Xuple action requires xuplespace name")
        }
        
        // Vérifier que le premier argument (xuplespace) est une constante string
        if xuplespaceArg, ok := actionCall.Arguments[0].(*ast.StringLiteral); ok {
            xuplespace := xuplespaceArg.Value
            
            // Vérifier que le xuple-space a été déclaré
            if _, exists := c.xupleSpaces[xuplespace]; !exists {
                return fmt.Errorf("xuple-space '%s' not declared (line %d)",
                    xuplespace, actionCall.Location.Line)
            }
        }
        // Si ce n'est pas une constante, on ne peut pas vérifier au compile-time
    }
    
    return nil
}
```

**Validation au runtime** (déjà faite dans executeXuple via CreateXuple qui retourne une erreur si le xuple-space n'existe pas).

**Livrables** :
- [ ] Validation compile-time si possible
- [ ] Validation runtime robuste
- [ ] Messages d'erreur clairs

### 7. Créer des tests d'intégration complets

**Objectif** : Tester l'intégration complète RETE → action Xuple → module xuples.

**Fichier à créer** : `tsd/tests/integration/xuples_integration_test.go`

**Tests attendus** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package integration

import (
    "testing"
    
    "tsd/compiler"
    "tsd/parser"
    "tsd/rete/actions"
)

func TestXupleAction_Integration(t *testing.T) {
    t.Log("🧪 TEST INTÉGRATION ACTION XUPLE")
    
    program := `
        // Déclaration du xuple-space
        xuple-space notifications {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        // Déclaration du fait
        fact Alert(level: string, message: string)
        
        // Règle qui crée un xuple
        rule "high-priority-alert" {
            when {
                alert: Alert(level == "high")
            }
            then {
                Log("High priority alert detected")
                Xuple("notifications", alert)
            }
        }
    `
    
    // Parser
    result, err := parser.ParseTSD(program)
    if err != nil {
        t.Fatalf("❌ Erreur parsing: %v", err)
    }
    
    // Compiler
    comp, err := compiler.NewCompiler()
    if err != nil {
        t.Fatalf("❌ Erreur création compilateur: %v", err)
    }
    
    err = comp.Compile(result)
    if err != nil {
        t.Fatalf("❌ Erreur compilation: %v", err)
    }
    
    // Instancier les xuple-spaces
    err = comp.InstantiateXupleSpaces()
    if err != nil {
        t.Fatalf("❌ Erreur instanciation xuple-spaces: %v", err)
    }
    
    // Construire le réseau RETE
    network := comp.BuildReteNetwork()
    
    // Configurer l'executor
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    // Insérer un fait qui déclenche la règle
    alert := &rete.Fact{
        ID:   "alert1",
        Type: "Alert",
        Attributes: map[string]interface{}{
            "level":   "high",
            "message": "System overload",
        },
    }
    
    err = network.InsertFact(alert)
    if err != nil {
        t.Fatalf("❌ Erreur insertion fait: %v", err)
    }
    
    // Vérifier que le xuple a été créé
    xupleManager := comp.GetXupleManager()
    xuplespace, err := xupleManager.GetXupleSpace("notifications")
    if err != nil {
        t.Fatalf("❌ Erreur récupération xuple-space: %v", err)
    }
    
    count := xuplespace.Count()
    if count != 1 {
        t.Errorf("❌ Attendu 1 xuple, reçu %d", count)
    }
    
    // Récupérer le xuple
    xuple, err := xuplespace.Retrieve("agent1")
    if err != nil {
        t.Fatalf("❌ Erreur récupération xuple: %v", err)
    }
    
    if xuple == nil {
        t.Fatal("❌ Xuple récupéré est nil")
    }
    
    // Vérifier le contenu
    if xuple.Fact.ID != "alert1" {
        t.Errorf("❌ Mauvais fait dans le xuple: %s", xuple.Fact.ID)
    }
    
    if len(xuple.TriggeringFacts) != 1 {
        t.Errorf("❌ Attendu 1 fait déclencheur, reçu %d", len(xuple.TriggeringFacts))
    }
    
    t.Log("✅ Intégration complète fonctionne")
}

func TestXupleAction_UndeclaredXupleSpace(t *testing.T) {
    t.Log("🧪 TEST XUPLE-SPACE NON DÉCLARÉ")
    
    program := `
        fact Alert(level: string)
        
        rule "bad-rule" {
            when {
                alert: Alert(level == "high")
            }
            then {
                Xuple("nonexistent", alert)
            }
        }
    `
    
    // ... parsing et compilation ...
    
    // Insérer un fait
    alert := &rete.Fact{ID: "a1", Type: "Alert"}
    
    // L'exécution devrait échouer (xuple-space non déclaré)
    // Vérifier via l'observer que l'action a échoué
    
    t.Log("✅ Erreur correctement détectée pour xuple-space non déclaré")
}

func TestXupleAction_MultipleRules(t *testing.T) {
    t.Log("🧪 TEST PLUSIEURS RÈGLES CRÉANT DES XUPLES")
    
    program := `
        xuple-space events {
            selection: lifo
            consumption: per-agent
            retention: unlimited
        }
        
        fact Event(type: string)
        fact Condition(status: string)
        
        rule "rule1" {
            when {
                e: Event(type == "start")
            }
            then {
                Xuple("events", e)
            }
        }
        
        rule "rule2" {
            when {
                c: Condition(status == "active")
            }
            then {
                Xuple("events", c)
            }
        }
    `
    
    // ... setup ...
    
    // Insérer plusieurs faits
    // Vérifier que plusieurs xuples sont créés
    
    t.Log("✅ Plusieurs règles peuvent créer des xuples")
}

func TestXupleAction_WithMultipleTriggeringFacts(t *testing.T) {
    t.Log("🧪 TEST XUPLES AVEC PLUSIEURS FAITS DÉCLENCHEURS")
    
    program := `
        xuple-space complex {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        fact Person(name: string, age: int)
        fact Department(name: string)
        fact Assignment(person: string, dept: string)
        
        rule "complex-rule" {
            when {
                p: Person(age >= 18)
                d: Department(name == "Engineering")
                a: Assignment(person == p.name, dept == d.name)
            }
            then {
                Xuple("complex", a)
            }
        }
    `
    
    // ... setup ...
    
    // Insérer les trois faits
    // Vérifier que le xuple contient les 3 faits déclencheurs
    
    if len(xuple.TriggeringFacts) != 3 {
        t.Errorf("❌ Attendu 3 faits déclencheurs, reçu %d", len(xuple.TriggeringFacts))
    }
    
    t.Log("✅ Faits déclencheurs multiples correctement extraits")
}
```

**Livrables** :
- [ ] Tests d'intégration complets
- [ ] Test du cas nominal
- [ ] Test d'erreur (xuple-space non déclaré)
- [ ] Test avec plusieurs règles
- [ ] Test avec plusieurs faits déclencheurs
- [ ] Couverture > 80%
- [ ] Tous les tests passent

### 8. Créer des exemples utilisateur

**Objectif** : Fournir des exemples complets d'utilisation.

**Fichiers à créer** :
- `tsd/examples/xuples/basic-xuple-creation.tsd`
- `tsd/examples/xuples/multi-policy-xuples.tsd`
- `tsd/examples/xuples/agent-workflow.tsd`

**Exemple 1 - Basique** :

```tsd
// Exemple basique de création de xuples

xuple-space tasks {
    selection: fifo
    consumption: once
    retention: unlimited
}

fact Task(id: string, priority: int)

rule "high-priority-task" {
    when {
        task: Task(priority >= 8)
    }
    then {
        Print("High priority task detected: " + task.id)
        Xuple("tasks", task)
        Log("Task added to xuple-space")
    }
}
```

**Exemple 2 - Politiques multiples** :

```tsd
// Exemple avec différentes politiques

xuple-space urgent {
    selection: lifo
    consumption: once
    retention: duration(5m)
}

xuple-space shared {
    selection: random
    consumption: per-agent
    retention: unlimited
}

xuple-space limited {
    selection: fifo
    consumption: limited(3)
    retention: duration(1h)
}

fact Message(content: string, urgent: bool)

rule "urgent-message" {
    when {
        msg: Message(urgent == true)
    }
    then {
        Xuple("urgent", msg)
    }
}

rule "shared-message" {
    when {
        msg: Message(urgent == false)
    }
    then {
        Xuple("shared", msg)
    }
}
```

**Livrables** :
- [ ] Exemples TSD créés et testés
- [ ] Commentaires explicatifs
- [ ] Cas d'usage variés
- [ ] Documentation associée

### 9. Documenter l'intégration pour les utilisateurs

**Objectif** : Créer la documentation utilisateur complète.

**Fichier à créer** : `tsd/docs/xuples/user-guide/using-xuples.md`

**Contenu attendu** :

```markdown
# Guide d'utilisation des Xuples

## Introduction

Les xuples permettent de créer des espaces de coordination entre le moteur de règles RETE et des agents externes.

## Déclaration d'un xuple-space

\```tsd
xuple-space <nom> {
    selection: <random|fifo|lifo>
    consumption: <once|per-agent|limited(n)>
    retention: <unlimited|duration(temps)>
}
\```

## Utilisation de l'action Xuple

L'action `Xuple` crée un xuple dans un xuple-space :

\```tsd
Xuple("<nom-xuple-space>", <fait>)
\```

Le xuple créé contient :
- Le fait passé en argument
- Tous les faits qui ont déclenché la règle

## Exemple complet

\```tsd
xuple-space notifications {
    selection: fifo
    consumption: once
    retention: unlimited
}

fact Alert(level: string, message: string)

rule "critical-alert" {
    when {
        alert: Alert(level == "critical")
    }
    then {
        Xuple("notifications", alert)
    }
}
\```

## Politiques

### Sélection
- `random` : Sélection aléatoire
- `fifo` : Premier entré, premier sorti
- `lifo` : Dernier entré, premier sorti

### Consommation
- `once` : Un seul consommateur total
- `per-agent` : Une fois par agent
- `limited(n)` : Maximum n consommations

### Rétention
- `unlimited` : Pas d'expiration
- `duration(temps)` : Expire après la durée (ex: 5m, 1h, 2d)
```

**Livrables** :
- [ ] Documentation utilisateur complète
- [ ] Exemples clairs
- [ ] Cas d'usage documentés

## 📁 Structure finale

```
tsd/
├── docs/xuples/
│   ├── implementation/
│   │   └── 07-rete-xuples-integration.md
│   └── user-guide/
│       └── using-xuples.md
├── examples/xuples/
│   ├── basic-xuple-creation.tsd
│   ├── multi-policy-xuples.tsd
│   └── agent-workflow.tsd
├── rete/actions/
│   └── builtin.go                    # Modifié
├── compiler/
│   └── compiler.go                   # Modifié
├── internal/servercmd/
│   └── servercmd.go                  # Modifié
└── tests/integration/
    └── xuples_integration_test.go    # Nouveau
```

## ✅ Critères de succès

- [ ] BuiltinActionExecutor intégré avec XupleManager
- [ ] Xuple-spaces instanciés depuis les déclarations
- [ ] Action Xuple fonctionnelle
- [ ] Extraction des faits déclencheurs correcte
- [ ] Validation compile-time et runtime
- [ ] Tests d'intégration complets avec couverture > 80%
- [ ] Tous les tests passent
- [ ] Exemples fonctionnels fournis
- [ ] Documentation utilisateur complète
- [ ] `make test-integration` passe
- [ ] Découplage maintenu entre RETE et xuples

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/design/` - Conception
- `tsd/docs/xuples/implementation/` - Documentation d'implémentation
- Effective Go - https://go.dev/doc/effective_go

## 🎯 Prochaine étape

Une fois l'action Xuple intégrée, passer au prompt **08-test-complete-system.md** pour tester l'ensemble du système de manière exhaustive.