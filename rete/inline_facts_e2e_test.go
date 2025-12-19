// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_InlineFact_SimpleXuple teste l'utilisation d'un fait inline simple dans une action Xuple
func TestE2E_InlineFact_SimpleXuple(t *testing.T) {
	t.Log("🧪 TEST E2E - FAIT INLINE SIMPLE DANS ACTION XUPLE")
	t.Log("==================================================")

	program := `
		type Sensor(#sensorId: string, temp: number)
		type Alert(level: string, sensorId: string)
		
		rule high_temp: {s: Sensor} / s.temp > 40.0 ==>
			Xuple("alerts", Alert(level: "HIGH", sensorId: "S001"))
		
		Sensor(sensorId: "S001", temp: 25.0)
		Sensor(sensorId: "S002", temp: 45.0)
	`

	// Créer un fichier temporaire
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsd")
	err := os.WriteFile(testFile, []byte(program), 0644)
	require.NoError(t, err, "❌ Création du fichier temporaire échouée")

	// Créer le réseau RETE et le storage
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Variable pour capturer les xuples créés
	createdXuples := make([]*Fact, 0)

	// Configurer le handler Xuple AVANT l'ingestion
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		t.Logf("✅ Xuple créé dans '%s': Type=%s, Fields=%+v", xuplespace, fact.Type, fact.Fields)
		createdXuples = append(createdXuples, fact)
		return nil
	})

	// Créer le pipeline et ingérer le programme
	pipeline := NewConstraintPipeline()
	network, _, err = pipeline.IngestFile(testFile, network, storage)
	require.NoError(t, err, "❌ Ingestion du programme échouée")

	// Vérifier qu'un xuple a été créé (seul S002 déclenche la règle car temp > 40)
	require.Len(t, createdXuples, 1, "❌ Devrait avoir créé 1 xuple")

	alert := createdXuples[0]
	assert.Equal(t, "Alert", alert.Type, "❌ Le type devrait être Alert")
	assert.Equal(t, "HIGH", alert.Fields["level"], "❌ Le niveau devrait être HIGH")
	assert.Equal(t, "S001", alert.Fields["sensorId"], "❌ Le sensorId devrait être S001")

	t.Log("✅ Test E2E fait inline simple réussi")
}

// TestE2E_InlineFact_WithFieldReferences teste les références aux champs de variables
func TestE2E_InlineFact_WithFieldReferences(t *testing.T) {
	t.Log("🧪 TEST E2E - FAIT INLINE AVEC RÉFÉRENCES AUX CHAMPS")
	t.Log("====================================================")

	program := `
		type Sensor(#sensorId: string, temperature: number)
		type Alert(level: string, sensorId: string, temperature: number)
		
		rule high_temp: {s: Sensor} / s.temperature > 40.0 ==>
			Xuple("alerts", Alert(
				level: "HIGH",
				sensorId: s.sensorId,
				temperature: s.temperature
			))
		
		Sensor(sensorId: "SENSOR-001", temperature: 25.0)
		Sensor(sensorId: "SENSOR-002", temperature: 45.0)
	`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsd")
	err := os.WriteFile(testFile, []byte(program), 0644)
	require.NoError(t, err)

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	createdXuples := make([]*Fact, 0)
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		t.Logf("✅ Xuple créé: Type=%s, Fields=%+v", fact.Type, fact.Fields)
		createdXuples = append(createdXuples, fact)
		return nil
	})

	pipeline := NewConstraintPipeline()
	network, _, err = pipeline.IngestFile(testFile, network, storage)
	require.NoError(t, err, "❌ Ingestion du programme échouée")

	require.Len(t, createdXuples, 1, "❌ Devrait avoir créé 1 xuple")

	alert := createdXuples[0]
	assert.Equal(t, "Alert", alert.Type)
	assert.Equal(t, "HIGH", alert.Fields["level"])
	assert.Equal(t, "SENSOR-002", alert.Fields["sensorId"], "❌ Devrait copier le sensorId du sensor")
	assert.Equal(t, 45.0, alert.Fields["temperature"], "❌ Devrait copier la température du sensor")

	t.Log("✅ Test E2E avec références aux champs réussi")
}

// TestE2E_InlineFact_MultipleActions teste plusieurs actions avec des faits inline
func TestE2E_InlineFact_MultipleActions(t *testing.T) {
	t.Log("🧪 TEST E2E - ACTIONS MULTIPLES AVEC FAITS INLINE")
	t.Log("==================================================")

	program := `
		type Sensor(#sensorId: string, temperature: number, location: string)
		type Alert(level: string, sensorId: string)
		type Command(action: string, target: string)
		
		rule critical_temp: {s: Sensor} / s.temperature > 40.0 ==>
			Xuple("alerts", Alert(level: "CRITICAL", sensorId: s.sensorId)),
			Xuple("commands", Command(action: "shutdown", target: s.location))
		
		Sensor(sensorId: "S001", temperature: 45.0, location: "ServerRoom")
	`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsd")
	err := os.WriteFile(testFile, []byte(program), 0644)
	require.NoError(t, err)

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	createdXuples := make(map[string][]*Fact)
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		t.Logf("✅ Xuple créé dans '%s': Type=%s", xuplespace, fact.Type)
		createdXuples[xuplespace] = append(createdXuples[xuplespace], fact)
		return nil
	})

	pipeline := NewConstraintPipeline()
	network, _, err = pipeline.IngestFile(testFile, network, storage)
	require.NoError(t, err, "❌ Ingestion du programme échouée")

	// Vérifier qu'on a créé 2 xuples (1 Alert + 1 Command)
	require.Len(t, createdXuples["alerts"], 1, "❌ Devrait avoir 1 alert")
	require.Len(t, createdXuples["commands"], 1, "❌ Devrait avoir 1 command")

	alert := createdXuples["alerts"][0]
	assert.Equal(t, "Alert", alert.Type)
	assert.Equal(t, "CRITICAL", alert.Fields["level"])
	assert.Equal(t, "S001", alert.Fields["sensorId"])

	command := createdXuples["commands"][0]
	assert.Equal(t, "Command", command.Type)
	assert.Equal(t, "shutdown", command.Fields["action"])
	assert.Equal(t, "ServerRoom", command.Fields["target"])

	t.Log("✅ Test E2E actions multiples réussi")
}

// TestE2E_InlineFact_WithExpressions teste les expressions dans les champs de faits inline
func TestE2E_InlineFact_WithExpressions(t *testing.T) {
	t.Log("🧪 TEST E2E - FAIT INLINE AVEC EXPRESSIONS")
	t.Log("===========================================")

	program := `
		type Sensor(#sensorId: string, temperature: number)
		type Alert(level: string, tempCelsius: number, tempFahrenheit: number)
		
		rule high_temp: {s: Sensor} / s.temperature > 40.0 ==>
			Xuple("alerts", Alert(
				level: "HIGH",
				tempCelsius: s.temperature,
				tempFahrenheit: s.temperature * 1.8 + 32.0
			))
		
		Sensor(sensorId: "S001", temperature: 45.0)
	`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsd")
	err := os.WriteFile(testFile, []byte(program), 0644)
	require.NoError(t, err)

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	createdXuples := make([]*Fact, 0)
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		createdXuples = append(createdXuples, fact)
		return nil
	})

	pipeline := NewConstraintPipeline()
	network, _, err = pipeline.IngestFile(testFile, network, storage)
	require.NoError(t, err, "❌ Ingestion du programme échouée")

	require.Len(t, createdXuples, 1, "❌ Devrait avoir créé 1 xuple")

	alert := createdXuples[0]
	assert.Equal(t, "Alert", alert.Type)
	assert.Equal(t, "HIGH", alert.Fields["level"])
	assert.Equal(t, 45.0, alert.Fields["tempCelsius"])

	// Vérifier la conversion Celsius → Fahrenheit
	expectedFahrenheit := 45.0*1.8 + 32.0
	actualFahrenheit, ok := alert.Fields["tempFahrenheit"].(float64)
	require.True(t, ok, "❌ tempFahrenheit devrait être un float64")
	assert.InDelta(t, expectedFahrenheit, actualFahrenheit, 0.01, "❌ La conversion C→F est incorrecte")

	t.Log("✅ Test E2E avec expressions réussi")
}

// TestE2E_InlineFact_NestedReferences teste les références imbriquées
func TestE2E_InlineFact_NestedReferences(t *testing.T) {
	t.Log("🧪 TEST E2E - RÉFÉRENCES MULTIPLES AUX VARIABLES")
	t.Log("=================================================")

	program := `
		type Sensor(#sensorId: string, temperature: number)
		type Threshold(sensorType: string, maxTemp: number)
		type Alert(sensorId: string, temp: number, threshold: number, excess: number)
		
		rule over_threshold: {s: Sensor, th: Threshold} / s.temperature > th.maxTemp ==>
			Xuple("alerts", Alert(
				sensorId: s.sensorId,
				temp: s.temperature,
				threshold: th.maxTemp,
				excess: s.temperature - th.maxTemp
			))
		
		Sensor(sensorId: "S001", temperature: 45.0)
		Threshold(sensorType: "standard", maxTemp: 40.0)
	`

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.tsd")
	err := os.WriteFile(testFile, []byte(program), 0644)
	require.NoError(t, err)

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	createdXuples := make([]*Fact, 0)
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		createdXuples = append(createdXuples, fact)
		return nil
	})

	pipeline := NewConstraintPipeline()
	network, _, err = pipeline.IngestFile(testFile, network, storage)
	require.NoError(t, err, "❌ Ingestion du programme échouée")

	require.Len(t, createdXuples, 1, "❌ Devrait avoir créé 1 xuple")

	alert := createdXuples[0]
	assert.Equal(t, "Alert", alert.Type)
	assert.Equal(t, "S001", alert.Fields["sensorId"])
	assert.Equal(t, 45.0, alert.Fields["temp"])
	assert.Equal(t, 40.0, alert.Fields["threshold"])
	assert.Equal(t, 5.0, alert.Fields["excess"], "❌ L'excès devrait être 45 - 40 = 5")

	t.Log("✅ Test E2E avec références multiples réussi")
}
