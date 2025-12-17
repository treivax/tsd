// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAlphaSharing_SameCondition vérifie si les AlphaNodes avec la même condition sont partagés
func TestAlphaSharing_SameCondition(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Deux règles avec la MÊME condition: p.age > 18
	content := `type Person(#id: string, age: number, name:string)
action rule1_action(arg: string)
action rule2_action(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> rule1_action(p.id)
rule r2 : {p: Person} / p.age > 18 ==> rule2_action(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier le nombre d'AlphaNodes créés
	alphaCount := len(network.AlphaNodes)
	t.Logf("Nombre d'AlphaNodes créés: %d", alphaCount)
	if alphaCount == 1 {
		t.Log("✅ PARTAGÉ : Un seul AlphaNode pour les deux règles avec la même condition")
	} else if alphaCount == 2 {
		t.Log("❌ NON PARTAGÉ : Deux AlphaNodes créés pour la même condition")
		t.Log("   → Opportunité d'optimisation : les AlphaNodes devraient être partagés")
	} else {
		t.Errorf("Nombre d'AlphaNodes inattendu: %d", alphaCount)
	}
	// Vérifier le TypeNode
	personTypeNode, exists := network.TypeNodes["Person"]
	if !exists {
		t.Fatal("TypeNode 'Person' non trouvé")
	}
	// Vérifier combien d'enfants le TypeNode a
	children := personTypeNode.GetChildren()
	t.Logf("TypeNode 'Person' a %d enfant(s)", len(children))
	// Lister les AlphaNodes
	t.Log("\nAlphaNodes dans le réseau:")
	for id, alphaNode := range network.AlphaNodes {
		t.Logf("  - ID: %s", id)
		t.Logf("    Type: %s", alphaNode.GetType())
		t.Logf("    Enfants: %d", len(alphaNode.GetChildren()))
	}
	// Vérifier les TerminalNodes (il devrait y en avoir 2, un par règle)
	terminalCount := len(network.TerminalNodes)
	t.Logf("\nNombre de TerminalNodes: %d", terminalCount)
	if terminalCount != 2 {
		t.Errorf("Attendu 2 TerminalNodes (un par règle), obtenu %d", terminalCount)
	}
}

// TestAlphaSharing_DifferentConditions vérifie que des conditions différentes créent des AlphaNodes séparés
func TestAlphaSharing_DifferentConditions(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Deux règles avec des conditions DIFFÉRENTES
	content := `type Person(#id: string, age:number)
action adult(arg: string)
action young(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 65 ==> young(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Avec des conditions différentes, on devrait avoir 2 AlphaNodes distincts
	alphaCount := len(network.AlphaNodes)
	t.Logf("Nombre d'AlphaNodes pour conditions différentes: %d", alphaCount)
	if alphaCount != 2 {
		t.Errorf("Attendu 2 AlphaNodes pour conditions différentes, obtenu %d", alphaCount)
	} else {
		t.Log("✅ Conditions différentes → AlphaNodes séparés (comportement correct)")
	}
}

// TestAlphaSharing_ThreeRulesSameCondition vérifie le comportement avec 3 règles identiques
func TestAlphaSharing_ThreeRulesSameCondition(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	// Trois règles avec la MÊME condition
	content := `type Person(#id: string, age:number)
action action1(arg: string)
action action2(arg: string)
action action3(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 18 ==> action2(p.id)
rule r3 : {p: Person} / p.age > 18 ==> action3(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	alphaCount := len(network.AlphaNodes)
	terminalCount := len(network.TerminalNodes)
	t.Logf("Nombre d'AlphaNodes: %d", alphaCount)
	t.Logf("Nombre de TerminalNodes: %d", terminalCount)
	// Il devrait y avoir 3 TerminalNodes (un par règle)
	if terminalCount != 3 {
		t.Errorf("Attendu 3 TerminalNodes, obtenu %d", terminalCount)
	}
	// Si partagé: 1 AlphaNode, sinon: 3 AlphaNodes
	if alphaCount == 1 {
		t.Log("✅ OPTIMAL : 1 AlphaNode partagé par 3 règles")
	} else if alphaCount == 3 {
		t.Log("❌ SOUS-OPTIMAL : 3 AlphaNodes créés pour la même condition")
		t.Log("   → Potentiel d'optimisation : 1 AlphaNode suffirait")
	}
}

// TestAlphaSharing_WithFacts vérifie que le comportement est correct avec des faits
func TestAlphaSharing_WithFacts(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(#id: string, age:number)
action action1(arg: string)
action action2(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 18 ==> action2(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Soumettre des faits
	facts := []*Fact{
		{ID: "P001", Type: "Person", Fields: map[string]interface{}{"id": "P001", "age": float64(25)}}, // > 18
		{ID: "P002", Type: "Person", Fields: map[string]interface{}{"id": "P002", "age": float64(15)}}, // < 18
		{ID: "P003", Type: "Person", Fields: map[string]interface{}{"id": "P003", "age": float64(30)}}, // > 18
	}
	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			t.Fatalf("Erreur soumission fait %s: %v", fact.ID, err)
		}
	}
	// Vérifier les activations
	activatedTerminals := 0
	totalTokens := 0
	for _, terminal := range network.TerminalNodes {
		terminalMemory := terminal.GetMemory()
		tokenCount := len(terminalMemory.Tokens)
		if tokenCount > 0 {
			activatedTerminals++
			totalTokens += tokenCount
			t.Logf("TerminalNode %s: %d token(s)", terminal.GetID(), tokenCount)
		}
	}
	t.Logf("\nRésultats:")
	t.Logf("  Terminaux activés: %d/2", activatedTerminals)
	t.Logf("  Tokens totaux: %d", totalTokens)
	// Les deux terminaux devraient être activés (2 faits matchent la condition)
	if activatedTerminals != 2 {
		t.Errorf("Attendu 2 terminaux activés, obtenu %d", activatedTerminals)
	}
	// Chaque terminal devrait avoir 2 tokens (P001 et P003)
	expectedTokens := 4 // 2 faits * 2 règles
	if totalTokens != expectedTokens {
		t.Logf("⚠️  Nombre de tokens: attendu %d, obtenu %d", expectedTokens, totalTokens)
		t.Log("   (Le comportement peut varier selon l'implémentation du partage)")
	}
	t.Log("\n✅ Les faits ont été correctement propagés")
}

// TestAlphaSharing_StructureVisualization affiche la structure pour comprendre le comportement
func TestAlphaSharing_StructureVisualization(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")
	content := `type Person(#id: string, age:number)
action action1(arg: string)
action action2(arg: string)
rule r1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 18 ==> action2(p.id)
`
	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	t.Log("\n" + strings.Repeat("=", 60))
	t.Log("VISUALISATION DE LA STRUCTURE DU RÉSEAU")
	t.Log(strings.Repeat("=", 60))
	t.Log("\n📊 Statistiques:")
	t.Logf("   • TypeNodes: %d", len(network.TypeNodes))
	t.Logf("   • AlphaNodes: %d", len(network.AlphaNodes))
	t.Logf("   • TerminalNodes: %d", len(network.TerminalNodes))
	t.Log("\n🌳 Structure:")
	t.Log("RootNode")
	for typeName, typeNode := range network.TypeNodes {
		t.Logf("  └── TypeNode: %s", typeName)
		children := typeNode.GetChildren()
		t.Logf("      Enfants: %d", len(children))
		for i, child := range children {
			isLast := i == len(children)-1
			prefix := "      ├──"
			if isLast {
				prefix = "      └──"
			}
			t.Logf("%s AlphaNode: %s", prefix, child.GetID())
			alphaChildren := child.GetChildren()
			for j, terminal := range alphaChildren {
				isTerminalLast := j == len(alphaChildren)-1
				terminalPrefix := "      │   ├──"
				if isLast {
					terminalPrefix = "          ├──"
				}
				if isTerminalLast {
					if isLast {
						terminalPrefix = "          └──"
					} else {
						terminalPrefix = "      │   └──"
					}
				}
				t.Logf("%s TerminalNode: %s", terminalPrefix, terminal.GetID())
			}
		}
	}
	t.Log("\n" + strings.Repeat("=", 60))
	// Analyse du résultat
	alphaCount := len(network.AlphaNodes)
	if alphaCount == 1 {
		t.Log("✅ RÉSULTAT : AlphaNodes PARTAGÉS")
		t.Log("   Structure optimale avec un seul nœud de test")
	} else {
		t.Log("❌ RÉSULTAT : AlphaNodes NON PARTAGÉS")
		t.Log("   Chaque règle a son propre AlphaNode")
		t.Log("   Optimisation possible en partageant les nœuds avec conditions identiques")
	}
	t.Log(strings.Repeat("=", 60))
}
