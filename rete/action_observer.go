// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "time"

// ActionObserver permet d'observer les exécutions d'actions.
//
// Cette interface implémente le pattern Observer pour découpler
// le moteur RETE de la collecte/traitement des activations.
//
// Utilisations typiques :
//   - Tests : capturer les exécutions pour assertions
//   - Xuples : publier vers xuple-spaces
//   - Métriques : collecter statistiques d'exécution
//   - Logging : journalisation centralisée
//   - Audit : traçabilité des actions
//
// Thread-Safety :
//   - Les implémentations DOIVENT être thread-safe
//   - Plusieurs terminal nodes peuvent notifier en parallèle
//   - L'observer ne doit PAS bloquer l'exécution
type ActionObserver interface {
	// OnActionExecuted est appelé après chaque exécution d'action.
	//
	// L'appel se fait de manière synchrone après l'exécution.
	// L'implémentation NE DOIT PAS bloquer longtemps.
	// Pour traitements longs, utiliser une goroutine interne.
	//
	// Paramètres :
	//   - result : résultat de l'exécution avec contexte complet
	OnActionExecuted(result ExecutionResult)
}

// ExecutionResult représente le résultat de l'exécution d'une action.
type ExecutionResult struct {
	Success   bool          // true si l'exécution a réussi
	Error     error         // erreur si l'exécution a échoué
	Duration  time.Duration // durée d'exécution
	Context   ActionContext // contexte d'exécution complet
	Arguments []interface{} // arguments évalués
}

// ActionContext contient le contexte d'exécution d'une action.
type ActionContext struct {
	ActionName string       // Nom de l'action (ex: "print", "insert")
	RuleName   string       // Nom de la règle qui a déclenché l'action
	Token      *Token       // Token déclencheur avec tous les faits
	Network    *ReteNetwork // Réseau RETE (pour insert/update/retract)
	Timestamp  time.Time    // Moment de l'exécution
}

// NoOpObserver est un observateur qui ne fait rien.
// Utilisé comme valeur par défaut pour éviter les vérifications nil.
type NoOpObserver struct{}

// OnActionExecuted ne fait rien (implémentation vide).
func (n *NoOpObserver) OnActionExecuted(result ExecutionResult) {
	// Intentionally empty
}
