# 📚 Prompts TSD - Guide d'Utilisation

## 🎯 Vue d'Ensemble

Ce répertoire contient les **prompts réutilisables** pour le développement du projet TSD.

Le système a été **simplifié** : au lieu de 19+ prompts spécifiques, nous utilisons maintenant **6 prompts universels** couvrant tous les besoins.

---

## 📋 Prompts Disponibles

### 🔧 [develop.md](./develop.md)
**Développement de code**
- Ajouter une fonctionnalité
- Modifier un comportement existant
- Corriger un bug

**Utilisation** :
```
Utilise le prompt "develop" pour ajouter une fonctionnalité de validation
Utilise le prompt "develop" pour corriger le bug dans le module X
```

---

### 🧪 [test.md](./test.md)
**Gestion des tests**
- Écrire des tests (unitaires, intégration, e2e)
- Exécuter et analyser les tests
- Déboguer un test qui échoue
- Analyser la couverture

**Utilisation** :
```
Utilise le prompt "test" pour écrire des tests pour la fonction Y
Utilise le prompt "test" pour déboguer le test TestProblematic
```

---

### 📚 [document.md](./document.md)
**Documentation**
- Écrire/mettre à jour la documentation
- Expliquer du code (niveaux débutant à expert)
- Générer des exemples (.tsd, code)
- README, GoDoc, guides techniques

**Utilisation** :
```
Utilise le prompt "document" pour expliquer le module RETE
Utilise le prompt "document" pour générer des exemples .tsd
```

---

### 🔍 [review.md](./review.md)
**Revue et qualité**
- Code review complète
- Refactoring (sans changer le comportement)
- Amélioration de la qualité

**Utilisation** :
```
Utilise le prompt "review" pour faire une code review de node_join.go
Utilise le prompt "review" pour refactoriser la fonction complexe
```

---

### 🔬 [analyze.md](./analyze.md)
**Analyse et diagnostic**
- Analyser une erreur
- Investiguer un comportement inattendu
- Valider un réseau RETE
- Analyser la performance (profiling)

**Utilisation** :
```
Utilise le prompt "analyze" pour analyser l'erreur "variable non liée"
Utilise le prompt "analyze" pour valider le réseau RETE test.tsd
```

---

### 🔧 [maintain.md](./maintain.md)
**Maintenance du projet**
- Migration (Go version, dépendances, API)
- Nettoyage (code mort, docs obsolètes)
- Vérification des licences
- Statistiques du code
- Optimisation de la performance

**Utilisation** :
```
Utilise le prompt "maintain" pour migrer vers Go 1.21
Utilise le prompt "maintain" pour générer les stats du projet
Utilise le prompt "maintain" pour vérifier les licences
```

---

## ⭐ Document de Référence

### [common.md](./common.md)
**Standards communs du projet - À TOUJOURS CONSULTER**

Ce document contient TOUS les standards du projet :
- 🔒 Licence et Copyright (obligatoire)
- ⚠️ Règles strictes - Code Go (interdictions, bonnes pratiques)
- 🧪 Standards de tests (structure, couverture > 80%)
- 📚 Documentation (organisation, langues, formats)
- 🔧 Outils et commandes (validation, profiling)
- 🎨 Conventions de nommage
- 📋 Checklist avant commit
- 🚀 Workflow de développement

**Tous les prompts référencent common.md** - C'est la source unique de vérité pour les standards du projet.

---

## 🚀 Démarrage Rapide

### 1. Première Utilisation

```bash
# Lire les standards du projet (OBLIGATOIRE)
cat .github/prompts/common.md

# Parcourir l'INDEX
cat .github/prompts/INDEX.md
```

### 2. Utilisation Quotidienne

```bash
# Format d'utilisation
"Utilise le prompt '[nom]' pour [action]"

# Exemples
"Utilise le prompt 'develop' pour ajouter une fonctionnalité de cache"
"Utilise le prompt 'test' pour écrire des tests pour le cache"
"Utilise le prompt 'review' pour faire une code review"
```

### 3. Workflow Type

1. **Consulter** [common.md](./common.md) pour les standards
2. **Choisir** le prompt adapté à votre besoin
3. **Préciser** votre demande clairement
4. **Suivre** les instructions du prompt
5. **Valider** avec la checklist de [common.md](./common.md)

---

## 🎯 Avantages du Système

### ✅ Simplicité
- 6 prompts universels au lieu de 19 spécifiques
- Un prompt par catégorie d'action
- Pas d'hésitation sur lequel utiliser

### ✅ Cohérence
- Tous référencent [common.md](./common.md)
- Standards unifiés et à jour
- Aucune redondance

### ✅ Maintenabilité
- Mise à jour centralisée dans [common.md](./common.md)
- Facile à faire évoluer
- Une seule source de vérité

### ✅ Exhaustivité
- Tous les besoins couverts
- Générique et adaptable
- Extensible facilement

---

## 📊 Comparaison Avant/Après

| Aspect | Avant | Après |
|--------|-------|-------|
| **Nombre de prompts** | 19 prompts spécifiques | 6 prompts universels |
| **Taille totale** | ~260 Ko | ~80 Ko |
| **Redondances** | Nombreuses | Aucune |
| **Maintenance** | Difficile (19 fichiers) | Simple (6 fichiers) |
| **Standards** | Dispersés | Centralisés dans common.md |
| **Apprentissage** | Long (19 prompts) | Rapide (6 prompts) |

---

## 🔍 Navigation

- **[INDEX.md](./INDEX.md)** - Navigation complète et exemples
- **[common.md](./common.md)** - Standards du projet (RÉFÉRENCE)
- **Prompts** - develop, test, document, review, analyze, maintain

---

## 📚 Ressources Externes

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

## 📝 Notes

### Anciens Prompts
Les 19 anciens prompts ont été **archivés** dans `REPORTS/prompts-optimization/old-prompts/` pour référence historique.

### Évolution
Le système de prompts est conçu pour évoluer. Si un besoin spécifique émerge qui ne peut pas être couvert par les 6 prompts universels, un nouveau prompt peut être ajouté en suivant le même principe : référencer [common.md](./common.md) et rester générique.

---

**Version** : 2.0  
**Date** : Décembre 2024  
**Statut** : ✅ Système simplifié et optimisé