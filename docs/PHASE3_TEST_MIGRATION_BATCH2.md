# Phase 3 - Test Migration Batch 2: Converting 20-30 Tests to TestEnvironment

## Objectif

Convertir 20-30 tests supplémentaires vers le pattern `TestEnvironment` pour améliorer l'isolation, la parallélisation et la maintenabilité des tests.

## Contexte

- **Batch 1 complété**: 31 tests de cohérence convertis + 16 tests unitaires TestEnvironment
- **Pattern établi**: `TestEnvironment` avec options fonctionnelles
- **Documentation**: `LOGGING_GUIDE.md` et exemples disponibles
- **Cible Batch 2**: 20-30 tests critiques dans les modules core (action, builder, network, storage)

## Tests Prioritaires à Convertir

### Priorité 1: Action Execution (5-7 tests)
**Fichier**: `rete/action_executor_test.go`

Tests ciblés:
- `TestActionExecutor_BasicExecution` - Exécution basique d'actions
- `TestActionExecutor_WithBindings` - Actions avec bindings de variables
- `TestActionExecutor_ErrorHandling` - Gestion d'erreurs d'exécution
- `TestActionExecutor_ComplexExpressions` - Expressions complexes dans actions
- Tests additionnels selon la couverture du fichier

**Justification**: Tests critiques pour l'exécution des actions, actuellement utilisent `NewMemoryStorage()` et `NewReteNetwork()` directement.

**Bénéfices**:
- Logger dédié pour tracer l'exécution des actions
- Isolation complète entre tests parallèles
- Cleanup automatique des ressources

---

### Priorité 2: Builder Tests (6-8 tests)
**Fichiers**: 
- `rete/builder_rules_test.go`
- `rete/builder_types_test.go`
- `rete/builder_alpha_rules_test.go`

Tests ciblés:
- Tests de construction de règles simples
- Tests de construction de types
- Tests de validation des structures alpha
- Tests d'erreurs de construction

**Justification**: Le builder est au cœur de la création du réseau RETE, ces tests doivent être robustes et isolés.

**Bénéfices**:
- Validation précise des logs de construction
- Métriques de performance du builder
- Tests parallèles sans interférence

---

### Priorité 3: Network Core (4-6 tests)
**Fichier**: Tests de base du réseau RETE

Tests ciblés:
- Tests de création et initialisation du réseau
- Tests de submission de faits de base
- Tests de propagation dans le réseau
- Tests de gestion du cycle de vie

**Justification**: Tests fondamentaux du moteur RETE.

**Bénéfices**:
- Traçabilité complète des opérations réseau
- Isolation des états du réseau entre tests
- Détection précoce des race conditions

---

### Priorité 4: Storage Operations (3-5 tests)
**Fichier**: Tests du storage (si non déjà convertis)

Tests ciblés:
- Tests CRUD sur les faits
- Tests de récupération de faits
- Tests de gestion des IDs internes
- Tests de cohérence storage-réseau

**Justification**: Le storage est critique pour la persistance et la cohérence.

**Bénéfices**:
- Logger pour auditer les opérations storage
- Tests de race conditions sur accès concurrents
- Validation des garanties transactionnelles

---

### Priorité 5: Evaluation & Conditions (4-6 tests)
**Fichiers**:
- `rete/condition_evaluator_test.go`
- `rete/evaluator_operators_test.go`

Tests ciblés:
- Tests d'évaluation de conditions simples
- Tests d'opérateurs (AND, OR, NOT)
- Tests de comparaisons
- Tests d'expressions arithmétiques

**Justification**: L'évaluateur est utilisé partout dans le moteur.

**Bénéfices**:
- Logs détaillés des évaluations
- Tests parallèles pour performance
- Isolation des contextes d'évaluation

---

## Pattern de Conversion

### Avant (Pattern actuel)
```go
func TestSomething(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    
    // Types
    network.Types = append(network.Types, TypeDefinition{...})
    
    // Facts
    fact := &Fact{...}
    network.AddFact(fact)
    
    // Assertions
    assert.Equal(t, expected, actual)
}
```

### Après (Pattern TestEnvironment)
```go
func TestSomething(t *testing.T) {
    t.Parallel() // Si approprié
    
    env := NewTestEnvironment(t,
        WithLogLevel(LogLevelInfo), // Ou LogLevelSilent pour tests lourds
    )
    defer env.Cleanup()
    
    // Types via l'environment
    env.Network.Types = append(env.Network.Types, TypeDefinition{...})
    
    // Facts via l'environment
    fact := &Fact{...}
    env.Network.AddFact(fact)
    
    // Assertions + vérification des logs
    assert.Equal(t, expected, actual)
    env.AssertNoErrors(t) // Vérifie qu'aucune erreur n'a été loggée
    
    // Optionnel: inspecter les logs
    logs := env.GetLogs()
    assert.Contains(t, logs, "expected log message")
}
```

---

## Checklist de Conversion par Test

Pour chaque test converti:

- [ ] Remplacer `NewMemoryStorage()` + `NewReteNetwork()` par `NewTestEnvironment(t)`
- [ ] Ajouter `defer env.Cleanup()`
- [ ] Remplacer `storage` par `env.Storage`
- [ ] Remplacer `network` par `env.Network`
- [ ] Ajouter `t.Parallel()` si le test est indépendant
- [ ] Choisir le bon `LogLevel`:
  - `LogLevelInfo`: tests avec peu d'opérations (< 10 faits)
  - `LogLevelWarn`: tests moyens (10-50 faits)
  - `LogLevelSilent`: tests lourds (> 50 faits) ou boucles concurrentes
- [ ] Ajouter `env.AssertNoErrors(t)` à la fin si approprié
- [ ] Optionnel: vérifier les logs spécifiques avec `env.GetLogs()`
- [ ] Exécuter le test avec `-race`: `go test -race -run TestName`
- [ ] Vérifier que le test passe et n'a pas de race conditions

---

## Métriques de Succès

### Objectifs Quantitatifs
- **Minimum**: 20 tests convertis
- **Cible**: 25 tests convertis
- **Optimal**: 30 tests convertis

### Objectifs Qualitatifs
- ✅ Tous les tests convertis passent avec `-race`
- ✅ Aucune régression de couverture de code
- ✅ Temps d'exécution de la suite maintenu ou amélioré
- ✅ Logs exploitables pour le debugging
- ✅ Pattern cohérent et maintenable

---

## Estimation de Temps

- **Par test simple**: 5-10 minutes
- **Par test complexe**: 15-20 minutes
- **Documentation et validation**: 30 minutes
- **Total estimé**: 2h30 - 4h00 pour 20-30 tests

---

## Ordre d'Exécution

1. **Phase 1** (30-45 min): Convertir Priorité 1 (Action Execution)
2. **Phase 2** (45-60 min): Convertir Priorité 2 (Builder Tests)
3. **Phase 3** (30-45 min): Convertir Priorité 3 (Network Core)
4. **Phase 4** (20-30 min): Convertir Priorité 4 (Storage Operations)
5. **Phase 5** (30-40 min): Convertir Priorité 5 (Evaluation & Conditions)
6. **Validation finale** (20-30 min): Run complet avec `-race`, vérifier métriques

---

## Livrables

À la fin de cette batch:

1. **Tests convertis**: 20-30 fichiers de test mis à jour
2. **Rapport de conversion**: Liste des tests convertis avec statut
3. **Métriques**: 
   - Nombre total de tests convertis (Batch 1 + Batch 2)
   - Couverture de code maintenue/améliorée
   - Résultats `-race` clean
4. **Commit Git**: Message structuré documentant les changements

---

## Critères de Validation

### Pour chaque test converti:
```bash
# Test individuel avec race detector
go test -race -run TestName ./rete

# Vérifier les logs (pas d'erreurs inattendues)
# Vérifier le passage (exit code 0)
```

### Pour la suite complète:
```bash
# Tous les tests du package rete
go test -race -count=1 ./rete

# Vérifier:
# - Tous les tests passent (PASS)
# - Aucune race condition détectée
# - Temps d'exécution raisonnable (< 5 min pour la suite)
```

---

## Prochaines Étapes Après Batch 2

Une fois ce batch complété:

1. **Batch 3 (optionnel)**: Convertir tests de constraint package
2. **Batch 4 (optionnel)**: Convertir tests d'intégration E2E
3. **Phase 4**: Implémenter modes de cohérence (Strong/Relaxed/Eventual)
4. **Phase 5**: Métriques Prometheus et dashboards Grafana

---

## Notes Importantes

### Race Conditions à Surveiller
- Accès concurrents au logger buffer
- `Network.SetTransaction()` appelé depuis plusieurs goroutines
- Partage d'instances de Network/Storage entre goroutines

### Solutions
- Utiliser `LogLevelSilent` pour tests avec goroutines lourdes
- Créer un `TestEnvironment` par goroutine si nécessaire
- Éviter `t.Parallel()` pour tests qui doivent être séquentiels

### Référence
- **Guide complet**: `docs/LOGGING_GUIDE.md`
- **Exemples**: `rete/coherence_testenv_example_test.go`
- **Tests déjà convertis**: `rete/coherence_test.go`, `rete/coherence_phase2_test.go`

---

**Statut**: 🚀 Prêt à démarrer
**Responsable**: À assigner
**Date de début prévue**: 2025-12-04
**Date de fin estimée**: 2025-12-04 (même journée)