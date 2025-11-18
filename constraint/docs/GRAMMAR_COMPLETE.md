# Grammaire PEG Complète - Module Constraint

## Vue d'ensemble

Cette grammaire PEG unique et complète assure une **cohérence totale** entre les constructs de langage et les nœuds du réseau RETE. Elle supporte 100% des fichiers de contraintes existants avec parsing réel et validation sémantique.

## Fichiers du Module

```
constraint/
├── grammar/
│   └── constraint.peg          # Grammaire PEG unique et complète
├── parser.go                   # Parser généré par pigeon
├── api.go                      # API publique du module
├── constraint_types.go         # Types de données
├── constraint_utils.go         # Utilitaires
└── test/integration/           # Fichiers de test
```

## Génération du Parser

```bash
cd constraint/grammar
pigeon -o ../parser.go constraint.peg
```

## Cohérence PEG ↔ RETE

| **Construct PEG** | **Nœud RETE** | **Exemple** |
|---|---|---|
| `type Name : <fields>` | RootNode | Types de données de base |
| `{var: Type}` | AlphaNode | Variables typées |
| `field == value` | AlphaNode | Conditions simples |
| `expr1 AND expr2` | BetaNode/JoinNode | Jointures entre faits |
| `NOT(expression)` | NotNode | Négation de conditions |
| `EXISTS(var: Type / cond)` | ExistsNode | Quantification existentielle |
| `COUNT/SUM/AVG(expr)` | AccumulateNode | Agrégation de données |
| `complex_expr` | ProductionNode | Expressions logiques complexes |
| `==> action(args)` | TerminalNode | Actions finales |

## Constructs Supportés

### 1. Définitions de Types
```
type Transaction : <id: string, amount: number, status: string>
type Account : <id: string, balance: number, active: bool>
```

### 2. Expressions de Base
```
{t: Transaction} / t.amount > 1000
{a: Account} / a.active == true AND a.balance >= 0
```

### 3. Opérateurs Avancés
```
{t: Transaction} / t.status IN ["pending", "approved", "rejected"]
{a: Account} / a.type LIKE "premium%"
{t: Transaction} / t.id MATCHES "TX[0-9]+"
```

### 4. Fonctions Intégrées
```
{t: Transaction} / LENGTH(t.id) == 8
{t: Transaction} / UPPER(t.status) == "APPROVED"
{a: Account} / ABS(a.balance) > 1000
```

### 5. Négation (NotNode)
```
{u: User} / NOT (u.last_login > 1700000000)
{u: User} / u.active == true AND NOT (u.name CONTAINS "admin")
```

### 6. Existence (ExistsNode)
```
{a: Account} / EXISTS (t: Transaction / t.account_id == a.id AND t.amount > 10000)
```

### 7. Agrégation (AccumulateNode)
```
{u: User, s: SecurityEvent} / u.id == s.user_id AND COUNT(s.user_id) > 3
{p: Portfolio} / SUM(a.value WHERE a.portfolio_id == p.id) > 1000000
```

### 8. Actions (TerminalNode)
```
{a: Alarm} / a.severity == "critical" ==> alert_team(a.id, a.source)
{s: System} / s.cpu_usage > 90 ==> restart_service(s.id)
```

### 9. Commentaires
```
// Commentaire simple ligne
/* Commentaire
   multi-lignes */
```

## Validation et Tests

### Tests d'Intégration
- **6 fichiers** de contraintes complexes : 100% de réussite
- **Parsing réel** avec structures complètes extraites
- **Validation sémantique** des références de types

### Exécution des Tests
```bash
go test -run TestFlexibleParserIntegration -v advanced_integration_test.go
```

### Résultats Attendus
```
✅ alpha_conditions.constraint: 2 types, 12 expressions
✅ beta_joins.constraint: 3 types, 9 expressions
✅ negation.constraint: 3 types, 8 expressions
✅ exists.constraint: 3 types, 10 expressions
✅ aggregation.constraint: 3 types, 14 expressions
✅ actions.constraint: 3 types, 10 expressions

📊 Success rate: 6/6 (100.0%)
```

## Architecture Technique

### Parser PEG
- **Grammaire flexible** supportant syntaxe entrelacée (types + expressions)
- **Validation sémantique** intégrée au parsing
- **Support complet** des commentaires
- **Gestion d'erreurs** précise avec positions

### Structure de Sortie
```json
{
  "types": [
    {
      "type": "typeDefinition",
      "name": "Transaction",
      "fields": [...]
    }
  ],
  "expressions": [
    {
      "type": "expression",
      "set": {...},
      "constraints": {...},
      "action": {...}  // optionnel
    }
  ]
}
```

### API Publique
- `ParseConstraint(filename, input)` - Parsing depuis bytes
- `ParseConstraintFile(filename)` - Parsing depuis fichier
- `ValidateConstraintProgram(result)` - Validation post-parsing

## Cohérence Garantie

Cette grammaire unique garantit :
1. **Couverture complète** de tous les nœuds RETE
2. **Parsing réel** (pas seulement validation syntaxique)
3. **Compatibilité** avec tous les fichiers de contraintes existants
4. **Extensibilité** pour futurs constructs RETE
5. **Documentation** complète des correspondances PEG ↔ RETE

## Maintenance

Pour ajouter de nouveaux constructs :
1. Mettre à jour `constraint.peg`
2. Régénérer avec `pigeon -o ../parser.go constraint.peg`
3. Ajouter les tests d'intégration correspondants
4. Documenter la correspondance RETE

La grammaire est maintenant **complète, cohérente et testée à 100%** ! ✅
