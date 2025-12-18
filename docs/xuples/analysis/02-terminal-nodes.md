# Analyse des Terminal Nodes - TSD

## 📋 Vue d'Ensemble

Ce document analyse en profondeur l'implémentation actuelle des Terminal Nodes dans le réseau RETE, leur rôle dans le stockage des tokens et l'exécution des actions.

## 🎯 Objectif

Comprendre comment les Terminal Nodes fonctionnent actuellement, comment ils stockent les tokens matchés, et comment ils déclenchent l'exécution des actions.

---

## 1. Architecture Actuelle des Terminal Nodes

### 1.1 Structure du TerminalNode

**Emplacement** : `rete/node_terminal.go` lignes 12-30

```go
type TerminalNode struct {
	BaseNode
	Action *Action `json:"action"`
}

// NewTerminalNode crée un nouveau nœud terminal
func NewTerminalNode(nodeID string, action *Action, storage Storage) *TerminalNode {
	return &TerminalNode{
		BaseNode: BaseNode{
			ID:        nodeID,
			Type:      "terminal",
			Memory:    &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children:  make([]Node, 0), // Les nœuds terminaux n'ont pas d'enfants
			Storage:   storage,
			createdAt: time.Now(),
		},
		Action: action,
	}
}
```

**Composition** :
- **BaseNode** : Hérite des fonctionnalités de base (Memory, ID, Type, Children, Storage)
- **Action** : Pointeur vers l'action à exécuter (structure `*Action` de constraint)

### 1.2 BaseNode (hérité)

Le `BaseNode` contient les éléments essentiels :

```go
type BaseNode struct {
	ID        string
	Type      string
	Memory    *WorkingMemory
	Children  []Node
	Storage   Storage
	createdAt time.Time
	mutex     sync.RWMutex
	// ... autres champs métier
}
```

**WorkingMemory** : Structure clé pour stocker Facts et Tokens

---

## 2. Cycle de Vie d'un Token

### 2.1 Activation par ActivateLeft

**Emplacement** : `rete/node_terminal.go` lignes 32-62

```go
// ActivateLeft déclenche l'action lorsqu'un token arrive.
//
// Process :
//  1. Stocke le token dans la mémoire du nœud
//  2. Exécute l'action associée avec le contexte du token
//
// Le token contient tous les bindings (via BindingChain) nécessaires
// pour l'évaluation des arguments de l'action.
//
// Paramètres :
//   - token : token contenant les faits et bindings déclencheurs
//
// Retourne :
//   - error : erreur si l'exécution de l'action échoue
func (tn *TerminalNode) ActivateLeft(token *Token) error {
	// Enregistrer l'activation
	tn.recordActivation()

	// Stocker le token
	tn.mutex.Lock()
	if tn.Memory.Tokens == nil {
		tn.Memory.Tokens = make(map[string]*Token)
	}
	tn.Memory.Tokens[token.ID] = token
	tn.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Déclencher l'action
	return tn.executeAction(token)
}
```

**Étapes** :
1. **Enregistrement** : `recordActivation()` incrémente un compteur d'activations
2. **Stockage** : Le token est ajouté à `Memory.Tokens` avec protection mutex
3. **Persistance** : Commentaire indique que c'est désactivé pour performances
4. **Exécution** : Appel immédiat de `executeAction(token)`

### 2.2 Stockage des Tokens

**Structure WorkingMemory** : `rete/fact_token.go` lignes 100-105

```go
type WorkingMemory struct {
	NodeID string            `json:"node_id"`
	Facts  map[string]*Fact  `json:"facts"`
	Tokens map[string]*Token `json:"tokens"`
}
```

**Clés utilisées** :
- `Tokens` : map avec `token.ID` comme clé
- Type : `map[string]*Token`
- **Indexation** : Par ID du token uniquement

**Thread-Safety** :
- Protection par `mutex.Lock()` / `mutex.Unlock()`
- Mutex de type `sync.RWMutex` dans BaseNode

### 2.3 Rétractation (ActivateRetract)

**Emplacement** : `rete/node_terminal.go` lignes 64-85

```go
// ActivateRetract retrait des tokens contenant le fait rétracté
// factID doit être l'identifiant interne (Type_ID)
func (tn *TerminalNode) ActivateRetract(factID string) error {
	tn.mutex.Lock()
	var tokensToRemove []string
	for tokenID, token := range tn.Memory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				tokensToRemove = append(tokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range tokensToRemove {
		delete(tn.Memory.Tokens, tokenID)
	}
	tn.mutex.Unlock()
	if len(tokensToRemove) > 0 {
		fmt.Printf("🗑️  [TERMINAL_%s] Rétractation: %d tokens retirés\n", tn.ID, len(tokensToRemove))
	}
	return nil
}
```

**Process** :
1. Parcourt tous les tokens stockés
2. Pour chaque token, vérifie si un de ses faits correspond au factID rétracté
3. Collecte les IDs des tokens à supprimer
4. Supprime les tokens identifiés
5. Affiche un message de log

**Complexité** : O(n*m) où n=nombre de tokens, m=nombre de faits par token

---

## 3. Exécution des Actions

### 3.1 Méthode executeAction

**Emplacement** : `rete/node_terminal.go` lignes 109-172

```go
// executeAction exécute l'action avec le contexte du token.
//
// Process :
//  1. Vérifie qu'une action est définie
//  2. Affiche l'action dans le tuple-space (pour compatibilité)
//  3. Délègue l'exécution au ActionExecutor du réseau
//
// Le ActionExecutor crée un ExecutionContext avec token.Bindings,
// permettant l'accès aux variables via BindingChain.
//
// Paramètres :
//   - token : token contenant les faits et bindings
//
// Retourne :
//   - error : erreur si l'exécution échoue
func (tn *TerminalNode) executeAction(token *Token) error {
	// Les actions sont maintenant obligatoires dans la grammaire
	// Mais nous gardons cette vérification par sécurité
	if tn.Action == nil {
		return fmt.Errorf("aucune action définie pour le nœud %s", tn.ID)
	}

	// Afficher aussi dans tuple-space pour compatibilité
	actionName := "action"
	jobs := tn.Action.GetJobs()
	if len(jobs) > 0 {
		actionName = jobs[0].Name
	}

	// Affichage direct (fmt est déjà thread-safe)
	fmt.Printf("🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: %s", actionName)

	// Afficher les faits déclencheurs entre parenthèses
	if len(token.Facts) > 0 {
		fmt.Print(" (")
		for i, fact := range token.Facts {
			if i > 0 {
				fmt.Print(", ")
			}
			// Format compact : Type(id:value, field:value, ...)
			fmt.Printf("%s(", fact.Type)
			fieldCount := 0
			for key, value := range fact.Fields {
				if fieldCount > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s:%v", key, value)
				fieldCount++
			}
			fmt.Print(")")
		}
		fmt.Print(")")
	}

	fmt.Print("\n")

	// Exécuter réellement l'action avec l'ActionExecutor
	network := tn.BaseNode.GetNetwork()
	if network != nil && network.ActionExecutor != nil {
		return network.ActionExecutor.ExecuteAction(tn.Action, token)
	}

	return nil
}
```

**Étapes** :
1. **Validation** : Vérifie que `tn.Action != nil`
2. **Logging tuple-space** : Affiche l'action disponible (compatibilité ancien système)
3. **Affichage faits** : Format lisible des faits déclencheurs
4. **Délégation** : Appel à `network.ActionExecutor.ExecuteAction()`

**Note importante** : **Deux comportements coexistent** :
- Affichage "tuple-space" (héritage ancien système)
- Exécution réelle via ActionExecutor (système actuel)

### 3.2 Interface avec ActionExecutor

**Emplacement** : `rete/action_executor.go` lignes 100-144

```go
// ExecuteAction exécute une action avec les faits fournis par le token.
//
// Process :
//  1. Valide les paramètres (action et token non nil)
//  2. Récupère tous les jobs de l'action
//  3. Crée un contexte d'exécution avec les bindings du token
//  4. Exécute chaque job en séquence avec récupération sur panic
//
// Thread-Safety :
//   - Cette méthode est thread-safe
//   - Le contexte d'exécution est isolé par appel
//   - Les panics sont récupérés et convertis en erreurs
//
// Paramètres :
//   - action : action à exécuter (peut contenir plusieurs jobs)
//   - token : token contenant les faits et bindings disponibles
//
// Retourne :
//   - error : erreur si l'exécution échoue ou si paramètres invalides
func (ae *ActionExecutor) ExecuteAction(action *Action, token *Token) error {
	if action == nil {
		return fmt.Errorf("action is nil")
	}
	if token == nil {
		return fmt.Errorf("token is nil")
	}

	// Obtenir tous les jobs à exécuter
	jobs := action.GetJobs()

	// Créer un contexte d'exécution avec les faits disponibles
	ctx := NewExecutionContext(token, ae.network)
	if ctx == nil {
		return fmt.Errorf("échec création contexte d'exécution")
	}

	// Exécuter chaque job en séquence
	for i, job := range jobs {
		if err := ae.executeJob(job, ctx, i); err != nil {
			return fmt.Errorf("erreur exécution job %s (index %d): %w", job.Name, i, err)
		}
	}

	return nil
}
```

**Flux d'exécution** :
```
TerminalNode.executeAction(token)
    ↓
ActionExecutor.ExecuteAction(action, token)
    ↓
NewExecutionContext(token, network)
    ↓
Pour chaque job:
    ActionExecutor.executeJob(job, ctx, i)
        ↓
    ActionHandler.Execute(args, ctx)
```

---

## 4. Structure Token et Métadonnées

### 4.1 Structure Token

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

**Champs clés** :
- **ID** : Identifiant unique du token
- **Facts** : Liste des faits associés au token
- **NodeID** : ID du nœud qui a créé le token
- **Parent** : Token parent (chaînage pour historique)
- **Bindings** : Chaîne immuable de bindings (variable → fact)
- **IsJoinResult** : Flag indiquant si c'est le résultat d'une jointure
- **Metadata** : Métadonnées de traçage

### 4.2 TokenMetadata

**Emplacement** : `rete/fact_token.go` lignes 78-84

```go
type TokenMetadata struct {
	CreatedAt    string   `json:"created_at,omitempty"`    // Timestamp de création
	CreatedBy    string   `json:"created_by,omitempty"`    // ID du nœud créateur
	JoinLevel    int      `json:"join_level,omitempty"`    // Niveau de jointure (0 = fait initial, 1+ = jointures)
	ParentTokens []string `json:"parent_tokens,omitempty"` // IDs des tokens parents (pour jointures)
}
```

**Utilité** :
- **Traçage** : Permet de comprendre l'origine d'un token
- **Debug** : Facilite le debug des jointures complexes
- **Audit** : Historique de construction du token

### 4.3 BindingChain (Immuable)

**Concept** : Au lieu d'une map mutable, utilise une chaîne immuable de bindings

**Avantages** :
- **Immutabilité** : Pas de perte de bindings lors des jointures
- **Partage structurel** : Plusieurs tokens peuvent partager une même chaîne
- **Thread-safe** : Pas besoin de synchronisation

**Méthodes** (sur Token) :
```go
func (t *Token) GetBinding(variable string) *Fact
func (t *Token) HasBinding(variable string) bool
func (t *Token) GetVariables() []string
```

**Référence** : `rete/fact_token.go` lignes 282-325

---

## 5. Récupération des Tokens

### 5.1 GetTriggeredActions

**Emplacement** : `rete/node_terminal.go` lignes 87-97

```go
// GetTriggeredActions retourne les actions déclenchées (pour les tests)
func (tn *TerminalNode) GetTriggeredActions() []*Action {
	tn.mutex.RLock()
	defer tn.mutex.RUnlock()

	actions := make([]*Action, 0, len(tn.Memory.Tokens))
	for range tn.Memory.Tokens {
		actions = append(actions, tn.Action)
	}
	return actions
}
```

**Comportement** :
- Retourne l'action du TerminalNode répétée autant de fois qu'il y a de tokens
- **Utilisation** : Tests uniquement
- **Limitation** : Ne retourne pas les tokens eux-mêmes, juste le nombre d'activations

### 5.2 Accès aux Tokens via WorkingMemory

**Méthodes disponibles** sur `WorkingMemory` :

```go
func (wm *WorkingMemory) GetTokens() []*Token
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token
func (wm *WorkingMemory) AddToken(token *Token)
func (wm *WorkingMemory) RemoveToken(tokenID string)
```

**Référence** : `rete/fact_token.go` lignes 159-228

**GetTokensByVariable** :
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

---

## 6. Utilisation Actuelle : collectActivations

### 6.1 Fonction collectActivations

**Emplacement** : `internal/servercmd/servercmd.go` (référencé dans prompt)

**Comportement hypothétique** (basé sur contexte) :
```go
func (network *ReteNetwork) collectActivations() []Activation {
    activations := []Activation{}
    
    // Parcourir tous les terminal nodes
    for _, terminalNode := range network.GetTerminalNodes() {
        // Récupérer tous les tokens stockés
        tokens := terminalNode.Memory.GetTokens()
        
        for _, token := range tokens {
            activation := Activation{
                RuleID: terminalNode.ID,
                Token:  token,
                Action: terminalNode.Action,
            }
            activations = append(activations, activation)
        }
    }
    
    return activations
}
```

**Usage** : 
- Collecte toutes les activations de tous les terminal nodes
- Utilisé pour obtenir l'état global des règles activées
- **Comportement tuple-space actuel** : Les tokens restent en mémoire

---

## 7. Diagramme de Séquence

```
┌─────────┐         ┌──────────────┐       ┌────────────────┐       ┌────────────────┐
│JoinNode │         │TerminalNode  │       │ActionExecutor  │       │ActionHandler   │
└────┬────┘         └──────┬───────┘       └────────┬───────┘       └────────┬───────┘
     │                     │                        │                        │
     │ ActivateLeft(token) │                        │                        │
     ├────────────────────>│                        │                        │
     │                     │                        │                        │
     │                     │ recordActivation()     │                        │
     │                     ├───────────┐            │                        │
     │                     │           │            │                        │
     │                     │<──────────┘            │                        │
     │                     │                        │                        │
     │                     │ mutex.Lock()           │                        │
     │                     ├───────────┐            │                        │
     │                     │           │            │                        │
     │                     │<──────────┘            │                        │
     │                     │                        │                        │
     │                     │ Memory.Tokens[id]=token│                        │
     │                     ├───────────┐            │                        │
     │                     │           │            │                        │
     │                     │<──────────┘            │                        │
     │                     │                        │                        │
     │                     │ mutex.Unlock()         │                        │
     │                     ├───────────┐            │                        │
     │                     │           │            │                        │
     │                     │<──────────┘            │                        │
     │                     │                        │                        │
     │                     │ executeAction(token)   │                        │
     │                     ├───────────┐            │                        │
     │                     │           │            │                        │
     │                     │   Affichage tuple-space│                        │
     │                     │           │            │                        │
     │                     │           │ ExecuteAction(action, token)        │
     │                     │           ├───────────>│                        │
     │                     │           │            │                        │
     │                     │           │            │ NewExecutionContext()  │
     │                     │           │            ├────────────┐           │
     │                     │           │            │            │           │
     │                     │           │            │<───────────┘           │
     │                     │           │            │                        │
     │                     │           │            │ executeJob(job, ctx)   │
     │                     │           │            ├────────────┐           │
     │                     │           │            │            │           │
     │                     │           │            │  evaluateArgument()    │
     │                     │           │            │            │           │
     │                     │           │            │            │ Execute(args, ctx)
     │                     │           │            │            ├──────────>│
     │                     │           │            │            │           │
     │                     │           │            │            │  Exec...  │
     │                     │           │            │            │           │
     │                     │           │            │            │<──────────│
     │                     │           │            │<───────────┘           │
     │                     │           │<───────────┤                        │
     │                     │<──────────┘            │                        │
     │<────────────────────┤                        │                        │
     │                     │                        │                        │
```

---

## 8. Points d'Intervention pour la Refonte

### 8.1 Stockage des Tokens

**État actuel** :
- Tokens stockés indéfiniment dans `Memory.Tokens`
- Pas de nettoyage automatique (sauf rétractation)
- Tous les tokens d'un TerminalNode sont conservés

**Opportunités d'intervention** :
1. **Après exécution** : Décider si le token doit être conservé ou supprimé
2. **Stratégie de rétention** : Configurer durée de vie des tokens
3. **Migration vers xuples** : Copier tokens vers tuple-space externe avant suppression

### 8.2 Exécution des Actions

**État actuel** :
- Exécution immédiate lors de `ActivateLeft`
- Pas de file d'attente d'actions
- Pas de priorisation

**Opportunités d'intervention** :
1. **Découplage** : Séparer stockage token et exécution action
2. **File d'attente** : Ajouter queue d'actions à exécuter
3. **Mode différé** : Permettre exécution batch

### 8.3 Interface ActionExecutor

**État actuel** :
```go
type ActionExecutor interface {
    ExecuteAction(action *Action, token *Token) error
}
```

**Propositions** :
1. Conserver cette interface (très propre)
2. Ajouter callbacks pour notification de fin d'exécution
3. Permettre mode synchrone/asynchrone

### 8.4 Récupération des Tokens (collectActivations)

**État actuel** :
- Parcourt tous les TerminalNodes
- Récupère tous les tokens de chaque Memory
- Retourne tableau d'activations

**Propositions pour xuples** :
1. **Garder tokens en mémoire** : Comportement actuel OK pour tuple-space
2. **Ajouter flag "consumed"** : Marquer tokens exécutés sans les supprimer
3. **Index par action** : Faciliter recherche des activations par type d'action

---

## 9. Problèmes Identifiés

### 9.1 Croissance Mémoire

❌ **Problème** : Tokens jamais supprimés (sauf rétractation)  
⚠️ **Impact** : Croissance continue de la mémoire  
✅ **Solution** : Implémentation de stratégies de rétention dans xuples

### 9.2 Pas de Séparation RETE/Tuple-Space

❌ **Problème** : Terminal nodes font double emploi (RETE + tuple-space)  
⚠️ **Impact** : Confusion des responsabilités  
✅ **Solution** : Module xuples dédié au tuple-space

### 9.3 Affichage Console Hardcodé

❌ **Problème** : `fmt.Printf` dans executeAction (ligne 139)  
⚠️ **Impact** : Couplage avec sortie console  
✅ **Solution** : Utiliser logger configurablepro ou callbacks

---

## 10. Recommandations pour Xuples

### 10.1 Architecture Proposée

```
TerminalNode
    ↓ ActivateLeft(token)
    ├─> Store token in Memory (RETE interne)
    ├─> Execute action
    └─> Publish to xuples (nouveau)
            ↓
        XuplesSpace
            ├─> Store activation (action, token, metadata)
            ├─> Emit events
            └─> Allow queries
```

### 10.2 Conservation

✅ **Garder** :
- Structure Token (excellente)
- BindingChain immutable (parfait pour multi-threading)
- Metadata de traçage (très utile)
- Interface ActionExecutor (propre et extensible)

### 10.3 Modifications Minimales

📝 **Propositions** :
1. Ajouter hook dans `executeAction` pour publication vers xuples
2. Ajouter flag de configuration pour activer/désactiver tuple-space
3. Créer interface `TupleSpacePublisher` pour découplage

```go
type TupleSpacePublisher interface {
    PublishActivation(nodeID string, action *Action, token *Token) error
}
```

---

## 11. Fichiers de Référence

| Fichier | Description | Lignes clés |
|---------|-------------|-------------|
| `rete/node_terminal.go` | Implémentation TerminalNode | 12-30 (struct), 32-62 (ActivateLeft), 109-172 (executeAction) |
| `rete/fact_token.go` | Structures Token, WorkingMemory | 86-98 (Token), 100-105 (WorkingMemory), 78-84 (Metadata) |
| `rete/action_executor.go` | Exécution des actions | 37-69 (struct), 100-144 (ExecuteAction) |
| `rete/action_handler.go` | Registry et interface handlers | 12-25 (interface), 27-135 (registry) |

---

**Date de création** : 2025-12-17  
**Auteur** : Analyse automatique pour refonte xuples  
**Statut** : ✅ Complet
