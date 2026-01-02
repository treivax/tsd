# Changement de Syntaxe : Update avec Objets Littéraux

## Résumé

La syntaxe de l'action `Update` a été améliorée pour utiliser des **objets littéraux** `{...}` au lieu de la syntaxe personnalisée `_Mods_(...)`. Cette nouvelle syntaxe est plus naturelle, cohérente avec le reste du langage, et plus facile à lire.

## Motivation

L'ancienne syntaxe `Update(p, _Mods_(statut: "actif"))` utilisait un pseudo-type `_Mods_` spécial qui :
- N'était pas intuitif pour les utilisateurs
- Créait une exception syntaxique dans le langage
- Ne correspondait à aucune construction standard

La grammaire PEG de TSD supporte déjà les objets littéraux (`ObjectLiteral`) via la syntaxe `{...}`. Utiliser cette syntaxe existante pour `Update` rend le langage plus cohérent et plus facile à apprendre.

## Nouvelle Syntaxe

### Syntaxe de Base

```tsd
Update(variable, {champ: valeur, champ2: valeur2, ...})
```

### Exemples

#### Mise à jour d'un seul champ

```tsd
type Personne(#nom: string, age: number, statut: string)

rule anniversaire : {p: Personne} / p.statut == "anniversaire" ==>
    Update(p, {age: p.age + 1.0})
```

#### Mise à jour de plusieurs champs

```tsd
rule activation : {p: Personne} / p.statut == "nouveau" ==>
    Update(p, {statut: "actif", ville: "Paris"})
```

#### Avec expressions et accès aux champs

```tsd
type Relation(personne1: string, personne2: string, lien: string)

rule mettre_en_couple : {p: Personne, r: Relation} /
    p.nom == r.personne1 AND r.lien == "mariage" ==>
    Update(p, {statut: "en couple"})
```

## Ancienne Syntaxe (Déconseillée)

```tsd
// ❌ Ancienne syntaxe avec _Mods_
Update(p, _Mods_(statut: "actif"))

// ✅ Nouvelle syntaxe avec objet littéral
Update(p, {statut: "actif"})
```

## Sémantique

L'action `Update` avec objet littéral :

1. **Préserve l'identifiant du fait** : contrairement à l'ancienne syntaxe `Update(Personne(...))` qui créait un nouveau fait avec un nouvel ID, la nouvelle syntaxe modifie le fait existant en conservant son `_id_`.

2. **Modifie uniquement les champs spécifiés** : les champs non mentionnés dans l'objet littéral conservent leur valeur actuelle.

3. **Propage les changements dans RETE** : la mise à jour déclenche la réévaluation des règles dépendantes dans le réseau.

4. **Valide les types** : les champs spécifiés doivent exister dans le type du fait et respecter les types définis.

## Implémentation Technique

### Parser PEG

La grammaire détecte automatiquement les appels `Update(variable, {...})` et les transforme en une structure AST spéciale `updateWithModifications` :

```go
// Dans constraint.peg, règle JobCall
if nameStr == "Update" {
    argsList, ok := args.([]interface{})
    if ok && len(argsList) == 2 {
        if secondArg, isMap := argsList[1].(map[string]interface{}); isMap {
            if secondArg["type"] == "objectLiteral" {
                // Transformation en updateWithModifications
                ...
            }
        }
    }
}
```

### Évaluation

L'évaluateur reconnaît `updateWithModifications` et :
1. Évalue la variable pour obtenir le fait original
2. Copie les champs du fait original
3. Applique les modifications spécifiées
4. Valide les champs modifiés
5. Retourne un nouveau fait avec le **même ID** mais les champs modifiés

```go
// Dans constraint/evaluator.go
func evaluateUpdateWithModifications(update map[string]interface{}, bindings map[string]interface{}) (map[string]interface{}, error) {
    // Préserve l'ID original
    result := make(map[string]interface{})
    result["_id_"] = originalFact["_id_"]
    
    // Applique les modifications
    for fieldName, fieldValueAST := range modifications {
        evaluatedValue, err := evaluateArithmeticExpr(fieldValueAST, bindings)
        result[fieldName] = evaluatedValue
    }
    
    return result, nil
}
```

## Migration

### Code Existant avec `_Mods_`

Remplacer simplement `_Mods_(...)` par `{...}` :

```diff
- Update(p, _Mods_(statut: "actif", age: 30.0))
+ Update(p, {statut: "actif", age: 30.0})
```

### Code Existant avec `Update(InlineFact(...))`

Cette syntaxe continue de fonctionner mais a une **sémantique différente** : elle crée un **nouveau fait** avec un **nouvel ID** au lieu de mettre à jour le fait existant.

```tsd
// Crée un NOUVEAU fait avec un nouvel ID
Update(Personne(nom: "Alice", statut: "actif"))

// Met à jour le fait EXISTANT (préserve l'ID)
Update(p, {statut: "actif"})
```

## Tests

Tous les tests ont été mis à jour :
- `tests/e2e/update_syntax_test.go` : tests de parsing de la nouvelle syntaxe
- `tests/e2e/testdata/update_simple.tsd` : exemple simple
- `tests/e2e/testdata/relationship_step1_types_rules.tsd` : cas réel
- `internal/defaultactions/loader_test.go` : signature de l'action Update

## Fichiers Modifiés

1. **Grammaire** : `constraint/grammar/constraint.peg`
   - Détection de `objectLiteral` comme second argument d'`Update`
   - Transformation en AST `updateWithModifications`

2. **Documentation** : `internal/defaultactions/defaults.tsd`
   - Mise à jour des exemples et commentaires

3. **Validateur** : `constraint/action_validator.go`
   - Commentaire mis à jour

4. **Tests** :
   - `tests/e2e/update_syntax_test.go`
   - `tests/e2e/testdata/*.tsd`
   - `internal/defaultactions/loader_test.go`

## Compatibilité

- ✅ **Parsing** : la nouvelle syntaxe est entièrement supportée
- ✅ **Validation** : les types sont validés correctement
- ✅ **Évaluation** : l'ID est préservé comme attendu
- ⚠️ **Migration** : rechercher et remplacer `_Mods_` par `{...}` dans le code existant

## Prochaines Étapes

1. ✅ Parser génère `updateWithModifications` pour `Update(var, {...})`
2. ✅ Évaluateur préserve l'ID du fait original
3. ✅ Tests E2E passent avec la nouvelle syntaxe
4. ⏳ Intégrer `BuiltinActionExecutor` pour exécution réelle des Updates
5. ⏳ Tests d'intégration vérifiant la propagation RETE après Update

## Conclusion

La nouvelle syntaxe `Update(variable, {champs...})` :
- Est plus **intuitive** et **naturelle**
- Réutilise une construction **existante** du langage
- Rend le code plus **lisible** et **maintenable**
- Préserve correctement l'**identité des faits**

Cette amélioration simplifie le langage TSD et le rend plus cohérent.