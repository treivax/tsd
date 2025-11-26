# 📖 Expliquer du Code

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu as besoin de comprendre comment fonctionne une partie spécifique du code, son rôle dans l'architecture, et comment l'utiliser ou le modifier.

## Objectif

Obtenir une explication claire, détaillée et pédagogique d'une portion de code, incluant son contexte, son fonctionnement interne, et des exemples d'utilisation.

## Instructions

### 1. Identifier le Code à Expliquer

**Précise** :
- **Fichier** : Chemin complet (ex: `rete/node_join.go`)
- **Fonction/Struct** : Nom exact (ex: `evaluateJoinConditions`)
- **Lignes** : Numéros de lignes si pertinent (ex: L240-L290)
- **Contexte** : Pourquoi tu veux comprendre ce code ?

**Exemple** :
```
Fichier : rete/node_join.go
Fonction : evaluateJoinConditions
Lignes : 240-290
Contexte : Je veux comprendre comment les conditions de jointure sont évaluées
```

### 2. Niveau d'Explication Souhaité

**Choisis le niveau** :
- [ ] **Vue d'ensemble** : Rôle général, sans détails techniques
- [ ] **Détaillé** : Explication ligne par ligne si nécessaire
- [ ] **Architecture** : Comment ça s'intègre dans le système global
- [ ] **Utilisation** : Comment l'utiliser dans mon propre code
- [ ] **Modification** : Comment le modifier/étendre

### 3. Questions Spécifiques (optionnel)

**Pose des questions précises** :
- Comment cette fonction gère-t-elle les erreurs ?
- Pourquoi utilise-t-on une map ici plutôt qu'un slice ?
- Quelle est la complexité algorithmique ?
- Comment ce code s'exécute-t-il en mode concurrent ?

## Format de Réponse Attendu

```
=== EXPLICATION DE CODE ===

📍 Localisation
Fichier : rete/node_join.go
Fonction : evaluateJoinConditions
Lignes : 240-290

🎯 Rôle et Responsabilité

Cette fonction évalue si les conditions de jointure entre plusieurs faits
sont satisfaites dans le réseau RETE. Elle est appelée lors de la propagation
de tokens dans les nœuds de jointure (JoinNode) pour déterminer si deux ou
plusieurs faits peuvent être combinés.

📐 Architecture et Contexte

Dans l'architecture RETE :
- Alpha Network : Filtre les faits individuels
- Beta Network : Joint les faits entre eux ← NOUS SOMMES ICI
- Terminal Nodes : Exécute les actions

Le JoinNode utilise evaluateJoinConditions pour :
1. Vérifier que les variables nécessaires sont disponibles
2. Évaluer la condition complète si toutes les variables sont liées
3. Sinon, évaluer partiellement les conditions disponibles

🔍 Fonctionnement Détaillé

Étape 1 : Affichage des bindings (L242-247)
```go
for varName, fact := range bindings {
    fmt.Printf("    %s -> %s (ID: %s)\n", varName, fact.Type, fact.ID)
}
```
Affiche les variables disponibles pour le debug.
Exemple : u -> User (ID: U1), o -> Order (ID: O1)

Étape 2 : Vérification du nombre de variables (L250-253)
```go
if len(bindings) < 2 {
    return false
}
```
Une jointure nécessite au minimum 2 variables différentes.

Étape 3 : Évaluation de la condition complète (L256-286)
```go
if jn.Condition != nil {
    requiredVars := jn.extractRequiredVariables(jn.Condition)
    allVarsAvailable := true
    for _, reqVar := range requiredVars {
        if _, exists := bindings[reqVar]; !exists {
            allVarsAvailable = false
            break
        }
    }
    
    if allVarsAvailable {
        evaluator := NewAlphaConditionEvaluator()
        // ... évaluation complète
    }
}
```

Cette section :
- Extrait les variables requises de la condition (u, o, p, etc.)
- Vérifie si TOUTES sont disponibles
- Si oui → évalue la condition complète
- Si non → passe à l'évaluation partielle (fallback)

📊 Exemple Concret

Scénario : Jointure 3-way (User ⋈ Order ⋈ Product)
Condition : u.id == o.user_id AND o.product_id == p.id AND u.age >= 18

Cas 1 : Seulement u et o disponibles
```
bindings = {u: User{id:U1, age:25}, o: Order{id:O1, user_id:U1, product_id:P1}}
requiredVars = [u, o, p]
allVarsAvailable = false (p manque)
→ Évaluation partielle : u.id == o.user_id évaluée, reste skippé
→ Résultat : true si U1 == U1
```

Cas 2 : Tous disponibles
```
bindings = {u: User{...}, o: Order{...}, p: Product{id:P1}}
requiredVars = [u, o, p]
allVarsAvailable = true
→ Évaluation complète : u.id == o.user_id AND o.product_id == p.id AND u.age >= 18
→ Résultat : true si toutes les conditions satisfaites
```

💡 Concepts Clés

1. **Propagation Incrémentale**
   Les faits arrivent un par un. La fonction doit gérer les cas où
   certaines variables ne sont pas encore disponibles.

2. **Évaluation Partielle**
   Si p manque, on évalue quand même u.id == o.user_id pour filtrer
   les combinaisons invalides tôt.

3. **Variables Liées**
   Une variable est "liée" quand le fait correspondant est disponible
   dans les bindings.

🔧 Comment l'Utiliser

```go
// Créer un JoinNode avec une condition
joinNode := NewJoinNode(
    "join_1",
    condition,        // Condition à évaluer
    []string{"u"},    // Variables gauche
    []string{"o"},    // Variables droite
    varTypes,         // Mapping variable -> type
    storage,
)

// La fonction est appelée automatiquement lors de la propagation
// Elle retourne true si la jointure est valide, false sinon
```

⚙️ Comment le Modifier

Pour ajouter un nouveau type de condition :

1. Modifier extractRequiredVariables() pour extraire les variables
2. Ajouter le support dans collectVariablesFromExpression()
3. Mettre à jour l'évaluateur de conditions si nécessaire
4. Ajouter des tests pour le nouveau type

Exemple : Support des opérateurs de chaînes
```go
case "stringOp":
    // Extraire les variables utilisées dans l'opération
    if obj, ok := expr["object"].(string); ok {
        vars[obj] = true
    }
```

⚠️ Points d'Attention

1. **Concurrence** : La fonction utilise jn.Condition qui est partagée.
   S'assurer que les accès sont thread-safe.

2. **Performance** : extractRequiredVariables() est appelée à chaque
   évaluation. Envisager de mettre en cache le résultat.

3. **Cas edge** : Que se passe-t-il si bindings est vide ? Actuellement
   retourne false, mais devrait-on log un warning ?

4. **Mémoire** : collectVariablesFromExpression() crée une map à chaque
   appel. Pourrait bénéficier d'un sync.Pool.

🔗 Fonctions Liées

- `extractRequiredVariables()` : Extrait les variables de la condition
- `collectVariablesFromExpression()` : Parcourt l'AST récursivement
- `NewAlphaConditionEvaluator()` : Crée l'évaluateur de conditions
- `performJoinWithTokens()` : Appelante qui utilise cette fonction

📚 Pour Aller Plus Loin

- Lire la documentation sur l'algorithme RETE
- Examiner les tests dans node_join_test.go
- Étudier la structure des conditions dans l'AST
- Comprendre le pattern Visitor pour le parcours d'AST
```

## Niveaux d'Explication

### 🔰 Débutant
- Analogies et métaphores
- Explication du "pourquoi" avant le "comment"
- Exemples concrets et simples
- Éviter le jargon technique

### 🎓 Intermédiaire
- Détails techniques modérés
- Diagrammes et schémas
- Comparaisons avec d'autres approches
- Références aux patterns de conception

### 🚀 Expert
- Détails d'implémentation complets
- Analyse de performance
- Trade-offs et alternatives
- Code assembly si pertinent

## Types d'Explications

### 📝 Explication Conceptuelle
Focus sur le "pourquoi" et le rôle dans l'architecture globale.

### 🔬 Explication Technique
Focus sur le "comment" avec détails d'implémentation.

### 🎯 Explication Pratique
Focus sur l'utilisation concrète et les exemples.

### 🏗️ Explication Architecturale
Focus sur l'intégration dans le système global.

## Commandes Utiles

```bash
# Afficher une fonction spécifique
grep -A 50 "func evaluateJoinConditions" rete/node_join.go

# Trouver toutes les utilisations d'une fonction
grep -r "evaluateJoinConditions" .

# Voir l'historique git d'une fonction
git log -p -S "evaluateJoinConditions" rete/node_join.go

# Générer la documentation
godoc -http=:6060
# Puis ouvrir http://localhost:6060/pkg/github.com/treivax/tsd/rete/
```

## Exemple d'Utilisation

```
Je ne comprends pas comment fonctionne la fonction evaluateJoinConditions
dans rete/node_join.go. Peux-tu m'expliquer en utilisant le prompt 
"explain-code" ?

Je voudrais comprendre :
1. Pourquoi on vérifie les variables disponibles
2. Comment fonctionne l'évaluation partielle
3. Quel est l'impact sur la performance

Niveau : Intermédiaire
```

## Checklist pour Bonne Explication

- [ ] Contexte et rôle clairement définis
- [ ] Fonctionnement détaillé expliqué
- [ ] Exemples concrets fournis
- [ ] Diagrammes si complexe
- [ ] Cas d'usage documentés
- [ ] Points d'attention signalés
- [ ] Références aux ressources complémentaires

## Templates par Type de Code

### Pour une Fonction
```
Nom : functionName
Signature : func(params) (returns)
Rôle : Ce qu'elle fait
Paramètres : Explication de chaque param
Retour : Ce qu'elle retourne
Algorithme : Étapes du traitement
Complexité : O(n), O(log n), etc.
Exemple : Code d'utilisation
```

### Pour une Struct
```
Nom : StructName
Rôle : Ce qu'elle représente
Champs : Explication de chaque champ
Méthodes : Liste des méthodes principales
Utilisation : Comment l'instancier et l'utiliser
Relations : Avec quelles autres structs elle interagit
```

### Pour un Package
```
Package : packagename
Rôle : Responsabilité du package
Exports : Types/fonctions publiques principales
Architecture : Organisation interne
Dépendances : Autres packages utilisés
Usage : Exemples d'utilisation
```

## Ressources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Blog](https://go.dev/blog/)
- [RETE Algorithm](https://en.wikipedia.org/wiki/Rete_algorithm)
- [Architecture du projet](../../docs/)

---

**Rappel** : Une bonne explication permet de comprendre non seulement le "comment" mais aussi le "pourquoi" !