# 🔧 Prompt 06 - Refactor et Cleanup

---

## ⚠️ CONTRAINTE ARCHITECTURALE STRICTE

**RÈGLE ABSOLUE**: Il est **STRICTEMENT INTERDIT** de créer des xuples directement dans les tests ou en batch (appel direct à `XupleManager.Create()`, `space.Add()`, etc.).

**Les xuples DOIVENT IMPÉRATIVEMENT être générés à partir de faits soumis au réseau RETE via des règles.**

✅ **CORRECT**:
```go
// Soumettre un fait au réseau RETE
network.Assert(ctx, fact)
// Le réseau évalue les règles, exécute l'action Xuple() qui crée le xuple
```

❌ **INTERDIT**:
```go
// NE JAMAIS faire ça dans les tests ou le code métier
xupleManager.Create(ctx, "space", fact)
space.Add(fact)
```

**Justification**:
- Garantit que tous les xuples passent par le réseau RETE
- Assure l'évaluation complète des règles et conditions
- Préserve la traçabilité et l'auditabilité
- Évite les contournements du pipeline qui introduisent des incohérences

Cette règle s'applique à **tous les tests, benchmarks, exemples et code de production**.

---

## 🎯 Objectif

**Nettoyer le code obsolète, refactorer les composants pour éliminer les patterns temporaires, et finaliser l'architecture** après la migration complète vers le package `api` et l'automatisation des xuples.

### Contexte

À ce stade du projet :
- ✅ Parser complet avec support des faits inline (Prompt 01)
- ✅ Package `api` fonctionnel (Prompt 02)
- ✅ Xuple-spaces créés automatiquement (Prompt 03)
- ✅ Actions Xuple automatiques (Prompt 04)
- ✅ Tests E2E migrés (Prompt 05)

**MAIS** : Le code contient encore :
- Ancien code de factory pluggable (obsolète)
- Méthodes de configuration manuelle (non utilisées)
- Code de workaround temporaire
- Imports inutilisés
- Documentation obsolète
- Duplication de logique

L'objectif est de **nettoyer complètement le code** pour n'avoir que l'implémentation finale, propre et maintenable.

### Prérequis

- ✅ Prompts 01-05 complétés
- ✅ Tous les tests migrés et passant
- ✅ Aucun test n'utilise l'ancien pattern

### Résultat Attendu Final

**Code avant (avec patterns obsolètes) :**

```go
// tsd/internal/rete/constraint_pipeline.go
type ConstraintPipeline struct {
    network              *Network
    storage              Storage
    xupleSpaceDefinitions map[string]map[string]interface{}
    xupleSpaceFactory    XupleSpaceFactory  // ❌ OBSOLETE
    xupleActionHandler   XupleActionHandler // ❌ OBSOLETE
}

// ❌ OBSOLETE - À supprimer
func (cp *ConstraintPipeline) SetXupleSpaceFactory(factory XupleSpaceFactory) {
    cp.xupleSpaceFactory = factory
}

// ❌ OBSOLETE - À supprimer
func (cp *ConstraintPipeline) createXupleSpaces() error {
    if cp.xupleSpaceFactory == nil {
        return nil
    }
    // ... code de factory ...
}
```

**Code après (nettoyé) :**

```go
// tsd/internal/rete/constraint_pipeline.go
type ConstraintPipeline struct {
    network              *Network
    storage              Storage
    xupleSpaceDefinitions map[string]map[string]interface{}
    // Factories et handlers supprimés - gérés par le package api
}

// Méthodes de configuration manuelle supprimées
// La création des xuple-spaces est gérée par api.Pipeline
```

---

## 📋 Analyse Préliminaire

### 1. Identifier le Code Obsolète

**Fichiers à examiner pour code obsolète :**

```
tsd/internal/rete/
├── constraint_pipeline.go    # Factory patterns, méthodes de config
├── xuple_action.go           # Si existe (à déplacer ou supprimer)
├── network.go                # Méthodes de config manuelle
└── action.go                 # Code temporaire

tsd/internal/constraint/
├── rete_converter.go         # Code de workaround
└── parser.go                 # Code temporaire de parsing

tsd/docs/
├── XUPLES_E2E_AUTOMATIC.md   # Sections obsolètes
└── ARCHITECTURE.md           # Diagrammes à mettre à jour
```

**Questions à résoudre :**

1. **Quelles méthodes publiques ne sont plus utilisées ?**
   - `SetXupleSpaceFactory()`
   - `SetXupleActionHandler()`
   - Méthodes de configuration manuelle

2. **Quels types/interfaces sont obsolètes ?**
   - `XupleSpaceFactory`
   - `XupleActionHandler`
   - Types intermédiaires temporaires

3. **Quel code de workaround peut être supprimé ?**
   - Création manuelle de xuples dans les tests
   - Configuration conditionnelle (if factory != nil)
   - Code de fallback

### 2. Analyser les Dépendances

**Créer un graphe de dépendances pour identifier les cycles résiduels :**

```
Avant cleanup:
    constraint → rete (OK)
    rete → xuples (via factory) ❌ CYCLE potentiel
    api → rete (OK)
    api → xuples (OK)
    api → constraint (OK)

Après cleanup:
    constraint → rete (OK)
    rete ↛ xuples (pas de dépendance directe) ✅
    api → rete (OK)
    api → xuples (OK)
    api → constraint (OK)
```

### 3. Planifier le Refactoring

**Stratégie :**

1. **Phase 1** : Identifier et documenter tout le code à supprimer
2. **Phase 2** : Supprimer les méthodes publiques obsolètes
3. **Phase 3** : Supprimer les types/interfaces obsolètes
4. **Phase 4** : Nettoyer les imports inutilisés
5. **Phase 5** : Simplifier la logique (éliminer les if/else temporaires)
6. **Phase 6** : Refactorer pour améliorer la cohérence
7. **Phase 7** : Mettre à jour la documentation

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Supprimer le Pattern Factory des Xuple-Spaces

**Fichier :** `tsd/internal/rete/constraint_pipeline.go`

**Objectif :** Supprimer tout le code lié à la factory pluggable.

#### 1.1 Identifier le Code à Supprimer

```go
// ❌ SUPPRIMER - Types obsolètes
type XupleSpaceFactory func(name string, props map[string]interface{}) error

type XupleActionHandler interface {
    HandleXupleAction(spaceName string, fact *Fact) error
}
```

#### 1.2 Nettoyer ConstraintPipeline

```go
// AVANT
type ConstraintPipeline struct {
    network              *Network
    storage              Storage
    xupleSpaceDefinitions map[string]map[string]interface{}
    xupleSpaceFactory    XupleSpaceFactory  // ❌ À supprimer
    xupleActionHandler   XupleActionHandler // ❌ À supprimer
    mu                   sync.RWMutex
}

// APRÈS (nettoyé)
type ConstraintPipeline struct {
    network              *Network
    storage              Storage
    xupleSpaceDefinitions map[string]map[string]interface{}
    mu                   sync.RWMutex
}
```

#### 1.3 Supprimer les Méthodes de Configuration

```go
// ❌ SUPPRIMER - Méthodes obsolètes

// SetXupleSpaceFactory configure la factory pour créer des xuple-spaces.
// OBSOLETE: Les xuple-spaces sont maintenant créés automatiquement par api.Pipeline.
func (cp *ConstraintPipeline) SetXupleSpaceFactory(factory XupleSpaceFactory) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    cp.xupleSpaceFactory = factory
}

// SetXupleActionHandler configure le handler pour les actions Xuple.
// OBSOLETE: Les actions Xuple sont maintenant gérées automatiquement.
func (cp *ConstraintPipeline) SetXupleActionHandler(handler XupleActionHandler) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    cp.xupleActionHandler = handler
}

// createXupleSpaces appelle la factory pour créer les xuple-spaces.
// OBSOLETE: Cette méthode n'est plus appelée.
func (cp *ConstraintPipeline) createXupleSpaces() error {
    if cp.xupleSpaceFactory == nil {
        return nil
    }
    
    for name, props := range cp.xupleSpaceDefinitions {
        if err := cp.xupleSpaceFactory(name, props); err != nil {
            return fmt.Errorf("creating xuple-space '%s': %w", name, err)
        }
    }
    
    return nil
}
```

#### 1.4 Version Finale de ConstraintPipeline

```go
// ConstraintPipeline gère le parsing et la construction du réseau RETE
// à partir de fichiers TSD.
//
// Note: La création des xuple-spaces et l'enregistrement des actions Xuple
// sont gérés par le package api, pas par ce pipeline.
type ConstraintPipeline struct {
    network              *Network
    storage              Storage
    xupleSpaceDefinitions map[string]map[string]interface{}
    mu                   sync.RWMutex
}

// NewConstraintPipeline crée un nouveau pipeline de contraintes.
func NewConstraintPipeline(network *Network, storage Storage) *ConstraintPipeline {
    return &ConstraintPipeline{
        network:              network,
        storage:              storage,
        xupleSpaceDefinitions: make(map[string]map[string]interface{}),
    }
}

// IngestFile parse un fichier TSD et construit le réseau RETE.
func (cp *ConstraintPipeline) IngestFile(filepath string) error {
    // Parse le fichier
    parser := constraint.NewParserFromFile(filepath)
    program, err := parser.Parse()
    if err != nil {
        return fmt.Errorf("parsing file: %w", err)
    }
    
    // Convertit l'AST en réseau RETE
    converter := constraint.NewASTConverter(cp)
    if err := converter.Convert(program); err != nil {
        return fmt.Errorf("converting AST: %w", err)
    }
    
    return nil
}

// GetXupleSpaceDefinitions retourne les définitions de xuple-spaces
// extraites du fichier TSD.
//
// Ces définitions sont utilisées par api.Pipeline pour créer
// automatiquement les xuple-spaces.
func (cp *ConstraintPipeline) GetXupleSpaceDefinitions() map[string]map[string]interface{} {
    cp.mu.RLock()
    defer cp.mu.RUnlock()
    
    // Retourner une copie pour éviter les mutations
    result := make(map[string]map[string]interface{})
    for name, props := range cp.xupleSpaceDefinitions {
        propsCopy := make(map[string]interface{})
        for k, v := range props {
            propsCopy[k] = v
        }
        result[name] = propsCopy
    }
    
    return result
}

// Network retourne le réseau RETE.
func (cp *ConstraintPipeline) Network() *Network {
    return cp.network
}

// Storage retourne le storage utilisé.
func (cp *ConstraintPipeline) Storage() Storage {
    return cp.storage
}
```

---

### Tâche 2: Nettoyer le Package RETE

**Fichier :** `tsd/internal/rete/network.go`

**Objectif :** Supprimer les méthodes de configuration manuelle.

#### 2.1 Supprimer les Méthodes Obsolètes

```go
// ❌ SUPPRIMER - Méthodes obsolètes de configuration manuelle

// RegisterXupleSpace enregistre un xuple-space manuellement.
// OBSOLETE: Les xuple-spaces sont maintenant créés automatiquement.
func (n *Network) RegisterXupleSpace(name string, space XupleSpace) {
    // ...
}

// SetXupleManager configure le XupleManager.
// OBSOLETE: Le XupleManager est géré par api.Pipeline.
func (n *Network) SetXupleManager(manager *xuples.XupleManager) {
    // ...
}
```

#### 2.2 Version Finale de Network (extrait)

```go
// Network représente le réseau RETE.
type Network struct {
    types           map[string]*TypeDefinition
    rules           map[string]*Rule
    facts           map[string]*Fact
    actionFactories map[string]ActionFactory
    executor        *ActionExecutor
    mu              sync.RWMutex
}

// NewNetwork crée un nouveau réseau RETE.
func NewNetwork() *Network {
    return &Network{
        types:           make(map[string]*TypeDefinition),
        rules:           make(map[string]*Rule),
        facts:           make(map[string]*Fact),
        actionFactories: make(map[string]ActionFactory),
    }
}

// RegisterAction enregistre une factory d'action.
func (n *Network) RegisterAction(name string, factory ActionFactory) {
    n.mu.Lock()
    defer n.mu.Unlock()
    n.actionFactories[name] = factory
}

// CreateAction crée une action à partir de son nom et de ses arguments.
func (n *Network) CreateAction(name string, args []ActionArgument) (Action, error) {
    n.mu.RLock()
    factory, ok := n.actionFactories[name]
    n.mu.RUnlock()
    
    if !ok {
        return nil, fmt.Errorf("unknown action: %s", name)
    }
    
    return factory(args)
}

// SetActionExecutor définit l'exécuteur d'actions pour ce réseau.
func (n *Network) SetActionExecutor(executor *ActionExecutor) {
    n.mu.Lock()
    defer n.mu.Unlock()
    n.executor = executor
}

// GetActionExecutor retourne l'exécuteur d'actions.
func (n *Network) GetActionExecutor() *ActionExecutor {
    n.mu.RLock()
    defer n.mu.RUnlock()
    return n.executor
}

// ... autres méthodes (RegisterType, CreateFact, Assert, etc.)
```

---

### Tâche 3: Supprimer le Code de Workaround

**Fichiers multiples**

**Objectif :** Éliminer tous les workarounds temporaires.

#### 3.1 Identifier les Workarounds

**Pattern typique de workaround :**

```go
// ❌ Workaround à supprimer
if cp.xupleSpaceFactory != nil {
    // Code temporaire
    cp.createXupleSpaces()
} else {
    // Fallback
    log.Warn("No factory configured, xuple-spaces not created")
}
```

#### 3.2 Nettoyer les Conditions Temporaires

**AVANT :**

```go
// Dans api/pipeline.go (hypothétique ancien code)
func (p *Pipeline) IngestFile(filepath string) (*Result, error) {
    // ...
    
    // Workaround: création conditionnelle
    if p.config.EnableXuples { // ❌ Condition inutile
        if err := p.createXupleSpaces(); err != nil {
            return nil, err
        }
    }
    
    // ...
}
```

**APRÈS :**

```go
// Dans api/pipeline.go (version finale)
func (p *Pipeline) IngestFile(filepath string) (*Result, error) {
    startTime := time.Now()
    
    // Parse et build le réseau RETE
    parseStart := time.Now()
    if err := p.retePipeline.IngestFile(filepath); err != nil {
        return nil, p.wrapError("parse", err)
    }
    parseDuration := time.Since(parseStart)
    
    // Créer automatiquement les xuple-spaces (toujours, pas de condition)
    createSpacesStart := time.Now()
    if err := p.createXupleSpaces(); err != nil {
        return nil, p.wrapError("create-xuple-spaces", err)
    }
    createSpacesDuration := time.Since(createSpacesStart)
    
    // Construire le résultat
    // ...
}
```

---

### Tâche 4: Simplifier les Imports

**Tous les fichiers**

**Objectif :** Supprimer les imports inutilisés.

#### 4.1 Script de Détection

**Créer :** `tsd/scripts/maj-xuples/cleanup_imports.sh`

```bash
#!/bin/bash

# Script pour détecter et nettoyer les imports inutilisés

echo "🔍 Détection des imports inutilisés..."

# Utiliser goimports pour formater et nettoyer
find tsd -name "*.go" -type f | while read -r file; do
    echo "Nettoyage: $file"
    goimports -w "$file"
done

echo "✅ Imports nettoyés"

# Vérifier les imports cycliques
echo ""
echo "🔍 Vérification des imports cycliques..."
go list -f '{{ join .DepsErrors "\n" }}' ./... 2>&1 | grep -i "import cycle" || echo "✅ Aucun import cyclique détecté"

# Rapport des packages inutilisés
echo ""
echo "📊 Packages potentiellement inutilisés:"
go list -f '{{if gt (len .Imports) 0}}{{.ImportPath}}: {{join .Imports ", "}}{{end}}' ./... | \
    grep -E "(xuples/.*factory|rete/.*handler)" || echo "✅ Aucun package obsolète importé"
```

#### 4.2 Exécution

```bash
chmod +x tsd/scripts/maj-xuples/cleanup_imports.sh
./tsd/scripts/maj-xuples/cleanup_imports.sh
```

---

### Tâche 5: Refactorer pour la Cohérence

**Fichiers multiples**

**Objectif :** Améliorer la cohérence du code (nommage, structure, patterns).

#### 5.1 Standardiser le Nommage

**Règles de nommage à appliquer :**

1. **Préfixe "Xuple"** : Majuscule (XupleSpace, XupleManager, XupleAction)
2. **Méthodes publiques** : CapitalCase
3. **Méthodes privées** : camelCase
4. **Variables locales** : camelCase
5. **Constantes** : CamelCase ou UPPER_SNAKE_CASE selon contexte

**Exemple de refactoring :**

```go
// AVANT (incohérent)
type xuple_Space struct { // ❌ Mauvais nommage
    name string
}

func (xs *xuple_Space) Get_All() []*Xuple { // ❌ Underscore inutile
    // ...
}

// APRÈS (cohérent)
type XupleSpace struct {
    name string
}

func (xs *XupleSpace) GetAll() []*Xuple {
    // ...
}
```

#### 5.2 Standardiser la Gestion d'Erreurs

**Pattern uniforme pour les erreurs :**

```go
// AVANT (incohérent)
func someFunc() error {
    if err := doSomething(); err != nil {
        return errors.New("error doing something: " + err.Error()) // ❌ Concaténation
    }
    
    if err := doOther(); err != nil {
        log.Printf("Error: %v", err) // ❌ Log + return
        return err
    }
    
    return nil
}

// APRÈS (cohérent)
func someFunc() error {
    if err := doSomething(); err != nil {
        return fmt.Errorf("doing something: %w", err) // ✅ Wrapping
    }
    
    if err := doOther(); err != nil {
        return fmt.Errorf("doing other: %w", err) // ✅ Cohérent
    }
    
    return nil
}
```

#### 5.3 Standardiser les Commentaires GoDoc

**Template pour les commentaires :**

```go
// AVANT (incomplet)
// CreateXupleSpace crée un space
func (m *XupleManager) CreateXupleSpace(name string) (*XupleSpace, error) {
    // ...
}

// APRÈS (complet)
// CreateXupleSpace crée un nouveau xuple-space avec les politiques spécifiées.
//
// Paramètres:
//   - name: Nom unique du xuple-space
//   - selection: Politique de sélection (FIFO, LIFO, Random)
//   - consumption: Politique de consommation (Once, PerAgent)
//   - retention: Politique de rétention (Unlimited, Duration)
//
// Retourne:
//   - Le xuple-space créé
//   - Une erreur si le nom existe déjà ou si les paramètres sont invalides
//
// Exemple:
//   space, err := manager.CreateXupleSpace(
//       "alerts",
//       xuples.SelectionFIFO,
//       xuples.ConsumptionOnce,
//       xuples.RetentionUnlimited,
//   )
func (m *XupleManager) CreateXupleSpace(
    name string,
    selection SelectionPolicy,
    consumption ConsumptionPolicy,
    retention RetentionPolicy,
) (*XupleSpace, error) {
    // ...
}
```

---

### Tâche 6: Optimiser les Performances

**Fichiers critiques : `network.go`, `xuplespace.go`, `manager.go`**

**Objectif :** Identifier et corriger les goulots d'étranglement.

#### 6.1 Profiling

**Créer :** `tsd/test/benchmark/xuples_bench_test.go`

```go
package benchmark

import (
    "testing"
    
    "github.com/resinsec/tsd/api"
)

// BenchmarkXupleCreation mesure la performance de création de xuples.
func BenchmarkXupleCreation(b *testing.B) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id": i,
        })
        result.Network().Assert(fact)
    }
}

// BenchmarkXupleRetrieval mesure la performance de récupération.
func BenchmarkXupleRetrieval(b *testing.B) {
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    // Pré-remplir avec 1000 xuples
    for i := 0; i < 1000; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id": i,
        })
        result.Network().Assert(fact)
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        result.Retrieve("test", nil, "agent")
    }
}

// BenchmarkLargeXupleSpace mesure la performance avec un grand espace.
func BenchmarkLargeXupleSpace(b *testing.B) {
    tsdContent := `
xuple-space test { selection: fifo, max-size: 100000 }
type T { id: int, data: string }
type X { id: int, data: string }
rule R { when { t: T() } then { Xuple("test", X(id: t.id, data: t.data)) } }
`
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id":   i,
            "data": "some data here",
        })
        result.Network().Assert(fact)
    }
}
```

#### 6.2 Optimisations Potentielles

**Si le profiling révèle des problèmes :**

1. **Pool d'objets** pour réduire les allocations
2. **Synchronisation optimisée** (RWMutex vs Mutex)
3. **Indexation** des xuples pour recherche rapide
4. **Batching** des opérations réseau

**Exemple d'optimisation :**

```go
// AVANT (lock global)
func (m *XupleManager) GetSpace(name string) *XupleSpace {
    m.mu.Lock()         // ❌ Lock exclusif
    defer m.mu.Unlock()
    return m.spaces[name]
}

// APRÈS (lock en lecture)
func (m *XupleManager) GetSpace(name string) *XupleSpace {
    m.mu.RLock()        // ✅ Lock partagé (lecture)
    defer m.mu.RUnlock()
    return m.spaces[name]
}
```

---

### Tâche 7: Nettoyer la Documentation

**Fichiers : `docs/*.md`**

**Objectif :** Supprimer les sections obsolètes, mettre à jour les exemples.

#### 7.1 Fichiers à Mettre à Jour

**Liste :**

1. `docs/XUPLES_E2E_AUTOMATIC.md` - Supprimer sections factory
2. `docs/ARCHITECTURE.md` - Mettre à jour diagrammes
3. `docs/API_USAGE.md` - Supprimer exemples de config manuelle
4. `docs/TSD_LANGUAGE.md` - Vérifier cohérence
5. `README.md` - Mettre à jour exemples

#### 7.2 XUPLES_E2E_AUTOMATIC.md

**Sections à supprimer :**

```markdown
<!-- ❌ SUPPRIMER cette section -->
## Configuration Manuelle (Legacy)

Si vous utilisez l'ancien pattern, vous devez configurer manuellement:

\`\`\`go
pipeline.SetXupleSpaceFactory(func(name string, props map[string]interface{}) error {
    // ...
})
\`\`\`

**Note:** Ce pattern est obsolète. Utilisez le package `api` à la place.
```

**Nouvelle version (simplifiée) :**

```markdown
## Utilisation

L'intégration des xuples est **entièrement automatique** :

\`\`\`go
import "github.com/resinsec/tsd/api"

pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")

// Les xuple-spaces sont créés automatiquement
// Les actions Xuple fonctionnent immédiatement
\`\`\`

C'est tout ! Aucune configuration nécessaire.
```

#### 7.3 ARCHITECTURE.md

**Mettre à jour le diagramme :**

```markdown
## Architecture Finale

\`\`\`
┌─────────────────────────────────────────────────┐
│                 Package API                      │
│  (Point d'entrée unifié - api.Pipeline)         │
│                                                  │
│  - NewPipeline()                                │
│  - IngestFile() / IngestString()                │
│  - Configuration automatique                    │
└─────────────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
        ▼           ▼           ▼
┌──────────┐  ┌──────────┐  ┌──────────┐
│   RETE   │  │ Xuples   │  │Constraint│
│ Network  │  │ Manager  │  │ Parser   │
│          │  │          │  │          │
│ - Facts  │  │ - Spaces │  │ - AST    │
│ - Rules  │  │ - Xuples │  │ - Convert│
│ - Actions│  │ - Policies│ │          │
└──────────┘  └──────────┘  └──────────┘

Note: Aucun cycle d'import
      Tout est orchestré par le package api
\`\`\`
```

---

## 🧪 Tests de Validation

### Test 1: Vérification du Cleanup

**Fichier :** `tsd/test/cleanup_validation_test.go`

```go
package test

import (
    "go/ast"
    "go/parser"
    "go/token"
    "testing"
    
    "github.com/stretchr/testify/assert"
)

// TestNoObsoleteFunctions vérifie qu'aucune fonction obsolète n'existe.
func TestNoObsoleteFunctions(t *testing.T) {
    obsoleteFunctions := []string{
        "SetXupleSpaceFactory",
        "SetXupleActionHandler",
        "createXupleSpaces", // Fonction privée mais obsolète
    }
    
    // Parser tous les fichiers Go du projet
    fset := token.NewFileSet()
    packages, err := parser.ParseDir(fset, "tsd/internal/rete", nil, 0)
    assert.NoError(t, err)
    
    for _, pkg := range packages {
        for _, file := range pkg.Files {
            ast.Inspect(file, func(n ast.Node) bool {
                if fn, ok := n.(*ast.FuncDecl); ok {
                    funcName := fn.Name.Name
                    for _, obsolete := range obsoleteFunctions {
                        assert.NotEqual(t, obsolete, funcName,
                            "Obsolete function '%s' still exists", obsolete)
                    }
                }
                return true
            })
        }
    }
}

// TestNoObsoleteTypes vérifie qu'aucun type obsolète n'existe.
func TestNoObsoleteTypes(t *testing.T) {
    obsoleteTypes := []string{
        "XupleSpaceFactory",
        "XupleActionHandler",
    }
    
    fset := token.NewFileSet()
    packages, err := parser.ParseDir(fset, "tsd/internal/rete", nil, 0)
    assert.NoError(t, err)
    
    for _, pkg := range packages {
        for _, file := range pkg.Files {
            ast.Inspect(file, func(n ast.Node) bool {
                if ts, ok := n.(*ast.TypeSpec); ok {
                    typeName := ts.Name.Name
                    for _, obsolete := range obsoleteTypes {
                        assert.NotEqual(t, obsolete, typeName,
                            "Obsolete type '%s' still exists", obsolete)
                    }
                }
                return true
            })
        }
    }
}
```

### Test 2: Vérification des Imports

```go
// TestNoCircularImports vérifie qu'il n'y a pas d'imports cycliques.
func TestNoCircularImports(t *testing.T) {
    // Exécuter go list pour détecter les cycles
    cmd := exec.Command("go", "list", "-f", "{{ join .DepsErrors \"\\n\" }}", "./...")
    output, err := cmd.CombinedOutput()
    
    // S'il y a une erreur, vérifier que ce n'est pas un cycle
    if err != nil {
        assert.NotContains(t, string(output), "import cycle",
            "Circular import detected: %s", string(output))
    }
}

// TestNoUnusedImports vérifie qu'il n'y a pas d'imports inutilisés.
func TestNoUnusedImports(t *testing.T) {
    cmd := exec.Command("goimports", "-l", ".")
    output, err := cmd.Output()
    assert.NoError(t, err)
    
    // goimports -l ne devrait rien retourner (tous les fichiers sont formatés)
    assert.Empty(t, string(output),
        "Files with formatting issues (unused imports): %s", string(output))
}
```

### Test 3: Validation des Performances

```go
// TestPerformanceNoRegression vérifie qu'il n'y a pas de régression de performance.
func TestPerformanceNoRegression(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance test in short mode")
    }
    
    tsdContent := `
xuple-space test { selection: fifo }
type T { id: int }
type X { id: int }
rule R { when { t: T() } then { Xuple("test", X(id: t.id)) } }
`
    
    pipeline, _ := api.NewPipeline()
    result, _ := pipeline.IngestString(tsdContent)
    
    // Mesurer le temps pour créer 10000 xuples
    start := time.Now()
    for i := 0; i < 10000; i++ {
        fact := result.Network().CreateFact("T", map[string]interface{}{
            "id": i,
        })
        result.Network().Assert(fact)
    }
    duration := time.Since(start)
    
    // Ne devrait pas prendre plus de 1 seconde pour 10k xuples
    assert.Less(t, duration, 1*time.Second,
        "Creating 10k xuples took %v (expected < 1s)", duration)
}
```

---

## ✅ Checklist de Validation

### Code Obsolète Supprimé

- [ ] `XupleSpaceFactory` type supprimé
- [ ] `XupleActionHandler` interface supprimée
- [ ] `SetXupleSpaceFactory()` méthode supprimée
- [ ] `SetXupleActionHandler()` méthode supprimée
- [ ] `createXupleSpaces()` (ancienne version) supprimée
- [ ] Tous les workarounds temporaires supprimés

### Code Nettoyé

- [ ] Imports inutilisés supprimés
- [ ] Conditions temporaires (if factory != nil) supprimées
- [ ] Code mort (dead code) supprimé
- [ ] Variables non utilisées supprimées

### Cohérence

- [ ] Nommage uniforme (XupleSpace, XupleManager, etc.)
- [ ] Gestion d'erreurs standardisée (fmt.Errorf avec %w)
- [ ] Commentaires GoDoc complets et cohérents
- [ ] Patterns de code uniformes

### Performance

- [ ] Benchmarks exécutés
- [ ] Pas de régression de performance (< 10%)
- [ ] Locks optimisés (RWMutex où approprié)
- [ ] Allocations minimisées

### Documentation

- [ ] Sections obsolètes supprimées
- [ ] Exemples mis à jour
- [ ] Diagrammes d'architecture mis à jour
- [ ] README à jour

### Tests

- [ ] Tests de validation du cleanup passent
- [ ] Tests d'imports cycliques passent
- [ ] Tests de performance passent
- [ ] Couverture de code maintenue (> 80%)

### Standards

- [ ] Code formaté (`gofmt`)
- [ ] Pas de warnings du linter
- [ ] `go vet` passe sans erreurs
- [ ] `golangci-lint` passe

---

## 📝 Documentation à Mettre à Jour

### 1. Changelog

**Créer :** `CHANGELOG.md`

```markdown
# Changelog

## [v2.0.0] - 2024-XX-XX

### 🎉 Features Majeures

- **Automatisation complète des xuples**
  - Xuple-spaces créés automatiquement lors du parsing
  - Actions `Xuple()` fonctionnent sans configuration
  - Package `api` comme point d'entrée unifié

### ✨ Améliorations

- Parser TSD étendu avec support des faits inline
- Références aux champs dans les actions (`t.sensorId`)
- Actions multiples séparées par virgules
- Politiques de xuple-spaces configurables dans TSD

### 🔧 Changements Techniques

- **BREAKING:** Suppression du pattern factory pour xuple-spaces
- **BREAKING:** `SetXupleSpaceFactory()` supprimée (obsolète)
- **BREAKING:** `SetXupleActionHandler()` supprimée (obsolète)
- Refactoring complet de l'architecture interne
- Élimination des imports cycliques

### 🗑️ Suppressions

- Pattern factory pluggable (remplacé par automatisation)
- Méthodes de configuration manuelle
- Code de workaround temporaire

### 📚 Documentation

- Guide d'utilisation simplifié
- Nouveaux exemples E2E
- Architecture mise à jour

### 🐛 Corrections

- Correction des imports cycliques entre `rete` et `xuples`
- Amélioration de la gestion d'erreurs
- Optimisation des performances (locks, allocations)

### 🔒 Migration

Pour migrer depuis v1.x :

**Avant :**
\`\`\`go
network := rete.NewNetwork()
pipeline := constraint.NewConstraintPipeline(network, storage)
xupleManager := xuples.NewXupleManager()
pipeline.SetXupleSpaceFactory(...)
\`\`\`

**Après :**
\`\`\`go
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")
\`\`\`

Voir [MIGRATION.md](docs/MIGRATION.md) pour les détails.
```

### 2. Guide de Migration

**Créer :** `docs/MIGRATION.md`

```markdown
# Guide de Migration v1.x → v2.0

## Vue d'Ensemble

La version 2.0 introduit l'automatisation complète des xuples, éliminant
toute configuration manuelle. Ce guide vous aide à migrer votre code.

## Changements Incompatibles

### 1. Point d'Entrée Unifié

**v1.x (ancien) :**
\`\`\`go
import (
    "github.com/resinsec/tsd/internal/rete"
    "github.com/resinsec/tsd/internal/constraint"
    "github.com/resinsec/tsd/xuples"
)

network := rete.NewNetwork()
storage := rete.NewMemoryStorage()
pipeline := constraint.NewConstraintPipeline(network, storage)
xupleManager := xuples.NewXupleManager()
\`\`\`

**v2.0 (nouveau) :**
\`\`\`go
import "github.com/resinsec/tsd/api"

pipeline, err := api.NewPipeline()
\`\`\`

### 2. Suppression de SetXupleSpaceFactory

**v1.x :**
\`\`\`go
pipeline.SetXupleSpaceFactory(func(name string, props map[string]interface{}) error {
    space, err := xupleManager.CreateXupleSpace(name, ...)
    return err
})
\`\`\`

**v2.0 :**
Cette méthode n'existe plus. Les xuple-spaces sont créés automatiquement
à partir des définitions dans le fichier TSD.

\`\`\`tsd
xuple-space alerts {
    selection: fifo,
    consumption: once
}
\`\`\`

### 3. Actions Xuple Automatiques

**v1.x :**
\`\`\`go
pipeline.SetXupleActionHandler(myHandler)
\`\`\`

**v2.0 :**
Les actions `Xuple()` fonctionnent automatiquement. Aucune configuration nécessaire.

## Migration Étape par Étape

### Étape 1: Mettre à Jour les Imports

Remplacer :
\`\`\`go
import (
    "github.com/resinsec/tsd/internal/rete"
    "github.com/resinsec/tsd/internal/constraint"
)
\`\`\`

Par :
\`\`\`go
import "github.com/resinsec/tsd/api"
\`\`\`

### Étape 2: Simplifier l'Initialisation

Remplacer tout le code d'initialisation par :
\`\`\`go
pipeline, err := api.NewPipeline()
if err != nil {
    // Gérer l'erreur
}
\`\`\`

### Étape 3: Utiliser Result au Lieu de Composants Séparés

**v1.x :**
\`\`\`go
pipeline.IngestFile("rules.tsd")
xuples := xupleManager.GetSpace("alerts").GetAll()
\`\`\`

**v2.0 :**
\`\`\`go
result, err := pipeline.IngestFile("rules.tsd")
xuples := result.GetXuples("alerts")
\`\`\`

### Étape 4: Supprimer les Configurations Manuelles

Supprimer tous les appels à :
- `SetXupleSpaceFactory()`
- `SetXupleActionHandler()`
- Création manuelle de xuple-spaces

## Compatibilité

### Packages Internes

Les packages suivants sont maintenant **internes** et ne doivent plus
être importés directement :
- `internal/rete`
- `internal/constraint`

Utilisez uniquement :
- `api` (point d'entrée)
- `xuples` (types publics)

### Tests

Migrer les tests en utilisant les helpers :
\`\`\`go
import "github.com/resinsec/tsd/test/testutil"

_, result := testutil.CreatePipelineFromTSD(t, tsdContent)
\`\`\`

## Aide

Pour toute question, voir :
- [Documentation API](API_USAGE.md)
- [Exemples](examples/)
- [Issues GitHub](https://github.com/resinsec/tsd/issues)
```

---

## 🎯 Résultat Attendu

### Métriques de Code

**Avant cleanup :**
- Lignes de code total : ~15,000
- Code obsolète : ~2,000 (13%)
- Imports inutilisés : ~50
- Duplication : ~5%

**Après cleanup :**
- Lignes de code total : ~13,000
- Code obsolète : 0 (0%)
- Imports inutilisés : 0
- Duplication : < 2%

### Qualité

- **Couverture de tests** : > 85%
- **Complexité cyclomatique** : < 15 (moyenne)
- **Maintenabilité** : Score A (selon golangci-lint)
- **Dette technique** : < 5% (selon SonarQube)

---

## 🔗 Dépendances

### Entrantes

- ✅ Prompts 01-05 complétés
- ✅ Tous les tests migrés

### Sortantes

- ➡️ Prompt 07 : Documentation finale et release

---

## 🚀 Stratégie d'Implémentation

1. **Phase 1: Identification** (2h)
   - Lister tout le code obsolète
   - Créer un plan de suppression
   - Documenter les changements

2. **Phase 2: Suppression** (3h)
   - Supprimer les types/interfaces obsolètes
   - Supprimer les méthodes de configuration
   - Nettoyer les workarounds

3. **Phase 3: Refactoring** (2h)
   - Standardiser le nommage
   - Améliorer la gestion d'erreurs
   - Optimiser les performances

4. **Phase 4: Tests** (2h)
   - Tests de validation du cleanup
   - Tests de performance
   - Tests de régression

5. **Phase 5: Documentation** (2h)
   - Changelog
   - Guide de migration
   - Mise à jour des docs

**Estimation totale : 11-13 heures**

---

## 📊 Critères de Succès

- [ ] Zéro code obsolète dans la codebase
- [ ] Zéro imports inutilisés
- [ ] Zéro imports cycliques
- [ ] Nommage 100% cohérent
- [ ] Gestion d'erreurs standardisée
- [ ] GoDoc complet sur tous les exports publics
- [ ] Tests de cleanup passent (100%)
- [ ] Couverture maintenue ou améliorée (> 85%)
- [ ] Performance maintenue ou améliorée
- [ ] Documentation complète et à jour
- [ ] Changelog détaillé
- [ ] Guide de migration clair

---

**FIN DU PROMPT 06**