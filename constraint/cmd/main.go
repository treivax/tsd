// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/constraint/internal/config"
)

const (
	// Exit codes
	ExitSuccess       = 0
	ExitUsageError    = 1
	ExitRuntimeError  = 2
	ExitInvalidConfig = 3

	// Application metadata
	AppVersion = "1.0.0"
	AppName    = "constraint-parser"

	// Default values
	DefaultConfigPath   = ""
	DefaultOutputFormat = "json"
	StdinPlaceholder    = "-"
)

// CLIConfig contient la configuration de la ligne de commande
type CLIConfig struct {
	InputFile    string
	ConfigFile   string
	OutputFormat string
	Debug        bool
	Version      bool
	Help         bool
}

func main() {
	exitCode := Run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

// Run executes the constraint parser CLI and returns an exit code
// This function is testable and doesn't call os.Exit
func Run(args []string, stdout, stderr io.Writer) int {
	// Parse command-line flags
	cliConfig, err := parseFlags(args, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return ExitUsageError
	}

	// Handle --version
	if cliConfig.Version {
		fmt.Fprintf(stdout, "%s version %s\n", AppName, AppVersion)
		return ExitSuccess
	}

	// Handle --help or no input
	if cliConfig.Help || cliConfig.InputFile == "" {
		PrintHelp(stderr)
		if cliConfig.Help {
			return ExitSuccess
		}
		return ExitUsageError
	}

	// Load configuration with priority: default < file < env < flags
	cfg, err := loadConfiguration(cliConfig)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur configuration: %v\n", err)
		return ExitInvalidConfig
	}

	// Set debug mode from CLI flag if provided
	if cliConfig.Debug {
		cfg.SetDebug(true)
	}

	// Parse input file
	result, err := ParseInput(cliConfig.InputFile)
	if err != nil {
		fmt.Fprintf(stderr, "Erreur: %v\n", err)
		return ExitRuntimeError
	}

	// Output according to format
	if err := OutputResult(result, cliConfig.OutputFormat, stdout); err != nil {
		fmt.Fprintf(stderr, "Erreur sortie: %v\n", err)
		return ExitRuntimeError
	}

	return ExitSuccess
}

// parseFlags parse les arguments de ligne de commande
func parseFlags(args []string, errOutput io.Writer) (*CLIConfig, error) {
	fs := flag.NewFlagSet(AppName, flag.ContinueOnError)
	fs.SetOutput(errOutput)

	cfg := &CLIConfig{}

	fs.StringVar(&cfg.ConfigFile, "config", DefaultConfigPath, "Chemin vers le fichier de configuration")
	fs.StringVar(&cfg.OutputFormat, "output", DefaultOutputFormat, "Format de sortie (json)")
	fs.BoolVar(&cfg.Debug, "debug", false, "Activer le mode debug")
	fs.BoolVar(&cfg.Version, "version", false, "Afficher la version")
	fs.BoolVar(&cfg.Help, "help", false, "Afficher l'aide")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// Get positional argument (input file)
	if fs.NArg() > 0 {
		cfg.InputFile = fs.Arg(0)
	}

	return cfg, nil
}

// loadConfiguration charge la configuration avec priorité: défaut < fichier < env < flags
func loadConfiguration(cliConfig *CLIConfig) (*config.ConfigManager, error) {
	// Déterminer le chemin du fichier de config (flag > env > défaut)
	configPath := cliConfig.ConfigFile
	if configPath == "" {
		configPath = config.GetConfigFilePath(DefaultConfigPath)
	}

	// Créer le gestionnaire avec défauts
	cm := config.NewConfigManager(configPath)

	// Charger depuis fichier si spécifié
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			if err := cm.LoadFromFile(); err != nil {
				return nil, fmt.Errorf("chargement fichier config: %w", err)
			}
		}
	}

	// Surcharger avec variables d'environnement
	if err := cm.LoadFromEnv(); err != nil {
		return nil, fmt.Errorf("chargement variables environnement: %w", err)
	}

	// Valider la configuration finale
	if err := cm.Validate(); err != nil {
		return nil, fmt.Errorf("validation config: %w", err)
	}

	return cm, nil
}

// ParseInput lit et parse une entrée (fichier ou stdin)
func ParseInput(inputSource string) (interface{}, error) {
	var input []byte
	var err error
	var filename string

	if inputSource == StdinPlaceholder {
		// Lire depuis stdin
		filename = "stdin"
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("lecture stdin: %w", err)
		}
	} else {
		// Lire depuis fichier
		filename = inputSource
		input, err = os.ReadFile(inputSource)
		if err != nil {
			return nil, fmt.Errorf("lecture fichier: %w", err)
		}
	}

	// Parse the input
	result, err := constraint.ParseConstraint(filename, input)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	// Validate the program
	if err := constraint.ValidateConstraintProgram(result); err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	return result, nil
}

// OutputResult convertit le résultat au format demandé et l'écrit
func OutputResult(result interface{}, format string, w io.Writer) error {
	switch format {
	case "json":
		return OutputJSON(result, w)
	default:
		return fmt.Errorf("format de sortie non supporté: %s", format)
	}
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
	fmt.Fprintf(w, "  %s [options] <input-file>\n", AppName)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Arguments:")
	fmt.Fprintf(w, "  <input-file>    Fichier de contraintes à parser (utilisez '%s' pour stdin)\n", StdinPlaceholder)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  --config PATH   Chemin vers le fichier de configuration")
	fmt.Fprintln(w, "  --output FORMAT Format de sortie (défaut: json)")
	fmt.Fprintln(w, "  --debug         Activer le mode debug")
	fmt.Fprintln(w, "  --version       Afficher la version")
	fmt.Fprintln(w, "  --help          Afficher cette aide")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Variables d'environnement:")
	fmt.Fprintf(w, "  %s   Fichier de configuration\n", config.EnvConfigFile)
	fmt.Fprintf(w, "  %s  Nombre max d'expressions\n", config.EnvMaxExpressions)
	fmt.Fprintf(w, "  %s         Profondeur max de validation\n", config.EnvMaxDepth)
	fmt.Fprintf(w, "  %s              Mode debug (true/false)\n", config.EnvDebug)
	fmt.Fprintf(w, "  %s         Niveau de log (debug/info/warn/error)\n", config.EnvLogLevel)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Exemples:")
	fmt.Fprintf(w, "  %s constraints.tsd\n", AppName)
	fmt.Fprintf(w, "  %s --debug --config myconfig.json constraints.tsd\n", AppName)
	fmt.Fprintf(w, "  cat constraints.tsd | %s %s\n", AppName, StdinPlaceholder)
	fmt.Fprintf(w, "  %s=true %s constraints.tsd\n", config.EnvDebug, AppName)
}

// ParseFile wrapper pour compatibilité avec tests existants
// TODO: Les tests doivent être migrés pour utiliser ParseInput au lieu de ParseFile
func ParseFile(inputFile string) (interface{}, error) {
	return ParseInput(inputFile)
}
