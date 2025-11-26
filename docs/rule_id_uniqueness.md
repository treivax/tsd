# Validation de l'Unicité des Identifiants de Règles

## 📋 Vue d'ensemble

Le système TSD valide automatiquement l'unicité des identifiants de règles lors du parsing. Cette validation garantit qu'aucune règle ne peut avoir le même identifiant qu'une règle précédemment parsée, sauf après un `reset`.

## 🎯 Comportement

### Règles de base

1. **Unicité obligatoire** : Chaque identifiant de règle doit être unique dans le contexte courant
2. **Erreur non-bloquante** : Les règles avec ID dupliqué sont **ignorées** mais ne font pas échouer le parsing
3. **Suivi permanent** : Les IDs utilisés sont mémorisés jusqu'à un `reset`
4. **Reset autorise la réutilisation** : Après un `reset`, tous les IDs peuvent être réutilisés

## ✅ Comportement Valide

### Exemple 1 : Règles avec IDs uniques

```
type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)
rule r3 : {p: Person} / p.age == 18 ==> exactly_eighteen(p.id)
```

**Résultat :** ✅ 3 règles acceptées (r1, r2, r3)

### Exemple 2 : Réutilisation après reset

**Fichier 1 (before_reset.constraint) :**
```
type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)
```

**Fichier 2 (after_reset.constraint) :**
```
reset

type Product : <id: string, price: number>

rule r1 : {prod: Product} / prod.price > 100 ==> expensive(prod.id)
rule r2 : {prod: Product} / prod.price < 50 ==> cheap(prod.id)
```

**Résultat :**
- ✅ Fichier 1 : 2 règles (r1, r2) pour Person
- ✅ Fichier 2 : 2 règles (r1, r2) pour Product
- ✅ Aucune erreur : le reset permet de réutiliser les IDs

## ❌ Comportement Invalide (Non-bloquant)

### Exemple 1 : ID dupliqué dans le même fichier

```
type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)
rule r1 : {p: Person} / p.age == 18 ==> exactly_eighteen(p.id)  // ⚠️ DUPLIQUÉ
```

**Résultat :**
- ✅ 2 règles acceptées : r1 (première occurrence), r2
- ⚠️ 1 règle ignorée : r1 (deuxième occurrence, ligne 6)
- ⚠️ Avertissement affiché :
  ```
  ⚠️  Skipping duplicate rule ID in file.constraint: rule ID 'r1' already used, ignoring duplicate rule
  ```
- ⚠️ Erreur enregistrée dans `ProgramState.Errors`

### Exemple 2 : ID dupliqué entre fichiers

**Fichier 1 :**
```
type Person : <id: string, age: number>
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
```

**Fichier 2 :**
```
rule r1 : {p: Person} / p.age < 18 ==> minor(p.id)  // ⚠️ DUPLIQUÉ
rule r2 : {p: Person} / p.age == 18 ==> exactly_eighteen(p.id)
```

**Résultat :**
- ✅ Fichier 1 : 1 règle acceptée (r1)
- ✅ Fichier 2 : 1 règle acceptée (r2)
- ⚠️ Fichier 2 : 1 règle ignorée (r1 dupliqué)
- ⚠️ Total : 2 règles dans le système (r1 du fichier 1, r2 du fichier 2)

## 🔍 Détection et Traçabilité

### Message d'avertissement

Lors de la détection d'un ID dupliqué, le système affiche :

```
⚠️  Skipping duplicate rule ID in <fichier>: rule ID '<id>' already used, ignoring duplicate rule
```

### Erreur enregistrée

Une entrée est ajoutée dans `ProgramState.Errors` :

```go
ValidationError{
    File:    "fichier.constraint",
    Type:    "rule",
    Message: "rule ID 'r1' already used, ignoring duplicate rule",
    Line:    0,
}
```

### Accès aux erreurs

```go
ps := constraint.NewProgramState()
err := ps.ParseAndMerge("rules.constraint")

// Vérifier les erreurs non-bloquantes
if len(ps.Errors) > 0 {
    fmt.Printf("Avertissements : %d\n", len(ps.Errors))
    for _, e := range ps.Errors {
        fmt.Printf("  - %s: %s\n", e.File, e.Message)
    }
}

// Le parsing réussit malgré les duplicates
fmt.Printf("Règles acceptées : %d\n", len(ps.Rules))
```

## 🔄 Comportement du Reset

### Effacement complet

La commande `reset` efface :
- ✅ Tous les types (`Types`)
- ✅ Toutes les règles (`Rules`)
- ✅ Tous les faits (`Facts`)
- ✅ **Tous les IDs de règles** (`RuleIDs`)
- ✅ Toutes les erreurs (`Errors`)

### Réutilisation autorisée

Après un `reset`, **tous les identifiants** peuvent être réutilisés :

```
// Contexte 1
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)

reset

// Contexte 2 : r1 peut être réutilisé
rule r1 : {prod: Product} / prod.price > 100 ==> expensive(prod.id)
```

### Reset dans un fichier unique

**Important :** Si un reset est présent dans un fichier, il efface l'état **avant** de parser le reste du fichier :

```
type Person : <id: string>
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)

reset

type Product : <id: string>
rule r1 : {prod: Product} / prod.price > 100 ==> expensive(prod.id)  // ✅ VALIDE
```

**Résultat :** Seule la règle Product existe après le parsing (Person et sa règle ont été effacées par reset).

## 💡 Cas Particuliers

### IDs vides

Les règles avec identifiant vide (`""`) sont **toujours acceptées** et ne déclenchent pas la validation d'unicité :

```go
rule1 := Expression{RuleId: ""}  // ✅ Accepté
rule2 := Expression{RuleId: ""}  // ✅ Accepté aussi (pas de vérification)
```

**Note :** Depuis la v2.0, tous les identifiants doivent être non-vides selon la grammaire PEG. Les IDs vides ne peuvent exister que dans du code programmatique.

### Multiple duplicates

Si plusieurs règles ont le même ID, **seule la première** est acceptée :

```
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)      // ✅ ACCEPTÉE
rule r1 : {p: Person} / p.age < 18 ==> minor(p.id)      // ⚠️ IGNORÉE
rule r1 : {p: Person} / p.age == 18 ==> eighteen(p.id)  // ⚠️ IGNORÉE
```

**Résultat :** 1 règle acceptée, 2 règles ignorées, 2 erreurs enregistrées.

## 🧪 Exemples de Tests

### Test unitaire

```go
func TestRuleIdUniqueness(t *testing.T) {
    ps := constraint.NewProgramState()
    
    // Ajouter le type
    ps.Types["Person"] = &constraint.TypeDefinition{
        Name: "Person",
        Fields: []constraint.Field{
            {Name: "id", Type: "string"},
            {Name: "age", Type: "number"},
        },
    }
    
    // Parse first file with r1
    ps.ParseAndMerge("file1.constraint")
    
    // Parse second file with duplicate r1
    ps.ParseAndMerge("file2.constraint")
    
    // Verify only first r1 was accepted
    assert.Equal(t, 1, len(ps.Rules))
    assert.Equal(t, 1, len(ps.Errors))
}
```

### Test d'intégration

```go
func TestResetAllowsReuseIntegration(t *testing.T) {
    ps := constraint.NewProgramState()
    
    // Parse file with rules
    ps.ParseAndMerge("before_reset.constraint")
    assert.Equal(t, 2, len(ps.Rules))
    
    // Parse file with reset and reused IDs
    ps.ParseAndMerge("after_reset.constraint")
    
    // After reset, IDs can be reused
    assert.Equal(t, 2, len(ps.Rules))
    assert.Equal(t, 0, len(ps.Errors))  // No errors
}
```

## 📊 API de Suivi

### Structure ProgramState

```go
type ProgramState struct {
    Types       map[string]*TypeDefinition
    Rules       []*Expression
    Facts       []*Fact
    FilesParsed []string
    Errors      []ValidationError
    RuleIDs     map[string]bool  // Suivi des IDs utilisés
}
```

### Méthodes de validation

```go
// mergeRules valide l'unicité lors de la fusion
func (ps *ProgramState) mergeRules(newRules []Expression, filename string) error {
    for _, rule := range newRules {
        // Vérifier l'unicité
        if rule.RuleId != "" && ps.RuleIDs[rule.RuleId] {
            // ID dupliqué : enregistrer erreur et ignorer
            ps.Errors = append(ps.Errors, ValidationError{
                File:    filename,
                Type:    "rule",
                Message: fmt.Sprintf("rule ID '%s' already used...", rule.RuleId),
            })
            continue
        }
        
        // Marquer l'ID comme utilisé
        if rule.RuleId != "" {
            ps.RuleIDs[rule.RuleId] = true
        }
        
        ps.Rules = append(ps.Rules, &rule)
    }
    return nil
}

// Reset efface tout, y compris les IDs
func (ps *ProgramState) Reset() {
    ps.Types = make(map[string]*TypeDefinition)
    ps.Rules = make([]*Expression, 0)
    ps.Facts = make([]*Fact, 0)
    ps.RuleIDs = make(map[string]bool)  // Réinitialiser les IDs
    ps.Errors = make([]ValidationError, 0)
}
```

## 🎯 Bonnes Pratiques

### 1. Utiliser des IDs descriptifs uniques

```
✅ BON
rule validate_adult_age : {p: Person} / p.age >= 18 ==> adult(p.id)
rule validate_senior_age : {p: Person} / p.age >= 65 ==> senior(p.id)

❌ MAUVAIS
rule check : {p: Person} / p.age >= 18 ==> adult(p.id)
rule check : {p: Person} / p.age >= 65 ==> senior(p.id)  // Dupliqué !
```

### 2. Préfixer par domaine pour grands projets

```
rule person_adult_check : ...
rule person_senior_check : ...
rule order_premium_check : ...
rule order_bulk_check : ...
```

### 3. Vérifier les erreurs après parsing

```go
ps := constraint.NewProgramState()
err := ps.ParseAndMerge("rules.constraint")

if len(ps.Errors) > 0 {
    log.Warn("Parsing succeeded with warnings:")
    for _, e := range ps.Errors {
        log.Warnf("  %s: %s", e.File, e.Message)
    }
}
```

### 4. Documenter les resets

```
// ===============================
// RESET : Nouveau contexte métier
// ===============================
reset

// Les règles suivantes utilisent un nouveau domaine
type NewDomain : <...>
rule r1 : ...  // OK, reset permet réutilisation
```

## 📚 Voir aussi

- [Identifiants de Règles](./rule_identifiers.md) - Syntaxe complète
- [Grammaire PEG](../constraint/grammar/constraint.peg) - Définition formelle
- [Tests de validation](../constraint/rule_id_validation_test.go) - Tests unitaires
- [CHANGELOG v2.0.0](../CHANGELOG.md) - Notes de version

## 🔗 Références

- **Fichiers de test** :
  - `constraint/rule_id_validation_test.go` - Tests unitaires complets
  - `constraint/rule_id_integration_test.go` - Tests d'intégration
  - `constraint/test/integration/duplicate_rule_ids.constraint` - Exemple de duplicates
  - `constraint/test/integration/reset_rule_ids.constraint` - Exemple avec reset

- **Code source** :
  - `constraint/program_state.go` - Logique de validation
  - `constraint/program_state_methods.go` - Méthode Reset()

---

**Version** : 2.0.0  
**Date** : Janvier 2025  
**Statut** : ✅ Implémenté et testé