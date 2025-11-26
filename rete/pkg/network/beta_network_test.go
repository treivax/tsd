package network

import (
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// MockLogger implémente domain.Logger pour les tests
type MockLogger struct {
	debugCalls []LogCall
	infoCalls  []LogCall
	warnCalls  []LogCall
	errorCalls []LogCall
}

type LogCall struct {
	Message string
	Fields  map[string]interface{}
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugCalls: make([]LogCall, 0),
		infoCalls:  make([]LogCall, 0),
		warnCalls:  make([]LogCall, 0),
		errorCalls: make([]LogCall, 0),
	}
}

func (ml *MockLogger) Debug(msg string, fields map[string]interface{}) {
	ml.debugCalls = append(ml.debugCalls, LogCall{Message: msg, Fields: fields})
}

func (ml *MockLogger) Info(msg string, fields map[string]interface{}) {
	ml.infoCalls = append(ml.infoCalls, LogCall{Message: msg, Fields: fields})
}

func (ml *MockLogger) Warn(msg string, fields map[string]interface{}) {
	ml.warnCalls = append(ml.warnCalls, LogCall{Message: msg, Fields: fields})
}

func (ml *MockLogger) Error(msg string, err error, fields map[string]interface{}) {
	ml.errorCalls = append(ml.errorCalls, LogCall{Message: msg, Fields: fields})
}

// ===== BetaNetworkBuilder Tests =====

func TestNewBetaNetworkBuilder(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	if builder == nil {
		t.Fatal("NewBetaNetworkBuilder should not return nil")
	}
	if builder.logger == nil {
		t.Error("Logger should be set")
	}
	if builder.betaNodes == nil {
		t.Error("betaNodes map should be initialized")
	}
	if len(builder.betaNodes) != 0 {
		t.Errorf("Expected 0 nodes initially, got %d", len(builder.betaNodes))
	}
}

func TestBetaNetworkBuilder_CreateJoinNode(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	// Create join conditions
	condition := domain.NewBasicJoinCondition("field1", "field2", "==")
	conditions := []domain.JoinCondition{condition}

	// Create join node
	joinNode := builder.CreateJoinNode("join1", conditions)

	if joinNode == nil {
		t.Fatal("CreateJoinNode should not return nil")
	}
	if joinNode.ID() != "join1" {
		t.Errorf("Expected node ID 'join1', got '%s'", joinNode.ID())
	}

	// Verify conditions are set
	nodeConditions := joinNode.GetJoinConditions()
	if len(nodeConditions) != 1 {
		t.Errorf("Expected 1 condition, got %d", len(nodeConditions))
	}

	// Verify node is registered
	node, exists := builder.GetBetaNode("join1")
	if !exists {
		t.Error("Join node should be registered")
	}
	if node.ID() != "join1" {
		t.Errorf("Expected registered node ID 'join1', got '%s'", node.ID())
	}

	// Verify logging
	if len(logger.infoCalls) != 1 {
		t.Errorf("Expected 1 info log call, got %d", len(logger.infoCalls))
	}
}

func TestBetaNetworkBuilder_CreateJoinNode_MultipleConditions(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	conditions := []domain.JoinCondition{
		domain.NewBasicJoinCondition("field1", "field2", "=="),
		domain.NewBasicJoinCondition("field3", "field4", ">"),
		domain.NewBasicJoinCondition("field5", "field6", "<="),
	}

	joinNode := builder.CreateJoinNode("join1", conditions)

	if joinNode == nil {
		t.Fatal("CreateJoinNode should not return nil")
	}

	nodeConditions := joinNode.GetJoinConditions()
	if len(nodeConditions) != 3 {
		t.Errorf("Expected 3 conditions, got %d", len(nodeConditions))
	}
}

func TestBetaNetworkBuilder_CreateJoinNode_EmptyConditions(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	joinNode := builder.CreateJoinNode("join1", []domain.JoinCondition{})

	if joinNode == nil {
		t.Fatal("CreateJoinNode should not return nil")
	}

	nodeConditions := joinNode.GetJoinConditions()
	if len(nodeConditions) != 0 {
		t.Errorf("Expected 0 conditions, got %d", len(nodeConditions))
	}
}

func TestBetaNetworkBuilder_CreateBetaNode(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	betaNode := builder.CreateBetaNode("beta1")

	if betaNode == nil {
		t.Fatal("CreateBetaNode should not return nil")
	}
	if betaNode.ID() != "beta1" {
		t.Errorf("Expected node ID 'beta1', got '%s'", betaNode.ID())
	}

	// Verify node is registered
	node, exists := builder.GetBetaNode("beta1")
	if !exists {
		t.Error("Beta node should be registered")
	}
	if node.ID() != "beta1" {
		t.Errorf("Expected registered node ID 'beta1', got '%s'", node.ID())
	}

	// Verify logging
	if len(logger.infoCalls) != 1 {
		t.Errorf("Expected 1 info log call, got %d", len(logger.infoCalls))
	}
}

func TestBetaNetworkBuilder_CreateMultipleNodes(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	builder.CreateJoinNode("join1", []domain.JoinCondition{})
	builder.CreateJoinNode("join2", []domain.JoinCondition{})
	builder.CreateBetaNode("beta1")
	builder.CreateBetaNode("beta2")

	nodes := builder.ListBetaNodes()
	if len(nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(nodes))
	}

	expectedIDs := []string{"join1", "join2", "beta1", "beta2"}
	for _, id := range expectedIDs {
		if _, exists := nodes[id]; !exists {
			t.Errorf("Node '%s' should exist", id)
		}
	}
}

func TestBetaNetworkBuilder_ConnectNodes(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	parent := builder.CreateBetaNode("parent")
	_ = builder.CreateBetaNode("child")

	err := builder.ConnectNodes("parent", "child")
	if err != nil {
		t.Fatalf("ConnectNodes should not return error: %v", err)
	}

	// Verify connection
	children := parent.GetChildren()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}
	if len(children) > 0 && children[0].ID() != "child" {
		t.Errorf("Expected child ID 'child', got '%s'", children[0].ID())
	}

	// Verify logging
	foundConnectLog := false
	for _, call := range logger.infoCalls {
		if call.Message == "connected nodes" {
			foundConnectLog = true
			break
		}
	}
	if !foundConnectLog {
		t.Error("Expected 'connected nodes' log message")
	}
}

func TestBetaNetworkBuilder_ConnectNodes_ParentNotFound(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	builder.CreateBetaNode("child")

	err := builder.ConnectNodes("nonexistent", "child")
	if err == nil {
		t.Error("ConnectNodes should return error when parent not found")
	}
	if err != nil && !contains(err.Error(), "parent node not found") {
		t.Errorf("Expected 'parent node not found' error, got: %v", err)
	}
}

func TestBetaNetworkBuilder_ConnectNodes_ChildNotFound(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	builder.CreateBetaNode("parent")

	err := builder.ConnectNodes("parent", "nonexistent")
	if err == nil {
		t.Error("ConnectNodes should return error when child not found")
	}
	if err != nil && !contains(err.Error(), "child node not found") {
		t.Errorf("Expected 'child node not found' error, got: %v", err)
	}
}

func TestBetaNetworkBuilder_ConnectNodes_BothNotFound(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	err := builder.ConnectNodes("nonexistent_parent", "nonexistent_child")
	if err == nil {
		t.Error("ConnectNodes should return error when both nodes not found")
	}
}

func TestBetaNetworkBuilder_GetBetaNode(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	builder.CreateBetaNode("beta1")

	// Existing node
	node, exists := builder.GetBetaNode("beta1")
	if !exists {
		t.Error("GetBetaNode should return true for existing node")
	}
	if node == nil {
		t.Error("GetBetaNode should return non-nil node")
	}
	if node != nil && node.ID() != "beta1" {
		t.Errorf("Expected node ID 'beta1', got '%s'", node.ID())
	}

	// Non-existing node
	node, exists = builder.GetBetaNode("nonexistent")
	if exists {
		t.Error("GetBetaNode should return false for non-existing node")
	}
	if node != nil {
		t.Error("GetBetaNode should return nil for non-existing node")
	}
}

func TestBetaNetworkBuilder_ListBetaNodes(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	// Empty list
	nodes := builder.ListBetaNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected empty list, got %d nodes", len(nodes))
	}

	// Add nodes
	builder.CreateJoinNode("join1", []domain.JoinCondition{})
	builder.CreateBetaNode("beta1")
	builder.CreateBetaNode("beta2")

	nodes = builder.ListBetaNodes()
	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(nodes))
	}

	// Verify all nodes are present
	if _, exists := nodes["join1"]; !exists {
		t.Error("join1 should be in list")
	}
	if _, exists := nodes["beta1"]; !exists {
		t.Error("beta1 should be in list")
	}
	if _, exists := nodes["beta2"]; !exists {
		t.Error("beta2 should be in list")
	}
}

func TestBetaNetworkBuilder_ClearNetwork(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	// Add nodes
	builder.CreateJoinNode("join1", []domain.JoinCondition{})
	builder.CreateBetaNode("beta1")
	builder.CreateBetaNode("beta2")

	nodes := builder.ListBetaNodes()
	if len(nodes) != 3 {
		t.Fatalf("Expected 3 nodes before clear, got %d", len(nodes))
	}

	// Clear network
	builder.ClearNetwork()

	nodes = builder.ListBetaNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes after clear, got %d", len(nodes))
	}

	// Verify logging
	foundClearLog := false
	for _, call := range logger.infoCalls {
		if call.Message == "cleared beta network" {
			foundClearLog = true
			break
		}
	}
	if !foundClearLog {
		t.Error("Expected 'cleared beta network' log message")
	}
}

func TestBetaNetworkBuilder_BuildMultiJoinNetwork(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	pattern := MultiJoinPattern{
		PatternID: "pattern1",
		JoinSpecs: []JoinSpecification{
			{
				LeftType:   "Person",
				RightType:  "Order",
				Conditions: []domain.JoinCondition{domain.NewBasicJoinCondition("id", "personId", "==")},
				NodeID:     "join1",
			},
			{
				LeftType:   "Order",
				RightType:  "Product",
				Conditions: []domain.JoinCondition{domain.NewBasicJoinCondition("productId", "id", "==")},
				NodeID:     "join2",
			},
		},
		FinalAction: "process",
	}

	nodes, err := builder.BuildMultiJoinNetwork(pattern)
	if err != nil {
		t.Fatalf("BuildMultiJoinNetwork should not return error: %v", err)
	}

	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	// Verify first node
	if nodes[0].ID() != "join1" {
		t.Errorf("Expected first node ID 'join1', got '%s'", nodes[0].ID())
	}

	// Verify second node
	if nodes[1].ID() != "join2" {
		t.Errorf("Expected second node ID 'join2', got '%s'", nodes[1].ID())
	}

	// Verify nodes are connected
	children := nodes[0].GetChildren()
	if len(children) != 1 {
		t.Errorf("Expected first node to have 1 child, got %d", len(children))
	}
	if len(children) > 0 && children[0].ID() != "join2" {
		t.Errorf("Expected child ID 'join2', got '%s'", children[0].ID())
	}

	// Verify nodes are registered
	_, exists := builder.GetBetaNode("join1")
	if !exists {
		t.Error("join1 should be registered")
	}
	_, exists = builder.GetBetaNode("join2")
	if !exists {
		t.Error("join2 should be registered")
	}
}

func TestBetaNetworkBuilder_BuildMultiJoinNetwork_EmptyPattern(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	pattern := MultiJoinPattern{
		PatternID:   "empty_pattern",
		JoinSpecs:   []JoinSpecification{},
		FinalAction: "none",
	}

	nodes, err := builder.BuildMultiJoinNetwork(pattern)
	if err != nil {
		t.Fatalf("BuildMultiJoinNetwork should not return error: %v", err)
	}

	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes for empty pattern, got %d", len(nodes))
	}
}

func TestBetaNetworkBuilder_BuildMultiJoinNetwork_AutoGenerateIDs(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	pattern := MultiJoinPattern{
		PatternID: "pattern1",
		JoinSpecs: []JoinSpecification{
			{
				LeftType:   "Person",
				RightType:  "Order",
				Conditions: []domain.JoinCondition{},
				NodeID:     "", // Empty - should auto-generate
			},
			{
				LeftType:   "Order",
				RightType:  "Product",
				Conditions: []domain.JoinCondition{},
				NodeID:     "", // Empty - should auto-generate
			},
		},
		FinalAction: "process",
	}

	nodes, err := builder.BuildMultiJoinNetwork(pattern)
	if err != nil {
		t.Fatalf("BuildMultiJoinNetwork should not return error: %v", err)
	}

	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}

	// Verify auto-generated IDs
	expectedID1 := "pattern1_join_0"
	expectedID2 := "pattern1_join_1"

	if nodes[0].ID() != expectedID1 {
		t.Errorf("Expected auto-generated ID '%s', got '%s'", expectedID1, nodes[0].ID())
	}
	if nodes[1].ID() != expectedID2 {
		t.Errorf("Expected auto-generated ID '%s', got '%s'", expectedID2, nodes[1].ID())
	}
}

func TestBetaNetworkBuilder_BuildMultiJoinNetwork_SingleJoin(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	pattern := MultiJoinPattern{
		PatternID: "single",
		JoinSpecs: []JoinSpecification{
			{
				LeftType:   "Person",
				RightType:  "Order",
				Conditions: []domain.JoinCondition{domain.NewBasicJoinCondition("id", "personId", "==")},
				NodeID:     "only_join",
			},
		},
		FinalAction: "action",
	}

	nodes, err := builder.BuildMultiJoinNetwork(pattern)
	if err != nil {
		t.Fatalf("BuildMultiJoinNetwork should not return error: %v", err)
	}

	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}

	// Single node should have no children initially
	children := nodes[0].GetChildren()
	if len(children) != 0 {
		t.Errorf("Expected 0 children for single node, got %d", len(children))
	}
}

func TestBetaNetworkBuilder_NetworkStatistics_Empty(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	stats := builder.NetworkStatistics()

	if stats.TotalNodes != 0 {
		t.Errorf("Expected 0 total nodes, got %d", stats.TotalNodes)
	}
	if stats.JoinNodes != 0 {
		t.Errorf("Expected 0 join nodes, got %d", stats.JoinNodes)
	}
	if stats.SimpleBetaNodes != 0 {
		t.Errorf("Expected 0 simple beta nodes, got %d", stats.SimpleBetaNodes)
	}
	if stats.TotalTokens != 0 {
		t.Errorf("Expected 0 total tokens, got %d", stats.TotalTokens)
	}
	if stats.TotalFacts != 0 {
		t.Errorf("Expected 0 total facts, got %d", stats.TotalFacts)
	}
	if len(stats.MemoryStats) != 0 {
		t.Errorf("Expected 0 memory stats entries, got %d", len(stats.MemoryStats))
	}
}

func TestBetaNetworkBuilder_NetworkStatistics_WithNodes(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	// Create various nodes
	builder.CreateJoinNode("join1", []domain.JoinCondition{})
	builder.CreateJoinNode("join2", []domain.JoinCondition{})
	builder.CreateBetaNode("beta1")
	builder.CreateBetaNode("beta2")
	builder.CreateBetaNode("beta3")

	stats := builder.NetworkStatistics()

	if stats.TotalNodes != 5 {
		t.Errorf("Expected 5 total nodes, got %d", stats.TotalNodes)
	}
	if stats.JoinNodes != 2 {
		t.Errorf("Expected 2 join nodes, got %d", stats.JoinNodes)
	}
	if stats.SimpleBetaNodes != 3 {
		t.Errorf("Expected 3 simple beta nodes, got %d", stats.SimpleBetaNodes)
	}
	if len(stats.MemoryStats) != 5 {
		t.Errorf("Expected 5 memory stats entries, got %d", len(stats.MemoryStats))
	}

	// Verify all nodes have memory stats
	expectedNodeIDs := []string{"join1", "join2", "beta1", "beta2", "beta3"}
	for _, nodeID := range expectedNodeIDs {
		if _, exists := stats.MemoryStats[nodeID]; !exists {
			t.Errorf("Memory stats should exist for node '%s'", nodeID)
		}
	}
}

func TestBetaNetworkBuilder_NetworkStatistics_OnlyJoinNodes(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	builder.CreateJoinNode("join1", []domain.JoinCondition{})
	builder.CreateJoinNode("join2", []domain.JoinCondition{})

	stats := builder.NetworkStatistics()

	if stats.TotalNodes != 2 {
		t.Errorf("Expected 2 total nodes, got %d", stats.TotalNodes)
	}
	if stats.JoinNodes != 2 {
		t.Errorf("Expected 2 join nodes, got %d", stats.JoinNodes)
	}
	if stats.SimpleBetaNodes != 0 {
		t.Errorf("Expected 0 simple beta nodes, got %d", stats.SimpleBetaNodes)
	}
}

func TestBetaNetworkBuilder_NetworkStatistics_OnlyBetaNodes(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	builder.CreateBetaNode("beta1")
	builder.CreateBetaNode("beta2")
	builder.CreateBetaNode("beta3")

	stats := builder.NetworkStatistics()

	if stats.TotalNodes != 3 {
		t.Errorf("Expected 3 total nodes, got %d", stats.TotalNodes)
	}
	if stats.JoinNodes != 0 {
		t.Errorf("Expected 0 join nodes, got %d", stats.JoinNodes)
	}
	if stats.SimpleBetaNodes != 3 {
		t.Errorf("Expected 3 simple beta nodes, got %d", stats.SimpleBetaNodes)
	}
}

func TestBetaNetworkBuilder_ConcurrentAccess(t *testing.T) {
	logger := NewMockLogger()
	builder := NewBetaNetworkBuilder(logger)

	// Create nodes concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			nodeID := string(rune('A' + id))
			builder.CreateBetaNode(nodeID)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Read nodes concurrently
	for i := 0; i < 10; i++ {
		go func(id int) {
			nodeID := string(rune('A' + id))
			builder.GetBetaNode(nodeID)
			builder.ListBetaNodes()
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all nodes were created
	nodes := builder.ListBetaNodes()
	if len(nodes) != 10 {
		t.Errorf("Expected 10 nodes after concurrent creation, got %d", len(nodes))
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
