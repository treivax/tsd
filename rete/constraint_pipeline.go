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

// IngestFileWithMetrics est un wrapper qui collecte les métriques
// IngestFileWithMetrics ingère un fichier avec collecte de métriques
// IMPORTANT: Cette fonction utilise TOUJOURS les transactions avec auto-commit/auto-rollback.
// En cas d'erreur, la transaction est automatiquement annulée (rollback).
// En cas de succès, la transaction est automatiquement validée (commit).
func (cp *ConstraintPipeline) IngestFileWithMetrics(filename string, network *ReteNetwork, storage Storage) (*ReteNetwork, *IngestionMetrics, error) {
	metrics := NewMetricsCollector()
	resultNetwork, err := cp.ingestFileWithMetrics(filename, network, storage, metrics)
	finalMetrics := metrics.Finalize()
	return resultNetwork, finalMetrics, err
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
func (cp *ConstraintPipeline) IngestFile(filename string, network *ReteNetwork, storage Storage) (*ReteNetwork, error) {
	return cp.ingestFileWithMetrics(filename, network, storage, nil)
}

// ingestFileWithMetrics est l'implémentation interne avec support optionnel des métriques
// IMPORTANT: Gère les transactions automatiquement (TOUJOURS activées)
func (cp *ConstraintPipeline) ingestFileWithMetrics(filename string, network *ReteNetwork, storage Storage, metrics *MetricsCollector) (*ReteNetwork, error) {
	cp.logger.Info("========================================")
	cp.logger.Info("📁 Ingestion incrémentale: %s", filename)

	// ÉTAPE 1: Parsing du fichier
	parsingStart := time.Now()
	parsedAST, err := constraint.ParseConstraintFile(filename)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur parsing fichier %s: %w", filename, err)
	}
	if metrics != nil {
		metrics.RecordParsingDuration(time.Since(parsingStart))
	}
	cp.logger.Info("✅ Parsing réussi")

	// ÉTAPE 2: Vérifier la présence d'une commande reset
	resultMap, ok := parsedAST.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ Format AST non reconnu: %T", parsedAST)
	}

	hasResets := false
	if resetsData, exists := resultMap["resets"]; exists {
		if resets, ok := resetsData.([]interface{}); ok && len(resets) > 0 {
			hasResets = true
			cp.GetLogger().Info("🔄 Commande reset détectée - Réinitialisation complète du réseau")
		}
	}

	// Si reset détecté, faire un GC de l'ancien réseau puis créer un nouveau
	if hasResets {
		cp.GetLogger().Info("🔄 Commande reset détectée - Garbage Collection de l'ancien réseau")

		// OPTIMISATION 2: Garbage Collection automatique après reset
		if network != nil {
			cp.GetLogger().Debug("🗑️ GC du réseau existant...")
			network.GarbageCollect()
			cp.GetLogger().Debug("✅ GC terminé")
		}

		cp.GetLogger().Info("🆕 Création d'un nouveau réseau RETE")
		network = NewReteNetwork(storage)
		if metrics != nil {
			metrics.SetWasReset(true)
		}
	}

	// ÉTAPE 2.5: Démarrer une transaction (OBLIGATOIRE) une fois que le réseau est défini
	var tx *Transaction
	if network != nil {
		tx = network.BeginTransaction()
		network.SetTransaction(tx)
		cp.GetLogger().Info("🔒 Transaction démarrée automatiquement: %s", tx.ID)
	}

	// Fonction de rollback en cas d'erreur
	rollbackOnError := func(err error) (*ReteNetwork, error) {
		if tx != nil && tx.IsActive {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				cp.GetLogger().Error("❌ Erreur rollback: %v", rollbackErr)
				return network, fmt.Errorf("erreur ingestion: %w; erreur rollback: %v", err, rollbackErr)
			}
			cp.GetLogger().Warn("🔙 Rollback automatique effectué")
		}
		return network, err
	}

	// ÉTAPE 3: Validation sémantique
	// OPTIMISATION 1: Validation incrémentale avec contexte (systématiquement activée)
	validationStart := time.Now()
	if network == nil || hasResets {
		// Validation standard pour la création initiale ou après reset
		err = constraint.ValidateConstraintProgram(parsedAST)
		if err != nil {
			return rollbackOnError(fmt.Errorf("❌ Erreur validation sémantique: %w", err))
		}
		cp.GetLogger().Info("✅ Validation sémantique réussie")
		if metrics != nil {
			metrics.RecordValidationDuration(time.Since(validationStart))
			metrics.SetValidationSkipped(false)
		}
	} else {
		// Validation incrémentale avec contexte du réseau existant
		cp.GetLogger().Info("🔍 Validation sémantique incrémentale avec contexte...")
		validator := NewIncrementalValidator(network)
		err = validator.ValidateWithContext(parsedAST)
		if err != nil {
			return rollbackOnError(fmt.Errorf("❌ Erreur validation incrémentale: %w", err))
		}
		cp.GetLogger().Info("✅ Validation incrémentale réussie (%d types en contexte)", len(network.Types))
		if metrics != nil {
			metrics.RecordValidationDuration(time.Since(validationStart))
			metrics.SetValidationSkipped(false)
			metrics.SetWasIncremental(true)
		}
	}

	// Convertir en programme
	program, err := constraint.ConvertResultToProgram(parsedAST)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur conversion programme: %w", err)
	}

	// ÉTAPE 4: Créer ou étendre le réseau
	if network == nil {
		cp.GetLogger().Info("🆕 Création d'un nouveau réseau RETE")
		network = NewReteNetwork(storage)
	} else if !hasResets {
		cp.GetLogger().Info("🔄 Extension du réseau RETE existant")
	}

	// Convertir au format RETE
	reteProgram := constraint.ConvertToReteProgram(program)
	reteResultMap, ok := reteProgram.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ Format programme RETE invalide: %T", reteProgram)
	}

	// ÉTAPE 5: Extraire et ajouter les nouveaux types
	types, expressions, err := cp.extractComponents(reteResultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur extraction composants: %w", err)
	}
	cp.GetLogger().Info("✅ Trouvé %d types et %d expressions dans le fichier", len(types), len(expressions))

	// Ajouter les types au réseau (évite les doublons automatiquement)
	typeCreationStart := time.Now()
	if len(types) > 0 {
		err = cp.createTypeNodes(network, types, storage)
		if err != nil {
			return nil, fmt.Errorf("❌ Erreur ajout types: %w", err)
		}
		cp.GetLogger().Info("✅ Types ajoutés/mis à jour dans le réseau")
		if metrics != nil {
			metrics.RecordTypeCreationDuration(time.Since(typeCreationStart))
			metrics.SetTypesAdded(len(types))
		}
	}

	// ÉTAPE 5.5: Extraire et stocker les définitions d'actions
	err = cp.extractAndStoreActions(network, reteResultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur extraction actions: %w", err)
	}

	// ÉTAPE 6: Collecter tous les faits existants dans le réseau AVANT d'ajouter les nouvelles règles
	// (sauf si reset car le réseau vient d'être créé vide)
	var existingFacts []*Fact
	var existingFactsByType map[string][]*Fact
	collectionStart := time.Now()
	if !hasResets {
		existingFacts = cp.collectExistingFacts(network)
		existingFactsByType = cp.organizeFactsByType(existingFacts)
		cp.GetLogger().Debug("📊 Faits préexistants dans le réseau: %d", len(existingFacts))
		if metrics != nil {
			metrics.RecordFactCollectionDuration(time.Since(collectionStart))
			metrics.SetExistingFactsCollected(len(existingFacts))
		}
	} else {
		cp.GetLogger().Debug("📊 Réseau réinitialisé - pas de faits préexistants")
	}

	// ÉTAPE 7: Identifier les terminaux existants avant l'ajout de règles
	existingTerminals := make(map[string]bool)
	for terminalID := range network.TerminalNodes {
		existingTerminals[terminalID] = true
	}

	// ÉTAPE 8: Ajouter les nouvelles règles
	ruleCreationStart := time.Now()
	if len(expressions) > 0 {
		err = cp.createRuleNodes(network, expressions, storage)
		if err != nil {
			return nil, fmt.Errorf("❌ Erreur ajout règles: %w", err)
		}
		cp.GetLogger().Info("✅ Règles ajoutées au réseau")
		if metrics != nil {
			metrics.RecordRuleCreationDuration(time.Since(ruleCreationStart))
			metrics.SetRulesAdded(len(expressions))
		}
	}

	// ÉTAPE 9: Traiter les suppressions de règles (si présentes)
	err = cp.processRuleRemovals(network, reteResultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur traitement suppressions de règles: %w", err)
	}

	// ÉTAPE 10: Propager les faits existants vers les nouvelles règles uniquement
	newTerminals := cp.identifyNewTerminals(network, existingTerminals)

	if len(newTerminals) > 0 && len(existingFacts) > 0 {
		cp.GetLogger().Info("🔄 Propagation ciblée de faits vers %d nouvelle(s) règle(s)", len(newTerminals))

		// Propager de manière ciblée pour chaque nouveau terminal
		propagationStart := time.Now()
		propagatedCount := cp.propagateToNewTerminals(network, newTerminals, existingFactsByType)

		if metrics != nil {
			metrics.RecordPropagationDuration(time.Since(propagationStart))
			metrics.SetFactsPropagated(propagatedCount)
			metrics.SetNewTerminalsAdded(len(newTerminals))
			metrics.SetPropagationTargets(len(newTerminals))
		}

		cp.GetLogger().Info("✅ Propagation rétroactive terminée (%d fait(s) propagé(s))", propagatedCount)
	}

	// ÉTAPE 10: Soumettre les nouveaux faits du fichier
	var factsForRete []map[string]interface{}
	if len(program.Facts) > 0 {
		factsForRete = constraint.ConvertFactsToReteFormat(*program)
		cp.GetLogger().Info("📥 Soumission de %d nouveaux faits", len(factsForRete))

		submissionStart := time.Now()
		err := network.SubmitFactsFromGrammar(factsForRete)
		if err != nil {
			return rollbackOnError(fmt.Errorf("❌ Erreur soumission faits: %w", err))
		}
		cp.GetLogger().Info("✅ Nouveaux faits soumis")
		if metrics != nil {
			metrics.RecordFactSubmissionDuration(time.Since(submissionStart))
			metrics.SetFactsSubmitted(len(factsForRete))
		}
	}

	// ÉTAPE 11: Validation finale
	err = cp.validateNetwork(network)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur validation réseau: %w", err)
	}
	cp.GetLogger().Info("✅ Validation réussie")

	// Enregistrer l'état final du réseau dans les métriques
	if metrics != nil {
		metrics.RecordNetworkState(network)
	}

	cp.GetLogger().Info("🎯 INGESTION INCRÉMENTALE TERMINÉE")
	cp.GetLogger().Info("   - Total TypeNodes: %d", len(network.TypeNodes))
	cp.GetLogger().Info("   - Total TerminalNodes: %d", len(network.TerminalNodes))

	// ÉTAPE 12: Vérification de cohérence avant commit
	if tx != nil && tx.IsActive && len(factsForRete) > 0 {
		cp.GetLogger().Info("🔍 Vérification de cohérence pré-commit...")

		// Vérifier que tous les faits soumis sont bien dans le storage
		expectedFactCount := len(factsForRete)
		actualFactCount := 0
		missingFacts := make([]string, 0)

		for i, factMap := range factsForRete {
			var factID string
			if id, ok := factMap["id"].(string); ok {
				factID = id
			} else {
				// Générer le même ID que dans SubmitFactsFromGrammar
				factID = fmt.Sprintf("fact_%d", i)
			}

			// Extraire le type du fait
			factType := "unknown"
			if typ, ok := factMap["type"].(string); ok {
				factType = typ
			} else if typ, ok := factMap["reteType"].(string); ok {
				factType = typ
			}

			// Construire l'ID interne (Type_ID) comme dans GetInternalID()
			internalID := fmt.Sprintf("%s_%s", factType, factID)

			if storage.GetFact(internalID) != nil {
				actualFactCount++
			} else {
				missingFacts = append(missingFacts, internalID)
			}
		}

		if expectedFactCount != actualFactCount {
			cp.GetLogger().Error("❌ Incohérence détectée: %d faits attendus, %d trouvés", expectedFactCount, actualFactCount)
			cp.GetLogger().Error("   Faits manquants: %v", missingFacts)
			return rollbackOnError(fmt.Errorf(
				"incohérence pré-commit: %d faits attendus mais %d trouvés dans le storage",
				expectedFactCount, actualFactCount))
		}

		cp.GetLogger().Info("✅ Cohérence vérifiée: %d/%d faits présents", actualFactCount, expectedFactCount)

		// Synchroniser le storage pour garantir la durabilité
		cp.GetLogger().Info("💾 Synchronisation du storage...")
		if err := storage.Sync(); err != nil {
			return rollbackOnError(fmt.Errorf("❌ Erreur sync storage: %w", err))
		}
		cp.GetLogger().Info("✅ Storage synchronisé")
	}

	// ÉTAPE 13: Commit de la transaction (OBLIGATOIRE)
	if tx != nil && tx.IsActive {
		commitErr := tx.Commit()
		if commitErr != nil {
			return rollbackOnError(fmt.Errorf("❌ Erreur commit transaction: %w", commitErr))
		}
		cp.GetLogger().Info("✅ Transaction committée: %d changements", tx.GetCommandCount())
	}

	cp.GetLogger().Info("🎯 INGESTION TERMINÉE")
	cp.GetLogger().Info("========================================")

	return network, nil
}
