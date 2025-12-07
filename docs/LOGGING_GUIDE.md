# Logging Guide

Guide complet pour l'utilisation du système de logging dans TSD.

## Table des matières

- [Vue d'ensemble](#vue-densemble)
- [Niveaux de log](#niveaux-de-log)
- [Configuration de base](#configuration-de-base)
- [Utilisation dans les tests](#utilisation-dans-les-tests)
- [Bonnes pratiques](#bonnes-pratiques)
- [Exemples avancés](#exemples-avancés)
- [Dépannage](#dépannage)

---

## Vue d'ensemble

Le système de logging de TSD fournit un logger thread-safe avec plusieurs niveaux de verbosité, configurable et facile à utiliser dans le code de production et les tests.

**Caractéristiques principales :**

- 🔒 **Thread-safe** : Peut être utilisé en toute sécurité depuis plusieurs goroutines
- 🎯 **Niveaux configurables** : Silent, Error, Warn, Info, Debug
- 📝 **Formatage flexible** : Timestamps, préfixes, couleurs optionnels
- 🧪 **Optimisé pour les tests** : Capture des logs et isolation complète

---

## Niveaux de log

Le logger supporte 5 niveaux de log, du plus silencieux au plus verbeux :

### LogLevelSilent (0)

Aucune sortie. Utile pour les tests ou la production silencieuse.

```go
logger := NewLogger(LogLevelSilent, os.Stdout)
logger.Info("Ce message ne sera pas affiché")
```

### LogLevelError (1)

Uniquement les erreurs critiques qui nécessitent une attention immédiate.

```go
logger.Error("❌ Échec de connexion à la base de données: %v", err)
```

**Cas d'usage :**
- Échecs de connexion
- Erreurs de validation critiques
- Panics récupérés
- Corruptions de données

### LogLevelWarn (2)

Avertissements sur des situations potentiellement problématiques.

```go
logger.Warn("⚠️  Cache miss - performance dégradée")
```

**Cas d'usage :**
- Dépréciations
- Configurations sous-optimales
- Ressources approchant des limites
- Fallbacks activés

### LogLevelInfo (3) - Défaut

Informations générales sur le flux d'exécution du programme.

```go
logger.Info("✅ Phase 2 - Synchronisation complète: %d/%d faits persistés", success, total)
```

**Cas d'usage :**
- Démarrage/arrêt de composants
- Résultats d'opérations majeures
- Statistiques de performance
- Transitions d'état importantes

### LogLevelDebug (4)

Informations détaillées pour le débogage et le développement.

```go
logger.Debug("🔍 Vérification fait %s: tentative %d/%d", factID, retry, maxRetries)
```

**Cas d'usage :**
- Traces de flux d'exécution
- Valeurs de variables intermédiaires
- Appels de fonctions internes
- Informations de cache

---

## Configuration de base

### Création d'un logger

```go
import (
    "os"
    "github.com/treivax/tsd/rete"
)

// Logger par défaut (niveau Info, sortie stdout)
logger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)

// Logger silencieux
logger := rete.NewLogger(rete.LogLevelSilent, os.Stdout)

// Logger debug avec buffer pour capture
var buf bytes.Buffer
logger := rete.NewLogger(rete.LogLevelDebug, &buf)
```

### Configuration du logger

```go
// Changer le niveau de log
logger.SetLevel(rete.LogLevelDebug)

// Activer/désactiver les timestamps
logger.SetTimestamps(true)  // 2025-12-04 16:19:30.181 [RETE] [INFO] ...
logger.SetTimestamps(false) // [RETE] [INFO] ...

// Changer la sortie
logger.SetOutput(os.Stderr)

// Personnaliser le préfixe
logger.SetPrefix("MYAPP")  // [MYAPP] [INFO] ...
```

### Utilisation avec les composants

```go
// Configuration du ReteNetwork
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)
network.SetLogger(logger)

// Configuration du ConstraintPipeline
pipeline := rete.NewConstraintPipeline()
pipeline.SetLogger(logger)
```

---

## Utilisation dans les tests

### Approche recommandée : TestEnvironment

Le helper `TestEnvironment` fournit une isolation complète avec capture automatique des logs :

```go
func TestMyFeature(t *testing.T) {
    t.Parallel() // Safe avec TestEnvironment !

    // Créer un environnement isolé
    env := rete.NewTestEnvironment(t,
        rete.WithLogLevel(rete.LogLevelDebug),
        rete.WithTimestamps(false),
    )
    defer env.Cleanup()

    // Utiliser les composants
    env.Network.SubmitFact(fact)

    // Inspecter les logs
    logs := env.GetLogs()
    assert.Contains(t, logs, "✅ Fait persisté")
    
    // Vérifier l'absence d'erreurs
    env.AssertNoErrors(t)
}
```

### Niveaux de log pour les tests

**Tests unitaires rapides :**
```go
env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelSilent))
// Pas de sortie - exécution la plus rapide
```

**Tests de débogage :**
```go
env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelDebug))
logs := env.GetLogs()
t.Logf("Debug logs:\n%s", logs)
```

**Tests d'intégration :**
```go
env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelInfo))
// Logs Info pour tracer le flux sans détails excessifs
```

### Tests concurrents

Pour les tests qui utilisent des goroutines, utilisez un logger silencieux pour éviter les races sur le buffer partagé :

```go
func TestConcurrent(t *testing.T) {
    t.Parallel()

    // Logger silencieux pour éviter les races
    env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelSilent))
    defer env.Cleanup()

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Pas de logs = pas de race sur le buffer
            env.Network.SubmitFact(fact)
        }()
    }
    wg.Wait()
}
```

### Assertions sur les logs

```go
// Vérifier la présence d'un message
logs := env.GetLogs()
assert.Contains(t, logs, "expected message")

// Vérifier plusieurs conditions
assert.Contains(t, logs, "Phase 1")
assert.Contains(t, logs, "Phase 2")
assert.NotContains(t, logs, "ERROR")

// Compter les occurrences
errorCount := strings.Count(logs, "[ERROR]")
assert.Equal(t, 0, errorCount, "Aucune erreur attendue")

// Helper pour vérifier l'absence d'erreurs
env.AssertNoErrors(t) // Échoue si [ERROR] trouvé dans les logs
```

---

## Bonnes pratiques

### 1. Choisir le bon niveau

**❌ Mauvais :**
```go
logger.Info("i=%d", i)  // Trop détaillé pour Info
logger.Error("Cache miss")  // Pas une erreur critique
```

**✅ Bon :**
```go
logger.Debug("Iteration i=%d", i)
logger.Warn("⚠️  Cache miss - performance impact")
```

### 2. Messages informatifs

**❌ Mauvais :**
```go
logger.Info("Done")  // Trop vague
logger.Error("Error: %v", err)  // Pas de contexte
```

**✅ Bon :**
```go
logger.Info("✅ Phase 2 - Synchronisation complète: %d faits persistés", count)
logger.Error("❌ Échec de persistance du fait %s: %v", factID, err)
```

### 3. Utiliser les emojis pour la lisibilité

- ✅ Succès / Opération complétée
- ❌ Erreur / Échec
- ⚠️  Avertissement
- 🔍 Débogage / Inspection
- 🔥 Action importante
- ⚙️  Configuration / Setup
- 📊 Statistiques / Métriques
- 🔒 Sécurité / Verrouillage

### 4. Logging structuré

```go
// Grouper les informations liées
logger.Info("🔥 Démarrage RETE Network")
logger.Info("   Storage: %s", storageType)
logger.Info("   Rules: %d", ruleCount)
logger.Info("   Coherence: %s", coherenceMode)
```

### 5. Éviter les logs excessifs dans les boucles

**❌ Mauvais :**
```go
for i := 0; i < 10000; i++ {
    logger.Debug("Processing item %d", i)  // 10k lignes !
}
```

**✅ Bon :**
```go
logger.Debug("🔍 Traitement de %d items...", len(items))
for i, item := range items {
    if i % 1000 == 0 {
        logger.Debug("   Progression: %d/%d", i, len(items))
    }
    // Process item
}
logger.Info("✅ %d items traités", len(items))
```

### 6. Logging dans les tests

```go
// Pour les tests normaux : Silent ou Info
env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelInfo))

// Pour déboguer UN test spécifique : Debug + affichage
if testing.Verbose() {
    env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelDebug))
    defer func() {
        t.Logf("Logs:\n%s", env.GetLogs())
    }()
}
```

---

## Exemples avancés

### Logger avec rotation de fichier

```go
import (
    "log"
    "os"
)

file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    log.Fatal(err)
}
defer file.Close()

logger := rete.NewLogger(rete.LogLevelInfo, file)
logger.SetTimestamps(true)
```

### Logger multiple (stdout + fichier)

```go
import "io"

file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
multiWriter := io.MultiWriter(os.Stdout, file)

logger := rete.NewLogger(rete.LogLevelInfo, multiWriter)
```

### Logger conditionnel (dev vs prod)

```go
var logLevel rete.LogLevel
if os.Getenv("ENV") == "production" {
    logLevel = rete.LogLevelInfo
} else {
    logLevel = rete.LogLevelDebug
}

logger := rete.NewLogger(logLevel, os.Stdout)
```

### Capture et analyse des logs

```go
var buf bytes.Buffer
logger := rete.NewLogger(rete.LogLevelInfo, &buf)

// ... opérations ...

logs := buf.String()
errorCount := strings.Count(logs, "[ERROR]")
warnCount := strings.Count(logs, "[WARN]")

fmt.Printf("Statistiques: %d erreurs, %d avertissements\n", errorCount, warnCount)
```

### Logger par composant

```go
// Logger spécialisé pour chaque composant
networkLogger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)
networkLogger.SetPrefix("NETWORK")

storageLogger := rete.NewLogger(rete.LogLevelDebug, os.Stdout)
storageLogger.SetPrefix("STORAGE")

network.SetLogger(networkLogger)
storage.SetLogger(storageLogger)
```

---

## Dépannage

### Problème : Pas de logs affichés

**Vérifications :**
1. Le niveau de log est-il supérieur à `LogLevelSilent` ?
2. La sortie est-elle correctement configurée ?
3. Les timestamps sont-ils désactivés dans les tests ?

```go
// Diagnostic
logger.SetLevel(rete.LogLevelDebug)
logger.SetOutput(os.Stdout)
logger.Info("Test de logging")  // Devrait s'afficher
```

### Problème : Trop de logs

**Solutions :**
```go
// Réduire le niveau de verbosité
logger.SetLevel(rete.LogLevelInfo)  // Au lieu de Debug

// Ou désactiver complètement
logger.SetLevel(rete.LogLevelSilent)
```

### Problème : Race conditions dans les tests

**Cause :** Plusieurs goroutines écrivent dans le même buffer de log.

**Solution 1 :** Logger silencieux
```go
env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelSilent))
```

**Solution 2 :** Environnements séparés
```go
for i := 0; i < 10; i++ {
    go func() {
        // Chaque goroutine a son propre environnement
        env := rete.NewTestEnvironment(t, rete.WithLogLevel(rete.LogLevelSilent))
        defer env.Cleanup()
        // ...
    }()
}
```

### Problème : Logs trop verbeux dans CI

```go
// Dans les tests CI, utiliser Silent par défaut
func getTestLogLevel() rete.LogLevel {
    if os.Getenv("CI") == "true" {
        return rete.LogLevelSilent
    }
    return rete.LogLevelInfo
}

env := rete.NewTestEnvironment(t, rete.WithLogLevel(getTestLogLevel()))
```

---

## Statistiques de logging (Production)

Analyse effectuée sur la base de code actuelle :

**Répartition par niveau :**
- Info : 54% (99 appels)
- Debug : 27% (49 appels)
- Warn : 18% (32 appels)
- Error : 4% (8 appels)

**Conclusion :** La répartition est appropriée avec une majorité de logs informatifs, suffisamment de debug pour le développement, et des erreurs/warnings en proportions raisonnables.

---

## Ressources complémentaires

- [PHASE3_ACTION_PLAN.md](./PHASE3_ACTION_PLAN.md) - Plan d'action Phase 3
- [LOGGING_STANDARDIZATION_REPORT.md](./LOGGING_STANDARDIZATION_REPORT.md) - Rapport de standardisation
- [test_environment.go](./rete/test_environment.go) - Code source du helper de test
- [test_environment_test.go](./rete/test_environment_test.go) - Tests unitaires du helper

---

## Changelog

### 2025-12-04
- ✅ Création du guide de logging
- ✅ Documentation des niveaux et bonnes pratiques
- ✅ Exemples d'utilisation avec TestEnvironment
- ✅ Section dépannage et FAQ

---

**Maintenu par :** TSD Contributors  
**Dernière mise à jour :** 2025-12-04