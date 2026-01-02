// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// Action names
const (
	// ActionPrint nom de l'action Print
	ActionPrint = "Print"

	// ActionLog nom de l'action Log
	ActionLog = "Log"

	// ActionUpdate nom de l'action Update
	ActionUpdate = "Update"

	// ActionInsert nom de l'action Insert
	ActionInsert = "Insert"

	// ActionRetract nom de l'action Retract
	ActionRetract = "Retract"

	// ActionXuple nom de l'action Xuple
	ActionXuple = "Xuple"
)

// Argument counts
const (
	// ArgsCountPrint nombre d'arguments attendus pour Print
	ArgsCountPrint = 1

	// ArgsCountLog nombre d'arguments attendus pour Log
	ArgsCountLog = 1

	// ArgsCountUpdate nombre d'arguments attendus pour Update
	ArgsCountUpdate = 1

	// ArgsCountInsert nombre d'arguments attendus pour Insert
	ArgsCountInsert = 1

	// ArgsCountRetract nombre d'arguments attendus pour Retract
	ArgsCountRetract = 1

	// ArgsCountXuple nombre d'arguments attendus pour Xuple
	ArgsCountXuple = 2
)

// Log format
const (
	// LogPrefix préfixe pour les messages de log TSD
	LogPrefix = "[TSD] %s"
)

// BuiltinActionExecutor exécute les actions par défaut du système TSD.
//
// Il s'agit d'une implémentation centralisée de toutes les actions système :
// Print, Log, Update, Insert, Retract, et Xuple.
//
// Thread-Safety:
//   - Les méthodes Execute* sont thread-safe si le réseau RETE l'est
//   - XupleManager doit être thread-safe
type BuiltinActionExecutor struct {
	network      *rete.ReteNetwork
	xupleManager xuples.XupleManager
	output       io.Writer
	logger       *log.Logger
}

// NewBuiltinActionExecutor crée un nouvel exécuteur d'actions natives.
//
// Paramètres:
//   - network: réseau RETE pour les opérations Insert/Update/Retract
//   - xupleManager: gestionnaire de xuples pour l'action Xuple (peut être nil si non utilisé)
//   - output: writer pour l'action Print (utilise os.Stdout si nil)
//   - logger: logger pour l'action Log (utilise log.Default() si nil)
//
// Retourne:
//   - *BuiltinActionExecutor: exécuteur initialisé
func NewBuiltinActionExecutor(network *rete.ReteNetwork, xupleManager xuples.XupleManager, output io.Writer, logger *log.Logger) *BuiltinActionExecutor {
	if output == nil {
		output = os.Stdout
	}
	if logger == nil {
		logger = log.Default()
	}

	return &BuiltinActionExecutor{
		network:      network,
		xupleManager: xupleManager,
		output:       output,
		logger:       logger,
	}
}

// Execute exécute une action par défaut.
//
// Paramètres:
//   - actionName: nom de l'action (Print, Log, Update, Insert, Retract, Xuple)
//   - args: arguments de l'action (déjà évalués)
//   - token: token contenant les faits déclencheurs (utilisé par Xuple)
//
// Retourne:
//   - error: erreur si l'action échoue ou si le nom est inconnu
func (e *BuiltinActionExecutor) Execute(actionName string, args []interface{}, token *rete.Token) error {
	switch actionName {
	case ActionPrint:
		return e.executePrint(args)
	case ActionLog:
		return e.executeLog(args)
	case ActionUpdate:
		return e.executeUpdate(args)
	case ActionInsert:
		return e.executeInsert(args)
	case ActionRetract:
		return e.executeRetract(args)
	case ActionXuple:
		return e.executeXuple(args, token)
	default:
		return fmt.Errorf("unknown builtin action: %s", actionName)
	}
}

// executePrint implémente l'action Print(message: string).
// Affiche une chaîne de caractères sur la sortie standard.
func (e *BuiltinActionExecutor) executePrint(args []interface{}) error {
	if len(args) != ArgsCountPrint {
		return fmt.Errorf("action Print expects %d argument, got %d", ArgsCountPrint, len(args))
	}

	message, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("action Print expects string argument, got %T", args[0])
	}

	_, err := fmt.Fprintln(e.output, message)
	if err != nil {
		return fmt.Errorf("action Print failed: %w", err)
	}

	return nil
}

// executeLog implémente l'action Log(message: string).
// Génère une trace dans le système de logging.
func (e *BuiltinActionExecutor) executeLog(args []interface{}) error {
	if len(args) != ArgsCountLog {
		return fmt.Errorf("action Log expects %d argument, got %d", ArgsCountLog, len(args))
	}

	message, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("action Log expects string argument, got %T", args[0])
	}

	e.logger.Printf(LogPrefix, message)
	return nil
}

// executeUpdate implémente l'action Update(fact: any).
// Modifie un fait existant et met à jour les tokens liés dans RETE.
//
// Cette action utilise la méthode UpdateFact du réseau RETE qui :
// 1. Vérifie que le fait existe
// 2. Rétracte l'ancien fait (propage la suppression)
// 3. Insère le fait mis à jour (propage l'ajout)
// 4. Garantit la cohérence du réseau RETE
func (e *BuiltinActionExecutor) executeUpdate(args []interface{}) error {
	if len(args) != ArgsCountUpdate {
		return fmt.Errorf("action Update expects %d argument, got %d", ArgsCountUpdate, len(args))
	}

	fact, ok := args[0].(*rete.Fact)
	if !ok {
		return fmt.Errorf("action Update expects fact argument, got %T", args[0])
	}

	if fact == nil {
		return fmt.Errorf("action Update: fact is nil")
	}

	// Déléguer au réseau RETE
	return e.network.UpdateFact(fact)
}

// executeInsert implémente l'action Insert(fact: any).
// Crée un nouveau fait et l'insère dans le réseau RETE.
//
// Cette action utilise la méthode InsertFact du réseau RETE qui :
// 1. Valide le fait (type, ID, champs)
// 2. Vérifie qu'il n'existe pas déjà
// 3. L'ajoute au storage
// 4. Propage l'insertion dans le réseau RETE
func (e *BuiltinActionExecutor) executeInsert(args []interface{}) error {
	if len(args) != ArgsCountInsert {
		return fmt.Errorf("action Insert expects %d argument, got %d", ArgsCountInsert, len(args))
	}

	fact, ok := args[0].(*rete.Fact)
	if !ok {
		return fmt.Errorf("action Insert expects fact argument, got %T", args[0])
	}

	if fact == nil {
		return fmt.Errorf("action Insert: fact is nil")
	}

	// Déléguer au réseau RETE
	return e.network.InsertFact(fact)
}

// executeRetract implémente l'action Retract(fact: Fact).
// Supprime un fait du réseau RETE ainsi que tous les tokens liés.
//
// Cette action utilise la méthode RetractFact du réseau RETE qui :
// 1. Extrait l'ID interne du fait (basé sur la clé primaire)
// 2. Vérifie que le fait existe
// 3. Le supprime du storage
// 4. Propage la rétraction dans le réseau RETE
//
// Important: L'argument doit être le fait entier (ex: Retract(p)),
// pas un champ du fait (ex: Retract(p.id) est incorrect).
func (e *BuiltinActionExecutor) executeRetract(args []interface{}) error {
	if len(args) != ArgsCountRetract {
		return fmt.Errorf("action Retract expects %d argument, got %d", ArgsCountRetract, len(args))
	}

	fact, ok := args[0].(*rete.Fact)
	if !ok {
		return fmt.Errorf("action Retract expects Fact argument, got %T (hint: use Retract(fact), not Retract(fact.id))", args[0])
	}

	if fact == nil {
		return fmt.Errorf("action Retract: fact is nil")
	}

	// Extraire l'ID interne du fait (Type~clé_primaire)
	internalID := fact.GetInternalID()
	if internalID == "" {
		return fmt.Errorf("action Retract: fact has empty internal ID")
	}

	// Déléguer au réseau RETE
	return e.network.RetractFact(internalID)
}

// executeXuple implémente l'action Xuple(xuplespace: string, fact: any).
// Crée un xuple dans le xuple-space spécifié.
func (e *BuiltinActionExecutor) executeXuple(args []interface{}, token *rete.Token) error {
	if len(args) != ArgsCountXuple {
		return fmt.Errorf("action Xuple expects %d arguments, got %d", ArgsCountXuple, len(args))
	}

	xuplespace, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("action Xuple expects string as first argument, got %T", args[0])
	}

	if xuplespace == "" {
		return fmt.Errorf("action Xuple: xuplespace is empty")
	}

	fact, ok := args[1].(*rete.Fact)
	if !ok {
		return fmt.Errorf("action Xuple expects fact as second argument, got %T", args[1])
	}

	if fact == nil {
		return fmt.Errorf("action Xuple: fact is nil")
	}

	// Vérifier que le XupleManager est disponible
	if e.xupleManager == nil {
		return fmt.Errorf("action Xuple requires XupleManager to be configured")
	}

	// Extraire les faits déclencheurs du token
	triggeringFacts := e.extractTriggeringFacts(token)

	// Déléguer au XupleManager
	return e.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
}

// extractTriggeringFacts extrait tous les faits d'un token.
//
// Cette méthode parcourt la chaîne de tokens pour extraire tous les faits
// qui ont déclenché l'activation de la règle.
//
// Paramètres:
//   - token: token contenant les faits
//
// Retourne:
//   - []*rete.Fact: liste des faits déclencheurs dans l'ordre chronologique
func (e *BuiltinActionExecutor) extractTriggeringFacts(token *rete.Token) []*rete.Fact {
	if token == nil {
		return []*rete.Fact{}
	}

	var facts []*rete.Fact

	// Parcourir la chaîne de tokens via Parent
	for t := token; t != nil; t = t.Parent {
		// Ajouter tous les faits du token actuel
		if t.Facts != nil {
			facts = append(facts, t.Facts...)
		}
	}

	// Inverser pour avoir l'ordre chronologique (plus ancien → plus récent)
	for i := 0; i < len(facts)/2; i++ {
		facts[i], facts[len(facts)-1-i] = facts[len(facts)-1-i], facts[i]
	}

	return facts
}

// SetOutput permet de changer le writer de sortie pour l'action Print.
// Utile pour les tests.
func (e *BuiltinActionExecutor) SetOutput(output io.Writer) {
	if output != nil {
		e.output = output
	}
}

// SetLogger permet de changer le logger pour l'action Log.
// Utile pour les tests.
func (e *BuiltinActionExecutor) SetLogger(logger *log.Logger) {
	if logger != nil {
		e.logger = logger
	}
}
