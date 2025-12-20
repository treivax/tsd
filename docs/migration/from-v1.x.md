# Guide de Migration v1.x → v2.0

## Vue d'Ensemble

La version 2.0 de TSD introduit des changements majeurs dans la gestion des identifiants et des relations entre faits.

⚠️ **Breaking Changes** - Ce guide est **obligatoire** pour migrer depuis v1.x.

---

## 🚨 Changements Principaux

### 1. Champ `id` → `_id_` Interne

| Aspect | v1.x | v2.0 |
|--------|------|------|
| **Nom du champ** | `id` (visible) | `_id_` (caché) |
| **Affectation manuelle** | ✅ Possible | ❌ **INTERDITE** |
| **Accès dans expressions** | ✅ Accessible | ❌ **INTERDIT** |
| **Génération** | Optionnelle | ✅ Obligatoire et automatique |
| **Utilisation** | Visible partout | Interne uniquement |

**Impact** : Le champ `_id_` est désormais **strictement réservé au système** et **jamais accessible** dans les expressions TSD.

### 2. Affectations de Variables (NOUVEAU)

**Nouveau en v2.0** : Possibilité d'affecter des faits à des variables pour les réutiliser.

```tsd
// v2.0 - NOUVEAU
alice = User("Alice", "alice@example.com", 30)
order1 = Order(alice, "ORD-001", 150.00)
```

### 3. Comparaisons de Faits (NOUVEAU)

**Nouveau en v2.0** : Comparaisons directes entre faits dans les règles.

```tsd
// v2.0 - NOUVEAU
{u: User, o: Order} / o.customer == u ==> 
    Log("Order for user: " + u.username)
```

### 4. Types de Faits dans les Champs (NOUVEAU)

**Nouveau en v2.0** : Les champs peuvent être d'un type de fait (pas seulement primitifs).

```tsd
// v2.0 - NOUVEAU
type Order(customer: Customer, #orderNumber: string, total: number)
```

---

## 📋 Migration Étape par Étape

### Étape 1 : Supprimer les Affectations Manuelles d'ID

#### ❌ v1.x (OBSOLÈTE)

```tsd
type Person(name: string, age: number)
assert Person(id: "person_1", name: "Alice", age: 30)

rule findPerson : {p: Person} / p.id == "person_1"
    ==> print("Found: " + p.name)
```

**Problèmes** :
- Affectation manuelle de `id: "person_1"` (interdit en v2.0)
- Accès à `p.id` dans la condition (interdit en v2.0)

#### ✅ v2.0 (CORRECT)

```tsd
type Person(#name: string, age: number)
alice = Person("Alice", 30)
// ID généré automatiquement en interne: "Person~Alice"

rule findPerson : {p: Person} / p.name == "Alice"
    ==> print("Found: " + p.name)
```

**Solutions** :
1. ✅ Ajouter `#` devant `name` pour en faire une clé primaire
2. ✅ Utiliser une affectation `alice = Person(...)`
3. ✅ Retirer `id:` de l'assertion
4. ✅ Comparer sur `p.name` au lieu de `p.id`

---

### Étape 2 : Convertir les Relations entre Faits

#### ❌ v1.x (OBSOLÈTE)

```tsd
type User(#email: string, name: string)
type Login(userEmail: string, #sessionId: string, timestamp: number)

assert User(email: "alice@ex.com", name: "Alice")
assert Login(userEmail: "alice@ex.com", sessionId: "SES-001", timestamp: 1704067200)

rule userLogin : {u: User, l: Login} / l.userEmail == u.email
    ==> print("Login for: " + u.name)
```

**Problèmes** :
- Duplication de `alice@ex.com`
- Relation via string (pas type-safe)
- Condition de jointure manuelle

#### ✅ v2.0 (CORRECT)

```tsd
type User(#email: string, name: string)
type Login(user: User, #sessionId: string, timestamp: number)

alice = User("alice@ex.com", "Alice")
Login(alice, "SES-001", 1704067200)

rule userLogin : {u: User, l: Login} / l.user == u
    ==> print("Login for: " + u.name)
```

**Solutions** :
1. ✅ Changer `userEmail: string` en `user: User`
2. ✅ Utiliser une variable `alice` pour éviter duplication
3. ✅ Passer `alice` directement au lieu de dupliquer l'email
4. ✅ Simplifier la condition : `l.user == u` au lieu de `l.userEmail == u.email`

---

### Étape 3 : Adapter les Définitions de Types

#### ❌ v1.x (OBSOLÈTE)

```tsd
type Order(orderId: string, customerId: string, total: number)
type OrderLine(orderId: string, productId: string, quantity: number)

assert Order(id: "1", orderId: "ORD-001", customerId: "CUST-001", total: 150.00)
assert OrderLine(id: "2", orderId: "ORD-001", productId: "PROD-123", quantity: 2)

rule orderDetails : {o: Order, ol: OrderLine} / ol.orderId == o.orderId
    ==> print("Order " + o.orderId + " has product " + ol.productId)
```

**Problèmes** :
- IDs manuels (`id: "1"`, `id: "2"`)
- Relations via strings (duplication)
- Pas de clés primaires définies

#### ✅ v2.0 (CORRECT)

```tsd
type Customer(#customerId: string, name: string)
type Product(#productId: string, name: string, price: number)
type Order(customer: Customer, #orderNumber: string, total: number)
type OrderLine(order: Order, product: Product, quantity: number)

cust1 = Customer("CUST-001", "Alice Johnson")
prod1 = Product("PROD-123", "Laptop", 1200.00)

order1 = Order(cust1, "ORD-001", 150.00)
OrderLine(order1, prod1, 2)

rule orderDetails : {o: Order, ol: OrderLine} / ol.order == o
    ==> print("Order " + o.orderNumber + " has " + ol.quantity + " " + ol.product.name)
```

**Solutions** :
1. ✅ Définir types manquants (`Customer`, `Product`)
2. ✅ Ajouter clés primaires avec `#`
3. ✅ Utiliser types de faits pour les relations
4. ✅ Utiliser affectations pour éviter duplication
5. ✅ Simplifier les conditions de jointure

---

## 🔄 Cas de Migration Courants

### Cas 1 : IDs Séquentiels

#### ❌ v1.x

```tsd
type Entity(name: string, description: string)
assert Entity(id: "1", name: "First", description: "The first entity")
assert Entity(id: "2", name: "Second", description: "The second entity")
```

#### ✅ v2.0

```tsd
type Entity(#entityId: string, name: string, description: string)
Entity("1", "First", "The first entity")
Entity("2", "Second", "The second entity")
```

**Note** : Si les IDs séquentiels sont importants, créez un champ explicite `#entityId`.

---

### Cas 2 : Relations Many-to-Many

#### ❌ v1.x

```tsd
type Student(#studentId: string, name: string)
type Course(#courseId: string, title: string)
type Enrollment(studentId: string, courseId: string, grade: string)

assert Student(studentId: "S001", name: "Alice")
assert Course(courseId: "C001", title: "Math 101")
assert Enrollment(id: "E001", studentId: "S001", courseId: "C001", grade: "A")

rule studentGrades : {s: Student, e: Enrollment, c: Course} / 
    e.studentId == s.studentId && e.courseId == c.courseId
    ==> print(s.name + " got " + e.grade + " in " + c.title)
```

#### ✅ v2.0

```tsd
type Student(#studentId: string, name: string)
type Course(#courseId: string, title: string)
type Enrollment(student: Student, course: Course, grade: string)

alice = Student("S001", "Alice")
math = Course("C001", "Math 101")
Enrollment(alice, math, "A")

rule studentGrades : {s: Student, e: Enrollment, c: Course} / 
    e.student == s && e.course == c
    ==> print(s.name + " got " + e.grade + " in " + c.title)
```

**Avantages** :
- ✅ Pas de duplication d'IDs
- ✅ Relations type-safe
- ✅ Conditions simplifiées
- ✅ Clé primaire composite automatique pour `Enrollment`

---

### Cas 3 : Logs et Événements (Sans Clé Primaire)

#### ❌ v1.x

```tsd
type LogEvent(timestamp: number, level: string, message: string)
assert LogEvent(id: "log_1", timestamp: 1704067200, level: "INFO", message: "Started")
assert LogEvent(id: "log_2", timestamp: 1704067260, level: "ERROR", message: "Failed")
```

#### ✅ v2.0

```tsd
type LogEvent(timestamp: number, level: string, message: string)
// Pas de clé primaire → hash automatique
LogEvent(1704067200, "INFO", "Application started")
LogEvent(1704067260, "ERROR", "Operation failed")
// IDs générés: "LogEvent~<hash1>", "LogEvent~<hash2>"
```

**Note** : Pour les événements/logs, un hash est souvent approprié car il n'y a pas d'identifiant naturel.

---

### Cas 4 : Accès aux IDs dans les Actions

#### ❌ v1.x

```tsd
type User(#username: string, email: string)
assert User(username: "alice", email: "alice@example.com")

rule showUser : {u: User} / u.username == "alice"
    ==> print("User ID: " + u.id + ", Name: " + u.username)
```

#### ✅ v2.0

```tsd
type User(#username: string, email: string)
alice = User("alice", "alice@example.com")

rule showUser : {u: User} / u.username == "alice"
    ==> print("Username: " + u.username + ", Email: " + u.email)
```

**Solution** :
- ❌ NE PAS accéder à `u.id` ou `u._id_`
- ✅ Utiliser les champs métier (`u.username`, `u.email`)
- ℹ️ L'ID interne existe mais est caché

---

## ✅ Checklist de Migration

### Préparation

- [ ] Lire ce guide complet
- [ ] Identifier tous les programmes TSD à migrer
- [ ] Sauvegarder les versions actuelles (git commit/tag)
- [ ] Tester sur un programme simple d'abord

### Modifications des Types

- [ ] Identifier les identifiants naturels pour chaque type
- [ ] Ajouter `#` devant les champs servant de clé primaire
- [ ] Changer les champs de relation de `string` vers type de fait
- [ ] Retirer toute définition de champ nommé `id` ou `_id_`
- [ ] Vérifier qu'aucun champ ne s'appelle `_id_`

### Modifications des Faits

- [ ] Retirer tous les `id: "..."` des assertions
- [ ] Créer des affectations pour les faits importants
- [ ] Utiliser les variables dans les relations
- [ ] Vérifier l'ordre (définir avant utiliser)
- [ ] Simplifier en évitant la duplication

### Modifications des Règles

- [ ] Remplacer accès à `p.id` par accès aux champs métier
- [ ] Simplifier les jointures (utiliser `==` entre faits)
- [ ] Vérifier qu'aucune règle n'accède à `_id_`
- [ ] Tester que les règles matchent toujours correctement

### Tests et Validation

- [ ] Parser le programme migré sans erreur
- [ ] Exécuter et vérifier les résultats
- [ ] Comparer avec le comportement v1.x
- [ ] Tester les cas limites
- [ ] Vérifier les performances (benchmarks si nécessaire)

---

## 🛠️ Outils de Migration

### Script de Vérification

Créez ce script `check_migration.sh` :

```bash
#!/bin/bash
# Vérifier un fichier TSD pour problèmes de migration

file="$1"

if [ -z "$file" ]; then
    echo "Usage: $0 <fichier.tsd>"
    exit 1
fi

echo "🔍 Vérification de $file"
echo "========================"
echo ""

has_errors=0

# Chercher 'id:' dans les faits
if grep -n 'id:' "$file" | grep -v '^[[:space:]]*//'; then
    echo "❌ Affectations manuelles d'ID trouvées (lignes ci-dessus)"
    has_errors=1
    echo ""
fi

# Chercher accès à .id (mais pas .identifier, etc.)
if grep -n '\\.id[^a-zA-Z]' "$file" | grep -v '^[[:space:]]*//'; then
    echo "❌ Accès à .id trouvés (lignes ci-dessus)"
    has_errors=1
    echo ""
fi

# Chercher _id_
if grep -n '_id_' "$file" | grep -v '^[[:space:]]*//'; then
    echo "❌ Utilisation de _id_ trouvée (lignes ci-dessus)"
    has_errors=1
    echo ""
fi

# Chercher assert avec 'id:' en paramètre
if grep -n 'assert.*id:' "$file" | grep -v '^[[:space:]]*//'; then
    echo "❌ Assertions avec 'id:' trouvées (lignes ci-dessus)"
    has_errors=1
    echo ""
fi

if [ $has_errors -eq 0 ]; then
    echo "✅ Aucun problème détecté"
else
    echo ""
    echo "⚠️  Des problèmes ont été détectés - migration nécessaire"
    exit 1
fi
```

Utilisation :
```bash
chmod +x check_migration.sh
./check_migration.sh mon_programme.tsd
```

### Validation avec TSD

```bash
# Valider un programme migré
tsd validate mon_programme.tsd

# Parser et afficher les types
tsd parse mon_programme.tsd --show-types

# Exécuter le programme
tsd run mon_programme.tsd
```

---

## 📖 Exemples de Migration Complets

### Exemple 1 : Système de Blog

#### ❌ v1.x

```tsd
type User(name: string, email: string)
type Post(userId: string, title: string, content: string)
type Comment(postId: string, authorId: string, text: string)

assert User(id: "user_1", name: "Alice", email: "alice@example.com")
assert User(id: "user_2", name: "Bob", email: "bob@example.com")

assert Post(id: "post_1", userId: "user_1", title: "Hello World", content: "First post!")
assert Comment(id: "com_1", postId: "post_1", authorId: "user_2", text: "Great post!")

rule postComments : {p: Post, c: Comment, u: User} / 
    c.postId == p.id && c.authorId == u.id
    ==> print(u.name + " commented on: " + p.title)
```

#### ✅ v2.0

```tsd
type User(#username: string, email: string)
type Post(author: User, #postId: string, title: string, content: string)
type Comment(post: Post, author: User, #commentId: string, text: string)

alice = User("alice", "alice@example.com")
bob = User("bob", "bob@example.com")

post1 = Post(alice, "post_1", "Hello World", "First post!")
Comment(post1, bob, "com_1", "Great post!")

rule postComments : {p: Post, c: Comment, u: User} / 
    c.post == p && c.author == u
    ==> print(u.username + " commented on: " + p.title)
```

**Changements** :
1. Types référencent des faits (`author: User` au lieu de `userId: string`)
2. Affectations de variables (`alice`, `bob`, `post1`)
3. Pas d'IDs manuels
4. Conditions simplifiées (`c.post == p` au lieu de `c.postId == p.id`)

---

### Exemple 2 : E-Commerce

#### ❌ v1.x

```tsd
type Customer(customerId: string, name: string, vip: boolean)
type Product(sku: string, name: string, price: number)
type Order(orderId: string, customerId: string, total: number)
type OrderLine(orderId: string, productSku: string, quantity: number)

assert Customer(id: "1", customerId: "CUST-001", name: "Alice", vip: true)
assert Product(id: "2", sku: "LAPTOP-001", name: "Gaming Laptop", price: 1200.00)

assert Order(id: "3", orderId: "ORD-001", customerId: "CUST-001", total: 2400.00)
assert OrderLine(id: "4", orderId: "ORD-001", productSku: "LAPTOP-001", quantity: 2)

rule vipOrders : {c: Customer, o: Order, ol: OrderLine, p: Product} /
    o.customerId == c.customerId && 
    ol.orderId == o.orderId && 
    ol.productSku == p.sku &&
    c.vip == true
    ==> print("VIP order: " + o.orderId + " with " + p.name)
```

#### ✅ v2.0

```tsd
type Customer(#customerId: string, name: string, vip: boolean)
type Product(#sku: string, name: string, price: number)
type Order(customer: Customer, #orderNumber: string, total: number)
type OrderLine(order: Order, product: Product, quantity: number)

alice = Customer("CUST-001", "Alice Johnson", true)
laptop = Product("LAPTOP-001", "Gaming Laptop", 1200.00)

order1 = Order(alice, "ORD-001", 2400.00)
OrderLine(order1, laptop, 2)

rule vipOrders : {c: Customer, o: Order, ol: OrderLine, p: Product} /
    o.customer == c && ol.order == o && ol.product == p && c.vip == true
    ==> print("VIP order: " + o.orderNumber + " - " + ol.quantity + "x " + p.name)
```

**Avantages** :
- 📉 Code réduit de ~30%
- ✅ Relations explicites et type-safe
- ✅ Pas de duplication d'IDs
- ✅ Conditions plus lisibles
- ✅ Moins d'erreurs possibles

---

## ❓ FAQ Migration

### Q1 : Puis-je encore utiliser `id` comme nom de champ ?

**R** : ❌ Non. Le nom `id` est réservé (en fait `_id_` mais `id` est déconseillé pour éviter confusion).

Si vous avez besoin d'un identifiant métier, utilisez un nom explicite :
```tsd
// ❌ INTERDIT
type User(id: string, name: string)

// ✅ BON
type User(#userId: string, name: string)
```

### Q2 : Comment référencer un fait spécifique si je ne peux plus utiliser `id` ?

**R** : Utilisez une **affectation de variable** :

```tsd
// ✅ Affecter le fait à une variable
alice = User("alice", "alice@example.com")

// ✅ Réutiliser la variable
Login(alice, "SES-001")
Order(alice, "ORD-001", 150.00)

// ✅ Dans les règles
{u: User, o: Order} / o.customer == u && u.username == "alice"
    ==> print("Found Alice's order")
```

### Q3 : Les IDs sont-ils encore déterministes ?

**R** : ✅ Oui, absolument. Les IDs sont toujours déterministes :
- Avec clés primaires : basés sur les valeurs des clés
- Sans clé primaire : basés sur un hash déterministe de tous les champs

```tsd
type User(#username: string, email: string)

// Ces deux assertions génèrent le même ID interne
alice1 = User("alice", "alice@example.com")
alice2 = User("alice", "alice@example.com")
// Les deux ont l'ID interne: "User~alice"
```

### Q4 : Que se passe-t-il si j'essaie d'accéder à `_id_` ?

**R** : ❌ Erreur de validation. Le parser rejettera le programme :

```tsd
// ❌ ERREUR
{u: User} / u._id_ == "something"
    ==> print("...")

// Erreur: "Le champ '_id_' est réservé et ne peut pas être utilisé dans les expressions"
```

### Q5 : Comment migrer un gros projet ?

**R** : Approche incrémentale recommandée :

1. **Phase 1** : Migrer les types
   - Ajouter clés primaires
   - Changer relations string → types

2. **Phase 2** : Migrer les faits
   - Retirer IDs manuels
   - Ajouter affectations

3. **Phase 3** : Migrer les règles
   - Adapter conditions
   - Retirer accès à `.id`

4. **Phase 4** : Tester
   - Tests unitaires
   - Tests d'intégration
   - Validation complète

### Q6 : Y a-t-il un impact sur les performances ?

**R** : ✅ Généralement **meilleures performances** :
- Génération d'IDs optimisée
- Moins de duplications en mémoire
- Comparaisons de faits plus efficaces

Benchmarks disponibles dans `rete/id_formats_benchmark_test.go`.

### Q7 : Quelle est la meilleure stratégie de clés primaires ?

**R** : Règles générales :

1. **Clé naturelle unique** : Privilégier
   ```tsd
   type User(#username: string, ...)  // ✅ BON
   ```

2. **Clé composite** : Si nécessaire
   ```tsd
   type OrderLine(#orderId: string, #productId: string, ...)  // ✅ OK
   ```

3. **Pas de clé** : Pour événements/logs
   ```tsd
   type LogEvent(timestamp: number, ...)  // ✅ Hash auto
   ```

4. **Clé stable** : Ne change pas
   ```tsd
   type User(#username: string, ...)  // ✅ username stable
   type User(#email: string, ...)     // ⚠️ email peut changer
   ```

---

## 📚 Ressources Complémentaires

### Documentation

- [Identifiants Internes](../internal-ids.md) - Documentation complète du système `_id_`
- [Affectations de Faits](../user-guide/fact-assignments.md) - Guide des affectations
- [Comparaisons de Faits](../user-guide/fact-comparisons.md) - Guide des comparaisons
- [Système de Types](../user-guide/type-system.md) - Types de faits dans les champs
- [README Principal](../../README.md) - Vue d'ensemble du projet

### Support

- [Issues GitHub](https://github.com/chrlesur/tsd/issues) - Rapporter des problèmes
- [Documentation Complète](../README.md) - Index de toute la documentation
- [Exemples](../../examples/) - Exemples de programmes TSD

---

## 🎯 Résumé

### Points Clés

1. ⚠️ **`_id_` est caché** - Jamais accessible dans les expressions
2. ✅ **Affectations** - Utiliser `variable = Type(...)` pour réutiliser
3. ✅ **Comparaisons** - `fact1 == fact2` directement
4. ✅ **Types dans champs** - Relations type-safe
5. ✅ **Clés primaires** - Utiliser `#` pour identifiants naturels

### Estimation Migration

| Taille Projet | Durée Estimée |
|---------------|---------------|
| Petit (< 500 lignes) | 1-2 heures |
| Moyen (500-2000 lignes) | 4-8 heures |
| Grand (> 2000 lignes) | 1-3 jours |

### Bénéfices

✅ **Code plus propre** - Moins de duplication  
✅ **Type safety** - Relations explicites  
✅ **Maintenabilité** - Code plus lisible  
✅ **Performance** - Optimisations internes  
✅ **Fiabilité** - Moins d'erreurs possibles  

---

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-19  
**Auteur** : Équipe TSD

---

*Ce guide est en constante évolution. N'hésitez pas à contribuer via pull request.*
