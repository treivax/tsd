package nodes

import (
	"fmt"
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// NotNodeImpl implémente l'interface NotNode pour la négation
type NotNodeImpl struct {
	*BaseBetaNode
	negationCondition interface{}
	mu                sync.RWMutex
}

// NewNotNode crée un nouveau nœud de négation
func NewNotNode(id string, logger domain.Logger) *NotNodeImpl {
	baseBeta := NewBaseBetaNode(id, "NotNode", logger)
	return &NotNodeImpl{
		BaseBetaNode: baseBeta,
	}
}

// SetNegationCondition définit la condition de négation
func (n *NotNodeImpl) SetNegationCondition(condition interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.negationCondition = condition
}

// GetNegationCondition retourne la condition de négation
func (n *NotNodeImpl) GetNegationCondition() interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.negationCondition
}

// ProcessNegation évalue la négation d'une condition
func (n *NotNodeImpl) ProcessNegation(token *domain.Token, fact *domain.Fact) bool {
	n.mu.RLock()
	condition := n.negationCondition
	n.mu.RUnlock()

	if condition == nil {
		return false
	}

	// Évaluer la condition et retourner sa négation
	result, err := n.evaluateCondition(condition, token, fact)
	if err != nil {
		n.logger.Error("Erreur évaluation condition négation", err, map[string]interface{}{
			"node_id": n.id,
			"token":   token.ID,
			"fact":    fact.ID,
		})
		return false
	}

	return !result // Négation du résultat
}

// ProcessLeftToken traite un token venant de la gauche
func (n *NotNodeImpl) ProcessLeftToken(token *domain.Token) error {
	n.logger.Debug("processing token in NotNode", map[string]interface{}{
		"node_id":    n.id,
		"token_id":   token.ID,
		"node_type":  "NotNode",
		"action":     "left_input",
		"fact_count": len(token.Facts),
	})

	// Stocker le token dans la mémoire gauche
	n.betaMemory.StoreToken(token)

	// Vérifier la négation contre tous les faits de droite
	rightFacts := n.betaMemory.GetFacts()
	shouldPropagate := true

	for _, fact := range rightFacts {
		if n.ProcessNegation(token, fact) {
			// Si la négation est vraie (condition originale fausse), continuer
			continue
		} else {
			// Si la négation est fausse (condition originale vraie), bloquer la propagation
			shouldPropagate = false
			break
		}
	}

	// Si aucun fait de droite ne satisfait la condition, propager le token (négation réussie)
	if shouldPropagate && len(rightFacts) > 0 {
		return n.propagateTokenToChildren(token)
	}

	// Si pas de faits de droite, propager également (négation par défaut)
	if len(rightFacts) == 0 {
		return n.propagateTokenToChildren(token)
	}

	return nil
}

// ProcessRightFact traite un fait venant de la droite
func (n *NotNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	n.logger.Debug("processing fact in NotNode", map[string]interface{}{
		"node_id":   n.id,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"node_type": "NotNode",
		"action":    "right_input",
	})

	// Stocker le fait dans la mémoire droite
	n.betaMemory.StoreFact(fact)

	// Vérifier tous les tokens de gauche
	leftTokens := n.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if !n.ProcessNegation(token, fact) {
			// Si la négation échoue (condition vraie), retirer le token s'il était propagé
			n.logger.Debug("negation failed, blocking token", map[string]interface{}{
				"node_id":  n.id,
				"token_id": token.ID,
				"fact_id":  fact.ID,
			})
		}
	}

	return nil
}

// evaluateCondition évalue une condition (méthode helper)
func (n *NotNodeImpl) evaluateCondition(condition interface{}, token *domain.Token, fact *domain.Fact) (bool, error) {
	// Pour l'instant, implémentation basique
	// TODO: Intégrer avec l'évaluateur d'expressions du package constraint
	return true, nil
}

// ExistsNodeImpl implémente l'interface ExistsNode pour la quantification existentielle
type ExistsNodeImpl struct {
	*BaseBetaNode
	existenceVariable  domain.TypedVariable
	existenceCondition interface{}
	mu                 sync.RWMutex
}

// NewExistsNode crée un nouveau nœud EXISTS
func NewExistsNode(id string, logger domain.Logger) *ExistsNodeImpl {
	baseBeta := NewBaseBetaNode(id, "ExistsNode", logger)
	return &ExistsNodeImpl{
		BaseBetaNode: baseBeta,
	}
}

// SetExistenceCondition définit la condition d'existence
func (e *ExistsNodeImpl) SetExistenceCondition(variable domain.TypedVariable, condition interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.existenceVariable = variable
	e.existenceCondition = condition
}

// GetExistenceCondition retourne la condition d'existence
func (e *ExistsNodeImpl) GetExistenceCondition() (domain.TypedVariable, interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.existenceVariable, e.existenceCondition
}

// CheckExistence vérifie l'existence d'au moins un fait satisfaisant la condition
func (e *ExistsNodeImpl) CheckExistence(token *domain.Token) bool {
	e.mu.RLock()
	condition := e.existenceCondition
	e.mu.RUnlock()

	if condition == nil {
		return false
	}

	// Vérifier tous les faits de droite
	rightFacts := e.betaMemory.GetFacts()
	for _, fact := range rightFacts {
		// Si le fait correspond au type de la variable d'existence
		if fact.Type == e.existenceVariable.DataType {
			// Évaluer la condition
			if result, err := e.evaluateExistenceCondition(condition, token, fact); err == nil && result {
				return true // Au moins un fait satisfait la condition
			}
		}
	}

	return false // Aucun fait ne satisfait la condition
}

// ProcessLeftToken traite un token venant de la gauche
func (e *ExistsNodeImpl) ProcessLeftToken(token *domain.Token) error {
	e.logger.Debug("processing token in ExistsNode", map[string]interface{}{
		"node_id":    e.id,
		"token_id":   token.ID,
		"node_type":  "ExistsNode",
		"action":     "left_input",
		"fact_count": len(token.Facts),
	})

	// Stocker le token dans la mémoire gauche
	e.betaMemory.StoreToken(token)

	// Vérifier l'existence
	if e.CheckExistence(token) {
		e.logger.Debug("existence condition satisfied", map[string]interface{}{
			"node_id":  e.id,
			"token_id": token.ID,
		})
		return e.propagateTokenToChildren(token)
	}

	e.logger.Debug("existence condition not satisfied", map[string]interface{}{
		"node_id":  e.id,
		"token_id": token.ID,
	})

	return nil
}

// ProcessRightFact traite un fait venant de la droite
func (e *ExistsNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	e.logger.Debug("processing fact in ExistsNode", map[string]interface{}{
		"node_id":   e.id,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"node_type": "ExistsNode",
		"action":    "right_input",
	})

	// Stocker le fait dans la mémoire droite
	e.betaMemory.StoreFact(fact)

	// Vérifier tous les tokens de gauche pour voir si l'existence est maintenant satisfaite
	leftTokens := e.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if e.CheckExistence(token) {
			e.logger.Debug("existence now satisfied by new fact", map[string]interface{}{
				"node_id":  e.id,
				"token_id": token.ID,
				"fact_id":  fact.ID,
			})
			// Propager le token s'il n'était pas déjà propagé
			e.propagateTokenToChildren(token)
		}
	}

	return nil
}

// evaluateExistenceCondition évalue une condition d'existence
func (e *ExistsNodeImpl) evaluateExistenceCondition(condition interface{}, token *domain.Token, fact *domain.Fact) (bool, error) {
	// Pour l'instant, implémentation basique
	// TODO: Intégrer avec l'évaluateur d'expressions du package constraint
	return true, nil
}

// AccumulateNodeImpl implémente l'interface AccumulateNode pour les fonctions d'agrégation
type AccumulateNodeImpl struct {
	*BaseBetaNode
	accumulator       domain.AccumulateFunction
	accumulatedValues map[string]interface{} // Stockage des valeurs agrégées par token
	mu                sync.RWMutex
}

// NewAccumulateNode crée un nouveau nœud d'accumulation
func NewAccumulateNode(id string, accumulator domain.AccumulateFunction, logger domain.Logger) *AccumulateNodeImpl {
	baseBeta := NewBaseBetaNode(id, "AccumulateNode", logger)
	return &AccumulateNodeImpl{
		BaseBetaNode:      baseBeta,
		accumulator:       accumulator,
		accumulatedValues: make(map[string]interface{}),
	}
}

// SetAccumulator définit la fonction d'accumulation
func (a *AccumulateNodeImpl) SetAccumulator(accumulator domain.AccumulateFunction) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.accumulator = accumulator
}

// GetAccumulator retourne la fonction d'accumulation
func (a *AccumulateNodeImpl) GetAccumulator() domain.AccumulateFunction {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.accumulator
}

// ComputeAggregate calcule l'agrégation pour un token donné
func (a *AccumulateNodeImpl) ComputeAggregate(token *domain.Token, facts []*domain.Fact) (interface{}, error) {
	a.mu.RLock()
	accumulator := a.accumulator
	a.mu.RUnlock()

	if accumulator.FunctionType == "" {
		return nil, fmt.Errorf("no accumulator function defined")
	}

	switch accumulator.FunctionType {
	case "SUM":
		return a.computeSum(facts, accumulator.Field)
	case "COUNT":
		return len(facts), nil
	case "AVG":
		return a.computeAverage(facts, accumulator.Field)
	case "MIN":
		return a.computeMin(facts, accumulator.Field)
	case "MAX":
		return a.computeMax(facts, accumulator.Field)
	default:
		return nil, fmt.Errorf("unsupported accumulator function: %s", accumulator.FunctionType)
	}
}

// ProcessLeftToken traite un token venant de la gauche
func (a *AccumulateNodeImpl) ProcessLeftToken(token *domain.Token) error {
	a.logger.Debug("processing token in AccumulateNode", map[string]interface{}{
		"node_id":    a.id,
		"token_id":   token.ID,
		"node_type":  "AccumulateNode",
		"action":     "left_input",
		"fact_count": len(token.Facts),
	})

	// Stocker le token dans la mémoire gauche
	a.betaMemory.StoreToken(token)

	// Obtenir tous les faits de droite pour l'agrégation
	rightFacts := a.betaMemory.GetFacts()

	// Calculer l'agrégation
	result, err := a.ComputeAggregate(token, rightFacts)
	if err != nil {
		a.logger.Error("failed to compute aggregate", err, map[string]interface{}{
			"node_id":  a.id,
			"token_id": token.ID,
		})
		return err
	}

	// Stocker le résultat
	a.mu.Lock()
	a.accumulatedValues[token.ID] = result
	a.mu.Unlock()

	// Créer un nouveau token avec le résultat d'agrégation
	newToken := &domain.Token{
		ID: fmt.Sprintf("%s_agg", token.ID),
		Facts: append(token.Facts, &domain.Fact{
			ID:   fmt.Sprintf("agg_%s", token.ID),
			Type: "AggregateResult",
			Fields: map[string]interface{}{
				"function": a.accumulator.FunctionType,
				"value":    result,
			},
		}),
	}

	// Propager le token enrichi
	return a.propagateTokenToChildren(newToken)
}

// ProcessRightFact traite un fait venant de la droite
func (a *AccumulateNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	a.logger.Debug("processing fact in AccumulateNode", map[string]interface{}{
		"node_id":   a.id,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"node_type": "AccumulateNode",
		"action":    "right_input",
	})

	// Stocker le fait dans la mémoire droite
	a.betaMemory.StoreFact(fact)

	// Recalculer l'agrégation pour tous les tokens de gauche
	leftTokens := a.betaMemory.GetTokens()
	for _, token := range leftTokens {
		// Obtenir tous les faits de droite
		rightFacts := a.betaMemory.GetFacts()

		// Recalculer l'agrégation
		result, err := a.ComputeAggregate(token, rightFacts)
		if err != nil {
			a.logger.Error("failed to recompute aggregate", err, map[string]interface{}{
				"node_id":  a.id,
				"token_id": token.ID,
				"fact_id":  fact.ID,
			})
			continue
		}

		// Mettre à jour le résultat
		a.mu.Lock()
		oldResult, existed := a.accumulatedValues[token.ID]
		a.accumulatedValues[token.ID] = result
		a.mu.Unlock()

		// Si le résultat a changé, propager la mise à jour
		if !existed || oldResult != result {
			a.logger.Debug("aggregate result updated", map[string]interface{}{
				"node_id":    a.id,
				"token_id":   token.ID,
				"old_result": oldResult,
				"new_result": result,
			})

			// Créer un token mis à jour
			newToken := &domain.Token{
				ID: fmt.Sprintf("%s_agg", token.ID),
				Facts: append(token.Facts, &domain.Fact{
					ID:   fmt.Sprintf("agg_%s", token.ID),
					Type: "AggregateResult",
					Fields: map[string]interface{}{
						"function": a.accumulator.FunctionType,
						"value":    result,
					},
				}),
			}

			// Propager le token mis à jour
			a.propagateTokenToChildren(newToken)
		}
	}

	return nil
}

// Fonctions d'aide pour les différentes agrégations

func (a *AccumulateNodeImpl) computeSum(facts []*domain.Fact, field string) (float64, error) {
	var sum float64
	count := 0

	for _, fact := range facts {
		if value, exists := fact.Fields[field]; exists {
			switch v := value.(type) {
			case int:
				sum += float64(v)
				count++
			case int64:
				sum += float64(v)
				count++
			case float64:
				sum += v
				count++
			case float32:
				sum += float64(v)
				count++
			default:
				// Ignorer les valeurs non numériques
				continue
			}
		}
	}

	return sum, nil
}

func (a *AccumulateNodeImpl) computeAverage(facts []*domain.Fact, field string) (float64, error) {
	sum, err := a.computeSum(facts, field)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, fact := range facts {
		if _, exists := fact.Fields[field]; exists {
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return sum / float64(count), nil
}

func (a *AccumulateNodeImpl) computeMin(facts []*domain.Fact, field string) (interface{}, error) {
	var minFloat float64
	var minString string
	var minOther interface{}
	foundNumeric := false
	foundString := false
	foundOther := false

	for _, fact := range facts {
		if value, exists := fact.Fields[field]; exists {
			switch v := value.(type) {
			case int:
				floatVal := float64(v)
				if !foundNumeric || floatVal < minFloat {
					minFloat = floatVal
					foundNumeric = true
				}
			case int64:
				floatVal := float64(v)
				if !foundNumeric || floatVal < minFloat {
					minFloat = floatVal
					foundNumeric = true
				}
			case float32:
				floatVal := float64(v)
				if !foundNumeric || floatVal < minFloat {
					minFloat = floatVal
					foundNumeric = true
				}
			case float64:
				if !foundNumeric || v < minFloat {
					minFloat = v
					foundNumeric = true
				}
			case string:
				if !foundString || v < minString {
					minString = v
					foundString = true
				}
			default:
				if !foundOther {
					minOther = v
					foundOther = true
				}
			}
		}
	}

	// Retourner le type le plus approprié
	if foundNumeric {
		return minFloat, nil
	}
	if foundString {
		return minString, nil
	}
	if foundOther {
		return minOther, nil
	}

	return nil, fmt.Errorf("no values found for field %s", field)
}

func (a *AccumulateNodeImpl) computeMax(facts []*domain.Fact, field string) (interface{}, error) {
	var maxFloat float64
	var maxString string
	var maxOther interface{}
	foundNumeric := false
	foundString := false
	foundOther := false

	for _, fact := range facts {
		if value, exists := fact.Fields[field]; exists {
			switch v := value.(type) {
			case int:
				floatVal := float64(v)
				if !foundNumeric || floatVal > maxFloat {
					maxFloat = floatVal
					foundNumeric = true
				}
			case int64:
				floatVal := float64(v)
				if !foundNumeric || floatVal > maxFloat {
					maxFloat = floatVal
					foundNumeric = true
				}
			case float32:
				floatVal := float64(v)
				if !foundNumeric || floatVal > maxFloat {
					maxFloat = floatVal
					foundNumeric = true
				}
			case float64:
				if !foundNumeric || v > maxFloat {
					maxFloat = v
					foundNumeric = true
				}
			case string:
				if !foundString || v > maxString {
					maxString = v
					foundString = true
				}
			default:
				if !foundOther {
					maxOther = v
					foundOther = true
				}
			}
		}
	}

	// Retourner le type le plus approprié
	if foundNumeric {
		return maxFloat, nil
	}
	if foundString {
		return maxString, nil
	}
	if foundOther {
		return maxOther, nil
	}

	return nil, fmt.Errorf("no values found for field %s", field)
}
