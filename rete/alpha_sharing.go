// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

// AlphaSharingRegistry gère le partage des AlphaNodes entre plusieurs règles
// qui ont des conditions alpha identiques
type AlphaSharingRegistry struct {
	// Map[ConditionHash] -> AlphaNode
	sharedAlphaNodes map[string]*AlphaNode
	mutex            sync.RWMutex
}

// NewAlphaSharingRegistry crée un nouveau registre de partage d'AlphaNodes
func NewAlphaSharingRegistry() *AlphaSharingRegistry {
	return &AlphaSharingRegistry{
		sharedAlphaNodes: make(map[string]*AlphaNode),
	}
}

// ConditionHash calcule un hash unique pour une condition alpha
// Deux conditions identiques produiront le même hash
func ConditionHash(condition interface{}, variableName string) (string, error) {
	// Normaliser la condition pour assurer un hash cohérent
	normalized, err := normalizeCondition(condition)
	if err != nil {
		return "", fmt.Errorf("erreur normalisation condition: %w", err)
	}

	// Créer une structure canonique incluant le nom de variable
	// car une condition sur "p" vs "q" est différente
	canonical := map[string]interface{}{
		"condition": normalized,
		"variable":  variableName,
	}

	// Sérialiser en JSON avec clés triées pour cohérence
	jsonBytes, err := json.Marshal(canonical)
	if err != nil {
		return "", fmt.Errorf("erreur sérialisation condition: %w", err)
	}

	// Calculer le hash SHA-256
	hash := sha256.Sum256(jsonBytes)
	return fmt.Sprintf("alpha_%x", hash[:8]), nil // Utiliser les 8 premiers octets pour l'ID
}

// normalizeCondition normalise une condition pour assurer un hash cohérent
func normalizeCondition(condition interface{}) (interface{}, error) {
	if condition == nil {
		return map[string]interface{}{"type": "simple"}, nil
	}

	switch cond := condition.(type) {
	case map[string]interface{}:
		// Copier la map et normaliser récursivement
		normalized := make(map[string]interface{})
		for key, value := range cond {
			normalizedValue, err := normalizeCondition(value)
			if err != nil {
				return nil, err
			}
			normalized[key] = normalizedValue
		}
		return normalized, nil

	case []interface{}:
		// Normaliser chaque élément du slice
		normalized := make([]interface{}, len(cond))
		for i, item := range cond {
			normalizedItem, err := normalizeCondition(item)
			if err != nil {
				return nil, err
			}
			normalized[i] = normalizedItem
		}
		return normalized, nil

	default:
		// Types primitifs: retourner tel quel
		return condition, nil
	}
}

// GetOrCreateAlphaNode récupère un AlphaNode existant ou en crée un nouveau
// Si un AlphaNode avec la même condition existe, il est réutilisé
func (asr *AlphaSharingRegistry) GetOrCreateAlphaNode(
	condition interface{},
	variableName string,
	storage Storage,
) (*AlphaNode, string, bool, error) {
	// Calculer le hash de la condition
	hash, err := ConditionHash(condition, variableName)
	if err != nil {
		return nil, "", false, fmt.Errorf("erreur calcul hash condition: %w", err)
	}

	// Vérifier si un AlphaNode existe déjà pour cette condition
	asr.mutex.RLock()
	existingNode, exists := asr.sharedAlphaNodes[hash]
	asr.mutex.RUnlock()

	if exists {
		// AlphaNode partagé trouvé
		return existingNode, hash, true, nil
	}

	// Créer un nouveau AlphaNode
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	// Double-check après avoir acquis le verrou d'écriture
	if existingNode, exists := asr.sharedAlphaNodes[hash]; exists {
		return existingNode, hash, true, nil
	}

	// Créer le nouveau nœud avec l'ID basé sur le hash
	alphaNode := NewAlphaNode(hash, condition, variableName, storage)
	asr.sharedAlphaNodes[hash] = alphaNode

	return alphaNode, hash, false, nil
}

// GetAlphaNode récupère un AlphaNode partagé par son hash
func (asr *AlphaSharingRegistry) GetAlphaNode(hash string) (*AlphaNode, bool) {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	node, exists := asr.sharedAlphaNodes[hash]
	return node, exists
}

// RemoveAlphaNode supprime un AlphaNode du registre
// Cette méthode doit être appelée uniquement quand plus aucune règle n'utilise ce nœud
func (asr *AlphaSharingRegistry) RemoveAlphaNode(hash string) error {
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	if _, exists := asr.sharedAlphaNodes[hash]; !exists {
		return fmt.Errorf("AlphaNode %s non trouvé dans le registre", hash)
	}

	delete(asr.sharedAlphaNodes, hash)
	return nil
}

// GetStats retourne des statistiques sur le partage des AlphaNodes
func (asr *AlphaSharingRegistry) GetStats() map[string]interface{} {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	totalNodes := len(asr.sharedAlphaNodes)

	// Compter le nombre total de règles utilisant ces nœuds
	totalRuleReferences := 0
	childCounts := make([]int, 0, totalNodes)

	for _, node := range asr.sharedAlphaNodes {
		childCount := len(node.GetChildren())
		childCounts = append(childCounts, childCount)
		totalRuleReferences += childCount
	}

	// Calculer la moyenne
	avgSharing := 0.0
	if totalNodes > 0 {
		avgSharing = float64(totalRuleReferences) / float64(totalNodes)
	}

	return map[string]interface{}{
		"total_shared_alpha_nodes": totalNodes,
		"total_rule_references":    totalRuleReferences,
		"average_sharing_ratio":    avgSharing,
	}
}

// ListSharedAlphaNodes retourne la liste de tous les AlphaNodes partagés
func (asr *AlphaSharingRegistry) ListSharedAlphaNodes() []string {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	hashes := make([]string, 0, len(asr.sharedAlphaNodes))
	for hash := range asr.sharedAlphaNodes {
		hashes = append(hashes, hash)
	}

	// Trier pour avoir un ordre déterministe
	sort.Strings(hashes)
	return hashes
}

// Reset réinitialise complètement le registre
func (asr *AlphaSharingRegistry) Reset() {
	asr.mutex.Lock()
	defer asr.mutex.Unlock()

	asr.sharedAlphaNodes = make(map[string]*AlphaNode)
}

// GetSharedAlphaNodeDetails retourne les détails d'un AlphaNode partagé
func (asr *AlphaSharingRegistry) GetSharedAlphaNodeDetails(hash string) map[string]interface{} {
	asr.mutex.RLock()
	defer asr.mutex.RUnlock()

	node, exists := asr.sharedAlphaNodes[hash]
	if !exists {
		return nil
	}

	children := node.GetChildren()
	childIDs := make([]string, len(children))
	for i, child := range children {
		childIDs[i] = child.GetID()
	}

	return map[string]interface{}{
		"hash":          hash,
		"node_id":       node.GetID(),
		"variable_name": node.VariableName,
		"condition":     node.Condition,
		"child_count":   len(children),
		"child_ids":     childIDs,
	}
}
