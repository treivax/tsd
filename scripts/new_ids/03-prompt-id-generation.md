# Prompt 03 - Génération d'IDs avec Types de Faits

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Adapter la génération d'IDs pour gérer les champs de type fait (User, Login, etc.) en plus des types primitifs.

Les IDs doivent être générés en utilisant les IDs des faits référencés pour les champs non-primitifs.

---

## 📋 Contexte

### État Actuel

```go
// Génération actuelle - uniquement types primitifs
type User(#name: string, age: number)
User("Alice", 30) → ID: "User~Alice"

// Types primitifs dans les clés primaires
type Product(#category: string, #sku: string, price: number)
Product("Electronics", "LAPTOP-001", 1200) → ID: "Product~Electronics_LAPTOP-001"
```

### État Cible

```go
// Génération avec types de faits
type User(#name: string, age: number)
type Login(user: User, #email: string, password: string)

alice = User("Alice", 30)          → ID interne: "User~Alice"
Login(alice, "alice@ex.com", "pw") → ID interne: "Login~User~Alice_alice@ex.com"
                                      (utilise l'ID de alice dans l'ID de Login)
```

---

## 📝 Tâches à Réaliser

### 1. Analyser la Génération Actuelle

#### Fichier : `constraint/id_generator.go`

**Fonctions existantes** :

```go
// GenerateFactID génère l'ID d'un fait
func GenerateFactID(fact Fact, typeDef TypeDefinition) (string, error)

// generateIDFromPrimaryKey génère l'ID depuis les clés primaires
func generateIDFromPrimaryKey(fact Fact, typeDef TypeDefinition) (string, error)

// generateIDFromHash génère l'ID par hash
func generateIDFromHash(fact Fact, typeDef TypeDefinition) (string, error)

// escapeIDValue échappe les caractères spéciaux
func escapeIDValue(value string) string
```

**Comprendre** :
1. Comment les valeurs de clés primaires sont extraites
2. Comment les valeurs sont converties en string
3. Comment les IDs composites sont construits
4. Format actuel : `TypeName~value1_value2`

### 2. Résoudre les Variables de Faits

#### Nouvelle Fonction : Résolution de Variables

**Problème** : Quand on a `Login(alice, "email", "pw")`, il faut :
1. Résoudre `alice` → `User("Alice", 30)`
2. Obtenir l'ID de `alice` → `"User~Alice"`
3. Utiliser cet ID pour générer l'ID de Login

**Solution** : Créer un contexte de résolution

```go
// FactContext contient le contexte pour la génération d'IDs
// Il permet de résoudre les références de variables vers leurs IDs
type FactContext struct {
    // VariableIDs map les noms de variables vers les IDs de faits
    VariableIDs map[string]string
    
    // TypeMap map les noms de types vers leurs définitions
    TypeMap map[string]TypeDefinition
}

// NewFactContext crée un nouveau contexte de génération d'IDs
func NewFactContext(types []TypeDefinition) *FactContext {
    typeMap := make(map[string]TypeDefinition)
    for _, t := range types {
        typeMap[t.Name] = t
    }
    
    return &FactContext{
        VariableIDs: make(map[string]string),
        TypeMap:     typeMap,
    }
}

// RegisterVariable enregistre l'ID d'une variable de fait
func (fc *FactContext) RegisterVariable(varName, factID string) {
    fc.VariableIDs[varName] = factID
}

// ResolveVariable résout une variable vers son ID
func (fc *FactContext) ResolveVariable(varName string) (string, error) {
    id, exists := fc.VariableIDs[varName]
    if !exists {
        return "", fmt.Errorf("variable '%s' non définie dans le contexte", varName)
    }
    return id, nil
}
```

#### Tests

```go
func TestFactContext(t *testing.T) {
    t.Log("🧪 TEST FACT CONTEXT")
    t.Log("====================")
    
    ctx := NewFactContext([]TypeDefinition{
        {Name: "User", Fields: []Field{{Name: "name", Type: "string", IsPrimaryKey: true}}},
    })
    
    // Enregistrer une variable
    ctx.RegisterVariable("alice", "User~Alice")
    
    // Résoudre la variable
    id, err := ctx.ResolveVariable("alice")
    if err != nil {
        t.Fatalf("❌ Erreur de résolution: %v", err)
    }
    
    if id != "User~Alice" {
        t.Errorf("❌ ID attendu 'User~Alice', reçu '%s'", id)
    }
    
    // Variable non définie
    _, err = ctx.ResolveVariable("bob")
    if err == nil {
        t.Error("❌ Attendu une erreur pour variable non définie")
    }
    
    t.Log("✅ Contexte fonctionne correctement")
}
```

### 3. Modifier la Génération d'IDs

#### Mise à Jour de la Signature

**Avant** :
```go
func GenerateFactID(fact Fact, typeDef TypeDefinition) (string, error)
```

**Après** :
```go
func GenerateFactID(fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error)
```

**Note** : Maintenir une version sans contexte pour backward compatibility temporaire :

```go
// GenerateFactIDWithoutContext generates an ID without variable resolution
// Deprecated: Use GenerateFactID with FactContext instead
func GenerateFactIDWithoutContext(fact Fact, typeDef TypeDefinition) (string, error) {
    return GenerateFactID(fact, typeDef, NewFactContext(nil))
}
```

#### Modification de `generateIDFromPrimaryKey`

**Nouvelle logique** :

```go
// generateIDFromPrimaryKey génère l'ID depuis les clés primaires
// Gère les types primitifs ET les références de faits
func generateIDFromPrimaryKey(fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
    if !typeDef.HasPrimaryKey() {
        return "", errors.New("le type n'a pas de clé primaire définie")
    }
    
    // Récupérer les champs de clé primaire
    pkFields := typeDef.GetPrimaryKeyFields()
    
    var pkValues []string
    for _, pkField := range pkFields {
        // Trouver la valeur du champ dans le fait
        value, err := getFactFieldValue(fact, pkField.Name)
        if err != nil {
            return "", fmt.Errorf("champ de clé primaire '%s': %v", pkField.Name, err)
        }
        
        // Convertir la valeur en string selon son type
        strValue, err := convertFieldValueToString(value, pkField, ctx)
        if err != nil {
            return "", fmt.Errorf("conversion du champ '%s': %v", pkField.Name, err)
        }
        
        // Échapper les caractères spéciaux
        escapedValue := escapeIDValue(strValue)
        pkValues = append(pkValues, escapedValue)
    }
    
    // Construire l'ID : TypeName~value1_value2_...
    if len(pkValues) == 1 {
        return fmt.Sprintf("%s%s%s", typeDef.Name, IDSeparatorType, pkValues[0]), nil
    }
    
    return fmt.Sprintf("%s%s%s", typeDef.Name, IDSeparatorType, 
        strings.Join(pkValues, IDSeparatorValue)), nil
}
```

#### Nouvelle Fonction : Conversion de Valeurs

```go
// convertFieldValueToString convertit une valeur de champ en string
// Gère les types primitifs ET les références de faits
func convertFieldValueToString(value FactValue, field Field, ctx *FactContext) (string, error) {
    switch value.Type {
    case ValueTypeString:
        // Type primitif string
        str, ok := value.Value.(string)
        if !ok {
            return "", fmt.Errorf("valeur string attendue, reçu %T", value.Value)
        }
        return str, nil
        
    case ValueTypeNumber:
        // Type primitif number
        num, ok := value.Value.(float64)
        if !ok {
            return "", fmt.Errorf("valeur number attendue, reçu %T", value.Value)
        }
        return formatNumber(num), nil
        
    case ValueTypeBoolean, ValueTypeBool:
        // Type primitif boolean
        b, ok := value.Value.(bool)
        if !ok {
            return "", fmt.Errorf("valeur boolean attendue, reçu %T", value.Value)
        }
        if b {
            return "true", nil
        }
        return "false", nil
        
    case "variableReference":
        // Référence de variable (nouveau)
        varName, ok := value.Value.(string)
        if !ok {
            return "", fmt.Errorf("nom de variable attendu, reçu %T", value.Value)
        }
        
        // Résoudre la variable vers son ID
        if ctx == nil {
            return "", errors.New("contexte requis pour résoudre les variables")
        }
        
        id, err := ctx.ResolveVariable(varName)
        if err != nil {
            return "", fmt.Errorf("résolution de variable '%s': %v", varName, err)
        }
        
        // Utiliser l'ID du fait référencé
        return id, nil
        
    default:
        return "", fmt.Errorf("type de valeur non supporté: %s", value.Type)
    }
}

// formatNumber formate un nombre pour l'ID (sans .0 pour les entiers)
func formatNumber(n float64) string {
    if n == float64(int64(n)) {
        return fmt.Sprintf("%d", int64(n))
    }
    return fmt.Sprintf("%g", n)
}
```

#### Tests de Conversion

```go
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
                t.Logf("✅ Conversion correcte: %s", got)
            }
        })
    }
}
```

### 4. Gérer les IDs Composites avec Faits

#### Tests d'Intégration

```go
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
    
    // L'ID devrait utiliser l'ID de alice dans sa clé primaire
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
    expectedOrderID := "Order~User~Bob_1001"
    if orderID != expectedOrderID {
        t.Errorf("❌ ID Order attendu '%s', reçu '%s'", expectedOrderID, orderID)
    }
    t.Logf("✅ ID Order composite généré: %s", orderID)
}
```

### 5. Modifier la Conversion de Faits en Format RETE

#### Fichier : `constraint/constraint_facts.go`

**Fonction actuelle** :
```go
func ConvertFactsToReteFormat(program Program) ([]map[string]interface{}, error)
```

**Modifications nécessaires** :

```go
// ConvertFactsToReteFormat convertit les faits d'un programme en format RETE
// Gère les affectations de variables et les références entre faits
func ConvertFactsToReteFormat(program Program) ([]map[string]interface{}, error) {
    // Créer le contexte avec les types
    ctx := NewFactContext(program.Types)
    
    // Créer le type map
    typeMap := make(map[string]TypeDefinition)
    for _, t := range program.Types {
        typeMap[t.Name] = t
    }
    
    var reteFacts []map[string]interface{}
    
    // 1. Traiter d'abord les affectations de variables
    for i, assignment := range program.FactAssignments {
        typeDef, exists := typeMap[assignment.Fact.TypeName]
        if !exists {
            return nil, fmt.Errorf("affectation %d: type '%s' non défini", i+1, assignment.Fact.TypeName)
        }
        
        reteFact := createReteFact(assignment.Fact, typeDef)
        factID, err := ensureFactID(reteFact, assignment.Fact, typeDef, ctx)
        if err != nil {
            return nil, fmt.Errorf("affectation %d: %v", i+1, err)
        }
        
        reteFact[FieldNameInternalID] = factID
        reteFact[FieldNameReteType] = assignment.Fact.TypeName
        
        // Enregistrer la variable dans le contexte
        ctx.RegisterVariable(assignment.Variable, factID)
        
        reteFacts = append(reteFacts, reteFact)
    }
    
    // 2. Traiter les faits normaux (peuvent référencer les variables)
    for i, fact := range program.Facts {
        typeDef, exists := typeMap[fact.TypeName]
        if !exists {
            return nil, fmt.Errorf("fait %d: type '%s' non défini", i+1, fact.TypeName)
        }
        
        reteFact := createReteFact(fact, typeDef)
        factID, err := ensureFactID(reteFact, fact, typeDef, ctx)
        if err != nil {
            return nil, fmt.Errorf("fait %d: %v", i+1, err)
        }
        
        reteFact[FieldNameInternalID] = factID
        reteFact[FieldNameReteType] = fact.TypeName
        
        reteFacts = append(reteFacts, reteFact)
    }
    
    return reteFacts, nil
}

// ensureFactID génère l'ID d'un fait avec support des variables
func ensureFactID(reteFact map[string]interface{}, fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
    // Vérifier que _id_ n'a PAS été fourni manuellement
    if _, exists := reteFact[FieldNameInternalID]; exists {
        return "", fmt.Errorf(
            "le champ '%s' ne peut pas être défini manuellement pour le type '%s'",
            FieldNameInternalID,
            fact.TypeName,
        )
    }
    
    // Générer l'ID avec le contexte
    id, err := GenerateFactID(fact, typeDef, ctx)
    if err != nil {
        return "", fmt.Errorf("génération d'ID pour le fait de type '%s': %v", fact.TypeName, err)
    }
    
    return id, nil
}
```

#### Tests de Conversion

```go
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
```

### 6. Gérer le Hash avec Faits

#### Modification de `generateIDFromHash`

```go
// generateIDFromHash génère un ID par hash de tous les champs
// Gère les références de faits
func generateIDFromHash(fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
    // Construire une représentation canonique du fait
    var parts []string
    
    // Trier les champs par nom pour garantir le déterminisme
    sortedFields := make([]FactField, len(fact.Fields))
    copy(sortedFields, fact.Fields)
    sort.Slice(sortedFields, func(i, j int) bool {
        return sortedFields[i].Name < sortedFields[j].Name
    })
    
    for _, field := range sortedFields {
        // Trouver le type du champ
        var fieldType Field
        for _, f := range typeDef.Fields {
            if f.Name == field.Name {
                fieldType = f
                break
            }
        }
        
        // Convertir la valeur en string
        strValue, err := convertFieldValueToString(field.Value, fieldType, ctx)
        if err != nil {
            return "", fmt.Errorf("champ '%s': %v", field.Name, err)
        }
        
        parts = append(parts, fmt.Sprintf("%s=%s", field.Name, strValue))
    }
    
    // Construire la chaîne à hasher
    dataToHash := strings.Join(parts, "&")
    
    // Calculer le hash SHA-256
    hash := sha256.Sum256([]byte(dataToHash))
    hashStr := hex.EncodeToString(hash[:8]) // 16 caractères hex
    
    return fmt.Sprintf("%s%s%s", typeDef.Name, IDSeparatorType, hashStr), nil
}
```

---

## ✅ Critères de Succès

### Compilation et Tests

```bash
# Code compile
go build ./constraint

# Tests passent
go test ./constraint -run TestGenerate -v
go test ./constraint -run TestConvert -v

# Couverture > 80%
go test ./constraint -cover
```

### Fonctionnalités

- [ ] `FactContext` créé et fonctionnel
- [ ] Variables enregistrées et résolues
- [ ] `convertFieldValueToString` gère tous les types
- [ ] `GenerateFactID` accepte le contexte
- [ ] IDs générés avec références de faits
- [ ] IDs composites avec faits fonctionnels
- [ ] Hash fonctionne avec références
- [ ] `ConvertFactsToReteFormat` gère affectations

### Validation

```bash
make format
make lint
make validate
```

---

## 📊 Tests Requis

### Tests Unitaires Minimaux

- [ ] `TestFactContext`
- [ ] `TestConvertFieldValueToString`
- [ ] `TestGenerateFactID_WithFactReference`
- [ ] `TestGenerateFactID_CompositeKeyWithFact`
- [ ] `TestGenerateFactIDFromHash_WithFacts`
- [ ] `TestConvertFactsToReteFormat_WithAssignments`

### Tests d'Intégration

```go
func TestCompleteFlow_FactReferences(t *testing.T) {
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string, password: string)
        
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        
        Login(alice, "alice@ex.com", "pw1")
        Login(bob, "bob@ex.com", "pw2")
    `
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }
    
    reteFacts, err := ConvertFactsToReteFormat(program)
    if err != nil {
        t.Fatalf("Convert error: %v", err)
    }
    
    // Vérifier que 4 faits sont créés (2 User + 2 Login)
    if len(reteFacts) != 4 {
        t.Errorf("Expected 4 facts, got %d", len(reteFacts))
    }
    
    t.Log("✅ Flow complet fonctionne")
}
```

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Créer `FactContext`
2. ✅ Implémenter `convertFieldValueToString`
3. ✅ Modifier `GenerateFactID` avec contexte
4. ✅ Modifier `generateIDFromPrimaryKey`
5. ✅ Modifier `generateIDFromHash`
6. ✅ Modifier `ConvertFactsToReteFormat`
7. ✅ Tests complets
8. ✅ Validation

### Commandes

```bash
go test ./constraint -run TestFactContext -v
go test ./constraint -run TestConvert -v
go test ./constraint -run TestGenerate -v
make validate
```

---

## 📚 Références

- `scripts/new_ids/02-prompt-parser-syntax.md` - Syntaxe parser
- `constraint/id_generator.go` - Code actuel
- `docs/ID_RULES_COMPLETE.md` - Règles IDs

---

## 📝 Notes

### Points d'Attention

1. **Ordre de traitement** : Les affectations DOIVENT être traitées avant les faits qui les référencent

2. **Résolution de variables** : Le contexte doit être passé partout

3. **Déterminisme** : Le hash doit donner le même résultat pour les mêmes valeurs

---

## 🎯 Résultat Attendu

```go
// alice = User("Alice", 30) → ID: "User~Alice"
// Login(alice, "a@ex.com", "pw") → ID: "Login~a@ex.com"
// Le champ user de Login contient "User~Alice"
```

---

**Prompt suivant** : `04-prompt-evaluation.md`

**Durée estimée** : 4-6 heures

**Complexité** : 🔴 Élevée