// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// AlphaChain représente une chaîne d'AlphaNodes construite pour un ensemble de conditions.
//
// Une chaîne alpha est une séquence ordonnée de nœuds alpha qui évaluent des conditions
// successives sur une même variable. Chaque nœud évalue une condition et propage les faits
// correspondants au nœud suivant dans la chaîne.
//
// Structure de chaîne typique:
//
//	TypeNode(Person)
//	  └── AlphaNode(p.age >= 18)
//	       └── AlphaNode(p.city == "Paris")
//	            └── TerminalNode(rule_terminal)
//
// Propriétés:
//   - len(Nodes) == len(Hashes) (toujours maintenu)
//   - FinalNode == Nodes[len(Nodes)-1] (si non vide)
//   - Ordre des nœuds correspond à l'ordre des conditions dans la règle
//
// Exemple d'utilisation:
//
//	chain, err := builder.BuildChain(conditions, "p", parentNode, "myRule")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Chaîne construite: %d nœuds\n", len(chain.Nodes))
//	stats := builder.GetChainStats(chain)
//	fmt.Printf("Nœuds partagés: %d/%d\n", stats["shared_nodes"], len(chain.Nodes))
type AlphaChain struct {
	Nodes     []*AlphaNode `json:"nodes"`      // Liste ordonnée des nœuds alpha dans la chaîne
	Hashes    []string     `json:"hashes"`     // Hashes correspondants pour chaque nœud
	FinalNode *AlphaNode   `json:"final_node"` // Le dernier nœud de la chaîne
	RuleID    string       `json:"rule_id"`    // ID de la règle pour laquelle la chaîne a été construite
}

// AlphaChainBuilder construit des chaînes d'AlphaNodes avec partage automatique.
//
// Le builder coordonne la construction de chaînes alpha en réutilisant intelligemment
// les nœuds existants et en maintenant un cache des connexions parent→child pour
// éviter les duplications.
//
// Fonctionnalités principales:
//   - Construction séquentielle de chaînes condition par condition
//   - Partage automatique via AlphaSharingRegistry
//   - Cache de connexions pour éviter duplications
//   - Collection de métriques détaillées
//   - Thread-safe avec sync.RWMutex
//
// Flux de construction:
//
//	Pour chaque condition:
//	  1. Calculer hash (avec cache LRU)
//	  2. Chercher nœud existant via hash
//	  3. Si trouvé: réutiliser (refcount++)
//	     Si non: créer nouveau nœud
//	  4. Vérifier connexion parent→child (avec cache)
//	  5. Connecter si nécessaire
//	  6. Enregistrer dans LifecycleManager
//	  7. Nœud devient parent pour suivant
//
// Exemple d'utilisation:
//
//	builder := NewAlphaChainBuilder(network, storage)
//	conditions := []SimpleCondition{
//	    NewSimpleCondition("binaryOperation", "p.age", ">", 18),
//	    NewSimpleCondition("binaryOperation", "p.name", "==", "Alice"),
//	}
//	chain, err := builder.BuildChain(conditions, "p", typeNode, "rule1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Accéder aux métriques
//	metrics := builder.GetMetrics()
//	fmt.Printf("Sharing ratio: %.1f%%\n", metrics.SharingRatio * 100)
//
// Thread-safety:
//   - Toutes les opérations publiques sont thread-safe
//   - Le cache de connexions est protégé par mutex
//   - Peut être utilisé concurremment par plusieurs goroutines
type AlphaChainBuilder struct {
	network         *ReteNetwork
	storage         Storage
	connectionCache map[string]bool // Cache pour les connexions existantes (parentID_childID -> bool)
	metrics         *ChainBuildMetrics
	mutex           sync.RWMutex
}

// NewAlphaChainBuilder crée un nouveau constructeur de chaînes alpha avec des métriques neuves.
//
// Cette fonction initialise un builder avec un objet de métriques local. Pour partager
// les métriques entre plusieurs builders (recommandé), utilisez NewAlphaChainBuilderWithMetrics.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//
// Retourne:
//   - Un nouveau builder prêt à l'emploi
//
// Exemple:
//
//	storage := NewMemoryStorage()
//	network := NewReteNetwork(storage)
//	builder := NewAlphaChainBuilder(network, storage)
func NewAlphaChainBuilder(network *ReteNetwork, storage Storage) *AlphaChainBuilder {
	return &AlphaChainBuilder{
		network:         network,
		storage:         storage,
		connectionCache: make(map[string]bool),
		metrics:         NewChainBuildMetrics(),
	}
}

// NewAlphaChainBuilderWithMetrics crée un constructeur avec des métriques partagées.
//
// Recommandé quand le réseau RETE crée son propre builder, permettant de partager
// les métriques entre le builder et d'autres composants.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//   - metrics: Objet de métriques partagé (non nil)
//
// Retourne:
//   - Un nouveau builder utilisant les métriques fournies
//
// Exemple:
//
//	metrics := NewChainBuildMetrics()
//	builder := NewAlphaChainBuilderWithMetrics(network, storage, metrics)
//	// Les métriques sont accessibles via builder.GetMetrics() et directement via 'metrics'
func NewAlphaChainBuilderWithMetrics(network *ReteNetwork, storage Storage, metrics *ChainBuildMetrics) *AlphaChainBuilder {
	return &AlphaChainBuilder{
		network:         network,
		storage:         storage,
		connectionCache: make(map[string]bool),
		metrics:         metrics,
	}
}

// BuildChain construit une chaîne d'AlphaNodes pour un ensemble de conditions
// avec partage automatique des nœuds identiques entre règles.
//
// Cette méthode est le point d'entrée principal pour la construction de chaînes.
// Elle itère sur chaque condition, tente de réutiliser un nœud existant, sinon
// en crée un nouveau, et maintient les connexions parent→child appropriées.
//
// Algorithme:
//
//	Pour chaque condition dans la liste:
//	  1. Convertir SimpleCondition en map
//	  2. Appeler AlphaSharingRegistry.GetOrCreateAlphaNode()
//	     → Calcule hash (avec cache LRU)
//	     → Cherche nœud existant
//	     → Crée nouveau si inexistant
//	  3. Si nœud réutilisé:
//	     - Vérifier connexion avec parent (cache)
//	     - Connecter si nécessaire
//	  4. Si nœud créé:
//	     - Connecter au parent
//	     - Ajouter au réseau
//	     - Mettre en cache la connexion
//	  5. Enregistrer dans LifecycleManager
//	  6. Nœud devient parent pour itération suivante
//
// Paramètres:
//   - conditions: liste de conditions simples dans l'ordre d'évaluation
//   - variableName: nom de la variable (ex: "p", "u") - utilisé dans le hash
//   - parentNode: nœud parent (généralement TypeNode) auquel connecter le premier nœud
//   - ruleID: identifiant unique de la règle pour le lifecycle management
//
// Retourne:
//   - *AlphaChain: la chaîne construite avec tous les nœuds et leurs hashes
//   - error: erreur si conditions vides, parent nil, ou problème de création
//
// Exemple simple:
//
//	conditions := []SimpleCondition{
//	    NewSimpleCondition("binaryOperation", "p.age", ">", 18),
//	}
//	chain, err := builder.BuildChain(conditions, "p", typeNode, "rule_adult")
//	// → Crée: TypeNode → AlphaNode(p.age>18) → Terminal
//
// Exemple avec partage:
//
//	// Règle 1
//	chain1, _ := builder.BuildChain(
//	    []SimpleCondition{NewSimpleCondition("binaryOperation", "p.age", ">", 18)},
//	    "p", typeNode, "rule1")
//	// → Crée nouveau nœud alpha_abc123
//
//	// Règle 2 (même condition)
//	chain2, _ := builder.BuildChain(
//	    []SimpleCondition{NewSimpleCondition("binaryOperation", "p.age", ">", 18)},
//	    "p", typeNode, "rule2")
//	// → Réutilise alpha_abc123 (RefCount=2)
//
// Logs générés:
//
//	🆕 [AlphaChainBuilder] Nouveau nœud alpha alpha_abc123 créé pour la règle rule1 (condition 1/1)
//	🔗 [AlphaChainBuilder] Connexion du nœud alpha_abc123 au parent type_person
//	♻️  [AlphaChainBuilder] Réutilisation du nœud alpha alpha_abc123 pour la règle rule2 (condition 1/1)
//	✓  [AlphaChainBuilder] Nœud alpha_abc123 déjà connecté au parent type_person
func (acb *AlphaChainBuilder) BuildChain(
	conditions []SimpleCondition,
	variableName string,
	parentNode Node,
	ruleID string,
) (*AlphaChain, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("impossible de construire une chaîne sans conditions")
	}

	if parentNode == nil {
		return nil, fmt.Errorf("le nœud parent ne peut pas être nil")
	}

	if acb.network.AlphaSharingManager == nil {
		return nil, fmt.Errorf("AlphaSharingManager non initialisé dans le réseau")
	}

	if acb.network.LifecycleManager == nil {
		return nil, fmt.Errorf("LifecycleManager non initialisé dans le réseau")
	}

	// Démarrer le chronomètre pour les métriques
	startTime := time.Now()
	nodesCreated := 0
	nodesReused := 0
	hashesGenerated := make([]string, 0, len(conditions))

	chain := &AlphaChain{
		Nodes:  make([]*AlphaNode, 0, len(conditions)),
		Hashes: make([]string, 0, len(conditions)),
		RuleID: ruleID,
	}

	currentParent := parentNode

	// Construire la chaîne condition par condition
	for i, condition := range conditions {
		// Convertir SimpleCondition en map pour la condition du nœud alpha
		conditionMap := map[string]interface{}{
			"type":     condition.Type,
			"left":     condition.Left,
			"operator": condition.Operator,
			"right":    condition.Right,
		}

		// Obtenir ou créer l'AlphaNode via le gestionnaire de partage
		alphaNode, hash, reused, err := acb.network.AlphaSharingManager.GetOrCreateAlphaNode(
			conditionMap,
			variableName,
			acb.storage,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la création/récupération du nœud alpha %d: %w", i, err)
		}

		// Ajouter le nœud et son hash à la chaîne
		chain.Nodes = append(chain.Nodes, alphaNode)
		chain.Hashes = append(chain.Hashes, hash)
		hashesGenerated = append(hashesGenerated, hash)

		if reused {
			nodesReused++
			// Nœud réutilisé - vérifier la connexion au parent
			log.Printf("♻️  [AlphaChainBuilder] Réutilisation du nœud alpha %s pour la règle %s (condition %d/%d)",
				alphaNode.ID, ruleID, i+1, len(conditions))

			if !acb.isAlreadyConnectedCached(currentParent, alphaNode) {
				// Connecter au parent si pas déjà connecté
				currentParent.AddChild(alphaNode)
				log.Printf("🔗 [AlphaChainBuilder] Connexion du nœud réutilisé %s au parent %s",
					alphaNode.ID, currentParent.GetID())
			} else {
				log.Printf("✓  [AlphaChainBuilder] Nœud %s déjà connecté au parent %s",
					alphaNode.ID, currentParent.GetID())
			}
		} else {
			nodesCreated++
			// Nouveau nœud - le connecter au parent et l'ajouter au réseau
			currentParent.AddChild(alphaNode)
			acb.network.AlphaNodes[alphaNode.ID] = alphaNode

			// Mettre à jour le cache de connexion
			acb.updateConnectionCache(currentParent.GetID(), alphaNode.ID, true)

			log.Printf("🆕 [AlphaChainBuilder] Nouveau nœud alpha %s créé pour la règle %s (condition %d/%d)",
				alphaNode.ID, ruleID, i+1, len(conditions))
			log.Printf("🔗 [AlphaChainBuilder] Connexion du nœud %s au parent %s",
				alphaNode.ID, currentParent.GetID())
		}

		// Enregistrer le nœud dans le LifecycleManager avec la règle
		lifecycle := acb.network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
		lifecycle.AddRuleReference(ruleID, "") // RuleName peut être ajouté plus tard si nécessaire

		if reused {
			log.Printf("📊 [AlphaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)",
				alphaNode.ID, lifecycle.GetRefCount())
		}

		// Le nœud actuel devient le parent pour le prochain nœud
		currentParent = alphaNode
	}

	// Le dernier nœud de la chaîne est le nœud final
	chain.FinalNode = chain.Nodes[len(chain.Nodes)-1]

	log.Printf("✅ [AlphaChainBuilder] Chaîne alpha complète construite pour la règle %s: %d nœud(s)",
		ruleID, len(chain.Nodes))

	// Enregistrer les métriques
	if acb.metrics != nil {
		buildTime := time.Since(startTime)
		detail := ChainMetricDetail{
			RuleID:          ruleID,
			ChainLength:     len(chain.Nodes),
			NodesCreated:    nodesCreated,
			NodesReused:     nodesReused,
			BuildTime:       buildTime,
			Timestamp:       time.Now(),
			HashesGenerated: hashesGenerated,
		}
		acb.metrics.RecordChainBuild(detail)
	}

	return chain, nil
}

// BuildDecomposedChain constructs an alpha chain from decomposed conditions with full metadata.
// This method sets ResultName, IsAtomic, and Dependencies on each AlphaNode for
// context-aware evaluation with intermediate result propagation.
func (acb *AlphaChainBuilder) BuildDecomposedChain(
	conditions []DecomposedCondition,
	variableName string,
	parentNode Node,
	ruleID string,
) (*AlphaChain, error) {
	if len(conditions) == 0 {
		return nil, fmt.Errorf("impossible de construire une chaîne sans conditions")
	}

	if parentNode == nil {
		return nil, fmt.Errorf("le nœud parent ne peut pas être nil")
	}

	if acb.network.AlphaSharingManager == nil {
		return nil, fmt.Errorf("AlphaSharingManager non initialisé dans le réseau")
	}

	if acb.network.LifecycleManager == nil {
		return nil, fmt.Errorf("LifecycleManager non initialisé dans le réseau")
	}

	// Démarrer le chronomètre pour les métriques
	startTime := time.Now()
	nodesCreated := 0
	nodesReused := 0
	hashesGenerated := make([]string, 0, len(conditions))

	chain := &AlphaChain{
		Nodes:  make([]*AlphaNode, 0, len(conditions)),
		Hashes: make([]string, 0, len(conditions)),
		RuleID: ruleID,
	}

	currentParent := parentNode

	// Construire la chaîne condition par condition
	for i, decomposedCond := range conditions {
		// Convertir DecomposedCondition en map pour la condition du nœud alpha
		conditionMap := map[string]interface{}{
			"type":     decomposedCond.Type,
			"left":     decomposedCond.Left,
			"operator": decomposedCond.Operator,
			"right":    decomposedCond.Right,
		}

		// Obtenir ou créer l'AlphaNode via le gestionnaire de partage
		alphaNode, hash, reused, err := acb.network.AlphaSharingManager.GetOrCreateAlphaNode(
			conditionMap,
			variableName,
			acb.storage,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la création/récupération du nœud alpha %d: %w", i, err)
		}

		// SET DECOMPOSITION METADATA - This is the key enhancement
		alphaNode.ResultName = decomposedCond.ResultName
		alphaNode.IsAtomic = decomposedCond.IsAtomic
		alphaNode.Dependencies = decomposedCond.Dependencies

		// Ajouter le nœud et son hash à la chaîne
		chain.Nodes = append(chain.Nodes, alphaNode)
		chain.Hashes = append(chain.Hashes, hash)
		hashesGenerated = append(hashesGenerated, hash)

		if reused {
			nodesReused++
			// Nœud réutilisé - vérifier la connexion au parent
			log.Printf("♻️  [AlphaChainBuilder] Réutilisation du nœud alpha %s (decomposed: %s) pour la règle %s (condition %d/%d)",
				alphaNode.ID, alphaNode.ResultName, ruleID, i+1, len(conditions))

			if !acb.isAlreadyConnectedCached(currentParent, alphaNode) {
				// Connecter au parent si pas déjà connecté
				currentParent.AddChild(alphaNode)
				log.Printf("🔗 [AlphaChainBuilder] Connexion du nœud réutilisé %s au parent %s",
					alphaNode.ID, currentParent.GetID())
			} else {
				log.Printf("✓  [AlphaChainBuilder] Nœud %s déjà connecté au parent %s",
					alphaNode.ID, currentParent.GetID())
			}
		} else {
			nodesCreated++
			// Nouveau nœud - le connecter au parent et l'ajouter au réseau
			currentParent.AddChild(alphaNode)
			acb.network.AlphaNodes[alphaNode.ID] = alphaNode

			// Mettre à jour le cache de connexion
			acb.updateConnectionCache(currentParent.GetID(), alphaNode.ID, true)

			log.Printf("🆕 [AlphaChainBuilder] Nouveau nœud alpha %s créé (decomposed: %s, deps: %v) pour la règle %s (condition %d/%d)",
				alphaNode.ID, alphaNode.ResultName, alphaNode.Dependencies, ruleID, i+1, len(conditions))
			log.Printf("🔗 [AlphaChainBuilder] Connexion du nœud %s au parent %s",
				alphaNode.ID, currentParent.GetID())
		}

		// Enregistrer le nœud dans le LifecycleManager avec la règle
		lifecycle := acb.network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
		lifecycle.AddRuleReference(ruleID, "") // RuleName peut être ajouté plus tard si nécessaire

		if reused {
			log.Printf("📊 [AlphaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)",
				alphaNode.ID, lifecycle.GetRefCount())
		}

		// Le nœud actuel devient le parent pour le prochain nœud
		currentParent = alphaNode
	}

	// Le dernier nœud de la chaîne est le nœud final
	chain.FinalNode = chain.Nodes[len(chain.Nodes)-1]

	log.Printf("✅ [AlphaChainBuilder] Chaîne alpha décomposée complète construite pour la règle %s: %d nœud(s) atomiques",
		ruleID, len(chain.Nodes))

	// Enregistrer les métriques
	if acb.metrics != nil {
		buildTime := time.Since(startTime)
		detail := ChainMetricDetail{
			RuleID:          ruleID,
			ChainLength:     len(chain.Nodes),
			NodesCreated:    nodesCreated,
			NodesReused:     nodesReused,
			BuildTime:       buildTime,
			Timestamp:       time.Now(),
			HashesGenerated: hashesGenerated,
		}
		acb.metrics.RecordChainBuild(detail)
	}

	return chain, nil
}

// isAlreadyConnectedCached vérifie si un nœud enfant est déjà connecté à un nœud parent avec cache
func (acb *AlphaChainBuilder) isAlreadyConnectedCached(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	parentID := parent.GetID()
	childID := child.GetID()
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)

	// Vérifier le cache
	acb.mutex.RLock()
	if connected, exists := acb.connectionCache[cacheKey]; exists {
		acb.mutex.RUnlock()
		if acb.metrics != nil {
			acb.metrics.RecordConnectionCacheHit()
		}
		return connected
	}
	acb.mutex.RUnlock()

	// Cache miss - vérifier réellement
	if acb.metrics != nil {
		acb.metrics.RecordConnectionCacheMiss()
	}

	connected := isAlreadyConnected(parent, child)

	// Mettre à jour le cache
	acb.updateConnectionCache(parentID, childID, connected)

	return connected
}

// updateConnectionCache met à jour le cache de connexion
func (acb *AlphaChainBuilder) updateConnectionCache(parentID, childID string, connected bool) {
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)
	acb.mutex.Lock()
	acb.connectionCache[cacheKey] = connected
	acb.mutex.Unlock()
}

// ClearConnectionCache vide le cache de connexions.
//
// Utile pour libérer de la mémoire après suppression de nombreuses règles,
// ou pour forcer une réévaluation complète des connexions.
//
// Thread-safe: peut être appelé à tout moment.
//
// Exemple:
//
//	// Après suppression de beaucoup de règles
//	for _, ruleID := range oldRules {
//	    network.RemoveRule(ruleID)
//	}
//	builder.ClearConnectionCache() // Libérer mémoire
func (acb *AlphaChainBuilder) ClearConnectionCache() {
	acb.mutex.Lock()
	defer acb.mutex.Unlock()
	acb.connectionCache = make(map[string]bool)
}

// GetConnectionCacheSize retourne la taille actuelle du cache de connexions.
//
// Utile pour monitoring et diagnostic de l'utilisation mémoire.
//
// Retourne:
//   - Nombre d'entrées dans le cache (une par connexion parent→child unique)
//
// Exemple:
//
//	size := builder.GetConnectionCacheSize()
//	fmt.Printf("Cache de connexions: %d entrées\n", size)
//	if size > 10000 {
//	    builder.ClearConnectionCache() // Nettoyage si trop grand
//	}
func (acb *AlphaChainBuilder) GetConnectionCacheSize() int {
	acb.mutex.RLock()
	defer acb.mutex.RUnlock()
	return len(acb.connectionCache)
}

// GetMetrics retourne les métriques de performance du builder.
//
// Les métriques incluent:
//   - Nombre de chaînes construites
//   - Nœuds créés vs réutilisés
//   - Ratio de partage
//   - Hit rate du cache de hash
//   - Temps moyen de construction
//
// Retourne:
//   - Pointeur vers l'objet de métriques (non nil)
//
// Exemple:
//
//	metrics := builder.GetMetrics()
//	fmt.Printf("Chaînes construites: %d\n", metrics.TotalChainsBuilt)
//	fmt.Printf("Ratio de partage: %.1f%%\n", metrics.SharingRatio * 100)
//	fmt.Printf("Cache hit rate: %.1f%%\n",
//	    float64(metrics.HashCacheHits) /
//	    float64(metrics.HashCacheHits + metrics.HashCacheMisses) * 100)
func (acb *AlphaChainBuilder) GetMetrics() *ChainBuildMetrics {
	return acb.metrics
}

// isAlreadyConnected vérifie si un nœud enfant est déjà connecté à un nœud parent
func isAlreadyConnected(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	children := parent.GetChildren()
	childID := child.GetID()

	for _, c := range children {
		if c.GetID() == childID {
			return true
		}
	}

	return false
}

// GetChainInfo retourne des informations détaillées sur la chaîne alpha.
//
// Utile pour debugging, logging, et visualisation de la structure de la chaîne.
//
// Retourne:
//   - Map contenant: rule_id, node_count, node_ids, hashes, final_node_id
//   - Map avec clé "error" si chaîne nil
//
// Exemple:
//
//	info := chain.GetChainInfo()
//	fmt.Printf("Chaîne pour règle: %s\n", info["rule_id"])
//	fmt.Printf("Longueur: %d nœuds\n", info["node_count"])
//	fmt.Printf("IDs: %v\n", info["node_ids"])
//	fmt.Printf("Hashes: %v\n", info["hashes"])
func (ac *AlphaChain) GetChainInfo() map[string]interface{} {
	if ac == nil {
		return map[string]interface{}{
			"error": "chain is nil",
		}
	}

	nodeIDs := make([]string, len(ac.Nodes))
	for i, node := range ac.Nodes {
		nodeIDs[i] = node.ID
	}

	finalNodeID := ""
	if ac.FinalNode != nil {
		finalNodeID = ac.FinalNode.ID
	}

	return map[string]interface{}{
		"rule_id":       ac.RuleID,
		"node_count":    len(ac.Nodes),
		"node_ids":      nodeIDs,
		"hashes":        ac.Hashes,
		"final_node_id": finalNodeID,
	}
}

// ValidateChain vérifie que la chaîne alpha est valide et cohérente.
//
// Vérifie:
//   - Chaîne non nil
//   - Au moins un nœud présent
//   - len(Nodes) == len(Hashes)
//   - FinalNode correspond au dernier élément de Nodes
//   - Tous les nœuds ont un ID non vide
//
// Retourne:
//   - nil si chaîne valide
//   - error décrivant le problème si invalide
//
// Exemple:
//
//	chain, err := builder.BuildChain(...)
//	if err := chain.ValidateChain(); err != nil {
//	    log.Fatalf("Chaîne invalide: %v", err)
//	}
func (ac *AlphaChain) ValidateChain() error {
	if ac == nil {
		return fmt.Errorf("chaîne alpha nil")
	}

	if len(ac.Nodes) == 0 {
		return fmt.Errorf("chaîne alpha vide")
	}

	if len(ac.Nodes) != len(ac.Hashes) {
		return fmt.Errorf("incohérence: %d nœuds mais %d hashes", len(ac.Nodes), len(ac.Hashes))
	}

	if ac.FinalNode == nil {
		return fmt.Errorf("nœud final nil")
	}

	// Vérifier que le nœud final est bien le dernier de la liste
	if ac.FinalNode != ac.Nodes[len(ac.Nodes)-1] {
		return fmt.Errorf("le nœud final ne correspond pas au dernier nœud de la liste")
	}

	// Vérifier que tous les nœuds sont non-nil
	for i, node := range ac.Nodes {
		if node == nil {
			return fmt.Errorf("nœud %d est nil", i)
		}
	}

	return nil
}

// CountSharedNodes retourne le nombre de nœuds partagés dans la chaîne
// (nœuds avec plus d'une référence dans le LifecycleManager)
func (acb *AlphaChainBuilder) CountSharedNodes(chain *AlphaChain) int {
	if chain == nil || acb.network.LifecycleManager == nil {
		return 0
	}

	sharedCount := 0
	for _, node := range chain.Nodes {
		if lifecycle, exists := acb.network.LifecycleManager.GetNodeLifecycle(node.ID); exists {
			if lifecycle.GetRefCount() > 1 {
				sharedCount++
			}
		}
	}

	return sharedCount
}

// GetChainStats retourne des statistiques détaillées sur une chaîne alpha.
//
// Calcule et retourne:
//   - chain_length: Nombre total de nœuds dans la chaîne
//   - shared_nodes: Nœuds avec RefCount > 1
//   - new_nodes: Nœuds avec RefCount == 1
//   - sharing_ratio: Proportion de nœuds partagés (0.0 à 1.0)
//   - node_details: Liste des infos par nœud (ID, RefCount, is_shared)
//
// Paramètres:
//   - chain: Chaîne alpha à analyser
//
// Retourne:
//   - Map avec statistiques détaillées
//   - Map avec clé "error" si chaîne nil
//
// Exemple:
//
//	chain, _ := builder.BuildChain(...)
//	stats := builder.GetChainStats(chain)
//	fmt.Printf("Longueur: %d\n", stats["chain_length"])
//	fmt.Printf("Partagés: %d\n", stats["shared_nodes"])
//	fmt.Printf("Nouveaux: %d\n", stats["new_nodes"])
//	fmt.Printf("Ratio: %.1f%%\n", stats["sharing_ratio"].(float64) * 100)
//
//	// Détails par nœud
//	for _, detail := range stats["node_details"].([]map[string]interface{}) {
//	    fmt.Printf("  Nœud %s: RefCount=%d, Partagé=%v\n",
//	        detail["node_id"], detail["ref_count"], detail["is_shared"])
//	}
func (acb *AlphaChainBuilder) GetChainStats(chain *AlphaChain) map[string]interface{} {
	if chain == nil {
		return map[string]interface{}{
			"error": "chain is nil",
		}
	}

	sharedNodes := acb.CountSharedNodes(chain)
	newNodes := len(chain.Nodes) - sharedNodes

	stats := map[string]interface{}{
		"total_nodes":  len(chain.Nodes),
		"shared_nodes": sharedNodes,
		"new_nodes":    newNodes,
		"rule_id":      chain.RuleID,
	}

	// Ajouter les détails de chaque nœud
	nodeDetails := make([]map[string]interface{}, len(chain.Nodes))
	for i, node := range chain.Nodes {
		refCount := 0
		if lifecycle, exists := acb.network.LifecycleManager.GetNodeLifecycle(node.ID); exists {
			refCount = lifecycle.GetRefCount()
		}

		nodeDetails[i] = map[string]interface{}{
			"index":     i,
			"node_id":   node.ID,
			"hash":      chain.Hashes[i],
			"ref_count": refCount,
			"is_shared": refCount > 1,
			"is_final":  node == chain.FinalNode,
		}
	}
	stats["node_details"] = nodeDetails

	return stats
}
