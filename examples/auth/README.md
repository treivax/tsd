# Exemples d'Authentification TSD

Ce répertoire contient des exemples d'utilisation du serveur TSD avec authentification.

## 📋 Contenu

- `client_auth_key.py` - Exemple complet avec authentification par clé API
- `client_jwt.py` - Exemple complet avec authentification JWT

## 🚀 Démarrage Rapide

### Prérequis

```bash
# Installer les dépendances Python
pip install requests

# Pour JWT (optionnel)
pip install PyJWT
```

### 1. Compiler les outils TSD

```bash
cd ../..  # Retour à la racine du projet
go build -o bin/tsd-server ./cmd/tsd-server
go build -o bin/tsd-client ./cmd/tsd-client
go build -o bin/tsd-auth ./cmd/tsd-auth
```

### 2. Authentification par Clé API

#### Terminal 1 : Démarrer le serveur

```bash
# Générer une clé API
export TSD_AUTH_KEYS=$(bin/tsd-auth generate-key -format json | jq -r '.keys[0]')

# Démarrer le serveur
bin/tsd-server -auth key -v
```

#### Terminal 2 : Exécuter l'exemple Python

```bash
# Utiliser la même clé que le serveur
export TSD_AUTH_TOKEN="votre-cle-api"

# Exécuter tous les exemples
python3 examples/auth/client_auth_key.py

# Ou un exemple spécifique
python3 examples/auth/client_auth_key.py --example 1
```

### 3. Authentification JWT

#### Terminal 1 : Démarrer le serveur

```bash
# Générer un secret JWT
export TSD_JWT_SECRET=$(openssl rand -base64 32)

# Démarrer le serveur
bin/tsd-server -auth jwt -v
```

#### Terminal 2 : Générer un JWT et exécuter l'exemple

```bash
# Générer un JWT
export TSD_JWT_SECRET="meme-secret-que-le-serveur"
export TSD_AUTH_TOKEN=$(bin/tsd-auth generate-jwt \
  -secret "$TSD_JWT_SECRET" \
  -username "alice" \
  -roles "admin,developer" \
  -format json | jq -r .token)

# Exécuter tous les exemples
python3 examples/auth/client_jwt.py

# Ou générer le JWT en Python
python3 examples/auth/client_jwt.py --generate --username alice --roles "admin,user"
```

## 📚 Exemples Détaillés

### Exemple Auth Key

```python
from client_auth_key import TSDAuthKeyClient

# Créer le client
client = TSDAuthKeyClient(
    server_url="http://localhost:8080",
    auth_token="votre-cle-api"
)

# Vérifier la connexion
health = client.health_check()
print(f"Serveur: {health['status']}")

# Exécuter du code TSD
result = client.execute("""
type Person : <
  id: string,
  name: string
>

Person("p1", "Alice")
""")

print(f"Succès: {result['success']}")
print(f"Faits: {result['results']['facts_count']}")
```

### Exemple JWT

```python
from client_jwt import TSDJWTClient, generate_jwt

# Option 1: Utiliser un JWT existant
client = TSDJWTClient(
    server_url="http://localhost:8080",
    jwt_token="eyJhbGciOi..."
)

# Option 2: Générer un JWT en Python
token = generate_jwt(
    secret="votre-secret-jwt",
    username="alice",
    roles=["admin", "user"],
    expiration_hours=24
)
client = TSDJWTClient(jwt_token=token)

# Utiliser le client
health = client.health_check()
result = client.execute("type Test : <id: string>\nTest(\"t1\")")
```

## 🎯 Cas d'Usage

### Script d'automatisation avec Auth Key

```python
#!/usr/bin/env python3
import os
from client_auth_key import TSDAuthKeyClient

# Configuration via variables d'environnement
client = TSDAuthKeyClient()

# Exécuter plusieurs programmes
programs = [
    "program1.tsd",
    "program2.tsd",
    "program3.tsd"
]

for program in programs:
    try:
        result = client.execute_file(program)
        if result["success"]:
            print(f"✅ {program}: OK")
        else:
            print(f"❌ {program}: {result['error']}")
    except Exception as e:
        print(f"❌ {program}: {e}")
```

### Application multi-utilisateurs avec JWT

```python
#!/usr/bin/env python3
from client_jwt import TSDJWTClient, generate_jwt

# Secret partagé avec le serveur
JWT_SECRET = os.getenv("TSD_JWT_SECRET")

def create_user_client(username, roles):
    """Crée un client pour un utilisateur"""
    token = generate_jwt(
        secret=JWT_SECRET,
        username=username,
        roles=roles,
        expiration_hours=1
    )
    return TSDJWTClient(jwt_token=token)

# Créer des clients pour différents utilisateurs
admin_client = create_user_client("admin", ["admin"])
dev_client = create_user_client("developer", ["developer"])
user_client = create_user_client("user", ["readonly"])

# Chaque client peut maintenant exécuter du code
# avec son propre contexte d'authentification
```

### Test de tokens expirés

```python
#!/usr/bin/env python3
import time
from client_jwt import TSDJWTClient, generate_jwt

# Générer un token avec expiration courte
token = generate_jwt(
    secret="mon-secret",
    username="test",
    expiration_hours=0.001  # ~3 secondes
)

client = TSDJWTClient(jwt_token=token)

# Test immédiat
try:
    client.health_check()
    print("✅ Token valide")
except:
    print("❌ Token invalide")

# Attendre l'expiration
time.sleep(5)

# Test après expiration
try:
    client.health_check()
    print("✅ Token valide")
except Exception as e:
    print(f"❌ Token expiré: {e}")
    
    # Régénérer un nouveau token
    new_token = generate_jwt(
        secret="mon-secret",
        username="test",
        expiration_hours=1
    )
    client.update_token(new_token)
    
    # Réessayer
    client.health_check()
    print("✅ Nouveau token valide")
```

## 🔧 Options des Scripts

### client_auth_key.py

```bash
# Afficher l'aide
python3 client_auth_key.py --help

# Utiliser un serveur distant
python3 client_auth_key.py --server https://tsd.example.com

# Passer le token en argument
python3 client_auth_key.py --token "votre-cle-api"

# Exécuter un exemple spécifique
python3 client_auth_key.py --example 1  # Utilisation basique
python3 client_auth_key.py --example 2  # Exécution de fichier
python3 client_auth_key.py --example 3  # Gestion d'erreurs
python3 client_auth_key.py --example 4  # Requêtes multiples
```

### client_jwt.py

```bash
# Afficher l'aide
python3 client_jwt.py --help

# Générer un JWT en Python
export TSD_JWT_SECRET="votre-secret"
python3 client_jwt.py --generate --username alice --roles "admin,user"

# Décoder un JWT
python3 client_jwt.py --decode "eyJhbGciOi..."

# Exécuter un exemple spécifique
python3 client_jwt.py --example 1  # Utilisation basique
python3 client_jwt.py --example 2  # Génération de JWT
python3 client_jwt.py --example 3  # Gestion expiration
python3 client_jwt.py --example 4  # Multi-utilisateurs
```

## 🐛 Dépannage

### Erreur: "Token d'authentification requis"

```bash
# Vérifier que la variable est définie
echo $TSD_AUTH_TOKEN

# Définir le token
export TSD_AUTH_TOKEN="votre-token"

# Ou passer en argument
python3 client_auth_key.py --token "votre-token"
```

### Erreur: "Authentification échouée"

```bash
# Vérifier que le serveur et le client utilisent le même token
# Pour Auth Key:
echo $TSD_AUTH_KEYS    # Serveur
echo $TSD_AUTH_TOKEN   # Client

# Pour JWT, vérifier le secret
echo $TSD_JWT_SECRET   # Doit être identique serveur et client
```

### Erreur: "Impossible de se connecter"

```bash
# Vérifier que le serveur est démarré
curl http://localhost:8080/health

# Avec authentification
curl -H "Authorization: Bearer $TSD_AUTH_TOKEN" http://localhost:8080/health
```

### Erreur: "PyJWT non installé"

```bash
# Installer PyJWT (nécessaire pour client_jwt.py avec --generate)
pip install PyJWT

# Ou utiliser tsd-auth pour générer les JWT
bin/tsd-auth generate-jwt -secret "votre-secret" -username alice
```

## 📖 Documentation

Pour plus d'informations, consultez :

- [Tutoriel complet d'authentification](../../docs/AUTHENTICATION_TUTORIAL.md)
- [Guide de démarrage rapide](../../docs/AUTHENTICATION_QUICKSTART.md)
- [Documentation du serveur](../../docs/SERVER_USAGE.md)

## 💡 Bonnes Pratiques

1. **Ne jamais hardcoder les secrets**
   ```python
   # ❌ Mauvais
   client = TSDAuthKeyClient(auth_token="ma-cle-secrete")
   
   # ✅ Bon
   client = TSDAuthKeyClient()  # Utilise TSD_AUTH_TOKEN
   ```

2. **Utiliser des variables d'environnement**
   ```bash
   # Stocker dans un fichier sécurisé
   echo "export TSD_AUTH_TOKEN='...'" > ~/.tsd_env
   chmod 600 ~/.tsd_env
   source ~/.tsd_env
   ```

3. **Gérer l'expiration des JWT**
   ```python
   try:
       result = client.execute(code)
   except Exception as e:
       if "expiré" in str(e).lower():
           # Régénérer le token
           new_token = generate_jwt(...)
           client.update_token(new_token)
           result = client.execute(code)
   ```

4. **Utiliser HTTPS en production**
   ```python
   client = TSDAuthKeyClient(
       server_url="https://tsd.prod.example.com"  # Pas http!
   )
   ```

## 🔒 Sécurité

- **Auth Key** : Utilisez des clés longues (générées par `tsd-auth generate-key`)
- **JWT Secret** : Minimum 32 caractères, aléatoire, jamais commité dans git
- **Rotation** : Changez régulièrement les clés et secrets
- **HTTPS** : Toujours en production pour éviter l'interception
- **Variables d'environnement** : Préférez-les aux arguments en ligne de commande

## 📝 Licence

Copyright (c) 2025 TSD Contributors - Licensed under MIT License