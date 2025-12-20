# Prompt 09 - Mise à Jour de la Documentation

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Mettre à jour toute la documentation du projet pour refléter la nouvelle gestion des identifiants :

1. **Documentation technique** - Centraliser dans `docs/`
2. **Guide utilisateur** - Nouvelles fonctionnalités
3. **Guide de migration** - Breaking changes
4. **Exemples** - Démonstrations complètes
5. **README** - Mise à jour des liens et contenus
6. **Supprimer l'obsolète** - Nettoyer les anciennes docs

---

## 📋 Contexte

### État Actuel

Documentation existante :
- `docs/ID_RULES_COMPLETE.md` - Règles actuelles des IDs
- `docs/primary-keys.md` - Documentation clés primaires
- `docs/MIGRATION_IDS.md` - Guide de migration existant
- README modules
- Docs obsolètes à supprimer

### État Cible

Documentation complète et à jour :
- **Guide de référence** : Syntaxe TSD complète
- **Guide utilisateur** : Utilisation des nouvelles fonctionnalités
- **Guide de migration** : Breaking changes et migration
- **Exemples** : Cas d'usage réels
- **Architecture** : Système interne
- **Obsolète supprimé** : Anciennes docs retirées

---

## 📝 Tâches à Réaliser

### 1. Analyser la Documentation Existante

#### Inventaire Complet

```bash
# Lister toute la documentation
find docs/ -name "*.md" -type f | sort

# Identifier les docs obsolètes
grep -r "id:" docs/ --include="*.md" | grep -v "_id_"
grep -r '"id"' docs/ --include="*.md" | grep -v FieldNameInternalID

# Identifier les références à l'ancienne syntaxe
grep -r "assert.*id:" docs/ --include="*.md"
```

**Créer rapport** : `REPORTS/new_ids_docs_audit.md`

```markdown
# Audit Documentation - Migration IDs

## Documentation à Mettre à Jour

### Priorité 1 - Critique
- [ ] docs/ID_RULES_COMPLETE.md - Réécrire complètement
- [ ] docs/primary-keys.md - Adapter nouvelle syntaxe
- [ ] docs/MIGRATION_IDS.md - Compléter
- [ ] README.md - Mise à jour générale

### Priorité 2 - Important
- [ ] docs/guides.md - Exemples à jour
- [ ] docs/reference.md - Référence complète
- [ ] docs/api.md - API publique
- [ ] constraint/README.md - Module constraint
- [ ] rete/README.md - Module RETE

### Priorité 3 - Secondaire
- [ ] docs/architecture.md - Diagrammes
- [ ] docs/tutorials/ - Tutoriels
- [ ] examples/README.md - Guide exemples

## Documentation à Supprimer

- [ ] Anciennes règles d'ID obsolètes
- [ ] Exemples avec syntaxe deprecated
- [ ] Docs contradictoires

## Nouvelle Documentation à Créer

- [ ] docs/user-guide/fact-assignments.md
- [ ] docs/user-guide/fact-comparisons.md
- [ ] docs/user-guide/type-system.md
- [ ] docs/migration/from-v1.x.md
- [ ] docs/examples/ (restructurer)
```

### 2. Réécrire la Documentation des IDs

#### Fichier : `docs/internal-ids.md` (nouveau)

**Remplace `ID_RULES_COMPLETE.md`**

```markdown
# Identifiants Internes - Documentation Complète

## Vue d'Ensemble

Dans TSD, chaque fait possède un **identifiant interne unique** (`_id_`) qui est :

1. **Généré automatiquement** - Jamais défini manuellement
2. **Déterministe** - Basé sur les clés primaires ou hash
3. **Caché** - Jamais accessible dans les expressions TSD
4. **Interne** - Utilisé uniquement par le moteur RETE

---

## ⚠️ Règle Fondamentale

Le champ `_id_` est **strictement réservé au système** :

❌ **INTERDIT**
- Définir `_id_` dans une définition de type
- Assigner `_id_` dans un fait
- Accéder à `_id_` dans une expression
- Comparer `_id_` explicitement

✅ **PERMIS**
- Les IDs sont générés automatiquement
- Les comparaisons de faits utilisent les IDs en interne
- Les références entre faits sont résolues automatiquement

---

## Génération Automatique

### Avec Clés Primaires

Les clés primaires (préfixées par `#`) déterminent l'ID :

```tsd
type User(#username: string, email: string, age: number)

// Définir un utilisateur
alice = User("alice", "alice@example.com", 30)
// ID généré en interne: "User~alice"
```

**Format** : `TypeName~valeur1_valeur2_...`

#### Clé Primaire Simple

```tsd
type Product(#sku: string, name: string, price: number)

laptop = Product("LAPTOP-001", "Gaming Laptop", 1200.00)
// ID: "Product~LAPTOP-001"
```

#### Clé Primaire Composite

```tsd
type OrderLine(#orderId: string, #productId: string, quantity: number)

line1 = OrderLine("ORD-001", "PROD-123", 2)
// ID: "OrderLine~ORD-001_PROD-123"
```

### Sans Clé Primaire (Hash)

Si aucun champ n'est marqué `#`, un hash déterministe est utilisé :

```tsd
type LogEvent(timestamp: number, level: string, message: string)

LogEvent(1704067200, "INFO", "Application started")
// ID: "LogEvent~a1b2c3d4e5f6g7h8" (hash 16 caractères)
```

---

## Caractères Spéciaux

Les caractères spéciaux dans les clés primaires sont encodés :

| Caractère | Encodage | Raison |
|-----------|----------|--------|
| `~` | `%7E` | Séparateur type/valeur |
| `_` | `%5F` | Séparateur clés composites |
| `%` | `%25` | Caractère d'échappement |
| ` ` | `%20` | Espace |

**Exemple** :

```tsd
type File(#path: string, size: number)

file1 = File("/home/user~backup_v1", 1024)
// ID: "File~%2Fhome%2Fuser%7Ebackup%5Fv1"
```

---

## Utilisation Interne

### Comparaisons de Faits

Les comparaisons de faits utilisent automatiquement les IDs :

```tsd
type User(#username: string, email: string)
type Login(user: User, #sessionId: string, timestamp: number)

alice = User("alice", "alice@example.com")
session1 = Login(alice, "SES-001", 1704067200)

// Comparaison via IDs internes (automatique)
{u: User, l: Login} / l.user == u ==> 
    Log("Login for user: " + u.username)
```

**Fonctionnement interne** :
1. `l.user` retourne l'ID interne : `"User~alice"`
2. `u` est résolu vers son ID : `"User~alice"`
3. Comparaison : `"User~alice" == "User~alice"` → `true`

### Références entre Faits

Les champs de type fait stockent l'ID en interne :

```tsd
type User(#username: string)
type Order(user: User, #orderNum: string, total: number)

alice = User("alice")
order1 = Order(alice, "ORD-001", 150.00)

// En interne, order1.user = "User~alice"
```

---

## Déterminisme

Les IDs sont **toujours les mêmes** pour les mêmes valeurs :

```tsd
type User(#username: string, age: number)

// Ces deux définitions génèrent le MÊME ID
alice1 = User("alice", 30)
alice2 = User("alice", 30)
// Les deux ont l'ID: "User~alice"
```

⚠️ **Attention** : Cela signifie que le deuxième fait **remplace** le premier (même identité).

---

## Bonnes Pratiques

### 1. Choisir des Clés Primaires Naturelles

```tsd
// ✅ BON - Clé naturelle
type User(#username: string, email: string)

// ❌ MOINS BON - Pas de clé (hash)
type User(username: string, email: string)
```

### 2. Clés Stables

```tsd
// ✅ BON - username change rarement
type User(#username: string, email: string)

// ❌ ÉVITER - email peut changer
type User(username: string, #email: string)
```

### 3. Clés Composites pour Relations N-N

```tsd
// ✅ BON - Clé composite pour table de jonction
type Enrollment(#studentId: string, #courseId: string, grade: string)
```

### 4. Hash pour Événements Temporels

```tsd
// ✅ BON - Événements n'ont pas de clé naturelle
type AuditLog(timestamp: number, userId: string, action: string)
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

// Créer des posts
post1 = Post(alice, "POST-001", "Introduction to TSD", "Welcome!")
post2 = Post(bob, "POST-002", "Advanced Features", "Deep dive...")

// Créer des commentaires
Comment(post1, bob, "COM-001", "Great post!")
Comment(post2, alice, "COM-002", "Very helpful!")

// Règles utilisant les comparaisons de faits
{p: Post, c: Comment, u: User} / 
    c.post == p && c.author == u && u.username == "alice" ==> 
    Log("Alice commented: " + c.text)
```

### Exemple 2 : E-Commerce

```tsd
type Customer(#customerId: string, name: string, vip: bool)
type Product(#sku: string, name: string, price: number)
type Order(customer: Customer, #orderNumber: string, total: number)
type OrderLine(order: Order, product: Product, quantity: number)

cust1 = Customer("CUST-001", "Alice Johnson", true)
prod1 = Product("LAPTOP-001", "Gaming Laptop", 1200.00)
prod2 = Product("MOUSE-001", "Wireless Mouse", 25.00)

order1 = Order(cust1, "ORD-001", 1250.00)

OrderLine(order1, prod1, 1)
OrderLine(order1, prod2, 2)

// Règle : Commandes VIP avec produits
{c: Customer, o: Order, ol: OrderLine, p: Product} / 
    o.customer == c && ol.order == o && ol.product == p && c.vip == true ==> 
    Log("VIP order " + o.orderNumber + " contains " + p.name)
```

---

## Migration

### Ancien Système (v1.x)

```tsd
// ❌ OBSOLÈTE - Ne fonctionne plus
type Person(name: string, age: number)
Person(id: "person_1", name: "Alice", age: 30)

{p: Person} / p.id == "person_1" ==> Log("Found")
```

### Nouveau Système (v2.0)

```tsd
// ✅ NOUVEAU
type Person(#name: string, age: number)
alice = Person("Alice", 30)

// Comparaison via variable, pas via ID
{p: Person} / p.name == "Alice" ==> Log("Found")
```

**Voir** : [Guide de Migration](migration/from-v1.x.md)

---

## Résumé

| Aspect | Comportement |
|--------|--------------|
| **Nom** | `_id_` (caché) |
| **Génération** | Automatique, toujours |
| **Format** | `Type~value` ou `Type~hash` |
| **Accès** | ❌ Jamais dans expressions TSD |
| **Comparaisons** | ✅ Automatiques en interne |
| **Affectation** | ❌ Interdite |

---

## Références

- [Guide Utilisateur - Affectations](user-guide/fact-assignments.md)
- [Guide Utilisateur - Comparaisons](user-guide/fact-comparisons.md)
- [Système de Types](user-guide/type-system.md)
- [Guide de Migration](migration/from-v1.x.md)
- [Architecture Interne](architecture/id-generation.md)

---

**Note** : Cette documentation décrit le système à partir de la version 2.0.
```

### 3. Créer Guide Utilisateur - Affectations

#### Fichier : `docs/user-guide/fact-assignments.md` (nouveau)

```markdown
# Guide Utilisateur - Affectations de Faits

## Introduction

Les **affectations de faits** permettent de nommer des faits pour les référencer dans d'autres définitions.

---

## Syntaxe

```tsd
variable = TypeName(champ1, champ2, ...)
```

**Exemples** :

```tsd
alice = User("Alice", 30)
laptop = Product("LAPTOP-001", "Gaming Laptop", 1200.00)
order1 = Order(alice, "ORD-001", 150.00)
```

---

## Utilisation

### Référencer dans d'Autres Faits

```tsd
type User(#username: string, email: string)
type Login(user: User, #sessionId: string, timestamp: number)

// Affecter un utilisateur
alice = User("alice", "alice@example.com")

// Utiliser la variable dans un autre fait
Login(alice, "SES-001", 1704067200)
```

**Fonctionnement** :
1. `alice` est créé et son ID interne est `"User~alice"`
2. Le `Login` stocke en interne `user: "User~alice"`
3. La relation est maintenue automatiquement

### Réutilisation de Variables

```tsd
alice = User("alice", "alice@example.com")

// Utiliser alice plusieurs fois
Login(alice, "SES-001", 1704067200)
Login(alice, "SES-002", 1704068000)
```

### Chaînes de Références

```tsd
type User(#username: string)
type Order(user: User, #orderNum: string)
type Payment(order: Order, #paymentId: string, amount: number)

alice = User("alice")
order1 = Order(alice, "ORD-001")
Payment(order1, "PAY-001", 150.00)
```

---

## Règles et Contraintes

### Variables Uniques

Une variable ne peut être définie qu'une seule fois :

```tsd
alice = User("alice", "alice@example.com")
alice = User("bob", "bob@example.com")  // ❌ ERREUR : alice déjà définie
```

### Variables Avant Utilisation

Une variable doit être définie avant d'être utilisée :

```tsd
Login(alice, "SES-001", 1704067200)  // ❌ ERREUR : alice non définie
alice = User("alice", "alice@example.com")
```

### Types Compatibles

La variable doit être du bon type :

```tsd
type User(#username: string)
type Login(user: User, #sessionId: string)

alice = User("alice")
product1 = Product("LAPTOP-001", "Laptop", 1200.00)

Login(product1, "SES-001")  // ❌ ERREUR : product1 n'est pas un User
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

// Enregistrer des actions
Action(adminSession, "create_user", 1704067300)
Action(aliceSession, "view_dashboard", 1704067320)
```

### Exemple 2 : Hiérarchie Organisationnelle

```tsd
type Company(#companyId: string, name: string)
type Department(company: Company, #deptId: string, name: string)
type Employee(dept: Department, #empId: string, name: string, salary: number)

// Créer une entreprise
acme = Company("COMP-001", "ACME Corp")

// Créer des départements
engineering = Department(acme, "DEPT-001", "Engineering")
marketing = Department(acme, "DEPT-002", "Marketing")

// Créer des employés
alice = Employee(engineering, "EMP-001", "Alice Johnson", 120000.00)
bob = Employee(engineering, "EMP-002", "Bob Smith", 90000.00)
charlie = Employee(marketing, "EMP-003", "Charlie Brown", 85000.00)
```

---

## Bonnes Pratiques

### 1. Noms de Variables Descriptifs

```tsd
// ✅ BON
adminUser = User("admin", "admin@example.com", "administrator")
currentSession = Session(adminUser, "SES-001", 1704067200)

// ❌ MOINS BON
u = User("admin", "admin@example.com", "administrator")
s = Session(u, "SES-001", 1704067200)
```

### 2. Grouper par Entité

```tsd
// ✅ BON - Regrouper logiquement
// Utilisateurs
alice = User("alice", "alice@example.com")
bob = User("bob", "bob@example.com")

// Produits
laptop = Product("LAPTOP-001", "Laptop", 1200.00)
mouse = Product("MOUSE-001", "Mouse", 25.00)

// Commandes
order1 = Order(alice, "ORD-001", 1200.00)
order2 = Order(bob, "ORD-002", 25.00)
```

### 3. Cohérence de Nommage

```tsd
// ✅ BON - Convention cohérente
user1 = User("alice", "alice@example.com")
user2 = User("bob", "bob@example.com")

session1 = Session(user1, "SES-001", 1704067200)
session2 = Session(user2, "SES-002", 1704067260)
```

---

## Utilisation avec Règles

Les affectations facilitent les règles complexes :

```tsd
type User(#username: string, active: bool)
type Order(user: User, #orderNum: string, total: number)
type Payment(order: Order, #paymentId: string, amount: number)

alice = User("alice", true)
order1 = Order(alice, "ORD-001", 150.00)
Payment(order1, "PAY-001", 150.00)

// Règle utilisant la chaîne User -> Order -> Payment
{u: User, o: Order, p: Payment} / 
    o.user == u && p.order == o && u.active == true ==> 
    Log("Payment " + p.paymentId + " for active user " + u.username)
```

---

## Résumé

| Aspect | Description |
|--------|-------------|
| **Syntaxe** | `variable = Type(...)` |
| **Utilisation** | Référencer dans d'autres faits |
| **Contrainte** | Définir avant utiliser |
| **Avantage** | Lisibilité et réutilisation |

---

## Voir Aussi

- [Comparaisons de Faits](fact-comparisons.md)
- [Système de Types](type-system.md)
- [Identifiants Internes](../internal-ids.md)
```

### 4. Créer Guide de Migration

#### Fichier : `docs/migration/from-v1.x.md` (nouveau)

```markdown
# Guide de Migration v1.x → v2.0

## Vue d'Ensemble

La version 2.0 de TSD introduit des changements majeurs dans la gestion des identifiants.

⚠️ **Breaking Changes** - Ce guide est **obligatoire** pour migrer.

---

## Changements Principaux

### 1. Champ `id` → `_id_` Interne

| Aspect | v1.x | v2.0 |
|--------|------|------|
| **Nom** | `id` (visible) | `_id_` (caché) |
| **Affectation** | Possible | ❌ Interdite |
| **Accès** | Dans expressions | ❌ Interdit |
| **Génération** | Optionnelle | Obligatoire |

### 2. Affectations de Variables

**Nouveau** : Possibilité d'affecter des faits à des variables

```tsd
// v2.0 - NOUVEAU
alice = User("Alice", 30)
Login(alice, "alice@example.com")
```

### 3. Comparaisons de Faits

**Nouveau** : Comparaisons directes de faits

```tsd
// v2.0 - NOUVEAU
{u: User, l: Login} / l.user == u ==> Log("Match")
```

---

## Migration Étape par Étape

### Étape 1 : Supprimer Affectations Manuelles d'ID

**v1.x** :
```tsd
type Person(name: string, age: number)
Person(id: "person_1", name: "Alice", age: 30)
```

**v2.0** :
```tsd
type Person(#name: string, age: number)
alice = Person("Alice", 30)
// ID généré automatiquement : "Person~Alice"
```

**Actions** :
1. Retirer tous les champs `id:` des faits
2. Ajouter `#` aux champs servant d'identifiant naturel
3. Utiliser des affectations pour nommer les faits importants

### Étape 2 : Remplacer Accès à `id`

**v1.x** :
```tsd
{p: Person} / p.id == "person_1" ==> Log("Found")
```

**v2.0** :
```tsd
// Option 1 : Comparer sur le champ
{p: Person} / p.name == "Alice" ==> Log("Found")

// Option 2 : Utiliser une variable
alice = Person("Alice", 30)
{p: Person} / p == alice ==> Log("Found")
```

### Étape 3 : Migrer les Relations

**v1.x** :
```tsd
type User(#email: string, name: string)
type Login(userEmail: string, #sessionId: string)

User(email: "alice@ex.com", name: "Alice")
Login(userEmail: "alice@ex.com", sessionId: "SES-001")

{u: User, l: Login} / l.userEmail == u.email ==> ...
```

**v2.0** :
```tsd
type User(#email: string, name: string)
type Login(user: User, #sessionId: string)

alice = User("alice@ex.com", "Alice")
Login(alice, "SES-001")

{u: User, l: Login} / l.user == u ==> ...
```

**Actions** :
1. Changer le type du champ de relation de `string` à `Type`
2. Utiliser des variables au lieu de dupliquer les valeurs
3. Simplifier les conditions de jointure

### Étape 4 : Adapter les Types

**v1.x** :
```tsd
type Order(orderId: string, userId: string, total: number)
```

**v2.0** :
```tsd
type User(#userId: string, name: string)
type Order(user: User, #orderId: string, total: number)
```

---

## Cas de Migration Courants

### Cas 1 : IDs Séquentiels

**v1.x** :
```tsd
type Entity(name: string)
Entity(id: "1", name: "First")
Entity(id: "2", name: "Second")
```

**v2.0** :
```tsd
type Entity(#entityId: string, name: string)
Entity("1", "First")
Entity("2", "Second")
```

### Cas 2 : Relations N-N

**v1.x** :
```tsd
type Student(#studentId: string, name: string)
type Course(#courseId: string, title: string)
type Enrollment(studentId: string, courseId: string, grade: string)

Student(studentId: "S001", name: "Alice")
Course(courseId: "C001", title: "Math")
Enrollment(studentId: "S001", courseId: "C001", grade: "A")

{s: Student, e: Enrollment, c: Course} / 
    e.studentId == s.studentId && e.courseId == c.courseId ==> ...
```

**v2.0** :
```tsd
type Student(#studentId: string, name: string)
type Course(#courseId: string, title: string)
type Enrollment(student: Student, course: Course, grade: string)

alice = Student("S001", "Alice")
math = Course("C001", "Math")
Enrollment(alice, math, "A")

{s: Student, e: Enrollment, c: Course} / 
    e.student == s && e.course == c ==> ...
```

### Cas 3 : Logs et Événements

**v1.x** :
```tsd
type LogEvent(timestamp: number, level: string, message: string)
LogEvent(id: "log_1", timestamp: 1704067200, level: "INFO", message: "Started")
```

**v2.0** :
```tsd
type LogEvent(timestamp: number, level: string, message: string)
// Pas de clé primaire → hash automatique
LogEvent(1704067200, "INFO", "Started")
// ID: "LogEvent~a1b2c3d4..." (hash)
```

---

## Checklist de Migration

### Préparation

- [ ] Lire ce guide complet
- [ ] Identifier tous les programmes TSD à migrer
- [ ] Sauvegarder les versions actuelles
- [ ] Tester sur un programme simple

### Modifications des Types

- [ ] Ajouter `#` aux champs servant d'identifiant
- [ ] Changer les champs de relation (`string` → `Type`)
- [ ] Retirer les définitions de champ `id`
- [ ] Vérifier qu'aucun champ ne s'appelle `_id_`

### Modifications des Faits

- [ ] Retirer tous les `id:` des faits
- [ ] Créer des affectations pour faits importants
- [ ] Utiliser les variables dans les relations
- [ ] Vérifier l'ordre (définir avant utiliser)

### Modifications des Règles

- [ ] Remplacer accès à `p.id` par `p.naturelKey`
- [ ] Simplifier les jointures (`l.user == u` au lieu de `l.userId == u.userId`)
- [ ] Vérifier qu'aucune règle n'accède à `_id_`

### Tests

- [ ] Parser le programme migré
- [ ] Exécuter et vérifier les résultats
- [ ] Comparer avec le comportement v1.x
- [ ] Tester les cas limites

---

## Outils de Migration

### Script de Vérification

```bash
#!/bin/bash
# Vérifier un fichier TSD pour problèmes de migration

file="$1"

echo "🔍 Vérification de $file"
echo ""

# Chercher 'id:' dans les faits
if grep -q 'id:' "$file"; then
    echo "❌ Affectations manuelles d'ID trouvées"
    grep -n 'id:' "$file"
    echo ""
fi

# Chercher accès à .id
if grep -q '\.id' "$file"; then
    echo "❌ Accès à .id trouvés"
    grep -n '\.id' "$file"
    echo ""
fi

# Chercher _id_
if grep -q '_id_' "$file"; then
    echo "❌ Utilisation de _id_ trouvée"
    grep -n '_id_' "$file"
    echo ""
fi

echo "✅ Vérification terminée"
```

### Validation

```bash
# Valider un programme migré
go run cmd/tsd/main.go validate mon_programme.tsd
```

---

## Exemples de Migration

### Exemple Complet : Système de Blog

**v1.x** :
```tsd
type User(name: string, email: string)
type Post(userId: string, title: string, content: string)

User(id: "user_1", name: "Alice", email: "alice@example.com")
Post(id: "post_1", userId: "user_1", title: "Hello", content: "World")

{u: User, p: Post} / p.userId == u.id ==> Log("Post by " + u.name)
```

**v2.0** :
```tsd
type User(#username: string, email: string)
type Post(author: User, #postId: string, title: string, content: string)

alice = User("alice", "alice@example.com")
Post(alice, "post_1", "Hello", "World")

{u: User, p: Post} / p.author == u ==> Log("Post by " + u.username)
```

---

## FAQ

### Q: Puis-je encore utiliser `id` comme nom de champ ?

**R** : Non, `id` est désormais réservé (en fait, c'est `_id_` qui est réservé, mais `id` est déconseillé pour éviter la confusion).

### Q: Comment référencer un fait spécifique ?

**R** : Utilisez une affectation :
```tsd
alice = User("alice", "alice@example.com")
// Utilisez alice partout où vous avez besoin de cet utilisateur
```

### Q: Les IDs sont-ils encore déterministes ?

**R** : Oui, les IDs sont toujours déterministes, basés sur les clés primaires ou un hash.

### Q: Quelle est la meilleure stratégie de migration ?

**R** :
1. Identifier les "entités" principales
2. Définir des clés primaires naturelles
3. Utiliser des affectations pour ces entités
4. Convertir les relations en références de types

---

## Support

### Ressources

- [Documentation Complète](../README.md)
- [Guide des IDs Internes](../internal-ids.md)
- [Guide des Affectations](../user-guide/fact-assignments.md)
- [Exemples](../../examples/)

### Problèmes Courants

| Erreur | Solution |
|--------|----------|
| `le champ '_id_' est réservé` | Retirer `_id_` de la définition |
| `variable 'x' non définie` | Définir la variable avant de l'utiliser |
| `type 'X' non trouvé` | Vérifier que le type existe |

---

## Conclusion

La migration vers v2.0 nécessite de repenser les relations entre types, mais offre :

✅ **Avantages**
- Syntaxe plus naturelle et lisible
- Relations explicites entre types
- Moins de duplication de données
- Comparaisons simplifiées

⚠️ **Breaking Changes**
- IDs ne sont plus accessibles
- Syntaxe des faits modifiée
- Relations à redéfinir

**Estimation** : 1-4 heures pour un projet moyen

---

**Version** : 2.0.0
**Dernière mise à jour** : 2025-01-XX
```

### 5. Mettre à Jour README Principal

#### Fichier : `README.md` (modifications)

```markdown
# TSD - Règles et Contraintes avec RETE

> Solution générale de synchronisation utilisant un moteur de règles RETE avec système de contraintes en Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/dl/)

---

## 🚀 Nouveautés v2.0

### Affectations de Variables

```tsd
alice = User("Alice", 30)
bob = User("Bob", 25)

Login(alice, "alice@example.com", "secret")
Login(bob, "bob@example.com", "password")
```

### Comparaisons de Faits

```tsd
type User(#name: string, age: number)
type Login(user: User, #email: string, password: string)

{u: User, l: Login} / l.user == u ==> 
    Log("Login for " + u.name)
```

### Types de Faits dans les Champs

```tsd
type Order(customer: Customer, product: Product, quantity: number)
```

**⚠️ Breaking Changes** : Voir [Guide de Migration](docs/migration/from-v1.x.md)

---

## 📚 Documentation

### Guides Utilisateur

- [🎯 Affectations de Faits](docs/user-guide/fact-assignments.md)
- [🔗 Comparaisons de Faits](docs/user-guide/fact-comparisons.md)
- [📐 Système de Types](docs/user-guide/type-system.md)

### Référence Technique

- [🔑 Identifiants Internes](docs/internal-ids.md)
- [📖 Référence Complète](docs/reference.md)
- [🏗️ Architecture](docs/architecture.md)

### Migration

- [📦 Guide de Migration v1.x → v2.0](docs/migration/from-v1.x.md)

### Exemples

- [💡 Exemples Simples](examples/)
- [🎓 Tutoriels](docs/tutorials/)

---

## 🎯 Démarrage Rapide

### Installation

```bash
go get github.com/resinsec/tsd
```

### Exemple Simple

```tsd
// Définir des types
type User(#username: string, email: string, role: string)
type Session(user: User, #sessionId: string, loginTime: number)

// Créer des utilisateurs
admin = User("admin", "admin@example.com", "administrator")
alice = User("alice", "alice@example.com", "user")

// Créer des sessions
Session(admin, "SES-001", 1704067200)
Session(alice, "SES-002", 1704067260)

// Règle : Logger les logins d'admin
{u: User, s: Session} / s.user == u && u.role == "administrator" ==> 
    Log("Admin login: " + u.username + " at " + s.loginTime)
```

### Exécution

```bash
go run cmd/tsd/main.go run mon_programme.tsd
```

---

## 🏗️ Architecture

```
tsd/
├── constraint/        # Parser et validation
├── rete/             # Moteur RETE
├── api/              # API publique
├── tsdio/            # Structures I/O
├── xuples/           # Gestion des tuples
├── docs/             # Documentation centralisée
├── examples/         # Exemples TSD
└── tests/            # Tests (unit, integration, e2e)
```

---

## 📋 Fonctionnalités

- ✅ **Moteur RETE** - Pattern matching efficace
- ✅ **Système de types** - Types utilisateur + primitifs
- ✅ **Clés primaires** - Génération automatique d'IDs
- ✅ **Affectations** - Variables pour réutiliser des faits
- ✅ **Comparaisons de faits** - Relations naturelles
- ✅ **Règles** - Conditions et actions
- ✅ **Validation** - Type-checking statique
- ✅ **API** - Interface Go propre

---

## 🧪 Tests

```bash
# Tests unitaires
make test-unit

# Tests d'intégration
make test-integration

# Tests end-to-end
make test-e2e

# Tous les tests
make test-complete

# Couverture
make test-coverage
```

---

## 📖 Exemples

### Gestion d'Utilisateurs

```tsd
type User(#username: string, email: string, active: bool)
type Login(user: User, #sessionId: string, ipAddress: string)

alice = User("alice", "alice@example.com", true)
bob = User("bob", "bob@example.com", false)

Login(alice, "SES-001", "192.168.1.10")
Login(bob, "SES-002", "192.168.1.11")

{u: User, l: Login} / l.user == u && u.active == false ==> 
    Log("ALERT: Inactive user login attempt: " + u.username)
```

Plus d'exemples dans [examples/](examples/)

---

## 🤝 Contribution

Les contributions sont les bienvenues ! Voir [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📄 Licence

MIT License - voir [LICENSE](LICENSE)

---

## 🔗 Liens

- [Documentation Complète](docs/)
- [Guide de Migration](docs/migration/from-v1.x.md)
- [Exemples](examples/)
- [Issues](https://github.com/resinsec/tsd/issues)

---

**Version actuelle** : 2.0.0
```

### 6. Supprimer Documentation Obsolète

#### Script de Nettoyage

```bash
#!/bin/bash
# Nettoyer la documentation obsolète

echo "🧹 NETTOYAGE DOCUMENTATION OBSOLÈTE"
echo "===================================="
echo ""

# Sauvegarder avant suppression
mkdir -p docs/archive/pre-v2.0
cp -r docs/* docs/archive/pre-v2.0/ 2>/dev/null || true

# Supprimer les anciennes règles d'ID
if [ -f "docs/ID_RULES_COMPLETE.md" ]; then
    echo "📦 Archivage de ID_RULES_COMPLETE.md"
    mv docs/ID_RULES_COMPLETE.md docs/archive/
fi

# Supprimer MIGRATION_IDS.md (remplacé par migration/from-v1.x.md)
if [ -f "docs/MIGRATION_IDS.md" ]; then
    echo "📦 Archivage de MIGRATION_IDS.md"
    mv docs/MIGRATION_IDS.md docs/archive/
fi

# Supprimer primary-keys.md (intégré dans internal-ids.md)
if [ -f "docs/primary-keys.md" ]; then
    echo "📦 Archivage de primary-keys.md"
    mv docs/primary-keys.md docs/archive/
fi

echo ""
echo "✅ Nettoyage terminé"
echo ""
echo "Archive créée dans : docs/archive/pre-v2.0"
```

### 7. Créer Index de Documentation

#### Fichier : `docs/README.md` (mise à jour)

```markdown
# Documentation TSD

Documentation centralisée du projet TSD.

---

## 🎯 Par Où Commencer ?

### Nouveaux Utilisateurs

1. [Démarrage Rapide](../README.md#démarrage-rapide)
2. [Affectations de Faits](user-guide/fact-assignments.md)
3. [Comparaisons de Faits](user-guide/fact-comparisons.md)
4. [Exemples](../examples/)

### Migration depuis v1.x

1. [Guide de Migration](migration/from-v1.x.md) ⚠️ **IMPORTANT**
2. [Nouveautés v2.0](../README.md#nouveautés-v20)
3. [Identifiants Internes](internal-ids.md)

---

## 📚 Documentation Utilisateur

### Guides

- [Affectations de Faits](user-guide/fact-assignments.md)
- [Comparaisons de Faits](user-guide/fact-comparisons.md)
- [Système de Types](user-guide/type-system.md)
- [Clés Primaires](user-guide/primary-keys.md)
- [Règles et Actions](user-guide/rules-and-actions.md)

### Référence

- [Référence Syntaxe TSD](reference.md)
- [Identifiants Internes](internal-ids.md)
- [Types de Données](reference/data-types.md)
- [Fonctions Intégrées](reference/built-in-functions.md)

### Tutoriels

- [Tutoriel 1 : Premier Programme](tutorials/01-first-program.md)
- [Tutoriel 2 : Relations entre Types](tutorials/02-type-relationships.md)
- [Tutoriel 3 : Règles Complexes](tutorials/03-complex-rules.md)

---

## 🔧 Documentation Technique

### Architecture

- [Vue d'Ensemble](architecture.md)
- [Moteur RETE](architecture/rete-engine.md)
- [Génération d'IDs](architecture/id-generation.md)
- [Système de Validation](architecture/validation.md)

### API

- [API Publique](api.md)
- [Package constraint](api/constraint.md)
- [Package rete](api/rete.md)
- [Package tsdio](api/tsdio.md)

---

## 📦 Migration et Mises à Jour

- [Guide de Migration v1.x → v2.0](migration/from-v1.x.md) ⚠️
- [CHANGELOG](../CHANGELOG.md)
- [Breaking Changes](migration/breaking-changes.md)

---

## 💡 Exemples

- [Exemples Simples](../examples/)
- [Cas d'Usage Avancés](examples/)
- [Patterns Courants](patterns/)

---

## 🤝 Contribution

- [Guide de Contribution](../CONTRIBUTING.md)
- [Standards de Code](../.github/prompts/common.md)
- [Standards de Développement](../.github/prompts/develop.md)

---

## 🔍 Index

### Par Fonctionnalité

- **Affectations** : [Guide](user-guide/fact-assignments.md)
- **Comparaisons** : [Guide](user-guide/fact-comparisons.md)
- **Types** : [Guide](user-guide/type-system.md)
- **Identifiants** : [Documentation](internal-ids.md)
- **Clés Primaires** : [Guide](user-guide/primary-keys.md)

### Par Niveau

- **Débutant** : [Démarrage Rapide](../README.md), [Tutoriels](tutorials/)
- **Intermédiaire** : [Guides](user-guide/), [Exemples](../examples/)
- **Avancé** : [Architecture](architecture/), [API](api/)

---

## 📞 Support

- [FAQ](faq.md)
- [Problèmes Courants](troubleshooting.md)
- [Issues GitHub](https://github.com/resinsec/tsd/issues)

---

**Version** : 2.0.0
**Dernière mise à jour** : 2025-01-XX
```

---

## ✅ Critères de Succès

### Documentation

```bash
# Vérifier que tous les fichiers existent
ls -la docs/internal-ids.md
ls -la docs/user-guide/fact-assignments.md
ls -la docs/migration/from-v1.x.md
ls -la docs/README.md

# Vérifier les liens
grep -r "docs/" README.md | grep -v "^#"
```

### Checklist

- [ ] `internal-ids.md` créé
- [ ] `user-guide/fact-assignments.md` créé
- [ ] `user-guide/fact-comparisons.md` créé
- [ ] `migration/from-v1.x.md` créé
- [ ] `README.md` mis à jour
- [ ] `docs/README.md` mis à jour
- [ ] Documentation obsolète archivée
- [ ] Tous les liens fonctionnent
- [ ] Exemples à jour
- [ ] Pas de contradiction

### Validation

```bash
# Vérifier absence de références à l'ancien système
grep -r '"id":' docs/ --include="*.md"
grep -r 'FieldNameID[^I]' docs/ --include="*.md"

# Vérifier cohérence
make validate-docs  # si script existe
```

---

## 📊 Livrables

### Nouveaux Fichiers

- [ ] `docs/internal-ids.md`
- [ ] `docs/user-guide/fact-assignments.md`
- [ ] `docs/user-guide/fact-comparisons.md`
- [ ] `docs/user-guide/type-system.md`
- [ ] `docs/migration/from-v1.x.md`
- [ ] `docs/migration/breaking-changes.md`

### Fichiers Mis à Jour

- [ ] `README.md`
- [ ] `docs/README.md`
- [ ] `docs/reference.md`
- [ ] `docs/api.md`
- [ ] Module READMEs

### Fichiers Archivés

- [ ] `docs/archive/ID_RULES_COMPLETE.md`
- [ ] `docs/archive/MIGRATION_IDS.md`
- [ ] `docs/archive/primary-keys.md`

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Analyser documentation existante
2. ✅ Créer `internal-ids.md`
3. ✅ Créer guides utilisateur
4. ✅ Créer guide de migration
5. ✅ Mettre à jour README
6. ✅ Créer index documentation
7. ✅ Archiver obsolète
8. ✅ Vérifier liens
9. ✅ Validation finale

### Commandes

```bash
# Créer répertoires
mkdir -p docs/user-guide
mkdir -p docs/migration
mkdir -p docs/archive
mkdir -p docs/architecture
mkdir -p docs/api
mkdir -p docs/tutorials

# Créer les fichiers
touch docs/internal-ids.md
touch docs/user-guide/fact-assignments.md
touch docs/user-guide/fact-comparisons.md
touch docs/migration/from-v1.x.md

# Archiver ancien
mkdir -p docs/archive/pre-v2.0
mv docs/ID_RULES_COMPLETE.md docs/archive/ 2>/dev/null || true
mv docs/MIGRATION_IDS.md docs/archive/ 2>/dev/null || true

# Vérifier
find docs/ -name "*.md" | sort
```

---

## 📚 Références

- `scripts/new_ids/08-prompt-tests-integration.md` - Tests E2E
- `scripts/new_ids/07-prompt-tests-unit.md` - Tests unitaires
- `docs/` - Documentation actuelle
- `examples/` - Exemples actuels

---

## 📝 Notes

### Points d'Attention

1. **Cohérence** : Toute la documentation doit être cohérente

2. **Clarté** : Exemples simples et clairs

3. **Exhaustivité** : Couvrir tous les cas d'usage

4. **Migration** : Guide complet et détaillé obligatoire

### Bonnes Pratiques

```markdown
<!-- ✅ BON - Exemple clair avec résultat -->
```tsd
alice = User("Alice", 30)
Login(alice, "alice@example.com")
// ID généré: "User~Alice"
```

<!-- ❌ MAUVAIS - Exemple sans contexte -->
```tsd
alice = User("Alice", 30)
```
```

---

## 🎯 Résultat Attendu

Après ce prompt :

```
docs/
├── README.md                          # ✅ Index complet
├── internal-ids.md                    # ✅ Documentation IDs
├── reference.md                       # ✅ Référence complète
├── architecture.md                    # ✅ Architecture
├── api.md                            # ✅ API
├── user-guide/                       # ✅ Guides utilisateur
│   ├── fact-assignments.md
│   ├── fact-comparisons.md
│   ├── type-system.md
│   └── primary-keys.md
├── migration/                        # ✅ Guides migration
│   ├── from-v1.x.md
│   └── breaking-changes.md
├── tutorials/                        # ✅ Tutoriels
├── examples/                         # ✅ Exemples
├── architecture/                     # ✅ Architecture détaillée
├── api/                             # ✅ Documentation API
└── archive/                         # ✅ Archives
    └── pre-v2.0/
```

---

**Prompt suivant** : `10-prompt-finalisation.md`

**Durée estimée** : 6-8 heures

**Complexité** : ⚠️ Moyenne (beaucoup de rédaction)