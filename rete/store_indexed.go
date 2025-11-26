// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"
)

// IndexedFactStorage fournit un stockage indexé pour les faits
type IndexedFactStorage struct {
	// Index principal par ID de fait
	factsByID map[string]*Fact

	// Index par type de fait
	factsByType map[string]map[string]*Fact

	// Index par propriétés de fait (champ -> valeur -> faits)
	factsByField map[string]map[interface{}]map[string]*Fact

	// Index composite pour les jointures fréquentes
	compositeIndex map[string]map[string]*Fact

	// Statistiques d'utilisation pour optimiser les index
	accessStats map[string]int64

	// Verrou pour la concurrence
	mutex sync.RWMutex

	// Configuration des index
	config IndexConfig
}

// IndexConfig configure les options d'indexation
type IndexConfig struct {
	// Champs à indexer automatiquement
	IndexedFields []string

	// Taille maximale du cache
	MaxCacheSize int

	// TTL pour les entrées de cache
	CacheTTL time.Duration

	// Activer les index composites
	EnableCompositeIndex bool

	// Seuil pour créer des index automatiques
	AutoIndexThreshold int64
}

// NewIndexedFactStorage crée un nouveau stockage indexé
func NewIndexedFactStorage(config IndexConfig) *IndexedFactStorage {
	return &IndexedFactStorage{
		factsByID:      make(map[string]*Fact),
		factsByType:    make(map[string]map[string]*Fact),
		factsByField:   make(map[string]map[interface{}]map[string]*Fact),
		compositeIndex: make(map[string]map[string]*Fact),
		accessStats:    make(map[string]int64),
		config:         config,
	}
}

// StoreFact stocke un fait avec indexation automatique
func (ifs *IndexedFactStorage) StoreFact(fact *Fact) error {
	ifs.mutex.Lock()
	defer ifs.mutex.Unlock()

	// Stocker dans l'index principal
	ifs.factsByID[fact.ID] = fact

	// Indexer par type
	if ifs.factsByType[fact.Type] == nil {
		ifs.factsByType[fact.Type] = make(map[string]*Fact)
	}
	ifs.factsByType[fact.Type][fact.ID] = fact

	// Indexer par champs configurés
	for _, fieldName := range ifs.config.IndexedFields {
		if value, exists := fact.Fields[fieldName]; exists {
			ifs.indexFieldValue(fieldName, value, fact)
		}
	}

	// Indexer par tous les champs si activé
	for fieldName, value := range fact.Fields {
		ifs.indexFieldValue(fieldName, value, fact)
	}

	// Créer des index composites si activé
	if ifs.config.EnableCompositeIndex {
		ifs.createCompositeIndexes(fact)
	}

	return nil
}

// indexFieldValue indexe une valeur de champ
func (ifs *IndexedFactStorage) indexFieldValue(fieldName string, value interface{}, fact *Fact) {
	if ifs.factsByField[fieldName] == nil {
		ifs.factsByField[fieldName] = make(map[interface{}]map[string]*Fact)
	}
	if ifs.factsByField[fieldName][value] == nil {
		ifs.factsByField[fieldName][value] = make(map[string]*Fact)
	}
	ifs.factsByField[fieldName][value][fact.ID] = fact
}

// createCompositeIndexes crée des index composites pour les jointures fréquentes
func (ifs *IndexedFactStorage) createCompositeIndexes(fact *Fact) {
	// Créer des clés composites pour des combinaisons communes
	if id, hasID := fact.Fields["id"]; hasID {
		if name, hasName := fact.Fields["name"]; hasName {
			compositeKey := fmt.Sprintf("id_name:%v_%v", id, name)
			if ifs.compositeIndex[compositeKey] == nil {
				ifs.compositeIndex[compositeKey] = make(map[string]*Fact)
			}
			ifs.compositeIndex[compositeKey][fact.ID] = fact
		}
	}
}

// GetFactByID récupère un fait par son ID
func (ifs *IndexedFactStorage) GetFactByID(id string) (*Fact, bool) {
	ifs.mutex.RLock()
	defer ifs.mutex.RUnlock()

	ifs.recordAccess("id:" + id)

	fact, exists := ifs.factsByID[id]
	return fact, exists
}

// GetFactsByType récupère tous les faits d'un type donné
func (ifs *IndexedFactStorage) GetFactsByType(factType string) []*Fact {
	ifs.mutex.RLock()
	defer ifs.mutex.RUnlock()

	ifs.recordAccess("type:" + factType)

	factsMap, exists := ifs.factsByType[factType]
	if !exists {
		return []*Fact{}
	}

	facts := make([]*Fact, 0, len(factsMap))
	for _, fact := range factsMap {
		facts = append(facts, fact)
	}

	return facts
}

// GetFactsByField récupère des faits par valeur de champ
func (ifs *IndexedFactStorage) GetFactsByField(fieldName string, value interface{}) []*Fact {
	ifs.mutex.RLock()
	defer ifs.mutex.RUnlock()

	accessKey := fmt.Sprintf("field:%s:%v", fieldName, value)
	ifs.recordAccess(accessKey)

	fieldIndex, exists := ifs.factsByField[fieldName]
	if !exists {
		return []*Fact{}
	}

	factsMap, exists := fieldIndex[value]
	if !exists {
		return []*Fact{}
	}

	facts := make([]*Fact, 0, len(factsMap))
	for _, fact := range factsMap {
		facts = append(facts, fact)
	}

	return facts
}

// GetFactsByCompositeKey récupère des faits par clé composite
func (ifs *IndexedFactStorage) GetFactsByCompositeKey(key string) []*Fact {
	ifs.mutex.RLock()
	defer ifs.mutex.RUnlock()

	ifs.recordAccess("composite:" + key)

	factsMap, exists := ifs.compositeIndex[key]
	if !exists {
		return []*Fact{}
	}

	facts := make([]*Fact, 0, len(factsMap))
	for _, fact := range factsMap {
		facts = append(facts, fact)
	}

	return facts
}

// RemoveFact supprime un fait et met à jour les index
func (ifs *IndexedFactStorage) RemoveFact(factID string) bool {
	ifs.mutex.Lock()
	defer ifs.mutex.Unlock()

	fact, exists := ifs.factsByID[factID]
	if !exists {
		return false
	}

	// Supprimer de l'index principal
	delete(ifs.factsByID, factID)

	// Supprimer de l'index par type
	if typeMap := ifs.factsByType[fact.Type]; typeMap != nil {
		delete(typeMap, factID)
		if len(typeMap) == 0 {
			delete(ifs.factsByType, fact.Type)
		}
	}

	// Supprimer des index par champ
	for fieldName, value := range fact.Fields {
		if fieldIndex := ifs.factsByField[fieldName]; fieldIndex != nil {
			if valueMap := fieldIndex[value]; valueMap != nil {
				delete(valueMap, factID)
				if len(valueMap) == 0 {
					delete(fieldIndex, value)
					if len(fieldIndex) == 0 {
						delete(ifs.factsByField, fieldName)
					}
				}
			}
		}
	}

	// Supprimer des index composites
	ifs.removeFromCompositeIndexes(fact)

	return true
}

// removeFromCompositeIndexes supprime le fait des index composites
func (ifs *IndexedFactStorage) removeFromCompositeIndexes(fact *Fact) {
	// Supprimer des clés composites
	if id, hasID := fact.Fields["id"]; hasID {
		if name, hasName := fact.Fields["name"]; hasName {
			compositeKey := fmt.Sprintf("id_name:%v_%v", id, name)
			if compositeMap := ifs.compositeIndex[compositeKey]; compositeMap != nil {
				delete(compositeMap, fact.ID)
				if len(compositeMap) == 0 {
					delete(ifs.compositeIndex, compositeKey)
				}
			}
		}
	}
}

// recordAccess enregistre un accès pour les statistiques
func (ifs *IndexedFactStorage) recordAccess(key string) {
	ifs.accessStats[key]++

	// Créer automatiquement des index pour les accès fréquents
	if ifs.accessStats[key] > ifs.config.AutoIndexThreshold {
		// Logique pour créer des index automatiques basés sur les patterns d'accès
	}
}

// GetAccessStats retourne les statistiques d'accès
func (ifs *IndexedFactStorage) GetAccessStats() map[string]int64 {
	ifs.mutex.RLock()
	defer ifs.mutex.RUnlock()

	stats := make(map[string]int64)
	for key, count := range ifs.accessStats {
		stats[key] = count
	}

	return stats
}

// OptimizeIndexes optimise les index basés sur les statistiques d'usage
func (ifs *IndexedFactStorage) OptimizeIndexes() {
	ifs.mutex.Lock()
	defer ifs.mutex.Unlock()

	// Analyser les patterns d'accès fréquents
	for _, count := range ifs.accessStats {
		if count > ifs.config.AutoIndexThreshold {
			// Créer des index optimisés pour cet accès
			// Cette logique peut être étendue selon les besoins
		}
	}
}

// Clear vide tous les index et statistiques
func (ifs *IndexedFactStorage) Clear() {
	ifs.mutex.Lock()
	defer ifs.mutex.Unlock()

	ifs.factsByID = make(map[string]*Fact)
	ifs.factsByType = make(map[string]map[string]*Fact)
	ifs.factsByField = make(map[string]map[interface{}]map[string]*Fact)
	ifs.compositeIndex = make(map[string]map[string]*Fact)
	ifs.accessStats = make(map[string]int64)
}

// Size retourne le nombre total de faits stockés
func (ifs *IndexedFactStorage) Size() int {
	ifs.mutex.RLock()
	defer ifs.mutex.RUnlock()

	return len(ifs.factsByID)
}
