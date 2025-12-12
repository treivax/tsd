# Prompt 01 : Diagnostic Approfondi des Jointures Multi-Variables

**Session** : 1/12  
**Durée estimée** : 1-2 heures  
**Pré-requis** : Avoir lu `00_PLAN_ACTION.md` et `RESOLUTION_TESTS_E2E.md`

---

## 🎯 Objectif de cette Session

Identifier **exactement** où et comment les bindings sont perdus dans les jointures à 3+ variables en :
1. Traçant le flux complet d'un fait à travers la cascade
2. Identifiant le point précis de perte des bindings
3. Documentant l'architecture actuelle pour guider le refactoring

**Livrable final** : `tsd/docs/architecture/BINDINGS_ANALYSIS.md` (500-1000 lignes)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Comprendre le Cas d'Étude (20 min)

#### 1.1 Analyser la fixture de test

**Fichier** : `tsd/tests/fixtures/beta/join_multi_variable_complex.tsd`

**Actions** :
1. Lire le fichier complet
2. Identifier la règle r2 qui échoue (3 variables : User, Team, Task)
3. Noter les faits soumis et leurs valeurs
4. Identifier quel fait arrive en dernier (probablement celui qui déclenche l'action)

**Questions à répondre** :
- Combien de faits sont soumis ?
- Quelles sont les conditions de jointure entre les 3 variables ?
- Quelle action devrait être déclenchée ?
- Quel est le message d'erreur exact observé ?

**Documenter dans** : Section "1. Cas d'Étude" de BINDINGS_ANALYSIS.md

---

#### 1.2 Comprendre l'architecture attendue

**Diagramme à produire** (ASCII art ou Mermaid) :

```
Pour la règle : {u: User, t: Team, task: Task}

Architecture attendue :
┌─────────────┐
│ TypeNode    │
│   (User)    │
└──────┬──────┘
       │
       ├─────────────┐
       │             │
┌──────▼──────┐  ┌──▼──────────┐
│ AlphaNode   │  │ TypeNode    │
│   (User)    │  │   (Team)    │
└──────┬──────┘  └──────┬──────┘
       │                │
       └────────┬───────┘
                │
         ┌──────▼──────────┐
         │ JoinNode1       │
         │ Left: [u]       │
         │ Right: [t]      │
         │ All: [u, t]     │
         └──────┬──────────┘
                │
                ├─────────────┐
                │             │
         ┌──────▼──────┐  ┌──▼──────────┐
         │ AlphaNode?  │  │ TypeNode    │
         │             │  │   (Task)    │
         └──────┬──────┘  └──────┬──────┘
                │                │
                └────────┬───────┘
                         │
                  ┌──────▼──────────┐
                  │ JoinNode2       │
                  │ Left: [u, t]    │
                  │ Right: [task]   │
                  │ All: [u, t, task]│
                  └──────┬──────────┘
                         │
                  ┌──────▼──────────┐
                  │ TerminalNode    │
                  │ Action: ...     │
                  └─────────────────┘
```

**Questions** :
- Y a-t-il un AlphaNode entre JoinNode1 et JoinNode2 ?
- Combien de JoinNodes au total ?
- Quels sont les IDs des nœuds ?

---

### Tâche 2 : Instrumentation Temporaire (30 min)

#### 2.1 Ajouter du logging détaillé

**⚠️ IMPORTANT** : Ces modifications sont **TEMPORAIRES** et **NE DOIVENT PAS être committées**.

**Fichier** : `tsd/rete/node_join.go`

**Ajouter en haut du fichier** :
```go
import (
    // ... imports existants ...
    "fmt"
    "sort"
)

// Helper pour le diagnostic
func getBindingKeys(bindings map[string]*Fact) []string {
    if bindings == nil {
        return []string{}
    }
    keys := make([]string, 0, len(bindings))
    for k := range bindings {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    return keys
}
```

**Modifier `ActivateLeft`** (ajouter au début de la fonction) :
```go
func (jn *JoinNode) ActivateLeft(token *Token) error {
    fmt.Printf("\n🔍 [JOIN_%s] ActivateLeft CALLED\n", jn.ID)
    fmt.Printf("   Token ID: %s\n", token.ID)
    fmt.Printf("   Token Bindings: %v\n", getBindingKeys(token.Bindings))
    fmt.Printf("   Token NodeID: %s\n", token.NodeID)
    fmt.Printf("   JoinNode Config:\n")
    fmt.Printf("     - LeftVariables: %v\n", jn.LeftVariables)
    fmt.Printf("     - RightVariables: %v\n", jn.RightVariables)
    fmt.Printf("     - AllVariables: %v\n", jn.AllVariables)
    
    // ... code existant ...
```

**Modifier `ActivateRight`** (ajouter au début) :
```go
func (jn *JoinNode) ActivateRight(fact *Fact) error {
    fmt.Printf("\n🔍 [JOIN_%s] ActivateRight CALLED\n", jn.ID)
    fmt.Printf("   Fact ID: %s\n", fact.ID)
    fmt.Printf("   Fact Type: %s\n", fact.Type)
    fmt.Printf("   Fact Attributes: %v\n", fact.Attributes)
    
    variable := jn.getVariableForFact(fact)
    fmt.Printf("   Variable detected for fact: '%s'\n", variable)
    if variable == "" {
        fmt.Printf("   ⚠️  WARNING: No variable found for fact type %s\n", fact.Type)
        fmt.Printf("   RightVariables: %v\n", jn.RightVariables)
        fmt.Printf("   VariableTypes: %v\n", jn.VariableTypes)
    }
    
    // ... code existant ...
```

**Modifier `performJoinWithTokens`** (ajouter logging détaillé) :
```go
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
    fmt.Printf("\n🔗 [JOIN_%s] performJoinWithTokens CALLED\n", jn.ID)
    fmt.Printf("   Token1 ID: %s, Bindings: %v\n", token1.ID, getBindingKeys(token1.Bindings))
    fmt.Printf("   Token2 ID: %s, Bindings: %v\n", token2.ID, getBindingKeys(token2.Bindings))
    
    // ... code existant pour créer combinedBindings ...
    
    fmt.Printf("   Combined bindings: %v\n", getBindingKeys(combinedBindings))
    
    // ... code de vérification des conditions ...
    
    if !jn.evaluateJoinConditions(combinedBindings) {
        fmt.Printf("   ❌ Join conditions FAILED\n")
        return nil
    }
    
    fmt.Printf("   ✅ Join conditions PASSED\n")
    
    // ... code de création du token joiné ...
    
    fmt.Printf("   Joined token created: ID=%s, Bindings=%v\n", 
        joinedToken.ID, getBindingKeys(joinedToken.Bindings))
    
    return joinedToken
}
```

**Fichier** : `tsd/rete/node_base.go`

**Modifier `PropagateToChildren`** :
```go
func (bn *BaseNode) PropagateToChildren(fact *Fact, token *Token) error {
    if token != nil {
        fmt.Printf("\n📤 [%s] PropagateToChildren with TOKEN\n", bn.ID)
        fmt.Printf("   Token Bindings: %v\n", getBindingKeys(token.Bindings))
        fmt.Printf("   Number of children: %d\n", len(bn.Children))
        for _, child := range bn.Children {
            fmt.Printf("   - Child: %s (type: %T)\n", child.GetID(), child)
        }
    } else if fact != nil {
        fmt.Printf("\n📤 [%s] PropagateToChildren with FACT\n", bn.ID)
        fmt.Printf("   Fact: %s (Type: %s)\n", fact.ID, fact.Type)
    }
    
    // ... code existant ...
}
```

**Fichier** : `tsd/rete/node_terminal.go`

**Modifier `ActivateLeft`** (au début) :
```go
func (tn *TerminalNode) ActivateLeft(token *Token) error {
    fmt.Printf("\n🎯 [TERMINAL_%s] ActivateLeft CALLED\n", tn.ID)
    fmt.Printf("   Token ID: %s\n", token.ID)
    fmt.Printf("   Token Bindings: %v\n", getBindingKeys(token.Bindings))
    fmt.Printf("   Rule: %s\n", tn.Rule.Name)
    
    // Afficher les variables attendues dans l'action
    if tn.Rule != nil && tn.Rule.Action != nil {
        fmt.Printf("   Action name: %s\n", tn.Rule.Action.Name)
        fmt.Printf("   Action expects arguments: ")
        for i, arg := range tn.Rule.Action.Arguments {
            fmt.Printf("\n     [%d] Type: %v", i, arg)
        }
        fmt.Printf("\n")
    }
    
    // ... code existant ...
}
```

---

#### 2.2 Exécuter le test avec capture de la trace

**Commande** :
```bash
cd tsd
go test -tags=e2e -v ./tests/e2e/... -run "TestBetaFixtures/join_multi_variable_complex" 2>&1 | tee diagnostic_output.log
```

**Résultat attendu** : Un fichier `diagnostic_output.log` avec une trace détaillée.

---

### Tâche 3 : Analyser la Trace (40 min)

#### 3.1 Examiner le fichier de trace

**Ouvrir** : `tsd/diagnostic_output.log`

**Chercher les sections clés** :

1. **Activation du TypeNode(Task)** - dernier fait soumis
2. **Propagation vers JoinNode2** - ActivateRight
3. **Jointure dans JoinNode2** - performJoinWithTokens
4. **Propagation vers TerminalNode** - PropagateToChildren
5. **Réception par TerminalNode** - ActivateLeft
6. **Erreur** - message "variable 'task' non trouvée"

#### 3.2 Identifier le point de perte

**Questions critiques à répondre** :

**Q1** : Dans `performJoinWithTokens` du JoinNode2, est-ce que `combinedBindings` contient bien [u, t, task] ?
- Si OUI → Le problème est après la création du token
- Si NON → Le problème est dans la combinaison des bindings

**Q2** : Est-ce que `joinedToken.Bindings` (juste après création) contient [u, t, task] ?
- Si OUI → Le problème est dans la propagation
- Si NON → Le problème est dans la création du token

**Q3** : Dans `PropagateToChildren` du JoinNode2, est-ce que le token propagé a [u, t, task] ?
- Si OUI → Le problème est entre JoinNode2 et TerminalNode
- Si NON → Le problème est dans JoinNode2

**Q4** : Dans `ActivateLeft` du TerminalNode, combien de bindings le token a-t-il ?
- Si 3 → Le problème est dans l'extraction des bindings pour l'action
- Si 2 → Le problème est dans la propagation vers le terminal

**Documenter** : Section "2. Analyse de la Trace" avec extraits de log

---

#### 3.3 Formuler des hypothèses

**Hypothèse A** : Le token joint est créé correctement mais modifié ensuite
- **Test** : Comparer bindings juste après création vs dans PropagateToChildren
- **Cause possible** : Mutation du map Bindings quelque part

**Hypothèse B** : Le token joint n'inclut jamais la 3ème variable
- **Test** : Vérifier `combinedBindings` dans performJoinWithTokens
- **Cause possible** : token2.Bindings est vide ou ne contient pas 'task'

**Hypothèse C** : getVariableForFact retourne "" pour le fait Task
- **Test** : Vérifier la sortie "Variable detected for fact"
- **Cause possible** : RightVariables ne contient pas 'task'

**Hypothèse D** : Le fait Task n'arrive jamais à JoinNode2
- **Test** : Vérifier si ActivateRight est appelé pour JoinNode2
- **Cause possible** : Problème de routage dans le réseau

**Documenter** : Section "3. Hypothèses et Vérifications"

---

### Tâche 4 : Analyser le Code de Construction (30 min)

#### 4.1 Examiner BetaChainBuilder

**Fichiers à lire** :
- `tsd/rete/builder_join_rules_cascade.go`
- `tsd/rete/builder_beta_chain.go`

**Chercher la fonction** : `buildJoinPatterns` ou similaire

**Questions** :
1. Comment les JoinPatterns sont-ils construits pour [u, t, task] ?
2. Le pattern 1 a-t-il : LeftVars=[u], RightVars=[t], AllVars=[u,t] ?
3. Le pattern 2 a-t-il : LeftVars=[u,t], RightVars=[task], AllVars=[u,t,task] ?
4. Les conditions de jointure sont-elles correctement assignées ?

**Vérifier le code** :
```go
// Chercher cette logique ou équivalent
for i := 2; i < len(variableNames); i++ {
    // Comment AllVars est-il construit ?
    // Est-ce que LeftVars contient TOUTES les variables précédentes ?
}
```

**Documenter** : Section "4. Construction de la Cascade"

---

#### 4.2 Vérifier la configuration des JoinNodes

**Dans la trace, extraire** :

Pour JoinNode1 :
- LeftVariables: ?
- RightVariables: ?
- AllVariables: ?

Pour JoinNode2 :
- LeftVariables: ?
- RightVariables: ?
- AllVariables: ?

**Vérifier** :
- JoinNode2.RightVariables contient-il ["task"] ?
- JoinNode2.AllVariables contient-il ["u", "t", "task"] ?
- JoinNode2.VariableTypes["task"] == "Task" ?

**Documenter** : Section "5. Configuration des JoinNodes"

---

### Tâche 5 : Rédiger le Document d'Analyse (30 min)

#### 5.1 Créer le fichier BINDINGS_ANALYSIS.md

**Chemin** : `tsd/docs/architecture/BINDINGS_ANALYSIS.md`

**Structure obligatoire** :

```markdown
# Analyse des Bindings - Jointures Multi-Variables

**Date** : [DATE]  
**Auteur** : Analyse de diagnostic (Prompt 01)  
**Objectif** : Identifier la cause racine de la perte de bindings dans les cascades de jointures

---

## 1. Résumé Exécutif

### 1.1 Problème Observé
[Description du bug : variables manquantes dans le token final]

### 1.2 Cause Racine Identifiée
[Réponse après analyse : où et pourquoi les bindings sont perdus]

### 1.3 Impact
[Quels tests échouent, quel comportement attendu]

---

## 2. Cas d'Étude : join_multi_variable_complex.tsd

### 2.1 Règle Testée
```tsd
[Copier la règle r2 ici]
```

### 2.2 Faits Soumis
1. User : [attributs]
2. Team : [attributs]
3. Task : [attributs]

### 2.3 Comportement Attendu
[Action attendue avec ses arguments]

### 2.4 Comportement Observé
[Erreur : variable 'task' non trouvée]

---

## 3. Architecture Actuelle

### 3.1 Diagramme de Flux
[Diagramme ASCII montrant TypeNodes → JoinNodes → TerminalNode]

### 3.2 Configuration des Nœuds

#### JoinNode1
- ID: [ID du nœud]
- LeftVariables: [u]
- RightVariables: [t]
- AllVariables: [u, t]

#### JoinNode2
- ID: [ID du nœud]
- LeftVariables: [u, t]
- RightVariables: [task]
- AllVariables: [u, t, task]

### 3.3 Flux de Propagation
[Ordre des activations observées]

---

## 4. Trace d'Exécution Détaillée

### 4.1 Soumission du Fait Task
```
[Extraits de log : TypeNode → AlphaNode → JoinNode2.ActivateRight]
```

### 4.2 Jointure dans JoinNode2
```
[Extraits de performJoinWithTokens : token1, token2, combinedBindings]
```

### 4.3 Propagation vers TerminalNode
```
[Extraits de PropagateToChildren : quels bindings sont propagés ?]
```

### 4.4 Réception par TerminalNode
```
[Extraits de TerminalNode.ActivateLeft : combien de bindings ?]
```

### 4.5 Erreur Levée
```
[Message d'erreur exact avec variables disponibles]
```

---

## 5. Analyse : Point de Perte des Bindings

### 5.1 Observation Clé
[Description précise du moment où les bindings passent de 3 à 2]

### 5.2 Code Problématique
**Fichier** : [nom du fichier]  
**Fonction** : [nom de la fonction]  
**Ligne** : [numéro de ligne approximatif]

```go
[Extrait de code pertinent]
```

### 5.3 Explication
[Pourquoi ce code cause la perte de bindings]

---

## 6. Hypothèses Vérifiées

### Hypothèse A : [Description]
- **Status** : ✅ Confirmée / ❌ Réfutée
- **Preuve** : [Extrait de log ou code]

### Hypothèse B : [Description]
- **Status** : ✅ Confirmée / ❌ Réfutée
- **Preuve** : [Extrait de log ou code]

[etc.]

---

## 7. Construction de la Cascade (BetaChainBuilder)

### 7.1 Patterns Générés
**Pattern 1** : [u, t]
- LeftVars: [u]
- RightVars: [t]
- AllVars: [u, t]

**Pattern 2** : [u, t, task]
- LeftVars: [u, t]
- RightVars: [task]
- AllVars: [u, t, task]

### 7.2 Analyse du Code de Construction
[Est-ce que le builder crée correctement les patterns ?]

---

## 8. Implications pour le Refactoring

### 8.1 Ce qui doit être changé
1. [Point 1]
2. [Point 2]
3. [Point 3]

### 8.2 Ce qui doit être préservé
1. [Comportement à maintenir]
2. [Compatibilité à garder]

### 8.3 Points d'Attention Critiques
- [Attention 1]
- [Attention 2]

---

## 9. Recommandations pour Prompt 02 (Design)

### 9.1 Focus Areas
1. [Zone à concentrer les efforts]
2. [Aspect critique du design]

### 9.2 Contraintes à Respecter
1. [Contrainte 1]
2. [Contrainte 2]

### 9.3 Opportunités d'Amélioration
1. [Amélioration au-delà du bug fix]
2. [Simplification possible]

---

## 10. Conclusion

### 10.1 Cause Racine Finale
[Réponse définitive : pourquoi les bindings sont perdus]

### 10.2 Prochaines Étapes
[Ce qui doit être fait dans Prompt 02]

---

**Annexes** :
- Trace complète : `tsd/diagnostic_output.log` (non committé)
- Code instrumenté : Changements temporaires dans `rete/node_*.go` (à supprimer)
```

---

#### 5.2 Remplir toutes les sections

**Instructions** :
1. Copier le template ci-dessus dans `docs/architecture/BINDINGS_ANALYSIS.md`
2. Remplir chaque section avec les informations collectées
3. Utiliser des extraits de log réels (pas d'invention)
4. Être précis : numéros de ligne, noms de variables, valeurs exactes
5. Ajouter des diagrammes ASCII pour clarifier

**Qualité requise** :
- Chaque affirmation doit être prouvée par un extrait de log ou de code
- Les hypothèses doivent être marquées clairement (confirmée/réfutée)
- La cause racine doit être identifiée sans ambiguïté
- Les recommandations doivent être actionnables

---

### Tâche 6 : Nettoyage et Validation (20 min)

#### 6.1 Supprimer TOUT le code de diagnostic

**Fichiers modifiés temporairement** :
- `tsd/rete/node_join.go`
- `tsd/rete/node_base.go`
- `tsd/rete/node_terminal.go`

**Actions** :
1. **Revenir à la version originale** de ces fichiers (git checkout)
2. Vérifier qu'aucun log de debug ne reste :
   ```bash
   grep -r "🔍\|🔗\|📤\|🎯" tsd/rete/*.go
   # Cette commande ne doit RIEN retourner
   ```

**⚠️ IMPORTANT** : Le fichier `diagnostic_output.log` peut être gardé pour référence mais ne doit PAS être committé.

---

#### 6.2 Vérifier le livrable

**Checklist** :
- [ ] Le fichier `docs/architecture/BINDINGS_ANALYSIS.md` existe
- [ ] Toutes les sections du template sont remplies
- [ ] La cause racine est clairement identifiée
- [ ] Des extraits de log sont présents comme preuves
- [ ] Des recommandations pour Prompt 02 sont listées
- [ ] Le code source est revenu à l'état initial (pas de logs debug)
- [ ] `diagnostic_output.log` n'est pas ajouté au git staging

**Validation finale** :
```bash
cd tsd
git status  # Vérifier que seul BINDINGS_ANALYSIS.md est nouveau
git diff    # Vérifier qu'aucun fichier source n'est modifié
```

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Livrables
- [ ] ✅ Fichier `docs/architecture/BINDINGS_ANALYSIS.md` complet (500-1000 lignes)
- [ ] ✅ Fichier `diagnostic_output.log` présent mais non committé
- [ ] ✅ Code source nettoyé (aucune modification restante)

### Connaissances Acquises
- [ ] Point exact de perte des bindings identifié (fichier + fonction + ligne)
- [ ] Cause racine comprise (pourquoi ça échoue)
- [ ] Architecture actuelle documentée (diagrammes + configurations)
- [ ] Hypothèses testées et vérifiées

### Qualité
- [ ] Analyse factuelle (preuves par logs et code)
- [ ] Pas de suppositions non vérifiées
- [ ] Recommandations claires pour Prompt 02
- [ ] Document structuré et lisible

---

## 📊 Questions Clés - Réponses Attendues

À la fin de cette session, vous devez pouvoir répondre :

1. **OÙ** : Dans quelle fonction les bindings sont-ils perdus ?
   - Réponse attendue : [Nom de fichier].[Nom de fonction], ligne ~[N]

2. **QUAND** : À quelle étape de la propagation ?
   - Réponse attendue : Entre [Nœud A] et [Nœud B]

3. **COMMENT** : Quel mécanisme cause la perte ?
   - Réponse attendue : [Mutation du map / Non-copie / Autre]

4. **POURQUOI** : Quelle est la raison fondamentale ?
   - Réponse attendue : [Design flaw dans la structure Token / Builder / Autre]

5. **SCOPE** : Est-ce que les jointures à 2 variables sont affectées ?
   - Réponse attendue : OUI / NON avec justification

---

## 🎯 Prochaine Étape

Une fois ce diagnostic **terminé et validé**, passer au **Prompt 02 - Design du Système Immuable**.

Le Prompt 02 utilisera les findings de cette analyse pour concevoir la solution architecturale complète.

---

## 💡 Conseils Pratiques

### Pour Gagner du Temps
1. Commencer par exécuter le test et capturer la trace
2. Analyser la trace en parallèle de la lecture du code
3. Utiliser des grep pour trouver rapidement les fonctions clés
4. Faire des hypothèses et les tester une par une

### Pour Éviter les Erreurs
1. Ne pas inventer de données - tout doit venir des logs réels
2. Ne pas supposer - vérifier chaque hypothèse
3. Ne pas committer le code de diagnostic
4. Ne pas passer au Prompt 02 sans avoir identifié la cause racine

### Pour un Bon Document
1. Utiliser des extraits de log courts et pertinents
2. Ajouter des diagrammes pour clarifier
3. Être précis dans les numéros de ligne et noms de fonction
4. Expliquer le "pourquoi" pas juste le "quoi"

---

**Note** : Cette session est **purement analytique**. Aucun code de production ne doit être modifié de façon permanente. Le but est de **COMPRENDRE**, pas encore de **CORRIGER**.