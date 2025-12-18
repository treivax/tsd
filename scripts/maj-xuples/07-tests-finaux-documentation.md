# 🔧 Prompt 07 - Tests Finaux et Documentation

---

## 🎯 Objectif

**Finaliser le projet avec une suite de tests exhaustive, une documentation complète et de haute qualité, et préparer la release finale** après l'implémentation complète de l'automatisation des xuples.

### Contexte

À ce stade du projet :
- ✅ Parser complet avec faits inline et références (Prompt 01)
- ✅ Package `api` fonctionnel et unifié (Prompt 02)
- ✅ Xuple-spaces créés automatiquement (Prompt 03)
- ✅ Actions Xuple automatiques (Prompt 04)
- ✅ Tests E2E migrés (Prompt 05)
- ✅ Code nettoyé et refactoré (Prompt 06)

**Maintenant** : Assurer la qualité finale avant release :
- Tests exhaustifs (unitaires, intégration, E2E, performance)
- Documentation complète (guides, tutoriels, API reference)
- Exemples d'utilisation concrets
- Validation de la couverture de tests
- Préparation de la release

L'objectif est de **livrer un produit de qualité production** avec une documentation exemplaire et une confiance totale dans le code.

### Prérequis

- ✅ Prompts 01-06 complétés
- ✅ Code nettoyé et refactoré
- ✅ Tous les tests passent

### Résultat Attendu Final

**Une release complète comprenant :**

1. **Couverture de tests > 85%** sur tout le code
2. **Documentation exhaustive** (guides, tutoriels, API reference)
3. **Exemples concrets** d'utilisation (IoT, e-commerce, monitoring)
4. **Benchmarks** et rapports de performance
5. **Guide de contribution** pour les développeurs
6. **Release notes** détaillées
7. **CI/CD** configuré et fonctionnel

---

## 📋 Analyse Préliminaire

### 1. État Actuel de la Couverture de Tests

**Commande d'analyse :**

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | grep total
```

**Objectifs de couverture par package :**

| Package                  | Objectif | Actuel | Action |
|--------------------------|----------|--------|--------|
| `api`                    | > 90%    | ?      | Tests d'intégration |
| `internal/rete`          | > 85%    | ?      | Tests unitaires |
| `internal/constraint`    | > 85%    | ?      | Tests parser |
| `xuples`                 | > 90%    | ?      | Tests complets |
| **Total Projet**         | **> 85%**| ?      | **Combler les gaps** |

### 2. Documentation Actuelle

**Inventaire :**

```
docs/
├── ARCHITECTURE.md           # ✅ Existe, à compléter
├── TSD_LANGUAGE.md           # ✅ Existe, à compléter
├── API_USAGE.md              # ✅ Existe, à compléter
├── XUPLES_E2E_AUTOMATIC.md   # ✅ Existe, à nettoyer
├── TESTING.md                # ❓ À créer
├── CONTRIBUTING.md           # ❓ À créer
├── EXAMPLES.md               # ❓ À créer
├── PERFORMANCE.md            # ❓ À créer
└── FAQ.md                    # ❓ À créer
```

### 3. Exemples à Créer

**Liste d'exemples concrets :**

1. **IoT - Surveillance de capteurs**
2. **E-commerce - Gestion de commandes**
3. **Monitoring - Alertes système**
4. **Workflow - Orchestration de tâches**
5. **Gaming - Système d'événements**

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Tests Unitaires Exhaustifs

**Objectif :** Atteindre > 85% de couverture sur tous les packages.

#### 1.1 Tests du Package `api`

**Fichier :** `tsd/api/pipeline_comprehensive_test.go`

```go
package api

import (
    "os"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/resinsec/tsd/xuples"
)

// TestPipeline_NewPipeline teste la création d'un pipeline avec config par défaut.
func TestPipeline_NewPipeline(t *testing.T) {
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    require.NotNil(t, pipeline)
    
    assert.NotNil(t, pipeline.config)
    assert.NotNil(t, pipeline.network)
    assert.NotNil(t, pipeline.xupleManager)
}

// TestPipeline_NewPipelineWithConfig teste la création avec config personnalisée.
func TestPipeline_NewPipelineWithConfig(t *testing.T) {
    config := &Config{
        LogLevel:         LogLevelDebug,
        EnableMetrics:    true,
        MaxFactsInMemory: 10000,
        XupleSpaceDefaults: XupleSpaceDefaults{
            Selection:   xuples.SelectionLIFO,
            Consumption: xuples.ConsumptionPerAgent,
        },
    }
    
    pipeline, err := NewPipelineWithConfig(config)
    require.NoError(t, err)
    require.NotNil(t, pipeline)
    
    assert.Equal(t, LogLevelDebug, pipeline.config.LogLevel)
    assert.True(t, pipeline.config.EnableMetrics)
}

// TestPipeline_IngestFile_ValidTSD teste l'ingestion d'un fichier TSD valide.
func TestPipeline_IngestFile_ValidTSD(t *testing.T) {
    tsdContent := `
xuple-space test {
    selection: fifo
}

type T {
    id: string
}

type X {
    id: string
}

rule R {
    when {
        t: T()
    }
    then {
        Xuple("test", X(id: t.id))
    }
}
`
    
    tmpfile, err := os.CreateTemp("", "test*.tsd")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())
    
    _, err = tmpfile.WriteString(tsdContent)
    require.NoError(t, err)
    tmpfile.Close()
    
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(tmpfile.Name())
    require.NoError(t, err)
    require.NotNil(t, result)
    
    // Vérifications
    assert.Equal(t, 1, result.TypeCount())
    assert.Equal(t, 1, result.RuleCount())
    assert.Equal(t, 1, result.XupleSpaceCount())
}

// TestPipeline_IngestFile_FileNotFound teste l'erreur si fichier inexistant.
func TestPipeline_IngestFile_FileNotFound(t *testing.T) {
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    _, err = pipeline.IngestFile("/nonexistent/file.tsd")
    require.Error(t, err)
    assert.Contains(t, err.Error(), "no such file")
}

// TestPipeline_IngestString teste l'ingestion depuis une chaîne.
func TestPipeline_IngestString(t *testing.T) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
`
    
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestString(tsdContent)
    require.NoError(t, err)
    require.NotNil(t, result)
    
    assert.Equal(t, 1, result.XupleSpaceCount())
}

// TestPipeline_Reset teste la réinitialisation du pipeline.
func TestPipeline_Reset(t *testing.T) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
`
    
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestString(tsdContent)
    require.NoError(t, err)
    assert.Equal(t, 1, result.XupleSpaceCount())
    
    // Reset
    err = pipeline.Reset()
    require.NoError(t, err)
    
    // Le pipeline devrait être vide
    result, err = pipeline.IngestString("type Empty { }")
    require.NoError(t, err)
    assert.Equal(t, 0, result.XupleSpaceCount())
}

// TestResult_GetXuples teste la récupération de xuples.
func TestResult_GetXuples(t *testing.T) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`
    
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestString(tsdContent)
    require.NoError(t, err)
    
    // Soumettre des faits
    for i := 1; i <= 3; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id": i,
        })
        result.Network().Assert(fact)
    }
    
    // Récupérer les xuples
    xuples := result.GetXuples("test")
    require.Len(t, xuples, 3)
    
    ids := []int{}
    for _, x := range xuples {
        ids = append(ids, x.Get("id").(int))
    }
    assert.ElementsMatch(t, []int{1, 2, 3}, ids)
}

// TestResult_Retrieve teste la consommation de xuples.
func TestResult_Retrieve(t *testing.T) {
    tsdContent := `
xuple-space test { 
    selection: fifo,
    consumption: once
}
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`
    
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestString(tsdContent)
    require.NoError(t, err)
    
    // Créer 3 xuples
    for i := 1; i <= 3; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id": i,
        })
        result.Network().Assert(fact)
    }
    
    // Consommer le premier (FIFO)
    xuple, err := result.Retrieve("test", nil, "agent-01")
    require.NoError(t, err)
    require.NotNil(t, xuple)
    assert.Equal(t, 1, xuple.Get("id"))
    
    // Il reste 2 xuples
    xuples := result.GetXuples("test")
    assert.Len(t, xuples, 2)
}

// TestResult_Metrics teste les métriques du résultat.
func TestResult_Metrics(t *testing.T) {
    tsdContent := `
xuple-space s1 { selection: fifo }
xuple-space s2 { selection: lifo }
type T1 { id: int }
type T2 { id: int }
rule R1 { when { t: T1() } then { Xuple("s1", T1(id: t.id)) } }
rule R2 { when { t: T2() } then { Xuple("s2", T2(id: t.id)) } }
`
    
    pipeline, err := NewPipeline()
    require.NoError(t, err)
    
    start := time.Now()
    result, err := pipeline.IngestString(tsdContent)
    require.NoError(t, err)
    duration := time.Since(start)
    
    metrics := result.Metrics()
    
    assert.Equal(t, 2, metrics.TypeCount)
    assert.Equal(t, 2, metrics.RuleCount)
    assert.Equal(t, 2, metrics.XupleSpaceCount)
    assert.Greater(t, metrics.TotalDuration, time.Duration(0))
    assert.LessOrEqual(t, metrics.TotalDuration, duration*2) // Marge
}

// TestConfig_Validate teste la validation de la configuration.
func TestConfig_Validate(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
        errMsg  string
    }{
        {
            name: "valid config",
            config: &Config{
                LogLevel:         LogLevelInfo,
                MaxFactsInMemory: 1000,
            },
            wantErr: false,
        },
        {
            name: "invalid log level",
            config: &Config{
                LogLevel: 999,
            },
            wantErr: true,
            errMsg:  "invalid log level",
        },
        {
            name: "negative max facts",
            config: &Config{
                MaxFactsInMemory: -1,
            },
            wantErr: true,
            errMsg:  "MaxFactsInMemory must be >= 0",
        },
        {
            name: "invalid selection policy",
            config: &Config{
                XupleSpaceDefaults: XupleSpaceDefaults{
                    Selection: 999,
                },
            },
            wantErr: true,
            errMsg:  "invalid selection policy",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.config.Validate()
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

#### 1.2 Tests du Package `xuples`

**Fichier :** `tsd/xuples/comprehensive_test.go`

```go
package xuples

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestXupleManager_CreateXupleSpace teste la création de xuple-spaces.
func TestXupleManager_CreateXupleSpace(t *testing.T) {
    manager := NewXupleManager()
    
    space, err := manager.CreateXupleSpace(
        "test",
        SelectionFIFO,
        ConsumptionOnce,
        RetentionUnlimited,
    )
    
    require.NoError(t, err)
    require.NotNil(t, space)
    assert.Equal(t, "test", space.Name())
}

// TestXupleManager_CreateDuplicateSpace teste l'erreur de duplication.
func TestXupleManager_CreateDuplicateSpace(t *testing.T) {
    manager := NewXupleManager()
    
    _, err := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    require.NoError(t, err)
    
    // Créer à nouveau devrait échouer
    _, err = manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    require.Error(t, err)
    assert.Contains(t, err.Error(), "already exists")
}

// TestXupleSpace_Deposit teste le dépôt de xuples.
func TestXupleSpace_Deposit(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    
    xuple := NewXuple("TestType", map[string]interface{}{
        "id":   1,
        "name": "test",
    })
    
    err := space.Deposit(xuple)
    require.NoError(t, err)
    
    xuples := space.GetAll()
    assert.Len(t, xuples, 1)
}

// TestXupleSpace_Retrieve_FIFO teste la récupération FIFO.
func TestXupleSpace_Retrieve_FIFO(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    
    // Déposer 3 xuples
    for i := 1; i <= 3; i++ {
        xuple := NewXuple("T", map[string]interface{}{"id": i})
        space.Deposit(xuple)
    }
    
    // Récupérer (FIFO = premier déposé)
    xuple, err := space.Retrieve(nil, "agent")
    require.NoError(t, err)
    require.NotNil(t, xuple)
    assert.Equal(t, 1, xuple.Get("id"))
}

// TestXupleSpace_Retrieve_LIFO teste la récupération LIFO.
func TestXupleSpace_Retrieve_LIFO(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionLIFO, ConsumptionOnce, RetentionUnlimited)
    
    for i := 1; i <= 3; i++ {
        xuple := NewXuple("T", map[string]interface{}{"id": i})
        space.Deposit(xuple)
    }
    
    // Récupérer (LIFO = dernier déposé)
    xuple, err := space.Retrieve(nil, "agent")
    require.NoError(t, err)
    require.NotNil(t, xuple)
    assert.Equal(t, 3, xuple.Get("id"))
}

// TestXupleSpace_Retrieve_WithTemplate teste la récupération avec template.
func TestXupleSpace_Retrieve_WithTemplate(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    
    // Déposer divers xuples
    space.Deposit(NewXuple("T", map[string]interface{}{"id": 1, "status": "active"}))
    space.Deposit(NewXuple("T", map[string]interface{}{"id": 2, "status": "inactive"}))
    space.Deposit(NewXuple("T", map[string]interface{}{"id": 3, "status": "active"}))
    
    // Template: status = "active"
    template := map[string]interface{}{"status": "active"}
    
    xuple, err := space.Retrieve(template, "agent")
    require.NoError(t, err)
    require.NotNil(t, xuple)
    assert.Equal(t, 1, xuple.Get("id"))
    assert.Equal(t, "active", xuple.Get("status"))
}

// TestXupleSpace_ConsumptionOnce teste la politique de consommation "once".
func TestXupleSpace_ConsumptionOnce(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    
    space.Deposit(NewXuple("T", map[string]interface{}{"id": 1}))
    
    // Premier agent récupère
    xuple1, err := space.Retrieve(nil, "agent-01")
    require.NoError(t, err)
    require.NotNil(t, xuple1)
    
    // Deuxième agent ne devrait rien récupérer (consommé)
    xuple2, err := space.Retrieve(nil, "agent-02")
    require.NoError(t, err)
    assert.Nil(t, xuple2)
}

// TestXupleSpace_ConsumptionPerAgent teste la politique "per-agent".
func TestXupleSpace_ConsumptionPerAgent(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionPerAgent, RetentionUnlimited)
    
    space.Deposit(NewXuple("T", map[string]interface{}{"id": 1}))
    
    // Agent 1 récupère
    xuple1, err := space.Retrieve(nil, "agent-01")
    require.NoError(t, err)
    require.NotNil(t, xuple1)
    
    // Agent 2 peut aussi récupérer
    xuple2, err := space.Retrieve(nil, "agent-02")
    require.NoError(t, err)
    require.NotNil(t, xuple2)
    
    // Agent 1 ne peut plus récupérer (déjà consommé par lui)
    xuple3, err := space.Retrieve(nil, "agent-01")
    require.NoError(t, err)
    assert.Nil(t, xuple3)
}

// TestXupleSpace_RetentionDuration teste la rétention par durée.
func TestXupleSpace_RetentionDuration(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionDuration)
    space.SetRetentionDuration(100 * time.Millisecond)
    
    space.Deposit(NewXuple("T", map[string]interface{}{"id": 1}))
    
    // Devrait être disponible immédiatement
    assert.Len(t, space.GetAll(), 1)
    
    // Attendre l'expiration
    time.Sleep(150 * time.Millisecond)
    
    // Devrait être expiré
    assert.Len(t, space.GetAll(), 0)
}

// TestXupleSpace_MaxSize teste la limitation de taille.
func TestXupleSpace_MaxSize(t *testing.T) {
    manager := NewXupleManager()
    space, _ := manager.CreateXupleSpace("test", SelectionFIFO, ConsumptionOnce, RetentionUnlimited)
    space.SetMaxSize(3)
    
    // Déposer 5 xuples (max = 3)
    for i := 1; i <= 5; i++ {
        err := space.Deposit(NewXuple("T", map[string]interface{}{"id": i}))
        require.NoError(t, err)
    }
    
    // Seulement 3 devraient être conservés
    xuples := space.GetAll()
    assert.Len(t, xuples, 3)
}

// TestXuple_Get teste l'accès aux champs.
func TestXuple_Get(t *testing.T) {
    xuple := NewXuple("TestType", map[string]interface{}{
        "id":     123,
        "name":   "test",
        "active": true,
        "score":  98.5,
    })
    
    assert.Equal(t, 123, xuple.Get("id"))
    assert.Equal(t, "test", xuple.Get("name"))
    assert.Equal(t, true, xuple.Get("active"))
    assert.Equal(t, 98.5, xuple.Get("score"))
    assert.Nil(t, xuple.Get("nonexistent"))
}

// TestXuple_Type teste le type du xuple.
func TestXuple_Type(t *testing.T) {
    xuple := NewXuple("MyType", map[string]interface{}{})
    assert.Equal(t, "MyType", xuple.Type())
}
```

---

### Tâche 2: Tests de Performance et Benchmarks

**Fichier :** `tsd/test/benchmark/comprehensive_bench_test.go`

```go
package benchmark

import (
    "fmt"
    "testing"
    
    "github.com/resinsec/tsd/api"
)

// BenchmarkPipeline_IngestFile mesure la performance d'ingestion.
func BenchmarkPipeline_IngestFile(b *testing.B) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int, value: float }
type X { id: int, value: float }
rule R { when { t: T() } then { Xuple("test", X(id: t.id, value: t.value)) } }
`
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        pipeline, _ := api.NewPipeline()
        _, _ = pipeline.IngestString(tsdContent)
    }
}

// BenchmarkXupleCreation_Small mesure la création de xuples (petit volume).
func BenchmarkXupleCreation_Small(b *testing.B) {
    benchmarkXupleCreation(b, 100)
}

// BenchmarkXupleCreation_Medium mesure la création (volume moyen).
func BenchmarkXupleCreation_Medium(b *testing.B) {
    benchmarkXupleCreation(b, 1000)
}

// BenchmarkXupleCreation_Large mesure la création (grand volume).
func BenchmarkXupleCreation_Large(b *testing.B) {
    benchmarkXupleCreation(b, 10000)
}

func benchmarkXupleCreation(b *testing.B, count int) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        for j := 0; j < count; j++ {
            fact := result.Network().CreateFact("T", map[string]interface{}{
                "id": j,
            })
            result.Network().Assert(fact)
        }
    }
    
    b.ReportMetric(float64(count*b.N), "xuples")
}

// BenchmarkXupleRetrieval_FIFO mesure la récupération FIFO.
func BenchmarkXupleRetrieval_FIFO(b *testing.B) {
    benchmarkRetrieval(b, "fifo")
}

// BenchmarkXupleRetrieval_LIFO mesure la récupération LIFO.
func BenchmarkXupleRetrieval_LIFO(b *testing.B) {
    benchmarkRetrieval(b, "lifo")
}

// BenchmarkXupleRetrieval_Random mesure la récupération aléatoire.
func BenchmarkXupleRetrieval_Random(b *testing.B) {
    benchmarkRetrieval(b, "random")
}

func benchmarkRetrieval(b *testing.B, policy string) {
    tsdContent := fmt.Sprintf(`
xuple-space test { selection: %s }
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`, policy)
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    // Pré-remplir avec 1000 xuples
    for i := 0; i < 1000; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id": i,
        })
        result.Network().Assert(fact)
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        result.Retrieve("test", nil, fmt.Sprintf("agent-%d", i))
    }
}

// BenchmarkComplexRules mesure les performances avec règles complexes.
func BenchmarkComplexRules(b *testing.B) {
    tsdContent := `
xuple-space alerts { selection: fifo }
xuple-space logs { selection: lifo }

type SensorReading {
    sensorId: string,
    temp: float,
    humidity: float
}

type Alert {
    sensorId: string,
    type: string,
    value: float
}

type Log {
    message: string
}

rule HighTemp {
    when {
        s: SensorReading(temp > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: s.sensorId,
            type: "temperature",
            value: s.temp
        )),
        Xuple("logs", Log(message: "High temp detected"))
    }
}

rule HighHumidity {
    when {
        s: SensorReading(humidity > 80.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: s.sensorId,
            type: "humidity",
            value: s.humidity
        ))
    }
}
`
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        fact := result.Network().CreateFact("SensorReading", map[string]interface{}{
            "sensorId": fmt.Sprintf("sensor-%d", i),
            "temp":     35.0,
            "humidity": 85.0,
        })
        result.Network().Assert(fact)
    }
}
```

---

### Tâche 3: Documentation Complète

#### 3.1 Guide Utilisateur Complet

**Fichier :** `docs/USER_GUIDE.md`

```markdown
# Guide Utilisateur TSD - Xuples E2E

## Introduction

Ce guide vous accompagne dans l'utilisation complète du système TSD avec
l'automatisation des xuples, de l'installation jusqu'à des cas d'usage avancés.

## Installation

### Prérequis

- Go 1.21 ou supérieur
- Git

### Installation via go get

\`\`\`bash
go get github.com/resinsec/tsd
\`\`\`

### Compilation depuis les sources

\`\`\`bash
git clone https://github.com/resinsec/tsd.git
cd tsd
go build ./...
go test ./...
\`\`\`

## Premiers Pas

### Exemple Minimal

**Fichier :** `hello.tsd`

\`\`\`tsd
xuple-space greetings {
    selection: fifo,
    consumption: once
}

type Person {
    name: string
}

type Greeting {
    message: string
}

rule SayHello {
    when {
        p: Person()
    }
    then {
        Xuple("greetings", Greeting(
            message: "Hello, " + p.name + "!"
        ))
    }
}
\`\`\`

**Code Go :**

\`\`\`go
package main

import (
    "fmt"
    "log"
    
    "github.com/resinsec/tsd/api"
)

func main() {
    // Créer le pipeline
    pipeline, err := api.NewPipeline()
    if err != nil {
        log.Fatal(err)
    }
    
    // Charger le fichier TSD
    result, err := pipeline.IngestFile("hello.tsd")
    if err != nil {
        log.Fatal(err)
    }
    
    // Soumettre un fait
    person := result.Network().CreateFact("Person", map[string]interface{}{
        "name": "Alice",
    })
    result.Network().Assert(person)
    
    // Récupérer les xuples
    xuples := result.GetXuples("greetings")
    for _, x := range xuples {
        fmt.Println(x.Get("message"))
    }
}
\`\`\`

**Sortie :**
\`\`\`
Hello, Alice!
\`\`\`

## Concepts Fondamentaux

### Xuple-Spaces

Les xuple-spaces sont des espaces de stockage configurables pour les xuples.

#### Propriétés

| Propriété           | Type      | Valeurs              | Défaut      |
|---------------------|-----------|----------------------|-------------|
| `selection`         | string    | fifo, lifo, random   | fifo        |
| `consumption`       | string    | once, per-agent      | once        |
| `retention`         | duration  | unlimited, durée     | unlimited   |
| `max-size`          | int       | 0 (illimité) ou > 0  | 0           |

#### Exemples

\`\`\`tsd
// FIFO, consommation unique, conservation 24h
xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

// LIFO, consommation par agent, max 1000 xuples
xuple-space logs {
    selection: lifo,
    consumption: per-agent,
    max-size: 1000
}

// Random, défauts pour le reste
xuple-space events {
    selection: random
}
\`\`\`

### Types

Définissez vos structures de données :

\`\`\`tsd
type Temperature {
    sensorId: string,
    value: float,
    unit: string,
    timestamp: int
}

type Alert {
    sensorId: string,
    severity: string,
    message: string
}
\`\`\`

### Règles

Les règles lient conditions et actions :

\`\`\`tsd
rule HighTemperature {
    when {
        t: Temperature(value > 30.0, unit == "celsius")
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            severity: "warning",
            message: "Temperature above threshold"
        ))
    }
}
\`\`\`

### Action Xuple

L'action `Xuple()` crée un xuple dans un xuple-space :

\`\`\`tsd
Xuple("<space-name>", FactType(...))
\`\`\`

#### Syntaxe Avancée

**Références aux champs :**

\`\`\`tsd
Xuple("space", Alert(
    sensorId: t.sensorId,        // Référence au champ
    temp: t.value,                // Référence
    message: "Static message"     // Valeur littérale
))
\`\`\`

**Actions multiples :**

\`\`\`tsd
rule MultiAction {
    when {
        e: Event(severity == "critical")
    }
    then {
        Xuple("alerts", Alert(eventId: e.id, level: "critical")),
        Xuple("logs", LogEntry(source: e.source, message: e.desc))
    }
}
\`\`\`

## API Go

### Pipeline

#### Création

\`\`\`go
// Avec configuration par défaut
pipeline, err := api.NewPipeline()

// Avec configuration personnalisée
config := &api.Config{
    LogLevel: api.LogLevelDebug,
    MaxFactsInMemory: 10000,
}
pipeline, err := api.NewPipelineWithConfig(config)
\`\`\`

#### Ingestion

\`\`\`go
// Depuis un fichier
result, err := pipeline.IngestFile("rules.tsd")

// Depuis une chaîne
result, err := pipeline.IngestString(tsdContent)
\`\`\`

#### Réinitialisation

\`\`\`go
err := pipeline.Reset()
\`\`\`

### Result

Le résultat d'une ingestion :

\`\`\`go
// Accès au réseau RETE
network := result.Network()

// Accès au XupleManager
manager := result.XupleManager()

// Métriques
metrics := result.Metrics()
fmt.Printf("Parsed in %v\n", metrics.ParseDuration)

// Résumé
summary := result.Summary()
fmt.Println(summary)
\`\`\`

#### Gestion des Xuples

\`\`\`go
// Récupérer tous les xuples d'un space
xuples := result.GetXuples("alerts")

// Récupérer (consommer) un xuple
xuple, err := result.Retrieve("alerts", nil, "agent-01")

// Avec template de recherche
template := map[string]interface{}{
    "severity": "critical",
}
xuple, err := result.Retrieve("alerts", template, "agent-01")

// Lister les xuple-spaces
spaces := result.XupleSpaceNames()

// Compter les xuples
count := result.XupleCount("alerts")
\`\`\`

### Soumettre des Faits

\`\`\`go
// Créer un fait
fact := result.Network().CreateFact("Temperature", map[string]interface{}{
    "sensorId": "sensor-01",
    "value":    35.5,
    "unit":     "celsius",
    "timestamp": time.Now().Unix(),
})

// Soumettre au réseau
err := result.Network().Assert(fact)
\`\`\`

## Cas d'Usage

### 1. IoT - Surveillance de Capteurs

**Objectif :** Surveiller des capteurs et générer des alertes.

\`\`\`tsd
xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

xuple-space analytics {
    selection: lifo,
    consumption: per-agent,
    max-size: 10000
}

type SensorReading {
    sensorId: string,
    temperature: float,
    humidity: float,
    battery: int,
    timestamp: int
}

type Alert {
    sensorId: string,
    alertType: string,
    severity: string,
    value: float,
    timestamp: int
}

type AnalyticsEvent {
    sensorId: string,
    eventType: string,
    metadata: string
}

rule HighTemperature {
    when {
        r: SensorReading(temperature > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: r.sensorId,
            alertType: "temperature",
            severity: "warning",
            value: r.temperature,
            timestamp: r.timestamp
        )),
        Xuple("analytics", AnalyticsEvent(
            sensorId: r.sensorId,
            eventType: "high_temp",
            metadata: "Temperature exceeded 30°C"
        ))
    }
}

rule CriticalTemperature {
    when {
        r: SensorReading(temperature > 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: r.sensorId,
            alertType: "temperature",
            severity: "critical",
            value: r.temperature,
            timestamp: r.timestamp
        ))
    }
}

rule LowBattery {
    when {
        r: SensorReading(battery < 20)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: r.sensorId,
            alertType: "battery",
            severity: "warning",
            value: r.battery,
            timestamp: r.timestamp
        ))
    }
}
\`\`\`

### 2. E-commerce - Traitement de Commandes

\`\`\`tsd
xuple-space orders {
    selection: fifo,
    consumption: once
}

xuple-space notifications {
    selection: fifo,
    consumption: per-agent
}

type Order {
    orderId: string,
    customerId: string,
    amount: float,
    status: string
}

type PaymentNotification {
    orderId: string,
    customerId: string,
    amount: float,
    message: string
}

rule ProcessPayment {
    when {
        o: Order(status == "pending_payment", amount > 0.0)
    }
    then {
        Xuple("notifications", PaymentNotification(
            orderId: o.orderId,
            customerId: o.customerId,
            amount: o.amount,
            message: "Payment required"
        ))
    }
}
\`\`\`

## Configuration Avancée

### Options du Pipeline

\`\`\`go
config := &api.Config{
    // Niveau de log (Silent, Error, Warn, Info, Debug)
    LogLevel: api.LogLevelInfo,
    
    // Activer les métriques de performance
    EnableMetrics: true,
    
    // Limite de faits en mémoire
    MaxFactsInMemory: 100000,
    
    // Valeurs par défaut pour les xuple-spaces
    XupleSpaceDefaults: api.XupleSpaceDefaults{
        Selection:   xuples.SelectionFIFO,
        Consumption: xuples.ConsumptionOnce,
        Retention:   xuples.RetentionUnlimited,
    },
}
\`\`\`

## Bonnes Pratiques

### 1. Nommage

- **Xuple-spaces** : Pluriel, descriptif (`alerts`, `notifications`)
- **Types** : PascalCase (`SensorReading`, `Alert`)
- **Champs** : camelCase (`sensorId`, `timestamp`)
- **Règles** : Descriptif de l'action (`HighTemperature`, `ProcessPayment`)

### 2. Organisation

- Un fichier TSD par domaine métier
- Grouper les xuple-spaces liés
- Documenter les règles complexes avec des commentaires

### 3. Performance

- Limiter la taille des xuple-spaces avec `max-size`
- Utiliser `retention` pour nettoyer automatiquement
- Préférer `consumption: once` si applicable

### 4. Tests

- Tester chaque règle individuellement
- Utiliser `testutil.CreatePipelineFromTSD()` pour les tests
- Vérifier les métriques de performance

## Dépannage

### Erreurs Courantes

**Erreur : "xuple-space 'X' does not exist"**

→ Vérifier que le xuple-space est défini dans le fichier TSD

**Erreur : "parsing file: syntax error at line X"**

→ Vérifier la syntaxe TSD (virgules, accolades)

**Performance dégradée**

→ Vérifier la taille des xuple-spaces, utiliser `max-size`

### Debug

Activer les logs détaillés :

\`\`\`go
config := &api.Config{
    LogLevel: api.LogLevelDebug,
}
pipeline, _ := api.NewPipelineWithConfig(config)
\`\`\`

## Ressources

- [API Reference](API_REFERENCE.md)
- [TSD Language Specification](TSD_LANGUAGE.md)
- [Examples](../examples/)
- [FAQ](FAQ.md)
```

#### 3.2 FAQ

**Fichier :** `docs/FAQ.md`

```markdown
# FAQ - Questions Fréquentes

## Général

### Q: Qu'est-ce qu'un xuple ?

**R:** Un xuple est un tuple enrichi stocké dans un xuple-space. C'est une
structure de données qui peut être déposée, récupérée et consommée selon
des politiques configurables.

### Q: Quelle est la différence entre un fait et un xuple ?

**R:** 
- Un **fait** est une donnée soumise au réseau RETE pour activer des règles
- Un **xuple** est une donnée produite par une règle et stockée dans un xuple-space

Les règles transforment des faits en xuples.

### Q: Dois-je configurer manuellement les xuple-spaces ?

**R:** Non ! Les xuple-spaces sont créés automatiquement à partir des
définitions dans votre fichier TSD. Aucune configuration Go nécessaire.

## Configuration

### Q: Comment choisir entre FIFO, LIFO et Random ?

**R:**
- **FIFO** : Pour traiter les xuples dans l'ordre d'arrivée (ex: files d'attente)
- **LIFO** : Pour traiter les plus récents en premier (ex: stack)
- **Random** : Pour équilibrer la charge entre agents

### Q: Quelle est la différence entre "once" et "per-agent" ?

**R:**
- **once** : Un xuple ne peut être consommé qu'une seule fois, par n'importe quel agent
- **per-agent** : Chaque agent peut consommer le xuple une fois

Utilisez "once" pour des tâches uniques, "per-agent" pour des notifications.

### Q: Comment limiter la taille d'un xuple-space ?

**R:** Utilisez la propriété `max-size` :

\`\`\`tsd
xuple-space logs {
    selection: fifo,
    max-size: 1000  // Garde seulement les 1000 plus récents
}
\`\`\`

## Développement

### Q: Comment tester mes règles ?

**R:** Utilisez les helpers de test :

\`\`\`go
import "github.com/resinsec/tsd/test/testutil"

func TestMyRule(t *testing.T) {
    _, result := testutil.CreatePipelineFromTSD(t, tsdContent)
    testutil.SubmitFact(t, result, "MyType", data)
    assert.Equal(t, 1, result.XupleCount("myspace"))
}
\`\`\`

### Q: Comment débugger une règle qui ne se déclenche pas ?

**R:**
1. Activer les logs debug :
   \`\`\`go
   config := &api.Config{LogLevel: api.LogLevelDebug}
   \`\`\`
2. Vérifier que le type est enregistré
3. Vérifier que les conditions de la règle sont satisfaites
4. Vérifier que le xuple-space existe

### Q: Puis-je utiliser TSD sans xuples ?

**R:** Oui ! Les xuples sont optionnels. Si vous ne définissez pas de
xuple-spaces, vous pouvez utiliser TSD comme moteur de règles classique.

## Performance

### Q: Quelle est la performance attendue ?

**R:** Sur une machine moderne :
- ~10,000 xuples/sec créés
- ~50,000 xuples/sec récupérés (FIFO/LIFO)
- < 1ms latence pour ingestion d'un fichier TSD simple

### Q: Comment optimiser les performances ?

**R:**
1. Limiter `max-size` des xuple-spaces
2. Utiliser `retention` pour nettoyer automatiquement
3. Préférer `consumption: once` si possible
4. Éviter les règles avec conditions trop complexes

### Q: Les xuples sont-ils persistés sur disque ?

**R:** Non, actuellement les xuples sont en mémoire uniquement.
La persistance est prévue dans une version future.

## Erreurs

### Q: "import cycle detected" lors de la compilation

**R:** N'importez jamais directement `internal/rete` ou `internal/constraint`.
Utilisez uniquement le package `api`.

### Q: "action 'Xuple' not found"

**R:** Cela ne devrait plus arriver dans la v2.0+. Vérifiez que vous utilisez
`api.NewPipeline()` et non l'ancienne API.

### Q: Mes xuples disparaissent

**R:** Vérifiez :
1. La politique de rétention (`retention`)
2. La taille maximale (`max-size`)
3. Que les xuples ne sont pas consommés par erreur

## Migration

### Q: Comment migrer depuis la v1.x ?

**R:** Consultez le [Guide de Migration](MIGRATION.md). En résumé :
- Remplacer les imports par `api`
- Supprimer toute configuration de factory
- Utiliser `api.NewPipeline()` comme point d'entrée

### Q: Mon ancien code avec SetXupleSpaceFactory() ne compile plus

**R:** Normal ! Cette méthode est obsolète. Les xuple-spaces sont maintenant
créés automatiquement. Consultez le guide de migration.

## Support

### Q: Où poser mes questions ?

**R:**
- [GitHub Issues](https://github.com/resinsec/tsd/issues)
- [GitHub Discussions](https://github.com/resinsec/tsd/discussions)
- Documentation : `docs/`

### Q: Comment contribuer ?

**R:** Consultez [CONTRIBUTING.md](../CONTRIBUTING.md)

### Q: J'ai trouvé un bug, que faire ?

**R:** Créez une issue GitHub avec :
1. Description du problème
2. Code minimal pour reproduire
3. Comportement attendu vs observé
4. Version de TSD (`go list -m github.com/resinsec/tsd`)
```

---

### Tâche 4: Exemples Concrets

**Créer :** `examples/iot-monitoring/main.go`

```go
package main

import (
    "fmt"
    "log"
    "math/rand"
    "time"
    
    "github.com/resinsec/tsd/api"
)

func main() {
    // Créer le pipeline
    pipeline, err := api.NewPipeline()
    if err != nil {
        log.Fatal(err)
    }
    
    // Charger les règles
    result, err := pipeline.IngestFile("rules.tsd")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("✅ Pipeline initialized")
    fmt.Printf("   - %d xuple-spaces created\n", result.XupleSpaceCount())
    fmt.Printf("   - %d rules loaded\n", result.RuleCount())
    
    // Simuler des lectures de capteurs
    fmt.Println("\n🌡️  Simulating sensor readings...")
    
    sensors := []string{"sensor-01", "sensor-02", "sensor-03"}
    
    for i := 0; i < 10; i++ {
        sensorId := sensors[rand.Intn(len(sensors))]
        temperature := 20.0 + rand.Float64()*30.0
        humidity := 40.0 + rand.Float64()*50.0
        
        // Soumettre la lecture
        reading := result.Network().CreateFact("SensorReading", map[string]interface{}{
            "sensorId":    sensorId,
            "temperature": temperature,
            "humidity":    humidity,
            "timestamp":   time.Now().Unix(),
        })
        
        result.Network().Assert(reading)
        
        fmt.Printf("   [%s] Temp: %.1f°C, Humidity: %.1f%%\n",
            sensorId, temperature, humidity)
        
        time.Sleep(500 * time.Millisecond)
    }
    
    // Afficher les alertes générées
    fmt.Println("\n🚨 Alerts generated:")
    
    alerts := result.GetXuples("alerts")
    for i, alert := range alerts {
        fmt.Printf("   %d. [%s] %s alert: %.1f\n",
            i+1,
            alert.Get("sensorId"),
            alert.Get("severity"),
            alert.Get("value"))
    }
    
    fmt.Printf("\n📊 Total: %d alerts\n", len(alerts))
    
    // Récupérer une alerte (simulation d'un agent de monitoring)
    fmt.Println("\n🤖 Monitoring agent retrieving alerts...")
    
    retrieved, err := result.Retrieve("alerts", nil, "monitoring-agent-01")
    if err != nil {
        log.Fatal(err)
    }
    
    if retrieved != nil {
        fmt.Printf("   Retrieved: [%s] %s\n",
            retrieved.Get("sensorId"),
            retrieved.Get("alertType"))
    } else {
        fmt.Println("   No alerts to process")
    }
    
    fmt.Println("\n✅ Done!")
}
```

**Fichier :** `examples/iot-monitoring/rules.tsd`

```tsd
// IoT Monitoring System

xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

xuple-space analytics {
    selection: lifo,
    consumption: per-agent,
    max-size: 10000
}

type SensorReading {
    sensorId: string,
    temperature: float,
    humidity: float,
    timestamp: int
}

type Alert {
    sensorId: string,
    alertType: string,
    severity: string,
    value: float,
    timestamp: int
}

type AnalyticsEvent {
    sensorId: string,
    eventType: string,
    metadata: string
}

rule HighTemperature {
    when {
        r: SensorReading(temperature > 30.0, temperature <= 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: r.sensorId,
            alertType: "temperature",
            severity: "warning",
            value: r.temperature,
            timestamp: r.timestamp
        )),
        Xuple("analytics", AnalyticsEvent(
            sensorId: r.sensorId,
            eventType: "high_temp",
            metadata: "Temperature above 30°C"
        ))
    }
}

rule CriticalTemperature {
    when {
        r: SensorReading(temperature > 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: r.sensorId,
            alertType: "temperature",
            severity: "critical",
            value: r.temperature,
            timestamp: r.timestamp
        ))
    }
}

rule HighHumidity {
    when {
        r: SensorReading(humidity > 80.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: r.sensorId,
            alertType: "humidity",
            severity: "warning",
            value: r.humidity,
            timestamp: r.timestamp
        ))
    }
}
```

---

### Tâche 5: CI/CD et Automatisation

**Fichier :** `.github/workflows/test.yml`

```yaml
name: Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Cache dependencies
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: unittests
    
    - name: Check coverage
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Total coverage: $COVERAGE%"
        if (( $(echo "$COVERAGE < 85.0" | bc -l) )); then
          echo "❌ Coverage below 85%"
          exit 1
        fi
        echo "✅ Coverage OK"
  
  lint:
    name: Lint
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22
    
    - name: golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
  
  benchmark:
    name: Benchmark
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.22
    
    - name: Run benchmarks
      run: go test -bench=. -benchmem ./test/benchmark/
```

---

## ✅ Checklist de Validation

### Tests

- [ ] Couverture totale > 85%
- [ ] Tous les packages ont > 80% de couverture
- [ ] Tests unitaires pour toutes les fonctions publiques
- [ ] Tests d'intégration pour les flux E2E
- [ ] Tests de performance (benchmarks)
- [ ] Tests de régression
- [ ] Tests de cas limites (edge cases)

### Documentation

- [ ] Guide utilisateur complet
- [ ] API Reference complète
- [ ] FAQ exhaustive
- [ ] Guide de contribution
- [ ] Guide de migration (v1 → v2)
- [ ] CHANGELOG détaillé
- [ ] README mis à jour

### Exemples

- [ ] Au moins 3 exemples complets
- [ ] Exemples documentés
- [ ] Exemples fonctionnels (testés)
- [ ] README dans chaque exemple

### CI/CD

- [ ] GitHub Actions configuré
- [ ] Tests automatisés sur PR
- [ ] Vérification de couverture
- [ ] Linting automatique
- [ ] Benchmarks automatiques

### Release

- [ ] Version taggée (semver)
- [ ] Release notes rédigées
- [ ] Binaries compilés (si applicable)
- [ ] Documentation publiée
- [ ] Annonce de release

---

## 📝 Documentation Finale

### Release Notes

**Fichier :** `RELEASE_NOTES.md`

```markdown
# Release Notes - TSD v2.0.0

## 🎉 Automatisation Complète des Xuples

La version 2.0 marque une évolution majeure avec l'automatisation totale
de la gestion des xuples. Fini la configuration manuelle !

### ✨ Nouveautés Majeures

#### 1. Xuple-Spaces Automatiques

Les xuple-spaces sont maintenant créés automatiquement lors du parsing :

\`\`\`tsd
xuple-space alerts {
    selection: fifo,
    consumption: once
}
\`\`\`

Aucune configuration Go nécessaire !

#### 2. Actions Xuple Automatiques

L'action `Xuple()` fonctionne immédiatement dans les règles :

\`\`\`tsd
rule CreateAlert {
    when {
        t: Temperature(value > 30)
    }
    then {
        Xuple("alerts", Alert(sensorId: t.sensorId, temp: t.value))
    }
}
\`\`\`

#### 3. Package API Unifié

Un seul point d'entrée pour tout :

\`\`\`go
import "github.com/resinsec/tsd/api"

pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")
xuples := result.GetXuples("alerts")
\`\`\`

#### 4. Parser Étendu

- Faits inline dans les actions
- Références aux champs (`t.sensorId`)
- Actions multiples par règle
- Support complet de la syntaxe TSD

### 🔧 Changements Incompatibles

⚠️ **BREAKING