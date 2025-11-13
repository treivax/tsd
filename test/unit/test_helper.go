package unit

import (
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
