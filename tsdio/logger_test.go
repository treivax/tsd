// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
)

// Test constants
const (
	TestMessage     = "test message"
	TestFormat      = "formatted: %s"
	TestValue       = "value"
	TestConcurrency = 100
	TestIterations  = 10
)

// TestLogger_Printf tests Printf functionality
func TestLogger_Printf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple message",
			format:   "hello",
			args:     nil,
			expected: "hello",
		},
		{
			name:     "formatted message",
			format:   "hello %s",
			args:     []interface{}{"world"},
			expected: "hello world",
		},
		{
			name:     "multiple arguments",
			format:   "%s %d %v",
			args:     []interface{}{"test", 42, true},
			expected: "test 42 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewLogger(buf)

			logger.Printf(tt.format, tt.args...)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("Printf() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestLogger_Println tests Println functionality
func TestLogger_Println(t *testing.T) {
	tests := []struct {
		name     string
		args     []interface{}
		expected string
	}{
		{
			name:     "single argument",
			args:     []interface{}{"hello"},
			expected: "hello\n",
		},
		{
			name:     "multiple arguments",
			args:     []interface{}{"hello", "world", 42},
			expected: "hello world 42\n",
		},
		{
			name:     "empty",
			args:     []interface{}{},
			expected: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewLogger(buf)

			logger.Println(tt.args...)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("Println() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestLogger_Print tests Print functionality
func TestLogger_Print(t *testing.T) {
	tests := []struct {
		name     string
		args     []interface{}
		expected string
	}{
		{
			name:     "single argument",
			args:     []interface{}{"hello"},
			expected: "hello",
		},
		{
			name:     "multiple arguments",
			args:     []interface{}{"hello", " ", "world"},
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewLogger(buf)

			logger.Print(tt.args...)

			got := buf.String()
			if got != tt.expected {
				t.Errorf("Print() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestLogger_SetOutput tests SetOutput functionality
func TestLogger_SetOutput(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	logger := NewLogger(buf1)

	// Write to buf1
	logger.Println("message1")
	if !strings.Contains(buf1.String(), "message1") {
		t.Errorf("message1 not written to buf1")
	}

	// Switch to buf2
	logger.SetOutput(buf2)
	logger.Println("message2")

	if strings.Contains(buf1.String(), "message2") {
		t.Errorf("message2 should not be in buf1")
	}
	if !strings.Contains(buf2.String(), "message2") {
		t.Errorf("message2 not written to buf2")
	}
}

// TestLogger_GetOutput tests GetOutput functionality
func TestLogger_GetOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf)

	output := logger.GetOutput()
	if output != buf {
		t.Errorf("GetOutput() returned wrong writer")
	}
}

// TestLogger_Mute tests Mute functionality
func TestLogger_Mute(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf)

	// Write before mute
	logger.Println("before")
	if !strings.Contains(buf.String(), "before") {
		t.Errorf("message before mute should be written")
	}

	// Mute and write
	logger.Mute()
	logger.Println("during")

	// Should not contain message during mute
	if strings.Contains(buf.String(), "during") {
		t.Errorf("message during mute should not be written")
	}
}

// TestLogger_Unmute tests Unmute functionality
func TestLogger_Unmute(t *testing.T) {
	logger := NewLogger(io.Discard)

	// Start muted
	logger.Println("muted")

	// Unmute (restores to dynamic os.Stdout, which we can't easily capture in tests)
	// So we verify the output is set to nil
	logger.Unmute()

	output := logger.GetOutput()
	// After Unmute, output should be os.Stdout (the default)
	if output == io.Discard {
		t.Errorf("Unmute() should restore output from Discard")
	}
}

// TestLogger_WithMutex tests WithMutex functionality
func TestLogger_WithMutex(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf)

	// Execute atomic operation using direct writes (not logger methods to avoid deadlock)
	executed := false
	logger.WithMutex(func() {
		executed = true
	})

	if !executed {
		t.Errorf("WithMutex() did not execute function")
	}
}

// TestLogger_ConcurrentWrites tests thread-safety of logger
func TestLogger_ConcurrentWrites(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf)

	var wg sync.WaitGroup
	wg.Add(TestConcurrency)

	// Launch concurrent writes
	for i := 0; i < TestConcurrency; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < TestIterations; j++ {
				logger.Printf("goroutine-%d-iteration-%d\n", id, j)
			}
		}(i)
	}

	wg.Wait()

	// Verify all messages were written
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	expectedLines := TestConcurrency * TestIterations
	if len(lines) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(lines))
	}

	// Verify no corrupted lines (each line should be complete)
	for i, line := range lines {
		if !strings.HasPrefix(line, "goroutine-") {
			t.Errorf("Line %d corrupted: %q", i, line)
		}
	}
}

// TestLogger_ConcurrentSetOutput tests concurrent SetOutput operations
func TestLogger_ConcurrentSetOutput(t *testing.T) {
	logger := NewLogger(&bytes.Buffer{})
	var wg sync.WaitGroup
	wg.Add(TestConcurrency)

	// Launch concurrent SetOutput calls
	for i := 0; i < TestConcurrency; i++ {
		go func() {
			defer wg.Done()
			buf := &bytes.Buffer{}
			logger.SetOutput(buf)
			logger.Println("test")
		}()
	}

	wg.Wait()
	// If we get here without race detector errors, test passes
}

// TestLogger_AddCaptureHook tests capture hook functionality
func TestLogger_AddCaptureHook(t *testing.T) {
	originalBuf := &bytes.Buffer{}
	logger := NewLogger(originalBuf)

	captureBuf := &bytes.Buffer{}
	captureActive := false

	// Add capture hook
	logger.AddCaptureHook(func() io.Writer {
		if captureActive {
			return captureBuf
		}
		return nil
	})

	// Write without capture
	logger.Println("original")
	if !strings.Contains(originalBuf.String(), "original") {
		t.Errorf("message should go to original buffer")
	}
	if strings.Contains(captureBuf.String(), "original") {
		t.Errorf("message should not be captured yet")
	}

	// Activate capture and write
	captureActive = true
	logger.Println("captured")

	if !strings.Contains(captureBuf.String(), "captured") {
		t.Errorf("message should be captured")
	}
	if strings.Contains(originalBuf.String(), "captured") {
		t.Errorf("message should not go to original buffer when captured")
	}
}

// TestLogger_ClearCaptureHooks tests clearing capture hooks
func TestLogger_ClearCaptureHooks(t *testing.T) {
	originalBuf := &bytes.Buffer{}
	logger := NewLogger(originalBuf)

	captureBuf := &bytes.Buffer{}

	// Add capture hook that always captures
	logger.AddCaptureHook(func() io.Writer {
		return captureBuf
	})

	// Write with hook
	logger.Println("captured")
	if !strings.Contains(captureBuf.String(), "captured") {
		t.Errorf("message should be captured")
	}

	// Clear hooks
	logger.ClearCaptureHooks()
	logger.Println("original")

	if !strings.Contains(originalBuf.String(), "original") {
		t.Errorf("message should go to original buffer after clearing hooks")
	}
}

// TestLogger_MultipleHooks tests multiple capture hooks
func TestLogger_MultipleHooks(t *testing.T) {
	originalBuf := &bytes.Buffer{}
	logger := NewLogger(originalBuf)

	captureBuf1 := &bytes.Buffer{}
	captureBuf2 := &bytes.Buffer{}

	hook1Active := false
	hook2Active := false

	// Add two hooks
	logger.AddCaptureHook(func() io.Writer {
		if hook1Active {
			return captureBuf1
		}
		return nil
	})

	logger.AddCaptureHook(func() io.Writer {
		if hook2Active {
			return captureBuf2
		}
		return nil
	})

	// Activate hook1 (first hook should win)
	hook1Active = true
	logger.Println("hook1")

	if !strings.Contains(captureBuf1.String(), "hook1") {
		t.Errorf("message should go to captureBuf1")
	}
	if strings.Contains(captureBuf2.String(), "hook1") {
		t.Errorf("message should not go to captureBuf2")
	}

	// Activate both (first hook should still win)
	hook2Active = true
	logger.Println("both")

	if !strings.Contains(captureBuf1.String(), "both") {
		t.Errorf("message should go to captureBuf1 when both active")
	}
}

// TestGlobalFunctions tests package-level functions
func TestGlobalFunctions(t *testing.T) {
	// Save original output
	originalOutput := GetOutput()
	defer SetOutput(originalOutput)

	buf := &bytes.Buffer{}
	SetOutput(buf)

	// Test Printf
	Printf("hello %s", "world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("Printf() not working")
	}

	buf.Reset()

	// Test Println
	Println("test", "message")
	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("Println() not working")
	}

	buf.Reset()

	// Test Print
	Print("no", "newline")
	if !strings.Contains(buf.String(), "nonewline") {
		t.Errorf("Print() not working")
	}
}

// TestGlobalMute tests global Mute/Unmute
func TestGlobalMute(t *testing.T) {
	// Save original output
	originalOutput := GetOutput()
	defer SetOutput(originalOutput)

	buf := &bytes.Buffer{}
	SetOutput(buf)

	// Write before mute
	Println("before")
	if !strings.Contains(buf.String(), "before") {
		t.Errorf("message before mute should be written")
	}

	// Mute
	Mute()
	Println("during")

	// Should not contain message during mute
	if strings.Contains(buf.String(), "during") {
		t.Errorf("message during mute should not be written")
	}

	// Unmute and restore
	Unmute()
	SetOutput(buf)
	Println("after")

	if !strings.Contains(buf.String(), "after") {
		t.Errorf("message after unmute should be written")
	}
}

// TestWithMutex_Global tests global WithMutex function
func TestWithMutex_Global(t *testing.T) {
	// Test that WithMutex executes the function
	executed := false
	WithMutex(func() {
		executed = true
	})

	if !executed {
		t.Errorf("WithMutex() did not execute function")
	}
}

// TestLockUnlockStdout tests LockStdout/UnlockStdout
func TestLockUnlockStdout(t *testing.T) {
	// This test verifies the functions exist and don't panic
	LockStdout()
	// Do something while locked
	UnlockStdout()
}

// TestGetGlobalLogger tests GetGlobalLogger
func TestGetGlobalLogger(t *testing.T) {
	logger := GetGlobalLogger()
	if logger == nil {
		t.Errorf("GetGlobalLogger() returned nil")
	}

	// Verify it's the same instance
	logger2 := GetGlobalLogger()
	if logger != logger2 {
		t.Errorf("GetGlobalLogger() should return same instance")
	}
}

// TestLogger_LogPrintf tests LogPrintf with timestamp
func TestLogger_LogPrintf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf)

	logger.LogPrintf("test %s", "message")

	output := buf.String()
	// Should contain the message
	if !strings.Contains(output, "test message") {
		t.Errorf("LogPrintf() should write message")
	}
	// Should have timestamp format (contains date/time)
	if !strings.Contains(output, "/") && !strings.Contains(output, ":") {
		t.Errorf("LogPrintf() should include timestamp")
	}
}

// TestGlobalLogPrintf tests global LogPrintf
func TestGlobalLogPrintf(t *testing.T) {
	// Save original output
	originalOutput := GetOutput()
	defer SetOutput(originalOutput)

	buf := &bytes.Buffer{}
	SetOutput(buf)

	LogPrintf("test %s", "log")

	output := buf.String()
	if !strings.Contains(output, "test log") {
		t.Errorf("LogPrintf() should write message")
	}
}

// TestGlobalAddCaptureHook tests global AddCaptureHook
func TestGlobalAddCaptureHook(t *testing.T) {
	// Clean up after test
	defer ClearCaptureHooks()

	// Save original output
	originalOutput := GetOutput()
	defer SetOutput(originalOutput)

	originalBuf := &bytes.Buffer{}
	SetOutput(originalBuf)

	captureBuf := &bytes.Buffer{}
	captureActive := false

	AddCaptureHook(func() io.Writer {
		if captureActive {
			return captureBuf
		}
		return nil
	})

	// Write without capture
	Println("original")
	if !strings.Contains(originalBuf.String(), "original") {
		t.Errorf("message should go to original buffer")
	}

	// Activate capture
	captureActive = true
	Println("captured")

	if !strings.Contains(captureBuf.String(), "captured") {
		t.Errorf("message should be captured")
	}
}

// TestGlobalClearCaptureHooks tests global ClearCaptureHooks
func TestGlobalClearCaptureHooks(t *testing.T) {
	// Save original output
	originalOutput := GetOutput()
	defer SetOutput(originalOutput)

	originalBuf := &bytes.Buffer{}
	SetOutput(originalBuf)

	captureBuf := &bytes.Buffer{}

	// Add hook that always captures
	AddCaptureHook(func() io.Writer {
		return captureBuf
	})

	// Clear hooks
	ClearCaptureHooks()

	// Write should go to original buffer
	Println("after-clear")
	if !strings.Contains(originalBuf.String(), "after-clear") {
		t.Errorf("message should go to original buffer after clearing hooks")
	}
}

// TestNewLogger tests NewLogger constructor
func TestNewLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf)

	if logger == nil {
		t.Errorf("NewLogger() returned nil")
	}

	// Verify it writes to the buffer
	logger.Println("test")
	if !strings.Contains(buf.String(), "test") {
		t.Errorf("NewLogger() created logger that doesn't write to buffer")
	}
}

// TestLogger_NilOutput tests logger with nil output (should use os.Stdout)
func TestLogger_NilOutput(t *testing.T) {
	logger := &Logger{
		output: nil,
	}

	// This should not panic
	output := logger.resolveOutput()
	if output == nil {
		t.Errorf("resolveOutput() should not return nil")
	}
}
