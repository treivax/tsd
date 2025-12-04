# Phase 2 : Barrière de Synchronisation - Document de Conception

## Date
2025-01-XX

## Contexte
Suite à l'implémentation réussie de la Phase 1 (transactions implicites renforcées), nous passons maintenant à la Phase 2 qui vise à ajouter une barrière de synchronisation explicite pour garantir la visibilité immédiate des faits après leur soumission.

## Objectifs de la Phase 2

### Objectif Principal
Garantir qu'après le retour de `SubmitFactsFromGrammar()`, tous les faits ont été :
1. Soumis au réseau RETE
2. Persistés dans le storage
3. Propagés dans les chaînes alpha/beta
4. Visibles pour les lectures ultérieures

### Objectifs Secondaires
- Ajouter un mécanisme de retry avec backoff pour la vérification de persistance
- Implémenter un timeout configurable pour éviter les blocages
- Maintenir la compatibilité avec le code existant
- Éviter d'introduire de nouveaux problèmes de concurrence

## Analyse de la Situation Actuelle

### Code Existant (Phase 1)
```go
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
    factsSubmitted := 0
    factsPersisted := 0
    
    for i, factMap := range facts {
        // Conversion map -> Fact
        fact := convertToFact(factMap)
        
        // Soumission
        if err := rn.SubmitFact(fact); err != nil {
            return err
        }
        factsSubmitted++
        
        // Vérification immédiate
        internalID := fact.GetInternalID()
        if rn.Storage.GetFact(internalID) != nil {
            factsPersisted++
        } else {
            // ⚠️ Warning mais pas d'erreur
            tsdio.Printf("⚠️  Fait %s non persisté immédiatement\n", fact.ID)
        }
    }
    
    // Vérification finale
    if factsSubmitted != factsPersisted {
        return fmt.Errorf("incohérence: %d soumis, %d persistés", 
            factsSubmitted, factsPersisted)
    }
    
    return nil
}
```

### Problèmes Résiduels
1. **Vérification unique** : Un seul appel à `GetFact()` peut manquer un fait qui n'est pas encore visible
2. **Pas de retry** : Si un fait n'est pas immédiatement visible, on échoue directement
3. **Timeout absent** : Aucune protection contre des attentes infinies
4. **Synchronisation faible** : Pas de garantie que la propagation RETE est complète

## Conception Phase 2

### Architecture Proposée

#### 1. Configuration du Timeout
```go
// network.go - Ajouter dans ReteNetwork
type ReteNetwork struct {
    // ... champs existants ...
    
    // Configuration de synchronisation (Phase 2)
    SubmissionTimeout time.Duration  // Timeout global pour soumission de faits
    VerifyRetryDelay  time.Duration  // Délai entre tentatives de vérification
    MaxVerifyRetries  int            // Nombre max de tentatives de vérification
}

// Valeurs par défaut
const (
    DefaultSubmissionTimeout = 30 * time.Second
    DefaultVerifyRetryDelay  = 10 * time.Millisecond
    DefaultMaxVerifyRetries  = 10
)
```

#### 2. Fonction d'Attente avec Retry
```go
// waitForFactPersistence attend qu'un fait soit persisté avec retry + backoff
func (rn *ReteNetwork) waitForFactPersistence(fact *Fact, timeout time.Duration) error {
    internalID := fact.GetInternalID()
    deadline := time.Now().Add(timeout)
    attempt := 0
    
    for time.Now().Before(deadline) {
        attempt++
        
        // Vérifier si le fait est persisté
        if storedFact := rn.Storage.GetFact(internalID); storedFact != nil {
            // ✅ Fait trouvé
            if attempt > 1 {
                tsdio.Printf("✅ Fait %s persisté après %d tentative(s)\n", 
                    fact.ID, attempt)
            }
            return nil
        }
        
        // Backoff exponentiel (limité)
        if attempt < rn.MaxVerifyRetries {
            backoff := rn.VerifyRetryDelay * time.Duration(1<<uint(attempt-1))
            if backoff > 500*time.Millisecond {
                backoff = 500 * time.Millisecond
            }
            time.Sleep(backoff)
        } else {
            // Dernier essai : attendre le reste du timeout
            time.Sleep(100 * time.Millisecond)
        }
    }
    
    // ❌ Timeout dépassé
    return fmt.Errorf("timeout: fait %s (ID interne: %s) non persisté après %v",
        fact.ID, internalID, timeout)
}
```

#### 3. SubmitFactsFromGrammar Améliorée
```go
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
    if len(facts) == 0 {
        return nil
    }
    
    // Timeout par fait : timeout total divisé par nombre de faits
    // Minimum 1 seconde par fait
    timeoutPerFact := rn.SubmissionTimeout / time.Duration(len(facts))
    if timeoutPerFact < 1*time.Second {
        timeoutPerFact = 1 * time.Second
    }
    
    factsSubmitted := 0
    factsPersisted := 0
    startTime := time.Now()
    
    for i, factMap := range facts {
        // 1. Conversion map -> Fact
        factID := fmt.Sprintf("fact_%d", i)
        if id, ok := factMap["id"].(string); ok {
            factID = id
        }
        
        factType := "unknown"
        if typ, ok := factMap["type"].(string); ok {
            factType = typ
        } else if typ, ok := factMap["reteType"].(string); ok {
            factType = typ
        }
        
        fact := &Fact{
            ID:     factID,
            Type:   factType,
            Fields: make(map[string]interface{}),
        }
        
        for key, value := range factMap {
            if key != "type" && key != "reteType" {
                fact.Fields[key] = value
            }
        }
        
        // 2. Soumission au réseau RETE
        if err := rn.SubmitFact(fact); err != nil {
            return fmt.Errorf("erreur soumission fait %s: %w", fact.ID, err)
        }
        factsSubmitted++
        
        // 3. Barrière de synchronisation : attendre la persistance
        if err := rn.waitForFactPersistence(fact, timeoutPerFact); err != nil {
            return fmt.Errorf("échec synchronisation fait %s: %w", fact.ID, err)
        }
        factsPersisted++
    }
    
    duration := time.Since(startTime)
    
    // 4. Vérification finale de cohérence
    if factsSubmitted != factsPersisted {
        return fmt.Errorf("incohérence détectée: %d faits soumis mais seulement %d persistés",
            factsSubmitted, factsPersisted)
    }
    
    tsdio.Printf("✅ Phase 2 - Synchronisation complète: %d/%d faits persistés en %v\n", 
        factsPersisted, factsSubmitted, duration)
    
    return nil
}
```

#### 4. Initialisation avec Valeurs par Défaut
```go
// NewReteNetworkWithConfig - Ajouter initialisation des nouveaux champs
func NewReteNetworkWithConfig(storage Storage, config *ReteConfig) *ReteNetwork {
    network := &ReteNetwork{
        // ... champs existants ...
        
        // Phase 2: Configuration de synchronisation
        SubmissionTimeout: DefaultSubmissionTimeout,
        VerifyRetryDelay:  DefaultVerifyRetryDelay,
        MaxVerifyRetries:  DefaultMaxVerifyRetries,
    }
    
    // Si config spécifie des valeurs personnalisées, les appliquer
    if config != nil {
        if config.SubmissionTimeout > 0 {
            network.SubmissionTimeout = config.SubmissionTimeout
        }
        if config.VerifyRetryDelay > 0 {
            network.VerifyRetryDelay = config.VerifyRetryDelay
        }
        if config.MaxVerifyRetries > 0 {
            network.MaxVerifyRetries = config.MaxVerifyRetries
        }
    }
    
    return network
}
```

#### 5. Extension de ReteConfig (Optionnel)
```go
// config.go
type ReteConfig struct {
    // ... champs existants ...
    
    // Phase 2: Configuration de synchronisation
    SubmissionTimeout time.Duration
    VerifyRetryDelay  time.Duration
    MaxVerifyRetries  int
}
```

### Stratégie Séquentielle vs Parallèle

**Décision : Rester Séquentiel**

Raisons :
1. **Thread-safety** : Le réseau RETE n'est pas conçu pour des soumissions concurrentes sans verrous complexes
2. **Ordre de propagation** : Les règles peuvent dépendre de l'ordre de soumission des faits
3. **Transactions** : Les transactions sont actuellement séquentielles
4. **Simplicité** : Moins de risques de race conditions

Si nécessaire ultérieurement, on pourra :
- Ajouter un mode "parallèle" opt-in
- Utiliser un pool de workers avec file d'attente
- Coordonner avec des channels

### Mécanisme de Retry avec Backoff

**Stratégie Adoptée : Backoff Exponentiel Limité**

```
Tentative 1 : Vérification immédiate (0ms)
Tentative 2 : Attente 10ms
Tentative 3 : Attente 20ms
Tentative 4 : Attente 40ms
Tentative 5 : Attente 80ms
Tentative 6 : Attente 160ms
Tentative 7 : Attente 320ms
Tentative 8+ : Attente 500ms (plafond)
```

**Avantages** :
- Détection rapide des faits déjà persistés
- Réduction de la charge CPU (pas de busy-wait)
- Protection contre les attentes infinies

## Plan d'Implémentation

### Étape 1 : Ajouter les Champs de Configuration
- Fichier : `rete/network.go`
- Ajouter `SubmissionTimeout`, `VerifyRetryDelay`, `MaxVerifyRetries` à `ReteNetwork`
- Ajouter les constantes par défaut

### Étape 2 : Implémenter waitForFactPersistence
- Fichier : `rete/network.go`
- Nouvelle méthode privée avec retry + backoff
- Logging approprié

### Étape 3 : Modifier SubmitFactsFromGrammar
- Fichier : `rete/network.go`
- Intégrer `waitForFactPersistence` après chaque `SubmitFact`
- Calculer timeout par fait
- Améliorer les logs

### Étape 4 : Initialiser les Valeurs par Défaut
- Fichier : `rete/network.go`
- Modifier `NewReteNetwork` et `NewReteNetworkWithConfig`
- Assurer la rétro-compatibilité

### Étape 5 : Étendre ReteConfig (Optionnel)
- Fichier : `rete/config.go`
- Ajouter champs de configuration si besoin

### Étape 6 : Tests
- Fichier : `rete/coherence_phase2_test.go`
- Tests de synchronisation
- Tests de timeout
- Tests de retry

## Tests de Validation

### Test 1 : Synchronisation Basique
```go
func TestPhase2_BasicSynchronization(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    
    facts := []map[string]interface{}{
        {"id": "fact1", "type": "Product", "name": "Item1"},
        {"id": "fact2", "type": "Product", "name": "Item2"},
    }
    
    err := network.SubmitFactsFromGrammar(facts)
    require.NoError(t, err)
    
    // Vérification immédiate : tous les faits doivent être visibles
    for _, fm := range facts {
        id := fm["id"].(string)
        typ := fm["type"].(string)
        internalID := typ + "_" + id
        
        fact := storage.GetFact(internalID)
        assert.NotNil(t, fact, "Fait %s doit être immédiatement visible", id)
    }
}
```

### Test 2 : Retry Avec Délai
```go
func TestPhase2_RetryMechanism(t *testing.T) {
    storage := NewMemoryStorageWithDelay(50 * time.Millisecond) // Simule délai
    network := NewReteNetwork(storage)
    network.VerifyRetryDelay = 10 * time.Millisecond
    network.MaxVerifyRetries = 5
    
    facts := []map[string]interface{}{
        {"id": "delayed_fact", "type": "Product", "name": "Delayed"},
    }
    
    start := time.Now()
    err := network.SubmitFactsFromGrammar(facts)
    duration := time.Since(start)
    
    require.NoError(t, err)
    assert.Greater(t, duration, 50*time.Millisecond, "Devrait attendre le délai")
    assert.Less(t, duration, 200*time.Millisecond, "Ne devrait pas attendre trop longtemps")
}
```

### Test 3 : Timeout
```go
func TestPhase2_Timeout(t *testing.T) {
    storage := NewBrokenStorage() // Ne persiste jamais les faits
    network := NewReteNetwork(storage)
    network.SubmissionTimeout = 100 * time.Millisecond
    network.VerifyRetryDelay = 10 * time.Millisecond
    
    facts := []map[string]interface{}{
        {"id": "never_persisted", "type": "Product", "name": "Ghost"},
    }
    
    err := network.SubmitFactsFromGrammar(facts)
    require.Error(t, err)
    assert.Contains(t, err.Error(), "timeout")
}
```

### Test 4 : Backoff Exponentiel
```go
func TestPhase2_ExponentialBackoff(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    
    attempts := []time.Duration{}
    
    // Mock pour capturer les délais
    originalSleep := time.Sleep
    time.Sleep = func(d time.Duration) {
        attempts = append(attempts, d)
        originalSleep(d)
    }
    defer func() { time.Sleep = originalSleep }()
    
    // Test avec storage qui retarde la persistance
    // Vérifier que les délais augmentent exponentiellement
}
```

### Test 5 : Concurrent Reads After Write
```go
func TestPhase2_ConcurrentReadsAfterWrite(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    
    facts := []map[string]interface{}{
        {"id": "concurrent_fact", "type": "Product", "name": "Concurrent"},
    }
    
    // Soumettre le fait
    err := network.SubmitFactsFromGrammar(facts)
    require.NoError(t, err)
    
    // Lancer plusieurs lectures concurrentes immédiatement après
    var wg sync.WaitGroup
    errors := make(chan error, 10)
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            fact := storage.GetFact("Product_concurrent_fact")
            if fact == nil {
                errors <- fmt.Errorf("fait non visible dans goroutine")
            }
        }()
    }
    
    wg.Wait()
    close(errors)
    
    for err := range errors {
        require.NoError(t, err)
    }
}
```

## Métriques de Succès

### Critères d'Acceptation
- ✅ Tous les tests Phase 1 continuent de passer
- ✅ Tous les nouveaux tests Phase 2 passent
- ✅ `go test -race ./rete/...` passe sans erreurs
- ✅ Tests d'intégration passent de manière déterministe
- ✅ Performance acceptable (overhead < 10% pour fichiers normaux)
- ✅ Timeout fonctionne correctement (pas de blocages infinis)

### Benchmarks
```bash
# Avant Phase 2
BenchmarkIngestFile-8         100      12.5 ms/op

# Après Phase 2 (objectif)
BenchmarkIngestFile-8         100      13.5 ms/op    (< 10% overhead)
```

## Risques et Mitigations

### Risque 1 : Overhead de Performance
**Probabilité** : Moyenne  
**Impact** : Faible

**Mitigation** :
- Backoff intelligent (pas de busy-wait)
- Timeout par fait proportionnel
- Mode "fast path" si fait visible immédiatement

### Risque 2 : Faux Positifs de Timeout
**Probabilité** : Faible  
**Impact** : Moyen

**Mitigation** :
- Timeout généreux par défaut (30s total)
- Configurable via ReteConfig
- Logs détaillés pour diagnostic

### Risque 3 : Régression des Tests Existants
**Probabilité** : Faible  
**Impact** : Élevé

**Mitigation** :
- Valeurs par défaut conservatrices
- Rétro-compatibilité totale
- Tests avant/après

## Compatibilité

### Rétro-Compatibilité
- ✅ Aucun changement d'interface publique
- ✅ Comportement par défaut conservateur
- ✅ Pas de breaking changes

### Migration
Aucune migration nécessaire : les changements sont transparents.

## Documentation

### À Mettre à Jour
- [ ] `COHERENCE_FIX_SUMMARY.md` - Ajouter section Phase 2
- [ ] Commentaires dans `network.go`
- [ ] CHANGELOG.md - Nouvelle entrée

### Nouveaux Documents
- [x] `COHERENCE_FIX_PHASE2_DESIGN.md` - Ce document
- [ ] `COHERENCE_FIX_PHASE2_IMPLEMENTATION.md` - Notes d'implémentation

## Chronologie

- **J1 Matin** : Ajouter champs de configuration + constantes
- **J1 Après-midi** : Implémenter `waitForFactPersistence()`
- **J2 Matin** : Modifier `SubmitFactsFromGrammar()`
- **J2 Après-midi** : Tests Phase 2
- **J3** : Validation, benchmarks, documentation

## Validation Finale

- [ ] Code review complet
- [ ] Tous les tests passent (Phase 1 + Phase 2)
- [ ] Tests d'intégration passent de manière déterministe
- [ ] Benchmarks documentés
- [ ] Documentation mise à jour
- [ ] Commit + Push

## Notes Additionnelles

- Cette phase est conservatrice : elle ajoute des garanties sans changer l'architecture
- Le mécanisme de retry est essentiel pour gérer les variations de timing système
- Le timeout protège contre les deadlocks et les bugs de persistance
- La Phase 3 (métriques détaillées) pourra enrichir cette base avec plus d'observabilité