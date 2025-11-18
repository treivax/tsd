# Système TUPLE-SPACE pour le Réseau RETE

## 📋 Vue d'Ensemble

Le système RETE a été modifié pour fonctionner comme un **tuple-space** où :
- Les **nœuds terminaux** stockent les ensembles de faits déclencheurs au lieu d'exécuter directement les actions
- Les **agents** du tuple-space peuvent "prendre" ces tuples pour déclencher les actions de manière asynchrone
- Chaque **action disponible** est affichée avec ses **faits déclencheurs** entre parenthèses

## 🏗️ Architecture Modifiée

### Avant (Exécution Directe)
```
Fait → TypeNode → AlphaNode → TerminalNode → [EXÉCUTION IMMÉDIATE]
```

### Après (Tuple-Space)
```
Fait → TypeNode → AlphaNode → TerminalNode → [STOCKAGE + AFFICHAGE]
                                                    ↓
                                            Agents Tuple-Space
                                            "prennent" les tuples
                                                    ↓
                                            [EXÉCUTION DIFFÉRÉE]
```

## 🔧 Modifications Techniques

### 1. Fonction `executeAction` Modifiée

**Fichier :** `rete/rete.go` (lignes 453-481)

**Avant :**
```go
func (tn *TerminalNode) executeAction(token *Token) error {
    // Exécution silencieuse de l'action
    // (logs désactivés)
    return nil
}
```

**Après :**
```go
func (tn *TerminalNode) executeAction(token *Token) error {
    // === VERSION TUPLE-SPACE ===
    // Au lieu d'exécuter l'action, on l'affiche avec les faits déclencheurs

    actionName := tn.Action.Job.Name
    fmt.Printf("🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: %s", actionName)

    // Afficher les faits déclencheurs entre parenthèses
    if len(token.Facts) > 0 {
        fmt.Print(" (")
        for i, fact := range token.Facts {
            if i > 0 {
                fmt.Print(", ")
            }
            fmt.Printf("%s[", fact.Type)
            for key, value := range fact.Fields {
                fmt.Printf("%s=%v", key, value)
            }
            fmt.Print("]")
        }
        fmt.Print(")")
    }
    fmt.Println()

    return nil
}
```

### 2. Support Évaluateur Étendu

**Fichier :** `rete/evaluator.go`

**Ajouts pour la compatibilité :**
- Support de `"binary_op"` en plus de `"binaryOperation"`
- Support de `"logical_op"` en plus de `"logicalExpression"`
- Support de `"field_access"` en plus de `"fieldAccess"`
- Support du format `"op"` en plus de `"operator"`
- Support du format `"variable"` en plus de `"object"` pour l'accès aux champs

## 🧪 Test Validé

**Fichier :** `tests/real_parsing_test.go` (fonction `TestTupleSpaceTerminalNodes`)

### Scénario de Test
1. **Client majeur (age=25)** → Déclenche `authorize_customer`
2. **Client mineur (age=16)** → Ne déclenche rien
3. **Autre client majeur (age=30)** → Déclenche `authorize_customer`

### Résultats Attendus ✅
```
🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: authorize_customer (Customer[id=C001, age=25, vip=true])
🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: authorize_customer (Customer[id=C003, age=30, vip=false])

📋 ANALYSE DU TUPLE-SPACE:
  Terminal: terminal_authorize (Action: authorize_customer)
  Tokens stockés: 2
    Token 1: 1 faits déclencheurs - Client C001 (age=25)
    Token 2: 1 faits déclencheurs - Client C003 (age=30)
```

## 🎯 Fonctionnalité Tuple-Space

### Stockage des Ensembles de Faits
- Chaque **TerminalNode** maintient une mémoire (`Memory.Tokens`) des ensembles de faits déclencheurs
- Chaque **Token** contient un ou plusieurs **Facts** qui ont satisfait les conditions de la règle
- Les tokens sont **persistés** jusqu'à ce qu'un agent les "prenne"

### Format d'Affichage
```
🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: <nom_action> (<fait1>, <fait2>, ...)
```

**Exemple avec un seul fait :**
```
🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: authorize_customer (Customer[id=C001, age=25, vip=true])
```

**Exemple avec plusieurs faits (jointures) :**
```
🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: process_order (Customer[id=C001, age=25], Order[id=O001, total=1500])
```

## 🔄 Flux de Fonctionnement

### 1. Soumission de Fait
```go
network.SubmitFact(customerFact)
```

### 2. Propagation dans le Réseau
```
RootNode → TypeNode → AlphaNode → TerminalNode
```

### 3. Évaluation des Conditions
- **AlphaNode** évalue les conditions sur le fait
- Si **conditions satisfaites** → création d'un **Token**
- **Token** propagé vers le **TerminalNode**

### 4. Stockage Tuple-Space
- **TerminalNode** reçoit le token via `ActivateLeft(token)`
- Token **stocké** dans `Memory.Tokens[token.ID] = token`
- **Action affichée** avec faits déclencheurs

### 5. Consommation par les Agents (À implémenter)
Les agents du tuple-space peuvent :
- **Lister** les actions disponibles dans `network.TerminalNodes`
- **Prendre** un token spécifique
- **Exécuter** l'action correspondante
- **Supprimer** le token du tuple-space

## 📊 Avantages du Système

### ✅ Séparation des Préoccupations
- **Moteur RETE** : Évaluation des règles et détection des patterns
- **Tuple-Space** : Stockage temporaire des actions déclenchées
- **Agents** : Exécution asynchrone et distribuée

### ✅ Flexibilité d'Exécution
- **Priorités** : Les agents peuvent traiter les actions par priorité
- **Parallélisme** : Plusieurs agents peuvent consommer simultanément
- **Résilience** : Actions persistées jusqu'à traitement complet

### ✅ Monitoring et Debug
- **Visibilité** : Toutes les actions déclenchées sont visibles
- **Traçabilité** : Chaque action est liée à ses faits déclencheurs
- **État** : État du tuple-space consultable à tout moment

## 🚀 Étapes Suivantes

### Phase 2 - Agents Consommateurs
1. **Interface Agent** : Définir l'API pour les agents
2. **Take Operation** : Implémenter la prise de tuples atomique
3. **Concurrence** : Gestion des accès concurrents au tuple-space

### Phase 3 - Distribution
1. **Réseau** : Distribution des tuples sur plusieurs nœuds
2. **Persistence** : Sauvegarde des tuples non traités
3. **Récupération** : Gestion des pannes et reprises

---

## ✅ Status Actuel

**IMPLÉMENTATION COMPLÈTE** de la première étape :
- ✅ Stockage des ensembles de faits déclencheurs
- ✅ Affichage des actions avec faits en format tuple-space
- ✅ Tests validés avec règles simples et complexes
- ✅ Architecture prête pour l'ajout des agents consommateurs
