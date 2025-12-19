// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

/*
Package api fournit une interface simplifiée pour utiliser le moteur de règles TSD.

Ce package est le point d'entrée recommandé pour toutes les applications utilisant TSD.
Il intègre automatiquement les packages rete, xuples, et constraint, et gère toute
la configuration nécessaire.

# Utilisation Basique

La manière la plus simple d'utiliser TSD est via le Pipeline:

import "github.com/treivax/tsd/api"

func main() {
// Créer un pipeline
pipeline := api.NewPipeline()

// Ingérer un programme TSD
result, err := pipeline.IngestFile("program.tsd")
if err != nil {
log.Fatal(err)
}

// Utiliser les résultats
fmt.Printf("Types définis: %d\n", result.TypeCount())
fmt.Printf("Règles actives: %d\n", result.RuleCount())
fmt.Printf("Faits dans le réseau: %d\n", result.FactCount())
fmt.Printf("Xuple-spaces créés: %d\n", result.XupleSpaceCount())
}

# Accès aux Xuples

Les xuples créés par les règles sont accessibles via le résultat:

result, _ := pipeline.IngestFile("monitoring.tsd")

// Récupérer tous les xuples d'un xuple-space
alerts := result.GetXuples("critical_alerts")
for _, xuple := range alerts {
fmt.Printf("Alert: %v\n", xuple.Fact.Fields)
}

// Consommer un xuple (retrieve)
xuple, err := result.Retrieve("critical_alerts", "agent1")
if err == nil {
fmt.Printf("Consumed: %v\n", xuple.Fact.Fields)
}

# Configuration Avancée

Pour une configuration personnalisée:

config := &api.Config{
LogLevel:          api.LogLevelDebug,
EnableMetrics:     true,
MaxFactsInMemory:  100000,
XupleSpaceDefaults: &api.XupleSpaceDefaults{
Selection:   api.SelectionFIFO,
Consumption: api.ConsumptionOnce,
Retention:   api.RetentionUnlimited,
},
}

pipeline := api.NewPipelineWithConfig(config)
result, err := pipeline.IngestFile("program.tsd")

# Ingestion Incrémentale

Le pipeline supporte l'ingestion incrémentale de plusieurs fichiers:

pipeline := api.NewPipeline()

// Charger les types et règles
_, err := pipeline.IngestFile("types.tsd")
if err != nil {
log.Fatal(err)
}

// Ajouter plus de règles
_, err = pipeline.IngestFile("additional-rules.tsd")
if err != nil {
log.Fatal(err)
}

// Soumettre des faits
result, err := pipeline.IngestFile("facts.tsd")

Le réseau RETE est étendu de manière incrémentale et tous les faits précédents
sont automatiquement propagés aux nouvelles règles.

# Thread Safety

Le Pipeline est thread-safe. Plusieurs goroutines peuvent appeler IngestFile
en parallèle, mais notez que l'ordre d'exécution des règles peut varier.

Pour un contrôle strict de l'ordre, utilisez un seul goroutine ou synchronisez
explicitement les appels.

# Métriques

Les métriques d'ingestion sont disponibles dans le résultat:

result, _ := pipeline.IngestFile("program.tsd")
metrics := result.Metrics()

fmt.Printf("Temps de parsing: %v\n", metrics.ParseDuration)
fmt.Printf("Temps de construction réseau: %v\n", metrics.BuildDuration)
fmt.Printf("Nombre de propagations: %d\n", metrics.PropagationCount)

# Gestion d'Erreurs

Les erreurs sont détaillées et incluent la position dans le fichier source:

_, err := pipeline.IngestFile("invalid.tsd")
if err != nil {
if parseErr, ok := err.(*api.ParseError); ok {
fmt.Printf("Erreur ligne %d, colonne %d: %s\n",
parseErr.Line, parseErr.Column, parseErr.Message)
}
}

# Architecture

Le package api est construit au-dessus de:
  - constraint: Parser TSD (PEG)
  - rete: Moteur de règles (algorithme RETE)
  - xuples: Gestion des xuple-spaces et xuples

Il gère automatiquement:
  - Création du réseau RETE
  - Initialisation du XupleManager
  - Création des xuple-spaces à partir des définitions
  - Enregistrement des actions (Xuple, Print, etc.)
  - Configuration des handlers
  - Propagation des faits

L'utilisateur n'a besoin de connaître aucun détail d'implémentation.
*/
package api
