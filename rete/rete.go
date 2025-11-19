package rete

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// ========== TYPES DE BASE ==========

// Fact représente un fait dans le système RETE
type Fact struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp time.Time              `json:"timestamp"`
}

// String retourne la représentation string d'un fait
func (f *Fact) String() string {
	return fmt.Sprintf("Fact{ID:%s, Type:%s, Fields:%v}", f.ID, f.Type, f.Fields)
}

// GetField retourne la valeur d'un champ
func (f *Fact) GetField(fieldName string) (interface{}, bool) {
	value, exists := f.Fields[fieldName]
	return value, exists
}

// Token représente un token dans le réseau RETE
type Token struct {
	ID           string           `json:"id"`
	Facts        []*Fact          `json:"facts"`
	NodeID       string           `json:"node_id"`
	Parent       *Token           `json:"parent,omitempty"`
	Bindings     map[string]*Fact `json:"bindings,omitempty"`       // Nouveau: bindings pour jointures
	IsJoinResult bool             `json:"is_join_result,omitempty"` // Indique si c'est un token de jointure réussie
}

// WorkingMemory représente la mémoire de travail d'un nœud
type WorkingMemory struct {
	NodeID string            `json:"node_id"`
	Facts  map[string]*Fact  `json:"facts"`
	Tokens map[string]*Token `json:"tokens"`
}

// AddFact ajoute un fait à la mémoire
func (wm *WorkingMemory) AddFact(fact *Fact) {
	if wm.Facts == nil {
		wm.Facts = make(map[string]*Fact)
	}
	wm.Facts[fact.ID] = fact
}

// RemoveFact supprime un fait de la mémoire
func (wm *WorkingMemory) RemoveFact(factID string) {
	delete(wm.Facts, factID)
}

// GetFacts retourne tous les faits de la mémoire
func (wm *WorkingMemory) GetFacts() []*Fact {
	facts := make([]*Fact, 0, len(wm.Facts))
	for _, fact := range wm.Facts {
		facts = append(facts, fact)
	}
	return facts
}

// AddToken ajoute un token à la mémoire
func (wm *WorkingMemory) AddToken(token *Token) {
	if wm.Tokens == nil {
		wm.Tokens = make(map[string]*Token)
	}
	wm.Tokens[token.ID] = token
}

// RemoveToken supprime un token de la mémoire
func (wm *WorkingMemory) RemoveToken(tokenID string) {
	delete(wm.Tokens, tokenID)
}

// GetTokens retourne tous les tokens de la mémoire
func (wm *WorkingMemory) GetTokens() []*Token {
	tokens := make([]*Token, 0, len(wm.Tokens))
	for _, token := range wm.Tokens {
		tokens = append(tokens, token)
	}
	return tokens
}

// GetFactsByVariable retourne les faits associés aux variables spécifiées
func (wm *WorkingMemory) GetFactsByVariable(variables []string) []*Fact {
	// Pour l'instant, retourne tous les faits (implémentation simplifiée)
	return wm.GetFacts()
}

// GetTokensByVariable retourne les tokens associés aux variables spécifiées
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token {
	// Pour l'instant, retourne tous les tokens (implémentation simplifiée)
	return wm.GetTokens()
}

// ========== INTERFACES ==========

// Node interface pour tous les nœuds du réseau RETE
type Node interface {
	GetID() string
	GetType() string
	GetMemory() *WorkingMemory
	ActivateLeft(token *Token) error
	ActivateRight(fact *Fact) error
	AddChild(child Node)
	GetChildren() []Node
}

// Storage interface pour la persistance
type Storage interface {
	SaveMemory(nodeID string, memory *WorkingMemory) error
	LoadMemory(nodeID string) (*WorkingMemory, error)
	DeleteMemory(nodeID string) error
	ListNodes() ([]string, error)
}

// ========== TYPES POUR COMPATIBILITÉ AST ==========

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TypeDefinition struct {
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type JobCall struct {
	Type string        `json:"type"`
	Name string        `json:"name"`
	Args []interface{} `json:"args"`
}

type Action struct {
	Type string  `json:"type"`
	Job  JobCall `json:"job"`
}

type TypedVariable struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

type Set struct {
	Type      string          `json:"type"`
	Variables []TypedVariable `json:"variables"`
}

type Expression struct {
	Type        string      `json:"type"`
	Set         Set         `json:"set"`
	Constraints interface{} `json:"constraints"`
	Action      *Action     `json:"action,omitempty"`
}

type Program struct {
	Types       []TypeDefinition `json:"types"`
	Expressions []Expression     `json:"expressions"`
}

// ========== IMPLÉMENTATION DES NŒUDS ==========

// BaseNode implémente les fonctionnalités communes à tous les nœuds
type BaseNode struct {
	ID       string         `json:"id"`
	Type     string         `json:"type"`
	Memory   *WorkingMemory `json:"memory"`
	Children []Node         `json:"children"`
	Storage  Storage        `json:"-"`
	mutex    sync.RWMutex   `json:"-"`
}

// GetID retourne l'ID du nœud
func (bn *BaseNode) GetID() string {
	return bn.ID
}

// GetType retourne le type du nœud
func (bn *BaseNode) GetType() string {
	return bn.Type
}

// GetMemory retourne la mémoire de travail du nœud
func (bn *BaseNode) GetMemory() *WorkingMemory {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return bn.Memory
}

// AddChild ajoute un nœud enfant
func (bn *BaseNode) AddChild(child Node) {
	bn.mutex.Lock()
	defer bn.mutex.Unlock()
	bn.Children = append(bn.Children, child)
}

// GetChildren retourne les nœuds enfants
func (bn *BaseNode) GetChildren() []Node {
	bn.mutex.RLock()
	defer bn.mutex.RUnlock()
	return bn.Children
}

// PropagateToChildren propage un fait ou token aux enfants
func (bn *BaseNode) PropagateToChildren(fact *Fact, token *Token) error {
	for _, child := range bn.GetChildren() {
		if fact != nil {
			if err := child.ActivateRight(fact); err != nil {
				return fmt.Errorf("erreur propagation fait vers %s: %w", child.GetID(), err)
			}
		}
		if token != nil {
			if err := child.ActivateLeft(token); err != nil {
				return fmt.Errorf("erreur propagation token vers %s: %w", child.GetID(), err)
			}
		}
	}
	return nil
}

// SaveMemory sauvegarde la mémoire du nœud
func (bn *BaseNode) SaveMemory() error {
	if bn.Storage != nil {
		return bn.Storage.SaveMemory(bn.ID, bn.Memory)
	}
	return nil
}

// ========== NŒUD RACINE ==========

// RootNode est le nœud racine qui reçoit tous les faits
type RootNode struct {
	BaseNode
}

// NewRootNode crée un nouveau nœud racine
func NewRootNode(storage Storage) *RootNode {
	return &RootNode{
		BaseNode: BaseNode{
			ID:       "root",
			Type:     "root",
			Memory:   &WorkingMemory{NodeID: "root", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
	}
}

// ActivateLeft (non utilisé pour le nœud racine)
func (rn *RootNode) ActivateLeft(token *Token) error {
	return fmt.Errorf("le nœud racine ne peut pas recevoir de tokens")
}

// ActivateRight distribue les faits aux nœuds de type
func (rn *RootNode) ActivateRight(fact *Fact) error {
	rn.mutex.Lock()
	rn.Memory.AddFact(fact)
	rn.mutex.Unlock()

	// Log désactivé pour les performances
	// fmt.Printf("[ROOT] Reçu fait: %s\n", fact.String())

	// Persistance désactivée pour les performances

	// Propager aux enfants (TypeNodes)
	return rn.PropagateToChildren(fact, nil)
}

// ========== NŒUD DE TYPE ==========

// TypeNode filtre les faits selon leur type
type TypeNode struct {
	BaseNode
	TypeName       string         `json:"type_name"`
	TypeDefinition TypeDefinition `json:"type_definition"`
}

// NewTypeNode crée un nouveau nœud de type
func NewTypeNode(typeName string, typeDef TypeDefinition, storage Storage) *TypeNode {
	return &TypeNode{
		BaseNode: BaseNode{
			ID:       fmt.Sprintf("type_%s", typeName),
			Type:     "type",
			Memory:   &WorkingMemory{NodeID: fmt.Sprintf("type_%s", typeName), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		TypeName:       typeName,
		TypeDefinition: typeDef,
	}
}

// ActivateLeft (non utilisé pour les nœuds de type)
func (tn *TypeNode) ActivateLeft(token *Token) error {
	return fmt.Errorf("les nœuds de type ne reçoivent pas de tokens")
}

// ActivateRight filtre les faits par type et les propage
func (tn *TypeNode) ActivateRight(fact *Fact) error {
	// Vérifier si le fait correspond au type de ce nœud
	if fact.Type != tn.TypeName {
		return nil // Ignorer silencieusement les faits d'autres types
	}

	// Log désactivé pour les performances
	// fmt.Printf("[TYPE_%s] Reçu fait: %s\n", tn.TypeName, fact.String())

	// Valider les champs du fait
	if err := tn.validateFact(fact); err != nil {
		return fmt.Errorf("validation du fait échouée: %w", err)
	}

	tn.mutex.Lock()
	tn.Memory.AddFact(fact)
	tn.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Propager aux enfants (AlphaNodes)
	return tn.PropagateToChildren(fact, nil)
}

// validateFact valide qu'un fait respecte la définition de type
func (tn *TypeNode) validateFact(fact *Fact) error {
	for _, field := range tn.TypeDefinition.Fields {
		value, exists := fact.Fields[field.Name]
		if !exists {
			return fmt.Errorf("champ manquant: %s", field.Name)
		}

		// Validation basique des types
		if !tn.isValidType(value, field.Type) {
			return fmt.Errorf("type invalide pour le champ %s: attendu %s", field.Name, field.Type)
		}
	}
	return nil
}

// isValidType vérifie si une valeur correspond au type attendu
func (tn *TypeNode) isValidType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		switch value.(type) {
		case int, int32, int64, float32, float64:
			return true
		}
		return false
	case "bool":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

// ========== NŒUD ALPHA (CONDITIONS SIMPLES) ==========

// AlphaNode teste une condition sur un fait
type AlphaNode struct {
	BaseNode
	Condition    interface{} `json:"condition"`
	VariableName string      `json:"variable_name"`
}

// NewAlphaNode crée un nouveau nœud alpha
func NewAlphaNode(nodeID string, condition interface{}, variableName string, storage Storage) *AlphaNode {
	return &AlphaNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "alpha",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Condition:    condition,
		VariableName: variableName,
	}
}

// ActivateLeft (non utilisé pour les nœuds alpha)
func (an *AlphaNode) ActivateLeft(token *Token) error {
	return fmt.Errorf("les nœuds alpha ne reçoivent pas de tokens")
}

// ActivateRight teste la condition sur le fait
func (an *AlphaNode) ActivateRight(fact *Fact) error {
	// Log désactivé pour les performances
	// fmt.Printf("[ALPHA_%s] Test condition sur fait: %s\n", an.ID, fact.String())

	// Cas spécial: passthrough pour les JoinNodes - pas de filtrage
	if an.Condition != nil {
		if condMap, ok := an.Condition.(map[string]interface{}); ok {
			if condType, exists := condMap["type"].(string); exists && condType == "passthrough" {
				// Mode pass-through: convertir le fait en token et propager selon le côté
				an.mutex.Lock()
				an.Memory.AddFact(fact)
				an.mutex.Unlock()

				// Créer un token pour le fait avec la variable correspondante
				token := &Token{
					ID:       fmt.Sprintf("alpha_token_%s_%s", an.ID, fact.ID),
					Facts:    []*Fact{fact},
					NodeID:   an.ID,
					Bindings: map[string]*Fact{an.VariableName: fact},
				}

				// Déterminer le côté et propager selon l'architecture RETE
				side, sideExists := condMap["side"].(string)
				if sideExists && side == "left" {
					fmt.Printf("🔗 ALPHA PASSTHROUGH[%s]: Propagation LEFT pour JoinNode\n", an.ID)
					return an.PropagateToChildren(nil, token) // ActivateLeft
				} else {
					fmt.Printf("🔗 ALPHA PASSTHROUGH[%s]: Propagation RIGHT pour JoinNode\n", an.ID)
					return an.PropagateToChildren(fact, nil) // ActivateRight
				}
			}
		}
	}

	// Évaluation normale de condition Alpha
	if an.Condition != nil {
		evaluator := NewAlphaConditionEvaluator()
		passed, err := evaluator.EvaluateCondition(an.Condition, fact, an.VariableName)
		if err != nil {
			return fmt.Errorf("erreur évaluation condition Alpha: %w", err)
		}

		// Si la condition n'est pas satisfaite, ignorer le fait
		if !passed {
			// Log désactivé pour les performances
			// fmt.Printf("[ALPHA_%s] Condition non satisfaite pour le fait: %s\n", an.ID, fact.String())
			return nil
		}
	}

	// Log désactivé pour les performances
	// fmt.Printf("[ALPHA_%s] Condition satisfaite pour le fait: %s\n", an.ID, fact.String())

	an.mutex.Lock()
	an.Memory.AddFact(fact)
	an.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Créer un token et le propager
	token := &Token{
		ID:     fmt.Sprintf("token_%s_%s", an.ID, fact.ID),
		Facts:  []*Fact{fact},
		NodeID: an.ID,
	}

	return an.PropagateToChildren(nil, token)
}

// ========== NŒUD TERMINAL (ACTIONS) ==========

// TerminalNode déclenche une action
type TerminalNode struct {
	BaseNode
	Action *Action `json:"action"`
}

// NewTerminalNode crée un nouveau nœud terminal
func NewTerminalNode(nodeID string, action *Action, storage Storage) *TerminalNode {
	return &TerminalNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "terminal",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0), // Les nœuds terminaux n'ont pas d'enfants
			Storage:  storage,
		},
		Action: action,
	}
}

// ActivateLeft déclenche l'action
func (tn *TerminalNode) ActivateLeft(token *Token) error {
	// Log désactivé pour les performances
	// fmt.Printf("[TERMINAL_%s] Déclenchement action avec token: %s\n", tn.ID, token.ID)

	// Stocker le token
	tn.mutex.Lock()
	if tn.Memory.Tokens == nil {
		tn.Memory.Tokens = make(map[string]*Token)
	}
	tn.Memory.Tokens[token.ID] = token
	tn.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Déclencher l'action
	return tn.executeAction(token)
}

// ActivateRight (non utilisé pour les nœuds terminaux)
func (tn *TerminalNode) ActivateRight(fact *Fact) error {
	return fmt.Errorf("les nœuds terminaux ne reçoivent pas de faits directement")
}

// executeAction affiche l'action déclenchée avec les faits déclencheurs (version tuple-space)
func (tn *TerminalNode) executeAction(token *Token) error {
	// Les actions sont maintenant obligatoires dans la grammaire
	// Mais nous gardons cette vérification par sécurité
	if tn.Action == nil {
		return fmt.Errorf("aucune action définie pour le nœud %s", tn.ID)
	}

	// === VERSION TUPLE-SPACE ===
	// Au lieu d'exécuter l'action, on l'affiche avec les faits déclencheurs
	// Les agents du tuple-space viendront "prendre" ces tuples plus tard

	actionName := tn.Action.Job.Name
	fmt.Printf("🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: %s", actionName)

	// Afficher les faits déclencheurs entre parenthèses
	if len(token.Facts) > 0 {
		fmt.Print(" (")
		for i, fact := range token.Facts {
			if i > 0 {
				fmt.Print(", ")
			}
			// Format compact : Type(id:value, field:value, ...)
			fmt.Printf("%s(", fact.Type)
			fieldCount := 0
			for key, value := range fact.Fields {
				if fieldCount > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s:%v", key, value)
				fieldCount++
			}
			fmt.Print(")")
		}
		fmt.Print(")")
	}
	fmt.Println()

	return nil
}

// ========== NŒUD DE JOINTURE (BETA) ==========

// JoinNode effectue des jointures entre tuples basées sur des conditions d'égalité
type JoinNode struct {
	BaseNode
	Condition      map[string]interface{} `json:"condition"`
	LeftVariables  []string               `json:"left_variables"`
	RightVariables []string               `json:"right_variables"`
	AllVariables   []string               `json:"all_variables"`
	VariableTypes  map[string]string      `json:"variable_types"` // Nouveau: mapping variable -> type
	JoinConditions []JoinCondition        `json:"join_conditions"`
	mutex          sync.RWMutex
	// Mémoires séparées pour architecture RETE propre
	LeftMemory   *WorkingMemory // Tokens venant de la gauche
	RightMemory  *WorkingMemory // Tokens venant de la droite
	ResultMemory *WorkingMemory // Tokens de jointure réussie
}

// JoinCondition représente une condition de jointure entre variables
type JoinCondition struct {
	LeftField  string `json:"left_field"`  // p.id
	RightField string `json:"right_field"` // o.customer_id
	LeftVar    string `json:"left_var"`    // p
	RightVar   string `json:"right_var"`   // o
	Operator   string `json:"operator"`    // ==
}

// NewJoinNode crée un nouveau nœud de jointure
func NewJoinNode(nodeID string, condition map[string]interface{}, leftVars []string, rightVars []string, varTypes map[string]string, storage Storage) *JoinNode {
	allVars := append(leftVars, rightVars...)

	return &JoinNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "join",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Condition:      condition,
		LeftVariables:  leftVars,
		RightVariables: rightVars,
		AllVariables:   allVars,
		VariableTypes:  varTypes,
		JoinConditions: extractJoinConditions(condition),
		// Initialiser les mémoires séparées
		LeftMemory:   &WorkingMemory{NodeID: nodeID + "_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:  &WorkingMemory{NodeID: nodeID + "_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory: &WorkingMemory{NodeID: nodeID + "_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
	}
}

// ActivateLeft traite les tokens de la gauche (généralement des AlphaNodes)
func (jn *JoinNode) ActivateLeft(token *Token) error {
	fmt.Printf("🔍 JOINNODE[%s]: ActivateLeft - token %s\n", jn.ID, token.ID)

	// Stocker le token dans la mémoire gauche
	jn.mutex.Lock()
	jn.LeftMemory.AddToken(token)
	jn.mutex.Unlock()

	fmt.Printf("🔍 JOINNODE[%s]: Mémoire gauche: %d tokens\n", jn.ID, len(jn.LeftMemory.GetTokens()))

	// Essayer de joindre avec tous les tokens de la mémoire droite
	rightTokens := jn.RightMemory.GetTokens()
	fmt.Printf("🔍 JOINNODE[%s]: Mémoire droite: %d tokens\n", jn.ID, len(rightTokens))

	for _, rightToken := range rightTokens {
		fmt.Printf("🔍 JOINNODE[%s]: Tentative jointure LEFT[%s] + RIGHT[%s]\n", jn.ID, token.ID, rightToken.ID)
		if joinedToken := jn.performJoinWithTokens(token, rightToken); joinedToken != nil {
			fmt.Printf("🔍 JOINNODE[%s]: Jointure réussie! LEFT[%s] + RIGHT[%s]\n", jn.ID, token.ID, rightToken.ID)

			// Stocker uniquement les tokens de jointure réussie
			joinedToken.IsJoinResult = true
			jn.mutex.Lock()
			jn.ResultMemory.AddToken(joinedToken)
			jn.Memory.AddToken(joinedToken) // Pour compatibilité avec le comptage
			jn.mutex.Unlock()

			if err := jn.PropagateToChildren(nil, joinedToken); err != nil {
				return err
			}
		}
	}
	return nil
}

// ActivateRight traite les faits de la droite (nouveau fait injecté via AlphaNode)
func (jn *JoinNode) ActivateRight(fact *Fact) error {
	fmt.Printf("🔍 JOINNODE[%s]: ActivateRight - %s\n", jn.ID, fact.Type)

	// Convertir le fait en token pour la mémoire droite
	factVar := jn.getVariableForFact(fact)
	if factVar == "" {
		fmt.Printf("🔍 JOINNODE[%s]: Fait %s non applicable (variable introuvable)\n", jn.ID, fact.ID)
		return nil // Fait non applicable à ce JoinNode
	}

	factToken := &Token{
		ID:       fmt.Sprintf("right_token_%s_%s", jn.ID, fact.ID),
		Facts:    []*Fact{fact},
		NodeID:   jn.ID,
		Bindings: map[string]*Fact{factVar: fact},
	}

	// Stocker le token dans la mémoire droite
	jn.mutex.Lock()
	jn.RightMemory.AddToken(factToken)
	jn.mutex.Unlock()

	fmt.Printf("🔍 JOINNODE[%s]: Mémoire droite: %d tokens\n", jn.ID, len(jn.RightMemory.GetTokens()))

	// Essayer de joindre avec tous les tokens de la mémoire gauche
	leftTokens := jn.LeftMemory.GetTokens()
	fmt.Printf("🔍 JOINNODE[%s]: Mémoire gauche: %d tokens\n", jn.ID, len(leftTokens))

	for _, leftToken := range leftTokens {
		fmt.Printf("🔍 JOINNODE[%s]: Tentative jointure LEFT[%s] + RIGHT[%s]\n", jn.ID, leftToken.ID, factToken.ID)
		if joinedToken := jn.performJoinWithTokens(leftToken, factToken); joinedToken != nil {
			fmt.Printf("🔍 JOINNODE[%s]: Jointure réussie! LEFT[%s] + RIGHT[%s]\n", jn.ID, leftToken.ID, factToken.ID)

			// Stocker uniquement les tokens de jointure réussie
			joinedToken.IsJoinResult = true
			jn.mutex.Lock()
			jn.ResultMemory.AddToken(joinedToken)
			jn.Memory.AddToken(joinedToken) // Pour compatibilité avec le comptage
			jn.mutex.Unlock()

			if err := jn.PropagateToChildren(nil, joinedToken); err != nil {
				return err
			}
		}
	}
	return nil
}

// performJoinWithTokens effectue la jointure entre deux tokens
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
	// Vérifier que les tokens ont des variables différentes
	if !jn.tokensHaveDifferentVariables(token1, token2) {
		return nil
	}

	// Combiner les bindings des deux tokens
	combinedBindings := make(map[string]*Fact)

	// Copier les bindings du premier token
	for varName, varFact := range token1.Bindings {
		combinedBindings[varName] = varFact
	}

	// Copier les bindings du second token
	for varName, varFact := range token2.Bindings {
		combinedBindings[varName] = varFact
	}

	// Valider les conditions de jointure
	if !jn.evaluateJoinConditions(combinedBindings) {
		return nil // Jointure échoue
	}

	// Créer et retourner le token joint
	return &Token{
		ID:       fmt.Sprintf("%s_JOIN_%s", token1.ID, token2.ID),
		Bindings: combinedBindings,
		NodeID:   jn.ID,
		Facts:    append(token1.Facts, token2.Facts...),
	}
}

// tokensHaveDifferentVariables vérifie que les tokens représentent des variables différentes
func (jn *JoinNode) tokensHaveDifferentVariables(token1 *Token, token2 *Token) bool {
	for var1 := range token1.Bindings {
		for var2 := range token2.Bindings {
			if var1 == var2 {
				return false // Même variable = pas de jointure possible
			}
		}
	}
	return true
}

// getVariableForFact détermine la variable associée à un fait basé sur son type
func (jn *JoinNode) getVariableForFact(fact *Fact) string {
	// Utiliser le mapping variable -> type du JoinNode
	for _, varName := range jn.AllVariables {
		if expectedType, exists := jn.VariableTypes[varName]; exists {
			if expectedType == fact.Type {
				fmt.Printf("🔍 JOINNODE[%s]: Variable %s trouvée pour fait %s (type: %s)\n", jn.ID, varName, fact.ID, fact.Type)
				return varName
			}
		}
	}

	fmt.Printf("❌ JOINNODE[%s]: Aucune variable trouvée pour fait %s (type: %s)\n", jn.ID, fact.ID, fact.Type)
	fmt.Printf("   Variables disponibles: %v\n", jn.AllVariables)
	fmt.Printf("   Types attendus: %v\n", jn.VariableTypes)
	return ""
}

// evaluateJoinConditions vérifie si toutes les conditions de jointure sont respectées
func (jn *JoinNode) evaluateJoinConditions(bindings map[string]*Fact) bool {
	fmt.Printf("🔍 JOINNODE[%s]: Évaluation conditions jointure\n", jn.ID)
	fmt.Printf("  📊 Bindings: %d variables\n", len(bindings))
	for varName, fact := range bindings {
		fmt.Printf("    %s -> %s (ID: %s)\n", varName, fact.Type, fact.ID)
	}
	fmt.Printf("  📊 Conditions: %d à vérifier\n", len(jn.JoinConditions))
	for i, condition := range jn.JoinConditions {
		fmt.Printf("    Condition %d: %s.%s %s %s.%s\n", i,
			condition.LeftVar, condition.LeftField, condition.Operator,
			condition.RightVar, condition.RightField)
	}

	// Vérifier qu'on a au moins 2 variables différentes
	if len(bindings) < 2 {
		fmt.Printf("  ❌ Pas assez de variables (%d < 2)\n", len(bindings))
		return false
	}

	// Évaluer chaque condition de jointure
	for i, joinCondition := range jn.JoinConditions {
		leftFact := bindings[joinCondition.LeftVar]
		rightFact := bindings[joinCondition.RightVar]

		if leftFact == nil || rightFact == nil {
			fmt.Printf("  ❌ Condition %d: variable manquante (%s ou %s)\n", i, joinCondition.LeftVar, joinCondition.RightVar)
			return false // Une variable manque
		}

		// Récupérer les valeurs des champs
		leftValue := leftFact.Fields[joinCondition.LeftField]
		rightValue := rightFact.Fields[joinCondition.RightField]

		fmt.Printf("  🔍 Condition %d: %v %s %v\n", i, leftValue, joinCondition.Operator, rightValue)

		// Évaluer l'opérateur
		switch joinCondition.Operator {
		case "==":
			if leftValue != rightValue {
				fmt.Printf("  ❌ Condition %d échoue: %v != %v\n", i, leftValue, rightValue)
				return false
			}
			fmt.Printf("  ✅ Condition %d réussie: %v == %v\n", i, leftValue, rightValue)
		case "!=":
			if leftValue == rightValue {
				fmt.Printf("  ❌ Condition %d échoue: %v == %v\n", i, leftValue, rightValue)
				return false
			}
			fmt.Printf("  ✅ Condition %d réussie: %v != %v\n", i, leftValue, rightValue)
		case "<":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat >= rightFloat {
						return false
					}
				} else {
					return false // Comparaison numérique impossible
				}
			} else {
				return false
			}
		case ">":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat <= rightFloat {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		case "<=":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat > rightFloat {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		case ">=":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat < rightFloat {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		default:
			return false // Opérateur non supporté
		}
	}

	return true // Toutes les conditions sont satisfaites
}

// convertToFloat64 tente de convertir une valeur en float64
func convertToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// extractJoinConditions extrait les conditions de jointure d'une condition complexe
func extractJoinConditions(condition map[string]interface{}) []JoinCondition {
	fmt.Printf("🔍 EXTRACT JOIN CONDITIONS: analyzing condition\n")
	fmt.Printf("  🔧 Condition type: %T\n", condition)
	for key, value := range condition {
		fmt.Printf("    %s: %v (type: %T)\n", key, value, value)
	}

	var joinConditions []JoinCondition

	// Cas 1: condition wrappée dans un type "constraint"
	if conditionType, exists := condition["type"].(string); exists && conditionType == "constraint" {
		fmt.Printf("  🔧 Condition wrappée détectée - extraction de la sous-condition\n")
		if innerCondition, ok := condition["constraint"].(map[string]interface{}); ok {
			fmt.Printf("  ✅ Sous-condition extraite, analyse récursive\n")
			return extractJoinConditions(innerCondition)
		}
	}

	// Cas 2: condition EXISTS avec array de conditions
	if conditionType, exists := condition["type"].(string); exists && conditionType == "exists" {
		fmt.Printf("  🔧 Condition EXISTS détectée - extraction des sous-conditions\n")
		if conditionsData, ok := condition["conditions"].([]map[string]interface{}); ok {
			fmt.Printf("  ✅ Array de conditions EXISTS trouvé: %d conditions\n", len(conditionsData))
			for i, subCondition := range conditionsData {
				fmt.Printf("  🔍 Analyse condition EXISTS %d: %+v\n", i, subCondition)
				subJoinConditions := extractJoinConditions(subCondition)
				joinConditions = append(joinConditions, subJoinConditions...)
			}
			return joinConditions
		}
	}

	// Cas 3: condition directe de comparaison
	if conditionType, exists := condition["type"].(string); exists && conditionType == "comparison" {
		fmt.Printf("  ✅ Condition de comparaison détectée\n")
		if left, leftOk := condition["left"].(map[string]interface{}); leftOk {
			if right, rightOk := condition["right"].(map[string]interface{}); rightOk {
				fmt.Printf("  ✅ Left et Right extraits\n")
				if leftType, _ := left["type"].(string); leftType == "fieldAccess" {
					if rightType, _ := right["type"].(string); rightType == "fieldAccess" {
						// Condition de jointure détectée
						fmt.Printf("  ✅ Condition de jointure fieldAccess détectée\n")
						leftObj, _ := left["object"].(string)
						leftField, _ := left["field"].(string)
						rightObj, _ := right["object"].(string)
						rightField, _ := right["field"].(string)
						operator, _ := condition["operator"].(string)

						fmt.Printf("    📌 %s.%s %s %s.%s\n", leftObj, leftField, operator, rightObj, rightField)

						joinConditions = append(joinConditions, JoinCondition{
							LeftField:  leftField,
							RightField: rightField,
							LeftVar:    leftObj,
							RightVar:   rightObj,
							Operator:   operator,
						})
					}
				}
			}
		}
	}

	return joinConditions
}

// ========== NŒUD EXISTS ==========

// ExistsNode représente un nœud d'existence dans le réseau RETE
type ExistsNode struct {
	BaseNode
	Condition       map[string]interface{} `json:"condition"`
	MainVariable    string                 `json:"main_variable"`    // Variable principale (p)
	ExistsVariable  string                 `json:"exists_variable"`  // Variable d'existence (o)
	VariableTypes   map[string]string      `json:"variable_types"`   // Mapping variable -> type
	ExistsCondition []JoinCondition        `json:"exists_condition"` // Condition d'existence (o.customer_id == p.id)
	mutex           sync.RWMutex
	// Mémoires pour architecture RETE
	MainMemory   *WorkingMemory // Faits de la variable principale
	ExistsMemory *WorkingMemory // Faits pour vérification d'existence
	ResultMemory *WorkingMemory // Tokens avec existence vérifiée
}

// NewExistsNode crée un nouveau nœud d'existence
func NewExistsNode(nodeID string, condition map[string]interface{}, mainVar string, existsVar string, varTypes map[string]string, storage Storage) *ExistsNode {
	return &ExistsNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "exists",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Condition:       condition,
		MainVariable:    mainVar,
		ExistsVariable:  existsVar,
		VariableTypes:   varTypes,
		ExistsCondition: extractJoinConditions(condition),
		// Initialiser les mémoires séparées
		MainMemory:   &WorkingMemory{NodeID: nodeID + "_main", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ExistsMemory: &WorkingMemory{NodeID: nodeID + "_exists", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory: &WorkingMemory{NodeID: nodeID + "_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
	}
}

// ActivateLeft traite les faits de la variable principale
func (en *ExistsNode) ActivateLeft(token *Token) error {
	fmt.Printf("🔍 EXISTSNODE[%s]: ActivateLeft - token %s\n", en.ID, token.ID)

	// Stocker le token dans la mémoire principale
	en.mutex.Lock()
	en.MainMemory.AddToken(token)
	en.mutex.Unlock()

	// Vérifier s'il existe des faits correspondants
	if en.checkExistence(token) {
		fmt.Printf("🔍 EXISTSNODE[%s]: Existence vérifiée pour %s\n", en.ID, token.ID)

		// Stocker le token avec existence vérifiée
		token.IsJoinResult = true // Marquer comme résultat validé
		en.mutex.Lock()
		en.ResultMemory.AddToken(token)
		en.Memory.AddToken(token) // Pour compatibilité avec le comptage
		en.mutex.Unlock()

		// Propager le token
		if err := en.PropagateToChildren(nil, token); err != nil {
			return err
		}
	} else {
		fmt.Printf("🔍 EXISTSNODE[%s]: Aucune existence trouvée pour %s\n", en.ID, token.ID)
	}

	return nil
}

// ActivateRight traite les faits pour vérification d'existence
func (en *ExistsNode) ActivateRight(fact *Fact) error {
	fmt.Printf("🔍 EXISTSNODE[%s]: ActivateRight - %s\n", en.ID, fact.Type)

	// Stocker le fait dans la mémoire d'existence
	en.mutex.Lock()
	en.ExistsMemory.AddFact(fact)
	en.mutex.Unlock()

	// Re-vérifier tous les tokens principaux avec ce nouveau fait
	mainTokens := en.MainMemory.GetTokens()
	for _, mainToken := range mainTokens {
		if en.checkExistence(mainToken) && !en.isAlreadyValidated(mainToken) {
			fmt.Printf("🔍 EXISTSNODE[%s]: Nouvelle existence vérifiée pour %s\n", en.ID, mainToken.ID)

			// Stocker le token avec existence vérifiée
			validatedToken := &Token{
				ID:           mainToken.ID + "_validated",
				Facts:        mainToken.Facts,
				NodeID:       en.ID,
				Bindings:     mainToken.Bindings,
				IsJoinResult: true,
			}

			en.mutex.Lock()
			en.ResultMemory.AddToken(validatedToken)
			en.Memory.AddToken(validatedToken)
			en.mutex.Unlock()

			// Propager le token
			if err := en.PropagateToChildren(nil, validatedToken); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkExistence vérifie si un token principal a des faits correspondants
func (en *ExistsNode) checkExistence(mainToken *Token) bool {
	existsFacts := en.ExistsMemory.GetFacts()

	// Récupérer le fait principal du token
	if len(mainToken.Facts) == 0 {
		return false
	}
	mainFact := mainToken.Facts[0]

	// Vérifier les conditions d'existence
	for _, existsFact := range existsFacts {
		if en.evaluateExistsCondition(mainFact, existsFact) {
			fmt.Printf("🔍 EXISTSNODE[%s]: Condition EXISTS satisfaite: %s ↔ %s\n", en.ID, mainFact.ID, existsFact.ID)
			return true
		}
	}

	return false
}

// evaluateExistsCondition évalue la condition d'existence entre deux faits
func (en *ExistsNode) evaluateExistsCondition(mainFact *Fact, existsFact *Fact) bool {
	fmt.Printf("🔍 EXISTSNODE[%s]: Évaluation condition EXISTS\n", en.ID)
	fmt.Printf("  📊 MainFact: %s (ID: %s)\n", mainFact.Type, mainFact.ID)
	fmt.Printf("  📊 ExistsFact: %s (ID: %s)\n", existsFact.Type, existsFact.ID)
	fmt.Printf("  📊 Conditions: %d à vérifier\n", len(en.ExistsCondition))

	for i, condition := range en.ExistsCondition {
		fmt.Printf("    Condition %d: %s.%s %s %s.%s\n", i,
			condition.LeftVar, condition.LeftField, condition.Operator,
			condition.RightVar, condition.RightField)

		// Déterminer quel fait correspond à quelle variable
		var leftFact, rightFact *Fact

		if condition.LeftVar == en.MainVariable {
			leftFact = mainFact
			rightFact = existsFact
			fmt.Printf("    → MainFact comme LeftVar (%s), ExistsFact comme RightVar (%s)\n", condition.LeftVar, condition.RightVar)
		} else if condition.LeftVar == en.ExistsVariable {
			leftFact = existsFact
			rightFact = mainFact
			fmt.Printf("    → ExistsFact comme LeftVar (%s), MainFact comme RightVar (%s)\n", condition.LeftVar, condition.RightVar)
		} else {
			fmt.Printf("    ❌ Variable %s non trouvée dans MainVariable:%s ou ExistsVariable:%s\n", condition.LeftVar, en.MainVariable, en.ExistsVariable)
			continue
		}

		leftValue := leftFact.Fields[condition.LeftField]
		rightValue := rightFact.Fields[condition.RightField]

		fmt.Printf("    🔍 Condition %d: %v %s %v\n", i, leftValue, condition.Operator, rightValue)

		switch condition.Operator {
		case "==":
			if leftValue != rightValue {
				fmt.Printf("    ❌ Condition %d échoue: %v != %v\n", i, leftValue, rightValue)
				return false
			}
			fmt.Printf("    ✅ Condition %d réussie: %v == %v\n", i, leftValue, rightValue)
		case "!=":
			if leftValue == rightValue {
				fmt.Printf("    ❌ Condition %d échoue: %v == %v\n", i, leftValue, rightValue)
				return false
			}
			fmt.Printf("    ✅ Condition %d réussie: %v != %v\n", i, leftValue, rightValue)
		default:
			fmt.Printf("    ❌ Opérateur non supporté: %s\n", condition.Operator)
			return false
		}
	}

	fmt.Printf("  ✅ Toutes les conditions EXISTS satisfaites\n")
	return true
}

// isAlreadyValidated vérifie si un token a déjà été validé
func (en *ExistsNode) isAlreadyValidated(token *Token) bool {
	validatedTokens := en.ResultMemory.GetTokens()
	for _, validatedToken := range validatedTokens {
		if validatedToken.ID == token.ID+"_validated" || validatedToken.ID == token.ID {
			return true
		}
	}
	return false
}
