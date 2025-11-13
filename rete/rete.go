package rete

import (
	"fmt"
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
	ID     string  `json:"id"`
	Facts  []*Fact `json:"facts"`
	NodeID string  `json:"node_id"`
	Parent *Token  `json:"parent,omitempty"`
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
	Type string   `json:"type"`
	Name string   `json:"name"`
	Args []string `json:"args"`
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

	// Évaluer la condition Alpha sur le fait
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
			// Format compact : Type[id=value, field=value, ...]
			fmt.Printf("%s[", fact.Type)
			fieldCount := 0
			for key, value := range fact.Fields {
				if fieldCount > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s=%v", key, value)
				fieldCount++
			}
			fmt.Print("]")
		}
		fmt.Print(")")
	}
	fmt.Println()

	return nil
}
