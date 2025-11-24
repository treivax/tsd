package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/treivax/tsd/constraint"
)

func main() {
	var (
		constraintFile = flag.String("constraint", "", "Fichier de contraintes (.constraint)")
		constraintText = flag.String("text", "", "Texte de contrainte directement (alternative à -constraint)")
		stdin          = flag.Bool("stdin", false, "Lire les contraintes depuis stdin")
		factsFile      = flag.String("facts", "", "Fichier de faits (.facts)")
		verbose        = flag.Bool("v", false, "Mode verbeux")
		version        = flag.Bool("version", false, "Afficher la version")
		help           = flag.Bool("h", false, "Afficher l'aide")
	)

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *version {
		fmt.Println("TSD (Type System Development) v1.0")
		fmt.Println("Moteur de règles basé sur l'algorithme RETE")
		return
	}

	// Compter les sources d'entrée
	sourcesCount := 0
	if *constraintFile != "" {
		sourcesCount++
	}
	if *constraintText != "" {
		sourcesCount++
	}
	if *stdin {
		sourcesCount++
	}

	if sourcesCount == 0 {
		fmt.Fprintf(os.Stderr, "Erreur: spécifiez une source (-constraint, -text, ou -stdin)\n\n")
		printHelp()
		os.Exit(1)
	}

	if sourcesCount > 1 {
		fmt.Fprintf(os.Stderr, "Erreur: spécifiez une seule source d'entrée\n\n")
		printHelp()
		os.Exit(1)
	}

	var result interface{}
	var err error
	var sourceName string

	if *stdin {
		// Lire depuis stdin
		sourceName = "<stdin>"
		if *verbose {
			fmt.Printf("🚀 TSD - Analyse des contraintes\n")
			fmt.Printf("===============================\n")
			fmt.Printf("Source: stdin\n\n")
		}
		stdinContent, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "Erreur lecture stdin: %v\n", readErr)
			os.Exit(1)
		}
		result, err = constraint.ParseConstraint(sourceName, stdinContent)
	} else if *constraintText != "" {
		// Parser du texte directement
		sourceName = "<text>"
		if *verbose {
			fmt.Printf("🚀 TSD - Analyse des contraintes\n")
			fmt.Printf("===============================\n")
			fmt.Printf("Source: texte direct\n\n")
		}
		result, err = constraint.ParseConstraint(sourceName, []byte(*constraintText))
	} else {
		// Parser depuis un fichier
		sourceName = *constraintFile
		// Vérifier que le fichier constraint existe
		if _, statErr := os.Stat(*constraintFile); os.IsNotExist(statErr) {
			fmt.Fprintf(os.Stderr, "Fichier contrainte non trouvé: %s\n", *constraintFile)
			os.Exit(1)
		}

		if *verbose {
			fmt.Printf("🚀 TSD - Analyse des contraintes\n")
			fmt.Printf("===============================\n")
			fmt.Printf("Fichier: %s\n\n", *constraintFile)
		}

		result, err = constraint.ParseConstraintFile(*constraintFile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur de parsing: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("✅ Parsing réussi\n")
		fmt.Printf("📋 Validation du programme...\n")
	}

	// Valider le programme
	if err := constraint.ValidateConstraintProgram(result); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur de validation: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Contraintes validées avec succès\n")

	if *verbose {
		fmt.Printf("\n🎉 Analyse terminée!\n")
		fmt.Printf("Les contraintes sont syntaxiquement correctes.\n")
	}

	// TODO: Intégration avec le moteur RETE pour l'exécution complète
	if *factsFile != "" {
		fmt.Printf("ℹ️ Fichier faits spécifié: %s (intégration RETE à venir)\n", *factsFile)
	}
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
