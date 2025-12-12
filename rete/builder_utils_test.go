// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBuilderUtils(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()
	utils := NewBuilderUtils(env.Storage)
	assert.NotNil(t, utils)
	assert.NotNil(t, utils.storage)
	assert.Equal(t, env.Storage, utils.storage)
}
func TestBuilderUtils_CreatePassthroughAlphaNode(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()
	utils := NewBuilderUtils(env.Storage)
	tests := []struct {
		name     string
		ruleID   string
		varName  string
		side     string
		expected string
	}{
		{
			name:     "Without side",
			ruleID:   "rule1",
			varName:  "p",
			side:     "",
			expected: "rule1_pass_p",
		},
		{
			name:     "With left side",
			ruleID:   "rule2",
			varName:  "e",
			side:     "left",
			expected: "rule2_pass_e",
		},
		{
			name:     "With right side",
			ruleID:   "rule3",
			varName:  "d",
			side:     "right",
			expected: "rule3_pass_d",
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			node := utils.CreatePassthroughAlphaNode(tt.ruleID, tt.varName, tt.side)
			assert.NotNil(t, node)
			assert.Equal(t, tt.expected, node.ID)
			assert.Equal(t, tt.varName, node.VariableName)
			// Vérifier la condition
			assert.NotNil(t, node.Condition)
			condMap, ok := node.Condition.(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, ConditionTypePassthrough, condMap["type"])
			if tt.side != "" {
				assert.Equal(t, tt.side, condMap["side"])
			}
		})
	}
}
func TestBuilderUtils_ConnectTypeNodeToBetaNode(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()
	utils := NewBuilderUtils(env.Storage)
	// Créer un TypeNode
	typeDef := TypeDefinition{
		Type:   "type",
		Name:   "Person",
		Fields: []Field{{Name: "age", Type: "number"}},
	}
	typeNode := NewTypeNode("Person", typeDef, env.Storage)
	env.Network.TypeNodes["Person"] = typeNode
	env.Network.RootNode.AddChild(typeNode)
	// Créer un BetaNode mock (JoinNode)
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}
	condition := map[string]interface{}{"type": "comparison"}
	joinNode := NewJoinNode("test_join", condition, leftVars, rightVars, varTypes, env.Storage)
	// Connecter
	utils.ConnectTypeNodeToBetaNode(env.Network, "test_rule", "p", "Person", joinNode, "left")
	// Vérifier que le TypeNode a un enfant (AlphaNode passthrough)
	require.Equal(t, 1, len(typeNode.Children))
	// Vérifier que l'AlphaNode est bien créé
	alphaNode, ok := typeNode.Children[0].(*AlphaNode)
	require.True(t, ok)
	// With per-rule passthroughs, the ID now includes the rule ID and variable name
	assert.Equal(t, "passthrough_test_rule_p_Person_left", alphaNode.ID)
	// Vérifier que l'AlphaNode a le BetaNode comme enfant
	require.Equal(t, 1, len(alphaNode.Children))
	assert.Equal(t, joinNode, alphaNode.Children[0])
	env.AssertNoErrors(t)
}
func TestBuilderUtils_GetStringField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		input        map[string]interface{}
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "Existing string field",
			input:        map[string]interface{}{"name": "Alice"},
			key:          "name",
			defaultValue: "default",
			expected:     "Alice",
		},
		{
			name:         "Missing field",
			input:        map[string]interface{}{"age": 25},
			key:          "name",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Non-string field",
			input:        map[string]interface{}{"age": 25},
			key:          "age",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "Empty string",
			input:        map[string]interface{}{"name": ""},
			key:          "name",
			defaultValue: "default",
			expected:     "",
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GetStringField(tt.input, tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBuilderUtils_GetIntField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		input        map[string]interface{}
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "Existing int field",
			input:        map[string]interface{}{"count": 42},
			key:          "count",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "Existing float64 field",
			input:        map[string]interface{}{"count": 42.0},
			key:          "count",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "Missing field",
			input:        map[string]interface{}{"name": "Alice"},
			key:          "count",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "Non-numeric field",
			input:        map[string]interface{}{"name": "Alice"},
			key:          "name",
			defaultValue: 10,
			expected:     10,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GetIntField(tt.input, tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBuilderUtils_GetBoolField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		input        map[string]interface{}
		key          string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "Existing true field",
			input:        map[string]interface{}{"active": true},
			key:          "active",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "Existing false field",
			input:        map[string]interface{}{"active": false},
			key:          "active",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "Missing field",
			input:        map[string]interface{}{"name": "Alice"},
			key:          "active",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "Non-boolean field",
			input:        map[string]interface{}{"count": 42},
			key:          "count",
			defaultValue: true,
			expected:     true,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GetBoolField(tt.input, tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBuilderUtils_GetMapField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    map[string]interface{}
		key      string
		expected map[string]interface{}
		found    bool
	}{
		{
			name:     "Existing map field",
			input:    map[string]interface{}{"meta": map[string]interface{}{"version": "1.0"}},
			key:      "meta",
			expected: map[string]interface{}{"version": "1.0"},
			found:    true,
		},
		{
			name:     "Missing field",
			input:    map[string]interface{}{"name": "Alice"},
			key:      "meta",
			expected: nil,
			found:    false,
		},
		{
			name:     "Non-map field",
			input:    map[string]interface{}{"count": 42},
			key:      "count",
			expected: nil,
			found:    false,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, found := GetMapField(tt.input, tt.key)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBuilderUtils_GetListField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    map[string]interface{}
		key      string
		expected []interface{}
		found    bool
	}{
		{
			name:     "Existing list field",
			input:    map[string]interface{}{"items": []interface{}{"a", "b", "c"}},
			key:      "items",
			expected: []interface{}{"a", "b", "c"},
			found:    true,
		},
		{
			name:     "Empty list",
			input:    map[string]interface{}{"items": []interface{}{}},
			key:      "items",
			expected: []interface{}{},
			found:    true,
		},
		{
			name:     "Missing field",
			input:    map[string]interface{}{"name": "Alice"},
			key:      "items",
			expected: nil,
			found:    false,
		},
		{
			name:     "Non-list field",
			input:    map[string]interface{}{"count": 42},
			key:      "count",
			expected: nil,
			found:    false,
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, found := GetListField(tt.input, tt.key)
			assert.Equal(t, tt.found, found)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBuilderUtils_CreateTerminalNode(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()
	utils := NewBuilderUtils(env.Storage)
	action := &Action{
		Type: "action",
		Job: &JobCall{
			Type: "jobCall",
			Name: "print",
			Args: []interface{}{"Test message"},
		},
	}
	terminalNode := utils.CreateTerminalNode(env.Network, "test_rule", action)
	require.NotNil(t, terminalNode)
	assert.Equal(t, "test_rule_terminal", terminalNode.ID)
	assert.Equal(t, action, terminalNode.Action)
	// Vérifier que le terminal est enregistré dans le réseau
	assert.Equal(t, terminalNode, env.Network.TerminalNodes[terminalNode.ID])
	// Vérifier l'enregistrement dans le LifecycleManager si disponible
	if env.Network.LifecycleManager != nil {
		t.Log("LifecycleManager present")
	}
	env.AssertNoErrors(t)
}
func TestBuilderUtils_BuildVarTypesMap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		variableNames []string
		variableTypes []string
		expected      map[string]string
	}{
		{
			name:          "Equal length arrays",
			variableNames: []string{"p", "e", "d"},
			variableTypes: []string{"Person", "Employee", "Department"},
			expected: map[string]string{
				"p": "Person",
				"e": "Employee",
				"d": "Department",
			},
		},
		{
			name:          "Types shorter than names",
			variableNames: []string{"p", "e", "d"},
			variableTypes: []string{"Person", "Employee"},
			expected: map[string]string{
				"p": "Person",
				"e": "Employee",
			},
		},
		{
			name:          "Empty arrays",
			variableNames: []string{},
			variableTypes: []string{},
			expected:      map[string]string{},
		},
		{
			name:          "Single variable",
			variableNames: []string{"p"},
			variableTypes: []string{"Person"},
			expected: map[string]string{
				"p": "Person",
			},
		},
	}
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := BuildVarTypesMap(tt.variableNames, tt.variableTypes)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestBuilderUtils_ConnectTypeNodeToBetaNode_TypeNotFound(t *testing.T) {
	t.Parallel()
	env := NewTestEnvironment(t, WithLogLevel(LogLevelSilent))
	defer env.Cleanup()
	utils := NewBuilderUtils(env.Storage)
	// Créer un BetaNode sans créer le TypeNode correspondant
	leftVars := []string{"p"}
	rightVars := []string{"o"}
	varTypes := map[string]string{"p": "Person", "o": "Order"}
	condition := map[string]interface{}{"type": "comparison"}
	joinNode := NewJoinNode("test_join", condition, leftVars, rightVars, varTypes, env.Storage)
	// Essayer de connecter avec un type qui n'existe pas
	// La fonction ne devrait pas paniquer, elle ne fait simplement rien
	utils.ConnectTypeNodeToBetaNode(env.Network, "test_rule", "p", "NonExistent", joinNode, "left")
	// Pas d'assertion - on vérifie juste que ça ne panique pas
	t.Log("No panic when TypeNode not found")
	env.AssertNoErrors(t)
}
