# Analyse de l'ActionExecutor et son Interface - TSD

## 📋 Vue d'Ensemble

Ce document analyse en profondeur l'ActionExecutor, son interface, et comment les actions sont exécutées dans le réseau RETE.

## 🎯 Objectif

Comprendre le contrat d'exécution des actions, l'architecture de l'ActionExecutor, et comment il s'intègre avec le réseau RETE.

---

## 1. Interface ActionHandler

### 1.1 Définition

**Emplacement** : `rete/action_handler.go` lignes 12-25

```go
// ActionHandler définit l'interface pour les gestionnaires d'actions personnalisées.
// Chaque action peut avoir son propre handler qui définit son comportement.
type ActionHandler interface {
	// Execute exécute l'action avec les arguments évalués fournis.
	// Retourne une erreur si l'exécution échoue.
	Execute(args []interface{}, ctx *ExecutionContext) error

	// GetName retourne le nom de l'action gérée par ce handler.
	GetName() string

	// Validate valide que les arguments sont corrects pour cette action.
	// Cette validation est optionnelle et peut retourner nil si aucune validation spécifique n'est nécessaire.
	Validate(args []interface{}) error
}
```

**Méthodes** :
1. **Execute** : Exécute l'action avec arguments évalués
2. **GetName** : Retourne nom de l'action (pour registry)
3. **Validate** : Validation optionnelle des arguments

**Contrat** :
- ✅ Arguments déjà évalués (pas de parsing à faire)
- ✅ Contexte fourni pour accès aux variables si nécessaire
- ✅ Validation séparée de l'exécution
- ✅ Gestion d'erreur explicite

### 1.2 Exemple d'Implémentation : PrintAction

**Emplacement** : `rete/action_print.go`

```go
// PrintAction implémente une action d'affichage simple
type PrintAction struct {
	logger *log.Logger
}

func NewPrintAction(logger *log.Logger) *PrintAction {
	if logger == nil {
		logger = log.Default()
	}
	return &PrintAction{logger: logger}
}

func (pa *PrintAction) GetName() string {
	return "print"
}

func (pa *PrintAction) Validate(args []interface{}) error {
	// Print accepte n'importe quel nombre d'arguments
	return nil
}

func (pa *PrintAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	// Convertir tous les arguments en string et afficher
	var parts []string
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%v", arg))
	}
	message := strings.Join(parts, " ")
	pa.logger.Printf("🖨️  PRINT: %s", message)
	return nil
}
```

**Référence** : `rete/action_print.go` lignes 10-100 (hypothétique basé sur contexte)

---

## 2. ActionRegistry

### 2.1 Structure

**Emplacement** : `rete/action_handler.go` lignes 27-38

```go
// ActionRegistry gère l'enregistrement et la récupération des handlers d'actions.
type ActionRegistry struct {
	handlers map[string]ActionHandler
	mu       sync.RWMutex
}

// NewActionRegistry crée un nouveau registry d'actions.
func NewActionRegistry() *ActionRegistry {
	return &ActionRegistry{
		handlers: make(map[string]ActionHandler),
	}
}
```

**Caractéristiques** :
- **Thread-safe** : Utilise `sync.RWMutex`
- **Indexation** : Par nom d'action (string → ActionHandler)
- **Dynamique** : Enregistrement/désenregistrement à runtime

### 2.2 Méthodes du Registry

**Emplacement** : `rete/action_handler.go` lignes 40-135

```go
// Register enregistre un handler d'action.
// Si un handler existe déjà pour ce nom, il est remplacé.
func (ar *ActionRegistry) Register(handler ActionHandler) error

// Unregister supprime un handler d'action du registry.
func (ar *ActionRegistry) Unregister(actionName string)

// Get récupère un handler d'action par son nom.
// Retourne nil si aucun handler n'est enregistré pour ce nom.
func (ar *ActionRegistry) Get(actionName string) ActionHandler

// Has vérifie si un handler est enregistré pour une action donnée.
func (ar *ActionRegistry) Has(actionName string) bool

// GetAll retourne une copie de tous les handlers enregistrés.
func (ar *ActionRegistry) GetAll() map[string]ActionHandler

// Count retourne le nombre de handlers enregistrés.
func (ar *ActionRegistry) Count() int

// Clear supprime tous les handlers du registry.
func (ar *ActionRegistry) Clear()

// RegisterMultiple enregistre plusieurs handlers en une seule opération.
func (ar *ActionRegistry) RegisterMultiple(handlers []ActionHandler) error

// GetRegisteredNames retourne la liste des noms d'actions enregistrées.
func (ar *ActionRegistry) GetRegisteredNames() []string
```

**Particularités** :
- ✅ RWMutex permet lectures concurrentes
- ✅ Méthodes atomiques (lock/unlock par opération)
- ✅ Remplace handler existant sans erreur (facilite hot-reload)

---

## 3. ActionExecutor

### 3.1 Structure

**Emplacement** : `rete/action_executor.go` lignes 37-42

```go
type ActionExecutor struct {
	network       *ReteNetwork
	logger        *log.Logger
	enableLogging bool
	registry      *ActionRegistry
}
```

**Champs** :
- **network** : Référence au réseau RETE (pour types, etc.)
- **logger** : Logger pour journalisation des exécutions
- **enableLogging** : Flag pour activer/désactiver logs
- **registry** : Registry des handlers d'actions disponibles

### 3.2 Constructeur

**Emplacement** : `rete/action_executor.go` lignes 44-69

```go
// NewActionExecutor crée un nouveau exécuteur d'actions.
//
// Initialise le registry et enregistre les actions par défaut (print, etc.).
//
// Paramètres :
//   - network : réseau RETE
//   - logger : logger pour journalisation (utilise log.Default() si nil)
//
// Retourne :
//   - *ActionExecutor : exécuteur initialisé
func NewActionExecutor(network *ReteNetwork, logger *log.Logger) *ActionExecutor {
	if logger == nil {
		logger = log.Default()
	}
	ae := &ActionExecutor{
		network:       network,
		logger:        logger,
		enableLogging: true,
		registry:      NewActionRegistry(),
	}

	// Enregistrer les actions par défaut
	ae.RegisterDefaultActions()

	return ae
}
```

**Actions par défaut** :
- **print** : Affichage de valeurs

### 3.3 RegisterDefaultActions

**Emplacement** : `rete/action_executor.go` lignes 71-83

```go
// RegisterDefaultActions enregistre les actions par défaut disponibles.
//
// Actions enregistrées :
//   - print : affichage de valeurs
//
// Cette méthode est appelée automatiquement par NewActionExecutor.
func (ae *ActionExecutor) RegisterDefaultActions() {
	// Enregistrer l'action print
	printAction := NewPrintAction(nil)
	if err := ae.registry.Register(printAction); err != nil {
		ae.logger.Printf("⚠️  Erreur enregistrement action print: %v", err)
	}
}
```

**Actions par défaut actuelles** :
- ✅ print

**Actions proposées pour implémentation future** :
- assert (ajouter un fait)
- retract (retirer un fait)
- modify (modifier un fait)
- halt (arrêter le moteur)
- log (journaliser)

---

## 4. Exécution d'Actions

### 4.1 ExecuteAction (Point d'Entrée)

**Emplacement** : `rete/action_executor.go` lignes 100-144

```go
// ExecuteAction exécute une action avec les faits fournis par le token.
//
// Process :
//  1. Valide les paramètres (action et token non nil)
//  2. Récupère tous les jobs de l'action
//  3. Crée un contexte d'exécution avec les bindings du token
//  4. Exécute chaque job en séquence avec récupération sur panic
//
// Thread-Safety :
//   - Cette méthode est thread-safe
//   - Le contexte d'exécution est isolé par appel
//   - Les panics sont récupérés et convertis en erreurs
//
// Paramètres :
//   - action : action à exécuter (peut contenir plusieurs jobs)
//   - token : token contenant les faits et bindings disponibles
//
// Retourne :
//   - error : erreur si l'exécution échoue ou si paramètres invalides
func (ae *ActionExecutor) ExecuteAction(action *Action, token *Token) error {
	if action == nil {
		return fmt.Errorf("action is nil")
	}
	if token == nil {
		return fmt.Errorf("token is nil")
	}

	// Obtenir tous les jobs à exécuter
	jobs := action.GetJobs()

	// Créer un contexte d'exécution avec les faits disponibles
	ctx := NewExecutionContext(token, ae.network)
	if ctx == nil {
		return fmt.Errorf("échec création contexte d'exécution")
	}

	// Exécuter chaque job en séquence
	for i, job := range jobs {
		if err := ae.executeJob(job, ctx, i); err != nil {
			return fmt.Errorf("erreur exécution job %s (index %d): %w", job.Name, i, err)
		}
	}

	return nil
}
```

**Flux** :
1. Validation paramètres
2. Extraction jobs (supporte multi-jobs)
3. Création ExecutionContext
4. Exécution séquentielle des jobs

### 4.2 executeJob (Exécution Individuelle)

**Emplacement** : `rete/action_executor.go` lignes 146-212

```go
// executeJob exécute un job individuel avec récupération sur panic.
//
// Process :
//  1. Log l'action (si activé)
//  2. Évalue tous les arguments
//  3. Recherche le handler dans le registry
//  4. Valide les arguments (si handler définit une validation)
//  5. Exécute le handler avec récupération sur panic
//
// Thread-safety :
//   - La méthode est thread-safe grâce au RWMutex du registry
//   - Le panic dans un handler est converti en erreur
//
// Paramètres :
//   - job : job à exécuter
//   - ctx : contexte d'exécution
//   - jobIndex : index du job dans la séquence (pour debug)
//
// Retourne :
//   - error : erreur si l'exécution échoue ou si panic
func (ae *ActionExecutor) executeJob(job JobCall, ctx *ExecutionContext, jobIndex int) (err error) {
	// Récupération sur panic
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic dans exécution action '%s': %v", job.Name, r)
			ae.logger.Printf("❌ PANIC RÉCUPÉRÉ dans action '%s': %v", job.Name, r)
		}
	}()

	// Logger l'action
	if ae.enableLogging {
		ae.logAction(job, ctx)
	}

	// Évaluer les arguments
	evaluatedArgs := make([]interface{}, 0, len(job.Args))
	for i, arg := range job.Args {
		evaluated, err := ae.evaluateArgument(arg, ctx)
		if err != nil {
			return fmt.Errorf("erreur évaluation argument %d de l'action '%s': %w", i, job.Name, err)
		}
		evaluatedArgs = append(evaluatedArgs, evaluated)
	}

	// Vérifier si un handler est enregistré pour cette action
	handler := ae.registry.Get(job.Name)
	if handler != nil {
		// Valider les arguments (optionnel)
		if err := handler.Validate(evaluatedArgs); err != nil {
			return fmt.Errorf("validation échouée pour action '%s': %w", job.Name, err)
		}

		// Exécuter l'action via son handler
		if err := handler.Execute(evaluatedArgs, ctx); err != nil {
			return fmt.Errorf("exécution échouée pour action '%s': %w", job.Name, err)
		}

		// Logger le succès
		ae.logger.Printf("🎯 ACTION EXÉCUTÉE: %s(%v)", job.Name, formatArgs(evaluatedArgs))
	} else {
		// Aucun handler défini : comportement par défaut (simple log)
		ae.logger.Printf("📋 ACTION NON DÉFINIE (log uniquement): %s(%v)", job.Name, formatArgs(evaluatedArgs))
	}

	return nil
}
```

**Étapes détaillées** :
1. **Panic recovery** : defer/recover pour sécurité
2. **Logging** : Si activé
3. **Évaluation arguments** : Via `evaluateArgument`
4. **Recherche handler** : Dans registry
5. **Validation** : `handler.Validate(args)`
6. **Exécution** : `handler.Execute(args, ctx)`
7. **Log succès** : Confirmation

**Comportement sans handler** :
- ⚠️ Log uniquement (pas d'erreur)
- Permet actions "fantômes" pour debug
- Facilite développement incrémental

---

## 5. ExecutionContext

### 5.1 Structure

**Emplacement** : `rete/action_executor_context.go` lignes 26-30

```go
type ExecutionContext struct {
	token    *Token
	network  *ReteNetwork
	bindings *BindingChain
}
```

**Champs** :
- **token** : Token source (contient Facts et Metadata)
- **network** : Réseau RETE (accès types, etc.)
- **bindings** : BindingChain immuable (variables → facts)

### 5.2 Constructeur

**Emplacement** : `rete/action_executor_context.go` lignes 32-63

```go
// NewExecutionContext crée un nouveau contexte d'exécution.
//
// Le contexte référence directement la chaîne de bindings du token,
// sans copie, garantissant l'immutabilité et la performance.
//
// Validation :
//   - Si token est nil, crée un contexte avec bindings vides
//   - Si network est nil, les fonctionnalités dépendant du network (type validation, etc.) ne seront pas disponibles
//
// Note : Un network nil est acceptable pour les tests unitaires simples,
// mais dans un contexte de production, le network devrait toujours être fourni.
//
// Paramètres :
//   - token : token contenant les faits et bindings (peut être nil)
//   - network : réseau RETE pour accès aux types (peut être nil pour tests simples)
//
// Retourne :
//   - *ExecutionContext : contexte d'exécution initialisé
func NewExecutionContext(token *Token, network *ReteNetwork) *ExecutionContext {
	ctx := &ExecutionContext{
		token:    token,
		network:  network,
		bindings: nil,
	}

	// Référencer directement la chaîne de bindings du token si disponible
	if token != nil {
		ctx.bindings = token.Bindings
	}

	return ctx
}
```

**Points importants** :
- ✅ Pas de copie : référence directe à BindingChain
- ✅ Accepte token nil (tests)
- ✅ Accepte network nil (tests simples)

### 5.3 GetVariable

**Emplacement** : `rete/action_executor_context.go` lignes 65-90

```go
// GetVariable récupère un fait par nom de variable.
//
// Utilise la BindingChain pour rechercher le fait lié à la variable.
// Retourne nil si la variable n'existe pas dans le contexte.
//
// Complexité : O(n) où n est le nombre de bindings (typiquement < 10)
//
// Paramètres :
//   - name : nom de la variable (ex: "user", "order", "task")
//
// Retourne :
//   - *Fact : pointeur vers le fait si trouvé, nil sinon
//
// Exemple :
//
//	user := ctx.GetVariable("user")
//	if user == nil {
//	    return fmt.Errorf("variable user non trouvée")
//	}
//	userName := user.Fields["name"]
func (ctx *ExecutionContext) GetVariable(name string) *Fact {
	if ctx.bindings == nil {
		return nil
	}
	return ctx.bindings.Get(name)
}
```

---

## 6. Évaluation des Arguments

### 6.1 evaluateArgument (Méthode Principale)

**Emplacement** : `rete/action_executor_evaluation.go` lignes 11-130

**Types supportés** :
1. **Littéraux simples** : string, float64, bool, int, int64
2. **Littéraux typés** : `{type: "string", value: "..."}`
3. **Variables** : `{type: "variable", name: "u"}`
4. **FieldAccess** : `{type: "fieldAccess", object: "u", field: "name"}`
5. **FactCreation** : `{type: "factCreation", ...}`
6. **FactModification** : `{type: "factModification", ...}`
7. **BinaryOperation** : `{type: "addition", left: ..., right: ...}`
8. **Cast** : `{type: "cast", ...}`

**Signature** :
```go
func (ae *ActionExecutor) evaluateArgument(arg interface{}, ctx *ExecutionContext) (interface{}, error)
```

### 6.2 Évaluation FieldAccess

**Exemple** :
```go
case "fieldAccess":
	objectName, ok := argMap["object"].(string)
	if !ok {
		return nil, fmt.Errorf("nom d'objet invalide dans fieldAccess")
	}
	fieldName, ok := argMap["field"].(string)
	if !ok {
		return nil, fmt.Errorf("nom de champ invalide dans fieldAccess")
	}

	fact := ctx.GetVariable(objectName)
	if fact == nil {
		// Message d'erreur détaillé avec liste des variables disponibles
		availableVars := []string{}
		if ctx.bindings != nil {
			availableVars = ctx.bindings.Variables()
		}
		return nil, fmt.Errorf(
			"❌ Erreur d'exécution d'action:\n"+
				"   Variable '%s' non trouvée dans le contexte\n"+
				"   Variables disponibles: %v\n"+
				"   Vérifiez que la règle déclare bien cette variable dans sa clause de pattern",
			objectName, availableVars,
		)
	}
	
	// Accès au champ
	value, exists := fact.GetField(fieldName)
	if !exists {
		return nil, fmt.Errorf("champ '%s' non trouvé dans le fait de type '%s'", fieldName, fact.Type)
	}
	return value, nil
```

**Messages d'erreur** :
- ✅ Très détaillés
- ✅ Liste variables disponibles
- ✅ Suggestions pour corriger

### 6.3 Évaluation BinaryOperation

**Support** :
- Arithmétique : `+`, `-`, `*`, `/`, `%`
- Comparaison : `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logique : `AND`, `OR`

**Évaluation récursive** :
```go
case "addition", "subtraction", "multiplication", "division", "modulo":
	left, err := ae.evaluateArgument(argMap["left"], ctx)
	if err != nil {
		return nil, err
	}
	right, err := ae.evaluateArgument(argMap["right"], ctx)
	if err != nil {
		return nil, err
	}
	return performArithmeticOperation(argType, left, right)
```

---

## 7. Intégration avec le Réseau RETE

### 7.1 ReteNetwork.ActionExecutor

**Stockage** :
```go
type ReteNetwork struct {
	// ... autres champs
	ActionExecutor *ActionExecutor
}
```

**Initialisation** :
```go
network := NewReteNetwork(storage)
network.ActionExecutor = NewActionExecutor(network, logger)
```

### 7.2 TerminalNode → ActionExecutor

**Flux** :
```
TerminalNode.ActivateLeft(token)
    ↓
TerminalNode.executeAction(token)
    ↓
network.ActionExecutor.ExecuteAction(action, token)
    ↓
ActionExecutor.executeJob(job, ctx, index)
    ↓
handler.Execute(evaluatedArgs, ctx)
```

**Code** (`rete/node_terminal.go` lignes 164-170) :
```go
// Exécuter réellement l'action avec l'ActionExecutor
network := tn.BaseNode.GetNetwork()
if network != nil && network.ActionExecutor != nil {
	return network.ActionExecutor.ExecuteAction(tn.Action, token)
}

return nil
```

---

## 8. Thread-Safety et Concurrence

### 8.1 ActionExecutor

✅ **Thread-safe** :
- Registry utilise `sync.RWMutex`
- Pas d'état mutable (sauf registry)
- Logger thread-safe (log.Logger)

### 8.2 ExecutionContext

✅ **Thread-safe** :
- Immutable (BindingChain immuable)
- Pas de modification de token
- Lecture seule

### 8.3 Panic Recovery

✅ **Robuste** :
```go
defer func() {
	if r := recover(); r != nil {
		err = fmt.Errorf("panic dans exécution action '%s': %v", job.Name, r)
		ae.logger.Printf("❌ PANIC RÉCUPÉRÉ dans action '%s': %v", job.Name, r)
	}
}()
```

**Garantie** : Un panic dans un handler ne crash pas le moteur RETE

---

## 9. Propositions pour Actions par Défaut

### 9.1 Actions Standard Proposées

```go
// assert : Ajouter un fait dans le réseau
type AssertAction struct {
	network *ReteNetwork
}

func (a *AssertAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	// args[0] = fact à ajouter
	fact, ok := args[0].(*Fact)
	if !ok {
		return fmt.Errorf("assert attend un fait en argument")
	}
	return a.network.AddFact(fact)
}

// retract : Retirer un fait du réseau
type RetractAction struct {
	network *ReteNetwork
}

func (r *RetractAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	// args[0] = variable ou ID du fait à retirer
	// Implementation...
}

// modify : Modifier un fait existant
type ModifyAction struct {
	network *ReteNetwork
}

func (m *ModifyAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	// args[0] = fait à modifier, args[1] = champs à modifier
	// Implementation...
}

// halt : Arrêter le moteur
type HaltAction struct {
	network *ReteNetwork
}

func (h *HaltAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	// Signaler arrêt du moteur
	// Implementation...
}
```

### 9.2 Enregistrement

```go
func (ae *ActionExecutor) RegisterDefaultActions() {
	// Actions actuelles
	ae.registry.Register(NewPrintAction(nil))
	
	// Actions proposées
	ae.registry.Register(NewAssertAction(ae.network))
	ae.registry.Register(NewRetractAction(ae.network))
	ae.registry.Register(NewModifyAction(ae.network))
	ae.registry.Register(NewHaltAction(ae.network))
	ae.registry.Register(NewLogAction(ae.logger))
}
```

---

## 10. Points d'Intervention pour Xuples

### 10.1 Conservation

✅ **Garder** :
- Interface `ActionHandler` (parfaite)
- `ActionRegistry` (très bien conçu)
- `ActionExecutor` (architecture solide)
- `ExecutionContext` (propre et efficace)
- Évaluation arguments (complète et robuste)

### 10.2 Propositions d'Amélioration

1. **Ajouter callback post-exécution** :
```go
type ActionCallback func(actionName string, args []interface{}, err error)

type ActionExecutor struct {
	// ... champs existants
	callbacks []ActionCallback
}

func (ae *ActionExecutor) AddCallback(cb ActionCallback) {
	ae.callbacks = append(ae.callbacks, cb)
}
```

**Usage pour xuples** :
```go
executor.AddCallback(func(actionName string, args []interface{}, err error) {
	// Notifier xuples de l'exécution
	xupleSpace.NotifyExecution(actionName, args, err)
})
```

2. **Mode asynchrone** :
```go
func (ae *ActionExecutor) ExecuteActionAsync(action *Action, token *Token) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- ae.ExecuteAction(action, token)
	}()
	return errChan
}
```

3. **Métriques d'exécution** :
```go
type ExecutionMetrics struct {
	ActionName    string
	ExecutionTime time.Duration
	Success       bool
	Error         error
}

func (ae *ActionExecutor) GetMetrics() []ExecutionMetrics {
	// Retourner historique des exécutions
}
```

---

## 11. Synthèse

### 11.1 Points Forts

✅ **Interface claire** : ActionHandler bien défini  
✅ **Registry flexible** : Enregistrement dynamique, thread-safe  
✅ **Évaluation robuste** : Support types complexes, messages d'erreur détaillés  
✅ **Thread-safe** : Aucun problème de concurrence identifié  
✅ **Panic recovery** : Robustesse garantie  
✅ **Extensible** : Facile d'ajouter nouvelles actions  
✅ **Testable** : ExecutionContext peut être mocké

### 11.2 Recommandations

1. **Conserver l'architecture actuelle** : Excellente conception
2. **Ajouter actions par défaut** : assert, retract, modify, halt, log
3. **Implémenter callbacks** : Pour intégration xuples
4. **Ajouter métriques** : Pour monitoring et debug
5. **Mode asynchrone optionnel** : Pour performances

---

## 12. Fichiers de Référence

| Fichier | Description | Lignes clés |
|---------|-------------|-------------|
| `rete/action_handler.go` | Interface et Registry | 12-25 (interface), 27-135 (registry) |
| `rete/action_executor.go` | Exécuteur principal | 37-42 (struct), 100-212 (execution) |
| `rete/action_executor_context.go` | Contexte d'exécution | 26-90 (ExecutionContext) |
| `rete/action_executor_evaluation.go` | Évaluation arguments | 11-300+ (evaluateArgument) |
| `rete/action_print.go` | Exemple de handler | 10-100 (PrintAction) |

---

**Date de création** : 2025-12-17  
**Auteur** : Analyse automatique pour refonte xuples  
**Statut** : ✅ Complet
