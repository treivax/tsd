# Prompt 09 : Mise à jour des exemples et fixtures

**Objectif** : Mettre à jour tous les exemples, fixtures et fichiers de démonstration du projet pour utiliser la nouvelle syntaxe des clés primaires et refléter la génération automatique des IDs.

**Prérequis** : Prompts 01-08 complétés et validés.

---

## Contexte

La nouvelle fonctionnalité de clés primaires introduit :
- La syntaxe `#` pour marquer les champs de clé primaire
- La génération automatique des IDs
- L'interdiction de spécifier `id` explicitement dans les assertions

Tous les exemples et fixtures du projet doivent être mis à jour pour :
1. Utiliser la nouvelle syntaxe `#` là où approprié
2. Retirer toute spécification explicite de `id`
3. Démontrer les différents cas d'usage (PK simple, composite, sans PK)
4. Documenter le comportement de génération d'IDs dans les commentaires

---

## Tâches

### 9.1. Inventaire des fichiers à mettre à jour

**Action** : Trouver tous les fichiers `.tsd`, exemples et fixtures dans le projet.

**Commandes** :
```bash
cd /home/resinsec/dev/tsd
find . -name "*.tsd" -type f
find . -path "*/examples/*" -name "*.tsd" -o -path "*/testdata/*" -name "*.tsd"
find . -path "*/fixtures/*" -type f
```

**Créer une liste** : Noter tous les fichiers trouvés et leur emplacement.

### 9.2. Catégoriser les fichiers

**Action** : Pour chaque fichier trouvé, déterminer la stratégie de mise à jour :

1. **Ajouter des PK simples** : Pour les types avec un identifiant naturel unique
2. **Ajouter des PK composites** : Pour les types avec plusieurs champs formant une clé
3. **Laisser sans PK** : Pour les types sans identifiant naturel (utiliseront le hash)
4. **Mise à jour minimale** : Retirer les `id` explicites uniquement

### 9.3. Mise à jour des exemples dans `examples/`

**Rechercher** : Les fichiers dans le répertoire `examples/` (s'il existe).

**Exemple 1** : `examples/basic.tsd` (ou similaire)

**Avant** :
```
type Person(nom: string, age: number)

assert Person(id: "person_1", nom: "Alice", age: 30)
assert Person(id: "person_2", nom: "Bob", age: 25)

rule Adults {
    when {
        p: Person()
        p.age >= 18
    }
    then {
        print(p.nom + " is an adult")
    }
}
```

**Après** :
```
// Exemple basique : type avec clé primaire simple
// Les IDs seront générés automatiquement au format "Person~nom"
type Person(#nom: string, age: number)

assert Person(nom: "Alice", age: 30)  // ID: Person~Alice
assert Person(nom: "Bob", age: 25)    // ID: Person~Bob

rule Adults {
    when {
        p: Person()
        p.age >= 18
    }
    then {
        print(p.nom + " is an adult (ID: " + p.id + ")")
    }
}
```

**Exemple 2** : `examples/composite_key.tsd` (créer si n'existe pas)

```
// Exemple de clé primaire composite
// Les IDs seront générés au format "TypeName~valeur1_valeur2_..."
type Produit(#categorie: string, #nom: string, prix: number, stock: number)

assert Produit(categorie: "Electronique", nom: "Laptop", prix: 1200, stock: 5)
// ID généré: Produit~Electronique_Laptop

assert Produit(categorie: "Electronique", nom: "Souris", prix: 25, stock: 50)
// ID généré: Produit~Electronique_Souris

assert Produit(categorie: "Livre", nom: "TSD Guide", prix: 30, stock: 100)
// ID généré: Produit~Livre_TSD%20Guide

rule StockFaible {
    when {
        p: Produit()
        p.stock < 10
    }
    then {
        print("ALERTE: Stock faible pour " + p.categorie + "/" + p.nom)
        print("  ID du produit: " + p.id)
    }
}
```

**Exemple 3** : `examples/no_primary_key.tsd` (créer si n'existe pas)

```
// Exemple sans clé primaire (utilise un hash)
// Les IDs seront générés au format "TypeName~<hash>"
type LogEntry(timestamp: number, level: string, message: string)

assert LogEntry(timestamp: 1704067200, level: "INFO", message: "Application started")
// ID généré: LogEntry~a1b2c3d4e5f6g7h8 (exemple de hash)

assert LogEntry(timestamp: 1704067201, level: "WARN", message: "High memory usage")
// ID généré: LogEntry~i9j0k1l2m3n4o5p6 (hash différent)

assert LogEntry(timestamp: 1704067202, level: "ERROR", message: "Connection failed")
// ID généré: LogEntry~q7r8s9t0u1v2w3x4 (hash différent)

rule ErrorLogs {
    when {
        log: LogEntry()
        log.level == "ERROR"
    }
    then {
        print("Error detected (ID: " + log.id + "): " + log.message)
    }
}
```

**Exemple 4** : `examples/relationships.tsd` (créer ou mettre à jour)

```
// Exemple de relations entre types avec IDs
type User(#username: string, email: string, created_at: number)
type Post(#post_id: number, author_username: string, title: string, content: string)
type Comment(#comment_id: number, post_id: number, author_username: string, text: string)

// Utilisateurs
assert User(username: "alice", email: "alice@example.com", created_at: 1704067200)
// ID: User~alice

assert User(username: "bob", email: "bob@example.com", created_at: 1704067300)
// ID: User~bob

// Posts
assert Post(post_id: 1, author_username: "alice", title: "Hello World", content: "My first post")
// ID: Post~1

assert Post(post_id: 2, author_username: "bob", title: "TSD is awesome", content: "Check this out")
// ID: Post~2

// Commentaires
assert Comment(comment_id: 1, post_id: 1, author_username: "bob", text: "Great post!")
// ID: Comment~1

assert Comment(comment_id: 2, post_id: 1, author_username: "alice", text: "Thanks!")
// ID: Comment~2

// Règle: Trouver les posts avec leurs auteurs
rule PostsWithAuthors {
    when {
        p: Post()
        u: User()
        p.author_username == u.username
    }
    then {
        print("Post '" + p.title + "' by " + u.email)
        print("  Post ID: " + p.id + ", User ID: " + u.id)
    }
}

// Règle: Commentaires sur les posts
rule CommentsOnPosts {
    when {
        p: Post()
        c: Comment()
        c.post_id == p.post_id
    }
    then {
        print("Comment on '" + p.title + "': " + c.text)
    }
}
```

**Exemple 5** : `examples/escaping.tsd` (créer)

```
// Exemple d'échappement de caractères spéciaux dans les clés primaires
// Les caractères ~, _, et % sont échappés dans les IDs
type Resource(#path: string, size: number)

assert Resource(path: "/home/user", size: 1024)
// ID: Resource~%2Fhome%2Fuser

assert Resource(path: "/tmp/file~backup", size: 2048)
// ID: Resource~%2Ftmp%2Ffile%7Ebackup (~ devient %7E)

assert Resource(path: "/var/log_archive", size: 4096)
// ID: Resource~%2Fvar%2Flog%5Farchive (_ devient %5F)

rule LargeResources {
    when {
        r: Resource()
        r.size > 2000
    }
    then {
        print("Large resource: " + r.path + " (ID: " + r.id + ")")
    }
}
```

### 9.4. Mise à jour des fixtures de test

**Rechercher** : Les fichiers dans `testdata/`, `fixtures/`, ou similaires.

**Pour chaque fichier** :

1. **Identifier les types** qui devraient avoir des clés primaires
2. **Ajouter `#`** aux champs appropriés
3. **Retirer les `id` explicites** des assertions
4. **Ajouter des commentaires** expliquant les IDs générés (pour la documentation)

**Exemple** : `testdata/sample.tsd`

**Avant** :
```
type Employee(id: string, name: string, department: string, salary: number)

assert Employee(id: "emp_001", name: "Alice", department: "Engineering", salary: 80000)
assert Employee(id: "emp_002", name: "Bob", department: "Sales", salary: 60000)
```

**Après** :
```
type Employee(#employee_id: string, name: string, department: string, salary: number)

assert Employee(employee_id: "emp_001", name: "Alice", department: "Engineering", salary: 80000)
assert Employee(employee_id: "emp_002", name: "Bob", department: "Sales", salary: 60000)
```

Ou bien, si l'ID n'est pas un champ métier :

```
type Employee(#name: string, department: string, salary: number)

assert Employee(name: "Alice", department: "Engineering", salary: 80000)
// ID: Employee~Alice

assert Employee(name: "Bob", department: "Sales", salary: 60000)
// ID: Employee~Bob
```

### 9.5. Mise à jour de la documentation embarquée

**Fichiers potentiels** :
- `README.md` (sections d'exemples)
- `docs/examples.md`
- `docs/syntax.md`
- Commentaires dans les fichiers de test

**Action** : Rechercher et mettre à jour toute documentation qui montre des exemples de code TSD.

**Exemple de mise à jour dans README.md** :

**Avant** :
```markdown
## Example

```tsd
type Person(name: string, age: number)
assert Person(id: "1", name: "Alice", age: 30)
```
```

**Après** :
```markdown
## Example

```tsd
// Define a type with a primary key
type Person(#name: string, age: number)

// IDs are generated automatically from primary keys
assert Person(name: "Alice", age: 30)  // ID: Person~Alice
assert Person(name: "Bob", age: 25)    // ID: Person~Bob

// You can access the generated ID in rules
rule PrintPerson {
    when {
        p: Person()
    }
    then {
        print(p.name + " (ID: " + p.id + ")")
    }
}
```

### Primary Keys

TSD supports automatic ID generation based on primary keys:

- **Simple primary key**: `type Person(#name: string, age: number)`
  - Generated ID format: `Person~Alice`
  
- **Composite primary key**: `type Product(#category: string, #name: string, price: number)`
  - Generated ID format: `Product~Electronics_Laptop`
  
- **No primary key**: `type Event(timestamp: number, message: string)`
  - Generated ID format: `Event~a1b2c3d4e5f6g7h8` (hash-based)

The `id` field is always available in rules and is of type string.
```

### 9.6. Créer un guide de migration

**Fichier** : `docs/MIGRATION_IDS.md` (créer)

```markdown
# Migration Guide: Automatic ID Generation

This guide helps you migrate existing TSD programs to use the new automatic ID generation feature.

## What Changed

### Before (Old Syntax)
```tsd
type Person(name: string, age: number)
assert Person(id: "person_1", name: "Alice", age: 30)
```

### After (New Syntax)
```tsd
type Person(#name: string, age: number)
assert Person(name: "Alice", age: 30)  // ID auto-generated: Person~Alice
```

## Migration Steps

### 1. Identify Natural Keys

For each type, identify if there's a natural unique identifier:
- User → username or email
- Product → SKU or (category + name)
- Order → order_number
- Event → May not have a natural key (use hash)

### 2. Add Primary Key Markers

Add `#` before field names that form the primary key:

```tsd
// Simple key
type User(#username: string, email: string)

// Composite key
type Product(#category: string, #name: string, price: number)

// No natural key (will use hash)
type Event(timestamp: number, message: string)
```

### 3. Remove Explicit IDs

Remove all `id: "..."` from your assertions:

**Before:**
```tsd
assert Person(id: "p1", name: "Alice", age: 30)
```

**After:**
```tsd
assert Person(name: "Alice", age: 30)
```

### 4. Update ID References in Rules

If your rules reference `id`, they will continue to work but now use the generated IDs:

```tsd
rule CheckID {
    when {
        p: Person()
        p.id == "Person~Alice"  // Use new format
    }
    then {
        print("Found Alice")
    }
}
```

## ID Format Reference

### With Primary Key

Format: `TypeName~value` or `TypeName~value1_value2_...`

Examples:
- `Person~Alice`
- `Product~Electronics_Laptop`
- `Order~2024_12345`

Special characters (`~`, `_`, `%`) are escaped:
- `/` → `%2F`
- `~` → `%7E`
- `_` → `%5F`
- `%` → `%25`

### Without Primary Key (Hash)

Format: `TypeName~<16-char-hex-hash>`

Example:
- `Event~a1b2c3d4e5f6g7h8`

The hash is deterministic: same field values always produce the same hash.

## Common Patterns

### Pattern 1: User Management
```tsd
type User(#username: string, email: string, role: string)
type Session(#session_id: string, username: string, active: bool)

assert User(username: "alice", email: "alice@example.com", role: "admin")
assert Session(session_id: "sess_123", username: "alice", active: true)

rule ActiveAdminSessions {
    when {
        u: User()
        s: Session()
        u.username == s.username
        u.role == "admin"
        s.active == true
    }
    then {
        print("Active admin session: " + u.username)
    }
}
```

### Pattern 2: Product Catalog
```tsd
type Product(#sku: string, name: string, price: number, stock: number)
type Order(#order_id: string, product_sku: string, quantity: number)

assert Product(sku: "LAPTOP-001", name: "Gaming Laptop", price: 1200, stock: 5)
assert Order(order_id: "ORD-2024-001", product_sku: "LAPTOP-001", quantity: 2)

rule CheckInventory {
    when {
        p: Product()
        o: Order()
        o.product_sku == p.sku
        p.stock < o.quantity
    }
    then {
        print("Insufficient stock for order " + o.order_id)
    }
}
```

### Pattern 3: Event Logging (No Primary Key)
```tsd
type LogEntry(timestamp: number, level: string, message: string)

assert LogEntry(timestamp: 1704067200, level: "INFO", message: "App started")
assert LogEntry(timestamp: 1704067201, level: "ERROR", message: "Connection failed")

rule ErrorHandler {
    when {
        log: LogEntry()
        log.level == "ERROR"
    }
    then {
        print("Error: " + log.message + " (ID: " + log.id + ")")
    }
}
```

## Backward Compatibility

Programs without primary key markers (`#`) will continue to work:
- All IDs will be generated using the hash method
- Existing programs remain valid
- No breaking changes for programs that don't use explicit `id`

## Troubleshooting

### Error: "field 'id' is reserved"
You tried to set `id` explicitly in an assertion. Remove it:
```tsd
// ❌ Wrong
assert Person(id: "custom", name: "Alice", age: 30)

// ✅ Correct
assert Person(name: "Alice", age: 30)
```

### Error: "primary key field 'X' not found in fact"
You marked a field as primary key but didn't provide it in the assertion:
```tsd
type Person(#name: string, age: number)

// ❌ Wrong
assert Person(age: 30)

// ✅ Correct
assert Person(name: "Alice", age: 30)
```

### Error: "primary key field must be a primitive type"
Primary keys can only be string, number, or bool:
```tsd
// ❌ Wrong
type Person(#data: object, age: number)

// ✅ Correct
type Person(#id: string, age: number)
```

## Questions?

For more information, see:
- [Syntax Documentation](syntax.md)
- [Examples](../examples/)
- [Test Cases](../testdata/integration/)
```

### 9.7. Vérification et validation

**Action** : Pour chaque fichier mis à jour, valider qu'il parse et fonctionne correctement.

**Commande** :
```bash
# Pour chaque fichier .tsd mis à jour
for file in examples/*.tsd testdata/**/*.tsd; do
    echo "Testing $file"
    ./tsd validate "$file" || echo "FAILED: $file"
done
```

Ou bien, si pas de CLI :
```bash
# Créer un script de test
cat > test_updated_files.sh << 'EOF'
#!/bin/bash
cd /home/resinsec/dev/tsd

for file in $(find . -name "*.tsd"); do
    echo "Checking $file..."
    # Utiliser un test Go pour parser le fichier
    go test ./constraint -run TestParseFile -args "$file"
done
EOF

chmod +x test_updated_files.sh
./test_updated_files.sh
```

---

## Validation

### Étape 1 : Lister tous les fichiers à mettre à jour

```bash
cd /home/resinsec/dev/tsd
find . -name "*.tsd" -type f > files_to_update.txt
cat files_to_update.txt
```

### Étape 2 : Mettre à jour les fichiers par catégorie

1. Mettre à jour `examples/`
2. Mettre à jour `testdata/`
3. Mettre à jour `fixtures/` (si existe)
4. Créer les nouveaux exemples

### Étape 3 : Valider chaque fichier

```bash
# Option 1 : Avec le CLI TSD (si disponible)
for file in $(cat files_to_update.txt); do
    ./tsd validate "$file" || echo "ERROR: $file"
done

# Option 2 : Avec les tests Go
go test ./constraint -run TestParseValidExamples -v
```

### Étape 4 : Mettre à jour la documentation

1. Éditer `README.md`
2. Éditer les fichiers dans `docs/`
3. Créer `docs/MIGRATION_IDS.md`

### Étape 5 : Vérifier que tous les tests passent

```bash
make test-complete
```

### Étape 6 : Validation finale

```bash
make validate
```

---

## Checklist

- [ ] Inventaire de tous les fichiers `.tsd` effectué
- [ ] Fichiers catégorisés par stratégie de mise à jour
- [ ] Exemples dans `examples/` mis à jour
- [ ] Nouvel exemple `composite_key.tsd` créé
- [ ] Nouvel exemple `no_primary_key.tsd` créé
- [ ] Nouvel exemple `relationships.tsd` créé
- [ ] Nouvel exemple `escaping.tsd` créé
- [ ] Fixtures dans `testdata/` mises à jour
- [ ] Documentation dans `README.md` mise à jour
- [ ] Documentation dans `docs/` mise à jour
- [ ] Guide de migration `docs/MIGRATION_IDS.md` créé
- [ ] Tous les fichiers `.tsd` validés (parsing réussit)
- [ ] Commentaires ajoutés pour expliquer les IDs générés
- [ ] `make test-complete` réussit
- [ ] `make validate` réussit

---

## Rapport

Une fois toutes les tâches complétées :

1. Lister tous les fichiers mis à jour
2. Lister les nouveaux fichiers créés
3. Documenter les choix de clés primaires pour chaque type
4. Copier la sortie des validations réussies
5. Commit :

```bash
git add examples/ testdata/ docs/ README.md
git commit -m "docs(examples): mise à jour pour clés primaires et génération automatique d'IDs

Mise à jour des exemples:
- examples/basic.tsd: ajout de PK simple sur nom
- examples/composite_key.tsd: nouveau (PK composite)
- examples/no_primary_key.tsd: nouveau (génération par hash)
- examples/relationships.tsd: nouveau (relations entre types)
- examples/escaping.tsd: nouveau (échappement de caractères)

Mise à jour des fixtures:
- testdata/: suppression des id explicites, ajout de PK
- Tous les fichiers .tsd valident correctement

Documentation:
- README.md: exemples mis à jour, section PK ajoutée
- docs/MIGRATION_IDS.md: guide de migration créé
- Commentaires ajoutés dans tous les exemples

Tous les fichiers ont été testés et valident correctement.

Refs #<issue_number>"
```

---

## Dépendances

- **Bloque** : Prompt 10 (documentation finale)
- **Bloqué par** : Prompts 01-08

---

## Notes importantes

1. **Cohérence** : Assurez-vous que tous les exemples utilisent une convention cohérente pour nommer les champs de clés primaires (préférer les noms significatifs comme `username`, `sku`, `order_id` plutôt que juste `id`).

2. **Commentaires** : Ajoutez des commentaires dans les exemples pour expliquer :
   - Quels champs sont des clés primaires et pourquoi
   - Quel sera le format de l'ID généré
   - Comment accéder à l'ID dans les règles

3. **Variété** : Créez des exemples montrant tous les cas d'usage :
   - PK simple (1 champ)
   - PK composite (2+ champs)
   - Sans PK (hash)
   - Différents types primitifs (string, number, bool)
   - Échappement de caractères spéciaux

4. **Réalisme** : Les exemples doivent refléter des cas d'usage réels pour aider les utilisateurs à comprendre quand utiliser quelle stratégie.

5. **Migration** : Le guide de migration doit être clair et pratique, avec des exemples avant/après pour chaque cas.

6. **Validation** : Ne committez que les fichiers qui parsent et valident correctement. Si un fichier a des problèmes, documentez-les et corrigez-les avant de continuer.

7. **Backward Compatibility** : Gardez au moins un exemple montrant qu'un type sans `#` fonctionne toujours (utilise le hash automatiquement).

---

**Prêt à passer au prompt 10 après validation complète.**