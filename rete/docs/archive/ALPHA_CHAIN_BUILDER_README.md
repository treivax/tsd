# Alpha Chain Builder - Documentation

## Vue d'ensemble

Le **Alpha Chain Builder** est un composant du réseau RETE qui construit automatiquement des chaînes d'AlphaNodes avec **partage intelligent** entre règles. Il optimise l'utilisation de la mémoire et les performances en réutilisant les nœuds alpha identiques entre plusieurs règles.

## Architecture

### Composants principaux

#### 1. `AlphaChain`

Représente une chaîne complète d'AlphaNodes construite pour un ensemble de conditions.

```go
type AlphaChain struct {
    Nodes     []*AlphaNode  // Liste ordonnée des nœuds alpha
    Hashes    []string      // Hashes correspondants pour chaque nœud
    FinalNode *AlphaNode    // Le dernier nœud de la chaîne
    RuleID    string        // ID de la règle
}
```

#### 2. `AlphaChainBuilder`

Le constructeur qui gère la création et le partage des chaînes.

```go
type AlphaChainBuilder struct {
    network *ReteNetwork
    storage Storage
}
```

## Fonctionnalités

### 1. Construction de chaînes avec partage automatique

Le builder construit des chaînes d'AlphaNodes en réutilisant automatiquement les nœuds existants lorsque des conditions identiques sont rencontrées.

**Avantages :**
- 🔄 **Réutilisation automatique** : partage transparent des nœuds entre règles
- 💾 **Optimisation mémoire** : évite la duplication de nœuds identiques
- ⚡ **Performances améliorées** : moins de nœuds = évaluation plus rapide
- 📊 **Tracking du cycle de vie** : gestion automatique des références

### 2. Partage partiel et complet

Le builder supporte trois modes de partage :

#### Partage complet
Toutes les conditions sont identiques → tous les nœuds sont partagés

```
Rule1: age > 18 AND name == "Alice"
Rule2: age > 18 AND name == "Alice"
→ Les 2 nœuds sont partagés
```

#### Partage partiel
Certaines conditions sont identiques → partage des préfixes communs

```
Rule1: age > 18 AND name == "Alice"
Rule2: age > 18 AND city == "Paris"
→ Le 1er nœud est partagé, le 2ème est distinct
```

#### Pas de partage
Conditions complètement différentes → nouveaux nœuds créés

```
Rule1: age > 18
Rule2: salary > 50000
→ Aucun nœud partagé
```

## Utilisation

### 1. Création du builder

```go
// Créer un réseau RETE
storage := NewMemoryStorage()
network := NewReteNetwork(storage)

// Créer le builder
builder := NewAlphaChainBuilder(network, storage)
```

### 2. Construction d'une chaîne

```go
// Définir les conditions (normalisées)
conditions := []SimpleCondition{
    NewSimpleCondition("comparison", "p.age", ">", 18),
    NewSimpleCondition("comparison", "p.name", "==", "Alice"),
}

// Définir le nœud parent (TypeNode)
typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
parentNode := NewTypeNode("person", typeDef, storage)
network.TypeNodes["Person"] = parentNode

// Construire la chaîne
chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
if err != nil {
    log.Fatalf("Erreur: %v", err)
}

fmt.Printf("Chaîne construite avec %d nœuds\n", len(chain.Nodes))
```

### 3. Validation de la chaîne

```go
// Valider la cohérence de la chaîne
if err := chain.ValidateChain(); err != nil {
    log.Fatalf("Chaîne invalide: %v", err)
}
```

### 4. Récupération d'informations

```go
// Informations basiques
info := chain.GetChainInfo()
fmt.Printf("Rule ID: %s\n", info["rule_id"])
fmt.Printf("Nœuds: %d\n", info["node_count"])
fmt.Printf("Final node: %s\n", info["final_node_id"])

// Statistiques détaillées
stats := builder.GetChainStats(chain)
fmt.Printf("Total nodes: %d\n", stats["total_nodes"])
fmt.Printf("Shared nodes: %d\n", stats["shared_nodes"])
fmt.Printf("New nodes: %d\n", stats["new_nodes"])

// Détails par nœud
nodeDetails := stats["node_details"].([]map[string]interface{})
for _, detail := range nodeDetails {
    fmt.Printf("  Node %d: %s (refs=%d, shared=%v)\n",
        detail["index"],
        detail["node_id"],
        detail["ref_count"],
        detail["is_shared"])
}
```

### 5. Comptage des nœuds partagés

```go
sharedCount := builder.CountSharedNodes(chain)
fmt.Printf("Nœuds partagés: %d/%d\n", sharedCount, len(chain.Nodes))
```

## Algorithme de construction

Le builder suit cet algorithme pour chaque condition :

1. **Hash de la condition** : calculer un hash unique via `AlphaSharingManager`
2. **Recherche** : vérifier si un nœud avec ce hash existe déjà
3. **Si existant** :
   - Réutiliser le nœud
   - Vérifier la connexion au parent
   - Ajouter la référence de règle
   - Logger la réutilisation
4. **Si nouveau** :
   - Créer un nouveau AlphaNode
   - Connecter au parent
   - Ajouter au réseau
   - Enregistrer dans le LifecycleManager
   - Logger la création
5. **Suivant** : le nœud actuel devient parent pour le prochain

## Logging

Le builder fournit un logging détaillé pour le debugging :

```
🆕 [AlphaChainBuilder] Nouveau nœud alpha alpha_053115d3 créé pour rule1 (1/2)
🔗 [AlphaChainBuilder] Connexion du nœud alpha_053115d3 au parent type_person
♻️  [AlphaChainBuilder] Réutilisation du nœud alpha alpha_053115d3 pour rule2 (1/2)
✓  [AlphaChainBuilder] Nœud alpha_053115d3 déjà connecté au parent type_person
📊 [AlphaChainBuilder] Nœud alpha_053115d3 maintenant utilisé par 2 règle(s)
✅ [AlphaChainBuilder] Chaîne alpha complète construite pour rule1: 2 nœud(s)
```

## Intégration avec le réseau RETE

### Dépendances requises

Le builder nécessite que le `ReteNetwork` ait :

1. **AlphaSharingManager** : pour le partage automatique des nœuds
2. **LifecycleManager** : pour le tracking des références de règles
3. **Storage** : pour la persistance des nœuds

Ces composants sont automatiquement initialisés par `NewReteNetwork()`.

### Cycle de vie des nœuds

- **Création** : enregistré dans `network.AlphaNodes` et `LifecycleManager`
- **Partage** : incrémentation du compteur de références
- **Suppression** : décrémentation lors de `RemoveRule()`, suppression si ref_count == 0

## Tests

### Tests unitaires disponibles

- `TestBuildChain_SingleCondition` : chaîne d'un seul nœud
- `TestBuildChain_TwoConditions_New` : deux conditions nouvelles
- `TestBuildChain_TwoConditions_Reuse` : réutilisation complète
- `TestBuildChain_PartialReuse` : partage partiel
- `TestBuildChain_CompleteReuse` : partage par 5 règles
- `TestBuildChain_MultipleRules_SharedSubchain` : sous-chaînes partagées
- `TestBuildChain_EmptyConditions` : validation d'erreur
- `TestBuildChain_NilParent` : validation d'erreur
- `TestAlphaChain_ValidateChain` : validation de chaîne
- `TestAlphaChain_GetChainInfo` : extraction d'informations
- `TestAlphaChainBuilder_CountSharedNodes` : comptage de partage
- `TestAlphaChainBuilder_GetChainStats` : statistiques détaillées
- `TestIsAlreadyConnected` : helper de connexion

### Exécuter les tests

```bash
cd tsd/rete
go test -v -run TestBuildChain
go test -v -run TestAlphaChain
```

## Exemple complet

```go
package main

import (
    "fmt"
    "log"
    "github.com/treivax/tsd/rete"
)

func main() {
    // 1. Initialiser le réseau
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    
    // 2. Créer le type de données
    typeDef := rete.TypeDefinition{
        Type: "type",
        Name: "Person",
        Fields: []rete.Field{
            {Name: "age", Type: "number"},
            {Name: "name", Type: "string"},
            {Name: "city", Type: "string"},
        },
    }
    parentNode := rete.NewTypeNode("person", typeDef, storage)
    network.TypeNodes["Person"] = parentNode
    
    // 3. Créer le builder
    builder := rete.NewAlphaChainBuilder(network, storage)
    
    // 4. Définir les règles
    rule1Conditions := []rete.SimpleCondition{
        rete.NewSimpleCondition("comparison", "p.age", ">", 18),
        rete.NewSimpleCondition("comparison", "p.name", "==", "Alice"),
    }
    
    rule2Conditions := []rete.SimpleCondition{
        rete.NewSimpleCondition("comparison", "p.age", ">", 18),
        rete.NewSimpleCondition("comparison", "p.city", "==", "Paris"),
    }
    
    // 5. Construire les chaînes
    chain1, err := builder.BuildChain(rule1Conditions, "p", parentNode, "rule1")
    if err != nil {
        log.Fatal(err)
    }
    
    chain2, err := builder.BuildChain(rule2Conditions, "p", parentNode, "rule2")
    if err != nil {
        log.Fatal(err)
    }
    
    // 6. Analyser les statistiques
    fmt.Println("=== Rule 1 ===")
    stats1 := builder.GetChainStats(chain1)
    fmt.Printf("Nodes: %d, Shared: %d, New: %d\n",
        stats1["total_nodes"], stats1["shared_nodes"], stats1["new_nodes"])
    
    fmt.Println("\n=== Rule 2 ===")
    stats2 := builder.GetChainStats(chain2)
    fmt.Printf("Nodes: %d, Shared: %d, New: %d\n",
        stats2["total_nodes"], stats2["shared_nodes"], stats2["new_nodes"])
    
    // 7. Statistiques réseau globales
    netStats := network.GetNetworkStats()
    fmt.Printf("\n=== Network Stats ===\n")
    fmt.Printf("Total alpha nodes: %d\n", netStats["alpha_nodes"])
    fmt.Printf("Shared alpha nodes: %d\n", netStats["sharing_total_shared_alpha_nodes"])
}
```

**Sortie attendue :**
```
=== Rule 1 ===
Nodes: 2, Shared: 0, New: 2

=== Rule 2 ===
Nodes: 2, Shared: 1, New: 1

=== Network Stats ===
Total alpha nodes: 3
Shared alpha nodes: 3
```

## Bonnes pratiques

### 1. Normalisation des conditions

Toujours normaliser les conditions avant de construire la chaîne :

```go
// Normaliser l'expression
normalized, err := rete.NormalizeExpression(expression)
if err != nil {
    return err
}

// Extraire les conditions normalisées
conditions, operator, err := rete.ExtractConditions(normalized)
if err != nil {
    return err
}

// Construire la chaîne
chain, err := builder.BuildChain(conditions, variableName, parentNode, ruleID)
```

### 2. Gestion des erreurs

```go
chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
if err != nil {
    // Erreurs possibles :
    // - Liste de conditions vide
    // - Parent nil
    // - AlphaSharingManager non initialisé
    // - LifecycleManager non initialisé
    log.Fatalf("Échec construction chaîne: %v", err)
}
```

### 3. Validation post-construction

```go
if err := chain.ValidateChain(); err != nil {
    log.Fatalf("Chaîne invalide: %v", err)
}
```

### 4. Monitoring des performances

```go
// Avant d'ajouter des règles
initialStats := network.GetNetworkStats()
initialNodes := initialStats["alpha_nodes"].(int)

// Ajouter règles...

// Après
finalStats := network.GetNetworkStats()
finalNodes := finalStats["alpha_nodes"].(int)
nodesAdded := finalNodes - initialNodes

fmt.Printf("Nœuds ajoutés: %d (économie: %.1f%%)\n",
    nodesAdded,
    100.0 * (1.0 - float64(nodesAdded)/float64(totalConditions)))
```

## Compatibilité

- **License** : MIT
- **Go version** : 1.18+
- **Dépendances** : Package `tsd/rete` uniquement

## Limitations et évolutions futures

### Limitations actuelles

- Le builder ne gère que les conditions alpha (mono-fait)
- Pas de support pour la normalisation automatique (doit être faite en amont)
- Pas de cache de hashes des conditions (recalcul à chaque fois)

### Évolutions prévues

- [ ] Support des expressions complexes avec normalization intégrée
- [ ] Cache des hashes de conditions pour performances
- [ ] Métriques détaillées de partage par règle
- [ ] Visualisation graphique des chaînes partagées
- [ ] Support du partage inter-variables (avec vérification de compatibilité)

## Support et contributions

Pour toute question ou contribution :
- Voir `TEST_README.md` pour les instructions de test
- Voir `NORMALIZATION_README.md` pour la normalisation des conditions
- Voir `ALPHA_NODE_SHARING.md` pour les détails du mécanisme de partage

---

**Copyright (c) 2025 TSD Contributors**  
**Licensed under the MIT License**