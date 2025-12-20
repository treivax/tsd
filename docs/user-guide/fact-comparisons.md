# Guide Utilisateur - Comparaisons de Faits

## Introduction

Les **comparaisons de faits** permettent de comparer directement deux faits dans les conditions des règles, établissant des relations de manière naturelle et type-safe.

Cette fonctionnalité est **nouvelle en v2.0** et simplifie considérablement l'écriture des règles impliquant des relations entre faits.

---

## 📋 Table des Matières

1. [Syntaxe](#syntaxe)
2. [Fonctionnement Interne](#fonctionnement-interne)
3. [Cas d'Usage](#cas-dusage)
4. [Opérateurs Disponibles](#opérateurs-disponibles)
5. [Comparaisons avec Champs](#comparaisons-avec-champs)
6. [Exemples Pratiques](#exemples-pratiques)
7. [Bonnes Pratiques](#bonnes-pratiques)

---

## Syntaxe

### Forme Générale

```tsd
{var1: Type1, var2: Type2} / var1.champFait == var2 ==> Action
```

**Composants** :
- `var1.champFait` : Champ de type fait
- `==` : Opérateur de comparaison
- `var2` : Variable représentant un fait

### Exemple Basique

```tsd
type User(#username: string, email: string)
type Login(user: User, #sessionId: string, timestamp: number)

alice = User("alice", "alice@example.com")
Login(alice, "SES-001", 1704067200)

// Règle avec comparaison de faits
rule userLogins : {u: User, l: Login} / l.user == u ==> 
    Log("Login for: " + u.username)
```

**Lecture** : "Pour chaque paire (User, Login) où le champ `user` du Login est égal au User `u`..."

---

## Fonctionnement Interne

### Mécanisme de Comparaison

Les comparaisons de faits utilisent les **identifiants internes** (`_id_`) de manière transparente :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)

alice = User("alice", "alice@example.com")
order1 = Order(alice, "ORD-001", 150.00)

{u: User, o: Order} / o.customer == u ==> 
    Log("Match!")
```

**En interne (transparent pour vous)** :
1. `alice` reçoit l'ID interne : `"User~alice"`
2. `order1.customer` stocke la référence : `"User~alice"`
3. La comparaison `o.customer == u` compare : `"User~alice" == "User~alice"` → `true`

**Important** : Vous n'avez jamais à manipuler les IDs. Le système gère tout automatiquement.

---

### Équivalence avec v1.x

Pour mieux comprendre, voici la différence avec l'ancienne approche :

#### ❌ v1.x (Obsolète)

```tsd
type User(#userId: string, name: string)
type Login(userId: string, #sessionId: string)

assert User(userId: "U001", name: "Alice")
assert Login(userId: "U001", sessionId: "SES-001")

rule userLogins : {u: User, l: Login} / l.userId == u.userId ==> 
    Log("Login for: " + u.name)
```

**Problèmes** :
- Duplication de `userId`
- Pas de type-safety
- Risque d'incohérence

#### ✅ v2.0 (Correct)

```tsd
type User(#userId: string, name: string)
type Login(user: User, #sessionId: string)

alice = User("U001", "Alice")
Login(alice, "SES-001")

rule userLogins : {u: User, l: Login} / l.user == u ==> 
    Log("Login for: " + u.name)
```

**Avantages** :
- ✅ Pas de duplication
- ✅ Type-safe
- ✅ Plus simple et lisible
- ✅ Cohérence garantie

---

## Cas d'Usage

### 1. Relations One-to-Many

Un utilisateur peut avoir plusieurs commandes :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)

alice = User("alice", "alice@example.com")

Order(alice, "ORD-001", 100.00)
Order(alice, "ORD-002", 150.00)
Order(alice, "ORD-003", 75.00)

// Règle : Trouver toutes les commandes d'un utilisateur
rule userOrders : {u: User, o: Order} / o.customer == u && u.username == "alice" ==> 
    Log("Alice's order: " + o.orderNum + " for $" + o.total)
```

---

### 2. Relations Many-to-Many

Étudiants et cours (via une table de jonction) :

```tsd
type Student(#studentId: string, name: string)
type Course(#courseId: string, title: string)
type Enrollment(student: Student, course: Course, grade: string)

alice = Student("S001", "Alice")
bob = Student("S002", "Bob")

math = Course("C001", "Math 101")
physics = Course("C002", "Physics 101")

Enrollment(alice, math, "A")
Enrollment(alice, physics, "B")
Enrollment(bob, math, "B")

// Règle : Afficher les notes des étudiants
rule studentGrades : {s: Student, e: Enrollment, c: Course} /
    e.student == s && e.course == c
    ==> Log(s.name + " got " + e.grade + " in " + c.title)
```

---

### 3. Chaînes de Relations

Naviguer dans plusieurs niveaux de relations :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)
type Payment(order: Order, #paymentId: string, amount: number)

alice = User("alice", "alice@example.com")
order1 = Order(alice, "ORD-001", 150.00)
Payment(order1, "PAY-001", 150.00)

// Règle : Chaîne User → Order → Payment
rule completeChain : {u: User, o: Order, p: Payment} /
    o.customer == u && p.order == o
    ==> Log("Payment " + p.paymentId + " for " + u.username + "'s order " + o.orderNum)
```

**Chaîne de relations** : `Payment → Order → User`

---

### 4. Comparaisons Multiples

Plusieurs comparaisons dans une même règle :

```tsd
type User(#username: string, role: string)
type Session(user: User, #sessionId: string, active: boolean)
type Action(session: Session, action: string, timestamp: number)

admin = User("admin", "administrator")
adminSession = Session(admin, "SES-001", true)
Action(adminSession, "delete_user", 1704067200)

// Règle : Actions admin actives
rule activeAdminActions : {u: User, s: Session, a: Action} /
    s.user == u && a.session == s && u.role == "administrator" && s.active == true
    ==> Log("Admin action: " + a.action + " by " + u.username)
```

**Deux comparaisons de faits** :
- `s.user == u` : Session liée à User
- `a.session == s` : Action liée à Session

---

## Opérateurs Disponibles

### Égalité (`==`)

Compare deux faits pour vérifier s'ils sont identiques :

```tsd
{u: User, o: Order} / o.customer == u ==> ...
```

**Retourne `true` si** : Les deux faits ont le même ID interne.

---

### Inégalité (`!=`)

Compare deux faits pour vérifier s'ils sont différents :

```tsd
type Order(customer: User, #orderNum: string)

alice = User("alice", "alice@example.com")
bob = User("bob", "bob@example.com")

Order(alice, "ORD-001")
Order(bob, "ORD-002")

// Règle : Commandes de clients différents
rule differentCustomers : {o1: Order, o2: Order} /
    o1 != o2 && o1.customer != o2.customer
    ==> Log("Different customers: " + o1.orderNum + " and " + o2.orderNum)
```

**Retourne `true` si** : Les deux faits ont des IDs internes différents.

---

## Comparaisons avec Champs

### Accéder aux Champs via les Références

Vous pouvez accéder aux champs d'un fait référencé :

```tsd
type User(#username: string, email: string, active: boolean)
type Order(customer: User, #orderNum: string, total: number)

alice = User("alice", "alice@example.com", true)
Order(alice, "ORD-001", 150.00)

// Règle : Accéder aux champs du customer via la référence
rule activeCustomerOrders : {o: Order} /
    o.customer.active == true
    ==> Log("Order " + o.orderNum + " from active user " + o.customer.username)
```

**Navigation** : `o.customer.active` suit la référence et accède au champ `active`.

---

### Comparaison Mixte (Fait + Champ)

Combiner comparaisons de faits et de champs :

```tsd
type User(#username: string, vip: boolean)
type Order(customer: User, #orderNum: string, total: number)

alice = User("alice", true)
bob = User("bob", false)

Order(alice, "ORD-001", 1000.00)
Order(bob, "ORD-002", 50.00)

// Règle : VIP avec commandes > $500
rule vipHighValue : {u: User, o: Order} /
    o.customer == u && u.vip == true && o.total > 500.00
    ==> Log("VIP high-value order: " + o.orderNum)
```

**Combinaison** :
- `o.customer == u` : Comparaison de faits
- `u.vip == true` : Comparaison de champ booléen
- `o.total > 500.00` : Comparaison de champ numérique

---

## Exemples Pratiques

### Exemple 1 : Système de Blog

```tsd
type User(#username: string, email: string, bio: string)
type Post(author: User, #postId: string, title: string, content: string, published: boolean)
type Comment(post: Post, author: User, #commentId: string, text: string, approved: boolean)

// Utilisateurs
alice = User("alice", "alice@example.com", "Software engineer")
bob = User("bob", "bob@example.com", "Tech enthusiast")
charlie = User("charlie", "charlie@example.com", "Blogger")

// Posts
post1 = Post(alice, "POST-001", "Introduction to TSD", "Welcome to TSD!", true)
post2 = Post(bob, "POST-002", "Advanced RETE", "Deep dive...", true)
post3 = Post(alice, "POST-003", "Draft Post", "Not ready yet", false)

// Commentaires
Comment(post1, bob, "COM-001", "Great introduction!", true)
Comment(post1, charlie, "COM-002", "Very helpful!", true)
Comment(post2, alice, "COM-003", "Excellent article!", true)
Comment(post2, charlie, "COM-004", "Spam comment", false)

// Règle : Posts publiés avec commentaires approuvés
rule publishedPostComments : {p: Post, c: Comment, author: User} /
    c.post == p && c.author == author && p.published == true && c.approved == true
    ==> Log(author.username + " commented on published post: " + p.title)

// Règle : Auteurs qui commentent leurs propres posts
rule selfComments : {p: Post, c: Comment, u: User} /
    p.author == u && c.author == u && c.post == p
    ==> Log("Self-comment by " + u.username + " on: " + p.title)

// Règle : Posts avec plusieurs commentaires
rule popularPosts : {p: Post, c1: Comment, c2: Comment} /
    c1.post == p && c2.post == p && c1 != c2 && p.published == true
    ==> Log("Popular post: " + p.title + " has multiple comments")
```

---

### Exemple 2 : Système de Réservation

```tsd
type Customer(#customerId: string, name: string, email: string, loyaltyLevel: string)
type Room(#roomNum: string, type: string, pricePerNight: number, available: boolean)
type Reservation(customer: Customer, room: Room, #reservationId: string, 
                 checkIn: string, checkOut: string, status: string)
type Payment(reservation: Reservation, #paymentId: string, amount: number, method: string)

// Clients
vipCustomer = Customer("C001", "Alice Premium", "alice@example.com", "VIP")
goldCustomer = Customer("C002", "Bob Gold", "bob@example.com", "Gold")
regularCustomer = Customer("C003", "Charlie Normal", "charlie@example.com", "Regular")

// Chambres
suite = Room("R101", "Suite", 300.00, false)
deluxe = Room("R102", "Deluxe", 200.00, true)
standard = Room("R103", "Standard", 100.00, true)

// Réservations
res1 = Reservation(vipCustomer, suite, "RES-001", "2024-12-20", "2024-12-25", "confirmed")
res2 = Reservation(goldCustomer, deluxe, "RES-002", "2024-12-21", "2024-12-23", "confirmed")
res3 = Reservation(regularCustomer, standard, "RES-003", "2024-12-22", "2024-12-24", "pending")

// Paiements
Payment(res1, "PAY-001", 1500.00, "credit_card")
Payment(res2, "PAY-002", 400.00, "paypal")

// Règles métier

// VIP avec suites
rule vipSuites : {c: Customer, res: Reservation, r: Room} /
    res.customer == c && res.room == r && c.loyaltyLevel == "VIP" && r.type == "Suite"
    ==> Log("VIP " + c.name + " reserved suite " + r.roomNum)

// Réservations confirmées payées
rule confirmedPaidReservations : {res: Reservation, p: Payment, c: Customer} /
    p.reservation == res && res.customer == c && res.status == "confirmed"
    ==> Log("Confirmed paid reservation " + res.reservationId + " for " + c.name)

// Réservations en attente de paiement
rule unpaidReservations : {res: Reservation, c: Customer} /
    res.customer == c && res.status == "pending" &&
    not exists {p: Payment} / p.reservation == res
    ==> Log("Unpaid reservation " + res.reservationId + " for " + c.name)

// Chambres réservées (indisponibles)
rule reservedRooms : {r: Room, res: Reservation} /
    res.room == r && res.status == "confirmed" && r.available == false
    ==> Log("Room " + r.roomNum + " is reserved (unavailable)")
```

---

### Exemple 3 : Workflow d'Approbation

```tsd
type User(#username: string, role: string, department: string)
type Document(author: User, #docId: string, title: string, status: string)
type Approval(document: Document, approver: User, #approvalId: string, 
              decision: string, comment: string)

// Utilisateurs
employee = User("alice", "employee", "Engineering")
manager = User("bob", "manager", "Engineering")
director = User("charlie", "director", "Engineering")

// Documents
doc1 = Document(employee, "DOC-001", "Technical Proposal", "pending")
doc2 = Document(employee, "DOC-002", "Budget Request", "pending")

// Approbations
Approval(doc1, manager, "APP-001", "approved", "Looks good")
Approval(doc1, director, "APP-002", "approved", "Approved")
Approval(doc2, manager, "APP-003", "rejected", "Needs revision")

// Règles

// Documents approuvés par tous les niveaux
rule fullyApprovedDocs : {d: Document, a1: Approval, a2: Approval, u1: User, u2: User} /
    a1.document == d && a2.document == d && a1 != a2 &&
    a1.approver == u1 && a2.approver == u2 &&
    u1.role == "manager" && u2.role == "director" &&
    a1.decision == "approved" && a2.decision == "approved"
    ==> Log("Fully approved: " + d.title)

// Documents rejetés
rule rejectedDocs : {d: Document, a: Approval, author: User, approver: User} /
    d.author == author && a.document == d && a.approver == approver &&
    a.decision == "rejected"
    ==> Log("Rejected: " + d.title + " by " + approver.username + " - " + a.comment)

// Documents en attente d'approbation d'un directeur
rule pendingDirectorApproval : {d: Document, a: Approval, mgr: User} /
    a.document == d && a.approver == mgr && mgr.role == "manager" && 
    a.decision == "approved" &&
    not exists {a2: Approval, dir: User} / a2.document == d && a2.approver == dir && dir.role == "director"
    ==> Log("Pending director approval: " + d.title)
```

---

## Bonnes Pratiques

### 1. Comparaisons Explicites

```tsd
// ✅ BON - Comparaison explicite et claire
{u: User, o: Order} / o.customer == u ==> ...

// ❌ ÉVITER - Comparaison implicite (n'existe pas en TSD)
{u: User, o: Order} / ... // Pas de comparaison automatique
```

**Raison** : Soyez explicite sur les relations.

---

### 2. Ordre des Comparaisons

```tsd
// ✅ BON - Comparaisons de faits en premier
{u: User, o: Order} / 
    o.customer == u && u.active == true && o.total > 100.00
    ==> ...

// ⚠️ ACCEPTABLE mais moins lisible
{u: User, o: Order} / 
    u.active == true && o.total > 100.00 && o.customer == u
    ==> ...
```

**Conseil** : Placez les comparaisons de faits en premier pour la lisibilité.

---

### 3. Utiliser `!=` pour Exclusions

```tsd
// ✅ BON - Exclure les auto-références
{o1: Order, o2: Order} / 
    o1 != o2 && o1.customer == o2.customer
    ==> Log("Same customer has multiple orders")

// ❌ MAUVAIS - Peut inclure les auto-comparaisons
{o1: Order, o2: Order} / 
    o1.customer == o2.customer
    ==> Log("...")  // o1 == o2 possible
```

**Raison** : `!=` évite les faux positifs.

---

### 4. Décomposer les Règles Complexes

```tsd
// ✅ BON - Règles séparées et claires
rule paidOrders : {o: Order, p: Payment} /
    p.order == o
    ==> Log("Paid: " + o.orderNum)

rule unpaidOrders : {o: Order} /
    not exists {p: Payment} / p.order == o
    ==> Log("Unpaid: " + o.orderNum)

// ❌ ÉVITER - Règle unique trop complexe
rule allOrders : {o: Order} /
    (exists {p: Payment} / p.order == o) || 
    (not exists {p: Payment} / p.order == o)
    ==> Log("...")  // Trop complexe
```

**Raison** : Règles simples sont plus maintenables et debuggables.

---

### 5. Nommer les Variables de Façon Descriptive

```tsd
// ✅ BON - Noms descriptifs
rule customerOrders : {customer: User, order: Order} /
    order.customer == customer
    ==> Log(customer.username + " ordered " + order.orderNum)

// ❌ ÉVITER - Noms génériques
rule r1 : {u: User, o: Order} /
    o.customer == u
    ==> Log(u.username + " ordered " + o.orderNum)
```

**Raison** : Noms descriptifs améliorent la lisibilité.

---

## Résumé

### Syntaxe Essentielle

```tsd
{var1: Type1, var2: Type2} / var1.factField == var2 ==> Action
```

### Opérateurs

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| `==` | Égalité | `o.customer == u` |
| `!=` | Inégalité | `o1 != o2` |

### Avantages

| Aspect | Bénéfice |
|--------|----------|
| **Simplicité** | Syntaxe naturelle et intuitive |
| **Type-safety** | Erreurs détectées au parsing |
| **Performance** | Optimisations internes du moteur RETE |
| **Lisibilité** | Code plus clair que les jointures manuelles |
| **Maintenabilité** | Moins de code, moins d'erreurs |

### Points Clés

1. Les comparaisons utilisent les IDs internes (transparence totale)
2. Syntaxe simple : `fact1 == fact2`
3. Fonctionne avec toutes les relations (1-N, N-N, chaînes)
4. Accès aux champs via navigation : `fact.ref.field`
5. Combiner avec autres conditions : `&&`, `||`, `not`

### Workflow Type

```tsd
// 1. Définir types avec relations
type User(#username: string, ...)
type Order(customer: User, ...)

// 2. Créer faits avec affectations
alice = User("alice", ...)
order1 = Order(alice, ...)

// 3. Comparer dans les règles
{u: User, o: Order} / o.customer == u ==> 
    Log("Match!")
```

---

## Voir Aussi

- [Affectations de Faits](fact-assignments.md) - Créer et nommer des faits
- [Système de Types](type-system.md) - Types de faits dans les champs
- [Identifiants Internes](../internal-ids.md) - Fonctionnement des IDs
- [Guide de Migration](../migration/from-v1.x.md) - Migrer depuis v1.x

---

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-19  
**Auteur** : Équipe TSD
