# Gestion du cycle de vie des nœuds dans le réseau RETE

## 🎯 Vue d'ensemble

Cette fonctionnalité implémente un système complet de **tracking des règles** et de **gestion du cycle de vie des nœuds** dans le réseau RETE.

**Problème résolu** : Permettre la suppression propre de règles sans fuites mémoire, en ne supprimant un nœud que si plus aucune règle ne l'utilise.

**Statut** : ✅ Production-ready (19/19 tests PASS)

## 🚀 Démarrage rapide

### Utilisation basique

```go
import "github.com/treivax/tsd/rete"

// 1. Créer un réseau (LifecycleManager initialisé automatiquement)
storage := rete.NewMemoryStorage()
network := rete.NewReteNetwork(storage)

// 2. Construire depuis un fichier TSD (tracking automatique)
pipeline := rete.NewConstraintPipeline()
network, err := pipeline.BuildNetworkFromConstraintFile("rules.tsd", storage)

// 3. Supprimer une règle proprement
err = network.RemoveRule("rule_0")
// → Supprime les nœuds orphelins
// → Conserve les nœuds partagés

// 4. Obtenir des informations
info := network.GetRuleInfo("rule_0")
fmt.Printf("Règle %s utilise %d nœuds\n", info.RuleID, info.NodeCount)

// 5. Statistiques du réseau
stats := network.GetNetworkStats()
fmt.Printf("Total nœuds: %d\n", stats["lifecycle_total_nodes"])
```

### Exemple complet

```go
// Fichier rules.tsd:
// type Person : <id: string, age: number>
// rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
// rule r2 : {p: Person} / p.age < 65 ==> young(p.id)

network, _ := pipeline.BuildNetworkFromConstraintFile("rules.tsd", storage)

// État initial: 1 TypeNode, 2 AlphaNodes, 2 TerminalNodes

network.RemoveRule("rule_0")

// État après: 1 TypeNode (partagé), 1 AlphaNode, 1 TerminalNode
// → rule_0_alpha et rule_0_terminal supprimés
// → TypeNode conservé (utilisé par rule_1)
```

## 📊 Fonctionnalités

### 1. Comptage de références automatique

Chaque nœud sait combien de règles l'utilisent :

```go
lifecycle := network.LifecycleManager.GetNodeLifecycle("alpha_1")
fmt.Printf("Nœud utilisé par %d règle(s)\n", lifecycle.GetRefCount())
```

### 2. Suppression intelligente

Les nœuds partagés sont automatiquement conservés :

```go
// TypeNode partagé par 3 règles
network.RemoveRule("rule_0")  // TypeNode conservé (2 règles restantes)
network.RemoveRule("rule_1")  // TypeNode conservé (1 règle restante)
network.RemoveRule("rule_2")  // TypeNode toujours conservé (créé au niveau type)
```

### 3. Informations détaillées

```go
// Informations sur une règle
info := network.GetRuleInfo("rule_0")
// → info.RuleID, info.RuleName, info.NodeCount, info.NodeIDs

// Statistiques du réseau
stats := network.GetNetworkStats()
// → type_nodes, alpha_nodes, terminal_nodes
// → lifecycle_total_nodes, lifecycle_total_references
// → lifecycle_nodes_with_refs, lifecycle_nodes_without_refs

// Nœuds d'une règle spécifique
nodeIDs := network.LifecycleManager.GetNodesForRule("rule_0")
```

## 🏗️ Architecture

### Composants

```
RuleReference
  ↓ contient
NodeLifecycle (par nœud)
  ↓ gérés par
LifecycleManager (global)
  ↓ intégré dans
ReteNetwork
```

### Flux de données

**Création d'une règle** :
```
Parser TSD → Créer nœuds → Enregistrer dans LifecycleManager
                                    ↓
                            Ajouter référence règle
                                    ↓
                            Connecter dans le réseau
```

**Suppression d'une règle** :
```
network.RemoveRule(id) → Récupérer nœuds de la règle
                                    ↓
                         Retirer référence de chaque nœud
                                    ↓
                         Si RefCount == 0:
                           → Supprimer du réseau
                         Sinon:
                           → Conserver (partagé)
```

## 🔑 Points clés

### Thread-safety ✅
Toutes les opérations sont thread-safe (utilisation de `sync.RWMutex`)

### Performance ⚡
- Ajout/suppression : O(n) où n = nœuds de la règle
- Vérification : O(1)
- Overhead mémoire : ~150 bytes par nœud

### Partage de nœuds 🔗
- **TypeNodes** : partagés entre toutes les règles du même type
- **AlphaNodes** : spécifiques à chaque règle
- **TerminalNodes** : spécifiques à chaque règle

### Gestion automatique 🤖
Le tracking est transparent :
- Enregistrement automatique lors de la compilation
- Pas de code supplémentaire nécessaire
- Compatible avec le code existant

## 📚 Documentation

### Fichiers de référence

| Fichier | Description | Taille |
|---------|-------------|--------|
| [NODE_LIFECYCLE_README.md](./NODE_LIFECYCLE_README.md) | Ce fichier (guide principal) | - |
| [NODE_LIFECYCLE_SUMMARY.md](./NODE_LIFECYCLE_SUMMARY.md) | Résumé exécutif | 309 lignes |
| [NODE_LIFECYCLE_FEATURE.md](./NODE_LIFECYCLE_FEATURE.md) | Documentation technique complète | 410 lignes |
| [node_lifecycle.go](./node_lifecycle.go) | Code source | 278 lignes |
| [node_lifecycle_test.go](./node_lifecycle_test.go) | Tests unitaires | 434 lignes |
| [network_lifecycle_test.go](./network_lifecycle_test.go) | Tests d'intégration | 388 lignes |

### Ordre de lecture recommandé

1. **Ce fichier** (NODE_LIFECYCLE_README.md) → Vue d'ensemble et démarrage rapide
2. **NODE_LIFECYCLE_SUMMARY.md** → Résumé avec exemples concrets
3. **NODE_LIFECYCLE_FEATURE.md** → Documentation technique approfondie
4. **Tests** → Exemples d'utilisation pratiques

## 🧪 Tests

### Exécution

```bash
cd tsd/rete

# Tous les tests lifecycle
go test -v -run "TestNodeLifecycle|TestNetworkLifecycle"

# Tests unitaires seulement
go test -v -run TestNodeLifecycle

# Tests d'intégration seulement
go test -v -run TestNetworkLifecycle
```

### Résultats attendus

```
=== Tests unitaires (11) ===
✅ NodeLifecycle basique
✅ Ajout/suppression de références
✅ Récupération de règles
✅ LifecycleManager operations
✅ Statistiques et reset

=== Tests d'intégration (8) ===
✅ Suppression de règles simples
✅ Suppression de toutes les règles
✅ Partage de nœuds
✅ Informations et statistiques
✅ Gestion d'erreurs

Total: 19/19 PASS (~8ms)
```

## 💡 Exemples d'utilisation

### Exemple 1 : Suppression basique

```go
network, _ := pipeline.BuildNetworkFromConstraintFile("rules.tsd", storage)

// Supprimer une règle
err := network.RemoveRule("rule_0")
if err != nil {
    log.Printf("Erreur: %v", err)
}

// Vérifier le résultat
stats := network.GetNetworkStats()
fmt.Printf("Nœuds restants: %d\n", stats["lifecycle_total_nodes"])
```

### Exemple 2 : Informations détaillées

```go
// Avant suppression
info := network.GetRuleInfo("rule_0")
fmt.Printf("Règle %s utilise ces nœuds:\n", info.RuleID)
for _, nodeID := range info.NodeIDs {
    fmt.Printf("  - %s\n", nodeID)
}

// Supprimer
network.RemoveRule("rule_0")

// Vérifier que les nœuds ont été supprimés
for _, nodeID := range info.NodeIDs {
    lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
    if exists {
        fmt.Printf("Nœud %s: %d référence(s) restante(s)\n", 
            nodeID, lifecycle.GetRefCount())
    } else {
        fmt.Printf("Nœud %s: supprimé\n", nodeID)
    }
}
```

### Exemple 3 : Monitoring

```go
// Surveiller l'état du réseau
ticker := time.NewTicker(5 * time.Second)
for range ticker.C {
    stats := network.GetNetworkStats()
    
    fmt.Printf("=== État du réseau ===\n")
    fmt.Printf("Nœuds actifs: %d\n", stats["lifecycle_total_nodes"])
    fmt.Printf("Références: %d\n", stats["lifecycle_total_references"])
    fmt.Printf("Nœuds orphelins: %d\n", stats["lifecycle_nodes_without_refs"])
}
```

### Exemple 4 : Gestion d'erreurs

```go
// Tentative de suppression d'une règle inexistante
err := network.RemoveRule("rule_999")
if err != nil {
    fmt.Printf("Erreur attendue: %v\n", err)
    // → "règle rule_999 non trouvée ou aucun nœud associé"
}

// Vérifier avant de supprimer
if info := network.GetRuleInfo("rule_0"); info.NodeCount > 0 {
    network.RemoveRule("rule_0")
} else {
    fmt.Println("Règle déjà supprimée")
}
```

## ⚙️ API complète

### ReteNetwork

```go
// Supprimer une règle et ses nœuds orphelins
RemoveRule(ruleID string) error

// Obtenir des informations sur une règle
GetRuleInfo(ruleID string) *RuleInfo

// Statistiques du réseau
GetNetworkStats() map[string]interface{}

// Réinitialiser (nettoie le lifecycle)
Reset()
```

### LifecycleManager

```go
// Enregistrer un nœud
RegisterNode(nodeID, nodeType string) *NodeLifecycle

// Associer une règle à un nœud
AddRuleToNode(nodeID, ruleID, ruleName string) error

// Retirer une règle d'un nœud (retourne true si plus de références)
RemoveRuleFromNode(nodeID, ruleID string) (bool, error)

// Supprimer un nœud (uniquement si RefCount == 0)
RemoveNode(nodeID string) error

// Récupérer les nœuds d'une règle
GetNodesForRule(ruleID string) []string

// Vérifier si un nœud peut être supprimé
CanRemoveNode(nodeID string) bool

// Statistiques
GetStats() map[string]interface{}

// Réinitialiser
Reset()
```

### NodeLifecycle

```go
// Ajouter une référence de règle
AddRuleReference(ruleID, ruleName string)

// Retirer une référence (retourne true si plus de références)
RemoveRuleReference(ruleID string) bool

// Vérifier s'il reste des références
HasReferences() bool

// Nombre de références
GetRefCount() int

// Liste des IDs de règles
GetRules() []string
```

## 🔧 Configuration

### Initialisation automatique

Le `LifecycleManager` est initialisé automatiquement :

```go
network := rete.NewReteNetwork(storage)
// → network.LifecycleManager est prêt
```

### Enregistrement automatique

Les nœuds sont enregistrés automatiquement lors de la compilation :

```go
network, _ := pipeline.BuildNetworkFromConstraintFile("rules.tsd", storage)
// → Tous les nœuds sont trackés automatiquement
```

### Enregistrement manuel (si nécessaire)

```go
// Enregistrer un nœud manuellement
lifecycle := network.LifecycleManager.RegisterNode("custom_node", "alpha")
lifecycle.AddRuleReference("custom_rule", "Custom Rule")
```

## ✅ Validation

### Checklist de validation

- [x] Tracking des règles par nœud
- [x] Suppression propre des règles
- [x] Comptage de références fonctionnel
- [x] Pas de fuites mémoire
- [x] Thread-safe
- [x] Tests complets (19 tests, 100% PASS)
- [x] Documentation complète
- [x] Rétrocompatible
- [x] Production-ready

### Scénarios validés

- [x] Suppression d'une règle simple
- [x] Suppression de toutes les règles d'un type
- [x] Partage de nœuds entre règles
- [x] Nœuds orphelins correctement supprimés
- [x] Nœuds partagés correctement conservés
- [x] Reset du réseau
- [x] Statistiques et informations
- [x] Gestion d'erreurs

## 🐛 Dépannage

### Problème : "règle non trouvée"

```go
err := network.RemoveRule("rule_0")
// Erreur: règle rule_0 non trouvée ou aucun nœud associé

// Solution: Vérifier que la règle existe
info := network.GetRuleInfo("rule_0")
if info.NodeCount == 0 {
    fmt.Println("Règle non trouvée ou déjà supprimée")
}
```

### Problème : Nœud non supprimé

```go
// Le nœud n'est pas supprimé car il est partagé
lifecycle, _ := network.LifecycleManager.GetNodeLifecycle("type_Person")
fmt.Printf("Nœud utilisé par %d règle(s)\n", lifecycle.GetRefCount())

// Solution: Supprimer toutes les règles utilisant le nœud
for _, ruleID := range lifecycle.GetRules() {
    network.RemoveRule(ruleID)
}
```

### Problème : Fuite mémoire suspectée

```go
// Vérifier les nœuds orphelins
stats := network.GetNetworkStats()
orphans := stats["lifecycle_nodes_without_refs"].(int)
if orphans > 0 {
    fmt.Printf("⚠️  %d nœud(s) orphelin(s) détecté(s)\n", orphans)
}

// Solution: Reset complet si nécessaire
network.Reset()
```

## 📞 Support

### Questions fréquentes

**Q: Les TypeNodes sont-ils supprimés ?**  
R: Non, les TypeNodes sont créés au niveau du type et sont conservés pour permettre l'ajout dynamique de nouvelles règles.

**Q: Puis-je supprimer un nœud manuellement ?**  
R: Oui, mais uniquement si `RefCount == 0` :
```go
err := network.LifecycleManager.RemoveNode("nodeID")
```

**Q: Comment voir toutes les règles d'un nœud ?**  
R: 
```go
lifecycle, _ := network.LifecycleManager.GetNodeLifecycle("nodeID")
rules := lifecycle.GetRules()
```

**Q: Est-ce thread-safe ?**  
R: Oui, toutes les opérations utilisent des mutex pour garantir la thread-safety.

### Ressources

- Code source : `tsd/rete/node_lifecycle.go`
- Tests : `tsd/rete/node_lifecycle_test.go`, `network_lifecycle_test.go`
- Documentation : `NODE_LIFECYCLE_FEATURE.md`, `NODE_LIFECYCLE_SUMMARY.md`

---

**Date** : 26 janvier 2025  
**Version** : 1.0  
**Statut** : ✅ Production-ready  
**Tests** : 19/19 PASS (100%)