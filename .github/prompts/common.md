# 📋 Standards Communs - Projet TSD

## 🎯 Contexte du Projet

**TSD** : Solution générale de synchronisation utilisant un moteur de règles RETE avec système de contraintes en Go.

*Note : TSD n'est pas un acronyme, c'est simplement le nom du projet.*

---

## 🔧 GRAMMAIRE PEG ET PARSER

### Règles Critiques de Génération du Parser

**EMPLACEMENT DES FICHIERS** :
- Fichier source PEG : `constraint/grammar/constraint.peg`
- Parser généré : `constraint/parser.go` (à la racine du package constraint)
- ⚠️ **JAMAIS** dans `constraint/grammar/parser.go`

**Commande de génération** :
```bash
cd constraint/grammar
pigeon -o ../parser.go constraint.peg
```

**IMPORTANT** :
- Le parser DOIT être généré à la racine du package `constraint`
- Le fichier `.peg` reste dans `constraint/grammar/` pour l'organisation
- Tous les autres fichiers Go du package `constraint` sont à la racine
- Ne JAMAIS créer `constraint/grammar/parser.go`

**Vérification** :
```bash
# Bon emplacement
ls -l constraint/parser.go

# Mauvais emplacement (à supprimer si existe)
ls -l constraint/grammar/parser.go
```

---

## 🔒 LICENCE ET COPYRIGHT

### Vérification de Compatibilité

**AVANT toute utilisation de code externe, bibliothèque ou algorithme** :

| Statut | Licences | Action |
|--------|----------|--------|
| ✅ **Acceptées** | MIT, BSD, Apache-2.0, ISC | Utilisation autorisée |
| ⚠️ **À éviter** | GPL, AGPL, LGPL (copyleft) | Incompatible avec MIT |
| ❌ **Interdites** | Code sans licence, propriétaire | NE PAS UTILISER |

**Documentation obligatoire** :
- Code inspiré/adapté → Commentaire avec source
- Bibliothèque tierce → Mise à jour `go.mod` + `THIRD_PARTY_LICENSES.md`
- Algorithme connu → Citation académique

```go
// Algorithm based on: Dijkstra, E. W. (1959). "A note on two problems 
// in connexion with graphs". Numerische Mathematik, 1(1), 269-271.
// Implementation is original.
```

### En-tête de Copyright OBLIGATOIRE

**Tous les nouveaux fichiers `.go` doivent commencer par** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package [nom_du_package]
```

**Vérification avant commit** :
```bash
for file in $(find . -name "*.go" -type f ! -path "./.git/*"); do
    if ! head -1 "$file" | grep -q "Copyright\|Code generated"; then
        echo "⚠️  EN-TÊTE MANQUANT: $file"
    fi
done
```

---

## 🔑 IDENTIFIANTS DE FAITS (RETE)

### Principe Fondamental

**Il n'existe qu'UN SEUL identifiant pour un fait : l'ID interne.**

#### ID Interne (Identifiant Réel)

L'**ID interne** est l'identifiant unique et réel d'un fait dans le réseau RETE :

- **Format** : `Type~valeur` ou `Type~val1_val2_...` pour clés composites
- **Génération automatique** :
  - Si le type a des **clés primaires** définies → `Type~<valeurs_clés_primaires>`
  - Si le type n'a **pas de clés primaires** → `Type_<compteur>` (auto-incrémenté)
- **Accès** : Via `Fact.ID` (champ système)
- **Utilisation** : Indexation, recherche, rétraction dans le réseau RETE

**Exemple** :
```go
type Product(#id: string, name: string, price: number)
// Fait : Product(id: "prod_123", name: "Laptop", price: 999.99)
// ID interne : "Product~prod_123"

type Alert(level: string, message: string)  // Pas de clé primaire
// Fait : Alert(level: "HIGH", message: "Temperature warning")
// ID interne : "Alert_1" (auto-généré)
```

#### Attributs 'id' (Champs Ordinaires)

Les attributs nommés `id` (ou `ID`, `Id`, etc.) **n'ont RIEN de particulier** :

- Ce sont des **champs ordinaires** comme n'importe quel autre attribut
- Ils **ne sont PAS des identifiants** au sens système
- Ils peuvent être des **clés primaires** (avec `#id`) ou de simples attributs
- **Accès** : Via `Fact.Fields["id"]` (valeur d'attribut)

**Exemple** :
```go
type Person(#id: string, name: string, age: number)

person := &Fact{
    ID:     "Person~p1",           // ← ID INTERNE (identifiant réel)
    Type:   "Person",
    Fields: map[string]interface{}{
        "id":   "p1",              // ← Attribut 'id' (simple valeur de champ)
        "name": "Alice",
        "age":  30,
    },
}

// ✅ CORRECT : Accès à l'ID interne
internalID := person.ID  // "Person~p1"

// ✅ CORRECT : Accès à l'attribut 'id'
idField := person.Fields["id"]  // "p1"

// ❌ INCORRECT : Confondre les deux
// L'attribut 'id' n'est PAS l'identifiant du fait !
```

#### Règles Importantes

1. **Ne JAMAIS confondre** `Fact.ID` (identifiant interne) et `Fields["id"]` (attribut)
2. **Toujours utiliser** `Fact.ID` pour les opérations système (recherche, rétraction)
3. **Les clés primaires** (`#id`) servent à générer l'ID interne, mais ne sont pas l'ID interne
4. **Un type sans clé primaire** aura quand même un ID interne (auto-généré)

#### Dans les Tests

```go
// ✅ BON - Test de l'ID interne
assert.Equal(t, "Product~prod_123", fact.ID)
assert.Contains(t, fact.ID, "Product~")

// ✅ BON - Test de l'attribut 'id' (clé primaire)
assert.Equal(t, "prod_123", fact.Fields["id"])

// ❌ MAUVAIS - Confusion entre ID interne et attribut
// assert.Equal(t, fact.Fields["id"], fact.ID)  // FAUX !
```

---

## ⚠️ RÈGLES STRICTES - CODE GO

### 🚫 Interdictions Absolues

#### 1. AUCUN HARDCODING

❌ **Interdit** :
- Valeurs en dur dans le code
- "Magic numbers" ou "magic strings"
- Chemins de fichiers hardcodés
- Configurations hardcodées
- Code spécifique à un seul cas d'usage

✅ **Obligatoire** :
- Constantes nommées et explicites
- Variables de configuration
- Paramètres de fonction
- Interfaces pour abstraction
- Code générique et réutilisable

**Exemple** :

```go
// ❌ MAUVAIS - Hardcodé
func ProcessOrder(id string) error {
    if id == "special-customer-123" { // Hardcodé !
        discount = 0.25
    }
    timeout := 30 * time.Second // Magic number !
}

// ✅ BON - Générique
const DefaultTimeout = 30 * time.Second

type DiscountRule interface {
    ApplyDiscount(customerID string) float64
}

func ProcessOrder(id string, timeout time.Duration, rule DiscountRule) error {
    discount := rule.ApplyDiscount(id)
    // ... code générique
}
```

#### Code généré - OBLIGATOIRE
Ne modifie jamais le code généré par un outil tiers (typiquement `constraint/parser.go` généré par pigeon).

#### 2. TESTS Fonctionnels RÉELS

❌ **Interdit** :
- Simulation de résultats
- Mocks (sauf si explicitement demandé)
- Suppositions sur les résultats
- Tests non-déterministes (flaky)
- Dépendances entre tests
- Tests qui passent toujours
- Tests sans assertions

✅ **Obligatoire** :
- Extraction des résultats réels obtenus
- Pour RETE : interroger les TerminalNodes
- Pour RETE : inspecter les mémoires (Left/Right/Result)
- Tests isolés et indépendants
- Constantes nommées pour valeurs de test
- Assertions claires et explicites
- Messages d'erreur descriptifs

### ✅ Standards de Code Go

#### Conventions (Obligatoires)

| Aspect | Règle |
|--------|-------|
| **Style** | Effective Go, go fmt, goimports |
| **Nommage** | MixedCaps pour exports, idiomatique |
| **Erreurs** | Gestion explicite, pas de panic (sauf critique) |
| **Documentation** | GoDoc pour exports, commentaires inline si complexe |
| **Complexité** | Cyclomatique < 15 |
| **Fonctions** | < 50 lignes (sauf justification) |
| **Imbrication** | < 4 niveaux |
| **Duplication** | DRY - Don't Repeat Yourself |

#### Principes Architecturaux

- **Single Responsibility Principle** - Une fonction, une responsabilité
- **Open/Closed** - Extensible sans modification
- **Dependency Injection** - Pas de dépendances globales
- **Composition over Inheritance** - Interfaces et embedding
- **Interfaces** - Petites, focalisées, cohérentes
- **Découplage fort** - Couplage faible, cohésion forte

#### Qualité

- Code auto-documenté (noms explicites)
- Pas de "God Objects"
- Pas de code mort (dead code)
- Validation d'entrée systématique
- Gestion des cas nil/vides
- Pas de race conditions
- Pas de fuites mémoires

#### Visibilité et Encapsulation

- **Variables et fonctions privées par défaut** - Tout est privé (non exporté) sauf nécessité
- **Minimiser les exports publics** - N'exporter que ce qui fait partie du contrat d'interface
- **Respecter strictement les contrats** - L'API publique est un engagement
- **Préférer les interfaces** - Exposer des interfaces plutôt que des types concrets

---

## 🧪 STANDARDS DE TESTS

### Structure

```
project/
├── module/
│   ├── feature.go
│   ├── feature_test.go       # Tests unitaires
│   └── testdata/             # Données de test
│       └── test.tsd          # Fichiers TSD
└── tests/                     # Répertoire de tests racine (structure extensible)
    ├── e2e/                  # Tests end-to-end
    ├── fixtures/             # Fixtures partagées pour tests
    ├── integration/          # Tests d'intégration entre modules
    ├── performance/          # Tests de performance et benchmarks
    └── [autres types]/       # Structure extensible - ajoutez d'autres catégories selon les besoins
```

**Note importante** : La structure `tests/` est **extensible et non limitative**. 
Les sous-répertoires listés ci-dessus sont des exemples déjà présents dans le projet, 
mais vous pouvez ajouter d'autres catégories de tests selon les besoins 
(ex: `security/`, `stress/`, `acceptance/`, etc.).

### Checklist Tests

- [ ] **Couverture > 80%** (obligatoire)
- [ ] Cas nominaux testés
- [ ] Cas limites testés
- [ ] Cas d'erreur testés
- [ ] Table-driven tests si applicable
- [ ] Sous-tests (t.Run) si pertinent
- [ ] Noms explicites (*_test.go)
- [ ] Tests déterministes
- [ ] Tests isolés
- [ ] Messages clairs avec émojis (✅ ❌ ⚠️)
- [ ] Setup/teardown propre
- [ ] Pas de dépendances entre tests

### Template de Test

```go
func TestFeature(t *testing.T) {
    t.Log("🧪 TEST FEATURE")
    t.Log("================")
    
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"cas nominal", "input", "output", false},
        {"cas erreur", "", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            // Act
            result, err := Feature(tt.input)
            
            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("❌ Erreur = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("❌ Attendu '%s', reçu '%s'", tt.expected, result)
            }
            t.Log("✅ Test réussi")
        })
    }
}
```

---

## 📚 DOCUMENTATION

### Organisation Centralisée

```
tsd/
├── docs/                      # Documentation centralisée
│   ├── architecture/
│   ├── api/
│   └── [module]/             # Docs spécifiques aux modules
├── REPORTS/                   # Rapports (non versionnés)
├── scripts/                   # Scripts centralisés
├── README.md                  # Racine du projet
└── [module]/
    └── README.md             # README du module
```

### Standards

| Type | Langue | Format | Emplacement |
|------|--------|--------|-------------|
| **GoDoc** | Anglais | Commentaires Go | Dans le code |
| **Commentaires internes** | Français | Inline | Dans le code |
| **Documentation technique** | Français | Markdown | `docs/` |
| **README modules** | Français | Markdown | Racine module |
| **Exemples** | Code + Commentaires | `.go`, `.tsd` | `testdata/` |

### Checklist Documentation

- [ ] GoDoc pour toutes les fonctions exportées
- [ ] Commentaires inline pour code complexe
- [ ] Exemples d'utilisation testables
- [ ] README mis à jour si nécessaire
- [ ] CHANGELOG.md avec entrée si applicable
- [ ] TODO/FIXME documentés si nécessaire
- [ ] Pas de commentaires obsolètes

---

## 🔧 OUTILS ET COMMANDES

### Validation du Code

```bash
# Formattage (obligatoire avant commit)
go fmt ./...
goimports -w .

# Analyse statique (obligatoire)
go vet ./...
staticcheck ./...
errcheck ./...
gosec ./...

# Linting
golangci-lint run

# Tests
make test                    # Tests unitaires
make test-coverage          # Avec couverture
make test-integration       # Tests d'intégration

# Vérifications avancées
go test -race ./...         # Race conditions
gocyclo -over 15 .         # Complexité cyclomatique
```

### Makefile

Se référer au [Makefile](../../Makefile) pour toutes les commandes disponibles :

**Tests** :
- `make test` (alias de `test-unit`) - Tests unitaires uniquement
- `make test-unit` - Tests unitaires (rapides)
- `make test-fixtures` - Tests des fixtures partagées
- `make test-integration` - Tests d'intégration
- `make test-e2e` - Tests end-to-end
- `make test-performance` - Tests de performance
- `make test-all` - Tous les tests standards (unit + fixtures + integration + e2e + performance)
- `make test-complete` - **TOUS les tests** (complet, recommandé avant commit)
- `make test-coverage` - Rapport de couverture complet

**Validation** :
- `make validate` - Validation complète (format + lint + build + test-complete)
- `make quick-check` - Validation rapide sans tests
- `make ci` - Validation pour CI/CD

**Autres** :
- `make build` - Compilation
- `make clean` - Nettoyage
- `make lint` - Analyse statique
- `make format` - Formatage du code
- `make help` - Liste complète des commandes

---

## 🏗️ ARCHITECTURE ET ORGANISATION

### Principes

- **Évolution incrémentale** > Réécriture complète
- **Commencer simple** > Optimiser ensuite
- **Tests d'abord** (TDD encouragé)
- **NE PAS maintenir rétrocompatibilité** - Supprimer anciennes versions
- **Mise à jour documentation** - Supprimer docs obsolètes

### Performance

- Complexité algorithmique acceptable (O(n), O(n log n))
- Pas de boucles inutiles
- Pas de calculs redondants
- Slices/maps dimensionnés correctement
- Réutilisation d'objets si pertinent
- Benchmarks si optimisation nécessaire

### Concurrence

- Synchronisation correcte (mutex, channels)
- Pas de race conditions
- Pas de goroutine leaks
- Channels fermés proprement

### Sécurité

- Validation de toutes les entrées
- Pas d'injection possible
- Gestion des cas nil/vides
- Erreurs propagées correctement
- Messages informatifs sans informations sensibles
- Dépendances minimales et spécifiées

---

## 📦 GESTION DES DÉPENDANCES

### Règles

- Préférer les bibliothèques Go standard
- Versions spécifiées dans `go.mod`
- Pas de dépendances non nécessaires
- Licence vérifiée (voir section Licence)
- Documentation dans `THIRD_PARTY_LICENSES.md`

### Mise à jour

```bash
go mod tidy           # Nettoyer
go mod verify         # Vérifier
go mod download       # Télécharger
```

---

## 🎨 CONVENTIONS DE NOMMAGE

| Élément | Convention | Exemple |
|---------|------------|---------|
| **Packages** | lowercase, singulier | `rete`, `constraint` |
| **Fichiers** | lowercase_underscore | `node_join.go` |
| **Tests** | *_test.go | `node_join_test.go` |
| **Constantes** | MixedCaps/UPPER | `MaxNodes`, `DefaultTimeout` |
| **Variables** | camelCase | `nodeCount`, `resultToken` |
| **Fonctions** | MixedCaps | `ProcessToken`, `EvaluateCondition` |
| **Types** | MixedCaps | `AlphaNode`, `TokenMemory` |
| **Interfaces** | MixedCaps + "er" | `Evaluator`, `Processor` |

---

## 📋 CHECKLIST AVANT COMMIT

- [ ] **Copyright** : En-tête présent dans tous les nouveaux fichiers `.go`
- [ ] **Licence** : Code externe vérifié et documenté
- [ ] **Hardcoding** : Aucun hardcoding (valeurs, chemins, configs)
- [ ] **Généricité** : Code générique avec paramètres/interfaces
- [ ] **Constantes** : Toutes les valeurs ont des constantes nommées
- [ ] **Formattage** : `go fmt` + `goimports` appliqués
- [ ] **Linting** : `go vet` + `staticcheck` + `errcheck` sans erreur
- [ ] **Tests** : Tests écrits et passent (couverture > 80%)
- [ ] **Documentation** : GoDoc + README mis à jour
- [ ] **Validation** : `make validate` passe (inclut test-complete)
- [ ] **Non-régression** : Tous les tests passent (`make test-complete`)

---

## 🚀 WORKFLOW DE DÉVELOPPEMENT

1. **Analyse** - Comprendre le besoin et l'architecture existante
2. **Conception** - Planifier l'implémentation (interfaces, structures)
3. **Tests** - Écrire les tests d'abord (TDD)
4. **Implémentation** - Coder en respectant les standards
5. **Validation** - Tests + Linting + Formattage
6. **Documentation** - GoDoc + README + Exemples
7. **Revue** - Auto-revue avec checklist
8. **Commit** - Message clair et descriptif

---

## 📚 RESSOURCES

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Documentation TSD](../../docs/) - Documentation technique
- [Makefile](../../Makefile) - Commandes du projet

---

**Note** : Ce document définit les standards communs à tous les prompts. Chaque prompt spécifique peut ajouter des règles supplémentaires mais ne doit pas contredire ces standards.
