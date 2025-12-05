# Exemples pour le Serveur TSD

Ce répertoire contient des exemples de programmes TSD pour tester le serveur et le client.

## Fichiers d'exemple

### simple.tsd

Un exemple simple démontrant les fonctionnalités de base:
- Définition de types (Person, Order)
- Définition d'actions (notify, process_order, send_welcome)
- Règles avec conditions
- Assertions de faits

**Utilisation:**

```bash
# Avec le client TSD
tsd-client simple.tsd

# Avec curl
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d @- <<EOF
{
  "source": "$(cat simple.tsd | sed 's/"/\\"/g' | tr '\n' ' ')",
  "source_name": "simple.tsd"
}
EOF
```

**Résultat attendu:**
- 6 faits injectés (3 Person + 3 Order)
- 5 activations:
  - 1 action notify (Alice)
  - 2 actions process_order (commandes > 100)
  - 2 actions send_welcome (personnes < 30 ans)

### multiple_activations.tsd

Un exemple plus complexe avec plusieurs types et relations:
- Définition de types (Employee, Project, Assignment)
- Règles avec jointures entre plusieurs faits
- Multiples conditions et activations

**Utilisation:**

```bash
# Avec le client TSD
tsd-client multiple_activations.tsd

# En mode verbeux pour voir les détails
tsd-client -v multiple_activations.tsd

# Format JSON pour intégration
tsd-client -format json multiple_activations.tsd
```

**Résultat attendu:**
- 13 faits injectés (5 Employee + 4 Project + 4 Assignment)
- Plusieurs activations:
  - Actions alert_high_salary pour salaires > 80000
  - Actions approve_project pour budgets > 100000
  - Actions assign_task pour les assignments

## Tester les exemples

### Prérequis

1. **Démarrer le serveur:**

```bash
# Depuis la racine du projet
go build -o bin/tsd-server ./cmd/tsd-server
./bin/tsd-server
```

Le serveur démarre sur http://localhost:8080

2. **Compiler le client:**

```bash
go build -o bin/tsd-client ./cmd/tsd-client
```

### Exécution

#### Test rapide

```bash
# Test de santé du serveur
./bin/tsd-client -health

# Exécuter un exemple
./bin/tsd-client examples/server/simple.tsd
```

#### Mode verbeux

```bash
./bin/tsd-client -v examples/server/simple.tsd
```

Affiche:
- Les faits déclencheurs pour chaque activation
- Les détails des arguments
- Les attributs des faits

#### Format JSON

```bash
./bin/tsd-client -format json examples/server/simple.tsd
```

Retourne la réponse complète en JSON, utile pour l'intégration avec d'autres outils.

#### Avec stdin

```bash
cat examples/server/simple.tsd | ./bin/tsd-client -stdin
```

#### Serveur distant

```bash
./bin/tsd-client -server http://remote-server:8080 examples/server/simple.tsd
```

## Script de test automatique

Un script de test complet est disponible:

```bash
./scripts/test_server_client.sh
```

Ce script:
1. Compile le serveur et le client
2. Démarre le serveur
3. Exécute tous les tests (health check, fichiers, stdin, JSON, etc.)
4. Vérifie les résultats
5. Arrête le serveur

## Créer vos propres exemples

### Template de base

```tsd
// Définir vos types
type MyType : <field1: string, field2: number>

// Définir vos actions
action my_action : <param1: string>

// Définir vos règles
rule my_rule : {t: MyType} / t.field2 > 10 ==> my_action(t.field1)

// Asserter vos faits
MyType("id1", 15)
MyType("id2", 5)
```

### Conseils

1. **Types**: Définissez d'abord tous vos types
2. **Actions**: Déclarez les actions que vous allez utiliser
3. **Règles**: Écrivez les règles avec des conditions claires
4. **Faits**: Ajoutez les faits qui déclencheront les règles

### Validation locale

Avant de soumettre au serveur, validez localement:

```bash
./bin/tsd votre_programme.tsd
```

## Intégration programmatique

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "os"
)

func main() {
    // Lire le fichier
    source, _ := os.ReadFile("examples/server/simple.tsd")
    
    // Créer la requête
    req := map[string]interface{}{
        "source": string(source),
        "source_name": "simple.tsd",
    }
    
    jsonData, _ := json.Marshal(req)
    
    // Envoyer au serveur
    resp, _ := http.Post(
        "http://localhost:8080/api/v1/execute",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    defer resp.Body.Close()
    
    // Lire la réponse
    body, _ := io.ReadAll(resp.Body)
    println(string(body))
}
```

### Python

```python
import requests

# Lire le fichier
with open('examples/server/simple.tsd', 'r') as f:
    source = f.read()

# Créer la requête
payload = {
    'source': source,
    'source_name': 'simple.tsd'
}

# Envoyer au serveur
response = requests.post(
    'http://localhost:8080/api/v1/execute',
    json=payload
)

# Afficher la réponse
print(response.json())
```

### JavaScript

```javascript
const fs = require('fs');
const axios = require('axios');

// Lire le fichier
const source = fs.readFileSync('examples/server/simple.tsd', 'utf8');

// Créer la requête
const payload = {
    source: source,
    source_name: 'simple.tsd'
};

// Envoyer au serveur
axios.post('http://localhost:8080/api/v1/execute', payload)
    .then(response => {
        console.log(JSON.stringify(response.data, null, 2));
    })
    .catch(error => {
        console.error('Error:', error.message);
    });
```

### cURL

```bash
# Avec un fichier
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d "{\"source\": \"$(cat examples/server/simple.tsd | jq -sR .)\"}"

# Avec code inline
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{
    "source": "type Person : <id: string>\nPerson(\"p1\")",
    "source_name": "test"
  }'
```

## Dépannage

### Le serveur ne répond pas

```bash
# Vérifier que le serveur est démarré
curl http://localhost:8080/health

# Vérifier les logs
tail -f /tmp/tsd-server.log
```

### Erreurs de parsing

```bash
# Valider le programme localement
./bin/tsd examples/server/simple.tsd

# Utiliser le mode verbeux
./bin/tsd-client -v examples/server/simple.tsd
```

### Timeout

```bash
# Augmenter le timeout
./bin/tsd-client -timeout 60s examples/server/simple.tsd
```

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License