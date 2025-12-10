# 🧪 Tests - Prompt Universel

> **📋 Standards** : Ce prompt respecte les règles de [common.md](./common.md)

## 🎯 Objectif

Gérer les tests du projet TSD : écrire, exécuter, déboguer, ou analyser les tests.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter [common.md](./common.md) :
- [🧪 Standards Tests](./common.md#standards-de-tests) - Structure, couverture > 80%
- [⚠️ Tests Fonctionnels RÉELS](./common.md#tests-fonctionnels-réels) - Aucun mock, résultats réels
- [📋 Checklist](./common.md#checklist-tests) - Points de validation

---

## 📋 Instructions

### 1. Définir l'Action

**Précise** :
- **Type** : [ ] Écrire tests  [ ] Exécuter tests  [ ] Déboguer test  [ ] Analyser couverture
- **Cible** : Module/fonction/fichier concerné
- **Contexte** : Nouveauté, régression, optimisation ?

### 2. Écrire des Tests

#### Template de Base

```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   interface{}
        want    interface{}
        wantErr bool
    }{
        {"cas nominal", validInput, expectedOutput, false},
        {"cas erreur", invalidInput, nil, true},
        {"cas limite", edgeInput, edgeOutput, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Feature(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("❌ Erreur = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("❌ Attendu %v, reçu %v", tt.want, got)
            }
        })
    }
}
```

#### Principes Tests

- ✅ **Tests déterministes** : Mêmes entrées = mêmes sorties
- ✅ **Tests isolés** : Aucune dépendance entre tests
- ✅ **Résultats réels** : Pas de mocks (sauf explicitement nécessaire)
- ✅ **Couverture > 80%** : Cas nominaux + limites + erreurs
- ✅ **Messages clairs** : Émojis ✅ ❌ ⚠️ pour visibilité
- ✅ **Constantes nommées** : Pas de valeurs hardcodées

#### Structure Tests

```
module/
├── feature.go
├── feature_test.go          # Tests unitaires
└── testdata/
    └── cases.tsd            # Fichiers de test TSD

tests/
├── integration/             # Tests d'intégration
├── e2e/                     # Tests end-to-end
├── performance/             # Tests de performance
└── fixtures/                # Fixtures partagées
```

### 3. Exécuter des Tests

```bash
# Tests unitaires
go test ./...

# Avec couverture
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Tests spécifiques
go test -v -run TestFeature ./module

# Tests avec race detection
go test -race ./...

# Verbose avec détails
go test -v ./...

# Validation complète
make validate
```

### 4. Déboguer un Test

#### Si test échoue

1. **Lire le message d'erreur**
   - Identifier la ligne qui échoue
   - Comprendre l'assertion

2. **Vérifier les entrées**
   - Valeurs de test correctes ?
   - Cas edge couverts ?

3. **Examiner les résultats**
   - Résultat obtenu vs attendu
   - Erreur retournée vs attendue

4. **Isoler le problème**
   ```bash
   # Exécuter uniquement le test qui échoue
   go test -v -run TestProbleme ./module
   
   # Ajouter logs pour debugging
   t.Logf("🔍 Valeur intermédiaire: %v", value)
   ```

5. **Vérifier non-régression**
   - Le test est-il correct ?
   - Le code a-t-il régressé ?
   - L'environnement a-t-il changé ?

#### Si test flaky (non-déterministe)

- ❌ **Problème** : Concurrence, timing, aléatoire
- ✅ **Solution** : Rendre déterministe ou supprimer

### 5. Analyser la Couverture

```bash
# Générer rapport couverture
go test -coverprofile=coverage.out ./...

# Visualiser en HTML
go tool cover -html=coverage.out

# Par fonction
go tool cover -func=coverage.out

# Objectif : > 80% globalement
```

**Priorités couverture** :
1. Code critique (logique métier)
2. Gestion d'erreurs
3. Cas limites
4. Code public (API)

---

## 📝 Types de Tests

### Tests Unitaires

- **Localisation** : `*_test.go` à côté du code
- **Portée** : Fonction/méthode isolée
- **Vitesse** : Rapide (< 1s)
- **Objectif** : Comportement fonctionnel

```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("❌ Add(2,3) = %d, want 5", result)
    }
}
```

### Tests d'Intégration

- **Localisation** : `tests/integration/`
- **Portée** : Plusieurs modules ensemble
- **Vitesse** : Moyen (quelques secondes)
- **Objectif** : Intégration entre composants

### Tests E2E

- **Localisation** : `tests/e2e/`
- **Portée** : Système complet
- **Vitesse** : Lent (minutes)
- **Objectif** : Scénarios utilisateur

### Tests de Performance

- **Localisation** : `tests/performance/`
- **Type** : Benchmarks
- **Objectif** : Performance, non-régression

```go
func BenchmarkFeature(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Feature(input)
    }
}
```

---

## ✅ Checklist Tests

Voir [common.md#checklist-tests](./common.md#checklist-tests) :

- [ ] Couverture > 80%
- [ ] Cas nominaux testés
- [ ] Cas limites testés
- [ ] Cas d'erreur testés
- [ ] Tests déterministes
- [ ] Tests isolés
- [ ] Messages clairs avec émojis
- [ ] Pas de hardcoding dans tests
- [ ] Constantes nommées
- [ ] Tests passent localement

---

## 🎯 Bonnes Pratiques

1. **TDD** : Tests d'abord, code ensuite
2. **AAA** : Arrange, Act, Assert
3. **Table-driven** : Plusieurs cas dans un test
4. **Sous-tests** : `t.Run()` pour organisation
5. **Nommage** : `Test<Feature>_<Scenario>`
6. **Messages** : Descriptifs avec contexte
7. **Isolation** : Aucune dépendance entre tests
8. **Cleanup** : `t.Cleanup()` pour ressources

---

## 🚫 Anti-Patterns

- ❌ Tests qui passent toujours (faux positifs)
- ❌ Tests sans assertions
- ❌ Tests non-déterministes (flaky)
- ❌ Dépendances entre tests
- ❌ Mocks abusifs (tester les vrais résultats)
- ❌ Tests trop complexes
- ❌ Hardcoding de valeurs
- ❌ Coverage pour coverage (privilégier qualité)

---

## 📊 Commandes Utiles

```bash
# Tests
go test ./...                              # Tous
go test -v ./...                           # Verbose
go test -run TestName ./...                # Spécifique
go test -short ./...                       # Tests courts uniquement

# Couverture
go test -cover ./...                       # Pourcentage
go test -coverprofile=coverage.out ./...   # Rapport
go tool cover -html=coverage.out           # Visualisation

# Performance
go test -bench=. ./...                     # Benchmarks
go test -benchmem ./...                    # Avec mémoire

# Validation
go test -race ./...                        # Race conditions
make test                                  # Tests standard
make test-coverage                         # Avec couverture
make validate                              # Validation complète
```

---

## 📚 Ressources

- [common.md](./common.md) - Standards tests
- [Testing Package](https://pkg.go.dev/testing) - Documentation Go
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests) - Pattern
- [Makefile](../../Makefile) - Commandes projet

---

**Workflow** : Écrire → Exécuter → Vérifier → Déboguer → Valider → Commit