// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"math"
	"sync"
)

type AccumulatorNode struct {
	BaseNode
	AggregateFunc string                 `json:"aggregate_func"` // "AVG", "SUM", "COUNT", "MIN", "MAX"
	MainVariable  string                 `json:"main_variable"`  // Variable principale (ex: "e")
	MainType      string                 `json:"main_type"`      // Type principal (ex: "Employee")
	AggVariable   string                 `json:"agg_variable"`   // Variable à agréger (ex: "p")
	AggType       string                 `json:"agg_type"`       // Type à agréger (ex: "Performance")
	Field         string                 `json:"field"`          // Champ à agréger (ex: "score"), vide pour COUNT
	JoinField     string                 `json:"join_field"`     // Champ de jointure (ex: "employee_id")
	MainField     string                 `json:"main_field"`     // Champ principal pour jointure (ex: "id")
	Condition     map[string]interface{} `json:"condition"`      // Condition de comparaison du résultat
	MainFacts     map[string]*Fact       `json:"-"`              // Faits principaux indexés par ID
	AllFacts      map[string]*Fact       `json:"-"`              // Tous les faits (principaux + agrégés) par ID
	mutex         sync.RWMutex
}

// NewAccumulatorNode crée un nouveau nœud d'agrégation
func NewAccumulatorNode(id string, mainVar, mainType, aggVar, aggType, field, joinField, mainField, aggregateFunc string, condition map[string]interface{}, storage Storage) *AccumulatorNode {
	return &AccumulatorNode{
		BaseNode: BaseNode{
			ID:       id,
			Type:     "accumulator",
			Children: make([]Node, 0),
			Memory:   &WorkingMemory{Tokens: make(map[string]*Token), Facts: make(map[string]*Fact)},
		},
		AggregateFunc: aggregateFunc,
		MainVariable:  mainVar,
		MainType:      mainType,
		AggVariable:   aggVar,
		AggType:       aggType,
		Field:         field,
		JoinField:     joinField,
		MainField:     mainField,
		Condition:     condition,
		MainFacts:     make(map[string]*Fact),
		AllFacts:      make(map[string]*Fact),
	}
}

// Activate traite un fait dans le nœud d'agrégation
func (an *AccumulatorNode) Activate(fact *Fact, token *Token) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Stocker tous les faits
	an.AllFacts[fact.ID] = fact

	// Si c'est un fait principal, stocker et calculer l'agrégation
	if fact.Type == an.MainType {
		an.MainFacts[fact.ID] = fact
		fmt.Printf("📊 ACCUMULATOR[%s]: Fait principal reçu %s\n", an.ID, fact.ID)

		// Calculer l'agrégation pour ce fait principal
		return an.processMainFact(fact)
	}

	// Si c'est un fait à agréger, recalculer pour tous les faits principaux
	if fact.Type == an.AggType {
		fmt.Printf("📊 ACCUMULATOR[%s]: Fait agrégé reçu %s\n", an.ID, fact.ID)
		// Recalculer pour tous les faits principaux existants
		for _, mainFact := range an.MainFacts {
			if err := an.processMainFact(mainFact); err != nil {
				fmt.Printf("⚠️  ACCUMULATOR[%s]: Erreur recalcul pour %s: %v\n", an.ID, mainFact.ID, err)
			}
		}
	}

	return nil
}

// processMainFact calcule l'agrégation pour un fait principal donné
func (an *AccumulatorNode) processMainFact(mainFact *Fact) error {
	// Collecter les faits à agréger qui correspondent à ce fait principal
	aggregatedFacts := an.collectAggregatedFacts(mainFact)

	fmt.Printf("📊 ACCUMULATOR[%s]: %d faits agrégés trouvés pour %s\n", an.ID, len(aggregatedFacts), mainFact.ID)

	// Calculer l'agrégation
	aggregatedValue, err := an.calculateAggregateForFacts(aggregatedFacts)
	if err != nil {
		return fmt.Errorf("erreur calcul agrégation: %w", err)
	}

	fmt.Printf("📊 ACCUMULATOR[%s]: Valeur agrégée = %.2f pour %s\n", an.ID, aggregatedValue, mainFact.ID)

	// Évaluer la condition
	satisfied, err := an.evaluateCondition(aggregatedValue)
	if err != nil {
		return fmt.Errorf("erreur évaluation condition agrégation: %w", err)
	}

	if satisfied {
		fmt.Printf("✅ ACCUMULATOR[%s]: Condition satisfaite (%.2f) pour %s\n", an.ID, aggregatedValue, mainFact.ID)

		// Créer un token avec le fait et le résultat de l'agrégation
		newToken := &Token{
			ID:       fmt.Sprintf("accum_%s", mainFact.ID),
			Facts:    []*Fact{mainFact},
			Bindings: NewBindingChainWith(an.MainVariable, mainFact),
		}
		an.Memory.AddToken(newToken)

		// Propager aux enfants - ne passer que le token, pas le fait
		// car TerminalNode ne veut que des tokens
		return an.PropagateToChildren(nil, newToken)
	} else {
		fmt.Printf("❌ ACCUMULATOR[%s]: Condition NON satisfaite (%.2f) pour %s\n", an.ID, aggregatedValue, mainFact.ID)
	}

	return nil
}

// collectAggregatedFacts collecte les faits à agréger pour un fait principal
func (an *AccumulatorNode) collectAggregatedFacts(mainFact *Fact) []*Fact {
	collected := make([]*Fact, 0)

	// Obtenir la valeur du champ de jointure du fait principal
	mainValue, exists := mainFact.Fields[an.MainField]
	if !exists {
		// Essayer aussi dans fact.ID si c'est le champ "id"
		if an.MainField == "id" {
			mainValue = mainFact.ID
		} else {
			fmt.Printf("⚠️  ACCUMULATOR[%s]: Champ principal %s non trouvé dans %s\n", an.ID, an.MainField, mainFact.ID)
			return collected
		}
	}

	// Parcourir tous les faits pour trouver ceux qui correspondent
	for _, fact := range an.AllFacts {
		if fact.Type == an.AggType {
			// Vérifier la condition de jointure
			joinValue, exists := fact.Fields[an.JoinField]
			if exists && joinValue == mainValue {
				collected = append(collected, fact)
			}
		}
	}

	return collected
}

// calculateAggregateForFacts calcule la valeur agrégée pour une liste de faits
func (an *AccumulatorNode) calculateAggregateForFacts(facts []*Fact) (float64, error) {
	if len(facts) == 0 {
		// Pas de faits à agréger - retourner 0
		return 0, nil
	}

	switch an.AggregateFunc {
	case "COUNT":
		return float64(len(facts)), nil

	case "SUM":
		sum := 0.0
		for _, f := range facts {
			if val, ok := f.Fields[an.Field]; ok {
				numVal := an.toFloat64(val)
				sum += numVal
			}
		}
		return sum, nil

	case "AVG":
		sum := 0.0
		count := 0
		for _, f := range facts {
			if val, ok := f.Fields[an.Field]; ok {
				numVal := an.toFloat64(val)
				if numVal != 0 || val == 0 || val == 0.0 {
					sum += numVal
					count++
				}
			}
		}
		if count == 0 {
			return 0, nil
		}
		return sum / float64(count), nil

	case "MIN":
		minVal := math.MaxFloat64
		for _, f := range facts {
			if val, ok := f.Fields[an.Field]; ok {
				numVal := an.toFloat64(val)
				if numVal < minVal {
					minVal = numVal
				}
			}
		}
		if minVal == math.MaxFloat64 {
			return 0, nil
		}
		return minVal, nil

	case "MAX":
		maxVal := -math.MaxFloat64
		for _, f := range facts {
			if val, ok := f.Fields[an.Field]; ok {
				numVal := an.toFloat64(val)
				if numVal > maxVal {
					maxVal = numVal
				}
			}
		}
		if maxVal == -math.MaxFloat64 {
			return 0, nil
		}
		return maxVal, nil

	default:
		return 0, fmt.Errorf("fonction d'agrégation non supportée: %s", an.AggregateFunc)
	}
}

// toFloat64 converts various numeric types to float64
func (an *AccumulatorNode) toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	default:
		return 0
	}
}

// ActivateLeft traite un token venant de la gauche (compatible avec interface Node)
func (an *AccumulatorNode) ActivateLeft(token *Token) error {
	// Pour AccumulatorNode, on traite le premier fait du token
	if len(token.Facts) > 0 {
		return an.Activate(token.Facts[0], token)
	}
	return nil
}

// ActivateRight traite un fait venant de la droite
func (an *AccumulatorNode) ActivateRight(fact *Fact) error {
	return an.Activate(fact, nil)
}

// evaluateCondition évalue si la valeur agrégée satisfait la condition
func (an *AccumulatorNode) evaluateCondition(aggregatedValue float64) (bool, error) {
	if an.Condition == nil {
		return true, nil
	}

	condType, ok := an.Condition["type"].(string)
	if !ok || condType != "comparison" {
		return false, fmt.Errorf("type de condition invalide")
	}

	operator, ok := an.Condition["operator"].(string)
	if !ok {
		return false, fmt.Errorf("opérateur manquant")
	}

	threshold, ok := an.Condition["value"].(float64)
	if !ok {
		return false, fmt.Errorf("valeur de comparaison invalide")
	}

	switch operator {
	case ">=":
		return aggregatedValue >= threshold, nil
	case ">":
		return aggregatedValue > threshold, nil
	case "<=":
		return aggregatedValue <= threshold, nil
	case "<":
		return aggregatedValue < threshold, nil
	case "==":
		return aggregatedValue == threshold, nil
	case "!=":
		return aggregatedValue != threshold, nil
	default:
		return false, fmt.Errorf("opérateur non supporté: %s", operator)
	}
}

// ActivateRetract gère la rétractation dans le nœud d'agrégation
func (an *AccumulatorNode) ActivateRetract(factID string) error {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	// Retirer des faits principaux et de tous les faits
	delete(an.MainFacts, factID)
	delete(an.AllFacts, factID)

	// Retirer des tokens
	an.Memory.RemoveToken(factID)

	fmt.Printf("🗑️  [ACCUMULATOR_%s] Rétractation: fait %s retiré\n", an.ID, factID)
	return an.PropagateRetractToChildren(factID)
}
