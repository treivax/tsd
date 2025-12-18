# 🔧 Prompt 03 - Création Automatique des Xuple-Spaces

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

**Automatiser la création des xuple-spaces lors du parsing des fichiers TSD**, en détectant les définitions `xuple-space` dans le fichier source et en les instanciant automatiquement via le pipeline API, éliminant ainsi toute configuration manuelle.

### Contexte

Actuellement, même avec le package `api` et le parser amélioré (Prompts 01 et 02), la création des xuple-spaces nécessite une étape manuelle ou une configuration explicite via factory. L'objectif est que **le simple fait de déclarer un xuple-space dans un fichier TSD provoque automatiquement sa création lors de l'ingestion**.

### Prérequis

- ✅ Prompt 01 complété : parser supporte faits inline et références aux champs
- ✅ Prompt 02 complété : package `api` existe avec `Pipeline.IngestFile()`
- ✅ Le package `xuples` expose `NewXupleManager()` et `CreateXupleSpace()`

### Résultat Attendu Final

Après ingestion d'un fichier TSD contenant :

```tsd
xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

xuple-space notifications {
    selection: random,
    consumption: per-agent,
    max-size: 1000
}
```

Les xuple-spaces `alerts` et `notifications` sont **automatiquement créés et configurés** dans le `XupleManager`, sans aucune intervention manuelle du code appelant.

---

## 📋 Analyse Préliminaire

### 1. Comprendre le Flux Actuel

**Fichiers clés à examiner :**

```
tsd/internal/constraint/
├── ast.go                    # Définitions des nœuds AST
├── parser.go                 # Parser PEG actuel
└── rete_converter.go         # Conversion AST → RETE

tsd/internal/rete/
├── network.go                # Réseau RETE
└── constraint_pipeline.go    # ConstraintPipeline (point d'entrée actuel)

tsd/api/
├── pipeline.go               # Pipeline API (créé dans Prompt 02)
└── config.go                 # Configuration

tsd/xuples/
├── manager.go                # XupleManager
└── xuplespace.go             # XupleSpace
```

**Questions à résoudre :**

1. **Où sont stockées les définitions de xuple-space après parsing ?**
   - Actuellement dans `ConstraintPipeline.xupleSpaceDefinitions`
   - Structure : `map[string]map[string]interface{}`

2. **À quel moment du cycle de vie les créer ?**
   - Après le parsing, avant la propagation des faits initiaux
   - Dans `Pipeline.IngestFile()` après `retePipeline.IngestFile()`

3. **Comment mapper les propriétés TSD vers les types Go ?**
   - `selection` → `xuples.SelectionPolicy`
   - `consumption` → `xuples.ConsumptionPolicy`
   - `retention` → `xuples.RetentionPolicy`
   - `max-size` → `int`
   - `retention-duration` → `time.Duration`

### 2. Vérifier l'AST Actuel pour Xuple-Spaces

**Examiner `ast.go` :**

```go
// Chercher si XupleSpaceNode existe déjà
type XupleSpaceNode struct {
    Name       string
    Properties map[string]interface{}
    Location   SourceLocation
}
```

**Si absent, il faudra l'ajouter** (partie de ce prompt).

### 3. Comprendre la Conversion AST → RETE

**Examiner `rete_converter.go` :**

Le converter doit :
1. Détecter les `XupleSpaceNode` dans l'AST
2. Collecter leurs définitions
3. Les stocker dans `ConstraintPipeline.xupleSpaceDefinitions`

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Définir le Nœud AST pour Xuple-Space

**Fichier:** `tsd/internal/constraint/ast.go`

**Objectif:** Ajouter un nœud AST représentant une définition de xuple-space.

#### 1.1 Structure du Nœud

```go
// XupleSpaceNode représente une définition de xuple-space dans le fichier TSD.
//
// Exemple TSD:
//   xuple-space alerts {
//       selection: fifo,
//       consumption: once,
//       retention: 24h
//   }
type XupleSpaceNode struct {
    Name       string                 // Nom du xuple-space (ex: "alerts")
    Properties map[string]interface{} // Propriétés de configuration
    Location   SourceLocation         // Position dans le fichier source
}

// Implement Node interface
func (n *XupleSpaceNode) node() {}
func (n *XupleSpaceNode) GetLocation() SourceLocation {
    return n.Location
}
```

#### 1.2 Propriétés Supportées

| Propriété TSD        | Type Go           | Description                                    | Valeur par défaut |
|----------------------|-------------------|------------------------------------------------|-------------------|
| `selection`          | `string`          | Politique de sélection (fifo/lifo/random)      | `"fifo"`          |
| `consumption`        | `string`          | Politique de consommation (once/per-agent)     | `"once"`          |
| `retention`          | `string`/`duration` | Politique de rétention (unlimited/duration)  | `"unlimited"`     |
| `retention-duration` | `duration`        | Durée de rétention (ex: "24h", "7d")           | `0` (unlimited)   |
| `max-size`           | `int`             | Taille maximale du xuple-space                 | `0` (unlimited)   |

#### 1.3 Validation des Propriétés

```go
// ValidateXupleSpaceProperties valide les propriétés d'un xuple-space.
// Retourne une erreur si une propriété est invalide ou manquante.
func ValidateXupleSpaceProperties(name string, props map[string]interface{}) error {
    // Vérifier les valeurs de selection
    if sel, ok := props["selection"]; ok {
        selStr, ok := sel.(string)
        if !ok {
            return fmt.Errorf("xuple-space '%s': property 'selection' must be a string", name)
        }
        if selStr != "fifo" && selStr != "lifo" && selStr != "random" {
            return fmt.Errorf("xuple-space '%s': invalid selection policy '%s' (must be fifo, lifo, or random)", name, selStr)
        }
    }

    // Vérifier les valeurs de consumption
    if cons, ok := props["consumption"]; ok {
        consStr, ok := cons.(string)
        if !ok {
            return fmt.Errorf("xuple-space '%s': property 'consumption' must be a string", name)
        }
        if consStr != "once" && consStr != "per-agent" {
            return fmt.Errorf("xuple-space '%s': invalid consumption policy '%s' (must be once or per-agent)", name, consStr)
        }
    }

    // Vérifier les valeurs de retention
    if ret, ok := props["retention"]; ok {
        switch v := ret.(type) {
        case string:
            if v != "unlimited" {
                // Tenter de parser comme durée
                if _, err := parseDuration(v); err != nil {
                    return fmt.Errorf("xuple-space '%s': invalid retention value '%s': %w", name, v, err)
                }
            }
        case time.Duration:
            // OK
        default:
            return fmt.Errorf("xuple-space '%s': property 'retention' must be a string or duration", name)
        }
    }

    // Vérifier retention-duration
    if dur, ok := props["retention-duration"]; ok {
        switch v := dur.(type) {
        case string:
            if _, err := parseDuration(v); err != nil {
                return fmt.Errorf("xuple-space '%s': invalid retention-duration '%s': %w", name, v, err)
            }
        case time.Duration:
            // OK
        default:
            return fmt.Errorf("xuple-space '%s': property 'retention-duration' must be a string or duration", name)
        }
    }

    // Vérifier max-size
    if maxSize, ok := props["max-size"]; ok {
        switch v := maxSize.(type) {
        case int:
            if v < 0 {
                return fmt.Errorf("xuple-space '%s': max-size must be >= 0", name)
            }
        case float64:
            if v < 0 {
                return fmt.Errorf("xuple-space '%s': max-size must be >= 0", name)
            }
            props["max-size"] = int(v) // Convertir en int
        default:
            return fmt.Errorf("xuple-space '%s': property 'max-size' must be an integer", name)
        }
    }

    return nil
}

// parseDuration parse une chaîne de durée au format Go (24h, 7d, etc.)
// Supporte également des extensions comme "d" pour jours.
func parseDuration(s string) (time.Duration, error) {
    // Convertir "d" (jours) en "h" (heures)
    if strings.HasSuffix(s, "d") {
        days := strings.TrimSuffix(s, "d")
        d, err := strconv.Atoi(days)
        if err != nil {
            return 0, fmt.Errorf("invalid duration: %w", err)
        }
        return time.Duration(d) * 24 * time.Hour, nil
    }
    return time.ParseDuration(s)
}
```

---

### Tâche 2: Étendre le Parser

**Fichier:** `tsd/internal/constraint/parser.go`

**Objectif:** Parser la syntaxe `xuple-space <name> { ... }` et créer un `XupleSpaceNode`.

#### 2.1 Grammaire PEG (ajout)

```peg
Program         <- (TypeDef / RuleDef / XupleSpaceDef / FactDecl)* EOF

XupleSpaceDef   <- 'xuple-space' SPACE Identifier SPACE '{' XupleProps '}' SPACE*

XupleProps      <- (XupleProp (',' XupleProp)*)? ','?

XupleProp       <- SPACE* Identifier SPACE* ':' SPACE* XuplePropValue SPACE*

XuplePropValue  <- DurationLiteral / StringLiteral / IntegerLiteral / Identifier
```

#### 2.2 Fonction de Parsing

```go
// parseXupleSpace parse une définition de xuple-space.
//
// Syntaxe:
//   xuple-space <name> {
//       property1: value1,
//       property2: value2,
//       ...
//   }
func (p *Parser) parseXupleSpace() (*XupleSpaceNode, error) {
    start := p.pos

    // Consommer 'xuple-space'
    if !p.expectKeyword("xuple-space") {
        return nil, p.error("expected 'xuple-space'")
    }

    p.skipWhitespace()

    // Parser le nom
    name, err := p.parseIdentifier()
    if err != nil {
        return nil, p.wrapError("expected xuple-space name", err)
    }

    p.skipWhitespace()

    // Consommer '{'
    if !p.expect('{') {
        return nil, p.error("expected '{'")
    }

    // Parser les propriétés
    properties, err := p.parseXupleSpaceProperties()
    if err != nil {
        return nil, p.wrapError(fmt.Sprintf("in xuple-space '%s'", name), err)
    }

    // Consommer '}'
    if !p.expect('}') {
        return nil, p.error("expected '}'")
    }

    // Valider les propriétés
    if err := ValidateXupleSpaceProperties(name, properties); err != nil {
        return nil, err
    }

    return &XupleSpaceNode{
        Name:       name,
        Properties: properties,
        Location:   p.locationFrom(start),
    }, nil
}

// parseXupleSpaceProperties parse les propriétés d'un xuple-space.
//
// Format:
//   property1: value1,
//   property2: value2
func (p *Parser) parseXupleSpaceProperties() (map[string]interface{}, error) {
    props := make(map[string]interface{})

    p.skipWhitespace()

    // Si '}' immédiat, xuple-space vide (utilise les défauts)
    if p.peek() == '}' {
        return props, nil
    }

    for {
        p.skipWhitespace()

        // Vérifier fin de propriétés
        if p.peek() == '}' {
            break
        }

        // Parser le nom de la propriété
        propName, err := p.parseIdentifier()
        if err != nil {
            return nil, p.wrapError("expected property name", err)
        }

        p.skipWhitespace()

        // Consommer ':'
        if !p.expect(':') {
            return nil, p.errorf("expected ':' after property name '%s'", propName)
        }

        p.skipWhitespace()

        // Parser la valeur
        propValue, err := p.parseXupleSpacePropertyValue()
        if err != nil {
            return nil, p.wrapError(fmt.Sprintf("parsing value for property '%s'", propName), err)
        }

        // Vérifier que la propriété n'est pas dupliquée
        if _, exists := props[propName]; exists {
            return nil, p.errorf("duplicate property '%s'", propName)
        }

        props[propName] = propValue

        p.skipWhitespace()

        // Virgule optionnelle
        if p.peek() == ',' {
            p.advance()
            p.skipWhitespace()
        }
    }

    return props, nil
}

// parseXupleSpacePropertyValue parse une valeur de propriété.
// Supporte: durées ("24h"), entiers (1000), chaînes ("fifo"), identifiants (fifo).
func (p *Parser) parseXupleSpacePropertyValue() (interface{}, error) {
    p.skipWhitespace()

    c := p.peek()

    // Chaîne littérale
    if c == '"' || c == '\'' {
        return p.parseStringLiteral()
    }

    // Entier ou identifiant/durée
    start := p.pos
    for p.pos < len(p.input) && (isAlphaNumeric(p.input[p.pos]) || p.input[p.pos] == '-') {
        p.pos++
    }

    if p.pos == start {
        return nil, p.error("expected property value")
    }

    value := p.input[start:p.pos]

    // Tenter de parser comme entier
    if intVal, err := strconv.Atoi(value); err == nil {
        return intVal, nil
    }

    // Tenter de parser comme durée
    if dur, err := parseDuration(value); err == nil {
        return dur, nil
    }

    // Sinon, retourner comme chaîne (ex: "fifo", "once")
    return value, nil
}
```

#### 2.3 Intégration dans `parseProgram`

```go
// Modifier parseProgram pour appeler parseXupleSpace
func (p *Parser) parseProgram() (*ProgramNode, error) {
    nodes := []Node{}

    for p.pos < len(p.input) {
        p.skipWhitespace()
        if p.pos >= len(p.input) {
            break
        }

        // Détection du type de définition
        if p.peekKeyword("type") {
            node, err := p.parseTypeDef()
            if err != nil {
                return nil, err
            }
            nodes = append(nodes, node)
        } else if p.peekKeyword("rule") {
            node, err := p.parseRuleDef()
            if err != nil {
                return nil, err
            }
            nodes = append(nodes, node)
        } else if p.peekKeyword("xuple-space") {
            node, err := p.parseXupleSpace()
            if err != nil {
                return nil, err
            }
            nodes = append(nodes, node)
        } else {
            // Tentative de parser un fait
            node, err := p.parseFactDecl()
            if err != nil {
                return nil, err
            }
            nodes = append(nodes, node)
        }
    }

    return &ProgramNode{Nodes: nodes}, nil
}
```

---

### Tâche 3: Conversion AST → Définitions

**Fichier:** `tsd/internal/constraint/rete_converter.go`

**Objectif:** Extraire les définitions de xuple-spaces de l'AST et les stocker dans le pipeline.

#### 3.1 Méthode de Conversion

```go
// convertXupleSpace extrait la définition d'un xuple-space et l'ajoute
// à la liste des xuple-spaces à créer.
func (c *ASTConverter) convertXupleSpace(node *XupleSpaceNode) error {
    if c.pipeline.xupleSpaceDefinitions == nil {
        c.pipeline.xupleSpaceDefinitions = make(map[string]map[string]interface{})
    }

    // Vérifier que le xuple-space n'existe pas déjà
    if _, exists := c.pipeline.xupleSpaceDefinitions[node.Name]; exists {
        return c.errorf(node, "duplicate xuple-space definition '%s'", node.Name)
    }

    // Copier les propriétés (pour éviter les mutations)
    props := make(map[string]interface{})
    for k, v := range node.Properties {
        props[k] = v
    }

    c.pipeline.xupleSpaceDefinitions[node.Name] = props

    c.logf("Registered xuple-space definition: %s (properties: %v)", node.Name, props)

    return nil
}
```

#### 3.2 Intégration dans `Convert`

```go
// Modifier la méthode Convert pour gérer XupleSpaceNode
func (c *ASTConverter) Convert(program *ProgramNode) error {
    for _, node := range program.Nodes {
        switch n := node.(type) {
        case *TypeDefNode:
            if err := c.convertTypeDef(n); err != nil {
                return err
            }
        case *RuleDefNode:
            if err := c.convertRuleDef(n); err != nil {
                return err
            }
        case *XupleSpaceNode:
            if err := c.convertXupleSpace(n); err != nil {
                return err
            }
        case *FactDeclNode:
            if err := c.convertFactDecl(n); err != nil {
                return err
            }
        default:
            return c.errorf(nil, "unknown node type: %T", n)
        }
    }
    return nil
}
```

---

### Tâche 4: Création Automatique dans le Pipeline API

**Fichier:** `tsd/api/pipeline.go`

**Objectif:** Après ingestion du fichier TSD, créer automatiquement les xuple-spaces définis.

#### 4.1 Méthode de Création Automatique

```go
// createXupleSpaces crée automatiquement tous les xuple-spaces définis
// dans le fichier TSD ingéré.
//
// Cette méthode est appelée après l'ingestion du fichier TSD par le
// ConstraintPipeline, mais avant la propagation des faits initiaux.
func (p *Pipeline) createXupleSpaces() error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.retePipeline == nil {
        return fmt.Errorf("RETE pipeline not initialized")
    }

    // Récupérer les définitions de xuple-spaces depuis le ConstraintPipeline
    definitions := p.retePipeline.GetXupleSpaceDefinitions()

    if len(definitions) == 0 {
        // Aucun xuple-space défini, ce n'est pas une erreur
        return nil
    }

    if p.xupleManager == nil {
        return fmt.Errorf("xuple manager not initialized")
    }

    // Créer chaque xuple-space
    for name, props := range definitions {
        if err := p.createXupleSpaceFromDefinition(name, props); err != nil {
            return &XupleSpaceError{
                SpaceName: name,
                Operation: "create",
                Message:   "failed to create xuple-space",
                Cause:     err,
            }
        }
    }

    return nil
}

// createXupleSpaceFromDefinition crée un xuple-space à partir de sa définition.
func (p *Pipeline) createXupleSpaceFromDefinition(name string, props map[string]interface{}) error {
    // Appliquer les valeurs par défaut depuis la config
    defaults := p.config.XupleSpaceDefaults

    // Parser la politique de sélection
    selectionPolicy := defaults.Selection
    if sel, ok := props["selection"]; ok {
        if selStr, ok := sel.(string); ok {
            selectionPolicy = parseSelectionPolicy(selStr)
        }
    }

    // Parser la politique de consommation
    consumptionPolicy := defaults.Consumption
    if cons, ok := props["consumption"]; ok {
        if consStr, ok := cons.(string); ok {
            consumptionPolicy = parseConsumptionPolicy(consStr)
        }
    }

    // Parser la politique de rétention
    retentionPolicy := defaults.Retention
    retentionDuration := defaults.RetentionDuration

    if ret, ok := props["retention"]; ok {
        switch v := ret.(type) {
        case string:
            if v == "unlimited" {
                retentionPolicy = xuples.RetentionUnlimited
                retentionDuration = 0
            } else {
                // Parser comme durée
                dur, err := parseDuration(v)
                if err != nil {
                    return fmt.Errorf("invalid retention duration: %w", err)
                }
                retentionPolicy = xuples.RetentionDuration
                retentionDuration = dur
            }
        case time.Duration:
            retentionPolicy = xuples.RetentionDuration
            retentionDuration = v
        }
    }

    // Propriété retention-duration explicite (prioritaire)
    if dur, ok := props["retention-duration"]; ok {
        switch v := dur.(type) {
        case string:
            d, err := parseDuration(v)
            if err != nil {
                return fmt.Errorf("invalid retention-duration: %w", err)
            }
            retentionPolicy = xuples.RetentionDuration
            retentionDuration = d
        case time.Duration:
            retentionPolicy = xuples.RetentionDuration
            retentionDuration = v
        }
    }

    // Parser max-size
    maxSize := defaults.MaxSize
    if ms, ok := props["max-size"]; ok {
        switch v := ms.(type) {
        case int:
            maxSize = v
        case float64:
            maxSize = int(v)
        }
    }

    // Créer le xuple-space
    space, err := p.xupleManager.CreateXupleSpace(
        name,
        selectionPolicy,
        consumptionPolicy,
        retentionPolicy,
    )
    if err != nil {
        return err
    }

    // Configurer la durée de rétention si applicable
    if retentionPolicy == xuples.RetentionDuration && retentionDuration > 0 {
        space.SetRetentionDuration(retentionDuration)
    }

    // Configurer la taille maximale si spécifiée
    if maxSize > 0 {
        space.SetMaxSize(maxSize)
    }

    return nil
}
```

#### 4.2 Modification de `IngestFile` pour Appeler `createXupleSpaces`

```go
// IngestFile modifié pour créer automatiquement les xuple-spaces
func (p *Pipeline) IngestFile(filepath string) (*Result, error) {
    startTime := time.Now()

    // 1. Parse et build le réseau RETE
    parseStart := time.Now()
    if err := p.retePipeline.IngestFile(filepath); err != nil {
        return nil, p.wrapError("parse", err)
    }
    parseDuration := time.Since(parseStart)

    // 2. Créer automatiquement les xuple-spaces définis
    createSpacesStart := time.Now()
    if err := p.createXupleSpaces(); err != nil {
        return nil, p.wrapError("create-xuple-spaces", err)
    }
    createSpacesDuration := time.Since(createSpacesStart)

    // 3. Propagation (si des faits initiaux existent)
    propagateStart := time.Now()
    propagationCount := 0
    actionCount := 0
    // TODO: récupérer les compteurs réels depuis le network
    propagateDuration := time.Since(propagateStart)

    totalDuration := time.Since(startTime)

    // Construire le résultat
    result := &Result{
        network:      p.network,
        xupleManager: p.xupleManager,
        metrics: Metrics{
            TotalDuration:       totalDuration,
            ParseDuration:       parseDuration,
            BuildDuration:       parseDuration, // Approximation
            PropagationDuration: propagateDuration,
            TypeCount:           len(p.network.GetTypes()),
            RuleCount:           len(p.network.GetRules()),
            FactCount:           p.network.GetFactCount(),
            XupleSpaceCount:     len(p.xupleManager.ListSpaces()),
            PropagationCount:    propagationCount,
            ActionCount:         actionCount,
        },
    }

    return result, nil
}
```

#### 4.3 Ajout de `GetXupleSpaceDefinitions` dans `ConstraintPipeline`

**Fichier:** `tsd/internal/rete/constraint_pipeline.go`

```go
// GetXupleSpaceDefinitions retourne les définitions de xuple-spaces
// extraites du fichier TSD lors du parsing.
func (cp *ConstraintPipeline) GetXupleSpaceDefinitions() map[string]map[string]interface{} {
    cp.mu.Lock()
    defer cp.mu.Unlock()

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
```

---

### Tâche 5: Utilitaires de Mapping des Politiques

**Fichier:** `tsd/api/xuples_util.go`

**Objectif:** Mapper les chaînes TSD vers les constantes `xuples`.

```go
package api

import (
    "fmt"
    "time"

    "github.com/resinsec/tsd/xuples"
)

// parseSelectionPolicy convertit une chaîne en SelectionPolicy.
func parseSelectionPolicy(s string) xuples.SelectionPolicy {
    switch s {
    case "fifo":
        return xuples.SelectionFIFO
    case "lifo":
        return xuples.SelectionLIFO
    case "random":
        return xuples.SelectionRandom
    default:
        // Défaut: FIFO
        return xuples.SelectionFIFO
    }
}

// parseConsumptionPolicy convertit une chaîne en ConsumptionPolicy.
func parseConsumptionPolicy(s string) xuples.ConsumptionPolicy {
    switch s {
    case "once":
        return xuples.ConsumptionOnce
    case "per-agent":
        return xuples.ConsumptionPerAgent
    default:
        // Défaut: Once
        return xuples.ConsumptionOnce
    }
}

// parseRetentionPolicy convertit une chaîne en RetentionPolicy.
func parseRetentionPolicy(s string) xuples.RetentionPolicy {
    switch s {
    case "unlimited":
        return xuples.RetentionUnlimited
    case "duration":
        return xuples.RetentionDuration
    default:
        // Défaut: Unlimited
        return xuples.RetentionUnlimited
    }
}

// parseDuration parse une durée avec support des extensions (jours, etc.)
func parseDuration(s string) (time.Duration, error) {
    // Support des jours ("7d" -> 7 * 24h)
    if len(s) > 1 && s[len(s)-1] == 'd' {
        daysStr := s[:len(s)-1]
        var days int
        if _, err := fmt.Sscanf(daysStr, "%d", &days); err != nil {
            return 0, fmt.Errorf("invalid duration '%s': %w", s, err)
        }
        return time.Duration(days) * 24 * time.Hour, nil
    }

    // Déléguer au parser standard de Go
    return time.ParseDuration(s)
}

// FormatSelectionPolicy convertit une SelectionPolicy en chaîne.
func FormatSelectionPolicy(policy xuples.SelectionPolicy) string {
    switch policy {
    case xuples.SelectionFIFO:
        return "fifo"
    case xuples.SelectionLIFO:
        return "lifo"
    case xuples.SelectionRandom:
        return "random"
    default:
        return "unknown"
    }
}

// FormatConsumptionPolicy convertit une ConsumptionPolicy en chaîne.
func FormatConsumptionPolicy(policy xuples.ConsumptionPolicy) string {
    switch policy {
    case xuples.ConsumptionOnce:
        return "once"
    case xuples.ConsumptionPerAgent:
        return "per-agent"
    default:
        return "unknown"
    }
}

// FormatRetentionPolicy convertit une RetentionPolicy en chaîne.
func FormatRetentionPolicy(policy xuples.RetentionPolicy) string {
    switch policy {
    case xuples.RetentionUnlimited:
        return "unlimited"
    case xuples.RetentionDuration:
        return "duration"
    default:
        return "unknown"
    }
}
```

---

## 🧪 Tests à Implémenter

### Test 1: Parser - Xuple-Space Simple

**Fichier:** `tsd/internal/constraint/parser_xuplespace_test.go`

```go
func TestParser_XupleSpace_Simple(t *testing.T) {
    input := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}
`

    parser := NewParser(input)
    program, err := parser.Parse()
    require.NoError(t, err)
    require.Len(t, program.Nodes, 1)

    xupleSpace, ok := program.Nodes[0].(*XupleSpaceNode)
    require.True(t, ok)
    assert.Equal(t, "alerts", xupleSpace.Name)
    assert.Equal(t, "fifo", xupleSpace.Properties["selection"])
    assert.Equal(t, "once", xupleSpace.Properties["consumption"])
}
```

### Test 2: Parser - Propriétés avec Durées

```go
func TestParser_XupleSpace_WithDuration(t *testing.T) {
    input := `
xuple-space notifications {
    retention: 24h,
    max-size: 1000
}
`

    parser := NewParser(input)
    program, err := parser.Parse()
    require.NoError(t, err)
    require.Len(t, program.Nodes, 1)

    xupleSpace, ok := program.Nodes[0].(*XupleSpaceNode)
    require.True(t, ok)
    assert.Equal(t, "notifications", xupleSpace.Name)

    // Vérifier que la durée a été parsée
    retention, ok := xupleSpace.Properties["retention"].(time.Duration)
    require.True(t, ok)
    assert.Equal(t, 24*time.Hour, retention)

    assert.Equal(t, 1000, xupleSpace.Properties["max-size"])
}
```

### Test 3: Parser - Erreurs de Validation

```go
func TestParser_XupleSpace_InvalidSelection(t *testing.T) {
    input := `
xuple-space bad {
    selection: invalid_policy
}
`

    parser := NewParser(input)
    _, err := parser.Parse()
    require.Error(t, err)
    assert.Contains(t, err.Error(), "invalid selection policy")
}

func TestParser_XupleSpace_DuplicateProperty(t *testing.T) {
    input := `
xuple-space dup {
    selection: fifo,
    selection: lifo
}
`

    parser := NewParser(input)
    _, err := parser.Parse()
    require.Error(t, err)
    assert.Contains(t, err.Error(), "duplicate property")
}
```

### Test 4: Conversion - Extraction des Définitions

**Fichier:** `tsd/internal/constraint/converter_xuplespace_test.go`

```go
func TestConverter_XupleSpace_ExtractDefinitions(t *testing.T) {
    input := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}

xuple-space logs {
    retention: 7d,
    max-size: 5000
}
`

    parser := NewParser(input)
    program, err := parser.Parse()
    require.NoError(t, err)

    network := rete.NewNetwork()
    pipeline := &rete.ConstraintPipeline{Network: network}
    converter := NewASTConverter(pipeline)

    err = converter.Convert(program)
    require.NoError(t, err)

    defs := pipeline.GetXupleSpaceDefinitions()
    require.Len(t, defs, 2)

    // Vérifier "alerts"
    alertDef, exists := defs["alerts"]
    require.True(t, exists)
    assert.Equal(t, "fifo", alertDef["selection"])
    assert.Equal(t, "once", alertDef["consumption"])

    // Vérifier "logs"
    logsDef, exists := defs["logs"]
    require.True(t, exists)
    assert.Equal(t, 7*24*time.Hour, logsDef["retention"])
    assert.Equal(t, 5000, logsDef["max-size"])
}
```

### Test 5: API Pipeline - Création Automatique

**Fichier:** `tsd/api/pipeline_xuplespace_test.go`

```go
func TestPipeline_AutoCreateXupleSpaces(t *testing.T) {
    // Créer un fichier TSD temporaire
    tsdContent := `
xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

xuple-space notifications {
    selection: random,
    max-size: 1000
}

type Alert {
    id: string,
    message: string
}
`

    tmpfile, err := os.CreateTemp("", "test*.tsd")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())

    _, err = tmpfile.WriteString(tsdContent)
    require.NoError(t, err)
    tmpfile.Close()

    // Créer le pipeline
    pipeline, err := NewPipeline()
    require.NoError(t, err)

    // Ingérer le fichier
    result, err := pipeline.IngestFile(tmpfile.Name())
    require.NoError(t, err)

    // Vérifier que les xuple-spaces ont été créés
    spaces := result.XupleSpaceNames()
    require.Len(t, spaces, 2)
    assert.Contains(t, spaces, "alerts")
    assert.Contains(t, spaces, "notifications")

    // Vérifier les propriétés du xuple-space "alerts"
    alertSpace := result.XupleManager().GetSpace("alerts")
    require.NotNil(t, alertSpace)
    assert.Equal(t, xuples.SelectionFIFO, alertSpace.GetSelectionPolicy())
    assert.Equal(t, xuples.ConsumptionOnce, alertSpace.GetConsumptionPolicy())
    assert.Equal(t, xuples.RetentionDuration, alertSpace.GetRetentionPolicy())

    // Vérifier "notifications"
    notifSpace := result.XupleManager().GetSpace("notifications")
    require.NotNil(t, notifSpace)
    assert.Equal(t, xuples.SelectionRandom, notifSpace.GetSelectionPolicy())
    assert.Equal(t, 1000, notifSpace.GetMaxSize())
}
```

### Test 6: E2E - Xuple-Spaces et Actions Xuple

```go
func TestE2E_XupleSpaceAutoCreation_WithXupleAction(t *testing.T) {
    tsdContent := `
xuple-space alerts {
    selection: fifo,
    consumption: once
}

type Temperature {
    sensorId: string,
    value: float
}

type Alert {
    sensorId: string,
    message: string,
    temp: float
}

rule HighTemperature {
    when {
        t: Temperature(value > 30.0)
    }
    then {
        Xuple("alerts", Alert(
            sensorId: t.sensorId,
            message: "High temperature detected",
            temp: t.value
        ))
    }
}
`

    tmpfile, err := os.CreateTemp("", "test*.tsd")
    require.NoError(t, err)
    defer os.Remove(tmpfile.Name())

    _, err = tmpfile.WriteString(tsdContent)
    require.NoError(t, err)
    tmpfile.Close()

    // Créer le pipeline
    pipeline, err := NewPipeline()
    require.NoError(t, err)

    // Ingérer le fichier
    result, err := pipeline.IngestFile(tmpfile.Name())
    require.NoError(t, err)

    // Vérifier que le xuple-space "alerts" a été créé
    spaces := result.XupleSpaceNames()
    require.Contains(t, spaces, "alerts")

    // Soumettre un fait Temperature
    tempFact := result.Network().CreateFact("Temperature", map[string]interface{}{
        "sensorId": "sensor-01",
        "value":    35.5,
    })
    result.Network().Assert(tempFact)

    // Vérifier qu'un xuple a été créé dans "alerts"
    xuples := result.GetXuples("alerts")
    require.Len(t, xuples, 1)

    alert := xuples[0]
    assert.Equal(t, "Alert", alert.Type())
    assert.Equal(t, "sensor-01", alert.Get("sensorId"))
    assert.Equal(t, "High temperature detected", alert.Get("message"))
    assert.Equal(t, 35.5, alert.Get("temp"))
}
```

---

## ✅ Checklist de Validation

### Parser

- [ ] Grammaire PEG étendue avec `XupleSpaceDef`
- [ ] `parseXupleSpace()` implémentée et testée
- [ ] `parseXupleSpaceProperties()` gère toutes les propriétés supportées
- [ ] `parseXupleSpacePropertyValue()` parse durées, entiers, chaînes
- [ ] Validation des propriétés (`ValidateXupleSpaceProperties`)
- [ ] Support des durées avec extension "d" (jours)
- [ ] Erreurs claires pour propriétés invalides ou manquantes
- [ ] Tests unitaires pour syntaxe valide
- [ ] Tests unitaires pour syntaxe invalide

### AST

- [ ] `XupleSpaceNode` ajouté à `ast.go`
- [ ] Interface `Node` implémentée correctement
- [ ] `GetLocation()` retourne la position source
- [ ] Propriétés stockées dans `map[string]interface{}`

### Conversion

- [ ] `convertXupleSpace()` extrait définitions
- [ ] Définitions stockées dans `ConstraintPipeline.xupleSpaceDefinitions`
- [ ] Détection des doublons (xuple-space défini deux fois)
- [ ] `GetXupleSpaceDefinitions()` retourne copie (immutabilité)
- [ ] Tests de conversion avec multiples xuple-spaces

### API Pipeline

- [ ] `createXupleSpaces()` appelée après ingestion
- [ ] `createXupleSpaceFromDefinition()` mappe propriétés → xuples API
- [ ] Application des valeurs par défaut depuis `Config.XupleSpaceDefaults`
- [ ] Support de toutes les politiques (selection, consumption, retention)
- [ ] Configuration de `RetentionDuration` et `MaxSize`
- [ ] Gestion d'erreurs claire (XupleSpaceError)
- [ ] Tests unitaires pour création automatique
- [ ] Tests E2E avec actions Xuple

### Utilitaires

- [ ] `parseSelectionPolicy()` implémentée
- [ ] `parseConsumptionPolicy()` implémentée
- [ ] `parseRetentionPolicy()` implémentée
- [ ] `parseDuration()` supporte "d" (jours)
- [ ] Fonctions de formatage (pour debug/logs)

### Standards

- [ ] Code formaté (`gofmt`)
- [ ] Pas de warnings du linter
- [ ] Couverture de tests > 80%
- [ ] Commentaires GoDoc complets
- [ ] Exemples d'utilisation dans GoDoc

---

## 📝 Documentation à Mettre à Jour

### 1. Guide TSD (`docs/TSD_LANGUAGE.md`)

Ajouter section sur la syntaxe `xuple-space` :

```markdown
## Xuple-Spaces

Les xuple-spaces sont des espaces de stockage pour les xuples générés par les règles.

### Syntaxe

\`\`\`tsd
xuple-space <name> {
    selection: <fifo|lifo|random>,
    consumption: <once|per-agent>,
    retention: <unlimited|duration>,
    retention-duration: <duration>,
    max-size: <integer>
}
\`\`\`

### Propriétés

- **selection**: Politique de sélection (défaut: `fifo`)
  - `fifo`: Premier arrivé, premier servi
  - `lifo`: Dernier arrivé, premier servi
  - `random`: Sélection aléatoire

- **consumption**: Politique de consommation (défaut: `once`)
  - `once`: Chaque xuple ne peut être consommé qu'une fois
  - `per-agent`: Chaque agent peut consommer le xuple une fois

- **retention**: Politique de rétention (défaut: `unlimited`)
  - `unlimited`: Les xuples sont conservés indéfiniment
  - `<duration>`: Durée de rétention (ex: `24h`, `7d`)

- **max-size**: Taille maximale du xuple-space (défaut: illimité)

### Exemples

\`\`\`tsd
xuple-space alerts {
    selection: fifo,
    consumption: once,
    retention: 24h
}

xuple-space notifications {
    selection: random,
    consumption: per-agent,
    max-size: 1000
}
\`\`\`
```

### 2. Guide API (`docs/API_USAGE.md`)

Ajouter section sur la création automatique :

```markdown
## Création Automatique des Xuple-Spaces

Les xuple-spaces définis dans le fichier TSD sont **automatiquement créés**
lors de l'ingestion du fichier via `Pipeline.IngestFile()`.

### Exemple

Fichier TSD (`rules.tsd`):
\`\`\`tsd
xuple-space alerts {
    selection: fifo,
    retention: 24h
}
\`\`\`

Code Go:
\`\`\`go
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")

// Le xuple-space "alerts" existe maintenant
spaces := result.XupleSpaceNames()
// => ["alerts"]
\`\`\`

Aucune configuration manuelle nécessaire !
```

### 3. Documentation Xuples E2E (`docs/XUPLES_E2E_AUTOMATIC.md`)

Mettre à jour la section "Création des Xuple-Spaces" pour indiquer que c'est automatique.

---

## 🎯 Résultat Attendu

### Avant (manuel)

```go
// Test E2E - avant
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")

// Création manuelle du xuple-space
manager := result.XupleManager()
manager.CreateXupleSpace("alerts", xuples.SelectionFIFO, xuples.ConsumptionOnce, xuples.RetentionUnlimited)

// Utilisation...
```

### Après (automatique)

```go
// Test E2E - après
pipeline, _ := api.NewPipeline()
result, _ := pipeline.IngestFile("rules.tsd")

// Le xuple-space "alerts" existe déjà !
xuples := result.GetXuples("alerts")
```

Le fichier TSD contient :

```tsd
xuple-space alerts {
    selection: fifo,
    consumption: once
}
```

**Aucune autre étape nécessaire.** ✨

---

## 🔗 Dépendances

### Entrantes

- ✅ Prompt 01 : Parser supporte faits inline (pour actions Xuple complètes)
- ✅ Prompt 02 : Package `api` existe avec Pipeline

### Sortantes

- ➡️ Prompt 04 : Automatisation des actions Xuple (utilisera les xuple-spaces créés)
- ➡️ Prompt 05 : Migration des tests E2E (bénéficiera de la création automatique)

---

## 🚀 Stratégie d'Implémentation

1. **Phase 1: Parser** (1-2h)
   - Ajouter `XupleSpaceNode` à l'AST
   - Implémenter `parseXupleSpace()` et helpers
   - Tests unitaires du parser

2. **Phase 2: Conversion** (30min)
   - Implémenter `convertXupleSpace()`
   - Ajouter `GetXupleSpaceDefinitions()` au pipeline
   - Tests de conversion

3. **Phase 3: API Pipeline** (1-2h)
   - Implémenter `createXupleSpaces()`
   - Modifier `IngestFile()` pour appeler création automatique
   - Utilitaires de mapping des politiques

4. **Phase 4: Tests E2E** (1h)
   - Tests d'intégration complets
   - Tests E2E avec actions Xuple

5. **Phase 5: Documentation** (30min)
   - Mise à jour des guides
   - Exemples GoDoc

**Estimation totale: 4-6 heures**

---

## 📊 Critères de Succès

- [ ] Parser reconnaît la syntaxe `xuple-space { ... }`
- [ ] Toutes les propriétés sont supportées et validées
- [ ] Conversion AST → définitions fonctionne
- [ ] `Pipeline.IngestFile()` crée automatiquement les xuple-spaces
- [ ] Aucune configuration manuelle nécessaire dans les tests
- [ ] Tests unitaires passent (couverture > 80%)
- [ ] Tests E2E passent avec actions Xuple
- [ ] Documentation à jour
- [ ] Pas de régression dans les tests existants

---

**FIN DU PROMPT 03**