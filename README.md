# Projet TSD (Type System & Data)

Projet combinant un client etcd Go et un module de système de contraintes avec grammaire PEG personnalisée.

## Composants du projet

### 🔧 Client etcd
Programme Go qui se connecte à un service etcd et gère les clés/valeurs.

### 📝 Module Constraint  
Système de contraintes avec :
- Grammaire PEG personnalisée
- Parser automatique 
- Validation de types et règles
- API Go réutilisable

## Prérequis

1. **Go installé** (version 1.21 ou supérieure)
   ```bash
   # Sur Ubuntu/Debian
   sudo apt update
   sudo apt install golang-go
   
   # Sur CentOS/RHEL
   sudo yum install golang
   
   # Ou télécharger depuis https://golang.org/dl/
   ```

2. **Service etcd en cours d'exécution**
   ```bash
   # Installation d'etcd (exemple sur Ubuntu)
   sudo apt install etcd
   
   # Ou utiliser Docker
   docker run -d -p 2379:2379 -p 2380:2380 \
     --name etcd-server \
     quay.io/coreos/etcd:latest \
     etcd \
     --advertise-client-urls http://0.0.0.0:2379 \
     --listen-client-urls http://0.0.0.0:2379
   ```

## Configuration

Par défaut, le programme se connecte à etcd sur `localhost:2379`. 

Pour modifier l'adresse de connexion, éditez le fichier `main.go` et changez la ligne :
```go
Endpoints: []string{"localhost:2379"},
```

## Structure du projet

```
tsd/
├── README.md              # Cette documentation
├── build.sh              # Script de build principal
├── go.mod               # Dépendances du client etcd
├── main.go              # Client etcd principal
├── operations.go        # Opérations etcd
├── put.go              # Opérations PUT
├── take.go             # Opérations TAKE
└── constraint/         # 📁 MODULE CONSTRAINT
    ├── README.md           # Documentation du module
    ├── build.sh           # Build spécifique au module
    ├── go.mod            # Configuration du module
    ├── api.go            # API publique
    ├── constraint_types.go  # Types d'AST
    ├── constraint_utils.go  # Utilitaires de validation
    ├── parser.go         # Parser généré (ne pas modifier)
    ├── cmd/              # Exécutable
    │   ├── go.mod
    │   └── main.go
    ├── grammar/          # Grammaires
    │   ├── constraint.peg    # Grammaire PEG source
    │   └── SetConstraint.g4  # Grammaire ANTLR
    ├── tests/            # Fichiers de test
    │   ├── test_input.txt
    │   ├── test_type_valid.txt
    │   └── ... (autres tests)
    └── docs/             # Documentation
        ├── GUIDE_CONTRAINTES.md
        ├── TUTORIEL_CONTRAINTES.md
        └── PARSER_README.md
```

## Installation et compilation

### Prérequis globaux

1. **Go 1.21+**
2. **Pigeon** (pour le module constraint)
   ```bash
   go install github.com/mna/pigeon@latest
   ```
3. **etcd** (pour le client)

### Build complet

```bash
# Build de tous les composants
./build.sh

# Ou builds séparés
cd constraint && ./build.sh  # Module constraint
go build -o etcd-client main.go operations.go put.go take.go  # Client etcd

# Exécuter le programme
./etcd-client
```

Ou directement :
```bash
go run main.go
```

## Utilisation

Le programme :
1. Se connecte au service etcd
2. Recherche toutes les clés avec le préfixe `/a/b/c`
3. Affiche les détails de chaque clé trouvée (clé, valeur, version, etc.)

## Exemple de sortie

```
Connexion réussie à etcd!
Récupération des clés avec le préfixe '/a/b/c'...

Nombre de clés trouvées avec le préfixe '/a/b/c': 3

Clés trouvées:
==============
1. Clé: /a/b/c/key1
   Valeur: value1
   Version: 1
   Créée à: 12345
   Modifiée à: 12345
   ---
2. Clé: /a/b/c/config/setting
   Valeur: {"enabled": true}
   Version: 2
   Créée à: 12346
   Modifiée à: 12347
   ---
```

## Test avec des données d'exemple

Pour tester le programme, vous pouvez ajouter des clés de test dans etcd :

```bash
# Installer etcdctl
sudo apt install etcd-client

# Ajouter des clés de test
etcdctl put /a/b/c/test1 "valeur de test 1"
etcdctl put /a/b/c/test2 "valeur de test 2"
etcdctl put /a/b/c/config/debug "true"
etcdctl put /other/key "cette clé ne sera pas listée"

# Exécuter le programme
go run main.go
```

## Gestion des erreurs

Le programme gère plusieurs types d'erreurs :
- Connexion impossible à etcd
- Timeout de connexion (5 secondes par défaut)
- Erreurs lors de la récupération des clés

En cas d'erreur, le programme affichera un message d'erreur explicite et se terminera.