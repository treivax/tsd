# Plan d'Implémentation: Partage de Nœuds avec Décomposition en Chaînes
# Phase 2 - Tous les Opérateurs

## Objectif Global
Implémenter le partage maximal de nœuds RETE en décomposant les expressions complexes en chaînes d'AlphaNodes réutilisables, pour TOUS les types d'opérateurs (logiques, arithmétiques, comparaisons, etc.).

## Durée Estimée
2 semaines (14 jours de développement)

---

## PROMPTS À LANCER SUCCESSIVEMENT

### ✅ Prompt 1: Analyse et Extraction des Conditions (Jour 1)
```
Crée le fichier `tsd/rete/alpha_chain_extractor.go` qui extrait et analyse les conditions de n'importe quelle expression (AND, OR, comparaisons, opérations arithmétiques, etc.).

Implémente les fonctions suivantes:

1. `ExtractConditions(expr interface{}) ([]SimpleCondition, string, error)`
   - Extrait toutes les conditions simples d'une expression complexe
   - Retourne: liste de conditions, type d'opérateur principal (AND/OR/etc.), erreur
   - Gère les expressions imbriquées récursivement

2. `SimpleCondition` struct:
   - Type: string (binaryOperation, comparison, arithmetic, etc.)
   - Left: interface{} (opérande gauche)
   - Operator: string
   - Right: interface{} (opérande droite)
   - Hash: string (calculé automatiquement)

3. `CanonicalString(condition SimpleCondition) string`
   - Génère une représentation textuelle unique et déterministe
   - Format: "type(left,operator,right)"
   - Exemples:
     * p.age > 18 → "binaryOp(fieldAccess(p,age),>,literal(18))"
     * p.salary + 100 → "arithmetic(fieldAccess(p,salary),+,literal(100))"

4. Tests unitaires:
   - TestExtractConditions_SimpleComparison
   - TestExtractConditions_LogicalAND
   - TestExtractConditions_NestedExpressions
   - TestExtractConditions_ArithmeticOperations
   - TestCanonicalString_Deterministic
   - TestCanonicalString_Uniqueness

Critères de succès:
- Tous les tests passent
- Gère correctement les expressions imbriquées
- CanonicalString est déterministe (même condition → même string)
```

---

### ✅ Prompt 2: Normalisation Canonique (Jour 2)
```
Dans `tsd/rete/alpha_chain_extractor.go`, ajoute les fonctions de normalisation qui ordonnent les conditions de manière canonique.

Implémente:

1. `NormalizeConditions(conditions []SimpleCondition, operator string) []SimpleCondition`
   - Trie les conditions dans un ordre canonique déterministe
   - Respecte les règles de commutativité selon l'opérateur:
     * AND: commutatif → trier
     * OR: commutatif → trier
     * Opérations séquentielles: préserver l'ordre
   - Retourne la liste triée

2. `IsCommutative(operator string) bool`
   - Retourne true si l'opérateur est commutatif (AND, OR, +, *, etc.)
   - Retourne false pour les opérateurs non-commutatifs (-, /, séquences, etc.)

3. `NormalizeExpression(expr interface{}) (interface{}, error)`
   - Point d'entrée principal
   - Détecte le type d'expression
   - Applique la normalisation appropriée
   - Gère les cas spéciaux (OR, expressions mixtes, etc.)

4. Tests:
   - TestNormalizeConditions_AND_OrderIndependent
   - TestNormalizeConditions_OR_OrderIndependent
   - TestNormalizeConditions_NonCommutative_PreserveOrder
   - TestNormalizeExpression_ComplexNested
   - TestIsCommutative_AllOperators

Critères de succès:
- `A AND B` et `B AND A` normalisent au même ordre
- Les opérateurs non-commutatifs préservent l'ordre
- Tous les tests passent
```

---

### ✅ Prompt 3: Constructeur de Chaînes d'AlphaNodes (Jours 3-4)
```
Crée `tsd/rete/alpha_chain_builder.go` qui construit des chaînes d'AlphaNodes avec partage automatique.

Implémente:

1. Type `AlphaChain`:
   ```go
   type AlphaChain struct {
       Nodes       []*AlphaNode
       Hashes      []string
       FinalNode   *AlphaNode
       RuleID      string
   }
   ```

2. Type `AlphaChainBuilder`:
   ```go
   type AlphaChainBuilder struct {
       network  *ReteNetwork
       storage  Storage
   }
   ```

3. `NewAlphaChainBuilder(network *ReteNetwork, storage Storage) *AlphaChainBuilder`

4. `(acb *AlphaChainBuilder) BuildChain(
       conditions []SimpleCondition,
       variableName string,
       parentNode Node,
       ruleID string,
   ) (*AlphaChain, error)`
   
   Algorithme:
   - Pour chaque condition dans l'ordre normalisé:
     * Appeler `network.AlphaSharingManager.GetOrCreateAlphaNode()`
     * Si nouveau: connecter au parent, ajouter au réseau
     * Si réutilisé: vérifier connexion, logger le partage
     * Enregistrer dans LifecycleManager avec la règle
     * Le nœud devient parent pour le suivant
   - Retourner la chaîne complète

5. Helper `isAlreadyConnected(parent Node, child Node) bool`

6. Tests:
   - TestBuildChain_SingleCondition
   - TestBuildChain_TwoConditions_New
   - TestBuildChain_TwoConditions_Reuse
   - TestBuildChain_PartialReuse
   - TestBuildChain_CompleteReuse
   - TestBuildChain_MultipleRules_SharedSubchain

Critères de succès:
- Partage automatique des nœuds identiques
- Partage partiel fonctionne correctement
- Logging clair (nouveau vs réutilisé)
- Tous les tests passent
```

---

### ✅ Prompt 4: Détection et Décomposition des Expressions (Jour 5)
```
Crée `tsd/rete/expression_analyzer.go` qui analyse une condition et décide comment la traiter.

Implémente:

1. Type `ExpressionType`:
   ```go
   type ExpressionType int
   const (
       SimpleCondition ExpressionType = iota
       ANDExpression
       ORExpression
       MixedExpression
       ArithmeticChain
   )
   ```

2. `AnalyzeExpression(expr interface{}) (ExpressionType, error)`
   - Identifie le type d'expression
   - Retourne le type approprié

3. `CanDecompose(exprType ExpressionType) bool`
   - Retourne true si l'expression peut être décomposée en chaîne
   - true pour: SimpleCondition, ANDExpression, ArithmeticChain (commutatif)
   - false pour: ORExpression (nécessite traitement spécial), MixedExpression

4. `ShouldNormalize(exprType ExpressionType) bool`
   - Détermine si la normalisation est nécessaire

5. Tests:
   - TestAnalyzeExpression_Simple
   - TestAnalyzeExpression_AND
   - TestAnalyzeExpression_OR
   - TestAnalyzeExpression_Mixed_AND_OR
   - TestCanDecompose_AllTypes
   - TestShouldNormalize_AllTypes

Critères de succès:
- Détection correcte de tous les types d'expressions
- Gestion appropriée des cas edge
- Tous les tests passent
```

---

### ✅ Prompt 5: Intégration dans le Pipeline (Jours 6-7)
```
Modifie `tsd/rete/constraint_pipeline_helpers.go` pour intégrer la décomposition en chaînes.

1. Renomme `createAlphaNodeWithTerminal` en `createSimpleAlphaNodeWithTerminal`

2. Crée la nouvelle fonction `createAlphaNodeWithTerminal` qui:
   - Appelle `AnalyzeExpression(condition)` pour identifier le type
   - Si `CanDecompose()` == true:
     * Appelle `ExtractConditions()` puis `NormalizeConditions()`
     * Construit une chaîne avec `BuildChain()`
     * Attache le TerminalNode à la fin de la chaîne
   - Sinon:
     * Appelle `createSimpleAlphaNodeWithTerminal()` (comportement actuel)

3. Ajoute logging détaillé:
   ```
   - "🔗 Décomposition en chaîne: X conditions détectées"
   - "✨ Nouveau AlphaNode créé: [hash]"
   - "♻️  AlphaNode partagé réutilisé: [hash]"
   - "✅ Chaîne construite: X nœud(s), Y partagé(s)"
   ```

4. Gère les cas spéciaux:
   - Expressions OR: créer un seul AlphaNode normalisé
   - Conditions simples: comportement actuel inchangé
   - Erreurs d'extraction: fallback vers comportement actuel

5. Tests d'intégration:
   - TestPipeline_SimpleCondition_NoChange
   - TestPipeline_AND_CreatesChain
   - TestPipeline_OR_SingleNode
   - TestPipeline_TwoRules_ShareChain
   - TestPipeline_Logging_Correct

Critères de succès:
- Backward compatible (conditions simples fonctionnent comme avant)
- Chaînes créées pour expressions AND
- Logging informatif
- Tous les tests passent
```

---

### ✅ Prompt 6: Gestion du Lifecycle pour les Chaînes (Jours 8-9)
```
Modifie `tsd/rete/network.go` pour gérer la suppression correcte des chaînes d'AlphaNodes.

1. Modifie `removeNodeFromNetwork()`:
   - Détecte si un AlphaNode fait partie d'une chaîne
   - Lors de la suppression:
     * Ne supprime QUE si RefCount == 0
     * Déconnecte des parents (TypeNode ou autre AlphaNode)
     * Supprime du registre AlphaSharingManager
     * Supprime du LifecycleManager

2. Crée `removeAlphaChain(ruleID string) error`:
   - Récupère tous les AlphaNodes de la règle via LifecycleManager
   - Remonte la chaîne en ordre inverse (depuis le terminal)
   - Pour chaque nœud:
     * Décrémenter RefCount
     * Si RefCount == 0: supprimer
     * Si RefCount > 0: arrêter (nœuds parents forcément partagés)
   - Log chaque action

3. Améliore `RemoveRule()`:
   - Utilise `removeAlphaChain()` pour les règles avec chaînes
   - Conserve le comportement actuel pour les règles simples

4. Ajoute des helpers:
   - `isPartOfChain(nodeID string) bool`
   - `getChainParent(alphaNode *AlphaNode) Node`

5. Tests:
   - TestRemoveChain_AllNodesUnique_DeletesAll
   - TestRemoveChain_PartialSharing_DeletesOnlyUnused
   - TestRemoveChain_CompleteSharing_DeletesNone
   - TestRemoveRule_WithChain_CorrectCleanup
   - TestRemoveRule_MultipleChains_IndependentCleanup

Critères de succès:
- Suppression correcte sans orphelins
- Nœuds partagés préservés
- Logging détaillé des suppressions
- Tous les tests passent
```

---

### ✅ Prompt 7: Tests End-to-End - Scénarios Réels (Jours 10-11)
```
Crée `tsd/rete/alpha_chain_integration_test.go` avec des tests complets sur des rulesets réels.

Implémente les tests suivants:

1. `TestAlphaChain_TwoRules_SameConditions_DifferentOrder`
   ```constraint
   rule r1: {p: Person} / p.age > 18 AND p.name='toto' => print('A')
   rule r2: {p: Person} / p.name='toto' AND p.age > 18 => print('B')
   ```
   Vérifie: 2 AlphaNodes partagés, 2 TerminalNodes

2. `TestAlphaChain_PartialSharing_ThreeRules`
   ```constraint
   rule r1: {p: Person} / p.age > 18 => print('A')
   rule r2: {p: Person} / p.age > 18 AND p.name='toto' => print('B')
   rule r3: {p: Person} / p.age > 18 AND p.name='toto' AND p.salary > 1000 => print('C')
   ```
   Vérifie: 3 AlphaNodes, partage partiel correct

3. `TestAlphaChain_FactPropagation_ThroughChain`
   - Soumet un fait qui satisfait toute la chaîne
   - Vérifie que tous les TerminalNodes concernés sont activés
   - Vérifie que chaque condition n'est évaluée qu'UNE fois

4. `TestAlphaChain_RuleRemoval_PreservesShared`
   - Crée 3 règles avec partage
   - Supprime la règle du milieu
   - Vérifie que les nœuds partagés restent

5. `TestAlphaChain_ComplexScenario_FraudDetection`
   ```constraint
   type Transaction : <id: string, amount: number, country: string, risk: number>
   
   rule fraud_low: {t: Transaction} / t.amount > 1000 AND t.country = 'XX' => alert('LOW')
   rule fraud_med: {t: Transaction} / t.amount > 1000 AND t.country = 'XX' AND t.risk > 50 => alert('MED')
   rule fraud_high: {t: Transaction} / t.amount > 1000 AND t.country = 'XX' AND t.risk > 80 => alert('HIGH')
   rule large: {t: Transaction} / t.amount > 1000 => log('LARGE')
   ```
   Vérifie: Partage optimal (amount partagé par 4 règles, etc.)

6. `TestAlphaChain_OR_NotDecomposed`
   ```constraint
   rule r1: {p: Person} / p.age > 18 OR p.status='VIP' => print('A')
   ```
   Vérifie: Un seul AlphaNode (pas de décomposition)

7. `TestAlphaChain_NetworkStats_Accurate`
   - Vérifie que `GetNetworkStats()` reporte correctement:
     * Nombre d'AlphaNodes uniques
     * Nombre de références
     * Ratio de partage

Chaque test doit:
- Créer le fichier .constraint
- Builder le réseau avec ConstraintPipeline
- Vérifier la structure du réseau
- Tester la propagation de faits
- Vérifier les statistiques

Critères de succès:
- Tous les scénarios passent
- Partage vérifié dans chaque cas
- Propagation de faits correcte
```

---

### ✅ Prompt 8: Gestion Spéciale des Opérateurs OR (Jour 12)
```
Améliore la gestion des expressions OR dans les fichiers existants.

1. Dans `expression_analyzer.go`, améliore `AnalyzeExpression()`:
   - Détecte les expressions OR pures
   - Détecte les expressions mixtes (AND + OR)
   - Retourne le type approprié

2. Dans `alpha_chain_extractor.go`, ajoute:
   ```go
   func NormalizeORExpression(expr interface{}) (interface{}, error)
   ```
   - Extrait les termes OR
   - Les trie dans l'ordre canonique
   - Reconstruit l'expression normalisée (sans décomposer)

3. Dans `constraint_pipeline_helpers.go`, améliore le traitement:
   - Si ORExpression: normaliser mais créer un seul AlphaNode
   - Si MixedExpression (AND + OR): 
     * Option A: Créer un seul AlphaNode normalisé
     * Option B: Décomposer par groupes (plus complexe)
   - Choisir Option A pour simplicité

4. Tests:
   - TestOR_SingleNode_NotDecomposed
   - TestOR_Normalization_OrderIndependent
   - TestMixedAND_OR_SingleNode
   - TestOR_FactPropagation_Correct

Critères de succès:
- OR n'est pas décomposé en chaîne
- OR est quand même normalisé pour le partage
- Comportement correct avec faits
- Tous les tests passent
```

---

### ✅ Prompt 9: Optimisation des Performances (Jour 13)
```
Optimise les performances de la décomposition en chaînes.

1. Dans `alpha_sharing.go`, améliore `ConditionHash()`:
   - Cache les hash calculés (map[condition]→hash)
   - Évite les recalculs inutiles

2. Dans `alpha_chain_builder.go`, optimise `BuildChain()`:
   - Cache la détection de connexions existantes
   - Réutilise les résultats de normalisation

3. Ajoute des métriques de performance:
   ```go
   type ChainBuildMetrics struct {
       TotalChainsBuilt      int
       TotalNodesCreated     int
       TotalNodesReused      int
       AverageChainLength    float64
       SharingRatio          float64
   }
   ```

4. Dans `network.go`, ajoute:
   ```go
   func (rn *ReteNetwork) GetChainMetrics() *ChainBuildMetrics
   ```

5. Tests de performance:
   - Benchmark avec 100 règles similaires
   - Benchmark avec 1000 règles variées
   - Comparer avant/après optimisations

6. Tests:
   - TestPerformance_LargeRuleset_100Rules
   - TestPerformance_LargeRuleset_1000Rules
   - TestMetrics_Accurate

Critères de succès:
- Cache fonctionne correctement
- Amélioration de performance mesurable
- Métriques précises
- Benchmarks passent
```

---

### ✅ Prompt 10: Documentation Complète (Jour 14)
```
Crée la documentation complète de la fonctionnalité de chaînes d'AlphaNodes.

1. Crée `tsd/rete/ALPHA_CHAINS_USER_GUIDE.md`:
   - Introduction et bénéfices
   - Comment ça marche (avec diagrammes)
   - Exemples d'utilisation
   - Scénarios de partage
   - Guide de débogage

2. Crée `tsd/rete/ALPHA_CHAINS_TECHNICAL_GUIDE.md`:
   - Architecture détaillée
   - Algorithmes de normalisation et construction
   - Lifecycle management
   - Gestion des cas edge
   - API reference

3. Crée `tsd/rete/ALPHA_CHAINS_EXAMPLES.md`:
   - 10+ exemples concrets avec résultats attendus
   - Visualisation des chaînes créées
   - Métriques de partage

4. Mets à jour `tsd/rete/ALPHA_NODE_SHARING.md`:
   - Ajoute section sur les chaînes
   - Mise à jour des exemples
   - Lien vers les nouveaux documents

5. Ajoute des commentaires de code:
   - Docstrings pour toutes les fonctions publiques
   - Exemples dans les commentaires
   - Diagrammes ASCII dans les fichiers complexes

6. Crée `tsd/rete/ALPHA_CHAINS_MIGRATION.md`:
   - Impact sur le code existant (aucun si tout est backward compatible)
   - Comment activer/désactiver les chaînes (si option)
   - Troubleshooting

Critères de succès:
- Documentation complète et claire
- Exemples exécutables
- Diagrammes visuels
- Guide de migration détaillé
```

---

### ✅ Prompt 11: Tests de Régression Complets (Bonus - si temps)
```
Vérifie que toutes les fonctionnalités existantes fonctionnent toujours correctement.

1. Exécute toute la suite de tests RETE:
   ```bash
   cd tsd/rete && go test -v
   ```

2. Si des tests échouent:
   - Identifier la cause (régression vs test obsolète)
   - Corriger le code ou adapter le test
   - Re-tester jusqu'à 100% de succès

3. Teste les scénarios de la conversation précédente:
   - TypeNode sharing (doit toujours fonctionner)
   - Lifecycle management (doit toujours fonctionner)
   - Removal de règles simples (doit toujours fonctionner)

4. Ajoute des tests de régression spécifiques:
   - TestBackwardCompatibility_SimpleRules
   - TestBackwardCompatibility_ExistingBehavior
   - TestNoRegression_AllPreviousTests

Critères de succès:
- 100% des tests existants passent
- Aucune régression détectée
- Backward compatible confirmé
```

---

## RÉSUMÉ DES LIVRABLES

### Code (Production)
- ✅ `alpha_chain_extractor.go` - Extraction et normalisation
- ✅ `alpha_chain_builder.go` - Construction de chaînes
- ✅ `expression_analyzer.go` - Analyse d'expressions
- ✅ Modifications dans `constraint_pipeline_helpers.go` - Intégration
- ✅ Modifications dans `network.go` - Lifecycle pour chaînes
- ✅ Modifications dans `alpha_sharing.go` - Optimisations

### Tests
- ✅ Tests unitaires pour chaque composant (8+ fichiers de tests)
- ✅ `alpha_chain_integration_test.go` - Tests end-to-end
- ✅ Benchmarks de performance
- ✅ Tests de régression

### Documentation
- ✅ `ALPHA_CHAINS_USER_GUIDE.md` - Guide utilisateur
- ✅ `ALPHA_CHAINS_TECHNICAL_GUIDE.md` - Guide technique
- ✅ `ALPHA_CHAINS_EXAMPLES.md` - Exemples
- ✅ `ALPHA_CHAINS_MIGRATION.md` - Guide de migration
- ✅ Mise à jour des docs existantes

---

## MÉTRIQUES DE SUCCÈS

### Fonctionnalité
- ✅ Décomposition en chaînes fonctionne pour expressions AND
- ✅ Partage partiel et complet fonctionne
- ✅ Normalisation rend l'ordre indépendant
- ✅ Expressions OR gérées correctement
- ✅ Backward compatible avec règles simples

### Qualité
- ✅ 100% des tests unitaires passent
- ✅ 100% des tests d'intégration passent
- ✅ Aucune régression sur tests existants
- ✅ Code coverage > 80%

### Performance
- ✅ Partage mesurable (ratio > 1.0)
- ✅ Pas de dégradation pour règles simples
- ✅ Amélioration pour rulesets avec conditions communes

### Documentation
- ✅ Documentation complète et claire
- ✅ Exemples exécutables
- ✅ Guide de migration disponible

---

## ORDRE D'EXÉCUTION

**Lancer les prompts dans l'ordre numérique (1 → 11)**

Chaque prompt est conçu pour:
- Être autonome et testable
- Produire un résultat vérifiable
- S'appuyer sur les résultats des prompts précédents

**Validation à chaque étape**:
Avant de passer au prompt suivant, vérifier que:
1. Le code compile sans erreur
2. Les tests du prompt passent
3. Aucune régression sur tests existants

---

## NOTES IMPORTANTES

### Gestion des Opérateurs
Le plan couvre TOUS les opérateurs:
- **Logiques**: AND, OR, NOT
- **Comparaisons**: >, <, >=, <=, =, !=
- **Arithmétiques**: +, -, *, / (si commutatifs: décomposables)
- **Chaînes**: LIKE, CONTAINS, MATCHES
- **Listes**: IN, CONTAINS

### Commutativité
Seuls les opérateurs commutatifs sont décomposés:
- ✅ AND (commutatif)
- ✅ + et * (arithmétique commutative)
- ❌ OR (traité spécialement)
- ❌ -, / (non-commutatif, ordre important)

### Extensibilité
L'architecture permet d'ajouter facilement:
- Nouveaux types d'opérateurs
- Nouvelles stratégies de normalisation
- Optimisations supplémentaires

---

## SUPPORT ET CONTACT

En cas de blocage ou question:
1. Consulter la documentation technique
2. Examiner les tests existants pour des exemples
3. Revenir à un prompt précédent si nécessaire

**Bonne implémentation! 🚀**

---

**Créé**: Janvier 2025
**Version**: 1.0
**Statut**: Prêt pour exécution séquentielle