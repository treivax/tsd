# Identifiants Internes - Documentation Complète

## Vue d'Ensemble

Dans TSD, chaque fait possède un **identifiant interne unique** (`_id_`) qui est :

1. **Généré automatiquement** - Jamais défini manuellement
2. **Déterministe** - Basé sur les clés primaires ou hash du contenu
3. **Caché** - Jamais accessible dans les expressions TSD
4. **Interne** - Utilisé uniquement par le moteur RETE

---

## ⚠️ Règle Fondamentale

Le champ `_id_` est **strictement réservé au système** :

### ❌ INTERDIT

- Définir `_id_` dans une définition de type
- Assigner `_id_` dans un fait
- Accéder à `_id_` dans une expression ou règle
- Comparer `_id_` explicitement dans les conditions

### ✅ PERMIS

- Les IDs sont générés automatiquement par le système
- Les comparaisons de faits utilisent les IDs en interne (transparence)
- Les références entre faits sont résolues automatiquement
- Le moteur RETE utilise `_id_` pour identifier les faits

**Exemple du fonctionnement** :

```tsd
type User(#username: string, email: string, age: number)

// Définir un utilisateur
alice = User("alice", "alice@example.com", 30)
// En interne, le système génère: _id_ = "User~alice"
// Mais vous ne voyez JAMAIS ce champ dans vos expressions

// ❌ INTERDIT - Erreur de parsing
{u: User} / u._id_ == "User~alice" ==> Log("Found")

// ✅ CORRECT - Utiliser les champs métier
{u: User} / u.username == "alice" ==> Log("Found")
```

---

## 📋 Table des Matières

1. [Génération Automatique](#génération-automatique)
2. [Format des IDs](#format-des-ids)
3. [Clés Primaires](#clés-primaires)
4. [Caractères Spéciaux](#caractères-spéciaux)
5. [Utilisation Interne](#utilisation-interne)
6. [Déterminisme](#déterminisme)
7. [Bonnes Pratiques](#bonnes-pratiques)
8. [Exemples Complets](#exemples-complets)

---

## Génération Automatique

### Avec Clés Primaires

Les clés primaires (préfixées par `#`) déterminent l'ID généré :

```tsd
type User(#username: string, email: string, age: number)

// Définir un utilisateur
alice = User("alice", "alice@example.com", 30)
// ID généré en interne: "User~alice"
// Mais _id_ n'est PAS accessible dans vos expressions
```

**Format Interne** : `TypeName~valeur1_valeur2_...`

L'ID est généré mais **caché**. Vous n'y accédez jamais directement.

#### Clé Primaire Simple

Une seule clé primaire :

```tsd
type Product(#sku: string, name: string, price: number)

laptop = Product("LAPTOP-001", "Gaming Laptop", 1200.00)
// ID interne: "Product~LAPTOP-001"
```

#### Clé Primaire Composite

Plusieurs champs forment la clé :

```tsd
type OrderLine(#orderId: string, #productId: string, quantity: number)

line1 = OrderLine("ORD-001", "PROD-123", 2)
// ID interne: "OrderLine~ORD-001_PROD-123"
```

L'ordre des champs dans la définition détermine l'ordre dans l'ID.

---

### Sans Clé Primaire (Hash)

Si aucun champ n'est marqué `#`, un **hash déterministe** est généré :

```tsd
type LogEvent(timestamp: number, level: string, message: string)

LogEvent(1704067200, "INFO", "Application started")
// ID interne: "LogEvent~a1b2c3d4e5f6g7h8" (hash de 16 caractères)
```

Le hash est calculé à partir de **tous les champs** du fait, garantissant :
- ✅ Déterminisme : mêmes valeurs → même hash
- ✅ Unicité : valeurs différentes → hash différent (avec très haute probabilité)
- ✅ Performance : hash rapide et efficace

**Quand utiliser** :
- Événements temporels (logs, audit)
- Faits éphémères
- Pas d'identifiant naturel évident

---

## Format des IDs

### Structure Générale

```
TypeName~valeur1_valeur2_..._valeurN
```

**Composants** :
- `TypeName` : Nom du type du fait
- `~` : Séparateur entre type et valeurs
- `valeur1, valeur2, ...` : Valeurs des clés primaires (encodées)
- `_` : Séparateur entre valeurs de clés composites

### Exemples de Format

| Type | Clés Primaires | Exemple d'ID Interne |
|------|----------------|----------------------|
| User(#username) | Simple | `User~alice` |
| Product(#category, #sku) | Composite | `Product~Electronics_LAP123` |
| LogEvent(...) | Aucune (hash) | `LogEvent~a1b2c3d4e5f6g7h8` |
| Order(#year, #num) | Composite | `Order~2024_ORD001` |

**Note** : Ces IDs sont internes et transparents pour vous.

---

## Clés Primaires

### Définir une Clé Primaire

Utilisez le préfixe `#` devant les champs servant d'identifiant :

```tsd
type User(#username: string, email: string, role: string)
```

### Règles

1. **Un ou plusieurs champs** : Clé simple ou composite
2. **Types compatibles** : string, number, boolean
3. **Ordre important** : Détermine l'ordre dans l'ID
4. **Valeurs non-null** : Les clés primaires ne peuvent pas être null

### Choix de Clés Primaires

#### Critères de Sélection

✅ **Bonne clé primaire** :
- Unique pour chaque instance
- Stable (ne change pas)
- Simple si possible
- Signification métier

❌ **Mauvaise clé primaire** :
- Valeur changeante (ex: email qui peut changer)
- Calculée ou dérivée
- Trop complexe

#### Exemples

```tsd
// ✅ BON - username est stable et unique
type User(#username: string, email: string)

// ⚠️ ATTENTION - email peut changer
type User(username: string, #email: string)

// ✅ BON - Clé composite pour relation N-N
type Enrollment(#studentId: string, #courseId: string, grade: string)

// ✅ BON - SKU est un identifiant métier standard
type Product(#sku: string, name: string, price: number)

// ❌ ÉVITER - timestamp seul peut avoir des collisions
type Event(#timestamp: number, description: string)

// ✅ MIEUX - Pas de clé → hash de tous les champs
type Event(timestamp: number, description: string, source: string)
```

---

## Caractères Spéciaux

Les caractères spéciaux dans les valeurs de clés primaires sont **automatiquement encodés** :

### Table d'Encodage

| Caractère | Encodage | Raison |
|-----------|----------|--------|
| `~` | `%7E` | Séparateur type/valeur |
| `_` | `%5F` | Séparateur clés composites |
| `%` | `%25` | Caractère d'échappement |
| ` ` (espace) | `%20` | Lisibilité |
| `/` | `%2F` | Chemins |
| `\` | `%5C` | Chemins Windows |

### Exemple

```tsd
type File(#path: string, size: number)

file1 = File("/home/user~backup_v1", 1024)
// ID interne: "File~%2Fhome%2Fuser%7Ebackup%5Fv1"
// L'encodage est automatique et transparent
```

**Important** : Vous n'avez pas à vous soucier de l'encodage. C'est géré automatiquement.

---

## Utilisation Interne

### Comparaisons de Faits

Les comparaisons de faits utilisent automatiquement les IDs internes :

```tsd
type User(#username: string, email: string)
type Login(user: User, #sessionId: string, timestamp: number)

alice = User("alice", "alice@example.com")
session1 = Login(alice, "SES-001", 1704067200)

// Règle : Comparer les faits
{u: User, l: Login} / l.user == u ==> 
    Log("Login for user: " + u.username)
```

**Fonctionnement interne (transparent pour vous)** :
1. `alice` est créé avec ID interne : `"User~alice"`
2. `session1.user` stocke l'ID : `"User~alice"`
3. Dans la règle, `l.user == u` compare : `"User~alice" == "User~alice"` → `true`

Vous écrivez simplement `l.user == u` et le système gère le reste.

---

### Références entre Faits

Les champs de type fait stockent l'ID interne (transparent) :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)

alice = User("alice", "alice@example.com")
order1 = Order(alice, "ORD-001", 150.00)

// En interne:
// alice._id_ = "User~alice"
// order1.customer = "User~alice" (référence)
// order1._id_ = "Order~ORD-001"
```

**Transparence** : Vous manipulez `alice` et `order1` naturellement, sans voir les IDs.

---

### Résolution de Références

Le moteur RETE résout automatiquement les références :

```tsd
type User(#username: string)
type Post(author: User, #postId: string, title: string)
type Comment(post: Post, author: User, #commentId: string, text: string)

alice = User("alice")
bob = User("bob")

post1 = Post(alice, "P1", "Hello World")
Comment(post1, bob, "C1", "Nice post!")

// Règle avec chaîne de références
{p: Post, c: Comment, u: User} / 
    c.post == p && c.author == u ==> 
    Log(u.username + " commented on post by " + p.author.username)
```

Le moteur navigue dans la chaîne : `Comment → Post → User`

---

## Déterminisme

Les IDs internes sont **toujours les mêmes** pour les mêmes valeurs :

```tsd
type User(#username: string, age: number)

// Ces deux définitions génèrent le MÊME ID interne
alice1 = User("alice", 30)
alice2 = User("alice", 30)
// Les deux ont l'ID interne: "User~alice"
```

⚠️ **Attention** : Cela signifie que le deuxième fait **remplace** le premier dans le réseau RETE (même identité).

### Implications

```tsd
type Product(#sku: string, name: string, price: number)

// Première assertion
Product("LAP-001", "Old Laptop", 1000.00)

// Deuxième assertion avec même SKU
Product("LAP-001", "New Laptop", 1200.00)

// Résultat: Le premier est remplacé par le second
// Car même ID: "Product~LAP-001"
```

**Utilité** : Permet les mises à jour naturelles (upsert).

---

## Bonnes Pratiques

### 1. Choisir des Clés Primaires Naturelles

```tsd
// ✅ BON - Clé naturelle et stable
type User(#username: string, email: string, role: string)

// ❌ ÉVITER - Pas de clé (hash moins prévisible)
type User(username: string, email: string, role: string)
```

**Avantage** : IDs prévisibles facilitent le débogage.

---

### 2. Clés Stables

```tsd
// ✅ BON - username change rarement
type User(#username: string, email: string)

// ⚠️ ATTENTION - email peut changer
type User(username: string, #email: string)
```

Si l'email change, l'ID change → nouveau fait créé.

---

### 3. Clés Composites pour Relations N-N

```tsd
// ✅ BON - Clé composite pour table de jonction
type Enrollment(#studentId: string, #courseId: string, grade: string, year: number)
```

Garantit unicité de la paire (étudiant, cours).

---

### 4. Hash pour Événements Temporels

```tsd
// ✅ BON - Événements sans identifiant naturel
type AuditLog(timestamp: number, userId: string, action: string, details: string)
```

Pas de clé primaire → hash automatique.

---

### 5. Ne Jamais Accéder à `_id_`

```tsd
// ❌ INTERDIT - Erreur de parsing
{u: User} / u._id_ == "User~alice" ==> Log("Found")

// ✅ CORRECT - Utiliser les champs métier
{u: User} / u.username == "alice" ==> Log("Found")
```

---

## Exemples Complets

### Exemple 1 : Système de Blog

```tsd
type User(#username: string, email: string, bio: string)
type Post(author: User, #postId: string, title: string, content: string)
type Comment(post: Post, author: User, #commentId: string, text: string)

// Créer des utilisateurs
alice = User("alice", "alice@example.com", "Software engineer")
bob = User("bob", "bob@example.com", "Tech enthusiast")

// IDs internes générés:
// alice._id_ = "User~alice"
// bob._id_ = "User~bob"

// Créer des posts
post1 = Post(alice, "POST-001", "Introduction to TSD", "Welcome to TSD!")
post2 = Post(bob, "POST-002", "Advanced Features", "Deep dive into RETE...")

// IDs internes:
// post1._id_ = "Post~POST-001"
// post1.author = "User~alice" (référence)

// Créer des commentaires
Comment(post1, bob, "COM-001", "Great post!")
Comment(post2, alice, "COM-002", "Very helpful, thanks!")

// IDs internes:
// Comment._id_ = "Comment~COM-001"
// Comment.post = "Post~POST-001" (référence)
// Comment.author = "User~bob" (référence)

// Règles utilisant les comparaisons automatiques
rule userComments : {p: Post, c: Comment, u: User} / 
    c.post == p && c.author == u && u.username == "alice" 
    ==> Log("Alice commented: " + c.text)

rule postsByAuthor : {p: Post, u: User} /
    p.author == u && u.username == "bob"
    ==> Log("Post by Bob: " + p.title)
```

**Transparence** : Vous ne voyez jamais les IDs, tout fonctionne naturellement.

---

### Exemple 2 : E-Commerce

```tsd
type Customer(#customerId: string, name: string, vip: boolean)
type Product(#sku: string, name: string, price: number)
type Order(customer: Customer, #orderNumber: string, date: string, total: number)
type OrderLine(order: Order, product: Product, quantity: number, subtotal: number)

// Créer des clients
cust1 = Customer("CUST-001", "Alice Johnson", true)
cust2 = Customer("CUST-002", "Bob Smith", false)

// IDs internes: "Customer~CUST-001", "Customer~CUST-002"

// Créer des produits
prod1 = Product("LAPTOP-001", "Gaming Laptop", 1200.00)
prod2 = Product("MOUSE-001", "Wireless Mouse", 25.00)
prod3 = Product("KEYBOARD-001", "Mechanical Keyboard", 150.00)

// IDs internes: "Product~LAPTOP-001", etc.

// Créer des commandes
order1 = Order(cust1, "ORD-001", "2024-12-19", 1250.00)
order2 = Order(cust2, "ORD-002", "2024-12-19", 300.00)

// IDs internes: "Order~ORD-001", "Order~ORD-002"
// order1.customer = "Customer~CUST-001" (référence)

// Créer des lignes de commande
OrderLine(order1, prod1, 1, 1200.00)
OrderLine(order1, prod2, 2, 50.00)
OrderLine(order2, prod3, 2, 300.00)

// Pas de clé primaire sur OrderLine → hash automatique
// IDs internes: "OrderLine~<hash1>", "OrderLine~<hash2>", etc.

// Règles : Analyser les commandes VIP
rule vipOrders : {c: Customer, o: Order, ol: OrderLine, p: Product} / 
    o.customer == c && ol.order == o && ol.product == p && c.vip == true 
    ==> Log("VIP order " + o.orderNumber + " contains " + ol.quantity + "x " + p.name)

// Règle : Produits populaires
rule popularProducts : {ol1: OrderLine, ol2: OrderLine, p: Product} /
    ol1.product == p && ol2.product == p && ol1 != ol2
    ==> Log("Product " + p.name + " ordered multiple times")
```

---

### Exemple 3 : Système Hiérarchique

```tsd
type Company(#companyId: string, name: string, country: string)
type Department(company: Company, #deptId: string, name: string, budget: number)
type Employee(dept: Department, #empId: string, name: string, salary: number)
type Project(dept: Department, #projectId: string, name: string, deadline: string)

// Créer une entreprise
acme = Company("COMP-001", "ACME Corp", "USA")
// ID interne: "Company~COMP-001"

// Créer des départements
engineering = Department(acme, "DEPT-ENG", "Engineering", 1000000.00)
marketing = Department(acme, "DEPT-MKT", "Marketing", 500000.00)

// IDs internes: "Department~DEPT-ENG", "Department~DEPT-MKT"
// engineering.company = "Company~COMP-001" (référence)

// Créer des employés
alice = Employee(engineering, "EMP-001", "Alice Johnson", 120000.00)
bob = Employee(engineering, "EMP-002", "Bob Smith", 90000.00)
charlie = Employee(marketing, "EMP-003", "Charlie Brown", 85000.00)

// IDs internes: "Employee~EMP-001", etc.
// alice.dept = "Department~DEPT-ENG" (référence)

// Créer des projets
proj1 = Project(engineering, "PROJ-001", "New Platform", "2024-12-31")
proj2 = Project(marketing, "PROJ-002", "Campaign 2024", "2024-06-30")

// Règles : Analyser la structure
rule companyEmployees : {comp: Company, dept: Department, emp: Employee} /
    dept.company == comp && emp.dept == dept
    ==> Log(emp.name + " works at " + comp.name + " in " + dept.name)

rule highSalaries : {emp: Employee, dept: Department, comp: Company} /
    emp.dept == dept && dept.company == comp && emp.salary > 100000.00
    ==> Log("High salary: " + emp.name + " earns " + emp.salary + " at " + comp.name)

rule departmentProjects : {dept: Department, proj: Project, emp: Employee} /
    proj.dept == dept && emp.dept == dept
    ==> Log("Dept " + dept.name + ": " + emp.name + " may work on " + proj.name)
```

---

## Résumé

### Caractéristiques Principales

| Aspect | Comportement |
|--------|--------------|
| **Nom** | `_id_` (caché, interne) |
| **Génération** | Automatique, toujours |
| **Format** | `Type~value` ou `Type~hash` |
| **Accès** | ❌ Jamais dans expressions TSD |
| **Comparaisons** | ✅ Automatiques en interne |
| **Affectation** | ❌ Interdite |
| **Utilisation** | Interne au moteur RETE uniquement |

### Points Clés

1. **Caché** : `_id_` n'est jamais visible dans vos expressions
2. **Automatique** : Généré par le système, pas par vous
3. **Déterministe** : Mêmes valeurs → même ID
4. **Transparent** : Comparaisons de faits fonctionnent automatiquement
5. **Interdit** : Ne jamais essayer d'accéder à `_id_`

### Workflow Typique

```tsd
// 1. Définir types avec clés primaires
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)

// 2. Créer faits avec affectations
alice = User("alice", "alice@example.com")
order1 = Order(alice, "ORD-001", 150.00)
// IDs générés automatiquement en interne (cachés)

// 3. Écrire règles naturellement
{u: User, o: Order} / o.customer == u ==> 
    Log("Order " + o.orderNum + " for " + u.username)
// Comparaisons utilisent IDs internes automatiquement
```

---

## Références

- [Guide Utilisateur - Affectations](user-guide/fact-assignments.md)
- [Guide Utilisateur - Comparaisons](user-guide/fact-comparisons.md)
- [Système de Types](user-guide/type-system.md)
- [Guide de Migration v1.x → v2.0](migration/from-v1.x.md)
- [Architecture - Génération d'IDs](architecture/id-generation.md)

---

**Note** : Cette documentation décrit le système à partir de la version 2.0.

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-19  
**Auteur** : Équipe TSD
