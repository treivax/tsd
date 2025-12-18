// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"testing"
	"time"
)

// Tests for Selection Policies

func TestFIFOSelectionPolicy(t *testing.T) {
	t.Log("🧪 TEST FIFO SELECTION POLICY")

	policy := NewFIFOSelectionPolicy()

	// Test avec liste vide
	selected := policy.Select([]*Xuple{})
	if selected != nil {
		t.Error("❌ Select sur liste vide devrait retourner nil")
	}

	// Test avec plusieurs xuples
	now := time.Now()
	x1 := &Xuple{ID: "x1", CreatedAt: now.Add(-2 * time.Hour)}
	x2 := &Xuple{ID: "x2", CreatedAt: now.Add(-1 * time.Hour)}
	x3 := &Xuple{ID: "x3", CreatedAt: now}

	xuples := []*Xuple{x2, x3, x1} // Ordre mélangé

	selected = policy.Select(xuples)
	if selected == nil {
		t.Fatal("❌ Selected ne devrait pas être nil")
	}

	if selected.ID != "x1" {
		t.Errorf("❌ Devrait sélectionner le plus ancien (x1), reçu %s", selected.ID)
	}

	if policy.Name() != "fifo" {
		t.Errorf("❌ Name devrait être 'fifo', reçu '%s'", policy.Name())
	}

	t.Log("✅ FIFO selection policy fonctionne correctement")
}

func TestLIFOSelectionPolicy(t *testing.T) {
	t.Log("🧪 TEST LIFO SELECTION POLICY")

	policy := NewLIFOSelectionPolicy()

	// Test avec liste vide
	selected := policy.Select([]*Xuple{})
	if selected != nil {
		t.Error("❌ Select sur liste vide devrait retourner nil")
	}

	// Test avec plusieurs xuples
	now := time.Now()
	x1 := &Xuple{ID: "x1", CreatedAt: now.Add(-2 * time.Hour)}
	x2 := &Xuple{ID: "x2", CreatedAt: now.Add(-1 * time.Hour)}
	x3 := &Xuple{ID: "x3", CreatedAt: now}

	xuples := []*Xuple{x2, x3, x1} // Ordre mélangé

	selected = policy.Select(xuples)
	if selected == nil {
		t.Fatal("❌ Selected ne devrait pas être nil")
	}

	if selected.ID != "x3" {
		t.Errorf("❌ Devrait sélectionner le plus récent (x3), reçu %s", selected.ID)
	}

	if policy.Name() != "lifo" {
		t.Errorf("❌ Name devrait être 'lifo', reçu '%s'", policy.Name())
	}

	t.Log("✅ LIFO selection policy fonctionne correctement")
}

func TestRandomSelectionPolicy(t *testing.T) {
	t.Log("🧪 TEST RANDOM SELECTION POLICY")

	policy := NewRandomSelectionPolicy()

	// Test avec liste vide
	selected := policy.Select([]*Xuple{})
	if selected != nil {
		t.Error("❌ Select sur liste vide devrait retourner nil")
	}

	// Test avec plusieurs xuples
	xuples := []*Xuple{
		{ID: "x1"},
		{ID: "x2"},
		{ID: "x3"},
		{ID: "x4"},
		{ID: "x5"},
	}

	// Vérifier que ça retourne bien un des xuples
	selected = policy.Select(xuples)
	if selected == nil {
		t.Fatal("❌ Selected ne devrait pas être nil")
	}

	// Vérifier que c'est bien un des xuples de la liste
	found := false
	for _, x := range xuples {
		if x.ID == selected.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("❌ Xuple sélectionné (%s) n'est pas dans la liste", selected.ID)
	}

	if policy.Name() != "random" {
		t.Errorf("❌ Name devrait être 'random', reçu '%s'", policy.Name())
	}

	// Test du caractère aléatoire (statistique)
	counts := make(map[string]int)
	for i := 0; i < 1000; i++ {
		selected = policy.Select(xuples)
		counts[selected.ID]++
	}

	// Chaque xuple devrait être sélectionné au moins quelques fois
	for _, x := range xuples {
		if counts[x.ID] == 0 {
			t.Errorf("❌ Xuple %s jamais sélectionné (pas vraiment aléatoire)", x.ID)
		}
	}

	t.Log("✅ Random selection policy fonctionne correctement")
}

// Tests for Consumption Policies

func TestOnceConsumptionPolicy(t *testing.T) {
	t.Log("🧪 TEST ONCE CONSUMPTION POLICY")

	policy := NewOnceConsumptionPolicy()

	xuple := &Xuple{
		ID: "x1",
		Metadata: XupleMetadata{
			ConsumptionCount: 0,
			ConsumedBy:       make(map[string]time.Time),
		},
	}

	// Première consommation : OK
	if !policy.CanConsume(xuple, "agent1") {
		t.Error("❌ Devrait pouvoir consommer la première fois")
	}

	// OnConsumed devrait retourner true (complètement consommé)
	if !policy.OnConsumed(xuple, "agent1") {
		t.Error("❌ OnConsumed devrait retourner true pour policy Once")
	}

	// Après consommation : non consommable
	xuple.Metadata.ConsumptionCount = 1
	if policy.CanConsume(xuple, "agent2") {
		t.Error("❌ Ne devrait pas pouvoir consommer après once")
	}

	if policy.Name() != "once" {
		t.Errorf("❌ Name devrait être 'once', reçu '%s'", policy.Name())
	}

	t.Log("✅ Once consumption policy fonctionne correctement")
}

func TestPerAgentConsumptionPolicy(t *testing.T) {
	t.Log("🧪 TEST PER-AGENT CONSUMPTION POLICY")

	policy := NewPerAgentConsumptionPolicy()

	xuple := &Xuple{
		ID: "x1",
		Metadata: XupleMetadata{
			ConsumptionCount: 0,
			ConsumedBy:       make(map[string]time.Time),
		},
	}

	// Agent1 peut consommer
	if !policy.CanConsume(xuple, "agent1") {
		t.Error("❌ Agent1 devrait pouvoir consommer")
	}

	// Marquer consommé par agent1
	xuple.Metadata.ConsumedBy["agent1"] = time.Now()
	xuple.Metadata.ConsumptionCount = 1

	// Agent1 ne peut plus consommer
	if policy.CanConsume(xuple, "agent1") {
		t.Error("❌ Agent1 ne devrait pas pouvoir consommer deux fois")
	}

	// Agent2 peut consommer
	if !policy.CanConsume(xuple, "agent2") {
		t.Error("❌ Agent2 devrait pouvoir consommer")
	}

	// OnConsumed ne devrait jamais marquer comme complètement consommé
	if policy.OnConsumed(xuple, "agent2") {
		t.Error("❌ OnConsumed ne devrait pas marquer comme complètement consommé pour per-agent")
	}

	if policy.Name() != "per-agent" {
		t.Errorf("❌ Name devrait être 'per-agent', reçu '%s'", policy.Name())
	}

	t.Log("✅ Per-agent consumption policy fonctionne correctement")
}

func TestLimitedConsumptionPolicy(t *testing.T) {
	t.Log("🧪 TEST LIMITED CONSUMPTION POLICY")

	policy := NewLimitedConsumptionPolicy(3)

	xuple := &Xuple{
		ID: "x1",
		Metadata: XupleMetadata{
			ConsumptionCount: 0,
			ConsumedBy:       make(map[string]time.Time),
		},
	}

	// Première consommation : OK
	if !policy.CanConsume(xuple, "agent1") {
		t.Error("❌ Devrait pouvoir consommer (1/3)")
	}
	xuple.Metadata.ConsumptionCount = 1

	// Deuxième consommation : OK
	if !policy.CanConsume(xuple, "agent2") {
		t.Error("❌ Devrait pouvoir consommer (2/3)")
	}
	xuple.Metadata.ConsumptionCount = 2

	// Troisième consommation : OK mais atteint limite
	if !policy.CanConsume(xuple, "agent3") {
		t.Error("❌ Devrait pouvoir consommer (3/3)")
	}

	// OnConsumed devrait retourner true (limite atteinte)
	xuple.Metadata.ConsumptionCount = 3
	if !policy.OnConsumed(xuple, "agent3") {
		t.Error("❌ OnConsumed devrait retourner true quand limite atteinte")
	}

	// Quatrième consommation : NON
	if policy.CanConsume(xuple, "agent4") {
		t.Error("❌ Ne devrait pas pouvoir consommer après limite")
	}

	if policy.Name() != "limited" {
		t.Errorf("❌ Name devrait être 'limited', reçu '%s'", policy.Name())
	}

	// Test avec limite invalide (devrait mettre 1)
	policy2 := NewLimitedConsumptionPolicy(0)
	if policy2.MaxConsumptions != 1 {
		t.Errorf("❌ Limite invalide devrait être corrigée à 1, reçu %d", policy2.MaxConsumptions)
	}

	t.Log("✅ Limited consumption policy fonctionne correctement")
}

// Tests for Retention Policies

func TestUnlimitedRetentionPolicy(t *testing.T) {
	t.Log("🧪 TEST UNLIMITED RETENTION POLICY")

	policy := NewUnlimitedRetentionPolicy()

	now := time.Now()
	expiration := policy.ComputeExpiration(now)

	// Devrait retourner zero time
	if !expiration.IsZero() {
		t.Error("❌ ComputeExpiration devrait retourner zero time")
	}

	xuple := &Xuple{
		ID: "x1",
		Metadata: XupleMetadata{
			ExpiresAt: time.Time{},
		},
	}

	// Devrait toujours retenir
	if !policy.ShouldRetain(xuple) {
		t.Error("❌ ShouldRetain devrait toujours retourner true")
	}

	if policy.Name() != "unlimited" {
		t.Errorf("❌ Name devrait être 'unlimited', reçu '%s'", policy.Name())
	}

	t.Log("✅ Unlimited retention policy fonctionne correctement")
}

func TestDurationRetentionPolicy(t *testing.T) {
	t.Log("🧪 TEST DURATION RETENTION POLICY")

	duration := 100 * time.Millisecond
	policy := NewDurationRetentionPolicy(duration)

	now := time.Now()
	expiration := policy.ComputeExpiration(now)

	// Expiration devrait être now + duration
	expectedExpiration := now.Add(duration)
	diff := expiration.Sub(expectedExpiration)
	if diff > 1*time.Millisecond || diff < -1*time.Millisecond {
		t.Errorf("❌ Expiration incorrecte, différence: %v", diff)
	}

	// Test ShouldRetain avec xuple non expiré
	xuple := &Xuple{
		ID: "x1",
		Metadata: XupleMetadata{
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}

	if !policy.ShouldRetain(xuple) {
		t.Error("❌ Devrait retenir xuple non expiré")
	}

	// Test ShouldRetain avec xuple expiré
	xuple.Metadata.ExpiresAt = time.Now().Add(-1 * time.Hour)
	if policy.ShouldRetain(xuple) {
		t.Error("❌ Ne devrait pas retenir xuple expiré")
	}

	// Test ShouldRetain avec xuple sans expiration
	xuple.Metadata.ExpiresAt = time.Time{}
	if !policy.ShouldRetain(xuple) {
		t.Error("❌ Devrait retenir xuple sans expiration")
	}

	if policy.Name() != "duration" {
		t.Errorf("❌ Name devrait être 'duration', reçu '%s'", policy.Name())
	}

	// Test avec durée invalide (devrait mettre 1 heure)
	policy2 := NewDurationRetentionPolicy(0)
	if policy2.Duration != 1*time.Hour {
		t.Errorf("❌ Durée invalide devrait être corrigée à 1h, reçu %v", policy2.Duration)
	}

	t.Log("✅ Duration retention policy fonctionne correctement")
}

// Test PolicyType String

func TestPolicyTypeString(t *testing.T) {
	t.Log("🧪 TEST POLICY TYPE STRING")

	tests := []struct {
		policyType PolicyType
		expected   string
	}{
		{PolicyTypeSelection, "selection"},
		{PolicyTypeConsumption, "consumption"},
		{PolicyTypeRetention, "retention"},
		{PolicyType(999), "unknown"},
	}

	for _, tt := range tests {
		result := tt.policyType.String()
		if result != tt.expected {
			t.Errorf("❌ PolicyType(%d).String() = '%s', attendu '%s'",
				tt.policyType, result, tt.expected)
		}
	}

	t.Log("✅ PolicyType.String() fonctionne correctement")
}
