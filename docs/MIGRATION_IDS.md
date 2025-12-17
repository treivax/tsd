# Guide de Migration : Génération Automatique des IDs

## Vue d'ensemble

Ce guide vous aide à migrer vos programmes TSD existants pour utiliser la nouvelle fonctionnalité de génération automatique d'identifiants basée sur des clés primaires.

## Qu'est-ce qui a changé ?

### Avant (Ancienne Syntaxe)

```tsd
type Person(name: string, age: number)
assert Person(id: "person_1", name: "Alice", age: 30)
```

**Problèmes** :
- IDs gérés manuellement
- Risque de collisions d'IDs
- IDs non standardisés
- Duplication de code pour générer des IDs

### Après (Nouvelle Syntaxe)

```tsd
type Person(#name: string, age: number)
assert Person(name: "Alice", age: 30)  // ID auto-généré: Person~Alice
```

**Avantages** :
- IDs générés automatiquement
- Format standardisé et prévisible
- Pas de collisions (basé sur clés primaires)
- Code plus propre et maintenable

---

## Étapes de Migration

### Étape 1 : Identifier les Identifiants Naturels

Pour chaque type, déterminez s'il existe un ou plusieurs champs qui forment un identifiant unique naturel.

**Questions à se poser** :
- Quel(s) champ(s) rend(ent) ce fait unique ?
- Y a-t-il un identifiant métier naturel ?
- Plusieurs champs sont-ils nécessaires pour l'unicité ?

**Exemples** :

| Type | Identifiant Naturel | Type de Clé |
|------|---------------------|-------------|
| User | username ou email | Simple |
| Product | SKU | Simple |
| Order | year + orderNumber | Composite |
| Location | country + city | Composite |
| LogEvent | (aucun) | Aucune (hash) |

### Étape 2 : Marquer les Clés Primaires

Ajoutez le préfixe `#` devant les champs qui forment la clé primaire.

#### Clé Primaire Simple

Un seul champ forme l'identifiant unique.

```tsd
// Avant
type User(username: string, email: string, age: number)

// Après
type User(#username: string, email: string, age: number)
```

**ID généré** : `User~alice`

#### Clé Primaire Composite

Plusieurs champs forment ensemble l'identifiant unique.

```tsd
// Avant
type Product(category: string, name: string, price: number)

// Après
type Product(#category: string, #name: string, price: number)
```

**ID généré** : `Product~Electronics_Laptop`

#### Aucune Clé Primaire (Hash)

Pas d'identifiant naturel → utilisation d'un hash.

```tsd
// Pas de changement dans la définition
type LogEvent(timestamp: number, level: string, message: string)
```

**ID généré** : `LogEvent~a1b2c3d4e5f6g7h8` (hash déterministe)

### Étape 3 : Retirer les IDs Explicites

Supprimez tous les champs `id: "..."` de vos assertions.

**Avant** :
```tsd
assert Person(id: "p1", name: "Alice", age: 30)
assert Person(id: "p2", name: "Bob", age: 25)
```

**Après** :
```tsd
assert Person(name: "Alice", age: 30)
assert Person(name: "Bob", age: 25)
```

### Étape 4 : Mettre à Jour les Références aux IDs

Si vos règles référencent des IDs, mettez-les à jour pour utiliser le nouveau format.

**Avant** :
```tsd
rule checkSpecificUser : {p: Person} / p.id == "p1"
    ==> print("Found user p1")
```

**Après** :
```tsd
rule checkSpecificUser : {p: Person} / p.id == "Person~Alice"
    ==> print("Found user Alice")
```

**Ou mieux, utilisez la clé primaire directement** :
```tsd
rule checkSpecificUser : {p: Person} / p.name == "Alice"
    ==> print("Found user Alice")
```

---

## Format des IDs Générés

### Avec Clé Primaire Simple

**Format** : `TypeName~valeur`

**Exemples** :
- `User~alice`
- `Product~LAPTOP-001`
- `Country~FR`
- `Student~2024001`

### Avec Clé Primaire Composite

**Format** : `TypeName~valeur1_valeur2_valeur3`

**Exemples** :
- `Product~Electronics_Laptop`
- `Order~2024_1001`
- `Location~France_Paris`
- `Enrollment~S2024001_CS101`

**Note** : Les valeurs sont séparées par underscore (`_`) dans l'ordre de déclaration.

### Sans Clé Primaire (Hash)

**Format** : `TypeName~<hash-16-caractères-hex>`

**Exemples** :
- `LogEvent~a1b2c3d4e5f6g7h8`
- `Metric~fedcba9876543210`

**Note** : Le hash est déterministe (mêmes valeurs → même hash).

### Échappement des Caractères Spéciaux

Certains caractères sont échappés en format URL-encoding :

| Caractère | Échappement | Raison |
|-----------|-------------|--------|
| `~` (tilde) | `%7E` | Séparateur type/valeur |
| `_` (underscore) | `%5F` | Séparateur composite |
| `%` (percent) | `%25` | Caractère d'échappement |
| ` ` (espace) | `%20` | Caractère spécial |
| `/` (slash) | `%2F` | Caractère spécial |

**Exemples** :
- `User~alice` → `User~alice` (pas d'échappement)
- `User~user~admin` → `User~user%7Eadmin`
- `Product~Books_Go Programming` → `Product~Books_Go%20Programming`
- `File~/home/user` → `File~%2Fhome%2Fuser`

---

## Exemples de Migration

### Exemple 1 : Gestion d'Utilisateurs

**Avant** :
```tsd
type User(name: string, email: string, role: string)

assert User(id: "u1", name: "Alice", email: "alice@example.com", role: "admin")
assert User(id: "u2", name: "Bob", email: "bob@example.com", role: "user")

rule adminUsers : {u: User} / u.role == "admin"
    ==> notify(u.id)
```

**Après** :
```tsd
type User(#username: string, email: string, role: string)

assert User(username: "alice", email: "alice@example.com", role: "admin")
assert User(username: "bob", email: "bob@example.com", role: "user")

rule adminUsers : {u: User} / u.role == "admin"
    ==> notify(u.id)  // u.id vaut maintenant "User~alice"
```

### Exemple 2 : Catalogue de Produits

**Avant** :
```tsd
type Product(sku: string, name: string, price: number, stock: number)

assert Product(id: "prod_1", sku: "LAPTOP-001", name: "Gaming Laptop", price: 1200, stock: 5)
assert Product(id: "prod_2", sku: "MOUSE-042", name: "Wireless Mouse", price: 25, stock: 100)

rule lowStock : {p: Product} / p.stock < 10
    ==> alert(p.id)
```

**Après** :
```tsd
type Product(#sku: string, name: string, price: number, stock: number)

assert Product(sku: "LAPTOP-001", name: "Gaming Laptop", price: 1200, stock: 5)
assert Product(sku: "MOUSE-042", name: "Wireless Mouse", price: 25, stock: 100)

rule lowStock : {p: Product} / p.stock < 10
    ==> alert(p.id)  // p.id vaut maintenant "Product~LAPTOP-001"
```

### Exemple 3 : Événements de Log (Sans Clé Primaire)

**Avant** :
```tsd
type LogEvent(timestamp: number, level: string, message: string)

assert LogEvent(id: "log_1", timestamp: 1704067200, level: "INFO", message: "App started")
assert LogEvent(id: "log_2", timestamp: 1704067201, level: "ERROR", message: "Connection failed")

rule errorLogs : {log: LogEvent} / log.level == "ERROR"
    ==> process(log.id)
```

**Après** :
```tsd
// Pas de clé primaire → utilisation du hash
type LogEvent(timestamp: number, level: string, message: string)

assert LogEvent(timestamp: 1704067200, level: "INFO", message: "App started")
assert LogEvent(timestamp: 1704067201, level: "ERROR", message: "Connection failed")

rule errorLogs : {log: LogEvent} / log.level == "ERROR"
    ==> process(log.id)  // log.id vaut maintenant "LogEvent~<hash>"
```

### Exemple 4 : Relations Entre Types

**Avant** :
```tsd
type User(username: string, email: string)
type Order(orderId: string, userId: string, total: number)

assert User(id: "u1", username: "alice", email: "alice@example.com")
assert Order(id: "o1", orderId: "ORD-001", userId: "u1", total: 1500)

rule userOrders : {u: User, o: Order} / u.id == o.userId
    ==> process(u.id, o.id)
```

**Après** :
```tsd
type User(#username: string, email: string)
type Order(#orderId: string, username: string, total: number)

assert User(username: "alice", email: "alice@example.com")
assert Order(orderId: "ORD-001", username: "alice", total: 1500)

// Jointure sur le champ username (clé primaire de User)
rule userOrders : {u: User, o: Order} / u.username == o.username
    ==> process(u.id, o.id)  // IDs: "User~alice", "Order~ORD-001"
```

---

## Compatibilité Descendante

### Programmes Sans Clés Primaires

Les programmes qui ne définissent pas de clés primaires (pas de `#`) continuent de fonctionner :
- Tous les IDs seront générés par hash
- Aucun changement de comportement
- Migration progressive possible

### Champ `id` Réservé

**IMPORTANT** : Le champ `id` est désormais réservé et généré automatiquement.

**Erreur** :
```tsd
// ❌ INTERDIT
type Person(id: string, name: string, age: number)
assert Person(id: "custom_id", name: "Alice", age: 30)
```

**Solution** :
```tsd
// ✅ CORRECT - Renommer le champ
type Person(personId: string, name: string, age: number)
assert Person(personId: "custom_id", name: "Alice", age: 30)

// ✅ MIEUX - Utiliser une clé primaire
type Person(#personId: string, name: string, age: number)
assert Person(personId: "custom_id", name: "Alice", age: 30)
```

---

## Dépannage

### Erreur : "field 'id' is reserved"

**Cause** : Vous avez essayé de définir ou d'assigner un champ nommé `id`.

**Solution** :
1. Renommez le champ (ex: `personId`, `userId`, `recordId`)
2. Ou utilisez-le comme clé primaire avec `#`

### Erreur : "primary key field 'X' not found in fact"

**Cause** : Vous avez marqué un champ comme clé primaire (`#X`) mais ne l'avez pas fourni dans l'assertion.

**Exemple** :
```tsd
type Person(#name: string, age: number)

// ❌ Erreur - 'name' manquant
assert Person(age: 30)

// ✅ Correct
assert Person(name: "Alice", age: 30)
```

### Erreur : "primary key field must be a primitive type"

**Cause** : Les clés primaires doivent être de type primitif (string, number, bool).

**Exemple** :
```tsd
// ❌ Erreur - 'data' est de type object
type Person(#data: object, age: number)

// ✅ Correct - utiliser un champ primitif
type Person(#id: string, age: number)
```

### IDs Non Lisibles (Hash)

**Problème** : Les IDs générés sont des hash non lisibles.

**Cause** : Aucune clé primaire n'est définie.

**Solution** : Identifiez un champ qui peut servir de clé primaire et marquez-le avec `#`.

**Avant** :
```tsd
type User(username: string, email: string)
// ID généré: User~a1b2c3d4e5f6g7h8
```

**Après** :
```tsd
type User(#username: string, email: string)
// ID généré: User~alice
```

---

## Bonnes Pratiques

### 1. Choisir de Bonnes Clés Primaires

**✅ Faire** :
- Utiliser des identifiants métier naturels (username, SKU, code ISO)
- Choisir des valeurs stables (qui ne changent pas)
- Préférer des valeurs courtes et lisibles
- Éviter les caractères spéciaux si possible

**❌ Éviter** :
- Champs volatiles (timestamps, compteurs)
- Valeurs avec beaucoup de caractères spéciaux
- Champs non uniques
- Trop de champs dans une clé composite (limiter à 2-3)

### 2. Documenter les Choix

Ajoutez des commentaires pour expliquer pourquoi un champ est (ou n'est pas) une clé primaire.

```tsd
// Username est unique et stable - utilisé comme clé primaire
type User(#username: string, email: string, role: string)

// Pas d'identifiant naturel - les événements sont éphémères
// IDs générés par hash
type LogEvent(timestamp: number, level: string, message: string)

// Clé composite : un produit est unique par (catégorie, nom)
type Product(#category: string, #name: string, price: number)
```

### 3. Tester avec des Données Réalistes

Testez vos types avec des valeurs représentatives, notamment :
- Valeurs avec caractères spéciaux
- Valeurs avec espaces
- Valeurs unicode
- Valeurs longues

### 4. Utiliser le Champ `id` dans les Règles

Le champ `id` est toujours disponible pour le logging et la traçabilité.

```tsd
rule processUser : {u: User} / u.age >= 18
    ==> logActivity("Processing user " + u.username + " (ID: " + u.id + ")")
```

### 5. Migration Progressive

Vous pouvez migrer vos types progressivement :
1. Commencez par les types simples avec identifiants naturels évidents
2. Continuez avec les types à clés composites
3. Laissez les types sans identifiant naturel utiliser le hash
4. Testez chaque étape avant de continuer

---

## Checklist de Migration

- [ ] Inventaire de tous les types définis
- [ ] Identification des clés primaires pour chaque type
- [ ] Ajout des marqueurs `#` sur les champs de clés primaires
- [ ] Suppression des `id:` explicites dans les assertions
- [ ] Mise à jour des références aux IDs dans les règles
- [ ] Renommage des champs nommés `id` (si existants)
- [ ] Ajout de commentaires documentant les choix
- [ ] Tests avec données réalistes
- [ ] Validation que tous les tests passent
- [ ] Mise à jour de la documentation

---

## Ressources

- **Exemples** : Consultez le répertoire `examples/` pour voir des cas d'usage variés
  - `examples/pk_simple.tsd` - Clés primaires simples
  - `examples/pk_composite.tsd` - Clés primaires composites
  - `examples/pk_none.tsd` - Génération par hash
  - `examples/pk_special_chars.tsd` - Échappement de caractères
  - `examples/pk_relationships.tsd` - Relations entre types

- **Tests** : Consultez `tests/fixtures/integration/` pour des exemples de tests

- **Documentation** : 
  - `README.md` - Documentation générale
  - `docs/syntax.md` - Syntaxe complète (si disponible)

---

## Support

Si vous rencontrez des problèmes lors de la migration :

1. Vérifiez que votre syntaxe est correcte
2. Consultez les exemples fournis
3. Testez avec un fichier minimal pour isoler le problème
4. Vérifiez les messages d'erreur (ils sont explicites)

---

**Note** : Cette migration améliore la maintenabilité et la cohérence de vos programmes TSD. Prenez le temps de bien identifier vos clés primaires pour bénéficier pleinement de cette fonctionnalité.