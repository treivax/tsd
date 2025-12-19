# Guide Utilisateur - Affectations de Faits

## Introduction

Les **affectations de faits** permettent de nommer des faits pour les référencer dans d'autres définitions de faits et dans les règles.

Cette fonctionnalité est **nouvelle en v2.0** et change fondamentalement la façon dont vous construisez des relations entre faits.

---

## 📋 Table des Matières

1. [Syntaxe](#syntaxe)
2. [Utilisation](#utilisation)
3. [Règles et Contraintes](#règles-et-contraintes)
4. [Exemples Pratiques](#exemples-pratiques)
5. [Bonnes Pratiques](#bonnes-pratiques)
6. [Utilisation avec Règles](#utilisation-avec-règles)

---

## Syntaxe

### Forme Générale

```tsd
variable = TypeName(valeur1, valeur2, ...)
```

**Composants** :
- `variable` : Nom de la variable (identificateur)
- `=` : Opérateur d'affectation
- `TypeName(...)` : Création d'un fait du type spécifié

### Exemples de Base

```tsd
alice = User("alice", "alice@example.com", 30)
laptop = Product("LAPTOP-001", "Gaming Laptop", 1200.00)
order1 = Order(alice, "ORD-001", 150.00)
```

**Convention de nommage** :
- camelCase pour les variables
- Noms descriptifs et signifiants
- Pas de caractères spéciaux (sauf `_`)

---

## Utilisation

### Référencer dans d'Autres Faits

L'utilisation principale des affectations est de créer des **relations entre faits** :

```tsd
type User(#username: string, email: string)
type Login(user: User, #sessionId: string, timestamp: number)

// Affecter un utilisateur à une variable
alice = User("alice", "alice@example.com")

// Utiliser la variable dans un autre fait
Login(alice, "SES-001", 1704067200)
Login(alice, "SES-002", 1704068000)
```

**Fonctionnement** :
1. `alice` est créé et reçoit un ID interne automatique
2. Les faits `Login` stockent une référence à `alice`
3. La relation est maintenue automatiquement par le système

---

### Réutilisation de Variables

Une fois définie, une variable peut être réutilisée autant de fois que nécessaire :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)
type Payment(order: Order, #paymentId: string, amount: number)

// Définir une fois
alice = User("alice", "alice@example.com")

// Réutiliser plusieurs fois
order1 = Order(alice, "ORD-001", 150.00)
order2 = Order(alice, "ORD-002", 75.00)
order3 = Order(alice, "ORD-003", 200.00)

Payment(order1, "PAY-001", 150.00)
Payment(order2, "PAY-002", 75.00)
```

**Avantages** :
- ✅ Pas de duplication de valeurs
- ✅ Maintenabilité (modifier une seule fois)
- ✅ Lisibilité accrue
- ✅ Cohérence garantie

---

### Chaînes de Références

Les affectations permettent de construire des **hiérarchies de faits** :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)
type Payment(order: Order, #paymentId: string, amount: number, method: string)

// Niveau 1: Utilisateur
alice = User("alice", "alice@example.com")

// Niveau 2: Commande référence utilisateur
order1 = Order(alice, "ORD-001", 150.00)

// Niveau 3: Paiement référence commande (qui référence utilisateur)
Payment(order1, "PAY-001", 150.00, "credit_card")
```

La chaîne : `Payment → Order → User`

Le moteur RETE peut naviguer dans cette chaîne automatiquement.

---

## Règles et Contraintes

### 1. Variables Uniques

Une variable ne peut être définie qu'**une seule fois** dans un programme :

```tsd
// ✅ CORRECT
alice = User("alice", "alice@example.com")
bob = User("bob", "bob@example.com")

// ❌ ERREUR : alice est déjà définie
alice = User("alice2", "alice2@example.com")
```

**Erreur** : `variable 'alice' is already defined`

---

### 2. Variables Avant Utilisation

Une variable doit être **définie avant** d'être utilisée :

```tsd
// ❌ ERREUR : alice n'est pas encore définie
Login(alice, "SES-001", 1704067200)

alice = User("alice", "alice@example.com")
```

**Erreur** : `undefined variable 'alice'`

**✅ CORRECT** :

```tsd
alice = User("alice", "alice@example.com")
Login(alice, "SES-001", 1704067200)
```

**Ordre d'exécution** : Les faits sont traités dans l'ordre de définition.

---

### 3. Types Compatibles

La variable doit être du **type attendu** :

```tsd
type User(#username: string, email: string)
type Product(#sku: string, name: string, price: number)
type Login(user: User, #sessionId: string)

alice = User("alice", "alice@example.com")
laptop = Product("LAP-001", "Laptop", 1200.00)

// ✅ CORRECT : alice est un User
Login(alice, "SES-001")

// ❌ ERREUR : laptop est un Product, pas un User
Login(laptop, "SES-002")
```

**Erreur** : `type mismatch: expected User, got Product`

---

### 4. Portée des Variables

Les variables ont une portée **globale** dans le programme :

```tsd
// Définition globale
alice = User("alice", "alice@example.com")

// Utilisable partout dans le programme
order1 = Order(alice, "ORD-001", 100.00)

rule showOrders : {u: User, o: Order} / o.customer == u ==>
    Log("Order for: " + u.username)
    
// alice est également utilisable ici
```

---

## Exemples Pratiques

### Exemple 1 : Gestion d'Utilisateurs

```tsd
type User(#username: string, email: string, role: string)
type Session(user: User, #sessionId: string, loginTime: number)
type Action(session: Session, action: string, timestamp: number)

// Créer des utilisateurs
admin = User("admin", "admin@example.com", "administrator")
alice = User("alice", "alice@example.com", "user")
bob = User("bob", "bob@example.com", "user")

// Créer des sessions
adminSession = Session(admin, "SES-001", 1704067200)
aliceSession = Session(alice, "SES-002", 1704067260)
bobSession = Session(bob, "SES-003", 1704067320)

// Enregistrer des actions
Action(adminSession, "create_user", 1704067300)
Action(adminSession, "delete_user", 1704067400)
Action(aliceSession, "view_dashboard", 1704067320)
Action(bobSession, "edit_profile", 1704067380)

// Règle : Auditer les actions admin
rule adminActions : {u: User, s: Session, a: Action} /
    s.user == u && a.session == s && u.role == "administrator"
    ==> Log("Admin " + u.username + " performed: " + a.action)
```

**Avantages** :
- Pas de duplication de données utilisateur
- Relations explicites et type-safe
- Facile à maintenir et modifier

---

### Exemple 2 : Hiérarchie Organisationnelle

```tsd
type Company(#companyId: string, name: string, country: string)
type Department(company: Company, #deptId: string, name: string, budget: number)
type Employee(dept: Department, #empId: string, name: string, salary: number)
type Project(dept: Department, #projectId: string, name: string, deadline: string)

// Niveau 1: Entreprise
acme = Company("COMP-001", "ACME Corp", "USA")
techno = Company("COMP-002", "Techno Inc", "Canada")

// Niveau 2: Départements
acmeEngineering = Department(acme, "DEPT-001", "Engineering", 1000000.00)
acmeMarketing = Department(acme, "DEPT-002", "Marketing", 500000.00)
technoRD = Department(techno, "DEPT-003", "R&D", 800000.00)

// Niveau 3: Employés
alice = Employee(acmeEngineering, "EMP-001", "Alice Johnson", 120000.00)
bob = Employee(acmeEngineering, "EMP-002", "Bob Smith", 90000.00)
charlie = Employee(acmeMarketing, "EMP-003", "Charlie Brown", 85000.00)
david = Employee(technoRD, "EMP-004", "David Lee", 110000.00)

// Niveau 3: Projets (parallèle aux employés)
platformProject = Project(acmeEngineering, "PROJ-001", "New Platform", "2024-12-31")
campaignProject = Project(acmeMarketing, "PROJ-002", "Campaign 2024", "2024-06-30")

// Règles d'analyse
rule companyEmployees : {comp: Company, dept: Department, emp: Employee} /
    dept.company == comp && emp.dept == dept
    ==> Log(emp.name + " works at " + comp.name + " in " + dept.name)

rule highSalaries : {emp: Employee, dept: Department} /
    emp.dept == dept && emp.salary > 100000.00
    ==> Log("High salary: " + emp.name + " (" + dept.name + ")")
```

---

### Exemple 3 : E-Commerce

```tsd
type Customer(#customerId: string, name: string, email: string, vip: boolean)
type Product(#sku: string, name: string, category: string, price: number)
type Order(customer: Customer, #orderNumber: string, date: string, status: string)
type OrderLine(order: Order, product: Product, quantity: number, subtotal: number)
type Payment(order: Order, #paymentId: string, amount: number, method: string)

// Clients
vipCustomer = Customer("CUST-001", "Alice Johnson", "alice@example.com", true)
regularCustomer = Customer("CUST-002", "Bob Smith", "bob@example.com", false)

// Catalogue produits
laptop = Product("LAP-001", "Gaming Laptop", "Electronics", 1200.00)
mouse = Product("MOUSE-001", "Wireless Mouse", "Accessories", 25.00)
keyboard = Product("KEY-001", "Mechanical Keyboard", "Accessories", 150.00)
monitor = Product("MON-001", "4K Monitor", "Electronics", 500.00)

// Commandes
order1 = Order(vipCustomer, "ORD-001", "2024-12-19", "confirmed")
order2 = Order(regularCustomer, "ORD-002", "2024-12-19", "pending")

// Lignes de commande
OrderLine(order1, laptop, 1, 1200.00)
OrderLine(order1, mouse, 2, 50.00)
OrderLine(order2, keyboard, 1, 150.00)
OrderLine(order2, monitor, 1, 500.00)

// Paiements
Payment(order1, "PAY-001", 1250.00, "credit_card")

// Règles business
rule vipOrderDiscount : {c: Customer, o: Order} /
    o.customer == c && c.vip == true
    ==> Log("VIP order: " + o.orderNumber + " for " + c.name)

rule highValueOrders : {o: Order, ol: OrderLine, p: Product} /
    ol.order == o && ol.product == p && ol.subtotal > 1000.00
    ==> Log("High value item in " + o.orderNumber + ": " + p.name)

rule unpaidOrders : {o: Order, p: Payment} /
    o.status == "confirmed" && not exists {p2: Payment} / p2.order == o
    ==> Log("Unpaid order: " + o.orderNumber)
```

---

## Bonnes Pratiques

### 1. Noms de Variables Descriptifs

```tsd
// ✅ BON - Noms descriptifs
adminUser = User("admin", "admin@example.com", "administrator")
currentSession = Session(adminUser, "SES-001", 1704067200)
loginAction = Action(currentSession, "login", 1704067200)

// ❌ ÉVITER - Noms cryptiques
u = User("admin", "admin@example.com", "administrator")
s = Session(u, "SES-001", 1704067200)
a = Action(s, "login", 1704067200)
```

**Raison** : La lisibilité est essentielle pour la maintenance.

---

### 2. Grouper par Entité

```tsd
// ✅ BON - Regroupement logique

// Utilisateurs
alice = User("alice", "alice@example.com")
bob = User("bob", "bob@example.com")
charlie = User("charlie", "charlie@example.com")

// Produits
laptop = Product("LAP-001", "Laptop", 1200.00)
mouse = Product("MOUSE-001", "Mouse", 25.00)
keyboard = Product("KEY-001", "Keyboard", 150.00)

// Commandes
order1 = Order(alice, "ORD-001", 1200.00)
order2 = Order(bob, "ORD-002", 25.00)
order3 = Order(charlie, "ORD-003", 150.00)
```

**Raison** : Organisation claire facilite la compréhension.

---

### 3. Cohérence de Nommage

```tsd
// ✅ BON - Convention cohérente
user1 = User("alice", "alice@example.com")
user2 = User("bob", "bob@example.com")
user3 = User("charlie", "charlie@example.com")

session1 = Session(user1, "SES-001", 1704067200)
session2 = Session(user2, "SES-002", 1704067260)
session3 = Session(user3, "SES-003", 1704067320)

// ❌ ÉVITER - Incohérence
alice = User("alice", "alice@example.com")
bobUser = User("bob", "bob@example.com")
u3 = User("charlie", "charlie@example.com")
```

**Raison** : Patterns cohérents réduisent la charge cognitive.

---

### 4. Éviter les Variables Temporaires Inutiles

```tsd
// ✅ BON - Variable réutilisée plusieurs fois
alice = User("alice", "alice@example.com")
Order(alice, "ORD-001", 150.00)
Order(alice, "ORD-002", 75.00)
Login(alice, "SES-001")

// ⚠️ ACCEPTABLE - Variable utilisée une seule fois mais améliore lisibilité
admin = User("admin", "admin@example.com", "administrator")
Session(admin, "SES-ADMIN", 1704067200)

// ❌ INUTILE - Jamais réutilisé, pas de gain de lisibilité
temp = Product("TEMP-001", "Temporary", 0.00)
```

**Règle** : Utilisez des variables si :
- Réutilisation multiple, OU
- Amélioration de la lisibilité

---

### 5. Documentation avec Commentaires

```tsd
// ✅ BON - Commentaires explicatifs

// Admin system user
admin = User("admin", "admin@example.com", "administrator")

// VIP customers
vip1 = Customer("VIP-001", "Alice Premium", true)
vip2 = Customer("VIP-002", "Bob Elite", true)

// Regular customers
regular1 = Customer("REG-001", "Charlie Normal", false)
regular2 = Customer("REG-002", "David Standard", false)

// High-value products
laptop = Product("LAP-001", "Gaming Laptop", 1200.00)
workstation = Product("WKS-001", "Pro Workstation", 2500.00)
```

---

## Utilisation avec Règles

Les affectations facilitent l'écriture de **règles complexes** :

```tsd
type User(#username: string, email: string, active: boolean)
type Order(customer: User, #orderNum: string, total: number, status: string)
type Payment(order: Order, #paymentId: string, amount: number, method: string)

// Définir des faits
alice = User("alice", "alice@example.com", true)
bob = User("bob", "bob@example.com", false)

order1 = Order(alice, "ORD-001", 150.00, "confirmed")
order2 = Order(bob, "ORD-002", 75.00, "pending")

Payment(order1, "PAY-001", 150.00, "credit_card")

// Règles utilisant les relations
rule activeUserOrders : {u: User, o: Order} / 
    o.customer == u && u.active == true 
    ==> Log("Active user order: " + o.orderNum + " for " + u.username)

rule paidOrders : {o: Order, p: Payment} /
    p.order == o && o.status == "confirmed"
    ==> Log("Order " + o.orderNum + " paid with " + p.method)

rule unpaidUserOrders : {u: User, o: Order} /
    o.customer == u && 
    not exists {p: Payment} / p.order == o
    ==> Log("Unpaid order " + o.orderNum + " for user " + u.username)

rule customerOrderChain : {u: User, o: Order, p: Payment} /
    o.customer == u && p.order == o
    ==> Log("Complete chain: " + u.username + " → " + o.orderNum + " → " + p.paymentId)
```

**Avantages** :
- Relations explicites dans les règles
- Pas de duplication de code de jointure
- Type-safe (erreurs détectées au parsing)
- Navigation dans les chaînes de références

---

## Résumé

### Syntaxe

```tsd
variable = Type(valeur1, valeur2, ...)
```

### Avantages

| Aspect | Bénéfice |
|--------|----------|
| **Réutilisation** | Pas de duplication de données |
| **Lisibilité** | Code plus clair et expressif |
| **Maintenabilité** | Modifications centralisées |
| **Type-safety** | Erreurs détectées au parsing |
| **Relations** | Références explicites entre faits |

### Contraintes

| Règle | Description |
|-------|-------------|
| **Unicité** | Variable définie une seule fois |
| **Ordre** | Définir avant utiliser |
| **Type** | Compatibilité de types obligatoire |
| **Portée** | Globale dans le programme |

### Workflow Type

```tsd
// 1. Définir les types avec relations
type User(#username: string, ...)
type Order(customer: User, ...)

// 2. Créer les faits avec affectations
alice = User("alice", ...)
order1 = Order(alice, ...)

// 3. Réutiliser dans les règles
{u: User, o: Order} / o.customer == u ==> ...
```

---

## Voir Aussi

- [Comparaisons de Faits](fact-comparisons.md) - Utiliser `==` entre faits
- [Système de Types](type-system.md) - Types de faits dans les champs
- [Identifiants Internes](../internal-ids.md) - Comment fonctionnent les IDs
- [Guide de Migration](../migration/from-v1.x.md) - Migrer depuis v1.x

---

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-19  
**Auteur** : Équipe TSD
