# Rapport d'implémentation : Règles sans condition

**Date** : 2025-01-15  
**Auteur** : Assistant IA  
**Type** : Nouvelle fonctionnalité  
**Module** : `constraint` (parser, grammaire PEG)

---

## 📋 Résumé exécutif

Implémentation de la fonctionnalité permettant d'écrire des règles TSD sans condition de filtrage. Les règles se déclenchent automatiquement dès qu'un fait correspondant au pattern est asserté, sans nécessiter de contrainte explicite.

**Syntaxe** :
```tsd
rule nom_regle : {variable: Type} / ==> action(...)
```

**Exemple** :
```tsd
rule assertion_user : {p: Person} / ==> log("nouvel utilisateur : " + p.name)
```

---

## 🎯 Objectif

Permettre l'écriture de règles sans condition pour :
- Logger automatiquement tous les événements d'un certain type
- Auditer toutes les opérations
- Déclencher des webhooks/notifications systématiques
- Implémenter des patterns Event Sourcing / CQRS
- Simplifier les règles qui s'appliquent à tous les faits d'un type

---

## 🔧 Modifications apportées

### 1. Grammaire PEG (`constraint/grammar/constraint.peg`)

**Modification de la règle `Expression`** :

Ajout d'une alternative pour supporter les règles sans contrainte :

```peg
Expression <- "rule" _ ruleId:IdentName _ ":" _ patterns:PatternBlocks _ "/" _ constraints:Constraints _ "==>" _ action:Action {
    // ... code existant (règle avec contrainte)
} / "rule" _ ruleId:IdentName _ ":" _ patterns:PatternBlocks _ "/" _ "==>" _ action:Action {
    // Règle sans contrainte - nouveau cas
    patternList := patterns.([]interface{})
    if len(patternList) == 1 {
        return map[string]interface{}{
            "type": "expression",
            "ruleId": ruleId,
            "set": patternList[0],
            "constraints": nil,  // ← contrainte nil
            "action": action,
        }, nil
    }

    return map[string]interface{}{
        "type": "expression",
        "ruleId": ruleId,
        "patterns": patterns,
        "constraints": nil,  // ← contrainte nil
        "action": action,
    }, nil
}
```

**Principe** :
- La grammaire PEG essaie d'abord de parser une règle avec contrainte
- Si ça échoue, elle essaie de parser une règle sans contrainte (juste `/ ==>`)
- Dans le second cas, le champ `constraints` est mis à `nil`

### 2. Parser généré (`constraint/parser.go`)

Régénération automatique via :
```bash
cd constraint
pigeon -o parser.go grammar/constraint.peg
```

Le parser Go généré contient maintenant la logique pour parser les deux formes de règles.

### 3. Tests (`constraint/no_condition_rules_test.go`)

**Nouveau fichier de tests** avec couverture complète :

- `TestParser_NoConditionRules` : 10 scénarios de parsing
  - Règle simple sans condition
  - Règles multiples sans conditions
  - Règle avec actions multiples
  - Accès à plusieurs champs
  - Règles mixtes (avec/sans condition)
  - Expressions arithmétiques dans actions
  - Concaténation de strings
  - Multi-patterns sans condition
  - Champs booléens et numériques

- `TestParser_NoConditionRulesValidation` : Vérification de la structure AST
  - Vérifie que `constraints` est `nil` pour les règles sans condition

- `TestParser_NoConditionRulesWithSpaces` : Tolérance aux espaces
  - Espaces minimaux, multiples, newlines, tabs

- `TestParser_NoConditionRulesErrorCases` : Cas d'erreur
  - Action manquante
  - Flèche `==>` manquante
  - Pattern manquant

- `TestNoConditionRulesIntegration` : Test d'intégration
  - Parsing complet avec types, actions, règles et faits

- `TestNoConditionRulesWithComplexActions` : Actions complexes
  - Multi-actions avec expressions arithmétiques

- `TestNoConditionRulesWithAggregation` : Multi-patterns
  - Règles d'agrégation sans conditions explicites

**Résultats** : ✅ Tous les tests passent

### 4. Exemple (`examples/no_condition_rules.tsd`)

**Fichier d'exemple complet** avec :
- 4 types (Person, Product, Order, Event)
- 5 actions (log, notify, track, audit, webhook)
- 13 règles dont 8 sans condition
- 10 faits de test
- Documentation inline complète

**Points clés de l'exemple** :
- Règles sans condition pour logging, audit, tracking
- Règles mixtes (avec/sans condition) pour comparaison
- Multi-actions en séquence
- Démonstration des cas d'usage réels

### 5. Documentation (`docs/no-condition-rules.md`)

**Documentation complète** (357 lignes) couvrant :

1. **Vue d'ensemble** : Concept et cas d'usage
2. **Syntaxe** : Forme générale et comparaison avec règles classiques
3. **Exemples** : 5 exemples progressifs
4. **Sémantique** : Activation, pattern matching, contraintes implicites
5. **Cas d'usage** : Event Sourcing, Monitoring, CQRS, Notifications
6. **Bonnes pratiques** : À faire, à éviter
7. **Performance** : Considérations et optimisations
8. **Architecture interne** : AST, compilation RETE, évaluation
9. **Limitations** : Contraintes et solutions
10. **Migration** : Depuis règles classiques

---

## 🧪 Validation

### Tests unitaires

```bash
cd constraint
go test -run TestParser_NoConditionRules -v
```

**Résultat** : ✅ PASS (tous les sous-tests passent)

### Tests d'intégration

```bash
cd constraint
go test -v
```

**Résultat** : ✅ PASS (0.163s, aucun test cassé)

### Validation de l'exemple

```bash
cd constraint
go run cmd/main.go ../examples/no_condition_rules.tsd
```

**Résultat** : ✅ Programme valide avec 4 types, 13 expressions, 10 faits

### Compilation

```bash
go build ./constraint/...
```

**Résultat** : ✅ Aucune erreur de compilation

### Formatage et linting

```bash
go fmt ./constraint/...
goimports -w ./constraint
go vet ./constraint/...
```

**Résultat** : ✅ Aucun problème détecté

---

## 📊 Statistiques

| Métrique | Valeur |
|----------|--------|
| Fichiers modifiés | 4 |
| Fichiers créés | 3 |
| Lignes de code ajoutées | ~1000 |
| Tests ajoutés | 7 fonctions de test |
| Scénarios de test | 30+ |
| Couverture documentation | 357 lignes |
| Exemples | 13 règles, 10 faits |

---

## 🔍 Détails techniques

### Représentation AST

Les règles sans condition dans le JSON AST :

```json
{
  "type": "expression",
  "ruleId": "log_all",
  "set": {
    "type": "set",
    "variables": [{"type": "typedVariable", "name": "p", "dataType": "Person"}]
  },
  "constraints": null,  // ← Champ clé : null pour règles sans condition
  "action": {
    "type": "sequenceAction",
    "jobs": [...]
  }
}
```

### Compilation RETE

Les règles sans condition dans le réseau RETE :

1. **AlphaNode** : Test de type du fait
2. **Pas de BetaNode** : Aucune contrainte à évaluer
3. **TerminalNode** : Action à exécuter directement
4. **Activation** : Une activation par fait correspondant

**Flux** :
```
Fait → AlphaNode (type check) → TerminalNode → Action exécutée
```

### Rétrocompatibilité

✅ **100% rétrocompatible**

- Les règles existantes avec contraintes fonctionnent exactement comme avant
- Aucun changement dans la structure AST pour les règles avec contraintes
- Le parser essaie d'abord la forme avec contrainte (comportement par défaut)
- La nouvelle forme est une alternative PEG qui n'affecte pas l'existant

---

## 📝 Exemples d'utilisation

### Cas 1 : Logging automatique

```tsd
type User(#userId: string, name: string)
action log(message: string)

rule log_users : {u: User} / ==> log(u.name)

User(userId: "1", name: "Alice")  // ← Déclenche automatiquement log("Alice")
```

### Cas 2 : Audit de sécurité

```tsd
type SecurityEvent(#eventId: string, type: string, userId: string)
action audit(eventId: string, type: string)

rule audit_security : {e: SecurityEvent} / ==>
    audit(e.eventId, e.type)

SecurityEvent(eventId: "e1", type: "login_failed", userId: "u123")
// ← Automatiquement audité
```

### Cas 3 : Webhook systématique

```tsd
type Order(#orderId: string, customerId: string, amount: number)
action webhook(url: string, orderId: string)

rule notify_orders : {o: Order} / ==>
    webhook("https://api.example.com/orders", o.orderId)

Order(orderId: "o1", customerId: "c1", amount: 100)
// ← Webhook automatiquement envoyé
```

---

## ⚠️ Points d'attention

### 1. Performance sur gros volumes

Les règles sans condition créent une activation pour **chaque** fait du type correspondant.

**Exemple** :
- 1 million de faits `Person`
- 1 règle sans condition sur `Person`
- = 1 million d'activations

**Solution** : Utiliser des règles avec conditions pour filtrer en amont.

### 2. Multi-patterns sans conditions

Attention au produit cartésien :

```tsd
// ⚠️ DANGER : n × m activations !
rule dangerous : {o: Order} / {c: Customer} / ==> match(o.orderId, c.customerId)
```

**Solution** : Ajouter une condition de jointure :

```tsd
// ✅ Mieux : filtrage sur jointure
rule safe : {o: Order} / {c: Customer} / o.customerId == c.customerId ==>
    match(o.orderId, c.customerId)
```

### 3. Boucles infinies

Ne pas créer de faits du même type dans l'action :

```tsd
// ⚠️ BOUCLE INFINIE !
rule loop : {p: Person} / ==> create_person(p.name)
```

---

## 🚀 Prochaines étapes

### Court terme
- [ ] Tester avec des cas d'usage réels (production)
- [ ] Monitorer les performances sur gros volumes
- [ ] Documenter les métriques RETE pour règles sans condition

### Moyen terme
- [ ] Optimisations spécifiques pour règles sans condition
- [ ] Batching automatique des activations
- [ ] Métriques de performance dédiées

### Long terme
- [ ] Support de patterns plus complexes
- [ ] Mode "sampling" (ne déclencher que 10% des faits, par exemple)
- [ ] Intégration avec des systèmes de streaming (Kafka, etc.)

---

## 📚 Références

- **Grammaire PEG** : `constraint/grammar/constraint.peg` (lignes 148-192)
- **Tests** : `constraint/no_condition_rules_test.go`
- **Documentation** : `docs/no-condition-rules.md`
- **Exemple** : `examples/no_condition_rules.tsd`
- **Standard développement** : `.github/prompts/develop.md`

---

## ✅ Checklist finale

- [x] En-tête copyright présent dans tous les fichiers
- [x] Aucun hardcoding
- [x] Code générique avec paramètres
- [x] Variables/fonctions privées par défaut
- [x] `go fmt` + `goimports` appliqués
- [x] `go vet` sans erreur
- [x] Tests écrits et passent
- [x] Documentation complète
- [x] Exemple fonctionnel
- [x] Rétrocompatibilité assurée

---

## 📧 Contact

Pour questions ou suggestions : voir documentation projet TSD.

---

**FIN DU RAPPORT**