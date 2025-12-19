// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestFactContext(t *testing.T) {
	t.Log("🧪 TEST FACT CONTEXT")
	t.Log("====================")

	userType := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "name", Type: "string", IsPrimaryKey: true},
		},
	}

	ctx := NewFactContext([]TypeDefinition{userType})

	// Test 1: Enregistrer une variable
	ctx.RegisterVariable("alice", "User~Alice")

	// Test 2: Résoudre la variable
	id, err := ctx.ResolveVariable("alice")
	if err != nil {
		t.Fatalf("❌ Erreur de résolution: %v", err)
	}

	if id != "User~Alice" {
		t.Errorf("❌ ID attendu 'User~Alice', reçu '%s'", id)
	}
	t.Logf("✅ Variable résolue correctement: alice → %s", id)

	// Test 3: Variable non définie
	_, err = ctx.ResolveVariable("bob")
	if err == nil {
		t.Error("❌ Attendu une erreur pour variable non définie")
	} else {
		t.Logf("✅ Erreur attendue pour variable non définie: %v", err)
	}

	// Test 4: Vérifier le TypeMap
	if len(ctx.TypeMap) != 1 {
		t.Errorf("❌ TypeMap devrait contenir 1 type, contient %d", len(ctx.TypeMap))
	}

	if _, exists := ctx.TypeMap["User"]; !exists {
		t.Error("❌ Type 'User' devrait être dans TypeMap")
	} else {
		t.Log("✅ TypeMap correctement initialisé")
	}

	t.Log("✅ Contexte fonctionne correctement")
}

func TestConvertFieldValueToString(t *testing.T) {
	t.Log("🧪 TEST CONVERSION VALEURS")
	t.Log("==========================")

	ctx := NewFactContext(nil)
	ctx.RegisterVariable("alice", "User~Alice")

	tests := []struct {
		name    string
		value   FactValue
		field   Field
		ctx     *FactContext
		want    string
		wantErr bool
	}{
		{
			name:  "string primitive",
			value: FactValue{Type: ValueTypeString, Value: "test"},
			field: Field{Name: "name", Type: "string"},
			want:  "test",
		},
		{
			name:  "identifier as string",
			value: FactValue{Type: ValueTypeIdentifier, Value: "identifier_value"},
			field: Field{Name: "code", Type: "string"},
			want:  "identifier_value",
		},
		{
			name:  "number entier",
			value: FactValue{Type: ValueTypeNumber, Value: float64(42)},
			field: Field{Name: "age", Type: "number"},
			want:  "42",
		},
		{
			name:  "number décimal",
			value: FactValue{Type: ValueTypeNumber, Value: 3.14},
			field: Field{Name: "price", Type: "number"},
			want:  "3.14",
		},
		{
			name:  "boolean true",
			value: FactValue{Type: ValueTypeBoolean, Value: true},
			field: Field{Name: "active", Type: "bool"},
			want:  "true",
		},
		{
			name:  "boolean false",
			value: FactValue{Type: ValueTypeBoolean, Value: false},
			field: Field{Name: "active", Type: "bool"},
			want:  "false",
		},
		{
			name:  "variable reference",
			value: FactValue{Type: "variableReference", Value: "alice"},
			field: Field{Name: "user", Type: "User"},
			ctx:   ctx,
			want:  "User~Alice",
		},
		{
			name:    "variable non définie",
			value:   FactValue{Type: "variableReference", Value: "bob"},
			field:   Field{Name: "user", Type: "User"},
			ctx:     ctx,
			wantErr: true,
		},
		{
			name:    "variable sans contexte",
			value:   FactValue{Type: "variableReference", Value: "alice"},
			field:   Field{Name: "user", Type: "User"},
			ctx:     nil,
			wantErr: true,
		},
		{
			name:    "type non supporté",
			value:   FactValue{Type: "unknown_type", Value: "value"},
			field:   Field{Name: "field", Type: "string"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertFieldValueToString(tt.value, tt.field, tt.ctx)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("❌ Erreur inattendue: %v", err)
			}

			if got != tt.want {
				t.Errorf("❌ Attendu '%s', reçu '%s'", tt.want, got)
			} else {
				t.Logf("✅ Conversion correcte: %s → %s", tt.name, got)
			}
		})
	}

	t.Log("✅ Toutes les conversions validées")
}

func TestGenerateFactID_WithFactReference(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION ID - RÉFÉRENCE DE FAIT")
	t.Log("==========================================")

	// Définir les types
	userType := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "name", Type: "string", IsPrimaryKey: true},
			{Name: "age", Type: "number"},
		},
	}

	loginType := TypeDefinition{
		Name: "Login",
		Fields: []Field{
			{Name: "user", Type: "User"},
			{Name: "email", Type: "string", IsPrimaryKey: true},
			{Name: "password", Type: "string"},
		},
	}

	// Créer le contexte
	ctx := NewFactContext([]TypeDefinition{userType, loginType})

	// Créer le fait User
	userFact := Fact{
		TypeName: "User",
		Fields: []FactField{
			{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
			{Name: "age", Value: FactValue{Type: "number", Value: float64(30)}},
		},
	}

	// Générer l'ID de User
	userID, err := GenerateFactID(userFact, userType, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur génération ID User: %v", err)
	}

	expectedUserID := "User~Alice"
	if userID != expectedUserID {
		t.Errorf("❌ ID User attendu '%s', reçu '%s'", expectedUserID, userID)
	}
	t.Logf("✅ ID User généré: %s", userID)

	// Enregistrer la variable alice
	ctx.RegisterVariable("alice", userID)

	// Créer le fait Login qui référence alice
	loginFact := Fact{
		TypeName: "Login",
		Fields: []FactField{
			{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
			{Name: "email", Value: FactValue{Type: "string", Value: "alice@example.com"}},
			{Name: "password", Value: FactValue{Type: "string", Value: "secret"}},
		},
	}

	// Générer l'ID de Login
	loginID, err := GenerateFactID(loginFact, loginType, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur génération ID Login: %v", err)
	}

	// L'ID devrait utiliser l'email dans sa clé primaire (email est la PK, pas user)
	expectedLoginID := "Login~alice@example.com"
	if loginID != expectedLoginID {
		t.Errorf("❌ ID Login attendu '%s', reçu '%s'", expectedLoginID, loginID)
	}
	t.Logf("✅ ID Login généré: %s", loginID)
}

func TestGenerateFactID_CompositeKeyWithFact(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION ID - CLÉ COMPOSITE AVEC FAIT")
	t.Log("=================================================")

	// Type User
	userType := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "name", Type: "string", IsPrimaryKey: true},
		},
	}

	// Type Order avec clé composite incluant une référence
	orderType := TypeDefinition{
		Name: "Order",
		Fields: []Field{
			{Name: "user", Type: "User", IsPrimaryKey: true},
			{Name: "orderNum", Type: "number", IsPrimaryKey: true},
			{Name: "total", Type: "number"},
		},
	}

	ctx := NewFactContext([]TypeDefinition{userType, orderType})

	// Créer User
	userFact := Fact{
		TypeName: "User",
		Fields: []FactField{
			{Name: "name", Value: FactValue{Type: "string", Value: "Bob"}},
		},
	}

	userID, _ := GenerateFactID(userFact, userType, ctx)
	ctx.RegisterVariable("bob", userID)
	t.Logf("✅ User créé avec ID: %s", userID)

	// Créer Order avec clé composite
	orderFact := Fact{
		TypeName: "Order",
		Fields: []FactField{
			{Name: "user", Value: FactValue{Type: "variableReference", Value: "bob"}},
			{Name: "orderNum", Value: FactValue{Type: "number", Value: float64(1001)}},
			{Name: "total", Value: FactValue{Type: "number", Value: 150.50}},
		},
	}

	orderID, err := GenerateFactID(orderFact, orderType, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur génération ID Order: %v", err)
	}

	// L'ID devrait combiner l'ID de bob + le numéro
	// L'ID de bob contient ~ qui est échappé en %7E
	expectedOrderID := "Order~User%7EBob_1001"
	if orderID != expectedOrderID {
		t.Errorf("❌ ID Order attendu '%s', reçu '%s'", expectedOrderID, orderID)
	}
	t.Logf("✅ ID Order composite généré: %s", orderID)
}

func TestGenerateFactIDFromHash_WithFacts(t *testing.T) {
	t.Log("🧪 TEST GÉNÉRATION ID HASH - AVEC FAITS")
	t.Log("========================================")

	// Type sans clé primaire
	userType := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "name", Type: "string"},
		},
	}

	logType := TypeDefinition{
		Name: "Log",
		Fields: []Field{
			{Name: "user", Type: "User"},
			{Name: "message", Type: "string"},
		},
	}

	ctx := NewFactContext([]TypeDefinition{userType, logType})

	// Créer User (sans PK, génère hash)
	userFact := Fact{
		TypeName: "User",
		Fields: []FactField{
			{Name: "name", Value: FactValue{Type: "string", Value: "Charlie"}},
		},
	}

	userID, err := GenerateFactID(userFact, userType, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur génération ID User: %v", err)
	}

	// L'ID devrait être un hash
	if !strings.HasPrefix(userID, "User~") {
		t.Errorf("❌ ID devrait commencer par 'User~', reçu '%s'", userID)
	}

	hashPart := strings.TrimPrefix(userID, "User~")
	if len(hashPart) != IDHashLength {
		t.Errorf("❌ Hash devrait avoir longueur %d, reçu %d", IDHashLength, len(hashPart))
	}
	t.Logf("✅ ID User hash généré: %s", userID)

	// Enregistrer la variable
	ctx.RegisterVariable("charlie", userID)

	// Créer Log qui référence charlie (sans PK, génère hash)
	logFact := Fact{
		TypeName: "Log",
		Fields: []FactField{
			{Name: "user", Value: FactValue{Type: "variableReference", Value: "charlie"}},
			{Name: "message", Value: FactValue{Type: "string", Value: "Test message"}},
		},
	}

	logID, err := GenerateFactID(logFact, logType, ctx)
	if err != nil {
		t.Fatalf("❌ Erreur génération ID Log: %v", err)
	}

	// L'ID devrait être un hash qui inclut l'ID de charlie
	if !strings.HasPrefix(logID, "Log~") {
		t.Errorf("❌ ID devrait commencer par 'Log~', reçu '%s'", logID)
	}
	t.Logf("✅ ID Log hash généré: %s", logID)

	// Test déterminisme : même fait doit générer même ID
	logID2, _ := GenerateFactID(logFact, logType, ctx)
	if logID != logID2 {
		t.Errorf("❌ Hash non déterministe: '%s' != '%s'", logID, logID2)
	} else {
		t.Log("✅ Hash déterministe vérifié")
	}
}

func TestConvertFactsToReteFormat_WithAssignments(t *testing.T) {
	t.Log("🧪 TEST CONVERSION RETE - AVEC AFFECTATIONS")
	t.Log("============================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: "string", IsPrimaryKey: true},
					{Name: "age", Type: "number"},
				},
			},
			{
				Name: "Login",
				Fields: []Field{
					{Name: "user", Type: "User"},
					{Name: "email", Type: "string", IsPrimaryKey: true},
				},
			},
		},
		FactAssignments: []FactAssignment{
			{
				Variable: "alice",
				Fact: Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
						{Name: "age", Value: FactValue{Type: "number", Value: float64(30)}},
					},
				},
			},
		},
		Facts: []Fact{
			{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
					{Name: "email", Value: FactValue{Type: "string", Value: "alice@example.com"}},
				},
			},
		},
	}

	reteFacts, err := ConvertFactsToReteFormat(program)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	if len(reteFacts) != 2 {
		t.Fatalf("❌ Attendu 2 faits RETE, reçu %d", len(reteFacts))
	}

	// Vérifier le fait User
	userFact := reteFacts[0]
	userID, ok := userFact[FieldNameInternalID].(string)
	if !ok || userID != "User~Alice" {
		t.Errorf("❌ ID User attendu 'User~Alice', reçu '%v'", userID)
	}
	t.Logf("✅ Fait User: ID = %s", userID)

	// Vérifier le fait Login
	loginFact := reteFacts[1]
	loginID, ok := loginFact[FieldNameInternalID].(string)
	if !ok || loginID != "Login~alice@example.com" {
		t.Errorf("❌ ID Login attendu 'Login~alice@example.com', reçu '%v'", loginID)
	}

	// Vérifier que le champ user du Login contient l'ID de alice
	userField, ok := loginFact["user"].(string)
	if !ok || userField != "User~Alice" {
		t.Errorf("❌ Champ user attendu 'User~Alice', reçu '%v'", userField)
	}
	t.Logf("✅ Fait Login: ID = %s, user = %s", loginID, userField)
}

func TestCompleteFlow_FactReferences(t *testing.T) {
	t.Log("🧪 TEST FLOW COMPLET - RÉFÉRENCES DE FAITS")
	t.Log("===========================================")

	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: "string", IsPrimaryKey: true},
					{Name: "age", Type: "number"},
				},
			},
			{
				Name: "Login",
				Fields: []Field{
					{Name: "user", Type: "User"},
					{Name: "email", Type: "string", IsPrimaryKey: true},
					{Name: "password", Type: "string"},
				},
			},
		},
		FactAssignments: []FactAssignment{
			{
				Variable: "alice",
				Fact: Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
						{Name: "age", Value: FactValue{Type: "number", Value: float64(30)}},
					},
				},
			},
			{
				Variable: "bob",
				Fact: Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: "name", Value: FactValue{Type: "string", Value: "Bob"}},
						{Name: "age", Value: FactValue{Type: "number", Value: float64(25)}},
					},
				},
			},
		},
		Facts: []Fact{
			{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
					{Name: "email", Value: FactValue{Type: "string", Value: "alice@ex.com"}},
					{Name: "password", Value: FactValue{Type: "string", Value: "pw1"}},
				},
			},
			{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "bob"}},
					{Name: "email", Value: FactValue{Type: "string", Value: "bob@ex.com"}},
					{Name: "password", Value: FactValue{Type: "string", Value: "pw2"}},
				},
			},
		},
	}

	reteFacts, err := ConvertFactsToReteFormat(program)
	if err != nil {
		t.Fatalf("❌ Erreur de conversion: %v", err)
	}

	// Vérifier que 4 faits sont créés (2 User + 2 Login)
	if len(reteFacts) != 4 {
		t.Errorf("❌ Attendu 4 faits, reçu %d", len(reteFacts))
	}

	// Vérifier les IDs générés
	expectedIDs := map[string]string{
		"User~Alice":         "User",
		"User~Bob":           "User",
		"Login~alice@ex.com": "Login",
		"Login~bob@ex.com":   "Login",
	}

	foundIDs := make(map[string]bool)
	for _, fact := range reteFacts {
		id := fact[FieldNameInternalID].(string)
		expectedType, exists := expectedIDs[id]
		if !exists {
			t.Errorf("❌ ID inattendu: %s", id)
			continue
		}
		if fact[FieldNameReteType] != expectedType {
			t.Errorf("❌ Type attendu '%s' pour ID '%s', reçu '%s'", expectedType, id, fact[FieldNameReteType])
		}
		foundIDs[id] = true
		t.Logf("✅ Fait trouvé: %s (type: %s)", id, expectedType)
	}

	if len(foundIDs) != 4 {
		t.Errorf("❌ Attendu 4 IDs uniques, trouvé %d", len(foundIDs))
	}

	// Vérifier les références dans les Login
	for _, fact := range reteFacts {
		if fact[FieldNameReteType] == "Login" {
			userRef, ok := fact["user"].(string)
			if !ok {
				t.Errorf("❌ Champ user devrait être un string")
				continue
			}
			if userRef != "User~Alice" && userRef != "User~Bob" {
				t.Errorf("❌ Référence user invalide: %s", userRef)
			} else {
				t.Logf("✅ Référence user correcte: %s", userRef)
			}
		}
	}

	t.Log("✅ Flow complet fonctionne correctement")
}

func TestBackwardCompatibility(t *testing.T) {
	t.Log("🧪 TEST RÉTROCOMPATIBILITÉ")
	t.Log("===========================")

	// Test sans contexte (ancienne méthode)
	userType := TypeDefinition{
		Name: "User",
		Fields: []Field{
			{Name: "id", Type: "string", IsPrimaryKey: true},
		},
	}

	userFact := Fact{
		TypeName: "User",
		Fields: []FactField{
			{Name: "id", Value: FactValue{Type: "string", Value: "U001"}},
		},
	}

	// Appel sans contexte (nil) - devrait fonctionner
	id1, err := GenerateFactID(userFact, userType, nil)
	if err != nil {
		t.Fatalf("❌ Erreur avec contexte nil: %v", err)
	}
	t.Logf("✅ GenerateFactID avec ctx=nil: %s", id1)

	// Appel avec GenerateFactIDWithoutContext (deprecated)
	id2, err := GenerateFactIDWithoutContext(userFact, userType)
	if err != nil {
		t.Fatalf("❌ Erreur avec fonction deprecated: %v", err)
	}
	t.Logf("✅ GenerateFactIDWithoutContext: %s", id2)

	// Les deux devraient donner le même résultat
	if id1 != id2 {
		t.Errorf("❌ IDs différents: '%s' != '%s'", id1, id2)
	}

	expectedID := "User~U001"
	if id1 != expectedID {
		t.Errorf("❌ ID attendu '%s', reçu '%s'", expectedID, id1)
	}

	t.Log("✅ Rétrocompatibilité assurée")
}
