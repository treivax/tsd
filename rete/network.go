// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"

	"github.com/treivax/tsd/rete/delta"
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

	// Propagation delta (Prompt 06)
	DeltaPropagator        *delta.DeltaPropagator   `json:"-"` // Propagateur delta pour optimisation Update
	DependencyIndex        *delta.DependencyIndex   `json:"-"` // Index des dépendances champs→nœuds
	EnableDeltaPropagation bool                     `json:"-"` // Active/désactive la propagation delta
	IntegrationHelper      *delta.IntegrationHelper `json:"-"` // Helper d'intégration delta

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

// InitializeDeltaPropagation initialise le système de propagation delta.
//
// Cette méthode construit l'index de dépendances depuis le réseau existant
// et crée le DeltaPropagator configuré avec les callbacks vers le réseau RETE.
//
// Doit être appelée après la construction complète du réseau RETE.
//
// L'initialisation comprend :
//  1. Construction de l'index de dépendances
//  2. Indexation des nœuds alpha (conditions sur champs)
//  3. Indexation des nœuds beta (jointures)
//  4. Indexation des nœuds terminaux (actions)
//  5. Création du propagateur avec configuration
//
// Retourne une erreur si l'initialisation échoue.
func (rn *ReteNetwork) InitializeDeltaPropagation() error {
	// 1. Créer l'index de dépendances
	rn.DependencyIndex = delta.NewDependencyIndex()

	// 2. Créer le builder d'index
	indexBuilder := delta.NewIndexBuilder()
	indexBuilder.EnableDiagnostics()

	// 3. Indexer les nœuds alpha
	for nodeID, alphaNode := range rn.AlphaNodes {
		if alphaNode.Condition == nil {
			continue
		}

		// Extraire le type de fait depuis le nœud (via le TypeNode parent)
		factType := rn.inferFactTypeForNode(nodeID)

		err := indexBuilder.BuildFromAlphaNode(
			rn.DependencyIndex,
			nodeID,
			factType,
			alphaNode.Condition,
		)
		if err != nil {
			rn.logger.Debug("Failed to index alpha node %s: %v", nodeID, err)
		}
	}

	// 4. Indexer les nœuds beta (jointures)
	// Note: L'indexation des nœuds beta nécessite l'extraction des conditions
	// de jointure. Implémentation simplifiée pour l'instant.
	// Future: Indexation avancée basée sur les conditions de jointure pour
	// optimiser la propagation incrémentale des deltas.

	// 5. Indexer les nœuds terminaux (actions)
	for nodeID, terminalNode := range rn.TerminalNodes {
		if terminalNode.Action == nil {
			continue
		}

		// Extraire le type de fait depuis le nœud
		factType := rn.inferFactTypeForNode(nodeID)

		// Convertir l'action en slice pour l'indexation
		actions := []interface{}{terminalNode.Action}

		err := indexBuilder.BuildFromTerminalNode(
			rn.DependencyIndex,
			nodeID,
			factType,
			actions,
		)
		if err != nil {
			rn.logger.Debug("Failed to index terminal node %s: %v", nodeID, err)
		}
	}

	// 6. Créer le détecteur de delta
	detector := delta.NewDeltaDetectorWithConfig(delta.DefaultDetectorConfig())

	// 7. Créer le propagateur
	propagator, err := delta.NewDeltaPropagatorBuilder().
		WithDetector(detector).
		WithIndex(rn.DependencyIndex).
		WithStrategy(&delta.SequentialStrategy{}).
		WithConfig(delta.DefaultPropagationConfig()).
		WithPropagateCallback(rn.propagateDeltaToNode).
		Build()

	if err != nil {
		return fmt.Errorf("failed to build delta propagator: %w", err)
	}

	rn.DeltaPropagator = propagator

	// 8. Créer le helper d'intégration
	callbacks := &reteNetworkCallbacks{network: rn}
	rn.IntegrationHelper = delta.NewIntegrationHelper(
		rn.DeltaPropagator,
		rn.DependencyIndex,
		callbacks,
	)

	rn.logger.Info("Delta propagation initialized successfully")
	stats := rn.DependencyIndex.GetStats()
	rn.logger.Info("Index stats: %d nodes, %d field dependencies", stats.NodeCount, stats.FieldCount)

	return nil
}

// propagateDeltaToNode est le callback pour propager un delta vers un nœud RETE.
//
// Cette méthode est appelée par le DeltaPropagator pour chaque nœud affecté
// par un changement de fait.
//
// État actuel: Infrastructure de base en place avec indexation des nœuds.
// La propagation incrémentale complète est planifiée pour une future version.
//
// Fonctionnalités planifiées pour la propagation incrémentale:
// - Alpha nodes: Ré-évaluation de la condition avec le fait modifié
// - Beta nodes: Ré-évaluation des jointures affectées uniquement
// - Terminal nodes: Déclenchement conditionnel des actions
// - Optimisation: Éviter la re-propagation complète du réseau
//
// Pour l'instant, le système utilise une approche de re-propagation complète
// qui est fonctionnelle mais moins optimisée.
func (rn *ReteNetwork) propagateDeltaToNode(nodeID string, factDelta *delta.FactDelta) error {
	// Rechercher le type de nœud
	if _, exists := rn.AlphaNodes[nodeID]; exists {
		// Propagation incrémentale non implémentée - utilise re-propagation complète
		rn.logger.Debug("Delta propagation to alpha node %s (using full re-propagation)", nodeID)
		return nil
	}

	if _, exists := rn.BetaNodes[nodeID]; exists {
		// Propagation incrémentale non implémentée - utilise re-propagation complète
		rn.logger.Debug("Delta propagation to beta node %s (using full re-propagation)", nodeID)
		return nil
	}

	if _, exists := rn.TerminalNodes[nodeID]; exists {
		// Propagation incrémentale non implémentée - utilise re-propagation complète
		rn.logger.Debug("Delta propagation to terminal node %s (using full re-propagation)", nodeID)
		return nil
	}

	return fmt.Errorf("node not found: %s", nodeID)
}

// reteNetworkCallbacks implémente delta.NetworkCallbacks pour ReteNetwork.
type reteNetworkCallbacks struct {
	network *ReteNetwork
}

// newReteNetworkCallbacks crée une nouvelle instance de callbacks.
func newReteNetworkCallbacks(network *ReteNetwork) *reteNetworkCallbacks {
	return &reteNetworkCallbacks{
		network: network,
	}
}

// PropagateToNode propage un delta vers un nœud.
func (cb *reteNetworkCallbacks) PropagateToNode(nodeID string, factDelta *delta.FactDelta) error {
	return cb.network.propagateDeltaToNode(nodeID, factDelta)
}

// GetNode récupère un nœud par son ID.
func (cb *reteNetworkCallbacks) GetNode(nodeID string) (interface{}, error) {
	if node, exists := cb.network.AlphaNodes[nodeID]; exists {
		return node, nil
	}
	if node, exists := cb.network.BetaNodes[nodeID]; exists {
		return node, nil
	}
	if node, exists := cb.network.TerminalNodes[nodeID]; exists {
		return node, nil
	}
	return nil, fmt.Errorf("node not found: %s", nodeID)
}

// UpdateStorage met à jour le storage avec le fait modifié.
func (cb *reteNetworkCallbacks) UpdateStorage(factID string, newFact map[string]interface{}) error {
	if cb.network.Storage == nil {
		return fmt.Errorf("storage not initialized")
	}

	existingFact := cb.network.Storage.GetFact(factID)
	if existingFact == nil {
		return fmt.Errorf("fact not found in storage: %s", factID)
	}

	// Mettre à jour les champs du fait
	for key, value := range newFact {
		existingFact.Fields[key] = value
	}

	return nil
}

// inferFactTypeForNode tente de déduire le type de fait pour un nœud.
//
// Cette méthode remonte l'arbre des nœuds pour trouver le TypeNode parent.
// Si aucun TypeNode n'est trouvé, retourne "Unknown".
//
// Note: Cette méthode utilise une heuristique simple. Une future amélioration
// pourrait utiliser un index type->nœuds pour des performances optimales.
func (rn *ReteNetwork) inferFactTypeForNode(nodeID string) string {
	// Parcourir les TypeNodes pour trouver celui qui contient ce nœud
	for factType, typeNode := range rn.TypeNodes {
		// Vérifier si le nœud est dans les enfants du TypeNode
		if rn.isNodeDescendantOf(nodeID, typeNode) {
			return factType
		}
	}

	// Pas trouvé - retourner "Unknown"
	return "Unknown"
}

// isNodeDescendantOf vérifie si un nœud est un descendant d'un autre nœud.
//
// Note: Implémentation simplifiée retournant toujours false.
// Future: Implémenter la recherche récursive dans l'arbre RETE si nécessaire
// pour des optimisations avancées de propagation delta.
func (rn *ReteNetwork) isNodeDescendantOf(nodeID string, parent Node) bool {
	// Pour l'instant, implémentation simplifiée
	// Dans une version complète, il faudrait parcourir l'arbre
	return false
}
