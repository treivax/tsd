# TODO - Migration des Tests - Nouvelle Gestion des IDs

> **Statut** : 🔴 EN ATTENTE - Migration complexe nécessitant refactoring architectural
> 
> **Créé le** : 2025-12-19
> 
> **Contexte** : Execution du prompt `/home/resinsec/dev/tsd/scripts/new_ids/07-prompt-tests-unit.md`

## 🎯 Objectif

Migrer tous les tests unitaires pour utiliser la nouvelle gestion des identifiants internes (`_id_`) de manière cohérente et conforme à l'architecture production.

## ⚠️ BLOQUEURS IDENTIFIÉS

### 1. Architecture Hybride (CRITIQUE)

**Problème** : Deux APIs incompatibles coexistent dans le code :

```go
// API 1 : Tests (ANCIEN - attend ID pré-rempli)
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
    // Attend que fact.ID soit déjà défini
}

// API 2 : Production (NOUVEAU - génère ID automatiquement)
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
    // Extrait "id" de Fields et le met dans fact.ID
}
```

**Impact** : Impossible de migrer proprement les tests sans d'abord uniformiser l'API.

**Action requise** : 
```bash
# TODO: Créer issue GitHub
# Titre: "Uniformiser l'API de soumission de faits (SubmitFact vs SubmitFactsFromGrammar)"
# Labels: architecture, refactoring, breaking-change
# Priorité: HIGH
```

### 2. Validation TypeNode Incohérente

**Fichier** : `rete/node_type.go:114-118`

**Problème** : La validation attend "id" dans `fact.ID`, pas dans `fact.Fields["id"]`

**Action requise** :
- [ ] Modifier `validateFact()` pour accepter "id" dans `Fields`
- [ ] Copier automatiquement de `Fields["id"]` vers `fact.ID`
- [ ] Ajouter tests de validation pour les deux cas

### 3. Comportement TerminalNode Non Documenté

**Observation** : Les actions s'exécutent mais `TerminalNode.Memory.Tokens` reste vide.

**Investigation nécessaire** :
- [ ] Documenter le cycle de vie des tokens dans TerminalNode
- [ ] Vérifier si c'est un comportement intentionnel (tokens consommés après action)
- [ ] Mettre à jour les tests pour refléter le comportement réel

## 📋 PLAN D'ACTION DÉTAILLÉ

### Phase 0 : Préparation (AVANT toute migration de tests)

#### A. Corriger l'Architecture

```bash
# 1. Créer une branche dédiée
git checkout -b fix/uniform-fact-submission-api

# 2. Implémenter les changements suivants
```

**Fichier : `rete/network_manager.go`**
```go
// Refactorer SubmitFact pour uniformiser avec SubmitFactsFromGrammar
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
    // Si l'ID n'est pas fourni, essayer de l'extraire de Fields
    if fact.ID == "" {
        if id, ok := fact.Fields["id"].(string); ok {
            fact.ID = id
            // NE PAS supprimer de Fields - laisser pour validation
        } else {
            // Si pas d'ID dans Fields, générer par hash?
            // OU retourner erreur selon la politique
            return fmt.Errorf("impossible de déterminer l'ID du fait de type %s", fact.Type)
        }
    }
    
    // Reste du code existant...
}
```

**Fichier : `rete/node_type.go`**
```go
// Modifier validateFact pour accepter "id" dans Fields OU fact.ID
func (tn *TypeNode) validateFact(fact *Fact) error {
    for _, field := range tn.TypeDefinition.Fields {
        if field.Name == "id" {
            // Vérifier fact.ID d'abord
            if fact.ID != "" {
                continue // OK, ID fourni
            }
            
            // Sinon, chercher dans Fields
            if id, ok := fact.Fields["id"].(string); ok {
                fact.ID = id // Copier dans fact.ID
                continue
            }
            
            // Ni dans fact.ID ni dans Fields
            return fmt.Errorf("champ manquant: %s", field.Name)
        }
        
        // Validation des autres champs...
        value, exists := fact.Fields[field.Name]
        if !exists {
            return fmt.Errorf("champ manquant: %s", field.Name)
        }
        
        if !tn.isValidType(value, field.Type) {
            return fmt.Errorf("type invalide pour le champ %s: attendu %s", field.Name, field.Type)
        }
    }
    return nil
}
```

**Créer tests de non-régression** :
```go
// Fichier: rete/node_type_id_migration_test.go
func TestTypeNode_ValidateFact_IDInFields(t *testing.T) {
    // Test 1: ID dans fact.ID (ancien pattern)
    // Test 2: ID dans fact.Fields["id"] (nouveau pattern)
    // Test 3: ID dans les deux (vérifier cohérence)
    // Test 4: ID absent partout (doit échouer)
}
```

#### B. Créer Helper de Test

**Fichier : `rete/test_helpers.go`** (nouveau)
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// NewTestFact crée un fait pour les tests unitaires.
// Cette fonction encapsule la logique de création de faits
// pour faciliter la migration progressive vers la nouvelle API.
//
// Deprecated: À terme, les tests devraient utiliser SubmitFactsFromGrammar
// Cette fonction est temporaire pour la période de transition.
func NewTestFact(id, factType string, fields map[string]interface{}) *Fact {
    return &Fact{
        ID:     id,
        Type:   factType,
        Fields: fields,
    }
}

// SubmitTestFact soumet un fait créé pour les tests.
// Utilise le nouveau pattern unifié.
func (rn *ReteNetwork) SubmitTestFact(id, factType string, fields map[string]interface{}) error {
    fact := &Fact{
        ID:     id,
        Type:   factType,
        Fields: fields,
    }
    return rn.SubmitFact(fact)
}
```

#### C. Documentation

**Créer : `docs/migration/test-ids-migration.md`**
```markdown
# Guide de Migration - Tests et Identifiants

## Contexte

La gestion des identifiants de faits a évolué pour cacher le champ `ID` 
et utiliser `_id_` en interne uniquement.

## Pattern Actuel (à éviter)

```go
fact := Fact{ID: "p1", Type: "Person", Fields: {...}}
network.SubmitFact(&fact)
```

## Nouveau Pattern (recommandé)

```go
fact := Fact{Type: "Person", Fields: map[string]interface{}{
    "id": "p1",
    ...
}}
network.SubmitFact(&fact)
```

## Helper Temporaire

```go
fact := NewTestFact("p1", "Person", map[string]interface{}{...})
network.SubmitFact(fact)
```

## TODO

Cette migration est temporaire. À terme, tous les tests devraient utiliser
`SubmitFactsFromGrammar()` comme le code de production.
```

### Phase 1 : Tests Critiques (constraint/)

#### Fichiers à migrer en priorité :

- [ ] `constraint/id_generator_test.go` - Tests de génération d'IDs
- [ ] `constraint/id_generator_fact_references_test.go` - Références de faits
- [ ] `constraint/id_generator_edge_cases_test.go` - Cas limites
- [ ] `constraint/primary_key_validation_test.go` - Validation clés primaires
- [ ] `constraint/fact_validator_test.go` - Validateur de faits
- [ ] `constraint/type_system_test.go` - Système de types
- [ ] `constraint/internal_id_test.go` - IDs internes

**Checklist par fichier** :
```bash
# Pour chaque fichier :
1. Vérifier si utilise directement ID: "..."
2. Remplacer par Fields: {"id": "..."}
3. Vérifier messages d'erreur attendus
4. Exécuter : go test -v -run TestNomDuFichier
5. Vérifier couverture : go test -cover
6. Commit atomique : git commit -m "test: migrate <file> to new ID pattern"
```

### Phase 2 : Tests RETE (rete/)

**Stratégie** : Migrer par catégorie, pas par fichier individuel

#### Catégorie A : Tests de Base
- [ ] `rete/fact_token_test.go`
- [ ] `rete/field_resolver_test.go`
- [ ] `rete/comparison_evaluator_test.go`
- [ ] `rete/network_test.go`

#### Catégorie B : Tests Alpha
- [ ] `rete/node_alpha_test.go`
- [ ] `rete/alpha_chain_*_test.go` (tous)
- [ ] `rete/alpha_sharing_*_test.go` (tous)

#### Catégorie C : Tests Beta
- [ ] `rete/beta_chain_*_test.go` (tous)
- [ ] `rete/beta_sharing_*_test.go` (tous)
- [ ] `rete/node_join_*_test.go` (tous)

#### Catégorie D : Tests Backward Compatibility
- [ ] `rete/backward_compatibility_test.go`
- [ ] Investiguer problème TerminalNode.Memory.Tokens vide
- [ ] Documenter comportement attendu

### Phase 3 : Tests API et TSDIO

- [ ] `api/simple_test.go`
- [ ] `api/xuple_action_automatic_test.go`
- [ ] `api/xuplespace_e2e_test.go`
- [ ] `tsdio/api_test.go`

### Phase 4 : Validation et Nettoyage

- [ ] Exécuter `go test ./...` - tous les tests doivent passer
- [ ] Vérifier couverture : `go test ./... -cover` > 80% partout
- [ ] Exécuter `make validate`
- [ ] Exécuter `golangci-lint run`
- [ ] Supprimer code deprecated :
  - [ ] Supprimer `FieldNameIDLegacy` de `constraint/constraint_constants.go`
  - [ ] Supprimer `GenerateFactIDWithoutContext` de `constraint/id_generator.go`
  - [ ] Supprimer `valueToString` de `constraint/id_generator.go`
- [ ] Mettre à jour CHANGELOG.md
- [ ] Mettre à jour documentation

## 📊 Métriques de Succès

### Critères d'Acceptation

- [ ] **Tests** : 100% des tests passent (`make test-complete`)
- [ ] **Couverture** : > 80% dans tous les modules
- [ ] **Aucun hardcoding** : Recherche `ID:.*"[A-Za-z0-9]+"` retourne 0 résultats dans tests
- [ ] **Validation** : `make validate` passe sans erreur
- [ ] **Documentation** : Guide de migration créé et à jour
- [ ] **Performance** : Aucune régression (benchmarks)

### Commandes de Validation

```bash
# Vérifier que tous les tests passent
make test-complete

# Vérifier la couverture
go test ./constraint -cover
go test ./rete -cover
go test ./api -cover  
go test ./tsdio -cover

# Rechercher les patterns à migrer
grep -r 'ID:.*"' rete/ --include="*_test.go" | wc -l
# Résultat attendu : 0

# Validation complète
make validate

# Linting
golangci-lint run
```

## ⚠️ NOTES IMPORTANTES

### Ce qui NE doit PAS être fait

❌ **Ne pas** migrer les tests sans d'abord corriger l'architecture (Phase 0)
❌ **Ne pas** modifier le comportement des tests, seulement la façon de créer les faits
❌ **Ne pas** sacrifier la couverture de code pour accélérer la migration
❌ **Ne pas** merger des changements partiels - tout doit passer

### Gestion des Erreurs Attendues

Certains tests peuvent avoir des assertions sur les messages d'erreur qui incluent "id". Ces messages devront peut-être être mis à jour :

```go
// Avant
"champ de clé primaire 'id' manquant"

// Après (potentiellement)
"champ de clé primaire '_id_' manquant" // OU
"identifiant interne manquant"
```

## 🔄 Workflow de Développement

Pour chaque fichier de test migré :

1. **Créer une branche** : `git checkout -b test/migrate-<module>-<file>`
2. **Modifier le fichier** selon le nouveau pattern
3. **Tester** : `go test -v -run <TestName>`
4. **Vérifier couverture** : `go test -cover`
5. **Commit atomique** avec message clair
6. **Push et PR** pour review
7. **Merger** après validation

## 📅 Estimation Temporelle

| Phase | Durée Estimée | Priorité |
|-------|---------------|----------|
| Phase 0 : Préparation | 2-3 jours | CRITIQUE |
| Phase 1 : Tests Critiques | 3-5 jours | HAUTE |
| Phase 2 : Tests RETE | 10-15 jours | HAUTE |
| Phase 3 : Tests API/TSDIO | 2-3 jours | MOYENNE |
| Phase 4 : Validation | 2-3 jours | CRITIQUE |
| **TOTAL** | **19-29 jours** | - |

## 👥 Responsabilités

- **Architecte** : Approuver les changements de Phase 0
- **Dev Lead** : Review des PRs de migration
- **QA** : Validation des tests après migration
- **DevOps** : Mise à jour CI/CD si nécessaire

## 📚 Références

- Prompt d'origine : `/home/resinsec/dev/tsd/scripts/new_ids/07-prompt-tests-unit.md`
- Rapport de review : `/home/resinsec/dev/tsd/REPORTS/new_ids_review.md`
- Standards : `/home/resinsec/dev/tsd/.github/prompts/common.md`
- Architecture : `/home/resinsec/dev/tsd/docs/architecture/`

---

**Dernière mise à jour** : 2025-12-19
**Statut** : 🔴 BLOQUÉ - En attente de correction architecture (Phase 0)
