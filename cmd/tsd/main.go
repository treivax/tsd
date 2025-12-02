// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
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

func main() {
	exitCode := Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

// Run executes the TSD CLI with the given arguments and returns an exit code
// This function is testable and doesn't call os.Exit
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	config, err := ParseFlags(args)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return 1
	}

	if config.ShowHelp {
		PrintHelp(stdout)
		return 0
	}

	if config.ShowVersion {
		PrintVersion(stdout)
		return 0
	}

	if err := ValidateConfig(config); err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n\n", err)
		PrintHelp(stderr)
		return 1
	}

	result, sourceName, err := ParseConstraintSource(config, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur de parsing: %v\n", err)
		return 1
	}

	if config.Verbose {
		fmt.Fprintf(stdout, "✅ Parsing réussi\n")
		fmt.Fprintf(stdout, "📋 Validation du programme...\n")
	}

	if err := constraint.ValidateConstraintProgram(result); err != nil {
		fmt.Fprintf(stderr, "Erreur de validation: %v\n", err)
		return 1
	}

	if config.Verbose {
		fmt.Fprintf(stdout, "✅ Contraintes validées avec succès\n")
	}

	if config.FactsFile != "" {
		return RunWithFacts(config, sourceName, stdout, stderr)
	}

	return RunValidationOnly(config, stdout)
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

// ValidateConfig validates that exactly one input source is specified
func ValidateConfig(config *Config) error {
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
		return fmt.Errorf("aucune source spécifiée (-file, -text ou -stdin)")
	}

	if sourcesCount > 1 {
		return fmt.Errorf("une seule source autorisée (-file, -text ou -stdin)")
	}

	return nil
}

// ParseConstraintSource parses constraints from the configured source
func ParseConstraintSource(config *Config, stdin io.Reader) (interface{}, string, error) {
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
	sourceName := "<stdin>"

	stdinContent, err := io.ReadAll(stdin)
	if err != nil {
		return nil, "", fmt.Errorf("lecture stdin: %w", err)
	}

	result, err := constraint.ParseConstraint(sourceName, stdinContent)
	return result, sourceName, err
}

// parseFromText parses constraints from a text string
func parseFromText(config *Config) (interface{}, string, error) {
	sourceName := "<text>"

	result, err := constraint.ParseConstraint(sourceName, []byte(config.ConstraintText))
	return result, sourceName, err
}

// parseFromFile parses constraints from a file
func parseFromFile(config *Config) (interface{}, string, error) {
	sourceName := config.File

	if _, err := os.Stat(config.File); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("fichier non trouvé: %s", config.File)
	}

	result, err := constraint.ParseConstraintFile(config.File)
	return result, sourceName, err
}

// RunValidationOnly runs in validation-only mode (no facts file)
func RunValidationOnly(config *Config, stdout io.Writer) int {
	fmt.Fprintf(stdout, "✅ Contraintes validées avec succès\n")

	if config.Verbose {
		fmt.Fprintf(stdout, "\n🎉 Validation terminée!\n")
		fmt.Fprintf(stdout, "Les contraintes sont syntaxiquement correctes.\n")
		fmt.Fprintf(stdout, "ℹ️  Utilisez -facts <file> pour exécuter le pipeline RETE complet.\n")
	}

	return 0
}

// RunWithFacts runs the full RETE pipeline with facts and returns exit code
func RunWithFacts(config *Config, sourceName string, stdout, stderr io.Writer) int {
	if config.Verbose {
		fmt.Fprintf(stdout, "\n🔧 PIPELINE RETE COMPLET\n")
		fmt.Fprintf(stdout, "========================\n")
		fmt.Fprintf(stdout, "Fichier faits: %s\n\n", config.FactsFile)
	}

	if _, err := os.Stat(config.FactsFile); os.IsNotExist(err) {
		fmt.Fprintf(stderr, "Fichier faits non trouvé: %s\n", config.FactsFile)
		return 1
	}

	result, err := ExecutePipeline(sourceName, config.FactsFile)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur pipeline RETE: %v\n", err)
		return 1
	}

	PrintResults(config, result, stdout)
	return 0
}

// ExecutePipeline executes the RETE pipeline and returns the result
func ExecutePipeline(constraintSource, factsFile string) (*Result, error) {
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	// Ingest constraint file
	network, err := pipeline.IngestFile(constraintSource, nil, storage)
	if err != nil {
		return nil, err
	}

	// Ingest facts file
	network, err = pipeline.IngestFile(factsFile, network, storage)
	if err != nil {
		return nil, err
	}

	// Collect facts from storage
	facts := storage.GetAllFacts()

	activations := CountActivations(network)

	return &Result{
		Network:     network,
		Facts:       facts,
		Activations: activations,
	}, nil
}

// PrintResults prints the RETE pipeline execution results
func PrintResults(config *Config, result *Result, stdout io.Writer) {
	if config.Verbose {
		fmt.Fprintf(stdout, "\n📊 RÉSULTATS\n")
		fmt.Fprintf(stdout, "============\n")
		fmt.Fprintf(stdout, "Faits injectés: %d\n", len(result.Facts))
	}

	if result.Activations > 0 {
		fmt.Fprintf(stdout, "\n🎯 ACTIONS DISPONIBLES: %d\n", result.Activations)
		if config.Verbose {
			PrintActivationDetails(result.Network, stdout)
		}
	} else {
		fmt.Fprintf(stdout, "\nℹ️  Aucune action déclenchée\n")
	}

	if config.Verbose {
		fmt.Fprintf(stdout, "\n✅ Pipeline RETE exécuté avec succès\n")
	}
}

// CountActivations counts the total number of activations in the network
func CountActivations(network *rete.ReteNetwork) int {
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

// PrintActivationDetails prints detailed information about activations
func PrintActivationDetails(network *rete.ReteNetwork, stdout io.Writer) {
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

// PrintVersion prints the version information
func PrintVersion(w io.Writer) {
	fmt.Fprintln(w, "TSD (Type System Development) v1.0")
	fmt.Fprintln(w, "Moteur de règles basé sur l'algorithme RETE")
}

// PrintHelp prints the help message
func PrintHelp(w io.Writer) {
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
