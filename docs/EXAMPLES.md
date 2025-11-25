# Exemples TSD - Guide pratique

Cette page contient des exemples complets et commentés pour différents domaines d'application de TSD.

## Table des matières

1. [Exemples de base](#exemples-de-base)
2. [E-commerce](#e-commerce)
3. [IoT et monitoring](#iot-et-monitoring)
4. [Exemples du projet](#exemples-du-projet)

---

## Exemples de base

### Exemple 1 : Gestion d'utilisateurs simple

```tsd
# Types
type User : <id: string, name: string, age: number, active: bool>
type Account : <user_id: string, balance: number, verified: bool>

# Règles
{u: User} / u.age >= 18 AND u.active == true ==> adult_active_user(u.id)

{u: User, a: Account} /
    u.id == a.user_id AND
    a.verified == true AND
    a.balance > 1000
    ==> premium_account(u.name, a.balance)

# Faits
User(id="U1", name="Alice", age=25, active=true)
User(id="U2", name="Bob", age=17, active=true)

Account(user_id="U1", balance=5000, verified=true)
Account(user_id="U2", balance=500, verified=false)
```

**Actions déclenchées :**
- `adult_active_user("U1")` - Alice a 25 ans et est active
- `premium_account("Alice", 5000)` - Compte vérifié avec solde élevé

---

## E-commerce

### Système de fidélité complet

```tsd
# ========================================
# Types
# ========================================

type Customer : <
    id: string,
    email: string,
    loyalty_tier: string,
    join_date: string
>

type Product : <
    id: string,
    name: string,
    category: string,
    price: number,
    stock: number
>

type Order : <
    id: string,
    customer_id: string,
    product_id: string,
    quantity: number,
    total: number,
    status: string
>

# ========================================
# Règles - Fidélité
# ========================================

# Upgrade Gold (1000€+ dépensés)
{c: Customer} /
    c.loyalty_tier != "gold" AND
    SUM(o: Order / o.customer_id == c.id AND o.status == "completed"; o.total) >= 1000
    ==> upgrade_to_gold(c.id, c.email)

# Client VIP avec commandes nombreuses
{c: Customer} /
    COUNT(o: Order / o.customer_id == c.id AND o.status == "completed") > 20
    ==> vip_customer_recognition(c.id, c.email)

# ========================================
# Règles - Gestion des stocks
# ========================================

# Stock critique pour produits populaires
{p: Product} /
    p.stock < 10 AND
    COUNT(o: Order / o.product_id == p.id) > 50
    ==> urgent_restock(p.id, p.name, 100)

# Produit en rupture avec demande active
{p: Product} /
    p.stock == 0 AND
    COUNT(o: Order / o.product_id == p.id AND o.status == "pending") > 5
    ==> out_of_stock_high_demand(p.id, p.name)

# ========================================
# Règles - Promotions
# ========================================

# Livraison gratuite pour gros paniers
{c: Customer, o: Order} /
    c.id == o.customer_id AND
    o.status == "pending" AND
    o.total > 200
    ==> free_shipping(o.id, c.email)

# Réduction fidélité
{c: Customer, o: Order} /
    c.id == o.customer_id AND
    c.loyalty_tier IN ["gold", "platinum"] AND
    o.total > 100
    ==> loyalty_discount(o.id, 0.10)

# ========================================
# Faits
# ========================================

Customer(id="C001", email="alice@email.com", loyalty_tier="silver", join_date="2023-01-15")
Customer(id="C002", email="bob@email.com", loyalty_tier="gold", join_date="2021-06-20")

Product(id="P001", name="Laptop Pro 15", category="electronics", price=1299, stock=5)
Product(id="P002", name="Wireless Mouse", category="accessories", price=29, stock=0)

Order(id="O001", customer_id="C001", product_id="P001", quantity=1, total=1299, status="completed")
Order(id="O002", customer_id="C002", product_id="P002", quantity=2, total=58, status="pending")
```

**Résultats attendus :**
- `upgrade_to_gold("C001", "alice@email.com")` - Alice dépasse 1000€
- `free_shipping("O002", "bob@email.com")` - Panier > 200€ (si applicable)
- `urgent_restock("P001", "Laptop Pro 15", 100)` - Stock faible pour produit populaire
- `out_of_stock_high_demand("P002", "Wireless Mouse")` - Rupture avec demande

---

## IoT et monitoring

### Monitoring de capteurs avec alertes

```tsd
# ========================================
# Types
# ========================================

type Sensor : <
    id: string,
    location: string,
    type: string,
    status: string,
    battery_level: number
>

type Reading : <
    sensor_id: string,
    value: number,
    unit: string,
    timestamp: string,
    quality: string
>

type Threshold : <
    sensor_type: string,
    min_value: number,
    max_value: number,
    critical: bool
>

# ========================================
# Règles - Alertes critiques
# ========================================

# Température critique
{s: Sensor, r: Reading, th: Threshold} /
    s.id == r.sensor_id AND
    s.type == th.sensor_type AND
    s.type == "temperature" AND
    (r.value < th.min_value OR r.value > th.max_value) AND
    th.critical == true
    ==> critical_temperature_alert(s.id, s.location, r.value)

# Humidité anormale
{s: Sensor, r: Reading} /
    s.id == r.sensor_id AND
    s.type == "humidity" AND
    (r.value < 20 OR r.value > 80)
    ==> humidity_warning(s.id, s.location, r.value)

# ========================================
# Règles - Maintenance
# ========================================

# Batterie faible
{s: Sensor} /
    s.battery_level < 20 AND
    s.status == "active"
    ==> low_battery_alert(s.id, s.location, s.battery_level)

# Capteur défaillant (pas de lecture récente)
{s: Sensor} /
    s.status == "active" AND
    NOT (EXISTS (r: Reading / r.sensor_id == s.id AND r.timestamp >= "2024-11-15T00:00:00"))
    ==> sensor_malfunction(s.id, s.location, s.type)

# Qualité dégradée
{s: Sensor} /
    COUNT(r: Reading / r.sensor_id == s.id AND r.quality == "poor") > 5
    ==> sensor_calibration_needed(s.id, s.location)

# ========================================
# Règles - Patterns et anomalies
# ========================================

# Moyenne anormale
{s: Sensor, th: Threshold} /
    s.type == th.sensor_type AND
    AVG(r: Reading / r.sensor_id == s.id; r.value) > th.max_value
    ==> abnormal_average_alert(s.id, s.location, th.max_value)

# Plusieurs capteurs défaillants (même zone)
{s: Sensor} /
    s.status != "active" AND
    COUNT(sen: Sensor / sen.location == s.location AND sen.status != "active") >= 3
    ==> zone_failure_alert(s.location)

# ========================================
# Faits
# ========================================

Sensor(id="S001", location="Server Room A", type="temperature", status="active", battery_level=85)
Sensor(id="S002", location="Server Room A", type="humidity", status="active", battery_level=15)
Sensor(id="S003", location="Warehouse B", type="temperature", status="inactive", battery_level=0)

Reading(sensor_id="S001", value=55, unit="celsius", timestamp="2024-11-15T10:00:00", quality="good")
Reading(sensor_id="S001", value=58, unit="celsius", timestamp="2024-11-15T10:05:00", quality="good")
Reading(sensor_id="S002", value=85, unit="percent", timestamp="2024-11-15T10:00:00", quality="poor")

Threshold(sensor_type="temperature", min_value=-10, max_value=50, critical=true)
Threshold(sensor_type="humidity", min_value=30, max_value=70, critical=false)
```

**Résultats attendus :**
- `critical_temperature_alert("S001", "Server Room A", 55)` - Température > 50°C
- `low_battery_alert("S002", "Server Room A", 15)` - Batterie < 20%
- `humidity_warning("S002", "Server Room A", 85)` - Humidité > 80%
- `sensor_calibration_needed("S002", "Server Room A")` - Qualité dégradée

---

## Exemples du projet

Le projet TSD contient de nombreux exemples dans différents répertoires :

### Tests Beta (Beta Coverage)

Le répertoire [`beta_coverage_tests/`](../beta_coverage_tests/) contient **42 fichiers de tests** couvrant tous les aspects de la grammaire :

**Tests de jointures :**
- `join_simple.constraint` / `.facts` - Jointures basiques
- `join_multi_variable_complex.constraint` / `.facts` - Jointures multiples
- `join_comparison_operators.constraint` / `.facts` - Tous les opérateurs de comparaison
- `join_arithmetic_operators.constraint` / `.facts` - Opérations arithmétiques
- `join_in_contains_operators.constraint` / `.facts` - Opérateurs IN et CONTAINS
- `join_and_operator.constraint` / `.facts` - Conditions avec AND
- `join_or_operator.constraint` / `.facts` - Conditions avec OR

**Tests EXISTS :**
- `exists_simple.constraint` / `.facts` - EXISTS basique
- `exists_complex_operator.constraint` / `.facts` - EXISTS avec conditions complexes
- `beta_exists_complex.constraint` / `.facts` - EXISTS avancé dans réseau Beta

**Tests NOT :**
- `not_simple.constraint` / `.facts` - NOT basique
- `not_complex_operator.constraint` / `.facts` - NOT avec conditions complexes
- `beta_not_complex.constraint` / `.facts` - NOT avancé dans réseau Beta

**Tests d'agrégation :**
- `beta_accumulate_count.constraint` / `.facts` - COUNT
- `beta_accumulate_sum.constraint` / `.facts` - SUM
- `beta_accumulate_avg.constraint` / `.facts` - AVG
- `beta_accumulate_minmax.constraint` / `.facts` - MIN/MAX

**Tests combinés :**
- `complex_not_exists_combination.constraint` / `.facts` - NOT + EXISTS
- `beta_pattern_complex.constraint` / `.facts` - Patterns complexes
- `beta_join_complex.constraint` / `.facts` - Jointures complexes multi-variables

### Comment utiliser les exemples

```bash
# Lister tous les exemples
ls -1 beta_coverage_tests/*.constraint

# Voir un exemple spécifique
cat beta_coverage_tests/join_simple.constraint
cat beta_coverage_tests/join_simple.facts

# Exécuter un exemple (si le runner est disponible)
./bin/universal-rete-runner beta_coverage_tests/join_simple.constraint beta_coverage_tests/join_simple.facts
```

### Structure des fichiers

Chaque exemple se compose de **deux fichiers** :

1. **`.constraint`** - Définitions de types et règles
2. **`.facts`** - Instances de données (faits)

**Exemple de structure :**

```tsd
# Fichier: example.constraint
type Person : <name: string, age: number>
{p: Person} / p.age >= 18 ==> adult(p.name)
```

```tsd
# Fichier: example.facts
Person(name="Alice", age=25)
Person(name="Bob", age=17)
```

---

## Exercices pratiques

Pour approfondir vos connaissances, essayez de créer vos propres exemples :

### Exercice 1 : Finance
Créez un système de détection de fraude bancaire avec :
- Types : Account, Transaction, Alert
- Règles : montants inhabituels, transactions multiples rapides, comptes suspendus
- Agrégations : somme quotidienne, moyenne des transactions

### Exercice 2 : Santé
Créez un système d'alertes médicales avec :
- Types : Patient, Vital, Medication
- Règles : valeurs critiques (tension, glucose, fréquence cardiaque)
- Quantification : interactions médicamenteuses, allergies

### Exercice 3 : Logistique
Créez un système de gestion d'entrepôt avec :
- Types : Warehouse, Product, Shipment
- Règles : capacité, répartition des stocks, expéditions en retard
- Agrégations : stock total par entrepôt, produits critiques

---

## Ressources complémentaires

- **Guide de référence** : [GRAMMAR_GUIDE.md](./GRAMMAR_GUIDE.md) - Syntaxe complète
- **Tutoriel** : [TUTORIAL.md](./TUTORIAL.md) - Apprentissage progressif
- **Tests d'intégration** : `beta_coverage_tests/` - 42 exemples complets
- **Documentation RETE** : [`rete/README.md`](../rete/README.md) - Moteur d'inférence
- **Grammaire PEG** : [`constraint/grammar/constraint.peg`](../constraint/grammar/constraint.peg) - Spécification formelle

---

## Contribution

Pour ajouter vos propres exemples :

1. Créez deux fichiers : `your_example.constraint` et `your_example.facts`
2. Placez-les dans `beta_coverage_tests/` ou créez votre propre répertoire
3. Testez-les avec le runner TSD
4. Soumettez une pull request avec des commentaires explicatifs

Bon codage avec TSD ! 🚀
