package rete

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// TestPrintNetworkStructure teste la fonction PrintNetworkStructure
func TestPrintNetworkStructure(t *testing.T) {
	// Capturer la sortie standard
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Créer un réseau RETE simple
	network := NewReteNetwork(NewMemoryStorage())

	// Créer un AST de test
	ast := &Program{
		Types: []TypeDefinition{
			{
				Type: "typeDefinition",
				Name: "TestEvent",
				Fields: []Field{
					{Name: "priority", Type: "int"},
					{Name: "status", Type: "string"},
				},
			},
		},
		Expressions: []Expression{
			{
				Type: "expression",
				Set: Set{
					Type: "set",
					Variables: []TypedVariable{
						{
							Type:     "typedVariable",
							Name:     "event",
							DataType: "TestEvent",
						},
					},
				},
				Constraints: map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "==",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "event",
						"field":  "priority",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": float64(1),
					},
				},
				Action: &Action{
					Type: "action",
					Job: JobCall{
						Type: "job",
						Name: "print",
						Args: []string{"High priority event detected"},
					},
				},
			},
		},
	}

	// Charger l'AST dans le réseau
	err := network.LoadFromAST(ast)
	if err != nil {
		t.Fatalf("Erreur lors du chargement de l'AST: %v", err)
	}

	// Appeler PrintNetworkStructure
	network.PrintNetworkStructure()

	// Fermer le writer et récupérer la sortie
	w.Close()
	os.Stdout = oldStdout

	// Lire la sortie capturée
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Vérifications
	if output == "" {
		t.Fatal("PrintNetworkStructure n'a produit aucune sortie")
	}

	// Vérifier que la sortie contient les éléments attendus
	expectedElements := []string{
		"📊 STRUCTURE DU RÉSEAU RETE:",
		"Root:",
		"TypeNode[TestEvent]:",
		"AlphaNode:",
		"TerminalNode:",
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Fatalf("La sortie ne contient pas '%s'. Sortie: %s", element, output)
		}
	}

	// Vérifier la structure hiérarchique (présence de caractères de structure)
	if !strings.Contains(output, "├──") || !strings.Contains(output, "└──") {
		t.Fatalf("La sortie ne contient pas la structure hiérarchique attendue. Sortie: %s", output)
	}

	fmt.Printf("Sortie capturée de PrintNetworkStructure:\n%s", output)
}

// TestPrintNetworkStructureEmpty teste PrintNetworkStructure avec un réseau vide
func TestPrintNetworkStructureEmpty(t *testing.T) {
	// Capturer la sortie standard
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Créer un réseau RETE vide
	network := NewReteNetwork(NewMemoryStorage())

	// Appeler PrintNetworkStructure sur un réseau vide
	network.PrintNetworkStructure()

	// Fermer le writer et récupérer la sortie
	w.Close()
	os.Stdout = oldStdout

	// Lire la sortie capturée
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Vérifications pour un réseau vide
	if output == "" {
		t.Fatal("PrintNetworkStructure n'a produit aucune sortie même pour un réseau vide")
	}

	if !strings.Contains(output, "📊 STRUCTURE DU RÉSEAU RETE:") {
		t.Fatalf("La sortie ne contient pas l'en-tête attendu. Sortie: %s", output)
	}

	if !strings.Contains(output, "Root:") {
		t.Fatalf("La sortie ne contient pas l'information sur le nœud racine. Sortie: %s", output)
	}
}
