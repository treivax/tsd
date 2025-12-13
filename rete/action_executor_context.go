// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// ExecutionContext contient le contexte d'exécution d'une action.
//
// Le contexte fournit l'accès aux faits disponibles via BindingChain,
// permettant l'évaluation des arguments d'action référençant des variables.
//
// Structure :
//   - token : token contenant les faits déclencheurs
//   - network : réseau RETE pour accès aux types et autres ressources
//   - bindings : chaîne immuable de bindings variable → fact
//
// Utilisation :
//
//	ctx := NewExecutionContext(token, network)
//	fact := ctx.GetVariable("user")  // Récupère le fait lié à "user"
//
// Propriétés :
//   - Thread-safe grâce à l'immutabilité de BindingChain
//   - Accès O(n) aux variables (n = nombre de bindings, typiquement < 10)
//   - Pas de copie des faits, seulement des références
type ExecutionContext struct {
	token    *Token
	network  *ReteNetwork
	bindings *BindingChain
}

// NewExecutionContext crée un nouveau contexte d'exécution.
//
// Le contexte référence directement la chaîne de bindings du token,
// sans copie, garantissant l'immutabilité et la performance.
//
// Validation :
//   - Si token est nil, crée un contexte avec bindings vides
//   - Si network est nil, les fonctionnalités dépendant du network (type validation, etc.) ne seront pas disponibles
//
// Note : Un network nil est acceptable pour les tests unitaires simples,
// mais dans un contexte de production, le network devrait toujours être fourni.
//
// Paramètres :
//   - token : token contenant les faits et bindings (peut être nil)
//   - network : réseau RETE pour accès aux types (peut être nil pour tests simples)
//
// Retourne :
//   - *ExecutionContext : contexte d'exécution initialisé
func NewExecutionContext(token *Token, network *ReteNetwork) *ExecutionContext {
	ctx := &ExecutionContext{
		token:    token,
		network:  network,
		bindings: nil,
	}

	// Référencer directement la chaîne de bindings du token si disponible
	if token != nil {
		ctx.bindings = token.Bindings
	}

	return ctx
}

// GetVariable récupère un fait par nom de variable.
//
// Utilise la BindingChain pour rechercher le fait lié à la variable.
// Retourne nil si la variable n'existe pas dans le contexte.
//
// Complexité : O(n) où n est le nombre de bindings (typiquement < 10)
//
// Paramètres :
//   - name : nom de la variable (ex: "user", "order", "task")
//
// Retourne :
//   - *Fact : pointeur vers le fait si trouvé, nil sinon
//
// Exemple :
//
//	user := ctx.GetVariable("user")
//	if user == nil {
//	    return fmt.Errorf("variable user non trouvée")
//	}
//	userName := user.Fields["name"]
func (ctx *ExecutionContext) GetVariable(name string) *Fact {
	if ctx.bindings == nil {
		return nil
	}
	return ctx.bindings.Get(name)
}
