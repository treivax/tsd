# Ingestion Incrémentale du Réseau RETE

## Vue d'ensemble

Ce document décrit l'implémentation de l'ingestion incrémentale pour le réseau RETE, qui permet de construire et d'étendre le réseau de manière progressive en chargeant des fichiers multiples contenant des types, des règles et des faits.

Le mode incrémental supporte également la commande `reset` qui permet de réinitialiser complètement le réseau RETE à tout moment.

## Fonctionnalité Principale : `IngestFile`

### Signature

```go
func (cp *ConstraintPipeline) IngestFile(filename string, network *ReteNetwork, storage Storage) (*ReteNetwork, error)
```

### Description

`IngestFile` est la fonction unique et incrémentale pour étendre le réseau RETE. Elle remplace toutes les anciennes variantes de pipeline (`BuildNetworkFromConstraintFile`, `BuildNetworkFromMultipleFiles`, `BuildNetworkFromIterativeParser`, `BuildNetworkFromConstraintFileWithFacts`).

### Caractéristiques

1. **Création ou Extension** : Si `network == nil`, crée un nouveau réseau. Sinon, étend le réseau existant.

2. **Parsing et Validation** : 
   - Parse le fichier spécifié
   - Détecte la présence d'une commande `reset`
   - Validation sémantique lors de la création initiale ou après un reset
   - En mode incrémental (sans reset), la validation est ignorée car les types peuvent être définis dans des fichiers précédemment ingérés

3. **Gestion du Reset** :
   - Si une commande `reset` est détectée, crée un nouveau réseau vide
   - Supprime complètement l'ancien réseau (types, règles, faits, tokens, actions)
   - Le reset est généralement placé en première instruction d'un fichier

4. **Ajout de Types** : Ajoute les nouveaux types au réseau (évite automatiquement les doublons)

5. **Propagation Rétroactive** :
   - Collecte tous les faits existants dans le réseau avant d'ajouter les nouvelles règles
   - Après l'ajout de nouvelles règles, propage les faits existants vers ces nouvelles règles
   - Utilise `RepropagateExistingFact` pour éviter les erreurs de duplication

6. **Soumission de Nouveaux Faits** : Soumet les nouveaux faits définis dans le fichier

### Étapes du Pipeline

1. **Parsing** : Parse le fichier source
2. **Détection Reset** : Vérifie la présence d'une commande `reset`
3. **Reset (si nécessaire)** : Crée un nouveau réseau vide si reset détecté
4. **Validation** : Validation sémantique (pour réseau initial ou après reset)
5. **Conversion** : Convertit l'AST en programme
6. **Création/Extension** : Crée un nouveau réseau ou étend l'existant
7. **Ajout de Types** : Ajoute les types au réseau
8. **Collection** : Collecte les faits existants (ignoré après reset)
9. **Ajout de Règles** : Ajoute les nouvelles règles
10. **Propagation Rétroactive** : Propage les faits existants vers les nouvelles règles (ignoré après reset)
11. **Soumission** : Soumet les nouveaux faits
12. **Validation** : Validation finale du réseau

## Méthode de Propagation : `RepropagateExistingFact`

### Signature

```go
func (rn *ReteNetwork) RepropagateExistingFact(fact *Fact) error
```

### Description

Propage un fait déjà existant dans le réseau vers les nouveaux nœuds sans le rajouter au `RootNode` ou `TypeNode`. Cette méthode est essentielle pour le mode incrémental.

### Fonctionnement

1. Vérifie que le type du fait existe dans le réseau
2. Crée un token pour le fait
3. Propage directement aux enfants du TypeNode (AlphaNodes)

## Collection des Faits Existants : `collectExistingFacts`

### Sources de Collection

La fonction parcourt tous les nœuds du réseau pour collecter les faits existants :

- **RootNode** : Faits en attente de distribution
- **TypeNodes** : Faits typés
- **AlphaNodes** : Faits filtrés
- **BetaNodes** (JoinNodes, ExistsNodes, AccumulatorNodes) : Faits en mémoire de jointure

### Retour

Retourne une slice de faits uniques (dédupliqués par ID).

## Validation en Mode Incrémental

### Validation Sémantique

- **Réseau initial** : Validation complète du programme
- **Après reset** : Validation complète du programme (nouveau réseau vide)
- **Extension incrémentale** : Validation ignorée car les types peuvent provenir de fichiers précédents

### Validation du Réseau

La validation du réseau a été assouplie pour permettre :
- Réseaux sans TypeNodes (au tout début)
- Réseaux sans TerminalNodes (quand seuls des types sont chargés)

Les nœuds terminaux existants sont toujours validés pour s'assurer qu'ils ont des actions.

## Fonctions de Compatibilité (Deprecated)

Pour maintenir la compatibilité avec le code existant, les anciennes fonctions sont conservées mais marquées comme deprecated :

```go
// Deprecated: Utilisez IngestFile à la place
func (cp *ConstraintPipeline) BuildNetworkFromConstraintFile(constraintFile string, storage Storage) (*ReteNetwork, error)

// Deprecated: Utilisez IngestFile à la place avec plusieurs appels
func (cp *ConstraintPipeline) BuildNetworkFromMultipleFiles(filenames []string, storage Storage) (*ReteNetwork, error)

// Deprecated: Utilisez IngestFile à la place
func (cp *ConstraintPipeline) BuildNetworkFromIterativeParser(parser *constraint.IterativeParser, storage Storage) (*ReteNetwork, error)

// Deprecated: Utilisez IngestFile à la place
func (cp *ConstraintPipeline) BuildNetworkFromConstraintFileWithFacts(constraintFile, factsFile string, storage Storage) (*ReteNetwork, []*Fact, error)
```

Ces fonctions utilisent maintenant `IngestFile` en interne.

## Cas d'Usage

### Cas 1 : Chargement Initial Complet

```go
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()

// Fichier contenant types, règles et faits
network, err := pipeline.IngestFile("complete.tsd", nil, storage)
```

### Cas 2 : Chargement Incrémental (Types → Règles → Faits)

```go
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()

// Étape 1 : Types
network, err := pipeline.IngestFile("types.tsd", nil, storage)

// Étape 2 : Règles (les faits existants seront propagés)
network, err = pipeline.IngestFile("rules.tsd", network, storage)

// Étape 3 : Faits supplémentaires
network, err = pipeline.IngestFile("facts.tsd", network, storage)
```

### Cas 3 : Faits Avant Règles

```go
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()

// Étape 1 : Types et faits
network, err := pipeline.IngestFile("types_and_facts.tsd", nil, storage)

// Étape 2 : Règles (propagation rétroactive automatique)
network, err = pipeline.IngestFile("rules.tsd", network, storage)
```

### Cas 4 : Utilisation de Reset

```go
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()

// Réseau initial
network, err := pipeline.IngestFile("initial.tsd", nil, storage)

// Fichier avec reset - supprime tout et repart de zéro
network, err = pipeline.IngestFile("reset_and_new.tsd", network, storage)
// Le réseau ne contient plus que ce qui est dans reset_and_new.tsd

// Continuer en mode incrémental après reset
network, err = pipeline.IngestFile("more_data.tsd", network, storage)
```

**Exemple de fichier avec reset** :
```
reset

type NewType(id: string, value: number)

action new_action(id: string)

rule new_rule: {n: NewType} / n.value > 0 ==> new_action(n.id)

NewType(id: N001, value: 42)
```

### Cas 5 : Extension de Types

```go
pipeline := rete.NewConstraintPipeline()
storage := rete.NewMemoryStorage()

// Type Person avec règles
network, err := pipeline.IngestFile("person.tsd", nil, storage)

// Ajouter type Company avec nouvelles règles
network, err = pipeline.IngestFile("company.tsd", network, storage)
```

## Tests

Les tests d'intégration se trouvent dans `test/integration/incremental/` :

- `TestIncrementalIngestion_FactsBeforeRules` : Vérifie la propagation rétroactive
- `TestIncrementalIngestion_MultipleRules` : Vérifie l'ajout de règles multiples
- `TestIncrementalIngestion_TypeExtension` : Vérifie l'extension de types
- `TestIncrementalIngestion_Reset` : Vérifie le comportement de la commande reset

## Limitations Connues

1. **Avertissements AlphaNode** : Des avertissements "les nœuds alpha ne reçoivent pas de tokens" peuvent apparaître lors de la propagation rétroactive. Ces avertissements sont bénins - la propagation fonctionne correctement via `PropagateToChildren`.

2. **Validation Sémantique** : La validation sémantique est désactivée en mode incrémental pour permettre des références à des types définis dans d'autres fichiers. Elle est réactivée après un reset car le réseau repart de zéro.

3. **Position du Reset** : La commande `reset` est généralement placée en première instruction d'un fichier. Elle supprime immédiatement l'intégralité du réseau existant avant de traiter les instructions suivantes du fichier.

## Bénéfices

1. **Simplicité** : Une seule API publique au lieu de multiples variantes
2. **Flexibilité** : Supporte tous les scénarios de chargement (complet, incrémental, faits avant règles, reset, etc.)
3. **Propagation Automatique** : Les faits existants sont automatiquement propagés vers les nouvelles règles
4. **Extension Progressive** : Permet de construire le réseau progressivement sans reconstruction complète
5. **Reset Simple** : La commande `reset` permet de repartir de zéro à tout moment en créant un nouveau réseau vide
6. **Compatibilité** : Les anciennes fonctions restent disponibles pour la migration progressive

## Migration

Pour migrer du code existant :

```go
// Ancien code
network, err := pipeline.BuildNetworkFromConstraintFile(file, storage)

// Nouveau code
network, err := pipeline.IngestFile(file, nil, storage)
```

```go
// Ancien code
network, err := pipeline.BuildNetworkFromMultipleFiles(files, storage)

// Nouveau code
var network *rete.ReteNetwork
for _, file := range files {
    network, err = pipeline.IngestFile(file, network, storage)
    if err != nil {
        return nil, err
    }
}
```

## Commande Reset

### Description

La commande `reset` permet de réinitialiser complètement le réseau RETE. Lorsqu'elle est rencontrée lors du parsing d'un fichier, elle provoque :

1. **Suppression totale** du réseau existant (types, règles, faits, tokens, actions)
2. **Création** d'un nouveau réseau vide
3. **Traitement** des instructions suivantes du fichier dans le nouveau réseau

### Utilisation

La commande `reset` doit généralement être placée en première ligne d'un fichier :

```
reset

type Person(id: string, name: string)
action greet(name: string)
rule r1: {p: Person} / p.name != "" ==> greet(p.name)
```

### Comportement

- **Avant reset** : Le réseau peut contenir des types, règles et faits
- **Après reset** : Le réseau est complètement vide, comme si on venait de le créer
- **Mode incrémental** : Après un reset, on peut continuer à ingérer d'autres fichiers de manière incrémentale

### Cas d'usage

- **Rechargement complet** : Remplacer totalement la configuration du système
- **Tests** : Repartir d'un état propre entre différents scénarios
- **Reconfiguration** : Changer radicalement les règles et types du système

## Prochaines Étapes

1. **Optimisation** : Améliorer l'identification des nœuds nouvellement créés pour une propagation plus ciblée
2. **Métriques** : Ajouter des métriques de performance pour la propagation rétroactive
3. **Validation Incrémentale** : Implémenter une validation sémantique incrémentale qui prend en compte les types déjà chargés
4. **Documentation API** : Compléter la documentation GoDoc pour toutes les fonctions publiques
5. **Garbage Collection** : Implémenter la libération de mémoire pour les réseaux supprimés par reset