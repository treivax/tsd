# Conception du système d'actions par défaut

## 🎯 Objectif

Implémenter un système d'actions par défaut :
- **Non hardcodé** : définitions dans un fichier `.tsd` embarqué
- **Automatique** : chargement à l'initialisation
- **Protégé** : interdiction de redéfinition
- **Complet** : toutes les actions système (Print, Log, Update, Insert, Retract, Xuple)

## 📋 Actions par défaut requises

| Action | Signature | Description |
|--------|-----------|-------------|
| **Print** | `Print(message: string)` | Affiche sur stdout |
| **Log** | `Log(message: string)` | Génère une trace système |
| **Update** | `Update(fact: any)` | Met à jour un fait existant |
| **Insert** | `Insert(fact: any)` | Insère un nouveau fait |
| **Retract** | `Retract(id: string)` | Supprime un fait par ID |
| **Xuple** | `Xuple(xuplespace: string, fact: any)` | Crée un xuple |

## 🏗️ Architecture proposée

### 1. Fichier de définitions embarqué

**Fichier** : `internal/defaultactions/defaults.tsd`

- Contient les 6 actions par défaut
- Embarqué dans le binaire via `//go:embed`
- Parsé au démarrage via le parser TSD existant
- Format standard TSD (pas de syntaxe spéciale)

### 2. Package defaultactions

**Package** : `internal/defaultactions`

Responsabilités :
- Charger et parser `defaults.tsd`
- Marquer les actions comme "par défaut"
- Fournir la liste des noms d'actions système
- Exposer une fonction de chargement

### 3. Marquage des actions système

**Modification** : `constraint/constraint_types.go`

Ajouter un champ `IsDefault bool` à `ActionDefinition` :

```go
type ActionDefinition struct {
    Type       string      `json:"type"`
    Name       string      `json:"name"`
    Parameters []Parameter `json:"parameters"`
    IsDefault  bool        `json:"isDefault,omitempty"` // ← NOUVEAU
}
```

### 4. Validation des redéfinitions

**Modification** : `constraint/action_validator.go`

- Vérifier si une action est marquée `IsDefault`
- Retourner une erreur explicite en cas de tentative de redéfinition
- Message : `"cannot redefine default action 'X'"`

### 5. Implémentations natives

**Nouveau package** : `rete/actions`

Implémentations concrètes des 6 actions :
- `BuiltinActionExecutor` : dispatche vers la bonne implémentation
- Une méthode par action (executeprint, executeLog, etc.)
- Délégation au réseau RETE pour Insert/Update/Retract
- Délégation au XupleManager pour Xuple

## 📁 Structure des fichiers

```
tsd/
├── docs/xuples/implementation/
│   ├── 03-current-action-system.md     ← Analyse (créé)
│   └── 04-default-actions-design.md    ← Ce document
│
├── internal/defaultactions/
│   ├── defaults.tsd                    ← Définitions TSD
│   ├── loader.go                       ← Chargement
│   └── loader_test.go                  ← Tests
│
├── rete/actions/
│   ├── builtin.go                      ← Implémentations
│   └── builtin_test.go                 ← Tests
│
├── constraint/
│   ├── constraint_types.go             ← Modifié (IsDefault)
│   └── action_validator.go             ← Modifié (validation)
│
└── rete/
    └── action_executor.go              ← Modifié (chargement)
```

## 🔄 Flux de chargement

### Diagramme de séquence

```
┌──────────┐         ┌──────────────┐         ┌─────────────┐         ┌──────────────┐
│   Main   │         │ActionExecutor│         │defaultactions│         │    Parser    │
└────┬─────┘         └──────┬───────┘         └──────┬──────┘         └──────┬───────┘
     │                      │                        │                        │
     │ NewActionExecutor()  │                        │                        │
     │─────────────────────>│                        │                        │
     │                      │                        │                        │
     │                      │ LoadDefaultActions()   │                        │
     │                      │───────────────────────>│                        │
     │                      │                        │                        │
     │                      │                        │ ParseTSD(defaults.tsd) │
     │                      │                        │───────────────────────>│
     │                      │                        │                        │
     │                      │                        │  Program (6 actions)   │
     │                      │                        │<───────────────────────│
     │                      │                        │                        │
     │                      │                        │ Marquer IsDefault=true │
     │                      │                        │──────────┐             │
     │                      │                        │          │             │
     │                      │                        │<─────────┘             │
     │                      │                        │                        │
     │                      │   []ActionDefinition   │                        │
     │                      │<───────────────────────│                        │
     │                      │                        │                        │
     │ Enregistrer dans     │                        │                        │
     │ ActionValidator      │                        │                        │
     │──────────────────────┤                        │                        │
     │                      │                        │                        │
```

### Étapes détaillées

1. **Initialisation** (`main` ou test setup)
   - Création de `ActionExecutor`

2. **Chargement automatique** (`NewActionExecutor`)
   - Appel `defaultactions.LoadDefaultActions()`
   - Parser le fichier embarqué `defaults.tsd`
   - Marquer chaque action `IsDefault = true`
   - Retourner `[]ActionDefinition`

3. **Enregistrement**
   - Ajouter les actions au `ActionValidator`
   - Les actions sont maintenant connues du système

4. **Validation utilisateur**
   - Parse du fichier utilisateur `.tsd`
   - Si action système déjà définie → **ERREUR**
   - Sinon → enregistrement normal

## 📄 Contenu de defaults.tsd

```tsd
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// ============================================================================
// ACTIONS PAR DÉFAUT DU SYSTÈME TSD
// ============================================================================
//
// Ces actions sont automatiquement disponibles dans tous les programmes TSD.
// Elles ne nécessitent pas de déclaration explicite.
//
// Toute tentative de redéfinition provoquera une erreur de compilation.
// ============================================================================

// Print affiche une chaîne de caractères sur la sortie standard
// Paramètres:
//   - message: la chaîne à afficher
action Print(message: string)

// Log génère une trace dans le système de logging
// Paramètres:
//   - message: la chaîne à tracer
action Log(message: string)

// Update modifie un fait existant et met à jour les tokens liés dans RETE
// Paramètres:
//   - fact: le fait à modifier (doit exister dans le réseau)
// Notes:
//   - Déclenche la propagation des mises à jour dans le réseau RETE
//   - Le fait doit avoir le même type qu'un fait existant
action Update(fact: any)

// Insert crée un nouveau fait et l'insère dans le réseau RETE
// Paramètres:
//   - fact: le nouveau fait à créer
// Notes:
//   - Le fait est propagé dans le réseau RETE
//   - Peut déclencher l'activation de nouvelles règles
action Insert(fact: any)

// Retract supprime un fait du réseau RETE ainsi que tous les tokens liés
// Paramètres:
//   - id: l'identifiant du fait à supprimer
// Notes:
//   - Tous les tokens dépendant de ce fait sont invalidés
//   - La suppression se propage dans tout le réseau
action Retract(id: string)

// Xuple crée un xuple dans le xuple-space spécifié
// Paramètres:
//   - xuplespace: nom du xuple-space cible
//   - fact: le fait principal du xuple
// Notes:
//   - Les faits déclencheurs sont automatiquement extraits du token
//   - Le xuple-space doit avoir été déclaré via 'xuple-space'
//   - Le xuple est soumis aux politiques du xuple-space
action Xuple(xuplespace: string, fact: any)
```

## 🔧 Code loader.go

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package defaultactions

import (
    _ "embed"
    "fmt"
    
    "tsd/constraint"
)

// defaults.tsd est embarqué dans le binaire via go:embed
//go:embed defaults.tsd
var defaultActionsTSD string

// DefaultActionNames contient les noms de toutes les actions par défaut
var DefaultActionNames = []string{
    "Print",
    "Log",
    "Update",
    "Insert",
    "Retract",
    "Xuple",
}

// LoadDefaultActions parse le fichier defaults.tsd et retourne les actions
func LoadDefaultActions() ([]constraint.ActionDefinition, error) {
    // Parser le fichier embarqué
    result, err := constraint.ParseConstraintProgram(defaultActionsTSD)
    if err != nil {
        return nil, fmt.Errorf("failed to parse default actions: %w", err)
    }
    
    // Vérifier que toutes les actions attendues sont présentes
    if len(result.Actions) != len(DefaultActionNames) {
        return nil, fmt.Errorf("expected %d default actions, got %d",
            len(DefaultActionNames), len(result.Actions))
    }
    
    // Marquer chaque action comme "par défaut"
    for i := range result.Actions {
        result.Actions[i].IsDefault = true
    }
    
    return result.Actions, nil
}

// IsDefaultAction vérifie si un nom correspond à une action par défaut
func IsDefaultAction(name string) bool {
    for _, defaultName := range DefaultActionNames {
        if name == defaultName {
            return true
        }
    }
    return false
}
```

## 🔧 Modification ActionValidator

### Nouvelle méthode : ValidateNonRedefinition

```go
// ValidateNonRedefinition vérifie qu'aucune action par défaut n'est redéfinie
func (av *ActionValidator) ValidateNonRedefinition(newActions []ActionDefinition) error {
    for _, newAction := range newActions {
        if existing, exists := av.actions[newAction.Name]; exists {
            if existing.IsDefault {
                return fmt.Errorf(
                    "cannot redefine default action '%s' (default actions: %v)",
                    newAction.Name,
                    defaultactions.DefaultActionNames,
                )
            }
            // Doublon non-système
            return fmt.Errorf("action '%s' already defined", newAction.Name)
        }
    }
    return nil
}
```

### Intégration

Dans le flux de parsing/validation :

```go
// 1. Charger actions par défaut
defaultActs, err := defaultactions.LoadDefaultActions()
if err != nil {
    return nil, err
}

// 2. Créer le validator avec les actions par défaut
validator := NewActionValidator(defaultActs, types)

// 3. Parser le fichier utilisateur
userProgram, err := ParseConstraintProgram(userInput)
if err != nil {
    return nil, err
}

// 4. Valider qu'aucune action système n'est redéfinie
if err := validator.ValidateNonRedefinition(userProgram.Actions); err != nil {
    return nil, err
}

// 5. Ajouter les actions utilisateur
for _, action := range userProgram.Actions {
    validator.AddAction(action)
}
```

## 🔧 Implémentation BuiltinActionExecutor

### Structure

```go
// BuiltinActionExecutor exécute les actions par défaut du système
type BuiltinActionExecutor struct {
    network      *rete.ReteNetwork
    xupleManager XupleManager
    output       io.Writer
    logger       *log.Logger
}

// XupleManager interface vers le module xuples
type XupleManager interface {
    CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error
}
```

### Méthode Execute

```go
func (e *BuiltinActionExecutor) Execute(actionName string, args []interface{}, ctx *rete.ExecutionContext) error {
    switch actionName {
    case "Print":
        return e.executePrint(args)
    case "Log":
        return e.executeLog(args)
    case "Update":
        return e.executeUpdate(args, ctx)
    case "Insert":
        return e.executeInsert(args, ctx)
    case "Retract":
        return e.executeRetract(args, ctx)
    case "Xuple":
        return e.executeXuple(args, ctx)
    default:
        return fmt.Errorf("unknown builtin action: %s", actionName)
    }
}
```

### Implémentations individuelles

```go
func (e *BuiltinActionExecutor) executePrint(args []interface{}) error {
    if len(args) != 1 {
        return fmt.Errorf("Print expects 1 argument, got %d", len(args))
    }
    message, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Print expects string argument")
    }
    fmt.Fprintln(e.output, message)
    return nil
}

func (e *BuiltinActionExecutor) executeLog(args []interface{}) error {
    if len(args) != 1 {
        return fmt.Errorf("Log expects 1 argument, got %d", len(args))
    }
    message, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Log expects string argument")
    }
    e.logger.Printf("[TSD] %s", message)
    return nil
}

func (e *BuiltinActionExecutor) executeUpdate(args []interface{}, ctx *rete.ExecutionContext) error {
    if len(args) != 1 {
        return fmt.Errorf("Update expects 1 argument, got %d", len(args))
    }
    fact, ok := args[0].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Update expects fact argument")
    }
    return e.network.UpdateFact(fact)
}

func (e *BuiltinActionExecutor) executeInsert(args []interface{}, ctx *rete.ExecutionContext) error {
    if len(args) != 1 {
        return fmt.Errorf("Insert expects 1 argument, got %d", len(args))
    }
    fact, ok := args[0].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Insert expects fact argument")
    }
    return e.network.InsertFact(fact)
}

func (e *BuiltinActionExecutor) executeRetract(args []interface{}, ctx *rete.ExecutionContext) error {
    if len(args) != 1 {
        return fmt.Errorf("Retract expects 1 argument, got %d", len(args))
    }
    id, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Retract expects string argument")
    }
    return e.network.RetractFact(id)
}

func (e *BuiltinActionExecutor) executeXuple(args []interface{}, ctx *rete.ExecutionContext) error {
    if len(args) != 2 {
        return fmt.Errorf("Xuple expects 2 arguments, got %d", len(args))
    }
    
    xuplespace, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Xuple expects string as first argument")
    }
    
    fact, ok := args[1].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Xuple expects fact as second argument")
    }
    
    // Extraire les faits déclencheurs du contexte
    triggeringFacts := extractTriggeringFacts(ctx.Token)
    
    return e.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
}
```

## 🧪 Tests

### Tests du loader

```go
func TestLoadDefaultActions(t *testing.T) {
    actions, err := LoadDefaultActions()
    if err != nil {
        t.Fatalf("Failed to load default actions: %v", err)
    }
    
    // Vérifier le nombre
    if len(actions) != 6 {
        t.Errorf("Expected 6 actions, got %d", len(actions))
    }
    
    // Vérifier qu'elles sont marquées IsDefault
    for _, action := range actions {
        if !action.IsDefault {
            t.Errorf("Action %s should be marked IsDefault", action.Name)
        }
    }
    
    // Vérifier les noms
    expectedNames := map[string]bool{
        "Print": true, "Log": true, "Update": true,
        "Insert": true, "Retract": true, "Xuple": true,
    }
    
    for _, action := range actions {
        if !expectedNames[action.Name] {
            t.Errorf("Unexpected action: %s", action.Name)
        }
    }
}
```

### Tests de redéfinition

```go
func TestCannotRedefineDefaultAction(t *testing.T) {
    input := `
        action Print(msg: string)
    `
    
    _, err := CompileProgram(input)
    if err == nil {
        t.Fatal("Expected error when redefining default action")
    }
    
    if !strings.Contains(err.Error(), "cannot redefine default action") {
        t.Errorf("Wrong error message: %v", err)
    }
}
```

### Tests d'exécution

```go
func TestBuiltinActions(t *testing.T) {
    tests := []struct {
        name   string
        action string
        args   []interface{}
        wantErr bool
    }{
        {"Print OK", "Print", []interface{}{"Hello"}, false},
        {"Print missing arg", "Print", []interface{}{}, true},
        {"Log OK", "Log", []interface{}{"Info"}, false},
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            executor := NewBuiltinActionExecutor(network, nil, nil, logger)
            err := executor.Execute(tt.action, tt.args, ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 📊 Métriques attendues

- **Couverture tests** : > 80% pour tous les nouveaux fichiers
- **Complexité cyclomatique** : < 15 pour toutes les fonctions
- **Longueur fonctions** : < 50 lignes (sauf justification)
- **Pas de duplication** : DRY respecté
- **Aucun hardcoding** : Toutes les valeurs dans des constantes

## ✅ Validation

### Checklist

- [ ] `defaults.tsd` créé avec copyright
- [ ] Package `defaultactions` créé
- [ ] `IsDefault` ajouté à `ActionDefinition`
- [ ] Chargement automatique implémenté
- [ ] Validation de non-redéfinition
- [ ] 6 actions implémentées
- [ ] Tests complets (> 80% couverture)
- [ ] `make validate` passe
- [ ] Documentation complète

### Commandes de validation

```bash
# Tests unitaires
go test ./internal/defaultactions/...
go test ./rete/actions/...

# Tests d'intégration
go test ./constraint/... -run TestDefaultActions

# Validation complète
make validate

# Couverture
make test-coverage
```

## 🔄 Migration

### Rétrocompatibilité

- ✅ Le code existant continue de fonctionner
- ✅ L'action `print` existante est remplacée par la nouvelle
- ✅ Pas de breaking change dans l'API publique

### TODO pour l'appelant

Si le nouveau code n'est pas compatible avec l'existant :

```go
// TODO: Adapter les tests existants qui utilisent directement PrintAction
// Avant:
//   printAction := NewPrintAction(nil)
//   executor.RegisterAction(printAction)
//
// Après:
//   Les actions par défaut sont chargées automatiquement.
//   Plus besoin d'enregistrer Print manuellement.

// TODO: Mettre à jour les appels directs à RegisterDefaultActions()
// Avant:
//   executor.RegisterDefaultActions()
//
// Après:
//   Le chargement est automatique dans NewActionExecutor.
//   Supprimer les appels explicites.
```

## 📚 Références

- [common.md](../../.github/prompts/common.md) - Standards
- [Effective Go](https://go.dev/doc/effective_go)
- [go:embed documentation](https://pkg.go.dev/embed)
- Spécification TSD (docs internes)
