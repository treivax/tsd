# Prompt 08 - Tests complets du système xuples

## 🎯 Objectif

Tester de manière exhaustive l'ensemble du système xuples, incluant :
- Le parsing de la commande `xuple-space`
- Les actions par défaut (notamment Xuple)
- L'exécution immédiate des actions par RETE
- Le module xuples avec toutes ses politiques
- L'intégration complète RETE ↔ xuples
- Les scénarios de bout en bout

Cette phase de test garantit que tout fonctionne ensemble de manière cohérente et robuste.

## 📋 Tâches

### 1. Créer une suite de tests end-to-end

**Objectif** : Tester des scénarios complets du parsing à l'exécution.

**Fichier à créer** : `tsd/tests/e2e/xuples_e2e_test.go`

**Tests attendus** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
    "testing"
    "time"
    
    "tsd/compiler"
    "tsd/parser"
)

func TestXuplesSystem_EndToEnd_BasicScenario(t *testing.T) {
    t.Log("🧪 TEST E2E - SCÉNARIO BASIQUE XUPLES")
    
    program := `
        // Xuple-space pour les tâches
        xuple-space tasks {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        // Types de faits
        fact Task(id: string, priority: int, description: string)
        fact User(name: string)
        
        // Règle simple
        rule "create-task-xuple" {
            when {
                task: Task(priority >= 5)
            }
            then {
                Print("Task " + task.id + " added to xuple-space")
                Xuple("tasks", task)
            }
        }
    `
    
    // 1. Parser le programme
    t.Log("  📝 Parsing...")
    parseResult, err := parser.ParseTSD(program)
    if err != nil {
        t.Fatalf("❌ Parsing failed: %v", err)
    }
    
    // Vérifier que le xuple-space est parsé
    if len(parseResult.XupleSpaces) != 1 {
        t.Fatalf("❌ Expected 1 xuple-space, got %d", len(parseResult.XupleSpaces))
    }
    
    if parseResult.XupleSpaces[0].Name != "tasks" {
        t.Errorf("❌ Wrong xuple-space name: %s", parseResult.XupleSpaces[0].Name)
    }
    
    // 2. Compiler
    t.Log("  ⚙️  Compiling...")
    comp, err := compiler.NewCompiler()
    if err != nil {
        t.Fatalf("❌ Compiler creation failed: %v", err)
    }
    
    err = comp.Compile(parseResult)
    if err != nil {
        t.Fatalf("❌ Compilation failed: %v", err)
    }
    
    // 3. Instancier les xuple-spaces
    t.Log("  🔧 Instantiating xuple-spaces...")
    err = comp.InstantiateXupleSpaces()
    if err != nil {
        t.Fatalf("❌ Xuple-space instantiation failed: %v", err)
    }
    
    // 4. Construire le réseau RETE et configurer
    t.Log("  🕸️  Building RETE network...")
    network := comp.BuildReteNetwork()
    
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    // 5. Insérer des faits
    t.Log("  📥 Inserting facts...")
    
    task1 := &rete.Fact{
        ID:   "task1",
        Type: "Task",
        Attributes: map[string]interface{}{
            "id":          "task1",
            "priority":    8,
            "description": "High priority task",
        },
    }
    
    task2 := &rete.Fact{
        ID:   "task2",
        Type: "Task",
        Attributes: map[string]interface{}{
            "id":          "task2",
            "priority":    3,
            "description": "Low priority task",
        },
    }
    
    err = network.InsertFact(task1)
    if err != nil {
        t.Fatalf("❌ Failed to insert task1: %v", err)
    }
    
    err = network.InsertFact(task2)
    if err != nil {
        t.Fatalf("❌ Failed to insert task2: %v", err)
    }
    
    // 6. Vérifier les xuples créés
    t.Log("  🔍 Verifying xuples...")
    
    xupleManager := comp.GetXupleManager()
    xuplespace, err := xupleManager.GetXupleSpace("tasks")
    if err != nil {
        t.Fatalf("❌ Failed to get xuple-space: %v", err)
    }
    
    // Seul task1 devrait créer un xuple (priority >= 5)
    count := xuplespace.Count()
    if count != 1 {
        t.Errorf("❌ Expected 1 xuple, got %d", count)
    }
    
    // Récupérer le xuple
    xuple, err := xuplespace.Retrieve("agent1")
    if err != nil {
        t.Fatalf("❌ Failed to retrieve xuple: %v", err)
    }
    
    if xuple.Fact.ID != "task1" {
        t.Errorf("❌ Wrong fact in xuple: %s", xuple.Fact.ID)
    }
    
    t.Log("✅ E2E basic scenario passed")
}

func TestXuplesSystem_EndToEnd_MultiplePolicies(t *testing.T) {
    t.Log("🧪 TEST E2E - POLITIQUES MULTIPLES")
    
    program := `
        xuple-space fifo-once {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        xuple-space lifo-peragent {
            selection: lifo
            consumption: per-agent
            retention: duration(1h)
        }
        
        xuple-space random-limited {
            selection: random
            consumption: limited(3)
            retention: duration(30m)
        }
        
        fact Event(id: string, category: string)
        
        rule "fifo-events" {
            when {
                e: Event(category == "fifo")
            }
            then {
                Xuple("fifo-once", e)
            }
        }
        
        rule "lifo-events" {
            when {
                e: Event(category == "lifo")
            }
            then {
                Xuple("lifo-peragent", e)
            }
        }
        
        rule "random-events" {
            when {
                e: Event(category == "random")
            }
            then {
                Xuple("random-limited", e)
            }
        }
    `
    
    // Setup complet
    parseResult, _ := parser.ParseTSD(program)
    comp, _ := compiler.NewCompiler()
    comp.Compile(parseResult)
    comp.InstantiateXupleSpaces()
    network := comp.BuildReteNetwork()
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    // Insérer des événements de chaque catégorie
    categories := []string{"fifo", "lifo", "random"}
    for i, cat := range categories {
        for j := 0; j < 3; j++ {
            fact := &rete.Fact{
                ID:   fmt.Sprintf("event-%s-%d", cat, j),
                Type: "Event",
                Attributes: map[string]interface{}{
                    "id":       fmt.Sprintf("event-%s-%d", cat, j),
                    "category": cat,
                },
            }
            network.InsertFact(fact)
        }
    }
    
    // Vérifier chaque xuple-space
    xupleManager := comp.GetXupleManager()
    
    // FIFO-once : devrait avoir 3 xuples
    fifoSpace, _ := xupleManager.GetXupleSpace("fifo-once")
    if fifoSpace.Count() != 3 {
        t.Errorf("❌ fifo-once: expected 3 xuples, got %d", fifoSpace.Count())
    }
    
    // Tester la politique once
    xuple1, _ := fifoSpace.Retrieve("agent1")
    fifoSpace.MarkConsumed(xuple1.ID, "agent1")
    
    // Après consommation, le xuple ne devrait plus être disponible
    if fifoSpace.Count() != 2 {
        t.Errorf("❌ After once consumption, expected 2 available, got %d", fifoSpace.Count())
    }
    
    // LIFO-peragent : devrait avoir 3 xuples
    lifoSpace, _ := xupleManager.GetXupleSpace("lifo-peragent")
    if lifoSpace.Count() != 3 {
        t.Errorf("❌ lifo-peragent: expected 3 xuples, got %d", lifoSpace.Count())
    }
    
    // Tester la politique per-agent
    xuple2, _ := lifoSpace.Retrieve("agent1")
    lifoSpace.MarkConsumed(xuple2.ID, "agent1")
    
    // Agent1 a consommé, mais agent2 peut encore le consommer
    xuple3, err := lifoSpace.Retrieve("agent2")
    if err != nil || xuple3.ID != xuple2.ID {
        t.Error("❌ per-agent policy should allow different agents to consume same xuple")
    }
    
    // Random-limited : devrait avoir 3 xuples
    randomSpace, _ := xupleManager.GetXupleSpace("random-limited")
    if randomSpace.Count() != 3 {
        t.Errorf("❌ random-limited: expected 3 xuples, got %d", randomSpace.Count())
    }
    
    t.Log("✅ E2E multiple policies passed")
}

func TestXuplesSystem_EndToEnd_ComplexRule(t *testing.T) {
    t.Log("🧪 TEST E2E - RÈGLE COMPLEXE AVEC MULTIPLES FAITS")
    
    program := `
        xuple-space assignments {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        fact Person(name: string, age: int)
        fact Department(name: string, budget: int)
        fact Assignment(person: string, dept: string)
        
        rule "valid-assignment" {
            when {
                p: Person(age >= 18)
                d: Department(budget > 1000)
                a: Assignment(person == p.name, dept == d.name)
            }
            then {
                Print("Valid assignment: " + p.name + " -> " + d.name)
                Xuple("assignments", a)
            }
        }
    `
    
    // Setup
    parseResult, _ := parser.ParseTSD(program)
    comp, _ := compiler.NewCompiler()
    comp.Compile(parseResult)
    comp.InstantiateXupleSpaces()
    network := comp.BuildReteNetwork()
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    // Insérer les faits dans l'ordre
    person := &rete.Fact{
        ID:   "p1",
        Type: "Person",
        Attributes: map[string]interface{}{
            "name": "Alice",
            "age":  25,
        },
    }
    
    dept := &rete.Fact{
        ID:   "d1",
        Type: "Department",
        Attributes: map[string]interface{}{
            "name":   "Engineering",
            "budget": 5000,
        },
    }
    
    assignment := &rete.Fact{
        ID:   "a1",
        Type: "Assignment",
        Attributes: map[string]interface{}{
            "person": "Alice",
            "dept":   "Engineering",
        },
    }
    
    network.InsertFact(person)
    network.InsertFact(dept)
    network.InsertFact(assignment)
    
    // Vérifier le xuple créé
    xupleManager := comp.GetXupleManager()
    xuplespace, _ := xupleManager.GetXupleSpace("assignments")
    
    if xuplespace.Count() != 1 {
        t.Fatalf("❌ Expected 1 xuple, got %d", xuplespace.Count())
    }
    
    xuple, _ := xuplespace.Retrieve("agent1")
    
    // Le xuple devrait contenir le fait assignment
    if xuple.Fact.ID != "a1" {
        t.Errorf("❌ Wrong main fact: %s", xuple.Fact.ID)
    }
    
    // Le xuple devrait contenir les 3 faits déclencheurs
    if len(xuple.TriggeringFacts) != 3 {
        t.Errorf("❌ Expected 3 triggering facts, got %d", len(xuple.TriggeringFacts))
    }
    
    // Vérifier que les 3 faits sont présents
    factIDs := make(map[string]bool)
    for _, f := range xuple.TriggeringFacts {
        factIDs[f.ID] = true
    }
    
    if !factIDs["p1"] || !factIDs["d1"] || !factIDs["a1"] {
        t.Error("❌ Missing triggering facts")
    }
    
    t.Log("✅ E2E complex rule passed")
}

func TestXuplesSystem_EndToEnd_Retention(t *testing.T) {
    t.Log("🧪 TEST E2E - POLITIQUE DE RÉTENTION")
    
    program := `
        xuple-space short-lived {
            selection: fifo
            consumption: once
            retention: duration(1s)
        }
        
        fact Message(content: string)
        
        rule "ephemeral-message" {
            when {
                m: Message()
            }
            then {
                Xuple("short-lived", m)
            }
        }
    `
    
    // Setup
    parseResult, _ := parser.ParseTSD(program)
    comp, _ := compiler.NewCompiler()
    comp.Compile(parseResult)
    comp.InstantiateXupleSpaces()
    network := comp.BuildReteNetwork()
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    // Insérer un message
    msg := &rete.Fact{
        ID:   "msg1",
        Type: "Message",
        Attributes: map[string]interface{}{
            "content": "Temporary message",
        },
    }
    
    network.InsertFact(msg)
    
    xupleManager := comp.GetXupleManager()
    xuplespace, _ := xupleManager.GetXupleSpace("short-lived")
    
    // Immédiatement après, le xuple devrait être disponible
    if xuplespace.Count() != 1 {
        t.Fatalf("❌ Expected 1 xuple immediately, got %d", xuplespace.Count())
    }
    
    // Attendre l'expiration
    time.Sleep(1500 * time.Millisecond)
    
    // Nettoyer les xuples expirés
    cleaned := xuplespace.Cleanup()
    if cleaned != 1 {
        t.Errorf("❌ Expected to clean 1 xuple, cleaned %d", cleaned)
    }
    
    // Le xuple ne devrait plus être disponible
    if xuplespace.Count() != 0 {
        t.Errorf("❌ Expected 0 xuples after expiration, got %d", xuplespace.Count())
    }
    
    t.Log("✅ E2E retention policy passed")
}

func TestXuplesSystem_EndToEnd_ErrorHandling(t *testing.T) {
    t.Log("🧪 TEST E2E - GESTION D'ERREURS")
    
    // Test 1: Xuple-space non déclaré
    t.Run("undeclared-xuplespace", func(t *testing.T) {
        program := `
            fact Event(id: string)
            
            rule "bad-rule" {
                when {
                    e: Event()
                }
                then {
                    Xuple("nonexistent", e)
                }
            }
        `
        
        parseResult, _ := parser.ParseTSD(program)
        comp, _ := compiler.NewCompiler()
        comp.Compile(parseResult)
        comp.InstantiateXupleSpaces()
        network := comp.BuildReteNetwork()
        
        executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
        network.SetActionExecutor(executor)
        
        // Observer pour capturer les erreurs
        observer := NewTestActionObserver(t)
        network.SetActionObserver(observer)
        
        event := &rete.Fact{ID: "e1", Type: "Event"}
        network.InsertFact(event)
        
        // L'action devrait échouer
        executions := observer.GetExecutions()
        if len(executions) == 0 {
            t.Fatal("❌ Expected action execution")
        }
        
        if executions[0].Success {
            t.Error("❌ Expected action to fail for undeclared xuple-space")
        }
        
        if executions[0].Error == nil {
            t.Error("❌ Expected error for undeclared xuple-space")
        }
        
        t.Log("  ✅ Undeclared xuple-space error handled correctly")
    })
    
    // Test 2: Redéfinition de xuple-space
    t.Run("duplicate-xuplespace", func(t *testing.T) {
        program := `
            xuple-space myspace {
                selection: fifo
                consumption: once
                retention: unlimited
            }
            
            xuple-space myspace {
                selection: lifo
                consumption: once
                retention: unlimited
            }
        `
        
        _, err := parser.ParseTSD(program)
        if err == nil {
            t.Error("❌ Expected parsing error for duplicate xuple-space")
        }
        
        t.Log("  ✅ Duplicate xuple-space detected")
    })
    
    // Test 3: Redéfinition d'action par défaut
    t.Run("redefine-default-action", func(t *testing.T) {
        program := `
            action Xuple(space: string, f: any) {
                // Tentative de redéfinition
            }
        `
        
        parseResult, _ := parser.ParseTSD(program)
        comp, _ := compiler.NewCompiler()
        err := comp.Compile(parseResult)
        
        if err == nil {
            t.Error("❌ Expected compilation error for redefining default action")
        }
        
        t.Log("  ✅ Default action redefinition prevented")
    })
    
    t.Log("✅ E2E error handling passed")
}
```

**Livrables** :
- [ ] Suite complète de tests E2E
- [ ] Test scénario basique
- [ ] Test politiques multiples
- [ ] Test règles complexes
- [ ] Test rétention temporelle
- [ ] Test gestion d'erreurs
- [ ] Tous les tests passent

### 2. Créer des tests de performance

**Objectif** : S'assurer que le système est performant.

**Fichier à créer** : `tsd/tests/performance/xuples_perf_test.go`

**Tests attendus** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package performance

import (
    "fmt"
    "testing"
    
    "tsd/compiler"
    "tsd/parser"
)

func BenchmarkXupleCreation(b *testing.B) {
    program := `
        xuple-space bench {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        fact Event(id: string)
        
        rule "create-xuple" {
            when {
                e: Event()
            }
            then {
                Xuple("bench", e)
            }
        }
    `
    
    parseResult, _ := parser.ParseTSD(program)
    comp, _ := compiler.NewCompiler()
    comp.Compile(parseResult)
    comp.InstantiateXupleSpaces()
    network := comp.BuildReteNetwork()
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        event := &rete.Fact{
            ID:   fmt.Sprintf("e%d", i),
            Type: "Event",
            Attributes: map[string]interface{}{
                "id": fmt.Sprintf("e%d", i),
            },
        }
        network.InsertFact(event)
    }
}

func BenchmarkXupleRetrieval(b *testing.B) {
    // Setup avec beaucoup de xuples
    program := `
        xuple-space bench {
            selection: fifo
            consumption: per-agent
            retention: unlimited
        }
        
        fact Event(id: string)
        
        rule "create-xuple" {
            when {
                e: Event()
            }
            then {
                Xuple("bench", e)
            }
        }
    `
    
    parseResult, _ := parser.ParseTSD(program)
    comp, _ := compiler.NewCompiler()
    comp.Compile(parseResult)
    comp.InstantiateXupleSpaces()
    network := comp.BuildReteNetwork()
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    // Créer 1000 xuples
    for i := 0; i < 1000; i++ {
        event := &rete.Fact{
            ID:   fmt.Sprintf("e%d", i),
            Type: "Event",
        }
        network.InsertFact(event)
    }
    
    xupleManager := comp.GetXupleManager()
    xuplespace, _ := xupleManager.GetXupleSpace("bench")
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        agentID := fmt.Sprintf("agent%d", i%100)
        xuplespace.Retrieve(agentID)
    }
}

func BenchmarkPolicyEvaluation(b *testing.B) {
    program := `
        xuple-space fifo {
            selection: fifo
            consumption: once
            retention: unlimited
        }
        
        xuple-space lifo {
            selection: lifo
            consumption: per-agent
            retention: duration(1h)
        }
        
        xuple-space random {
            selection: random
            consumption: limited(5)
            retention: unlimited
        }
        
        fact Event(id: string, target: string)
        
        rule "fifo-rule" {
            when {
                e: Event(target == "fifo")
            }
            then {
                Xuple("fifo", e)
            }
        }
        
        rule "lifo-rule" {
            when {
                e: Event(target == "lifo")
            }
            then {
                Xuple("lifo", e)
            }
        }
        
        rule "random-rule" {
            when {
                e: Event(target == "random")
            }
            then {
                Xuple("random", e)
            }
        }
    `
    
    parseResult, _ := parser.ParseTSD(program)
    comp, _ := compiler.NewCompiler()
    comp.Compile(parseResult)
    comp.InstantiateXupleSpaces()
    network := comp.BuildReteNetwork()
    executor := actions.NewBuiltinActionExecutor(network, comp.GetXupleManager())
    network.SetActionExecutor(executor)
    
    targets := []string{"fifo", "lifo", "random"}
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        event := &rete.Fact{
            ID:   fmt.Sprintf("e%d", i),
            Type: "Event",
            Attributes: map[string]interface{}{
                "id":     fmt.Sprintf("e%d", i),
                "target": targets[i%3],
            },
        }
        network.InsertFact(event)
    }
}
```

**Livrables** :
- [ ] Benchmarks de création de xuples
- [ ] Benchmarks de récupération
- [ ] Benchmarks d'évaluation des politiques
- [ ] Résultats de performance documentés

### 3. Créer des tests de concurrence

**Objectif** : Vérifier la thread-safety du module xuples.

**Fichier à créer** : `tsd/xuples/concurrent_test.go`

**Tests attendus** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
    "fmt"
    "sync"
    "testing"
    
    "tsd/rete"
)

func TestXupleSpace_ConcurrentInsert(t *testing.T) {
    t.Log("🧪 TEST CONCURRENCE - INSERT")
    
    config := XupleSpaceConfig{
        Name:              "concurrent",
        SelectionPolicy:   NewFIFOSelectionPolicy(),
        ConsumptionPolicy: NewOnceConsumptionPolicy(),
        RetentionPolicy:   NewUnlimitedRetentionPolicy(),
    }
    
    xs := NewXupleSpace(config)
    
    const numGoroutines = 100
    const insertsPerGoroutine = 10
    
    var wg sync.WaitGroup
    wg.Add(numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < insertsPerGoroutine; j++ {
                xuple := &Xuple{
                    Fact: &rete.Fact{
                        ID: fmt.Sprintf("f-%d-%d", id, j),
                    },
                    CreatedAt: time.Now(),
                    Metadata: XupleMetadata{
                        State:      XupleStateAvailable,
                        ConsumedBy: make(map[string]time.Time),
                    },
                }
                
                err := xs.Insert(xuple)
                if err != nil {
                    t.Errorf("❌ Insert failed: %v", err)
                }
            }
        }(i)
    }
    
    wg.Wait()
    
    expectedCount := numGoroutines * insertsPerGoroutine
    actualCount := xs.Count()
    
    if actualCount != expectedCount {
        t.Errorf("❌ Expected %d xuples, got %d", expectedCount, actualCount)
    }
    
    t.Log("✅ Concurrent insert test passed")
}

func TestXupleSpace_ConcurrentRetrieveAndConsume(t *testing.T) {
    t.Log("🧪 TEST CONCURRENCE - RETRIEVE & CONSUME")
    
    config := XupleSpaceConfig{
        Name:              "concurrent",
        SelectionPolicy:   NewFIFOSelectionPolicy(),
        ConsumptionPolicy: NewPerAgentConsumptionPolicy(),
        RetentionPolicy:   NewUnlimitedRetentionPolicy(),
    }
    
    xs := NewXupleSpace(config)
    
    // Insérer des xuples
    for i := 0; i < 100; i++ {
        xuple := &Xuple{
            Fact:      &rete.Fact{ID: fmt.Sprintf("f%d", i)},
            CreatedAt: time.Now(),
            Metadata: XupleMetadata{
                State:      XupleStateAvailable,
                ConsumedBy: make(map[string]time.Time),
            },
        }
        xs.Insert(xuple)
    }
    
    const numAgents = 50
    
    var wg sync.WaitGroup
    wg.Add(numAgents)
    
    errors := make(chan error, numAgents)
    
    for i := 0; i < numAgents; i++ {
        go func(agentID string) {
            defer wg.Done()
            
            xuple, err := xs.Retrieve(agentID)
            if err != nil {
                errors <- err
                return
            }
            
            err = xs.MarkConsumed(xuple.ID, agentID)
            if err != nil {
                errors <- err
            }
        }(fmt.Sprintf("agent%d", i))
    }
    
    wg.Wait()
    close(errors)
    
    for err := range errors {
        t.Errorf("❌ Concurrent operation failed: %v", err)
    }
    
    t.Log("✅ Concurrent retrieve & consume test passed")
}

func TestXupleManager_ConcurrentCreateXupleSpace(t *testing.T) {
    t.Log("🧪 TEST CONCURRENCE - CREATE XUPLE-SPACE")
    
    manager := NewXupleManager()
    
    const numGoroutines = 100
    
    var wg sync.WaitGroup
    wg.Add(numGoroutines)
    
    successCount := int32(0)
    
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            defer wg.Done()
            
            config := XupleSpaceConfig{
                SelectionPolicy:   NewFIFOSelectionPolicy(),
                ConsumptionPolicy: NewOnceConsumptionPolicy(),
                RetentionPolicy:   NewUnlimitedRetentionPolicy(),
            }
            
            err := manager.CreateXupleSpace(fmt.Sprintf("space%d", id), config)
            if err == nil {
                atomic.AddInt32(&successCount, 1)
            }
        }(i)
    }
    
    wg.Wait()
    
    if int(successCount) != numGoroutines {
        t.Errorf("❌ Expected %d successful creations, got %d", numGoroutines, successCount)
    }
    
    t.Log("✅ Concurrent create xuple-space test passed")
}
```

**Livrables** :
- [ ] Tests de concurrence pour XupleSpace
- [ ] Tests de concurrence pour XupleManager
- [ ] Tests avec go test -race
- [ ] Aucune race condition détectée

### 4. Documenter les résultats de tests

**Objectif** : Créer un rapport de tests complet.

**Fichier à créer** : `tsd/docs/xuples/testing/test-report.md`

**Contenu attendu** :

```markdown
# Rapport de Tests - Système Xuples

## Résumé Exécutif

Date : [DATE]
Version : [VERSION]

- ✅ Tests unitaires : PASS (couverture XX%)
- ✅ Tests d'intégration : PASS
- ✅ Tests E2E : PASS
- ✅ Tests de performance : PASS
- ✅ Tests de concurrence : PASS

## Tests Unitaires

### Module xuples

- Structures de données : 100% couverture
- Politiques de sélection : 100% couverture
- Politiques de consommation : 100% couverture
- Politiques de rétention : 100% couverture
- XupleSpace : 95% couverture
- XupleManager : 95% couverture

### Module rete/actions

- BuiltinActionExecutor : 90% couverture
- Action Xuple : 100% couverture
- Extraction faits déclencheurs : 100% couverture

### Parser

- Parsing xuple-space : 100% couverture
- Validation : 100% couverture
- Détection doublons : 100% couverture

## Tests d'Intégration

| Test | Statut | Remarques |
|------|--------|-----------|
| RETE → Xuple action → xuples | ✅ PASS | - |
| Multiples xuple-spaces | ✅ PASS | - |
| Faits déclencheurs multiples | ✅ PASS | - |
| Validation compile-time | ✅ PASS | - |
| Validation runtime | ✅ PASS | - |

## Tests End-to-End

| Scénario | Statut | Temps |
|----------|--------|-------|
| Scénario basique | ✅ PASS | 50ms |
| Politiques multiples | ✅ PASS | 120ms |
| Règle complexe | ✅ PASS | 80ms |
| Rétention temporelle | ✅ PASS | 1.5s |
| Gestion d'erreurs | ✅ PASS | 100ms |

## Tests de Performance

### Benchmarks

```
BenchmarkXupleCreation-8          100000   12000 ns/op
BenchmarkXupleRetrieval-8         500000    3500 ns/op
BenchmarkPolicyEvaluation-8       200000    8000 ns/op
```

### Analyse

- Création de xuple : ~12μs (acceptable)
- Récupération : ~3.5μs (excellent)
- Évaluation politique : ~8μs (acceptable)

## Tests de Concurrence

- ✅ Insertion concurrente : PASS
- ✅ Récupération/consommation concurrente : PASS
- ✅ Création xuple-space concurrente : PASS
- ✅ go test -race : Aucune race condition détectée

## Couverture Globale

```
tsd/xuples              92.5%
tsd/rete/actions        88.0%
tsd/parser/ast          95.0%
tsd/compiler            85.0%
tsd/internal/defaultactions  100.0%
---
TOTAL                   90.1%
```

## Problèmes Identifiés

Aucun problème critique identifié.

## Recommandations

1. Continuer à maintenir la couverture > 80%
2. Ajouter plus de benchmarks pour les cas limites
3. Tests de charge avec des milliers de xuples
```

**Livrables** :
- [ ] Rapport de tests complet
- [ ] Métriques de couverture
- [ ] Résultats de performance
- [ ] Analyse des résultats

### 5. Créer un script de validation complète

**Objectif** : Script pour exécuter tous les tests du système xuples.

**Fichier à créer** : `tsd/scripts/test-xuples.sh`

**Contenu attendu** :

```bash
#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License

set -e

echo "🧪 VALIDATION COMPLÈTE DU SYSTÈME XUPLES"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Fonction pour afficher les résultats
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        exit 1
    fi
}

# 1. Tests unitaires du module xuples
echo "📦 Tests unitaires module xuples..."
go test -v -cover ./xuples/... > /tmp/xuples-unit.log 2>&1
print_result $? "Tests unitaires module xuples"

# 2. Tests unitaires rete/actions
echo "⚙️  Tests unitaires rete/actions..."
go test -v -cover ./rete/actions/... > /tmp/actions-unit.log 2>&1
print_result $? "Tests unitaires rete/actions"

# 3. Tests unitaires parser (xuple-space)
echo "📝 Tests parser xuple-space..."
go test -v -cover ./parser/... -run ".*Xuple.*" > /tmp/parser-unit.log 2>&1
print_result $? "Tests parser xuple-space"

# 4. Tests unitaires compiler
echo "🔧 Tests compiler..."
go test -v -cover ./compiler/... > /tmp/compiler-unit.log 2>&1
print_result $? "Tests compiler"

# 5. Tests unitaires defaultactions
echo "🎬 Tests default actions..."
go test -v -cover ./internal/defaultactions/... > /tmp/defaultactions-unit.log 2>&1
print_result $? "Tests default actions"

# 6. Tests d'intégration
echo "🔗 Tests d'intégration..."
go test -v ./tests/integration/... -run ".*Xuple.*" > /tmp/integration.log 2>&1
print_result $? "Tests d'intégration"

# 7. Tests E2E
echo "🌐 Tests end-to-end..."
go test -v ./tests/e2e/... > /tmp/e2e.log 2>&1
print_result $? "Tests E2E"

# 8. Tests de performance
echo "⚡ Tests de performance..."
go test -bench=. -benchmem ./tests/performance/... > /tmp/perf.log 2>&1
print_result $? "Tests de performance"

# 9. Tests de concurrence avec race detector
echo "🏃 Tests de concurrence (race detector)..."
go test -race ./xuples/... > /tmp/race.log 2>&1
print_result $? "Tests de concurrence"

# 10. Couverture globale
echo "📊 Calcul de la couverture globale..."
go test -coverprofile=/tmp/xuples-coverage.out ./xuples/... ./rete/actions/... ./parser/... ./compiler/... ./internal/defaultactions/... > /dev/null 2>&1
COVERAGE=$(go tool cover -func=/tmp/xuples-coverage.out | grep total | awk '{print $3}')
echo -e "${YELLOW}📊 Couverture globale: $COVERAGE${NC}"

# Vérifier seuil de couverture (80%)
COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
if (( $(echo "$COVERAGE_NUM < 80" | bc -l) )); then
    echo -e "${RED}❌ Couverture insuffisante (< 80%)${NC}"
    exit 1
else
    echo -e "${GREEN}✅ Couverture satisfaisante (>= 80%)${NC}"
fi

echo ""
echo "========================================"
echo -e "${GREEN}🎉 TOUTES LES VALIDATIONS SONT PASSÉES${NC}"
echo "========================================"
echo ""
echo "Logs disponibles dans /tmp/xuples-*.log"
```

**Livrables** :
- [ ] Script de validation créé
- [ ] Permissions d'exécution configurées
- [ ] Intégré dans le Makefile
- [ ] Documentation du script

### 6. Mettre à jour le Makefile

**Objectif** : Ajouter des cibles pour tester le système xuples.

**Fichier à modifier** : `tsd/Makefile`

**Ajouts attendus** :

```makefile
# Tests xuples
.PHONY: test-xuples test-xuples-unit test-xuples-integration test-xuples-e2e test-xuples-perf

test-xuples: ## Run all xuples tests
	@echo "🧪 Running all xuples tests..."
	@./scripts/test-xuples.sh

test-xuples-unit: ## Run xuples unit tests
	@echo "📦 Running xuples unit tests..."
	@go test -v -cover ./xuples/...
	@go test -v -cover ./rete/actions/...
	@go test -v -cover ./internal/defaultactions/...

test-xuples-integration: ## Run xuples integration tests
	@echo "🔗 Running xuples integration tests..."
	@go test -v ./tests/integration/... -run ".*Xuple.*"

test-xuples-e2e: ## Run xuples E2E tests
	@echo "🌐 Running xuples E2E tests..."
	@go test -v ./tests/e2e/...

test-xuples-perf: ## Run xuples performance tests
	@echo "⚡ Running xuples performance tests..."
	@go test -bench=. -benchmem ./tests/performance/...

test-xuples-race: ## Run xuples tests with race detector
	@echo "🏃 Running xuples tests with race detector..."
	@go test -race ./xuples/...
```

**Livrables** :
- [ ] Makefile mis à jour
- [ ] Toutes les cibles fonctionnelles
- [ ] Documentation dans `make help`

## 📁 Structure finale

```
tsd/
├── scripts/
│   └── test-xuples.sh                 # Nouveau
├── tests/
│   ├── e2e/
│   │   └── xuples_e2e_test.go         # Nouveau
│   └── performance/
│       └── xuples_perf_test.go        # Nouveau
├── xuples/
│   └── concurrent_test.go             # Nouveau
├── docs/xuples/
│   └── testing/
│       └── test-report.md             # Nouveau
└── Makefile                           # Modifié
```

## ✅ Critères de succès

- [ ] Suite complète de tests E2E créée
- [ ] Tests de performance (benchmarks)
- [ ] Tests de concurrence avec race detector
- [ ] Couverture globale > 80%
- [ ] Tous les tests passent
- [ ] Aucune race condition détectée
- [ ] Script de validation fonctionnel
- [ ] Makefile mis à jour
- [ ] Rapport de tests documenté
- [ ] `make test-xuples` passe sans erreur
- [ ] `make test-complete` inclut les tests xuples

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/design/` - Conception
- `tsd/docs/xuples/implementation/` - Documentation implémentation
- Go Testing Best Practices
- Go Benchmarking Guide

## 🎯 Prochaine étape

Une fois tous les tests passant, passer au prompt **09-finalize-documentation.md** pour finaliser toute la documentation du système xuples.