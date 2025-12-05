# Guide de Démarrage Rapide : Authentification TSD

Guide ultra-rapide pour sécuriser votre serveur TSD en 5 minutes.

## 🚀 Démarrage Rapide

### Option 1 : Auth Key (Le plus simple)

```bash
# 1. Compiler les outils
go build -o bin/tsd-auth ./cmd/tsd-auth
go build -o bin/tsd-server ./cmd/tsd-server
go build -o bin/tsd-client ./cmd/tsd-client

# 2. Générer une clé API
API_KEY=$(bin/tsd-auth generate-key -format json | jq -r '.keys[0]')
echo "Votre clé API: $API_KEY"

# 3. Démarrer le serveur
export TSD_AUTH_KEYS="$API_KEY"
bin/tsd-server -auth key

# 4. Utiliser le client (nouveau terminal)
export TSD_AUTH_TOKEN="$API_KEY"
bin/tsd-client -health
bin/tsd-client example.tsd
```

### Option 2 : JWT (Multi-utilisateurs)

```bash
# 1. Générer un secret JWT
JWT_SECRET=$(openssl rand -base64 32)
echo "Votre secret JWT: $JWT_SECRET"

# 2. Démarrer le serveur
export TSD_JWT_SECRET="$JWT_SECRET"
bin/tsd-server -auth jwt

# 3. Générer un JWT pour un utilisateur (nouveau terminal)
TOKEN=$(bin/tsd-auth generate-jwt \
  -secret "$JWT_SECRET" \
  -username "alice" \
  -format json | jq -r .token)

# 4. Utiliser le client
export TSD_AUTH_TOKEN="$TOKEN"
bin/tsd-client -health
bin/tsd-client example.tsd
```

## 📝 Exemples Concrets

### Créer un fichier TSD de test

```bash
cat > example.tsd << 'EOF'
type Person : <
  id: string,
  name: string,
  age: int
>

Person("p1", "Alice", 30)
Person("p2", "Bob", 25)
EOF
```

### Client CLI avec Auth Key

```bash
# Définir votre clé
export TSD_AUTH_TOKEN="Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"

# Exécuter
bin/tsd-client example.tsd -v
```

### Client Python avec Auth Key

```python
import requests
import os

# Configuration
API_KEY = os.getenv("TSD_AUTH_TOKEN")
SERVER_URL = "http://localhost:8080"

# En-têtes avec authentification
headers = {
    "Authorization": f"Bearer {API_KEY}",
    "Content-Type": "application/json"
}

# Exécuter du code TSD
response = requests.post(
    f"{SERVER_URL}/api/v1/execute",
    headers=headers,
    json={
        "source": 'type Person : <id: string, name: string>\nPerson("p1", "Alice")',
        "source_name": "test.tsd"
    }
)

print(response.json())
```

### Client Python avec JWT

```python
import requests
import os

# Configuration
JWT_TOKEN = os.getenv("TSD_AUTH_TOKEN")
SERVER_URL = "http://localhost:8080"

# En-têtes avec JWT
headers = {
    "Authorization": f"Bearer {JWT_TOKEN}",
    "Content-Type": "application/json"
}

# Health check
health = requests.get(f"{SERVER_URL}/health", headers=headers)
print(f"Serveur: {health.json()['status']}")

# Exécuter du code
response = requests.post(
    f"{SERVER_URL}/api/v1/execute",
    headers=headers,
    json={"source": "type Test : <id: string>\nTest(\"t1\")"}
)

result = response.json()
if result["success"]:
    print(f"✅ Succès: {result['results']['facts_count']} faits")
else:
    print(f"❌ Erreur: {result['error']}")
```

## 🔧 Commandes Utiles

### Génération de tokens

```bash
# Générer une clé API
bin/tsd-auth generate-key

# Générer plusieurs clés
bin/tsd-auth generate-key -count 5

# Générer un JWT
bin/tsd-auth generate-jwt \
  -secret "mon-secret-jwt" \
  -username "alice" \
  -roles "admin,user" \
  -expiration 24h

# Mode interactif (ne pas exposer le secret)
bin/tsd-auth generate-jwt -i -username alice
```

### Validation de tokens

```bash
# Valider une clé API
bin/tsd-auth validate \
  -type key \
  -token "votre-cle" \
  -keys "cle1,cle2,cle3"

# Valider un JWT
bin/tsd-auth validate \
  -type jwt \
  -token "eyJhbGciOi..." \
  -secret "votre-secret"

# Mode interactif
bin/tsd-auth validate -i
```

### Démarrage du serveur

```bash
# Sans authentification (développement)
bin/tsd-server

# Avec Auth Key
export TSD_AUTH_KEYS="key1,key2,key3"
bin/tsd-server -auth key

# Avec JWT
export TSD_JWT_SECRET="votre-secret-32-chars-minimum"
bin/tsd-server -auth jwt -jwt-expiration 24h

# Mode verbeux
bin/tsd-server -auth key -v
```

### Utilisation du client

```bash
# Health check
bin/tsd-client -health

# Exécuter un fichier
bin/tsd-client program.tsd

# Exécuter du code directement
bin/tsd-client -text 'type Test : <id: string>\nTest("t1")'

# Depuis stdin
cat program.tsd | bin/tsd-client -stdin

# Format JSON
bin/tsd-client program.tsd -format json

# Avec token
bin/tsd-client -token "votre-token" program.tsd

# Via variable d'environnement
export TSD_AUTH_TOKEN="votre-token"
bin/tsd-client program.tsd
```

## 🔒 Sécurité

### ✅ À FAIRE

```bash
# Utiliser des variables d'environnement
export TSD_AUTH_TOKEN="votre-token"

# Stocker les secrets dans des fichiers sécurisés
chmod 600 ~/.tsd/secrets
source ~/.tsd/secrets

# Utiliser HTTPS en production
bin/tsd-client -server https://tsd.example.com

# Rotation régulière des clés
```

### ❌ À NE JAMAIS FAIRE

```bash
# Ne jamais commiter dans git
git add .env  # NON !

# Ne jamais logger les secrets
echo "Token: $TSD_AUTH_TOKEN"  # NON !

# Ne jamais hardcoder
bin/tsd-server -auth-keys "ma-cle-secrete"  # NON ! (visible dans ps)

# Utiliser HTTP en production
bin/tsd-client -server http://api.prod.example.com  # NON !
```

## 🐛 Dépannage Rapide

### Erreur : "token invalide"

```bash
# Vérifier le token
bin/tsd-auth validate -type key -token "$TSD_AUTH_TOKEN" -keys "$TSD_AUTH_KEYS"

# Vérifier les espaces
echo -n "$TSD_AUTH_TOKEN" | wc -c
```

### Erreur : "token expiré" (JWT)

```bash
# Générer un nouveau token
export TSD_AUTH_TOKEN=$(bin/tsd-auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "alice" \
  -format json | jq -r .token)
```

### Erreur : Connexion refusée

```bash
# Vérifier que le serveur est démarré
curl http://localhost:8080/health

# Vérifier avec authentification
curl -H "Authorization: Bearer $TSD_AUTH_TOKEN" http://localhost:8080/health
```

### Debug complet

```bash
# Serveur en mode verbeux
bin/tsd-server -auth key -v

# Tester avec curl
curl -v \
  -H "Authorization: Bearer votre-token" \
  -H "Content-Type: application/json" \
  -d '{"source":"type Test : <id: string>\nTest(\"t1\")"}' \
  http://localhost:8080/api/v1/execute
```

## 📚 Documentation Complète

Pour plus de détails, consultez :

- [Tutoriel complet d'authentification](./AUTHENTICATION_TUTORIAL.md)
- [Documentation du serveur](./SERVER_USAGE.md)
- [Documentation du client](./CLIENT_USAGE.md)

## 🎯 Exemples par Cas d'Usage

### Script Bash

```bash
#!/bin/bash
export TSD_AUTH_TOKEN="your-api-key"

result=$(bin/tsd-client program.tsd -format json)
success=$(echo "$result" | jq -r .success)

if [ "$success" = "true" ]; then
    echo "✅ Succès"
else
    echo "❌ Erreur: $(echo "$result" | jq -r .error)"
    exit 1
fi
```

### CI/CD (GitHub Actions)

```yaml
- name: Execute TSD
  env:
    TSD_AUTH_TOKEN: ${{ secrets.TSD_API_KEY }}
  run: |
    ./bin/tsd-client -server https://tsd.prod.example.com program.tsd
```

### Docker

```dockerfile
FROM alpine:latest
COPY bin/tsd-server /usr/local/bin/
ENV TSD_AUTH_KEYS=""
EXPOSE 8080
CMD ["tsd-server", "-auth", "key"]
```

```bash
# Démarrer avec Docker
docker run -d \
  -e TSD_AUTH_KEYS="your-api-key" \
  -p 8080:8080 \
  tsd-server:latest
```

### Kubernetes

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tsd-auth
type: Opaque
stringData:
  auth-keys: "key1,key2,key3"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tsd-server
spec:
  template:
    spec:
      containers:
      - name: tsd-server
        image: tsd-server:latest
        args: ["-auth", "key"]
        env:
        - name: TSD_AUTH_KEYS
          valueFrom:
            secretKeyRef:
              name: tsd-auth
              key: auth-keys
```

## ⚡ Résumé Ultra-Rapide

```bash
# Auth Key (3 commandes)
API_KEY=$(bin/tsd-auth generate-key -format json | jq -r '.keys[0]')
TSD_AUTH_KEYS="$API_KEY" bin/tsd-server -auth key &
TSD_AUTH_TOKEN="$API_KEY" bin/tsd-client -health

# JWT (4 commandes)
JWT_SECRET=$(openssl rand -base64 32)
TSD_JWT_SECRET="$JWT_SECRET" bin/tsd-server -auth jwt &
TOKEN=$(bin/tsd-auth generate-jwt -secret "$JWT_SECRET" -username "alice" -format json | jq -r .token)
TSD_AUTH_TOKEN="$TOKEN" bin/tsd-client -health
```

C'est tout ! 🎉

---

**© 2025 TSD Contributors - Licensed under MIT License**