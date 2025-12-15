# Installation et Démarrage Rapide

**Documentation TSD** - Guide complet d'installation et premiers pas

---

## Table des Matières

1. [Prérequis](#prérequis)
2. [Installation](#installation)
3. [Vérification](#vérification)
4. [Démarrage Rapide (5 minutes)](#démarrage-rapide-5-minutes)
5. [Concepts Fondamentaux](#concepts-fondamentaux)
6. [Patterns Courants](#patterns-courants)
7. [Fonctionnalités Avancées](#fonctionnalités-avancées)
8. [Modes d'Exécution](#modes-dexécution)
9. [Configuration](#configuration)
10. [Dépannage](#dépannage)
11. [Prochaines Étapes](#prochaines-étapes)

---

## Prérequis

### Obligatoire

- **Go 1.21 ou supérieur**
  ```bash
  go version  # Doit afficher 1.21 ou plus
  ```

### Optionnel

- **Docker** (pour déploiement conteneurisé)
- **Make** (pour commandes de commodité)

---

## Installation

### Méthode 1 : Depuis les Sources (Recommandé)

Installation la plus flexible, recommandée pour le développement.

#### 1. Cloner le Dépôt

```bash
git clone https://github.com/treivax/tsd.git
cd tsd
```

#### 2. Compiler le Binaire

```bash
# Avec Make (recommandé)
make build

# Ou avec Go directement
go build -o bin/tsd ./cmd/tsd
```

#### 3. Vérifier l'Installation

```bash
./bin/tsd --version
```

#### 4. Installation Système (Optionnel)

```bash
# Linux/macOS
sudo cp bin/tsd /usr/local/bin/

# Ou ajouter au PATH
export PATH=$PATH:$(pwd)/bin
```

### Méthode 2 : Via Go Install

Installation directe depuis le dépôt :

```bash
go install github.com/treivax/tsd/cmd/tsd@latest
```

Le binaire sera installé dans `$GOPATH/bin` (typiquement `~/go/bin`).

Assurez-vous que `$GOPATH/bin` est dans votre PATH :

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Méthode 3 : Docker

#### Construire depuis les Sources

```bash
# Depuis la racine du projet
docker build -t tsd:local .
```

#### Exécuter le Conteneur

```bash
# Lancer le serveur TSD
docker run -p 8080:8080 tsd:local server

# Lancer le compilateur (monter répertoire local)
docker run -v $(pwd):/workspace tsd:local /workspace/program.tsd
```

---

## Vérification

### Vérification de Base

```bash
# Vérifier la version
tsd --version

# Afficher l'aide
tsd --help
```

### Lancer les Tests

```bash
# Lancer tous les tests
go test ./...

# Avec couverture
go test -cover ./...

# Tests d'un package spécifique
go test ./rete
```

### Test avec un Exemple

```bash
# Créer un fichier de test simple
cat > test.tsd << 'EOF'
type Person(name: string, age: number)
action greet(name: string)

rule hello : {p: Person} / p.age >= 18 ==> greet(p.name)

Person(name: "Alice", age: 25)
EOF

# Exécuter le programme
tsd test.tsd
```

**Sortie attendue :**
```
🎯 ACTION EXÉCUTÉE: greet("Alice")
```

---

## Démarrage Rapide (5 minutes)

### Votre Premier Programme TSD

#### 1. Créer une Règle Simple

Créez un fichier nommé `hello.tsd` :

```tsd
// Définir un type
type Person(name: string, age: number)

// Définir une action
action greet(name: string)

// Définir une règle
rule welcome : {p: Person} / p.age >= 18 ==> greet(p.name)

// Ajouter des faits
Person(name: "Alice", age: 25)
Person(name: "Bob", age: 16)
```

#### 2. Exécuter le Programme

```bash
tsd hello.tsd
```

**Sortie :**
```
🎯 ACTION EXÉCUTÉE: greet("Alice")
```

Seule la salutation d'Alice est exécutée car Bob a moins de 18 ans !

---

## Concepts Fondamentaux

### 1. Types

Définissez la structure de vos données :

```tsd
type Product(name: string, price: number, inStock: bool)
type Order(id: string, quantity: number, total: number)
```

### 2. Faits

Créez des instances de vos types :

```tsd
Product(name: "Laptop", price: 999.99, inStock: true)
Product(name: "Mouse", price: 29.99, inStock: false)
```

### 3. Règles

Définissez la logique métier avec du pattern matching :

```tsd
// Structure : rule nom : {pattern} / condition ==> action
rule expensive : {p: Product} / p.price > 500 ==> markAsPremium(p.name)
```

**Pattern :** `{p: Product}` - Correspond aux faits de type Product, liés à la variable `p`  
**Condition :** `p.price > 500` - Filtre pour les produits chers  
**Action :** `markAsPremium(p.name)` - Exécute l'action avec le nom du produit

### 4. Actions

Déclarez les actions que les règles peuvent déclencher :

```tsd
action markAsPremium(name: string)
action sendEmail(to: string, subject: string, body: string)
action createInvoice(orderId: string, amount: number)
```

---

## Patterns Courants

### Pattern 1 : Conditions Multiples

Combinez les conditions avec `AND` et `OR` :

```tsd
type User(name: string, age: number, premium: bool)
action sendOffer(name: string)

rule targetUser : {u: User} / 
    u.age >= 18 AND u.age <= 65 AND u.premium == true 
    ==> sendOffer(u.name)

User(name: "Alice", age: 30, premium: true)
```

### Pattern 2 : Faits Multiples

Associez plusieurs faits ensemble :

```tsd
type Customer(id: string, name: string, vip: bool)
type Order(customerId: string, total: number)
action applyDiscount(customerName: string, orderId: string)

rule vipDiscount : {c: Customer, o: Order} / 
    c.id == o.customerId AND c.vip == true AND o.total > 100 
    ==> applyDiscount(c.name, o.customerId)

Customer(id: "C001", name: "Alice", vip: true)
Order(customerId: "C001", total: 250.00)
```

### Pattern 3 : Opérations sur Chaînes

Utilisez les opérateurs de chaînes pour le pattern matching :

```tsd
type Email(address: string, subject: string)
action flagAsSpam(address: string)

// Vérifier si le sujet contient "URGENT"
rule spamFilter : {e: Email} / 
    e.subject CONTAINS "URGENT" 
    ==> flagAsSpam(e.address)

Email(address: "spam@example.com", subject: "URGENT: Act now!")
```

### Pattern 4 : Conversion de Types (Type Casting)

Convertissez les valeurs entre types :

```tsd
type Product(name: string, price: number, quantity: number)
action notify(message: string)

// Convertir nombre en chaîne pour concaténation
rule priceAlert : {p: Product} / p.price > 100 ==> 
    notify("High price: $" + (string)p.price)

Product(name: "Laptop", price: 999.99, quantity: 5)
```

**Conversions disponibles :**
- `(number)value` - Convertir en nombre
- `(string)value` - Convertir en chaîne
- `(bool)value` - Convertir en booléen

### Pattern 5 : Opérations Arithmétiques

Effectuez des calculs dans les conditions et actions :

```tsd
type Order(id: string, price: number, quantity: number)
action createInvoice(orderId: string, total: number)

rule calculateTotal : {o: Order} / o.quantity > 0 ==> 
    createInvoice(o.id, o.price * o.quantity)

Order(id: "ORD001", price: 50.00, quantity: 3)
```

**Opérateurs supportés :** `+`, `-`, `*`, `/`, `%`

---

## Fonctionnalités Avancées

### Négation (NOT)

Correspond quand une condition est fausse :

```tsd
type User(email: string, verified: bool)
action sendVerificationEmail(email: string)

rule needsVerification : {u: User} / 
    NOT(u.verified) 
    ==> sendVerificationEmail(u.email)

User(email: "user@example.com", verified: false)
```

### Patterns de Chaînes

Utilisez `LIKE` pour des patterns style SQL ou `MATCHES` pour regex :

```tsd
type File(name: string, path: string)
action processImage(name: string)
action processConfig(name: string)

// LIKE: % = n'importe quels caractères, _ = caractère unique
rule imageFiles : {f: File} / 
    f.name LIKE "%.png" OR f.name LIKE "%.jpg" 
    ==> processImage(f.name)

// MATCHES: support regex complet
rule configFiles : {f: File} / 
    f.path MATCHES "^/etc/.+\\.conf$" 
    ==> processConfig(f.name)

File(name: "photo.png", path: "/images/photo.png")
```

### Opérations sur Collections

Vérifiez l'appartenance avec `IN` :

```tsd
type User(name: string, role: string)
action grantAccess(name: string)

rule adminAccess : {u: User} / 
    u.role IN ["admin", "superuser", "root"] 
    ==> grantAccess(u.name)

User(name: "Alice", role: "admin")
```

---

## Modes d'Exécution

### Mode Compilateur (Par Défaut)

```bash
# Exécuter un programme TSD
tsd program.tsd

# Avec sortie verbeuse
tsd -v program.tsd

# Avec logging debug
TSD_LOG_LEVEL=debug tsd program.tsd
```

### Mode Serveur

Démarrer TSD comme serveur HTTP :

```bash
# Démarrer le serveur
tsd server --port 8080

# Avec authentification
tsd server --port 8080 --auth-key-file api-key.txt

# Avec TLS
tsd server --port 8443 --tls-cert cert.pem --tls-key key.pem
```

### Mode Client

Envoyer des programmes à un serveur TSD :

```bash
# Envoyer un programme au serveur
tsd client --url http://localhost:8080 program.tsd

# Avec clé API
tsd client --url http://localhost:8080 --api-key YOUR_KEY program.tsd
```

### Mode Authentification

Gérer les clés API :

```bash
# Générer une clé API
tsd auth generate-key --output api-key.txt

# Générer un JWT
tsd auth generate-jwt --user admin --output token.txt

# Valider un token
tsd auth validate-token --token YOUR_TOKEN
```

---

## Configuration

### Rôles du Binaire

TSD est un binaire unifié avec plusieurs rôles :

```bash
# Compilateur/Exécuteur (par défaut)
tsd program.tsd

# Gestion de l'authentification
tsd auth generate-key --output api-key.txt

# Serveur HTTP
tsd server --port 8080

# Client HTTP
tsd client --url http://localhost:8080 program.tsd
```

### Variables d'Environnement

```bash
# Définir le niveau de log
export TSD_LOG_LEVEL=debug

# Définir le port du serveur
export TSD_PORT=8080

# Définir l'authentification
export TSD_API_KEY=your-api-key-here

# Activer TLS
export TSD_TLS_CERT=/path/to/cert.pem
export TSD_TLS_KEY=/path/to/key.pem
```

### Fichiers de Configuration

Créez un fichier de configuration pour des paramètres persistants :

```yaml
# config.yaml
server:
  port: 8080
  host: 0.0.0.0
  
logging:
  level: info
  output: stdout
  
authentication:
  enabled: true
  key_file: /etc/tsd/api-key.txt
  
tls:
  enabled: false
  cert_file: /etc/tsd/cert.pem
  key_file: /etc/tsd/key.pem
```

Charger la configuration :

```bash
tsd server --config config.yaml
```

---

## Dépannage

### Problèmes de Compilation

#### Version de Go Trop Ancienne

```
Error: go version go1.20.x is too old
```

**Solution :** Mettre à jour Go vers 1.21 ou supérieur :

```bash
# Linux/macOS
go install golang.org/dl/go1.21.0@latest
go1.21.0 download
```

#### Dépendances Manquantes

```
Error: cannot find package ...
```

**Solution :** Télécharger les dépendances :

```bash
go mod download
go mod tidy
```

### Problèmes d'Exécution

#### Permission Refusée

```
Error: permission denied when writing output
```

**Solution :** Vérifier les permissions ou exécuter avec les privilèges appropriés :

```bash
chmod +w output-directory/
# Ou
sudo tsd program.tsd
```

#### Port Déjà Utilisé

```
Error: bind: address already in use
```

**Solution :** Utiliser un port différent ou tuer le processus utilisant le port :

```bash
# Trouver le processus utilisant le port 8080
lsof -i :8080
# Ou
netstat -tulpn | grep 8080

# Tuer le processus
kill -9 <PID>

# Ou utiliser un port différent
tsd server --port 8081
```

#### Problèmes de Certificat TLS

```
Error: tls: failed to verify certificate
```

**Solution :** Vérifier les chemins et la validité du certificat :

```bash
# Vérifier le certificat
openssl x509 -in cert.pem -text -noout

# Générer un certificat auto-signé pour les tests
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

### Problèmes de Tests

#### Tests Échouent à Cause de Conditions de Course

```bash
# Lancer les tests avec le détecteur de race conditions
go test -race ./...
```

#### Nettoyer les Artefacts de Test

```bash
# Supprimer le cache de test
go clean -testcache

# Supprimer les artefacts de compilation
make clean
```

### Problèmes Courants

#### Règle Non Déclenchée

**Problème :** Votre règle ne s'exécute pas même si vous vous y attendez.

**Solutions :**
1. Vérifier la syntaxe de la condition (utiliser `==` et non `=`)
2. Vérifier que les types de faits correspondent au pattern
3. Ajouter le logging debug : `TSD_LOG_LEVEL=debug tsd program.tsd`
4. Vérifier les incompatibilités de types (nombres vs chaînes)

#### Erreurs de Type

**Problème :** Erreurs "type mismatch" ou "invalid operation".

**Solutions :**
1. Utiliser des conversions explicites : `(string)numberValue`
2. Vérifier les types de champs dans les définitions de types
3. Vérifier que les opérations arithmétiques utilisent des nombres
4. Pour la concaténation de chaînes, les deux opérandes doivent être des chaînes

#### Pattern Ne Correspond Pas

**Problème :** Un pattern multi-faits ne correspond pas.

**Solutions :**
1. S'assurer que tous les faits référencés existent
2. Vérifier les conditions de jointure (égalité de variables)
3. Vérifier que les types de faits sont corrects
4. Tester chaque composant du pattern séparément

### Obtenir de l'Aide

1. **Consulter la Documentation :**
   - [Guides](guides.md)
   - [Architecture](architecture.md)
   - [Référence](reference.md)

2. **Voir les Exemples :**
   ```bash
   ls examples/
   ```

3. **Activer le Logging Debug :**
   ```bash
   TSD_LOG_LEVEL=debug tsd program.tsd
   ```

4. **Signaler des Problèmes :**
   - GitHub Issues: https://github.com/treivax/tsd/issues
   - Inclure : version, OS, messages d'erreur, cas de reproduction minimal

---

## Prochaines Étapes

Après l'installation :

1. **Lire les [Guides](guides.md)** - Tutoriels détaillés et cas d'usage
2. **Explorer [Architecture](architecture.md)** - Comprendre le fonctionnement interne
3. **Consulter [Configuration](configuration.md)** - Configuration avancée du système
4. **Voir [API](api.md)** - API programmatique Go
5. **Consulter [Référence](reference.md)** - Référence complète (API HTTP, grammaire, auth)

---

## Structure de Projet

Organisez les projets plus importants :

```
my-project/
├── types/
│   ├── user.tsd       # Définitions de types
│   └── order.tsd
├── rules/
│   ├── validation.tsd # Règles métier
│   └── pricing.tsd
├── facts/
│   └── initial.tsd    # Données initiales
└── main.tsd           # Programme principal
```

Exécuter avec plusieurs fichiers :

```bash
tsd types/*.tsd rules/*.tsd facts/*.tsd main.tsd
```

---

## Désinstallation

### Supprimer le Binaire

```bash
# Si installé système
sudo rm /usr/local/bin/tsd

# Si installé via go install
rm $(go env GOPATH)/bin/tsd
```

### Supprimer les Sources

```bash
# Supprimer le dépôt cloné
rm -rf /path/to/tsd
```

### Supprimer les Images Docker

```bash
docker rmi tsd:local
```

### Nettoyer le Cache Go

```bash
go clean -cache -modcache -i -r
```

---

## Aide-Mémoire

```tsd
// Types
type Name(field: string, count: number, active: bool)

// Faits
Name(field: "value", count: 42, active: true)

// Actions
action doSomething(arg1: string, arg2: number)

// Règles
rule name : {x: Type} / condition ==> action(x.field)

// Opérateurs
x == y          // Égal
x != y          // Différent
x < y           // Inférieur
x > y           // Supérieur
x <= y          // Inférieur ou égal
x >= y          // Supérieur ou égal
x AND y         // ET logique
x OR y          // OU logique
NOT(x)          // NON logique
x + y           // Addition (nombres) ou concaténation (chaînes)
x - y           // Soustraction
x * y           // Multiplication
x / y           // Division
x % y           // Modulo
x CONTAINS y    // Chaîne contient
x IN [a, b]     // Dans collection
x LIKE "%.txt"  // Pattern style SQL
x MATCHES "^a"  // Pattern regex

// Conversions
(number)value   // Vers nombre
(string)value   // Vers chaîne
(bool)value     // Vers booléen

// Commentaires
// Commentaire une ligne
/* Commentaire
   multi-lignes */
```

---

**Bon développement avec TSD ! 🚀**