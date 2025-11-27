# Corrections AlphaNode Sharing - Janvier 2025

## Vue d'ensemble

Cette correction résout deux problèmes identifiés lors de la session de debugging :
1. **Build failures** dans `rete/examples` (multiples fonctions `main`)
2. **Non-partage des AlphaNodes** entre règles simples et chaînes AND

## Problème #1 : Build Failures `rete/examples`

### Description
Le package `rete/examples` contenait plusieurs fichiers avec des fonctions `main()`, causant des erreurs de compilation lors de `go test ./...`.

### Solution
Ajout de la directive `//go:build ignore` en en-tête de chaque fichier d'exemple pour les exclure du build normal tout en permettant leur exécution individuelle via `go run`.

### Fichiers modifiés
- `rete/examples/alpha_chain_builder_example.go`
- `rete/examples/alpha_chain_extractor_example.go`
- `rete/examples/constraint_pipeline_chain_example.go`
- `rete/examples/expression_analyzer_example.go`

### Impact
✅ `go test ./...` passe maintenant sans erreur  
✅ Les exemples restent exécutables avec `go run`

---

## Problème #2 : Non-Partage AlphaNodes Règles Simples/Chaînes

### Description
Les règles simples et les chaînes ne partageaient pas les AlphaNodes pour des conditions identiques, créant des doublons inutiles.

**Exemple** :
```constraint
rule large : {t: Transaction} / t.amount > 1000 ==> print("LARGE")
rule fraud : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' ==> print("FRAUD")
```

**Avant** : 2 AlphaNodes distincts pour `t.amount > 1000`  
**Après** : 1 AlphaNode partagé

### Cause Racine
Les règles simples et les chaînes représentaient les conditions différemment :
- **Règles simples** : `{type: "constraint", constraint: {type: "comparison", ...}}`
- **Chaînes** : `{type: "binaryOperation", ...}`

Le hashing produisait des valeurs différentes, empêchant le partage.

### Solution
Implémentation d'une fonction `normalizeConditionForSharing()` dans `alpha_sharing.go` qui :

1. **Déballe les conditions wrappées**
   ```go
   // {type: "constraint", constraint: X} → X
   ```

2. **Normalise les types équivalents**
   ```go
   // "comparison" → "binaryOperation"
   ```

3. **Normalise récursivement** les structures imbriquées (maps, slices)

### Fichiers modifiés

#### `rete/alpha_sharing.go`
- Ajout de `normalizeConditionForSharing(condition interface{}) interface{}`
- Modification de `ConditionHash()` pour appeler la normalisation avant hashing

#### `rete/alpha_chain_integration_test.go`
- Mise à jour des assertions pour refléter le partage optimal
- `TestAlphaChain_ComplexScenario_FraudDetection` : 5 → 4 AlphaNodes
- `TestAlphaChain_PartialSharing_ThreeRules` : 4 → 3 AlphaNodes
- `TestAlphaChain_RuleRemoval_PreservesShared` : 3 → 2 AlphaNodes
- `TestAlphaChain_NetworkStats_Accurate` : 7 → 5 AlphaNodes
- `TestAlphaChain_MixedConditions_ComplexSharing` : 6 → 4 AlphaNodes

### Résultats Mesurés

| Test | Avant | Après | Amélioration |
|------|-------|-------|--------------|
| ComplexScenario_FraudDetection | 5 | 4 | -20% |
| PartialSharing_ThreeRules | 4 | 3 | -25% |
| RuleRemoval_PreservesShared | 3 | 2 | -33% |
| NetworkStats_Accurate | 7 | 5 | -28.6% |
| MixedConditions_ComplexSharing | 6 | 4 | -33.3% |

**Moyenne** : Réduction de 28% du nombre d'AlphaNodes

### Exemple Concret

**Scénario** : Détection de fraude bancaire
```constraint
type Transaction : <id: string, amount: number, country: string, risk: number>

rule fraud_low : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' ==> print("LOW")
rule fraud_med : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' AND t.risk > 50 ==> print("MED")
rule fraud_high : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' AND t.risk > 80 ==> print("HIGH")
rule large : {t: Transaction} / t.amount > 1000 ==> print("LARGE")
```

**Structure du réseau AVANT** :
```
TypeNode(Transaction)
├── AlphaNode(alpha_xxx: t.amount > 1000)  ← règle "large" (simple)
│   └── TerminalNode(large_terminal)
└── AlphaNode(alpha_yyy: t.amount > 1000)  ← règles "fraud_*" (chaînes) ❌ DOUBLON
    ├── AlphaNode(alpha_zzz: t.country == 'XX')
    │   ├── TerminalNode(fraud_low_terminal)
    │   └── AlphaNode(...)
    ...
```

**Structure du réseau APRÈS** :
```
TypeNode(Transaction)
└── AlphaNode(alpha_xxx: t.amount > 1000)  ← PARTAGÉ par toutes les règles ✅
    ├── TerminalNode(large_terminal)
    └── AlphaNode(alpha_zzz: t.country == 'XX')  ← PARTAGÉ par fraud_*
        ├── TerminalNode(fraud_low_terminal)
        └── AlphaNode(...)
        ...
```

**Résultat** : 5 AlphaNodes → 4 AlphaNodes (-20%)

### Impact

#### Performance
- ✅ **Moins d'évaluations** : Chaque condition unique n'est évaluée qu'une seule fois
- ✅ **Propagation optimale** : Un fait traverse moins de nœuds
- ✅ **Scalabilité** : Amélioration linéaire avec le nombre de règles

#### Mémoire
- ✅ **Réduction 20-33%** du nombre d'AlphaNodes
- ✅ **Moins de structures** : Maps, mutex, mémoires de travail

#### Conformité
- ✅ **Algorithme RETE classique** : Le partage de nœuds est une optimisation standard
- ✅ **Transparence** : Aucun changement d'API, correction invisible pour l'utilisateur

### Tests de Validation

Tous les tests passent avec succès :
```bash
$ go test ./rete
ok  	github.com/treivax/tsd/rete	0.112s

$ go test ./...
ok  	github.com/treivax/tsd/cmd/tsd	0.603s
ok  	github.com/treivax/tsd/constraint	(cached)
ok  	github.com/treivax/tsd/rete	(cached)
[...tous les packages passent...]
```

Tests spécifiques au partage :
- ✅ `TestAlphaSharingIntegration_*` (5 tests)
- ✅ `TestAlphaChain_*` (10 tests)
- ✅ `TestConditionHash` (normalisation)
- ✅ `TestTypeNodeSharing_*` (3 tests)

---

## Code Ajouté

### `normalizeConditionForSharing()` dans `alpha_sharing.go`

```go
// normalizeConditionForSharing déballe les conditions wrappées pour permettre le partage
// entre règles simples (qui wrappent dans {"type": "constraint", "constraint": X})
// et chaînes (qui utilisent directement la condition décomposée)
func normalizeConditionForSharing(condition interface{}) interface{} {
	// Si la condition est une map
	if condMap, ok := condition.(map[string]interface{}); ok {
		// Vérifier si c'est une condition wrappée dans un type "constraint"
		if condType, hasType := condMap["type"]; hasType {
			if condTypeStr, ok := condType.(string); ok && condTypeStr == "constraint" {
				// Déballer la condition interne
				if innerCond, hasConstraint := condMap["constraint"]; hasConstraint {
					// Récursion pour déballer plusieurs niveaux si nécessaire
					return normalizeConditionForSharing(innerCond)
				}
			}
		}

		// Normaliser les types équivalents pour le partage
		// "comparison" et "binaryOperation" sont des synonymes
		normalized := make(map[string]interface{})
		for key, value := range condMap {
			if key == "type" {
				if typeStr, ok := value.(string); ok {
					// Normaliser "comparison" vers "binaryOperation"
					if typeStr == "comparison" {
						normalized[key] = "binaryOperation"
					} else {
						normalized[key] = value
					}
				} else {
					normalized[key] = value
				}
			} else {
				// Normaliser récursivement les valeurs imbriquées
				normalized[key] = normalizeConditionForSharing(value)
			}
		}
		return normalized
	}

	// Si c'est un slice, normaliser chaque élément
	if slice, ok := condition.([]interface{}); ok {
		normalized := make([]interface{}, len(slice))
		for i, item := range slice {
			normalized[i] = normalizeConditionForSharing(item)
		}
		return normalized
	}

	// Sinon, retourner la condition telle quelle
	return condition
}
```

### Modification de `ConditionHash()`

```go
func ConditionHash(condition interface{}, variableName string) (string, error) {
	// Déballer la condition si elle est wrappée (pour le partage entre règles simples et chaînes)
	unwrapped := normalizeConditionForSharing(condition)

	// Normaliser la condition pour assurer un hash cohérent
	normalized, err := normalizeCondition(unwrapped)
	if err != nil {
		return "", fmt.Errorf("erreur normalisation condition: %w", err)
	}

	// ... suite du code (calcul du hash)
}
```

---

## Vérification

### Avant correction
```bash
$ go test ./rete/examples/...
# github.com/treivax/tsd/rete/examples
rete/examples/alpha_chain_extractor_example.go:16:6: main redeclared in this block
FAIL	github.com/treivax/tsd/rete/examples [build failed]
```

### Après correction
```bash
$ go test ./rete/examples/...
?   	github.com/treivax/tsd/rete/examples/normalization	[no test files]

$ go run rete/examples/alpha_chain_builder_example.go
=== Alpha Chain Builder - Exemple d'utilisation ===
[...succès...]
```

### Partage AlphaNodes - Logs

**Avant** (règle `large`) :
```
✨ Nouveau AlphaNode partageable créé: alpha_22024b423dba910f (hash: alpha_22024b423dba910f)
```

**Après** (règle `large`) :
```
♻️  AlphaNode partagé réutilisé: alpha_e554bda722b2b37a (hash: alpha_e554bda722b2b37a)
✓ Règle large attachée à l'AlphaNode partagé alpha_e554bda722b2b37a via terminal large_terminal
```

✅ **Le même hash `alpha_e554bda722b2b37a` est utilisé par les chaînes `fraud_*` !**

---

## Recommandations Futures

### Court terme
- ✅ Commit et push des changements
- ✅ Mise à jour de `ALPHA_NODE_SHARING.md`
- ⚠️ Ajouter test unitaire spécifique pour `normalizeConditionForSharing()`

### Moyen terme
- 🔄 Considérer l'extension du partage aux BetaNodes (jointures)
- 🔄 Ajouter métriques de monitoring sur le taux de partage effectif
- 🔄 Documenter les équivalences de types (`comparison` ↔ `binaryOperation`)

### Long terme
- 🔄 Condition subsumption (ex: `age > 18` subsume `age > 21`)
- 🔄 Partage de sous-expressions logiques (ex: `(A AND B) OR (A AND C)` → partage de `A`)

---

## Références

- Issue/Thread: `3480a406-cabf-4f8e-8645-791e2ba5dad4`
- Documentation: `ALPHA_NODE_SHARING.md`
- Tests: `alpha_chain_integration_test.go`, `alpha_sharing_test.go`
- Algorithme: Forgy, C. L. (1982). "Rete: A Fast Algorithm for the Many Pattern/Many Object Pattern Match Problem"

---

**Date de correction** : 2025-01-XX  
**Auteur** : Session de debugging + corrections automatiques  
**Statut** : ✅ COMPLÉTÉ ET VALIDÉ  
**Version** : rete v1.3.0 (après merge)