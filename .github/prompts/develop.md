# 🔧 Développement - Prompt Universel

> **📋 Standards** : Ce prompt respecte les règles de [common.md](./common.md)

## 🎯 Objectif

Développer du code pour le projet TSD : ajouter une fonctionnalité, modifier un comportement, ou corriger un bug.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter [common.md](./common.md) :
- [🔒 Licence et Copyright](./common.md#licence-et-copyright) - En-tête obligatoire
- [⚠️ Interdictions](./common.md#interdictions-absolues) - Aucun hardcoding, code générique
- [🧪 Standards Tests](./common.md#standards-de-tests) - Couverture > 80%
- [📋 Checklist Commit](./common.md#checklist-avant-commit) - Validation finale

---

## 📋 Instructions

### 1. Définir le Besoin

**Précise clairement** :
- **Type** : [ ] Nouvelle fonctionnalité  [ ] Modification  [ ] Correction bug
- **Description** : Que faut-il faire ?
- **Motivation** : Pourquoi ?
- **Portée** : Modules/fichiers impactés
- **Contraintes** : Limites, dépendances, compatibilité

**Si bug** :
- Comportement observé vs attendu
- Étapes de reproduction
- Logs/erreurs si disponibles

### 2. Analyser l'Existant

1. **Examiner le code concerné**
   - Comprendre l'architecture actuelle
   - Identifier les patterns utilisés
   - Repérer les conventions du module

2. **Vérifier les dépendances**
   - Impacts sur autres modules
   - Interfaces à respecter
   - Tests existants à maintenir

3. **Valider l'approche**
   - Solution la plus simple
   - Éviter la sur-ingénierie
   - Réutiliser l'existant si possible

### 3. Concevoir

1. **Signature** : Fonctions, interfaces, types
2. **Visibilité** : Tout privé sauf exports nécessaires (voir [common.md](./common.md))
3. **Tests** : Scénarios de test (TDD encouragé)
4. **Documentation** : GoDoc, commentaires

### 4. Implémenter

**Ordre recommandé** :
1. **En-tête copyright** (obligatoire - voir [common.md](./common.md#en-tête-de-copyright-obligatoire))
2. **Tests d'abord** (TDD)
3. **Code minimal** fonctionnel
4. **Refactoring** si nécessaire
5. **Documentation** (GoDoc + exemples)

**Points d'attention** :
- ✅ Code générique avec paramètres (pas de hardcoding)
- ✅ Constantes nommées pour toutes valeurs
- ✅ Tout privé par défaut (minimiser exports)
- ✅ Gestion d'erreurs explicite
- ✅ Validation des entrées
- ✅ Messages d'erreur descriptifs

### 5. Valider

```bash
# Formattage
go fmt ./...
goimports -w .

# Validation
go vet ./...
staticcheck ./...
errcheck ./...

# Tests
go test ./...                    # Tous les tests
go test -cover ./...             # Avec couverture
go test -race ./...              # Race conditions

# Validation complète
make validate
```

### 6. Documenter

- [ ] GoDoc pour exports
- [ ] Commentaires inline si code complexe
- [ ] Exemples de tests `.tsd` si applicable
- [ ] README module si changement majeur
- [ ] CHANGELOG.md si pertinent

---

## 📝 Template de Code

### Nouveau Fichier

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package monpackage

// privateHelper fait quelque chose (privé par défaut)
func privateHelper(param string) string {
    // Implémentation
    return result
}

// PublicFunction est exportée car fait partie de l'API publique
// Description détaillée du comportement, paramètres, retours.
func PublicFunction(param string) (result string, err error) {
    // Validation entrée
    if param == "" {
        return "", errors.New("param ne peut pas être vide")
    }
    
    // Traitement
    result = privateHelper(param)
    return result, nil
}
```

### Tests (TDD)

```go
func TestPublicFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"cas nominal", "input", "expected", false},
        {"entrée vide", "", "", true},
        {"cas limite", "x", "y", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := PublicFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("❌ Erreur = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("❌ Attendu '%s', reçu '%s'", tt.want, got)
            }
        })
    }
}
```

---

## ✅ Checklist Finale

Avant de commit, vérifier [common.md#checklist-avant-commit](./common.md#checklist-avant-commit) :

- [ ] En-tête copyright présent
- [ ] Aucun hardcoding (valeurs, chemins, configs)
- [ ] Code générique avec paramètres/interfaces
- [ ] Constantes nommées pour toutes valeurs
- [ ] Variables/fonctions privées par défaut
- [ ] Exports publics minimaux et justifiés
- [ ] `go fmt` + `goimports` appliqués
- [ ] `go vet` + `staticcheck` + `errcheck` sans erreur
- [ ] Tests écrits et passent (couverture > 80%)
- [ ] GoDoc complet pour exports
- [ ] `make validate` passe
- [ ] Tous les tests passent

---

## 🎯 Principes

1. **Simplicité** : La solution la plus simple qui fonctionne
2. **Généricité** : Code réutilisable, pas de cas spécifiques hardcodés
3. **Encapsulation** : Privé par défaut, API publique minimale
4. **Testabilité** : Tests d'abord, couverture > 80%
5. **Lisibilité** : Code auto-documenté, noms explicites
6. **Robustesse** : Validation entrées, gestion erreurs

---

## 🚫 Anti-Patterns

- ❌ Hardcoding de valeurs, chemins, configurations
- ❌ Sur-ingénierie, complexité inutile
- ❌ Exports publics non nécessaires
- ❌ Tests absents ou insuffisants
- ❌ Code dupliqué (DRY)
- ❌ Magic numbers, magic strings
- ❌ Gestion d'erreurs négligée
- ❌ Documentation absente

---

## 📚 Ressources

- [common.md](./common.md) - Standards du projet
- [Makefile](../../Makefile) - Commandes disponibles
- [Documentation](../../docs/) - Architecture et guides
- [Effective Go](https://go.dev/doc/effective_go)

---

**Workflow** : Analyser → Concevoir → Tests → Implémenter → Valider → Documenter → Commit
