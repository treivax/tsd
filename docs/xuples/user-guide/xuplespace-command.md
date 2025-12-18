# Commande xuple-space - Guide Utilisateur

## 📋 Vue d'Ensemble

La commande `xuple-space` permet de déclarer des espaces de xuples (xuple-spaces) dans les fichiers TSD. Un xuple-space est un espace partagé où les activations de règles RETE sont publiées et peuvent être consommées par des agents externes selon des politiques configurables.

## 🔤 Syntaxe

```tsd
xuple-space <nom> {
    selection: <random|fifo|lifo>
    consumption: <once|per-agent|limited(n)>
    retention: <unlimited|duration(temps)>
}
```

### Paramètres

- **nom** : Identifiant unique du xuple-space (caractères alphanumériques, tirets et underscores)
- **selection** : Politique de sélection des xuples
- **consumption** : Politique de consommation des xuples
- **retention** : Politique de rétention/expiration des xuples

## 📊 Politiques

### 1. Selection Policy (Sélection)

Détermine quel xuple est sélectionné parmi plusieurs disponibles lors d'une consommation.

#### `fifo` - First-In-First-Out (Par défaut)

Sélectionne le xuple le plus ancien (créé en premier).

**Cas d'usage** :
- Files de travail équitables
- Traitement dans l'ordre chronologique
- Garantie d'ordre de traitement

**Exemple** :
```tsd
xuple-space job-queue {
    selection: fifo
    consumption: once
    retention: unlimited
}
```

#### `lifo` - Last-In-First-Out

Sélectionne le xuple le plus récent (créé en dernier).

**Cas d'usage** :
- Pile de traitement (stack)
- Priorité aux événements les plus récents
- Systèmes d'alerte (traiter les alertes récentes d'abord)

**Exemple** :
```tsd
xuple-space alert-stack {
    selection: lifo
    consumption: once
    retention: duration(5m)
}
```

#### `random` - Sélection Aléatoire

Sélectionne un xuple au hasard parmi les disponibles.

**Cas d'usage** :
- Load balancing (répartition de charge)
- Distribution équitable sans ordre
- Réduire les contentions

**Exemple** :
```tsd
xuple-space load-balancer {
    selection: random
    consumption: once
    retention: unlimited
}
```

### 2. Consumption Policy (Consommation)

Détermine combien de fois un xuple peut être consommé.

#### `once` - Consommation Unique (Par défaut)

Le xuple ne peut être consommé qu'une seule fois au total.

**Comportement** :
- Statut : `Pending` → `Consumed` après consommation
- Le xuple n'est plus sélectionnable après consommation

**Cas d'usage** :
- Commandes à exécution unique
- Traitement de job
- Actions non-répétables

**Exemple** :
```tsd
xuple-space commands {
    selection: fifo
    consumption: once
    retention: duration(1h)
}
```

#### `per-agent` - Une Fois par Agent

Chaque agent peut consommer le xuple une fois.

**Comportement** :
- Le xuple reste `Pending` après consommation
- Chaque agent peut le consommer une seule fois
- Idéal pour pattern publish-subscribe

**Cas d'usage** :
- Notifications broadcast
- Événements multi-agents
- Synchronisation distribuée

**Exemple** :
```tsd
xuple-space notifications {
    selection: random
    consumption: per-agent
    retention: duration(10m)
}
```

#### `limited(n)` - Consommation Limitée

Le xuple peut être consommé jusqu'à `n` fois.

**Paramètres** :
- `n` : Nombre maximum de consommations (entier positif > 0)

**Comportement** :
- Statut : `Pending` → `Consumed` après n consommations
- Compteur de consommations incrémenté à chaque consommation

**Cas d'usage** :
- Cache avec quota de lectures
- Réplication limitée
- Partage de ressources contrôlé

**Exemple** :
```tsd
xuple-space cache {
    selection: fifo
    consumption: limited(10)
    retention: duration(5m)
}
```

### 3. Retention Policy (Rétention)

Détermine combien de temps un xuple est conservé avant expiration.

#### `unlimited` - Rétention Illimitée (Par défaut)

Les xuples ne sont jamais supprimés automatiquement.

**Comportement** :
- Aucune expiration temporelle
- Les xuples restent jusqu'à consommation ou suppression manuelle
- Attention à la consommation mémoire

**Cas d'usage** :
- Archivage permanent
- Audit trail
- Données historiques

**Exemple** :
```tsd
xuple-space archive {
    selection: fifo
    consumption: once
    retention: unlimited
}
```

#### `duration(temps)` - Expiration Temporelle

Les xuples expirent après une durée spécifiée.

**Format de durée** :
- `s` : secondes
- `m` : minutes
- `h` : heures
- `d` : jours

**Comportement** :
- Statut : `Pending`/`Consumed` → `Expired` après expiration
- Nettoyage automatique via `Cleanup()`
- Le champ `ExpiresAt` est défini lors de la création

**Cas d'usage** :
- Cache temporaire
- TTL (Time To Live)
- Données éphémères
- Prévention de fuite mémoire

**Exemples** :
```tsd
// Cache court (30 secondes)
xuple-space short-cache {
    selection: random
    consumption: limited(5)
    retention: duration(30s)
}

// Cache moyen (5 minutes)
xuple-space medium-cache {
    selection: fifo
    consumption: per-agent
    retention: duration(5m)
}

// Données journalières (1 heure)
xuple-space hourly-data {
    selection: fifo
    consumption: once
    retention: duration(1h)
}

// Archive hebdomadaire (7 jours)
xuple-space weekly-archive {
    selection: fifo
    consumption: once
    retention: duration(7d)
}
```

## 🎯 Patterns d'Utilisation Recommandés

### File de Travail (Job Queue)

```tsd
xuple-space job-queue {
    selection: fifo
    consumption: once
    retention: duration(24h)
}
```

**Caractéristiques** :
- FIFO pour traitement équitable
- Une seule exécution par job
- Expiration après 24h (prévention de blocage)

### Publish-Subscribe

```tsd
xuple-space pubsub {
    selection: random
    consumption: per-agent
    retention: duration(10m)
}
```

**Caractéristiques** :
- Sélection aléatoire (pas d'ordre imposé)
- Chaque agent peut lire
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
- FIFO pour cohérence
- Limite de lectures
- TTL d'1 heure

### Système d'Alerte

```tsd
xuple-space alerts {
    selection: lifo
    consumption: per-agent
    retention: duration(2m)
}
```

**Caractéristiques** :
- LIFO pour traiter alertes récentes
- Broadcast aux agents
- Courte durée de vie

## 📝 Exemples Complets

### Exemple 1 : Système Simple

```tsd
// Types
type Task(#id: string, title: string, priority: number)

// Actions
action processTask(taskId: string, title: string)

// Xuple-space
xuple-space tasks {
    selection: fifo
    consumption: once
    retention: unlimited
}

// Règle
rule process_task: {t: Task} / t.priority > 5 ==> processTask(t.id, t.title)

// Faits
Task(id: "T001", title: "Urgent", priority: 10)
```

### Exemple 2 : Multi-Agent Notifications

```tsd
// Types
type Notification(#id: string, message: string, severity: number)

// Actions
action notify(message: string)

// Xuple-space pour broadcast
xuple-space notifications {
    selection: random
    consumption: per-agent
    retention: duration(5m)
}

// Règle
rule send_notification: {n: Notification} / n.severity >= 5 ==> notify(n.message)

// Faits
Notification(id: "N001", message: "System update required", severity: 8)
```

### Exemple 3 : Cache avec Quota

```tsd
// Types
type CacheEntry(#key: string, value: string)

// Actions
action accessCache(key: string, value: string)

// Xuple-space avec limite de lectures
xuple-space cache {
    selection: fifo
    consumption: limited(10)
    retention: duration(1h)
}

// Règle
rule cache_access: {c: CacheEntry} / ==> accessCache(c.key, c.value)

// Faits
CacheEntry(key: "user:123", value: "John Doe")
```

## ⚠️ Bonnes Pratiques

### Nommage

- Utiliser des noms descriptifs : `job-queue`, `user-notifications`, `cache-entries`
- Éviter les noms génériques : `space1`, `temp`, `data`
- Préférer kebab-case : `my-xuple-space`

### Choix des Politiques

#### Selection Policy
- **FIFO** : Quand l'ordre chronologique est important
- **LIFO** : Pour priorité aux événements récents
- **Random** : Pour load balancing sans ordre

#### Consumption Policy
- **once** : Pour actions non-répétables
- **per-agent** : Pour broadcast/pub-sub
- **limited(n)** : Pour contrôle de quota

#### Retention Policy
- **unlimited** : Pour archivage ou audit
- **duration(...)** : Pour données éphémères ou cache
  - Courte (< 1m) : Cache très volatile
  - Moyenne (1m-1h) : Cache applicatif
  - Longue (> 1h) : Archivage temporaire

### Performance

- Éviter `unlimited` avec `per-agent` sans `duration` (risque de fuite mémoire)
- Préférer des durées courtes pour cache
- Utiliser `limited()` pour contrôler la charge

### Sécurité

- Ne pas stocker de données sensibles sans expiration
- Limiter la consommation avec `limited()` si nécessaire
- Utiliser `duration()` pour auto-nettoyage

## 🔍 Validation

Le parser valide automatiquement :

✅ **Validations syntaxiques** :
- Nom du xuple-space unique dans le fichier
- Politique de sélection valide (random, fifo, lifo)
- Limite de consommation > 0
- Durée de rétention > 0
- Unité de temps valide (s, m, h, d)

❌ **Erreurs détectées** :
```tsd
// ❌ Politique invalide
xuple-space bad { selection: priority }

// ❌ Limite zéro
xuple-space bad { consumption: limited(0) }

// ❌ Durée négative
xuple-space bad { retention: duration(-5m) }

// ❌ Unité invalide
xuple-space bad { retention: duration(5x) }
```

## 🛠️ TODO : Actions Ultérieures

> **Note** : Cette implémentation concerne uniquement le **parsing** de la commande `xuple-space`. Les étapes suivantes nécessaires pour rendre le système fonctionnel sont :

1. **Intégration Compilateur** : Ajouter la gestion des xuple-spaces dans le contexte de compilation
2. **Validation Unicité** : Vérifier l'unicité des noms de xuple-spaces à la compilation
3. **Création Runtime** : Instancier les xuple-spaces déclarés lors de l'exécution
4. **Actions Par Défaut** : Implémenter les actions `xuple:put`, `xuple:take`, `xuple:read`
5. **Intégration RETE** : Modifier le réseau RETE pour publier dans les xuple-spaces

Voir `/home/resinsec/dev/tsd/scripts/xuples/04-implement-default-actions.md` pour la suite.

## 📚 Références

- [Parser Analysis](../implementation/01-parser-analysis.md) - Analyse technique du parser
- [Syntax Specification](../implementation/02-xuplespace-syntax.md) - Spécification complète
- [Examples](../../examples/xuples/) - Exemples de code TSD
- [Xuples Module](../../xuples/) - Implémentation Go du module xuples

---

*Dernière mise à jour : 2025-12-17*
