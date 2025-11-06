package rete

import (
	"fmt"
	"testing"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
	"github.com/treivax/tsd/rete/pkg/nodes"
)

// MockLogger implémente l'interface Logger pour les tests
type MockAdvancedLogger struct{}

func (ml *MockAdvancedLogger) Debug(msg string, fields map[string]interface{})            {}
func (ml *MockAdvancedLogger) Info(msg string, fields map[string]interface{})             {}
func (ml *MockAdvancedLogger) Warn(msg string, fields map[string]interface{})             {}
func (ml *MockAdvancedLogger) Error(msg string, err error, fields map[string]interface{}) {}

// MockStorage implémente l'interface Storage pour les tests
type MockStorage struct{}

func (m *MockStorage) SaveMemory(nodeID string, memory *WorkingMemory) error {
	return nil
}

func (m *MockStorage) LoadMemory(nodeID string) (*WorkingMemory, error) {
	return &WorkingMemory{Facts: make(map[string]*Fact)}, nil
}

func (m *MockStorage) DeleteMemory(nodeID string) error {
	return nil
}

func (m *MockStorage) ListNodes() ([]string, error) {
	return []string{}, nil
}

// Test intégré des nœuds avancés : Détection de fraude bancaire sophistiquée
func TestAdvancedNodesIntegration_ComplexFraudDetection(t *testing.T) {
	logger := &MockAdvancedLogger{}

	// 1. Créer un nœud NOT pour "Pas de transaction légitime récente"
	notNode := nodes.NewNotNode("fraud_not_1", logger)
	notNode.SetNegationCondition("type == 'legitimate' AND timestamp > recent")

	// 2. Créer un nœud EXISTS pour "Il existe des transactions suspectes"
	existsNode := nodes.NewExistsNode("fraud_exists_1", logger)
	variable := domain.TypedVariable{
		Name:     "suspicious_tx",
		DataType: "Transaction",
	}
	existsNode.SetExistenceCondition(variable, "amount > 10000 AND location != home")

	// 3. Créer un nœud d'accumulation pour "Somme des transactions dans la journée"
	accumulator := domain.AccumulateFunction{
		FunctionType: "SUM",
		Field:        "amount",
	}
	accNode := nodes.NewAccumulateNode("fraud_sum_1", accumulator, logger)

	// === Scénario de test : Détection de fraude complexe ===

	// Ajouter des transactions légitimes (pour le NOT node)
	legitimateTx1 := &domain.Fact{
		ID:   "tx_legit_1",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"type":      "legitimate",
			"amount":    50.0,
			"location":  "home",
			"timestamp": time.Now().Add(-1 * time.Hour), // Ancienne
		},
		Timestamp: time.Now(),
	}

	legitimateTx2 := &domain.Fact{
		ID:   "tx_legit_2",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"type":      "legitimate",
			"amount":    100.0,
			"location":  "home",
			"timestamp": time.Now().Add(-30 * time.Minute), // Récente
		},
		Timestamp: time.Now(),
	}

	// Ajouter des transactions suspectes (pour le EXISTS node)
	suspiciousTx1 := &domain.Fact{
		ID:   "tx_suspicious_1",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"type":     "withdrawal",
			"amount":   15000.0, // Gros montant
			"location": "foreign_country",
			"risk":     "high",
		},
		Timestamp: time.Now(),
	}

	suspiciousTx2 := &domain.Fact{
		ID:   "tx_suspicious_2",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"type":     "transfer",
			"amount":   25000.0, // Très gros montant
			"location": "unknown",
			"risk":     "high",
		},
		Timestamp: time.Now(),
	}

	// Token représentant un compte à analyser
	accountToken := &domain.Token{
		ID: "account_analysis_1",
		Facts: []*domain.Fact{
			{
				ID:   "account_1",
				Type: "Account",
				Fields: map[string]interface{}{
					"id":      "ACC123456",
					"balance": 50000.0,
					"status":  "active",
				},
				Timestamp: time.Now(),
			},
		},
	}

	t.Run("NotNode_NoRecentLegitimateTransactions", func(t *testing.T) {
		// Test du nœud NOT : Vérifier qu'il n'y a pas de transactions légitimes récentes

		// Ajouter seulement la transaction ancienne
		err := notNode.ProcessRightFact(legitimateTx1)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		// Traiter le token (devrait passer car la transaction ancienne ne satisfait pas "recent")
		err = notNode.ProcessLeftToken(accountToken)
		if err != nil {
			t.Errorf("Failed to process left token: %v", err)
		}

		// Vérifier que le token est propagé (négation réussie)
		tokens := notNode.GetLeftMemory()
		if len(tokens) != 1 {
			t.Errorf("Expected 1 token in NOT node memory, got %d", len(tokens))
		}

		fmt.Println("✓ NOT Node: Pas de transaction légitime récente détectée")
	})

	t.Run("ExistsNode_SuspiciousTransactionsExist", func(t *testing.T) {
		// Test du nœud EXISTS : Vérifier l'existence de transactions suspectes

		// Ajouter les transactions suspectes
		err := existsNode.ProcessRightFact(suspiciousTx1)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		err = existsNode.ProcessRightFact(suspiciousTx2)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		// Traiter le token
		err = existsNode.ProcessLeftToken(accountToken)
		if err != nil {
			t.Errorf("Failed to process left token: %v", err)
		}

		// Vérifier l'existence (devrait détecter les transactions suspectes)
		exists := existsNode.CheckExistence(accountToken)
		if !exists {
			t.Errorf("Expected existence check to return true for suspicious transactions")
		}

		fmt.Println("✓ EXISTS Node: Transactions suspectes détectées")
	})

	t.Run("AccumulateNode_DailyTransactionSum", func(t *testing.T) {
		// Test du nœud d'accumulation : Calculer la somme des transactions

		// Ajouter toutes les transactions pour le calcul
		err := accNode.ProcessRightFact(legitimateTx1)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		err = accNode.ProcessRightFact(legitimateTx2)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		err = accNode.ProcessRightFact(suspiciousTx1)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		err = accNode.ProcessRightFact(suspiciousTx2)
		if err != nil {
			t.Errorf("Failed to process right fact: %v", err)
		}

		// Traiter le token pour déclencher l'accumulation
		err = accNode.ProcessLeftToken(accountToken)
		if err != nil {
			t.Errorf("Failed to process left token: %v", err)
		}

		// Calculer manuellement la somme attendue
		expectedSum := 50.0 + 100.0 + 15000.0 + 25000.0 // 40150.0

		// Vérifier la somme calculée
		facts := []*domain.Fact{legitimateTx1, legitimateTx2, suspiciousTx1, suspiciousTx2}
		result, err := accNode.ComputeAggregate(accountToken, facts)
		if err != nil {
			t.Errorf("Failed to compute aggregate: %v", err)
		}

		if result != expectedSum {
			t.Errorf("Expected sum %v, got %v", expectedSum, result)
		}

		fmt.Printf("✓ ACCUMULATE Node: Somme des transactions = %.2f\n", result.(float64))
	})

	t.Run("IntegratedFraudDetection", func(t *testing.T) {
		// Test intégré : Tous les nœuds ensemble pour détecter la fraude

		fmt.Println("\n=== DÉTECTION DE FRAUDE INTÉGRÉE ===")

		// Conditions de fraude détectées :
		// 1. PAS de transactions légitimes récentes (NOT node)
		// 2. Présence de transactions suspectes (EXISTS node)
		// 3. Somme élevée des transactions (ACCUMULATE node > 10000)

		fraudScore := 0
		reasons := []string{}

		// Vérifier la négation (pas de transactions légitimes récentes)
		notResult := notNode.ProcessNegation(accountToken, legitimateTx1) // Transaction ancienne
		if notResult {
			fraudScore += 30
			reasons = append(reasons, "Absence de transactions légitimes récentes")
		}

		// Vérifier l'existence (transactions suspectes présentes)
		existsResult := existsNode.CheckExistence(accountToken)
		if existsResult {
			fraudScore += 50
			reasons = append(reasons, "Présence de transactions suspectes")
		}

		// Vérifier l'accumulation (somme élevée)
		facts := []*domain.Fact{legitimateTx1, legitimateTx2, suspiciousTx1, suspiciousTx2}
		sumResult, err := accNode.ComputeAggregate(accountToken, facts)
		if err == nil {
			if sumResult.(float64) > 10000.0 {
				fraudScore += 20
				reasons = append(reasons, fmt.Sprintf("Somme élevée des transactions: %.2f", sumResult.(float64)))
			}
		}

		// Évaluation finale
		if fraudScore >= 70 {
			fmt.Printf("🚨 FRAUDE DÉTECTÉE - Score: %d/100\n", fraudScore)
			for _, reason := range reasons {
				fmt.Printf("   • %s\n", reason)
			}
		} else {
			fmt.Printf("✅ Compte sécurisé - Score: %d/100\n", fraudScore)
		}

		// Le test passe si on détecte bien la fraude
		if fraudScore < 70 {
			t.Errorf("Expected fraud to be detected with score >= 70, got %d", fraudScore)
		}
	})
}

// Test des différentes fonctions d'agrégation
func TestAccumulateNodeAggregationFunctions(t *testing.T) {
	logger := &MockAdvancedLogger{}

	// Créer des faits de test avec différents montants
	transactions := []*domain.Fact{
		{
			ID:   "tx1",
			Type: "Transaction",
			Fields: map[string]interface{}{
				"amount": 100,
				"type":   "purchase",
			},
			Timestamp: time.Now(),
		},
		{
			ID:   "tx2",
			Type: "Transaction",
			Fields: map[string]interface{}{
				"amount": 250.5,
				"type":   "transfer",
			},
			Timestamp: time.Now(),
		},
		{
			ID:   "tx3",
			Type: "Transaction",
			Fields: map[string]interface{}{
				"amount": 75.25,
				"type":   "withdrawal",
			},
			Timestamp: time.Now(),
		},
		{
			ID:   "tx4",
			Type: "Transaction",
			Fields: map[string]interface{}{
				"amount": 500.0,
				"type":   "deposit",
			},
			Timestamp: time.Now(),
		},
	}

	token := &domain.Token{
		ID:    "test_token",
		Facts: []*domain.Fact{transactions[0]},
	}

	t.Run("SUM_Aggregation", func(t *testing.T) {
		accumulator := domain.AccumulateFunction{
			FunctionType: "SUM",
			Field:        "amount",
		}
		accNode := nodes.NewAccumulateNode("sum_node", accumulator, logger)

		result, err := accNode.ComputeAggregate(token, transactions)
		if err != nil {
			t.Errorf("SUM aggregation failed: %v", err)
		}

		expectedSum := 925.75 // 100 + 250.5 + 75.25 + 500.0
		if result != expectedSum {
			t.Errorf("Expected SUM %v, got %v", expectedSum, result)
		}

		fmt.Printf("SUM: %.2f\n", result.(float64))
	})

	t.Run("COUNT_Aggregation", func(t *testing.T) {
		accumulator := domain.AccumulateFunction{
			FunctionType: "COUNT",
			Field:        "",
		}
		accNode := nodes.NewAccumulateNode("count_node", accumulator, logger)

		result, err := accNode.ComputeAggregate(token, transactions)
		if err != nil {
			t.Errorf("COUNT aggregation failed: %v", err)
		}

		expectedCount := 4
		if result != expectedCount {
			t.Errorf("Expected COUNT %v, got %v", expectedCount, result)
		}

		fmt.Printf("COUNT: %d\n", result.(int))
	})

	t.Run("AVG_Aggregation", func(t *testing.T) {
		accumulator := domain.AccumulateFunction{
			FunctionType: "AVG",
			Field:        "amount",
		}
		accNode := nodes.NewAccumulateNode("avg_node", accumulator, logger)

		result, err := accNode.ComputeAggregate(token, transactions)
		if err != nil {
			t.Errorf("AVG aggregation failed: %v", err)
		}

		expectedAvg := 231.4375 // 925.75 / 4
		if result != expectedAvg {
			t.Errorf("Expected AVG %v, got %v", expectedAvg, result)
		}

		fmt.Printf("AVG: %.4f\n", result.(float64))
	})

	t.Run("MIN_Aggregation", func(t *testing.T) {
		accumulator := domain.AccumulateFunction{
			FunctionType: "MIN",
			Field:        "amount",
		}
		accNode := nodes.NewAccumulateNode("min_node", accumulator, logger)

		result, err := accNode.ComputeAggregate(token, transactions)
		if err != nil {
			t.Errorf("MIN aggregation failed: %v", err)
		}

		expectedMin := 75.25
		if result != expectedMin {
			t.Errorf("Expected MIN %v, got %v", expectedMin, result)
		}

		fmt.Printf("MIN: %.2f\n", result.(float64))
	})

	t.Run("MAX_Aggregation", func(t *testing.T) {
		accumulator := domain.AccumulateFunction{
			FunctionType: "MAX",
			Field:        "amount",
		}
		accNode := nodes.NewAccumulateNode("max_node", accumulator, logger)

		result, err := accNode.ComputeAggregate(token, transactions)
		if err != nil {
			t.Errorf("MAX aggregation failed: %v", err)
		}

		expectedMax := 500.0
		if result != expectedMax {
			t.Errorf("Expected MAX %v, got %v", expectedMax, result)
		}

		fmt.Printf("MAX: %.2f\n", result.(float64))
	})
}

// Test d'intégration réseau avec nœuds avancés
func TestReteNetwork_AdvancedNodesIntegration(t *testing.T) {
	// Créer un réseau RETE avec stockage mock
	storage := &MockStorage{}
	network := NewReteNetwork(storage)

	t.Run("CreateAdvancedNodes", func(t *testing.T) {
		// Créer des nœuds avancés dans le réseau
		err := network.CreateNotNode("not_fraud", "amount < 100")
		if err != nil {
			t.Errorf("Failed to create NOT node: %v", err)
		}

		err = network.CreateExistsNode("exists_suspicious", "suspicious_tx", "Transaction", "risk == 'high'")
		if err != nil {
			t.Errorf("Failed to create EXISTS node: %v", err)
		}

		err = network.CreateAccumulateNode("sum_daily", "SUM", "amount", nil)
		if err != nil {
			t.Errorf("Failed to create ACCUMULATE node: %v", err)
		}

		// Vérifier que les nœuds sont créés
		if len(network.BetaNodes) != 3 {
			t.Errorf("Expected 3 beta nodes, got %d", len(network.BetaNodes))
		}
	})

	t.Run("AdvancedNodeStatistics", func(t *testing.T) {
		stats := network.GetAdvancedNodeStatistics()

		expectedStats := map[string]int{
			"notNodes":        1,
			"existsNodes":     1,
			"accumulateNodes": 1,
		}

		for key, expected := range expectedStats {
			if stats[key] != expected {
				t.Errorf("Expected %s = %d, got %v", key, expected, stats[key])
			}
		}

		if stats["advancedEnabled"] != true {
			t.Errorf("Expected advanced nodes to be enabled")
		}

		fmt.Printf("Advanced Node Statistics: %+v\n", stats)
	})
}
