// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

// Testing helpers for ProgramState
// These functions are exported for use in tests but should not be used in production code

// AddTypeForTesting adds a type definition directly (for testing only)
// This bypasses normal validation and should only be used in test code
func (ps *ProgramState) AddTypeForTesting(name string, typeDef *TypeDefinition) {
	if ps.types == nil {
		ps.types = make(map[string]*TypeDefinition)
	}
	ps.types[name] = typeDef
}

// AddRuleForTesting adds a rule directly (for testing only)
// This bypasses normal validation and should only be used in test code
func (ps *ProgramState) AddRuleForTesting(rule *Expression) {
	if ps.rules == nil {
		ps.rules = make([]*Expression, 0)
	}
	ps.rules = append(ps.rules, rule)
}

// AddFactForTesting adds a fact directly (for testing only)
// This bypasses normal validation and should only be used in test code
func (ps *ProgramState) AddFactForTesting(fact *Fact) {
	if ps.facts == nil {
		ps.facts = make([]*Fact, 0)
	}
	ps.facts = append(ps.facts, fact)
}

// GetRuleIDsCountForTesting returns the number of tracked rule IDs (for testing only)
func (ps *ProgramState) GetRuleIDsCountForTesting() int {
	return len(ps.ruleIDs)
}

// SetRuleIDForTesting sets a rule ID as used (for testing only)
func (ps *ProgramState) SetRuleIDForTesting(ruleID string) {
	if ps.ruleIDs == nil {
		ps.ruleIDs = make(map[string]bool)
	}
	ps.ruleIDs[ruleID] = true
}

// HasRuleIDForTesting checks if a rule ID is tracked (for testing only)
func (ps *ProgramState) HasRuleIDForTesting(ruleID string) bool {
	return ps.ruleIDs[ruleID]
}
