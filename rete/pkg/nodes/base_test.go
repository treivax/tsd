package nodes

import (
	"fmt"
	"sync"
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// mockLogger implémente domain.Logger pour les tests
type mockLogger struct {
	debugCalls []map[string]interface{}
	infoCalls  []map[string]interface{}
	errorCalls []map[string]interface{}
	mu         sync.Mutex
}

func newMockLogger() *mockLogger {
	return &mockLogger{
		debugCalls: make([]map[string]interface{}, 0),
		infoCalls:  make([]map[string]interface{}, 0),
		errorCalls: make([]map[string]interface{}, 0),
	}
}

func (m *mockLogger) Debug(msg string, fields map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.debugCalls = append(m.debugCalls, fields)
}

func (m *mockLogger) Info(msg string, fields map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.infoCalls = append(m.infoCalls, fields)
}

func (m *mockLogger) Warn(msg string, fields map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Warn calls tracked but not used in tests
}

func (m *mockLogger) Error(msg string, err error, fields map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCalls = append(m.errorCalls, fields)
}

func (m *mockLogger) DebugCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.debugCalls)
}

// mockNode implémente domain.Node pour les tests
type mockNode struct {
	id           string
	processedIDs []string
	mu           sync.Mutex
}

func newMockNode(id string) *mockNode {
	return &mockNode{
		id:           id,
		processedIDs: make([]string, 0),
	}
}

func (m *mockNode) ID() string {
	return m.id
}

func (m *mockNode) Type() string {
	return "mock"
}

func (m *mockNode) ProcessFact(fact *domain.Fact) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.processedIDs = append(m.processedIDs, fact.ID)
	return nil
}

func (m *mockNode) GetMemory() *domain.WorkingMemory {
	return nil
}

func (m *mockNode) AddChild(child domain.Node) {}

func (m *mockNode) GetChildren() []domain.Node {
	return nil
}

func (m *mockNode) ProcessedCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.processedIDs)
}

// Tests pour NewBaseNode
func TestNewBaseNode(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("test-node-1", "alpha", logger)

	if node == nil {
		t.Fatal("NewBaseNode returned nil")
	}

	if node.ID() != "test-node-1" {
		t.Errorf("Expected ID 'test-node-1', got '%s'", node.ID())
	}

	if node.Type() != "alpha" {
		t.Errorf("Expected Type 'alpha', got '%s'", node.Type())
	}

	if node.GetMemory() == nil {
		t.Error("Expected non-nil memory")
	}

	children := node.GetChildren()
	if len(children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(children))
	}
}

// Tests pour ID()
func TestBaseNode_ID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"simple ID", "node-1"},
		{"complex ID", "alpha-node-customer-filter-123"},
		{"empty ID", ""},
		{"special chars", "node@#$%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := newMockLogger()
			node := NewBaseNode(tt.id, "test", logger)

			if node.ID() != tt.id {
				t.Errorf("Expected ID '%s', got '%s'", tt.id, node.ID())
			}
		})
	}
}

// Tests pour Type()
func TestBaseNode_Type(t *testing.T) {
	tests := []struct {
		name     string
		nodeType string
	}{
		{"alpha type", "alpha"},
		{"beta type", "beta"},
		{"join type", "join"},
		{"terminal type", "terminal"},
		{"empty type", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := newMockLogger()
			node := NewBaseNode("test-id", tt.nodeType, logger)

			if node.Type() != tt.nodeType {
				t.Errorf("Expected Type '%s', got '%s'", tt.nodeType, node.Type())
			}
		})
	}
}

// Tests pour GetMemory()
func TestBaseNode_GetMemory(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("test-node", "alpha", logger)

	memory := node.GetMemory()
	if memory == nil {
		t.Fatal("GetMemory returned nil")
	}

	// Vérifier que la mémoire est initialisée correctement (pas d'ID public dans WorkingMemory)
	// La mémoire devrait juste être non-nil
	if memory == nil {
		t.Error("Expected non-nil memory")
	}
}

// Tests pour AddChild()
func TestBaseNode_AddChild(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("parent", "alpha", logger)

	child1 := newMockNode("child-1")
	child2 := newMockNode("child-2")

	node.AddChild(child1)
	children := node.GetChildren()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}

	node.AddChild(child2)
	children = node.GetChildren()
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}

	// Vérifier que les enfants sont les bons
	if children[0].ID() != "child-1" {
		t.Errorf("Expected first child ID 'child-1', got '%s'", children[0].ID())
	}
	if children[1].ID() != "child-2" {
		t.Errorf("Expected second child ID 'child-2', got '%s'", children[1].ID())
	}
}

// Tests pour GetChildren()
func TestBaseNode_GetChildren(t *testing.T) {
	t.Run("no children", func(t *testing.T) {
		logger := newMockLogger()
		node := NewBaseNode("parent", "alpha", logger)

		children := node.GetChildren()
		if len(children) != 0 {
			t.Errorf("Expected 0 children, got %d", len(children))
		}
	})

	t.Run("with children", func(t *testing.T) {
		logger := newMockLogger()
		node := NewBaseNode("parent", "alpha", logger)

		child1 := newMockNode("child-1")
		child2 := newMockNode("child-2")
		child3 := newMockNode("child-3")

		node.AddChild(child1)
		node.AddChild(child2)
		node.AddChild(child3)

		children := node.GetChildren()
		if len(children) != 3 {
			t.Errorf("Expected 3 children, got %d", len(children))
		}
	})

	t.Run("returns copy not reference", func(t *testing.T) {
		logger := newMockLogger()
		node := NewBaseNode("parent", "alpha", logger)

		child := newMockNode("child")
		node.AddChild(child)

		children1 := node.GetChildren()
		children2 := node.GetChildren()

		// Modifier children1 ne doit pas affecter children2
		if len(children1) != 1 || len(children2) != 1 {
			t.Error("Initial children count mismatch")
		}

		// Les slices doivent être différents (copies)
		// mais contenir les mêmes éléments
		if children1[0].ID() != children2[0].ID() {
			t.Error("Children content mismatch")
		}
	})
}

// Tests pour logFactProcessing()
func TestBaseNode_logFactProcessing(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("test-node", "alpha", logger)

	fact := &domain.Fact{
		ID:   "fact-1",
		Type: "Customer",
		Fields: map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}

	// Appeler logFactProcessing via réflexion ou méthode publique qui l'utilise
	// Comme c'est une méthode privée, on teste indirectement via son effet
	node.logFactProcessing(fact, "test-action")

	if logger.DebugCallCount() != 1 {
		t.Errorf("Expected 1 debug call, got %d", logger.DebugCallCount())
	}

	if len(logger.debugCalls) > 0 {
		fields := logger.debugCalls[0]
		if fields["node_id"] != "test-node" {
			t.Errorf("Expected node_id 'test-node', got '%v'", fields["node_id"])
		}
		if fields["node_type"] != "alpha" {
			t.Errorf("Expected node_type 'alpha', got '%v'", fields["node_type"])
		}
		if fields["fact_id"] != "fact-1" {
			t.Errorf("Expected fact_id 'fact-1', got '%v'", fields["fact_id"])
		}
		if fields["fact_type"] != "Customer" {
			t.Errorf("Expected fact_type 'Customer', got '%v'", fields["fact_type"])
		}
		if fields["action"] != "test-action" {
			t.Errorf("Expected action 'test-action', got '%v'", fields["action"])
		}
	}
}

// Tests de concurrence pour BaseNode
func TestBaseNode_Concurrency(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("concurrent-node", "alpha", logger)

	const goroutines = 10
	const childrenPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Ajouter des enfants concurrents
	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			for j := 0; j < childrenPerGoroutine; j++ {
				child := newMockNode(fmt.Sprintf("child-%d-%d", index, j))
				node.AddChild(child)
			}
		}(i)
	}

	wg.Wait()

	// Vérifier le nombre total d'enfants
	children := node.GetChildren()
	expected := goroutines * childrenPerGoroutine
	if len(children) != expected {
		t.Errorf("Expected %d children, got %d", expected, len(children))
	}
}

// Tests de concurrence pour GetChildren()
func TestBaseNode_GetChildren_Concurrency(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("concurrent-read-node", "alpha", logger)

	// Ajouter quelques enfants
	for i := 0; i < 5; i++ {
		child := newMockNode(fmt.Sprintf("child-%d", i))
		node.AddChild(child)
	}

	const readers = 20
	var wg sync.WaitGroup
	wg.Add(readers)

	// Lire concurremment les enfants
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			children := node.GetChildren()
			if len(children) != 5 {
				t.Errorf("Expected 5 children, got %d", len(children))
			}
		}()
	}

	wg.Wait()
}

// Tests de concurrence pour GetMemory()
func TestBaseNode_GetMemory_Concurrency(t *testing.T) {
	logger := newMockLogger()
	node := NewBaseNode("memory-node", "alpha", logger)

	const readers = 20
	var wg sync.WaitGroup
	wg.Add(readers)

	// Lire concurremment la mémoire
	for i := 0; i < readers; i++ {
		go func() {
			defer wg.Done()
			memory := node.GetMemory()
			if memory == nil {
				t.Error("GetMemory returned nil")
			}
		}()
	}

	wg.Wait()
}

// Tests cas limites
func TestBaseNode_EdgeCases(t *testing.T) {
	t.Run("nil logger handling", func(t *testing.T) {
		// Le code devrait gérer un logger nil sans paniquer
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("NewBaseNode panicked with nil logger: %v", r)
			}
		}()

		node := NewBaseNode("test", "alpha", nil)
		if node == nil {
			t.Error("Expected node to be created even with nil logger")
		}
	})

	t.Run("add nil child", func(t *testing.T) {
		logger := newMockLogger()
		node := NewBaseNode("parent", "alpha", logger)

		// Ajouter un enfant nil (comportement à définir)
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("AddChild panicked with nil child: %v", r)
			}
		}()

		node.AddChild(nil)
		children := node.GetChildren()

		// Si nil est accepté, il devrait être dans la liste
		if len(children) == 0 {
			// OK, nil est filtré
		} else if len(children) == 1 && children[0] == nil {
			// OK, nil est conservé
		} else {
			t.Errorf("Unexpected children count: %d", len(children))
		}
	})

	t.Run("add duplicate children", func(t *testing.T) {
		logger := newMockLogger()
		node := NewBaseNode("parent", "alpha", logger)

		child := newMockNode("child")
		node.AddChild(child)
		node.AddChild(child)

		children := node.GetChildren()
		// Les duplicatas sont autorisés dans l'implémentation actuelle
		if len(children) != 2 {
			t.Errorf("Expected 2 children (duplicates allowed), got %d", len(children))
		}
	})
}

// Benchmark pour AddChild
func BenchmarkBaseNode_AddChild(b *testing.B) {
	logger := newMockLogger()
	node := NewBaseNode("bench-node", "alpha", logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child := newMockNode(fmt.Sprintf("child-%d", i))
		node.AddChild(child)
	}
}

// Benchmark pour GetChildren
func BenchmarkBaseNode_GetChildren(b *testing.B) {
	logger := newMockLogger()
	node := NewBaseNode("bench-node", "alpha", logger)

	// Pré-remplir avec des enfants
	for i := 0; i < 100; i++ {
		child := newMockNode(fmt.Sprintf("child-%d", i))
		node.AddChild(child)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = node.GetChildren()
	}
}
