# 06 - Design de l'exécution immédiate des actions

## 🎯 Objectif

Modifier le comportement du moteur RETE pour qu'il exécute les actions **immédiatement** lors de l'arrivée d'un token déclencheur, sans stocker les tokens dans la mémoire des terminal nodes.

## 🎨 Nouveau flux d'exécution

### Diagramme de séquence cible

```
┌──────────┐     ┌──────────────┐     ┌───────────────┐     ┌────────────────┐
│  Network │     │ TerminalNode │     │ ActionExecutor│     │ActionObserver  │
└────┬─────┘     └──────┬───────┘     └───────┬───────┘     └───────┬────────┘
     │                  │                     │                     │
     │ ActivateLeft     │                     │                     │
     ├─────────────────>│                     │                     │
     │                  │                     │                     │
     │                  │ recordActivation()  │                     │
     │                  ├──────────┐          │                     │
     │                  │          │          │                     │
     │                  │<─────────┘          │                     │
     │                  │                     │                     │
     │                  │ executeAction()     │                     │
     │                  ├──────────┐          │                     │
     │                  │          │          │                     │
     │                  │<─────────┘          │                     │
     │                  │                     │                     │
     │                  │ ExecuteAction()     │                     │
     │                  ├────────────────────>│                     │
     │                  │                     │                     │
     │                  │                     │ Execute handlers    │
     │                  │                     ├──────────┐          │
     │                  │                     │          │          │
     │                  │                     │<─────────┘          │
     │                  │                     │                     │
     │                  │<────────────────────┤                     │
     │                  │                     │                     │
     │                  │ OnActionExecuted()  │                     │
     │                  ├─────────────────────────────────────────>│
     │                  │                     │                     │
     │                  │                     │ Capture/Process     │
     │                  │                     │<────────────────────┤
     │                  │                     │                     │
     │                  │<──────────────────────────────────────────┤
     │                  │                     │                     │
     │<─────────────────┤                     │                     │
     │                  │                     │                     │
```

### Nouveau process détaillé

1. **Token arrive** au TerminalNode via `ActivateLeft(token)`
2. **Enregistrement** de l'activation (métriques existantes)
3. **⚠️ PAS DE STOCKAGE** du token dans Memory.Tokens
4. **Exécution immédiate** via `executeAction(token)`
5. **Délégation** à `ActionExecutor.ExecuteAction(action, token)`
6. **Exécution** des handlers d'action (print, log, insert, etc.)
7. **Notification** de l'observer via `OnActionExecuted(result)`
8. **Gestion des erreurs** et logging
9. **Retour** au réseau RETE

## 🏗️ Architecture proposée

### 1. Interface ActionObserver (nouveau)

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "time"

// ActionObserver permet d'observer les exécutions d'actions.
// 
// Cette interface implémente le pattern Observer pour découpler
// le moteur RETE de la collecte/traitement des activations.
//
// Utilisations typiques :
//   - Tests : capturer les exécutions pour assertions
//   - Xuples : publier vers xuple-spaces
//   - Métriques : collecter statistiques d'exécution
//   - Logging : journalisation centralisée
//   - Audit : traçabilité des actions
//
// Thread-Safety :
//   - Les implémentations DOIVENT être thread-safe
//   - Plusieurs terminal nodes peuvent notifier en parallèle
//   - L'observer ne doit PAS bloquer l'exécution
type ActionObserver interface {
	// OnActionExecuted est appelé après chaque exécution d'action.
	//
	// L'appel se fait de manière synchrone après l'exécution.
	// L'implémentation NE DOIT PAS bloquer longtemps.
	// Pour traitements longs, utiliser une goroutine interne.
	//
	// Paramètres :
	//   - result : résultat de l'exécution avec contexte complet
	OnActionExecuted(result ExecutionResult)
}

// ExecutionResult représente le résultat de l'exécution d'une action.
type ExecutionResult struct {
	Success   bool              // true si l'exécution a réussi
	Error     error             // erreur si l'exécution a échoué
	Duration  time.Duration     // durée d'exécution
	Context   ActionContext     // contexte d'exécution complet
	Arguments []interface{}     // arguments évalués
}

// ActionContext contient le contexte d'exécution d'une action.
type ActionContext struct {
	ActionName string        // Nom de l'action (ex: "print", "insert")
	RuleName   string        // Nom de la règle qui a déclenché l'action
	Token      *Token        // Token déclencheur avec tous les faits
	Network    *ReteNetwork  // Réseau RETE (pour insert/update/retract)
	Timestamp  time.Time     // Moment de l'exécution
}

// NoOpObserver est un observateur qui ne fait rien.
// Utilisé comme valeur par défaut pour éviter les vérifications nil.
type NoOpObserver struct{}

// OnActionExecuted ne fait rien (implémentation vide).
func (n *NoOpObserver) OnActionExecuted(result ExecutionResult) {
	// Intentionally empty
}
```

### 2. Modification de TerminalNode

#### Nouveaux champs

```go
type TerminalNode struct {
	BaseNode
	Action *Action `json:"action"`
	
	// Observer pour notification des exécutions
	observer ActionObserver
	
	// Statistiques d'exécution (pour debug/tests)
	lastExecutionResult *ExecutionResult
	executionCount      int64
	statsMutex          sync.RWMutex
}
```

#### Nouvelle implémentation de ActivateLeft

```go
// ActivateLeft exécute immédiatement l'action sans stocker le token.
//
// Process :
//  1. Enregistre l'activation (métriques)
//  2. Exécute l'action immédiatement
//  3. Notifie l'observer du résultat
//  4. NE STOCKE PAS le token
//
// Thread-Safety :
//   - Méthode thread-safe
//   - Les statistiques sont protégées par statsMutex
//   - L'observer DOIT être thread-safe
//
// Paramètres :
//   - token : token contenant les faits et bindings déclencheurs
//
// Retourne :
//   - error : erreur si l'exécution de l'action échoue
func (tn *TerminalNode) ActivateLeft(token *Token) error {
	// Enregistrer l'activation (métriques réseau)
	tn.recordActivation()

	// PAS DE STOCKAGE - Exécuter directement
	start := time.Now()
	err := tn.executeAction(token)
	duration := time.Since(start)

	// Créer le résultat d'exécution
	result := ExecutionResult{
		Success:  err == nil,
		Error:    err,
		Duration: duration,
		Context: ActionContext{
			ActionName: tn.getActionName(),
			RuleName:   tn.getRuleName(),
			Token:      token,
			Network:    tn.BaseNode.GetNetwork(),
			Timestamp:  start,
		},
		Arguments: tn.extractArguments(token),
	}

	// Mettre à jour les statistiques (pour debug/tests)
	tn.updateStats(result)

	// Notifier l'observer
	if tn.observer != nil {
		tn.observer.OnActionExecuted(result)
	}

	return err
}
```

#### Méthodes helper

```go
// SetObserver configure l'observateur d'actions.
func (tn *TerminalNode) SetObserver(observer ActionObserver) {
	tn.observer = observer
}

// GetExecutionCount retourne le nombre total d'exécutions.
// Utilisé principalement pour les tests.
func (tn *TerminalNode) GetExecutionCount() int64 {
	tn.statsMutex.RLock()
	defer tn.statsMutex.RUnlock()
	return tn.executionCount
}

// GetLastExecutionResult retourne le dernier résultat d'exécution.
// Utilisé principalement pour les tests.
func (tn *TerminalNode) GetLastExecutionResult() *ExecutionResult {
	tn.statsMutex.RLock()
	defer tn.statsMutex.RUnlock()
	if tn.lastExecutionResult == nil {
		return nil
	}
	// Retourner une copie pour éviter modifications concurrentes
	resultCopy := *tn.lastExecutionResult
	return &resultCopy
}

// ResetExecutionStats réinitialise les statistiques d'exécution.
// Utilisé principalement pour les tests.
func (tn *TerminalNode) ResetExecutionStats() {
	tn.statsMutex.Lock()
	defer tn.statsMutex.Unlock()
	tn.lastExecutionResult = nil
	tn.executionCount = 0
}

// updateStats met à jour les statistiques internes.
func (tn *TerminalNode) updateStats(result ExecutionResult) {
	tn.statsMutex.Lock()
	defer tn.statsMutex.Unlock()
	
	// Copier le résultat
	resultCopy := result
	tn.lastExecutionResult = &resultCopy
	tn.executionCount++
}

// getActionName retourne le nom de l'action.
func (tn *TerminalNode) getActionName() string {
	if tn.Action == nil {
		return "unknown"
	}
	jobs := tn.Action.GetJobs()
	if len(jobs) > 0 {
		return jobs[0].Name
	}
	return "unknown"
}

// getRuleName extrait le nom de la règle depuis l'ID du nœud.
func (tn *TerminalNode) getRuleName() string {
	// L'ID du terminal node contient le nom de la règle
	// Format: "terminal_<ruleName>"
	if len(tn.ID) > 9 && tn.ID[:9] == "terminal_" {
		return tn.ID[9:]
	}
	return tn.ID
}

// extractArguments extrait les arguments évalués de l'action.
func (tn *TerminalNode) extractArguments(token *Token) []interface{} {
	if tn.Action == nil {
		return nil
	}
	
	// Utiliser l'ActionExecutor pour évaluer les arguments
	network := tn.BaseNode.GetNetwork()
	if network == nil || network.ActionExecutor == nil {
		return nil
	}
	
	// Les arguments sont évalués dans ExecuteAction
	// Ici on retourne juste les arguments bruts
	jobs := tn.Action.GetJobs()
	if len(jobs) == 0 {
		return nil
	}
	
	return jobs[0].Args
}
```

### 3. Modification de ReteNetwork

#### Nouveaux champs et méthodes

```go
// Dans la structure ReteNetwork
type ReteNetwork struct {
	// ... champs existants ...
	
	ActionExecutor *ActionExecutor  `json:"-"` // Déjà existant
	actionObserver ActionObserver   `json:"-"` // NOUVEAU
}

// SetActionObserver configure l'observateur pour tous les terminal nodes.
//
// Cette méthode configure l'observer pour tous les terminal nodes
// existants ET futurs (via AddTerminalNode).
//
// Thread-Safety :
//   - Méthode thread-safe si appelée avant démarrage du réseau
//   - Si appelée pendant l'exécution, risque de race condition
//   - Recommandé : appeler pendant la phase d'initialisation
//
// Paramètres :
//   - observer : observateur à configurer (peut être nil pour désactiver)
func (rn *ReteNetwork) SetActionObserver(observer ActionObserver) {
	if observer == nil {
		observer = &NoOpObserver{}
	}
	
	rn.actionObserver = observer
	
	// Configurer tous les terminal nodes existants
	for _, terminal := range rn.TerminalNodes {
		terminal.SetObserver(observer)
	}
}

// GetActionObserver retourne l'observateur configuré.
func (rn *ReteNetwork) GetActionObserver() ActionObserver {
	if rn.actionObserver == nil {
		return &NoOpObserver{}
	}
	return rn.actionObserver
}
```

#### Modification de AddTerminalNode (si existe)

Si une méthode `AddTerminalNode` existe, elle doit configurer automatiquement l'observer :

```go
func (rn *ReteNetwork) AddTerminalNode(node *TerminalNode) {
	// Ajouter le nœud
	rn.TerminalNodes[node.ID] = node
	
	// Configurer automatiquement l'observer
	if rn.actionObserver != nil {
		node.SetObserver(rn.actionObserver)
	}
}
```

### 4. Collecteur d'activations (remplacement de collectActivations)

#### ExecutionStatsCollector (nouveau)

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"sync"
	
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tsdio"
)

// ExecutionStatsCollector collecte les statistiques d'exécution des actions.
//
// Implémente ActionObserver pour capturer toutes les exécutions
// et les convertir en format tsdio.Activation pour l'API.
//
// Thread-Safety :
//   - Thread-safe grâce au mutex interne
//   - Peut être utilisé par plusieurs terminal nodes en parallèle
type ExecutionStatsCollector struct {
	executions []rete.ExecutionResult
	mu         sync.RWMutex
}

// NewExecutionStatsCollector crée un nouveau collecteur.
func NewExecutionStatsCollector() *ExecutionStatsCollector {
	return &ExecutionStatsCollector{
		executions: make([]rete.ExecutionResult, 0),
	}
}

// OnActionExecuted capture un résultat d'exécution.
// Implémente ActionObserver.
func (c *ExecutionStatsCollector) OnActionExecuted(result rete.ExecutionResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executions = append(c.executions, result)
}

// GetExecutions retourne une copie de tous les résultats capturés.
func (c *ExecutionStatsCollector) GetExecutions() []rete.ExecutionResult {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// Retourner une copie pour éviter modifications concurrentes
	return append([]rete.ExecutionResult{}, c.executions...)
}

// GetActivations convertit les résultats en format tsdio.Activation.
func (c *ExecutionStatsCollector) GetActivations() []tsdio.Activation {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	activations := make([]tsdio.Activation, 0, len(c.executions))
	
	for _, exec := range c.executions {
		activation := tsdio.Activation{
			ActionName:      exec.Context.ActionName,
			Arguments:       formatArguments(exec.Arguments),
			TriggeringFacts: extractFacts(exec.Context.Token),
			BindingsCount:   len(exec.Context.Token.Facts),
			Success:         exec.Success,
			Duration:        exec.Duration,
			Error:           formatError(exec.Error),
		}
		activations = append(activations, activation)
	}
	
	return activations
}

// GetExecutionCount retourne le nombre d'exécutions capturées.
func (c *ExecutionStatsCollector) GetExecutionCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.executions)
}

// Reset réinitialise le collecteur.
func (c *ExecutionStatsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executions = make([]rete.ExecutionResult, 0)
}

// Helper functions

func formatArguments(args []interface{}) []interface{} {
	// Conversion si nécessaire
	return args
}

func extractFacts(token *rete.Token) []string {
	if token == nil {
		return []string{}
	}
	
	facts := make([]string, 0, len(token.Facts))
	for _, fact := range token.Facts {
		facts = append(facts, fact.GetInternalID())
	}
	return facts
}

func formatError(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
```

## 🧪 Stratégie de migration des tests

### 1. Tests utilisant Memory.Tokens

**Avant** :
```go
// Vérifier les tokens stockés
if len(terminal.Memory.Tokens) != 1 {
    t.Errorf("Expected 1 activation, got %d", len(terminal.Memory.Tokens))
}
```

**Après** :
```go
// Vérifier via les statistiques
if terminal.GetExecutionCount() != 1 {
    t.Errorf("Expected 1 execution, got %d", terminal.GetExecutionCount())
}

// Ou via un observer de test
observer := NewTestActionObserver(t)
network.SetActionObserver(observer)
// ... exécution ...
observer.AssertExecutionCount(1)
```

### 2. TestActionObserver pour tests

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete_test

import (
	"sync"
	"testing"
	
	"github.com/treivax/tsd/rete"
)

// TestActionObserver capture les exécutions pour assertions dans tests.
type TestActionObserver struct {
	t          *testing.T
	executions []rete.ExecutionResult
	mu         sync.RWMutex
}

// NewTestActionObserver crée un observateur de test.
func NewTestActionObserver(t *testing.T) *TestActionObserver {
	return &TestActionObserver{
		t:          t,
		executions: make([]rete.ExecutionResult, 0),
	}
}

// OnActionExecuted capture l'exécution.
func (o *TestActionObserver) OnActionExecuted(result rete.ExecutionResult) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	o.t.Logf("🎯 Action executed: %s (rule: %s, success: %v, duration: %v)",
		result.Context.ActionName,
		result.Context.RuleName,
		result.Success,
		result.Duration)
	
	o.executions = append(o.executions, result)
}

// GetExecutions retourne tous les résultats capturés.
func (o *TestActionObserver) GetExecutions() []rete.ExecutionResult {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return append([]rete.ExecutionResult{}, o.executions...)
}

// GetExecutionCount retourne le nombre d'exécutions.
func (o *TestActionObserver) GetExecutionCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.executions)
}

// AssertExecutionCount vérifie le nombre d'exécutions.
func (o *TestActionObserver) AssertExecutionCount(expected int) {
	o.mu.RLock()
	count := len(o.executions)
	o.mu.RUnlock()
	
	if count != expected {
		o.t.Errorf("❌ Expected %d executions, got %d", expected, count)
	}
}

// AssertActionExecuted vérifie qu'une action a été exécutée.
func (o *TestActionObserver) AssertActionExecuted(actionName string) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	for _, exec := range o.executions {
		if exec.Context.ActionName == actionName {
			return
		}
	}
	o.t.Errorf("❌ Action '%s' was not executed", actionName)
}

// AssertAllSuccessful vérifie que toutes les exécutions ont réussi.
func (o *TestActionObserver) AssertAllSuccessful() {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	for i, exec := range o.executions {
		if !exec.Success {
			o.t.Errorf("❌ Execution %d failed: %v", i, exec.Error)
		}
	}
}

// Reset réinitialise le collecteur.
func (o *TestActionObserver) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.executions = make([]rete.ExecutionResult, 0)
}
```

## 📊 Impact et bénéfices

### Performance

| Aspect | Avant | Après | Gain |
|--------|-------|-------|------|
| Stockage tokens | O(n) en mémoire | O(1) stats | -100% mémoire |
| collectActivations | O(n*m) parcours | O(1) accès | -100% CPU |
| Notification | Polling | Push immédiat | Temps réel |

### Architecture

| Aspect | Avant | Après |
|--------|-------|-------|
| Couplage | Servercmd ↔ TerminalNode | Observer pattern |
| Encapsulation | Violation (Memory.Tokens) | Respectée |
| Séparation | RETE + collecte mélangés | RETE pur |
| Extensibilité | Difficile | Observer chain |

### Testabilité

| Aspect | Avant | Après |
|--------|-------|-------|
| Vérification | Indirect via Memory | Direct via observer |
| Assertions | Complexes | Simples et claires |
| Isolation | Couplage fort | Découplage total |

## ✅ Checklist d'implémentation

- [ ] Créer `action_observer.go` avec interfaces
- [ ] Modifier `node_terminal.go` pour exécution immédiate
- [ ] Ajouter champs observer et stats à TerminalNode
- [ ] Modifier `ActivateLeft` pour ne plus stocker
- [ ] Ajouter méthodes SetObserver, GetExecutionCount, etc.
- [ ] Modifier `network.go` pour configurer observers
- [ ] Créer `ExecutionStatsCollector` dans servercmd
- [ ] Remplacer `collectActivations` par observer
- [ ] Créer `TestActionObserver` pour tests
- [ ] Migrer tous les tests utilisant Memory.Tokens
- [ ] Supprimer ou déprécier GetTriggeredActions()
- [ ] Mettre à jour la documentation
- [ ] Tests de non-régression
- [ ] Tests de performance
- [ ] Validation avec `make test-complete`

## 🚀 Plan de déploiement

### Phase 1 : Implémentation infrastructure
1. Créer interfaces et types
2. Modifier TerminalNode
3. Modifier Network
4. Tests unitaires

### Phase 2 : Migration du serveur
1. Créer ExecutionStatsCollector
2. Modifier executeTSDProgram
3. Supprimer collectActivations
4. Tests d'intégration

### Phase 3 : Migration des tests
1. Créer TestActionObserver
2. Migrer tests un par un
3. Vérifier couverture
4. Supprimer code obsolète

### Phase 4 : Validation
1. make test-complete
2. Tests de performance
3. Documentation
4. Review finale

## 📚 Références

- Observer Pattern : https://refactoring.guru/design-patterns/observer
- Go Concurrency Patterns : https://go.dev/blog/pipelines
- RETE Algorithm : Forgy, C. (1982)
- Common.md : Standards du projet
