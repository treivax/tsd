# Prompt 05 - Modification du moteur RETE pour exécution immédiate des actions

## 🎯 Objectif

Modifier le comportement du moteur RETE pour qu'il exécute les actions immédiatement lorsqu'un token déclencheur est produit, au lieu de stocker les tokens dans la mémoire des terminal nodes.

Cette modification est essentielle pour :
- Rétablir le comportement classique d'un moteur de règles RETE
- Séparer clairement RETE (moteur de règles) et xuples (système de coordination)
- Permettre aux actions de s'exécuter en temps réel

## 📋 Tâches

### 1. Analyser le comportement actuel des Terminal Nodes

**Objectif** : Comprendre précisément comment les terminal nodes stockent actuellement les tokens.

- [ ] Examiner `tsd/rete/node_terminal.go` en détail
- [ ] Identifier la méthode `ActivateLeft` qui stocke les tokens
- [ ] Analyser la méthode `executeAction` et son rôle actuel
- [ ] Comprendre comment `collectActivations` récupère les tokens stockés
- [ ] Identifier tous les usages de la mémoire des terminal nodes

**Livrables** :
- Créer `tsd/docs/xuples/implementation/05-terminal-node-current-behavior.md` documentant :
  - Flux actuel : token → stockage → récupération
  - Méthodes impliquées et leur responsabilité
  - Structure de la mémoire (WorkingMemory)
  - Points d'utilisation dans le code
  - Diagramme de séquence actuel

### 2. Concevoir le nouveau comportement d'exécution immédiate

**Objectif** : Définir comment les actions seront exécutées immédiatement.

**Nouveau flux attendu** :
1. Un token déclencheur arrive au terminal node
2. Le terminal node extrait les informations nécessaires du token
3. Le terminal node invoque l'ActionExecutor avec ces informations
4. L'action s'exécute immédiatement (Print, Log, Update, Insert, Retract, Xuple)
5. Les erreurs d'exécution sont propagées et loggées
6. **Le token n'est PAS stocké** (sauf pour debug/observabilité optionnelle)

**Considérations importantes** :
- Gestion des erreurs d'exécution
- Performance (éviter allocations inutiles)
- Observabilité (tracer les actions exécutées)
- Testabilité (injection de l'executor)
- Compatibilité avec les tests existants

**Livrables** :
- Créer `tsd/docs/xuples/implementation/06-immediate-execution-design.md` contenant :
  - Nouveau flux d'exécution détaillé
  - Diagramme de séquence du nouveau comportement
  - Interface ActionExecutor requise
  - Gestion des erreurs
  - Stratégie d'observabilité
  - Plan de migration des tests existants
  - Mesures de performance

### 3. Définir l'interface ActionExecutor

**Objectif** : Créer une interface claire pour l'exécution des actions.

- [ ] Définir le contrat d'exécution
- [ ] Prévoir la transmission du contexte (token, réseau)
- [ ] Gérer les erreurs
- [ ] Permettre l'injection pour tests

**Fichier à créer** :
- `tsd/rete/action_executor.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// ActionExecutor définit l'interface pour exécuter des actions
type ActionExecutor interface {
    // Execute exécute une action avec les arguments fournis et le token déclencheur
    // Paramètres:
    //   - actionName: nom de l'action à exécuter
    //   - args: arguments évalués de l'action
    //   - token: token déclencheur contenant les faits
    // Retour:
    //   - error si l'exécution échoue
    Execute(actionName string, args []interface{}, token *Token) error
}

// ActionContext contient le contexte d'exécution d'une action
type ActionContext struct {
    ActionName string        // Nom de l'action
    RuleName   string        // Nom de la règle qui a déclenché l'action
    Token      *Token        // Token déclencheur
    Network    *Network      // Réseau RETE (pour Update, Insert, Retract)
    Timestamp  time.Time     // Moment de l'exécution
}

// ExecutionResult représente le résultat de l'exécution d'une action
type ExecutionResult struct {
    Success   bool
    Error     error
    Duration  time.Duration
    Context   ActionContext
}

// ActionObserver permet d'observer les exécutions d'actions (pour debug/tests)
type ActionObserver interface {
    OnActionExecuted(result ExecutionResult)
}

// NoOpObserver est un observateur qui ne fait rien (défaut)
type NoOpObserver struct{}

func (n *NoOpObserver) OnActionExecuted(result ExecutionResult) {
    // Ne fait rien
}
```

**Livrables** :
- [ ] Interface ActionExecutor définie avec copyright
- [ ] Structures de contexte et résultat
- [ ] Interface ActionObserver pour observabilité
- [ ] Documentation GoDoc complète
- [ ] Exemples d'utilisation en commentaire

### 4. Modifier TerminalNode pour exécution immédiate

**Objectif** : Transformer le terminal node pour exécuter immédiatement au lieu de stocker.

- [ ] Modifier la méthode `ActivateLeft`
- [ ] Supprimer le stockage dans la mémoire (ou le rendre optionnel)
- [ ] Invoquer l'ActionExecutor immédiatement
- [ ] Gérer les erreurs d'exécution
- [ ] Ajouter l'observabilité
- [ ] Conserver les informations nécessaires pour les tests

**Fichier à modifier** :
- `tsd/rete/node_terminal.go`

**Modification attendue** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "fmt"
    "log"
    "time"
)

// TerminalNode représente un nœud terminal du réseau RETE
type TerminalNode struct {
    BaseNode
    RuleName       string
    ActionName     string
    ActionArgs     []ActionArgument
    executor       ActionExecutor
    observer       ActionObserver
    network        *Network
    
    // Pour debug/tests uniquement (optionnel)
    lastExecutionResult *ExecutionResult
    executionCount      int
}

// NewTerminalNode crée un nouveau terminal node
func NewTerminalNode(ruleName, actionName string, args []ActionArgument) *TerminalNode {
    return &TerminalNode{
        BaseNode:   BaseNode{},
        RuleName:   ruleName,
        ActionName: actionName,
        ActionArgs: args,
        observer:   &NoOpObserver{}, // Défaut
    }
}

// SetExecutor configure l'exécuteur d'actions
func (tn *TerminalNode) SetExecutor(executor ActionExecutor) {
    tn.executor = executor
}

// SetObserver configure l'observateur d'actions
func (tn *TerminalNode) SetObserver(observer ActionObserver) {
    tn.observer = observer
}

// SetNetwork configure le réseau RETE (pour le contexte)
func (tn *TerminalNode) SetNetwork(network *Network) {
    tn.network = network
}

// ActivateLeft est appelé quand un token arrive au terminal node
func (tn *TerminalNode) ActivateLeft(token *Token, binding *Binding) {
    // NE PLUS STOCKER LE TOKEN - EXÉCUTER IMMÉDIATEMENT
    
    if tn.executor == nil {
        log.Printf("⚠️  Terminal node '%s' has no executor, skipping action execution", tn.RuleName)
        return
    }
    
    // Évaluer les arguments de l'action à partir du token et du binding
    args, err := tn.evaluateActionArguments(token, binding)
    if err != nil {
        log.Printf("❌ Failed to evaluate action arguments for '%s': %v", tn.RuleName, err)
        return
    }
    
    // Créer le contexte d'exécution
    ctx := ActionContext{
        ActionName: tn.ActionName,
        RuleName:   tn.RuleName,
        Token:      token,
        Network:    tn.network,
        Timestamp:  time.Now(),
    }
    
    // Exécuter l'action IMMÉDIATEMENT
    start := time.Now()
    err = tn.executor.Execute(tn.ActionName, args, token)
    duration := time.Since(start)
    
    // Créer le résultat d'exécution
    result := ExecutionResult{
        Success:  err == nil,
        Error:    err,
        Duration: duration,
        Context:  ctx,
    }
    
    // Notifier l'observateur
    tn.observer.OnActionExecuted(result)
    
    // Logger les erreurs
    if err != nil {
        log.Printf("❌ Action '%s' failed in rule '%s': %v", tn.ActionName, tn.RuleName, err)
    } else {
        log.Printf("✅ Action '%s' executed successfully in rule '%s' (took %v)", 
            tn.ActionName, tn.RuleName, duration)
    }
    
    // Pour debug/tests : conserver le dernier résultat
    tn.lastExecutionResult = &result
    tn.executionCount++
}

// evaluateActionArguments évalue les arguments de l'action
func (tn *TerminalNode) evaluateActionArguments(token *Token, binding *Binding) ([]interface{}, error) {
    args := make([]interface{}, len(tn.ActionArgs))
    
    for i, arg := range tn.ActionArgs {
        value, err := arg.Evaluate(token, binding)
        if err != nil {
            return nil, fmt.Errorf("failed to evaluate argument %d: %w", i, err)
        }
        args[i] = value
    }
    
    return args, nil
}

// GetLastExecutionResult retourne le dernier résultat d'exécution (pour tests)
func (tn *TerminalNode) GetLastExecutionResult() *ExecutionResult {
    return tn.lastExecutionResult
}

// GetExecutionCount retourne le nombre d'exécutions (pour tests)
func (tn *TerminalNode) GetExecutionCount() int {
    return tn.executionCount
}

// ResetExecutionStats réinitialise les statistiques (pour tests)
func (tn *TerminalNode) ResetExecutionStats() {
    tn.lastExecutionResult = nil
    tn.executionCount = 0
}
```

**Livrables** :
- [ ] TerminalNode modifié pour exécution immédiate
- [ ] Stockage supprimé (ou rendu optionnel pour debug)
- [ ] Gestion d'erreurs robuste
- [ ] Logging des exécutions
- [ ] Méthodes pour tests (GetLastExecutionResult, etc.)
- [ ] Documentation mise à jour

### 5. Adapter Network pour configurer les executors

**Objectif** : Permettre au réseau RETE de configurer les executors des terminal nodes.

- [ ] Ajouter un champ executor au Network
- [ ] Configurer automatiquement tous les terminal nodes
- [ ] Permettre l'injection d'un executor personnalisé (tests)

**Fichier à modifier** :
- `tsd/rete/network.go`

**Modification attendue** :
```go
// Dans la structure Network
type Network struct {
    // ... champs existants ...
    
    executor ActionExecutor
    observer ActionObserver
}

// SetActionExecutor configure l'exécuteur d'actions pour tous les terminal nodes
func (n *Network) SetActionExecutor(executor ActionExecutor) {
    n.executor = executor
    
    // Configurer tous les terminal nodes existants
    for _, terminalNode := range n.TerminalNodes {
        terminalNode.SetExecutor(executor)
        terminalNode.SetNetwork(n)
    }
}

// SetActionObserver configure l'observateur pour tous les terminal nodes
func (n *Network) SetActionObserver(observer ActionObserver) {
    n.observer = observer
    
    for _, terminalNode := range n.TerminalNodes {
        terminalNode.SetObserver(observer)
    }
}

// AddTerminalNode ajoute un terminal node et le configure
func (n *Network) AddTerminalNode(node *TerminalNode) {
    n.TerminalNodes = append(n.TerminalNodes, node)
    
    // Configurer automatiquement avec l'executor du réseau
    if n.executor != nil {
        node.SetExecutor(n.executor)
        node.SetNetwork(n)
    }
    if n.observer != nil {
        node.SetObserver(n.observer)
    }
}
```

**Livrables** :
- [ ] Network modifié avec executor et observer
- [ ] Configuration automatique des terminal nodes
- [ ] Méthodes d'injection pour tests
- [ ] Documentation mise à jour

### 6. Supprimer ou adapter collectActivations

**Objectif** : Adapter ou supprimer la fonction `collectActivations` qui récupérait les tokens stockés.

**Options** :
1. **Supprimer complètement** si plus nécessaire
2. **Adapter** pour retourner les statistiques d'exécution
3. **Conserver** pour compatibilité avec observer pattern

**Approche recommandée** : Adapter pour retourner les statistiques via observer.

**Fichier à modifier** :
- `tsd/internal/servercmd/servercmd.go`

**Code attendu** :
```go
// Remplacer collectActivations par collectExecutionStats

// ExecutionStatsCollector collecte les statistiques d'exécution
type ExecutionStatsCollector struct {
    executions []rete.ExecutionResult
    mu         sync.Mutex
}

func NewExecutionStatsCollector() *ExecutionStatsCollector {
    return &ExecutionStatsCollector{
        executions: make([]rete.ExecutionResult, 0),
    }
}

func (c *ExecutionStatsCollector) OnActionExecuted(result rete.ExecutionResult) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.executions = append(c.executions, result)
}

func (c *ExecutionStatsCollector) GetExecutions() []rete.ExecutionResult {
    c.mu.Lock()
    defer c.mu.Unlock()
    return append([]rete.ExecutionResult{}, c.executions...)
}

// Dans executeTSDProgram
func executeTSDProgram(source string) (*tsdio.ExecuteResponse, error) {
    // ... parsing et compilation ...
    
    // Créer un collecteur de statistiques
    statsCollector := NewExecutionStatsCollector()
    network.SetActionObserver(statsCollector)
    
    // Créer et configurer l'executor
    executor := actions.NewBuiltinActionExecutor(network, xupleManager)
    network.SetActionExecutor(executor)
    
    // Insérer les faits (déclenche les règles et actions)
    for _, fact := range facts {
        network.InsertFact(fact)
    }
    
    // Récupérer les statistiques d'exécution
    executions := statsCollector.GetExecutions()
    
    // Retourner la réponse avec les statistiques
    return &tsdio.ExecuteResponse{
        Success:        true,
        ExecutionStats: convertExecutionStats(executions),
        // ...
    }, nil
}
```

**Livrables** :
- [ ] collectActivations supprimé ou adapté
- [ ] ExecutionStatsCollector implémenté
- [ ] Intégration dans executeTSDProgram
- [ ] Tests mis à jour

### 7. Migrer les tests existants

**Objectif** : Adapter tous les tests qui utilisaient la mémoire des terminal nodes.

- [ ] Identifier tous les tests utilisant `terminal.Memory.Tokens`
- [ ] Remplacer par des vérifications via observer ou méthodes de test
- [ ] Utiliser `GetLastExecutionResult()` et `GetExecutionCount()`
- [ ] Créer des helpers de test pour simplifier

**Helper de test attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete_test

import "testing"

// TestActionObserver pour capturer les exécutions dans les tests
type TestActionObserver struct {
    t          *testing.T
    executions []rete.ExecutionResult
}

func NewTestActionObserver(t *testing.T) *TestActionObserver {
    return &TestActionObserver{
        t:          t,
        executions: make([]rete.ExecutionResult, 0),
    }
}

func (o *TestActionObserver) OnActionExecuted(result rete.ExecutionResult) {
    o.t.Logf("🎯 Action executed: %s (rule: %s, success: %v)",
        result.Context.ActionName, result.Context.RuleName, result.Success)
    o.executions = append(o.executions, result)
}

func (o *TestActionObserver) GetExecutions() []rete.ExecutionResult {
    return o.executions
}

func (o *TestActionObserver) GetExecutionCount() int {
    return len(o.executions)
}

func (o *TestActionObserver) AssertExecutionCount(expected int) {
    if len(o.executions) != expected {
        o.t.Errorf("❌ Expected %d executions, got %d", expected, len(o.executions))
    }
}

func (o *TestActionObserver) AssertActionExecuted(actionName string) {
    for _, exec := range o.executions {
        if exec.Context.ActionName == actionName {
            return
        }
    }
    o.t.Errorf("❌ Action '%s' was not executed", actionName)
}
```

**Exemple de migration de test** :
```go
// AVANT (ancien comportement)
func TestRule_Activation(t *testing.T) {
    // ... setup ...
    
    // Vérifier les tokens stockés
    if len(terminal.Memory.Tokens) != 1 {
        t.Errorf("Expected 1 activation, got %d", len(terminal.Memory.Tokens))
    }
}

// APRÈS (nouveau comportement)
func TestRule_Execution(t *testing.T) {
    // ... setup ...
    
    // Créer un observateur de test
    observer := NewTestActionObserver(t)
    network.SetActionObserver(observer)
    
    // Créer un executor de test (ou mock)
    executor := NewTestActionExecutor(t)
    network.SetActionExecutor(executor)
    
    // ... exécution ...
    
    // Vérifier les exécutions
    observer.AssertExecutionCount(1)
    observer.AssertActionExecuted("Print")
    
    // Ou via le terminal node directement
    if terminal.GetExecutionCount() != 1 {
        t.Errorf("Expected 1 execution, got %d", terminal.GetExecutionCount())
    }
}
```

**Livrables** :
- [ ] Helper TestActionObserver créé
- [ ] Helper TestActionExecutor créé (mock)
- [ ] Tous les tests migrés
- [ ] Tous les tests passent
- [ ] Documentation des nouveaux helpers

### 8. Créer les tests du nouveau comportement

**Objectif** : Tester exhaustivement l'exécution immédiate.

**Fichier à créer** :
- `tsd/rete/terminal_node_execution_test.go`

**Tests attendus** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "testing"

func TestTerminalNode_ImmediateExecution(t *testing.T) {
    t.Log("🧪 TEST EXÉCUTION IMMÉDIATE DES ACTIONS")
    
    // Créer un terminal node
    terminal := NewTerminalNode("test-rule", "Print", []ActionArgument{
        {Type: "string", Value: "Hello"},
    })
    
    // Créer un mock executor
    mockExecutor := &MockActionExecutor{
        executed: make([]ExecutionRecord, 0),
    }
    terminal.SetExecutor(mockExecutor)
    
    // Créer un token
    token := NewToken(nil, NewFact("test", nil))
    
    // Activer le terminal node
    terminal.ActivateLeft(token, NewBinding())
    
    // Vérifier que l'action a été exécutée IMMÉDIATEMENT
    if len(mockExecutor.executed) != 1 {
        t.Fatalf("❌ Expected 1 execution, got %d", len(mockExecutor.executed))
    }
    
    exec := mockExecutor.executed[0]
    if exec.ActionName != "Print" {
        t.Errorf("❌ Expected action 'Print', got '%s'", exec.ActionName)
    }
    
    // Vérifier via le terminal node
    if terminal.GetExecutionCount() != 1 {
        t.Errorf("❌ Expected execution count 1, got %d", terminal.GetExecutionCount())
    }
    
    t.Log("✅ Action exécutée immédiatement")
}

func TestTerminalNode_ExecutionError(t *testing.T) {
    t.Log("🧪 TEST GESTION ERREUR EXÉCUTION")
    
    terminal := NewTerminalNode("test-rule", "FailingAction", nil)
    
    // Executor qui échoue
    mockExecutor := &MockActionExecutor{
        shouldFail: true,
    }
    terminal.SetExecutor(mockExecutor)
    
    token := NewToken(nil, NewFact("test", nil))
    
    // Activer (ne devrait pas paniquer)
    terminal.ActivateLeft(token, NewBinding())
    
    // Vérifier que l'erreur a été capturée
    result := terminal.GetLastExecutionResult()
    if result == nil {
        t.Fatal("❌ Expected execution result")
    }
    
    if result.Success {
        t.Error("❌ Expected execution to fail")
    }
    
    if result.Error == nil {
        t.Error("❌ Expected error to be set")
    }
    
    t.Log("✅ Erreur d'exécution correctement gérée")
}

func TestTerminalNode_Observer(t *testing.T) {
    t.Log("🧪 TEST OBSERVATEUR D'ACTIONS")
    
    terminal := NewTerminalNode("test-rule", "Print", nil)
    
    executor := &MockActionExecutor{}
    observer := NewTestActionObserver(t)
    
    terminal.SetExecutor(executor)
    terminal.SetObserver(observer)
    
    token := NewToken(nil, NewFact("test", nil))
    
    // Exécuter plusieurs fois
    terminal.ActivateLeft(token, NewBinding())
    terminal.ActivateLeft(token, NewBinding())
    terminal.ActivateLeft(token, NewBinding())
    
    // Vérifier que l'observateur a reçu toutes les notifications
    if observer.GetExecutionCount() != 3 {
        t.Errorf("❌ Expected 3 observations, got %d", observer.GetExecutionCount())
    }
    
    t.Log("✅ Observateur correctement notifié")
}

// MockActionExecutor pour les tests
type MockActionExecutor struct {
    executed   []ExecutionRecord
    shouldFail bool
}

type ExecutionRecord struct {
    ActionName string
    Args       []interface{}
    Token      *Token
}

func (m *MockActionExecutor) Execute(actionName string, args []interface{}, token *Token) error {
    m.executed = append(m.executed, ExecutionRecord{
        ActionName: actionName,
        Args:       args,
        Token:      token,
    })
    
    if m.shouldFail {
        return fmt.Errorf("mock execution failure")
    }
    
    return nil
}
```

**Livrables** :
- [ ] Tests d'exécution immédiate
- [ ] Tests de gestion d'erreurs
- [ ] Tests de l'observer pattern
- [ ] Tests de performance (pas de regression)
- [ ] Couverture > 80%
- [ ] Tous les tests passent

## 📁 Structure attendue

```
tsd/
├── docs/xuples/implementation/
│   ├── 05-terminal-node-current-behavior.md
│   └── 06-immediate-execution-design.md
├── rete/
│   ├── action_executor.go              # Nouveau
│   ├── node_terminal.go                # Modifié
│   ├── network.go                      # Modifié
│   ├── terminal_node_execution_test.go # Nouveau
│   └── test_helpers.go                 # Nouveau (helpers de test)
└── internal/servercmd/
    └── servercmd.go                    # Modifié (collectActivations)
```

## ✅ Critères de succès

- [ ] Terminal nodes exécutent immédiatement au lieu de stocker
- [ ] Interface ActionExecutor définie et implémentée
- [ ] Observer pattern implémenté pour observabilité
- [ ] Gestion d'erreurs robuste
- [ ] Logging des exécutions
- [ ] Méthodes de test (GetLastExecutionResult, etc.)
- [ ] collectActivations adapté ou supprimé
- [ ] Tous les tests existants migrés et passent
- [ ] Nouveaux tests complets avec couverture > 80%
- [ ] Aucune régression de performance
- [ ] `make test-complete` passe sans erreur
- [ ] Documentation complète

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/design/` - Conception du module
- `tsd/docs/xuples/implementation/` - Documentation d'implémentation
- Observer Pattern - Design Patterns
- Effective Go - https://go.dev/doc/effective_go

## 🎯 Prochaine étape

Une fois le moteur RETE modifié pour l'exécution immédiate, passer au prompt **06-implement-xuples-module.md** pour implémenter le module xuples complet avec les xuple-spaces et leurs politiques.