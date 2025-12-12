# Prompt 06 : JoinNode - Activation

**Session** : 6/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Prompt 05 complété, performJoinWithTokens optimisée

---

## 🎯 Objectif de cette Session

Réécrire les fonctions `ActivateLeft` et `ActivateRight` dans JoinNode pour :
1. Utiliser correctement BindingChain pour créer les tokens
2. Garantir que les tokens joints sont propagés avec TOUS les bindings
3. Assurer que les mémoires Left/Right fonctionnent correctement
4. Valider que les jointures à 2 variables continuent de fonctionner

**Livrable final** : `tsd/rete/node_join.go` (fonctions ActivateLeft/ActivateRight refactorées)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Analyser les Fonctions Actuelles (20 min)

#### 1.1 Lire ActivateLeft

**Fichier** : `tsd/rete/node_join.go`

**Questions à répondre** :
1. Comment le token est-il stocké dans LeftMemory ?
2. Comment les tokens de RightMemory sont-ils récupérés ?
3. Comment les tokens joints sont-ils créés et propagés ?
4. Y a-t-il des problèmes de gestion des bindings ?

**Documenter** : Noter les points à améliorer

---

#### 1.2 Lire ActivateRight

**Questions à répondre** :
1. Comment le fait est-il stocké dans RightMemory ?
2. Comment le token est-il créé pour ce fait ?
3. Quelle variable est utilisée pour lier le fait ?
4. Comment est-elle déterminée (getVariableForFact) ?

**Documenter** : Noter les points critiques

---

### Tâche 2 : Réécrire ActivateLeft (45 min)

#### 2.1 Nouvelle implémentation

**Fichier** : `tsd/rete/node_join.go`

**Code à implémenter** :

```go
// ActivateLeft est appelé quand un token arrive du côté gauche de la jointure.
// Le token contient les bindings des variables déjà jointes en amont.
func (jn *JoinNode) ActivateLeft(token *Token) error {
	if jn.Debug {
		fmt.Printf("\n🔍 [JOIN_%s] ActivateLeft CALLED\n", jn.ID)
		fmt.Printf("   Token ID: %s\n", token.ID)
		fmt.Printf("   Token Bindings: %v\n", token.GetVariables())
		fmt.Printf("   LeftVariables: %v\n", jn.LeftVariables)
	}
	
	// Stocker le token dans Left Memory
	jn.LeftMemory = append(jn.LeftMemory, token)
	
	if jn.Debug {
		fmt.Printf("   Left Memory size: %d\n", len(jn.LeftMemory))
		fmt.Printf("   Right Memory size: %d\n", len(jn.RightMemory))
	}
	
	// Tenter de joindre avec tous les faits de Right Memory
	for _, rightFact := range jn.RightMemory {
		// Créer un token pour le fait du côté droit
		rightToken, err := jn.createTokenForRightFact(rightFact)
		if err != nil {
			if jn.Debug {
				fmt.Printf("   ⚠️  Cannot create token for fact %s: %v\n", rightFact.ID, err)
			}
			continue
		}
		
		// Joindre les deux tokens
		joinedToken := jn.performJoinWithTokens(token, rightToken)
		
		if joinedToken != nil {
			if jn.Debug {
				fmt.Printf("   ✅ Join successful, propagating token with bindings: %v\n", 
					joinedToken.GetVariables())
			}
			
			// Propager le token joint aux enfants
			err := jn.PropagateToChildren(nil, joinedToken)
			if err != nil {
				return fmt.Errorf("error propagating joined token: %w", err)
			}
		} else {
			if jn.Debug {
				fmt.Printf("   ❌ Join failed (conditions not met)\n")
			}
		}
	}
	
	return nil
}
```

---

#### 2.2 Implémenter createTokenForRightFact

**Nouvelle fonction helper** :

```go
// createTokenForRightFact crée un token pour un fait du côté droit.
// Le token contiendra un seul binding : variable → fait.
func (jn *JoinNode) createTokenForRightFact(fact *Fact) (*Token, error) {
	// Déterminer quelle variable ce fait représente
	variable := jn.getVariableForFact(fact)
	if variable == "" {
		return nil, fmt.Errorf("no variable found for fact type %s (RightVariables: %v, VariableTypes: %v)",
			fact.Type, jn.RightVariables, jn.VariableTypes)
	}
	
	// Créer un token avec un seul binding
	token := NewTokenWithFact(fact, variable, jn.ID)
	
	if jn.Debug {
		fmt.Printf("   Created right token: variable=%s, fact=%s\n", variable, fact.ID)
	}
	
	return token, nil
}
```

---

#### 2.3 Vérifier getVariableForFact

**S'assurer que cette fonction est correcte** :

```go
// getVariableForFact retourne le nom de la variable pour un fait donné,
// en cherchant dans RightVariables et en matchant le type.
func (jn *JoinNode) getVariableForFact(fact *Fact) string {
	if fact == nil {
		return ""
	}
	
	// Chercher dans RightVariables
	for _, varName := range jn.RightVariables {
		// Vérifier le type attendu pour cette variable
		if expectedType, exists := jn.VariableTypes[varName]; exists {
			if expectedType == fact.Type {
				return varName
			}
		}
	}
	
	// Pas trouvé
	if jn.Debug {
		fmt.Printf("   ⚠️  getVariableForFact: no match for fact type %s\n", fact.Type)
		fmt.Printf("      RightVariables: %v\n", jn.RightVariables)
		fmt.Printf("      VariableTypes: %v\n", jn.VariableTypes)
	}
	
	return ""
}
```

**Point clé** : Cette fonction DOIT chercher dans `RightVariables` uniquement, pas dans toutes les variables.

---

### Tâche 3 : Réécrire ActivateRight (45 min)

#### 3.1 Nouvelle implémentation

**Code à implémenter** :

```go
// ActivateRight est appelé quand un fait arrive du côté droit de la jointure.
// Le fait représente une nouvelle variable à joindre avec les tokens existants.
func (jn *JoinNode) ActivateRight(fact *Fact) error {
	if jn.Debug {
		fmt.Printf("\n🔍 [JOIN_%s] ActivateRight CALLED\n", jn.ID)
		fmt.Printf("   Fact ID: %s\n", fact.ID)
		fmt.Printf("   Fact Type: %s\n", fact.Type)
		fmt.Printf("   RightVariables: %v\n", jn.RightVariables)
	}
	
	// Stocker le fait dans Right Memory
	jn.RightMemory = append(jn.RightMemory, fact)
	
	// Créer un token pour ce fait
	rightToken, err := jn.createTokenForRightFact(fact)
	if err != nil {
		// Logging de l'erreur mais continuer (le fait est stocké en mémoire)
		if jn.Debug {
			fmt.Printf("   ⚠️  Cannot create token for fact: %v\n", err)
		}
		return nil // Ne pas retourner d'erreur, juste ignorer
	}
	
	if jn.Debug {
		fmt.Printf("   Right token created with bindings: %v\n", rightToken.GetVariables())
		fmt.Printf("   Left Memory size: %d\n", len(jn.LeftMemory))
	}
	
	// Tenter de joindre avec tous les tokens de Left Memory
	for _, leftToken := range jn.LeftMemory {
		// Joindre les deux tokens
		joinedToken := jn.performJoinWithTokens(leftToken, rightToken)
		
		if joinedToken != nil {
			if jn.Debug {
				fmt.Printf("   ✅ Join successful with left token %s\n", leftToken.ID)
				fmt.Printf("      Joined token bindings: %v\n", joinedToken.GetVariables())
			}
			
			// Propager le token joint aux enfants
			err := jn.PropagateToChildren(nil, joinedToken)
			if err != nil {
				return fmt.Errorf("error propagating joined token: %w", err)
			}
		} else {
			if jn.Debug {
				fmt.Printf("   ❌ Join failed with left token %s (conditions not met)\n", 
					leftToken.ID)
			}
		}
	}
	
	return nil
}
```

---

### Tâche 4 : Tests Unitaires (50 min)

#### 4.1 Test ActivateLeft avec bindings multiples

**Fichier** : `tsd/rete/node_join_test.go`

**Test à ajouter** :

```go
func TestJoinNode_ActivateLeft_PreservesAllBindings(t *testing.T) {
	t.Log("🧪 TEST JoinNode - ActivateLeft préserve tous les bindings")
	
	// Setup : Faits
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Attributes: map[string]interface{}{"id": 1, "name": "Alice"},
	}
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Attributes: map[string]interface{}{"id": 100, "user_id": 1},
	}
	productFact := &Fact{
		ID:   "p1",
		Type: "Product",
		Attributes: map[string]interface{}{"id": 200},
	}
	
	// Setup : JoinNode configuré pour [user, order] → [user, order, product]
	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_test",
			Children: []Node{},
		},
		LeftVariables:  []string{"user", "order"},
		RightVariables: []string{"product"},
		AllVariables:   []string{"user", "order", "product"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     []*Token{},
		RightMemory:    []*Fact{},
		JoinConditions: nil, // Pas de conditions pour ce test
		Debug:          true,
	}
	
	// Ajouter le fait Product dans Right Memory
	joinNode.RightMemory = append(joinNode.RightMemory, productFact)
	
	// Créer un token du côté gauche avec [user, order]
	leftToken := &Token{
		ID:    "t_left",
		Facts: []*Fact{userFact, orderFact},
		Bindings: NewBindingChain().
			Add("user", userFact).
			Add("order", orderFact),
		NodeID: "upstream_node",
		Metadata: TokenMetadata{JoinLevel: 1},
	}
	
	// Mock du nœud enfant pour capturer le token propagé
	var capturedToken *Token
	mockChild := &MockNode{
		OnActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.Children = append(joinNode.Children, mockChild)
	
	// Act
	err := joinNode.ActivateLeft(leftToken)
	
	// Assert
	if err != nil {
		t.Fatalf("❌ ActivateLeft retourné erreur: %v", err)
	}
	
	if capturedToken == nil {
		t.Fatal("❌ Aucun token propagé au nœud enfant")
	}
	
	// Vérifier que le token propagé contient TOUS les bindings
	if capturedToken.Bindings.Len() != 3 {
		t.Errorf("❌ Attendu 3 bindings, got %d", capturedToken.Bindings.Len())
	}
	
	expectedVars := []string{"user", "order", "product"}
	for _, v := range expectedVars {
		if !capturedToken.HasBinding(v) {
			t.Errorf("❌ Variable '%s' manquante dans le token propagé", v)
		}
	}
	
	// Vérifier les valeurs
	if capturedToken.GetBinding("user") != userFact {
		t.Errorf("❌ Binding 'user' incorrect")
	}
	if capturedToken.GetBinding("order") != orderFact {
		t.Errorf("❌ Binding 'order' incorrect")
	}
	if capturedToken.GetBinding("product") != productFact {
		t.Errorf("❌ Binding 'product' incorrect")
	}
	
	t.Log("✅ ActivateLeft préserve bien tous les bindings")
}

// Mock pour capturer les tokens propagés
type MockNode struct {
	BaseNode
	OnActivateLeft func(*Token) error
}

func (m *MockNode) ActivateLeft(token *Token) error {
	if m.OnActivateLeft != nil {
		return m.OnActivateLeft(token)
	}
	return nil
}

func (m *MockNode) ActivateRight(fact *Fact) error {
	return nil
}
```

---

#### 4.2 Test ActivateRight avec bindings

**Test à ajouter** :

```go
func TestJoinNode_ActivateRight_CreatesCorrectToken(t *testing.T) {
	t.Log("🧪 TEST JoinNode - ActivateRight crée le bon token")
	
	// Setup
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Attributes: map[string]interface{}{"id": 1},
	}
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Attributes: map[string]interface{}{"id": 100, "user_id": 1},
	}
	
	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_test",
			Children: []Node{},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":  "User",
			"order": "Order",
		},
		LeftMemory:     []*Token{},
		RightMemory:    []*Fact{},
		JoinConditions: nil,
		Debug:          true,
	}
	
	// Ajouter un token dans Left Memory
	leftToken := NewTokenWithFact(userFact, "user", "upstream")
	joinNode.LeftMemory = append(joinNode.LeftMemory, leftToken)
	
	// Mock pour capturer
	var capturedToken *Token
	mockChild := &MockNode{
		OnActivateLeft: func(token *Token) error {
			capturedToken = token
			return nil
		},
	}
	joinNode.Children = append(joinNode.Children, mockChild)
	
	// Act
	err := joinNode.ActivateRight(orderFact)
	
	// Assert
	if err != nil {
		t.Fatalf("❌ ActivateRight retourné erreur: %v", err)
	}
	
	if capturedToken == nil {
		t.Fatal("❌ Aucun token propagé")
	}
	
	if capturedToken.Bindings.Len() != 2 {
		t.Errorf("❌ Attendu 2 bindings, got %d", capturedToken.Bindings.Len())
	}
	
	if !capturedToken.HasBinding("user") || !capturedToken.HasBinding("order") {
		t.Errorf("❌ Bindings manquants")
	}
	
	t.Log("✅ ActivateRight fonctionne correctement")
}
```

---

#### 4.3 Exécuter les tests

**Commandes** :

```bash
cd tsd

# Tests spécifiques
go test -v -run "TestJoinNode_Activate" ./rete/

# Tous les tests de JoinNode
go test -v ./rete/node_join_test.go

# Tous les tests rete
go test -v ./rete/...
```

**Résultat attendu** : Tous les tests passent ✅

---

### Tâche 5 : Tests d'Intégration (30 min)

#### 5.1 Tester avec une cascade réelle

**Créer un test d'intégration** :

```go
func TestJoinNode_Cascade_2Variables_Integration(t *testing.T) {
	t.Log("🧪 TEST JoinNode - Cascade 2 variables (intégration)")
	
	// Ce test simule une cascade TypeNode → JoinNode → TerminalNode
	// avec 2 variables : User et Order
	
	// Créer le réseau
	network := NewNetwork()
	
	// TypeNodes
	userTypeNode := &TypeNode{Type: "User"}
	orderTypeNode := &TypeNode{Type: "Order"}
	
	// JoinNode
	joinNode := &JoinNode{
		BaseNode: BaseNode{ID: "join_user_order"},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":  "User",
			"order": "Order",
		},
		LeftMemory:  []*Token{},
		RightMemory: []*Fact{},
		Debug:       true,
	}
	
	// Connecter
	userTypeNode.AddChild(joinNode)
	orderTypeNode.AddChild(joinNode)
	
	// TerminalNode mock
	var terminalToken *Token
	terminal := &MockNode{
		OnActivateLeft: func(token *Token) error {
			terminalToken = token
			return nil
		},
	}
	joinNode.AddChild(terminal)
	
	// Soumettre les faits
	userFact := &Fact{ID: "u1", Type: "User", Attributes: map[string]interface{}{"id": 1}}
	orderFact := &Fact{ID: "o1", Type: "Order", Attributes: map[string]interface{}{"user_id": 1}}
	
	network.AddFact(userFact)
	network.AddFact(orderFact)
	
	// Vérifier
	if terminalToken == nil {
		t.Fatal("❌ Terminal n'a pas reçu de token")
	}
	
	if terminalToken.Bindings.Len() != 2 {
		t.Errorf("❌ Attendu 2 bindings, got %d", terminalToken.Bindings.Len())
	}
	
	t.Log("✅ Cascade 2 variables fonctionne")
}
```

---

#### 5.2 Exécuter les tests d'intégration

**Commande** :

```bash
make test-integration
```

**Résultat attendu** : Tests d'intégration passent

---

### Tâche 6 : Validation et Nettoyage (20 min)

#### 6.1 Vérifier la non-régression

**Exécuter tous les tests** :

```bash
# Tests unitaires
go test ./rete/...

# Tests d'intégration
make test-integration

# Compilation
go build ./...
```

**Critère** : Aucune régression, tous les tests existants passent.

---

#### 6.2 Désactiver le logging de debug

**Option 1** : Supprimer les blocs `if jn.Debug { ... }`

**Option 2** : Garder mais s'assurer que `Debug = false` par défaut

**Recommandation** : Garder le logging mais désactivé, utile pour debugging futur.

---

#### 6.3 Vérifier la qualité du code

```bash
go fmt ./rete/node_join.go
go vet ./rete/node_join.go
```

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Code
- [ ] ✅ `ActivateLeft` utilise `createTokenForRightFact` et propage correctement
- [ ] ✅ `ActivateRight` crée un token avec le bon binding
- [ ] ✅ `getVariableForFact` cherche dans `RightVariables`
- [ ] ✅ Les tokens joints contiennent TOUS les bindings
- [ ] ✅ Gestion correcte des mémoires Left/Right

### Tests
- [ ] ✅ `TestJoinNode_ActivateLeft_PreservesAllBindings` passe
- [ ] ✅ `TestJoinNode_ActivateRight_CreatesCorrectToken` passe
- [ ] ✅ Tests d'intégration 2 variables passent
- [ ] ✅ Aucune régression sur les tests existants

### Qualité
- [ ] ✅ Code formatté et sans warnings
- [ ] ✅ Logging de debug désactivé ou supprimé
- [ ] ✅ GoDoc présent

---

## 🎯 Prochaine Étape

Une fois ActivateLeft/ActivateRight **refactorées et validées**, passer au **Prompt 07 - BetaChainBuilder**.

Le Prompt 07 s'assurera que les cascades sont construites avec les bonnes configurations (AllVariables, LeftVariables, RightVariables) à chaque niveau.

---

## 💡 Conseils Pratiques

### Pour les Tests
1. **Utiliser des mocks** : Facilite la capture des tokens propagés
2. **Tester les deux sens** : ActivateLeft et ActivateRight
3. **Vérifier les bindings** : Nombre + présence de chaque variable

### Pour le Debug
1. **Activer Debug = true dans les tests** : Voir exactement ce qui se passe
2. **Logger les IDs** : Facilite le suivi des tokens
3. **Vérifier les mémoires** : Left/Right doivent être remplies correctement

---

**Note** : Cette session garantit que les tokens sont correctement créés et propagés lors des activations. Le prochain prompt (07) s'assurera que la construction des cascades configure correctement les JoinNodes.