// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/treivax/tsd/constraint"
)

// Ce fichier contient les fonctions de génération de représentations canoniques
// et de calcul de hash pour les conditions des chaînes alpha du système RETE.
//
// La représentation canonique garantit qu'une même condition logique produit
// toujours la même chaîne de caractères, permettant ainsi la comparaison et
// la déduplication des conditions.

// CanonicalString génère une représentation textuelle unique et déterministe d'une condition
// Format: "type(left,operator,right)"
func CanonicalString(condition SimpleCondition) string {
	leftStr := canonicalValue(condition.Left)
	rightStr := canonicalValue(condition.Right)

	return fmt.Sprintf("%s(%s,%s,%s)",
		condition.Type,
		leftStr,
		condition.Operator,
		rightStr,
	)
}

// canonicalValue convertit une valeur en représentation canonique
func canonicalValue(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		return canonicalMap(v)

	case constraint.FieldAccess:
		return fmt.Sprintf("fieldAccess(%s,%s)", v.Object, v.Field)

	case constraint.NumberLiteral:
		return fmt.Sprintf("literal(%v)", v.Value)

	case constraint.StringLiteral:
		return fmt.Sprintf("literal(%s)", v.Value)

	case constraint.BooleanLiteral:
		return fmt.Sprintf("literal(%t)", v.Value)

	case constraint.BinaryOperation:
		leftStr := canonicalValue(v.Left)
		rightStr := canonicalValue(v.Right)
		return fmt.Sprintf("binaryOp(%s,%s,%s)", leftStr, v.Operator, rightStr)

	case constraint.LogicalExpression:
		// Pour les expressions logiques, créer une représentation ordonnée
		parts := []string{canonicalValue(v.Left)}
		for _, op := range v.Operations {
			parts = append(parts, fmt.Sprintf("%s:%s", op.Op, canonicalValue(op.Right)))
		}
		return fmt.Sprintf("logical(%s)", strings.Join(parts, ","))

	case string:
		return fmt.Sprintf("string(%s)", v)

	case int, int32, int64:
		return fmt.Sprintf("int(%v)", v)

	case float32, float64:
		return fmt.Sprintf("float(%v)", v)

	case bool:
		return fmt.Sprintf("bool(%v)", v)

	case nil:
		return "nil"

	default:
		return fmt.Sprintf("unknown(%T:%v)", v, v)
	}
}

// canonicalMap convertit une map en représentation canonique triée
func canonicalMap(m map[string]interface{}) string {
	// Récupérer le type
	mapType, ok := m["type"].(string)
	if !ok {
		mapType = "map"
	}

	switch mapType {
	case "fieldAccess":
		obj, _ := m["object"].(string)
		field, _ := m["field"].(string)
		return fmt.Sprintf("fieldAccess(%s,%s)", obj, field)

	case "literal", "numberLiteral", "stringLiteral", "booleanLiteral":
		value := m["value"]
		return fmt.Sprintf("literal(%v)", value)

	case "binaryOperation", "binary_op", "comparison":
		leftStr := canonicalValue(m["left"])
		operator, _ := m["operator"].(string)
		if operator == "" {
			operator, _ = m["op"].(string)
		}
		rightStr := canonicalValue(m["right"])
		return fmt.Sprintf("binaryOp(%s,%s,%s)", leftStr, operator, rightStr)

	case "logicalExpression", "logical_op", "logicalExpr":
		leftStr := canonicalValue(m["left"])
		operations, ok := m["operations"].([]interface{})
		if !ok {
			return fmt.Sprintf("logical(%s)", leftStr)
		}

		parts := []string{leftStr}
		for _, opInterface := range operations {
			opMap, ok := opInterface.(map[string]interface{})
			if !ok {
				continue
			}
			op, _ := opMap["op"].(string)
			right := opMap["right"]
			parts = append(parts, fmt.Sprintf("%s:%s", op, canonicalValue(right)))
		}
		return fmt.Sprintf("logical(%s)", strings.Join(parts, ","))

	default:
		// Pour les maps génériques, trier les clés pour un ordre déterministe
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, fmt.Sprintf("%s:%s", k, canonicalValue(m[k])))
		}

		return fmt.Sprintf("%s{%s}", mapType, strings.Join(pairs, ","))
	}
}

// computeHash calcule le hash SHA-256 d'une condition
// Ce hash est utilisé pour l'égalité sémantique et la déduplication
func computeHash(condition SimpleCondition) string {
	canonical := CanonicalString(condition)
	hash := sha256.Sum256([]byte(canonical))
	return fmt.Sprintf("%x", hash)
}
