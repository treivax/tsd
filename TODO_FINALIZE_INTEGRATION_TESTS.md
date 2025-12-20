# TODO - Finalisation Tests d'Intégration et E2E

Date: 2025-12-19
Priorité: HAUTE
Statut: ⚠️ À Compléter

---

## 🎯 Contexte

Le refactoring du système de gestion des IDs est **complété et validé**.
Les **tests d'intégration et E2E** doivent maintenant être finalisés.

**Travail accompli**:
- ✅ Revue de code complète
- ✅ Refactoring avec réduction 60% complexité
- ✅ Tous tests existants passent
- ✅ Fixtures TSD créées
- ✅ Exemples créés

**Travail restant**:
- ⚠️ Tests d'intégration à finaliser
- ⚠️ Tests E2E à créer
- ⚠️ Benchmarks à ajouter

---

## 📋 Actions Immédiate (2-3h)

### 1. Vérifier API RETE

**Objectif**: Identifier les fonctions exactes pour les tests

```bash
# Rechercher les fonctions de parsing
cd /home/resinsec/dev/tsd
grep -r "func.*Parse" constraint/*.go | grep -v test | grep "Program"

# Rechercher les fonctions RETE
grep -r "func.*NewNetwork" rete/*.go
grep -r "func.*Compile" rete/*.go | head -20
grep -r "func.*Assert" rete/*.go | head -10

# Trouver exemples de tests existants
grep -r "NewNetwork\|CompileConstraints\|AssertFact" tests/ --include="*.go" | head -20
```

**Notes**:
- `constraint.ParseConstraint("inline", []byte(input))` retourne interface{}
- Convertir avec `constraint.ConvertResultToProgram(parsedResult)` → Program
- Valider avec `constraint.ValidateProgram(parsedResult)`

### 2. Créer Tests d'Intégration

**Fichier**: `tests/integration/fact_lifecycle_test.go`

**Template de base** (à adapter selon API):

```go
//go:build integration

package integration

import (
    "testing"
    "github.com/treivax/tsd/constraint"
    "github.com/treivax/tsd/rete"
)

func TestFactLifecycle_Complete(t *testing.T) {
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string)
        
        alice = User("Alice", 30)
        Login(alice, "alice@example.com", "pass")
        
        {u: User, l: Login} / l.user == u ==> Log("test")
    `
    
    // 1. Parser
    parsedResult, err := constraint.ParseConstraint("inline", []byte(input))
    if err != nil {
        t.Fatalf("Parsing: %v", err)
    }
    
    // 2. Convertir
    program, err := constraint.ConvertResultToProgram(parsedResult)
    if err != nil {
        t.Fatalf("Conversion: %v", err)
    }
    
    // 3. Valider
    err = constraint.ValidateProgram(parsedResult)
    if err != nil {
        t.Fatalf("Validation: %v", err)
    }
    
    // 4. Convertir faits RETE
    reteFacts, err := constraint.ConvertFactsToReteFormat(program)
    if err != nil {
        t.Fatalf("Conversion RETE: %v", err)
    }
    
    // 5. Créer réseau RETE (ADAPTER ICI)
    network := rete.NewNetwork()
    
    // Compiler règles - VÉRIFIER LA MÉTHODE EXACTE
    // Possibilités: CompileConstraints, CompileExpression, AddRule...
    for _, expr := range program.Expressions {
        err := network.CompileConstraints(program, expr) // À VÉRIFIER
        if err != nil {
            t.Fatalf("Compilation: %v", err)
        }
    }
    
    // Asserter faits
    for _, fact := range reteFacts {
        network.AssertFact(fact) // À VÉRIFIER
    }
    
    // Vérifier résultats
    // À ADAPTER selon API RETE
    if len(network.TerminalNodes) == 0 {
        t.Error("Aucun TerminalNode")
    }
}
```

**Commandes pour développer**:
```bash
# Copier template
cat > tests/integration/fact_lifecycle_test.go << 'EOF'
[Coller template ci-dessus]
EOF

# Tester compilation
go test -tags=integration ./tests/integration/fact_lifecycle_test.go -v -run TestFactLifecycle_Complete
```

### 3. Créer Tests E2E

**Fichier**: `tests/e2e/user_scenarios_test.go`

**Contenu**:
```go
//go:build e2e

package e2e

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/treivax/tsd/constraint"
    "github.com/treivax/tsd/rete"
)

func TestE2E_UserLoginScenario(t *testing.T) {
    // Lire fichier TSD
    tsdFile := filepath.Join("testdata", "user_login.tsd")
    content, err := os.ReadFile(tsdFile)
    if err != nil {
        t.Fatalf("Lecture fichier: %v", err)
    }
    
    // Parser
    parsedResult, err := constraint.ParseConstraint(tsdFile, content)
    if err != nil {
        t.Fatalf("Parsing: %v", err)
    }
    
    // Valider
    err = constraint.ValidateProgram(parsedResult)
    if err != nil {
        t.Fatalf("Validation: %v", err)
    }
    
    program, _ := constraint.ConvertResultToProgram(parsedResult)
    reteFacts, _ := constraint.ConvertFactsToReteFormat(program)
    
    // Créer et exécuter réseau
    network := rete.NewNetwork()
    // ... (même logique que tests d'intégration)
    
    // Vérifier résultats attendus
    // (selon le fichier user_login.tsd)
}

func TestE2E_AllExamples(t *testing.T) {
    examplesDir := "../../examples"
    files, err := filepath.Glob(filepath.Join(examplesDir, "*.tsd"))
    if err != nil {
        t.Fatalf("Recherche: %v", err)
    }
    
    for _, file := range files {
        t.Run(filepath.Base(file), func(t *testing.T) {
            content, err := os.ReadFile(file)
            if err != nil {
                t.Fatalf("Lecture: %v", err)
            }
            
            parsedResult, err := constraint.ParseConstraint(file, content)
            if err != nil {
                t.Fatalf("Parsing: %v", err)
            }
            
            err = constraint.ValidateProgram(parsedResult)
            if err != nil {
                t.Fatalf("Validation: %v", err)
            }
            
            t.Logf("✅ %s validé", filepath.Base(file))
        })
    }
}
```

---

## 📋 Actions Court Terme (2-3h)

### 4. Créer Benchmarks

**Fichier**: `tests/performance/id_generation_benchmark_test.go`

```go
package performance

import (
    "testing"
    "github.com/treivax/tsd/constraint"
)

func BenchmarkFactGeneration(b *testing.B) {
    typeDef := constraint.TypeDefinition{
        Name: "User",
        Fields: []constraint.Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
            {Name: "age", Type: "number"},
        },
    }
    
    fact := constraint.Fact{
        TypeName: "User",
        Fields: []constraint.FactField{
            {Name: "name", Value: constraint.FactValue{Type: "string", Value: "Alice"}},
            {Name: "age", Value: constraint.FactValue{Type: "number", Value: 30.0}},
        },
    }
    
    ctx := constraint.NewFactContext([]constraint.TypeDefinition{typeDef})
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := constraint.GenerateFactID(fact, typeDef, ctx)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// Ajouter BenchmarkFactGenerationWithReference
// Ajouter BenchmarkProgramParsing
// Ajouter BenchmarkCompleteFlow
```

### 5. Créer Script Global

**Fichier**: `scripts/run-e2e-tests.sh`

```bash
#!/bin/bash
set -e

echo "🧪 TESTS END-TO-END - NOUVELLE GESTION IDS"
echo "==========================================="

# Tests d'intégration
echo "📝 Tests d'intégration..."
go test -tags=integration ./tests/integration/... -v

# Tests E2E
echo "📝 Tests E2E..."
go test -tags=e2e ./tests/e2e/... -v

# Benchmarks
echo "📝 Benchmarks..."
go test ./tests/performance/... -bench=. -benchtime=1s

# Validation exemples
echo "📝 Validation exemples..."
for file in examples/*.tsd; do
    if [ -f "$file" ]; then
        echo "  - $(basename $file)"
        # Adapter selon CLI disponible
        # go run cmd/tsd/main.go validate "$file"
    fi
done

echo ""
echo "✅ TOUS LES TESTS E2E RÉUSSIS"
```

```bash
chmod +x scripts/run-e2e-tests.sh
```

---

## 📋 Actions Optionnelles (1-2h)

### 6. Nettoyer Code Déprécié

**Rechercher utilisations**:
```bash
grep -r "GenerateFactIDWithoutContext" --include="*.go" .
grep -r "valueToString" --include="*.go" . | grep -v "test"
grep -r "convertFactFieldValue" --include="*.go" . | grep -v "test"
```

**Si non utilisés**: Supprimer de `constraint/id_generator.go` et `constraint/constraint_facts.go`

### 7. Mettre à Jour Documentation

**Fichiers à mettre à jour**:
- `README.md` - Ajouter exemples nouvelles fonctionnalités
- `constraint/README.md` - Documenter nouvelles fonctions
- `tests/README.md` - Compléter avec nouveaux tests

---

## ✅ Checklist de Validation

Avant de considérer le travail terminé:

- [ ] Tests d'intégration créés et passent
- [ ] Tests E2E créés et passent
- [ ] Benchmarks exécutables
- [ ] Script global fonctionne
- [ ] Tous les tests existants passent toujours
- [ ] `make validate` passe
- [ ] Exemples validés manuellement
- [ ] Documentation mise à jour
- [ ] Code déprécié supprimé (si applicable)

---

## 🔗 Ressources

### Documents de Référence
- [Code Review](./REPORTS/code_review_id_system.md)
- [Refactoring Summary](./REPORTS/refactoring_id_system_summary.md)
- [Final Report](./REPORTS/final_report_review_refactoring.md)
- [Test Inventory](./REPORTS/new_ids_integration_tests_inventory.md)

### Fichiers Créés
- `tests/e2e/testdata/user_login.tsd` ✅
- `tests/e2e/testdata/order_management.tsd` ✅
- `tests/e2e/testdata/circular_reference_error.tsd` ✅
- `tests/e2e/testdata/undefined_variable_error.tsd` ✅
- `examples/new_syntax_demo.tsd` ✅
- `examples/advanced_relationships.tsd` ✅

### Tests Existants à Consulter
```bash
# Exemples de tests d'intégration
cat tests/integration/primary_key_e2e_test.go

# Exemples de tests E2E
cat tests/e2e/xuples_e2e_test.go

# Voir comment RETE est utilisé
grep -A 10 "NewNetwork" tests/integration/*.go
```

---

## 🚀 Commandes Rapides

```bash
# Développement tests
cd /home/resinsec/dev/tsd

# Créer tests d'intégration
vi tests/integration/fact_lifecycle_test.go

# Tester
go test -tags=integration ./tests/integration/... -v

# Créer tests E2E
vi tests/e2e/user_scenarios_test.go

# Tester E2E
go test -tags=e2e ./tests/e2e/... -v

# Créer benchmarks
vi tests/performance/id_generation_benchmark_test.go

# Exécuter benchmarks
go test ./tests/performance/... -bench=. -benchmem

# Validation globale
make test-complete
make validate
```

---

**Date**: 2025-12-19
**Priorité**: HAUTE
**Estimation**: 6-9h pour finalisation complète
**Prochaine action**: Vérifier API RETE et créer fact_lifecycle_test.go
