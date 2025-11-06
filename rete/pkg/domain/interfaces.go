package domain

// Node interface pour tous les nœuds du réseau RETE.
// Utilise le principe de ségrégation d'interfaces (ISP) pour réduire le couplage.
type Node interface {
	ID() string
	Type() string
	ProcessFact(*Fact) error
}

// MemoryNode interface pour les nœuds qui maintiennent une mémoire de travail.
type MemoryNode interface {
	Node
	GetMemory() *WorkingMemory
}

// ParentNode interface pour les nœuds qui ont des enfants.
type ParentNode interface {
	Node
	AddChild(Node)
	GetChildren() []Node
}

// Storage interface pour la persistance de l'état des nœuds.
type Storage interface {
	Save(nodeID string, data interface{}) error
	Load(nodeID string, result interface{}) error
	Delete(nodeID string) error
	List() ([]string, error)
}

// Logger interface pour le logging structuré.
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
}

// BetaNode interface pour les nœuds beta du réseau RETE.
// Gère les jointures multi-faits et la mémoire beta.
type BetaNode interface {
	MemoryNode
	ParentNode
	ProcessLeftToken(*Token) error
	ProcessRightFact(*Fact) error
	GetLeftMemory() []*Token
	GetRightMemory() []*Fact
	ClearMemory()
}

// JoinNode interface pour les nœuds de jointure.
// Applique le principe ISP (Interface Segregation Principle).
type JoinNode interface {
	BetaNode
	SetJoinConditions(conditions []JoinCondition)
	GetJoinConditions() []JoinCondition
	EvaluateJoin(*Token, *Fact) bool
}

// NotNode interface pour les nœuds de négation.
// Implémente la logique NOT dans l'algorithme RETE.
type NotNode interface {
	BetaNode
	SetNegationCondition(condition interface{})
	GetNegationCondition() interface{}
	ProcessNegation(*Token, *Fact) bool
}

// ExistsNode interface pour les nœuds de quantification existentielle.
// Vérifie l'existence d'au moins un fait satisfaisant une condition.
type ExistsNode interface {
	BetaNode
	SetExistenceCondition(variable TypedVariable, condition interface{})
	GetExistenceCondition() (TypedVariable, interface{})
	CheckExistence(*Token) bool
}

// AccumulateNode interface pour les nœuds d'agrégation.
// Effectue des calculs d'agrégation (SUM, COUNT, AVG, MIN, MAX).
type AccumulateNode interface {
	BetaNode
	SetAggregateFunction(function string, expression interface{})
	GetAggregateFunction() (string, interface{})
	ComputeAggregate([]*Fact) (interface{}, error)
	UpdateAccumulation(*Fact, bool) // true pour ajout, false pour suppression
}

// TypedVariable représente une variable typée dans les contraintes
type TypedVariable struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
}

// BetaMemory interface pour la gestion de la mémoire des nœuds beta.
type BetaMemory interface {
	StoreToken(*Token)
	RemoveToken(tokenID string) bool
	GetTokens() []*Token
	StoreFact(*Fact)
	RemoveFact(factID string) bool
	GetFacts() []*Fact
	Clear()
	Size() (tokens int, facts int)
}

// JoinCondition interface pour les conditions de jointure.
type JoinCondition interface {
	Evaluate(*Token, *Fact) bool
	GetLeftField() string
	GetRightField() string
	GetOperator() string
}

// EventListener interface pour les événements du réseau RETE.
type EventListener interface {
	OnFactSubmitted(fact *Fact)
	OnActionTriggered(action string, facts []*Fact)
	OnNodeActivated(nodeID string, fact *Fact)
	OnTokenCreated(token *Token, nodeID string)
	OnJoinPerformed(leftToken *Token, rightFact *Fact, resultToken *Token)
}

// AccumulateFunction définit une fonction d'agrégation
type AccumulateFunction struct {
	FunctionType string      // SUM, COUNT, AVG, MIN, MAX
	Field        string      // Champ sur lequel appliquer la fonction
	Condition    interface{} // Condition optionnelle
}
