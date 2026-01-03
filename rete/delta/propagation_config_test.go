// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestPropagationMode_String(t *testing.T) {
	tests := []struct {
		mode PropagationMode
		want string
	}{
		{PropagationModeDelta, "Delta"},
		{PropagationModeClassic, "Classic"},
		{PropagationModeAuto, "Auto"},
		{PropagationMode(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("PropagationMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultPropagationConfig(t *testing.T) {
	config := DefaultPropagationConfig()

	if config.DefaultMode != PropagationModeAuto {
		t.Errorf("Expected Auto mode, got %v", config.DefaultMode)
	}

	if !config.EnableDeltaPropagation {
		t.Error("Expected EnableDeltaPropagation = true")
	}

	if config.DeltaThreshold != 0.5 {
		t.Errorf("Expected DeltaThreshold = 0.5, got %v", config.DeltaThreshold)
	}

	if config.MinFieldsForDelta != 3 {
		t.Errorf("Expected MinFieldsForDelta = 3, got %d", config.MinFieldsForDelta)
	}

	if config.MaxAffectedNodesForDelta != 100 {
		t.Errorf("Expected MaxAffectedNodesForDelta = 100, got %d", config.MaxAffectedNodesForDelta)
	}

	if config.PropagationTimeout != 30*time.Second {
		t.Errorf("Expected timeout = 30s, got %v", config.PropagationTimeout)
	}
}

func TestPropagationConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    PropagationConfig
		wantError bool
	}{
		{
			name:      "valid default",
			config:    DefaultPropagationConfig(),
			wantError: false,
		},
		{
			name: "invalid delta threshold (negative)",
			config: PropagationConfig{
				DeltaThreshold:            -0.1,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid delta threshold (> 1)",
			config: PropagationConfig{
				DeltaThreshold:            1.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid min fields (negative)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         -1,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid max nodes (zero)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  0,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid timeout (negative)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        -time.Second,
				MaxConcurrentPropagations: 10,
			},
			wantError: true,
		},
		{
			name: "invalid max concurrent (zero)",
			config: PropagationConfig{
				DeltaThreshold:            0.5,
				MinFieldsForDelta:         3,
				MaxAffectedNodesForDelta:  100,
				PropagationTimeout:        time.Second,
				MaxConcurrentPropagations: 0,
			},
			wantError: true,
		},
		{
			name: "edge case: zero threshold",
			config: PropagationConfig{
				DeltaThreshold:            0.0,
				MinFieldsForDelta:         0,
				MaxAffectedNodesForDelta:  1,
				PropagationTimeout:        0,
				MaxConcurrentPropagations: 1,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestPropagationConfig_ShouldUseDelta_FeatureFlagDisabled(t *testing.T) {
	config := DefaultPropagationConfig()
	config.EnableDeltaPropagation = false

	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("field1", "old", "new")

	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when feature flag disabled")
	}
}

func TestPropagationConfig_ShouldUseDelta_ForcedMode(t *testing.T) {
	t.Run("forced classic", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.DefaultMode = PropagationModeClassic

		delta := NewFactDelta("Test~1", "Test")
		delta.FieldCount = 10
		delta.AddFieldChange("field1", "old", "new")

		if config.ShouldUseDelta(delta, 5) {
			t.Error("Expected false when mode forced to Classic")
		}
	})

	t.Run("forced delta", func(t *testing.T) {
		config := DefaultPropagationConfig()
		config.DefaultMode = PropagationModeDelta

		delta := NewFactDelta("Test~1", "Test")
		delta.FieldCount = 1
		delta.AddFieldChange("field1", "old", "new")

		if !config.ShouldUseDelta(delta, 5) {
			t.Error("Expected true when mode forced to Delta")
		}
	})
}

func TestPropagationConfig_ShouldUseDelta_MinFields(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.MinFieldsForDelta = 5

	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 3
	delta.AddFieldChange("field1", "old", "new")

	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when field count below threshold")
	}
}

func TestPropagationConfig_ShouldUseDelta_ChangeRatio(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.DeltaThreshold = 0.3
	config.MinFieldsForDelta = 0

	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10

	for i := 0; i < 4; i++ {
		delta.AddFieldChange("field"+string(rune('0'+i)), "old", "new")
	}

	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when change ratio exceeds threshold")
	}
}

func TestPropagationConfig_ShouldUseDelta_AffectedNodes(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.MaxAffectedNodesForDelta = 10
	config.MinFieldsForDelta = 0

	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("field1", "old", "new")

	if config.ShouldUseDelta(delta, 15) {
		t.Error("Expected false when affected nodes exceed limit")
	}
}

func TestPropagationConfig_ShouldUseDelta_PrimaryKeyChange(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.AllowPrimaryKeyChange = false
	config.PrimaryKeyFields = []string{"id", "pk"}
	config.MinFieldsForDelta = 0

	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("id", "123", "456")

	if config.ShouldUseDelta(delta, 5) {
		t.Error("Expected false when PK changed and not allowed")
	}
}

func TestPropagationConfig_ShouldUseDelta_AllConditionsPass(t *testing.T) {
	config := DefaultPropagationConfig()
	config.DefaultMode = PropagationModeAuto
	config.MinFieldsForDelta = 5
	config.DeltaThreshold = 0.5
	config.MaxAffectedNodesForDelta = 50

	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = 10
	delta.AddFieldChange("field1", "old", "new")

	if !config.ShouldUseDelta(delta, 20) {
		t.Error("Expected true when all conditions pass")
	}
}

func TestPropagationConfig_hasPrimaryKeyChange(t *testing.T) {
	tests := []struct {
		name       string
		pkFields   []string
		delta      *FactDelta
		wantChange bool
	}{
		{
			name:     "no PK fields configured",
			pkFields: []string{},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("id", "1", "2")
				return d
			}(),
			wantChange: false,
		},
		{
			name:     "PK field changed",
			pkFields: []string{"id"},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("id", "1", "2")
				return d
			}(),
			wantChange: true,
		},
		{
			name:     "PK field not changed",
			pkFields: []string{"id"},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("name", "old", "new")
				return d
			}(),
			wantChange: false,
		},
		{
			name:     "multiple PK fields, one changed",
			pkFields: []string{"id", "tenant_id"},
			delta: func() *FactDelta {
				d := NewFactDelta("Test~1", "Test")
				d.AddFieldChange("tenant_id", "A", "B")
				return d
			}(),
			wantChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultPropagationConfig()
			config.PrimaryKeyFields = tt.pkFields

			got := config.hasPrimaryKeyChange(tt.delta)
			if got != tt.wantChange {
				t.Errorf("hasPrimaryKeyChange() = %v, want %v", got, tt.wantChange)
			}
		})
	}
}

func TestPropagationConfig_Clone(t *testing.T) {
	original := PropagationConfig{
		DefaultMode:                 PropagationModeDelta,
		EnableDeltaPropagation:      true,
		DeltaThreshold:              0.3,
		MinFieldsForDelta:           5,
		MaxAffectedNodesForDelta:    50,
		AllowPrimaryKeyChange:       true,
		PrimaryKeyFields:            []string{"id", "pk"},
		EnableMetrics:               true,
		PropagationTimeout:          10 * time.Second,
		RetryOnError:                true,
		MaxConcurrentPropagations:   20,
		EnableOptimisticPropagation: true,
		LogPropagationDetails:       true,
	}

	cloned := original.Clone()

	if cloned.DefaultMode != original.DefaultMode {
		t.Error("DefaultMode not cloned")
	}
	if cloned.DeltaThreshold != original.DeltaThreshold {
		t.Error("DeltaThreshold not cloned")
	}

	if len(cloned.PrimaryKeyFields) != len(original.PrimaryKeyFields) {
		t.Error("PrimaryKeyFields length mismatch")
	}

	cloned.PrimaryKeyFields[0] = "modified"
	if original.PrimaryKeyFields[0] == "modified" {
		t.Error("Clone not independent (slice mutation affected original)")
	}
}
