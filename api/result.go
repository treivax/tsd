// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"fmt"
	"time"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// Result contient le résultat d'une ingestion de programme TSD
type Result struct {
	network      *rete.ReteNetwork
	xupleManager xuples.XupleManager
	metrics      *Metrics
}

// Metrics contient les métriques d'ingestion
type Metrics struct {
	TotalDuration       time.Duration
	ParseDuration       time.Duration
	BuildDuration       time.Duration
	PropagationDuration time.Duration
	TypeCount           int
	RuleCount           int
	FactCount           int
	XupleSpaceCount     int
	PropagationCount    int
	ActionCount         int
}

// Network retourne le réseau RETE sous-jacent
func (r *Result) Network() *rete.ReteNetwork {
	return r.network
}

// XupleManager retourne le XupleManager
func (r *Result) XupleManager() xuples.XupleManager {
	return r.xupleManager
}

// Metrics retourne les métriques d'ingestion
func (r *Result) Metrics() *Metrics {
	return r.metrics
}

// TypeCount retourne le nombre de types définis
func (r *Result) TypeCount() int {
	return r.metrics.TypeCount
}

// RuleCount retourne le nombre de règles actives
func (r *Result) RuleCount() int {
	return r.metrics.RuleCount
}

// FactCount retourne le nombre de faits dans le réseau
func (r *Result) FactCount() int {
	return r.metrics.FactCount
}

// XupleSpaceCount retourne le nombre de xuple-spaces créés
func (r *Result) XupleSpaceCount() int {
	return r.metrics.XupleSpaceCount
}

// GetXuples retourne tous les xuples d'un xuple-space
func (r *Result) GetXuples(spaceName string) ([]*xuples.Xuple, error) {
	if r.xupleManager == nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "GetXuples",
			Message:   "XupleManager non initialisé",
		}
	}

	space, err := r.xupleManager.GetXupleSpace(spaceName)
	if err != nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "GetXuples",
			Message:   "xuple-space non trouvé",
			Cause:     err,
		}
	}

	return space.ListAll(), nil
}

// Retrieve récupère et consomme un xuple d'un xuple-space selon sa politique
func (r *Result) Retrieve(spaceName string, agentID string) (*xuples.Xuple, error) {
	if r.xupleManager == nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "Retrieve",
			Message:   "XupleManager non initialisé",
		}
	}

	space, err := r.xupleManager.GetXupleSpace(spaceName)
	if err != nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "Retrieve",
			Message:   "xuple-space non trouvé",
			Cause:     err,
		}
	}

	xuple, err := space.Retrieve(agentID)
	if err != nil {
		return nil, &XupleSpaceError{
			SpaceName: spaceName,
			Operation: "Retrieve",
			Message:   "échec de récupération",
			Cause:     err,
		}
	}

	return xuple, nil
}

// XupleSpaceNames retourne les noms de tous les xuple-spaces
func (r *Result) XupleSpaceNames() []string {
	if r.xupleManager == nil {
		return []string{}
	}
	return r.xupleManager.ListXupleSpaces()
}

// XupleCount retourne le nombre de xuples dans un xuple-space
func (r *Result) XupleCount(spaceName string) (int, error) {
	xuples, err := r.GetXuples(spaceName)
	if err != nil {
		return 0, err
	}
	return len(xuples), nil
}

// Summary retourne un résumé texte du résultat
func (r *Result) Summary() string {
	summary := fmt.Sprintf("=== Résultat d'Ingestion TSD ===\n")
	summary += fmt.Sprintf("Types définis:        %d\n", r.TypeCount())
	summary += fmt.Sprintf("Règles actives:       %d\n", r.RuleCount())
	summary += fmt.Sprintf("Faits dans réseau:    %d\n", r.FactCount())
	summary += fmt.Sprintf("Xuple-spaces créés:   %d\n", r.XupleSpaceCount())

	if r.xupleManager != nil {
		for _, spaceName := range r.XupleSpaceNames() {
			count, _ := r.XupleCount(spaceName)
			summary += fmt.Sprintf("  - %s: %d xuples\n", spaceName, count)
		}
	}

	summary += fmt.Sprintf("\nMétriques de Performance:\n")
	summary += fmt.Sprintf("Durée totale:         %v\n", r.metrics.TotalDuration)
	summary += fmt.Sprintf("  - Parsing:          %v\n", r.metrics.ParseDuration)
	summary += fmt.Sprintf("  - Construction:     %v\n", r.metrics.BuildDuration)
	summary += fmt.Sprintf("  - Propagation:      %v\n", r.metrics.PropagationDuration)
	summary += fmt.Sprintf("Propagations:         %d\n", r.metrics.PropagationCount)
	summary += fmt.Sprintf("Actions exécutées:    %d\n", r.metrics.ActionCount)

	return summary
}
