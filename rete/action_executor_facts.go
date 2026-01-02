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

// evaluateInlineFact évalue un fait inline créé dans une action.
//
// Les faits inline utilisent le format du parser PEG avec une liste de champs
// contenant des maps avec "name" et "value".
//
// Validation :
//   - Vérifie que le contexte contient un network valide
//   - Vérifie que le type existe dans le network
//   - Valide tous les champs selon la définition de type
//   - Évalue récursivement toutes les expressions dans les valeurs de champs
//
// Paramètres :
//   - argMap : map contenant "typeName" et "fields" (liste de {name, value})
//   - ctx : contexte d'exécution avec bindings des variables
//
// Retourne :
//   - interface{} : nouveau fait créé avec valeurs évaluées
//   - error : erreur si validation ou évaluation échoue
//
// Exemple d'utilisation dans une règle TSD :
//
//	rule alert: {s: Sensor} / s.temp > 40.0 ==>
//	    Xuple("alerts", Alert(
//	        level: "HIGH",
//	        sensorId: s.id,
//	        temperature: s.temp
//	    ))
func (ae *ActionExecutor) evaluateInlineFact(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("contexte d'exécution nil")
	}
	if ctx.network == nil {
		return nil, fmt.Errorf("network non disponible dans le contexte")
	}

	typeName, ok := argMap["typeName"].(string)
	if !ok {
		return nil, fmt.Errorf("typeName manquant dans inlineFact")
	}

	// Vérifier que le type existe
	typeDef := ctx.network.GetTypeDefinition(typeName)
	if typeDef == nil {
		return nil, fmt.Errorf("type '%s' non défini", typeName)
	}

	// Extraire la liste de champs (format du parser PEG)
	fieldsArray, ok := argMap["fields"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("champs manquants ou format invalide dans inlineFact")
	}

	// Convertir le format liste en map et évaluer les valeurs
	evaluatedFields := make(map[string]interface{})
	for _, fieldItem := range fieldsArray {
		fieldMap, ok := fieldItem.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("format de champ invalide dans inlineFact")
		}

		fieldName, ok := fieldMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("nom de champ manquant dans inlineFact")
		}

		fieldValue, ok := fieldMap["value"]
		if !ok {
			return nil, fmt.Errorf("valeur de champ manquante pour '%s' dans inlineFact", fieldName)
		}

		// Évaluer la valeur du champ (peut être une expression, référence, etc.)
		evaluated, err := ae.evaluateArgument(fieldValue, ctx)
		if err != nil {
			return nil, fmt.Errorf("erreur évaluation champ '%s': %w", fieldName, err)
		}
		evaluatedFields[fieldName] = evaluated
	}

	// Valider la cohérence avec la définition de type
	if err := ae.validateFactFields(typeDef, evaluatedFields); err != nil {
		return nil, fmt.Errorf("validation inline fact: %w", err)
	}

	// Générer l'ID interne basé sur les clés primaires si elles existent,
	// sinon utiliser un générateur d'ID unique
	internalID, err := ae.generateInternalIDFromFields(typeDef, evaluatedFields)
	if err != nil {
		return nil, fmt.Errorf("erreur génération ID: %w", err)
	}

	// Créer le nouveau fait avec l'ID basé sur les clés primaires
	newFact := &Fact{
		ID:     internalID,
		Type:   typeName,
		Fields: evaluatedFields,
	}

	return newFact, nil
}

// generateInternalIDFromFields génère un ID interne à partir des clés primaires d'un type.
//
// Si le type n'a pas de clés primaires définies (par exemple pour les inline facts
// créés dynamiquement dans les actions comme Xuple), un ID unique est généré
// automatiquement enParamètres :
//   - typeDef : définition du type contenant les champs et clés primaires
//   - fields : map des champs évalués du fait
//
// Retourne :
//   - string : ID interne au format "Type~valeur" ou "Type~val1_val2_..." pour clés composées
//   - error : si les clés primaires sont manquantes ou invalides
func (ae *ActionExecutor) generateInternalIDFromFields(typeDef *TypeDefinition, fields map[string]interface{}) (string, error) {
	// Récupérer les clés primaires
	primaryKeys := []string{}
	for _, field := range typeDef.Fields {
		if field.IsPrimaryKey {
			primaryKeys = append(primaryKeys, field.Name)
		}
	}

	// Si le type n'a pas de clés primaires (par exemple pour les inline facts
	// créés dynamiquement comme Alert, Command dans les xuples), générer un ID unique
	if len(primaryKeys) == 0 {
		if ae.network != nil {
			return ae.network.GenerateFactID(typeDef.Name), nil
		}
		// Fallback si network n'est pas disponible (ne devrait pas arriver)
		return fmt.Sprintf("%s_%d", typeDef.Name, len(fields)), nil
	}

	// Extraire les valeurs des clés primaires
	keyValues := []string{}
	for _, pkName := range primaryKeys {
		value, exists := fields[pkName]
		if !exists {
			return "", fmt.Errorf("clé primaire '%s' manquante dans les champs", pkName)
		}

		// Convertir la valeur en string
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case float64:
			strValue = fmt.Sprintf("%v", v)
		case int, int64:
			strValue = fmt.Sprintf("%d", v)
		case bool:
			strValue = fmt.Sprintf("%v", v)
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		keyValues = append(keyValues, strValue)
	}

	// Construire l'ID interne : Type~valeur ou Type~val1_val2_...
	if len(keyValues) == 1 {
		return fmt.Sprintf("%s~%s", typeDef.Name, keyValues[0]), nil
	}
	return fmt.Sprintf("%s~%s", typeDef.Name, joinKeys(keyValues)), nil
}

// joinKeys joint les valeurs de clés avec underscores, en gérant l'encodage URL si nécessaire.
func joinKeys(keys []string) string {
	// Pour l'instant, simple join avec underscore
	// TODO: gérer l'encodage URL si les valeurs contiennent des caractères spéciaux
	result := ""
	for i, key := range keys {
		if i > 0 {
			result += "_"
		}
		result += key
	}
	return result
}

// evaluateUpdateWithModifications évalue une mise à jour de fait avec modifications de champs.
//
// Cette fonction implémente la nouvelle sémantique d'Update qui préserve l'ID du fait original.
// Syntaxe: Update(variable, field: value, field: value, ...)
//
// Validation :
//   - Vérifie que le contexte contient un network valide
//   - Vérifie que la variable existe dans le contexte et référence un fait
//   - Valide chaque modification de champ
//   - Valide le fait modifié selon sa définition de type
//
// Comportement :
//   - Récupère le fait original via la variable
//   - Préserve l'ID et le type du fait original
//   - Copie tous les champs existants
//   - Applique les modifications spécifiées
//   - Retourne un nouveau fait avec l'ID préservé
//
// Paramètres :
//   - argMap : map contenant "variable" (string) et "modifications" (map[string]interface{})
//   - ctx : contexte d'exécution avec bindings des variables
//
// Retourne :
//   - interface{} : fait modifié avec ID préservé
//   - error : erreur si validation échoue
//
// Exemple d'utilisation dans une règle TSD :
//
//	rule mettre_en_couple : {p: Personne, r: Relation} /
//	    p.nom == r.personne1 AND r.lien == "mariage" ==>
//	    Update(p, statut: "en couple")
func (ae *ActionExecutor) evaluateUpdateWithModifications(argMap map[string]interface{}, ctx *ExecutionContext) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("contexte d'exécution nil")
	}
	if ctx.network == nil {
		return nil, fmt.Errorf("network non disponible dans le contexte")
	}

	// Récupérer et évaluer la variable (peut être map[name:p type:variable])
	variableArg, ok := argMap["variable"]
	if !ok {
		return nil, fmt.Errorf("variable manquante dans updateWithModifications")
	}

	// Évaluer l'argument variable pour obtenir le fait
	originalFact, err := ae.evaluateArgument(variableArg, ctx)
	if err != nil {
		return nil, fmt.Errorf("erreur évaluation variable: %w", err)
	}

	// Vérifier que c'est bien un fait
	fact, ok := originalFact.(*Fact)
	if !ok {
		return nil, fmt.Errorf("la variable ne référence pas un fait valide (type: %T)", originalFact)
	}

	// Récupérer les modifications
	modifications, ok := argMap["modifications"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("modifications manquantes ou invalides dans updateWithModifications")
	}

	// Créer une copie du fait avec l'ID préservé
	modifiedFact := &Fact{
		ID:     fact.ID, // PRÉSERVATION DE L'ID
		Type:   fact.Type,
		Fields: make(map[string]interface{}),
	}

	// Copier tous les champs existants
	for fieldName, fieldValue := range fact.Fields {
		modifiedFact.Fields[fieldName] = fieldValue
	}

	// Appliquer les modifications et détecter si quelque chose a réellement changé
	hasChanges := false
	for fieldName, newValue := range modifications {
		// Évaluer la nouvelle valeur (peut être une expression)
		evaluated, err := ae.evaluateArgument(newValue, ctx)
		if err != nil {
			return nil, fmt.Errorf("erreur évaluation champ '%s': %w", fieldName, err)
		}

		// Vérifier si la valeur a réellement changé
		oldValue, exists := modifiedFact.Fields[fieldName]
		if !exists || !areValuesEqualForUpdate(oldValue, evaluated) {
			hasChanges = true
		}

		modifiedFact.Fields[fieldName] = evaluated
	}

	// Si aucune valeur n'a changé, retourner le fait original pour éviter une boucle infinie
	if !hasChanges {
		return fact, nil
	}

	// Valider le fait modifié selon sa définition de type
	// Exclure _id_ qui est un champ système
	typeDef := ctx.network.GetTypeDefinition(fact.Type)
	if typeDef != nil {
		fieldsToValidate := make(map[string]interface{})
		for k, v := range modifiedFact.Fields {
			if k != "_id_" {
				fieldsToValidate[k] = v
			}
		}
		if err := ae.validateFactFields(typeDef, fieldsToValidate); err != nil {
			return nil, fmt.Errorf("validation update with modifications: %w", err)
		}
	}

	return modifiedFact, nil
}

// areValuesEqualForUpdate compare deux valeurs pour détecter si elles sont égales
// Gère les conversions de types numériques (int, int64, float64)
func areValuesEqualForUpdate(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Comparaison directe pour les types identiques
	if a == b {
		return true
	}

	// Comparaison spéciale pour les nombres (int, int64, float64)
	aNum, aIsNum := toFloat64ForUpdate(a)
	bNum, bIsNum := toFloat64ForUpdate(b)
	if aIsNum && bIsNum {
		return aNum == bNum
	}

	return false
}

// toFloat64ForUpdate convertit un nombre en float64 si possible
func toFloat64ForUpdate(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}
