# Visualisation Détaillée des Nœuds RETE et Analyse de Partage

**Date**: Décembre 2025  
**Status**: ✅ **IMPLÉMENTÉ**  
**Feature**: Visualisation détaillée des opérations de chaque nœud et analyse du partage

---

## Résumé Exécutif

Ce document décrit l'amélioration apportée au test E2E `TestArithmeticExpressionsE2E` pour afficher en détail :

1. **Les opérations traitées par chaque nœud** (opérateurs, opérandes)
2. **L'analyse précise du partage** entre les règles
3. **La correspondance exacte entre expressions TSD et nœuds RETE**
4. **Un diagramme ASCII détaillé** montrant chaque opération par nœud

### Objectif

Permettre de **visualiser exactement quelle partie d'une expression** correspond à chaque nœud du réseau RETE, et comprendre **quels nœuds sont partagés** entre plusieurs règles, avec une **représentation graphique détaillée** de toutes les opérations arithmétiques et logiques.

---

## Cas de Test : 3 Règles avec Partage

### Fichier TSD : `testdata/arithmetic_e2e.tsd`

```tsd
// Règle 1: Condition (c.qte * 23 - 10 > 0)
rule calcul_facture_base : {p: Produit, c: Commande} /
    c.produit_id == p.id AND c.qte * 23 - 10 > 0
    ==> facture_calculee(...)

// Règle 2: Condition inversée (c.qte * 23 - 10 < 0)
rule calcul_facture_speciale : {p: Produit, c: Commande} /
    c.produit_id == p.id AND c.qte * 23 - 10 < 0
    ==> facture_speciale(...)

// Règle 3: MÊME condition que Règle 1 (c.qte * 23 - 10 > 0)
rule calcul_facture_premium : {p: Produit, c: Commande} /
    c.produit_id == p.id AND c.qte * 23 - 10 > 0
    ==> facture_speciale("Commande premium", ...)
```

### Analyse Attendue

**Nœuds partagés** :
- ✅ **TypeNodes** : `Produit` et `Commande` (partagés par les 3 règles)
- ✅ **PassthroughAlphaNodes** : `passthrough_Produit_left` et `passthrough_Commande_right` (partagés par les 3 règles)

**Nœuds séparés** :
- **JoinNodes** : 3 nœuds distincts (conditions différentes ou pas encore optimisés)
  - `calcul_facture_base_join` : condition `> 0`
  - `calcul_facture_speciale_join` : condition `< 0`
  - `calcul_facture_premium_join` : condition `> 0` (IDENTIQUE à règle 1, mais pas encore partagé)
- **TerminalNodes** : 3 nœuds (un par règle, comportement attendu)

---

## Fonctionnalités Ajoutées

### 1. Formatage Détaillé des Conditions

**Fonction** : `formatCondition(cond interface{}, indent string) string`

Affiche les conditions de manière lisible avec des symboles visuels :

```
🔀 PASSTHROUGH (side: left)
🔢 OPÉRATION BINAIRE: c.qte * 23
⚖️  COMPARAISON: (c.qte * 23 - 10) > 0
🔗 EXPRESSION LOGIQUE:
    ├─ c.produit_id == p.id
    └─ AND
       └─ (c.qte * 23 - 10) > 0
📍 ACCÈS CHAMP: c.qte
🔢 NOMBRE: 23
```

**Types supportés** :
- `passthrough` : Nœuds de routage sans filtrage
- `binaryOperation` : Opérations arithmétiques (`+`, `-`, `*`, `/`)
- `comparison` : Comparaisons (`==`, `>`, `<`, `>=`, `<=`, `!=`)
- `logicalExpr` : Expressions logiques (`AND`, `OR`)
- `fieldAccess` : Accès aux champs d'objets (`p.prix`, `c.qte`)
- `number` : Valeurs numériques constantes

### 2. Visualisation des AlphaNodes avec Détails

**Avant** :
```
AlphaNode: passthrough_Produit_left [passthrough]
   └─ 3 enfant(s)
```

**Après** :
```
🔹 AlphaNode: passthrough_Produit_left [passthrough]
   └─ Type parent: Produit
   └─ Condition:
      🔀 PASSTHROUGH (side: left)
   └─ 3 enfant(s) (JoinNodes)
      └─> calcul_facture_base_join (join)
      └─> calcul_facture_speciale_join (join)
      └─> calcul_facture_premium_join (join)
```

### 3. Visualisation des JoinNodes avec Conditions Complètes

**Avant** :
```
JoinNode: calcul_facture_base_join
   └─ Parent gauche: passthrough_Produit_left
   └─ Parent droite: passthrough_Commande_right
```

**Après** :
```
🔶 JoinNode: calcul_facture_base_join
   └─ Parent gauche: passthrough_Produit_left
   └─ Parent droite: passthrough_Commande_right
   └─ Condition de jointure:
      🔗 EXPRESSION LOGIQUE:
         ├─ ⚖️ COMPARAISON: c.produit_id == p.id
         └─ AND
            └─ ⚖️ COMPARAISON: (c.qte * 23 - 10) > 0
               ├─ Left: 🔢 OPÉRATION BINAIRE: (c.qte * 23) - 10
               │   ├─ Left: 🔢 OPÉRATION BINAIRE: c.qte * 23
               │   └─ Right: 🔢 NOMBRE: 10
               └─ Right: 🔢 NOMBRE: 0
   └─ 1 enfant(s) (TerminalNodes)
      └─> calcul_facture_base_terminal (terminal)
```

### 4. Analyse Détaillée du Partage par Expression

**Nouvelle section** :
```
📊 Analyse détaillée du partage par expression:

   🔍 Conditions de jointure et partage:
      ✓ Condition PARTAGEABLE (utilisée par 2 règles):
         🔗 EXPRESSION LOGIQUE:
            ├─ c.produit_id == p.id
            └─ AND
               └─ (c.qte * 23 - 10) > 0
         Utilisée par:
            - calcul_facture_base_join
            - calcul_facture_premium_join
```

Cette section identifie automatiquement les conditions **identiques** qui pourraient être partagées dans une optimisation future.

### 5. Diagramme ASCII Détaillé avec Opérations

**Nouvelle structure ultra-détaillée** montrant chaque opération :

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              RÉSEAU RETE COMPLET                            │
│                          3 Règles, 2 Types, 6 Faits                         │
└─────────────────────────────────────────────────────────────────────────────┘

                            ┌──────────────────┐
                            │   ROOT NODE      │
                            │  (point entrée)  │
                            └────────┬─────────┘
                                     │
                 ┌───────────────────┼───────────────────┐
                 │                   │                   │
                 ▼                   ▼                   ▼
          ┌────────────┐      ┌────────────┐      ┌────────────┐
          │[T] Produit │ ✅   │[T] Commande│ ✅   │[T] Client  │
          │            │      │            │      │  (unused)  │
          │  p: {...}  │      │  c: {...}  │      └────────────┘
          └──────┬─────┘      └──────┬─────┘
                 │                   │
                 ▼                   ▼
    ┌────────────────────────┐   ┌────────────────────────┐
    │ [α] passthrough_       │ ✅│ [α] passthrough_       │ ✅
    │     Produit_left       │   │     Commande_right     │
    │                        │   │                        │
    │ Opération: ROUTAGE     │   │ Opération: ROUTAGE     │
    │ Side: LEFT             │   │ Side: RIGHT            │
    │ Partagé par: 3 règles  │   │ Partagé par: 3 règles  │
    └───────┬────────────────┘   └────────────┬───────────┘
            │                                 │
            └────────┬───────────┬────────────┘
                     │           │            │
         ┌───────────┘           │            └───────────┐
         │                       │                        │
         ▼                       ▼                        ▼
┌──────────────────────┐ ┌──────────────────────┐ ┌──────────────────────┐
│[β] JoinNode: base    │⚠│[β] JoinNode: special │⚠│[β] JoinNode: premium │⚠
│                      │ │                      │ │                      │
│ CONDITION DE JOINTURE│ │ CONDITION DE JOINTURE│ │ CONDITION DE JOINTURE│
│ ==================== │ │ ==================== │ │ ==================== │
│                      │ │                      │ │                      │
│ 1️⃣  c.produit_id ==   │ │ 1️⃣  c.produit_id ==   │ │ 1️⃣  c.produit_id ==   │
│     p.id             │ │     p.id             │ │     p.id             │
│     (équijointure)   │ │     (équijointure)   │ │     (équijointure)   │
│                      │ │                      │ │                      │
│ 2️⃣  AND              │ │ 2️⃣  AND              │ │ 2️⃣  AND              │
│                      │ │                      │ │                      │
│ 3️⃣  Expression:       │ │ 3️⃣  Expression:       │ │ 3️⃣  Expression:       │
│     c.qte * 23 - 10  │ │     c.qte * 23 - 10  │ │     c.qte * 23 - 10  │
│     > 0              │ │     < 0 (inversée!)  │ │     > 0 (identique!) │
│                      │ │                      │ │                      │
│ DÉCOMPOSITION:       │ │ DÉCOMPOSITION:       │ │ DÉCOMPOSITION:       │
│   🔢 c.qte (accès)   │ │   🔢 c.qte (accès)   │ │   🔢 c.qte (accès)   │
│   ✖️  * 23           │ │   ✖️  * 23           │ │   ✖️  * 23           │
│   ➖ - 10            │ │   ➖ - 10            │ │   ➖ - 10            │
│   ⚖️  > 0            │ │   ⚖️  < 0            │ │   ⚖️  > 0            │
│                      │ │                      │ │                      │
│ Parents:             │ │ Parents:             │ │ Parents:             │
│  - LEFT: p (Produit) │ │  - LEFT: p (Produit) │ │  - LEFT: p (Produit) │
│  - RIGHT: c (Commande│ │  - RIGHT: c (Commande│ │  - RIGHT: c (Commande│
└──────────┬───────────┘ └──────────┬───────────┘ └──────────┬───────────┘
           │                        │                        │
           ▼                        ▼                        ▼
  ┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
  │[⚡] Terminal:    │      │[⚡] Terminal:    │      │[⚡] Terminal:    │
  │    base          │      │    special       │      │    premium       │
  │                  │      │                  │      │                  │
  │ ACTION:          │      │ ACTION:          │      │ ACTION:          │
  │ facture_calculee │      │ facture_speciale │      │ facture_speciale │
  │                  │      │                  │      │                  │
  │ Args (6):        │      │ Args (3):        │      │ Args (3):        │
  │  - c.id          │      │  - c.id          │      │  - c.id          │
  │  - p.prix * c.qte│      │  - "Commande     │      │  - "Premium"    │
  │  - remise calc   │      │    speciale"     │      │  - montant 1.2x  │
  │  - ...           │      │  - montant 1.1x  │      │                  │
  └──────────────────┘      └──────────────────┘      └──────────────────┘
```

**Légende des symboles** :
- `[T]` = TypeNode (routage par type)
- `[α]` = AlphaNode passthrough (routage sans filtrage)
- `[β]` = JoinNode (jointure + évaluation de conditions)
- `[⚡]` = TerminalNode (exécution d'action)
- ✅ = Nœud PARTAGÉ entre plusieurs règles
- ⚠️ = Nœud DÉDIÉ (pourrait être partagé)

**Observations détaillées** :
1. ✅ **TypeNodes PARTAGÉS** : Produit et Commande utilisés par 3 règles
2. ✅ **AlphaNodes PARTAGÉS** : 2 passthrough partagés (chacun connecté à 3 JoinNodes)
3. ⚠️ **JoinNodes SÉPARÉS** : 3 JoinNodes avec conditions **visibles dans le diagramme**
   - `base` et `premium` : même condition `> 0` (partageables !)
   - `special` : condition différente `< 0`
4. ✅ **TerminalNodes DÉDIÉS** : Un par règle avec actions différentes

**Optimisation mesurée** :
- **Avant** : 6 AlphaNodes (2 par règle × 3 règles)
- **Après** : 2 AlphaNodes partagés
- **Réduction** : 67% 🎉

**Innovation** : Chaque JoinNode affiche maintenant :
- La condition complète de jointure
- La décomposition de l'expression arithmétique en étapes
- Les opérateurs utilisés (✖️ multiplication, ➖ soustraction, ⚖️ comparaison)
- Les différences entre conditions similaires

---

## Résultats du Test

### Structure du Réseau Construit

```
✅ Réseau RETE construit avec succès
   - TypeNodes: 3
   - AlphaNodes: 0 (tous passthrough, dans PassthroughRegistry)
   - BetaNodes: 3 (un par règle)
   - TerminalNodes: 3 (un par règle)
   - PassthroughRegistry: 2 ✅ PARTAGÉS!
```

### Analyse du Partage

```
📊 AlphaNodes (partage des filtres et passthrough):
   ✓ PARTAGÉ: passthrough_Produit_left [passthrough] → utilisé par 3 JoinNode(s)
   ✓ PARTAGÉ: passthrough_Commande_right [passthrough] → utilisé par 3 JoinNode(s)

   Résumé AlphaNodes: 2 partagé(s), 0 dédié(s)
   └─ Passthrough: 2 partagé(s), 0 dédié(s)
   └─ Filtres: 0 partagé(s), 0 dédié(s)

   ✅ EXCELLENT: Les nœuds passthrough sont PARTAGÉS entre les règles!
```

### Tokens Générés

```
📈 RÉSUMÉ
✅ Total de tokens générés: 6
✅ Actions exécutées: 6

📊 Tokens par règle:
   - calcul_facture_base: 3 tokens
   - calcul_facture_speciale: 0 tokens (condition < 0 jamais vraie)
   - calcul_facture_premium: 3 tokens

✅ Règle 'calcul_facture_base': 3 tokens
✅ Règle 'calcul_facture_speciale': 0 tokens
✅ Règle 'calcul_facture_premium': 3 tokens (mêmes conditions que règle 1!)
```

---

## Hiérarchie des Nœuds et Correspondance avec les Expressions

### Expression TSD : `c.produit_id == p.id AND c.qte * 23 - 10 > 0`

**Décomposition complète en nœuds RETE** :

```
Expression TSD:
   rule calcul_facture_base : {p: Produit, c: Commande} /
       c.produit_id == p.id AND c.qte * 23 - 10 > 0

Décomposition en nœuds RETE:

   p: Produit
   └─→ TypeNode[type_Produit]              ← Route les faits Produit
       └─→ AlphaNode[passthrough_Produit_left] ✅ PARTAGÉ
           └─→ JoinNode[calcul_facture_base_join]

   c: Commande
   └─→ TypeNode[type_Commande]             ← Route les faits Commande
       └─→ AlphaNode[passthrough_Commande_right] ✅ PARTAGÉ
           └─→ JoinNode[calcul_facture_base_join]

   c.produit_id == p.id
   └─→ Évalué dans JoinNode[calcul_facture_base_join]
       └─→ Comparaison des champs après jointure

   c.qte * 23 - 10 > 0
   └─→ Évalué dans JoinNode[calcul_facture_base_join]
       ├─→ 🔢 Accès: c.qte              (lecture du champ)
       ├─→ ✖️  Multiplication: * 23     (opération binaire)
       ├─→ ➖ Soustraction: - 10        (opération binaire)
       └─→ ⚖️  Comparaison: > 0         (test booléen)

   ==> facture_calculee(...)
   └─→ TerminalNode[calcul_facture_base_terminal]
       └─→ Exécution de l'action avec évaluation des expressions
```

**Visualisation dans le diagramme ASCII** :

Le diagramme ASCII détaillé affiche maintenant **directement dans chaque JoinNode** :
- La condition complète : `c.produit_id == p.id AND c.qte * 23 - 10 > 0`
- La décomposition de l'expression arithmétique en 4 étapes :
  1. `c.qte` (accès au champ)
  2. `* 23` (multiplication)
  3. `- 10` (soustraction)
  4. `> 0` (comparaison)

Cela permet de **voir immédiatement** quelle opération est traitée par quel nœud !

### Nœuds Partagés vs Dédiés

| Expression | Nœud RETE | Type | Partagé? | Raison |
|-----------|-----------|------|----------|--------|
| `p: Produit` | `type_Produit` | TypeNode | ✅ OUI | Type partagé par toutes les règles |
| `c: Commande` | `type_Commande` | TypeNode | ✅ OUI | Type partagé par toutes les règles |
| (routage vers JoinNode) | `passthrough_Produit_left` | AlphaNode | ✅ OUI | Même type, même côté |
| (routage vers JoinNode) | `passthrough_Commande_right` | AlphaNode | ✅ OUI | Même type, même côté |
| `c.produit_id == p.id AND c.qte * 23 - 10 > 0` | `calcul_facture_base_join` | JoinNode | ❌ NON | Conditions identiques mais pas encore optimisé |
| `c.produit_id == p.id AND c.qte * 23 - 10 > 0` | `calcul_facture_premium_join` | JoinNode | ❌ NON | **MÊME** condition mais nœud séparé |
| `c.produit_id == p.id AND c.qte * 23 - 10 < 0` | `calcul_facture_speciale_join` | JoinNode | ❌ NON | Condition **différente** (`<` au lieu de `>`) |
| `facture_calculee(...)` | `calcul_facture_base_terminal` | TerminalNode | ❌ NON | Actions toujours dédiées |

---

## Opportunités d'Optimisation Identifiées

### 1. ✅ Passthrough AlphaNodes : **IMPLÉMENTÉ**

**Avant** :
```
calcul_facture_base_pass_p
calcul_facture_speciale_pass_p
calcul_facture_premium_pass_p
calcul_facture_base_pass_c
calcul_facture_speciale_pass_c
calcul_facture_premium_pass_c
```
**Total** : 6 nœuds

**Après** :
```
passthrough_Produit_left
passthrough_Commande_right
```
**Total** : 2 nœuds (✅ **67% de réduction**)

### 2. ⚠️ JoinNodes : **NON IMPLÉMENTÉ (Future Work)**

**Condition identique détectée** :
- `calcul_facture_base_join` : `c.produit_id == p.id AND c.qte * 23 - 10 > 0`
- `calcul_facture_premium_join` : `c.produit_id == p.id AND c.qte * 23 - 10 > 0` (IDENTIQUE!)

**Optimisation possible** :
- Partager le même JoinNode entre les deux règles
- Chaque JoinNode aurait 2 enfants TerminalNodes
- **Réduction potentielle** : 3 JoinNodes → 2 JoinNodes (33% de réduction)

**Complexité** :
- Nécessite normalisation des conditions AST
- Gestion de la commutativité (`A AND B` ≡ `B AND A`)
- Équivalence d'expressions arithmétiques (`x * 2` ≡ `2 * x`)

---

## Utilisation

### Exécuter le Test

```bash
cd tsd
go test -v ./rete -run TestArithmeticExpressionsE2E
```

### Sortie Attendue

Le test affiche :

1. **Construction du réseau** : parsing et création des nœuds
2. **Visualisation détaillée** :
   - Niveau 1 : TypeNodes avec leurs enfants
   - Niveau 2 : AlphaNodes avec conditions détaillées
   - Niveau 3 : JoinNodes avec conditions de jointure complètes
   - Niveau 4 : TerminalNodes avec actions
3. **Analyse du partage** :
   - Comptage des nœuds partagés vs dédiés
   - Identification des conditions partageables
   - Métriques d'optimisation
4. **Diagramme ASCII** : structure visuelle du réseau
5. **Résultats d'exécution** : tokens générés et actions exécutées

---

## Code Source

### Fichiers Modifiés

1. **`rete/action_arithmetic_e2e_test.go`**
   - Ajout de `formatCondition()` pour affichage détaillé des conditions
   - Ajout de `formatOperand()` pour affichage récursif des opérandes
   - Amélioration de la visualisation des AlphaNodes et JoinNodes
   - Ajout de l'analyse détaillée du partage par expression
   - Mise à jour du diagramme ASCII

2. **`rete/testdata/arithmetic_e2e.tsd`**
   - Ajout d'une 3ème règle (`calcul_facture_premium`) avec condition identique à la règle 1
   - Conditions différentes pour démontrer le partage vs non-partage

### Fonctions Clés

```go
// Formate une condition de manière lisible avec symboles visuels
func formatCondition(cond interface{}, indent string) string

// Formate un opérande (récursif pour expressions imbriquées)
func formatOperand(operand interface{}) string

// Caractères de boîte pour affichage arborescent
func getBoxChar(i, total int) string
func getIndent(i, total int) string
```

---

## Bénéfices

### 1. **Compréhension du Réseau**

- ✅ Voir exactement quelle opération est traitée par quel nœud
- ✅ Comprendre la décomposition des expressions complexes
- ✅ Identifier les nœuds partagés vs dédiés

### 2. **Debugging**

- ✅ Tracer une expression TSD jusqu'aux nœuds RETE
- ✅ Identifier pourquoi une règle ne déclenche pas
- ✅ Vérifier que les conditions sont correctement décomposées

### 3. **Optimisation**

- ✅ Mesurer l'impact du partage (67% de réduction pour passthrough)
- ✅ Identifier les opportunités d'optimisation future (JoinNodes)
- ✅ Visualiser l'évolution avant/après optimisation

### 4. **Documentation**

- ✅ Diagrammes ASCII ultra-détaillés générés automatiquement
- ✅ Analyse détaillée du réseau pour chaque test
- ✅ Métriques de performance (nombre de nœuds, taux de partage)
- ✅ Correspondance exacte expression TSD → nœuds RETE
- ✅ Visualisation des opérations arithmétiques dans les JoinNodes

---

## Limitations et Travail Futur

### Limitations Actuelles

1. **Affichage JSON brut pour conditions complexes** : Les conditions très imbriquées peuvent être difficiles à lire
2. **Pas de graphe visuel interactif** : Sortie texte uniquement
3. **JoinNodes pas encore partagés** : Duplication pour conditions identiques

### Travail Futur

1. **Format HTML interactif** : 
   - Graphe cliquable avec zoom
   - Coloration syntaxique des expressions
   - Filtrage par règle

2. **Partage de JoinNodes** :
   - Normalisation des AST de conditions
   - Détection d'équivalence sémantique
   - Métriques de partage Beta

3. **Analyse de performance en temps réel** :
   - Comptage de propagations par nœud
   - Temps d'évaluation par condition
   - Hotspots du réseau

4. **Export vers formats standards** :
   - GraphViz DOT
   - Mermaid diagrams
   - PlantUML

---

## Conclusion

La visualisation détaillée des nœuds RETE permet de :

✅ **Comprendre précisément** comment les expressions TSD sont traduites en réseau RETE  
✅ **Vérifier le partage** des nœuds entre règles multiples  
✅ **Mesurer l'optimisation** du partage des passthrough (67% de réduction)  
✅ **Identifier les opportunités** d'optimisation future (partage de JoinNodes)  
✅ **Déboguer efficacement** les règles complexes  

**Résultat final** : Le test `TestArithmeticExpressionsE2E` passe avec 3 règles générant 6 tokens, et montre clairement le partage de 2 AlphaNodes passthrough entre les 3 règles.

---

**Auteur** : TSD RETE Engine Team  
**Date** : Décembre 2025  
**Version** : 1.0  
**Status** : ✅ Implémenté et testé