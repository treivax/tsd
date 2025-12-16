## 🎯 Description

<!-- Décrire clairement les changements apportés -->

**Type de changement** :
- [ ] 🐛 Bug fix (correction non-breaking)
- [ ] ✨ New feature (ajout de fonctionnalité non-breaking)
- [ ] 💥 Breaking change (fix ou feature qui casse la compatibilité)
- [ ] 📝 Documentation
- [ ] 🧪 Tests
- [ ] 🔧 Refactoring (pas de changement fonctionnel)
- [ ] ⚡ Performance

## 🔗 Issue Liée

Closes # <!-- Numéro de l'issue -->

## 📝 Changements

<!-- Lister les changements principaux -->

- 
- 
- 

## 🧪 Tests

**Tests ajoutés** :
- [ ] Tests unitaires
- [ ] Tests d'intégration
- [ ] Tests E2E
- [ ] N/A (changement sans code)

**Comment tester** :
```bash
# Commandes pour tester les changements
make test
```

**Couverture** :
- Couverture avant : XX%
- Couverture après : XX%

## 📸 Captures d'Écran (si applicable)

<!-- Ajouter captures d'écran pour changements visuels -->

## ✅ Checklist Avant Soumission

### Code

- [ ] Mon code suit les standards du projet ([common.md](../.github/prompts/common.md))
- [ ] J'ai ajouté les en-têtes copyright sur les nouveaux fichiers
- [ ] Aucun hardcoding (valeurs en constantes nommées)
- [ ] Code générique et réutilisable
- [ ] Variables/fonctions privées par défaut
- [ ] `go fmt` et `goimports` appliqués
- [ ] `go vet`, `staticcheck`, `errcheck` passent sans erreur

### Tests

- [ ] J'ai écrit des tests pour mes changements
- [ ] Tous les tests passent (`make test-complete`)
- [ ] Couverture > 80% maintenue
- [ ] Tests déterministes (pas de flaky tests)
- [ ] Pas de dépendances entre tests

### Documentation

- [ ] GoDoc ajouté pour les exports
- [ ] README mis à jour (si nécessaire)
- [ ] CHANGELOG.md mis à jour (section Unreleased)
- [ ] Documentation technique mise à jour (si nécessaire)

### Validation

- [ ] `make validate` passe sans erreur
- [ ] Branch à jour avec `main`
- [ ] Commits suivent convention ([type]: description)
- [ ] Pas de conflits de merge

## 📋 Contexte Additionnel

<!-- Toute information supplémentaire pour les reviewers -->

## 🔍 Points d'Attention pour Review

<!-- Signaler les parties nécessitant une attention particulière -->

---

**Pour les reviewers** : Vérifier que la checklist est complète avant approval.
