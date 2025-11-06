package nodes

import (
	"fmt"
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// BetaMemoryImpl implémente l'interface BetaMemory pour la gestion de la mémoire beta.
type BetaMemoryImpl struct {
	tokens map[string]*domain.Token
	facts  map[string]*domain.Fact
	mutex  sync.RWMutex
}

// NewBetaMemory crée une nouvelle instance de mémoire beta.
func NewBetaMemory() *BetaMemoryImpl {
	return &BetaMemoryImpl{
		tokens: make(map[string]*domain.Token),
		facts:  make(map[string]*domain.Fact),
	}
}

// StoreToken stocke un token dans la mémoire beta.
func (bm *BetaMemoryImpl) StoreToken(token *domain.Token) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	bm.tokens[token.ID] = token
}

// RemoveToken supprime un token de la mémoire beta.
func (bm *BetaMemoryImpl) RemoveToken(tokenID string) bool {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if _, exists := bm.tokens[tokenID]; exists {
		delete(bm.tokens, tokenID)
		return true
	}
	return false
}

// GetTokens retourne tous les tokens stockés.
func (bm *BetaMemoryImpl) GetTokens() []*domain.Token {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	tokens := make([]*domain.Token, 0, len(bm.tokens))
	for _, token := range bm.tokens {
		tokens = append(tokens, token)
	}
	return tokens
}

// StoreFact stocke un fait dans la mémoire beta.
func (bm *BetaMemoryImpl) StoreFact(fact *domain.Fact) {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()
	bm.facts[fact.ID] = fact
}

// RemoveFact supprime un fait de la mémoire beta.
func (bm *BetaMemoryImpl) RemoveFact(factID string) bool {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if _, exists := bm.facts[factID]; exists {
		delete(bm.facts, factID)
		return true
	}
	return false
}

// GetFacts retourne tous les faits stockés.
func (bm *BetaMemoryImpl) GetFacts() []*domain.Fact {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	facts := make([]*domain.Fact, 0, len(bm.facts))
	for _, fact := range bm.facts {
		facts = append(facts, fact)
	}
	return facts
}

// Clear vide la mémoire beta.
func (bm *BetaMemoryImpl) Clear() {
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	bm.tokens = make(map[string]*domain.Token)
	bm.facts = make(map[string]*domain.Fact)
}

// Size retourne le nombre de tokens et de faits stockés.
func (bm *BetaMemoryImpl) Size() (tokens int, facts int) {
	bm.mutex.RLock()
	defer bm.mutex.RUnlock()

	return len(bm.tokens), len(bm.facts)
}

// BaseBetaNode implémente les fonctionnalités de base d'un nœud beta.
type BaseBetaNode struct {
	*BaseNode
	betaMemory domain.BetaMemory
}

// NewBaseBetaNode crée un nouveau nœud beta de base.
func NewBaseBetaNode(id, nodeType string, logger domain.Logger) *BaseBetaNode {
	return &BaseBetaNode{
		BaseNode:   NewBaseNode(id, nodeType, logger),
		betaMemory: NewBetaMemory(),
	}
}

// ProcessFact implémente l'interface Node (délègue à ProcessRightFact).
func (bbn *BaseBetaNode) ProcessFact(fact *domain.Fact) error {
	return bbn.ProcessRightFact(fact)
}

// ProcessLeftToken traite un token venant du côté gauche.
func (bbn *BaseBetaNode) ProcessLeftToken(token *domain.Token) error {
	bbn.logTokenProcessing(token, "left_input")
	bbn.betaMemory.StoreToken(token)

	// Essayer de faire des jointures avec les faits du côté droit
	rightFacts := bbn.betaMemory.GetFacts()
	for _, fact := range rightFacts {
		if err := bbn.tryJoin(token, fact); err != nil {
			return err
		}
	}

	return nil
}

// ProcessRightFact traite un fait venant du côté droit.
func (bbn *BaseBetaNode) ProcessRightFact(fact *domain.Fact) error {
	bbn.logFactProcessing(fact, "right_input")
	bbn.betaMemory.StoreFact(fact)

	// Essayer de faire des jointures avec les tokens du côté gauche
	leftTokens := bbn.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if err := bbn.tryJoin(token, fact); err != nil {
			return err
		}
	}

	return nil
}

// tryJoin essaie de faire une jointure entre un token et un fait.
// Cette méthode de base accepte toute combinaison (pas de condition).
func (bbn *BaseBetaNode) tryJoin(token *domain.Token, fact *domain.Fact) error {
	// Créer un nouveau token combinant le token existant et le nouveau fait
	newFacts := make([]*domain.Fact, len(token.Facts)+1)
	copy(newFacts, token.Facts)
	newFacts[len(token.Facts)] = fact

	newTokenID := fmt.Sprintf("%s_%s", token.ID, fact.ID)
	newToken := domain.NewToken(newTokenID, bbn.ID(), newFacts)
	newToken.Parent = token

	bbn.logJoin(token, fact, newToken)

	// Propager le nouveau token vers les enfants
	return bbn.propagateTokenToChildren(newToken)
}

// GetLeftMemory retourne les tokens de la mémoire gauche.
func (bbn *BaseBetaNode) GetLeftMemory() []*domain.Token {
	return bbn.betaMemory.GetTokens()
}

// GetRightMemory retourne les faits de la mémoire droite.
func (bbn *BaseBetaNode) GetRightMemory() []*domain.Fact {
	return bbn.betaMemory.GetFacts()
}

// ClearMemory vide la mémoire beta.
func (bbn *BaseBetaNode) ClearMemory() {
	bbn.betaMemory.Clear()
}

// propagateTokenToChildren propage un token vers tous les enfants.
func (bbn *BaseBetaNode) propagateTokenToChildren(token *domain.Token) error {
	children := bbn.GetChildren()

	for _, child := range children {
		// Si l'enfant est un BetaNode, propager comme token gauche
		if betaChild, ok := child.(domain.BetaNode); ok {
			if err := betaChild.ProcessLeftToken(token); err != nil {
				bbn.logger.Error("failed to propagate token to beta child", err, map[string]interface{}{
					"parent_node": bbn.ID(),
					"child_node":  child.ID(),
					"token_id":    token.ID,
				})
				return err
			}
		} else {
			// Pour les autres types de nœuds, traiter comme fait (dernier fait du token)
			if len(token.Facts) > 0 {
				lastFact := token.Facts[len(token.Facts)-1]
				if err := child.ProcessFact(lastFact); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// logTokenProcessing enregistre le traitement d'un token.
func (bbn *BaseBetaNode) logTokenProcessing(token *domain.Token, action string) {
	bbn.logger.Debug("processing token", map[string]interface{}{
		"node_id":    bbn.ID(),
		"node_type":  bbn.Type(),
		"token_id":   token.ID,
		"fact_count": len(token.Facts),
		"action":     action,
	})
}

// logJoin enregistre une jointure réussie.
func (bbn *BaseBetaNode) logJoin(leftToken *domain.Token, rightFact *domain.Fact, resultToken *domain.Token) {
	bbn.logger.Debug("join performed", map[string]interface{}{
		"node_id":         bbn.ID(),
		"left_token_id":   leftToken.ID,
		"right_fact_id":   rightFact.ID,
		"result_token_id": resultToken.ID,
		"fact_count":      len(resultToken.Facts),
	})
}

// JoinNodeImpl implémente un nœud de jointure avec conditions.
type JoinNodeImpl struct {
	*BaseBetaNode
	joinConditions []domain.JoinCondition
	mutex          sync.RWMutex
}

// NewJoinNode crée un nouveau nœud de jointure.
func NewJoinNode(id string, logger domain.Logger) *JoinNodeImpl {
	return &JoinNodeImpl{
		BaseBetaNode:   NewBaseBetaNode(id, "JoinNode", logger),
		joinConditions: make([]domain.JoinCondition, 0),
	}
}

// SetJoinConditions définit les conditions de jointure.
func (jn *JoinNodeImpl) SetJoinConditions(conditions []domain.JoinCondition) {
	jn.mutex.Lock()
	defer jn.mutex.Unlock()
	jn.joinConditions = make([]domain.JoinCondition, len(conditions))
	copy(jn.joinConditions, conditions)
}

// GetJoinConditions retourne les conditions de jointure.
func (jn *JoinNodeImpl) GetJoinConditions() []domain.JoinCondition {
	jn.mutex.RLock()
	defer jn.mutex.RUnlock()

	conditions := make([]domain.JoinCondition, len(jn.joinConditions))
	copy(conditions, jn.joinConditions)
	return conditions
}

// EvaluateJoin évalue si un token et un fait peuvent être joints.
func (jn *JoinNodeImpl) EvaluateJoin(token *domain.Token, fact *domain.Fact) bool {
	conditions := jn.GetJoinConditions()

	// Si aucune condition, accepter la jointure (comportement par défaut)
	if len(conditions) == 0 {
		return true
	}

	// Toutes les conditions doivent être satisfaites (AND logique)
	for _, condition := range conditions {
		if !condition.Evaluate(token, fact) {
			return false
		}
	}

	return true
}

// ProcessLeftToken surcharge la méthode de base pour utiliser notre logique de jointure.
func (jn *JoinNodeImpl) ProcessLeftToken(token *domain.Token) error {
	jn.logTokenProcessing(token, "left_input")
	jn.betaMemory.StoreToken(token)

	// Essayer de faire des jointures avec les faits du côté droit
	rightFacts := jn.betaMemory.GetFacts()
	for _, fact := range rightFacts {
		if err := jn.tryJoin(token, fact); err != nil {
			return err
		}
	}

	return nil
}

// ProcessRightFact surcharge la méthode de base pour utiliser notre logique de jointure.
func (jn *JoinNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	jn.logFactProcessing(fact, "right_input")
	jn.betaMemory.StoreFact(fact)

	// Essayer de faire des jointures avec les tokens du côté gauche
	leftTokens := jn.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if err := jn.tryJoin(token, fact); err != nil {
			return err
		}
	}

	return nil
}

// tryJoin surcharge la méthode de base pour inclure l'évaluation des conditions.
func (jn *JoinNodeImpl) tryJoin(token *domain.Token, fact *domain.Fact) error {
	if !jn.EvaluateJoin(token, fact) {
		// Log de la jointure rejetée pour debug
		jn.logger.Debug("join rejected", map[string]interface{}{
			"node_id":       jn.ID(),
			"left_token_id": token.ID,
			"right_fact_id": fact.ID,
			"conditions":    len(jn.joinConditions),
		})
		return nil // Pas d'erreur, juste pas de jointure
	}

	// Utiliser la logique de jointure de la classe parent
	return jn.BaseBetaNode.tryJoin(token, fact)
}
