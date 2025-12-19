# Prompt 08 - Tests d'Intégration et End-to-End

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Créer et migrer les tests d'intégration et end-to-end pour valider le système complet avec la nouvelle gestion des identifiants :

1. **Tests d'intégration** - Interaction entre modules
2. **Tests end-to-end** - Scénarios utilisateur complets
3. **Tests de performance** - Benchmarks si nécessaire
4. **Tests de non-régression** - Garantir compatibilité fonctionnelle
5. **Exemples TSD** - Fichiers `.tsd` de démonstration

---

## 📋 Contexte

### État Actuel

Les tests existants :
- Tests unitaires isolés par module
- Quelques tests d'intégration basiques
- Tests e2e dans `tests/e2e/`
- Exemples TSD dans `examples/`

### État Cible

Les tests doivent couvrir :
- **Intégration complète** : Parser → Validation → RETE → API
- **Scénarios réels** : Programmes TSD complets avec affectations
- **Performance** : Pas de dégradation
- **Exemples** : Démonstrations des nouvelles fonctionnalités

---

## 📝 Tâches à Réaliser

### 1. Analyser les Tests d'Intégration Existants

#### Rechercher les Tests

```bash
# Tests d'intégration
find tests/integration -name "*.go" -type f 2>/dev/null || echo "Pas de tests/integration"
find . -name "*integration*test.go" -type f

# Tests e2e
find tests/e2e -name "*.go" -o -name "*.tsd" -type f 2>/dev/null

# Exemples TSD
find examples/ -name "*.tsd" -type f
```

**Créer inventaire** : `REPORTS/new_ids_integration_tests_inventory.md`

```markdown
# Inventaire Tests d'Intégration - Migration IDs

## Tests d'Intégration Existants

### constraint/
- integration_test.go : XX scénarios
- [...]

### rete/
- [...]

### tests/integration/
- [...]

## Tests End-to-End

### tests/e2e/
- Fichiers .tsd : XX
- Tests Go : XX
- [...]

## Exemples

### examples/
- basic.tsd
- primary_keys.tsd
- [...]

## Action Requise

- Tests à migrer : XX
- Nouveaux tests : XX
- Exemples à créer : XX
```

### 2. Créer Tests d'Intégration Complète

#### Fichier : `tests/integration/fact_lifecycle_test.go` (nouveau)

**Cycle de vie complet d'un fait** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package integration

import (
	"testing"
	"github.com/resinsec/tsd/constraint"
	"github.com/resinsec/tsd/rete"
)

func TestFactLifecycle_Complete(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - CYCLE DE VIE COMPLET")
	t.Log("===========================================")
	
	// 1. Parser un programme TSD
	input := `
		type User(#name: string, age: number)
		type Login(user: User, #email: string, password: string)
		
		alice = User("Alice", 30)
		bob = User("Bob", 25)
		
		Login(alice, "alice@example.com", "secret1")
		Login(bob, "bob@example.com", "secret2")
		
		{u: User, l: Login} / l.user == u && u.age > 25 ==> 
			Log("Senior user: " + u.name + " - " + l.email)
	`
	
	t.Log("📝 Étape 1: Parsing...")
	program, err := constraint.ParseProgram(input)
	if err != nil {
		t.Fatalf("❌ Erreur de parsing: %v", err)
	}
	t.Logf("✅ Programme parsé: %d types, %d affectations, %d faits, %d règles",
		len(program.Types), len(program.FactAssignments), len(program.Facts), len(program.Expressions))
	
	// 2. Valider le programme
	t.Log("📝 Étape 2: Validation...")
	validator := constraint.NewProgramValidator()
	if err := validator.Validate(*program); err != nil {
		t.Fatalf("❌ Erreur de validation: %v", err)
	}
	t.Log("✅ Programme validé")
	
	// 3. Convertir en format RETE
	t.Log("📝 Étape 3: Conversion RETE...")
	reteFacts, err := constraint.ConvertFactsToReteFormat(*program)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}
	t.Logf("✅ %d faits convertis", len(reteFacts))
	
	// Vérifier les IDs générés
	for i, fact := range reteFacts {
		id, ok := fact[constraint.FieldNameInternalID].(string)
		if !ok || id == "" {
			t.Errorf("❌ Fait %d: ID manquant ou invalide", i)
		}
		t.Logf("   - Fait %d: %s (ID: %s)", i, fact[constraint.FieldNameReteType], id)
	}
	
	// 4. Créer le réseau RETE
	t.Log("📝 Étape 4: Création réseau RETE...")
	network := rete.NewNetwork()
	
	// Compiler les règles
	for _, expr := range program.Expressions {
		if err := network.CompileExpression(expr, program.Types); err != nil {
			t.Fatalf("❌ Erreur de compilation règle: %v", err)
		}
	}
	t.Log("✅ Règles compilées")
	
	// 5. Asserter les faits
	t.Log("📝 Étape 5: Assertion des faits...")
	var activations int
	for i, fact := range reteFacts {
		n := network.AssertFact(fact)
		activations += n
		t.Logf("   - Fait %d asserté: %d activations", i, n)
	}
	
	// Vérifier les activations
	// On attend 1 activation (alice a plus de 25 ans, pas bob)
	expectedActivations := 1
	if activations != expectedActivations {
		t.Errorf("❌ Attendu %d activations, reçu %d", expectedActivations, activations)
	} else {
		t.Logf("✅ %d activations détectées (correct)", activations)
	}
	
	// 6. Vérifier les résultats
	t.Log("📝 Étape 6: Vérification résultats...")
	results := network.GetResults()
	if len(results) != expectedActivations {
		t.Errorf("❌ Attendu %d résultats, reçu %d", expectedActivations, len(results))
	} else {
		t.Logf("✅ %d résultats obtenus", len(results))
	}
	
	t.Log("")
	t.Log("🎉 CYCLE DE VIE COMPLET RÉUSSI")
}

func TestFactLifecycle_WithMultipleTypes(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - TYPES MULTIPLES")
	t.Log("======================================")
	
	input := `
		type User(#username: string, email: string, age: number)
		type Order(user: User, #orderNum: number, total: number, status: string)
		type Payment(order: Order, #paymentId: string, amount: number, method: string)
		
		alice = User("alice", "alice@example.com", 30)
		bob = User("bob", "bob@example.com", 25)
		
		order1 = Order(alice, 1001, 150.50, "pending")
		order2 = Order(bob, 1002, 75.00, "pending")
		
		Payment(order1, "PAY-001", 150.50, "card")
		Payment(order2, "PAY-002", 75.00, "paypal")
		
		{u: User, o: Order, p: Payment} / 
			o.user == u && p.order == o && o.status == "pending" ==> 
			Log("Payment " + p.paymentId + " for user " + u.username)
	`
	
	t.Log("📝 Parsing et validation...")
	program, err := constraint.ParseProgram(input)
	if err != nil {
		t.Fatalf("❌ Parsing: %v", err)
	}
	
	validator := constraint.NewProgramValidator()
	if err := validator.Validate(*program); err != nil {
		t.Fatalf("❌ Validation: %v", err)
	}
	
	t.Log("📝 Conversion et assertion...")
	reteFacts, err := constraint.ConvertFactsToReteFormat(*program)
	if err != nil {
		t.Fatalf("❌ Conversion: %v", err)
	}
	
	network := rete.NewNetwork()
	for _, expr := range program.Expressions {
		if err := network.CompileExpression(expr, program.Types); err != nil {
			t.Fatalf("❌ Compilation: %v", err)
		}
	}
	
	var activations int
	for _, fact := range reteFacts {
		activations += network.AssertFact(fact)
	}
	
	// On attend 2 activations (alice et bob)
	if activations != 2 {
		t.Errorf("❌ Attendu 2 activations, reçu %d", activations)
	} else {
		t.Logf("✅ %d activations (correct)", activations)
	}
	
	t.Log("🎉 Test réussi avec chaîne de 3 types")
}

func TestFactLifecycle_ErrorHandling(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION - GESTION ERREURS")
	t.Log("======================================")
	
	tests := []struct {
		name    string
		input   string
		errStep string // parsing, validation, conversion
	}{
		{
			name: "variable non définie",
			input: `
				type User(#name: string)
				type Login(user: User, #email: string)
				Login(unknownUser, "test@example.com")
			`,
			errStep: "validation",
		},
		{
			name: "_id_ manuel",
			input: `
				type User(#name: string)
				User(_id_: "manual", name: "Alice")
			`,
			errStep: "parsing",
		},
		{
			name: "type inexistant",
			input: `
				type Login(user: UnknownType, #email: string)
				Login(something, "test@example.com")
			`,
			errStep: "validation",
		},
		{
			name: "référence circulaire",
			input: `
				type A(b: B)
				type B(a: A)
			`,
			errStep: "validation",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, parseErr := constraint.ParseProgram(tt.input)
			
			if tt.errStep == "parsing" {
				if parseErr == nil {
					t.Errorf("❌ Attendu erreur de parsing, reçu nil")
				} else {
					t.Logf("✅ Erreur de parsing détectée: %v", parseErr)
				}
				return
			}
			
			if parseErr != nil {
				t.Fatalf("❌ Erreur de parsing inattendue: %v", parseErr)
			}
			
			if tt.errStep == "validation" {
				validator := constraint.NewProgramValidator()
				if err := validator.Validate(*program); err == nil {
					t.Errorf("❌ Attendu erreur de validation, reçu nil")
				} else {
					t.Logf("✅ Erreur de validation détectée: %v", err)
				}
				return
			}
			
			if tt.errStep == "conversion" {
				if _, err := constraint.ConvertFactsToReteFormat(*program); err == nil {
					t.Errorf("❌ Attendu erreur de conversion, reçu nil")
				} else {
					t.Logf("✅ Erreur de conversion détectée: %v", err)
				}
			}
		})
	}
}
```

### 3. Créer Tests End-to-End

#### Fichier : `tests/e2e/complete_program_test.go` (nouveau)

**Scénarios utilisateur complets** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/resinsec/tsd/constraint"
	"github.com/resinsec/tsd/rete"
)

func TestE2E_UserLoginScenario(t *testing.T) {
	t.Log("🧪 TEST E2E - SCÉNARIO USER/LOGIN")
	t.Log("==================================")
	
	// Lire le fichier TSD
	tsdFile := filepath.Join("testdata", "user_login.tsd")
	content, err := os.ReadFile(tsdFile)
	if err != nil {
		t.Fatalf("❌ Lecture fichier: %v", err)
	}
	
	t.Logf("📄 Fichier: %s (%d bytes)", tsdFile, len(content))
	
	// Parser
	program, err := constraint.ParseProgram(string(content))
	if err != nil {
		t.Fatalf("❌ Parsing: %v", err)
	}
	
	// Valider
	validator := constraint.NewProgramValidator()
	if err := validator.Validate(*program); err != nil {
		t.Fatalf("❌ Validation: %v", err)
	}
	
	// Convertir
	reteFacts, err := constraint.ConvertFactsToReteFormat(*program)
	if err != nil {
		t.Fatalf("❌ Conversion: %v", err)
	}
	
	// Créer réseau et compiler
	network := rete.NewNetwork()
	for _, expr := range program.Expressions {
		if err := network.CompileExpression(expr, program.Types); err != nil {
			t.Fatalf("❌ Compilation: %v", err)
		}
	}
	
	// Asserter
	for _, fact := range reteFacts {
		network.AssertFact(fact)
	}
	
	// Vérifier résultats
	results := network.GetResults()
	if len(results) == 0 {
		t.Error("❌ Aucun résultat obtenu")
	} else {
		t.Logf("✅ %d résultats obtenus", len(results))
	}
	
	t.Log("🎉 Scénario E2E réussi")
}

func TestE2E_OrderManagement(t *testing.T) {
	t.Log("🧪 TEST E2E - GESTION COMMANDES")
	t.Log("================================")
	
	tsdFile := filepath.Join("testdata", "order_management.tsd")
	content, err := os.ReadFile(tsdFile)
	if err != nil {
		t.Fatalf("❌ Lecture fichier: %v", err)
	}
	
	program, err := constraint.ParseProgram(string(content))
	if err != nil {
		t.Fatalf("❌ Parsing: %v", err)
	}
	
	validator := constraint.NewProgramValidator()
	if err := validator.Validate(*program); err != nil {
		t.Fatalf("❌ Validation: %v", err)
	}
	
	reteFacts, err := constraint.ConvertFactsToReteFormat(*program)
	if err != nil {
		t.Fatalf("❌ Conversion: %v", err)
	}
	
	network := rete.NewNetwork()
	for _, expr := range program.Expressions {
		if err := network.CompileExpression(expr, program.Types); err != nil {
			t.Fatalf("❌ Compilation: %v", err)
		}
	}
	
	for _, fact := range reteFacts {
		network.AssertFact(fact)
	}
	
	results := network.GetResults()
	t.Logf("✅ %d résultats obtenus", len(results))
	
	t.Log("🎉 Test réussi")
}

func TestE2E_AllExamples(t *testing.T) {
	t.Log("🧪 TEST E2E - TOUS LES EXEMPLES")
	t.Log("================================")
	
	// Trouver tous les fichiers .tsd dans examples/
	examplesDir := filepath.Join("..", "..", "examples")
	files, err := filepath.Glob(filepath.Join(examplesDir, "*.tsd"))
	if err != nil {
		t.Fatalf("❌ Recherche fichiers: %v", err)
	}
	
	if len(files) == 0 {
		t.Skip("⚠️  Aucun fichier exemple trouvé")
	}
	
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("❌ Lecture %s: %v", file, err)
			}
			
			program, err := constraint.ParseProgram(string(content))
			if err != nil {
				t.Fatalf("❌ Parsing %s: %v", file, err)
			}
			
			validator := constraint.NewProgramValidator()
			if err := validator.Validate(*program); err != nil {
				t.Fatalf("❌ Validation %s: %v", file, err)
			}
			
			reteFacts, err := constraint.ConvertFactsToReteFormat(*program)
			if err != nil {
				t.Fatalf("❌ Conversion %s: %v", file, err)
			}
			
			network := rete.NewNetwork()
			for _, expr := range program.Expressions {
				if err := network.CompileExpression(expr, program.Types); err != nil {
					t.Fatalf("❌ Compilation %s: %v", file, err)
				}
			}
			
			for _, fact := range reteFacts {
				network.AssertFact(fact)
			}
			
			t.Logf("✅ %s: validé et exécuté", filepath.Base(file))
		})
	}
	
	t.Logf("🎉 Tous les exemples validés (%d fichiers)", len(files))
}
```

### 4. Créer Fichiers TSD de Test

#### Fichier : `tests/e2e/testdata/user_login.tsd`

```tsd
// Scénario de test: Utilisateurs et logins
// Démontre les affectations de variables et comparaisons de faits

type User(#username: string, email: string, age: number, active: bool)
type Login(user: User, #sessionId: string, timestamp: number, ipAddress: string)
type AuditLog(login: Login, action: string, timestamp: number)

// Créer des utilisateurs
alice = User("alice", "alice@example.com", 30, true)
bob = User("bob", "bob@example.com", 25, true)
charlie = User("charlie", "charlie@example.com", 35, false)

// Créer des sessions
session1 = Login(alice, "SES-001", 1704067200, "192.168.1.10")
session2 = Login(bob, "SES-002", 1704067260, "192.168.1.11")
session3 = Login(charlie, "SES-003", 1704067320, "192.168.1.12")

// Créer des logs d'audit
AuditLog(session1, "login_success", 1704067200)
AuditLog(session2, "login_success", 1704067260)
AuditLog(session3, "login_failed", 1704067320)

// Règle 1: Identifier les logins d'utilisateurs actifs
{u: User, l: Login} / l.user == u && u.active == true ==> 
    Log("Active user login: " + u.username + " from " + l.ipAddress)

// Règle 2: Auditer les sessions avec utilisateur et log
{u: User, l: Login, a: AuditLog} / 
    l.user == u && a.login == l && a.action == "login_success" ==> 
    Log("Audit: " + u.username + " logged in successfully")

// Règle 3: Alerter pour utilisateurs inactifs
{u: User, l: Login} / l.user == u && u.active == false ==> 
    Log("ALERT: Inactive user attempted login: " + u.username)

// Règle 4: Utilisateurs seniors (>= 30 ans)
{u: User, l: Login} / l.user == u && u.age >= 30 ==> 
    Log("Senior user login: " + u.username + " (age: " + u.age + ")")
```

#### Fichier : `tests/e2e/testdata/order_management.tsd`

```tsd
// Scénario de test: Gestion de commandes
// Démontre les chaînes de références entre types

type Customer(#customerId: string, name: string, vipStatus: bool)
type Product(#productId: string, name: string, price: number, stock: number)
type Order(customer: Customer, #orderNumber: string, totalAmount: number, status: string)
type OrderLine(order: Order, product: Product, quantity: number, subtotal: number)
type Payment(order: Order, #paymentId: string, amount: number, method: string)

// Créer des clients
alice = Customer("CUST-001", "Alice Johnson", true)
bob = Customer("CUST-002", "Bob Smith", false)

// Créer des produits
laptop = Product("PROD-001", "Laptop", 1200.00, 10)
mouse = Product("PROD-002", "Mouse", 25.00, 50)
keyboard = Product("PROD-003", "Keyboard", 75.00, 30)

// Créer des commandes
order1 = Order(alice, "ORD-001", 1300.00, "pending")
order2 = Order(bob, "ORD-002", 100.00, "pending")

// Créer des lignes de commande
OrderLine(order1, laptop, 1, 1200.00)
OrderLine(order1, mouse, 4, 100.00)
OrderLine(order2, keyboard, 1, 75.00)
OrderLine(order2, mouse, 1, 25.00)

// Créer des paiements
Payment(order1, "PAY-001", 1300.00, "credit_card")
Payment(order2, "PAY-002", 100.00, "paypal")

// Règle 1: Commandes VIP
{c: Customer, o: Order} / o.customer == c && c.vipStatus == true ==> 
    Log("VIP Order: " + o.orderNumber + " from " + c.name)

// Règle 2: Commandes complètes (avec paiement)
{o: Order, p: Payment} / p.order == o && p.amount == o.totalAmount ==> 
    Log("Order " + o.orderNumber + " fully paid via " + p.method)

// Règle 3: Vérification stock
{o: Order, ol: OrderLine, prod: Product} / 
    ol.order == o && ol.product == prod && prod.stock < ol.quantity ==> 
    Log("ALERT: Insufficient stock for " + prod.name + " in order " + o.orderNumber)

// Règle 4: Commandes importantes (> 1000)
{c: Customer, o: Order} / o.customer == c && o.totalAmount > 1000.00 ==> 
    Log("Large order: " + o.orderNumber + " (" + o.totalAmount + ") from " + c.name)
```

#### Fichier : `tests/e2e/testdata/circular_reference_error.tsd`

```tsd
// Scénario de test: Erreur de référence circulaire
// Ce programme DOIT échouer à la validation

type Node(value: string, next: Node)

// Ceci devrait être rejeté à la validation
// car Node référence Node (cycle)
```

#### Fichier : `tests/e2e/testdata/undefined_variable_error.tsd`

```tsd
// Scénario de test: Erreur de variable non définie
// Ce programme DOIT échouer à la validation

type User(#name: string, age: number)
type Login(user: User, #email: string)

// alice n'est pas définie
Login(alice, "test@example.com")

// Ceci devrait être rejeté à la validation
```

### 5. Créer Exemples de Démonstration

#### Fichier : `examples/new_syntax_demo.tsd`

```tsd
// Démonstration de la nouvelle syntaxe TSD
// Affectations de variables et comparaisons de faits

// Définition des types
type User(#username: string, email: string, role: string)
type Post(author: User, #postId: string, title: string, content: string, likes: number)
type Comment(post: Post, author: User, #commentId: string, text: string)

// Créer des utilisateurs avec affectation
alice = User("alice", "alice@example.com", "admin")
bob = User("bob", "bob@example.com", "user")
charlie = User("charlie", "charlie@example.com", "moderator")

// Créer des posts
post1 = Post(alice, "POST-001", "Welcome to TSD", "This is the new syntax!", 42)
post2 = Post(bob, "POST-002", "Question about types", "How do fact types work?", 15)

// Créer des commentaires
Comment(post1, bob, "COM-001", "Great feature!")
Comment(post1, charlie, "COM-002", "Very useful!")
Comment(post2, alice, "COM-003", "Check the documentation")

// Règle 1: Posts d'admin
{u: User, p: Post} / p.author == u && u.role == "admin" ==> 
    Log("Admin post: " + p.title + " by " + u.username)

// Règle 2: Commentaires sur posts populaires
{p: Post, c: Comment, u: User} / 
    c.post == p && c.author == u && p.likes > 20 ==> 
    Log(u.username + " commented on popular post: " + p.title)

// Règle 3: Auto-commentaires
{p: Post, c: Comment} / c.post == p && c.author == p.author ==> 
    Log("Self-comment detected on post: " + p.title)
```

#### Fichier : `examples/advanced_relationships.tsd`

```tsd
// Exemple avancé: Relations complexes entre types
// Démontre les chaînes de références

type Organization(#orgId: string, name: string, country: string)
type Department(org: Organization, #deptId: string, name: string, budget: number)
type Employee(dept: Department, #empId: string, name: string, salary: number, position: string)
type Project(dept: Department, #projectId: string, name: string, status: string)
type Assignment(employee: Employee, project: Project, role: string, hours: number)

// Créer une organisation
acme = Organization("ORG-001", "ACME Corp", "USA")

// Créer des départements
engineering = Department(acme, "DEPT-001", "Engineering", 1000000.00)
marketing = Department(acme, "DEPT-002", "Marketing", 500000.00)

// Créer des employés
alice = Employee(engineering, "EMP-001", "Alice Johnson", 120000.00, "Senior Engineer")
bob = Employee(engineering, "EMP-002", "Bob Smith", 90000.00, "Engineer")
charlie = Employee(marketing, "EMP-003", "Charlie Brown", 85000.00, "Marketing Manager")

// Créer des projets
project1 = Project(engineering, "PROJ-001", "New Platform", "active")
project2 = Project(marketing, "PROJ-002", "Product Launch", "planning")

// Créer des affectations
Assignment(alice, project1, "Tech Lead", 40)
Assignment(bob, project1, "Developer", 40)
Assignment(charlie, project2, "Campaign Manager", 30)

// Règle 1: Projets actifs avec employés
{e: Employee, a: Assignment, p: Project} / 
    a.employee == e && a.project == p && p.status == "active" ==> 
    Log(e.name + " assigned to active project: " + p.name)

// Règle 2: Employés bien payés (>100k) sur projets
{e: Employee, a: Assignment, p: Project} / 
    a.employee == e && a.project == p && e.salary > 100000.00 ==> 
    Log("High-paid employee on project: " + e.name + " (" + p.name + ")")

// Règle 3: Projets et départements de la même org
{o: Organization, d: Department, p: Project} / 
    d.org == o && p.dept == d && o.country == "USA" ==> 
    Log("USA project: " + p.name + " in " + d.name)

// Règle 4: Budget et salaires
{d: Department, e: Employee} / 
    e.dept == d && e.salary > (d.budget / 10.0) ==> 
    Log("WARNING: " + e.name + " salary is >10% of department budget")
```

### 6. Créer Tests de Performance

#### Fichier : `tests/performance/benchmark_test.go` (nouveau)

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package performance

import (
	"testing"
	"github.com/resinsec/tsd/constraint"
	"github.com/resinsec/tsd/rete"
)

func BenchmarkFactGeneration(b *testing.B) {
	typeDef := constraint.TypeDefinition{
		Name: "User",
		Fields: []constraint.Field{
			{Name: "name", Type: "string", IsPrimaryKey: true},
			{Name: "age", Type: "number"},
		},
	}
	
	fact := constraint.Fact{
		TypeName: "User",
		Fields: []constraint.FactField{
			{Name: "name", Value: constraint.FactValue{Type: "string", Value: "Alice"}},
			{Name: "age", Value: constraint.FactValue{Type: "number", Value: 30.0}},
		},
	}
	
	ctx := constraint.NewFactContext([]constraint.TypeDefinition{typeDef})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := constraint.GenerateFactID(fact, typeDef, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFactGenerationWithReference(b *testing.B) {
	userType := constraint.TypeDefinition{
		Name: "User",
		Fields: []constraint.Field{
			{Name: "name", Type: "string", IsPrimaryKey: true},
		},
	}
	
	loginType := constraint.TypeDefinition{
		Name: "Login",
		Fields: []constraint.Field{
			{Name: "user", Type: "User"},
			{Name: "email", Type: "string", IsPrimaryKey: true},
		},
	}
	
	ctx := constraint.NewFactContext([]constraint.TypeDefinition{userType, loginType})
	ctx.RegisterVariable("alice", "User~Alice")
	
	loginFact := constraint.Fact{
		TypeName: "Login",
		Fields: []constraint.FactField{
			{Name: "user", Value: constraint.FactValue{Type: "variableReference", Value: "alice"}},
			{Name: "email", Value: constraint.FactValue{Type: "string", Value: "alice@ex.com"}},
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := constraint.GenerateFactID(loginFact, loginType, ctx)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProgramParsing(b *testing.B) {
	program := `
		type User(#name: string, age: number)
		alice = User("Alice", 30)
		bob = User("Bob", 25)
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := constraint.ParseProgram(program)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCompleteFlow(b *testing.B) {
	program := `
		type User(#name: string, age: number)
		type Login(user: User, #email: string)
		
		alice = User("Alice", 30)
		Login(alice, "alice@example.com")
		
		{u: User, l: Login} / l.user == u ==> Log("test")
	`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parsed, err := constraint.ParseProgram(program)
		if err != nil {
			b.Fatal(err)
		}
		
		validator := constraint.NewProgramValidator()
		if err := validator.Validate(*parsed); err != nil {
			b.Fatal(err)
		}
		
		reteFacts, err := constraint.ConvertFactsToReteFormat(*parsed)
		if err != nil {
			b.Fatal(err)
		}
		
		network := rete.NewNetwork()
		for _, expr := range parsed.Expressions {
			if err := network.CompileExpression(expr, parsed.Types); err != nil {
				b.Fatal(err)
			}
		}
		
		for _, fact := range reteFacts {
			network.AssertFact(fact)
		}
	}
}
```

### 7. Créer Script de Test E2E

#### Fichier : `scripts/run-e2e-tests.sh` (nouveau)

```bash
#!/bin/bash
# Script pour exécuter tous les tests E2E

set -e

echo "🧪 TESTS END-TO-END - NOUVELLE GESTION IDS"
echo "==========================================="
echo ""

# Fonction de log
log_success() {
    echo "✅ $1"
}

log_error() {
    echo "❌ $1"
}

log_info() {
    echo "ℹ️  $1"
}

# 1. Tests d'intégration
log_info "Étape 1/4: Tests d'intégration..."
if go test ./tests/integration/... -v -timeout 5m; then
    log_success "Tests d'intégration passent"
else
    log_error "Tests d'intégration échouent"
    exit 1
fi

echo ""

# 2. Tests E2E
log_info "Étape 2/4: Tests end-to-end..."
if go test ./tests/e2e/... -v -timeout 5m; then
    log_success "Tests E2E passent"
else
    log_error "Tests E2E échouent"
    exit 1
fi

echo ""

# 3. Tests de performance
log_info "Étape 3/4: Tests de performance..."
if go test ./tests/performance/... -bench=. -benchtime=1s; then
    log_success "Benchmarks exécutés"
else
    log_error "Benchmarks échouent"
    exit 1
fi

echo ""

# 4. Validation des exemples
log_info "Étape 4/4: Validation des exemples..."
example_count=0
failed_count=0

for file in examples/*.tsd; do
    if [ -f "$file" ]; then
        example_count=$((example_count + 1))
        basename=$(basename "$file")
        
        if go run cmd/tsd/main.go validate "$file" > /dev/null 2>&1; then
            log_success "Exemple $basename validé"
        else
            log_error "Exemple $basename invalide"
            failed_count=$((failed_count + 1))
        fi
    fi
done

echo ""

if [ $failed_count -eq 0 ]; then
    log_success "Tous les exemples validés ($example_count fichiers)"
else
    log_error "$failed_count/$example_count exemples invalides"
    exit 1
fi

# Résumé
echo ""
echo "==========================================="
log_success "TOUS LES TESTS E2E RÉUSSIS"
echo "==========================================="
echo ""
echo "Résumé:"
echo "  - Tests d'intégration: ✅"
echo "  - Tests E2E: ✅"
echo "  - Benchmarks: ✅"
echo "  - Exemples ($example_count): ✅"
echo ""
```

**Rendre exécutable** :
```bash
chmod +x scripts/run-e2e-tests.sh
```

### 8. Créer Documentation des Tests

#### Fichier : `tests/README.md` (mise à jour)

```markdown
# Tests TSD - Documentation

## Structure

```
tests/
├── integration/       # Tests d'intégration entre modules
├── e2e/              # Tests end-to-end avec fichiers .tsd
│   └── testdata/     # Fichiers .tsd de test
├── performance/      # Benchmarks
└── fixtures/         # Données partagées
```

## Tests d'Intégration

### Exécution

```bash
go test ./tests/integration/... -v
```

### Description

Les tests d'intégration valident l'interaction entre modules :
- Parser → Validation → RETE
- API → Constraint → RETE
- Cycle de vie complet des faits

### Fichiers

- `fact_lifecycle_test.go` - Cycle de vie complet
- Autres tests selon besoins

## Tests End-to-End

### Exécution

```bash
go test ./tests/e2e/... -v
```

### Description

Les tests E2E utilisent des programmes TSD complets pour valider des scénarios réels.

### Fichiers de Test

- `testdata/user_login.tsd` - Scénario utilisateurs/logins
- `testdata/order_management.tsd` - Gestion de commandes
- `testdata/circular_reference_error.tsd` - Test d'erreur
- `testdata/undefined_variable_error.tsd` - Test d'erreur

## Tests de Performance

### Exécution

```bash
go test ./tests/performance/... -bench=. -benchmem
```

### Benchmarks

- `BenchmarkFactGeneration` - Génération d'IDs
- `BenchmarkFactGenerationWithReference` - Avec références
- `BenchmarkProgramParsing` - Parsing
- `BenchmarkCompleteFlow` - Flow complet

## Script Global

```bash
# Exécuter tous les tests E2E
./scripts/run-e2e-tests.sh
```

## Exemples

Les fichiers dans `examples/` servent aussi de tests :

```bash
# Valider tous les exemples
for f in examples/*.tsd; do
    go run cmd/tsd/main.go validate "$f"
done
```

## Nouvelle Gestion des IDs

### Points Testés

1. **Affectations** : `alice = User(...)`
2. **Références** : `Login(alice, ...)`
3. **Comparaisons** : `l.user == u`
4. **Validation** : Interdiction de `_id_`
5. **Génération** : IDs automatiques

### Exemples de Tests

```go
// Test d'affectation
alice = User("Alice", 30)
Login(alice, "alice@example.com")

// Test de comparaison
{u: User, l: Login} / l.user == u ==> Log("Match")

// Test d'erreur
User(_id_: "manual") // ❌ Doit échouer
```

## Contribution

### Ajouter un Test E2E

1. Créer fichier `.tsd` dans `tests/e2e/testdata/`
2. Ajouter test Go dans `tests/e2e/`
3. Exécuter : `go test ./tests/e2e/... -v`

### Ajouter un Exemple

1. Créer fichier `.tsd` dans `examples/`
2. Documenter le scénario en commentaires
3. Valider : `go run cmd/tsd/main.go validate examples/mon_exemple.tsd`
```

---

## ✅ Critères de Succès

### Tests

```bash
# Tests d'intégration passent
go test ./tests/integration/... -v

# Tests E2E passent
go test ./tests/e2e/... -v

# Benchmarks s'exécutent
go test ./tests/performance/... -bench=.

# Script global réussit
./scripts/run-e2e-tests.sh
```

### Checklist

- [ ] Tests d'intégration créés
- [ ] Tests E2E créés
- [ ] Fichiers .tsd de test créés
- [ ] Exemples de démonstration créés
- [ ] Benchmarks ajoutés
- [ ] Script global créé
- [ ] Documentation mise à jour
- [ ] Tous les tests passent
- [ ] Pas de dégradation de performance

### Validation

```bash
make test-integration
make test-e2e
make test-performance
make test-complete
```

---

## 📊 Métriques Attendues

### Couverture E2E

- Scénarios utilisateur : > 5
- Cas d'erreur : > 3
- Exemples : > 3

### Performance

- Pas de régression > 10%
- Génération ID : < 1ms
- Parsing : < 10ms pour programmes typiques

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Analyser tests existants
2. ✅ Créer tests d'intégration
3. ✅ Créer tests E2E
4. ✅ Créer fichiers .tsd de test
5. ✅ Créer exemples de démonstration
6. ✅ Créer benchmarks
7. ✅ Créer script global
8. ✅ Mettre à jour documentation
9. ✅ Valider tous les tests

### Commandes

```bash
# Créer les répertoires
mkdir -p tests/integration
mkdir -p tests/e2e/testdata
mkdir -p tests/performance

# Créer les fichiers
touch tests/integration/fact_lifecycle_test.go
touch tests/e2e/complete_program_test.go
touch tests/performance/benchmark_test.go

# Créer les exemples
touch examples/new_syntax_demo.tsd
touch examples/advanced_relationships.tsd

# Créer le script
touch scripts/run-e2e-tests.sh
chmod +x scripts/run-e2e-tests.sh

# Exécuter
go test ./tests/integration/... -v
go test ./tests/e2e/... -v
go test ./tests/performance/... -bench=.
./scripts/run-e2e-tests.sh
```

---

## 📚 Références

- `scripts/new_ids/07-prompt-tests-unit.md` - Tests unitaires
- `scripts/new_ids/06-prompt-api-tsdio.md` - API
- `tests/` - Tests actuels
- `examples/` - Exemples actuels

---

## 📝 Notes

### Points d'Attention

1. **Tests réels** : Pas de mocks, programmes TSD complets

2. **Scénarios variés** : Couvrir cas nominaux et erreurs

3. **Performance** : Benchmarks pour détecter régressions

4. **Documentation** : Exemples servent de documentation

### Bonnes Pratiques

```go
// ✅ BON - Test E2E complet
func TestE2E_Scenario(t *testing.T) {
    t.Log("🧪 TEST E2E - SCÉNARIO")
    
    // 1. Parser
    program, err := constraint.ParseProgram(tsdCode)
    require.NoError(t, err)
    
    // 2. Valider
    validator := constraint.NewProgramValidator()
    require.NoError(t, validator.Validate(*program))
    
    // 3. Convertir
    reteFacts, err := constraint.ConvertFactsToReteFormat(*program)
    require.NoError(t, err)
    
    // 4. Exécuter
    network := rete.NewNetwork()
    // ... assertions et vérifications
    
    t.Log("✅ Scénario complet réussi")
}
```

---

## 🎯 Résultat Attendu

Après ce prompt :

```bash
# Tous les tests passent
$ ./scripts/run-e2e-tests.sh
🧪 TESTS END-TO-END - NOUVELLE GESTION IDS
===========================================

ℹ️  Étape 1/4: Tests d'intégration...
✅ Tests d'intégration passent

ℹ️  Étape 2/4: Tests end-to-end...
✅ Tests E2E passent

ℹ️  Étape 3/4: Tests de performance...
✅ Benchmarks exécutés

ℹ️  Étape 4/4: Validation des exemples...
✅ Exemple new_syntax_demo.tsd validé
✅ Exemple advanced_relationships.tsd validé
✅ Tous les exemples validés (5 fichiers)

===========================================
✅ TOUS LES TESTS E2E RÉUSSIS
===========================================

Résumé:
  - Tests d'intégration: ✅
  - Tests E2E: ✅
  - Benchmarks: ✅
  - Exemples (5): ✅
```

---

**Prompt suivant** : `09-prompt-documentation.md`

**Durée estimée** : 6-8 heures

**Complexité** : ⚠️ Moyenne-Élevée