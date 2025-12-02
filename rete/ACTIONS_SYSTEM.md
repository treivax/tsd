# Système d'Actions Personnalisables

## Vue d'ensemble

Le système d'actions de TSD permet de définir des comportements personnalisés qui s'exécutent lorsque les règles sont déclenchées. Chaque action peut avoir son propre comportement défini via un `ActionHandler`.

Le système supporte également l'utilisation d'**expressions arithmétiques** pour calculer dynamiquement des valeurs lors de la création ou modification de faits. Voir la section [Expressions Arithmétiques](#expressions-arithmétiques) pour plus de détails.

## Architecture

### Composants principaux

1. **ActionHandler** : Interface pour définir le comportement d'une action
2. **ActionRegistry** : Gestionnaire d'enregistrement des handlers
3. **ActionExecutor** : Exécuteur d'actions intégré au moteur RETE
4. **PrintAction** : Première action intégrée (affichage)

### Flux d'exécution

```
Règle déclenchée → ActionExecutor → Registry → ActionHandler → Exécution
                                   ↓
                            Si non trouvé → Log uniquement
```

## Utilisation

### Utiliser l'action print (intégrée)

```go
// Créer un réseau RETE
storage := NewMemoryStorage()
network := NewReteNetwork(storage)

// L'action print est automatiquement enregistrée
// Elle peut être utilisée dans les règles

// Créer une action print
action := &Action{
    Type: "action",
    Jobs: []JobCall{
        {
            Type: "jobCall",
            Name: "print",
            Args: []interface{}{
                map[string]interface{}{
                    "type":  "string",
                    "value": "Hello, World!",
                },
            },
        },
    },
}

// Exécuter l'action
token := &Token{
    ID:    "token1",
    Facts: []*Fact{fact},
    Bindings: map[string]*Fact{
        "p": fact,
    },
}

err := network.ActionExecutor.ExecuteAction(action, token)
```

### Créer une action personnalisée

```go
// 1. Définir la constante pour le nom
const ActionNameCustom = "custom_action"

// 2. Implémenter l'interface ActionHandler
type CustomAction struct {
    config Config
}

func NewCustomAction(config Config) *CustomAction {
    return &CustomAction{config: config}
}

func (ca *CustomAction) Execute(args []interface{}, ctx *ExecutionContext) error {
    // Votre logique métier ici
    if len(args) == 0 {
        return fmt.Errorf("custom_action requires at least one argument")
    }
    
    // Traiter les arguments
    value := args[0]
    
    // Effectuer l'action
    fmt.Printf("Custom action executed with: %v\n", value)
    
    return nil
}

func (ca *CustomAction) GetName() string {
    return ActionNameCustom
}

func (ca *CustomAction) Validate(args []interface{}) error {
    if len(args) == 0 {
        return fmt.Errorf("custom_action requires at least one argument")
    }
    return nil
}

// 3. Enregistrer l'action
customAction := NewCustomAction(config)
err := network.ActionExecutor.RegisterAction(customAction)
if err != nil {
    log.Fatal(err)
}
```

## Actions non définies

Les actions qui n'ont pas de handler enregistré sont simplement loguées sans causer d'erreur. Cela permet :

- De tester des règles avant d'implémenter les actions
- De maintenir la compatibilité avec des règles utilisant des actions non encore implémentées
- De déboguer facilement les actions appelées

Exemple de log pour une action non définie :
```
📋 ACTION NON DÉFINIE (log uniquement): my_undefined_action("param1", 42)
```

## Action Print

### Description

L'action `print` affiche une chaîne de caractères sur la sortie standard (ou un writer personnalisé).

### Signature

```
print(message: any)
```

### Types d'arguments supportés

- **string** : Affiché tel quel
- **number** : Converti en chaîne
- **boolean** : Converti en "true" ou "false"
- **Fact** : Affiché avec sa structure (type, ID, champs)
- **fieldAccess** : Valeur du champ extraite et affichée
- **variable** : Fait complet affiché

### Exemples

#### Afficher une chaîne littérale

```go
action := &Action{
    Jobs: []JobCall{
        {
            Name: "print",
            Args: []interface{}{
                map[string]interface{}{
                    "type":  "string",
                    "value": "Hello, World!",
                },
            },
        },
    },
}
```

## Expressions Arithmétiques

### Vue d'ensemble

Le système d'actions supporte l'utilisation d'expressions arithmétiques directement dans les arguments d'actions. Cela permet de calculer dynamiquement des valeurs lors de la création ou modification de faits en utilisant les variables liées par la règle.

### Opérateurs supportés

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| `+` | Addition | `a.age + 5` |
| `-` | Soustraction | `a.age - e.age` |
| `*` | Multiplication | `p.price * p.quantity` |
| `/` | Division | `total / count` |
| `%` | Modulo | `value % 10` |

### Cas d'utilisation

#### 1. Création de fait avec calcul

```tsd
{ a: Adulte, e: Enfant } / a.age > e.age AND e.pere = a.ID 
==> setFact(
    Naissance(
        id: e.ID,
        parent: a.ID,
        ageParentALaNaissance: a.age - e.age
    )
)
```

#### 2. Modification de fait avec calcul

```tsd
{ p: Person } / p.age < 30 
==> setFact(p[bonus] = p.salary * 0.1)
```

#### 3. Expressions imbriquées

```tsd
{ prod: Product } / prod.available = true
==> setFact(
    Invoice(
        productId: prod.id,
        subtotal: prod.price * prod.quantity,
        total: (prod.price * prod.quantity) * 1.20
    )
)
```

### Format interne

Les expressions arithmétiques utilisent le type `"binaryOperation"` avec la structure suivante :

```json
{
    "type": "binaryOperation",
    "operator": "-",
    "left": {
        "type": "fieldAccess",
        "object": "a",
        "field": "age"
    },
    "right": {
        "type": "fieldAccess",
        "object": "e",
        "field": "age"
    }
}
```

### Gestion des erreurs

Le système gère automatiquement :
- **Division par zéro** : erreur levée lors de l'exécution
- **Modulo par zéro** : erreur levée lors de l'exécution
- **Types incompatibles** : erreur si les opérandes ne sont pas numériques
- **Validation de type** : les résultats doivent correspondre au type attendu du champ

### Documentation complète

Pour une documentation détaillée sur les expressions arithmétiques dans les actions, consultez :
- [ARITHMETIC_IN_ACTIONS.md](../docs/ARITHMETIC_IN_ACTIONS.md)

## Références

- [Guide rapide des actions](ACTIONS_QUICKSTART.md)
- [Résumé des fonctionnalités](ACTIONS_FEATURE_SUMMARY.md)
- [README des actions](ACTIONS_README.md)
- [Expressions arithmétiques](../docs/ARITHMETIC_IN_ACTIONS.md)

```go
{
    Type: "jobCall",
    Name: "print",
    Args: []interface{}{
        map[string]interface{}{
            "type":  "string",
            "value": "Hello, World!",
        },
    },
}
// Sortie: Hello, World!
```

#### Afficher un champ d'un fait

```go
{
    Type: "jobCall",
    Name: "print",
    Args: []interface{}{
        map[string]interface{}{
            "type":   "fieldAccess",
            "object": "p",
            "field":  "name",
        },
    },
}
// Sortie: Alice (si p.name == "Alice")
```

#### Afficher un fait complet

```go
{
    Type: "jobCall",
    Name: "print",
    Args: []interface{}{
        map[string]interface{}{
            "type": "variable",
            "name": "p",
        },
    },
}
// Sortie: Person{id: person_1, name: "Alice", age: 25}
```

#### Afficher un nombre

```go
{
    Type: "jobCall",
    Name: "print",
    Args: []interface{}{
        map[string]interface{}{
            "type":  "number",
            "value": 42.5,
        },
    },
}
// Sortie: 42.5
```

### Personnaliser la sortie

```go
// Créer un buffer ou un fichier
var output bytes.Buffer
printAction := NewPrintAction(&output)

// Enregistrer dans le registry
network.ActionExecutor.GetRegistry().Register(printAction)

// Ou changer dynamiquement
printAction.SetOutput(os.Stderr)
```

## ActionRegistry

### Méthodes disponibles

#### Register

Enregistre un handler d'action.

```go
err := registry.Register(handler)
```

#### Unregister

Supprime un handler du registry.

```go
registry.Unregister(actionName)
```

#### Get

Récupère un handler par son nom.

```go
handler := registry.Get(actionName)
```

#### Has

Vérifie si un handler est enregistré.

```go
if registry.Has(actionName) {
    // L'action est disponible
}
```

#### GetAll

Récupère tous les handlers enregistrés.

```go
allHandlers := registry.GetAll()
```

#### GetRegisteredNames

Récupère la liste des noms d'actions enregistrées.

```go
names := registry.GetRegisteredNames()
```

#### Clear

Supprime tous les handlers du registry.

```go
registry.Clear()
```

#### RegisterMultiple

Enregistre plusieurs handlers en une seule opération.

```go
handlers := []ActionHandler{action1, action2, action3}
err := registry.RegisterMultiple(handlers)
```

## ExecutionContext

Le contexte d'exécution fourni aux handlers contient :

- **token** : Le token RETE contenant les faits matchés
- **network** : Référence au réseau RETE
- **varCache** : Cache des variables disponibles

### Méthodes

#### GetVariable

Récupère un fait par nom de variable.

```go
fact := ctx.GetVariable("p")
if fact != nil {
    // Utiliser le fait
}
```

## Logging

Toutes les actions sont loguées automatiquement :

### Actions définies

```
📋 ACTION: print(p.name)
🎯 ACTION EXÉCUTÉE: print("Alice")
```

### Actions non définies

```
📋 ACTION: custom_action(p.id)
📋 ACTION NON DÉFINIE (log uniquement): custom_action("123")
```

### Désactiver le logging

```go
network.ActionExecutor.SetLogging(false)
```

## Validation

Les handlers peuvent implémenter une validation optionnelle via la méthode `Validate`.

```go
func (h *MyHandler) Validate(args []interface{}) error {
    if len(args) < 2 {
        return fmt.Errorf("my_action requires at least 2 arguments")
    }
    
    // Validation spécifique
    firstArg := args[0]
    if _, ok := firstArg.(string); !ok {
        return fmt.Errorf("first argument must be a string")
    }
    
    return nil
}
```

La validation est appelée automatiquement avant l'exécution.

## Bonnes pratiques

### 1. Utiliser des constantes pour les noms

```go
const ActionNameNotify = "notify"

func (na *NotifyAction) GetName() string {
    return ActionNameNotify
}
```

### 2. Valider les arguments

Toujours implémenter `Validate` pour vérifier les arguments avant l'exécution.

### 3. Gestion des erreurs

Retourner des erreurs explicites et descriptives.

```go
if value == nil {
    return fmt.Errorf("argument cannot be nil for action %s", h.GetName())
}
```

### 4. Thread-safety

Le registry est thread-safe. Les handlers doivent aussi l'être s'ils sont utilisés de manière concurrente.

### 5. Pas de hardcoding

Utiliser des paramètres pour toute configuration.

```go
// ❌ Mauvais
func (h *EmailAction) Execute(args []interface{}, ctx *ExecutionContext) error {
    smtp := "smtp.example.com:587" // Hardcodé !
    // ...
}

// ✅ Bon
type EmailAction struct {
    smtpHost string
    smtpPort int
}

func NewEmailAction(smtpHost string, smtpPort int) *EmailAction {
    return &EmailAction{
        smtpHost: smtpHost,
        smtpPort: smtpPort,
    }
}
```

### 6. Documentation

Documenter chaque handler avec des exemples d'utilisation.

## Tests

### Tester un handler

```go
func TestMyAction_Execute(t *testing.T) {
    // Arrange
    action := NewMyAction()
    ctx := NewExecutionContext(nil, nil)
    args := []interface{}{"test"}
    
    // Act
    err := action.Execute(args, ctx)
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
}
```

### Tester avec un mock

```go
type MockActionHandler struct {
    executeCalled bool
    lastArgs      []interface{}
}

func (m *MockActionHandler) Execute(args []interface{}, ctx *ExecutionContext) error {
    m.executeCalled = true
    m.lastArgs = args
    return nil
}

func (m *MockActionHandler) GetName() string {
    return "mock"
}

func (m *MockActionHandler) Validate(args []interface{}) error {
    return nil
}
```

## Exemples d'actions utiles

### Action de notification

```go
type NotifyAction struct {
    notifier Notifier
}

func (na *NotifyAction) Execute(args []interface{}, ctx *ExecutionContext) error {
    if len(args) < 2 {
        return fmt.Errorf("notify requires 2 arguments: channel and message")
    }
    
    channel := args[0].(string)
    message := args[1].(string)
    
    return na.notifier.Send(channel, message)
}
```

### Action d'assertion de fait

```go
type AssertAction struct {
    network *ReteNetwork
}

func (aa *AssertAction) Execute(args []interface{}, ctx *ExecutionContext) error {
    if len(args) == 0 {
        return fmt.Errorf("assert requires at least one fact")
    }
    
    fact, ok := args[0].(*Fact)
    if !ok {
        return fmt.Errorf("argument must be a Fact")
    }
    
    return aa.network.AssertFact(fact)
}
```

### Action de retrait de fait

```go
type RetractAction struct {
    network *ReteNetwork
}

func (ra *RetractAction) Execute(args []interface{}, ctx *ExecutionContext) error {
    if len(args) == 0 {
        return fmt.Errorf("retract requires a fact")
    }
    
    fact, ok := args[0].(*Fact)
    if !ok {
        return fmt.Errorf("argument must be a Fact")
    }
    
    return ra.network.RetractFact(fact)
}
```

## Références

- **action_handler.go** : Interface et registry
- **action_print.go** : Implémentation de l'action print
- **action_executor.go** : Exécuteur d'actions
- **action_handler_test.go** : Tests unitaires
- **action_print_integration_test.go** : Tests d'intégration

## Feuille de route

Actions futures à implémenter :

- [ ] `assert(fact)` : Assertion de nouveau fait
- [ ] `retract(fact)` : Retrait de fait
- [ ] `modify(fact, field, value)` : Modification de fait
- [ ] `log(level, message)` : Logging avec niveaux
- [ ] `http(method, url, body)` : Appel HTTP
- [ ] `emit(event, data)` : Émission d'événement
- [ ] `delay(duration, action)` : Action différée