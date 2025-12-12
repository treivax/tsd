// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// ExecutionContext contient le contexte d'exécution d'une action,
// incluant les tokens et faits disponibles pour évaluation.
//
// Utilise maintenant BindingChain directement pour un accès immuable aux variables.
type ExecutionContext struct {
	token    *Token
	network  *ReteNetwork
	bindings *BindingChain
}

// NewExecutionContext crée un nouveau contexte d'exécution
func NewExecutionContext(token *Token, network *ReteNetwork) *ExecutionContext {
	ctx := &ExecutionContext{
		token:    token,
		network:  network,
		bindings: nil,
	}

	// Référencer directement la chaîne de bindings du token
	if token != nil {
		ctx.bindings = token.Bindings
	}

	return ctx
}

// GetVariable récupère un fait par nom de variable
func (ctx *ExecutionContext) GetVariable(name string) *Fact {
	if ctx.bindings == nil {
		return nil
	}
	return ctx.bindings.Get(name)
}
