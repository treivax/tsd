// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"
)

// ReteNetwork représente le réseau RETE complet
type ReteNetwork struct {
	RootNode              *RootNode                `json:"root_node"`
	TypeNodes             map[string]*TypeNode     `json:"type_nodes"`
	AlphaNodes            map[string]*AlphaNode    `json:"alpha_nodes"`
	BetaNodes             map[string]interface{}   `json:"beta_nodes"` // Nœuds Beta pour les jointures multi-faits
	TerminalNodes         map[string]*TerminalNode `json:"terminal_nodes"`
	Storage               Storage                  `json:"-"`
	Types                 []TypeDefinition         `json:"types"`
	Actions               []ActionDefinition       `json:"actions"` // Action definitions for incremental validation
	BetaBuilder           interface{}              `json:"-"`       // Constructeur de réseau Beta (deprecated, use BetaChainBuilder)
	LifecycleManager      *LifecycleManager        `json:"-"`       // Gestionnaire du cycle de vie des nœuds
	AlphaSharingManager   *AlphaSharingRegistry    `json:"-"`       // Gestionnaire du partage des AlphaNodes
	AlphaChainBuilder     *AlphaChainBuilder       `json:"-"`       // Constructeur de chaînes alpha avec décomposition
	PassthroughRegistry   map[string]*AlphaNode    `json:"-"`       // Registre de partage des AlphaNodes passthrough
	BetaSharingRegistry   BetaSharingRegistry      `json:"-"`       // Gestionnaire du partage des JoinNodes
	BetaChainBuilder      *BetaChainBuilder        `json:"-"`       // Constructeur de chaînes beta avec partage
	ChainMetrics          *ChainBuildMetrics       `json:"-"`       // Métriques de performance pour la construction des chaînes
	Config                *ChainPerformanceConfig  `json:"-"`       // Configuration de performance
	ActionExecutor        *ActionExecutor          `json:"-"`       // Exécuteur d'actions
	ArithmeticResultCache *ArithmeticResultCache   `json:"-"`       // Cache global des résultats arithmétiques intermédiaires
	currentTx             *Transaction             `json:"-"`       // Transaction courante (si en cours)
	txMutex               sync.RWMutex             `json:"-"`       // Mutex pour accès concurrent à la transaction
	logger                *Logger                  `json:"-"`       // Logger structuré pour instrumentation
	factIDCounter         int64                    `json:"-"`       // Compteur thread-safe pour génération d'IDs de faits
	factIDMutex           sync.Mutex               `json:"-"`       // Mutex pour génération d'IDs de faits

	// Phase 2: Configuration de synchronisation pour garanties de cohérence
	SubmissionTimeout time.Duration `json:"-"` // Timeout global pour soumission de faits
	VerifyRetryDelay  time.Duration `json:"-"` // Délai entre tentatives de vérification
	MaxVerifyRetries  int           `json:"-"` // Nombre max de tentatives de vérification
}

// Valeurs par défaut pour la synchronisation Phase 2
const (
	DefaultSubmissionTimeout = 30 * time.Second
	DefaultVerifyRetryDelay  = 10 * time.Millisecond
	DefaultMaxVerifyRetries  = 10
)

// GetChainMetrics retourne les métriques de performance pour la construction des chaînes alpha
func (rn *ReteNetwork) GetChainMetrics() *ChainBuildMetrics {
	if rn.ChainMetrics == nil {
		rn.ChainMetrics = NewChainBuildMetrics()
	}
	return rn.ChainMetrics
}

// GetBetaSharingStats retourne les statistiques de partage des JoinNodes
func (rn *ReteNetwork) GetBetaSharingStats() *BetaSharingStats {
	if rn.BetaSharingRegistry == nil {
		return nil
	}
	return rn.BetaSharingRegistry.GetSharingStats()
}

// GetBetaChainMetrics retourne les métriques de construction des chaînes beta
func (rn *ReteNetwork) GetBetaChainMetrics() *BetaChainMetrics {
	if rn.BetaChainBuilder == nil {
		return nil
	}
	return rn.BetaChainBuilder.GetMetrics()
}

// GetConfig retourne la configuration de performance
func (rn *ReteNetwork) GetConfig() *ChainPerformanceConfig {
	if rn.Config == nil {
		rn.Config = DefaultChainPerformanceConfig()
	}
	return rn.Config
}

// ResetChainMetrics réinitialise toutes les métriques de performance
func (rn *ReteNetwork) ResetChainMetrics() {
	if rn.ChainMetrics != nil {
		rn.ChainMetrics.Reset()
	}
	if rn.BetaChainBuilder != nil {
		rn.BetaChainBuilder.ResetMetrics()
	}
}

// SetLogger configure le logger pour le réseau RETE
func (rn *ReteNetwork) SetLogger(logger *Logger) {
	if logger != nil {
		rn.logger = logger
	}
}

// GetLogger retourne le logger actuel du réseau
func (rn *ReteNetwork) GetLogger() *Logger {
	if rn.logger == nil {
		rn.logger = NewLogger(LogLevelInfo, nil)
	}
	return rn.logger
}

// SetTransaction active une transaction pour toutes les opérations suivantes
func (rn *ReteNetwork) SetTransaction(tx *Transaction) {
	rn.txMutex.Lock()
	defer rn.txMutex.Unlock()
	rn.currentTx = tx
}

// GetTransaction retourne la transaction courante (ou nil)
func (rn *ReteNetwork) GetTransaction() *Transaction {
	rn.txMutex.RLock()
	defer rn.txMutex.RUnlock()
	return rn.currentTx
}

// GetTypeDefinition retourne la définition d'un type par son nom
func (rn *ReteNetwork) GetTypeDefinition(typeName string) *TypeDefinition {
	for i := range rn.Types {
		if rn.Types[i].Name == typeName {
			return &rn.Types[i]
		}
	}
	return nil
}

// GetRuleInfo retourne les informations d'une règle
func (rn *ReteNetwork) GetRuleInfo(ruleID string) *RuleInfo {
	if rn.LifecycleManager == nil {
		return &RuleInfo{
			RuleID:    ruleID,
			NodeIDs:   []string{},
			NodeCount: 0,
		}
	}
	return rn.LifecycleManager.GetRuleInfo(ruleID)
}

// GetNetworkStats retourne des statistiques sur le réseau
func (rn *ReteNetwork) GetNetworkStats() map[string]interface{} {
	stats := map[string]interface{}{
		"type_nodes":     len(rn.TypeNodes),
		"alpha_nodes":    len(rn.AlphaNodes),
		"beta_nodes":     len(rn.BetaNodes),
		"terminal_nodes": len(rn.TerminalNodes),
	}

	if rn.LifecycleManager != nil {
		lifecycleStats := rn.LifecycleManager.GetStats()
		for k, v := range lifecycleStats {
			stats["lifecycle_"+k] = v
		}
	}

	if rn.AlphaSharingManager != nil {
		alphaStats := rn.AlphaSharingManager.GetStats()
		for k, v := range alphaStats {
			stats["sharing_"+k] = v
		}
	}

	return stats
}

// GenerateFactID génère un ID unique thread-safe pour un nouveau fait.
//
// Cette méthode utilise un compteur atomique pour garantir l'unicité
// même en cas d'exécution concurrente.
//
// Paramètres :
//   - typeName : nom du type de fait
//
// Retourne :
//   - string : ID unique au format "typeName_N"
func (rn *ReteNetwork) GenerateFactID(typeName string) string {
	rn.factIDMutex.Lock()
	defer rn.factIDMutex.Unlock()

	rn.factIDCounter++
	return fmt.Sprintf("%s_%d", typeName, rn.factIDCounter)
}
