# 📊 Statistiques du Code (Code Stats)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux obtenir des statistiques détaillées sur le code du projet : nombre de lignes par module, fichiers les plus volumineux, fonctions les plus longues, répartition du code, métriques de qualité, couverture de tests, etc.

**⚠️ IMPORTANT** : Les statistiques principales doivent concerner **uniquement le code fonctionnel manuel** (hors tests, hors code généré). Les tests et le code généré sont analysés dans des sections séparées.

## Objectif

Générer un rapport complet de statistiques sur le code source du projet pour évaluer la taille, la complexité, la qualité, et identifier les zones nécessitant potentiellement du refactoring.

## Scope des Statistiques

### ✅ CODE MANUEL (Statistiques Principales)
- Code Go dans `rete/`, `constraint/`, `cmd/`, `internal/`
- Fichiers `.go` (hors `*_test.go`, hors générés)
- Code écrit manuellement par l'équipe
- Fonctions, méthodes, structures

### 📝 CODE GÉNÉRÉ (Section Séparée)
- Fichiers avec marqueur `// Code generated` ou `DO NOT EDIT`
- Parser PEG (`constraint/parser.go`)
- Protobuf, gRPC, mocks générés
- **⚠️ Pas de recommandations de refactoring** (non modifiable)

### 🧪 TESTS (Section Séparée)
- Fichiers `*_test.go`
- Tests unitaires, d'intégration, benchmarks
- Code de test helpers et fixtures

### ❌ Exclus Complètement
- Vendor / dépendances externes
- Documentation (markdown, comments seuls)
- Scripts auxiliaires (bash, python) - sauf si demandé explicitement

## Instructions

### PHASE 1 : IDENTIFICATION DES FICHIERS

#### 1.1 Identifier Code Généré

**Commandes** :
```bash
# Trouver fichiers générés (marqueurs)
grep -l "^// Code generated\|DO NOT EDIT" --include="*.go" -r . 2>/dev/null

# Lister avec tailles
find . -name "*.go" -not -path "*/vendor/*" -exec grep -l "^// Code generated\|DO NOT EDIT" {} \; \
  -exec wc -l {} \; 2>/dev/null

# Parser spécifique (si connu)
ls -lh constraint/parser.go 2>/dev/null
```

**Rapport attendu** :
```markdown
## 🔍 IDENTIFICATION FICHIERS

### Code Généré Détecté
- `constraint/parser.go` (5,230 lignes) - Pigeon PEG parser
- Total code généré : X,XXX lignes

### Tests Détectés
- XX fichiers `*_test.go`
- Total code tests : X,XXX lignes

### Code Manuel
- XX fichiers fonctionnels
- Total code manuel : X,XXX lignes
```

### PHASE 2 : STATISTIQUES CODE MANUEL (PRINCIPAL)

#### 2.1 Comptage Global Code Manuel

**Commandes** :
```bash
# Total lignes code manuel (hors tests, hors généré)
find . -name "*.go" \
  -not -name "*_test.go" \
  -not -path "*/vendor/*" \
  -not -path "*/testdata/*" \
  -exec grep -L "^// Code generated\|DO NOT EDIT" {} \; 2>/dev/null | \
  xargs wc -l 2>/dev/null | tail -1

# Avec tokei (si disponible)
tokei --exclude '**/*_test.go' --exclude 'constraint/parser.go' \
      --exclude 'vendor/**' --exclude 'testdata/**'

# Détails par type de lignes
find . -name "*.go" \
  -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | while read f; do
  total=$(wc -l < "$f")
  code=$(grep -v "^\s*$" "$f" | grep -v "^\s*//" | wc -l)
  comments=$(grep "^\s*//" "$f" | wc -l)
  blanks=$((total - code - comments))
  echo "$f: $total total, $code code, $comments comments, $blanks blanks"
done
```

**Rapport attendu** :
```markdown
## 📊 STATISTIQUES CODE MANUEL (PRINCIPAL)

**Date** : 2025-11-26
**Commit** : abc123def
**Scope** : Code fonctionnel manuel uniquement (hors tests, hors généré)

### Lignes de Code Totales
- **Code Go fonctionnel** : X,XXX lignes
- **Commentaires** : XXX lignes
- **Lignes vides** : XXX lignes
- **Total** : X,XXX lignes

### Répartition
- Code : XX.X%
- Commentaires : XX.X%
- Lignes vides : XX.X%

### Fichiers
- **Nombre de fichiers Go** : XX fichiers
- **Moyenne lignes/fichier** : XXX lignes
- **Médiane lignes/fichier** : XXX lignes
```

#### 2.2 Éléments du Code Manuel

**Analyser** :
```bash
# Nombre de fonctions (hors tests, hors généré)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs grep -h "^func " 2>/dev/null | wc -l

# Nombre de structures
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs grep -h "^type.*struct" 2>/dev/null | wc -l

# Nombre d'interfaces
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs grep -h "^type.*interface" 2>/dev/null | wc -l

# Nombre de méthodes (approximation)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs grep -h "^func (.*)" 2>/dev/null | wc -l
```

**Rapport** :
```markdown
### Éléments du Code
- **Fonctions** : XXX (dont XXX méthodes)
- **Structures (struct)** : XXX
- **Interfaces** : XXX
- **Types custom** : XXX
- **Constantes** : XXX
```

### PHASE 3 : STATISTIQUES PAR MODULE

#### 3.1 Comptage par Répertoire

**Commande** :
```bash
# Lignes par module (hors tests, hors généré)
for dir in rete constraint cmd internal; do
  if [ -d "$dir" ]; then
    files=$(find "$dir" -name "*.go" -not -name "*_test.go" \
      -exec grep -L "^// Code generated" {} \; 2>/dev/null)
    if [ ! -z "$files" ]; then
      lines=$(echo "$files" | xargs cat 2>/dev/null | wc -l)
      count=$(echo "$files" | wc -l)
      funcs=$(echo "$files" | xargs grep -h "^func " 2>/dev/null | wc -l)
      echo "$dir: $lines lignes, $count fichiers, $funcs fonctions"
    fi
  fi
done
```

**Rapport** :
```markdown
## 📁 STATISTIQUES PAR MODULE (CODE MANUEL)

| Module | Lignes | Fichiers | % Total | Fonctions | Lignes/Fichier | Qualité |
|--------|--------|----------|---------|-----------|----------------|---------|
| `rete/` | X,XXX | XX | XX% | XXX | XXX | ✅ |
| `constraint/` | XXX | XX | XX% | XX | XXX | ✅ |
| `cmd/` | XXX | XX | XX% | XX | XXX | ✅ |
| `internal/` | XXX | XX | XX% | XX | XXX | ⚠️ |
| **TOTAL** | **X,XXX** | **XX** | **100%** | **XXX** | **XXX** | |

### Visualisation ASCII
```
rete/        ████████████████████████████████████ 65% (X,XXX lignes)
constraint/  ████████████ 22% (XXX lignes)
cmd/         ████ 8% (XXX lignes)
internal/    ██ 5% (XXX lignes)
```

### Analyse
- **Module principal** : `rete/` (XX% du code)
- **Module le plus dense** : `rete/` (XXX lignes/fichier)
- **Répartition** : Équilibrée ✅ / Déséquilibrée ⚠️
- **Modularité** : Bonne ✅ / À améliorer ⚠️
```

### PHASE 4 : FICHIERS LES PLUS VOLUMINEUX

#### 4.1 Top Fichiers Manuel

**Commande** :
```bash
# Top 10 fichiers manuels les plus gros
find . -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | while read f; do
  lines=$(wc -l < "$f")
  funcs=$(grep -c "^func " "$f" 2>/dev/null || echo 0)
  echo "$lines $funcs $f"
done | sort -rn | head -10
```

**Rapport** :
```markdown
## 📄 TOP 10 FICHIERS LES PLUS VOLUMINEUX (CODE MANUEL)

| # | Fichier | Lignes | Fonctions | Fonc/Fichier | État | Action |
|---|---------|--------|-----------|--------------|------|--------|
| 1 | `rete/constraint_pipeline.go` | 1,039 | 19 | 54.7 | 🔴 | Refactoring |
| 2 | `rete/evaluator.go` | 1,011 | 41 | 24.7 | ⚠️ | Surveiller |
| 3 | `rete/pkg/nodes/advanced_beta.go` | 726 | 30 | 24.2 | ⚠️ | Surveiller |
| 4 | `constraint/constraint_utils.go` | 586 | 18 | 32.6 | ✅ | OK |
| 5 | `rete/node_join.go` | 544 | 15 | 36.3 | ✅ | OK |
| 6 | `constraint/program_state.go` | 420 | 15 | 28.0 | ✅ | OK |
| 7 | `constraint/pkg/validator/types.go` | 340 | 12 | 28.3 | ✅ | OK |
| 8 | `rete/pkg/nodes/beta.go` | 338 | 27 | 12.5 | ✅ | OK |
| 9 | `rete/store_indexed.go` | 312 | 15 | 20.8 | ✅ | OK |
| 10 | `rete/node_accumulate.go` | 293 | 13 | 22.5 | ✅ | OK |

### Seuils d'Évaluation
- ✅ **OK** : < 500 lignes par fichier
- ⚠️ **Surveiller** : 500-800 lignes (acceptable mais à surveiller)
- 🔴 **Refactoring** : > 800 lignes (découpage recommandé)

### Fichiers Nécessitant Attention

#### 🔴 **REFACTORING RECOMMANDÉ** (> 800 lignes)
- `rete/constraint_pipeline.go` (1,039 lignes, 19 fonctions)
  - **Problème** : Responsabilités multiples (parsing, validation, création réseau)
  - **Solution** : Découper en modules : `pipeline_parser.go`, `pipeline_builder.go`, `pipeline_validator.go`
  - **Impact** : -70% complexité, meilleure testabilité

#### ⚠️ **À SURVEILLER** (500-800 lignes)
- `rete/evaluator.go` (1,011 lignes) : Envisager extraction de sous-évaluateurs
- `rete/pkg/nodes/advanced_beta.go` (726 lignes) : Acceptable pour nœuds complexes
```

### PHASE 5 : FONCTIONS LES PLUS VOLUMINEUSES

#### 5.1 Identifier Fonctions Longues

**Script d'analyse** :
```bash
# Trouver fonctions longues (code manuel uniquement)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | while read file; do
  awk -v file="$file" '
  /^func / {
    if (func_name != "") {
      lines = NR - func_start
      if (lines > 30) {
        print lines "\t" file ":" func_name "\tL" func_start
      }
    }
    func_name = $0
    func_start = NR
  }
  END {
    if (func_name != "") {
      lines = NR - func_start + 1
      if (lines > 30) {
        print lines "\t" file ":" func_name "\tL" func_start
      }
    }
  }
  ' "$file"
done | sort -rn | head -20
```

**Analyse avec gocyclo** :
```bash
# Installer gocyclo si nécessaire
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Top fonctions complexes (exclure généré)
gocyclo -over 10 -ignore "constraint/parser.go" . | sort -rn | head -20

# Moyenne par module
for dir in rete constraint cmd internal; do
  gocyclo -avg -ignore "constraint/parser.go" "$dir/" 2>/dev/null
done
```

**Rapport** :
```markdown
## 🔧 TOP 15 FONCTIONS LES PLUS VOLUMINEUSES (CODE MANUEL)

| # | Fonction | Fichier | Lignes | Complexité | État | Action |
|---|----------|---------|--------|------------|------|--------|
| 1 | `main()` | `cmd/tsd/main.go` | 189 | 15 | 🔴 | Refactoring urgent |
| 2 | `createSingleRule()` | `constraint_pipeline.go` | 178 | 22 | 🔴 | Refactoring urgent |
| 3 | `main()` | `universal-rete-runner/main.go` | 141 | 12 | 🔴 | Refactoring |
| 4 | `evaluateJoinConditions()` | `node_join.go` | 135 | 18 | 🔴 | Refactoring |
| 5 | `createExistsRule()` | `constraint_pipeline.go` | 126 | 16 | 🔴 | Refactoring |
| 6 | `evaluateValueFromMap()` | `evaluator.go` | 122 | 20 | 🔴 | Refactoring |
| 7 | `extractAggregationInfo()` | `constraint_pipeline.go` | 91 | 12 | ⚠️ | Surveiller |
| 8 | `extractJoinConditions()` | `node_join.go` | 78 | 10 | ⚠️ | Surveiller |
| 9 | `ValidateTypes()` | `validator.go` | 76 | 9 | ⚠️ | Surveiller |
| 10 | `createAccumulatorRule()` | `constraint_pipeline.go` | 74 | 11 | ⚠️ | Surveiller |
| 11 | `calculateAggregateForFacts()` | `node_accumulate.go` | 73 | 8 | ⚠️ | Surveiller |
| 12 | `ActivateRight()` | `node_alpha.go` | 72 | 9 | ⚠️ | Surveiller |
| 13 | `ConvertFactsToReteFormat()` | `constraint_utils.go` | 68 | 7 | ✅ | OK |
| 14 | `parseConstraintFile()` | `validate_coherence.go` | 66 | 8 | ✅ | OK |
| 15 | `BuildNetworkFromConstraintFileWithFacts()` | `constraint_pipeline.go` | 65 | 7 | ✅ | OK |

### Seuils d'Évaluation
- ✅ **OK** : < 50 lignes ET complexité < 10
- ⚠️ **Surveiller** : 50-100 lignes OU complexité 10-15
- 🔴 **Refactoring** : > 100 lignes OU complexité > 15

### Fonctions Nécessitant Refactoring Urgent

#### 🔴 **PRIORITÉ 1** (> 100 lignes OU complexité > 15)

1. **`cmd/tsd/main.go:main()`** (189 lignes, complexité 15)
   - **Problème** : Logique applicative mélangée avec CLI
   - **Solution** : Extraire `parseArgs()`, `initApp()`, `runCommand()`, `handleError()`
   - **Prompt suggéré** : `refactor.md`

2. **`constraint_pipeline.go:createSingleRule()`** (178 lignes, complexité 22)
   - **Problème** : Trop de responsabilités (parsing + validation + construction)
   - **Solution** : Pipeline en 3 étapes : parse → validate → build
   - **Impact** : -60% complexité

3. **`node_join.go:evaluateJoinConditions()`** (135 lignes, complexité 18)
   - **Problème** : Switch/if imbriqués, nombreux cas
   - **Solution** : Pattern Strategy ou table de dispatch par type de condition
   - **Impact** : -50% complexité

4. **`evaluator.go:evaluateValueFromMap()`** (122 lignes, complexité 20)
   - **Problème** : Grosse fonction switch avec logique imbriquée
   - **Solution** : Map de fonctions évaluatrices par type
   - **Impact** : -65% complexité, extensibilité ++

#### ⚠️ **PRIORITÉ 2** (50-100 lignes OU complexité 10-15)
- 7 fonctions identifiées (lignes 7-13 du tableau)
- Action : Surveiller lors de modifications, refactoring si ajout de logique
```

### PHASE 6 : MÉTRIQUES DE QUALITÉ

#### 6.1 Ratio Code/Commentaires

**Analyse** :
```bash
# Ratio par module (code manuel uniquement)
for dir in rete constraint cmd internal; do
  if [ -d "$dir" ]; then
    files=$(find "$dir" -name "*.go" -not -name "*_test.go" \
      -exec grep -L "^// Code generated" {} \; 2>/dev/null)
    if [ ! -z "$files" ]; then
      code=$(echo "$files" | xargs cat 2>/dev/null | grep -v "^\s*//" | grep -v "^\s*$" | wc -l)
      comments=$(echo "$files" | xargs cat 2>/dev/null | grep "^\s*//" | wc -l)
      if [ "$code" -gt 0 ]; then
        ratio=$(echo "scale=2; $comments * 100 / $code" | bc 2>/dev/null || echo "0")
        echo "$dir: $code lignes code, $comments commentaires ($ratio%)"
      fi
    fi
  fi
done
```

**Rapport** :
```markdown
## 📈 MÉTRIQUES DE QUALITÉ (CODE MANUEL)

### Ratio Code/Commentaires
| Module | Code | Commentaires | Ratio | Évaluation |
|--------|------|--------------|-------|------------|
| `rete/` | 6,811 | 550 | 8.1% | ⚠️ Insuffisant |
| `constraint/` | 3,073 | 245 | 8.0% | ⚠️ Insuffisant |
| `cmd/` | 387 | 35 | 9.0% | ⚠️ Insuffisant |
| `internal/` | 150 | 20 | 13.3% | ⚠️ Insuffisant |
| **TOTAL** | **10,421** | **850** | **8.2%** | ⚠️ |

**Seuils** :
- ✅ **Excellent** : > 20% commentaires
- 🟢 **Bon** : 15-20% commentaires
- ⚠️ **Insuffisant** : 10-15% commentaires
- 🔴 **Faible** : < 10% commentaires

**Recommandation** : Augmenter à minimum 15% (ajouter ~700 lignes de commentaires)
- Focus sur fonctions publiques (GoDoc)
- Documenter algorithmes complexes
- Ajouter exemples d'utilisation
```

#### 6.2 Complexité Cyclomatique

**Analyse** :
```bash
# Complexité moyenne par module
for dir in rete constraint cmd internal; do
  avg=$(gocyclo -avg -ignore "constraint/parser.go" "$dir/" 2>/dev/null | \
    grep "Average" | awk '{print $3}')
  max=$(gocyclo -ignore "constraint/parser.go" "$dir/" 2>/dev/null | \
    sort -rn | head -1 | awk '{print $1}')
  echo "$dir: moyenne=$avg, max=$max"
done

# Distribution des complexités
gocyclo -ignore "constraint/parser.go" . 2>/dev/null | \
  awk '{print $1}' | sort -n | uniq -c | \
  awk '{print "Complexité " $2 ": " $1 " fonctions"}'
```

**Rapport** :
```markdown
### Complexité Cyclomatique
| Module | Moyenne | Maximum | Fonctions > 10 | Fonctions > 15 | Qualité |
|--------|---------|---------|----------------|----------------|---------|
| `rete/` | 4.8 | 22 | 12 | 4 | ⚠️ |
| `constraint/` | 3.2 | 16 | 3 | 1 | ✅ |
| `cmd/` | 5.5 | 15 | 2 | 1 | ⚠️ |
| `internal/` | 2.1 | 8 | 0 | 0 | ✅ |
| **GLOBAL** | **4.2** | **22** | **17** | **6** | ⚠️ |

**Distribution Complexité** :
```
1-5:   ████████████████████████████████ 425 fonctions (85%)
6-10:  ████ 45 fonctions (9%)
11-15: ██ 11 fonctions (2%)
16-20: █ 4 fonctions (0.8%)
21+:   █ 2 fonctions (0.4%)
```

**Seuils** :
- ✅ **Excellent** : Moyenne < 5, Max < 10
- 🟢 **Bon** : Moyenne < 8, Max < 15
- ⚠️ **Acceptable** : Moyenne < 10, Max < 20
- 🔴 **Problématique** : Moyenne > 10 OU Max > 20

**Actions** :
- 6 fonctions avec complexité > 15 nécessitent refactoring urgent
- 11 fonctions avec complexité 11-15 à surveiller
```

#### 6.3 Longueur Moyenne des Fonctions

**Analyse** :
```bash
# Calculer longueur moyenne par module
for dir in rete constraint cmd internal; do
  files=$(find "$dir" -name "*.go" -not -name "*_test.go" \
    -exec grep -L "^// Code generated" {} \; 2>/dev/null)
  if [ ! -z "$files" ]; then
    total_lines=0
    total_funcs=0
    for file in $files; do
      func_count=$(grep -c "^func " "$file" 2>/dev/null || echo 0)
      file_lines=$(wc -l < "$file")
      total_funcs=$((total_funcs + func_count))
      total_lines=$((total_lines + file_lines))
    done
    if [ "$total_funcs" -gt 0 ]; then
      avg=$((total_lines / total_funcs))
      echo "$dir: $avg lignes/fonction (approximatif)"
    fi
  fi
done
```

**Rapport** :
```markdown
### Longueur des Fonctions
| Module | Moyenne | Médiane | > 50 lignes | > 100 lignes | Qualité |
|--------|---------|---------|-------------|--------------|---------|
| `rete/` | 28 | 22 | 18 | 6 | ⚠️ |
| `constraint/` | 19 | 15 | 3 | 0 | ✅ |
| `cmd/` | 35 | 20 | 4 | 2 | ⚠️ |
| `internal/` | 15 | 12 | 0 | 0 | ✅ |
| **GLOBAL** | **24** | **18** | **25** | **8** | ⚠️ |

**Seuils** :
- ✅ **Excellent** : < 25 lignes/fonction en moyenne
- 🟢 **Bon** : 25-40 lignes/fonction
- ⚠️ **Acceptable** : 40-60 lignes/fonction
- 🔴 **Problématique** : > 60 lignes/fonction

**Distribution** :
```
0-25:   ████████████████████████████ 380 fonctions (76%)
26-50:  █████ 70 fonctions (14%)
51-100: ██ 25 fonctions (5%)
101+:   █ 8 fonctions (1.6%)
```

**Actions** :
- 8 fonctions > 100 lignes nécessitent refactoring urgent
- 25 fonctions > 50 lignes à surveiller et simplifier si possible
```

#### 6.4 Duplication de Code

**Analyse** :
```bash
# Utiliser simian, jscpd ou dupl
go install github.com/mibk/dupl@latest

# Détecter duplications
dupl -t 50 -ignore "constraint/parser.go" ./...

# Ou avec simple grep pour patterns répétés
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs -I {} sh -c 'echo "=== {} ===" && grep -E "if err != nil|return nil, err" {} | head -5'
```

**Rapport** :
```markdown
### Duplication de Code
| Type | Occurrences | Impact | Action |
|------|-------------|--------|--------|
| **Error handling patterns** | ~120 | Moyen | Helpers possibles |
| **Type assertions répétées** | ~45 | Faible | Acceptable |
| **Blocs similaires > 50 lignes** | 3 | Élevé | Extraire fonctions |

**Duplication Détectée** :
- 3 blocs de code quasi-identiques (50-70 lignes) dans `constraint_pipeline.go`
  - `createSingleRule()`, `createExistsRule()`, `createAccumulatorRule()`
  - **Solution** : Extraire logique commune dans `createRuleBase()`
  - **Impact** : -150 lignes, meilleure maintenabilité

**Seuils** :
- ✅ **Excellent** : < 5% duplication
- 🟢 **Bon** : 5-10% duplication
- ⚠️ **Acceptable** : 10-15% duplication
- 🔴 **Problématique** : > 15% duplication

**Évaluation** : ~8% duplication (acceptable) ✅
```

### PHASE 7 : COUVERTURE ET TESTS

#### 7.1 Statistiques Tests

**Commandes** :
```bash
# Compter lignes de tests
find . -name "*_test.go" -not -path "*/vendor/*" -exec wc -l {} + | tail -1

# Nombre de fonctions de test
grep -r "^func Test\|^func Benchmark" --include="*_test.go" . | wc -l

# Ratio tests/code
code_lines=$(find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | xargs cat | wc -l)
test_lines=$(find . -name "*_test.go" -exec cat {} \; | wc -l)
ratio=$(echo "scale=2; $test_lines * 100 / $code_lines" | bc)
echo "Ratio tests/code: $ratio%"
```

**Rapport** :
```markdown
## 🧪 STATISTIQUES TESTS

### Volume Tests
- **Lignes de tests** : 5,241 lignes
- **Fichiers de test** : 21 fichiers
- **Fonctions de test** : 145 tests
- **Benchmarks** : 8 benchmarks
- **Ratio tests/code** : 50.3% (excellent ✅)

### Répartition Tests par Module
| Module | Fichiers | Lignes | Tests | Ratio Local |
|--------|----------|--------|-------|-------------|
| `rete/` | 10 | 2,450 | 68 | 36% |
| `constraint/` | 6 | 1,890 | 52 | 61% |
| `cmd/` | 1 | 120 | 3 | 31% |
| `test/integration/` | 4 | 780 | 22 | N/A |

**Évaluation** : Couverture tests excellente en volume ✅
```

#### 7.2 Couverture Fonctionnelle

**Commandes** :
```bash
# Exécuter tests avec couverture
go test -coverprofile=/tmp/coverage.out ./...

# Afficher couverture globale
go tool cover -func=/tmp/coverage.out | tail -1

# Couverture par package
go tool cover -func=/tmp/coverage.out | grep "^.*\.go:" | \
  awk '{package=$1; sub(/\/[^\/]+$/, "", package); sum[package]+=$3; count[package]++} 
       END {for (p in sum) print p ": " sum[p]/count[p] "%"}'

# Identifier fichiers sans couverture
go tool cover -func=/tmp/coverage.out | awk '$3 == "0.0%" {print $1}' | sort | uniq
```

**Rapport** :
```markdown
### Couverture de Tests (Coverage)

**Couverture Globale** : 41.8% des statements

| Package | Coverage | État | Priorité |
|---------|----------|------|----------|
| `github.com/treivax/tsd/constraint` | 59.2% | 🟢 Bon | Maintenir |
| `github.com/treivax/tsd/rete` | 34.3% | ⚠️ Insuffisant | Améliorer |
| `github.com/treivax/tsd/test/integration` | 29.4% | ⚠️ Insuffisant | Améliorer |
| `github.com/treivax/tsd/cmd/tsd` | 0.0% | 🔴 Aucune | Urgent |
| `github.com/treivax/tsd/cmd/universal-rete-runner` | 0.0% | 🔴 Aucune | Urgent |
| `github.com/treivax/tsd/constraint/pkg/validator` | 0.0% | 🔴 Aucune | Urgent |
| `github.com/treivax/tsd/rete/pkg/nodes` | 0.0% | 🔴 Aucune | Urgent |

### Visualisation Coverage
```
constraint     ███████████████████████████████████████████████████████████ 59.2%
rete           ██████████████████████████████████ 34.3%
integration    █████████████████████████████ 29.4%
cmd            ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 0.0%
validator      ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ 0.0%
```

**Seuils** :
- ✅ **Excellent** : > 80% coverage
- 🟢 **Bon** : 60-80% coverage
- ⚠️ **Insuffisant** : 40-60% coverage
- 🔴 **Faible** : < 40% coverage

**Évaluation Globale** : 41.8% (insuffisant) ⚠️

### Fichiers Sans Couverture (0%)
```
- cmd/tsd/main.go
- cmd/universal-rete-runner/main.go
- constraint/pkg/validator/validator.go
- rete/pkg/nodes/advanced_beta.go
- rete/pkg/nodes/beta.go
- rete/pkg/domain/*.go
- constraint/internal/config/config.go
```

### Actions Recommandées

#### 🔴 **URGENT** - Packages sans tests (0% coverage)
1. **Ajouter tests CLI** (`cmd/*`)
   - Tests d'intégration CLI avec entrées/sorties mockées
   - Objectif : 40% coverage minimum
   - Prompt : `add-test.md`

2. **Tester validateur** (`constraint/pkg/validator`)
   - Tests unitaires pour chaque règle de validation
   - Objectif : 70% coverage
   - Impact critique : validation des contraintes

3. **Tester nœuds RETE** (`rete/pkg/nodes`)
   - Tests unitaires pour chaque type de nœud
   - Objectif : 60% coverage
   - Impact critique : cœur du moteur

#### ⚠️ **PRIORITÉ 2** - Augmenter coverage insuffisante
1. **rete/** (34.3% → 60%)
   - Focus sur : `evaluator.go`, `node_join.go`, `constraint_pipeline.go`
   - Ajouter tests pour cas edge et erreurs

2. **test/integration/** (29.4% → 50%)
   - Compléter scénarios d'intégration
   - Tests end-to-end avec fichiers .constraint/.facts

#### 🟢 **MAINTIEN** - Coverage acceptable
1. **constraint/** (59.2%)
   - Maintenir niveau actuel
   - Ajouter tests pour nouvelles features
```

#### 7.3 Qualité des Tests

**Analyse** :
```bash
# Identifier tests sans assertions
grep -r "^func Test" --include="*_test.go" -A 20 . | \
  grep -v "t.Error\|t.Fatal\|assert\|require\|if.*!=" | \
  grep "^func Test" | head -10

# Compter mocks et fixtures
find . -name "*_test.go" -exec grep -l "mock\|Mock\|stub\|Stub" {} \; | wc -l

# Tests d'intégration vs unitaires
unit_tests=$(find . -name "*_test.go" -not -path "*/test/integration/*" | wc -l)
integration_tests=$(find ./test/integration -name "*_test.go" 2>/dev/null | wc -l)
echo "Tests unitaires: $unit_tests, Tests intégration: $integration_tests"
```

**Rapport** :
```markdown
### Qualité des Tests

| Métrique | Valeur | Évaluation |
|----------|--------|------------|
| **Tests unitaires** | 17 fichiers | 🟢 |
| **Tests d'intégration** | 4 fichiers | ✅ |
| **Tests end-to-end** | Présents | ✅ |
| **Benchmarks** | 8 | 🟢 |
| **Mocks utilisés** | 3 packages | ⚠️ |
| **Table-driven tests** | Majoritaire | ✅ |
| **Tests sans assertions** | 0 détectés | ✅ |

**Points Forts** :
- ✅ Bonne utilisation de table-driven tests
- ✅ Tests d'intégration avec fichiers .constraint/.facts réels
- ✅ Benchmarks pour mesurer performances

**Points d'Amélioration** :
- ⚠️ Peu de mocks (interfaces non testées en isolation)
- ⚠️ Manque tests pour validation et nœuds
- ⚠️ Coverage CLI inexistante

**Évaluation Globale** : 🟢 Bonne (avec amélioration needed)
```

### PHASE 8 : CODE GÉNÉRÉ (SECTION SÉPARÉE)

#### 8.1 Statistiques Code Généré

**Commandes** :
```bash
# Identifier et compter code généré
find . -name "*.go" -not -path "*/vendor/*" \
  -exec grep -l "^// Code generated\|DO NOT EDIT" {} \; | while read f; do
  lines=$(wc -l < "$f")
  funcs=$(grep -c "^func " "$f" 2>/dev/null || echo 0)
  generator=$(head -5 "$f" | grep -i "generated by\|Code generated" | sed 's/.*generated by //' | sed 's/;.*//')
  echo "$lines $funcs $generator $f"
done | sort -rn
```

**Rapport** :
```markdown
## 🤖 CODE GÉNÉRÉ (NON MODIFIABLE)

**⚠️ Important** : Le code généré ne peut pas être modifié manuellement. Les recommandations ne s'appliquent pas à ces fichiers.

### Fichiers Générés Détectés

| Fichier | Lignes | Fonctions | Générateur | % du Total |
|---------|--------|-----------|------------|------------|
| `constraint/parser.go` | 5,230 | 216 | Pigeon PEG | 50.2% code total |

### Statistiques Globales Code Généré
- **Total lignes générées** : 5,230 lignes
- **Fichiers générés** : 1 fichier
- **% du projet** : 50.2% du code Go total (incluant généré)
- **% sans généré** : Le code manuel représente 66.6% du projet

### Impact du Code Généré

**Répartition Réelle du Projet** :
```
Code manuel      ████████████████████████████████████████ 10,421 lignes (66.6%)
Code généré      ████████████████████████ 5,230 lignes (33.4%)
Tests            ████████████████████████ 5,241 lignes (ratio 50.3%)
```

### Analyse

- **Parser PEG (`constraint/parser.go`)** :
  - Généré automatiquement par Pigeon depuis grammaire PEG
  - Fichier volumineux (5,230 lignes) mais c'est normal pour un parser généré
  - **Ne pas modifier manuellement** - modifier la grammaire source à la place
  - Exclu des statistiques de qualité (complexité, commentaires, etc.)
  - Source probable : `constraint/grammar.peg` ou similaire

**Note** : Le code généré est exclu de toutes les recommandations de refactoring et des métriques de qualité.
```

### PHASE 9 : TENDANCES ET ÉVOLUTION

#### 9.1 Historique Git

**Commandes** :
```bash
# Évolution sur 6 derniers mois
for i in {0..5}; do
  date=$(date -d "$((6-i)) months ago" +%Y-%m-01)
  commit=$(git rev-list -n 1 --before="$date" HEAD 2>/dev/null)
  if [ ! -z "$commit" ]; then
    lines=$(git ls-tree -r $commit --name-only | \
      grep "\.go$" | grep -v "_test\.go$" | \
      xargs git show $commit:{} 2>/dev/null | wc -l 2>/dev/null || echo "0")
    commits=$(git rev-list --count --since="$date" --before="$(date -d "$((6-i-1)) months ago" +%Y-%m-01)" HEAD 2>/dev/null || echo "0")
    echo "$date: $lines lignes, $commits commits"
  fi
done

# Contributeurs principaux
git shortlog -sn --no-merges --since="6 months ago" | head -5

# Fichiers les plus modifiés
git log --since="6 months ago" --pretty=format: --name-only -- "*.go" | \
  grep -v "^$" | sort | uniq -c | sort -rn | head -10
```

**Rapport** :
```markdown
## 📊 TENDANCES ET ÉVOLUTION

### Évolution Volume Code (6 derniers mois)

| Mois | Lignes Code | Variation | Commits | Activité |
|------|-------------|-----------|---------|----------|
| 2025-06 | 8,500 | - | 45 | 🟢 |
| 2025-07 | 9,200 | +8.2% | 38 | 🟢 |
| 2025-08 | 9,850 | +7.1% | 52 | 🟢 |
| 2025-09 | 10,100 | +2.5% | 28 | 🟡 |
| 2025-10 | 10,350 | +2.5% | 31 | 🟡 |
| 2025-11 | 10,421 | +0.7% | 35 | 🟡 |

**Croissance Globale** : +22.6% en 6 mois (acceptable pour projet en développement)

### Visualisation Évolution
```
11K |                                                        *
10K |                                              *    *    *
 9K |                                    *    *
 8K |                          *
 7K |                     *
 6K +----+----+----+----+----+----+----
    Jun  Jul  Aug  Sep  Oct  Nov  Déc

Tendance: Croissance ralentie (stabilisation) ✅
```

### Activité Récente (30 derniers jours)

| Semaine | Commits | Lignes +/- | Fichiers modifiés |
|---------|---------|------------|-------------------|
| Sem 48 | 9 | +145/-68 | 12 |
| Sem 47 | 12 | +234/-89 | 18 |
| Sem 46 | 8 | +89/-45 | 9 |
| Sem 45 | 6 | +56/-23 | 7 |

### Contributeurs (6 mois)

| Contributeur | Commits | Lignes | % Activité |
|--------------|---------|--------|------------|
| Développeur A | 145 | +3,200 | 65% |
| Développeur B | 68 | +1,450 | 30% |
| Développeur C | 16 | +320 | 5% |

### Fichiers les Plus Modifiés

| Fichier | Modifications | Impact |
|---------|---------------|--------|
| `rete/constraint_pipeline.go` | 34 fois | 🔴 Hotspot |
| `rete/evaluator.go` | 28 fois | 🔴 Hotspot |
| `rete/node_join.go` | 22 fois | ⚠️ Fréquent |
| `constraint/program_state.go` | 18 fois | ⚠️ Fréquent |
| `rete/store_indexed.go` | 15 fois | ⚠️ Fréquent |

**Hotspots Détectés** : Les 2 fichiers les plus modifiés sont aussi les plus volumineux → **Refactoring urgent recommandé**

### Commits Impactants (3 mois)

| Date | Commit | Impact | Description |
|------|--------|--------|-------------|
| 2025-11-18 | `a3f2b1` | +450 lignes | Ajout opérateurs chaînes |
| 2025-10-22 | `d5e8c3` | +380 lignes | Support agrégation COUNT |
| 2025-09-15 | `7b9a2f` | -220 lignes | Refactoring store_indexed |
| 2025-08-30 | `2c4d6e` | +520 lignes | Nœuds beta avancés |

**Observation** : Bonne alternance entre ajouts de features et refactoring ✅
```

#### 9.2 Vélocité et Productivité

**Analyse** :
```bash
# Lignes ajoutées/supprimées par mois
for i in {0..5}; do
  month=$(date -d "$((6-i)) months ago" +%Y-%m)
  since=$(date -d "$((6-i)) months ago" +%Y-%m-01)
  until=$(date -d "$((6-i-1)) months ago" +%Y-%m-01)
  stats=$(git log --since="$since" --until="$until" --numstat --pretty=format: -- "*.go" | \
    awk '{added+=$1; removed+=$2} END {print added " " removed " " added-removed}')
  echo "$month: $stats"
done

# Taille moyenne des commits
git log --since="6 months ago" --shortstat -- "*.go" | \
  grep "files changed" | \
  awk '{files+=$1; inserted+=$4; deleted+=$6; count++} 
       END {print "Moyenne: " files/count " fichiers, " inserted/count " lignes ajoutées, " deleted/count " lignes supprimées par commit"}'
```

**Rapport** :
```markdown
### Vélocité Développement

| Mois | Lignes + | Lignes - | Net | Productivité |
|------|----------|----------|-----|--------------|
| 2025-06 | 2,450 | 1,120 | +1,330 | 🟢 Élevée |
| 2025-07 | 1,880 | 980 | +900 | 🟢 Élevée |
| 2025-08 | 2,100 | 1,450 | +650 | 🟢 Bonne |
| 2025-09 | 890 | 640 | +250 | 🟡 Modérée |
| 2025-10 | 780 | 530 | +250 | 🟡 Modérée |
| 2025-11 | 420 | 349 | +71 | 🟡 Faible |

**Tendance** : Ralentissement normal (projet mature) ✅

### Taille Moyenne des Commits
- **Fichiers modifiés** : 3.2 fichiers/commit
- **Lignes ajoutées** : 45 lignes/commit
- **Lignes supprimées** : 28 lignes/commit
- **Net** : +17 lignes/commit

**Évaluation** : Commits de taille raisonnable (ni micro-commits, ni commits géants) ✅
```

## Format de Réponse Complet

```markdown
# 📊 RAPPORT STATISTIQUES CODE - TSD

**Date** : 2025-11-26  
**Commit** : `b4e9916` (2025-11-25 19:49:44)  
**Branche** : main  
**Scope** : Code manuel uniquement (hors tests, hors généré)

---

## 📈 RÉSUMÉ EXÉCUTIF

### Vue d'Ensemble
- **Lignes de code manuel** : 10,421 lignes (66.6% du projet)
- **Lignes de code généré** : 5,230 lignes (33.4% du projet)
- **Lignes de tests** : 5,241 lignes (ratio 50.3% - excellent)
- **Fichiers Go fonctionnels** : 49 fichiers
- **Fonctions/Méthodes** : 487 fonctions

### Indicateurs Qualité
| Indicateur | Valeur | Cible | État |
|------------|--------|-------|------|
| **Lignes/Fichier** | 213 | < 400 | ✅ |
| **Lignes/Fonction** | 24 | < 50 | ✅ |
| **Complexité Moyenne** | 4.2 | < 8 | ✅ |
| **Ratio Commentaires** | 8.2% | > 15% | ⚠️ |
| **Coverage Tests** | 41.8% | > 70% | ⚠️ |
| **Fichiers > 800 lignes** | 1 | 0 | ⚠️ |
| **Fonctions > 100 lignes** | 6 | 0 | ⚠️ |

### 🎯 Priorités
1. 🔴 **Urgent** : Refactoriser `constraint_pipeline.go` (1,039 lignes, 22 complexité max)
2. 🔴 **Urgent** : Ajouter tests pour packages à 0% coverage (cmd, validator, nodes)
3. ⚠️ **Important** : Augmenter commentaires de 8.2% à 15% (+700 lignes)
4. ⚠️ **Important** : Simplifier 6 fonctions avec complexité > 15

---

## 📊 STATISTIQUES CODE MANUEL (PRINCIPAL)

[Inclure ici toutes les sections de la PHASE 2 à PHASE 6]

---

## 🧪 STATISTIQUES TESTS

[Inclure ici PHASE 7 complète]

---

## 🤖 CODE GÉNÉRÉ (NON MODIFIABLE)

[Inclure ici PHASE 8 complète]

---

## 📊 TENDANCES ET ÉVOLUTION

[Inclure ici PHASE 9 complète]

---

## 🎯 RECOMMANDATIONS DÉTAILLÉES

### 🔴 PRIORITÉ 1 - URGENT (À faire cette semaine)

#### 1. Refactoriser `rete/constraint_pipeline.go`
- **Problème** : 1,039 lignes, complexité max 22, modifié 34 fois (hotspot)
- **Impact** : Maintenabilité critique, bugs fréquents
- **Solution** :
  ```
  constraint_pipeline.go (1039 lignes)
  ↓ Découper en ↓
  ├── pipeline_parser.go (~300 lignes) - Parsing expressions
  ├── pipeline_validator.go (~250 lignes) - Validation règles
  ├── pipeline_builder.go (~350 lignes) - Construction réseau
  └── pipeline_helpers.go (~150 lignes) - Utilitaires
  ```
- **Prompts suggérés** : `refactor.md`, `deep-clean.md`
- **Estimation** : 4-6h de travail

#### 2. Ajouter tests pour packages critiques à 0%
- **Packages concernés** :
  - `cmd/tsd` et `cmd/universal-rete-runner` (CLI)
  - `constraint/pkg/validator` (validation contraintes)
  - `rete/pkg/nodes` (cœur du moteur RETE)
- **Objectif** : Minimum 40% coverage pour chaque package
- **Impact** : Fiabilité critique du système
- **Prompt suggéré** : `add-test.md`
- **Estimation** : 8-12h de travail

#### 3. Simplifier fonctions avec complexité > 15
- **Fonctions concernées** :
  1. `createSingleRule()` - complexité 22
  2. `evaluateValueFromMap()` - complexité 20
  3. `evaluateJoinConditions()` - complexité 18
  4. `createExistsRule()` - complexité 16
- **Solution** : Décomposer en sous-fonctions, pattern Strategy
- **Prompt suggéré** : `refactor.md`
- **Estimation** : 6-8h de travail

### ⚠️ PRIORITÉ 2 - IMPORTANT (À faire ce sprint)

#### 4. Augmenter documentation (8.2% → 15%)
- **Actions** :
  - Ajouter GoDoc sur toutes fonctions publiques (règle golint)
  - Documenter algorithmes complexes dans `evaluator.go` et `node_join.go`
  - Ajouter exemples d'utilisation dans packages principaux
- **Volume** : ~700 lignes de commentaires à ajouter
- **Prompt suggéré** : `update-docs.md`
- **Estimation** : 4-6h de travail

#### 5. Augmenter coverage tests (41.8% → 70%)
- **Modules prioritaires** :
  - `rete/` : 34.3% → 60% (+750 lignes tests)
  - `test/integration/` : 29.4% → 50% (+400 lignes tests)
- **Focus** : Cas edge, gestion erreurs, scénarios complexes
- **Prompt suggéré** : `add-test.md`
- **Estimation** : 10-15h de travail

#### 6. Refactoriser fonctions longues (> 100 lignes)
- **Fonctions concernées** : 6 fonctions identifiées
- **Objectif** : < 50 lignes par fonction
- **Impact** : Lisibilité, testabilité
- **Estimation** : 4-6h de travail

### 💡 PRIORITÉ 3 - AMÉLIORATION CONTINUE

#### 7. Réduire duplication de code
- **Cibles** : 3 blocs similaires dans `constraint_pipeline.go`
- **Solution** : Extraire logique commune
- **Impact** : -150 lignes, meilleure maintenabilité

#### 8. Implémenter linting continu
- **Outils** : `golangci-lint`, `gocyclo`, `golines`
- **Seuils CI/CD** :
  - Complexité max : 15
  - Lignes/fonction max : 80
  - Coverage min : 60%
- **Estimation** : 2-3h setup

#### 9. Monitoring métriques qualité
- **Mettre en place** : Dashboard métriques (SonarQube, CodeClimate)
- **Tracker** : Évolution coverage, complexité, duplication
- **Fréquence** : Analyse mensuelle automatique

---

## 🔗 PROMPTS SUGGÉRÉS

Pour agir sur ces statistiques :
- **Refactoring urgent** → [`refactor.md`](refactor.md)
- **Nettoyage profond** → [`deep-clean.md`](deep-clean.md)
- **Ajouter tests** → [`add-test.md`](add-test.md)
- **Améliorer docs** → [`update-docs.md`](update-docs.md)
- **Review qualité** → [`code-review.md`](code-review.md)
- **Analyser erreur** → [`analyze-error.md`](analyze-error.md)

---

## 📌 NOTES TECHNIQUES

### Méthodologie
- **Code manuel** : Calculé en excluant fichiers avec `// Code generated` et `*_test.go`
- **Complexité** : Mesurée avec `gocyclo` (complexité cyclomatique)
- **Coverage** : Mesurée avec `go test -cover`
- **Duplication** : Détectée avec `dupl` (seuil 50 tokens)

### Seuils de Référence
Basés sur bonnes pratiques Go (Effective Go, Go Code Review Comments) :
- **Fichier** : < 500 lignes (idéal), < 800 (acceptable)
- **Fonction** : < 50 lignes (idéal), < 80 (acceptable)
- **Complexité** : < 10 (idéal), < 15 (acceptable)
- **Commentaires** : > 15% (idéal), > 10% (minimum)
- **Coverage** : > 80% (idéal), > 60% (acceptable)

### Exclusions Importantes
- ⚠️ **Parser généré** (`constraint/parser.go`) exclu de toutes statistiques de qualité
- ⚠️ **Tests** (`*_test.go`) comptés séparément
- ⚠️ **Vendor** et **testdata** toujours exclus

**Prochaine analyse recommandée** : Dans 1 mois (après refactoring priorité 1)

---

**📊 Rapport généré avec prompt `stats-code.md`**  
**Version** : 2.0  
**Dernière mise à jour** : Novembre 2025
```

## Exemple d'Utilisation

```
Je veux connaître les statistiques complètes du code du projet TSD :
- Combien de lignes de code Go fonctionnel manuel (hors tests, hors généré) ?
- Quels sont les fichiers les plus gros ?
- Quelles fonctions sont trop longues ou complexes ?
- Quelle est la couverture de tests ?
- Y a-t-il du code nécessitant refactoring ?

Utilise le prompt "stats-code" pour générer un rapport complet.
```

## Commandes Utiles

### Identification Fichiers

```bash
# Trouver code généré
grep -rl "^// Code generated\|DO NOT EDIT" --include="*.go" .

# Lister tous les types de fichiers
echo "=== Code Manuel ===" && \
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | wc -l && \
echo "=== Tests ===" && \
find . -name "*_test.go" | wc -l && \
echo "=== Généré ===" && \
grep -rl "^// Code generated" --include="*.go" . | wc -l
```

### Comptage Code Manuel

```bash
# Lignes code manuel (hors tests, hors généré)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs cat 2>/dev/null | wc -l

# Par module
for dir in rete constraint cmd internal; do
  lines=$(find "$dir" -name "*.go" -not -name "*_test.go" \
    -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
    xargs cat 2>/dev/null | wc -l)
  echo "$dir: $lines lignes"
done
```

### Comptage Tests

```bash
# Lignes de tests
find . -name "*_test.go" -exec cat {} + | wc -l

# Nombre de fonctions de test
grep -r "^func Test\|^func Benchmark" --include="*_test.go" . | wc -l

# Ratio tests/code
code=$(find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs cat 2>/dev/null | wc -l)
test=$(find . -name "*_test.go" -exec cat {} + | wc -l)
ratio=$(echo "scale=2; $test * 100 / $code" | bc 2>/dev/null || echo "0")
echo "Ratio tests/code: $ratio%"
```

### Couverture Tests

```bash
# Exécuter tests avec couverture
go test -coverprofile=/tmp/coverage.out ./...

# Afficher couverture globale
go tool cover -func=/tmp/coverage.out | tail -1

# Couverture par package
go tool cover -func=/tmp/coverage.out | grep -v "total:" | \
  awk '{print $1 ": " $3}'

# Fichiers sans couverture
go tool cover -func=/tmp/coverage.out | awk '$3 == "0.0%" {print $1}'
```

### Analyse Complexité

```bash
# Installer gocyclo
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Complexité > 10 (hors généré)
gocyclo -over 10 -ignore "constraint/parser.go" . | sort -rn

# Moyenne par module
for dir in rete constraint cmd internal; do
  gocyclo -avg -ignore "constraint/parser.go" "$dir/" 2>/dev/null
done

# Top 20 fonctions complexes
gocyclo -ignore "constraint/parser.go" . | sort -rn | head -20
```

### Fichiers Volumineux

```bash
# Top 10 fichiers (code manuel)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | while read f; do
  lines=$(wc -l < "$f")
  echo "$lines $f"
done | sort -rn | head -10

# Avec détails (lignes + fonctions)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | while read f; do
  lines=$(wc -l < "$f")
  funcs=$(grep -c "^func " "$f" 2>/dev/null || echo 0)
  echo "$lines $funcs $f"
done | sort -rn | head -10
```

### Fonctions Longues

```bash
# Fonctions > 50 lignes (code manuel uniquement)
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | while read file; do
  awk -v file="$file" '
  /^func / {
    if (name) {
      len = NR - start
      if (len > 50) print len "\t" file ":" name "\tL" start
    }
    name = $0; start = NR
  }
  END {
    if (name) {
      len = NR - start + 1
      if (len > 50) print len "\t" file ":" name "\tL" start
    }
  }
  ' "$file"
done | sort -rn | head -20
```

### Duplication Code

```bash
# Installer dupl
go install github.com/mibk/dupl@latest

# Détecter duplications (seuil 50 tokens)
dupl -t 50 -ignore "constraint/parser.go" ./...

# Patterns répétés
find . -name "*.go" -not -name "*_test.go" \
  -exec grep -L "^// Code generated" {} \; 2>/dev/null | \
  xargs grep -h "if err != nil" | sort | uniq -c | sort -rn | head -10
```

## Outils Recommandés

### Comptage de Lignes

- **[tokei](https://github.com/XAMPPRocky/tokei)** - Moderne, rapide, coloré (Rust)
  ```bash
  cargo install tokei
  tokei --exclude '**/*_test.go' --exclude 'constraint/parser.go'
  ```

- **[cloc](https://github.com/AlDanial/cloc)** - Classique, fiable, multi-langages
  ```bash
  cloc --exclude-dir=vendor,testdata --not-match-f='_test\.go$|parser\.go$' .
  ```

- **[scc](https://github.com/boyter/scc)** - Très rapide, statistiques détaillées
  ```bash
  scc --exclude-dir vendor,testdata --not-match '_test\.go$|parser\.go$'
  ```

- **[gocloc](https://github.com/hhatto/gocloc)** - Spécialisé Go
  ```bash
  go install github.com/hhatto/gocloc/cmd/gocloc@latest
  gocloc --not-match='_test\.go$|parser\.go$' .
  ```

### Analyse de Complexité

- **[gocyclo](https://github.com/fzipp/gocyclo)** - Complexité cyclomatique
  ```bash
  go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
  gocyclo -over 10 .
  ```

- **[gocognit](https://github.com/uudashr/gocognit)** - Complexité cognitive
  ```bash
  go install github.com/uudashr/gocognit/cmd/gocognit@latest
  gocognit -over 15 .
  ```

- **[golangci-lint](https://golangci-lint.run/)** - Lint avec métriques multiples
  ```bash
  golangci-lint run --enable-all
  ```

### Détection Duplication

- **[dupl](https://github.com/mibk/dupl)** - Détection duplication Go
  ```bash
  go install github.com/mibk/dupl@latest
  dupl -t 50 ./...
  ```

- **[jscpd](https://github.com/kucherenko/jscpd)** - Copy-paste detector multi-langages
  ```bash
  npm install -g jscpd
  jscpd --threshold 50 .
  ```

### Couverture Tests

- **Built-in Go** - go test -cover
  ```bash
  go test -coverprofile=coverage.out ./...
  go tool cover -html=coverage.out
  ```

- **[gocov](https://github.com/axw/gocov)** - Coverage JSON
  ```bash
  go install github.com/axw/gocov/gocov@latest
  gocov test ./... | gocov report
  ```

- **[go-carpet](https://github.com/msoap/go-carpet)** - Heatmap coverage
  ```bash
  go install github.com/msoap/go-carpet@latest
  go-carpet
  ```

### Visualisation & Rapports

- **[gopherbadger](https://github.com/jpoles1/gopherbadger)** - Badges coverage pour README
- **[octommander](https://github.com/alexandear/octommander)** - Stats git avancées
- **[SonarQube](https://www.sonarqube.org/)** - Plateforme qualité complète
- **[CodeClimate](https://codeclimate.com/)** - Analyse qualité continue

## Bonnes Pratiques

### Fréquence d'Analyse

- **Hebdomadaire** : Coverage tests (automatique en CI/CD)
- **Mensuel** : Stats globales, complexité, tendances
- **Trimestriel** : Analyse détaillée, refactoring planifié
- **Annuel** : Audit complet, revue architecture

### Seuils de Qualité Recommandés

| Métrique | Excellent | Bon | Acceptable | Problématique |
|----------|-----------|-----|------------|---------------|
| **Fichier** | < 300 lignes | < 500 | < 800 | > 800 |
| **Fonction** | < 30 lignes | < 50 | < 80 | > 100 |
| **Complexité** | < 5 | < 10 | < 15 | > 15 |
| **Commentaires** | > 20% | > 15% | > 10% | < 10% |
| **Coverage** | > 80% | > 70% | > 60% | < 60% |
| **Duplication** | < 3% | < 5% | < 10% | > 10% |

### Actions Basées sur Stats

| Situation | Action | Priorité |
|-----------|--------|----------|
| **Fichier > 800 lignes** | Découpage en modules | 🔴 Urgente |
| **Fonction > 100 lignes** | Refactoring | 🔴 Urgente |
| **Complexité > 15** | Simplification | 🔴 Urgente |
| **Coverage < 40%** | Ajouter tests | 🔴 Urgente |
| **Fichier 500-800 lignes** | Surveillance | ⚠️ Importante |
| **Fonction 50-100 lignes** | Révision | ⚠️ Importante |
| **Commentaires < 10%** | Documentation | ⚠️ Importante |
| **Duplication > 10%** | Factorisation | ⚠️ Importante |

### Intégration CI/CD

Ajouter vérifications automatiques dans pipeline :

```yaml
# .github/workflows/quality.yml
name: Code Quality

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Check complexity
        run: |
          go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
          gocyclo -over 15 . && echo "✅ Complexity OK" || exit 1
      
      - name: Check coverage
        run: |
          go test -coverprofile=coverage.out ./...
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 60" | bc -l) )); then
            echo "❌ Coverage $coverage% < 60%"
            exit 1
          fi
          echo "✅ Coverage $coverage% >= 60%"
      
      - name: Lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
```

## Anti-Patterns à Éviter

### ❌ Stats Sans Action

```
❌ Générer rapport et ne rien faire
✅ Stats → Analyse → Plan d'action → Exécution → Mesure
```

### ❌ Optimisation Prématurée

```
❌ Refactoriser systématiquement tout code > 50 lignes
✅ Prioriser code critique, fréquemment modifié, ou bugué
```

### ❌ Ignorer le Contexte

```
❌ "Ce fichier fait 800 lignes, il FAUT le découper"
✅ Analyser si découpage logique existe et apporte vraie valeur
```

### ❌ Inclure Code Généré

```
❌ Parser généré de 5000 lignes compte dans stats qualité
✅ Exclure systématiquement code généré des métriques qualité
```

### ❌ Coverage à 100% Aveugle

```
❌ Viser 100% coverage sans discernement
✅ Focus sur code critique, ignorer code trivial (getters, etc.)
```

### ❌ Complexité vs Clarté

```
❌ Réduire complexité cyclomatique en ajoutant abstractions obscures
✅ Privilégier code clair même si légèrement plus complexe
```

## Critères de Succès

### ✅ Rapport Complet

- [ ] Identification claire : code manuel vs généré vs tests
- [ ] Statistiques globales précises (lignes, fichiers, fonctions)
- [ ] Répartition par module détaillée
- [ ] Top 10 fichiers volumineux identifiés
- [ ] Top 15 fonctions volumineuses/complexes identifiées
- [ ] Métriques qualité (commentaires, complexité, longueur)
- [ ] Coverage tests avec analyse par package
- [ ] Code généré documenté séparément (sans recommandations)
- [ ] Tendances et évolution (si git disponible)

### ✅ Analyse Pertinente

- [ ] Code généré exclu des recommandations
- [ ] Tests analysés séparément avec métriques propres
- [ ] Seuils de qualité clairement définis
- [ ] Recommandations priorisées (Urgent/Important/Amélioration)
- [ ] Hotspots identifiés (fichiers fréquemment modifiés)
- [ ] Duplication détectée et quantifiée
- [ ] Actions concrètes et estimées en temps

### ✅ Rapport Actionnable

- [ ] Tableaux formatés et lisibles
- [ ] Graphiques ASCII pour visualisation
- [ ] Code couleur (✅ 🟢 ⚠️ 🔴)
- [ ] Liens vers prompts suggérés
- [ ] Estimations temps pour chaque action
- [ ] Priorisation claire (1/2/3)
- [ ] Contexte et justifications fournies

### ✅ Rapport Professionnel

- [ ] Format markdown structuré
- [ ] Résumé exécutif en tête
- [ ] Sections clairement séparées
- [ ] Commandes reproductibles fournies
- [ ] Références aux outils et ressources
- [ ] Date et commit pour traçabilité

## Ressources

### Documentation Officielle

- [Effective Go](https://go.dev/doc/effective_go) - Guide officiel Go
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Best practices
- [Go Testing](https://go.dev/doc/tutorial/add-a-test) - Guide tests officiels

### Articles & Guides

- [Cyclomatic Complexity Explained](https://en.wikipedia.org/wiki/Cyclomatic_complexity)
- [Test Coverage Best Practices](https://martinfowler.com/bliki/TestCoverage.html)
- [Code Metrics for Go](https://github.com/edgurgel/awesome-go#code-analysis)

### Outils Open Source

- [tokei](https://github.com/XAMPPRocky/tokei) - Compteur de lignes rapide
- [gocyclo](https://github.com/fzipp/gocyclo) - Complexité cyclomatique
- [dupl](https://github.com/mibk/dupl) - Détection duplication
- [golangci-lint](https://github.com/golangci/golangci-lint) - Linter complet

### Communauté

- [r/golang](https://reddit.com/r/golang) - Communauté Reddit Go
- [Gophers Slack](https://gophers.slack.com/) - Chat communauté Go
- [Go Forum](https://forum.golangbridge.org/) - Forum officiel

---

## Changelog

### Version 2.0 (Novembre 2025)
- ✅ Séparation code manuel / généré / tests
- ✅ Exclusion parser.go des statistiques qualité
- ✅ Ajout section couverture tests complète
- ✅ Ajout métriques qualité (duplication, hotspots)
- ✅ Ajout tendances et évolution git
- ✅ Recommandations détaillées avec estimations
- ✅ Commandes pour tous les outils
- ✅ Guide bonnes pratiques et anti-patterns

### Version 1.0 (Novembre 2025)
- Création prompt initial
- Stats basiques par module
- Top fichiers et fonctions
- Métriques complexité

---

**Version** : 2.0  
**Dernière mise à jour** : 26 novembre 2025  
**Mainteneur** : Équipe TSD