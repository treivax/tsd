# Prompt 00 - Analyse Préliminaire

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Analyser en profondeur l'état actuel de la gestion des identifiants dans TSD pour préparer la migration vers le nouveau système avec `_id_` interne.

Cette analyse servira de base pour tous les prompts suivants.

---

## 📋 Contexte

### Modification Demandée

Nous devons transformer la gestion des identifiants de faits :

| Aspect | Actuellement | Cible |
|--------|--------------|-------|
| **Nom du champ** | `id` (visible) | `_id_` (caché) |
| **Affectation manuelle** | Possible (backward compat) | ❌ Interdite |
| **Visibilité** | Accessible dans expressions | ❌ Jamais visible |
| **Comparaisons** | `p.user == u.user` | `p.user == u` |
| **Types de champs** | Primitifs uniquement | Primitifs + Faits |
| **Affectation de faits** | Non supportée | `a = User(...)` |

---

## 📝 Tâches à Réaliser

### 1. Cartographie des Structures

#### 1.1 Identifier toutes les structures `Fact`

Chercher et documenter :

```bash
# Structures Fact dans le projet
grep -r "type Fact struct" --include="*.go"
```

**Pour chaque structure `Fact` trouvée** :
- Localisation (fichier, package)
- Champs actuels
- Utilisation du champ `id` ou `ID`
- Dépendances (qui l'utilise)

**Fichiers attendus** :
- `constraint/constraint_types.go` - Structure parser
- `rete/fact_token.go` - Structure RETE
- `tsdio/api.go` - Structure API
- `api/result.go` - Structure résultats

#### 1.2 Identifier les constantes d'ID

```bash
# Rechercher les constantes liées aux IDs
grep -r "FieldNameID\|FieldName.*ID" --include="*.go"
```

**Documenter** :
- Nom de la constante
- Valeur actuelle
- Où elle est définie
- Où elle est utilisée

**Fichier attendu** :
- `constraint/constraint_constants.go`

### 2. Analyse du Parser

#### 2.1 Grammaire PEG

Lire et analyser :
- `constraint/grammar/constraint.peg` - Grammaire actuelle

**Identifier** :
- Règles pour définition de types
- Règles pour définition de faits
- Règles pour field access
- Règles pour comparaisons

**Questions à répondre** :
1. Comment sont parsés les types actuellement ?
2. Comment sont parsées les définitions de faits ?
3. Comment sont parsés les field access (`p.user`) ?
4. Y a-t-il déjà support pour affectation de variables ?

#### 2.2 Code Parser

Analyser :
- `constraint/parser.go` - Parser généré (NE PAS MODIFIER)
- `constraint/api.go` - API du parser

**Identifier** :
- Fonctions de parsing des types
- Fonctions de parsing des faits
- Fonctions de parsing des expressions
- Points d'entrée principaux

### 3. Analyse de la Génération d'IDs

#### 3.1 Algorithme Actuel

Lire et documenter :
- `constraint/id_generator.go`

**Fonctions critiques** :
- `GenerateFactID(fact Fact, typeDef TypeDefinition)`
- `generateIDFromPrimaryKey()`
- `generateIDFromHash()`
- `escapeIDValue()`

**Format actuel des IDs** :
- Avec clé primaire simple : `TypeName~value`
- Avec clé composite : `TypeName~value1_value2`
- Sans clé primaire : `TypeName~<hash>`

#### 3.2 Utilisation des IDs

Chercher où les IDs sont :
- Générés
- Stockés
- Comparés
- Affichés

```bash
grep -r "GenerateFactID\|ensureFactID" --include="*.go"
grep -r "FieldNameID" --include="*.go"
grep -r '\.ID\s*=\|\.id\s*=' --include="*.go"
```

### 4. Analyse des Comparaisons

#### 4.1 Évaluation des Contraintes

Analyser :
- `rete/` - Package RETE
- Rechercher évaluation de comparaisons
- Rechercher field access

**Identifier** :
- Comment `p.user == u.user` est évalué
- Comment les champs sont résolus
- Où les types sont vérifiés

#### 4.2 Type Checking

Analyser :
- `constraint/constraint_type_checking.go`
- `constraint/constraint_field_validation.go`

**Questions** :
1. Comment les types de champs sont validés ?
2. Comment les comparaisons sont type-checkées ?
3. Où ajouter le support pour types de faits ?

### 5. Analyse de la Validation

#### 5.1 Validation des Faits

Lire :
- `constraint/constraint_facts.go`
- `constraint/primary_key_validation.go`

**Identifier** :
- Validation des champs de faits
- Validation des clés primaires
- Interdictions actuelles (champ `id` manuel)

#### 5.2 Validation des Types

Lire :
- `constraint/constraint_type_validation.go`

**Identifier** :
- Validation des définitions de types
- Types de champs autorisés
- Restrictions actuelles

### 6. Analyse des Tests

#### 6.1 Tests Critiques

Lister les fichiers de tests critiques :

```bash
find constraint/ -name "*test.go" | grep -i "id\|primary\|fact"
find rete/ -name "*test.go" | grep -i "fact\|token"
find tests/e2e/ -name "*.tsd"
```

**Catégoriser** :
- Tests unitaires de génération d'IDs
- Tests de validation de faits
- Tests de clés primaires
- Tests d'intégration
- Tests E2E

#### 6.2 Exemples TSD

Lister :
```bash
find examples/ -name "*.tsd"
find tests/ -name "*.tsd"
```

**Identifier** :
- Exemples utilisant clés primaires
- Exemples utilisant comparaisons
- Exemples à migrer

### 7. Analyse de la Documentation

#### 7.1 Documentation Technique

Lire :
- `docs/ID_RULES_COMPLETE.md`
- `docs/primary-keys.md`
- `docs/MIGRATION_IDS.md`

**Extraire** :
- Règles actuelles complètes
- Exemples documentés
- Cas d'usage

#### 7.2 Documentation API

Lire :
- `docs/api/`
- `docs/architecture/`

**Identifier** :
- Documentation à mettre à jour
- Diagrammes à modifier
- Exemples à changer

---

## 📊 Livrables

### Rapport d'Analyse

Créer : `REPORTS/new_ids_analysis.md`

**Structure attendue** :

```markdown
# Analyse Préliminaire - Migration des IDs

## 1. Cartographie des Structures

### Structures Fact
- constraint/constraint_types.go
  - Champs : Type, TypeName, Fields []FactField
  - Pas de champ ID direct
  
- rete/fact_token.go
  - Champs : ID string, Type string, Fields map[string]interface{}
  - ID est explicite
  
[...]

### Constantes
- FieldNameID = "id" (constraint/constraint_constants.go)
- [...]

## 2. Parser et Grammaire

### Grammaire PEG
- TypeDefinition : [détails]
- FactDefinition : [détails]
- FieldAccess : [détails]

### Points d'Extension
- [Où ajouter support pour User type dans fields]
- [Où ajouter support pour affectation]

## 3. Génération d'IDs

### Algorithme Actuel
- [Description détaillée]
- Fonctions : GenerateFactID, generateIDFromPrimaryKey, etc.

### Points de Modification
- [Liste des endroits à changer]

## 4. Comparaisons et Évaluation

### Évaluation Actuelle
- [Comment p.user == u.user fonctionne]

### Modification Requise
- [Comment implémenter p.user == u]
- [Résolution de type nécessaire]

## 5. Validation

### Validations Actuelles
- [Liste complète]

### Nouvelles Validations
- Interdire _id_ dans définitions
- Valider types de faits dans champs
- [...]

## 6. Tests Impactés

### Tests à Migrer
- [Liste par catégorie]
- constraint/ : XX fichiers
- rete/ : XX fichiers
- tests/e2e/ : XX fichiers

### Nouveaux Tests Requis
- Test affectation de faits
- Test comparaison p.user == u
- [...]

## 7. Documentation

### Fichiers à Mettre à Jour
- docs/ID_RULES_COMPLETE.md
- docs/primary-keys.md
- [...]

### Nouveaux Exemples
- [Liste]
```

### Matrice d'Impact

Créer : `REPORTS/new_ids_impact_matrix.md`

**Format** :

| Module | Fichier | Impact | Complexité | Priorité | Notes |
|--------|---------|--------|------------|----------|-------|
| constraint | constraint_types.go | ⚠️ Majeur | Élevée | 1 | Structures de base |
| constraint | id_generator.go | ⚠️ Majeur | Moyenne | 2 | Génération |
| constraint | parser.go | ⚠️ Majeur | Élevée | 1 | Parser (généré) |
| [...] | [...] | [...] | [...] | [...] | [...] |

**Légende Impact** :
- 🔴 Critique - Réécriture complète
- ⚠️ Majeur - Modifications importantes
- 🟡 Modéré - Quelques changements
- 🟢 Mineur - Ajustements légers
- ⚪ Aucun - Non impacté

**Légende Complexité** :
- Élevée (> 3j)
- Moyenne (1-3j)
- Faible (< 1j)

### Checklist de Décisions

Créer : `REPORTS/new_ids_decisions.md`

**Questions à trancher** :

```markdown
# Décisions Architecture - Migration IDs

## 1. Nom du Champ Interne

- [ ] Option A : `_id_` (recommandé)
- [ ] Option B : `__id__`
- [ ] Option C : autre ?

**Décision** : `_id_`
**Raison** : [...]

## 2. Constante pour le Champ

- [ ] `FieldNameInternalID = "_id_"`
- [ ] Garder `FieldNameID` mais changer valeur
- [ ] Autre ?

**Décision** : [...]
**Raison** : [...]

## 3. Validation du Champ _id_

Comment interdire _id_ dans les définitions ?

- [ ] Dans le parser (grammaire PEG)
- [ ] Dans la validation post-parsing
- [ ] Les deux (défense en profondeur)

**Décision** : [...]
**Raison** : [...]

[...]
```

---

## ✅ Critères de Succès

- [ ] Rapport d'analyse complet généré
- [ ] Matrice d'impact créée
- [ ] Checklist de décisions complétée
- [ ] Tous les fichiers critiques identifiés
- [ ] Toutes les dépendances documentées
- [ ] Points de modification listés
- [ ] Complexité estimée par module
- [ ] Ordre d'exécution des prompts validé

---

## 🚀 Exécution

### Commandes

```bash
# Se positionner sur la branche
cd tsd
git checkout feature/new-id-management

# Créer le répertoire REPORTS si nécessaire
mkdir -p REPORTS

# Lancer l'analyse (manuelle, via exploration)
# Utiliser grep, find, read_file pour explorer

# Générer les rapports
# Créer les fichiers markdown listés dans Livrables
```

### Validation

```bash
# Vérifier que les rapports sont créés
ls -la REPORTS/new_ids_*

# Vérifier le contenu
cat REPORTS/new_ids_analysis.md
cat REPORTS/new_ids_impact_matrix.md
cat REPORTS/new_ids_decisions.md
```

---

## 📚 Références

- `.github/prompts/common.md` - Standards du projet
- `.github/prompts/develop.md` - Standards de développement
- `docs/ID_RULES_COMPLETE.md` - Règles actuelles
- `docs/primary-keys.md` - Documentation clés primaires
- `scripts/new_ids/README.md` - Vue d'ensemble du plan

---

## 📝 Notes

### Points d'Attention

1. **Parser généré** : `constraint/parser.go` est généré depuis `constraint/grammar/constraint.peg`. Ne JAMAIS modifier parser.go directement.

2. **Backward compatibility** : Cette migration est un breaking change assumé. Pas de rétrocompatibilité à maintenir.

3. **Performance** : Attention à ne pas dégrader les performances lors des comparaisons de faits.

4. **Sécurité** : Le champ `_id_` ne doit jamais pouvoir être défini par l'utilisateur (risque d'injection/collision).

### Questions Ouvertes

À discuter/trancher lors de l'analyse :

1. Faut-il maintenir un alias `id` en lecture seule pour debug ?
2. Comment gérer la sérialisation JSON (cacher `_id_` ?) ?
3. Les comparaisons `p.user == u` nécessitent-elles un cache de résolution ?
4. Comment optimiser les jointures sur types de faits ?

---

**Prompt suivant** : `01-prompt-structures-base.md`

**Durée estimée** : 2-4 heures

**Complexité** : 🟡 Moyenne (analyse, pas de code)