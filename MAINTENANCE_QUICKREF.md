# 🔧 Référence Rapide - Maintenance TSD

> Commandes essentielles pour la maintenance quotidienne du projet

---

## 🚀 Démarrage Rapide

```bash
# Validation complète du projet
./scripts/validate-maintenance.sh

# Tests complets
make test

# Validation CI (ce qui tourne en CI/CD)
make validate
```

---

## 🧹 Nettoyage

### Fichiers Temporaires
```bash
# Nettoyer builds et caches
go clean -cache -testcache -modcache

# Supprimer fichiers de profiling/coverage
find . -type f \( -name "*.prof" -o -name "*.out" -o -name "*.test" \) \
  ! -path "./.git/*" -exec rm {} \;

# Tout nettoyer (via Makefile si disponible)
make clean
```

### Dépendances
```bash
# Nettoyer dépendances non utilisées
go mod tidy

# Vérifier intégrité des modules
go mod verify

# Voir pourquoi une dépendance est nécessaire
go mod why github.com/example/package

# Mettre à jour toutes les dépendances
go get -u ./...
go mod tidy
```

### Code
```bash
# Formater tout le code
go fmt ./...

# Organiser les imports
goimports -w .

# Fix basiques avec gofmt
gofmt -w -s .
```

---

## 🔍 Analyse

### Analyse Statique
```bash
# staticcheck (linter complet)
staticcheck ./...

# go vet (analyseur officiel Go)
go vet ./...

# Vérifier problèmes courants
go vet -composites=false ./...
```

### Complexité
```bash
# Complexité cyclomatique (seuil: 15)
gocyclo -over 15 .

# Top 10 fonctions les plus complexes
gocyclo -top 10 .

# Complexité moyenne
gocyclo . | awk '{sum+=$1; count++} END {print "Moyenne:", sum/count}'
```

### Code Non Utilisé
```bash
# Trouver code mort
deadcode ./...

# Installer deadcode si nécessaire
go install golang.org/x/tools/cmd/deadcode@latest
```

### Sécurité
```bash
# Vérifier vulnérabilités connues
govulncheck ./...

# Installer govulncheck si nécessaire
go install golang.org/x/vuln/cmd/govulncheck@latest

# Scanner avec gosec
gosec ./...
```

---

## 🧪 Tests

### Exécution
```bash
# Tous les tests
go test ./...

# Tests courts seulement
go test -short ./...

# Tests avec verbose
go test -v ./...

# Tests d'un package spécifique
go test ./rete/...

# Test spécifique
go test -run TestNomDuTest ./package

# Tests avec race detector
go test -race ./...

# Tests parallèles (défaut: GOMAXPROCS)
go test -parallel 4 ./...
```

### Couverture
```bash
# Couverture globale
go test -cover ./...

# Couverture détaillée (fichier)
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Couverture HTML (navigateur)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Couverture d'un package spécifique
go test -coverprofile=coverage.out ./rete
go tool cover -func=coverage.out
```

### Benchmarks
```bash
# Exécuter benchmarks
go test -bench=. ./...

# Benchmarks avec allocations mémoire
go test -bench=. -benchmem ./...

# Benchmark spécifique
go test -bench=BenchmarkNom ./package

# Comparer benchmarks
go test -bench=. -benchmem > bench-new.txt
benchcmp bench-old.txt bench-new.txt
```

---

## 📊 Profiling

### CPU
```bash
# Profiler CPU
go test -cpuprofile=cpu.prof -bench=.

# Analyser profil CPU
go tool pprof cpu.prof
# Dans pprof: top, list, web

# Profil CPU avec flamegraph
go tool pprof -http=:8080 cpu.prof
```

### Mémoire
```bash
# Profiler mémoire
go test -memprofile=mem.prof -bench=.

# Analyser profil mémoire
go tool pprof mem.prof

# Allocations mémoire
go test -benchmem -bench=.
```

### Trace
```bash
# Générer trace d'exécution
go test -trace=trace.out

# Visualiser trace
go tool trace trace.out
```

---

## 📝 TODOs et FIXMEs

```bash
# Trouver tous les TODOs
grep -rn "TODO" --include="*.go" .

# Compter TODOs
grep -r "TODO\|FIXME\|XXX\|HACK" --include="*.go" . | wc -l

# TODOs par fichier
grep -r "TODO" --include="*.go" . | cut -d: -f1 | sort | uniq -c | sort -rn

# TODOs récents (git)
git log --all --oneline --grep="TODO"
```

---

## 📈 Statistiques

### Lignes de Code
```bash
# Total lignes Go
find . -name "*.go" -not -path "*/vendor/*" | xargs wc -l | tail -1

# Lignes par package
find . -name "*.go" -not -path "*/vendor/*" -exec dirname {} \; | \
  sort -u | while read dir; do
    echo "$dir: $(find "$dir" -maxdepth 1 -name "*.go" | xargs wc -l | tail -1)"
  done
```

### Packages et Fichiers
```bash
# Nombre de packages
go list ./... | wc -l

# Nombre de fichiers Go
find . -name "*.go" -not -path "*/vendor/*" | wc -l

# Nombre de fichiers de test
find . -name "*_test.go" -not -path "*/vendor/*" | wc -l
```

### Dépendances
```bash
# Liste toutes les dépendances
go list -m all

# Graphe de dépendances
go mod graph

# Dépendances directes seulement
go list -m -json all | jq -r 'select(.Main != true) | .Path'

# Vérifier dépendances obsolètes
go list -u -m all
```

---

## 🔧 Outils à Installer

```bash
# Outils essentiels
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Outils avancés
go install golang.org/x/tools/cmd/deadcode@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Benchmarking
go install golang.org/x/perf/cmd/benchstat@latest
```

---

## 🎯 Workflows Courants

### Avant un Commit
```bash
# 1. Formater
go fmt ./...
goimports -w .

# 2. Vérifier
go vet ./...
staticcheck ./...

# 3. Tests
go test ./...

# 4. Validation complète
./scripts/validate-maintenance.sh
```

### Avant une Release
```bash
# 1. Tous les tests (incluant race)
go test -race ./...

# 2. Couverture
go test -cover ./... | tee coverage-report.txt

# 3. Benchmarks
go test -bench=. -benchmem ./... > benchmarks.txt

# 4. Sécurité
govulncheck ./...

# 5. Build
go build ./...

# 6. Validation CI
make validate
```

### Debug de Performance
```bash
# 1. Benchmark baseline
go test -bench=. -benchmem > bench-before.txt

# 2. Faire les changements

# 3. Nouveau benchmark
go test -bench=. -benchmem > bench-after.txt

# 4. Comparer
benchstat bench-before.txt bench-after.txt

# 5. Profiler si nécessaire
go test -cpuprofile=cpu.prof -bench=.
go tool pprof -http=:8080 cpu.prof
```

### Nettoyage Hebdomadaire
```bash
#!/bin/bash
# weekly-maintenance.sh

echo "🧹 Nettoyage hebdomadaire TSD"

# Nettoyage
go clean -cache -testcache -modcache
go mod tidy

# Formatage
goimports -w .

# Validation
./scripts/validate-maintenance.sh

# Rapport
echo "✅ Nettoyage terminé - $(date)"
```

---

## 📋 Checklist Quotidienne

```markdown
- [ ] `go test ./...` passe
- [ ] `go vet ./...` sans erreurs
- [ ] `goimports -w .` exécuté
- [ ] Pas de nouveaux TODOs non documentés
- [ ] Couverture ≥ 80% pour nouveaux packages
```

---

## 📋 Checklist Hebdomadaire

```markdown
- [ ] `./scripts/validate-maintenance.sh` OK
- [ ] `staticcheck ./...` vérifié
- [ ] `govulncheck ./...` sans vulnérabilités
- [ ] Dépendances à jour (`go list -u -m all`)
- [ ] CHANGELOG.md mis à jour si nécessaire
```

---

## 📋 Checklist Mensuelle

```markdown
- [ ] Rapport maintenance généré (voir maintain.md)
- [ ] TODOs triés et priorisés (MAINTENANCE_TODO.md)
- [ ] Métriques enregistrées (couverture, complexité)
- [ ] Code mort nettoyé (`deadcode ./...`)
- [ ] Benchmarks exécutés et archivés
- [ ] Documentation mise à jour
```

---

## 🆘 Dépannage

### Tests qui échouent
```bash
# Verbose pour voir ce qui se passe
go test -v ./package

# Un seul test à la fois
go test -run TestNom -v ./package

# Désactiver cache si problème
go clean -testcache
go test ./...
```

### Build qui échoue
```bash
# Nettoyer tout
go clean -cache -modcache
go mod tidy
go build ./...

# Vérifier versions Go
go version
cat go.mod | grep "^go "
```

### Dépendances cassées
```bash
# Réinitialiser dépendances
rm go.sum
go mod tidy
go mod verify
```

---

## 📚 Références

- **Guide complet** : `.github/prompts/maintain.md`
- **Standards** : `.github/prompts/common.md`
- **Rapports** : `REPORTS/MAINTENANCE_20251220.md`
- **TODOs** : `REPORTS/MAINTENANCE_TODO.md`
- **Santé projet** : `REPORTS/PROJECT_HEALTH_20251220.md`

---

## 🔗 Liens Utiles

- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [Go Profiling](https://go.dev/blog/pprof)
- [staticcheck docs](https://staticcheck.io/docs/)
- [gocyclo](https://github.com/fzipp/gocyclo)
- [govulncheck](https://go.dev/blog/vuln)

---

**Dernière mise à jour** : 2025-12-20  
**Version** : 1.0