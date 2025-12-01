# Actions Multiples dans les Règles RETE

## Vue d'ensemble

Les règles RETE supportent maintenant **plusieurs actions** dans la partie action d'une règle. Cela permet d'exécuter une séquence d'opérations lorsque les conditions d'une règle sont satisfaites, sans avoir à créer plusieurs règles identiques.

## Syntaxe

### Actions multiples séparées par des virgules

```
rule <nom> : <patterns> / <contraintes> ==> action1(...), action2(...), action3(...)
```

### Exemples

#### Exemple simple - deux actions
```
type Person : <id: string, name: string, age: number>

rule adult_check : {p: Person} / p.age >= 18 ==> mark_adult(p.id), log("Adult detected")
```

#### Exemple avec trois actions
```
type Person : <id: string, salary: number>

rule high_earner : {p: Person} / p.salary > 50000 ==> 
    flag_high_earner(p.id), 
    update_stats(p.salary), 
    notify_manager("High earner found")
```

#### Exemple avec agrégation
```
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_stats : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> 
    print_stats(d.name, avg_sal), 
    update_dashboard(d.id, avg_sal), 
    alert_hr(d.id)
```

#### Exemple avec arguments complexes
```
type Person : <id: string, name: string, age: number, salary: number>

rule bonus_calculation : {p: Person} / p.age > 18 AND p.salary < 40000 ==> 
    calculate_bonus(p.id, p.salary * 1.1), 
    log("Bonus calculated", p.name), 
    notify_payroll(p.id)
```

## Ordre d'exécution

Les actions sont exécutées **dans l'ordre** où elles apparaissent dans la règle, de gauche à droite.

```
rule example : {p: Person} / p.age > 18 ==> action1(p.id), action2(p.id), action3(p.id)
```

Dans cet exemple :
1. `action1(p.id)` est exécutée en premier
2. `action2(p.id)` est exécutée ensuite
3. `action3(p.id)` est exécutée en dernier

## Rétrocompatibilité

Les règles avec une seule action continuent de fonctionner exactement comme avant :

```
rule simple : {p: Person} / p.age > 18 ==> adult(p.id)
```

Cette syntaxe est entièrement compatible avec les versions antérieures.

## Arguments des actions

Chaque action peut avoir ses propres arguments, qui peuvent être :

### Variables complètes
```
rule r1 : {p: Person} / p.age > 18 ==> process(p), log(p)
```

### Accès aux champs
```
rule r2 : {p: Person} / p.age > 18 ==> process(p.id), log(p.name)
```

### Constantes
```
rule r3 : {p: Person} / p.age > 18 ==> log("Adult"), notify("admin")
```

### Expressions arithmétiques
```
rule r4 : {p: Person} / p.salary > 50000 ==> calculate_bonus(p.salary * 1.1), log("Bonus")
```

### Mélange de types
```
rule r5 : {p: Person} / p.age > 18 ==> process(p.id, "adult", p.age + 1)
```

## Structure interne (JSON)

### Format avec actions multiples (nouveau)
```json
{
  "action": {
    "type": "action",
    "jobs": [
      {
        "type": "jobCall",
        "name": "action1",
        "args": ["arg1"]
      },
      {
        "type": "jobCall",
        "name": "action2",
        "args": ["arg2"]
      }
    ]
  }
}
```

### Format avec action unique (rétrocompatibilité)
```json
{
  "action": {
    "type": "action",
    "job": {
      "type": "jobCall",
      "name": "action1",
      "args": ["arg1"]
    }
  }
}
```

## API Go

### Type Action

```go
type Action struct {
    Type string    `json:"type"`
    Job  *JobCall  `json:"job,omitempty"`  // Action unique (rétrocompatibilité)
    Jobs []JobCall `json:"jobs,omitempty"` // Actions multiples (nouveau format)
}
```

### Méthode GetJobs()

Pour obtenir la liste des jobs, utilisez la méthode `GetJobs()` qui gère automatiquement les deux formats :

```go
jobs := action.GetJobs()
for _, job := range jobs {
    fmt.Printf("Job: %s, Args: %v\n", job.Name, job.Args)
}
```

Cette méthode :
- Retourne `Jobs` si le champ est présent et non vide (nouveau format)
- Retourne `[]JobCall{*Job}` si seul `Job` est présent (ancien format)
- Retourne `[]JobCall{}` si aucun job n'est présent

## Cas d'usage

### Logging et notification
```
rule critical_alert : {s: Server} / s.cpu > 90 ==> 
    log_alert(s.id, s.cpu), 
    notify_admin(s.id), 
    trigger_auto_scale(s.id)
```

### Mise à jour de plusieurs systèmes
```
rule new_customer : {c: Customer} / c.status == "new" ==> 
    create_account(c.id), 
    send_welcome_email(c.email), 
    update_crm(c.id), 
    notify_sales_team(c.id)
```

### Calculs et reporting
```
rule monthly_report : {d: Department, total: SUM(e.sales)} / {e: Employee} / e.deptId == d.id ==> 
    calculate_commission(d.id, total), 
    generate_report(d.id, total), 
    send_to_management(d.id)
```

## Validation

Le système valide automatiquement que :
- Toutes les variables utilisées dans les actions existent dans les patterns
- Les actions sont correctement formatées
- Les arguments sont valides

Erreur de validation exemple :
```
Erreur: action process: argument contient la variable 'x' qui ne correspond à aucune variable de l'expression
```

## Limites

Il n'y a **pas de limite** au nombre d'actions par règle, mais pour des raisons de lisibilité et de maintenabilité, il est recommandé de :
- Limiter à 3-5 actions par règle maximum
- Grouper les actions logiquement liées
- Créer des règles séparées pour des séquences d'actions très différentes

## Bonnes pratiques

### ✅ Recommandé
```
rule process_order : {o: Order} / o.status == "pending" ==> 
    validate_order(o.id), 
    charge_payment(o.customerId, o.amount), 
    send_confirmation(o.customerId)
```

### ⚠️ À éviter
```
// Trop d'actions, difficile à maintenir
rule process_order : {o: Order} / o.status == "pending" ==> 
    validate_order(o.id), charge_payment(o.customerId, o.amount), 
    send_confirmation(o.customerId), update_inventory(o.items), 
    log_transaction(o.id), notify_warehouse(o.id), 
    update_analytics(o.id), send_sms(o.customerId), 
    create_invoice(o.id), archive_data(o.id)
```

## Migration depuis l'ancien format

Les règles existantes continuent de fonctionner sans modification :

```
// Ancien format (toujours supporté)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)

// Nouveau format (avec une seule action)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)

// Nouveau format (avec plusieurs actions)
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id), log("Adult detected")
```

Aucune modification du code existant n'est nécessaire.

## Références

- [Grammaire PEG](../constraint/grammar/constraint.peg) - Définition de la syntaxe
- [Types d'actions](../constraint/constraint_types.go) - Structures de données
- [Tests d'actions multiples](../constraint/multiple_actions_test.go) - Exemples de tests
- [Exemple complet](../examples/multiple_actions_example.tsd) - Fichier d'exemple

## Copyright

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License