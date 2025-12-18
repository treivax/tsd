// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// SubmitFact soumet un nouveau fait au réseau RETE
// Si une transaction est active, la commande est enregistrée pour rollback
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
	rn.logger.Debug("🔥 Soumission fait: %s", fact.String())

	// Debug logging for E2E debugging
	debugLogger := GetDebugLogger()
	debugLogger.LogFactSubmission(fact.Type, fact.ID, fact.Fields)

	// Vérifier si une transaction est active
	tx := rn.GetTransaction()
	if tx != nil && tx.IsActive {
		// Mode transactionnel : enregistrer la commande
		cmd := NewAddFactCommand(rn.Storage, fact)
		if err := tx.RecordAndExecute(cmd); err != nil {
			return err
		}
		// Propager le fait dans le réseau
		return rn.RootNode.ActivateRight(fact)
	}

	// Mode normal : exécution directe
	if err := rn.Storage.AddFact(fact); err != nil {
		return err
	}
	return rn.RootNode.ActivateRight(fact)
}

// RemoveFact supprime un fait du réseau
// Si une transaction est active, la commande est enregistrée pour rollback
func (rn *ReteNetwork) RemoveFact(factID string) error {
	tx := rn.GetTransaction()
	if tx != nil && tx.IsActive {
		cmd := NewRemoveFactCommand(rn.Storage, factID)
		return tx.RecordAndExecute(cmd)
	}

	return rn.Storage.RemoveFact(factID)
}

// InsertFact insère dynamiquement un nouveau fait dans le réseau RETE.
// Cette méthode valide le fait, l'ajoute au storage et le propage dans le réseau.
//
// Paramètres:
//   - fact: le fait à insérer
//
// Retourne:
//   - error: erreur si le fait est invalide ou s'il existe déjà
func (rn *ReteNetwork) InsertFact(fact *Fact) error {
	// Validation du fait
	if fact == nil {
		return fmt.Errorf("fact cannot be nil")
	}
	if fact.Type == "" {
		return fmt.Errorf("fact type cannot be empty")
	}
	if fact.ID == "" {
		return fmt.Errorf("fact ID cannot be empty")
	}

	// Vérifier si le fait existe déjà
	internalID := fact.GetInternalID()
	if existingFact := rn.Storage.GetFact(internalID); existingFact != nil {
		return fmt.Errorf("fact with ID '%s' and type '%s' already exists", fact.ID, fact.Type)
	}

	// Utiliser SubmitFact qui gère déjà le storage et la propagation
	return rn.SubmitFact(fact)
}

// UpdateFact met à jour dynamiquement un fait existant dans le réseau RETE.
// Cette méthode remplace les champs du fait existant et propage les changements.
//
// Paramètres:
//   - fact: le fait avec les nouvelles valeurs
//
// Retourne:
//   - error: erreur si le fait est invalide ou n'existe pas
func (rn *ReteNetwork) UpdateFact(fact *Fact) error {
	// Validation du fait
	if fact == nil {
		return fmt.Errorf("fact cannot be nil")
	}
	if fact.Type == "" {
		return fmt.Errorf("fact type cannot be empty")
	}
	if fact.ID == "" {
		return fmt.Errorf("fact ID cannot be empty")
	}

	// Vérifier que le fait existe
	internalID := fact.GetInternalID()
	existingFact := rn.Storage.GetFact(internalID)
	if existingFact == nil {
		return fmt.Errorf("fact with ID '%s' and type '%s' not found", fact.ID, fact.Type)
	}

	rn.logger.Debug("🔄 Mise à jour du fait: %s", internalID)

	// Stratégie: Retract puis Insert pour garantir la cohérence
	// Cela propage correctement la suppression puis l'ajout dans le réseau

	// 1. Rétracter l'ancien fait (propage la suppression)
	if err := rn.RetractFact(internalID); err != nil {
		return fmt.Errorf("failed to retract old fact: %w", err)
	}

	// 2. Insérer le nouveau fait avec les valeurs mises à jour (propage l'ajout)
	if err := rn.SubmitFact(fact); err != nil {
		return fmt.Errorf("failed to submit updated fact: %w", err)
	}

	return nil
}

// RepropagateExistingFact propage un fait déjà existant dans le réseau vers les nouveaux nœuds
// Cette fonction est utilisée en mode incrémental pour propager les faits existants
// vers les nouvelles règles qui viennent d'être ajoutées au réseau
func (rn *ReteNetwork) RepropagateExistingFact(fact *Fact) error {
	// Ne pas ajouter le fait au RootNode ou TypeNode (il y est déjà)
	// Propager directement aux enfants du TypeNode (AlphaNodes, etc.)
	typeNode, exists := rn.TypeNodes[fact.Type]
	if !exists {
		return fmt.Errorf("type %s non trouvé dans le réseau", fact.Type)
	}

	// Créer un token pour ce fait
	token := &Token{
		ID:     fmt.Sprintf("token_%s_%s", fact.Type, fact.ID),
		NodeID: typeNode.GetID(),
		Facts:  []*Fact{fact},
	}

	// Propager directement aux enfants du TypeNode sans ajouter à sa mémoire
	return typeNode.PropagateToChildren(fact, token)
}

// waitForFactPersistence attend qu'un fait soit persisté avec retry + backoff exponentiel
// Cette fonction implémente la barrière de synchronisation de la Phase 2
func (rn *ReteNetwork) waitForFactPersistence(fact *Fact, timeout time.Duration) error {
	return rn.waitForFactPersistenceWithMetrics(fact, timeout, nil)
}

// waitForFactPersistenceWithMetrics attend la persistance d'un fait avec collecte de métriques optionnelle
func (rn *ReteNetwork) waitForFactPersistenceWithMetrics(fact *Fact, timeout time.Duration, metricsCollector *CoherenceMetricsCollector) error {
	internalID := fact.GetInternalID()
	deadline := time.Now().Add(timeout)
	attempt := 0

	for time.Now().Before(deadline) {
		attempt++

		// Enregistrer la tentative de vérification
		if metricsCollector != nil {
			metricsCollector.RecordVerifyAttempt()
		}

		// Vérifier si le fait est persisté
		if storedFact := rn.Storage.GetFact(internalID); storedFact != nil {
			// ✅ Fait trouvé
			if attempt > 1 {
				rn.logger.Info("✅ Fait %s persisté après %d tentative(s)", fact.ID, attempt)
				if metricsCollector != nil {
					metricsCollector.RecordFactRetried()
					metricsCollector.RecordRetry(attempt - 1)
				}
			}
			return nil
		}

		// Si on n'a pas dépassé le nombre max de retries, utiliser backoff exponentiel
		if attempt < rn.MaxVerifyRetries {
			// Backoff exponentiel: 10ms, 20ms, 40ms, 80ms, 160ms, 320ms, max 500ms
			backoff := rn.VerifyRetryDelay * time.Duration(1<<uint(attempt-1))
			if backoff > 500*time.Millisecond {
				backoff = 500 * time.Millisecond
			}
			time.Sleep(backoff)
		} else {
			// Après max retries, attendre un peu avant de vérifier à nouveau
			time.Sleep(100 * time.Millisecond)
		}
	}

	// ❌ Timeout dépassé
	return fmt.Errorf("timeout: fait %s (ID interne: %s) non persisté après %v",
		fact.ID, internalID, timeout)
}

// SubmitFactsFromGrammar soumet une liste de faits au réseau RETE avec garanties de synchronisation (Phase 2)
// Cette fonction garantit que tous les faits soumis sont persistés et visibles avant de retourner
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
	return rn.submitFactsFromGrammarWithMetrics(facts, nil)
}

// SubmitFactsFromGrammarWithMetrics soumet des faits avec collecte de métriques de cohérence
func (rn *ReteNetwork) SubmitFactsFromGrammarWithMetrics(facts []map[string]interface{}, metricsCollector *CoherenceMetricsCollector) error {
	return rn.submitFactsFromGrammarWithMetrics(facts, metricsCollector)
}

// submitFactsFromGrammarWithMetrics est l'implémentation interne avec support optionnel des métriques
func (rn *ReteNetwork) submitFactsFromGrammarWithMetrics(facts []map[string]interface{}, metricsCollector *CoherenceMetricsCollector) error {
	if len(facts) == 0 {
		return nil
	}

	// Debug: dump network structure before fact submission
	debugLogger := GetDebugLogger()
	if debugLogger.IsEnabled() {
		debugLogger.LogNetworkStructure(rn)
	}

	// Démarrer la phase de soumission si collecteur disponible
	if metricsCollector != nil {
		metricsCollector.StartPhase("fact_submission")
		defer func() {
			metricsCollector.EndPhase("fact_submission", len(facts), true)
		}()
	}

	// Timeout par fait : timeout total divisé par nombre de faits
	// Minimum 1 seconde par fait pour éviter les timeouts prématurés
	timeoutPerFact := rn.SubmissionTimeout / time.Duration(len(facts))
	if timeoutPerFact < 1*time.Second {
		timeoutPerFact = 1 * time.Second
	}

	// Compteurs pour garantir la cohérence
	factsSubmitted := 0
	factsPersisted := 0
	startTime := time.Now()

	for i, factMap := range facts {
		// 1. Convertir le map en Fact
		factID := fmt.Sprintf("fact_%d", i)
		if id, ok := factMap["id"].(string); ok {
			factID = id
		}

		factType := "unknown"
		// Chercher "type" ou "reteType" (ConvertFactsToReteFormat utilise "reteType")
		if typ, ok := factMap["type"].(string); ok {
			factType = typ
		} else if typ, ok := factMap["reteType"].(string); ok {
			factType = typ
		}

		fact := &Fact{
			ID:     factID,
			Type:   factType,
			Fields: make(map[string]interface{}),
		}

		// Copier tous les champs
		for key, value := range factMap {
			if key != "type" && key != "reteType" {
				fact.Fields[key] = value
			}
		}

		// 2. Soumettre le fait au réseau RETE
		if metricsCollector != nil {
			metricsCollector.RecordFactSubmitted()
		}

		if err := rn.SubmitFact(fact); err != nil {
			if metricsCollector != nil {
				metricsCollector.RecordFactFailed()
			}
			return fmt.Errorf("erreur soumission fait %s: %w", fact.ID, err)
		}
		factsSubmitted++

		// 3. Barrière de synchronisation Phase 2 : attendre la persistance avec retry
		waitStart := time.Now()
		err := rn.waitForFactPersistenceWithMetrics(fact, timeoutPerFact, metricsCollector)
		waitDuration := time.Since(waitStart)

		if metricsCollector != nil {
			metricsCollector.RecordWaitTime(waitDuration)
		}

		if err != nil {
			if metricsCollector != nil {
				metricsCollector.RecordTimeout()
				metricsCollector.RecordFactFailed()
			}
			return fmt.Errorf("échec synchronisation fait %s: %w", fact.ID, err)
		}

		if metricsCollector != nil {
			metricsCollector.RecordFactPersisted()
		}
		factsPersisted++
	}

	duration := time.Since(startTime)

	if metricsCollector != nil {
		metricsCollector.RecordSubmissionTime(duration)
	}

	// 4. Vérification finale de cohérence
	if factsSubmitted != factsPersisted {
		return fmt.Errorf("incohérence détectée: %d faits soumis mais seulement %d persistés dans le storage",
			factsSubmitted, factsPersisted)
	}

	rn.logger.Info("✅ Phase 2 - Synchronisation complète: %d/%d faits persistés en %v", factsPersisted, factsSubmitted, duration)

	return nil
}

// RetractFact supprime dynamiquement un fait du réseau RETE.
// Cette méthode retire le fait du storage et propage la suppression.
//
// Paramètres:
//   - factID: l'identifiant interne du fait (format: Type_ID)
//
// Retourne:
//   - error: erreur si l'ID est vide ou si le fait n'existe pas
func (rn *ReteNetwork) RetractFact(factID string) error {
	// Validation de l'ID
	if factID == "" {
		return fmt.Errorf("fact ID cannot be empty")
	}

	rn.logger.Info("🗑️ Rétractation du fait: %s", factID)

	// Vérifier que le fait existe
	existingFact := rn.Storage.GetFact(factID)
	if existingFact == nil {
		return fmt.Errorf("fact with ID '%s' not found", factID)
	}

	// Utiliser RemoveFact qui gère le storage et les transactions
	if err := rn.RemoveFact(factID); err != nil {
		return fmt.Errorf("failed to remove fact from storage: %w", err)
	}

	// Propager la rétractation dans le réseau
	return rn.RootNode.ActivateRetract(factID)
}

// Reset clears the entire RETE network and resets it to an empty state.
// This removes all facts, rules, types, and network nodes.
// After calling Reset, the network is ready to accept new definitions from scratch.
func (rn *ReteNetwork) Reset() {
	rn.logger.Info("🧹 Réinitialisation complète du réseau RETE")

	// Clear all node collections
	rn.TypeNodes = make(map[string]*TypeNode)
	rn.AlphaNodes = make(map[string]*AlphaNode)
	rn.BetaNodes = make(map[string]interface{})
	rn.TerminalNodes = make(map[string]*TerminalNode)
	rn.Types = make([]TypeDefinition, 0)
	rn.BetaBuilder = nil

	// Reset lifecycle manager
	if rn.LifecycleManager != nil {
		rn.LifecycleManager.Reset()
	} else {
		rn.LifecycleManager = NewLifecycleManager()
	}

	// Reset alpha sharing manager
	if rn.AlphaSharingManager != nil {
		rn.AlphaSharingManager.Reset()
	} else {
		rn.AlphaSharingManager = NewAlphaSharingRegistry()
	}

	// Reset passthrough registry
	rn.PassthroughRegistry = make(map[string]*AlphaNode)

	// Recreate a fresh root node with the existing storage
	rn.RootNode = NewRootNode(rn.Storage)

	rn.logger.Info("✅ Réseau RETE réinitialisé avec succès")
}

// ClearMemory efface uniquement les mémoires (faits et tokens) de tous les nœuds
// sans détruire la structure du réseau
func (rn *ReteNetwork) ClearMemory() {
	rn.logger.Info("🧹 Nettoyage de la mémoire du réseau RETE")

	// Clear TypeNode memories
	for _, typeNode := range rn.TypeNodes {
		typeNode.mutex.Lock()
		typeNode.Memory.Facts = make(map[string]*Fact)
		typeNode.Memory.Tokens = make(map[string]*Token)
		typeNode.mutex.Unlock()
	}

	// Clear AlphaNode memories
	for _, alphaNode := range rn.AlphaNodes {
		alphaNode.mutex.Lock()
		alphaNode.Memory.Facts = make(map[string]*Fact)
		alphaNode.Memory.Tokens = make(map[string]*Token)
		alphaNode.mutex.Unlock()
	}

	// Clear BetaNode memories (JoinNodes, etc.)
	for _, betaNode := range rn.BetaNodes {
		if node, ok := betaNode.(Node); ok {
			node.GetMemory().Facts = make(map[string]*Fact)
			node.GetMemory().Tokens = make(map[string]*Token)
		}
	}

	// Clear TerminalNode memories
	for _, terminalNode := range rn.TerminalNodes {
		terminalNode.mutex.Lock()
		terminalNode.Memory.Facts = make(map[string]*Fact)
		terminalNode.Memory.Tokens = make(map[string]*Token)
		terminalNode.mutex.Unlock()
	}

	rn.logger.Info("✅ Mémoire du réseau RETE nettoyée avec succès")
}

// GarbageCollect nettoie et libère les ressources du réseau
func (rn *ReteNetwork) GarbageCollect() {
	// 1. Vider les caches
	if rn.ArithmeticResultCache != nil {
		rn.ArithmeticResultCache.Clear()
	}

	if rn.BetaSharingRegistry != nil {
		rn.BetaSharingRegistry.Clear()
	}

	if rn.AlphaSharingManager != nil {
		rn.AlphaSharingManager.Clear()
	}

	// 2. Nettoyer les nœuds et supprimer les références
	// TypeNodes
	for _, node := range rn.TypeNodes {
		if node != nil && node.Memory != nil {
			node.Memory.Facts = make(map[string]*Fact)
			node.Memory.Tokens = make(map[string]*Token)
		}
		if node != nil {
			node.Children = nil
		}
	}
	rn.TypeNodes = make(map[string]*TypeNode)

	// AlphaNodes
	for _, node := range rn.AlphaNodes {
		if node != nil && node.Memory != nil {
			node.Memory.Facts = make(map[string]*Fact)
			node.Memory.Tokens = make(map[string]*Token)
		}
		if node != nil {
			node.Children = nil
		}
	}
	rn.AlphaNodes = make(map[string]*AlphaNode)

	// BetaNodes
	rn.BetaNodes = make(map[string]interface{})

	// TerminalNodes
	for _, node := range rn.TerminalNodes {
		if node != nil && node.Memory != nil {
			node.Memory.Facts = make(map[string]*Fact)
			node.Memory.Tokens = make(map[string]*Token)
		}
		if node != nil {
			node.Children = nil
		}
	}
	rn.TerminalNodes = make(map[string]*TerminalNode)

	// 3. Vider les types
	rn.Types = make([]TypeDefinition, 0)

	// 4. Vider le PassthroughRegistry
	rn.PassthroughRegistry = make(map[string]*AlphaNode)

	// 5. Nettoyer le LifecycleManager
	if rn.LifecycleManager != nil {
		rn.LifecycleManager.Cleanup()
	}

	// 6. Nettoyer l'ActionExecutor
	if rn.ActionExecutor != nil {
		// ActionExecutor n'a pas de méthode Cleanup pour l'instant
		// mais on pourrait en ajouter une si nécessaire
	}

	// 7. Nettoyer le Storage
	if rn.Storage != nil {
		rn.Storage.Clear()
	}

	// 8. Réinitialiser le RootNode
	if rn.RootNode != nil && rn.RootNode.Memory != nil {
		rn.RootNode.Memory.Facts = make(map[string]*Fact)
		rn.RootNode.Memory.Tokens = make(map[string]*Token)
	}
}
