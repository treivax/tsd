# ✨ Ajouter une Nouvelle Fonctionnalité

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux ajouter une nouvelle fonctionnalité au projet TSD, qu'il s'agisse d'un nouvel opérateur, d'un nouveau type de nœud RETE, d'une amélioration du parseur, ou de toute autre fonctionnalité.

## Objectif

Implémenter proprement une nouvelle fonctionnalité en respectant l'architecture existante et les conventions du projet.

## ⚠️ RÈGLES STRICTES - CODE GOLANG

### 🚫 INTERDICTIONS ABSOLUES

1. **AUCUN HARDCODING** :
   - ❌ Pas de valeurs en dur dans le code
   - ❌ Pas de "magic numbers" ou "magic strings"
   - ❌ Pas de chemins de fichiers hardcodés
   - ❌ Pas de configurations hardcodées
   - ✅ Utiliser des constantes nommées
   - ✅ Utiliser des variables de configuration
   - ✅ Utiliser des paramètres de fonction

2. **CODE TOUJOURS GÉNÉRIQUE** :
   - ✅ Fonctions réutilisables avec paramètres
   - ✅ Types génériques quand approprié
   - ✅ Interfaces pour abstraction
   - ✅ Code extensible sans modification
   - ❌ Pas de code spécifique à un cas d'usage

### ✅ BONNES PRATIQUES GO OBLIGATOIRES

1. **Conventions Go** :
   - Respect de Effective Go
   - Nommage idiomatique (MixedCaps pour export)
   - Gestion explicite des erreurs (pas de panic sauf critique)
   - go fmt et goimports appliqués
   - Commentaires GoDoc pour exports

2. **Architecture** :
   - Single Responsibility Principle
   - Interfaces petites et focalisées
   - Composition over inheritance
   - Dependency injection
   - Découplage fort

3. **Qualité** :
   - Code auto-documenté
   - Complexité cyclomatique < 15
   - Fonctions < 50 lignes (sauf justification)
   - Pas de duplication (DRY)
   - Tests unitaires obligatoires

**Exemples** :

❌ **MAUVAIS - Hardcodé** :
```go
func ProcessOrder(id string) error {
    if id == "special-customer-123" {  // Hardcodé !
        discount = 0.25
    }
    timeout := 30 * time.Second  // Magic number !
}
```

✅ **BON - Générique** :
```go
type DiscountRule interface {
    ApplyDiscount(customerID string) float64
}

func ProcessOrder(id string, timeout time.Duration, rule DiscountRule) error {
    discount := rule.ApplyDiscount(id)
    // ... code générique
}
```

## Instructions

### 1. Définir la Fonctionnalité

**Spécifie clairement** :
- **Nom de la fonctionnalité** : Ex. "Support des opérateurs de comparaison de chaînes"
- **Description** : Ce que la fonctionnalité doit faire
- **Cas d'usage** : Exemples concrets d'utilisation
- **Portée** : Modules affectés (rete, constraint, test, etc.)

### 2. Analyser l'Architecture Existante

1. **Examiner les composants similaires** :
   - Y a-t-il déjà quelque chose de similaire ?
   - Comment est-ce implémenté actuellement ?
   - Quelles sont les conventions de code ?

2. **Identifier les points d'intégration** :
   - Parseur (grammaire PEG)
   - Nœuds RETE (Alpha, Beta, Join, etc.)
   - Évaluateurs de conditions
   - Tests

3. **Vérifier les dépendances** :
   - Quels modules doivent être modifiés ?
   - Y a-t-il des impacts sur l'API existante ?

### 3. Concevoir l'Implémentation

1. **Architecture** :
   - Quels fichiers créer/modifier ?
   - Quelle structure de données utiliser ?
   - Comment s'intégrer avec l'existant ?

2. **API** :
   - Quelles fonctions/méthodes exposer ?
   - Quelle signature de fonctions ?
   - Quelles interfaces implémenter ?

3. **Tests** :
   - Quels tests unitaires ajouter ?
   - Quels tests d'intégration créer ?
   - Quels fichiers `.constraint` et `.facts` créer ?

### 4. Implémenter la Fonctionnalité

**Suivre l'ordre** :

1. **Commencer par les tests** (TDD) :
   ```go
   func TestNouvelleFeature(t *testing.T) {
       // Test de la nouvelle fonctionnalité
   }
   ```

2. **Implémenter le code minimal** :
   - ⚠️ **VÉRIFIER** : Aucun hardcoding introduit
   - ⚠️ **VÉRIFIER** : Code générique et réutilisable
   - Créer les structures nécessaires avec constantes nommées
   - Implémenter les fonctions de base avec paramètres
   - Faire passer les tests
   - Valider avec go vet et golangci-lint

3. **Ajouter la documentation** :
   - Commentaires GoDoc
   - Exemples dans les tests
   - Mise à jour du README si nécessaire

4. **Intégrer avec l'existant** :
   - Connecter aux autres modules
   - Gérer les cas limites
   - Ajouter la validation d'erreurs

### 5. Tester et Valider

1. **Tests unitaires** :
   ```bash
   go test -v -run TestNouvelleFeature ./rete
   ```

2. **Tests d'intégration** :
   ```bash
   make test-integration
   ```

3. **Validation complète** :
   ```bash
   make validate  # format + lint + build + test
   make rete-unified  # Runner universel
   ```

4. **Tests de performance** (si applicable) :
   ```bash
   go test -bench=. -benchmem ./rete
   ```

### 6. Documenter

1. **Code** :
   - Commentaires clairs sur les nouvelles fonctions
   - Exemples d'utilisation dans GoDoc

2. **Tests** :
   - Fichiers `.constraint` d'exemple
   - Fichiers `.facts` de test

3. **Documentation projet** :
   - Mettre à jour `README.md` si nécessaire
   - Ajouter une entrée dans `CHANGELOG.md`
   - Créer une doc technique dans `docs/` si importante

## Critères de Succès

✅ La fonctionnalité est implémentée et fonctionne
✅ Tests unitaires passent (100% de couverture si possible)
✅ Tests d'intégration passent
✅ Runner universel passe (58/58)
✅ Aucune régression introduite
✅ Code documenté (GoDoc)
✅ Suit les conventions du projet
✅ Performance acceptable

## Structure de Fichiers Typique

```
tsd/
├── rete/
│   ├── nouvelle_feature.go         # Implémentation
│   ├── nouvelle_feature_test.go    # Tests unitaires
│   └── testdata/                   # Données de test
│       ├── feature_test.constraint
│       └── feature_test.facts
├── constraint/
│   └── grammar.peg                 # Si modification du parseur
├── test/integration/
│   └── nouvelle_feature_test.go    # Tests d'intégration
└── docs/
    └── nouvelle_feature.md         # Documentation détaillée
```

## Exemple d'Utilisation

```
Je veux ajouter le support des opérateurs de comparaison de chaînes 
(startsWith, endsWith, contains) dans les AlphaNodes. 

Exemple d'utilisation :
{p: Person} / p.name startsWith "Alice" ==> action(p)

Peux-tu utiliser le prompt "add-feature" pour m'aider à implémenter ça ?
```

## Template de Code

### Nouvelle Fonction/Méthode

```go
// NouvelleFeature fait quelque chose d'utile.
// 
// Paramètres:
//   - param1: description du paramètre
//   - param2: description du paramètre
//
// Retourne:
//   - result: description du résultat
//   - error: erreur si problème
//
// Exemple:
//   result, err := NouvelleFeature("valeur")
//   if err != nil {
//       log.Fatal(err)
//   }
func NouvelleFeature(param1 string, param2 int) (result string, err error) {
    // Implémentation
    return
}
```

### Nouveau Test

```go
func TestNouvelleFeature(t *testing.T) {
    t.Log("🧪 TEST NOUVELLE FONCTIONNALITÉ")
    t.Log("================================")
    
    // Arrange
    input := "test"
    expected := "result"
    
    // Act
    result, err := NouvelleFeature(input, 42)
    
    // Assert
    if err != nil {
        t.Fatalf("❌ Erreur inattendue: %v", err)
    }
    
    if result != expected {
        t.Errorf("❌ Attendu '%s', reçu '%s'", expected, result)
    }
    
    t.Log("✅ Test réussi")
}
```

### Nouveau Type

```go
// NouveauType représente [description].
type NouveauType struct {
    // Champs avec documentation
    Field1 string `json:"field1"` // Description du champ
    Field2 int    `json:"field2"` // Description du champ
}

// NewNouveauType crée une nouvelle instance de NouveauType.
func NewNouveauType(field1 string) *NouveauType {
    return &NouveauType{
        Field1: field1,
        Field2: 0,
    }
}

// Methode fait quelque chose avec NouveauType.
func (n *NouveauType) Methode() error {
    // Implémentation
    return nil
}
```

## Checklist Avant de Commencer

- [ ] J'ai bien compris ce que je veux implémenter
- [ ] J'ai vérifié qu'il n'existe pas déjà
- [ ] J'ai analysé l'architecture existante
- [ ] J'ai conçu l'implémentation
- [ ] J'ai préparé les tests

## Checklist Après Implémentation

- [ ] **AUCUN HARDCODING** vérifié et validé
- [ ] **CODE GÉNÉRIQUE** vérifié (paramètres, interfaces)
- [ ] **CONSTANTES NOMMÉES** pour toutes les valeurs
- [ ] Tests unitaires écrits et passent
- [ ] Tests d'intégration écrits et passent
- [ ] Aucune régression (make test && make rete-unified)
- [ ] Code formaté (go fmt, goimports)
- [ ] Code linté (go vet, golangci-lint)
- [ ] Documentation GoDoc ajoutée
- [ ] Exemples d'utilisation fournis
- [ ] CHANGELOG.md mis à jour
- [ ] README.md mis à jour si nécessaire

## Bonnes Pratiques

### Code Go
- **OBLIGATOIRE** : Aucun hardcoding (valeurs, chemins, configs)
- **OBLIGATOIRE** : Code générique avec paramètres/interfaces
- **OBLIGATOIRE** : Constantes nommées pour toutes les valeurs
- Suivre les conventions Go (Effective Go)
- Utiliser les types et interfaces appropriés
- Gérer les erreurs explicitement (pas de panic)
- Utiliser des noms descriptifs et idiomatiques
- go vet et golangci-lint sans erreur

### Tests
- Un test = un cas d'usage
- Tests déterministes (pas d'aléatoire)
- Tests isolés (pas de dépendances entre tests)
- Messages d'erreur clairs avec émojis (✅ ❌ ⚠️)
- Utiliser des sous-tests (t.Run) si nécessaire

### Documentation
- Commentaires en français pour cohérence projet
- GoDoc en anglais pour compatibilité Go
- Exemples concrets et testables
- Diagrammes si architecture complexe

## Ressources

- [Makefile](../../Makefile) - Commandes disponibles
- [Architecture RETE](../../docs/) - Documentation technique
- [Grammaire PEG](../../constraint/grammar.peg) - Syntaxe des contraintes
- [Tests existants](../../rete/) - Exemples de tests

## Notes

- Préférer l'évolution incrémentale à la réécriture complète
- Commencer simple, optimiser ensuite si nécessaire
- Demander une revue de code si changement important
- Penser à la rétrocompatibilité