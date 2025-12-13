# 🔍 Revue RETE - Prompt 08: Pipeline et Validation

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/review.md](../../.github/prompts/review.md)

**Priorité:** Haute  
**Durée estimée:** 2-3 heures  
**Fichiers concernés:** ~6 fichiers (~2,000 lignes)  
**Date:** 2024-12-15

---

## 📋 Vue d'ensemble

Le module Pipeline est responsable de :
- L'ingestion et le parsing des fichiers de contraintes (.tsd)
- La validation des contraintes (syntaxe, sémantique, types)
- La construction du réseau RETE à partir des contraintes
- La gestion du flux complet (parsing → validation → construction)
- Le type checking et la cohérence sémantique
- Les métriques de cohérence

**⚠️ ATTENTION CRITIQUE:** Ce module contient `IngestFile` avec complexité ~48. Décomposition OBLIGATOIRE.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter obligatoirement :
- [⚠️ Standards Code Go](../../.github/prompts/common.md#standards-de-code-go) - Conventions, complexité, qualité
- [🎨 Conventions Nommage](../../.github/prompts/common.md#conventions-de-nommage) - Standards projet
- [📋 Checklist Commit](../../.github/prompts/common.md#checklist-avant-commit) - Validation
- [🔍 Revue Code](../../.github/prompts/review.md) - Process et techniques

---

## 🎯 Objectifs de cette revue

### 1. ⚠️ CRITIQUE: Décomposer IngestFile (complexité ~48)
- ✅ Identifier la fonction monolithique
- ✅ Décomposer en 5-7 sous-fonctions cohérentes (<12 chacune)
- ✅ Améliorer testabilité radicalement
- ✅ Target <10 pour chaque sous-fonction

### 2. Améliorer gestion erreurs pipeline
- ✅ Erreurs propagées avec contexte (fichier, ligne, colonne)
- ✅ Messages informatifs pour l'utilisateur
- ✅ Pas de suppression silencieuse
- ✅ Recovery sur panic parser

### 3. Valider tous cas edge
- ✅ Fichiers vides
- ✅ Fichiers malformés
- ✅ Types invalides
- ✅ Règles incohérentes
- ✅ Dépendances circulaires

### 4. Optimiser validation contraintes
- ✅ Validation incrémentale si possible
- ✅ Pas de revalidation inutile
- ✅ Court-circuit sur première erreur critique
- ✅ Validation parallèle si applicable

### 5. Clarifier flux pipeline
- ✅ Étapes claires et séquentielles
- ✅ Chaque étape testable indépendamment
- ✅ Rollback possible sur erreur
- ✅ Transactionnel si applicable

### 6. Garantir encapsulation et généricité
- ✅ Minimiser exports publics (privé par défaut)
- ✅ Éliminer tout hardcoding
- ✅ Rendre le code générique et réutilisable

---

## 📂 Périmètre des fichiers

```
rete/constraint_pipeline.go                 # ⚠️ IngestFile COMPLEXITÉ 48!
rete/constraint_pipeline_validator.go       # Validation contraintes
rete/constraint_pipeline_ingest.go          # Ingestion fichiers
rete/type_checker.go                        # Type checking
rete/coherence_mode.go                      # Mode cohérence
rete/coherence_metrics.go                   # Métriques cohérence
```

---

## ✅ Checklist détaillée

### 🏗️ Architecture et Design (Principes SOLID)

- [ ] **Single Responsibility Principle**
  - Pipeline → orchestration seulement
  - Validator → validation seulement
  - TypeChecker → types seulement
  - Parser → parsing seulement
  - Pas de "God Pipeline"

- [ ] **Open/Closed Principle**
  - Extensible sans modifier existant
  - Nouvelles étapes ajoutables
  - Interfaces pour abstraction

- [ ] **Liskov Substitution Principle**
  - Toutes implémentations respectent contrats
  - Pas de comportements surprenants

- [ ] **Interface Segregation Principle**
  - Interfaces petites et focalisées
  - Pas d'interface monolithique
  - Clients dépendent du minimum

- [ ] **Dependency Inversion Principle**
  - Dépendances sur interfaces
  - Injection de dépendances (parser, validator, builder)
  - Pas de dépendances hardcodées

### 🔒 Encapsulation et Visibilité

- [ ] **Variables et fonctions privées par défaut**
  - Tous symboles privés sauf nécessité absolue
  - Exports publics justifiés et documentés
  - Étapes internes pipeline cachées

- [ ] **Minimiser exports publics**
  - Interface Pipeline exportée
  - Méthode IngestFile exportée
  - Implémentation étapes privée
  - Types internes privés

- [ ] **Contrats d'interface respectés**
  - API publique stable
  - Breaking changes documentés

### 🚫 Anti-Hardcoding (CRITIQUE)

- [ ] **Aucune valeur hardcodée**
  - Pas de magic numbers
  - Pas de magic strings
  - Pas de chemins hardcodés
  - Pas de limites hardcodées (taille fichier, nombre règles, etc.)

- [ ] **Constantes nommées et explicites**
  ```go
  // ❌ MAUVAIS
  if fileSize > 1048576 { return errTooBig }
  if len(rules) > 1000 { return errTooMany }
  
  // ✅ BON
  const (
      MaxFileSize      = 1 * 1024 * 1024  // 1 MB
      MaxRulesPerFile  = 1000
  )
  if fileSize > MaxFileSize { 
      return fmt.Errorf("file too large: %d > %d bytes", fileSize, MaxFileSize)
  }
  if len(rules) > MaxRulesPerFile { return errTooMany }
  ```

- [ ] **Code générique et paramétrable**
  - Pas de code spécifique à un fichier
  - Configuration via options
  - Extensible

### 🧪 Tests Fonctionnels RÉELS (CRITIQUE)

- [ ] **Pas de simulation/mocks**
  - Tests ingèrent vraiment des fichiers .tsd
  - Parsing et validation réels
  - Vérification réseau construit
  - Pas de suppositions

- [ ] **Tests déterministes et isolés**
  - Chaque test indépendant
  - Setup/teardown propre (fichiers temporaires nettoyés)
  - Reproductibles

- [ ] **Couverture > 85%**
  - Cas nominaux
  - Cas d'erreur (parsing, validation)
  - Edge cases (fichiers vides, malformés, très gros)
  - Tous chemins d'erreur

- [ ] **Tests par étape**
  - Tests parsing isolé
  - Tests validation isolée
  - Tests type checking isolé
  - Tests intégration pipeline complet

### 📋 Qualité du Code

- [ ] **Complexité cyclomatique < 15**
  - ⚠️ CRITIQUE: Décomposer IngestFile (48 → <10)
  - Toutes autres fonctions <15 (idéalement <10)
  - Extract Function massivement

- [ ] **Fonctions < 50 lignes**
  - Sauf justification documentée
  - IngestFile doit être <30 lignes après décomposition
  - Une fonction = une étape claire

- [ ] **Imbrication < 4 niveaux**
  - Pas de deep nesting
  - Early return
  - Extract Function

- [ ] **Pas de duplication (DRY)**
  - Patterns communs extraits
  - Helpers partagés
  - Constantes pour valeurs répétées

- [ ] **Noms explicites et idiomatiques**
  - Variables: camelCase descriptif
  - Fonctions: MixedCaps, verbes (parseFile, validateTypes)
  - Types: MixedCaps, noms (ConstraintPipeline, TypeChecker)
  - Constantes: MixedCaps ou UPPER_CASE
  - Pas d'abréviations obscures

- [ ] **Code auto-documenté**
  - Flux pipeline clair à la lecture
  - Logique évidente
  - Commentaires si algorithme complexe

### 🔐 Sécurité et Robustesse

- [ ] **Validation des entrées**
  - Fichiers nulles/vides gérées
  - Chemins validés (pas de path traversal)
  - Taille fichier limitée (DoS)
  - Encoding validé
  - Pas de panic

- [ ] **Gestion d'erreurs robuste**
  - Erreurs propagées avec contexte (fichier, ligne, colonne)
  - Messages informatifs pour utilisateur
  - Pas de suppression silencieuse
  - Return early on error
  - Wrap errors avec position

- [ ] **Recovery sur panic parser**
  - Panic parser catchée
  - Convertie en erreur avec position
  - Logged avec contexte
  - Exécution arrêtée proprement

- [ ] **Validation tous cas edge**
  - Fichier vide
  - Fichier très gros (>100 MB)
  - Fichier malformé (syntaxe invalide)
  - Types inconsistants
  - Dépendances circulaires
  - Règles sans patterns
  - Patterns sans conditions

- [ ] **Ressources libérées proprement**
  - Fichiers fermés (defer)
  - Pas de fuites mémoire
  - Cleanup sur erreur
  - Context pour timeout si long

### 📚 Documentation

- [ ] **En-tête copyright présent**
  ```go
  // Copyright (c) 2025 TSD Contributors
  // Licensed under the MIT License
  // See LICENSE file in the project root for full license text
  ```

- [ ] **GoDoc pour tous exports**
  - Pipeline documenté
  - IngestFile documenté avec exemples
  - Erreurs possibles documentées
  - Limites documentées (taille max, etc.)

- [ ] **Commentaires inline si complexe**
  - Étapes pipeline expliquées
  - Justification validations
  - Références si algorithmes connus

- [ ] **Pas de commentaires obsolètes**
  - Supprimer code commenté
  - MAJ après changements
  - Pas de redondance

### ⚡ Performance

- [ ] **Parsing efficace**
  - Pas de reparsing inutile
  - Streaming si gros fichiers
  - Bufferisation appropriée

- [ ] **Validation efficace**
  - Validation incrémentale si possible
  - Court-circuit sur erreur critique
  - Pas de validation redondante

- [ ] **Construction efficace**
  - Pas de reconstruction inutile
  - Réutilisation composants (sharing)
  - Allocations minimisées

- [ ] **Type checking efficace**
  - Cache résultats si répétés
  - Pas de recalculs

### 🎨 Pipeline (Spécifique)

- [ ] **Flux clair et séquentiel**
  ```
  Fichier → Parse → Validate → TypeCheck → Build → Done
  ```

- [ ] **Chaque étape indépendante et testable**
  - parseConstraintFile()
  - validateConstraints()
  - checkTypes()
  - buildNetwork()

- [ ] **Rollback possible sur erreur**
  - État avant pipeline restaurable
  - Ou transaction atomique
  - Ou fail fast sans effets

- [ ] **Erreurs avec position**
  - Fichier, ligne, colonne
  - Contexte (quelle règle, quel pattern)
  - Message clair pour utilisateur

- [ ] **Métriques cohérence**
  - Collectées sans ralentir
  - Optionnelles
  - Bien documentées

---

## 🚫 Anti-Patterns à Détecter et Éliminer

- [ ] **God Object** - Pipeline qui fait tout
  - ⚠️ IngestFile = God Method
  - Décomposer massivement
  - Délégation

- [ ] **Long Method** - IngestFile >100-200 lignes
  - ⚠️ CRITIQUE: décomposer
  - Extract Function
  - Orchestration simple

- [ ] **Long Parameter List** - >5 paramètres
  - Utiliser Options/Config
  - Grouper paramètres

- [ ] **Magic Numbers/Strings** - Hardcoding
  - Extract Constant
  - Constantes nommées

- [ ] **Duplicate Code** - Validation répétée
  - Extract Function
  - Helpers

- [ ] **Dead Code** - Code inutilisé
  - Supprimer

- [ ] **Deep Nesting** - >4 niveaux
  - Early return
  - Extract Function

- [ ] **Error Swallowing** - Erreurs ignorées
  - Propager avec contexte
  -Logger minimum

---

## 🔧 Commandes de validation

### Tests

```bash
# Tests pipeline
go test -v ./rete -run "TestPipeline"
go test -v ./rete -run "TestIngest"

# Tests validation
go test -v ./rete -run "TestValidat"
go test -v ./rete -run "TestTypeCheck"

# Tests cohérence
go test -v ./rete -run "TestCoherence"

# Tous tests avec couverture
go test -coverprofile=coverage_pipeline.out ./rete -run "TestPipeline|TestIngest|TestValidat|TestCoherence"
go tool cover -func=coverage_pipeline.out
go tool cover -html=coverage_pipeline.out -o coverage_pipeline.html

# Race detector
go test -race ./rete -run "TestPipeline"
```

### Performance

```bash
# Benchmarks pipeline
go test -bench=BenchmarkIngest -benchmem ./rete
go test -bench=BenchmarkPipeline -benchmem ./rete

# Profiling
go test -bench=BenchmarkIngest -cpuprofile=cpu_pipeline.prof ./rete
go tool pprof -http=:8080 cpu_pipeline.prof
```

### Qualité

```bash
# Complexité (CRITIQUE: trouver IngestFile à 48)
gocyclo -over 15 rete/constraint*.go rete/type_checker.go rete/coherence*.go
gocyclo -top 20 rete/constraint*.go

# Vérifications statiques
go vet ./rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go
staticcheck ./rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go
errcheck ./rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go
gosec ./rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go

# Formatage
gofmt -l rete/constraint*.go rete/type_checker.go rete/coherence*.go
go fmt ./rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go
goimports -w rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go

# Linting
golangci-lint run ./rete/constraint*.go ./rete/type_checker.go ./rete/coherence*.go

# Validation complète
make validate
```

### Vérification Copyright

```bash
for file in rete/constraint*.go rete/type_checker.go rete/coherence*.go; do
    if [ -f "$file" ] && ! head -1 "$file" | grep -q "Copyright"; then
        echo "⚠️  COPYRIGHT MANQUANT: $file"
    fi
done
```

---

## 🔄 Processus de revue recommandé

### Phase 1: Analyse initiale (30-45 min)

1. **Mesurer baseline**
   ```bash
   mkdir -p REPORTS/review-rete
   
   # Complexité (TROUVER IngestFile à 48!)
   gocyclo -over 10 rete/constraint*.go > REPORTS/review-rete/pipeline_complexity_before.txt
   echo "=== TOP COMPLEXITÉ ==="
   gocyclo -top 20 rete/constraint*.go rete/type_checker.go
   
   # Couverture
   go test -coverprofile=REPORTS/review-rete/pipeline_coverage_before.out ./rete -run "TestPipeline|TestIngest" 2>/dev/null
   go tool cover -func=REPORTS/review-rete/pipeline_coverage_before.out > REPORTS/review-rete/pipeline_coverage_before.txt
   
   # Benchmarks
   go test -bench=BenchmarkIngest -benchmem ./rete > REPORTS/review-rete/pipeline_benchmarks_before.txt 2>&1
   ```

2. **Lire fichiers dans ordre logique**
   - `coherence_mode.go` (modes)
   - `coherence_metrics.go` (métriques)
   - `type_checker.go` (type checking)
   - `constraint_pipeline_validator.go` (validation)
   - `constraint_pipeline_ingest.go` (ingestion)
   - `constraint_pipeline.go` (⚠️ INGESTFILE!)

3. **Pour chaque fichier, vérifier**
   - [ ] Copyright présent?
   - [ ] Exports minimaux?
   - [ ] Aucun hardcoding?
   - [ ] Code générique?
   - [ ] Complexité <15? (⚠️ identifier IngestFile)
   - [ ] Gestion erreurs avec position?
   - [ ] Tests présents?
   - [ ] GoDoc complet?

### Phase 2: Identification des problèmes (30-45 min)

**Créer liste priorisée dans** `REPORTS/review-rete/08_pipeline_issues.md`:

```markdown
# Problèmes Identifiés - Pipeline et Validation

## P0 - BLOQUANT

### 1. [Si bugs détectés]
- **Fichier:** constraint_pipeline.go:XXX
- **Type:** Validation incorrecte / Parsing échoue
- **Impact:** Réseau incorrect ou crash
- **Solution:** ...

## P1 - IMPORTANT

### 1. Complexité 48 dans constraint_pipeline.go
- **Fichier:** constraint_pipeline.go:XXX
- **Fonction:** `IngestFile` (ou similaire)
- **Complexité:** 48
- **Impact:** Impossible à maintenir, tester, comprendre
- **Solution:** Extract Function - décomposer en 6-8 étapes
- **Cible:** Max 10 par fonction

### 2. Gestion erreurs sans position
- **Fichiers:** Multiples
- **Type:** Erreurs sans contexte (ligne, colonne)
- **Impact:** Débogage impossible pour utilisateur
- **Solution:** Wrap errors avec position

### 3. Hardcoding limites
- **Fichiers:** Multiples
- **Type:** Magic numbers (taille max, nombre règles)
- **Impact:** Pas configurable
- **Solution:** Extract Constant

## P2 - SOUHAITABLE
...
```

**Problèmes à chercher:**

**P0:**
- Validation incorrecte (règles invalides acceptées)
- Parsing échoue sur fichiers valides
- Type checking bugué
- Panic non catchée
- Path traversal possible

**P1:**
- **Complexité 48 IngestFile (PRIORITÉ ABSOLUE)**
- Erreurs sans position fichier
- Hardcoding limites/tailles
- Exports non justifiés
- Couverture <70%
- Missing copyright

**P2:**
- Complexité 10-15
- Optimisations mineures
- Refactoring clarté

### Phase 3: Corrections (75-90 min)

#### 3.1 Fixer P0 (bloquants)

**Exemple: Validation incorrecte**

```go
// AVANT - règle invalide acceptée
func validateRule(rule *ast.Rule) error {
    if rule.Name == "" {
        return errors.New("empty name")
    }
    // ❌ Ne valide pas patterns
    return nil
}

// APRÈS - validation complète
func validateRule(rule *ast.Rule, pos ast.Position) error {
    if rule.Name == "" {
        return fmt.Errorf("%s: rule must have a name", pos)
    }
    if len(rule.Patterns) == 0 {
        return fmt.Errorf("%s: rule %s has no patterns", pos, rule.Name)
    }
    for i, pattern := range rule.Patterns {
        if err := validatePattern(pattern, pos); err != nil {
            return fmt.Errorf("%s: rule %s pattern %d: %w", pos, rule.Name, i, err)
        }
    }
    return nil
}
```

**Commit:**
```bash
git commit -m "[Review-08/Pipeline] fix(P0): validation complète règles avec position

- Valide présence patterns
- Valide chaque pattern
- Erreurs avec position fichier
- Messages informatifs
- Tests edge cases ajoutés

Resolves: P0-pipeline-validation-incomplete
Refs: scripts/review-rete/08_pipeline.md"
```

#### 3.2 Décomposer IngestFile (P1 PRIORITÉ ABSOLUE)

**Identifier la fonction:**
```bash
gocyclo -over 40 rete/constraint*.go
# Probablement IngestFile dans constraint_pipeline.go
```

**Pattern de décomposition:**

```go
// AVANT - complexité 48, ~200 lignes
func (p *ConstraintPipeline) IngestFile(filepath string) error {
    // 30 lignes lecture fichier
    // 40 lignes parsing
    // 50 lignes extraction (types, rules, facts)
    // 40 lignes validation
    // 40 lignes construction réseau
}

// APRÈS - décomposer en étapes claires

func (p *ConstraintPipeline) IngestFile(filepath string) error {
    // Orchestration simple - complexité ~8
    
    // Étape 1: Lecture fichier
    content, err := p.readConstraintFile(filepath)
    if err != nil {
        return fmt.Errorf("read file %s: %w", filepath, err)
    }
    
    // Étape 2: Parsing
    parsed, err := p.parseConstraints(content, filepath)
    if err != nil {
        return fmt.Errorf("parse file %s: %w", filepath, err)
    }
    
    // Étape 3: Extraction composants
    types, rules, facts := p.extractComponents(parsed)
    
    // Étape 4: Validation
    if err := p.validateConstraints(types, rules, facts, filepath); err != nil {
        return fmt.Errorf("validate file %s: %w", filepath, err)
    }
    
    // Étape 5: Type checking
    if err := p.checkTypes(types, rules, filepath); err != nil {
        return fmt.Errorf("type check file %s: %w", filepath, err)
    }
    
    // Étape 6: Construction réseau
    if err := p.buildNetworkFromConstraints(types, rules, facts); err != nil {
        return fmt.Errorf("build network from %s: %w", filepath, err)
    }
    
    // Étape 7: Ingestion facts initiaux
    if err := p.ingestInitialFacts(facts); err != nil {
        return fmt.Errorf("ingest facts from %s: %w", filepath, err)
    }
    
    return nil
}

// Chaque sous-fonction <10 complexité

func (p *ConstraintPipeline) readConstraintFile(filepath string) ([]byte, error) {
    // Complexité ~6
    // Validation path, lecture, limite taille
}

func (p *ConstraintPipeline) parseConstraints(content []byte, filepath string) (*ast.File, error) {
    // Complexité ~8
    // Parsing avec recovery panic
    // Erreurs avec position
}

func (p *ConstraintPipeline) extractComponents(parsed *ast.File) (types, rules, facts) {
    // Complexité ~7
    // Extraction types, rules, facts du AST
}

func (p *ConstraintPipeline) validateConstraints(types, rules, facts, filepath) error {
    // Complexité ~9
    // Validation sémantique complète
}

func (p *ConstraintPipeline) checkTypes(types, rules, filepath) error {
    // Complexité ~8
    // Type checking complet
}

func (p *ConstraintPipeline) buildNetworkFromConstraints(types, rules, facts) error {
    // Complexité ~9
    // Construction réseau RETE
}

func (p *ConstraintPipeline) ingestInitialFacts(facts []Fact) error {
    // Complexité ~6
    // Insertion facts initiaux
}
```

**Tests pour chaque sous-fonction:**
```go
func TestReadConstraintFile(t *testing.T) { /* ... */ }
func TestParseConstraints(t *testing.T) { /* ... */ }
func TestExtractComponents(t *testing.T) { /* ... */ }
func TestValidateConstraints(t *testing.T) { /* ... */ }
func TestCheckTypes(t *testing.T) { /* ... */ }
func TestBuildNetworkFromConstraints(t *testing.T) { /* ... */ }
func TestIngestInitialFacts(t *testing.T) { /* ... */ }

// Tests intégration
func TestIngestFile_Complete(t *testing.T) { /* ... */ }
```

**Commit:**
```bash
git commit -m "[Review-08/Pipeline] refactor(P1): décompose IngestFile (48→8)

- Extrait readConstraintFile() (6)
- Extrait parseConstraints() (8)
- Extrait extractComponents() (7)
- Extrait validateConstraints() (9)
- Extrait checkTypes() (8)
- Extrait buildNetworkFromConstraints() (9)
- Extrait ingestInitialFacts() (6)
- IngestFile orchestration: 8
- 7 tests unitaires ajoutés
- 1 test intégration maintenu
- Améliore testabilité radicalement

Resolves: P1-pipeline-complexity-48
Refs: scripts/review-rete/08_pipeline.md"
```

#### 3.3 Améliorer erreurs avec position (P1)

```go
// AVANT - pas de position
return errors.New("invalid type")

// APRÈS - position complète
return fmt.Errorf("%s:%d:%d: invalid type '%s': %w", 
    filepath, line, column, typeName, underlyingErr)
```

#### 3.4 Éliminer hardcoding (P1)

```go
// AVANT
if fileSize > 10485760 { return errTooBig }  // 10 MB
if len(rules) > 500 { return errTooMany }

// APRÈS
const (
    MaxConstraintFileSize = 10 * 1024 * 1024  // 10 MB
    MaxRulesPerFile       = 500
)

if fileSize > MaxConstraintFileSize {
    return fmt.Errorf("file %s too large: %d > %d bytes", 
        filepath, fileSize, MaxConstraintFileSize)
}
if len(rules) > MaxRulesPerFile { return errTooMany }
```

### Phase 4: Validation finale (15-30 min)

```bash
#!/bin/bash
echo "=== VALIDATION FINALE PIPELINE ==="

# 1. Tests
echo "🧪 Tests..."
go test -v ./rete -run "TestPipeline|TestIngest|TestValidat"
TESTS=$?

# 2. Race detector
echo "🏁 Race detector..."
go test -race ./rete -run "TestPipeline"
RACE=$?

# 3. Complexité (CRITIQUE: IngestFile doit être <10)
echo "📊 Complexité..."
COMPLEX=$(gocyclo -over 15 rete/constraint*.go rete/type_checker.go | wc -l)
INGESTFILE_COMPLEX=$(gocyclo rete/constraint_pipeline.go | grep IngestFile | awk '{print $1}')

# 4. Couverture
echo "📈 Couverture..."
go test -coverprofile=pipeline_final.out ./rete -run "TestPipeline|TestIngest" 2>/dev/null
COVERAGE=$(go tool cover -func=pipeline_final.out | tail -1 | awk '{print $3}' | sed 's/%//')

# 5. Copyright
echo "©️  Copyright..."
MISSING_COPYRIGHT=0
for file in rete/constraint*.go rete/type_checker.go rete/coherence*.go; do
    if [ -f "$file" ] && ! head -1 "$file" | grep -q "Copyright"; then
        MISSING_COPYRIGHT=$((MISSING_COPYRIGHT + 1))
        echo "  ⚠️  $file"
    fi
done

# 6. Validation
echo "✅ Validation..."
make validate
VALIDATE=$?

# Résumé
echo ""
echo "=== RÉSULTATS ==="
[ $TESTS -eq 0 ] && echo "✅ Tests: PASS" || echo "❌ Tests: FAIL"
[ $RACE -eq 0 ] && echo "✅ Race: PASS" || echo "❌ Race: FAIL"
[ $COMPLEX -eq 0 ] && echo "✅ Complexité: OK (0 >15)" || echo "❌ Complexité: $COMPLEX >15"
[ -n "$INGESTFILE_COMPLEX" ] && echo "  IngestFile: $INGESTFILE_COMPLEX" || echo "  IngestFile: N/A"
[ $(echo "$COVERAGE >= 85" | bc -l) -eq 1 ] && echo "✅ Couverture: $COVERAGE%" || echo "❌ Couverture: $COVERAGE%"
[ $MISSING_COPYRIGHT -eq 0 ] && echo "✅ Copyright: OK" || echo "❌ Copyright: $MISSING_COPYRIGHT manquants"
[ $VALIDATE -eq 0 ] && echo "✅ Validation: PASS" || echo "❌ Validation: FAIL"

# Verdict
if [ $TESTS -eq 0 ] && [ $RACE -eq 0 ] && [ $COMPLEX -eq 0 ] && [ $(echo "$COVERAGE >= 85" | bc -l) -eq 1 ] && [ $MISSING_COPYRIGHT -eq 0 ] && [ $VALIDATE -eq 0 ]; then
    echo ""
    echo "🎉 VALIDATION RÉUSSIE - Prêt pour Prompt 09!"
    exit 0
else
    echo ""
    echo "❌ VALIDATION ÉCHOUÉE"
    exit 1
fi
```

---

## 📝 Livrables attendus

### 1. Rapport d'analyse

**Créer:** `REPORTS/review-rete/08_pipeline_report.md`

**Structure obligatoire:**

```markdown
# 🔍 Revue de Code : Pipeline et Validation

**Date:** 2024-12-XX  
**Réviseur:** [Nom]  
**Durée:** Xh Ym

---

## 📊 Vue d'Ensemble

- **Fichiers analysés:** 6
- **Lignes de code:** ~2,000
- **Complexité IngestFile avant:** 48
- **Complexité IngestFile après:** <10
- **Couverture avant:** X%
- **Couverture après:** Y%

---

## ✅ Points Forts

- Pipeline identifié clairement
- Validation présente
- ...

---

## ❌ Problèmes Identifiés et Corrigés

### P0 - BLOQUANT

#### 1. [Si applicable]
...

### P1 - IMPORTANT

#### 1. Complexité 48 dans IngestFile
- **Fonction:** IngestFile
- **Avant:** 48, ~200 lignes
- **Après:** 8, ~40 lignes
- **Décomposition:** 7 sous-fonctions (<10 chacune)
- **Tests:** 7 tests unitaires + 1 intégration
- **Commit:** abc1234

#### 2. Erreurs sans position
- **Solution:** Position ajoutée à toutes erreurs
- **Format:** fichier:ligne:colonne: message
- **Commit:** def5678

#### 3. Hardcoding limites
- **Constantes créées:** 6
- **Commit:** ghi9012

---

## 🔧 Changements Apportés

### Refactoring

1. **Décomposition IngestFile**
   - 1 fonction monolithique → 7 fonctions claires
   - Complexité 48 → max 9
   - Tests unitaires: 7
   - Testabilité: radicalement améliorée

2. **Erreurs avec position**
   - Toutes erreurs incluent fichier:ligne:colonne
   - Messages utilisateur informatifs
   - Contexte complet

3. **Constantes nommées**
   - 10 magic numbers → constantes
   - 4 magic strings → constantes

---

## 📈 Métriques Avant/Après

| Métrique | Avant | Après | Amélioration |
|----------|-------|-------|--------------|
| Complexité IngestFile | 48 | 8 | ✅ -83% |
| Fonctions >15 | 3 | 0 | ✅ 100% |
| Couverture | 67% | 89% | ✅ +22% |
| Magic numbers | 10 | 0 | ✅ 100% |

---

## 💡 Recommandations Futures

### Court terme
1. Tests fichiers .tsd réels variés
2. Benchmarks sur gros fichiers
3. Documentation format .tsd

### Moyen terme
1. Validation incrémentale
2. Parsing parallèle si multi-fichiers
3. Cache parsing AST

---

## 🏁 Verdict

✅ **APPROUVÉ**

IngestFile décomposé, flux clair, erreurs informatives, standards respectés.
Prêt pour Prompt 09 (Métriques).

---

**Prochaines étapes:**
1. Merge commits
2. Lancer Prompt 09
3. Documenter format .tsd
```

### 2. Commits atomiques

**Format:**
```
[Review-08/Pipeline] <type>(scope): <description courte>

- Détail 1
- Détail 2
- Resolves: <issue>

Refs: scripts/review-rete/08_pipeline.md
```

---

## 📊 Métriques de succès

| Métrique | Valeur Actuelle | Cible | Critique |
|----------|----------------|-------|----------|
| Complexité IngestFile | À mesurer (48?) | <10 | ⚠️ OUI! |
| Fonctions >15 | À mesurer | 0 | ⚠️ Oui |
| Couverture tests | À mesurer | >85% | Oui |
| Erreurs avec position | À vérifier | 100% | ⚠️ Oui |
| Exports publics | À mesurer | Minimal | Oui |
| Magic numbers | À mesurer | 0 | Oui |
| Copyright | À mesurer | 100% | Oui |

---

## 🎓 Ressources et références

### Standards Projet
- [common.md](../../.github/prompts/common.md)
- [review.md](../../.github/prompts/review.md)
- [Makefile](../../Makefile)

### Parsing & Validation
- AST (Abstract Syntax Tree)
- Semantic analysis
- Type systems

### Error Handling
- Error positions in compilers
- User-friendly error messages

---

## ✅ Checklist finale avant Prompt 09

**Validation technique:**
- [ ] Tous tests pipeline passent
- [ ] Race detector clean
- [ ] IngestFile complexité <10 (CRITIQUE!)
- [ ] Aucune autre fonction >15
- [ ] Couverture >85%
- [ ] `make validate` passe

**Qualité code:**
- [ ] Aucun hardcoding
- [ ] Code générique
- [ ] Exports minimaux
- [ ] Flux pipeline clair (7-8 étapes)
- [ ] Chaque étape testable indépendamment
- [ ] Pas de duplication

**Robustesse:**
- [ ] Erreurs avec position (fichier:ligne:colonne)
- [ ] Validation edge cases
- [ ] Recovery panic parser
- [ ] Ressources libérées (defer)

**Tests:**
- [ ] Tests par étape
- [ ] Tests intégration
- [ ] Tests fichiers réels .tsd
- [ ] Tests fichiers malformés

**Documentation:**
- [ ] Copyright 100%
- [ ] GoDoc complet
- [ ] Flux documenté
- [ ] Limites documentées

---

## 🚀 Script d'analyse rapide

```bash
#!/bin/bash
# scripts/review-rete/analyze_pipeline.sh

set -e
echo "=== ANALYSE PIPELINE ==="
echo ""

mkdir -p REPORTS/review-rete

# Baseline
echo "📊 Mesure baseline..."
gocyclo -over 10 rete/constraint*.go rete/type_checker.go > REPORTS/review-rete/pipeline_complexity_before.txt
go test -coverprofile=REPORTS/review-rete/pipeline_coverage_before.out ./rete -run "TestPipeline|TestIngest" 2>/dev/null
go tool cover -func=REPORTS/review-rete/pipeline_coverage_before.out > REPORTS/review-rete/pipeline_coverage_before.txt

echo "✅ Baseline sauvegardée"
echo ""

# CRITIQUE: Trouver IngestFile
echo "🚨 RECHERCHE IngestFile (complexité 48?)..."
gocyclo rete/constraint_pipeline.go | grep -i ingest || echo "  (IngestFile non trouvé)"
echo ""

# TOP complexité
echo "📈 TOP COMPLEXITÉ..."
gocyclo -top 20 rete/constraint*.go rete/type_checker.go | head -15
echo ""

# Copyright
echo "©️  COPYRIGHT..."
MISSING=0
for file in rete/constraint*.go rete/type_checker.go rete/coherence*.go; do
    if [ -f "$file" ] && ! head -1 "$file" | grep -q "Copyright"; then
        echo "  ❌ $file"
        MISSING=$((MISSING + 1))
    fi
done
[ $MISSING -eq 0 ] && echo "  ✓ OK"

echo ""
echo "=== Analyse terminée ==="
echo "Créer REPORTS/review-rete/08_pipeline_issues.md"
```

**Lancer:**
```bash
chmod +x scripts/review-rete/analyze_pipeline.sh
./scripts/review-rete/analyze_pipeline.sh
```

---

**⚠️ PRIORITÉ ABSOLUE:** Décomposer IngestFile (48 → <10). C'est le refactoring le plus critique du prompt 08.

**Prêt à commencer?** 🚀

Bonne revue! Respecter scrupuleusement les standards common.md et review.md.