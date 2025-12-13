# TODOs Actifs - Liste Consolidée

**Date:** 2024-12-15  
**Status:** Tous non-critiques (améliorations futures)  
**Total:** 7 items

---

## 📋 TODOs Non-Critiques (7 items)

### 1. Migration ParseInput (Compatibilité)

**Fichier:** `constraint/cmd/main.go:248`  
**Priorité:** 🟢 Basse  
**Impact:** Aucun (wrapper de compatibilité)  
**Effort:** ~2h

```go
// TODO: Les tests doivent être migrés pour utiliser ParseInput au lieu de ParseFile
func ParseFile(inputFile string) (interface{}, error) {
    return ParseInput(inputFile)
}
```

**Description:**  
Les tests existants utilisent `ParseFile()` qui est maintenant un simple wrapper autour de `ParseInput()`. Migration recommandée pour cohérence du code.

**Action recommandée:**  
- Identifier tous les appels à `ParseFile()` dans les tests
- Remplacer par `ParseInput()`
- Supprimer la fonction wrapper

---

### 2. Validation Types Personnalisés

**Fichier:** `constraint/constraint_facts.go:71`  
**Priorité:** 🟡 Moyenne  
**Impact:** Extensibilité future  
**Effort:** ~4h

```go
if !ValidPrimitiveTypes[expectedType] {
    // Type personnalisé ou non standard accepté pour extensibilité
    // TODO: Valider que le type personnalisé existe dans le programme
    return nil
}
```

**Description:**  
Actuellement, les types personnalisés sont acceptés sans validation. Il serait utile de vérifier que le type existe réellement dans le programme.

**Action recommandée:**  
- Créer un registre des types définis dans le programme
- Valider les types personnalisés lors de la validation des faits
- Retourner erreur explicite si type inconnu

---

### 3. Génération Table de Règles

**Fichier:** `constraint/parser.go:6034`  
**Priorité:** 🟢 Basse  
**Impact:** Performance marginale  
**Effort:** ~3h

```go
func (p *parser) parse(g *grammar) (val any, err error) {
    // ...
    // TODO : not super critical but this could be generated
    p.buildRulesTable(g)
    // ...
}
```

**Description:**  
La construction de la table de règles pourrait être générée au lieu d'être construite à runtime, ce qui améliorerait légèrement les performances.

**Action recommandée:**  
- Évaluer le gain de performance potentiel
- Si significatif, générer la table lors de la compilation du parser
- Sinon, conserver l'implémentation actuelle

---

### 4. Support Opérateur Modulo (%)

**Fichier:** `rete/arithmetic_alpha_extraction_test.go:317`  
**Priorité:** 🟡 Moyenne  
**Impact:** Feature manquante  
**Effort:** ~6h

```go
// TODO: Enable when parser supports % operator
// {
//     name: "modulo operation",
//     expression: "150 % 100",
//     expected: 50,
// },
```

**Description:**  
L'opérateur modulo (%) n'est pas encore supporté par le parser. Le test est commenté en attendant l'implémentation.

**Action recommandée:**  
- Ajouter le support de l'opérateur `%` dans le parser
- Implémenter l'évaluation dans `ArithmeticEvaluator`
- Décommenter et valider le test

**Dépendances:**  
- Modification du parser (constraint/parser.peg)
- Régénération du parser
- Tests de régression

---

### 5. Comparaison Profonde Conditions

**Fichier:** `rete/beta_sharing_interface.go:444`  
**Priorité:** 🟡 Moyenne  
**Impact:** Partage plus agressif (optimisation)  
**Effort:** ~8h

```go
func CanShareJoinNodes(sig1, sig2 CanonicalJoinSignature) bool {
    // ... validations existantes ...
    
    // TODO: Deep comparison of normalized conditions
    
    return true
}
```

**Description:**  
Actuellement, la comparaison des conditions pour déterminer si deux JoinNodes peuvent être partagés est basique. Une comparaison profonde permettrait un partage plus agressif.

**Action recommandée:**  
- Implémenter normalisation des conditions (ordre, équivalences)
- Comparer les conditions normalisées
- Ajouter tests de validation du partage amélioré

**Attention:**  
- Ne pas introduire de partage incorrect (régression)
- Tests exhaustifs requis

---

### 6. Métriques Détaillées JoinNode

**Fichier:** `rete/beta_sharing_stats.go:135-136`  
**Priorité:** 🟢 Basse  
**Impact:** Observabilité (non-bloquant)  
**Effort:** ~2h

```go
SharedJoinNodeDetails{
    RightMemorySize:  len(node.RightMemory.Tokens),
    ResultMemorySize: len(node.ResultMemory.Tokens),
    CreatedAt:        time.Time{}, // TODO: Track creation time
    LastAccessedAt:   time.Now(),
    ActivationCount:  0,            // TODO: Track activation count
}
```

**Description:**  
Les métriques `CreatedAt` et `ActivationCount` ne sont pas encore trackées, ce qui limite l'observabilité du système.

**Action recommandée:**  
- Ajouter champs dans structure `JoinNode`
- Incrémenter `ActivationCount` dans `ActivateLeft/Right`
- Initialiser `CreatedAt` lors de la création
- Exporter métriques pour monitoring

---

### 7. AlphaConditionEvaluator Arithmétique

**Fichier:** `rete/condition_splitter.go:86`  
**Priorité:** 🟡 Moyenne  
**Impact:** Performance (conditions arithmétiques en alpha)  
**Effort:** ~10h

```go
if isSimpleArithmetic {
    // ...
} else {
    // Complex arithmetic expression - keep in beta for now
    // TODO: Enhance AlphaConditionEvaluator to handle arithmetic
    splitCond.Type = ConditionTypeBeta
    betaConditions = append(betaConditions, splitCond)
}
```

**Description:**  
Les expressions arithmétiques complexes sont actuellement évaluées dans les nœuds beta. Les déplacer vers alpha améliorerait les performances en filtrant plus tôt.

**Action recommandée:**  
- Étendre `AlphaConditionEvaluator` pour supporter arithmétique
- Implémenter cache d'évaluation si nécessaire
- Benchmarker les gains de performance
- Migrer progressivement les conditions

**Complexité:**  
- Nécessite refactoring de l'évaluateur alpha
- Tests exhaustifs requis
- Impact performance à valider

---

## 📊 Statistiques

```
Total TODOs:           7
Priorité Haute:        0  ✅
Priorité Moyenne:      4  🟡
Priorité Basse:        3  🟢

Effort Total Estimé:   ~35h
```

## 🎯 Priorisation Recommandée

### Phase 1 (Court terme - ~8h)
1. ✅ **TODO #1** - Migration ParseInput (2h)
2. ✅ **TODO #6** - Métriques JoinNode (2h)
3. ✅ **TODO #3** - Génération table règles (3h) - Si gain mesuré

### Phase 2 (Moyen terme - ~12h)
4. ✅ **TODO #2** - Validation types personnalisés (4h)
5. ✅ **TODO #4** - Opérateur modulo (6h)

### Phase 3 (Long terme - ~18h)
6. ✅ **TODO #5** - Comparaison profonde conditions (8h)
7. ✅ **TODO #7** - Alpha arithmétique (10h)

## 📝 Notes Importantes

### Aucun TODO n'est bloquant pour la production ✅

- Le système fonctionne parfaitement sans ces améliorations
- Tous les tests passent (100%)
- Aucune régression détectée
- Code de qualité production

### Ces TODOs sont des optimisations et améliorations futures

- **Performance** : #3, #7
- **Features** : #4
- **Qualité** : #1, #2, #5
- **Observabilité** : #6

### Validation requise après chaque implémentation

- ✅ Tests unitaires passent
- ✅ Tests E2E passent
- ✅ Aucune régression introduite
- ✅ Documentation mise à jour
- ✅ Benchmarks si pertinent

---

## 🔗 Références

- **Rapport de validation** : `SYNTHESE_VALIDATION_FINALE.md`
- **Détails techniques** : `VALIDATION_FINALE_POST_FIX.md`
- **Architecture** : `docs/architecture/BINDINGS_DESIGN.md`
- **Changelog** : `CHANGELOG.md`

---

## 📅 Historique

| Date | Action | Détails |
|------|--------|---------|
| 2024-12-15 | Création | Liste initiale des 7 TODOs non-critiques |

---

*Dernière révision: 2024-12-15*  
*Status: ✅ Tous les TODOs critiques résolus - Production Ready*