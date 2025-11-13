package rete

import (
	"fmt"

	"github.com/treivax/tsd/constraint"
)

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
	fmt.Printf("🔧 PIPELINE CONSTRAINT → RETE\n")
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichier: %s\n", constraintFile)

	// ÉTAPE 1: Parsing avec le vrai parseur PEG
	fmt.Printf("🔍 Étape 1/4: Parsing PEG du fichier .constraint...\n")
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
	fmt.Printf("🔍 Étape 2/4: Extraction types et expressions...\n")
	types, expressions, err := cp.extractComponents(resultMap)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur extraction composants: %w", err)
	}
	fmt.Printf("✅ Trouvé %d types et %d expressions\n", len(types), len(expressions))

	// ÉTAPE 3: Construction du réseau RETE
	fmt.Printf("🔍 Étape 3/4: Construction réseau RETE...\n")
	network, err := cp.buildNetwork(storage, types, expressions)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur construction réseau: %w", err)
	}
	fmt.Printf("✅ Réseau construit avec %d nœuds terminaux\n", len(network.TerminalNodes))

	// ÉTAPE 4: Validation finale
	fmt.Printf("🔍 Étape 4/4: Validation réseau...\n")
	err = cp.validateNetwork(network)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur validation réseau: %w", err)
	}
	fmt.Printf("✅ Validation réussie\n")

	fmt.Printf("🎯 PIPELINE TERMINÉ AVEC SUCCÈS\n")
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

// createSingleRule crée une règle unique (Alpha + Terminal pour l'instant)
func (cp *ConstraintPipeline) createSingleRule(network *ReteNetwork, ruleID string, exprMap map[string]interface{}, storage Storage) error {
	// Pour cette implémentation de base, on crée des nœuds Alpha simples
	// Une implémentation future pourrait analyser les contraintes pour créer des Beta joints
	
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

	// Créer un nœud Alpha simple pour cette règle
	// TODO: Analyser les contraintes pour déterminer le type exact de nœud
	condition := map[string]interface{}{
		"type": "simple", // Condition simple pour ce pipeline de base
	}

	alphaNode := NewAlphaNode(ruleID+"_alpha", condition, "x", storage)
	
	// Connecter à un type node (prendre le premier disponible pour l'instant)
	if len(network.TypeNodes) > 0 {
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

	return nil
}

// createAction crée une action RETE à partir d'un map parsé
func (cp *ConstraintPipeline) createAction(actionMap map[string]interface{}) *Action {
	actionName := "default_action"
	var args []string

	// Extraire les données du job depuis la structure PEG: action.job.name et action.job.args
	if jobData, hasJob := actionMap["job"]; hasJob {
		if jobMap, ok := jobData.(map[string]interface{}); ok {
			// Extraire le nom de l'action depuis job.name
			if nameData, hasName := jobMap["name"]; hasName {
				if name, ok := nameData.(string); ok {
					actionName = name
				}
			}
			
			// Extraire les arguments depuis job.args
			if argsData, hasArgs := jobMap["args"]; hasArgs {
				if argsList, ok := argsData.([]interface{}); ok {
					for _, arg := range argsList {
						if argStr, ok := arg.(string); ok {
							args = append(args, argStr)
						}
					}
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
	fmt.Printf("🔧 PIPELINE CONSTRAINT + FAITS → RETE\n")
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichier contraintes: %s\n", constraintFile)
	fmt.Printf("📁 Fichier faits: %s\n", factsFile)
	
	// Étape 1-4: Construction du réseau RETE normal
	network, err := cp.BuildNetworkFromConstraintFile(constraintFile, storage)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur construction réseau RETE: %w", err)
	}
	
	fmt.Printf("\n🔍 Étape 5/6: Parsing et validation fichier faits...\n")
	
	// Extraire les définitions de types du réseau pour validation des faits
	typeDefinitions := make(map[string]TypeDefinition)
	for typeName, typeNode := range network.TypeNodes {
		typeDefinitions[typeName] = typeNode.TypeDefinition
	}
	
	// Parser les faits
	factsParser := NewFactsParser()
	facts, err := factsParser.ParseFactsFile(factsFile, typeDefinitions)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur parsing faits: %w", err)
	}
	
	// Afficher les métadonnées du fichier faits
	metadata := factsParser.GetMetadata()
	if len(metadata) > 0 {
		fmt.Printf("📋 Métadonnées fichier faits:\n")
		for key, value := range metadata {
			fmt.Printf("   %s: %s\n", key, value)
		}
	}
	
	fmt.Printf("✅ %d faits parsés et validés\n", len(facts))
	
	fmt.Printf("\n🔍 Étape 6/6: Injection des faits dans le réseau RETE...\n")
	
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