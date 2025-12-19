// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// XupleAction implémente l'action Xuple(xuplespace: string, fact: any).
//
// Cette action crée un xuple dans le xuple-space spécifié en utilisant
// le handler configuré dans le réseau RETE via SetXupleHandler().
//
// Enregistrement automatique :
//   - L'action est automatiquement enregistrée par ActionExecutor.RegisterDefaultActions()
//   - Elle est disponible dès qu'un XupleHandler est configuré via network.SetXupleHandler()
//   - Aucune configuration manuelle n'est nécessaire
//
// Utilisation dans une règle TSD :
//
//	rule HighTemperature {
//	    when {
//	        t: Temperature(value > 30.0)
//	    }
//	    then {
//	        Xuple("alerts", Alert(
//	            sensorId: t.sensorId,
//	            message: "High temperature detected",
//	            temp: t.value
//	        ))
//	    }
//	}
//
// Thread-Safety:
//   - Cette action est thread-safe si le XupleManager sous-jacent l'est
//   - L'exécution est isolée par token (contexte d'exécution)
//
// Arguments attendus:
//   - xuplespace (string): nom du xuple-space cible (doit exister)
//   - fact (*Fact): fait à stocker dans le xuple (peut être un fait inline)
//
// Erreurs:
//   - XupleHandler non configuré : levée si aucun handler n'est défini
//   - Xuple-space inexistant : levée si le space n'a pas été créé
//   - Validation du fait : levée si le type du fait n'est pas défini
type XupleAction struct {
	network *ReteNetwork
}

// NewXupleAction crée un nouveau handler pour l'action Xuple.
//
// Paramètres:
//   - network: réseau RETE contenant le XupleHandler configuré
//
// Retourne:
//   - *XupleAction: handler initialisé
func NewXupleAction(network *ReteNetwork) *XupleAction {
	return &XupleAction{
		network: network,
	}
}

// GetName retourne le nom de l'action.
func (a *XupleAction) GetName() string {
	return "Xuple"
}

// Validate valide les arguments de l'action Xuple.
//
// Attend exactement 2 arguments:
//   - arg[0]: string (nom du xuple-space, non vide)
//   - arg[1]: *Fact (fait à stocker, non nil)
//
// Retourne :
//   - error: erreur si les arguments sont invalides
func (a *XupleAction) Validate(args []interface{}) error {
	if len(args) != 2 {
		return fmt.Errorf(
			"❌ Action Xuple: nombre d'arguments incorrect\n"+
				"   Attendu: 2 (xuplespace, fact)\n"+
				"   Reçu: %d\n"+
				"   Syntaxe: Xuple(\"space_name\", FactType(...))",
			len(args),
		)
	}

	// Valider le premier argument (xuplespace)
	xuplespace, ok := args[0].(string)
	if !ok {
		return fmt.Errorf(
			"❌ Action Xuple: premier argument invalide\n"+
				"   Attendu: string (nom du xuple-space)\n"+
				"   Reçu: %T\n"+
				"   Syntaxe: Xuple(\"space_name\", ...)",
			args[0],
		)
	}

	if xuplespace == "" {
		return fmt.Errorf(
			"❌ Action Xuple: nom de xuple-space vide\n" +
				"   Le nom du xuple-space ne peut pas être vide\n" +
				"   Syntaxe: Xuple(\"space_name\", ...)",
		)
	}

	// Valider le second argument (fact)
	fact, ok := args[1].(*Fact)
	if !ok {
		return fmt.Errorf(
			"❌ Action Xuple: second argument invalide\n"+
				"   Attendu: *Fact (fait ou fait inline)\n"+
				"   Reçu: %T\n"+
				"   Syntaxe: Xuple(\"space\", FactType(...))",
			args[1],
		)
	}

	if fact == nil {
		return fmt.Errorf(
			"❌ Action Xuple: fait nil\n" +
				"   Le fait à stocker ne peut pas être nil\n" +
				"   Syntaxe: Xuple(\"space\", FactType(...))",
		)
	}

	return nil
}

// Execute exécute l'action Xuple.
//
// Cette méthode:
//  1. Valide les arguments (xuplespace et fact)
//  2. Vérifie que le handler Xuple est configuré
//  3. Extrait les faits déclencheurs du contexte d'exécution
//  4. Appelle le handler Xuple configuré dans le réseau
//
// Paramètres :
//   - args: arguments évalués [xuplespace, fact]
//   - ctx: contexte d'exécution contenant le token avec les faits déclencheurs
//
// Retourne :
//   - error: erreur si la création du xuple échoue
//
// Erreurs possibles :
//   - Arguments invalides (nombre, type)
//   - Handler Xuple non configuré
//   - Xuple-space inexistant
//   - Échec de la création du xuple
func (a *XupleAction) Execute(args []interface{}, ctx *ExecutionContext) error {
	// Validation (déjà faite par Validate, mais on sécurise)
	if len(args) != 2 {
		return fmt.Errorf("❌ Action Xuple: attend 2 arguments (xuplespace, fact), reçu %d", len(args))
	}

	xuplespace, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("❌ Action Xuple: le premier argument doit être un string (nom du xuple-space), reçu %T", args[0])
	}

	fact, ok := args[1].(*Fact)
	if !ok {
		return fmt.Errorf("❌ Action Xuple: le second argument doit être un fait (*Fact), reçu %T", args[1])
	}

	// Vérifier que le handler Xuple est configuré
	handler := a.network.GetXupleHandler()
	if handler == nil {
		return fmt.Errorf(
			"❌ Action Xuple non configurée:\n" +
				"   L'action Xuple requiert un XupleHandler configuré dans le réseau RETE.\n" +
				"   Assurez-vous que le pipeline API est correctement initialisé.\n" +
				"   Utilisez: network.SetXupleHandler(handler)",
		)
	}

	// Extraire les faits déclencheurs du contexte
	triggeringFacts := a.extractTriggeringFacts(ctx)

	// Appeler le handler pour créer le xuple
	if err := handler(xuplespace, fact, triggeringFacts); err != nil {
		return fmt.Errorf(
			"❌ Échec création xuple dans '%s':\n"+
				"   Type du fait: %s\n"+
				"   Erreur: %w\n"+
				"   Vérifiez que le xuple-space '%s' existe et est correctement configuré",
			xuplespace, fact.Type, err, xuplespace,
		)
	}

	return nil
}

// extractTriggeringFacts extrait tous les faits du token dans le contexte d'exécution.
//
// Cette méthode parcourt la chaîne de tokens pour collecter tous les faits
// qui ont déclenché l'activation de la règle.
//
// Paramètres:
//   - ctx: contexte d'exécution contenant le token
//
// Retourne:
//   - []*Fact: liste des faits déclencheurs (peut être vide)
func (a *XupleAction) extractTriggeringFacts(ctx *ExecutionContext) []*Fact {
	if ctx == nil {
		return []*Fact{}
	}

	// Accéder au token via le champ (privé, donc on doit utiliser une approche indirecte)
	// Pour l'instant, on retourne une liste vide
	// TODO: Ajouter une méthode GetToken() ou GetFacts() dans ExecutionContext

	// Approche temporaire: essayer d'obtenir les bindings
	bindings := ctx.GetBindings()
	if bindings == nil {
		return []*Fact{}
	}

	facts := make([]*Fact, 0)
	current := bindings
	for current != nil {
		if current.Fact != nil {
			facts = append(facts, current.Fact)
		}
		current = current.Parent
	}

	return facts
}
