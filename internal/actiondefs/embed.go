// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package actiondefs fournit les définitions des actions par défaut du système TSD.
//
// Ce package intermédiaire existe pour résoudre les dépendances circulaires :
// - Il contient uniquement le fichier defaults.tsd embarqué
// - Il peut être importé par constraint (pour validation) et defaultactions (pour chargement)
// - Il ne dépend d'aucun autre package interne
package actiondefs

import (
	_ "embed"
)

// DefaultActionsTSD contient le contenu du fichier defaults.tsd embarqué dans le binaire.
// Ce fichier définit les 6 actions système : Print, Log, Update, Insert, Retract, Xuple.
//
//go:embed defaults.tsd
var DefaultActionsTSD string

// DefaultActionNames liste les noms de toutes les actions par défaut.
// Cette liste est utilisée pour la validation et les vérifications de complétude.
var DefaultActionNames = []string{
	"Print",
	"Log",
	"Update",
	"Insert",
	"Retract",
	"Xuple",
}

// IsDefaultAction vérifie si un nom correspond à une action par défaut.
//
// Paramètres:
//   - name: nom de l'action à vérifier
//
// Retourne:
//   - true si l'action fait partie des actions système
//   - false sinon
func IsDefaultAction(name string) bool {
	for _, defaultName := range DefaultActionNames {
		if name == defaultName {
			return true
		}
	}
	return false
}
