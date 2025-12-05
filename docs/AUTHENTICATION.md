# Authentification TSD

Ce document présente le système d'authentification pour le serveur TSD.

## 🔒 Vue d'ensemble

Le serveur TSD supporte trois modes d'authentification :

- **`none`** : Aucune authentification (mode développement uniquement)
- **`key`** : Authentification par clé API statique (simple et efficace)
- **`jwt`** : Authentification par JSON Web Token (avancé avec expiration)

## 📚 Documentation

### Guides Complets

- **[Tutoriel Complet](./AUTHENTICATION_TUTORIAL.md)** - Guide détaillé avec tous les cas d'usage
  - Installation et configuration
  - Authentification par clé API (Auth Key)
  - Authentification JWT
  - Exemples de sessions complètes
  - Bonnes pratiques de sécurité
  - Dépannage

- **[Guide de Démarrage Rapide](./AUTHENTICATION_QUICKSTART.md)** - Mise en place en 5 minutes
  - Démarrage rapide Auth Key
  - Démarrage rapide JWT
  - Exemples concrets
  - Commandes essentielles

### Exemples de Code

- **[Exemples Python](../examples/auth/)** - Code prêt à l'emploi
  - `client_auth_key.py` - Client Python avec Auth Key
  - `client_jwt.py` - Client Python avec JWT
  - README avec cas d'usage pratiques

## 🚀 Démarrage Rapide

### Option 1 : Auth Key (Recommandé pour débuter)

```bash
# 1. Compiler les outils
go build -o bin/tsd-auth ./cmd/tsd-auth
go build -o bin/tsd-server ./cmd/tsd-server
go build -o bin/tsd-client ./cmd/tsd-client

# 2. Générer une clé API
API_KEY=$(bin/tsd-auth generate-key -format json | jq -r '.keys[0]')

# 3. Démarrer le serveur
export TSD_AUTH_KEYS="$API_KEY"
bin/tsd-server -auth key

# 4. Utiliser le client (nouveau terminal)
export TSD_AUTH_TOKEN="$API_KEY"
bin/tsd-client -health
bin/tsd-client example.tsd
```

### Option 2 : JWT (Pour multi-utilisateurs)

```bash
# 1. Générer un secret JWT
JWT_SECRET=$(openssl rand -base64 32)

# 2. Démarrer le serveur
export TSD_JWT_SECRET="$JWT_SECRET"
bin/tsd-server -auth jwt

# 3. Générer un JWT (nouveau terminal)
TOKEN=$(bin/tsd-auth generate-jwt \
  -secret "$JWT_SECRET" \
  -username "alice" \
  -format json | jq -r .token)

# 4. Utiliser le client
export TSD_AUTH_TOKEN="$TOKEN"
bin/tsd-client -health
bin/tsd-client example.tsd
```

## 🔧 Outils

### tsd-auth - Outil de gestion d'authentification

```bash
# Générer une clé API
tsd-auth generate-key

# Générer plusieurs clés
tsd-auth generate-key -count 5

# Générer un JWT
tsd-auth generate-jwt \
  -secret "mon-secret" \
  -username "alice" \
  -roles "admin,user" \
  -expiration 24h

# Mode interactif (secret masqué)
tsd-auth generate-jwt -i -username alice

# Valider un token
tsd-auth validate -type key -token "..." -keys "..."
tsd-auth validate -type jwt -token "..." -secret "..."
```

### tsd-server - Options d'authentification

```bash
# Sans authentification (développement)
tsd-server

# Avec Auth Key
tsd-server -auth key -auth-keys "cle1,cle2,cle3"
# Ou via variable d'environnement
export TSD_AUTH_KEYS="cle1,cle2,cle3"
tsd-server -auth key

# Avec JWT
tsd-server -auth jwt -jwt-secret "mon-secret-32-chars-min"
# Ou via variable d'environnement
export TSD_JWT_SECRET="mon-secret"
tsd-server -auth jwt

# Options JWT avancées
tsd-server -auth jwt \
  -jwt-expiration 48h \
  -jwt-issuer "my-company"
```

### tsd-client - Utilisation avec authentification

```bash
# Via flag
tsd-client -token "votre-token" program.tsd

# Via variable d'environnement (recommandé)
export TSD_AUTH_TOKEN="votre-token"
tsd-client program.tsd

# Health check
tsd-client -health

# Format JSON
tsd-client program.tsd -format json
```

## 🐍 Utilisation avec Python

### Auth Key

```python
import requests
import os

API_KEY = os.getenv("TSD_AUTH_TOKEN")

headers = {
    "Authorization": f"Bearer {API_KEY}",
    "Content-Type": "application/json"
}

response = requests.post(
    "http://localhost:8080/api/v1/execute",
    headers=headers,
    json={"source": "type Test : <id: string>\nTest(\"t1\")"}
)

print(response.json())
```

### JWT

```python
import requests
import os

JWT_TOKEN = os.getenv("TSD_AUTH_TOKEN")

headers = {
    "Authorization": f"Bearer {JWT_TOKEN}",
    "Content-Type": "application/json"
}

response = requests.post(
    "http://localhost:8080/api/v1/execute",
    headers=headers,
    json={"source": "type Test : <id: string>\nTest(\"t1\")"}
)

result = response.json()
if result["success"]:
    print(f"✅ {result['results']['facts_count']} faits créés")
else:
    print(f"❌ {result['error']}")
```

Voir **[exemples complets](../examples/auth/)** pour plus de détails.

## 🔐 Quelle méthode choisir ?

| Critère | Auth Key | JWT |
|---------|----------|-----|
| **Simplicité** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| **Sécurité** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Expiration auto** | ❌ | ✅ |
| **Multi-utilisateurs** | ⚠️ Limité | ✅ Excellent |
| **Métadonnées** | ❌ | ✅ (username, rôles) |
| **Révocation** | ⚠️ Manuelle | ⚠️ Attendre expiration |
| **Cas d'usage** | Scripts, CI/CD | Applications, API |

### Recommandations

- **Développement local** : `none` (pas d'auth)
- **Scripts / CI/CD** : `key` (simple et efficace)
- **API publique** : `jwt` (expiration automatique)
- **Multi-tenant** : `jwt` (isolation par utilisateur)
- **Production simple** : `key` (moins de complexité)

## ⚠️ Sécurité

### ✅ À FAIRE

- Utiliser des variables d'environnement pour les secrets
- Utiliser HTTPS en production
- Générer des clés longues (32+ caractères)
- Rotation régulière des clés/secrets
- Limiter les permissions (principe du moindre privilège)
- Logger les tentatives d'authentification échouées

### ❌ À NE JAMAIS FAIRE

- Commiter les secrets dans git
- Hardcoder les tokens dans le code
- Utiliser HTTP en production
- Partager les mêmes tokens entre environnements
- Logger les tokens/secrets
- Utiliser des clés courtes ou prévisibles

## 📊 Variables d'Environnement

### Serveur

| Variable | Description | Exemple |
|----------|-------------|---------|
| `TSD_AUTH_KEYS` | Clés API (séparées par virgules) | `key1,key2,key3` |
| `TSD_JWT_SECRET` | Secret JWT (32+ chars) | `mon-secret-securise...` |

### Client

| Variable | Description | Exemple |
|----------|-------------|---------|
| `TSD_AUTH_TOKEN` | Token d'authentification | `eyJhbGciOi...` ou clé API |

## 🧪 Tests

```bash
# Exécuter les tests d'authentification
./scripts/test_auth.sh

# Tests unitaires
go test ./auth/...
go test ./cmd/tsd-server/...
go test ./cmd/tsd-client/...
```

## 🐛 Dépannage Rapide

### "token invalide"

```bash
# Vérifier le token
tsd-auth validate -type key -token "$TSD_AUTH_TOKEN" -keys "$TSD_AUTH_KEYS"

# Générer un nouveau token
API_KEY=$(tsd-auth generate-key -format json | jq -r '.keys[0]')
export TSD_AUTH_TOKEN="$API_KEY"
```

### "token expiré" (JWT)

```bash
# Générer un nouveau JWT
export TSD_AUTH_TOKEN=$(tsd-auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "alice" \
  -format json | jq -r .token)
```

### Connexion refusée

```bash
# Vérifier que le serveur est démarré
curl http://localhost:8080/health

# Avec authentification
curl -H "Authorization: Bearer $TSD_AUTH_TOKEN" \
  http://localhost:8080/health
```

## 📖 Ressources

- [Tutoriel complet](./AUTHENTICATION_TUTORIAL.md) - 1000+ lignes de documentation
- [Guide rapide](./AUTHENTICATION_QUICKSTART.md) - Démarrage en 5 minutes
- [Exemples Python](../examples/auth/) - Code prêt à l'emploi
- [Documentation du serveur](./SERVER_USAGE.md)
- [Documentation du client](./CLIENT_USAGE.md)
- [JWT.io](https://jwt.io) - Décoder et debugger des JWT
- [RFC 7519](https://tools.ietf.org/html/rfc7519) - Spécification JWT

## 💬 Support

Pour toute question :

1. Consultez le [tutoriel complet](./AUTHENTICATION_TUTORIAL.md)
2. Vérifiez la section [Dépannage](./AUTHENTICATION_TUTORIAL.md#7-dépannage)
3. Exécutez les tests : `./scripts/test_auth.sh`
4. Testez avec curl pour isoler le problème

---

**© 2025 TSD Contributors - Licensed under MIT License**