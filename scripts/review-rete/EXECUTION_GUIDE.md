# 🚀 Guide d'Exécution - Revue Module RETE

**Objectif:** Guide pratique pour exécuter la série de prompts de revue du module `rete`

---

## 📋 Vue Rapide

**11 prompts** numérotés (00-10) pour une revue systématique et complète du module `rete`.

**Durée totale estimée:** 20-28 heures  
**Format:** Compatible contexte Zed (128k tokens)  
**Résultat attendu:** Module rete qualité production++

---

## 🎯 Avant de Commencer

### 1. Préparer l'Environnement

```bash
# Naviguer au projet
cd /home/resinsec/dev/tsd

# Installer outils nécessaires
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/kisielk/errcheck@latest

# Vérifier installation
gocyclo -h
staticcheck -h
```

### 2. Établir Baseline

```bash
# Sauvegarder état actuel
git checkout -b review-rete-baseline
git add -A
git commit -m "chore: baseline avant revue rete"
git tag review-rete-start

# Créer branche de travail
git checkout -b review-rete-work

# Vérifier tests passent
go test ./rete/... -v
# Doit afficher: ok  github.com/treivax/tsd/rete
```

### 3. Créer Structure Rapports

```bash
mkdir -p REPORTS/review-rete
```

---

## 🔄 Exécution Prompt par Prompt

### Prompt 00 - Overview (15 minutes)

**Action:** Lecture seule, comprendre la structure

```bash
# Ouvrir dans Zed
zed scripts/review-rete/00_overview_and_plan.md

# Ou lire directement
cat scripts/review-rete/00_overview_and_plan.md

# Comprendre:
# - Structure de la revue (11 prompts)
# - Découpage par domaine fonctionnel
# - Métriques cibles
# - Workflow d'exécution
```

**Livrable:** Compréhension du plan ✅

---

### Prompt 01 - Core RETE Nodes (2-3h) ⚠️ CRITIQUE

**Fichiers:** `network.go`, `node*.go`, `memory.go`, `token.go`

**Session Zed:**
```bash
# Ouvrir prompt
zed scripts/review-rete/01_core_rete_nodes.md

# Ajouter au contexte:
# - rete/network.go
# - rete/node.go
# - rete/node_alpha.go
# - rete/node_beta.go
# - rete/node_join.go
# - rete/node_terminal.go
# - rete/memory.go
# - rete/token.go
```

**Focus Principal:** Décomposer `evaluateSimpleJoinConditions` (complexité 26 → <15)

**Validation:**
```bash
# Tests
go test -v ./rete -run "TestNode"
go test -v ./rete -run "TestMemory"
go test -v ./rete -run "TestToken"

# Complexité
gocyclo rete/node*.go rete/memory.go rete/token.go
# Cible: Toutes fonctions <15

# Benchmarks
go test -bench=BenchmarkNode -benchmem ./rete
go test -bench=BenchmarkJoin -benchmem ./rete
```

**Commit:**
```bash
git add rete/node*.go rete/memory.go rete/token.go
git commit -m "refactor(rete): core nodes - reduce complexity

- Decompose evaluateSimpleJoinConditions (26 → X)
- Improve memory management
- Enhance thread-safety
- Add missing tests

Prompt: 01_core_rete_nodes.md"
```

**Rapport:** Créer `REPORTS/review-rete/01_report.md`

---

### Prompt 02 - Bindings & Chains (2-3h) ⚠️ CRITIQUE

**Fichiers:** `binding_chain.go`, `beta_chain*.go`, `chain_*.go`, `token_metadata.go`

**Session Zed:**
```bash
zed scripts/review-rete/02_bindings_chains.md
# + Ajouter fichiers du périmètre au contexte
```

**Focus Principal:** 
- Valider immuabilité complète
- Vérifier correction bug partage JoinNode

**Validation:**
```bash
# Tests immuabilité
go test -v ./rete -run "TestBindingChain"
go test -v ./rete -run "TestBetaChain"

# Tests régression (post-fix bug)
go test -v ./rete -run "TestBetaJoinComplex"
go test -v ./rete -run "TestJoinMultiVariable"
go test -v ./rete -run "TestBetaExhaustive"

# Thread-safety
go test -race ./rete -run "BindingChain"

# Benchmarks
go test -bench=BenchmarkBindingChain -benchmem ./rete
```

**Commit + Rapport**

---

### Prompts 03-10 - Voir Fichier Condensé

Pour les prompts 03-10, utiliser le fichier condensé:

```bash
cat scripts/review-rete/03-10_prompts_condensed.md
```

**Format par prompt:**
1. Lire section du prompt dans fichier condensé
2. Charger fichiers du périmètre dans Zed
3. Appliquer checklist de revue
4. Implémenter changements
5. Valider (tests, complexité, benchmarks)
6. Commit atomique
7. Générer rapport
8. Passer au suivant

---

## 📊 Validation Continue

### Après Chaque Prompt

```bash
# Tests module complet
go test ./rete/... -v

# Complexité
gocyclo -over 15 rete/

# Formatage
go fmt ./rete/...

# Vérifications
go vet ./rete/...
staticcheck ./rete/...

# Si tout OK → Commit
```

### Points de Contrôle (Après Prompts 3, 6, 9)

```bash
# Couverture
go test -coverprofile=coverage_checkpoint.out ./rete/...
go tool cover -func=coverage_checkpoint.out | tail -1

# Métriques
echo "Complexité max:" $(gocyclo -top 1 rete/ | head -1)
echo "Fonctions >15:" $(gocyclo -over 15 rete/ | wc -l)

# État git
git log --oneline --graph | head -20
```

---

## 🎯 Validation Finale (Après Prompt 10)

### Tests Exhaustifs

```bash
# Tous tests rete
go test ./rete/... -v -count=1

# Couverture finale
go test -coverprofile=coverage_final.out ./rete/...
go tool cover -func=coverage_final.out
go tool cover -html=coverage_final.out -o coverage_rete.html

# Cible: >85%
go tool cover -func=coverage_final.out | grep total
```

### Métriques Finales

```bash
# Complexité (CIBLE: 0 fonctions >15)
gocyclo -over 15 rete/
gocyclo -top 20 rete/

# Statistiques
echo "=== STATISTIQUES FINALES ==="
echo "Fichiers source:" $(find rete -name "*.go" -not -name "*_test.go" | wc -l)
echo "Fichiers tests:" $(find rete -name "*_test.go" | wc -l)
echo "Lignes code:" $(find rete -name "*.go" -not -name "*_test.go" | xargs wc -l | tail -1)
echo "Lignes tests:" $(find rete -name "*_test.go" | xargs wc -l | tail -1)

# Vérifications complètes
go vet ./rete/...
staticcheck ./rete/...
errcheck ./rete/...

# Race detector
go test -race ./rete/...
```

### Rapport Final

```bash
# Créer rapport synthèse
cat > REPORTS/review-rete/FINAL_REPORT.md << 'EOF'
# 🔍 Rapport Final - Revue Module RETE

## Résumé Exécutif
- Durée totale: X heures
- Prompts exécutés: 11/11
- Commits: X

## Métriques Avant/Après

| Métrique | Avant | Après | Amélioration |
|----------|-------|-------|--------------|
| Complexité max | 48 | X | ↓ Y% |
| Fonctions >15 | ~50 | X | ↓ Y |
| Couverture | 80.8% | X% | ↑ Y% |
| Warnings | 0 | 0 | ✅ |

## Changements Majeurs
1. [Liste des refactorings importants]
2. ...

## Problèmes Résolus
1. [Liste des bugs/issues fixés]
2. ...

## Recommandations Futures
1. [Suggestions d'amélioration]
2. ...

## Verdict
✅ Module rete - Qualité production++
EOF
```

---

## 📝 Templates Utiles

### Template Rapport Prompt

```markdown
# 🔍 Rapport - Prompt XX: [Domaine]

**Date:** YYYY-MM-DD
**Durée:** Xh

## Fichiers Analysés
- [Liste]

## Problèmes Identifiés

### Critiques
1. [Description] - fichier.go:ligne

### Majeurs
1. [Description] - fichier.go:ligne

### Mineurs
1. [Description] - fichier.go:ligne

## Changements Effectués
- [Liste des modifications]

## Métriques

| Métrique | Avant | Après |
|----------|-------|-------|
| Complexité max | X | Y |
| Tests | X | Y |
| Couverture | X% | Y% |

## Validation
- [ ] Tests passent
- [ ] Complexité <15
- [ ] Pas de régression

## Recommandations
1. [Si applicable]
```

### Template Commit

```
refactor(rete): [domaine] - [résumé court]

[Description détaillée des changements]

Changes:
- Change 1
- Change 2
- Fixes complexity/duplication/etc.

Metrics:
- Complexity: X → Y
- Coverage: A% → B%

Prompt: XX_domain.md
```

---

## 🔧 Commandes Utiles

### Analyse Rapide

```bash
# Top 10 complexité
gocyclo -top 10 rete/

# Fichiers avec >15
gocyclo -over 15 rete/

# Couverture rapide
go test -cover ./rete/... | grep coverage

# Warnings
go vet ./rete/... 2>&1 | grep -v "no Go files"
```

### Navigation Fichiers

```bash
# Lister nœuds
grep -r "type.*Node struct" rete/

# Lister builders
ls rete/builder*.go

# Lister métriques
ls rete/*_metrics.go rete/*_stats.go

# Trouver fonction spécifique
grep -n "func.*IngestFile" rete/*.go
```

### Debug

```bash
# Tests verbose avec logs
go test -v ./rete -run TestSpecific

# Tests avec race detector
go test -race ./rete -run TestConcurrent

# Benchmark comparaison
go test -bench=. -benchmem ./rete > bench_before.txt
# (après changements)
go test -bench=. -benchmem ./rete > bench_after.txt
# Comparer manuellement
```

---

## ⚠️ Points d'Attention

### À Faire Systématiquement
- ✅ Commit après chaque prompt
- ✅ Tests avant commit
- ✅ Rapport pour chaque prompt
- ✅ Sauvegarder progression régulièrement

### À Éviter
- ❌ Changer comportement fonctionnel
- ❌ Optimiser sans benchmarker
- ❌ Refactorer sans tests
- ❌ Commits trop gros
- ❌ Ignorer warnings

### En Cas de Problème

**Tests cassés:**
```bash
# Revenir au commit précédent
git diff HEAD
git checkout -- [fichier]

# Ou reset doux
git reset --soft HEAD~1
```

**Complexité ne baisse pas:**
- Décomposer davantage
- Extract method
- Extract function
- Simplify conditions

**Couverture baisse:**
- Ajouter tests manquants
- Tests edge cases
- Tests nouvelles fonctions

---

## 🎯 Checklist Globale

### Avant Démarrage
- [ ] Outils installés
- [ ] Baseline créée
- [ ] Tests passent
- [ ] Structure rapports créée

### Pendant Revue (Par Prompt)
- [ ] Prompt lu et compris
- [ ] Fichiers chargés dans contexte
- [ ] Analyse effectuée
- [ ] Changements implémentés
- [ ] Tests validés
- [ ] Commit effectué
- [ ] Rapport généré

### Après Revue
- [ ] Tous prompts complétés
- [ ] Tests 100% passent
- [ ] Complexité <15 partout
- [ ] Couverture >85%
- [ ] Rapport final créé
- [ ] Code review ready

---

## 📚 Références Rapides

- **Standards:** `.github/prompts/review.md`
- **Conventions:** `.github/prompts/common.md`
- **Overview:** `scripts/review-rete/00_overview_and_plan.md`
- **Condensé:** `scripts/review-rete/03-10_prompts_condensed.md`
- **README:** `scripts/review-rete/README.md`

---

**Bon courage ! 🚀**

**Questions:** Consulter les fichiers de référence ci-dessus

**Prêt à commencer ?** → Prompt 00, puis 01, puis 02...