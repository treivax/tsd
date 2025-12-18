// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package xuples implémente le système de xuple-spaces pour TSD.
//
// Un xuple est un tuple étendu contenant un fait principal et les faits
// déclencheurs qui ont conduit à sa création. Les xuples sont stockés dans
// des xuple-spaces qui appliquent des politiques de sélection, consommation
// et rétention.
//
// # Architecture
//
// Le module xuples est composé de plusieurs éléments :
//
//   - Xuple : tuple étendu avec fait principal et faits déclencheurs
//   - XupleSpace : espace nommé gérant des xuples avec politiques
//   - XupleManager : gestionnaire global de multiples xuple-spaces
//   - Policies : stratégies configurables (sélection, consommation, rétention)
//
// # Découplage RETE
//
// Ce package est totalement découplé du moteur RETE. Il ne dépend que du
// type rete.Fact et peut être utilisé indépendamment pour tout système
// nécessitant une gestion de tuples étendus avec politiques.
//
// # Thread-Safety
//
// Toutes les opérations sont thread-safe :
//   - Synchronisation via sync.RWMutex
//   - Génération d'IDs unique thread-safe
//   - Opérations atomiques pour compteurs
//
// # Exemple d'Utilisation
//
//	// Créer un manager
//	manager := xuples.NewXupleManager()
//
//	// Configurer les politiques
//	config := xuples.XupleSpaceConfig{
//	    Name:              "alerts",
//	    SelectionPolicy:   xuples.NewFIFOSelectionPolicy(),
//	    ConsumptionPolicy: xuples.NewOnceConsumptionPolicy(),
//	    RetentionPolicy:   xuples.NewUnlimitedRetentionPolicy(),
//	}
//
//	// Créer un xuple-space
//	err := manager.CreateXupleSpace("alerts", config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Récupérer le xuple-space
//	space, err := manager.GetXupleSpace("alerts")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Créer un xuple
//	fact := &rete.Fact{ID: "f1", Type: "Alert"}
//	triggering := []*rete.Fact{{ID: "t1", Type: "Trigger"}}
//	err = manager.CreateXuple("alerts", fact, triggering)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Récupérer et consommer un xuple
//	xuple, err := space.Retrieve("agent-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = space.MarkConsumed(xuple.ID, "agent-1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Nettoyer les xuples expirés
//	cleaned := space.Cleanup()
//	log.Printf("Cleaned %d xuples", cleaned)
//
// # Politiques
//
// Le module fournit trois types de politiques configurables :
//
// Politiques de Sélection :
//   - FIFO : First-In-First-Out (plus ancien d'abord)
//   - LIFO : Last-In-First-Out (plus récent d'abord)
//   - Random : Sélection aléatoire
//
// Politiques de Consommation :
//   - Once : Une seule consommation au total
//   - PerAgent : Une consommation par agent
//   - Limited : Nombre limité de consommations
//
// Politiques de Rétention :
//   - Unlimited : Conservation illimitée
//   - Duration : Expiration après une durée
package xuples
