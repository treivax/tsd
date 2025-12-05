// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// isExistsConstraint vérifie si une contrainte est de type EXISTS
func (cp *ConstraintPipeline) isExistsConstraint(constraintsData interface{}) bool {
	if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
		if constraintType, exists := constraintMap["type"].(string); exists && constraintType == "existsConstraint" {
			return true
		}
	}
	return false
}

// getStringField extrait un champ string d'une map avec une valeur par défaut
func getStringField(m map[string]interface{}, key, defaultValue string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return defaultValue
}
