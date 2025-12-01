// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"log"
	"strings"
)

// ActionExecutor gère l'exécution des actions déclenchées par les règles
type ActionExecutor struct {
	network       *ReteNetwork
	logger        *log.Logger
	enableLogging bool
}

// NewActionExecutor crée un nouveau exécuteur d'actions
func NewActionExecutor(network *ReteNetwork, logger *log.Logger) *ActionExecutor {
	if logger == nil {
		logger = log.Default()
	}
	return &ActionExecutor{
		network:       network,
		logger:        logger,
		enableLogging: true,
	}
}

// SetLogging active ou désactive le logging des actions
func (ae *ActionExecutor) SetLogging(enabled bool) {
	ae.enableLogging = enabled
}

// ExecuteAction exécute une action avec les faits fournis par le token
func (ae *ActionExecutor) ExecuteAction(action *Action, token *Token) error {
	if action == nil {
		return fmt.Errorf("action is nil")
	}

	// Obtenir tous les jobs à exécuter
	jobs := action.GetJobs()

	// Créer un contexte d'exécution avec les faits disponibles
	ctx := NewExecutionContext(token, ae.network)

	// Exécuter chaque job en séquence
	for i, job := range jobs {
		if err := ae.executeJob(job, ctx, i); err != nil {
			return fmt.Errorf("erreur exécution job %s: %w", job.Name, err)
		}
	}

	return nil
}

// executeJob exécute un job individuel
func (ae *ActionExecutor) executeJob(job JobCall, ctx *ExecutionContext, jobIndex int) error {
	// Logger l'action
	if ae.enableLogging {
		ae.logAction(job, ctx)
	}

	// Évaluer les arguments
	evaluatedArgs := make([]interface{}, 0, len(job.Args))
	for i, arg := range job.Args {
		evaluated, err := ae.evaluateArgument(arg, ctx)
		if err != nil {
			return fmt.Errorf("erreur évaluation argument %d: %w", i, err)
		}
		evaluatedArgs = append(evaluatedArgs, evaluated)
	}

	// Exécuter l'action (actuellement, on se contente de logger)
	// Dans une implémentation complète, on pourrait dispatcher vers différents handlers
	ae.logger.Printf("🎯 ACTION EXÉCUTÉE: %s(%v)", job.Name, formatArgs(evaluatedArgs))

	return nil
}

// evaluateArgument évalue un argument selon son type
func (ae *ActionExecutor) evaluateArgument(arg interface{}, ctx *ExecutionContext) (interface{}, error) {
	// Cas 1: Valeur littérale simple (string, number, bool)
	switch v := arg.(type) {
	case string, float64, bool, int, int64:
		return v, nil
	}

	// Cas 2: Map (objet structuré du parser)
	argMap, ok := arg.(map[string]interface{})
	if !ok {
		return arg, nil // Retourner tel quel si on ne peut pas le comprendre
	}

	argType, hasType := argMap["type"].(string)
	if !hasType {
		return arg, nil
	}

	switch argType {
	case "string", "number", "bool":
		// Valeur littérale typée
		if value, ok := argMap["value"]; ok {
			return value, nil
		}
		return arg, nil

	case "variable":
		// Cas 2: Fait complet référencé par variable
		varName, ok := argMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de variable invalide")
		}
		fact := ctx.GetVariable(varName)
		if fact == nil {
			return nil, fmt.Errorf("variable '%s' non trouvée", varName)
		}
		return fact, nil

	case "fieldAccess":
		// Cas 3: Attribut de fait (variable.attribut)
		objectName, ok := argMap["object"].(string)
		if !ok {
			return nil, fmt.Errorf("nom d'objet invalide dans fieldAccess")
		}
		fieldName, ok := argMap["field"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de champ invalide dans fieldAccess")
		}

		fact := ctx.GetVariable(objectName)
		if fact == nil {
			return nil, fmt.Errorf("variable '%s' non trouvée", objectName)
		}

		value, exists := fact.Fields[fieldName]
		if !exists {
			return nil, fmt.Errorf("champ '%s' non trouvé dans le fait %s", fieldName, objectName)
		}
		return value, nil

	case "factCreation":
		// Cas 4: Création de nouveau fait
		return ae.evaluateFactCreation(argMap, ctx)

	case "factModification":
		// Cas 5: Modification de fait
		return ae.evaluateFactModification(argMap, ctx)

	case "arithmetic":
		// Expression arithmétique
		return ae.evaluateArithmetic(argMap, ctx)

	default:
		return arg, nil
	}
}

// evaluateFactCreation évalue une création de fait
func (ae *ActionExecutor) evaluateFactCreation(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	typeName, ok := argMap["typeName"].(string)
	if !ok {
		return nil, fmt.Errorf("typeName manquant dans factCreation")
	}

	// Vérifier que le type existe
	typeDef := ctx.network.GetTypeDefinition(typeName)
	if typeDef == nil {
		return nil, fmt.Errorf("type '%s' non défini", typeName)
	}

	// Extraire les champs
	fieldsData, ok := argMap["fields"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("champs manquants dans factCreation")
	}

	// Évaluer chaque valeur de champ
	evaluatedFields := make(map[string]interface{})
	for fieldName, fieldValue := range fieldsData {
		evaluated, err := ae.evaluateArgument(fieldValue, ctx)
		if err != nil {
			return nil, fmt.Errorf("erreur évaluation champ '%s': %w", fieldName, err)
		}
		evaluatedFields[fieldName] = evaluated
	}

	// Valider la cohérence avec la définition de type
	if err := ae.validateFactFields(typeDef, evaluatedFields); err != nil {
		return nil, fmt.Errorf("validation fact creation: %w", err)
	}

	// Créer le nouveau fait
	newFact := &Fact{
		ID:     generateFactID(typeName),
		Type:   typeName,
		Fields: evaluatedFields,
	}

	return newFact, nil
}

// evaluateFactModification évalue une modification de fait
func (ae *ActionExecutor) evaluateFactModification(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	varName, ok := argMap["variable"].(string)
	if !ok {
		return nil, fmt.Errorf("variable manquante dans factModification")
	}

	fieldName, ok := argMap["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field manquant dans factModification")
	}

	newValue := argMap["value"]
	if newValue == nil {
		return nil, fmt.Errorf("value manquante dans factModification")
	}

	// Récupérer le fait
	fact := ctx.GetVariable(varName)
	if fact == nil {
		return nil, fmt.Errorf("variable '%s' non trouvée", varName)
	}

	// Vérifier que le type existe
	typeDef := ctx.network.GetTypeDefinition(fact.Type)
	if typeDef == nil {
		return nil, fmt.Errorf("type '%s' non défini", fact.Type)
	}

	// Évaluer la nouvelle valeur
	evaluatedValue, err := ae.evaluateArgument(newValue, ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation nouvelle valeur: %w", err)
	}

	// Valider la cohérence avec la définition du champ
	if err := ae.validateFieldValue(typeDef, fieldName, evaluatedValue); err != nil {
		return nil, fmt.Errorf("validation field modification: %w", err)
	}

	// Créer une copie modifiée du fait
	modifiedFact := &Fact{
		ID:     fact.ID,
		Type:   fact.Type,
		Fields: make(map[string]interface{}),
	}
	for k, v := range fact.Fields {
		modifiedFact.Fields[k] = v
	}
	modifiedFact.Fields[fieldName] = evaluatedValue

	return modifiedFact, nil
}

// evaluateArithmetic évalue une expression arithmétique
func (ae *ActionExecutor) evaluateArithmetic(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	operator, ok := argMap["operator"].(string)
	if !ok {
		return nil, fmt.Errorf("opérateur manquant dans arithmetic")
	}

	left, err := ae.evaluateArgument(argMap["left"], ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation left: %w", err)
	}

	right, err := ae.evaluateArgument(argMap["right"], ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation right: %w", err)
	}

	// Convertir en nombres
	leftNum, okL := toNumber(left)
	rightNum, okR := toNumber(right)
	if !okL || !okR {
		return nil, fmt.Errorf("opération arithmétique nécessite des nombres")
	}

	switch operator {
	case "+":
		return leftNum + rightNum, nil
	case "-":
		return leftNum - rightNum, nil
	case "*":
		return leftNum * rightNum, nil
	case "/":
		if rightNum == 0 {
			return nil, fmt.Errorf("division par zéro")
		}
		return leftNum / rightNum, nil
	default:
		return nil, fmt.Errorf("opérateur arithmétique inconnu: %s", operator)
	}
}

// validateFactFields valide que tous les champs requis sont présents et corrects
func (ae *ActionExecutor) validateFactFields(typeDef *TypeDefinition, fields map[string]interface{}) error {
	// Vérifier que tous les champs définis sont présents
	for _, fieldDef := range typeDef.Fields {
		value, exists := fields[fieldDef.Name]
		if !exists {
			return fmt.Errorf("champ requis '%s' manquant", fieldDef.Name)
		}

		if err := ae.validateFieldType(fieldDef.Type, value); err != nil {
			return fmt.Errorf("champ '%s': %w", fieldDef.Name, err)
		}
	}

	// Vérifier qu'il n'y a pas de champs non définis
	for fieldName := range fields {
		found := false
		for _, fieldDef := range typeDef.Fields {
			if fieldDef.Name == fieldName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("champ '%s' non défini dans le type", fieldName)
		}
	}

	return nil
}

// validateFieldValue valide qu'une valeur de champ est cohérente avec sa définition
func (ae *ActionExecutor) validateFieldValue(typeDef *TypeDefinition, fieldName string, value interface{}) error {
	// Trouver la définition du champ
	var fieldDef *Field
	for i := range typeDef.Fields {
		if typeDef.Fields[i].Name == fieldName {
			fieldDef = &typeDef.Fields[i]
			break
		}
	}

	if fieldDef == nil {
		return fmt.Errorf("champ '%s' non défini dans le type '%s'", fieldName, typeDef.Name)
	}

	return ae.validateFieldType(fieldDef.Type, value)
}

// validateFieldType valide qu'une valeur correspond au type attendu
func (ae *ActionExecutor) validateFieldType(expectedType string, value interface{}) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("type attendu: string, reçu: %T", value)
		}
	case "number":
		if _, ok := toNumber(value); !ok {
			return fmt.Errorf("type attendu: number, reçu: %T", value)
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("type attendu: bool, reçu: %T", value)
		}
	default:
		// Type personnalisé ou non reconnu
		return nil
	}
	return nil
}

// logAction log une action avec ses arguments
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

// toNumber convertit une valeur en float64
func toNumber(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	default:
		return 0, false
	}
}

// generateFactID génère un ID unique pour un nouveau fait
var factCounter = 0

func generateFactID(typeName string) string {
	factCounter++
	return fmt.Sprintf("%s_%d", typeName, factCounter)
}

// ExecutionContext contient le contexte d'exécution d'une action
type ExecutionContext struct {
	token    *Token
	network  *ReteNetwork
	varCache map[string]*Fact
}

// NewExecutionContext crée un nouveau contexte d'exécution
func NewExecutionContext(token *Token, network *ReteNetwork) *ExecutionContext {
	ctx := &ExecutionContext{
		token:    token,
		network:  network,
		varCache: make(map[string]*Fact),
	}

	// Construire le cache des variables depuis le token
	if token != nil && token.Bindings != nil {
		for varName, fact := range token.Bindings {
			ctx.varCache[varName] = fact
		}
	}

	return ctx
}

// GetVariable récupère un fait par nom de variable
func (ctx *ExecutionContext) GetVariable(name string) *Fact {
	return ctx.varCache[name]
}
