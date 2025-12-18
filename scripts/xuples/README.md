# Scripts du Plan d'Action Xuples

Ce répertoire contient le plan d'action complet pour l'implémentation du système xuples dans TSD.

## 📋 Vue d'Ensemble

Le système xuples sépare clairement le moteur de règles RETE et le système de xuple-spaces, avec :
- Exécution immédiate des actions (pas de stockage dans les terminal nodes)
- Actions par défaut (Print, Log, Update, Insert, Retract, Xuple)
- Module xuples découplé avec politiques configurables
- Commande `xuple-space` dans le langage TSD

## 🗂️ Structure

```
scripts/xuples/
├── README.md                           # Ce fichier
├── 00-INDEX.md                         # Vue d'ensemble du plan complet
├── 01-analyze-existing-actions.md      # Analyse de l'existant
├── 02-design-xuples-architecture.md    # Conception architecture
├── 03-extend-parser-xuplespace.md      # Parser xuple-space
├── 04-implement-default-actions.md     # Actions par défaut
├── 05-modify-rete-immediate-execution.md # RETE exécution immédiate
├── 06-implement-xuples-module.md       # Module xuples core
├── 07-integrate-xuple-action.md        # Intégration action Xuple
├── 08-test-complete-system.md          # Tests exhaustifs
├── 09-finalize-documentation.md        # Documentation finale
└── 10-final-validation.md              # Validation et intégration
```

## 🚀 Démarrage

### Prérequis

1. **Lire les standards du projet** :
   ```bash
   cat .github/prompts/common.md
   ```

2. **Installer les outils de développement** (si pas déjà fait) :
   ```bash
   go install honnef.co/go/tools/cmd/staticcheck@latest
   go install github.com/kisielk/errcheck@latest
   go install golang.org/x/tools/cmd/goimports@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   go install github.com/securego/gosec/v2/cmd/gosec@latest
   go install golang.org/x/vuln/cmd/govulncheck@latest
   ```

3. **Vérifier l'environnement** :
   ```bash
   make help
   go version
   git status
   ```

### Utilisation

1. **Commencer par lire l'INDEX** :
   ```bash
   cat scripts/xuples/00-INDEX.md
   ```

2. **Exécuter les prompts dans l'ordre** :
   - Lire complètement le prompt avant de commencer
   - Suivre toutes les tâches
   - Créer tous les livrables
   - Valider les critères de succès
   - Passer au prompt suivant

3. **Validation à chaque étape** :
   ```bash
   # Après chaque prompt majeur
   make test-unit
   make lint
   git status
   ```

## 📊 Progression

### Phase 1 : Analyse et Conception (01-02)
**Durée** : 5-7h | **Status** : ⬜ À faire

- [ ] Analyse de l'existant (01)
- [ ] Conception de l'architecture (02)

### Phase 2 : Extensions Langage (03-04)
**Durée** : 6-8h | **Status** : ⬜ À faire

- [ ] Parser xuple-space (03)
- [ ] Actions par défaut (04)

### Phase 3 : Modification RETE (05)
**Durée** : 3-4h | **Status** : ⬜ À faire

- [ ] Exécution immédiate (05)

### Phase 4 : Module Xuples (06-07)
**Durée** : 6-8h | **Status** : ⬜ À faire

- [ ] Module xuples (06)
- [ ] Intégration action Xuple (07)

### Phase 5 : Finalisation (08-10)
**Durée** : 7-10h | **Status** : ⬜ À faire

- [ ] Tests complets (08)
- [ ] Documentation (09)
- [ ] Validation finale (10)

## ✅ Checklist Globale

### Fonctionnalités
- [ ] Commande `xuple-space` parsable
- [ ] 6 actions par défaut implémentées
- [ ] Exécution immédiate des actions
- [ ] Module xuples avec politiques
- [ ] Intégration RETE ↔ xuples

### Qualité
- [ ] Tests > 80% couverture
- [ ] 0 race condition
- [ ] 0 hardcoding
- [ ] 0 erreur lint
- [ ] 0 vulnérabilité

### Documentation
- [ ] Guide utilisateur
- [ ] Documentation architecture
- [ ] Exemples fonctionnels
- [ ] GoDoc complet
- [ ] Guide migration

## 🎯 Objectifs par Prompt

| Prompt | Objectif Principal | Durée | Livrables |
|--------|-------------------|-------|-----------|
| 01 | Analyser l'existant | 2-3h | 6 docs analyse |
| 02 | Concevoir l'architecture | 3-4h | 8 docs conception |
| 03 | Parser xuple-space | 3-4h | Grammaire, AST, tests |
| 04 | Actions par défaut | 3-4h | defaults.tsd, loader |
| 05 | RETE immédiat | 3-4h | TerminalNode modifié |
| 06 | Module xuples | 4-5h | Package complet |
| 07 | Action Xuple | 2-3h | Intégration |
| 08 | Tests exhaustifs | 3-4h | E2E, perf, concurrence |
| 09 | Documentation | 2-3h | Guides, exemples |
| 10 | Validation finale | 2-3h | Rapport, release |

## 📚 Ressources

### Standards (OBLIGATOIRE)
- `.github/prompts/common.md` - Standards du projet

### Documentation
- `tsd/rete/docs/TUPLE_SPACE_IMPLEMENTATION.md` - Implémentation actuelle
- Thread précédent - Analyse tuple-space

### Références Externes
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Keep a Changelog](https://keepachangelog.com/)

## ⚠️ Important

### À FAIRE
✅ Lire `.github/prompts/common.md` AVANT de commencer  
✅ Suivre les prompts dans l'ordre (01 → 10)  
✅ Valider chaque prompt avant de passer au suivant  
✅ Commiter fréquemment avec messages clairs  
✅ Tester en continu  
✅ Documenter au fur et à mesure  

### À ÉVITER
❌ Sauter des prompts  
❌ Hardcoder des valeurs  
❌ Créer du couplage fort  
❌ Oublier les copyrights  
❌ Négliger les tests  
❌ Ignorer le linting  

## 🆘 Aide

En cas de problème :

1. **Relire le prompt** concerné attentivement
2. **Consulter l'INDEX** (`00-INDEX.md`)
3. **Vérifier les standards** (`.github/prompts/common.md`)
4. **Examiner la conception** (prompts 01-02)
5. **Demander de l'aide** avec contexte précis

## 📞 Contact

Pour toute question sur ce plan d'action, consulter :
- L'INDEX complet (`00-INDEX.md`)
- Les standards du projet (`.github/prompts/common.md`)
- La documentation TSD existante

## 🎉 Conclusion

Ce plan d'action guidé vous mènera de l'analyse à la production d'un système xuples complet, testé, documenté et prêt à l'emploi.

**Bonne implémentation ! 🚀**

---

**Version** : 1.0.0  
**Date** : 2025  
**Licence** : MIT