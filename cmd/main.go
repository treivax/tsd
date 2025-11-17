package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/treivax/tsd/constraint"
)

func main() {
	var (
		constraintFile = flag.String("constraint", "", "Fichier de contraintes (.constraint)")
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

	if *constraintFile == "" {
		fmt.Fprintf(os.Stderr, "Erreur: fichier constraint requis\n\n")
		printHelp()
		os.Exit(1)
	}

	// Vérifier que le fichier constraint existe
	if _, err := os.Stat(*constraintFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Fichier contrainte non trouvé: %s\n", *constraintFile)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("🚀 TSD - Analyse des contraintes\n")
		fmt.Printf("===============================\n")
		fmt.Printf("Fichier: %s\n\n", *constraintFile)
	}

	// Parser le fichier constraint
	result, err := constraint.ParseConstraintFile(*constraintFile)
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
		fmt.Printf("Le fichier de contraintes est syntaxiquement correct.\n")
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
	fmt.Println("")
	fmt.Println("OPTIONS:")
	fmt.Println("  -constraint <file>  Fichier de règles/contraintes (requis)")
	fmt.Println("  -facts <file>       Fichier de faits (optionnel, pour futur usage)")
	fmt.Println("  -v                  Mode verbeux")
	fmt.Println("  -version            Afficher la version")
	fmt.Println("  -h                  Afficher cette aide")
	fmt.Println("")
	fmt.Println("EXEMPLES:")
	fmt.Println("  tsd -constraint rules.constraint")
	fmt.Println("  tsd -constraint rules.constraint -v")
	fmt.Println("")
	fmt.Println("FORMATS DE FICHIERS:")
	fmt.Println("  .constraint : Règles en syntaxe TSD")
	fmt.Println("  .facts      : Faits en format structuré (support futur)")
}
