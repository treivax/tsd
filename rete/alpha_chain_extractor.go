// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient des utilitaires pour extraire et analyser les conditions d'expressions complexes.
//
// Exemple d'utilisation:
//
//	// Expression AND: p.age > 18 AND p.salary >= 50000
//	expr := constraint.LogicalExpression{
//		Type: "logicalExpr",
//		Left: constraint.BinaryOperation{
//			Type:     "binaryOperation",
//			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "age"},
//			Operator: ">",
//			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 18},
//		},
//		Operations: []constraint.LogicalOperation{
//			{
//				Op: "AND",
//				Right: constraint.BinaryOperation{
//					Type:     "binaryOperation",
//					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "salary"},
//					Operator: ">=",
//					Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 50000},
//				},
//			},
//		},
//	}
//
//	// Extraire les conditions
//	conditions, opType, err := ExtractConditions(expr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Type d'opérateur: %s\n", opType) // Output: AND
//	fmt.Printf("Nombre de conditions: %d\n", len(conditions)) // Output: 2
//
//	// Générer des représentations canoniques
//	for _, cond := range conditions {
//		canonical := CanonicalString(cond)
//		fmt.Printf("Condition: %s\n", canonical)
//		fmt.Printf("Hash: %s\n", cond.Hash)
//	}
//
//	// Dédupliquer les conditions
//	uniqueConditions := DeduplicateConditions(conditions)
package rete

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/treivax/tsd/constraint"
)

// SimpleCondition représente une condition atomique extraite d'une expression complexe
type SimpleCondition struct {
	Type     string      `json:"type"`     // Type de condition: binaryOperation, comparison, arithmetic, etc.
	Left     interface{} `json:"left"`     // Opérande gauche
	Operator string      `json:"operator"` // Opérateur
	Right    interface{} `json:"right"`    // Opérande droite
	Hash     string      `json:"hash"`     // Hash unique calculé automatiquement
}

// DecomposedCondition extends SimpleCondition with decomposition metadata
// for supporting intermediate result propagation in alpha chains
type DecomposedCondition struct {
	SimpleCondition
	ResultName   string   `json:"result_name,omitempty"`  // Name of intermediate result produced (e.g., "temp_1")
	Dependencies []string `json:"dependencies,omitempty"` // Required intermediate results
	IsAtomic     bool     `json:"is_atomic,omitempty"`    // true if atomic operation
}

// NewSimpleCondition crée une nouvelle condition simple avec hash calculé
func NewSimpleCondition(condType string, left interface{}, operator string, right interface{}) SimpleCondition {
	cond := SimpleCondition{
		Type:     condType,
		Left:     left,
		Operator: operator,
		Right:    right,
	}
	cond.Hash = computeHash(cond)
	return cond
}

// computeHash calcule le hash SHA-256 d'une condition
func computeHash(condition SimpleCondition) string {
	canonical := CanonicalString(condition)
	hash := sha256.Sum256([]byte(canonical))
	return fmt.Sprintf("%x", hash)
}

// ExtractConditions extrait toutes les conditions simples d'une expression complexe
// Retourne: liste de conditions, type d'opérateur principal (AND/OR/SINGLE), erreur
func ExtractConditions(expr interface{}) ([]SimpleCondition, string, error) {
	switch e := expr.(type) {
	case map[string]interface{}:
		return extractFromMap(e)

	case constraint.BinaryOperation:
		cond := NewSimpleCondition("binaryOperation", e.Left, e.Operator, e.Right)
		return []SimpleCondition{cond}, "SINGLE", nil

	case constraint.LogicalExpression:
		return extractFromLogicalExpression(e)

	case constraint.Constraint:
		return extractFromConstraint(e)

	case constraint.FieldAccess:
		// Un accès de champ seul n'est pas une condition
		return []SimpleCondition{}, "NONE", nil

	case constraint.NumberLiteral, constraint.StringLiteral, constraint.BooleanLiteral:
		// Les littéraux seuls ne sont pas des conditions
		return []SimpleCondition{}, "NONE", nil

	default:
		return nil, "", fmt.Errorf("type d'expression non supporté: %T", expr)
	}
}

// extractFromMap extrait les conditions d'une expression sous forme de map
func extractFromMap(expr map[string]interface{}) ([]SimpleCondition, string, error) {
	exprType, ok := expr["type"].(string)
	if !ok {
		return nil, "", fmt.Errorf("type d'expression manquant")
	}

	switch exprType {
	case "binaryOperation", "binary_op", "comparison":
		operator, ok := expr["operator"].(string)
		if !ok {
			if operator, ok = expr["op"].(string); !ok {
				return nil, "", fmt.Errorf("opérateur manquant")
			}
		}
		cond := NewSimpleCondition("binaryOperation", expr["left"], operator, expr["right"])
		return []SimpleCondition{cond}, "SINGLE", nil

	case "logicalExpression", "logical_op", "logicalExpr":
		return extractFromLogicalExpressionMap(expr)

	case "constraint":
		// Les contraintes dans les maps peuvent avoir différentes structures
		// On essaie d'extraire left/operator/right directement
		if left, ok := expr["left"]; ok {
			operator, _ := expr["operator"].(string)
			right := expr["right"]
			cond := NewSimpleCondition("constraint", left, operator, right)
			return []SimpleCondition{cond}, "SINGLE", nil
		}
		return []SimpleCondition{}, "NONE", nil

	case "fieldAccess":
		return []SimpleCondition{}, "NONE", nil

	case "literal", "numberLiteral", "stringLiteral", "booleanLiteral":
		return []SimpleCondition{}, "NONE", nil

	default:
		return nil, "", fmt.Errorf("type d'expression map non supporté: %s", exprType)
	}
}

// extractFromLogicalExpression extrait les conditions d'une expression logique
func extractFromLogicalExpression(expr constraint.LogicalExpression) ([]SimpleCondition, string, error) {
	allConditions := []SimpleCondition{}
	operatorType := ""

	// Extraire les conditions du côté gauche
	leftConds, _, err := ExtractConditions(expr.Left)
	if err != nil {
		return nil, "", fmt.Errorf("erreur extraction left: %w", err)
	}
	allConditions = append(allConditions, leftConds...)

	// Traiter toutes les opérations
	for _, op := range expr.Operations {
		rightConds, _, err := ExtractConditions(op.Right)
		if err != nil {
			return nil, "", fmt.Errorf("erreur extraction right: %w", err)
		}
		allConditions = append(allConditions, rightConds...)

		// Déterminer le type d'opérateur principal
		if operatorType == "" {
			operatorType = op.Op
		} else if operatorType != op.Op {
			// Mélange d'opérateurs - retourner "MIXED"
			operatorType = "MIXED"
		}
	}

	if operatorType == "" {
		operatorType = "SINGLE"
	}

	return allConditions, operatorType, nil
}

// extractFromLogicalExpressionMap extrait les conditions d'une expression logique (format map)
func extractFromLogicalExpressionMap(expr map[string]interface{}) ([]SimpleCondition, string, error) {
	allConditions := []SimpleCondition{}
	operatorType := ""

	// Extraire les conditions du côté gauche
	left, ok := expr["left"]
	if !ok {
		return nil, "", fmt.Errorf("left manquant dans logicalExpression")
	}

	leftConds, _, err := ExtractConditions(left)
	if err != nil {
		return nil, "", fmt.Errorf("erreur extraction left: %w", err)
	}
	allConditions = append(allConditions, leftConds...)

	// Traiter toutes les opérations
	operations, ok := expr["operations"]
	if !ok {
		return allConditions, "SINGLE", nil
	}

	// Supporter []interface{}, []map[string]interface{} et []constraint.LogicalOperation
	// Essayer []map[string]interface{} en premier (type le plus courant du parser)
	if opsMapList, ok := operations.([]map[string]interface{}); ok {
		for _, opMap := range opsMapList {
			op, ok := opMap["op"].(string)
			if !ok {
				return nil, "", fmt.Errorf("op manquant dans operation")
			}

			right, ok := opMap["right"]
			if !ok {
				return nil, "", fmt.Errorf("right manquant dans operation")
			}

			rightConds, _, err := ExtractConditions(right)
			if err != nil {
				return nil, "", fmt.Errorf("erreur extraction right: %w", err)
			}
			allConditions = append(allConditions, rightConds...)

			// Déterminer le type d'opérateur principal
			if operatorType == "" {
				operatorType = op
			} else if operatorType != op {
				operatorType = "MIXED"
			}
		}
	} else if opsList, ok := operations.([]interface{}); ok {
		for _, opInterface := range opsList {
			opMap, ok := opInterface.(map[string]interface{})
			if !ok {
				return nil, "", fmt.Errorf("operation doit être une map")
			}

			op, ok := opMap["op"].(string)
			if !ok {
				return nil, "", fmt.Errorf("op manquant dans operation")
			}

			right, ok := opMap["right"]
			if !ok {
				return nil, "", fmt.Errorf("right manquant dans operation")
			}

			rightConds, _, err := ExtractConditions(right)
			if err != nil {
				return nil, "", fmt.Errorf("erreur extraction right: %w", err)
			}
			allConditions = append(allConditions, rightConds...)

			// Déterminer le type d'opérateur principal
			if operatorType == "" {
				operatorType = op
			} else if operatorType != op {
				operatorType = "MIXED"
			}
		}
	} else if logicalOps, ok := operations.([]constraint.LogicalOperation); ok {
		// Supporter le type constraint.LogicalOperation directement
		for _, op := range logicalOps {
			rightConds, _, err := ExtractConditions(op.Right)
			if err != nil {
				return nil, "", fmt.Errorf("erreur extraction right: %w", err)
			}
			allConditions = append(allConditions, rightConds...)

			// Déterminer le type d'opérateur principal
			if operatorType == "" {
				operatorType = op.Op
			} else if operatorType != op.Op {
				operatorType = "MIXED"
			}
		}
	} else {
		return nil, "", fmt.Errorf("operations doit être un tableau ([]interface{}, []map[string]interface{} ou []LogicalOperation)")
	}

	if operatorType == "" {
		operatorType = "SINGLE"
	}

	return allConditions, operatorType, nil
}

// extractFromNOTConstraint extrait les conditions d'une contrainte NOT
func extractFromNOTConstraint(expr constraint.NotConstraint) ([]SimpleCondition, string, error) {
	// Pour les contraintes NOT, on retourne une condition spéciale
	// qui sera gérée différemment par le constructeur de chaîne
	cond := NewSimpleCondition("not", expr.Expression, "NOT", nil)
	return []SimpleCondition{cond}, "NOT", nil
}

// extractFromNOTConstraintMap extrait les conditions d'une contrainte NOT (format map)
func extractFromNOTConstraintMap(expr map[string]interface{}) ([]SimpleCondition, string, error) {
	expression, ok := expr["expression"]
	if !ok {
		return nil, "", fmt.Errorf("expression manquant dans notConstraint")
	}

	// Pour les contraintes NOT, on retourne une condition spéciale
	cond := NewSimpleCondition("not", expression, "NOT", nil)
	return []SimpleCondition{cond}, "NOT", nil
}

// extractFromConstraint extrait les conditions d'une contrainte
func extractFromConstraint(c constraint.Constraint) ([]SimpleCondition, string, error) {
	// Les contraintes peuvent avoir left/operator/right directement
	if c.Left != nil && c.Operator != "" {
		cond := NewSimpleCondition("constraint", c.Left, c.Operator, c.Right)
		return []SimpleCondition{cond}, "SINGLE", nil
	}
	return []SimpleCondition{}, "NONE", nil
}

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

// CompareConditions compare deux conditions pour l'égalité
func CompareConditions(c1, c2 SimpleCondition) bool {
	return c1.Hash == c2.Hash
}

// DeduplicateConditions supprime les conditions dupliquées d'une liste
func DeduplicateConditions(conditions []SimpleCondition) []SimpleCondition {
	seen := make(map[string]bool)
	result := []SimpleCondition{}

	for _, cond := range conditions {
		if !seen[cond.Hash] {
			seen[cond.Hash] = true
			result = append(result, cond)
		}
	}

	return result
}

// IsCommutative retourne true si l'opérateur est commutatif
// Les opérateurs commutatifs (AND, OR, +, *, ==, !=) peuvent être réordonnés
// Les opérateurs non-commutatifs (-, /, <, >, <=, >=, séquences) doivent préserver l'ordre
func IsCommutative(operator string) bool {
	commutativeOps := map[string]bool{
		"AND": true,
		"OR":  true,
		"&&":  true,
		"||":  true,
		"+":   true,
		"*":   true,
		"==":  true,
		"!=":  true,
		"<>":  true,
	}
	return commutativeOps[operator]
}

// NormalizeConditions trie les conditions dans un ordre canonique déterministe
// Si l'opérateur est commutatif (AND, OR), les conditions sont triées
// Si l'opérateur est non-commutatif, l'ordre est préservé
func NormalizeConditions(conditions []SimpleCondition, operator string) []SimpleCondition {
	// Si pas de conditions ou une seule condition, retourner tel quel
	if len(conditions) <= 1 {
		return conditions
	}

	// Si l'opérateur n'est pas commutatif, préserver l'ordre original
	if !IsCommutative(operator) {
		return conditions
	}

	// Créer une copie pour ne pas modifier l'original
	normalized := make([]SimpleCondition, len(conditions))
	copy(normalized, conditions)

	// Trier par représentation canonique pour un ordre déterministe
	sort.Slice(normalized, func(i, j int) bool {
		return CanonicalString(normalized[i]) < CanonicalString(normalized[j])
	})

	return normalized
}

// NormalizeExpression normalise une expression en appliquant un ordre canonique
// aux conditions quand l'opérateur est commutatif
// NormalizeORExpression normalise une expression OR en triant ses termes dans un ordre canonique.
// Cette fonction extrait les termes OR, les trie et reconstruit l'expression normalisée
// SANS la décomposer en chaîne - l'expression OR reste un seul nœud atomique.
//
// Paramètres:
//   - expr: expression à normaliser (LogicalExpression ou map avec operations OR)
//
// Retourne:
//   - interface{}: expression OR normalisée avec termes triés
//   - error: erreur si l'expression n'est pas une expression OR valide
//
// Exemple:
//
//	Input:  p.status == "VIP" OR p.age > 18
//	Output: p.age > 18 OR p.status == "VIP"  (ordre alphabétique des champs)
func NormalizeORExpression(expr interface{}) (interface{}, error) {
	if expr == nil {
		return nil, fmt.Errorf("expression nil")
	}

	// Vérifier que c'est bien une expression OR
	exprType, err := AnalyzeExpression(expr)
	if err != nil {
		return nil, fmt.Errorf("erreur analyse expression: %w", err)
	}

	if exprType != ExprTypeOR && exprType != ExprTypeMixed {
		return nil, fmt.Errorf("expression n'est pas de type OR ou Mixed: %s", exprType)
	}

	switch e := expr.(type) {
	case constraint.LogicalExpression:
		return normalizeORLogicalExpression(e)
	case map[string]interface{}:
		return normalizeORExpressionMap(e)
	default:
		return nil, fmt.Errorf("type d'expression non supporté pour normalisation OR: %T", expr)
	}
}

// normalizeORLogicalExpression normalise une LogicalExpression contenant des OR
func normalizeORLogicalExpression(expr constraint.LogicalExpression) (constraint.LogicalExpression, error) {
	// Extraire tous les termes OR
	terms := []interface{}{expr.Left}

	for _, op := range expr.Operations {
		opStr := strings.ToUpper(op.Op)
		if opStr == "OR" || opStr == "||" {
			terms = append(terms, op.Right)
		}
	}

	// Convertir chaque terme en string canonique pour le tri
	type termWithCanonical struct {
		term      interface{}
		canonical string
	}

	termsWithCanonical := make([]termWithCanonical, len(terms))
	for i, term := range terms {
		// Créer une représentation canonique pour le tri
		canonical := canonicalValue(term)
		termsWithCanonical[i] = termWithCanonical{
			term:      term,
			canonical: canonical,
		}
	}

	// Trier par représentation canonique
	sort.Slice(termsWithCanonical, func(i, j int) bool {
		return termsWithCanonical[i].canonical < termsWithCanonical[j].canonical
	})

	// Reconstruire l'expression avec les termes triés
	if len(termsWithCanonical) == 0 {
		return constraint.LogicalExpression{}, fmt.Errorf("aucun terme trouvé")
	}

	normalized := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: termsWithCanonical[0].term,
	}

	// Reconstruire les opérations OR
	for i := 1; i < len(termsWithCanonical); i++ {
		normalized.Operations = append(normalized.Operations, constraint.LogicalOperation{
			Op:    "OR",
			Right: termsWithCanonical[i].term,
		})
	}

	return normalized, nil
}

// normalizeORExpressionMap normalise une expression OR au format map
func normalizeORExpressionMap(expr map[string]interface{}) (map[string]interface{}, error) {
	// Extraire tous les termes OR
	terms := []interface{}{}

	if left, ok := expr["left"]; ok {
		terms = append(terms, left)
	}

	if operations, ok := expr["operations"]; ok {
		// Supporter différents formats d'opérations
		if opsList, ok := operations.([]interface{}); ok {
			for _, opInterface := range opsList {
				if opMap, ok := opInterface.(map[string]interface{}); ok {
					if op, ok := opMap["op"].(string); ok {
						opStr := strings.ToUpper(op)
						if opStr == "OR" || opStr == "||" {
							if right, ok := opMap["right"]; ok {
								terms = append(terms, right)
							}
						}
					}
				}
			}
		} else if opsMapList, ok := operations.([]map[string]interface{}); ok {
			for _, opMap := range opsMapList {
				if op, ok := opMap["op"].(string); ok {
					opStr := strings.ToUpper(op)
					if opStr == "OR" || opStr == "||" {
						if right, ok := opMap["right"]; ok {
							terms = append(terms, right)
						}
					}
				}
			}
		}
	}

	if len(terms) == 0 {
		return nil, fmt.Errorf("aucun terme OR trouvé")
	}

	// Convertir chaque terme en string canonique pour le tri
	type termWithCanonical struct {
		term      interface{}
		canonical string
	}

	termsWithCanonical := make([]termWithCanonical, len(terms))
	for i, term := range terms {
		canonical := canonicalValue(term)
		termsWithCanonical[i] = termWithCanonical{
			term:      term,
			canonical: canonical,
		}
	}

	// Trier par représentation canonique
	sort.Slice(termsWithCanonical, func(i, j int) bool {
		return termsWithCanonical[i].canonical < termsWithCanonical[j].canonical
	})

	// Reconstruire l'expression map avec les termes triés
	normalized := map[string]interface{}{
		"type": expr["type"],
		"left": termsWithCanonical[0].term,
	}

	// Reconstruire les opérations OR
	normalizedOps := make([]map[string]interface{}, 0)
	for i := 1; i < len(termsWithCanonical); i++ {
		normalizedOps = append(normalizedOps, map[string]interface{}{
			"op":    "OR",
			"right": termsWithCanonical[i].term,
		})
	}

	if len(normalizedOps) > 0 {
		normalized["operations"] = normalizedOps
	}

	return normalized, nil
}

func NormalizeExpression(expr interface{}) (interface{}, error) {
	switch e := expr.(type) {
	case constraint.LogicalExpression:
		return normalizeLogicalExpression(e)

	case constraint.BinaryOperation:
		// Les opérations binaires simples ne nécessitent pas de normalisation
		// mais on peut normaliser récursivement les sous-expressions
		return e, nil

	case constraint.Constraint:
		// Les contraintes simples ne nécessitent pas de normalisation
		return e, nil

	case map[string]interface{}:
		return normalizeExpressionMap(e)

	default:
		// Pour les autres types (literals, field access), retourner tel quel
		return expr, nil
	}
}

// normalizeLogicalExpression normalise une expression logique
func normalizeLogicalExpression(expr constraint.LogicalExpression) (constraint.LogicalExpression, error) {
	// Si pas d'opérations, retourner tel quel
	if len(expr.Operations) == 0 {
		return expr, nil
	}

	// Déterminer si tous les opérateurs sont identiques et commutatifs
	firstOp := expr.Operations[0].Op
	allSame := true
	for _, op := range expr.Operations {
		if op.Op != firstOp {
			allSame = false
			break
		}
	}

	// Si les opérateurs ne sont pas tous identiques ou si non-commutatif, retourner tel quel
	if !allSame || !IsCommutative(firstOp) {
		return expr, nil
	}

	// Extraire toutes les conditions
	conditions, _, err := extractFromLogicalExpression(expr)
	if err != nil {
		return expr, err
	}

	// Normaliser l'ordre des conditions
	normalized := NormalizeConditions(conditions, firstOp)

	// Reconstruire l'expression logique avec les conditions normalisées
	if len(normalized) == 0 {
		return expr, nil
	}

	// Reconstruire l'expression avec les conditions normalisées
	rebuiltExpr, err := rebuildLogicalExpression(normalized, firstOp)
	if err != nil {
		return expr, err
	}

	return rebuiltExpr, nil
}

// rebuildLogicalExpression reconstruit une expression logique à partir de conditions normalisées
func rebuildLogicalExpression(conditions []SimpleCondition, operator string) (constraint.LogicalExpression, error) {
	if len(conditions) == 0 {
		return constraint.LogicalExpression{}, fmt.Errorf("cannot rebuild expression from empty conditions")
	}

	// Cas simple : une seule condition
	if len(conditions) == 1 {
		cond := conditions[0]
		return constraint.LogicalExpression{
			Type:       "logicalExpr",
			Left:       rebuildConditionAsExpression(cond),
			Operations: []constraint.LogicalOperation{},
		}, nil
	}

	// Cas avec plusieurs conditions : créer une chaîne d'opérations
	// Le premier élément devient Left, les autres deviennent Operations
	rebuiltExpr := constraint.LogicalExpression{
		Type:       "logicalExpr",
		Left:       rebuildConditionAsExpression(conditions[0]),
		Operations: make([]constraint.LogicalOperation, 0, len(conditions)-1),
	}

	// Ajouter les conditions restantes comme opérations
	for i := 1; i < len(conditions); i++ {
		rebuiltExpr.Operations = append(rebuiltExpr.Operations, constraint.LogicalOperation{
			Op:    operator,
			Right: rebuildConditionAsExpression(conditions[i]),
		})
	}

	return rebuiltExpr, nil
}

// rebuildConditionAsExpression convertit une SimpleCondition en expression utilisable
func rebuildConditionAsExpression(cond SimpleCondition) interface{} {
	// Créer une BinaryOperation à partir de la SimpleCondition
	return constraint.BinaryOperation{
		Type:     cond.Type,
		Left:     cond.Left,
		Operator: cond.Operator,
		Right:    cond.Right,
	}
}

// rebuildLogicalExpressionMap reconstruit une expression map à partir de conditions normalisées
func rebuildLogicalExpressionMap(conditions []SimpleCondition, operator string) (map[string]interface{}, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("cannot rebuild expression from empty conditions")
	}

	// Cas simple : une seule condition
	if len(conditions) == 1 {
		cond := conditions[0]
		return map[string]interface{}{
			"type":       "logicalExpr",
			"left":       rebuildConditionAsMap(cond),
			"operations": []interface{}{},
		}, nil
	}

	// Cas avec plusieurs conditions
	operations := make([]interface{}, 0, len(conditions)-1)
	for i := 1; i < len(conditions); i++ {
		operations = append(operations, map[string]interface{}{
			"op":    operator,
			"right": rebuildConditionAsMap(conditions[i]),
		})
	}

	return map[string]interface{}{
		"type":       "logicalExpr",
		"left":       rebuildConditionAsMap(conditions[0]),
		"operations": operations,
	}, nil
}

// rebuildConditionAsMap convertit une SimpleCondition en map
func rebuildConditionAsMap(cond SimpleCondition) map[string]interface{} {
	return map[string]interface{}{
		"type":     cond.Type,
		"left":     cond.Left,
		"operator": cond.Operator,
		"right":    cond.Right,
	}
}

// normalizeExpressionMap normalise une expression sous forme de map
func normalizeExpressionMap(expr map[string]interface{}) (map[string]interface{}, error) {
	exprType, ok := expr["type"].(string)
	if !ok {
		return expr, nil
	}

	switch exprType {
	case "logicalExpression", "logical_op", "logicalExpr":
		// Extraire les conditions
		conditions, opType, err := extractFromLogicalExpressionMap(expr)
		if err != nil {
			return expr, err
		}

		// Si l'opérateur est commutatif, normaliser et reconstruire
		if IsCommutative(opType) {
			normalized := NormalizeConditions(conditions, opType)
			rebuiltExpr, err := rebuildLogicalExpressionMap(normalized, opType)
			if err != nil {
				return expr, err
			}
			return rebuiltExpr, nil
		}
		return expr, nil

	case "binaryOperation", "binary_op", "comparison":
		// Les opérations binaires simples ne nécessitent pas de normalisation
		return expr, nil

	default:
		return expr, nil
	}
}
