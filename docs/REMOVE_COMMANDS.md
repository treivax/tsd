# Commandes de Suppression (Remove Commands)

## Vue d'ensemble

Le langage TSD supporte deux commandes de suppression dynamiques qui permettent de modifier le réseau RETE et la mémoire de travail pendant l'exécution :

1. **`remove fact`** : Supprime un fait de la mémoire de travail
2. **`remove rule`** : Supprime une règle du réseau RETE

## Commande `remove fact`

### Syntaxe

```tsd
remove fact <TypeName> <FactID>
```

### Description

Supprime un fait spécifique de la mémoire de travail. Le fait est identifié par son type et son identifiant unique.

### Paramètres

- **TypeName** : Le nom du type du fait (ex: `Person`, `Order`)
- **FactID** : L'identifiant unique du fait à supprimer

### Exemples

```tsd
// Définir un type
type Person : <id: string, name: string, age: number>

// Ajouter un fait
Person(id: "P1", name: "Alice", age: 25)

// Supprimer le fait
remove fact Person P1
```

### Cas d'usage

- **Rétractation de faits** : Retirer des informations obsolètes
- **Mise à jour** : Supprimer puis réinsérer un fait modifié
- **Nettoyage** : Libérer la mémoire en supprimant des faits non nécessaires

### Comportement

- Le fait est retiré de la mémoire de travail
- Les tokens associés dans les nœuds du réseau sont propagés en rétractation
- Les activations qui dépendaient de ce fait sont supprimées
- Si le fait n'existe pas, l'opération échoue silencieusement

### Notes importantes

⚠️ **Changement de syntaxe** : Avant cette fonctionnalité, la commande était `remove <TypeName> <FactID>`. La nouvelle syntaxe explicite `remove fact` améliore la clarté et permet l'ajout de la commande `remove rule`.

## Commande `remove rule`

### Syntaxe

```tsd
remove rule <RuleID>
```

### Description

Supprime une règle complète du réseau RETE, incluant tous ses nœuds (alpha, beta, terminal) qui ne sont plus utilisés par d'autres règles.

### Paramètres

- **RuleID** : L'identifiant de la règle à supprimer (défini lors de la déclaration de la règle)

### Exemples

```tsd
// Définir un type
type Person : <id: string, name: string, age: number>

// Définir des règles
rule adult_check : {p: Person} / p.age >= 18 ==> notify(p.id)
rule senior_check : {p: Person} / p.age >= 65 ==> alert(p.id)

// Supprimer une règle
remove rule adult_check
```

### Cas d'usage

- **Désactivation de règles** : Désactiver temporairement certaines règles
- **Optimisation** : Supprimer des règles non utilisées pour améliorer les performances
- **Reconfiguration dynamique** : Adapter le comportement du système à l'exécution
- **Tests** : Isoler l'exécution de certaines règles

### Comportement

Lorsqu'une règle est supprimée :

1. **Terminal Node** : Le nœud terminal de la règle est supprimé
2. **Alpha Nodes** : Les nœuds alpha sont supprimés **uniquement** s'ils ne sont pas partagés
3. **Reference Counting** : Le système utilise un compteur de références pour gérer le partage
4. **Nettoyage progressif** : La suppression remonte la chaîne de nœuds de manière sécurisée

#### Partage de nœuds

Le réseau RETE optimise l'utilisation de la mémoire en **partageant** les nœuds alpha entre règles qui ont des conditions identiques.

**Exemple de partage :**

```tsd
type Person : <id: string, age: number>

// Ces deux règles partagent le nœud alpha pour "p.age >= 18"
rule can_vote : {p: Person} / p.age >= 18 ==> allow_vote(p.id)
rule is_adult : {p: Person} / p.age >= 18 ==> mark_adult(p.id)
```

Si vous supprimez `can_vote` :
- ✅ Son nœud terminal est supprimé
- ✅ Le compteur de références du nœud alpha `p.age >= 18` est décrémenté
- ❌ Le nœud alpha n'est **pas** supprimé (encore utilisé par `is_adult`)

Si vous supprimez ensuite `is_adult` :
- ✅ Son nœud terminal est supprimé
- ✅ Le compteur de références du nœud alpha devient 0
- ✅ Le nœud alpha est maintenant supprimé du réseau

### Gestion des erreurs

Si la règle n'existe pas :
- Un avertissement est loggé
- L'exécution continue (pas d'erreur fatale)
- Le réseau reste dans un état cohérent

### Notes importantes

⚠️ **Suppression irréversible** : Une fois supprimée, la règle ne peut pas être restaurée sans re-parser le fichier source.

⚠️ **Impact sur les activations** : Les activations existantes associées à la règle sont perdues.

## Ordre d'exécution

Les commandes sont exécutées dans l'ordre où elles apparaissent dans le fichier `.tsd` :

```tsd
// 1. Définitions de types
type Person : <id: string, name: string, age: number>

// 2. Définitions de règles
rule r1 : {p: Person} / p.age > 18 ==> action1(p.id)
rule r2 : {p: Person} / p.age > 65 ==> action2(p.id)

// 3. Assertions de faits
Person(id: "P1", name: "Alice", age: 25)
Person(id: "P2", name: "Bob", age: 70)

// 4. Suppressions de faits
remove fact Person P1

// 5. Suppressions de règles
remove rule r2

// 6. Nouveaux faits
Person(id: "P3", name: "Charlie", age: 30)
```

## Exemple complet

```tsd
// Système de gestion des commandes avec suppression dynamique

type Customer : <id: string, name: string, status: string>
type Order : <id: string, customer_id: string, amount: number>

// Règles de traitement
rule vip_discount : {c: Customer} / c.status == "VIP" ==> apply_discount(c.id, 20)
rule regular_discount : {c: Customer} / c.status == "REGULAR" ==> apply_discount(c.id, 5)
rule large_order : {o: Order} / o.amount >= 1000 ==> flag_for_review(o.id)

// Faits initiaux
Customer(id: "C1", name: "Alice", status: "VIP")
Customer(id: "C2", name: "Bob", status: "REGULAR")
Order(id: "O1", customer_id: "C1", amount: 1500)
Order(id: "O2", customer_id: "C2", amount: 500)

// Changement de stratégie : suppression des remises régulières
remove rule regular_discount

// Client révoqué
remove fact Customer C2

// Commande traitée
remove fact Order O2

// Nouveaux clients et commandes
Customer(id: "C3", name: "Charlie", status: "VIP")
Order(id: "O3", customer_id: "C3", amount: 2000)
```

## Avantages

### `remove fact`
- ✅ Gestion dynamique de la mémoire de travail
- ✅ Support des systèmes temps réel
- ✅ Mise à jour incrémentale des connaissances

### `remove rule`
- ✅ Reconfiguration dynamique du comportement
- ✅ Optimisation de la performance (moins de règles = évaluation plus rapide)
- ✅ Isolation pour les tests
- ✅ Gestion intelligente du partage de nœuds

## Limitations

- Les suppressions sont **locales** au fichier `.tsd` en cours d'exécution
- Une règle supprimée ne peut pas être restaurée sans re-parser
- Les suppressions ne sont **pas** transactionnelles (pas de rollback)

## Logs et debugging

Le système affiche des logs détaillés lors des suppressions :

```
🗑️  Suppression de la règle: adult_check
   📊 Nœuds associés à la règle: 2
   ✓ Nœud alpha_21ee82570d6f8f0e marqué pour suppression (plus de références)
   ✓ Nœud adult_check_terminal marqué pour suppression (plus de références)
   🔗 AlphaNode alpha_21ee82570d6f8f0e déconnecté de son parent type_Person
   ✓ AlphaNode alpha_21ee82570d6f8f0e supprimé du AlphaSharingManager
   🗑️  Nœud alpha_21ee82570d6f8f0e supprimé du réseau
   🗑️  Nœud adult_check_terminal supprimé du réseau
✅ Règle adult_check supprimée avec succès (2 nœud(s) supprimé(s))
```

## Migration depuis l'ancienne syntaxe

Si vous utilisez l'ancienne syntaxe `remove <TypeName> <FactID>`, vous devez mettre à jour vers :

```tsd
// ❌ Ancienne syntaxe (ne fonctionne plus)
remove Person P1

// ✅ Nouvelle syntaxe
remove fact Person P1
```

## Références

- [Algorithme RETE](./RETE_ALGORITHM.md)
- [Gestion du cycle de vie des nœuds](./NODE_LIFECYCLE.md)
- [Partage des nœuds Alpha](./ALPHA_SHARING.md)
- [Décomposition en chaînes](./CONSTRAINT_PIPELINE_CHAIN_DECOMPOSITION.md)

## License

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License