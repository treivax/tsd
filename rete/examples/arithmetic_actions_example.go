// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

//go:build ignore
// +build ignore

// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/treivax/tsd/rete"
)

// ArithmeticActionsExample démontre l'utilisation d'expressions arithmétiques dans les actions
func ArithmeticActionsExample() {
	fmt.Println("🧮 EXEMPLE: Expressions Arithmétiques dans les Actions")
	fmt.Println("======================================================")
	fmt.Println()

	// Créer le réseau RETE
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Capturer les sorties de l'action print
	var output bytes.Buffer
	printAction := rete.NewPrintAction(&output)
	network.ActionExecutor.GetRegistry().Register(printAction)

	// Définir les types
	defineTypes(network)

	// Scénario 1: Calcul de l'âge du parent à la naissance
	fmt.Println("📊 Scénario 1: Calcul d'âge parent/enfant")
	fmt.Println("------------------------------------------")
	scenario1_ParentChildAge(network, &output)
	fmt.Println()

	// Scénario 2: Calcul de facture avec TVA
	fmt.Println("📊 Scénario 2: Calcul de facture avec TVA")
	fmt.Println("------------------------------------------")
	scenario2_InvoiceCalculation(network, &output)
	fmt.Println()

	// Scénario 3: Bonus salarial calculé
	fmt.Println("📊 Scénario 3: Calcul de bonus salarial")
	fmt.Println("----------------------------------------")
	scenario3_SalaryBonus(network, &output)
	fmt.Println()

	// Scénario 4: Opérations arithmétiques complexes
	fmt.Println("📊 Scénario 4: Calculs complexes")
	fmt.Println("---------------------------------")
	scenario4_ComplexCalculations(network, &output)
	fmt.Println()

	fmt.Println("✅ Tous les exemples ont été exécutés avec succès!")
}

func main() {
	ArithmeticActionsExample()
}

func defineTypes(network *rete.ReteNetwork) {
	network.Types = []rete.TypeDefinition{
		{
			Type: "typeDefinition",
			Name: "Adulte",
			Fields: []rete.Field{
				{Name: "ID", Type: "string"},
				{Name: "nom", Type: "string"},
				{Name: "age", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Enfant",
			Fields: []rete.Field{
				{Name: "ID", Type: "string"},
				{Name: "nom", Type: "string"},
				{Name: "pere", Type: "string"},
				{Name: "age", Type: "number"},
				{Name: "differenceAgeParent", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Naissance",
			Fields: []rete.Field{
				{Name: "id", Type: "string"},
				{Name: "enfantNom", Type: "string"},
				{Name: "parent", Type: "string"},
				{Name: "ageParentALaNaissance", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Product",
			Fields: []rete.Field{
				{Name: "id", Type: "string"},
				{Name: "name", Type: "string"},
				{Name: "price", Type: "number"},
				{Name: "quantity", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Invoice",
			Fields: []rete.Field{
				{Name: "id", Type: "string"},
				{Name: "productId", Type: "string"},
				{Name: "subtotal", Type: "number"},
				{Name: "tax", Type: "number"},
				{Name: "total", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Employee",
			Fields: []rete.Field{
				{Name: "id", Type: "string"},
				{Name: "name", Type: "string"},
				{Name: "salary", Type: "number"},
				{Name: "bonus", Type: "number"},
				{Name: "performance", Type: "number"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Calculator",
			Fields: []rete.Field{
				{Name: "id", Type: "string"},
				{Name: "a", Type: "number"},
				{Name: "b", Type: "number"},
				{Name: "c", Type: "number"},
				{Name: "result", Type: "number"},
			},
		},
	}
}

func scenario1_ParentChildAge(network *rete.ReteNetwork, output *bytes.Buffer) {
	// Créer les faits
	adulte := &rete.Fact{
		ID:   "adult1",
		Type: "Adulte",
		Fields: map[string]interface{}{
			"ID":  "A001",
			"nom": "Jean",
			"age": float64(45),
		},
	}

	enfant := &rete.Fact{
		ID:   "child1",
		Type: "Enfant",
		Fields: map[string]interface{}{
			"ID":                  "E001",
			"nom":                 "Pierre",
			"pere":                "A001",
			"age":                 float64(18),
			"differenceAgeParent": float64(0),
		},
	}

	fmt.Printf("👤 Adulte: %s (ID: %s, âge: %.0f ans)\n",
		adulte.Fields["nom"], adulte.Fields["ID"], adulte.Fields["age"])
	fmt.Printf("👶 Enfant: %s (ID: %s, âge: %.0f ans, père: %s)\n",
		enfant.Fields["nom"], enfant.Fields["ID"], enfant.Fields["age"], enfant.Fields["pere"])

	// Créer le token
	token := &rete.Token{
		ID:    "token1",
		Facts: []*rete.Fact{adulte, enfant},
		Bindings: map[string]*rete.Fact{
			"a": adulte,
			"e": enfant,
		},
	}

	// Action 1: Créer un fait Naissance avec calcul de l'âge
	action1 := &rete.Action{
		Jobs: []rete.JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factCreation",
						"typeName": "Naissance",
						"fields": map[string]interface{}{
							"id": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "e",
								"field":  "ID",
							},
							"enfantNom": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "e",
								"field":  "nom",
							},
							"parent": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "a",
								"field":  "ID",
							},
							"ageParentALaNaissance": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "-",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "a",
									"field":  "age",
								},
								"right": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "e",
									"field":  "age",
								},
							},
						},
					},
				},
			},
		},
	}

	// Action 2: Modifier l'enfant pour ajouter la différence d'âge
	action2 := &rete.Action{
		Jobs: []rete.JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factModification",
						"variable": "e",
						"field":    "differenceAgeParent",
						"value": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "-",
							"left": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "a",
								"field":  "age",
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "e",
								"field":  "age",
							},
						},
					},
				},
			},
		},
	}

	// Exécuter les actions
	err := network.ActionExecutor.ExecuteAction(action1, token)
	if err != nil {
		log.Printf("❌ Erreur action1: %v", err)
	} else {
		fmt.Printf("✅ Calcul: Âge du parent à la naissance = %.0f - %.0f = %.0f ans\n",
			adulte.Fields["age"], enfant.Fields["age"],
			adulte.Fields["age"].(float64)-enfant.Fields["age"].(float64))
	}

	err = network.ActionExecutor.ExecuteAction(action2, token)
	if err != nil {
		log.Printf("❌ Erreur action2: %v", err)
	} else {
		fmt.Printf("✅ Champ 'differenceAgeParent' calculé et ajouté à l'enfant\n")
	}
}

func scenario2_InvoiceCalculation(network *rete.ReteNetwork, output *bytes.Buffer) {
	// Créer un produit
	product := &rete.Fact{
		ID:   "prod1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":       "P001",
			"name":     "Laptop",
			"price":    float64(1000),
			"quantity": float64(3),
		},
	}

	fmt.Printf("🛒 Produit: %s (Prix: %.2f€, Quantité: %.0f)\n",
		product.Fields["name"], product.Fields["price"], product.Fields["quantity"])

	token := &rete.Token{
		ID:       "token2",
		Facts:    []*rete.Fact{product},
		Bindings: map[string]*rete.Fact{"prod": product},
	}

	// Créer une facture avec calculs de sous-total, TVA et total
	action := &rete.Action{
		Jobs: []rete.JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factCreation",
						"typeName": "Invoice",
						"fields": map[string]interface{}{
							"id": map[string]interface{}{
								"type":  "string",
								"value": "INV001",
							},
							"productId": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "prod",
								"field":  "id",
							},
							// Sous-total = prix * quantité
							"subtotal": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "prod",
									"field":  "price",
								},
								"right": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "prod",
									"field":  "quantity",
								},
							},
							// TVA = (prix * quantité) * 0.20
							"tax": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "*",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "prod",
										"field":  "price",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "prod",
										"field":  "quantity",
									},
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(0.20),
								},
							},
							// Total = (prix * quantité) * 1.20
							"total": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "*",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "prod",
										"field":  "price",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "prod",
										"field":  "quantity",
									},
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(1.20),
								},
							},
						},
					},
				},
			},
		},
	}

	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		log.Printf("❌ Erreur: %v", err)
	} else {
		price := product.Fields["price"].(float64)
		qty := product.Fields["quantity"].(float64)
		subtotal := price * qty
		tax := subtotal * 0.20
		total := subtotal * 1.20

		fmt.Printf("✅ Facture créée:\n")
		fmt.Printf("   - Sous-total: %.2f€ (%.2f × %.0f)\n", subtotal, price, qty)
		fmt.Printf("   - TVA (20%%): %.2f€\n", tax)
		fmt.Printf("   - Total TTC: %.2f€\n", total)
	}
}

func scenario3_SalaryBonus(network *rete.ReteNetwork, output *bytes.Buffer) {
	// Créer un employé
	employee := &rete.Fact{
		ID:   "emp1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":          "E001",
			"name":        "Alice",
			"salary":      float64(50000),
			"performance": float64(1.15), // 115% de performance
			"bonus":       float64(0),
		},
	}

	fmt.Printf("👔 Employé: %s (Salaire: %.2f€, Performance: %.0f%%)\n",
		employee.Fields["name"], employee.Fields["salary"],
		employee.Fields["performance"].(float64)*100)

	token := &rete.Token{
		ID:       "token3",
		Facts:    []*rete.Fact{employee},
		Bindings: map[string]*rete.Fact{"emp": employee},
	}

	// Calculer le bonus = salaire * (performance - 1) * 0.5
	action := &rete.Action{
		Jobs: []rete.JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factModification",
						"variable": "emp",
						"field":    "bonus",
						"value": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "*",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "emp",
									"field":  "salary",
								},
								"right": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "-",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "emp",
										"field":  "performance",
									},
									"right": map[string]interface{}{
										"type":  "number",
										"value": float64(1.0),
									},
								},
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": float64(0.5),
							},
						},
					},
				},
			},
		},
	}

	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		log.Printf("❌ Erreur: %v", err)
	} else {
		salary := employee.Fields["salary"].(float64)
		perf := employee.Fields["performance"].(float64)
		bonus := salary * (perf - 1) * 0.5

		fmt.Printf("✅ Bonus calculé:\n")
		fmt.Printf("   - Formule: %.2f × (%.2f - 1) × 0.5 = %.2f€\n", salary, perf, bonus)
	}
}

func scenario4_ComplexCalculations(network *rete.ReteNetwork, output *bytes.Buffer) {
	// Créer un fait avec plusieurs valeurs
	calc := &rete.Fact{
		ID:   "calc1",
		Type: "Calculator",
		Fields: map[string]interface{}{
			"id":     "C001",
			"a":      float64(10),
			"b":      float64(5),
			"c":      float64(2),
			"result": float64(0),
		},
	}

	fmt.Printf("🧮 Valeurs: a=%.0f, b=%.0f, c=%.0f\n",
		calc.Fields["a"], calc.Fields["b"], calc.Fields["c"])

	token := &rete.Token{
		ID:       "token4",
		Facts:    []*rete.Fact{calc},
		Bindings: map[string]*rete.Fact{"x": calc},
	}

	// Calcul complexe: result = (a * b) + c
	action := &rete.Action{
		Jobs: []rete.JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factModification",
						"variable": "x",
						"field":    "result",
						"value": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "x",
									"field":  "a",
								},
								"right": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "x",
									"field":  "b",
								},
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "x",
								"field":  "c",
							},
						},
					},
				},
			},
		},
	}

	err := network.ActionExecutor.ExecuteAction(action, token)
	if err != nil {
		log.Printf("❌ Erreur: %v", err)
	} else {
		a := calc.Fields["a"].(float64)
		b := calc.Fields["b"].(float64)
		c := calc.Fields["c"].(float64)
		result := (a * b) + c

		fmt.Printf("✅ Résultat: (%.0f × %.0f) + %.0f = %.0f\n", a, b, c, result)
	}
}
