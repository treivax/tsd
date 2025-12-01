# Exécution des Actions dans TSD

## Vue d'ensemble

Le système TSD implémente maintenant l'exécution réelle des actions déclenchées par les règles RETE. Chaque action est systématiquement loggée avec son nom et ses arguments, et supporte cinq types d'arguments différents avec validation de cohérence des types.

## Fonctionnalités

### 1. Logging Systématique

Toutes les actions sont automatiquement loggées lors de leur exécution :

```
📋 ACTION: log("Adult detected")
🎯 ACTION EXÉCUTÉE: log("Adult detected")

📋 ACTION: notify(p.name)
🎯 ACTION EXÉCUTÉE: notify("Alice")

📋 ACTION: process(p)
🎯 ACTION EXÉCUTÉE: process(Person{p1})
```

### 2. Types d'Arguments Supportés

#### 2.1 Valeurs Littérales

Chaînes de caractères, nombres et booléens directement dans l'action.

**Syntaxe :**
```
rule r1 : {p: Person} / p.age > 18 ==> log("Adult detected"), notify("admin")
```

**Exemple :**
```tsd
type Person : <id: string, age: number>

rule log_adults : {p: Person} / p.age >= 18 ==> log("Adult person detected")

Person(id:"p1", age:25)
```

**Output :**
```
📋 ACTION: log("Adult person detected")
🎯 ACTION EXÉCUTÉE: log("Adult person detected")
```

#### 2.2 Fait Complet (Variable)

Référence à un fait complet via son nom de variable dans la règle.

**Syntaxe :**
```
rule r2 : {p: Person} / p.age > 18 ==> process(p)
```

**Exemple :**
```tsd
type Person : <id: string, name: string, age: number>

rule process_adult : {p: Person} / p.age > 21 ==> process(p), archive(p)

Person(id:"p1", name:"Alice", age:25)
```

**Output :**
```
📋 ACTION: process(p)
🎯 ACTION EXÉCUTÉE: process(Person{p1})
📋 ACTION: archive(p)
🎯 ACTION EXÉCUTÉE: archive(Person{p1})
```

#### 2.3 Attribut de Fait (variable.attribut)

Accès à un champ spécifique d'un fait.

**Syntaxe :**
```
rule r3 : {p: Person} / p.salary > 50000 ==> notify_hr(p.name), log_salary(p.salary)
```

**Exemple :**
```tsd
type Person : <id: string, name: string, salary: number>

rule high_earner : {p: Person} / p.salary > 50000 ==> notify_hr(p.name), log_salary(p.salary)

Person(id:"p1", name:"Bob", salary:65000)
```

**Output :**
```
📋 ACTION: notify_hr(p.name)
🎯 ACTION EXÉCUTÉE: notify_hr("Bob")
📋 ACTION: log_salary(p.salary)
🎯 ACTION EXÉCUTÉE: log_salary(65000)
```

#### 2.4 Expressions Arithmétiques

Calculs mathématiques dans les arguments.

**Syntaxe :**
```
rule r4 : {p: Person} / p.salary < 40000 ==> calculate_bonus(p.id, p.salary * 1.1)
```

**Opérateurs supportés :**
- `+` : Addition
- `-` : Soustraction
- `*` : Multiplication
- `/` : Division

**Exemple :**
```tsd
type Person : <id: string, salary: number>

rule bonus : {p: Person} / p.salary < 40000 ==> 
    calculate_bonus(p.id, p.salary * 1.1),
    log("Bonus calculated")

Person(id:"p1", salary:35000)
```

**Output :**
```
📋 ACTION: calculate_bonus(p.id, p.salary * 1.1)
🎯 ACTION EXÉCUTÉE: calculate_bonus("p1", 38500)
📋 ACTION: log("Bonus calculated")
🎯 ACTION EXÉCUTÉE: log("Bonus calculated")
```

#### 2.5 Arguments Mixtes

Combinaison de différents types d'arguments dans une même action.

**Exemple :**
```tsd
type Person : <id: string, name: string, age: number, salary: number>

rule mixed_args : {p: Person} / p.age > 25 ==> 
    process(p, p.name, "active", p.salary * 1.05)

Person(id:"p1", name:"Charlie", age:30, salary:50000)
```

**Output :**
```
📋 ACTION: process(p, p.name, "active", p.salary * 1.05)
🎯 ACTION EXÉCUTÉE: process(Person{p1}, "Charlie", "active", 52500)
```

## Validation de Cohérence

### Validation des Variables

Les variables utilisées dans les actions doivent être définies dans les patterns de la règle.

**✅ Valide :**
```tsd
rule r1 : {p: Person} / p.age > 18 ==> process(p.name)
```

**❌ Invalide :**
```tsd
rule r1 : {p: Person} / p.age > 18 ==> process(x.name)
// Erreur: variable 'x' non trouvée
```

### Validation des Attributs

Les attributs utilisés doivent exister dans la définition du type.

**✅ Valide :**
```tsd
type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> notify(p.name)
```

**❌ Invalide :**
```tsd
type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age > 18 ==> notify(p.email)
// Erreur: champ 'email' non trouvé dans le type Person
```

### Validation des Types de Valeurs

Les valeurs doivent correspondre au type défini dans le schéma.

**Définition de type :**
```tsd
type Person : <id: string, name: string, age: number>
```

**Validation automatique :**
- `id` doit être une chaîne de caractères
- `name` doit être une chaîne de caractères
- `age` doit être un nombre

**Erreurs détectées :**
```
❌ type attendu: string, reçu: number
❌ type attendu: number, reçu: string
❌ champ requis 'id' manquant
```

## Actions Multiples en Séquence

Les actions multiples sont exécutées dans l'ordre où elles apparaissent.

**Exemple :**
```tsd
type Person : <id: string, name: string, status: string>

rule onboarding : {p: Person} / p.status == "new" ==>
    create_account(p.id),
    send_welcome_email(p.name),
    log("Onboarding complete")

Person(id:"p1", name:"Alice", status:"new")
```

**Output (dans l'ordre) :**
```
📋 ACTION: create_account(p.id)
🎯 ACTION EXÉCUTÉE: create_account("p1")
📋 ACTION: send_welcome_email(p.name)
🎯 ACTION EXÉCUTÉE: send_welcome_email("Alice")
📋 ACTION: log("Onboarding complete")
🎯 ACTION EXÉCUTÉE: log("Onboarding complete")
```

## Actions avec Agrégation

Les actions peuvent utiliser les variables d'agrégation.

**Exemple :**
```tsd
type Department : <id: string, name: string>
type Employee : <id: string, deptId: string, salary: number>

rule dept_stats : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==>
    print_stats(d.name, avg_sal),
    update_dashboard(d.id, avg_sal),
    alert_if_high(avg_sal)

Department(id:"d1", name:"Engineering")
Employee(id:"e1", deptId:"d1", salary:60000)
Employee(id:"e2", deptId:"d1", salary:70000)
```

**Output :**
```
📋 ACTION: print_stats(d.name, avg_sal)
🎯 ACTION EXÉCUTÉE: print_stats("Engineering", 65000)
📋 ACTION: update_dashboard(d.id, avg_sal)
🎯 ACTION EXÉCUTÉE: update_dashboard("d1", 65000)
📋 ACTION: alert_if_high(avg_sal)
🎯 ACTION EXÉCUTÉE: alert_if_high(65000)
```

## Gestion des Erreurs

### Erreurs de Variables

```
❌ erreur exécution job process: erreur évaluation argument 0: variable 'unknown' non trouvée
```

### Erreurs de Champs

```
❌ erreur exécution job notify: erreur évaluation argument 0: champ 'email' non trouvé dans le fait p
```

### Erreurs Arithmétiques

```
❌ erreur exécution job calculate: erreur évaluation argument 0: division par zéro
❌ erreur exécution job calculate: opération arithmétique nécessite des nombres
```

### Erreurs de Type

```
❌ validation field modification: champ 'age': type attendu: number, reçu: string
❌ validation fact creation: champ requis 'id' manquant
```

## Architecture Interne

### ActionExecutor

Le composant principal qui gère l'exécution des actions.

**Structure :**
```go
type ActionExecutor struct {
    network       *ReteNetwork
    logger        *log.Logger
    enableLogging bool
}
```

**Méthodes principales :**
- `ExecuteAction(action *Action, token *Token) error` - Exécute une action
- `SetLogging(enabled bool)` - Active/désactive le logging
- `evaluateArgument(arg interface{}, ctx *ExecutionContext) (interface{}, error)` - Évalue un argument

### ExecutionContext

Contexte d'exécution contenant les faits disponibles pour l'action.

**Structure :**
```go
type ExecutionContext struct {
    token    *Token
    network  *ReteNetwork
    varCache map[string]*Fact
}
```

**Méthodes :**
- `GetVariable(name string) *Fact` - Récupère un fait par nom de variable
- `NewExecutionContext(token *Token, network *ReteNetwork) *ExecutionContext`

### Flux d'Exécution

1. **Règle déclenchée** → Token créé avec faits correspondants
2. **TerminalNode.executeAction()** → Appelle l'ActionExecutor
3. **ActionExecutor.ExecuteAction()** → Parcourt tous les jobs
4. **Pour chaque job** :
   - Logger l'action (si activé)
   - Évaluer chaque argument
   - Valider les types
   - Exécuter l'action
   - Logger le résultat

## Configuration

### Activer/Désactiver le Logging

```go
network := NewReteNetwork(storage)
network.ActionExecutor.SetLogging(false) // Désactiver le logging
```

### Logger Personnalisé

```go
import "log"

customLogger := log.New(os.Stdout, "[ACTIONS] ", log.LstdFlags)
network.ActionExecutor = NewActionExecutor(network, customLogger)
```

## Exemples Complets

### Exemple 1 : Système de Notifications

```tsd
type User : <id: string, name: string, email: string, age: number>
type Alert : <id: string, userId: string, severity: string>

rule adult_notification : {u: User} / u.age >= 18 ==>
    send_email(u.email, "Welcome"),
    log("Email sent to adult"),
    create_profile(u.id)

rule critical_alert : {a: Alert, u: User} / a.userId == u.id AND a.severity == "critical" ==>
    notify_immediately(u.email),
    log("Critical alert processed"),
    escalate(a.id)

User(id:"u1", name:"Alice", email:"alice@example.com", age:25)
Alert(id:"a1", userId:"u1", severity:"critical")
```

### Exemple 2 : Calcul de Bonus

```tsd
type Employee : <id: string, name: string, salary: number, performance: number>

rule bonus_calculation : {e: Employee} / e.performance > 0.8 AND e.salary < 60000 ==>
    calculate_bonus(e.id, e.salary * 0.15),
    notify_payroll(e.name, e.salary * 1.15),
    log("Bonus calculated for high performer")

Employee(id:"e1", name:"Bob", salary:50000, performance:0.9)
```

### Exemple 3 : Gestion Multi-Départements

```tsd
type Department : <id: string, name: string, budget: number>
type Employee : <id: string, name: string, deptId: string, salary: number>

rule budget_check : {d: Department, total_sal: SUM(e.salary)} / {e: Employee} / e.deptId == d.id ==>
    check_budget(d.name, total_sal, d.budget),
    alert_if_over(d.id, total_sal),
    update_report(d.id)

Department(id:"d1", name:"Engineering", budget:500000)
Employee(id:"e1", name:"Alice", deptId:"d1", salary:80000)
Employee(id:"e2", name:"Bob", deptId:"d1", salary:70000)
```

## Bonnes Pratiques

### ✅ Recommandé

1. **Utiliser des noms d'actions descriptifs**
   ```tsd
   ==> send_welcome_email(p.email), create_user_profile(p.id)
   ```

2. **Grouper les actions logiquement liées**
   ```tsd
   ==> validate_order(o.id), charge_payment(o.amount), send_confirmation(o.email)
   ```

3. **Gérer les erreurs explicitement**
   ```tsd
   rule safe_division : {p: Person} / p.hours > 0 ==> 
       calculate_rate(p.salary / p.hours),
       log("Rate calculated")
   ```

4. **Utiliser des expressions arithmétiques simples**
   ```tsd
   ==> calculate_bonus(p.salary * 1.1)
   ```

### ⚠️ À Éviter

1. **Actions trop nombreuses dans une seule règle**
   ```tsd
   // ❌ Difficile à maintenir
   ==> action1(), action2(), action3(), action4(), action5(), action6()
   ```

2. **Expressions arithmétiques complexes**
   ```tsd
   // ❌ Difficile à lire
   ==> calculate(p.a * p.b / (p.c + p.d) - p.e)
   ```

3. **Dépendance entre actions**
   ```tsd
   // ❌ Si action1 échoue, action2 peut être incohérente
   ==> action1(p.id), action2(p.id)
   ```

## Tests

### Test d'Exécution Basique

```go
func TestActionExecution(t *testing.T) {
    network := NewReteNetwork(storage)
    
    // Définir les types
    network.Types = append(network.Types, personType)
    
    // Créer un fait et un token
    fact := &Fact{ID: "p1", Type: "Person", Fields: map[string]interface{}{
        "id": "p1", "name": "Alice", "age": 25.0,
    }}
    
    token := &Token{
        ID: "token1",
        Facts: []*Fact{fact},
        Bindings: map[string]*Fact{"p": fact},
    }
    
    // Créer et exécuter une action
    action := &Action{
        Type: "action",
        Jobs: []JobCall{{
            Name: "process",
            Args: []interface{}{
                map[string]interface{}{"type": "variable", "name": "p"},
            },
        }},
    }
    
    err := network.ActionExecutor.ExecuteAction(action, token)
    assert.NoError(t, err)
}
```

## Références

- [Actions Multiples](multiple_actions.md) - Documentation des actions multiples
- [Grammaire PEG](../constraint/grammar/constraint.peg) - Syntaxe complète
- [Types de Contraintes](../constraint/constraint_types.go) - Structures de données
- [Exemples](../examples/action_execution_example.tsd) - Exemples complets

## Copyright

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License