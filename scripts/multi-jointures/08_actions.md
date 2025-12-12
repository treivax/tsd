# Prompt 08 : ExecutionContext et Actions

**Session** : 8/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Prompt 07 complété, BetaChainBuilder validé

---

## 🎯 Objectif de cette Session

Adapter l'exécution des actions pour utiliser BindingChain :
1. ExecutionContext utilise *BindingChain au lieu de map
2. Résolution de variables via BindingChain
3. Messages d'erreur clairs listant les variables disponibles
4. TerminalNode propage correctement les bindings aux actions

**Livrables** :
- `tsd/rete/action_executor_context.go` (modifié)
- `tsd/rete/action_executor_evaluation.go` (modifié)
- `tsd/rete/node_terminal.go` (modifié)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Modifier ExecutionContext (30 min)

#### 1.1 Mettre à jour la structure

**Fichier** : `tsd/rete/action_executor_context.go`

**Chercher** :
```go
type ExecutionContext struct {
    varCache map[string]*Fact  // ❌ ANCIEN
    // ... autres champs
}
```

**Remplacer par** :
```go
type ExecutionContext struct {
    bindings *BindingChain     // ✅ NOUVEAU
    // ... autres champs conservés
}
```

#### 1.2 Adapter le constructeur

```go
func NewExecutionContext(token *Token, network *Network, ...) *ExecutionContext {
    return &ExecutionContext{
        bindings: token.Bindings,  // Référence à la chaîne immuable
        network:  network,
        // ... autres champs
    }
}
```

---

### Tâche 2 : Adapter la Résolution de Variables (40 min)

#### 2.1 Mettre à jour resolveVariable

**Fichier** : `tsd/rete/action_executor_evaluation.go`

**Remplacer** :
```go
func (ctx *ExecutionContext) resolveVariable(name string) (interface{}, error) {
    // ANCIEN
    if fact, ok := ctx.varCache[name]; ok {
        return fact, nil
    }
    return nil, fmt.Errorf("variable '%s' non trouvée", name)
}
```

**Par** :
```go
func (ctx *ExecutionContext) resolveVariable(name string) (interface{}, error) {
    // NOUVEAU : Utiliser BindingChain
    if ctx.bindings != nil && ctx.bindings.Has(name) {
        fact := ctx.bindings.Get(name)
        return fact, nil
    }
    
    // Message d'erreur TRÈS clair avec liste des variables disponibles
    available := []string{}
    if ctx.bindings != nil {
        available = ctx.bindings.Variables()
    }
    
    return nil, fmt.Errorf(
        "❌ Erreur d'exécution d'action:\n"+
        "   Variable '%s' non trouvée dans le contexte\n"+
        "   Variables disponibles: %v\n"+
        "   Vérifiez que la règle déclare bien cette variable dans sa clause de pattern",
        name, available,
    )
}
```

**Points clés** :
- Utiliser `bindings.Has()` et `bindings.Get()`
- Message d'erreur détaillé avec liste complète des variables
- Aide l'utilisateur à comprendre le problème

---

#### 2.2 Adapter evaluateArgument

**Chercher toutes les fonctions qui utilisent varCache** :

```bash
grep -n "varCache" rete/action_executor_evaluation.go
```

**Pour chaque occurrence, remplacer** :
```go
// ANCIEN
fact := ctx.varCache[varName]

// NOUVEAU
fact := ctx.bindings.Get(varName)
```

**Exemple dans evaluateArgument** :
```go
func (ctx *ExecutionContext) evaluateArgument(arg interface{}) (interface{}, error) {
    switch v := arg.(type) {
    case map[string]interface{}:
        argType, _ := v["type"].(string)
        
        switch argType {
        case "variable":
            varName, _ := v["name"].(string)
            // Utiliser resolveVariable (qui utilise BindingChain)
            return ctx.resolveVariable(varName)
            
        case "attribute":
            varName, _ := v["variable"].(string)
            attrName, _ := v["attribute"].(string)
            
            // Récupérer via BindingChain
            if !ctx.bindings.Has(varName) {
                return nil, fmt.Errorf("variable '%s' non trouvée", varName)
            }
            fact := ctx.bindings.Get(varName)
            
            if fact == nil {
                return nil, fmt.Errorf("variable '%s' est nil", varName)
            }
            
            // Accéder à l'attribut
            value, exists := fact.Attributes[attrName]
            if !exists {
                return nil, fmt.Errorf("attribut '%s' non trouvé dans %s", 
                    attrName, varName)
            }
            return value, nil
            
        // ... autres cas
        }
    }
    
    return arg, nil
}
```

---

### Tâche 3 : Adapter TerminalNode (30 min)

#### 3.1 Vérifier ActivateLeft

**Fichier** : `tsd/rete/node_terminal.go`

**S'assurer que** :

```go
func (tn *TerminalNode) ActivateLeft(token *Token) error {
    // Logging pour debug (TEMPORAIRE)
    fmt.Printf("\n🎯 [TERMINAL_%s] Received token\n", tn.ID)
    fmt.Printf("   Token ID: %s\n", token.ID)
    fmt.Printf("   Token Bindings: %v\n", token.GetVariables())
    fmt.Printf("   Rule: %s\n", tn.Rule.Name)
    
    // Créer le contexte d'exécution avec BindingChain
    ctx := NewExecutionContext(token, tn.Network, tn.Rule)
    
    // Exécuter l'action
    err := tn.executeAction(ctx)
    if err != nil {
        return fmt.Errorf("erreur lors de l'exécution de l'action '%s': %w", 
            tn.Rule.Action.Name, err)
    }
    
    return nil
}
```

**Points de vérification** :
- Le token reçu contient bien tous les bindings (vérifier avec logs)
- ExecutionContext est créé avec `token.Bindings`
- Les erreurs sont bien propagées

---

#### 3.2 Vérifier executeAction

**S'assurer que executeAction utilise bien le contexte** :

```go
func (tn *TerminalNode) executeAction(ctx *ExecutionContext) error {
    action := tn.Rule.Action
    
    // Évaluer les arguments de l'action
    evaluatedArgs := make([]interface{}, len(action.Arguments))
    for i, arg := range action.Arguments {
        value, err := ctx.evaluateArgument(arg)
        if err != nil {
            return fmt.Errorf("erreur évaluation argument %d: %w", i, err)
        }
        evaluatedArgs[i] = value
    }
    
    // Exécuter l'action avec les arguments évalués
    return tn.ActionExecutor.Execute(action.Name, evaluatedArgs, ctx)
}
```

---

### Tâche 4 : Tests Unitaires (50 min)

#### 4.1 Test de resolveVariable

**Fichier** : `tsd/rete/action_executor_evaluation_test.go`

```go
func TestExecutionContext_ResolveVariable(t *testing.T) {
    t.Log("🧪 TEST ExecutionContext - resolveVariable avec BindingChain")
    
    // Setup
    userFact := &Fact{ID: "u1", Type: "User", Attributes: map[string]interface{}{"id": 1}}
    orderFact := &Fact{ID: "o1", Type: "Order", Attributes: map[string]interface{}{"id": 100}}
    
    token := &Token{
        ID: "t1",
        Bindings: NewBindingChain().
            Add("user", userFact).
            Add("order", orderFact),
    }
    
    ctx := NewExecutionContext(token, nil, nil)
    
    // Test 1: Variable existante
    result, err := ctx.resolveVariable("user")
    if err != nil {
        t.Errorf("❌ Erreur inattendue: %v", err)
    }
    if result != userFact {
        t.Errorf("❌ Fact incorrect")
    }
    
    // Test 2: Variable inexistante
    _, err = ctx.resolveVariable("product")
    if err == nil {
        t.Errorf("❌ Devrait retourner une erreur")
    }
    
    // Vérifier que le message d'erreur contient les variables disponibles
    if !strings.Contains(err.Error(), "user") || !strings.Contains(err.Error(), "order") {
        t.Errorf("❌ Message d'erreur devrait lister les variables disponibles")
    }
    
    t.Log("✅ resolveVariable fonctionne correctement")
}
```

---

#### 4.2 Test d'intégration action complète

```go
func TestTerminalNode_ExecuteAction_AllVariablesAvailable(t *testing.T) {
    t.Log("🧪 TEST TerminalNode - Exécution action avec 3 variables")
    
    // Setup : Faits
    userFact := &Fact{
        ID: "u1", Type: "User",
        Attributes: map[string]interface{}{"id": 1, "name": "Alice"},
    }
    orderFact := &Fact{
        ID: "o1", Type: "Order",
        Attributes: map[string]interface{}{"id": 100, "user_id": 1},
    }
    productFact := &Fact{
        ID: "p1", Type: "Product",
        Attributes: map[string]interface{}{"id": 200, "name": "Laptop"},
    }
    
    // Token avec 3 bindings
    token := &Token{
        ID: "t_final",
        Facts: []*Fact{userFact, orderFact, productFact},
        Bindings: NewBindingChain().
            Add("user", userFact).
            Add("order", orderFact).
            Add("product", productFact),
    }
    
    // Règle avec action utilisant les 3 variables
    rule := &Rule{
        Name: "test_rule",
        Action: &Action{
            Name: "log_purchase",
            Arguments: []interface{}{
                map[string]interface{}{"type": "attribute", "variable": "user", "attribute": "name"},
                map[string]interface{}{"type": "attribute", "variable": "product", "attribute": "name"},
                map[string]interface{}{"type": "attribute", "variable": "order", "attribute": "id"},
            },
        },
    }
    
    // Mock ActionExecutor
    var capturedArgs []interface{}
    mockExecutor := &MockActionExecutor{
        ExecuteFunc: func(name string, args []interface{}, ctx *ExecutionContext) error {
            capturedArgs = args
            return nil
        },
    }
    
    // TerminalNode
    terminal := &TerminalNode{
        BaseNode: BaseNode{ID: "terminal_test"},
        Rule: rule,
        ActionExecutor: mockExecutor,
    }
    
    // Act
    err := terminal.ActivateLeft(token)
    
    // Assert
    if err != nil {
        t.Fatalf("❌ Erreur d'exécution: %v", err)
    }
    
    if len(capturedArgs) != 3 {
        t.Errorf("❌ Attendu 3 arguments, got %d", len(capturedArgs))
    }
    
    // Vérifier les valeurs
    if capturedArgs[0] != "Alice" {
        t.Errorf("❌ Arg 0: attendu 'Alice', got %v", capturedArgs[0])
    }
    if capturedArgs[1] != "Laptop" {
        t.Errorf("❌ Arg 1: attendu 'Laptop', got %v", capturedArgs[1])
    }
    if capturedArgs[2] != 100 {
        t.Errorf("❌ Arg 2: attendu 100, got %v", capturedArgs[2])
    }
    
    t.Log("✅ Action exécutée avec toutes les variables")
}

type MockActionExecutor struct {
    ExecuteFunc func(string, []interface{}, *ExecutionContext) error
}

func (m *MockActionExecutor) Execute(name string, args []interface{}, ctx *ExecutionContext) error {
    return m.ExecuteFunc(name, args, ctx)
}
```

---

### Tâche 5 : Validation (30 min)

#### 5.1 Exécuter les tests

```bash
cd tsd

# Tests unitaires
go test -v ./rete/action_executor_*_test.go

# Tests de TerminalNode
go test -v ./rete/node_terminal_test.go

# Tous les tests rete
go test -v ./rete/...
```

#### 5.2 Vérifier avec tests d'intégration

```bash
make test-integration
```

#### 5.3 Supprimer le logging de debug

**Dans TerminalNode.ActivateLeft**, supprimer ou désactiver les `fmt.Printf`.

---

## ✅ Critères de Validation

### Code
- [ ] ExecutionContext utilise *BindingChain
- [ ] resolveVariable utilise bindings.Has() et bindings.Get()
- [ ] Messages d'erreur listent les variables disponibles
- [ ] TerminalNode propage correctement les bindings
- [ ] Tous les accès à varCache remplacés

### Tests
- [ ] TestExecutionContext_ResolveVariable passe
- [ ] TestTerminalNode_ExecuteAction_AllVariablesAvailable passe
- [ ] Tests d'intégration passent
- [ ] Pas de régression

### Qualité
- [ ] Code formatté et sans warnings
- [ ] Logging debug supprimé
- [ ] GoDoc présent

---

## 🎯 Prochaine Étape

Passer au **Prompt 09 - Tests Cascades Multi-Variables**.

Les tests du Prompt 09 valideront que tout le système fonctionne correctement pour N variables.