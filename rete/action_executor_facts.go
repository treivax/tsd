// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// generateFactID génère un ID unique pour un nouveau fait.
//
// Cette fonction est appelée par evaluateFactCreation pour créer des IDs uniques.
// L'implémentation actuelle utilise un compteur global thread-safe géré par le network.
//
// Paramètres :
//   - typeName : nom du type de fait
//   - network : réseau RETE contenant le générateur d'IDs
//
// Retourne :
//   - string : ID unique au format "typeName_N"
func (ae *ActionExecutor) generateFactID(typeName string) string {
	return ae.network.GenerateFactID(typeName)
}

// evaluateFactCreation évalue une création de fait.
//
// Validation :
//   - Vérifie que le contexte contient un network valide
//   - Vérifie que le type existe dans le network
//   - Valide tous les champs selon la définition de type
//
// Paramètres :
//   - argMap : map contenant "typeName" et "fields"
//   - ctx : contexte d'exécution
//
// Retourne :
//   - interface{} : nouveau fait créé
//   - error : erreur si validation échoue
func (ae *ActionExecutor) evaluateFactCreation(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("contexte d'exécution nil")
	}
	if ctx.network == nil {
		return nil, fmt.Errorf("network non disponible dans le contexte")
	}

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

	// Créer le nouveau fait avec un ID unique
	newFact := &Fact{
		ID:     ae.generateFactID(typeName),
		Type:   typeName,
		Fields: evaluatedFields,
	}

	return newFact, nil
}

// evaluateFactModification évalue une modification de fait.
//
// Validation :
//   - Vérifie que le contexte contient un network valide
//   - Vérifie que la variable existe dans le contexte
//   - Valide le fait modifié selon sa définition de type
//
// Paramètres :
//   - argMap : map contenant "variable", "field", "value"
//   - ctx : contexte d'exécution
//
// Retourne :
//   - interface{} : fait modifié (copie du fait original avec nouvelle valeur)
//   - error : erreur si validation échoue
func (ae *ActionExecutor) evaluateFactModification(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("contexte d'exécution nil")
	}

	// Récupérer le fait cible
	varName, ok := argMap["variable"].(string)
	if !ok {
		return nil, fmt.Errorf("variable manquante dans factModification")
	}

	fact := ctx.GetVariable(varName)
	if fact == nil {
		return nil, fmt.Errorf("variable '%s' non trouvée", varName)
	}

	// Récupérer le nom du champ à modifier
	fieldName, ok := argMap["field"].(string)
	if !ok {
		return nil, fmt.Errorf("field manquant dans factModification")
	}

	// Évaluer la nouvelle valeur
	newValue, ok := argMap["value"]
	if !ok {
		return nil, fmt.Errorf("value manquante dans factModification")
	}

	evaluated, err := ae.evaluateArgument(newValue, ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation nouvelle valeur: %w", err)
	}

	// Créer une copie modifiée du fait
	modifiedFact := &Fact{
		ID:     fact.ID,
		Type:   fact.Type,
		Fields: make(map[string]interface{}),
	}

	// Copier tous les champs existants
	for k, v := range fact.Fields {
		modifiedFact.Fields[k] = v
	}

	// Appliquer la modification
	modifiedFact.Fields[fieldName] = evaluated

	// Valider le fait modifié si le network est disponible
	if ctx.network != nil {
		typeDef := ctx.network.GetTypeDefinition(fact.Type)
		if typeDef != nil {
			if err := ae.validateFactFields(typeDef, modifiedFact.Fields); err != nil {
				return nil, fmt.Errorf("validation fact modification: %w", err)
			}
		}
	}

	return modifiedFact, nil
}
