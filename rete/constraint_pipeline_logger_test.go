// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestConstraintPipeline_SilentLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelSilent, &buf)
	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(logger)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SetLogger(logger)
	// Create a simple test constraint file content
	constraintContent := `type Person(name: string, age: number)
action print(message: string)
rule Adults : {p: Person} / p.age >= 18 ==> print(p.name)
`
	// Write to temp file
	tmpFile := createTempLoggerTestFile(t, constraintContent)
	// Ingest with silent logging
	_, _, err := pipeline.IngestFile(tmpFile, network, storage)
	require.NoError(t, err)
	// Assert no output was produced
	output := buf.String()
	assert.Empty(t, output, "Silent mode should produce no log output")
}
func TestConstraintPipeline_DebugLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, &buf)
	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(logger)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SetLogger(logger)
	constraintContent := `type Person(name: string, age: number)
action print(message: string)
rule TestRule : {p: Person} / p.age > 20 ==> print(p.name)
`
	tmpFile := createTempLoggerTestFile(t, constraintContent)
	_, _, err := pipeline.IngestFile(tmpFile, network, storage)
	require.NoError(t, err)
	output := buf.String()
	// Debug mode should contain detailed traces
	assert.NotEmpty(t, output, "Debug mode should produce output")
	assert.Contains(t, output, "[DEBUG]", "Should contain DEBUG level logs")
	// Should log various stages of processing
	// (exact log messages may vary, but we expect some processing details)
	logLines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Greater(t, len(logLines), 3, "Debug mode should produce multiple log lines")
}
func TestConstraintPipeline_InfoLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(logger)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SetLogger(logger)
	constraintContent := `type Employee(id: number, name: string, salary: number)
action print(message: string)
rule HighEarners : {e: Employee} / e.salary > 100000 ==> print(e.name)
`
	tmpFile := createTempLoggerTestFile(t, constraintContent)
	_, _, err := pipeline.IngestFile(tmpFile, network, storage)
	require.NoError(t, err)
	output := buf.String()
	// Info mode should show milestones but not debug details
	assert.NotEmpty(t, output, "Info mode should produce output")
	assert.Contains(t, output, "[INFO]", "Should contain INFO level logs")
	assert.NotContains(t, output, "[DEBUG]", "Should not contain DEBUG level logs")
}
func TestConstraintPipeline_SetGetLogger(t *testing.T) {
	var buf bytes.Buffer
	customLogger := NewLogger(LogLevelInfo, &buf)
	customLogger.SetPrefix("[CUSTOM]")
	customLogger.SetTimestamps(false)
	pipeline := NewConstraintPipeline()
	// Initially, pipeline has default logger
	defaultLogger := pipeline.GetLogger()
	require.NotNil(t, defaultLogger)
	// Set custom logger
	pipeline.SetLogger(customLogger)
	// GetLogger should return the custom logger
	retrievedLogger := pipeline.GetLogger()
	assert.Equal(t, customLogger, retrievedLogger, "GetLogger should return the set logger")
	// Verify logs go to custom buffer with custom prefix
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SetLogger(customLogger)
	constraintContent := `type Product(name: string, price: number)
action print(message: string)
rule ExpensiveProducts : {p: Product} / p.price > 50 ==> print(p.name)
`
	tmpFile := createTempLoggerTestFile(t, constraintContent)
	_, _, err := pipeline.IngestFile(tmpFile, network, storage)
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "[CUSTOM]", "Logs should use custom prefix")
}
func TestConstraintPipeline_LazyLoggerInit(t *testing.T) {
	// Create pipeline without explicitly setting logger
	pipeline := &ConstraintPipeline{}
	// GetLogger should initialize a default logger
	logger1 := pipeline.GetLogger()
	require.NotNil(t, logger1, "GetLogger should initialize logger if nil")
	// Subsequent calls should return same instance
	logger2 := pipeline.GetLogger()
	assert.Equal(t, logger1, logger2, "GetLogger should return same instance")
	// Logger should be functional
	var buf bytes.Buffer
	logger1.SetOutput(&buf)
	logger1.SetTimestamps(false)
	logger1.Info("test message")
	output := buf.String()
	assert.Contains(t, output, "test message", "Lazy-initialized logger should be functional")
}
func TestConstraintPipeline_LoggerNilSafety(t *testing.T) {
	pipeline := NewConstraintPipeline()
	// SetLogger with nil should be a no-op (doesn't panic)
	pipeline.SetLogger(nil)
	// GetLogger should still return valid logger
	logger := pipeline.GetLogger()
	require.NotNil(t, logger, "GetLogger should never return nil")
}
func TestConstraintPipeline_MultipleFilesLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	logger.SetTimestamps(false)
	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(logger)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SetLogger(logger)
	// First file
	file1Content := `type Person(name: string, age: number)
action print(message: string)
rule Adults : {p: Person} / p.age >= 18 ==> print(p.name)
`
	tmpFile1 := createTempLoggerTestFile(t, file1Content)
	_, _, err := pipeline.IngestFile(tmpFile1, network, storage)
	require.NoError(t, err)
	// Second file (extending network)
	file2Content := `type Employee(id: number, name: string, salary: number)
rule AllEmployees : {e: Employee} / e.salary > 0 ==> print(e.name)
`
	tmpFile2 := createTempLoggerTestFile(t, file2Content)
	_, _, err = pipeline.IngestFile(tmpFile2, network, storage)
	require.NoError(t, err)
	output := buf.String()
	// Both ingestions should be logged
	assert.Contains(t, output, "[INFO]", "Should contain info logs")
	// Should have multiple log entries (at least one per file)
	logLines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Greater(t, len(logLines), 1, "Should have logs from both ingestions")
}
func TestConstraintPipeline_LoggerIsolation(t *testing.T) {
	// Test that different pipelines can have different loggers
	var buf1 bytes.Buffer
	logger1 := NewLogger(LogLevelInfo, &buf1)
	logger1.SetPrefix("[PIPELINE1]")
	logger1.SetTimestamps(false)
	var buf2 bytes.Buffer
	logger2 := NewLogger(LogLevelInfo, &buf2)
	logger2.SetPrefix("[PIPELINE2]")
	logger2.SetTimestamps(false)
	pipeline1 := NewConstraintPipeline()
	pipeline1.SetLogger(logger1)
	pipeline2 := NewConstraintPipeline()
	pipeline2.SetLogger(logger2)
	storage := NewMemoryStorage()
	network1 := NewReteNetwork(storage)
	network1.SetLogger(logger1)
	network2 := NewReteNetwork(storage)
	network2.SetLogger(logger2)
	constraintContent := `type Item(id: number)
action print(message: string)
rule AllItems : {i: Item} / i.id > 0 ==> print("Item found")
`
	tmpFile := createTempLoggerTestFile(t, constraintContent)
	// Process with pipeline1
	_, _, err := pipeline1.IngestFile(tmpFile, network1, storage)
	require.NoError(t, err)
	// Process with pipeline2
	_, _, err = pipeline2.IngestFile(tmpFile, network2, storage)
	require.NoError(t, err)
	output1 := buf1.String()
	output2 := buf2.String()
	// Each buffer should only contain its own prefix
	assert.Contains(t, output1, "[PIPELINE1]", "Buffer 1 should have PIPELINE1 prefix")
	assert.NotContains(t, output1, "[PIPELINE2]", "Buffer 1 should not have PIPELINE2 prefix")
	assert.Contains(t, output2, "[PIPELINE2]", "Buffer 2 should have PIPELINE2 prefix")
	assert.NotContains(t, output2, "[PIPELINE1]", "Buffer 2 should not have PIPELINE1 prefix")
}
func TestConstraintPipeline_ErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelError, &buf)
	logger.SetTimestamps(false)
	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(logger)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Invalid constraint file (syntax error)
	invalidContent := `type Person(name: string, age: number)
action print(message: string)
rule InvalidRule : {p: Person} / p.age >= 18 AND p.name == ==> print("Missing operand")
`
	tmpFile := createTempLoggerTestFile(t, invalidContent)
	// This should produce an error
	_, _, err := pipeline.IngestFile(tmpFile, network, storage)
	// We expect an error
	assert.Error(t, err, "Invalid syntax should produce error")
	output := buf.String()
	// At Error level, only errors should be logged (if any error logs exist)
	// Note: depending on implementation, errors might be returned without logging
	// or they might be logged. We just verify no DEBUG/INFO/WARN appears.
	assert.NotContains(t, output, "[DEBUG]", "Error level should not show DEBUG")
	assert.NotContains(t, output, "[INFO]", "Error level should not show INFO")
	assert.NotContains(t, output, "[WARN]", "Error level should not show WARN")
}
func TestConstraintPipeline_ContextualLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	logger.SetTimestamps(false)
	// Create logger with context
	contextLogger := logger.WithContext("TestContext")
	pipeline := NewConstraintPipeline()
	pipeline.SetLogger(contextLogger)
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.SetLogger(contextLogger)
	constraintContent := `type Event(id: number)
action print(message: string)
rule AllEvents : {e: Event} / e.id > 0 ==> print("Event found")
`
	tmpFile := createTempLoggerTestFile(t, constraintContent)
	_, _, err := pipeline.IngestFile(tmpFile, network, storage)
	require.NoError(t, err)
	output := buf.String()
	// Should contain the context in logs
	assert.Contains(t, output, "[TestContext]", "Logs should include context")
}

// Helper function to create temporary constraint file for logger tests
func createTempLoggerTestFile(t *testing.T, content string) string {
	t.Helper()
	tmpFile := t.TempDir() + "/test.constraint"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err, "Failed to create temp file")
	return tmpFile
}
