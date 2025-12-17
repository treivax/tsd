# Règles sans condition

## Vue d'ensemble

Les règles sans condition permettent de déclencher automatiquement des actions dès qu'un fait correspondant au pattern est asserté dans le système RETE, sans aucune condition de filtrage.

Cette fonctionnalité est particulièrement utile pour :
- Logger automatiquement tous les événements d'un certain type
- Auditer toutes les opérations
- Déclencher des webhooks ou notifications systématiques
- Tracker l'activité du système
- Implémenter des patterns d'événements (Event Sourcing, CQRS)

## Syntaxe

### Forme générale

```tsd
rule nom_regle : {variable: Type} / ==> action(...)
```

La différence avec une règle classique est l'absence de contrainte entre le `/` et le `==>`.

### Comparaison

**Règle classique (avec condition)** :
```tsd
rule adultes : {p: Person} / p.age >= 18 ==> notify(p.email)
```

**Règle sans condition** :
```tsd
rule toutes_personnes : {p: Person} / ==> log(p.name)
```

## Exemples

### Exemple 1 : Logging automatique

```tsd
type Person(#personId: string, name: string, age: number)

action log(message: string)

// Logger toute nouvelle personne
rule log_person : {p: Person} / ==> log(p.name)

// Faits
Person(personId: "1", name: "Alice", age: 30)
Person(personId: "2", name: "Bob", age: 25)
```

**Résultat** : Les deux personnes sont loggées automatiquement.

### Exemple 2 : Audit et tracking

```tsd
type Order(#orderId: string, customerId: string, amount: number)

action audit(action: string, details: string)
action track(eventType: string, entityId: string)

// Auditer toutes les commandes
rule audit_orders : {o: Order} / ==>
    audit("order_created", o.orderId),
    track("order", o.orderId)

// Faits
Order(orderId: "o1", customerId: "c1", amount: 100)
Order(orderId: "o2", customerId: "c2", amount: 200)
```

**Résultat** : Chaque commande est auditée et trackée automatiquement.

### Exemple 3 : Webhook automatique

```tsd
type Event(#eventId: string, type: string, source: string)

action webhook(url: string, payload: string)
action log(message: string)

// Envoyer tous les événements à un webhook
rule notify_webhook : {e: Event} / ==>
    webhook("https://api.example.com/events", e.eventId),
    log(e.type)

// Faits
Event(eventId: "e1", type: "login", source: "web")
Event(eventId: "e2", type: "logout", source: "mobile")
```

### Exemple 4 : Règles mixtes (avec et sans condition)

```tsd
type Product(#productId: string, name: string, price: number)

action log(message: string)
action notify(recipient: string, message: string)

// Règle SANS condition : logger tous les produits
rule log_all_products : {pr: Product} / ==> log(pr.name)

// Règle AVEC condition : notifier uniquement les produits chers
rule alert_expensive : {pr: Product} / pr.price > 1000 ==>
    notify("manager@example.com", pr.name)

// Faits
Product(productId: "p1", name: "Laptop", price: 1200)  // Déclenche les 2 règles
Product(productId: "p2", name: "Mouse", price: 25)     // Déclenche uniquement log_all_products
```

### Exemple 5 : Multi-patterns sans condition

```tsd
type Order(#orderId: string, customerId: string)
type Customer(#customerId: string, name: string)

action match(orderId: string, customerName: string)

// Joindre tous les ordres avec tous les clients (produit cartésien)
rule match_all : {o: Order} / {c: Customer} / ==> match(o.orderId, c.name)
```

**Note** : Sans condition de jointure, cette règle créera une activation pour chaque combinaison ordre-client.

## Sémantique

### Activation

Une règle sans condition s'active pour **chaque** fait qui correspond au pattern, sans aucun filtrage.

### Pattern matching

Le pattern matching fonctionne de la même manière que pour les règles classiques :
- Le type du fait doit correspondre au type déclaré dans le pattern
- Les variables sont liées aux faits correspondants
- Les champs du fait sont accessibles via la variable

### Contraintes implicites

Même sans contrainte explicite, le système RETE :
- Vérifie que le fait existe dans la working memory
- Assure la cohérence des types
- Gère les retraits de faits (les activations sont révoquées si le fait est retiré)

## Cas d'usage

### 1. Event Sourcing

```tsd
type DomainEvent(#eventId: string, aggregateId: string, type: string, timestamp: string)

action persist(eventId: string)
action publish(eventId: string, type: string)

rule store_all_events : {e: DomainEvent} / ==>
    persist(e.eventId),
    publish(e.eventId, e.type)
```

### 2. Monitoring et observabilité

```tsd
type Metric(#metricId: string, name: string, value: number, timestamp: string)

action send_to_prometheus(name: string, value: number)
action log_metric(name: string, value: number)

rule monitor_all : {m: Metric} / ==>
    send_to_prometheus(m.name, m.value),
    log_metric(m.name, m.value)
```

### 3. CQRS - Command handling

```tsd
type Command(#commandId: string, type: string, payload: string)

action validate(commandId: string)
action dispatch(commandId: string, type: string)

rule process_commands : {c: Command} / ==>
    validate(c.commandId),
    dispatch(c.commandId, c.type)
```

### 4. Notification systématique

```tsd
type User(#userId: string, email: string, name: string)

action send_welcome_email(email: string, name: string)
action create_default_preferences(userId: string)

rule onboard_user : {u: User} / ==>
    send_welcome_email(u.email, u.name),
    create_default_preferences(u.userId)
```

## Bonnes pratiques

### ✅ À faire

1. **Utiliser pour les side-effects systématiques**
   ```tsd
   rule audit_all : {e: Event} / ==> audit(e.eventId)
   ```

2. **Combiner avec des règles conditionnelles**
   ```tsd
   rule log_all : {o: Order} / ==> log(o.orderId)
   rule alert_large : {o: Order} / o.amount > 10000 ==> alert(o.orderId)
   ```

3. **Documenter clairement l'intention**
   ```tsd
   // Enregistrer TOUS les événements pour audit de sécurité
   rule security_audit : {e: SecurityEvent} / ==> persist(e.eventId)
   ```

### ❌ À éviter

1. **Éviter les règles multi-patterns sans conditions**
   ```tsd
   // ⚠️ Crée un produit cartésien - peut être très coûteux !
   rule bad_join : {o: Order} / {c: Customer} / ==> match(o.orderId, c.customerId)
   ```

2. **Ne pas créer de boucles infinies**
   ```tsd
   // ⚠️ DANGER : Si l'action crée un nouveau Person, boucle infinie !
   rule dangerous : {p: Person} / ==> create_person(p.name)
   ```

3. **Attention aux performances sur gros volumes**
   ```tsd
   // ⚠️ Si des millions de faits, des millions d'activations !
   rule expensive : {f: Fact} / ==> expensive_operation(f.id)
   ```

## Performance

### Considérations

Les règles sans condition peuvent avoir un impact significatif sur les performances :

1. **Volume d'activations** : Chaque fait déclenche une activation
2. **Pas de filtrage précoce** : Aucune optimisation alpha node possible
3. **Coût des actions** : Chaque action est exécutée pour chaque fait

### Optimisations recommandées

1. **Actions légères** : Préférer des actions simples et rapides
   ```tsd
   // Bon : action légère
   rule log_fast : {e: Event} / ==> log(e.eventId)
   
   // Mauvais : action coûteuse
   rule slow : {e: Event} / ==> complex_computation(e.data)
   ```

2. **Batching** : Grouper les opérations quand possible
   ```tsd
   // Utiliser des aggregations si applicable
   rule batch_events : {e: Event, count: COUNT(e: Event)} / ==>
       process_batch(count)
   ```

3. **Monitoring** : Surveiller le nombre d'activations
   ```tsd
   action monitor_activations(ruleId: string)
   
   rule track : {e: Event} / ==>
       process(e.eventId),
       monitor_activations("track")
   ```

## Architecture interne

### Représentation AST

Dans le JSON AST, les règles sans condition ont `constraints: null` :

```json
{
  "type": "expression",
  "ruleId": "log_all",
  "set": {
    "type": "set",
    "variables": [
      {
        "type": "typedVariable",
        "name": "p",
        "dataType": "Person"
      }
    ]
  },
  "constraints": null,
  "action": {
    "type": "sequenceAction",
    "jobs": [...]
  }
}
```

### Compilation RETE

Les règles sans condition :
1. Créent un **AlphaNode** pour le type matching
2. **Pas de BetaNode** pour les contraintes (il n'y en a pas)
3. Connectent directement au **TerminalNode**
4. S'activent pour chaque fait du type correspondant

### Évaluation

Lors de l'assertion d'un fait :
1. Le fait passe le test du type dans l'AlphaNode
2. Directement propagé au TerminalNode (pas de test de contrainte)
3. Création d'une activation
4. Exécution de l'action

## Limitations

1. **Pas de condition** : Impossible de filtrer au niveau de la règle
   - Solution : Utiliser une règle classique avec condition

2. **Multi-patterns sans conditions** : Produit cartésien
   - Solution : Ajouter au moins une condition de jointure

3. **Performance sur gros volumes** : Beaucoup d'activations
   - Solution : Filtrer en amont ou utiliser des règles conditionnelles

## Migration

### Depuis les règles classiques

Si vous avez une règle qui devrait s'appliquer à tous les faits :

**Avant** :
```tsd
// Condition toujours vraie - mauvaise pratique
rule log_all : {p: Person} / p.age > 0 OR p.age <= 0 ==> log(p.name)
```

**Après** :
```tsd
// Règle sans condition - plus claire
rule log_all : {p: Person} / ==> log(p.name)
```

## Voir aussi

- [Documentation des règles](./rules.md)
- [Architecture RETE](./architecture/rete.md)
- [Actions](./actions.md)
- [Exemples complets](../examples/no_condition_rules.tsd)