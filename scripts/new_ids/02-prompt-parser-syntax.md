# Prompt 02 - Modification du Parser et Syntaxe

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Modifier la grammaire PEG et le parser pour supporter la nouvelle syntaxe TSD :

1. **Types comme valeurs de champs** : `type Login(user: User, ...)`
2. **Affectation de faits** : `a = User("Alice", 30)`
3. **Interdiction de `_id_`** dans les expressions
4. **Comparaisons simplifiées** : préparation pour `p.user == u`

---

## 📋 Contexte

### État Actuel

```tsd
// Syntaxe actuelle
type User(#name: string, age: number)
type Login(#email: string, password: string)

// Pas d'affectation possible
User("Alice", 30)
Login("alice@example.com", "pass123")

// Comparaisons explicites
{u: User, l: Login} / l.userEmail == u.email ==> ...
```

### État Cible

```tsd
// Nouvelle syntaxe
type User(#name: string, age: number)
type Login(user: User, #email: string, password: string)

// Affectation de faits
alice = User("Alice", 30)
bob = User("Bob", 25)

// Utilisation de variables de faits
Login(alice, "alice@example.com", "pass123")
Login(bob, "bob@example.com", "secret")

// Comparaisons simplifiées
{u: User, l: Login} / l.user == u ==> 
    Log("Login for " + u.name)
```

---

## 📝 Tâches à Réaliser

### 1. Analyser la Grammaire Actuelle

#### Fichier : `constraint/grammar/constraint.peg`

**Lire et comprendre** :

1. **Règles de types** :
   ```peg
   TypeDefinition <- ...
   Field <- ...
   FieldType <- ...
   ```

2. **Règles de faits** :
   ```peg
   Fact <- ...
   FactField <- ...
   FactValue <- ...
   ```

3. **Règles d'expressions** :
   ```peg
   Expression <- ...
   Constraint <- ...
   FieldAccess <- ...
   ```

4. **Identifiants et types** :
   ```peg
   Identifier <- ...
   TypeReference <- ...
   ```

**Questions à répondre** :
- Où sont définis les types de champs autorisés ?
- Comment sont parsées les valeurs de faits ?
- Y a-t-il déjà un mécanisme d'affectation ?
- Comment sont parsés les field access ?

### 2. Modifier la Grammaire pour Types de Faits

#### Objectif

Permettre : `type Login(user: User, #email: string, ...)`

#### Modifications de la Grammaire

**Avant** :
```peg
FieldType <- ("string" / "number" / "bool" / "boolean")
```

**Après** :
```peg
FieldType <- PrimitiveType / UserDefinedType

PrimitiveType <- ("string" / "number" / "bool" / "boolean")

UserDefinedType <- !ReservedWord Identifier
```

**Détails** :

1. **Ajouter `UserDefinedType`** :
   ```peg
   UserDefinedType <- !ReservedWord Identifier {
       // Retourner le nom du type
       return string(c.text), nil
   }
   ```

2. **Mettre à jour `Field`** :
   ```peg
   Field <- PrimaryKeyMarker? _ name:Identifier _ ":" _ type:FieldType {
       return map[string]interface{}{
           "name": name,
           "type": type,
           "isPrimaryKey": // ... selon PrimaryKeyMarker
       }, nil
   }
   ```

3. **Ajouter validation** :
   - Le type doit exister (validation post-parsing)
   - Pas de récursion circulaire (User → Login → User)

#### Tests de Parsing

```go
func TestParseTypeWithUserDefinedField(t *testing.T) {
    input := `type Login(user: User, #email: string, password: string)`
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Erreur de parsing: %v", err)
    }
    
    if len(program.Types) != 1 {
        t.Fatalf("Attendu 1 type, reçu %d", len(program.Types))
    }
    
    typeDef := program.Types[0]
    if typeDef.Name != "Login" {
        t.Errorf("Nom attendu 'Login', reçu '%s'", typeDef.Name)
    }
    
    if len(typeDef.Fields) != 3 {
        t.Fatalf("Attendu 3 champs, reçu %d", len(typeDef.Fields))
    }
    
    // Vérifier le champ user de type User
    userField := typeDef.Fields[0]
    if userField.Name != "user" {
        t.Errorf("Champ 0: attendu 'user', reçu '%s'", userField.Name)
    }
    if userField.Type != "User" {
        t.Errorf("Champ 0: type attendu 'User', reçu '%s'", userField.Type)
    }
    if userField.IsPrimaryKey {
        t.Error("Champ user ne devrait pas être clé primaire")
    }
    
    // Vérifier email (clé primaire)
    emailField := typeDef.Fields[1]
    if emailField.Name != "email" {
        t.Errorf("Champ 1: attendu 'email', reçu '%s'", emailField.Name)
    }
    if emailField.Type != "string" {
        t.Errorf("Champ 1: type attendu 'string', reçu '%s'", emailField.Type)
    }
    if !emailField.IsPrimaryKey {
        t.Error("Champ email devrait être clé primaire")
    }
}
```

### 3. Modifier la Grammaire pour Affectation de Faits

#### Objectif

Permettre : `alice = User("Alice", 30)`

#### Modifications de la Grammaire

**Ajouter une nouvelle règle `Statement`** :

```peg
Statement <- FactAssignment / Fact / Expression / TypeDefinition / ActionDefinition / ...

FactAssignment <- variable:Identifier _ "=" _ fact:Fact {
    return map[string]interface{}{
        "type": "factAssignment",
        "variable": variable,
        "fact": fact,
    }, nil
}

Fact <- typeName:Identifier "(" _ fields:FactFieldList? _ ")" {
    return map[string]interface{}{
        "type": "fact",
        "typeName": typeName,
        "fields": fields,
    }, nil
}
```

**Détails** :

1. **Structure `FactAssignment`** :
   ```go
   type FactAssignment struct {
       Type     string `json:"type"`     // "factAssignment"
       Variable string `json:"variable"` // Nom de la variable
       Fact     Fact   `json:"fact"`     // Le fait assigné
   }
   ```

2. **Ajouter à `Program`** :
   ```go
   type Program struct {
       Types           []TypeDefinition        `json:"types"`
       Actions         []ActionDefinition      `json:"actions"`
       XupleSpaces     []XupleSpaceDeclaration `json:"xupleSpaces"`
       Expressions     []Expression            `json:"expressions"`
       Facts           []Fact                  `json:"facts"`
       FactAssignments []FactAssignment        `json:"factAssignments"` // NOUVEAU
       Resets          []Reset                 `json:"resets"`
       RuleRemovals    []RuleRemoval           `json:"ruleRemovals"`
   }
   ```

3. **Parser doit distinguer** :
   - `User("Alice", 30)` → Fait simple
   - `alice = User("Alice", 30)` → Affectation

#### Tests de Parsing

```go
func TestParseFactAssignment(t *testing.T) {
    input := `alice = User("Alice", 30)`
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Erreur de parsing: %v", err)
    }
    
    if len(program.FactAssignments) != 1 {
        t.Fatalf("Attendu 1 affectation, reçu %d", len(program.FactAssignments))
    }
    
    assignment := program.FactAssignments[0]
    if assignment.Variable != "alice" {
        t.Errorf("Variable attendue 'alice', reçu '%s'", assignment.Variable)
    }
    
    fact := assignment.Fact
    if fact.TypeName != "User" {
        t.Errorf("Type attendu 'User', reçu '%s'", fact.TypeName)
    }
    
    if len(fact.Fields) != 2 {
        t.Fatalf("Attendu 2 champs, reçu %d", len(fact.Fields))
    }
}

func TestParseMultipleFactAssignments(t *testing.T) {
    input := `
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        Login(alice, "alice@example.com", "pass")
    `
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Erreur de parsing: %v", err)
    }
    
    if len(program.FactAssignments) != 2 {
        t.Errorf("Attendu 2 affectations, reçu %d", len(program.FactAssignments))
    }
    
    if len(program.Facts) != 1 {
        t.Errorf("Attendu 1 fait direct, reçu %d", len(program.Facts))
    }
}
```

### 4. Modifier la Grammaire pour Valeurs de Faits

#### Objectif

Permettre : `Login(alice, ...)` où `alice` est une variable

#### Modifications de la Grammaire

**Avant** :
```peg
FactValue <- StringLiteral / NumberLiteral / BooleanLiteral
```

**Après** :
```peg
FactValue <- StringLiteral / NumberLiteral / BooleanLiteral / VariableReference

VariableReference <- !ReservedWord Identifier {
    return map[string]interface{}{
        "type": "variableReference",
        "name": string(c.text),
    }, nil
}
```

**Structure** :

```go
type VariableReference struct {
    Type string `json:"type"` // "variableReference"
    Name string `json:"name"` // Nom de la variable
}
```

**Mise à jour de `FactValue`** :

```go
type FactValue struct {
    Type  string      `json:"type"`  // "string", "number", "bool", "variableReference"
    Value interface{} `json:"value"` // Valeur ou nom de variable
}
```

#### Tests

```go
func TestParseFactWithVariableReference(t *testing.T) {
    input := `Login(alice, "alice@example.com", "pass")`
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Erreur de parsing: %v", err)
    }
    
    fact := program.Facts[0]
    if len(fact.Fields) != 3 {
        t.Fatalf("Attendu 3 champs, reçu %d", len(fact.Fields))
    }
    
    // Premier champ est une référence de variable
    userField := fact.Fields[0]
    if userField.Value.Type != "variableReference" {
        t.Errorf("Type attendu 'variableReference', reçu '%s'", userField.Value.Type)
    }
    
    varName, ok := userField.Value.Value.(string)
    if !ok || varName != "alice" {
        t.Errorf("Nom de variable attendu 'alice', reçu '%v'", userField.Value.Value)
    }
}
```

### 5. Interdire `_id_` dans la Grammaire

#### Objectif

Rejeter au niveau du parser toute utilisation de `_id_`

#### Modifications de la Grammaire

**Ajouter aux mots réservés** :

```peg
ReservedWord <- (
    "type" /
    "action" /
    "rule" /
    "when" /
    "then" /
    "_id_" /     // NOUVEAU : Interdire _id_
    // ... autres mots réservés
) !IdentifierChar
```

**Validation dans `Field`** :

```peg
Field <- PrimaryKeyMarker? _ name:Identifier _ ":" _ type:FieldType {
    // Vérifier que le nom n'est pas _id_
    if name == "_id_" {
        return nil, errors.New("le champ '_id_' est réservé et ne peut pas être utilisé")
    }
    
    return map[string]interface{}{
        "name": name,
        "type": type,
        "isPrimaryKey": // ...
    }, nil
}
```

**Validation dans `FactField`** :

```peg
FactField <- name:Identifier _ ":" _ value:FactValue {
    // Vérifier que le nom n'est pas _id_
    if name == "_id_" {
        return nil, errors.New("le champ '_id_' est réservé et ne peut pas être assigné")
    }
    
    return map[string]interface{}{
        "name": name,
        "value": value,
    }, nil
}
```

#### Tests

```go
func TestParseType_InternalIDForbidden(t *testing.T) {
    input := `type User(_id_: string, name: string)`
    
    _, err := ParseProgram(input)
    
    if err == nil {
        t.Fatal("Attendu une erreur pour champ _id_")
    }
    
    if !strings.Contains(err.Error(), "réservé") {
        t.Errorf("Message d'erreur attendu contenant 'réservé', reçu: %v", err)
    }
}

func TestParseFact_InternalIDForbidden(t *testing.T) {
    input := `User(_id_: "manual", name: "Alice")`
    
    _, err := ParseProgram(input)
    
    if err == nil {
        t.Fatal("Attendu une erreur pour affectation _id_")
    }
    
    if !strings.Contains(err.Error(), "réservé") {
        t.Errorf("Message d'erreur attendu contenant 'réservé', reçu: %v", err)
    }
}

func TestParseFieldAccess_InternalIDForbidden(t *testing.T) {
    input := `{u: User} / u._id_ == "test" ==> Log("test")`
    
    _, err := ParseProgram(input)
    
    // Note: Cette validation peut aussi être post-parsing
    // Selon l'approche choisie
    if err == nil {
        t.Log("⚠️ Validation de _id_ en field access peut être post-parsing")
    }
}
```

### 6. Préparer Comparaisons Simplifiées

#### Objectif

Parser `p.user == u` (résolution de type post-parsing)

#### Note Importante

Le parser doit **accepter** cette syntaxe, mais la **résolution** se fera dans le prochain prompt (évaluation).

**Pas de modification de grammaire nécessaire** - la grammaire actuelle devrait déjà accepter :
- `p.user` (field access)
- `u` (variable)
- `==` (comparateur)

**Vérification** :

```go
func TestParseFactComparison(t *testing.T) {
    input := `{u: User, l: Login} / l.user == u ==> Log("test")`
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Erreur de parsing: %v", err)
    }
    
    if len(program.Expressions) != 1 {
        t.Fatalf("Attendu 1 expression, reçu %d", len(program.Expressions))
    }
    
    expr := program.Expressions[0]
    
    // Vérifier que la contrainte est parsée
    // l.user == u
    // Left: FieldAccess{Object: "l", Field: "user"}
    // Right: Variable{Name: "u"}
    
    t.Log("✅ Syntaxe p.user == u acceptée par le parser")
}
```

### 7. Régénérer le Parser

#### Commandes

```bash
# Se positionner dans le répertoire constraint
cd constraint

# Vérifier que pigeon est installé
which pigeon || go install github.com/mna/pigeon@latest

# Régénérer le parser
pigeon -o parser.go grammar/constraint.peg

# Vérifier la génération
ls -lh parser.go

# Compiler pour vérifier
go build ./...
```

#### Vérifications Post-Génération

```bash
# Le parser doit compiler
go build ./constraint

# Tests de parsing doivent passer
go test ./constraint -run TestParse -v

# Vérifier taille du parser (ne doit pas exploser)
wc -l constraint/parser.go
```

### 8. Mettre à Jour l'API de Parsing

#### Fichier : `constraint/api.go`

**Ajouter fonctions pour nouveaux types** :

```go
// ParseFactAssignment parses a fact assignment statement
func ParseFactAssignment(input string) (*FactAssignment, error) {
    program, err := ParseProgram(input)
    if err != nil {
        return nil, err
    }
    
    if len(program.FactAssignments) == 0 {
        return nil, errors.New("aucune affectation de fait trouvée")
    }
    
    return &program.FactAssignments[0], nil
}

// ValidateVariableReference validates that a variable reference is defined
func ValidateVariableReference(varName string, assignments []FactAssignment) error {
    for _, assignment := range assignments {
        if assignment.Variable == varName {
            return nil
        }
    }
    return fmt.Errorf("variable '%s' non définie", varName)
}
```

### 9. Mise à Jour de la Validation

#### Fichier : `constraint/constraint_validation.go` (ou nouveau fichier)

**Créer validations post-parsing** :

```go
// ValidateProgram validates a complete program after parsing
func ValidateProgram(program Program) error {
    // 1. Valider que les types référencés dans les champs existent
    if err := validateTypeReferences(program); err != nil {
        return err
    }
    
    // 2. Valider que les variables référencées sont définies
    if err := validateVariableReferences(program); err != nil {
        return err
    }
    
    // 3. Valider qu'il n'y a pas de récursion circulaire
    if err := validateNoCircularReferences(program); err != nil {
        return err
    }
    
    return nil
}

func validateTypeReferences(program Program) error {
    typeMap := make(map[string]bool)
    for _, typeDef := range program.Types {
        typeMap[typeDef.Name] = true
    }
    
    primitiveTypes := map[string]bool{
        "string": true,
        "number": true,
        "bool": true,
        "boolean": true,
    }
    
    for _, typeDef := range program.Types {
        for _, field := range typeDef.Fields {
            // Si ce n'est pas un type primitif, vérifier qu'il existe
            if !primitiveTypes[field.Type] && !typeMap[field.Type] {
                return fmt.Errorf(
                    "type '%s': champ '%s' référence un type inconnu '%s'",
                    typeDef.Name,
                    field.Name,
                    field.Type,
                )
            }
        }
    }
    
    return nil
}

func validateVariableReferences(program Program) error {
    // Créer un map des variables définies par affectation
    varMap := make(map[string]string) // varName -> typeName
    
    for _, assignment := range program.FactAssignments {
        varMap[assignment.Variable] = assignment.Fact.TypeName
    }
    
    // Vérifier les faits qui utilisent des variables
    for i, fact := range program.Facts {
        for j, field := range fact.Fields {
            if field.Value.Type == "variableReference" {
                varName, ok := field.Value.Value.(string)
                if !ok {
                    return fmt.Errorf("fait %d, champ %d: référence de variable invalide", i+1, j+1)
                }
                
                if _, exists := varMap[varName]; !exists {
                    return fmt.Errorf(
                        "fait %d, champ %d: variable '%s' non définie",
                        i+1,
                        j+1,
                        varName,
                    )
                }
            }
        }
    }
    
    return nil
}

func validateNoCircularReferences(program Program) error {
    typeGraph := make(map[string][]string) // type -> types qu'il référence
    
    for _, typeDef := range program.Types {
        for _, field := range typeDef.Fields {
            // Si le champ est un type utilisateur
            if isUserDefinedType(field.Type, program.Types) {
                typeGraph[typeDef.Name] = append(typeGraph[typeDef.Name], field.Type)
            }
        }
    }
    
    // Détection de cycle (DFS)
    for typeName := range typeGraph {
        visited := make(map[string]bool)
        if hasCycle(typeName, typeGraph, visited, make(map[string]bool)) {
            return fmt.Errorf("référence circulaire détectée impliquant le type '%s'", typeName)
        }
    }
    
    return nil
}

func hasCycle(node string, graph map[string][]string, visited, recStack map[string]bool) bool {
    visited[node] = true
    recStack[node] = true
    
    for _, neighbor := range graph[node] {
        if !visited[neighbor] {
            if hasCycle(neighbor, graph, visited, recStack) {
                return true
            }
        } else if recStack[neighbor] {
            return true
        }
    }
    
    recStack[node] = false
    return false
}

func isUserDefinedType(typeName string, types []TypeDefinition) bool {
    primitives := map[string]bool{
        "string": true,
        "number": true,
        "bool": true,
        "boolean": true,
    }
    
    if primitives[typeName] {
        return false
    }
    
    for _, t := range types {
        if t.Name == typeName {
            return true
        }
    }
    
    return false
}
```

#### Tests de Validation

```go
func TestValidateTypeReferences(t *testing.T) {
    tests := []struct {
        name    string
        program Program
        wantErr bool
    }{
        {
            name: "type valide référencé",
            program: Program{
                Types: []TypeDefinition{
                    {
                        Name: "User",
                        Fields: []Field{
                            {Name: "name", Type: "string", IsPrimaryKey: true},
                        },
                    },
                    {
                        Name: "Login",
                        Fields: []Field{
                            {Name: "user", Type: "User"},
                            {Name: "password", Type: "string"},
                        },
                    },
                },
            },
            wantErr: false,
        },
        {
            name: "type inconnu référencé",
            program: Program{
                Types: []TypeDefinition{
                    {
                        Name: "Login",
                        Fields: []Field{
                            {Name: "user", Type: "UnknownType"},
                        },
                    },
                },
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateTypeReferences(tt.program)
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr %v, got err %v", tt.wantErr, err)
            }
        })
    }
}

func TestValidateCircularReferences(t *testing.T) {
    program := Program{
        Types: []TypeDefinition{
            {
                Name: "A",
                Fields: []Field{
                    {Name: "b", Type: "B"},
                },
            },
            {
                Name: "B",
                Fields: []Field{
                    {Name: "a", Type: "A"},
                },
            },
        },
    }
    
    err := validateNoCircularReferences(program)
    if err == nil {
        t.Error("Attendu une erreur pour référence circulaire")
    }
}
```

---

## ✅ Critères de Succès

### Compilation et Génération

```bash
# Parser se régénère sans erreur
pigeon -o parser.go grammar/constraint.peg

# Code compile
go build ./constraint

# Tests de parsing passent
go test ./constraint -run TestParse -v
```

### Nouvelles Fonctionnalités

- [ ] Parser accepte `type Login(user: User, ...)`
- [ ] Parser accepte `alice = User("Alice", 30)`
- [ ] Parser accepte `Login(alice, ...)`
- [ ] Parser rejette `_id_` comme nom de champ
- [ ] Parser accepte `l.user == u`
- [ ] Validation détecte types inconnus
- [ ] Validation détecte variables non définies
- [ ] Validation détecte références circulaires

### Tests et Couverture

```bash
# Tous les tests passent
go test ./constraint -v

# Couverture > 80%
go test ./constraint -cover

# Validation complète
make validate
```

---

## 📊 Tests Requis

### Tests de Parsing Minimaux

- [ ] `TestParseTypeWithUserDefinedField`
- [ ] `TestParseFactAssignment`
- [ ] `TestParseMultipleFactAssignments`
- [ ] `TestParseFactWithVariableReference`
- [ ] `TestParseType_InternalIDForbidden`
- [ ] `TestParseFact_InternalIDForbidden`
- [ ] `TestParseFactComparison`

### Tests de Validation Minimaux

- [ ] `TestValidateTypeReferences`
- [ ] `TestValidateVariableReferences`
- [ ] `TestValidateCircularReferences`

### Tests d'Intégration

```go
func TestParseAndValidate_Complete(t *testing.T) {
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string, password: string)
        
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        
        Login(alice, "alice@example.com", "pass123")
        Login(bob, "bob@example.com", "secret")
    `
    
    program, err := ParseProgram(input)
    if err != nil {
        t.Fatalf("Erreur de parsing: %v", err)
    }
    
    err = ValidateProgram(program)
    if err != nil {
        t.Fatalf("Erreur de validation: %v", err)
    }
    
    // Vérifications
    if len(program.Types) != 2 {
        t.Errorf("Attendu 2 types, reçu %d", len(program.Types))
    }
    
    if len(program.FactAssignments) != 2 {
        t.Errorf("Attendu 2 affectations, reçu %d", len(program.FactAssignments))
    }
    
    if len(program.Facts) != 2 {
        t.Errorf("Attendu 2 faits, reçu %d", len(program.Facts))
    }
    
    t.Log("✅ Programme complet parsé et validé")
}
```

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Analyser grammaire actuelle
2. ✅ Modifier grammaire PEG (types utilisateur)
3. ✅ Modifier grammaire PEG (affectations)
4. ✅ Modifier grammaire PEG (variables)
5. ✅ Modifier grammaire PEG (interdire _id_)
6. ✅ Régénérer parser
7. ✅ Mettre à jour structures
8. ✅ Implémenter validations
9. ✅ Tests complets
10. ✅ Validation finale

### Commandes

```bash
# 1. Modifier la grammaire
vim constraint/grammar/constraint.peg

# 2. Régénérer le parser
cd constraint
pigeon -o parser.go grammar/constraint.peg

# 3. Vérifier compilation
go build ./...

# 4. Tester
go test ./constraint -v

# 5. Validation
cd ..
make validate
```

---

## 📚 Références

- `.github/prompts/common.md` - Standards
- `.github/prompts/develop.md` - Développement
- `scripts/new_ids/01-prompt-structures-base.md` - Structures de base
- `constraint/grammar/constraint.peg` - Grammaire actuelle
- [PEG Parser Documentation](https://github.com/mna/pigeon)

---

## 📝 Notes

### Points d'Attention

1. **Parser généré** : `parser.go` est généré depuis `constraint.peg`. **Modifier uniquement la grammaire**, puis régénérer.

2. **Rétrocompatibilité** : Cette modification CASSE la compatibilité. Les anciens programmes doivent être migrés.

3. **Validation en deux temps** :
   - Parser : validation syntaxique
   - Post-parsing : validation sémantique (types, variables)

4. **Performance** : La grammaire PEG est récursive, attention aux grammaires ambiguës.

### Questions Résolues

Q: Valider `_id_` dans le parser ou post-parsing ?
R: Les deux (défense en profondeur). Parser rejette syntaxiquement, validation vérifie aussi.

Q: Comment gérer les types récursifs (ex: Tree) ?
R: Interdire pour l'instant (validation détecte cycles). À supporter plus tard si besoin.

---

## 🎯 Résultat Attendu

Après ce prompt :

```tsd
// ✅ Syntaxe acceptée
type User(#name: string, age: number)
type Login(user: User, #email: string, password: string)

alice = User("Alice", 30)
Login(alice, "alice@example.com", "pass")

{u: User, l: Login} / l.user == u ==> Log("OK")

// ❌ Syntaxe rejetée
type Bad(_id_: string)           // Erreur: _id_ réservé
User(_id_: "manual")             // Erreur: _id_ réservé
Login(unknown_var, "test")       // Erreur: variable non définie
type Loop(field: Loop)           // Erreur: référence circulaire
```

---

**Prompt suivant** : `03-prompt-id-generation.md`

**Durée estimée** : 6-8 heures

**Complexité** : 🔴 Élevée (modification grammaire + régénération)