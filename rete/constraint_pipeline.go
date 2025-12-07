// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
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
	metrics := NewMetricsCollector()
	resultNetwork, err := cp.ingestFileWithMetrics(filename, network, storage, metrics)
	finalMetrics := metrics.Finalize()
	return resultNetwork, finalMetrics, err
}

// ingestFileWithMetrics est l'implémentation interne avec support optionnel des métriques
// IMPORTANT: Gère les transactions automatiquement (TOUJOURS activées)
func (cp *ConstraintPipeline) ingestFileWithMetrics(filename string, network *ReteNetwork, storage Storage, metrics *MetricsCollector) (*ReteNetwork, error) {
	cp.logger.Info("========================================")
	cp.logger.Info("📁 Ingestion incrémentale: %s", filename)

	// Initialiser le contexte d'ingestion
	ctx := &ingestionContext{
		filename: filename,
		network:  network,
		storage:  storage,
		metrics:  metrics,
	}

	// ÉTAPE 1: Parsing et détection reset
	if err := cp.parseAndDetectReset(ctx); err != nil {
		return nil, err
	}

	// ÉTAPE 2: Initialisation réseau (GC si reset)
	if err := cp.initializeNetworkWithReset(ctx); err != nil {
		return nil, err
	}

	// ÉTAPE 3: Démarrer transaction
	if err := ctx.beginIngestionTransaction(cp); err != nil {
		return nil, err
	}

	// Wrapper pour rollback automatique en cas d'erreur
	handleError := func(err error) (*ReteNetwork, error) {
		return ctx.network, ctx.rollbackIngestionOnError(cp, err)
	}

	// ÉTAPE 4: Validation sémantique
	if err := cp.validateConstraintProgram(ctx); err != nil {
		return handleError(err)
	}

	// ÉTAPE 5: Conversion en programme RETE
	if err := cp.convertToReteProgram(ctx); err != nil {
		return handleError(err)
	}

	// ÉTAPE 6: Ajout types et actions
	if err := cp.addTypesAndActions(ctx); err != nil {
		return handleError(err)
	}

	// ÉTAPE 7: Collection faits existants
	cp.collectExistingFactsForPropagation(ctx)

	// ÉTAPE 8: Gestion des règles (ajout + suppression)
	if err := cp.manageRules(ctx); err != nil {
		return handleError(err)
	}

	// ÉTAPE 9: Propagation rétroactive vers nouvelles règles
	cp.propagateFactsToNewRules(ctx)

	// ÉTAPE 10: Soumission nouveaux faits
	if err := cp.submitNewFacts(ctx); err != nil {
		return handleError(err)
	}

	// ÉTAPE 11: Validation finale et cohérence
	if err := cp.validateNetworkAndCoherence(ctx); err != nil {
		return handleError(err)
	}

	// ÉTAPE 12: Commit transaction
	if err := ctx.commitIngestionTransaction(cp); err != nil {
		return handleError(err)
	}

	cp.logger.Info("🎯 INGESTION TERMINÉE")
	cp.logger.Info("========================================")

	return ctx.network, nil
}
