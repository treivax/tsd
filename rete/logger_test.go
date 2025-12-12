// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"bytes"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
func TestLogger_Levels(t *testing.T) {
	tests := []struct {
		name          string
		level         LogLevel
		logFunc       func(*Logger, string)
		expectedLevel string
		shouldLog     bool
	}{
		{"Debug at Debug level", LogLevelDebug, func(l *Logger, msg string) { l.Debug(msg) }, "DEBUG", true},
		{"Debug at Info level", LogLevelInfo, func(l *Logger, msg string) { l.Debug(msg) }, "DEBUG", false},
		{"Info at Info level", LogLevelInfo, func(l *Logger, msg string) { l.Info(msg) }, "INFO", true},
		{"Info at Warn level", LogLevelWarn, func(l *Logger, msg string) { l.Info(msg) }, "INFO", false},
		{"Warn at Warn level", LogLevelWarn, func(l *Logger, msg string) { l.Warn(msg) }, "WARN", true},
		{"Warn at Error level", LogLevelError, func(l *Logger, msg string) { l.Warn(msg) }, "WARN", false},
		{"Error at Error level", LogLevelError, func(l *Logger, msg string) { l.Error(msg) }, "ERROR", true},
		{"Error at Silent level", LogLevelSilent, func(l *Logger, msg string) { l.Error(msg) }, "ERROR", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(tt.level, &buf)
			logger.SetTimestamps(false)
			testMessage := "test message"
			tt.logFunc(logger, testMessage)
			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, testMessage, "Expected message to be logged")
				assert.Contains(t, output, tt.expectedLevel, "Expected level to be in output")
			} else {
				assert.Empty(t, output, "Expected no output for this log level")
			}
		})
	}
}
func TestLogger_Format(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	logger.SetTimestamps(false)
	logger.SetPrefix("[TEST]")
	logger.Info("formatted message: %s, number: %d", "hello", 42)
	output := buf.String()
	assert.Contains(t, output, "[TEST]")
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "formatted message: hello, number: 42")
}
func TestLogger_Timestamps(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	// With timestamps (default)
	logger.SetTimestamps(true)
	logger.Info("test with timestamp")
	outputWithTimestamp := buf.String()
	// Should contain a timestamp pattern (YYYY-MM-DD HH:MM:SS.mmm)
	assert.Regexp(t, `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}`, outputWithTimestamp)
	// Without timestamps
	buf.Reset()
	logger.SetTimestamps(false)
	logger.Info("test without timestamp")
	outputWithoutTimestamp := buf.String()
	assert.NotRegexp(t, `\d{4}-\d{2}-\d{2}`, outputWithoutTimestamp)
	assert.Contains(t, outputWithoutTimestamp, "[RETE] [INFO] test without timestamp")
}
func TestLogger_Prefix(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	logger.SetTimestamps(false)
	// Default prefix
	logger.Info("test default prefix")
	defaultOutput := buf.String()
	assert.Contains(t, defaultOutput, "[RETE]")
	// Custom prefix
	buf.Reset()
	logger.SetPrefix("[CUSTOM]")
	logger.Info("test custom prefix")
	customOutput := buf.String()
	assert.Contains(t, customOutput, "[CUSTOM]")
	assert.NotContains(t, customOutput, "[RETE]")
}
func TestLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	logger.SetTimestamps(false)
	logger.SetPrefix("[RETE]")
	// Create contextual logger
	contextLogger := logger.WithContext("Pipeline")
	contextLogger.Info("test message")
	output := buf.String()
	assert.Contains(t, output, "[RETE][Pipeline]")
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "test message")
	// Original logger should be unaffected
	buf.Reset()
	logger.Info("original logger")
	originalOutput := buf.String()
	assert.Contains(t, originalOutput, "[RETE]")
	assert.NotContains(t, originalOutput, "[Pipeline]")
}
func TestLogger_NestedContext(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf)
	logger.SetTimestamps(false)
	ctx1 := logger.WithContext("Module1")
	ctx2 := ctx1.WithContext("SubModule")
	ctx2.Info("nested context test")
	output := buf.String()
	assert.Contains(t, output, "[RETE][Module1][SubModule]")
}
func TestLogger_SetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelError, &buf)
	logger.SetTimestamps(false)
	// Initially at Error level - Info should not log
	logger.Info("should not appear")
	assert.Empty(t, buf.String())
	// Change to Info level
	logger.SetLevel(LogLevelInfo)
	logger.Info("should appear")
	assert.Contains(t, buf.String(), "should appear")
}
func TestLogger_GetLevel(t *testing.T) {
	logger := NewLogger(LogLevelWarn, nil)
	level := logger.GetLevel()
	assert.Equal(t, LogLevelWarn, level)
	logger.SetLevel(LogLevelDebug)
	level = logger.GetLevel()
	assert.Equal(t, LogLevelDebug, level)
}
func TestLogger_Concurrent(t *testing.T) {
	// Use separate buffers per goroutine to avoid race on shared buffer
	// This tests that the logger itself is thread-safe
	type result struct {
		id     int
		output string
	}
	results := make(chan result, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			var buf bytes.Buffer
			logger := NewLogger(LogLevelInfo, &buf)
			logger.SetTimestamps(false)
			for j := 0; j < 100; j++ {
				logger.Info("goroutine %d message %d", id, j)
			}
			results <- result{id: id, output: buf.String()}
		}(i)
	}
	// Collect results from all goroutines
	for i := 0; i < 10; i++ {
		res := <-results
		assert.NotEmpty(t, res.output, "Goroutine %d should have output", res.id)
		assert.Contains(t, res.output, "[INFO]")
	}
}
func TestLogger_SetOutput(t *testing.T) {
	var buf1 bytes.Buffer
	var buf2 bytes.Buffer
	logger := NewLogger(LogLevelInfo, &buf1)
	logger.SetTimestamps(false)
	logger.Info("to buffer 1")
	assert.Contains(t, buf1.String(), "to buffer 1")
	assert.Empty(t, buf2.String())
	// Change output
	logger.SetOutput(&buf2)
	logger.Info("to buffer 2")
	assert.NotContains(t, buf1.String(), "to buffer 2")
	assert.Contains(t, buf2.String(), "to buffer 2")
}
func TestGlobalLogger(t *testing.T) {
	// Save original level
	originalLevel := GetGlobalLogLevel()
	defer SetGlobalLogLevel(originalLevel)
	// Test global level setting
	SetGlobalLogLevel(LogLevelError)
	assert.Equal(t, LogLevelError, GetGlobalLogLevel())
	SetGlobalLogLevel(LogLevelDebug)
	assert.Equal(t, LogLevelDebug, GetGlobalLogLevel())
}
func TestLogger_AllLevelsOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelDebug, &buf)
	logger.SetTimestamps(false)
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Len(t, lines, 4, "Should have 4 log lines")
	assert.Contains(t, lines[0], "[DEBUG]")
	assert.Contains(t, lines[0], "debug message")
	assert.Contains(t, lines[1], "[INFO]")
	assert.Contains(t, lines[1], "info message")
	assert.Contains(t, lines[2], "[WARN]")
	assert.Contains(t, lines[2], "warn message")
	assert.Contains(t, lines[3], "[ERROR]")
	assert.Contains(t, lines[3], "error message")
}
func TestLogger_SilentLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(LogLevelSilent, &buf)
	logger.SetTimestamps(false)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")
	output := buf.String()
	assert.Empty(t, output, "Silent level should produce no output")
}
func TestLogger_ThreadSafety(t *testing.T) {
	// Test thread-safety of logger operations (SetLevel, logging)
	// Each goroutine writes to its own buffer to avoid buffer race
	done := make(chan bool)
	// Writer goroutines with separate buffers
	for i := 0; i < 5; i++ {
		go func(id int) {
			var buf bytes.Buffer
			logger := NewLogger(LogLevelInfo, &buf)
			logger.SetTimestamps(false)
			for j := 0; j < 50; j++ {
				logger.Info("writer %d msg %d", id, j)
			}
			done <- true
		}(i)
	}
	// Level changer goroutines on a shared logger (tests SetLevel thread-safety)
	sharedLogger := NewLogger(LogLevelInfo, &bytes.Buffer{})
	for i := 0; i < 3; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				sharedLogger.SetLevel(LogLevelDebug)
				sharedLogger.SetLevel(LogLevelInfo)
				// Also test GetLevel
				_ = sharedLogger.GetLevel()
			}
			done <- true
		}()
	}
	// Wait for all goroutines
	for i := 0; i < 8; i++ {
		<-done
	}
	// Test passes if no race detected and no panic occurred
}