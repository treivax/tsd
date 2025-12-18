# Prompt 04 - Implémentation du système d'actions par défaut

## 🎯 Objectif

Implémenter le système d'actions par défaut (Print, Log, Update, Insert, Retract, Xuple) de manière configurable et non hardcodée.

Les actions par défaut doivent :
- Être chargées automatiquement à l'initialisation
- Pouvoir être redéfinies (avec erreur si tentative)
- Ne pas être hardcodées dans le code
- Être implémentées via un fichier de définition parsé

## 📋 Tâches

### 1. Analyser le système actuel de déclaration d'actions

**Objectif** : Comprendre comment les actions sont actuellement déclarées et gérées.

- [ ] Examiner la commande `action` dans le parser
- [ ] Comprendre la structure AST des actions
- [ ] Analyser le registre d'actions dans le compilateur
- [ ] Identifier la validation des doublons
- [ ] Comprendre l'interface ActionExecutor

**Livrables** :
- Créer `tsd/docs/xuples/implementation/03-current-action-system.md` documentant :
  - Syntaxe actuelle de la commande `action`
  - Structure AST de ActionDeclaration
  - Registre des actions dans le compilateur
  - Mécanisme de validation des doublons
  - Interface d'exécution des actions

### 2. Concevoir le système de chargement d'actions par défaut

**Objectif** : Définir comment les actions par défaut sont chargées sans hardcoding.

**Approche recommandée** :
1. Créer un fichier `tsd/internal/defaultactions/defaults.tsd` contenant les définitions
2. Parser ce fichier à l'initialisation du compilateur
3. Enregistrer les actions comme si elles avaient été parsées depuis l'entrée utilisateur
4. Marquer les actions par défaut pour détecter les redéfinitions

**Fichier defaults.tsd attendu** :
```tsd
// Actions par défaut du système TSD
// Ces actions sont disponibles automatiquement dans tous les programmes

action Print(message: string) {
    // Implémentation native
}

action Log(message: string) {
    // Implémentation native
}

action Update(fact: any) {
    // Implémentation native
}

action Insert(fact: any) {
    // Implémentation native
}

action Retract(id: string) {
    // Implémentation native
}

action Xuple(xuplespace: string, fact: any) {
    // Implémentation native
}
```

**Livrables** :
- Créer `tsd/docs/xuples/implementation/04-default-actions-design.md` contenant :
  - Stratégie de chargement des actions par défaut
  - Format du fichier defaults.tsd
  - Mécanisme de marquage des actions par défaut
  - Détection des redéfinitions
  - Gestion des erreurs
  - Diagramme de séquence du chargement

### 3. Créer le fichier de définitions par défaut

**Objectif** : Créer le fichier TSD définissant les actions par défaut.

- [ ] Créer `tsd/internal/defaultactions/defaults.tsd`
- [ ] Définir la signature de chaque action
- [ ] Ajouter des commentaires explicatifs
- [ ] Assurer la cohérence avec la spécification

**Fichier à créer** :
- `tsd/internal/defaultactions/defaults.tsd`

**Contenu complet attendu** :
```tsd
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// ============================================================================
// ACTIONS PAR DÉFAUT DU SYSTÈME TSD
// ============================================================================
//
// Ces actions sont automatiquement disponibles dans tous les programmes TSD.
// Elles ne nécessitent pas de déclaration explicite via la commande 'action'.
//
// Toute tentative de redéfinition de ces actions provoquera une erreur de
// compilation, exactement comme pour toute action déclarée deux fois.
//
// ============================================================================

// Print affiche une chaîne de caractères sur la sortie standard
// Paramètres:
//   - message: la chaîne à afficher
action Print(message: string) {
    // Implémentation native (voir rete/actions/builtin.go)
}

// Log génère une trace dans le système de logging
// Paramètres:
//   - message: la chaîne à tracer
action Log(message: string) {
    // Implémentation native (voir rete/actions/builtin.go)
}

// Update modifie un fait existant et met à jour les tokens liés dans RETE
// Paramètres:
//   - fact: le fait à modifier (doit exister dans le réseau)
// Notes:
//   - Déclenche la propagation des mises à jour dans le réseau RETE
//   - Le fait doit avoir le même type qu'un fait existant
action Update(fact: any) {
    // Implémentation native (voir rete/actions/builtin.go)
}

// Insert crée un nouveau fait et l'insère dans le réseau RETE
// Paramètres:
//   - fact: le nouveau fait à créer
// Notes:
//   - Le fait est propagé dans le réseau RETE
//   - Peut déclencher l'activation de nouvelles règles
action Insert(fact: any) {
    // Implémentation native (voir rete/actions/builtin.go)
}

// Retract supprime un fait du réseau RETE ainsi que tous les tokens liés
// Paramètres:
//   - id: l'identifiant du fait à supprimer
// Notes:
//   - Tous les tokens dépendant de ce fait sont invalidés
//   - La suppression se propage dans tout le réseau
action Retract(id: string) {
    // Implémentation native (voir rete/actions/builtin.go)
}

// Xuple crée un xuple dans le xuple-space spécifié
// Paramètres:
//   - xuplespace: nom du xuple-space cible
//   - fact: le fait principal du xuple
// Notes:
//   - Les faits déclencheurs sont automatiquement extraits du token
//   - Le xuple-space doit avoir été déclaré via 'xuple-space'
//   - Le xuple est soumis aux politiques du xuple-space
action Xuple(xuplespace: string, fact: any) {
    // Implémentation native (voir rete/actions/builtin.go)
}
```

**Livrables** :
- [ ] Fichier defaults.tsd créé avec copyright
- [ ] Commentaires complets et clairs
- [ ] Signatures cohérentes avec la spec

### 4. Implémenter le chargement des actions par défaut

**Objectif** : Charger et enregistrer les actions par défaut à l'initialisation.

- [ ] Créer le package `tsd/internal/defaultactions`
- [ ] Implémenter la fonction de chargement
- [ ] Marquer les actions comme "par défaut"
- [ ] Intégrer dans l'initialisation du compilateur
- [ ] Gérer les erreurs de chargement

**Fichier à créer** :
- `tsd/internal/defaultactions/loader.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package defaultactions

import (
    _ "embed"
    "fmt"
    
    "tsd/parser"
    "tsd/parser/ast"
)

// defaults.tsd est embarqué dans le binaire via go:embed
//go:embed defaults.tsd
var defaultActionsTSD string

// DefaultActionNames contient les noms de toutes les actions par défaut
var DefaultActionNames = []string{
    "Print",
    "Log",
    "Update",
    "Insert",
    "Retract",
    "Xuple",
}

// LoadDefaultActions parse le fichier defaults.tsd et retourne les actions
func LoadDefaultActions() ([]*ast.ActionDeclaration, error) {
    // Parser le fichier embarqué
    result, err := parser.ParseTSD(defaultActionsTSD)
    if err != nil {
        return nil, fmt.Errorf("failed to parse default actions: %w", err)
    }
    
    // Vérifier que toutes les actions attendues sont présentes
    if len(result.Actions) != len(DefaultActionNames) {
        return nil, fmt.Errorf("expected %d default actions, got %d",
            len(DefaultActionNames), len(result.Actions))
    }
    
    // Marquer chaque action comme "par défaut"
    for _, action := range result.Actions {
        action.IsDefault = true // Nouveau champ dans ActionDeclaration
    }
    
    return result.Actions, nil
}

// IsDefaultAction vérifie si un nom correspond à une action par défaut
func IsDefaultAction(name string) bool {
    for _, defaultName := range DefaultActionNames {
        if name == defaultName {
            return true
        }
    }
    return false
}
```

**Livrables** :
- [ ] Package defaultactions créé
- [ ] Fonction LoadDefaultActions implémentée
- [ ] Fichier embarqué via go:embed
- [ ] Validation du contenu chargé
- [ ] Gestion d'erreurs robuste

### 5. Modifier ActionDeclaration pour supporter le marquage

**Objectif** : Ajouter un champ pour distinguer les actions par défaut.

- [ ] Ajouter le champ `IsDefault bool` à ActionDeclaration
- [ ] Mettre à jour les méthodes existantes si nécessaire
- [ ] Documenter le nouveau champ

**Fichier à modifier** :
- `tsd/parser/ast/action.go` (ou équivalent)

**Modification attendue** :
```go
// ActionDeclaration représente une déclaration d'action
type ActionDeclaration struct {
    Name       string
    Parameters []Parameter
    Body       []Statement
    Location   Location
    IsDefault  bool  // true si action par défaut du système
}
```

**Livrables** :
- [ ] Champ IsDefault ajouté
- [ ] Documentation mise à jour
- [ ] Tests mis à jour si nécessaire

### 6. Intégrer le chargement dans le compilateur

**Objectif** : Charger les actions par défaut lors de l'initialisation du compilateur.

- [ ] Modifier le constructeur du compilateur
- [ ] Charger les actions par défaut automatiquement
- [ ] Valider les doublons entre actions par défaut et utilisateur
- [ ] Gérer les erreurs de chargement
- [ ] Assurer l'ordre : défaut → utilisateur

**Fichier à modifier** :
- `tsd/compiler/compiler.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package compiler

import (
    "fmt"
    
    "tsd/internal/defaultactions"
    "tsd/parser/ast"
)

// NewCompiler crée un nouveau compilateur
func NewCompiler() (*Compiler, error) {
    c := &Compiler{
        actions:     make(map[string]*ast.ActionDeclaration),
        xupleSpaces: make(map[string]*ast.XupleSpaceDeclaration),
        // ... autres champs ...
    }
    
    // Charger les actions par défaut
    defaultActs, err := defaultactions.LoadDefaultActions()
    if err != nil {
        return nil, fmt.Errorf("failed to load default actions: %w", err)
    }
    
    // Enregistrer les actions par défaut
    for _, action := range defaultActs {
        c.actions[action.Name] = action
    }
    
    return c, nil
}

// RegisterAction enregistre une action (utilisateur)
func (c *Compiler) RegisterAction(action *ast.ActionDeclaration) error {
    // Vérifier si l'action existe déjà
    if existing, exists := c.actions[action.Name]; exists {
        // Message d'erreur différent si c'est une action par défaut
        if existing.IsDefault {
            return fmt.Errorf("cannot redefine default action '%s' at line %d (default actions: %v)",
                action.Name, action.Location.Line, defaultactions.DefaultActionNames)
        }
        return fmt.Errorf("action '%s' already declared at line %d",
            action.Name, existing.Location.Line)
    }
    
    c.actions[action.Name] = action
    return nil
}
```

**Livrables** :
- [ ] Chargement automatique implémenté
- [ ] Validation des doublons améliorée
- [ ] Messages d'erreur clairs et spécifiques
- [ ] Gestion d'erreurs robuste

### 7. Implémenter les exécuteurs d'actions natives

**Objectif** : Créer les implémentations réelles des actions par défaut.

- [ ] Créer le package `tsd/rete/actions`
- [ ] Implémenter chaque action par défaut
- [ ] Respecter l'interface ActionExecutor
- [ ] Gérer les erreurs spécifiques à chaque action
- [ ] Documenter chaque implémentation

**Fichier à créer** :
- `tsd/rete/actions/builtin.go`

**Code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
    "fmt"
    "log"
    
    "tsd/rete"
)

// BuiltinActionExecutor exécute les actions par défaut du système
type BuiltinActionExecutor struct {
    network      *rete.Network
    xupleManager XupleManager // Interface vers le module xuples
}

// NewBuiltinActionExecutor crée un nouvel exécuteur d'actions natives
func NewBuiltinActionExecutor(network *rete.Network, xupleManager XupleManager) *BuiltinActionExecutor {
    return &BuiltinActionExecutor{
        network:      network,
        xupleManager: xupleManager,
    }
}

// Execute exécute une action par défaut
func (e *BuiltinActionExecutor) Execute(actionName string, args []interface{}, token *rete.Token) error {
    switch actionName {
    case "Print":
        return e.executePrint(args)
    case "Log":
        return e.executeLog(args)
    case "Update":
        return e.executeUpdate(args, token)
    case "Insert":
        return e.executeInsert(args)
    case "Retract":
        return e.executeRetract(args)
    case "Xuple":
        return e.executeXuple(args, token)
    default:
        return fmt.Errorf("unknown builtin action: %s", actionName)
    }
}

// executePrint implémente l'action Print
func (e *BuiltinActionExecutor) executePrint(args []interface{}) error {
    if len(args) != 1 {
        return fmt.Errorf("Print expects 1 argument, got %d", len(args))
    }
    
    message, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Print expects string argument, got %T", args[0])
    }
    
    fmt.Println(message)
    return nil
}

// executeLog implémente l'action Log
func (e *BuiltinActionExecutor) executeLog(args []interface{}) error {
    if len(args) != 1 {
        return fmt.Errorf("Log expects 1 argument, got %d", len(args))
    }
    
    message, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Log expects string argument, got %T", args[0])
    }
    
    log.Printf("[TSD] %s", message)
    return nil
}

// executeUpdate implémente l'action Update
func (e *BuiltinActionExecutor) executeUpdate(args []interface{}, token *rete.Token) error {
    if len(args) != 1 {
        return fmt.Errorf("Update expects 1 argument, got %d", len(args))
    }
    
    // L'argument doit être un fait
    fact, ok := args[0].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Update expects fact argument, got %T", args[0])
    }
    
    // Déléguer au réseau RETE
    return e.network.UpdateFact(fact)
}

// executeInsert implémente l'action Insert
func (e *BuiltinActionExecutor) executeInsert(args []interface{}) error {
    if len(args) != 1 {
        return fmt.Errorf("Insert expects 1 argument, got %d", len(args))
    }
    
    fact, ok := args[0].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Insert expects fact argument, got %T", args[0])
    }
    
    // Déléguer au réseau RETE
    return e.network.InsertFact(fact)
}

// executeRetract implémente l'action Retract
func (e *BuiltinActionExecutor) executeRetract(args []interface{}) error {
    if len(args) != 1 {
        return fmt.Errorf("Retract expects 1 argument, got %d", len(args))
    }
    
    id, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Retract expects string argument, got %T", args[0])
    }
    
    // Déléguer au réseau RETE
    return e.network.RetractFact(id)
}

// executeXuple implémente l'action Xuple
func (e *BuiltinActionExecutor) executeXuple(args []interface{}, token *rete.Token) error {
    if len(args) != 2 {
        return fmt.Errorf("Xuple expects 2 arguments, got %d", len(args))
    }
    
    xuplespace, ok := args[0].(string)
    if !ok {
        return fmt.Errorf("Xuple expects string as first argument, got %T", args[0])
    }
    
    fact, ok := args[1].(*rete.Fact)
    if !ok {
        return fmt.Errorf("Xuple expects fact as second argument, got %T", args[1])
    }
    
    // Extraire les faits déclencheurs du token
    triggeringFacts := e.extractTriggeringFacts(token)
    
    // Déléguer au XupleManager
    return e.xupleManager.CreateXuple(xuplespace, fact, triggeringFacts)
}

// extractTriggeringFacts extrait tous les faits d'un token combiné
func (e *BuiltinActionExecutor) extractTriggeringFacts(token *rete.Token) []*rete.Fact {
    var facts []*rete.Fact
    
    // Parcourir la chaîne de tokens
    for t := token; t != nil; t = t.Parent {
        if t.Fact != nil {
            facts = append(facts, t.Fact)
        }
    }
    
    // Inverser pour avoir l'ordre chronologique
    for i := 0; i < len(facts)/2; i++ {
        facts[i], facts[len(facts)-1-i] = facts[len(facts)-1-i], facts[i]
    }
    
    return facts
}

// XupleManager interface vers le module xuples (définie ailleurs)
type XupleManager interface {
    CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error
}
```

**Livrables** :
- [ ] Package actions créé avec copyright
- [ ] Toutes les actions implémentées
- [ ] Validation des arguments
- [ ] Gestion d'erreurs robuste
- [ ] Extraction des faits déclencheurs
- [ ] Documentation GoDoc complète

### 8. Créer les tests des actions par défaut

**Objectif** : Tester le chargement et l'exécution des actions par défaut.

**Fichiers à créer** :
- `tsd/internal/defaultactions/loader_test.go`
- `tsd/rete/actions/builtin_test.go`

**Tests attendus** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package defaultactions

import "testing"

func TestLoadDefaultActions(t *testing.T) {
    t.Log("🧪 TEST CHARGEMENT ACTIONS PAR DÉFAUT")
    
    actions, err := LoadDefaultActions()
    if err != nil {
        t.Fatalf("❌ Erreur chargement: %v", err)
    }
    
    // Vérifier le nombre d'actions
    expectedCount := len(DefaultActionNames)
    if len(actions) != expectedCount {
        t.Errorf("❌ Attendu %d actions, reçu %d", expectedCount, len(actions))
    }
    
    // Vérifier que chaque action est marquée comme par défaut
    for _, action := range actions {
        if !action.IsDefault {
            t.Errorf("❌ Action '%s' devrait être marquée IsDefault", action.Name)
        }
    }
    
    // Vérifier que toutes les actions attendues sont présentes
    actionMap := make(map[string]bool)
    for _, action := range actions {
        actionMap[action.Name] = true
    }
    
    for _, name := range DefaultActionNames {
        if !actionMap[name] {
            t.Errorf("❌ Action par défaut manquante: %s", name)
        }
    }
    
    t.Log("✅ Toutes les actions par défaut chargées correctement")
}

func TestIsDefaultAction(t *testing.T) {
    t.Log("🧪 TEST IsDefaultAction")
    
    tests := []struct {
        name     string
        expected bool
    }{
        {"Print", true},
        {"Log", true},
        {"Update", true},
        {"Insert", true},
        {"Retract", true},
        {"Xuple", true},
        {"CustomAction", false},
        {"Unknown", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsDefaultAction(tt.name)
            if result != tt.expected {
                t.Errorf("❌ IsDefaultAction('%s') = %v, attendu %v",
                    tt.name, result, tt.expected)
            }
        })
    }
    
    t.Log("✅ IsDefaultAction fonctionne correctement")
}
```

**Livrables** :
- [ ] Tests du loader complets
- [ ] Tests de chaque action implémentée
- [ ] Tests d'erreurs (mauvais arguments, etc.)
- [ ] Couverture > 80%
- [ ] Tous les tests passent

### 9. Tester l'intégration avec le compilateur

**Objectif** : Vérifier que les actions par défaut sont chargées et validées correctement.

**Fichier à créer** :
- `tsd/compiler/defaultactions_test.go`

**Tests attendus** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package compiler

import (
    "strings"
    "testing"
)

func TestCompiler_DefaultActionsLoaded(t *testing.T) {
    t.Log("🧪 TEST CHARGEMENT ACTIONS PAR DÉFAUT DANS COMPILATEUR")
    
    compiler, err := NewCompiler()
    if err != nil {
        t.Fatalf("❌ Erreur création compilateur: %v", err)
    }
    
    // Vérifier que les actions par défaut sont présentes
    defaultNames := []string{"Print", "Log", "Update", "Insert", "Retract", "Xuple"}
    
    for _, name := range defaultNames {
        if _, exists := compiler.actions[name]; !exists {
            t.Errorf("❌ Action par défaut '%s' non chargée", name)
        }
    }
    
    t.Log("✅ Toutes les actions par défaut chargées dans le compilateur")
}

func TestCompiler_CannotRedefineDefaultAction(t *testing.T) {
    t.Log("🧪 TEST INTERDICTION REDÉFINITION ACTIONS PAR DÉFAUT")
    
    input := `
        action Print(msg: string) {
            // Tentative de redéfinition
        }
    `
    
    _, err := CompileTSD(input)
    if err == nil {
        t.Fatal("❌ Erreur attendue lors de la redéfinition d'action par défaut")
    }
    
    if !strings.Contains(err.Error(), "cannot redefine default action") {
        t.Errorf("❌ Message d'erreur incorrect: %v", err)
    }
    
    t.Log("✅ Redéfinition d'action par défaut correctement interdite")
}

func TestCompiler_DefaultActionsUsableInRules(t *testing.T) {
    t.Log("🧪 TEST UTILISATION ACTIONS PAR DÉFAUT DANS RÈGLES")
    
    input := `
        fact Person(name: string, age: int)
        
        rule "print-adult" {
            when {
                p: Person(age >= 18)
            }
            then {
                Print("Adult: " + p.name)
                Log("Found adult person")
            }
        }
    `
    
    _, err := CompileTSD(input)
    if err != nil {
        t.Fatalf("❌ Erreur compilation: %v", err)
    }
    
    t.Log("✅ Actions par défaut utilisables dans les règles")
}
```

**Livrables** :
- [ ] Tests d'intégration complets
- [ ] Test du chargement automatique
- [ ] Test de l'interdiction de redéfinition
- [ ] Test de l'utilisation dans les règles
- [ ] Tous les tests passent

## 📁 Structure attendue

```
tsd/
├── docs/xuples/implementation/
│   ├── 03-current-action-system.md
│   └── 04-default-actions-design.md
├── internal/defaultactions/
│   ├── defaults.tsd                 # Définitions embarquées
│   ├── loader.go                    # Chargement
│   └── loader_test.go               # Tests
├── rete/actions/
│   ├── builtin.go                   # Implémentations
│   └── builtin_test.go              # Tests
├── parser/ast/
│   └── action.go                    # Modifié (IsDefault)
└── compiler/
    ├── compiler.go                  # Modifié (chargement)
    └── defaultactions_test.go       # Tests intégration
```

## ✅ Critères de succès

- [ ] Fichier defaults.tsd créé avec copyright
- [ ] Chargement automatique implémenté
- [ ] Aucun hardcoding des actions
- [ ] Fichier embarqué via go:embed
- [ ] Toutes les actions implémentées
- [ ] Validation des doublons fonctionnelle
- [ ] Messages d'erreur clairs et spécifiques
- [ ] Tests complets avec couverture > 80%
- [ ] Tous les tests passent
- [ ] `make validate` passe sans erreur
- [ ] Documentation complète

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/design/` - Conception du module
- Effective Go - https://go.dev/doc/effective_go
- go:embed documentation

## 🎯 Prochaine étape

Une fois les actions par défaut implémentées, passer au prompt **05-modify-rete-immediate-execution.md** pour modifier le moteur RETE afin qu'il exécute les actions immédiatement au lieu de stocker les tokens.