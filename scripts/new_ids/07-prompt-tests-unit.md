# Prompt 07 - Migration des Tests Unitaires

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Migrer tous les tests unitaires existants pour supporter la nouvelle gestion des identifiants :

1. **Adapter les tests existants** - Remplacer `id` par `_id_`
2. **Ajouter nouveaux tests** - Couvrir nouvelles fonctionnalités
3. **Maintenir la couverture** - > 80% obligatoire
4. **Tests fonctionnels** - Pas de mocks, résultats réels
5. **Messages clairs** - Émojis et descriptions

---

## 📋 Contexte

### État Actuel

Les tests existants utilisent :
- Champ `id` visible et accessible
- Affectations manuelles d'IDs possibles
- Comparaisons via champs primitifs uniquement
- Pas de tests pour affectations de variables

### État Cible

Les tests doivent vérifier :
- Champ `_id_` caché et inaccessible
- Génération automatique d'IDs obligatoire
- Comparaisons de faits via IDs internes
- Affectations de variables fonctionnelles
- Erreurs claires pour utilisations interdites

---

## 📝 Tâches à Réaliser

### 1. Inventorier les Tests Existants

#### Rechercher Tous les Tests

```bash
# Tests unitaires dans constraint/
find constraint/ -name "*_test.go" -type f | sort

# Tests unitaires dans rete/
find rete/ -name "*_test.go" -type f | sort

# Tests unitaires dans api/ et tsdio/
find api/ tsdio/ -name "*_test.go" -type f | sort

# Compter le nombre total de tests
grep -r "^func Test" constraint/ rete/ api/ tsdio/ --include="*_test.go" | wc -l
```

**Créer un rapport** : `REPORTS/new_ids_tests_inventory.md`

```markdown
# Inventaire des Tests - Migration IDs

## Tests par Module

### constraint/
- constraint_test.go : XX tests
- id_generator_test.go : XX tests
- primary_key_validation_test.go : XX tests
- [...]

### rete/
- [...]

### api/
- [...]

### tsdio/
- [...]

## Total
- Nombre de fichiers : XXX
- Nombre de tests : XXX
- Estimation : XXX heures de migration
```

### 2. Migrer les Tests de Génération d'IDs

#### Fichier : `constraint/id_generator_test.go`

**Modifications nécessaires** :

1. **Constantes** : Utiliser `FieldNameInternalID`
2. **Contexte** : Ajouter `FactContext` aux appels
3. **Format IDs** : Vérifier nouveau format

**Avant** :
```go
func TestGenerateFactID(t *testing.T) {
    fact := Fact{
        TypeName: "User",
        Fields: []FactField{
            {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
        },
    }
    
    typeDef := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    }
    
    id, err := GenerateFactID(fact, typeDef)
    // ...
}
```

**Après** :
```go
func TestGenerateFactID(t *testing.T) {
    t.Log("🧪 TEST GÉNÉRATION ID")
    t.Log("====================")
    
    fact := Fact{
        TypeName: "User",
        Fields: []FactField{
            {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
        },
    }
    
    typeDef := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    }
    
    // Créer le contexte
    ctx := NewFactContext([]TypeDefinition{typeDef})
    
    // Générer l'ID avec le contexte
    id, err := GenerateFactID(fact, typeDef, ctx)
    if err != nil {
        t.Fatalf("❌ Erreur inattendue: %v", err)
    }
    
    expectedID := "User~Alice"
    if id != expectedID {
        t.Errorf("❌ ID attendu '%s', reçu '%s'", expectedID, id)
    } else {
        t.Logf("✅ ID généré correctement: %s", id)
    }
}
```

**Nouveaux tests à ajouter** :

```go
func TestGenerateFactID_WithVariableReference(t *testing.T) {
    t.Log("🧪 TEST GÉNÉRATION ID - RÉFÉRENCE VARIABLE")
    t.Log("===========================================")
    
    // Définir les types
    userType := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    }
    
    loginType := TypeDefinition{
        Name: "Login",
        Fields: []Field{
            {Name: "user", Type: "User"},
            {Name: "email", Type: "string", IsPrimaryKey: true},
        },
    }
    
    ctx := NewFactContext([]TypeDefinition{userType, loginType})
    
    // Créer et générer ID pour User
    userFact := Fact{
        TypeName: "User",
        Fields: []FactField{
            {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
        },
    }
    
    userID, err := GenerateFactID(userFact, userType, ctx)
    if err != nil {
        t.Fatalf("❌ Erreur génération User: %v", err)
    }
    t.Logf("✅ User ID: %s", userID)
    
    // Enregistrer la variable
    ctx.RegisterVariable("alice", userID)
    
    // Créer Login avec référence
    loginFact := Fact{
        TypeName: "Login",
        Fields: []FactField{
            {Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
            {Name: "email", Value: FactValue{Type: "string", Value: "alice@example.com"}},
        },
    }
    
    loginID, err := GenerateFactID(loginFact, loginType, ctx)
    if err != nil {
        t.Fatalf("❌ Erreur génération Login: %v", err)
    }
    
    expectedLoginID := "Login~alice@example.com"
    if loginID != expectedLoginID {
        t.Errorf("❌ Login ID attendu '%s', reçu '%s'", expectedLoginID, loginID)
    } else {
        t.Logf("✅ Login ID généré: %s", loginID)
    }
}

func TestGenerateFactID_VariableNotDefined(t *testing.T) {
    t.Log("🧪 TEST GÉNÉRATION ID - VARIABLE NON DÉFINIE")
    t.Log("=============================================")
    
    loginType := TypeDefinition{
        Name: "Login",
        Fields: []Field{
            {Name: "user", Type: "User"},
            {Name: "email", Type: "string", IsPrimaryKey: true},
        },
    }
    
    ctx := NewFactContext([]TypeDefinition{loginType})
    
    // Créer Login avec variable non définie
    loginFact := Fact{
        TypeName: "Login",
        Fields: []FactField{
            {Name: "user", Value: FactValue{Type: "variableReference", Value: "unknown"}},
            {Name: "email", Value: FactValue{Type: "string", Value: "test@example.com"}},
        },
    }
    
    _, err := GenerateFactID(loginFact, loginType, ctx)
    if err == nil {
        t.Error("❌ Attendu une erreur pour variable non définie")
    } else {
        t.Logf("✅ Erreur attendue: %v", err)
    }
}
```

### 3. Migrer les Tests de Validation

#### Fichier : `constraint/primary_key_validation_test.go`

**Modifications** :

1. **Interdire `_id_`** : Vérifier rejet du champ `_id_`
2. **Messages d'erreur** : Mettre à jour les messages attendus

**Nouveaux tests** :

```go
func TestValidateFactPrimaryKey_InternalIDForbidden(t *testing.T) {
    t.Log("🧪 TEST VALIDATION - _id_ INTERDIT")
    t.Log("===================================")
    
    typeDef := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    }
    
    tests := []struct {
        name    string
        fact    Fact
        wantErr bool
        errMsg  string
    }{
        {
            name: "_id_ dans les champs",
            fact: Fact{
                TypeName: "User",
                Fields: []FactField{
                    {Name: FieldNameInternalID, Value: FactValue{Type: "string", Value: "manual"}},
                    {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
                },
            },
            wantErr: true,
            errMsg:  "réservé",
        },
        {
            name: "fait valide sans _id_",
            fact: Fact{
                TypeName: "User",
                Fields: []FactField{
                    {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
                },
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateFactPrimaryKey(tt.fact, typeDef)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("❌ Attendu une erreur, reçu nil")
                } else {
                    if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
                        t.Errorf("❌ Message attendu contenant '%s', reçu: %v", tt.errMsg, err)
                    } else {
                        t.Logf("✅ Erreur attendue: %v", err)
                    }
                }
            } else {
                if err != nil {
                    t.Errorf("❌ Erreur inattendue: %v", err)
                } else {
                    t.Logf("✅ Validation réussie")
                }
            }
        })
    }
}
```

### 4. Migrer les Tests de Conversion de Faits

#### Fichier : `constraint/constraint_facts_test.go` (si existe)

**Modifications** :

1. **Contexte** : Ajouter support `FactContext`
2. **Affectations** : Tester conversion avec affectations
3. **Références** : Tester résolution de variables

**Nouveaux tests** :

```go
func TestConvertFactsToReteFormat_WithAssignments(t *testing.T) {
    t.Log("🧪 TEST CONVERSION - AVEC AFFECTATIONS")
    t.Log("=======================================")
    
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
                        {Name: "age", Value: FactValue{Type: "number", Value: 30.0}},
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
                },
            },
        },
    }
    
    reteFacts, err := ConvertFactsToReteFormat(program)
    if err != nil {
        t.Fatalf("❌ Erreur de conversion: %v", err)
    }
    
    // Vérifier nombre de faits
    if len(reteFacts) != 2 {
        t.Fatalf("❌ Attendu 2 faits, reçu %d", len(reteFacts))
    }
    
    // Vérifier le fait User
    userFact := reteFacts[0]
    userID, ok := userFact[FieldNameInternalID].(string)
    if !ok {
        t.Fatal("❌ ID User manquant ou invalide")
    }
    
    if userID != "User~Alice" {
        t.Errorf("❌ User ID attendu 'User~Alice', reçu '%s'", userID)
    } else {
        t.Logf("✅ User ID: %s", userID)
    }
    
    // Vérifier le fait Login
    loginFact := reteFacts[1]
    loginID, ok := loginFact[FieldNameInternalID].(string)
    if !ok {
        t.Fatal("❌ Login ID manquant ou invalide")
    }
    
    if loginID != "Login~alice@ex.com" {
        t.Errorf("❌ Login ID attendu 'Login~alice@ex.com', reçu '%s'", loginID)
    } else {
        t.Logf("✅ Login ID: %s", loginID)
    }
    
    // Vérifier que le champ user contient l'ID de alice
    userField, ok := loginFact["user"].(string)
    if !ok {
        t.Fatal("❌ Champ user manquant ou invalide")
    }
    
    if userField != "User~Alice" {
        t.Errorf("❌ Champ user attendu 'User~Alice', reçu '%s'", userField)
    } else {
        t.Logf("✅ Champ user résolu: %s", userField)
    }
}

func TestConvertFactsToReteFormat_InternalIDForbidden(t *testing.T) {
    t.Log("🧪 TEST CONVERSION - _id_ INTERDIT")
    t.Log("===================================")
    
    program := Program{
        Types: []TypeDefinition{
            {
                Name: "User",
                Fields: []Field{
                    {Name: "name", Type: "string", IsPrimaryKey: true},
                },
            },
        },
        Facts: []Fact{
            {
                TypeName: "User",
                Fields: []FactField{
                    {Name: FieldNameInternalID, Value: FactValue{Type: "string", Value: "manual"}},
                    {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
                },
            },
        },
    }
    
    _, err := ConvertFactsToReteFormat(program)
    if err == nil {
        t.Error("❌ Attendu une erreur pour _id_ manuel")
    } else {
        t.Logf("✅ Erreur attendue: %v", err)
    }
}
```

### 5. Migrer les Tests de Parsing

#### Fichier : `constraint/parser_test.go` (ou fichiers de parsing)

**Modifications** :

1. **Affectations** : Tester parsing de `alice = User(...)`
2. **Variables** : Tester parsing de références
3. **Interdictions** : Tester rejet de `_id_`

**Nouveaux tests** :

```go
func TestParseFactAssignment(t *testing.T) {
    t.Log("🧪 TEST PARSE - AFFECTATION")
    t.Log("============================")
    
    tests := []struct {
        name         string
        input        string
        wantVariable string
        wantType     string
        wantErr      bool
    }{
        {
            name:         "affectation simple",
            input:        `alice = User("Alice", 30)`,
            wantVariable: "alice",
            wantType:     "User",
            wantErr:      false,
        },
        {
            name:         "affectation avec underscore",
            input:        `user_1 = User("Bob", 25)`,
            wantVariable: "user_1",
            wantType:     "User",
            wantErr:      false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            program, err := ParseProgram(tt.input)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("❌ Attendu une erreur, reçu nil")
                } else {
                    t.Logf("✅ Erreur attendue: %v", err)
                }
                return
            }
            
            if err != nil {
                t.Fatalf("❌ Erreur de parsing: %v", err)
            }
            
            if len(program.FactAssignments) != 1 {
                t.Fatalf("❌ Attendu 1 affectation, reçu %d", len(program.FactAssignments))
            }
            
            assignment := program.FactAssignments[0]
            if assignment.Variable != tt.wantVariable {
                t.Errorf("❌ Variable attendue '%s', reçu '%s'", tt.wantVariable, assignment.Variable)
            }
            
            if assignment.Fact.TypeName != tt.wantType {
                t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.wantType, assignment.Fact.TypeName)
            }
            
            t.Logf("✅ Affectation parsée: %s = %s", assignment.Variable, assignment.Fact.TypeName)
        })
    }
}

func TestParseFieldAccess_InternalIDForbidden(t *testing.T) {
    t.Log("🧪 TEST PARSE - _id_ INTERDIT")
    t.Log("==============================")
    
    input := `{u: User} / u._id_ == "test" ==> Log("test")`
    
    _, err := ParseProgram(input)
    
    // Selon l'implémentation, peut être rejeté au parsing ou à la validation
    if err != nil {
        t.Logf("✅ _id_ rejeté au parsing: %v", err)
    } else {
        t.Log("⚠️  _id_ accepté au parsing, sera rejeté à la validation")
    }
}
```

### 6. Migrer les Tests RETE

#### Fichier : `rete/fact_token_test.go` (ou similaires)

**Modifications** :

1. **Structure Fact** : Utiliser nouveau tag JSON `_id_`
2. **Résolveur** : Tester `FieldResolver`
3. **Comparaisons** : Tester `ComparisonEvaluator`

**Nouveaux tests** :

```go
func TestFact_InternalIDNotExposed(t *testing.T) {
    t.Log("🧪 TEST FACT - ID INTERNE NON EXPOSÉ")
    t.Log("=====================================")
    
    fact := &Fact{
        ID:   "User~Alice",
        Type: "User",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    }
    
    // Sérialiser en JSON
    jsonData, err := json.Marshal(fact)
    if err != nil {
        t.Fatalf("❌ Erreur de sérialisation: %v", err)
    }
    
    jsonStr := string(jsonData)
    t.Logf("JSON: %s", jsonStr)
    
    // Vérifier que ID n'est pas dans le JSON avec le tag 'id'
    var parsed map[string]interface{}
    if err := json.Unmarshal(jsonData, &parsed); err != nil {
        t.Fatalf("❌ Erreur de désérialisation: %v", err)
    }
    
    // Vérifier que _id_ est dans le JSON (nouveau tag)
    if _, exists := parsed["_id_"]; !exists {
        t.Error("❌ _id_ devrait être dans le JSON (tag _id_)")
    }
    
    // Mais pas avec l'ancien tag 'id'
    if _, exists := parsed["id"]; exists {
        t.Error("❌ 'id' ne devrait plus être dans le JSON")
    }
    
    t.Log("✅ Sérialisation correcte avec tag _id_")
}
```

### 7. Créer Tests pour Nouvelles Fonctionnalités

#### Fichier : `constraint/type_system_test.go` (déjà créé dans prompt 05)

Vérifier que tous les tests du prompt 05 sont présents :

- [ ] `TestTypeSystem_TypeChecks`
- [ ] `TestTypeSystem_GetFieldType`
- [ ] `TestTypeSystem_Variables`
- [ ] `TestTypeSystem_CircularReferences`
- [ ] `TestTypeSystem_TypeCompatibility`

#### Fichier : `constraint/fact_validator_test.go` (déjà créé dans prompt 05)

Vérifier :

- [ ] `TestFactValidator_ValidateFact`

#### Fichier : `constraint/program_validator_test.go` (déjà créé dans prompt 05)

Vérifier :

- [ ] `TestProgramValidator_ValidProgram`
- [ ] `TestProgramValidator_InvalidPrograms`

#### Fichier : `rete/field_resolver_test.go` (déjà créé dans prompt 04)

Vérifier :

- [ ] `TestFieldResolver_ResolveFieldValue`
- [ ] `TestFieldResolver_ResolveFactID`

#### Fichier : `rete/comparison_evaluator_test.go` (déjà créé dans prompt 04)

Vérifier :

- [ ] `TestComparisonEvaluator_CompareFactIDs`
- [ ] `TestComparisonEvaluator_ComparePrimitives`
- [ ] `TestComparisonEvaluator_EvaluateComparison`

### 8. Vérifier la Couverture de Code

#### Commandes de Couverture

```bash
# Couverture par module
go test ./constraint -cover -coverprofile=constraint_coverage.out
go test ./rete -cover -coverprofile=rete_coverage.out
go test ./api -cover -coverprofile=api_coverage.out
go test ./tsdio -cover -coverprofile=tsdio_coverage.out

# Afficher le rapport détaillé
go tool cover -func=constraint_coverage.out
go tool cover -func=rete_coverage.out
go tool cover -func=api_coverage.out
go tool cover -func=tsdio_coverage.out

# Générer rapport HTML
go tool cover -html=constraint_coverage.out -o constraint_coverage.html
go tool cover -html=rete_coverage.out -o rete_coverage.html
go tool cover -html=api_coverage.out -o api_coverage.html
go tool cover -html=tsdio_coverage.out -o tsdio_coverage.html
```

**Objectif** : > 80% de couverture dans chaque module

#### Identifier les Zones Non Couvertes

```bash
# Trouver les fonctions avec couverture < 80%
go tool cover -func=constraint_coverage.out | awk '$3 < 80.0 {print $1, $2, $3}'
```

**Ajouter des tests** pour les zones non couvertes.

### 9. Créer un Script de Validation des Tests

#### Fichier : `scripts/validate-tests.sh` (nouveau)

```bash
#!/bin/bash
# Script de validation des tests pour la migration des IDs

set -e

echo "🧪 VALIDATION DES TESTS - MIGRATION IDs"
echo "========================================"
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

# 1. Vérifier que tous les tests passent
log_info "Étape 1/5: Exécution de tous les tests..."
if go test ./... -v > test_output.log 2>&1; then
    log_success "Tous les tests passent"
else
    log_error "Des tests échouent. Voir test_output.log"
    exit 1
fi

# 2. Vérifier la couverture
log_info "Étape 2/5: Vérification de la couverture..."
go test ./constraint -cover -coverprofile=constraint_coverage.out > /dev/null 2>&1
go test ./rete -cover -coverprofile=rete_coverage.out > /dev/null 2>&1
go test ./api -cover -coverprofile=api_coverage.out > /dev/null 2>&1
go test ./tsdio -cover -coverprofile=tsdio_coverage.out > /dev/null 2>&1

# Extraire le pourcentage de couverture
constraint_coverage=$(go tool cover -func=constraint_coverage.out | grep total | awk '{print $3}' | sed 's/%//')
rete_coverage=$(go tool cover -func=rete_coverage.out | grep total | awk '{print $3}' | sed 's/%//')
api_coverage=$(go tool cover -func=api_coverage.out | grep total | awk '{print $3}' | sed 's/%//')
tsdio_coverage=$(go tool cover -func=tsdio_coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo ""
echo "Couverture de code:"
echo "  - constraint: ${constraint_coverage}%"
echo "  - rete: ${rete_coverage}%"
echo "  - api: ${api_coverage}%"
echo "  - tsdio: ${tsdio_coverage}%"
echo ""

# Vérifier que chaque module a > 80%
if (( $(echo "$constraint_coverage < 80.0" | bc -l) )); then
    log_error "Couverture constraint < 80% (${constraint_coverage}%)"
    exit 1
fi

if (( $(echo "$rete_coverage < 80.0" | bc -l) )); then
    log_error "Couverture rete < 80% (${rete_coverage}%)"
    exit 1
fi

if (( $(echo "$api_coverage < 80.0" | bc -l) )); then
    log_error "Couverture api < 80% (${api_coverage}%)"
    exit 1
fi

if (( $(echo "$tsdio_coverage < 80.0" | bc -l) )); then
    log_error "Couverture tsdio < 80% (${tsdio_coverage}%)"
    exit 1
fi

log_success "Couverture > 80% dans tous les modules"

# 3. Vérifier absence de _id_ manuel dans les tests
log_info "Étape 3/5: Vérification de l'usage de _id_..."
if grep -r "\"id\":" constraint/ rete/ api/ tsdio/ --include="*_test.go" | grep -v "// Allowed" | grep -v "FieldNameInternalID"; then
    log_error "Utilisation de \"id\": trouvée dans les tests (devrait être FieldNameInternalID)"
    exit 1
fi
log_success "Aucune utilisation incorrecte de 'id'"

# 4. Vérifier que FieldNameInternalID est utilisé
log_info "Étape 4/5: Vérification de l'usage de FieldNameInternalID..."
if ! grep -r "FieldNameInternalID" constraint/ --include="*_test.go" > /dev/null; then
    log_error "FieldNameInternalID non utilisé dans les tests constraint/"
    exit 1
fi
log_success "FieldNameInternalID utilisé correctement"

# 5. Vérifier présence de tests pour nouvelles fonctionnalités
log_info "Étape 5/5: Vérification des tests de nouvelles fonctionnalités..."

# Tests d'affectations
if ! grep -r "TestFactAssignment" constraint/ --include="*_test.go" > /dev/null; then
    log_error "Tests d'affectations manquants"
    exit 1
fi

# Tests de résolution de variables
if ! grep -r "TestGenerateFactID.*Variable" constraint/ --include="*_test.go" > /dev/null; then
    log_error "Tests de résolution de variables manquants"
    exit 1
fi

# Tests de comparaison de faits
if ! grep -r "TestComparison.*Fact" rete/ --include="*_test.go" > /dev/null; then
    log_error "Tests de comparaison de faits manquants"
    exit 1
fi

log_success "Tests de nouvelles fonctionnalités présents"

# Résumé
echo ""
echo "=========================================="
log_success "VALIDATION RÉUSSIE"
echo "=========================================="
echo ""
echo "Statistiques:"
echo "  - Tests passés: $(grep -c "PASS" test_output.log || echo "N/A")"
echo "  - Couverture moyenne: $(echo "scale=2; ($constraint_coverage + $rete_coverage + $api_coverage + $tsdio_coverage) / 4" | bc)%"
echo ""

# Nettoyer
rm -f constraint_coverage.out rete_coverage.out api_coverage.out tsdio_coverage.out
```

**Rendre exécutable** :
```bash
chmod +x scripts/validate-tests.sh
```

### 10. Créer un Rapport de Migration

#### Fichier : `REPORTS/new_ids_tests_migration.md`

```markdown
# Rapport de Migration des Tests - Nouvelle Gestion IDs

Date: [DATE]

## Résumé

### Tests Migrés
- constraint/ : XX fichiers, XXX tests
- rete/ : XX fichiers, XXX tests
- api/ : XX fichiers, XXX tests
- tsdio/ : XX fichiers, XXX tests

### Nouveaux Tests Ajoutés
- Tests d'affectations : XX
- Tests de résolution de variables : XX
- Tests de comparaison de faits : XX
- Tests de validation : XX

### Total
- Anciens tests : XXX
- Nouveaux tests : XXX
- Tests migrés : XXX
- Total final : XXX

## Couverture de Code

### Avant Migration
- constraint: XX%
- rete: XX%
- api: XX%
- tsdio: XX%

### Après Migration
- constraint: XX% (✅ > 80%)
- rete: XX% (✅ > 80%)
- api: XX% (✅ > 80%)
- tsdio: XX% (✅ > 80%)

## Modifications Principales

### 1. Constantes
- Remplacement de `FieldNameID` par `FieldNameInternalID`
- Utilisation de `"_id_"` au lieu de `"id"`

### 2. Contexte
- Ajout de `FactContext` à tous les appels de génération d'ID
- Tests de résolution de variables

### 3. Validation
- Tests d'interdiction de `_id_` manuel
- Tests de validation de types de faits

### 4. Comparaisons
- Tests de comparaisons de faits via IDs
- Tests de compatibilité de types

## Problèmes Rencontrés

### Problème 1: [Description]
**Solution**: [Solution appliquée]

### Problème 2: [Description]
**Solution**: [Solution appliquée]

## Tests Critiques

### Tests Ajoutés
1. `TestGenerateFactID_WithVariableReference` - Génération ID avec variables
2. `TestValidateFactPrimaryKey_InternalIDForbidden` - Interdiction _id_
3. `TestConvertFactsToReteFormat_WithAssignments` - Conversion avec affectations
4. `TestComparisonEvaluator_CompareFactIDs` - Comparaison de faits
5. [...]

### Tests Modifiés
1. `TestGenerateFactID` - Ajout contexte
2. `TestConvertFactsToReteFormat` - Support affectations
3. [...]

## Validation

### Commandes Exécutées
```bash
go test ./... -v
go test ./constraint -cover
go test ./rete -cover
scripts/validate-tests.sh
make test-complete
```

### Résultats
- ✅ Tous les tests passent
- ✅ Couverture > 80% partout
- ✅ Aucune régression détectée
- ✅ Nouvelles fonctionnalités couvertes

## Recommandations

1. Maintenir la couverture > 80%
2. Ajouter tests pour cas limites si identifiés
3. Documenter les tests complexes
4. Utiliser émojis pour lisibilité

## Conclusion

Migration réussie. Tous les tests sont à jour et la couverture est maintenue.
```

---

## ✅ Critères de Succès

### Tests

```bash
# Tous les tests passent
go test ./... -v

# Couverture > 80%
go test ./constraint -cover
go test ./rete -cover
go test ./api -cover
go test ./tsdio -cover

# Script de validation
scripts/validate-tests.sh
```

### Checklist

- [ ] Tous les anciens tests migrés
- [ ] Nouveaux tests ajoutés
- [ ] Couverture > 80% partout
- [ ] Constantes utilisées partout
- [ ] Contexte ajouté aux générateurs
- [ ] Affectations testées
- [ ] Comparaisons testées
- [ ] Validation testée
- [ ] Messages clairs avec émojis
- [ ] Rapport de migration créé

### Validation

```bash
make test-unit
make test-coverage
make validate
```

---

## 📊 Métriques Attendues

### Couverture Minimale
- constraint/ : > 80%
- rete/ : > 80%
- api/ : > 80%
- tsdio/ : > 80%

### Nombre de Tests
- Avant migration : ~XXX tests
- Après migration : ~XXX tests (augmentation de XX%)

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Inventorier tests existants
2. ✅ Migrer tests de génération d'IDs
3. ✅ Migrer tests de validation
4. ✅ Migrer tests de conversion
5. ✅ Migrer tests de parsing
6. ✅ Migrer tests RETE
7. ✅ Créer tests pour nouvelles fonctionnalités
8. ✅ Vérifier couverture
9. ✅ Créer script de validation
10. ✅ Générer rapport de migration

### Commandes

```bash
# Lister tous les tests
find . -name "*_test.go" -type f | sort

# Exécuter les tests
go test ./... -v

# Vérifier la couverture
go test ./constraint -cover
go test ./rete -cover
go test ./api -cover
go test ./tsdio -cover

# Valider
scripts/validate-tests.sh

# Générer rapport HTML
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## 📚 Références

- `scripts/new_ids/06-prompt-api-tsdio.md` - API
- `scripts/new_ids/05-prompt-types-validation.md` - Validation
- `scripts/new_ids/04-prompt-evaluation.md` - Évaluation
- `.github/prompts/common.md` - Standards tests

---

## 📝 Notes

### Points d'Attention

1. **Tests fonctionnels** : Pas de mocks, résultats réels obligatoires

2. **Messages clairs** : Utiliser émojis (✅ ❌ ⚠️) pour lisibilité

3. **Couverture** : Ne pas sacrifier la qualité pour atteindre 80%

4. **Isolation** : Tests doivent être indépendants

### Bonnes Pratiques

```go
// ✅ BON - Test avec émojis et logs clairs
func TestFeature(t *testing.T) {
    t.Log("🧪 TEST FEATURE")
    t.Log("===============")
    
    result, err := Feature()
    if err != nil {
        t.Fatalf("❌ Erreur: %v", err)
    }
    
    if result != expected {
        t.Errorf("❌ Attendu %v, reçu %v", expected, result)
    }
    
    t.Log("✅ Test réussi")
}

// ❌ MAUVAIS - Test sans logs ni émojis
func TestFeature(t *testing.T) {
    result, _ := Feature()
    if result != expected {
        t.Error("failed")
    }
}
```

---

## 🎯 Résultat Attendu

Après ce prompt :

```bash
# Tous les tests passent
$ go test ./...
ok      github.com/resinsec/tsd/constraint  X.XXXs  coverage: XX.X% of statements
ok      github.com/resinsec/tsd/rete        X.XXXs  coverage: XX.X% of statements
ok      github.com/resinsec/tsd/api         X.XXXs  coverage: XX.X% of statements
ok      github.com/resinsec/tsd/tsdio       X.XXXs  coverage: XX.X% of statements

# Validation réussie
$ scripts/validate-tests.sh
🧪 VALIDATION DES TESTS - MIGRATION IDs
========================================

ℹ️  Étape 1/5: Exécution de tous les tests...
✅ Tous les tests passent
ℹ️  Étape 2/5: Vérification de la couverture...

Couverture de code:
  - constraint: 85.2%
  - rete: 82.7%
  - api: 81.4%
  - tsdio: 83.9%

✅ Couverture > 80% dans tous les modules
ℹ️  Étape 3/5: Vérification de l'usage de _id_...
✅ Aucune utilisation incorrecte de 'id'
ℹ️  Étape 4/5: Vérification de l'usage de FieldNameInternalID...
✅ FieldNameInternalID utilisé correctement
ℹ️  Étape 5/5: Vérification des tests de nouvelles fonctionnalités...
✅ Tests de nouvelles fonctionnalités présents

==========================================
✅ VALIDATION RÉUSSIE
==========================================
```

---

**Prompt suivant** : `08-prompt-tests-integration.md`

**Durée estimée** : 6-8 heures

**Complexité** : 🔴 Élevée (migration exhaustive)