# Progression du Refactoring - constraint_pipeline_builder.go

**Date de début:** 2024
**Fichier cible:** `rete/constraint_pipeline_builder.go`
**Statut:** 🟡 En cours (Phase 1-4 complétées)

---

## 📊 État Actuel

### Fichiers Créés

✅ **Phase 1-2: Infrastructure de base**
- `rete/builders/` - Nouveau package créé
- `rete/builders/utils.go` - Utilitaires communs (154 lignes)
- `rete/builders/types.go` - Builder pour les types (97 lignes)
- `rete/builders/alpha_rules.go` - Builder pour règles alpha (101 lignes)

### Métriques

| Métrique | Avant | Actuel | Cible | Progression |
|----------|-------|--------|-------|-------------|
| **Fichier principal** | 1,030 lignes | 1,030 lignes | 200 lignes | 0% |
| **Nouveaux builders** | 0 | 3 fichiers | 7 fichiers | 43% |
| **Lignes extraites** | 0 | ~352 lignes | ~830 lignes | 42% |
| **Complexité max** | 18 | 18 | 10 | 0% |

---

## ✅ Phases Complétées

### Phase 1: Préparation (✅ Complété)
- [x] Création du package `rete/builders/`
- [x] Structure de fichiers définie
- [x] Plan de refactoring documenté

### Phase 2: Extraction des Utilitaires (✅ Complété)
- [x] `builders/utils.go` créé avec:
  - Constants (ConditionType*, NodeSide*)
  - BuilderUtils struct
  - CreatePassthroughAlphaNode()
  - ConnectTypeNodeToBetaNode()
  - CreateTerminalNode()
  - Fonctions helper (GetStringField, GetMapField, etc.)
  - BuildVarTypesMap()

### Phase 3: Extraction des Types (✅ Complété)
- [x] `builders/types.go` créé avec:
  - TypeBuilder struct
  - CreateTypeNodes()
  - CreateTypeDefinition()

### Phase 4: Extraction des Règles Alpha (✅ Complété)
- [x] `builders/alpha_rules.go` créé avec:
  - AlphaRuleBuilder struct
  - CreateAlphaRule()
  - getVariableInfo()
  - createAlphaNodeWithTerminal()

---

## 🔄 Phases en Cours

### Phase 5: Extraction des Règles EXISTS (⏳ À faire)
Fichier: `builders/exists_rules.go`

**Fonctions à extraire:**
- [ ] ExistsRuleBuilder struct
- [ ] CreateExistsRule() (51 lignes)
- [ ] ExtractExistsVariables() (44 lignes)
- [ ] ExtractExistsConditions() (28 lignes)
- [ ] ConnectExistsNodeToTypeNodes() (17 lignes)

**Total:** ~140 lignes

### Phase 6: Extraction des Règles de Jointure (⏳ À faire)
Fichier: `builders/join_rules.go`

**Fonctions à extraire:**
- [ ] JoinRuleBuilder struct
- [ ] CreateJoinRule() (28 lignes)
- [ ] CreateBinaryJoinRule() (80 lignes)
- [ ] CreateCascadeJoinRule() (99 lignes)
- [ ] CreateCascadeJoinRuleWithBuilder() (95 lignes)

**Total:** ~302 lignes

**Refactoring nécessaire:**
- Décomposer CreateCascadeJoinRuleWithBuilder en 3 fonctions:
  - buildJoinPatterns()
  - buildChainWithBuilder()
  - connectChainToNetwork()

### Phase 7: Extraction des Règles d'Accumulation (⏳ À faire)
Fichier: `builders/accumulator_rules.go`

**Fonctions à extraire:**
- [ ] AccumulatorRuleBuilder struct
- [ ] IsMultiSourceAggregation() (48 lignes)
- [ ] CreateMultiSourceAccumulatorRule() (**154 lignes** ⚠️)
- [ ] CreateAccumulatorRule() (69 lignes)

**Total:** ~271 lignes

**Refactoring critique (priorité haute):**
CreateMultiSourceAccumulatorRule doit être décomposé en:
- [ ] createJoinChainForSources() (~50 lignes)
- [ ] createMultiSourceAccumulatorNode() (~40 lignes)
- [ ] connectAccumulatorChainToTerminal() (~40 lignes)
- [ ] Fonction principale simplifiée (~24 lignes)

### Phase 8: Orchestration Centrale (⏳ À faire)
Fichier: `builders/rules.go`

**À créer:**
- [ ] RuleBuilder struct (agrège tous les builders)
- [ ] CreateRuleNodes() (25 lignes)
- [ ] CreateSingleRule() - simplifié à ~50 lignes (actuellement 82)

**Total:** ~75 lignes

---

## 🚧 Prochaines Étapes

### Immédiat (Priorité 1)
1. **Créer `builders/exists_rules.go`**
   - Extraire les 4 fonctions EXISTS
   - Ajouter tests unitaires

2. **Créer `builders/join_rules.go`**
   - Extraire les 4 fonctions de jointure
   - Refactorer CreateCascadeJoinRuleWithBuilder
   - Ajouter tests unitaires

3. **Créer `builders/accumulator_rules.go`**
   - Extraire les 3 fonctions d'accumulation
   - **CRITIQUE:** Décomposer CreateMultiSourceAccumulatorRule
   - Ajouter tests unitaires

### Court terme (Priorité 2)
4. **Créer `builders/rules.go`**
   - Créer RuleBuilder qui orchestre tous les builders
   - Simplifier CreateSingleRule

5. **Refactorer `constraint_pipeline_builder.go`**
   - Intégrer les builders
   - Réduire à ~200 lignes
   - Déléguer toute la logique aux builders

### Moyen terme (Priorité 3)
6. **Tests d'intégration**
   - Vérifier que tous les tests existants passent
   - Ajouter tests spécifiques aux builders

7. **Benchmarks de performance**
   - Comparer performances avant/après
   - S'assurer qu'il n'y a pas de régression

8. **Documentation**
   - Documenter chaque builder
   - Mettre à jour les exemples
   - Guide d'utilisation

---

## 📈 Avancement Global

### Code

```
Progression: ████████░░░░░░░░░░░░ 42% (352/830 lignes extraites)
```

| Phase | Status | Lignes | Progression |
|-------|--------|--------|-------------|
| Utils | ✅ | 154 | 100% |
| Types | ✅ | 97 | 100% |
| Alpha | ✅ | 101 | 100% |
| EXISTS | ⏳ | 0/140 | 0% |
| Join | ⏳ | 0/302 | 0% |
| Accumulator | ⏳ | 0/271 | 0% |
| Orchestration | ⏳ | 0/75 | 0% |

### Tests

```
Progression: ░░░░░░░░░░░░░░░░░░░░ 0% (0/7 fichiers de tests)
```

- [ ] builders/utils_test.go
- [ ] builders/types_test.go
- [ ] builders/alpha_rules_test.go
- [ ] builders/exists_rules_test.go
- [ ] builders/join_rules_test.go
- [ ] builders/accumulator_rules_test.go
- [ ] builders/rules_test.go

---

## 🎯 Objectifs de Qualité

### Métriques Cibles

| Métrique | Avant | Après (cible) | Statut |
|----------|-------|---------------|--------|
| Fichier principal | 1,030 lignes | 200 lignes | ⏳ 0% |
| Fonctions >100 lignes | 3 | 0 | ⏳ 0% |
| Fonctions >80 lignes | 5 | 0 | ⏳ 0% |
| Complexité max | 18 | ≤10 | ⏳ 0% |
| Maintenabilité | 72/100 | 85/100 | ⏳ 72/100 |
| Couverture tests | ? | >80% | ⏳ 0% |

### Fonctions Critiques à Refactorer

| Fonction | Lignes | Complexité | Priorité | Statut |
|----------|--------|------------|----------|--------|
| createMultiSourceAccumulatorRule | 154 | 18 | 🔴 Haute | ⏳ |
| createCascadeJoinRuleWithBuilder | 95 | 16 | 🟡 Moyenne | ⏳ |
| createSingleRule | 82 | 14 | 🟡 Moyenne | ⏳ |
| createCascadeJoinRule | 99 | 12 | 🟢 Basse | ⏳ |
| createBinaryJoinRule | 80 | 10 | 🟢 Basse | ⏳ |

---

## 💡 Décisions Techniques

### Architecture Choisie

**Pattern:** Builder + Composition
- Chaque type de règle a son propre builder
- Les builders partagent un BuilderUtils commun
- Un RuleBuilder orchestre tous les builders spécialisés

**Avantages:**
- ✅ Séparation claire des responsabilités
- ✅ Testabilité maximale
- ✅ Réutilisabilité du code
- ✅ Facilite l'ajout de nouveaux types de règles

**Inconvénients:**
- ⚠️ Plus de fichiers à maintenir
- ⚠️ Navigation légèrement plus complexe
- ⚠️ Overhead minimal (création d'objets builders)

### Conventions de Code

```go
// Nommage des builders
type <Type>Builder struct {
    utils *BuilderUtils
}

// Méthodes publiques: Create/Extract/Connect
func (b *Builder) CreateRule(...)
func (b *Builder) ExtractVariables(...)
func (b *Builder) ConnectNodes(...)

// Méthodes privées: build/prepare/setup
func (b *Builder) buildPattern(...)
func (b *Builder) prepareCondition(...)
func (b *Builder) setupConnection(...)
```

---

## 🐛 Problèmes Rencontrés

### Aucun problème majeur pour l'instant

Les phases 1-4 se sont déroulées sans incident.

---

## 📝 Notes de Travail

### Points d'Attention

1. **Rétrocompatibilité**
   - L'API publique de ConstraintPipeline ne doit pas changer
   - Les tests existants doivent continuer à passer

2. **Performance**
   - Aucune dégradation de performance acceptable
   - Benchmarks obligatoires avant merge

3. **Documentation**
   - Chaque builder doit être documenté
   - Exemples d'utilisation nécessaires

### Dépendances Externes

Les builders utilisent les types du package `rete`:
- `rete.ReteNetwork`
- `rete.Storage`
- `rete.Action`
- `rete.Node`
- Tous les types de nœuds (AlphaNode, JoinNode, etc.)

Aucune modification de ces types n'est nécessaire.

---

## 🔗 Ressources

### Documents de Référence
- [Plan de Refactoring Détaillé](./REFACTORING_PLAN_PIPELINE_BUILDER.md)
- [Statistiques du Code Manuel](./MANUAL_CODE_STATISTICS.md)
- [Rapport de Deep Clean](./DEEP_CLEAN_REPORT.md)

### Commandes Utiles

```bash
# Lancer les tests du package builders
go test ./rete/builders -v

# Vérifier la complexité
gocyclo -over 10 ./rete/builders

# Vérifier le formatage
go fmt ./rete/builders/...

# Lancer tous les tests RETE
go test ./rete/... -v
```

---

## 📅 Timeline

### Temps Estimé par Phase Restante

| Phase | Temps Estimé | Complexité |
|-------|--------------|------------|
| Phase 5 (EXISTS) | 1h | Moyenne |
| Phase 6 (Join) | 2h | Élevée |
| Phase 7 (Accumulator) | 3h | Très élevée |
| Phase 8 (Orchestration) | 1h | Moyenne |
| Phase 9 (Refactoring main) | 1h | Moyenne |
| Phase 10 (Tests) | 2h | Moyenne |
| **Total restant** | **10h** | |

### Phases Déjà Complétées

| Phase | Temps Réel |
|-------|------------|
| Phase 1 (Préparation) | 15min |
| Phase 2 (Utils) | 30min |
| Phase 3 (Types) | 20min |
| Phase 4 (Alpha) | 25min |
| **Total complété** | **1.5h** |

**Progression totale:** 13% (1.5h / 11.5h)

---

## ✅ Checklist Finale

### Avant de Merger

- [ ] Toutes les phases complétées
- [ ] Tous les tests passent (existants + nouveaux)
- [ ] Benchmarks exécutés (pas de régression)
- [ ] Complexité cyclomatique réduite (max 10)
- [ ] Fichier principal réduit à ~200 lignes
- [ ] Documentation à jour
- [ ] Revue de code effectuée
- [ ] Exemples mis à jour si nécessaire

### Validation Qualité

- [ ] `go test ./...` passe à 100%
- [ ] `go vet ./...` sans erreurs
- [ ] `gocyclo -over 10 .` conforme
- [ ] `go fmt` appliqué
- [ ] Couverture de tests >80% pour les builders

---

**Dernière mise à jour:** Phase 1-4 complétées
**Prochaine action:** Créer `builders/exists_rules.go`
**Responsable:** À définir
**Statut général:** 🟡 En progression (42% du code extrait, 0% des tests)