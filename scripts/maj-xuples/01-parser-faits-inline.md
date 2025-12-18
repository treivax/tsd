# 🔧 Prompt 01 - Parser TSD: Support Faits Inline dans Actions

> **Objectif**: Étendre le parser TSD pour supporter complètement les faits inline dans les actions  
> **Dépendances**: Aucune (peut être exécuté en premier)  
> **Contexte max**: 128k tokens  
> **Durée estimée**: 1 session

---

## 🎯 Objectif

Actuellement, le parser TSD ne supporte pas complètement la création de faits inline dans les actions, notamment:
- La syntaxe multi-ligne `Xuple("space", Alert(...))`
- Les références aux champs des faits déclencheurs (ex: `s.sensorId`, `s.temperature`)
- Les faits inline imbriqués dans les appels d'actions

**Comportement actuel** (NON supporté):
```tsd
rule critical_alert: {s: Sensor} / s.temperature > 40.0 ==>
    Xuple("critical_alerts", Alert(
        level: "CRITICAL",
        message: "Temperature too high",
        sensorId: s.sensorId,
        temperature: s.temperature
    ))
```

**Comportement cible** (DOIT être supporté):
```tsd
// Simple (une ligne)
rule alert1: {s: Sensor} / s.temp > 40.0 ==> Xuple("alerts", Alert(level: "HIGH", id: s.id))

// Multi-ligne (formatage lisible)
rule alert2: {s: Sensor} / s.temp > 40.0 ==>
    Xuple("alerts", Alert(
        level: "CRITICAL",
        message: "Too hot",
        sensorId: s.sensorId,
        temperature: s.temperature,
        location: s.location
    ))

// Multiple actions
rule complex: {s: Sensor} / s.temp > 40.0 ==>
    Print("Alert!"),
    Xuple("alerts", Alert(level: "HIGH", id: s.id)),
    Xuple("commands", Command(action: "cool", target: s.location))

// Références aux variables
rule with_refs: {s: Sensor, a: Alert} / s.id == a.sensorId ==>
    Xuple("correlated", Event(
        sensor: s.id,
        alert: a.level,
        combined: s.temperature + a.priority
    ))
```

---

## 📋 Analyse Préliminaire

### 1. Identifier le Parser Actuel

Examiner la structure du parser:

```bash
# Localiser les fichiers du parser
find tsd/constraint -name "*.go" | grep -E "(parser|peg)"

# Identifier le type de parser
# - Parser PEG (Parsing Expression Grammar) ?
# - Parser hand-written ?
# - Parser généré (yacc, antlr, etc.) ?
```

**Questions à répondre**:
- [ ] Quel système de parsing est utilisé ? (PEG, hand-written, généré)
- [ ] Où sont définies les règles de grammaire ?
- [ ] Comment sont représentés les faits inline actuellement dans l'AST ?
- [ ] Existe-t-il déjà un support partiel qu'on peut étendre ?

### 2. Examiner l'AST Actuel

Lire les structures AST pour comprendre la représentation actuelle:

```bash
# Examiner les nœuds AST
cat tsd/constraint/ast.go | grep -A 5 "type.*Node"
```

**Questions à répondre**:
- [ ] Comment les actions sont-elles représentées dans l'AST ?
- [ ] Existe-t-il un nœud pour les appels de fonction avec arguments ?
- [ ] Comment les références aux variables sont-elles gérées ?
- [ ] Y a-t-il déjà un nœud pour les faits inline ?

### 3. Comprendre la Conversion AST → RETE

Examiner comment l'AST est converti en structures RETE:

```bash
# Chercher la conversion d'actions
grep -r "Action" tsd/rete/*.go | grep -i convert
```

**Questions à répondre**:
- [ ] Où se fait la conversion AST → Actions RETE ?
- [ ] Comment les paramètres d'actions sont-ils traités ?
- [ ] Comment résoudre les références aux variables (binding context) ?

---

## 🛠️ Tâches à Réaliser

### Tâche 1: Étendre la Grammaire

**Objectif**: Modifier la grammaire du parser pour accepter:
- Faits inline dans les paramètres d'actions
- Multi-ligne (avec espaces/retours à la ligne)
- Références aux champs de variables (`var.field`)

**Fichiers concernés**:
- `constraint/parser.peg` (si PEG)
- `constraint/parser.go` (si hand-written)
- Fichier de grammaire approprié selon le système utilisé

**Modifications attendues**:

```peg
// Exemple pour un parser PEG (adapter selon le système réel)

Action ← ActionCall / Print / ...

ActionCall ← Identifier "(" ArgumentList ")" Spacing

ArgumentList ← Argument ("," Argument)*

Argument ← InlineFact / FieldReference / StringLiteral / NumberLiteral / Identifier

InlineFact ← TypeName "(" FieldAssignmentList ")"

FieldAssignmentList ← FieldAssignment ("," Spacing? FieldAssignment)*

FieldAssignment ← Identifier ":" Expression

FieldReference ← Identifier "." Identifier ("." Identifier)*

Expression ← FieldReference / Literal / BinaryOp / ...
```

**Validation**:
```bash
# Regénérer le parser si nécessaire
cd tsd/constraint
go generate ./...

# Vérifier compilation
go build ./...
```

### Tâche 2: Étendre l'AST

**Objectif**: Ajouter les nœuds AST nécessaires pour représenter les nouvelles constructions

**Fichiers concernés**:
- `constraint/ast.go`

**Modifications attendues**:

```go
// Nœud pour un fait inline dans une action
type InlineFactNode struct {
    TypeName   string                    // Type du fait (ex: "Alert")
    Fields     map[string]ExpressionNode // Assignations de champs
    Location   SourceLocation            // Pour messages d'erreur
}

// Nœud pour une référence à un champ de variable
type FieldReferenceNode struct {
    Variable string   // Nom de la variable (ex: "s")
    Path     []string // Chemin de champs (ex: ["sensorId"])
    Location SourceLocation
}

// Expression (peut être littéral, référence, opération, etc.)
type ExpressionNode interface {
    isExpression()
    GetLocation() SourceLocation
}

// Implémentations
func (*InlineFactNode) isExpression()     {}
func (*FieldReferenceNode) isExpression() {}
func (*LiteralNode) isExpression()        {}
// ... etc.

// Nœud d'action modifié pour accepter des expressions
type ActionNode struct {
    Name      string           // "Xuple", "Print", etc.
    Arguments []ExpressionNode // Arguments (peuvent être faits inline, références, etc.)
    Location  SourceLocation
}
```

**Validation**:
```bash
# Compiler pour vérifier les interfaces
go build ./constraint/...
```

### Tâche 3: Implémenter le Parsing

**Objectif**: Implémenter la logique de parsing pour construire les nouveaux nœuds AST

**Fichiers concernés**:
- `constraint/parser.go` (ou fichier approprié)

**Modifications attendues**:

```go
// Parser une action avec arguments complexes
func (p *Parser) parseAction() (*ActionNode, error) {
    name := p.parseIdentifier()
    
    if !p.expect("(") {
        return nil, p.error("expected '(' after action name")
    }
    
    args := []ExpressionNode{}
    for !p.check(")") {
        arg, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        args = append(args, arg)
        
        if !p.check(")") {
            if !p.expect(",") {
                return nil, p.error("expected ',' or ')' after argument")
            }
            p.skipWhitespace() // Supporter multi-ligne
        }
    }
    
    p.expect(")")
    
    return &ActionNode{
        Name:      name,
        Arguments: args,
        Location:  p.currentLocation(),
    }, nil
}

// Parser une expression (peut être fait inline, référence, littéral, etc.)
func (p *Parser) parseExpression() (ExpressionNode, error) {
    // Essayer de parser un fait inline (TypeName(...))
    if p.isTypeStart() {
        return p.parseInlineFact()
    }
    
    // Essayer de parser une référence de champ (var.field)
    if p.isIdentifier() && p.peekAhead(".") {
        return p.parseFieldReference()
    }
    
    // Littéraux (string, number, etc.)
    return p.parseLiteral()
}

// Parser un fait inline: Alert(level: "HIGH", id: s.id)
func (p *Parser) parseInlineFact() (*InlineFactNode, error) {
    typeName := p.parseIdentifier()
    
    if !p.expect("(") {
        return nil, p.error("expected '(' after type name")
    }
    
    fields := make(map[string]ExpressionNode)
    
    for !p.check(")") {
        fieldName := p.parseIdentifier()
        
        if !p.expect(":") {
            return nil, p.error("expected ':' after field name")
        }
        
        p.skipWhitespace() // Supporter multi-ligne
        
        value, err := p.parseExpression()
        if err != nil {
            return nil, err
        }
        
        fields[fieldName] = value
        
        if !p.check(")") {
            if !p.expect(",") {
                return nil, p.error("expected ',' or ')' after field assignment")
            }
            p.skipWhitespace() // Supporter multi-ligne
        }
    }
    
    p.expect(")")
    
    return &InlineFactNode{
        TypeName: typeName,
        Fields:   fields,
        Location: p.currentLocation(),
    }, nil
}

// Parser une référence de champ: s.sensorId ou s.location.city
func (p *Parser) parseFieldReference() (*FieldReferenceNode, error) {
    variable := p.parseIdentifier()
    path := []string{}
    
    for p.expect(".") {
        field := p.parseIdentifier()
        path = append(path, field)
    }
    
    if len(path) == 0 {
        return nil, p.error("expected field name after '.'")
    }
    
    return &FieldReferenceNode{
        Variable: variable,
        Path:     path,
        Location: p.currentLocation(),
    }, nil
}
```

**Validation**:
```bash
# Tester le parsing de syntaxes simples
go test ./constraint/... -v -run TestParse
```

### Tâche 4: Conversion AST → RETE

**Objectif**: Convertir les nouveaux nœuds AST en structures utilisables par RETE

**Fichiers concernés**:
- `rete/ast_converter.go` (ou fichier approprié)
- `rete/actions.go` (pour la représentation des actions)

**Modifications attendues**:

```go
// Convertir une action AST en action RETE
func (c *ASTConverter) convertAction(node *constraint.ActionNode, bindingContext map[string]*BoundVariable) (*ActionDefinition, error) {
    args := make([]ActionArgument, len(node.Arguments))
    
    for i, argNode := range node.Arguments {
        arg, err := c.convertExpression(argNode, bindingContext)
        if err != nil {
            return nil, fmt.Errorf("erreur conversion argument %d: %w", i, err)
        }
        args[i] = arg
    }
    
    return &ActionDefinition{
        Name:      node.Name,
        Arguments: args,
    }, nil
}

// Convertir une expression AST en argument d'action
func (c *ASTConverter) convertExpression(node constraint.ExpressionNode, bindingContext map[string]*BoundVariable) (ActionArgument, error) {
    switch n := node.(type) {
    case *constraint.InlineFactNode:
        return c.convertInlineFact(n, bindingContext)
    case *constraint.FieldReferenceNode:
        return c.convertFieldReference(n, bindingContext)
    case *constraint.LiteralNode:
        return c.convertLiteral(n)
    default:
        return nil, fmt.Errorf("type d'expression non supporté: %T", node)
    }
}

// Convertir un fait inline
func (c *ASTConverter) convertInlineFact(node *constraint.InlineFactNode, bindingContext map[string]*BoundVariable) (*InlineFactArgument, error) {
    fields := make(map[string]ActionArgument)
    
    for fieldName, exprNode := range node.Fields {
        value, err := c.convertExpression(exprNode, bindingContext)
        if err != nil {
            return nil, fmt.Errorf("erreur champ '%s': %w", fieldName, err)
        }
        fields[fieldName] = value
    }
    
    return &InlineFactArgument{
        TypeName: node.TypeName,
        Fields:   fields,
    }, nil
}

// Convertir une référence de champ (avec résolution du binding)
func (c *ASTConverter) convertFieldReference(node *constraint.FieldReferenceNode, bindingContext map[string]*BoundVariable) (*FieldReferenceArgument, error) {
    // Vérifier que la variable existe dans le contexte
    boundVar, exists := bindingContext[node.Variable]
    if !exists {
        return nil, fmt.Errorf("variable '%s' non définie dans le contexte", node.Variable)
    }
    
    // Vérifier que le chemin de champs est valide pour le type
    if err := c.validateFieldPath(boundVar.Type, node.Path); err != nil {
        return nil, err
    }
    
    return &FieldReferenceArgument{
        Variable:  node.Variable,
        FieldPath: node.Path,
        BoundVar:  boundVar, // Référence au fait bound pour résolution runtime
    }, nil
}

// Types pour représenter les arguments d'actions
type ActionArgument interface {
    isActionArgument()
}

type InlineFactArgument struct {
    TypeName string
    Fields   map[string]ActionArgument
}

type FieldReferenceArgument struct {
    Variable  string
    FieldPath []string
    BoundVar  *BoundVariable
}

type LiteralArgument struct {
    Value interface{}
}

func (*InlineFactArgument) isActionArgument()     {}
func (*FieldReferenceArgument) isActionArgument() {}
func (*LiteralArgument) isActionArgument()        {}
```

**Validation**:
```bash
# Tester la conversion
go test ./rete/... -v -run TestConvertAction
```

### Tâche 5: Résolution Runtime des Références

**Objectif**: Lors de l'exécution d'une action, résoudre les références aux champs des faits déclencheurs

**Fichiers concernés**:
- `rete/action_executor.go` (ou fichier approprié)
- `rete/terminal_node.go` (exécution des actions)

**Modifications attendues**:

```go
// Résoudre un argument d'action en valeur concrète
func (e *ActionExecutor) resolveArgument(arg ActionArgument, context *ExecutionContext) (interface{}, error) {
    switch a := arg.(type) {
    case *InlineFactArgument:
        return e.resolveInlineFact(a, context)
    case *FieldReferenceArgument:
        return e.resolveFieldReference(a, context)
    case *LiteralArgument:
        return a.Value, nil
    default:
        return nil, fmt.Errorf("type d'argument inconnu: %T", arg)
    }
}

// Résoudre un fait inline (créer le fait avec champs résolus)
func (e *ActionExecutor) resolveInlineFact(arg *InlineFactArgument, context *ExecutionContext) (*Fact, error) {
    fields := make(map[string]interface{})
    
    for fieldName, fieldArg := range arg.Fields {
        value, err := e.resolveArgument(fieldArg, context)
        if err != nil {
            return nil, fmt.Errorf("erreur résolution champ '%s': %w", fieldName, err)
        }
        fields[fieldName] = value
    }
    
    // Créer le fait
    fact := &Fact{
        ID:     generateFactID(arg.TypeName),
        Type:   arg.TypeName,
        Fields: fields,
    }
    
    return fact, nil
}

// Résoudre une référence de champ (extraire la valeur du fait bound)
func (e *ActionExecutor) resolveFieldReference(arg *FieldReferenceArgument, context *ExecutionContext) (interface{}, error) {
    // Récupérer le fait bound depuis le contexte
    fact := context.GetBoundFact(arg.Variable)
    if fact == nil {
        return nil, fmt.Errorf("variable '%s' non trouvée dans le contexte", arg.Variable)
    }
    
    // Naviguer le chemin de champs
    value := fact.Fields
    for _, fieldName := range arg.FieldPath {
        // Extraire la valeur du champ
        if fieldMap, ok := value.(map[string]interface{}); ok {
            value = fieldMap[fieldName]
        } else {
            return nil, fmt.Errorf("impossible d'accéder au champ '%s'", fieldName)
        }
    }
    
    return value, nil
}

// Contexte d'exécution contenant les faits bound
type ExecutionContext struct {
    TriggeringFacts []*Fact
    Bindings        map[string]*Fact // variable name → fait bound
}

func (ctx *ExecutionContext) GetBoundFact(variable string) *Fact {
    return ctx.Bindings[variable]
}
```

**Validation**:
```bash
# Tester l'exécution d'actions avec références
go test ./rete/... -v -run TestActionExecution
```

---

## 🧪 Tests à Implémenter

### Test 1: Parser - Syntaxe Simple

**Fichier**: `constraint/parser_inline_test.go`

```go
func TestParser_InlineFact_Simple(t *testing.T) {
    input := `rule test: {s: Sensor} / s.temp > 40.0 ==> 
        Xuple("alerts", Alert(level: "HIGH", id: "A001"))`
    
    ast, err := ParseTSD(input)
    require.NoError(t, err)
    require.NotNil(t, ast)
    
    // Vérifier la structure
    require.Len(t, ast.Rules, 1)
    rule := ast.Rules[0]
    require.Len(t, rule.Actions, 1)
    
    action := rule.Actions[0]
    assert.Equal(t, "Xuple", action.Name)
    require.Len(t, action.Arguments, 2)
    
    // Premier argument: string "alerts"
    arg0, ok := action.Arguments[0].(*LiteralNode)
    require.True(t, ok)
    assert.Equal(t, "alerts", arg0.Value)
    
    // Deuxième argument: fait inline Alert(...)
    arg1, ok := action.Arguments[1].(*InlineFactNode)
    require.True(t, ok)
    assert.Equal(t, "Alert", arg1.TypeName)
    assert.Len(t, arg1.Fields, 2)
}
```

### Test 2: Parser - Syntaxe Multi-ligne

```go
func TestParser_InlineFact_Multiline(t *testing.T) {
    input := `rule test: {s: Sensor} / s.temp > 40.0 ==> 
        Xuple("alerts", Alert(
            level: "CRITICAL",
            message: "Temperature too high",
            sensorId: "S001",
            temperature: 45.5
        ))`
    
    ast, err := ParseTSD(input)
    require.NoError(t, err)
    
    action := ast.Rules[0].Actions[0]
    fact := action.Arguments[1].(*InlineFactNode)
    
    assert.Len(t, fact.Fields, 4)
    assert.Contains(t, fact.Fields, "level")
    assert.Contains(t, fact.Fields, "message")
    assert.Contains(t, fact.Fields, "sensorId")
    assert.Contains(t, fact.Fields, "temperature")
}
```

### Test 3: Parser - Références aux Variables

```go
func TestParser_FieldReference(t *testing.T) {
    input := `rule test: {s: Sensor} / s.temp > 40.0 ==> 
        Xuple("alerts", Alert(
            level: "HIGH",
            sensorId: s.sensorId,
            temperature: s.temperature
        ))`
    
    ast, err := ParseTSD(input)
    require.NoError(t, err)
    
    action := ast.Rules[0].Actions[0]
    fact := action.Arguments[1].(*InlineFactNode)
    
    // Vérifier référence s.sensorId
    sensorIdField := fact.Fields["sensorId"]
    ref, ok := sensorIdField.(*FieldReferenceNode)
    require.True(t, ok)
    assert.Equal(t, "s", ref.Variable)
    assert.Equal(t, []string{"sensorId"}, ref.Path)
    
    // Vérifier référence s.temperature
    tempField := fact.Fields["temperature"]
    ref2, ok := tempField.(*FieldReferenceNode)
    require.True(t, ok)
    assert.Equal(t, "s", ref2.Variable)
    assert.Equal(t, []string{"temperature"}, ref2.Path)
}
```

### Test 4: Parser - Actions Multiples

```go
func TestParser_MultipleActions(t *testing.T) {
    input := `rule test: {s: Sensor} / s.temp > 40.0 ==> 
        Print("Alert!"),
        Xuple("alerts", Alert(level: "HIGH", id: s.id)),
        Xuple("commands", Command(action: "cool", target: s.location))`
    
    ast, err := ParseTSD(input)
    require.NoError(t, err)
    
    rule := ast.Rules[0]
    assert.Len(t, rule.Actions, 3)
    assert.Equal(t, "Print", rule.Actions[0].Name)
    assert.Equal(t, "Xuple", rule.Actions[1].Name)
    assert.Equal(t, "Xuple", rule.Actions[2].Name)
}
```

### Test 5: Conversion AST → RETE

**Fichier**: `rete/ast_converter_inline_test.go`

```go
func TestConverter_InlineFact(t *testing.T) {
    // Créer un AST avec fait inline
    action := &constraint.ActionNode{
        Name: "Xuple",
        Arguments: []constraint.ExpressionNode{
            &constraint.LiteralNode{Value: "alerts"},
            &constraint.InlineFactNode{
                TypeName: "Alert",
                Fields: map[string]constraint.ExpressionNode{
                    "level": &constraint.LiteralNode{Value: "HIGH"},
                    "id":    &constraint.LiteralNode{Value: "A001"},
                },
            },
        },
    }
    
    // Convertir
    converter := NewASTConverter()
    bindingContext := make(map[string]*BoundVariable)
    
    actionDef, err := converter.convertAction(action, bindingContext)
    require.NoError(t, err)
    require.NotNil(t, actionDef)
    
    assert.Equal(t, "Xuple", actionDef.Name)
    require.Len(t, actionDef.Arguments, 2)
    
    // Vérifier le fait inline
    inlineFact, ok := actionDef.Arguments[1].(*InlineFactArgument)
    require.True(t, ok)
    assert.Equal(t, "Alert", inlineFact.TypeName)
    assert.Len(t, inlineFact.Fields, 2)
}
```

### Test 6: Résolution Runtime

**Fichier**: `rete/action_executor_inline_test.go`

```go
func TestExecutor_ResolveFieldReference(t *testing.T) {
    // Créer un fait sensor
    sensorFact := &Fact{
        ID:   "S001",
        Type: "Sensor",
        Fields: map[string]interface{}{
            "sensorId":    "S001",
            "temperature": 45.5,
            "location":    "RoomA",
        },
    }
    
    // Créer un contexte d'exécution
    ctx := &ExecutionContext{
        Bindings: map[string]*Fact{
            "s": sensorFact,
        },
    }
    
    // Créer une référence de champ
    ref := &FieldReferenceArgument{
        Variable:  "s",
        FieldPath: []string{"temperature"},
    }
    
    // Résoudre
    executor := NewActionExecutor()
    value, err := executor.resolveFieldReference(ref, ctx)
    
    require.NoError(t, err)
    assert.Equal(t, 45.5, value)
}

func TestExecutor_ResolveInlineFact(t *testing.T) {
    // Créer un fait sensor pour le contexte
    sensorFact := &Fact{
        ID:   "S001",
        Type: "Sensor",
        Fields: map[string]interface{}{
            "sensorId":    "S001",
            "temperature": 45.5,
        },
    }
    
    ctx := &ExecutionContext{
        Bindings: map[string]*Fact{
            "s": sensorFact,
        },
    }
    
    // Créer un fait inline avec références
    inlineFact := &InlineFactArgument{
        TypeName: "Alert",
        Fields: map[string]ActionArgument{
            "level": &LiteralArgument{Value: "HIGH"},
            "sensorId": &FieldReferenceArgument{
                Variable:  "s",
                FieldPath: []string{"sensorId"},
            },
            "temperature": &FieldReferenceArgument{
                Variable:  "s",
                FieldPath: []string{"temperature"},
            },
        },
    }
    
    // Résoudre
    executor := NewActionExecutor()
    fact, err := executor.resolveInlineFact(inlineFact, ctx)
    
    require.NoError(t, err)
    assert.Equal(t, "Alert", fact.Type)
    assert.Equal(t, "HIGH", fact.Fields["level"])
    assert.Equal(t, "S001", fact.Fields["sensorId"])
    assert.Equal(t, 45.5, fact.Fields["temperature"])
}
```

### Test 7: E2E - Action Xuple avec Fait Inline

**Fichier**: `rete/integration_inline_test.go`

```go
func TestE2E_XupleActionWithInlineFact(t *testing.T) {
    // Programme TSD complet
    program := `
        type Sensor(id: string, temp: number)
        type Alert(level: string, sensorId: string, temperature: number)
        
        rule high_temp: {s: Sensor} / s.temp > 40.0 ==>
            Xuple("alerts", Alert(
                level: "HIGH",
                sensorId: s.id,
                temperature: s.temp
            ))
        
        Sensor(id: "S001", temp: 25.0)
        Sensor(id: "S002", temp: 45.0)
    `
    
    // Parser et exécuter
    ast, err := constraint.ParseTSD(program)
    require.NoError(t, err)
    
    network := NewReteNetwork(NewMemoryStorage())
    // Configurer XupleManager et handler...
    
    pipeline := NewConstraintPipeline()
    result, _, err := pipeline.IngestFromAST(ast, network)
    require.NoError(t, err)
    
    // Vérifier que l'action Xuple a créé un xuple
    xupleManager := result.GetXupleManager()
    space := xupleManager.GetXupleSpace("alerts")
    xuples := space.ListAll()
    
    require.Len(t, xuples, 1)
    alert := xuples[0].Fact
    assert.Equal(t, "Alert", alert.Type)
    assert.Equal(t, "HIGH", alert.Fields["level"])
    assert.Equal(t, "S002", alert.Fields["sensorId"])
    assert.Equal(t, 45.0, alert.Fields["temperature"])
}
```

---

## ✅ Checklist de Validation

Avant de considérer cette tâche terminée, vérifier:

### Parser
- [ ] La grammaire accepte les faits inline dans les actions
- [ ] Les syntaxes simple et multi-ligne sont supportées
- [ ] Les références aux champs (`var.field`) sont reconnues
- [ ] Les actions multiples séparées par virgules fonctionnent
- [ ] Les messages d'erreur sont clairs en cas de syntaxe invalide

### AST
- [ ] Les nouveaux nœuds (`InlineFactNode`, `FieldReferenceNode`) sont définis
- [ ] L'interface `ExpressionNode` est correctement implémentée
- [ ] Les nœuds contiennent toutes les informations nécessaires
- [ ] Les locations (source positions) sont correctes pour les erreurs

### Conversion
- [ ] Les faits inline sont convertis en `InlineFactArgument`
- [ ] Les références de champs sont converties en `FieldReferenceArgument`
- [ ] Le binding context est correctement utilisé
- [ ] La validation des types est effectuée

### Runtime
- [ ] Les références sont résolues correctement lors de l'exécution
- [ ] Les faits inline sont créés avec les bonnes valeurs
- [ ] Le contexte d'exécution contient les faits déclencheurs
- [ ] Les erreurs de résolution sont bien gérées

### Tests
- [ ] Tous les tests passent: `go test ./constraint/... ./rete/...`
- [ ] Couverture de code > 80% pour le nouveau code
- [ ] Tests unitaires pour chaque composant
- [ ] Test d'intégration E2E complet
- [ ] Tests de cas d'erreur (syntaxe invalide, variable non définie, etc.)

### Standards
- [ ] En-tête copyright sur les nouveaux fichiers
- [ ] GoDoc complet pour les exports
- [ ] Pas de hardcoding (valeurs, chemins, etc.)
- [ ] `go fmt` et `goimports` appliqués
- [ ] `go vet`, `staticcheck`, `errcheck` sans erreur
- [ ] `make validate` passe

---

## 📝 Documentation à Mettre à Jour

Après implémentation, mettre à jour:

1. **Documentation syntaxe TSD** (`docs/syntax.md` ou équivalent):
   - Ajouter section sur les faits inline dans les actions
   - Documenter la syntaxe des références aux variables
   - Fournir des exemples complets

2. **Documentation développeur**:
   - Expliquer l'architecture AST étendue
   - Documenter le processus de résolution des références
   - Ajouter des exemples de code

3. **CHANGELOG.md**:
   - Ajouter entrée pour cette nouvelle fonctionnalité

---

## 🎯 Résultat Attendu

À la fin de cette tâche, le code suivant doit fonctionner parfaitement:

```tsd
type Sensor(id: string, location: string, temperature: number, humidity: number)
type Alert(level: string, message: string, sensorId: string, temperature: number)
type Command(action: string, target: string, priority: number, reason: string)

xuple-space critical_alerts {
    selection: lifo
    consumption: per-agent
    retention: unlimited
}

rule critical_temperature: {s: Sensor} / s.temperature > 40.0 ==>
    Print("CRITICAL: Sensor ", s.id, " at ", s.temperature, "°C"),
    Xuple("critical_alerts", Alert(
        level: "CRITICAL",
        message: "Temperature exceeds safe threshold",
        sensorId: s.id,
        temperature: s.temperature
    ))

rule critical_combined: {s: Sensor} / s.temperature > 40.0 AND s.humidity > 80.0 ==>
    Xuple("critical_alerts", Alert(
        level: "EMERGENCY",
        message: "Critical conditions detected",
        sensorId: s.id,
        temperature: s.temperature
    )),
    Xuple("commands", Command(
        action: "shutdown",
        target: s.location,
        priority: 999,
        reason: "Emergency conditions"
    ))

Sensor(id: "S001", location: "ServerRoom", temperature: 42.0, humidity: 85.0)
```

**Comportement attendu**:
1. Le fichier est parsé sans erreur
2. Les règles sont construites correctement
3. Lors du déclenchement de `critical_temperature`, un xuple `Alert` est créé avec `sensorId: "S001"` et `temperature: 42.0`
4. Lors du déclenchement de `critical_combined`, deux xuples sont créés (Alert + Command)
5. Les valeurs sont correctement extraites du fait `Sensor` déclencheur

---

**Prochaine étape**: Après validation de ce prompt, passer au **Prompt 02 - Package API Pipeline**