# 05 - Comportement actuel des Terminal Nodes

## 🎯 Vue d'ensemble

Ce document analyse le comportement actuel des terminal nodes dans le moteur RETE, avant la modification pour l'exécution immédiate.

## 📊 Architecture actuelle

### Flux de traitement actuel

```
Token arrive → TerminalNode.ActivateLeft()
                    ↓
            Stockage dans Memory.Tokens
                    ↓
            Exécution de l'action via executeAction()
                    ↓
            ActionExecutor.ExecuteAction()
                    ↓
            Actions handlers (print, etc.)
```

### Composants impliqués

#### 1. TerminalNode (rete/node_terminal.go)

**Structure** :
```go
type TerminalNode struct {
    BaseNode
    Action *Action
}
```

**Méthodes clés** :
- `ActivateLeft(token *Token)` : Point d'entrée quand un token arrive
- `executeAction(token *Token)` : Exécute l'action associée
- `GetTriggeredActions()` : Retourne les actions déclenchées (pour tests)
- `ActivateRetract(factID string)` : Gère la rétractation de faits

#### 2. WorkingMemory (dans BaseNode)

**Stockage** :
```go
type WorkingMemory struct {
    NodeID string
    Facts  map[string]*Fact
    Tokens map[string]*Token  // ← Stockage actuel des tokens
}
```

#### 3. ActionExecutor (rete/action_executor.go)

**Interface actuelle** :
```go
type ActionExecutor struct {
    network       *ReteNetwork
    logger        *log.Logger
    enableLogging bool
    registry      *ActionRegistry
}

func (ae *ActionExecutor) ExecuteAction(action *Action, token *Token) error
```

**Process** :
1. Valide l'action et le token
2. Obtient les jobs à exécuter
3. Crée un ExecutionContext avec les bindings
4. Exécute chaque job en séquence avec récupération sur panic

## 🔍 Analyse détaillée du flux actuel

### 1. ActivateLeft - Point d'entrée

**Code actuel** (lignes 46-62) :
```go
func (tn *TerminalNode) ActivateLeft(token *Token) error {
    // 1. Enregistrer l'activation
    tn.recordActivation()

    // 2. STOCKER LE TOKEN dans la mémoire
    tn.mutex.Lock()
    if tn.Memory.Tokens == nil {
        tn.Memory.Tokens = make(map[string]*Token)
    }
    tn.Memory.Tokens[token.ID] = token
    tn.mutex.Unlock()

    // 3. Exécuter l'action
    return tn.executeAction(token)
}
```

**Comportement** :
- ✅ Enregistre l'activation (métriques)
- ⚠️ **STOCKE le token** dans la mémoire du nœud
- ✅ Exécute l'action immédiatement

**Problème identifié** :
- Le stockage n'est PAS nécessaire pour l'exécution
- Il sert uniquement à `collectActivations()` pour récupérer les activations après coup
- Cela crée un couplage entre RETE (moteur) et la récupération d'activations

### 2. executeAction - Exécution de l'action

**Code actuel** (lignes 128-153) :
```go
func (tn *TerminalNode) executeAction(token *Token) error {
    if tn.Action == nil {
        return fmt.Errorf("aucune action définie pour le nœud %s", tn.ID)
    }

    // TODO(xuples): Publier vers XupleSpace si configuré
    // [Code commenté pour future intégration xuples]

    // Exécuter l'action avec l'ActionExecutor
    network := tn.BaseNode.GetNetwork()
    if network != nil && network.ActionExecutor != nil {
        return network.ActionExecutor.ExecuteAction(tn.Action, token)
    }

    return nil
}
```

**Comportement** :
- ✅ Vérifie qu'une action existe
- ✅ Délègue à `ActionExecutor.ExecuteAction()`
- ⚠️ TODO pour intégration xuples (non implémenté)

### 3. collectActivations - Récupération des activations

**Code actuel** (internal/servercmd/servercmd.go) :
```go
func (s *Server) collectActivations(network *rete.ReteNetwork) []tsdio.Activation {
    if network == nil {
        return []tsdio.Activation{}
    }

    activations := []tsdio.Activation{}

    // Parcourt TOUS les terminal nodes
    for _, terminal := range network.TerminalNodes {
        if terminal.Memory == nil || terminal.Memory.Tokens == nil {
            continue
        }

        actionName := "unknown"
        if terminal.Action != nil && terminal.Action.Job != nil {
            actionName = terminal.Action.Job.Name
        }

        // Parcourt TOUS les tokens stockés
        for _, token := range terminal.Memory.Tokens {
            activation := tsdio.Activation{
                ActionName:      actionName,
                Arguments:       s.extractArguments(terminal, token),
                TriggeringFacts: s.extractFacts(token),
                BindingsCount:   len(token.Facts),
            }
            activations = append(activations, activation)
        }
    }

    return activations
}
```

**Comportement** :
- Parcourt tous les terminal nodes du réseau
- Lit `terminal.Memory.Tokens` pour chaque terminal node
- Extrait les informations d'activation
- Retourne la liste complète des activations

**Problème** :
- **Couplage fort** : Le serveur dépend de la structure interne de TerminalNode
- **Violation de l'encapsulation** : Accès direct à `Memory.Tokens`
- **Redondance** : Les informations sont déjà disponibles au moment de l'exécution

## 📈 Diagramme de séquence actuel

```
┌──────────┐     ┌──────────────┐     ┌───────────────┐     ┌────────────────┐
│  Network │     │ TerminalNode │     │ ActionExecutor│     │ Server/Tests   │
└────┬─────┘     └──────┬───────┘     └───────┬───────┘     └───────┬────────┘
     │                  │                     │                     │
     │ ActivateLeft     │                     │                     │
     ├─────────────────>│                     │                     │
     │                  │                     │                     │
     │                  │ recordActivation()  │                     │
     │                  ├──────────┐          │                     │
     │                  │          │          │                     │
     │                  │<─────────┘          │                     │
     │                  │                     │                     │
     │                  │ STORE token in      │                     │
     │                  │ Memory.Tokens       │                     │
     │                  ├──────────┐          │                     │
     │                  │          │          │                     │
     │                  │<─────────┘          │                     │
     │                  │                     │                     │
     │                  │ executeAction()     │                     │
     │                  ├──────────┐          │                     │
     │                  │          │          │                     │
     │                  │<─────────┘          │                     │
     │                  │                     │                     │
     │                  │ ExecuteAction()     │                     │
     │                  ├────────────────────>│                     │
     │                  │                     │                     │
     │                  │                     │ Execute handlers    │
     │                  │                     ├──────────┐          │
     │                  │                     │          │          │
     │                  │                     │<─────────┘          │
     │                  │                     │                     │
     │                  │<────────────────────┤                     │
     │                  │                     │                     │
     │<─────────────────┤                     │                     │
     │                  │                     │                     │
     .                  .                     .                     .
     .                  .                     .                     .
     │                  │                     │                     │
     │                  │                     │ collectActivations()│
     │                  │<────────────────────────────────────────┤
     │                  │                     │                     │
     │                  │ READ Memory.Tokens  │                     │
     │                  ├──────────┐          │                     │
     │                  │          │          │                     │
     │                  │<─────────┘          │                     │
     │                  │                     │                     │
     │                  ├─────────────────────────────────────────>│
     │                  │ Return activations  │                     │
     │                  │                     │                     │
```

## ⚠️ Points d'attention identifiés

### 1. Stockage redondant
- **Problème** : Les tokens sont stockés mais l'information pourrait être capturée au moment de l'exécution
- **Impact** : Consommation mémoire inutile, complexité ajoutée
- **Solution** : Observer pattern pour capturer les exécutions sans stockage

### 2. Couplage Server ↔ TerminalNode
- **Problème** : `collectActivations()` accède directement à la structure interne
- **Impact** : Violation de l'encapsulation, difficile à maintenir
- **Solution** : Interface d'observation propre

### 3. Récupération a posteriori
- **Problème** : Les activations sont collectées après l'exécution
- **Impact** : Impossible de filtrer ou transformer en temps réel
- **Solution** : Observer pattern avec callback immédiat

### 4. Absence d'observabilité en temps réel
- **Problème** : Pas de notification immédiate des exécutions
- **Impact** : Impossible d'intégrer avec xuples ou autres systèmes
- **Solution** : Observer pattern + hooks d'exécution

## 🎯 Usages actuels de Memory.Tokens

### 1. Tests
```bash
$ grep -r "Memory.Tokens" rete/*.go | wc -l
15
```

**Fichiers concernés** :
- Tests unitaires de `TerminalNode`
- Tests d'intégration de règles
- Vérifications d'activations

### 2. Serveur
- `internal/servercmd/servercmd.go` : `collectActivations()`
- `internal/servercmd/servercmd_test.go` : Tests du serveur

### 3. GetTriggeredActions()
**Code actuel** (lignes 88-97) :
```go
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

**Problème** :
- Retourne une copie de l'action pour chaque token
- Utilisé uniquement dans les tests

## 💡 Opportunités d'amélioration

### 1. Séparation des responsabilités
- **RETE** : Moteur de règles pur, exécution immédiate
- **Xuples** : Système de coordination, gestion des activations
- **Observer** : Pont entre les deux, découplage

### 2. Performance
- Éliminer le stockage des tokens = moins de mémoire
- Éliminer `collectActivations()` = moins de parcours du réseau
- Observer pattern = notification directe, pas de polling

### 3. Extensibilité
- Observer pattern permet multiples observateurs
- Facile d'ajouter logging, métriques, xuples, etc.
- Découplage total des consommateurs

### 4. Testabilité
- Mock observers pour tests
- Vérification immédiate des exécutions
- Statistiques accessibles via observer

## 📊 Métriques actuelles

### Complexité
- `ActivateLeft` : **8 lignes** (simple)
- `executeAction` : **15 lignes** (simple)
- `collectActivations` : **25 lignes** (peut être remplacé)

### Dépendances
- TerminalNode → ActionExecutor ✅
- TerminalNode → Network ✅
- Server → TerminalNode.Memory ⚠️ (couplage fort)
- Tests → TerminalNode.Memory ⚠️ (couplage fort)

## 🚀 Prochaine étape

Voir `06-immediate-execution-design.md` pour la conception du nouveau comportement avec exécution immédiate et observer pattern.
