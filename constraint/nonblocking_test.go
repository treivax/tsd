package constraint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNonBlockingValidation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nonblocking_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	mixedFile := filepath.Join(tempDir, "mixed.constraint")
	mixedContent := "type Person : <id: string, name: string, age: number>\n\nPerson(id: \"P001\", name: \"Alice\", age: 30)\nUnknownType(id: \"U001\", field: \"value\")\nPerson(id: \"P002\", name: \"Bob\", age: 25)\n"

	err = os.WriteFile(mixedFile, []byte(mixedContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write mixed file: %v", err)
	}

	ps := NewProgramState()
	err = ps.ParseAndMerge(mixedFile)

	if err != nil {
		t.Fatalf("ParseAndMerge should not fail with validation errors: %v", err)
	}

	if len(ps.Facts) != 2 {
		t.Errorf("Expected 2 valid facts, got %d", len(ps.Facts))
	}

	if !ps.HasErrors() {
		t.Error("Expected validation errors to be recorded")
	}

	errorCount := ps.GetErrorCount()
	if errorCount < 1 {
		t.Errorf("Expected at least 1 error, got %d", errorCount)
	}

	t.Logf("Non-blocking: %d types, %d valid facts, %d errors", len(ps.Types), len(ps.Facts), errorCount)
}
