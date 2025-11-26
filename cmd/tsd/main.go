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

func main() {
	config := parseFlags()

	if config.ShowHelp {
		printHelp()
		return
	}

	if config.ShowVersion {
		printVersion()
		return
	}

	if err := validateConfig(config); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n\n", err)
		printHelp()
		os.Exit(1)
	}

	result, sourceName, err := parseConstraintSource(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur de parsing: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Printf("✅ Parsing réussi\n")
		fmt.Printf("📋 Validation du programme...\n")
	}

	if err := constraint.ValidateConstraintProgram(result); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur de validation: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Printf("✅ Contraintes validées avec succès\n")
	}

	if config.FactsFile != "" {
		runWithFacts(config, sourceName)
	} else {
		runValidationOnly(config)
	}
}

// parseFlags parses command-line flags and returns a Config
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.ConstraintFile, "constraint", "", "Fichier de contraintes (.constraint)")
	flag.StringVar(&config.ConstraintText, "text", "", "Texte de contrainte directement (alternative à -constraint)")
	flag.BoolVar(&config.UseStdin, "stdin", false, "Lire les contraintes depuis stdin")
	flag.StringVar(&config.FactsFile, "facts", "", "Fichier de faits (.facts)")
	flag.BoolVar(&config.Verbose, "v", false, "Mode verbeux")
	flag.BoolVar(&config.ShowVersion, "version", false, "Afficher la version")
	flag.BoolVar(&config.ShowHelp, "h", false, "Afficher l'aide")

	flag.Parse()

	return config
}

// validateConfig validates that exactly one input source is specified
func validateConfig(config *Config) error {
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

// parseConstraintSource parses constraints from the configured source
func parseConstraintSource(config *Config) (interface{}, string, error) {
	if config.UseStdin {
		return parseFromStdin(config)
	}

	if config.ConstraintText != "" {
		return parseFromText(config)
	}

	return parseFromFile(config)
}

// parseFromStdin reads and parses constraints from stdin
func parseFromStdin(config *Config) (interface{}, string, error) {
	sourceName := "<stdin>"

	if config.Verbose {
		printParsingHeader("stdin")
	}

	stdinContent, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, "", fmt.Errorf("lecture stdin: %w", err)
	}

	result, err := constraint.ParseConstraint(sourceName, stdinContent)
	return result, sourceName, err
}

// parseFromText parses constraints from a text string
func parseFromText(config *Config) (interface{}, string, error) {
	sourceName := "<text>"

	if config.Verbose {
		printParsingHeader("texte direct")
	}

	result, err := constraint.ParseConstraint(sourceName, []byte(config.ConstraintText))
	return result, sourceName, err
}

// parseFromFile parses constraints from a file
func parseFromFile(config *Config) (interface{}, string, error) {
	sourceName := config.ConstraintFile

	if _, err := os.Stat(config.ConstraintFile); os.IsNotExist(err) {
		return nil, "", fmt.Errorf("fichier contrainte non trouvé: %s", config.ConstraintFile)
	}

	if config.Verbose {
		fmt.Printf("🚀 TSD - Analyse des contraintes\n")
		fmt.Printf("===============================\n")
		fmt.Printf("Fichier: %s\n\n", config.ConstraintFile)
	}

	result, err := constraint.ParseConstraintFile(config.ConstraintFile)
	return result, sourceName, err
}

// printParsingHeader prints the header for parsing operations
func printParsingHeader(source string) {
	fmt.Printf("🚀 TSD - Analyse des contraintes\n")
	fmt.Printf("===============================\n")
	fmt.Printf("Source: %s\n\n", source)
}

// runValidationOnly runs in validation-only mode (no facts file)
func runValidationOnly(config *Config) {
	fmt.Printf("✅ Contraintes validées avec succès\n")

	if config.Verbose {
		fmt.Printf("\n🎉 Validation terminée!\n")
		fmt.Printf("Les contraintes sont syntaxiquement correctes.\n")
		fmt.Printf("ℹ️  Utilisez -facts <file> pour exécuter le pipeline RETE complet.\n")
	}
}

// runWithFacts runs the full RETE pipeline with facts
func runWithFacts(config *Config, sourceName string) {
	if config.Verbose {
		fmt.Printf("\n🔧 PIPELINE RETE COMPLET\n")
		fmt.Printf("========================\n")
		fmt.Printf("Fichier faits: %s\n\n", config.FactsFile)
	}

	if _, err := os.Stat(config.FactsFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Fichier faits non trouvé: %s\n", config.FactsFile)
		os.Exit(1)
	}

	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts(
		sourceName,
		config.FactsFile,
		storage,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur pipeline RETE: %v\n", err)
		os.Exit(1)
	}

	printResults(config, network, facts)
}

// printResults prints the RETE pipeline execution results
func printResults(config *Config, network *rete.ReteNetwork, facts []*rete.Fact) {
	if config.Verbose {
		fmt.Printf("\n📊 RÉSULTATS\n")
		fmt.Printf("============\n")
		fmt.Printf("Faits injectés: %d\n", len(facts))
	}

	activations := countActivations(network)

	if activations > 0 {
		fmt.Printf("\n🎯 ACTIONS DISPONIBLES: %d\n", activations)
		if config.Verbose {
			printActivationDetails(network)
		}
	} else {
		fmt.Printf("\nℹ️  Aucune action déclenchée\n")
	}

	if config.Verbose {
		fmt.Printf("\n✅ Pipeline RETE exécuté avec succès\n")
	}
}

// countActivations counts the total number of activations in the network
func countActivations(network *rete.ReteNetwork) int {
	count := 0
	for _, terminal := range network.TerminalNodes {
		if terminal.Memory != nil && terminal.Memory.Tokens != nil {
			count += len(terminal.Memory.Tokens)
		}
	}
	return count
}

// printActivationDetails prints detailed information about activations
func printActivationDetails(network *rete.ReteNetwork) {
	count := 0
	for _, terminal := range network.TerminalNodes {
		if terminal.Memory != nil && terminal.Memory.Tokens != nil {
			actionName := "unknown"
			if terminal.Action != nil {
				actionName = terminal.Action.Job.Name
			}
			for _, token := range terminal.Memory.Tokens {
				count++
				fmt.Printf("  %d. %s() - %d bindings\n", count, actionName, len(token.Facts))
			}
		}
	}
}

// printVersion prints the version information
func printVersion() {
	fmt.Println("TSD (Type System Development) v1.0")
	fmt.Println("Moteur de règles basé sur l'algorithme RETE")
}

func printHelp() {
	fmt.Println("TSD - Type System Development")
	fmt.Println("Moteur de règles basé sur l'algorithme RETE")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  tsd -constraint <file.constraint> [options]")
	fmt.Println("  tsd -text \"<constraint text>\" [options]")
	fmt.Println("  tsd -stdin [options]")
	fmt.Println("  echo \"<constraint>\" | tsd -stdin")
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -constraint <file>  Fichier de règles/contraintes")
	fmt.Println("  -text <string>      Texte de contrainte directement")
	fmt.Println("  -stdin              Lire les contraintes depuis stdin")
	fmt.Println("  -facts <file>       Fichier de faits (optionnel, pour futur usage)")
	fmt.Println("  -v                  Mode verbeux")
	fmt.Println("  -version            Afficher la version")
	fmt.Println("  -h                  Afficher cette aide")
	fmt.Println("")
	fmt.Println("EXEMPLES:")
	fmt.Println("  tsd -constraint rules.constraint")
	fmt.Println("  tsd -constraint rules.constraint -v")
	fmt.Println("  tsd -text 'type Person : <id: string, name: string>'")
	fmt.Println("  echo 'type Person : <id: string>' | tsd -stdin")
	fmt.Println("  cat rules.constraint | tsd -stdin -v")
	fmt.Println("")
	fmt.Println("FORMATS DE FICHIERS:")
	fmt.Println("  .constraint : Règles en syntaxe TSD")
	fmt.Println("  .facts      : Faits en format structuré (support futur)")
}
