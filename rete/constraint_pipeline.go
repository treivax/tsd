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
	logger *Logger // Logger structuré pour instrumentation
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

	// Initialiser la collecte de métriques
	metrics := NewMetricsCollector()

	// Initialiser le contexte d'ingestion
	ctx := &ingestionContext{
		filename: filename,
		network:  network,
		storage:  storage,
		metrics:  metrics,
	}

	// ÉTAPE 1: Parsing et détection reset
	parsingStart := time.Now()
	parsedAST, err := constraint.ParseConstraintFile(ctx.filename)
	if err != nil {
		metrics.RecordParsingDuration(time.Since(parsingStart))
		return nil, metrics.Finalize(), fmt.Errorf("❌ Erreur parsing fichier %s: %w", ctx.filename, err)
	}
	ctx.parsedAST = parsedAST
	metrics.RecordParsingDuration(time.Since(parsingStart))
	cp.logger.Info("✅ Parsing réussi")

	// Détecter reset
	resultMap, ok := parsedAST.(map[string]interface{})
	if !ok {
		return nil, metrics.Finalize(), fmt.Errorf("❌ Format AST non reconnu: %T", parsedAST)
	}

	if resetsData, exists := resultMap["resets"]; exists {
		if resets, ok := resetsData.([]interface{}); ok && len(resets) > 0 {
			ctx.hasResets = true
			cp.logger.Info("🔄 Commande reset détectée - Réinitialisation complète du réseau")
		}
	}

	// ÉTAPE 2: Initialisation réseau (GC si reset)
	if ctx.hasResets {
		cp.logger.Info("🔄 Commande reset détectée - Garbage Collection de l'ancien réseau")

		if ctx.network != nil {
			cp.logger.Debug("🗑️ GC du réseau existant...")
			ctx.network.GarbageCollect()
			cp.logger.Debug("✅ GC terminé")
		}

		cp.logger.Info("🆕 Création d'un nouveau réseau RETE")
		ctx.network = NewReteNetwork(ctx.storage)
		metrics.SetWasReset(true)
	}

	// ÉTAPE 3: Démarrer transaction
	if ctx.network != nil {
		ctx.tx = ctx.network.BeginTransaction()
		ctx.network.SetTransaction(ctx.tx)
		cp.logger.Info("🔒 Transaction démarrée automatiquement: %s", ctx.tx.ID)
	}

	// Wrapper pour rollback automatique en cas d'erreur
	handleError := func(err error) (*ReteNetwork, *IngestionMetrics, error) {
		if ctx.tx != nil && ctx.tx.IsActive {
			rollbackErr := ctx.tx.Rollback()
			if rollbackErr != nil {
				cp.logger.Error("❌ Erreur rollback: %v", rollbackErr)
				return ctx.network, metrics.Finalize(), fmt.Errorf("erreur ingestion: %w; erreur rollback: %v", err, rollbackErr)
			}
			cp.logger.Warn("🔙 Rollback automatique effectué")
		}
		return ctx.network, metrics.Finalize(), err
	}

	// ÉTAPE 4: Validation sémantique
	validationStart := time.Now()
	if ctx.network == nil || ctx.hasResets {
		// Validation standard
		if err := constraint.ValidateConstraintProgram(ctx.parsedAST); err != nil {
			return handleError(fmt.Errorf("❌ Erreur validation sémantique: %w", err))
		}
		cp.logger.Info("✅ Validation sémantique réussie")
		metrics.SetValidationSkipped(false)
	} else {
		// Validation incrémentale
		cp.logger.Info("🔍 Validation sémantique incrémentale avec contexte...")
		validator := NewIncrementalValidator(ctx.network)
		if err := validator.ValidateWithContext(ctx.parsedAST); err != nil {
			return handleError(fmt.Errorf("❌ Erreur validation incrémentale: %w", err))
		}
		cp.logger.Info("✅ Validation incrémentale réussie (%d types en contexte)", len(ctx.network.Types))
		metrics.SetValidationSkipped(false)
		metrics.SetWasIncremental(true)
	}
	metrics.RecordValidationDuration(time.Since(validationStart))

	// ÉTAPE 5: Conversion en programme RETE
	program, err := constraint.ConvertResultToProgram(ctx.parsedAST)
	if err != nil {
		return handleError(fmt.Errorf("❌ Erreur conversion programme: %w", err))
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
		return handleError(fmt.Errorf("❌ Erreur conversion programme RETE: %w", err))
	}
	ctx.reteProgram = reteProgram
	reteResultMap, ok := ctx.reteProgram.(map[string]interface{})
	if !ok {
		return handleError(fmt.Errorf("❌ Format programme RETE invalide: %T", ctx.reteProgram))
	}

	// Extraire les composants
	types, expressions, err := cp.extractComponents(reteResultMap)
	if err != nil {
		return handleError(fmt.Errorf("❌ Erreur extraction composants: %w", err))
	}
	ctx.types = types
	ctx.expressions = expressions
	cp.logger.Info("✅ Trouvé %d types et %d expressions dans le fichier", len(types), len(expressions))

	// ÉTAPE 6: Ajout types et actions
	if len(ctx.types) > 0 {
		typeCreationStart := time.Now()
		if err := cp.createTypeNodes(ctx.network, ctx.types, ctx.storage); err != nil {
			return handleError(fmt.Errorf("❌ Erreur ajout types: %w", err))
		}
		cp.logger.Info("✅ Types ajoutés/mis à jour dans le réseau")
		metrics.RecordTypeCreationDuration(time.Since(typeCreationStart))
		metrics.SetTypesAdded(len(ctx.types))
	}

	// Extraire et stocker les actions
	if err := cp.extractAndStoreActions(ctx.network, reteResultMap); err != nil {
		return handleError(fmt.Errorf("❌ Erreur extraction actions: %w", err))
	}

	// ÉTAPE 7: Collection faits existants
	if ctx.hasResets {
		cp.logger.Debug("📊 Réseau réinitialisé - pas de faits préexistants")
	} else {
		collectionStart := time.Now()
		ctx.existingFacts = cp.collectExistingFacts(ctx.network)
		ctx.factsByType = cp.organizeFactsByType(ctx.existingFacts)
		cp.logger.Debug("📊 Faits préexistants dans le réseau: %d", len(ctx.existingFacts))
		metrics.RecordFactCollectionDuration(time.Since(collectionStart))
		metrics.SetExistingFactsCollected(len(ctx.existingFacts))
	}

	// ÉTAPE 8: Gestion des règles (ajout + suppression)
	ctx.existingTerminals = make(map[string]bool)
	for terminalID := range ctx.network.TerminalNodes {
		ctx.existingTerminals[terminalID] = true
	}

	// Ajouter les nouvelles règles
	if len(ctx.expressions) > 0 {
		ruleCreationStart := time.Now()
		if err := cp.createRuleNodes(ctx.network, ctx.expressions, ctx.storage); err != nil {
			return handleError(fmt.Errorf("❌ Erreur ajout règles: %w", err))
		}
		cp.logger.Info("✅ Règles ajoutées au réseau")
		metrics.RecordRuleCreationDuration(time.Since(ruleCreationStart))
		metrics.SetRulesAdded(len(ctx.expressions))
	}

	// Traiter les suppressions de règles
	if err := cp.processRuleRemovals(ctx.network, reteResultMap); err != nil {
		return handleError(fmt.Errorf("❌ Erreur traitement suppressions de règles: %w", err))
	}

	// ÉTAPE 9: Propagation rétroactive vers nouvelles règles
	ctx.newTerminals = cp.identifyNewTerminals(ctx.network, ctx.existingTerminals)
	if len(ctx.newTerminals) > 0 && len(ctx.existingFacts) > 0 {
		cp.logger.Info("🔄 Propagation ciblée de faits vers %d nouvelle(s) règle(s)", len(ctx.newTerminals))
		propagationStart := time.Now()
		propagatedCount := cp.propagateToNewTerminals(ctx.network, ctx.newTerminals, ctx.factsByType)
		metrics.RecordPropagationDuration(time.Since(propagationStart))
		metrics.SetFactsPropagated(propagatedCount)
		metrics.SetNewTerminalsAdded(len(ctx.newTerminals))
		metrics.SetPropagationTargets(len(ctx.newTerminals))
		cp.logger.Info("✅ Propagation rétroactive terminée (%d fait(s) propagé(s))", propagatedCount)
	}

	// ÉTAPE 10: Soumission nouveaux faits
	if len(ctx.program.Facts) > 0 {
		ctx.factsForRete = constraint.ConvertFactsToReteFormat(*ctx.program)
		cp.logger.Info("📥 Soumission de %d nouveaux faits", len(ctx.factsForRete))
		submissionStart := time.Now()
		if err := ctx.network.SubmitFactsFromGrammar(ctx.factsForRete); err != nil {
			return handleError(fmt.Errorf("❌ Erreur soumission faits: %w", err))
		}
		cp.logger.Info("✅ Nouveaux faits soumis")
		metrics.RecordFactSubmissionDuration(time.Since(submissionStart))
		metrics.SetFactsSubmitted(len(ctx.factsForRete))
	}

	// ÉTAPE 11: Validation finale et cohérence
	if err := cp.validateNetwork(ctx.network); err != nil {
		return handleError(fmt.Errorf("❌ Erreur validation réseau: %w", err))
	}
	cp.logger.Info("✅ Validation réussie")

	// Enregistrer l'état du réseau
	metrics.RecordNetworkState(ctx.network)
	cp.logger.Info("🎯 INGESTION INCRÉMENTALE TERMINÉE")
	cp.logger.Info("   - Total TypeNodes: %d", len(ctx.network.TypeNodes))
	cp.logger.Info("   - Total TerminalNodes: %d", len(ctx.network.TerminalNodes))

	// Vérification de cohérence pré-commit
	if ctx.tx != nil && ctx.tx.IsActive && len(ctx.factsForRete) > 0 {
		cp.logger.Info("🔍 Vérification de cohérence pré-commit...")

		expectedFactCount := len(ctx.factsForRete)
		actualFactCount := 0
		missingFacts := make([]string, 0)

		for i, factMap := range ctx.factsForRete {
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

			internalID := fmt.Sprintf("%s_%s", factType, factID)

			if ctx.storage.GetFact(internalID) != nil {
				actualFactCount++
			} else {
				missingFacts = append(missingFacts, internalID)
			}
		}

		if expectedFactCount != actualFactCount {
			cp.logger.Error("❌ Incohérence détectée: %d faits attendus, %d trouvés", expectedFactCount, actualFactCount)
			cp.logger.Error("   Faits manquants: %v", missingFacts)
			return handleError(fmt.Errorf(
				"incohérence pré-commit: %d faits attendus mais %d trouvés dans le storage",
				expectedFactCount, actualFactCount))
		}

		cp.logger.Info("✅ Cohérence vérifiée: %d/%d faits présents", actualFactCount, expectedFactCount)

		// Synchroniser le storage
		cp.logger.Info("💾 Synchronisation du storage...")
		if err := ctx.storage.Sync(); err != nil {
			return handleError(fmt.Errorf("❌ Erreur sync storage: %w", err))
		}
		cp.logger.Info("✅ Storage synchronisé")
	}

	// ÉTAPE 12: Commit transaction
	if ctx.tx != nil && ctx.tx.IsActive {
		if err := ctx.tx.Commit(); err != nil {
			return handleError(fmt.Errorf("❌ Erreur commit transaction: %w", err))
		}
		cp.logger.Info("✅ Transaction committée: %d changements", ctx.tx.GetCommandCount())
	}

	cp.logger.Info("🎯 INGESTION TERMINÉE")
	cp.logger.Info("========================================")

	return ctx.network, metrics.Finalize(), nil
}
