# Analyse du système actuel de déclaration d'actions

## 📋 Vue d'ensemble

Le système TSD actuel gère les actions à deux niveaux :
1. **Définition** : via `ActionDefinition` dans le package `constraint` (syntaxe TSD)
2. **Exécution** : via `ActionHandler` interface dans le package `rete` (implémentation Go)

## 🎯 Syntaxe de déclaration d'action

### Format TSD

```tsd
action nom_action(param1: type1, param2: type2, param3?: type3 = valeur_par_defaut)
```

### Exemples

```tsd
// Action simple
action notify(recipient: string, message: string)

// Action avec paramètre optionnel et valeur par défaut
action alert(severity: string, message: string, priority: number = 1)

// Action avec types utilisateur
action process_order(order: Order, customer: Customer)
```

## 📐 Structure AST

### ActionDefinition

**Fichier** : `constraint/constraint_types.go`

```go
type ActionDefinition struct {
    Type       string      `json:"type"`       // Toujours "actionDefinition"
    Name       string      `json:"name"`       // Nom de l'action
    Parameters []Parameter `json:"parameters"` // Liste des paramètres
}
```

### Parameter

```go
type Parameter struct {
    Name         string      `json:"name"`                   // Nom du paramètre
    Type         string      `json:"type"`                   // Type (primitif ou utilisateur)
    Optional     bool        `json:"optional"`               // Paramètre optionnel ?
    DefaultValue interface{} `json:"defaultValue,omitempty"` // Valeur par défaut
}
```

## 🔍 Parsing

### Grammar PEG

**Fichier** : `constraint/grammar/constraint.peg`

Le parser PEG génère automatiquement les structures `ActionDefinition` lors du parsing de fichiers `.tsd`.

### Code généré

**Fichier** : `constraint/parser.go` (généré par pigeon, NE PAS MODIFIER)

Le parser intègre `ActionDefinition` dans la structure `Program` :

```go
type Program struct {
    Types        []TypeDefinition
    Actions      []ActionDefinition      // ← Actions déclarées
    XupleSpaces  []XupleSpaceDeclaration
    Expressions  []Expression
    Facts        []Fact
    Resets       []Reset
    RuleRemovals []RuleRemoval
}
```

## ✅ Validation

### ActionValidator

**Fichier** : `constraint/action_validator.go`

Valide que :
- Les appels d'actions correspondent à des actions déclarées
- Le nombre d'arguments est correct (min/max selon paramètres requis/optionnels)
- Les types des arguments sont compatibles avec les types des paramètres
- Les types des paramètres existent (primitifs ou définis par l'utilisateur)
- Les valeurs par défaut correspondent au type du paramètre

```go
type ActionValidator struct {
    actions          map[string]*ActionDefinition
    types            map[string]*TypeDefinition
    functionRegistry *FunctionRegistry
}
```

### Méthodes principales

```go
// Valider un appel d'action
func (av *ActionValidator) ValidateActionCall(jobCall *JobCall, ruleVariables map[string]string) error

// Valider les définitions d'actions
func (av *ActionValidator) ValidateActionDefinitions() []error

// Récupérer une définition d'action
func (av *ActionValidator) GetActionDefinition(name string) (*ActionDefinition, bool)
```

## ⚙️ Exécution

### ActionHandler Interface

**Fichier** : `rete/action_handler.go`

```go
type ActionHandler interface {
    Execute(args []interface{}, ctx *ExecutionContext) error
    GetName() string
    Validate(args []interface{}) error
}
```

### ActionRegistry

Registre thread-safe des handlers d'actions :

```go
type ActionRegistry struct {
    handlers map[string]ActionHandler
    mu       sync.RWMutex
}

func (ar *ActionRegistry) Register(handler ActionHandler) error
func (ar *ActionRegistry) Get(actionName string) ActionHandler
func (ar *ActionRegistry) Has(actionName string) bool
```

### ActionExecutor

**Fichier** : `rete/action_executor.go`

Coordonne l'exécution des actions :

```go
type ActionExecutor struct {
    network       *ReteNetwork
    logger        *log.Logger
    enableLogging bool
    registry      *ActionRegistry
}

func NewActionExecutor(network *ReteNetwork, logger *log.Logger) *ActionExecutor
func (ae *ActionExecutor) RegisterDefaultActions()
func (ae *ActionExecutor) ExecuteAction(action Action, token *Token) error
```

## 🔧 Actions par défaut actuelles

### Hardcodage dans RegisterDefaultActions()

**Fichier** : `rete/action_executor.go:77-82`

```go
func (ae *ActionExecutor) RegisterDefaultActions() {
    // ❌ SEULE action par défaut actuellement - HARDCODÉE
    printAction := NewPrintAction(nil)
    if err := ae.registry.Register(printAction); err != nil {
        ae.logger.Printf("⚠️  Erreur enregistrement action print: %v", err)
    }
}
```

### PrintAction implémentation

**Fichier** : `rete/action_print.go`

```go
type PrintAction struct {
    output io.Writer
}

func (pa *PrintAction) Execute(args []interface{}, ctx *ExecutionContext) error
func (pa *PrintAction) GetName() string  // Retourne "print"
func (pa *PrintAction) Validate(args []interface{}) error
```

## ❌ Problèmes identifiés

### 1. Hardcoding des actions par défaut

- ❌ Seule `print` est enregistrée, directement dans le code
- ❌ Pas de définition centralisée des actions système
- ❌ Actions manquantes : `Log`, `Update`, `Insert`, `Retract`, `Xuple`

### 2. Incohérence définition/implémentation

- ⚠️ Les actions système ne sont pas déclarées via `ActionDefinition`
- ⚠️ Aucune validation de signature pour les actions natives
- ⚠️ Impossible de redéfinir une action système (pas de protection)

### 3. Manque d'extensibilité

- ❌ Ajout d'une nouvelle action système nécessite modification du code
- ❌ Pas de fichier de définition parsable
- ❌ Couplage fort entre définition et implémentation

## 💡 Besoins identifiés

### Actions système à implémenter

Selon la spécification :

1. **Print(message: string)** - Affichage console (✅ existe)
2. **Log(message: string)** - Logging système (❌ manquante)
3. **Update(fact: any)** - Mise à jour de fait (❌ manquante)
4. **Insert(fact: any)** - Insertion de fait (❌ manquante)
5. **Retract(id: string)** - Suppression de fait (❌ manquante)
6. **Xuple(xuplespace: string, fact: any)** - Création de xuple (❌ manquante)

### Mécanisme requis

1. **Chargement automatique** à l'initialisation
2. **Fichier de définition** parsé (pas de hardcoding)
3. **Marquage** des actions système (interdire redéfinition)
4. **Validation** de non-duplication
5. **Implémentations natives** liées aux handlers

## 🏗️ Architecture actuelle

```
┌─────────────────────────────────────────────────────────────┐
│                      Fichier .tsd                           │
│  action notify(recipient: string, message: string)          │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      │ Parsing (PEG)
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              constraint.Program                             │
│  Actions: []ActionDefinition                                │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      │ Validation
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              ActionValidator                                │
│  ValidateActionCall()                                       │
│  ValidateActionDefinitions()                                │
└─────────────────────┬───────────────────────────────────────┘
                      │
                      │ Exécution
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              rete.ActionExecutor                            │
│  registry: ActionRegistry                                   │
│    ├─ "print" → PrintAction (❌ hardcodé)                   │
│    └─ [actions utilisateur enregistrées dynamiquement]      │
└─────────────────────────────────────────────────────────────┘
```

## 📊 Métriques

- **Complexité** : Moyenne (séparation claire parsing/validation/exécution)
- **Couverture tests** :
  - `action_validator.go` : ✅ Bien testé
  - `action_handler.go` : ✅ Bien testé
  - `action_executor.go` : ✅ Bien testé
  - `action_print.go` : ✅ Bien testé
- **Lignes de code** :
  - action_validator.go : ~340 lignes
  - action_executor.go : ~200 lignes (+ fichiers séparés)
  - action_handler.go : ~150 lignes

## 🎯 Points forts

- ✅ Séparation claire des responsabilités
- ✅ Interface `ActionHandler` bien conçue
- ✅ Validation robuste des actions
- ✅ Support des paramètres optionnels et valeurs par défaut
- ✅ Thread-safety du registry
- ✅ Gestion d'erreurs complète

## ⚠️ Points à améliorer

- ❌ Actions système hardcodées
- ❌ Pas de fichier de définition centralisé
- ❌ Actions manquantes (Log, Update, Insert, Retract, Xuple)
- ⚠️ Pas de mécanisme de marquage des actions système
- ⚠️ Pas de protection contre redéfinition

## 📚 Références

- `constraint/constraint_types.go` - Structures AST
- `constraint/action_validator.go` - Validation
- `constraint/parser.go` - Parser généré
- `constraint/grammar/constraint.peg` - Grammaire
- `rete/action_handler.go` - Interface et registry
- `rete/action_executor.go` - Exécuteur
- `rete/action_print.go` - Implémentation Print
