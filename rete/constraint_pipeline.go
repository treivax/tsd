package rete

import (
	"fmt"
	"strings"

	"github.com/treivax/tsd/constraint"
)

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
}

// ConstraintPipeline implémente le pipeline complet :
// fichier .constraint → parseur PEG → conversion AST → réseau RETE
type ConstraintPipeline struct{}

// NewConstraintPipeline crée une nouvelle instance du pipeline
func NewConstraintPipeline() *ConstraintPipeline {
	return &ConstraintPipeline{}
}

// BuildNetworkFromConstraintFile construit un réseau RETE complet à partir d'un fichier .constraint
// Cette fonction implémente le pipeline unique utilisé par TOUS les tests
func (cp *ConstraintPipeline) BuildNetworkFromConstraintFile(constraintFile string, storage Storage) (*ReteNetwork, error) {
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichier: %s\n", constraintFile)

	// ÉTAPE 1: Parsing avec le vrai parseur PEG
	parsedAST, err := constraint.ParseConstraintFile(constraintFile)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur parsing fichier %s: %w", constraintFile, err)
	}
	fmt.Printf("✅ Parsing réussi\n")

	// Valider que c'est un map[string]interface{}
	resultMap, ok := parsedAST.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ Format AST non reconnu: %T", parsedAST)
	}

	// ÉTAPE 2: Extraction et validation des composants
	types, expressions, err := cp.extractComponents(resultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur extraction composants: %w", err)
	}
	fmt.Printf("✅ Trouvé %d types et %d expressions\n", len(types), len(expressions))

	// ÉTAPE 3: Construction du réseau RETE
	network, err := cp.buildNetwork(storage, types, expressions)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur construction réseau: %w", err)
	}
	fmt.Printf("✅ Réseau construit avec %d nœuds terminaux\n", len(network.TerminalNodes))

	// ÉTAPE 4: Validation finale
	err = cp.validateNetwork(network)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur validation réseau: %w", err)
	}
	fmt.Printf("✅ Validation réussie\n")

	fmt.Printf("🎯 PIPELINE TERMINÉ AVEC SUCCÈS\n")
	fmt.Printf("========================================\n\n")

	return network, nil
}

// BuildNetworkFromMultipleFiles construit un réseau RETE en parsant plusieurs fichiers de manière itérative
// Cette fonction permet de parser des types, règles et faits répartis dans différents fichiers
func (cp *ConstraintPipeline) BuildNetworkFromMultipleFiles(filenames []string, storage Storage) (*ReteNetwork, error) {
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichiers: %v\n", filenames)

	// Créer un parser itératif
	parser := constraint.NewIterativeParser()

	// Parser tous les fichiers de manière itérative
	for i, filename := range filenames {
		fmt.Printf("  📄 Parsing fichier %d/%d: %s\n", i+1, len(filenames), filename)
		err := parser.ParseFile(filename)
		if err != nil {
			return nil, fmt.Errorf("❌ Erreur parsing fichier %s: %w", filename, err)
		}
	}
	fmt.Printf("✅ Parsing itératif réussi\n")

	// Obtenir le programme combiné
	program := parser.GetProgram()

	// Statistiques

	// Convertir au format RETE
	reteProgram := constraint.ConvertToReteProgram(program)
	resultMap, ok := reteProgram.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ Format programme RETE invalide: %T", reteProgram)
	}

	// Extraire les composants
	types, expressions, err := cp.extractComponents(resultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur extraction composants: %w", err)
	}
	fmt.Printf("✅ Trouvé %d types et %d expressions\n", len(types), len(expressions))

	// Construction du réseau RETE
	network, err := cp.buildNetwork(storage, types, expressions)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur construction réseau: %w", err)
	}
	fmt.Printf("✅ Réseau construit avec %d nœuds terminaux\n", len(network.TerminalNodes))

	// Injection des faits dans le réseau
	if len(program.Facts) > 0 {
		factsForRete := constraint.ConvertFactsToReteFormat(*program)

		err := network.SubmitFactsFromGrammar(factsForRete)
		if err != nil {
			fmt.Printf("❌ Erreur injection faits: %v\n", err)
		} else {
			fmt.Printf("✅ Injection terminée: %d faits injectés\n", len(factsForRete))
		}
	}

	fmt.Printf("🎯 PIPELINE MULTIFILES TERMINÉ\n")
	fmt.Printf("========================================\n\n")

	return network, nil
}

// BuildNetworkFromIterativeParser construit un réseau RETE à partir d'un parser itératif existant
// Cette méthode est utile quand le parsing a déjà été fait et qu'on veut juste construire le réseau
func (cp *ConstraintPipeline) BuildNetworkFromIterativeParser(parser *constraint.IterativeParser, storage Storage) (*ReteNetwork, error) {
	fmt.Printf("========================================\n")

	// Obtenir le programme combiné
	program := parser.GetProgram()

	// Statistiques

	// Convertir au format RETE
	reteProgram := constraint.ConvertToReteProgram(program)
	resultMap, ok := reteProgram.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ Format programme RETE invalide: %T", reteProgram)
	}

	// Extraire les composants
	types, expressions, err := cp.extractComponents(resultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur extraction composants: %w", err)
	}
	fmt.Printf("✅ Trouvé %d types et %d expressions\n", len(types), len(expressions))

	// Construction du réseau RETE
	network, err := cp.buildNetwork(storage, types, expressions)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur construction réseau: %w", err)
	}
	fmt.Printf("✅ Réseau construit avec %d nœuds terminaux\n", len(network.TerminalNodes))

	// Injection des faits dans le réseau
	if len(program.Facts) > 0 {
		factsForRete := constraint.ConvertFactsToReteFormat(*program)

		err := network.SubmitFactsFromGrammar(factsForRete)
		if err != nil {
			fmt.Printf("❌ Erreur injection faits: %v\n", err)
		} else {
			fmt.Printf("✅ Injection terminée: %d faits injectés\n", len(factsForRete))
		}
	}

	fmt.Printf("🎯 PIPELINE DEPUIS PARSER TERMINÉ\n")
	fmt.Printf("========================================\n\n")

	return network, nil
}

// extractComponents extrait les types et expressions du map parsé
func (cp *ConstraintPipeline) extractComponents(resultMap map[string]interface{}) ([]interface{}, []interface{}, error) {
	// Extraire les types
	typesData, hasTypes := resultMap["types"]
	if !hasTypes {
		return nil, nil, fmt.Errorf("aucun type trouvé dans le fichier")
	}
	types, ok := typesData.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("format types invalide: %T", typesData)
	}

	// Extraire les expressions
	exprsData, hasExprs := resultMap["expressions"]
	if !hasExprs {
		return nil, nil, fmt.Errorf("aucune expression trouvée dans le fichier")
	}
	expressions, ok := exprsData.([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("format expressions invalide: %T", exprsData)
	}

	return types, expressions, nil
}

// buildNetwork construit le réseau RETE à partir des composants extraits
func (cp *ConstraintPipeline) buildNetwork(storage Storage, types []interface{}, expressions []interface{}) (*ReteNetwork, error) {
	network := NewReteNetwork(storage)

	// Créer les types de données
	err := cp.createTypeNodes(network, types, storage)
	if err != nil {
		return nil, fmt.Errorf("erreur création types: %w", err)
	}

	// Créer les nœuds pour les règles
	err = cp.createRuleNodes(network, expressions, storage)
	if err != nil {
		return nil, fmt.Errorf("erreur création règles: %w", err)
	}

	return network, nil
}

// createTypeNodes crée les nœuds de type à partir des définitions parsées
func (cp *ConstraintPipeline) createTypeNodes(network *ReteNetwork, types []interface{}, storage Storage) error {
	for i, typeData := range types {
		typeMap, ok := typeData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format type %d invalide: %T", i, typeData)
		}

		// Extraire le nom du type
		nameData, hasName := typeMap["name"]
		if !hasName {
			return fmt.Errorf("nom manquant pour type %d", i)
		}
		typeName, ok := nameData.(string)
		if !ok {
			return fmt.Errorf("nom type %d invalide: %T", i, nameData)
		}

		// Créer une définition de type RETE
		typeDef := cp.createTypeDefinition(typeName, typeMap)

		// Créer et ajouter le nœud de type
		typeNode := NewTypeNode(typeName, typeDef, storage)
		network.TypeNodes[typeName] = typeNode
		network.RootNode.AddChild(typeNode)

		fmt.Printf("   ✓ TypeNode créé: %s\n", typeName)
	}

	return nil
}

// createTypeDefinition crée une définition de type RETE à partir d'un map parsé
func (cp *ConstraintPipeline) createTypeDefinition(typeName string, typeMap map[string]interface{}) TypeDefinition {
	typeDef := TypeDefinition{
		Name:   typeName,
		Fields: []Field{},
	}

	// Extraire les champs si disponibles
	if fieldsData, hasFields := typeMap["fields"]; hasFields {
		if fieldsList, ok := fieldsData.([]interface{}); ok {
			for _, fieldData := range fieldsList {
				if fieldMap, ok := fieldData.(map[string]interface{}); ok {
					field := Field{
						Name: getStringField(fieldMap, "name", ""),
						Type: getStringField(fieldMap, "type", "string"),
					}
					typeDef.Fields = append(typeDef.Fields, field)
				}
			}
		}
	}

	// Si pas de champs définis, créer des champs par défaut selon le nom du type
	if len(typeDef.Fields) == 0 {
		switch typeName {
		case "Utilisateur":
			typeDef.Fields = []Field{
				{Name: "id", Type: "string"},
				{Name: "nom", Type: "string"},
				{Name: "prenom", Type: "string"},
				{Name: "age", Type: "number"},
			}
		case "Adresse":
			typeDef.Fields = []Field{
				{Name: "utilisateur_id", Type: "string"},
				{Name: "rue", Type: "string"},
				{Name: "ville", Type: "string"},
			}
		default:
			// Type générique
			typeDef.Fields = []Field{
				{Name: "id", Type: "string"},
				{Name: "value", Type: "string"},
			}
		}
	}

	return typeDef
}

// createRuleNodes crée les nœuds de règles à partir des expressions parsées
func (cp *ConstraintPipeline) createRuleNodes(network *ReteNetwork, expressions []interface{}, storage Storage) error {
	for i, exprData := range expressions {
		ruleID := fmt.Sprintf("rule_%d", i)

		exprMap, ok := exprData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("format expression %d invalide: %T", i, exprData)
		}

		err := cp.createSingleRule(network, ruleID, exprMap, storage)
		if err != nil {
			return fmt.Errorf("erreur création règle %s: %w", ruleID, err)
		}

		fmt.Printf("   ✓ Règle créée: %s\n", ruleID)
	}

	return nil
}

// createSingleRule crée une règle unique (Alpha + Terminal avec support des contraintes NOT)
func (cp *ConstraintPipeline) createSingleRule(network *ReteNetwork, ruleID string, exprMap map[string]interface{}, storage Storage) error {
	// Extraire l'action
	actionData, hasAction := exprMap["action"]
	if !hasAction {
		return fmt.Errorf("aucune action trouvée pour règle %s", ruleID)
	}

	actionMap, ok := actionData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("format action invalide pour règle %s: %T", ruleID, actionData)
	}

	// Créer l'action RETE
	action := cp.createAction(actionMap)

	// Analyser les contraintes pour détecter les négations et EXISTS
	constraintsData, hasConstraints := exprMap["constraints"]
	var condition map[string]interface{}
	var isExistsConstraint bool
	var hasAggregation bool

	if hasConstraints {
		// Vérifier si la contrainte contient une agrégation
		if constraintStr := fmt.Sprintf("%v", constraintsData); constraintStr != "" {
			hasAggregation = strings.Contains(constraintStr, "AVG") ||
				strings.Contains(constraintStr, "SUM") ||
				strings.Contains(constraintStr, "COUNT") ||
				strings.Contains(constraintStr, "MIN") ||
				strings.Contains(constraintStr, "MAX") ||
				strings.Contains(constraintStr, "ACCUMULATE")
		}

		// Si c'est une agrégation, créer un passthrough
		if hasAggregation {
			condition = map[string]interface{}{
				"type": "passthrough",
			}
		} else {
			// Analyser et créer la condition appropriée
			isNegation, negatedCondition, err := cp.analyzeConstraints(constraintsData)
			if err != nil {
				return fmt.Errorf("erreur analyse contraintes pour règle %s: %w", ruleID, err)
			}

			// Vérifier si c'est une contrainte EXISTS
			if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
				if constraintType, exists := constraintMap["type"].(string); exists && constraintType == "existsConstraint" {
					isExistsConstraint = true
				}
			}

			if isNegation {
				fmt.Printf("   🚫 Détection contrainte NOT - création d'un AlphaNode de négation\n")
				condition = map[string]interface{}{
					"type":      "negation",
					"negated":   true,
					"condition": negatedCondition,
				}
			} else {
				condition = map[string]interface{}{
					"type":       "constraint",
					"constraint": constraintsData,
				}
			}
		}
	} else {
		condition = map[string]interface{}{
			"type": "simple",
		}
	}

	// Analyser les variables pour déterminer si c'est une jointure ou une règle Alpha
	variables := []map[string]interface{}{}
	variableNames := []string{}
	variableTypes := []string{}

	if setData, hasSet := exprMap["set"]; hasSet {
		if setMap, ok := setData.(map[string]interface{}); ok {
			if varsData, hasVars := setMap["variables"]; hasVars {
				if varsList, ok := varsData.([]interface{}); ok && len(varsList) > 0 {
					// Extraire toutes les variables
					for _, varInterface := range varsList {
						if varMap, ok := varInterface.(map[string]interface{}); ok {
							variables = append(variables, varMap)

							if name, ok := varMap["name"].(string); ok {
								variableNames = append(variableNames, name)
							}

							// Extraire le type de la variable
							var varType string
							if dataType, ok := varMap["dataType"].(string); ok {
								varType = dataType
							} else if typeField, ok := varMap["type"].(string); ok {
								varType = typeField
							}
							variableTypes = append(variableTypes, varType)
						}
					}
				}
			}
		}
	}

	// Si c'est une contrainte EXISTS, créer un ExistsNode
	if isExistsConstraint {
		return cp.createExistsRule(network, ruleID, exprMap, condition, action, storage)
	}

	// Si c'est une agrégation, forcer la création d'un JoinNode même avec 1 variable
	if hasAggregation {

		// Extraire les informations d'agrégation
		aggInfo, err := cp.extractAggregationInfo(constraintsData)
		if err != nil {
			fmt.Printf("   ⚠️  Impossible d'extraire info agrégation: %v, utilisation JoinNode standard\n", err)
			return cp.createJoinRule(network, ruleID, variables, variableNames, variableTypes, condition, action, storage)
		}

		return cp.createAccumulatorRule(network, ruleID, variables, variableNames, variableTypes, aggInfo, action, storage)
	}

	// Si plus d'une variable, c'est une jointure Beta - créer un JoinNode
	if len(variables) > 1 {
		fmt.Printf("   📍 Règle multi-variables détectée (%d variables): %v\n", len(variables), variableNames)

		return cp.createJoinRule(network, ruleID, variables, variableNames, variableTypes, condition, action, storage)
	}

	// Sinon, traitement Alpha normal avec une seule variable
	variableName := "p" // défaut
	variableType := ""

	if len(variables) > 0 {
		if name, ok := variables[0]["name"].(string); ok {
			variableName = name
		}
		variableType = variableTypes[0]
	}

	// Créer un nœud Alpha avec la condition appropriée
	alphaNode := NewAlphaNode(ruleID+"_alpha", condition, variableName, storage)

	// Connecter seulement au type node correspondant selon le type de variable
	if variableType != "" {
		// Les TypeNodes sont stockés avec leur nom direct, pas avec "type_" préfixe
		if typeNode, exists := network.TypeNodes[variableType]; exists {
			typeNode.AddChild(alphaNode)
			fmt.Printf("   ✓ AlphaNode %s connecté au TypeNode %s\n", alphaNode.ID, variableType)
		} else {
			fmt.Printf("   ⚠️  TypeNode %s non trouvé pour variable %s\n", variableType, variableName)
			// Fallback: connecter au premier type node trouvé
			for _, typeNode := range network.TypeNodes {
				typeNode.AddChild(alphaNode)
				break
			}
		}
	} else {
		fmt.Printf("   ⚠️  Type de variable non trouvé pour %s, fallback\n", variableName)
		// Fallback: connecter au premier type node trouvé
		for _, typeNode := range network.TypeNodes {
			typeNode.AddChild(alphaNode)
			break
		}
	}
	network.AlphaNodes[alphaNode.ID] = alphaNode

	// Créer le terminal
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	alphaNode.AddChild(terminalNode)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	if condition["type"] == "negation" {
		fmt.Printf("   ✓ AlphaNode de négation créé: %s -> %s\n", alphaNode.ID, terminalNode.ID)
	}

	return nil
}

// analyzeConstraints analyse les contraintes pour détecter les négations
func (cp *ConstraintPipeline) analyzeConstraints(constraints interface{}) (bool, interface{}, error) {
	constraintMap, ok := constraints.(map[string]interface{})
	if !ok {
		return false, nil, fmt.Errorf("format contraintes invalide: %T", constraints)
	}

	// Vérifier si c'est une contrainte NOT
	if constraintType, hasType := constraintMap["type"]; hasType {
		if constraintType == "notConstraint" {
			// Extraire l'expression niée
			if expression, hasExpr := constraintMap["expression"]; hasExpr {
				fmt.Printf("   📍 Contrainte NOT détectée: %+v\n", expression)
				return true, expression, nil
			}
		}
		if constraintType == "existsConstraint" {
			// Contrainte EXISTS détectée - créer un ExistsNode
			fmt.Printf("   📍 Contrainte EXISTS détectée: %+v\n", constraintMap)
			return true, constraintMap, nil // Marquer comme complexe pour traitement spécial
		}
	}

	// Pour les autres types de contraintes, retourner false
	return false, nil, nil
}

// createAction crée une action RETE à partir d'un map parsé
func (cp *ConstraintPipeline) createAction(actionMap map[string]interface{}) *Action {
	actionName := "default_action"
	var args []interface{}

	// Extraire les données du job depuis la structure PEG: action.job.name et action.job.args
	if jobData, hasJob := actionMap["job"]; hasJob {
		if jobMap, ok := jobData.(map[string]interface{}); ok {
			// Extraire le nom de l'action depuis job.name
			if nameData, hasName := jobMap["name"]; hasName {
				if name, ok := nameData.(string); ok {
					actionName = name
				}
			}

			// Extraire les arguments depuis job.args (maintenant []interface{})
			if argsData, hasArgs := jobMap["args"]; hasArgs {
				if argsList, ok := argsData.([]interface{}); ok {
					args = argsList // Garder les objets complexes
				}
			}
		}
	}

	return &Action{
		Type: "action",
		Job: JobCall{
			Name: actionName,
			Args: args,
		},
	}
}

// BuildNetworkFromConstraintFileWithFacts construit un réseau RETE et injecte des faits massifs
func (cp *ConstraintPipeline) BuildNetworkFromConstraintFileWithFacts(constraintFile, factsFile string, storage Storage) (*ReteNetwork, []*Fact, error) {
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichier contraintes: %s\n", constraintFile)
	fmt.Printf("📁 Fichier faits: %s\n", factsFile)

	// Étape 1-4: Construction du réseau RETE normal
	network, err := cp.BuildNetworkFromConstraintFile(constraintFile, storage)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur construction réseau RETE: %w", err)
	}

	// Parser les faits avec la grammaire PEG étendue
	factsResult, err := constraint.ParseFactsFile(factsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur parsing fichier faits: %w", err)
	}

	// Extraire les faits parsés
	parsedFactsData, err := constraint.ExtractFactsFromProgram(factsResult)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur extraction faits: %w", err)
	}

	// Convertir au format RETE Facts
	facts := make([]*Fact, 0, len(parsedFactsData))
	for _, factData := range parsedFactsData {
		fact := &Fact{
			ID:     factData["id"].(string),
			Type:   factData["reteType"].(string), // Utiliser le type RETE
			Fields: make(map[string]interface{}),
		}

		// Copier tous les champs sauf id et reteType
		for key, value := range factData {
			if key != "id" && key != "reteType" {
				fact.Fields[key] = value
			}
		}

		facts = append(facts, fact)
	}

	fmt.Printf("✅ %d faits parsés et validés\n", len(facts))

	// Injecter tous les faits
	successCount := 0
	errorCount := 0

	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			errorCount++
			// Log des erreurs mais continuer
			fmt.Printf("⚠️ Erreur injection fait %s: %v\n", fact.ID, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("✅ Injection terminée: %d succès, %d erreurs\n", successCount, errorCount)
	fmt.Printf("🎯 PIPELINE CONSTRAINT + FAITS TERMINÉ\n")
	fmt.Printf("========================================\n\n")

	return network, facts, nil
}

// validateNetwork effectue une validation basique du réseau construit
func (cp *ConstraintPipeline) validateNetwork(network *ReteNetwork) error {
	if len(network.TypeNodes) == 0 {
		return fmt.Errorf("aucun type défini dans le réseau")
	}

	if len(network.TerminalNodes) == 0 {
		return fmt.Errorf("aucune règle définie dans le réseau")
	}

	// Validation additionnelle
	for typeName, typeNode := range network.TypeNodes {
		if typeNode == nil {
			return fmt.Errorf("type node null pour %s", typeName)
		}
	}

	for terminalID, terminal := range network.TerminalNodes {
		if terminal == nil {
			return fmt.Errorf("terminal node null pour %s", terminalID)
		}
		if terminal.Action == nil {
			return fmt.Errorf("action manquante pour terminal %s", terminalID)
		}
	}

	return nil
}

// getStringField extrait un champ string d'un map avec valeur par défaut
func getStringField(m map[string]interface{}, key, defaultValue string) string {
	if value, exists := m[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// createJoinRule crée une règle Beta avec JoinNode pour les règles multi-variables
func (cp *ConstraintPipeline) createJoinRule(network *ReteNetwork, ruleID string, variables []map[string]interface{}, variableNames []string, variableTypes []string, condition map[string]interface{}, action *Action, storage Storage) error {

	// Créer le nœud terminal pour cette règle
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Créer le JoinNode
	leftVars := []string{variableNames[0]} // Variable primaire
	rightVars := variableNames[1:]         // Variables secondaires

	// Créer le mapping variable -> type
	varTypes := make(map[string]string)
	for i, varName := range variableNames {
		varTypes[varName] = variableTypes[i]
	}

	joinNode := NewJoinNode(ruleID+"_join", condition, leftVars, rightVars, varTypes, storage)
	joinNode.AddChild(terminalNode)

	// Stocker le JoinNode dans les BetaNodes du réseau
	network.BetaNodes[joinNode.ID] = joinNode

	// Créer des AlphaNodes pass-through qui ne filtrent pas mais transfèrent vers JoinNode
	for i, varName := range variableNames {
		varType := variableTypes[i]
		if varType != "" {
			if typeNode, exists := network.TypeNodes[varType]; exists {

				// Déterminer le côté (gauche/droite) selon l'architecture RETE
				side := "right"
				if i == 0 {
					side = "left" // Première variable va vers la gauche
				}

				// Créer un AlphaNode pass-through avec indication de côté
				passCondition := map[string]interface{}{
					"type": "passthrough",
					"side": side, // "left" ou "right"
				}
				alphaNode := NewAlphaNode(ruleID+"_pass_"+varName, passCondition, varName, storage)

				// Connecter TypeNode -> AlphaPassthrough -> JoinNode
				typeNode.AddChild(alphaNode)
				alphaNode.AddChild(joinNode)

				fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> JoinNode_%s\n", varType, varName, ruleID)
			} else {
				fmt.Printf("   ⚠️ TypeNode %s introuvable!\n", varType)
			}
		} else {
			fmt.Printf("   ⚠️ Type vide pour variable %s\n", varName)
		}
	}

	fmt.Printf("   ✅ JoinNode %s créé pour jointure %s\n", joinNode.ID, strings.Join(variableNames, " ⋈ "))
	return nil
}

// createExistsRule crée une règle EXISTS avec ExistsNode
func (cp *ConstraintPipeline) createExistsRule(network *ReteNetwork, ruleID string, exprMap map[string]interface{}, condition map[string]interface{}, action *Action, storage Storage) error {
	// Créer le nœud terminal pour cette règle
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Extraire les variables depuis l'expression
	var mainVariable, existsVariable string
	var mainVarType, existsVarType string

	// Extraire la variable principale depuis "set"
	if setData, hasSet := exprMap["set"]; hasSet {
		if setMap, ok := setData.(map[string]interface{}); ok {
			if varsData, hasVars := setMap["variables"]; hasVars {
				if varsList, ok := varsData.([]interface{}); ok && len(varsList) > 0 {
					if varMap, ok := varsList[0].(map[string]interface{}); ok {
						if name, ok := varMap["name"].(string); ok {
							mainVariable = name
						}
						if dataType, ok := varMap["dataType"].(string); ok {
							mainVarType = dataType
						}
					}
				}
			}
		}
	}

	// Extraire la variable d'existence depuis les contraintes
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			if variable, hasVar := constraintMap["variable"]; hasVar {
				if varMap, ok := variable.(map[string]interface{}); ok {
					if name, ok := varMap["name"].(string); ok {
						existsVariable = name
					}
					if dataType, ok := varMap["dataType"].(string); ok {
						existsVarType = dataType
					}
				}
			}
		}
	}

	if mainVariable == "" || existsVariable == "" {
		return fmt.Errorf("variables EXISTS non trouvées: main=%s, exists=%s", mainVariable, existsVariable)
	}

	// Extraire les conditions d'EXISTS depuis exprMap["constraints"]["condition"]
	var existsConditions []map[string]interface{}
	if constraintsData, hasConstraints := exprMap["constraints"]; hasConstraints {
		if constraintMap, ok := constraintsData.(map[string]interface{}); ok {
			// Essayer d'abord "condition" (au singulier)
			if conditionData, hasCondition := constraintMap["condition"]; hasCondition {
				if conditionObj, ok := conditionData.(map[string]interface{}); ok {
					existsConditions = append(existsConditions, conditionObj)
				}
			}
			// Puis essayer "conditions" (au pluriel) si pas trouvé
			if len(existsConditions) == 0 {
				if conditionsData, hasConditions := constraintMap["conditions"]; hasConditions {
					if conditionsList, ok := conditionsData.([]interface{}); ok {
						for _, conditionData := range conditionsList {
							if conditionObj, ok := conditionData.(map[string]interface{}); ok {
								existsConditions = append(existsConditions, conditionObj)
							}
						}
					}
				}
			}
		}
	}

	// Créer l'objet condition pour l'ExistsNode avec les conditions extraites
	existsConditionObj := map[string]interface{}{
		"type":       "exists",
		"conditions": existsConditions,
	}

	// Créer le mapping variable -> type pour l'ExistsNode
	varTypes := make(map[string]string)
	varTypes[mainVariable] = mainVarType
	varTypes[existsVariable] = existsVarType

	// Créer l'ExistsNode avec les vraies conditions
	existsNode := NewExistsNode(ruleID+"_exists", existsConditionObj, mainVariable, existsVariable, varTypes, storage)
	existsNode.AddChild(terminalNode)

	// Stocker l'ExistsNode dans les BetaNodes du réseau
	network.BetaNodes[existsNode.ID] = existsNode

	// Créer des AlphaNodes pass-through pour les deux variables
	// Variable principale → ActivateLeft
	if mainVarType != "" {
		if typeNode, exists := network.TypeNodes[mainVarType]; exists {
			mainAlphaCondition := map[string]interface{}{
				"type": "passthrough",
				"side": "left",
			}
			mainAlphaNode := NewAlphaNode(ruleID+"_pass_"+mainVariable, mainAlphaCondition, mainVariable, storage)

			typeNode.AddChild(mainAlphaNode)
			mainAlphaNode.AddChild(existsNode)

			fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> ExistsNode_%s (LEFT)\n", mainVarType, mainVariable, ruleID)
		}
	}

	// Variable d'existence → ActivateRight
	if existsVarType != "" {
		if typeNode, exists := network.TypeNodes[existsVarType]; exists {
			existsAlphaCondition := map[string]interface{}{
				"type": "passthrough",
				"side": "right",
			}
			existsAlphaNode := NewAlphaNode(ruleID+"_pass_"+existsVariable, existsAlphaCondition, existsVariable, storage)

			typeNode.AddChild(existsAlphaNode)
			existsAlphaNode.AddChild(existsNode)

			fmt.Printf("   ✓ %s -> PassthroughAlpha_%s -> ExistsNode_%s (RIGHT)\n", existsVarType, existsVariable, ruleID)
		}
	}

	fmt.Printf("   ✅ ExistsNode %s créé pour %s EXISTS %s\n", existsNode.ID, mainVariable, existsVariable)
	return nil
}

// extractAggregationInfo extrait les informations d'une agrégation depuis les contraintes
func (cp *ConstraintPipeline) extractAggregationInfo(constraintsData interface{}) (*AggregationInfo, error) {
	constraintMap, ok := constraintsData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("contrainte invalide: %T", constraintsData)
	}

	// Vérifier que c'est bien une accumulateConstraint
	constraintType, ok := constraintMap["type"].(string)
	if !ok || constraintType != "accumulateConstraint" {
		return nil, fmt.Errorf("pas une accumulateConstraint: %s", constraintType)
	}

	info := &AggregationInfo{}

	// Extraire la fonction d'agrégation
	function, ok := constraintMap["function"].(string)
	if !ok {
		return nil, fmt.Errorf("fonction d'agrégation manquante")
	}
	info.Function = function

	// Extraire la variable à agréger (ex: p: Performance)
	variableData, ok := constraintMap["variable"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("variable d'agrégation manquante")
	}

	aggVarName, ok := variableData["name"].(string)
	if !ok {
		return nil, fmt.Errorf("nom de variable d'agrégation manquant")
	}
	info.AggVariable = aggVarName

	aggVarType, ok := variableData["dataType"].(string)
	if !ok {
		return nil, fmt.Errorf("type de variable d'agrégation manquant")
	}
	info.AggType = aggVarType

	// Extraire le champ à agréger (ex: score)
	if field, ok := constraintMap["field"].(string); ok {
		info.Field = field
	}
	// Pour COUNT, le champ peut être vide

	// Extraire l'opérateur et le seuil
	operator, ok := constraintMap["operator"].(string)
	if !ok {
		return nil, fmt.Errorf("opérateur manquant")
	}
	info.Operator = operator

	thresholdData, ok := constraintMap["threshold"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("seuil manquant")
	}
	threshold, ok := thresholdData["value"].(float64)
	if !ok {
		return nil, fmt.Errorf("valeur de seuil invalide")
	}
	info.Threshold = threshold

	// Extraire la condition de jointure (ex: p.employee_id == e.id)
	conditionData, ok := constraintMap["condition"].(map[string]interface{})
	if ok {
		info.JoinCondition = conditionData

		// Extraire les champs de jointure depuis la condition
		if condType, ok := conditionData["type"].(string); ok && condType == "comparison" {
			// Left side: p.employee_id
			if leftData, ok := conditionData["left"].(map[string]interface{}); ok {
				if leftType, ok := leftData["type"].(string); ok && leftType == "fieldAccess" {
					if joinField, ok := leftData["field"].(string); ok {
						info.JoinField = joinField // "employee_id"
					}
				}
			}

			// Right side: e.id
			if rightData, ok := conditionData["right"].(map[string]interface{}); ok {
				if rightType, ok := rightData["type"].(string); ok && rightType == "fieldAccess" {
					if mainField, ok := rightData["field"].(string); ok {
						info.MainField = mainField // "id"
					}
				}
			}
		}
	}

	return info, nil
}

// createAccumulatorRule crée une règle avec AccumulatorNode
func (cp *ConstraintPipeline) createAccumulatorRule(network *ReteNetwork, ruleID string, variables []map[string]interface{}, variableNames []string, variableTypes []string, aggInfo *AggregationInfo, action *Action, storage Storage) error {

	// Extraire la variable principale et son type depuis variables
	if len(variables) == 0 || len(variableTypes) == 0 {
		return fmt.Errorf("aucune variable principale trouvée")
	}

	mainVariable := variableNames[0]
	mainType := variableTypes[0]

	// Stocker dans aggInfo
	aggInfo.MainVariable = mainVariable
	aggInfo.MainType = mainType

	// Créer le nœud terminal
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Créer la condition de comparaison
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": aggInfo.Operator,
		"value":    aggInfo.Threshold,
	}

	// Créer l'AccumulatorNode avec tous les paramètres
	accumNode := NewAccumulatorNode(
		ruleID+"_accum",
		aggInfo.MainVariable, // "e"
		aggInfo.MainType,     // "Employee"
		aggInfo.AggVariable,  // "p"
		aggInfo.AggType,      // "Performance"
		aggInfo.Field,        // "score"
		aggInfo.JoinField,    // "employee_id"
		aggInfo.MainField,    // "id"
		aggInfo.Function,     // "AVG"
		condition,
		storage,
	)
	accumNode.AddChild(terminalNode)
	network.BetaNodes[accumNode.ID] = accumNode

	// Connecter le TypeNode principal (Employee) à l'AccumulatorNode
	if typeNode, exists := network.TypeNodes[mainType]; exists {
		// Créer un AlphaNode passthrough pour la variable principale
		passCondition := map[string]interface{}{
			"type": "passthrough",
		}
		alphaNode := NewAlphaNode(ruleID+"_pass_"+mainVariable, passCondition, mainVariable, storage)

		typeNode.AddChild(alphaNode)
		alphaNode.AddChild(accumNode)

		fmt.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s]\n", mainType, aggInfo.Function)
	}

	// CRUCIAL: Connecter aussi le TypeNode des faits à agréger (Performance) à l'AccumulatorNode
	if aggTypeNode, exists := network.TypeNodes[aggInfo.AggType]; exists {
		// Créer un AlphaNode passthrough pour la variable d'agrégation
		passConditionAgg := map[string]interface{}{
			"type": "passthrough",
		}
		alphaNodeAgg := NewAlphaNode(ruleID+"_pass_"+aggInfo.AggVariable, passConditionAgg, aggInfo.AggVariable, storage)

		aggTypeNode.AddChild(alphaNodeAgg)
		alphaNodeAgg.AddChild(accumNode)

		fmt.Printf("   ✓ %s -> PassthroughAlpha -> AccumulatorNode[%s] (pour agrégation)\n", aggInfo.AggType, aggInfo.Function)
	}

	fmt.Printf("   ✅ AccumulatorNode %s créé pour %s(%s.%s) %s %.2f\n",
		accumNode.ID, aggInfo.Function, aggInfo.AggVariable, aggInfo.Field, aggInfo.Operator, aggInfo.Threshold)
	return nil
}
