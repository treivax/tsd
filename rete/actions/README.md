# Package actions - Actions Natives TSD

## Vue d'ensemble

Ce package implémente les actions natives (built-in) du système TSD. Ces actions sont automatiquement disponibles dans toutes les règles TSD sans nécessiter de définition préalable.

## Actions Implémentées ✅

### Print(message: string)
Affiche un message sur la sortie standard.

```tsd
rule "example" {
    when { ... }
    then {
        Print("Hello, World!")
    }
}
```

**Statut**: ✅ Implémenté et testé
**Couverture**: 100%

### Log(message: string)
Génère une entrée dans le système de logging.

```tsd
rule "example" {
    when { ... }
    then {
        Log("Event processed")
    }
}
```

**Statut**: ✅ Implémenté et testé
**Couverture**: 100%

### Xuple(xuplespace: string, fact: any)
Crée un xuple dans le xuple-space spécifié.

```tsd
rule "example" {
    when {
        e: Event()
    }
    then {
        Xuple("events", e)
    }
}
```

**Statut**: ✅ Implémenté et testé
**Couverture**: 100%
**Prérequis**: Nécessite un XupleManager configuré

## Implémentation des Actions CRUD ✅

Les actions suivantes modifient dynamiquement les faits dans le réseau RETE et propagent automatiquement les changements.

### Update(fact: any)
**Statut**: ✅ IMPLÉMENTÉ

Modifie un fait existant et propage les changements dans le réseau RETE.

**Fonctionnement**:
1. Vérifie que le fait existe dans le réseau
2. Rétracte l'ancien fait (propage la suppression)
3. Insère le fait mis à jour (propage l'ajout)
4. Garantit la cohérence du réseau RETE

**Exemple**:
```tsd
rule update_age : {p: Person} / p.age == 30 ==> Update(Person(id: p.id, name: p.name, age: 31))
```

### Insert(fact: any)
**Statut**: ✅ IMPLÉMENTÉ

Crée et insère dynamiquement un nouveau fait dans le réseau RETE.

**Fonctionnement**:
1. Valide le fait (type, ID, champs)
2. Vérifie qu'il n'existe pas déjà
3. L'ajoute au storage
4. Propage l'insertion dans le réseau RETE

**Exemple**:
```tsd
rule create_admin : {u: User} / u.role == "manager" ==> Insert(Admin(id: u.id, level: "high"))
```

### Retract(id: string)
**Statut**: ✅ IMPLÉMENTÉ

Supprime dynamiquement un fait du réseau RETE et tous les tokens dépendants.

**Fonctionnement**:
1. Valide l'ID du fait
2. Vérifie que le fait existe
3. Le supprime du storage
4. Propage la rétraction dans le réseau RETE
5. Nettoie les références et tokens associés

**Exemple**:
```tsd
rule remove_inactive : {u: User} / u.active == false ==> Retract("User_" + u.id)
```

## Architecture

### BuiltinActionExecutor

Point central d'exécution de toutes les actions natives.

```go
executor := actions.NewBuiltinActionExecutor(
    network,      // *rete.ReteNetwork
    xupleManager, // xuples.XupleManager (peut être nil)
    output,       // io.Writer pour Print (os.Stdout par défaut)
    logger,       // *log.Logger pour Log (log.Default() par défaut)
)

err := executor.Execute(actionName, args, token)
```

### Thread-Safety

- ✅ Toutes les méthodes sont thread-safe
- ✅ Aucune variable globale mutable
- ✅ Synchronisation déléguée au ReteNetwork et XupleManager

## Tests

**Couverture actuelle**: 91.5%

Tests complets pour:
- ✅ Validation des arguments
- ✅ Cas nominaux et cas d'erreur
- ✅ Toutes les actions implémentées (Print, Log, Update, Insert, Retract, Xuple)
- ✅ Extraction des faits déclencheurs
- ✅ Configuration du output et logger
- ✅ Tests d'intégration avec le réseau RETE

## Utilisation

### Configuration Standard

```go
import (
    "github.com/treivax/tsd/rete"
    "github.com/treivax/tsd/rete/actions"
    "github.com/treivax/tsd/xuples"
)

// Créer le réseau RETE
network := rete.NewReteNetwork()

// Créer le gestionnaire de xuples (optionnel)
xupleManager := xuples.NewXupleManager()

// Créer l'exécuteur d'actions
executor := actions.NewBuiltinActionExecutor(
    network,
    xupleManager,
    nil, // Utilise os.Stdout
    nil, // Utilise log.Default()
)

// Utiliser dans les règles
err := executor.Execute("Print", []interface{}{"Hello"}, token)
```

### Tests Unitaires

```go
import (
    "bytes"
    "log"
    "testing"
)

func TestMyAction(t *testing.T) {
    output := &bytes.Buffer{}
    logOutput := &bytes.Buffer{}
    logger := log.New(logOutput, "", 0)
    
    executor := actions.NewBuiltinActionExecutor(
        network,
        nil,
        output,
        logger,
    )
    
    err := executor.Execute("Print", []interface{}{"test"}, nil)
    // Assertions...
}
```

## Feuille de Route

### Court terme (v1.1)
- [x] Implémenter Update() dans rete ✅
- [x] Implémenter Insert() dans rete ✅
- [x] Implémenter Retract() dans rete ✅

### Moyen terme (v1.2)
- [ ] Ajouter action Debug() pour inspection
- [ ] Ajouter action Assert() pour validation
- [ ] Support des actions asynchrones

### Long terme (v2.0)
- [ ] Actions personnalisées utilisateur
- [ ] Plugin system pour actions
- [ ] Actions distribuées

## Contribution

Pour ajouter une nouvelle action:

1. Ajouter la signature dans `internal/defaultactions/defaults.tsd`
2. Implémenter la méthode `execute*` dans `builtin.go`
3. Ajouter le case dans `Execute()` 
4. Implémenter les méthodes RETE nécessaires si besoin
5. Ajouter des tests complets
6. Mettre à jour cette documentation

## Références

- [defaults.tsd](../../internal/defaultactions/defaults.tsd) - Définitions TSD des actions
- [loader.go](../../internal/defaultactions/loader.go) - Chargement des définitions
- [Documentation xuples](../../xuples/README.md) - Pour l'action Xuple
