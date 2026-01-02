# 📋 TODO - Nettoyage Documentation

Suite au refactoring effectué (voir REFACTORING_CLEANUP_REPORT.md), les fichiers de documentation suivants contiennent encore des références à l'ancien pattern factory qui a été supprimé.

---

## 🔴 URGENT - Documentation Obsolète

### Fichiers à Mettre à Jour

1. **XUPLES_E2E_AUTOMATIC.md**
   - [ ] Supprimer toutes les sections mentionnant `SetXupleSpaceFactory()`
   - [ ] Supprimer les exemples de configuration manuelle de la factory
   - [ ] Mettre à jour les exemples pour montrer l'automatisation complète
   - [ ] Ajouter une note de migration depuis v1.x

2. **MIGRATION_E2E_RESUME.md**
   - [ ] Supprimer les références à `XupleSpaceFactory`
   - [ ] Mettre à jour les exemples de code
   - [ ] Documenter le nouveau flow avec callback (si pertinent pour les utilisateurs externes)

3. **scripts/02-package-api-pipeline.md**
   - [ ] Mettre à jour les diagrammes d'architecture
   - [ ] Supprimer les sections sur la configuration de la factory
   - [ ] Documenter le callback pattern (usage interne uniquement)

4. **scripts/03-creation-auto-xuplespaces.md**
   - [ ] Mettre à jour pour refléter l'implémentation finale
   - [ ] Supprimer les workarounds temporaires mentionnés
   - [ ] Documenter le flow final

5. **RAPPORT_API_PACKAGE.md**
   - [ ] Mettre à jour les exemples de code
   - [ ] Supprimer les références à l'ancien pattern

6. **RAPPORT_REFACTORING_XUPLE_ACTION.md**
   - [ ] Vérifier et mettre à jour si nécessaire
   - [ ] Marquer comme obsolète si totalement dépassé

---

## 🟡 RECOMMANDÉ - Nouvelles Documentation

### À Créer

1. **MIGRATION_GUIDE_v2.md**
   - [ ] Guide de migration v1.x → v2.0
   - [ ] Changements incompatibles (breaking changes)
   - [ ] Exemples avant/après
   - [ ] FAQ

2. **API_REFERENCE.md**
   - [ ] Documentation complète de l'API publique du package `api`
   - [ ] Tous les types, fonctions, méthodes exportés
   - [ ] Exemples d'utilisation pour chaque fonction principale

3. **ARCHITECTURE.md** (mise à jour)
   - [ ] Diagramme d'architecture mis à jour
   - [ ] Flow de données explicite
   - [ ] Responsabilités de chaque package

---

## 🟢 OPTIONNEL - Amélioration Continue

### Exemples et Tutoriels

1. **examples/01-basic-xuples/main.go**
   - [ ] Exemple minimal utilisant le package `api`
   - [ ] Démonstration de création automatique de xuple-spaces

2. **examples/02-custom-policies/main.go**
   - [ ] Exemple avec policies personnalisées
   - [ ] Configuration avancée

3. **examples/03-multiple-spaces/main.go**
   - [ ] Exemple avec plusieurs xuple-spaces
   - [ ] Utilisation avancée

### Tests de Documentation

1. **Créer des tests "example"**
   - [ ] Tests qui servent de documentation exécutable
   - [ ] Apparaissent dans GoDoc

---

## 📝 Commandes de Recherche

Pour trouver toutes les occurrences restantes de l'ancien pattern :

```bash
# Rechercher dans la documentation
grep -r "SetXupleSpaceFactory\|XupleSpaceFactory\|factory pluggable" --include="*.md" .

# Rechercher dans le code
grep -r "SetXupleSpaceFactory\|XupleSpaceFactory" --include="*.go" .

# Rechercher "factory" dans les commentaires
grep -r "// .*factory\|/\* .*factory" --include="*.go" .
```

---

## ✅ Critères de Complétion

La tâche sera considérée comme terminée quand :

- [ ] Aucune référence à `SetXupleSpaceFactory` dans la documentation
- [ ] Aucune référence à `XupleSpaceFactory` dans la documentation (sauf dans les notes de migration)
- [ ] Tous les exemples de code utilisent la nouvelle API
- [ ] Un guide de migration v1.x → v2.0 existe
- [ ] Les diagrammes d'architecture sont à jour
- [ ] La documentation API est complète

---

## 🚀 Priorité d'Exécution

1. **URGENT** : Mettre à jour `XUPLES_E2E_AUTOMATIC.md` (documentation principale)
2. **IMPORTANT** : Créer `MIGRATION_GUIDE_v2.md`
3. **RECOMMANDÉ** : Mettre à jour les autres fichiers de documentation
4. **OPTIONNEL** : Créer de nouveaux exemples

---

**Créé**: 2025-12-18  
**Contexte**: Après refactoring pour suppression du pattern factory obsolète
