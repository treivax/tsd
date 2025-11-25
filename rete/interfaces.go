package rete

type Node interface {
	GetID() string
	GetType() string
	GetMemory() *WorkingMemory
	ActivateLeft(token *Token) error
	ActivateRight(fact *Fact) error
	ActivateRetract(factID string) error
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
