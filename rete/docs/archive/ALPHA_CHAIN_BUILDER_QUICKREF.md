# Alpha Chain Builder - Référence rapide

## 🚀 Utilisation en 3 étapes

```go
// 1. Créer le builder
builder := rete.NewAlphaChainBuilder(network, storage)

// 2. Définir les conditions
conditions := []rete.SimpleCondition{
    rete.NewSimpleCondition("comparison", "p.age", ">", 18),
    rete.NewSimpleCondition("comparison", "p.name", "==", "Alice"),
}

// 3. Construire la chaîne
chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
```

## 📊 Statistiques

```go
// Informations basiques
info := chain.GetChainInfo()
fmt.Printf("Nœuds: %d\n", info["node_count"])

// Statistiques détaillées
stats := builder.GetChainStats(chain)
fmt.Printf("Partagés: %d/%d\n", stats["shared_nodes"], stats["total_nodes"])

// Compter les nœuds partagés
sharedCount := builder.CountSharedNodes(chain)
```

## ✅ Validation

```go
if err := chain.ValidateChain(); err != nil {
    log.Fatalf("Chaîne invalide: %v", err)
}
```

## 🔑 Concepts clés

| Concept | Description |
|---------|-------------|
| **Partage automatique** | Nœuds identiques réutilisés entre règles |
| **Partage partiel** | Préfixes communs partagés |
| **AlphaChain** | Structure représentant la chaîne complète |
| **Hash** | Identifiant unique pour chaque condition |
| **LifecycleManager** | Tracking des références de règles |

## 📝 Structures

```go
type AlphaChain struct {
    Nodes     []*AlphaNode  // Liste des nœuds
    Hashes    []string      // Hashes des conditions
    FinalNode *AlphaNode    // Dernier nœud
    RuleID    string        // ID de la règle
}

type AlphaChainBuilder struct {
    network *ReteNetwork
    storage Storage
}
```

## 🎯 Exemple complet

```go
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)
builder := rete.NewAlphaChainBuilder(network, storage)

typeDef := rete.TypeDefinition{
    Type: "type",
    Name: "Person",
    Fields: []rete.Field{
        {Name: "age", Type: "number"},
        {Name: "name", Type: "string"},
    },
}
parentNode := rete.NewTypeNode("person", typeDef, storage)

conditions := []rete.SimpleCondition{
    rete.NewSimpleCondition("comparison", "p.age", ">", 18),
    rete.NewSimpleCondition("comparison", "p.name", "==", "Alice"),
}

chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
if err != nil {
    log.Fatal(err)
}

stats := builder.GetChainStats(chain)
fmt.Printf("Chaîne construite: %d nœuds (%d partagés)\n",
    stats["total_nodes"], stats["shared_nodes"])
```

## 🔍 Messages de log

| Emoji | Message | Signification |
|-------|---------|---------------|
| 🆕 | Nouveau nœud créé | Nœud alpha créé pour la première fois |
| ♻️ | Réutilisation | Nœud existant réutilisé |
| 🔗 | Connexion | Nœud connecté à son parent |
| ✓ | Déjà connecté | Connexion existante vérifiée |
| 📊 | Statistiques | Compteur de références mis à jour |
| ✅ | Chaîne complète | Construction terminée avec succès |

## 📈 Performance

**Économie typique** : 30-50% de mémoire sur des règles similaires

**Exemple** :
- Sans partage : 7 nœuds
- Avec partage : 4 nœuds
- **Économie : 42.9%**

## 🧪 Tests

```bash
# Tous les tests du builder
go test -v -run TestBuildChain

# Tests de statistiques
go test -v -run TestAlphaChain

# Exécuter l'exemple
go run examples/alpha_chain_builder_example.go
```

## ⚠️ Erreurs courantes

| Erreur | Cause | Solution |
|--------|-------|----------|
| "liste vide" | Aucune condition | Passer au moins 1 condition |
| "parent nil" | Parent non initialisé | Créer le TypeNode avant |
| "AlphaSharingManager non initialisé" | Réseau mal créé | Utiliser NewReteNetwork() |
| "LifecycleManager non initialisé" | Réseau mal créé | Utiliser NewReteNetwork() |

## 📚 Documentation complète

- **README** : Guide d'utilisation détaillé → `ALPHA_CHAIN_BUILDER_README.md`
- **SUMMARY** : Résumé technique → `ALPHA_CHAIN_BUILDER_SUMMARY.md`
- **CHANGELOG** : Historique → `ALPHA_CHAIN_BUILDER_CHANGELOG.md`
- **INDEX** : Vue d'ensemble → `ALPHA_CHAIN_BUILDER_INDEX.md`
- **Exemple** : Code complet → `examples/alpha_chain_builder_example.go`

## 🔗 Workflow d'intégration

```
Expression → Normalize → Extract → BuildChain → AlphaChain
             (normalization) (extractor) (builder)   (partagée)
```

## ✅ Checklist avant utilisation

- [ ] Réseau RETE créé avec `NewReteNetwork()`
- [ ] Storage initialisé (MemoryStorage ou autre)
- [ ] TypeNode créé et ajouté au réseau
- [ ] Conditions normalisées (via `NormalizeExpression`)
- [ ] Variable name définie (ex: "p", "order", etc.)

## 🎓 Bonnes pratiques

1. ✅ **Toujours normaliser** les conditions avant BuildChain
2. ✅ **Valider** la chaîne après construction
3. ✅ **Monitorer** les statistiques de partage
4. ✅ **Réutiliser** le même builder pour toutes les règles
5. ✅ **Logger** pour debug (logs détaillés fournis)

## 💡 Tips

- Le partage fonctionne sur des conditions **exactement identiques**
- Les nœuds sont identifiés par **hash SHA-256** des conditions
- Le **premier nœud** d'une chaîne a souvent le plus de partage
- Utiliser `GetChainStats()` pour optimiser les règles
- Les **compteurs de références** permettent la suppression sûre

## 🚦 Status

- **Version** : 1.0.0
- **Status** : ✅ Production Ready
- **Tests** : 13/13 (100%)
- **License** : MIT

---

**Copyright (c) 2025 TSD Contributors**  
**Licensed under the MIT License**