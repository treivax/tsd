// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestEnvironment provides isolated test resources with automatic cleanup.
// It enables safe parallel testing by ensuring each test has its own isolated
// network, storage, logger, and temporary directory.
type TestEnvironment struct {
	Network   *ReteNetwork
	Storage   Storage
	Pipeline  *ConstraintPipeline
	Logger    *Logger
	LogBuffer *bytes.Buffer
	TempDir   string

	t       *testing.T
	cleanup []func()
	mu      sync.Mutex
}

// TestEnvOption is a functional option for configuring TestEnvironment.
type TestEnvOption func(*TestEnvironment)

// WithLogLevel sets the log level for the test environment.
func WithLogLevel(level LogLevel) TestEnvOption {
	return func(te *TestEnvironment) {
		te.Logger.SetLevel(level)
	}
}

// WithTempStorage creates a temporary file-based storage instead of in-memory.
func WithTempStorage() TestEnvOption {
	return func(te *TestEnvironment) {
		// For now, we'll use memory storage
		// Future: implement file-based storage if needed
		te.Storage = NewMemoryStorage()
	}
}

// WithCustomStorage allows providing a custom storage implementation.
func WithCustomStorage(storage Storage) TestEnvOption {
	return func(te *TestEnvironment) {
		te.Storage = storage
	}
}

// WithLogOutput sets a custom output writer for the logger.
func WithLogOutput(w io.Writer) TestEnvOption {
	return func(te *TestEnvironment) {
		te.Logger.SetOutput(w)
	}
}

// WithTimestamps enables or disables timestamps in logs.
func WithTimestamps(enabled bool) TestEnvOption {
	return func(te *TestEnvironment) {
		te.Logger.SetTimestamps(enabled)
	}
}

// WithLogPrefix sets a custom prefix for the logger.
func WithLogPrefix(prefix string) TestEnvOption {
	return func(te *TestEnvironment) {
		te.Logger.SetPrefix(prefix)
	}
}

// NewTestEnvironment creates a new isolated test environment with automatic cleanup.
//
// Usage:
//
//	func TestMyFeature(t *testing.T) {
//	    t.Parallel() // Now safe!
//
//	    env := NewTestEnvironment(t,
//	        WithLogLevel(LogLevelDebug),
//	        WithTimestamps(false),
//	    )
//	    defer env.Cleanup()
//
//	    // Use env.Network, env.Pipeline, etc.
//	}
func NewTestEnvironment(t *testing.T, opts ...TestEnvOption) *TestEnvironment {
	t.Helper()

	// Create temporary directory
	tempDir := t.TempDir()

	// Create isolated logger with buffer
	logBuffer := &bytes.Buffer{}
	logger := NewLogger(LogLevelInfo, logBuffer)
	logger.SetTimestamps(false) // Default: no timestamps for cleaner test output

	// Create isolated components
	storage := NewMemoryStorage()

	network := NewReteNetwork(storage)
	network.SetLogger(logger)

	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(logger)

	te := &TestEnvironment{
		Network:   network,
		Storage:   storage,
		Pipeline:  pipeline,
		Logger:    logger,
		LogBuffer: logBuffer,
		TempDir:   tempDir,
		t:         t,
		cleanup:   make([]func(), 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(te)
	}

	// Re-apply logger to components in case options changed it
	network.SetLogger(te.Logger)
	pipeline.SetLogger(te.Logger)

	return te
}

// Cleanup releases all resources associated with the test environment.
// It should be called with defer after creating the environment.
func (te *TestEnvironment) Cleanup() {
	te.mu.Lock()
	defer te.mu.Unlock()

	// Run cleanup functions in reverse order (LIFO)
	for i := len(te.cleanup) - 1; i >= 0; i-- {
		te.cleanup[i]()
	}

	// TempDir is automatically cleaned up by testing.T
}

// AddCleanup registers a cleanup function to be called during Cleanup().
func (te *TestEnvironment) AddCleanup(fn func()) {
	te.mu.Lock()
	defer te.mu.Unlock()
	te.cleanup = append(te.cleanup, fn)
}

// GetLogs returns all logs captured in the buffer as a string.
func (te *TestEnvironment) GetLogs() string {
	return te.LogBuffer.String()
}

// ClearLogs clears the log buffer. Useful when testing multiple operations
// and you want to check logs for each operation separately.
func (te *TestEnvironment) ClearLogs() {
	te.LogBuffer.Reset()
}

// AssertNoErrors checks that the log buffer doesn't contain ERROR level logs.
func (te *TestEnvironment) AssertNoErrors(t *testing.T) {
	t.Helper()
	logs := te.GetLogs()
	if bytes.Contains([]byte(logs), []byte("[ERROR]")) {
		t.Errorf("Expected no errors in logs, but found:\n%s", logs)
	}
}

// AssertContainsLog checks that the log buffer contains the expected string.
func (te *TestEnvironment) AssertContainsLog(t *testing.T, expected string) {
	t.Helper()
	logs := te.GetLogs()
	if !bytes.Contains([]byte(logs), []byte(expected)) {
		t.Errorf("Expected logs to contain %q, but got:\n%s", expected, logs)
	}
}

// AssertNotContainsLog checks that the log buffer does not contain the string.
func (te *TestEnvironment) AssertNotContainsLog(t *testing.T, unexpected string) {
	t.Helper()
	logs := te.GetLogs()
	if bytes.Contains([]byte(logs), []byte(unexpected)) {
		t.Errorf("Expected logs to NOT contain %q, but got:\n%s", unexpected, logs)
	}
}

// CreateConstraintFile creates a temporary constraint file in the test directory.
func (te *TestEnvironment) CreateConstraintFile(name string, content string) string {
	te.t.Helper()

	path := filepath.Join(te.TempDir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		te.t.Fatalf("Failed to create constraint file: %v", err)
	}

	return path
}

// IngestFile is a convenience method to ingest a constraint file using the
// environment's pipeline, network, and storage.
func (te *TestEnvironment) IngestFile(filename string) (*ReteNetwork, error) {
	return te.Pipeline.IngestFile(filename, te.Network, te.Storage)
}

// IngestFileContent creates a temporary constraint file with the given content
// and ingests it, returning the results.
func (te *TestEnvironment) IngestFileContent(content string) (*ReteNetwork, error) {
	te.t.Helper()

	filename := te.CreateConstraintFile("test.constraint", content)
	return te.IngestFile(filename)
}

// RequireIngestFile ingests a file and fails the test if there's an error.
func (te *TestEnvironment) RequireIngestFile(filename string) *ReteNetwork {
	te.t.Helper()

	network, err := te.IngestFile(filename)
	if err != nil {
		te.t.Fatalf("Failed to ingest file %s: %v\nLogs:\n%s", filename, err, te.GetLogs())
	}

	return network
}

// RequireIngestFileContent creates and ingests a file, failing the test on error.
func (te *TestEnvironment) RequireIngestFileContent(content string) *ReteNetwork {
	te.t.Helper()

	network, err := te.IngestFileContent(content)
	if err != nil {
		te.t.Fatalf("Failed to ingest content: %v\nLogs:\n%s", err, te.GetLogs())
	}

	return network
}

// SubmitFact submits a fact to the network using a transaction.
func (te *TestEnvironment) SubmitFact(fact Fact) error {
	tx := te.Network.GetTransaction()
	if tx == nil {
		te.t.Fatal("No transaction available")
	}

	cmd := NewAddFactCommand(te.Storage, &fact)
	return tx.RecordAndExecute(cmd)
}

// RequireSubmitFact submits a fact and fails the test if there's an error.
func (te *TestEnvironment) RequireSubmitFact(fact Fact) {
	te.t.Helper()

	err := te.SubmitFact(fact)
	if err != nil {
		te.t.Fatalf("Failed to submit fact: %v\nLogs:\n%s", err, te.GetLogs())
	}
}

// GetFactCount returns the number of facts in storage.
func (te *TestEnvironment) GetFactCount() int {
	facts := te.Storage.GetAllFacts()
	return len(facts)
}

// GetFactsByType returns all facts of a given type.
func (te *TestEnvironment) GetFactsByType(factType string) []Fact {
	allFacts := te.Storage.GetAllFacts()
	facts := make([]Fact, 0)
	for _, factPtr := range allFacts {
		if factPtr != nil && factPtr.Type == factType {
			facts = append(facts, *factPtr)
		}
	}
	return facts
}

// SetLogLevel changes the log level for the environment.
func (te *TestEnvironment) SetLogLevel(level LogLevel) {
	te.Logger.SetLevel(level)
}

// NewSubEnvironment creates a new test environment that shares the storage
// but has its own network, pipeline, and logger. Useful for testing
// multi-network scenarios.
func (te *TestEnvironment) NewSubEnvironment(opts ...TestEnvOption) *TestEnvironment {
	te.t.Helper()

	subLogBuffer := &bytes.Buffer{}
	subLogger := NewLogger(LogLevelInfo, subLogBuffer)
	subLogger.SetTimestamps(false)

	subNetwork := NewReteNetwork(te.Storage)
	subNetwork.SetLogger(subLogger)

	subPipeline := NewConstraintPipeline()
	subPipeline.SetLogger(subLogger)

	subEnv := &TestEnvironment{
		Network:   subNetwork,
		Storage:   te.Storage, // Shared storage
		Pipeline:  subPipeline,
		Logger:    subLogger,
		LogBuffer: subLogBuffer,
		TempDir:   te.TempDir, // Shared temp dir
		t:         te.t,
		cleanup:   make([]func(), 0),
	}

	// Apply options
	for _, opt := range opts {
		opt(subEnv)
	}

	subNetwork.SetLogger(subEnv.Logger)
	subPipeline.SetLogger(subEnv.Logger)

	// Register cleanup with parent
	te.AddCleanup(subEnv.Cleanup)

	return subEnv
}
