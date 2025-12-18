# TODO - Refactoring Xuples - Actions Restantes

**Date de création** : 2025-12-17  
**Dernière mise à jour** : 2025-12-17  
**Statut global** : Phase 1 ✅ Terminée | Phase 2 📋 En attente

---

## ✅ Phase 1 : Corrections et Améliorations (TERMINÉE)

- [x] Analyser l'implémentation actuelle (docs/xuples/analysis/)
- [x] Fix thread-safety de `generateTokenID()`
- [x] Refactoring `executeAction()` dans TerminalNode
- [x] Documentation avec TODOs pour xuples
- [x] Vérification non-régression (tests)
- [x] Création rapport Phase 1

---

## 📋 Phase 2 : Suite des Corrections (PROCHAINE ÉTAPE)

### A. Tests Manquants

#### 1. Test Concurrence generateTokenID ⏱️ 30 min

**Fichier à créer** : `rete/fact_token_concurrent_test.go`

**Code à ajouter** :
```go
func TestGenerateTokenID_Concurrent(t *testing.T) {
	t.Log("🧪 TEST CONCURRENCE - GÉNÉRATION TOKEN IDS")
	t.Log("==========================================")
	
	const goroutines = 1000
	ids := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := generateTokenID()
			mu.Lock()
			if ids[id] {
				t.Errorf("❌ ID dupliqué: %s", id)
			}
			ids[id] = true
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	if len(ids) != goroutines {
		t.Errorf("❌ Attendu %d IDs uniques, reçu %d", goroutines, len(ids))
	}
	
	t.Logf("✅ %d IDs uniques générés en concurrence", len(ids))
}

func TestGenerateTokenID_Sequential(t *testing.T) {
	t.Log("🧪 TEST SÉQUENTIEL - GÉNÉRATION TOKEN IDS")
	t.Log("=========================================")
	
	// Réinitialiser compteur (pour test isolé)
	// Note: En production, tokenCounter ne doit jamais être réinitialisé
	
	id1 := generateTokenID()
	id2 := generateTokenID()
	id3 := generateTokenID()
	
	if id1 == id2 || id2 == id3 || id1 == id3 {
		t.Error("❌ IDs ne sont pas uniques")
	}
	
	// Vérifier format
	if !strings.HasPrefix(id1, "token_") {
		t.Errorf("❌ Format incorrect: %s", id1)
	}
	
	t.Log("✅ IDs générés séquentiellement sont uniques")
}
```

**Estimation** : 30 minutes

#### 2. Tests Unitaires TerminalNode ⏱️ 2 heures

**Fichier à créer** : `rete/node_terminal_test.go`

**Tests à implémenter** :
```go
func TestTerminalNode_ActivateLeft_StoresToken(t *testing.T)
func TestTerminalNode_ActivateRetract_RemovesTokens(t *testing.T)
func TestTerminalNode_GetTriggeredActions(t *testing.T)
func TestTerminalNode_ExecuteAction_WithNilAction(t *testing.T)
func TestTerminalNode_ExecuteAction_WithValidAction(t *testing.T)
func TestTerminalNode_ConcurrentActivations(t *testing.T)
func TestTerminalNode_LogTupleSpaceActivation(t *testing.T)
func TestTerminalNode_FormatFact(t *testing.T)
```

**Priorité** : ⭐⭐⭐ Haute  
**Estimation** : 2 heures

#### 3. Tests Thread-Safety ActionRegistry ⏱️ 1 heure

**Fichier à créer** : `rete/action_handler_concurrent_test.go`

**Tests à implémenter** :
```go
func TestActionRegistry_ConcurrentRegister(t *testing.T)
func TestActionRegistry_ConcurrentGet(t *testing.T)
func TestActionRegistry_ConcurrentUnregister(t *testing.T)
func TestActionRegistry_MixedOperations(t *testing.T)
```

**Priorité** : ⭐⭐ Moyenne  
**Estimation** : 1 heure

### B. Validation Unicité Actions ⏱️ 1-2 heures

#### 1. Créer ActionRegistry pour ActionDefinitions

**Fichier à créer** : `constraint/action_registry.go`

**Fonctionnalités** :
- Indexation des `ActionDefinition` par nom
- Validation unicité lors de l'ajout
- Recherche par nom
- Thread-safe

**Code proposé** :
```go
type ActionDefinitionRegistry struct {
	definitions map[string]ActionDefinition
	mu          sync.RWMutex
}

func NewActionDefinitionRegistry() *ActionDefinitionRegistry {
	return &ActionDefinitionRegistry{
		definitions: make(map[string]ActionDefinition),
	}
}

func (r *ActionDefinitionRegistry) Add(def ActionDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.definitions[def.Name]; exists {
		return fmt.Errorf("action '%s' est déjà définie", def.Name)
	}
	
	r.definitions[def.Name] = def
	return nil
}

func (r *ActionDefinitionRegistry) Get(name string) (ActionDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	def, exists := r.definitions[name]
	return def, exists
}
```

**Priorité** : ⭐⭐⭐ Haute  
**Estimation** : 1-2 heures

#### 2. Intégrer avec Parsing

**Fichier à modifier** : `constraint/api.go` (ou similaire)

**Modifications** :
- Créer registry lors du parsing
- Valider unicité des actions
- Retourner erreur si doublon

**Estimation** : 30 minutes

### C. Actions Par Défaut ⏱️ 3-4 heures

#### 1. Action Assert

**Fichier à créer** : `rete/action_assert.go`

**Fonctionnalité** : Ajouter un fait dans le réseau RETE

**Code proposé** :
```go
type AssertAction struct {
	network *ReteNetwork
	logger  *log.Logger
}

func NewAssertAction(network *ReteNetwork, logger *log.Logger) *AssertAction {
	if logger == nil {
		logger = log.Default()
	}
	return &AssertAction{
		network: network,
		logger:  logger,
	}
}

func (a *AssertAction) GetName() string {
	return "assert"
}

func (a *AssertAction) Validate(args []interface{}) error {
	if len(args) != 1 {
		return fmt.Errorf("assert attend exactement 1 argument (fait à ajouter), reçu %d", len(args))
	}
	
	// Vérifier que c'est un Fact
	if _, ok := args[0].(*Fact); !ok {
		return fmt.Errorf("assert attend un fait, reçu %T", args[0])
	}
	
	return nil
}

func (a *AssertAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	fact := args[0].(*Fact)
	
	if err := a.network.AddFact(fact); err != nil {
		return fmt.Errorf("échec ajout fait: %w", err)
	}
	
	a.logger.Printf("✅ ASSERT: Fait %s ajouté", fact.GetInternalID())
	return nil
}
```

**Priorité** : ⭐⭐⭐ Haute  
**Estimation** : 1 heure

#### 2. Action Retract

**Fichier à créer** : `rete/action_retract.go`

**Fonctionnalité** : Retirer un fait du réseau RETE

**Estimation** : 1 heure

#### 3. Action Modify

**Fichier à créer** : `rete/action_modify.go`

**Fonctionnalité** : Modifier un fait existant

**Estimation** : 1-2 heures

#### 4. Action Halt

**Fichier à créer** : `rete/action_halt.go`

**Fonctionnalité** : Arrêter le moteur RETE

**Estimation** : 30 minutes

#### 5. Enregistrement dans ActionExecutor

**Fichier à modifier** : `rete/action_executor.go`

**Modifications** :
```go
func (ae *ActionExecutor) RegisterDefaultActions() {
	// Actions actuelles
	ae.registry.Register(NewPrintAction(nil))
	
	// Actions proposées
	ae.registry.Register(NewAssertAction(ae.network, ae.logger))
	ae.registry.Register(NewRetractAction(ae.network, ae.logger))
	ae.registry.Register(NewModifyAction(ae.network, ae.logger))
	ae.registry.Register(NewHaltAction(ae.network, ae.logger))
}
```

**Estimation** : 15 minutes

---

## 📅 Phase 3 : Architecture Xuples (À PLANIFIER)

### A. Création Package Xuples ⏱️ 1 jour

**Tâches** :
- [ ] Créer `tsd/xuples/` package
- [ ] Définir structures `Xuple`, `XupleSpace`
- [ ] Implémenter indexation multi-critères
- [ ] Définir lifecycle (pending, executing, executed, failed)

### B. Interface TupleSpacePublisher ⏱️ 2-3 heures

**Tâches** :
- [ ] Définir interface `TupleSpacePublisher`
- [ ] Implémenter publication d'activations
- [ ] Ajouter callbacks/events

### C. Intégration RETE ↔ Xuples ⏱️ 2-3 jours

**Tâches** :
- [ ] Modifier `TerminalNode.executeAction()`
- [ ] Remplacer `logTupleSpaceActivation()` par `publisher.Publish()`
- [ ] Ajouter configuration enable/disable tuple-space
- [ ] Migrer `collectActivations()` vers xuples
- [ ] Tests d'intégration

### D. Tests et Validation ⏱️ 2-3 jours

**Tâches** :
- [ ] Tests unitaires xuples (couverture > 90%)
- [ ] Tests d'intégration RETE + xuples
- [ ] Benchmarks de performance
- [ ] Tests de régression complets
- [ ] Documentation complète

---

## 🎯 Priorités Immédiates (Cette Semaine)

### Priorité 1 : Tests ⭐⭐⭐

1. **Tests concurrence** (1-2 heures)
   - `generateTokenID()`
   - `ActionRegistry`

2. **Tests unitaires TerminalNode** (2 heures)
   - Couverture > 80%
   - Cas nominaux et d'erreur

### Priorité 2 : Validation Actions ⭐⭐⭐

1. **ActionDefinitionRegistry** (1-2 heures)
   - Créer registry
   - Validation unicité
   - Tests

### Priorité 3 : Actions Par Défaut ⭐⭐

1. **Assert** (1 heure)
2. **Retract** (1 heure)
3. **Tests** (1 heure)

**Temps total estimé** : 8-10 heures (1-2 jours)

---

## 📊 Progression Globale

| Phase | Tâches | Terminées | En cours | À faire | Progression |
|-------|--------|-----------|----------|---------|-------------|
| **Phase 1** | 6 | 6 | 0 | 0 | 100% ✅ |
| **Phase 2** | 12 | 0 | 0 | 12 | 0% 📋 |
| **Phase 3** | 15 | 0 | 0 | 15 | 0% 📅 |
| **TOTAL** | 33 | 6 | 0 | 27 | 18% |

---

## 🚀 Quick Start pour Développeurs

### Pour commencer Phase 2 :

1. **Créer branche** :
   ```bash
   git checkout -b feature/xuples-phase2
   ```

2. **Commencer par les tests** :
   ```bash
   touch rete/fact_token_concurrent_test.go
   # Implémenter TestGenerateTokenID_Concurrent
   ```

3. **Exécuter** :
   ```bash
   go test ./rete/... -run Concurrent -v
   ```

4. **Valider** :
   ```bash
   make test-complete
   make validate
   ```

---

## 📝 Notes et Observations

### Décisions Techniques

1. **Thread-safety** : Utiliser `sync/atomic` plutôt que `sync.Mutex` pour compteurs simples
2. **Affichage legacy** : Conservé dans Phase 2, sera remplacé en Phase 3
3. **Tests** : Privilégier table-driven tests pour meilleure couverture

### Risques Identifiés

1. **Migration xuples** : Impacte beaucoup de code, nécessite stratégie progressive
2. **Performance** : Actions par défaut pourraient impacter performances (à benchmarker)
3. **Rétrocompatibilité** : Pas requise selon prompt, mais documenter breaking changes

---

## 📚 Références

- [00-INDEX.md](docs/xuples/analysis/00-INDEX.md) - Stratégie complète
- [refactoring-phase1-2025-12-17.md](REPORTS/refactoring-phase1-2025-12-17.md) - Rapport Phase 1
- [common.md](.github/prompts/common.md) - Standards projet
- [review.md](.github/prompts/review.md) - Process revue

---

**Dernière mise à jour** : 2025-12-17 13:15 UTC  
**Prochain review** : Après Phase 2 terminée
