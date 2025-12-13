// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// DebugLogger provides thread-safe debug logging to stderr for E2E debugging
type DebugLogger struct {
	enabled bool
	mu      sync.Mutex
}

var (
	globalDebugLogger *DebugLogger
	once              sync.Once
)

// GetDebugLogger returns the singleton debug logger instance
func GetDebugLogger() *DebugLogger {
	once.Do(func() {
		// Check environment variable to enable debug logging
		enabled := os.Getenv("TSD_DEBUG_BINDINGS") == "1"
		globalDebugLogger = &DebugLogger{
			enabled: enabled,
		}
	})
	return globalDebugLogger
}

// Enable turns on debug logging
func (dl *DebugLogger) Enable() {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	dl.enabled = true
}

// Disable turns off debug logging
func (dl *DebugLogger) Disable() {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	dl.enabled = false
}

// IsEnabled returns whether debug logging is enabled
func (dl *DebugLogger) IsEnabled() bool {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	return dl.enabled
}

// Log prints a debug message to stderr if enabled
func (dl *DebugLogger) Log(format string, args ...interface{}) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()
	fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
}

// LogJoinNode logs detailed information about a join node
func (dl *DebugLogger) LogJoinNode(nodeID string, operation string, details map[string]interface{}) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	var parts []string
	parts = append(parts, fmt.Sprintf("🔗 [JOIN_%s] %s", nodeID, operation))

	for k, v := range details {
		parts = append(parts, fmt.Sprintf("  %s: %v", k, v))
	}

	fmt.Fprintf(os.Stderr, "%s\n", strings.Join(parts, "\n"))
}

// LogBindings logs the current bindings in a token
func (dl *DebugLogger) LogBindings(prefix string, bindings *BindingChain) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	vars := bindings.Variables()
	fmt.Fprintf(os.Stderr, "[DEBUG] %s - Bindings: [%s]\n", prefix, strings.Join(vars, ", "))

	for _, varName := range vars {
		fact := bindings.Get(varName)
		fmt.Fprintf(os.Stderr, "[DEBUG]   %s -> %v\n", varName, fact)
	}
}

// LogJoinConditionEvaluation logs the evaluation of a single join condition
func (dl *DebugLogger) LogJoinConditionEvaluation(nodeID string, condIndex int, leftVar, rightVar string, leftFound, rightFound bool, leftVal, rightVal interface{}, result bool) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	status := "✓ PASS"
	if !result {
		status = "✗ FAIL"
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] 🔍 [JOIN_%s] Condition[%d] %s\n", nodeID, condIndex, status)
	fmt.Fprintf(os.Stderr, "[DEBUG]   LeftVar=%s (found=%v, val=%v)\n", leftVar, leftFound, leftVal)
	fmt.Fprintf(os.Stderr, "[DEBUG]   RightVar=%s (found=%v, val=%v)\n", rightVar, rightFound, rightVal)
	fmt.Fprintf(os.Stderr, "[DEBUG]   Result: %v\n", result)
}

// LogMemorySizes logs the sizes of join node memories
func (dl *DebugLogger) LogMemorySizes(nodeID string, leftSize, rightSize, resultSize int) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	fmt.Fprintf(os.Stderr, "[DEBUG] 📊 [JOIN_%s] Memory sizes: Left=%d, Right=%d, Result=%d\n",
		nodeID, leftSize, rightSize, resultSize)
}

// LogFactSubmission logs when a fact is submitted to the network
func (dl *DebugLogger) LogFactSubmission(factType string, factID string, factData interface{}) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	fmt.Fprintf(os.Stderr, "[DEBUG] 📥 FACT SUBMITTED: Type=%s, ID=%s, Data=%v\n",
		factType, factID, factData)
}

// LogTokenCreation logs when a new token is created
func (dl *DebugLogger) LogTokenCreation(nodeID string, tokenType string, bindings *BindingChain) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	vars := bindings.Variables()
	fmt.Fprintf(os.Stderr, "[DEBUG] 🎫 [%s] Token created (%s): vars=[%s]\n",
		nodeID, tokenType, strings.Join(vars, ", "))
}

// LogNetworkStructure logs the structure of the network (called after build)
func (dl *DebugLogger) LogNetworkStructure(network *ReteNetwork) {
	if !dl.IsEnabled() {
		return
	}
	dl.mu.Lock()
	defer dl.mu.Unlock()

	fmt.Fprintf(os.Stderr, "\n[DEBUG] ========== NETWORK STRUCTURE ==========\n")

	// Log type nodes
	fmt.Fprintf(os.Stderr, "[DEBUG] TypeNodes:\n")
	for typeName, typeNode := range network.TypeNodes {
		childCount := len(typeNode.GetChildren())
		fmt.Fprintf(os.Stderr, "[DEBUG]   - %s (children: %d)\n", typeName, childCount)
		for _, child := range typeNode.GetChildren() {
			fmt.Fprintf(os.Stderr, "[DEBUG]       -> %s (%s)\n", child.GetID(), child.GetType())
		}
	}

	// Log passthrough alphas
	fmt.Fprintf(os.Stderr, "[DEBUG] Passthrough Alphas:\n")
	for key, alpha := range network.PassthroughRegistry {
		childCount := len(alpha.GetChildren())
		var side string
		if condMap, ok := alpha.Condition.(map[string]interface{}); ok {
			if s, exists := condMap["side"].(string); exists {
				side = s
			}
		}
		fmt.Fprintf(os.Stderr, "[DEBUG]   - %s (side=%s, children: %d)\n", key, side, childCount)
		for _, child := range alpha.GetChildren() {
			fmt.Fprintf(os.Stderr, "[DEBUG]       -> %s (%s)\n", child.GetID(), child.GetType())
		}
	}

	// Log beta (join) nodes
	fmt.Fprintf(os.Stderr, "[DEBUG] Beta (Join) Nodes:\n")
	for nodeID, nodeInterface := range network.BetaNodes {
		if joinNode, ok := nodeInterface.(*JoinNode); ok {
			childCount := len(joinNode.GetChildren())
			fmt.Fprintf(os.Stderr, "[DEBUG]   - %s\n", nodeID)
			fmt.Fprintf(os.Stderr, "[DEBUG]       LeftVars: %v\n", joinNode.LeftVariables)
			fmt.Fprintf(os.Stderr, "[DEBUG]       RightVars: %v\n", joinNode.RightVariables)
			fmt.Fprintf(os.Stderr, "[DEBUG]       AllVars: %v\n", joinNode.AllVariables)
			fmt.Fprintf(os.Stderr, "[DEBUG]       JoinConditions: %d\n", len(joinNode.JoinConditions))
			for i, cond := range joinNode.JoinConditions {
				fmt.Fprintf(os.Stderr, "[DEBUG]         [%d] %s.%s %s %s.%s\n",
					i, cond.LeftVar, cond.LeftField, cond.Operator, cond.RightVar, cond.RightField)
			}
			fmt.Fprintf(os.Stderr, "[DEBUG]       Children: %d\n", childCount)
			for _, child := range joinNode.GetChildren() {
				fmt.Fprintf(os.Stderr, "[DEBUG]         -> %s (%s)\n", child.GetID(), child.GetType())
			}
		}
	}

	// Log terminal nodes
	fmt.Fprintf(os.Stderr, "[DEBUG] Terminal Nodes:\n")
	for nodeID, termNode := range network.TerminalNodes {
		actionName := "unknown"
		if termNode.Action != nil && termNode.Action.Job != nil {
			actionName = termNode.Action.Job.Name
		}
		fmt.Fprintf(os.Stderr, "[DEBUG]   - %s (action: %s)\n", nodeID, actionName)
	}

	fmt.Fprintf(os.Stderr, "[DEBUG] =======================================\n\n")
}
