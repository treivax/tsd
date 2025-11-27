# Changelog - Commandes de Suppression

## [2025-01-XX] - Ajout des commandes `remove fact` et `remove rule`

### ✨ Nouvelles fonctionnalités

#### Commande `remove fact`
- **Changement de syntaxe** : La commande de suppression de fait passe de `remove <TypeName> <FactID>` à `remove fact <TypeName> <FactID>`
- **Amélioration de la clarté** : La syntaxe explicite `remove fact` rend le code plus lisible et prévisible
- **Rétrocompatibilité** : ⚠️ **BREAKING CHANGE** - L'ancienne syntaxe ne fonctionne plus

#### Commande `remove rule` (NOUVEAU)
- **Suppression dynamique de règles** : Nouvelle commande `remove rule <RuleID>` pour supprimer une règle du réseau RETE
- **Gestion intelligente du partage** : Les nœuds alpha partagés entre règles sont préservés
- **Comptage de références** : Utilise le `LifecycleManager` pour gérer le cycle de vie des nœuds
- **Nettoyage automatique** : Supprime automatiquement les nœuds alpha qui ne sont plus référencés

### 🔧 Modifications techniques

#### Grammaire PEG (`constraint.peg`)
- Ajout de la règle `RemoveRule` : `"remove" _ "rule" _ ruleID:IdentName`
- Modification de `RemoveFact` : `"remove" _ "fact" _ typeName:IdentName _ factID:FactID`
- Ajout du type de statement `ruleRemoval` dans l'AST
- Modification de `Start` pour gérer `ruleRemovals` en plus de `retractions`

#### Pipeline de contraintes (`constraint_pipeline.go`)
- Nouvelle fonction `processRuleRemovals()` pour traiter les commandes de suppression de règles
- Intégration dans `BuildNetworkFromConstraintFile()` après la construction du réseau
- Gestion des erreurs avec continuation (une erreur n'arrête pas le traitement des autres suppressions)

#### Construction du réseau (`constraint_pipeline_builder.go`)
- **FIX** : Utilisation du `ruleId` de l'expression au lieu d'un index numérique généré
- Les règles sont maintenant identifiées par leur nom déclaré (ex: `adult_check`) et non par `rule_0`, `rule_1`, etc.
- Fallback vers un index numérique si `ruleId` n'est pas présent (compatibilité)

#### Parser
- Régénération du parser avec `pigeon` pour intégrer les modifications de grammaire

### 🧪 Tests ajoutés

#### Tests de parsing (`remove_rule_test.go`)
- `TestParseRemoveFactNewSyntax` : Vérifie la nouvelle syntaxe `remove fact`
- `TestParseRemoveRule` : Vérifie le parsing de `remove rule`
- `TestParseMultipleRemoveCommandsMixed` : Teste plusieurs commandes mixtes
- `TestParseRemoveRuleWithComplexID` : Teste les identifiants complexes de règles
- `TestParseRemoveRuleFromFile` : Teste le parsing depuis un fichier
- `TestOldRemoveSyntaxShouldFail` : Vérifie que l'ancienne syntaxe échoue

#### Tests d'intégration (`remove_rule_integration_test.go`)
- `TestRemoveRuleCommand_ParseAndExecute` : Test de bout en bout de la suppression de règle
- `TestRemoveRuleCommand_MultipleRules` : Suppression de plusieurs règles
- `TestRemoveRuleCommand_WithSharedAlphaNodes` : Vérifie la préservation des nœuds partagés
- `TestRemoveRuleCommand_NonExistentRule` : Gestion des règles inexistantes
- `TestRemoveRuleCommand_AfterFactSubmission` : Suppression après soumission de faits

#### Tests mis à jour
- `remove_fact_test.go` : Migration vers la nouvelle syntaxe `remove fact`
- `network_lifecycle_test.go` : Utilisation des vrais IDs de règles (`r1`, `r2`) au lieu de `rule_0`, `rule_1`
- `alpha_sharing_integration_test.go` : Utilisation des vrais IDs de règles

### 📚 Documentation

#### Nouveaux fichiers
- `docs/REMOVE_COMMANDS.md` : Documentation complète des commandes de suppression
- `examples/remove_commands.tsd` : Exemple d'utilisation des commandes
- `docs/CHANGELOG_REMOVE_COMMANDS.md` : Ce fichier

### 🎯 Cas d'usage

#### `remove fact`
- Rétractation de faits obsolètes
- Mise à jour de données (supprimer puis réinsérer)
- Nettoyage de la mémoire de travail

#### `remove rule`
- Désactivation temporaire de règles
- Optimisation des performances (moins de règles = évaluation plus rapide)
- Reconfiguration dynamique du comportement du système
- Isolation pour les tests

### 📊 Exemple d'utilisation

```tsd
// Définir des types et règles
type Person : <id: string, name: string, age: number>

rule adult_check : {p: Person} / p.age >= 18 ==> notify(p.id)
rule senior_check : {p: Person} / p.age >= 65 ==> alert(p.id)

// Ajouter des faits
Person(id: "P1", name: "Alice", age: 25)
Person(id: "P2", name: "Bob", age: 70)

// Supprimer un fait (nouvelle syntaxe)
remove fact Person P1

// Supprimer une règle (nouvelle fonctionnalité)
remove rule senior_check
```

### ⚠️ Breaking Changes

1. **Syntaxe `remove fact`** : L'ancienne syntaxe `remove <TypeName> <FactID>` ne fonctionne plus. Utiliser `remove fact <TypeName> <FactID>`.

2. **Identifiants de règles** : Les règles utilisent maintenant leur `ruleId` déclaré au lieu d'un index numérique. Cela affecte :
   - `network.RemoveRule(ruleID)` : Utiliser le nom de la règle (`"adult_check"`) et non `"rule_0"`
   - Les nœuds terminaux : Identifiés par `<ruleID>_terminal` (ex: `"adult_check_terminal"`)

### 🔍 Détails d'implémentation

#### Gestion du partage de nœuds
Le système utilise un **compteur de références** (`RefCount`) pour gérer le partage des nœuds alpha :
- Chaque règle utilisant un nœud incrémente son `RefCount`
- La suppression d'une règle décrémente le `RefCount`
- Le nœud n'est physiquement supprimé que quand `RefCount == 0`

#### Logs de débogage
```
🗑️  Suppression de la règle: adult_check
   📊 Nœuds associés à la règle: 2
   ✓ Nœud alpha_21ee82570d6f8f0e marqué pour suppression
   ✓ Nœud adult_check_terminal marqué pour suppression
   🔗 AlphaNode déconnecté de son parent
   🗑️  Nœud supprimé du réseau
✅ Règle adult_check supprimée avec succès (2 nœud(s) supprimé(s))
```

### 🚀 Prochaines étapes

- [ ] Support de la suppression de faits dans le pipeline (actuellement uniquement parsé)
- [ ] Métriques Prometheus pour les suppressions
- [ ] Mode dry-run pour simuler les suppressions
- [ ] Support des transactions (rollback)
- [ ] Interface REPL pour tests interactifs

### 📝 Notes

- Tous les tests existants ont été mis à jour et passent (100% success)
- La fonctionnalité respecte les conventions du projet (MIT License, GoDoc, etc.)
- Le code suit les bonnes pratiques Go (no hardcoding, code générique, etc.)
- La documentation est complète et en français pour cohérence avec le projet

---

**Contributeurs** : Implémenté selon le prompt "add-feature"  
**License** : MIT  
**Date** : 2025-01-XX