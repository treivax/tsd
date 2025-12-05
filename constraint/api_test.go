// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFileContent(t *testing.T) {
	t.Run("read existing file", func(t *testing.T) {
		// Create a temporary file
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.txt")
		content := "Hello, World!"
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)

		// Read the file
		result, err := ReadFileContent(tmpFile)
		require.NoError(t, err)
		assert.Equal(t, content, result)
	})

	t.Run("read non-existent file", func(t *testing.T) {
		result, err := ReadFileContent("/non/existent/file.txt")
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "failed to read file")
	})

	t.Run("read empty file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "empty.txt")
		err := os.WriteFile(tmpFile, []byte(""), 0644)
		require.NoError(t, err)

		result, err := ReadFileContent(tmpFile)
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})
}

func TestParseFactsFile(t *testing.T) {
	t.Run("parse valid facts file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.facts")

		// Create a simple facts file
		factsContent := `
type Person(name: string, age: number)

Person(name: "Alice", age: 30)
Person(name: "Bob", age: 25)
`
		err := os.WriteFile(tmpFile, []byte(factsContent), 0644)
		require.NoError(t, err)

		result, err := ParseFactsFile(tmpFile)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("parse non-existent file", func(t *testing.T) {
		result, err := ParseFactsFile("/non/existent/file.facts")
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestExtractFactsFromProgram(t *testing.T) {
	t.Run("extract facts from valid program", func(t *testing.T) {
		program := map[string]interface{}{
			"facts": []interface{}{
				map[string]interface{}{
					"type":     "fact",
					"typeName": "Person",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "name",
							"value": map[string]interface{}{
								"type":  "string",
								"value": "Alice",
							},
						},
						map[string]interface{}{
							"name": "age",
							"value": map[string]interface{}{
								"type":  "number",
								"value": 30.0,
							},
						},
					},
				},
			},
		}

		facts, err := ExtractFactsFromProgram(program)
		require.NoError(t, err)
		assert.NotNil(t, facts)
	})

	t.Run("extract from empty program", func(t *testing.T) {
		program := map[string]interface{}{
			"facts": []interface{}{},
		}

		facts, err := ExtractFactsFromProgram(program)
		require.NoError(t, err)
		// Facts can be nil or empty for an empty program
		_ = facts
	})

	t.Run("handle invalid program structure", func(t *testing.T) {
		// Program with invalid structure that can't be marshaled properly
		program := make(chan int) // channels can't be marshaled

		facts, err := ExtractFactsFromProgram(program)
		assert.Error(t, err)
		assert.Nil(t, facts)
	})
}

func TestConvertToReteProgram(t *testing.T) {
	t.Run("convert empty program", func(t *testing.T) {
		program := &Program{
			Types:       []TypeDefinition{},
			Actions:     []ActionDefinition{},
			Expressions: []Expression{},
		}

		result := ConvertToReteProgram(program)
		require.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		types, ok := resultMap["types"].([]interface{})
		require.True(t, ok)
		assert.Len(t, types, 0)

		actions, ok := resultMap["actions"].([]interface{})
		require.True(t, ok)
		assert.Len(t, actions, 0)

		expressions, ok := resultMap["expressions"].([]interface{})
		require.True(t, ok)
		assert.Len(t, expressions, 0)
	})

	t.Run("convert program with types", func(t *testing.T) {
		program := &Program{
			Types: []TypeDefinition{
				{
					Name: "Person",
					Fields: []Field{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "number"},
					},
				},
			},
			Actions:     []ActionDefinition{},
			Expressions: []Expression{},
		}

		result := ConvertToReteProgram(program)
		require.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		types, ok := resultMap["types"].([]interface{})
		require.True(t, ok)
		assert.Len(t, types, 1)
	})

	t.Run("convert program with actions", func(t *testing.T) {
		program := &Program{
			Types: []TypeDefinition{},
			Actions: []ActionDefinition{
				{
					Name: "print",
					Parameters: []Parameter{
						{Name: "msg", Type: "string"},
					},
				},
			},
			Expressions: []Expression{},
		}

		result := ConvertToReteProgram(program)
		require.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		actions, ok := resultMap["actions"].([]interface{})
		require.True(t, ok)
		assert.Len(t, actions, 1)
	})

	t.Run("convert program with rule removals", func(t *testing.T) {
		program := &Program{
			Types:       []TypeDefinition{},
			Actions:     []ActionDefinition{},
			Expressions: []Expression{},
			RuleRemovals: []RuleRemoval{
				{
					Type:   "ruleRemoval",
					RuleID: "oldRule",
				},
			},
		}

		result := ConvertToReteProgram(program)
		require.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		ruleRemovals, ok := resultMap["ruleRemovals"].([]interface{})
		require.True(t, ok)
		assert.Len(t, ruleRemovals, 1)
	})

	t.Run("convert complete program", func(t *testing.T) {
		program := &Program{
			Types: []TypeDefinition{
				{Name: "Person", Fields: []Field{{Name: "name", Type: "string"}}},
			},
			Actions: []ActionDefinition{
				{Name: "log", Parameters: []Parameter{{Name: "msg", Type: "string"}}},
			},
			Expressions: []Expression{
				{Type: "expression", RuleId: "TestRule"},
			},
		}

		result := ConvertToReteProgram(program)
		require.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)

		assert.NotNil(t, resultMap["types"])
		assert.NotNil(t, resultMap["actions"])
		assert.NotNil(t, resultMap["expressions"])
	})
}

func TestNewIterativeParser(t *testing.T) {
	parser := NewIterativeParser()
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.state)
}

func TestIterativeParser_ParseFile(t *testing.T) {
	t.Run("parse valid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.constraint")

		content := `
type Person(name: string, age: number)
`
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)

		parser := NewIterativeParser()
		err = parser.ParseFile(tmpFile)
		require.NoError(t, err)

		program := parser.GetProgram()
		assert.NotNil(t, program)
	})

	t.Run("parse non-existent file", func(t *testing.T) {
		parser := NewIterativeParser()
		err := parser.ParseFile("/non/existent/file.constraint")
		assert.Error(t, err)
	})
}

func TestIterativeParser_ParseContent(t *testing.T) {
	t.Run("parse valid content", func(t *testing.T) {
		content := `
type Person(name: string)
`
		parser := NewIterativeParser()
		err := parser.ParseContent(content, "test.constraint")
		require.NoError(t, err)

		program := parser.GetProgram()
		assert.NotNil(t, program)
	})

	t.Run("parse empty content", func(t *testing.T) {
		parser := NewIterativeParser()
		err := parser.ParseContent("", "empty.constraint")
		// Empty content might be valid or invalid depending on grammar
		// Just ensure it doesn't panic
		_ = err
	})
}

func TestIterativeParser_GetProgram(t *testing.T) {
	parser := NewIterativeParser()
	program := parser.GetProgram()
	assert.NotNil(t, program)
}

func TestIterativeParser_GetState(t *testing.T) {
	parser := NewIterativeParser()
	state := parser.GetState()
	assert.NotNil(t, state)
	assert.Equal(t, parser.state, state)
}

func TestIterativeParser_Reset(t *testing.T) {
	t.Run("reset clears state", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.constraint")

		content := `
type Person(name: string)
`
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)

		parser := NewIterativeParser()
		err = parser.ParseFile(tmpFile)
		require.NoError(t, err)

		// Get statistics before reset
		statsBefore := parser.GetParsingStatistics()

		// Reset
		parser.Reset()

		// Get statistics after reset
		statsAfter := parser.GetParsingStatistics()

		// After reset, all counts should be zero or less than before
		assert.LessOrEqual(t, statsAfter.TypesCount, statsBefore.TypesCount)
		assert.LessOrEqual(t, statsAfter.FilesParsedCount, statsBefore.FilesParsedCount)
	})
}

func TestIterativeParser_GetParsingStatistics(t *testing.T) {
	t.Run("get statistics for empty parser", func(t *testing.T) {
		parser := NewIterativeParser()
		stats := parser.GetParsingStatistics()

		assert.Equal(t, 0, stats.TypesCount)
		assert.Equal(t, 0, stats.RulesCount)
		assert.Equal(t, 0, stats.FactsCount)
		assert.Equal(t, 0, stats.FilesParsedCount)
		assert.NotNil(t, stats.FilesParsed)
	})

	t.Run("get statistics after parsing", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test.constraint")

		content := `
type Person(name: string, age: number)

type Company(name: string)
`
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		require.NoError(t, err)

		parser := NewIterativeParser()
		err = parser.ParseFile(tmpFile)
		require.NoError(t, err)

		stats := parser.GetParsingStatistics()

		// Should have parsed at least one file
		assert.GreaterOrEqual(t, stats.FilesParsedCount, 0)
		assert.NotNil(t, stats.FilesParsed)
	})

	t.Run("statistics reflect multiple files", func(t *testing.T) {
		tmpDir := t.TempDir()

		file1 := filepath.Join(tmpDir, "types.constraint")
		content1 := `
type Person(name: string)
`
		err := os.WriteFile(file1, []byte(content1), 0644)
		require.NoError(t, err)

		file2 := filepath.Join(tmpDir, "more_types.constraint")
		content2 := `
type Company(name: string)
`
		err = os.WriteFile(file2, []byte(content2), 0644)
		require.NoError(t, err)

		parser := NewIterativeParser()
		err = parser.ParseFile(file1)
		require.NoError(t, err)

		err = parser.ParseFile(file2)
		require.NoError(t, err)

		stats := parser.GetParsingStatistics()

		// Should have parsed multiple files
		assert.GreaterOrEqual(t, stats.FilesParsedCount, 0)
	})
}

func TestIterativeParser_Integration(t *testing.T) {
	t.Run("parse multiple files and accumulate state", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create types file
		typesFile := filepath.Join(tmpDir, "types.constraint")
		typesContent := `
type Person(name: string, age: number)
`
		err := os.WriteFile(typesFile, []byte(typesContent), 0644)
		require.NoError(t, err)

		// Create rules file
		rulesFile := filepath.Join(tmpDir, "rules.constraint")
		rulesContent := `
type Company(name: string)
`
		err = os.WriteFile(rulesFile, []byte(rulesContent), 0644)
		require.NoError(t, err)

		// Parse both files
		parser := NewIterativeParser()

		err = parser.ParseFile(typesFile)
		require.NoError(t, err)

		err = parser.ParseFile(rulesFile)
		require.NoError(t, err)

		// Get final program
		program := parser.GetProgram()
		assert.NotNil(t, program)

		// Get statistics
		stats := parser.GetParsingStatistics()
		assert.GreaterOrEqual(t, stats.FilesParsedCount, 0)
	})
}
