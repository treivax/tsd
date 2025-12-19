# Règles Complètes de Gestion des IDs dans TSD

## Vue d'Ensemble

Dans TSD, **chaque fait possède automatiquement un champ `id` de type `string`** qui est généré automatiquement par le système. Ce document décrit exhaustivement toutes les règles de fonctionnement de ces identifiants.

---

## 📋 Table des Matières

1. [Principes Fondamentaux](#principes-fondamentaux)
2. [Le Champ `id` - Champ Réservé](#le-champ-id---champ-réservé)
3. [Clés Primaires](#clés-primaires)
4. [Génération d'IDs](#génération-dids)
5. [Format des IDs](#format-des-ids)
6. [Échappement de Caractères](#échappement-de-caractères)
7. [Cas Particuliers](#cas-particuliers)
8. [Utilisation dans les Règles](#utilisation-dans-les-règles)
9. [Erreurs Courantes](#erreurs-courantes)
10. [Exemples Complets](#exemples-complets)

---

## Principes Fondamentaux

### 1. Le Champ `id` est TOUJOURS Présent

```tsd
type Person(name: string, age: number)
assert Person(name: "Alice", age: 30)

rule showPerson : {p: Person} / true
    ==> print(p.id)  // ✅ TOUJOURS disponible
```

**Garanties :**
- ✅ Tous les faits ont un champ `id`
- ✅ Le champ `id` est toujours de type `string`
- ✅ Le champ `id` est généré automatiquement
- ✅ Le champ `id` est accessible dans toutes les expressions

### 2. L'ID est Unique par Fait

Chaque instance de fait possède un ID unique dans le réseau RETE.

```tsd
type User(#username: string, email: string)

assert User(username: "alice", email: "alice@example.com")
// ID généré: "User~alice"

assert User(username: "bob", email: "bob@example.com")
// ID généré: "User~bob"
```

### 3. L'ID est Déterministe

Pour des valeurs identiques, l'ID généré sera toujours le même.

```tsd
type Product(#sku: string, name: string)

// Premier assert
assert Product(sku: "ABC123", name: "Laptop")
// ID: "Product~ABC123"

// Deuxième assert avec les mêmes valeurs
assert Product(sku: "ABC123", name: "Laptop")
// ID: "Product~ABC123" (identique)
```

**⚠️ Important :** Avec des clés primaires, deux faits avec les mêmes valeurs de clé primaire auront le même ID (et seront considérés comme le même fait dans RETE).

---

## Le Champ `id` - Champ Réservé

### Règle Absolue : Le Champ `id` est RÉSERVÉ

Le nom `id` est un **champ réservé** du système. Vous **NE POUVEZ PAS** :
- ❌ Définir un champ nommé `id` dans un type
- ❌ Assigner manuellement une valeur au champ `id`
- ❌ Utiliser `id` comme nom de clé primaire sans le déclarer explicitement

### ❌ INTERDIT - Définir un Champ `id`

```tsd
// ❌ ERREUR : Le champ 'id' est réservé
type Person(id: string, name: string, age: number)

// Erreur lors du parsing :
// "field 'id' is reserved and cannot be defined manually"
```

### ❌ INTERDIT - Assigner Manuellement un ID

```tsd
type Person(name: string, age: number)

// ❌ ERREUR : On ne peut pas définir 'id' manuellement
assert Person(id: "custom-id", name: "Alice", age: 30)

// Erreur lors de la validation :
// "fait de type 'Person': le champ 'id' ne peut pas être défini manuellement
//  (il est généré automatiquement)"
```

### ✅ EXCEPTION - Déclarer `id` comme Clé Primaire

La **seule exception** est de déclarer explicitement `id` comme clé primaire :

```tsd
// ✅ AUTORISÉ : 'id' déclaré comme clé primaire
type Person(#id: string, name: string, age: number)

assert Person(id: "person-001", name: "Alice", age: 30)
// ID généré: "Person~person-001"
```

**Note :** Dans ce cas, vous DEVEZ fournir la valeur de `id` dans chaque assertion.

### Type du Champ `id`

Le champ `id` virtuel est **TOUJOURS de type `string`**, même si la clé primaire est d'un autre type :

```tsd
type Product(#productId: number, name: string)
assert Product(productId: 123, name: "Laptop")

rule checkProduct : {p: Product} / true ==> {
    print(p.id)          // "Product~123" (string, pas number)
    print(p.productId)   // 123 (number)
}
```

**Implications :**

```tsd
// ❌ ERREUR : Comparaison de string avec number
rule invalid : {p: Product} / p.id > 0 ==> print("found")
// Erreur : "type incompatibility in comparison: string vs number"

// ✅ CORRECT : Comparer la clé primaire directement
rule valid : {p: Product} / p.productId > 0 ==> print("found")

// ✅ CORRECT : Vérifier que l'id n'est pas vide
rule valid2 : {p: Product} / p.id != "" ==> print("found")
```

---

## Clés Primaires

### Définition des Clés Primaires

Une clé primaire est un ou plusieurs champs qui identifient de manière **unique** un fait.

**Syntaxe :** Préfixer le nom du champ avec `#`

```tsd
type TypeName(#field1: type1, field2: type2, ...)
```

### Clé Primaire Simple

Un seul champ forme l'identifiant unique.

```tsd
type User(#username: string, email: string, age: number)
```

**Caractéristiques :**
- ✅ Un seul champ marqué avec `#`
- ✅ Ce champ doit être de type primitif (string, number, bool)
- ✅ Ce champ DOIT être fourni dans chaque assertion

### Clé Primaire Composite

Plusieurs champs forment ensemble l'identifiant unique.

```tsd
type Product(#category: string, #name: string, price: number)
```

**Caractéristiques :**
- ✅ Plusieurs champs marqués avec `#`
- ✅ Tous les champs doivent être de type primitif
- ✅ TOUS les champs DOIVENT être fournis dans chaque assertion
- ✅ L'ordre des champs dans le type définit l'ordre dans l'ID

### Absence de Clé Primaire

Si aucun champ n'est marqué comme clé primaire, l'ID est généré par **hash**.

```tsd
type LogEvent(timestamp: number, level: string, message: string)
```

**Caractéristiques :**
- ✅ Aucun champ marqué avec `#`
- ✅ ID généré à partir du hash de toutes les valeurs
- ✅ ID de forme : `TypeName~<hash-16-chars-hex>`

### Types Autorisés pour les Clés Primaires

**✅ Types Primitifs AUTORISÉS :**
- `string`
- `number` (int ou float)
- `bool`

**❌ Types NON AUTORISÉS :**
- `object`
- Types complexes/composites

```tsd
// ❌ ERREUR : 'data' est de type object
type Document(#data: object, title: string)
// Erreur : "primary key field must be a primitive type"

// ✅ CORRECT : Utiliser un champ primitif
type Document(#documentId: string, title: string)
```

---

## Génération d'IDs

### Algorithme de Génération

```
SI le type a une ou plusieurs clés primaires
    ALORS génération par clé primaire
    SINON génération par hash
```

### Génération par Clé Primaire

**Fonction :** `generateIDFromPrimaryKey(fact, typeDef)`

**Algorithme :**

1. Récupérer tous les champs marqués comme clé primaire (dans l'ordre de définition)
2. Pour chaque champ de clé primaire :
   - Extraire la valeur du fait
   - Convertir la valeur en string
   - Échapper les caractères spéciaux
3. Construire l'ID : `TypeName~value1_value2_..._valueN`

**Exemple :**

```tsd
type Product(#category: string, #name: string, price: number)
assert Product(category: "Electronics", name: "Laptop", price: 1200)

// Étapes :
// 1. Clés primaires : ["category", "name"]
// 2. Valeurs : ["Electronics", "Laptop"]
// 3. Échappement : ["Electronics", "Laptop"] (pas de caractères spéciaux)
// 4. ID : "Product~Electronics_Laptop"
```

### Génération par Hash

**Fonction :** `generateIDFromHash(fact, typeDef)`

**Algorithme :**

1. Récupérer TOUS les champs du type (dans l'ordre de définition)
2. Pour chaque champ avec une valeur non-nulle :
   - Créer la chaîne `fieldName=value`
3. Concaténer toutes les chaînes avec `|` comme séparateur
4. Calculer le hash MD5 de la chaîne concaténée
5. Tronquer le hash à 16 caractères hexadécimaux
6. Construire l'ID : `TypeName~<hash>`

**Exemple :**

```tsd
type LogEvent(timestamp: number, level: string, message: string)
assert LogEvent(timestamp: 1704067200, level: "ERROR", message: "Connection failed")

// Étapes :
// 1. Chaînes : ["timestamp=1704067200", "level=ERROR", "message=Connection failed"]
// 2. Concaténation : "timestamp=1704067200|level=ERROR|message=Connection failed"
// 3. Hash MD5 : "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
// 4. Tronqué : "a1b2c3d4e5f6g7h8"
// 5. ID : "LogEvent~a1b2c3d4e5f6g7h8"
```

### Déterminisme du Hash

Le hash est **déterministe** : les mêmes valeurs produisent toujours le même hash.

```tsd
type Event(timestamp: number, data: string)

// Premier assert
assert Event(timestamp: 100, data: "test")
// ID : "Event~a1b2c3d4e5f6g7h8"

// Deuxième assert avec les MÊMES valeurs
assert Event(timestamp: 100, data: "test")
// ID : "Event~a1b2c3d4e5f6g7h8" (identique)

// Troisième assert avec des valeurs DIFFÉRENTES
assert Event(timestamp: 200, data: "test")
// ID : "Event~f8e7d6c5b4a39281" (différent)
```

### Conversion de Valeurs en String

**Fonction :** `valueToString(value)`

| Type Go   | Conversion                                    | Exemple                  |
|-----------|-----------------------------------------------|--------------------------|
| `string`  | Aucune conversion                             | `"Alice"` → `"Alice"`    |
| `int`     | `strconv.Itoa()`                              | `123` → `"123"`          |
| `int64`   | `strconv.FormatInt(v, 10)`                    | `456` → `"456"`          |
| `float64` | `strconv.FormatFloat(v, 'f', -1, 64)`         | `12.5` → `"12.5"`        |
| `bool`    | `strconv.FormatBool()`                        | `true` → `"true"`        |
| Autre     | `fmt.Sprintf("%v")`                           | `<custom>` → `"<custom>"`|

**Note sur les floats :** Le format `'f'` avec précision `-1` garantit un format déterministe sans notation scientifique.

```go
// Exemples :
12.5     → "12.5"
100.0    → "100"
0.00001  → "0.00001"
```

---

## Format des IDs

### Structure Générale

Tous les IDs suivent le pattern :

```
TypeName~<valeurs>
```

Où `~` (tilde) est le **séparateur type/valeur**.

### Avec Clé Primaire Simple

**Format :** `TypeName~valeur`

**Exemples :**

| Type | Clé Primaire | Valeur | ID Généré |
|------|--------------|--------|-----------|
| User | `#username: string` | `"alice"` | `User~alice` |
| Product | `#sku: string` | `"LAPTOP-001"` | `Product~LAPTOP-001` |
| Country | `#code: string` | `"FR"` | `Country~FR` |
| Student | `#studentId: number` | `2024001` | `Student~2024001` |
| Active | `#isActive: bool` | `true` | `Active~true` |

### Avec Clé Primaire Composite

**Format :** `TypeName~valeur1_valeur2_..._valeurN`

Le séparateur entre valeurs est `_` (underscore).

**Exemples :**

| Type | Clés Primaires | Valeurs | ID Généré |
|------|----------------|---------|-----------|
| Product | `#category: string`<br>`#name: string` | `"Electronics"`<br>`"Laptop"` | `Product~Electronics_Laptop` |
| Order | `#year: number`<br>`#orderNum: number` | `2024`<br>`1001` | `Order~2024_1001` |
| Location | `#country: string`<br>`#city: string` | `"France"`<br>`"Paris"` | `Location~France_Paris` |
| Enrollment | `#studentId: string`<br>`#courseId: string` | `"S2024001"`<br>`"CS101"` | `Enrollment~S2024001_CS101` |

**Ordre des Valeurs :** L'ordre dans l'ID correspond à l'ordre de déclaration des champs dans le type.

```tsd
// Ordre : category puis name
type Product(#category: string, #name: string, price: number)
assert Product(category: "Books", name: "Go Programming", price: 50)
// ID : "Product~Books_Go Programming"

// ⚠️ Si on inverse l'ordre dans la définition, l'ID change
type ProductV2(#name: string, #category: string, price: number)
assert ProductV2(name: "Go Programming", category: "Books", price: 50)
// ID : "ProductV2~Go Programming_Books"  (ordre inversé)
```

### Sans Clé Primaire (Hash)

**Format :** `TypeName~<hash-16-chars-hex>`

**Exemples :**

```tsd
type LogEvent(timestamp: number, level: string, message: string)
assert LogEvent(timestamp: 1704067200, level: "INFO", message: "Started")
// ID : "LogEvent~a1b2c3d4e5f6g7h8"

type Metric(value: number, unit: string)
assert Metric(value: 42, unit: "ms")
// ID : "Metric~fedcba9876543210"
```

**Caractéristiques du Hash :**
- ✅ 16 caractères hexadécimaux (0-9, a-f)
- ✅ MD5 tronqué
- ✅ Déterministe
- ✅ Basé sur TOUTES les valeurs du fait

---

## Échappement de Caractères

### Pourquoi l'Échappement ?

Certains caractères ont une signification spéciale dans le format d'ID et doivent être échappés.

### Caractères Échappés

| Caractère | Encodage | Raison |
|-----------|----------|--------|
| `~` (tilde) | `%7E` | Séparateur type/valeur |
| `_` (underscore) | `%5F` | Séparateur de valeurs composites |
| `%` (percent) | `%25` | Caractère d'échappement lui-même |
| ` ` (espace) | `%20` | Caractère spécial |
| `/` (slash) | `%2F` | Caractère spécial |

### Fonction d'Échappement

```go
func escapeIDValue(value string) string {
    value = strings.ReplaceAll(value, "%", "%25")   // En premier
    value = strings.ReplaceAll(value, "~", "%7E")
    value = strings.ReplaceAll(value, "_", "%5F")
    value = strings.ReplaceAll(value, " ", "%20")
    return value
}
```

**Ordre Important :** Le `%` est échappé en premier pour éviter de double-échapper.

### Exemples d'Échappement

```tsd
type User(#username: string, email: string)

// Pas de caractères spéciaux
assert User(username: "alice", email: "alice@example.com")
// ID : "User~alice"

// Avec underscore
assert User(username: "john_doe", email: "john@example.com")
// ID : "User~john%5Fdoe"

// Avec tilde
assert User(username: "user~admin", email: "admin@example.com")
// ID : "User~user%7Eadmin"

// Avec espace
assert User(username: "John Doe", email: "john@example.com")
// ID : "User~John%20Doe"

// Avec slash
assert User(username: "admin/root", email: "root@example.com")
// ID : "User~admin%2Froot"

// Avec plusieurs caractères spéciaux
assert User(username: "test_user~123", email: "test@example.com")
// ID : "User~test%5Fuser%7E123"
```

### Clé Primaire Composite avec Échappement

```tsd
type File(#directory: string, #filename: string)

assert File(directory: "/home/user", filename: "my_file.txt")
// directory: "/home/user" → "%2Fhome%2Fuser"
// filename: "my_file.txt" → "my%5Ffile.txt"
// ID : "File~%2Fhome%2Fuser_my%5Ffile.txt"
```

**Note :** Le séparateur `_` entre les valeurs n'est PAS échappé (c'est le séparateur structurel).

### Déséchappement

La fonction inverse `unescapeIDValue()` permet de récupérer les valeurs originales :

```go
func unescapeIDValue(value string) string {
    value = strings.ReplaceAll(value, "%20", " ")
    value = strings.ReplaceAll(value, "%5F", "_")
    value = strings.ReplaceAll(value, "%7E", "~")
    value = strings.ReplaceAll(value, "%25", "%")   // En dernier
    return value
}
```

---

## Cas Particuliers

### Valeurs Nulles dans les Clés Primaires

**Règle :** Les champs de clé primaire **NE PEUVENT PAS** être nuls ou absents.

```tsd
type Product(#sku: string, name: string)

// ❌ ERREUR : Champ de clé primaire manquant
assert Product(name: "Laptop")
// Erreur : "champ de clé primaire 'sku' manquant dans le fait"
```

### Valeurs Vides (String)

Une string vide `""` est une valeur valide pour une clé primaire :

```tsd
type Tag(#label: string)

// ✅ Valide : string vide
assert Tag(label: "")
// ID : "Tag~" (TypeName~ suivi de rien)
```

### Valeurs Booléennes

Les booléens sont convertis en `"true"` ou `"false"` :

```tsd
type Flag(#isActive: bool, description: string)

assert Flag(isActive: true, description: "Active flag")
// ID : "Flag~true"

assert Flag(isActive: false, description: "Inactive flag")
// ID : "Flag~false"
```

### Valeurs Numériques

Les nombres sont convertis en string sans notation scientifique :

```tsd
type Measurement(#sensorId: number, value: number)

assert Measurement(sensorId: 123, value: 45.6)
// ID : "Measurement~123"

assert Measurement(sensorId: 0, value: 0.0)
// ID : "Measurement~0"

// Float dans la clé primaire
type Data(#temperature: number)

assert Data(temperature: 12.5)
// ID : "Data~12.5"

assert Data(temperature: 100.0)
// ID : "Data~100"  (pas "100.0")
```

### Types Sans Champs

Un type sans champs (ou avec uniquement des clés primaires) :

```tsd
type EmptyType()
assert EmptyType()
// ID : "EmptyType~<hash>" (hash de chaîne vide)

type SingletonType(#marker: bool)
assert SingletonType(marker: true)
// ID : "SingletonType~true"
```

### Collisions d'IDs

Avec des clés primaires, deux faits avec les mêmes valeurs de clé primaire ont le même ID :

```tsd
type User(#username: string, email: string)

assert User(username: "alice", email: "alice@example.com")
// ID : "User~alice"

assert User(username: "alice", email: "alice@newdomain.com")
// ID : "User~alice" (MÊME ID)
// ⚠️ Dans RETE, cela mettra à jour le fait existant ou sera considéré comme une duplication
```

**Comportement RETE :** Si un fait avec le même ID existe déjà, le comportement dépend de l'implémentation (mise à jour ou rejet).

Avec hash, les collisions sont théoriquement possibles mais extrêmement rares (MD5 sur 16 chars hex = 2^64 possibilités).

---

## Utilisation dans les Règles

### Accès au Champ `id`

Le champ `id` est accessible comme n'importe quel autre champ :

```tsd
type User(#username: string, age: number)

rule logUser : {u: User} / true
    ==> print("User ID: " + u.id)
    // Affiche : "User ID: User~alice"
```

### Comparaison d'IDs

```tsd
type Person(#name: string, age: number)

// Vérifier qu'un ID n'est pas vide
rule hasId : {p: Person} / p.id != ""
    ==> print("Person has ID")

// Comparer avec un ID spécifique
rule isAlice : {p: Person} / p.id == "Person~Alice"
    ==> print("Found Alice")

// ⚠️ Attention : Préférer comparer la clé primaire directement
rule isAliceBetter : {p: Person} / p.name == "Alice"
    ==> print("Found Alice")
```

### Jointures sur IDs

```tsd
type User(#username: string, email: string)
type Order(#orderId: string, username: string, total: number)

// Jointure via le username (clé primaire de User)
rule userOrders : {u: User, o: Order} / u.username == o.username
    ==> process(u.id, o.id)
    // u.id : "User~alice"
    // o.id : "Order~ORD-001"
```

### Pattern Matching sur IDs

```tsd
type Product(#category: string, #name: string, price: number)

// ⚠️ Difficile : utiliser des regex/patterns sur les IDs
rule electronicProducts : {p: Product} / p.id matches "Product~Electronics_.*"
    ==> print("Electronic product")

// ✅ Mieux : utiliser la clé primaire directement
rule electronicProductsBetter : {p: Product} / p.category == "Electronics"
    ==> print("Electronic product")
```

### Logging et Debugging

Le champ `id` est très utile pour le logging :

```tsd
type Event(timestamp: number, data: string)

rule logEvent : {e: Event} / true
    ==> log("Processing event: " + e.id + " at " + e.timestamp)
    // Log : "Processing event: Event~a1b2c3d4e5f6g7h8 at 1704067200"
```

---

## Erreurs Courantes

### Erreur 1 : Définir un Champ `id`

```tsd
// ❌ ERREUR
type Person(id: string, name: string)

// Erreur : "field 'id' is reserved and cannot be defined manually"
```

**Solution :**
```tsd
// ✅ Option 1 : Renommer le champ
type Person(personId: string, name: string)

// ✅ Option 2 : Déclarer 'id' comme clé primaire
type Person(#id: string, name: string)
```

### Erreur 2 : Assigner Manuellement un ID

```tsd
type Person(name: string)

// ❌ ERREUR
assert Person(id: "custom-id", name: "Alice")

// Erreur : "le champ 'id' ne peut pas être défini manuellement"
```

**Solution :**
```tsd
// ✅ Laisser le système générer l'ID
assert Person(name: "Alice")
// ID auto-généré : "Person~<hash>"

// ✅ Ou utiliser une vraie clé primaire
type Person(#personId: string, name: string)
assert Person(personId: "custom-id", name: "Alice")
```

### Erreur 3 : Clé Primaire Manquante

```tsd
type Product(#sku: string, name: string)

// ❌ ERREUR : 'sku' manquant
assert Product(name: "Laptop")

// Erreur : "champ de clé primaire 'sku' manquant dans le fait"
```

**Solution :**
```tsd
// ✅ Fournir TOUS les champs de clé primaire
assert Product(sku: "ABC123", name: "Laptop")
```

### Erreur 4 : Comparer `id` avec un Number

```tsd
type Item(#itemId: number, name: string)

// ❌ ERREUR : 'id' est string, pas number
rule checkItem : {i: Item} / i.id > 0
    ==> print("found")

// Erreur : "type incompatibility in comparison: string vs number"
```

**Solution :**
```tsd
// ✅ Option 1 : Comparer la clé primaire directement
rule checkItem : {i: Item} / i.itemId > 0
    ==> print("found")

// ✅ Option 2 : Vérifier que l'id existe
rule checkItem : {i: Item} / i.id != ""
    ==> print("found")
```

### Erreur 5 : Clé Primaire de Type Non-Primitif

```tsd
// ❌ ERREUR : 'data' est de type object
type Document(#data: object, title: string)

// Erreur : "primary key field must be a primitive type"
```

**Solution :**
```tsd
// ✅ Utiliser un champ primitif comme clé primaire
type Document(#documentId: string, title: string, data: object)
```

---

## Exemples Complets

### Exemple 1 : Gestion d'Utilisateurs

```tsd
// Définition du type avec clé primaire simple
type User(#username: string, email: string, role: string)

// Assertions
assert User(username: "alice", email: "alice@example.com", role: "admin")
assert User(username: "bob", email: "bob@example.com", role: "user")
assert User(username: "charlie", email: "charlie@example.com", role: "user")

// IDs générés :
// "User~alice"
// "User~bob"
// "User~charlie"

// Règle utilisant l'ID
rule logAdmins : {u: User} / u.role == "admin"
    ==> log("Admin user: " + u.id + " (" + u.username + ")")
    // Log : "Admin user: User~alice (alice)"
```

### Exemple 2 : Catalogue de Produits avec Clé Composite

```tsd
// Clé primaire composite : catégorie + nom
type Product(#category: string, #name: string, price: number, stock: number)

// Assertions
assert Product(category: "Electronics", name: "Laptop", price: 1200, stock: 5)
assert Product(category: "Electronics", name: "Mouse", price: 25, stock: 100)
assert Product(category: "Books", name: "Go Programming", price: 50, stock: 20)

// IDs générés :
// "Product~Electronics_Laptop"
// "Product~Electronics_Mouse"
// "Product~Books_Go Programming"  (espace conservé)

// Règle : Stock bas
rule lowStock : {p: Product} / p.stock < 10
    ==> alert("Low stock for product: " + p.id)
    // Alert : "Low stock for product: Product~Electronics_Laptop"
```

### Exemple 3 : Événements de Log sans Clé Primaire

```tsd
// Pas de clé primaire → génération par hash
type LogEvent(timestamp: number, level: string, message: string, source: string)

// Assertions
assert LogEvent(
    timestamp: 1704067200,
    level: "INFO",
    message: "Application started",
    source: "main"
)

assert LogEvent(
    timestamp: 1704067201,
    level: "ERROR",
    message: "Connection failed",
    source: "network"
)

// IDs générés (exemples) :
// "LogEvent~a1b2c3d4e5f6g7h8"
// "LogEvent~f8e7d6c5b4a39281"

// Règle : Filtrer les erreurs
rule errorLogs : {log: LogEvent} / log.level == "ERROR"
    ==> process("Error detected: " + log.id + " - " + log.message)
    // Process : "Error detected: LogEvent~f8e7d6c5b4a39281 - Connection failed"
```

### Exemple 4 : Relations Entre Types

```tsd
// Type User avec clé primaire
type User(#username: string, email: string, department: string)

// Type Task avec clé primaire composite
type Task(#taskId: string, assignee: string, status: string, priority: number)

// Assertions Users
assert User(username: "alice", email: "alice@example.com", department: "Engineering")
assert User(username: "bob", email: "bob@example.com", department: "Sales")

// Assertions Tasks (assignee référence username de User)
assert Task(taskId: "TASK-001", assignee: "alice", status: "in-progress", priority: 1)
assert Task(taskId: "TASK-002", assignee: "alice", status: "todo", priority: 2)
assert Task(taskId: "TASK-003", assignee: "bob", status: "done", priority: 1)

// IDs :
// Users : "User~alice", "User~bob"
// Tasks : "Task~TASK-001", "Task~TASK-002", "Task~TASK-003"

// Règle : Tâches prioritaires d'Alice
rule aliceHighPriority : {u: User, t: Task} /
    u.username == "alice" AND
    t.assignee == u.username AND
    t.priority == 1
    ==> notify(u.id + " has high priority task " + t.id)
    // Notify : "User~alice has high priority task Task~TASK-001"
```

### Exemple 5 : Échappement de Caractères Spéciaux

```tsd
type File(#path: string, size: number)

// Divers chemins avec caractères spéciaux
assert File(path: "/home/user/file.txt", size: 1024)
assert File(path: "C:\\Users\\Admin\\file.txt", size: 2048)
assert File(path: "my_document~v2.pdf", size: 512)
assert File(path: "folder/sub folder/file.txt", size: 256)

// IDs générés :
// "File~%2Fhome%2Fuser%2Ffile.txt"
// "File~C%3A%5CUsers%5CAdmin%5Cfile.txt"
// "File~my%5Fdocument%7Ev2.pdf"
// "File~folder%2Fsub%20folder%2Ffile.txt"

// Règle : Logs tous les fichiers
rule logFiles : {f: File} / true
    ==> log("File: " + f.id + " (" + f.size + " bytes)")
    // Les IDs sont échappés dans les logs
```

### Exemple 6 : Type avec `id` comme Clé Primaire

```tsd
// Exception : 'id' déclaré explicitement comme clé primaire
type Record(#id: string, data: string, timestamp: number)

// On DOIT fournir 'id' dans chaque assertion
assert Record(id: "rec-001", data: "Sample data", timestamp: 1704067200)
assert Record(id: "rec-002", data: "More data", timestamp: 1704067201)

// IDs générés :
// "Record~rec-001"
// "Record~rec-002"

// Règle
rule processRecords : {r: Record} / true
    ==> process(r.id + ": " + r.data)
    // Process : "Record~rec-001: Sample data"
```

---

## Résumé des Règles

### ✅ À FAIRE

1. **Toujours déclarer les clés primaires naturelles** avec `#`
2. **Utiliser des types primitifs** pour les clés primaires
3. **Fournir toutes les valeurs** de clé primaire dans les assertions
4. **Accéder au champ `id`** pour logging et traçabilité
5. **Comparer les clés primaires directement** plutôt que les IDs
6. **Documenter** pourquoi un type a ou n'a pas de clé primaire

### ❌ À ÉVITER

1. **Ne jamais définir un champ nommé `id`** (sauf comme clé primaire)
2. **Ne jamais assigner manuellement `id`** dans les assertions
3. **Ne pas comparer `id` avec des types non-string**
4. **Ne pas utiliser de types complexes** comme clés primaires
5. **Ne pas compter sur l'ordre des hash** (ils sont non prédictibles)
6. **Ne pas oublier d'échapper** les caractères spéciaux si vous parsez des IDs

---

## Références

- **Code Source :** `tsd/constraint/id_generator.go`
- **Tests :** `tsd/constraint/id_generator_test.go`
- **Constantes :** `tsd/constraint/constraint_constants.go`
- **Validation :** `tsd/constraint/primary_key_validation.go`
- **Documentation :** `tsd/docs/MIGRATION_IDS.md`
- **README :** `tsd/README.md`

---

## Constantes du Système

```go
// Séparateur entre type et valeurs
const IDSeparatorType = "~"

// Séparateur entre valeurs de clé composite
const IDSeparatorValue = "_"

// Longueur du hash (16 caractères hex)
const IDHashLength = 16

// Nom du champ ID réservé
const FieldNameID = "id"

// Nom du champ type RETE réservé
const FieldNameReteType = "reteType"
```

---

**Version du Document :** 1.0  
**Dernière Mise à Jour :** 2025-01-XX  
**Auteur :** TSD Team