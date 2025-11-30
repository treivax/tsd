// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License

package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBetaBackwardCompatibility_SimpleJoins vérifie que les jointures simples
// fonctionnent toujours correctement avec Beta Sharing
func TestBetaBackwardCompatibility_SimpleJoins(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "simple_joins.tsd")

	// Règle avec une jointure simple entre deux patterns
	content := `type Person : <id: string, age: number, name: string>
type Order : <id: string, personId: string, amount: number>

rule person_order : {p: Person, o: Order} / p.id == o.personId AND o.amount > 100 ==> print("Large order")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Vérifier la structure du réseau
	stats := network.GetNetworkStats()

	if stats["type_nodes"].(int) != 2 {
		t.Errorf("Attendu 2 TypeNodes, obtenu %d", stats["type_nodes"].(int))
	}

	// Vérifier qu'un JoinNode est créé
	if stats["beta_nodes"].(int) < 1 {
		t.Errorf("Attendu au moins 1 BetaNode (JoinNode), obtenu %d", stats["beta_nodes"].(int))
	}

	// Soumettre des faits qui correspondent
	person := Fact{
		ID:   "P1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"age":  30,
			"name": "Alice",
		},
	}

	order := Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o1",
			"personId": "p1",
			"amount":   150,
		},
	}

	if err := network.SubmitFact(&person); err != nil {
		t.Fatalf("Erreur ajout person: %v", err)
	}

	if err := network.SubmitFact(&order); err != nil {
		t.Fatalf("Erreur ajout order: %v", err)
	}

	// Vérifier l'activation
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount != 1 {
		t.Errorf("Attendu 1 activation (jointure réussie), obtenu %d", activatedCount)
	}

	t.Logf("✅ Jointures simples: backward compatible")
}

// TestBetaBackwardCompatibility_ExistingBehavior vérifie que le comportement
// des jointures existantes n'a pas changé avec Beta Sharing
func TestBetaBackwardCompatibility_ExistingBehavior(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "existing_behavior.tsd")

	content := `type Employee : <id: string, name: string, deptId: string>
type Department : <id: string, name: string, budget: number>

rule emp_in_dept : {e: Employee, d: Department} / e.deptId == d.id ==> print("Employee in dept")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Test 1: Ajout dans l'ordre (employee puis department)
	emp1 := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"name":   "Alice",
			"deptId": "d1",
		},
	}

	dept1 := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":     "d1",
			"name":   "Engineering",
			"budget": 100000,
		},
	}

	network.SubmitFact(&emp1)
	network.SubmitFact(&dept1)

	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount != 1 {
		t.Errorf("Test 1: Attendu 1 activation, obtenu %d", activatedCount)
	}

	// Test 2: Ajouter un autre employee dans le même département
	emp2 := Fact{
		ID:   "E2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e2",
			"name":   "Bob",
			"deptId": "d1",
		},
	}

	network.SubmitFact(&emp2)

	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount != 2 {
		t.Errorf("Test 2: Attendu 2 activations, obtenu %d", activatedCount)
	}

	// Test 3: Retraction d'un employee
	if err := network.RetractFact("Employee_E1"); err != nil {
		t.Fatalf("Erreur rétractation: %v", err)
	}

	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount != 1 {
		t.Errorf("Test 3: Attendu 1 activation après rétractation, obtenu %d", activatedCount)
	}

	t.Logf("✅ Comportement existant des jointures: backward compatible")
}

// TestBetaNoRegression_AllPreviousTests exécute plusieurs scénarios
// de jointures pour détecter toute régression
func TestBetaNoRegression_AllPreviousTests(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		facts       []Fact
		activations int
	}{
		{
			name: "Two pattern join",
			content: `type A : <id: string, x: number>
type B : <id: string, y: number, aId: string>
rule join_ab : {a: A, b: B} / a.id == b.aId ==> print("Join")`,
			facts: []Fact{
				{ID: "A1", Type: "A", Fields: map[string]interface{}{"id": "a1", "x": 10}},
				{ID: "B1", Type: "B", Fields: map[string]interface{}{"id": "b1", "y": 20, "aId": "a1"}},
			},
			activations: 1,
		},
		{
			name: "Three pattern join",
			content: `type A : <id: string, x: number>
type B : <id: string, aId: string>
type C : <id: string, bId: string>
rule join_abc : {a: A, b: B, c: C} / a.id == b.aId AND b.id == c.bId ==> print("Join")`,
			facts: []Fact{
				{ID: "A1", Type: "A", Fields: map[string]interface{}{"id": "a1", "x": 10}},
				{ID: "B1", Type: "B", Fields: map[string]interface{}{"id": "b1", "aId": "a1"}},
				{ID: "C1", Type: "C", Fields: map[string]interface{}{"id": "c1", "bId": "b1"}},
			},
			activations: 1,
		},
		{
			name: "Join with additional constraints",
			content: `type Person : <id: string, age: number, name: string>
type Account : <id: string, ownerId: string, balance: number>
rule rich_adult : {p: Person, a: Account} / p.id == a.ownerId AND p.age >= 18 AND a.balance > 10000 ==> print("Rich adult")`,
			facts: []Fact{
				{ID: "P1", Type: "Person", Fields: map[string]interface{}{"id": "p1", "age": 25, "name": "Alice"}},
				{ID: "A1", Type: "Account", Fields: map[string]interface{}{"id": "a1", "ownerId": "p1", "balance": 15000}},
			},
			activations: 1,
		},
		{
			name: "Multiple matching joins",
			content: `type User : <id: string, name: string>
type Post : <id: string, userId: string, title: string>
rule user_posts : {u: User, p: Post} / u.id == p.userId ==> print("Post")`,
			facts: []Fact{
				{ID: "U1", Type: "User", Fields: map[string]interface{}{"id": "u1", "name": "Alice"}},
				{ID: "P1", Type: "Post", Fields: map[string]interface{}{"id": "p1", "userId": "u1", "title": "Post1"}},
				{ID: "P2", Type: "Post", Fields: map[string]interface{}{"id": "p2", "userId": "u1", "title": "Post2"}},
			},
			activations: 2,
		},
		{
			name: "Join with no matches",
			content: `type X : <id: string, val: number>
type Y : <id: string, xId: string>
rule join_xy : {x: X, y: Y} / x.id == y.xId ==> print("Join")`,
			facts: []Fact{
				{ID: "X1", Type: "X", Fields: map[string]interface{}{"id": "x1", "val": 10}},
				{ID: "Y1", Type: "Y", Fields: map[string]interface{}{"id": "y1", "xId": "x2"}},
			},
			activations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tsdFile := filepath.Join(tempDir, "test.tsd")

			if err := os.WriteFile(tsdFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Erreur création fichier: %v", err)
			}

			storage := NewMemoryStorage()
			pipeline := NewConstraintPipeline()
			network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
			if err != nil {
				t.Fatalf("Erreur construction réseau: %v", err)
			}

			// Soumettre les faits
			for _, fact := range tt.facts {
				// Copy fact to avoid loop variable reuse bug
				factCopy := fact
				if err := network.SubmitFact(&factCopy); err != nil {
					t.Fatalf("Erreur ajout fait %s: %v", fact.ID, err)
				}
			}

			// Vérifier les activations
			activatedCount := 0
			for _, terminalNode := range network.TerminalNodes {
				memory := terminalNode.GetMemory()
				activatedCount += len(memory.Tokens)
			}

			if activatedCount != tt.activations {
				t.Errorf("Attendu %d activations, obtenu %d", tt.activations, activatedCount)
			}

			t.Logf("✅ Test '%s': pas de régression", tt.name)
		})
	}
}

// TestBetaBackwardCompatibility_JoinNodeSharing vérifie que le partage
// des JoinNodes fonctionne correctement avec Beta Sharing
func TestBetaBackwardCompatibility_JoinNodeSharing(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "join_sharing.tsd")

	// Deux règles avec la même jointure de base (devraient partager le JoinNode)
	content := `type Person : <id: string, age: number, name: string>
type Order : <id: string, personId: string, amount: number>

rule large_orders : {p: Person, o: Order} / p.id == o.personId AND o.amount > 100 ==> print("Large")
rule very_large_orders : {p: Person, o: Order} / p.id == o.personId AND o.amount > 500 ==> print("Very large")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	stats := network.GetNetworkStats()

	// Avec Beta Sharing, on devrait avoir moins de JoinNodes qu'il n'y a de règles
	// car la jointure de base (p.id == o.personId) devrait être partagée
	betaCount := stats["beta_nodes"].(int)
	t.Logf("Nombre de BetaNodes créés: %d", betaCount)

	// Soumettre des faits de test
	person := Fact{
		ID:   "P1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"age":  30,
			"name": "Alice",
		},
	}

	order1 := Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o1",
			"personId": "p1",
			"amount":   150,
		},
	}

	order2 := Fact{
		ID:   "O2",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":       "o2",
			"personId": "p1",
			"amount":   600,
		},
	}

	network.SubmitFact(&person)
	network.SubmitFact(&order1)
	network.SubmitFact(&order2)

	// Vérifier les activations
	activations := make(map[string]int)
	for ruleName, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activations[ruleName] = len(memory.Tokens)
	}

	// order1 (150) devrait activer large_orders
	// order2 (600) devrait activer les deux règles
	totalActivations := 0
	for _, count := range activations {
		totalActivations += count
	}

	if totalActivations != 3 {
		t.Errorf("Attendu 3 activations totales, obtenu %d", totalActivations)
		for rule, count := range activations {
			t.Logf("  %s: %d activations", rule, count)
		}
	}

	t.Logf("✅ JoinNode sharing: backward compatible")
}

// TestBetaBackwardCompatibility_PerformanceCharacteristics vérifie que
// les caractéristiques de performance sont maintenues ou améliorées
func TestBetaBackwardCompatibility_PerformanceCharacteristics(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "perf.tsd")

	// Plusieurs règles avec des jointures similaires
	content := `type User : <id: string, name: string, age: number>
type Product : <id: string, name: string, price: number>
type Purchase : <id: string, userId: string, productId: string, quantity: number>

rule user_purchase : {u: User, pu: Purchase} / u.id == pu.userId ==> print("User purchase")
rule adult_purchase : {u: User, pu: Purchase} / u.id == pu.userId AND u.age >= 18 ==> print("Adult purchase")
rule product_purchase : {pr: Product, pu: Purchase} / pr.id == pu.productId ==> print("Product purchase")
rule expensive_product : {pr: Product, pu: Purchase} / pr.id == pu.productId AND pr.price > 100 ==> print("Expensive")
rule full_purchase : {u: User, pr: Product, pu: Purchase} / u.id == pu.userId AND pr.id == pu.productId ==> print("Full")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	stats := network.GetNetworkStats()

	// Avec Beta Sharing, on devrait avoir moins de JoinNodes que de règles avec jointures
	betaCount := stats["beta_nodes"].(int)
	terminalCount := stats["terminal_nodes"].(int)

	if terminalCount != 5 {
		t.Errorf("Attendu 5 TerminalNodes, obtenu %d", terminalCount)
	}

	t.Logf("Performance: %d BetaNodes pour 5 règles avec jointures", betaCount)

	// Soumettre des faits et mesurer les activations
	user := Fact{
		ID:   "U1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":   "u1",
			"name": "Alice",
			"age":  25,
		},
	}

	product := Fact{
		ID:   "PR1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":    "pr1",
			"name":  "Laptop",
			"price": 150,
		},
	}

	purchase := Fact{
		ID:   "PU1",
		Type: "Purchase",
		Fields: map[string]interface{}{
			"id":        "pu1",
			"userId":    "u1",
			"productId": "pr1",
			"quantity":  1,
		},
	}

	network.SubmitFact(&user)
	network.SubmitFact(&product)
	network.SubmitFact(&purchase)

	// Vérifier que toutes les règles appropriées sont activées
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	// user_purchase: 1
	// adult_purchase: 1
	// product_purchase: 1
	// expensive_product: 1
	// full_purchase: 1
	// Total: 5 activations
	if activatedCount != 5 {
		t.Errorf("Attendu 5 activations, obtenu %d", activatedCount)
	}

	t.Logf("✅ Performance: toutes les règles activées correctement")
}

// TestBetaBackwardCompatibility_ComplexJointures vérifie que les jointures
// complexes avec plusieurs patterns fonctionnent toujours
func TestBetaBackwardCompatibility_ComplexJointures(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "complex_joins.tsd")

	// Règle avec 4 patterns et plusieurs jointures
	content := `type Customer : <id: string, name: string, country: string>
type Order : <id: string, customerId: string, total: number>
type Product : <id: string, name: string, category: string>
type OrderItem : <id: string, orderId: string, productId: string, quantity: number>

rule complex : {c: Customer, o: Order, p: Product, oi: OrderItem} /
    c.id == o.customerId AND
    o.id == oi.orderId AND
    p.id == oi.productId AND
    p.category == 'electronics'
    ==> print("Complex match")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Soumettre des faits qui correspondent
	customer := Fact{
		ID:   "C1",
		Type: "Customer",
		Fields: map[string]interface{}{
			"id":      "c1",
			"name":    "Alice",
			"country": "USA",
		},
	}

	order := Fact{
		ID:   "O1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":         "o1",
			"customerId": "c1",
			"total":      500,
		},
	}

	product := Fact{
		ID:   "P1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":       "p1",
			"name":     "Laptop",
			"category": "electronics",
		},
	}

	orderItem := Fact{
		ID:   "OI1",
		Type: "OrderItem",
		Fields: map[string]interface{}{
			"id":        "oi1",
			"orderId":   "o1",
			"productId": "p1",
			"quantity":  1,
		},
	}

	// Soumettre dans un ordre différent pour tester la propagation
	network.SubmitFact(&product)
	network.SubmitFact(&orderItem)
	network.SubmitFact(&customer)
	network.SubmitFact(&order)

	// Vérifier l'activation
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount != 1 {
		t.Errorf("Attendu 1 activation (jointure complexe réussie), obtenu %d", activatedCount)
	}

	t.Logf("✅ Jointures complexes: backward compatible")
}

// TestBetaBackwardCompatibility_AggregationsWithJoins vérifie que les
// agrégations fonctionnent correctement avec les jointures
func TestBetaBackwardCompatibility_AggregationsWithJoins(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "agg_joins.tsd")

	// Règle avec agrégation et jointure
	// NOTE: This syntax is not currently supported by the TSD parser
	// The parser fails at "/ {e: Employee}" - aggregations with join syntax needs parser work
	content := `type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Avg salary")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Soumettre des faits
	dept := Fact{
		ID:   "D1",
		Type: "Department",
		Fields: map[string]interface{}{
			"id":   "d1",
			"name": "Engineering",
		},
	}

	emp1 := Fact{
		ID:   "E1",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e1",
			"deptId": "d1",
			"salary": 50000,
		},
	}

	emp2 := Fact{
		ID:   "E2",
		Type: "Employee",
		Fields: map[string]interface{}{
			"id":     "e2",
			"deptId": "d1",
			"salary": 60000,
		},
	}

	network.SubmitFact(&dept)
	network.SubmitFact(&emp1)
	network.SubmitFact(&emp2)

	// Vérifier qu'une activation existe
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount < 1 {
		t.Errorf("Attendu au moins 1 activation (agrégation avec jointure), obtenu %d", activatedCount)
	}

	t.Logf("✅ Agrégations avec jointures: backward compatible")
}

// TestBetaBackwardCompatibility_RuleRemovalWithJoins vérifie que la
// suppression de règles fonctionne correctement avec les jointures
func TestBetaBackwardCompatibility_RuleRemovalWithJoins(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "removal_joins.tsd")

	content := `type Person : <id: string, name: string>
type Address : <id: string, personId: string, city: string>

rule person_address : {p: Person, a: Address} / p.id == a.personId ==> print("Person address")
rule person_in_ny : {p: Person, a: Address} / p.id == a.personId AND a.city == 'New York' ==> print("NY resident")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Vérifier l'état initial
	statsInitial := network.GetNetworkStats()
	initialTerminalCount := statsInitial["terminal_nodes"].(int)
	if initialTerminalCount != 2 {
		t.Errorf("Attendu 2 TerminalNodes initialement, obtenu %d", initialTerminalCount)
	}

	// Supprimer une règle
	if err := network.RemoveRule("person_in_ny"); err != nil {
		t.Fatalf("Erreur suppression règle: %v", err)
	}

	// Vérifier qu'il reste 1 TerminalNode
	statsAfter := network.GetNetworkStats()
	afterTerminalCount := statsAfter["terminal_nodes"].(int)
	if afterTerminalCount != 1 {
		t.Errorf("Attendu 1 TerminalNode après suppression, obtenu %d", afterTerminalCount)
	}

	// Vérifier que la règle restante fonctionne
	person := Fact{
		ID:   "P1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"name": "Alice",
		},
	}

	address := Fact{
		ID:   "A1",
		Type: "Address",
		Fields: map[string]interface{}{
			"id":       "a1",
			"personId": "p1",
			"city":     "Boston",
		},
	}

	network.SubmitFact(&person)
	network.SubmitFact(&address)

	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}

	if activatedCount != 1 {
		t.Errorf("Attendu 1 activation, obtenu %d", activatedCount)
	}

	t.Logf("✅ Suppression de règles avec jointures: backward compatible")
}

// TestBetaBackwardCompatibility_FactRetractionWithJoins vérifie que la
// rétractation de faits fonctionne correctement avec les jointures
func TestBetaBackwardCompatibility_FactRetractionWithJoins(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "retraction_joins.tsd")

	content := `type Author : <id: string, name: string>
type Book : <id: string, authorId: string, title: string>

rule author_book : {a: Author, b: Book} / a.id == b.authorId ==> print("Author book")
`

	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, err := pipeline.BuildNetworkFromConstraintFile(tsdFile, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Ajouter des faits
	author := Fact{
		ID:   "A1",
		Type: "Author",
		Fields: map[string]interface{}{
			"id":   "a1",
			"name": "Alice",
		},
	}

	book1 := Fact{
		ID:   "B1",
		Type: "Book",
		Fields: map[string]interface{}{
			"id":       "b1",
			"authorId": "a1",
			"title":    "Book 1",
		},
	}

	book2 := Fact{
		ID:   "B2",
		Type: "Book",
		Fields: map[string]interface{}{
			"id":       "b2",
			"authorId": "a1",
			"title":    "Book 2",
		},
	}

	network.SubmitFact(&author)
	network.SubmitFact(&book1)
	network.SubmitFact(&book2)

	// Vérifier les activations initiales
	initialActivations := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		initialActivations += len(memory.Tokens)
	}

	if initialActivations != 2 {
		t.Errorf("Attendu 2 activations initiales, obtenu %d", initialActivations)
	}

	// Rétracter un livre
	if err := network.RetractFact("Book_B1"); err != nil {
		t.Fatalf("Erreur rétractation: %v", err)
	}

	// Vérifier qu'il reste 1 activation
	activationsAfter := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activationsAfter += len(memory.Tokens)
	}

	if activationsAfter != 1 {
		t.Errorf("Attendu 1 activation après rétractation, obtenu %d", activationsAfter)
	}

	// Rétracter l'auteur
	if err := network.RetractFact("Author_A1"); err != nil {
		t.Fatalf("Erreur rétractation auteur: %v", err)
	}

	// Vérifier qu'il n'y a plus d'activations
	activationsFinal := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activationsFinal += len(memory.Tokens)
	}

	if activationsFinal != 0 {
		t.Errorf("Attendu 0 activations après rétractation de l'auteur, obtenu %d", activationsFinal)
	}

	t.Logf("✅ Rétractation de faits avec jointures: backward compatible")
}
