// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/treivax/tsd/rete"
)

// TestHelper contient des fonctions utilitaires pour les tests RETE
type TestHelper struct {
	pipeline *rete.ConstraintPipeline
}

// NewTestHelper crée un nouvel helper de test
func NewTestHelper() *TestHelper {
	return &TestHelper{
		pipeline: rete.NewConstraintPipeline(),
	}
}

// BuildNetworkFromConstraintFile utilise le pipeline unique pour construire un réseau RETE
// Cette fonction DOIT être utilisée par TOUS les tests qui utilisent des fichiers .constraint
func (th *TestHelper) BuildNetworkFromConstraintFile(t *testing.T, constraintFile string) (*rete.ReteNetwork, rete.Storage) {
	storage := rete.NewMemoryStorage()

	network, err := th.pipeline.BuildNetworkFromConstraintFile(constraintFile, storage)
	if err != nil {
		t.Fatalf("❌ Erreur pipeline constraint → RETE: %v", err)
	}

	return network, storage
}

// BuildNetworkFromConstraintFileWithMassiveFacts utilise le pipeline avec fichiers .constraint et .facts
func (th *TestHelper) BuildNetworkFromConstraintFileWithMassiveFacts(t *testing.T, constraintFile, factsFile string) (*rete.ReteNetwork, []*rete.Fact, rete.Storage) {
	storage := rete.NewMemoryStorage()

	network, facts, err := th.pipeline.BuildNetworkFromConstraintFileWithFacts(constraintFile, factsFile, storage)
	if err != nil {
		t.Fatalf("❌ Erreur pipeline constraint + faits → RETE: %v", err)
	}

	return network, facts, storage
}

// CreateUserFact crée un fait utilisateur standardisé pour les tests
func (th *TestHelper) CreateUserFact(id, nom, prenom string, age float64) *rete.Fact {
	return &rete.Fact{
		ID:   "fact_u_" + id,
		Type: "Utilisateur",
		Fields: map[string]interface{}{
			"id":     id,
			"nom":    nom,
			"prenom": prenom,
			"age":    age,
		},
	}
}

// CreateAddressFact crée un fait adresse standardisé pour les tests
func (th *TestHelper) CreateAddressFact(userID, rue, ville string) *rete.Fact {
	return &rete.Fact{
		ID:   "fact_a_" + userID,
		Type: "Adresse",
		Fields: map[string]interface{}{
			"utilisateur_id": userID,
			"rue":            rue,
			"ville":          ville,
		},
	}
}

// CreateCustomerFact crée un fait customer standardisé pour les tests
func (th *TestHelper) CreateCustomerFact(id string, age float64, vip bool) *rete.Fact {
	return &rete.Fact{
		ID:   "fact_c_" + id,
		Type: "Customer",
		Fields: map[string]interface{}{
			"id":  id,
			"age": age,
			"vip": vip,
		},
	}
}

// SubmitFactsAndAnalyze soumet des faits et analyse le tuple-space
func (th *TestHelper) SubmitFactsAndAnalyze(t *testing.T, network *rete.ReteNetwork, facts []*rete.Fact) int {
	totalActions := 0

	// Soumettre les faits
	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			t.Logf("⚠️ Erreur soumission fait %s: %v", fact.ID, err)
		}
	}

	// Analyser les résultats
	for terminalID, terminal := range network.TerminalNodes {
		tokenCount := len(terminal.Memory.Tokens)
		totalActions += tokenCount

		t.Logf("Terminal %s: %d tuples stockés", terminalID, tokenCount)
	}

	return totalActions
}

// ShowFactDetails affiche les détails complets d'un fait avec tous ses attributs
func (th *TestHelper) ShowFactDetails(fact *rete.Fact, index int) string {
	if fact == nil || fact.Fields == nil {
		return ""
	}

	// Construire une représentation complète du fait avec id en premier
	var sortedAttrs []string
	if id, exists := fact.Fields["id"]; exists {
		sortedAttrs = append(sortedAttrs, fmt.Sprintf("id=%v", id))
	}
	for key, value := range fact.Fields {
		if key != "id" {
			sortedAttrs = append(sortedAttrs, fmt.Sprintf("%s=%v", key, value))
		}
	}

	return fmt.Sprintf("[%d] %s{%s}", index, fact.Type, strings.Join(sortedAttrs, ", "))
}

// ShowActionDetailsWithAllAttributes affiche les détails d'une action avec tous les attributs des faits
func (th *TestHelper) ShowActionDetailsWithAllAttributes(actionName string, terminal *rete.TerminalNode, maxResults int) {
	count := 0
	for _, token := range terminal.Memory.Tokens {
		if count >= maxResults {
			break
		}

		// Extraire et afficher tous les faits du token avec leurs attributs complets
		if len(token.Facts) > 0 {
			fmt.Printf("   → %s:\n", actionName)
			for i, fact := range token.Facts {
				if fact != nil {
					fmt.Printf("     %s\n", th.ShowFactDetails(fact, i+1))
				}
			}
		}
		count++
	}

	if len(terminal.Memory.Tokens) > maxResults {
		fmt.Printf("   ... et %d autres résultats\n", len(terminal.Memory.Tokens)-maxResults)
	}
}
