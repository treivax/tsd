# Identifiants de Règles (Rule Identifiers)

## Vue d'ensemble

À partir de cette version, **toutes les règles doivent obligatoirement posséder un identifiant unique**. Cette fonctionnalité permet une meilleure gestion des règles, notamment pour leur suppression, leur suivi et leur maintenance.

## Syntaxe

### Ancienne syntaxe (obsolète)

```
{p: Person} / p.age > 18 ==> adult(p.id)
```

### Nouvelle syntaxe (obligatoire)

```
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
```

## Format de la règle

Une règle suit maintenant le format suivant :

```
rule <IDENTIFIANT> : <ENSEMBLE_VARIABLES> / <CONDITIONS> ==> <ACTION>
```

Où :
- **`rule`** : Mot-clé obligatoire
- **`<IDENTIFIANT>`** : Identifiant unique de la règle (ex: `r1`, `check_age`, `fraud_detection`)
- **`:`** : Séparateur obligatoire
- **`<ENSEMBLE_VARIABLES>`** : Ensemble des variables typées entre accolades `{...}`
- **`/`** : Séparateur entre variables et conditions
- **`<CONDITIONS>`** : Conditions logiques à évaluer
- **`==>`** : Séparateur avant l'action
- **`<ACTION>`** : Action à exécuter si les conditions sont satisfaites

## Règles de nommage des identifiants

### Contraintes techniques

1. **Unicité** : Chaque identifiant de règle doit être unique dans le programme
2. **Format** : Un identifiant est un nom valide (lettres, chiffres, underscore)
3. **Commence par une lettre ou underscore** : `r1`, `_rule`, `checkAge` sont valides
4. **Sensible à la casse** : `R1` et `r1` sont considérés comme différents

### Conventions recommandées

#### Convention simple numérique
```
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)
rule r3 : {o: Order} / o.amount > 1000 ==> premium_order(o.id)
```

#### Convention descriptive
```
rule check_adult_age : {p: Person} / p.age >= 18 ==> adult(p.id)
rule validate_order : {o: Order} / o.amount > 0 ==> process_order(o.id)
rule detect_fraud : {t: Transaction} / t.amount > 10000 ==> flag_suspicious(t.id)
```

#### Convention par domaine
```
rule person_adult_check : {p: Person} / p.age >= 18 ==> adult(p.id)
rule person_senior_check : {p: Person} / p.age >= 65 ==> senior(p.id)
rule order_premium : {o: Order} / o.amount > 1000 ==> premium_order(o.id)
rule order_bulk : {o: Order} / o.quantity > 100 ==> bulk_order(o.id)
```

## Exemples complets

### Règle simple (Alpha)

```
type Person : <id: string, name: string, age: number>

rule r1 : {p: Person} / p.age >= 18 ==> adult(p.id, p.name)
```

### Règle avec jointure (Beta)

```
type Person : <id: string, name: string>
type Order : <id: string, customer_id: string, amount: number>

rule r1 : {p: Person, o: Order} / p.id == o.customer_id AND o.amount > 100 ==> premium_order(p.id, o.id)
```

### Règle avec négation (NOT)

```
type Person : <id: string, age: number>

rule r1 : {p: Person} / p.age >= 18 AND NOT(p.age > 65) ==> working_age(p.id)
```

### Règle avec quantificateur existentiel (EXISTS)

```
type Person : <id: string, name: string>
type Order : <id: string, customer_id: string, amount: number>

rule r1 : {p: Person} / EXISTS (o: Order / o.customer_id == p.id AND o.amount > 100) ==> has_premium_orders(p.id)
```

### Règle avec agrégation

```
type Person : <id: string, name: string>
type Order : <id: string, customer_id: string, amount: number>

rule r1 : {p: Person} / SUM(o: Order / o.customer_id == p.id ; o.amount) >= 1000 ==> vip_customer(p.id)
```

### Règle complexe multi-variables

```
type User : <id: string, age: number>
type Order : <id: string, user_id: string, product_id: string>
type Product : <id: string, name: string, price: number>

rule r1 : {u: User, o: Order, p: Product} / 
    u.id == o.user_id AND 
    o.product_id == p.id AND 
    u.age >= 18 AND 
    p.price > 100 
    ==> eligible_purchase(u.id, o.id, p.id)
```

## Avantages des identifiants de règles

### 1. Gestion des règles

Les identifiants permettent de :
- **Supprimer** une règle spécifique par son identifiant
- **Modifier** une règle en la rechargeant avec le même identifiant
- **Suivre** l'exécution d'une règle dans les logs
- **Déboguer** plus facilement en identifiant quelle règle pose problème

### 2. Traçabilité

```json
{
  "ruleId": "r1",
  "type": "expression",
  "set": { ... },
  "constraints": { ... },
  "action": { ... }
}
```

Le champ `ruleId` est maintenant présent dans toutes les structures JSON des règles.

### 3. Opérations futures

Les identifiants ouvrent la voie à de nouvelles fonctionnalités :
- **Suppression dynamique** : `remove rule r1`
- **Statistiques** : nombre d'activations par règle
- **Profiling** : temps d'exécution par règle
- **Priorisation** : ordre d'exécution des règles

## Migration des règles existantes

### Script de migration automatique

Un script de migration est fourni pour mettre à jour automatiquement tous vos fichiers `.constraint` :

```bash
cd tsd
bash scripts/add_rule_ids.sh
```

Ce script :
1. Parcourt tous les fichiers `.constraint` du projet
2. Détecte les règles sans identifiant
3. Ajoute automatiquement des identifiants séquentiels (`r1`, `r2`, etc.)
4. Préserve les règles déjà migrées

### Exemple de migration

**Avant :**
```
{p: Person} / p.age > 18 ==> adult(p.id)
{p: Person} / p.age < 18 ==> minor(p.id)
```

**Après (automatique) :**
```
rule r1 : {p: Person} / p.age > 18 ==> adult(p.id)
rule r2 : {p: Person} / p.age < 18 ==> minor(p.id)
```

**Après (manuel avec noms descriptifs) :**
```
rule check_adult : {p: Person} / p.age > 18 ==> adult(p.id)
rule check_minor : {p: Person} / p.age < 18 ==> minor(p.id)
```

## Validation et erreurs

### Erreur : Règle sans identifiant

```
Erreur: parsing: 3:1 (41): no match found, expected: "rule"
```

**Solution :** Ajouter le préfixe `rule <ID> :` avant la règle.

### Erreur : Identifiant dupliqué

Si deux règles ont le même identifiant, le parser peut réussir mais cela peut causer des conflits lors de la gestion des règles.

**Bonne pratique :** Utiliser des identifiants uniques et significatifs.

### Avertissement : Ancien format détecté

Si vous utilisez l'ancien format sans identifiant, le parser affichera une erreur explicite.

## API JSON

### Structure d'une règle avec identifiant

```json
{
  "type": "expression",
  "ruleId": "r1",
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
  "constraints": {
    "type": "comparison",
    "left": {
      "type": "fieldAccess",
      "object": "p",
      "field": "age"
    },
    "operator": ">",
    "right": {
      "type": "number",
      "value": 18
    }
  },
  "action": {
    "type": "action",
    "job": {
      "type": "jobCall",
      "name": "adult",
      "args": [
        {
          "type": "fieldAccess",
          "object": "p",
          "field": "id"
        }
      ]
    }
  }
}
```

## Compatibilité

### Version introduite

Cette fonctionnalité a été introduite dans la version **2.0.0** du projet TSD.

### Rétrocompatibilité

⚠️ **Attention** : Cette modification **casse la rétrocompatibilité** avec l'ancienne syntaxe.

Tous les fichiers `.constraint` existants doivent être migrés pour utiliser la nouvelle syntaxe avec identifiants de règles.

### Outils de migration

- **Script automatique** : `scripts/add_rule_ids.sh`
- **Documentation** : Ce fichier
- **Tests** : 344 règles migrées avec succès dans la suite de tests

## Bonnes pratiques

### ✅ Faire

1. **Utiliser des identifiants descriptifs** pour les règles complexes
2. **Préfixer par domaine** pour grouper les règles similaires
3. **Documenter** les règles importantes avec des commentaires
4. **Maintenir** une liste des identifiants utilisés dans chaque fichier

### ❌ Éviter

1. **Identifiants génériques** comme `rule1`, `rule2` sans contexte
2. **Identifiants trop longs** qui rendent la règle difficile à lire
3. **Caractères spéciaux** non standards dans les identifiants
4. **Duplication** d'identifiants entre fichiers différents si possible

## Exemples d'organisation

### Petit projet (< 20 règles)

```
rule r1 : ...
rule r2 : ...
rule r3 : ...
```

### Projet moyen (20-100 règles)

```
// Règles de validation
rule validate_person : ...
rule validate_order : ...

// Règles métier
rule business_discount : ...
rule business_vip : ...

// Règles de sécurité
rule security_fraud : ...
rule security_limit : ...
```

### Grand projet (> 100 règles)

Utiliser plusieurs fichiers avec préfixes :

**validation_rules.constraint**
```
rule val_person_age : ...
rule val_person_email : ...
rule val_order_amount : ...
```

**business_rules.constraint**
```
rule biz_discount_vip : ...
rule biz_discount_bulk : ...
rule biz_loyalty_points : ...
```

**security_rules.constraint**
```
rule sec_fraud_high_amount : ...
rule sec_fraud_foreign : ...
rule sec_rate_limit : ...
```

## Référence rapide

| Élément | Obligatoire | Format | Exemple |
|---------|-------------|--------|---------|
| Mot-clé `rule` | ✅ Oui | `rule` | `rule` |
| Identifiant | ✅ Oui | `[a-zA-Z_][a-zA-Z0-9_]*` | `r1`, `check_age` |
| Séparateur `:` | ✅ Oui | `:` | `:` |
| Variables | ✅ Oui | `{var: Type, ...}` | `{p: Person}` |
| Séparateur `/` | ✅ Oui | `/` | `/` |
| Conditions | ✅ Oui | Expression logique | `p.age > 18` |
| Séparateur `==>` | ✅ Oui | `==>` | `==>` |
| Action | ✅ Oui | `fonction(args)` | `adult(p.id)` |

## Support et questions

Pour toute question ou problème lié aux identifiants de règles :

1. Consultez les [exemples de tests](../beta_coverage_tests/)
2. Consultez le [CHANGELOG](../CHANGELOG.md)
3. Ouvrez une issue sur le dépôt GitHub

## Voir aussi

- [Grammaire PEG complète](../constraint/grammar/constraint.peg)
- [Types de contraintes](./constraints.md)
- [Guide de migration](./migration_guide.md)
- [Architecture RETE](./rete_architecture.md)