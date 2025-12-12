package rete

import (
	"strings"
	"testing"
)

// TestAlphaChainBuilderGetMetrics tests the GetMetrics function
func TestAlphaChainBuilderGetMetrics(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)
	metrics := builder.GetMetrics()
	if metrics == nil {
		t.Fatal("expected non-nil metrics")
	}
	// Metrics should be initialized
	if metrics.TotalChainsBuilt < 0 {
		t.Error("expected non-negative TotalChainsBuilt")
	}
	if metrics.TotalNodesReused < 0 {
		t.Error("expected non-negative TotalNodesReused")
	}
	if metrics.TotalNodesCreated < 0 {
		t.Error("expected non-negative TotalNodesCreated")
	}
}

// TestAlphaChainBuilderGetConnectionCacheSize tests the GetConnectionCacheSize function
func TestAlphaChainBuilderGetConnectionCacheSize(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)
	size := builder.GetConnectionCacheSize()
	if size < 0 {
		t.Errorf("expected non-negative cache size, got %d", size)
	}
	// Initially should be 0 or small
	if size > 1000 {
		t.Errorf("unexpected initial cache size: %d", size)
	}
}

// TestAlphaChainBuilderClearConnectionCache tests the ClearConnectionCache function
func TestAlphaChainBuilderClearConnectionCache(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)
	// Clear cache (should not panic even if empty)
	builder.ClearConnectionCache()
	// Size should be 0 after clearing
	size := builder.GetConnectionCacheSize()
	if size != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", size)
	}
}

// TestBetaChainBuilderGetConnectionCacheSize tests the GetConnectionCacheSize function on BetaChainBuilder
func TestBetaChainBuilderGetConnectionCacheSize(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	betaBuilder := network.BetaChainBuilder
	size := betaBuilder.GetConnectionCacheSize()
	if size < 0 {
		t.Errorf("expected non-negative cache size, got %d", size)
	}
}

// TestComputeJoinNodeHash tests the ComputeJoinNodeHash function
func TestComputeJoinNodeHash(t *testing.T) {
	tests := []struct {
		name      string
		canonical *CanonicalJoinSignature
		wantError bool
	}{
		{
			name: "valid signature",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{"x", "y"},
				RightVars: []string{"z"},
				AllVars:   []string{"x", "y", "z"},
				Condition: map[string]interface{}{"op": "eq"},
				VarTypes: []VariableTypeMapping{
					{VarName: "x", TypeName: "TypeX"},
					{VarName: "y", TypeName: "TypeY"},
					{VarName: "z", TypeName: "TypeZ"},
				},
			},
			wantError: false,
		},
		{
			name: "empty signature",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{},
				RightVars: []string{},
				AllVars:   []string{},
				Condition: map[string]interface{}{},
				VarTypes:  []VariableTypeMapping{},
			},
			wantError: false,
		},
		{
			name: "complex nested condition",
			canonical: &CanonicalJoinSignature{
				Version:   "1.0",
				LeftVars:  []string{"a", "b", "c"},
				RightVars: []string{"d", "e"},
				AllVars:   []string{"a", "b", "c", "d", "e"},
				Condition: map[string]interface{}{
					"type": "and",
					"conditions": []interface{}{
						map[string]interface{}{"op": "eq", "field": "name"},
						map[string]interface{}{"op": "gt", "field": "age"},
					},
				},
				VarTypes: []VariableTypeMapping{
					{VarName: "a", TypeName: "TypeA"},
					{VarName: "b", TypeName: "TypeB"},
					{VarName: "c", TypeName: "TypeC"},
					{VarName: "d", TypeName: "TypeD"},
					{VarName: "e", TypeName: "TypeE"},
				},
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := ComputeJoinNodeHash(tt.canonical)
			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !strings.HasPrefix(hash, "join_") {
				t.Errorf("expected hash to start with 'join_', got: %s", hash)
			}
			if len(hash) <= len("join_") {
				t.Error("expected hash to have content after prefix")
			}
		})
	}
}

// TestComputeSHA256Hash tests the ComputeSHA256Hash function
func TestComputeSHA256Hash(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected int // expected length of hex string
	}{
		{
			name:     "simple data",
			data:     []byte("hello world"),
			expected: 16, // 8 bytes * 2 hex chars per byte
		},
		{
			name:     "empty data",
			data:     []byte{},
			expected: 16,
		},
		{
			name:     "complex data",
			data:     []byte(`{"type":"join","vars":["x","y","z"]}`),
			expected: 16,
		},
		{
			name:     "numeric data",
			data:     []byte{1, 2, 3, 4, 5},
			expected: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := ComputeSHA256Hash(tt.data)
			if len(hash) != tt.expected {
				t.Errorf("expected hash length %d, got %d", tt.expected, len(hash))
			}
			// Hash should be hex string (only 0-9 and a-f)
			for _, c := range hash {
				if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
					t.Errorf("hash contains non-hex character: %c", c)
				}
			}
			// Same input should produce same hash
			hash2 := ComputeSHA256Hash(tt.data)
			if hash != hash2 {
				t.Error("expected consistent hash for same input")
			}
		})
	}
}

// TestComputeJoinNodeHashConsistency verifies that identical signatures produce identical hashes
func TestComputeJoinNodeHashConsistency(t *testing.T) {
	canonical := &CanonicalJoinSignature{
		Version:   "1.0",
		LeftVars:  []string{"a", "b"},
		RightVars: []string{"c"},
		AllVars:   []string{"a", "b", "c"},
		Condition: map[string]interface{}{
			"operator": "equals",
			"field":    "id",
		},
		VarTypes: []VariableTypeMapping{
			{VarName: "a", TypeName: "TypeA"},
			{VarName: "b", TypeName: "TypeB"},
			{VarName: "c", TypeName: "TypeC"},
		},
	}
	// Compute hash multiple times
	hashes := make([]string, 10)
	for i := 0; i < 10; i++ {
		hash, err := ComputeJoinNodeHash(canonical)
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}
		hashes[i] = hash
	}
	// All hashes should be identical
	firstHash := hashes[0]
	for i := 1; i < len(hashes); i++ {
		if hashes[i] != firstHash {
			t.Errorf("hash inconsistency: iteration 0 produced %s, iteration %d produced %s", firstHash, i, hashes[i])
		}
	}
}

// TestComputeSHA256HashDifferentInputs verifies that different inputs produce different hashes
func TestComputeSHA256HashDifferentInputs(t *testing.T) {
	data1 := []byte("input1")
	data2 := []byte("input2")
	data3 := []byte("different")
	hash1 := ComputeSHA256Hash(data1)
	hash2 := ComputeSHA256Hash(data2)
	hash3 := ComputeSHA256Hash(data3)
	if hash1 == hash2 {
		t.Error("expected different hashes for different inputs (1 vs 2)")
	}
	if hash1 == hash3 {
		t.Error("expected different hashes for different inputs (1 vs 3)")
	}
	if hash2 == hash3 {
		t.Error("expected different hashes for different inputs (2 vs 3)")
	}
}

// TestComputeSHA256HashEmptyVsNonEmpty verifies empty and non-empty produce different hashes
func TestComputeSHA256HashEmptyVsNonEmpty(t *testing.T) {
	emptyHash := ComputeSHA256Hash([]byte{})
	nonEmptyHash := ComputeSHA256Hash([]byte("x"))
	if emptyHash == nonEmptyHash {
		t.Error("expected different hashes for empty vs non-empty input")
	}
}
