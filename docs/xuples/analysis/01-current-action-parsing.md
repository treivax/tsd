# Analyse du Parsing des Actions - TSD

## 📋 Vue d'Ensemble

Ce document analyse en profondeur comment les actions sont actuellement définies et parsées dans le langage TSD.

## 🎯 Objectif

Comprendre le flux complet du parsing des actions, de la grammaire PEG jusqu'aux structures de données internes utilisées par le réseau RETE.

---

## 1. Grammaire PEG des Actions

### 1.1 Fichier Source

**Emplacement** : `constraint/grammar/constraint.peg`  
**Taille** : 692 lignes  
**Générateur** : pigeon (commande : `pigeon -o parser.go constraint.peg`)

### 1.2 Définition d'Action (ActionDefinition)

Les actions peuvent être pré-définies avec une signature (optionnel) :

```peg
ActionDefinition <- "action" _ name:IdentName _ "(" _ params:ParameterList? _ ")" {
    if params == nil {
        params = []interface{}{}
    }
    return map[string]interface{}{
        "type": "actionDefinition",
        "name": name,
        "parameters": params,
    }, nil
}
```

**Format TSD** :
```tsd
action notify(recipient: string, message: string, priority: number = 1)
```

**Structure AST produite** :
```json
{
  "type": "actionDefinition",
  "name": "notify",
  "parameters": [
    {"name": "recipient", "type": "string", "optional": false},
    {"name": "message", "type": "string", "optional": false},
    {"name": "priority", "type": "number", "optional": true, "defaultValue": 1}
  ]
}
```

**Référence** : `constraint/grammar/constraint.peg` lignes 113-122

### 1.3 Action dans une Règle

Les actions sont utilisées dans les expressions de règles :

```peg
Action <- first:JobCall rest:(_ "," _ JobCall)* {
    jobs := []interface{}{first}
    if rest != nil {
        for _, item := range rest.([]interface{}) {
            jobs = append(jobs, item.([]interface{})[3])
        }
    }
    return map[string]interface{}{
        "type": "action",
        "jobs": jobs,
    }, nil
}
```

**Référence** : `constraint/grammar/constraint.peg` lignes 442-453

### 1.4 JobCall (Appel de Job)

Un JobCall représente l'appel d'une fonction/action avec ses arguments :

```peg
JobCall <- name:IdentName _ "(" _ args:ArgumentList? _ ")" {
    if args == nil {
        args = []interface{}{}
    }
    return map[string]interface{}{
        "type": "jobCall",
        "name": name,
        "args": args,
    }, nil
}
```

**Format TSD** :
```tsd
print(user.name, "a", user.age, "ans")
```

**Structure AST produite** :
```json
{
  "type": "jobCall",
  "name": "print",
  "args": [
    {"type": "fieldAccess", "object": "user", "field": "name"},
    {"type": "stringLiteral", "value": "a"},
    {"type": "fieldAccess", "object": "user", "field": "age"},
    {"type": "stringLiteral", "value": "ans"}
  ]
}
```

**Référence** : `constraint/grammar/constraint.peg` lignes 455-464

### 1.5 ArgumentList

Les arguments supportent des expressions arithmétiques complètes :

```peg
ArgumentList <- first:ArithmeticExpr rest:(_ "," _ ArithmeticExpr)* {
    arguments := []interface{}{first}
    if rest != nil {
        for _, item := range rest.([]interface{}) {
            arguments = append(arguments, item.([]interface{})[3])
        }
    }
    return arguments, nil
}
```

**Référence** : `constraint/grammar/constraint.peg` lignes 466-474

---

## 2. Structures de Données AST

### 2.1 Type ActionDefinition

**Emplacement** : `constraint/constraint_types.go` lignes 34-40

```go
// ActionDefinition represents a user-defined action with its signature.
// Example: action notify(recipient: string, message: string, priority: number = 1)
type ActionDefinition struct {
	Type       string      `json:"type"`       // Always "actionDefinition"
	Name       string      `json:"name"`       // The action name (e.g., "notify")
	Parameters []Parameter `json:"parameters"` // List of parameters for the action
}
```

### 2.2 Type Parameter

**Emplacement** : `constraint/constraint_types.go` lignes 42-49

```go
// Parameter represents a single parameter within an action definition.
// It contains the parameter name, type, whether it's optional, and an optional default value.
type Parameter struct {
	Name         string      `json:"name"`                   // Parameter name (e.g., "recipient", "priority")
	Type         string      `json:"type"`                   // Parameter type (e.g., "string", "number", "bool", or a user-defined type like "Person")
	Optional     bool        `json:"optional"`               // Whether the parameter is optional (marked with ?)
	DefaultValue interface{} `json:"defaultValue,omitempty"` // Default value if provided
}
```

### 2.3 Type Action

**Emplacement** : `constraint/constraint_types.go` lignes 191-211

```go
// Action represents an action to execute when constraints are satisfied.
// It defines what job(s) should be performed and with what parameters.
// Supports both single action (Job field, for backward compatibility) and
// multiple actions (Jobs field, new format).
type Action struct {
	Type string    `json:"type"`           // Always "action"
	Job  *JobCall  `json:"job,omitempty"`  // Single job (backward compatibility)
	Jobs []JobCall `json:"jobs,omitempty"` // Multiple jobs (new format)
}

// GetJobs returns the list of jobs to execute.
// It handles both the old format (single Job) and new format (multiple Jobs).
func (a *Action) GetJobs() []JobCall {
	if len(a.Jobs) > 0 {
		return a.Jobs
	}
	if a.Job != nil {
		return []JobCall{*a.Job}
	}
	return []JobCall{}
}
```

**Note importante** : L'Action supporte deux formats pour rétrocompatibilité :
- **Ancien format** : `Job *JobCall` (une seule action)
- **Nouveau format** : `Jobs []JobCall` (actions multiples)

### 2.4 Type JobCall

**Emplacement** : `constraint/constraint_types.go` lignes 213-219

```go
// JobCall represents a specific job/function call within an action.
// It specifies the job name and arguments to pass.
type JobCall struct {
	Type string        `json:"type"` // Always "jobCall"
	Name string        `json:"name"` // Job/function name
	Args []interface{} `json:"args"` // Arguments to pass to the job
}
```

---

## 3. Flux de Parsing

### 3.1 Étapes du Parsing

```
Fichier .tsd
    ↓
[Pigeon Parser]
    ↓
AST brut (map[string]interface{})
    ↓
[Conversion en structures Go]
    ↓
constraint.Program
    ├── Types []TypeDefinition
    ├── Actions []ActionDefinition
    ├── Expressions []Expression
    │       └── Action *Action
    │               └── Jobs []JobCall
    └── Facts []Fact
```

### 3.2 Exemple de Parsing Complet

**Input TSD** :
```tsd
type User(#id: string, name: string, age: number)

action notify(recipient: string, message: string)

rule user_adult: {u: User} / u.age >= 18 ==> notify(u.name, "Vous êtes majeur")

User(id: "U001", name: "Alice", age: 25)
```

**Output AST** :
```json
{
  "types": [
    {
      "type": "typeDefinition",
      "name": "User",
      "fields": [
        {"name": "id", "type": "string", "isPrimaryKey": true},
        {"name": "name", "type": "string", "isPrimaryKey": false},
        {"name": "age", "type": "number", "isPrimaryKey": false}
      ]
    }
  ],
  "actions": [
    {
      "type": "actionDefinition",
      "name": "notify",
      "parameters": [
        {"name": "recipient", "type": "string", "optional": false},
        {"name": "message", "type": "string", "optional": false}
      ]
    }
  ],
  "expressions": [
    {
      "type": "expression",
      "ruleId": "user_adult",
      "set": {
        "type": "set",
        "variables": [
          {"type": "typedVariable", "name": "u", "dataType": "User"}
        ]
      },
      "constraints": {
        "type": "comparison",
        "left": {"type": "fieldAccess", "object": "u", "field": "age"},
        "operator": ">=",
        "right": {"type": "numberLiteral", "value": 18}
      },
      "action": {
        "type": "action",
        "jobs": [
          {
            "type": "jobCall",
            "name": "notify",
            "args": [
              {"type": "fieldAccess", "object": "u", "field": "name"},
              {"type": "stringLiteral", "value": "Vous êtes majeur"}
            ]
          }
        ]
      }
    }
  ],
  "facts": [
    {
      "type": "fact",
      "typeName": "User",
      "fields": [
        {"name": "id", "value": {"type": "string", "value": "U001"}},
        {"name": "name", "value": {"type": "string", "value": "Alice"}},
        {"name": "age", "value": {"type": "number", "value": 25}}
      ]
    }
  ]
}
```

---

## 4. Stockage des Définitions d'Actions

### 4.1 Dans constraint.Program

Les actions définies sont stockées dans `Program.Actions` :

**Emplacement** : `constraint/constraint_types.go` lignes 7-16

```go
type Program struct {
	Types        []TypeDefinition   `json:"types"`        // Type definitions declared in the program
	Actions      []ActionDefinition `json:"actions"`      // Action definitions with their signatures
	Expressions  []Expression       `json:"expressions"`  // Constraint expressions/rules
	Facts        []Fact             `json:"facts"`        // Facts parsed from the program
	Resets       []Reset            `json:"resets"`       // Reset instructions to clear the system
	RuleRemovals []RuleRemoval      `json:"ruleRemovals"` // Rule removal commands
}
```

### 4.2 Dans Expression

Chaque expression contient son action :

```go
type Expression struct {
	Type        string      `json:"type"`               // Always "expression"
	RuleId      string      `json:"ruleId"`             // Unique identifier for the rule
	Set         Set         `json:"set,omitempty"`      // Set of variables (single pattern, backward compatibility)
	Patterns    []Set       `json:"patterns,omitempty"` // Multiple pattern blocks (aggregation with joins)
	Constraints interface{} `json:"constraints"`        // Constraints to evaluate
	Action      *Action     `json:"action,omitempty"`   // Action to execute when constraints match
}
```

**Référence** : `constraint/constraint_types.go` lignes 51-61

---

## 5. Validations Existantes

### 5.1 Validation des Variables dans les Actions

**Emplacement** : `constraint/constraint_actions.go`

La fonction `ValidateAction` vérifie que tous les arguments d'action référencent des variables valides définies dans l'expression :

```go
func ValidateAction(program Program, action Action, expressionIndex int) error {
	if expressionIndex >= len(program.Expressions) {
		return fmt.Errorf("index d'expression invalide: %d", expressionIndex)
	}

	expression := program.Expressions[expressionIndex]

	// Créer une map des variables disponibles dans l'expression
	availableVars := make(map[string]bool)

	// Ajouter les variables du Set principal (ancien format, rétrocompatibilité)
	for _, variable := range expression.Set.Variables {
		availableVars[variable.Name] = true
	}

	// Ajouter les variables des Patterns multiples (nouveau format avec agrégation)
	for _, pattern := range expression.Patterns {
		for _, variable := range pattern.Variables {
			availableVars[variable.Name] = true
		}
	}

	// Obtenir tous les jobs (supporte ancien et nouveau format)
	jobs := action.GetJobs()

	// Vérifier que tous les arguments de chaque job référencent des variables valides
	for _, job := range jobs {
		for _, arg := range job.Args {
			// Extraire les variables utilisées dans l'argument
			vars := extractVariablesFromArg(arg)
			for _, varName := range vars {
				if !availableVars[varName] {
					return fmt.Errorf("action %s: argument contient la variable '%s' qui ne correspond à aucune variable de l'expression", job.Name, varName)
				}
			}
		}
	}

	return nil
}
```

**Référence** : `constraint/constraint_actions.go` lignes 11-51

### 5.2 Extraction de Variables

Le système supporte plusieurs types d'arguments :

```go
func extractVariablesByType(argType string, argMap map[string]interface{}) []string {
	switch argType {
	case "fieldAccess":
		return extractFromFieldAccess(argMap)
	case "variable":
		return extractFromVariable(argMap)
	case ArgTypeStringLiteral, "string", ArgTypeNumberLiteral, "number", ArgTypeBoolLiteral, ValueTypeBoolean:
		return []string{} // Literals ne contiennent pas de variables
	case ArgTypeFunctionCall:
		return extractFromFunctionCall(argMap)
	default:
		if isBinaryOperationType(argType) {
			return extractFromBinaryOp(argMap)
		}
		return []string{}
	}
}
```

**Référence** : `constraint/constraint_actions.go` lignes 70-87

**Types d'arguments supportés** :
- `fieldAccess` : Accès à un champ (ex: `user.name`)
- `variable` : Variable simple (ex: `x`)
- `stringLiteral`, `numberLiteral`, `booleanLiteral` : Littéraux
- `functionCall` : Appel de fonction (ex: `LENGTH(user.name)`)
- Opérations binaires (ex: `user.age + 5`)

### 5.3 Détection de Doublons

**Question** : Y a-t-il validation des doublons d'actions ?

**Réponse** : Actuellement, **NON**. Il n'y a pas de validation de doublons pour les `ActionDefinition` dans le code de parsing ou validation analysé.

Les `ActionDefinition` sont simplement accumulées dans `Program.Actions` sans vérification d'unicité du nom.

**Fichiers examinés** :
- `constraint/constraint_actions.go` - Validation variables uniquement
- `constraint/action_validator.go` - Validation arguments d'action
- `constraint/parser.go` - Code généré, pas de validation custom

---

## 6. Points d'Intervention pour la Refonte

### 6.1 Actions Parsées

**État actuel** :
- Les actions sont parsées et stockées dans `constraint.Program.Actions`
- Chaque expression a son `Action *Action` avec les jobs à exécuter
- Pas de validation de doublons d'actions
- Pas de registry centralisé des actions définies

**Points d'intervention** :
1. **Après parsing** : Valider unicité des noms d'actions
2. **Conversion vers RETE** : Créer un registry d'actions disponibles
3. **Terminal nodes** : Référencer les actions par nom plutôt que par structure inline

### 6.2 Structure Action dans Expression

**État actuel** :
```go
type Expression struct {
    // ...
    Action *Action `json:"action,omitempty"` // Action inline
}
```

**Proposition pour xuples** :
- Stocker uniquement le nom de l'action dans Expression
- Résoudre l'action via un registry au moment de la construction du réseau RETE
- Permettre des actions par défaut si non définies

### 6.3 JobCall et Arguments

**État actuel** :
- `JobCall.Args` est `[]interface{}` (très flexible)
- Arguments évalués au runtime par `ActionExecutor.evaluateArgument()`
- Supporte expressions complexes, accès champs, opérations arithmétiques

**Conservation pour xuples** :
- Cette flexibilité est nécessaire et doit être conservée
- L'évaluation runtime des arguments est correcte
- Aucune modification nécessaire de `JobCall`

---

## 7. Exemples de Code Pertinents

### 7.1 Parsing d'un Fichier TSD

```go
// Dans constraint/api.go (hypothétique, basé sur la structure)
func ParseConstraintFile(filename string) (*Program, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    // Parse avec pigeon
    result, err := Parse(filename, content)
    if err != nil {
        return nil, err
    }
    
    // Convertir map vers Program
    program := convertToProgram(result)
    
    return program, nil
}
```

### 7.2 Utilisation dans RETE

**Emplacement hypothétique** : Lors de la construction du réseau RETE

```go
// Construction du TerminalNode avec action
for _, expr := range program.Expressions {
    // Créer les nœuds RETE pour l'expression
    // ...
    
    // Créer le TerminalNode avec l'action
    terminalNode := rete.NewTerminalNode(nodeID, expr.Action, storage)
    
    // L'action est stockée directement dans le TerminalNode
}
```

**Référence** : `rete/node_terminal.go` lignes 12-30

---

## 8. Synthèse et Observations

### 8.1 Points Forts

✅ **Grammaire claire et extensible** : La grammaire PEG est bien structurée  
✅ **Support multi-actions** : Format `Jobs []JobCall` permet plusieurs actions  
✅ **Validation des variables** : Vérification que les variables existent  
✅ **Arguments flexibles** : Support expressions complexes, accès champs, etc.  
✅ **Rétrocompatibilité** : Support ancien format `Job` et nouveau `Jobs`

### 8.2 Points d'Amélioration

⚠️ **Pas de validation unicité** : Noms d'actions peuvent être dupliqués  
⚠️ **Pas de registry centralisé** : Actions définies mais non indexées  
⚠️ **Action inline dans Expression** : Couplage fort entre règle et action  
⚠️ **Pas de résolution d'action** : Pas de mécanisme pour référencer une action par nom

### 8.3 Recommandations pour Xuples

1. **Ajouter validation unicité** des noms d'actions lors du parsing
2. **Créer un ActionRegistry** après parsing pour indexer les `ActionDefinition`
3. **Permettre référence par nom** : `action: "notify"` au lieu de structure complète
4. **Actions par défaut** : Créer des actions intégrées (print, assert, retract, etc.)
5. **Conserver flexibilité** des arguments et leur évaluation runtime

---

## 9. Fichiers de Référence

| Fichier | Description | Lignes clés |
|---------|-------------|-------------|
| `constraint/grammar/constraint.peg` | Grammaire PEG complète | 113-122 (ActionDefinition), 442-464 (Action/JobCall) |
| `constraint/constraint_types.go` | Structures AST | 34-49 (ActionDefinition), 191-219 (Action/JobCall) |
| `constraint/constraint_actions.go` | Validation des actions | 11-127 (ValidateAction, extraction variables) |
| `constraint/action_validator.go` | Validateur d'actions | Validation avancée des arguments |
| `rete/node_terminal.go` | Utilisation des actions | 12-30 (NewTerminalNode avec Action) |

---

**Date de création** : 2025-12-17  
**Auteur** : Analyse automatique pour refonte xuples  
**Statut** : ✅ Complet
