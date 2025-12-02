package rete

import (
	"fmt"
	"sync"
)

// RuleRouterNode is an intermediate node between a shared JoinNode and TerminalNodes
// It ensures that tokens from a shared JoinNode are properly routed to the correct
// TerminalNode for each rule, preventing token duplication when multiple rules
// share the same JoinNode due to identical conditions.
//
// Architecture:
// SharedJoinNode -> RuleRouterNode (Rule1) -> TerminalNode (Rule1)
//
//	-> RuleRouterNode (Rule2) -> TerminalNode (Rule2)
//
// This ensures each rule receives tokens independently, even when sharing the same JoinNode.
type RuleRouterNode struct {
	BaseNode
	RuleID       string        // ID of the rule this router serves
	JoinNodeID   string        // ID of the shared JoinNode this router is connected to
	TerminalNode *TerminalNode // The terminal node for this rule
	storage      Storage
	mutex        sync.RWMutex
}

// NewRuleRouterNode creates a new RuleRouterNode for routing tokens to a specific rule
func NewRuleRouterNode(ruleID string, joinNodeID string, storage Storage) *RuleRouterNode {
	nodeID := fmt.Sprintf("router_%s", ruleID)

	router := &RuleRouterNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "RuleRouterNode",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
		},
		RuleID:     ruleID,
		JoinNodeID: joinNodeID,
		storage:    storage,
	}

	return router
}

// ActivateLeft receives tokens from the shared JoinNode and routes them to the rule's TerminalNode
func (rrn *RuleRouterNode) ActivateLeft(token *Token) error {
	// Store the token in memory for tracking
	rrn.mutex.Lock()
	rrn.Memory.AddToken(token)
	rrn.mutex.Unlock()

	// Route the token to the terminal node if it exists
	if rrn.TerminalNode != nil {
		return rrn.TerminalNode.ActivateLeft(token)
	}

	// Otherwise propagate to children (fallback for flexibility)
	return rrn.PropagateToChildren(nil, token)
}

// ActivateRight is not used for RuleRouterNode (only receives from left/JoinNode)
func (rrn *RuleRouterNode) ActivateRight(fact *Fact) error {
	// RuleRouterNode doesn't receive facts directly, only tokens from JoinNode
	return nil
}

// ActivateRetract handles retraction of facts by removing tokens containing the fact
func (rrn *RuleRouterNode) ActivateRetract(factID string) error {
	rrn.mutex.Lock()

	// Remove tokens containing the retracted fact
	var tokensToRemove []string
	for tokenID, token := range rrn.Memory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				tokensToRemove = append(tokensToRemove, tokenID)
				break
			}
		}
	}

	for _, tokenID := range tokensToRemove {
		delete(rrn.Memory.Tokens, tokenID)
	}

	rrn.mutex.Unlock()

	// Propagate retraction to terminal node
	if rrn.TerminalNode != nil {
		return rrn.TerminalNode.ActivateRetract(factID)
	}

	// Otherwise propagate to children
	return rrn.PropagateRetractToChildren(factID)
}

// SetTerminalNode connects the router to its terminal node
func (rrn *RuleRouterNode) SetTerminalNode(terminal *TerminalNode) {
	rrn.mutex.Lock()
	defer rrn.mutex.Unlock()
	rrn.TerminalNode = terminal
}

// GetRuleID returns the rule ID this router serves
func (rrn *RuleRouterNode) GetRuleID() string {
	return rrn.RuleID
}

// GetJoinNodeID returns the ID of the shared JoinNode
func (rrn *RuleRouterNode) GetJoinNodeID() string {
	return rrn.JoinNodeID
}

// String returns a string representation of the RuleRouterNode
func (rrn *RuleRouterNode) String() string {
	return fmt.Sprintf("RuleRouterNode[id=%s, rule=%s, joinNode=%s]",
		rrn.ID, rrn.RuleID, rrn.JoinNodeID)
}
