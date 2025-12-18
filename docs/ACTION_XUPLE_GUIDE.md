# Guide de l'Action Xuple

## Vue d'ensemble

L'action `Xuple` est une action prédéfinie de TSD qui permet de créer des **xuples** dans des **xuple-spaces** depuis des règles. Un xuple est un tuple enrichi qui combine :
- Un fait résultant d'une activation de règle
- Les faits déclencheurs qui ont causé cette activation
- Des métadonnées (timestamp, état, politiques)

Les xuple-spaces sont des espaces de coordination inspirés des **tuple spaces** de Linda, permettant une communication asynchrone et découplée entre agents.

## Syntaxe

```tsd
Xuple(xuplespace: string, fact: any)
```

### Paramètres

- **xuplespace** (string) : Nom du xuple-space cible (doit être déclaré au préalable)
- **fact** (any) : Fait à insérer dans le xuple-space

## Déclaration de xuple-spaces

Avant d'utiliser l'action `Xuple`, vous devez déclarer les xuple-spaces avec leurs politiques :

```tsd
xuple-space <name> {
    selection: <fifo|lifo|random>
    consumption: <once|per-agent|limited(N)>
    retention: <unlimited|duration(Xs|Xm|Xh|Xd)>
}
```

### Politiques disponibles

#### Selection Policy
- **fifo** : Premier arrivé, premier servi (queue)
- **lifo** : Dernier arrivé, premier servi (stack)
- **random** : Sélection aléatoire (load balancing)

#### Consumption Policy
- **once** : Consommé une seule fois globalement
- **per-agent** : Chaque agent peut consommer une fois (broadcast)
- **limited(N)** : Consommable N fois maximum

#### Retention Policy
- **unlimited** : Conservé indéfiniment
- **duration(Xs)** : Expire après X secondes
- **duration(Xm)** : Expire après X minutes
- **duration(Xh)** : Expire après X heures
- **duration(Xd)** : Expire après X jours

## Exemple complet

```tsd
// Types
type Sensor(#id: string, location: string, temperature: number)
type Alert(#id: string, level: string, message: string, sensorId: string)
type Command(#id: string, action: string, target: string, priority: number)

// Déclaration des xuple-spaces
xuple-space critical-alerts {
    selection: lifo
    consumption: per-agent
    retention: duration(10m)
}

xuple-space command-queue {
    selection: fifo
    consumption: once
    retention: duration(1h)
}

// Règles utilisant Xuple
rule critical_temperature: {s: Sensor} / s.temperature > 40 ==>
    Xuple("critical-alerts", Alert(
        id: s.id + "_alert",
        level: "CRITICAL",
        message: "Temperature critical at " + s.location,
        sensorId: s.id
    ))

rule alert_to_command: {a: Alert} / a.level == "CRITICAL" ==>
    Xuple("command-queue", Command(
        id: a.sensorId + "_cmd",
        action: "activate_cooling",
        target: a.sensorId,
        priority: 10
    ))

// Faits déclencheurs
Sensor(id: "S001", location: "Server-Room", temperature: 45.0)
```

## Fonctionnement interne

### 1. Validation
- Vérifie que le xuple-space existe
- Vérifie que le fait est valide

### 2. Extraction du contexte
L'action `Xuple` extrait automatiquement tous les faits déclencheurs du token de règle, préservant ainsi la **traçabilité causale**.

### 3. Création du xuple
Un xuple est créé avec :
- **ID unique** : Généré automatiquement (UUID)
- **Fact** : Le fait passé en paramètre
- **TriggeringFacts** : Tous les faits qui ont déclenché la règle
- **CreatedAt** : Timestamp de création
- **Metadata** : État, consommations, expiration

### 4. Application des politiques
Le xuple-space applique ses politiques :
- **Rétention** : Calcul de `ExpiresAt`
- **Capacité** : Vérification de `MaxSize` (si défini)

### 5. Disponibilité
Le xuple devient immédiatement disponible pour récupération par les agents via `Retrieve()`.

## Structure d'un Xuple

```go
type Xuple struct {
    ID              string        // UUID unique
    Fact            *Fact         // Fait principal
    TriggeringFacts []*Fact       // Faits déclencheurs
    CreatedAt       time.Time     // Timestamp de création
    Metadata        XupleMetadata // État et métadonnées
}

type XupleMetadata struct {
    State            string                  // available, consumed, expired
    ConsumedBy       map[string]time.Time    // agentID -> timestamp
    ConsumptionCount int                     // Nombre de consommations
    ExpiresAt        time.Time               // Date d'expiration
}
```

## Validation du fonctionnement

### Méthode 1 : Tests unitaires

```go
import (
    "testing"
    "time"
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/xuples"
)

func TestXupleCreation(t *testing.T) {
    // Setup
    xupleManager := xuples.NewXupleManager()
    config := xuples.XupleSpaceConfig{
        Name:              "test-space",
        SelectionPolicy:   xuples.NewFIFOSelectionPolicy(),
        ConsumptionPolicy: xuples.NewOnceConsumptionPolicy(),
        RetentionPolicy:   xuples.NewDurationRetentionPolicy(10 * time.Minute),
    }
    xupleManager.CreateXupleSpace("test-space", config)
    
    // Créer un xuple via l'action
    executor := NewBuiltinActionExecutor(network, xupleManager, nil, nil)
    fact := &rete.Fact{ID: "F001", Type: "Test"}
    token := &rete.Token{Facts: []*rete.Fact{fact}}
    
    err := executor.Execute("Xuple", []interface{}{"test-space", fact}, token)
    if err != nil {
        t.Fatalf("Failed to create xuple: %v", err)
    }
    
    // Vérifier la création
    space, _ := xupleManager.GetXupleSpace("test-space")
    xuples := space.ListAll()
    
    if len(xuples) != 1 {
        t.Errorf("Expected 1 xuple, got %d", len(xuples))
    }
    
    xuple := xuples[0]
    t.Logf("Xuple créé:")
    t.Logf("  ID: %s", xuple.ID)
    t.Logf("  Type: %s", xuple.Fact.Type)
    t.Logf("  State: %s", xuple.Metadata.State)
    t.Logf("  TriggeringFacts: %d", len(xuple.TriggeringFacts))
}
```

### Méthode 2 : Inspection via l'API

```go
// Obtenir un xuple-space
space, err := xupleManager.GetXupleSpace("critical-alerts")
if err != nil {
    log.Fatalf("Space not found: %v", err)
}

// Lister tous les xuples (pour debug/test)
xuples := space.ListAll()
fmt.Printf("Total xuples: %d\n", len(xuples))

for i, xuple := range xuples {
    fmt.Printf("Xuple %d:\n", i+1)
    fmt.Printf("  ID: %s\n", xuple.ID)
    fmt.Printf("  Type: %s\n", xuple.Fact.Type)
    fmt.Printf("  State: %s\n", xuple.Metadata.State)
    fmt.Printf("  Created: %s\n", xuple.CreatedAt)
    fmt.Printf("  Expires: %s\n", xuple.Metadata.ExpiresAt)
    fmt.Printf("  Triggering facts: %d\n", len(xuple.TriggeringFacts))
    fmt.Printf("  Consumed by: %d agents\n", len(xuple.Metadata.ConsumedBy))
}

// Compter les xuples disponibles
available := space.Count()
fmt.Printf("Available xuples: %d\n", available)
```

### Méthode 3 : Récupération avec politiques

```go
// Récupérer un xuple selon les politiques
xuple, err := space.Retrieve("agent1")
if err != nil {
    log.Printf("No xuple available: %v", err)
} else {
    fmt.Printf("Retrieved xuple: %s (Type: %s)\n", xuple.ID, xuple.Fact.Type)
    
    // Marquer comme consommé
    err = space.MarkConsumed(xuple.ID, "agent1")
    if err != nil {
        log.Printf("Failed to mark consumed: %v", err)
    }
}
```

## Cas d'usage

### 1. Alertes critiques (LIFO + per-agent)
Traiter les alertes les plus récentes en priorité, chaque agent doit les voir.

```tsd
xuple-space critical-alerts {
    selection: lifo
    consumption: per-agent
    retention: duration(10m)
}
```

### 2. File de commandes (FIFO + once)
Traiter les commandes dans l'ordre d'arrivée, chaque commande exécutée une fois.

```tsd
xuple-space command-queue {
    selection: fifo
    consumption: once
    retention: duration(1h)
}
```

### 3. Load balancing (Random + once)
Distribution aléatoire des tâches entre agents.

```tsd
xuple-space task-pool {
    selection: random
    consumption: once
    retention: unlimited
}
```

### 4. Cache distribué (Random + limited)
Données partagées avec limite de lecture.

```tsd
xuple-space cache {
    selection: random
    consumption: limited(100)
    retention: duration(5m)
}
```

### 5. Publish-Subscribe (Random + per-agent)
Diffusion d'événements à tous les agents.

```tsd
xuple-space events {
    selection: random
    consumption: per-agent
    retention: duration(15m)
}
```

## Gestion des erreurs

L'action `Xuple` retourne une erreur dans les cas suivants :

### Erreur : Xuple-space inexistant
```
Error: xuple-space not found
```
**Solution** : Déclarer le xuple-space avant de l'utiliser.

### Erreur : Xuple-space plein
```
Error: xuple-space full
```
**Solution** : Augmenter `MaxSize` ou nettoyer les xuples expirés.

### Erreur : Fait invalide
```
Error: fact is nil
```
**Solution** : Vérifier que le fait passé est valide.

### Erreur : XupleManager non configuré
```
Error: action Xuple requires XupleManager to be configured
```
**Solution** : Initialiser le XupleManager lors de la création de l'executor.

## Tests d'intégration

Le fichier `rete/actions/builtin_integration_test.go` contient un test complet :

```bash
go test -v ./rete/actions -run TestBuiltinActions_EndToEnd_XupleAction
```

Ce test vérifie :
- ✅ Création de xuple-spaces avec différentes politiques
- ✅ Création de xuples via l'action Xuple
- ✅ Extraction des faits déclencheurs
- ✅ Application des politiques de sélection (FIFO, LIFO)
- ✅ Application des politiques de consommation (once, per-agent)
- ✅ Inspection du contenu des xuple-spaces
- ✅ Récupération avec politiques
- ✅ Gestion des erreurs

## Exemples complets

Voir les fichiers d'exemple :
- `examples/xuples/xuple-action-example.tsd` : Exemple complet avec sensors/alerts/commands
- `examples/xuples/basic-xuplespace.tsd` : Exemple basique
- `examples/xuples/all-policies.tsd` : Démonstration de toutes les politiques

## Références

- [Documentation Xuples](../xuples/README.md)
- [Actions par défaut](ACTIONS_PAR_DEFAUT_SYNTHESE.md)
- [Tests unitaires](../rete/actions/builtin_test.go)
- [Tests d'intégration](../rete/actions/builtin_integration_test.go)

## Métriques et performance

Pour des performances optimales :
- Utiliser `duration` pour éviter l'accumulation de xuples expirés
- Limiter `MaxSize` pour éviter la croissance mémoire
- Nettoyer périodiquement avec `Cleanup()`
- Utiliser `Count()` pour monitorer la taille des xuple-spaces

```go
// Monitoring exemple
spaces := xupleManager.ListXupleSpaces()
for _, name := range spaces {
    space, _ := xupleManager.GetXupleSpace(name)
    count := space.Count()
    log.Printf("Space %s: %d available xuples", name, count)
    
    // Nettoyer les expirés
    cleaned := space.Cleanup()
    if cleaned > 0 {
        log.Printf("Cleaned %d expired xuples from %s", cleaned, name)
    }
}
```

## Conclusion

L'action `Xuple` offre un mécanisme puissant de coordination asynchrone entre règles et agents :
- **Découplage** : Les producteurs et consommateurs n'ont pas besoin de se connaître
- **Traçabilité** : Chaque xuple conserve ses faits déclencheurs
- **Flexibilité** : Politiques configurables pour différents cas d'usage
- **Robustesse** : Gestion automatique de l'expiration et de la consommation

Pour plus d'informations, consultez la documentation complète des xuples et les exemples fournis.