# 🔧 Prompt 05 - Migration des Tests E2E

---

## 🎯 Objectif

**Migrer tous les tests E2E existants pour utiliser exclusivement le package `api` avec `Pipeline.IngestFile()`**, éliminant toute configuration manuelle, tous les workarounds temporaires, et toute création manuelle de xuples ou de xuple-spaces.

### Contexte

À ce stade du projet :
- ✅ Le parser supporte les faits inline et les références aux champs (Prompt 01)
- ✅ Le package `api` fournit un point d'entrée unifié (Prompt 02)
- ✅ Les xuple-spaces sont créés automatiquement (Prompt 03)
- ✅ Les actions `Xuple()` fonctionnent automatiquement (Prompt 04)

**MAIS** : Les tests E2E existants utilisent encore :
- L'ancien code de configuration manuelle (factories, handlers)
- Des workarounds pour créer des xuples manuellement
- Des étapes intermédiaires qui ne devraient plus être nécessaires

L'objectif est de **simplifier radicalement les tests** pour qu'ils reflètent l'expérience utilisateur finale : un seul appel à `IngestFile()` et tout fonctionne.

### Prérequis

- ✅ Prompts 01-04 complétés
- ✅ Package `api` fonctionnel
- ✅ Actions Xuple automatiques
- ✅ Xuple-spaces créés automatiquement

### Résultat Attendu Final

**Avant (complexe, avec workarounds) :**

```go
// Test E2E - avant migration
func TestXuplesE2E_OldWay(t *testing.T) {
    // Configuration manuelle du network
    network := rete.NewNetwork()
    storage := rete.NewMemoryStorage()
    pipeline := constraint.NewConstraintPipeline(network, storage)
    
    // Configuration manuelle du XupleManager
    xupleManager := xuples.NewXupleManager()
    
    // Configuration manuelle de la factory
    pipeline.SetXupleSpaceFactory(func(name string, props map[string]interface{}) error {
        // ... code complexe ...
        return nil
    })
    
    // Ingestion
    pipeline.IngestFile("rules.tsd")
    
    // Création manuelle du xuple-space (workaround)
    xupleManager.CreateXupleSpace("alerts", xuples.SelectionFIFO, ...)
    
    // Création manuelle des xuples (workaround)
    xuple := xuples.NewXuple("Alert", map[string]interface{}{...})
    space.Deposit(xuple)
    
    // Assertions...
}
```

**Après (simple, automatique) :**

```go
// Test E2E - après migration
func TestXuplesE2E_NewWay(t *testing.T) {
    // Un seul appel !
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestFile("rules.tsd")
    
    // Soumettre un fait
    fact := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.5,
    })
    result.Network().Assert(fact)
    
    // Vérifier les xuples créés automatiquement
    xuples := result.GetXuples("alerts")
    assert.Len(t, xuples, 1)
    assert.Equal(t, "Alert", xuples[0].Type())
}
```

---

## 📋 Analyse Préliminaire

### 1. Identifier les Tests E2E Existants

**Fichiers à examiner :**

```
tsd/test/
├── e2e/
│   ├── xuples_e2e_test.go           # Tests E2E principaux
│   ├── xuples_realworld_test.go     # Tests scénarios réels
│   └── integration_test.go          # Tests d'intégration
├── internal/rete/
│   ├── xuple_action_test.go         # Tests unitaires actions
│   └── constraint_pipeline_test.go  # Tests pipeline
└── api/
    └── pipeline_test.go              # Tests API (nouveaux)
```

**Questions à résoudre :**

1. **Quels tests existent actuellement ?**
   - Lister tous les tests E2E liés aux xuples
   - Identifier les tests qui utilisent l'ancien pattern

2. **Quels workarounds sont présents ?**
   - Création manuelle de xuples
   - Configuration manuelle de factories
   - Étapes intermédiaires non nécessaires

3. **Quels tests peuvent être supprimés ?**
   - Tests de l'ancien mécanisme de factory (obsolète)
   - Tests de configuration manuelle
   - Tests redondants avec les nouveaux tests API

### 2. Comprendre la Structure Actuelle

**Pattern typique d'un ancien test :**

```go
func TestOldPattern(t *testing.T) {
    // 1. Setup manuel du network
    network := rete.NewNetwork()
    storage := rete.NewMemoryStorage()
    
    // 2. Setup manuel du pipeline
    pipeline := constraint.NewConstraintPipeline(network, storage)
    
    // 3. Configuration de la factory
    xupleManager := xuples.NewXupleManager()
    pipeline.SetXupleSpaceFactory(createFactory(xupleManager))
    
    // 4. Ingestion
    err := pipeline.IngestFile("testdata/rules.tsd")
    require.NoError(t, err)
    
    // 5. Création manuelle des xuple-spaces (workaround)
    xupleManager.CreateXupleSpace("alerts", ...)
    
    // 6. Tests...
}
```

**Pattern cible (nouveau) :**

```go
func TestNewPattern(t *testing.T) {
    // 1. Créer le pipeline API
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    // 2. Ingérer le fichier (tout est automatique)
    result, err := pipeline.IngestFile("testdata/rules.tsd")
    require.NoError(t, err)
    
    // 3. Tests directement
    xuples := result.GetXuples("alerts")
    // ...
}
```

### 3. Planifier la Migration

**Stratégie :**

1. **Phase 1** : Créer de nouveaux tests E2E avec le pattern API
2. **Phase 2** : Vérifier que les nouveaux tests couvrent tous les scénarios
3. **Phase 3** : Migrer les tests existants un par un
4. **Phase 4** : Supprimer les tests obsolètes
5. **Phase 5** : Nettoyer le code de test (helpers, fixtures, etc.)

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Créer un Répertoire de Tests E2E Migré

**Fichier :** `tsd/test/e2e/xuples_e2e_migrated_test.go`

**Objectif :** Créer une nouvelle suite de tests E2E utilisant le pattern API.

#### 1.1 Setup de Base

```go
package e2e_test

import (
    "os"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/resinsec/tsd/api"
    "github.com/resinsec/tsd/xuples"
)

// createTempTSDFile crée un fichier TSD temporaire pour les tests.
func createTempTSDFile(t *testing.T, content string) string {
    tmpfile, err := os.CreateTemp("", "test_*.tsd")
    require.NoError(t, err)
    
    t.Cleanup(func() {
        os.Remove(tmpfile.Name())
    })
    
    _, err = tmpfile.WriteString(content)
    require.NoError(t, err)
    
    err = tmpfile.Close()
    require.NoError(t, err)
    
    return tmpfile.Name()
}

// assertXuple vérifie qu'un xuple a les propriétés attendues.
func assertXuple(t *testing.T, xuple *xuples.Xuple, typeName string, fields map[string]interface{}) {
    assert.Equal(t, typeName, xuple.Type())
    for field, expectedValue := range fields {
        actualValue := xuple.Get(field)
        assert.Equal(t, expectedValue, actualValue, "field %s mismatch", field)
    }
}
```

#### 1.2 Test E2E Basique

```go
// TestE2E_Basic_XupleCreation teste la création automatique de xuples.
func TestE2E_Basic_XupleCreation(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    message: string,
    temp: float
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            message: "High temperature detected",
            temp: t.value
        ))
    }
}
`
    
    filepath := createTempTSDFile(t, tsdContent)
    
    // Créer le pipeline et ingérer le fichier
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(filepath)
    require.NoError(t, err)
    
    // Vérifier que le xuple-space a été créé
    spaces := result.XupleSpaceNames()
    require.Contains(t, spaces, "alerts")
    
    // Soumettre un fait Temperature
    tempFact := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.5,
    })
    err = result.Network().Assert(tempFact)
    require.NoError(t, err)
    
    // Vérifier qu'un xuple a été créé
    xuples := result.GetXuples("alerts")
    require.Len(t, xuples, 1)
    
    assertXuple(t, xuples[0], "Alert", map[string]interface{}{
        "sensorId": "sensor-01",
        "message":  "High temperature detected",
        "temp":     35.5,
    })
}
```

#### 1.3 Test E2E Multiples Règles

```go
// TestE2E_MultipleRules teste plusieurs règles créant des xuples.
func TestE2E_MultipleRules(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    level: string,
    temp: float
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0, value <= 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            level: "warning",
            temp: t.value
        ))
    }
}

rule CriticalTemperature {
    when {
        t: Temperature(value > 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            level: "critical",
            temp: t.value
        ))
    }
}
`
    
    filepath := createTempTSDFile(t, tsdContent)
    
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(filepath)
    require.NoError(t, err)
    
    // Soumettre un fait déclenchant la règle "warning"
    temp1 := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.0,
    })
    result.Network().Assert(temp1)
    
    xuples := result.GetXuples("alerts")
    require.Len(t, xuples, 1)
    assert.Equal(t, "warning", xuples[0].Get("level"))
    
    // Soumettre un fait déclenchant la règle "critical"
    temp2 := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-02",
        "value":    45.0,
    })
    result.Network().Assert(temp2)
    
    xuples = result.GetXuples("alerts")
    require.Len(t, xuples, 3) // warning + critical + critical (les 2 règles se déclenchent pour > 40)
    
    // Vérifier qu'on a 1 warning et 2 critical
    var warningCount, criticalCount int
    for _, x := range xuples {
        level := x.Get("level").(string)
        if level == "warning" {
            warningCount++
        } else if level == "critical" {
            criticalCount++
        }
    }
    assert.Equal(t, 1, warningCount)
    assert.Equal(t, 2, criticalCount)
}
```

#### 1.4 Test E2E Multiples Xuple-Spaces

```go
// TestE2E_MultipleSpaces teste plusieurs xuple-spaces avec différentes politiques.
func TestE2E_MultipleSpaces(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}

xuple-space logs {
    selection: lifo,
    consumption: per-agent,
    max-size: 100
}

type Event {
    id: string,
    severity: string,
    message: string
}

type Alert {
    eventId: string,
    level: string
}

type LogEntry {
    eventId: string,
    message: string
}

rule CriticalEvent {
    when {
        e: Event(severity == "critical")
    }
    then {
        Xuple("alerts", Alert(
            eventId: e.id,
            level: "critical"
        )),
        Xuple("logs", LogEntry(
            eventId: e.id,
            message: e.message
        ))
    }
}
`
    
    filepath := createTempTSDFile(t, tsdContent)
    
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(filepath)
    require.NoError(t, err)
    
    // Vérifier que les deux xuple-spaces existent
    spaces := result.XupleSpaceNames()
    require.Len(t, spaces, 2)
    require.Contains(t, spaces, "alerts")
    require.Contains(t, spaces, "logs")
    
    // Vérifier les politiques
    alertSpace := result.XupleManager().GetSpace("alerts")
    require.NotNil(t, alertSpace)
    assert.Equal(t, xuples.SelectionFIFO, alertSpace.GetSelectionPolicy())
    assert.Equal(t, xuples.ConsumptionOnce, alertSpace.GetConsumptionPolicy())
    
    logSpace := result.XupleManager().GetSpace("logs")
    require.NotNil(t, logSpace)
    assert.Equal(t, xuples.SelectionLIFO, logSpace.GetSelectionPolicy())
    assert.Equal(t, xuples.ConsumptionPerAgent, logSpace.GetConsumptionPolicy())
    assert.Equal(t, 100, logSpace.GetMaxSize())
    
    // Soumettre un événement critique
    event := result.Network().CreateFact("Event", map[string]interface{}{
        "id":       "evt-001",
        "severity": "critical",
        "message":  "System failure",
    })
    result.Network().Assert(event)
    
    // Vérifier qu'un xuple a été créé dans chaque space
    alerts := result.GetXuples("alerts")
    require.Len(t, alerts, 1)
    assertXuple(t, alerts[0], "Alert", map[string]interface{}{
        "eventId": "evt-001",
        "level":   "critical",
    })
    
    logs := result.GetXuples("logs")
    require.Len(t, logs, 1)
    assertXuple(t, logs[0], "LogEntry", map[string]interface{}{
        "eventId": "evt-001",
        "message": "System failure",
    })
}
```

#### 1.5 Test E2E Retrieve (Consommation)

```go
// TestE2E_Retrieve teste la récupération et consommation de xuples.
func TestE2E_Retrieve(t *testing.T) {
    tsdContent := `
xuple-space notifications {
    selection: fifo,
    consumption: once
}

type Notification {
    id: string,
    message: string
}

type Event {
    id: string
}

rule CreateNotification {
    when {
        e: Event()
    }
    then {
        Xuple("notifications", Notification(
            id: e.id,
            message: "Event processed"
        ))
    }
}
`
    
    filepath := createTempTSDFile(t, tsdContent)
    
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(filepath)
    require.NoError(t, err)
    
    // Créer plusieurs notifications
    for i := 1; i <= 3; i++ {
        event := result.Network().CreateFact("Event", map[string]interface{}{
            "id": fmt.Sprintf("evt-%03d", i),
        })
        result.Network().Assert(event)
    }
    
    // Vérifier qu'il y a 3 notifications
    xuples := result.GetXuples("notifications")
    require.Len(t, xuples, 3)
    
    // Récupérer (consommer) une notification
    retrieved, err := result.Retrieve("notifications", nil, "agent-01")
    require.NoError(t, err)
    require.NotNil(t, retrieved)
    assert.Equal(t, "Notification", retrieved.Type())
    
    // Il devrait maintenant rester 2 notifications (consumption: once)
    xuples = result.GetXuples("notifications")
    require.Len(t, xuples, 2)
    
    // Récupérer une deuxième fois avec le même agent (ne devrait rien retourner)
    retrieved2, err := result.Retrieve("notifications", nil, "agent-01")
    require.NoError(t, err)
    assert.Nil(t, retrieved2) // Déjà consommé
    
    // Avec un autre agent
    retrieved3, err := result.Retrieve("notifications", nil, "agent-02")
    require.NoError(t, err)
    require.NotNil(t, retrieved3)
    
    // Il reste 1 notification
    xuples = result.GetXuples("notifications")
    require.Len(t, xuples, 1)
}
```

---

### Tâche 2: Migrer les Tests Existants

**Fichier :** `tsd/test/e2e/xuples_e2e_test.go`

**Objectif :** Migrer les tests existants vers le nouveau pattern.

#### 2.1 Inventaire des Tests à Migrer

**Créer un fichier de tracking :**

```markdown
# Migration des Tests E2E - Checklist

## Tests à Migrer

### xuples_e2e_test.go
- [ ] TestXuplesE2E_BasicFlow
- [ ] TestXuplesE2E_RealWorld
- [ ] TestXuplesE2E_MultipleSpaces
- [ ] TestXuplesE2E_Policies

### xuples_realworld_test.go
- [ ] TestRealWorld_IoTSensors
- [ ] TestRealWorld_OrderProcessing
- [ ] TestRealWorld_EventNotifications

### integration_test.go
- [ ] TestIntegration_XupleSpaceCreation
- [ ] TestIntegration_XupleActionExecution

## Tests à Supprimer (obsolètes)
- [ ] TestXupleSpaceFactory_Manual
- [ ] TestXupleSpaceFactory_Configuration
- [ ] TestXupleAction_ManualSetup

## Nouveaux Tests à Ajouter
- [ ] TestE2E_ErrorHandling_InvalidSpace
- [ ] TestE2E_ErrorHandling_MissingType
- [ ] TestE2E_Performance_LargeVolume
```

#### 2.2 Template de Migration

**Pour chaque test à migrer, suivre ce template :**

```go
// AVANT (ancien pattern)
func TestOldPattern(t *testing.T) {
    // Setup manuel complexe
    network := rete.NewNetwork()
    storage := rete.NewMemoryStorage()
    pipeline := constraint.NewConstraintPipeline(network, storage)
    xupleManager := xuples.NewXupleManager()
    
    // Factory manuelle
    pipeline.SetXupleSpaceFactory(func(name string, props map[string]interface{}) error {
        // ... code ...
        return nil
    })
    
    // Ingestion
    err := pipeline.IngestFile("testdata/scenario1.tsd")
    require.NoError(t, err)
    
    // Workaround: création manuelle du space
    space, _ := xupleManager.CreateXupleSpace("alerts", ...)
    
    // Workaround: création manuelle du xuple
    xuple := xuples.NewXuple("Alert", map[string]interface{}{
        "message": "test",
    })
    space.Deposit(xuple)
    
    // Assertions
    xuples := space.GetAll()
    assert.Len(t, xuples, 1)
}

// APRÈS (nouveau pattern)
func TestNewPattern(t *testing.T) {
    // Lecture du fichier TSD
    tsdContent := `
xuple-space alerts {
    selection: fifo
}

type Alert {
    message: string
}

rule CreateAlert {
    when {
        trigger: Trigger()
    }
    then {
        Xuple("alerts", Alert(message: "test"))
    }
}
`
    
    filepath := createTempTSDFile(t, tsdContent)
    
    // Pipeline API (tout est automatique)
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(filepath)
    require.NoError(t, err)
    
    // Déclencher la règle
    trigger := result.Network().CreateFact("Trigger", map[string]interface{}{})
    result.Network().Assert(trigger)
    
    // Assertions (xuples créés automatiquement)
    xuples := result.GetXuples("alerts")
    assert.Len(t, xuples, 1)
    assert.Equal(t, "test", xuples[0].Get("message"))
}
```

#### 2.3 Migration du Test RealWorld

**Exemple de migration d'un test réel complexe :**

```go
// TestRealWorld_IoTSensors_Migrated teste un scénario IoT complet.
func TestRealWorld_IoTSensors_Migrated(t *testing.T) {
    tsdContent := `
// Définition des xuple-spaces
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

// Types de données
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

// Règles de détection
rule HighTemperature {
    when {
        reading: SensorReading(temperature > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: reading.sensorId,
            alertType: "temperature",
            severity: "warning",
            value: reading.temperature,
            timestamp: reading.timestamp
        )),
        Xuple("analytics", AnalyticsEvent(
            sensorId: reading.sensorId,
            eventType: "high_temp",
            metadata: "Temperature exceeded threshold"
        ))
    }
}

rule CriticalTemperature {
    when {
        reading: SensorReading(temperature > 40.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: reading.sensorId,
            alertType: "temperature",
            severity: "critical",
            value: reading.temperature,
            timestamp: reading.timestamp
        ))
    }
}

rule HighHumidity {
    when {
        reading: SensorReading(humidity > 80.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: reading.sensorId,
            alertType: "humidity",
            severity: "warning",
            value: reading.humidity,
            timestamp: reading.timestamp
        ))
    }
}
`
    
    filepath := createTempTSDFile(t, tsdContent)
    
    // Créer le pipeline
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(filepath)
    require.NoError(t, err)
    
    // Vérifier que les xuple-spaces ont été créés avec les bonnes politiques
    spaces := result.XupleSpaceNames()
    require.Len(t, spaces, 2)
    require.Contains(t, spaces, "alerts")
    require.Contains(t, spaces, "analytics")
    
    // Scénario 1: Température normale (pas d'alerte)
    reading1 := result.Network().CreateFact("SensorReading", map[string]interface{}{
        "sensorId":    "sensor-001",
        "temperature": 25.0,
        "humidity":    60.0,
        "timestamp":   1234567890,
    })
    result.Network().Assert(reading1)
    
    alerts := result.GetXuples("alerts")
    assert.Len(t, alerts, 0, "No alerts for normal readings")
    
    // Scénario 2: Température élevée (warning)
    reading2 := result.Network().CreateFact("SensorReading", map[string]interface{}{
        "sensorId":    "sensor-002",
        "temperature": 35.0,
        "humidity":    65.0,
        "timestamp":   1234567900,
    })
    result.Network().Assert(reading2)
    
    alerts = result.GetXuples("alerts")
    require.Len(t, alerts, 1)
    assertXuple(t, alerts[0], "Alert", map[string]interface{}{
        "sensorId":  "sensor-002",
        "alertType": "temperature",
        "severity":  "warning",
        "value":     35.0,
    })
    
    analytics := result.GetXuples("analytics")
    require.Len(t, analytics, 1)
    assert.Equal(t, "high_temp", analytics[0].Get("eventType"))
    
    // Scénario 3: Température critique (2 alertes: warning + critical)
    reading3 := result.Network().CreateFact("SensorReading", map[string]interface{}{
        "sensorId":    "sensor-003",
        "temperature": 45.0,
        "humidity":    70.0,
        "timestamp":   1234567910,
    })
    result.Network().Assert(reading3)
    
    alerts = result.GetXuples("alerts")
    require.Len(t, alerts, 3) // 1 from scenario 2 + 2 from scenario 3
    
    // Vérifier qu'on a bien 1 critical et 1 warning pour sensor-003
    sensor003Alerts := filterXuples(alerts, func(x *xuples.Xuple) bool {
        return x.Get("sensorId") == "sensor-003"
    })
    require.Len(t, sensor003Alerts, 2)
    
    severities := []string{
        sensor003Alerts[0].Get("severity").(string),
        sensor003Alerts[1].Get("severity").(string),
    }
    assert.Contains(t, severities, "warning")
    assert.Contains(t, severities, "critical")
    
    // Scénario 4: Humidité élevée
    reading4 := result.Network().CreateFact("SensorReading", map[string]interface{}{
        "sensorId":    "sensor-004",
        "temperature": 28.0,
        "humidity":    85.0,
        "timestamp":   1234567920,
    })
    result.Network().Assert(reading4)
    
    alerts = result.GetXuples("alerts")
    require.Len(t, alerts, 4)
    
    // Vérifier la dernière alerte (humidité)
    humidityAlerts := filterXuples(alerts, func(x *xuples.Xuple) bool {
        return x.Get("alertType") == "humidity"
    })
    require.Len(t, humidityAlerts, 1)
    assert.Equal(t, 85.0, humidityAlerts[0].Get("value"))
    
    // Test de récupération (simulation d'un agent de monitoring)
    retrieved, err := result.Retrieve("alerts", nil, "monitoring-agent-01")
    require.NoError(t, err)
    require.NotNil(t, retrieved)
    
    // L'agent devrait récupérer le premier (FIFO)
    assert.Equal(t, "sensor-002", retrieved.Get("sensorId"))
    
    // Il reste maintenant 3 alertes
    alerts = result.GetXuples("alerts")
    assert.Len(t, alerts, 3)
}

// filterXuples est un helper pour filtrer les xuples.
func filterXuples(xuples []*xuples.Xuple, predicate func(*xuples.Xuple) bool) []*xuples.Xuple {
    var result []*xuples.Xuple
    for _, x := range xuples {
        if predicate(x) {
            result = append(result, x)
        }
    }
    return result
}
```

---

### Tâche 3: Supprimer les Tests Obsolètes

**Fichier :** `tsd/test/internal/rete/xuple_factory_test.go` (à supprimer)

**Objectif :** Identifier et supprimer les tests qui testent du code obsolète.

#### 3.1 Liste des Tests à Supprimer

```go
// Ces tests doivent être SUPPRIMÉS car ils testent l'ancien pattern

// ❌ OBSOLETE - Supprimer
func TestXupleSpaceFactory_ManualConfiguration(t *testing.T) {
    // Test de la configuration manuelle de la factory
    // Obsolète car la factory est maintenant configurée automatiquement
}

// ❌ OBSOLETE - Supprimer
func TestXupleSpaceFactory_CreateFromDefinition(t *testing.T) {
    // Test de la création manuelle via factory
    // Obsolète car les xuple-spaces sont créés automatiquement
}

// ❌ OBSOLETE - Supprimer
func TestXupleAction_ManualSetup(t *testing.T) {
    // Test du setup manuel de l'action Xuple
    // Obsolète car l'action est enregistrée automatiquement
}

// ❌ OBSOLETE - Supprimer
func TestConstraintPipeline_SetXupleSpaceFactory(t *testing.T) {
    // Test de la méthode SetXupleSpaceFactory
    // Obsolète car cette méthode sera supprimée (Prompt 06)
}
```

#### 3.2 Commandes de Suppression

```bash
# Supprimer les fichiers de tests obsolètes
rm tsd/test/internal/rete/xuple_factory_test.go
rm tsd/test/internal/rete/manual_setup_test.go

# Commenter temporairement (pour vérification) puis supprimer
# Les tests suivants dans xuples_e2e_test.go:
# - TestXuplesE2E_OldPattern (ligne X-Y)
# - TestXuplesE2E_ManualFactory (ligne Z-W)
```

---

### Tâche 4: Mettre à Jour les Helpers de Test

**Fichier :** `tsd/test/testutil/helpers.go`

**Objectif :** Créer des helpers pour simplifier les tests avec le nouveau pattern.

```go
package testutil

import (
    "os"
    "testing"
    
    "github.com/stretchr/testify/require"
    
    "github.com/resinsec/tsd/api"
    "github.com/resinsec/tsd/xuples"
)

// CreatePipelineFromTSD crée un pipeline à partir d'un contenu TSD.
// Retourne le pipeline et le résultat de l'ingestion.
func CreatePipelineFromTSD(t *testing.T, tsdContent string) (*api.Pipeline, *api.Result) {
    tmpfile, err := os.CreateTemp("", "test_*.tsd")
    require.NoError(t, err)
    
    t.Cleanup(func() {
        os.Remove(tmpfile.Name())
    })
    
    _, err = tmpfile.WriteString(tsdContent)
    require.NoError(t, err)
    tmpfile.Close()
    
    pipeline, err := api.NewPipeline()
    require.NoError(t, err)
    
    result, err := pipeline.IngestFile(tmpfile.Name())
    require.NoError(t, err)
    
    return pipeline, result
}

// AssertXupleFields vérifie qu'un xuple a les champs attendus.
func AssertXupleFields(t *testing.T, xuple *xuples.Xuple, expectedType string, expectedFields map[string]interface{}) {
    require.NotNil(t, xuple, "xuple should not be nil")
    require.Equal(t, expectedType, xuple.Type(), "xuple type mismatch")
    
    for field, expectedValue := range expectedFields {
        actualValue := xuple.Get(field)
        require.Equal(t, expectedValue, actualValue, "field '%s' mismatch", field)
    }
}

// AssertXupleSpaceExists vérifie qu'un xuple-space existe avec les bonnes politiques.
func AssertXupleSpaceExists(t *testing.T, result *api.Result, spaceName string, selection xuples.SelectionPolicy, consumption xuples.ConsumptionPolicy) {
    spaces := result.XupleSpaceNames()
    require.Contains(t, spaces, spaceName, "xuple-space '%s' should exist", spaceName)
    
    space := result.XupleManager().GetSpace(spaceName)
    require.NotNil(t, space, "xuple-space '%s' should not be nil", spaceName)
    require.Equal(t, selection, space.GetSelectionPolicy(), "selection policy mismatch for space '%s'", spaceName)
    require.Equal(t, consumption, space.GetConsumptionPolicy(), "consumption policy mismatch for space '%s'", spaceName)
}

// SubmitFact crée et soumet un fait au réseau.
func SubmitFact(t *testing.T, result *api.Result, typeName string, data map[string]interface{}) {
    fact := result.Network().CreateFact(typeName, data)
    err := result.Network().Assert(fact)
    require.NoError(t, err, "failed to assert fact of type '%s'", typeName)
}

// GetXuplesCount retourne le nombre de xuples dans un xuple-space.
func GetXuplesCount(t *testing.T, result *api.Result, spaceName string) int {
    return result.XupleCount(spaceName)
}

// RetrieveAndAssert récupère un xuple et vérifie ses propriétés.
func RetrieveAndAssert(t *testing.T, result *api.Result, spaceName string, agentID string, expectedType string, expectedFields map[string]interface{}) {
    xuple, err := result.Retrieve(spaceName, nil, agentID)
    require.NoError(t, err, "failed to retrieve from space '%s'", spaceName)
    require.NotNil(t, xuple, "retrieved xuple should not be nil")
    
    AssertXupleFields(t, xuple, expectedType, expectedFields)
}
```

#### 4.2 Utilisation des Helpers

```go
// Exemple d'utilisation des helpers dans un test
func TestE2E_WithHelpers(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}

type Event {
    id: string,
    severity: string
}

type Alert {
    eventId: string,
    level: string
}

rule CriticalEvent {
    when {
        e: Event(severity == "critical")
    }
    then {
        Xuple("alerts", Alert(
            eventId: e.id,
            level: "critical"
        ))
    }
}
`
    
    // Créer le pipeline (1 ligne !)
    _, result := testutil.CreatePipelineFromTSD(t, tsdContent)
    
    // Vérifier le xuple-space
    testutil.AssertXupleSpaceExists(t, result, "alerts", xuples.SelectionFIFO, xuples.ConsumptionOnce)
    
    // Soumettre un événement
    testutil.SubmitFact(t, result, "Event", map[string]interface{}{
        "id":       "evt-001",
        "severity": "critical",
    })
    
    // Vérifier le nombre de xuples
    count := testutil.GetXuplesCount(t, result, "alerts")
    assert.Equal(t, 1, count)
    
    // Récupérer et vérifier
    testutil.RetrieveAndAssert(t, result, "alerts", "agent-01", "Alert", map[string]interface{}{
        "eventId": "evt-001",
        "level":   "critical",
    })
    
    // Vérifier qu'il ne reste plus de xuples
    count = testutil.GetXuplesCount(t, result, "alerts")
    assert.Equal(t, 0, count)
}
```

---

### Tâche 5: Mettre à Jour les Fixtures de Test

**Répertoire :** `tsd/test/testdata/`

**Objectif :** Mettre à jour les fichiers TSD de test pour utiliser la nouvelle syntaxe complète.

#### 5.1 Fichier de Test Basique

**Fichier :** `tsd/test/testdata/basic_xuple.tsd`

```tsd
// Test basique de création de xuples

xuple-space alerts {
    selection: fifo,
    consumption: once
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    message: string,
    temp: float
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            message: "High temperature detected",
            temp: t.value
        ))
    }
}
```

#### 5.2 Fichier de Test Complexe

**Fichier :** `tsd/test/testdata/complex_xuples.tsd`

```tsd
// Test complexe avec multiples spaces et règles

xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

xuple-space logs {
    selection: lifo,
    consumption: per-agent,
    max-size: 1000
}

xuple-space analytics {
    selection: random,
    consumption: per-agent,
    max-size: 5000
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

type LogEntry {
    sensorId: string,
    eventType: string,
    message: string,
    timestamp: int
}

type AnalyticsEvent {
    sensorId: string,
    metric: string,
    value: float,
    timestamp: int
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
        Xuple("logs", LogEntry(
            sensorId: r.sensorId,
            eventType: "high_temperature",
            message: "Temperature above normal threshold",
            timestamp: r.timestamp
        )),
        Xuple("analytics", AnalyticsEvent(
            sensorId: r.sensorId,
            metric: "temperature",
            value: r.temperature,
            timestamp: r.timestamp
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
        )),
        Xuple("logs", LogEntry(
            sensorId: r.sensorId,
            eventType: "critical_temperature",
            message: "CRITICAL: Temperature exceeds safety limit",
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

## 🧪 Tests de Validation

### Test 1: Vérification de la Migration

**Fichier :** `tsd/test/e2e/migration_validation_test.go`

```go
// TestMigrationValidation vérifie que tous les scénarios fonctionnent avec le nouveau pattern.
func TestMigrationValidation_AllScenariosWork(t *testing.T) {
    testCases := []struct {
        name        string
        tsdFile     string
        setupFacts  func(*api.Result)
        assertions  func(*testing.T, *api.Result)
    }{
        {
            name:    "Basic xuple creation",
            tsdFile: "testdata/basic_xuple.tsd",
            setupFacts: func(r *api.Result) {
                testutil.SubmitFact(t, r, "Temperature", map[string]interface{}{
                    "sensorId": "sensor-01",
                    "value":    35.0,
                })
            },
            assertions: func(t *testing.T, r *api.Result) {
                assert.Equal(t, 1, r.XupleCount("alerts"))
            },
        },
        {
            name:    "Complex multi-space scenario",
            tsdFile: "testdata/complex_xuples.tsd",
            setupFacts: func(r *api.Result) {
                testutil.SubmitFact(t, r, "SensorReading", map[string]interface{}{
                    "sensorId":    "sensor-02",
                    "temperature": 45.0,
                    "humidity":    85.0,
                    "timestamp":   1234567890,
                })
            },
            assertions: func(t *testing.T, r *api.Result) {
                // Température critique + humidité élevée
                assert.Equal(t, 3, r.XupleCount("alerts")) // 2 temp + 1 humidity
                assert.Equal(t, 2, r.XupleCount("logs"))   // 2 temp logs
                assert.Equal(t, 1, r.XupleCount("analytics")) // 1 analytics (only warning)
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            pipeline, err := api.NewPipeline()
            require.NoError(t, err)
            
            result, err := pipeline.IngestFile(tc.tsdFile)
            require.NoError(t, err)
            
            tc.setupFacts(result)
            tc.assertions(t, result)
        })
    }
}
```

### Test 2: Régression - Compatibilité Ascendante

```go
// TestNoRegression vérifie qu'il n'y a pas de régression par rapport aux tests existants.
func TestNoRegression_ExistingScenarios(t *testing.T) {
    // Ce test exécute les mêmes scénarios que les anciens tests
    // mais avec le nouveau pattern, pour vérifier qu'on obtient les mêmes résultats
    
    t.Run("Scenario 1: Single rule, single space", func(t *testing.T) {
        // Ancien résultat attendu: 1 xuple créé
        _, result := testutil.CreatePipelineFromTSD(t, `
xuple-space test { selection: fifo }
type T { id: string }
type X { id: string }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`)
        
        testutil.SubmitFact(t, result, "T", map[string]interface{}{"id": "123"})
        assert.Equal(t, 1, result.XupleCount("test"))
    })
    
    // Ajouter d'autres scénarios de régression...
}
```

---

## ✅ Checklist de Validation

### Migration des Tests

- [ ] Tous les tests E2E identifiés et listés
- [ ] Nouveaux tests créés avec le pattern API
- [ ] Tests existants migrés un par un
- [ ] Chaque test migré produit les mêmes résultats
- [ ] Tests obsolètes identifiés et supprimés
- [ ] Pas de tests redondants

### Helpers et Utilities

- [ ] `testutil.CreatePipelineFromTSD()` implémenté
- [ ] `testutil.AssertXupleFields()` implémenté
- [ ] `testutil.AssertXupleSpaceExists()` implémenté
- [ ] `testutil.SubmitFact()` implémenté
- [ ] Tous les helpers documentés avec GoDoc

### Fixtures

- [ ] Fichiers TSD de test mis à jour
- [ ] Syntaxe complète utilisée (xuple-space + règles)
- [ ] Fichiers organisés par complexité
- [ ] Commentaires explicatifs ajoutés

### Validation

- [ ] Tests de migration passent (100%)
- [ ] Tests de régression passent (100%)
- [ ] Couverture de code maintenue ou améliorée (> 80%)
- [ ] Aucun test ne contient de code obsolète

### Standards

- [ ] Code formaté (`gofmt`)
- [ ] Pas de warnings du linter
- [ ] Commentaires GoDoc complets
- [ ] Pas de code mort (dead code)

---

## 📝 Documentation à Mettre à Jour

### 1. Guide de Test (`docs/TESTING.md`)

```markdown
## Tests E2E - Xuples

### Pattern Recommandé

Tous les tests E2E doivent utiliser le package `api` :

\`\`\`go
import (
    "testing"
    "github.com/resinsec/tsd/api"
    "github.com/resinsec/tsd/test/testutil"
)

func TestMyScenario(t *testing.T) {
    tsdContent := `...`
    
    // Créer le pipeline
    _, result := testutil.CreatePipelineFromTSD(t, tsdContent)
    
    // Soumettre des faits
    testutil.SubmitFact(t, result, "MyType", map[string]interface{}{
        "field": "value",
    })
    
    // Assertions
    assert.Equal(t, 1, result.XupleCount("myspace"))
}
\`\`\`

### Helpers Disponibles

- `CreatePipelineFromTSD(t, content)` : Crée un pipeline depuis du TSD
- `AssertXupleFields(t, xuple, type, fields)` : Vérifie les champs d'un xuple
- `SubmitFact(t, result, type, data)` : Soumet un fait au réseau
- (voir `test/testutil/helpers.go` pour la liste complète)

### ❌ Patterns à Éviter

**N'utilisez PLUS ces patterns :**

\`\`\`go
// ❌ OBSOLETE - Ne pas faire
network := rete.NewNetwork()
pipeline := constraint.NewConstraintPipeline(network, ...)
pipeline.SetXupleSpaceFactory(...)
xupleManager.CreateXupleSpace(...)
\`\`\`

**Utilisez le package API à la place :**

\`\`\`go
// ✅ CORRECT
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")
\`\`\`
```

### 2. Guide de Contribution (`CONTRIBUTING.md`)

```markdown
## Écrire des Tests E2E

### Structure Recommandée

1. **Créer le contenu TSD inline** : Facilite la lecture du test
2. **Utiliser les helpers** : `testutil.CreatePipelineFromTSD()`
3. **Assertions claires** : Un assert par comportement vérifié
4. **Nommage descriptif** : `TestE2E_<Scenario>_<Behavior>`

### Exemple Complet

\`\`\`go
func TestE2E_TemperatureMonitoring_CreatesAlertsOnHighTemp(t *testing.T) {
    // 1. Définir le TSD
    tsdContent := `
        xuple-space alerts { selection: fifo }
        type Temperature { value: float }
        type Alert { message: string }
        rule High { 
            when { t: Temperature(value > 30) } 
            then { Xuple("alerts", Alert(message: "High!")) } 
        }
    `
    
    // 2. Créer le pipeline
    _, result := testutil.CreatePipelineFromTSD(t, tsdContent)
    
    // 3. Setup
    testutil.SubmitFact(t, result, "Temperature", map[string]interface{}{
        "value": 35.0,
    })
    
    // 4. Assertions
    xuples := result.GetXuples("alerts")
    require.Len(t, xuples, 1)
    assert.Equal(t, "High!", xuples[0].Get("message"))
}
\`\`\`
```

---

## 🎯 Résultat Attendu

### Métriques de Succès

**Avant la migration :**
- Nombre de tests E2E : ~15
- Tests avec configuration manuelle : ~10 (67%)
- Lignes de code de setup par test : ~30-50
- Tests avec workarounds : ~8 (53%)

**Après la migration :**
- Nombre de tests E2E : ~20 (+ nouveaux scénarios)
- Tests avec configuration manuelle : 0 (0%)
- Lignes de code de setup par test : ~5-10
- Tests avec workarounds : 0 (0%)
- Couverture de code : > 85%

### Bénéfices

1. **Simplicité** : Tests beaucoup plus courts et lisibles
2. **Maintenabilité** : Changements d'API impactent un seul endroit
3. **Robustesse** : Moins de dépendances entre composants
4. **Documentation** : Les tests servent d'exemples d'utilisation
5. **Confiance** : Couverture améliorée, moins de code mort

---

## 🔗 Dépendances

### Entrantes

- ✅ Prompt 01 : Parser complet
- ✅ Prompt 02 : Package API
- ✅ Prompt 03 : Xuple-spaces automatiques
- ✅ Prompt 04 : Actions automatiques

### Sortantes

- ➡️ Prompt 06 : Cleanup (supprimera le code obsolète testé ici)
- ➡️ Prompt 07 : Documentation finale (inclura les exemples des tests)

---

## 🚀 Stratégie d'Implémentation

1. **Phase 1: Nouveaux Tests** (2h)
   - Créer `xuples_e2e_migrated_test.go`
   - Implémenter 5-6 tests de base avec le nouveau pattern
   - Vérifier qu'ils passent tous

2. **Phase 2: Helpers** (1h)
   - Créer `testutil/helpers.go`
   - Implémenter tous les helpers listés
   - Refactorer les nouveaux tests pour utiliser les helpers

3. **Phase 3: Migration** (3-4h)
   - Migrer les tests existants un par un
   - Vérifier que chaque test migré passe
   - Commit après chaque migration réussie

4. **Phase 4: Suppression** (1h)
   - Identifier les tests obsolètes
   - Les supprimer
   - Vérifier que tous les tests passent encore

5. **Phase 5: Fixtures** (1h)
   - Mettre à jour les fichiers TSD de test
   - Ajouter des fichiers pour les nouveaux scénarios
   - Organiser par complexité

6. **Phase 6: Documentation** (1h)
   - Mettre à jour TESTING.md
   - Mettre à jour CONTRIBUTING.md
   - Ajouter des exemples

**Estimation totale : 9-11 heures**

---

## 📊 Critères de Succès

- [ ] Tous les tests E2E migrés vers le pattern API
- [ ] Zéro test utilisant l'ancien pattern (factory manuelle)
- [ ] Helpers de test créés et utilisés
- [ ] Fixtures TSD mises à jour
- [ ] Tests obsolètes supprimés
- [ ] Couverture de code maintenue (> 80%)
- [ ] Tous les tests passent (100%)
- [ ] Documentation à jour
- [ ] Pas de code mort dans les tests
- [ ] Temps d'exécution des tests identique ou amélioré

---

**FIN DU PROMPT 05**