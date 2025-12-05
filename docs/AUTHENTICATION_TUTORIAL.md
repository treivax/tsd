# Tutoriel Complet : Authentification TSD

Ce tutoriel vous guide à travers la configuration et l'utilisation de l'authentification pour le serveur TSD, avec deux méthodes : **Auth Key** (clés API) et **JWT** (JSON Web Tokens).

## Table des matières

1. [Introduction](#introduction)
2. [Installation et Configuration](#installation-et-configuration)
3. [Authentification par Clé API (Auth Key)](#authentification-par-clé-api-auth-key)
   - [Configuration du serveur](#31-configuration-du-serveur-avec-auth-key)
   - [Utilisation avec le CLI](#32-utilisation-avec-le-cli)
   - [Utilisation avec Python](#33-utilisation-avec-python)
4. [Authentification par JWT](#authentification-par-jwt)
   - [Configuration du serveur](#41-configuration-du-serveur-avec-jwt)
   - [Utilisation avec le CLI](#42-utilisation-avec-le-cli-1)
   - [Utilisation avec Python](#43-utilisation-avec-python-1)
5. [Exemples de Sessions Complètes](#5-exemples-de-sessions-complètes)
6. [Bonnes Pratiques de Sécurité](#6-bonnes-pratiques-de-sécurité)
7. [Dépannage](#7-dépannage)

---

## 1. Introduction

Le serveur TSD supporte trois modes d'authentification :

- **`none`** : Aucune authentification (mode développement)
- **`key`** : Authentification par clé API statique
- **`jwt`** : Authentification par JSON Web Token avec expiration

### Quelle méthode choisir ?

| Méthode | Cas d'usage | Avantages | Inconvénients |
|---------|-------------|-----------|---------------|
| **Auth Key** | Scripts, CI/CD, intégrations simples | Simple, pas d'expiration | Pas de métadonnées utilisateur, révocation manuelle |
| **JWT** | Applications multi-utilisateurs, API publiques | Expiration automatique, métadonnées, révocation | Plus complexe à mettre en place |
| **None** | Développement local uniquement | Aucune configuration | ⚠️ DANGEREUX en production |

---

## 2. Installation et Configuration

### 2.1. Compilation des outils

Compilez le binaire unique TSD :

```bash
# Se placer dans le répertoire du projet
cd tsd

# Compiler le binaire unique TSD
go build -o bin/tsd ./cmd/tsd

# Ajouter au PATH (optionnel)
export PATH=$PATH:$(pwd)/bin
```

Le binaire unique `tsd` gère tous les rôles (serveur, client, authentification) via des sous-commandes.

### 2.2. Vérification de l'installation

```bash
# Vérifier la version et l'aide
tsd -h
tsd version

# Voir l'aide pour chaque sous-commande
tsd server -h
tsd client -h
tsd auth -h
```

### 2.3. Configuration TLS/HTTPS (Requis)

⚠️ **Important** : TSD utilise HTTPS par défaut pour sécuriser les communications. Vous devez générer des certificats avant de démarrer le serveur.

#### Étape 1 : Générer des certificats pour le développement

```bash
# Générer des certificats auto-signés
tsd auth generate-cert

# Sortie :
# 🔐 Certificats TLS générés avec succès!
# =====================================
# 
# 📁 Répertoire: ./certs
# 
# 📄 Fichiers générés:
#    - ./certs/server.crt (certificat serveur)
#    - ./certs/server.key (clé privée serveur)
#    - ./certs/ca.crt (certificat CA pour clients)
```

⚠️ **Sécurité** : 
- Ne JAMAIS committer les certificats dans Git
- Ces certificats sont auto-signés (développement uniquement)
- En production, utilisez des certificats signés (Let's Encrypt, etc.)

#### Étape 2 : Options avancées de génération

```bash
# Personnaliser les hôtes autorisés
tsd auth generate-cert -hosts "localhost,127.0.0.1,192.168.1.100,myserver.local"

# Personnaliser la durée de validité
tsd auth generate-cert -valid-days 730

# Répertoire de sortie personnalisé
tsd auth generate-cert -output-dir ./my-certs

# Tout personnalisé
tsd auth generate-cert \
  -hosts "localhost,127.0.0.1" \
  -valid-days 365 \
  -org "My Company" \
  -output-dir ./certs
```

#### Étape 3 : Vérifier les certificats générés

```bash
# Lister les fichiers
ls -lh certs/

# Afficher les détails du certificat
openssl x509 -in certs/server.crt -text -noout
```

#### Configuration pour la Production

Pour un environnement de production, utilisez des certificats signés par une autorité de certification reconnue :

```bash
# Exemple avec Let's Encrypt (certbot)
sudo certbot certonly --standalone -d tsd.example.com

# Démarrer le serveur avec ces certificats
tsd server \
  --tls-cert /etc/letsencrypt/live/tsd.example.com/fullchain.pem \
  --tls-key /etc/letsencrypt/live/tsd.example.com/privkey.pem
```

#### Mode HTTP non sécurisé (Déconseillé)

Pour désactiver TLS en développement (NON recommandé) :

```bash
# Serveur en HTTP simple
tsd server --insecure

# Client vers serveur HTTP
tsd client program.tsd -server http://localhost:8080
```

⚠️ **Avertissement** : N'utilisez JAMAIS `--insecure` en production !

---

## 3. Authentification par Clé API (Auth Key)

L'authentification par clé API utilise des tokens statiques pré-partagés. C'est la méthode la plus simple pour sécuriser votre serveur TSD.

### 3.1. Configuration du serveur avec Auth Key

#### Étape 1 : Générer des certificats TLS (si pas déjà fait)

```bash
# Générer les certificats pour HTTPS
tsd auth generate-cert
```

#### Étape 2 : Générer des clés API

```bash
# 1. Générer les certificats TLS
$ tsd auth generate-cert
🔐 Certificats TLS générés avec succès!
=====================================

📁 Répertoire: ./certs

📄 Fichiers générés:
   - ./certs/server.crt (certificat serveur)
   - ./certs/server.key (clé privée serveur)
   - ./certs/ca.crt (certificat CA pour clients)

# 2. Générer une clé API
$ tsd auth generate-key
```

**Sortie :**
```
🔑 Clé(s) API générée(s):
========================

Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k
```

⚠️ **Important** : Sauvegardez cette clé dans un endroit sûr ! Elle ne peut pas être récupérée.

```bash
# Générer 3 clés (pour différents clients/services)
tsd auth generate-key -count 3

# Format JSON pour intégration
tsd auth generate-key -count 3 -format json > keys.json
```

#### Étape 4 : Démarrer le serveur avec JWT (HTTPS)

**Option A : Via ligne de commande**

```bash
tsd server \
  -auth key \
  -auth-keys "Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"
```

**Option B : Via variables d'environnement (recommandé)**

```bash
# Définir la clé API
export TSD_AUTH_KEYS="Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"

# Démarrer le serveur
tsd server -auth jwt
```

**Option C : Plusieurs clés (séparées par des virgules)**

```bash
export TSD_AUTH_KEYS="key1_xxxxxxx,key2_yyyyyyy,key3_zzzzzzz"
tsd server -auth key
```

**Sortie du serveur :**
```
[TSD-SERVER] 2025/01/15 10:00:00 🚀 Démarrage du serveur TSD sur https://0.0.0.0:8080
[TSD-SERVER] 2025/01/15 10:00:00 📊 Version: 1.0.0
[TSD-SERVER] 2025/01/15 10:00:00 🔒 TLS: activé
[TSD-SERVER] 2025/01/15 10:00:00    Certificat: ./certs/server.crt
[TSD-SERVER] 2025/01/15 10:00:00    Clé: ./certs/server.key
[TSD-SERVER] 2025/01/15 10:00:00 🔒 Authentification: activée (key)
[TSD-SERVER] 2025/01/15 10:00:00 🔗 Endpoints disponibles:
[TSD-SERVER] 2025/01/15 10:00:00    POST https://0.0.0.0:8080/api/v1/execute - Exécuter un programme TSD
[TSD-SERVER] 2025/01/15 10:00:00    GET  https://0.0.0.0:8080/health - Health check
[TSD-SERVER] 2025/01/15 10:00:00    GET  https://0.0.0.0:8080/api/v1/version - Version info
```

### 3.2. Utilisation avec le CLI

#### Test de connexion sans authentification

```bash
# Essayer sans token (devrait échouer)
tsd client health
```

**Sortie :**
```
❌ Erreur health check: requête health: ...
```

#### Utilisation avec le token

**Option A : Via flag**

```bash
tsd client health -token "Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"
```

**Option B : Via variable d'environnement (recommandé)**

```bash
# Définir le token
export TSD_AUTH_TOKEN="Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"

# Utiliser le client (le token sera automatiquement utilisé)
tsd client health
```

**Sortie :**
```
✅ Serveur TSD: ok
📊 Version: 1.0.0
⏱️  Uptime: 42s
🕐 Timestamp: 2025-01-15T10:31:27+01:00
```

#### Exécuter un programme TSD

```bash
# Créer un fichier TSD de test
cat > example.tsd << 'EOF'
type Person : <
  id: string,
  name: string,
  age: int
>

Person("p1", "Alice", 30)
Person("p2", "Bob", 25)
EOF

# Exécuter avec authentification
export TSD_AUTH_TOKEN="votre-cle-api"
tsd client execute example.tsd -v
```

**Sortie :**
```
📤 Envoi requête à http://localhost:8080/api/v1/execute...
🔒 Authentification: activée

✅ EXÉCUTION RÉUSSIE
===================
Temps d'exécution: 45ms
Faits injectés: 2
Activations: 0
```

### 3.3. Utilisation avec Python

#### Installation des dépendances

```bash
pip install requests
```

#### Exemple Python complet

```python
#!/usr/bin/env python3
"""
Exemple d'utilisation du serveur TSD avec authentification par clé API
"""

import os
import requests
import json

class TSDClient:
    """Client Python pour le serveur TSD avec Auth Key"""
    
    def __init__(self, server_url="http://localhost:8080", auth_token=None):
        self.server_url = server_url
        self.auth_token = auth_token or os.getenv("TSD_AUTH_TOKEN")
        
        if not self.auth_token:
            raise ValueError("Token d'authentification requis (auth_token ou TSD_AUTH_TOKEN)")
    
    def _get_headers(self):
        """Retourne les headers HTTP avec authentification"""
        return {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.auth_token}"
        }
    
    def health_check(self):
        """Vérifie la santé du serveur"""
        response = requests.get(
            f"{self.server_url}/health",
            headers=self._get_headers()
        )
        response.raise_for_status()
        return response.json()
    
    def execute(self, source, source_name="<python>", verbose=False):
        """Exécute un programme TSD"""
        payload = {
            "source": source,
            "source_name": source_name,
            "verbose": verbose
        }
        
        response = requests.post(
            f"{self.server_url}/api/v1/execute",
            headers=self._get_headers(),
            json=payload,
            timeout=30
        )
        response.raise_for_status()
        return response.json()

# Exemple d'utilisation
def main():
    # Définir votre clé API
    AUTH_KEY = "Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"
    
    # Ou utiliser la variable d'environnement
    # AUTH_KEY = None  # Utilisera TSD_AUTH_TOKEN
    
    # Créer le client
    client = TSDClient(
        server_url="http://localhost:8080",
        auth_token=AUTH_KEY
    )
    
    # Test de connexion
    print("🔍 Test de connexion...")
    health = client.health_check()
    print(f"✅ Serveur OK - Version: {health['version']}")
    print()
    
    # Exécuter un programme TSD
    print("📝 Exécution d'un programme TSD...")
    tsd_code = """
type Person : <
  id: string,
  name: string,
  age: int
>

Person("p1", "Alice", 30)
Person("p2", "Bob", 25)
Person("p3", "Charlie", 35)
"""
    
    result = client.execute(tsd_code, verbose=True)
    
    if result["success"]:
        print(f"✅ Succès!")
        print(f"   Faits: {result['results']['facts_count']}")
        print(f"   Activations: {result['results']['activations_count']}")
        print(f"   Temps: {result['execution_time_ms']}ms")
    else:
        print(f"❌ Erreur: {result['error']}")
        print(f"   Type: {result['error_type']}")

if __name__ == "__main__":
    main()
```

#### Exécution du script Python

```bash
# Définir la clé API
export TSD_AUTH_TOKEN="Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2uV5xY8zA3bD6eG9hJ2k"

# Exécuter le script
python3 tsd_client.py
```

**Sortie :**
```
🔍 Test de connexion...
✅ Serveur OK - Version: 1.0.0

📝 Exécution d'un programme TSD...
✅ Succès!
   Faits: 3
   Activations: 0
   Temps: 42ms
```

---

## 4. Authentification par JWT

L'authentification JWT est plus avancée et permet d'inclure des métadonnées utilisateur (nom, rôles) et une expiration automatique.

### 4.1. Configuration du serveur avec JWT

#### Étape 1 : Générer des certificats TLS (si pas déjà fait)

```bash
# Générer les certificats pour HTTPS
tsd auth generate-cert
```

#### Étape 2 : Générer un secret JWT

Le secret JWT doit être une chaîne aléatoire sécurisée d'au moins 32 caractères.

```bash
# Générer un secret aléatoire (Linux/macOS)
openssl rand -base64 32

# Ou utiliser tsd auth pour générer une clé sécurisée
tsd auth generate-key
```

**Exemple de secret :**
```
bW9uLXNlY3JldC1zdXBlci1zZWN1cmlzZS1kZS0zMi1jaGFyYWN0ZXJlcy1taW5pbXVt
```

#### Étape 2 : Démarrer le serveur avec JWT

```bash
# Définir le secret JWT (NE JAMAIS commiter dans git!)
export TSD_JWT_SECRET="bW9uLXNlY3JldC1zdXBlci1zZWN1cmlzZS1kZS0zMi1jaGFyYWN0ZXJlcy1taW5pbXVt"

# Démarrer le serveur en mode JWT
tsd server -auth jwt

# Avec expiration personnalisée (défaut: 24h)
tsd server -auth jwt -jwt-expiration 48h

# Avec émetteur personnalisé
tsd server -auth jwt -jwt-issuer "my-company-tsd"
```

**Sortie du serveur :**
**Sortie :**
```
[TSD-SERVER] 2025/01/15 10:30:00 🚀 Démarrage du serveur TSD sur https://0.0.0.0:8080
[TSD-SERVER] 2025/01/15 10:30:00 📊 Version: 1.0.0
[TSD-SERVER] 2025/01/15 10:30:00 🔒 TLS: activé
[TSD-SERVER] 2025/01/15 10:30:00    Certificat: ./certs/server.crt
[TSD-SERVER] 2025/01/15 10:30:00    Clé: ./certs/server.key
[TSD-SERVER] 2025/01/15 10:30:00 🔒 Authentification: activée (jwt)
```

### 4.2. Utilisation avec le CLI

#### Étape 1 : Générer un JWT

**Option A : Ligne de commande**

```bash
tsd auth generate-jwt \
  -secret "bW9uLXNlY3JldC1zdXBlci1zZWN1cmlzZS1kZS0zMi1jaGFyYWN0ZXJlcy1taW5pbXVt" \
  -username "alice" \
  -roles "admin,user" \
  -expiration 24h
```

**Option B : Mode interactif (recommandé pour ne pas exposer le secret)**

```bash
tsd auth generate-jwt -i -username alice
# Le système vous demandera le secret de manière sécurisée
```

**Sortie :**
```
🎫 JWT généré avec succès:
==========================
Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFsaWNlIiwicm9sZXMiOlsiYWRtaW4iLCJ1c2VyIl0sImV4cCI6MTcwNTQxMDAwMCwiaWF0IjoxNzA1MzIzNjAwLCJuYmYiOjE3MDUzMjM2MDAsImlzcyI6InRzZC1zZXJ2ZXIifQ.Xj8KpL9mN2qR5sT7vW0yZ3bC6dF8gH1jK4lM7nP0qS2

Utilisateur: alice
Rôles: admin, user
Expire dans: 24h0m0s
Expire le: 2025-01-16T11:00:00+01:00
Émetteur: tsd-server

⚠️  IMPORTANT: Conservez ce token en lieu sûr!
```

#### Étape 2 : Valider un JWT (optionnel)

```bash
# Valider le token
tsd auth validate \
  -type jwt \
  -token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -secret "bW9uLXNlY3JldC1zdXBlci1zZWN1cmlzZS1kZS0zMi1jaGFyYWN0ZXJlcy1taW5pbXVt"

# Mode interactif
tsd auth validate -i
```

**Sortie si valide :**
```
✅ Token valide
Type: jwt
Utilisateur: alice
Rôles: admin, user
```

#### Étape 3 : Utiliser le JWT avec le client

```bash
# Définir le token JWT
export TSD_AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Tester la connexion
tsd client health

# Exécuter un programme
tsd client execute example.tsd
```

### 4.3. Utilisation avec Python

#### Exemple Python avec JWT

```python
#!/usr/bin/env python3
"""
Exemple d'utilisation du serveur TSD avec authentification JWT
"""

import os
import requests
import json
from datetime import datetime, timedelta
import hmac
import hashlib
import base64

class TSDJWTClient:
    """Client Python pour le serveur TSD avec JWT"""
    
    def __init__(self, server_url="http://localhost:8080", jwt_token=None):
        self.server_url = server_url
        self.jwt_token = jwt_token or os.getenv("TSD_AUTH_TOKEN")
        
        if not self.jwt_token:
            raise ValueError("JWT token requis (jwt_token ou TSD_AUTH_TOKEN)")
    
    def _get_headers(self):
        """Retourne les headers HTTP avec JWT"""
        return {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.jwt_token}"
        }
    
    def health_check(self):
        """Vérifie la santé du serveur"""
        response = requests.get(
            f"{self.server_url}/health",
            headers=self._get_headers()
        )
        
        if response.status_code == 401:
            raise Exception("Authentification échouée - JWT invalide ou expiré")
        
        response.raise_for_status()
        return response.json()
    
    def execute(self, source, source_name="<python>", verbose=False):
        """Exécute un programme TSD"""
        payload = {
            "source": source,
            "source_name": source_name,
            "verbose": verbose
        }
        
        response = requests.post(
            f"{self.server_url}/api/v1/execute",
            headers=self._get_headers(),
            json=payload,
            timeout=30
        )
        
        if response.status_code == 401:
            raise Exception("Authentification échouée - JWT invalide ou expiré")
        
        response.raise_for_status()
        return response.json()

# Exemple d'utilisation avec génération de JWT en Python
def generate_jwt_python(secret, username, roles=None, expiration_hours=24):
    """
    Génère un JWT simple (pour démonstration)
    En production, utilisez une bibliothèque comme PyJWT
    """
    try:
        import jwt as pyjwt
        from datetime import datetime, timedelta
        
        payload = {
            "username": username,
            "roles": roles or [],
            "exp": datetime.utcnow() + timedelta(hours=expiration_hours),
            "iat": datetime.utcnow(),
            "nbf": datetime.utcnow(),
            "iss": "tsd-server"
        }
        
        token = pyjwt.encode(payload, secret, algorithm="HS256")
        return token
    except ImportError:
        print("⚠️  PyJWT non installé. Utilisez: pip install PyJWT")
        print("   Ou générez le JWT avec: tsd auth generate-jwt")
        return None

# Exemple d'utilisation
def main():
    # Option 1: Utiliser un JWT existant
    JWT_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    
    # Option 2: Générer un JWT en Python (nécessite PyJWT)
    # JWT_SECRET = "bW9uLXNlY3JldC1zdXBlci1zZWN1cmlzZS1kZS0zMi1jaGFyYWN0ZXJlcy1taW5pbXVt"
    # JWT_TOKEN = generate_jwt_python(JWT_SECRET, "alice", ["admin", "user"])
    
    # Option 3: Utiliser la variable d'environnement
    JWT_TOKEN = None  # Utilisera TSD_AUTH_TOKEN
    
    # Créer le client
    client = TSDJWTClient(
        server_url="http://localhost:8080",
        jwt_token=JWT_TOKEN
    )
    
    try:
        # Test de connexion
        print("🔍 Test de connexion avec JWT...")
        health = client.health_check()
        print(f"✅ Serveur OK - Version: {health['version']}")
        print()
        
        # Exécuter un programme TSD
        print("📝 Exécution d'un programme TSD...")
        tsd_code = """
type User : <
  id: string,
  username: string,
  role: string
>

User("u1", "alice", "admin")
User("u2", "bob", "user")
"""
        
        result = client.execute(tsd_code, verbose=True)
        
        if result["success"]:
            print(f"✅ Succès!")
            print(f"   Faits: {result['results']['facts_count']}")
            print(f"   Activations: {result['results']['activations_count']}")
            print(f"   Temps: {result['execution_time_ms']}ms")
        else:
            print(f"❌ Erreur: {result['error']}")
    
    except Exception as e:
        print(f"❌ Erreur: {e}")

if __name__ == "__main__":
    main()
```

#### Installation de PyJWT (optionnel)

```bash
# Pour générer des JWT en Python
pip install PyJWT
```

#### Exécution du script

```bash
# Générer un JWT
JWT_TOKEN=$(tsd auth generate-jwt \
  -secret "bW9uLXNlY3JldC1zdXBlci1zZWN1cmlzZS1kZS0zMi1jaGFyYWN0ZXJlcy1taW5pbXVt" \
  -username "alice" \
  -format json | jq -r .token)

# Définir le token
export TSD_AUTH_TOKEN="$JWT_TOKEN"

# Exécuter le script
python3 tsd_jwt_client.py
```

---

## 5. Exemples de Sessions Complètes

### 5.1. Session complète avec Auth Key

#### Terminal 1 : Serveur

```bash
# Générer une clé API
export MY_API_KEY=$(tsd auth generate-key -format json | jq -r '.keys[0]')
echo "Clé API: $MY_API_KEY"

# Démarrer le serveur
export TSD_AUTH_KEYS="$MY_API_KEY"
tsd server -auth key -v
```

#### Terminal 2 : Client CLI

```bash
# Définir la clé API (même que le serveur)
export TSD_AUTH_TOKEN="votre-cle-generee"

# Test de connexion
tsd client health

# Créer un programme TSD
cat > inventory.tsd << 'EOF'
type Product : <
  id: string,
  name: string,
  price: float,
  stock: int
>

type Order : <
  product_id: string,
  quantity: int
>

Product("p1", "Laptop", 999.99, 10)
Product("p2", "Mouse", 29.99, 50)
Product("p3", "Keyboard", 79.99, 25)

Order("p1", 2)
Order("p2", 5)
EOF

# Exécuter le programme
tsd client execute inventory.tsd -v

# Format JSON pour intégration
tsd client execute inventory.tsd -format json > result.json
cat result.json | jq .
```

#### Terminal 3 : Client Python

```bash
# Créer le script Python
cat > inventory_client.py << 'PYTHON'
import os
import requests

TSD_CODE = """
type Product : <
  id: string,
  name: string,
  price: float,
  stock: int
>

Product("p1", "Laptop", 999.99, 10)
Product("p2", "Mouse", 29.99, 50)
"""

headers = {
    "Authorization": f"Bearer {os.getenv('TSD_AUTH_TOKEN')}",
    "Content-Type": "application/json"
}

response = requests.post(
    "http://localhost:8080/api/v1/execute",
    headers=headers,
    json={"source": TSD_CODE, "source_name": "inventory.tsd"}
)

print(response.json())
PYTHON

# Exécuter
python3 inventory_client.py
```

### 5.2. Session complète avec JWT

#### Terminal 1 : Serveur

```bash
# Générer un secret JWT sécurisé
export TSD_JWT_SECRET=$(openssl rand -base64 32)
echo "Secret JWT: $TSD_JWT_SECRET"

# Démarrer le serveur avec JWT
tsd server -auth jwt -jwt-expiration 1h -v
```

#### Terminal 2 : Génération et utilisation de tokens

```bash
# Sauvegarder le même secret que le serveur
export TSD_JWT_SECRET="le-meme-secret-que-le-serveur"

# Générer un JWT pour Alice (admin)
ALICE_TOKEN=$(tsd auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "alice" \
  -roles "admin,developer" \
  -expiration 1h \
  -format json | jq -r .token)

echo "Token Alice: $ALICE_TOKEN"

# Générer un JWT pour Bob (utilisateur standard)
BOB_TOKEN=$(tsd auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "bob" \
  -roles "user" \
  -expiration 30m \
  -format json | jq -r .token)

echo "Token Bob: $BOB_TOKEN"

# Utiliser le token d'Alice
export TSD_AUTH_TOKEN="$ALICE_TOKEN"
tsd client health
tsd client execute example.tsd -v

# Changer pour le token de Bob
export TSD_AUTH_TOKEN="$BOB_TOKEN"
tsd client health

# Attendre l'expiration (30min pour Bob)
sleep 1800
tsd client health  # Échouera avec "token expiré"
```

### 5.3. Session multi-utilisateurs avec JWT

```bash
# Serveur avec JWT
export TSD_JWT_SECRET="my-super-secret-jwt-key-32-chars-min"
tsd server -auth jwt -v

# Terminal Alice (admin)
export TSD_AUTH_TOKEN=$(tsd auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "alice" \
  -roles "admin" \
  -format json | jq -r .token)

tsd client execute example.tsd -v

# Terminal Bob (developer)
export TSD_AUTH_TOKEN=$(tsd auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "bob" \
  -roles "developer" \
  -format json | jq -r .token)

tsd client execute example.tsd -v

# Terminal Charlie (lecture seule)
export TSD_AUTH_TOKEN=$(tsd auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "charlie" \
  -roles "readonly" \
  -format json | jq -r .token)

tsd client execute example.tsd -v
```

---

## 6. Bonnes Pratiques de Sécurité

### 6.1. Gestion des secrets

#### ❌ À NE JAMAIS FAIRE

```bash
# NE JAMAIS commiter dans git
git add .env
echo "TSD_JWT_SECRET=my-secret" >> config.yaml

# NE JAMAIS logger les secrets
echo "Secret: $TSD_JWT_SECRET"

# NE JAMAIS passer en ligne de commande visible
ps aux | grep tsd  # Les arguments sont visibles!
```

#### ✅ Bonnes pratiques

```bash
# Utiliser des variables d'environnement
export TSD_JWT_SECRET=$(cat /secure/path/jwt-secret.key)

# Utiliser des fichiers de configuration sécurisés
chmod 600 /etc/tsd/config.env
source /etc/tsd/config.env

# Utiliser un gestionnaire de secrets (production)
export TSD_JWT_SECRET=$(vault kv get -field=jwt_secret secret/tsd)
export TSD_JWT_SECRET=$(aws secretsmanager get-secret-value --secret-id tsd-jwt --query SecretString --output text)

# Utiliser des fichiers .env avec gitignore
echo "TSD_JWT_SECRET=..." > .env
echo ".env" >> .gitignore
source .env
```

### 6.2. Rotation des clés

#### Auth Keys

```bash
# Générer de nouvelles clés
NEW_KEY=$(tsd auth generate-key -format json | jq -r '.keys[0]')

# Période de transition : autoriser les deux clés
export TSD_AUTH_KEYS="$OLD_KEY,$NEW_KEY"
tsd server -auth key

# Après migration des clients, retirer l'ancienne clé
export TSD_AUTH_KEYS="$NEW_KEY"
```

#### JWT Secret

```bash
# Impossible de faire une rotation transparente avec JWT
# Il faut :
# 1. Générer un nouveau secret
# 2. Redémarrer le serveur
# 3. Demander aux utilisateurs de se reconnecter

# Alternative : utiliser un délai de grâce court
tsd server -auth jwt -jwt-expiration 1h
```

### 6.3. Sécurité réseau

```bash
# Toujours utiliser HTTPS en production
# Utiliser un reverse proxy (nginx, caddy, traefik)

# nginx.conf
server {
    listen 443 ssl;
    server_name tsd.example.com;
    
    ssl_certificate /etc/ssl/certs/tsd.crt;
    ssl_certificate_key /etc/ssl/private/tsd.key;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# Client avec HTTPS
tsd client execute -server https://tsd.example.com example.tsd
```

### 6.4. Limitation du taux (Rate Limiting)

```bash
# À implémenter au niveau du reverse proxy

# nginx rate limiting
limit_req_zone $binary_remote_addr zone=tsd:10m rate=10r/s;

server {
    location /api/v1/execute {
        limit_req zone=tsd burst=20 nodelay;
        proxy_pass http://localhost:8080;
    }
}
```

### 6.5. Logs et audit

```bash
# Activer les logs verbeux
tsd server -auth jwt -v

# Rediriger vers un fichier
tsd server -auth jwt -v 2>&1 | tee /var/log/tsd/server.log

# Analyser les tentatives d'authentification échouées
grep "Authentification échouée" /var/log/tsd/server.log
```

---

## 7. Dépannage

### Problème : "token invalide"

```bash
# Vérifier que le token est correct
tsd auth validate -type key -token "votre-token" -keys "cle-attendue"

# Pour JWT, vérifier le secret
tsd auth validate -type jwt -token "votre-jwt" -secret "votre-secret"

# Vérifier les espaces/newlines
echo -n "votre-token" | wc -c  # Pas de caractères parasites
```

### Problème : "token expiré" (JWT)

```bash
# Vérifier l'expiration du JWT
# Décoder le JWT (attention: ne valide pas la signature!)
echo "eyJhbGciOi..." | cut -d. -f2 | base64 -d | jq .

# Régénérer un nouveau JWT
tsd auth generate-jwt -secret "$TSD_JWT_SECRET" -username alice
```

### Problème : Serveur ne démarre pas avec authentification

```bash
# Vérifier la configuration
tsd server -auth key  # Erreur: "au moins une clé API doit être configurée"

# Solution
export TSD_AUTH_KEYS="une-cle-valide"
tsd server -auth key

# Pour JWT
export TSD_JWT_SECRET="secret-au-moins-32-caracteres-long"
tsd server -auth jwt
```

### Problème : Client ne peut pas se connecter

```bash
# Test basique sans authentification
curl http://localhost:8080/health

# Test avec Auth Key
curl -H "Authorization: Bearer votre-cle-api" http://localhost:8080/health

# Test avec JWT
curl -H "Authorization: Bearer votre-jwt" http://localhost:8080/health

# Vérifier que le serveur est en mode authentifié
# Le serveur doit afficher: "🔒 Authentification: activée (key|jwt)"
```

### Problème : "méthode de signature inattendue" (JWT)

```bash
# Le JWT doit utiliser HS256
# Vérifier l'algorithme du JWT
echo "votre-jwt" | cut -d. -f1 | base64 -d | jq .

# Devrait afficher: {"alg":"HS256","typ":"JWT"}
```

### Debug complet

```bash
# Activer tous les logs
export TSD_DEBUG=1
tsd server -auth key -v

# Utiliser curl pour tester
curl -v \
  -H "Authorization: Bearer votre-token" \
  -H "Content-Type: application/json" \
  -d '{"source":"type Test : <id: string>\nTest(\"t1\")","source_name":"test"}' \
  http://localhost:8080/api/v1/execute

# Vérifier les variables d'environnement
env | grep TSD

# Tester avec un token minimal
export TSD_AUTH_TOKEN="test"
tsd client health  # Devrait échouer avec "token invalide"
```

---

## Ressources supplémentaires

- [Documentation du serveur TSD](./SERVER_USAGE.md)
- [Documentation du client TSD](./CLIENT_USAGE.md)
- [API Reference](./API_REFERENCE.md)
- [JWT.io](https://jwt.io) - Pour décoder et débugger des JWT
- [RFC 7519](https://tools.ietf.org/html/rfc7519) - Spécification JWT

---

## Support

Pour toute question ou problème :

1. Vérifiez les logs du serveur avec `-v`
2. Testez avec `curl` pour isoler le problème
3. Utilisez `tsd auth validate` pour vérifier vos tokens
4. Consultez la section [Dépannage](#7-dépannage)

---

**© 2025 TSD Contributors - Licensed under MIT License**