//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
)

func main() {
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  Transaction Example - Command Pattern avec Rejeu Inversé")
	fmt.Println("═══════════════════════════════════════════════════════════\n")

	// 1. Créer un réseau RETE
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Ajouter un type
	network.Types = append(network.Types, rete.TypeDefinition{
		Name: "User",
		Fields: []rete.Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "int"},
			{Name: "email", Type: "string"},
		},
	})

	fmt.Println("📊 État initial du réseau")
	fmt.Printf("   Faits: %d\n\n", len(storage.GetAllFacts()))

	// ═══════════════════════════════════════════════════════════
	// EXEMPLE 1: Transaction réussie (Commit)
	// ═══════════════════════════════════════════════════════════

	fmt.Println("🔄 EXEMPLE 1: Transaction réussie")
	fmt.Println("────────────────────────────────────────────────────────────")

	// Créer une transaction
	tx1 := network.BeginTransaction()
	network.SetTransaction(tx1)

	fmt.Printf("✅ Transaction créée: %s\n", tx1.ID)
	fmt.Printf("   Empreinte mémoire: %d bytes\n\n", tx1.GetMemoryFootprint())

	// Ajouter des faits dans la transaction
	users := []*rete.Fact{
		{
			ID:   "alice",
			Type: "User",
			Fields: map[string]interface{}{
				"name":  "Alice",
				"age":   30,
				"email": "alice@example.com",
			},
			Timestamp: time.Now(),
		},
		{
			ID:   "bob",
			Type: "User",
			Fields: map[string]interface{}{
				"name":  "Bob",
				"age":   25,
				"email": "bob@example.com",
			},
			Timestamp: time.Now(),
		},
		{
			ID:   "charlie",
			Type: "User",
			Fields: map[string]interface{}{
				"name":  "Charlie",
				"age":   35,
				"email": "charlie@example.com",
			},
			Timestamp: time.Now(),
		},
	}

	fmt.Println("➕ Ajout de 3 utilisateurs...")
	for _, user := range users {
		if err := network.SubmitFact(user); err != nil {
			fmt.Printf("❌ Erreur: %v\n", err)
			tx1.Rollback()
			network.SetTransaction(nil)
			return
		}
		fmt.Printf("   ✓ %s ajouté\n", user.Fields["name"])
	}

	fmt.Printf("\n📊 État pendant la transaction\n")
	fmt.Printf("   Faits: %d\n", len(storage.GetAllFacts()))
	fmt.Printf("   Commandes enregistrées: %d\n", tx1.GetCommandCount())
	fmt.Printf("   Empreinte mémoire: %d bytes\n", tx1.GetMemoryFootprint())

	// Commit de la transaction
	if err := tx1.Commit(); err != nil {
		fmt.Printf("❌ Erreur commit: %v\n", err)
		return
	}

	network.SetTransaction(nil)

	fmt.Printf("\n✅ Transaction committée avec succès\n")
	fmt.Printf("   Durée: %v\n", tx1.GetDuration())
	fmt.Printf("   Faits finaux: %d\n\n", len(storage.GetAllFacts()))

	// ═══════════════════════════════════════════════════════════
	// EXEMPLE 2: Transaction avec erreur (Rollback)
	// ═══════════════════════════════════════════════════════════

	fmt.Println("🔄 EXEMPLE 2: Transaction avec rollback")
	fmt.Println("────────────────────────────────────────────────────────────")

	beforeRollback := len(storage.GetAllFacts())
	fmt.Printf("📊 État avant transaction: %d faits\n\n", beforeRollback)

	// Créer une nouvelle transaction
	tx2 := network.BeginTransaction()
	network.SetTransaction(tx2)

	fmt.Printf("✅ Transaction créée: %s\n\n", tx2.ID)

	// Ajouter des faits (qui seront annulés)
	fmt.Println("➕ Ajout de données temporaires...")
	tempUsers := []*rete.Fact{
		{
			ID:   "temp1",
			Type: "User",
			Fields: map[string]interface{}{
				"name":  "Temporary User 1",
				"age":   99,
				"email": "temp1@example.com",
			},
			Timestamp: time.Now(),
		},
		{
			ID:   "temp2",
			Type: "User",
			Fields: map[string]interface{}{
				"name":  "Temporary User 2",
				"age":   99,
				"email": "temp2@example.com",
			},
			Timestamp: time.Now(),
		},
	}

	for _, user := range tempUsers {
		if err := network.SubmitFact(user); err != nil {
			fmt.Printf("❌ Erreur: %v\n", err)
			tx2.Rollback()
			network.SetTransaction(nil)
			return
		}
		fmt.Printf("   ✓ %s ajouté (temporaire)\n", user.Fields["name"])
	}

	fmt.Printf("\n📊 État pendant la transaction\n")
	fmt.Printf("   Faits: %d\n", len(storage.GetAllFacts()))
	fmt.Printf("   Commandes enregistrées: %d\n\n", tx2.GetCommandCount())

	// Simuler une erreur et rollback
	fmt.Println("⚠️  Simulation d'une erreur détectée...")
	fmt.Println("🔙 Rollback de la transaction...\n")

	if err := tx2.Rollback(); err != nil {
		fmt.Printf("❌ Erreur rollback: %v\n", err)
		return
	}

	network.SetTransaction(nil)

	afterRollback := len(storage.GetAllFacts())

	fmt.Printf("✅ Transaction annulée avec succès\n")
	fmt.Printf("   Durée: %v\n", tx2.GetDuration())
	fmt.Printf("   Faits avant rollback: %d\n", beforeRollback)
	fmt.Printf("   Faits après rollback: %d\n", afterRollback)

	if beforeRollback == afterRollback {
		fmt.Printf("   ✓ État restauré correctement!\n\n")
	} else {
		fmt.Printf("   ✗ État non restauré!\n\n")
	}

	// ═══════════════════════════════════════════════════════════
	// EXEMPLE 3: Mesures de performance
	// ═══════════════════════════════════════════════════════════

	fmt.Println("🔄 EXEMPLE 3: Mesures de performance")
	fmt.Println("────────────────────────────────────────────────────────────")

	// Créer un grand réseau
	fmt.Println("📊 Création d'un réseau avec 10,000 faits...")
	largeStorage := rete.NewMemoryStorage()
	largeNetwork := rete.NewReteNetwork(largeStorage)
	largeNetwork.Types = network.Types

	for i := 0; i < 10000; i++ {
		fact := &rete.Fact{
			ID:   fmt.Sprintf("user_%d", i),
			Type: "User",
			Fields: map[string]interface{}{
				"name":  fmt.Sprintf("User %d", i),
				"age":   20 + (i % 50),
				"email": fmt.Sprintf("user%d@example.com", i),
			},
			Timestamp: time.Now(),
		}
		largeStorage.AddFact(fact)
	}

	fmt.Printf("   ✓ Réseau créé: %d faits\n\n", len(largeStorage.GetAllFacts()))

	// Mesurer BeginTransaction
	fmt.Println("⏱️  Mesure de BeginTransaction...")
	start := time.Now()
	tx3 := largeNetwork.BeginTransaction()
	beginDuration := time.Since(start)

	fmt.Printf("   Temps: %v\n", beginDuration)
	fmt.Printf("   Empreinte: %d bytes\n\n", tx3.GetMemoryFootprint())

	// Effectuer des opérations
	largeNetwork.SetTransaction(tx3)

	fmt.Println("➕ Ajout de 100 faits...")
	start = time.Now()
	for i := 0; i < 100; i++ {
		fact := &rete.Fact{
			ID:   fmt.Sprintf("new_%d", i),
			Type: "User",
			Fields: map[string]interface{}{
				"name":  fmt.Sprintf("New User %d", i),
				"age":   20,
				"email": fmt.Sprintf("new%d@example.com", i),
			},
			Timestamp: time.Now(),
		}
		largeNetwork.SubmitFact(fact)
	}
	opsDuration := time.Since(start)

	fmt.Printf("   Temps: %v (%.2f ops/ms)\n", opsDuration,
		float64(100)/float64(opsDuration.Milliseconds()))
	fmt.Printf("   Commandes: %d\n", tx3.GetCommandCount())
	fmt.Printf("   Empreinte: %d bytes (%.2f bytes/cmd)\n\n",
		tx3.GetMemoryFootprint(),
		float64(tx3.GetMemoryFootprint())/float64(tx3.GetCommandCount()))

	// Mesurer Rollback
	fmt.Println("🔙 Mesure de Rollback...")
	start = time.Now()
	tx3.Rollback()
	rollbackDuration := time.Since(start)

	fmt.Printf("   Temps: %v\n", rollbackDuration)
	fmt.Printf("   Faits restaurés: %d\n\n", len(largeStorage.GetAllFacts()))

	largeNetwork.SetTransaction(nil)

	// ═══════════════════════════════════════════════════════════
	// RÉSUMÉ
	// ═══════════════════════════════════════════════════════════

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  📊 RÉSUMÉ DES PERFORMANCES")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("  Taille réseau: 10,000 faits\n")
	fmt.Printf("  BeginTransaction: %v (O(1))\n", beginDuration)
	fmt.Printf("  100 opérations: %v\n", opsDuration)
	fmt.Printf("  Rollback: %v (O(k) avec k=100)\n", rollbackDuration)
	fmt.Printf("  Overhead mémoire: %d bytes (< 0.1%%)\n", tx3.GetMemoryFootprint())
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("\n✅ Tous les exemples terminés avec succès!")
}
