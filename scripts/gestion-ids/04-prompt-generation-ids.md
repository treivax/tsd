# Prompt 04 - Implémentation de la Génération d'Identifiants

> **📋 Standards** : Ce prompt respecte les règles de [common.md](../../.github/prompts/common.md) et [develop.md](../../.github/prompts/develop.md)

## 🎯 Objectif

Implémenter la logique de génération automatique des identifiants de faits basée sur les clés primaires ou le hash des valeurs.

## 📋 Contexte

Suite aux modifications précédentes (grammaire, structures, validation), nous devons maintenant implémenter la génération effective des ID selon les règles :

1. **Avec clé primaire** : `TypeName~value1_value2_..._valueN`
2. **Sans clé primaire** : `TypeName~<hash>`

### Exemples

```tsd
type Person(#firstName: string, #lastName: string, age: number)
Person(firstName: "Jean-Claude", lastName: "Pignon", age: 27)
# ID généré: "Person~Jean-Claude_Pignon"

type User(#login: string, name: string)
User(login: "jcp", name: "Jean-Claude Pignon")
# ID généré: "User~jcp"

type Document(title: string, content: string)
Document(title: "Doc1", content: "...")
# ID généré: "Document~a3f5b9c2e1d4f8a7" (hash)
```

## 🔍 Analyse Préliminaire

### Fichiers à Créer

1. **`constraint/id_generator.go`** - Logique de génération d'ID
2. **`constraint/id_generator_test.go`** - Tests unitaires

### Fichiers à Modifier

1. **`constraint/constraint_facts.go`** - Intégration de la génération d'ID
2. **`rete/constraint_pipeline_facts.go`** - Utilisation des ID générés (si nécessaire)

### Algorithmes Nécessaires

1. **Génération d'ID avec clé primaire** :
   - Extraire les valeurs des champs PK dans l'ordre
   - Convertir chaque valeur en string
   - Échapper les caractères spéciaux `~` et `_`
   - Concaténer avec le format `TypeName~value1_value2_...`

2. **Génération d'ID par hash** :
   - Concaténer toutes les valeurs des champs
   - Calculer le hash MD5 ou SHA-256
   - Tronquer à une longueur raisonnable (16 caractères hex)
   - Format: `TypeName~<hash>`

## 🔧 Implémentation

### Étape 1 : Créer le Générateur d'ID

**Fichier** : `constraint/id_generator.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Constantes pour la génération d'ID
const (
	// IDSeparatorType sépare le nom du type et les valeurs de clé primaire
	IDSeparatorType = "~"
	
	// IDSeparatorValue sépare les valeurs de clé primaire entre elles
	IDSeparatorValue = "_"
	
	// IDHashLength est la longueur du hash pour les ID sans clé primaire
	IDHashLength = 16
)

// GenerateFactID génère l'identifiant d'un fait selon ses clés primaires.
// Si le type a une clé primaire définie, l'ID est : TypeName~value1_value2_...
// Sinon, l'ID est : TypeName~<hash>
func GenerateFactID(fact Fact, typeDef TypeDefinition) (string, error) {
	// Vérifier si le type a une clé primaire
	if typeDef.HasPrimaryKey() {
		return generateIDFromPrimaryKey(fact, typeDef)
	}
	
	// Sinon, générer un ID par hash
	return generateIDFromHash(fact, typeDef)
}

// generateIDFromPrimaryKey génère un ID basé sur les valeurs de clé primaire.
// Format: TypeName~value1_value2_..._valueN
func generateIDFromPrimaryKey(fact Fact, typeDef TypeDefinition) (string, error) {
	pkFields := typeDef.GetPrimaryKeyFields()
	if len(pkFields) == 0 {
		return "", fmt.Errorf("type '%s' n'a pas de clé primaire définie", typeDef.Name)
	}
	
	// Créer une map des valeurs du fait
	factValues := make(map[string]FactValue)
	for _, field := range fact.Fields {
		factValues[field.Name] = field.Value
	}
	
	// Extraire les valeurs PK dans l'ordre de définition
	var pkValues []string
	for _, pkField := range pkFields {
		factValue, exists := factValues[pkField.Name]
		if !exists {
			return "", fmt.Errorf("champ de clé primaire '%s' manquant dans le fait", pkField.Name)
		}
		
		// Convertir la valeur en string
		strValue, err := valueToString(factValue.Value)
		if err != nil {
			return "", fmt.Errorf("impossible de convertir le champ '%s' en string: %v", pkField.Name, err)
		}
		
		// Échapper les caractères spéciaux
		escapedValue := escapeIDValue(strValue)
		pkValues = append(pkValues, escapedValue)
	}
	
	// Construire l'ID
	id := typeDef.Name + IDSeparatorType + strings.Join(pkValues, IDSeparatorValue)
	return id, nil
}

// generateIDFromHash génère un ID basé sur le hash de toutes les valeurs.
// Format: TypeName~<hash>
func generateIDFromHash(fact Fact, typeDef TypeDefinition) (string, error) {
	// Concaténer toutes les valeurs dans un ordre déterministe
	// Utiliser l'ordre de définition des champs du type
	var valueStrings []string
	
	// Créer une map des valeurs du fait
	factValues := make(map[string]FactValue)
	for _, field := range fact.Fields {
		factValues[field.Name] = field.Value
	}
	
	// Parcourir les champs dans l'ordre de définition du type
	for _, field := range typeDef.Fields {
		factValue, exists := factValues[field.Name]
		if exists && factValue.Value != nil {
			strValue, err := valueToString(factValue.Value)
			if err != nil {
				return "", fmt.Errorf("impossible de convertir le champ '%s' en string: %v", field.Name, err)
			}
			valueStrings = append(valueStrings, field.Name+"="+strValue)
		}
	}
	
	// Calculer le hash MD5
	concatenated := strings.Join(valueStrings, "|")
	hash := md5.Sum([]byte(concatenated))
	hashStr := hex.EncodeToString(hash[:])
	
	// Tronquer à la longueur souhaitée
	if len(hashStr) > IDHashLength {
		hashStr = hashStr[:IDHashLength]
	}
	
	// Construire l'ID
	id := typeDef.Name + IDSeparatorType + hashStr
	return id, nil
}

// valueToString convertit une valeur de fait en string.
func valueToString(value interface{}) (string, error) {
	if value == nil {
		return "", fmt.Errorf("valeur nulle")
	}
	
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		// Pour les floats, utiliser une précision fixe pour cohérence
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	default:
		// Fallback sur fmt.Sprintf
		return fmt.Sprintf("%v", v), nil
	}
}

// escapeIDValue échappe les caractères spéciaux dans une valeur d'ID.
// Remplace ~ par %7E et _ par %5F (URL encoding partiel)
func escapeIDValue(value string) string {
	value = strings.ReplaceAll(value, "%", "%25") // % en premier pour éviter double-escape
	value = strings.ReplaceAll(value, IDSeparatorType, "%7E")
	value = strings.ReplaceAll(value, IDSeparatorValue, "%5F")
	return value
}

// unescapeIDValue inverse l'échappement des caractères spéciaux.
func unescapeIDValue(value string) string {
	value = strings.ReplaceAll(value, "%7E", IDSeparatorType)
	value = strings.ReplaceAll(value, "%5F", IDSeparatorValue)
	value = strings.ReplaceAll(value, "%25", "%")
	return value
}

// ParseFactID décompose un ID de fait en type et valeurs.
// Retourne (typeName, pkValues, isHashID, error)
func ParseFactID(id string) (typeName string, pkValues []string, isHashID bool, err error) {
	parts := strings.SplitN(id, IDSeparatorType, 2)
	if len(parts) != 2 {
		return "", nil, false, fmt.Errorf("format d'ID invalide: '%s'", id)
	}
	
	typeName = parts[0]
	valuesPart := parts[1]
	
	// Déterminer si c'est un hash (16 caractères hexadécimaux)
	if len(valuesPart) == IDHashLength && isHexString(valuesPart) {
		return typeName, []string{valuesPart}, true, nil
	}
	
	// Sinon, c'est une clé primaire composite
	rawValues := strings.Split(valuesPart, IDSeparatorValue)
	pkValues = make([]string, len(rawValues))
	for i, raw := range rawValues {
		pkValues[i] = unescapeIDValue(raw)
	}
	
	return typeName, pkValues, false, nil
}

// isHexString vérifie si une string est une chaîne hexadécimale valide.
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
```

### Étape 2 : Créer les Tests

**Fichier** : `constraint/id_generator_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestGenerateFactIDWithPrimaryKey(t *testing.T) {
	t.Log("🧪 TEST GENERATE FACT ID WITH PRIMARY KEY")
	t.Log("==========================================")
	
	tests := []struct {
		name       string
		fact       Fact
		typeDef    TypeDefinition
		wantID     string
		wantErr    bool
	}{
		{
			name: "clé primaire simple",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "login", Value: FactValue{Type: "string", Value: "alice"}},
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "User",
				Fields: []Field{
					{Name: "login", Type: "string", IsPrimaryKey: true},
					{Name: "name", Type: "string", IsPrimaryKey: false},
				},
			},
			wantID:  "User~alice",
			wantErr: false,
		},
		{
			name: "clé primaire composite",
			fact: Fact{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "firstName", Value: FactValue{Type: "string", Value: "Jean-Claude"}},
					{Name: "lastName", Value: FactValue{Type: "string", Value: "Pignon"}},
					{Name: "age", Value: FactValue{Type: "number", Value: float64(27)}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Person",
				Fields: []Field{
					{Name: "firstName", Type: "string", IsPrimaryKey: true},
					{Name: "lastName", Type: "string", IsPrimaryKey: true},
					{Name: "age", Type: "number", IsPrimaryKey: false},
				},
			},
			wantID:  "Person~Jean-Claude_Pignon",
			wantErr: false,
		},
		{
			name: "clé primaire avec number",
			fact: Fact{
				TypeName: "Product",
				Fields: []FactField{
					{Name: "code", Value: FactValue{Type: "number", Value: float64(12345)}},
					{Name: "name", Value: FactValue{Type: "string", Value: "Widget"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Product",
				Fields: []Field{
					{Name: "code", Type: "number", IsPrimaryKey: true},
					{Name: "name", Type: "string", IsPrimaryKey: false},
				},
			},
			wantID:  "Product~12345",
			wantErr: false,
		},
		{
			name: "clé primaire avec bool",
			fact: Fact{
				TypeName: "Flag",
				Fields: []FactField{
					{Name: "active", Value: FactValue{Type: "bool", Value: true}},
					{Name: "label", Value: FactValue{Type: "string", Value: "Test"}},
				},
			},
			typeDef: TypeDefinition{
				Name: "Flag",
				Fields: []Field{
					{Name: "active", Type: "bool", IsPrimaryKey: true},
					{Name: "label", Type: "string", IsPrimaryKey: false},
				},
			},
			wantID:  "Flag~true",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GenerateFactID(tt.fact, tt.typeDef)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else if id != tt.wantID {
					t.Errorf("❌ ID: attendu '%s', reçu '%s'", tt.wantID, id)
				} else {
					t.Logf("✅ ID généré: %s", id)
				}
			}
		})
	}
}

func TestGenerateFactIDWithHash(t *testing.T) {
	t.Log("🧪 TEST GENERATE FACT ID WITH HASH")
	t.Log("===================================")
	
	typeDef := TypeDefinition{
		Name: "Document",
		Fields: []Field{
			{Name: "title", Type: "string", IsPrimaryKey: false},
			{Name: "content", Type: "string", IsPrimaryKey: false},
		},
	}
	
	fact := Fact{
		TypeName: "Document",
		Fields: []FactField{
			{Name: "title", Value: FactValue{Type: "string", Value: "Doc1"}},
			{Name: "content", Value: FactValue{Type: "string", Value: "Content"}},
		},
	}
	
	id, err := GenerateFactID(fact, typeDef)
	if err != nil {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}
	
	// Vérifier le format: Document~<16 caractères hex>
	if !strings.HasPrefix(id, "Document~") {
		t.Errorf("❌ ID devrait commencer par 'Document~', reçu '%s'", id)
	}
	
	hashPart := strings.TrimPrefix(id, "Document~")
	if len(hashPart) != IDHashLength {
		t.Errorf("❌ Hash devrait avoir %d caractères, reçu %d", IDHashLength, len(hashPart))
	}
	
	if !isHexString(hashPart) {
		t.Errorf("❌ Hash devrait être hexadécimal, reçu '%s'", hashPart)
	}
	
	// Vérifier la reproductibilité (même fait = même hash)
	id2, err := GenerateFactID(fact, typeDef)
	if err != nil {
		t.Fatalf("❌ Erreur inattendue: %v", err)
	}
	
	if id != id2 {
		t.Errorf("❌ Hash non reproductible: '%s' != '%s'", id, id2)
	}
	
	t.Logf("✅ ID généré avec hash: %s", id)
}

func TestEscapeIDValue(t *testing.T) {
	t.Log("🧪 TEST ESCAPE ID VALUE")
	t.Log("========================")
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "pas de caractères spéciaux",
			input:    "alice",
			expected: "alice",
		},
		{
			name:     "avec tilde",
			input:    "user~123",
			expected: "user%7E123",
		},
		{
			name:     "avec underscore",
			input:    "first_last",
			expected: "first%5Flast",
		},
		{
			name:     "avec les deux",
			input:    "user~name_123",
			expected: "user%7Ename%5F123",
		},
		{
			name:     "avec pourcent",
			input:    "discount%20",
			expected: "discount%2520",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeIDValue(tt.input)
			if result != tt.expected {
				t.Errorf("❌ Attendu '%s', reçu '%s'", tt.expected, result)
			} else {
				t.Logf("✅ '%s' → '%s'", tt.input, result)
			}
			
			// Vérifier l'unescape
			unescaped := unescapeIDValue(result)
			if unescaped != tt.input {
				t.Errorf("❌ Unescape: attendu '%s', reçu '%s'", tt.input, unescaped)
			}
		})
	}
}

func TestParseFactID(t *testing.T) {
	t.Log("🧪 TEST PARSE FACT ID")
	t.Log("======================")
	
	tests := []struct {
		name           string
		id             string
		wantTypeName   string
		wantPKValues   []string
		wantIsHashID   bool
		wantErr        bool
	}{
		{
			name:         "clé primaire simple",
			id:           "User~alice",
			wantTypeName: "User",
			wantPKValues: []string{"alice"},
			wantIsHashID: false,
			wantErr:      false,
		},
		{
			name:         "clé primaire composite",
			id:           "Person~Jean-Claude_Pignon",
			wantTypeName: "Person",
			wantPKValues: []string{"Jean-Claude", "Pignon"},
			wantIsHashID: false,
			wantErr:      false,
		},
		{
			name:         "hash ID",
			id:           "Document~a3f5b9c2e1d4f8a7",
			wantTypeName: "Document",
			wantPKValues: []string{"a3f5b9c2e1d4f8a7"},
			wantIsHashID: true,
			wantErr:      false,
		},
		{
			name:    "format invalide",
			id:      "InvalidIDFormat",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeName, pkValues, isHashID, err := ParseFactID(tt.id)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}
			
			if err != nil {
				t.Errorf("❌ Erreur inattendue: %v", err)
				return
			}
			
			if typeName != tt.wantTypeName {
				t.Errorf("❌ Type: attendu '%s', reçu '%s'", tt.wantTypeName, typeName)
			}
			
			if isHashID != tt.wantIsHashID {
				t.Errorf("❌ IsHashID: attendu %v, reçu %v", tt.wantIsHashID, isHashID)
			}
			
			if len(pkValues) != len(tt.wantPKValues) {
				t.Errorf("❌ Nombre de valeurs: attendu %d, reçu %d", len(tt.wantPKValues), len(pkValues))
			} else {
				for i, want := range tt.wantPKValues {
					if pkValues[i] != want {
						t.Errorf("❌ Valeur[%d]: attendu '%s', reçu '%s'", i, want, pkValues[i])
					}
				}
			}
			
			t.Log("✅ Test réussi")
		})
	}
}

func TestValueToString(t *testing.T) {
	t.Log("🧪 TEST VALUE TO STRING")
	t.Log("========================")
	
	tests := []struct {
		name     string
		value    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "string",
			value:    "test",
			expected: "test",
			wantErr:  false,
		},
		{
			name:     "int",
			value:    42,
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "int64",
			value:    int64(123),
			expected: "123",
			wantErr:  false,
		},
		{
			name:     "float64",
			value:    float64(3.14),
			expected: "3.14",
			wantErr:  false,
		},
		{
			name:     "bool true",
			value:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "bool false",
			value:    false,
			expected: "false",
			wantErr:  false,
		},
		{
			name:    "nil",
			value:   nil,
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := valueToString(tt.value)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else if result != tt.expected {
					t.Errorf("❌ Attendu '%s', reçu '%s'", tt.expected, result)
				} else {
					t.Logf("✅ %v → '%s'", tt.value, result)
				}
			}
		})
	}
}
```

### Étape 3 : Intégrer dans constraint_facts.go

**Fichier** : `constraint/constraint_facts.go`

Modifier la fonction `ensureFactID` pour utiliser le générateur :

```go
// ensureFactID ensures a fact has an ID, generating one if necessary using primary keys or hash.
func ensureFactID(reteFact map[string]interface{}, fact Fact, typeDef TypeDefinition) (string, error) {
	// Check if ID was explicitly provided (should be prevented by validation)
	if id, exists := reteFact[FieldNameID]; exists {
		if idStr, ok := id.(string); ok && idStr != "" {
			// ID was provided, this should have been caught by validation
			// but we allow it for backward compatibility in some cases
			return idStr, nil
		}
	}

	// Generate ID based on primary key or hash
	id, err := GenerateFactID(fact, typeDef)
	if err != nil {
		return "", fmt.Errorf("génération d'ID pour le fait de type '%s': %v", fact.TypeName, err)
	}

	return id, nil
}
```

Modifier `ConvertFactsToReteFormat` pour passer le `TypeDefinition` :

```go
// ConvertFactsToReteFormat convertit les faits du Program vers le format RETE
func ConvertFactsToReteFormat(program Program) ([]map[string]interface{}, error) {
	reteFacts := []map[string]interface{}{}
	
	// Créer une map des types pour lookup rapide
	typeMap := make(map[string]TypeDefinition)
	for _, typeDef := range program.Types {
		typeMap[typeDef.Name] = typeDef
	}

	for i, fact := range program.Facts {
		reteFact := map[string]interface{}{
			FieldNameReteType: fact.TypeName,
		}

		// Convertir les champs
		convertFactFieldsToMap(fact.Fields, reteFact)

		// Récupérer la définition du type
		typeDef, exists := typeMap[fact.TypeName]
		if !exists {
			return nil, fmt.Errorf("fait %d: type '%s' non défini", i+1, fact.TypeName)
		}

		// Générer l'ID du fait
		factID, err := ensureFactID(reteFact, fact, typeDef)
		if err != nil {
			return nil, fmt.Errorf("fait %d: %v", i+1, err)
		}
		reteFact[FieldNameID] = factID

		reteFacts = append(reteFacts, reteFact)
	}

	return reteFacts, nil
}
```

## ✅ Validation

### Tests Automatiques

```bash
# Formattage
go fmt ./constraint/...
goimports -w constraint/

# Tests du générateur
go test -v ./constraint/ -run TestGenerateFactID
go test -v ./constraint/ -run TestEscapeIDValue
go test -v ./constraint/ -run TestParseFactID

# Tests complets
cd constraint
go test -v ./...

# Validation complète
cd ..
make test-unit
make validate
```

### Tests d'Intégration

Créer un fichier de test complet :

**Fichier** : `constraint/test/id_generation_integration.tsd`

```tsd
# Types avec clés primaires
type User(#login: string, name: string, age: number)
type Person(#firstName: string, #lastName: string, age: number)
type Product(#code: number, name: string, price: number)

# Type sans clé primaire
type Document(title: string, content: string)

# Faits
User(login: "alice", name: "Alice", age: 30)
User(login: "bob", name: "Bob", age: 25)

Person(firstName: "Jean-Claude", lastName: "Pignon", age: 27)
Person(firstName: "Marie", lastName: "Dupont", age: 32)

Product(code: 12345, name: "Widget", price: 9.99)

Document(title: "Doc1", content: "Some content")
Document(title: "Doc2", content: "Other content")
```

Vérifier que les ID générés sont corrects en inspectant le résultat du parsing.

### Checklist

- [ ] Fichier `id_generator.go` créé avec toutes les fonctions
- [ ] Tests complets dans `id_generator_test.go`
- [ ] Intégration dans `constraint_facts.go`
- [ ] Gestion des caractères spéciaux (escape/unescape)
- [ ] Hash MD5 implémenté et testé
- [ ] Reproductibilité du hash garantie
- [ ] `make validate` passe sans erreur
- [ ] Code formatté
- [ ] En-tête copyright présent
- [ ] Aucun hardcoding (constantes nommées)
- [ ] Messages d'erreur clairs

## 📝 Notes Importantes

### Choix de l'Algorithme de Hash

**MD5** a été choisi pour :
- Rapidité de calcul
- Déterminisme garanti
- Taille du hash (128 bits, 32 hex chars, tronqué à 16)
- Pas besoin de sécurité cryptographique ici

Alternative : SHA-256 si collision MD5 devient un problème.

### Gestion des Caractères Spéciaux

Les caractères `~` et `_` sont échappés car ils sont utilisés comme séparateurs.
L'échappement utilise le format URL encoding (`%XX`).

**Important** : Le caractère `%` lui-même doit être échappé en premier.

### Performance

- La génération d'ID doit être rapide (O(n) avec n = nombre de champs)
- Le hash MD5 est très rapide (< 1µs par fait)
- Pas d'allocation excessive de strings

### Reproductibilité

**CRITIQUE** : Le hash doit être identique pour les mêmes valeurs.

Garanties :
- Ordre des champs : ordre de définition du type (déterministe)
- Format des nombres : précision fixe pour float64
- Encodage : UTF-8 (Go par défaut)

## 🔄 Prochaines Étapes

Après validation de ce prompt :
1. Commit les changements :
   ```bash
   git add constraint/id_generator.go constraint/id_generator_test.go
   git add constraint/constraint_facts.go
   git commit -m "feat: implement automatic fact ID generation with primary keys and hash"
   ```
2. Passer au prompt **05-prompt-rete-integration.md**

## 📚 Références

- Génération d'ID : `constraint/id_generator.go`
- Conversion des faits : `constraint/constraint_facts.go`
- Standards du projet : `.github/prompts/common.md`
- Crypto MD5 : https://pkg.go.dev/crypto/md5

---

**Type** : Nouvelle fonctionnalité (génération d'ID)  
**Module** : constraint  
**Complexité** : Élevée  
**Temps estimé** : 60-90 minutes