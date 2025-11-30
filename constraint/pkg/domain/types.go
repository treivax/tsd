// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// Program représente un programme complet de contraintes
type Program struct {
	Types       []TypeDefinition `json:"types"`
	Expressions []Expression     `json:"expressions"`
	Metadata    *Metadata        `json:"metadata,omitempty"`
}

// Metadata contient les informations sur le programme
type Metadata struct {
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	Author      string    `json:"author,omitempty"`
	Description string    `json:"description,omitempty"`
}

// TypeDefinition définit un type personnalisé avec ses champs
type TypeDefinition struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

// Field représente un champ dans une définition de type
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Expression représente une règle métier complète
// Chaque expression doit avoir un identifiant unique pour la gestion (ex: suppression)
type Expression struct {
	Type        string      `json:"type"`
	RuleId      string      `json:"ruleId"`             // Identifiant unique de la règle
	Set         Set         `json:"set,omitempty"`      // Set of variables (single pattern, backward compatibility)
	Patterns    []Set       `json:"patterns,omitempty"` // Multiple pattern blocks (aggregation with joins)
	Constraints interface{} `json:"constraints"`
	Action      *Action     `json:"action,omitempty"`
}

// Set représente l'ensemble de variables typées dans une règle
type Set struct {
	Type      string          `json:"type"`
	Variables []TypedVariable `json:"variables"`
}

// TypedVariable représente une variable avec son type
type TypedVariable struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

// Constraint représente une contrainte (condition)
type Constraint struct {
	Type     string      `json:"type"`
	Left     interface{} `json:"left,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Right    interface{} `json:"right,omitempty"`
}

// Action représente une action à exécuter quand les conditions sont remplies
type Action struct {
	Type string  `json:"type"`
	Job  JobCall `json:"job"`
}

// JobCall représente l'appel d'une fonction/job
type JobCall struct {
	Type string        `json:"type"`
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}

// FieldAccess représente l'accès à un champ d'une variable
type FieldAccess struct {
	Type   string `json:"type"`
	Object string `json:"object"`
	Field  string `json:"field"`
}

// Variable représente une référence à une variable
type Variable struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Literal types pour les valeurs constantes
type BooleanLiteral struct {
	Type  string `json:"type"`
	Value bool   `json:"value"`
}

type StringLiteral struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type NumberLiteral struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type IntegerLiteral struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

// Constructeurs pour faciliter la création des structures

// NewProgram crée un nouveau programme
func NewProgram() *Program {
	return &Program{
		Types:       make([]TypeDefinition, 0),
		Expressions: make([]Expression, 0),
		Metadata: &Metadata{
			Version:   "1.0",
			CreatedAt: time.Now(),
		},
	}
}

// NewTypeDefinition crée une nouvelle définition de type
func NewTypeDefinition(name string) *TypeDefinition {
	return &TypeDefinition{
		Type:   "typeDefinition",
		Name:   name,
		Fields: make([]Field, 0),
	}
}

// AddField ajoute un champ à la définition de type
func (td *TypeDefinition) AddField(name, fieldType string) {
	td.Fields = append(td.Fields, Field{
		Name: name,
		Type: fieldType,
	})
}

// NewExpression crée une nouvelle expression/règle
func NewExpression() *Expression {
	return &Expression{
		Type: "expression",
		Set: Set{
			Type:      "set",
			Variables: make([]TypedVariable, 0),
		},
	}
}

// AddVariable ajoute une variable au set de l'expression
func (e *Expression) AddVariable(name, dataType string) {
	e.Set.Variables = append(e.Set.Variables, TypedVariable{
		Type:     "typedVariable",
		Name:     name,
		DataType: dataType,
	})
}

// NewConstraint crée une nouvelle contrainte binaire
func NewConstraint(left interface{}, operator string, right interface{}) *Constraint {
	return &Constraint{
		Type:     "constraint",
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// NewFieldAccess crée un nouvel accès de champ
func NewFieldAccess(object, field string) *FieldAccess {
	return &FieldAccess{
		Type:   "fieldAccess",
		Object: object,
		Field:  field,
	}
}

// NewAction crée une nouvelle action
func NewAction(jobName string, args ...interface{}) *Action {
	return &Action{
		Type: "action",
		Job: JobCall{
			Type: "jobCall",
			Name: jobName,
			Args: args,
		},
	}
}

// Méthodes utilitaires

// String retourne une représentation JSON formatée du programme
func (p *Program) String() string {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Sprintf("Program{Types: %d, Expressions: %d}", len(p.Types), len(p.Expressions))
	}
	return string(data)
}

// GetTypeByName trouve une définition de type par son nom
func (p *Program) GetTypeByName(name string) *TypeDefinition {
	for i := range p.Types {
		if p.Types[i].Name == name {
			return &p.Types[i]
		}
	}
	return nil
}

// GetFieldByName trouve un champ dans la définition de type par son nom
func (td *TypeDefinition) GetFieldByName(name string) *Field {
	for i := range td.Fields {
		if td.Fields[i].Name == name {
			return &td.Fields[i]
		}
	}
	return nil
}

// HasField vérifie si un type a un champ donné
func (td *TypeDefinition) HasField(name string) bool {
	return td.GetFieldByName(name) != nil
}

// Validation helpers

// IsValidOperator vérifie si un opérateur est valide
func IsValidOperator(op string) bool {
	validOps := map[string]bool{
		"==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
		"AND": true, "OR": true, "NOT": true,
		"+": true, "-": true, "*": true, "/": true, "%": true,
	}
	return validOps[op]
}

// IsValidType vérifie si un type est valide
func IsValidType(t string) bool {
	validTypes := map[string]bool{
		"string": true, "number": true, "bool": true, "integer": true,
	}
	return validTypes[t]
}
