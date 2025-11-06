package main

import (
	"fmt"
	"log"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/rete/pkg/domain"
	"github.com/treivax/tsd/rete/pkg/network"
)

// Exemple d'utilisation des nœuds Beta dans le réseau RETE
func main() {
	fmt.Println("🚀 Démonstration des nœuds Beta dans le réseau RETE")

	// 1. Créer un storage en mémoire
	storage := rete.NewMemoryStorage()

	// 2. Créer le réseau RETE
	reteNetwork := rete.NewReteNetwork(storage)

	// 3. Activer le support des nœuds Beta
	err := reteNetwork.EnableBetaNodes()
	if err != nil {
		log.Fatalf("Erreur activation nœuds Beta: %v", err)
	}

	// 4. Créer un constructeur de réseau Beta indépendant pour la démonstration
	logger := &ConsoleLogger{}
	betaBuilder := network.NewBetaNetworkBuilder(logger)

	// 5. Exemple de pattern multi-jointures : Person -> Address -> Company
	demonstrateBetaJoinPattern(betaBuilder)

	// 6. Démonstration d'intégration avec le réseau principal
	demonstrateReteIntegration(reteNetwork)

	fmt.Println("\n✅ Démonstration terminée avec succès!")
}

// ConsoleLogger implémente une interface de logging simple
type ConsoleLogger struct{}

func (c *ConsoleLogger) Debug(msg string, fields map[string]interface{}) {
	fmt.Printf("DEBUG: %s %v\n", msg, fields)
}

func (c *ConsoleLogger) Info(msg string, fields map[string]interface{}) {
	fmt.Printf("INFO: %s %v\n", msg, fields)
}

func (c *ConsoleLogger) Warn(msg string, fields map[string]interface{}) {
	fmt.Printf("WARN: %s %v\n", msg, fields)
}

func (c *ConsoleLogger) Error(msg string, err error, fields map[string]interface{}) {
	fmt.Printf("ERROR: %s - %v %v\n", msg, err, fields)
}

// demonstrateBetaJoinPattern montre comment créer un pattern de jointures multiples
func demonstrateBetaJoinPattern(builder *network.BetaNetworkBuilder) {
	fmt.Println("\n🔗 Démonstration du pattern de jointures Beta")

	// Créer les conditions de jointure
	personAddressCondition := domain.NewBasicJoinCondition("address_id", "id", "==")
	addressCompanyCondition := domain.NewBasicJoinCondition("company_id", "id", "==")

	// Définir le pattern multi-jointures
	pattern := network.MultiJoinPattern{
		PatternID: "employee_complete_info",
		JoinSpecs: []network.JoinSpecification{
			{
				LeftType:   "Person",
				RightType:  "Address",
				Conditions: []domain.JoinCondition{personAddressCondition},
				NodeID:     "person_address_join",
			},
			{
				LeftType:   "PersonAddress",
				RightType:  "Company",
				Conditions: []domain.JoinCondition{addressCompanyCondition},
				NodeID:     "address_company_join",
			},
		},
		FinalAction: "create_employee_record",
	}

	// Construire le réseau de jointures
	createdNodes, err := builder.BuildMultiJoinNetwork(pattern)
	if err != nil {
		log.Printf("Erreur construction réseau: %v", err)
		return
	}

	fmt.Printf("✅ Pattern créé avec %d nœuds de jointure\n", len(createdNodes))

	// Tester avec des données d'exemple
	testBetaJoinWithSampleData(builder, createdNodes)

	// Afficher les statistiques
	stats := builder.NetworkStatistics()
	fmt.Printf("📊 Statistiques du réseau Beta:\n")
	fmt.Printf("   - Nœuds totaux: %d\n", stats.TotalNodes)
	fmt.Printf("   - Nœuds Beta simples: %d\n", stats.SimpleBetaNodes)
	fmt.Printf("   - Nœuds de jointure: %d\n", stats.JoinNodes)
	fmt.Printf("   - Tokens totaux: %d\n", stats.TotalTokens)
	fmt.Printf("   - Faits totaux: %d\n", stats.TotalFacts)
}

// testBetaJoinWithSampleData teste les jointures avec des données d'exemple
func testBetaJoinWithSampleData(builder *network.BetaNetworkBuilder, nodes []domain.BetaNode) {
	fmt.Println("\n🧪 Test avec des données d'exemple")

	if len(nodes) == 0 {
		fmt.Println("Aucun nœud à tester")
		return
	}

	// Créer des faits d'exemple
	personFact := domain.NewFact("person_1", "Person", map[string]interface{}{
		"id":         "p1",
		"name":       "Jean Dupont",
		"address_id": "a1",
	})

	addressFact := domain.NewFact("address_1", "Address", map[string]interface{}{
		"id":         "a1",
		"street":     "123 Rue de la Paix",
		"city":       "Paris",
		"company_id": "c1",
	})

	companyFact := domain.NewFact("company_1", "Company", map[string]interface{}{
		"id":   "c1",
		"name": "Tech Corp",
		"type": "Technology",
	})

	// Traiter les faits dans le premier nœud (person_address_join)
	firstNode := nodes[0]
	fmt.Printf("Traitement des faits dans le nœud: %s\n", firstNode.ID())

	// Créer un token pour la partie gauche (Person)
	personToken := domain.NewToken("token_1", "person_source", []*domain.Fact{personFact})

	// Traiter le token et les faits
	firstNode.ProcessLeftToken(personToken)
	firstNode.ProcessRightFact(addressFact)

	// Si on a un deuxième nœud, traiter aussi le fait Company
	if len(nodes) > 1 {
		secondNode := nodes[1]
		fmt.Printf("Traitement du fait Company dans le nœud: %s\n", secondNode.ID())
		secondNode.ProcessRightFact(companyFact)
	}

	fmt.Println("✅ Données d'exemple traitées")
}

// demonstrateReteIntegration montre l'intégration avec le réseau RETE principal
func demonstrateReteIntegration(network *rete.ReteNetwork) {
	fmt.Println("\n🌐 Démonstration d'intégration avec le réseau RETE")

	// Créer une jointure Beta dans le réseau principal
	conditions := []interface{}{
		map[string]string{"field": "id", "operator": "=="},
	}

	err := network.CreateBetaJoin("alpha_person", "alpha_address", "person_address_beta", conditions)
	if err != nil {
		log.Printf("Erreur création jointure: %v", err)
		return
	}

	// Obtenir les statistiques
	stats := network.GetBetaNodeStatistics()
	fmt.Printf("📊 Statistiques d'intégration:\n")
	fmt.Printf("   - Nœuds Beta dans le réseau: %v\n", stats["totalBetaNodes"])
	fmt.Printf("   - Support Beta activé: %v\n", stats["betaEnabled"])

	// Afficher la structure du réseau
	network.PrintNetworkStructure()
}
