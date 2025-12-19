// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestValidateXupleSpaceDeclaration_Valid(t *testing.T) {
	t.Log("🧪 TEST: Validation - Déclaration xuple-space valide")

	tests := []struct {
		name string
		decl *XupleSpaceDeclaration
	}{
		{
			name: "configuration complète valide",
			decl: &XupleSpaceDeclaration{
				Type:            "xupleSpaceDeclaration",
				Name:            "alerts",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type:  ConsumptionOnce,
					Limit: 0,
				},
				RetentionPolicy: XupleRetentionPolicyConf{
					Type:     RetentionUnlimited,
					Duration: 0,
				},
				MaxSize: 1000,
			},
		},
		{
			name: "avec rétention par durée",
			decl: &XupleSpaceDeclaration{
				Name:            "logs",
				SelectionPolicy: SelectionRandom,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type: ConsumptionPerAgent,
				},
				RetentionPolicy: XupleRetentionPolicyConf{
					Type:     RetentionDuration,
					Duration: 3600, // 1 heure
				},
			},
		},
		{
			name: "consommation limitée",
			decl: &XupleSpaceDeclaration{
				Name:            "tasks",
				SelectionPolicy: SelectionLIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type:  ConsumptionLimited,
					Limit: 5,
				},
				RetentionPolicy: XupleRetentionPolicyConf{
					Type: RetentionUnlimited,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateXupleSpaceDeclaration(tt.decl)
			if err != nil {
				t.Errorf("❌ Validation should pass, got error: %v", err)
			} else {
				t.Log("✅ Validation réussie")
			}
		})
	}
}

func TestValidateXupleSpaceDeclaration_Invalid(t *testing.T) {
	t.Log("🧪 TEST: Validation - Déclarations xuple-space invalides")

	tests := []struct {
		name        string
		decl        *XupleSpaceDeclaration
		expectedErr string
	}{
		{
			name:        "déclaration nulle",
			decl:        nil,
			expectedErr: "cannot be nil",
		},
		{
			name: "nom vide",
			decl: &XupleSpaceDeclaration{
				Name:            "",
				SelectionPolicy: SelectionFIFO,
			},
			expectedErr: "name cannot be empty",
		},
		{
			name: "politique de sélection invalide",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: "invalid",
			},
			expectedErr: "invalid selection policy",
		},
		{
			name: "politique de consommation invalide",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type: "invalid",
				},
			},
			expectedErr: "invalid consumption policy type",
		},
		{
			name: "consommation limitée sans limite",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type:  ConsumptionLimited,
					Limit: 0,
				},
			},
			expectedErr: "consumption limit must be greater than 0",
		},
		{
			name: "consommation limitée avec limite négative",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type:  ConsumptionLimited,
					Limit: -5,
				},
			},
			expectedErr: "consumption limit must be greater than 0",
		},
		{
			name: "politique de rétention invalide",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type: ConsumptionOnce,
				},
				RetentionPolicy: XupleRetentionPolicyConf{
					Type: "invalid",
				},
			},
			expectedErr: "invalid retention policy type",
		},
		{
			name: "rétention par durée sans durée",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type: ConsumptionOnce,
				},
				RetentionPolicy: XupleRetentionPolicyConf{
					Type:     RetentionDuration,
					Duration: 0,
				},
			},
			expectedErr: "retention duration must be greater than 0",
		},
		{
			name: "rétention par durée avec durée négative",
			decl: &XupleSpaceDeclaration{
				Name:            "bad",
				SelectionPolicy: SelectionFIFO,
				ConsumptionPolicy: XupleConsumptionPolicyConf{
					Type: ConsumptionOnce,
				},
				RetentionPolicy: XupleRetentionPolicyConf{
					Type:     RetentionDuration,
					Duration: -100,
				},
			},
			expectedErr: "retention duration must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateXupleSpaceDeclaration(tt.decl)
			if err == nil {
				t.Errorf("❌ Expected validation error containing '%s'", tt.expectedErr)
			} else if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("❌ Expected error containing '%s', got: %v", tt.expectedErr, err)
			} else {
				t.Logf("✅ Erreur attendue: %v", err)
			}
		})
	}
}

func TestValidateXupleSpaceProperties_Valid(t *testing.T) {
	t.Log("🧪 TEST: Validation - Propriétés xuple-space valides")

	tests := []struct {
		name  string
		props map[string]interface{}
	}{
		{
			name: "toutes propriétés valides",
			props: map[string]interface{}{
				"selectionPolicy": "fifo",
				"consumptionPolicy": map[string]interface{}{
					"type":  "once",
					"limit": 0,
				},
				"retentionPolicy": map[string]interface{}{
					"type":     "unlimited",
					"duration": 0,
				},
				"maxSize": 1000,
			},
		},
		{
			name: "max-size à zéro (illimité)",
			props: map[string]interface{}{
				"maxSize": 0,
			},
		},
		{
			name: "propriétés minimales",
			props: map[string]interface{}{
				"selectionPolicy": "random",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateXupleSpaceProperties("test-space", tt.props)
			if err != nil {
				t.Errorf("❌ Validation should pass, got error: %v", err)
			} else {
				t.Log("✅ Validation réussie")
			}
		})
	}
}

func TestValidateXupleSpaceProperties_Invalid(t *testing.T) {
	t.Log("🧪 TEST: Validation - Propriétés xuple-space invalides")

	tests := []struct {
		name        string
		spaceName   string
		props       map[string]interface{}
		expectedErr string
	}{
		{
			name:        "nom vide",
			spaceName:   "",
			props:       map[string]interface{}{},
			expectedErr: "name cannot be empty",
		},
		{
			name:      "selection invalide (type incorrect)",
			spaceName: "test",
			props: map[string]interface{}{
				"selectionPolicy": 123, // devrait être string
			},
			expectedErr: "selectionPolicy must be a string",
		},
		{
			name:      "consumption invalide (type incorrect)",
			spaceName: "test",
			props: map[string]interface{}{
				"consumptionPolicy": "invalid", // devrait être map
			},
			expectedErr: "consumptionPolicy must be a map",
		},
		{
			name:      "retention invalide (type incorrect)",
			spaceName: "test",
			props: map[string]interface{}{
				"retentionPolicy": "invalid", // devrait être map
			},
			expectedErr: "retentionPolicy must be a map",
		},
		{
			name:      "max-size négatif",
			spaceName: "test",
			props: map[string]interface{}{
				"maxSize": -100,
			},
			expectedErr: "max-size must be >= 0",
		},
		{
			name:      "max-size type invalide",
			spaceName: "test",
			props: map[string]interface{}{
				"maxSize": "invalid", // devrait être int
			},
			expectedErr: "max-size must be an integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateXupleSpaceProperties(tt.spaceName, tt.props)
			if err == nil {
				t.Errorf("❌ Expected validation error containing '%s'", tt.expectedErr)
			} else if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("❌ Expected error containing '%s', got: %v", tt.expectedErr, err)
			} else {
				t.Logf("✅ Erreur attendue: %v", err)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	t.Log("🧪 TEST: Parsing de durées")

	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "durée vide",
			input:    "",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "format Go standard - 24h",
			input:    "24h",
			expected: 24 * SecondsPerHour,
			wantErr:  false,
		},
		{
			name:     "format Go standard - 5m",
			input:    "5m",
			expected: 5 * SecondsPerMinute,
			wantErr:  false,
		},
		{
			name:     "format Go standard - 30s",
			input:    "30s",
			expected: 30,
			wantErr:  false,
		},
		{
			name:     "format personnalisé - 7 jours",
			input:    "7d",
			expected: 7 * SecondsPerDay,
			wantErr:  false,
		},
		{
			name:     "format personnalisé - 2 semaines",
			input:    "2w",
			expected: 2 * SecondsPerWeek,
			wantErr:  false,
		},
		{
			name:    "format invalide",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "durée négative",
			input:   "-5d",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("❌ ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("❌ ParseDuration() = %d, want %d", result, tt.expected)
			} else if !tt.wantErr {
				t.Logf("✅ ParseDuration('%s') = %d seconds", tt.input, result)
			} else {
				t.Logf("✅ Erreur attendue pour '%s': %v", tt.input, err)
			}
		})
	}
}
