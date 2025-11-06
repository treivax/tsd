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

// EventListener interface pour les événements du réseau RETE.
type EventListener interface {
OnFactSubmitted(fact *Fact)
OnActionTriggered(action string, facts []*Fact)
OnNodeActivated(nodeID string, fact *Fact)
}
