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
	ConstraintFile string
	ConstraintText string
	UseStdin       bool
	FactsFile      string
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

	flagSet.StringVar(&config.ConstraintFile, "constraint", "", "Fichier de contraintes (.constraint)")
	flagSet.StringVar(&config.ConstraintText, "text", "", "Texte de contrainte directement (alternative à -constraint)")
	flagSet.BoolVar(&config.UseStdin, "stdin", false, "Lire les contraintes depuis stdin")
	flagSet.StringVar(&config.FactsFile, "facts", "", "Fichier de faits (.facts)")
	flagSet.BoolVar(&config.Verbose, "v", false, "Mode verbeux")
	flagSet.BoolVar(&config.ShowVersion, "version", false, "Afficher la version")
	flagSet.BoolVar(&config.ShowHelp, "h", false, "Afficher l'aide")

	if err := flagSet.Parse(args); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfig validates that exactly one input source is specified
func ValidateConfig(config *Config) error {
	sourcesCount := 0
	if config.ConstraintFile != "" {
		sourcesCount++
	}
	if config.ConstraintText != "" {
		sourcesCount++
	}
	if config.UseStdin {
		sourcesCount++
	}

	if sourcesCount == 0 {
		return fmt.Errorf("spécifiez une source (-constraint, -text, ou -stdin)")
	}

	if sourcesCount > 1 {
		return fmt.Errorf("spécifiez une seule source d'entrée")
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
	sourceName := config.ConstraintFile

	if _, err := os.Stat(config.ConstraintFile); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("fichier contrainte non trouvé: %s", config.ConstraintFile)
	}

	result, err := constraint.ParseConstraintFile(config.ConstraintFile)
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

	network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts(
		constraintSource,
		factsFile,
		storage,
	)

	if err != nil {
		return nil, err
	}

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
	fmt.Fprintln(w, "  tsd -constraint <file.constraint> [options]")
	fmt.Fprintln(w, "  tsd -text \"<constraint text>\" [options]")
	fmt.Fprintln(w, "  tsd -stdin [options]")
	fmt.Fprintln(w, "  echo \"<constraint>\" | tsd -stdin")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "OPTIONS:")
	fmt.Fprintln(w, "  -constraint <file>  Fichier de règles/contraintes")
	fmt.Fprintln(w, "  -text <string>      Texte de contrainte directement")
	fmt.Fprintln(w, "  -stdin              Lire les contraintes depuis stdin")
	fmt.Fprintln(w, "  -facts <file>       Fichier de faits (optionnel, pour futur usage)")
	fmt.Fprintln(w, "  -v                  Mode verbeux")
	fmt.Fprintln(w, "  -version            Afficher la version")
	fmt.Fprintln(w, "  -h                  Afficher cette aide")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "EXEMPLES:")
	fmt.Fprintln(w, "  tsd -constraint rules.constraint")
	fmt.Fprintln(w, "  tsd -constraint rules.constraint -v")
	fmt.Fprintln(w, "  tsd -text 'type Person : <id: string, name: string>'")
	fmt.Fprintln(w, "  echo 'type Person : <id: string>' | tsd -stdin")
	fmt.Fprintln(w, "  cat rules.constraint | tsd -stdin -v")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "FORMATS DE FICHIERS:")
	fmt.Fprintln(w, "  .constraint : Règles en syntaxe TSD")
	fmt.Fprintln(w, "  .facts      : Faits en format structuré (support futur)")
}
