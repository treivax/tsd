# Instruction `reset` - Documentation

## Vue d'ensemble

L'instruction `reset` permet de réinitialiser complètement le système TSD, supprimant tous les types, règles, faits et le réseau RETE en une seule commande.

## Syntaxe

```
reset
```

L'instruction `reset` est un mot-clé simple qui ne prend aucun paramètre.

## Comportement

Lorsqu'une instruction `reset` est exécutée, elle effectue les actions suivantes :

1. **Supprime tous les types** définis précédemment
2. **Supprime toutes les règles** (expressions)
3. **Supprime tous les faits** présents dans le système
4. **Réinitialise le réseau RETE** complètement
5. **Nettoie l'état du programme** (ProgramState)

Après un `reset`, le système est dans un état vide, comme s'il venait d'être démarré. Vous pouvez alors définir de nouveaux types, règles et faits.

## Cas d'utilisation

### 1. Redémarrage complet du système

Lorsque vous voulez recommencer avec une configuration entièrement nouvelle sans redémarrer l'application :

```
type User : <name: string, age: number>
type Order : <id: number, amount: number>

# ... définitions et traitement ...

reset

# Nouveau système avec types différents
type Product : <id: string, price: number>
type Category : <name: string>
```

### 2. Tests et développement

Utile pour réinitialiser l'état entre différents scénarios de test :

```
# Scénario 1
type TestType1 : <field1: string>
# ... tests ...

reset

# Scénario 2 (complètement isolé)
type TestType2 : <field2: number>
# ... tests ...
```

### 3. Changement de contexte

Lorsque vous passez d'un contexte métier à un autre :

```
# Contexte e-commerce
type Customer : <id: string, name: string>
type Product : <id: string, price: number>
# ... règles e-commerce ...

reset

# Contexte RH
type Employee : <id: string, name: string, salary: number>
type Department : <name: string, budget: number>
# ... règles RH ...
```

## Exemples complets

### Exemple 1 : Reset simple

```
type User : <name: string, age: number>

reset

type Order : <id: number>
```

**Résultat** : Le type `User` est supprimé, seul le type `Order` existe après le reset.

### Exemple 2 : Multiples resets

```
type TypeA : <field: string>

reset

type TypeB : <field: number>

reset

type TypeC : <field: boolean>
```

**Résultat** : À la fin, seul le type `TypeC` existe. Les types `TypeA` et `TypeB` ont été supprimés par les resets successifs.

### Exemple 3 : Reset avec fichiers séparés

**fichier1.constraint :**
```
type User : <name: string>
type Order : <id: number>
```

**fichier2.constraint :**
```
reset

type Product : <name: string, price: number>
```

Si vous parsez ces fichiers dans l'ordre avec un `IterativeParser`, après le parsing de `fichier2.constraint`, seul le type `Product` existera.

## API de programmation

### Constraint Package

#### ProgramState.Reset()

Réinitialise l'état du programme (types, règles, faits) :

```go
ps := constraint.NewProgramState()

// Ajouter des données
ps.Types["User"] = &constraint.TypeDefinition{Name: "User"}

// Reset
ps.Reset()

// ps.Types est maintenant vide
```

#### IterativeParser.Reset()

Réinitialise le parser itératif :

```go
parser := constraint.NewIterativeParser()

// Parser du contenu
parser.ParseContent(content1, "file1.constraint")

// Reset
parser.Reset()

// L'état est vide, prêt pour nouveau parsing
parser.ParseContent(content2, "file2.constraint")
```

### RETE Package

#### ReteNetwork.Reset()

Réinitialise le réseau RETE :

```go
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)

// Construire le réseau
// ... ajouter types, rules, facts ...

// Reset
network.Reset()

// Le réseau est vide, avec un nouveau RootNode
// Le Storage est préservé
```

## Détails d'implémentation

### Grammaire PEG

L'instruction `reset` a été ajoutée à la grammaire comme une nouvelle règle :

```peg
Statement <- TypeDefinition / Expression / RemoveFact / Fact / Reset

Reset <- "reset" {
    return map[string]interface{}{
        "type": "reset",
    }, nil
}
```

### Structure de données

```go
type Reset struct {
    Type string `json:"type"` // Toujours "reset"
}

type Program struct {
    Types       []TypeDefinition `json:"types"`
    Expressions []Expression     `json:"expressions"`
    Facts       []Fact           `json:"facts"`
    Resets      []Reset          `json:"resets"`
}
```

### Comportement du reset

1. **ProgramState.Reset()** :
   - Vide la map `Types`
   - Vide les slices `Rules`, `Facts`, `FilesParsed`
   - Vide les erreurs de validation
   - Réinitialise avec des structures vides (non nil)

2. **ReteNetwork.Reset()** :
   - Vide toutes les maps de nœuds (`TypeNodes`, `AlphaNodes`, `BetaNodes`, `TerminalNodes`)
   - Vide le slice `Types`
   - Recrée un nouveau `RootNode` avec le même `Storage`
   - Réinitialise `BetaBuilder` à `nil`

## Considérations importantes

### ⚠️ Perte de données

**ATTENTION** : L'instruction `reset` est **destructive** et **irréversible**. Toutes les données (types, règles, faits) sont définitivement supprimées.

- ✅ Utilisez `reset` seulement quand vous êtes certain de vouloir tout effacer
- ✅ Sauvegardez vos définitions importantes avant un reset
- ✅ Documentez l'utilisation de reset dans vos fichiers

### Performance

L'instruction `reset` est **très rapide** car elle :
- Ne parcourt pas les données existantes
- Crée simplement de nouvelles structures vides
- Préserve la référence au Storage (pas de réallocation)

### Ordre de parsing

Si vous utilisez des fichiers multiples avec un `IterativeParser`, l'ordre compte :

```go
parser := constraint.NewIterativeParser()
parser.ParseFile("types.constraint")    // Types A, B, C définis
parser.ParseFile("reset.constraint")     // Contient "reset"
parser.ParseFile("newtypes.constraint") // Types X, Y, Z définis

// Résultat : seuls X, Y, Z existent
```

## Tests

Des tests complets ont été créés pour valider le comportement :

- **constraint/reset_test.go** : Tests du parsing et du ProgramState
- **rete/reset_test.go** : Tests du réseau RETE
- **beta_coverage_tests/reset_example.constraint** : Exemple d'utilisation

Pour exécuter les tests :

```bash
# Tests de l'instruction reset
go test ./constraint -run TestResetInstruction
go test ./constraint -run TestProgramStateReset
go test ./constraint -run TestIterativeParserReset

# Tests du réseau RETE
go test ./rete -run TestReteNetworkReset

# Tous les tests
make test
```

## Changelog

**Version** : 2.3.1  
**Date** : 26 novembre 2025  
**Type** : Nouvelle fonctionnalité

### Modifications apportées

1. **Grammaire** : Ajout de l'instruction `reset` dans `constraint/grammar/constraint.peg`
2. **Types** : Ajout de la structure `Reset` dans `constraint/constraint_types.go`
3. **API Constraint** : Méthode `Reset()` ajoutée à `ProgramState` et `IterativeParser`
4. **API RETE** : Méthode `Reset()` ajoutée à `ReteNetwork`
5. **Tests** : Suite de tests complète pour toutes les composantes
6. **Documentation** : Ce document

## Références

- [Grammaire PEG](../constraint/grammar/constraint.peg)
- [Types de contraintes](../constraint/constraint_types.go)
- [API Constraint](../constraint/api.go)
- [Réseau RETE](../rete/network.go)
- [Exemples](../beta_coverage_tests/reset_example.constraint)

## Support

Pour toute question ou problème concernant l'instruction `reset`, consultez :

1. Les tests dans `constraint/reset_test.go` et `rete/reset_test.go`
2. L'exemple dans `beta_coverage_tests/reset_example.constraint`
3. Cette documentation

---

**Auteur** : TSD Contributors  
**Licence** : MIT  
**Dernière mise à jour** : 26 novembre 2025