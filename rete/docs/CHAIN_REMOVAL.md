// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

# Chain Removal - Gestion de la Suppression des Chaînes d'AlphaNodes

## 🎯 Vue d'Ensemble

La fonctionnalité de **suppression de chaînes** permet de supprimer correctement les règles RETE qui utilisent des chaînes d'AlphaNodes, tout en préservant les nœuds partagés entre plusieurs règles.

### Problématique

Lorsqu'une règle avec une expression AND est décomposée en chaîne d'AlphaNodes :
```
Rule 1: p.age > 18 AND p.salary >= 50000
→ TypeNode → AlphaNode(age) → AlphaNode(salary) → Terminal
```

Si une autre règle partage certains nœuds :
```
Rule 2: p.age > 18 AND p.experience > 5
→ TypeNode → AlphaNode(age) → AlphaNode(experience) → Terminal
```

Le nœud `AlphaNode(age)` est partagé. Lors de la suppression de Rule 1, ce nœud doit être **conservé** car Rule 2 l'utilise encore.

### Solution

Un algorithme de suppression intelligent qui :
1. Détecte les chaînes d'AlphaNodes
2. Remonte la chaîne en ordre inverse
3. Décrémenter le RefCount de chaque nœud
4. Supprime uniquement les nœuds avec RefCount == 0
5. S'arrête dès qu'un nœud partagé est rencontré

## 🏗️ Architecture

### Flux de Suppression

```
RemoveRule(ruleID)
    ↓
Détection de chaîne ? ───No──→ removeSimpleRule() (comportement classique)
    ↓ Yes
removeAlphaChain(ruleID)
    ↓
1. Identifier les nœuds (Terminal, Alpha, autres)
    ↓
2. Supprimer Terminal
    ↓
3. Ordonner AlphaNodes en ordre inverse
    ↓
4. Pour chaque AlphaNode (du terminal vers TypeNode):
    ├─ Décrémenter RefCount
    ├─ Si RefCount == 0: Supprimer
    ├─ Si RefCount > 0: Marquer arrêt suppressions
    └─ Continuer décrémentation des parents
    ↓
5. Supprimer autres nœuds (TypeNodes, etc.)
```

### Composants Clés

#### 1. `RemoveRule(ruleID string) error`

Fonction principale de suppression :
- Détecte si la règle utilise une chaîne
- Délègue à `removeAlphaChain()` ou `removeSimpleRule()`
- Point d'entrée unique pour toutes les suppressions

#### 2. `removeAlphaChain(ruleID string) error`

Gestion spécialisée des chaînes :
- Identifie les nœuds par type
- Ordonne les AlphaNodes en ordre inverse
- Supprime intelligemment avec gestion du partage
- Log détaillé de chaque opération

#### 3. `removeSimpleRule(ruleID string, nodeIDs []string) error`

Comportement classique pour règles simples :
- Parcourt tous les nœuds
- Supprime si RefCount == 0
- Maintient la rétrocompatibilité

#### 4. `isPartOfChain(nodeID string) bool`

Détection de chaînes :
- Vérifie si un AlphaNode a un parent AlphaNode
- Vérifie si un AlphaNode a un enfant AlphaNode
- Retourne true si l'une des conditions est vraie

#### 5. `getChainParent(alphaNode *AlphaNode) Node`

Récupération du parent :
- Cherche dans les TypeNodes
- Cherche dans les autres AlphaNodes
- Retourne le nœud parent ou nil

#### 6. `orderAlphaNodesReverse(alphaNodeIDs []string) []string`

Ordonnancement intelligent :
- Construit un graphe parent→enfant
- Trouve le nœud terminal de la chaîne
- Remonte la chaîne en ordre inverse
- Gère les cas dégénérés

## 📊 Exemples d'Utilisation

### Exemple 1 : Chaîne Unique (Suppression Complète)

```go
// Règle: p.age > 18 AND p.salary >= 50000
// Aucun partage avec d'autres règles

err := network.RemoveRule("rule_unique")

// Résultat:
// ✅ 2 AlphaNodes supprimés
// ✅ 1 TerminalNode supprimé
// ✅ Nœuds supprimés de AlphaSharingManager et LifecycleManager
```

**Logs**:
```
🗑️  Suppression de la règle: rule_unique
   📊 Nœuds associés à la règle: 3
   🔗 Chaîne d'AlphaNodes détectée, utilisation de la suppression optimisée
   🗑️  TerminalNode rule_unique_terminal supprimé
   🔗 AlphaNode alpha_xxx déconnecté de son parent
   ✓ AlphaNode alpha_xxx supprimé du AlphaSharingManager
   🗑️  AlphaNode alpha_xxx supprimé (position 2 dans la chaîne)
   🔗 AlphaNode alpha_yyy déconnecté de son parent
   ✓ AlphaNode alpha_yyy supprimé du AlphaSharingManager
   🗑️  AlphaNode alpha_yyy supprimé (position 1 dans la chaîne)
✅ Règle rule_unique avec chaîne supprimée avec succès (3 nœud(s) supprimé(s))
```

### Exemple 2 : Partage Partiel

```go
// Rule 1: p.age > 18 AND p.salary >= 50000
// Rule 2: p.age > 18 AND p.experience > 5
// Partage: AlphaNode(age > 18)

err := network.RemoveRule("rule_1")

// Résultat:
// ✅ AlphaNode(salary) supprimé (RefCount 0)
// ♻️  AlphaNode(age) conservé (RefCount 1)
// ✅ TerminalNode supprimé
```

**Logs**:
```
🗑️  Suppression de la règle: rule_1
   📊 Nœuds associés à la règle: 3
   🔗 Chaîne d'AlphaNodes détectée, utilisation de la suppression optimisée
   🗑️  TerminalNode rule_1_terminal supprimé
   🗑️  AlphaNode alpha_salary supprimé (position 2 dans la chaîne)
   ♻️  AlphaNode alpha_age conservé (1 référence(s) restante(s)) - arrêt des suppressions
   ℹ️  Décrémentation du RefCount des nœuds parents partagés
✅ Règle rule_1 avec chaîne supprimée avec succès (2 nœud(s) supprimé(s))
```

### Exemple 3 : Partage Complet

```go
// Rule 1: p.age > 18 AND p.salary >= 50000
// Rule 2: p.age > 18 AND p.salary >= 50000 (même condition)
// Partage: Tous les AlphaNodes

err := network.RemoveRule("rule_1")

// Résultat:
// ♻️  Tous les AlphaNodes conservés (RefCount décrémenté)
// ✅ Seul le TerminalNode supprimé
```

**Logs**:
```
🗑️  Suppression de la règle: rule_1
   📊 Nœuds associés à la règle: 3
   🔗 Chaîne d'AlphaNodes détectée, utilisation de la suppression optimisée
   🗑️  TerminalNode rule_1_terminal supprimé
   ♻️  AlphaNode alpha_salary conservé (1 référence(s) restante(s)) - arrêt des suppressions
   ℹ️  Décrémentation du RefCount des nœuds parents partagés
   ♻️  AlphaNode alpha_age: RefCount décrémenté (1 référence(s) restante(s))
✅ Règle rule_1 avec chaîne supprimée avec succès (1 nœud(s) supprimé(s))
```

### Exemple 4 : Chaînes Indépendantes

```go
// Rule 1: p.age > 18 AND p.salary >= 50000
// Rule 2: p.name == "John" AND p.city == "NYC"
// Aucun partage

// Supprimer Rule 1
err := network.RemoveRule("rule_1")
// ✅ Nœuds de Rule 1 supprimés

// Vérifier Rule 2
nodes := network.LifecycleManager.GetNodesForRule("rule_2")
// ✅ Tous les nœuds de Rule 2 intacts

// Supprimer Rule 2
err = network.RemoveRule("rule_2")
// ✅ Nœuds de Rule 2 supprimés
// ✅ Réseau complètement nettoyé
```

## 🔍 Détection de Chaînes

### Algorithme `isPartOfChain()`

Un AlphaNode fait partie d'une chaîne si :

1. **Son parent est un AlphaNode** :
   ```
   AlphaNode(parent) → AlphaNode(current) → ...
   ```

2. **Un de ses enfants est un AlphaNode** :
   ```
   ... → AlphaNode(current) → AlphaNode(child)
   ```

### Exemples

**Chaîne détectée** :
```
TypeNode → AlphaNode(age) → AlphaNode(salary) → Terminal
           ↑ parent != Alpha    ↑ enfant = Alpha
           → Fait partie         → Fait partie
```

**Pas de chaîne** :
```
TypeNode → AlphaNode(age) → Terminal
           ↑ parent != Alpha  ↑ enfant != Alpha
           → PAS de chaîne
```

## 🔄 Ordonnancement Inverse

### Pourquoi l'ordre inverse ?

Pour supprimer correctement une chaîne, on doit :
1. Commencer par les nœuds terminaux (enfants)
2. Remonter vers les nœuds sources (parents)
3. S'arrêter dès qu'un nœud partagé est trouvé

### Algorithme `orderAlphaNodesReverse()`

```
1. Construire graphe parent→enfant
   Pour chaque AlphaNode:
     - Identifier son parent AlphaNode
     - Enregistrer la relation

2. Trouver le nœud terminal
   - Chercher un nœud qui n'est parent de personne
   - C'est le dernier nœud de la chaîne

3. Remonter la chaîne
   - Partir du terminal
   - Suivre les relations parent
   - Construire la liste ordonnée
```

### Exemple

**Chaîne** :
```
TypeNode → A → B → C → Terminal
```

**Graphe parent→enfant** :
```
B → A (B est enfant de A)
C → B (C est enfant de B)
```

**Nœud terminal** : C (n'est parent de personne)

**Ordre inverse** : [C, B, A]

**Suppression** :
```
1. C: RefCount 0 → Supprimer
2. B: RefCount 0 → Supprimer
3. A: RefCount 1 → Conserver, arrêter suppressions, décrémenter
```

## 📝 Gestion des RefCounts

### Principe

Chaque nœud a un `RefCount` qui indique combien de règles l'utilisent.

**Création** :
```go
// Règle 1 utilise le nœud
RefCount = 1

// Règle 2 utilise le même nœud (partage)
RefCount = 2
```

**Suppression** :
```go
// Supprimer Règle 1
RefCount = 2 → 1 (décrémenté)
// Nœud conservé (RefCount > 0)

// Supprimer Règle 2
RefCount = 1 → 0 (décrémenté)
// Nœud supprimé (RefCount == 0)
```

### Décrémentation Continue

**Important** : Même quand on arrête les suppressions (nœud partagé trouvé), on continue à **décrémenter** les RefCounts des parents.

**Pourquoi ?**

```
Règle 1: A → B → C
Règle 2: A → B → D

Supprimer Règle 1:
1. C: RefCount 1→0, Supprimer ✓
2. B: RefCount 2→1, Conserver, ARRÊTER suppressions
3. A: RefCount 2→1, DÉCRÉMENTER quand même (sinon RefCount incorrect!)

Si on ne décrémentait pas A:
- A aurait RefCount = 2 alors que seule Règle 2 l'utilise
- Problème lors de la suppression de Règle 2
```

## 🧪 Tests

### Suite de Tests

| Test | Objectif | Vérifications |
|------|----------|---------------|
| `TestRemoveChain_AllNodesUnique_DeletesAll` | Chaîne unique | Tous les nœuds supprimés |
| `TestRemoveChain_PartialSharing_DeletesOnlyUnused` | Partage partiel | Nœuds non partagés supprimés |
| `TestRemoveChain_CompleteSharing_DeletesNone` | Partage complet | Aucun AlphaNode supprimé |
| `TestRemoveRule_WithChain_CorrectCleanup` | Nettoyage complet | Tous les registres nettoyés |
| `TestRemoveRule_MultipleChains_IndependentCleanup` | Chaînes indépendantes | Suppressions isolées |
| `TestRemoveRule_SimpleCondition_BackwardCompatibility` | Rétrocompatibilité | Règles simples OK |

### Exécution

```bash
# Tous les tests de suppression
go test ./rete -v -run "TestRemove"

# Tests spécifiques aux chaînes
go test ./rete -v -run "TestRemoveChain_"

# Test spécifique
go test ./rete -v -run "TestRemoveChain_PartialSharing"
```

### Résultats Attendus

```
=== RUN   TestRemoveChain_AllNodesUnique_DeletesAll
    ✓ Chaîne unique supprimée complètement
--- PASS: TestRemoveChain_AllNodesUnique_DeletesAll (0.00s)

=== RUN   TestRemoveChain_PartialSharing_DeletesOnlyUnused
    ✓ Suppression partielle correcte
--- PASS: TestRemoveChain_PartialSharing_DeletesOnlyUnused (0.00s)

=== RUN   TestRemoveChain_CompleteSharing_DeletesNone
    ✓ Partage complet: RefCount correctement décrémenté
--- PASS: TestRemoveChain_CompleteSharing_DeletesNone (0.00s)

PASS
```

## 🔒 Garanties

### Sécurité

✅ **Pas d'orphelins** : Aucun nœud n'est laissé sans référence  
✅ **Pas de suppressions intempestives** : Les nœuds partagés sont préservés  
✅ **RefCount cohérent** : Toujours synchronisé avec l'état réel  
✅ **Nettoyage complet** : Suppression de tous les registres  

### Cohérence

✅ **AlphaNodes** : Supprimés de `network.AlphaNodes`  
✅ **TerminalNodes** : Supprimés de `network.TerminalNodes`  
✅ **AlphaSharingManager** : Nœuds désindexés  
✅ **LifecycleManager** : Nœuds et références supprimés  
✅ **Connexions** : Parents/enfants déconnectés proprement  

### Performance

✅ **Complexité** : O(n) où n = nombre de nœuds de la chaîne  
✅ **Arrêt précoce** : Dès qu'un nœud partagé est trouvé  
✅ **Pas de parcours inutile** : Ordonnancement intelligent  

## 🐛 Débogage

### Logging Détaillé

Chaque opération est loguée avec un emoji identifiant :

| Emoji | Signification |
|-------|---------------|
| 🗑️ | Suppression de règle/nœud |
| 🔗 | Détection/gestion de chaîne |
| ✓ | Opération réussie |
| ♻️ | Nœud partagé conservé |
| ℹ️ | Information |
| ⚠️ | Avertissement |
| ✅ | Succès final |
| 📊 | Statistiques |

### Problèmes Courants

#### 1. Nœud non supprimé

**Symptôme** :
```
Nœud encore présent après RemoveRule()
```

**Diagnostic** :
```bash
# Vérifier le RefCount
lifecycle, _ := network.LifecycleManager.GetNodeLifecycle(nodeID)
fmt.Printf("RefCount: %d\n", lifecycle.GetRefCount())
```

**Cause probable** : Nœud partagé avec une autre règle

#### 2. RefCount incorrect

**Symptôme** :
```
RefCount ne correspond pas au nombre de règles
```

**Diagnostic** :
```bash
# Lister les règles utilisant le nœud
rules := lifecycle.GetRules()
fmt.Printf("Règles: %v\n", rules)
```

**Cause probable** : Bug dans l'incrémentation/décrémentation

#### 3. Orphelins

**Symptôme** :
```
Nœuds présents dans network.AlphaNodes mais pas dans LifecycleManager
```

**Diagnostic** :
```bash
# Vérifier la cohérence
for nodeID := range network.AlphaNodes {
    if _, exists := network.LifecycleManager.GetNodeLifecycle(nodeID); !exists {
        fmt.Printf("ORPHELIN: %s\n", nodeID)
    }
}
```

**Cause probable** : Bug dans la synchronisation des registres

## 📚 API Reference

### `RemoveRule(ruleID string) error`

Supprime une règle et tous ses nœuds non partagés.

**Paramètres** :
- `ruleID` : Identifiant unique de la règle

**Retour** :
- `error` : nil si succès, erreur sinon

**Exemple** :
```go
err := network.RemoveRule("rule_123")
if err != nil {
    log.Fatalf("Erreur suppression: %v", err)
}
```

### `isPartOfChain(nodeID string) bool`

Détecte si un nœud fait partie d'une chaîne d'AlphaNodes.

**Paramètres** :
- `nodeID` : Identifiant du nœud

**Retour** :
- `bool` : true si partie d'une chaîne

**Exemple** :
```go
if network.isPartOfChain("alpha_xxx") {
    fmt.Println("Ce nœud est dans une chaîne")
}
```

### `getChainParent(alphaNode *AlphaNode) Node`

Récupère le nœud parent d'un AlphaNode.

**Paramètres** :
- `alphaNode` : Pointeur vers l'AlphaNode

**Retour** :
- `Node` : Nœud parent ou nil

**Exemple** :
```go
parent := network.getChainParent(alphaNode)
if parent != nil {
    fmt.Printf("Parent: %s\n", parent.GetID())
}
```

## 🎯 Bonnes Pratiques

### 1. Toujours utiliser RemoveRule()

❌ **Mauvais** :
```go
// Suppression manuelle
delete(network.AlphaNodes, nodeID)
```

✅ **Bon** :
```go
// Suppression via l'API
err := network.RemoveRule(ruleID)
```

### 2. Vérifier les erreurs

❌ **Mauvais** :
```go
network.RemoveRule("rule_123")
// Pas de vérification
```

✅ **Bon** :
```go
if err := network.RemoveRule("rule_123"); err != nil {
    log.Printf("Erreur suppression: %v", err)
}
```

### 3. Utiliser les logs pour diagnostiquer

✅ **Bon** :
```go
// Les logs détaillés permettent de suivre la suppression
network.RemoveRule("rule_123")
// Lire les logs pour comprendre ce qui s'est passé
```

### 4. Tester avec différents scénarios

✅ **Bon** :
```go
// Test 1: Chaîne unique
// Test 2: Partage partiel
// Test 3: Partage complet
// Test 4: Chaînes indépendantes
```

## 🔮 Évolutions Futures

### Court Terme
- [ ] Métriques de suppression (temps, nœuds supprimés)
- [ ] Mode dry-run (simuler sans supprimer)
- [ ] Validation avant suppression

### Moyen Terme
- [ ] Suppression batch de plusieurs règles
- [ ] Récupération automatique d'orphelins
- [ ] Statistiques d'utilisation des nœuds

### Long Terme
- [ ] Garbage collector pour nœuds non référencés
- [ ] Compaction automatique du réseau
- [ ] Optimisation des chaînes lors de la suppression

## 📄 Licence

```
Copyright (c) 2025 TSD Contributors
Licensed under the MIT License
```

Voir le fichier [LICENSE](../../../LICENSE) pour les détails complets.

---

**Version** : 1.0.0  
**Date** : 2025-01-27  
**Status** : ✅ Production Ready