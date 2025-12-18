# Syntaxe de la Commande xuple-space

## 🎯 Objectif

Définir la syntaxe précise et complète de la commande `xuple-space` pour le langage TSD.

## 📝 Syntaxe Complète

### Format Général

```ebnf
XupleSpaceDeclaration ::= "xuple-space" Identifier "{" XupleSpaceBody "}"

XupleSpaceBody ::= (XupleSpaceProperty)*

XupleSpaceProperty ::= SelectionPolicy | ConsumptionPolicy | RetentionPolicy

SelectionPolicy ::= "selection" ":" SelectionValue

SelectionValue ::= "random" | "fifo" | "lifo"

ConsumptionPolicy ::= "consumption" ":" ConsumptionValue

ConsumptionValue ::= "once" | "per-agent" | "limited" "(" Integer ")"

RetentionPolicy ::= "retention" ":" RetentionValue

RetentionValue ::= "unlimited" | "duration" "(" Duration ")"

Duration ::= Integer TimeUnit

TimeUnit ::= "s" | "m" | "h" | "d"

Identifier ::= [a-zA-Z_][a-zA-Z0-9_-]*

Integer ::= [0-9]+
```

## 🔤 Syntaxe Textuelle

```tsd
xuple-space <nom> {
    selection: <random|fifo|lifo>
    consumption: <once|per-agent|limited(<n>)>
    retention: <unlimited|duration(<temps>)>
}
```

## 📋 Descriptions des Politiques

### 1. Selection Policy (Politique de Sélection)

Détermine quel xuple est sélectionné parmi plusieurs disponibles.

| Valeur | Description | Comportement |
|--------|-------------|--------------|
| `random` | Sélection aléatoire | Choisit un xuple au hasard |
| `fifo` | First-In-First-Out | Choisit le xuple le plus ancien (CreatedAt min) |
| `lifo` | Last-In-First-Out | Choisit le xuple le plus récent (CreatedAt max) |

**Valeur par défaut** : `fifo`

### 2. Consumption Policy (Politique de Consommation)

Détermine combien de fois un xuple peut être consommé.

| Valeur | Description | Comportement | Transition Statut |
|--------|-------------|--------------|-------------------|
| `once` | Une seule consommation totale | Le xuple ne peut être consommé qu'une fois au total | Pending → Consumed |
| `per-agent` | Une fois par agent | Chaque agent peut consommer le xuple une fois | Pending → Pending |
| `limited(n)` | n consommations max | Le xuple peut être consommé n fois au total | Pending → Consumed (après n) |

**Paramètres** :
- `limited(n)` : `n` doit être un entier positif > 0

**Valeur par défaut** : `once`

### 3. Retention Policy (Politique de Rétention)

Détermine combien de temps un xuple est conservé.

| Valeur | Description | Comportement |
|--------|-------------|--------------|
| `unlimited` | Conservation illimitée | Les xuples ne sont jamais supprimés automatiquement |
| `duration(temps)` | Expire après un délai | Les xuples sont marqués expirés après le délai |

**Format durée** :
- `s` : secondes
- `m` : minutes  
- `h` : heures
- `d` : jours

**Exemples** :
- `duration(30s)` : 30 secondes
- `duration(5m)` : 5 minutes
- `duration(1h)` : 1 heure
- `duration(7d)` : 7 jours

**Valeur par défaut** : `unlimited`

## ✅ Exemples Valides

### Exemple 1 : Configuration Minimale (Valeurs par Défaut)

```tsd
xuple-space simple {
    selection: fifo
    consumption: once
    retention: unlimited
}
```

**Comportement** :
- Sélection du xuple le plus ancien
- Une seule consommation
- Conservation permanente

### Exemple 2 : File d'Attente de Commandes

```tsd
xuple-space agents-commands {
    selection: fifo
    consumption: once
    retention: duration(1h)
}
```

**Comportement** :
- FIFO : traitement dans l'ordre d'arrivée
- Une commande consommée une seule fois
- Expiration après 1 heure

### Exemple 3 : Notifications Multi-Agents

```tsd
xuple-space notifications {
    selection: random
    consumption: per-agent
    retention: duration(5m)
}
```

**Comportement** :
- Sélection aléatoire
- Chaque agent peut consulter la notification
- Expiration après 5 minutes

### Exemple 4 : Traitement Limité

```tsd
xuple-space limited-processing {
    selection: lifo
    consumption: limited(3)
    retention: unlimited
}
```

**Comportement** :
- Traitement du xuple le plus récent
- Maximum 3 consommations
- Conservation permanente

### Exemple 5 : Données Éphémères

```tsd
xuple-space cache {
    selection: fifo
    consumption: limited(10)
    retention: duration(30s)
}
```

**Comportement** :
- FIFO pour cohérence temporelle
- Jusqu'à 10 lectures
- Expiration rapide (30 secondes)

## ❌ Exemples Invalides

### Erreur 1 : Politique de sélection invalide

```tsd
xuple-space bad {
    selection: priority  // ❌ Invalide - doit être random|fifo|lifo
}
```

**Message d'erreur attendu** :
```
Parse error at line 2: invalid selection policy 'priority', expected random, fifo or lifo
```

### Erreur 2 : Limite zéro

```tsd
xuple-space bad {
    consumption: limited(0)  // ❌ Invalide - doit être > 0
}
```

**Message d'erreur attendu** :
```
Parse error at line 2: consumption limit must be greater than zero
```

### Erreur 3 : Durée négative

```tsd
xuple-space bad {
    retention: duration(-5m)  // ❌ Invalide - doit être > 0
}
```

**Message d'erreur attendu** :
```
Parse error at line 2: duration must be positive
```

### Erreur 4 : Unité de temps invalide

```tsd
xuple-space bad {
    retention: duration(5x)  // ❌ Invalide - doit être s|m|h|d
}
```

**Message d'erreur attendu** :
```
Parse error at line 2: invalid time unit 'x', expected s, m, h or d
```

### Erreur 5 : Nom de xuple-space dupliqué

```tsd
xuple-space myspace {
    selection: fifo
}

xuple-space myspace {  // ❌ Invalide - nom déjà utilisé
    selection: lifo
}
```

**Message d'erreur attendu** :
```
Compilation error: xuple-space 'myspace' already declared at line 1
```

### Erreur 6 : Propriété manquante

```tsd
xuple-space incomplete {
    selection: fifo
    // ❌ Manque consumption et retention
}
```

**Comportement** : Les propriétés manquantes utilisent les valeurs par défaut.

### Erreur 7 : Syntaxe de paramètre incorrecte

```tsd
xuple-space bad {
    consumption: limited 5  // ❌ Invalide - doit être limited(5)
}
```

**Message d'erreur attendu** :
```
Parse error at line 2: expected '(' after 'limited'
```

## 🔄 Valeurs par Défaut

Si une propriété n'est pas spécifiée, les valeurs par défaut s'appliquent :

```go
const (
    DefaultSelectionPolicy   = "fifo"
    DefaultConsumptionPolicy = "once"
    DefaultRetentionPolicy   = "unlimited"
)
```

**Exemple minimal valide** :

```tsd
xuple-space minimal {}
```

**Équivaut à** :

```tsd
xuple-space minimal {
    selection: fifo
    consumption: once
    retention: unlimited
}
```

## 🎨 Cas d'Usage Recommandés

### File de Travail (Job Queue)

```tsd
xuple-space job-queue {
    selection: fifo
    consumption: once
    retention: duration(24h)
}
```

**Caractéristiques** :
- Traitement FIFO pour équité
- Job consommé une fois
- Nettoyage après 24h

### Publish-Subscribe

```tsd
xuple-space pubsub {
    selection: random
    consumption: per-agent
    retention: duration(10m)
}
```

**Caractéristiques** :
- Pas d'ordre garanti
- Tous les abonnés peuvent lire
- Expiration pour libérer mémoire

### Cache Distribué

```tsd
xuple-space distributed-cache {
    selection: fifo
    consumption: limited(100)
    retention: duration(1h)
}
```

**Caractéristiques** :
- FIFO pour LRU-like
- Limite de lectures
- Expiration temporelle

### Stack de Traitement

```tsd
xuple-space processing-stack {
    selection: lifo
    consumption: once
    retention: unlimited
}
```

**Caractéristiques** :
- LIFO pour traiter les plus récents
- Traitement unique
- Conservation pour audit

## 🔍 Validation Syntaxique

### Règles de Validation

1. **Nom** : Doit être un identifiant valide (`[a-zA-Z_][a-zA-Z0-9_-]*`)
2. **Unicité** : Pas de duplication de noms de xuple-space
3. **Selection** : Doit être `random`, `fifo` ou `lifo`
4. **Consumption** :
   - `once` ou `per-agent` : pas de paramètre
   - `limited(n)` : `n` entier > 0
5. **Retention** :
   - `unlimited` : pas de paramètre
   - `duration(temps)` : temps entier > 0 avec unité valide (s, m, h, d)

### Validation au Runtime

Lors de la compilation, vérifier :
- [ ] Nom du xuple-space unique
- [ ] Valeurs des politiques dans les ensembles autorisés
- [ ] Paramètres numériques > 0
- [ ] Unités de temps valides

## 📊 Mapping vers Structures Go

### Structure AST

```go
type XupleSpaceDeclaration struct {
    Type              string                   `json:"type"`              // "xupleSpaceDeclaration"
    Name              string                   `json:"name"`              // Nom du xuple-space
    SelectionPolicy   string                   `json:"selectionPolicy"`   // "random", "fifo", "lifo"
    ConsumptionPolicy ConsumptionPolicyConfig  `json:"consumptionPolicy"` // Configuration consommation
    RetentionPolicy   RetentionPolicyConfig    `json:"retentionPolicy"`   // Configuration rétention
}

type ConsumptionPolicyConfig struct {
    Type  string `json:"type"`            // "once", "per-agent", "limited"
    Limit int    `json:"limit,omitempty"` // Pour "limited", sinon 0
}

type RetentionPolicyConfig struct {
    Type     string `json:"type"`               // "unlimited", "duration"
    Duration int    `json:"duration,omitempty"` // En secondes, pour "duration"
}
```

### Conversion des Durées

```go
const (
    SecondUnit = 1
    MinuteUnit = 60
    HourUnit   = 3600
    DayUnit    = 86400
)

func ParseDuration(value int, unit string) (int, error) {
    switch unit {
    case "s":
        return value * SecondUnit, nil
    case "m":
        return value * MinuteUnit, nil
    case "h":
        return value * HourUnit, nil
    case "d":
        return value * DayUnit, nil
    default:
        return 0, fmt.Errorf("invalid time unit '%s'", unit)
    }
}
```

## ✅ Checklist Implémentation

- [ ] Règles PEG définies
- [ ] Structures Go créées avec copyright
- [ ] Validation des valeurs implémentée
- [ ] Messages d'erreur clairs
- [ ] Tests pour tous les cas valides
- [ ] Tests pour tous les cas d'erreur
- [ ] Documentation utilisateur créée
- [ ] Exemples fournis

## 📚 Références

- `01-parser-analysis.md` - Analyse du parser
- `constraint/grammar/constraint.peg` - Grammaire PEG
- `constraint/constraint_types.go` - Types AST
- `xuples/policies.go` - Implémentation des politiques
