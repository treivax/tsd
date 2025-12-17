// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/constraint"
)

// ingestionContext encapsule l'état d'une ingestion de fichier
type ingestionContext struct {
	filename          string
	network           *ReteNetwork
	storage           Storage
	metrics           *MetricsCollector
	parsedAST         interface{}
	program           *constraint.Program
	reteProgram       interface{}
	types             []interface{}
	expressions       []interface{}
	factsForRete      []map[string]interface{}
	existingFacts     []*Fact
	factsByType       map[string][]*Fact
	existingTerminals map[string]bool
	newTerminals      []*TerminalNode
	hasResets         bool
	tx                *Transaction
}

// beginIngestionTransaction démarre une transaction pour l'ingestion
func (ctx *ingestionContext) beginIngestionTransaction(cp *ConstraintPipeline) error {
	if ctx.network == nil {
		return nil
	}

	ctx.tx = ctx.network.BeginTransaction()
	ctx.network.SetTransaction(ctx.tx)
	cp.GetLogger().Info("🔒 Transaction démarrée automatiquement: %s", ctx.tx.ID)
	return nil
}

// rollbackIngestionOnError effectue un rollback en cas d'erreur
func (ctx *ingestionContext) rollbackIngestionOnError(cp *ConstraintPipeline, err error) error {
	if ctx.tx != nil && ctx.tx.IsActive {
		rollbackErr := ctx.tx.Rollback()
		if rollbackErr != nil {
			cp.GetLogger().Error("❌ Erreur rollback: %v", rollbackErr)
			return fmt.Errorf("erreur ingestion: %w; erreur rollback: %v", err, rollbackErr)
		}
		cp.GetLogger().Warn("🔙 Rollback automatique effectué")
	}
	return err
}

// commitIngestionTransaction commit la transaction après vérifications
func (ctx *ingestionContext) commitIngestionTransaction(cp *ConstraintPipeline) error {
	if ctx.tx == nil || !ctx.tx.IsActive {
		return nil
	}

	commitErr := ctx.tx.Commit()
	if commitErr != nil {
		return fmt.Errorf("❌ Erreur commit transaction: %w", commitErr)
	}
	cp.GetLogger().Info("✅ Transaction committée: %d changements", ctx.tx.GetCommandCount())
	return nil
}

// parseAndDetectReset parse le fichier et détecte les commandes reset
func (cp *ConstraintPipeline) parseAndDetectReset(ctx *ingestionContext) error {
	parsingStart := time.Now()

	parsedAST, err := constraint.ParseConstraintFile(ctx.filename)
	if err != nil {
		return fmt.Errorf("❌ Erreur parsing fichier %s: %w", ctx.filename, err)
	}
	ctx.parsedAST = parsedAST

	if ctx.metrics != nil {
		ctx.metrics.RecordParsingDuration(time.Since(parsingStart))
	}
	cp.GetLogger().Info("✅ Parsing réussi")

	// Détecter reset
	resultMap, ok := parsedAST.(map[string]interface{})
	if !ok {
		return fmt.Errorf("❌ Format AST non reconnu: %T", parsedAST)
	}

	if resetsData, exists := resultMap["resets"]; exists {
		if resets, ok := resetsData.([]interface{}); ok && len(resets) > 0 {
			ctx.hasResets = true
			cp.GetLogger().Info("🔄 Commande reset détectée - Réinitialisation complète du réseau")
		}
	}

	return nil
}

// initializeNetworkWithReset gère la réinitialisation ou création du réseau
func (cp *ConstraintPipeline) initializeNetworkWithReset(ctx *ingestionContext) error {
	if !ctx.hasResets {
		return nil
	}

	cp.GetLogger().Info("🔄 Commande reset détectée - Garbage Collection de l'ancien réseau")

	if ctx.network != nil {
		cp.GetLogger().Debug("🗑️ GC du réseau existant...")
		ctx.network.GarbageCollect()
		cp.GetLogger().Debug("✅ GC terminé")
	}

	cp.GetLogger().Info("🆕 Création d'un nouveau réseau RETE")
	ctx.network = NewReteNetwork(ctx.storage)

	if ctx.metrics != nil {
		ctx.metrics.SetWasReset(true)
	}

	return nil
}

// validateConstraintProgram effectue la validation sémantique
func (cp *ConstraintPipeline) validateConstraintProgram(ctx *ingestionContext) error {
	validationStart := time.Now()

	if ctx.network == nil || ctx.hasResets {
		// Validation standard
		err := constraint.ValidateConstraintProgram(ctx.parsedAST)
		if err != nil {
			return fmt.Errorf("❌ Erreur validation sémantique: %w", err)
		}
		cp.GetLogger().Info("✅ Validation sémantique réussie")

		if ctx.metrics != nil {
			ctx.metrics.RecordValidationDuration(time.Since(validationStart))
			ctx.metrics.SetValidationSkipped(false)
		}
	} else {
		// Validation incrémentale
		cp.GetLogger().Info("🔍 Validation sémantique incrémentale avec contexte...")
		validator := NewIncrementalValidator(ctx.network)
		err := validator.ValidateWithContext(ctx.parsedAST)
		if err != nil {
			return fmt.Errorf("❌ Erreur validation incrémentale: %w", err)
		}
		cp.GetLogger().Info("✅ Validation incrémentale réussie (%d types en contexte)", len(ctx.network.Types))

		if ctx.metrics != nil {
			ctx.metrics.RecordValidationDuration(time.Since(validationStart))
			ctx.metrics.SetValidationSkipped(false)
			ctx.metrics.SetWasIncremental(true)
		}
	}

	return nil
}

// convertToReteProgram convertit l'AST en programme RETE et extrait les composants
func (cp *ConstraintPipeline) convertToReteProgram(ctx *ingestionContext) error {
	// Convertir en programme
	program, err := constraint.ConvertResultToProgram(ctx.parsedAST)
	if err != nil {
		return fmt.Errorf("❌ Erreur conversion programme: %w", err)
	}
	ctx.program = program

	// Créer ou étendre le réseau
	if ctx.network == nil {
		cp.GetLogger().Info("🆕 Création d'un nouveau réseau RETE")
		ctx.network = NewReteNetwork(ctx.storage)
	} else if !ctx.hasResets {
		cp.GetLogger().Info("🔄 Extension du réseau RETE existant")
	}

	// Convertir au format RETE
	reteProgram, err := constraint.ConvertToReteProgram(program)
	if err != nil {
		return fmt.Errorf("❌ Erreur conversion RETE: %w", err)
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

	cp.GetLogger().Info("✅ Trouvé %d types et %d expressions dans le fichier", len(types), len(expressions))

	return nil
}

// addTypesAndActions ajoute les types et actions au réseau
func (cp *ConstraintPipeline) addTypesAndActions(ctx *ingestionContext) error {
	// Ajouter les types
	if len(ctx.types) > 0 {
		typeCreationStart := time.Now()

		err := cp.createTypeNodes(ctx.network, ctx.types, ctx.storage)
		if err != nil {
			return fmt.Errorf("❌ Erreur ajout types: %w", err)
		}
		cp.GetLogger().Info("✅ Types ajoutés/mis à jour dans le réseau")

		if ctx.metrics != nil {
			ctx.metrics.RecordTypeCreationDuration(time.Since(typeCreationStart))
			ctx.metrics.SetTypesAdded(len(ctx.types))
		}
	}

	// Extraire et stocker les actions
	reteResultMap, ok := ctx.reteProgram.(map[string]interface{})
	if !ok {
		return fmt.Errorf("❌ Format programme RETE invalide pour actions")
	}

	err := cp.extractAndStoreActions(ctx.network, reteResultMap)
	if err != nil {
		return fmt.Errorf("❌ Erreur extraction actions: %w", err)
	}

	return nil
}

// collectExistingFactsForPropagation collecte les faits existants (sauf si reset)
func (cp *ConstraintPipeline) collectExistingFactsForPropagation(ctx *ingestionContext) {
	if ctx.hasResets {
		cp.GetLogger().Debug("📊 Réseau réinitialisé - pas de faits préexistants")
		return
	}

	collectionStart := time.Now()
	ctx.existingFacts = cp.collectExistingFacts(ctx.network)
	ctx.factsByType = cp.organizeFactsByType(ctx.existingFacts)

	cp.GetLogger().Debug("📊 Faits préexistants dans le réseau: %d", len(ctx.existingFacts))

	if ctx.metrics != nil {
		ctx.metrics.RecordFactCollectionDuration(time.Since(collectionStart))
		ctx.metrics.SetExistingFactsCollected(len(ctx.existingFacts))
	}
}

// manageRules gère l'ajout et la suppression de règles
func (cp *ConstraintPipeline) manageRules(ctx *ingestionContext) error {
	// Identifier les terminaux existants
	ctx.existingTerminals = make(map[string]bool)
	for terminalID := range ctx.network.TerminalNodes {
		ctx.existingTerminals[terminalID] = true
	}

	// Ajouter les nouvelles règles
	if len(ctx.expressions) > 0 {
		ruleCreationStart := time.Now()

		err := cp.createRuleNodes(ctx.network, ctx.expressions, ctx.storage)
		if err != nil {
			return fmt.Errorf("❌ Erreur ajout règles: %w", err)
		}
		cp.GetLogger().Info("✅ Règles ajoutées au réseau")

		if ctx.metrics != nil {
			ctx.metrics.RecordRuleCreationDuration(time.Since(ruleCreationStart))
			ctx.metrics.SetRulesAdded(len(ctx.expressions))
		}
	}

	// Traiter les suppressions de règles
	reteResultMap, ok := ctx.reteProgram.(map[string]interface{})
	if !ok {
		return fmt.Errorf("❌ Format programme RETE invalide pour suppressions")
	}

	err := cp.processRuleRemovals(ctx.network, reteResultMap)
	if err != nil {
		return fmt.Errorf("❌ Erreur traitement suppressions de règles: %w", err)
	}

	return nil
}

// propagateFactsToNewRules propage les faits existants vers les nouvelles règles
func (cp *ConstraintPipeline) propagateFactsToNewRules(ctx *ingestionContext) {
	ctx.newTerminals = cp.identifyNewTerminals(ctx.network, ctx.existingTerminals)

	if len(ctx.newTerminals) == 0 || len(ctx.existingFacts) == 0 {
		return
	}

	cp.GetLogger().Info("🔄 Propagation ciblée de faits vers %d nouvelle(s) règle(s)", len(ctx.newTerminals))

	propagationStart := time.Now()
	propagatedCount := cp.propagateToNewTerminals(ctx.network, ctx.newTerminals, ctx.factsByType)

	if ctx.metrics != nil {
		ctx.metrics.RecordPropagationDuration(time.Since(propagationStart))
		ctx.metrics.SetFactsPropagated(propagatedCount)
		ctx.metrics.SetNewTerminalsAdded(len(ctx.newTerminals))
		ctx.metrics.SetPropagationTargets(len(ctx.newTerminals))
	}

	cp.GetLogger().Info("✅ Propagation rétroactive terminée (%d fait(s) propagé(s))", propagatedCount)
}

// submitFactsInternal soumet les nouveaux faits du fichier au réseau (legacy - utiliser submitNewFacts)
func (cp *ConstraintPipeline) submitFactsInternal(ctx *ingestionContext) error {
	if len(ctx.program.Facts) == 0 {
		return nil
	}

	factsForRete, err := constraint.ConvertFactsToReteFormat(*ctx.program)
	if err != nil {
		return fmt.Errorf("❌ Erreur conversion faits: %w", err)
	}
	ctx.factsForRete = factsForRete
	cp.GetLogger().Info("📥 Soumission de %d nouveaux faits", len(ctx.factsForRete))

	submissionStart := time.Now()
	err = ctx.network.SubmitFactsFromGrammar(ctx.factsForRete)
	if err != nil {
		return fmt.Errorf("❌ Erreur soumission faits: %w", err)
	}

	cp.GetLogger().Info("✅ Nouveaux faits soumis")

	if ctx.metrics != nil {
		ctx.metrics.RecordFactSubmissionDuration(time.Since(submissionStart))
		ctx.metrics.SetFactsSubmitted(len(ctx.factsForRete))
	}

	return nil
}

// validateNetworkAndCoherence effectue la validation finale et la vérification de cohérence
func (cp *ConstraintPipeline) validateNetworkAndCoherence(ctx *ingestionContext) error {
	// Validation réseau
	err := cp.validateNetwork(ctx.network)
	if err != nil {
		return fmt.Errorf("❌ Erreur validation réseau: %w", err)
	}
	cp.GetLogger().Info("✅ Validation réussie")

	// Enregistrer l'état du réseau
	if ctx.metrics != nil {
		ctx.metrics.RecordNetworkState(ctx.network)
	}

	cp.GetLogger().Info("🎯 INGESTION INCRÉMENTALE TERMINÉE")
	cp.GetLogger().Info("   - Total TypeNodes: %d", len(ctx.network.TypeNodes))
	cp.GetLogger().Info("   - Total TerminalNodes: %d", len(ctx.network.TerminalNodes))

	// Vérification de cohérence pré-commit
	if ctx.tx != nil && ctx.tx.IsActive && len(ctx.factsForRete) > 0 {
		cp.GetLogger().Info("🔍 Vérification de cohérence pré-commit...")

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
			cp.GetLogger().Error("❌ Incohérence détectée: %d faits attendus, %d trouvés", expectedFactCount, actualFactCount)
			cp.GetLogger().Error("   Faits manquants: %v", missingFacts)
			return fmt.Errorf(
				"incohérence pré-commit: %d faits attendus mais %d trouvés dans le storage",
				expectedFactCount, actualFactCount)
		}

		cp.GetLogger().Info("✅ Cohérence vérifiée: %d/%d faits présents", actualFactCount, expectedFactCount)

		// Synchroniser le storage
		cp.GetLogger().Info("💾 Synchronisation du storage...")
		if err := ctx.storage.Sync(); err != nil {
			return fmt.Errorf("❌ Erreur sync storage: %w", err)
		}
		cp.GetLogger().Info("✅ Storage synchronisé")
	}

	return nil
}
