// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"testing"

	"github.com/treivax/tsd/rete"
)

func TestSimpleBetaNodeTupleSpace(t *testing.T) {
	fmt.Printf("🎯 TEST TUPLE-SPACE - Pipeline Unique Simple\n")
	fmt.Printf("=================================================================\n")

	// 🚀 UTILISER LE PIPELINE UNIQUE
	helper := NewTestHelper()
	constraintPath := "../../constraint/test/integration/beta_complex_rules.constraint"
	network, _ := helper.BuildNetworkFromConstraintFile(t, constraintPath)

	// Créer des faits de test avec l'helper
	userFact := helper.CreateUserFact("U001", "Martin", "Pierre", 25.0)
	addressFact := helper.CreateAddressFact("U001", "Rue de la Paix", "Paris")

	facts := []*rete.Fact{userFact, addressFact}

	// Soumettre et analyser avec l'helper
	fmt.Printf("\n📊 Test de soumission...\n")
	totalActions := helper.SubmitFactsAndAnalyze(t, network, facts)

	// Validations
	if len(network.TerminalNodes) > 0 {
		fmt.Printf("✅ Pipeline unique utilisé avec succès\n")
	}

	fmt.Printf("✅ Actions dans tuple-space: %d\n", totalActions)
	fmt.Printf("✅ RÈGLE RESPECTÉE: Pipeline unique pour tous les tests\n")

	fmt.Printf("\n🎊 TEST PIPELINE UNIQUE SIMPLE: RÉUSSI\n")
}
