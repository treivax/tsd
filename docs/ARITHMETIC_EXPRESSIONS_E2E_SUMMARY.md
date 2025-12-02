# Support des Expressions Arithmétiques dans les Actions - Résumé

## Vue d'ensemble

Ce document résume l'implémentation complète du support des expressions arithmétiques dans les actions du système RETE, incluant le test end-to-end (E2E) qui valide l'ensemble du pipeline.

## Fonctionnalités implémentées

### 1. Expressions arithmétiques dans les actions

Les actions peuvent maintenant contenir des expressions arithmétiques complexes utilisant:
- **Opérateurs arithmétiques**: `+`, `-`, `*`, `/`, `%`
- **Expressions imbriquées**: `(a + b) * (c - d)`
- **Accès aux champs**: `variable.champ`
- **Littéraux numériques**: entiers et décimaux

### 2. Exemple de syntaxe

```tsd
type Produit(id: string, prix: number, quantite: number, poids: number)
type Commande(id: string, produit_id: string, qte: number, remise: number)

action facture_calculee(
    commande_id: string,
    total_brut: number,
    montant_remise: number,
    total_net: number,
    prix_par_unite: number,
    cout_livraison: number
)

rule calcul_facture_complete : {p: Produit, c: Commande} /
    c.produit_id == p.id AND c.qte > 0
    ==> facture_calculee(
        c.id,
        p.prix * c.qte,
        (p.prix * c.qte) * (c.remise / 100),
        (p.prix * c.qte) - ((p.prix * c.qte) * (c.remise / 100)),
        p.prix,
        (p.poids * c.qte) + 10
    )

Produit(id: "PROD001", prix: 100, quantite: 50, poids: 2)
Commande(id: "CMD001", produit_id: "PROD001", qte: 5, remise: 10)
```

## Corrections et améliorations apportées

### 1. Décodage des opérateurs Base64

**Problème**: Le parser encode les opérateurs en base64 (ex: `+` devient `Kw==`)

**Solution**: Décodage automatique dans:
- `constraint/action_validator.go` - `inferArgumentType()`
- `rete/action_executor.go` - `evaluateBinaryOperation()`

```go
// The operator might be base64 encoded, try to decode it
if decoded, err := base64.StdEncoding.DecodeString(operator); err == nil {
    operator = string(decoded)
}
```

### 2. Décodage des opérateurs Base64 dans l'évaluateur

**Problème**: L'AlphaConditionEvaluator utilisé pour évaluer les conditions dans les prémisses ne décodait pas les opérateurs base64, causant le rejet ou l'acceptation incorrecte des conditions arithmétiques complexes.

**Solution**: Ajout du décodage dans `evaluator_values.go` et `evaluator_expressions.go`:

```go
// Dans evaluator_values.go - cas binaryOp
// The operator might be base64 encoded, try to decode it
if decoded, err := base64.StdEncoding.DecodeString(operator); err == nil {
    operator = string(decoded)
}

// Dans evaluator_expressions.go - evaluateBinaryOperationMap
// The operator might be base64 encoded, try to decode it
if decoded, err := base64.StdEncoding.DecodeString(operator); err == nil {
    operator = string(decoded)
}
```

### 3. Support des variantes de type dans l'inférence

**Problème**: Le parser génère différentes variantes: `binaryOp`, `binaryOperation`, `binary_operation`

**Solution**: Support de toutes les variantes dans `action_validator.go`:

```go
case "binaryOp", "binaryOperation", "binary_operation":
    // For binary operations, infer from operands
    op, ok := v["operator"].(string)
    if !ok {
        return "", fmt.Errorf("binaryOp missing operator")
    }
    // Arithmetic operations return number
    if op == "+" || op == "-" || op == "*" || op == "/" || op == "%" {
        return "number", nil
    }
```

### 4. Gestion du champ `id` dans les faits

**Problème**: Le champ `id` n'était pas copié dans `fact.Fields`, donc inaccessible dans les conditions de jointure

**Solution**: Modification de `SubmitFactsFromGrammar()` dans `network.go`:

```go
// Copier tous les champs (y compris "id")
for key, value := range factMap {
    if key != "type" && key != "reteType" {
        fact.Fields[key] = value
    }
}
```

### 5. Support du champ `reteType`

**Problème**: `ConvertFactsToReteFormat` utilise `"reteType"` mais `SubmitFactsFromGrammar` cherchait seulement `"type"`

**Solution**: Recherche des deux variantes:

```go
factType := "unknown"
// Chercher "type" ou "reteType"
if typ, ok := factMap["type"].(string); ok {
    factType = typ
} else if typ, ok := factMap["reteType"].(string); ok {
    factType = typ
}
```

### 6. Vérification de l'existence des champs dans les jointures

**Problème**: Les conditions de jointure pouvaient échouer silencieusement si un champ n'existait pas

**Solution**: Vérification explicite dans `evaluateSimpleJoinConditions()`:

```go
leftValue, leftExists := leftFact.Fields[joinCondition.LeftField]
rightValue, rightExists := rightFact.Fields[joinCondition.RightField]

// Vérifier que les champs existent
if !leftExists || !rightExists {
    return false
}
```

## Structure du test E2E

### Fichiers créés

1. **`rete/testdata/arithmetic_e2e.tsd`**
   - Définitions de types (Produit, Commande, Client)
   - Définition d'action avec expressions complexes
   - Règle avec conditions de jointure
   - Faits de test (3 produits, 3 commandes, 2 clients)

2. **`rete/action_arithmetic_e2e_test.go`**
   - Test complet du pipeline unique
   - Vérification de la construction du réseau
   - Validation de l'exécution des actions
   - Affichage des résultats détaillés

### Déroulement du test

1. **Parsing**: Lecture du fichier `.tsd` avec types, règles et faits
2. **Construction**: Création du réseau RETE avec TypeNodes, AlphaNodes, JoinNodes, TerminalNodes
3. **Injection**: Soumission des faits au réseau
4. **Propagation**: 
   - Faits → TypeNodes → AlphaNodes (passthrough) → JoinNodes
   - Évaluation des conditions de jointure
   - Création des tokens résultants
5. **Exécution**: Activation des TerminalNodes avec évaluation des expressions arithmétiques
6. **Validation**: Vérification que 3 tokens sont générés et 3 actions exécutées

### Résultat du test

```
✅ Total de tokens générés: 3
✅ Actions exécutées: 3

📦 Token #1: CMD001 × PROD001
   - total_brut: 100 * 5 = 500
   - montant_remise: 500 * 0.10 = 50
   - total_net: 500 - 50 = 450
   - cout_livraison: (2 * 5) + 10 = 20

📦 Token #2: CMD002 × PROD002
   - total_brut: 250 * 2 = 500
   - montant_remise: 500 * 0.15 = 75
   - total_net: 500 - 75 = 425
   - cout_livraison: (5 * 2) + 10 = 20

📦 Token #3: CMD003 × PROD003
   - total_brut: 50 * 10 = 500
   - montant_remise: 500 * 0.05 = 25
   - total_net: 500 - 25 = 475
   - cout_livraison: (1 * 10) + 10 = 20
```

## Composants modifiés

### Fichiers du moteur RETE

1. **`rete/action_executor.go`**
   - Ajout du décodage base64 des opérateurs pour les actions
   - Support des variantes de type

2. **`rete/evaluator_values.go`**
   - **CRUCIAL**: Ajout du décodage base64 des opérateurs dans `evaluateValueFromMap()` pour le cas `binaryOp`
   - Permet l'évaluation correcte des expressions arithmétiques dans les prémisses

3. **`rete/evaluator_expressions.go`**
   - **CRUCIAL**: Ajout du décodage base64 dans `evaluateBinaryOperationMap()`
   - Assure l'évaluation correcte des comparaisons avec expressions arithmétiques

4. **`rete/network.go`**
   - Modification de `SubmitFactsFromGrammar()` pour gérer `reteType` et inclure `id`

5. **`rete/node_join.go`**
   - Ajout de la vérification d'existence des champs dans les jointures
   - Évaluation correcte des conditions alpha avec expressions arithmétiques

### Fichiers de validation

6. **`constraint/action_validator.go`**
   - Ajout du décodage base64 des opérateurs pour la validation de types
   - Support de `binaryOp`, `binaryOperation`, `binary_operation`
   - Support de l'opérateur modulo `%`

## Architecture du pipeline

```
Fichier .tsd
    ↓
Parser (IterativeParser)
    ↓
Programme AST
    ↓
ConvertToReteProgram
    ↓
Réseau RETE:
    RootNode
      ├─→ TypeNode[Produit] → AlphaNode[passthrough,left] ─┐
      └─→ TypeNode[Commande] → AlphaNode[passthrough,right] ┴→ JoinNode
                                                                  ↓
                                                              TerminalNode
                                                                  ↓
                                                              ActionExecutor
                                                                  ↓
                                                    Évaluation des expressions
```

## Tests disponibles

### Tests unitaires
- `rete/action_arithmetic_test.go` - Tests unitaires des expressions arithmétiques
- `rete/action_arithmetic_complex_test.go` - Tests d'expressions complexes

### Tests d'intégration
- `rete/evaluator_complex_expressions_test.go` - Tests des expressions dans les contraintes

### Test E2E
- `rete/action_arithmetic_e2e_test.go` - Test du pipeline complet

### Commande de test

```bash
# Test E2E uniquement
go test -v -run TestArithmeticExpressionsE2E ./rete

# Tous les tests arithmétiques
go test -v -run Arithmetic ./rete

# Tous les tests
go test -v ./rete
```

## Documentation associée

- **`docs/ARITHMETIC_IN_ACTIONS.md`** - Guide d'utilisation détaillé
- **`docs/ACTIONS_SYSTEM.md`** - Documentation du système d'actions (section Expressions Arithmétiques)
- **`rete/examples/arithmetic_actions_example.go`** - Exemple exécutable

## Points clés de l'implémentation

### ✅ Avantages

1. **Transparence**: Le décodage base64 est automatique et transparent pour l'utilisateur
2. **Compatibilité**: Support de toutes les variantes de types générées par le parser
3. **Robustesse**: Vérification de l'existence des champs avant évaluation
4. **Testabilité**: Test E2E complet avec fichier `.tsd` réaliste
5. **Maintenabilité**: Code bien structuré avec séparation des responsabilités

### ⚠️ Points d'attention

1. **Le parser encode les opérateurs en base64** - nécessite un décodage dans TROIS endroits:
   - `action_executor.go` - pour l'évaluation des actions
   - `evaluator_values.go` - pour l'évaluation des valeurs dans les prémisses (expressions arithmétiques)
   - `evaluator_expressions.go` - pour l'évaluation des opérations binaires dans les conditions
   
2. Les faits doivent avoir leur champ `id` dans `Fields` pour les jointures

3. Le champ de type peut être `"type"` ou `"reteType"` selon le contexte

4. Les passthrough AlphaNodes avec `side` nécessitent une propagation appropriée (LEFT vs RIGHT)

5. **Les conditions arithmétiques dans les prémisses sont évaluées au niveau du JoinNode**, pas au niveau des AlphaNodes

## Validation avec condition inversée

Un test supplémentaire a été effectué en inversant la condition `c.qte * 3 * 2 + 1 > 0` en `c.qte * 3 * 2 + 1 < 0`:

### Résultat attendu
- Toutes les commandes ont `qte > 0`, donc l'expression `c.qte * 3 * 2 + 1` est toujours positive
- La condition `< 0` devrait être fausse pour toutes les commandes
- **Aucun token ne devrait être généré**

### Résultat obtenu
✅ **0 tokens générés, 0 actions exécutées** - Comportement correct !

Le système évalue correctement les expressions arithmétiques complexes dans les prémisses et rejette les jointures lorsque les conditions ne sont pas satisfaites.

## Conclusion

L'implémentation du support des expressions arithmétiques est complète et validée end-to-end. Le système peut maintenant:

- Parser des expressions arithmétiques complexes dans les **actions** ET les **prémisses**
- Valider les types des expressions
- **Évaluer correctement les expressions avec scalaires dans les conditions de jointure**
- Gérer les expressions imbriquées avec plusieurs opérateurs
- Supporter l'ensemble des opérateurs arithmétiques standard
- **Filtrer correctement les tokens selon les conditions arithmétiques**

Le test E2E avec conditions inversées démontre que l'ensemble du pipeline fonctionne correctement, incluant l'évaluation des expressions arithmétiques dans les prémisses pour filtrer les faits qui ne satisfont pas les conditions.