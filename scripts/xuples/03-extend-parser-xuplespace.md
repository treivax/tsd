# Prompt 03 - Extension du parser pour la commande xuple-space

## 🎯 Objectif

Étendre le langage TSD et son parser pour supporter la nouvelle commande `xuple-space` qui permet de déclarer des xuple-spaces avec leurs politiques.

Cette commande doit permettre de définir :
- Le nom du xuple-space
- La politique de sélection (random, FIFO, LIFO)
- La politique de consommation (once, per-agent, limited)
- La politique de rétention (unlimited, duration)

## 📋 Tâches

### 1. Analyser la grammaire actuelle et les commandes existantes

**Objectif** : Comprendre comment ajouter une nouvelle commande au langage TSD.

- [ ] Examiner la grammaire PEG existante (fichier `.peg`)
- [ ] Analyser comment sont définies les commandes existantes (fact, rule, action, etc.)
- [ ] Comprendre la structure AST générée
- [ ] Identifier le pattern à suivre pour ajouter `xuple-space`

**Livrables** :
- Créer `tsd/docs/xuples/implementation/01-parser-analysis.md` documentant :
  - Localisation de la grammaire
  - Structure actuelle des commandes
  - Pattern de parsing à suivre
  - Exemples de commandes similaires

### 2. Concevoir la syntaxe de la commande xuple-space

**Objectif** : Définir la syntaxe exacte de la nouvelle commande.

- [ ] Concevoir une syntaxe claire et cohérente avec le langage TSD
- [ ] Définir les mots-clés pour les politiques
- [ ] Prévoir les paramètres optionnels et obligatoires
- [ ] Assurer la lisibilité et l'expressivité

**Syntaxe proposée** :
```tsd
xuple-space <name> {
    selection: <random|fifo|lifo>
    consumption: <once|per-agent|limited(<count>)>
    retention: <unlimited|duration(<time>)>
}
```

**Exemples d'utilisation** :
```tsd
xuple-space agents-commands {
    selection: fifo
    consumption: once
    retention: unlimited
}

xuple-space notifications {
    selection: random
    consumption: per-agent
    retention: duration(5m)
}

xuple-space shared-data {
    selection: lifo
    consumption: limited(3)
    retention: duration(1h)
}
```

**Livrables** :
- Créer `tsd/docs/xuples/implementation/02-xuplespace-syntax.md` contenant :
  - Syntaxe complète avec BNF/EBNF
  - Tous les cas d'utilisation
  - Valeurs par défaut si paramètres omis
  - Exemples valides et invalides
  - Messages d'erreur attendus

### 3. Étendre la grammaire PEG

**Objectif** : Ajouter la règle `xuple-space` à la grammaire.

- [ ] Ajouter la règle `XupleSpaceDeclaration` dans le fichier PEG
- [ ] Définir les sous-règles pour les politiques
- [ ] Ajouter les mots-clés nécessaires
- [ ] Respecter les conventions de la grammaire existante
- [ ] Valider la syntaxe de la grammaire

**Fichiers à modifier** :
- `tsd/parser/grammar.peg` (ou équivalent)

**Règle PEG attendue (exemple)** :
```peg
XupleSpaceDeclaration <- "xuple-space" _ Identifier _ '{' _ XupleSpaceBody _ '}' _

XupleSpaceBody <- (XupleSpaceProperty _)*

XupleSpaceProperty <- SelectionPolicy / ConsumptionPolicy / RetentionPolicy

SelectionPolicy <- "selection:" _ SelectionValue

SelectionValue <- "random" / "fifo" / "lifo"

ConsumptionPolicy <- "consumption:" _ ConsumptionValue

ConsumptionValue <- "once" / "per-agent" / ("limited" _ '(' _ Integer _ ')')

RetentionPolicy <- "retention:" _ RetentionValue

RetentionValue <- "unlimited" / ("duration" _ '(' _ Duration _ ')')

Duration <- Integer TimeUnit

TimeUnit <- "s" / "m" / "h" / "d"
```

**Livrables** :
- [ ] Grammaire étendue avec la nouvelle commande
- [ ] Fichier PEG modifié et validé
- [ ] Documentation des modifications dans le commit

### 4. Définir les structures AST pour xuple-space

**Objectif** : Créer les types Go représentant la commande parsée.

- [ ] Créer la structure `XupleSpaceDeclaration`
- [ ] Créer les structures pour chaque type de politique
- [ ] Implémenter les interfaces AST nécessaires
- [ ] Ajouter les méthodes de validation
- [ ] Respecter les conventions de nommage du projet

**Structures attendues** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package ast

// XupleSpaceDeclaration représente une déclaration de xuple-space
type XupleSpaceDeclaration struct {
    Name              string
    SelectionPolicy   SelectionPolicyType
    ConsumptionPolicy ConsumptionPolicyConfig
    RetentionPolicy   RetentionPolicyConfig
    Location          Location // Position dans le fichier source
}

// SelectionPolicyType représente le type de politique de sélection
type SelectionPolicyType int

const (
    SelectionRandom SelectionPolicyType = iota
    SelectionFIFO
    SelectionLIFO
)

// ConsumptionPolicyConfig configure la politique de consommation
type ConsumptionPolicyConfig struct {
    Type  ConsumptionPolicyType
    Limit int // Pour limited(n), sinon 0
}

type ConsumptionPolicyType int

const (
    ConsumptionOnce ConsumptionPolicyType = iota
    ConsumptionPerAgent
    ConsumptionLimited
)

// RetentionPolicyConfig configure la politique de rétention
type RetentionPolicyConfig struct {
    Type     RetentionPolicyType
    Duration time.Duration // Pour duration(x), sinon 0
}

type RetentionPolicyType int

const (
    RetentionUnlimited RetentionPolicyType = iota
    RetentionDuration
)
```

**Fichiers à créer/modifier** :
- `tsd/parser/ast/xuplespace.go` (nouveau)
- `tsd/parser/ast/ast.go` (modification si nécessaire)

**Livrables** :
- [ ] Structures AST complètes avec copyright
- [ ] Méthodes de validation
- [ ] Méthodes String() pour debug
- [ ] Tests unitaires des structures

### 5. Implémenter la transformation PEG → AST

**Objectif** : Convertir le résultat du parsing PEG en structures AST Go.

- [ ] Implémenter les fonctions de transformation
- [ ] Gérer tous les cas de politiques
- [ ] Valider les valeurs parsées (durées, limites, etc.)
- [ ] Gérer les erreurs de parsing avec messages clairs
- [ ] Conserver les informations de localisation (ligne, colonne)

**Fichiers à modifier** :
- `tsd/parser/parser.go` (ou fichier de transformation PEG)

**Exemple de code attendu** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

func (p *Parser) buildXupleSpaceDeclaration(node *peg.Node) (*ast.XupleSpaceDeclaration, error) {
    // Extraction du nom
    name := extractIdentifier(node)
    
    // Parsing des politiques avec valeurs par défaut
    selection := ast.SelectionFIFO // Défaut
    consumption := ast.ConsumptionPolicyConfig{Type: ast.ConsumptionOnce}
    retention := ast.RetentionPolicyConfig{Type: ast.RetentionUnlimited}
    
    // Parser chaque propriété
    for _, prop := range node.Children {
        switch prop.Type {
        case "SelectionPolicy":
            selection = p.parseSelectionPolicy(prop)
        case "ConsumptionPolicy":
            consumption = p.parseConsumptionPolicy(prop)
        case "RetentionPolicy":
            retention = p.parseRetentionPolicy(prop)
        }
    }
    
    return &ast.XupleSpaceDeclaration{
        Name:              name,
        SelectionPolicy:   selection,
        ConsumptionPolicy: consumption,
        RetentionPolicy:   retention,
        Location:          extractLocation(node),
    }, nil
}
```

**Livrables** :
- [ ] Fonctions de transformation complètes
- [ ] Gestion d'erreurs robuste
- [ ] Messages d'erreur clairs et localisés
- [ ] Validation des valeurs (durées positives, limites > 0, etc.)

### 6. Intégrer xuple-space dans le processus de compilation

**Objectif** : Faire en sorte que les déclarations xuple-space soient traitées lors de la compilation.

- [ ] Ajouter `XupleSpaceDeclaration` à la liste des déclarations possibles
- [ ] Créer un registre de xuple-spaces dans le compilateur
- [ ] Valider l'unicité des noms de xuple-spaces
- [ ] Détecter les doublons et générer des erreurs claires
- [ ] Intégrer dans le flux de compilation existant

**Fichiers à modifier** :
- `tsd/compiler/compiler.go` (ou équivalent)
- `tsd/compiler/context.go` (pour le registre)

**Structures attendues** :
```go
// Dans le contexte de compilation
type CompilationContext struct {
    // ... champs existants ...
    XupleSpaces map[string]*ast.XupleSpaceDeclaration
}

// Validation lors de l'ajout
func (ctx *CompilationContext) RegisterXupleSpace(decl *ast.XupleSpaceDeclaration) error {
    if _, exists := ctx.XupleSpaces[decl.Name]; exists {
        return fmt.Errorf("xuple-space '%s' already declared at line %d", 
            decl.Name, decl.Location.Line)
    }
    ctx.XupleSpaces[decl.Name] = decl
    return nil
}
```

**Livrables** :
- [ ] Registre de xuple-spaces dans le contexte
- [ ] Validation des doublons
- [ ] Messages d'erreur cohérents avec le reste du compilateur
- [ ] Intégration dans le flux de compilation

### 7. Créer les tests du parser

**Objectif** : Tester exhaustivement le parsing de xuple-space.

- [ ] Tests de parsing valide (tous les cas)
- [ ] Tests de parsing invalide (erreurs attendues)
- [ ] Tests de détection de doublons
- [ ] Tests des valeurs par défaut
- [ ] Tests des cas limites (durées, limites)

**Fichier à créer** :
- `tsd/parser/xuplespace_test.go`

**Tests attendus (exemples)** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package parser

import "testing"

func TestParseXupleSpace_Valid(t *testing.T) {
    t.Log("🧪 TEST PARSING XUPLE-SPACE VALID")
    
    tests := []struct {
        name     string
        input    string
        expected ast.XupleSpaceDeclaration
    }{
        {
            name: "xuple-space complet",
            input: `xuple-space myspace {
                selection: fifo
                consumption: once
                retention: unlimited
            }`,
            expected: ast.XupleSpaceDeclaration{
                Name:              "myspace",
                SelectionPolicy:   ast.SelectionFIFO,
                ConsumptionPolicy: ast.ConsumptionPolicyConfig{Type: ast.ConsumptionOnce},
                RetentionPolicy:   ast.RetentionPolicyConfig{Type: ast.RetentionUnlimited},
            },
        },
        {
            name: "avec duration",
            input: `xuple-space timed {
                selection: random
                consumption: per-agent
                retention: duration(5m)
            }`,
            expected: ast.XupleSpaceDeclaration{
                Name:              "timed",
                SelectionPolicy:   ast.SelectionRandom,
                ConsumptionPolicy: ast.ConsumptionPolicyConfig{Type: ast.ConsumptionPerAgent},
                RetentionPolicy:   ast.RetentionPolicyConfig{
                    Type:     ast.RetentionDuration,
                    Duration: 5 * time.Minute,
                },
            },
        },
        {
            name: "avec limited",
            input: `xuple-space limited {
                selection: lifo
                consumption: limited(3)
                retention: unlimited
            }`,
            expected: ast.XupleSpaceDeclaration{
                Name:              "limited",
                SelectionPolicy:   ast.SelectionLIFO,
                ConsumptionPolicy: ast.ConsumptionPolicyConfig{
                    Type:  ast.ConsumptionLimited,
                    Limit: 3,
                },
                RetentionPolicy: ast.RetentionPolicyConfig{Type: ast.RetentionUnlimited},
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseTSD(tt.input)
            if err != nil {
                t.Fatalf("❌ Erreur parsing: %v", err)
            }
            
            // Validation du résultat
            if len(result.XupleSpaces) != 1 {
                t.Fatalf("❌ Attendu 1 xuple-space, reçu %d", len(result.XupleSpaces))
            }
            
            xs := result.XupleSpaces[0]
            if xs.Name != tt.expected.Name {
                t.Errorf("❌ Nom: attendu '%s', reçu '%s'", tt.expected.Name, xs.Name)
            }
            
            // ... autres assertions ...
            
            t.Log("✅ Test réussi")
        })
    }
}

func TestParseXupleSpace_Invalid(t *testing.T) {
    t.Log("🧪 TEST PARSING XUPLE-SPACE INVALID")
    
    tests := []struct {
        name        string
        input       string
        expectedErr string
    }{
        {
            name: "politique de sélection invalide",
            input: `xuple-space bad {
                selection: invalid
            }`,
            expectedErr: "invalid selection policy",
        },
        {
            name: "durée négative",
            input: `xuple-space bad {
                retention: duration(-5m)
            }`,
            expectedErr: "duration must be positive",
        },
        {
            name: "limite zéro",
            input: `xuple-space bad {
                consumption: limited(0)
            }`,
            expectedErr: "limit must be greater than zero",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := ParseTSD(tt.input)
            if err == nil {
                t.Fatal("❌ Erreur attendue mais parsing réussi")
            }
            
            if !strings.Contains(err.Error(), tt.expectedErr) {
                t.Errorf("❌ Erreur attendue '%s', reçu '%s'", tt.expectedErr, err.Error())
            }
            
            t.Log("✅ Erreur correctement détectée")
        })
    }
}

func TestParseXupleSpace_Duplicates(t *testing.T) {
    t.Log("🧪 TEST DÉTECTION DOUBLONS XUPLE-SPACE")
    
    input := `
        xuple-space myspace {
            selection: fifo
        }
        
        xuple-space myspace {
            selection: lifo
        }
    `
    
    _, err := ParseTSD(input)
    if err == nil {
        t.Fatal("❌ Erreur de doublon attendue")
    }
    
    if !strings.Contains(err.Error(), "already declared") {
        t.Errorf("❌ Message d'erreur incorrect: %v", err)
    }
    
    t.Log("✅ Doublon correctement détecté")
}
```

**Livrables** :
- [ ] Tests complets avec couverture > 80%
- [ ] Tests de tous les cas valides
- [ ] Tests de tous les cas d'erreur
- [ ] Tests de cas limites
- [ ] Messages de test clairs avec émojis

### 8. Créer des exemples et documentation

**Objectif** : Documenter la nouvelle fonctionnalité pour les utilisateurs.

- [ ] Créer des exemples TSD complets
- [ ] Documenter la syntaxe dans le guide utilisateur
- [ ] Créer un guide de référence des politiques
- [ ] Ajouter des exemples de cas d'usage

**Fichiers à créer** :
- `tsd/docs/xuples/user-guide/xuplespace-command.md`
- `tsd/examples/xuples/basic-xuplespace.tsd`
- `tsd/examples/xuples/all-policies.tsd`

**Exemple de documentation utilisateur** :
```markdown
# Commande xuple-space

La commande `xuple-space` permet de déclarer un espace de xuples avec ses politiques.

## Syntaxe

\```tsd
xuple-space <nom> {
    selection: <random|fifo|lifo>
    consumption: <once|per-agent|limited(n)>
    retention: <unlimited|duration(temps)>
}
\```

## Politiques

### Selection
- `random` : Sélection aléatoire
- `fifo` : Premier entré, premier sorti
- `lifo` : Dernier entré, premier sorti

### Consumption
- `once` : Un seul consommateur au total
- `per-agent` : Une fois par agent
- `limited(n)` : Maximum n consommations

### Retention
- `unlimited` : Pas d'expiration
- `duration(temps)` : Expire après la durée (ex: 5m, 1h, 2d)

## Exemples

Voir `examples/xuples/` pour des exemples complets.
```

**Livrables** :
- [ ] Documentation utilisateur complète
- [ ] Exemples TSD fonctionnels
- [ ] Guide de référence des politiques

## 📁 Structure attendue

```
tsd/
├── docs/xuples/
│   ├── implementation/
│   │   ├── 01-parser-analysis.md
│   │   └── 02-xuplespace-syntax.md
│   └── user-guide/
│       └── xuplespace-command.md
├── examples/xuples/
│   ├── basic-xuplespace.tsd
│   └── all-policies.tsd
├── parser/
│   ├── grammar.peg                  # Modifié
│   ├── ast/
│   │   └── xuplespace.go            # Nouveau
│   ├── parser.go                    # Modifié
│   └── xuplespace_test.go           # Nouveau
└── compiler/
    ├── compiler.go                  # Modifié
    └── context.go                   # Modifié
```

## ✅ Critères de succès

- [ ] Grammaire PEG étendue et validée
- [ ] Structures AST complètes avec copyright
- [ ] Parsing fonctionnel pour tous les cas
- [ ] Validation des doublons implémentée
- [ ] Tests complets avec couverture > 80%
- [ ] Tous les tests passent
- [ ] Aucun hardcoding
- [ ] Messages d'erreur clairs et localisés
- [ ] Documentation utilisateur complète
- [ ] Exemples fonctionnels fournis
- [ ] `make validate` passe sans erreur

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `tsd/docs/xuples/design/` - Conception du module (prompt 02)
- Grammaire PEG existante
- Documentation parser existante
- Effective Go - https://go.dev/doc/effective_go

## 🎯 Prochaine étape

Une fois le parsing de `xuple-space` terminé et testé, passer au prompt **04-implement-default-actions.md** pour implémenter le système d'actions par défaut.