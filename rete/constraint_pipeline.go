// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

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
// Si le fichier contient des instructions reset, utilise l'IterativeParser pour appliquer
// correctement la sémantique de reset (seuls les types/règles après le dernier reset sont conservés)
func (cp *ConstraintPipeline) BuildNetworkFromConstraintFile(constraintFile string, storage Storage) (*ReteNetwork, error) {
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichier: %s\n", constraintFile)

	// ÉTAPE 1: Parsing initial pour détecter les resets
	parsedAST, err := constraint.ParseConstraintFile(constraintFile)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur parsing fichier %s: %w", constraintFile, err)
	}
	fmt.Printf("✅ Parsing réussi\n")

	// ÉTAPE 1.5: Validation sémantique du programme
	err = constraint.ValidateConstraintProgram(parsedAST)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur validation sémantique: %w", err)
	}
	fmt.Printf("✅ Validation sémantique réussie\n")

	// Valider que c'est un map[string]interface{}
	resultMap, ok := parsedAST.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ Format AST non reconnu: %T", parsedAST)
	}

	// Vérifier si le fichier contient des instructions reset
	hasResets := false
	if resetsData, exists := resultMap["resets"]; exists {
		if resets, ok := resetsData.([]interface{}); ok && len(resets) > 0 {
			hasResets = true
			fmt.Printf("⚠️  Instructions reset détectées (%d) - Utilisation de l'IterativeParser\n", len(resets))
		}
	}

	// Si des resets sont présents, utiliser l'IterativeParser pour appliquer la sémantique correcte
	if hasResets {
		return cp.buildNetworkWithResetSemantics(constraintFile, storage)
	}

	// ÉTAPE 2: Extraction et validation des composants (cas sans reset)
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

// buildNetworkWithResetSemantics construit un réseau en appliquant correctement la sémantique reset
// Analyse le fichier pour déterminer quels types/expressions viennent après le dernier reset
func (cp *ConstraintPipeline) buildNetworkWithResetSemantics(constraintFile string, storage Storage) (*ReteNetwork, error) {
	// Lire le fichier pour analyser la structure
	fileContent, err := constraint.ReadFileContent(constraintFile)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur lecture fichier: %w", err)
	}

	// Parser le fichier complet pour obtenir tous les éléments
	parsedAST, err := constraint.ParseConstraintFile(constraintFile)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur parsing: %w", err)
	}

	// Valider
	err = constraint.ValidateConstraintProgram(parsedAST)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur validation: %w", err)
	}

	// Convertir en programme
	program, err := constraint.ConvertResultToProgram(parsedAST)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur conversion: %w", err)
	}

	// Analyser où se trouve le dernier reset dans le fichier
	lastResetPosition := cp.findLastResetPosition(fileContent)

	// Filtrer les types et expressions pour ne garder que ceux après le dernier reset
	filteredTypes, filteredExpressions := cp.filterAfterReset(
		program.Types, program.Expressions, fileContent, lastResetPosition)

	fmt.Printf("✅ Après application des resets: %d type(s), %d expression(s)\n",
		len(filteredTypes), len(filteredExpressions))

	// Convertir au format RETE
	filteredProgram := &constraint.Program{
		Types:       filteredTypes,
		Expressions: filteredExpressions,
		Facts:       []constraint.Fact{},  // Les faits seront ajoutés séparément
		Resets:      []constraint.Reset{}, // Plus de resets après filtrage
	}

	reteProgram := constraint.ConvertToReteProgram(filteredProgram)
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

	// Validation finale
	err = cp.validateNetwork(network)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur validation réseau: %w", err)
	}
	fmt.Printf("✅ Validation réussie\n")

	fmt.Printf("🎯 PIPELINE AVEC RESET TERMINÉ AVEC SUCCÈS\n")
	fmt.Printf("========================================\n\n")

	return network, nil
}

// findLastResetPosition trouve la position du dernier mot "reset" dans le fichier
// Retourne la ligne (0-based) où se trouve le dernier reset, ou -1 si aucun
func (cp *ConstraintPipeline) findLastResetPosition(content string) int {
	lines := strings.Split(content, "\n")
	lastResetLine := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "reset" {
			lastResetLine = i
		}
	}

	return lastResetLine
}

// filterAfterReset filtre les types et expressions pour ne garder que ceux définis après le reset
// Stratégie: compte combien de définitions de types et d'expressions apparaissent avant le reset
// dans le fichier source, puis ne garde que les éléments après ces positions dans les slices
func (cp *ConstraintPipeline) filterAfterReset(
	types []constraint.TypeDefinition,
	expressions []constraint.Expression,
	fileContent string,
	resetLine int,
) ([]constraint.TypeDefinition, []constraint.Expression) {

	if resetLine < 0 {
		// Pas de reset trouvé, retourner tout
		return types, expressions
	}

	lines := strings.Split(fileContent, "\n")

	// Compter combien de "type " apparaissent avant le reset
	typesBeforeReset := 0
	for i := 0; i < resetLine && i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "type ") && strings.Contains(trimmed, ":") {
			typesBeforeReset++
		}
	}

	// Compter combien de règles (lignes avec "==>") apparaissent avant le reset
	expressionsBeforeReset := 0
	for i := 0; i < resetLine && i < len(lines); i++ {
		if strings.Contains(lines[i], "==>") {
			expressionsBeforeReset++
		}
	}

	// Filtrer les types: garder seulement ceux après l'index typesBeforeReset
	var filteredTypes []constraint.TypeDefinition
	if typesBeforeReset < len(types) {
		filteredTypes = types[typesBeforeReset:]
	}

	// Filtrer les expressions: garder seulement celles après l'index expressionsBeforeReset
	var filteredExpressions []constraint.Expression
	if expressionsBeforeReset < len(expressions) {
		filteredExpressions = expressions[expressionsBeforeReset:]
	}

	return filteredTypes, filteredExpressions
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

// BuildNetworkFromConstraintFileWithFacts construit un réseau et soumet immédiatement des faits
func (cp *ConstraintPipeline) BuildNetworkFromConstraintFileWithFacts(constraintFile, factsFile string, storage Storage) (*ReteNetwork, []*Fact, error) {
	fmt.Printf("========================================\n")
	fmt.Printf("📁 Fichier contraintes: %s\n", constraintFile)
	fmt.Printf("📁 Fichier faits: %s\n", factsFile)

	// ÉTAPE 1: Construire le réseau depuis le fichier de contraintes
	network, err := cp.BuildNetworkFromConstraintFile(constraintFile, storage)
	if err != nil {
		return nil, nil, fmt.Errorf("❌ Erreur construction réseau: %w", err)
	}

	// ÉTAPE 2: Parser et soumettre les faits
	fmt.Printf("📊 Parsing des faits depuis %s\n", factsFile)

	parsedFacts, err := constraint.ParseFactsFile(factsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("❌ Erreur parsing faits: %w", err)
	}

	// Extraire les faits du programme parsé
	factsList, err := constraint.ExtractFactsFromProgram(parsedFacts)
	if err != nil {
		return nil, nil, fmt.Errorf("❌ Erreur extraction faits: %w", err)
	}

	// Convertir et soumettre chaque fait
	submittedFacts := []*Fact{}
	for _, factMap := range factsList {
		// ExtractFactsFromProgram retourne des maps avec 'reteType' et tous les champs directement
		factID := getStringField(factMap, "id", "")
		factType := getStringField(factMap, "reteType", "") // Utiliser 'reteType' au lieu de 'type'

		if factID == "" || factType == "" {
			fmt.Printf("⚠️ Fait ignoré: id='%s', type='%s'\n", factID, factType)
			continue
		}

		// Les champs sont directement dans factMap (pas de sous-clé 'fields')
		fields := make(map[string]interface{})
		for key, value := range factMap {
			// Exclure les métadonnées RETE (id, reteType)
			if key != "id" && key != "reteType" {
				fields[key] = value
			}
		}

		fact := &Fact{
			ID:     factID,
			Type:   factType,
			Fields: fields,
		}

		err := network.SubmitFact(fact)
		if err != nil {
			fmt.Printf("⚠️ Erreur soumission fait %s: %v\n", factID, err)
		}
		submittedFacts = append(submittedFacts, fact)
	}

	fmt.Printf("✅ %d faits soumis au réseau\n", len(submittedFacts))
	fmt.Printf("🎯 PIPELINE AVEC FAITS TERMINÉ\n")
	fmt.Printf("========================================\n\n")

	return network, submittedFacts, nil
}
