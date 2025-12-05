// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// AggregationInfo contient les informations extraites d'une agrégation
type AggregationInfo struct {
	Function      string      // AVG, SUM, COUNT, MIN, MAX
	MainVariable  string      // Variable principale (ex: "e" pour Employee)
	MainType      string      // Type principal (ex: "Employee")
	AggVariable   string      // Variable à agréger (ex: "p" pour Performance)
	AggType       string      // Type à agréger (ex: "Performance")
	Field         string      // Champ à agréger (ex: "score")
	Operator      string      // Opérateur de comparaison (>=, >, etc.)
	Threshold     float64     // Valeur de seuil
	JoinField     string      // Champ de jointure dans faits agrégés (ex: "employee_id")
	MainField     string      // Champ de jointure dans fait principal (ex: "id")
	JoinCondition interface{} // Condition de jointure complète

	// Multi-source aggregation support
	AggregationVars []AggregationVariable // Multiple aggregation variables
	SourcePatterns  []SourcePattern       // Multiple source patterns to join
	JoinConditions  []JoinCondition       // Join conditions between patterns
}

// AggregationVariable represents a single aggregation variable
type AggregationVariable struct {
	Name      string  // Variable name (ex: "avg_sal")
	Function  string  // AVG, SUM, COUNT, MIN, MAX
	SourceVar string  // Source variable (ex: "e")
	Field     string  // Field to aggregate (ex: "salary")
	Operator  string  // Threshold operator (>=, >, etc.)
	Threshold float64 // Threshold value
}

// SourcePattern represents a pattern block in multi-source aggregation
type SourcePattern struct {
	Variable string // Variable name (ex: "e")
	Type     string // Type name (ex: "Employee")
}
