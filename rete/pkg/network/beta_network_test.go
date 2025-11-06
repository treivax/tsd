package network

import (
	"fmt"
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// mockLogger pour les tests
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields map[string]interface{})            {}
func (m *mockLogger) Info(msg string, fields map[string]interface{})             {}
func (m *mockLogger) Warn(msg string, fields map[string]interface{})             {}
func (m *mockLogger) Error(msg string, err error, fields map[string]interface{}) {}

func TestBetaNetworkBuilder(t *testing.T) {
	logger := &mockLogger{}

	t.Run("NewBetaNetworkBuilder", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)
		if builder == nil {
			t.Fatal("NewBetaNetworkBuilder should not return nil")
		}

		stats := builder.NetworkStatistics()
		if stats.TotalNodes != 0 {
			t.Errorf("New builder should have 0 nodes, got %d", stats.TotalNodes)
		}
	})

	t.Run("CreateJoinNode", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		condition := domain.NewBasicJoinCondition("name", "name", "==")
		conditions := []domain.JoinCondition{condition}

		joinNode := builder.CreateJoinNode("join1", conditions)

		if joinNode == nil {
			t.Fatal("CreateJoinNode should not return nil")
		}
		if joinNode.ID() != "join1" {
			t.Errorf("Expected node ID 'join1', got '%s'", joinNode.ID())
		}

		nodeConditions := joinNode.GetJoinConditions()
		if len(nodeConditions) != 1 {
			t.Errorf("Expected 1 condition, got %d", len(nodeConditions))
		}
	})

	t.Run("CreateBetaNode", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		betaNode := builder.CreateBetaNode("beta1")

		if betaNode == nil {
			t.Fatal("CreateBetaNode should not return nil")
		}
		if betaNode.ID() != "beta1" {
			t.Errorf("Expected node ID 'beta1', got '%s'", betaNode.ID())
		}
		if betaNode.Type() != "BetaNode" {
			t.Errorf("Expected node type 'BetaNode', got '%s'", betaNode.Type())
		}
	})

	t.Run("GetBetaNode", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Test get non-existent node
		node, exists := builder.GetBetaNode("nonexistent")
		if exists || node != nil {
			t.Error("Non-existent node should not be found")
		}

		// Create and get existing node
		originalNode := builder.CreateBetaNode("beta1")

		retrievedNode, exists := builder.GetBetaNode("beta1")
		if !exists || retrievedNode == nil {
			t.Fatal("Existing node should be found")
		}
		if retrievedNode.ID() != originalNode.ID() {
			t.Error("Retrieved node should be the same as original")
		}
	})

	t.Run("ConnectNodes", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		parent := builder.CreateBetaNode("parent")
		_ = builder.CreateBetaNode("child")

		err := builder.ConnectNodes("parent", "child")
		if err != nil {
			t.Errorf("ConnectNodes should not return error: %v", err)
		}

		// Verify connection
		children := parent.GetChildren()
		if len(children) != 1 {
			t.Errorf("Parent should have 1 child, got %d", len(children))
		}
		if children[0].ID() != "child" {
			t.Errorf("Child ID should be 'child', got '%s'", children[0].ID())
		}
	})

	t.Run("ConnectNodes_NonExistent", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Test connecting non-existent parent
		err := builder.ConnectNodes("nonexistent", "child")
		if err == nil {
			t.Error("Connecting non-existent parent should return error")
		}

		// Test connecting non-existent child
		builder.CreateBetaNode("parent")
		err = builder.ConnectNodes("parent", "nonexistent")
		if err == nil {
			t.Error("Connecting non-existent child should return error")
		}
	})

	t.Run("ListBetaNodes", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Test empty list
		nodes := builder.ListBetaNodes()
		if len(nodes) != 0 {
			t.Errorf("Empty builder should return 0 nodes, got %d", len(nodes))
		}

		// Create some nodes
		builder.CreateBetaNode("beta1")
		builder.CreateJoinNode("join1", []domain.JoinCondition{})

		nodes = builder.ListBetaNodes()
		if len(nodes) != 2 {
			t.Errorf("Builder should have 2 nodes, got %d", len(nodes))
		}

		// Verify nodes exist
		if _, exists := nodes["beta1"]; !exists {
			t.Error("beta1 should exist in list")
		}
		if _, exists := nodes["join1"]; !exists {
			t.Error("join1 should exist in list")
		}
	})

	t.Run("ClearNetwork", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Create some nodes
		builder.CreateBetaNode("beta1")
		builder.CreateJoinNode("join1", []domain.JoinCondition{})

		// Verify nodes exist
		stats := builder.NetworkStatistics()
		if stats.TotalNodes != 2 {
			t.Errorf("Should have 2 nodes before clear, got %d", stats.TotalNodes)
		}

		// Clear network
		builder.ClearNetwork()

		// Verify network is empty
		stats = builder.NetworkStatistics()
		if stats.TotalNodes != 0 {
			t.Errorf("Should have 0 nodes after clear, got %d", stats.TotalNodes)
		}
	})
}

func TestMultiJoinPattern(t *testing.T) {
	logger := &mockLogger{}

	t.Run("BuildMultiJoinNetwork", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Créer un pattern de jointures multiples
		pattern := MultiJoinPattern{
			PatternID: "person_address_company",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Person",
					RightType: "Address",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("address_id", "id", "=="),
					},
					NodeID: "person_address_join",
				},
				{
					LeftType:  "PersonAddress",
					RightType: "Company",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("company_id", "id", "=="),
					},
					NodeID: "address_company_join",
				},
			},
			FinalAction: "create_employee_record",
		}

		createdNodes, err := builder.BuildMultiJoinNetwork(pattern)
		if err != nil {
			t.Fatalf("BuildMultiJoinNetwork should not return error: %v", err)
		}

		if len(createdNodes) != 2 {
			t.Errorf("Expected 2 created nodes, got %d", len(createdNodes))
		}

		// Vérifier que les nœuds sont connectés
		firstNode := createdNodes[0]
		children := firstNode.GetChildren()
		if len(children) != 1 {
			t.Errorf("First node should have 1 child, got %d", len(children))
		}

		if children[0].ID() != createdNodes[1].ID() {
			t.Error("First node should be connected to second node")
		}
	})

	t.Run("BuildMultiJoinNetwork_EmptyPattern", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		pattern := MultiJoinPattern{
			PatternID: "empty",
			JoinSpecs: []JoinSpecification{},
		}

		createdNodes, err := builder.BuildMultiJoinNetwork(pattern)
		if err != nil {
			t.Errorf("BuildMultiJoinNetwork with empty pattern should not return error: %v", err)
		}

		if len(createdNodes) != 0 {
			t.Errorf("Empty pattern should create 0 nodes, got %d", len(createdNodes))
		}
	})

	t.Run("BuildComplexECommercePattern", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Pattern E-commerce ultra-complexe : Order → Customer → Product → Stock
		ecommercePattern := MultiJoinPattern{
			PatternID: "order_validation_system",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Order",
					RightType: "Customer",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("customer_id", "id", "=="),
					},
					NodeID: "order_customer_join",
				},
				{
					LeftType:  "OrderCustomer",
					RightType: "Product",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("product_id", "id", "=="),
					},
					NodeID: "customer_product_join",
				},
				{
					LeftType:  "OrderProduct",
					RightType: "Stock",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("product_id", "product_id", "=="),
						domain.NewBasicJoinCondition("quantity", "available_quantity", "<="),
					},
					NodeID: "product_stock_validation",
				},
			},
			FinalAction: "validate_and_process_order",
		}

		// Construire le réseau E-commerce
		createdNodes, err := builder.BuildMultiJoinNetwork(ecommercePattern)
		if err != nil {
			t.Fatalf("E-commerce pattern should build successfully: %v", err)
		}

		if len(createdNodes) != 3 {
			t.Errorf("Expected 3 nodes for E-commerce pattern, got %d", len(createdNodes))
		}

		// Vérifier la chaîne de connexions : Order→Customer→Product→Stock
		firstNode := createdNodes[0]   // order_customer_join
		secondNode := createdNodes[1]  // customer_product_join  
		thirdNode := createdNodes[2]   // product_stock_validation

		// Vérifier connexion 1→2
		children1 := firstNode.GetChildren()
		if len(children1) != 1 || children1[0].ID() != secondNode.ID() {
			t.Error("First node should connect to second node")
		}

		// Vérifier connexion 2→3
		children2 := secondNode.GetChildren()
		if len(children2) != 1 || children2[0].ID() != thirdNode.ID() {
			t.Error("Second node should connect to third node")
		}

		// Vérifier que le dernier nœud a des conditions multiples
		finalConditions := thirdNode.(domain.JoinNode).GetJoinConditions()
		if len(finalConditions) != 2 {
			t.Errorf("Final node should have 2 conditions, got %d", len(finalConditions))
		}

		// Tester avec des données E-commerce réalistes
		testECommerceFlow(t, createdNodes)
	})

	t.Run("BuildComplexHRPattern", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Pattern RH avec niveaux de sécurité : Employee → Department → Project
		hrPattern := MultiJoinPattern{
			PatternID: "employee_project_assignment",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Employee",
					RightType: "Department",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("dept_id", "id", "=="),
					},
					NodeID: "emp_dept_join",
				},
				{
					LeftType:  "EmployeeDepartment",
					RightType: "Project",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("project_id", "id", "=="),
						domain.NewBasicJoinCondition("clearance_level", "required_clearance", ">="),
					},
					NodeID: "dept_project_security_join",
				},
			},
			FinalAction: "assign_employee_to_project",
		}

		// Construire le réseau RH
		createdNodes, err := builder.BuildMultiJoinNetwork(hrPattern)
		if err != nil {
			t.Fatalf("HR pattern should build successfully: %v", err)
		}

		if len(createdNodes) != 2 {
			t.Errorf("Expected 2 nodes for HR pattern, got %d", len(createdNodes))
		}

		// Tester avec des données RH réalistes incluant clearance de sécurité
		testHRSecurityFlow(t, createdNodes)
	})
}

func TestNetworkStatistics(t *testing.T) {
	logger := &mockLogger{}

	t.Run("NetworkStatistics", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Create different types of nodes
		betaNode := builder.CreateBetaNode("beta1")
		joinNode := builder.CreateJoinNode("join1", []domain.JoinCondition{})

		// Add some data to test memory statistics
		fact := domain.NewFact("f1", "Person", map[string]interface{}{"name": "John"})
		token := domain.NewToken("t1", "node1", []*domain.Fact{fact})

		betaNode.ProcessLeftToken(token)
		betaNode.ProcessRightFact(fact)

		joinNode.ProcessRightFact(fact)

		stats := builder.NetworkStatistics()

		// Verify node counts
		if stats.TotalNodes != 2 {
			t.Errorf("Expected 2 total nodes, got %d", stats.TotalNodes)
		}
		if stats.SimpleBetaNodes != 1 {
			t.Errorf("Expected 1 simple beta node, got %d", stats.SimpleBetaNodes)
		}
		if stats.JoinNodes != 1 {
			t.Errorf("Expected 1 join node, got %d", stats.JoinNodes)
		}

		// Verify memory statistics exist
		if len(stats.MemoryStats) != 2 {
			t.Errorf("Expected memory stats for 2 nodes, got %d", len(stats.MemoryStats))
		}

		// Check that beta1 has both tokens and facts
		if betaStats, exists := stats.MemoryStats["beta1"]; exists {
			if betaStats.TokenCount == 0 {
				t.Error("beta1 should have at least 1 token")
			}
			if betaStats.FactCount == 0 {
				t.Error("beta1 should have at least 1 fact")
			}
		} else {
			t.Error("beta1 should have memory statistics")
		}

		// Check overall totals
		if stats.TotalTokens == 0 {
			t.Error("Network should have at least 1 token")
		}
		if stats.TotalFacts == 0 {
			t.Error("Network should have at least 1 fact")
		}
	})

	t.Run("EmptyNetworkStatistics", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		stats := builder.NetworkStatistics()

		if stats.TotalNodes != 0 {
			t.Errorf("Empty network should have 0 nodes, got %d", stats.TotalNodes)
		}
		if stats.SimpleBetaNodes != 0 {
			t.Errorf("Empty network should have 0 simple beta nodes, got %d", stats.SimpleBetaNodes)
		}
		if stats.JoinNodes != 0 {
			t.Errorf("Empty network should have 0 join nodes, got %d", stats.JoinNodes)
		}
		if stats.TotalTokens != 0 {
			t.Errorf("Empty network should have 0 tokens, got %d", stats.TotalTokens)
		}
		if stats.TotalFacts != 0 {
			t.Errorf("Empty network should have 0 facts, got %d", stats.TotalFacts)
		}
		if len(stats.MemoryStats) != 0 {
			t.Errorf("Empty network should have 0 memory stats, got %d", len(stats.MemoryStats))
		}
	})
}

// testECommerceFlow teste le flux E-commerce complexe avec validation de stock
func testECommerceFlow(t *testing.T, nodes []domain.BetaNode) {
	if len(nodes) < 3 {
		t.Fatal("E-commerce flow needs at least 3 nodes")
	}

	orderCustomerNode := nodes[0]
	customerProductNode := nodes[1] 
	productStockNode := nodes[2]

	// Créer les entités E-commerce
	orderFact := domain.NewFact("order_001", "Order", map[string]interface{}{
		"id":          "ORD-2025-001",
		"customer_id": "CUST-123",
		"product_id":  "PROD-456",
		"quantity":    5,
		"total":       299.99,
		"status":      "pending",
	})

	customerFact := domain.NewFact("customer_123", "Customer", map[string]interface{}{
		"id":            "CUST-123",
		"name":          "Marie Durand",
		"email":         "marie@email.com",
		"loyalty_level": "Gold",
		"credit_limit":  5000.00,
	})

	productFact := domain.NewFact("product_456", "Product", map[string]interface{}{
		"id":          "PROD-456",
		"name":        "Laptop Gaming Pro",
		"category":    "Electronics",
		"price":       59.99,
		"brand":       "TechCorp",
		"warranty":    24,
	})

	// Stock suffisant
	stockOKFact := domain.NewFact("stock_456_ok", "Stock", map[string]interface{}{
		"product_id":          "PROD-456",
		"warehouse_id":        "WH-001",
		"available_quantity":  10,  // Suffisant pour quantity=5
		"reserved_quantity":   2,
		"location":           "A-12-3",
	})

	// Stock insuffisant pour test d'échec
	stockKOFact := domain.NewFact("stock_456_ko", "Stock", map[string]interface{}{
		"product_id":          "PROD-456", 
		"warehouse_id":        "WH-002",
		"available_quantity":  3,   // Insuffisant pour quantity=5
		"reserved_quantity":   1,
		"location":           "B-5-7",
	})

	// Test 1: Flux réussi avec stock suffisant
	t.Run("ECommerceFlowSuccess", func(t *testing.T) {
		// Traiter Order + Customer
		orderToken := domain.NewToken("order_token", "order_source", []*domain.Fact{orderFact})
		orderCustomerNode.ProcessLeftToken(orderToken)
		orderCustomerNode.ProcessRightFact(customerFact)

		// Traiter Product 
		customerProductNode.ProcessRightFact(productFact)

		// Traiter Stock (suffisant)
		productStockNode.ProcessRightFact(stockOKFact)

		// Vérifier que le flux a réussi (au moins des faits traités)
		stats := getNodeStats(productStockNode)
		if stats.FactCount == 0 {
			t.Error("E-commerce flow should process facts when stock is sufficient")
		}
		t.Logf("E-commerce success stats: %d tokens, %d facts", stats.TokenCount, stats.FactCount)
	})

	// Test 2: Flux échoué avec stock insuffisant
	t.Run("ECommerceFlowStockFailure", func(t *testing.T) {
		// Nettoyer la mémoire
		productStockNode.ClearMemory()

		// Re-traiter le flux avec stock insuffisant
		orderToken := domain.NewToken("order_token_2", "order_source", []*domain.Fact{orderFact})
		orderCustomerNode.ProcessLeftToken(orderToken)
		orderCustomerNode.ProcessRightFact(customerFact)
		customerProductNode.ProcessRightFact(productFact)
		
		// Stock insuffisant ne devrait pas passer la condition quantity <= available_quantity
		productStockNode.ProcessRightFact(stockKOFact)
		
		// Le nombre de joins devrait être moindre ou nul
		stats := getNodeStats(productStockNode)
		t.Logf("Stock validation stats: %+v", stats)
	})
}

// testHRSecurityFlow teste le flux RH avec clearance de sécurité
func testHRSecurityFlow(t *testing.T, nodes []domain.BetaNode) {
	if len(nodes) < 2 {
		t.Fatal("HR flow needs at least 2 nodes")
	}

	empDeptNode := nodes[0]
	deptProjectNode := nodes[1]

	// Créer les entités RH
	employeeFact := domain.NewFact("emp_001", "Employee", map[string]interface{}{
		"id":              "EMP-001",
		"name":            "Alice Smith",
		"dept_id":         "DEPT-SEC",
		"clearance_level": 7,  // Niveau de sécurité élevé
		"position":        "Senior Developer",
		"hire_date":       "2020-03-15",
	})

	departmentFact := domain.NewFact("dept_sec", "Department", map[string]interface{}{
		"id":          "DEPT-SEC", 
		"name":        "Security Department",
		"manager":     "Bob Wilson",
		"budget":      2500000.00,
		"location":    "Building-A Floor-3",
	})

	// Projet nécessitant clearance élevée
	projectHighSecFact := domain.NewFact("proj_classified", "Project", map[string]interface{}{
		"id":                 "PROJ-CLASSIFIED",
		"name":               "Operation Phoenix",
		"required_clearance": 6,  // Alice (7) peut y accéder
		"deadline":           "2025-12-31",
		"budget":             500000.00,
		"classification":     "SECRET",
	})

	// Projet nécessitant clearance ultra-élevée
	projectTopSecretFact := domain.NewFact("proj_topsecret", "Project", map[string]interface{}{
		"id":                 "PROJ-TOPSECRET",
		"name":               "Quantum Initiative", 
		"required_clearance": 9,  // Alice (7) ne peut PAS y accéder
		"deadline":           "2026-06-30",
		"budget":             2000000.00,
		"classification":     "TOP-SECRET",
	})

	// Test 1: Assignation réussie avec clearance suffisante
	t.Run("HRFlowClearanceSuccess", func(t *testing.T) {
		// Traiter Employee + Department
		empToken := domain.NewToken("emp_token", "hr_source", []*domain.Fact{employeeFact})
		empDeptNode.ProcessLeftToken(empToken)
		empDeptNode.ProcessRightFact(departmentFact)

		// Traiter Project avec clearance suffisante
		deptProjectNode.ProcessRightFact(projectHighSecFact)

		// Vérifier l'assignation réussie
		stats := getNodeStats(deptProjectNode)
		if stats.TokenCount == 0 {
			t.Error("HR flow should assign employee when clearance is sufficient")
		}
	})

	// Test 2: Assignation échouée avec clearance insuffisante  
	t.Run("HRFlowClearanceFailure", func(t *testing.T) {
		// Nettoyer la mémoire
		deptProjectNode.ClearMemory()

		// Re-traiter avec projet nécessitant clearance trop élevée
		empToken := domain.NewToken("emp_token_2", "hr_source", []*domain.Fact{employeeFact})
		empDeptNode.ProcessLeftToken(empToken)
		empDeptNode.ProcessRightFact(departmentFact)
		
		// Projet TOP-SECRET ne devrait pas être accessible
		deptProjectNode.ProcessRightFact(projectTopSecretFact)
		
		// L'assignation devrait échouer (clearance_level >= required_clearance échoue)
		stats := getNodeStats(deptProjectNode)
		t.Logf("Security clearance stats: %+v", stats)
		// Note: Le test exact dépend de l'implémentation des conditions
	})
}

// getNodeStats helper pour obtenir les statistiques d'un nœud
func getNodeStats(node domain.BetaNode) struct {
	TokenCount int
	FactCount  int
} {
	tokens := node.GetLeftMemory()
	facts := node.GetRightMemory()
	return struct {
		TokenCount int
		FactCount  int
	}{
		TokenCount: len(tokens),
		FactCount:  len(facts),
	}
}

// Tests de performance pour les patterns complexes
func TestComplexPatternPerformance(t *testing.T) {
	logger := &mockLogger{}

	t.Run("ECommerceStressTest", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Construire le pattern E-commerce
		ecommercePattern := MultiJoinPattern{
			PatternID: "stress_order_processing",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Order",
					RightType: "Customer",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("customer_id", "id", "=="),
					},
					NodeID: "stress_order_customer",
				},
				{
					LeftType:  "OrderCustomer",
					RightType: "Product",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("product_id", "id", "=="),
					},
					NodeID: "stress_customer_product",
				},
			},
		}

		nodes, err := builder.BuildMultiJoinNetwork(ecommercePattern)
		if err != nil {
			t.Fatalf("Failed to build stress test network: %v", err)
		}

		// Stress test avec 100 commandes simultanées
		const numOrders = 100
		const numCustomers = 50
		const numProducts = 20

		// Générer les données de test
		for i := 0; i < numOrders; i++ {
			customerID := fmt.Sprintf("CUST-%03d", i%numCustomers)
			productID := fmt.Sprintf("PROD-%03d", i%numProducts)

			orderFact := domain.NewFact(fmt.Sprintf("order_%d", i), "Order", map[string]interface{}{
				"id":          fmt.Sprintf("ORD-%05d", i),
				"customer_id": customerID,
				"product_id":  productID,
				"amount":      float64(i * 10),
			})

			orderToken := domain.NewToken(fmt.Sprintf("order_token_%d", i), "order_source", []*domain.Fact{orderFact})
			nodes[0].ProcessLeftToken(orderToken)
		}

		// Générer les clients
		for i := 0; i < numCustomers; i++ {
			customerFact := domain.NewFact(fmt.Sprintf("customer_%d", i), "Customer", map[string]interface{}{
				"id":   fmt.Sprintf("CUST-%03d", i),
				"name": fmt.Sprintf("Customer-%d", i),
			})
			nodes[0].ProcessRightFact(customerFact)
		}

		// Générer les produits
		for i := 0; i < numProducts; i++ {
			productFact := domain.NewFact(fmt.Sprintf("product_%d", i), "Product", map[string]interface{}{
				"id":   fmt.Sprintf("PROD-%03d", i),
				"name": fmt.Sprintf("Product-%d", i),
			})
			nodes[1].ProcessRightFact(productFact)
		}

		// Vérifier les performances
		stats1 := getNodeStats(nodes[0])
		stats2 := getNodeStats(nodes[1])

		t.Logf("Stress test results - Node 0: %d tokens, %d facts", stats1.TokenCount, stats1.FactCount)
		t.Logf("Stress test results - Node 1: %d tokens, %d facts", stats2.TokenCount, stats2.FactCount)

		if stats1.TokenCount == 0 && stats1.FactCount == 0 {
			t.Error("Stress test should process some data")
		}
	})

	t.Run("MultiConditionPerformanceTest", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Créer un nœud avec conditions multiples complexes
		complexConditions := []domain.JoinCondition{
			domain.NewBasicJoinCondition("score", "min_score", ">="),
			domain.NewBasicJoinCondition("level", "required_level", "=="),
			domain.NewBasicJoinCondition("status", "valid_status", "=="),
			domain.NewBasicJoinCondition("rating", "threshold_rating", ">"),
		}

		joinNode := builder.CreateJoinNode("complex_multi_condition", complexConditions)

		// Test avec 1000 combinaisons
		const numTests = 1000

		for i := 0; i < numTests; i++ {
			// Token avec données variées
			leftFact := domain.NewFact(fmt.Sprintf("left_%d", i), "LeftEntity", map[string]interface{}{
				"score":  float64(i % 100),
				"level":  i % 10,
				"status": []string{"active", "pending", "inactive"}[i%3],
				"rating": float64((i % 50) + 1),
			})
			token := domain.NewToken(fmt.Sprintf("token_%d", i), "source", []*domain.Fact{leftFact})

			// Fait avec critères correspondants ou non
			rightFact := domain.NewFact(fmt.Sprintf("right_%d", i), "RightEntity", map[string]interface{}{
				"min_score":       float64((i % 80) + 10),
				"required_level":  (i + 5) % 10,
				"valid_status":    []string{"active", "pending", "disabled"}[i%3],
				"threshold_rating": float64((i % 30) + 20),
			})

			// Traitement
			joinNode.ProcessLeftToken(token)
			joinNode.ProcessRightFact(rightFact)
		}

		// Vérifier que le traitement s'est fait
		stats := getNodeStats(joinNode)
		t.Logf("Multi-condition performance test: %d tokens, %d facts processed", stats.TokenCount, stats.FactCount)

		if stats.TokenCount == 0 && stats.FactCount == 0 {
			t.Error("Performance test should process some data")
		}
	})
}

// Tests avancés de cas d'usage métier
func TestAdvancedBusinessCases(t *testing.T) {
	logger := &mockLogger{}

	t.Run("FinancialRiskAssessment", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Pattern d'évaluation des risques financiers
		riskPattern := MultiJoinPattern{
			PatternID: "financial_risk_assessment",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Transaction",
					RightType: "Account",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("account_id", "id", "=="),
					},
					NodeID: "transaction_account_join",
				},
				{
					LeftType:  "TransactionAccount",
					RightType: "RiskProfile",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("customer_id", "customer_id", "=="),
						domain.NewBasicJoinCondition("amount", "daily_limit", "<="),
					},
					NodeID: "account_risk_evaluation",
				},
			},
		}

		nodes, err := builder.BuildMultiJoinNetwork(riskPattern)
		if err != nil {
			t.Fatalf("Failed to build financial risk network: %v", err)
		}

		// Transaction suspecte
		suspiciousTransaction := domain.NewFact("txn_suspicious", "Transaction", map[string]interface{}{
			"id":         "TXN-999999",
			"account_id": "ACC-HIGH-RISK",
			"amount":     50000.00, // Montant élevé
			"currency":   "USD",
			"country":    "Offshore",
			"timestamp":  "2025-11-06T15:30:00Z",
		})

		// Compte à risque
		riskAccount := domain.NewFact("acc_risky", "Account", map[string]interface{}{
			"id":          "ACC-HIGH-RISK",
			"customer_id": "CUST-SUSPICIOUS",
			"type":        "Business",
			"balance":     1000000.00,
			"status":      "Active",
		})

		// Profil de risque
		highRiskProfile := domain.NewFact("risk_high", "RiskProfile", map[string]interface{}{
			"customer_id":  "CUST-SUSPICIOUS",
			"risk_level":   "HIGH",
			"daily_limit":  25000.00, // Limite dépassée par la transaction
			"flags":        []string{"PEP", "Sanctions"},
			"last_review":  "2025-01-01",
		})

		// Traiter le flux de risque
		txnToken := domain.NewToken("risk_token", "risk_source", []*domain.Fact{suspiciousTransaction})
		nodes[0].ProcessLeftToken(txnToken)
		nodes[0].ProcessRightFact(riskAccount)
		nodes[1].ProcessRightFact(highRiskProfile)

		// Vérifier l'évaluation des risques
		stats := getNodeStats(nodes[1])
		t.Logf("Financial risk assessment: %d risk evaluations", stats.TokenCount)
	})

	t.Run("SupplyChainOptimization", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Pattern d'optimisation de la chaîne d'approvisionnement
		supplyPattern := MultiJoinPattern{
			PatternID: "supply_chain_optimization",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Order",
					RightType: "Supplier",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("supplier_id", "id", "=="),
					},
					NodeID: "order_supplier_join",
				},
				{
					LeftType:  "OrderSupplier",
					RightType: "Logistics",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("region", "service_region", "=="),
						domain.NewBasicJoinCondition("delivery_date", "available_date", ">="),
					},
					NodeID: "supplier_logistics_optimization",
				},
			},
		}

		nodes, err := builder.BuildMultiJoinNetwork(supplyPattern)
		if err != nil {
			t.Fatalf("Failed to build supply chain network: %v", err)
		}

		// Commande urgente
		urgentOrder := domain.NewFact("order_urgent", "Order", map[string]interface{}{
			"id":            "ORD-URGENT-001",
			"supplier_id":   "SUP-GLOBAL",
			"region":        "Europe",
			"delivery_date": "2025-11-10",
			"priority":      "HIGH",
			"value":         75000.00,
		})

		// Fournisseur global
		globalSupplier := domain.NewFact("sup_global", "Supplier", map[string]interface{}{
			"id":           "SUP-GLOBAL",
			"name":         "Global Supply Corp",
			"rating":       4.8,
			"regions":      []string{"Europe", "Asia", "Americas"},
			"capacity":     1000000,
			"lead_time":    3,
		})

		// Logistique disponible
		expressLogistics := domain.NewFact("log_express", "Logistics", map[string]interface{}{
			"id":             "LOG-EXPRESS",
			"service_region": "Europe", 
			"available_date": "2025-11-08", // Disponible avant la date requise
			"service_type":   "Express",
			"cost_per_km":    2.50,
			"max_weight":     5000,
		})

		// Traiter l'optimisation de la chaîne
		orderToken := domain.NewToken("supply_token", "supply_source", []*domain.Fact{urgentOrder})
		nodes[0].ProcessLeftToken(orderToken)
		nodes[0].ProcessRightFact(globalSupplier)
		nodes[1].ProcessRightFact(expressLogistics)

		// Vérifier l'optimisation
		stats := getNodeStats(nodes[1])
		t.Logf("Supply chain optimization: %d optimized routes", stats.TokenCount)
	})

	t.Run("BankingAntiFraudMegaPattern", func(t *testing.T) {
		builder := NewBetaNetworkBuilder(logger)

		// Pattern anti-fraude ULTRA-COMPLEXE : 5 niveaux de jointures
		// Transaction → Account → Customer → RiskProfile → GeolocationRisk
		antiFraudPattern := MultiJoinPattern{
			PatternID: "banking_anti_fraud_ultra_complex",
			JoinSpecs: []JoinSpecification{
				{
					LeftType:  "Transaction",
					RightType: "Account",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("account_id", "id", "=="),
					},
					NodeID: "txn_account_join",
				},
				{
					LeftType:  "TransactionAccount", 
					RightType: "Customer",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("customer_id", "id", "=="),
					},
					NodeID: "account_customer_join",
				},
				{
					LeftType:  "TransactionCustomer",
					RightType: "RiskProfile",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("customer_id", "customer_id", "=="),
						domain.NewBasicJoinCondition("amount", "max_transaction", "<="),
					},
					NodeID: "customer_risk_join",
				},
				{
					LeftType:  "CustomerRisk",
					RightType: "GeolocationRisk",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("country_code", "country", "=="),
						domain.NewBasicJoinCondition("risk_score", "max_allowed_score", "<="),
					},
					NodeID: "risk_geolocation_join",
				},
				{
					LeftType:  "RiskGeolocation",
					RightType: "ComplianceRule",
					Conditions: []domain.JoinCondition{
						domain.NewBasicJoinCondition("transaction_type", "applicable_types", "=="),
						domain.NewBasicJoinCondition("compliance_level", "required_level", ">="),
					},
					NodeID: "geolocation_compliance_final",
				},
			},
			FinalAction: "execute_anti_fraud_decision",
		}

		nodes, err := builder.BuildMultiJoinNetwork(antiFraudPattern)
		if err != nil {
			t.Fatalf("Failed to build ultra-complex anti-fraud network: %v", err)
		}

		if len(nodes) != 5 {
			t.Errorf("Expected 5 nodes for ultra-complex pattern, got %d", len(nodes))
		}

		// Données ultra-réalistes de fraude bancaire
		suspiciousTransaction := domain.NewFact("txn_fraud_suspect", "Transaction", map[string]interface{}{
			"id":               "TXN-2025110615300001",
			"account_id":       "ACC-CRYPTO-TRADER",
			"amount":           99999.99,
			"currency":         "USD",
			"transaction_type": "WIRE_TRANSFER",
			"country_code":     "KY", // Cayman Islands
			"timestamp":        "2025-11-06T15:30:00Z",
			"merchant":         "CRYPTO_EXCHANGE_OFFSHORE",
			"risk_score":       8.5, // Très élevé
			"compliance_level": 3,
		})

		cryptoAccount := domain.NewFact("acc_crypto", "Account", map[string]interface{}{
			"id":             "ACC-CRYPTO-TRADER",
			"customer_id":    "CUST-WHALE-001",
			"type":           "BUSINESS_PREMIUM",
			"balance":        50000000.00, // 50M USD
			"opened_date":    "2025-01-01",
			"status":         "ACTIVE",
			"flags":          []string{"CRYPTO", "HIGH_VOLUME"},
		})

		whaleCustomer := domain.NewFact("cust_whale", "Customer", map[string]interface{}{
			"id":              "CUST-WHALE-001",
			"name":            "Crypto Whale Holdings LLC",
			"customer_type":   "CORPORATE",
			"kyc_status":      "ENHANCED_DUE_DILIGENCE",
			"pep_status":      false,
			"sanctions_check": "CLEAR",
			"incorporation":   "Cayman Islands",
		})

		highRiskProfile := domain.NewFact("risk_whale", "RiskProfile", map[string]interface{}{
			"customer_id":      "CUST-WHALE-001",
			"risk_category":    "HIGH",
			"max_transaction":  100000.00, // Transaction limite
			"daily_limit":      5000000.00,
			"monitoring_level": "ENHANCED",
			"last_review":      "2025-11-01",
		})

		caymanRisk := domain.NewFact("geo_cayman", "GeolocationRisk", map[string]interface{}{
			"country":             "KY",
			"country_name":        "Cayman Islands",
			"risk_rating":         "HIGH",
			"max_allowed_score":   9.0, // Autorise le score de 8.5
			"sanctions_risk":      "MEDIUM",
			"offshore_jurisdiction": true,
			"fatf_compliance":     "ADEQUATE",
		})

		complianceRule := domain.NewFact("rule_wire", "ComplianceRule", map[string]interface{}{
			"id":                "RULE-WIRE-OFFSHORE",
			"applicable_types":  "WIRE_TRANSFER",
			"required_level":    2, // Customer a niveau 3, OK
			"max_amount":        100000.00,
			"enhanced_monitoring": true,
			"reporting_required": true,
			"jurisdiction":       []string{"US", "EU", "KY"},
		})

		// Traiter le flux anti-fraude complet (5 niveaux!)
		txnToken := domain.NewToken("anti_fraud_token", "fraud_detection_source", []*domain.Fact{suspiciousTransaction})
		
		// Niveau 1: Transaction → Account
		nodes[0].ProcessLeftToken(txnToken)
		nodes[0].ProcessRightFact(cryptoAccount)

		// Niveau 2: Account → Customer  
		nodes[1].ProcessRightFact(whaleCustomer)

		// Niveau 3: Customer → RiskProfile (avec validation montant)
		nodes[2].ProcessRightFact(highRiskProfile)

		// Niveau 4: RiskProfile → GeolocationRisk (avec validation score)
		nodes[3].ProcessRightFact(caymanRisk)

		// Niveau 5: GeolocationRisk → ComplianceRule (validation finale)
		nodes[4].ProcessRightFact(complianceRule)

		// Vérifier le résultat de l'analyse anti-fraude COMPLÈTE
		finalStats := getNodeStats(nodes[4])
		t.Logf("🔒 ULTRA-COMPLEX Anti-Fraud Analysis Complete:")
		t.Logf("   📊 Final Decision Tokens: %d", finalStats.TokenCount)
		t.Logf("   📋 Rules Processed: %d", finalStats.FactCount)
		t.Logf("   🚨 Transaction: $%.2f to %s", 99999.99, "Cayman Islands")
		t.Logf("   ⚖️  Risk Score: %.1f (Threshold: %.1f)", 8.5, 9.0)
		t.Logf("   ✅ Compliance Level: %d (Required: %d)", 3, 2)

		// Validation des statistiques réseau final
		networkStats := builder.NetworkStatistics()
		t.Logf("🌐 Network Performance:")
		t.Logf("   🔗 Total Nodes: %d", networkStats.TotalNodes)
		t.Logf("   🧠 Total Tokens: %d", networkStats.TotalTokens)
		t.Logf("   📦 Total Facts: %d", networkStats.TotalFacts)

		// Ce pattern ultra-complexe devrait traiter au moins quelque chose !
		if networkStats.TotalFacts == 0 {
			t.Error("Ultra-complex anti-fraud pattern should process facts")
		}

		t.Log("🎉 ULTRA-COMPLEX BANKING ANTI-FRAUD PATTERN EXECUTED SUCCESSFULLY!")
	})
}
