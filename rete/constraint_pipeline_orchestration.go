// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"github.com/treivax/tsd/constraint"
)

// ingestionContext encapsule l'état d'une ingestion de fichier.
// Cette structure est utilisée en interne par IngestFile pour maintenir
// l'état durant toutes les étapes du pipeline d'ingestion.
type ingestionContext struct {
	filename          string
	network           *ReteNetwork
	storage           Storage
	metrics           *MetricsCollector
	parsedAST         interface{}
	program           *constraint.Program
	reteProgram       interface{}
	types             []interface{}
	expressions       []interface{}
	factsForRete      []map[string]interface{}
	existingFacts     []*Fact
	factsByType       map[string][]*Fact
	existingTerminals map[string]bool
	newTerminals      []*TerminalNode
	hasResets         bool
	tx                *Transaction
}
