// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"testing"
)

// TestFactReferenceJoin teste la jointure fact-to-fact avec c.produit == p
func TestFactReferenceJoin(t *testing.T) {
	t.Log("🧪 TEST: Jointure fact-to-fact (c.produit == p)")

	pipeline := NewConstraintPipeline()
	storage := NewMemoryStorage()

	// Programme TSD minimal avec affectation et référence de fait
	program := `
type Produit(#id: string, nom: string, prix: number)
type Commande(#id: string, produit: Produit, quantite: number)

action commande_trouvee(cmd_id: string, prod_nom: string)

rule match_commande : {p: Produit, c: Commande} / c.produit == p
    ==> commande_trouvee(c.id, p.nom)

// Créer les faits avec affectations
p1 = Produit(id: "PROD001", nom: "Laptop", prix: 1000.0)
p2 = Produit(id: "PROD002", nom: "Mouse", prix: 25.0)

Commande(id: "CMD001", produit: p1, quantite: 2)
Commande(id: "CMD002", produit: p2, quantite: 5)
`

	// Créer un fichier temporaire
	tmpFile, err := os.CreateTemp("", "test_fact_ref_*.tsd")
	if err != nil {
		t.Fatalf("❌ Erreur création fichier temporaire: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(program); err != nil {
		t.Fatalf("❌ Erreur écriture fichier temporaire: %v", err)
	}
	tmpFile.Close()

	network, _, err := pipeline.IngestFile(tmpFile.Name(), nil, storage)
	if err != nil {
		t.Fatalf("❌ Erreur construction réseau: %v", err)
	}

	t.Logf("✅ Réseau construit: %d TypeNodes, %d BetaNodes, %d TerminalNodes",
		len(network.TypeNodes), len(network.BetaNodes), len(network.TerminalNodes))

	// Debug: afficher les clés des TypeNodes
	t.Logf("🔍 Clés des TypeNodes:")
	for key := range network.TypeNodes {
		t.Logf("  - %s", key)
	}

	// Vérifier les TypeNodes
	produitNode, okProduit := network.TypeNodes["Produit"]
	commandeNode, okCommande := network.TypeNodes["Commande"]

	if !okProduit || !okCommande {
		t.Fatalf("❌ TypeNodes manquants: Produit=%v, Commande=%v", okProduit, okCommande)
	}

	produitFacts := produitNode.Memory.GetFacts()
	commandeFacts := commandeNode.Memory.GetFacts()

	t.Logf("📦 Faits Produit: %d", len(produitFacts))
	for id, fact := range produitFacts {
		t.Logf("  - %v: %+v", id, fact.Fields)
	}

	t.Logf("📦 Faits Commande: %d", len(commandeFacts))
	for id, fact := range commandeFacts {
		t.Logf("  - %v: ID=%s, produit=%v, quantite=%v",
			id, fact.ID, fact.Fields["produit"], fact.Fields["quantite"])
	}

	// Vérifier les JoinNodes
	t.Logf("🔗 JoinNodes: %d", len(network.BetaNodes))
	for id, node := range network.BetaNodes {
		if jn, ok := node.(*JoinNode); ok {
			leftTokens := jn.LeftMemory.GetTokens()
			rightTokens := jn.RightMemory.GetTokens()
			resultTokens := jn.Memory.GetTokens()

			t.Logf("  JoinNode %v:", id)
			t.Logf("    - LeftMemory: %d tokens", len(leftTokens))
			for tid, token := range leftTokens {
				t.Logf("      • %v (vars=%v)", tid, token.GetVariables())
			}
			t.Logf("    - RightMemory: %d tokens", len(rightTokens))
			for tid, token := range rightTokens {
				t.Logf("      • %v (vars=%v)", tid, token.GetVariables())
			}
			t.Logf("    - ResultMemory: %d tokens", len(resultTokens))
			t.Logf("    - JoinConditions: %+v", jn.JoinConditions)

			if len(jn.JoinConditions) > 0 {
				cond := jn.JoinConditions[0]
				t.Logf("    - Condition détaillée: %s.%s %s %s.%s",
					cond.LeftVar, cond.LeftField, cond.Operator, cond.RightVar, cond.RightField)
			}
		}
	}

	// Vérifier les TerminalNodes
	if len(network.TerminalNodes) != 1 {
		t.Fatalf("❌ Attendu 1 TerminalNode, got %d", len(network.TerminalNodes))
	}

	var terminalNode *TerminalNode
	for _, tn := range network.TerminalNodes {
		terminalNode = tn
		break
	}

	tokens := terminalNode.Memory.GetTokens()
	execCount := terminalNode.GetExecutionCount()

	t.Logf("🎯 TerminalNode: %d tokens, %d exécutions", len(tokens), execCount)

	// On attend 2 tokens/exécutions (une par commande)
	expectedCount := int64(2)
	if execCount != expectedCount {
		t.Errorf("❌ Attendu %d exécutions, got %d", expectedCount, execCount)
		t.Logf("📊 Détails des tokens:")
		for i, token := range tokens {
			t.Logf("  Token %d: %+v", i, token.Bindings)
		}
	} else {
		t.Logf("✅ Test réussi: %d jointures correctes", execCount)
	}
}
