# Plan de Migration - Nouvelle Gestion des Identifiants

## 🎯 Objectif

Modifier la gestion des identifiants dans TSD pour :

1. **Identifiant interne caché** (`_id_`) jamais accessible dans les expressions TSD
2. **Génération automatique** obligatoire, plus d'affectation manuelle
3. **Comparaisons simplifiées** : `p.user == u` au lieu de `p.user == u.user`
4. **Types comme valeurs** : permettre `type Login(user: User, ...)` 
5. **Affectation de faits** : `a = User(...); Login(a, ...)`

---

## 📋 Vue d'Ensemble

### Changements Majeurs

| Aspect | Avant | Après |
|--------|-------|-------|
| **Champ ID** | `id` (visible) | `_id_` (caché) |
| **Définition manuelle** | Possible (backward compat) | ❌ Interdite |
| **Comparaison** | `p.user == u.user` | `p.user == u` |
| **Type de champ** | Types primitifs uniquement | Primitifs + Types faits |
| **Affectation** | Non supportée | `a = User(...)` |

### Modules Impactés

```
constraint/     ← Parser, types, validation, génération ID
rete/          ← Structures, évaluation, comparaisons
tsdio/         ← API publique
api/           ← Interface externe
tests/         ← Tous les tests
docs/          ← Documentation complète
```

---

## 📂 Structure du Plan

Les prompts sont numérotés pour exécution séquentielle :

```
scripts/new_ids/
├── README.md                          # Ce fichier
├── 00-prompt-analyse.md              # Analyse préliminaire
├── 01-prompt-structures-base.md      # Structures Fact et constantes
├── 02-prompt-parser-syntax.md        # Parser - Nouvelle syntaxe
├── 03-prompt-id-generation.md        # Génération d'IDs
├── 04-prompt-evaluation.md           # Évaluation et comparaisons
├── 05-prompt-types-validation.md     # Types et validation
├── 06-prompt-api-tsdio.md            # API et tsdio
├── 07-prompt-tests-unit.md           # Tests unitaires
├── 08-prompt-tests-integration.md    # Tests d'intégration
├── 09-prompt-tests-e2e.md            # Tests end-to-end
└── 10-prompt-documentation.md        # Documentation
```

---

## 🔄 Ordre d'Exécution

### Phase 1 : Préparation et Analyse
1. **00-prompt-analyse.md** - Comprendre l'état actuel

### Phase 2 : Core Structures
2. **01-prompt-structures-base.md** - Modifier Fact, FieldNameID → _id_
3. **02-prompt-parser-syntax.md** - Parser pour nouvelle syntaxe
4. **03-prompt-id-generation.md** - Génération automatique d'IDs

### Phase 3 : Logique Métier
5. **04-prompt-evaluation.md** - Comparaisons et évaluation
6. **05-prompt-types-validation.md** - Types complexes et validation

### Phase 4 : API et Interface
7. **06-prompt-api-tsdio.md** - Adapter API publique

### Phase 5 : Tests
8. **07-prompt-tests-unit.md** - Tests unitaires
9. **08-prompt-tests-integration.md** - Tests d'intégration  
10. **09-prompt-tests-e2e.md** - Tests end-to-end

### Phase 6 : Documentation
11. **10-prompt-documentation.md** - Mise à jour complète docs/

---

## ✅ Checklist par Prompt

Chaque prompt DOIT :

- [ ] Respecter strictement `.github/prompts/common.md`
- [ ] Respecter strictement `.github/prompts/develop.md`
- [ ] Être auto-suffisant (contexte complet)
- [ ] Inclure des tests (couverture > 80%)
- [ ] Ne pas dépasser 128k tokens de contexte
- [ ] Être exécutable automatiquement
- [ ] Inclure validation (`make validate`)
- [ ] Mettre à jour CHANGELOG si pertinent

---

## 🎯 Principes de Migration

### 1. Pas de Hardcoding
- ❌ Aucune valeur en dur
- ✅ Constantes nommées (`FieldNameInternalID = "_id_"`)

### 2. Généricité
- ❌ Cas spécifiques codés en dur
- ✅ Code paramétré et extensible

### 3. Tests d'Abord
- Écrire tests AVANT implémentation
- Couverture > 80% obligatoire

### 4. Documentation Synchrone
- Mettre à jour docs/ au fur et à mesure
- Supprimer docs obsolètes

### 5. Validation Continue
- `make validate` après chaque prompt
- `make test-complete` avant passage au suivant

---

## 🚀 Workflow d'Exécution

### Pour chaque prompt :

1. **Lire le prompt** complet
2. **Créer une branche** si nécessaire
3. **Exécuter** les modifications
4. **Valider** :
   ```bash
   make format
   make validate
   make test-complete
   ```
5. **Commit** avec message descriptif
6. **Passer** au prompt suivant

### Commandes Utiles

```bash
# Validation complète
make validate

# Tests complets
make test-complete

# Couverture
make test-coverage

# Format
make format

# Vérifier diagnostics
make lint
```

---

## 📝 Convention de Commit

```
feat(ids): [Prompt XX] Description courte

- Détail 1
- Détail 2
- Détail 3

Refs: scripts/new_ids/XX-prompt-name.md
```

**Exemple** :
```
feat(ids): [Prompt 01] Modifier structures de base

- Renommer FieldNameID → FieldNameInternalID
- Valeur "id" → "_id_"
- Interdire champ _id_ dans définitions de faits
- Tests associés

Refs: scripts/new_ids/01-prompt-structures-base.md
```

---

## ⚠️ Points d'Attention

### Compatibilité Ascendante
- ❌ **Pas de rétrocompatibilité** - Breaking change assumé
- Documentation de migration pour utilisateurs
- Exemples avant/après

### Performance
- Pas de dégradation de performance
- Benchmarks si optimisations

### Sécurité
- Validation stricte des entrées
- Pas d'injection possible
- Erreurs informatives sans fuites

---

## 📚 Références

### Documents Projet
- `.github/prompts/common.md` - Standards communs
- `.github/prompts/develop.md` - Standards développement
- `docs/ID_RULES_COMPLETE.md` - Règles actuelles des IDs
- `docs/primary-keys.md` - Documentation clés primaires

### Code Critique
- `constraint/constraint_types.go` - Structures de base
- `constraint/parser.go` - Parser (généré)
- `constraint/id_generator.go` - Génération d'IDs
- `rete/fact_token.go` - Structures RETE
- `tsdio/api.go` - API publique

### Tests de Référence
- `constraint/id_generator_test.go`
- `constraint/primary_key_validation_test.go`
- `tests/e2e/`

---

## 🎓 Formation

### Avant de Commencer

Lire attentivement :
1. `.github/prompts/common.md` (OBLIGATOIRE)
2. `.github/prompts/develop.md` (OBLIGATOIRE)
3. `docs/ID_RULES_COMPLETE.md`
4. `docs/primary-keys.md`

### Comprendre le Contexte

- Architecture RETE
- Système de contraintes
- Génération d'IDs actuelle
- Types et validation

---

## ✨ Résultat Attendu

### Nouvelle Syntaxe TSD

```tsd
type User(#name: string, #firstname: string, age: number)
type Login(user: User, #password: string, #email: string)

// Définition avec affectation
jean = User("Dupont", "Jean", 23)
amelie = User("Poulain", "Amélie", 19)
Login(amelie, "pass123", "ap@gmail.com")

// Règle avec comparaison simplifiée
{l: Login, u: User} / l.user == u ==> 
  Log("Compte " + l.email + " pour " + u.name + " " + u.firstname)
```

### Comportement Interne

- `_id_` calculé automatiquement (jamais visible)
- Comparaisons `p.user == u` résolues via `_id_`
- Affectations de faits stockées et référencées
- Validation stricte des types

---

## 📊 Métriques de Succès

- [ ] Tous les tests passent (`make test-complete`)
- [ ] Couverture > 80% maintenue
- [ ] `make validate` sans erreurs
- [ ] Documentation à jour dans `docs/`
- [ ] Exemples fonctionnels dans `examples/`
- [ ] Aucun hardcoding
- [ ] Code générique et extensible

---

**Branche** : `feature/new-id-management`

**Statut** : 🚧 En cours

**Dernière mise à jour** : 2025-01-XX