# Référence Complète TSD

**Documentation de référence** - API HTTP, Grammaire, Authentification, Logging, Contribution

---

## Table des Matières

1. [API HTTP/REST](#api-httprest)
2. [Grammaire TSD](#grammaire-tsd)
3. [Authentification](#authentification)
4. [Logging](#logging)
5. [Contribution](#contribution)

---

## API HTTP/REST

### Vue d'Ensemble

Le serveur TSD expose une API REST pour compiler et exécuter des programmes TSD à distance.

### Démarrage du Serveur

```bash
# Serveur HTTP basique
tsd server --port 8080

# Serveur avec authentification
tsd server --port 8080 --auth-key-file api-key.txt

# Serveur HTTPS avec TLS
tsd server --port 8443 --tls-cert cert.pem --tls-key key.pem

# Serveur HTTPS avec JWT
tsd server --port 8443 --tls-cert cert.pem --tls-key key.pem --jwt-secret-file secret.txt
```

### Endpoints

#### POST /compile

Compile et exécute un programme TSD.

**Request:**
```http
POST /compile HTTP/1.1
Host: localhost:8080
Content-Type: text/plain
X-API-Key: your-api-key-here

type Person(name: string, age: number)
action greet(name: string)

rule adult : {p: Person} / p.age >= 18 ==> greet(p.name)

Person(name: "Alice", age: 25)
```

**Response Success (200 OK):**
```json
{
  "status": "success",
  "actions": [
    {
      "name": "greet",
      "args": ["Alice"]
    }
  ],
  "execution_time_ms": 15,
  "facts_processed": 1,
  "rules_evaluated": 1
}
```

**Response Error (400 Bad Request):**
```json
{
  "status": "error",
  "error": "syntax error at line 3: unexpected token",
  "line": 3,
  "column": 12
}
```

#### GET /health

Vérification de santé du serveur.

**Request:**
```http
GET /health HTTP/1.1
Host: localhost:8080
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime_seconds": 3600
}
```

#### GET /metrics

Métriques Prometheus pour monitoring.

**Request:**
```http
GET /metrics HTTP/1.1
Host: localhost:8080
```

**Response (200 OK):**
```
# HELP tsd_rules_total Total number of rules in the network
# TYPE tsd_rules_total gauge
tsd_rules_total 42

# HELP tsd_facts_total Total number of facts asserted
# TYPE tsd_facts_total counter
tsd_facts_total 1523

# HELP tsd_alpha_nodes_total Total number of alpha nodes
# TYPE tsd_alpha_nodes_total gauge
tsd_alpha_nodes_total 38

# HELP tsd_beta_nodes_total Total number of beta nodes
# TYPE tsd_beta_nodes_total gauge
tsd_beta_nodes_total 15
```

### Codes de Statut HTTP

| Code | Description | Quand |
|------|-------------|-------|
| 200 | OK | Compilation réussie |
| 400 | Bad Request | Erreur de syntaxe ou validation |
| 401 | Unauthorized | Authentification manquante ou invalide |
| 403 | Forbidden | Authentification valide mais accès refusé |
| 500 | Internal Server Error | Erreur serveur interne |
| 503 | Service Unavailable | Serveur surchargé ou en maintenance |

### Authentification

#### API Key (Header)

```http
POST /compile HTTP/1.1
Host: localhost:8080
X-API-Key: your-api-key-here
Content-Type: text/plain
```

#### JWT (Bearer Token)

```http
POST /compile HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: text/plain
```

### Exemples cURL

#### Compiler un Programme

```bash
curl -X POST http://localhost:8080/compile \
  -H "Content-Type: text/plain" \
  -H "X-API-Key: your-key" \
  --data-binary @program.tsd
```

#### Vérifier la Santé

```bash
curl http://localhost:8080/health
```

#### Récupérer les Métriques

```bash
curl http://localhost:8080/metrics
```

---

## Grammaire TSD

### Vue d'Ensemble

TSD utilise une grammaire PEG (Parsing Expression Grammar) pour définir sa syntaxe.

### Structure d'un Programme

```ebnf
Program         = Statement*
Statement       = TypeDef | ActionDecl | Rule | FactAssertion | Comment
Comment         = "//" [^\n]* | "/*" .* "*/"
```

### Types

#### Définition de Type

```ebnf
TypeDef         = "type" Identifier "(" FieldList ")"
FieldList       = Field ("," Field)*
Field           = Identifier ":" TypeName
TypeName        = "string" | "number" | "bool" | Identifier
```

**Exemples:**
```tsd
type Person(name: string, age: number)
type Product(id: string, price: number, available: bool)
type Order(customerId: string, items: number, total: number)
```

#### Types Primitifs

| Type | Description | Exemples |
|------|-------------|----------|
| `string` | Chaîne de caractères | `"hello"`, `"123"`, `""` |
| `number` | Nombre (int/float) | `42`, `3.14`, `-10`, `0.0` |
| `bool` | Booléen | `true`, `false` |

### Actions

```ebnf
ActionDecl      = "action" Identifier "(" ParamList? ")"
ParamList       = Param ("," Param)*
Param           = Identifier ":" TypeName
```

**Exemples:**
```tsd
action notify(message: string)
action processOrder(orderId: string, amount: number)
action sendEmail(to: string, subject: string, body: string)
```

### Règles

```ebnf
Rule            = "rule" Identifier ":" Pattern "/" Condition "==>" Action
Pattern         = "{" VarBinding ("," VarBinding)* "}"
VarBinding      = Identifier ":" TypeName
Condition       = Expression
Action          = Identifier "(" ArgList? ")"
ArgList         = Expr ("," Expr)*
```

**Exemples:**
```tsd
rule simple : {p: Person} / p.age >= 18 ==> greet(p.name)

rule complex : {c: Customer, o: Order} / 
    c.id == o.customerId AND o.total > 100 
    ==> applyDiscount(c.id, o.total * 0.1)
```

### Expressions

#### Littéraux

```ebnf
Literal         = String | Number | Boolean
String          = '"' [^"]* '"'
Number          = [0-9]+ ("." [0-9]+)?
Boolean         = "true" | "false"
```

#### Accès aux Champs

```ebnf
FieldAccess     = Identifier "." Identifier
```

**Exemples:**
```tsd
p.name
order.total
user.email
```

#### Opérateurs de Comparaison

```ebnf
Comparison      = Expr CompOp Expr
CompOp          = "==" | "!=" | "<" | ">" | "<=" | ">="
```

**Exemples:**
```tsd
age >= 18
price < 100
status != "pending"
```

#### Opérateurs Logiques

```ebnf
LogicalExpr     = Expr LogicalOp Expr | "NOT" "(" Expr ")"
LogicalOp       = "AND" | "OR"
```

**Exemples:**
```tsd
age >= 18 AND age <= 65
premium == true OR vip == true
NOT(verified)
```

#### Opérateurs Arithmétiques

```ebnf
ArithExpr       = Expr ArithOp Expr
ArithOp         = "+" | "-" | "*" | "/" | "%"
```

**Priorité (haut → bas):**
1. `*`, `/`, `%`
2. `+`, `-`

**Exemples:**
```tsd
price * quantity
total + tax
discount / 100
age % 10
```

#### Opérateurs de Chaînes

```ebnf
StringOp        = "CONTAINS" | "LIKE" | "MATCHES" | "IN"
```

**CONTAINS:**
```tsd
subject CONTAINS "urgent"
```

**LIKE (SQL-style):**
```tsd
filename LIKE "%.txt"        // % = any chars
code LIKE "ABC___"           // _ = single char
```

**MATCHES (Regex):**
```tsd
email MATCHES "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
```

**IN (Collection):**
```tsd
role IN ["admin", "superuser", "root"]
status IN ["active", "pending"]
```

#### Type Casting

```ebnf
Cast            = "(" TypeName ")" Expr
```

**Exemples:**
```tsd
(string)42              // "42"
(number)"123"           // 123
(bool)1                 // true
(string)price           // Convertir price en string
```

### Faits

```ebnf
FactAssertion   = TypeName "(" FieldAssignments ")"
FieldAssignments = FieldAssign ("," FieldAssign)*
FieldAssign     = Identifier ":" Literal
```

**Exemples:**
```tsd
Person(name: "Alice", age: 25)
Product(id: "P001", price: 99.99, available: true)
```

### Identifiants

```ebnf
Identifier      = [a-zA-Z_] [a-zA-Z0-9_]*
```

**Règles:**
- Commence par lettre ou underscore
- Contient lettres, chiffres, underscores
- Support Unicode (UTF-8)
- Case-sensitive (sauf mots-clés)

**Valides:**
```tsd
myVariable
_private
camelCase
PascalCase
snake_case
var123
αβγ
名前
```

**Invalides:**
```tsd
123start        // Commence par chiffre
my-var          // Tiret non autorisé
my.var          // Point réservé pour accès champ
```

### Mots-Clés Réservés

**Insensibles à la casse:**
```
type, action, rule
true, false
AND, OR, NOT
CONTAINS, LIKE, MATCHES, IN
string, number, bool
```

### Commentaires

```ebnf
LineComment     = "//" [^\n]*
BlockComment    = "/*" .* "*/"
```

**Exemples:**
```tsd
// Commentaire sur une ligne

/* Commentaire
   sur plusieurs
   lignes */

type Person(name: string)  // Commentaire en fin de ligne
```

### EBNF Complet

```ebnf
(* Programme TSD *)
Program         ::= Statement*
Statement       ::= TypeDef | ActionDecl | Rule | FactAssertion | Comment

(* Types *)
TypeDef         ::= "type" Identifier "(" FieldList ")"
FieldList       ::= Field ("," Field)*
Field           ::= Identifier ":" TypeName
TypeName        ::= "string" | "number" | "bool" | Identifier

(* Actions *)
ActionDecl      ::= "action" Identifier "(" ParamList? ")"
ParamList       ::= Param ("," Param)*
Param           ::= Identifier ":" TypeName

(* Règles *)
Rule            ::= "rule" Identifier ":" Pattern "/" Condition "==>" Action
Pattern         ::= "{" VarBinding ("," VarBinding)* "}"
VarBinding      ::= Identifier ":" TypeName
Condition       ::= Expression

(* Expressions *)
Expression      ::= LogicalExpr | Comparison | ArithExpr | StringExpr | Primary
LogicalExpr     ::= Expression ("AND" | "OR") Expression | "NOT" "(" Expression ")"
Comparison      ::= Expression ("==" | "!=" | "<" | ">" | "<=" | ">=") Expression
ArithExpr       ::= Expression ("+" | "-" | "*" | "/" | "%") Expression
StringExpr      ::= Expression ("CONTAINS" | "LIKE" | "MATCHES" | "IN") Expression
Primary         ::= Literal | FieldAccess | Cast | "(" Expression ")"

(* Termes *)
Literal         ::= String | Number | Boolean | Collection
String          ::= '"' [^"]* '"'
Number          ::= [0-9]+ ("." [0-9]+)?
Boolean         ::= "true" | "false"
Collection      ::= "[" (Literal ("," Literal)*)? "]"
FieldAccess     ::= Identifier ("." Identifier)*
Cast            ::= "(" TypeName ")" Expression

(* Actions *)
Action          ::= Identifier "(" ArgList? ")"
ArgList         ::= Expression ("," Expression)*

(* Faits *)
FactAssertion   ::= TypeName "(" FieldAssignments ")"
FieldAssignments ::= FieldAssign ("," FieldAssign)*
FieldAssign     ::= Identifier ":" Literal

(* Lexique *)
Identifier      ::= [a-zA-Z_] [a-zA-Z0-9_]*
Whitespace      ::= [ \t\n\r]+
Comment         ::= "//" [^\n]* | "/*" .* "*/"
```

---

## Authentification

### Vue d'Ensemble

TSD supporte deux mécanismes d'authentification :
1. **API Keys** : Clés statiques pour authentification simple
2. **JWT (JSON Web Tokens)** : Tokens signés pour authentification avancée

### API Keys

#### Générer une Clé API

```bash
# Générer et sauvegarder dans fichier
tsd auth generate-key --output api-key.txt

# Afficher dans stdout
tsd auth generate-key
```

**Format généré:**
```
tsd_ak_1a2b3c4d5e6f7g8h9i0j
```

**Structure:**
- Préfixe: `tsd_ak_` (identifie une API key TSD)
- Body: 20 caractères alphanumériques
- Longueur totale: 27 caractères

#### Valider une Clé API

```bash
# Valider depuis fichier
tsd auth validate-key --key-file api-key.txt

# Valider clé directe
tsd auth validate-key --key tsd_ak_1a2b3c4d5e6f7g8h9i0j
```

**Sortie Success:**
```
✅ API Key valide
```

**Sortie Erreur:**
```
❌ API Key invalide: format incorrect
```

#### Utiliser une API Key (Serveur)

```bash
# Démarrer serveur avec API key
tsd server --port 8080 --auth-key-file api-key.txt

# Ou avec variable d'environnement
export TSD_API_KEY=$(cat api-key.txt)
tsd server --port 8080
```

#### Utiliser une API Key (Client)

**HTTP Header:**
```http
POST /compile HTTP/1.1
Host: localhost:8080
X-API-Key: tsd_ak_1a2b3c4d5e6f7g8h9i0j
Content-Type: text/plain
```

**cURL:**
```bash
curl -X POST http://localhost:8080/compile \
  -H "X-API-Key: $(cat api-key.txt)" \
  --data-binary @program.tsd
```

**Client TSD:**
```bash
tsd client --url http://localhost:8080 \
  --api-key $(cat api-key.txt) \
  program.tsd
```

### JWT (JSON Web Tokens)

#### Générer un JWT

```bash
# Avec secret depuis fichier
tsd auth generate-jwt --secret-file secret.txt --user admin --output token.txt

# Avec secret direct
tsd auth generate-jwt --secret "my-secret-key" --user admin --expiry 24h
```

**Options:**
- `--user` : Nom d'utilisateur (requis)
- `--expiry` : Durée de validité (défaut: 24h)
  - Formats: `1h`, `24h`, `7d`, `30d`
- `--claims` : Claims additionnels (JSON)

**Exemple avec claims:**
```bash
tsd auth generate-jwt \
  --secret-file secret.txt \
  --user admin \
  --expiry 12h \
  --claims '{"role":"admin","department":"IT"}'
```

**Format JWT généré:**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhZG1pbiIsImV4cCI6MTcwNjI4MDAwMH0.signature
```

**Structure:**
```
Header.Payload.Signature
```

**Payload décodé:**
```json
{
  "sub": "admin",
  "exp": 1706280000,
  "iat": 1706193600,
  "role": "admin",
  "department": "IT"
}
```

#### Valider un JWT

```bash
# Valider depuis fichier
tsd auth validate-jwt --secret-file secret.txt --token-file token.txt

# Valider token direct
tsd auth validate-jwt --secret "my-secret-key" --token "eyJhbGc..."
```

**Sortie Success:**
```
✅ JWT valide
User: admin
Expires: 2024-12-20 15:30:00
Claims:
  role: admin
  department: IT
```

**Sortie Erreur:**
```
❌ JWT invalide: token expiré
```

#### Utiliser un JWT (Serveur)

```bash
# Démarrer serveur avec JWT
tsd server --port 8443 \
  --tls-cert cert.pem \
  --tls-key key.pem \
  --jwt-secret-file secret.txt
```

#### Utiliser un JWT (Client)

**HTTP Header (Bearer Token):**
```http
POST /compile HTTP/1.1
Host: localhost:8443
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: text/plain
```

**cURL:**
```bash
curl -X POST https://localhost:8443/compile \
  -H "Authorization: Bearer $(cat token.txt)" \
  --data-binary @program.tsd
```

**Client TSD:**
```bash
tsd client --url https://localhost:8443 \
  --jwt-token $(cat token.txt) \
  program.tsd
```

### Comparaison API Key vs JWT

| Critère | API Key | JWT |
|---------|---------|-----|
| **Simplicité** | ✅ Très simple | ⚠️ Plus complexe |
| **Expiration** | ❌ Pas d'expiration | ✅ Expiration configurée |
| **Claims** | ❌ Pas de métadonnées | ✅ Claims personnalisés |
| **Révocation** | ✅ Facile (supprimer clé) | ⚠️ Nécessite blacklist |
| **Sécurité** | ⚠️ Clé statique | ✅ Signature cryptographique |
| **Use Case** | Scripts, CI/CD | Applications, multi-users |

### Bonnes Pratiques

#### Stockage des Secrets

```bash
# ❌ MAUVAIS : Stocker en clair dans code
api_key = "tsd_ak_1a2b3c4d5e6f7g8h9i0j"

# ✅ BON : Variables d'environnement
export TSD_API_KEY=$(cat /secure/api-key.txt)

# ✅ BON : Fichiers protégés
chmod 600 /secure/api-key.txt
chown app:app /secure/api-key.txt
```

#### Rotation des Clés

```bash
# Générer nouvelle clé
tsd auth generate-key --output new-key.txt

# Déployer nouvelle clé
# Garder ancienne clé active pendant transition
# Supprimer ancienne clé après migration clients
```

#### Gestion des JWT

```bash
# Courte durée pour tokens utilisateur
tsd auth generate-jwt --expiry 1h

# Longue durée pour tokens service
tsd auth generate-jwt --expiry 30d

# Refresh tokens avant expiration
```

### Sécurité

#### TLS/HTTPS Obligatoire en Production

```bash
# ❌ MAUVAIS : HTTP en production
tsd server --port 8080 --auth-key-file key.txt

# ✅ BON : HTTPS en production
tsd server --port 8443 \
  --tls-cert cert.pem \
  --tls-key key.pem \
  --auth-key-file key.txt
```

#### Secrets Forts

```bash
# API Key: 20+ caractères aléatoires
tsd auth generate-key

# JWT Secret: 32+ caractères aléatoires
openssl rand -base64 32 > secret.txt
```

#### Limiter les Permissions

```bash
# Fichier secrets: lecture seule propriétaire
chmod 400 secret.txt
chown app:app secret.txt

# Répertoire secrets: pas de listing
chmod 700 /secure/
```

---

## Logging

### Niveaux de Log

TSD utilise des niveaux de log hiérarchiques :

| Niveau | Description | Quand utiliser |
|--------|-------------|----------------|
| `ERROR` | Erreurs critiques | Échecs bloquants |
| `WARN` | Avertissements | Situations anormales mais non-bloquantes |
| `INFO` | Informations générales | Événements importants (démarrage, arrêt) |
| `DEBUG` | Détails de debugging | Diagnostic détaillé |
| `TRACE` | Trace complète | Debug très détaillé (performance) |

### Configuration du Niveau de Log

#### Via Variable d'Environnement

```bash
# Niveau ERROR (production)
export TSD_LOG_LEVEL=error
tsd program.tsd

# Niveau INFO (défaut)
export TSD_LOG_LEVEL=info
tsd program.tsd

# Niveau DEBUG (développement)
export TSD_LOG_LEVEL=debug
tsd program.tsd

# Niveau TRACE (debug détaillé)
export TSD_LOG_LEVEL=trace
tsd program.tsd
```

#### Via Flag CLI

```bash
# Équivalent
tsd --log-level debug program.tsd
```

#### Via Code Go

```go
import "github.com/treivax/tsd/tsdio"

// Créer logger
logger := tsdio.NewLogger()
logger.SetLevel(tsdio.LevelDebug)

// Utiliser
logger.Info("Démarrage du réseau")
logger.Debug("Condition normalisée: %s", normalized)
logger.Error("Échec transaction: %v", err)
```

### Format des Logs

#### Format par Défaut

```
2024-12-20 15:30:45 [INFO] Démarrage du serveur TSD sur port 8080
2024-12-20 15:30:46 [DEBUG] Network construit: 42 rules, 156 nodes
2024-12-20 15:30:47 [WARN] Cache alpha plein, éviction LRU
2024-12-20 15:30:48 [ERROR] Transaction échouée: timeout après 30s
```

**Structure:**
```
YYYY-MM-DD HH:MM:SS [LEVEL] Message
```

#### Format JSON (Structured Logging)

```bash
export TSD_LOG_FORMAT=json
tsd server --port 8080
```

**Sortie:**
```json
{"time":"2024-12-20T15:30:45Z","level":"INFO","msg":"Démarrage du serveur TSD","port":8080}
{"time":"2024-12-20T15:30:46Z","level":"DEBUG","msg":"Network construit","rules":42,"nodes":156}
{"time":"2024-12-20T15:30:47Z","level":"WARN","msg":"Cache alpha plein","eviction":"LRU"}
```

### Logging par Composant

#### Network RETE

```
[DEBUG] Building RETE network from 42 rules
[DEBUG] Alpha nodes created: 156
[DEBUG] Beta nodes created: 89
[DEBUG] Terminal nodes: 42
[INFO] Network ready (build time: 250ms)
```

#### Alpha Chains

```
[DEBUG] Alpha chain hash: 0x1a2b3c
[DEBUG] Alpha sharing enabled: reusing node alpha_Person_age>18
[TRACE] Condition normalized: age >= 18 -> normalized(age, 18, GTE)
```

#### Beta Chains

```
[DEBUG] Join node created: join_PersonOrder
[DEBUG] Join condition: p.id == o.customerId
[TRACE] Left memory size: 45 tokens
[TRACE] Right memory size: 123 tokens
[DEBUG] Join produced: 67 tokens
```

#### Transactions

```
[INFO] Transaction started: tx_1234
[DEBUG] Submitting 10 facts
[DEBUG] Verification started (retry 1/5)
[INFO] Transaction committed (duration: 150ms)
```

```
[WARN] Transaction retry 2/5: verification failed
[ERROR] Transaction rolled back: max retries exceeded
```

#### Storage

```
[DEBUG] Fact added: Person_123
[DEBUG] Fact retrieved: Person_123
[TRACE] Storage sync: 1523 facts persisted
```

### Exemples par Cas d'Usage

#### Développement

```bash
export TSD_LOG_LEVEL=debug
tsd program.tsd
```

**Sortie:**
```
[DEBUG] Parsing program.tsd
[DEBUG] Types found: 3 (Person, Order, Product)
[DEBUG] Actions found: 5
[DEBUG] Rules found: 12
[DEBUG] Building RETE network
[DEBUG] Alpha nodes: 24, Beta nodes: 18
[INFO] Asserting initial facts
[DEBUG] Fact asserted: Person_1
[DEBUG] Fact asserted: Order_1
[INFO] Execution complete (42 actions executed)
```

#### Production (Minimal)

```bash
export TSD_LOG_LEVEL=warn
tsd server --port 8080
```

**Sortie:**
```
[INFO] Server started on :8080
[WARN] Connection pool exhausted, queueing request
[ERROR] Failed to compile program: syntax error
```

#### Debug Problème

```bash
export TSD_LOG_LEVEL=trace
tsd program.tsd 2>&1 | grep "join_PersonOrder"
```

**Sortie:**
```
[TRACE] Activating join_PersonOrder with left token_45
[TRACE] Left bindings: {p: Person_1}
[TRACE] Right memory: 123 tokens
[TRACE] Testing join: p.id (1) == o.customerId (1)
[TRACE] Join successful, creating token_456
[TRACE] Result bindings: {p: Person_1, o: Order_1}
```

### Redirection des Logs

#### Fichier

```bash
# Redirection stdout
tsd program.tsd > output.log 2>&1

# Ou avec tee (afficher + sauvegarder)
tsd program.tsd 2>&1 | tee output.log
```

#### Syslog (Linux)

```bash
tsd program.tsd 2>&1 | logger -t tsd
```

#### Logrotate

```
/var/log/tsd/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
}
```

### Performance du Logging

#### Impact selon Niveau

| Niveau | Impact CPU | Impact I/O | Recommandé |
|--------|-----------|-----------|------------|
| ERROR | ~0.1% | Minimal | Production |
| WARN | ~0.5% | Faible | Production |
| INFO | ~1% | Moyen | Staging |
| DEBUG | ~5-10% | Élevé | Développement |
| TRACE | ~15-25% | Très élevé | Debug ponctuel |

#### Bonnes Pratiques

```bash
# ✅ Production: logs minimaux
export TSD_LOG_LEVEL=warn

# ✅ Staging: logs informatifs
export TSD_LOG_LEVEL=info

# ✅ Développement: logs détaillés
export TSD_LOG_LEVEL=debug

# ⚠️ Debug ponctuel seulement
export TSD_LOG_LEVEL=trace
```

---

## Contribution

### Introduction

Merci de votre intérêt pour contribuer à TSD ! Ce guide vous aidera à contribuer efficacement.

### Prérequis

#### Outils Requis

- **Go 1.21+** : Langage principal
- **Make** : Build et tâches automatisées
- **Git** : Gestion de version
- **golangci-lint** : Linting (optionnel mais recommandé)

#### Installation

```bash
# Go
go version  # Vérifier version >= 1.21

# Make
make --version

# golangci-lint (optionnel)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Premier Démarrage

#### 1. Fork et Clone

```bash
# Fork sur GitHub (via interface web)
# Puis cloner votre fork
git clone https://github.com/VOTRE_USERNAME/tsd.git
cd tsd

# Ajouter upstream
git remote add upstream https://github.com/treivax/tsd.git
```

#### 2. Compiler et Tester

```bash
# Compiler
make build

# Lancer les tests
make test

# Vérifier la couverture
make coverage

# Linter
make lint
```

### Standards de Code

#### Style Go

Suivre les conventions Go standard :

```go
// ✅ BON
func ProcessOrder(order *Order) error {
    if order == nil {
        return errors.New("order cannot be nil")
    }
    // ...
}

// ❌ MAUVAIS
func process_order(Order *Order) error {  // snake_case non idiomatique
    if Order == nil {  // Variable avec majuscule non idiomatique
        return errors.New("order cannot be nil")
    }
}
```

#### Nommage

**Packages:**
- Minuscules, un seul mot
- Éviter underscores

```go
// ✅ BON
package rete
package constraint

// ❌ MAUVAIS
package rete_engine
package constraint_parser
```

**Fonctions exportées:**
- PascalCase
- Commence par verbe

```go
// ✅ BON
func BuildNetwork() *Network
func ValidateRule() error

// ❌ MAUVAIS
func network_builder()
func rulevalidate()
```

**Variables:**
- camelCase (locales)
- PascalCase (exportées)

```go
// ✅ BON
var networkConfig Config
var MaxRetries int = 5

// ❌ MAUVAIS
var NetworkConfig Config  // Local ne devrait pas être exporté
var max_retries = 5       // snake_case
```

#### Documentation

**Fonctions exportées** : GoDoc en anglais

```go
// BuildNetwork constructs a RETE network from the given rules.
// It validates the rules, creates the node structure, and optimizes
// the network for performance.
//
// Returns an error if the rules are invalid or if construction fails.
func BuildNetwork(rules []Rule) (*Network, error) {
    // ...
}
```

**Commentaires internes** : Français

```go
func buildAlphaChain(cond Condition) *AlphaNode {
    // Normaliser la condition avant de créer le nœud
    normalized := normalize(cond)
    
    // Vérifier si un nœud existe déjà pour cette condition
    if existing := cache.Get(normalized); existing != nil {
        return existing
    }
    
    // ...
}
```

#### Tests

**Nommage:**
```go
func TestBuildNetwork(t *testing.T) { /* ... */ }
func TestBuildNetwork_WithInvalidRules(t *testing.T) { /* ... */ }
func TestAlphaNode_Activate(t *testing.T) { /* ... */ }
```

**Structure:**
```go
func TestProcessOrder(t *testing.T) {
    // Arrange (Setup)
    order := &Order{ID: "123", Total: 100}
    processor := NewProcessor()
    
    // Act (Execute)
    err := processor.Process(order)
    
    // Assert (Verify)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if order.Status != "processed" {
        t.Errorf("expected status 'processed', got %q", order.Status)
    }
}
```

**Table-driven tests:**
```go
func TestValidateAge(t *testing.T) {
    tests := []struct {
        name    string
        age     int
        wantErr bool
    }{
        {"valid adult", 25, false},
        {"valid senior", 70, false},
        {"invalid negative", -5, true},
        {"invalid zero", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateAge(tt.age)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateAge(%d) error = %v, wantErr %v", 
                    tt.age, err, tt.wantErr)
            }
        })
    }
}
```

### Workflow de Contribution

#### 1. Créer une Branche

```bash
# Synchroniser avec upstream
git checkout main
git pull upstream main

# Créer branche feature
git checkout -b feature/ma-nouvelle-fonctionnalite

# Ou branche bugfix
git checkout -b fix/corriger-bug-xyz
```

#### 2. Développer

```bash
# Faire vos modifications
# Tester régulièrement
make test

# Committer régulièrement
git add .
git commit -m "feat: ajouter support pour X"
```

#### 3. Tester

```bash
# Tests unitaires
make test

# Tests avec race detector
make test-race

# Couverture
make coverage

# Lint
make lint
```

#### 4. Préparer la Pull Request

```bash
# Rebaser sur main
git fetch upstream
git rebase upstream/main

# Pousser
git push origin feature/ma-nouvelle-fonctionnalite
```

#### 5. Créer la Pull Request

Sur GitHub :
1. Aller sur votre fork
2. Cliquer "Compare & pull request"
3. Remplir le template
4. Soumettre

### Messages de Commit

Format : **Conventional Commits**

```
<type>(<scope>): <description>

[corps optionnel]

[footer optionnel]
```

**Types:**
- `feat`: Nouvelle fonctionnalité
- `fix`: Correction de bug
- `docs`: Documentation
- `test`: Tests
- `refactor`: Refactoring
- `perf`: Amélioration performance
- `chore`: Tâches maintenance

**Exemples:**

```bash
# Feature
git commit -m "feat(rete): add support for accumulator nodes"

# Bugfix
git commit -m "fix(parser): handle empty constraint lists"

# Documentation
git commit -m "docs(api): update REST endpoint documentation"

# Tests
git commit -m "test(alpha): add tests for chain normalization"

# Breaking change
git commit -m "feat(storage)!: remove PostgreSQL backend

BREAKING CHANGE: PostgreSQL storage is no longer supported.
Use memory storage or file-based persistence instead."
```

### Pull Request Template

```markdown
## Description

[Description claire de ce que fait la PR]

## Type de Changement

- [ ] Bugfix (changement non-breaking corrigeant un bug)
- [ ] Feature (changement non-breaking ajoutant une fonctionnalité)
- [ ] Breaking change (fix ou feature cassant la compatibilité)
- [ ] Documentation

## Tests

- [ ] Tests unitaires ajoutés/mis à jour
- [ ] Tous les tests passent (`make test`)
- [ ] Couverture maintenue/améliorée
- [ ] Tests manuels effectués

## Checklist

- [ ] Code suit les standards du projet
- [ ] Commentaires ajoutés pour code complexe
- [ ] Documentation mise à jour
- [ ] Aucun warning de lint
- [ ] Commit messages suivent Conventional Commits

## Screenshots (si applicable)

[Captures d'écran]
```

### Review Process

#### Pour les Contributeurs

1. **Répondre aux commentaires** rapidement
2. **Faire les changements demandés** dans de nouveaux commits
3. **Ne pas forcer-push** après review (sauf demande explicite)
4. **Être patient et respectueux**

#### Pour les Reviewers

1. **Être constructif** et bienveillant
2. **Expliquer le "pourquoi"** des suggestions
3. **Approuver rapidement** si tout est OK
4. **Proposer des améliorations** sans bloquer

### Bonnes Pratiques

#### Code Quality

```bash
# Avant de committer
make lint        # Vérifier style
make test        # Lancer tests
make test-race   # Détecter race conditions
```

#### Performance

```bash
# Benchmarker les changements
go test -bench=. -benchmem ./...

# Profiler
go test -cpuprofile=cpu.prof -memprofile=mem.prof
go tool pprof cpu.prof
```

#### Documentation

- **GoDoc** pour toutes les fonctions exportées
- **README** pour nouveaux packages
- **Exemples** dans la doc quand pertinent
- **CHANGELOG** pour changements importants

### Communication

#### Channels

- **GitHub Issues**: Bugs, feature requests
- **Pull Requests**: Discussions de code
- **Discussions**: Questions générales

#### Reporting Bugs

Template issue :

```markdown
## Description

[Description claire et concise du bug]

## To Reproduce

Steps to reproduce:
1. Go to '...'
2. Click on '....'
3. See error

## Expected Behavior

[Ce qui devrait se passer]

## Actual Behavior

[Ce qui se passe réellement]

## Environment

- OS: [e.g. Ubuntu 22.04]
- Go version: [e.g. 1.21.5]
- TSD version: [e.g. 1.0.0]

## Additional Context

[Tout contexte additionnel, logs, screenshots]
```

### Ressources

#### Documentation

- [Installation](installation.md)
- [Guides](guides.md)
- [Architecture](architecture.md)
- [Configuration](configuration.md)

#### Liens Externes

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

**Merci de contribuer à TSD ! 🚀**