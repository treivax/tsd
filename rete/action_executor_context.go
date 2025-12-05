// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// ExecutionContext contient le contexte d'exécution d'une action,
// incluant les tokens et faits disponibles pour évaluation.
type ExecutionContext struct {
	token    *Token
	network  *ReteNetwork
	varCache map[string]*Fact
}

// NewExecutionContext crée un nouveau contexte d'exécution
func NewExecutionContext(token *Token, network *ReteNetwork) *ExecutionContext {
	ctx := &ExecutionContext{
		token:    token,
		network:  network,
		varCache: make(map[string]*Fact),
	}

	// Construire le cache des variables depuis le token
	if token != nil && token.Bindings != nil {
		for varName, fact := range token.Bindings {
			ctx.varCache[varName] = fact
		}
	}

	return ctx
}

// GetVariable récupère un fait par nom de variable
func (ctx *ExecutionContext) GetVariable(name string) *Fact {
	return ctx.varCache[name]
}
