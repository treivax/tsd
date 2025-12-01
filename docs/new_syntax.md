# Nouvelle syntaxe pour les types et actions

## Vue d'ensemble

Ce document décrit la nouvelle syntaxe pour la définition des types et des actions dans TSD. Cette syntaxe a été conçue pour être plus naturelle et intuitive, tout en supportant des fonctionnalités avancées comme les paramètres optionnels et les valeurs par défaut.

## Syntaxe des types

### Ancienne syntaxe (obsolète)

```tsd
type Person : <name: string, age: number, active: bool>
```

### Nouvelle syntaxe (recommandée)

```tsd
type Person(name: string, age: number, active: bool)
```

### Description

La nouvelle syntaxe utilise des **parenthèses** au lieu d'accolades angulaires, ce qui la rend plus proche d'une signature de fonction et plus naturelle pour les développeurs.

**Format:**
```
type <TypeName>(<field1>: <type1>, <field2>: <type2>, ...)
```

**Types primitifs supportés:**
- `string` - Chaîne de caractères
- `number` - Nombre (entier ou décimal)
- `bool` - Booléen (true/false)

### Exemples

```tsd
// Type simple avec un seul champ
type Counter(value: number)

// Type avec plusieurs champs
type Person(name: string, age: number, active: bool)

// Type pour un produit e-commerce
type Product(id: number, title: string, price: number, inStock: bool)

// Type pour une commande
type Order(orderId: number, customerName: string, total: number, paid: bool)

// Type pour un événement système
type Event(eventName: string, timestamp: number, severity: string, handled: bool)
```

## Syntaxe des actions

### Description

Les actions définissent les opérations qui peuvent être exécutées lorsqu'une règle est déclenchée. La nouvelle syntaxe permet de déclarer explicitement la signature des actions avec support pour:

- **Types primitifs** (string, number, bool)
- **Types personnalisés** (définis avec `type`)
- **Paramètres optionnels** (marqués avec `?`)
- **Valeurs par défaut** (avec `= valeur`)

### Format général

```
action <ActionName>(<param1>: <type1>, <param2>: <type2>?, <param3>: <type3> = <default>, ...)
```

### Exemples de base

```tsd
// Action simple avec un argument primitif
action log(message: string)

// Action avec plusieurs arguments primitifs
action notify(recipient: string, message: string, priority: number)

// Action sans arguments
action triggerBatchProcess()
```

### Paramètres optionnels

Un paramètre optionnel est marqué avec `?` après le type. Il peut être omis lors de l'appel de l'action.

```tsd
type User(id: number, name: string, email: string)

// L'argument 'active' est optionnel
action updateUser(user: User, active: bool?)

// Appels valides:
// updateUser(u)          ✓ (active non fourni)
// updateUser(u, true)    ✓ (active fourni)
```

### Valeurs par défaut

Un paramètre peut avoir une valeur par défaut avec `= valeur`. Si l'argument n'est pas fourni, la valeur par défaut est utilisée.

```tsd
// Le paramètre 'priority' a une valeur par défaut de 1
action notify(recipient: string, message: string, priority: number = 1)

// Appels valides:
// notify(user, "Hello")        ✓ (priority = 1 par défaut)
// notify(user, "Alert", 5)     ✓ (priority = 5 explicite)
```

### Types personnalisés dans les actions

Les actions peuvent accepter des types personnalisés définis avec `type`. Le type doit être défini **avant** l'action.

```tsd
// Définition du type
type Person(name: string, age: number, email: string)

// Action acceptant un type personnalisé
action savePerson(person: Person)

// Utilisation dans une règle
rule r1 : {p: Person} / p.age > 18 ==> savePerson(p)
```

### Combinaison de fonctionnalités

```tsd
type Order(id: number, total: number, customerId: string)
type Customer(id: string, name: string, vip: bool)

// Action avec types mixtes, optionnels et valeurs par défaut
action processOrder(
    order: Order,
    discount: number?,
    sendNotification: bool = true,
    priority: number = 1
)

// Action avec plusieurs types personnalisés
action linkOrderToCustomer(order: Order, customer: Customer, verified: bool = false)
```

## Validation lors du parsing

### Validation des définitions d'actions

Les définitions d'actions sont validées lors du parsing:

1. **Types des paramètres**: Tous les types doivent être valides (primitifs ou définis)
2. **Valeurs par défaut**: Les valeurs par défaut doivent correspondre au type du paramètre
3. **Cohérence**: Les types personnalisés référencés doivent exister

**Exemple d'erreur:**

```tsd
type Person(name: string, age: number)

// ❌ ERREUR: Le type 'User' n'est pas défini
action saveUser(user: User)
```

### Validation des appels d'actions

Les appels d'actions dans les règles sont validés:

1. **Existence de l'action**: L'action doit être définie
2. **Nombre d'arguments**: Doit respecter les paramètres requis/optionnels
3. **Types des arguments**: Doivent correspondre aux types des paramètres
4. **Variables**: Les variables utilisées doivent exister dans le contexte de la règle

**Exemples:**

```tsd
type Person(name: string, age: number)

action log(message: string)
action notify(recipient: string, message: string, priority: number = 1)

// ✓ Valide: type correct (p.name est string)
rule r1 : {p: Person} / p.age > 18 ==> log(p.name)

// ✓ Valide: valeur par défaut pour priority
rule r2 : {p: Person} / p.age > 18 ==> notify(p.name, "Welcome")

// ❌ ERREUR: type incorrect (p.age est number, attendu string)
rule r3 : {p: Person} / p.age > 18 ==> log(p.age)

// ❌ ERREUR: action non définie
rule r4 : {p: Person} / p.age > 18 ==> unknownAction(p)

// ❌ ERREUR: nombre d'arguments insuffisant
rule r5 : {p: Person} / p.age > 18 ==> notify(p.name)
```

## Programme complet avec nouvelle syntaxe

Voici un exemple complet utilisant toutes les fonctionnalités:

```tsd
// Définition des types
type Person(name: string, age: number, active: bool, email: string)
type Order(orderId: number, customerEmail: string, total: number, paid: bool)
type Alert(level: string, message: string, timestamp: number)

// Définition des actions
action log(message: string)
action sendEmail(recipient: string, subject: string, body: string)
action notify(recipient: string, message: string, priority: number = 1)
action createAlert(level: string, message: string)
action processOrder(order: Order, discount: number?, notifyCustomer: bool = true)
action updatePerson(person: Person, active: bool = true)

// Règles utilisant les actions
rule checkAge : {p: Person} / p.age > 18 AND p.active == true 
    ==> log(p.name)

rule seniorDiscount : {p: Person} / p.age > 65 
    ==> notify(p.email, "Senior discount available", 2)

rule largeOrder : {o: Order} / o.paid == false AND o.total > 100 
    ==> processOrder(o, 10)

rule vipCustomer : {p: Person, o: Order} / p.email == o.customerEmail AND o.total > 500 
    ==> sendEmail(p.email, "VIP Order", "Thank you for your order!"), 
        createAlert("INFO", "VIP order processed")

// Faits
Person(name: "Alice", age: 30, active: true, email: "alice@example.com")
Person(name: "Bob", age: 70, active: true, email: "bob@example.com")
Order(orderId: 1001, customerEmail: "alice@example.com", total: 600, paid: false)
```

## Migration de l'ancienne syntaxe

### Script de conversion

Un script de conversion automatique est disponible pour migrer l'ancienne syntaxe vers la nouvelle:

```bash
./scripts/convert_syntax.sh
```

Ce script:
- Trouve tous les fichiers `.tsd` et `.constraint`
- Convertit `type Name : <...>` en `type Name(...)`
- Crée des backups automatiquement
- Affiche un rapport de conversion

### Compatibilité

- ✅ **Rétrocompatibilité**: Les règles existantes continuent de fonctionner
- ✅ **Validation améliorée**: La nouvelle syntaxe ajoute des validations supplémentaires
- ⚠️ **Actions**: Les actions doivent maintenant être **définies** avant utilisation

### Points d'attention lors de la migration

1. **Définir toutes les actions**: Ajoutez des définitions `action` pour toutes les actions utilisées
2. **Types personnalisés**: Assurez-vous que tous les types sont définis avant utilisation
3. **Validation stricte**: Les erreurs de type sont maintenant détectées au parsing

## Avantages de la nouvelle syntaxe

### Pour les types

- ✅ **Plus naturelle**: Ressemble à une signature de fonction
- ✅ **Plus lisible**: Moins de caractères spéciaux
- ✅ **Cohérence**: Même style que les actions et fonctions

### Pour les actions

- ✅ **Documentation explicite**: Les signatures servent de contrat d'interface
- ✅ **Validation au parsing**: Détection précoce des erreurs
- ✅ **Types forts**: Vérification de compatibilité des types
- ✅ **Flexibilité**: Support pour optionnels et valeurs par défaut
- ✅ **Réutilisabilité**: Les actions sont clairement définies et documentées

## Référence rapide

### Types

```tsd
type TypeName(field1: type1, field2: type2, ...)
```

Types primitifs: `string`, `number`, `bool`

### Actions

```tsd
action ActionName(param1: type1, param2: type2?, param3: type3 = default, ...)
```

- `?` après le type = paramètre optionnel
- `= valeur` = valeur par défaut
- Types: primitifs ou personnalisés

### Règles avec actions

```tsd
rule ruleName : {var: Type} / condition ==> action1(args), action2(args)
```

Multiple actions séparées par des virgules.

## Exemples additionnels

### Système de monitoring

```tsd
type Server(id: string, name: string, cpu: number, memory: number)
type Alert(serverId: string, type: string, severity: number)

action sendAlert(server: Server, message: string, severity: number = 1)
action restartService(serverId: string, serviceName: string = "main")
action notifyAdmin(message: string, urgent: bool = false)

rule highCPU : {s: Server} / s.cpu > 90 
    ==> sendAlert(s, "High CPU usage", 3), restartService(s.id)

rule criticalMemory : {s: Server} / s.memory > 95 
    ==> sendAlert(s, "Critical memory", 5), notifyAdmin("Server critical", true)
```

### E-commerce

```tsd
type Product(id: number, name: string, price: number, stock: number)
type Cart(userId: string, total: number, itemCount: number)
type User(id: string, email: string, vip: bool)

action applyDiscount(cart: Cart, percentage: number)
action sendPromoEmail(user: User, promoCode: string, discount: number = 10)
action notifyLowStock(product: Product, threshold: number = 10)

rule vipDiscount : {u: User, c: Cart} / u.vip == true AND c.total > 100 
    ==> applyDiscount(c, 15), sendPromoEmail(u, "VIP15", 15)

rule lowStockAlert : {p: Product} / p.stock < 10 
    ==> notifyLowStock(p)
```

## Support et questions

Pour toute question ou problème avec la nouvelle syntaxe, veuillez:
- Consulter les tests dans `constraint/new_syntax_test.go`
- Consulter les exemples dans `examples/`
- Ouvrir une issue sur le dépôt GitHub