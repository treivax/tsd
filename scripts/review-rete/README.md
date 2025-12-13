# 🔍 Série de Prompts - Revue Complète Module RETE

**Objectif:** Revue systématique et complète du module `rete` selon `.github/prompts/review.md`

**Organisation:** 11 prompts numérotés (00-10) pour revue incrémentale compatible avec contexte Zed (128k tokens)

---

## 📋 Vue d'Ensemble

### Structure de la Série

| Prompt | Domaine | Fichiers | Lignes | Priorité | Durée Estimée |
|--------|---------|----------|--------|----------|---------------|
| **00** | Overview & Plan | - | - | Info | 15 min |
| **01** | Core RETE Nodes | 8 | ~2,000 | ⚠️ Critique | 2-3h |
| **02** | Bindings & Chains | 6 | ~1,500 | ⚠️ Critique | 2-3h |
| **03** | Alpha Network | 10 | ~2,500 | Haute | 2-3h |
| **04** | Beta Network | 8 | ~2,200 | ⚠️ Critique | 2-3h |
| **05** | Arithmétique | 8 | ~2,800 | Moyenne | 2-3h |
| **06** | Builders | 12 | ~3,000 | Haute | 3-4h |
| **07** | Actions | 8 | ~1,800 | Moyenne | 2h |
| **08** | Pipeline | 6 | ~2,000 | Haute | 2-3h |
| **09** | Métriques | 10 | ~2,500 | Basse | 1-2h |
| **10** | Utilitaires | ~10 | ~1,500 | Basse | 1-2h |

**Total estimé:** 20-28 heures de revue approfondie

---

## 🚀 Utilisation

### Pré-requis

1. **Outils installés:**
   ```bash
   go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
   go install honnef.co/go/tools/cmd/staticcheck@latest
   go install github.com/kisielk/errcheck@latest
   ```

2. **Tests passants:**
   ```bash
   cd /home/resinsec/dev/tsd
   go test ./rete/... -v
   # Doit afficher: ok  github.com/treivax/tsd/rete
   ```

3. **Baseline établie:**
   ```bash
   # Sauvegarder état actuel
   git checkout -b review-rete-baseline
   git add -A
   git commit -m "chore: baseline avant revue rete"
   
   # Créer branche de travail
   git checkout -b review-rete-work
   ```

### Workflow d'Exécution

#### Pour Chaque Prompt (01-10)

1. **Ouvrir dans Zed** (nouvelle session)
   - Charger le prompt `0X_domain.md`
   - Ajouter au contexte les fichiers mentionnés dans le périmètre

2. **Exécuter la revue** selon le prompt
   - Analyser selon checklist
   - Identifier problèmes
   - Proposer corrections
   - Implémenter changements
   - Valider tests

3. **Validation intermédiaire**
   ```bash
   # Tests du module
   go test ./rete/... -v
   
   # Complexité
   gocyclo -over 15 rete/
   
   # Vérifications
   go vet ./rete/...
   staticcheck ./rete/...
   
   # Formatage
   go fmt ./rete/...
   ```

4. **Commit atomique**
   ```bash
   git add rete/
   git commit -m "refactor(rete): [domaine] - [description courte]
   
   - Changement 1
   - Changement 2
   - Fixes complexité/duplication/etc.
   
   Prompt: 0X_domain.md"
   ```

5. **Générer rapport** dans `REPORTS/review-rete/`
   ```bash
   mkdir -p REPORTS/review-rete
   # Créer 0X_report.md avec résultats
   ```

6. **Passer au prompt suivant**

### Validation Globale (Après Prompt 10)

```bash
# Tests complets
go test ./... -v

# Couverture
go test -coverprofile=coverage.out ./rete/...
go tool cover -func=coverage.out

# Complexité finale
gocyclo -over 15 rete/ | wc -l
# Cible: 0 fonctions

# Métriques
gocyclo -top 10 rete/

# Vérifications complètes
go vet ./rete/...
staticcheck ./rete/...
errcheck ./rete/...

# Validation finale
make validate  # Si Makefile disponible
```

---

## 📊 Objectifs Globaux

### Métriques Cibles (Fin de Revue)

| Métrique | Actuel | Cible | Critique |
|----------|--------|-------|----------|
| **Complexité max** | 48 | <20 | ⚠️ Oui |
| **Fonctions >15** | ~50 | 0 | ⚠️ Oui |
| **Couverture** | 80.8% | >85% | Non |
| **Duplication** | ? | <5% | Oui |
| **GoDoc exports** | ~90% | 100% | Non |
| **Warnings** | 0 | 0 | ✅ OK |

### Qualité Globale

- ✅ Tous tests passent (100%)
- ✅ Aucune régression fonctionnelle
- ✅ Performance préservée ou améliorée
- ✅ Architecture SOLID respectée
- ✅ Code auto-documenté
- ✅ Encapsulation rigoureuse

---

## 🎯 Focus par Prompt

### Prompt 01 - Core RETE Nodes ⚠️
**Critique:** Fondations du système
- Décomposer `evaluateSimpleJoinConditions` (complexité 26)
- Valider thread-safety des nœuds
- Optimiser gestion mémoire

### Prompt 02 - Bindings ⚠️
**Critique:** Post-fix bug partage
- Valider immuabilité complète
- Vérifier correction bug JoinNode
- Optimiser allocations

### Prompt 03 - Alpha Network
**Important:** Partage et normalisation
- Valider mécanismes de partage
- Optimiser normalisation conditions
- Vérifier cache efficacité

### Prompt 04 - Beta Network ⚠️
**Critique:** Jointures et cascades
- Valider partage JoinNode (post-fix)
- Optimiser construction chaînes
- Vérifier cascadeLevel usage

### Prompt 05 - Arithmétique
**Moyen:** Décomposition expressions
- Optimiser cache résultats
- Valider décomposition correcte
- Améliorer performance évaluation

### Prompt 06 - Builders
**Important:** Construction réseau
- Valider séparation responsabilités
- Optimiser orchestration
- Réduire complexité builders

### Prompt 07 - Actions
**Moyen:** Exécution et handlers
- Valider gestion erreurs
- Vérifier thread-safety
- Optimiser contexte exécution

### Prompt 08 - Pipeline
**Important:** Validation et robustesse
- Améliorer gestion erreurs
- Valider tous cas edge
- Optimiser validation

### Prompt 09 - Métriques
**Bas:** Observabilité
- Compléter métriques manquantes
- Optimiser overhead collection
- Documenter exposition

### Prompt 10 - Utilitaires
**Bas:** Helpers et utils
- Valider généricité
- Éliminer duplication
- Améliorer réutilisabilité

---

## 📝 Livrables Finaux

### Code
- [ ] 100% fichiers rete/ revus
- [ ] Complexité <15 partout
- [ ] Duplication <5%
- [ ] Tests >85% couverture
- [ ] GoDoc 100% exports
- [ ] Aucun warning

### Documentation
- [ ] Rapport par prompt (11 rapports)
- [ ] Rapport synthèse final
- [ ] Métriques avant/après
- [ ] Guide patterns identifiés
- [ ] Recommandations futures

### Validation
- [ ] Tous tests passent
- [ ] Aucune régression
- [ ] Performance préservée
- [ ] Architecture validée
- [ ] Code review ready

---

## 🔄 Gestion Interruptions

### Sauvegarder Progression

```bash
# Sauvegarder après chaque prompt
git add -A
git commit -m "wip: revue rete - prompt 0X en cours"
git push origin review-rete-work
```

### Reprendre Plus Tard

```bash
# Vérifier où vous en êtes
git log --oneline | head -10

# Consulter dernier rapport
ls -lt REPORTS/review-rete/

# Reprendre au prompt suivant
# Ouvrir 0X_domain.md et continuer
```

---

## 📚 Références

- **Standards:** `.github/prompts/review.md`
- **Conventions:** `.github/prompts/common.md`
- **État actuel:** `REPORTS/MAINTENANCE_REPORT.md`
- **Architecture:** `docs/architecture/`
- **Baseline:** Tag `review-rete-baseline`

---

## ⚠️ Avertissements

### À Faire
- ✅ Commits atomiques fréquents
- ✅ Tests après chaque changement
- ✅ Valider métriques régulièrement
- ✅ Documenter décisions
- ✅ Sauvegarder progression

### À Ne PAS Faire
- ❌ Changer comportement fonctionnel
- ❌ Optimiser sans mesurer
- ❌ Refactorer sans tests
- ❌ Tout changer d'un coup
- ❌ Ignorer warnings/erreurs

---

## 🎓 Principes

### Refactoring
1. **Incrémental** - Petites étapes
2. **Testé** - Validation continue
3. **Documenté** - Explications claires
4. **Réversible** - Commits atomiques

### Qualité
1. **Simplicité** - Solution la plus simple
2. **Lisibilité** - Code auto-documenté
3. **Maintenabilité** - Facile à modifier
4. **Performance** - Pas de dégradation

---

## 📞 Support

**Questions ou problèmes:**
1. Consulter `.github/prompts/review.md`
2. Vérifier `REPORTS/MAINTENANCE_REPORT.md`
3. Relire prompt 00 (overview)
4. Documenter dans rapport si décision difficile

---

**Bon courage pour la revue ! 🚀**

**Durée totale estimée:** 20-28 heures  
**Résultat attendu:** Module rete de qualité production++  
**Date création:** 2024-12-15