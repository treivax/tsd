# TODO - Intégration Système de Validation

**Date**: 2025-12-19  
**Contexte**: Implémentation complète du système de validation de types (TypeSystem, FactValidator, ProgramValidator)

---

## ✅ Complété

- [x] Création de `TypeSystem` avec gestion complète des types
- [x] Création de `FactValidator` pour validation des faits
- [x] Création de `ProgramValidator` pour orchestration complète
- [x] Tests unitaires complets (47 tests, 100% PASS)
- [x] Validation code (go fmt, go vet, staticcheck)
- [x] Couverture > 80% (84.8%)
- [x] Documentation complète (GoDoc + rapport)

---

## 📋 Actions à Réaliser pour Intégration Complète

### 1. Intégration dans `constraint/api.go`

**Fichier** : `constraint/api.go`

**Action** : Ajouter une fonction wrapper pour faciliter l'utilisation

```go
// ParseAndValidateProgram parse et valide un programme complet.
// Cette fonction combine le parsing et la validation en une seule étape.
func ParseAndValidateProgram(input string) (*Program, error) {
    // Parser le programme
    program, err := ParseProgram(input)
    if err != nil {
        return nil, fmt.Errorf("erreur de parsing: %v", err)
    }
    
    // Valider avec le nouveau système
    validator := NewProgramValidator()
    if err := validator.Validate(*program); err != nil {
        return nil, fmt.Errorf("erreur de validation: %v", err)
    }
    
    return program, nil
}
```

**Priorité** : 🔴 Haute (facilite l'utilisation)

**Impact** : Aucune modification du code existant - fonction additionnelle

---

### 2. Intégration dans `constraint/constraint_program.go`

**Fichier** : `constraint/constraint_program.go`

**Action** : Remplacer/compléter `ValidateProgram` pour utiliser le nouveau système

**Option 1 - Wrapper (recommandé)** :
```go
// ValidateProgram effectue une validation complète du programme parsé.
// Utilise le nouveau ProgramValidator pour une validation exhaustive.
func ValidateProgram(result interface{}) error {
    // Convertir le résultat en structure Program
    program, err := convertResultToProgram(result)
    if err != nil {
        return err
    }
    
    // Normaliser les types de valeurs de faits
    normalizeFactValueTypes(&program)
    
    // Utiliser le nouveau ProgramValidator
    validator := NewProgramValidator()
    if err := validator.Validate(program); err != nil {
        return err
    }
    
    // Validations supplémentaires existantes (xuple-spaces, etc.)
    if err := validateXupleSpaces(program); err != nil {
        return fmt.Errorf("erreur validation xuple-spaces: %v", err)
    }
    
    tsdio.Printf("✓ Programme valide avec %d type(s), %d expression(s), %d fait(s), %d affectation(s) et %d xuple-space(s)\n",
        len(program.Types), len(program.Expressions), len(program.Facts), len(program.FactAssignments), len(program.XupleSpaces))
    return nil
}
```

**Option 2 - Remplacement complet** :
- Supprimer les anciennes fonctions de validation
- Migrer toute la logique vers ProgramValidator
- Plus invasif mais plus propre à long terme

**Priorité** : 🔴 Haute (point d'entrée principal)

**Impact** : 
- Aucun si Option 1 (wrapper)
- Moyen si Option 2 (refactoring complet)

---

### 3. Mise à Jour de la Validation des Actions

**Fichier** : `constraint/action_validator.go`

**Action** : Intégrer TypeSystem dans ActionValidator

```go
type ActionValidator struct {
    actions          map[string]*ActionDefinition
    types            map[string]*TypeDefinition
    functionRegistry *FunctionRegistry
    typeSystem       *TypeSystem  // NOUVEAU
}

func NewActionValidator(actions []ActionDefinition, types []TypeDefinition) *ActionValidator {
    av := &ActionValidator{
        actions:          make(map[string]*ActionDefinition),
        types:            make(map[string]*TypeDefinition),
        functionRegistry: DefaultFunctionRegistry,
        typeSystem:       NewTypeSystem(types),  // NOUVEAU
    }
    // ... rest of initialization
    return av
}

// Utiliser typeSystem au lieu de manipuler types directement
func (av *ActionValidator) isTypeCompatible(argType, paramType string) bool {
    if paramType == "any" {
        return true
    }
    
    return av.typeSystem.AreTypesCompatible(argType, paramType, "==")
}
```

**Priorité** : 🟡 Moyenne (amélioration mais non critique)

**Impact** : Faible - amélioration de cohérence

---

### 4. Tests d'Intégration

**Fichier** : Créer `constraint/integration_validation_test.go`

**Action** : Tester l'intégration complète du système

```go
func TestProgramValidation_CompleteIntegration(t *testing.T) {
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string, password: string)
        type Audit(login: Login, timestamp: number, action: string)
        
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        
        login1 = Login(alice, "alice@ex.com", "pw1")
        login2 = Login(bob, "bob@ex.com", "pw2")
        
        Audit(login1, 1234567890, "login")
        
        {u: User, l: Login} / l.user == u && u.age > 25 ==> 
            Log("Senior user login: " + u.name)
    `
    
    program, err := ParseAndValidateProgram(input)
    if err != nil {
        t.Fatalf("Erreur: %v", err)
    }
    
    if program == nil {
        t.Fatal("Programme nil")
    }
    
    // Vérifier que tout est validé correctement
    // ...
}
```

**Priorité** : 🟡 Moyenne (validation mais non bloquant)

**Impact** : Aucun - tests additionnels

---

### 5. Migration des Tests Existants

**Fichiers** : 
- `constraint/validation_test.go`
- `constraint/comprehensive_validation_test.go`
- Autres tests utilisant validation

**Action** : Vérifier compatibilité et migrer si nécessaire

**Étapes** :
1. Exécuter tous les tests existants : `go test ./constraint -v`
2. Identifier les tests qui échouent avec le nouveau système
3. Migrer les tests pour utiliser le nouveau système
4. Supprimer les tests obsolètes/redondants

**Priorité** : 🟢 Basse (les tests existants passent déjà)

**Impact** : Aucun si pas de migration

---

### 6. Documentation Utilisateur

**Fichier** : Créer `docs/validation/README.md`

**Action** : Documenter le système de validation pour les utilisateurs

**Contenu** :
- Guide d'utilisation du système de validation
- Exemples de programmes valides/invalides
- Messages d'erreur courants et solutions
- Référence API (TypeSystem, FactValidator, ProgramValidator)

**Priorité** : 🟡 Moyenne (important pour adoption)

**Impact** : Aucun - documentation additionnelle

---

### 7. Amélioration des Messages d'Erreur (Optionnel)

**Fichier** : Créer `constraint/validation_errors.go`

**Action** : Structurer les erreurs de validation

```go
type ValidationError struct {
    Type     string // "type", "fact", "expression"
    Element  string // Nom de l'élément
    Field    string // Champ concerné (optionnel)
    Location string // Ligne/position (si disponible)
    Message  string // Message d'erreur
}

func (e *ValidationError) Error() string {
    if e.Field != "" {
        return fmt.Sprintf("%s '%s', champ '%s': %s", 
            e.Type, e.Element, e.Field, e.Message)
    }
    return fmt.Sprintf("%s '%s': %s", e.Type, e.Element, e.Message)
}
```

**Priorité** : 🟢 Basse (amélioration future)

**Impact** : Faible - meilleure UX

---

### 8. Support des Champs Optionnels (Future)

**Action** : Étendre TypeSystem pour supporter les champs optionnels

```go
type Field struct {
    Name         string `json:"name"`
    Type         string `json:"type"`
    IsPrimaryKey bool   `json:"isPrimaryKey,omitempty"`
    Optional     bool   `json:"optional,omitempty"` // NOUVEAU
}
```

**Priorité** : 🟢 Basse (fonctionnalité future)

**Impact** : Extension de fonctionnalité

---

## 🎯 Ordre Recommandé d'Exécution

1. **Étape 1** : Intégration dans `api.go` (fonction wrapper)
2. **Étape 2** : Intégration dans `constraint_program.go` (Option 1)
3. **Étape 4** : Tests d'intégration
4. **Étape 6** : Documentation utilisateur
5. **Étape 3** : Amélioration ActionValidator (si temps disponible)
6. **Étape 5** : Migration tests (si nécessaire)
7. **Étape 7** : Messages d'erreur structurés (amélioration future)
8. **Étape 8** : Champs optionnels (feature v2)

---

## 📝 Notes Importantes

### Rétrocompatibilité

Le nouveau système est **entièrement rétrocompatible** :
- Aucun fichier existant n'a été modifié
- Tous les tests existants passent
- Les nouvelles fonctions sont additionnelles

### Performance

Le système de validation ajoute un overhead minimal :
- Validation effectuée une seule fois au parsing
- Algorithmes efficaces (DFS O(V+E) pour cycles)
- Caching des types et variables

### Extensibilité

Le système est conçu pour être facilement extensible :
- Nouvelles règles de validation : ajouter dans ProgramValidator
- Nouveaux types de contraintes : ajouter dans validateConstraints
- Nouveaux opérateurs : ajouter dans AreTypesCompatible

---

## ✅ Checklist de Validation Post-Intégration

Après chaque intégration, vérifier :

- [ ] `go build ./constraint` - OK
- [ ] `go test ./constraint` - OK  
- [ ] `go test ./constraint -cover` - > 80%
- [ ] `make format` - OK
- [ ] `make lint` - OK
- [ ] Tests d'intégration existants - OK
- [ ] Documentation à jour

---

## 🔗 Références

- **Rapport d'implémentation** : `RAPPORT_TYPE_VALIDATION_SYSTEM.md`
- **Prompt source** : `scripts/new_ids/05-prompt-types-validation.md`
- **Standards projet** : `.github/prompts/common.md`
- **Guide de revue** : `.github/prompts/review.md`

---

**Statut** : 🟢 Prêt pour intégration progressive
