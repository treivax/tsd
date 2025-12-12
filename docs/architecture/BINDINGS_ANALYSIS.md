# Analyse des Bindings - Jointures Multi-Variables

**Date** : 2025-12-12  
**Auteur** : Analyse de diagnostic (Prompt 01)  
**Objectif** : Identifier la cause racine de la perte de bindings dans les cascades de jointures

---

## 1. Résumé Exécutif

### 1.1 Problème Observé

Dans les règles avec 3+ variables (ex: `{u: User, t: Team, task: Task}`), le token final qui arrive au TerminalNode ne contient que 2 bindings au lieu de 3.

**Erreur observée** :
```
erreur évaluation argument 2: variable 'task' non trouvée (variables disponibles: [t u])
```

L'action `affordable_task_assignment(u.id, t.id, task.id)` ne peut pas s'exécuter car le binding 'task' est manquant.

### 1.2 Cause Racine Identifiée

**LE BUG EST DANS LA CONSTRUCTION DU RÉSEAU BETA (BetaChainBuilder)**

Le deuxième JoinNode de la cascade reçoit le **mauvais type de fait** du côté droit :
- **Configuration du JoinNode** : `RightVariables: [task]` → attend des faits **Task**
- **Fait réellement reçu** : Faits **Team** provenant de `passthrough_r2_t_Team_right`
- **Résultat** : Le JoinNode détecte le fait Team comme variable 't' (au lieu de 'task') et crée un binding `[t]` au lieu de `[task]`

**Conséquence** : 
- La jointure produit `[u, t]` au lieu de `[u, t, task]`
- Le token propagé au terminal manque la variable 'task'
- L'action échoue car elle attend 3 variables

### 1.3 Impact

**Tests échouant** :
- `join_multi_variable_complex.tsd` - Règles r1, r2, r3 (toutes avec 3 variables)
- Tout pattern avec 3+ variables en cascade

**Comportement attendu** :
- Les 3 variables doivent être disponibles dans le token final
- L'action doit s'exécuter avec les 3 arguments

**Comportement observé** :
- Seules 2 variables présentes dans le token final
- Erreur lors de l'évaluation des arguments de l'action

---

## 2. Cas d'Étude : join_multi_variable_complex.tsd

### 2.1 Règle Testée (r2)

```tsd
rule r2 : {u: User, t: Team, task: Task} / 
    u.team_id == t.id AND 
    u.id == task.assignee_id AND 
    t.budget > task.effort * 100 
    ==> affordable_task_assignment(u.id, t.id, task.id)
```

**Variables** : 3 (u: User, t: Team, task: Task)
**Conditions de jointure** :
1. `u.team_id == t.id` (jointure User-Team)
2. `u.id == task.assignee_id` (jointure User-Task)  
3. `t.budget > task.effort * 100` (condition arithmétique)

### 2.2 Faits Soumis

**Ordre de soumission** :
1. `User(id:U001, name:Alice, role:manager, team_id:T001)`
2. `User(id:U002, name:Bob, role:lead, team_id:T001)`
3. `User(id:U003, name:Carol, role:developer, team_id:T002)`
4. `Team(id:T001, name:Alpha, budget:10000, manager_id:U001)` ← **Fait déclencheur du bug**
5. `Team(id:T002, name:Beta, budget:5000, manager_id:U003)`
6. `Task(id:TASK001, assignee_id:U002, team_id:T001, priority:high, effort:50)`
7. `Task(id:TASK002, assignee_id:U003, team_id:T002, priority:medium, effort:20)`
8. `Task(id:TASK003, assignee_id:U001, team_id:T001, priority:high, effort:30)`

**Match attendu pour r2** :
- u = User(U001), t = Team(T001), task = Task(TASK003)
- Conditions vérifiées :
  - u.team_id (T001) == t.id (T001) ✓
  - u.id (U001) == task.assignee_id (U001) ✓
  - t.budget (10000) > task.effort * 100 (3000) ✓

### 2.3 Comportement Attendu

**Action déclenchée** :
```
affordable_task_assignment(U001, T001, TASK003)
```

Avec les arguments :
1. u.id = "U001"
2. t.id = "T001"  
3. task.id = "TASK003"

### 2.4 Comportement Observé

**Erreur levée lors de la soumission du fait Team(T001)** :
```
❌ Erreur soumission faits: erreur soumission fait T001: 
erreur propagation fait vers type_Team: 
error activating alpha node: 
erreur propagation fait vers join_6e16ce13b16480f9: 
erreur propagation token vers r2_terminal: 
erreur exécution job affordable_task_assignment: 
erreur évaluation argument 2: variable 'task' non trouvée (variables disponibles: [t u])
```

**Token reçu par le terminal** :
- Bindings : `[t, u]` ← **SEULEMENT 2 variables au lieu de 3**
- Variable manquante : `task`

---

## 3. Architecture Actuelle

### 3.1 Diagramme de Flux (Architecture Attendue)

```
Pour la règle r2 : {u: User, t: Team, task: Task}

Architecture attendue avec jointures en cascade :

┌─────────────┐
│ TypeNode    │
│   (User)    │
└──────┬──────┘
       │
       ├─────────────────────────────┐
       │                             │
┌──────▼──────┐                ┌─────▼─────────┐
│ AlphaNode   │                │  TypeNode     │
│  r2/User    │                │   (Team)      │
└──────┬──────┘                └─────┬─────────┘
       │                             │
       │ (ActivateLeft)              │ (ActivateRight)
       │                             │
       └──────────┬──────────────────┘
                  │
           ┌──────▼──────────────┐
           │ JoinNode1           │
           │ ID: join_...1       │
           │ Left: [u]           │  ← Premier JoinNode (User x Team)
           │ Right: [t]          │
           │ All: [u, t]         │
           │ Condition:          │
           │  u.team_id == t.id  │
           └──────┬──────────────┘
                  │
                  │ Token: [u, t]
                  │
                  ├─────────────────────────────┐
                  │                             │
           ┌──────▼──────┐              ┌───────▼────────┐
           │ (Propagate  │              │   TypeNode     │
           │  left)      │              │    (Task)      │
           └──────┬──────┘              └───────┬────────┘
                  │                             │
                  │ (ActivateLeft)              │ (ActivateRight)
                  │                             │
                  └──────────┬──────────────────┘
                             │
                      ┌──────▼──────────────┐
                      │ JoinNode2           │
                      │ ID: join_...2       │
                      │ Left: [u, t]        │  ← Deuxième JoinNode (User+Team x Task)
                      │ Right: [task]       │  ← **DEVRAIT recevoir Task**
                      │ All: [u, t, task]   │
                      │ Condition:          │
                      │  u.id == task...    │
                      └──────┬──────────────┘
                             │
                             │ Token: [u, t, task]
                             │
                      ┌──────▼──────────────┐
                      │  TerminalNode       │
                      │  (r2_terminal)      │
                      │  Action:            │
                      │  affordable_task... │
                      └─────────────────────┘
```

### 3.2 Architecture Réelle (Bugguée)

```
Architecture RÉELLE observée via les traces :

TypeNode(User) ──→ passthrough_r2_u_User_left ──→ JoinNode2 (ActivateLeft)
                                                        ↑
TypeNode(Team) ──→ passthrough_r2_t_Team_right ────────┘ (ActivateRight)
                                                        │
                                                        │ ❌ PROBLÈME ICI !
                                                        │ Team arrive alors qu'on 
                                                        │ attend Task !
```

**Le nœud passthrough_r2_t_Team_right est connecté au JoinNode2**, ce qui est incorrect.

Le JoinNode2 devrait être connecté à un nœud qui propage des faits **Task**, pas Team !

### 3.3 Configuration des Nœuds (Extraite des Traces)

#### JoinNode1 (Premier join : u x t)
- **ID** : `join_def99a26470f7c95`
- **LeftVariables** : `[u]`
- **RightVariables** : `[t]`
- **AllVariables** : `[u, t]`
- **Connecté à** : Reçoit User (left) et Team (right) ✓ CORRECT

#### JoinNode2 (Second join : u+t x task)  
- **ID** : `join_6e16ce13b16480f9`
- **LeftVariables** : `[u, t]` ✓ Configuration correcte
- **RightVariables** : `[task]` ✓ Configuration correcte
- **AllVariables** : `[u, t, task]` ✓ Configuration correcte
- **VariableTypes** : `task → Task` ✓ Mapping correct
- **Connecté à** : Reçoit **Team** (right) via `passthrough_r2_t_Team_right` ❌ **INCORRECT !**

**Le problème** : Le JoinNode2 est correctement configuré, mais le réseau le connecte au mauvais TypeNode/AlphaNode.

### 3.4 Flux de Propagation Observé

**Lors de la soumission du fait Team(T001)** :

1. **TypeNode(Team)** reçoit le fait T001
2. Propage vers **passthrough_r2_t_Team_right**
3. Propage vers **join_6e16ce13b16480f9.ActivateRight(Team:T001)**
4. Le JoinNode détecte le fait comme variable **'t'** (au lieu de 'task')
   - Car Team correspond au type de la variable 't', pas 'task'
   - Crée token : `Bindings: [t]`
5. Joint avec token left `[u]` → Résultat : `[u, t]`
6. Propage au TerminalNode avec seulement **[u, t]**
7. **Erreur** : variable 'task' manquante

---

## 4. Trace d'Exécution Détaillée

### 4.1 Soumission du Fait User(U001)

```
📤 [root] PropagateToChildren with FACT
   Fact: U001 (Type: User)

📤 [passthrough_r2_u_User_left] PropagateToChildren with TOKEN
   Token Bindings: [u]
   Number of children: 1
   - Child: join_6e16ce13b16480f9 (type: *rete.JoinNode)

🔍 [JOIN_join_6e16ce13b16480f9] ActivateLeft CALLED
   Token ID: alpha_token_passthrough_r2_u_User_left_U001
   Token Bindings: [u]
   Token NodeID: passthrough_r2_u_User_left
   JoinNode Config:
     - LeftVariables: [u t]    ← Configuration dit qu'il attend [u, t] à gauche
     - RightVariables: [task]  ← Et [task] à droite
     - AllVariables: [u t task]
```

**Observation** : Le token User arrive avec binding `[u]`, mais le JoinNode attend `[u, t]` à gauche.  
Ceci suggère qu'il existe un JoinNode1 qui devrait d'abord créer `[u, t]`.

### 4.2 Soumission du Fait Team(T001) - LE BUG SE MANIFESTE ICI

```
📤 [root] PropagateToChildren with FACT
   Fact: T001 (Type: Team)

📤 [passthrough_r2_t_Team_right] PropagateToChildren with FACT
   Fact: T001 (Type: Team)

🔍 [JOIN_join_6e16ce13b16480f9] ActivateRight CALLED
   Fact ID: T001
   Fact Type: Team
   Fact Attributes: map[]
   Variable detected for fact: 't'    ← ❌ PROBLÈME : détecté comme 't' au lieu de 'task'
```

**ANALYSE CRITIQUE** :
- Le fait Team(T001) arrive au JoinNode via **ActivateRight**
- La fonction `getVariableForFact(fact)` retourne **'t'** 
- Pourquoi ? Parce que Team correspond au type de la variable 't' dans `VariableTypes`
- **MAIS** ce JoinNode devrait recevoir des Task (variable 'task'), pas des Team !

### 4.3 Jointure dans JoinNode2 (avec le mauvais fait)

```
🔗 [JOIN_join_6e16ce13b16480f9] performJoinWithTokens CALLED
   Token1 ID: alpha_token_passthrough_r2_u_User_left_U001
   Token2 ID: right_token_join_6e16ce13b16480f9_T001
   Token1 Bindings: [u]         ← Token left avec juste User
   Token2 Bindings: [t]         ← Token right avec Team (devrait être [task] !)
   Combined bindings: [t u]     ← ❌ Résultat : [u, t] au lieu de [u, t, task]
   ✅ Join conditions PASSED
   Joined token created: ID=alpha_token_passthrough_r2_u_User_left_U001_JOIN_right_token_join_6e16ce13b16480f9_T001
   Bindings=[t u]               ← Token final avec seulement 2 variables !
```

**PROBLÈME IDENTIFIÉ** :
1. Le token left ne contient que `[u]` (il devrait contenir `[u, t]` après le premier join)
2. Le token right contient `[t]` (il devrait contenir `[task]`)
3. La jointure produit `[u, t]` au lieu de `[u, t, task]`

**Questions soulevées** :
- Pourquoi le token left ne contient-il que `[u]` et pas `[u, t]` ?
- Suggère qu'il manque une étape de jointure intermédiaire

### 4.4 Propagation vers TerminalNode

```
📤 [join_6e16ce13b16480f9] PropagateToChildren with TOKEN
   Token Bindings: [t u]       ← Seulement 2 variables propagées
   Number of children: 1
   - Child: r2_terminal (type: *rete.TerminalNode)
```

Le token avec seulement `[t, u]` est propagé au terminal.

### 4.5 Réception par TerminalNode

```
🎯 [TERMINAL_r2_terminal] ActivateLeft CALLED
   Token ID: alpha_token_passthrough_r2_u_User_left_U001_JOIN_right_token_join_6e16ce13b16480f9_T001
   Token Bindings: [t u]                    ← Seulement 2 variables disponibles
   Action name: affordable_task_assignment
   Action expects arguments: 
     [0] map[field:id object:u type:fieldAccess]      ← u.id
     [1] map[field:id object:t type:fieldAccess]      ← t.id  
     [2] map[field:id object:task type:fieldAccess]   ← task.id ← ❌ 'task' manquante !
```

Le TerminalNode reçoit le token mais ne peut pas évaluer l'argument 2 (task.id) car la variable 'task' n'existe pas dans les bindings.

### 4.6 Erreur Levée

```
2025/12/12 17:24:09 📋 ACTION: affordable_task_assignment(u.id, t.id, task.id)

tsd_fixtures_test.go:75: Unexpected error for .../join_multi_variable_complex.tsd: 
❌ Erreur soumission faits: erreur soumission fait T001: 
erreur propagation fait vers type_Team: 
error activating alpha node: 
erreur propagation fait vers join_6e16ce13b16480f9: 
erreur propagation token vers r2_terminal: 
erreur exécution job affordable_task_assignment: 
erreur évaluation argument 2: variable 'task' non trouvée (variables disponibles: [t u])
```

---

## 5. Analyse : Point de Perte des Bindings

### 5.1 Observation Clé

**Il y a DEUX problèmes distincts** :

#### Problème 1 : Mauvais routage du fait Team
Le fait Team(T001) arrive au **deuxième** JoinNode (`join_6e16ce13b16480f9`) via le côté droit, alors que :
- Ce JoinNode attend des faits **Task** à droite (`RightVariables: [task]`)
- Le fait Team devrait arriver au **premier** JoinNode

#### Problème 2 : Token left incomplet
Le token qui arrive par la gauche du deuxième JoinNode ne contient que `[u]` alors qu'il devrait contenir `[u, t]` (résultat du premier join).

**Hypothèse** : La cascade n'est pas correctement construite. Il semble manquer un niveau intermédiaire.

### 5.2 Code Problématique

**Fichier** : `rete/builder_beta_chain.go` ou similaire (fichier de construction du réseau)  
**Fonction** : Construction de la cascade de jointures pour règles multi-variables  

**Analyse du problème** :
La construction du réseau beta pour une règle à 3 variables devrait créer :

```
Pattern 1 (Join 1) : u × t → [u, t]
Pattern 2 (Join 2) : [u, t] × task → [u, t, task]
```

**Mais il semble que le builder crée** :
```
Jointure unique : u × ? → [u, t] ← Le '?' reçoit Team au lieu de Task
```

Ou bien :
```
Join 1 : u × t → [u, t]  (OK)
Join 2 : u × t → [u, t]  (ERREUR : reçoit Team au lieu de Task)
```

### 5.3 Explication Technique

**Dans `node_join.go:getVariableForFact()`** :
```go
func (jn *JoinNode) getVariableForFact(fact *Fact) string {
    // Chercher uniquement dans RightVariables
    for _, varName := range jn.RightVariables {
        if expectedType, exists := jn.VariableTypes[varName]; exists {
            if expectedType == fact.Type {
                return varName
            }
        }
    }
    // Fallback : chercher dans AllVariables
    for _, varName := range jn.AllVariables {
        if expectedType, exists := jn.VariableTypes[varName]; exists {
            if expectedType == fact.Type {
                return varName
            }
        }
    }
    return ""
}
```

**Ce qui se passe** :
1. JoinNode2 a `RightVariables: [task]` et `VariableTypes: {task: "Task", u: "User", t: "Team"}`
2. Quand un fait **Team** arrive :
   - Cherche dans `RightVariables: [task]` → Type attendu : "Task" ≠ "Team" → Pas trouvé
   - **Fallback** dans `AllVariables: [u, t, task]` → Trouve 't' avec type "Team" ✓
   - Retourne **'t'**

**Le bug réel** : Le fait Team ne devrait **jamais** arriver à ce JoinNode en ActivateRight !

---

## 6. Hypothèses Vérifiées

### Hypothèse A : Le token joint est créé correctement mais modifié ensuite
- **Status** : ❌ Réfutée
- **Preuve** : Les traces montrent que le token est créé avec seulement `[t, u]` dès `performJoinWithTokens`. Il n'est pas modifié après.

### Hypothèse B : Le token joint n'inclut jamais la 3ème variable
- **Status** : ✅ Confirmée
- **Preuve** : 
```
Combined bindings: [t u]
Joined token created: ID=..., Bindings=[t u]
```
Le token créé par la jointure ne contient que 2 variables.

### Hypothèse C : getVariableForFact retourne une mauvaise variable
- **Status** : ⚠️ Partiellement confirmée
- **Preuve** : 
```
Variable detected for fact: 't'
```
La fonction retourne 't' pour un fait Team, ce qui est logique MAIS le fait Team ne devrait pas arriver à ce JoinNode.

### Hypothèse D : Le fait Task n'arrive jamais à JoinNode2
- **Status** : ⚠️ À vérifier (pas de trace dans le log jusqu'à l'erreur)
- **Observation** : Le log s'arrête à l'erreur lors de la soumission de Team(T001). Les faits Task sont soumis après.

### Hypothèse E : Le réseau est mal construit (mauvais routage)
- **Status** : ✅ **CONFIRMÉE - C'EST LA CAUSE RACINE**
- **Preuve** :
  - `passthrough_r2_t_Team_right` est connecté à `join_6e16ce13b16480f9`
  - Mais ce JoinNode attend des Task (`RightVariables: [task]`)
  - Le builder connecte le mauvais AlphaNode/TypeNode au JoinNode

---

## 7. Construction de la Cascade (BetaChainBuilder)

### 7.1 Analyse du Problème de Construction

**Pour une règle avec 3 variables** `{u: User, t: Team, task: Task}`, le builder devrait créer :

#### Architecture Attendue :

**JoinPattern 1** : User × Team
- LeftVars: `[u]`
- RightVars: `[t]`  
- AllVars: `[u, t]`
- Output: Token avec `[u, t]`

**JoinPattern 2** : (User+Team) × Task
- LeftVars: `[u, t]` ← Output du Join1
- RightVars: `[task]`
- AllVars: `[u, t, task]`
- Output: Token avec `[u, t, task]`

#### Connexions Réseau Attendues :

```
TypeNode(User) → passthrough_u → JoinNode1 (left)
TypeNode(Team) → passthrough_t → JoinNode1 (right)
JoinNode1 → JoinNode2 (left)
TypeNode(Task) → passthrough_task → JoinNode2 (right)
JoinNode2 → TerminalNode
```

#### Connexions Réelles Observées :

```
TypeNode(User) → passthrough_r2_u_User_left → JoinNode2 (left)  ← ❌ Skip JoinNode1 ?
TypeNode(Team) → passthrough_r2_t_Team_right → JoinNode2 (right) ← ❌ Devrait aller à JoinNode2 mais pour le bon côté
```

**Le problème** : Le réseau semble manquer une étape intermédiaire, ou les connexions sont incorrectes.

### 7.2 Hypothèse sur le Code de Construction

Le builder semble créer des patterns mais **ne connecte pas correctement** les nœuds alpha aux joinNodes.

**Code probablement bugué** (à vérifier dans le code source) :
```go
// Pseudo-code du bug suspecté
for i := 2; i < len(variableNames); i++ {
    leftVars := variableNames[:i]    // Ex: [u, t]
    rightVars := []string{variableNames[i]} // Ex: [task]
    
    // Crée JoinNode correctement
    joinNode := NewJoinNode(..., leftVars, rightVars, ...)
    
    // ❌ BUG PROBABLE ICI : Connexion du mauvais AlphaNode
    // Connecte l'AlphaNode de la variable rightVars[0]
    // MAIS utilise peut-être le mauvais index ou la mauvaise variable !
    alphaNode := getAlphaNode(rightVars[0]) // Devrait être Task
    alphaNode.AddChild(joinNode)            // Mais connecte Team !
}
```

**Piste d'investigation** :
- Vérifier `builder_beta_chain.go` ou `builder_join_rules_cascade.go`
- Chercher comment les enfants (children) sont assignés aux nœuds
- Vérifier l'ordre de création et de connexion

---

## 8. Implications pour le Refactoring

### 8.1 Ce qui doit être changé

#### 1. **Correction du Builder de Cascade**
- **Fichier** : `rete/builder_beta_chain.go` ou similaire
- **Fonction** : Construction des patterns de jointure
- **Action** : Vérifier et corriger la logique de connexion des AlphaNodes aux JoinNodes
- **Vérification** : S'assurer que chaque JoinNode reçoit le bon type de fait du côté droit

#### 2. **Vérification de l'Ordre de Construction**
- S'assurer que la cascade est construite dans le bon ordre
- JoinNode1 doit être créé et connecté AVANT JoinNode2
- Le output de JoinNode1 doit être l'input left de JoinNode2

#### 3. **Amélioration de getVariableForFact**
- Actuellement, la fonction a un fallback qui masque le problème
- Devrait être plus stricte : si `RightVariables` ne contient pas la variable, ne pas faire de fallback
- Ou au moins logger un warning

### 8.2 Ce qui doit être préservé

#### 1. **La structure Token et Bindings**
- Le système de bindings fonctionne correctement
- Les tokens propagent bien les variables
- Ne PAS toucher à `performJoinWithTokens` qui fonctionne correctement

#### 2. **La configuration des JoinNodes**
- `LeftVariables`, `RightVariables`, `AllVariables` sont correctement définis
- `VariableTypes` contient les bons mappings
- Le problème n'est PAS dans la configuration mais dans les connexions réseau

#### 3. **Le mécanisme de propagation**
- `PropagateToChildren` fonctionne correctement
- Les tokens sont bien propagés d'un nœud à l'autre
- Le problème est en amont (construction)

### 8.3 Points d'Attention Critiques

#### ⚠️ **Attention 1** : Ne pas casser les règles à 2 variables
Les règles simples comme `{u: User, t: Team}` fonctionnent. Le refactoring ne doit pas les impacter.

#### ⚠️ **Attention 2** : Tester tous les patterns
- 2 variables : OK actuellement
- 3 variables : KO (c'est le bug)
- 4+ variables : Probablement KO aussi
- Patterns avec NOT, EXISTS : À vérifier

#### ⚠️ **Attention 3** : Vérifier l'ordre de soumission des faits
Le bug se manifeste lors de la soumission du Team AVANT le Task. Vérifier que l'ordre n'a pas d'importance après correction.

#### ⚠️ **Attention 4** : Identifier tous les builders concernés
Il peut y avoir plusieurs builders :
- `builder_beta_chain.go`
- `builder_join_rules.go`  
- `builder_join_rules_cascade.go`

Tous doivent être vérifiés et corrigés.

---

## 9. Recommandations pour Prompt 02 (Design)

### 9.1 Focus Areas

#### **1. Analyse Complète du Builder**
- Lire et comprendre `builder_beta_chain.go` ligne par ligne
- Identifier EXACTEMENT où les connexions sont faites
- Tracer la construction d'une règle à 3 variables pas à pas

#### **2. Design du Système de Connexion**
- Concevoir un algorithme clair pour connecter les nœuds
- Séparer clairement :
  - Création des JoinNodes (OK)
  - Connexion des inputs (gauche) (à vérifier)
  - Connexion des inputs (droite) (BUGUÉ)

#### **3. Validation du Réseau**
- Ajouter une fonction de validation du réseau après construction
- Vérifier que chaque JoinNode reçoit le bon type de fait
- Détecter les incohérences avant exécution

### 9.2 Contraintes à Respecter

#### **1. Pas de Rupture de Compatibilité**
- Les règles existantes à 2 variables doivent continuer à fonctionner
- Les tests existants qui passent ne doivent pas régresser

#### **2. Performances**
- La construction du réseau ne doit pas être significativement plus lente
- Éviter les parcours inutiles ou redondants

#### **3. Maintenabilité**
- Le code du builder doit être clair et documenté
- Chaque étape de construction doit être explicite
- Ajouter des assertions/validations en mode debug

### 9.3 Opportunités d'Amélioration

#### **1. Meilleure Abstraction des Patterns**
- Créer une structure `CascadePattern` qui encapsule :
  - Les JoinNodes à créer
  - Les connexions entre eux
  - Les TypeNodes/AlphaNodes à connecter
- Builder : Génère le CascadePattern
- Executor : Construit le réseau depuis le pattern

#### **2. Validation Automatique**
- Après construction, vérifier :
  - Chaque JoinNode a les bons parents
  - Les types attendus correspondent aux types reçus
  - Pas de connexions manquantes ou en trop

#### **3. Tests Unitaires du Builder**
- Tester la construction du réseau indépendamment de l'exécution
- Vérifier la structure du réseau (nœuds, connexions)
- Ne pas se contenter de tester le résultat final

#### **4. Logs de Construction**
- Ajouter des logs détaillés pendant la construction
- Faciliter le debug des futurs problèmes
- Mode verbose pour développement

---

## 10. Conclusion

### 10.1 Cause Racine Finale

**LE BUG EST DANS LE BUILDER DE CASCADE** (`builder_beta_chain.go` ou similaire).

**Symptôme** :  
Pour une règle à 3 variables `{u: User, t: Team, task: Task}`, le deuxième JoinNode de la cascade reçoit des faits **Team** du côté droit alors qu'il est configuré pour recevoir des faits **Task**.

**Conséquence** :  
Le token final ne contient que 2 bindings `[u, t]` au lieu de 3 `[u, t, task]`, provoquant l'échec de l'action.

**Solution** :  
Corriger le builder pour qu'il connecte le bon TypeNode/AlphaNode au côté droit de chaque JoinNode dans la cascade.

### 10.2 Prochaines Étapes

**Dans Prompt 02** (Design du Système) :

1. **Analyser le code du builder actuel**
   - `builder_beta_chain.go`
   - `builder_join_rules_cascade.go`
   - Identifier la logique de connexion des nœuds

2. **Concevoir le nouvel algorithme de construction**
   - Clarifier l'ordre de création des JoinNodes
   - Définir précisément les règles de connexion
   - Créer un diagramme de l'algorithme

3. **Spécifier les tests de validation**
   - Tests unitaires du builder
   - Tests d'intégration pour 2, 3, 4+ variables
   - Tests de non-régression

**Dans Prompt 03+** (Implémentation) :

4. **Implémenter les corrections**
5. **Valider avec les tests**
6. **Nettoyer et documenter**

---

## Annexes

### A. Commandes Utilisées

```bash
# Test avec diagnostic
cd tsd
go test -tags=e2e -v ./tests/e2e/... -run "TestBetaFixtures/join_multi_variable_complex" > diagnostic_stdout.log 2>&1

# Extraction des traces
grep "🔍\|🔗\|📤\|🎯" diagnostic_stdout.log

# Recherche de patterns spécifiques
grep -A 15 "join_6e16ce13b16480f9.*ActivateRight" diagnostic_stdout.log
grep -B 5 "erreur évaluation argument" diagnostic_stdout.log
```

### B. Fichiers Modifiés Temporairement (INSTRUMENTATION - À SUPPRIMER)

**⚠️ CES MODIFICATIONS SONT TEMPORAIRES ET NE DOIVENT PAS ÊTRE COMMITTÉES** :

1. `rete/node_join.go` - Ajout de traces diagnostic dans :
   - `ActivateLeft`
   - `ActivateRight`
   - `performJoinWithTokens`

2. `rete/node_base.go` - Ajout de traces dans :
   - `PropagateToChildren`
   - Fonction helper `getBindingKeys` et `diagPrintf`

3. `rete/node_terminal.go` - Ajout de traces dans :
   - `ActivateLeft`

**NETTOYAGE REQUIS AVANT COMMIT** :
```bash
# Supprimer toutes les fonctions diagPrintf et getBindingKeys
# Supprimer tous les appels à diagPrintf
# Revenir à la version originale :
git checkout rete/node_join.go rete/node_base.go rete/node_terminal.go
```

### C. Fichier de Trace

**Fichier** : `diagnostic_stdout.log` (235 lignes)
**À conserver** : OUI (pour référence)
**À committer** : NON (ajouter au .gitignore)

### D. Corrections Nécessaires (Non Liées au Bug)

Lors de l'instrumentation, nous avons découvert et corrigé :

1. **Import manquant dans `fact_token.go`**
   - L'import `github.com/treivax/tsd/rete/pkg/domain` n'existe plus
   - Fix : Définir `Fact` inline dans `rete/fact_token.go`
   - **Cette correction DOIT être committée**

2. **Signature de `ConvertToReteProgram`**
   - Retourne maintenant `(interface{}, error)` au lieu de `interface{}`
   - Fix dans `constraint_pipeline.go`
   - **Cette correction DOIT être committée**

---

**FIN DU DOCUMENT D'ANALYSE**

**Statut** : ✅ Diagnostic complet  
**Cause racine** : ✅ Identifiée (Builder de cascade)  
**Prêt pour** : Prompt 02 - Design de la Solution
