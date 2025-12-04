// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/constraint"
)

// AdvancedPipelineConfig configure les optimisations avancées du pipeline
// Note: La validation incrémentale, le GC et les TRANSACTIONS sont TOUJOURS activés (non configurables)
type AdvancedPipelineConfig struct {
	// Transactions (toujours activées)
	TransactionTimeout  time.Duration
	MaxTransactionSize  int64
	AutoCommit          bool // Commit automatique si pas d'erreur
	AutoRollbackOnError bool // Rollback automatique en cas d'erreur
}

// DefaultAdvancedPipelineConfig retourne la configuration par défaut
// Note: La validation incrémentale, le GC et les transactions sont toujours activés
func DefaultAdvancedPipelineConfig() *AdvancedPipelineConfig {
	return &AdvancedPipelineConfig{
		TransactionTimeout:  30 * time.Second,
		MaxTransactionSize:  100 * 1024 * 1024, // 100 MB
		AutoCommit:          false,
		AutoRollbackOnError: true,
	}
}

// AdvancedMetrics contient les métriques des optimisations avancées
type AdvancedMetrics struct {
	// Validation (toujours activée)
	ValidationWithContextDuration time.Duration
	TypesFoundInContext           int
	ValidationErrors              []string

	// Garbage Collection (toujours activée)
	GCDuration     time.Duration
	NodesCollected int
	MemoryFreed    int64
	GCPerformed    bool

	// Transaction (toujours utilisée)
	TransactionID        string
	TransactionFootprint int64
	ChangesTracked       int
	RollbackPerformed    bool
	RollbackDuration     time.Duration
	TransactionDuration  time.Duration
}

// IngestFileWithAdvancedFeatures ingère un fichier avec toutes les optimisations avancées
// Cette fonction combine validation incrémentale, GC et transactions (toujours activées)
func (cp *ConstraintPipeline) IngestFileWithAdvancedFeatures(
	filename string,
	network *ReteNetwork,
	storage Storage,
	config *AdvancedPipelineConfig,
) (*ReteNetwork, *AdvancedMetrics, error) {
	if config == nil {
		config = DefaultAdvancedPipelineConfig()
	}

	metrics := &AdvancedMetrics{}

	// Create network if nil
	if network == nil {
		network = NewReteNetwork(storage)
	}

	// Phase 1: Démarrer une transaction (OBLIGATOIRE)
	txStart := time.Now()
	tx := network.BeginTransaction()
	network.SetTransaction(tx)
	metrics.TransactionID = tx.ID
	metrics.TransactionFootprint = tx.GetMemoryFootprint()
	cp.GetLogger().Info("🔒 Transaction démarrée: %s (footprint: %.2f KB)",
		tx.ID, float64(metrics.TransactionFootprint)/1024)

	// Vérifier la taille de l'empreinte mémoire
	if config.MaxTransactionSize > 0 && metrics.TransactionFootprint > config.MaxTransactionSize {
		return nil, metrics, fmt.Errorf(
			"transaction trop volumineuse: %d bytes (max: %d)",
			metrics.TransactionFootprint,
			config.MaxTransactionSize,
		)
	}

	defer func() {
		metrics.TransactionDuration = time.Since(txStart)
	}()

	// Phase 2: Parser le fichier
	cp.GetLogger().Info("========================================")
	cp.GetLogger().Info("📁 Ingestion avancée: %s", filename)

	parsedAST, err := parseFile(filename)
	if err != nil {
		if config.AutoRollbackOnError {
			rollbackStart := time.Now()
			tx.Rollback()
			metrics.RollbackPerformed = true
			metrics.RollbackDuration = time.Since(rollbackStart)
			cp.GetLogger().Warn("🔙 Rollback effectué en %v", metrics.RollbackDuration)
		}
		return nil, metrics, fmt.Errorf("❌ Erreur parsing: %w", err)
	}

	// Phase 3: Vérifier reset et effectuer GC (toujours activé)
	hasReset := detectReset(parsedAST)
	if hasReset {
		cp.GetLogger().Info("🔄 Commande reset détectée")

		// Effectuer GC de l'ancien réseau (toujours activé)
		if network != nil {
			gcStart := time.Now()
			cp.GetLogger().Info("🗑️  Garbage Collection de l'ancien réseau...")

			// Compter les nœuds avant GC
			nodesBefore := len(network.TypeNodes) + len(network.AlphaNodes) +
				len(network.BetaNodes) + len(network.TerminalNodes)

			network.GarbageCollect()

			metrics.GCPerformed = true
			metrics.NodesCollected = nodesBefore
			metrics.GCDuration = time.Since(gcStart)

			cp.GetLogger().Info("✅ GC terminé: %d nœuds collectés en %v",
				metrics.NodesCollected, metrics.GCDuration)
		}

		// Créer nouveau réseau
		network = NewReteNetwork(storage)
		cp.GetLogger().Info("🆕 Nouveau réseau RETE créé")
	}

	// Phase 4: Validation sémantique (toujours activée)
	// Validation incrémentale avec contexte ou validation standard selon le cas
	validationStart := time.Now()
	if network != nil && !hasReset {
		// Validation incrémentale avec contexte (mode incrémental)
		cp.GetLogger().Info("🔍 Validation incrémentale avec contexte...")

		validator := NewIncrementalValidator(network)
		err = validator.ValidateWithContext(parsedAST)

		metrics.TypesFoundInContext = len(network.Types)
		metrics.ValidationWithContextDuration = time.Since(validationStart)

		if err != nil {
			metrics.ValidationErrors = append(metrics.ValidationErrors, err.Error())
			cp.GetLogger().Error("❌ Validation incrémentale échouée: %v", err)

			if config.AutoRollbackOnError {
				rollbackStart := time.Now()
				tx.Rollback()
				metrics.RollbackPerformed = true
				metrics.RollbackDuration = time.Since(rollbackStart)
				cp.GetLogger().Warn("🔙 Rollback effectué en %v", metrics.RollbackDuration)
			}

			return nil, metrics, fmt.Errorf("❌ Validation incrémentale: %w", err)
		}

		cp.GetLogger().Info("✅ Validation incrémentale réussie (%d types en contexte)",
			metrics.TypesFoundInContext)
	} else {
		// Validation standard (création initiale ou après reset)
		err = validateStandard(parsedAST, network, hasReset)
		if err != nil {
			if config.AutoRollbackOnError {
				tx.Rollback()
				metrics.RollbackPerformed = true
			}
			return nil, metrics, fmt.Errorf("❌ Validation standard: %w", err)
		}
	}

	// Phase 5: Ingestion normale via le pipeline standard
	ingestionNetwork, ingestionErr := cp.IngestFile(filename, network, storage)

	if ingestionErr != nil {
		cp.GetLogger().Error("❌ Erreur lors de l'ingestion: %v", ingestionErr)

		// Rollback automatique si configuré
		if config.AutoRollbackOnError {
			rollbackStart := time.Now()
			rollbackErr := tx.Rollback()
			metrics.RollbackPerformed = true
			metrics.RollbackDuration = time.Since(rollbackStart)

			if rollbackErr != nil {
				cp.GetLogger().Error("❌ Erreur rollback: %v", rollbackErr)
				return nil, metrics, fmt.Errorf(
					"erreur ingestion: %w; erreur rollback: %v",
					ingestionErr, rollbackErr,
				)
			}

			cp.GetLogger().Warn("🔙 Rollback effectué avec succès en %v", metrics.RollbackDuration)
		}

		return nil, metrics, ingestionErr
	}

	network = ingestionNetwork

	// Phase 6: Commit de la transaction
	metrics.ChangesTracked = tx.GetCommandCount()

	if config.AutoCommit {
		commitErr := tx.Commit()
		if commitErr != nil {
			cp.GetLogger().Error("❌ Erreur commit: %v", commitErr)
			return nil, metrics, fmt.Errorf("erreur commit: %w", commitErr)
		}
		cp.GetLogger().Info("✅ Transaction committée: %d changements", metrics.ChangesTracked)
	} else {
		cp.GetLogger().Info("⏸️  Transaction active, commit manuel requis")
	}

	cp.GetLogger().Info("🎯 INGESTION AVANCÉE TERMINÉE")
	cp.GetLogger().Info("========================================")

	return network, metrics, nil
}

// IngestFileTransactionalSafe ingère un fichier dans une transaction avec auto-rollback
// Note: Les transactions sont maintenant TOUJOURS utilisées
func (cp *ConstraintPipeline) IngestFileTransactionalSafe(
	filename string,
	network *ReteNetwork,
	storage Storage,
) (*ReteNetwork, *Transaction, error) {
	config := DefaultAdvancedPipelineConfig()
	config.AutoCommit = false
	config.AutoRollbackOnError = true

	resultNetwork, _, err := cp.IngestFileWithAdvancedFeatures(filename, network, storage, config)

	// Récupérer la transaction depuis le réseau (elle est toujours créée)
	tx := network.GetTransaction()

	return resultNetwork, tx, err
}

// Fonctions utilitaires internes

func parseFile(filename string) (interface{}, error) {
	// Import du package constraint pour le parsing
	return constraint.ParseConstraintFile(filename)
}

func detectReset(parsedAST interface{}) bool {
	resultMap, ok := parsedAST.(map[string]interface{})
	if !ok {
		return false
	}

	if resetsData, exists := resultMap["resets"]; exists {
		if resets, ok := resetsData.([]interface{}); ok && len(resets) > 0 {
			return true
		}
	}

	return false
}

func validateStandard(parsedAST interface{}, network *ReteNetwork, hasReset bool) error {
	// Import du package constraint pour la validation
	if network == nil || hasReset {
		return constraint.ValidateConstraintProgram(parsedAST)
	}
	// En mode incrémental sans validation avancée, on ne valide pas
	return nil
}

// PrintAdvancedMetrics affiche les métriques avancées de manière formatée
func PrintAdvancedMetrics(metrics *AdvancedMetrics) {
	if metrics == nil {
		return
	}

	fmt.Println("\n📊 MÉTRIQUES AVANCÉES")
	fmt.Println("═══════════════════════════════════════")

	// Validation (toujours activée)
	if metrics.ValidationWithContextDuration > 0 {
		fmt.Println("🔍 Validation incrémentale")
		fmt.Printf("   Durée: %v\n", metrics.ValidationWithContextDuration)
		fmt.Printf("   Types en contexte: %d\n", metrics.TypesFoundInContext)
		if len(metrics.ValidationErrors) > 0 {
			fmt.Printf("   Erreurs: %d\n", len(metrics.ValidationErrors))
		}
	}

	// Garbage Collection
	if metrics.GCPerformed {
		fmt.Println("\n🗑️  Garbage Collection")
		fmt.Printf("   Durée: %v\n", metrics.GCDuration)
		fmt.Printf("   Nœuds collectés: %d\n", metrics.NodesCollected)
		if metrics.MemoryFreed > 0 {
			fmt.Printf("   Mémoire libérée: %.2f MB\n", float64(metrics.MemoryFreed)/(1024*1024))
		}
	}

	// Transaction (toujours active)
	fmt.Println("\n🔒 Transaction")
	fmt.Printf("   ID: %s\n", metrics.TransactionID)
	fmt.Printf("   Durée: %v\n", metrics.TransactionDuration)
	fmt.Printf("   Empreinte mémoire: %.2f KB\n", float64(metrics.TransactionFootprint)/1024)
	fmt.Printf("   Changements trackés: %d\n", metrics.ChangesTracked)
	if metrics.RollbackPerformed {
		fmt.Printf("   ⚠️  Rollback effectué en %v\n", metrics.RollbackDuration)
	}

	fmt.Println("═══════════════════════════════════════")
}

// GetAdvancedMetricsSummary retourne un résumé textuel des métriques
func GetAdvancedMetricsSummary(metrics *AdvancedMetrics) string {
	if metrics == nil {
		return "Pas de métriques disponibles"
	}

	summary := "Métriques avancées:\n"

	if metrics.ValidationWithContextDuration > 0 {
		summary += fmt.Sprintf("- Validation incrémentale: %v (%d types)\n",
			metrics.ValidationWithContextDuration, metrics.TypesFoundInContext)
	}

	if metrics.GCPerformed {
		summary += fmt.Sprintf("- GC: %v (%d nœuds)\n",
			metrics.GCDuration, metrics.NodesCollected)
	}

	// Transaction (toujours active)
	status := "committée"
	if metrics.RollbackPerformed {
		status = "rolled back"
	}
	summary += fmt.Sprintf("- Transaction: %v (%s, %d changements)\n",
		metrics.TransactionDuration, status, metrics.ChangesTracked)

	return summary
}
