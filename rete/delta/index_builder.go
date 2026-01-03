// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"reflect"
)

// IndexBuilder construit un DependencyIndex depuis un réseau RETE complet.
//
// Cette structure orchestre l'extraction de champs depuis tous les nœuds
// et la construction de l'index de dépendances.
type IndexBuilder struct {
	alphaExtractor    *AlphaConditionExtractor
	betaExtractor     *BetaConditionExtractor
	actionExtractor   *ActionFieldExtractor
	enableDiagnostics bool
	diagnostics       *BuildDiagnostics
}

// BuildDiagnostics contient des informations de diagnostic sur la construction.
type BuildDiagnostics struct {
	NodesProcessed  int
	NodesSkipped    int
	FieldsExtracted int
	Errors          []string
	Warnings        []string
}

// NewIndexBuilder crée un nouveau constructeur d'index.
func NewIndexBuilder() *IndexBuilder {
	return &IndexBuilder{
		alphaExtractor:    &AlphaConditionExtractor{},
		betaExtractor:     &BetaConditionExtractor{},
		actionExtractor:   &ActionFieldExtractor{},
		enableDiagnostics: false,
		diagnostics:       &BuildDiagnostics{},
	}
}

// EnableDiagnostics active la collecte d'informations de diagnostic.
func (ib *IndexBuilder) EnableDiagnostics() {
	ib.enableDiagnostics = true
}

// GetDiagnostics retourne les diagnostics de la dernière construction.
func (ib *IndexBuilder) GetDiagnostics() *BuildDiagnostics {
	return ib.diagnostics
}

// BuildFromNetwork construit un index depuis un réseau RETE.
//
// Cette méthode parcourt tous les nœuds du réseau, extrait les champs
// utilisés, et construit l'index de dépendances.
//
// Paramètres :
//   - network : interface{} représentant le ReteNetwork
//     (typiquement *rete.ReteNetwork, mais on utilise interface{} pour
//     éviter les dépendances circulaires)
//
// Retourne un DependencyIndex complet et une erreur si échec.
//
// Note : Cette fonction utilise la reflection pour accéder aux champs
// du ReteNetwork afin d'éviter les dépendances circulaires.
func (ib *IndexBuilder) BuildFromNetwork(network interface{}) (*DependencyIndex, error) {
	idx := NewDependencyIndex()

	// Reset diagnostics
	ib.diagnostics = &BuildDiagnostics{}

	if network == nil {
		return idx, nil // Index vide pour réseau nil
	}

	// Utiliser reflection pour accéder aux champs sans dépendance circulaire
	networkValue := reflect.ValueOf(network)
	if networkValue.Kind() == reflect.Ptr {
		networkValue = networkValue.Elem()
	}

	if networkValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("network must be a struct, got %v", networkValue.Kind())
	}

	// Extraire AlphaNodes
	if err := ib.extractAlphaNodes(idx, networkValue); err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Failed to extract alpha nodes: %v", err),
			)
		}
		// Continue malgré l'erreur pour extraire les autres nœuds
	}

	// Extraire BetaNodes (si présents)
	if err := ib.extractBetaNodes(idx, networkValue); err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Failed to extract beta nodes: %v", err),
			)
		}
		// Continue malgré l'erreur
	}

	// Extraire TerminalNodes
	if err := ib.extractTerminalNodes(idx, networkValue); err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Failed to extract terminal nodes: %v", err),
			)
		}
		// Continue malgré l'erreur
	}

	return idx, nil
}

// extractAlphaNodes extrait les nœuds alpha du réseau.
func (ib *IndexBuilder) extractAlphaNodes(idx *DependencyIndex, networkValue reflect.Value) error {
	alphaNodesField := networkValue.FieldByName("AlphaNodes")
	if !alphaNodesField.IsValid() || alphaNodesField.IsNil() {
		return nil // Pas d'AlphaNodes, c'est OK
	}

	// AlphaNodes est un map[string]*AlphaNode
	if alphaNodesField.Kind() != reflect.Map {
		return fmt.Errorf("AlphaNodes is not a map")
	}

	iter := alphaNodesField.MapRange()
	for iter.Next() {
		nodeID := iter.Key().String()
		nodeValue := iter.Value()

		if nodeValue.Kind() == reflect.Ptr {
			nodeValue = nodeValue.Elem()
		}

		if !nodeValue.IsValid() {
			continue
		}

		// Extraire VariableName (type du fait)
		varNameField := nodeValue.FieldByName("VariableName")
		if !varNameField.IsValid() {
			if ib.enableDiagnostics {
				ib.diagnostics.Warnings = append(
					ib.diagnostics.Warnings,
					fmt.Sprintf("Alpha node %s: no VariableName", nodeID),
				)
			}
			continue
		}
		factType := varNameField.String()

		// Extraire Condition
		conditionField := nodeValue.FieldByName("Condition")
		if !conditionField.IsValid() {
			if ib.enableDiagnostics {
				ib.diagnostics.Warnings = append(
					ib.diagnostics.Warnings,
					fmt.Sprintf("Alpha node %s: no Condition", nodeID),
				)
			}
			continue
		}

		condition := conditionField.Interface()

		// Utiliser BuildFromAlphaNode qui gère les diagnostics
		if err := ib.BuildFromAlphaNode(idx, nodeID, factType, condition); err != nil {
			// Erreur déjà loggée par BuildFromAlphaNode
			continue
		}
	}

	return nil
}

// extractBetaNodes extrait les nœuds beta du réseau.
func (ib *IndexBuilder) extractBetaNodes(idx *DependencyIndex, networkValue reflect.Value) error {
	betaNodesField := networkValue.FieldByName("BetaNodes")
	if !betaNodesField.IsValid() {
		return nil // Pas de BetaNodes, c'est OK
	}

	// Check if the field can be nil before calling IsNil()
	if betaNodesField.Kind() == reflect.Ptr || betaNodesField.Kind() == reflect.Map ||
		betaNodesField.Kind() == reflect.Slice || betaNodesField.Kind() == reflect.Chan ||
		betaNodesField.Kind() == reflect.Func || betaNodesField.Kind() == reflect.Interface {
		if betaNodesField.IsNil() {
			return nil // BetaNodes is nil, c'est OK
		}
	}

	// BetaNodes est un map[string]interface{}
	if betaNodesField.Kind() != reflect.Map {
		return fmt.Errorf("BetaNodes is not a map")
	}

	iter := betaNodesField.MapRange()
	for iter.Next() {
		nodeID := iter.Key().String()
		nodeValue := iter.Value()

		if nodeValue.Kind() == reflect.Interface {
			nodeValue = nodeValue.Elem()
		}

		if nodeValue.Kind() == reflect.Ptr {
			nodeValue = nodeValue.Elem()
		}

		if !nodeValue.IsValid() {
			continue
		}

		// Pour les nœuds beta, on essaie d'extraire le type et les join conditions
		// Note: La structure exacte des BetaNodes peut varier
		// Pour l'instant, on skip car la structure est générique (interface{})
		if ib.enableDiagnostics {
			ib.diagnostics.NodesSkipped++
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Beta node %s: skipped (generic extraction not yet implemented)", nodeID),
			)
		}
	}

	return nil
}

// extractTerminalNodes extrait les nœuds terminaux du réseau.
func (ib *IndexBuilder) extractTerminalNodes(idx *DependencyIndex, networkValue reflect.Value) error {
	terminalNodesField := networkValue.FieldByName("TerminalNodes")
	if !terminalNodesField.IsValid() || terminalNodesField.IsNil() {
		return nil // Pas de TerminalNodes, c'est OK
	}

	// TerminalNodes est un map[string]*TerminalNode
	if terminalNodesField.Kind() != reflect.Map {
		return fmt.Errorf("TerminalNodes is not a map")
	}

	iter := terminalNodesField.MapRange()
	for iter.Next() {
		nodeID := iter.Key().String()
		nodeValue := iter.Value()

		if nodeValue.Kind() == reflect.Ptr {
			nodeValue = nodeValue.Elem()
		}

		if !nodeValue.IsValid() {
			continue
		}

		// Extraire Action
		actionField := nodeValue.FieldByName("Action")
		if !actionField.IsValid() || actionField.IsNil() {
			if ib.enableDiagnostics {
				ib.diagnostics.Warnings = append(
					ib.diagnostics.Warnings,
					fmt.Sprintf("Terminal node %s: no Action", nodeID),
				)
			}
			continue
		}

		actionValue := actionField
		if actionValue.Kind() == reflect.Ptr {
			actionValue = actionValue.Elem()
		}

		// Extraire VariableName (type du fait cible)
		varNameField := actionValue.FieldByName("VariableName")
		factType := ""
		if varNameField.IsValid() {
			factType = varNameField.String()
		}

		// Si pas de VariableName dans Action, essayer de trouver le type autrement
		if factType == "" {
			// Essayer d'obtenir le type depuis le TypeName dans Action
			typeNameField := actionValue.FieldByName("TypeName")
			if typeNameField.IsValid() {
				factType = typeNameField.String()
			}
		}

		if factType == "" {
			if ib.enableDiagnostics {
				ib.diagnostics.Warnings = append(
					ib.diagnostics.Warnings,
					fmt.Sprintf("Terminal node %s: cannot determine fact type", nodeID),
				)
			}
			// Utiliser un type par défaut basé sur le nodeID
			factType = "Unknown"
		}

		// Convertir Action en interface{} pour l'extracteur
		action := actionField.Interface()

		// Pour terminal nodes, on passe l'action comme un slice d'une action
		actions := []interface{}{action}

		// Utiliser BuildFromTerminalNode qui gère les diagnostics
		if err := ib.BuildFromTerminalNode(idx, nodeID, factType, actions); err != nil {
			// Erreur déjà loggée par BuildFromTerminalNode
			continue
		}
	}

	return nil
}

// BuildFromAlphaNode traite un nœud alpha et l'ajoute à l'index.
//
// Paramètres :
//   - idx : index de dépendances à remplir
//   - nodeID : identifiant du nœud
//   - factType : type de fait
//   - condition : condition du nœud (AST)
//
// Retourne une erreur si l'extraction de champs échoue.
func (ib *IndexBuilder) BuildFromAlphaNode(
	idx *DependencyIndex,
	nodeID, factType string,
	condition interface{},
) error {
	fields, err := ib.alphaExtractor.ExtractFields(condition)
	if err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Alpha node %s: %v", nodeID, err),
			)
			ib.diagnostics.NodesSkipped++
		}
		return err
	}

	if len(fields) == 0 {
		if ib.enableDiagnostics {
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Alpha node %s: no fields extracted", nodeID),
			)
		}
	}

	idx.AddAlphaNode(nodeID, factType, fields)

	if ib.enableDiagnostics {
		ib.diagnostics.NodesProcessed++
		ib.diagnostics.FieldsExtracted += len(fields)
	}

	return nil
}

// BuildFromBetaNode traite un nœud beta et l'ajoute à l'index.
//
// Paramètres :
//   - idx : index de dépendances
//   - nodeID : identifiant du nœud
//   - factType : type de fait
//   - joinCondition : condition de jointure
//
// Retourne une erreur si l'extraction échoue.
func (ib *IndexBuilder) BuildFromBetaNode(
	idx *DependencyIndex,
	nodeID, factType string,
	joinCondition interface{},
) error {
	fields, err := ib.betaExtractor.ExtractFields(joinCondition)
	if err != nil {
		if ib.enableDiagnostics {
			ib.diagnostics.Errors = append(
				ib.diagnostics.Errors,
				fmt.Sprintf("Beta node %s: %v", nodeID, err),
			)
			ib.diagnostics.NodesSkipped++
		}
		return err
	}

	if len(fields) == 0 {
		if ib.enableDiagnostics {
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Beta node %s: no fields extracted", nodeID),
			)
		}
	}

	idx.AddBetaNode(nodeID, factType, fields)

	if ib.enableDiagnostics {
		ib.diagnostics.NodesProcessed++
		ib.diagnostics.FieldsExtracted += len(fields)
	}

	return nil
}

// BuildFromTerminalNode traite un nœud terminal et l'ajoute à l'index.
//
// Paramètres :
//   - idx : index de dépendances
//   - nodeID : identifiant du nœud
//   - factType : type de fait
//   - actions : liste des actions (AST)
//
// Retourne une erreur si l'extraction échoue.
func (ib *IndexBuilder) BuildFromTerminalNode(
	idx *DependencyIndex,
	nodeID, factType string,
	actions []interface{},
) error {
	allFields := make(map[string]bool)

	// Extraire champs depuis toutes les actions
	for _, action := range actions {
		fields, err := ib.actionExtractor.ExtractFields(action)
		if err != nil {
			if ib.enableDiagnostics {
				ib.diagnostics.Errors = append(
					ib.diagnostics.Errors,
					fmt.Sprintf("Terminal node %s action: %v", nodeID, err),
				)
			}
			// Continue avec les autres actions
			continue
		}

		for _, field := range fields {
			allFields[field] = true
		}
	}

	// Convertir map en slice
	fieldSlice := make([]string, 0, len(allFields))
	for field := range allFields {
		fieldSlice = append(fieldSlice, field)
	}

	if len(fieldSlice) == 0 {
		if ib.enableDiagnostics {
			ib.diagnostics.Warnings = append(
				ib.diagnostics.Warnings,
				fmt.Sprintf("Terminal node %s: no fields extracted", nodeID),
			)
		}
	}

	idx.AddTerminalNode(nodeID, factType, fieldSlice)

	if ib.enableDiagnostics {
		ib.diagnostics.NodesProcessed++
		ib.diagnostics.FieldsExtracted += len(fieldSlice)
	}

	return nil
}
