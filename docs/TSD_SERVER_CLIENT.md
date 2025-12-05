# TSD Server & Client

Documentation pour le serveur HTTP TSD et son client CLI.

## Vue d'ensemble

Le serveur TSD permet d'exécuter des programmes TSD à distance via une API REST. Le client TSD (`tsd-client`) permet de soumettre des programmes au serveur et d'afficher les résultats.

### Architecture

```
┌──────────────┐         HTTP POST          ┌──────────────┐
│              │  /api/v1/execute           │              │
│  tsd-client  ├───────────────────────────►│  tsd-server  │
│              │                             │              │
│  (CLI)       │◄───────────────────────────┤  (API REST)  │
│              │  JSON Response             │              │
└──────────────┘                             └──────┬───────┘
                                                    │
                                                    │ utilise
                                                    ▼
                                             ┌──────────────┐
                                             │              │
                                             │  RETE Engine │
                                             │              │
                                             └──────────────┘
```

## Serveur TSD (`tsd-server`)

### Installation

```bash
# Compiler le serveur
cd cmd/tsd-server
go build -o tsd-server

# Ou depuis la racine du projet
go build -o bin/tsd-server ./cmd/tsd-server
```

### Utilisation

```bash
# Démarrer le serveur sur le port par défaut (8080)
./tsd-server

# Démarrer sur un port spécifique
./tsd-server -port 9000

# Démarrer sur une interface spécifique
./tsd-server -host 127.0.0.1 -port 8080

# Mode verbeux
./tsd-server -v
```

### Options

| Option    | Description                          | Défaut    |
|-----------|--------------------------------------|-----------|
| `-host`   | Hôte d'écoute du serveur            | `0.0.0.0` |
| `-port`   | Port d'écoute du serveur            | `8080`    |
| `-v`      | Mode verbeux (logs détaillés)       | `false`   |

### API REST

#### POST `/api/v1/execute`

Exécute un programme TSD et retourne les résultats.

**Requête:**

```json
{
  "source": "type Person : <id: string, name: string>\n\naction notify : <message: string>\n\nrule person_rule : {p: Person} / p.name == \"Alice\" ==> notify(p.id)\n\nPerson(\"p1\", \"Alice\")",
  "source_name": "example.tsd",
  "verbose": false
}
```

**Réponse (succès):**

```json
{
  "success": true,
  "results": {
    "facts_count": 1,
    "activations_count": 1,
    "activations": [
      {
        "action_name": "notify",
        "arguments": [
          {
            "position": 0,
            "value": "p1",
            "type": "string"
          }
        ],
        "triggering_facts": [
          {
            "id": "p1",
            "type": "Person",
            "attributes": {
              "id": "p1",
              "name": "Alice"
            }
          }
        ],
        "bindings_count": 1
      }
    ]
  },
  "execution_time_ms": 15
}
```

**Réponse (erreur):**

```json
{
  "success": false,
  "error": "Erreur de parsing: syntax error at line 1",
  "error_type": "parsing_error",
  "execution_time_ms": 5
}
```

**Types d'erreurs:**

- `parsing_error` : Erreur de parsing du code TSD
- `validation_error` : Erreur de validation du programme
- `execution_error` : Erreur lors de l'exécution
- `server_error` : Erreur serveur (requête invalide, etc.)

#### GET `/health`

Vérifie l'état du serveur.

**Réponse:**

```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime_seconds": 3600,
  "timestamp": "2025-01-15T10:30:00Z"
}
```

#### GET `/api/v1/version`

Retourne la version du serveur.

**Réponse:**

```json
{
  "version": "1.0.0",
  "go_version": "go1.21.0"
}
```

## Client TSD (`tsd-client`)

### Installation

```bash
# Compiler le client
cd cmd/tsd-client
go build -o tsd-client

# Ou depuis la racine du projet
go build -o bin/tsd-client ./cmd/tsd-client
```

### Utilisation

```bash
# Exécuter un fichier TSD
./tsd-client program.tsd

# Exécuter avec un serveur distant
./tsd-client -server http://tsd.example.com:8080 program.tsd

# Exécuter du code TSD directement
./tsd-client -text 'type Person : <id: string, name: string>'

# Lire depuis stdin
cat program.tsd | ./tsd-client -stdin

# Mode verbeux
./tsd-client -v program.tsd

# Format JSON
./tsd-client -format json program.tsd

# Vérifier la santé du serveur
./tsd-client -health
```

### Options

| Option      | Description                              | Défaut                    |
|-------------|------------------------------------------|---------------------------|
| `-server`   | URL du serveur TSD                       | `http://localhost:8080`   |
| `-file`     | Fichier TSD à exécuter                   | -                         |
| `-text`     | Code TSD directement                     | -                         |
| `-stdin`    | Lire depuis l'entrée standard            | `false`                   |
| `-v`        | Mode verbeux (affiche plus de détails)   | `false`                   |
| `-format`   | Format de sortie (`text` ou `json`)      | `text`                    |
| `-timeout`  | Timeout des requêtes                     | `30s`                     |
| `-health`   | Vérifier la santé du serveur             | `false`                   |
| `-h`        | Afficher l'aide                          | `false`                   |

### Exemples

#### Exécuter un fichier simple

```bash
# Fichier example.tsd
cat > example.tsd <<EOF
type Person : <id: string, name: string>
action notify : <message: string>
rule person_rule : {p: Person} / p.name == "Alice" ==> notify(p.id)
Person("p1", "Alice")
EOF

# Exécuter
./tsd-client example.tsd
```

**Sortie:**

```
✅ EXÉCUTION RÉUSSIE
===================
Temps d'exécution: 15ms
Faits injectés: 1
Activations: 1

🎯 ACTIONS DÉCLENCHÉES
======================

1. Action: notify
   Arguments:
     [0] p1 (string)
```

#### Mode verbeux

```bash
./tsd-client -v example.tsd
```

**Sortie:**

```
📤 Envoi requête à http://localhost:8080/api/v1/execute...

✅ EXÉCUTION RÉUSSIE
===================
Temps d'exécution: 15ms
Faits injectés: 1
Activations: 1

🎯 ACTIONS DÉCLENCHÉES
======================

1. Action: notify
   Arguments:
     [0] p1 (string)
   Faits déclencheurs:
     [0] Person (id: p1)
         id: p1
         name: Alice
```

#### Format JSON

```bash
./tsd-client -format json example.tsd
```

**Sortie:**

```json
{
  "success": true,
  "results": {
    "facts_count": 1,
    "activations_count": 1,
    "activations": [
      {
        "action_name": "notify",
        "arguments": [
          {
            "position": 0,
            "value": "p1",
            "type": "string"
          }
        ],
        "triggering_facts": [
          {
            "id": "p1",
            "type": "Person",
            "attributes": {
              "id": "p1",
              "name": "Alice"
            }
          }
        ],
        "bindings_count": 1
      }
    ]
  },
  "execution_time_ms": 15
}
```

#### Utiliser un serveur distant

```bash
# Démarrer le serveur sur une machine
ssh server1 'cd tsd && ./tsd-server -host 0.0.0.0 -port 8080'

# Depuis une autre machine
./tsd-client -server http://server1:8080 program.tsd
```

#### Pipeline avec stdin

```bash
# Générer dynamiquement du code TSD et l'exécuter
cat <<EOF | ./tsd-client -stdin
type Order : <id: string, amount: number>
action process : <order_id: string>
rule order_rule : {o: Order} / o.amount > 100 ==> process(o.id)
Order("o1", 150)
Order("o2", 50)
Order("o3", 200)
EOF
```

#### Health check

```bash
# Vérifier que le serveur est accessible
./tsd-client -health

# Avec format JSON
./tsd-client -health -format json
```

## Cas d'usage

### 1. Microservices

Déployer le serveur TSD comme microservice et l'utiliser depuis d'autres services:

```bash
# Démarrer le serveur
docker run -p 8080:8080 tsd-server

# Depuis un autre service
curl -X POST http://tsd-server:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{
    "source": "type Event : <name: string>\nEvent(\"test\")"
  }'
```

### 2. CI/CD Pipeline

Valider des règles TSD dans un pipeline CI/CD:

```bash
# .gitlab-ci.yml
validate-rules:
  script:
    - tsd-client -server http://tsd-server:8080 rules.tsd
    - if [ $? -ne 0 ]; then exit 1; fi
```

### 3. Monitoring et Alerting

Exécuter des règles de monitoring périodiquement:

```bash
# Cron job
*/5 * * * * /usr/local/bin/tsd-client -server http://localhost:8080 /etc/tsd/monitoring.tsd
```

### 4. API Gateway

Utiliser le serveur TSD derrière un API Gateway:

```bash
# nginx.conf
location /tsd/ {
    proxy_pass http://localhost:8080/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Intégration Programmatique

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "github.com/treivax/tsd/tsdio"
)

func executeTSD(source string) (*tsdio.ExecuteResponse, error) {
    req := tsdio.ExecuteRequest{
        Source:     source,
        SourceName: "api-call",
        Verbose:    false,
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        "http://localhost:8080/api/v1/execute",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result tsdio.ExecuteResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}
```

### Python

```python
import requests

def execute_tsd(source: str) -> dict:
    url = "http://localhost:8080/api/v1/execute"
    payload = {
        "source": source,
        "source_name": "python-client",
        "verbose": False
    }
    
    response = requests.post(url, json=payload)
    return response.json()

# Utilisation
result = execute_tsd("""
type Person : <id: string, name: string>
Person("p1", "Alice")
""")

print(f"Success: {result['success']}")
if result['success']:
    print(f"Activations: {result['results']['activations_count']}")
```

### JavaScript/Node.js

```javascript
const axios = require('axios');

async function executeTSD(source) {
    const url = 'http://localhost:8080/api/v1/execute';
    const payload = {
        source: source,
        source_name: 'js-client',
        verbose: false
    };
    
    const response = await axios.post(url, payload);
    return response.data;
}

// Utilisation
(async () => {
    const result = await executeTSD(`
        type Person : <id: string, name: string>
        Person("p1", "Alice")
    `);
    
    console.log(`Success: ${result.success}`);
    if (result.success) {
        console.log(`Activations: ${result.results.activations_count}`);
    }
})();
```

### cURL

```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{
    "source": "type Person : <id: string, name: string>\nPerson(\"p1\", \"Alice\")",
    "source_name": "curl-test",
    "verbose": false
  }' | jq
```

## Sécurité

### Recommandations

1. **Authentification**: Ajouter une couche d'authentification (JWT, API Key)
2. **Rate Limiting**: Limiter le nombre de requêtes par client
3. **Timeout**: Configurer des timeouts appropriés
4. **Validation**: Valider et sanitizer les entrées
5. **HTTPS**: Utiliser HTTPS en production
6. **Firewall**: Restreindre l'accès au serveur

### Exemple avec API Key (à implémenter)

```bash
# Client avec API key
./tsd-client -server http://localhost:8080 \
  -header "X-API-Key: your-secret-key" \
  program.tsd
```

## Performance

### Benchmarks

Sur une machine Intel i7 avec 16GB RAM:

- Parsing simple: ~5ms
- Exécution avec 10 faits: ~15ms
- Exécution avec 100 faits: ~50ms
- Exécution avec 1000 faits: ~200ms

### Optimisations

1. **Connection pooling**: Réutiliser les connexions HTTP
2. **Caching**: Cacher les programmes TSD fréquemment utilisés
3. **Batch processing**: Grouper plusieurs requêtes
4. **Load balancing**: Déployer plusieurs instances du serveur

## Dépannage

### Le serveur ne démarre pas

```bash
# Vérifier que le port n'est pas déjà utilisé
lsof -i :8080

# Changer de port
./tsd-server -port 9000
```

### Le client ne peut pas se connecter

```bash
# Vérifier que le serveur est accessible
curl http://localhost:8080/health

# Vérifier les logs du serveur
./tsd-server -v
```

### Erreurs de parsing

```bash
# Utiliser le mode verbeux pour voir les détails
./tsd-client -v program.tsd

# Valider le programme localement d'abord
./tsd program.tsd
```

## Développement

### Tests

```bash
# Tests du serveur
cd cmd/tsd-server
go test -v

# Tests du client
cd cmd/tsd-client
go test -v

# Tests d'intégration
go test -v ./...
```

### Build

```bash
# Build local
make build

# Build avec optimisations
go build -ldflags="-s -w" -o tsd-server ./cmd/tsd-server
go build -ldflags="-s -w" -o tsd-client ./cmd/tsd-client

# Cross-compilation
GOOS=linux GOARCH=amd64 go build -o tsd-server-linux ./cmd/tsd-server
GOOS=windows GOARCH=amd64 go build -o tsd-server.exe ./cmd/tsd-server
```

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License