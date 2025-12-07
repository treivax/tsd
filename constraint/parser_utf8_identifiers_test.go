// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

// TestBug_UTF8Support_Fixed vérifie le support UTF-8 dans les chaînes de caractères
func TestBug_UTF8Support_Fixed(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// Caractères accentués français
		{
			name: "French accents in string",
			input: `
type Person(name: string, city: string)
rule test1 : {p:Person} / p.name == "François" ==> action()
Person(name: "François", city: "Paris")
`,
			shouldParse: true,
			description: "Les accents français (é, è, à, ç) devraient fonctionner dans les chaînes",
		},
		{
			name: "French accents complex",
			input: `
type Person(name: string, description: string)
rule test2 : {p:Person} / p.description == "Étudiant à l'université" ==> action()
Person(name: "Marie", description: "Étudiant à l'université")
`,
			shouldParse: true,
			description: "Les phrases avec accents français devraient fonctionner",
		},
		// Caractères allemands
		{
			name: "German umlauts",
			input: `
type Person(name: string, city: string)
rule test3 : {p:Person} / p.city == "München" ==> action()
Person(name: "Hans", city: "München")
`,
			shouldParse: true,
			description: "Les trémas allemands (ä, ö, ü) devraient fonctionner",
		},
		// Caractères espagnols
		{
			name: "Spanish characters",
			input: `
type Person(name: string, greeting: string)
rule test4 : {p:Person} / p.greeting == "¡Hola señor!" ==> action()
Person(name: "José", greeting: "¡Hola señor!")
`,
			shouldParse: true,
			description: "Les caractères espagnols (ñ, ¡, ¿) devraient fonctionner",
		},
		// Caractères cyrilliques (russe)
		{
			name: "Russian cyrillic",
			input: `
type Person(name: string, city: string)
rule test5 : {p:Person} / p.city == "Москва" ==> action()
Person(name: "Иван", city: "Москва")
`,
			shouldParse: true,
			description: "Les caractères cyrilliques russes devraient fonctionner",
		},
		// Caractères chinois
		{
			name: "Chinese characters",
			input: `
type Person(name: string, city: string)
rule test6 : {p:Person} / p.city == "北京" ==> action()
Person(name: "李明", city: "北京")
`,
			shouldParse: true,
			description: "Les caractères chinois devraient fonctionner",
		},
		// Caractères japonais (hiragana, katakana, kanji)
		{
			name: "Japanese characters",
			input: `
type Person(name: string, city: string)
rule test7 : {p:Person} / p.city == "東京" ==> action()
Person(name: "田中さん", city: "東京")
`,
			shouldParse: true,
			description: "Les caractères japonais devraient fonctionner",
		},
		// Caractères arabes
		{
			name: "Arabic characters",
			input: `
type Person(name: string, city: string)
rule test8 : {p:Person} / p.city == "القاهرة" ==> action()
Person(name: "محمد", city: "القاهرة")
`,
			shouldParse: true,
			description: "Les caractères arabes devraient fonctionner",
		},
		// Emojis
		{
			name: "Emoji in strings",
			input: `
type Message(text: string, sentiment: string)
rule test9 : {m:Message} / m.sentiment == "😊" ==> action()
Message(text: "Hello", sentiment: "😊")
`,
			shouldParse: true,
			description: "Les emojis devraient fonctionner dans les chaînes",
		},
		// Caractères grecs
		{
			name: "Greek characters",
			input: `
type Symbol(name: string, value: string)
rule test10 : {s:Symbol} / s.value == "α β γ δ" ==> action()
Symbol(name: "greek", value: "α β γ δ")
`,
			shouldParse: true,
			description: "Les caractères grecs devraient fonctionner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if tt.shouldParse && err != nil {
				t.Errorf("%s: parsing devrait réussir mais a échoué: %v", tt.description, err)
			}
			if !tt.shouldParse && err == nil {
				t.Errorf("%s: parsing devrait échouer mais a réussi", tt.description)
			}
		})
	}
}

// TestBug_IdentifierStyles_Fixed vérifie le support de camelCase et snake_case dans les identifiants
func TestBug_IdentifierStyles_Fixed(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// camelCase dans les noms de types
		{
			name: "camelCase type name",
			input: `
type CustomerOrder(orderId: string, totalAmount: number)
rule test1 : {o:CustomerOrder} / o.totalAmount > 100 ==> action()
CustomerOrder(orderId: "O123", totalAmount: 150)
`,
			shouldParse: true,
			description: "camelCase devrait fonctionner pour les noms de types",
		},
		// snake_case dans les noms de types
		{
			name: "snake_case type name",
			input: `
type customer_order(order_id: string, total_amount: number)
rule test2 : {o:customer_order} / o.total_amount > 100 ==> action()
customer_order(order_id: "O123", total_amount: 150)
`,
			shouldParse: true,
			description: "snake_case devrait fonctionner pour les noms de types",
		},
		// camelCase dans les noms de champs
		{
			name: "camelCase field names",
			input: `
type Person(firstName: string, lastName: string, phoneNumber: string)
rule test3 : {p:Person} / LENGTH(p.firstName) > 3 ==> action()
Person(firstName: "John", lastName: "Doe", phoneNumber: "123456")
`,
			shouldParse: true,
			description: "camelCase devrait fonctionner pour les noms de champs",
		},
		// snake_case dans les noms de champs
		{
			name: "snake_case field names",
			input: `
type Person(first_name: string, last_name: string, phone_number: string)
rule test4 : {p:Person} / LENGTH(p.first_name) > 3 ==> action()
Person(first_name: "John", last_name: "Doe", phone_number: "123456")
`,
			shouldParse: true,
			description: "snake_case devrait fonctionner pour les noms de champs",
		},
		// camelCase dans les noms de règles
		{
			name: "camelCase rule name",
			input: `
type Person(name: string, age: number)
rule checkHighSalaryEmployee : {p:Person} / p.age > 18 ==> action()
Person(name: "Alice", age: 25)
`,
			shouldParse: true,
			description: "camelCase devrait fonctionner pour les noms de règles",
		},
		// snake_case dans les noms de règles
		{
			name: "snake_case rule name",
			input: `
type Person(name: string, age: number)
rule check_high_salary_employee : {p:Person} / p.age > 18 ==> action()
Person(name: "Alice", age: 25)
`,
			shouldParse: true,
			description: "snake_case devrait fonctionner pour les noms de règles",
		},
		// camelCase dans les noms d'actions
		{
			name: "camelCase action name",
			input: `
type Person(name: string, age: number)
action sendNotificationEmail(recipient: string)
rule test5 : {p:Person} / p.age > 18 ==> sendNotificationEmail(p.name)
Person(name: "Bob", age: 30)
`,
			shouldParse: true,
			description: "camelCase devrait fonctionner pour les noms d'actions",
		},
		// snake_case dans les noms d'actions
		{
			name: "snake_case action name",
			input: `
type Person(name: string, age: number)
action send_notification_email(recipient: string)
rule test6 : {p:Person} / p.age > 18 ==> send_notification_email(p.name)
Person(name: "Bob", age: 30)
`,
			shouldParse: true,
			description: "snake_case devrait fonctionner pour les noms d'actions",
		},
		// Mélange camelCase et snake_case
		{
			name: "mixed camelCase and snake_case",
			input: `
type CustomerOrder(order_id: string, totalAmount: number, customer_name: string)
action processOrder(orderId: string)
rule check_large_order : {o:CustomerOrder} / o.totalAmount > 1000 ==> processOrder(o.order_id)
CustomerOrder(order_id: "O999", totalAmount: 1500, customer_name: "Alice")
`,
			shouldParse: true,
			description: "Un mélange de camelCase et snake_case devrait fonctionner",
		},
		// Identifiants avec plusieurs underscores consécutifs
		{
			name: "multiple underscores",
			input: `
type Test(field__name: string)
rule test__rule : {t:Test} / t.field__name == "value" ==> action()
Test(field__name: "value")
`,
			shouldParse: true,
			description: "Plusieurs underscores consécutifs devraient être acceptés",
		},
		// Identifiants commençant par underscore
		{
			name: "leading underscore",
			input: `
type _InternalType(name: string)
rule _internal_rule : {t:_InternalType} / t.name == "test" ==> action()
_InternalType(name: "test")
`,
			shouldParse: true,
			description: "Les identifiants commençant par underscore devraient être acceptés",
		},
		// Identifiants avec chiffres
		{
			name: "identifiers with numbers",
			input: `
type Product2(name: string, version2: number)
rule check_version_2 : {p:Product2} / p.version2 > 1 ==> action()
Product2(name: "Widget", version2: 2)
`,
			shouldParse: true,
			description: "Les identifiants avec chiffres devraient fonctionner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if tt.shouldParse && err != nil {
				t.Errorf("%s: parsing devrait réussir mais a échoué: %v", tt.description, err)
			}
			if !tt.shouldParse && err == nil {
				t.Errorf("%s: parsing devrait échouer mais a réussi", tt.description)
			}
		})
	}
}

// TestBug_UTF8InIdentifiers_Fixed vérifie que les caractères UTF-8 (tous les scripts Unicode majeurs)
// sont acceptés dans les identifiants (noms de types, champs, règles, actions, variables)
func TestBug_UTF8InIdentifiers_Fixed(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// Identifiants avec accents français
		{
			name: "French accents in identifiers",
			input: `
type Personne(nom: string, prénom: string, âge: number)
rule règle1 : {p:Personne} / p.âge > 18 ==> action()
Personne(nom: "Dupont", prénom: "François", âge: 25)
`,
			shouldParse: true,
			description: "Les accents français dans les identifiants devraient fonctionner",
		},
		// Identifiants avec caractères chinois
		{
			name: "Chinese in identifiers",
			input: `
type 用户(姓名: string, 年龄: number)
rule 规则1 : {u:用户} / u.年龄 > 18 ==> action()
用户(姓名: "李明", 年龄: 25)
`,
			shouldParse: true,
			description: "Les caractères chinois dans les identifiants devraient fonctionner",
		},
		// Identifiants avec caractères cyrilliques
		{
			name: "Cyrillic in identifiers",
			input: `
type Пользователь(имя: string, возраст: number)
rule правило1 : {u:Пользователь} / u.возраст > 18 ==> action()
Пользователь(имя: "Иван", возраст: 25)
`,
			shouldParse: true,
			description: "Les caractères cyrilliques dans les identifiants devraient fonctionner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if tt.shouldParse && err != nil {
				t.Errorf("%s: parsing devrait réussir mais a échoué: %v", tt.description, err)
				t.Logf("Input:\n%s", tt.input)
			}
			if !tt.shouldParse && err == nil {
				t.Errorf("%s: parsing devrait échouer mais a réussi", tt.description)
			}
		})
	}
}
