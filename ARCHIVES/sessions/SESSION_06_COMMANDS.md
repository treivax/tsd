# 🚀 Session 06 - Commandes de Validation

Ce document liste toutes les commandes pour valider le refactoring de la session 06.

---

## ✅ Tests

### Tests unitaires JoinNode
```bash
cd /home/resinsec/dev/tsd
go test -v -run "TestJoinNode_Activate" ./rete/
```

### Tests performJoinWithTokens
```bash
go test -v -run "TestJoinNode_PerformJoinWithTokens" ./rete/
```

### Tests cascade
```bash
go test -v -run "TestJoinNodeCascade" ./rete/
```

### Tous les tests rete
```bash
go test -v ./rete/...
```

### Tests d'intégration
```bash
make test-integration
```

### Tous les tests du projet
```bash
make test
```

---

## 🔍 Vérifications Qualité

### Formatage
```bash
go fmt ./rete/node_join.go ./rete/node_join_activate_test.go
```

### Analyse statique
```bash
go vet ./rete/
```

### Compilation
```bash
go build ./...
```

### Validation complète
```bash
make validate
```

---

## 📊 Métriques

### Couverture tests
```bash
go test -cover ./rete/
```

### Couverture détaillée
```bash
go test -coverprofile=coverage.out ./rete/
go tool cover -html=coverage.out
```

---

## 🔧 Debug

### Activer le logging debug dans les tests
Les tests activent déjà `Debug = true` automatiquement.

### Voir les logs détaillés des jointures
```bash
go test -v -run "TestJoinNode_ActivateLeft_PreservesAllBindings" ./rete/
```

### Analyser une cascade spécifique
```bash
go test -v -run "TestJoinNodeCascade_ThreeVariables" ./rete/
```

---

## 📁 Fichiers Modifiés

### Voir les modifications
```bash
git diff rete/node_join.go
```

### Voir le nouveau fichier de tests
```bash
cat rete/node_join_activate_test.go
```

---

## 📖 Documentation

### Lire les rapports
```bash
ls -lh REPORTS/SESSION_06*.md
cat REPORTS/SESSION_06_README.md
```

### Documentation technique
```bash
cat REPORTS/SESSION_06_TECHNICAL_DOC.md
```

---

## ✅ Checklist Rapide

```bash
# 1. Tests passent
make test

# 2. Formatage OK
go fmt ./rete/...

# 3. Analyse OK
go vet ./rete/...

# 4. Build OK
go build ./...
```

Si toutes ces commandes réussissent, le refactoring est validé ! ✅

---

**Généré le 2025-12-12**
