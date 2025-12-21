// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"os"
	"time"

	"github.com/treivax/tsd/constraint"
)

// ConstraintPipeline implémente le pipeline complet :
// fichier .constraint → parseur PEG → conversion AST → réseau RETE
type ConstraintPipeline struct {
	logger                *Logger                                                     // Logger structuré pour instrumentation
	onXupleSpacesDetected func(network *ReteNetwork, definitions []interface{}) error // Callback appelé après détection des xuple-spaces
}

// GetLogger retourne le logger, en l'initialisant si nécessaire
func (cp *ConstraintPipeline) GetLogger() *Logger {
	if cp.logger == nil {
		cp.logger = NewLogger(LogLevelInfo, os.Stdout)
	}
	return cp.logger
}

// NewConstraintPipeline crée une nouvelle instance du pipeline
func NewConstraintPipeline() *ConstraintPipeline {
	return &ConstraintPipeline{
		logger: NewLogger(LogLevelInfo, os.Stdout), // Logger par défaut niveau Info
	}
}

// SetLogger configure le logger pour le pipeline
func (cp *ConstraintPipeline) SetLogger(logger *Logger) {
	if logger != nil {
		cp.logger = logger
	}
}

// SetOnXupleSpacesDetected configure le callback appelé après détection des xuple-spaces.
// Ce callback permet au package api de créer les xuple-spaces avant la soumission des faits inline.
func (cp *ConstraintPipeline) SetOnXupleSpacesDetected(callback func(network *ReteNetwork, definitions []interface{}) error) {
	cp.onXupleSpacesDetected = callback
}

// IngestFile est la fonction unique et incrémentale pour étendre le réseau RETE.
// Elle peut être appelée plusieurs fois avec des fichiers différents pour :
// - Parser le fichier (types, règles, faits)
// - Étendre le réseau RETE existant (ou créer un nouveau réseau si network == nil)
// - Propager les faits préexistants vers les nouvelles règles
// - Soumettre les nouveaux faits au réseau
//
// Cette fonction remplace toutes les anciennes variantes de pipeline.
// Le support de la commande 'reset' en milieu de fichier a été supprimé.
//
// IMPORTANT: Cette fonction utilise TOUJOURS les transactions avec auto-commit/auto-rollback.
// En cas d'erreur, la transaction est automatiquement annulée (rollback).
// En cas de succès, la transaction est automatiquement validée (commit).
//
// Les métriques sont toujours collectées et retournées (coût négligeable < 0.1%).
func (cp *ConstraintPipeline) IngestFile(filename string, network *ReteNetwork, storage Storage) (*ReteNetwork, *IngestionMetrics, error) {
	cp.logger.Info("========================================")
	cp.logger.Info("📁 Ingestion incrémentale: %s", filename)

	// Initialiser le contexte d'ingestion
	ctx := &ingestionContext{
		filename:              filename,
		network:               network,
		storage:               storage,
		metrics:               NewMetricsCollector(),
		xupleManager:          nil, // Sera créé si nécessaire lors de la détection de xuple-spaces
		onXupleSpacesDetected: cp.onXupleSpacesDetected,
	}

	// Exécuter le pipeline complet
	if err := cp.executePipeline(ctx); err != nil {
		return cp.handlePipelineError(ctx, err)
	}

	cp.logger.Info("🎯 INGESTION TERMINÉE")
	cp.logger.Info("========================================")

	return ctx.network, ctx.metrics.Finalize(), nil
}

// enrichProgramWithNetworkTypes merges types from the network into the program
// This is crucial for incremental validation when facts reference types defined in previous files
func (cp *ConstraintPipeline) enrichProgramWithNetworkTypes(program *constraint.Program, network *ReteNetwork) constraint.Program {
	// Create a copy of the program to avoid modifying the original
	enrichedProgram := *program

	// Build a map of types already in the program
	existingTypes := make(map[string]bool)
	for _, typeDef := range program.Types {
		existingTypes[typeDef.Name] = true
	}

	// Add types from the network that aren't already in the program
	for _, networkType := range network.Types {
		if !existingTypes[networkType.Name] {
			// Convert rete.TypeDefinition to constraint.TypeDefinition
			constraintType := constraint.TypeDefinition{
				Type:   networkType.Type,
				Name:   networkType.Name,
				Fields: make([]constraint.Field, len(networkType.Fields)),
			}

			// Convert fields
			for i, field := range networkType.Fields {
				constraintType.Fields[i] = constraint.Field{
					Name:         field.Name,
					Type:         field.Type,
					IsPrimaryKey: field.IsPrimaryKey,
				}
			}

			enrichedProgram.Types = append(enrichedProgram.Types, constraintType)
		}
	}

	return enrichedProgram
}

// executePipeline exécute toutes les étapes du pipeline d'ingestion
func (cp *ConstraintPipeline) executePipeline(ctx *ingestionContext) error {
	// Phase 1: Préparation
	if err := cp.prepareIngestion(ctx); err != nil {
		return err
	}

	// Phase 2: Construction réseau
	if err := cp.buildNetworkFromContext(ctx); err != nil {
		return err
	}

	// Phase 3: Gestion faits
	if err := cp.manageFacts(ctx); err != nil {
		return err
	}

	// Phase 4: Finalisation
	return cp.finalizeIngestion(ctx)
}

// prepareIngestion prépare le contexte d'ingestion (parsing, reset, transaction, validation)
func (cp *ConstraintPipeline) prepareIngestion(ctx *ingestionContext) error {
	if err := cp.readAndParseFile(ctx); err != nil {
		return err
	}

	if err := cp.detectReset(ctx); err != nil {
		return err
	}

	if err := cp.initializeOrResetNetwork(ctx); err != nil {
		return err
	}

	if err := cp.beginTransactionIfNeeded(ctx); err != nil {
		return err
	}

	return cp.validateConstraints(ctx)
}

// buildNetworkFromContext construit ou étend le réseau RETE à partir du contexte
func (cp *ConstraintPipeline) buildNetworkFromContext(ctx *ingestionContext) error {
	if err := cp.convertAndExtractComponents(ctx); err != nil {
		return err
	}

	if err := cp.extractXupleSpaces(ctx); err != nil {
		return err
	}

	if err := cp.createXupleSpaces(ctx); err != nil {
		return err
	}

	// Enregistrer l'action Xuple si un handler est configuré
	if err := cp.registerXupleActionIfNeeded(ctx); err != nil {
		return err
	}

	return cp.addTypesAndRules(ctx)
}

// manageFacts gère la collection et la propagation des faits
func (cp *ConstraintPipeline) manageFacts(ctx *ingestionContext) error {
	if err := cp.collectExistingFactsIfNeeded(ctx); err != nil {
		return err
	}

	if err := cp.propagateToNewRules(ctx); err != nil {
		return err
	}

	return cp.submitNewFacts(ctx)
}

// finalizeIngestion finalise l'ingestion avec validation et commit
func (cp *ConstraintPipeline) finalizeIngestion(ctx *ingestionContext) error {
	if err := cp.validateNetworkAndState(ctx); err != nil {
		return err
	}

	return cp.verifyAndCommit(ctx)
}

// handlePipelineError gère les erreurs du pipeline avec rollback automatique
func (cp *ConstraintPipeline) handlePipelineError(ctx *ingestionContext, err error) (*ReteNetwork, *IngestionMetrics, error) {
	if ctx.tx != nil && ctx.tx.IsActive {
		rollbackErr := ctx.tx.Rollback()
		if rollbackErr != nil {
			cp.logger.Error("❌ Erreur rollback: %v", rollbackErr)
			return ctx.network, ctx.metrics.Finalize(), fmt.Errorf("erreur ingestion: %w; erreur rollback: %v", err, rollbackErr)
		}
		cp.logger.Warn("🔙 Rollback automatique effectué")
	}
	return ctx.network, ctx.metrics.Finalize(), err
}

// readAndParseFile lit et parse le fichier de contraintes
func (cp *ConstraintPipeline) readAndParseFile(ctx *ingestionContext) error {
	parsingStart := time.Now()
	parsedAST, err := constraint.ParseConstraintFile(ctx.filename)
	if err != nil {
		ctx.metrics.RecordParsingDuration(time.Since(parsingStart))
		return fmt.Errorf("❌ Erreur parsing fichier %s: %w", ctx.filename, err)
	}
	ctx.parsedAST = parsedAST
	ctx.metrics.RecordParsingDuration(time.Since(parsingStart))
	cp.logger.Info("✅ Parsing réussi")
	return nil
}

// detectReset détecte la présence d'une commande reset dans l'AST
func (cp *ConstraintPipeline) detectReset(ctx *ingestionContext) error {
	resultMap, ok := ctx.parsedAST.(map[string]interface{})
	if !ok {
		return fmt.Errorf("❌ Format AST non reconnu: %T", ctx.parsedAST)
	}

	if resetsData, exists := resultMap["resets"]; exists {
		if resets, ok := resetsData.([]interface{}); ok && len(resets) > 0 {
			ctx.hasResets = true
			cp.logger.Info("🔄 Commande reset détectée - Réinitialisation complète du réseau")
		}
	}
	return nil
}

// initializeOrResetNetwork initialise ou réinitialise le réseau selon le contexte
func (cp *ConstraintPipeline) initializeOrResetNetwork(ctx *ingestionContext) error {
	if ctx.hasResets {
		cp.logger.Info("🔄 Commande reset détectée - Garbage Collection de l'ancien réseau")

		if ctx.network != nil {
			cp.logger.Debug("🗑️ GC du réseau existant...")
			ctx.network.GarbageCollect()
			cp.logger.Debug("✅ GC terminé")
		}

		cp.logger.Info("🆕 Création d'un nouveau réseau RETE")
		ctx.network = NewReteNetwork(ctx.storage)
		ctx.metrics.SetWasReset(true)
	}
	return nil
}

// beginTransactionIfNeeded démarre une transaction si le réseau existe
func (cp *ConstraintPipeline) beginTransactionIfNeeded(ctx *ingestionContext) error {
	if ctx.network != nil {
		ctx.tx = ctx.network.BeginTransaction()
		ctx.network.SetTransaction(ctx.tx)
		cp.logger.Info("🔒 Transaction démarrée automatiquement: %s", ctx.tx.ID)
	}
	return nil
}

// validateConstraints effectue la validation sémantique
func (cp *ConstraintPipeline) validateConstraints(ctx *ingestionContext) error {
	validationStart := time.Now()
	if ctx.network == nil || ctx.hasResets {
		// Validation standard
		if err := constraint.ValidateConstraintProgram(ctx.parsedAST); err != nil {
			return fmt.Errorf("❌ Erreur validation sémantique: %w", err)
		}
		cp.logger.Info("✅ Validation sémantique réussie")
		ctx.metrics.SetValidationSkipped(false)
	} else {
		// Validation incrémentale
		cp.logger.Info("🔍 Validation sémantique incrémentale avec contexte...")
		validator := NewIncrementalValidator(ctx.network)
		if err := validator.ValidateWithContext(ctx.parsedAST); err != nil {
			return fmt.Errorf("❌ Erreur validation incrémentale: %w", err)
		}
		cp.logger.Info("✅ Validation incrémentale réussie (%d types en contexte)", len(ctx.network.Types))
		ctx.metrics.SetValidationSkipped(false)
		ctx.metrics.SetWasIncremental(true)
	}
	ctx.metrics.RecordValidationDuration(time.Since(validationStart))
	return nil
}

// convertAndExtractComponents convertit l'AST en programme RETE et extrait les composants
func (cp *ConstraintPipeline) convertAndExtractComponents(ctx *ingestionContext) error {
	// Conversion en programme
	program, err := constraint.ConvertResultToProgram(ctx.parsedAST)
	if err != nil {
		return fmt.Errorf("❌ Erreur conversion programme: %w", err)
	}
	ctx.program = program

	// Créer ou étendre le réseau
	if ctx.network == nil {
		cp.logger.Info("🆕 Création d'un nouveau réseau RETE")
		ctx.network = NewReteNetwork(ctx.storage)
	} else if !ctx.hasResets {
		cp.logger.Info("🔄 Extension du réseau RETE existant")
	}

	// Convertir au format RETE
	reteProgram, err := constraint.ConvertToReteProgram(program)
	if err != nil {
		return fmt.Errorf("❌ Erreur conversion programme RETE: %w", err)
	}
	ctx.reteProgram = reteProgram

	reteResultMap, ok := ctx.reteProgram.(map[string]interface{})
	if !ok {
		return fmt.Errorf("❌ Format programme RETE invalide: %T", ctx.reteProgram)
	}

	// Extraire les composants
	types, expressions, err := cp.extractComponents(reteResultMap)
	if err != nil {
		return fmt.Errorf("❌ Erreur extraction composants: %w", err)
	}
	ctx.types = types
	ctx.expressions = expressions
	cp.logger.Info("✅ Trouvé %d types et %d expressions dans le fichier", len(types), len(expressions))

	return nil
}

// addTypesAndRules ajoute les types et les règles au réseau
func (cp *ConstraintPipeline) addTypesAndRules(ctx *ingestionContext) error {
	// Ajout types
	if len(ctx.types) > 0 {
		typeCreationStart := time.Now()
		if err := cp.createTypeNodes(ctx.network, ctx.types, ctx.storage); err != nil {
			return fmt.Errorf("❌ Erreur ajout types: %w", err)
		}
		cp.logger.Info("✅ Types ajoutés/mis à jour dans le réseau")
		ctx.metrics.RecordTypeCreationDuration(time.Since(typeCreationStart))
		ctx.metrics.SetTypesAdded(len(ctx.types))
	}

	// Extraire et stocker les actions
	reteResultMap, _ := ctx.reteProgram.(map[string]interface{})
	if err := cp.extractAndStoreActions(ctx.network, reteResultMap); err != nil {
		return fmt.Errorf("❌ Erreur extraction actions: %w", err)
	}

	// Identifier terminaux existants
	ctx.existingTerminals = make(map[string]bool)
	for terminalID := range ctx.network.TerminalNodes {
		ctx.existingTerminals[terminalID] = true
	}

	// Ajouter les nouvelles règles
	if len(ctx.expressions) > 0 {
		ruleCreationStart := time.Now()
		if err := cp.createRuleNodes(ctx.network, ctx.expressions, ctx.storage); err != nil {
			return fmt.Errorf("❌ Erreur ajout règles: %w", err)
		}
		cp.logger.Info("✅ Règles ajoutées au réseau")
		ctx.metrics.RecordRuleCreationDuration(time.Since(ruleCreationStart))
		ctx.metrics.SetRulesAdded(len(ctx.expressions))
	}

	// Traiter les suppressions de règles
	if err := cp.processRuleRemovals(ctx.network, reteResultMap); err != nil {
		return fmt.Errorf("❌ Erreur traitement suppressions de règles: %w", err)
	}

	return nil
}

// collectExistingFactsIfNeeded collecte les faits existants si nécessaire
func (cp *ConstraintPipeline) collectExistingFactsIfNeeded(ctx *ingestionContext) error {
	if ctx.hasResets {
		cp.logger.Debug("📊 Réseau réinitialisé - pas de faits préexistants")
	} else {
		collectionStart := time.Now()
		ctx.existingFacts = cp.collectExistingFacts(ctx.network)
		ctx.factsByType = cp.organizeFactsByType(ctx.existingFacts)
		cp.logger.Debug("📊 Faits préexistants dans le réseau: %d", len(ctx.existingFacts))
		ctx.metrics.RecordFactCollectionDuration(time.Since(collectionStart))
		ctx.metrics.SetExistingFactsCollected(len(ctx.existingFacts))
	}
	return nil
}

// propagateToNewRules propage les faits existants vers les nouvelles règles
func (cp *ConstraintPipeline) propagateToNewRules(ctx *ingestionContext) error {
	ctx.newTerminals = cp.identifyNewTerminals(ctx.network, ctx.existingTerminals)
	if len(ctx.newTerminals) > 0 && len(ctx.existingFacts) > 0 {
		cp.logger.Info("🔄 Propagation ciblée de faits vers %d nouvelle(s) règle(s)", len(ctx.newTerminals))
		propagationStart := time.Now()
		propagatedCount := cp.propagateToNewTerminals(ctx.network, ctx.newTerminals, ctx.factsByType)
		ctx.metrics.RecordPropagationDuration(time.Since(propagationStart))
		ctx.metrics.SetFactsPropagated(propagatedCount)
		ctx.metrics.SetNewTerminalsAdded(len(ctx.newTerminals))
		ctx.metrics.SetPropagationTargets(len(ctx.newTerminals))
		cp.logger.Info("✅ Propagation rétroactive terminée (%d fait(s) propagé(s))", propagatedCount)
	}
	return nil
}

// submitNewFacts soumet les nouveaux faits au réseau
func (cp *ConstraintPipeline) submitNewFacts(ctx *ingestionContext) error {
	if len(ctx.program.Facts) > 0 {
		// CRUCIAL: Merge network types into program for incremental validation
		// When loading facts from a separate file, the program only contains types
		// defined in that file. We need to merge in types from previous files.
		programWithAllTypes := cp.enrichProgramWithNetworkTypes(ctx.program, ctx.network)

		factsForRete, err := constraint.ConvertFactsToReteFormat(programWithAllTypes)
		if err != nil {
			return fmt.Errorf("❌ Erreur conversion faits: %w", err)
		}
		ctx.factsForRete = factsForRete
		cp.logger.Info("📥 Soumission de %d nouveaux faits", len(ctx.factsForRete))
		submissionStart := time.Now()
		if err := ctx.network.SubmitFactsFromGrammar(ctx.factsForRete); err != nil {
			return fmt.Errorf("❌ Erreur soumission faits: %w", err)
		}
		cp.logger.Info("✅ Nouveaux faits soumis")
		ctx.metrics.RecordFactSubmissionDuration(time.Since(submissionStart))
		ctx.metrics.SetFactsSubmitted(len(ctx.factsForRete))
	}
	return nil
}

// validateNetworkAndState valide le réseau et enregistre son état
func (cp *ConstraintPipeline) validateNetworkAndState(ctx *ingestionContext) error {
	if err := cp.validateNetwork(ctx.network); err != nil {
		return fmt.Errorf("❌ Erreur validation réseau: %w", err)
	}
	cp.logger.Info("✅ Validation réussie")

	// Enregistrer l'état du réseau
	ctx.metrics.RecordNetworkState(ctx.network)
	cp.logger.Info("🎯 INGESTION INCRÉMENTALE TERMINÉE")
	cp.logger.Info("   - Total TypeNodes: %d", len(ctx.network.TypeNodes))
	cp.logger.Info("   - Total TerminalNodes: %d", len(ctx.network.TerminalNodes))

	return nil
}

// extractXupleSpaces extrait les xuple-spaces depuis l'AST parsé
func (cp *ConstraintPipeline) extractXupleSpaces(ctx *ingestionContext) error {
	resultMap, ok := ctx.parsedAST.(map[string]interface{})
	if !ok {
		return fmt.Errorf("format AST invalide pour extraction xuple-spaces")
	}

	xupleSpacesData, exists := resultMap["xupleSpaces"]
	if !exists {
		// Pas de xuple-spaces dans ce fichier, ce n'est pas une erreur
		return nil
	}

	xupleSpacesList, ok := xupleSpacesData.([]interface{})
	if !ok {
		return fmt.Errorf("format xupleSpaces invalide: %T", xupleSpacesData)
	}

	ctx.xupleSpaces = xupleSpacesList
	cp.logger.Info("✅ Trouvé %d xuple-space(s) à créer", len(xupleSpacesList))

	return nil
}

// createXupleSpaces stocke les définitions de xuple-spaces détectées lors du parsing.
// Les définitions sont stockées dans le réseau et un callback optionnel est appelé
// pour permettre au package api de créer les xuple-spaces avant la soumission des faits inline.
//
// Note: La création effective des xuple-spaces est gérée par le callback (package api),
// pas par ce pipeline. Le pipeline se contente de parser et stocker les définitions.
func (cp *ConstraintPipeline) createXupleSpaces(ctx *ingestionContext) error {
	if len(ctx.xupleSpaces) == 0 {
		return nil
	}

	cp.logger.Info("───────────────────────────────────────────────────────────────")
	cp.logger.Info("📦 DÉTECTION DES XUPLE-SPACES")
	cp.logger.Info("───────────────────────────────────────────────────────────────")

	// Stocker les définitions dans le réseau pour utilisation par le package api
	ctx.network.SetXupleSpaceDefinitions(ctx.xupleSpaces)
	cp.logger.Info("   %d xuple-space(s) détecté(s) et stocké(s)", len(ctx.xupleSpaces))

	// Appeler le callback si configuré (permet au package api de créer les xuple-spaces immédiatement)
	if ctx.onXupleSpacesDetected != nil {
		cp.logger.Info("   Création des xuple-spaces via callback...")
		if err := ctx.onXupleSpacesDetected(ctx.network, ctx.xupleSpaces); err != nil {
			return fmt.Errorf("erreur callback création xuple-spaces: %w", err)
		}
		cp.logger.Info("   ✅ Xuple-spaces créés via callback")
	} else {
		cp.logger.Info("   Les xuple-spaces seront créés automatiquement par le package api")
	}

	cp.logger.Info("")
	return nil
}

// registerXupleActionIfNeeded enregistre l'action Xuple si un handler est configuré.
// Cette méthode est appelée après createXupleSpaces pour s'assurer que l'action
// est disponible même si aucun xuple-space n'est déclaré dans le fichier TSD.
//
// L'enregistrement est désormais délégué à ActionExecutor.RegisterXupleActionIfNeeded()
// qui gère automatiquement la vérification et l'enregistrement.
func (cp *ConstraintPipeline) registerXupleActionIfNeeded(ctx *ingestionContext) error {
	// Vérifier si un ActionExecutor est disponible
	if ctx.network == nil || ctx.network.ActionExecutor == nil {
		return nil
	}

	// Déléguer l'enregistrement à l'ActionExecutor
	if err := ctx.network.ActionExecutor.RegisterXupleActionIfNeeded(); err != nil {
		return fmt.Errorf("erreur enregistrement action Xuple: %w", err)
	}

	return nil
}

// verifyAndCommit vérifie la cohérence et commit la transaction
func (cp *ConstraintPipeline) verifyAndCommit(ctx *ingestionContext) error {
	// Vérification de cohérence pré-commit
	if ctx.tx != nil && ctx.tx.IsActive && len(ctx.factsForRete) > 0 {
		cp.logger.Info("🔍 Vérification de cohérence pré-commit...")

		expectedFactCount := len(ctx.factsForRete)
		actualFactCount := 0
		missingFacts := make([]string, 0)

		for i, factMap := range ctx.factsForRete {
			// Utiliser _id_ qui contient déjà le format Type~Value
			var internalID string
			if id, ok := factMap["_id_"].(string); ok {
				internalID = id
			} else {
				// Fallback: construire l'ID si _id_ n'existe pas
				var factID string
				if id, ok := factMap["id"].(string); ok {
					factID = id
				} else {
					factID = fmt.Sprintf("fact_%d", i)
				}

				factType := "unknown"
				if typ, ok := factMap["type"].(string); ok {
					factType = typ
				} else if typ, ok := factMap["reteType"].(string); ok {
					factType = typ
				}

				internalID = fmt.Sprintf("%s~%s", factType, factID)
			}

			if ctx.storage.GetFact(internalID) != nil {
				actualFactCount++
			} else {
				missingFacts = append(missingFacts, internalID)
			}
		}

		if expectedFactCount != actualFactCount {
			cp.logger.Error("❌ Incohérence détectée: %d faits attendus, %d trouvés", expectedFactCount, actualFactCount)
			cp.logger.Error("   Faits manquants: %v", missingFacts)
			return fmt.Errorf(
				"incohérence pré-commit: %d faits attendus mais %d trouvés dans le storage",
				expectedFactCount, actualFactCount)
		}

		cp.logger.Info("✅ Cohérence vérifiée: %d/%d faits présents", actualFactCount, expectedFactCount)

		// Synchroniser le storage
		cp.logger.Info("💾 Synchronisation du storage...")
		if err := ctx.storage.Sync(); err != nil {
			return fmt.Errorf("❌ Erreur sync storage: %w", err)
		}
		cp.logger.Info("✅ Storage synchronisé")
	}

	// Commit transaction
	if ctx.tx != nil && ctx.tx.IsActive {
		if err := ctx.tx.Commit(); err != nil {
			return fmt.Errorf("❌ Erreur commit transaction: %w", err)
		}
		cp.logger.Info("✅ Transaction committée: %d changements", ctx.tx.GetCommandCount())
	}

	return nil
}
