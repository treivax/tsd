// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestGetSuccessRate(t *testing.T) {
	t.Log("🧪 TEST getSuccessRate")
	t.Log("======================")

	tests := []struct {
		name                   string
		transactionCount       int
		successfulTransactions int
		expectedRate           float64
	}{
		{
			name:                   "100% de succès",
			transactionCount:       100,
			successfulTransactions: 100,
			expectedRate:           100.0,
		},
		{
			name:                   "75% de succès",
			transactionCount:       100,
			successfulTransactions: 75,
			expectedRate:           75.0,
		},
		{
			name:                   "50% de succès",
			transactionCount:       100,
			successfulTransactions: 50,
			expectedRate:           50.0,
		},
		{
			name:                   "0% de succès",
			transactionCount:       100,
			successfulTransactions: 0,
			expectedRate:           0.0,
		},
		{
			name:                   "cas limite - 0 transactions",
			transactionCount:       0,
			successfulTransactions: 0,
			expectedRate:           0.0,
		},
		{
			name:                   "petit échantillon",
			transactionCount:       3,
			successfulTransactions: 2,
			expectedRate:           66.66666666666666,
		},
		{
			name:                   "une seule transaction réussie",
			transactionCount:       1,
			successfulTransactions: 1,
			expectedRate:           100.0,
		},
		{
			name:                   "une seule transaction échouée",
			transactionCount:       1,
			successfulTransactions: 0,
			expectedRate:           0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				TransactionCount:       tt.transactionCount,
				SuccessfulTransactions: tt.successfulTransactions,
			}

			result := pm.getSuccessRate()

			if result != tt.expectedRate {
				t.Errorf("❌ getSuccessRate() = %.2f, attendu %.2f", result, tt.expectedRate)
				return
			}

			t.Logf("✅ Test réussi: %d/%d transactions = %.2f%%",
				tt.successfulTransactions, tt.transactionCount, result)
		})
	}
}

func TestGetFailureRate(t *testing.T) {
	t.Log("🧪 TEST getFailureRate")
	t.Log("======================")

	tests := []struct {
		name               string
		transactionCount   int
		failedTransactions int
		expectedRate       float64
	}{
		{
			name:               "100% d'échecs",
			transactionCount:   100,
			failedTransactions: 100,
			expectedRate:       100.0,
		},
		{
			name:               "25% d'échecs",
			transactionCount:   100,
			failedTransactions: 25,
			expectedRate:       25.0,
		},
		{
			name:               "50% d'échecs",
			transactionCount:   100,
			failedTransactions: 50,
			expectedRate:       50.0,
		},
		{
			name:               "0% d'échecs",
			transactionCount:   100,
			failedTransactions: 0,
			expectedRate:       0.0,
		},
		{
			name:               "cas limite - 0 transactions",
			transactionCount:   0,
			failedTransactions: 0,
			expectedRate:       0.0,
		},
		{
			name:               "petit échantillon",
			transactionCount:   4,
			failedTransactions: 1,
			expectedRate:       25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				TransactionCount:   tt.transactionCount,
				FailedTransactions: tt.failedTransactions,
			}

			result := pm.getFailureRate()

			if result != tt.expectedRate {
				t.Errorf("❌ getFailureRate() = %.2f, attendu %.2f", result, tt.expectedRate)
				return
			}

			t.Logf("✅ Test réussi: %d/%d échecs = %.2f%%",
				tt.failedTransactions, tt.transactionCount, result)
		})
	}
}

func TestGetFactPersistRate(t *testing.T) {
	t.Log("🧪 TEST getFactPersistRate")
	t.Log("==========================")

	tests := []struct {
		name                string
		totalFactsProcessed int
		totalFactsPersisted int
		expectedRate        float64
	}{
		{
			name:                "100% de persistance",
			totalFactsProcessed: 1000,
			totalFactsPersisted: 1000,
			expectedRate:        100.0,
		},
		{
			name:                "90% de persistance",
			totalFactsProcessed: 100,
			totalFactsPersisted: 90,
			expectedRate:        90.0,
		},
		{
			name:                "50% de persistance",
			totalFactsProcessed: 200,
			totalFactsPersisted: 100,
			expectedRate:        50.0,
		},
		{
			name:                "0% de persistance",
			totalFactsProcessed: 100,
			totalFactsPersisted: 0,
			expectedRate:        0.0,
		},
		{
			name:                "cas limite - 0 faits traités",
			totalFactsProcessed: 0,
			totalFactsPersisted: 0,
			expectedRate:        0.0,
		},
		{
			name:                "un seul fait persisté",
			totalFactsProcessed: 1,
			totalFactsPersisted: 1,
			expectedRate:        100.0,
		},
		{
			name:                "un fait sur trois",
			totalFactsProcessed: 3,
			totalFactsPersisted: 1,
			expectedRate:        33.33333333333333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				TotalFactsProcessed: tt.totalFactsProcessed,
				TotalFactsPersisted: tt.totalFactsPersisted,
			}

			result := pm.getFactPersistRate()

			if result != tt.expectedRate {
				t.Errorf("❌ getFactPersistRate() = %.2f, attendu %.2f", result, tt.expectedRate)
				return
			}

			t.Logf("✅ Test réussi: %d/%d faits persistés = %.2f%%",
				tt.totalFactsPersisted, tt.totalFactsProcessed, result)
		})
	}
}

func TestGetFactFailureRate(t *testing.T) {
	t.Log("🧪 TEST getFactFailureRate")
	t.Log("==========================")

	tests := []struct {
		name                string
		totalFactsProcessed int
		totalFactsFailed    int
		expectedRate        float64
	}{
		{
			name:                "100% d'échecs",
			totalFactsProcessed: 100,
			totalFactsFailed:    100,
			expectedRate:        100.0,
		},
		{
			name:                "10% d'échecs",
			totalFactsProcessed: 100,
			totalFactsFailed:    10,
			expectedRate:        10.0,
		},
		{
			name:                "50% d'échecs",
			totalFactsProcessed: 200,
			totalFactsFailed:    100,
			expectedRate:        50.0,
		},
		{
			name:                "0% d'échecs",
			totalFactsProcessed: 100,
			totalFactsFailed:    0,
			expectedRate:        0.0,
		},
		{
			name:                "cas limite - 0 faits traités",
			totalFactsProcessed: 0,
			totalFactsFailed:    0,
			expectedRate:        0.0,
		},
		{
			name:                "un fait échoué sur cinq",
			totalFactsProcessed: 5,
			totalFactsFailed:    1,
			expectedRate:        20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				TotalFactsProcessed: tt.totalFactsProcessed,
				TotalFactsFailed:    tt.totalFactsFailed,
			}

			result := pm.getFactFailureRate()

			if result != tt.expectedRate {
				t.Errorf("❌ getFactFailureRate() = %.2f, attendu %.2f", result, tt.expectedRate)
				return
			}

			t.Logf("✅ Test réussi: %d/%d faits échoués = %.2f%%",
				tt.totalFactsFailed, tt.totalFactsProcessed, result)
		})
	}
}

func TestGetVerifySuccessRate(t *testing.T) {
	t.Log("🧪 TEST getVerifySuccessRate")
	t.Log("============================")

	tests := []struct {
		name               string
		totalVerifications int
		successfulVerifies int
		expectedRate       float64
	}{
		{
			name:               "100% de vérifications réussies",
			totalVerifications: 1000,
			successfulVerifies: 1000,
			expectedRate:       100.0,
		},
		{
			name:               "95% de vérifications réussies",
			totalVerifications: 100,
			successfulVerifies: 95,
			expectedRate:       95.0,
		},
		{
			name:               "50% de vérifications réussies",
			totalVerifications: 200,
			successfulVerifies: 100,
			expectedRate:       50.0,
		},
		{
			name:               "0% de vérifications réussies",
			totalVerifications: 100,
			successfulVerifies: 0,
			expectedRate:       0.0,
		},
		{
			name:               "cas limite - 0 vérifications",
			totalVerifications: 0,
			successfulVerifies: 0,
			expectedRate:       0.0,
		},
		{
			name:               "deux vérifications sur trois",
			totalVerifications: 3,
			successfulVerifies: 2,
			expectedRate:       66.66666666666666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				TotalVerifications: tt.totalVerifications,
				SuccessfulVerifies: tt.successfulVerifies,
			}

			result := pm.getVerifySuccessRate()

			if result != tt.expectedRate {
				t.Errorf("❌ getVerifySuccessRate() = %.2f, attendu %.2f", result, tt.expectedRate)
				return
			}

			t.Logf("✅ Test réussi: %d/%d vérifications réussies = %.2f%%",
				tt.successfulVerifies, tt.totalVerifications, result)
		})
	}
}

func TestGetCommitSuccessRate(t *testing.T) {
	t.Log("🧪 TEST getCommitSuccessRate")
	t.Log("============================")

	tests := []struct {
		name              string
		totalCommits      int
		successfulCommits int
		expectedRate      float64
	}{
		{
			name:              "100% de commits réussis",
			totalCommits:      500,
			successfulCommits: 500,
			expectedRate:      100.0,
		},
		{
			name:              "99% de commits réussis",
			totalCommits:      100,
			successfulCommits: 99,
			expectedRate:      99.0,
		},
		{
			name:              "50% de commits réussis",
			totalCommits:      200,
			successfulCommits: 100,
			expectedRate:      50.0,
		},
		{
			name:              "0% de commits réussis",
			totalCommits:      100,
			successfulCommits: 0,
			expectedRate:      0.0,
		},
		{
			name:              "cas limite - 0 commits",
			totalCommits:      0,
			successfulCommits: 0,
			expectedRate:      0.0,
		},
		{
			name:              "premier commit réussi",
			totalCommits:      1,
			successfulCommits: 1,
			expectedRate:      100.0,
		},
		{
			name:              "trois commits sur quatre",
			totalCommits:      4,
			successfulCommits: 3,
			expectedRate:      75.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				TotalCommits:      tt.totalCommits,
				SuccessfulCommits: tt.successfulCommits,
			}

			result := pm.getCommitSuccessRate()

			if result != tt.expectedRate {
				t.Errorf("❌ getCommitSuccessRate() = %.2f, attendu %.2f", result, tt.expectedRate)
				return
			}

			t.Logf("✅ Test réussi: %d/%d commits réussis = %.2f%%",
				tt.successfulCommits, tt.totalCommits, result)
		})
	}
}

func TestPerformanceCalculationsConsistency(t *testing.T) {
	t.Log("🧪 TEST Cohérence des Calculs de Performance")
	t.Log("============================================")

	t.Run("succès + échecs = 100%", func(t *testing.T) {
		pm := &StrongModePerformanceMetrics{
			TransactionCount:       100,
			SuccessfulTransactions: 75,
			FailedTransactions:     25,
		}

		successRate := pm.getSuccessRate()
		failureRate := pm.getFailureRate()

		total := successRate + failureRate

		if total != 100.0 {
			t.Errorf("❌ Incohérence: succès(%.2f%%) + échecs(%.2f%%) = %.2f%%, attendu 100%%",
				successRate, failureRate, total)
			return
		}

		t.Logf("✅ Cohérence vérifiée: %.2f%% + %.2f%% = 100%%", successRate, failureRate)
	})

	t.Run("persistance + échecs faits = 100%", func(t *testing.T) {
		pm := &StrongModePerformanceMetrics{
			TotalFactsProcessed: 200,
			TotalFactsPersisted: 180,
			TotalFactsFailed:    20,
		}

		persistRate := pm.getFactPersistRate()
		failureRate := pm.getFactFailureRate()

		total := persistRate + failureRate

		if total != 100.0 {
			t.Errorf("❌ Incohérence: persistance(%.2f%%) + échecs(%.2f%%) = %.2f%%, attendu 100%%",
				persistRate, failureRate, total)
			return
		}

		t.Logf("✅ Cohérence vérifiée: %.2f%% + %.2f%% = 100%%", persistRate, failureRate)
	})

	t.Run("métriques vides retournent 0%", func(t *testing.T) {
		pm := &StrongModePerformanceMetrics{}

		rates := []struct {
			name  string
			value float64
		}{
			{"getSuccessRate", pm.getSuccessRate()},
			{"getFailureRate", pm.getFailureRate()},
			{"getFactPersistRate", pm.getFactPersistRate()},
			{"getFactFailureRate", pm.getFactFailureRate()},
			{"getVerifySuccessRate", pm.getVerifySuccessRate()},
			{"getCommitSuccessRate", pm.getCommitSuccessRate()},
		}

		for _, r := range rates {
			if r.value != 0.0 {
				t.Errorf("❌ %s devrait retourner 0%% pour métriques vides, reçu %.2f%%",
					r.name, r.value)
			}
		}

		t.Log("✅ Toutes les fonctions retournent 0% pour métriques vides")
	})
}

func TestGetHealthStatus(t *testing.T) {
	t.Log("🧪 TEST getHealthStatus")
	t.Log("=======================")

	tests := []struct {
		name           string
		isHealthy      bool
		expectedStatus string
	}{
		{
			name:           "système sain",
			isHealthy:      true,
			expectedStatus: "✅ Healthy",
		},
		{
			name:           "système nécessitant attention",
			isHealthy:      false,
			expectedStatus: "⚠️  Needs Attention",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &StrongModePerformanceMetrics{
				IsHealthy: tt.isHealthy,
			}

			result := pm.getHealthStatus()

			if result != tt.expectedStatus {
				t.Errorf("❌ getHealthStatus() = %q, attendu %q", result, tt.expectedStatus)
				return
			}

			t.Logf("✅ Test réussi: IsHealthy=%v → %q", tt.isHealthy, result)
		})
	}
}

func TestPerformanceCalculationsEdgeCases(t *testing.T) {
	t.Log("🧪 TEST Cas Limites des Calculs de Performance")
	t.Log("==============================================")

	t.Run("division par zéro prévenue", func(t *testing.T) {
		pm := &StrongModePerformanceMetrics{
			TransactionCount:    0,
			TotalFactsProcessed: 0,
			TotalVerifications:  0,
			TotalCommits:        0,
		}

		// Toutes ces fonctions devraient gérer la division par zéro
		successRate := pm.getSuccessRate()
		failureRate := pm.getFailureRate()
		factPersistRate := pm.getFactPersistRate()
		factFailureRate := pm.getFactFailureRate()
		verifySuccessRate := pm.getVerifySuccessRate()
		commitSuccessRate := pm.getCommitSuccessRate()

		if successRate != 0.0 || failureRate != 0.0 || factPersistRate != 0.0 ||
			factFailureRate != 0.0 || verifySuccessRate != 0.0 || commitSuccessRate != 0.0 {
			t.Error("❌ Division par zéro non gérée correctement")
			return
		}

		t.Log("✅ Division par zéro correctement gérée")
	})

	t.Run("très grands nombres", func(t *testing.T) {
		pm := &StrongModePerformanceMetrics{
			TransactionCount:       1000000,
			SuccessfulTransactions: 999999,
		}

		successRate := pm.getSuccessRate()
		expectedRate := 99.9999

		if successRate < 99.99 || successRate > 100.0 {
			t.Errorf("❌ Calcul incorrect pour grands nombres: %.4f%%", successRate)
			return
		}

		t.Logf("✅ Grands nombres gérés correctement: %.4f%% ≈ %.4f%%", successRate, expectedRate)
	})
}
