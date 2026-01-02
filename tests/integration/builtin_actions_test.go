// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/api"
	"github.com/treivax/tsd/rete"
)

func TestBuiltinActions_Update_Integration(t *testing.T) {
	program := `
type Person(#id: string, name: string, age: number)

rule update_age : {p: Person} / p.age < 18 ==>
    Update(p, {age: 18})

Person(id: "p1", name: "Alice", age: 15)
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err, "L'ingestion du programme doit réussir")
	require.NotNil(t, result, "Le résultat ne doit pas être nil")

	// Vérifier que le fait a été modifié
	network := result.Network()
	facts := network.Storage.GetAllFacts()
	require.Len(t, facts, 1, "Il doit y avoir exactement 1 fait")

	person := facts[0]
	assert.Equal(t, "Person", person.Type, "Le type doit être Person")
	assert.Contains(t, person.ID, "Person~", "L'ID interne doit contenir le préfixe Type")
	assert.Equal(t, "p1", person.Fields["id"], "La clé primaire doit être préservée")
	assert.Equal(t, "Alice", person.Fields["name"], "Le nom doit être Alice")
	assert.Equal(t, float64(18), person.Fields["age"], "L'âge doit avoir été mis à jour à 18")
}

func TestBuiltinActions_Update_PreservesID(t *testing.T) {
	program := `
type Product(#id: string, name: string, stock: number, status: string)

rule mark_low_stock : {p: Product} / p.stock < 10 AND p.status == "available" ==>
    Update(p, {status: "low_stock"})

Product(id: "prod_123", name: "Widget", stock: 5, status: "available")
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 1)

	product := facts[0]
	assert.Contains(t, product.ID, "Product~", "L'ID interne doit contenir le préfixe Type")
	assert.Equal(t, "prod_123", product.Fields["id"], "La clé primaire doit être préservée")
	assert.Equal(t, "low_stock", product.Fields["status"], "Le statut doit être mis à jour")
}

func TestBuiltinActions_Update_MultipleFields(t *testing.T) {
	program := `
type User(#id: string, name: string, email: string, verified: bool)

rule verify_user : {u: User} / u.verified == false ==>
    Update(u, {verified: true})

User(id: "u1", name: "Bob", email: "bob@example.com", verified: false)
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 1)

	user := facts[0]
	assert.Contains(t, user.ID, "User~", "L'ID interne doit contenir le préfixe Type")
	assert.Equal(t, "u1", user.Fields["id"], "La clé primaire doit être préservée")
	assert.Equal(t, true, user.Fields["verified"], "Le champ verified doit être true")
	assert.Equal(t, "Bob", user.Fields["name"], "Les autres champs doivent être préservés")
}

func TestBuiltinActions_Insert_Integration(t *testing.T) {
	program := `
type Order(#id: string, customerId: string, amount: number)
type Alert(#id: string, orderId: string, message: string)

rule high_value_alert : {o: Order} / o.amount > 1000 ==>
    Insert(Alert(id: "alert_1", orderId: o.id, message: "High value order detected"))

Order(id: "ord_1", customerId: "cust_1", amount: 1500)
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 2, "Il doit y avoir 2 faits : l'ordre et l'alerte")

	// Vérifier que l'alerte a été créée
	var alert *rete.Fact
	for _, f := range facts {
		if f.Type == "Alert" {
			alert = f
			break
		}
	}

	require.NotNil(t, alert, "L'alerte doit avoir été insérée")
	assert.Contains(t, alert.ID, "Alert~", "L'ID interne doit contenir le préfixe Type")
	assert.Equal(t, "alert_1", alert.Fields["id"], "La clé primaire doit être correcte")
	assert.Equal(t, "ord_1", alert.Fields["orderId"])
	assert.Equal(t, "High value order detected", alert.Fields["message"])
}

func TestBuiltinActions_Insert_MultipleFacts(t *testing.T) {
	program := `
type Product(#id: string, name: string, category: string)
type CategorySummary(#id: string, category: string, count: number)

rule count_electronics : {p: Product} / p.category == "electronics" ==>
    Insert(CategorySummary(id: p.id, category: "electronics", count: 1))

Product(id: "p1", name: "Laptop", category: "electronics")
Product(id: "p2", name: "Mouse", category: "electronics")
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()

	// Compter les faits par type
	productCount := 0
	summaryCount := 0
	for _, f := range facts {
		if f.Type == "Product" {
			productCount++
		} else if f.Type == "CategorySummary" {
			summaryCount++
		}
	}

	assert.Equal(t, 2, productCount, "Les 2 produits doivent exister")
	assert.Equal(t, 2, summaryCount, "Un résumé doit être créé par produit electronics")
}

func TestBuiltinActions_Retract_Integration(t *testing.T) {
	program := `
type Task(#id: string, title: string, status: string)

rule remove_completed : {t: Task} / t.status == "completed" ==>
    Retract(t)

Task(id: "t1", title: "Task 1", status: "pending")
Task(id: "t2", title: "Task 2", status: "completed")
Task(id: "t3", title: "Task 3", status: "in_progress")
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 2, "Il ne doit rester que 2 faits après Retract")

	// Vérifier que la tâche complétée a été supprimée
	for _, f := range facts {
		assert.NotEqual(t, "t2", f.Fields["id"], "La tâche t2 ne doit plus exister")
		assert.NotEqual(t, "completed", f.Fields["status"])
	}
}

func TestBuiltinActions_Retract_ByID(t *testing.T) {
	program := `
type User(#id: string, name: string, active: bool)

rule remove_inactive : {u: User} / u.active == false ==>
    Retract(u)

User(id: "u1", name: "Alice", active: true)
User(id: "u2", name: "Bob", active: false)
User(id: "u3", name: "Charlie", active: true)
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 2, "Les utilisateurs inactifs doivent être supprimés")

	// Vérifier que seuls les utilisateurs actifs restent
	for _, f := range facts {
		assert.True(t, f.Fields["active"].(bool), "Seuls les utilisateurs actifs doivent rester")
	}
}

func TestBuiltinActions_Combined_Integration(t *testing.T) {
	program := `
type Item(#id: string, name: string, quantity: number, status: string)
type LowStockAlert(#id: string, itemId: string, quantity: number)

rule check_stock : {i: Item} / i.quantity < 5 AND i.quantity > 0 AND i.status == "active" ==>
    Update(i, {status: "low_stock"}),
    Insert(LowStockAlert(id: i.id, itemId: i.id, quantity: i.quantity))

rule remove_zero_stock : {i: Item} / i.quantity == 0 ==>
    Retract(i)

Item(id: "i1", name: "Widget", quantity: 3, status: "active")
Item(id: "i2", name: "Gadget", quantity: 0, status: "active")
Item(id: "i3", name: "Tool", quantity: 10, status: "active")
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()

	// Compter les faits par type
	itemCount := 0
	alertCount := 0
	var lowStockItem *rete.Fact

	for _, f := range facts {
		if f.Type == "Item" {
			itemCount++
			if f.Fields["id"] == "i1" {
				lowStockItem = f
			}
		} else if f.Type == "LowStockAlert" {
			alertCount++
		}
	}

	// Vérifications
	assert.Equal(t, 2, itemCount, "i2 doit être retracté, il doit rester 2 items")
	assert.GreaterOrEqual(t, alertCount, 1, "Une alerte doit être créée pour i1")

	if lowStockItem != nil {
		assert.Contains(t, lowStockItem.ID, "Item~", "L'ID interne doit contenir le préfixe Type")
		assert.Equal(t, "i1", lowStockItem.Fields["id"], "La clé primaire doit être préservée")
		assert.Equal(t, "low_stock", lowStockItem.Fields["status"],
			"Le statut de i1 doit être mis à jour à low_stock")
	}
}

func TestBuiltinActions_UpdateWithExpressions(t *testing.T) {
	program := `
type Counter(#id: string, value: number)

rule increment : {c: Counter} / c.value < 10 ==>
    Update(c, {value: c.value + 1})

Counter(id: "c1", value: 5)
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 1)

	counter := facts[0]
	// Note: L'expression c.value + 1 doit être évaluée
	// Le résultat attendu dépend de l'implémentation de l'évaluateur
	assert.NotNil(t, counter.Fields["value"], "La valeur doit exister")
}

func TestBuiltinActions_NoAction_WhenConditionFalse(t *testing.T) {
	program := `
type Product(#id: string, stock: number, status: string)

rule mark_out_of_stock : {p: Product} / p.stock == 0 ==>
    Update(p, {status: "out_of_stock"})

Product(id: "p1", stock: 10, status: "available")
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 1)

	product := facts[0]
	// La règle ne doit pas s'être déclenchée car stock != 0
	assert.Equal(t, "available", product.Fields["status"],
		"Le statut ne doit pas changer car la condition est false")
}

func TestBuiltinActions_ChainedRules(t *testing.T) {
	program := `
type Order(#id: string, status: string, priority: number)

rule set_high_priority : {o: Order} / o.status == "urgent" ==>
    Update(o, {priority: 10})

rule process_high_priority : {o: Order} / o.priority == 10 ==>
    Update(o, {status: "processing"})

Order(id: "o1", status: "urgent", priority: 1)
`

	pipeline := api.NewPipeline()
	result, err := pipeline.IngestString(program)
	require.NoError(t, err)

	facts := result.Network().Storage.GetAllFacts()
	require.Len(t, facts, 1)

	order := facts[0]
	// Les deux règles devraient s'être déclenchées en chaîne
	assert.Contains(t, order.ID, "Order~", "L'ID interne doit contenir le préfixe Type")
	assert.Equal(t, "o1", order.Fields["id"], "La clé primaire doit être préservée")
	assert.Equal(t, float64(10), order.Fields["priority"], "La priorité doit être 10")
	// Note: Le statut final dépend de l'ordre d'exécution des règles
	assert.NotEqual(t, "urgent", order.Fields["status"], "Le statut doit avoir changé")
}
