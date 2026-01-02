# Changements de Syntaxe TSD

Ce document regroupe tous les changements de syntaxe apportés au langage TSD.

**Version** : 2.0.0  
**Date** : 2025-01-02

---

## 1. Syntaxe Update avec Object Literal

**Date** : 2025-01-02  
**Version** : 2.0.0

### Problème Initial

Auparavant, l'action `Update` utilisait une syntaxe non-naturelle avec `_Mods_` :

```tsd
Update(_Mods_(variable, field1: value1, field2: value2))
```

### Solution

Utilisation de la syntaxe naturelle avec object literal `{...}` :

```tsd
Update(variable, { field1: value1, field2: value2 })
```

### Détails Techniques

La syntaxe de l'action `Update` a été améliorée pour utiliser des **objets littéraux** `{...}` au lieu de la syntaxe personnalisée `_Mods_(...)`. Cette nouvelle syntaxe est plus naturelle, cohérente avec le reste du langage, et plus facile à lire.

## Motivation

L'ancienne syntaxe `Update(p, _Mods_(statut: "actif"))` utilisait un pseudo-type `_Mods_` spécial qui :
- N'était pas intuitif pour les utilisateurs
- Créait une exception syntaxique dans le langage
- Ne correspondait à aucune construction standard

La grammaire PEG de TSD supporte déjà les objets littéraux (`ObjectLiteral`) via la syntaxe `{...}`. Utiliser cette syntaxe existante pour `Update` rend le langage plus cohérent et plus facile à apprendre.

## Nouvelle Syntaxe

### Syntaxe de Base

```tsd
Update(variable, {champ: valeur, champ2: valeur2, ...})
```

### Exemples

#### Mise à jour d'un seul champ

```tsd
type Personne(#nom: string, age: number, statut: string)

rule anniversaire : {p: Personne} / p.statut == "anniversaire" ==>
    Update(p, {age: p.age + 1.0})
```

#### Mise à jour de plusieurs champs

```tsd
rule activation : {p: Personne} / p.statut == "nouveau" ==>
    Update(p, {statut: "actif", ville: "Paris"})
```

#### Avec expressions et accès aux champs

```tsd
type Relation(personne1: string, personne2: string, lien: string)

rule mettre_en_couple : {p: Personne, r: Relation} /
    p.nom == r.personne1 AND r.lien == "mariage" ==>
    Update(p, {statut: "en couple"})
```

## Ancienne Syntaxe (Déconseillée)

```tsd
// ❌ Ancienne syntaxe avec _Mods_
Update(p, _Mods_(statut: "actif"))

// ✅ Nouvelle syntaxe avec objet littéral
Update(p, {statut: "actif"})
```

## Sémantique

L'action `Update` avec objet littéral :

1. **Préserve l'identifiant du fait** : contrairement à l'ancienne syntaxe `Update(Personne(...))` qui créait un nouveau fait avec un nouvel ID, la nouvelle syntaxe modifie le fait existant en conservant son `_id_`.

2. **Modifie uniquement les champs spécifiés** : les champs non mentionnés dans l'objet littéral conservent leur valeur actuelle.

3. **Propage les changements dans RETE** : la mise à jour déclenche la réévaluation des règles dépendantes dans le réseau.

4. **Valide les types** : les champs spécifiés doivent exister dans le type du fait et respecter les types définis.

## Implémentation Technique

### Parser PEG

La grammaire détecte automatiquement les appels `Update(variable, {...})` et les transforme en une structure AST spéciale `updateWithModifications` :

```go
// Dans constraint.peg, règle JobCall
if nameStr == "Update" {
    argsList, ok := args.([]interface{})
    if ok && len(argsList) == 2 {
        if secondArg, isMap := argsList[1].(map[string]interface{}); isMap {
            if secondArg["type"] == "objectLiteral" {
                // Transformation en updateWithModifications
                ...
            }
        }
    }
}
```

### Évaluation

L'évaluateur reconnaît `updateWithModifications` et :
1. Évalue la variable pour obtenir le fait original
2. Copie les champs du fait original
3. Applique les modifications spécifiées
4. Valide les champs modifiés
5. Retourne un nouveau fait avec le **même ID** mais les champs modifiés

```go
// Dans constraint/evaluator.go
func evaluateUpdateWithModifications(update map[string]interface{}, bindings map[string]interface{}) (map[string]interface{}, error) {
    // Préserve l'ID original
    result := make(map[string]interface{})
    result["_id_"] = originalFact["_id_"]
    
    // Applique les modifications
    for fieldName, fieldValueAST := range modifications {
        evaluatedValue, err := evaluateArithmeticExpr(fieldValueAST, bindings)
        result[fieldName] = evaluatedValue
    }
    
    return result, nil
}
```

## Migration

### Code Existant avec `_Mods_`

Remplacer simplement `_Mods_(...)` par `{...}` :

```diff
- Update(p, _Mods_(statut: "actif", age: 30.0))
+ Update(p, {statut: "actif", age: 30.0})
```

### Code Existant avec `Update(InlineFact(...))`

Cette syntaxe continue de fonctionner mais a une **sémantique différente** : elle crée un **nouveau fait** avec un **nouvel ID** au lieu de mettre à jour le fait existant.

```tsd
// Crée un NOUVEAU fait avec un nouvel ID
Update(Personne(nom: "Alice", statut: "actif"))

// Met à jour le fait EXISTANT (préserve l'ID)
Update(p, {statut: "actif"})
```

## Tests

Tous les tests ont été mis à jour :
- `tests/e2e/update_syntax_test.go` : tests de parsing de la nouvelle syntaxe
- `tests/e2e/testdata/update_simple.tsd` : exemple simple
- `tests/e2e/testdata/relationship_step1_types_rules.tsd` : cas réel
- `internal/defaultactions/loader_test.go` : signature de l'action Update

## Fichiers Modifiés

1. **Grammaire** : `constraint/grammar/constraint.peg`
   - Détection de `objectLiteral` comme second argument d'`Update`
   - Transformation en AST `updateWithModifications`

2. **Documentation** : `internal/defaultactions/defaults.tsd`
   - Mise à jour des exemples et commentaires

3. **Validateur** : `constraint/action_validator.go`
   - Commentaire mis à jour

4. **Tests** :
   - `tests/e2e/update_syntax_test.go`
   - `tests/e2e/testdata/*.tsd`
   - `internal/defaultactions/loader_test.go`

## Compatibilité

- ✅ **Parsing** : la nouvelle syntaxe est entièrement supportée
- ✅ **Validation** : les types sont validés correctement
- ✅ **Évaluation** : l'ID est préservé comme attendu
- ⚠️ **Migration** : rechercher et remplacer `_Mods_` par `{...}` dans le code existant

## Prochaines Étapes

1. ✅ Parser génère `updateWithModifications` pour `Update(var, {...})`
2. ✅ Évaluateur préserve l'ID du fait original
3. ✅ Tests E2E passent avec la nouvelle syntaxe
4. ⏳ Intégrer `BuiltinActionExecutor` pour exécution réelle des Updates
5. ⏳ Tests d'intégration vérifiant la propagation RETE après Update

## Conclusion

La nouvelle syntaxe `Update(variable, {champs...})` :
- Est plus **intuitive** et **naturelle**
- Réutilise une construction **existante** du langage
- Rend le code plus **lisible** et **maintenable**
- Préserve correctement l'**identité des faits**

Cette amélioration simplifie le langage TSD et le rend plus cohérent.

---

## 2. Syntaxe des Commentaires

**Date** : 2024-12-16  
**Version** : antérieur

### Syntaxe

```tsd
// Commentaire sur une ligne

/* Commentaire
   multi-lignes */
```

**Version** : TSD v2.x

---

## 🎯 Objectif

Modifier la grammaire TSD pour que le caractère `#` ne soit plus utilisé pour les commentaires. Les commentaires sont désormais introduits **uniquement** par `//` ou `/* */`.

Cette modification permet de réserver le caractère `#` pour marquer les **clés primaires** dans les définitions de types.

---

## 📝 Modifications apportées

### 1. Grammaire PEG (`constraint/grammar/constraint.peg`)

**Avant** :
```peg
Comment <- LineComment / BlockComment / EndOfLineComment

LineComment <- "//" CommentText:(![\r\n] .)* {
    return nil, nil
}

BlockComment <- "/*" CommentText:(!"*/" .)* "*/" {
    return nil, nil
}

EndOfLineComment <- "#" !IdentStart CommentText:(![\r\n] .)* {
    return nil, nil  // Support pour commentaires style shell/Python
}
```

**Après** :
```peg
Comment <- LineComment / BlockComment

LineComment <- "//" CommentText:(![\r\n] .)* {
    return nil, nil  // Les commentaires ne retournent rien
}

BlockComment <- "/*" CommentText:(!"*/" .)* "*/" {
    return nil, nil  // Les commentaires ne retournent rien
}
```

**Changement** : Suppression de la règle `EndOfLineComment` qui permettait les commentaires avec `#`.

### 2. Parser régénéré (`constraint/parser.go`)

Le parser a été régénéré avec la commande :
```bash
cd constraint/grammar
pigeon -o parser.go constraint.peg
cp parser.go ../parser.go
```

### 3. Fichiers `.tsd` mis à jour

**Tous les fichiers `.tsd` du projet** (136 fichiers) ont été automatiquement modifiés pour remplacer les commentaires `#` par `//`.

**Exemples de conversion** :

**Avant** :
```tsd
# Commentaire sur une ligne
type Person(id: string, name: string, age: number)

# Valid facts
Person(id: "P001", name: "Alice", age: 30)  # Commentaire en fin de ligne
```

**Après** :
```tsd
// Commentaire sur une ligne
type Person(id: string, name: string, age: number)

// Valid facts
Person(id: "P001", name: "Alice", age: 30)  // Commentaire en fin de ligne
```

### 4. Tests Go mis à jour

Fichiers de tests modifiés :
- `constraint/unicode_test.go` : Commentaire `#` remplacé par `//`
- `rete/constraint_pipeline_test.go` : Commentaire `#` remplacé par `//`
- `constraint/comment_changes_test.go` : Nouveaux tests ajoutés pour valider le changement

---

## ✅ Validation

### Tests de régression

```bash
# Tous les tests du module constraint passent
go test ./constraint/... -v
# PASS

# Tests spécifiques aux commentaires
go test ./constraint -run TestComments -v
# PASS
```

### Comportement vérifié

| Syntaxe | Avant | Après | Status |
|---------|-------|-------|--------|
| `// commentaire` | ✅ Fonctionne | ✅ Fonctionne | ✅ Conservé |
| `/* commentaire */` | ✅ Fonctionne | ✅ Fonctionne | ✅ Conservé |
| `# commentaire` | ✅ Fonctionne | ❌ **Rejeté** | ✅ Supprimé |
| `type T(#field: string)` | ❌ N/A | ✅ **Fonctionne** | ✅ Nouveau |

### Exemples de tests

```go
func TestCommentsWithSlashes(t *testing.T) {
    input := `
    // Commentaire ligne
    /* Commentaire bloc */
    type Person(name: string)
    `
    _, err := ParseConstraint("test.tsd", []byte(input))
    // ✅ PASS : Les commentaires // et /* */ fonctionnent
}

func TestHashAsCommentShouldFail(t *testing.T) {
    input := `# Ce commentaire ne devrait plus fonctionner`
    _, err := ParseConstraint("test.tsd", []byte(input))
    // ✅ PASS : Erreur de parsing attendue
}

func TestHashAsPrimaryKeyStillWorks(t *testing.T) {
    input := `type Person(#name: string, age: number)`
    result, err := ParseConstraint("test.tsd", []byte(input))
    // ✅ PASS : Le # pour les clés primaires fonctionne
}
```

---

## 🔄 Migration

### Pour les utilisateurs existants

Si vous avez des fichiers `.tsd` existants utilisant `#` pour les commentaires :

**Option 1 : Script automatique**
```bash
#!/bin/bash
find . -name "*.tsd" -type f | while read file; do
    sed -i 's/^\( *\)#\( .*\|$\)/\1\/\/\2/g' "$file"
done
```

**Option 2 : Remplacement manuel**
- Rechercher : `^\s*#`
- Remplacer : `//`

### Exemples de migration

**Avant** :
```tsd
# Définition des types
type User(id: string, name: string)

# Données de test
User(id: "U001", name: "Alice")  # Utilisateur admin
```

**Après** :
```tsd
// Définition des types
type User(id: string, name: string)

// Données de test
User(id: "U001", name: "Alice")  // Utilisateur admin
```

---

## 📋 Nouvelle utilisation de `#` : Clés primaires

Le caractère `#` est maintenant **réservé** pour marquer les champs de clé primaire :

```tsd
// Clé primaire simple
type User(#username: string, email: string, age: number)

// Clé primaire composite
type Product(#category: string, #name: string, price: number)

// Sans clé primaire (utilise un hash)
type Event(timestamp: number, message: string)
```

Les IDs de faits seront générés automatiquement :
- Avec PK : `User~alice`, `Product~Electronics_Laptop`
- Sans PK : `Event~a1b2c3d4e5f6g7h8` (hash MD5 tronqué)

Pour plus d'informations, voir la documentation sur les clés primaires.

---

## 🎯 Impact

### Breaking Changes

- ❌ **Les commentaires avec `#` ne fonctionnent plus**
- ✅ Migration simple : remplacer `#` par `//`
- ✅ Les commentaires `//` et `/* */` continuent de fonctionner

### Compatibilité ascendante

- ✅ Tous les programmes existants utilisant `//` et `/* */` fonctionnent sans modification
- ✅ Aucun impact sur la sémantique des règles, types ou faits

### Fichiers impactés

```
Fichiers modifiés :
- constraint/grammar/constraint.peg (grammaire)
- constraint/parser.go (parser régénéré)
- constraint/unicode_test.go (test)
- rete/constraint_pipeline_test.go (test)
- 136 fichiers .tsd (exemples, tests, fixtures)
```

---

## 📚 Références

- [Guide de migration clés primaires](../scripts/gestion-ids/09-prompt-maj-exemples.md)
- [Syntaxe des clés primaires](../scripts/gestion-ids/00-PLAN.md)
- [Documentation grammaire](constraint/grammar/constraint.peg)

---

## ✅ Checklist de validation

- [x] Grammaire PEG modifiée
- [x] Parser régénéré avec pigeon
- [x] Tous les fichiers `.tsd` migrés (136 fichiers)
- [x] Tests Go mis à jour
- [x] Nouveaux tests de validation ajoutés
- [x] Tous les tests constraint passent
- [x] Commentaires `//` et `/* */` fonctionnent
- [x] Commentaires `#` sont rejetés
- [x] Clés primaires avec `#` fonctionnent
- [x] Documentation créée

---

**Note** : Cette modification fait partie de la fonctionnalité plus large de gestion automatique des IDs basée sur les clés primaires. Voir `scripts/gestion-ids/` pour le plan complet.