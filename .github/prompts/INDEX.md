# 📇 INDEX - Prompts TSD

> Navigation rapide vers tous les prompts disponibles

---

## 🎯 Prompts Universels

Le projet TSD utilise **5 prompts génériques** couvrant tous les besoins de développement :

| Prompt | Description | Remplace |
|--------|-------------|----------|
| **[develop.md](./develop.md)** | Développement : ajouter fonctionnalité, modifier comportement, corriger bug | add-feature, modify-behavior, fix-bug |
| **[test.md](./test.md)** | Tests : écrire, exécuter, déboguer, analyser couverture | add-test, debug-test, run-tests |
| **[document.md](./document.md)** | Documentation : écrire/MAJ docs, expliquer code, générer exemples | update-docs, explain-code, generate-examples |
| **[review.md](./review.md)** | Revue et qualité : code review, refactoring, optimisation | code-review, refactor |
| **[analyze.md](./analyze.md)** | Analyse et diagnostic : erreurs, comportements, réseaux RETE, performance | analyze-error, investigate, validate-network |
| **[maintain.md](./maintain.md)** | Maintenance : migration, nettoyage, licence, stats, optimisation | migrate, deep-clean, verify-license-compliance, stats-code, optimize-performance |

---

## 📋 Document de Référence

| Document | Description | Quand l'utiliser |
|----------|-------------|------------------|
| **[common.md](./common.md)** | Standards communs du projet | **TOUJOURS** - Consulter avant tout développement |

Ce document contient :
- 🔒 Licence et Copyright (obligatoire)
- ⚠️ Règles strictes - Code Go
- 🧪 Standards de tests
- 📚 Documentation
- 🔧 Outils et commandes
- 🎨 Conventions de nommage
- 📋 Checklist avant commit
- 🚀 Workflow de développement

---

## 🚀 Utilisation

### Format Standard

```
Utilise le prompt "[nom]" pour [action]
```

### Exemples

```bash
# Développement
Utilise le prompt "develop" pour ajouter une fonctionnalité de validation
Utilise le prompt "develop" pour corriger le bug dans le module X

# Tests
Utilise le prompt "test" pour écrire des tests pour la fonction Y
Utilise le prompt "test" pour déboguer le test qui échoue

# Documentation
Utilise le prompt "document" pour expliquer le code du module Z
Utilise le prompt "document" pour générer des exemples .tsd

# Revue et Qualité
Utilise le prompt "review" pour faire une code review du fichier X
Utilise le prompt "review" pour refactoriser la fonction complexe Y

# Analyse
Utilise le prompt "analyze" pour analyser l'erreur "variable non liée"
Utilise le prompt "analyze" pour valider le réseau RETE

# Maintenance
Utilise le prompt "maintain" pour migrer vers Go 1.21
Utilise le prompt "maintain" pour générer les stats du projet
Utilise le prompt "maintain" pour vérifier les licences
```

---

## 🔍 Par Besoin

### Je veux développer
- **Ajouter une fonctionnalité** → [develop.md](./develop.md)
- **Modifier un comportement** → [develop.md](./develop.md)
- **Corriger un bug** → [develop.md](./develop.md)

### Je veux tester
- **Écrire des tests** → [test.md](./test.md)
- **Exécuter les tests** → [test.md](./test.md)
- **Déboguer un test** → [test.md](./test.md)
- **Analyser la couverture** → [test.md](./test.md)

### Je veux documenter
- **Mettre à jour la doc** → [document.md](./document.md)
- **Expliquer du code** → [document.md](./document.md)
- **Générer des exemples** → [document.md](./document.md)

### Je veux améliorer la qualité
- **Code review** → [review.md](./review.md)
- **Refactoriser** → [review.md](./review.md)

### Je veux analyser
- **Analyser une erreur** → [analyze.md](./analyze.md)
- **Investiguer un comportement** → [analyze.md](./analyze.md)
- **Valider un réseau RETE** → [analyze.md](./analyze.md)
- **Analyser la performance** → [analyze.md](./analyze.md)

### Je veux maintenir
- **Migrer version/dépendances** → [maintain.md](./maintain.md)
- **Nettoyer le projet** → [maintain.md](./maintain.md)
- **Vérifier les licences** → [maintain.md](./maintain.md)
- **Générer des stats** → [maintain.md](./maintain.md)
- **Optimiser la performance** → [maintain.md](./maintain.md)

---

## 📊 Statistiques

- **Prompts universels** : 6
- **Document de référence** : 1 (common.md)
- **Couverture** : 100% des besoins de développement

---

## 🎓 Parcours Recommandés

### 👶 Nouveau sur le Projet
1. Lire [common.md](./common.md) (standards du projet)
2. Explorer [develop.md](./develop.md) (développement)
3. Parcourir [test.md](./test.md) (tests)

### 👨‍💻 Développeur
1. [common.md](./common.md) - Toujours à portée de main
2. [develop.md](./develop.md) - Développement quotidien
3. [test.md](./test.md) - Tests systématiques
4. [review.md](./review.md) - Qualité du code

### 🐛 Debugger
1. [analyze.md](./analyze.md) - Diagnostiquer le problème
2. [develop.md](./develop.md) - Corriger le bug
3. [test.md](./test.md) - Valider la correction

### 📝 Documentation Writer
1. [document.md](./document.md) - Toute la documentation
2. [common.md](./common.md) - Standards à respecter

### 🔧 Mainteneur
1. [maintain.md](./maintain.md) - Toute la maintenance
2. [review.md](./review.md) - Qualité globale
3. [analyze.md](./analyze.md) - Diagnostic système

---

## 📚 Ressources

- [common.md](./common.md) - Standards du projet ⭐
- [README.md](./README.md) - Documentation du système de prompts
- [Makefile](../../Makefile) - Commandes du projet
- [Documentation](../../docs/) - Documentation technique

---

## ✨ Avantages des Prompts Universels

### Simplicité
- ✅ 6 prompts au lieu de 19
- ✅ Un prompt par catégorie d'action
- ✅ Pas d'hésitation sur lequel choisir

### Cohérence
- ✅ Tous référencent [common.md](./common.md)
- ✅ Standards unifiés
- ✅ Pas de redondance

### Maintenabilité
- ✅ Mise à jour centralisée
- ✅ Une seule source de vérité
- ✅ Facile à faire évoluer

### Exhaustivité
- ✅ Tous les besoins couverts
- ✅ Générique et adaptable
- ✅ Extensible facilement

---

## 🔗 Liens Utiles

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

**Version** : 2.0  
**Dernière mise à jour** : Décembre 2024  
**Prompts** : Universels et génériques