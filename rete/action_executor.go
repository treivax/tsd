// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"log"
)

// ActionExecutor gère l'exécution des actions déclenchées par les règles.
//
// Responsabilités :
//   - Évaluation des arguments d'action (variables, expressions, etc.)
//   - Validation des arguments selon le type attendu
//   - Exécution via le registry d'actions (handlers personnalisés)
//   - Logging des actions exécutées
//
// Architecture :
//   - registry : enregistre les handlers d'actions disponibles
//   - network : accès au réseau RETE pour types et autres ressources
//   - logger : journalisation des exécutions
//
// Utilisation :
//
//	executor := NewActionExecutor(network, logger)
//	executor.RegisterAction(customHandler)
//	err := executor.ExecuteAction(action, token)
type ActionExecutor struct {
	network       *ReteNetwork
	logger        *log.Logger
	enableLogging bool
	registry      *ActionRegistry
}

// NewActionExecutor crée un nouveau exécuteur d'actions.
//
// Initialise le registry et enregistre les actions par défaut (print, etc.).
//
// Paramètres :
//   - network : réseau RETE
//   - logger : logger pour journalisation (utilise log.Default() si nil)
//
// Retourne :
//   - *ActionExecutor : exécuteur initialisé
func NewActionExecutor(network *ReteNetwork, logger *log.Logger) *ActionExecutor {
	if logger == nil {
		logger = log.Default()
	}
	ae := &ActionExecutor{
		network:       network,
		logger:        logger,
		enableLogging: true,
		registry:      NewActionRegistry(),
	}

	// Enregistrer les actions par défaut
	ae.RegisterDefaultActions()

	return ae
}

// RegisterDefaultActions enregistre les actions par défaut disponibles.
//
// Actions enregistrées :
//   - print : affichage de valeurs
//
// Cette méthode est appelée automatiquement par NewActionExecutor.
func (ae *ActionExecutor) RegisterDefaultActions() {
	// Enregistrer l'action print
	printAction := NewPrintAction(nil)
	if err := ae.registry.Register(printAction); err != nil {
		ae.logger.Printf("⚠️  Erreur enregistrement action print: %v", err)
	}
}

// GetRegistry retourne le registry d'actions.
func (ae *ActionExecutor) GetRegistry() *ActionRegistry {
	return ae.registry
}

// RegisterAction enregistre une action personnalisée.
func (ae *ActionExecutor) RegisterAction(handler ActionHandler) error {
	return ae.registry.Register(handler)
}

// SetLogging active ou désactive le logging des actions
func (ae *ActionExecutor) SetLogging(enabled bool) {
	ae.enableLogging = enabled
}

// ExecuteAction exécute une action avec les faits fournis par le token.
//
// Process :
//  1. Récupère tous les jobs de l'action
//  2. Crée un contexte d'exécution avec les bindings du token
//  3. Exécute chaque job en séquence
//
// Paramètres :
//   - action : action à exécuter (peut contenir plusieurs jobs)
//   - token : token contenant les faits et bindings disponibles
//
// Retourne :
//   - error : erreur si l'exécution échoue
func (ae *ActionExecutor) ExecuteAction(action *Action, token *Token) error {
	if action == nil {
		return fmt.Errorf("action is nil")
	}

	// Obtenir tous les jobs à exécuter
	jobs := action.GetJobs()

	// Créer un contexte d'exécution avec les faits disponibles
	ctx := NewExecutionContext(token, ae.network)

	// Exécuter chaque job en séquence
	for i, job := range jobs {
		if err := ae.executeJob(job, ctx, i); err != nil {
			return fmt.Errorf("erreur exécution job %s: %w", job.Name, err)
		}
	}

	return nil
}

// executeJob exécute un job individuel.
//
// Process :
//  1. Log l'action (si activé)
//  2. Évalue tous les arguments
//  3. Recherche le handler dans le registry
//  4. Valide les arguments (si handler définit une validation)
//  5. Exécute le handler
//
// Paramètres :
//   - job : job à exécuter
//   - ctx : contexte d'exécution
//   - jobIndex : index du job dans la séquence (pour debug)
//
// Retourne :
//   - error : erreur si l'exécution échoue
func (ae *ActionExecutor) executeJob(job JobCall, ctx *ExecutionContext, jobIndex int) error {
	// Logger l'action
	if ae.enableLogging {
		ae.logAction(job, ctx)
	}

	// Évaluer les arguments
	evaluatedArgs := make([]interface{}, 0, len(job.Args))
	for i, arg := range job.Args {
		evaluated, err := ae.evaluateArgument(arg, ctx)
		if err != nil {
			return fmt.Errorf("erreur évaluation argument %d: %w", i, err)
		}
		evaluatedArgs = append(evaluatedArgs, evaluated)
	}

	// Vérifier si un handler est enregistré pour cette action
	handler := ae.registry.Get(job.Name)
	if handler != nil {
		// Valider les arguments (optionnel)
		if err := handler.Validate(evaluatedArgs); err != nil {
			return fmt.Errorf("validation échouée pour action '%s': %w", job.Name, err)
		}

		// Exécuter l'action via son handler
		if err := handler.Execute(evaluatedArgs, ctx); err != nil {
			return fmt.Errorf("exécution échouée pour action '%s': %w", job.Name, err)
		}

		// Logger le succès
		ae.logger.Printf("🎯 ACTION EXÉCUTÉE: %s(%v)", job.Name, formatArgs(evaluatedArgs))
	} else {
		// Aucun handler défini : comportement par défaut (simple log)
		ae.logger.Printf("📋 ACTION NON DÉFINIE (log uniquement): %s(%v)", job.Name, formatArgs(evaluatedArgs))
	}

	return nil
}
