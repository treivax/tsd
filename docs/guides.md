# Guides Utilisateur TSD

**Documentation complète** - Du débutant à l'expert

---

## Table des Matières

1. [Guide Débutant](#guide-débutant)
2. [Guide Développeur](#guide-développeur)
3. [Guide Avancé](#guide-avancé)
4. [Cas d'Usage Pratiques](#cas-dusage-pratiques)
5. [Bonnes Pratiques](#bonnes-pratiques)

---

## Guide Débutant

### Introduction

TSD (Type System Development) est un moteur de règles basé sur l'algorithme RETE qui permet de définir des règles métier avec une syntaxe déclarative.

**Caractéristiques principales :**
- Pattern matching sur des combinaisons de faits
- Typage fort avec validation
- Algorithme RETE pour évaluation efficace
- Optimisation réseau (partage de conditions)
- Opérations avancées (casting, regex, arithmétique)

### Structure d'un Programme TSD

Un programme TSD contient quatre éléments principaux :

```tsd
// 1. Définitions de types
type Person(name: string, age: number)

// 2. Déclarations d'actions
action greet(name: string)

// 3. Définitions de règles
rule adult : {p: Person} / p.age >= 18 ==> greet(p.name)

// 4. Assertions de faits
Person(name: "Alice", age: 25)
```

### Étape 1 : Votre Premier Type

Créons un type simple représentant une personne :

```tsd
type Person(name: string, age: number)
```

**Explication :**
- `type Person` : Déclare un nouveau type nommé "Person"
- `name: string` : Champ "name" de type chaîne
- `age: number` : Champ "age" de type nombre

**Types primitifs disponibles :**
- `string` : Chaîne de caractères
- `number` : Nombre (entier ou décimal)
- `bool` : Booléen (true/false)

### Étape 2 : Créer des Faits

Créons quelques instances (faits) de ce type :

```tsd
Person(name: "Alice", age: 25)
Person(name: "Bob", age: 17)
Person(name: "Charlie", age: 30)
```

**Important :** Les valeurs doivent correspondre aux types définis.

### Étape 3 : Votre Première Règle

Identifions les adultes (18 ans ou plus) :

```tsd
action adult(name: string)

rule welcome : {p: Person} / p.age >= 18 ==> adult(p.name)
```

**Structure d'une règle :**
- `{p: Person}` : Pattern - lie les faits Person à la variable `p`
- `/` : Séparateur entre pattern et condition
- `p.age >= 18` : Condition à vérifier
- `==>` : Si vrai, exécuter l'action
- `adult(p.name)` : Appelle l'action avec le nom

**Résultat :** `adult("Alice")` et `adult("Charlie")` seront exécutés, mais pas pour Bob.

### Étape 4 : Conditions Multiples

Combinez des conditions avec `AND` et `OR` :

```tsd
type Product(name: string, price: number, inStock: bool)
action recommend(name: string)

rule affordableAndAvailable : {p: Product} / 
    p.price <= 50 AND p.inStock == true 
    ==> recommend(p.name)

Product(name: "Mouse", price: 25, inStock: true)
Product(name: "Keyboard", price: 45, inStock: false)
```

**Opérateurs de comparaison :**
- `==` : Égal
- `!=` : Différent
- `<` : Inférieur
- `>` : Supérieur
- `<=` : Inférieur ou égal
- `>=` : Supérieur ou égal

**Opérateurs logiques :**
- `AND` : ET logique
- `OR` : OU logique
- `NOT(...)` : NON logique

### Étape 5 : Jointures (Faits Multiples)

Associez plusieurs faits ensemble :

```tsd
type Customer(id: string, name: string, vip: bool)
type Order(customerId: string, total: number)
action applyDiscount(customerName: string, amount: number)

rule vipDiscount : {c: Customer, o: Order} / 
    c.id == o.customerId AND c.vip == true AND o.total > 100 
    ==> applyDiscount(c.name, o.total * 0.1)

Customer(id: "C001", name: "Alice", vip: true)
Order(customerId: "C001", total: 250.00)
```

**Explication :**
- `{c: Customer, o: Order}` : Match deux faits simultanément
- `c.id == o.customerId` : Condition de jointure (lier les faits)
- `o.total * 0.1` : Calcul arithmétique dans l'action

---

## Guide Développeur

### Syntaxe du Langage

#### Commentaires

```tsd
// Commentaire sur une ligne

/* Commentaire
   sur plusieurs lignes */
```

#### Sensibilité à la Casse

TSD est **insensible à la casse** pour les mots-clés mais **sensible** pour les identifiants :

```tsd
// Équivalent (mots-clés)
TYPE Person(name: string)
type Person(name: string)
Type Person(name: string)

// DIFFÉRENT (identifiants)
Person(name: "Alice")    // Type "Person"
person(name: "Bob")      // Type "person" (différent!)
```

#### Identifiants Valides

```tsd
// Valides
myVariable
_underscore
camelCase
PascalCase
snake_case
with123Numbers
αβγ              // Support Unicode
名前              // Support UTF-8

// Invalides
123start         // Ne peut pas commencer par un chiffre
my-var           // Pas de tirets
my.var           // Pas de points (sauf accès champ)
```

### Système de Types

#### Définition de Types

```tsd
// Type simple
type Person(name: string, age: number)

// Type avec plusieurs champs
type Order(
    id: string,
    customerId: string,
    total: number,
    paid: bool,
    createdAt: string
)

// Type avec booléens
type Settings(enabled: bool, verbose: bool)
```

#### Types Primitifs

| Type | Description | Exemples |
|------|-------------|----------|
| `string` | Chaîne de caractères | `"hello"`, `"123"`, `""` |
| `number` | Nombre (int/float) | `42`, `3.14`, `-10` |
| `bool` | Booléen | `true`, `false` |

### Opérations sur Chaînes

#### CONTAINS : Contient

Vérifie si une chaîne contient une sous-chaîne :

```tsd
type Email(address: string, subject: string)
action flagSpam(address: string)

rule spamFilter : {e: Email} / 
    e.subject CONTAINS "URGENT" OR e.subject CONTAINS "Click here"
    ==> flagSpam(e.address)

Email(address: "spam@test.com", subject: "URGENT: Act now!")
```

#### LIKE : Pattern SQL

Pattern matching style SQL (`%` = n'importe quels caractères, `_` = caractère unique) :

```tsd
type File(name: string, path: string)
action processImage(name: string)

rule imageFiles : {f: File} / 
    f.name LIKE "%.png" OR f.name LIKE "%.jpg"
    ==> processImage(f.name)

File(name: "photo.png", path: "/images/photo.png")
```

#### MATCHES : Regex

Pattern matching avec expressions régulières complètes :

```tsd
type Log(message: string, level: string)
action alert(message: string)

rule errorPattern : {l: Log} / 
    l.message MATCHES "^ERROR.*database.*$"
    ==> alert(l.message)

Log(message: "ERROR: database connection failed", level: "error")
```

#### IN : Appartenance

Vérifie l'appartenance à une collection :

```tsd
type User(name: string, role: string)
action grantAccess(name: string)

rule adminAccess : {u: User} / 
    u.role IN ["admin", "superuser", "root"]
    ==> grantAccess(u.name)

User(name: "Alice", role: "admin")
```

### Conversion de Types (Type Casting)

Conversions explicites entre types :

```tsd
type Product(name: string, price: number, quantity: number)
action notify(message: string)

// Convertir nombre en chaîne pour concaténation
rule priceAlert : {p: Product} / p.price > 100 ==> 
    notify("Produit cher: " + p.name + " coûte " + (string)p.price + "€")

Product(name: "Laptop", price: 999.99, quantity: 5)
```

**Conversions disponibles :**

| Cast | Description | Exemple |
|------|-------------|---------|
| `(string)value` | Vers chaîne | `(string)42` → `"42"` |
| `(number)value` | Vers nombre | `(number)"123"` → `123` |
| `(bool)value` | Vers booléen | `(bool)1` → `true` |

### Opérations Arithmétiques

Effectuez des calculs dans les conditions et actions :

```tsd
type Order(id: string, price: number, quantity: number, tax: number)
action createInvoice(orderId: string, total: number)

rule calculateTotal : {o: Order} / o.quantity > 0 ==> 
    createInvoice(o.id, (o.price * o.quantity) + o.tax)

Order(id: "ORD001", price: 50.00, quantity: 3, tax: 15.00)
```

**Opérateurs supportés :**
- `+` : Addition (ou concaténation de chaînes)
- `-` : Soustraction
- `*` : Multiplication
- `/` : Division
- `%` : Modulo

**Priorité des opérateurs :** `*`, `/`, `%` > `+`, `-`

**Utiliser des parenthèses** pour contrôler l'ordre d'évaluation :

```tsd
price * (quantity + bonus)  // Correct
price * quantity + bonus    // Différent!
```

### Intégration Go

#### Utilisation de Base

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/constraint"
)

func main() {
    // Parser le programme TSD
    program, err := constraint.ParseInput("program.tsd")
    if err != nil {
        log.Fatal(err)
    }
    
    // Créer le réseau RETE
    network := rete.NewReteNetwork()
    
    // Compiler le programme
    if err := network.Compile(program); err != nil {
        log.Fatal(err)
    }
    
    // Exécuter
    results := network.Run()
    
    fmt.Printf("Actions exécutées: %d\n", len(results))
}
```

#### Avec Configuration Personnalisée

```go
import "github.com/treivax/tsd/rete"

// Configuration personnalisée
config := &rete.ChainPerformanceConfig{
    MaxChainLength:        100,
    MaxRecursionDepth:     50,
    EnableBetaSharing:     true,
    EnableAlphaSharing:    true,
}

network := rete.NewReteNetworkWithConfig(config)
```

#### Avec Transactions

```go
import "time"

opts := &rete.TransactionOptions{
    SubmissionTimeout: 10 * time.Second,
    VerifyRetryDelay:  10 * time.Millisecond,
    MaxVerifyRetries:  5,
    VerifyOnCommit:    true,
}

if err := network.SubmitFactsWithOptions(facts, opts); err != nil {
    log.Fatal(err)
}
```

### Intégration HTTP

#### Serveur

```bash
# Démarrer le serveur
tsd server --port 8080

# Avec authentification
tsd server --port 8080 --auth-key-file api-key.txt

# Avec TLS
tsd server --port 8443 --tls-cert cert.pem --tls-key key.pem
```

#### Client HTTP (curl)

```bash
# Compiler un programme
curl -X POST http://localhost:8080/compile \
  -H "Content-Type: text/plain" \
  --data-binary @program.tsd

# Avec authentification
curl -X POST http://localhost:8080/compile \
  -H "X-API-Key: your-api-key" \
  --data-binary @program.tsd

# Health check
curl http://localhost:8080/health

# Métriques Prometheus
curl http://localhost:8080/metrics
```

#### Client Go

```go
import (
    "bytes"
    "io"
    "net/http"
)

func compileTSD(program string, apiKey string) error {
    url := "http://localhost:8080/compile"
    
    req, err := http.NewRequest("POST", url, bytes.NewBufferString(program))
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "text/plain")
    if apiKey != "" {
        req.Header.Set("X-API-Key", apiKey)
    }
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Println(string(body))
    
    return nil
}
```

---

## Guide Avancé

### Négation (NOT)

Match quand une condition est fausse :

```tsd
type User(email: string, verified: bool)
action sendVerification(email: string)

rule needsVerification : {u: User} / 
    NOT(u.verified)
    ==> sendVerification(u.email)

User(email: "user@example.com", verified: false)
```

### Jointures Complexes (3+ Variables)

Associez plusieurs types de faits :

```tsd
type User(id: string, name: string)
type Team(id: string, name: string, budget: number)
type Task(id: string, teamId: string, cost: number)
action assign(userId: string, teamId: string, taskId: string)

rule affordableAssignment : {u: User, t: Team, task: Task} /
    task.teamId == t.id AND
    task.cost <= t.budget
    ==> assign(u.id, t.id, task.id)

User(id: "U001", name: "Alice")
Team(id: "T001", name: "DevOps", budget: 5000)
Task(id: "TSK001", teamId: "T001", cost: 3000)
```

### Patterns Conditionnels Avancés

Combinez plusieurs conditions complexes :

```tsd
type Transaction(id: string, amount: number, country: string, suspicious: bool)
action flagForReview(id: string)
action autoApprove(id: string)

rule fraudDetection : {t: Transaction} /
    (t.amount > 10000 OR t.country IN ["XX", "YY"]) AND
    NOT(t.suspicious == false)
    ==> flagForReview(t.id)

rule normalTransaction : {t: Transaction} /
    t.amount <= 10000 AND
    t.country NOT IN ["XX", "YY"] AND
    t.suspicious == false
    ==> autoApprove(t.id)
```

### Calculs Complexes

Utilisez des expressions arithmétiques élaborées :

```tsd
type Invoice(id: string, subtotal: number, taxRate: number, discount: number)
action processPayment(invoiceId: string, finalAmount: number)

rule calculateFinal : {i: Invoice} / i.subtotal > 0 ==>
    processPayment(
        i.id,
        (i.subtotal * (1 + i.taxRate)) - i.discount
    )

Invoice(id: "INV001", subtotal: 100.00, taxRate: 0.20, discount: 10.00)
// Résultat: processPayment("INV001", 110.00)
```

### Optimisation et Performance

#### Partage de Nœuds Alpha

TSD partage automatiquement les conditions alpha identiques :

```tsd
// Ces deux règles partagent la condition "p.price > 100"
rule expensive1 : {p: Product} / p.price > 100 ==> action1(p.name)
rule expensive2 : {p: Product} / p.price > 100 ==> action2(p.name)
```

#### Partage de Nœuds Beta

Les jointures identiques sont partagées :

```tsd
// Ces règles partagent la jointure Customer-Order
rule vip1 : {c: Customer, o: Order} / c.id == o.customerId AND c.vip ==> action1()
rule vip2 : {c: Customer, o: Order} / c.id == o.customerId AND c.vip ==> action2()
```

#### Bonnes Pratiques de Performance

1. **Placer les conditions sélectives en premier** :
   ```tsd
   // Bon : condition sélective d'abord
   {p: Product} / p.price > 1000 AND p.inStock
   
   // Moins optimal
   {p: Product} / p.inStock AND p.price > 1000
   ```

2. **Éviter les calculs redondants** :
   ```tsd
   // Bon : calculer une fois
   {o: Order} / (o.price * o.quantity) > 1000 ==> process(o.price * o.quantity)
   
   // Moins bon : calcul dupliqué
   {o: Order} / (o.price * o.quantity) > 1000 ==> process((o.price * o.quantity))
   ```

3. **Utiliser les bons opérateurs** :
   ```tsd
   // CONTAINS pour recherche simple (plus rapide)
   subject CONTAINS "urgent"
   
   // MATCHES seulement si regex nécessaire
   subject MATCHES "^URGENT:.*important$"
   ```

---

## Cas d'Usage Pratiques

### E-Commerce : Gestion de Promotions

```tsd
type Product(id: string, name: string, price: number, category: string)
type Customer(id: string, name: string, loyaltyPoints: number)
type Cart(customerId: string, productId: string, quantity: number)
action applyDiscount(customerId: string, discountPercent: number)
action sendPromoCode(customerId: string, code: string)

// Remise fidélité
rule loyaltyDiscount : {c: Customer} /
    c.loyaltyPoints > 1000
    ==> applyDiscount(c.id, 15)

// Promotion catégorie
rule electronicsPromo : {cart: Cart, p: Product} /
    cart.productId == p.id AND
    p.category == "electronics" AND
    p.price > 500
    ==> sendPromoCode(cart.customerId, "TECH20")

// Achats en gros
rule bulkDiscount : {cart: Cart, p: Product} /
    cart.productId == p.id AND
    cart.quantity >= 10
    ==> applyDiscount(cart.customerId, 10)
```

### IoT : Surveillance Système

```tsd
type Sensor(id: string, location: string, temperature: number, humidity: number)
type Alert(sensorId: string, level: string, timestamp: string)
action notifyAdmin(message: string, level: string)
action shutdownSystem(location: string)

// Température critique
rule criticalTemp : {s: Sensor} /
    s.temperature > 80
    ==> shutdownSystem(s.location)

// Alertes multiples
rule multipleAlerts : {s: Sensor, a1: Alert, a2: Alert} /
    a1.sensorId == s.id AND
    a2.sensorId == s.id AND
    a1.level == "high" AND
    a2.level == "high"
    ==> notifyAdmin("Sensor " + s.id + " has multiple alerts", "critical")

// Conditions environnementales
rule environmentalRisk : {s: Sensor} /
    s.temperature > 70 AND s.humidity < 30
    ==> notifyAdmin("Fire risk at " + s.location, "high")
```

### Validation de Données

```tsd
type User(email: string, age: number, country: string)
type Subscription(userId: string, plan: string, price: number)
action rejectUser(email: string, reason: string)
action approveSubscription(userId: string, plan: string)

// Validation email
rule validateEmail : {u: User} /
    NOT(u.email MATCHES "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
    ==> rejectUser(u.email, "Invalid email format")

// Validation âge légal
rule legalAge : {u: User} /
    u.age < 18 AND u.country NOT IN ["US", "UK"]
    ==> rejectUser(u.email, "Must be 18 or older")

// Validation abonnement
rule premiumEligibility : {u: User, s: Subscription} /
    s.userId == u.email AND
    s.plan == "premium" AND
    u.age >= 18
    ==> approveSubscription(u.email, s.plan)
```

### Workflow Métier

```tsd
type Document(id: string, status: string, author: string, approvals: number)
type Approver(name: string, role: string, department: string)
action requestApproval(docId: string, approver: string)
action publishDocument(docId: string)

// Approbation hiérarchique
rule managerApproval : {d: Document, a: Approver} /
    d.status == "pending" AND
    a.role == "manager" AND
    d.approvals < 2
    ==> requestApproval(d.id, a.name)

// Publication automatique
rule autoPublish : {d: Document} /
    d.approvals >= 3 AND
    d.status != "published"
    ==> publishDocument(d.id)
```

---

## Bonnes Pratiques

### Organisation du Code

#### Structure Multi-Fichiers

```
project/
├── types/
│   ├── user.tsd
│   ├── order.tsd
│   └── product.tsd
├── rules/
│   ├── validation.tsd
│   ├── pricing.tsd
│   └── promotions.tsd
├── facts/
│   └── initial-data.tsd
└── main.tsd
```

Exécuter :
```bash
tsd types/*.tsd rules/*.tsd facts/*.tsd main.tsd
```

#### Nommage

```tsd
// Types : PascalCase
type UserAccount(...)
type OrderItem(...)

// Variables : camelCase
{user: UserAccount, order: OrderItem}

// Actions : camelCase avec verbe
action processPayment(...)
action sendNotification(...)

// Règles : camelCase descriptif
rule calculateTotalPrice : ...
rule applyLoyaltyDiscount : ...
```

### Gestion des Erreurs

#### Validation des Données

```tsd
// Valider avant traitement
type Input(value: string, validated: bool)
action process(value: string)
action reject(value: string)

rule validateFirst : {i: Input} /
    i.validated == false
    ==> reject(i.value)

rule processValid : {i: Input} /
    i.validated == true
    ==> process(i.value)
```

#### Logging et Debug

```bash
# Activer le logging debug
export TSD_LOG_LEVEL=debug
tsd program.tsd

# Logging détaillé
export TSD_LOG_LEVEL=trace
tsd program.tsd
```

### Sécurité

#### Authentification API

```bash
# Générer une clé API
tsd auth generate-key --output api-key.txt

# Utiliser la clé
tsd server --auth-key-file api-key.txt
```

#### HTTPS/TLS

```bash
# Générer certificat self-signed (dev)
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# Démarrer serveur HTTPS
tsd server --port 8443 --tls-cert cert.pem --tls-key key.pem
```

### Tests

#### Tests Unitaires

Créez des fichiers de test avec résultats attendus :

```tsd
// test-pricing.tsd
type Product(name: string, price: number)
action expensive(name: string)

rule highPrice : {p: Product} / p.price > 100 ==> expensive(p.name)

// Cas de test
Product(name: "Laptop", price: 999)   // Doit déclencher
Product(name: "Mouse", price: 25)     // Ne doit pas déclencher
```

```bash
# Exécuter et vérifier
tsd test-pricing.tsd | grep "expensive"
```

#### Tests d'Intégration

```bash
# Script de test
#!/bin/bash
OUTPUT=$(tsd program.tsd)
EXPECTED="ACTION EXÉCUTÉE: action1"

if echo "$OUTPUT" | grep -q "$EXPECTED"; then
    echo "✅ Test passed"
    exit 0
else
    echo "❌ Test failed"
    exit 1
fi
```

### Documentation

#### Commentaires Descriptifs

```tsd
// Type représentant un utilisateur du système
// Tous les champs sont obligatoires
type User(
    email: string,    // Format: user@domain.com
    age: number,      // Doit être >= 0
    active: bool      // true si compte actif
)

// Règle métier : Vérification d'éligibilité
// Conditions:
//   - Age >= 18
//   - Compte actif
//   - Email valide
rule checkEligibility : {u: User} /
    u.age >= 18 AND
    u.active == true AND
    u.email MATCHES "^.+@.+\\..+$"
    ==> approveUser(u.email)
```

---

## Dépannage

### Problèmes Courants

#### Règle Non Déclenchée

**Symptôme :** La règle ne s'exécute pas.

**Vérifications :**
1. Vérifier la syntaxe des conditions (`==` pas `=`)
2. Vérifier que les types de faits correspondent
3. Activer le debug : `TSD_LOG_LEVEL=debug tsd program.tsd`
4. Vérifier les types (nombres vs chaînes)

#### Erreurs de Type

**Symptôme :** "type mismatch" ou "invalid operation"

**Solutions :**
1. Utiliser des casts explicites : `(string)value`
2. Vérifier les définitions de types
3. Pour concaténation : convertir en string

#### Pattern Ne Correspond Pas

**Symptôme :** Jointure multi-faits ne matche pas

**Solutions :**
1. Vérifier que tous les faits existent
2. Vérifier les conditions de jointure
3. Tester chaque pattern séparément

### Performance

#### Règles Lentes

**Diagnostic :**
```bash
# Profiling
TSD_LOG_LEVEL=debug tsd program.tsd 2>&1 | grep "evaluation time"
```

**Optimisations :**
1. Placer conditions sélectives en premier
2. Éviter calculs complexes dans conditions
3. Utiliser CONTAINS au lieu de MATCHES si possible

---

## Ressources

### Documentation

- [Installation](installation.md) - Guide d'installation
- [Architecture](architecture.md) - Comprendre l'algorithme RETE
- [Configuration](configuration.md) - Configuration avancée
- [API](api.md) - API programmatique Go
- [Référence](reference.md) - Référence complète

### Exemples

```bash
# Explorer les exemples
ls examples/
tsd examples/basic-rules.tsd
tsd examples/type-casting.tsd
```

### Aide

- **GitHub Issues :** https://github.com/treivax/tsd/issues
- **Debug :** `TSD_LOG_LEVEL=debug tsd program.tsd`
- **Documentation :** `tsd --help`

---

**Bon développement avec TSD ! 🚀**