# Tutoriel : Écriture des Actions dans les Contraintes

## Introduction

Dans le système de contraintes, les **actions** sont déclenchées lorsqu'une règle de contrainte est satisfaite. Elles permettent d'exécuter des opérations spécifiques en réponse aux conditions détectées par le moteur RETE.

## Syntaxe Fondamentale

### Structure de base d'une action

```
{variables} / condition ==> action_name(arguments)
```

L'opérateur `==>` sépare la condition de l'action à exécuter.

## Types d'Actions

### 1. Actions Simples

#### Action sans paramètres
```constraint
{u: User} / u.active == false ==> deactivate_account()
```

#### Action avec un paramètre
```constraint
{s: System} / s.cpu_usage > 90 ==> restart_service(s.id)
```

#### Action avec plusieurs paramètres
```constraint
{a: Alarm} / a.severity == "critical" ==> alert_team(a.id, a.source)
```

### 2. Actions avec Accès aux Champs (Field Access)

Les arguments peuvent référencer directement les champs des objets impliqués dans la contrainte :

```constraint
{p: TestPerson} / p.name == "Alice" ==> test_string_equality_success(p.id)
```

```constraint
{t: Transaction} / t.amount > 1000 ==> flag_large_transaction(t.id, t.amount)
```

```constraint
{s: System, i: Incident} / s.id == i.system_id AND s.memory_usage > 95 ==> backup_and_restart(s.id, s.name, i.id)
```

### 3. Actions avec Chaînes Littérales

Les actions peuvent inclure des valeurs littérales comme des chaînes de caractères :

```constraint
{u: Utilisateur, a: Adresse} / u.id == a.utilisateur_id AND u.age < 18 AND a.ville == "Lille" ==> alert_mineur_lille(u.id, u.nom, u.prenom, a.rue)
```

```constraint
{ba: BankAccount, c: Customer} / ba.customer_id == c.id AND ba.balance < -1000 ==> freeze_account(ba.id, c.id, "fraud_detected")
```

## Exemples Pratiques par Contexte

### Alpha Nodes (Un seul type d'objet)

```constraint
// Action simple avec un argument
{t: Transaction} / t.status == "approved" ==> process_transaction(t.id)

// Action avec multiple arguments et valeur littérale
{a: Account} / a.active == true AND a.balance >= 0 ==> activate_account_services(a.id)

// Action avec fonction intégrée dans la condition
{t: Transaction} / LENGTH(t.id) == 8 ==> validate_id_length(t.id)
```

### Beta Nodes (Jointures entre objets)

```constraint
// Jointure simple avec action
{c: Customer, o: Order} / c.id == o.customer_id ==> link_customer_order(c.id, o.id)

// Jointure complexe avec conditions multiples
{c: Customer, o: Order} / c.vip == true AND o.total > 1000 AND c.age >= 18 ==> apply_vip_benefits(c.id, o.id)

// Jointure avec calculs arithmétiques
{c: Customer, o: Order} / c.age * 10 >= o.total / 100 ==> calculate_age_discount(c.id, o.id)
```

### Nodes avec Négation (NOT)

```constraint
// Négation simple
{u: User} / NOT (u.last_login > 1700000000) ==> flag_inactive_user(u.id)

// Négation avec jointure
{u: User, l: Login} / u.id == l.user_id AND NOT (l.success == true AND l.timestamp > 1700000000) ==> notify_login_failure(u.id)
```

### Nodes Existentiels (EXISTS)

```constraint
// Vérification d'existence simple
{a: Account} / EXISTS (t: Transaction / t.account_id == a.id AND t.amount > 10000) ==> flag_suspicious_account(a.id)

// Existence avec conditions complexes
{a: Account} / EXISTS (t: Transaction / t.account_id == a.id AND t.amount > 5000 AND t.type == "withdrawal") ==> monitor_large_withdrawal(a.id)
```

### Nodes d'Agrégation

```constraint
// Agrégation SUM
{p: Portfolio, a: Asset} / p.id == a.portfolio_id AND SUM(a.value) > 100000 ==> flag_high_value_portfolio(p.id)

// Agrégation COUNT
{p: Portfolio, a: Asset} / p.id == a.portfolio_id AND COUNT(a.id) >= 10 ==> apply_diversification_bonus(p.id)

// Agrégation AVG
{p: Portfolio, a: Asset} / p.id == a.portfolio_id AND AVG(a.risk_level) < 5 ==> recommend_conservative_strategy(p.id)

// Agrégations multiples
{p: Portfolio, a: Asset} / p.id == a.portfolio_id AND SUM(a.value) > 100000 AND AVG(a.risk_level) < 7 AND COUNT(a.id) >= 5 ==> award_premium_status(p.id)
```

## Règles de Nommage et Bonnes Pratiques

### 1. Noms d'Actions

Les noms d'actions doivent être :
- **Descriptifs** : `process_transaction`, `flag_large_transaction`
- **En snake_case** : `apply_vip_benefits`, `calculate_age_discount`
- **Verbes d'action** : `validate`, `process`, `notify`, `escalate`

### 2. Arguments

Les arguments peuvent être :
- **Champs d'objets** : `p.id`, `o.amount`, `s.name`
- **Chaînes littérales** : `"fraud_detected"`, `"critical"`
- **Nombres** : `1000`, `3.14`
- **Booléens** : `true`, `false`

### 3. Ordre des Arguments

Conventionnellement :
1. **Identifiants principaux** en premier : `(p.id, o.id)`
2. **Valeurs contextuelles** ensuite : `(p.id, p.name, o.total)`
3. **Chaînes de statut** en dernier : `(s.id, s.name, "emergency")`

## Exemples d'Actions Complexes

### Action avec de nombreux paramètres

```constraint
{u: Utilisateur, a: Adresse} / u.id == a.utilisateur_id AND u.age >= 18 AND a.ville == "Paris" ==> process_majeur_paris(u.id, u.nom, u.prenom, a.rue)
```

### Action avec calculs dans les paramètres

```constraint
{s: System, i: Incident} / s.id == i.system_id AND s.memory_usage > 95 ==> backup_and_restart(s.id, s.name, i.id)
```

### Action avec conditions très complexes

```constraint
{p: Portfolio, a: Asset, t: Trade} / p.id == a.portfolio_id AND a.id == t.asset_id AND SUM(a.value) > AVG(t.amount * t.price) * 10 ==> calculate_performance_bonus(p.id)
```

## Grammaire Formelle

Selon la grammaire PEG du système :

```peg
Action       = "==>" _ job:JobCall
JobCall      = name:IdentName args:ArgumentList?
ArgumentList = "(" args:ArgumentsContent? ")"
ArgumentsContent = first:Argument rest:("," _ arg:Argument)*
Argument     = Number / String / Boolean / FieldAccess / IdentName
FieldAccess  = object:IdentName "." field:IdentName
```

## Gestion des Types d'Arguments

Le système gère automatiquement différents types d'arguments :

- **String** : `"fraud_detected"`, `"critical"`
- **Number** : `1000`, `3.14159`, `-500`
- **Boolean** : `true`, `false`
- **FieldAccess** : `p.id`, `o.customer_id`, `s.cpu_usage`

Les arguments de type `FieldAccess` sont résolus automatiquement lors de l'exécution pour extraire la valeur du champ correspondant de l'objet.

## Conseils de Performance

1. **Minimisez le nombre d'arguments** pour réduire la surcharge
2. **Utilisez des noms explicites** pour faciliter la maintenance
3. **Groupez les actions logiquement liées** dans le même fichier de contraintes
4. **Évitez les calculs complexes dans les arguments** - faites-les dans l'action elle-même

## Conclusion

Les actions dans les contraintes offrent une interface flexible pour déclencher des opérations métier en réponse aux patterns détectés par le moteur RETE. La syntaxe `==>` fournit une séparation claire entre la logique de détection et l'exécution, permettant des systèmes de règles maintenables et expressifs.