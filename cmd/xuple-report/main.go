package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/rete/actions"
	"github.com/treivax/tsd/xuples"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: xuple-report <fichier.tsd>")
		os.Exit(1)
	}

	filename := os.Args[1]

	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	fmt.Println("  RAPPORT D'EXÉCUTION E2E - SYSTÈME XUPLE-SPACE")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("📁 Fichier: %s\n", filename)
	fmt.Printf("⏰ Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Lire le fichier
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("❌ Erreur lecture fichier: %v", err)
	}

	// Parser le programme
	fmt.Println("───────────────────────────────────────────────────────────────────────────")
	fmt.Println("ÉTAPE 1 : PARSING DU PROGRAMME")
	fmt.Println("───────────────────────────────────────────────────────────────────────────")

	program, err := constraint.ParseProgram(string(content), filename)
	if err != nil {
		log.Fatalf("❌ Erreur parsing: %v", err)
	}

	fmt.Printf("✅ Parsing réussi\n\n")

	// Afficher les types
	fmt.Println("📋 TYPES DÉFINIS:")
	for i, typ := range program.Types {
		fmt.Printf("  %d. %s\n", i+1, typ.Name)
		fmt.Printf("     Champs: ")
		fields := []string{}
		for _, field := range typ.Fields {
			prefix := ""
			if field.IsPrimaryKey {
				prefix = "#"
			}
			fields = append(fields, fmt.Sprintf("%s%s: %s", prefix, field.Name, field.Type))
		}
		fmt.Printf("%s\n", strings.Join(fields, ", "))
	}
	fmt.Println()

	// Afficher les xuple-spaces
	fmt.Println("🗄️  XUPLE-SPACES DÉCLARÉS:")
	for i, xs := range program.XupleSpaces {
		fmt.Printf("  %d. %s\n", i+1, xs.Name)
		fmt.Printf("     • Sélection: %s\n", xs.SelectionPolicy)
		fmt.Printf("     • Consommation: %s", xs.ConsumptionPolicy.Type)
		if xs.ConsumptionPolicy.Limit > 0 {
			fmt.Printf(" (limite: %d)", xs.ConsumptionPolicy.Limit)
		}
		fmt.Println()
		fmt.Printf("     • Rétention: %s", xs.RetentionPolicy.Type)
		if xs.RetentionPolicy.Duration > 0 {
			fmt.Printf(" (durée: %d secondes)", xs.RetentionPolicy.Duration)
		}
		fmt.Println()
	}
	fmt.Println()

	// Afficher les actions
	fmt.Println("⚡ ACTIONS DÉFINIES:")
	for i, action := range program.Actions {
		fmt.Printf("  %d. %s(", i+1, action.Name)
		params := []string{}
		for _, param := range action.Parameters {
			params = append(params, fmt.Sprintf("%s: %s", param.Name, param.Type))
		}
		fmt.Printf("%s)\n", strings.Join(params, ", "))
	}
	fmt.Println()

	// Afficher les règles
	fmt.Println("📜 RÈGLES DÉFINIES:")
	for i, expr := range program.Expressions {
		if expr.Type == "expression" {
			fmt.Printf("  %d. %s\n", i+1, expr.Name)

			// Afficher les patterns
			if len(expr.Patterns) > 0 {
				fmt.Printf("     Patterns: ")
				patterns := []string{}
				for _, p := range expr.Patterns {
					patterns = append(patterns, fmt.Sprintf("{%s: %s}", p.Alias, p.Type))
				}
				fmt.Printf("%s\n", strings.Join(patterns, ", "))
			}

			// Afficher les actions
			if len(expr.Actions) > 0 {
				fmt.Printf("     Actions: ")
				actionNames := []string{}
				for _, a := range expr.Actions {
					actionNames = append(actionNames, a.Name)
				}
				fmt.Printf("%s\n", strings.Join(actionNames, ", "))
			}
		}
	}
	fmt.Println()

	// Afficher les faits
	fmt.Println("📊 FAITS INJECTÉS:")
	factsByType := make(map[string][]constraint.Fact)
	for _, fact := range program.Facts {
		factsByType[fact.Type] = append(factsByType[fact.Type], fact)
	}

	for typeName, facts := range factsByType {
		fmt.Printf("  %s (%d fait(s)):\n", typeName, len(facts))
		for i, fact := range facts {
			fmt.Printf("    %d. ", i+1)
			fields := []string{}
			for key, value := range fact.Fields {
				fields = append(fields, fmt.Sprintf("%s: %v", key, formatValue(value)))
			}
			fmt.Printf("%s\n", strings.Join(fields, ", "))
		}
	}
	fmt.Println()

	// Créer le réseau RETE
	fmt.Println("───────────────────────────────────────────────────────────────────────────")
	fmt.Println("ÉTAPE 2 : CRÉATION DU RÉSEAU RETE ET EXÉCUTION")
	fmt.Println("───────────────────────────────────────────────────────────────────────────")

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	xupleManager := xuples.NewXupleManager()
	executor := actions.NewBuiltinActionExecutor(network, xupleManager, os.Stdout, log.Default())

	// Enregistrer l'exécuteur d'actions
	network.SetBuiltinActionExecutor(executor)

	// Créer les xuple-spaces
	for _, xs := range program.XupleSpaces {
		policies := xuples.XupleSpacePolicies{
			Selection: xs.SelectionPolicy,
			Consumption: xuples.ConsumptionPolicy{
				Type:  xs.ConsumptionPolicy.Type,
				Limit: xs.ConsumptionPolicy.Limit,
			},
			Retention: xuples.RetentionPolicy{
				Type:     xs.RetentionPolicy.Type,
				Duration: time.Duration(xs.RetentionPolicy.Duration) * time.Second,
			},
		}
		err := xupleManager.CreateXupleSpace(xs.Name, policies)
		if err != nil {
			log.Fatalf("❌ Erreur création xuple-space %s: %v", xs.Name, err)
		}
		fmt.Printf("✅ Xuple-space créé: %s\n", xs.Name)
	}
	fmt.Println()

	// Ingérer le programme
	fmt.Println("🔄 Ingestion du programme dans le réseau RETE...")
	err = network.IngestTSDFile(filename)
	if err != nil {
		log.Fatalf("❌ Erreur ingestion: %v", err)
	}
	fmt.Println("✅ Programme ingéré avec succès")
	fmt.Println()

	// Afficher les statistiques du réseau
	fmt.Println("───────────────────────────────────────────────────────────────────────────")
	fmt.Println("ÉTAPE 3 : RÉSULTATS DE L'EXÉCUTION")
	fmt.Println("───────────────────────────────────────────────────────────────────────────")

	fmt.Println("📈 STATISTIQUES DU RÉSEAU RETE:")
	stats := network.Statistics()
	fmt.Printf("  • Faits dans le working memory: %d\n", stats.FactCount)
	fmt.Printf("  • Activations générées: %d\n", stats.ActivationCount)
	fmt.Printf("  • TypeNodes: %d\n", stats.TypeNodeCount)
	fmt.Printf("  • AlphaNodes: %d\n", stats.AlphaNodeCount)
	fmt.Printf("  • BetaNodes: %d\n", stats.BetaNodeCount)
	fmt.Printf("  • TerminalNodes: %d\n", stats.TerminalNodeCount)
	fmt.Println()

	// Afficher les xuples créés
	fmt.Println("🎯 XUPLES GÉNÉRÉS DANS LES XUPLE-SPACES:")
	fmt.Println()

	totalXuples := 0
	for _, xs := range program.XupleSpaces {
		xupleList, err := xupleManager.ListXuples(xs.Name)
		if err != nil {
			fmt.Printf("⚠️  Erreur lecture xuple-space %s: %v\n", xs.Name, err)
			continue
		}

		fmt.Printf("📦 Xuple-space: %s\n", xs.Name)
		fmt.Printf("   Politique: selection=%s, consumption=%s, retention=%s\n",
			xs.SelectionPolicy, xs.ConsumptionPolicy.Type, xs.RetentionPolicy.Type)
		fmt.Printf("   Nombre de xuples: %d\n", len(xupleList))

		if len(xupleList) > 0 {
			for i, xuple := range xupleList {
				fmt.Printf("\n   Xuple #%d:\n", i+1)
				fmt.Printf("     ID: %s\n", xuple.ID)
				fmt.Printf("     Fait: Type=%s, ID=%s\n", xuple.Fact.Type, xuple.Fact.ID)
				fmt.Printf("     Champs:\n")
				for key, value := range xuple.Fact.Fields {
					fmt.Printf("       • %s: %v\n", key, formatValue(value))
				}
				fmt.Printf("     Créé: %s\n", xuple.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("     Consommé: %v\n", xuple.IsConsumed())
				if len(xuple.TriggeringFacts) > 0 {
					fmt.Printf("     Faits déclencheurs: %d\n", len(xuple.TriggeringFacts))
					for j, tf := range xuple.TriggeringFacts {
						fmt.Printf("       %d. %s (ID: %s)\n", j+1, tf.Type, tf.ID)
					}
				}
			}
		}
		fmt.Println()
		totalXuples += len(xupleList)
	}

	// Résumé final
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	fmt.Println("RÉSUMÉ FINAL")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	fmt.Printf("✓ Types définis: %d\n", len(program.Types))
	fmt.Printf("✓ Xuple-spaces déclarés: %d\n", len(program.XupleSpaces))
	fmt.Printf("✓ Actions définies: %d\n", len(program.Actions))
	fmt.Printf("✓ Règles définies: %d\n", len(program.Expressions))
	fmt.Printf("✓ Faits injectés: %d\n", len(program.Facts))
	fmt.Printf("✓ Activations générées: %d\n", stats.ActivationCount)
	fmt.Printf("✓ Xuples créés: %d\n", totalXuples)
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%.2f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "null"
	default:
		data, _ := json.Marshal(v)
		return string(data)
	}
}
