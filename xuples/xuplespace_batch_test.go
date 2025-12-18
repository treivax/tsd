// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"fmt"
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
)

// TestRetrieveMultiple_BasicFunctionality teste le comportement de base de RetrieveMultiple
func TestRetrieveMultiple_BasicFunctionality(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: RetrieveMultiple - Fonctionnalité de Base")
	t.Log("═══════════════════════════════════════════════════════════════")

	tests := []struct {
		name           string
		numXuples      int
		requestCount   int
		policy         ConsumptionPolicy
		expectedCount  int
		expectedRemain int
		wantErr        bool
	}{
		{
			name:           "récupérer moins que disponible",
			numXuples:      10,
			requestCount:   5,
			policy:         NewOnceConsumptionPolicy(),
			expectedCount:  5,
			expectedRemain: 5,
			wantErr:        false,
		},
		{
			name:           "récupérer exactement le nombre disponible",
			numXuples:      5,
			requestCount:   5,
			policy:         NewOnceConsumptionPolicy(),
			expectedCount:  5,
			expectedRemain: 0,
			wantErr:        false,
		},
		{
			name:           "récupérer plus que disponible",
			numXuples:      3,
			requestCount:   5,
			policy:         NewOnceConsumptionPolicy(),
			expectedCount:  3,
			expectedRemain: 0,
			wantErr:        false,
		},
		{
			name:           "récupérer 0 xuples",
			numXuples:      5,
			requestCount:   0,
			policy:         NewOnceConsumptionPolicy(),
			expectedCount:  0,
			expectedRemain: 5,
			wantErr:        false,
		},
		{
			name:           "nombre négatif",
			numXuples:      5,
			requestCount:   -1,
			policy:         NewOnceConsumptionPolicy(),
			expectedCount:  0,
			expectedRemain: 5,
			wantErr:        true,
		},
		{
			name:           "per-agent policy",
			numXuples:      5,
			requestCount:   3,
			policy:         NewPerAgentConsumptionPolicy(),
			expectedCount:  3,
			expectedRemain: 5, // Avec per-agent, les xuples restent disponibles pour d'autres
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := XupleSpaceConfig{
				Name:              "test_batch",
				SelectionPolicy:   NewFIFOSelectionPolicy(),
				ConsumptionPolicy: tt.policy,
				RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
				MaxSize:           0,
			}
			space := NewXupleSpace(config)

			// Insérer les xuples
			for i := 0; i < tt.numXuples; i++ {
				xuple := &Xuple{
					ID: fmt.Sprintf("xuple-%03d", i),
					Fact: &rete.Fact{
						Type: "TestFact",
						Fields: map[string]interface{}{
							"id":    fmt.Sprintf("fact-%d", i),
							"index": i,
						},
					},
					TriggeringFacts: nil,
					CreatedAt:       time.Now(),
					Metadata: XupleMetadata{
						State:      XupleStateAvailable,
						ConsumedBy: make(map[string]time.Time),
					},
				}
				if err := space.Insert(xuple); err != nil {
					t.Fatalf("❌ Insert échoué: %v", err)
				}
			}

			countBefore := space.Count()
			t.Logf("   Count avant RetrieveMultiple: %d", countBefore)

			// Act
			xuples, err := space.RetrieveMultiple("agent1", tt.requestCount)

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("❌ RetrieveMultiple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(xuples) != tt.expectedCount {
					t.Errorf("❌ Nombre de xuples récupérés = %d, attendu %d", len(xuples), tt.expectedCount)
				} else {
					t.Logf("✅ Récupéré %d xuples comme attendu", len(xuples))
				}

				countAfter := space.Count()
				if countAfter != tt.expectedRemain {
					t.Errorf("❌ Count après = %d, attendu %d", countAfter, tt.expectedRemain)
				} else {
					t.Logf("✅ Count après: %d (correct)", countAfter)
				}

				// Vérifier que tous les xuples sont marqués comme consommés
				for i, xuple := range xuples {
					if xuple.Metadata.ConsumptionCount < 1 {
						t.Errorf("❌ Xuple[%d] pas marqué comme consommé", i)
					}
					if _, consumed := xuple.Metadata.ConsumedBy["agent1"]; !consumed {
						t.Errorf("❌ Xuple[%d] agent1 pas dans ConsumedBy", i)
					}
				}
				if len(xuples) > 0 {
					t.Logf("✅ Tous les xuples correctement marqués comme consommés")
				}
			}
		})
	}
}

// TestRetrieveMultiple_SelectionPolicy teste les différentes politiques de sélection
func TestRetrieveMultiple_SelectionPolicy(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: RetrieveMultiple - Politiques de Sélection")
	t.Log("═══════════════════════════════════════════════════════════════")

	const numXuples = 5
	const requestCount = 3

	tests := []struct {
		name            string
		selectionPolicy SelectionPolicy
		expectedOrder   []int // Indices attendus dans l'ordre
	}{
		{
			name:            "FIFO - premiers insérés",
			selectionPolicy: NewFIFOSelectionPolicy(),
			expectedOrder:   []int{0, 1, 2},
		},
		{
			name:            "LIFO - derniers insérés",
			selectionPolicy: NewLIFOSelectionPolicy(),
			expectedOrder:   []int{4, 3, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := XupleSpaceConfig{
				Name:              "test_selection",
				SelectionPolicy:   tt.selectionPolicy,
				ConsumptionPolicy: NewOnceConsumptionPolicy(),
				RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
				MaxSize:           0,
			}
			space := NewXupleSpace(config)

			// Insérer xuples avec indices
			for i := 0; i < numXuples; i++ {
				xuple := &Xuple{
					ID: fmt.Sprintf("xuple-%03d", i),
					Fact: &rete.Fact{
						Type: "TestFact",
						Fields: map[string]interface{}{
							"id":    fmt.Sprintf("fact-%d", i),
							"index": i,
						},
					},
					TriggeringFacts: nil,
					CreatedAt:       time.Now(),
					Metadata: XupleMetadata{
						State:      XupleStateAvailable,
						ConsumedBy: make(map[string]time.Time),
					},
				}
				if err := space.Insert(xuple); err != nil {
					t.Fatalf("❌ Insert échoué: %v", err)
				}
			}

			// Act
			xuples, err := space.RetrieveMultiple("agent1", requestCount)
			if err != nil {
				t.Fatalf("❌ RetrieveMultiple échoué: %v", err)
			}

			// Assert
			if len(xuples) != requestCount {
				t.Fatalf("❌ Attendu %d xuples, reçu %d", requestCount, len(xuples))
			}

			for i, xuple := range xuples {
				index, ok := xuple.Fact.Fields["index"].(int)
				if !ok {
					t.Errorf("❌ Impossible de récupérer l'index du xuple %d", i)
					continue
				}

				if index != tt.expectedOrder[i] {
					t.Errorf("❌ Xuple[%d]: index = %d, attendu %d", i, index, tt.expectedOrder[i])
				}
			}

			t.Logf("✅ Ordre de sélection correct: %v", tt.expectedOrder)
		})
	}
}

// TestRetrieveMultiple_ConsumptionPolicy teste les différentes politiques de consommation
func TestRetrieveMultiple_ConsumptionPolicy(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: RetrieveMultiple - Politiques de Consommation")
	t.Log("═══════════════════════════════════════════════════════════════")

	const numXuples = 5
	const firstRequest = 3
	const secondRequest = 2

	tests := []struct {
		name                 string
		policy               ConsumptionPolicy
		expectedFirst        int
		expectedSecondAgent1 int
		expectedSecondAgent2 int
	}{
		{
			name:                 "once - consommation unique",
			policy:               NewOnceConsumptionPolicy(),
			expectedFirst:        3,
			expectedSecondAgent1: 2, // Agent1 récupère les 2 restants (pas encore consommés par lui)
			expectedSecondAgent2: 0, // Agent2 ne peut rien récupérer (tous consommés globalement)
		},
		{
			name:                 "per-agent - par agent",
			policy:               NewPerAgentConsumptionPolicy(),
			expectedFirst:        3,
			expectedSecondAgent1: 2, // Agent1 récupère les 2 restants (pas encore consommés par lui)
			expectedSecondAgent2: 3, // Agent2 peut récupérer 3 xuples (pas encore consommés par lui)
		},
		{
			name:                 "limited(2) - limité à 2 consommations",
			policy:               NewLimitedConsumptionPolicy(2),
			expectedFirst:        3,
			expectedSecondAgent1: 2, // Agent1 récupère les 2 restants
			expectedSecondAgent2: 3, // Chaque xuple peut être consommé 2 fois (1 par agent1, 1 par agent2)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := XupleSpaceConfig{
				Name:              "test_consumption",
				SelectionPolicy:   NewFIFOSelectionPolicy(),
				ConsumptionPolicy: tt.policy,
				RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
				MaxSize:           0,
			}
			space := NewXupleSpace(config)

			// Insérer xuples
			for i := 0; i < numXuples; i++ {
				xuple := &Xuple{
					ID: fmt.Sprintf("xuple-%03d", i),
					Fact: &rete.Fact{
						Type: "TestFact",
						Fields: map[string]interface{}{
							"id": fmt.Sprintf("fact-%d", i),
						},
					},
					TriggeringFacts: nil,
					CreatedAt:       time.Now(),
					Metadata: XupleMetadata{
						State:      XupleStateAvailable,
						ConsumedBy: make(map[string]time.Time),
					},
				}
				if err := space.Insert(xuple); err != nil {
					t.Fatalf("❌ Insert échoué: %v", err)
				}
			}

			// Act - Première récupération par agent1
			xuples1, err := space.RetrieveMultiple("agent1", firstRequest)
			if err != nil {
				t.Fatalf("❌ Premier RetrieveMultiple échoué: %v", err)
			}

			if len(xuples1) != tt.expectedFirst {
				t.Errorf("❌ Premier appel: attendu %d, reçu %d", tt.expectedFirst, len(xuples1))
			} else {
				t.Logf("✅ Agent1 première récupération: %d xuples", len(xuples1))
			}

			// Act - Deuxième récupération par agent1
			xuples2, err := space.RetrieveMultiple("agent1", secondRequest)
			if err != nil && tt.expectedSecondAgent1 > 0 {
				t.Errorf("❌ Deuxième RetrieveMultiple agent1 échoué: %v", err)
			}

			if len(xuples2) != tt.expectedSecondAgent1 {
				t.Errorf("❌ Deuxième appel agent1: attendu %d, reçu %d", tt.expectedSecondAgent1, len(xuples2))
			} else {
				t.Logf("✅ Agent1 deuxième récupération: %d xuples", len(xuples2))
			}

			// Act - Récupération par agent2
			xuples3, err := space.RetrieveMultiple("agent2", secondRequest+1) // +1 pour tester la limite
			if err != nil && tt.expectedSecondAgent2 > 0 {
				t.Errorf("❌ RetrieveMultiple agent2 échoué: %v", err)
			}

			if len(xuples3) != tt.expectedSecondAgent2 {
				t.Errorf("❌ Agent2: attendu %d, reçu %d", tt.expectedSecondAgent2, len(xuples3))
			} else {
				t.Logf("✅ Agent2 récupération: %d xuples", len(xuples3))
			}
		})
	}
}

// TestRetrieveMultiple_EmptySpace teste le comportement avec un espace vide
func TestRetrieveMultiple_EmptySpace(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: RetrieveMultiple - Espace Vide")
	t.Log("═══════════════════════════════════════════════════════════════")

	config := XupleSpaceConfig{
		Name:              "test_empty",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}
	space := NewXupleSpace(config)

	// Act
	xuples, err := space.RetrieveMultiple("agent1", 5)

	// Assert
	if err != nil {
		t.Errorf("❌ RetrieveMultiple sur espace vide ne devrait pas retourner d'erreur, reçu: %v", err)
	}

	if len(xuples) != 0 {
		t.Errorf("❌ Attendu 0 xuples, reçu %d", len(xuples))
	} else {
		t.Logf("✅ Espace vide: 0 xuple retourné correctement")
	}
}

// TestRetrieveMultiple_Concurrent teste la récupération concurrente
func TestRetrieveMultiple_Concurrent(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: RetrieveMultiple - Accès Concurrent")
	t.Log("═══════════════════════════════════════════════════════════════")

	const numXuples = 20
	const numAgents = 5
	const requestPerAgent = 5

	config := XupleSpaceConfig{
		Name:              "test_concurrent",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}
	space := NewXupleSpace(config)

	// Insérer xuples
	for i := 0; i < numXuples; i++ {
		xuple := &Xuple{
			ID: fmt.Sprintf("xuple-%03d", i),
			Fact: &rete.Fact{
				Type: "TestFact",
				Fields: map[string]interface{}{
					"id": fmt.Sprintf("fact-%d", i),
				},
			},
			TriggeringFacts: nil,
			CreatedAt:       time.Now(),
			Metadata: XupleMetadata{
				State:      XupleStateAvailable,
				ConsumedBy: make(map[string]time.Time),
			},
		}
		if err := space.Insert(xuple); err != nil {
			t.Fatalf("❌ Insert échoué: %v", err)
		}
	}

	// Lancer plusieurs agents concurrents
	results := make(chan int, numAgents)
	done := make(chan struct{})

	for i := 0; i < numAgents; i++ {
		agentID := fmt.Sprintf("agent-%d", i)
		go func(id string) {
			xuples, _ := space.RetrieveMultiple(id, requestPerAgent)
			results <- len(xuples)
		}(agentID)
	}

	// Collecter résultats
	go func() {
		totalRetrieved := 0
		for i := 0; i < numAgents; i++ {
			count := <-results
			totalRetrieved += count
		}
		if totalRetrieved != numXuples {
			t.Errorf("❌ Total récupéré = %d, attendu %d", totalRetrieved, numXuples)
		} else {
			t.Logf("✅ Total récupéré correctement: %d xuples", totalRetrieved)
		}
		close(done)
	}()

	// Attendre fin
	<-done

	// Vérifier que l'espace est vide
	countAfter := space.Count()
	if countAfter != 0 {
		t.Errorf("❌ Count final = %d, attendu 0", countAfter)
	} else {
		t.Logf("✅ Espace correctement vidé après récupération concurrente")
	}
}

// TestRetrieveMultiple_WithExpiration teste avec xuples expirés
func TestRetrieveMultiple_WithExpiration(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: RetrieveMultiple - Avec Expiration")
	t.Log("═══════════════════════════════════════════════════════════════")

	const shortDuration = 50 * time.Millisecond

	config := XupleSpaceConfig{
		Name:              "test_expiration",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(shortDuration),
		MaxSize:           0,
	}
	space := NewXupleSpace(config)

	// Insérer 5 xuples
	for i := 0; i < 5; i++ {
		xuple := &Xuple{
			ID: fmt.Sprintf("xuple-%03d", i),
			Fact: &rete.Fact{
				Type: "TestFact",
				Fields: map[string]interface{}{
					"id": fmt.Sprintf("fact-%d", i),
				},
			},
			TriggeringFacts: nil,
			CreatedAt:       time.Now(),
			Metadata: XupleMetadata{
				State:      XupleStateAvailable,
				ConsumedBy: make(map[string]time.Time),
			},
		}
		if err := space.Insert(xuple); err != nil {
			t.Fatalf("❌ Insert échoué: %v", err)
		}
	}

	countBefore := space.Count()
	t.Logf("   Count avant expiration: %d", countBefore)

	// Attendre expiration
	time.Sleep(shortDuration + 20*time.Millisecond)

	// Essayer de récupérer
	xuples, err := space.RetrieveMultiple("agent1", 5)
	if err != nil {
		t.Errorf("❌ RetrieveMultiple ne devrait pas échouer avec xuples expirés: %v", err)
	}

	if len(xuples) != 0 {
		t.Errorf("❌ Attendu 0 xuples (tous expirés), reçu %d", len(xuples))
	} else {
		t.Logf("✅ Xuples expirés correctement ignorés")
	}

	countAfter := space.Count()
	if countAfter != 0 {
		t.Logf("   Count après expiration: %d (xuples expirés marqués)", countAfter)
	}
}
