# Tutoriel TSD : De zéro à héros

Ce tutoriel vous guide pas à pas dans l'apprentissage du langage TSD, du plus simple au plus avancé.

## Table des matières

1. [Premiers pas](#premiers-pas)
2. [Règles simples](#règles-simples)
3. [Jointures et relations](#jointures-et-relations)
4. [Conditions avancées](#conditions-avancées)
5. [Agrégations](#agrégations)
6. [Projet complet](#projet-complet)

---

## Premiers pas

### Étape 1 : Votre premier type

Commençons par créer un type simple représentant une personne :

```tsd
type Person : <name: string, age: number>
```

**Explication :**
- `type Person` : Déclare un nouveau type nommé "Person"
- `: <...>` : Définit les champs du type
- `name: string` : Champ "name" de type chaîne de caractères
- `age: number` : Champ "age" de type nombre

### Étape 2 : Créer des faits

Maintenant, créons quelques instances (faits) de ce type :

```tsd
Person(name="Alice", age=25)
Person(name="Bob", age=17)
Person(name="Charlie", age=30)
```

**Explication :**
- `Person(...)` : Instancie le type Person
- `name="Alice"` : Assigne la valeur "Alice" au champ name
- Les valeurs doivent correspondre aux types définis

### Étape 3 : Votre première règle

Créons une règle qui identifie les adultes (18 ans ou plus) :

```tsd
{p: Person} / p.age >= 18 ==> adult(p.name)
```

**Explication :**
- `{p: Person}` : Déclare une variable "p" de type Person
- `/` : Sépare les variables des conditions
- `p.age >= 18` : Condition à vérifier
- `==>` : Si la condition est vraie, exécuter l'action
- `adult(p.name)` : Appelle la fonction "adult" avec le nom de la personne

**Résultat :** Cette règle déclenchera `adult("Alice")` et `adult("Charlie")`, mais pas pour Bob (17 ans).

### Programme complet

```tsd
# Définition du type
type Person : <name: string, age: number>

# Règle
{p: Person} / p.age >= 18 ==> adult(p.name)

# Faits
Person(name="Alice", age=25)
Person(name="Bob", age=17)
Person(name="Charlie", age=30)
```

---

## Règles simples

### Exemple 1 : Validation de prix

```tsd
type Product : <id: string, name: string, price: number, available: bool>

# Règle : Produits chers (> 100€)
{p: Product} / p.price > 100 ==> expensive_product(p.id, p.name)

# Règle : Produits disponibles et abordables
{p: Product} / p.available == true AND p.price <= 50 ==> affordable_option(p.id)

# Faits
Product(id="P1", name="Laptop", price=999, available=true)
Product(id="P2", name="Mouse", price=25, available=true)
Product(id="P3", name="Monitor", price=300, available=false)
```

**Résultats :**
- `expensive_product("P1", "Laptop")` ✓
- `expensive_product("P3", "Monitor")` ✓
- `affordable_option("P2")` ✓

### Exemple 2 : Utilisation de fonctions

```tsd
type User : <username: string, email: string, password: string>

# Règle : Mot de passe trop court
{u: User} / LENGTH(u.password) < 8 ==> weak_password_alert(u.username)

# Règle : Email valide
{u: User} / u.email CONTAINS "@" AND u.email CONTAINS "." ==> valid_email(u.username)

# Faits
User(username="alice", email="alice@example.com", password="secret")
User(username="bob", email="bob@mail.co", password="verysecurepassword123")
```

**Résultats :**
- `weak_password_alert("alice")` ✓ (6 caractères)
- `valid_email("alice")` ✓
- `valid_email("bob")` ✓

### Exemple 3 : Opérateur IN

```tsd
type Order : <id: string, status: string, priority: string>

# Règle : Commandes nécessitant un suivi
{o: Order} / 
    o.status IN ["pending", "processing", "shipped"]
    ==> track_order(o.id, o.status)

# Règle : Commandes prioritaires
{o: Order} /
    o.priority IN ["urgent", "high"] AND o.status == "pending"
    ==> prioritize_order(o.id)

# Faits
Order(id="O1", status="pending", priority="urgent")
Order(id="O2", status="delivered", priority="normal")
Order(id="O3", status="processing", priority="high")
```

---

## Jointures et relations

### Exemple 1 : Relation simple

```tsd
type User : <id: string, name: string>
type Order : <id: string, user_id: string, total: number>

# Règle : Associer utilisateurs et commandes
{u: User, o: Order} /
    u.id == o.user_id
    ==> user_order(u.name, o.id, o.total)

# Faits
User(id="U1", name="Alice")
User(id="U2", name="Bob")

Order(id="O1", user_id="U1", total=150)
Order(id="O2", user_id="U1", total=200)
Order(id="O3", user_id="U2", total=75)
```

**Résultats :**
- `user_order("Alice", "O1", 150)` ✓
- `user_order("Alice", "O2", 200)` ✓
- `user_order("Bob", "O3", 75)` ✓

### Exemple 2 : Jointure triple

```tsd
type Customer : <id: string, name: string, vip: bool>
type Order : <id: string, customer_id: string, product_id: string>
type Product : <id: string, name: string, price: number>

# Règle : Clients VIP achetant des produits chers
{c: Customer, o: Order, p: Product} /
    c.id == o.customer_id AND
    o.product_id == p.id AND
    c.vip == true AND
    p.price > 500
    ==> vip_expensive_purchase(c.name, p.name, p.price)

# Faits
Customer(id="C1", name="Alice", vip=true)
Customer(id="C2", name="Bob", vip=false)

Order(id="O1", customer_id="C1", product_id="P1")
Order(id="O2", customer_id="C2", product_id="P2")

Product(id="P1", name="Laptop", price=999)
Product(id="P2", name="Mouse", price=25)
```

**Résultat :**
- `vip_expensive_purchase("Alice", "Laptop", 999)` ✓

### Exemple 3 : Calculs dans les jointures

```tsd
type Item : <id: string, base_price: number>
type Discount : <item_id: string, percentage: number>

# Règle : Calculer le prix final avec réduction
{i: Item, d: Discount} /
    i.id == d.item_id
    ==> final_price(i.id, i.base_price * (1 - d.percentage / 100))

# Faits
Item(id="I1", base_price=100)
Item(id="I2", base_price=200)

Discount(item_id="I1", percentage=20)
Discount(item_id="I2", percentage=10)
```

**Résultats :**
- `final_price("I1", 80)` ✓ (100 * 0.8)
- `final_price("I2", 180)` ✓ (200 * 0.9)

---

## Conditions avancées

### NOT - Négation

```tsd
type Employee : <id: string, name: string, on_leave: bool, suspended: bool>

# Règle : Employés actifs (ni en congé, ni suspendus)
{e: Employee} /
    NOT (e.on_leave == true OR e.suspended == true)
    ==> active_employee(e.name)

# Faits
Employee(id="E1", name="Alice", on_leave=false, suspended=false)
Employee(id="E2", name="Bob", on_leave=true, suspended=false)
Employee(id="E3", name="Charlie", on_leave=false, suspended=true)
```

**Résultat :**
- `active_employee("Alice")` ✓ seulement

### EXISTS - Quantification

```tsd
type Customer : <id: string, name: string>
type Purchase : <customer_id: string, amount: number, completed: bool>

# Règle : Clients ayant au moins un achat complété
{c: Customer} /
    EXISTS (p: Purchase / p.customer_id == c.id AND p.completed == true)
    ==> active_customer(c.name)

# Règle : Clients ayant un gros achat (> 500)
{c: Customer} /
    EXISTS (p: Purchase / p.customer_id == c.id AND p.amount > 500)
    ==> high_value_customer(c.name)

# Faits
Customer(id="C1", name="Alice")
Customer(id="C2", name="Bob")

Purchase(customer_id="C1", amount=150, completed=true)
Purchase(customer_id="C1", amount=600, completed=false)
Purchase(customer_id="C2", amount=50, completed=false)
```

**Résultats :**
- `active_customer("Alice")` ✓
- `high_value_customer("Alice")` ✓ (600 > 500 même si non complété)

### Combinaisons complexes

```tsd
type Account : <id: string, balance: number, frozen: bool, verified: bool>

# Règle : Comptes éligibles aux transferts
{a: Account} /
    a.balance > 100 AND
    a.verified == true AND
    NOT (a.frozen == true)
    ==> can_transfer(a.id)

# Règle : Comptes nécessitant une vérification
{a: Account} /
    (a.balance > 10000 OR a.frozen == true) AND
    a.verified == false
    ==> needs_verification(a.id)

# Faits
Account(id="A1", balance=5000, frozen=false, verified=true)
Account(id="A2", balance=15000, frozen=false, verified=false)
Account(id="A3", balance=200, frozen=true, verified=true)
```

**Résultats :**
- `can_transfer("A1")` ✓
- `needs_verification("A2")` ✓

---

## Agrégations

### COUNT - Compter

```tsd
type Author : <id: string, name: string>
type Book : <id: string, author_id: string, published: bool>

# Règle : Auteurs prolifiques (> 3 livres)
{a: Author} /
    COUNT(b: Book / b.author_id == a.id) > 3
    ==> prolific_author(a.name)

# Règle : Auteurs avec au moins 2 livres publiés
{a: Author} /
    COUNT(b: Book / b.author_id == a.id AND b.published == true) >= 2
    ==> published_author(a.name)

# Faits
Author(id="A1", name="Alice Smith")
Author(id="A2", name="Bob Jones")

Book(id="B1", author_id="A1", published=true)
Book(id="B2", author_id="A1", published=true)
Book(id="B3", author_id="A1", published=false)
Book(id="B4", author_id="A1", published=true)
Book(id="B5", author_id="A2", published=true)
```

**Résultats :**
- `prolific_author("Alice Smith")` ✓ (4 livres)
- `published_author("Alice Smith")` ✓ (3 publiés)

### SUM - Somme

```tsd
type Customer : <id: string, name: string>
type Invoice : <id: string, customer_id: string, amount: number, paid: bool>

# Règle : Clients avec dette > 500
{c: Customer} /
    SUM(i: Invoice / i.customer_id == c.id AND i.paid == false; i.amount) > 500
    ==> high_debt_alert(c.name)

# Règle : Clients fidèles (total payé > 2000)
{c: Customer} /
    SUM(i: Invoice / i.customer_id == c.id AND i.paid == true; i.amount) > 2000
    ==> loyal_customer(c.name)

# Faits
Customer(id="C1", name="Alice")
Customer(id="C2", name="Bob")

Invoice(id="I1", customer_id="C1", amount=300, paid=false)
Invoice(id="I2", customer_id="C1", amount=400, paid=false)
Invoice(id="I3", customer_id="C1", amount=1500, paid=true)
Invoice(id="I4", customer_id="C2", amount=100, paid=false)
```

**Résultats :**
- `high_debt_alert("Alice")` ✓ (300 + 400 = 700)
- `loyal_customer("Alice")` ✓ (1500 payé)

### AVG - Moyenne

```tsd
type Product : <id: string, name: string>
type Review : <product_id: string, rating: number, verified: bool>

# Règle : Produits très bien notés (moyenne >= 4.5)
{p: Product} /
    AVG(r: Review / r.product_id == p.id; r.rating) >= 4.5
    ==> highly_rated(p.name)

# Règle : Produits avec avis vérifiés positifs (moyenne >= 4.0)
{p: Product} /
    AVG(r: Review / r.product_id == p.id AND r.verified == true; r.rating) >= 4.0
    ==> verified_quality(p.name)

# Faits
Product(id="P1", name="Laptop")
Product(id="P2", name="Mouse")

Review(product_id="P1", rating=5.0, verified=true)
Review(product_id="P1", rating=4.5, verified=true)
Review(product_id="P1", rating=4.0, verified=false)

Review(product_id="P2", rating=3.0, verified=true)
Review(product_id="P2", rating=3.5, verified=true)
```

**Résultats :**
- `highly_rated("Laptop")` ✓ (moyenne totale = 4.5)
- `verified_quality("Laptop")` ✓ (moyenne vérifiés = 4.75)

### MIN/MAX

```tsd
type Store : <id: string, name: string>
type StockLevel : <store_id: string, product: string, quantity: number>

# Règle : Magasins avec stock minimum critique (< 5)
{s: Store} /
    MIN(sl: StockLevel / sl.store_id == s.id; sl.quantity) < 5
    ==> critical_stock_level(s.name)

# Règle : Magasins bien approvisionnés (tous produits > 20)
{s: Store} /
    MIN(sl: StockLevel / sl.store_id == s.id; sl.quantity) > 20
    ==> well_stocked(s.name)

# Faits
Store(id="S1", name="Store A")
Store(id="S2", name="Store B")

StockLevel(store_id="S1", product="Product1", quantity=3)
StockLevel(store_id="S1", product="Product2", quantity=50)

StockLevel(store_id="S2", product="Product1", quantity=25)
StockLevel(store_id="S2", product="Product2", quantity=30)
```

**Résultats :**
- `critical_stock_level("Store A")` ✓ (min = 3)
- `well_stocked("Store B")` ✓ (min = 25)

---

## Projet complet

Créons un système de gestion de bibliothèque complet.

```tsd
# ========================================
# TYPES
# ========================================

type Member : <
    id: string,
    name: string,
    email: string,
    premium: bool,
    registration_year: number
>

type Book : <
    id: string,
    title: string,
    author: string,
    year: number,
    available: bool,
    popularity_score: number
>

type Loan : <
    id: string,
    member_id: string,
    book_id: string,
    days: number,
    returned: bool
>

type Review : <
    member_id: string,
    book_id: string,
    rating: number
>

# ========================================
# RÈGLES - Gestion des membres
# ========================================

# R1: Membres actifs (au moins 1 emprunt en cours)
{m: Member} /
    EXISTS (l: Loan / l.member_id == m.id AND l.returned == false)
    ==> active_member(m.name)

# R2: Membres fidèles (plus de 10 emprunts complétés)
{m: Member} /
    COUNT(l: Loan / l.member_id == m.id AND l.returned == true) > 10
    ==> loyal_member(m.name, m.email)

# R3: Promotion premium pour membres anciens et actifs
{m: Member} /
    m.premium == false AND
    m.registration_year <= 2020 AND
    COUNT(l: Loan / l.member_id == m.id) >= 5
    ==> offer_premium(m.id, m.email)

# ========================================
# RÈGLES - Gestion des livres
# ========================================

# R4: Livres populaires (bien notés et souvent empruntés)
{b: Book} /
    AVG(r: Review / r.book_id == b.id; r.rating) >= 4.0 AND
    COUNT(l: Loan / l.book_id == b.id) > 15
    ==> popular_book(b.title, b.author)

# R5: Livres nécessitant plus d'exemplaires
{b: Book} /
    b.available == false AND
    COUNT(l: Loan / l.book_id == b.id AND l.returned == false) > 5
    ==> purchase_more_copies(b.id, b.title)

# R6: Livres à archiver (vieux et peu populaires)
{b: Book} /
    b.year < 1990 AND
    COUNT(l: Loan / l.book_id == b.id) < 3 AND
    b.popularity_score < 2.0
    ==> archive_book(b.id, b.title)

# ========================================
# RÈGLES - Gestion des emprunts
# ========================================

# R7: Emprunts en retard (> 30 jours)
{m: Member, l: Loan} /
    m.id == l.member_id AND
    l.returned == false AND
    l.days > 30
    ==> late_return_notice(m.email, l.id, l.days)

# R8: Prolongation automatique pour membres premium
{m: Member, l: Loan} /
    m.id == l.member_id AND
    m.premium == true AND
    l.returned == false AND
    l.days >= 20 AND
    l.days <= 25
    ==> auto_extend_loan(m.name, l.id)

# R9: Emprunt massif détecté (alerte)
{m: Member} /
    COUNT(l: Loan / l.member_id == m.id AND l.returned == false) > 5 AND
    m.premium == false
    ==> excessive_loans_alert(m.id, m.name)

# ========================================
# RÈGLES - Recommandations
# ========================================

# R10: Recommander livre à membre basé sur historique
{m: Member, b: Book} /
    b.available == true AND
    EXISTS (l: Loan / l.member_id == m.id AND l.book_id != b.id) AND
    AVG(r: Review / r.book_id == b.id; r.rating) >= 4.5
    ==> recommend_book(m.email, b.title, b.author)

# R11: Membres inactifs à relancer (aucun emprunt récent)
{m: Member} /
    NOT (EXISTS (l: Loan / l.member_id == m.id AND l.returned == false)) AND
    COUNT(l: Loan / l.member_id == m.id) > 0
    ==> reengage_member(m.email, m.name)

# ========================================
# FAITS - Membres
# ========================================

Member(id="M001", name="Alice Martin", email="alice@mail.com", premium=false, registration_year=2019)
Member(id="M002", name="Bob Durand", email="bob@mail.com", premium=true, registration_year=2015)
Member(id="M003", name="Claire Dubois", email="claire@mail.com", premium=false, registration_year=2023)

# ========================================
# FAITS - Livres
# ========================================

Book(id="B001", title="1984", author="George Orwell", year=1949, available=true, popularity_score=4.8)
Book(id="B002", title="Clean Code", author="Robert Martin", year=2008, available=false, popularity_score=4.5)
Book(id="B003", title="Old Forgotten Book", author="Unknown", year=1985, available=true, popularity_score=1.2)

# ========================================
# FAITS - Emprunts
# ========================================

Loan(id="L001", member_id="M001", book_id="B001", days=5, returned=false)
Loan(id="L002", member_id="M001", book_id="B002", days=35, returned=false)
Loan(id="L003", member_id="M002", book_id="B001", days=22, returned=false)

# Historique d'emprunts complétés pour M001
Loan(id="L004", member_id="M001", book_id="B001", days=14, returned=true)
Loan(id="L005", member_id="M001", book_id="B002", days=10, returned=true)
Loan(id="L006", member_id="M001", book_id="B003", days=7, returned=true)
Loan(id="L007", member_id="M001", book_id="B001", days=12, returned=true)
Loan(id="L008", member_id="M001", book_id="B002", days=9, returned=true)

# ========================================
# FAITS - Avis
# ========================================

Review(member_id="M001", book_id="B001", rating=5.0)
Review(member_id="M002", book_id="B001", rating=4.5)
Review(member_id="M001", book_id="B002", rating=4.0)
```

### Résultats attendus

Ce programme déclenchera les actions suivantes :

1. **`active_member("Alice Martin")`** - Alice a des emprunts en cours
2. **`active_member("Bob Durand")`** - Bob a un emprunt en cours
3. **`offer_premium("M001", "alice@mail.com")`** - Alice est membre depuis 2019 avec 5+ emprunts
4. **`late_return_notice("alice@mail.com", "L002", 35)`** - Emprunt L002 en retard (35 jours)
5. **`auto_extend_loan("Bob Durand", "L003")`** - Bob est premium et l'emprunt est entre 20-25 jours
6. **`popular_book("1984", "George Orwell")`** - Bien noté (4.75) mais besoin de vérifier le nombre d'emprunts
7. **`archive_book("B003", "Old Forgotten Book")`** - Livre ancien (1985), peu populaire
8. **`reengage_member("claire@mail.com", "Claire Dubois")`** - Claire n'a aucun emprunt

---

## Exercices pratiques

### Exercice 1 : E-commerce

Créez un système de gestion de promotions pour un site e-commerce :

**Types nécessaires :**
- `Customer` : id, name, email, total_spent, registration_date
- `Product` : id, name, price, category, in_stock
- `Order` : id, customer_id, product_id, quantity, status

**Règles à implémenter :**
1. Remise 10% pour clients ayant dépensé > 1000€
2. Alerte stock faible si produit commandé plus de 10 fois et stock < 5
3. Upgrade gratuit vers livraison express pour commandes > 200€

### Exercice 2 : Gestion RH

Créez un système d'évaluation des employés :

**Types nécessaires :**
- `Employee` : id, name, department, salary, years_service
- `Performance` : employee_id, score, year
- `Training` : employee_id, course_name, completed

**Règles à implémenter :**
1. Éligibilité promotion si score moyen > 4.0 et ancienneté > 3 ans
2. Formation obligatoire si score < 3.0 deux années consécutives
3. Bonus annuel si score >= 4.5 et au moins 2 formations complétées

### Exercice 3 : IoT et monitoring

Créez un système d'alertes pour des capteurs :

**Types nécessaires :**
- `Sensor` : id, location, type, status
- `Reading` : sensor_id, value, timestamp
- `Alert` : sensor_id, level, message

**Règles à implémenter :**
1. Alerte critique si valeur > seuil maximum 3 fois consécutives
2. Capteur défaillant si pas de lecture depuis > 1 heure
3. Moyenne anormale si moyenne des 10 dernières lectures > 2x la normale

---

## Prochaines étapes

Maintenant que vous maîtrisez TSD :

1. **Explorez les exemples** : Consultez [`beta_coverage_tests/`](../beta_coverage_tests/) pour plus d'exemples
2. **Lisez la référence** : Consultez [GRAMMAR_GUIDE.md](./GRAMMAR_GUIDE.md) pour tous les détails
3. **Testez vos règles** : Utilisez `tsd` pour valider vos programmes
4. **Optimisez** : Apprenez les bonnes pratiques du moteur RETE dans [`rete/README.md`](../rete/README.md)

Bon codage avec TSD ! 🚀
