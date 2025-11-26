// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
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
	Type string  `json:"type"`
	Job  JobCall `json:"job"`
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
