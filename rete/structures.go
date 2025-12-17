// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// Field represents a single field within a type definition.
// It contains the field name, its type, and whether it's part of the primary key.
type Field struct {
	Name         string `json:"name"`                   // Field name (e.g., "id", "name")
	Type         string `json:"type"`                   // Field type (e.g., "string", "number", "bool")
	IsPrimaryKey bool   `json:"isPrimaryKey,omitempty"` // True if field is part of primary key (marked with #)
}

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ActionDefinition struct {
	Type       string      `json:"type"`       // Always "actionDefinition"
	Name       string      `json:"name"`       // The action name (e.g., "notify")
	Parameters []Parameter `json:"parameters"` // List of parameters for the action
}

type TypeDefinition struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type JobCall struct {
	Type string        `json:"type"`
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}

type Action struct {
	Type string    `json:"type"`
	Job  *JobCall  `json:"job,omitempty"`  // Single job (backward compatibility)
	Jobs []JobCall `json:"jobs,omitempty"` // Multiple jobs (new format)
}

// GetJobs returns the list of jobs to execute.
// It handles both the old format (single Job) and new format (multiple Jobs).
func (a *Action) GetJobs() []JobCall {
	if len(a.Jobs) > 0 {
		return a.Jobs
	}
	if a.Job != nil {
		return []JobCall{*a.Job}
	}
	return []JobCall{}
}

type TypedVariable struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type Set struct {
	Type      string          `json:"type"`
	Variables []TypedVariable `json:"variables"`
}

type Expression struct {
	Type        string      `json:"type"`
	Set         Set         `json:"set"`
	Constraints interface{} `json:"constraints"`
	Action      *Action     `json:"action,omitempty"`
}

type Program struct {
	Types       []TypeDefinition `json:"types"`
	Expressions []Expression     `json:"expressions"`
}

// Clone crée une copie profonde de TypeDefinition.
// Tous les champs incluant IsPrimaryKey sont copiés.
func (td TypeDefinition) Clone() TypeDefinition {
	clone := TypeDefinition{
		Type:   td.Type,
		Name:   td.Name,
		Fields: make([]Field, len(td.Fields)),
	}

	// Copier les champs (copy() copie tous les champs de la struct incluant IsPrimaryKey)
	copy(clone.Fields, td.Fields)

	return clone
}

// Clone crée une copie profonde d Action
func (a *Action) Clone() *Action {
	clone := &Action{
		Type: a.Type,
		Jobs: make([]JobCall, len(a.Jobs)),
	}

	// Copier les jobs
	copy(clone.Jobs, a.Jobs)

	// Copier le job unique si présent
	if a.Job != nil {
		jobCopy := *a.Job
		clone.Job = &jobCopy
	}

	return clone
}
