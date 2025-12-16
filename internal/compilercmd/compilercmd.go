// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package compilercmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

// Constants for input validation and security
const (
	// MaxInputSize limits the maximum size of input that can be parsed
	// This prevents potential DoS attacks via extremely large inputs
	MaxInputSize = 10 * 1024 * 1024 // 10 MB

	// MaxStdinRead limits the amount of data read from stdin
	MaxStdinRead = MaxInputSize

	// Version information
	ApplicationName        = "TSD (Type System Development)"
	ApplicationVersion     = "v1.0"
	ApplicationDescription = "Moteur de règles basé sur l'algorithme RETE"

	// Source name constants
	SourceNameStdin = "<stdin>"
	SourceNameText  = "<text>"

	// Exit codes
	ExitSuccess         = 0
	ExitErrorGeneric    = 1
	ExitErrorParsing    = 1
	ExitErrorValidation = 1
	ExitErrorFileAccess = 1
	ExitErrorExecution  = 1
)

// Error messages
var (
	ErrNoSource        = errors.New("aucune source spécifiée (-file, -text ou -stdin)")
	ErrMultipleSources = errors.New("une seule source autorisée (-file, -text ou -stdin)")
	ErrFileNotFound    = errors.New("fichier non trouvé")
	ErrInputTooLarge   = errors.New("entrée trop volumineuse")
	ErrInvalidPath     = errors.New("chemin de fichier non valide")
	ErrPathTraversal   = errors.New("tentative de traversée de répertoire interdite")
)

// Config holds the CLI configuration
type Config struct {
	File           string // Unified .tsd file
	ConstraintFile string // Deprecated: use File instead
	ConstraintText string
	UseStdin       bool
	FactsFile      string // Deprecated: use File instead
	Verbose        bool
	ShowVersion    bool
	ShowHelp       bool
}

// Result holds the execution result
type Result struct {
	Network     *rete.ReteNetwork
	Facts       []*rete.Fact
	Activations int
	Error       error
}

// Run executes the TSD compiler with the given arguments and returns an exit code
// This function is the main entry point for the compiler command
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	config, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return ExitErrorGeneric
	}

	if config.ShowHelp {
		printHelp(stdout)
		return ExitSuccess
	}

	if config.ShowVersion {
		printVersion(stdout)
		return ExitSuccess
	}

	if err := validateConfig(config); err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n\n", err)
		printHelp(stderr)
		return ExitErrorGeneric
	}

	result, sourceName, err := parseConstraintSource(config, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur de parsing: %v\n", err)
		return ExitErrorParsing
	}

	if config.Verbose {
		fmt.Fprintf(stdout, "✅ Parsing réussi\n")
		fmt.Fprintf(stdout, "📋 Validation du programme...\n")
	}

	if err := constraint.ValidateConstraintProgram(result); err != nil {
		fmt.Fprintf(stderr, "Erreur de validation: %v\n", err)
		return ExitErrorValidation
	}

	if config.Verbose {
		fmt.Fprintf(stdout, "✅ Contraintes validées avec succès\n")
	}

	if config.FactsFile != "" {
		return runWithFacts(config, sourceName, stdout, stderr)
	}

	return runValidationOnly(config, stdout)
}

// ParseFlags parses command-line flags and returns a Config
func ParseFlags(args []string) (*Config, error) {
	config := &Config{}
	flagSet := flag.NewFlagSet("tsd", flag.ContinueOnError)

	flagSet.StringVar(&config.File, "file", "", "Fichier TSD (.tsd)")
	flagSet.StringVar(&config.ConstraintFile, "constraint", "", "Deprecated: use -file instead (fichier .constraint)")
	flagSet.StringVar(&config.ConstraintText, "text", "", "Texte de contrainte directement (alternative à -file)")
	flagSet.BoolVar(&config.UseStdin, "stdin", false, "Lire depuis stdin")
	flagSet.StringVar(&config.FactsFile, "facts", "", "Deprecated: use -file instead (fichier .facts)")
	flagSet.BoolVar(&config.Verbose, "v", false, "Mode verbeux")
	flagSet.BoolVar(&config.ShowVersion, "version", false, "Afficher la version")
	flagSet.BoolVar(&config.ShowHelp, "h", false, "Afficher l'aide")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	// Handle backward compatibility: map old flags to new File field
	if config.ConstraintFile != "" && config.File == "" {
		fmt.Fprintln(os.Stderr, "⚠️  Warning: -constraint flag is deprecated, use -file instead")
		config.File = config.ConstraintFile
	}

	// Handle positional argument as file
	if config.File == "" && len(flagSet.Args()) > 0 {
		config.File = flagSet.Args()[0]
	}

	return config, nil
}

// validateConfig validates that exactly one input source is specified
func validateConfig(config *Config) error {
	sourcesCount := 0
	if config.File != "" {
		sourcesCount++
	}
	if config.ConstraintText != "" {
		sourcesCount++
	}
	if config.UseStdin {
		sourcesCount++
	}

	if sourcesCount == 0 {
		return ErrNoSource
	}

	if sourcesCount > 1 {
		return ErrMultipleSources
	}

	return nil
}

// parseConstraintSource parses constraints from the configured source
func parseConstraintSource(config *Config, stdin io.Reader) (interface{}, string, error) {
	if config.UseStdin {
		return parseFromStdin(config, stdin)
	}

	if config.ConstraintText != "" {
		return parseFromText(config)
	}

	return parseFromFile(config)
}

// parseFromStdin reads and parses constraints from stdin
func parseFromStdin(config *Config, stdin io.Reader) (interface{}, string, error) {
	sourceName := SourceNameStdin

	// Use LimitReader to prevent reading unbounded data from stdin
	limitedReader := io.LimitReader(stdin, MaxStdinRead+1)
	stdinContent, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, "", fmt.Errorf("lecture stdin: %w", err)
	}

	// Check if input exceeds maximum size
	if len(stdinContent) > MaxStdinRead {
		return nil, "", fmt.Errorf("%w: maximum %d bytes", ErrInputTooLarge, MaxStdinRead)
	}

	result, err := constraint.ParseConstraint(sourceName, stdinContent)
	return result, sourceName, err
}

// parseFromText parses constraints from a text string
func parseFromText(config *Config) (interface{}, string, error) {
	sourceName := SourceNameText

	// Validate text input size
	if len(config.ConstraintText) > MaxInputSize {
		return nil, "", fmt.Errorf("%w: maximum %d bytes", ErrInputTooLarge, MaxInputSize)
	}

	result, err := constraint.ParseConstraint(sourceName, []byte(config.ConstraintText))
	return result, sourceName, err
}

// parseFromFile parses constraints from a file
func parseFromFile(config *Config) (interface{}, string, error) {
	sourceName := config.File

	// Validate and sanitize file path
	if err := validateFilePath(config.File); err != nil {
		return nil, "", err
	}

	// Check file existence with proper error message
	fileInfo, err := os.Stat(config.File)
	if os.IsNotExist(err) {
		return nil, "", fmt.Errorf("%w: %s", ErrFileNotFound, config.File)
	}
	if err != nil {
		return nil, "", fmt.Errorf("erreur accès fichier: %w", err)
	}

	// Validate file size to prevent potential DoS
	if fileInfo.Size() > MaxInputSize {
		return nil, "", fmt.Errorf("%w: fichier %s dépasse %d bytes", ErrInputTooLarge, config.File, MaxInputSize)
	}

	result, err := constraint.ParseConstraintFile(config.File)
	return result, sourceName, err
}

// validateFilePath validates a file path for security
// - Prevents path traversal attacks
// - Validates path format
// - Ensures path is clean and safe
func validateFilePath(path string) error {
	if path == "" {
		return fmt.Errorf("%w: chemin vide", ErrInvalidPath)
	}

	// Clean the path (removes .., ., // etc.)
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	// After cleaning, ".." should not appear in the path
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("%w: '..' détecté dans %s", ErrPathTraversal, path)
	}

	// Additional check: ensure the path doesn't escape when made absolute
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("%w: impossible de résoudre le chemin absolu: %w", ErrInvalidPath, err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		// If we can't get cwd, we can't validate, but we've already done basic checks
		// This is acceptable as we've blocked obvious traversal attempts
		return nil
	}

	// If path is not absolute, ensure it stays within current directory tree
	if !filepath.IsAbs(path) {
		// The absolute path should be within or equal to cwd or its subdirectories
		// We allow files in current directory and below
		if !strings.HasPrefix(absPath, cwd) {
			return fmt.Errorf("%w: le fichier sort du répertoire courant", ErrPathTraversal)
		}
	}

	return nil
}

// runValidationOnly runs in validation-only mode (no facts file)
func runValidationOnly(config *Config, stdout io.Writer) int {
	fmt.Fprintf(stdout, "✅ Contraintes validées avec succès\n")

	if config.Verbose {
		fmt.Fprintf(stdout, "\n🎉 Validation terminée!\n")
		fmt.Fprintf(stdout, "Les contraintes sont syntaxiquement correctes.\n")
		fmt.Fprintf(stdout, "ℹ️  Utilisez -facts <file> pour exécuter le pipeline RETE complet.\n")
	}

	return ExitSuccess
}

// runWithFacts runs the full RETE pipeline with facts and returns exit code
func runWithFacts(config *Config, sourceName string, stdout, stderr io.Writer) int {
	if config.Verbose {
		fmt.Fprintf(stdout, "\n🔧 PIPELINE RETE COMPLET\n")
		fmt.Fprintf(stdout, "========================\n")
		fmt.Fprintf(stdout, "Fichier faits: %s\n\n", config.FactsFile)
	}

	// Validate facts file path
	if err := validateFilePath(config.FactsFile); err != nil {
		fmt.Fprintf(stderr, "Erreur validation chemin faits: %v\n", err)
		return ExitErrorFileAccess
	}

	// Check file existence with proper error handling
	fileInfo, err := os.Stat(config.FactsFile)
	if os.IsNotExist(err) {
		fmt.Fprintf(stderr, "%v: %s\n", ErrFileNotFound, config.FactsFile)
		return ExitErrorFileAccess
	}
	if err != nil {
		fmt.Fprintf(stderr, "Erreur accès fichier faits: %v\n", err)
		return ExitErrorFileAccess
	}

	// Validate facts file size
	if fileInfo.Size() > MaxInputSize {
		fmt.Fprintf(stderr, "Fichier faits trop volumineux: %s (max %d bytes)\n", config.FactsFile, MaxInputSize)
		return ExitErrorFileAccess
	}

	result, err := executePipeline(sourceName, config.FactsFile)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur pipeline RETE: %v\n", err)
		return ExitErrorExecution
	}

	printResults(config, result, stdout)
	return ExitSuccess
}

// executePipeline executes the RETE pipeline and returns the result
func executePipeline(constraintSource, factsFile string) (*Result, error) {
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	// Ingest constraint file
	network, _, err := pipeline.IngestFile(constraintSource, nil, storage)
	if err != nil {
		return nil, err
	}

	// Ingest facts file only if it's different from the constraint source
	// (to avoid double-ingesting the same file)
	if factsFile != constraintSource {
		network, _, err = pipeline.IngestFile(factsFile, network, storage)
		if err != nil {
			return nil, err
		}
	}

	// Collect facts from storage
	facts := storage.GetAllFacts()

	activations := countActivations(network)

	return &Result{
		Network:     network,
		Facts:       facts,
		Activations: activations,
	}, nil
}

// printResults prints the RETE pipeline execution results
func printResults(config *Config, result *Result, stdout io.Writer) {
	if config.Verbose {
		fmt.Fprintf(stdout, "\n📊 RÉSULTATS\n")
		fmt.Fprintf(stdout, "============\n")
		fmt.Fprintf(stdout, "Faits injectés: %d\n", len(result.Facts))
	}

	if result.Activations > 0 {
		fmt.Fprintf(stdout, "\n🎯 ACTIONS DISPONIBLES: %d\n", result.Activations)
		if config.Verbose {
			printActivationDetails(result.Network, stdout)
		}
	} else {
		fmt.Fprintf(stdout, "\nℹ️  Aucune action déclenchée\n")
	}

	if config.Verbose {
		fmt.Fprintf(stdout, "\n✅ Pipeline RETE exécuté avec succès\n")
	}
}

// countActivations counts the total number of activations in the network
func countActivations(network *rete.ReteNetwork) int {
	if network == nil {
		return 0
	}
	count := 0
	for _, terminal := range network.TerminalNodes {
		if terminal.Memory != nil && terminal.Memory.Tokens != nil {
			count += len(terminal.Memory.Tokens)
		}
	}
	return count
}

// printActivationDetails prints detailed information about activations
func printActivationDetails(network *rete.ReteNetwork, stdout io.Writer) {
	if network == nil {
		return
	}
	count := 0
	for _, terminal := range network.TerminalNodes {
		if terminal.Memory != nil && terminal.Memory.Tokens != nil {
			actionName := "unknown"
			if terminal.Action != nil {
				actionName = terminal.Action.Job.Name
			}
			for _, token := range terminal.Memory.Tokens {
				count++
				fmt.Fprintf(stdout, "  %d. %s() - %d bindings\n", count, actionName, len(token.Facts))
			}
		}
	}
}

// printVersion prints the version information
func printVersion(w io.Writer) {
	fmt.Fprintf(w, "%s %s\n", ApplicationName, ApplicationVersion)
	fmt.Fprintln(w, ApplicationDescription)
}

// printHelp prints the help message
func printHelp(w io.Writer) {
	fmt.Fprintln(w, "TSD - Type System Development")
	fmt.Fprintln(w, "Moteur de règles basé sur l'algorithme RETE")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "USAGE:")
	fmt.Fprintln(w, "  tsd <file.tsd> [options]")
	fmt.Fprintln(w, "  tsd -file <file.tsd> [options]")
	fmt.Fprintln(w, "  tsd -text \"<tsd code>\" [options]")
	fmt.Fprintln(w, "  tsd -stdin [options]")
	fmt.Fprintln(w, "  echo \"<tsd code>\" | tsd -stdin")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "OPTIONS:")
	fmt.Fprintln(w, "  -file <file>        Fichier TSD (.tsd)")
	fmt.Fprintln(w, "  -text <text>        Code TSD directement")
	fmt.Fprintln(w, "  -stdin              Lire depuis l'entrée standard")
	fmt.Fprintln(w, "  -facts <file>       [DEPRECATED] Use -file instead")
	fmt.Fprintln(w, "  -constraint <file>  [DEPRECATED] Use -file instead")
	fmt.Fprintln(w, "  -v                  Mode verbeux (affiche plus de détails)")
	fmt.Fprintln(w, "  -version            Afficher la version")
	fmt.Fprintln(w, "  -h                  Afficher cette aide")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "EXEMPLES:")
	fmt.Fprintln(w, "  tsd program.tsd")
	fmt.Fprintln(w, "  tsd -file program.tsd -v")
	fmt.Fprintln(w, "  tsd -text 'type Person : <id: string, name: string>'")
	fmt.Fprintln(w, "  echo 'type Person : <id: string>' | tsd -stdin")
	fmt.Fprintln(w, "  cat program.tsd | tsd -stdin -v")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "FORMAT DE FICHIER:")
	fmt.Fprintln(w, "  .tsd : Fichiers TSD (types, facts, rules)")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Un fichier .tsd peut contenir:")
	fmt.Fprintln(w, "  - Définitions de types: type Person : <id: string, name: string>")
	fmt.Fprintln(w, "  - Assertions de faits: Person(\"p1\", \"Alice\")")
	fmt.Fprintln(w, "  - Règles: rule r1 : {p: Person} / p.name == \"Alice\" ==> match(p.id)")
}
