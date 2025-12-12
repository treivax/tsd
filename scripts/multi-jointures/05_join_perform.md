# Prompt 05 : JoinNode - performJoinWithTokens

**Session** : 5/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Prompt 04 complété, Token refactoré avec BindingChain

---

## 🎯 Objectif de cette Session

Optimiser la fonction `performJoinWithTokens` dans JoinNode pour garantir que :
1. La composition des bindings via BindingChain est correcte
2. TOUS les bindings des deux tokens sont préservés dans le token joint
3. La logique de jointure est claire et traçable

**Livrable final** : `tsd/rete/node_join.go` (fonction `performJoinWithTokens` optimisée)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Analyser la Fonction Actuelle (20 min)

#### 1.1 Lire la fonction performJoinWithTokens

**Fichier** : `tsd/rete/node_join.go`

**Questions à répondre** :
1. Est-ce que `Merge()` est utilisée pour combiner les bindings ?
2. Est-ce que le token joint contient bien les Facts des deux tokens ?
3. Est-ce que les métadonnées sont correctement remplies ?
4. Y a-t-il des cas où des bindings pourraient être perdus ?

**Documenter** : Noter les problèmes identifiés

---

#### 1.2 Vérifier evaluateJoinConditions

**Fonction** : `func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool`

**Vérifier** :
- La signature utilise bien `*BindingChain` (pas `map[string]*Fact`)
- Tous les accès aux bindings utilisent `.Get(variable)`
- Les conditions de jointure sont correctement évaluées

**Si la signature est encore `map[string]*Fact`** : La mettre à jour maintenant.

---

### Tâche 2 : Optimiser performJoinWithTokens (60 min)

#### 2.1 Réécrire la fonction complète

**Fichier** : `tsd/rete/node_join.go`

**Code à implémenter** :

```go
// performJoinWithTokens combine deux tokens en vérifiant les conditions de jointure.
// Retourne un nouveau token avec TOUS les bindings des deux tokens parents,
// ou nil si les conditions de jointure ne sont pas satisfaites.
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
	// Étape 1 : Composer les chaînes de bindings (immuable)
	var newBindings *BindingChain
	
	if token1.Bindings == nil {
		newBindings = token2.Bindings
	} else if token2.Bindings == nil {
		newBindings = token1.Bindings
	} else {
		// Merge : tous les bindings de token1 + tous les bindings de token2
		newBindings = token1.Bindings.Merge(token2.Bindings)
	}
	
	// Logging pour debug (TEMPORAIRE - à supprimer après validation)
	if jn.Debug {
		fmt.Printf("🔗 [JOIN_%s] performJoinWithTokens\n", jn.ID)
		fmt.Printf("   Token1: ID=%s, Bindings=%v\n", token1.ID, token1.GetVariables())
		fmt.Printf("   Token2: ID=%s, Bindings=%v\n", token2.ID, token2.GetVariables())
		fmt.Printf("   Merged: Bindings=%v\n", newBindings.Variables())
	}
	
	// Étape 2 : Vérifier les conditions de jointure
	if !jn.evaluateJoinConditions(newBindings) {
		if jn.Debug {
			fmt.Printf("   ❌ Join conditions FAILED\n")
		}
		return nil
	}
	
	if jn.Debug {
		fmt.Printf("   ✅ Join conditions PASSED\n")
	}
	
	// Étape 3 : Combiner les facts
	combinedFacts := make([]*Fact, 0, len(token1.Facts)+len(token2.Facts))
	combinedFacts = append(combinedFacts, token1.Facts...)
	combinedFacts = append(combinedFacts, token2.Facts...)
	
	// Étape 4 : Créer le token joint
	joinedToken := &Token{
		ID:       generateTokenID(),
		Facts:    combinedFacts,
		Bindings: newBindings, // ✅ Chaîne complète avec TOUS les bindings
		NodeID:   jn.ID,
		Metadata: TokenMetadata{
			CreatedAt:    time.Now(),
			CreatedBy:    jn.ID,
			JoinLevel:    maxInt(token1.Metadata.JoinLevel, token2.Metadata.JoinLevel) + 1,
			ParentTokens: []string{token1.ID, token2.ID},
		},
	}
	
	if jn.Debug {
		fmt.Printf("   Created token: ID=%s, Bindings=%v, Facts=%d\n",
			joinedToken.ID, joinedToken.GetVariables(), len(joinedToken.Facts))
	}
	
	return joinedToken
}

// Helper function
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
```

**Points clés** :
- Utilisation de `Merge()` pour garantir que tous les bindings sont combinés
- Gestion des cas nil
- Métadonnées complètes (JoinLevel, ParentTokens)
- Logging temporaire pour debugging

---

#### 2.2 Ajouter le flag Debug dans JoinNode

**Si pas déjà présent**, ajouter dans la structure JoinNode :

```go
type JoinNode struct {
	BaseNode
	// ... champs existants ...
	Debug bool // Flag pour logging temporaire
}
```

**Activer le debug dans les tests** :

```go
joinNode.Debug = true  // Pour voir les traces
```

---

### Tâche 3 : Adapter evaluateJoinConditions (40 min)

#### 3.1 Mettre à jour la signature

**Si ce n'est pas déjà fait** :

```go
// Ancienne signature
func (jn *JoinNode) evaluateJoinConditions(bindings map[string]*Fact) bool

// Nouvelle signature
func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool
```

---

#### 3.2 Remplacer tous les accès aux bindings

**Pattern à chercher et remplacer** :

```go
// ANCIEN
fact := bindings[variable]
if fact == nil {
    return false
}

// NOUVEAU
if !bindings.Has(variable) {
    return false
}
fact := bindings.Get(variable)
```

**Pattern pour itération** :

```go
// ANCIEN
for variable, fact := range bindings {
    // ...
}

// NOUVEAU
for _, variable := range bindings.Variables() {
    fact := bindings.Get(variable)
    // ...
}
```

---

#### 3.3 Exemple complet de evaluateJoinConditions

```go
func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool {
	// Pas de conditions = toujours vrai
	if jn.JoinConditions == nil || len(jn.JoinConditions) == 0 {
		return true
	}
	
	// Vérifier que toutes les variables requises sont présentes
	for _, variable := range jn.AllVariables {
		if !bindings.Has(variable) {
			if jn.Debug {
				fmt.Printf("   ⚠️  Variable '%s' manquante dans bindings\n", variable)
			}
			return false
		}
	}
	
	// Évaluer chaque condition de jointure
	for _, condition := range jn.JoinConditions {
		// Extraire les variables de la condition
		leftVar, rightVar := jn.extractVariablesFromCondition(condition)
		
		// Récupérer les facts via BindingChain
		leftFact := bindings.Get(leftVar)
		rightFact := bindings.Get(rightVar)
		
		if leftFact == nil || rightFact == nil {
			return false
		}
		
		// Évaluer la condition
		result, err := jn.ConditionEvaluator.EvaluateJoinCondition(
			condition,
			leftFact,
			rightFact,
		)
		
		if err != nil || !result {
			if jn.Debug {
				fmt.Printf("   ❌ Condition failed: %v (err: %v)\n", condition, err)
			}
			return false
		}
	}
	
	return true
}
```

**Note** : Adapter cette fonction selon votre implémentation actuelle.

---

### Tâche 4 : Tests et Validation (40 min)

#### 4.1 Créer un test unitaire pour performJoinWithTokens

**Fichier** : `tsd/rete/node_join_test.go`

**Test à ajouter** :

```go
func TestJoinNode_PerformJoinWithTokens_PreservesAllBindings(t *testing.T) {
	t.Log("🧪 TEST JoinNode - performJoinWithTokens préserve tous les bindings")
	
	// Setup
	userFact := &Fact{ID: "u1", Type: "User", Attributes: map[string]interface{}{"id": 1}}
	orderFact := &Fact{ID: "o1", Type: "Order", Attributes: map[string]interface{}{"user_id": 1}}
	productFact := &Fact{ID: "p1", Type: "Product", Attributes: map[string]interface{}{"id": 100}}
	
	// Token1 : [user, order]
	token1 := &Token{
		ID:    "t1",
		Facts: []*Fact{userFact, orderFact},
		Bindings: NewBindingChain().
			Add("user", userFact).
			Add("order", orderFact),
		NodeID: "test_node_1",
		Metadata: TokenMetadata{JoinLevel: 1},
	}
	
	// Token2 : [product]
	token2 := &Token{
		ID:    "t2",
		Facts: []*Fact{productFact},
		Bindings: NewBindingChain().
			Add("product", productFact),
		NodeID: "test_node_2",
		Metadata: TokenMetadata{JoinLevel: 0},
	}
	
	// JoinNode sans conditions (toujours vrai)
	joinNode := &JoinNode{
		BaseNode: BaseNode{ID: "join_test"},
		AllVariables: []string{"user", "order", "product"},
		JoinConditions: nil, // Pas de conditions pour ce test
		Debug: true,
	}
	
	// Act
	result := joinNode.performJoinWithTokens(token1, token2)
	
	// Assert
	if result == nil {
		t.Fatal("❌ performJoinWithTokens devrait retourner un token, got nil")
	}
	
	// Vérifier le nombre de bindings
	if result.Bindings.Len() != 3 {
		t.Errorf("❌ Attendu 3 bindings, got %d", result.Bindings.Len())
	}
	
	// Vérifier que chaque variable est présente
	expectedVars := []string{"user", "order", "product"}
	for _, v := range expectedVars {
		if !result.HasBinding(v) {
			t.Errorf("❌ Variable '%s' manquante dans le token joint", v)
		}
	}
	
	// Vérifier les valeurs
	if result.GetBinding("user") != userFact {
		t.Errorf("❌ Binding 'user' incorrect")
	}
	if result.GetBinding("order") != orderFact {
		t.Errorf("❌ Binding 'order' incorrect")
	}
	if result.GetBinding("product") != productFact {
		t.Errorf("❌ Binding 'product' incorrect")
	}
	
	// Vérifier les facts
	if len(result.Facts) != 3 {
		t.Errorf("❌ Attendu 3 facts, got %d", len(result.Facts))
	}
	
	// Vérifier les métadonnées
	if result.Metadata.JoinLevel != 2 {
		t.Errorf("❌ JoinLevel attendu 2, got %d", result.Metadata.JoinLevel)
	}
	
	if len(result.Metadata.ParentTokens) != 2 {
		t.Errorf("❌ Attendu 2 parents, got %d", len(result.Metadata.ParentTokens))
	}
	
	t.Log("✅ performJoinWithTokens préserve bien tous les bindings")
}
```

---

#### 4.2 Exécuter le test

**Commande** :

```bash
cd tsd
go test -v -run "TestJoinNode_PerformJoinWithTokens" ./rete/
```

**Résultat attendu** : Test passe ✅

**Si échec** :
- Analyser le message d'erreur
- Vérifier les logs de debug
- Corriger la fonction
- Re-tester

---

#### 4.3 Tester avec les tests existants

**Commande** :

```bash
# Tous les tests de JoinNode
go test -v ./rete/node_join_test.go

# Tous les tests du module rete
go test -v ./rete/...
```

**Attendu** : Les tests existants doivent continuer à passer (non-régression)

---

### Tâche 5 : Validation et Nettoyage (20 min)

#### 5.1 Vérifier les jointures 2 variables

**S'assurer** que les jointures à 2 variables fonctionnent toujours :

```bash
# Tests d'intégration
make test-integration
```

**Si régression** : Identifier et fixer.

---

#### 5.2 Supprimer le logging de debug

**Une fois validé**, désactiver ou supprimer le logging :

Option 1 : Supprimer complètement
```go
// Supprimer tous les blocs if jn.Debug { ... }
```

Option 2 : Garder mais désactiver par défaut
```go
// S'assurer que Debug = false par défaut
```

---

#### 5.3 Vérifier la qualité du code

**Commandes** :

```bash
# Formattage
go fmt ./rete/node_join.go

# Analyse statique
go vet ./rete/node_join.go

# Vérifier la complexité
gocyclo -over 15 ./rete/node_join.go
```

**Critères** :
- Code formaté correctement
- Pas de warnings
- Complexité cyclomatique acceptable
- Commentaires GoDoc présents

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Code
- [ ] ✅ `performJoinWithTokens` utilise `Merge()` pour combiner les bindings
- [ ] ✅ TOUS les bindings sont préservés (prouvé par tests)
- [ ] ✅ `evaluateJoinConditions` utilise `*BindingChain`
- [ ] ✅ Métadonnées correctement remplies (JoinLevel, ParentTokens)
- [ ] ✅ Gestion des cas nil

### Tests
- [ ] ✅ Test unitaire `TestJoinNode_PerformJoinWithTokens_PreservesAllBindings` passe
- [ ] ✅ Tests existants de JoinNode passent (non-régression)
- [ ] ✅ Tests d'intégration passent

### Qualité
- [ ] ✅ Code formatté (`go fmt`)
- [ ] ✅ Pas de warnings (`go vet`)
- [ ] ✅ Complexité acceptable
- [ ] ✅ GoDoc présent
- [ ] ✅ Logging de debug supprimé ou désactivé

---

## 🎯 Résultats Attendus

### Comportement
- Un token joint contient TOUS les bindings de ses deux parents
- Les conditions de jointure sont évaluées sur la chaîne complète
- Les métadonnées permettent de tracer la provenance du token

### Performance
- Pas de régression (Merge est O(m) où m = taille de token2)
- Overhead acceptable pour les cas typiques (n < 10)

---

## 🎯 Prochaine Étape

Une fois `performJoinWithTokens` **optimisée et validée**, passer au **Prompt 06 - JoinNode Activation**.

Le Prompt 06 réécritera `ActivateLeft` et `ActivateRight` pour garantir la propagation correcte des tokens joints.

---

## 💡 Conseils Pratiques

### Pour le Debug
1. **Activer le logging temporairement** : Très utile pour voir le flux
2. **Vérifier les bindings à chaque étape** : Avant merge, après merge, dans le token joint
3. **Logger les IDs des tokens** : Facilite le traçage

### Pour la Validation
1. **Tester avec 2 variables d'abord** : Vérifier la non-régression
2. **Tester avec 3 variables** : Vérifier que tous les bindings sont présents
3. **Tester les cas limites** : Token avec bindings nil, conditions complexes

### Pour la Qualité
1. **Commentaires clairs** : Expliquer le "pourquoi" de chaque étape
2. **Gestion d'erreurs** : Retourner nil si les conditions échouent
3. **Invariants** : S'assurer que le token joint est valide

---

**Note** : Cette session se concentre sur la **fonction de jointure elle-même**. L'activation (quand cette fonction est appelée) sera traitée dans le Prompt 06. Restez focalisé sur la logique de composition des bindings.