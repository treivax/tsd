# Commandes de Validation - Clés Primaires

## 🧪 Tests

### Tests Unitaires du Module Constraint
```bash
cd /home/resinsec/dev/tsd
go test ./constraint/ -v -count=1
```

### Tests Spécifiques aux Clés Primaires
```bash
cd /home/resinsec/dev/tsd
go test ./constraint/ -v -run "PrimaryKey"
```

### Tests d'Intégration
```bash
cd /home/resinsec/dev/tsd
go test ./constraint/ -v -run "TestPrimaryKeyIntegration"
```

### Couverture de Tests
```bash
cd /home/resinsec/dev/tsd
go test ./constraint/... -cover
```

### Couverture Détaillée
```bash
cd /home/resinsec/dev/tsd
go test ./constraint/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🔍 Validation du Code

### Formatage
```bash
cd /home/resinsec/dev/tsd
go fmt ./constraint/
goimports -w constraint/
```

### Analyse Statique
```bash
cd /home/resinsec/dev/tsd
go vet ./constraint/
```

### Linting (si golangci-lint est installé)
```bash
cd /home/resinsec/dev/tsd
golangci-lint run ./constraint/
```

### Vérification des Headers
```bash
cd /home/resinsec/dev/tsd
for file in constraint/primary_key*.go; do
    if ! head -1 "$file" | grep -q "Copyright"; then
        echo "⚠️  EN-TÊTE MANQUANT: $file"
    else
        echo "✅ $file"
    fi
done
```

---

## 📊 Métriques

### Complexité Cyclomatique
```bash
cd /home/resinsec/dev/tsd
gocyclo -over 15 constraint/primary_key*.go
```

### Compte de Lignes
```bash
cd /home/resinsec/dev/tsd
wc -l constraint/primary_key*.go
```

### Statistiques du Module
```bash
cd /home/resinsec/dev/tsd
echo "=== Fichiers de validation PK ==="
ls -lh constraint/primary_key*.go
echo ""
echo "=== Tests ==="
go test ./constraint/ -count=1 2>&1 | grep -E "^(ok|PASS)"
echo ""
echo "=== Couverture ==="
go test ./constraint/ -cover 2>&1 | grep coverage
```

---

## 🔬 Tests de Régression

### Tous les Tests du Projet
```bash
cd /home/resinsec/dev/tsd
go test ./... -count=1
```

### Tests avec Race Detector
```bash
cd /home/resinsec/dev/tsd
go test ./constraint/ -race -count=1
```

---

## 📝 Vérification des Modifications

### Fichiers Modifiés
```bash
cd /home/resinsec/dev/tsd
git status
```

### Différences
```bash
cd /home/resinsec/dev/tsd
git diff constraint/
```

### Fichiers Créés
```bash
cd /home/resinsec/dev/tsd
git ls-files --others --exclude-standard constraint/ | grep primary_key
```

---

## 🚀 Validation Complète (Make)

### Tests Unitaires
```bash
cd /home/resinsec/dev/tsd
make test-unit
```

### Validation Complète
```bash
cd /home/resinsec/dev/tsd
make validate
```

---

## 📋 Checklist de Validation Manuelle

- [ ] Les 3 nouveaux fichiers existent :
  ```bash
  ls -l constraint/primary_key_validation.go
  ls -l constraint/primary_key_validation_test.go
  ls -l constraint/primary_key_integration_test.go
  ```

- [ ] Les headers de copyright sont présents :
  ```bash
  head -3 constraint/primary_key*.go
  ```

- [ ] Les tests passent :
  ```bash
  go test ./constraint/ -run PrimaryKey -v
  ```

- [ ] La couverture est > 80% :
  ```bash
  go test ./constraint/ -cover | grep coverage
  ```

- [ ] Le formatage est correct :
  ```bash
  gofmt -l constraint/primary_key*.go
  # Doit retourner vide si tout est formaté
  ```

- [ ] Aucune erreur de vet :
  ```bash
  go vet ./constraint/
  # Doit retourner exit code 0
  ```

- [ ] Les rapports sont créés :
  ```bash
  ls -l REPORTS/primary-key-validation-*.md
  ```

---

## 🔍 Recherche de Code à Migrer

### Recherche de `id` dans les faits (hors constraint)
```bash
cd /home/resinsec/dev/tsd
grep -r "Name: \"id\", Value:" --include="*.go" --exclude-dir=constraint
```

### Recherche dans les fichiers TSD
```bash
cd /home/resinsec/dev/tsd
find . -name "*.tsd" -type f | xargs grep -l "id:" | grep -v constraint
```

---

## ✅ Validation Finale

### Script de Validation Complet
```bash
#!/bin/bash
cd /home/resinsec/dev/tsd

echo "=== 🧪 Tests ==="
go test ./constraint/ -count=1 || exit 1

echo ""
echo "=== 📊 Couverture ==="
go test ./constraint/ -cover

echo ""
echo "=== 🔍 Formatage ==="
if [ -n "$(gofmt -l constraint/primary_key*.go)" ]; then
    echo "❌ Fichiers non formatés"
    gofmt -l constraint/primary_key*.go
    exit 1
else
    echo "✅ Formatage OK"
fi

echo ""
echo "=== 🔬 Analyse Statique ==="
go vet ./constraint/ || exit 1
echo "✅ Vet OK"

echo ""
echo "=== ✅ VALIDATION COMPLÈTE RÉUSSIE ==="
```

---

## 📚 Documentation Générée

### Rapports
- `REPORTS/primary-key-validation-review.md` - Revue de code complète
- `REPORTS/primary-key-validation-TODO.md` - Actions et prochaines étapes

### Consulter les Rapports
```bash
cd /home/resinsec/dev/tsd
cat REPORTS/primary-key-validation-review.md
cat REPORTS/primary-key-validation-TODO.md
```

---

**Date de Création** : 2025-12-16  
**Statut** : ✅ Validation Complète  
**Couverture** : 84.1%  
**Tests** : Tous passent
