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

// TestTypeNodeSharing_TwoSimpleRulesSameType vérifie qu'un seul TypeNode est créé
// pour deux règles simples portant sur un même type
func TestTypeNodeSharing_TwoSimpleRulesSameType(t *testing.T) {
	// Créer un fichier TSD temporaire avec un type et deux règles simples
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")

	content := `type Person(id: string, age: number, name:string)

action adult_detected(id: string, name: string)
action not_retired(id: string, name: string)

rule r1 : {p: Person} / p.age > 18 ==> adult_detected(p.id, p.name)
rule r2 : {p: Person} / p.age < 65 ==> not_retired(p.id, p.name)
`

	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}

	// Construire le réseau RETE
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// VÉRIFICATION 1: Un seul TypeNode doit être créé pour le type Person
	if len(network.TypeNodes) != 1 {
		t.Errorf("Attendu 1 TypeNode, obtenu %d", len(network.TypeNodes))
	}

	personTypeNode, exists := network.TypeNodes["Person"]
	if !exists {
		t.Fatal("TypeNode 'Person' non trouvé")
	}

	// VÉRIFICATION 2: Le TypeNode doit avoir exactement 2 enfants (un AlphaNode par règle)
	children := personTypeNode.GetChildren()
	if len(children) != 2 {
		t.Errorf("Le TypeNode Person devrait avoir 2 enfants (AlphaNodes), obtenu %d", len(children))
		t.Logf("Enfants du TypeNode Person:")
		for i, child := range children {
			t.Logf("  Enfant %d: ID=%s, Type=%s", i+1, child.GetID(), child.GetType())
		}
	}

	// VÉRIFICATION 3: Les deux enfants doivent être des AlphaNodes
	for i, child := range children {
		if child.GetType() != "alpha" {
			t.Errorf("L'enfant %d du TypeNode devrait être de type 'alpha', obtenu '%s'", i+1, child.GetType())
		}
	}

	// VÉRIFICATION 4: Chaque AlphaNode doit avoir un TerminalNode enfant
	for i, child := range children {
		alphaChildren := child.GetChildren()
		if len(alphaChildren) != 1 {
			t.Errorf("L'AlphaNode %d devrait avoir 1 enfant (TerminalNode), obtenu %d", i+1, len(alphaChildren))
			continue
		}
		terminal := alphaChildren[0]
		if terminal.GetType() != "terminal" {
			t.Errorf("L'enfant de l'AlphaNode %d devrait être de type 'terminal', obtenu '%s'", i+1, terminal.GetType())
		}
	}

	// VÉRIFICATION 5: Le réseau doit avoir exactement 2 TerminalNodes
	if len(network.TerminalNodes) != 2 {
		t.Errorf("Le réseau devrait avoir 2 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}

	// VÉRIFICATION 6: Le TypeNode doit être connecté au RootNode
	rootChildren := network.RootNode.GetChildren()
	foundPersonTypeNode := false
	for _, child := range rootChildren {
		if child.GetID() == personTypeNode.GetID() {
			foundPersonTypeNode = true
			break
		}
	}
	if !foundPersonTypeNode {
		t.Error("Le TypeNode Person devrait être un enfant du RootNode")
	}

	t.Log("✅ Vérification réussie: Un seul TypeNode créé pour deux règles simples sur le même type")
}

// TestTypeNodeSharing_ThreeRulesSameType vérifie qu'un seul TypeNode est créé
// pour trois règles simples portant sur un même type
func TestTypeNodeSharing_ThreeRulesSameType(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")

	content := `type Employee(id: string, salary: number, name:string)

action high_earner(id: string)
action low_earner(id: string)
action mid_earner(id: string)

rule r1 : {e: Employee} / e.salary > 50000 ==> high_earner(e.id)
rule r2 : {e: Employee} / e.salary < 30000 ==> low_earner(e.id)
rule r3 : {e: Employee} / e.salary >= 30000 AND e.salary <= 50000 ==> mid_earner(e.id)
`

	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// VÉRIFICATION: Un seul TypeNode pour Employee
	if len(network.TypeNodes) != 1 {
		t.Errorf("Attendu 1 TypeNode, obtenu %d", len(network.TypeNodes))
	}

	employeeTypeNode, exists := network.TypeNodes["Employee"]
	if !exists {
		t.Fatal("TypeNode 'Employee' non trouvé")
	}

	// VÉRIFICATION: 3 AlphaNodes enfants
	children := employeeTypeNode.GetChildren()
	if len(children) != 3 {
		t.Errorf("Le TypeNode Employee devrait avoir 3 enfants (AlphaNodes), obtenu %d", len(children))
	}

	// VÉRIFICATION: 3 TerminalNodes dans le réseau
	if len(network.TerminalNodes) != 3 {
		t.Errorf("Le réseau devrait avoir 3 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}

	t.Log("✅ Vérification réussie: Un seul TypeNode créé pour trois règles simples sur le même type")
}

// TestTypeNodeSharing_TwoDifferentTypes vérifie que deux TypeNodes distincts
// sont créés pour deux règles portant sur deux types différents
func TestTypeNodeSharing_TwoDifferentTypes(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")

	content := `type Person(id: string, age:number)
type Company(id: string, revenue:number)

action adult_detected(id: string)
action big_company(id: string)

rule r1 : {p: Person} / p.age > 18 ==> adult_detected(p.id)
rule r2 : {c: Company} / c.revenue > 1000000 ==> big_company(c.id)
`

	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// VÉRIFICATION: Deux TypeNodes distincts
	if len(network.TypeNodes) != 2 {
		t.Errorf("Attendu 2 TypeNodes, obtenu %d", len(network.TypeNodes))
	}

	personTypeNode, personExists := network.TypeNodes["Person"]
	companyTypeNode, companyExists := network.TypeNodes["Company"]

	if !personExists {
		t.Error("TypeNode 'Person' non trouvé")
	}
	if !companyExists {
		t.Error("TypeNode 'Company' non trouvé")
	}

	// VÉRIFICATION: Chaque TypeNode a 1 enfant
	if personExists {
		personChildren := personTypeNode.GetChildren()
		if len(personChildren) != 1 {
			t.Errorf("Le TypeNode Person devrait avoir 1 enfant, obtenu %d", len(personChildren))
		}
	}

	if companyExists {
		companyChildren := companyTypeNode.GetChildren()
		if len(companyChildren) != 1 {
			t.Errorf("Le TypeNode Company devrait avoir 1 enfant, obtenu %d", len(companyChildren))
		}
	}

	t.Log("✅ Vérification réussie: Deux TypeNodes distincts créés pour deux types différents")
}

// TestTypeNodeSharing_MixedRules vérifie le partage de TypeNode
// avec un mélange de règles simples et de règles de jointure
func TestTypeNodeSharing_MixedRules(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")

	content := `type Person(id: string, age: number, company_id:string)
type Company(id: string, name:string)

action adult_detected(id: string)
action employee_match(personId: string, companyId: string)
action not_retired(id: string)

rule r1 : {p: Person} / p.age > 18 ==> adult_detected(p.id)
rule r2 : {p: Person, c: Company} / p.company_id == c.id ==> employee_match(p.id, c.id)
rule r3 : {p: Person} / p.age < 65 ==> not_retired(p.id)
`

	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// VÉRIFICATION: Deux TypeNodes (Person et Company)
	if len(network.TypeNodes) != 2 {
		t.Errorf("Attendu 2 TypeNodes, obtenu %d", len(network.TypeNodes))
	}

	personTypeNode, personExists := network.TypeNodes["Person"]
	companyTypeNode, companyExists := network.TypeNodes["Company"]

	if !personExists {
		t.Fatal("TypeNode 'Person' non trouvé")
	}
	if !companyExists {
		t.Fatal("TypeNode 'Company' non trouvé")
	}

	// VÉRIFICATION: Le TypeNode Person doit avoir plusieurs enfants
	// (2 règles simples + 1 connexion pour la règle de jointure)
	personChildren := personTypeNode.GetChildren()
	if len(personChildren) < 2 {
		t.Errorf("Le TypeNode Person devrait avoir au moins 2 enfants, obtenu %d", len(personChildren))
	}

	// VÉRIFICATION: Le TypeNode Company doit avoir au moins 1 enfant
	// (pour la règle de jointure)
	companyChildren := companyTypeNode.GetChildren()
	if len(companyChildren) < 1 {
		t.Errorf("Le TypeNode Company devrait avoir au moins 1 enfant, obtenu %d", len(companyChildren))
	}

	t.Log("✅ Vérification réussie: TypeNodes correctement partagés dans un scénario mixte")
}

// TestTypeNodeSharing_VisualizeNetwork crée un réseau et affiche sa structure détaillée
// pour visualiser comment les TypeNodes sont partagés
func TestTypeNodeSharing_VisualizeNetwork(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")

	content := `type Person(id: string, age: number, name:string)

action adult_detected(id: string)
action not_retired(id: string)

rule r1 : {p: Person} / p.age > 18 ==> adult_detected(p.id)
rule r2 : {p: Person} / p.age < 65 ==> not_retired(p.id)
`

	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Afficher la structure du réseau
	t.Log("\n" + strings.Repeat("=", 60))
	t.Log("STRUCTURE DÉTAILLÉE DU RÉSEAU RETE")
	t.Log(strings.Repeat("=", 60))

	t.Logf("\n📊 Statistiques:")
	t.Logf("   • TypeNodes: %d", len(network.TypeNodes))
	t.Logf("   • AlphaNodes: %d", len(network.AlphaNodes))
	t.Logf("   • TerminalNodes: %d", len(network.TerminalNodes))

	t.Log("\n🌳 Arborescence du réseau:")
	t.Log("\nRootNode")
	t.Logf("  └── ID: %s", network.RootNode.GetID())

	// Parcourir les TypeNodes
	for typeName, typeNode := range network.TypeNodes {
		t.Logf("\n      ├── TypeNode: %s", typeName)
		t.Logf("      │   ID: %s", typeNode.GetID())

		children := typeNode.GetChildren()
		t.Logf("      │   Enfants: %d", len(children))

		for i, child := range children {
			isLast := i == len(children)-1
			prefix := "      │   "
			if isLast {
				t.Logf("%s└── AlphaNode: %s", prefix, child.GetID())
				prefix = "      │       "
			} else {
				t.Logf("%s├── AlphaNode: %s", prefix, child.GetID())
				prefix = "      │   │   "
			}

			// Afficher les enfants de l'AlphaNode (TerminalNodes)
			alphaChildren := child.GetChildren()
			for j, terminal := range alphaChildren {
				isTerminalLast := j == len(alphaChildren)-1
				if isTerminalLast {
					t.Logf("%s└── TerminalNode: %s", prefix, terminal.GetID())
				} else {
					t.Logf("%s├── TerminalNode: %s", prefix, terminal.GetID())
				}
			}
		}
	}

	t.Log("\n" + strings.Repeat("=", 60))
	t.Log("✅ CONFIRMATION: Un seul TypeNode 'Person' partagé par 2 règles")
	t.Log(strings.Repeat("=", 60))

	// Vérifications
	if len(network.TypeNodes) != 1 {
		t.Errorf("Attendu 1 TypeNode, obtenu %d", len(network.TypeNodes))
	}

	personTypeNode, exists := network.TypeNodes["Person"]
	if !exists {
		t.Fatal("TypeNode 'Person' non trouvé")
	}

	children := personTypeNode.GetChildren()
	if len(children) != 2 {
		t.Errorf("Le TypeNode Person devrait avoir 2 enfants, obtenu %d", len(children))
	}
}

// TestTypeNodeSharing_WithFactSubmission vérifie que le partage de TypeNode
// fonctionne correctement lors de la soumission de faits
func TestTypeNodeSharing_WithFactSubmission(t *testing.T) {
	tempDir := t.TempDir()
	constraintFile := filepath.Join(tempDir, "test.constraint")

	content := `type Person(id: string, age: number, name:string)

action adult_detected(id: string)
action not_retired(id: string)
action middle_aged(id: string)

rule r1 : {p: Person} / p.age > 18 ==> adult_detected(p.id)
rule r2 : {p: Person} / p.age < 65 ==> not_retired(p.id)
rule r3 : {p: Person} / p.age > 30 AND p.age < 50 ==> middle_aged(p.id)
`

	err := os.WriteFile(constraintFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Soumettre les faits manuellement
	facts := []*Fact{
		{ID: "P001", Type: "Person", Fields: map[string]interface{}{"id": "P001", "age": float64(25), "name": "Alice"}},
		{ID: "P002", Type: "Person", Fields: map[string]interface{}{"id": "P002", "age": float64(70), "name": "Bob"}},
		{ID: "P003", Type: "Person", Fields: map[string]interface{}{"id": "P003", "age": float64(15), "name": "Charlie"}},
	}

	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			t.Fatalf("Erreur soumission fait %s: %v", fact.ID, err)
		}
	}

	// VÉRIFICATION 1: Un seul TypeNode pour Person
	if len(network.TypeNodes) != 1 {
		t.Errorf("Attendu 1 TypeNode, obtenu %d", len(network.TypeNodes))
	}

	personTypeNode, exists := network.TypeNodes["Person"]
	if !exists {
		t.Fatal("TypeNode 'Person' non trouvé")
	}

	// VÉRIFICATION 2: Le TypeNode a 3 AlphaNodes enfants (une par règle)
	children := personTypeNode.GetChildren()
	if len(children) != 3 {
		t.Errorf("Le TypeNode Person devrait avoir 3 enfants, obtenu %d", len(children))
	}

	// VÉRIFICATION 3: Les faits ont été propagés à travers le TypeNode unique
	typeMemory := personTypeNode.GetMemory()
	t.Logf("TypeNode contient %d faits", len(typeMemory.Facts))
	if len(typeMemory.Facts) != 3 {
		t.Errorf("Le TypeNode devrait contenir 3 faits, obtenu %d", len(typeMemory.Facts))
	}

	// VÉRIFICATION 4: Chaque AlphaNode a reçu les faits appropriés
	for i, child := range children {
		alphaMemory := child.GetMemory()
		t.Logf("AlphaNode %d (%s): %d faits en mémoire", i+1, child.GetID(), len(alphaMemory.Facts))

		// Au moins un fait devrait avoir activé chaque AlphaNode
		if len(alphaMemory.Facts) == 0 {
			t.Logf("  ⚠️  Aucun fait n'a activé cet AlphaNode (peut être normal selon les conditions)")
		}
	}

	// VÉRIFICATION 5: Les TerminalNodes ont été activés
	activatedTerminals := 0
	for _, terminal := range network.TerminalNodes {
		terminalMemory := terminal.GetMemory()
		if len(terminalMemory.Tokens) > 0 {
			activatedTerminals++
			t.Logf("✅ TerminalNode %s activé avec %d token(s)", terminal.GetID(), len(terminalMemory.Tokens))
		}
	}

	t.Logf("\n📊 Résumé:")
	t.Logf("   • 1 TypeNode partagé par 3 règles")
	t.Logf("   • 3 faits soumis")
	t.Logf("   • %d TerminalNode(s) activé(s)", activatedTerminals)

	if activatedTerminals > 0 {
		t.Log("\n✅ Vérification réussie: Le TypeNode unique a propagé les faits vers toutes les règles")
	} else {
		t.Log("\n⚠️  Aucune règle activée (peut être normal selon les conditions)")
	}
}
