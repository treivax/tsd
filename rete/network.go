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
	actionObserver        ActionObserver           `json:"-"`       // Observateur d'actions (nouveau)
	ArithmeticResultCache *ArithmeticResultCache   `json:"-"`       // Cache global des résultats arithmétiques intermédiaires
	currentTx             *Transaction             `json:"-"`       // Transaction courante (si en cours)
	txMutex               sync.RWMutex             `json:"-"`       // Mutex pour accès concurrent à la transaction
	logger                *Logger                  `json:"-"`       // Logger structuré pour instrumentation
	factIDCounter         int64                    `json:"-"`       // Compteur thread-safe pour génération d'IDs de faits
	factIDMutex           sync.Mutex               `json:"-"`       // Mutex pour génération d'IDs de faits
	XupleManager          interface{}              `json:"-"`       // Gestionnaire de xuples (xuples.XupleManager)
	xupleHandlerFunc      XupleHandlerFunc         `json:"-"`       // Fonction handler pour l'action Xuple
	xupleSpaceDefinitions []interface{}            `json:"-"`       // Définitions des xuple-spaces parsées depuis TSD

	// Phase 2: Configuration de synchronisation pour garanties de cohérence
	SubmissionTimeout time.Duration `json:"-"` // Timeout global pour soumission de faits
	VerifyRetryDelay  time.Duration `json:"-"` // Délai entre tentatives de vérification
	MaxVerifyRetries  int           `json:"-"` // Nombre max de tentatives de vérification

	// Contexte de soumission en cours (pour tracking des rétractations)
	currentSubmission *SubmissionContext `json:"-"` // Contexte de la soumission active (nil si aucune)
	submissionMutex   sync.RWMutex       `json:"-"` // Mutex pour accès au contexte de soumission
}

// Valeurs par défaut pour la synchronisation Phase 2
const (
	DefaultSubmissionTimeout = 30 * time.Second
	DefaultVerifyRetryDelay  = 10 * time.Millisecond
	DefaultMaxVerifyRetries  = 10
)

// SubmissionContext représente le contexte d'une soumission de faits en cours.
// Il permet de tracker quels faits ont été soumis et rétractés pendant la même transaction.
type SubmissionContext struct {
	factsSubmitted map[string]bool // Faits soumis dans cette soumission
	factsRetracted map[string]bool // Faits rétractés pendant la propagation
	mutex          sync.RWMutex    // Mutex pour accès concurrent
}

// NewSubmissionContext crée un nouveau contexte de soumission
func NewSubmissionContext() *SubmissionContext {
	return &SubmissionContext{
		factsSubmitted: make(map[string]bool),
		factsRetracted: make(map[string]bool),
	}
}

// MarkSubmitted marque un fait comme soumis
func (sc *SubmissionContext) MarkSubmitted(factID string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.factsSubmitted[factID] = true
}

// MarkRetracted marque un fait comme rétracté pendant la propagation
func (sc *SubmissionContext) MarkRetracted(factID string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	sc.factsRetracted[factID] = true
}

// WasRetracted vérifie si un fait a été rétracté pendant cette soumission
func (sc *SubmissionContext) WasRetracted(factID string) bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.factsRetracted[factID]
}

// WasSubmitted vérifie si un fait a été soumis dans cette soumission
func (sc *SubmissionContext) WasSubmitted(factID string) bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.factsSubmitted[factID]
}

// XupleHandlerFunc est une fonction qui gère la création d'un xuple.
// Elle est appelée par l'action Xuple() dans les règles TSD.
type XupleHandlerFunc func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error

// GetChainMetrics retourne les métriques de performance pour la construction des chaînes alpha
func (rn *ReteNetwork) GetChainMetrics() *ChainBuildMetrics {
	if rn.ChainMetrics == nil {
		rn.ChainMetrics = NewChainBuildMetrics()
	}
	return rn.ChainMetrics
}

// GetBetaSharingStats retourne les statistiques de partage des JoinNodes
func (rn *ReteNetwork) GetBetaSharingStats() *BetaSharingStats {
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

// GetXupleManager retourne le gestionnaire de xuples
func (rn *ReteNetwork) GetXupleManager() interface{} {
	return rn.XupleManager
}

// SetXupleManager configure le gestionnaire de xuples
func (rn *ReteNetwork) SetXupleManager(xupleManager interface{}) {
	rn.XupleManager = xupleManager
}

// SetXupleHandler configure la fonction handler pour l'action Xuple
func (rn *ReteNetwork) SetXupleHandler(handler XupleHandlerFunc) {
	rn.xupleHandlerFunc = handler
}

// GetXupleHandler retourne la fonction handler pour l'action Xuple
func (rn *ReteNetwork) GetXupleHandler() XupleHandlerFunc {
	return rn.xupleHandlerFunc
}

// SetXupleSpaceDefinitions stocke les définitions de xuple-spaces parsées
func (rn *ReteNetwork) SetXupleSpaceDefinitions(definitions []interface{}) {
	rn.xupleSpaceDefinitions = definitions
}

// GetXupleSpaceDefinitions retourne les définitions de xuple-spaces stockées.
// Ces définitions sont utilisées par le package api pour créer automatiquement
// les xuple-spaces lors de l'ingestion.
func (rn *ReteNetwork) GetXupleSpaceDefinitions() []interface{} {
	return rn.xupleSpaceDefinitions
}

// GetLogger retourne le logger actuel du réseau
func (rn *ReteNetwork) GetLogger() *Logger {
	if rn.logger == nil {
		rn.logger = NewLogger(LogLevelInfo, nil)
	}
	return rn.logger
}

// SetActionObserver configure l'observateur pour tous les terminal nodes.
//
// Cette méthode configure l'observer pour tous les terminal nodes
// existants ET futurs (via les méthodes qui ajoutent des terminal nodes).
//
// Thread-Safety :
//   - Méthode thread-safe si appelée avant démarrage du réseau
//   - Si appelée pendant l'exécution, risque de race condition
//   - Recommandé : appeler pendant la phase d'initialisation
//
// Paramètres :
//   - observer : observateur à configurer (peut être nil pour désactiver)
func (rn *ReteNetwork) SetActionObserver(observer ActionObserver) {
	if observer == nil {
		observer = &NoOpObserver{}
	}

	rn.actionObserver = observer

	// Configurer tous les terminal nodes existants
	for _, terminal := range rn.TerminalNodes {
		terminal.SetObserver(observer)
	}
}

// GetActionObserver retourne l'observateur configuré.
func (rn *ReteNetwork) GetActionObserver() ActionObserver {
	if rn.actionObserver == nil {
		return &NoOpObserver{}
	}
	return rn.actionObserver
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
