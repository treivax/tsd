# Beta Node Sharing & Optimization - Prompts Suite

Ce dossier contient l'ensemble des prompts pour implémenter le partage et l'optimisation des nœuds Beta (JoinNodes) dans le moteur RETE TSD.

## 📋 Vue d'Ensemble

Le système de Beta Sharing est l'équivalent pour les JoinNodes de ce qui a été fait pour les AlphaNodes. Il permet de :
- Partager les JoinNodes identiques entre règles
- Construire des chaînes optimisées (BetaChains)
- Utiliser un cache LRU pour les opérations de jointure
- Collecter des métriques de performance
- Exporter vers Prometheus

## 🎯 Objectifs

- **Réduction mémoire:** 40-70% pour règles avec patterns similaires
- **Performance:** Cache avec hit rate > 70%
- **Backward compatible:** 100% compatible avec le code existant
- **Production ready:** Tests complets, documentation exhaustive

## 📦 Structure des Prompts

### Phase 1 : Foundation (Prompts 1-3)
- **01-analyze-existing.md** - Analyse de l'implémentation actuelle
- **02-design-sharing.md** - Conception du BetaSharingRegistry
- **03-design-chains.md** - Conception des BetaChains

### Phase 2 : Core Implementation (Prompts 4-7)
- **04-implement-registry.md** - Implémentation du registre de partage
- **05-implement-builder.md** - Implémentation du BetaChainBuilder
- **06-implement-cache.md** - Intégration du cache LRU
- **07-integrate-network.md** - Intégration dans ReteNetwork

### Phase 3 : Metrics & Performance (Prompts 8-9)
- **08-implement-metrics.md** - Système de métriques
- **09-benchmark-optimize.md** - Benchmarks et optimisations

### Phase 4 : Documentation (Prompts 10-11)
- **10-write-technical-docs.md** - Documentation technique complète
- **11-write-examples-migration.md** - Exemples et guide de migration

### Phase 5 : Testing & Validation (Prompts 12-13)
- **12-write-integration-tests.md** - Tests d'intégration
- **13-validate-compatibility.md** - Validation backward compatibility

### Phase 6 : Finalization (Prompt 14)
- **14-finalize-cleanup.md** - Nettoyage final et synthèse

## 🚀 Ordre d'Exécution

### Séquence Obligatoire

```
1. Prompts 1-3 (Foundation)
   ↓
2. Prompts 4-5 (Registry + Builder)
   ↓
3. Prompt 7 (Intégration réseau)
   ↓ (en parallèle avec Prompt 6)
4. Prompt 6 (Cache LRU)
   ↓
5. Prompts 8-9 (Métriques)
   ↓
6. Prompts 10-11 (Documentation)
   ↓
7. Prompts 12-13 (Tests)
   ↓
8. Prompt 14 (Finalisation)
```

### Dépendances

- **Prompts 4-7** nécessitent les conceptions de 1-3
- **Prompt 7** nécessite 4 et 5
- **Prompts 8-9** nécessitent 7 (implémentation complète)
- **Prompts 10-11** peuvent être faits après une implémentation stable
- **Prompts 12-13** nécessitent l'implémentation complète
- **Prompt 14** doit être fait en dernier

## 📊 Durées Estimées

| Phase | Durée Totale | Détails |
|-------|--------------|---------|
| Foundation | 5-8h | 3 prompts d'analyse et conception |
| Core Implementation | 14-20h | 4 prompts d'implémentation |
| Metrics & Performance | 5-7h | 2 prompts de métriques et benchmarks |
| Documentation | 7-9h | 2 prompts de documentation |
| Testing & Validation | 5-7h | 2 prompts de tests |
| Finalization | 2-3h | 1 prompt de finalisation |
| **TOTAL** | **38-54h** | 14 prompts |

## 🔧 Utilisation

### Pour Exécuter un Prompt

1. Lire le prompt complet dans son fichier
2. Copier le contenu du prompt
3. Le donner à l'assistant IA avec le contexte du projet
4. Valider les livrables attendus
5. Passer au prompt suivant

### Exemple

```bash
# Lire le premier prompt
cat .github/prompts/beta/01-analyze-existing.md

# L'exécuter avec l'assistant
# Vérifier les livrables :
# - rete/docs/BETA_NODES_ANALYSIS.md créé
# - Diagrammes présents
# - Recommandations claires

# Passer au suivant
cat .github/prompts/beta/02-design-sharing.md
```

## 📚 Livrables Attendus (Vue d'Ensemble)

### Code (14 nouveaux fichiers Go)
- `rete/beta_sharing.go` + tests
- `rete/beta_chain_builder.go` + tests
- `rete/beta_chain_metrics.go` + tests
- `rete/beta_join_cache.go` + tests (optionnel)
- `rete/beta_chain_performance_test.go`
- `rete/beta_chain_integration_test.go`
- `rete/beta_backward_compatibility_test.go`
- Extensions de `rete/chain_config.go` et `rete/network.go`

### Documentation (12+ fichiers MD)
- `rete/BETA_CHAINS_TECHNICAL_GUIDE.md`
- `rete/BETA_CHAINS_USER_GUIDE.md`
- `rete/BETA_CHAINS_EXAMPLES.md`
- `rete/BETA_CHAINS_MIGRATION.md`
- `rete/BETA_CHAINS_INDEX.md`
- `rete/BETA_NODE_SHARING.md`
- `rete/docs/BETA_NODES_ANALYSIS.md`
- `rete/docs/BETA_SHARING_DESIGN.md`
- `rete/docs/BETA_CHAINS_DESIGN.md`
- `rete/docs/BETA_PERFORMANCE_REPORT.md`
- `rete/BETA_COMPATIBILITY_VALIDATION_REPORT.md`
- `rete/BETA_IMPLEMENTATION_SUMMARY.md`

### Exemples
- `examples/beta_chains/` (dossier complet avec exemple exécutable)

### Mise à jour
- `CHANGELOG.md`
- `README.md`
- `rete/README.md`

## ✅ Critères de Succès Globaux

### Code
- [ ] Tous les tests passent (100%)
- [ ] go vet sans erreur
- [ ] Couverture > 70% sur nouveaux fichiers
- [ ] Backward compatible (tests de régression passent)
- [ ] Thread-safe (tests -race passent)

### Fonctionnalités
- [ ] BetaSharingRegistry fonctionne
- [ ] BetaChainBuilder construit des chaînes correctes
- [ ] Partage des JoinNodes identiques vérifié
- [ ] Cache LRU efficace (hit rate > 70%)
- [ ] Métriques exportées correctement
- [ ] Intégration Prometheus fonctionnelle

### Performance
- [ ] Réduction mémoire mesurée (> 30%)
- [ ] Partage vérifié sur exemples réels
- [ ] Cache hit rate > 70%
- [ ] Pas de régression de temps d'exécution
- [ ] Benchmarks documentés

### Documentation
- [ ] Guide technique complet (30-40 pages)
- [ ] Guide utilisateur complet (20-30 pages)
- [ ] 15+ exemples concrets
- [ ] Guide de migration détaillé
- [ ] Index centralisé
- [ ] README mis à jour
- [ ] CHANGELOG à jour

## 🔗 Références

### Documentation du Plan Complet
- `rete/docs/BETA_NODES_OPTIMIZATION_PLAN.md` - Plan détaillé complet

### Implémentation Alpha (à réutiliser)
- `rete/alpha_chain_builder.go`
- `rete/alpha_sharing.go`
- `rete/lru_cache.go`
- `rete/chain_config.go`
- `rete/chain_metrics.go`

### Documentation Alpha (structure à suivre)
- `rete/ALPHA_CHAINS_TECHNICAL_GUIDE.md`
- `rete/ALPHA_CHAINS_USER_GUIDE.md`
- `rete/ALPHA_CHAINS_EXAMPLES.md`
- `rete/ALPHA_CHAINS_MIGRATION.md`

### Tests Alpha (patterns à réutiliser)
- `rete/alpha_chain_integration_test.go`
- `rete/backward_compatibility_test.go`
- `rete/chain_performance_test.go`

## 📈 Suivi de Progression

### Template de Checklist

Copier cette checklist dans un fichier `BETA_PROGRESS.md` à la racine pour suivre l'avancement :

```markdown
# Progression Beta Node Sharing

## Phase 1 : Foundation
- [ ] Prompt 1 : Analyse existant (BETA_NODES_ANALYSIS.md créé)
- [ ] Prompt 2 : Design sharing (BETA_SHARING_DESIGN.md créé)
- [ ] Prompt 3 : Design chains (BETA_CHAINS_DESIGN.md créé)

## Phase 2 : Core Implementation
- [ ] Prompt 4 : BetaSharingRegistry (beta_sharing.go + tests)
- [ ] Prompt 5 : BetaChainBuilder (beta_chain_builder.go + tests)
- [ ] Prompt 6 : Cache LRU (intégration cache)
- [ ] Prompt 7 : Intégration réseau (network.go modifié)

## Phase 3 : Metrics & Performance
- [ ] Prompt 8 : BetaChainMetrics (beta_chain_metrics.go)
- [ ] Prompt 9 : Benchmarks (beta_chain_performance_test.go)

## Phase 4 : Documentation
- [ ] Prompt 10 : Docs techniques (3 guides créés)
- [ ] Prompt 11 : Exemples & migration (exemples + migration guide)

## Phase 5 : Testing & Validation
- [ ] Prompt 12 : Tests intégration (beta_chain_integration_test.go)
- [ ] Prompt 13 : Validation backward (beta_backward_compatibility_test.go)

## Phase 6 : Finalization
- [ ] Prompt 14 : Cleanup & synthèse (docs finales + CHANGELOG)

## Validation Finale
- [ ] Tous les tests passent
- [ ] go vet sans erreur
- [ ] Couverture > 70%
- [ ] Documentation complète
- [ ] Benchmarks validés
- [ ] Prêt pour commit
```

## 🎯 Prochaines Étapes

1. **Commencer par la Phase 1** - Exécuter les prompts 1-3 pour comprendre et concevoir
2. **Implémenter le Core** - Prompts 4-7 pour l'implémentation de base
3. **Ajouter les Métriques** - Prompts 8-9 pour mesurer les performances
4. **Documenter** - Prompts 10-11 pour la documentation complète
5. **Tester** - Prompts 12-13 pour valider
6. **Finaliser** - Prompt 14 pour nettoyer et synthétiser

## 💡 Conseils

- **Ne pas sauter de prompts** - Chaque prompt construit sur le précédent
- **Valider chaque étape** - Vérifier les livrables avant de continuer
- **Tester fréquemment** - Lancer `go test` après chaque implémentation
- **Documenter au fur et à mesure** - Ne pas attendre la fin
- **S'inspirer des Alpha** - Réutiliser les patterns qui ont fonctionné
- **Mesurer les gains** - Benchmarker avant/après pour valider les optimisations

## 📧 Support

Pour toute question :
1. Consulter `BETA_NODES_OPTIMIZATION_PLAN.md` pour le plan détaillé
2. Regarder l'implémentation Alpha comme référence
3. Lire la littérature RETE citée dans le plan
4. Ouvrir une issue GitHub si nécessaire

---

**Version:** 1.0  
**Dernière mise à jour:** 2025-11-27  
**Licence:** MIT  
**Durée totale estimée:** 38-54 heures