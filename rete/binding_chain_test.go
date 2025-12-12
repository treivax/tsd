// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
)
// ============================================================================
// Tests de Création
// ============================================================================
func TestBindingChain_CreateEmpty(t *testing.T) {
	t.Log("🧪 TEST: Création d'une chaîne vide")
	t.Log("=====================================")
	chain := NewBindingChain()
	if chain != nil {
		t.Errorf("❌ Chaîne vide devrait être nil, got %v", chain)
	} else {
		t.Log("✅ Chaîne vide correctement représentée par nil")
	}
	// Les opérations sur nil doivent fonctionner
	if chain.Len() != 0 {
		t.Errorf("❌ Len() sur chaîne vide devrait retourner 0, got %d", chain.Len())
	}
	if chain.Variables() == nil || len(chain.Variables()) != 0 {
		t.Errorf("❌ Variables() sur chaîne vide devrait retourner slice vide, got %v", chain.Variables())
	}
	t.Log("✅ Test réussi: Création de chaîne vide")
}
func TestBindingChain_CreateWithBinding(t *testing.T) {
	t.Log("🧪 TEST: Création avec binding initial")
	t.Log("======================================")
	fact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{"name": "Alice"}}
	chain := NewBindingChainWith("u", fact)
	if chain == nil {
		t.Fatal("❌ Chaîne ne devrait pas être nil")
	}
	if chain.Variable != "u" {
		t.Errorf("❌ Variable devrait être 'u', got '%s'", chain.Variable)
	}
	if chain.Fact != fact {
		t.Errorf("❌ Fact devrait être %v, got %v", fact, chain.Fact)
	}
	if chain.Parent != nil {
		t.Errorf("❌ Parent devrait être nil, got %v", chain.Parent)
	}
	t.Log("✅ Test réussi: Création avec binding initial")
}
// ============================================================================
// Tests d'Ajout
// ============================================================================
func TestBindingChain_Add_Single(t *testing.T) {
	t.Log("🧪 TEST: Ajout d'un binding unique")
	t.Log("==================================")
	fact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	empty := NewBindingChain()
	chain := empty.Add("u", fact)
	if chain == nil {
		t.Fatal("❌ Chaîne résultante ne devrait pas être nil")
	}
	if chain.Variable != "u" {
		t.Errorf("❌ Variable devrait être 'u', got '%s'", chain.Variable)
	}
	if chain.Fact != fact {
		t.Errorf("❌ Fact devrait être %v, got %v", fact, chain.Fact)
	}
	if chain.Parent != nil {
		t.Errorf("❌ Parent devrait être nil (car ajouté sur chaîne vide), got %v", chain.Parent)
	}
	t.Log("✅ Test réussi: Ajout d'un binding unique")
}
func TestBindingChain_Add_Multiple(t *testing.T) {
	t.Log("🧪 TEST: Ajout de bindings multiples")
	t.Log("====================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	taskFact := &Fact{ID: "T001", Type: "Task", Fields: map[string]interface{}{}}
	chain := NewBindingChain()
	chain = chain.Add("u", userFact)
	chain = chain.Add("order", orderFact)
	chain = chain.Add("task", taskFact)
	if chain.Len() != 3 {
		t.Errorf("❌ Len devrait être 3, got %d", chain.Len())
	}
	if chain.Get("u") != userFact {
		t.Errorf("❌ Get('u') devrait retourner userFact")
	}
	if chain.Get("order") != orderFact {
		t.Errorf("❌ Get('order') devrait retourner orderFact")
	}
	if chain.Get("task") != taskFact {
		t.Errorf("❌ Get('task') devrait retourner taskFact")
	}
	t.Log("✅ Test réussi: Ajout de bindings multiples")
}
func TestBindingChain_Add_Preserves_Parent(t *testing.T) {
	t.Log("🧪 TEST: Add préserve la chaîne parente (immutabilité)")
	t.Log("=====================================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	chain1 := NewBindingChain().Add("u", userFact)
	chain2 := chain1.Add("order", orderFact)
	// Vérifier que chain1 n'a pas été modifié
	if chain1.Len() != 1 {
		t.Errorf("❌ chain1.Len devrait rester 1 après ajout dans chain2, got %d", chain1.Len())
	}
	if chain1.Has("order") {
		t.Errorf("❌ chain1 ne devrait pas avoir 'order'")
	}
	// Vérifier que chain2 a les deux bindings
	if chain2.Len() != 2 {
		t.Errorf("❌ chain2.Len devrait être 2, got %d", chain2.Len())
	}
	if !chain2.Has("u") {
		t.Errorf("❌ chain2 devrait avoir 'u'")
	}
	if !chain2.Has("order") {
		t.Errorf("❌ chain2 devrait avoir 'order'")
	}
	// Vérifier le partage structurel
	if chain2.Parent != chain1 {
		t.Errorf("❌ chain2.Parent devrait pointer vers chain1")
	}
	t.Log("✅ Test réussi: Immutabilité préservée")
}
// ============================================================================
// Tests de Lecture
// ============================================================================
func TestBindingChain_Get_Existing(t *testing.T) {
	t.Log("🧪 TEST: Get sur variable existante")
	t.Log("===================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	chain := NewBindingChain()
	chain = chain.Add("u", userFact)
	chain = chain.Add("order", orderFact)
	if chain.Get("u") != userFact {
		t.Errorf("❌ Get('u') devrait retourner userFact")
	}
	if chain.Get("order") != orderFact {
		t.Errorf("❌ Get('order') devrait retourner orderFact")
	}
	t.Log("✅ Test réussi: Get retourne les bons faits")
}
func TestBindingChain_Get_NotFound(t *testing.T) {
	t.Log("🧪 TEST: Get sur variable inexistante")
	t.Log("=====================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	chain := NewBindingChain().Add("u", userFact)
	result := chain.Get("task")
	if result != nil {
		t.Errorf("❌ Get('task') devrait retourner nil, got %v", result)
	}
	t.Log("✅ Test réussi: Get retourne nil pour variable inexistante")
}
func TestBindingChain_Has(t *testing.T) {
	t.Log("🧪 TEST: Has vérifie l'existence de variables")
	t.Log("============================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	chain := NewBindingChain()
	chain = chain.Add("u", userFact)
	chain = chain.Add("order", orderFact)
	if !chain.Has("u") {
		t.Errorf("❌ Has('u') devrait retourner true")
	}
	if !chain.Has("order") {
		t.Errorf("❌ Has('order') devrait retourner true")
	}
	if chain.Has("task") {
		t.Errorf("❌ Has('task') devrait retourner false")
	}
	t.Log("✅ Test réussi: Has fonctionne correctement")
}
func TestBindingChain_Len(t *testing.T) {
	t.Log("🧪 TEST: Len retourne le nombre correct de bindings")
	t.Log("===================================================")
	chain := NewBindingChain()
	if chain.Len() != 0 {
		t.Errorf("❌ Len() sur chaîne vide devrait être 0, got %d", chain.Len())
	}
	chain = chain.Add("u", &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}})
	if chain.Len() != 1 {
		t.Errorf("❌ Len() après 1 ajout devrait être 1, got %d", chain.Len())
	}
	chain = chain.Add("order", &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}})
	if chain.Len() != 2 {
		t.Errorf("❌ Len() après 2 ajouts devrait être 2, got %d", chain.Len())
	}
	chain = chain.Add("task", &Fact{ID: "T001", Type: "Task", Fields: map[string]interface{}{}})
	if chain.Len() != 3 {
		t.Errorf("❌ Len() après 3 ajouts devrait être 3, got %d", chain.Len())
	}
	t.Log("✅ Test réussi: Len retourne le nombre correct")
}
func TestBindingChain_Variables(t *testing.T) {
	t.Log("🧪 TEST: Variables retourne la liste des variables")
	t.Log("==================================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	taskFact := &Fact{ID: "T001", Type: "Task", Fields: map[string]interface{}{}}
	chain := NewBindingChain()
	chain = chain.Add("u", userFact)
	chain = chain.Add("order", orderFact)
	chain = chain.Add("task", taskFact)
	vars := chain.Variables()
	if len(vars) != 3 {
		t.Errorf("❌ Variables() devrait retourner 3 variables, got %d", len(vars))
	}
	expected := []string{"u", "order", "task"}
	for i, v := range expected {
		if i >= len(vars) || vars[i] != v {
			t.Errorf("❌ Variables()[%d] devrait être '%s', got '%s'", i, v, vars[i])
		}
	}
	t.Log("✅ Test réussi: Variables retourne la liste correcte")
}
// ============================================================================
// Tests de Conversion
// ============================================================================
func TestBindingChain_ToMap(t *testing.T) {
	t.Log("🧪 TEST: ToMap convertit la chaîne en map")
	t.Log("=========================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	taskFact := &Fact{ID: "T001", Type: "Task", Fields: map[string]interface{}{}}
	chain := NewBindingChain()
	chain = chain.Add("u", userFact)
	chain = chain.Add("order", orderFact)
	chain = chain.Add("task", taskFact)
	m := chain.ToMap()
	if len(m) != 3 {
		t.Errorf("❌ Map devrait avoir 3 entrées, got %d", len(m))
	}
	if m["u"] != userFact {
		t.Errorf("❌ Map['u'] devrait être userFact")
	}
	if m["order"] != orderFact {
		t.Errorf("❌ Map['order'] devrait être orderFact")
	}
	if m["task"] != taskFact {
		t.Errorf("❌ Map['task'] devrait être taskFact")
	}
	t.Log("✅ Test réussi: ToMap convertit correctement")
}
func TestBindingChain_ToMap_Empty(t *testing.T) {
	t.Log("🧪 TEST: ToMap sur chaîne vide")
	t.Log("==============================")
	chain := NewBindingChain()
	m := chain.ToMap()
	if m == nil {
		t.Errorf("❌ ToMap() ne devrait pas retourner nil")
	}
	if len(m) != 0 {
		t.Errorf("❌ Map devrait être vide, got %d entrées", len(m))
	}
	t.Log("✅ Test réussi: ToMap sur chaîne vide retourne map vide")
}
// ============================================================================
// Tests de Merge
// ============================================================================
func TestBindingChain_Merge(t *testing.T) {
	t.Log("🧪 TEST: Merge combine deux chaînes")
	t.Log("===================================")
	userFact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	orderFact := &Fact{ID: "O001", Type: "Order", Fields: map[string]interface{}{}}
	taskFact := &Fact{ID: "T001", Type: "Task", Fields: map[string]interface{}{}}
	chain1 := NewBindingChain().Add("u", userFact)
	chain2 := NewBindingChain().Add("order", orderFact).Add("task", taskFact)
	merged := chain1.Merge(chain2)
	if merged.Len() != 3 {
		t.Errorf("❌ Chaîne fusionnée devrait avoir 3 bindings, got %d", merged.Len())
	}
	if !merged.Has("u") {
		t.Errorf("❌ Chaîne fusionnée devrait avoir 'u'")
	}
	if !merged.Has("order") {
		t.Errorf("❌ Chaîne fusionnée devrait avoir 'order'")
	}
	if !merged.Has("task") {
		t.Errorf("❌ Chaîne fusionnée devrait avoir 'task'")
	}
	t.Log("✅ Test réussi: Merge combine correctement")
}
func TestBindingChain_Merge_Conflicts(t *testing.T) {
	t.Log("🧪 TEST: Merge avec conflits (priorité à 'other')")
	t.Log("=================================================")
	fact1 := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	fact2 := &Fact{ID: "U002", Type: "User", Fields: map[string]interface{}{}}
	chain1 := NewBindingChain().Add("u", fact1)
	chain2 := NewBindingChain().Add("u", fact2)
	merged := chain1.Merge(chain2)
	if merged.Get("u") != fact2 {
		t.Errorf("❌ En cas de conflit, priorité devrait être à 'other' (fact2)")
	}
	t.Log("✅ Test réussi: Merge gère les conflits correctement")
}
// ============================================================================
// Tests de Edge Cases
// ============================================================================
func TestBindingChain_Nil_Operations(t *testing.T) {
	t.Log("🧪 TEST: Opérations sur chaîne nil")
	t.Log("==================================")
	var chain *BindingChain = nil
	// Len sur nil
	if chain.Len() != 0 {
		t.Errorf("❌ Len() sur nil devrait retourner 0, got %d", chain.Len())
	}
	// Get sur nil
	if chain.Get("any") != nil {
		t.Errorf("❌ Get() sur nil devrait retourner nil")
	}
	// Has sur nil
	if chain.Has("any") {
		t.Errorf("❌ Has() sur nil devrait retourner false")
	}
	// Variables sur nil
	vars := chain.Variables()
	if vars == nil || len(vars) != 0 {
		t.Errorf("❌ Variables() sur nil devrait retourner slice vide")
	}
	// ToMap sur nil
	m := chain.ToMap()
	if m == nil || len(m) != 0 {
		t.Errorf("❌ ToMap() sur nil devrait retourner map vide")
	}
	// String sur nil
	s := chain.String()
	if s != "BindingChain{}" {
		t.Errorf("❌ String() sur nil devrait retourner 'BindingChain{}', got '%s'", s)
	}
	t.Log("✅ Test réussi: Toutes les opérations sur nil sont sûres")
}
func TestBindingChain_Long_Chain(t *testing.T) {
	t.Log("🧪 TEST: Chaîne longue (100 bindings)")
	t.Log("=====================================")
	chain := NewBindingChain()
	// Ajouter 100 bindings
	for i := 0; i < 100; i++ {
		varName := "var" + string(rune('A'+i%26)) + string(rune('0'+i/26))
		fact := &Fact{ID: varName, Type: "Test", Fields: map[string]interface{}{}}
		chain = chain.Add(varName, fact)
	}
	if chain.Len() != 100 {
		t.Errorf("❌ Len devrait être 100, got %d", chain.Len())
	}
	// Vérifier que les variables sont toutes accessibles
	vars := chain.Variables()
	if len(vars) != 100 {
		t.Errorf("❌ Variables() devrait retourner 100 variables, got %d", len(vars))
	}
	// Vérifier un accès au milieu
	if !chain.Has("varA2") {
		t.Errorf("❌ Devrait trouver la variable 'varA2'")
	}
	t.Log("✅ Test réussi: Chaîne longue fonctionne correctement")
}
// ============================================================================
// Tests de Benchmarks
// ============================================================================
func BenchmarkBindingChain_Add(b *testing.B) {
	fact := &Fact{ID: "U001", Type: "User", Fields: map[string]interface{}{}}
	chain := NewBindingChain()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain = chain.Add("u", fact)
	}
}
func BenchmarkBindingChain_Get(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		varName := "var" + string(rune('A'+i))
		fact := &Fact{ID: varName, Type: "Test", Fields: map[string]interface{}{}}
		chain = chain.Add(varName, fact)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Get("varE") // Au milieu de la chaîne
	}
}
func BenchmarkBindingChain_Get_DeepChain(b *testing.B) {
	chain := NewBindingChain()
	// Créer une chaîne de 100 bindings
	for i := 0; i < 100; i++ {
		varName := "var" + string(rune('A'+i%26)) + string(rune('0'+i/26))
		fact := &Fact{ID: varName, Type: "Test", Fields: map[string]interface{}{}}
		chain = chain.Add(varName, fact)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Get("varA2") // Au milieu de la chaîne
	}
}
func BenchmarkBindingChain_Variables(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		varName := "var" + string(rune('A'+i))
		fact := &Fact{ID: varName, Type: "Test", Fields: map[string]interface{}{}}
		chain = chain.Add(varName, fact)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Variables()
	}
}
func BenchmarkBindingChain_ToMap(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		varName := "var" + string(rune('A'+i))
		fact := &Fact{ID: varName, Type: "Test", Fields: map[string]interface{}{}}
		chain = chain.Add(varName, fact)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.ToMap()
	}
}