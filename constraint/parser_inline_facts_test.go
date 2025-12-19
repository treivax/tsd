// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_InlineFact_Simple teste le parsing d'un fait inline simple dans une action
func TestParser_InlineFact_Simple(t *testing.T) {
	t.Log("🧪 TEST PARSER - INLINE FACT SIMPLE")
	t.Log("====================================")

	input := `
		type Sensor(id: string, temp: number)
		type Alert(level: string, id: string)
		
		rule test: {s: Sensor} / s.temp > 40.0 ==> 
			Xuple("alerts", Alert(level: "HIGH", id: "A001"))
	`

	result, err := Parse("test.tsd", []byte(input))
	require.NoError(t, err, "❌ Le parsing a échoué")
	require.NotNil(t, result, "❌ Le résultat est nil")

	resultMap, ok := result.(map[string]interface{})
	require.True(t, ok, "❌ Le résultat n'est pas une map")

	// Vérifier les expressions (rules)
	expressions := resultMap["expressions"].([]interface{})
	require.Len(t, expressions, 1, "❌ Devrait avoir 1 règle")

	ruleMap := expressions[0].(map[string]interface{})
	assert.Equal(t, "expression", ruleMap["type"])
	assert.Equal(t, "test", ruleMap["ruleId"])

	// Vérifier l'action
	actionMap := ruleMap["action"].(map[string]interface{})
	assert.Equal(t, "action", actionMap["type"])

	jobs := actionMap["jobs"].([]interface{})
	require.Len(t, jobs, 1, "❌ Devrait avoir 1 job")

	jobMap := jobs[0].(map[string]interface{})
	assert.Equal(t, "jobCall", jobMap["type"])
	assert.Equal(t, "Xuple", jobMap["name"])

	// Vérifier les arguments
	args := jobMap["args"].([]interface{})
	require.Len(t, args, 2, "❌ Devrait avoir 2 arguments")

	// Premier argument: string "alerts"
	arg0Map := args[0].(map[string]interface{})
	assert.Equal(t, "string", arg0Map["type"])
	assert.Equal(t, "alerts", arg0Map["value"])

	// Deuxième argument: fait inline Alert(...)
	arg1Map := args[1].(map[string]interface{})
	require.Equal(t, "inlineFact", arg1Map["type"], "❌ Le deuxième argument devrait être un inlineFact")
	assert.Equal(t, "Alert", arg1Map["typeName"])

	fields := arg1Map["fields"].([]interface{})
	assert.Len(t, fields, 2, "❌ Le fait Alert devrait avoir 2 champs")

	t.Log("✅ Parsing d'un fait inline simple réussi")
}

// TestParser_InlineFact_Multiline teste le parsing d'un fait inline multi-ligne
func TestParser_InlineFact_Multiline(t *testing.T) {
	t.Log("🧪 TEST PARSER - INLINE FACT MULTI-LIGNE")
	t.Log("=========================================")

	input := `
		type Sensor(id: string, temp: number, location: string)
		type Alert(level: string, message: string, sensorId: string, temperature: number)
		
		rule test: {s: Sensor} / s.temp > 40.0 ==> 
			Xuple("alerts", Alert(
				level: "CRITICAL",
				message: "Temperature too high",
				sensorId: "S001",
				temperature: 45.5
			))
	`

	result, err := Parse("test.tsd", []byte(input))
	require.NoError(t, err, "❌ Le parsing a échoué")

	resultMap := result.(map[string]interface{})
	expressions := resultMap["expressions"].([]interface{})
	ruleMap := expressions[0].(map[string]interface{})
	actionMap := ruleMap["action"].(map[string]interface{})
	jobs := actionMap["jobs"].([]interface{})
	jobMap := jobs[0].(map[string]interface{})
	args := jobMap["args"].([]interface{})

	// Vérifier le fait inline
	factMap := args[1].(map[string]interface{})
	assert.Equal(t, "inlineFact", factMap["type"])
	assert.Equal(t, "Alert", factMap["typeName"])

	fields := factMap["fields"].([]interface{})
	assert.Len(t, fields, 4, "❌ Le fait Alert devrait avoir 4 champs")

	// Vérifier les noms de champs
	fieldNames := make(map[string]bool)
	for _, field := range fields {
		fieldMap := field.(map[string]interface{})
		fieldNames[fieldMap["name"].(string)] = true
	}

	assert.True(t, fieldNames["level"], "❌ Devrait avoir le champ 'level'")
	assert.True(t, fieldNames["message"], "❌ Devrait avoir le champ 'message'")
	assert.True(t, fieldNames["sensorId"], "❌ Devrait avoir le champ 'sensorId'")
	assert.True(t, fieldNames["temperature"], "❌ Devrait avoir le champ 'temperature'")

	t.Log("✅ Parsing d'un fait inline multi-ligne réussi")
}

// TestParser_InlineFact_WithFieldReferences teste les références aux champs de variables
func TestParser_InlineFact_WithFieldReferences(t *testing.T) {
	t.Log("🧪 TEST PARSER - INLINE FACT AVEC RÉFÉRENCES AUX CHAMPS")
	t.Log("========================================================")

	input := `
		type Sensor(id: string, sensorId: string, temperature: number)
		type Alert(level: string, sensorId: string, temperature: number)
		
		rule test: {s: Sensor} / s.temperature > 40.0 ==> 
			Xuple("alerts", Alert(
				level: "HIGH",
				sensorId: s.sensorId,
				temperature: s.temperature
			))
	`

	result, err := Parse("test.tsd", []byte(input))
	require.NoError(t, err, "❌ Le parsing a échoué")

	resultMap := result.(map[string]interface{})
	expressions := resultMap["expressions"].([]interface{})
	ruleMap := expressions[0].(map[string]interface{})
	actionMap := ruleMap["action"].(map[string]interface{})
	jobs := actionMap["jobs"].([]interface{})
	jobMap := jobs[0].(map[string]interface{})
	args := jobMap["args"].([]interface{})

	// Vérifier le fait inline
	factMap := args[1].(map[string]interface{})
	assert.Equal(t, "inlineFact", factMap["type"])

	fields := factMap["fields"].([]interface{})
	require.Len(t, fields, 3, "❌ Le fait Alert devrait avoir 3 champs")

	// Trouver les champs sensorId et temperature
	var sensorIdField, temperatureField map[string]interface{}
	for _, field := range fields {
		fieldMap := field.(map[string]interface{})
		if fieldMap["name"].(string) == "sensorId" {
			sensorIdField = fieldMap
		} else if fieldMap["name"].(string) == "temperature" {
			temperatureField = fieldMap
		}
	}

	require.NotNil(t, sensorIdField, "❌ Champ 'sensorId' non trouvé")
	require.NotNil(t, temperatureField, "❌ Champ 'temperature' non trouvé")

	// Vérifier que les valeurs sont des références de champs (fieldAccess)
	sensorIdValue := sensorIdField["value"].(map[string]interface{})
	assert.Equal(t, "fieldAccess", sensorIdValue["type"], "❌ sensorId devrait être un fieldAccess")
	assert.Equal(t, "s", sensorIdValue["object"], "❌ L'objet devrait être 's'")
	assert.Equal(t, "sensorId", sensorIdValue["field"], "❌ Le champ devrait être 'sensorId'")

	temperatureValue := temperatureField["value"].(map[string]interface{})
	assert.Equal(t, "fieldAccess", temperatureValue["type"], "❌ temperature devrait être un fieldAccess")
	assert.Equal(t, "s", temperatureValue["object"], "❌ L'objet devrait être 's'")
	assert.Equal(t, "temperature", temperatureValue["field"], "❌ Le champ devrait être 'temperature'")

	t.Log("✅ Parsing de références aux champs réussi")
}

// TestParser_MultipleActions teste plusieurs actions séparées par des virgules
func TestParser_MultipleActions_WithInlineFacts(t *testing.T) {
	t.Log("🧪 TEST PARSER - ACTIONS MULTIPLES AVEC FAITS INLINE")
	t.Log("=====================================================")

	input := `
		type Sensor(id: string, temperature: number, location: string)
		type Alert(level: string, id: string)
		type Command(action: string, target: string)
		
		rule test: {s: Sensor} / s.temperature > 40.0 ==> 
			Print("Alert!"),
			Xuple("alerts", Alert(level: "HIGH", id: s.id)),
			Xuple("commands", Command(action: "cool", target: s.location))
	`

	result, err := Parse("test.tsd", []byte(input))
	require.NoError(t, err, "❌ Le parsing a échoué")

	resultMap := result.(map[string]interface{})
	expressions := resultMap["expressions"].([]interface{})
	ruleMap := expressions[0].(map[string]interface{})
	actionMap := ruleMap["action"].(map[string]interface{})
	jobs := actionMap["jobs"].([]interface{})

	require.Len(t, jobs, 3, "❌ Devrait avoir 3 jobs")

	// Vérifier les noms d'actions
	job0 := jobs[0].(map[string]interface{})
	assert.Equal(t, "Print", job0["name"])

	job1 := jobs[1].(map[string]interface{})
	assert.Equal(t, "Xuple", job1["name"])

	job2 := jobs[2].(map[string]interface{})
	assert.Equal(t, "Xuple", job2["name"])

	// Vérifier que job1 a un fait inline Alert
	args1 := job1["args"].([]interface{})
	require.Len(t, args1, 2)
	fact1 := args1[1].(map[string]interface{})
	assert.Equal(t, "inlineFact", fact1["type"])
	assert.Equal(t, "Alert", fact1["typeName"])

	// Vérifier que job2 a un fait inline Command
	args2 := job2["args"].([]interface{})
	require.Len(t, args2, 2)
	fact2 := args2[1].(map[string]interface{})
	assert.Equal(t, "inlineFact", fact2["type"])
	assert.Equal(t, "Command", fact2["typeName"])

	t.Log("✅ Parsing d'actions multiples avec faits inline réussi")
}

// TestParser_InlineFact_WithExpressions teste les expressions dans les champs
func TestParser_InlineFact_WithExpressions(t *testing.T) {
	t.Log("🧪 TEST PARSER - INLINE FACT AVEC EXPRESSIONS")
	t.Log("==============================================")

	input := `
		type Sensor(id: string, temperature: number)
		type Alert(level: string, tempCelsius: number, tempFahrenheit: number)
		
		rule test: {s: Sensor} / s.temperature > 40.0 ==> 
			Xuple("alerts", Alert(
				level: "HIGH",
				tempCelsius: s.temperature,
				tempFahrenheit: s.temperature * 1.8 + 32.0
			))
	`

	result, err := Parse("test.tsd", []byte(input))
	require.NoError(t, err, "❌ Le parsing a échoué")

	resultMap := result.(map[string]interface{})
	expressions := resultMap["expressions"].([]interface{})
	ruleMap := expressions[0].(map[string]interface{})
	actionMap := ruleMap["action"].(map[string]interface{})
	jobs := actionMap["jobs"].([]interface{})
	jobMap := jobs[0].(map[string]interface{})
	args := jobMap["args"].([]interface{})

	// Vérifier le fait inline
	factMap := args[1].(map[string]interface{})
	assert.Equal(t, "inlineFact", factMap["type"])

	fields := factMap["fields"].([]interface{})
	require.Len(t, fields, 3)

	// Trouver le champ tempFahrenheit
	var tempFahrenheitField map[string]interface{}
	for _, field := range fields {
		fieldMap := field.(map[string]interface{})
		if fieldMap["name"].(string) == "tempFahrenheit" {
			tempFahrenheitField = fieldMap
			break
		}
	}

	require.NotNil(t, tempFahrenheitField, "❌ Champ 'tempFahrenheit' non trouvé")

	// Vérifier que la valeur est une expression binaire
	value := tempFahrenheitField["value"].(map[string]interface{})
	assert.Equal(t, "binaryOp", value["type"], "❌ La valeur devrait être une expression binaire")

	t.Log("✅ Parsing d'expressions dans les faits inline réussi")
}
