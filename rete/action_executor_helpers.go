// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// logAction logs an action before execution
func (ae *ActionExecutor) logAction(job JobCall, ctx *ExecutionContext) {
	var argsStr []string
	for _, arg := range job.Args {
		argsStr = append(argsStr, formatArgument(arg, ctx))
	}

	ae.logger.Printf("📋 ACTION: %s(%s)", job.Name, strings.Join(argsStr, ", "))
}

// formatArgument formate un argument pour le logging
func formatArgument(arg interface{}, ctx *ExecutionContext) string {
	switch v := arg.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case float64, int, int64:
		return fmt.Sprintf("%v", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case map[string]interface{}:
		argType, _ := v["type"].(string)
		switch argType {
		case "string", "number", "bool":
			if value, ok := v["value"]; ok {
				return formatArgument(value, ctx)
			}
		case "variable":
			if name, ok := v["name"].(string); ok {
				return name
			}
		case "fieldAccess":
			if obj, ok := v["object"].(string); ok {
				if field, ok := v["field"].(string); ok {
					return fmt.Sprintf("%s.%s", obj, field)
				}
			}
		case "factCreation":
			if typeName, ok := v["typeName"].(string); ok {
				return fmt.Sprintf("new %s(...)", typeName)
			}
		case "factModification":
			if varName, ok := v["variable"].(string); ok {
				if field, ok := v["field"].(string); ok {
					return fmt.Sprintf("%s[%s]=...", varName, field)
				}
			}
		}
	}
	return fmt.Sprintf("%v", arg)
}

// formatArgs formate une liste d'arguments évalués
func formatArgs(args []interface{}) string {
	var parts []string
	for _, arg := range args {
		switch v := arg.(type) {
		case *Fact:
			parts = append(parts, fmt.Sprintf("%s{%s}", v.Type, v.ID))
		case string:
			parts = append(parts, fmt.Sprintf("%q", v))
		default:
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(parts, ", ")
}
