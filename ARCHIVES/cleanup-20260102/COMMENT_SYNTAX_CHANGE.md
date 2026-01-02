# Modification de la syntaxe des commentaires

**Date** : 16 décembre 2024  
**Auteur** : Modifications automatiques via script  
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