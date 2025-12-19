// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import "fmt"

// ValidateXupleSpaceDeclaration valide une déclaration de xuple-space.
// Vérifie que toutes les propriétés sont valides et cohérentes.
func ValidateXupleSpaceDeclaration(decl *XupleSpaceDeclaration) error {
	if decl == nil {
		return fmt.Errorf("xuple-space declaration cannot be nil")
	}

	if decl.Name == "" {
		return fmt.Errorf("xuple-space name cannot be empty")
	}

	// Valider la politique de sélection
	if err := validateSelectionPolicy(decl.Name, decl.SelectionPolicy); err != nil {
		return err
	}

	// Valider la politique de consommation
	if err := validateConsumptionPolicy(decl.Name, decl.ConsumptionPolicy); err != nil {
		return err
	}

	// Valider la politique de rétention
	if err := validateRetentionPolicy(decl.Name, decl.RetentionPolicy); err != nil {
		return err
	}

	return nil
}

// validateSelectionPolicy valide la politique de sélection
func validateSelectionPolicy(name, policy string) error {
	switch policy {
	case SelectionFIFO, SelectionLIFO, SelectionRandom:
		return nil
	case "":
		return fmt.Errorf("xuple-space '%s': selection policy cannot be empty", name)
	default:
		return fmt.Errorf("xuple-space '%s': invalid selection policy '%s' (must be %s, %s, or %s)",
			name, policy, SelectionFIFO, SelectionLIFO, SelectionRandom)
	}
}

// validateConsumptionPolicy valide la politique de consommation
func validateConsumptionPolicy(name string, policy XupleConsumptionPolicyConf) error {
	switch policy.Type {
	case ConsumptionOnce, ConsumptionPerAgent:
		// Ces types ne nécessitent pas de limite
		return nil
	case ConsumptionLimited:
		if policy.Limit <= 0 {
			return fmt.Errorf("xuple-space '%s': consumption limit must be greater than 0 for 'limited' policy, got %d",
				name, policy.Limit)
		}
		return nil
	case "":
		return fmt.Errorf("xuple-space '%s': consumption policy type cannot be empty", name)
	default:
		return fmt.Errorf("xuple-space '%s': invalid consumption policy type '%s' (must be %s, %s, or %s)",
			name, policy.Type, ConsumptionOnce, ConsumptionPerAgent, ConsumptionLimited)
	}
}

// validateRetentionPolicy valide la politique de rétention
func validateRetentionPolicy(name string, policy XupleRetentionPolicyConf) error {
	switch policy.Type {
	case RetentionUnlimited:
		// Pas de durée nécessaire
		return nil
	case RetentionDuration:
		if policy.Duration <= 0 {
			return fmt.Errorf("xuple-space '%s': retention duration must be greater than 0 for 'duration' policy, got %d",
				name, policy.Duration)
		}
		return nil
	case "":
		return fmt.Errorf("xuple-space '%s': retention policy type cannot be empty", name)
	default:
		return fmt.Errorf("xuple-space '%s': invalid retention policy type '%s' (must be %s or %s)",
			name, policy.Type, RetentionUnlimited, RetentionDuration)
	}
}

// ValidateXupleSpaceProperties valide les propriétés d'un xuple-space depuis une map.
// Cette fonction est utilisée lors du parsing pour valider les propriétés parsées.
func ValidateXupleSpaceProperties(name string, props map[string]interface{}) error {
	if name == "" {
		return fmt.Errorf("xuple-space name cannot be empty")
	}

	// Valider selection si présente
	if sel, ok := props["selectionPolicy"]; ok {
		selStr, ok := sel.(string)
		if !ok {
			return fmt.Errorf("xuple-space '%s': selectionPolicy must be a string, got %T", name, sel)
		}
		if err := validateSelectionPolicy(name, selStr); err != nil {
			return err
		}
	}

	// Valider consumption si présente
	if cons, ok := props["consumptionPolicy"]; ok {
		consMap, ok := cons.(map[string]interface{})
		if !ok {
			return fmt.Errorf("xuple-space '%s': consumptionPolicy must be a map, got %T", name, cons)
		}

		consType, _ := consMap["type"].(string)
		consLimit := 0
		if limit, ok := consMap["limit"]; ok {
			switch v := limit.(type) {
			case int:
				consLimit = v
			case float64:
				consLimit = int(v)
			default:
				return fmt.Errorf("xuple-space '%s': consumption limit must be an integer, got %T", name, limit)
			}
		}

		policy := XupleConsumptionPolicyConf{
			Type:  consType,
			Limit: consLimit,
		}
		if err := validateConsumptionPolicy(name, policy); err != nil {
			return err
		}
	}

	// Valider retention si présente
	if ret, ok := props["retentionPolicy"]; ok {
		retMap, ok := ret.(map[string]interface{})
		if !ok {
			return fmt.Errorf("xuple-space '%s': retentionPolicy must be a map, got %T", name, ret)
		}

		retType, _ := retMap["type"].(string)
		retDuration := 0
		if dur, ok := retMap["duration"]; ok {
			switch v := dur.(type) {
			case int:
				retDuration = v
			case float64:
				retDuration = int(v)
			default:
				return fmt.Errorf("xuple-space '%s': retention duration must be an integer (seconds), got %T", name, dur)
			}
		}

		policy := XupleRetentionPolicyConf{
			Type:     retType,
			Duration: retDuration,
		}
		if err := validateRetentionPolicy(name, policy); err != nil {
			return err
		}
	}

	// Valider max-size si présente
	if maxSize, ok := props["maxSize"]; ok {
		switch v := maxSize.(type) {
		case int:
			if v < 0 {
				return fmt.Errorf("xuple-space '%s': max-size must be >= 0, got %d", name, v)
			}
		case float64:
			if v < 0 {
				return fmt.Errorf("xuple-space '%s': max-size must be >= 0, got %f", name, v)
			}
		default:
			return fmt.Errorf("xuple-space '%s': max-size must be an integer, got %T", name, maxSize)
		}
	}

	return nil
}
