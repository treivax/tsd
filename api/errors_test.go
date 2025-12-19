// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"strings"
	"testing"
)

func TestErrorTypes(t *testing.T) {
	t.Log("🧪 TEST ERROR TYPES")

	types := []ErrorType{
		ErrorTypeParse,
		ErrorTypeValidation,
		ErrorTypeExecution,
		ErrorTypeConfig,
		ErrorTypeIO,
		ErrorTypeInternal,
	}

	expected := []string{
		"parse",
		"validation",
		"execution",
		"config",
		"io",
		"internal",
	}

	for i, errType := range types {
		if string(errType) != expected[i] {
			t.Errorf("❌ Type %d attendu: %s, reçu: %s",
				i, expected[i], errType)
		}
	}

	t.Log("✅ Types d'erreur corrects")
}

func TestError_Error(t *testing.T) {
	t.Log("🧪 TEST ERROR ERROR METHOD")

	err := &Error{
		Type:    ErrorTypeParse,
		Message: "test error",
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "parse") {
		t.Errorf("❌ Message devrait contenir 'parse', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "test error") {
		t.Errorf("❌ Message devrait contenir 'test error', reçu: %s", errStr)
	}

	t.Log("✅ Méthode Error() fonctionne")
}

func TestError_WithCause(t *testing.T) {
	t.Log("🧪 TEST ERROR WITH CAUSE")

	cause := &ConfigError{Field: "test", Message: "invalid"}
	err := &Error{
		Type:    ErrorTypeConfig,
		Message: "config error",
		Cause:   cause,
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "config") {
		t.Errorf("❌ Message devrait contenir 'config', reçu: %s", errStr)
	}

	if err.Unwrap() != cause {
		t.Error("❌ Unwrap devrait retourner la cause")
	}

	t.Log("✅ Error avec cause fonctionne")
}

func TestParseError_Error(t *testing.T) {
	t.Log("🧪 TEST PARSE ERROR ERROR METHOD")

	err := &ParseError{
		Filename: "test.tsd",
		Line:     10,
		Column:   5,
		Message:  "syntax error",
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "test.tsd") {
		t.Errorf("❌ Message devrait contenir 'test.tsd', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "10") {
		t.Errorf("❌ Message devrait contenir '10', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "5") {
		t.Errorf("❌ Message devrait contenir '5', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "syntax error") {
		t.Errorf("❌ Message devrait contenir 'syntax error', reçu: %s", errStr)
	}

	t.Log("✅ ParseError.Error() fonctionne")
}

func TestParseError_WithCause(t *testing.T) {
	t.Log("🧪 TEST PARSE ERROR WITH CAUSE")

	cause := &Error{Type: ErrorTypeInternal, Message: "internal"}
	err := &ParseError{
		Filename: "test.tsd",
		Line:     1,
		Column:   1,
		Message:  "parse failed",
		Cause:    cause,
	}

	if err.Unwrap() != cause {
		t.Error("❌ Unwrap devrait retourner la cause")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "internal") {
		t.Errorf("❌ Message devrait inclure la cause, reçu: %s", errStr)
	}

	t.Log("✅ ParseError avec cause fonctionne")
}

func TestConfigError_Error(t *testing.T) {
	t.Log("🧪 TEST CONFIG ERROR ERROR METHOD")

	err := &ConfigError{
		Field:   "LogLevel",
		Message: "invalid value",
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "LogLevel") {
		t.Errorf("❌ Message devrait contenir 'LogLevel', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "invalid value") {
		t.Errorf("❌ Message devrait contenir 'invalid value', reçu: %s", errStr)
	}

	t.Log("✅ ConfigError.Error() fonctionne")
}

func TestXupleSpaceError_Error(t *testing.T) {
	t.Log("🧪 TEST XUPLE SPACE ERROR ERROR METHOD")

	err := &XupleSpaceError{
		SpaceName: "alerts",
		Operation: "Retrieve",
		Message:   "not found",
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "alerts") {
		t.Errorf("❌ Message devrait contenir 'alerts', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "Retrieve") {
		t.Errorf("❌ Message devrait contenir 'Retrieve', reçu: %s", errStr)
	}
	if !strings.Contains(errStr, "not found") {
		t.Errorf("❌ Message devrait contenir 'not found', reçu: %s", errStr)
	}

	t.Log("✅ XupleSpaceError.Error() fonctionne")
}

func TestXupleSpaceError_WithCause(t *testing.T) {
	t.Log("🧪 TEST XUPLE SPACE ERROR WITH CAUSE")

	cause := &Error{Type: ErrorTypeInternal, Message: "internal"}
	err := &XupleSpaceError{
		SpaceName: "events",
		Operation: "Insert",
		Message:   "failed",
		Cause:     cause,
	}

	if err.Unwrap() != cause {
		t.Error("❌ Unwrap devrait retourner la cause")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "internal") {
		t.Errorf("❌ Message devrait inclure la cause, reçu: %s", errStr)
	}

	t.Log("✅ XupleSpaceError avec cause fonctionne")
}
