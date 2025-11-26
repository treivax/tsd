package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/treivax/tsd/constraint"
)

func main() {
	exitCode := Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

// Run executes the constraint parser CLI and returns an exit code
// This function is testable and doesn't call os.Exit
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		PrintHelp(stderr)
		return 1
	}

	inputFile := args[0]

	result, err := ParseFile(inputFile)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return 1
	}

	// Output JSON
	if err := OutputJSON(result, stdout); err != nil {
		fmt.Fprintf(stderr, "Erreur JSON: %v\n", err)
		return 1
	}

	return 0
}

// ParseFile reads and parses a constraint file, then validates it
func ParseFile(inputFile string) (interface{}, error) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("lecture fichier: %w", err)
	}

	// Parse the input
	result, err := constraint.ParseConstraint("", input)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	// Validate the program
	if err := constraint.ValidateConstraintProgram(result); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	return result, nil
}

// OutputJSON converts the result to JSON and writes it to the writer
func OutputJSON(result interface{}, w io.Writer) error {
	jsonOutput, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, string(jsonOutput))
	return err
}

// PrintHelp prints the help message
func PrintHelp(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  constraint-parser <input-file>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Exemple:")
	fmt.Fprintln(w, "  constraint-parser ../tests/test_input.txt")
}
