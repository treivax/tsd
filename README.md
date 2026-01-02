# 🎯 TSD - Type System Development

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Coverage](https://img.shields.io/badge/coverage-81.2%25-brightgreen.svg)](#test-coverage)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](#tests)
[![Go Conventions](https://github.com/treivax/tsd/workflows/🔍%20Go%20Conventions%20Validation/badge.svg)](https://github.com/treivax/tsd/actions)

**Moteur de règles haute performance basé sur l'algorithme RETE avec système d'authentification**

TSD est un système de règles métier moderne qui permet l'évaluation efficace de conditions complexes sur des flux de données. Il supporte les expressions de négation, les fonctions avancées et les patterns de correspondance. TSD inclut également un serveur HTTP avec authentification (Auth Key + JWT) et un client HTTP pour l'exécution distante.

## ✨ Fonctionnalités

- 🚀 **Moteur RETE optimisé** - Algorithme de pattern matching haute performance
- 🧠 **Expressions complexes** - Support complet des négations (`NOT`) et conditions composées
- 🔍 **Opérateurs avancés** - `CONTAINS`, `LIKE`, `MATCHES`, `IN`, fonctions `LENGTH()`, `ABS()`, `UPPER()`
- 📊 **Types fortement typés** - Système de types robuste avec validation
- 🆔 **Génération automatique d'IDs** - Clés primaires et IDs déterministes basés sur les données métier
- 📁 **Chargement multi-fichiers** - Organisation modulaire (schéma/règles/données) avec fusion incrémentale
- 🎯 **81.2% de couverture** - 100% des modules de production >80%, tests robustes et maintenables
- ⚡ **Performance** - <1ms par règle, optimisé pour le traitement en temps réel
- 🏷️ **Identifiants de règles** - Gestion fine des règles avec identifiants obligatoires
- 🔗 **Beta Sharing System** - Partage intelligent des nœuds (60-80% réduction mémoire)
- 📈 **Agrégations multi-sources** - AVG, SUM, COUNT, MIN, MAX sur jointures complexes
- 🔒 **Authentification** - Support Auth Key et JWT pour sécuriser l'accès au serveur
- 🔐 **TLS/HTTPS par défaut** - Communication sécurisée avec génération de certificats
- 🌐 **Architecture Client/Serveur** - Serveur HTTPS et client pour exécution distante
- 🔧 **Binaire unique** - Un seul binaire `tsd` pour tous les rôles (compiler, auth, client, server)
- 💾 **Stockage In-Memory** - Architecture pure en mémoire avec cohérence forte

> **⚠️ Note Architecture:** TSD utilise exclusivement du **stockage en mémoire** avec garanties de cohérence forte. Toutes les données sont conservées en RAM pour des performances maximales (~10,000-50,000 faits/sec). La persistance se fait via export de fichiers `.tsd` et la réplication réseau via Raft est prévue pour les versions futures. Voir [docs/INMEMORY_ONLY_MIGRATION.md](docs/INMEMORY_ONLY_MIGRATION.md) pour plus de détails.

---

## 🆕 Nouveautés v2.0

### 🎯 Affectations de Variables

Nommez des faits pour les réutiliser dans d'autres définitions :

```tsd
alice = User("alice", "alice@example.com", 30)
bob = User("bob", "bob@example.com", 25)

Login(alice, "SES-001", 1704067200)
Login(bob, "SES-002", 1704067260)
```

### 🔗 Comparaisons de Faits

Comparez directement des faits dans les règles :

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)

alice = User("alice", "alice@example.com")
Order(alice, "ORD-001", 150.00)

rule customerOrders : {u: User, o: Order} / o.customer == u ==> 
    Log("Order for " + u.username)
```

### 🔄 Actions CRUD Dynamiques

TSD supporte des actions natives pour modifier les faits en cours d'exécution :

```tsd
type Product(#id: string, name: string, stock: number, status: string)

// Mettre à jour un ou plusieurs champs
rule mark_low_stock : {p: Product} / p.stock < 10 AND p.status == "available" ==>
    Update(p, {status: "low_stock"})

// Créer de nouveaux faits dynamiquement
rule create_alert : {p: Product} / p.stock == 0 ==>
    Insert(Alert(id: p.id, level: "critical", message: "Stock épuisé"))

// Supprimer des faits du réseau RETE
rule remove_obsolete : {p: Product} / p.status == "obsolete" ==>
    Retract(p)
```

**Actions disponibles :**
- **Update(fact, {field: value, ...})** - Modifier un ou plusieurs champs d'un fait existant
- **Insert(Type(...))** - Créer un nouveau fait et l'insérer dans le réseau RETE
- **Retract(fact)** - Supprimer un fait du réseau RETE

Voir [docs/actions/README.md](docs/actions/README.md) pour la documentation complète des actions.

### 📐 Types de Faits dans les Champs

Les champs peuvent référencer d'autres types de faits :

```tsd
type Customer(#customerId: string, name: string)
type Product(#sku: string, name: string, price: number)
type Order(customer: Customer, #orderNumber: string, date: string)
type OrderLine(order: Order, product: Product, quantity: number)
```

### 🔒 Identifiants Internes (`_id_`)

Les identifiants sont maintenant **cachés et internes** :
- ✅ Générés automatiquement (déterministes)
- ❌ **Jamais accessibles** dans les expressions TSD
- ✅ Utilisés en interne par le moteur RETE
- ✅ Comparaisons de faits automatiques

**⚠️ Breaking Change** : Le champ `id` est devenu `_id_` et n'est plus accessible. Voir [Guide de Migration](docs/migration/from-v1.x.md).

### 📚 Documentation

- [Guide de Migration v1.x → v2.0](docs/migration/from-v1.x.md) - **Obligatoire pour migrer**
- [Identifiants Internes](docs/internal-ids.md) - Système `_id_` complet
- [Affectations de Faits](docs/user-guide/fact-assignments.md) - Utiliser les variables
- [Comparaisons de Faits](docs/user-guide/fact-comparisons.md) - Comparer les faits
- [Système de Types](docs/user-guide/type-system.md) - Types dans les champs

---

## 📝 Syntaxe des Règles

### Format Obligatoire (v2.0+)

Toutes les règles doivent maintenant avoir un identifiant unique :

```
rule <identifiant> : {variables} / conditions ==> action
```

### Exemples

```go
// Règle simple
rule r1 : {p: Person} / p.age >= 18 ==> adult(p.id)

// Règle avec jointure
rule check_order : {p: Person, o: Order} / 
    p.id == o.customer_id AND o.amount > 100 
    ==> premium_order(p.id, o.id)

// Règle avec agrégation
rule vip_check : {p: Person} / 
    SUM(o: Order / o.customer_id == p.id ; o.amount) >= 1000 
    ==> vip_customer(p.id)
```

**📖 Documentation complète :** [docs/rule_identifiers.md](docs/rule_identifiers.md)

**🔄 Migration :** Pour migrer vos règles existantes, utilisez :
```bash
bash scripts/add_rule_ids.sh
```

## 🆔 Clés Primaires et Génération d'IDs

TSD génère automatiquement des identifiants internes uniques et déterministes (`_id_`) pour tous les faits, basés sur des clés primaires.

⚠️ **Important** : Le champ `_id_` est **caché et réservé au système**. Vous ne pouvez **jamais** y accéder dans vos expressions TSD.

### Définition de Clés Primaires

Marquez les champs de clé primaire avec le préfixe `#` :

```tsd
// Clé primaire simple
type User(#username: string, email: string, role: string)

// Clé primaire composite
type Product(#category: string, #name: string, price: number)

// Sans clé primaire (génération par hash)
type LogEvent(timestamp: number, level: string, message: string)
```

### Format des IDs Générés (Internes)

**Clé simple** : `TypeName~valeur`
```tsd
alice = User("alice", "alice@example.com", "admin")
// ID interne (_id_): "User~alice"
```

**Clé composite** : `TypeName~valeur1_valeur2`
```tsd
Product("Electronics", "Laptop", 1200.00)
// ID interne (_id_): "Product~Electronics_Laptop"
```

**Sans clé primaire** : `TypeName~<hash-16-chars>`
```tsd
LogEvent(1704067200, "ERROR", "Connection failed")
// ID interne (_id_): "LogEvent~a1b2c3d4e5f6g7h8"
```

### Utilisation dans les Règles

❌ **INTERDIT** - Accéder à `_id_` :

```tsd
// ❌ ERREUR : _id_ est réservé et inaccessible
rule showId : {u: User} / u._id_ == "User~alice" ==> Log("Found")
```

✅ **CORRECT** - Utiliser les affectations et comparaisons :

```tsd
alice = User("alice", "alice@example.com")

// Comparer sur les champs métier
rule showUser : {u: User} / u.username == "alice" ==> Log("Found: " + u.username)

// Comparer des faits directement
type Order(customer: User, #orderNum: string, total: number)
order1 = Order(alice, "ORD-001", 150.00)

rule customerOrders : {u: User, o: Order} / o.customer == u ==> 
    Log("Order for: " + u.username)
```

### Échappement des Caractères

Les caractères spéciaux sont automatiquement échappés en interne :
- `~` → `%7E` (séparateur type/valeur)
- `_` → `%5F` (séparateur composite)
- `%` → `%25` (caractère d'échappement)
- ` ` → `%20` (espace)

**📖 Documentation complète :** 
- [Identifiants Internes](docs/internal-ids.md) - Système `_id_` complet
- [Guide de Migration v1.x → v2.0](docs/migration/from-v1.x.md) - **Breaking changes**

**🔍 Exemples :** Consultez `examples/pk_*.tsd` et `examples/fact_*.tsd` pour tous les cas d'usage

## 🚀 Installation Rapide

```bash
# Cloner le projet
git clone https://github.com/treivax/tsd.git
cd tsd

# Installation complète avec dépendances
make install

# Ou build rapide
make build
```

Le binaire unique `tsd` sera créé dans `./bin/tsd` et supporte tous les rôles :
- **Compilateur/Runner** (comportement par défaut)
- **Authentification** (`tsd auth ...`) - Inclut génération de certificats TLS
- **Client HTTPS** (`tsd client ...`)
- **Serveur HTTPS** (`tsd server ...`)

### Commandes Disponibles

```bash
# Construire le binaire unique TSD
make build

# Exécuter tous les tests (53 tests Alpha+Beta+Integration)
make rete-unified

# Tests unitaires Go
make test

# Formatage et analyse
make format lint

# Validation complète (format+lint+build+test)
make validate
```

## 🔐 TLS/HTTPS (Nouveau)

TSD utilise **HTTPS par défaut** pour toutes les communications client-serveur. Pour commencer :

### 1. Générer des Certificats (Développement)

```bash
# Générer des certificats auto-signés pour développement
tsd auth generate-cert

# Génère automatiquement dans ./certs/ :
# - server.crt (certificat serveur)
# - server.key (clé privée serveur)
# - ca.crt (certificat CA pour les clients)
```

Options avancées :
```bash
# Personnaliser les hôtes
tsd auth generate-cert -hosts "localhost,127.0.0.1,192.168.1.100"

# Personnaliser la durée de validité
tsd auth generate-cert -valid-days 730

# Répertoire personnalisé
tsd auth generate-cert -output-dir ./my-certs
```

### 2. Démarrer le Serveur (HTTPS)

```bash
# Mode sécurisé (par défaut, cherche ./certs/server.{crt,key})
tsd server

# Certificats personnalisés
tsd server --tls-cert /path/to/cert.crt --tls-key /path/to/key.key

# Mode HTTP non sécurisé (développement uniquement, déconseillé)
tsd server --insecure
```

Variables d'environnement :
```bash
export TSD_TLS_CERT=/path/to/cert.crt
export TSD_TLS_KEY=/path/to/key.key
export TSD_INSECURE=true  # pour mode HTTP
```

### 3. Utiliser le Client (HTTPS)

```bash
# HTTPS par défaut (avec certificat auto-signé)
tsd client program.tsd -insecure

# Ou avec vérification du CA
tsd client program.tsd -tls-ca ./certs/ca.crt

# Serveur distant avec certificat valide
tsd client program.tsd -server https://tsd.example.com:8080
```

Variables d'environnement :
```bash
export TSD_TLS_CA=./certs/ca.crt
export TSD_CLIENT_INSECURE=true  # désactive la vérification TLS
```

### ⚠️ Important - Sécurité

- **Développement** : Utilisez les certificats auto-signés générés par `tsd auth generate-cert`
- **Production** : Utilisez des certificats signés par une CA reconnue (Let's Encrypt, etc.)
- **Ne JAMAIS committer** les certificats/clés dans Git (déjà dans `.gitignore`)
- Le flag `--insecure` ne doit être utilisé qu'en développement

### Production avec Let's Encrypt

```bash
# Obtenir un certificat Let's Encrypt (exemple avec certbot)
sudo certbot certonly --standalone -d tsd.example.com

# Démarrer le serveur avec le certificat
tsd server \
  --tls-cert /etc/letsencrypt/live/tsd.example.com/fullchain.pem \
  --tls-key /etc/letsencrypt/live/tsd.example.com/privkey.pem
```

## 🛡️ Sécurité

### ⚠️ Reporting de Vulnérabilités

**Vous avez trouvé une vulnérabilité de sécurité ?** Ne créez **PAS** d'issue publique.

Consultez notre **[Security Policy](SECURITY.md)** pour :
- 🚨 Reporter une vulnérabilité de manière privée
- 📋 Connaître les versions supportées
- 🔄 Comprendre notre processus de gestion
- 🛡️ Suivre les best practices de déploiement

### Scan de Vulnérabilités

TSD intègre plusieurs outils de sécurité pour garantir la qualité et la sûreté du code.

### Scan de Vulnérabilités

**govulncheck** scanne automatiquement les vulnérabilités CVE dans les dépendances Go :

```bash
# Installer les outils de sécurité
make deps-dev

# Scan complet de sécurité (gosec + govulncheck)
make security-scan

# Scan vulnérabilités uniquement
make security-vulncheck

# Analyse statique uniquement
make security-gosec
```

**Intégration CI** : Le scan de vulnérabilités s'exécute automatiquement à chaque commit via GitHub Actions.

**Documentation complète** : [docs/security/VULNERABILITY_SCANNING.md](docs/security/VULNERABILITY_SCANNING.md)

### Outils de Sécurité

| Outil | Fonction | Documentation |
|-------|----------|---------------|
| **govulncheck** | Scan CVE dans dépendances | [VULNERABILITY_SCANNING.md](docs/security/VULNERABILITY_SCANNING.md) |
| **gosec** | Analyse statique sécurité | `.github/workflows/go-conventions.yml` |
| **go vet** | Analyse statique standard | Exécuté par `make lint` |

### En Cas de Vulnérabilité

Si govulncheck détecte une vulnérabilité :

1. **Ne pas merger** tant que non corrigée
2. **Mettre à jour Go** vers la version corrigée (si stdlib)
3. **Mettre à jour dépendances** (si externe)
4. **Re-scanner** avec `make security-vulncheck`

Voir la documentation complète pour plus de détails : [docs/security/VULNERABILITY_SCANNING.md](docs/security/VULNERABILITY_SCANNING.md)

## 📋 Usage

### Binaire Unique TSD

Le binaire `tsd` est multifonction et change de comportement selon son premier argument :

```bash
# Afficher l'aide globale
tsd --help

# Afficher la version
tsd --version

# Compiler/exécuter un programme (comportement par défaut)
tsd program.tsd
tsd -file program.tsd -v

# Gestion d'authentification et certificats
tsd auth generate-key
tsd auth generate-jwt -secret "mon-secret" -username alice
tsd auth validate -type jwt -token "..." -secret "mon-secret"
tsd auth generate-cert  # Générer certificats TLS

# Client HTTPS (par défaut)
tsd client program.tsd
tsd client program.tsd -insecure  # dev avec certificats auto-signés
tsd client -health -server https://tsd.example.com:8080

# Serveur HTTPS (par défaut)
tsd auth generate-cert  # d'abord générer les certificats
tsd server
tsd server -port 8443 -auth jwt -jwt-secret "mon-secret"
tsd server --insecure  # HTTP non sécurisé (déconseillé)
```

### Aide Spécifique par Rôle

```bash
tsd --help          # Aide globale
tsd auth --help     # Aide pour l'authentification
tsd client --help   # Aide pour le client HTTP
tsd server --help   # Aide pour le serveur HTTP
```

### Compilateur/Runner (Mode par Défaut)

Lorsqu'aucun rôle n'est spécifié, `tsd` fonctionne comme compilateur et runner :

```bash
# Compiler et valider un fichier TSD
tsd program.tsd

# Mode verbeux
tsd program.tsd -v

# Lire depuis stdin
cat program.tsd | tsd -stdin

# Code TSD directement
tsd -text 'type Person : <id: string, name: string>'
```

### Format de Fichier Unifié

À partir de la v3.0.0, TSD utilise une extension unique `.tsd` pour tous les fichiers. Un fichier `.tsd` peut contenir:
- **Définitions de types**: `type Person : <id: string, name: string>`
- **Assertions de faits**: `Person(id:p1, name:Alice)`
- **Règles**: `rule r1 : {p: Person} / p.name == "Alice" ==> match(p.id)`

### CLI Application - Pipeline Complet

Le binaire `tsd` exécute automatiquement le **pipeline RETE complet** (parsing → construction réseau → injection faits → évaluation):

```bash
# Validation seule (parsing + validation syntaxique)
./bin/tsd program.tsd

# Argument positionnel
./bin/tsd program.tsd

# Avec flag explicite
./bin/tsd -file program.tsd

# Mode verbeux (détails du réseau et actions)
./bin/tsd program.tsd -v

# Exemple avec un test
./bin/tsd beta_coverage_tests/join_simple.tsd -v

# Rétrocompatibilité (deprecated)
./bin/tsd -constraint rules.tsd  # affiche un warning
```

**Sortie typique:**
```
✅ Contraintes validées avec succès

📊 RÉSULTATS
============
Faits injectés: 10

🎯 ACTIONS DISPONIBLES: 3
  1. alert_action() - 2 bindings
  2. process_order() - 3 bindings
  3. validate_user() - 1 bindings

✅ Validation réussie
```

### Tests

Pour exécuter la suite de tests:

```bash
# Tests unitaires (rapides)
make test-unit

# Tests E2E (fixtures TSD)
make test-e2e

# Tests d'intégration
make test-integration

# Tous les tests
make test-all

# Via Makefile (anciennement rete-unified)
make rete-unified  # Exécute les tests E2E
```

### Test Coverage

**🎯 Couverture Globale : 81.2%** (code de production uniquement)

Le projet maintient une couverture de tests exceptionnelle avec **100% des modules de production au-dessus de 80%**.

#### Couverture par Module

| Module | Couverture | Statut |
|--------|-----------|--------|
| tsdio | 100.0% | ✅ Excellent |
| rete/internal/config | 100.0% | ✅ Excellent |
| auth | 94.5% | ✅ Excellent |
| constraint/internal/config | 90.8% | ✅ Excellent |
| internal/compilercmd | 89.7% | ✅ Excellent |
| constraint/cmd | 86.8% | ✅ Excellent |
| internal/authcmd | 85.5% | ✅ Excellent |
| internal/clientcmd | 84.7% | ✅ Excellent |
| cmd/tsd | 84.4% | ✅ Excellent |
| internal/servercmd | 83.4% | ✅ Excellent |
| constraint | 82.5% | ✅ Excellent |
| constraint/pkg/validator | 80.7% | ✅ Excellent |
| rete | 80.6% | ✅ Excellent |

#### Commandes de Couverture

```bash
# Couverture code production (sans exemples)
make coverage-prod

# Rapport détaillé avec analyse par module
make coverage-report

# Couverture complète (incluant exemples)
make coverage

# Couverture tests unitaires uniquement
make coverage-unit

# Couverture tests E2E uniquement
make coverage-e2e
```

#### Standards de Tests

Tous les tests respectent les standards définis dans `.github/prompts/test.md` :

- ✅ Tests déterministes (pas de flakiness)
- ✅ Tests isolés (cleanup complet)
- ✅ Structure table-driven
- ✅ Messages clairs avec émojis (✅ ❌ ⚠️)
- ✅ Couverture complète (cas nominaux, limites, erreurs)
- ✅ Pas de hardcoding (constantes nommées)

Pour plus de détails, voir les rapports dans `REPORTS/`:
- `TEST_COVERAGE_IMPROVEMENT_2025-01-15.md` (Phase 1)
- `TEST_COVERAGE_IMPROVEMENT_PHASE2_2025-12-15.md` (Phase 2)
- `TEST_COVERAGE_PHASE3_ANALYSIS_2025-12-15.md` (Phase 3)

### Exemple de Règle

```go
// Fichier: rules.constraint
type Account : <id: string, balance: number, active: bool>

// Règle: Détecter les comptes inactifs avec solde élevé
{a: Account} / NOT(a.active == true) AND a.balance > 1000
    ==> suspicious_account_alert(a.id, a.balance)
```

### API Programmatique

```go
import "github.com/treivax/tsd/constraint"

// Parser des contraintes
result, err := constraint.ParseConstraintFile("rules.constraint")
if err != nil {
    log.Fatal(err)
}

// Valider le programme
err = constraint.ValidateConstraintProgram(result)
if err != nil {
    log.Fatal(err)
}
```

## 🔒 Strong Mode - Cohérence Garantie

TSD utilise le **Strong Mode** par défaut pour garantir une cohérence stricte des données. Toutes les lectures reflètent les écritures les plus récentes avec vérification synchrone.

### Utilisation de Base (Configuration par Défaut)

```go
import "github.com/treivax/tsd/rete"

// Créer un réseau RETE
network := rete.NewReteNetwork(storage, logger)

// Utiliser la configuration par défaut (Strong mode)
tx := network.BeginTransaction()
defer tx.Rollback()

// Ajouter des faits
tx.AddFact("User", map[string]interface{}{
    "id": "user-123",
    "name": "Alice",
    "age": 30,
})

// Commit avec vérification automatique
err := tx.Commit()
if err != nil {
    log.Fatal("Transaction failed:", err)
}
```

### Configuration Personnalisée

```go
import (
    "time"
    "github.com/treivax/tsd/rete"
)

// Créer des options personnalisées
opts := rete.DefaultTransactionOptions()
opts.SubmissionTimeout = 15 * time.Second      // Timeout pour la soumission
opts.VerifyRetryDelay = 20 * time.Millisecond  // Délai entre les retries
opts.MaxVerifyRetries = 5                       // Nombre max de retries
opts.VerifyOnCommit = true                      // Vérifier au commit

// Utiliser la configuration personnalisée
tx := network.BeginTransactionWithOptions(opts)
defer tx.Rollback()

// ... ajouter des faits ...

err := tx.Commit()
```

### Configurations pour Stockage In-Memory

#### Configuration par Défaut (Single-Node)
```go
opts := rete.DefaultTransactionOptions()
// SubmissionTimeout: 30s
// VerifyRetryDelay:  50ms
// MaxVerifyRetries:  10
// VerifyOnCommit:    true
// Performance: ~10,000-50,000 faits/sec
```

#### Configuration Basse Latence
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 5 * time.Second,
    VerifyRetryDelay:  5 * time.Millisecond,
    MaxVerifyRetries:  3,
    VerifyOnCommit:    true,
}
// Performance: ~20,000-50,000 faits/sec
```

#### Configuration pour Réplication Future (Raft)
```go
opts := &rete.TransactionOptions{
    SubmissionTimeout: 30 * time.Second,
    VerifyRetryDelay:  50 * time.Millisecond,
    MaxVerifyRetries:  10,
    VerifyOnCommit:    true,
}
// Performance: ~1,000-10,000 faits/sec (selon réseau)
```

### Monitoring de Performance

```go
import "github.com/treivax/tsd/rete"

// Créer un collecteur de métriques
perfMetrics := rete.NewStrongModePerformanceMetrics()

// Pour chaque transaction
start := time.Now()
tx := network.BeginTransaction()
// ... opérations ...
err := tx.Commit()
duration := time.Since(start)

// Enregistrer les métriques
coherenceMetrics := tx.GetCoherenceMetrics()
perfMetrics.RecordTransaction(duration, factCount, err == nil, coherenceMetrics)

// Générer un rapport
fmt.Println(perfMetrics.GetReport())

// Vérifier la santé du système
if !perfMetrics.IsHealthy {
    log.Warn("Strong mode needs tuning")
    for _, rec := range perfMetrics.Recommendations {
        log.Info("Recommendation:", rec)
    }
}
```

### Garanties du Strong Mode

✅ **Cohérence Lecture-après-Écriture**: Toute lecture reflète les écritures les plus récentes  
✅ **Vérification Synchrone**: Chaque fait est vérifié avant de continuer  
✅ **Mécanisme de Retry**: Tentatives automatiques avec backoff exponentiel  
✅ **Transactions Atomiques**: Tous les faits sont persistés ou aucun  
✅ **Aucune Perte de Données**: Les échecs de stockage causent des échecs de transaction  

### Performances Attendues

- **In-Memory (Single-Node)**: ~10,000-50,000 faits/sec
- **In-Memory (Basse Latence)**: ~20,000-50,000 faits/sec
- **Future - Réplication Raft**: ~1,000-10,000 faits/sec
- **Latence moyenne**: 1-10ms par transaction

### Architecture de Stockage

TSD utilise exclusivement du **stockage en mémoire** avec garanties de cohérence forte:
- ✅ Cohérence lecture-après-écriture
- ✅ Vérification synchrone des faits
- ✅ Transactions atomiques
- ✅ Aucune perte de données en cas d'échec

**Persistance**: Export vers fichiers `.tsd`  
**Réplication**: Via protocole Raft (à venir)

### Documentation Complète

Pour plus d'informations:
- 📖 **Guide Utilisateur**: [`docs/USER_GUIDE.md`](docs/USER_GUIDE.md)
- 🏗️ **Architecture**: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)
- 🚀 **Démarrage Rapide**: [`docs/QUICK_START.md`](docs/QUICK_START.md)
- ✅ **Completion Report**: [`docs/PHASE4_STRONG_MODE_COMPLETION.md`](docs/PHASE4_STRONG_MODE_COMPLETION.md)

## 🏗️ Architecture

```
tsd/
├── cmd/
│   └── tsd/                    # CLI principal (binaire unique)
├── internal/                   # Packages internes
│   ├── compilercmd/            # Compilateur/Runner
│   ├── authcmd/                # Gestion d'authentification
│   ├── clientcmd/              # Client HTTP
│   └── servercmd/              # Serveur HTTP
├── constraint/                 # Parser PEG et validation
│   ├── grammar/                # Grammaire PEG
│   ├── parser.go               # Parser principal
│   └── *_test.go               # Tests unitaires
├── rete/                       # Moteur RETE
│   ├── rete.go                 # Nœuds RETE
│   ├── constraint_pipeline.go  # Pipeline complet
│   ├── evaluator.go            # Évaluation de conditions
│   ├── network.go              # Réseau RETE
│   ├── logger.go               # Système de logging
│   └── *_test.go               # Tests unitaires
├── tests/                      # Suite de tests organisée
│   ├── e2e/                    # Tests E2E (83 fixtures)
│   ├── integration/            # Tests d'intégration
│   ├── performance/            # Tests de performance
│   ├── fixtures/               # Fixtures TSD (alpha/beta/integration)
│   └── shared/testutil/        # Utilitaires de test partagés
└── docs/                       # Documentation
```

## 🧪 Tests

TSD utilise l'outillage Go standard avec une suite de tests organisée et complète.

### Quick Start

```bash
# Tests unitaires (rapides, <1s par package)
make test-unit

# Tests E2E (83 fixtures TSD)
make test-e2e

# Tests d'intégration (modules)
make test-integration

# Tests de performance
make test-performance

# Tous les tests
make test-all

# Rapport de couverture
make coverage
```

### Organisation des Tests

Le projet suit les conventions Go avec des build tags pour organiser les tests :

- **Unit Tests** : Tests rapides co-localisés avec le code (constraint/, rete/, cmd/)
- **Integration Tests** : Tests d'intégration entre modules (tests/integration/)
  - Intégration Constraint + RETE
  - Pipeline complet de compilation et exécution
  - Scénarios multi-modules complexes
- **E2E Tests** : 83 fixtures TSD validées (tests/e2e/)
  - 26 fixtures Alpha (opérations arithmétiques, comparaisons)
  - 26 fixtures Beta (jointures, patterns complexes)
  - 31 fixtures Integration (scénarios complets)
- **Performance Tests** : Load tests et benchmarks (tests/performance/)

### Couverture Complète

**✅ 83 fixtures TSD validées (100%)**

- **Alpha Tests (26)** : abs, addition, soustraction, multiplication, division, modulo, etc.
- **Beta Tests (26)** : Jointures, patterns multi-variables, contraintes complexes
- **Integration Tests (31)** : Pipeline complet, agrégations (AVG, SUM, COUNT, MIN, MAX)

### Commandes Avancées

```bash
# Tests par catégorie
make test-e2e-alpha        # Fixtures alpha uniquement
make test-e2e-beta         # Fixtures beta uniquement
make test-e2e-integration  # Fixtures integration uniquement

# Tests d'intégration détaillés
make test-integration-verbose   # Avec logs détaillés
make test-integration-coverage  # Avec rapport de couverture

# Performance et profiling
make test-load             # Tests de charge avec profiling
make bench                 # Benchmarks
make bench-performance     # Benchmarks de performance

# Couverture par type
make coverage-unit         # Couverture tests unitaires
make coverage-e2e          # Couverture tests E2E

# Tests avec race detector
make test-race

# Tests parallèles (configurable)
TEST_PARALLEL=8 make test-parallel
```

📖 **Documentation complète des tests** :
- [tests/README.md](tests/README.md) - Organisation générale
- [tests/integration/README.md](tests/integration/README.md) - Guide des tests d'intégration

### 🔐 Certificats de Test TLS

Les tests TLS nécessitent des certificats auto-signés. Ces certificats sont générés automatiquement lors de l'exécution des tests, mais vous pouvez aussi les générer manuellement :

```bash
# Générer les certificats de test (une seule fois)
cd tests/fixtures/certs
./generate_certs.sh
```

**⚠️ IMPORTANT** : Ces certificats sont **uniquement pour les tests** et ne doivent **jamais** être utilisés en production.

**Caractéristiques** :
- Certificats auto-signés RSA 2048 bits avec SHA-256
- Valides 365 jours pour localhost
- Générés automatiquement si manquants lors des tests
- Ignorés par Git (sécurité)

📖 Voir [tests/fixtures/certs/README.md](tests/fixtures/certs/README.md) pour plus de détails.

## 📖 Documentation

- [🗺️ **Index de Navigation**](DOCUMENTATION_INDEX.md) - **Guide complet pour naviguer dans la documentation**
- [📋 Guide Complet](docs/README.md) - Documentation complète
- [🎓 Tutoriel](docs/TUTORIAL.md) - Guide pas à pas de zéro à héros
- [✨ Fonctionnalités](docs/FEATURES.md) - Toutes les fonctionnalités du projet
- [⚡ Optimisations](docs/OPTIMIZATIONS.md) - Guide complet des optimisations
- [📚 API Reference](docs/API_REFERENCE.md) - Référence complète de l'API
- [📝 Guide de Logging](LOGGING_GUIDE.md) - Système de logging thread-safe
- [🔧 Guide Développeur](docs/development_guidelines.md) - Standards et bonnes pratiques

### 🏗️ Architecture et Diagrammes

Pour comprendre l'architecture du système avec des diagrammes visuels :

- [📊 **Diagrammes d'Architecture**](docs/architecture/diagrams/) - Collection complète de diagrammes Mermaid
  - [Architecture Globale](docs/architecture/diagrams/01-global-architecture.md) - Vue système, couches, dépendances
  - [Flux de Données](docs/architecture/diagrams/02-data-flow.md) - Séquences, propagation, compilation
  - [Moteur RETE](docs/architecture/diagrams/03-rete-architecture.md) - Nœuds Alpha/Beta, optimisations
  - [Sécurité](docs/architecture/diagrams/04-security-flow.md) - Authentification, TLS, JWT
  - [Modèle de Données](docs/architecture/diagrams/05-data-model.md) - Types, règles, contraintes
- [🏗️ Vue d'Ensemble Système](docs/architecture/SYSTEM_OVERVIEW.md) - Documentation architecture textuelle
- [📐 Architecture Détaillée](docs/architecture.md) - Spécifications complètes

> **💡 Nouveaux contributeurs ?** Commencez par [Architecture Globale](docs/architecture/diagrams/01-global-architecture.md) puis [Flux de Données](docs/architecture/diagrams/02-data-flow.md)

> **Note** : Les rapports générés par l'assistant IA sont stockés dans `REPORTS/` (non versionné).

## 📝 Logging

TSD fournit un système de logging thread-safe avec plusieurs niveaux de verbosité, optimisé pour la production et les tests.

### Configuration Rapide

```go
import "github.com/treivax/tsd/rete"

// Logger par défaut (Info level)
logger := rete.NewLogger(rete.LogLevelInfo, os.Stdout)

// Configuration du network
network := rete.NewReteNetwork(storage)
network.SetLogger(logger)

// Personnalisation
logger.SetLevel(rete.LogLevelDebug)
logger.SetTimestamps(true)
```

### Niveaux de Log

- `LogLevelSilent` (0) - Aucune sortie
- `LogLevelError` (1) - Erreurs critiques uniquement
- `LogLevelWarn` (2) - Avertissements
- `LogLevelInfo` (3) - Informations générales (défaut)
- `LogLevelDebug` (4) - Détails de débogage

### Utilisation dans les Tests

```go
func TestMyFeature(t *testing.T) {
    t.Parallel() // Safe avec TestEnvironment !

    env := rete.NewTestEnvironment(t,
        rete.WithLogLevel(rete.LogLevelDebug),
        rete.WithTimestamps(false),
    )
    defer env.Cleanup()

    // Utiliser les composants
    env.Network.SubmitFact(fact)

    // Inspecter les logs
    logs := env.GetLogs()
    assert.Contains(t, logs, "✅ Fait persisté")
}
```

### Bonnes Pratiques

- ✅ **Info** : Opérations majeures et résultats
- 🔍 **Debug** : Détails d'exécution et traces
- ⚠️ **Warn** : Situations sous-optimales
- ❌ **Error** : Erreurs critiques uniquement

**📖 Documentation complète :** [LOGGING_GUIDE.md](LOGGING_GUIDE.md)

## 🎯 Cas d'Usage Validés

### Expressions de Négation Complexes ✅

```go
// Exemple validé : Détecter les anomalies utilisateur
rule detect_anomaly : {u: User} / NOT(u.age >= 18 AND u.status != "blocked")
    ==> user_anomaly_detected(u.id, u.age, u.status)
```

**Résultat :** 100% de conformité sur 26 tests Alpha

### Patterns Avancés ✅

```go
// Validation d'emails d'entreprise
rule check_company_email : {e: Email} / e.address LIKE "%@company.com"
    ==> company_email_found(e.address)

// Codes conformes au format
rule validate_code : {c: Code} / c.value MATCHES "CODE[0-9]+"
    ==> valid_code_detected(c.value)
```

## 🔗 Beta Sharing System

**Nouveau système de partage intelligent des nœuds pour des performances exceptionnelles.**

### Gains de Performance

- 🎯 **60-80% de réduction des nœuds** - Élimination automatique des nœuds de jointure dupliqués
- 💾 **40-60% d'économie mémoire** - Workloads de production typiques
- ⚡ **30-50% plus rapide** - Compilation des règles avec cache basé sur hash
- ✅ **100% rétrocompatible** - Aucun changement de code nécessaire

### Quick Start (5 minutes)

```go
// Le partage beta est activé par défaut
network := rete.NewReteNetwork()

// Ajoutez vos règles normalement
network.AddRule(rule1)
network.AddRule(rule2) // Partage automatique avec rule1 si patterns similaires!

// Vérifiez les métriques
metrics := network.GetBetaMetrics()
fmt.Printf("Ratio de partage: %.1f%%\n", metrics.SharingRatio*100)
fmt.Printf("Nœuds créés: %d\n", metrics.TotalNodesCreated)
fmt.Printf("Nœuds réutilisés: %d\n", metrics.TotalNodesReused)
```

### Agrégations Multi-Sources

Support des agrégations complexes avec conditions de jointure:

```tsd
RULE high_value_customers
WHEN
  customer: Customer() /
  order: Order(customerId == customer.id) /
  item: OrderItem(orderId == order.id)
  total_spent: SUM(item.price * item.quantity) > 10000
  order_count: COUNT(order.id) > 5
  avg_order: AVG(order.amount) > 500
THEN
  MarkAsVIP(customer)
```

**Fonctions d'agrégation:** AVG, SUM, COUNT, MIN, MAX

### Documentation Complète

- 📖 [Quick Start Guide](rete/BETA_CHAINS_QUICK_START.md) - Démarrage en 5 minutes
- 🏗️ [Architecture Guide](rete/docs/BETA_SHARING_SYSTEM.md) - Conception complète
- 🚀 [Performance Guide](rete/MULTI_SOURCE_PERFORMANCE_GUIDE.md) - Optimisation avancée
- 📊 [Implementation Summary](rete/docs/BETA_IMPLEMENTATION_SUMMARY.md) - Résumé complet
- 🔧 [Lifecycle Management](rete/RULE_REMOVAL_WITH_JOINS_FEATURE.md) - Gestion du cycle de vie

### Exemples Réels

```bash
cd examples/multi_source_aggregations

# Analyse e-commerce
cat ecommerce_analytics.tsd

# Monitoring supply chain
cat supply_chain_monitoring.tsd

# Corrélation de capteurs IoT
cat iot_sensor_monitoring.tsd
```

### Profiling Automatisé

```bash
cd rete
./scripts/profile_multi_source.sh
# Génère: cpu.prof, mem.prof, profile_report.txt
```

## 📊 Performance

| Métrique | Valeur | Statut |
|----------|--------|---------|
| **Tests Passés** | 53/53 | ✅ 100% |
| **Couverture RETE** | 69.2% | ✅ Excellent |
| **Temps/Règle** | <1ms | ✅ Optimal |
| **Mémoire/Fait** | <100B | ✅ Efficient |
| **Throughput** | >10K faits/s | ✅ Élevé |
| **Réduction Nœuds** | 60-80% | ✅ Beta Sharing |
| **Économie Mémoire** | 40-60% | ✅ Beta Sharing |

### Benchmarks Beta Sharing

| Scénario | Réduction Nœuds | Gain Temps | Économie Mémoire |
|----------|-----------------|------------|------------------|
| E-commerce (5 règles) | 60% | 38% | 60% |
| Complexe (20 règles) | 60% | 45% | 60% |
| IoT Monitoring | 70% | 48% | 62% |
| Supply Chain | 62% | 38% | 55% |

### Optimisations Implémentées

- **Beta Sharing System** : Partage automatique des nœuds de jointure avec hash SHA-256
- **Join Result Cache** : Cache LRU avec TTL pour résultats de jointure
- **Hash Cache** : Mémoïsation des patterns avec invalidation automatique
- **Logger configurable** : Contrôle de verbosité en production (Silent/Error/Warn/Info/Debug)
- **Propagation RETE** : Tokens propagés efficacement sans calculs redondants
- **Extraction AST dynamique** : Aucun hardcoding, valeurs extraites du AST
- **Mémoire de travail optimisée** : Indexation par ID pour accès O(1)
- **Lifecycle Management** : Gestion sûre de suppression avec référence counting

## 🛠️ Scripts Utilitaires

```bash
# Build complet et tests
./scripts/build.sh

# Nettoyage
./scripts/clean.sh

# Validation des conventions Go
./scripts/validate_conventions.sh
```

## 🤝 Contribution

Nous accueillons les contributions ! Consultez [CONTRIBUTING.md](CONTRIBUTING.md) pour :

- 🛠️ **Setup environnement** - Installation et configuration complète
- ✅ **Standards de code** - Règles strictes et conventions
- 🧪 **Standards de tests** - Couverture, structure, bonnes pratiques
- 📝 **Process de PR** - Workflow complet de contribution
- 🔍 **Guidelines de review** - Ce qui est vérifié en review

**Quick Start :**
```bash
# Fork et clone
git clone https://github.com/VOTRE_USERNAME/tsd.git
cd tsd

# Installation complète
make install

# Validation avant commit
make validate
```

**Nouveau contributeur ?** Cherchez les issues [good first issue](../../issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22).

**Standards projet :** [.github/prompts/common.md](.github/prompts/common.md) ⭐

## 📈 Statut du Projet

**🟢 Production Ready**

- ✅ API stable et documentée
- ✅ 53/53 tests passés (100%)
- ✅ Agrégations sémantiquement validées
- ✅ Rétractation de faits implémentée
- ✅ Pipeline complet sans hardcoding
- ✅ Logger configurable pour production
- ✅ Performance validée

## 🎯 Fonctionnalités Avancées

### Rétractation de Faits ✅
Retrait dynamique de faits avec propagation automatique dans tout le réseau RETE.

### Agrégations Dynamiques ✅
AVG, SUM, COUNT, MIN, MAX avec extraction automatique des paramètres depuis l'AST.

### Nœuds Conditionnels ✅
EXISTS, NOT avec conditions de jointure complexes.

### Pipeline Unifié ✅
Un seul pipeline pour parsing, construction réseau, et exécution.

## 📄 License

Ce projet est sous licence MIT. Voir [LICENSE](LICENSE) pour le texte complet de la licence.

### Third-Party Components

TSD utilise des composants open-source sous licences permissives. Voir [THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md) pour la liste complète des dépendances et leurs licences.

### Acknowledgments

- **Pigeon PEG Parser Generator** - Utilisé pour générer le parser de contraintes depuis la grammaire PEG (BSD-3-Clause)
- **Testify** - Framework de tests unitaires (MIT)
- **Algorithme RETE** - Développé par Charles Forgy (Carnegie Mellon University, 1974-1979)

Toutes les dépendances utilisent des licences permissives compatibles avec un usage commercial.

## 🏆 Réalisations

- **100% succès** sur 53 tests (Alpha + Beta + Integration)
- **Agrégations complètes** : AVG, SUM, COUNT, MIN, MAX validées sémantiquement
- **Rétractation de faits** : Propagation automatique dans tout le réseau
- **Zéro hardcoding** : Extraction dynamique depuis l'AST
- **Architecture RETE optimisée** : Propagation de tokens sans calculs redondants
- **Logger configurable** : 5 niveaux (Silent/Error/Warn/Info/Debug)
- **Pipeline unifié** : Construction réseau + injection de faits en une passe

---

**TSD v2.0** - Moteur de règles RETE complet avec agrégations 🚀
