// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// ============================================================================
// Alpha Node Creation Orchestration
// ============================================================================

// alphaNodeCreationContext holds the state for alpha node creation orchestration
type alphaNodeCreationContext struct {
	network      *ReteNetwork
	ruleID       string
	condition    interface{}
	variableName string
	variableType string
	action       *Action
	storage      Storage
	pipeline     *ConstraintPipeline

	// Processed state
	actualCondition  interface{}
	exprType         ExpressionType
	normalizedExpr   interface{}
	conditions       []SimpleCondition
	opType           string
	parentNode       Node
	normalizedConds  []SimpleCondition
	fallbackToSimple bool
	shouldDecompose  bool
	chain            *AlphaChain
	terminalNode     *TerminalNode
}

// newAlphaNodeCreationContext creates a new context for alpha node creation
func newAlphaNodeCreationContext(
	cp *ConstraintPipeline,
	network *ReteNetwork,
	ruleID string,
	condition interface{},
	variableName string,
	variableType string,
	action *Action,
	storage Storage,
) *alphaNodeCreationContext {
	return &alphaNodeCreationContext{
		network:      network,
		ruleID:       ruleID,
		condition:    condition,
		variableName: variableName,
		variableType: variableType,
		action:       action,
		storage:      storage,
		pipeline:     cp,
	}
}

// unwrapCondition unwraps the condition if it's wrapped in a map
func (ctx *alphaNodeCreationContext) unwrapCondition() {
	ctx.actualCondition = ctx.condition

	condMap, ok := ctx.condition.(map[string]interface{})
	if !ok {
		return
	}

	condType, hasType := condMap["type"]
	if !hasType {
		return
	}

	// Handle wrapped constraint
	if condType == "constraint" {
		if constraint, hasConstraint := condMap["constraint"]; hasConstraint {
			ctx.actualCondition = constraint
		}
		return
	}

	// Handle special types that should use simple behavior
	if condType == "negation" || condType == "simple" || condType == "passthrough" {
		ctx.fallbackToSimple = true
		return
	}
}

// analyzeExpressionType analyzes the expression to determine its type
func (ctx *alphaNodeCreationContext) analyzeExpressionType() error {
	exprType, err := AnalyzeExpression(ctx.actualCondition)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur analyse expression: %v, fallback vers comportement simple", err)
		ctx.fallbackToSimple = true
		return nil // Not a fatal error, just fallback
	}

	ctx.exprType = exprType
	return nil
}

// handleORAndMixedExpressions handles OR and mixed (AND+OR) expressions
func (ctx *alphaNodeCreationContext) handleORAndMixedExpressions() error {
	if ctx.exprType != ExprTypeOR && ctx.exprType != ExprTypeMixed {
		return nil // Not OR/Mixed, continue
	}

	if ctx.exprType == ExprTypeOR {
		ctx.pipeline.GetLogger().Debug("   ℹ️  Expression OR détectée, normalisation avancée et création d'un nœud alpha unique")
	} else {
		ctx.pipeline.GetLogger().Debug("   ℹ️  Expression mixte (AND+OR) détectée, normalisation avancée et création d'un nœud alpha unique")
	}

	// Analyze nested OR complexity
	analysis, err := AnalyzeNestedOR(ctx.actualCondition)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur analyse OR imbriqué: %v, fallback vers normalisation simple", err)
		return ctx.performSimpleORNormalization()
	}

	// Display analysis information
	ctx.pipeline.GetLogger().Debug("   📊 Analyse OR: Complexité=%v, Profondeur=%d, OR=%d, AND=%d",
		analysis.Complexity, analysis.NestingDepth, analysis.ORTermCount, analysis.ANDTermCount)

	if analysis.OptimizationHint != "" {
		ctx.pipeline.GetLogger().Info("   💡 Suggestion: %s", analysis.OptimizationHint)
	}

	// Use advanced normalization for complex expressions
	if analysis.RequiresFlattening || analysis.RequiresDNF {
		return ctx.performAdvancedORNormalization(analysis)
	}

	// Use standard normalization for simple expressions
	return ctx.performStandardORNormalization()
}

// performSimpleORNormalization performs simple OR normalization as fallback
func (ctx *alphaNodeCreationContext) performSimpleORNormalization() error {
	normalizedExpr, err := NormalizeORExpression(ctx.actualCondition)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur normalisation simple: %v, fallback vers comportement simple", err)
		ctx.fallbackToSimple = true
		return nil
	}

	ctx.normalizedExpr = normalizedExpr
	normalizedCondition := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalizedExpr,
	}

	return ctx.pipeline.createSimpleAlphaNodeWithTerminal(
		ctx.network, ctx.ruleID, normalizedCondition,
		ctx.variableName, ctx.variableType, ctx.action, ctx.storage,
	)
}

// performAdvancedORNormalization performs advanced OR normalization with flattening/DNF
func (ctx *alphaNodeCreationContext) performAdvancedORNormalization(analysis *NestedORAnalysis) error {
	ctx.pipeline.GetLogger().Info("   🔧 Application de la normalisation avancée (aplatissement=%v, DNF=%v)",
		analysis.RequiresFlattening, analysis.RequiresDNF)

	normalizedExpr, err := NormalizeNestedOR(ctx.actualCondition)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur normalisation avancée: %v, fallback vers normalisation simple", err)
		return ctx.performSimpleORNormalization()
	}

	ctx.pipeline.GetLogger().Info("   ✅ Normalisation avancée réussie")
	ctx.normalizedExpr = normalizedExpr

	normalizedCondition := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalizedExpr,
	}

	return ctx.pipeline.createSimpleAlphaNodeWithTerminal(
		ctx.network, ctx.ruleID, normalizedCondition,
		ctx.variableName, ctx.variableType, ctx.action, ctx.storage,
	)
}

// performStandardORNormalization performs standard OR normalization
func (ctx *alphaNodeCreationContext) performStandardORNormalization() error {
	ctx.pipeline.GetLogger().Info("   🔧 Application de la normalisation standard")

	normalizedExpr, err := NormalizeORExpression(ctx.actualCondition)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur normalisation: %v, utilisation expression originale", err)
		normalizedExpr = ctx.actualCondition
	}

	ctx.normalizedExpr = normalizedExpr

	normalizedCondition := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalizedExpr,
	}

	return ctx.pipeline.createSimpleAlphaNodeWithTerminal(
		ctx.network, ctx.ruleID, normalizedCondition,
		ctx.variableName, ctx.variableType, ctx.action, ctx.storage,
	)
}

// checkDecomposability checks if the expression can be decomposed
func (ctx *alphaNodeCreationContext) checkDecomposability() {
	if !CanDecompose(ctx.exprType) {
		ctx.pipeline.GetLogger().Debug("   ℹ️  Expression de type %s non décomposable, utilisation du nœud simple", ctx.exprType)
		ctx.fallbackToSimple = true
		return
	}

	// Simple or arithmetic expressions use simple behavior
	if ctx.exprType == ExprTypeSimple || ctx.exprType == ExprTypeArithmetic {
		ctx.fallbackToSimple = true
		return
	}

	// AND or NOT expressions can be decomposed
	ctx.pipeline.GetLogger().Debug("   🔍 Expression de type %s détectée, tentative de décomposition...", ctx.exprType)
	ctx.shouldDecompose = true
}

// extractAndNormalizeConditions extracts and normalizes conditions from the expression
func (ctx *alphaNodeCreationContext) extractAndNormalizeConditions() error {
	conditions, opType, err := ExtractConditions(ctx.actualCondition)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur extraction conditions: %v, fallback vers comportement simple", err)
		ctx.fallbackToSimple = true
		return nil
	}

	// Single condition doesn't need chain
	if len(conditions) <= 1 {
		ctx.pipeline.GetLogger().Debug("   ℹ️  Une seule condition extraite, utilisation du nœud simple")
		ctx.fallbackToSimple = true
		return nil
	}

	ctx.conditions = conditions
	ctx.opType = opType

	ctx.pipeline.GetLogger().Info("   🔗 Décomposition en chaîne: %d conditions détectées (opérateur: %s)", len(conditions), opType)

	// Normalize conditions
	normalizedConds := NormalizeConditions(conditions, opType)
	ctx.normalizedConds = normalizedConds
	ctx.pipeline.GetLogger().Info("   📋 Conditions normalisées: %d condition(s)", len(ctx.normalizedConds))

	return nil
}

// resolveParentNode finds the TypeNode parent to connect the chain
func (ctx *alphaNodeCreationContext) resolveParentNode() error {
	// Try to find TypeNode by variable type
	if ctx.variableType != "" {
		if typeNode, exists := ctx.network.TypeNodes[ctx.variableType]; exists {
			ctx.parentNode = typeNode
			return nil
		}
	}

	// Fallback: use first available TypeNode
	for _, typeNode := range ctx.network.TypeNodes {
		ctx.parentNode = typeNode
		return nil
	}

	ctx.pipeline.GetLogger().Warn("   ⚠️  Aucun TypeNode trouvé, fallback vers comportement simple")
	ctx.fallbackToSimple = true
	return nil
}

// buildAndValidateChain builds and validates the alpha node chain
func (ctx *alphaNodeCreationContext) buildAndValidateChain() error {
	chainBuilder := NewAlphaChainBuilder(ctx.network, ctx.storage)

	chain, err := chainBuilder.BuildChain(ctx.normalizedConds, ctx.variableName, ctx.parentNode, ctx.ruleID)
	if err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Erreur construction chaîne: %v, fallback vers comportement simple", err)
		ctx.fallbackToSimple = true
		return nil
	}

	// Validate chain
	if err := chain.ValidateChain(); err != nil {
		ctx.pipeline.GetLogger().Warn("   ⚠️  Chaîne invalide: %v, fallback vers comportement simple", err)
		ctx.fallbackToSimple = true
		return nil
	}

	ctx.chain = chain

	// Get and display chain statistics
	stats := chainBuilder.GetChainStats(chain)
	sharedCount := 0
	if sc, ok := stats["shared_nodes"].(int); ok {
		sharedCount = sc
	}

	ctx.pipeline.GetLogger().Info("   ✅ Chaîne construite: %d nœud(s), %d partagé(s)", len(chain.Nodes), sharedCount)

	// Log details of each node
	for i, node := range chain.Nodes {
		if i < sharedCount {
			ctx.pipeline.GetLogger().Info("   ♻️  AlphaNode partagé réutilisé: %s (hash: %s)", node.ID, chain.Hashes[i])
		} else {
			ctx.pipeline.GetLogger().Info("   ✨ Nouveau AlphaNode créé: %s (hash: %s)", node.ID, chain.Hashes[i])
		}
	}

	return nil
}

// createAndAttachTerminal creates and attaches the terminal node to the chain
func (ctx *alphaNodeCreationContext) createAndAttachTerminal() error {
	ctx.terminalNode = NewTerminalNode(ctx.ruleID+"_terminal", ctx.action, ctx.storage)
	ctx.terminalNode.SetNetwork(ctx.network)
	ctx.chain.FinalNode.AddChild(ctx.terminalNode)
	ctx.network.TerminalNodes[ctx.terminalNode.ID] = ctx.terminalNode

	// Configure observer if network has one
	if ctx.network.actionObserver != nil {
		ctx.terminalNode.SetObserver(ctx.network.actionObserver)
	}

	// Register terminal node with lifecycle manager
	if ctx.network.LifecycleManager != nil {
		ctx.network.LifecycleManager.RegisterNode(ctx.terminalNode.ID, "terminal")
		ctx.network.LifecycleManager.AddRuleToNode(ctx.terminalNode.ID, ctx.ruleID, ctx.ruleID)
	}

	ctx.pipeline.GetLogger().Debug("   ✓ TerminalNode %s attaché au nœud final %s de la chaîne", ctx.terminalNode.ID, ctx.chain.FinalNode.ID)

	return nil
}

// createAlphaNodeWithTerminalOrchestrated orchestrates the creation of an alpha node with terminal
// using the extract method pattern to separate concerns
func (cp *ConstraintPipeline) createAlphaNodeWithTerminalOrchestrated(
	network *ReteNetwork,
	ruleID string,
	condition interface{},
	variableName string,
	variableType string,
	action *Action,
	storage Storage,
) error {
	ctx := newAlphaNodeCreationContext(cp, network, ruleID, condition, variableName, variableType, action, storage)

	// Step 1: Unwrap condition
	ctx.unwrapCondition()
	if ctx.fallbackToSimple {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Step 2: Analyze expression type
	if err := ctx.analyzeExpressionType(); err != nil {
		return err
	}
	if ctx.fallbackToSimple {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Step 3: Handle OR/Mixed expressions (MUST be before decomposition check)
	if err := ctx.handleORAndMixedExpressions(); err != nil {
		return err
	}
	// OR/Mixed handling returns directly if applicable
	if ctx.exprType == ExprTypeOR || ctx.exprType == ExprTypeMixed {
		return nil // Already handled
	}

	// Step 4: Check decomposability
	ctx.checkDecomposability()
	if ctx.fallbackToSimple {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Step 5: Extract and normalize conditions
	if err := ctx.extractAndNormalizeConditions(); err != nil {
		return err
	}
	if ctx.fallbackToSimple {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Step 6: Resolve parent node
	if err := ctx.resolveParentNode(); err != nil {
		return err
	}
	if ctx.fallbackToSimple {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Step 7: Build and validate chain
	if err := ctx.buildAndValidateChain(); err != nil {
		return err
	}
	if ctx.fallbackToSimple {
		return cp.createSimpleAlphaNodeWithTerminal(network, ruleID, condition, variableName, variableType, action, storage)
	}

	// Step 8: Create and attach terminal
	return ctx.createAndAttachTerminal()
}
