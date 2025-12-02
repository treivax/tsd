# Guide de Démarrage Rapide - Système d'Actions TSD

## ⚡ En 2 minutes

### 1. Utiliser l'action print (déjà intégrée)

```go
package main

import "github.com/treivax/tsd/rete"

func main() {
    // Créer le réseau
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Créer un fait
    fact := &rete.Fact{
        ID:   "p1",
        Type: "Person",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  25.0,
        },
    }
    
    // Créer un token
    token := &rete.Token{
        ID:    "t1",
        Facts: []*rete.Fact{fact},
        Bindings: map[string]*rete.Fact{"p": fact},
    }
    
    // Créer une action print
    action := &rete.Action{
        Type: "action",
        Jobs: []rete.JobCall{{
            Type: "jobCall",
            Name: "print",
            Args: []interface{}{
                map[string]interface{}{
                    "type":   "fieldAccess",
                    "object": "p",
                    "field":  "name",
                },
            },
        }},
    }
    
    // Exécuter
    network.ActionExecutor.ExecuteAction(action, token)
    // Sortie: Alice
}
```

### 2. Créer votre première action personnalisée

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

// Définir votre action
type NotifyAction struct {
    channel string
}

func (na *NotifyAction) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
    message := fmt.Sprintf("%v", args[0])
    fmt.Printf("[%s] %s\n", na.channel, message)
    return nil
}

func (na *NotifyAction) GetName() string {
    return "notify"
}

func (na *NotifyAction) Validate(args []interface{}) error {
    if len(args) == 0 {
        return fmt.Errorf("notify requires an argument")
    }
    return nil
}

func main() {
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Enregistrer votre action
    notify := &NotifyAction{channel: "slack"}
    network.ActionExecutor.RegisterAction(notify)
    
    // Utiliser dans une règle
    // ... (voir exemples complets)
}
```

## 🎯 Cas d'usage courants

### Afficher une chaîne littérale

```go
action := &rete.Action{
    Jobs: []rete.JobCall{{
        Name: "print",
        Args: []interface{}{
            map[string]interface{}{
                "type":  "string",
                "value": "Hello, World!",
            },
        },
    }},
}
```

### Afficher un champ d'un fait

```go
action := &rete.Action{
    Jobs: []rete.JobCall{{
        Name: "print",
        Args: []interface{}{
            map[string]interface{}{
                "type":   "fieldAccess",
                "object": "p",
                "field":  "name",
            },
        },
    }},
}
```

### Afficher un nombre

```go
action := &rete.Action{
    Jobs: []rete.JobCall{{
        Name: "print",
        Args: []interface{}{
            map[string]interface{}{
                "type":   "fieldAccess",
                "object": "p",
                "field":  "age",
            },
        },
    }},
}
```

### Afficher un fait complet

```go
action := &rete.Action{
    Jobs: []rete.JobCall{{
        Name: "print",
        Args: []interface{}{
            map[string]interface{}{
                "type": "variable",
                "name": "p",
            },
        },
    }},
}
```

### Exécuter plusieurs actions

```go
action := &rete.Action{
    Jobs: []rete.JobCall{
        {
            Name: "print",
            Args: []interface{}{
                map[string]interface{}{
                    "type":  "string",
                    "value": "User detected:",
                },
            },
        },
        {
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
```

## 🔧 API Essentielle

### ActionExecutor

```go
// Enregistrer une action
network.ActionExecutor.RegisterAction(handler)

// Accéder au registry
registry := network.ActionExecutor.GetRegistry()

// Désactiver le logging
network.ActionExecutor.SetLogging(false)
```

### ActionRegistry

```go
// Vérifier si une action existe
if registry.Has("print") { ... }

// Lister toutes les actions
names := registry.GetRegisteredNames()

// Supprimer une action
registry.Unregister("custom_action")
```

## 📊 Exemples de sortie

### Action définie

```
📋 ACTION: print(p.name)
Alice
🎯 ACTION EXÉCUTÉE: print("Alice")
```

### Action non définie

```
📋 ACTION: send_email(p.email)
📋 ACTION NON DÉFINIE (log uniquement): send_email("alice@example.com")
```

## 🚀 Template pour action personnalisée

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

// 1. Définir la constante
const ActionNameMyAction = "my_action"

// 2. Définir la structure
type MyAction struct {
    // Vos paramètres de configuration
    config Config
}

// 3. Constructeur
func NewMyAction(config Config) *MyAction {
    return &MyAction{config: config}
}

// 4. Implémenter Execute
func (ma *MyAction) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
    // Vérifier les arguments
    if len(args) == 0 {
        return fmt.Errorf("my_action requires at least one argument")
    }
    
    // Récupérer les arguments
    value := args[0]
    
    // Votre logique métier
    fmt.Printf("Executing my_action with: %v\n", value)
    
    return nil
}

// 5. Implémenter GetName
func (ma *MyAction) GetName() string {
    return ActionNameMyAction
}

// 6. Implémenter Validate
func (ma *MyAction) Validate(args []interface{}) error {
    if len(args) == 0 {
        return fmt.Errorf("my_action requires at least one argument")
    }
    
    // Validations spécifiques
    // ...
    
    return nil
}

// 7. Enregistrer et utiliser
func main() {
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // Enregistrer
    action := NewMyAction(config)
    err := network.ActionExecutor.RegisterAction(action)
    if err != nil {
        panic(err)
    }
    
    // Utiliser dans vos règles
    // ...
}
```

## 📚 Ressources

- **Guide complet** : [ACTIONS_README.md](ACTIONS_README.md)
- **Documentation technique** : [ACTIONS_SYSTEM.md](ACTIONS_SYSTEM.md)
- **Exemple d'utilisation** : [examples/action_print_example.go](examples/action_print_example.go)
- **Tests** : [action_handler_test.go](action_handler_test.go)

## ⚠️ Pièges courants

### ❌ Oublier de créer le token avec bindings

```go
// Mauvais
token := &rete.Token{Facts: []*rete.Fact{fact}}

// Bon
token := &rete.Token{
    Facts: []*rete.Fact{fact},
    Bindings: map[string]*rete.Fact{"p": fact},
}
```

### ❌ Utiliser un nom d'action hardcodé

```go
// Mauvais
func (ma *MyAction) GetName() string {
    return "my_action"  // Hardcodé !
}

// Bon
const ActionNameMyAction = "my_action"

func (ma *MyAction) GetName() string {
    return ActionNameMyAction
}
```

### ❌ Ne pas valider les arguments

```go
// Mauvais
func (ma *MyAction) Validate(args []interface{}) error {
    return nil  // Pas de validation !
}

// Bon
func (ma *MyAction) Validate(args []interface{}) error {
    if len(args) < 2 {
        return fmt.Errorf("my_action requires 2 arguments")
    }
    return nil
}
```

## 🎉 Prêt à commencer !

Vous avez maintenant tout ce qu'il faut pour :
1. ✅ Utiliser l'action print
2. ✅ Créer vos propres actions
3. ✅ Gérer les actions dans vos règles

Pour des exemples plus avancés, consultez la [documentation complète](ACTIONS_SYSTEM.md).

## 💡 Besoin d'aide ?

- Voir les tests : `action_handler_test.go`
- Exécuter l'exemple : `go run rete/examples/action_print_example.go`
- Lire la doc : `ACTIONS_SYSTEM.md`

Happy coding! 🚀