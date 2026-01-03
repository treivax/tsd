// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// FieldExtractor extrait les noms de champs depuis diverses structures AST.
//
// Cette interface abstraite permet de supporter différents formats de conditions
// (alpha, beta, terminal) de manière uniforme.
type FieldExtractor interface {
	// ExtractFields extrait les champs depuis une condition/expression
	ExtractFields(condition interface{}) ([]string, error)
}

// AlphaConditionExtractor extrait les champs depuis les conditions de nœuds alpha.
type AlphaConditionExtractor struct{}

// ExtractFields extrait les champs depuis une condition alpha.
//
// Les conditions alpha sont typiquement des expressions binaires ou comparaisons
// sur des champs d'un fait.
//
// Exemples :
//   - "product.price > 100" → ["price"]
//   - "order.status == 'active' && order.total > 50" → ["status", "total"]
//
// Paramètres :
//   - condition : structure de condition (map[string]interface{} depuis parser)
//
// Retourne la liste des noms de champs (dédupliqués).
func (ace *AlphaConditionExtractor) ExtractFields(condition interface{}) ([]string, error) {
	fields := make(map[string]bool)

	if err := extractFieldsRecursive(condition, fields); err != nil {
		return nil, err
	}

	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}

	return result, nil
}

// BetaConditionExtractor extrait les champs depuis les conditions de jointure beta.
type BetaConditionExtractor struct{}

// ExtractFields extrait les champs depuis une condition de jointure beta.
//
// Les conditions beta sont des comparaisons entre champs de différents faits.
//
// Exemple :
//   - "order.customer_id == customer.id" → {"order": ["customer_id"], "customer": ["id"]}
//
// Note : Cette version retourne une liste plate de tous les champs.
// Une version future pourrait retourner une map factType → fields.
func (bce *BetaConditionExtractor) ExtractFields(condition interface{}) ([]string, error) {
	fields := make(map[string]bool)

	if err := extractFieldsRecursive(condition, fields); err != nil {
		return nil, err
	}

	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}

	return result, nil
}

// ActionFieldExtractor extrait les champs depuis les actions des nœuds terminaux.
type ActionFieldExtractor struct{}

// ExtractFields extrait les champs depuis une action.
//
// Exemples :
//   - Update(product, {price: product.price * 1.1}) → ["price"]
//   - Insert(Alert, {message: order.id}) → ["id"] (du contexte)
//
// Cette fonction examine les arguments et modifications pour extraire
// tous les champs référencés.
func (afe *ActionFieldExtractor) ExtractFields(action interface{}) ([]string, error) {
	fields := make(map[string]bool)

	if err := extractFieldsRecursive(action, fields); err != nil {
		return nil, err
	}

	result := make([]string, 0, len(fields))
	for field := range fields {
		result = append(result, field)
	}

	return result, nil
}

// extractFieldsRecursive est une fonction récursive privée qui parcourt
// une structure AST et collecte tous les champs référencés.
//
// Cette fonction est générique et fonctionne pour alpha, beta et terminal nodes.
func extractFieldsRecursive(node interface{}, fields map[string]bool) error {
	if node == nil {
		return nil
	}

	// Cas 1: Map (structure du parser)
	if nodeMap, ok := node.(map[string]interface{}); ok {
		return extractFieldsFromMap(nodeMap, fields)
	}

	// Cas 2: Slice (liste d'expressions)
	if nodeSlice, ok := node.([]interface{}); ok {
		return extractFieldsFromSlice(nodeSlice, fields)
	}

	// Cas 3: Types primitifs → ignorer (pas de champs)
	return nil
}

// extractFieldsFromMap extrait les champs depuis une map AST
func extractFieldsFromMap(nodeMap map[string]interface{}, fields map[string]bool) error {
	nodeType, hasType := nodeMap[ASTFieldNameType].(string)

	// Traiter les nœuds typés spécialement
	if hasType {
		if err := extractFieldsFromTypedNode(nodeType, nodeMap, fields); err != nil {
			return err
		}
	}

	// Récurser sur toutes les valeurs de la map
	for _, value := range nodeMap {
		if err := extractFieldsRecursive(value, fields); err != nil {
			return err
		}
	}

	return nil
}

// extractFieldsFromTypedNode traite les nœuds AST avec un type spécifique
func extractFieldsFromTypedNode(
	nodeType string,
	nodeMap map[string]interface{},
	fields map[string]bool,
) error {
	switch nodeType {
	case ASTNodeTypeFieldAccess:
		return extractFieldFromFieldAccess(nodeMap, fields)

	case ASTNodeTypeBinaryOp, ASTNodeTypeComparison:
		return extractFieldsFromBinaryNode(nodeMap, fields)

	case ASTNodeTypeUpdateWithModif:
		return extractFieldsFromUpdateNode(nodeMap, fields)

	case ASTNodeTypeFactCreation:
		return extractFieldsFromInsertNode(nodeMap, fields)
	}

	return nil
}

// extractFieldFromFieldAccess extrait le champ d'un nœud fieldAccess
func extractFieldFromFieldAccess(nodeMap map[string]interface{}, fields map[string]bool) error {
	if field, ok := nodeMap[ASTFieldNameField].(string); ok {
		fields[field] = true
	}
	return nil
}

// extractFieldsFromBinaryNode extrait les champs d'un nœud binaire (binaryOp, comparison)
func extractFieldsFromBinaryNode(nodeMap map[string]interface{}, fields map[string]bool) error {
	if err := extractFieldsRecursive(nodeMap[ASTFieldNameLeft], fields); err != nil {
		return err
	}
	if err := extractFieldsRecursive(nodeMap[ASTFieldNameRight], fields); err != nil {
		return err
	}
	return nil
}

// extractFieldsFromUpdateNode extrait les champs d'une action Update
func extractFieldsFromUpdateNode(nodeMap map[string]interface{}, fields map[string]bool) error {
	if modifications, ok := nodeMap[ASTFieldNameModifications].(map[string]interface{}); ok {
		for fieldName := range modifications {
			fields[fieldName] = true
		}
	}
	return nil
}

// extractFieldsFromInsertNode extrait les champs d'une action Insert
func extractFieldsFromInsertNode(nodeMap map[string]interface{}, fields map[string]bool) error {
	if factFields, ok := nodeMap[ASTFieldNameFields].(map[string]interface{}); ok {
		for fieldName := range factFields {
			fields[fieldName] = true
		}
	}
	return nil
}

// extractFieldsFromSlice extrait les champs depuis un slice AST
func extractFieldsFromSlice(nodeSlice []interface{}, fields map[string]bool) error {
	for _, item := range nodeSlice {
		if err := extractFieldsRecursive(item, fields); err != nil {
			return err
		}
	}
	return nil
}

// ExtractFieldsFromAlphaCondition est une fonction helper pour extraire
// les champs depuis une condition alpha.
func ExtractFieldsFromAlphaCondition(condition interface{}) ([]string, error) {
	extractor := &AlphaConditionExtractor{}
	return extractor.ExtractFields(condition)
}

// ExtractFieldsFromBetaCondition est une fonction helper pour extraire
// les champs depuis une condition beta.
func ExtractFieldsFromBetaCondition(condition interface{}) ([]string, error) {
	extractor := &BetaConditionExtractor{}
	return extractor.ExtractFields(condition)
}

// ExtractFieldsFromAction est une fonction helper pour extraire
// les champs depuis une action.
func ExtractFieldsFromAction(action interface{}) ([]string, error) {
	extractor := &ActionFieldExtractor{}
	return extractor.ExtractFields(action)
}
