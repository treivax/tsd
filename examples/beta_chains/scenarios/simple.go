// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package scenarios

import (
	"fmt"

	"github.com/treivax/tsd/rete"
)

// RunSimple creates a simple scenario with 5 rules sharing the same join
// This demonstrates high sharing ratio (60-80%)
func RunSimple(network *rete.ReteNetwork) int {
	fmt.Println("  📝 Creating 5 rules with shared Person ⋈ Order join...")

	rules := []struct {
		id          string
		name        string
		description string
	}{
		{"rule_high_spender", "High Spender", "Orders > $1000"},
		{"rule_vip_customer", "VIP Customer", "High value orders"},
		{"rule_discount_eligible", "Discount Eligible", "Qualify for discount"},
		{"rule_loyalty_points", "Loyalty Points", "Award loyalty points"},
		{"rule_send_confirmation", "Send Confirmation", "Send email confirmation"},
	}

	for i, ruleSpec := range rules {
		rule := createSimpleRule(ruleSpec.id, ruleSpec.name, ruleSpec.description)
		err := network.AddRule(rule)
		if err != nil {
			fmt.Printf("  ❌ Error adding rule %s: %v\n", ruleSpec.id, err)
			continue
		}
		fmt.Printf("  [%d/5] %-25s ✅\n", i+1, ruleSpec.name)
	}

	fmt.Println("  ✅ Simple scenario complete")
	return len(rules)
}

func createSimpleRule(id, name, description string) *rete.Rule {
	// All rules share the same join pattern:
	// Person ⋈ Order where p.id == o.personId AND o.amount > 1000
	return &rete.Rule{
		ID:   id,
		Name: name,
		Patterns: []rete.Pattern{
			{
				Type: "Person",
				Var:  "p",
				Conditions: []rete.Condition{
					{
						Type:  "field_exists",
						Field: "id",
					},
				},
			},
			{
				Type: "Order",
				Var:  "o",
				Conditions: []rete.Condition{
					{
						Type:  "field_compare",
						Field: "personId",
						Op:    "==",
						Value: "p.id",
					},
					{
						Type:  "field_compare",
						Field: "amount",
						Op:    ">",
						Value: 1000,
					},
				},
			},
		},
		Actions: []rete.Action{
			{
				Type: "log",
				Params: map[string]interface{}{
					"message": fmt.Sprintf("%s: %s", name, description),
				},
			},
		},
	}
}
