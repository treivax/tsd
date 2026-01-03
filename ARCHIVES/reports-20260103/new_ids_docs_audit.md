# Audit Documentation - Migration IDs

**Date**: 2025-12-19  
**Objectif**: Audit complet de la documentation pour la migration des identifiants

---

## 📊 État Actuel

### Statistiques
- **Fichiers markdown totaux**: 82
- **Références à l'ancienne syntaxe**: ~111 occurrences de `id:` dans docs/
- **Documentation à créer**: 6+ nouveaux documents
- **Documentation à archiver**: 3 documents obsolètes
- **Documentation à mettre à jour**: ~15 fichiers

---

## 📋 Documentation à Mettre à Jour

### Priorité 1 - Critique (Breaking Changes)

- [ ] **docs/ID_RULES_COMPLETE.md** 
  - **État**: Obsolète - Décrit l'ancien système avec `id` accessible
  - **Action**: Remplacer par `docs/internal-ids.md`
  - **Raison**: Le champ est maintenant `_id_` et caché

- [ ] **docs/MIGRATION_IDS.md**
  - **État**: Incomplet - Ne couvre pas tous les breaking changes
  - **Action**: Remplacer par `docs/migration/from-v1.x.md`
  - **Raison**: Besoin d'un guide complet avec affectations et comparaisons

- [ ] **docs/primary-keys.md**
  - **État**: Partiellement correct mais en anglais
  - **Action**: Intégrer dans `docs/internal-ids.md` (français)
  - **Raison**: Consolidation de la doc sur les IDs

- [ ] **README.md (racine)**
  - **État**: Ne mentionne pas les nouvelles fonctionnalités v2.0
  - **Action**: Ajouter section "Nouveautés v2.0" avec exemples
  - **Raison**: Point d'entrée principal du projet

### Priorité 2 - Important (Guides Utilisateur)

- [ ] **docs/guides.md**
  - **État**: Exemples avec ancienne syntaxe
  - **Action**: Mettre à jour tous les exemples
  - **Impact**: Utilisateurs suivront des exemples obsolètes

- [ ] **docs/reference.md**
  - **État**: Référence incomplète (pas d'affectations)
  - **Action**: Ajouter sections affectations et comparaisons
  - **Impact**: Documentation de référence incomplète

- [ ] **constraint/README.md**
  - **État**: À vérifier
  - **Action**: Mettre à jour si exemples obsolètes

- [ ] **rete/README.md**
  - **État**: À vérifier
  - **Action**: Vérifier cohérence avec nouveau système

### Priorité 3 - Secondaire (Architecture et Tutoriels)

- [ ] **docs/architecture.md**
  - **Action**: Ajouter diagramme génération d'IDs

- [ ] **docs/tutorials/**
  - **Action**: Créer tutoriels avec nouvelle syntaxe

---

## 📝 Documentation à Créer

### Nouveaux Documents Essentiels

1. **docs/internal-ids.md** (NOUVEAU)
   - Remplace `ID_RULES_COMPLETE.md`
   - Documentation complète du système `_id_`
   - Explication que `_id_` est caché et automatique
   - ~400 lignes

2. **docs/user-guide/fact-assignments.md** (NOUVEAU)
   - Guide complet des affectations
   - Syntaxe `variable = Type(...)`
   - Exemples et bonnes pratiques
   - ~350 lignes

3. **docs/user-guide/fact-comparisons.md** (NOUVEAU)
   - Guide des comparaisons de faits
   - Syntaxe `fact1 == fact2`
   - Comment ça fonctionne en interne
   - ~300 lignes

4. **docs/user-guide/type-system.md** (NOUVEAU)
   - Système de types complet
   - Types de faits dans les champs
   - ~250 lignes

5. **docs/migration/from-v1.x.md** (NOUVEAU)
   - **CRITIQUE** - Guide de migration complet
   - Breaking changes détaillés
   - Exemples avant/après
   - Script de vérification
   - ~600 lignes

6. **docs/README.md** (MISE À JOUR MAJEURE)
   - Index complet de la documentation
   - Navigation par niveau (débutant/intermédiaire/avancé)
   - ~200 lignes

---

## 🗑️ Documentation à Archiver

### Fichiers Obsolètes

1. **docs/ID_RULES_COMPLETE.md**
   - Raison: Décrit système obsolète (`id` au lieu de `_id_`)
   - Destination: `docs/archive/pre-v2.0/`

2. **docs/MIGRATION_IDS.md** (ancien)
   - Raison: Remplacé par guide complet
   - Destination: `docs/archive/pre-v2.0/`

3. **docs/primary-keys.md** (si en anglais)
   - Raison: Contenu intégré dans `internal-ids.md`
   - Destination: `docs/archive/pre-v2.0/`

---

## 🔍 Analyse des Problèmes Détectés

### Incohérences Majeures

1. **Champ `id` vs `_id_`**
   - Ancienne doc: `p.id` accessible
   - Nouveau système: `_id_` caché, jamais accessible
   - Impact: **BREAKING CHANGE majeur**

2. **Affectations manquantes**
   - Ancienne doc: Pas d'affectations de variables
   - Nouveau système: `alice = User(...)`
   - Impact: Fonctionnalité non documentée

3. **Comparaisons de faits**
   - Ancienne doc: Pas de comparaisons directes
   - Nouveau système: `fact1 == fact2`
   - Impact: Fonctionnalité non documentée

4. **Types de faits dans champs**
   - Ancienne doc: Seulement primitifs
   - Nouveau système: `Order(customer: Customer, ...)`
   - Impact: Fonctionnalité majeure non documentée

### Exemples Obsolètes

Exemples trouvés avec ancienne syntaxe:
```tsd
// ❌ OBSOLÈTE
type Person(name: string, age: number)
assert Person(id: "person_1", name: "Alice", age: 30)
{p: Person} / p.id == "person_1" ==> ...
```

Devrait être:
```tsd
// ✅ NOUVEAU
type Person(#name: string, age: number)
alice = Person("Alice", 30)
{p: Person} / p.name == "Alice" ==> ...
```

---

## 📈 Métriques

### Effort Estimé

| Tâche | Lignes | Temps |
|-------|--------|-------|
| Créer `internal-ids.md` | 400 | 2h |
| Créer guides utilisateur (3 fichiers) | 900 | 3h |
| Créer guide migration | 600 | 2h |
| Mettre à jour README principal | 200 | 1h |
| Créer index docs | 200 | 1h |
| Archiver obsolète | 50 | 0.5h |
| Vérifier liens | - | 0.5h |
| **TOTAL** | **~2350** | **~10h** |

### Priorités d'Exécution

1. ✅ Créer `docs/migration/from-v1.x.md` (CRITIQUE)
2. ✅ Créer `docs/internal-ids.md`
3. ✅ Créer guides utilisateur
4. ✅ Mettre à jour README.md
5. ✅ Créer index documentation
6. ✅ Archiver obsolète
7. ✅ Validation finale

---

## ✅ Critères de Validation

### Checklist Complète

- [ ] Tous les fichiers créés existent
- [ ] Aucune référence à `.id` (sauf dans migration)
- [ ] Aucune référence à `id:` dans les assertions
- [ ] Tous les exemples utilisent `_id_` (caché)
- [ ] Guide de migration complet et testé
- [ ] Tous les liens internes fonctionnent
- [ ] Documentation cohérente (français)
- [ ] Archives créées proprement

### Commandes de Vérification

```bash
# Vérifier absence de références obsolètes
grep -r '"id":' docs/ --include="*.md" | grep -v archive | grep -v _id_

# Vérifier présence de tous les nouveaux fichiers
ls -la docs/internal-ids.md
ls -la docs/user-guide/fact-assignments.md
ls -la docs/user-guide/fact-comparisons.md
ls -la docs/migration/from-v1.x.md
ls -la docs/README.md

# Vérifier archives
ls -la docs/archive/pre-v2.0/
```

---

## 📋 Plan d'Action Détaillé

### Phase 1: Préparation (30 min)
1. Créer structure de répertoires
2. Sauvegarder documentation actuelle
3. Créer ce rapport

### Phase 2: Documentation Critique (3h)
1. Créer `docs/internal-ids.md`
2. Créer `docs/migration/from-v1.x.md`
3. Créer guides utilisateur

### Phase 3: Mise à Jour (2h)
1. Mettre à jour README.md
2. Créer index documentation
3. Mettre à jour docs existantes

### Phase 4: Nettoyage (1h)
1. Archiver documentation obsolète
2. Vérifier tous les liens
3. Validation finale

---

## 🎯 Résultat Attendu

Documentation complète et cohérente pour TSD v2.0 avec:

✅ Nouveaux utilisateurs: Guide de démarrage clair  
✅ Migration v1.x: Guide détaillé avec exemples  
✅ Référence: Documentation technique complète  
✅ Cohérence: Aucune contradiction  
✅ Exemples: Tous à jour avec nouvelle syntaxe  

---

**Statut**: ⏳ EN COURS  
**Prochaine étape**: Créer structure de répertoires et commencer Phase 2
