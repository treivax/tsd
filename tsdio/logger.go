// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package tsdio provides thread-safe I/O operations for the TSD project.
// All stdout/stderr writes should use this package to prevent race conditions
// when TSD is used concurrently from multiple goroutines.
package tsdio

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// stdoutMutex protects modifications to os.Stdout/os.Stderr globally.
// This prevents race conditions when test frameworks redirect standard streams.
// This mutex is separate from the logger's mutex to avoid deadlocks.
var stdoutMutex sync.Mutex

// Logger provides thread-safe logging with mutex-protected stdout/stderr writes.
// This prevents race conditions when multiple goroutines use TSD concurrently.
//
// Usage:
//
//	tsdio.Printf("Processing rule: %s", ruleID)
//	tsdio.Println("Operation completed")
type Logger struct {
	output       io.Writer
	logger       *log.Logger
	captureHooks []func() io.Writer // Hooks to dynamically resolve output during capture
}

// globalLogger is the global instance used throughout TSD
// Note: output is set to nil to force dynamic resolution to os.Stdout
var globalLogger = &Logger{
	output: nil,
	logger: log.New(os.Stdout, "", log.LstdFlags),
}

// Printf writes a formatted message to stdout with mutex protection
func (l *Logger) Printf(format string, v ...interface{}) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	output := l.resolveOutput()
	fmt.Fprintf(output, format, v...)
}

// Println writes a message with newline to stdout with mutex protection
func (l *Logger) Println(v ...interface{}) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	output := l.resolveOutput()
	fmt.Fprintln(output, v...)
}

// Print writes a message to stdout with mutex protection
func (l *Logger) Print(v ...interface{}) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	output := l.resolveOutput()
	fmt.Fprint(output, v...)
}

// LogPrintf writes a formatted log message with timestamp
func (l *Logger) LogPrintf(format string, v ...interface{}) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	l.logger.Printf(format, v...)
}

// resolveOutput returns the current output destination, checking hooks first
// Must be called while holding the mutex
func (l *Logger) resolveOutput() io.Writer {
	// Check capture hooks first (for testing/redirection)
	for _, hook := range l.captureHooks {
		if captured := hook(); captured != nil {
			return captured
		}
	}

	// Use explicitly set output
	if l.output != nil {
		return l.output
	}

	// Default to dynamic os.Stdout
	return os.Stdout
}

// getOutput returns the current output destination (for GetOutput API)
// Must be called while holding the mutex
func (l *Logger) getOutput() io.Writer {
	if l.output == nil {
		return os.Stdout
	}
	return l.output
}

// SetOutput changes the output destination (useful for testing)
func (l *Logger) SetOutput(w io.Writer) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	l.output = w
	if w == nil {
		l.logger.SetOutput(os.Stdout)
	} else {
		l.logger.SetOutput(w)
	}
}

// GetOutput returns the current output destination
func (l *Logger) GetOutput() io.Writer {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	return l.getOutput()
}

// Mute disables all output (useful for testing or silent mode)
func (l *Logger) Mute() {
	l.SetOutput(io.Discard)
}

// Unmute restores output to stdout (dynamic)
func (l *Logger) Unmute() {
	l.SetOutput(nil) // nil means dynamic os.Stdout
}

// WithMutex executes a function while holding the stdout mutex.
// This is useful for atomic multi-line operations.
//
// Example:
//
//	tsdio.WithMutex(func() {
//	    fmt.Println("Line 1")
//	    fmt.Println("Line 2")
//	})
func (l *Logger) WithMutex(fn func()) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	fn()
}

// AddCaptureHook adds a function that will be called to resolve output dynamically.
// This is useful for test frameworks that redirect stdout.
// The hook should return nil if it doesn't want to capture, or an io.Writer to use.
// Hooks are checked in order until one returns non-nil.
func (l *Logger) AddCaptureHook(hook func() io.Writer) {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	l.captureHooks = append(l.captureHooks, hook)
}

// ClearCaptureHooks removes all capture hooks
func (l *Logger) ClearCaptureHooks() {
	stdoutMutex.Lock()
	defer stdoutMutex.Unlock()
	l.captureHooks = nil
}

// LockStdout acquires the global stdout mutex (for external synchronization of os.Stdout)
// Use this when you need to redirect os.Stdout/os.Stderr safely.
func LockStdout() {
	stdoutMutex.Lock()
}

// UnlockStdout releases the global stdout mutex
func UnlockStdout() {
	stdoutMutex.Unlock()
}

// NewLogger creates a new isolated Logger instance.
// Most code should use the global functions instead.
func NewLogger(output io.Writer) *Logger {
	return &Logger{
		output: output,
		logger: log.New(output, "", log.LstdFlags),
	}
}

// Package-level convenience functions using the global logger

// Printf is a thread-safe version of fmt.Printf
func Printf(format string, v ...interface{}) {
	globalLogger.Printf(format, v...)
}

// Println is a thread-safe version of fmt.Println
func Println(v ...interface{}) {
	globalLogger.Println(v...)
}

// Print is a thread-safe version of fmt.Print
func Print(v ...interface{}) {
	globalLogger.Print(v...)
}

// LogPrintf writes a formatted log message with timestamp (thread-safe)
func LogPrintf(format string, v ...interface{}) {
	globalLogger.LogPrintf(format, v...)
}

// SetOutput changes the global output destination
func SetOutput(w io.Writer) {
	globalLogger.SetOutput(w)
}

// GetOutput returns the current global output destination
func GetOutput() io.Writer {
	return globalLogger.GetOutput()
}

// Mute disables all global output
func Mute() {
	globalLogger.Mute()
}

// Unmute restores global output to stdout (dynamic)
func Unmute() {
	globalLogger.Unmute()
}

// WithMutex executes a function while holding the global logger mutex
func WithMutex(fn func()) {
	globalLogger.WithMutex(fn)
}

// AddCaptureHook adds a capture hook to the global logger
func AddCaptureHook(hook func() io.Writer) {
	globalLogger.AddCaptureHook(hook)
}

// ClearCaptureHooks removes all capture hooks from the global logger
func ClearCaptureHooks() {
	globalLogger.ClearCaptureHooks()
}

// GetGlobalLogger returns the global logger instance (for advanced usage)
func GetGlobalLogger() *Logger {
	return globalLogger
}
