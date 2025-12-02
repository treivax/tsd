# Système d'Actions TSD

## 🎯 Vue d'ensemble

Le système d'actions TSD permet de définir des comportements personnalisés qui s'exécutent lorsque les règles du moteur RETE sont déclenchées. Chaque action peut avoir son propre comportement défini via un handler personnalisable.

## ✨ Fonctionnalités

- ✅ **Actions personnalisables** : Créez vos propres actions avec des comportements spécifiques
- ✅ **Action print intégrée** : Action d'affichage prête à l'emploi
- ✅ **Registry thread-safe** : Gestion sécurisée des handlers d'actions
- ✅ **Logging automatique** : Toutes les actions sont loguées
- ✅ **Actions non définies tolérées** : Les actions sans handler sont simplement loguées
- ✅ **Validation optionnelle** : Validez les arguments avant exécution
- ✅ **Extensible** : Architecture ouverte pour ajouter facilement de nouvelles actions

## 🚀 Démarrage rapide

### Utiliser l'action print

```go
package main

import (
    "github.com/treivax/tsd/rete"
)

func main() {
    // Créer un réseau RETE
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Définir un type
    personType := rete.TypeDefinition{
        Type: "typeDefinition",
        Name: "Person",
        Fields: []rete.Field{
            {Name: "id", Type: "string"},
            {Name: "name", Type: "string"},
        },
    }
    network.Types = append(network.Types, personType)
    
    // Créer un fait
    fact := &rete.Fact{
        ID:   "person_1",
        Type: "Person",
        Fields: map[string]interface{}{
            "id":   "1",
            "name": "Alice",
        },
    }
    
    // Créer un token avec bindings
    token := &rete.Token{
        ID:    "token1",
        Facts: []*rete.Fact{fact},
        Bindings: map[string]*rete.Fact{
            "p": fact,
        },
    }
    
    // Créer une action print
    action := &rete.Action{
        Type: "action",
        Jobs: []rete.JobCall{
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
            },
        },
    }
    
    // Exécuter l'action
    network.ActionExecutor.ExecuteAction(action, token)
    // Sortie: Alice
}
```

### Créer une action personnalisée

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

// Constante pour le nom de l'action
const ActionNameNotify = "notify"

// Définir l'action
type NotifyAction struct {
    channel string
}

func NewNotifyAction(channel string) *NotifyAction {
    return &NotifyAction{channel: channel}
}

func (na *NotifyAction) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
    if len(args) == 0 {
        return fmt.Errorf("notify requires at least one argument")
    }
    
    message := fmt.Sprintf("%v", args[0])
    fmt.Printf("[%s] %s\n", na.channel, message)
    
    return nil
}

func (na *NotifyAction) GetName() string {
    return ActionNameNotify
}

func (na *NotifyAction) Validate(args []interface{}) error {
    if len(args) == 0 {
        return fmt.Errorf("notify requires at least one argument")
    }
    return nil
}

func main() {
    // Créer le réseau
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Enregistrer l'action personnalisée
    notifyAction := NewNotifyAction("slack")
    err := network.ActionExecutor.RegisterAction(notifyAction)
    if err != nil {
        panic(err)
    }
    
    // Utiliser l'action dans une règle
    // ... (voir exemples complets)
}
```

## 📚 Documentation

### Fichiers de documentation

- **[ACTIONS_SYSTEM.md](ACTIONS_SYSTEM.md)** : Documentation complète du système
- **[examples/action_print_example.go](examples/action_print_example.go)** : Exemples d'utilisation

### Fichiers sources

- **[action_handler.go](action_handler.go)** : Interface ActionHandler et ActionRegistry
- **[action_print.go](action_print.go)** : Implémentation de l'action print
- **[action_executor.go](action_executor.go)** : ActionExecutor intégré au moteur RETE

### Tests

- **[action_handler_test.go](action_handler_test.go)** : Tests unitaires
- **[action_print_integration_test.go](action_print_integration_test.go)** : Tests d'intégration

## 🎨 Action Print

L'action `print` est incluse par défaut et permet d'afficher des valeurs sur la sortie standard.

### Types supportés

| Type | Exemple | Sortie |
|------|---------|--------|
| String | `"Hello"` | `Hello` |
| Number | `42.5` | `42.5` |
| Boolean | `true` | `true` |
| Field Access | `p.name` | Valeur du champ |
| Variable | `p` | Représentation complète du fait |
| Fact | `*Fact{...}` | `Person{id: person_1, name: "Alice"}` |

### Exemples

```go
// Afficher une chaîne
print("Hello, World!")

// Afficher un champ
print(p.name)

// Afficher un nombre
print(p.age)

// Afficher un fait complet
print(p)
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                   ReteNetwork                       │
│                                                     │
│  ┌───────────────────────────────────────────────┐ │
│  │           ActionExecutor                      │ │
│  │                                               │ │
│  │  ┌─────────────────────────────────────────┐ │ │
│  │  │        ActionRegistry                   │ │ │
│  │  │                                         │ │ │
│  │  │  • print  → PrintAction                │ │ │
│  │  │  • notify → NotifyAction               │ │ │
│  │  │  • custom → CustomAction               │ │ │
│  │  │  • ...                                 │ │ │
│  │  └─────────────────────────────────────────┘ │ │
│  │                                               │ │
│  │  [Execute] → [Get Handler] → [Validate]      │ │
│  │                             → [Execute]       │ │
│  │                             → [Log]           │ │
│  └───────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

## 🔧 API Principale

### ActionHandler (Interface)

```go
type ActionHandler interface {
    Execute(args []interface{}, ctx *ExecutionContext) error
    GetName() string
    Validate(args []interface{}) error
}
```

### ActionRegistry

```go
// Enregistrer une action
registry.Register(handler)

// Supprimer une action
registry.Unregister(actionName)

// Récupérer une action
handler := registry.Get(actionName)

// Vérifier l'existence
if registry.Has(actionName) { ... }

// Lister toutes les actions
names := registry.GetRegisteredNames()

// Nettoyer le registry
registry.Clear()
```

### ActionExecutor

```go
// Enregistrer une action
network.ActionExecutor.RegisterAction(handler)

// Accéder au registry
registry := network.ActionExecutor.GetRegistry()

// Activer/désactiver le logging
network.ActionExecutor.SetLogging(false)
```

## 📊 Logging

Toutes les actions sont automatiquement loguées :

### Action définie et exécutée

```
📋 ACTION: print(p.name)
🎯 ACTION EXÉCUTÉE: print("Alice")
```

### Action non définie

```
📋 ACTION: send_email(p.email)
📋 ACTION NON DÉFINIE (log uniquement): send_email("alice@example.com")
```

## ✅ Tests

Exécuter les tests :

```bash
# Tous les tests d'actions
go test -v -run TestAction ./rete

# Tests du registry
go test -v -run TestActionRegistry ./rete

# Tests de l'action print
go test -v -run TestPrintAction ./rete

# Tests d'intégration
go test -v -run TestPrintActionIntegration ./rete
```

## 🎯 Exemples complets

### Exemple 1 : Action print simple

```bash
go run rete/examples/action_print_example.go
```

### Exemple 2 : Action personnalisée avec configuration

```go
type EmailAction struct {
    smtpHost string
    smtpPort int
    sender   string
}

func NewEmailAction(smtpHost string, smtpPort int, sender string) *EmailAction {
    return &EmailAction{
        smtpHost: smtpHost,
        smtpPort: smtpPort,
        sender:   sender,
    }
}

func (ea *EmailAction) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
    if len(args) < 2 {
        return fmt.Errorf("email action requires 2 arguments: recipient and subject")
    }
    
    recipient := args[0].(string)
    subject := args[1].(string)
    
    // Envoyer l'email (implémentation simplifiée)
    fmt.Printf("📧 Sending email to %s: %s\n", recipient, subject)
    
    return nil
}

func (ea *EmailAction) GetName() string {
    return "email"
}

func (ea *EmailAction) Validate(args []interface{}) error {
    if len(args) < 2 {
        return fmt.Errorf("email action requires 2 arguments")
    }
    return nil
}
```

## 🔒 Bonnes pratiques

### 1. Utiliser des constantes pour les noms

```go
const (
    ActionNamePrint  = "print"
    ActionNameNotify = "notify"
    ActionNameEmail  = "email"
)
```

### 2. Valider les arguments

Toujours implémenter `Validate` pour une détection précoce des erreurs.

### 3. Pas de hardcoding

```go
// ❌ Mauvais
func (h *MyAction) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
    url := "http://localhost:8080" // Hardcodé !
    // ...
}

// ✅ Bon
type MyAction struct {
    baseURL string
}

func NewMyAction(baseURL string) *MyAction {
    return &MyAction{baseURL: baseURL}
}
```

### 4. Thread-safety

Les handlers doivent être thread-safe s'ils maintiennent un état.

### 5. Gestion des erreurs

Retourner des erreurs claires et descriptives.

### 6. Documentation

Documenter chaque action avec des exemples.

## 🚧 Feuille de route

Actions futures à implémenter :

- [ ] `assert(fact)` : Assertion de nouveau fait
- [ ] `retract(fact)` : Retrait de fait
- [ ] `modify(fact, field, value)` : Modification de fait
- [ ] `log(level, message)` : Logging avec niveaux
- [ ] `http(method, url, body)` : Appel HTTP
- [ ] `emit(event, data)` : Émission d'événement
- [ ] `delay(duration, action)` : Action différée
- [ ] `aggregate(collection, operation)` : Opérations d'agrégation

## 🤝 Contribution

Pour ajouter une nouvelle action :

1. Créer le fichier `action_<nom>.go`
2. Implémenter l'interface `ActionHandler`
3. Ajouter les tests dans `action_<nom>_test.go`
4. Mettre à jour cette documentation
5. Ajouter un exemple dans `examples/`

## 📝 Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License

## 🔗 Ressources

- [Documentation complète](ACTIONS_SYSTEM.md)
- [Guide d'architecture RETE](README.md)
- [Exemples](examples/)
- [Tests](action_handler_test.go)