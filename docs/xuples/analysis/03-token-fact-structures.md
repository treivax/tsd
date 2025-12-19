# Analyse des Structures Token et Fact - TSD

> ⚠️ **ATTENTION** : Cette documentation décrit l'ancienne architecture (v1.x) avant la migration vers `_id_`.  
> Pour la version actuelle (v2.0+), voir [docs/internal-ids.md](../../internal-ids.md).

## 📋 Vue d'Ensemble

Ce document analyse en détail les structures `Token` et `Fact` qui portent les données dans le réseau RETE et comment elles sont liées.

## 🎯 Objectif

Comprendre comment les faits et tokens sont structurés, comment ils interagissent, et quelles sont les implications pour la création de xuples (tuples enrichis).

---

## 1. Structure Fact

### 1.1 Définition Complète

**Emplacement** : `rete/fact_token.go` lignes 16-26

```go
// Fact représente un fait dans le réseau RETE
type Fact struct {
	// ID est l'identifiant unique du fait.
	// Il est soit généré à partir des clés primaires, soit calculé comme hash.
	// Format: "TypeName~value1_value2..." ou "TypeName~<hash>"
	// Accessible dans les expressions via le champ spécial 'id'.
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Fields     map[string]interface{} `json:"fields"`
	Attributes map[string]interface{} `json:"attributes,omitempty"` // Alias pour Fields (compatibilité)
}
```

**Champs** :
- **ID** : Identifiant unique du fait (calculé ou fourni)
- **Type** : Type du fait (ex: "User", "Order", "Product")
- **Fields** : Map des champs et leurs valeurs
- **Attributes** : Alias de Fields (rétrocompatibilité)

### 1.2 Génération de l'ID

**Format de l'ID** :
- **Avec clés primaires** : `"TypeName~value1_value2_..."`
- **Sans clés primaires** : `"TypeName~<hash>"` (hash calculé des champs)

**Exemple avec clés primaires** :
```go
// Type défini: type User(#id: string, name: string, age: number)
// Fait: User(id: "U001", name: "Alice", age: 25)
fact := &Fact{
    ID:     "U001",
    Type:   "User",
    Fields: map[string]interface{}{
        "id": "U001",
        "name": "Alice",
        "age": 25,
    },
}
```

**Exemple sans clés primaires** :
```go
// Type défini: type Event(message: string, timestamp: number)
// Fait: Event(message: "Login", timestamp: 1234567890)
fact := &Fact{
    ID:     "<hash_calculé>",  // Hash des champs
    Type:   "Event",
    Fields: map[string]interface{}{
        "message": "Login",
        "timestamp": 1234567890,
    },
}
```

### 1.3 Identifiant Interne (Type_ID)

**Méthode** : `GetInternalID()`

**Emplacement** : `rete/fact_token.go` lignes 33-35

```go
// GetInternalID retourne l'identifiant interne unique (Type_ID)
func (f *Fact) GetInternalID() string {
	return fmt.Sprintf("%s_%s", f.Type, f.ID)
}
```

**Format** : `"Type_ID"`

**Exemples** :
- User avec ID "U001" → `"User_U001"`
- Order avec ID "O123" → `"Order_O123"`
- Event avec hash → `"Event_<hash>"`

**Utilité** :
- **Indexation** : Clé unique dans WorkingMemory
- **Évite collisions** : Deux types différents peuvent avoir le même ID simple
- **Rétractation** : Identification précise du fait à retirer

### 1.4 Champ Spécial "id" (Obsolète v1.x)

> ⚠️ **OBSOLÈTE** : Dans la version v2.0+, le champ `id` a été remplacé par `_id_` (interne et caché).  
> Voir [docs/internal-ids.md](../../internal-ids.md) pour la nouvelle architecture.

**Constante (v1.x)** : `FieldNameID = "id"`  
**Constante (v2.0+)** : `FieldNameInternalID = "_id_"`

**Emplacement** : `rete/fact_token.go` lignes 11-13 (v1.x)

```go
// V1.X - OBSOLÈTE
// FieldNameID est le nom du champ spécial pour l'identifiant du fait.
// Ce champ est accessible dans les expressions mais stocké dans Fact.ID, pas dans Fact.Fields.
const FieldNameID = "id"

// V2.0+ - ACTUEL
// FieldNameInternalID est le nom du champ interne pour l'identifiant du fait.
// Ce champ est CACHÉ et NON accessible dans les expressions TSD.
const FieldNameInternalID = "_id_"
```

**Changements v2.0** :
- L'ID est **caché** et inaccessible dans les expressions TSD
- Génération **automatique** obligatoire (pas d'affectation manuelle)
- Format : `"TypeName~value1_value2..."` ou `"TypeName~<hash>"`
- Utilisation de clés primaires (`#field`) pour identification
- Support des types de faits dans les champs

**Exemple d'utilisation (v1.x - OBSOLÈTE)** :
```tsd
rule check_user: {u: User} / u.id == "U001" ==> print(u.id)
                                 ↑ accès au champ spécial
```

**Exemple d'utilisation (v2.0+ - ACTUEL)** :
```tsd
type User(#email: string, name: string)
alice = User("alice@example.com", "Alice")

rule check_user: {u: User} / u.email == "alice@example.com" ==> print(u.name)
                                 ↑ utiliser les clés primaires, pas _id_
```

### 1.5 Méthodes Utiles

```go
// String retourne la représentation string d'un fait
func (f *Fact) String() string

// GetField retourne la valeur d'un champ
func (f *Fact) GetField(fieldName string) (interface{}, bool)

// Clone crée une copie profonde d'un fait
func (f *Fact) Clone() *Fact
```

**Référence** : `rete/fact_token.go` lignes 28-60

---

## 2. Structure Token

### 2.1 Définition Complète

**Emplacement** : `rete/fact_token.go` lignes 86-98

```go
// Token représente un token dans le réseau RETE avec bindings immuables.
//
// Changement majeur: Bindings utilise maintenant BindingChain au lieu de map[string]*Fact
// pour garantir l'immutabilité et éviter la perte de bindings lors des jointures en cascade.
type Token struct {
	ID           string        `json:"id"`
	Facts        []*Fact       `json:"facts"`
	NodeID       string        `json:"node_id"`
	Parent       *Token        `json:"parent,omitempty"`
	Bindings     *BindingChain `json:"-"`                        // Chaîne immuable de bindings (non sérialisable)
	IsJoinResult bool          `json:"is_join_result,omitempty"` // Indique si c'est un token de jointure réussie
	Metadata     TokenMetadata `json:"metadata,omitempty"`       // Métadonnées pour traçage
}
```

**Champs** :
- **ID** : Identifiant unique du token (format: `token_<counter>`)
- **Facts** : Liste des faits associés au token
- **NodeID** : ID du nœud qui a créé ce token
- **Parent** : Token parent dans la chaîne (historique)
- **Bindings** : Chaîne immuable de bindings (variable → fact)
- **IsJoinResult** : Flag indiquant si c'est un résultat de jointure
- **Metadata** : Métadonnées de traçage (timestamps, créateur, etc.)

### 2.2 Génération de l'ID

**Fonction** : `generateTokenID()`

**Emplacement** : `rete/fact_token.go` lignes 328-338

```go
// generateTokenID génère un ID unique pour un token.
//
// Format: "token_<timestamp>_<counter>"
// Cette fonction utilise un compteur atomique pour garantir l'unicité.
var tokenCounter uint64

func generateTokenID() string {
	// Utiliser un compteur atomique simple pour l'unicité
	// Dans une implémentation production, utiliser atomic.AddUint64
	tokenCounter++
	return fmt.Sprintf("token_%d", tokenCounter)
}
```

**Format** : `"token_<counter>"`

**Note** : Le commentaire mentionne `timestamp_counter` mais le code utilise seulement `counter`

**Thread-safety** : ⚠️ Pas thread-safe actuellement (devrait utiliser `atomic.AddUint64`)

### 2.3 BindingChain (Immuable)

**Concept** : Remplace l'ancienne `map[string]*Fact` par une structure immuable

**Avantages** :
1. **Immutabilité** : Pas de modification après création
2. **Partage structurel** : Plusieurs tokens peuvent partager une même chaîne
3. **Pas de perte de bindings** : Résout les bugs de jointures en cascade
4. **Thread-safe** : Pas besoin de synchronisation

**Structure conceptuelle** :
```
BindingChain (linked list immuable)
    │
    ├─ var1 → Fact1
    │   └─> next → BindingChain
    │              │
    │              ├─ var2 → Fact2
    │              │   └─> next → BindingChain
    │              │              │
    │              │              └─ var3 → Fact3
    │              │                  └─> next → nil
```

**Méthodes d'accès** (via Token) :

```go
// GetBinding retourne le fait lié à une variable
func (t *Token) GetBinding(variable string) *Fact

// HasBinding vérifie si une variable est liée dans ce token
func (t *Token) HasBinding(variable string) bool

// GetVariables retourne toutes les variables liées dans ce token
func (t *Token) GetVariables() []string
```

**Référence** : `rete/fact_token.go` lignes 282-325

**Implémentation BindingChain** : (fichier séparé, non fourni mais référencé)

### 2.4 TokenMetadata

**Emplacement** : `rete/fact_token.go` lignes 78-84

```go
type TokenMetadata struct {
	CreatedAt    string   `json:"created_at,omitempty"`    // Timestamp de création
	CreatedBy    string   `json:"created_by,omitempty"`    // ID du nœud créateur
	JoinLevel    int      `json:"join_level,omitempty"`    // Niveau de jointure (0 = fait initial, 1+ = jointures)
	ParentTokens []string `json:"parent_tokens,omitempty"` // IDs des tokens parents (pour jointures)
}
```

**Champs** :
- **CreatedAt** : Timestamp de création (format string)
- **CreatedBy** : ID du nœud RETE créateur
- **JoinLevel** : Profondeur de jointure (0 = token initial, 1+ = jointures)
- **ParentTokens** : IDs des tokens parents (utile pour débugger jointures)

**Exemple** :
```go
// Token créé par un TypeNode
token1 := Token{
    ID: "token_1",
    Facts: []*Fact{userFact},
    Metadata: TokenMetadata{
        CreatedBy: "type_node_user",
        JoinLevel: 0,  // Token initial
    },
}

// Token créé par un JoinNode
token2 := Token{
    ID: "token_2",
    Facts: []*Fact{userFact, orderFact},
    Metadata: TokenMetadata{
        CreatedBy: "join_node_123",
        JoinLevel: 1,  // Première jointure
        ParentTokens: []string{"token_1"},
    },
}
```

---

## 3. Relation entre Token et Fact

### 3.1 Token Contient Plusieurs Facts

**Principe** :
- Un Token peut contenir **1 ou plusieurs** faits
- Chaque fait correspond à une variable matchée
- Les faits s'accumulent lors des jointures

**Exemple** :

```tsd
rule user_order: {u: User, o: Order} / u.id == o.user_id ==> print(u.name, o.id)
```

**Token résultant** :
```go
token := &Token{
    ID: "token_123",
    Facts: []*Fact{
        &Fact{ID: "U001", Type: "User", Fields: {...}},   // Fait pour variable 'u'
        &Fact{ID: "O456", Type: "Order", Fields: {...}},  // Fait pour variable 'o'
    },
    Bindings: chainWith("u", userFact).Add("o", orderFact),
}
```

### 3.2 Bindings : Variable → Fact

**Principe** :
- Chaque variable de la règle est liée à un fait spécifique
- Les bindings sont stockés dans une BindingChain immuable
- Accessible via `token.GetBinding(varName)`

**Exemple d'utilisation** :
```go
// Dans ActionExecutor, évaluation de l'argument "u.name"
userFact := token.GetBinding("u")  // Récupère le fait lié à 'u'
if userFact != nil {
    name, _ := userFact.GetField("name")
    // Utiliser name...
}
```

### 3.3 Évolution d'un Token à Travers les Nœuds

**Scénario** : Règle avec 2 patterns

```tsd
rule user_order: {u: User, o: Order} / u.id == o.user_id ==> print(u.name, o.id)
```

**Flux** :

```
1. TypeNode (User)
   ├─> Crée Token1
   │   ├─ Facts: [UserFact]
   │   ├─ Bindings: {u → UserFact}
   │   └─ Metadata: {JoinLevel: 0, CreatedBy: "type_node_user"}
   │
2. JoinNode (User x Order)
   ├─> Reçoit Token1 (left) et OrderFact (right)
   ├─> Évalue condition: u.id == o.user_id
   ├─> Si match, crée Token2
   │   ├─ Facts: [UserFact, OrderFact]
   │   ├─ Bindings: {u → UserFact, o → OrderFact}
   │   ├─ Parent: Token1
   │   └─ Metadata: {JoinLevel: 1, CreatedBy: "join_node_123", ParentTokens: ["token_1"]}
   │
3. TerminalNode
   ├─> Reçoit Token2
   ├─> Stocke Token2 dans Memory.Tokens
   └─> Exécute action avec Token2
       ├─> Évalue u.name → "Alice" (via Bindings)
       └─> Évalue o.id → "O456" (via Bindings)
```

---

## 4. WorkingMemory : Stockage Facts et Tokens

### 4.1 Structure

**Emplacement** : `rete/fact_token.go` lignes 100-105

```go
type WorkingMemory struct {
	NodeID string            `json:"node_id"`
	Facts  map[string]*Fact  `json:"facts"`
	Tokens map[string]*Token `json:"tokens"`
}
```

**Champs** :
- **NodeID** : ID du nœud propriétaire de cette mémoire
- **Facts** : Map des faits indexés par identifiant interne (Type_ID)
- **Tokens** : Map des tokens indexés par token ID

### 4.2 Indexation des Facts

**Clé utilisée** : Identifiant interne `Type_ID`

**Méthode d'ajout** :

```go
// AddFact ajoute un fait à la mémoire en utilisant un identifiant interne unique (Type_ID)
// Retourne une erreur si un fait avec le même type et ID existe déjà
func (wm *WorkingMemory) AddFact(fact *Fact) error {
	if wm.Facts == nil {
		wm.Facts = make(map[string]*Fact)
	}

	// Utiliser l'identifiant interne (Type_ID) pour garantir l'unicité par type
	internalID := fact.GetInternalID()

	if existingFact, exists := wm.Facts[internalID]; exists {
		return fmt.Errorf("fait avec ID '%s' et type '%s' existe déjà dans la mémoire du nœud %s (champs existants: %v)",
			fact.ID, fact.Type, wm.NodeID, existingFact.Fields)
	}

	wm.Facts[internalID] = fact
	return nil
}
```

**Référence** : `rete/fact_token.go` lignes 107-124

**Garanties** :
- ✅ Unicité par type : `User_U001` ≠ `Order_U001`
- ✅ Détection de doublons : Erreur si déjà présent
- ✅ Thread-safety : À condition d'utiliser mutex au niveau appelant

### 4.3 Indexation des Tokens

**Clé utilisée** : `token.ID`

**Méthodes** :
```go
func (wm *WorkingMemory) AddToken(token *Token)
func (wm *WorkingMemory) RemoveToken(tokenID string)
func (wm *WorkingMemory) GetTokens() []*Token
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token
```

**Référence** : `rete/fact_token.go` lignes 160-228

**Particularité GetTokensByVariable** :
```go
// GetTokensByVariable retourne les tokens contenant au moins une des variables spécifiées.
// Si variables est vide ou nil, retourne tous les tokens.
//
// Le filtrage est basé sur Token.Bindings.Has() pour vérifier la présence de chaque variable.
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token {
	// Si pas de filtre, retourner tous les tokens
	if len(variables) == 0 {
		return wm.GetTokens()
	}

	// Filtrer les tokens qui contiennent au moins une des variables
	result := make([]*Token, 0)
	for _, token := range wm.Tokens {
		if token.Bindings != nil {
			for _, varName := range variables {
				if token.Bindings.Has(varName) {
					result = append(result, token)
					break // Token déjà ajouté, passer au suivant
				}
			}
		}
	}

	return result
}
```

**Utilité** : Permet de rechercher les tokens affectés par une variable spécifique

---

## 5. Création de Tokens

### 5.1 NewTokenWithFact (Token Initial)

**Emplacement** : `rete/fact_token.go` lignes 340-370

```go
// NewTokenWithFact crée un nouveau token avec un seul binding.
//
// Fonction utilitaire pour créer un token initial avec un fait unique,
// typiquement utilisé lors de la première activation d'un JoinNode.
//
// Paramètres:
//   - fact: pointeur vers le fait à lier
//   - variable: nom de la variable à lier au fait
//   - nodeID: ID du nœud créateur du token
//
// Retourne:
//   - *Token: nouveau token avec le binding spécifié
//
// Exemple:
//
//	userFact := &Fact{ID: "u1", Type: "User", Fields: map[string]interface{}{"id": 1}}
//	token := NewTokenWithFact(userFact, "user", "type_node_user")
//	fmt.Println(token.HasBinding("user"))  // true
//	fmt.Println(token.GetBinding("user") == userFact)  // true
func NewTokenWithFact(fact *Fact, variable string, nodeID string) *Token {
	return &Token{
		ID:       generateTokenID(),
		Facts:    []*Fact{fact},
		NodeID:   nodeID,
		Bindings: NewBindingChainWith(variable, fact),
		Metadata: TokenMetadata{
			CreatedBy: nodeID,
			JoinLevel: 0,
		},
	}
}
```

**Usage typique** : Création d'un token initial par un TypeNode ou AlphaNode

### 5.2 Clone d'un Token

**Emplacement** : `rete/fact_token.go` lignes 251-280

```go
// Clone crée une copie profonde d'un token.
//
// Note: BindingChain est immuable donc pas besoin de cloner la chaîne elle-même.
// On réutilise la même référence (partage structurel).
func (t *Token) Clone() *Token {
	clone := &Token{
		ID:           t.ID,
		Facts:        make([]*Fact, len(t.Facts)),
		NodeID:       t.NodeID,
		Bindings:     t.Bindings, // Immuable, pas besoin de cloner
		IsJoinResult: t.IsJoinResult,
		Metadata:     t.Metadata, // Copie de la structure
	}

	// Copier les faits
	for i, fact := range t.Facts {
		clone.Facts[i] = fact.Clone()
	}

	// Copier les ParentTokens si présents
	if len(t.Metadata.ParentTokens) > 0 {
		clone.Metadata.ParentTokens = make([]string, len(t.Metadata.ParentTokens))
		copy(clone.Metadata.ParentTokens, t.Metadata.ParentTokens)
	}

	// Note: Parent n'est pas cloné pour éviter récursion infinie
	// Note: Bindings n'est pas cloné car BindingChain est immuable

	return clone
}
```

**Points importants** :
- ✅ Copie profonde des Facts
- ✅ Partage de BindingChain (immuable)
- ✅ Parent non cloné (évite récursion)
- ✅ Copie des ParentTokens (slice)

---

## 6. Implications pour Xuples

### 6.1 Xuples = Tokens Enrichis ?

**Concept** :
- Un xuple pourrait être un Token avec métadonnées supplémentaires
- Conserver la structure Token (excellente base)
- Ajouter informations spécifiques au tuple-space

**Structure proposée** :
```go
type Xuple struct {
    Token    *Token                 // Token RETE original
    Action   *Action                // Action déclenchée
    RuleID   string                 // ID de la règle
    Status   XupleStatus            // Status: pending, executing, executed, failed
    Created  time.Time              // Timestamp de création
    Updated  time.Time              // Timestamp dernière modification
    Metadata map[string]interface{} // Métadonnées additionnelles
}

type XupleStatus string

const (
    XupleStatusPending   XupleStatus = "pending"
    XupleStatusExecuting XupleStatus = "executing"
    XupleStatusExecuted  XupleStatus = "executed"
    XupleStatusFailed    XupleStatus = "failed"
)
```

### 6.2 Schéma de Données

```
Xuple
  │
  ├─ Token
  │   ├─ ID: string
  │   ├─ Facts: []*Fact
  │   │   ├─ Fact 1
  │   │   │   ├─ ID: string
  │   │   │   ├─ Type: string
  │   │   │   └─ Fields: map[string]interface{}
  │   │   └─ Fact 2
  │   │       └─ ...
  │   ├─ Bindings: *BindingChain
  │   │   ├─ "u" → Fact 1
  │   │   └─ "o" → Fact 2
  │   └─ Metadata: TokenMetadata
  │
  ├─ Action: *Action
  │   └─ Jobs: []JobCall
  │       └─ JobCall
  │           ├─ Name: string
  │           └─ Args: []interface{}
  │
  ├─ RuleID: string
  ├─ Status: XupleStatus
  └─ Metadata: map[string]interface{}
```

### 6.3 Conversion Token → Xuple

```go
func NewXuple(token *Token, action *Action, ruleID string) *Xuple {
    return &Xuple{
        Token:    token,
        Action:   action,
        RuleID:   ruleID,
        Status:   XupleStatusPending,
        Created:  time.Now(),
        Updated:  time.Now(),
        Metadata: make(map[string]interface{}),
    }
}
```

### 6.4 Index et Recherche

**Indices proposés** :
- Par RuleID : `map[string][]*Xuple`
- Par Status : `map[XupleStatus][]*Xuple`
- Par ActionName : `map[string][]*Xuple`
- Par Variable : `map[string][]*Xuple` (via Token.Bindings)

**Requêtes possibles** :
```go
// Trouver tous les xuples d'une règle
xuples := space.GetByRule("user_order")

// Trouver tous les xuples pending
pending := space.GetByStatus(XupleStatusPending)

// Trouver xuples avec variable 'u'
withUser := space.GetByVariable("u")
```

---

## 7. Fonctions Utilitaires

### 7.1 MakeInternalID et ParseInternalID

**Emplacement** : `rete/fact_token.go` lignes 63-76

```go
// MakeInternalID construit un identifiant interne à partir d'un type et d'un ID
func MakeInternalID(factType, factID string) string {
	return fmt.Sprintf("%s_%s", factType, factID)
}

// ParseInternalID décompose un identifiant interne en type et ID
// Retourne (type, id, true) si le format est valide, sinon ("", "", false)
func ParseInternalID(internalID string) (string, string, bool) {
	for i := 0; i < len(internalID); i++ {
		if internalID[i] == '_' {
			return internalID[:i], internalID[i+1:], true
		}
	}
	return "", "", false
}
```

**Utilité** :
- Créer clés pour indexation
- Parser clés pour extraction type/ID
- Validation format

**Exemple** :
```go
internalID := MakeInternalID("User", "U001")  // "User_U001"
factType, factID, ok := ParseInternalID(internalID)
// factType = "User", factID = "U001", ok = true
```

---

## 8. Synthèse

### 8.1 Points Forts

✅ **Structure Fact** : Simple, efficace, bien documentée  
✅ **BindingChain immuable** : Résout problèmes de jointures, thread-safe  
✅ **TokenMetadata** : Excellente traçabilité  
✅ **Identifiants internes** : Évite collisions, unicité garantie  
✅ **WorkingMemory** : Indexation claire Facts et Tokens  
✅ **Méthodes utilitaires** : GetBinding, HasBinding, GetVariables

### 8.2 Observations

⚠️ **generateTokenID** : Pas thread-safe (devrait utiliser atomic)  
⚠️ **Champ "id" virtuel** : Documentation claire mais peut prêter à confusion  
⚠️ **Attributes vs Fields** : Alias pour compatibilité, pourrait être simplifié

### 8.3 Recommandations pour Xuples

1. **Conserver structures** : Token et Fact sont excellentes
2. **Enrichir avec Xuple** : Wrapper autour de Token
3. **Ajouter status et lifecycle** : pending, executing, executed, failed
4. **Indexer par multiples critères** : RuleID, Status, ActionName, Variables
5. **Conserver immutabilité** : Ne pas modifier Token après création
6. **Fix generateTokenID** : Utiliser `atomic.AddUint64` pour thread-safety

---

## 9. Fichiers de Référence

| Fichier | Description | Lignes clés |
|---------|-------------|-------------|
| `rete/fact_token.go` | Structures Fact, Token, WorkingMemory | 16-60 (Fact), 78-98 (Token/Metadata), 100-280 (WorkingMemory) |

---

**Date de création** : 2025-12-17  
**Auteur** : Analyse automatique pour refonte xuples  
**Statut** : ✅ Complet
