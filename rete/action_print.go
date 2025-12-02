// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"io"
	"os"
)

// Nom de l'action print
const ActionNamePrint = "print"

// PrintAction implémente l'action "print" qui affiche une chaîne de caractères.
type PrintAction struct {
	output io.Writer
}

// NewPrintAction crée une nouvelle action print.
// Si output est nil, utilise os.Stdout par défaut.
func NewPrintAction(output io.Writer) *PrintAction {
	if output == nil {
		output = os.Stdout
	}
	return &PrintAction{
		output: output,
	}
}

// Execute exécute l'action print en affichant la chaîne de caractères fournie.
func (pa *PrintAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	if len(args) == 0 {
		return fmt.Errorf("print action requires at least one argument")
	}

	// Récupérer le premier argument (la chaîne à afficher)
	arg := args[0]

	// Convertir l'argument en chaîne de caractères
	str, err := pa.convertToString(arg)
	if err != nil {
		return fmt.Errorf("failed to convert argument to string: %w", err)
	}

	// Afficher la chaîne
	_, err = fmt.Fprintln(pa.output, str)
	if err != nil {
		return fmt.Errorf("failed to print: %w", err)
	}

	return nil
}

// GetName retourne le nom de l'action.
func (pa *PrintAction) GetName() string {
	return ActionNamePrint
}

// Validate valide que les arguments sont corrects pour l'action print.
func (pa *PrintAction) Validate(args []interface{}) error {
	if len(args) == 0 {
		return fmt.Errorf("print action requires at least one argument")
	}

	// Vérifier que le premier argument peut être converti en chaîne
	arg := args[0]
	_, err := pa.convertToString(arg)
	if err != nil {
		return fmt.Errorf("argument cannot be converted to string: %w", err)
	}

	return nil
}

// convertToString convertit un argument en chaîne de caractères.
func (pa *PrintAction) convertToString(arg interface{}) (string, error) {
	if arg == nil {
		return "", fmt.Errorf("argument is nil")
	}

	switch v := arg.(type) {
	case string:
		return v, nil
	case float64:
		return fmt.Sprintf("%g", v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case *Fact:
		// Vérifier si le pointeur Fact est nil
		if v == nil {
			return "", fmt.Errorf("fact is nil")
		}
		// Pour un fait, afficher sa représentation JSON-like
		return pa.factToString(v), nil
	case map[string]interface{}:
		// Pour un objet structuré, afficher sa représentation
		return fmt.Sprintf("%v", v), nil
	default:
		// Pour tout autre type, utiliser fmt.Sprintf
		return fmt.Sprintf("%v", v), nil
	}
}

// factToString convertit un fait en chaîne de caractères lisible.
func (pa *PrintAction) factToString(fact *Fact) string {
	if fact == nil {
		return "null"
	}

	result := fmt.Sprintf("%s{id: %s", fact.Type, fact.ID)
	for key, value := range fact.Fields {
		switch v := value.(type) {
		case string:
			result += fmt.Sprintf(", %s: %q", key, v)
		default:
			result += fmt.Sprintf(", %s: %v", key, v)
		}
	}
	result += "}"
	return result
}

// SetOutput change le writer de sortie pour l'action print.
func (pa *PrintAction) SetOutput(output io.Writer) {
	if output != nil {
		pa.output = output
	}
}
