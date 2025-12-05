# 🔧 Binaire Unique TSD

## Vue d'ensemble

À partir de la version 1.0, TSD utilise un **binaire unique** multifonction qui remplace les anciens binaires séparés (`tsd`, `tsd-auth`, `tsd-client`, `tsd-server`).

Le binaire `tsd` change automatiquement de comportement selon son premier argument :

| Commande | Rôle | Description |
|----------|------|-------------|
| `tsd [fichier]` | **Compilateur/Runner** | Compile et exécute un programme TSD (comportement par défaut) |
| `tsd auth ...` | **Authentification** | Gestion des clés API et JWT |
| `tsd client ...` | **Client HTTP** | Client pour communiquer avec un serveur TSD distant |
| `tsd server ...` | **Serveur HTTP** | Serveur HTTP TSD avec authentification |

## Avantages

✅ **Simplicité** : Un seul binaire à installer et déployer  
✅ **Taille optimisée** : 12 MB (vs 31 MB pour les 4 binaires séparés)  
✅ **Facilité d'utilisation** : Interface cohérente avec dispatch automatique  
✅ **Maintenance** : Un seul point d'entrée à maintenir  
✅ **Distribution** : Packaging et distribution simplifiés  

## Installation

```bash
# Cloner et compiler
git clone https://github.com/treivax/tsd.git
cd tsd
make build

# Le binaire unique est créé dans ./bin/tsd
```

## Utilisation

### Aide Globale

```bash
# Afficher l'aide globale
tsd --help

# Afficher la version
tsd --version
```

### Aide Spécifique par Rôle

```bash
# Aide pour chaque rôle
tsd --help           # Aide du compilateur (comportement par défaut)
tsd auth --help      # Aide pour l'authentification
tsd client --help    # Aide pour le client HTTP
tsd server --help    # Aide pour le serveur HTTP
```

## Rôle 1 : Compilateur/Runner (Défaut)

Lorsqu'aucun rôle n'est spécifié, `tsd` fonctionne comme compilateur et runner de programmes TSD.

### Syntaxe

```bash
tsd <fichier.tsd> [options]
tsd -file <fichier.tsd> [options]
tsd -text "<code TSD>" [options]
tsd -stdin [options]
```

### Options

| Option | Description |
|--------|-------------|
| `-file <fichier>` | Fichier TSD à compiler |
| `-text <code>` | Code TSD directement en ligne de commande |
| `-stdin` | Lire le code depuis stdin |
| `-v` | Mode verbeux |
| `-version` | Afficher la version |

### Exemples

```bash
# Compiler un fichier
tsd program.tsd

# Mode verbeux
tsd program.tsd -v

# Lire depuis stdin
cat program.tsd | tsd -stdin

# Code TSD directement
tsd -text 'type Person : <id: string, name: string>'
```

## Rôle 2 : Authentification (auth)

Gestion des clés API et des JWT pour sécuriser le serveur TSD.

### Syntaxe

```bash
tsd auth <commande> [options]
```

### Commandes

| Commande | Description |
|----------|-------------|
| `generate-key` | Générer une ou plusieurs clés API |
| `generate-jwt` | Générer un JWT avec expiration |
| `validate` | Valider un token (Auth Key ou JWT) |
| `help` | Afficher l'aide |
| `version` | Afficher la version |

### Exemples

```bash
# Générer une clé API
tsd auth generate-key

# Générer plusieurs clés
tsd auth generate-key -count 3

# Générer un JWT
tsd auth generate-jwt \
  -secret "mon-secret-super-securise-de-32-chars" \
  -username alice \
  -roles "admin,user" \
  -expiration 48h

# Mode interactif (ne pas exposer le secret)
tsd auth generate-jwt -i -username alice

# Valider une clé API
tsd auth validate \
  -type key \
  -token "ma-cle-api" \
  -keys "cle1,cle2,cle3"

# Valider un JWT
tsd auth validate \
  -type jwt \
  -token "eyJhbG..." \
  -secret "mon-secret"

# Format JSON
tsd auth generate-key -format json
```

### Documentation Complète

Pour plus de détails sur l'authentification, consultez :
- [AUTHENTICATION.md](AUTHENTICATION.md) - Vue d'ensemble
- [AUTHENTICATION_TUTORIAL.md](AUTHENTICATION_TUTORIAL.md) - Tutoriel détaillé
- [AUTHENTICATION_QUICKSTART.md](AUTHENTICATION_QUICKSTART.md) - Guide rapide

## Rôle 3 : Client HTTP (client)

Client HTTP pour communiquer avec un serveur TSD distant.

### Syntaxe

```bash
tsd client <fichier.tsd> [options]
tsd client -file <fichier.tsd> [options]
tsd client -text "<code TSD>" [options]
tsd client -stdin [options]
```

### Options

| Option | Description |
|--------|-------------|
| `-server <url>` | URL du serveur TSD (défaut: http://localhost:8080) |
| `-file <fichier>` | Fichier TSD à exécuter |
| `-text <code>` | Code TSD directement |
| `-stdin` | Lire depuis stdin |
| `-token <token>` | Token d'authentification (Auth Key ou JWT) |
| `-format <format>` | Format de sortie: text ou json (défaut: text) |
| `-timeout <duration>` | Timeout des requêtes (défaut: 30s) |
| `-health` | Vérifier la santé du serveur |
| `-v` | Mode verbeux |

### Variables d'Environnement

| Variable | Description |
|----------|-------------|
| `TSD_AUTH_TOKEN` | Token d'authentification (alternative au flag `-token`) |

### Exemples

```bash
# Vérifier la santé du serveur
tsd client -health

# Exécuter un fichier TSD
tsd client program.tsd
tsd client -file program.tsd -v

# Exécuter du code TSD directement
tsd client -text 'type Person : <id: string, name: string>'

# Lire depuis stdin
echo 'type Person : <id: string>' | tsd client -stdin
cat program.tsd | tsd client -stdin -v

# Utiliser un serveur distant
tsd client -server http://tsd.example.com:8080 program.tsd

# Format JSON pour intégration
tsd client program.tsd -format json

# Avec authentification par clé API
tsd client program.tsd -token "votre-cle-api"

# Avec authentification JWT
tsd client program.tsd -token "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Via variable d'environnement
export TSD_AUTH_TOKEN="votre-token"
tsd client program.tsd
```

## Rôle 4 : Serveur HTTP (server)

Serveur HTTP TSD avec support de l'authentification.

### Syntaxe

```bash
tsd server [options]
```

### Options

| Option | Description |
|--------|-------------|
| `-host <host>` | Hôte du serveur (défaut: 0.0.0.0) |
| `-port <port>` | Port du serveur (défaut: 8080) |
| `-auth <type>` | Type d'authentification: none, key, jwt (défaut: none) |
| `-auth-keys <keys>` | Clés API (séparées par des virgules) |
| `-jwt-secret <secret>` | Secret pour JWT |
| `-jwt-expiration <duration>` | Durée de validité JWT (défaut: 24h) |
| `-jwt-issuer <issuer>` | Émetteur JWT (défaut: tsd-server) |
| `-v` | Mode verbeux |

### Variables d'Environnement

| Variable | Description |
|----------|-------------|
| `TSD_AUTH_KEYS` | Clés API (alternative au flag `-auth-keys`) |
| `TSD_JWT_SECRET` | Secret JWT (alternative au flag `-jwt-secret`) |

### Exemples

```bash
# Démarrer le serveur sans authentification (développement)
tsd server

# Serveur sur un port spécifique
tsd server -port 8080

# Serveur avec authentification par clé API
tsd server -auth key -auth-keys "cle1,cle2,cle3"

# Serveur avec authentification JWT
tsd server -auth jwt -jwt-secret "mon-secret-super-securise-de-32-chars"

# Configuration complète JWT
tsd server \
  -auth jwt \
  -jwt-secret "mon-secret" \
  -jwt-expiration 48h \
  -jwt-issuer "mon-entreprise" \
  -port 8080 \
  -v

# Via variables d'environnement (recommandé en production)
export TSD_JWT_SECRET="mon-secret-super-securise-de-32-chars"
tsd server -auth jwt -port 8080
```

### Endpoints Disponibles

| Endpoint | Méthode | Description |
|----------|---------|-------------|
| `/api/v1/execute` | POST | Exécuter un programme TSD |
| `/health` | GET | Health check |
| `/api/v1/version` | GET | Informations de version |

## Architecture Interne

Le binaire unique utilise une architecture modulaire avec dispatch dynamique :

```
cmd/tsd/main.go (point d'entrée unique)
├── determineRole() → Analyse le premier argument
└── dispatch() → Redirige vers le package approprié
    ├── internal/compilercmd/ → Compilateur/Runner
    ├── internal/authcmd/     → Gestion d'authentification
    ├── internal/clientcmd/   → Client HTTP
    └── internal/servercmd/   → Serveur HTTP
```

### Packages Internes

| Package | Description |
|---------|-------------|
| `internal/compilercmd/` | Logique du compilateur et runner TSD |
| `internal/authcmd/` | Gestion des clés API et JWT |
| `internal/clientcmd/` | Client HTTP pour communiquer avec le serveur |
| `internal/servercmd/` | Serveur HTTP avec authentification |

## Migration depuis les Binaires Séparés

Si vous utilisiez les anciens binaires séparés, voici comment migrer :

### Avant (binaires séparés)

```bash
# Compilateur
./bin/tsd program.tsd

# Authentification
./bin/tsd-auth generate-key

# Client
./bin/tsd-client program.tsd

# Serveur
./bin/tsd-server -port 8080
```

### Après (binaire unique)

```bash
# Compilateur (identique ou sans rôle)
./bin/tsd program.tsd

# Authentification (préfixe "auth")
./bin/tsd auth generate-key

# Client (préfixe "client")
./bin/tsd client program.tsd

# Serveur (préfixe "server")
./bin/tsd server -port 8080
```

### Scripts et CI/CD

Si vous avez des scripts qui utilisent les anciens binaires, il suffit d'ajouter le rôle approprié :

```bash
# Avant
tsd-auth generate-key

# Après
tsd auth generate-key
```

## Déploiement

### Docker

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bin/tsd /usr/local/bin/tsd

# Utiliser le rôle approprié au démarrage
CMD ["tsd", "server", "-port", "8080"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tsd-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: tsd-server
  template:
    metadata:
      labels:
        app: tsd-server
    spec:
      containers:
      - name: tsd
        image: tsd:latest
        command: ["tsd", "server"]
        args:
          - "-auth"
          - "jwt"
          - "-jwt-secret"
          - "$(JWT_SECRET)"
          - "-port"
          - "8080"
        env:
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: tsd-secrets
              key: jwt-secret
        ports:
        - containerPort: 8080
```

### Systemd

```ini
[Unit]
Description=TSD Server
After=network.target

[Service]
Type=simple
User=tsd
WorkingDirectory=/opt/tsd
ExecStart=/usr/local/bin/tsd server -auth jwt -jwt-secret "${JWT_SECRET}" -port 8080
Environment="JWT_SECRET=mon-secret-super-securise"
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

## Tests

Le binaire unique inclut des tests pour vérifier le dispatch des rôles :

```bash
# Tests unitaires du dispatcher
go test -v ./cmd/tsd

# Tests d'intégration
go test -v ./cmd/tsd -run TestDetermineRole
go test -v ./cmd/tsd -run TestDispatchLogic

# Build et tests complets
make validate
```

## Taille du Binaire

Comparaison de la taille avant et après la refactorisation :

| Configuration | Taille Totale | Détails |
|---------------|---------------|---------|
| **Avant** (4 binaires séparés) | **31 MB** | tsd (6.7 MB) + tsd-auth (4.7 MB) + tsd-client (8.5 MB) + tsd-server (11 MB) |
| **Après** (binaire unique) | **12 MB** | tsd (12 MB) - Tout inclus |
| **Réduction** | **-61%** | 19 MB économisés |

L'optimisation provient de la mutualisation du code commun entre les différents rôles.

## Compatibilité

Le binaire unique est **100% compatible** avec les anciennes fonctionnalités :

- ✅ Toutes les options sont préservées
- ✅ Variables d'environnement identiques
- ✅ Format de sortie inchangé
- ✅ APIs et protocoles compatibles
- ✅ Migration transparente

## Support et Documentation

- **README principal** : [README.md](../README.md)
- **Authentification** : [AUTHENTICATION.md](AUTHENTICATION.md)
- **Tutoriel Auth** : [AUTHENTICATION_TUTORIAL.md](AUTHENTICATION_TUTORIAL.md)
- **Guide Rapide Auth** : [AUTHENTICATION_QUICKSTART.md](AUTHENTICATION_QUICKSTART.md)
- **Changelog** : [CHANGELOG.md](../CHANGELOG.md)

## Contribution

Pour contribuer au développement du binaire unique :

1. Le code source est dans `cmd/tsd/main.go`
2. Les packages internes sont dans `internal/*/`
3. Les tests sont dans `cmd/tsd/*_test.go`
4. Suivez les conventions du prompt `.github/prompts/add-feature.md`

## Licence

Copyright (c) 2025 TSD Contributors  
Licence: MIT - Voir [LICENSE](../LICENSE)