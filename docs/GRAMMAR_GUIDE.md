# Guide Complet de la Grammaire TSD

**TSD (Type System Declarations)** est un langage déclaratif pour définir des règles métier avec un système de types fort et une intégration native avec un moteur RETE.

## Table des matières

1. [Vue d'ensemble](#vue-densemble)
2. [Types de base](#types-de-base)
3. [Définition de types](#définition-de-types)
4. [Variables typées](#variables-typées)
5. [Expressions et règles](#expressions-et-règles)
6. [Conditions et contraintes](#conditions-et-contraintes)
7. [Opérateurs](#opérateurs)
8. [Fonctions intégrées](#fonctions-intégrées)
9. [Actions](#actions)
10. [Faits](#faits)
11. [Exemples avancés](#exemples-avancés)

---

## Vue d'ensemble

Un programme TSD se compose de trois éléments principaux :

1. **Définitions de types** - Structure des données
2. **Règles (Expressions)** - Logique métier
3. **Faits** - Instances de données

### Structure minimale

\`\`\`tsd
# 1. Définir les types
type Person : <name: string, age: number>

# 2. Définir une règle
{p: Person} / p.age >= 18 ==> adult(p.name)

# 3. Créer des faits
Person(name="Alice", age=25)
\`\`\`

---

## Types de base

TSD supporte trois types atomiques :

| Type | Description | Exemples |
|------|-------------|----------|
| \`string\` | Chaîne de caractères | \`"Hello"\`, \`'World'\` |
| \`number\` | Nombre (entier ou décimal) | \`42\`, \`3.14\`, \`-10\` |
| \`bool\` | Booléen | \`true\`, \`false\` |

---

## Définition de types

### Syntaxe

\`\`\`tsd
type TypeName : <field1: type1, field2: type2, ...>
\`\`\`

### Exemples

\`\`\`tsd
# Type simple
type User : <id: string, name: string, age: number>

# Type pour e-commerce
type Product : <
    id: string,
    name: string,
    price: number,
    available: bool
>

# Type pour commande
type Order : <
    id: string,
    user_id: string,
    product_id: string,
    quantity: number,
    total: number
>
\`\`\`

### Bonnes pratiques

- Utilisez des noms descriptifs en **PascalCase** (ex: \`UserAccount\`, \`OrderItem\`)
- Nommez les champs en **snake_case** ou **camelCase** (ex: \`user_id\` ou \`userId\`)
- Gardez les types simples et focalisés sur une seule responsabilité

---

## Variables typées

Dans les règles, chaque variable doit être typée explicitement :

\`\`\`tsd
{variableName: TypeName}
\`\`\`

### Exemples

\`\`\`tsd
# Variable unique
{u: User} / u.age > 18 ==> process_adult(u.name)

# Plusieurs variables (jointure)
{u: User, o: Order} / u.id == o.user_id ==> user_has_order(u.name, o.id)

# Trois variables ou plus
{u: User, o: Order, p: Product} / 
    u.id == o.user_id AND o.product_id == p.id 
    ==> complete_order(u.name, p.name, o.quantity)
\`\`\`

---

## Expressions et règles

Une règle TSD a la structure suivante :

\`\`\`tsd
{variables} / conditions ==> action
\`\`\`

### Composants

1. **Set de variables** : \`{var1: Type1, var2: Type2}\`
2. **Séparateur** : \`/\`
3. **Conditions** : Expression booléenne
4. **Flèche d'action** : \`==>\`
5. **Action** : Appel de fonction

### Exemples

\`\`\`tsd
# Règle simple - Une seule variable
{p: Person} / p.age >= 18 ==> mark_as_adult(p.id)

# Règle avec jointure - Deux variables
{u: User, a: Account} / 
    u.id == a.user_id AND a.balance > 1000 
    ==> high_balance_alert(u.email, a.balance)

# Règle complexe - Trois variables
{c: Customer, o: Order, p: Product} /
    c.id == o.customer_id AND
    o.product_id == p.id AND
    o.quantity > 10 AND
    p.stock >= o.quantity
    ==> bulk_order_discount(c.id, o.id, 0.15)
\`\`\`

---

## Conditions et contraintes

### Conditions simples

\`\`\`tsd
{u: User} / u.age > 18 ==> action()
{p: Product} / p.available == true ==> action()
{o: Order} / o.total >= 100 ==> action()
\`\`\`

### Conditions logiques (AND/OR)

⚠️ **Priorité des opérateurs** : Les opérateurs logiques suivent la priorité standard :
- **NOT** (priorité maximale) - évalué en premier
- **AND** (priorité moyenne) - évalué après NOT
- **OR** (priorité minimale) - évalué en dernier

Sans parenthèses, `a OR b AND c` est évalué comme `a OR (b AND c)`, pas comme `(a OR b) AND c`.

⚠️ **Priorité des opérateurs** : Les opérateurs logiques suivent la priorité standard :
- **NOT** (priorité maximale) - évalué en premier
- **AND** (priorité moyenne) - évalué après NOT
- **OR** (priorité minimale) - évalué en dernier

Sans parenthèses, `a OR b AND c` est évalué comme `a OR (b AND c)`, pas comme `(a OR b) AND c`.

⚠️ **Priorité des opérateurs** : Les opérateurs logiques suivent la priorité standard :
- **NOT** (priorité maximale) - évalué en premier
- **AND** (priorité moyenne) - évalué après NOT
- **OR** (priorité minimale) - évalué en dernier

Sans parenthèses, `a OR b AND c` est évalué comme `a OR (b AND c)`, pas comme `(a OR b) AND c`.

\`\`\`tsd
# AND - Toutes les conditions doivent être vraies
{p: Person} / p.age >= 18 AND p.country == "FR" ==> french_adult(p.name)

# OR - Au moins une condition doit être vraie
{u: User} / u.vip == true OR u.orders > 10 ==> premium_customer(u.id)

# Priorité AND > OR : cette expression signifie "age < 18 OU (status == 'student' ET discount == true)"
{p: Person} / p.age < 18 OR p.status == "student" AND p.discount == true ==> apply_reduction(p.id)

# Utiliser des parenthèses pour clarifier ou changer la priorité
{p: Person} / 
    (p.age >= 18 AND p.age <= 65) OR p.retired == true 
    ==> eligible_for_program(p.id)
\`\`\`

### NOT - Négation

\`\`\`tsd
# Simple négation
{u: User} / NOT (u.banned == true) ==> allow_access(u.id)

# Négation complexe
{p: Product} / 
    NOT (p.price > 1000 OR p.restricted == true) 
    ==> can_purchase_without_approval(p.id)
\`\`\`

### EXISTS - Quantification existentielle

Vérifie qu'au moins un élément satisfait une condition :

\`\`\`tsd
# Un utilisateur qui a au moins une commande confirmée
{u: User} /
    EXISTS (o: Order / o.user_id == u.id AND o.status == "confirmed")
    ==> active_customer(u.id)

# Un produit vendu au moins une fois
{p: Product} /
    EXISTS (o: Order / o.product_id == p.id AND o.quantity > 0)
    ==> popular_product(p.id)
\`\`\`

### Fonctions d'agrégation

TSD supporte 5 fonctions d'agrégation : \`SUM\`, \`COUNT\`, \`AVG\`, \`MIN\`, \`MAX\`

#### COUNT - Compter les éléments

\`\`\`tsd
# Utilisateur avec plus de 5 commandes
{u: User} /
    COUNT(o: Order / o.user_id == u.id) > 5
    ==> frequent_buyer(u.id)
\`\`\`

#### SUM - Somme des valeurs

\`\`\`tsd
# Client avec un total d'achats > 1000
{c: Customer} /
    SUM(o: Order / o.customer_id == c.id; o.total) > 1000
    ==> vip_customer(c.id)
\`\`\`

#### AVG - Moyenne des valeurs

\`\`\`tsd
# Produit avec note moyenne > 4.5
{p: Product} /
    AVG(r: Review / r.product_id == p.id; r.rating) >= 4.5
    ==> highly_rated(p.id)
\`\`\`

#### MIN/MAX - Minimum/Maximum

\`\`\`tsd
# Utilisateur dont la commande minimale > 50
{u: User} /
    MIN(o: Order / o.user_id == u.id; o.total) > 50
    ==> high_value_customer(u.id)

# Produit dont le stock maximum < 10
{p: Product} /
    MAX(s: Stock / s.product_id == p.id; s.quantity) < 10
    ==> low_stock_alert(p.id)
\`\`\`

---

## Opérateurs

### Priorité des opérateurs

TSD respecte la hiérarchie standard des opérateurs (du plus prioritaire au moins prioritaire) :

1. **Parenthèses** : `( )`
2. **Accès aux champs** : `object.field`
3. **Fonctions** : `func(args)`
4. **Multiplication/Division** : `*`, `/`
5. **Addition/Soustraction** : `+`, `-`
6. **Comparaisons** : `==`, `!=`, `<`, `>`, `<=`, `>=`, `IN`, `LIKE`, `CONTAINS`, `MATCHES`
7. **NOT** : Négation logique
8. **AND** : Conjonction logique
9. **OR** : Disjonction logique

**Exemple** : `a == 1 OR b == 2 AND c == 3` est évalué comme `a == 1 OR (b == 2 AND c == 3)`, car AND a une priorité supérieure à OR.


| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| \`==\` | Égal | \`x == 10\` |
| \`!=\` | Différent | \`x != 0\` |
| \`<\` | Inférieur | \`x < 100\` |
| \`<=\` | Inférieur ou égal | \`x <= 100\` |
| \`>\` | Supérieur | \`x > 0\` |
| \`>=\` | Supérieur ou égal | \`x >= 18\` |

### Opérateurs spéciaux

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| \`IN\` | Appartenance à une liste | \`status IN ["active", "pending"]\` |
| \`CONTAINS\` | Contient une sous-chaîne | \`text CONTAINS "error"\` |
| \`LIKE\` | Correspondance de motif | \`email LIKE "%@gmail.com"\` |
| \`MATCHES\` | Expression régulière | \`code MATCHES "^[A-Z]{3}[0-9]{3}$"\` |

### Opérateurs arithmétiques

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| \`+\` | Addition | \`price + tax\` |
| \`-\` | Soustraction | \`total - discount\` |
| \`*\` | Multiplication | \`quantity * price\` |
| \`/\` | Division | \`total / quantity\` |

### Exemples d'utilisation

\`\`\`tsd
# IN - Liste de valeurs
{o: Order} / 
    o.status IN ["pending", "processing", "shipped"]
    ==> track_order(o.id)

# CONTAINS - Recherche dans chaîne
{p: Product} /
    p.description CONTAINS "organic"
    ==> organic_product(p.id)

# Arithmétique
{o: Order} /
    (o.subtotal + o.tax - o.discount) > 100
    ==> free_shipping(o.id)
\`\`\`

---

## Fonctions intégrées

### Fonctions de chaînes

| Fonction | Description | Exemple |
|----------|-------------|---------|
| \`LENGTH(str)\` | Longueur d'une chaîne | \`LENGTH(p.name) > 10\` |
| \`UPPER(str)\` | Convertir en majuscules | \`UPPER(u.country) == "FRANCE"\` |
| \`LOWER(str)\` | Convertir en minuscules | \`LOWER(p.category) == "electronics"\` |
| \`TRIM(str)\` | Supprimer espaces | \`TRIM(u.name) != ""\` |
| \`SUBSTRING(str, start, len)\` | Extraire sous-chaîne | \`SUBSTRING(code, 0, 3) == "ABC"\` |

### Fonctions mathématiques

| Fonction | Description | Exemple |
|----------|-------------|---------|
| \`ABS(num)\` | Valeur absolue | \`ABS(account.balance) > 100\` |
| \`ROUND(num)\` | Arrondir | \`ROUND(price) == 10\` |
| \`FLOOR(num)\` | Arrondir vers le bas | \`FLOOR(3.7) == 3\` |
| \`CEIL(num)\` | Arrondir vers le haut | \`CEIL(3.2) == 4\` |

### Exemples

\`\`\`tsd
# Validation de longueur
{u: User} /
    LENGTH(u.password) >= 8
    ==> valid_password(u.id)

# Conversion et comparaison
{p: Product} /
    UPPER(p.category) == "ELECTRONICS"
    ==> electronics_discount(p.id)

# Calculs mathématiques
{o: Order} /
    ROUND(o.total * 1.2) <= u.budget
    ==> affordable_order(o.id)
\`\`\`

---

## Actions

Les actions sont des appels de fonction déclenchés lorsque toutes les conditions sont satisfaites.

### Syntaxe

\`\`\`tsd
==> function_name(arg1, arg2, ...)
\`\`\`

### Types d'arguments

\`\`\`tsd
# Arguments simples - Accès aux champs
{u: User} / u.age >= 18 ==> notify_adult(u.name, u.email)

# Expressions arithmétiques
{o: Order} / o.quantity > 10 ==> bulk_discount(o.id, o.quantity * 0.9)

# Constantes littérales
{p: Product} / p.stock < 5 ==> restock_alert(p.id, "LOW_STOCK", 100)

# Combinaisons
{u: User, o: Order} /
    u.id == o.user_id
    ==> process_order(
        u.email,
        o.id,
        o.total * 1.2,
        "PRIORITY"
    )
\`\`\`

---

## Faits

Les faits sont des instances concrètes des types définis.

### Syntaxe

\`\`\`tsd
TypeName(field1=value1, field2=value2, ...)
\`\`\`

### Exemples

\`\`\`tsd
# Fait simple
User(id="U001", name="Alice", age=30)

# Fait avec tous les types
Product(
    id="P001",
    name="Laptop",
    price=999.99,
    available=true
)

# Plusieurs faits
Order(id="O001", user_id="U001", product_id="P001", quantity=2, total=1999.98)
Order(id="O002", user_id="U002", product_id="P002", quantity=1, total=49.99)

# Fait avec chaînes entre guillemets simples ou doubles
Person(name='Bob', email="bob@example.com", country="FR")
\`\`\`

### Types de valeurs

\`\`\`tsd
# String - guillemets simples ou doubles
User(name="Alice", nickname='Ali')

# Number - entiers ou décimaux
Product(price=99.99, stock=42, discount=-10)

# Boolean - true ou false
Account(active=true, verified=false)
\`\`\`

### Rétractation de faits

TSD supporte la rétractation (suppression) de faits avec la commande `remove` :

\`\`\`tsd
# Syntaxe : remove TypeName FactID
remove Person P1
remove Order O123
\`\`\`

**Exemple complet avec ajout et rétractation :**

\`\`\`tsd
type Person : <id:string, name:string, age:number>
type Order : <id:string, customer_id:string, total:number>

# Ajouter des faits
Person(id:P1, name:Alice, age:25)
Person(id:P2, name:Bob, age:35)
Order(id:O1, customer_id:P1, total:100)
Order(id:O2, customer_id:P2, total:200)

# Rétracter des faits (suppression du réseau RETE)
remove Person P1
remove Order O2

# Ajouter d'autres faits après les rétractions
Person(id:P3, name:Charlie, age:45)
\`\`\`

**Notes importantes :**
- L'ID du fait doit correspondre exactement à l'ID utilisé lors de la création
- La rétractation propage la suppression à travers tout le réseau RETE
- Les tokens contenant le fait rétracté sont automatiquement retirés
- La rétractation utilise les IDs internes (format `Type_ID`)

**Cas d'usage :**
- Mise à jour de données (rétracter puis ajouter avec nouvelles valeurs)
- Simulation de scénarios dynamiques
- Tests de comportement du réseau RETE
- Gestion d'événements temporaires

---

## Exemples avancés

### Exemple 1 : Système de fidélité e-commerce

\`\`\`tsd
# Définitions de types
type Customer : <id: string, name: string, email: string, vip: bool>
type Order : <id: string, customer_id: string, total: number, status: string>
type Product : <id: string, name: string, price: number, category: string>

# Règle 1: Client VIP automatique après 1000€ d'achats
{c: Customer} /
    SUM(o: Order / o.customer_id == c.id AND o.status == "completed"; o.total) > 1000
    ==> promote_to_vip(c.id, c.email)

# Règle 2: Réduction sur les commandes importantes
{c: Customer, o: Order} /
    c.id == o.customer_id AND
    o.total > 500 AND
    c.vip == true
    ==> apply_vip_discount(o.id, 0.15)

# Règle 3: Alerte stock faible pour produits populaires
{p: Product} /
    COUNT(o: Order / o.product_id == p.id) > 10 AND
    p.stock < 5
    ==> urgent_restock(p.id, p.name, 50)

# Faits
Customer(id="C001", name="Alice Martin", email="alice@email.com", vip=false)
Customer(id="C002", name="Bob Durand", email="bob@email.com", vip=true)

Order(id="O001", customer_id="C001", total=250.00, status="completed")
Order(id="O002", customer_id="C001", total=800.00, status="completed")
Order(id="O003", customer_id="C002", total=600.00, status="pending")

Product(id="P001", name="Laptop", price=999.99, stock=3, category="electronics")
\`\`\`

### Exemple 2 : Gestion des ressources humaines

\`\`\`tsd
# Types
type Employee : <
    id: string,
    name: string,
    age: number,
    department: string,
    salary: number,
    years_service: number
>

type Leave : <
    id: string,
    employee_id: string,
    days: number,
    type: string,
    status: string
>

type Performance : <
    employee_id: string,
    score: number,
    year: number
>

# Règle 1: Éligibilité à la promotion (seniors performants)
{e: Employee} /
    e.years_service >= 5 AND
    AVG(p: Performance / p.employee_id == e.id; p.score) >= 4.5
    ==> eligible_for_promotion(e.id, e.name)

# Règle 2: Validation automatique des congés (moins de 5 jours)
{e: Employee, l: Leave} /
    e.id == l.employee_id AND
    l.days <= 5 AND
    l.status == "pending" AND
    COUNT(lv: Leave / lv.employee_id == e.id AND lv.status == "approved") < 3
    ==> auto_approve_leave(l.id, e.name)

# Règle 3: Alerte département sureffectif
{e: Employee} /
    COUNT(emp: Employee / emp.department == e.department) > 20 AND
    AVG(emp: Employee / emp.department == e.department; emp.salary) > 60000
    ==> department_cost_alert(e.department)

# Règle 4: Bonus annuel basé sur performance
{e: Employee} /
    e.years_service >= 2 AND
    AVG(p: Performance / p.employee_id == e.id AND p.year >= 2023; p.score) >= 4.0
    ==> calculate_bonus(e.id, e.salary * 0.10)

# Faits
Employee(
    id="E001",
    name="Sophie Laurent",
    age=35,
    department="Engineering",
    salary=75000,
    years_service=7
)

Performance(employee_id="E001", score=4.8, year=2024)
Performance(employee_id="E001", score=4.6, year=2023)

Leave(id="L001", employee_id="E001", days=3, type="vacation", status="pending")
\`\`\`

### Exemple 3 : Système d'alertes IoT

\`\`\`tsd
# Types
type Sensor : <id: string, location: string, type: string, active: bool>
type Reading : <sensor_id: string, value: number, timestamp: number, unit: string>
type Alert : <id: string, sensor_id: string, severity: string, message: string>

# Règle 1: Température critique
{s: Sensor, r: Reading} /
    s.id == r.sensor_id AND
    s.type == "temperature" AND
    (r.value > 50 OR r.value < -10)
    ==> critical_temperature_alert(s.id, s.location, r.value)

# Règle 2: Capteur défaillant (pas de lecture récente)
{s: Sensor} /
    s.active == true AND
    NOT (EXISTS (r: Reading / r.sensor_id == s.id AND r.timestamp > 1700000000))
    ==> sensor_malfunction(s.id, s.location)

# Règle 3: Moyenne de lectures anormale
{s: Sensor} /
    AVG(r: Reading / r.sensor_id == s.id; r.value) > 100 OR
    AVG(r: Reading / r.sensor_id == s.id; r.value) < 0
    ==> abnormal_readings_pattern(s.id, s.location)

# Faits
Sensor(id="S001", location="Room A", type="temperature", active=true)
Sensor(id="S002", location="Room B", type="humidity", active=true)

Reading(sensor_id="S001", value=55.0, timestamp=1700050000, unit="celsius")
Reading(sensor_id="S001", value=52.0, timestamp=1700050060, unit="celsius")
Reading(sensor_id="S002", value=85.0, timestamp=1700050000, unit="percent")
\`\`\`

---

## Bonnes pratiques

### 1. Nommage cohérent

\`\`\`tsd
# ✅ BON - Noms descriptifs et cohérents
type CustomerAccount : <account_id: string, balance: number>
{ca: CustomerAccount} / ca.balance > 1000 ==> high_value_account(ca.account_id)

# ❌ MAUVAIS - Noms incohérents
type CA : <aid: string, b: number>
{x: CA} / x.b > 1000 ==> hva(x.aid)
\`\`\`

### 2. Décomposition des règles complexes

\`\`\`tsd
# ✅ BON - Règles simples et focalisées
{u: User} / u.age >= 18 ==> adult_user(u.id)
{u: User} / u.orders > 10 ==> frequent_buyer(u.id)

# ❌ MAUVAIS - Règle trop complexe
{u: User, o: Order, p: Product, a: Account} /
    u.id == o.user_id AND o.product_id == p.id AND
    u.account_id == a.id AND a.balance > 1000 AND
    (p.price > 500 OR p.category == "premium") AND
    COUNT(ord: Order / ord.user_id == u.id) > 10
    ==> complex_action(u.id, o.id, p.id, a.id)
\`\`\`

### 3. Utilisation appropriée des agrégations

\`\`\`tsd
# ✅ BON - Agrégation ciblée
{u: User} /
    SUM(o: Order / o.user_id == u.id AND o.status == "completed"; o.total) > 1000
    ==> vip_customer(u.id)

# ❌ MAUVAIS - Agrégation trop large
{u: User} /
    SUM(o: Order / true; o.total) > 1000  # Somme TOUTES les commandes !
    ==> wrong_vip(u.id)
\`\`\`

### 4. Parenthèses pour la clarté

TSD suit la priorité standard des opérateurs (NOT > AND > OR), mais les parenthèses améliorent la lisibilité :

```tsd
# ✅ BON - Clarté avec parenthèses explicites
{p: Person} /
    (p.age >= 18 AND p.age <= 65) OR p.retired == true
    ==> eligible(p.id)

# ✅ Acceptable - Priorité explicite respectée (AND > OR)
# Cette expression signifie : age >= 18 AND (age <= 65 OR retired == true)
{p: Person} /
    p.age >= 18 AND p.age <= 65 OR p.retired == true
    ==> eligible(p.id)

# ⚠️ Recommandé - Ajoutez des parenthèses pour améliorer la lisibilité
{p: Person} /
    p.age >= 18 AND (p.age <= 65 OR p.retired == true)
    ==> eligible(p.id)
```

---

## Génération du parser

Le parser TSD est généré à partir de la grammaire PEG :

\`\`\`bash
# Installer pigeon (générateur PEG pour Go)
go install github.com/mna/pigeon@latest

# Générer le parser
cd constraint/grammar
pigeon -o parser.go constraint.peg
\`\`\`

---

## Ressources

- **Grammaire PEG** : [\`constraint/grammar/constraint.peg\`](../constraint/grammar/constraint.peg)
- **Tests d'intégration** : [\`test/integration/\`](../test/integration/)
- **Exemples Beta** : [\`beta_coverage_tests/\`](../beta_coverage_tests/)

---

## Support et contributions

Pour toute question ou contribution :
- Ouvrir une issue sur GitHub
- Consulter les tests existants pour des exemples
- Lire la documentation du moteur RETE : [\`rete/README.md\`](../rete/README.md)
