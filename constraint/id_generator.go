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

	factValues := fact.BuildFieldMap()
	pkValues := make([]string, 0, len(pkFields))

	for _, pkField := range pkFields {
		factValue, exists := factValues[pkField.Name]
		if !exists {
			return "", fmt.Errorf("champ de clé primaire '%s' manquant dans le fait", pkField.Name)
		}

		strValue, err := valueToString(factValue.Value)
		if err != nil {
			return "", fmt.Errorf("impossible de convertir le champ '%s' en string: %v", pkField.Name, err)
		}

		escapedValue := escapeIDValue(strValue)
		pkValues = append(pkValues, escapedValue)
	}

	return buildIDString(typeDef.Name, pkValues), nil
}

// generateIDFromHash génère un ID basé sur le hash de toutes les valeurs.
// Format: TypeName~<hash>
func generateIDFromHash(fact Fact, typeDef TypeDefinition) (string, error) {
	factValues := fact.BuildFieldMap()
	valueStrings := make([]string, 0, len(typeDef.Fields))

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

	hashStr := computeHash(valueStrings)
	return buildIDString(typeDef.Name, []string{hashStr}), nil
}

// valueToString convertit une valeur de fait en string.
// Pour les floats, utilise le format 'f' avec précision automatique (-1)
// pour garantir un format déterministe sans notation scientifique.
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
