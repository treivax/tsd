# 🐛 BUG RETE-001: Root Cause Analysis

**Date**: 2025-12-01  
**Severity**: Majeure  
**Type**: Fonctionnel / Performance  
**Status**: Analysé, correction en cours

---

## 1. DESCRIPTION DU BUG

### Problème Identifié

Le builder RETE actuel **ne sépare pas** les conditions alpha (tests sur un seul fait) des conditions beta (tests entre plusieurs faits). Toutes les conditions sont placées dans le JoinNode, violant le principe fondamental de l'architecture RETE.

### Exemple Concret

Pour la règle suivante :
```tsd
rule test : {c: Commande, p: Produit} /
    c.produit_id == p.id AND c.qte > 5
    ==> resultat(c.id, p.id)
```

**Structure actuelle (BUGGUÉE)** :
```
TypeNode(Commande) → PassthroughAlpha → JoinNode(c.produit_id == p.id AND c.qte > 5)
                                              ⋈
TypeNode(Produit)  → PassthroughAlpha ────────┘
```

**Structure attendue (CORRECTE)** :
```
TypeNode(Commande) → AlphaNode(c.qte > 5) → PassthroughAlpha → JoinNode(c.produit_id == p.id)
                                                                      ⋈
TypeNode(Produit)  → PassthroughAlpha ────────────────────────────────┘
```

---

## 2. ROOT CAUSE ANALYSIS (5 Pourquoi)

### Pourquoi 1 : Les conditions alpha ne sont pas évaluées avant la jointure
**→** Parce qu'il n'y a pas d'AlphaNode filtrant créé

### Pourquoi 2 : Aucun AlphaNode filtrant n'est créé
**→** Parce que le builder ne décompose pas les conditions AND

### Pourquoi 3 : Le builder ne décompose pas les conditions AND
**→** Parce que `JoinRuleBuilder.CreateJoinRule()` passe la condition complète au JoinNode

### Pourquoi 4 : La condition complète est passée au JoinNode
**→** Parce qu'il n'existe pas de composant `ConditionSplitter` pour séparer alpha vs beta

### Pourquoi 5 : Pas de `ConditionSplitter`
**→** **CAUSE RACINE** : Manque d'implémentation de la décomposition des conditions dans l'architecture du builder

---

## 3. ANALYSE TECHNIQUE

### 3.1. Composants Impliqués

| Composant | Fichier | Rôle | Problème |
|-----------|---------|------|----------|
| `JoinRuleBuilder` | `builder_join_rules.go` | Créer les règles de jointure | Passe condition complète au JoinNode |
| `ConstraintPipeline` | `constraint_pipeline_builder.go` | Orchestration | Délègue sans décomposer |
| `BuilderUtils` | `builder_utils.go` | Utilitaires | Pas d'outil de décomposition |
| **MANQUANT** | N/A | **Séparation alpha/beta** | **Non implémenté** |

### 3.2. Flux Actuel

```
ParseConstraintFile
    ↓
ConvertToReteProgram (AST)
    ↓
buildNetwork(types, expressions)
    ↓
createRuleNodes(expressions)
    ↓
CreateJoinRule(condition) ← PROBLÈME ICI
    ↓
NewJoinNode(condition)  ← Reçoit la condition complète (alpha + beta)
```

### 3.3. Code Problématique

**Fichier**: `rete/builder_join_rules.go`, ligne ~51-82

```go
func (jrb *JoinRuleBuilder) createBinaryJoinRule(
    network *ReteNetwork,
    ruleID string,
    variableNames []string,
    variableTypes []string,
    condition map[string]interface{},  // ← Condition complète reçue
    terminalNode *TerminalNode,
) error {
    // ...
    joinNode := NewJoinNode(
        ruleID+"_join", 
        condition,  // ← Condition complète passée sans décomposition
        leftVars, 
        rightVars, 
        varTypes, 
        jrb.utils.storage
    )
    // ...
}
```

**Aucune séparation** des conditions n'est effectuée avant la création du JoinNode.

---

## 4. IMPACT

### 4.1. Performance

| Scénario | Sans Filtre Alpha | Avec Filtre Alpha | Économie |
|----------|-------------------|-------------------|----------|
| 3 Commandes × 2 Produits | 6 évaluations | 4 évaluations | **33%** |
| 10 × 10 | 100 évaluations | ~67 évaluations | **33%** |
| 100 × 100 | 10,000 évaluations | ~6,700 évaluations | **33%** |
| 1000 × 1000 | 1,000,000 évaluations | ~670,000 évaluations | **33%** |

**Plus il y a de faits, plus l'impact est important.**

### 4.2. Violations RETE

1. **Pas de filtrage précoce** : Les faits ne sont pas filtrés avant la jointure
2. **Évaluations redondantes** : La condition alpha est réévaluée pour chaque paire
3. **Pas de partage** : Les conditions alpha identiques entre règles ne sont pas partagées
4. **Architecture incorrecte** : Violation du principe de séparation alpha/beta

### 4.3. Utilisateurs Affectés

- ✅ **Toutes les règles multi-variables** avec conditions mixtes (alpha + beta)
- ✅ **Règles arithmétiques** (ex: `c.qte * 23 - 10 > 0`)
- ✅ **Règles avec filtres sur champs** (ex: `c.statut == "actif"`)
- ❌ **Règles alpha simples** (une seule variable) : non affectées
- ❌ **Règles beta pures** (seulement jointure) : non affectées

---

## 5. ANALYSE D'ALTERNATIVES

### Option 1: Correction Simple (Décomposition locale)
**Approche** : Décomposer dans `CreateJoinRule()` uniquement

**Avantages** :
- ✅ Correction ciblée
- ✅ Peu de changements

**Inconvénients** :
- ❌ Code non réutilisable
- ❌ Difficile à maintenir
- ❌ Pas extensible

**Verdict** : ❌ Non recommandé

### Option 2: Nouveau Composant `ConditionSplitter` (Recommandé)
**Approche** : Créer un composant dédié à la décomposition des conditions

**Avantages** :
- ✅ Séparation des responsabilités
- ✅ Réutilisable (alpha rules, exists rules, etc.)
- ✅ Testable indépendamment
- ✅ Extensible (normalisation AST future)

**Inconvénients** :
- ⚠️ Plus de code à écrire
- ⚠️ Nouveau composant à maintenir

**Verdict** : ✅ **RECOMMANDÉ**

### Option 3: Refonte Complète (Trop complexe)
**Approche** : Revoir toute l'architecture du builder

**Avantages** :
- ✅ Architecture parfaite

**Inconvénients** :
- ❌ Trop de changements
- ❌ Risque de régression élevé
- ❌ Temps de développement important

**Verdict** : ❌ Overkill pour ce bug

---

## 6. SOLUTION CHOISIE

### Architecture

**Nouveau composant** : `ConditionSplitter`

```
ConditionSplitter
    ├── SplitConditions(condition) → (alphaConditions, betaConditions)
    ├── ClassifyCondition(condition) → ConditionType (alpha/beta)
    ├── ExtractVariables(condition) → []string
    └── IsAlphaCondition(condition) → bool
```

### Flux Corrigé

```
CreateJoinRule(condition)
    ↓
[NEW] ConditionSplitter.SplitConditions(condition)
    ↓
alphaConditions, betaConditions
    ↓
CreateAlphaNodes(alphaConditions)  ← Filtrage précoce
    ↓
CreateJoinNode(betaConditions only) ← Seulement conditions inter-faits
    ↓
Chaîner: TypeNode → AlphaFilter → PassthroughAlpha → JoinNode
```

### Modifications Requises

1. **Nouveau fichier** : `rete/condition_splitter.go`
   - Implémentation de `ConditionSplitter`
   - Classification alpha vs beta
   - Extraction de variables

2. **Modification** : `rete/builder_join_rules.go`
   - Utiliser `ConditionSplitter` avant création JoinNode
   - Créer AlphaNodes filtrants pour conditions alpha
   - Chaîner correctement les nœuds

3. **Nouveau fichier** : `rete/condition_splitter_test.go`
   - Tests unitaires pour classification
   - Tests de décomposition AND
   - Tests avec conditions imbriquées

4. **Modification** : `rete/bug_rete001_alpha_beta_separation_test.go`
   - Ajouter test de vérification post-correction
   - Valider structure avec AlphaNodes filtrants

---

## 7. CRITÈRES DE SUCCÈS

### Tests de Non-Régression

1. ✅ **Test de reproduction** doit échouer après correction (bug résolu)
2. ✅ **AlphaNodes filtrants** doivent être créés pour conditions alpha
3. ✅ **JoinNodes** ne doivent contenir que conditions beta
4. ✅ **Chaînage correct** : TypeNode → AlphaFilter → JoinNode
5. ✅ **Résultats identiques** : les mêmes actions doivent se déclencher

### Performance

1. ✅ Réduction mesurable des évaluations de jointure
2. ✅ Filtrage précoce des faits ne satisfaisant pas les conditions alpha
3. ✅ Partage possible des AlphaNodes filtrants entre règles

### Code Quality

1. ✅ Tests unitaires pour `ConditionSplitter`
2. ✅ Tests E2E validant la structure complète
3. ✅ Documentation du nouveau composant
4. ✅ Pas de régression sur tests existants

---

## 8. PLAN D'IMPLÉMENTATION

### Phase 1: Création du Splitter (Priorité Haute)
- [ ] Créer `condition_splitter.go`
- [ ] Implémenter classification alpha/beta
- [ ] Tests unitaires

### Phase 2: Intégration Builder (Priorité Haute)
- [ ] Modifier `CreateJoinRule()` pour utiliser le splitter
- [ ] Créer AlphaNodes filtrants
- [ ] Chaîner correctement les nœuds

### Phase 3: Validation (Priorité Haute)
- [ ] Vérifier que test de reproduction échoue
- [ ] Tests E2E avec structure corrigée
- [ ] Validation performance

### Phase 4: Documentation (Priorité Moyenne)
- [ ] Documenter `ConditionSplitter`
- [ ] Mettre à jour CHANGELOG
- [ ] Exemple d'utilisation

---

## 9. RISQUES ET MITIGATIONS

| Risque | Probabilité | Impact | Mitigation |
|--------|-------------|--------|------------|
| Régression sur règles simples | Faible | Élevé | Tests exhaustifs avant/après |
| Mauvaise classification alpha/beta | Moyenne | Élevé | Tests unitaires robustes |
| Performance dégradée (overhead) | Faible | Moyen | Benchmarks comparatifs |
| Conditions edge cases non gérées | Moyenne | Moyen | Tests avec conditions complexes |

---

## 10. RÉFÉRENCES

- **Test de reproduction** : `rete/bug_rete001_alpha_beta_separation_test.go`
- **Fichier TSD de test** : `rete/testdata/bug_rete001_minimal.tsd`
- **Code source principal** : `rete/builder_join_rules.go`
- **Principe RETE** : Alpha (filtres unaires) vs Beta (jointures binaires)
- **Prompt utilisé** : `.github/prompts/fix-bug.md`

---

## 11. NOTES

- Ce bug existe probablement depuis l'implémentation initiale des règles multi-variables
- L'impact augmente exponentiellement avec le nombre de faits
- La correction améliore aussi la maintenabilité (séparation des responsabilités)
- Opportunité d'ajouter des optimisations futures (normalisation AST, détection commutativité)

---

**Document créé le** : 2025-12-01  
**Auteur** : TSD Engineering Team  
**Version** : 1.0