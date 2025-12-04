# Plan de Correction des Problèmes de Cohérence

## Date
2025-01-XX

## Contexte
Suite à la migration vers des transactions RETE thread-safe (Command Pattern), certains tests d'intégration révèlent des problèmes de cohérence dans la pipeline d'ingestion. Les faits ne sont pas toujours visibles immédiatement après l'appel à `IngestFile()`.

## Diagnostic

### Problème Racine Identifié
**Manque de garanties d'atomicité et de cohérence dans `IngestFile()`**

Actuellement, `IngestFile()` :
1. ✅ Démarre une transaction implicite
2. ✅ Parse et valide le fichier
3. ✅ Ajoute les types et règles
4. ✅ Soumet les faits via `SubmitFactsFromGrammar()`
5. ✅ Commit la transaction
6. ❌ **MAIS** ne garantit pas que tous les faits sont persistés et propagés avant de retourner

### Symptômes Observés
- Comptage incorrect de faits immédiatement après `IngestFile()`
- Tests d'intégration échouant de manière intermittente
- Problèmes de "read-after-write" dans les scénarios concurrents

### Cause Technique
`SubmitFactsFromGrammar()` itère et soumet les faits de manière séquentielle, mais :
- Aucune barrière de synchronisation explicite
- Aucune vérification que tous les faits ont été persistés avant le commit
- La propagation dans le réseau RETE peut être asynchrone dans certains cas

## Plan d'Action Priorisé

### Phase 1 : Transaction Implicite Renforcée (CRITIQUE)
**Impact** : Résout le problème racine  
**Complexité** : Moyenne  
**Risque** : Faible (infrastructure transaction existe déjà)

#### Actions
1. ✅ Transaction déjà créée dans `ingestFileWithMetrics()`
2. ⚠️ Ajouter des compteurs atomiques pour tracking
3. ⚠️ Vérifier cohérence avant commit
4. ⚠️ Implémenter `Storage.Sync()` pour garantir durabilité

#### Fichiers à Modifier
- `rete/constraint_pipeline.go` : Renforcer `ingestFileWithMetrics()`
- `rete/interfaces.go` : Ajouter `Sync()` à l'interface `Storage`
- `rete/store_base.go` : Implémenter `Sync()` pour `MemoryStorage`
- `rete/network.go` : Ajouter vérifications dans `SubmitFactsFromGrammar()`

### Phase 2 : Barrière de Synchronisation (IMPORTANT)
**Impact** : Garantie supplémentaire de sécurité  
**Complexité** : Faible  
**Risque** : Très faible

#### Actions
1. Ajouter `sync.WaitGroup` dans `SubmitFactsFromGrammar()`
2. Attendre confirmation de persistance pour tous les faits
3. Timeout configurable pour éviter blocages

#### Fichiers à Modifier
- `rete/network.go` : `SubmitFactsFromGrammar()` avec barrière

### Phase 3 : Audit et Validation (QUALITÉ)
**Impact** : Meilleure observabilité  
**Complexité** : Faible  
**Risque** : Nul (ajout uniquement)

#### Actions
1. Ajouter métriques internes (factsSubmitted, factsPersisted, factsPropagated)
2. Logs structurés pour debugging
3. Assertions de cohérence avant commit

#### Fichiers à Modifier
- `rete/constraint_pipeline.go` : Ajouter assertions
- `rete/metrics.go` : Enrichir les métriques

### Phase 4 : Storage Sync et Mode Cohérence (OPTIONNEL)
**Impact** : Architecture long terme  
**Complexité** : Moyenne-Haute  
**Risque** : Moyen (changements d'interface)

#### Actions
1. Introduire `ConsistencyMode` (Strong/Relaxed)
2. Configurer mode via options du pipeline
3. Implémenter `Storage.Flush()` pour persistance garantie

#### Fichiers à Modifier
- `rete/pipeline_options.go` : Nouvelle structure d'options
- `rete/constraint_pipeline.go` : Support des modes
- `rete/interfaces.go` : Interface étendue

## Implémentation Détaillée

### Phase 1.1 : Ajouter Storage.Sync()

```go
// interfaces.go
type Storage interface {
    // ... méthodes existantes ...
    Sync() error  // Garantit que toutes les écritures sont durables
}

// store_base.go
func (ms *MemoryStorage) Sync() error {
    ms.mutex.Lock()
    defer ms.mutex.Unlock()
    // Pour MemoryStorage, Sync() est un no-op car tout est en mémoire
    // Mais on pourrait vérifier la cohérence interne ici
    return nil
}
```

### Phase 1.2 : Compteurs Atomiques dans SubmitFactsFromGrammar()

```go
// network.go
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
    factsSubmitted := 0
    factsPersisted := 0
    
    for i, factMap := range facts {
        // ... conversion existante ...
        
        if err := rn.SubmitFact(fact); err != nil {
            return fmt.Errorf("erreur soumission fait %s: %w", fact.ID, err)
        }
        factsSubmitted++
        
        // Vérifier que le fait a été persisté
        if rn.Storage.GetFact(fact.ID) != nil {
            factsPersisted++
        }
    }
    
    // Garantir cohérence
    if factsSubmitted != factsPersisted {
        return fmt.Errorf("incohérence: %d faits soumis mais %d persistés", 
            factsSubmitted, factsPersisted)
    }
    
    return nil
}
```

### Phase 1.3 : Vérification Avant Commit

```go
// constraint_pipeline.go - dans ingestFileWithMetrics()
// AVANT le commit de la transaction
if tx != nil && tx.IsActive {
    // Vérifier que tous les faits soumis sont bien dans le storage
    expectedFactCount := len(factsForRete)
    actualFactCount := 0
    
    for _, factMap := range factsForRete {
        factID := factMap["id"].(string)
        if network.Storage.GetFact(factID) != nil {
            actualFactCount++
        }
    }
    
    if expectedFactCount != actualFactCount {
        return rollbackOnError(fmt.Errorf(
            "incohérence pré-commit: %d faits attendus mais %d trouvés",
            expectedFactCount, actualFactCount))
    }
    
    // Synchroniser le storage avant commit
    if err := storage.Sync(); err != nil {
        return rollbackOnError(fmt.Errorf("erreur sync storage: %w", err))
    }
    
    commitErr := tx.Commit()
    if commitErr != nil {
        return rollbackOnError(fmt.Errorf("erreur commit transaction: %w", commitErr))
    }
}
```

### Phase 2 : Barrière de Synchronisation

```go
// network.go
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(facts))
    
    for i, factMap := range facts {
        wg.Add(1)
        
        go func(idx int, fm map[string]interface{}) {
            defer wg.Done()
            
            // ... conversion et soumission ...
            
            if err := rn.SubmitFact(fact); err != nil {
                errChan <- fmt.Errorf("fact %d: %w", idx, err)
                return
            }
            
            // Vérifier persistance
            if rn.Storage.GetFact(fact.ID) == nil {
                errChan <- fmt.Errorf("fact %s not persisted", fact.ID)
            }
        }(i, factMap)
    }
    
    // Attendre avec timeout
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        close(errChan)
        // Vérifier les erreurs
        for err := range errChan {
            if err != nil {
                return err
            }
        }
        return nil
    case <-time.After(30 * time.Second):
        return fmt.Errorf("timeout waiting for fact submission")
    }
}
```

## Tests de Validation

### Test 1 : Read-After-Write
```go
func TestIngestFileReadAfterWrite(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    pipeline := NewConstraintPipeline()
    
    // Ingest file
    network, err := pipeline.IngestFile("test.tsd", network, storage)
    require.NoError(t, err)
    
    // Vérification IMMÉDIATE
    facts := storage.GetAllFacts()
    assert.Equal(t, expectedCount, len(facts), "Faits doivent être immédiatement visibles")
}
```

### Test 2 : Concurrent Ingestion
```go
func TestConcurrentIngestion(t *testing.T) {
    var wg sync.WaitGroup
    errors := make(chan error, 10)
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // Ingest concurrent
            if err := ingestAndVerify(id); err != nil {
                errors <- err
            }
        }(i)
    }
    
    wg.Wait()
    close(errors)
    
    for err := range errors {
        require.NoError(t, err)
    }
}
```

## Métriques de Succès

- ✅ Tous les tests unitaires passent avec `-race`
- ✅ Tous les tests d'intégration passent avec `-race -tags=integration`
- ✅ Pas de problèmes de comptage de faits
- ✅ Read-after-write garanti
- ✅ Performance non dégradée (< 5% overhead)

## Risques et Mitigations

### Risque 1 : Overhead de Performance
**Mitigation** : 
- Benchmarks avant/après
- Mode "Relaxed" optionnel pour cas non-critiques

### Risque 2 : Timeout sur Grosses Ingestions
**Mitigation** :
- Timeout configurable
- Logs détaillés pour diagnostic

### Risque 3 : Deadlock sur Erreurs
**Mitigation** :
- Defer unlock sur tous les mutex
- Rollback automatique en cas d'erreur

## Ordre d'Implémentation

1. **Jour 1** : Phase 1.1 (Storage.Sync)
2. **Jour 1** : Phase 1.2 (Compteurs atomiques)
3. **Jour 2** : Phase 1.3 (Vérification avant commit)
4. **Jour 2** : Tests de validation Phase 1
5. **Jour 3** : Phase 2 (Barrière de synchronisation)
6. **Jour 3** : Tests de validation Phase 2
7. **Jour 4** : Phase 3 (Audit et métriques)
8. **Jour 5** : Phase 4 (Modes de cohérence) - si nécessaire

## Validation Finale

- [ ] Code review complet
- [ ] Tests passent en local
- [ ] Tests passent en CI avec `-race`
- [ ] Benchmarks comparatifs documentés
- [ ] Documentation mise à jour
- [ ] CHANGELOG.md mis à jour

## Notes

- Ce plan s'appuie sur l'infrastructure existante (transactions Command Pattern)
- Les changements sont majoritairement additifs (faible risque)
- La stratégie est incrémentale : chaque phase peut être validée indépendamment