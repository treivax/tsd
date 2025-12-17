# ID Generator API Reference

This document provides API reference for developers working with TSD's ID generation system.

---

## Table of Contents

- [Package: constraint](#package-constraint)
  - [Types](#types)
  - [Functions](#functions)
  - [Validation](#validation)
- [Package: rete](#package-rete)
  - [Fact Structure](#fact-structure)
  - [Working Memory](#working-memory)
- [Examples](#examples)
- [Error Handling](#error-handling)

---

## Package: constraint

### Types

#### TypeDefinition

Represents a type definition with metadata about primary keys.

```go
type TypeDefinition struct {
    Name   string  `json:"name"`
    Fields []Field `json:"fields"`
}
```

**Methods**:

```go
// Returns true if at least one field is marked as primary key
func (td *TypeDefinition) HasPrimaryKey() bool

// Returns all fields marked as primary keys
func (td *TypeDefinition) GetPrimaryKeyFields() []Field

// Returns names of all primary key fields
func (td *TypeDefinition) GetPrimaryKeyFieldNames() []string
```

**Example**:

```go
typeDef := &TypeDefinition{
    Name: "User",
    Fields: []Field{
        {Name: "username", Type: "string", IsPrimaryKey: true},
        {Name: "email", Type: "string", IsPrimaryKey: false},
        {Name: "age", Type: "number", IsPrimaryKey: false},
    },
}

hasPK := typeDef.HasPrimaryKey()              // true
pkFields := typeDef.GetPrimaryKeyFields()     // [{username, string, true}]
pkNames := typeDef.GetPrimaryKeyFieldNames()  // ["username"]
```

---

#### Field

Represents a field in a type definition.

```go
type Field struct {
    Name         string `json:"name"`
    Type         string `json:"type"`
    IsPrimaryKey bool   `json:"isPrimaryKey,omitempty"`
}
```

**Valid Types**: `"string"`, `"number"`, `"bool"`, `"object"`, `"list"`

**Note**: Only `string`, `number`, and `bool` can be marked as primary keys.

---

### Functions

#### GenerateFactID

Generates a unique, deterministic ID for a fact.

```go
func GenerateFactID(typeName string, fact map[string]interface{}, typeDef *TypeDefinition) string
```

**Parameters**:
- `typeName` (string): Name of the type (e.g., "User", "Product")
- `fact` (map[string]interface{}): Fact data with field values
- `typeDef` (*TypeDefinition): Type definition with primary key metadata

**Returns**: `string` - Generated ID in format `TypeName~identifier`

**Behavior**:
- **With primary key**: Returns `TypeName~value` or `TypeName~val1_val2_val3`
- **Without primary key**: Returns `TypeName~<16-char-hex-hash>`

**Example**:

```go
// Simple primary key
typeDef := &TypeDefinition{
    Name: "User",
    Fields: []Field{
        {Name: "username", Type: "string", IsPrimaryKey: true},
    },
}
fact := map[string]interface{}{"username": "alice", "email": "alice@example.com"}
id := GenerateFactID("User", fact, typeDef)
// Result: "User~alice"

// Composite primary key
typeDef2 := &TypeDefinition{
    Name: "Product",
    Fields: []Field{
        {Name: "category", Type: "string", IsPrimaryKey: true},
        {Name: "name", Type: "string", IsPrimaryKey: true},
    },
}
fact2 := map[string]interface{}{
    "category": "Electronics",
    "name": "Laptop",
    "price": 1200,
}
id2 := GenerateFactID("Product", fact2, typeDef2)
// Result: "Product~Electronics_Laptop"

// No primary key (hash-based)
typeDef3 := &TypeDefinition{
    Name: "LogEvent",
    Fields: []Field{
        {Name: "timestamp", Type: "number"},
        {Name: "message", Type: "string"},
    },
}
fact3 := map[string]interface{}{
    "timestamp": 1704067200,
    "message": "App started",
}
id3 := GenerateFactID("LogEvent", fact3, typeDef3)
// Result: "LogEvent~a1b2c3d4e5f6g7h8" (deterministic hash)
```

**Special Character Handling**:

Primary key values are percent-encoded:
- `~` → `%7E`
- `_` → `%5F`
- `%` → `%25`
- ` ` (space) → `%20`
- `/` → `%2F`

```go
fact := map[string]interface{}{"path": "/home/user~backup"}
id := GenerateFactID("File", fact, fileTypeDef)
// Result: "File~%2Fhome%2Fuser%7Ebackup"
```

---

#### ParseFactID

Parses an ID string to extract type name and identifier.

```go
func ParseFactID(id string) (typeName string, identifier string, err error)
```

**Parameters**:
- `id` (string): ID to parse (e.g., "User~alice")

**Returns**:
- `typeName` (string): Extracted type name
- `identifier` (string): Extracted identifier part
- `err` (error): Error if format is invalid

**Example**:

```go
typeName, identifier, err := ParseFactID("User~alice")
// typeName = "User"
// identifier = "alice"
// err = nil

typeName, identifier, err := ParseFactID("Product~Electronics_Laptop")
// typeName = "Product"
// identifier = "Electronics_Laptop"
// err = nil

typeName, identifier, err := ParseFactID("invalid_format")
// err != nil (no ~ separator found)
```

---

### Validation

#### ValidatePrimaryKeyTypes

Validates that all primary key fields are of primitive types.

```go
func ValidatePrimaryKeyTypes(typeDef *TypeDefinition) error
```

**Parameters**:
- `typeDef` (*TypeDefinition): Type definition to validate

**Returns**: `error` if any primary key field is not a primitive type

**Valid Primary Key Types**: `string`, `number`, `bool`

**Example**:

```go
// Valid
typeDef := &TypeDefinition{
    Name: "User",
    Fields: []Field{
        {Name: "username", Type: "string", IsPrimaryKey: true},
        {Name: "age", Type: "number", IsPrimaryKey: false},
    },
}
err := ValidatePrimaryKeyTypes(typeDef)
// err == nil

// Invalid
typeDef2 := &TypeDefinition{
    Name: "Document",
    Fields: []Field{
        {Name: "metadata", Type: "object", IsPrimaryKey: true},
    },
}
err2 := ValidatePrimaryKeyTypes(typeDef2)
// err2 != nil: "primary key field 'metadata' must be a primitive type..."
```

---

#### ValidatePrimaryKeyFieldsPresent

Validates that all primary key fields are present in a fact.

```go
func ValidatePrimaryKeyFieldsPresent(fact map[string]interface{}, typeDef *TypeDefinition) error
```

**Parameters**:
- `fact` (map[string]interface{}): Fact data to validate
- `typeDef` (*TypeDefinition): Type definition with primary key metadata

**Returns**: `error` if any primary key field is missing

**Example**:

```go
typeDef := &TypeDefinition{
    Name: "User",
    Fields: []Field{
        {Name: "username", Type: "string", IsPrimaryKey: true},
        {Name: "email", Type: "string", IsPrimaryKey: false},
    },
}

// Valid
fact1 := map[string]interface{}{"username": "alice", "email": "alice@example.com"}
err := ValidatePrimaryKeyFieldsPresent(fact1, typeDef)
// err == nil

// Invalid - missing primary key field
fact2 := map[string]interface{}{"email": "bob@example.com"}
err2 := ValidatePrimaryKeyFieldsPresent(fact2, typeDef)
// err2 != nil: "primary key field 'username' not found in fact of type 'User'"
```

---

#### ValidateNoExplicitID

Validates that the reserved `id` field is not set explicitly.

```go
func ValidateNoExplicitID(fact map[string]interface{}) error
```

**Parameters**:
- `fact` (map[string]interface{}): Fact data to validate

**Returns**: `error` if `id` field is present

**Example**:

```go
// Valid
fact1 := map[string]interface{}{"username": "alice", "email": "alice@example.com"}
err := ValidateNoExplicitID(fact1)
// err == nil

// Invalid - explicit id field
fact2 := map[string]interface{}{"id": "custom_id", "username": "bob"}
err2 := ValidateNoExplicitID(fact2)
// err2 != nil: "field 'id' is reserved and cannot be set explicitly"
```

---

## Package: rete

### Fact Structure

The `Fact` struct represents a fact in the RETE network.

```go
type Fact struct {
    ID     string                 `json:"id"`
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}
```

**Fields**:
- `ID` (string): Generated unique identifier (read-only)
- `Type` (string): Type name (e.g., "User", "Product")
- `Fields` (map[string]interface{}): Field values

**Example**:

```go
fact := Fact{
    ID:   "User~alice",
    Type: "User",
    Fields: map[string]interface{}{
        "username": "alice",
        "email":    "alice@example.com",
        "age":      30,
    },
}
```

**Note**: The `ID` field is set automatically during fact creation and should not be modified.

---

### Working Memory

#### AddFact

Adds a fact to working memory.

```go
func (wm *WorkingMemory) AddFact(fact Fact) error
```

**Parameters**:
- `fact` (Fact): Fact to add

**Returns**: `error` if fact already exists with the same ID

**Example**:

```go
fact := Fact{
    ID:   "User~alice",
    Type: "User",
    Fields: map[string]interface{}{"username": "alice"},
}

err := workingMemory.AddFact(fact)
if err != nil {
    log.Fatalf("Failed to add fact: %v", err)
}
```

---

#### GetFact

Retrieves a fact by ID.

```go
func (wm *WorkingMemory) GetFact(id string) (Fact, bool)
```

**Parameters**:
- `id` (string): Fact ID to retrieve

**Returns**:
- `Fact`: The fact (if found)
- `bool`: True if found, false otherwise

**Example**:

```go
fact, found := workingMemory.GetFact("User~alice")
if found {
    fmt.Printf("Found user: %v\n", fact.Fields["username"])
} else {
    fmt.Println("User not found")
}
```

---

#### RemoveFact

Removes a fact from working memory.

```go
func (wm *WorkingMemory) RemoveFact(id string) error
```

**Parameters**:
- `id` (string): ID of fact to remove

**Returns**: `error` if fact not found

**Example**:

```go
err := workingMemory.RemoveFact("User~alice")
if err != nil {
    log.Printf("Failed to remove fact: %v", err)
}
```

---

## Examples

### Example 1: Creating Facts with IDs

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/constraint"
)

func main() {
    // Define type with primary key
    typeDef := &constraint.TypeDefinition{
        Name: "User",
        Fields: []constraint.Field{
            {Name: "username", Type: "string", IsPrimaryKey: true},
            {Name: "email", Type: "string"},
        },
    }

    // Create fact
    fact := map[string]interface{}{
        "username": "alice",
        "email":    "alice@example.com",
    }

    // Validate
    if err := constraint.ValidatePrimaryKeyTypes(typeDef); err != nil {
        panic(err)
    }
    if err := constraint.ValidatePrimaryKeyFieldsPresent(fact, typeDef); err != nil {
        panic(err)
    }
    if err := constraint.ValidateNoExplicitID(fact); err != nil {
        panic(err)
    }

    // Generate ID
    id := constraint.GenerateFactID("User", fact, typeDef)
    fmt.Println("Generated ID:", id)
    // Output: Generated ID: User~alice
}
```

---

### Example 2: Working with Composite Keys

```go
// Define type with composite primary key
typeDef := &constraint.TypeDefinition{
    Name: "Enrollment",
    Fields: []constraint.Field{
        {Name: "student_id", Type: "string", IsPrimaryKey: true},
        {Name: "course_id", Type: "string", IsPrimaryKey: true},
        {Name: "grade", Type: "string"},
    },
}

// Create fact
fact := map[string]interface{}{
    "student_id": "S2024001",
    "course_id":  "CS101",
    "grade":      "A",
}

// Generate ID
id := constraint.GenerateFactID("Enrollment", fact, typeDef)
fmt.Println("Generated ID:", id)
// Output: Generated ID: Enrollment~S2024001_CS101
```

---

### Example 3: Hash-Based IDs

```go
// Define type without primary key
typeDef := &constraint.TypeDefinition{
    Name: "LogEvent",
    Fields: []constraint.Field{
        {Name: "timestamp", Type: "number"},
        {Name: "level", Type: "string"},
        {Name: "message", Type: "string"},
    },
}

// Create fact
fact := map[string]interface{}{
    "timestamp": 1704067200,
    "level":     "ERROR",
    "message":   "Connection failed",
}

// Generate ID (hash-based)
id := constraint.GenerateFactID("LogEvent", fact, typeDef)
fmt.Println("Generated ID:", id)
// Output: Generated ID: LogEvent~a1b2c3d4e5f6g7h8
```

---

## Error Handling

### Common Errors

#### Primary Key Type Error

```go
err := ValidatePrimaryKeyTypes(typeDef)
// Error: "primary key field 'metadata' must be a primitive type (string, number, bool), got 'object'"
```

**Solution**: Only use primitive types for primary keys.

---

#### Missing Primary Key Field

```go
err := ValidatePrimaryKeyFieldsPresent(fact, typeDef)
// Error: "primary key field 'username' not found in fact of type 'User'"
```

**Solution**: Ensure all primary key fields are present in the fact.

---

#### Explicit ID Field

```go
err := ValidateNoExplicitID(fact)
// Error: "field 'id' is reserved and cannot be set explicitly"
```

**Solution**: Remove the `id` field from your fact. IDs are generated automatically.

---

#### Invalid ID Format (ParseFactID)

```go
typeName, identifier, err := ParseFactID("invalid_format")
// Error: "invalid ID format: missing type/identifier separator '~'"
```

**Solution**: Ensure ID follows format `TypeName~identifier`.

---

## Best Practices

1. **Always validate** before generating IDs:
   ```go
   if err := ValidatePrimaryKeyTypes(typeDef); err != nil { ... }
   if err := ValidatePrimaryKeyFieldsPresent(fact, typeDef); err != nil { ... }
   if err := ValidateNoExplicitID(fact); err != nil { ... }
   ```

2. **Cache type definitions**: Don't parse type definitions repeatedly.

3. **Use meaningful primary key field names**: Prefer `username`, `sku`, `order_number` over generic `id`.

4. **Document your ID strategy**: Explain why each type uses its particular primary key.

5. **Test with special characters**: Ensure your primary key values work correctly when escaped.

---

## See Also

- [Primary Keys User Guide](../primary-keys.md) - User documentation
- [Architecture Documentation](../architecture/id-generation.md) - Internal design
- [Migration Guide](../MIGRATION_IDS.md) - Migrating to new syntax
- [Examples](../../examples/pk_*.tsd) - Code examples

---

**Version**: 1.0  
**Last Updated**: 2024-12-17  
**Maintainer**: TSD Core Team