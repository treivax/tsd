# Action Xuple - Réponse aux Questions

## ✅ L'action Xuple est-elle implémentée ?

**OUI, l'action Xuple est complètement implémentée et testée.**

## Comment fonctionne-t-elle ?

### Syntaxe
```tsd
Xuple(xuplespace: string, fact: any)
```

### Fonctionnement
1. **Validation** : Vérifie que le xuple-space existe et que le fait est valide
2. **Extraction automatique** : Récupère tous les faits déclencheurs du token de règle
3. **Création** : Crée un xuple avec ID unique (UUID), timestamp, et métadonnées
4. **Stockage** : Insère dans le xuple-space en appliquant les politiques configurées

### Exemple d'utilisation
```tsd
type Sensor(#id: string, location: string, temperature: number)
type Alert(#id: string, level: string, message: string)

xuple-space critical-alerts {
    selection: lifo
    consumption: per-agent
    retention: duration(10m)
}

rule critical_temp: {s: Sensor} / s.temperature > 40 ==>
    Xuple("critical-alerts", Alert(
        id: s.id + "_alert",
        level: "CRITICAL",
        message: "Temperature critical at " + s.location
    ))

Sensor(id: "S001", location: "Server-Room", temperature: 45.0)
```

## Quels sont les tests et exemples ?

### Tests Unitaires
**Fichier:** `rete/actions/builtin_test.go`
```bash
go test -v ./rete/actions -run TestExecuteXuple_InvalidArgs
```
✅ Teste la validation des arguments (7 cas)

### Tests d'Intégration
**Fichier:** `rete/actions/builtin_integration_test.go`
```bash
go test -v ./rete/actions -run TestBuiltinActions_EndToEnd_XupleAction
```
✅ Teste un scénario complet :
- Création de 2 xuple-spaces avec politiques différentes
- Création de 4 xuples via l'action
- Vérification du contenu avec `ListAll()`
- Test des politiques FIFO/LIFO
- Test des politiques once/per-agent
- Gestion d'erreurs

**Résultat:** PASS (100%)

### Exemples TSD
1. **`examples/xuples/xuple-action-example.tsd`**
   - Exemple complet avec sensors, alerts, commands
   - 5 règles utilisant Xuple
   - 3 xuple-spaces avec politiques différentes
   - Démonstration de cascades de règles

2. **`examples/xuples/basic-xuplespace.tsd`**
   - Exemple minimal

3. **`examples/xuples/all-policies.tsd`**
   - Démonstration de toutes les politiques possibles

## Comment valider son fonctionnement ?

### Méthode 1 : Tests automatisés
```bash
# Lancer le test d'intégration complet
go test -v ./rete/actions -run TestBuiltinActions_EndToEnd_XupleAction
```

### Méthode 2 : Inspection des xuples créés

#### API ajoutée : ListAll()
J'ai ajouté la méthode `ListAll()` à l'interface `XupleSpace` pour permettre l'inspection :

```go
// Obtenir un xuple-space
space, err := xupleManager.GetXupleSpace("critical-alerts")
if err != nil {
    log.Fatalf("Space not found: %v", err)
}

// Lister TOUS les xuples (pour tests/debug)
xuples := space.ListAll()
fmt.Printf("Total xuples créés: %d\n", len(xuples))

// Afficher les détails
for i, xuple := range xuples {
    fmt.Printf("Xuple %d:\n", i+1)
    fmt.Printf("  ID: %s\n", xuple.ID)
    fmt.Printf("  Type: %s\n", xuple.Fact.Type)
    fmt.Printf("  State: %s\n", xuple.Metadata.State)
    fmt.Printf("  Créé le: %s\n", xuple.CreatedAt)
    fmt.Printf("  Expire le: %s\n", xuple.Metadata.ExpiresAt)
    fmt.Printf("  Faits déclencheurs: %d\n", len(xuple.TriggeringFacts))
    
    // Afficher les champs du fait
    for key, value := range xuple.Fact.Fields {
        fmt.Printf("    %s: %v\n", key, value)
    }
}

// Compter les xuples DISPONIBLES (selon les politiques)
available := space.Count()
fmt.Printf("Xuples disponibles: %d\n", available)
```

### Méthode 3 : Exemple de test dans le test d'intégration

Extrait du test qui affiche les xuples :
```go
// Vérifier critical-alerts
criticalSpace, _ := xupleManager.GetXupleSpace("critical-alerts")
criticalXuples := criticalSpace.ListAll()

t.Logf("critical-alerts contient %d xuples", len(criticalXuples))
for i, xuple := range criticalXuples {
    t.Logf("   Xuple %d: ID=%s, Type=%s, State=%s",
        i+1, xuple.Fact.ID, xuple.Fact.Type, xuple.Metadata.State)
}
```

**Sortie du test :**
```
✅ critical-alerts contient 2 xuples
   Xuple 1: ID=A001, Type=Alert, State=available
   Xuple 2: ID=A002, Type=Alert, State=available
✅ command-queue contient 2 xuples
   Xuple 1: ID=C001, Type=Command, Action=activate_cooling, Priority=10
   Xuple 2: ID=C002, Type=Command, Action=send_notification, Priority=5
```

### Méthode 4 : Statistiques globales

```go
// Lister tous les xuple-spaces
spaces := xupleManager.ListXupleSpaces()
fmt.Printf("Total xuple-spaces: %d\n", len(spaces))

for _, name := range spaces {
    space, _ := xupleManager.GetXupleSpace(name)
    all := space.ListAll()
    available := space.Count()
    
    fmt.Printf("%s:\n", name)
    fmt.Printf("  Total: %d xuples\n", len(all))
    fmt.Printf("  Disponibles: %d xuples\n", available)
    
    // Nettoyer les expirés
    cleaned := space.Cleanup()
    if cleaned > 0 {
        fmt.Printf("  Nettoyés: %d xuples expirés\n", cleaned)
    }
}
```

## Structure d'un Xuple

Chaque xuple créé contient :

```go
type Xuple struct {
    ID              string        // UUID unique généré automatiquement
    Fact            *rete.Fact    // Le fait passé à l'action Xuple
    TriggeringFacts []*rete.Fact  // TOUS les faits qui ont déclenché la règle
    CreatedAt       time.Time     // Timestamp de création
    Metadata        XupleMetadata // État et métadonnées
}

type XupleMetadata struct {
    State            string                  // "available", "consumed", "expired"
    ConsumedBy       map[string]time.Time    // agentID -> timestamp de consommation
    ConsumptionCount int                     // Nombre total de consommations
    ExpiresAt        time.Time               // Date d'expiration (calculée par RetentionPolicy)
}
```

## Documentation complète

- **Guide utilisateur :** `docs/ACTION_XUPLE_GUIDE.md` (386 lignes)
- **Rapport d'implémentation :** `docs/XUPLE_ACTION_IMPLEMENTATION.md` (509 lignes)
- **Exemples TSD :** `examples/xuples/xuple-action-example.tsd`

## Résumé

✅ **Implémentation :** Complète et fonctionnelle  
✅ **Tests :** Unitaires + Intégration (100% pass)  
✅ **Exemples :** 3 fichiers TSD avec cas d'usage variés  
✅ **Inspection :** Méthode `ListAll()` ajoutée pour afficher les xuples  
✅ **Documentation :** Guide complet + rapport technique  
✅ **Traçabilité :** Les faits déclencheurs sont automatiquement conservés  
✅ **Statut :** Production-ready