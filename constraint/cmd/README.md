# Constraint Parser CLI

CLI pour parser et valider des fichiers de contraintes TSD.

## Installation

```bash
cd constraint/cmd
go build -o constraint-parser .
```

## Utilisation

### Syntaxe de base

```bash
constraint-parser [options] <input-file>
```

### Arguments

- `<input-file>` : Fichier de contraintes à parser (utilisez `-` pour stdin)

### Options

- `--config PATH` : Chemin vers le fichier de configuration JSON
- `--output FORMAT` : Format de sortie (défaut: json)
- `--debug` : Activer le mode debug
- `--version` : Afficher la version
- `--help` : Afficher l'aide

## Configuration

### Fichier de configuration

Créez un fichier JSON avec la structure suivante :

```json
{
  "parser": {
    "max_expressions": 1000,
    "debug": false,
    "recover": true
  },
  "validator": {
    "strict_mode": true,
    "allowed_operators": ["==", "!=", "<", ">", "<=", ">=", "AND", "OR", "NOT"],
    "max_depth": 20
  },
  "logger": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  },
  "debug": false,
  "version": "1.0.0"
}
```

Voir `example-config.json` pour un exemple complet.

### Variables d'environnement

Les variables d'environnement suivantes sont supportées et surchargent les valeurs du fichier de configuration :

| Variable | Type | Description | Défaut |
|----------|------|-------------|--------|
| `CONSTRAINT_CONFIG_FILE` | string | Chemin du fichier de config | "" |
| `CONSTRAINT_MAX_EXPRESSIONS` | int | Nombre max d'expressions | 1000 |
| `CONSTRAINT_MAX_DEPTH` | int | Profondeur max de validation | 20 |
| `CONSTRAINT_DEBUG` | bool | Mode debug | false |
| `CONSTRAINT_STRICT_MODE` | bool | Mode strict validation | true |
| `CONSTRAINT_LOG_LEVEL` | string | Niveau de log (debug/info/warn/error) | info |
| `CONSTRAINT_LOG_FORMAT` | string | Format de log (json/text/plain) | json |
| `CONSTRAINT_LOG_OUTPUT` | string | Sortie des logs | stdout |

### Priorité de configuration

L'ordre de priorité (du plus faible au plus fort) :

1. Valeurs par défaut
2. Fichier de configuration
3. Variables d'environnement
4. Flags CLI

## Exemples

### Utilisation basique

```bash
# Parser un fichier
constraint-parser constraints.tsd

# Avec un fichier de configuration
constraint-parser --config myconfig.json constraints.tsd

# Mode debug
constraint-parser --debug constraints.tsd
```

### Support stdin

```bash
# Depuis stdin
cat constraints.tsd | constraint-parser -

# Depuis echo
echo 'type Person(id: string, name: string)' | constraint-parser -

# Pipeline
curl https://example.com/constraints.tsd | constraint-parser -
```

### Variables d'environnement

```bash
# Activer le mode debug
CONSTRAINT_DEBUG=true constraint-parser constraints.tsd

# Augmenter la limite d'expressions
CONSTRAINT_MAX_EXPRESSIONS=5000 constraint-parser constraints.tsd

# Changer le niveau de log
CONSTRAINT_LOG_LEVEL=debug constraint-parser constraints.tsd

# Configuration complète via ENV
export CONSTRAINT_CONFIG_FILE=/etc/constraint/config.json
export CONSTRAINT_DEBUG=true
export CONSTRAINT_MAX_EXPRESSIONS=2000
constraint-parser constraints.tsd
```

### Version et aide

```bash
# Afficher la version
constraint-parser --version

# Afficher l'aide complète
constraint-parser --help
```

## Codes de sortie

| Code | Signification | Description |
|------|---------------|-------------|
| 0 | Success | Parsing et validation réussis |
| 1 | Usage Error | Arguments invalides ou manquants |
| 2 | Runtime Error | Erreur de parsing ou validation |
| 3 | Config Error | Erreur de configuration |

## Format de sortie

Le CLI produit du JSON sur stdout :

```json
{
  "types": [
    {
      "name": "Person",
      "fields": [
        {"name": "id", "type": "string"},
        {"name": "name", "type": "string"}
      ]
    }
  ],
  "expressions": [],
  "facts": [],
  "actions": [],
  "resets": [],
  "retractions": [],
  "ruleRemovals": []
}
```

Les messages d'erreur et de validation sont envoyés sur stderr.

## Développement

### Tester

```bash
# Tests unitaires
go test ./...

# Tests avec couverture
go test -cover ./...

# Tests verbose
go test -v ./...
```

### Linter

```bash
# Vérification statique
go vet ./...

# Formater le code
go fmt ./...
```

## Architecture

### Configuration

La configuration utilise un système de couches avec priorité :

1. **Défauts** : Constantes définies dans `config.go`
2. **Fichier** : Chargé via `--config` ou `CONSTRAINT_CONFIG_FILE`
3. **ENV** : Variables d'environnement `CONSTRAINT_*`
4. **Flags** : Options CLI (ex: `--debug`)

Chaque couche surcharge la précédente.

### Workflow

```
Arguments CLI
    ↓
Parse flags
    ↓
Load config (défaut → fichier → ENV → flags)
    ↓
Validate config
    ↓
Parse input (file ou stdin)
    ↓
Validate constraint program
    ↓
Output JSON
```

## Conformité

Le CLI respecte les principes 12-factor app :

- ✅ Configuration externalisée (ENV + fichiers)
- ✅ Processes stateless
- ✅ Logs vers stdout/stderr
- ✅ Codes de sortie explicites
- ✅ Build/run/config séparés

## Support

Pour toute question ou problème, consulter :

- Documentation projet : `/docs`
- Standards : `.github/prompts/common.md`
- Tests : `cmd/*_test.go`

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License
