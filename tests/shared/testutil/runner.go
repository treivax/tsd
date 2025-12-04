// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tsdio"
)

// TSDFixture represents a .tsd test file
type TSDFixture struct {
	Name        string
	Path        string
	Category    string
	ShouldError bool
}

// TSDResult contains execution results
type TSDResult struct {
	TypeNodes     int
	TerminalNodes int
	Facts         int
	Activations   int
	Error         error
	Output        string
	Duration      time.Duration
}

// ExecutionOptions configures test execution
type ExecutionOptions struct {
	ExpectError     bool
	MinActivations  int
	MaxActivations  int
	ValidateNetwork bool
	CaptureOutput   bool
	Timeout         time.Duration
}

// DefaultExecutionOptions returns default options
func DefaultExecutionOptions() *ExecutionOptions {
	return &ExecutionOptions{
		ExpectError:     false,
		MinActivations:  -1, // No minimum
		MaxActivations:  -1, // No maximum
		ValidateNetwork: true,
		CaptureOutput:   true,
		Timeout:         30 * time.Second,
	}
}

// ExecuteTSDFile executes a .tsd file and returns results
func ExecuteTSDFile(t *testing.T, path string) *TSDResult {
	t.Helper()
	return ExecuteTSDFileWithOptions(t, path, DefaultExecutionOptions())
}

// ExecuteTSDFileWithOptions executes with custom configuration
func ExecuteTSDFileWithOptions(t *testing.T, path string, opts *ExecutionOptions) *TSDResult {
	t.Helper()

	if opts == nil {
		opts = DefaultExecutionOptions()
	}

	result := &TSDResult{}
	startTime := time.Now()

	// Create pipeline and storage
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	// Variable to hold network reference
	var network *rete.ReteNetwork
	var err error

	// Capture stdout if requested
	var capturedOutput string
	if opts.CaptureOutput {
		capturedOutput = captureOutput(func() {
			network, err = pipeline.IngestFile(path, nil, storage)
			result.Error = err
		})
		result.Output = capturedOutput
	} else {
		network, err = pipeline.IngestFile(path, nil, storage)
		result.Error = err
	}

	// Process results AFTER capture is complete to avoid race conditions
	if network != nil {
		result.TypeNodes = len(network.TypeNodes)
		result.TerminalNodes = len(network.TerminalNodes)

		// Count activations
		for _, terminal := range network.TerminalNodes {
			if terminal.Memory != nil && terminal.Memory.Tokens != nil {
				result.Activations += len(terminal.Memory.Tokens)
			}
		}
	}

	// Get facts from storage AFTER everything is done
	facts := storage.GetAllFacts()
	result.Facts = len(facts)

	result.Duration = time.Since(startTime)

	// Detect injection errors in output
	hasInjectionErrors := strings.Contains(result.Output, "⚠️ Erreur injection fait")

	// If expecting error, validate
	if opts.ExpectError {
		if result.Error == nil && !hasInjectionErrors {
			t.Errorf("Expected error for %s but got none", path)
		}
	} else {
		if result.Error != nil {
			t.Errorf("Unexpected error for %s: %v", path, result.Error)
		}
		if hasInjectionErrors {
			t.Errorf("Unexpected injection errors for %s", path)
		}
	}

	// Validate network if requested
	if opts.ValidateNetwork && result.Error == nil && !opts.ExpectError {
		if result.TypeNodes == 0 {
			t.Logf("Warning: No type nodes created for %s", path)
		}
		if result.TerminalNodes == 0 {
			t.Logf("Warning: No terminal nodes created for %s", path)
		}
	}

	// Validate activation constraints
	if opts.MinActivations >= 0 && result.Activations < opts.MinActivations {
		t.Errorf("Expected at least %d activations for %s, got %d",
			opts.MinActivations, path, result.Activations)
	}
	if opts.MaxActivations >= 0 && result.Activations > opts.MaxActivations {
		t.Errorf("Expected at most %d activations for %s, got %d",
			opts.MaxActivations, path, result.Activations)
	}

	return result
}

// RunTSDFile is a convenience wrapper that fails the test on error
func RunTSDFile(t *testing.T, path string) *TSDResult {
	t.Helper()
	result := ExecuteTSDFile(t, path)
	if result.Error != nil {
		t.Fatalf("Failed to execute %s: %v", path, result.Error)
	}
	return result
}

// captureOutput captures stdout during function execution
// This function is thread-safe by using tsdio's stdout mutex to prevent races on os.Stdout
func captureOutput(fn func()) string {
	// Lock only during os.Stdout modifications, not during fn() execution
	tsdio.LockStdout()
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	tsdio.UnlockStdout()

	// Channel to read the output
	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputChan <- buf.String()
	}()

	// Execute function (without holding mutex to avoid deadlock)
	fn()

	// Lock again to restore stdout
	tsdio.LockStdout()
	w.Close()
	os.Stdout = oldStdout
	tsdio.UnlockStdout()

	// Read captured output
	return <-outputChan
}

// AssertTSDResult validates a TSD execution result against expected values
func AssertTSDResult(t *testing.T, result *TSDResult, expected *TSDResult) {
	t.Helper()

	if expected.TypeNodes >= 0 && result.TypeNodes != expected.TypeNodes {
		t.Errorf("Expected %d type nodes, got %d", expected.TypeNodes, result.TypeNodes)
	}

	if expected.TerminalNodes >= 0 && result.TerminalNodes != expected.TerminalNodes {
		t.Errorf("Expected %d terminal nodes, got %d", expected.TerminalNodes, result.TerminalNodes)
	}

	if expected.Facts >= 0 && result.Facts != expected.Facts {
		t.Errorf("Expected %d facts, got %d", expected.Facts, result.Facts)
	}

	if expected.Activations >= 0 && result.Activations != expected.Activations {
		t.Errorf("Expected %d activations, got %d", expected.Activations, result.Activations)
	}

	if expected.Error != nil && result.Error == nil {
		t.Error("Expected an error but got none")
	}

	if expected.Error == nil && result.Error != nil {
		t.Errorf("Unexpected error: %v", result.Error)
	}

	if expected.Output != "" && !strings.Contains(result.Output, expected.Output) {
		t.Errorf("Expected output to contain %q, got: %s", expected.Output, result.Output)
	}
}

// CompareResults compares two TSD execution results
func CompareResults(t *testing.T, result1, result2 *TSDResult, label1, label2 string) {
	t.Helper()

	if result1.TypeNodes != result2.TypeNodes {
		t.Errorf("%s has %d type nodes, %s has %d",
			label1, result1.TypeNodes, label2, result2.TypeNodes)
	}

	if result1.TerminalNodes != result2.TerminalNodes {
		t.Errorf("%s has %d terminal nodes, %s has %d",
			label1, result1.TerminalNodes, label2, result2.TerminalNodes)
	}

	if result1.Facts != result2.Facts {
		t.Errorf("%s has %d facts, %s has %d",
			label1, result1.Facts, label2, result2.Facts)
	}

	if result1.Activations != result2.Activations {
		t.Errorf("%s has %d activations, %s has %d",
			label1, result1.Activations, label2, result2.Activations)
	}

	if (result1.Error != nil) != (result2.Error != nil) {
		t.Errorf("%s error: %v, %s error: %v",
			label1, result1.Error, label2, result2.Error)
	}
}
