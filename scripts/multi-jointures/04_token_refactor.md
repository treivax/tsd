# Prompt 04 : Refactoring de Token vers BindingChain

**Session** : 4/12  
**Durée estimée** : 2-4 heures  
**Pré-requis** : Prompt 03 complété, BindingChain implémentée et testée

---

## 🎯 Objectif de cette Session

Remplacer **complètement** l'ancienne structure Token pour utiliser BindingChain, en :
1. Modifiant `rete/fact_token.go` : remplacer `Bindings map[string]*Fact` par `Bindings *BindingChain`
2. Fixant **toutes** les erreurs de compilation dans le code source
3. Adaptant les tests existants
4. Validant que tout compile et que les tests passent

**Principe clé** : **Remplacement direct, pas de cohabitation ancien/nouveau code**

---

## 📋 Tâches à Réaliser

### Tâche 1 : Modifier la Structure Token (30 min)

#### 1.1 Sauvegarder et analyser l'ancien Token

**Lire** : `tsd/rete/fact_token.go`

**Ancienne structure** :
```go
type Token struct {
    ID       string
    Facts    []*Fact
    Bindings map[string]*Fact  // ❌ À remplacer
    NodeID   string
}
```

**Identifier tous les usages** :
```bash
cd tsd
grep -r "\.Bindings\[" rete/
grep -r "\.Bindings =" rete/
grep -r "range.*\.Bindings" rete/
grep -r "map\[string\]\*Fact" rete/
```

Noter tous les fichiers affectés.

---

#### 1.2 Remplacer la structure Token

**Fichier** : `tsd/rete/fact_token.go`

**Nouvelle structure** :
```go
// Token représente un ensemble de faits liés par des bindings immuables.
type Token struct {
    ID       string
    Facts    []*Fact
    Bindings *BindingChain  // ✅ Nouvelle structure immuable
    NodeID   string
    Metadata TokenMetadata
}

// TokenMetadata contient des informations de traçage pour le debugging.
type TokenMetadata struct {
    CreatedAt    time.Time
    CreatedBy    string   // ID du nœud créateur
    JoinLevel    int      // Profondeur de cascade (0 = fait initial)
    ParentTokens []string // IDs des tokens parents
}
```

**Ajouter les helpers** :
```go
// NewToken crée un token vide.
func NewToken(nodeID string) *Token {
    return &Token{
        ID:       generateTokenID(),
        Facts:    []*Fact{},
        Bindings: NewBindingChain(),
        NodeID:   nodeID,
        Metadata: TokenMetadata{
            CreatedAt: time.Now(),
            CreatedBy: nodeID,
            JoinLevel: 0,
        },
    }
}

// NewTokenWithFact crée un token avec un fait initial lié à une variable.
func NewTokenWithFact(fact *Fact, variable string, nodeID string) *Token {
    return &Token{
        ID:       generateTokenID(),
        Facts:    []*Fact{fact},
        Bindings: NewBindingChainWith(variable, fact),
        NodeID:   nodeID,
        Metadata: TokenMetadata{
            CreatedAt: time.Now(),
            CreatedBy: nodeID,
            JoinLevel: 0,
        },
    }
}

// GetBinding retourne le fait lié à une variable.
func (t *Token) GetBinding(variable string) *Fact {
    if t.Bindings == nil {
        return nil
    }
    return t.Bindings.Get(variable)
}

// HasBinding vérifie si une variable est liée.
func (t *Token) HasBinding(variable string) bool {
    if t.Bindings == nil {
        return false
    }
    return t.Bindings.Has(variable)
}

// GetVariables retourne toutes les variables liées.
func (t *Token) GetVariables() []string {
    if t.Bindings == nil {
        return []string{}
    }
    return t.Bindings.Variables()
}
```

---

### Tâche 2 : Fixer les Erreurs de Compilation (90 min)

#### 2.1 Identifier tous les fichiers à modifier

**Commande** :
```bash
go build ./rete/... 2>&1 | tee build_errors.log
```

**Fichiers probablement affectés** :
- `rete/node_join.go`
- `rete/node_alpha.go`
- `rete/node_terminal.go`
- `rete/action_executor_context.go`
- `rete/action_executor_evaluation.go`
- `rete/builder_*.go`
- `rete/network.go`
- Tous les tests `*_test.go`

---

#### 2.2 Patterns de remplacement

**Pattern 1** : Accès direct au map
```go
// ANCIEN
fact := token.Bindings[variable]

// NOUVEAU
fact := token.Bindings.Get(variable)
// OU
fact := token.GetBinding(variable)
```

**Pattern 2** : Affectation du map
```go
// ANCIEN
token.Bindings = make(map[string]*Fact)
token.Bindings[variable] = fact

// NOUVEAU
token.Bindings = NewBindingChain()
token.Bindings = token.Bindings.Add(variable, fact)
```

**Pattern 3** : Itération sur le map
```go
// ANCIEN
for variable, fact := range token.Bindings {
    // ...
}

// NOUVEAU
for _, variable := range token.Bindings.Variables() {
    fact := token.Bindings.Get(variable)
    // ...
}
```

**Pattern 4** : Vérification d'existence
```go
// ANCIEN
if fact, ok := token.Bindings[variable]; ok {
    // ...
}

// NOUVEAU
if token.Bindings.Has(variable) {
    fact := token.Bindings.Get(variable)
    // ...
}
```

**Pattern 5** : Copie/merge de bindings
```go
// ANCIEN
combinedBindings := make(map[string]*Fact)
for k, v := range token1.Bindings {
    combinedBindings[k] = v
}
for k, v := range token2.Bindings {
    combinedBindings[k] = v
}

// NOUVEAU
combinedBindings := token1.Bindings.Merge(token2.Bindings)
```

---

#### 2.3 Modifier node_join.go

**Fichier** : `tsd/rete/node_join.go`

**Fonction `performJoinWithTokens`** - Réécrire complètement :
```go
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
    // Composer les chaînes de bindings
    newBindings := token1.Bindings
    if token2.Bindings != nil {
        newBindings = newBindings.Merge(token2.Bindings)
    }
    
    // Vérifier les conditions de jointure
    if !jn.evaluateJoinConditions(newBindings) {
        return nil
    }
    
    // Créer le token joint
    combinedFacts := append([]*Fact{}, token1.Facts...)
    combinedFacts = append(combinedFacts, token2.Facts...)
    
    joinedToken := &Token{
        ID:       generateTokenID(),
        Facts:    combinedFacts,
        Bindings: newBindings,
        NodeID:   jn.ID,
        Metadata: TokenMetadata{
            CreatedAt:    time.Now(),
            CreatedBy:    jn.ID,
            JoinLevel:    max(token1.Metadata.JoinLevel, token2.Metadata.JoinLevel) + 1,
            ParentTokens: []string{token1.ID, token2.ID},
        },
    }
    
    return joinedToken
}
```

**Fonction `evaluateJoinConditions`** - Adapter la signature :
```go
// Ancienne signature
func (jn *JoinNode) evaluateJoinConditions(bindings map[string]*Fact) bool

// Nouvelle signature
func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool {
    if jn.JoinConditions == nil || len(jn.JoinConditions) == 0 {
        return true
    }
    
    // Adapter chaque accès bindings[var] → bindings.Get(var)
    // ...
}
```

---

#### 2.4 Modifier action_executor_context.go

**Fichier** : `tsd/rete/action_executor_context.go`

**Structure ExecutionContext** :
```go
// ANCIEN
type ExecutionContext struct {
    varCache map[string]*Fact
    // ...
}

// NOUVEAU
type ExecutionContext struct {
    bindings *BindingChain
    // ...
}
```

**Fonction de création** :
```go
func NewExecutionContext(token *Token, ...) *ExecutionContext {
    return &ExecutionContext{
        bindings: token.Bindings,  // Référence à la chaîne immuable
        // ...
    }
}
```

---

#### 2.5 Modifier action_executor_evaluation.go

**Fichier** : `tsd/rete/action_executor_evaluation.go`

**Résolution de variables** :
```go
func (ctx *ExecutionContext) resolveVariable(name string) (interface{}, error) {
    if ctx.bindings != nil && ctx.bindings.Has(name) {
        fact := ctx.bindings.Get(name)
        return fact, nil
    }
    
    // Message d'erreur amélioré
    available := []string{}
    if ctx.bindings != nil {
        available = ctx.bindings.Variables()
    }
    return nil, fmt.Errorf("variable '%s' non trouvée (variables disponibles: %v)", name, available)
}
```

---

#### 2.6 Modifier les autres fichiers

**Pour chaque fichier avec erreur de compilation** :
1. Ouvrir le fichier
2. Localiser l'erreur (accès à `.Bindings[...]`)
3. Appliquer le pattern de remplacement approprié
4. Compiler pour vérifier
5. Passer au suivant

**Commande de vérification continue** :
```bash
# Après chaque modification
go build ./rete/...
```

---

### Tâche 3 : Adapter les Tests (60 min)

#### 3.1 Modifier les tests de Token

**Fichier** : `tsd/rete/fact_token_test.go`

**Adapter les assertions** :
```go
// ANCIEN
if len(token.Bindings) != 2 {
    t.Errorf("...")
}
if token.Bindings["user"] == nil {
    t.Errorf("...")
}

// NOUVEAU
if token.Bindings.Len() != 2 {
    t.Errorf("...")
}
if !token.HasBinding("user") {
    t.Errorf("...")
}
```

---

#### 3.2 Modifier les tests de JoinNode

**Fichier** : `tsd/rete/node_join_test.go`

**Adapter la création de tokens de test** :
```go
// ANCIEN
token := &Token{
    ID: "t1",
    Bindings: map[string]*Fact{
        "user": userFact,
    },
}

// NOUVEAU
token := &Token{
    ID: "t1",
    Bindings: NewBindingChain().Add("user", userFact),
}
// OU
token := NewTokenWithFact(userFact, "user", "test_node")
```

**Adapter les vérifications** :
```go
// ANCIEN
if result.Bindings["user"] != userFact {
    t.Errorf("...")
}

// NOUVEAU
if result.GetBinding("user") != userFact {
    t.Errorf("...")
}
```

---

#### 3.3 Modifier les autres tests

**Pour chaque fichier de test** :
1. Identifier les créations de Token
2. Remplacer les maps par BindingChain
3. Adapter les assertions
4. Exécuter le test :
   ```bash
   go test -v ./rete/[fichier_test].go
   ```

---

### Tâche 4 : Validation Complète (30 min)

#### 4.1 Compiler tout le projet

**Commandes** :
```bash
cd tsd

# Compilation complète
go build ./...

# Vérifier qu'il n'y a aucune erreur
echo $?  # Doit retourner 0
```

**Si erreurs** : Retourner à la Tâche 2 et fixer les fichiers manqués.

---

#### 4.2 Exécuter les tests unitaires

**Commandes** :
```bash
# Tests du module rete
go test -v ./rete/...

# Tests avec couverture
go test -cover ./rete/...
```

**Résultat attendu** : Tous les tests passent (ou la plupart - ajustements mineurs possibles)

---

#### 4.3 Exécuter les tests d'intégration

**Commandes** :
```bash
make test-integration
```

**Si échecs** :
- Analyser les messages d'erreur
- Identifier les patterns manqués
- Fixer et re-tester

---

#### 4.4 Vérifier la qualité du code

**Commandes** :
```bash
# Formattage
go fmt ./rete/...

# Analyse statique
go vet ./rete/...

# Vérifier qu'il n'y a pas de TODO restants
grep -r "TODO\|FIXME" rete/*.go
```

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Compilation
- [ ] ✅ `go build ./...` passe sans erreur
- [ ] ✅ Aucun warning de compilation
- [ ] ✅ Aucun code mort (ancien Token supprimé)

### Tests
- [ ] ✅ `go test ./rete/...` passe
- [ ] ✅ Tests de Token adaptés et passent
- [ ] ✅ Tests de JoinNode adaptés et passent
- [ ] ✅ Tests d'intégration passent (au moins partiellement)

### Qualité
- [ ] ✅ Code formaté (`go fmt`)
- [ ] ✅ Pas de warnings (`go vet`)
- [ ] ✅ Pas de code temporaire ou commenté
- [ ] ✅ Bindings utilisent BindingChain partout

### Migration
- [ ] ✅ Ancien `map[string]*Fact` complètement supprimé
- [ ] ✅ Tous les fichiers utilisent la nouvelle structure
- [ ] ✅ Pas de cohabitation ancien/nouveau code

---

## 🎯 Résultats Attendus

### Fichiers Modifiés (liste non exhaustive)
- `rete/fact_token.go` - Structure Token modifiée
- `rete/node_join.go` - JoinNode adapté
- `rete/action_executor_context.go` - ExecutionContext adapté
- `rete/action_executor_evaluation.go` - Résolution adaptée
- `rete/node_terminal.go` - Terminal adapté
- Tous les tests `*_test.go` adaptés

### Tests Passants
- ✅ Tests de BindingChain (déjà validés au Prompt 03)
- ✅ Tests de Token (adaptés)
- ✅ Tests de JoinNode (adaptés, mais cascades 3+ peuvent encore échouer - normal)
- ⚠️ Tests E2E : certains peuvent échouer (seront fixés dans les prompts suivants)

---

## 🎯 Prochaine Étape

Une fois Token **refactoré et validé**, passer au **Prompt 05 - JoinNode performJoinWithTokens**.

Le Prompt 05 optimisera la logique de jointure pour garantir que tous les bindings sont préservés.

---

## 💡 Conseils Pratiques

### Pour la Migration
1. **Travailler fichier par fichier** : Ne pas essayer de tout fixer d'un coup
2. **Compiler fréquemment** : Après chaque fichier modifié
3. **Tester fréquemment** : Après chaque groupe de modifications
4. **Utiliser grep** : Pour trouver tous les usages restants

### Pour les Erreurs
1. **Lire attentivement** : Le compilateur Go donne des messages précis
2. **Chercher les patterns** : Souvent la même erreur se répète
3. **Documenter** : Noter les patterns de remplacement pour référence
4. **Ne pas paniquer** : C'est normal d'avoir beaucoup d'erreurs au début

### Pour les Tests
1. **Adapter progressivement** : Un fichier de test à la fois
2. **Vérifier les helpers** : Souvent des fonctions utilitaires à adapter
3. **Messages clairs** : Améliorer les messages d'erreur si besoin
4. **Couverture** : Vérifier que la couverture reste bonne

---

## ⚠️ Points d'Attention

### Risques
1. **Oubli de fichiers** : Utiliser grep pour trouver tous les usages
2. **Tests cassés** : Certains tests peuvent révéler des bugs - c'est bien !
3. **Performance** : À ce stade, ne pas optimiser - juste faire fonctionner

### Ne PAS Faire
- ❌ Garder l'ancien code "au cas où"
- ❌ Ajouter des flags pour basculer entre ancien/nouveau
- ❌ Ignorer les warnings du compilateur
- ❌ Committer du code qui ne compile pas

### À Faire
- ✅ Supprimer complètement l'ancien code
- ✅ Fixer TOUTES les erreurs de compilation
- ✅ Adapter TOUS les tests
- ✅ Vérifier que ça compile et teste avant de passer au suivant

---

**Note** : Cette session est la plus critique - c'est la **migration big bang**. Prenez le temps nécessaire, travaillez méthodiquement, et ne passez au Prompt 05 que quand tout compile et que les tests de base passent.