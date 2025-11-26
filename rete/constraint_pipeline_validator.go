package rete

import (
	"fmt"
)

// validateNetwork valide qu'un réseau RETE est bien formé
func (cp *ConstraintPipeline) validateNetwork(network *ReteNetwork) error {
	if network == nil {
		return fmt.Errorf("réseau est nil")
	}

	// Vérifier qu'on a au moins un TypeNode
	if len(network.TypeNodes) == 0 {
		return fmt.Errorf("aucun TypeNode dans le réseau")
	}

	// Vérifier qu'on a au moins un nœud terminal
	if len(network.TerminalNodes) == 0 {
		return fmt.Errorf("aucun nœud terminal dans le réseau")
	}

	// Vérifier que les nœuds terminaux ont des actions
	for id, terminal := range network.TerminalNodes {
		if terminal.Action == nil {
			return fmt.Errorf("nœud terminal %s sans action", id)
		}
	}

	return nil
}

// validateAction valide qu'une action est bien formée
func (cp *ConstraintPipeline) validateAction(actionMap map[string]interface{}) error {
	if actionMap == nil {
		return fmt.Errorf("action map est nil")
	}

	// Vérifier qu'on a un type d'action
	actionType, hasType := actionMap["type"].(string)
	if !hasType {
		return fmt.Errorf("type d'action non trouvé")
	}

	// Vérifier selon le type d'action
	switch actionType {
	case "print", "PRINT":
		// L'action print doit avoir un message ou une expression
		if _, hasMsg := actionMap["message"]; !hasMsg {
			if _, hasExpr := actionMap["expression"]; !hasExpr {
				return fmt.Errorf("action print sans message ni expression")
			}
		}
	case "assert", "ASSERT":
		// L'action assert doit avoir un fait à insérer
		if _, hasFact := actionMap["fact"]; !hasFact {
			return fmt.Errorf("action assert sans fait")
		}
	case "retract", "RETRACT":
		// L'action retract doit avoir un fait à retirer
		if _, hasFact := actionMap["fact"]; !hasFact {
			return fmt.Errorf("action retract sans fait")
		}
	}

	return nil
}

// validateRuleExpression valide qu'une expression de règle est bien formée
func (cp *ConstraintPipeline) validateRuleExpression(exprMap map[string]interface{}) error {
	if exprMap == nil {
		return fmt.Errorf("expression map est nil")
	}

	// Vérifier qu'on a une action
	actionData, hasAction := exprMap["action"]
	if !hasAction {
		return fmt.Errorf("aucune action trouvée dans la règle")
	}

	actionMap, ok := actionData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("format action invalide: %T", actionData)
	}

	// Valider l'action
	if err := cp.validateAction(actionMap); err != nil {
		return fmt.Errorf("action invalide: %w", err)
	}

	// Vérifier qu'on a des variables (set)
	if _, hasSet := exprMap["set"]; !hasSet {
		return fmt.Errorf("aucune variable (set) trouvée dans la règle")
	}

	return nil
}

// validateTypeDefinition valide qu'une définition de type est bien formée
func (cp *ConstraintPipeline) validateTypeDefinition(typeName string, typeMap map[string]interface{}) error {
	if typeName == "" {
		return fmt.Errorf("nom de type vide")
	}

	if typeMap == nil {
		return fmt.Errorf("définition de type nil pour %s", typeName)
	}

	// Vérifier qu'on a des champs
	fieldsData, hasFields := typeMap["fields"]
	if !hasFields {
		return fmt.Errorf("type %s sans champs", typeName)
	}

	fields, ok := fieldsData.([]interface{})
	if !ok {
		return fmt.Errorf("format fields invalide pour type %s: %T", typeName, fieldsData)
	}

	if len(fields) == 0 {
		return fmt.Errorf("type %s avec liste de champs vide", typeName)
	}

	// Vérifier chaque champ
	for i, fieldInterface := range fields {
		fieldMap, ok := fieldInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("champ %d du type %s invalide: %T", i, typeName, fieldInterface)
		}

		// Vérifier qu'on a un nom de champ
		fieldName, hasName := fieldMap["name"].(string)
		if !hasName || fieldName == "" {
			return fmt.Errorf("champ %d du type %s sans nom", i, typeName)
		}

		// Vérifier qu'on a un type de champ
		fieldType, hasType := fieldMap["type"].(string)
		if !hasType || fieldType == "" {
			return fmt.Errorf("champ %s du type %s sans type", fieldName, typeName)
		}
	}

	return nil
}

// validateAggregationInfo valide qu'une information d'agrégation est complète
func (cp *ConstraintPipeline) validateAggregationInfo(aggInfo *AggregationInfo) error {
	if aggInfo == nil {
		return fmt.Errorf("information d'agrégation nil")
	}

	// Vérifier la fonction
	if aggInfo.Function == "" {
		return fmt.Errorf("fonction d'agrégation vide")
	}

	validFunctions := map[string]bool{
		"AVG": true, "SUM": true, "COUNT": true,
		"MIN": true, "MAX": true, "ACCUMULATE": true,
	}
	if !validFunctions[aggInfo.Function] {
		return fmt.Errorf("fonction d'agrégation invalide: %s", aggInfo.Function)
	}

	// Vérifier l'opérateur
	if aggInfo.Operator == "" {
		return fmt.Errorf("opérateur de comparaison vide")
	}

	validOperators := map[string]bool{
		">=": true, "<=": true, ">": true, "<": true, "==": true, "!=": true,
	}
	if !validOperators[aggInfo.Operator] {
		return fmt.Errorf("opérateur de comparaison invalide: %s", aggInfo.Operator)
	}

	// Vérifier les champs de jointure si présents
	if aggInfo.JoinField != "" && aggInfo.MainField == "" {
		return fmt.Errorf("champ de jointure principal manquant")
	}
	if aggInfo.MainField != "" && aggInfo.JoinField == "" {
		return fmt.Errorf("champ de jointure agrégé manquant")
	}

	return nil
}

// validateJoinCondition valide qu'une condition de jointure est bien formée
func (cp *ConstraintPipeline) validateJoinCondition(condition map[string]interface{}) error {
	if condition == nil {
		return fmt.Errorf("condition de jointure nil")
	}

	condType, hasType := condition["type"].(string)
	if !hasType {
		return fmt.Errorf("type de condition non spécifié")
	}

	switch condType {
	case "simple", "passthrough":
		// Conditions simples, pas de validation supplémentaire
		return nil
	case "constraint":
		// Doit avoir une contrainte
		if _, hasConstraint := condition["constraint"]; !hasConstraint {
			return fmt.Errorf("condition de type constraint sans contrainte")
		}
	case "negation":
		// Doit avoir une condition niée
		if _, hasNegated := condition["negated"]; !hasNegated {
			return fmt.Errorf("condition de type negation sans flag negated")
		}
		if _, hasCondition := condition["condition"]; !hasCondition {
			return fmt.Errorf("condition de type negation sans condition niée")
		}
	default:
		return fmt.Errorf("type de condition inconnu: %s", condType)
	}

	return nil
}
