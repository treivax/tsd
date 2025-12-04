# Plan de Refactoring : Migration vers Logger Structuré

**Date** : 2025-12-04  
**Phase** : Phase 3 - Sous-partie 3  
**Estimation** : 2-3 heures  

---

## 🎯 Objectif

Remplacer tous les `tsdio.Printf/Println` par le logger structuré pour :
- Uniformiser l'instrumentation du code
- Ajouter des niveaux de log appropriés (Debug/Info/Warn/Error)
- Enrichir avec du contexte structuré
- Préparer pour la production (logs configurables)

---

## 📊 Audit des Logs Existants

### Fichiers Principaux

| Fichier | Nombre de logs | Priorité |
|---------|---------------|----------|
| `rete/constraint_pipeline.go` | ~40 | 🔥 Haute |
| `rete/network.go` | ~20 | 🔥 Haute |
| `rete/store_base.go` | ~5 | 🟡 Moyenne |
| Autres fichiers | ~10 | 🟢 Basse |
| **TOTAL** | **~75** | - |

---

## 🔄 Stratégie de Migration

### Phase 1 : Ajouter le Logger aux Structures (30 min)

#### 1.1 Ajouter Logger à `ReteNetwork`

```go
type ReteNetwork struct {
    // ... champs existants ...
    logger *Logger  // Nouveau champ
}

func NewReteNetwork(storage Storage) *ReteNetwork {
    return &ReteNetwork{
        // ... init existante ...
        logger: NewLogger(LogLevelInfo),  // Info par défaut
    }
}

// Setter optionnel pour configurer
func (rn *ReteNetwork) SetLogger(logger *Logger) {
    rn.logger = logger
}
```

#### 1.2 Ajouter Logger à `ConstraintPipeline`

```go
type ConstraintPipeline struct {
    logger *Logger  // Nouveau champ
}

func NewConstraintPipeline() *ConstraintPipeline {
    return &ConstraintPipeline{
        logger: NewLogger(LogLevelInfo),
    }
}

func (cp *ConstraintPipeline) SetLogger(logger *Logger) {
    cp.logger = logger
}
```

### Phase 2 : Mapper Niveaux de Log (15 min)

| Type de message | Ancien | Nouveau niveau | Justification |
|-----------------|--------|----------------|---------------|
| Opérations normales | `tsdio.Printf("✅ ...")` | `Info` | Flux principal |
| Détails techniques | `tsdio.Printf("📊 ...")` | `Debug` | Détails internes |
| Problèmes non-bloquants | `tsdio.Printf("⚠️ ...")` | `Warn` | Attention requise |
| Erreurs | `tsdio.Printf("❌ ...")` | `Error` | Échecs critiques |
| Étapes importantes | `tsdio.Printf("🔒 ...")` | `Info` | Jalons |

### Phase 3 : Pattern de Conversion (10 min)

#### Pattern Simple

```go
// AVANT
tsdio.Printf("✅ Parsing réussi\n")

// APRÈS
rn.logger.Info("Parsing réussi").Log()
```

#### Pattern avec Contexte

```go
// AVANT
tsdio.Printf("✅ Fait %s persisté après %d tentatives\n", factID, attempts)

// APRÈS
rn.logger.Info("Fait persisté avec retries").
    WithContext("fact_id", factID).
    WithContext("attempts", attempts).
    Log()
```

#### Pattern avec Emoji (conserver pour lisibilité)

```go
// AVANT
tsdio.Printf("🔥 Soumission d'un nouveau fait au réseau RETE: %s\n", fact.String())

// APRÈS
rn.logger.Debug("🔥 Soumission fait").
    WithContext("fact", fact.String()).
    Log()
```

---

## 📝 Refactoring par Fichier

### Fichier 1 : `rete/network.go` (45 min)

**Logs à convertir** : ~20

#### Zone 1 : `SubmitFact()` - 1 log

```go
// L252 - AVANT
tsdio.Printf("🔥 Soumission d'un nouveau fait au réseau RETE: %s\n", fact.String())

// APRÈS
rn.logger.Debug("🔥 Soumission fait").
    WithContext("fact", fact.String()).
    Log()
```

**Niveau** : Debug (détail interne, haute fréquence)

#### Zone 2 : `waitForFactPersistenceWithMetrics()` - 1 log

```go
// L347 - AVANT
tsdio.Printf("✅ Fait %s persisté après %d tentative(s)\n", fact.ID, attempt)

// APRÈS
rn.logger.Info("Fait persisté avec retries").
    WithContext("fact_id", fact.ID).
    WithContext("attempts", attempt).
    Log()
```

**Niveau** : Info (événement important)

#### Zone 3 : `submitFactsFromGrammarWithMetrics()` - 1 log

```go
// L486 - AVANT
tsdio.Printf("✅ Phase 2 - Synchronisation complète: %d/%d faits persistés en %v\n",
    factsPersisted, factsSubmitted, duration)

// APRÈS
rn.logger.Info("Phase 2 - Synchronisation complète").
    WithContext("facts_persisted", factsPersisted).
    WithContext("facts_submitted", factsSubmitted).
    WithContext("duration", duration.String()).
    Log()
```

**Niveau** : Info (jalon important)

#### Zone 4 : `RetractFact()` - 1 log

```go
// L496 - AVANT
tsdio.Printf("🗑️  Rétractation du fait: %s\n", factID)

// APRÈS
rn.logger.Info("🗑️ Rétractation fait").
    WithContext("fact_id", factID).
    Log()
```

**Niveau** : Info (opération importante)

#### Zone 5 : `Reset()` - 2 logs

```go
// L512 - AVANT
tsdio.Println("🧹 Réinitialisation complète du réseau RETE")

// APRÈS
rn.logger.Info("🧹 Réinitialisation réseau RETE").Log()

// L543 - AVANT
tsdio.Println("✅ Réseau RETE réinitialisé avec succès")

// APRÈS
rn.logger.Info("✅ Réseau RETE réinitialisé").Log()
```

**Niveau** : Info (événement majeur)

#### Zone 6 : `ClearMemory()` - 2 logs

```go
// L548 - AVANT
tsdio.Println("🧹 Nettoyage de la mémoire du réseau RETE")

// APRÈS
rn.logger.Info("🧹 Nettoyage mémoire").Log()

// L581 - AVANT
tsdio.Println("✅ Mémoire du réseau RETE nettoyée avec succès")

// APRÈS
rn.logger.Info("✅ Mémoire nettoyée").Log()
```

**Niveau** : Info

#### Zone 7 : `RemoveRule()` - ~12 logs

```go
// L587 - AVANT
tsdio.Printf("🗑️  Suppression de la règle: %s\n", ruleID)

// APRÈS
rn.logger.Info("🗑️ Suppression règle").
    WithContext("rule_id", ruleID).
    Log()

// L598 - AVANT
tsdio.Printf("   📊 Nœuds associés à la règle: %d\n", len(nodeIDs))

// APRÈS
rn.logger.Debug("Nœuds associés").
    WithContext("rule_id", ruleID).
    WithContext("node_count", len(nodeIDs)).
    Log()

// L616 - AVANT
tsdio.Printf("   🔗 JoinNodes détectés, utilisation de la suppression avec lifecycle\n")

// APRÈS
rn.logger.Debug("JoinNodes détectés - suppression lifecycle").
    WithContext("rule_id", ruleID).
    Log()

// Etc. pour les autres logs de RemoveRule...
```

**Niveaux** :
- Opération principale : Info
- Détails internes : Debug
- Avertissements : Warn

### Fichier 2 : `rete/constraint_pipeline.go` (60 min)

**Logs à convertir** : ~40

#### Zone 1 : `ingestFileWithMetrics()` - En-tête

```go
// L93-94 - AVANT
tsdio.Printf("========================================\n")
tsdio.Printf("📁 Ingestion incrémentale: %s\n", filename)

// APRÈS
cp.logger.Info("========================================").Log()
cp.logger.Info("📁 Ingestion incrémentale").
    WithContext("filename", filename).
    Log()
```

**Niveau** : Info

#### Zone 2 : Parsing

```go
// L103 - AVANT
tsdio.Printf("✅ Parsing réussi\n")

// APRÈS
cp.logger.Info("✅ Parsing réussi").Log()
```

**Niveau** : Info

#### Zone 3 : Reset Detection

```go
// L117 - AVANT
tsdio.Printf("🔄 Commande reset détectée - Réinitialisation complète du réseau\n")

// APRÈS
cp.logger.Info("🔄 Commande reset détectée").Log()

// L126 - AVANT
tsdio.Printf("🗑️  GC du réseau existant...\n")

// APRÈS
cp.logger.Debug("🗑️ GC réseau").Log()

// L128 - AVANT
tsdio.Printf("✅ GC terminé\n")

// APRÈS
cp.logger.Debug("✅ GC terminé").Log()

// L132 - AVANT
tsdio.Printf("🆕 Création d'un nouveau réseau RETE\n")

// APRÈS
cp.logger.Info("🆕 Création nouveau réseau RETE").Log()
```

#### Zone 4 : Transaction

```go
// L144 - AVANT
tsdio.Printf("🔒 Transaction démarrée automatiquement: %s\n", tx.ID)

// APRÈS
cp.logger.Info("🔒 Transaction démarrée").
    WithContext("transaction_id", tx.ID).
    Log()

// L152 - AVANT
tsdio.Printf("❌ Erreur rollback: %v\n", rollbackErr)

// APRÈS
cp.logger.Error("❌ Erreur rollback").
    WithContext("error", rollbackErr.Error()).
    Log()

// L155 - AVANT
tsdio.Printf("🔙 Rollback automatique effectué\n")

// APRÈS
cp.logger.Warn("🔙 Rollback automatique").Log()
```

**Niveaux** :
- Démarrage transaction : Info
- Erreur rollback : Error
- Rollback effectué : Warn (attention requise)

#### Zone 5 : Validation

```go
// L169 - AVANT
tsdio.Printf("✅ Validation sémantique réussie\n")

// APRÈS
cp.logger.Info("✅ Validation sémantique réussie").Log()

// L176 - AVANT
tsdio.Printf("🔍 Validation sémantique incrémentale avec contexte...\n")

// APRÈS
cp.logger.Info("🔍 Validation incrémentale").Log()

// L181 - AVANT
tsdio.Printf("✅ Validation incrémentale réussie (%d types en contexte)\n", len(network.Types))

// APRÈS
cp.logger.Info("✅ Validation incrémentale réussie").
    WithContext("type_count", len(network.Types)).
    Log()
```

#### Zone 6 : Création Réseau

```go
// L199 - AVANT
tsdio.Printf("🆕 Création d'un nouveau réseau RETE\n")

// APRÈS
cp.logger.Info("🆕 Création réseau RETE").Log()

// L201 - AVANT
tsdio.Printf("🔄 Extension du réseau RETE existant\n")

// APRÈS
cp.logger.Info("🔄 Extension réseau RETE").Log()
```

#### Zone 7 : Types et Règles

```go
// L216 - AVANT
tsdio.Printf("✅ Trouvé %d types et %d expressions dans le fichier\n", len(types), len(expressions))

// APRÈS
cp.logger.Info("Composants extraits").
    WithContext("type_count", len(types)).
    WithContext("expression_count", len(expressions)).
    Log()

// L226 - AVANT
tsdio.Printf("✅ Types ajoutés/mis à jour dans le réseau\n")

// APRÈS
cp.logger.Info("✅ Types ajoutés").Log()

// L246 - AVANT
tsdio.Printf("📊 Faits préexistants dans le réseau: %d\n", len(existingFacts))

// APRÈS
cp.logger.Debug("📊 Faits préexistants").
    WithContext("fact_count", len(existingFacts)).
    Log()

// L252 - AVANT
tsdio.Printf("📊 Réseau réinitialisé - pas de faits préexistants\n")

// APRÈS
cp.logger.Debug("📊 Réseau réinitialisé").Log()

// L269 - AVANT
tsdio.Printf("✅ Règles ajoutées au réseau\n")

// APRÈS
cp.logger.Info("✅ Règles ajoutées").Log()
```

#### Zone 8 : Propagation

```go
// L286 - AVANT
tsdio.Printf("🔄 Propagation ciblée de faits vers %d nouvelle(s) règle(s)\n", len(newTerminals))

// APRÈS
cp.logger.Info("🔄 Propagation ciblée").
    WithContext("new_terminal_count", len(newTerminals)).
    Log()

// L298 - AVANT
tsdio.Printf("✅ Propagation rétroactive terminée (%d fait(s) propagé(s))\n", propagatedCount)

// APRÈS
cp.logger.Info("✅ Propagation terminée").
    WithContext("propagated_count", propagatedCount).
    Log()
```

#### Zone 9 : Soumission Faits

```go
// Continue avec tous les autres logs de manière similaire...
```

### Fichier 3 : `rete/store_base.go` (15 min)

**Logs à convertir** : ~5

Suivre le même pattern que les fichiers précédents.

---

## 🧪 Tests et Validation (30 min)

### 1. Tester avec Logger Silent

```go
func TestWithSilentLogger(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    network.SetLogger(NewLogger(LogLevelSilent))
    
    // Les tests existants devraient passer sans sortie
    // ...
}
```

### 2. Tester avec Logger Debug

```go
func TestWithDebugLogger(t *testing.T) {
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    network.SetLogger(NewLogger(LogLevelDebug))
    
    // Doit afficher tous les logs
    // ...
}
```

### 3. Valider Tous les Tests Existants

```bash
# Tous les tests doivent passer avec logger par défaut (Info)
go test -race ./rete/... -v

# Tous les tests doivent passer en mode silent
RETE_LOG_LEVEL=silent go test -race ./rete/... -v
```

---

## ✅ Checklist de Migration

### Préparation
- [ ] Ajouter champ `logger *Logger` à `ReteNetwork`
- [ ] Ajouter champ `logger *Logger` à `ConstraintPipeline`
- [ ] Ajouter méthodes `SetLogger()` aux deux structures
- [ ] Initialiser avec `LogLevelInfo` par défaut

### Migration `rete/network.go`
- [ ] Zone 1 : SubmitFact (1 log)
- [ ] Zone 2 : waitForFactPersistence (1 log)
- [ ] Zone 3 : submitFactsFromGrammar (1 log)
- [ ] Zone 4 : RetractFact (1 log)
- [ ] Zone 5 : Reset (2 logs)
- [ ] Zone 6 : ClearMemory (2 logs)
- [ ] Zone 7 : RemoveRule (~12 logs)

### Migration `rete/constraint_pipeline.go`
- [ ] Zone 1 : En-tête ingestion (2 logs)
- [ ] Zone 2 : Parsing (1 log)
- [ ] Zone 3 : Reset detection (4 logs)
- [ ] Zone 4 : Transaction (3 logs)
- [ ] Zone 5 : Validation (3 logs)
- [ ] Zone 6 : Création réseau (2 logs)
- [ ] Zone 7 : Types et règles (6 logs)
- [ ] Zone 8 : Propagation (2 logs)
- [ ] Zone 9 : Soumission faits (reste)

### Migration `rete/store_base.go`
- [ ] Tous les logs (~5)

### Tests et Validation
- [ ] Créer test avec logger Silent
- [ ] Créer test avec logger Debug
- [ ] Valider tous les tests existants passent
- [ ] Valider avec race detector
- [ ] Vérifier aucun Printf restant (grep)

### Documentation
- [ ] Mettre à jour README avec configuration logger
- [ ] Documenter niveaux de log recommandés
- [ ] Ajouter exemples d'utilisation

---

## 🚀 Ordre d'Exécution Recommandé

1. **Session 1** (1h30) : Préparation + `rete/network.go`
   - Ajouter loggers aux structures
   - Migrer tous les logs de `network.go`
   - Tester que tout compile

2. **Session 2** (1h) : `rete/constraint_pipeline.go`
   - Migrer tous les logs du pipeline
   - Tester l'ingestion fonctionne

3. **Session 3** (30 min) : Finalisation
   - Migrer `store_base.go` et autres
   - Tests complets
   - Documentation

---

## 📊 Métriques de Succès

- [ ] 0 `tsdio.Printf` restants (grep verification)
- [ ] 100% tests passent avec logger par défaut
- [ ] 100% tests passent en mode Silent
- [ ] Aucune régression détectée
- [ ] Race detector clean
- [ ] Documentation à jour

---

## 🎯 Bénéfices Attendus

1. **Uniformité** : Tous les logs utilisent le même système
2. **Configurabilité** : Niveaux ajustables selon environnement
3. **Contexte structuré** : Facilite le debugging et le monitoring
4. **Production-ready** : Logger peut être redirigé vers systèmes externes
5. **Performance** : Mode Silent = zéro overhead
6. **Maintenabilité** : Pattern cohérent dans toute la codebase

---

**Estimation totale** : 2-3 heures  
**Complexité** : Moyenne (refactoring mécanique mais volumineux)  
**Risque** : Faible (pas de changement de logique, uniquement de présentation)