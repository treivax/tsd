# Intégration RETE → Xuples

## 📊 Vue d'ensemble

Ce document décrit l'architecture d'intégration entre le moteur de règles RETE et le module xuples, permettant aux règles de créer des xuples via l'action `Xuple`.

## 🏗️ Architecture

### Diagramme de composants

```
┌──────────────────────────────────────────────────────────────────┐
│                        Serveur TSD                                │
├──────────────────────────────────────────────────────────────────┤
│                                                                    │
│  ┌─────────────┐                        ┌──────────────┐         │
│  │   Parser    │──(AST)───────────────▶│   Compiler   │         │
│  └─────────────┘                        └──────┬───────┘         │
│                                                 │                 │
│                                    ┌────────────▼──────────┐      │
│                                    │  XupleSpaceDecl       │      │
│                                    └────────────┬──────────┘      │
│                                                 │                 │
│                                    ┌────────────▼──────────┐      │
│                                    │  XupleManager         │      │
│                                    │  ├─ notifications     │      │
│                                    │  ├─ events            │      │
│                                    │  └─ ...               │      │
│                                    └────────────┬──────────┘      │
│                                                 │                 │
│  ┌────────────┐                    ┌────────────▼──────────┐      │
│  │ RETE       │◀──(configure)─────▶│ BuiltinActionExecutor │      │
│  │ Network    │                    │  ├─ Print             │      │
│  │            │                    │  ├─ Log               │      │
│  │  ┌──────┐  │                    │  ├─ Update            │      │
│  │  │ Rule │──┼──(activation)────▶ │  ├─ Insert            │      │
│  │  └──────┘  │                    │  ├─ Retract           │      │
│  │            │                    │  └─ Xuple ────────────┼──┐   │
│  └────────────┘                    └───────────────────────┘  │   │
│                                                                │   │
│                                                                │   │
│                                    ┌───────────────────────────▼┐  │
│                                    │  XupleManager.CreateXuple   │  │
│                                    │    ├─ Fact principal        │  │
│                                    │    └─ TriggeringFacts       │  │
│                                    └─────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Flux de données

```
1. Parsing du programme TSD
   ├─> Extraction des xuple-space declarations
   └─> Extraction des règles avec actions Xuple

2. Compilation
   ├─> Création du XupleManager
   ├─> Instanciation des xuple-spaces déclarés
   │   ├─> Conversion AST → Config xuples
   │   ├─> Construction des policies
   │   └─> Création du space dans le manager
   └─> Configuration du réseau RETE

3. Exécution des règles
   ├─> Activation d'une règle
   ├─> Exécution de l'action Xuple
   │   ├─> Extraction des faits déclencheurs du Token
   │   ├─> Appel XupleManager.CreateXuple()
   │   └─> Insertion dans le xuple-space
   └─> Xuple disponible pour récupération
```

## 🔌 Points d'intégration

### 1. Parser → Compiler

**Fichiers** : `constraint/constraint_types.go`, `constraint/parser.go`

Le parser extrait les déclarations de xuple-spaces :

```go
type XupleSpaceDeclaration struct {
    Type              string                     // "xupleSpaceDeclaration"
    Name              string                     // Nom du xuple-space
    SelectionPolicy   string                     // "random", "fifo", "lifo"
    ConsumptionPolicy XupleConsumptionPolicyConf // Configuration consumption
    RetentionPolicy   XupleRetentionPolicyConf   // Configuration retention
}
```

### 2. Compiler → XupleManager

**Fichier** : `internal/servercmd/servercmd.go`

**Fonction** : `instantiateXupleSpaces()`

Cette fonction convertit les déclarations AST en configurations concrètes :

```go
func instantiateXupleSpaces(xupleManager xuples.XupleManager, 
                            declarations []constraint.XupleSpaceDeclaration) error {
    for _, decl := range declarations {
        // Construire la configuration
        config, err := buildXupleSpaceConfig(decl)
        if err != nil {
            return err
        }
        
        // Créer le xuple-space
        err = xupleManager.CreateXupleSpace(decl.Name, config)
        if err != nil {
            return err
        }
    }
    return nil
}
```

### 3. RETE → BuiltinActionExecutor

**Fichier** : `rete/actions/builtin.go`

L'action Xuple est implémentée dans le BuiltinActionExecutor :

```go
func (e *BuiltinActionExecutor) executeXuple(args []interface{}, token *rete.Token) error {
    // Valider les arguments
    xuplespace := args[0].(string)
    fact := args[1].(*rete.Fact)
    
    // Extraire les faits déclencheurs
    triggeringFacts := e.extractTriggeringFacts(token)
    
    // Déléguer au XupleManager
    return e.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
}
```

### 4. BuiltinActionExecutor → XupleManager

**Interface** : `xuples.XupleManager`

```go
type XupleManager interface {
    CreateXupleSpace(name string, config XupleSpaceConfig) error
    GetXupleSpace(name string) (XupleSpace, error)
    CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error
    ListXupleSpaces() []string
    Close() error
}
```

## 📦 Conversion des politiques

### SelectionPolicy

```go
func buildSelectionPolicy(policyName string) (xuples.SelectionPolicy, error) {
    switch policyName {
    case "random":
        return xuples.NewRandomSelectionPolicy(), nil
    case "fifo":
        return xuples.NewFIFOSelectionPolicy(), nil
    case "lifo":
        return xuples.NewLIFOSelectionPolicy(), nil
    default:
        return nil, fmt.Errorf("unknown selection policy: %s", policyName)
    }
}
```

### ConsumptionPolicy

```go
func buildConsumptionPolicy(conf constraint.XupleConsumptionPolicyConf) (xuples.ConsumptionPolicy, error) {
    switch conf.Type {
    case "once":
        return xuples.NewOnceConsumptionPolicy(), nil
    case "per-agent":
        return xuples.NewPerAgentConsumptionPolicy(), nil
    case "limited":
        return xuples.NewLimitedConsumptionPolicy(conf.Limit), nil
    default:
        return nil, fmt.Errorf("unknown consumption policy: %s", conf.Type)
    }
}
```

### RetentionPolicy

```go
func buildRetentionPolicy(conf constraint.XupleRetentionPolicyConf) (xuples.RetentionPolicy, error) {
    switch conf.Type {
    case "unlimited":
        return xuples.NewUnlimitedRetentionPolicy(), nil
    case "duration":
        duration := time.Duration(conf.Duration) * time.Second
        return xuples.NewDurationRetentionPolicy(duration), nil
    default:
        return nil, fmt.Errorf("unknown retention policy: %s", conf.Type)
    }
}
```

## 🔄 Extraction des faits déclencheurs

La méthode `extractTriggeringFacts` parcourt la chaîne de tokens RETE pour extraire tous les faits qui ont contribué à l'activation :

```go
func (e *BuiltinActionExecutor) extractTriggeringFacts(token *rete.Token) []*rete.Fact {
    if token == nil {
        return []*rete.Fact{}
    }
    
    var facts []*rete.Fact
    
    // Parcourir la chaîne de tokens via Parent
    for t := token; t != nil; t = t.Parent {
        if t.Facts != nil {
            facts = append(facts, t.Facts...)
        }
    }
    
    // Inverser pour avoir l'ordre chronologique
    for i := 0; i < len(facts)/2; i++ {
        facts[i], facts[len(facts)-1-i] = facts[len(facts)-1-i], facts[i]
    }
    
    return facts
}
```

## ⚠️ Gestion d'erreurs

### Erreurs de compilation

- **Xuple-space non déclaré** : Détecté au compile-time si le nom du xuple-space est une constante
- **Politique invalide** : Détecté lors de la conversion AST → Config

### Erreurs d'exécution

- **XupleManager non configuré** : `fmt.Errorf("Xuple action requires XupleManager to be configured")`
- **Xuple-space inexistant** : `ErrXupleSpaceNotFound` retourné par `GetXupleSpace()`
- **Arguments invalides** : Validation du nombre et des types d'arguments

## 🔐 Thread-Safety

- **XupleManager** : Thread-safe via `sync.RWMutex`
- **BuiltinActionExecutor** : Thread-safe si le réseau RETE l'est
- **XupleSpace** : Thread-safe via `sync.RWMutex`

## 🧪 Tests

Voir `tests/integration/xuples_integration_test.go` pour :
- Test nominal d'intégration complète
- Test d'erreur (xuple-space non déclaré)
- Test avec plusieurs règles
- Test avec plusieurs faits déclencheurs
- Test avec différentes politiques

## 📚 Documentation associée

- [Design xuples](../design/xuples-architecture.md)
- [User guide](../user-guide/using-xuples.md)
- [API Reference](../api/xuples-api.md)
