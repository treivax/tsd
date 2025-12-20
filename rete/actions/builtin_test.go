// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// mockXupleManager est un mock simple du XupleManager pour les tests
type mockXupleManager struct {
	lastXuplespace      string
	lastFact            *rete.Fact
	lastTriggeringFacts []*rete.Fact
	shouldError         bool
}

// Vérifier que mockXupleManager implémente xuples.XupleManager
var _ xuples.XupleManager = (*mockXupleManager)(nil)

func (m *mockXupleManager) CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
	if m.shouldError {
		return fmt.Errorf("mock error")
	}
	m.lastXuplespace = xuplespace
	m.lastFact = fact
	m.lastTriggeringFacts = triggeringFacts
	return nil
}

func (m *mockXupleManager) CreateXupleSpace(name string, config xuples.XupleSpaceConfig) error {
	return nil
}

func (m *mockXupleManager) GetXupleSpace(name string) (xuples.XupleSpace, error) {
	return nil, nil
}

func (m *mockXupleManager) ListXupleSpaces() []string {
	return []string{}
}

func (m *mockXupleManager) Close() error {
	return nil
}

func TestNewBuiltinActionExecutor(t *testing.T) {
	t.Log("🧪 TEST NewBuiltinActionExecutor")

	network := &rete.ReteNetwork{}
	xupleMgr := &mockXupleManager{}
	output := &bytes.Buffer{}
	logger := log.New(&bytes.Buffer{}, "", 0)

	executor := NewBuiltinActionExecutor(network, xupleMgr, output, logger)
	if executor == nil {
		t.Fatal("❌ NewBuiltinActionExecutor returned nil")
	}

	executor2 := NewBuiltinActionExecutor(network, nil, nil, nil)
	if executor2.output == nil {
		t.Error("❌ Output should default")
	}

	t.Log("✅ NewBuiltinActionExecutor OK")
}

func TestExecutePrint(t *testing.T) {
	t.Log("🧪 TEST Print")

	output := &bytes.Buffer{}
	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, output, nil)

	err := executor.executePrint([]interface{}{"Hello"})
	if err != nil || output.String() != "Hello\n" {
		t.Errorf("❌ Print failed")
	}

	t.Log("✅ Print OK")
}

func TestExecuteLog(t *testing.T) {
	t.Log("🧪 TEST Log")

	logOutput := &bytes.Buffer{}
	logger := log.New(logOutput, "", 0)
	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, logger)

	err := executor.executeLog([]interface{}{"test"})
	if err != nil || !strings.Contains(logOutput.String(), "[TSD] test") {
		t.Errorf("❌ Log failed")
	}

	t.Log("✅ Log OK")
}

func TestExecute_AllActions(t *testing.T) {
	t.Log("🧪 TEST Execute")

	output := &bytes.Buffer{}
	logOutput := &bytes.Buffer{}
	logger := log.New(logOutput, "", 0)
	xupleMgr := &mockXupleManager{}
	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, xupleMgr, output, logger)
	token := &rete.Token{}

	if err := executor.Execute("Print", []interface{}{"test"}, token); err != nil {
		t.Error("❌ Print failed")
	}
	if err := executor.Execute("Log", []interface{}{"test"}, token); err != nil {
		t.Error("❌ Log failed")
	}
	if err := executor.Execute("Xuple", []interface{}{"space", &rete.Fact{}}, token); err != nil {
		t.Error("❌ Xuple failed")
	}
	if err := executor.Execute("Unknown", []interface{}{}, token); err == nil {
		t.Error("❌ Should fail on unknown action")
	}

	t.Log("✅ Execute OK")
}

func TestExtractTriggeringFacts(t *testing.T) {
	t.Log("🧪 TEST extractTriggeringFacts")

	fact1 := &rete.Fact{Type: "T1", ID: "1"}
	fact2 := &rete.Fact{Type: "T2", ID: "2"}

	token := &rete.Token{
		Facts: []*rete.Fact{fact2},
		Parent: &rete.Token{
			Facts: []*rete.Fact{fact1},
		},
	}

	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)
	facts := executor.extractTriggeringFacts(token)

	if len(facts) != 2 || facts[0] != fact1 || facts[1] != fact2 {
		t.Error("❌ Extract failed or wrong order")
	}

	t.Log("✅ extractTriggeringFacts OK")
}

// Tests pour améliorer la couverture

func TestExecutePrint_InvalidArgs(t *testing.T) {
	t.Log("🧪 TEST Print invalid args")

	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)

	// Test: nombre d'arguments incorrect
	if err := executor.executePrint([]interface{}{}); err == nil {
		t.Error("❌ Should fail with no arguments")
	}

	if err := executor.executePrint([]interface{}{"a", "b"}); err == nil {
		t.Error("❌ Should fail with too many arguments")
	}

	// Test: type d'argument incorrect
	if err := executor.executePrint([]interface{}{123}); err == nil {
		t.Error("❌ Should fail with non-string argument")
	}

	t.Log("✅ Print validation OK")
}

func TestExecuteLog_InvalidArgs(t *testing.T) {
	t.Log("🧪 TEST Log invalid args")

	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)

	// Test: nombre d'arguments incorrect
	if err := executor.executeLog([]interface{}{}); err == nil {
		t.Error("❌ Should fail with no arguments")
	}

	// Test: type d'argument incorrect
	if err := executor.executeLog([]interface{}{123}); err == nil {
		t.Error("❌ Should fail with non-string argument")
	}

	t.Log("✅ Log validation OK")
}

func TestExecuteUpdate_Implemented(t *testing.T) {
	t.Log("🧪 TEST Update implemented")

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	executor := NewBuiltinActionExecutor(network, nil, nil, nil)

	// Setup: Ajouter un fait initial
	initialFact := &rete.Fact{
		ID:   "TestType~test",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
		},
	}
	storage.AddFact(initialFact)

	// Test: mise à jour réussie
	updatedFact := &rete.Fact{
		ID:   "TestType~test",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 100,
		},
	}
	err := executor.executeUpdate([]interface{}{updatedFact})
	if err != nil {
		t.Errorf("❌ Update should succeed, got error: %v", err)
	}

	// Vérifier la mise à jour
	storedFact := storage.GetFact("TestType~test")
	if storedFact == nil || storedFact.Fields["value"] != 100 {
		t.Error("❌ Fact should be updated")
	}

	// Test: validation des arguments
	if err := executor.executeUpdate([]interface{}{}); err == nil {
		t.Error("❌ Should fail with no arguments")
	}

	if err := executor.executeUpdate([]interface{}{"wrong"}); err == nil {
		t.Error("❌ Should fail with wrong type")
	}

	if err := executor.executeUpdate([]interface{}{nil}); err == nil {
		t.Error("❌ Should fail with nil fact")
	}

	t.Log("✅ Update validation OK")
}

func TestExecuteInsert_Implemented(t *testing.T) {
	t.Log("🧪 TEST Insert implemented")

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	executor := NewBuiltinActionExecutor(network, nil, nil, nil)

	// Test: insertion réussie
	newFact := &rete.Fact{
		ID:   "TestType~test",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
		},
	}
	err := executor.executeInsert([]interface{}{newFact})
	if err != nil {
		t.Errorf("❌ Insert should succeed, got error: %v", err)
	}

	// Vérifier l'insertion
	storedFact := storage.GetFact("TestType~test")
	if storedFact == nil {
		t.Error("❌ Fact should be inserted")
	}

	// Test: validation des arguments
	if err := executor.executeInsert([]interface{}{}); err == nil {
		t.Error("❌ Should fail with no arguments")
	}

	if err := executor.executeInsert([]interface{}{"wrong"}); err == nil {
		t.Error("❌ Should fail with wrong type")
	}

	if err := executor.executeInsert([]interface{}{nil}); err == nil {
		t.Error("❌ Should fail with nil fact")
	}

	t.Log("✅ Insert validation OK")
}

func TestExecuteRetract_Implemented(t *testing.T) {
	t.Log("🧪 TEST Retract implemented")

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	executor := NewBuiltinActionExecutor(network, nil, nil, nil)

	// Setup: Ajouter un fait
	fact := &rete.Fact{
		ID:   "TestType~test",
		Type: "TestType",
		Fields: map[string]interface{}{
			"value": 42,
		},
	}
	storage.AddFact(fact)

	// Test: rétractation réussie
	err := executor.executeRetract([]interface{}{"TestType~test"})
	if err != nil {
		t.Errorf("❌ Retract should succeed, got error: %v", err)
	}

	// Vérifier la suppression
	storedFact := storage.GetFact("TestType~test")
	if storedFact != nil {
		t.Error("❌ Fact should be removed")
	}

	// Test: validation des arguments
	if err := executor.executeRetract([]interface{}{}); err == nil {
		t.Error("❌ Should fail with no arguments")
	}

	if err := executor.executeRetract([]interface{}{123}); err == nil {
		t.Error("❌ Should fail with wrong type")
	}

	if err := executor.executeRetract([]interface{}{""}); err == nil {
		t.Error("❌ Should fail with empty ID")
	}

	t.Log("✅ Retract validation OK")
}

func TestExecuteXuple_InvalidArgs(t *testing.T) {
	t.Log("🧪 TEST Xuple invalid args")

	xupleMgr := &mockXupleManager{}
	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, xupleMgr, nil, nil)
	token := &rete.Token{}

	// Test: nombre d'arguments incorrect
	if err := executor.executeXuple([]interface{}{}, token); err == nil {
		t.Error("❌ Should fail with no arguments")
	}

	if err := executor.executeXuple([]interface{}{"space"}, token); err == nil {
		t.Error("❌ Should fail with only one argument")
	}

	// Test: type d'argument incorrect pour xuplespace
	if err := executor.executeXuple([]interface{}{123, &rete.Fact{}}, token); err == nil {
		t.Error("❌ Should fail with non-string xuplespace")
	}

	// Test: xuplespace vide
	if err := executor.executeXuple([]interface{}{"", &rete.Fact{}}, token); err == nil {
		t.Error("❌ Should fail with empty xuplespace")
	}

	// Test: type d'argument incorrect pour fact
	if err := executor.executeXuple([]interface{}{"space", "wrong"}, token); err == nil {
		t.Error("❌ Should fail with non-fact argument")
	}

	// Test: fact nil
	if err := executor.executeXuple([]interface{}{"space", nil}, token); err == nil {
		t.Error("❌ Should fail with nil fact")
	}

	// Test: XupleManager nil
	executor2 := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)
	if err := executor2.executeXuple([]interface{}{"space", &rete.Fact{}}, token); err == nil {
		t.Error("❌ Should fail with nil XupleManager")
	}

	// Test: XupleManager retourne une erreur
	xupleMgr.shouldError = true
	if err := executor.executeXuple([]interface{}{"space", &rete.Fact{}}, token); err == nil {
		t.Error("❌ Should propagate XupleManager error")
	}

	t.Log("✅ Xuple validation OK")
}

func TestExtractTriggeringFacts_EmptyToken(t *testing.T) {
	t.Log("🧪 TEST extractTriggeringFacts empty")

	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)

	// Test: token nil
	facts := executor.extractTriggeringFacts(nil)
	if len(facts) != 0 {
		t.Error("❌ Should return empty slice for nil token")
	}

	// Test: token sans faits
	token := &rete.Token{}
	facts = executor.extractTriggeringFacts(token)
	if len(facts) != 0 {
		t.Error("❌ Should return empty slice for token without facts")
	}

	t.Log("✅ extractTriggeringFacts empty OK")
}

func TestSetOutput(t *testing.T) {
	t.Log("🧪 TEST SetOutput")

	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)
	newOutput := &bytes.Buffer{}

	executor.SetOutput(newOutput)
	if err := executor.executePrint([]interface{}{"test"}); err != nil {
		t.Error("❌ Print failed after SetOutput")
	}

	if newOutput.String() != "test\n" {
		t.Error("❌ Output not written to new writer")
	}

	// Test: SetOutput avec nil ne doit pas planter
	executor.SetOutput(nil)

	t.Log("✅ SetOutput OK")
}

func TestSetLogger(t *testing.T) {
	t.Log("🧪 TEST SetLogger")

	executor := NewBuiltinActionExecutor(&rete.ReteNetwork{}, nil, nil, nil)
	newLogOutput := &bytes.Buffer{}
	newLogger := log.New(newLogOutput, "", 0)

	executor.SetLogger(newLogger)
	if err := executor.executeLog([]interface{}{"test"}); err != nil {
		t.Error("❌ Log failed after SetLogger")
	}

	if !strings.Contains(newLogOutput.String(), "test") {
		t.Error("❌ Log not written to new logger")
	}

	// Test: SetLogger avec nil ne doit pas planter
	executor.SetLogger(nil)

	t.Log("✅ SetLogger OK")
}
