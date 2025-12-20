// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
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

// FactContext contient le contexte pour la génération d'IDs.
// Il permet de résoudre les références de variables vers leurs IDs.
type FactContext struct {
	// VariableIDs map les noms de variables vers les IDs de faits
	VariableIDs map[string]string

	// TypeMap map les noms de types vers leurs définitions
	TypeMap map[string]TypeDefinition
}

// NewFactContext crée un nouveau contexte de génération d'IDs.
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

// RegisterVariable enregistre l'ID d'une variable de fait.
func (fc *FactContext) RegisterVariable(varName, factID string) {
	fc.VariableIDs[varName] = factID
}

// ResolveVariable résout une variable vers son ID.
func (fc *FactContext) ResolveVariable(varName string) (string, error) {
	id, exists := fc.VariableIDs[varName]
	if !exists {
		return "", fmt.Errorf("variable '%s' non définie dans le contexte", varName)
	}
	return id, nil
}

// GenerateFactID génère l'identifiant d'un fait selon ses clés primaires.
// Si le type a une clé primaire définie, l'ID est : TypeName~value1_value2_...
// Sinon, l'ID est : TypeName~<hash>
// Le contexte est optionnel pour la rétrocompatibilité.
func GenerateFactID(fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
	// Créer un contexte vide si non fourni (rétrocompatibilité)
	if ctx == nil {
		ctx = NewFactContext(nil)
	}

	// Vérifier si le type a une clé primaire
	if typeDef.HasPrimaryKey() {
		return generateIDFromPrimaryKey(fact, typeDef, ctx)
	}

	// Sinon, générer un ID par hash
	return generateIDFromHash(fact, typeDef, ctx)
}

// GenerateFactIDWithoutContext génère un ID sans contexte de variables.
// Deprecated: Utiliser GenerateFactID avec FactContext pour supporter les références de faits.
func GenerateFactIDWithoutContext(fact Fact, typeDef TypeDefinition) (string, error) {
	return GenerateFactID(fact, typeDef, nil)
}

// generateIDFromPrimaryKey génère un ID basé sur les valeurs de clé primaire.
// Format: TypeName~value1_value2_..._valueN
// Gère les types primitifs ET les références de faits.
func generateIDFromPrimaryKey(fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
	pkFields := typeDef.GetPrimaryKeyFields()
	if len(pkFields) == 0 {
		return "", fmt.Errorf("type '%s' n'a pas de clé primaire définie", typeDef.Name)
	}

	factValues := fact.BuildFieldMap()
	pkValues := make([]string, 0, len(pkFields))

	for _, pkField := range pkFields {
		factValue, exists := factValues[pkField.Name]
		if !exists {
			return "", fmt.Errorf("champ de clé primaire '%s' manquant dans le fait", pkField.Name)
		}

		strValue, err := convertFieldValueToString(factValue, pkField, ctx)
		if err != nil {
			return "", fmt.Errorf("conversion du champ '%s': %v", pkField.Name, err)
		}

		escapedValue := escapeIDValue(strValue)
		pkValues = append(pkValues, escapedValue)
	}

	return buildIDString(typeDef.Name, pkValues), nil
}

// generateIDFromHash génère un ID basé sur le hash de toutes les valeurs.
// Format: TypeName~<hash>
// Gère les références de faits.
func generateIDFromHash(fact Fact, typeDef TypeDefinition, ctx *FactContext) (string, error) {
	factValues := fact.BuildFieldMap()
	valueStrings := make([]string, 0, len(typeDef.Fields))

	// Trier les champs par nom pour garantir le déterminisme
	sortedFields := make([]Field, len(typeDef.Fields))
	copy(sortedFields, typeDef.Fields)
	sort.Slice(sortedFields, func(i, j int) bool {
		return sortedFields[i].Name < sortedFields[j].Name
	})

	for _, field := range sortedFields {
		factValue, exists := factValues[field.Name]
		if exists && factValue.Value != nil {
			strValue, err := convertFieldValueToString(factValue, field, ctx)
			if err != nil {
				return "", fmt.Errorf("champ '%s': %v", field.Name, err)
			}
			valueStrings = append(valueStrings, field.Name+"="+strValue)
		}
	}

	hashStr := computeHash(valueStrings)
	return buildIDString(typeDef.Name, []string{hashStr}), nil
}

// convertFieldValueToString convertit une valeur de champ en string.
// Gère les types primitifs ET les références de faits (variableReference).
func convertFieldValueToString(value FactValue, field Field, ctx *FactContext) (string, error) {
	actualValue := value.Unwrap()

	switch value.Type {
	case ValueTypeString, ValueTypeIdentifier:
		return convertStringValue(actualValue)
	case ValueTypeNumber:
		return convertNumberValue(actualValue)
	case ValueTypeBoolean, ValueTypeBool:
		return convertBooleanValue(actualValue)
	case ValueTypeVariableReference:
		return resolveVariableReference(actualValue, ctx)
	default:
		return "", fmt.Errorf("type de valeur non supporté: %s", value.Type)
	}
}

// convertStringValue converts a string value to its string representation
func convertStringValue(value interface{}) (string, error) {
	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("valeur string attendue, reçu %T", value)
	}
	return str, nil
}

// convertNumberValue converts a numeric value to its string representation
func convertNumberValue(value interface{}) (string, error) {
	switch num := value.(type) {
	case float64:
		return formatNumber(num), nil
	case int:
		return strconv.Itoa(num), nil
	case int64:
		return strconv.FormatInt(num, 10), nil
	default:
		return "", fmt.Errorf("valeur number attendue, reçu %T", value)
	}
}

// convertBooleanValue converts a boolean value to its string representation
func convertBooleanValue(value interface{}) (string, error) {
	b, ok := value.(bool)
	if !ok {
		return "", fmt.Errorf("valeur boolean attendue, reçu %T", value)
	}
	if b {
		return "true", nil
	}
	return "false", nil
}

// resolveVariableReference resolves a variable reference to its fact ID
func resolveVariableReference(value interface{}, ctx *FactContext) (string, error) {
	if ctx == nil {
		return "", errors.New("contexte requis pour résoudre les variables")
	}

	varName, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("nom de variable attendu, reçu %T", value)
	}

	id, err := ctx.ResolveVariable(varName)
	if err != nil {
		return "", fmt.Errorf("résolution de variable '%s': %v", varName, err)
	}

	return id, nil
}

// formatNumber formate un nombre pour l'ID (sans .0 pour les entiers).
func formatNumber(n float64) string {
	if n == float64(int64(n)) {
		return fmt.Sprintf("%d", int64(n))
	}
	return strconv.FormatFloat(n, 'f', -1, 64)
}

// valueToString convertit une valeur de fait en string.
// Pour les floats, utilise le format 'f' avec précision automatique (-1)
// pour garantir un format déterministe sans notation scientifique.
// Deprecated: Utiliser convertFieldValueToString pour supporter les références de faits.
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
		// Format: 'f' = decimal notation (no exponent)
		// Precision: -1 = smallest number of digits necessary
		// This ensures deterministic string representation
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// escapeIDValue échappe les caractères spéciaux dans une valeur d'ID.
// Remplace ~ par %7E, _ par %5F, et espace par %20 (URL encoding partiel)
func escapeIDValue(value string) string {
	value = strings.ReplaceAll(value, "%", "%25") // % en premier pour éviter double-escape
	value = strings.ReplaceAll(value, IDSeparatorType, "%7E")
	value = strings.ReplaceAll(value, IDSeparatorValue, "%5F")
	value = strings.ReplaceAll(value, " ", "%20") // Échapper les espaces
	return value
}

// unescapeIDValue inverse l'échappement des caractères spéciaux.
func unescapeIDValue(value string) string {
	value = strings.ReplaceAll(value, "%20", " ") // Unescaper les espaces
	value = strings.ReplaceAll(value, "%5F", IDSeparatorValue)
	value = strings.ReplaceAll(value, "%7E", IDSeparatorType)
	value = strings.ReplaceAll(value, "%25", "%") // % en dernier pour éviter double-unescape
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
	if typeName == "" {
		return "", nil, false, fmt.Errorf("nom de type vide dans l'ID: '%s'", id)
	}

	valuesPart := parts[1]
	if valuesPart == "" {
		return "", nil, false, fmt.Errorf("partie valeurs vide dans l'ID: '%s'", id)
	}

	if len(valuesPart) == IDHashLength && isHexString(valuesPart) {
		return typeName, []string{valuesPart}, true, nil
	}

	rawValues := strings.Split(valuesPart, IDSeparatorValue)
	pkValues = make([]string, len(rawValues))
	for i, raw := range rawValues {
		pkValues[i] = unescapeIDValue(raw)
	}

	return typeName, pkValues, false, nil
}

// computeHash calcule le hash MD5 des valeurs et le tronque.
func computeHash(valueStrings []string) string {
	concatenated := strings.Join(valueStrings, "|")
	hash := md5.Sum([]byte(concatenated))
	hashStr := hex.EncodeToString(hash[:])

	if len(hashStr) > IDHashLength {
		return hashStr[:IDHashLength]
	}
	return hashStr
}

// buildIDString construit l'ID final à partir du type et des valeurs.
func buildIDString(typeName string, values []string) string {
	return typeName + IDSeparatorType + strings.Join(values, IDSeparatorValue)
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
