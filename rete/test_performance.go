package rete

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// BenchmarkIndexedFactStorage teste les performances du stockage indexé
func BenchmarkIndexedFactStorage(b *testing.B) {
	config := IndexConfig{
		IndexedFields:        []string{"id", "name", "age", "type"},
		MaxCacheSize:         10000,
		CacheTTL:             5 * time.Minute,
		EnableCompositeIndex: true,
		AutoIndexThreshold:   100,
	}

	storage := NewIndexedFactStorage(config)

	// Pré-remplir avec des données de test
	facts := generateTestFacts(1000)
	for _, fact := range facts {
		storage.StoreFact(fact)
	}

	b.ResetTimer()

	b.Run("StoreFact", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fact := &Fact{
				ID:   fmt.Sprintf("test_%d", i),
				Type: "TestType",
				Fields: map[string]interface{}{
					"id":   i,
					"name": fmt.Sprintf("name_%d", i),
					"age":  rand.Intn(100),
				},
				Timestamp: time.Now(),
			}
			storage.StoreFact(fact)
		}
	})

	b.Run("GetFactByID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			factID := fmt.Sprintf("fact_%d", i%1000)
			storage.GetFactByID(factID)
		}
	})

	b.Run("GetFactsByType", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			factType := fmt.Sprintf("Type_%d", i%10)
			storage.GetFactsByType(factType)
		}
	})

	b.Run("GetFactsByField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			age := i % 100
			storage.GetFactsByField("age", age)
		}
	})
}

// BenchmarkHashJoinEngine teste les performances du moteur de jointure hash
func BenchmarkHashJoinEngine(b *testing.B) {
	config := JoinConfig{
		InitialHashSize:       1024,
		GrowthFactor:          2.0,
		OptimizationThreshold: 1000,
		EnableJoinCache:       true,
		JoinCacheTTL:          time.Minute,
		MaxCacheEntries:       1000,
	}

	engine := NewHashJoinEngine(config)
	cache := NewJoinCache(1000, time.Minute)

	// Condition de jointure de test
	joinCondition := &OptimizedJoinCondition{
		LeftField:  "id",
		RightField: "person_id",
		Operator:   "==",
	}

	// Préparer des données de test
	tokens := generateTestTokens(500)
	facts := generateTestFacts(500)

	// Ajouter aux hash tables
	for _, token := range tokens {
		engine.AddLeftToken(token, joinCondition)
	}
	for _, fact := range facts {
		domainFact := &domain.Fact{
			ID:        fact.ID,
			Type:      fact.Type,
			Fields:    fact.Fields,
			Timestamp: fact.Timestamp,
		}
		engine.AddRightFact(domainFact, joinCondition)
	}

	b.ResetTimer()

	b.Run("PerformHashJoin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			engine.PerformHashJoin(joinCondition, cache)
		}
	})

	b.Run("AddLeftToken", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			token := &domain.Token{
				ID:    fmt.Sprintf("token_%d", i),
				Facts: []*domain.Fact{generateRandomFact(i)},
			}
			engine.AddLeftToken(token, joinCondition)
		}
	})

	b.Run("AddRightFact", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fact := generateRandomFact(i)
			engine.AddRightFact(fact, joinCondition)
		}
	})
}

// BenchmarkEvaluationCache teste les performances du cache d'évaluation
func BenchmarkEvaluationCache(b *testing.B) {
	config := CacheConfig{
		MaxSize:              10000,
		DefaultTTL:           5 * time.Minute,
		CleanupInterval:      time.Minute,
		PrecomputeThreshold:  10,
		EnableKeyCompression: true,
		MaxKeyLength:         100,
	}

	cache := NewEvaluationCache(config)

	// Pré-remplir le cache
	for i := 0; i < 1000; i++ {
		key := &EvaluationKey{
			ConditionType: "binary_operation",
			FactType:      "TestType",
			FieldName:     "age",
			Operator:      ">=",
			Value:         i % 100,
			FactID:        fmt.Sprintf("fact_%d", i),
		}
		cache.Put(key, i%2 == 0, nil, time.Microsecond)
	}

	b.ResetTimer()

	b.Run("CacheHit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := &EvaluationKey{
				ConditionType: "binary_operation",
				FactType:      "TestType",
				FieldName:     "age",
				Operator:      ">=",
				Value:         i % 100,
				FactID:        fmt.Sprintf("fact_%d", i%1000),
			}
			cache.Get(key)
		}
	})

	b.Run("CacheMiss", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := &EvaluationKey{
				ConditionType: "logical_expression",
				FactType:      "NewType",
				FieldName:     "name",
				Operator:      "==",
				Value:         fmt.Sprintf("new_value_%d", i),
				FactID:        fmt.Sprintf("new_fact_%d", i),
			}
			cache.Get(key)
		}
	})

	b.Run("CachePut", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := &EvaluationKey{
				ConditionType: "function_call",
				FactType:      "TestType",
				FieldName:     "value",
				Operator:      "CONTAINS",
				Value:         fmt.Sprintf("test_%d", i),
				FactID:        fmt.Sprintf("fact_%d", i),
			}
			cache.Put(key, true, nil, time.Millisecond)
		}
	})
}

// BenchmarkTokenPropagationEngine teste les performances du moteur de propagation
func BenchmarkTokenPropagationEngine(b *testing.B) {
	config := PropagationConfig{
		NumWorkers:               4,
		BatchSize:                100,
		BatchTimeout:             10 * time.Millisecond,
		EnablePrioritization:     true,
		TimePriorityFactor:       0.001,
		ComplexityPriorityFactor: 0.1,
		MaxQueueSize:             10000,
	}

	engine := NewTokenPropagationEngine(config)

	b.ResetTimer()

	b.Run("EnqueueToken", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			token := &domain.Token{
				ID:    fmt.Sprintf("token_%d", i),
				Facts: []*domain.Fact{generateRandomFact(i)},
			}
			nodeID := fmt.Sprintf("node_%d", i%10)
			priority := float64(i % 100)

			engine.EnqueueToken(token, nodeID, priority)
		}
	})

	b.Run("DequeueHighestPriority", func(b *testing.B) {
		// Pré-remplir la queue
		for i := 0; i < 1000; i++ {
			token := &domain.Token{
				ID:    fmt.Sprintf("token_%d", i),
				Facts: []*domain.Fact{generateRandomFact(i)},
			}
			engine.EnqueueToken(token, "test_node", float64(i))
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			engine.dequeueHighestPriority()
		}
	})
}

// Tests de stress et de charge
func TestStressIndexedStorage(t *testing.T) {
	config := IndexConfig{
		IndexedFields:        []string{"id", "name", "age", "department"},
		MaxCacheSize:         50000,
		CacheTTL:             10 * time.Minute,
		EnableCompositeIndex: true,
		AutoIndexThreshold:   1000,
	}

	storage := NewIndexedFactStorage(config)

	// Test de charge avec de gros volumes
	numFacts := 100000
	facts := make([]*Fact, numFacts)

	// Mesurer le temps d'insertion
	start := time.Now()

	for i := 0; i < numFacts; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("fact_%d", i),
			Type: fmt.Sprintf("Type_%d", i%100),
			Fields: map[string]interface{}{
				"id":         i,
				"name":       fmt.Sprintf("name_%d", i),
				"age":        rand.Intn(100),
				"department": fmt.Sprintf("dept_%d", i%20),
				"salary":     rand.Intn(100000) + 30000,
			},
			Timestamp: time.Now(),
		}
		facts[i] = fact

		if err := storage.StoreFact(fact); err != nil {
			t.Fatalf("Failed to store fact %d: %v", i, err)
		}
	}

	insertTime := time.Since(start)
	t.Logf("Inserted %d facts in %v (%.2f facts/sec)",
		numFacts, insertTime, float64(numFacts)/insertTime.Seconds())

	// Test de recherche après insertion
	searchStart := time.Now()

	for i := 0; i < 1000; i++ {
		// Test différents types de recherche
		factType := fmt.Sprintf("Type_%d", rand.Intn(100))
		results := storage.GetFactsByType(factType)

		if len(results) == 0 && i < 100 { // Les premiers types devraient avoir des résultats
			t.Errorf("No results found for type %s", factType)
		}
	}

	searchTime := time.Since(searchStart)
	t.Logf("Performed 1000 searches in %v (%.2f searches/sec)",
		searchTime, 1000.0/searchTime.Seconds())

	// Vérifier les statistiques d'accès
	stats := storage.GetAccessStats()
	if len(stats) == 0 {
		t.Error("No access statistics recorded")
	}

	t.Logf("Access statistics: %d different access patterns", len(stats))
}

func TestStressHashJoinEngine(t *testing.T) {
	config := JoinConfig{
		InitialHashSize:       2048,
		GrowthFactor:          2.0,
		OptimizationThreshold: 5000,
		EnableJoinCache:       true,
		JoinCacheTTL:          5 * time.Minute,
		MaxCacheEntries:       5000,
	}

	engine := NewHashJoinEngine(config)
	cache := NewJoinCache(5000, 5*time.Minute)

	joinCondition := &OptimizedJoinCondition{
		LeftField:  "user_id",
		RightField: "id",
		Operator:   "==",
	}

	// Générer de gros volumes de données
	numTokens := 10000
	numFacts := 10000

	// Ajouter des tokens
	tokenStart := time.Now()
	for i := 0; i < numTokens; i++ {
		token := &domain.Token{
			ID: fmt.Sprintf("token_%d", i),
			Facts: []*domain.Fact{
				{
					ID:   fmt.Sprintf("user_fact_%d", i),
					Type: "UserAction",
					Fields: map[string]interface{}{
						"user_id": i % 1000, // Créer des overlaps pour les jointures
						"action":  fmt.Sprintf("action_%d", i%10),
					},
					Timestamp: time.Now(),
				},
			},
		}

		if err := engine.AddLeftToken(token, joinCondition); err != nil {
			t.Fatalf("Failed to add token %d: %v", i, err)
		}
	}
	tokenTime := time.Since(tokenStart)

	// Ajouter des faits
	factStart := time.Now()
	for i := 0; i < numFacts; i++ {
		fact := &domain.Fact{
			ID:   fmt.Sprintf("user_%d", i%1000),
			Type: "User",
			Fields: map[string]interface{}{
				"id":   i % 1000,
				"name": fmt.Sprintf("user_%d", i%1000),
				"role": fmt.Sprintf("role_%d", i%5),
			},
			Timestamp: time.Now(),
		}

		if err := engine.AddRightFact(fact, joinCondition); err != nil {
			t.Fatalf("Failed to add fact %d: %v", i, err)
		}
	}
	factTime := time.Since(factStart)

	t.Logf("Added %d tokens in %v, %d facts in %v",
		numTokens, tokenTime, numFacts, factTime)

	// Effectuer des jointures
	joinStart := time.Now()

	for i := 0; i < 100; i++ {
		results, err := engine.PerformHashJoin(joinCondition, cache)
		if err != nil {
			t.Fatalf("Join %d failed: %v", i, err)
		}

		if i == 0 {
			t.Logf("First join produced %d results", len(results))
		}
	}

	joinTime := time.Since(joinStart)
	t.Logf("Performed 100 joins in %v (%.2f joins/sec)",
		joinTime, 100.0/joinTime.Seconds())

	// Vérifier les statistiques
	stats := engine.GetStats()
	t.Logf("Join stats - Total: %d, Cache hits: %d, Cache misses: %d, Avg time: %v",
		stats.TotalJoins, stats.CacheHits, stats.CacheMisses, stats.AverageJoinTime)
}

// Fonctions utilitaires pour générer des données de test
func generateTestFacts(count int) []*Fact {
	facts := make([]*Fact, count)

	for i := 0; i < count; i++ {
		facts[i] = &Fact{
			ID:   fmt.Sprintf("fact_%d", i),
			Type: fmt.Sprintf("Type_%d", i%10),
			Fields: map[string]interface{}{
				"id":     i,
				"name":   fmt.Sprintf("name_%d", i),
				"age":    rand.Intn(100),
				"active": i%2 == 0,
			},
			Timestamp: time.Now().Add(-time.Duration(rand.Intn(3600)) * time.Second),
		}
	}

	return facts
}

func generateTestTokens(count int) []*domain.Token {
	tokens := make([]*domain.Token, count)

	for i := 0; i < count; i++ {
		tokens[i] = &domain.Token{
			ID: fmt.Sprintf("token_%d", i),
			Facts: []*domain.Fact{
				{
					ID:   fmt.Sprintf("fact_%d", i),
					Type: fmt.Sprintf("Type_%d", i%5),
					Fields: map[string]interface{}{
						"id":        i,
						"person_id": i % 100, // Pour les jointures
						"value":     rand.Intn(1000),
					},
					Timestamp: time.Now(),
				},
			},
		}
	}

	return tokens
}

func generateRandomFact(id int) *domain.Fact {
	return &domain.Fact{
		ID:   fmt.Sprintf("fact_%d", id),
		Type: fmt.Sprintf("Type_%d", id%10),
		Fields: map[string]interface{}{
			"id":        id,
			"person_id": id % 100,
			"value":     rand.Intn(1000),
			"active":    id%2 == 0,
		},
		Timestamp: time.Now(),
	}
}

// Benchmarks comparatifs entre approches optimisées et non-optimisées
func BenchmarkComparisonFactLookup(b *testing.B) {
	// Préparer des données communes
	facts := generateTestFacts(10000)

	b.Run("IndexedStorage", func(b *testing.B) {
		config := IndexConfig{
			IndexedFields:        []string{"id", "name", "age"},
			MaxCacheSize:         10000,
			EnableCompositeIndex: true,
		}
		storage := NewIndexedFactStorage(config)

		for _, fact := range facts {
			storage.StoreFact(fact)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			factType := fmt.Sprintf("Type_%d", i%10)
			storage.GetFactsByType(factType)
		}
	})

	b.Run("LinearSearch", func(b *testing.B) {
		// Simulation d'une recherche linéaire
		for i := 0; i < b.N; i++ {
			targetType := fmt.Sprintf("Type_%d", i%10)
			results := []*Fact{}

			for _, fact := range facts {
				if fact.Type == targetType {
					results = append(results, fact)
				}
			}
		}
	})
}

func BenchmarkComparisonJoinStrategies(b *testing.B) {
	numTokens := 1000
	numFacts := 1000

	tokens := generateTestTokens(numTokens)
	facts := generateTestFacts(numFacts)

	b.Run("HashJoin", func(b *testing.B) {
		config := JoinConfig{
			InitialHashSize: 1024,
			GrowthFactor:    2.0,
		}
		engine := NewHashJoinEngine(config)

		joinCondition := &OptimizedJoinCondition{
			LeftField:  "id",
			RightField: "person_id",
			Operator:   "==",
		}

		for _, token := range tokens {
			engine.AddLeftToken(token, joinCondition)
		}

		for _, fact := range facts {
			domainFact := &domain.Fact{
				ID:        fact.ID,
				Type:      fact.Type,
				Fields:    fact.Fields,
				Timestamp: fact.Timestamp,
			}
			engine.AddRightFact(domainFact, joinCondition)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			engine.PerformHashJoin(joinCondition, nil)
		}
	})

	b.Run("NestedLoopJoin", func(b *testing.B) {
		// Simulation d'une jointure par boucles imbriquées
		for i := 0; i < b.N; i++ {
			results := []*JoinResult{}

			for _, token := range tokens {
				for _, tokenFact := range token.Facts {
					leftValue, hasLeft := tokenFact.Fields["id"]
					if !hasLeft {
						continue
					}

					for _, fact := range facts {
						rightValue, hasRight := fact.Fields["person_id"]
						if hasRight && fmt.Sprintf("%v", leftValue) == fmt.Sprintf("%v", rightValue) {
							result := &JoinResult{
								LeftToken: token,
								RightFact: &domain.Fact{
									ID:        fact.ID,
									Type:      fact.Type,
									Fields:    fact.Fields,
									Timestamp: fact.Timestamp,
								},
							}
							results = append(results, result)
						}
					}
				}
			}
		}
	})
}
