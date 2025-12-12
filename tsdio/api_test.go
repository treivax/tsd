// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
	"testing"
	"time"
)

// Test constants
const (
	TestSource        = "type Person : <id: string>"
	TestSourceName    = "test.tsd"
	TestErrorMessage  = "test error"
	TestVersion       = "1.0.0"
	TestBuildTime     = "2025-01-01"
	TestGitCommit     = "abc123"
	TestGoVersion     = "go1.21"
	TestActionName    = "testAction"
	TestFactID        = "fact1"
	TestFactType      = "Person"
	TestUsername      = "testuser"
	TestExecutionTime = int64(100)
	TestFactsCount    = 5
	TestActivations   = 3
	TestBindingsCount = 2
	TestArgumentPos   = 0
	TestArgValue      = "test-value"
	TestArgumentType  = "string"
	TestUptimeSeconds = int64(3600)
)

// TestNewExecuteRequest tests NewExecuteRequest constructor
func TestNewExecuteRequest(t *testing.T) {
	req := NewExecuteRequest(TestSource)

	if req == nil {
		t.Fatal("NewExecuteRequest() returned nil")
	}

	if req.Source != TestSource {
		t.Errorf("Source = %q, want %q", req.Source, TestSource)
	}

	if req.SourceName != "<request>" {
		t.Errorf("SourceName = %q, want %q", req.SourceName, "<request>")
	}

	if req.Verbose {
		t.Errorf("Verbose = true, want false")
	}
}

// TestExecuteRequest_Fields tests ExecuteRequest field initialization
func TestExecuteRequest_Fields(t *testing.T) {
	tests := []struct {
		name        string
		req         *ExecuteRequest
		wantSource  string
		wantName    string
		wantVerbose bool
	}{
		{
			name: "custom fields",
			req: &ExecuteRequest{
				Source:     TestSource,
				SourceName: TestSourceName,
				Verbose:    true,
			},
			wantSource:  TestSource,
			wantName:    TestSourceName,
			wantVerbose: true,
		},
		{
			name: "empty source name",
			req: &ExecuteRequest{
				Source:     TestSource,
				SourceName: "",
				Verbose:    false,
			},
			wantSource:  TestSource,
			wantName:    "",
			wantVerbose: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Source != tt.wantSource {
				t.Errorf("Source = %q, want %q", tt.req.Source, tt.wantSource)
			}
			if tt.req.SourceName != tt.wantName {
				t.Errorf("SourceName = %q, want %q", tt.req.SourceName, tt.wantName)
			}
			if tt.req.Verbose != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", tt.req.Verbose, tt.wantVerbose)
			}
		})
	}
}

// TestNewSuccessResponse tests NewSuccessResponse constructor
func TestNewSuccessResponse(t *testing.T) {
	results := &ExecutionResults{
		FactsCount:       TestFactsCount,
		ActivationsCount: TestActivations,
	}

	resp := NewSuccessResponse(results, TestExecutionTime)

	if resp == nil {
		t.Fatal("NewSuccessResponse() returned nil")
	}

	if !resp.Success {
		t.Errorf("Success = false, want true")
	}

	if resp.Error != "" {
		t.Errorf("Error = %q, want empty", resp.Error)
	}

	if resp.ErrorType != "" {
		t.Errorf("ErrorType = %q, want empty", resp.ErrorType)
	}

	if resp.Results != results {
		t.Errorf("Results not set correctly")
	}

	if resp.ExecutionTimeMs != TestExecutionTime {
		t.Errorf("ExecutionTimeMs = %d, want %d", resp.ExecutionTimeMs, TestExecutionTime)
	}
}

// TestNewErrorResponse tests NewErrorResponse constructor
func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name          string
		errorType     string
		errorMsg      string
		executionTime int64
	}{
		{
			name:          "parsing error",
			errorType:     ErrorTypeParsingError,
			errorMsg:      "syntax error",
			executionTime: TestExecutionTime,
		},
		{
			name:          "validation error",
			errorType:     ErrorTypeValidationError,
			errorMsg:      "type mismatch",
			executionTime: TestExecutionTime,
		},
		{
			name:          "execution error",
			errorType:     ErrorTypeExecutionError,
			errorMsg:      "runtime error",
			executionTime: TestExecutionTime,
		},
		{
			name:          "server error",
			errorType:     ErrorTypeServerError,
			errorMsg:      "internal error",
			executionTime: TestExecutionTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewErrorResponse(tt.errorType, tt.errorMsg, tt.executionTime)

			if resp == nil {
				t.Fatal("NewErrorResponse() returned nil")
			}

			if resp.Success {
				t.Errorf("Success = true, want false")
			}

			if resp.Error != tt.errorMsg {
				t.Errorf("Error = %q, want %q", resp.Error, tt.errorMsg)
			}

			if resp.ErrorType != tt.errorType {
				t.Errorf("ErrorType = %q, want %q", resp.ErrorType, tt.errorType)
			}

			if resp.Results != nil {
				t.Errorf("Results should be nil for error response")
			}

			if resp.ExecutionTimeMs != tt.executionTime {
				t.Errorf("ExecutionTimeMs = %d, want %d", resp.ExecutionTimeMs, tt.executionTime)
			}
		})
	}
}

// TestErrorTypeConstants tests error type constants
func TestErrorTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "parsing error",
			constant: ErrorTypeParsingError,
			expected: "parsing_error",
		},
		{
			name:     "validation error",
			constant: ErrorTypeValidationError,
			expected: "validation_error",
		},
		{
			name:     "execution error",
			constant: ErrorTypeExecutionError,
			expected: "execution_error",
		},
		{
			name:     "server error",
			constant: ErrorTypeServerError,
			expected: "server_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("constant = %q, want %q", tt.constant, tt.expected)
			}
		})
	}
}

// TestExecutionResults tests ExecutionResults structure
func TestExecutionResults(t *testing.T) {
	activation := Activation{
		ActionName:    TestActionName,
		Arguments:     []ArgumentValue{{Position: 0, Value: "test", Type: "string"}},
		BindingsCount: TestBindingsCount,
	}

	results := &ExecutionResults{
		FactsCount:       TestFactsCount,
		ActivationsCount: TestActivations,
		Activations:      []Activation{activation},
	}

	if results.FactsCount != TestFactsCount {
		t.Errorf("FactsCount = %d, want %d", results.FactsCount, TestFactsCount)
	}

	if results.ActivationsCount != TestActivations {
		t.Errorf("ActivationsCount = %d, want %d", results.ActivationsCount, TestActivations)
	}

	if len(results.Activations) != 1 {
		t.Errorf("len(Activations) = %d, want 1", len(results.Activations))
	}

	if results.Activations[0].ActionName != TestActionName {
		t.Errorf("Activations[0].ActionName = %q, want %q", results.Activations[0].ActionName, TestActionName)
	}
}

// TestActivation tests Activation structure
func TestActivation(t *testing.T) {
	fact := Fact{
		ID:     TestFactID,
		Type:   TestFactType,
		Fields: map[string]interface{}{"name": "Alice"},
	}

	arg := ArgumentValue{
		Position: TestArgumentPos,
		Value:    TestArgValue,
		Type:     TestArgumentType,
	}

	activation := Activation{
		ActionName:      TestActionName,
		Arguments:       []ArgumentValue{arg},
		TriggeringFacts: []Fact{fact},
		BindingsCount:   TestBindingsCount,
	}

	if activation.ActionName != TestActionName {
		t.Errorf("ActionName = %q, want %q", activation.ActionName, TestActionName)
	}

	if len(activation.Arguments) != 1 {
		t.Errorf("len(Arguments) = %d, want 1", len(activation.Arguments))
	}

	if activation.Arguments[0].Position != TestArgumentPos {
		t.Errorf("Arguments[0].Position = %d, want %d", activation.Arguments[0].Position, TestArgumentPos)
	}

	if activation.Arguments[0].Value != TestArgValue {
		t.Errorf("Arguments[0].Value = %v, want %v", activation.Arguments[0].Value, TestArgValue)
	}

	if activation.Arguments[0].Type != TestArgumentType {
		t.Errorf("Arguments[0].Type = %q, want %q", activation.Arguments[0].Type, TestArgumentType)
	}

	if len(activation.TriggeringFacts) != 1 {
		t.Errorf("len(TriggeringFacts) = %d, want 1", len(activation.TriggeringFacts))
	}

	if activation.TriggeringFacts[0].ID != TestFactID {
		t.Errorf("TriggeringFacts[0].ID = %q, want %q", activation.TriggeringFacts[0].ID, TestFactID)
	}

	if activation.BindingsCount != TestBindingsCount {
		t.Errorf("BindingsCount = %d, want %d", activation.BindingsCount, TestBindingsCount)
	}
}

// TestArgumentValue tests ArgumentValue structure
func TestArgumentValue(t *testing.T) {
	tests := []struct {
		name     string
		arg      ArgumentValue
		wantPos  int
		wantVal  interface{}
		wantType string
	}{
		{
			name:     "string argument",
			arg:      ArgumentValue{Position: 0, Value: "test", Type: "string"},
			wantPos:  0,
			wantVal:  "test",
			wantType: "string",
		},
		{
			name:     "number argument",
			arg:      ArgumentValue{Position: 1, Value: 42, Type: "number"},
			wantPos:  1,
			wantVal:  42,
			wantType: "number",
		},
		{
			name:     "boolean argument",
			arg:      ArgumentValue{Position: 2, Value: true, Type: "bool"},
			wantPos:  2,
			wantVal:  true,
			wantType: "bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.arg.Position != tt.wantPos {
				t.Errorf("Position = %d, want %d", tt.arg.Position, tt.wantPos)
			}
			if tt.arg.Value != tt.wantVal {
				t.Errorf("Value = %v, want %v", tt.arg.Value, tt.wantVal)
			}
			if tt.arg.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", tt.arg.Type, tt.wantType)
			}
		})
	}
}

// TestFact tests Fact structure
func TestFact(t *testing.T) {
	attrs := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}

	fact := Fact{
		ID:     TestFactID,
		Type:   TestFactType,
		Fields: attrs,
	}

	if fact.ID != TestFactID {
		t.Errorf("ID = %q, want %q", fact.ID, TestFactID)
	}

	if fact.Type != TestFactType {
		t.Errorf("Type = %q, want %q", fact.Type, TestFactType)
	}

	if len(fact.Fields) != 2 {
		t.Errorf("len(Attributes) = %d, want 2", len(fact.Fields))
	}

	if fact.Fields["name"] != "Alice" {
		t.Errorf("Attributes[name] = %v, want Alice", fact.Fields["name"])
	}

	if fact.Fields["age"] != 30 {
		t.Errorf("Attributes[age] = %v, want 30", fact.Fields["age"])
	}
}

// TestHealthResponse tests HealthResponse structure
func TestHealthResponse(t *testing.T) {
	now := time.Now()
	health := HealthResponse{
		Status:        "ok",
		Version:       TestVersion,
		UptimeSeconds: TestUptimeSeconds,
		Timestamp:     now,
	}

	if health.Status != "ok" {
		t.Errorf("Status = %q, want ok", health.Status)
	}

	if health.Version != TestVersion {
		t.Errorf("Version = %q, want %q", health.Version, TestVersion)
	}

	if health.UptimeSeconds != TestUptimeSeconds {
		t.Errorf("UptimeSeconds = %d, want %d", health.UptimeSeconds, TestUptimeSeconds)
	}

	if health.Timestamp != now {
		t.Errorf("Timestamp mismatch")
	}
}

// TestVersionResponse tests VersionResponse structure
func TestVersionResponse(t *testing.T) {
	version := VersionResponse{
		Version:   TestVersion,
		BuildTime: TestBuildTime,
		GitCommit: TestGitCommit,
		GoVersion: TestGoVersion,
	}

	if version.Version != TestVersion {
		t.Errorf("Version = %q, want %q", version.Version, TestVersion)
	}

	if version.BuildTime != TestBuildTime {
		t.Errorf("BuildTime = %q, want %q", version.BuildTime, TestBuildTime)
	}

	if version.GitCommit != TestGitCommit {
		t.Errorf("GitCommit = %q, want %q", version.GitCommit, TestGitCommit)
	}

	if version.GoVersion != TestGoVersion {
		t.Errorf("GoVersion = %q, want %q", version.GoVersion, TestGoVersion)
	}
}

// TestExecuteResponse_SuccessCase tests success response scenario
func TestExecuteResponse_SuccessCase(t *testing.T) {
	results := &ExecutionResults{
		FactsCount:       TestFactsCount,
		ActivationsCount: TestActivations,
		Activations:      []Activation{},
	}

	resp := &ExecuteResponse{
		Success:         true,
		Results:         results,
		ExecutionTimeMs: TestExecutionTime,
	}

	if !resp.Success {
		t.Errorf("Success = false, want true")
	}

	if resp.Error != "" {
		t.Errorf("Error should be empty for success response")
	}

	if resp.ErrorType != "" {
		t.Errorf("ErrorType should be empty for success response")
	}

	if resp.Results == nil {
		t.Errorf("Results should not be nil for success response")
	}

	if resp.ExecutionTimeMs != TestExecutionTime {
		t.Errorf("ExecutionTimeMs = %d, want %d", resp.ExecutionTimeMs, TestExecutionTime)
	}
}

// TestExecuteResponse_ErrorCase tests error response scenario
func TestExecuteResponse_ErrorCase(t *testing.T) {
	resp := &ExecuteResponse{
		Success:         false,
		Error:           TestErrorMessage,
		ErrorType:       ErrorTypeParsingError,
		ExecutionTimeMs: TestExecutionTime,
	}

	if resp.Success {
		t.Errorf("Success = true, want false")
	}

	if resp.Error != TestErrorMessage {
		t.Errorf("Error = %q, want %q", resp.Error, TestErrorMessage)
	}

	if resp.ErrorType != ErrorTypeParsingError {
		t.Errorf("ErrorType = %q, want %q", resp.ErrorType, ErrorTypeParsingError)
	}

	if resp.Results != nil {
		t.Errorf("Results should be nil for error response")
	}

	if resp.ExecutionTimeMs != TestExecutionTime {
		t.Errorf("ExecutionTimeMs = %d, want %d", resp.ExecutionTimeMs, TestExecutionTime)
	}
}

// TestExecutionResults_EmptyActivations tests ExecutionResults with no activations
func TestExecutionResults_EmptyActivations(t *testing.T) {
	results := &ExecutionResults{
		FactsCount:       TestFactsCount,
		ActivationsCount: 0,
		Activations:      []Activation{},
	}

	if results.FactsCount != TestFactsCount {
		t.Errorf("FactsCount = %d, want %d", results.FactsCount, TestFactsCount)
	}

	if results.ActivationsCount != 0 {
		t.Errorf("ActivationsCount = %d, want 0", results.ActivationsCount)
	}

	if len(results.Activations) != 0 {
		t.Errorf("len(Activations) = %d, want 0", len(results.Activations))
	}
}

// TestFact_EmptyAttributes tests Fact with empty attributes
func TestFact_EmptyAttributes(t *testing.T) {
	fact := Fact{
		ID:     TestFactID,
		Type:   TestFactType,
		Fields: map[string]interface{}{},
	}

	if fact.ID != TestFactID {
		t.Errorf("ID = %q, want %q", fact.ID, TestFactID)
	}

	if fact.Type != TestFactType {
		t.Errorf("Type = %q, want %q", fact.Type, TestFactType)
	}

	if len(fact.Fields) != 0 {
		t.Errorf("len(Attributes) = %d, want 0", len(fact.Fields))
	}
}

// TestActivation_EmptyArrays tests Activation with empty arrays
func TestActivation_EmptyArrays(t *testing.T) {
	activation := Activation{
		ActionName:      TestActionName,
		Arguments:       []ArgumentValue{},
		TriggeringFacts: []Fact{},
		BindingsCount:   0,
	}

	if activation.ActionName != TestActionName {
		t.Errorf("ActionName = %q, want %q", activation.ActionName, TestActionName)
	}

	if len(activation.Arguments) != 0 {
		t.Errorf("len(Arguments) = %d, want 0", len(activation.Arguments))
	}

	if len(activation.TriggeringFacts) != 0 {
		t.Errorf("len(TriggeringFacts) = %d, want 0", len(activation.TriggeringFacts))
	}

	if activation.BindingsCount != 0 {
		t.Errorf("BindingsCount = %d, want 0", activation.BindingsCount)
	}
}
