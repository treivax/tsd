# 🚀 Quick Start - Système d'Indexation Delta

## 📦 Installation

Le package est déjà installé dans `rete/delta`. Aucune installation supplémentaire n'est nécessaire.

## ✅ Validation

```bash
cd /home/resinsec/dev/tsd
./rete/delta/validate.sh
```

## 🧪 Lancer les Tests

```bash
# Tous les tests
go test ./rete/delta/...

# Tests avec détails
go test -v ./rete/delta/...

# Test d'intégration
go test -v ./rete/delta/... -run TestIndexation

# Avec couverture
go test -cover ./rete/delta/...

# Benchmarks
go test -bench=. -benchmem ./rete/delta/...
```

## 💻 Utilisation Basique

### 1. Créer un Index

```go
import "github.com/treivax/tsd/rete/delta"

// Créer un nouvel index
idx := delta.NewDependencyIndex()
```

### 2. Indexer des Nœuds

```go
// Indexer un nœud alpha
idx.AddAlphaNode("alpha1", "Product", []string{"price", "status"})

// Indexer un nœud beta
idx.AddBetaNode("beta1", "Order", []string{"customer_id"})

// Indexer un nœud terminal
idx.AddTerminalNode("term1", "Product", []string{"price"})
```

### 3. Requêtes

```go
// Qui est affecté par Product.price ?
affected := idx.GetAffectedNodes("Product", "price")
for _, node := range affected {
    fmt.Printf("Nœud affecté: %s\n", node.String())
}

// Avec un FactDelta
delta := delta.NewFactDelta("Product~123", "Product")
delta.AddFieldChange("price", 100.0, 150.0)

affected := idx.GetAffectedNodesForDelta(delta)
```

### 4. Builder (construction automatique)

```go
builder := delta.NewIndexBuilder()
builder.EnableDiagnostics()

// AST d'une condition alpha
condition := map[string]interface{}{
    "type": "comparison",
    "left": map[string]interface{}{
        "type":  "fieldAccess",
        "field": "price",
    },
    "right": 100,
}

// Construire l'index
err := builder.BuildFromAlphaNode(idx, "alpha1", "Product", condition)

// Diagnostics
diag := builder.GetDiagnostics()
fmt.Printf("Nœuds traités: %d\n", diag.NodesProcessed)
```

## 📊 Statistiques

```go
stats := idx.GetStats()
fmt.Printf("Nœuds: %d, Champs: %d\n", stats.NodeCount, stats.FieldCount)
fmt.Printf("Alpha: %d, Beta: %d, Terminal: %d\n", 
    stats.AlphaNodeCount, stats.BetaNodeCount, stats.TerminalCount)
```

## 📚 Documentation Complète

- **README.md** : Guide complet d'utilisation
- **IMPLEMENTATION_REPORT_PROMPT03.md** : Rapport détaillé d'implémentation
- **EXECUTION_SUMMARY.md** : Résumé d'exécution

## 🔧 Scripts Utiles

```bash
# Validation complète
./rete/delta/validate.sh

# Tests avec race detector
go test -race ./rete/delta/...

# Couverture détaillée
go test -coverprofile=coverage.out ./rete/delta/...
go tool cover -html=coverage.out
```

## 🎯 Prochaines Étapes

Pour intégrer avec RETE (Prompt 06) :
1. Implémenter `BuildFromNetwork()` dans `IndexBuilder`
2. Connecter avec les structures RETE existantes
3. Construire l'index au moment de la compilation des règles

## 🆘 Support

Problème ? Vérifier :
1. Tests passent : `go test ./rete/delta/...`
2. Validation : `./rete/delta/validate.sh`
3. Documentation : `rete/delta/README.md`

## ✅ Checklist de Validation

- [x] Tous les tests passent (100%)
- [x] Couverture > 80% (83.8%)
- [x] Pas de race conditions
- [x] Code formaté (go fmt, goimports)
- [x] Analyse statique OK (go vet, staticcheck)
- [x] Documentation complète (GoDoc + README)

---

**Prêt à l'emploi!** 🚀
