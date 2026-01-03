// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestDefaultDetectorConfig(t *testing.T) {
	config := DefaultDetectorConfig()

	if config.FloatEpsilon != DefaultFloatEpsilon {
		t.Errorf("Expected FloatEpsilon = %v, got %v", DefaultFloatEpsilon, config.FloatEpsilon)
	}

	if !config.IgnoreInternalFields {
		t.Error("Expected IgnoreInternalFields = true")
	}

	if !config.TrackTypeChanges {
		t.Error("Expected TrackTypeChanges = true")
	}

	if !config.EnableDeepComparison {
		t.Error("Expected EnableDeepComparison = true")
	}

	if config.MaxNestingLevel != 10 {
		t.Errorf("Expected MaxNestingLevel = 10, got %d", config.MaxNestingLevel)
	}

	if config.CacheComparisons {
		t.Error("Expected CacheComparisons = false by default")
	}
}

func TestDetectorConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		config    DetectorConfig
		wantError bool
	}{
		{
			name:      "valid default config",
			config:    DefaultDetectorConfig(),
			wantError: false,
		},
		{
			name: "negative epsilon",
			config: DetectorConfig{
				FloatEpsilon:    -0.1,
				MaxNestingLevel: 10,
				CacheTTL:        time.Minute,
			},
			wantError: true,
		},
		{
			name: "zero nesting level",
			config: DetectorConfig{
				FloatEpsilon:    DefaultFloatEpsilon,
				MaxNestingLevel: 0,
				CacheTTL:        time.Minute,
			},
			wantError: true,
		},
		{
			name: "negative cache TTL",
			config: DetectorConfig{
				FloatEpsilon:    DefaultFloatEpsilon,
				MaxNestingLevel: 10,
				CacheTTL:        -time.Second,
			},
			wantError: true,
		},
		{
			name: "edge case: zero epsilon",
			config: DetectorConfig{
				FloatEpsilon:    0,
				MaxNestingLevel: 1,
				CacheTTL:        0,
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

func TestDetectorConfig_ShouldIgnoreField(t *testing.T) {
	tests := []struct {
		name      string
		config    DetectorConfig
		fieldName string
		want      bool
	}{
		{
			name:      "normal field",
			config:    DefaultDetectorConfig(),
			fieldName: "price",
			want:      false,
		},
		{
			name:      "internal field with underscore",
			config:    DefaultDetectorConfig(),
			fieldName: "_internal",
			want:      true,
		},
		{
			name: "internal field but not ignored",
			config: DetectorConfig{
				IgnoreInternalFields: false,
			},
			fieldName: "_internal",
			want:      false,
		},
		{
			name: "explicitly ignored field",
			config: DetectorConfig{
				IgnoredFields: []string{"timestamp", "updated_at"},
			},
			fieldName: "timestamp",
			want:      true,
		},
		{
			name: "not in ignored list",
			config: DetectorConfig{
				IgnoredFields: []string{"timestamp"},
			},
			fieldName: "price",
			want:      false,
		},
		{
			name: "multiple ignored fields",
			config: DetectorConfig{
				IgnoredFields: []string{"field1", "field2", "field3"},
			},
			fieldName: "field2",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.ShouldIgnoreField(tt.fieldName)
			if got != tt.want {
				t.Errorf("ShouldIgnoreField(%s) = %v, want %v", tt.fieldName, got, tt.want)
			}
		})
	}
}

func TestDetectorConfig_Clone(t *testing.T) {
	original := DetectorConfig{
		FloatEpsilon:         0.001,
		IgnoreInternalFields: true,
		IgnoredFields:        []string{"field1", "field2"},
		TrackTypeChanges:     true,
		EnableDeepComparison: true,
		MaxNestingLevel:      5,
		CacheComparisons:     true,
		CacheTTL:             2 * time.Minute,
	}

	cloned := original.Clone()

	// Vérifier égalité des valeurs
	if cloned.FloatEpsilon != original.FloatEpsilon {
		t.Error("FloatEpsilon not cloned correctly")
	}
	if cloned.IgnoreInternalFields != original.IgnoreInternalFields {
		t.Error("IgnoreInternalFields not cloned correctly")
	}
	if cloned.MaxNestingLevel != original.MaxNestingLevel {
		t.Error("MaxNestingLevel not cloned correctly")
	}

	// Vérifier que les slices sont des copies indépendantes
	if len(cloned.IgnoredFields) != len(original.IgnoredFields) {
		t.Error("IgnoredFields length mismatch")
	}

	// Modifier le clone ne doit pas affecter l'original
	cloned.IgnoredFields[0] = "modified"
	if original.IgnoredFields[0] == "modified" {
		t.Error("Clone is not independent (slice mutation affected original)")
	}

	cloned.FloatEpsilon = 999
	if original.FloatEpsilon == 999 {
		t.Error("Clone is not independent (field mutation affected original)")
	}
}

func TestInvalidConfigError_Error(t *testing.T) {
	err := newInvalidConfigError("TestField", "test message")
	expected := "invalid detector config [TestField]: test message"

	if err.Error() != expected {
		t.Errorf("Error() = %s, want %s", err.Error(), expected)
	}
}
