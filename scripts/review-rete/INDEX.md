# 📑 Index - Série de Prompts Revue RETE

**Organisation complète de la revue systématique du module rete**

---

## 📁 Fichiers Disponibles

### 1. Documentation Principale

| Fichier | Description | Priorité Lecture |
|---------|-------------|------------------|
| **INDEX.md** | Ce fichier - Point d'entrée | 1️⃣ Lire en premier |
| **README.md** | Vue d'ensemble complète | 2️⃣ Après index |
| **EXECUTION_GUIDE.md** | Guide pratique d'exécution | 3️⃣ Avant démarrage |

### 2. Prompts de Revue (Ordre d'Exécution)

| Prompt | Fichier | Domaine | Priorité | Durée |
|--------|---------|---------|----------|-------|
| **00** | `00_overview_and_plan.md` | Overview & Planification | Info | 15min |
| **01** | `01_core_rete_nodes.md` | Core RETE (Nœuds) | ⚠️ Critique | 2-3h |
| **02** | `02_bindings_chains.md` | Bindings Immuables | ⚠️ Critique | 2-3h |
| **03-10** | `03-10_prompts_condensed.md` | 8 prompts condensés | Variable | 12-20h |

### 3. Prompts Condensés (Prompt 03-10)

Dans `03-10_prompts_condensed.md`:
- **03** Alpha Network (Haute, 2-3h)
- **04** Beta Network (⚠️ Critique, 2-3h)
- **05** Arithmétique (Moyenne, 2-3h)
- **06** Builders (Haute, 3-4h)
- **07** Actions (Moyenne, 2h)
- **08** Pipeline (Haute, 2-3h)
- **09** Métriques (Basse, 1-2h)
- **10** Utilitaires (Basse, 1-2h)

---

## 🚀 Démarrage Rapide

### Pour Commencer MAINTENANT

```bash
# 1. Aller au répertoire
cd /home/resinsec/dev/tsd

# 2. Lire ce fichier (INDEX.md) ✅ Fait

# 3. Lire README.md
cat scripts/review-rete/README.md

# 4. Lire EXECUTION_GUIDE.md
cat scripts/review-rete/EXECUTION_GUIDE.md

# 5. Préparer environnement
git checkout -b review-rete-baseline
git commit -am "chore: baseline avant revue rete"
git checkout -b review-rete-work

# 6. Lancer Prompt 00 (overview)
cat scripts/review-rete/00_overview_and_plan.md

# 7. Puis Prompt 01 (core nodes)
# Ouvrir dans Zed avec fichiers contexte:
zed scripts/review-rete/01_core_rete_nodes.md
# + Ajouter: rete/network.go, rete/node*.go, etc.
```

---

## 📊 Vue d'Ensemble de la Revue

### Statistiques

- **Prompts total:** 11 (00-10)
- **Fichiers à revue:** ~152 fichiers .go (hors tests)
- **Lignes de code:** ~52,000 lignes
- **Durée estimée:** 20-28 heures
- **Complexité actuelle:** 98 fonctions >15

### Objectifs Globaux

| Métrique | Actuel | Cible | Critique |
|----------|--------|-------|----------|
| Complexité max | 48 | <20 | ⚠️ |
| Fonctions >15 | ~50 | 0 | ⚠️ |
| Couverture | 80.8% | >85% | Non |
| Duplication | ? | <5% | Oui |

---

## 📋 Workflow Recommandé

### Phase 1: Préparation (30 min)
1. Lire INDEX.md (ce fichier)
2. Lire README.md
3. Lire EXECUTION_GUIDE.md
4. Préparer environnement (git, outils)

### Phase 2: Prompts Critiques (6-9h)
1. **Prompt 00** - Overview (15min)
2. **Prompt 01** - Core Nodes ⚠️ (2-3h)
3. **Prompt 02** - Bindings ⚠️ (2-3h)
4. **Prompt 04** - Beta Network ⚠️ (2-3h)

### Phase 3: Prompts Importants (8-11h)
5. **Prompt 03** - Alpha Network (2-3h)
6. **Prompt 06** - Builders (3-4h)
7. **Prompt 08** - Pipeline (2-3h)

### Phase 4: Prompts Complémentaires (5-8h)
8. **Prompt 05** - Arithmétique (2-3h)
9. **Prompt 07** - Actions (2h)
10. **Prompt 09** - Métriques (1-2h)
11. **Prompt 10** - Utilitaires (1-2h)

### Phase 5: Finalisation (1h)
12. Validation finale
13. Rapport synthèse
14. Merge si satisfait

---

## 🎯 Points Clés par Prompt

### Critiques ⚠️

**Prompt 01 - Core Nodes**
- Décomposer `evaluateSimpleJoinConditions` (26 → <15)
- Valider thread-safety

**Prompt 02 - Bindings**
- Valider immuabilité
- Vérifier correction bug partage JoinNode

**Prompt 04 - Beta Network**
- Décomposer orchestration (48 → <20)
- Valider cascades 3+ variables

### Importants

**Prompt 03 - Alpha Network**
- Optimiser partage nœuds
- Valider normalisation

**Prompt 06 - Builders**
- Séparer responsabilités
- Réduire complexité

**Prompt 08 - Pipeline**
- Décomposer IngestFile (48 → <20)
- Améliorer gestion erreurs

---

## 📚 Références Externes

### Standards Projet
- `.github/prompts/review.md` - Standards de revue
- `.github/prompts/common.md` - Conventions projet

### État Actuel
- `REPORTS/MAINTENANCE_REPORT.md` - Rapport maintenance
- `SYNTHESE_VALIDATION_FINALE.md` - Validation post-fix

### Architecture
- `docs/architecture/` - Documentation architecture
- `rete/README.md` - README module rete

---

## ✅ Checklist Avant Démarrage

- [ ] INDEX.md lu ✅
- [ ] README.md lu
- [ ] EXECUTION_GUIDE.md lu
- [ ] Git baseline créée
- [ ] Outils installés (gocyclo, staticcheck)
- [ ] Tests passent (`go test ./rete/...`)
- [ ] Contexte Zed prêt (128k)

---

## 📞 Support

**En cas de question:**
1. Relire EXECUTION_GUIDE.md
2. Consulter prompt 00 (overview)
3. Vérifier standards (.github/prompts/)
4. Documenter décisions dans rapport

---

**Prêt à commencer ?** 🚀

1. ✅ INDEX.md lu
2. ⏭️  Lire README.md
3. ⏭️  Lire EXECUTION_GUIDE.md
4. ⏭️  Exécuter Prompt 00

**Bon courage pour la revue !**

---

**Date création:** 2024-12-15
**Version:** 1.0
**Auteur:** Équipe TSD
