// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

func TestNewTypeSyntax(t *testing.T) {
	t.Log("🧪 TEST NOUVELLE SYNTAXE DES TYPES")
	t.Log("====================================")

	tests := []struct {
		name     string
		input    string
		wantErr  bool
		typeName string
		fields   []string
	}{
		{
			name:     "Type simple avec un seul attribut",
			input:    "type Counter(value: number)",
			wantErr:  false,
			typeName: "Counter",
			fields:   []string{"value"},
		},
		{
			name:     "Type avec plusieurs attributs",
			input:    "type Person(name: string, age: number, active: bool)",
			wantErr:  false,
			typeName: "Person",
			fields:   []string{"name", "age", "active"},
		},
		{
			name:     "Type pour un produit",
			input:    "type Product(id: number, title: string, price: number, inStock: bool)",
			wantErr:  false,
			typeName: "Product",
			fields:   []string{"id", "title", "price", "inStock"},
		},
		{
			name:     "Type pour une commande",
			input:    "type Order(orderId: number, customerName: string, total: number, paid: bool)",
			wantErr:  false,
			typeName: "Order",
			fields:   []string{"orderId", "customerName", "total", "paid"},
		},
		{
			name:     "Type pour un événement",
			input:    "type Event(eventName: string, timestamp: number, severity: string, handled: bool)",
			wantErr:  false,
			typeName: "Event",
			fields:   []string{"eventName", "timestamp", "severity", "handled"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("  📝 Test: %s", tt.name)
			t.Logf("  📥 Input: %s", tt.input)

			result, err := Parse("test", []byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				t.Logf("  ✅ Erreur attendue: %v", err)
				return
			}

			// Convertir en Program
			program, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("❌ ConvertResultToProgram() error = %v", err)
			}

			// Vérifier qu'on a bien un type
			if len(program.Types) != 1 {
				t.Fatalf("❌ Expected 1 type, got %d", len(program.Types))
			}

			typeDef := program.Types[0]

			// Vérifier le nom du type
			if typeDef.Name != tt.typeName {
				t.Errorf("❌ Type name = %v, want %v", typeDef.Name, tt.typeName)
			}

			// Vérifier le nombre de champs
			if len(typeDef.Fields) != len(tt.fields) {
				t.Errorf("❌ Number of fields = %v, want %v", len(typeDef.Fields), len(tt.fields))
			}

			// Vérifier les noms des champs
			for i, expectedField := range tt.fields {
				if i >= len(typeDef.Fields) {
					break
				}
				if typeDef.Fields[i].Name != expectedField {
					t.Errorf("❌ Field[%d] name = %v, want %v", i, typeDef.Fields[i].Name, expectedField)
				}
			}

			t.Logf("  ✅ Type '%s' parsé avec succès avec %d champs", typeDef.Name, len(typeDef.Fields))
		})
	}
}

func TestNewActionSyntax(t *testing.T) {
	t.Log("🧪 TEST NOUVELLE SYNTAXE DES ACTIONS")
	t.Log("=====================================")

	tests := []struct {
		name       string
		input      string
		wantErr    bool
		actionName string
		paramCount int
	}{
		{
			name:       "Action simple avec un argument primitif",
			input:      "action log(message: string)",
			wantErr:    false,
			actionName: "log",
			paramCount: 1,
		},
		{
			name:       "Action avec plusieurs arguments primitifs",
			input:      "action notify(recipient: string, message: string, priority: number)",
			wantErr:    false,
			actionName: "notify",
			paramCount: 3,
		},
		{
			name:       "Action avec valeur par défaut",
			input:      "action notify(recipient: string, message: string, priority: number = 1)",
			wantErr:    false,
			actionName: "notify",
			paramCount: 3,
		},
		{
			name:       "Action avec argument optionnel",
			input:      "action updateUser(user: User, active: bool?)",
			wantErr:    false,
			actionName: "updateUser",
			paramCount: 2,
		},
		{
			name:       "Action avec optionnel et défaut",
			input:      "action processOrder(order: Order, discount: number?, notify: bool = false)",
			wantErr:    false,
			actionName: "processOrder",
			paramCount: 3,
		},
		{
			name:       "Action sans arguments",
			input:      "action triggerBatchProcess()",
			wantErr:    false,
			actionName: "triggerBatchProcess",
			paramCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("  📝 Test: %s", tt.name)
			t.Logf("  📥 Input: %s", tt.input)

			result, err := Parse("test", []byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				t.Logf("  ✅ Erreur attendue: %v", err)
				return
			}

			// Convertir en Program
			program, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("❌ ConvertResultToProgram() error = %v", err)
			}

			// Vérifier qu'on a bien une action
			if len(program.Actions) != 1 {
				t.Fatalf("❌ Expected 1 action, got %d", len(program.Actions))
			}

			actionDef := program.Actions[0]

			// Vérifier le nom de l'action
			if actionDef.Name != tt.actionName {
				t.Errorf("❌ Action name = %v, want %v", actionDef.Name, tt.actionName)
			}

			// Vérifier le nombre de paramètres
			if len(actionDef.Parameters) != tt.paramCount {
				t.Errorf("❌ Number of parameters = %v, want %v", len(actionDef.Parameters), tt.paramCount)
			}

			t.Logf("  ✅ Action '%s' parsée avec succès avec %d paramètres", actionDef.Name, len(actionDef.Parameters))
		})
	}
}

func TestActionValidation(t *testing.T) {
	t.Log("🧪 TEST VALIDATION DES ACTIONS")
	t.Log("===============================")

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name: "Validation réussie - action avec type défini",
			input: `
type Person(name: string, age: number)
action savePerson(person: Person)
rule r1 : {p: Person} / p.age > 18 ==> savePerson(p)
`,
			wantErr: false,
		},
		{
			name: "Validation réussie - action avec primitifs",
			input: `
type Person(name: string, age: number)
action log(message: string)
rule r1 : {p: Person} / p.age > 18 ==> log(p.name)
`,
			wantErr: false,
		},
		{
			name: "Erreur - action non définie",
			input: `
type Person(name: string, age: number)
rule r1 : {p: Person} / p.age > 18 ==> undefinedAction(p)
`,
			wantErr: true,
			errMsg:  "action 'undefinedAction' is not defined",
		},
		{
			name: "Erreur - type de paramètre incorrect",
			input: `
type Person(name: string, age: number)
action log(message: string)
rule r1 : {p: Person} / p.age > 18 ==> log(p.age)
`,
			wantErr: true,
			errMsg:  "type mismatch",
		},
		{
			name: "Erreur - nombre d'arguments incorrect",
			input: `
type Person(name: string, age: number)
action notify(recipient: string, message: string)
rule r1 : {p: Person} / p.age > 18 ==> notify(p.name)
`,
			wantErr: true,
			errMsg:  "requires at least",
		},
		{
			name: "Validation réussie - paramètre avec valeur par défaut",
			input: `
type Person(name: string, age: number)
action notify(recipient: string, message: string, priority: number = 1)
rule r1 : {p: Person} / p.age > 18 ==> notify(p.name, "Hello")
`,
			wantErr: false,
		},
		{
			name: "Erreur - type personnalisé non défini dans action",
			input: `
type Person(name: string, age: number)
action saveUser(user: User)
rule r1 : {p: Person} / p.age > 18 ==> saveUser(p)
`,
			wantErr: true,
			errMsg:  "type 'User' is not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("  📝 Test: %s", tt.name)

			result, err := Parse("test", []byte(tt.input))
			if err != nil {
				if tt.wantErr {
					t.Logf("  ✅ Erreur de parsing attendue: %v", err)
					return
				}
				t.Fatalf("❌ Parse() error = %v", err)
			}

			// Valider le programme
			err = ValidateConstraintProgram(result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ ValidateConstraintProgram() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				t.Logf("  📄 Message d'erreur: %v", err)
				if tt.errMsg != "" {
					// Vérifier que le message d'erreur contient le texte attendu
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("❌ Error message should contain '%s', got '%s'", tt.errMsg, err.Error())
					}
				}
				t.Logf("  ✅ Erreur de validation attendue")
				return
			}

			t.Logf("  ✅ Validation réussie")
		})
	}
}

func TestCompleteProgram(t *testing.T) {
	t.Log("🧪 TEST PROGRAMME COMPLET AVEC NOUVELLE SYNTAXE")
	t.Log("================================================")

	input := `
type Person(name: string, age: number, active: bool)
type Order(orderId: number, customerName: string, total: number, paid: bool)

action log(message: string)
action notify(recipient: string, message: string, priority: number = 1)
action savePerson(person: Person)
action processOrder(order: Order, discount: number?, notifyCustomer: bool = true)

rule r1 : {p: Person} / p.age > 18 AND p.active == true ==> log(p.name)
rule r2 : {p: Person} / p.age > 65 ==> notify(p.name, "Senior discount available", 2)
rule r3 : {o: Order} / o.paid == false AND o.total > 100 ==> processOrder(o, 10)
rule r4 : {p: Person, o: Order} / p.name == o.customerName AND o.total > 500 ==> notify(p.name, "VIP order")
`

	t.Logf("  📥 Parsing programme complet...")

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("❌ Parse() error = %v", err)
	}

	// Convertir en Program
	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ ConvertResultToProgram() error = %v", err)
	}

	t.Logf("  📊 Statistiques du programme:")
	t.Logf("    - Types définis: %d", len(program.Types))
	t.Logf("    - Actions définies: %d", len(program.Actions))
	t.Logf("    - Règles définies: %d", len(program.Expressions))

	// Vérifier les types
	if len(program.Types) != 2 {
		t.Errorf("❌ Expected 2 types, got %d", len(program.Types))
	}

	// Vérifier les actions
	if len(program.Actions) != 4 {
		t.Errorf("❌ Expected 4 actions, got %d", len(program.Actions))
	}

	// Vérifier les règles
	if len(program.Expressions) != 4 {
		t.Errorf("❌ Expected 4 rules, got %d", len(program.Expressions))
	}

	// Valider le programme
	err = ValidateConstraintProgram(result)
	if err != nil {
		t.Fatalf("❌ ValidateConstraintProgram() error = %v", err)
	}

	t.Logf("  ✅ Programme complet parsé et validé avec succès")
}

func TestActionWithComplexParameters(t *testing.T) {
	t.Log("🧪 TEST ACTIONS AVEC PARAMÈTRES COMPLEXES")
	t.Log("==========================================")

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "Action avec type personnalisé",
			input: `
type User(id: number, name: string, role: string)
action createUser(user: User)
`,
			wantErr: false,
		},
		{
			name: "Action avec plusieurs types personnalisés",
			input: `
type User(id: number, name: string)
type Role(name: string, permissions: string)
action assignRole(user: User, role: Role)
`,
			wantErr: false,
		},
		{
			name: "Action avec mix de types primitifs et personnalisés",
			input: `
type Product(id: number, name: string, price: number)
action updatePrice(product: Product, newPrice: number, reason: string = "price_change")
`,
			wantErr: false,
		},
		{
			name: "Action avec tous optionnels après requis",
			input: `
type Order(id: number, total: number)
action processPayment(order: Order, method: string, fee: number?, sendReceipt: bool = true)
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("  📝 Test: %s", tt.name)

			result, err := Parse("test", []byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				t.Logf("  ✅ Erreur attendue: %v", err)
				return
			}

			// Convertir et valider
			program, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("❌ ConvertResultToProgram() error = %v", err)
			}

			// Valider les définitions d'action
			validator := NewActionValidator(program.Actions, program.Types)
			errs := validator.ValidateActionDefinitions()
			if len(errs) > 0 {
				t.Fatalf("❌ ValidateActionDefinitions() errors = %v", errs)
			}

			t.Logf("  ✅ Actions avec paramètres complexes validées")
		})
	}
}
