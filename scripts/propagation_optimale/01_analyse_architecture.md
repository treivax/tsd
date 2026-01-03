# 🔍 Prompt 01 - Analyse Architecture et Conception Delta

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Analyser l'architecture RETE actuelle, identifier les points d'extension pour la propagation delta, et concevoir l'architecture complète du système de propagation incrémentale (RETE-II/TREAT).

**⚠️ IMPORTANT** : Ce prompt est **ANALYSE UNIQUEMENT** - Aucun code ne sera généré.

---

## 📋 Tâches

### Tâche 1 : Analyse de l'Architecture RETE Actuelle

#### 1.1 Examiner la Structure du Réseau

Analyser les fichiers suivants et documenter leur rôle :

```bash
# Structure principale
rete/network.go                    # ReteNetwork struct
rete/root_node.go                  # Point d'entrée réseau
rete/type_node.go                  # Filtrage par type
rete/alpha_node.go                 # Nœuds de condition
rete/beta_node.go                  # Nœuds de jointure
rete/terminal_node.go              # Nœuds terminaux (règles)

# Propagation actuelle
rete/propagation.go                # Mécanisme de propagation
rete/network_manager.go            # Gestion insertion/retract
```

**Questions à répondre** :

1. **Flux de propagation actuel** :
   - Comment un fait est-il inséré dans le réseau ?
   - Quels nœuds sont traversés et dans quel ordre ?
   - Comment fonctionne le Retract actuellement ?

2. **Structure des nœuds** :
   - Quels champs/attributs contient chaque type de nœud ?
   - Comment les nœuds stockent-ils leurs conditions/tests ?
   - Y a-t-il déjà des métadonnées exploitables ?

3. **Gestion des faits** :
   - Format interne d'un fait : `map[string]interface{}` ?
   - Comment l'ID interne est-il géré (`Type~values`) ?
   - Où sont stockés les faits actifs ?

**Livrable** : Document Markdown `REPORTS/analyse_rete_actuel.md` avec schémas ASCII.

#### 1.2 Analyser l'Action Update Actuelle

Examiner :
```bash
rete/action_executor_evaluation.go  # evaluateArgument, evaluateUpdate*
rete/action_executor_facts.go       # Gestion faits (Insert/Retract)
rete/actions/builtin_handlers.go    # Handler Update
```

**Questions** :

1. Comment `Update(variable, { field: value })` est-il traité actuellement ?
2. Quel est le flux exact : Retract → Insert → Propagation ?
3. Où sont les points d'insertion pour intercepter et optimiser ?
4. Y a-t-il déjà une détection de no-op (valeurs inchangées) ?

**Livrable** : Diagramme de séquence `REPORTS/sequence_update_actuel.md`

#### 1.3 Identifier les Métadonnées de Nœuds

Pour chaque type de nœud (Alpha, Beta, Terminal), identifier :

1. **Quels champs du fait sont testés/utilisés** :
   - Alpha : conditions sur champs (`field > 10`, `status == "active"`)
   - Beta : champs de jointure (`order.customer_id == customer.id`)
   - Terminal : tous champs utilisés dans actions

2. **Où cette information est-elle disponible** :
   - AST des conditions ?
   - Metadata stockée dans le nœud ?
   - Faut-il parser/extraire depuis les conditions ?

3. **Format des conditions** :
   - Structure de données utilisée
   - Possibilité d'extraction automatique des champs

**Livrable** : Tableau `REPORTS/metadata_noeuds.md` listant chaque type de nœud et ses métadonnées accessibles.

---

### Tâche 2 : Conception de l'Architecture Delta

#### 2.1 Modèle de Données - Structures Principales

Concevoir (sans implémenter) les structures suivantes :

##### 2.1.1 FieldDelta - Représentation d'un Changement

```go
// Concept : Représente le changement d'un champ spécifique
type FieldDelta struct {
    FieldName  string      // Nom du champ modifié
    OldValue   interface{} // Ancienne valeur
    NewValue   interface{} // Nouvelle valeur
    ValueType  string      // Type de la valeur (pour validation)
}

// Concept : Ensemble des changements pour un Update
type FactDelta struct {
    FactID       string                 // ID interne du fait
    FactType     string                 // Type du fait (ex: "Product")
    Fields       map[string]FieldDelta  // Map field -> delta
    Timestamp    time.Time              // Moment du changement
}
```

**Questions de conception** :

1. Faut-il stocker l'ancienne valeur (pour rollback/audit) ?
2. Comment gérer les champs nested (ex: `address.city`) ?
3. Format optimal pour sérialisation (JSON, msgpack) ?

##### 2.1.2 DependencyIndex - Index Nœuds par Champs

```go
// Concept : Index inversé champ → nœuds sensibles
type DependencyIndex struct {
    // Index par type de fait
    alphaIndex map[string]map[string][]*AlphaNode
    // alphaIndex["Product"]["price"] = [alpha1, alpha2, ...]
    
    betaIndex map[string]map[string][]*BetaNode
    // betaIndex["Order"]["customer_id"] = [beta1, beta2, ...]
    
    terminalIndex map[string]map[string][]*TerminalNode
    
    mutex sync.RWMutex // Thread-safety
}
```

**Questions de conception** :

1. Structure de données optimale : map vs trie vs bloom filter ?
2. Granularité : index par (Type, Field) ou global ?
3. Lazy vs Eager : construire à la demande ou au démarrage ?
4. Invalidation : quand reconstruire l'index (ajout de règles) ?

##### 2.1.3 DeltaPropagator - Moteur de Propagation

```go
// Concept : Gère la propagation sélective des deltas
type DeltaPropagator struct {
    network         *ReteNetwork
    dependencyIndex *DependencyIndex
    config          DeltaPropagationConfig
}

type DeltaPropagationConfig struct {
    EnableDeltaPropagation bool          // Feature flag
    MinFieldsForDelta      int           // Seuil : si < N champs, utiliser delta
    MaxIndexSize           int           // Limite mémoire index
    RebuildIndexInterval   time.Duration // Fréquence rebuild index
}
```

**Questions de conception** :

1. Quand utiliser delta vs Retract+Insert classique ?
   - Seuil : nombre de champs modifiés / total champs ?
   - Toujours utiliser delta sauf cas spéciaux ?
2. Comment gérer les cas limites :
   - Modification de clé primaire (changement d'ID) ?
   - Champs calculés/dérivés ?
   - Cascades de mises à jour ?

#### 2.2 Architecture des Composants

Définir les modules et leurs responsabilités :

```
rete/
├── delta/
│   ├── field_delta.go           # Structures FieldDelta, FactDelta
│   ├── dependency_index.go      # Index nœuds par champs
│   ├── index_builder.go         # Construction de l'index
│   ├── delta_detector.go        # Détection des changements
│   ├── delta_propagator.go      # Moteur de propagation
│   ├── config.go                # Configuration delta
│   └── metrics.go               # Métriques performance
│
├── network.go                   # Ajout : *DeltaPropagator
├── action_executor_facts.go     # Modification : intégration delta
└── propagation.go               # Extension : propagation sélective
```

**Responsabilités** :

| Module | Responsabilité | Dépendances |
|--------|----------------|-------------|
| `field_delta.go` | Modèle de données delta | Aucune |
| `dependency_index.go` | Stockage et requête index | `field_delta.go` |
| `index_builder.go` | Construction index depuis réseau | `dependency_index.go`, nœuds RETE |
| `delta_detector.go` | Comparaison faits, extraction delta | `field_delta.go` |
| `delta_propagator.go` | Orchestration propagation | Tous les modules delta |

#### 2.3 Algorithmes Clés

##### 2.3.1 Construction de l'Index

**Entrée** : ReteNetwork complet  
**Sortie** : DependencyIndex rempli

```
ALGORITHME BuildDependencyIndex(network):
    index = new DependencyIndex()
    
    POUR CHAQUE alphaNode DANS network.AlphaNodes:
        fields = ExtractFieldsFromCondition(alphaNode.condition)
        factType = alphaNode.factType
        POUR CHAQUE field DANS fields:
            index.alphaIndex[factType][field].append(alphaNode)
    
    POUR CHAQUE betaNode DANS network.BetaNodes:
        fields = ExtractFieldsFromJoinCondition(betaNode.joinCondition)
        factType = DetermineFactType(betaNode)
        POUR CHAQUE field DANS fields:
            index.betaIndex[factType][field].append(betaNode)
    
    POUR CHAQUE terminalNode DANS network.TerminalNodes:
        fields = ExtractFieldsFromActions(terminalNode.actions)
        factType = terminalNode.factType
        POUR CHAQUE field DANS fields:
            index.terminalIndex[factType][field].append(terminalNode)
    
    RETOURNER index
```

**Question** : Comment extraire les champs depuis les conditions/AST ?

##### 2.3.2 Détection de Delta

**Entrée** : Fait ancien, Fait nouveau  
**Sortie** : FactDelta

```
ALGORITHME DetectDelta(oldFact, newFact):
    delta = new FactDelta()
    delta.FactID = oldFact.ID
    delta.FactType = oldFact.Type
    
    allFields = UNION(keys(oldFact), keys(newFact))
    
    POUR CHAQUE field DANS allFields:
        oldValue = oldFact[field]
        newValue = newFact[field]
        
        SI oldValue != newValue:
            delta.Fields[field] = FieldDelta{
                FieldName: field,
                OldValue:  oldValue,
                NewValue:  newValue,
            }
    
    RETOURNER delta
```

**Questions** :
- Comparaison deep equality pour objets nested ?
- Tolérance pour floats (epsilon) ?

##### 2.3.3 Propagation Sélective

**Entrée** : FactDelta  
**Sortie** : Nœuds activés

```
ALGORITHME PropagateDelta(delta, index):
    affectedNodes = new Set()
    
    POUR CHAQUE fieldDelta DANS delta.Fields:
        field = fieldDelta.FieldName
        factType = delta.FactType
        
        // Trouver nœuds alpha sensibles
        alphas = index.alphaIndex[factType][field]
        affectedNodes.addAll(alphas)
        
        // Trouver nœuds beta sensibles
        betas = index.betaIndex[factType][field]
        affectedNodes.addAll(betas)
        
        // Trouver terminaux sensibles
        terminals = index.terminalIndex[factType][field]
        affectedNodes.addAll(terminals)
    
    // Propager uniquement vers nœuds affectés
    POUR CHAQUE node DANS affectedNodes:
        PropagateToNode(node, delta)
    
    RETOURNER affectedNodes
```

**Optimisation** : Ordre de traversée (topologique) pour éviter re-propagations.

#### 2.4 Intégration avec Update

Modifier le flux d'exécution de `Update(variable, { field: value })` :

**Flux actuel** :
```
1. Évaluer variable → oldFact
2. Créer newFact = merge(oldFact, modifications)
3. Retract(oldFact)
4. Insert(newFact)
5. Propagation complète
```

**Nouveau flux (delta)** :
```
1. Évaluer variable → oldFact
2. Créer newFact = merge(oldFact, modifications)
3. DetectDelta(oldFact, newFact) → delta
4. SI delta.isEmpty():
       RETURN (no-op, déjà implémenté)
5. SI ShouldUseDelta(delta):
       DeltaPropagator.Propagate(delta)
   SINON:
       Retract(oldFact) + Insert(newFact)  // Fallback
```

**Critères `ShouldUseDelta`** :
- Nombre de champs modifiés < seuil (ex: < 30% des champs)
- Pas de modification de clé primaire
- Index construit et disponible

---

### Tâche 3 : Plan de Migration et Compatibilité

#### 3.1 Stratégie de Déploiement

**Options** :

**Option A : Feature Flag (Recommandé)**
```go
type ReteNetwork struct {
    // ...
    EnableDeltaPropagation bool `json:"-"`
    DeltaPropagator        *DeltaPropagator `json:"-"`
}
```

- Activation opt-in via configuration
- Fallback automatique si problème
- Tests A/B possibles

**Option B : Activation Automatique**
- Delta activé par défaut
- Détection automatique des cas incompatibles
- Plus agressif, risque plus élevé

**Recommandation** : **Option A** pour sécurité et tests progressifs.

#### 3.2 Compatibilité Backward

**Garanties** :

1. **API publique inchangée** :
   - `Update(variable, { field: value })` fonctionne identiquement
   - Résultats sémantiquement identiques
   - Pas de breaking change

2. **Tests existants** :
   - 100% des tests doivent passer avec delta activé
   - Résultats identiques delta ON vs OFF

3. **Migration transparente** :
   - Pas de modification de règles existantes
   - Fichiers TSD compatibles sans changement

#### 3.3 Cas Limites à Gérer

| Cas Limite | Stratégie |
|------------|-----------|
| Modification clé primaire | Fallback Retract+Insert (ID change) |
| Champs calculés/dérivés | Tracking transitive dependencies |
| Update concurrent | Locking ou versioning optimiste |
| Réseau sans index | Lazy build ou désactivation auto |
| Règles très complexes | Seuil de complexité → fallback |

---

### Tâche 4 : Spécifications Techniques Détaillées

#### 4.1 Extraction de Champs depuis Conditions

**Alpha Nodes** - Conditions sur champs :

Exemple condition : `product.price > 100 && product.category == "Electronics"`

Algorithme d'extraction :
```
FONCTION ExtractFieldsFromCondition(condition):
    fields = new Set()
    
    SI condition.type == "binaryOp":
        // Opération binaire (AND, OR, >, <, ==, etc.)
        fields.addAll(ExtractFieldsFromCondition(condition.left))
        fields.addAll(ExtractFieldsFromCondition(condition.right))
    
    SINON SI condition.type == "fieldAccess":
        // Accès direct : variable.field
        fields.add(condition.field)
    
    SINON SI condition.type == "comparison":
        // Comparaison : field op value
        SI condition.left.type == "fieldAccess":
            fields.add(condition.left.field)
        SI condition.right.type == "fieldAccess":
            fields.add(condition.right.field)
    
    RETOURNER fields
```

**Beta Nodes** - Conditions de jointure :

Exemple : `order.customer_id == customer.id`

```
FONCTION ExtractFieldsFromJoinCondition(joinCondition):
    fields = new Map() // factType -> [fields]
    
    POUR CHAQUE test DANS joinCondition.tests:
        SI test.left.type == "fieldAccess":
            factType = test.left.variable  // Ex: "order"
            field = test.left.field         // Ex: "customer_id"
            fields[factType].add(field)
        
        SI test.right.type == "fieldAccess":
            factType = test.right.variable
            field = test.right.field
            fields[factType].add(field)
    
    RETOURNER fields
```

**Terminal Nodes** - Champs dans actions :

Exemple action : `Update(product, { price: product.price * 1.1 })`

```
FONCTION ExtractFieldsFromActions(actions):
    fields = new Set()
    
    POUR CHAQUE action DANS actions:
        SI action.type == "Update":
            // Champs lus dans variable
            fields.add(ALL_FIELDS(action.variable))
            
            // Champs modifiés
            fields.addAll(keys(action.modifications))
        
        SINON SI action.type == "Insert":
            fields.addAll(keys(action.factFields))
        
        SINON SI action.type == "Retract":
            fields.add(ALL_FIELDS(action.variable))
    
    RETOURNER fields
```

#### 4.2 Structure de Condition - Mapping AST

Documenter la structure exacte des conditions dans le code actuel :

**À examiner** :
- `constraint/ast.go` - Structures AST des conditions
- `rete/alpha_node.go` - Comment conditions sont stockées
- `rete/beta_node.go` - Format des join conditions

**Livrable** : Document `REPORTS/ast_conditions_mapping.md` avec :
- Schéma de l'AST des conditions
- Exemples concrets
- Code d'extraction de champs

#### 4.3 Gestion de la Concurrence

**Problématiques** :

1. **Construction index pendant propagation** :
   - Solution : RWMutex (readers = propagation, writer = index rebuild)

2. **Mise à jour index lors ajout de règles** :
   - Solution : Invalidation + rebuild lazy ou incrémental

3. **Updates concurrents sur même fait** :
   - Solution : Sérialisation par factID ou versioning optimiste

**Spécifications** :

```go
type DependencyIndex struct {
    mutex sync.RWMutex
    // ...
}

func (idx *DependencyIndex) GetAffectedAlphaNodes(factType, field string) []*AlphaNode {
    idx.mutex.RLock()
    defer idx.mutex.RUnlock()
    // ...
}

func (idx *DependencyIndex) RebuildIndex(network *ReteNetwork) {
    idx.mutex.Lock()
    defer idx.mutex.Unlock()
    // ...
}
```

#### 4.4 Métriques et Observabilité

**Métriques à collecter** :

```go
type DeltaPropagationMetrics struct {
    // Compteurs
    TotalUpdates           int64
    DeltaUpdates           int64
    FallbackUpdates        int64
    
    // Performance
    AvgFieldsChanged       float64
    AvgNodesAffected       float64
    AvgPropagationTimeMs   float64
    
    // Index
    IndexSize              int64
    IndexRebuildCount      int64
    LastIndexRebuildTime   time.Time
    
    // Gains
    NodeEvaluationsSaved   int64  // Nœuds évités grâce au delta
    EstimatedSpeedupRatio  float64
}
```

**Instrumentation** :
- Logs structurés pour chaque propagation delta
- Prometheus metrics pour monitoring
- Traces pour debugging (OpenTelemetry ?)

---

### Tâche 5 : Plan de Tests

#### 5.1 Stratégie de Test

**Niveaux de tests** :

1. **Tests unitaires** (Prompt 07) :
   - `field_delta_test.go` - Structures delta
   - `dependency_index_test.go` - Index et requêtes
   - `delta_detector_test.go` - Détection changements
   - `delta_propagator_test.go` - Propagation sélective

2. **Tests d'intégration** (Prompt 08) :
   - Scénarios Update complets avec delta
   - Comparaison résultats delta ON vs OFF
   - Tests de régression sur suite existante

3. **Tests de performance** (Prompt 09) :
   - Benchmarks Update delta vs Retract+Insert
   - Stress tests (1000+ règles, 10000+ faits)
   - Profiling CPU/Mémoire

#### 5.2 Cas de Test Critiques

**Tests fonctionnels** :

| Test | Description | Critère de succès |
|------|-------------|-------------------|
| `TestDelta_SingleFieldUpdate` | Update 1 champ sur 10 | Propagation 10x plus rapide |
| `TestDelta_MultiFieldUpdate` | Update 5 champs sur 10 | Gain proportionnel |
| `TestDelta_NoOpDetection` | Update avec valeurs identiques | Aucune propagation |
| `TestDelta_PKModification` | Update clé primaire | Fallback Retract+Insert |
| `TestDelta_ConcurrentUpdates` | Updates parallèles | Résultat cohérent |
| `TestDelta_IndexRebuild` | Ajout règle → rebuild index | Index à jour |
| `TestDelta_BackwardCompat` | Suite tests existante | 100% passing |

**Tests de performance** :

```go
func BenchmarkUpdate_DeltaVsClassic(b *testing.B) {
    // Setup : 100 règles, 1000 faits
    // Update 1 champ sur 20
    // Mesurer : temps, allocations, évaluations
}
```

#### 5.3 Validation de Régression

**Critères** :

1. **Sémantique identique** :
   - Même résultat final avec delta ON/OFF
   - Même ordre d'activation des règles (ou équivalent)

2. **Performance non-dégradée** :
   - Cas où delta non applicable : overhead < 5%
   - Pas de régression sur insertions/retracts classiques

3. **Stabilité** :
   - Tests de charge (24h continue)
   - Pas de fuites mémoire
   - Thread-safety (race detector)

---

## 📊 Livrables du Prompt 01

À la fin de ce prompt, vous devez avoir produit :

### Documents de Conception

1. **`REPORTS/analyse_rete_actuel.md`** :
   - Schéma architecture RETE actuelle
   - Flux de propagation Insert/Retract/Update
   - Structure de chaque type de nœud
   - Points d'extension identifiés

2. **`REPORTS/sequence_update_actuel.md`** :
   - Diagramme de séquence Update actuel
   - Stack trace typique
   - Points d'interception pour delta

3. **`REPORTS/metadata_noeuds.md`** :
   - Tableau : Type nœud → Métadonnées disponibles
   - Accessibilité des champs testés/utilisés
   - Stratégie d'extraction

4. **`REPORTS/ast_conditions_mapping.md`** :
   - Schéma AST des conditions
   - Exemples de parsing
   - Code d'extraction de champs

5. **`REPORTS/conception_delta_architecture.md`** :
   - Architecture complète du système delta
   - Structures de données détaillées
   - Algorithmes clés (pseudocode)
   - Plan de migration
   - Stratégie de tests

### Schémas et Diagrammes

Inclure dans les documents :

1. **Schéma architecture delta** (ASCII art) :
```
┌─────────────────────────────────────────────────────────┐
│                     ReteNetwork                         │
│  ┌───────────────────────────────────────────────────┐  │
│  │            DeltaPropagator                        │  │
│  │  ┌──────────────────┐  ┌──────────────────┐      │  │
│  │  │ DependencyIndex  │  │  DeltaDetector   │      │  │
│  │  └──────────────────┘  └──────────────────┘      │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌──────────┐  │
│  │ Alpha   │  │  Beta   │  │ Terminal│  │ Storage  │  │
│  │ Nodes   │  │ Nodes   │  │  Nodes  │  │          │  │
│  └─────────┘  └─────────┘  └─────────┘  └──────────┘  │
└─────────────────────────────────────────────────────────┘
```

2. **Flux de propagation delta** :
```
Update(product, {price: 150})
    ↓
DetectDelta(oldProduct, newProduct)
    ↓ delta = {Fields: {"price": {100 → 150}}}
    ↓
GetAffectedNodes(delta)
    ↓ index.alphaIndex["Product"]["price"] → [alpha1, alpha3]
    ↓ index.betaIndex["Product"]["price"] → [beta2]
    ↓
PropagateSelective([alpha1, alpha3, beta2])
    ↓
ActivateRules (uniquement règles concernées)
```

3. **Comparaison performance** :
```
Classique:           Delta:
Update               Update
  ↓                    ↓
Retract (100%)       DetectDelta
  ↓                    ↓
Insert (100%)        Propagate (10%)  ← 90% économie
  ↓                    ↓
Propagate (100%)     ActivateRules
```

---

## ✅ Critères de Validation

Avant de passer au Prompt 02, vérifier :

- [ ] **Analyse complète** : Architecture RETE actuelle documentée
- [ ] **Points d'extension identifiés** : Où injecter la logique delta
- [ ] **Structures de données spécifiées** : FieldDelta, DependencyIndex, DeltaPropagator
- [ ] **Algorithmes conçus** : Pseudocode pour index, détection, propagation
- [ ] **Extraction de champs résolue** : Méthode pour extraire champs depuis AST/conditions
- [ ] **Plan de migration défini** : Feature flag, compatibilité, fallback
- [ ] **Stratégie de tests établie** : Cas de test critiques identifiés
- [ ] **Documents livrables produits** : 5 fichiers REPORTS créés

---

## 🚀 Prochaines Étapes

Une fois ce prompt validé :

1. **Révision** : Faire relire la conception par un pair (si possible)
2. **Ajustements** : Itérer sur les points flous ou complexes
3. **Validation** : S'assurer que tout est clair et implémentable
4. **Transition** : Passer au **Prompt 02 - Modèle de Données**

---

## 📚 Références

- `.github/prompts/common.md` - Standards projet
- `.github/prompts/develop.md` - Workflow développement
- Forgy, C. (1982). "Rete: A Fast Algorithm"
- Miranker, D. (1990). "TREAT: A New Match Algorithm"
- Doorenbos, R. (1995). "Production Matching for Large Learning Systems"

---

**Durée estimée** : 2-3 heures  
**Difficulté** : Moyenne (analyse + conception)  
**Prérequis** : Compréhension RETE, lecture code existant  
**Livrable** : 5 documents de conception détaillés