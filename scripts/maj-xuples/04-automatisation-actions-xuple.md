# 🔧 Prompt 04 - Automatisation des Actions Xuple

---

## 🎯 Objectif

**Automatiser complètement l'exécution des actions `Xuple()` dans les règles TSD**, en s'appuyant sur le parser amélioré (Prompt 01) et les xuple-spaces créés automatiquement (Prompt 03), pour éliminer toute configuration manuelle des handlers et permettre l'utilisation directe de `Xuple()` dans les règles.

### Contexte

Actuellement, même avec les améliorations précédentes :
- Le parser supporte les faits inline dans les actions (Prompt 01)
- Les xuple-spaces sont créés automatiquement (Prompt 03)
- **MAIS** : l'action `Xuple()` nécessite encore un handler configuré manuellement

L'objectif est que **les actions `Xuple()` fonctionnent automatiquement** dès qu'elles sont déclarées dans une règle TSD.

### Prérequis

- ✅ Prompt 01 : Parser supporte `Xuple("space", Fact(...))` avec références aux champs
- ✅ Prompt 02 : Package `api` avec `Pipeline.IngestFile()`
- ✅ Prompt 03 : Xuple-spaces créés automatiquement

### Résultat Attendu Final

Un fichier TSD comme celui-ci :

```tsd
xuple-space alerts {
    selection: fifo,
    consumption: once
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    message: string,
    temp: float
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            message: "High temperature detected",
            temp: t.value
        ))
    }
}
```

Fonctionne **immédiatement** après `pipeline.IngestFile("rules.tsd")`, sans aucune configuration manuelle du handler Xuple.

---

## 📋 Analyse Préliminaire

### 1. Comprendre le Mécanisme Actuel des Actions

**Fichiers clés à examiner :**

```
tsd/internal/rete/
├── network.go                # Réseau RETE
├── action.go                 # Interface Action + ActionExecutor
├── constraint_pipeline.go    # Point d'entrée
└── xuple_action.go          # Implémentation XupleAction (si existe)

tsd/api/
├── pipeline.go               # Pipeline API
└── config.go                 # Configuration

tsd/xuples/
├── manager.go                # XupleManager
└── xuplespace.go             # XupleSpace
```

**Questions à résoudre :**

1. **Comment les actions sont-elles enregistrées dans le réseau RETE ?**
   - Via `Network.RegisterAction(name string, factory ActionFactory)`
   - Chaque règle référence des actions par leur nom

2. **Qu'est-ce qu'un ActionFactory ?**
   ```go
   type ActionFactory func(args []interface{}) (Action, error)
   ```
   - Crée une instance d'action à partir des arguments

3. **Comment l'action Xuple accède-t-elle au XupleManager ?**
   - Actuellement via injection/contexte
   - Besoin d'un mécanisme pour passer le `XupleManager` aux actions

4. **À quel moment enregistrer l'action Xuple ?**
   - Lors de l'initialisation du pipeline API
   - Avant la conversion AST → RETE

### 2. Identifier l'Implémentation Actuelle de XupleAction

**Vérifier si `xuple_action.go` existe :**

```go
// tsd/internal/rete/xuple_action.go (si existe)
type XupleAction struct {
    spaceName string
    fact      *Fact
    manager   *xuples.XupleManager  // ⚠️ Problème: cycle d'import!
}
```

**Problème identifié :** Import cycle entre `rete` et `xuples`.

**Solution retenue :** Définir l'action Xuple dans le package `api`, pas dans `rete`.

### 3. Comprendre le Flux d'Exécution des Actions

```
1. Parsing TSD
   ↓
2. Conversion AST → RETE
   - Règles contiennent des ActionNodes
   ↓
3. Build du réseau RETE
   - Actions créées via ActionFactory
   ↓
4. Propagation de faits
   ↓
5. Activation de règles
   ↓
6. Exécution des actions
   - ActionExecutor.Execute(action, bindings)
```

**Point clé :** Les actions doivent être enregistrées **avant l'étape 2** (conversion AST).

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Définir l'Interface Action dans RETE

**Fichier :** `tsd/internal/rete/action.go`

**Objectif :** S'assurer que l'interface Action est suffisamment flexible pour supporter les actions Xuple.

#### 1.1 Interface Action

```go
// Action représente une action exécutable par une règle.
// Les actions sont exécutées lorsqu'une règle est activée.
type Action interface {
    // Execute exécute l'action avec le contexte donné.
    // Le contexte contient les faits déclencheurs et leurs bindings.
    Execute(ctx *ActionContext) error

    // Name retourne le nom de l'action (pour debug/logging).
    Name() string
}

// ActionContext fournit le contexte d'exécution d'une action.
type ActionContext struct {
    // Network est le réseau RETE dans lequel l'action s'exécute.
    Network *Network

    // TriggeringFacts contient les faits qui ont déclenché la règle.
    // Map: variable name → Fact
    TriggeringFacts map[string]*Fact

    // Bindings contient les valeurs des variables capturées.
    // Map: variable name → value
    Bindings map[string]interface{}

    // RuleName est le nom de la règle qui a déclenché l'action.
    RuleName string

    // UserContext permet de passer des données arbitraires.
    // Utilisé notamment pour passer le XupleManager aux actions Xuple.
    UserContext map[string]interface{}
}

// ActionFactory est une fonction qui crée une instance d'action.
type ActionFactory func(args []ActionArgument) (Action, error)

// ActionArgument représente un argument passé à une action.
// Peut être une valeur littérale, une référence à un champ, ou un fait inline.
type ActionArgument interface {
    isActionArgument()
}

// LiteralArgument représente une valeur littérale (string, int, float, bool).
type LiteralArgument struct {
    Value interface{}
}

func (LiteralArgument) isActionArgument() {}

// FieldReferenceArgument représente une référence à un champ d'un fait.
// Exemple: t.sensorId
type FieldReferenceArgument struct {
    Variable  string   // "t"
    FieldPath []string // ["sensorId"]
}

func (FieldReferenceArgument) isActionArgument() {}

// InlineFactArgument représente un fait créé inline.
// Exemple: Alert(sensorId: t.sensorId, message: "...")
type InlineFactArgument struct {
    TypeName string
    Fields   map[string]ActionArgument
}

func (InlineFactArgument) isActionArgument() {}
```

#### 1.2 ActionExecutor

```go
// ActionExecutor gère l'exécution des actions avec résolution des arguments.
type ActionExecutor struct {
    network     *Network
    userContext map[string]interface{}
}

// NewActionExecutor crée un nouvel exécuteur d'actions.
func NewActionExecutor(network *Network) *ActionExecutor {
    return &ActionExecutor{
        network:     network,
        userContext: make(map[string]interface{}),
    }
}

// SetUserContext définit une valeur dans le contexte utilisateur.
// Utilisé pour passer le XupleManager aux actions Xuple.
func (e *ActionExecutor) SetUserContext(key string, value interface{}) {
    e.userContext[key] = value
}

// Execute exécute une action avec les faits déclencheurs.
func (e *ActionExecutor) Execute(action Action, triggeringFacts map[string]*Fact, bindings map[string]interface{}, ruleName string) error {
    ctx := &ActionContext{
        Network:         e.network,
        TriggeringFacts: triggeringFacts,
        Bindings:        bindings,
        RuleName:        ruleName,
        UserContext:     e.userContext,
    }

    return action.Execute(ctx)
}

// ResolveArgument résout un argument d'action en utilisant le contexte.
func (e *ActionExecutor) ResolveArgument(arg ActionArgument, ctx *ActionContext) (interface{}, error) {
    switch a := arg.(type) {
    case LiteralArgument:
        return a.Value, nil

    case FieldReferenceArgument:
        return e.resolveFieldReference(a, ctx)

    case InlineFactArgument:
        return e.resolveInlineFact(a, ctx)

    default:
        return nil, fmt.Errorf("unknown argument type: %T", arg)
    }
}

// resolveFieldReference résout une référence à un champ.
func (e *ActionExecutor) resolveFieldReference(ref FieldReferenceArgument, ctx *ActionContext) (interface{}, error) {
    // Récupérer le fait lié à la variable
    fact, ok := ctx.TriggeringFacts[ref.Variable]
    if !ok {
        return nil, fmt.Errorf("undefined variable: %s", ref.Variable)
    }

    // Naviguer dans les champs
    value := fact.Data
    for _, field := range ref.FieldPath {
        switch v := value.(type) {
        case map[string]interface{}:
            var ok bool
            value, ok = v[field]
            if !ok {
                return nil, fmt.Errorf("field not found: %s", field)
            }
        default:
            return nil, fmt.Errorf("cannot access field %s on non-struct type", field)
        }
    }

    return value, nil
}

// resolveInlineFact résout un fait créé inline.
func (e *ActionExecutor) resolveInlineFact(inlineFact InlineFactArgument, ctx *ActionContext) (*Fact, error) {
    // Résoudre tous les champs
    resolvedFields := make(map[string]interface{})
    for fieldName, fieldArg := range inlineFact.Fields {
        value, err := e.ResolveArgument(fieldArg, ctx)
        if err != nil {
            return nil, fmt.Errorf("resolving field %s: %w", fieldName, err)
        }
        resolvedFields[fieldName] = value
    }

    // Créer le fait
    fact := e.network.CreateFact(inlineFact.TypeName, resolvedFields)
    return fact, nil
}
```

---

### Tâche 2: Implémenter XupleAction dans le Package API

**Fichier :** `tsd/api/xuple_action.go`

**Objectif :** Implémenter l'action Xuple qui utilise le XupleManager pour créer des xuples.

#### 2.1 Structure XupleAction

```go
package api

import (
    "fmt"

    "github.com/resinsec/tsd/internal/rete"
    "github.com/resinsec/tsd/xuples"
)

// XupleAction représente une action qui crée un xuple dans un xuple-space.
//
// Syntaxe TSD:
//   Xuple("spaceName", FactType(...))
type XupleAction struct {
    // SpaceNameArg est l'argument représentant le nom du xuple-space.
    // Peut être une chaîne littérale ou une référence à un champ.
    SpaceNameArg rete.ActionArgument

    // FactArg est l'argument représentant le fait à stocker comme xuple.
    // Typiquement un InlineFactArgument.
    FactArg rete.ActionArgument

    // executor est utilisé pour résoudre les arguments.
    executor *rete.ActionExecutor
}

// NewXupleAction crée une nouvelle action Xuple.
func NewXupleAction(spaceNameArg, factArg rete.ActionArgument, executor *rete.ActionExecutor) *XupleAction {
    return &XupleAction{
        SpaceNameArg: spaceNameArg,
        FactArg:      factArg,
        executor:     executor,
    }
}

// Name retourne le nom de l'action.
func (a *XupleAction) Name() string {
    return "Xuple"
}

// Execute exécute l'action Xuple.
func (a *XupleAction) Execute(ctx *rete.ActionContext) error {
    // 1. Récupérer le XupleManager depuis le contexte utilisateur
    xupleManagerIface, ok := ctx.UserContext["xupleManager"]
    if !ok {
        return fmt.Errorf("xuple manager not available in action context")
    }

    xupleManager, ok := xupleManagerIface.(*xuples.XupleManager)
    if !ok {
        return fmt.Errorf("invalid xuple manager type in context")
    }

    // 2. Résoudre le nom du xuple-space
    spaceNameIface, err := a.executor.ResolveArgument(a.SpaceNameArg, ctx)
    if err != nil {
        return fmt.Errorf("resolving xuple-space name: %w", err)
    }

    spaceName, ok := spaceNameIface.(string)
    if !ok {
        return fmt.Errorf("xuple-space name must be a string, got %T", spaceNameIface)
    }

    // 3. Vérifier que le xuple-space existe
    space := xupleManager.GetSpace(spaceName)
    if space == nil {
        return fmt.Errorf("xuple-space '%s' does not exist", spaceName)
    }

    // 4. Résoudre le fait
    factIface, err := a.executor.ResolveArgument(a.FactArg, ctx)
    if err != nil {
        return fmt.Errorf("resolving fact argument: %w", err)
    }

    fact, ok := factIface.(*rete.Fact)
    if !ok {
        return fmt.Errorf("second argument to Xuple must be a fact, got %T", factIface)
    }

    // 5. Convertir le fait RETE en Xuple
    xuple := convertFactToXuple(fact)

    // 6. Déposer le xuple dans le xuple-space
    if err := space.Deposit(xuple); err != nil {
        return fmt.Errorf("depositing xuple in '%s': %w", spaceName, err)
    }

    return nil
}

// convertFactToXuple convertit un fait RETE en Xuple.
func convertFactToXuple(fact *rete.Fact) *xuples.Xuple {
    return xuples.NewXuple(fact.Type, fact.Data)
}
```

#### 2.2 Factory pour XupleAction

```go
// CreateXupleActionFactory crée une factory pour les actions Xuple.
// Cette factory est enregistrée dans le réseau RETE sous le nom "Xuple".
func CreateXupleActionFactory(executor *rete.ActionExecutor) rete.ActionFactory {
    return func(args []rete.ActionArgument) (rete.Action, error) {
        // Vérifier le nombre d'arguments
        if len(args) != 2 {
            return nil, fmt.Errorf("Xuple action expects 2 arguments (space name, fact), got %d", len(args))
        }

        // Le premier argument doit être le nom du xuple-space (string)
        spaceNameArg := args[0]

        // Le second argument doit être un fait (inline ou référence)
        factArg := args[1]

        return NewXupleAction(spaceNameArg, factArg, executor), nil
    }
}
```

---

### Tâche 3: Enregistrer Automatiquement l'Action Xuple dans le Pipeline

**Fichier :** `tsd/api/pipeline.go`

**Objectif :** Configurer automatiquement l'action Xuple lors de l'initialisation du pipeline.

#### 3.1 Modification de `NewPipelineWithConfig`

```go
// NewPipelineWithConfig crée un nouveau pipeline avec une configuration personnalisée.
func NewPipelineWithConfig(config *Config) (*Pipeline, error) {
    // Valider la configuration
    if err := config.Validate(); err != nil {
        return nil, &ConfigError{
            Field:   "config",
            Message: err.Error(),
        }
    }

    // 1. Créer le réseau RETE
    network := rete.NewNetwork()

    // 2. Créer le storage (si nécessaire)
    storage := rete.NewMemoryStorage()

    // 3. Créer le XupleManager
    xupleManager := xuples.NewXupleManager()

    // 4. Créer le ConstraintPipeline
    retePipeline := constraint.NewConstraintPipeline(network, storage)

    // 5. Créer l'ActionExecutor
    executor := rete.NewActionExecutor(network)

    // 6. Injecter le XupleManager dans le contexte de l'executor
    executor.SetUserContext("xupleManager", xupleManager)

    // 7. Enregistrer l'action Xuple
    xupleFactory := CreateXupleActionFactory(executor)
    network.RegisterAction("Xuple", xupleFactory)

    // 8. Créer le pipeline
    p := &Pipeline{
        config:       config,
        network:      network,
        storage:      storage,
        xupleManager: xupleManager,
        retePipeline: retePipeline,
    }

    return p, nil
}
```

#### 3.2 Assurer que l'Executor Utilise le UserContext

**Fichier :** `tsd/internal/rete/network.go`

```go
// Network doit stocker et utiliser l'ActionExecutor configuré

type Network struct {
    // ... champs existants ...
    
    executor *ActionExecutor
}

// SetActionExecutor définit l'exécuteur d'actions pour ce réseau.
func (n *Network) SetActionExecutor(executor *ActionExecutor) {
    n.executor = executor
}

// GetActionExecutor retourne l'exécuteur d'actions.
func (n *Network) GetActionExecutor() *ActionExecutor {
    return n.executor
}

// RegisterAction enregistre une factory d'action sous un nom donné.
func (n *Network) RegisterAction(name string, factory ActionFactory) {
    if n.actionFactories == nil {
        n.actionFactories = make(map[string]ActionFactory)
    }
    n.actionFactories[name] = factory
}

// CreateAction crée une action à partir de son nom et de ses arguments.
func (n *Network) CreateAction(name string, args []ActionArgument) (Action, error) {
    factory, ok := n.actionFactories[name]
    if !ok {
        return nil, fmt.Errorf("unknown action: %s", name)
    }
    return factory(args)
}
```

#### 3.3 Modification du Pipeline API

```go
// NewPipelineWithConfig (version complète)
func NewPipelineWithConfig(config *Config) (*Pipeline, error) {
    if err := config.Validate(); err != nil {
        return nil, &ConfigError{
            Field:   "config",
            Message: err.Error(),
        }
    }

    // Créer le réseau RETE
    network := rete.NewNetwork()

    // Créer le storage
    storage := rete.NewMemoryStorage()

    // Créer le XupleManager
    xupleManager := xuples.NewXupleManager()

    // Créer l'ActionExecutor
    executor := rete.NewActionExecutor(network)
    executor.SetUserContext("xupleManager", xupleManager)

    // Associer l'executor au network
    network.SetActionExecutor(executor)

    // Enregistrer l'action Xuple
    xupleFactory := CreateXupleActionFactory(executor)
    network.RegisterAction("Xuple", xupleFactory)

    // Créer le ConstraintPipeline
    retePipeline := constraint.NewConstraintPipeline(network, storage)

    p := &Pipeline{
        config:       config,
        network:      network,
        storage:      storage,
        xupleManager: xupleManager,
        retePipeline: retePipeline,
    }

    return p, nil
}
```

---

### Tâche 4: Conversion des Actions dans AST → RETE

**Fichier :** `tsd/internal/constraint/rete_converter.go`

**Objectif :** S'assurer que les actions `Xuple()` sont correctement converties en ActionArguments.

#### 4.1 Conversion d'une ActionNode

```go
// convertAction convertit un ActionNode en rete.Action.
// Cette méthode est appelée lors de la conversion d'une RuleDef.
func (c *ASTConverter) convertAction(node *ActionNode) (rete.Action, error) {
    // Convertir les arguments
    args := make([]rete.ActionArgument, len(node.Arguments))
    for i, argNode := range node.Arguments {
        arg, err := c.convertActionArgument(argNode)
        if err != nil {
            return nil, c.wrapError(node, fmt.Sprintf("converting argument %d", i), err)
        }
        args[i] = arg
    }

    // Créer l'action via le network
    action, err := c.network.CreateAction(node.Name, args)
    if err != nil {
        return nil, c.wrapError(node, fmt.Sprintf("creating action '%s'", node.Name), err)
    }

    return action, nil
}

// convertActionArgument convertit un nœud d'expression en ActionArgument.
func (c *ASTConverter) convertActionArgument(node ExpressionNode) (rete.ActionArgument, error) {
    switch n := node.(type) {
    case *LiteralNode:
        return rete.LiteralArgument{Value: n.Value}, nil

    case *FieldReferenceNode:
        return rete.FieldReferenceArgument{
            Variable:  n.Variable,
            FieldPath: n.Path,
        }, nil

    case *InlineFactNode:
        return c.convertInlineFactArgument(n)

    default:
        return nil, fmt.Errorf("unsupported argument type: %T", n)
    }
}

// convertInlineFactArgument convertit un InlineFactNode en InlineFactArgument.
func (c *ASTConverter) convertInlineFactArgument(node *InlineFactNode) (rete.ActionArgument, error) {
    fields := make(map[string]rete.ActionArgument)

    for fieldName, fieldNode := range node.Fields {
        arg, err := c.convertActionArgument(fieldNode)
        if err != nil {
            return nil, fmt.Errorf("converting field %s: %w", fieldName, err)
        }
        fields[fieldName] = arg
    }

    return rete.InlineFactArgument{
        TypeName: node.TypeName,
        Fields:   fields,
    }, nil
}
```

---

### Tâche 5: Exécution des Actions lors de l'Activation de Règles

**Fichier :** `tsd/internal/rete/network.go`

**Objectif :** S'assurer que les actions sont exécutées avec le bon contexte.

#### 5.1 Activation de Règle

```go
// Rule représente une règle dans le réseau RETE.
type Rule struct {
    Name       string
    Conditions []Condition
    Actions    []Action
    Priority   int
}

// Activate active une règle avec les faits déclencheurs.
func (n *Network) ActivateRule(rule *Rule, triggeringFacts map[string]*Fact, bindings map[string]interface{}) error {
    // Récupérer l'executor
    executor := n.GetActionExecutor()
    if executor == nil {
        return fmt.Errorf("action executor not configured")
    }

    // Exécuter chaque action
    for i, action := range rule.Actions {
        if err := executor.Execute(action, triggeringFacts, bindings, rule.Name); err != nil {
            return fmt.Errorf("executing action %d in rule '%s': %w", i, rule.Name, err)
        }
    }

    return nil
}
```

---

### Tâche 6: Gestion des Erreurs Spécifiques à Xuple

**Fichier :** `tsd/api/errors.go`

**Objectif :** Ajouter des types d'erreurs spécifiques aux actions Xuple.

```go
// XupleActionError représente une erreur lors de l'exécution d'une action Xuple.
type XupleActionError struct {
    SpaceName string
    RuleName  string
    Message   string
    Cause     error
}

func (e *XupleActionError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("xuple action error in rule '%s' for space '%s': %s: %v",
            e.RuleName, e.SpaceName, e.Message, e.Cause)
    }
    return fmt.Sprintf("xuple action error in rule '%s' for space '%s': %s",
        e.RuleName, e.SpaceName, e.Message)
}

func (e *XupleActionError) Unwrap() error {
    return e.Cause
}
```

Modifier `XupleAction.Execute()` pour utiliser cette erreur :

```go
func (a *XupleAction) Execute(ctx *rete.ActionContext) error {
    // ... code existant ...

    // En cas d'erreur lors du dépôt
    if err := space.Deposit(xuple); err != nil {
        return &XupleActionError{
            SpaceName: spaceName,
            RuleName:  ctx.RuleName,
            Message:   "failed to deposit xuple",
            Cause:     err,
        }
    }

    return nil
}
```

---

## 🧪 Tests à Implémenter

### Test 1: Action Xuple - Création et Enregistrement

**Fichier :** `tsd/api/xuple_action_test.go`

```go
func TestXupleAction_Create(t *testing.T) {
    network := rete.NewNetwork()
    executor := rete.NewActionExecutor(network)
    xupleManager := xuples.NewXupleManager()

    executor.SetUserContext("xupleManager", xupleManager)

    // Créer un xuple-space
    space, err := xupleManager.CreateXupleSpace("alerts", xuples.SelectionFIFO, xuples.ConsumptionOnce, xuples.RetentionUnlimited)
    require.NoError(t, err)

    // Créer l'action Xuple
    spaceNameArg := rete.LiteralArgument{Value: "alerts"}
    factArg := rete.InlineFactArgument{
        TypeName: "Alert",
        Fields: map[string]rete.ActionArgument{
            "message": rete.LiteralArgument{Value: "Test alert"},
            "level":   rete.LiteralArgument{Value: "high"},
        },
    }

    action := NewXupleAction(spaceNameArg, factArg, executor)
    require.NotNil(t, action)
    assert.Equal(t, "Xuple", action.Name())
}
```

### Test 2: Action Xuple - Exécution Simple

```go
func TestXupleAction_Execute_Simple(t *testing.T) {
    network := rete.NewNetwork()
    executor := rete.NewActionExecutor(network)
    xupleManager := xuples.NewXupleManager()

    executor.SetUserContext("xupleManager", xupleManager)

    // Créer un xuple-space
    space, err := xupleManager.CreateXupleSpace("alerts", xuples.SelectionFIFO, xuples.ConsumptionOnce, xuples.RetentionUnlimited)
    require.NoError(t, err)

    // Enregistrer le type Alert dans le network
    network.RegisterType("Alert", map[string]string{
        "message": "string",
        "level":   "string",
    })

    // Créer l'action
    spaceNameArg := rete.LiteralArgument{Value: "alerts"}
    factArg := rete.InlineFactArgument{
        TypeName: "Alert",
        Fields: map[string]rete.ActionArgument{
            "message": rete.LiteralArgument{Value: "Test alert"},
            "level":   rete.LiteralArgument{Value: "high"},
        },
    }

    action := NewXupleAction(spaceNameArg, factArg, executor)

    // Créer le contexte d'exécution
    ctx := &rete.ActionContext{
        Network:         network,
        TriggeringFacts: make(map[string]*rete.Fact),
        Bindings:        make(map[string]interface{}),
        RuleName:        "TestRule",
        UserContext:     executor.UserContext,
    }

    // Exécuter l'action
    err = action.Execute(ctx)
    require.NoError(t, err)

    // Vérifier que le xuple a été créé
    xuples := space.GetAll()
    require.Len(t, xuples, 1)
    assert.Equal(t, "Alert", xuples[0].Type())
    assert.Equal(t, "Test alert", xuples[0].Get("message"))
    assert.Equal(t, "high", xuples[0].Get("level"))
}
```

### Test 3: Action Xuple - Avec Références aux Champs

```go
func TestXupleAction_Execute_WithFieldReferences(t *testing.T) {
    network := rete.NewNetwork()
    executor := rete.NewActionExecutor(network)
    xupleManager := xuples.NewXupleManager()

    executor.SetUserContext("xupleManager", xupleManager)

    // Créer un xuple-space
    space, err := xupleManager.CreateXupleSpace("alerts", xuples.SelectionFIFO, xuples.ConsumptionOnce, xuples.RetentionUnlimited)
    require.NoError(t, err)

    // Enregistrer les types
    network.RegisterType("Temperature", map[string]string{
        "sensorId": "string",
        "value":    "float",
    })
    network.RegisterType("Alert", map[string]string{
        "sensorId": "string",
        "message":  "string",
        "temp":     "float",
    })

    // Créer un fait Temperature
    tempFact := network.CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.5,
    })

    // Créer l'action avec des références aux champs
    spaceNameArg := rete.LiteralArgument{Value: "alerts"}
    factArg := rete.InlineFactArgument{
        TypeName: "Alert",
        Fields: map[string]rete.ActionArgument{
            "sensorId": rete.FieldReferenceArgument{
                Variable:  "t",
                FieldPath: []string{"sensorId"},
            },
            "message": rete.LiteralArgument{Value: "High temperature"},
            "temp": rete.FieldReferenceArgument{
                Variable:  "t",
                FieldPath: []string{"value"},
            },
        },
    }

    action := NewXupleAction(spaceNameArg, factArg, executor)

    // Créer le contexte avec le fait déclencheur
    ctx := &rete.ActionContext{
        Network: network,
        TriggeringFacts: map[string]*rete.Fact{
            "t": tempFact,
        },
        Bindings: map[string]interface{}{
            "t": tempFact,
        },
        RuleName:    "HighTemperature",
        UserContext: map[string]interface{}{"xupleManager": xupleManager},
    }

    // Exécuter l'action
    err = action.Execute(ctx)
    require.NoError(t, err)

    // Vérifier le xuple créé
    xuples := space.GetAll()
    require.Len(t, xuples, 1)
    assert.Equal(t, "Alert", xuples[0].Type())
    assert.Equal(t, "sensor-01", xuples[0].Get("sensorId"))
    assert.Equal(t, "High temperature", xuples[0].Get("message"))
    assert.Equal(t, 35.5, xuples[0].Get("temp"))
}
```

### Test 4: Pipeline - Action Xuple Automatique

```go
func TestPipeline_XupleAction_Automatic(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    message: string,
    temp: float
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            message: "High temperature detected",
            temp: t.value
        ))
    }
}
`

    tmpfile, err := os.CreateTemp("", "test*.tsd")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())

    _, err = tmpfile.WriteString(tsdContent)
    require.NoError(t, err)
    tmpfile.Close()

    // Créer le pipeline
    pipeline, err := NewPipeline()
    require.NoError(t, err)

    // Ingérer le fichier
    result, err := pipeline.IngestFile(tmpfile.Name())
    require.NoError(t, err)

    // Vérifier que le xuple-space existe
    spaces := result.XupleSpaceNames()
    require.Contains(t, spaces, "alerts")

    // Soumettre un fait Temperature
    tempFact := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.5,
    })
    result.Network().Assert(tempFact)

    // Vérifier qu'un xuple a été créé automatiquement
    xuples := result.GetXuples("alerts")
    require.Len(t, xuples, 1)

    alert := xuples[0]
    assert.Equal(t, "Alert", alert.Type())
    assert.Equal(t, "sensor-01", alert.Get("sensorId"))
    assert.Equal(t, "High temperature detected", alert.Get("message"))
    assert.Equal(t, 35.5, alert.Get("temp"))
}
```

### Test 5: Erreurs - Xuple-Space Inexistant

```go
func TestXupleAction_Error_SpaceNotFound(t *testing.T) {
    network := rete.NewNetwork()
    executor := rete.NewActionExecutor(network)
    xupleManager := xuples.NewXupleManager()

    executor.SetUserContext("xupleManager", xupleManager)

    // NE PAS créer le xuple-space "nonexistent"

    // Créer l'action qui référence un space inexistant
    spaceNameArg := rete.LiteralArgument{Value: "nonexistent"}
    factArg := rete.InlineFactArgument{
        TypeName: "Alert",
        Fields: map[string]rete.ActionArgument{
            "message": rete.LiteralArgument{Value: "Test"},
        },
    }

    action := NewXupleAction(spaceNameArg, factArg, executor)

    // Contexte
    ctx := &rete.ActionContext{
        Network:         network,
        TriggeringFacts: make(map[string]*rete.Fact),
        Bindings:        make(map[string]interface{}),
        RuleName:        "TestRule",
        UserContext:     map[string]interface{}{"xupleManager": xupleManager},
    }

    // Exécuter (doit échouer)
    err := action.Execute(ctx)
    require.Error(t, err)
    assert.Contains(t, err.Error(), "does not exist")
}
```

### Test 6: E2E - Multiples Règles, Multiples Xuple-Spaces

```go
func TestE2E_MultipleRules_MultipleSpaces(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo
}

xuple-space logs {
    selection: lifo,
    max-size: 100
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    level: string,
    temp: float
}

type LogEntry {
    source: string,
    message: string
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            level: "high",
            temp: t.value
        )),
        Xuple("logs", LogEntry(
            source: t.sensorId,
            message: "High temp recorded"
        ))
    }
}

rule VeryHighTemperature {
    when {
        t: Temperature(value > 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            level: "critical",
            temp: t.value
        ))
    }
}
`

    tmpfile, err := os.CreateTemp("", "test*.tsd")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())

    _, err = tmpfile.WriteString(tsdContent)
    require.NoError(t, err)
    tmpfile.Close()

    // Pipeline
    pipeline, err := NewPipeline()
    require.NoError(t, err)

    result, err := pipeline.IngestFile(tmpfile.Name())
    require.NoError(t, err)

    // Vérifier les xuple-spaces
    spaces := result.XupleSpaceNames()
    require.Len(t, spaces, 2)
    require.Contains(t, spaces, "alerts")
    require.Contains(t, spaces, "logs")

    // Soumettre un fait Temperature = 35°
    temp1 := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.0,
    })
    result.Network().Assert(temp1)

    // Vérifier: 1 alert (high), 1 log
    alerts := result.GetXuples("alerts")
    require.Len(t, alerts, 1)
    assert.Equal(t, "high", alerts[0].Get("level"))

    logs := result.GetXuples("logs")
    require.Len(t, logs, 1)

    // Soumettre Temperature = 45° (déclenche les 2 règles)
    temp2 := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-02",
        "value":    45.0,
    })
    result.Network().Assert(temp2)

    // Vérifier: 3 alerts au total (1 high + 1 high + 1 critical)
    alerts = result.GetXuples("alerts")
    require.Len(t, alerts, 3)

    // Vérifier: 2 logs
    logs = result.GetXuples("logs")
    require.Len(t, logs, 2)
}
```

---

## ✅ Checklist de Validation

### Interface Action

- [ ] `Action` interface définie avec `Execute(ctx *ActionContext) error`
- [ ] `ActionContext` contient `UserContext map[string]interface{}`
- [ ] `ActionExecutor` implémente `SetUserContext(key, value)`
- [ ] `ActionExecutor.ResolveArgument()` gère tous les types d'arguments
- [ ] `ActionArgument` interface avec 3 implémentations (Literal, FieldReference, InlineFact)

### XupleAction

- [ ] `XupleAction` implémentée dans `tsd/api/xuple_action.go`
- [ ] `NewXupleAction()` crée l'action correctement
- [ ] `Execute()` récupère le XupleManager depuis le contexte
- [ ] `Execute()` résout le nom du xuple-space
- [ ] `Execute()` résout le fait inline
- [ ] `Execute()` vérifie l'existence du xuple-space
- [ ] `Execute()` dépose le xuple dans le space
- [ ] Conversion Fact → Xuple fonctionne
- [ ] Erreurs claires et typées (`XupleActionError`)

### Enregistrement Automatique

- [ ] `CreateXupleActionFactory()` crée une factory valide
- [ ] `NewPipelineWithConfig()` crée l'ActionExecutor
- [ ] Pipeline injecte le XupleManager dans le UserContext
- [ ] Pipeline enregistre l'action "Xuple" dans le network
- [ ] `Network.RegisterAction()` fonctionne
- [ ] `Network.CreateAction()` utilise les factories

### Conversion AST → RETE

- [ ] `convertAction()` convertit les ActionNode
- [ ] `convertActionArgument()` gère tous les types d'expressions
- [ ] `convertInlineFactArgument()` convertit les faits inline récursivement
- [ ] Arguments correctement passés à la factory

### Exécution

- [ ] `Network.ActivateRule()` utilise l'ActionExecutor
- [ ] Contexte correctement passé aux actions
- [ ] Multiples actions dans une règle fonctionnent
- [ ] Erreurs d'exécution remontées correctement

### Tests

- [ ] Tests unitaires de `XupleAction` (création, exécution simple)
- [ ] Tests avec références aux champs
- [ ] Tests d'erreurs (space inexistant, arguments invalides)
- [ ] Tests d'intégration avec le pipeline
- [ ] Tests E2E complets (multiples règles, multiples spaces)
- [ ] Couverture > 80%

### Standards

- [ ] Code formaté (`gofmt`)
- [ ] Pas de warnings du linter
- [ ] Commentaires GoDoc complets
- [ ] Exemples d'utilisation

---

## 📝 Documentation à Mettre à Jour

### 1. Guide TSD (`docs/TSD_LANGUAGE.md`)

Ajouter section sur l'action Xuple :

```markdown
## Actions

### Action Xuple

L'action `Xuple()` crée un xuple dans un xuple-space.

#### Syntaxe

\`\`\`tsd
Xuple("<space-name>", FactType(...))
\`\`\`

#### Paramètres

1. **space-name** (string) : Nom du xuple-space cible
2. **fact** : Fait à stocker comme xuple (inline ou référence)

#### Exemples

**Exemple 1 : Fait inline simple**

\`\`\`tsd
rule CreateAlert {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            message: "Temperature too high",
            level: "warning"
        ))
    }
}
\`\`\`

**Exemple 2 : Avec références aux champs**

\`\`\`tsd
rule SensorAlert {
    when {
        s: Sensor(status == "error")
    }
    then {
        Xuple("notifications", Notification(
            sensorId: s.id,
            message: "Sensor error detected",
            timestamp: s.lastUpdate
        ))
    }
}
\`\`\`

**Exemple 3 : Multiples actions**

\`\`\`tsd
rule CriticalEvent {
    when {
        e: Event(severity == "critical")
    }
    then {
        Xuple("alerts", Alert(event: e.id, level: "critical")),
        Xuple("logs", LogEntry(source: e.source, message: e.description))
    }
}
\`\`\`

#### Notes

- Le xuple-space doit être défini avant utilisation
- Le type de fait doit être enregistré
- Les références aux champs utilisent la notation pointée (ex: `s.id`)
```

### 2. Guide API (`docs/API_USAGE.md`)

```markdown
## Actions Automatiques

Les actions définies dans les règles TSD sont **automatiquement configurées**
lors de l'ingestion du fichier.

### Action Xuple

L'action `Xuple()` est automatiquement disponible dans toutes les règles.

**Aucune configuration nécessaire** - elle fonctionne immédiatement après
`Pipeline.IngestFile()`.

#### Exemple Complet

Fichier TSD (`rules.tsd`):
\`\`\`tsd
xuple-space alerts {
    selection: fifo
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    temp: float
}

rule HighTemp {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(sensorId: t.sensorId, temp: t.value))
    }
}
\`\`\`

Code Go:
\`\`\`go
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")

// Soumettre un fait
temp := result.Network().CreateFact("Temperature", map[string]interface{}{
    "sensorId": "sensor-01",
    "value":    35.5,
})
result.Network().Assert(temp)

// Récupérer les xuples créés automatiquement
xuples := result.GetXuples("alerts")
// xuples contient 1 Alert avec sensorId="sensor-01", temp=35.5
\`\`\`

**C'est tout !** Aucun code supplémentaire requis.
```

### 3. Architecture (`docs/ARCHITECTURE.md`)

Ajouter section sur le flux d'exécution des actions :

```markdown
## Flux d'Exécution des Actions

\`\`\`
1. Parsing TSD
   - Détection des actions dans les règles
   - Création d'ActionNode dans l'AST
   ↓
2. Conversion AST → RETE
   - Conversion des arguments (literal, field ref, inline fact)
   - Création des actions via ActionFactory
   ↓
3. Enregistrement dans le Network
   - Actions stockées dans les règles
   ↓
4. Propagation de faits
   ↓
5. Activation de règles
   ↓
6. Exécution des actions
   - ActionExecutor.Execute(action, context)
   - Résolution des arguments (champs, faits inline)
   - Exécution de l'action (ex: dépôt de xuple)
\`\`\`

### Configuration Automatique (Action Xuple)

- **Où** : `api.NewPipelineWithConfig()`
- **Quand** : À l'initialisation du pipeline
- **Comment** :
  1. Création de l'ActionExecutor
  2. Injection du XupleManager dans UserContext
  3. Enregistrement de la factory "Xuple"
  4. Association de l'executor au network
```

---

## 🎯 Résultat Attendu

### Avant (avec configuration manuelle)

```go
// Configuration manuelle requise
network := rete.NewNetwork()
executor := rete.NewActionExecutor(network)
xupleManager := xuples.NewXupleManager()

executor.SetUserContext("xupleManager", xupleManager)
network.RegisterAction("Xuple", CreateXupleActionFactory(executor))

// Utilisation...
```

### Après (automatique)

```go
// Rien à configurer !
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")

// L'action Xuple fonctionne immédiatement
// dans toutes les règles du fichier TSD
```

**Bénéfice principal :** L'utilisateur n'a **aucune connaissance** à avoir de
l'implémentation interne des actions. Il écrit simplement `Xuple(...)` dans
ses règles et ça fonctionne.

---

## 🔗 Dépendances

### Entrantes

- ✅ Prompt 01 : Parser supporte `Xuple("space", Fact(...))`
- ✅ Prompt 02 : Package `api` avec `Pipeline`
- ✅ Prompt 03 : Xuple-spaces créés automatiquement

### Sortantes

- ➡️ Prompt 05 : Migration des tests E2E (utilisera les actions automatiques)
- ➡️ Prompt 06 : Cleanup (supprimera l'ancien code de configuration manuelle)

---

## 🚀 Stratégie d'Implémentation

1. **Phase 1: Interface Action** (1h)
   - Définir `ActionContext` avec `UserContext`
   - Implémenter `ActionExecutor.SetUserContext()`
   - Ajouter `ResolveArgument()` avec tous les types

2. **Phase 2: XupleAction** (1-2h)
   - Créer `tsd/api/xuple_action.go`
   - Implémenter `XupleAction.Execute()`
   - Factory `CreateXupleActionFactory()`
   - Conversion Fact → Xuple

3. **Phase 3: Enregistrement Automatique** (1h)
   - Modifier `NewPipelineWithConfig()`
   - Créer et injecter l'ActionExecutor
   - Enregistrer l'action "Xuple"
   - Modifier `Network` pour stocker l'executor

4. **Phase 4: Conversion AST** (30min)
   - `convertAction()` dans le converter
   - `convertActionArgument()` pour tous les types
   - Tests de conversion

5. **Phase 5: Tests** (2h)
   - Tests unitaires de XupleAction
   - Tests d'intégration avec pipeline
   - Tests E2E complets

6. **Phase 6: Documentation** (30min)
   - Mise à jour des guides
   - Exemples GoDoc

**Estimation totale : 6-7 heures**

---

## 📊 Critères de Succès

- [ ] `XupleAction` implémentée et fonctionnelle
- [ ] Action enregistrée automatiquement dans le pipeline
- [ ] Aucune configuration manuelle requise
- [ ] Résolution des arguments fonctionne (literal, field ref, inline fact)
- [ ] Multiples actions dans une règle supportées
- [ ] Erreurs claires et typées
- [ ] Tests unitaires passent (couverture > 80%)
- [ ] Tests E2E passent
- [ ] Documentation complète et à jour
- [ ] Pas de régression dans les tests existants

---

**FIN DU PROMPT 04**