// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// CircularDependencyDetector détecte les cycles dans les graphes de dépendances
// des expressions arithmétiques décomposées. Utilise un algorithme DFS avec
// marquage tricolore (blanc/gris/noir) pour détecter les cycles.
type CircularDependencyDetector struct {
	// graph représente le graphe de dépendances : resultName -> liste de dépendances
	graph map[string][]string

	// colors stocke l'état de visite de chaque nœud (white/gray/black)
	colors map[string]nodeColor

	// parent stocke le parent de chaque nœud dans le parcours DFS
	parent map[string]string

	// cyclePath stocke le chemin du cycle s'il est détecté
	cyclePath []string

	// metadata stocke des informations supplémentaires sur chaque nœud
	metadata map[string]*NodeMetadata
}

// nodeColor représente l'état de visite d'un nœud pendant le DFS
type nodeColor int

const (
	white nodeColor = iota // Non visité
	gray                   // En cours de visite (dans la pile DFS)
	black                  // Visite terminée
)

// NodeMetadata contient des informations supplémentaires sur un nœud
type NodeMetadata struct {
	ResultName  string
	Depth       int      // Profondeur dans le graphe
	Expression  string   // Expression textuelle (pour debug)
	IsAtomic    bool     // Si le nœud représente une opération atomique
	Descendants []string // Tous les descendants (transitivement)
}

// ValidationReport contient le résultat de la validation
type ValidationReport struct {
	Valid           bool
	HasCircularDeps bool
	CyclePath       []string
	MaxDepth        int
	TotalNodes      int
	IsolatedNodes   []string // Nœuds sans dépendances ni dépendants
	ErrorMessage    string
	Warnings        []string
}

// NewCircularDependencyDetector crée un nouveau détecteur
func NewCircularDependencyDetector() *CircularDependencyDetector {
	return &CircularDependencyDetector{
		graph:     make(map[string][]string),
		colors:    make(map[string]nodeColor),
		parent:    make(map[string]string),
		cyclePath: make([]string, 0),
		metadata:  make(map[string]*NodeMetadata),
	}
}

// AddNode ajoute un nœud au graphe avec ses dépendances
func (cdd *CircularDependencyDetector) AddNode(resultName string, dependencies []string) {
	if dependencies == nil {
		dependencies = make([]string, 0)
	}
	cdd.graph[resultName] = dependencies
	cdd.colors[resultName] = white

	// Initialiser metadata si elle n'existe pas
	if _, exists := cdd.metadata[resultName]; !exists {
		cdd.metadata[resultName] = &NodeMetadata{
			ResultName:  resultName,
			IsAtomic:    true,
			Descendants: make([]string, 0),
		}
	}
}

// AddNodeWithMetadata ajoute un nœud avec des métadonnées complètes
func (cdd *CircularDependencyDetector) AddNodeWithMetadata(resultName string, dependencies []string, metadata *NodeMetadata) {
	cdd.AddNode(resultName, dependencies)
	cdd.metadata[resultName] = metadata
}

// DetectCycles détecte s'il y a des cycles dans le graphe
// Retourne true si un cycle est détecté
func (cdd *CircularDependencyDetector) DetectCycles() bool {
	// Réinitialiser les couleurs et parents
	for node := range cdd.graph {
		cdd.colors[node] = white
		cdd.parent[node] = ""
	}
	cdd.cyclePath = make([]string, 0)

	// Lancer DFS depuis chaque nœud non visité
	for node := range cdd.graph {
		if cdd.colors[node] == white {
			if cdd.dfs(node) {
				return true // Cycle détecté
			}
		}
	}

	return false // Pas de cycle
}

// dfs effectue un parcours en profondeur depuis le nœud donné
// Retourne true si un cycle est détecté
func (cdd *CircularDependencyDetector) dfs(node string) bool {
	// Marquer comme en cours de visite
	cdd.colors[node] = gray

	// Explorer les dépendances
	for _, dep := range cdd.graph[node] {
		if cdd.colors[dep] == white {
			// Nœud non visité, continuer DFS
			cdd.parent[dep] = node
			if cdd.dfs(dep) {
				return true
			}
		} else if cdd.colors[dep] == gray {
			// Nœud gris : cycle détecté !
			cdd.buildCyclePath(node, dep)
			return true
		}
		// Si noir, déjà visité complètement, pas de cycle
	}

	// Marquer comme visité complètement
	cdd.colors[node] = black
	return false
}

// buildCyclePath reconstruit le chemin du cycle
func (cdd *CircularDependencyDetector) buildCyclePath(from, to string) {
	cdd.cyclePath = []string{to}

	// Remonter depuis from jusqu'à to
	current := from
	for current != "" && current != to {
		cdd.cyclePath = append([]string{current}, cdd.cyclePath...)
		current = cdd.parent[current]
	}

	// Ajouter to à la fin pour fermer le cycle
	cdd.cyclePath = append(cdd.cyclePath, to)
}

// GetCyclePath retourne le chemin du cycle détecté
func (cdd *CircularDependencyDetector) GetCyclePath() []string {
	return cdd.cyclePath
}

// Validate effectue une validation complète du graphe et retourne un rapport détaillé
func (cdd *CircularDependencyDetector) Validate() ValidationReport {
	report := ValidationReport{
		Valid:           true,
		HasCircularDeps: false,
		MaxDepth:        0,
		TotalNodes:      len(cdd.graph),
		IsolatedNodes:   make([]string, 0),
		Warnings:        make([]string, 0),
	}

	// Détecter les cycles
	if cdd.DetectCycles() {
		report.Valid = false
		report.HasCircularDeps = true
		report.CyclePath = cdd.GetCyclePath()
		report.ErrorMessage = fmt.Sprintf("Circular dependency detected: %s", cdd.FormatCyclePath())
		return report
	}

	// Calculer la profondeur maximale
	report.MaxDepth = cdd.calculateMaxDepth()

	// Détecter les nœuds isolés
	report.IsolatedNodes = cdd.findIsolatedNodes()
	if len(report.IsolatedNodes) > 0 {
		report.Warnings = append(report.Warnings,
			fmt.Sprintf("Found %d isolated nodes: %v", len(report.IsolatedNodes), report.IsolatedNodes))
	}

	// Détecter les dépendances manquantes
	missing := cdd.findMissingDependencies()
	if len(missing) > 0 {
		report.Valid = false
		report.ErrorMessage = fmt.Sprintf("Missing dependencies: %v", missing)
		return report
	}

	// Warnings pour profondeur excessive
	if report.MaxDepth > 10 {
		report.Warnings = append(report.Warnings,
			fmt.Sprintf("Maximum depth is %d (>10), may impact performance", report.MaxDepth))
	}

	return report
}

// calculateMaxDepth calcule la profondeur maximale du graphe
func (cdd *CircularDependencyDetector) calculateMaxDepth() int {
	depths := make(map[string]int)

	// Initialiser les profondeurs à -1
	for node := range cdd.graph {
		depths[node] = -1
	}

	var calculateDepth func(string) int
	calculateDepth = func(node string) int {
		if depths[node] != -1 {
			return depths[node]
		}

		maxChildDepth := 0
		for _, dep := range cdd.graph[node] {
			childDepth := calculateDepth(dep)
			if childDepth+1 > maxChildDepth {
				maxChildDepth = childDepth + 1
			}
		}

		depths[node] = maxChildDepth
		if cdd.metadata[node] != nil {
			cdd.metadata[node].Depth = maxChildDepth
		}
		return maxChildDepth
	}

	maxDepth := 0
	for node := range cdd.graph {
		depth := calculateDepth(node)
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

// findIsolatedNodes trouve les nœuds sans dépendances ni dépendants
func (cdd *CircularDependencyDetector) findIsolatedNodes() []string {
	hasDependents := make(map[string]bool)

	// Marquer les nœuds qui sont dépendances d'autres nœuds
	for _, deps := range cdd.graph {
		for _, dep := range deps {
			hasDependents[dep] = true
		}
	}

	isolated := make([]string, 0)
	for node, deps := range cdd.graph {
		if len(deps) == 0 && !hasDependents[node] {
			isolated = append(isolated, node)
		}
	}

	return isolated
}

// findMissingDependencies trouve les dépendances référencées mais non définies
func (cdd *CircularDependencyDetector) findMissingDependencies() []string {
	missing := make([]string, 0)
	seen := make(map[string]bool)

	for _, deps := range cdd.graph {
		for _, dep := range deps {
			if _, exists := cdd.graph[dep]; !exists {
				if !seen[dep] {
					missing = append(missing, dep)
					seen[dep] = true
				}
			}
		}
	}

	return missing
}

// FormatCyclePath retourne une représentation formatée du chemin du cycle
func (cdd *CircularDependencyDetector) FormatCyclePath() string {
	if len(cdd.cyclePath) == 0 {
		return "(no cycle)"
	}
	return strings.Join(cdd.cyclePath, " → ")
}

// GetDependencyChain retourne la chaîne complète de dépendances pour un nœud donné
func (cdd *CircularDependencyDetector) GetDependencyChain(node string) []string {
	if _, exists := cdd.graph[node]; !exists {
		return []string{}
	}

	visited := make(map[string]bool)
	chain := make([]string, 0)

	var buildChain func(string)
	buildChain = func(n string) {
		if visited[n] {
			return
		}
		visited[n] = true
		chain = append(chain, n)

		for _, dep := range cdd.graph[n] {
			buildChain(dep)
		}
	}

	buildChain(node)
	return chain
}

// GetTopologicalSort retourne un tri topologique du graphe (si pas de cycle)
func (cdd *CircularDependencyDetector) GetTopologicalSort() ([]string, error) {
	if cdd.DetectCycles() {
		return nil, fmt.Errorf("cannot perform topological sort: graph contains cycles")
	}

	visited := make(map[string]bool)
	result := make([]string, 0, len(cdd.graph))

	var visit func(string)
	visit = func(node string) {
		if visited[node] {
			return
		}
		visited[node] = true

		// Visiter les dépendances d'abord (post-order)
		for _, dep := range cdd.graph[node] {
			visit(dep)
		}

		result = append(result, node)
	}

	// Visiter tous les nœuds
	for node := range cdd.graph {
		visit(node)
	}

	return result, nil
}

// GetStatistics retourne des statistiques sur le graphe
func (cdd *CircularDependencyDetector) GetStatistics() map[string]interface{} {
	totalEdges := 0
	for _, deps := range cdd.graph {
		totalEdges += len(deps)
	}

	return map[string]interface{}{
		"total_nodes":       len(cdd.graph),
		"total_edges":       totalEdges,
		"average_outdegree": float64(totalEdges) / float64(len(cdd.graph)),
		"max_depth":         cdd.calculateMaxDepth(),
		"has_cycles":        cdd.DetectCycles(),
	}
}

// Clear réinitialise le détecteur
func (cdd *CircularDependencyDetector) Clear() {
	cdd.graph = make(map[string][]string)
	cdd.colors = make(map[string]nodeColor)
	cdd.parent = make(map[string]string)
	cdd.cyclePath = make([]string, 0)
	cdd.metadata = make(map[string]*NodeMetadata)
}

// ValidateDecomposedConditions valide un slice de conditions décomposées
func (cdd *CircularDependencyDetector) ValidateDecomposedConditions(decomposed []DecomposedCondition) error {
	// Construire le graphe depuis les conditions décomposées
	for _, step := range decomposed {
		cdd.AddNode(step.ResultName, step.Dependencies)
	}

	// Valider
	report := cdd.Validate()
	if !report.Valid {
		return fmt.Errorf("validation failed: %s", report.ErrorMessage)
	}

	if report.HasCircularDeps {
		return fmt.Errorf("circular dependency detected: %s", cdd.FormatCyclePath())
	}

	return nil
}

// ValidateAlphaChain valide une chaîne d'alpha nodes
func (cdd *CircularDependencyDetector) ValidateAlphaChain(nodes []*AlphaNode) ValidationReport {
	// Construire le graphe depuis les alpha nodes
	for _, node := range nodes {
		if node.ResultName != "" {
			cdd.AddNode(node.ResultName, node.Dependencies)
		}
	}

	return cdd.Validate()
}

// String retourne une représentation textuelle du graphe
func (cdd *CircularDependencyDetector) String() string {
	var sb strings.Builder
	sb.WriteString("Dependency Graph:\n")

	for node, deps := range cdd.graph {
		sb.WriteString(fmt.Sprintf("  %s -> [%s]\n", node, strings.Join(deps, ", ")))
	}

	if len(cdd.cyclePath) > 0 {
		sb.WriteString(fmt.Sprintf("\nCycle detected: %s\n", cdd.FormatCyclePath()))
	}

	return sb.String()
}
