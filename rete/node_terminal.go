// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"
)

type TerminalNode struct {
	BaseNode
	Action *Action `json:"action"`

	// Observer pour notification des exécutions
	observer ActionObserver

	// Statistiques d'exécution (pour debug/tests)
	lastExecutionResult *ExecutionResult
	executionCount      int64
	statsMutex          sync.RWMutex
}

// NewTerminalNode crée un nouveau nœud terminal
func NewTerminalNode(nodeID string, action *Action, storage Storage) *TerminalNode {
	return &TerminalNode{
		BaseNode: BaseNode{
			ID:        nodeID,
			Type:      "terminal",
			Memory:    &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children:  make([]Node, 0), // Les nœuds terminaux n'ont pas d'enfants
			Storage:   storage,
			createdAt: time.Now(),
		},
		Action:         action,
		observer:       &NoOpObserver{}, // Observer par défaut
		executionCount: 0,
	}
}

// ActivateLeft exécute immédiatement l'action sans stocker le token.
//
// Process :
//  1. Enregistre l'activation (métriques)
//  2. Exécute l'action immédiatement
//  3. Notifie l'observer du résultat
//  4. NE STOCKE PAS le token (exécution immédiate)
//
// Le token contient tous les bindings (via BindingChain) nécessaires
// pour l'évaluation des arguments de l'action.
//
// Thread-Safety :
//   - Méthode thread-safe
//   - Les statistiques sont protégées par statsMutex
//   - L'observer DOIT être thread-safe
//
// Paramètres :
//   - token : token contenant les faits et bindings déclencheurs
//
// Retourne :
//   - error : erreur si l'exécution de l'action échoue
func (tn *TerminalNode) ActivateLeft(token *Token) error {
	// Enregistrer l'activation (métriques réseau)
	tn.recordActivation()

	// PAS DE STOCKAGE - Exécuter directement
	start := time.Now()
	err := tn.executeAction(token)
	duration := time.Since(start)

	// Créer le résultat d'exécution
	result := ExecutionResult{
		Success:  err == nil,
		Error:    err,
		Duration: duration,
		Context: ActionContext{
			ActionName: tn.getActionName(),
			RuleName:   tn.getRuleName(),
			Token:      token,
			Network:    tn.BaseNode.GetNetwork(),
			Timestamp:  start,
		},
		Arguments: tn.extractArguments(token),
	}

	// Mettre à jour les statistiques (pour debug/tests)
	tn.updateStats(result)

	// Notifier l'observer
	if tn.observer != nil {
		tn.observer.OnActionExecuted(result)
	}

	return err
}

// ActivateRetract retrait des tokens contenant le fait rétracté
// factID doit être l'identifiant interne (Type_ID)
func (tn *TerminalNode) ActivateRetract(factID string) error {
	tn.mutex.Lock()
	var tokensToRemove []string
	for tokenID, token := range tn.Memory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				tokensToRemove = append(tokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range tokensToRemove {
		delete(tn.Memory.Tokens, tokenID)
	}
	tn.mutex.Unlock()
	if len(tokensToRemove) > 0 {
		fmt.Printf("🗑️  [TERMINAL_%s] Rétractation: %d tokens retirés\n", tn.ID, len(tokensToRemove))
	}
	return nil
}

// GetTriggeredActions retourne les actions déclenchées (pour les tests)
// DEPRECATED: Utiliser GetExecutionCount() et observer pattern à la place
func (tn *TerminalNode) GetTriggeredActions() []*Action {
	tn.statsMutex.RLock()
	defer tn.statsMutex.RUnlock()

	// Retourner autant de copies de l'action que d'exécutions
	actions := make([]*Action, 0, tn.executionCount)
	for i := int64(0); i < tn.executionCount; i++ {
		actions = append(actions, tn.Action)
	}
	return actions
}

// ActivateRight (non utilisé pour les nœuds terminaux)
func (tn *TerminalNode) ActivateRight(fact *Fact) error {
	return fmt.Errorf("les nœuds terminaux ne reçoivent pas de faits directement")
}

// SetNetwork définit la référence au réseau RETE
func (tn *TerminalNode) SetNetwork(network *ReteNetwork) {
	tn.BaseNode.SetNetwork(network)
}

// executeAction exécute l'action avec le contexte du token.
//
// Process :
//  1. Vérifie qu'une action est définie
//  2. Publie l'activation vers le xuple-space si configuré (xuples)
//  3. Délègue l'exécution au ActionExecutor du réseau
//
// Le ActionExecutor crée un ExecutionContext avec token.Bindings,
// permettant l'accès aux variables via BindingChain.
//
// Note: L'affichage console a été supprimé (violation principe NO HARDCODING).
// Les activations sont maintenant gérées via le module xuples et peuvent être
// récupérées programmatiquement.
//
// Paramètres :
//   - token : token contenant les faits et bindings
//
// Retourne :
//   - error : erreur si l'exécution échoue
func (tn *TerminalNode) executeAction(token *Token) error {
	// Les actions sont maintenant obligatoires dans la grammaire
	// Mais nous gardons cette vérification par sécurité
	if tn.Action == nil {
		return fmt.Errorf("aucune action définie pour le nœud %s", tn.ID)
	}

	// TODO(xuples): Publier vers XupleSpace si configuré
	// Exemple d'intégration future :
	//
	// network := tn.BaseNode.GetNetwork()
	// if network != nil && network.XuplePublisher != nil {
	//     if err := network.XuplePublisher.Publish(tn.Action, token, token.Facts); err != nil {
	//         // Log l'erreur mais ne bloque pas l'exécution
	//         network.Logger.Printf("⚠️  Erreur publication xuple: %v", err)
	//     }
	// }

	// Exécuter réellement l'action avec l'ActionExecutor
	network := tn.BaseNode.GetNetwork()
	if network != nil && network.ActionExecutor != nil {
		return network.ActionExecutor.ExecuteAction(tn.Action, token)
	}

	return nil
}

// Clone crée une copie profonde du TerminalNode
func (tn *TerminalNode) Clone() *TerminalNode {
	clone := &TerminalNode{
		BaseNode: BaseNode{
			ID:       tn.ID,
			Type:     tn.Type,
			Memory:   tn.Memory.Clone(),
			Children: make([]Node, len(tn.Children)),
			Storage:  tn.Storage,
		},
		Action:         tn.Action.Clone(),
		observer:       &NoOpObserver{}, // Ne pas cloner l'observer
		executionCount: 0,               // Réinitialiser les stats
	}

	// Copier les enfants
	copy(clone.Children, tn.Children)

	return clone
}

// SetObserver configure l'observateur d'actions.
func (tn *TerminalNode) SetObserver(observer ActionObserver) {
	if observer == nil {
		observer = &NoOpObserver{}
	}
	tn.observer = observer
}

// GetExecutionCount retourne le nombre total d'exécutions.
// Utilisé principalement pour les tests.
func (tn *TerminalNode) GetExecutionCount() int64 {
	tn.statsMutex.RLock()
	defer tn.statsMutex.RUnlock()
	return tn.executionCount
}

// GetLastExecutionResult retourne le dernier résultat d'exécution.
// Utilisé principalement pour les tests.
func (tn *TerminalNode) GetLastExecutionResult() *ExecutionResult {
	tn.statsMutex.RLock()
	defer tn.statsMutex.RUnlock()
	if tn.lastExecutionResult == nil {
		return nil
	}
	// Retourner une copie pour éviter modifications concurrentes
	resultCopy := *tn.lastExecutionResult
	return &resultCopy
}

// ResetExecutionStats réinitialise les statistiques d'exécution.
// Utilisé principalement pour les tests.
func (tn *TerminalNode) ResetExecutionStats() {
	tn.statsMutex.Lock()
	defer tn.statsMutex.Unlock()
	tn.lastExecutionResult = nil
	tn.executionCount = 0
}

// updateStats met à jour les statistiques internes.
func (tn *TerminalNode) updateStats(result ExecutionResult) {
	tn.statsMutex.Lock()
	defer tn.statsMutex.Unlock()

	// Copier le résultat
	resultCopy := result
	tn.lastExecutionResult = &resultCopy
	tn.executionCount++
}

// getActionName retourne le nom de l'action.
func (tn *TerminalNode) getActionName() string {
	if tn.Action == nil {
		return "unknown"
	}
	jobs := tn.Action.GetJobs()
	if len(jobs) > 0 {
		return jobs[0].Name
	}
	return "unknown"
}

// getRuleName extrait le nom de la règle depuis l'ID du nœud.
func (tn *TerminalNode) getRuleName() string {
	// L'ID du terminal node contient le nom de la règle
	// Format: "terminal_<ruleName>"
	if len(tn.ID) > 9 && tn.ID[:9] == "terminal_" {
		return tn.ID[9:]
	}
	return tn.ID
}

// extractArguments extrait les arguments bruts de l'action.
func (tn *TerminalNode) extractArguments(token *Token) []interface{} {
	if tn.Action == nil {
		return nil
	}

	jobs := tn.Action.GetJobs()
	if len(jobs) == 0 {
		return nil
	}

	return jobs[0].Args
}
