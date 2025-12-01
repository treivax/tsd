# Phase 9: Intégration des Builders

## Objectif
Modifier `constraint_pipeline_builder.go` pour utiliser les builders du package `rete/builders` et réduire le fichier à ~200 lignes.

## État actuel
- Fichier actuel: `rete/constraint_pipeline_builder.go` (~1030 lignes)
- Builders créés dans: `rete/builders/` (7 fichiers, ~1437 lignes)
- Tous les builders compilent et sont fonctionnels

## Problème identifié: Import cyclique
- Package `rete` ne peut pas importer `rete/builders`
- Package `rete/builders` importe `rete` pour utiliser les types (ReteNetwork, Node, etc.)
- Solution: Garder les builders dans le package `rete` mais dans des fichiers séparés

## Stratégie d'intégration

### Option A: Déplacer builders dans package rete (RECOMMANDÉ)
1. Renommer `rete/builders/*.go` avec préfixe `builder_`
2. Les déplacer dans `rete/`
3. Changer `package builders` → `package rete`
4. Supprimer imports `github.com/treivax/tsd/rete`
5. Utiliser directement les types sans préfixe

Avantages:
- Pas de cycle d'import
- Même package, juste fichiers séparés
- Compilation simple

### Option B: Package séparé avec interface
Créer un package intermédiaire pour les types partagés.

Désavantage: Trop complexe pour le bénéfice

## Plan d'exécution (Option A)

### Étape 1: Réorganiser les fichiers
```bash
cd rete/builders
for f in *.go; do
  mv "$f" "../builder_$f"
done
cd ..
rmdir builders
```

### Étape 2: Adapter le package
```bash
cd rete
for f in builder_*.go; do
  # Changer le package
  sed -i 's/^package builders$/package rete/' "$f"
  
  # Supprimer les imports rete
  sed -i '/github.com\/treivax\/tsd\/rete/d' "$f"
  
  # Enlever les préfixes rete.
  sed -i 's/rete\.//g' "$f"
done
```

### Étape 3: Vérifier les constantes en double
Les constantes suivantes sont définies dans `builder_utils.go`:
- ConditionTypePassthrough, ConditionTypeSimple, ConditionTypeExists, ConditionTypeComparison
- NodeSideLeft, NodeSideRight

Vérifier si elles existent déjà dans `constraint_pipeline_builder.go` et supprimer les doublons.

### Étape 4: Simplifier constraint_pipeline_builder.go

Le fichier doit simplement instancier et déléguer:

```go
// buildNetwork construit le réseau RETE
func (cp *ConstraintPipeline) buildNetwork(storage Storage, types []interface{}, expressions []interface{}) (*ReteNetwork, error) {
    network := NewReteNetwork(storage)
    
    err := cp.createTypeNodes(network, types, storage)
    if err != nil {
        return nil, fmt.Errorf("erreur création TypeNodes: %w", err)
    }
    
    err = cp.createRuleNodes(network, expressions, storage)
    if err != nil {
        return nil, fmt.Errorf("erreur création règles: %w", err)
    }
    
    return network, nil
}

// createTypeNodes délègue au TypeBuilder
func (cp *ConstraintPipeline) createTypeNodes(network *ReteNetwork, types []interface{}, storage Storage) error {
    utils := NewBuilderUtils(storage)
    typeBuilder := NewTypeBuilder(utils)
    return typeBuilder.CreateTypeNodes(network, types, storage)
}

// createRuleNodes délègue au RuleBuilder
func (cp *ConstraintPipeline) createRuleNodes(network *ReteNetwork, expressions []interface{}, storage Storage) error {
    utils := NewBuilderUtils(storage)
    alphaBuilder := NewAlphaRuleBuilder(utils)
    existsBuilder := NewExistsRuleBuilder(utils)
    joinBuilder := NewJoinRuleBuilder(utils)
    accumulatorBuilder := NewAccumulatorRuleBuilder(utils)
    ruleBuilder := NewRuleBuilder(utils, alphaBuilder, existsBuilder, joinBuilder, accumulatorBuilder)
    
    return ruleBuilder.CreateRuleNodes(network, expressions)
}

// Les anciennes méthodes deviennent des wrappers simples
```

### Étape 5: Supprimer le code redondant

Toutes les fonctions suivantes peuvent être supprimées car déléguées:
- `createTypeDefinition` → dans `builder_types.go`
- `createAlphaRule` → dans `builder_alpha_rules.go`
- `createJoinRule` → dans `builder_join_rules.go`
- `createExistsRule` → dans `builder_exists_rules.go`
- `createAccumulatorRule` → dans `builder_accumulator_rules.go`
- `createMultiSourceAccumulatorRule` → dans `builder_accumulator_rules.go`
- `createBinaryJoinRule` → dans `builder_join_rules.go`
- `createCascadeJoinRule` → dans `builder_join_rules.go`
- Toutes les fonctions helper (extract*, is*, connect*, create*)

### Étape 6: Résultat attendu

Fichier `constraint_pipeline_builder.go` final (~150-200 lignes):
- Constants (si nécessaires)
- `buildNetwork()` - orchestration principale
- `createTypeNodes()` - délégation TypeBuilder
- `createRuleNodes()` - délégation RuleBuilder  
- `createSingleRule()` - délégation RuleBuilder
- Quelques wrappers legacy pour compatibilité

## Validation

### Tests de compilation
```bash
go build ./rete/...
```

### Tests unitaires
```bash
go test ./rete/... -v
```

### Vérification de la réduction
```bash
wc -l rete/constraint_pipeline_builder.go
# Doit être < 250 lignes
```

## Commit
```bash
git add rete/
git commit -m "refactor(rete): Phase 9 - Integration builders, reduce constraint_pipeline_builder to ~200 lines"
```

## Risques et mitigations

### Risque 1: Signatures incompatibles
**Mitigation**: Adapter les signatures des builders pour matcher exactement les appels existants

### Risque 2: Tests qui cassent
**Mitigation**: Exécuter tous les tests après chaque étape et corriger immédiatement

### Risque 3: Régression de performance
**Mitigation**: Exécuter benchmarks avant/après (Phase 10)

## Prochaine étape: Phase 10 - Tests
Une fois l'intégration terminée et compilant, passer à la Phase 10 pour:
- Tests unitaires de chaque builder
- Validation des tests existants
- Benchmarks de performance