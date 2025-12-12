# 🔍 Review Session 6 - Config & CLI

**Module** : `constraint/`  
**Priorité** : 🟢 BASSE  
**Fichiers** : Infrastructure et commande  
**Lignes** : ~373 lignes

---

## 📋 Contexte

Cette session audite l'**infrastructure** : configuration, CLI, injection dépendances.

---

## 🎯 Objectifs

- Auditer configuration
- Valider CLI si applicable
- Vérifier 12-factor app
- Analyser injection dépendances
- Évaluer externalisation config

---

## 📂 Fichiers à Reviewer

```
constraint/internal/config/config.go  (223 lignes)
constraint/cmd/main.go                (~150 lignes estimé)
```

**Total** : ~373 lignes

---

## ❓ Questions Clés

- [ ] Configuration externalisée ?
- [ ] Pas de hardcoding ?
- [ ] CLI robuste ?
- [ ] Injection dépendances propre ?
- [ ] 12-factor respecté ?
- [ ] Variables d'environnement gérées ?

---

## ✅ Checklist Review

### Configuration
- [ ] Configuration externalisée (fichier/env)
- [ ] Pas de hardcoding valeurs
- [ ] Valeurs par défaut sensibles
- [ ] Validation configuration au démarrage
- [ ] Documentation config complète

### 12-Factor App
- [ ] Config dans environnement
- [ ] Séparation build/run/config
- [ ] Logs vers stdout/stderr
- [ ] Processes stateless
- [ ] Port binding configuré

### CLI
- [ ] Arguments bien définis
- [ ] Help/usage clair
- [ ] Gestion erreurs robuste
- [ ] Exit codes appropriés
- [ ] Flags cohérents

### Injection Dépendances
- [ ] Pas de globals mutables
- [ ] Dépendances injectées
- [ ] Testabilité facilitée
- [ ] Wire-up propre (main)
- [ ] Lifecycle clair

### Infrastructure
- [ ] Logging configuré
- [ ] Metrics si applicable
- [ ] Health checks si pertinent
- [ ] Graceful shutdown
- [ ] Resource cleanup

---

## 📊 Métriques

- Variables config
- Hardcoding détecté
- Flags CLI
- Globals mutables
- Coverage tests infra

---

## 📝 Format Rapport

Créer : `REPORTS/REVIEW_CONSTRAINT_SESSION_6_CONFIG_CLI.md`

```markdown
# 🔍 Review Constraint - Session 6 : Config & CLI

**Date** : YYYY-MM-DD
**Fichiers** : internal/config/config.go, cmd/main.go
**Lignes** : ~373

## 📊 Vue d'Ensemble
- Configuration : Bonne/Moyenne/Faible
- 12-factor : Respecté/Partiel/Non
- CLI : Robuste/Acceptable/Faible

## ✅ Points Forts
- ...

## ⚠️ Points d'Attention
- ...

## ❌ Problèmes
### 🔴 Critiques
- ...
### 🟡 Majeurs
- ...
### 🟢 Mineurs
- ...

## 💡 Recommandations
1. ...

## 📈 Métriques
- Variables config : X
- Hardcoding : X occurrences
- Flags CLI : X
- 12-factor score : X/12

## 🏁 Verdict
✅ / ⚠️ / ❌

## 🔜 Actions
- [ ] ...
```

---

## 🚨 Points d'Attention Spécifiques

### Configuration
- Toutes configs externalisables
- Pas de secrets en dur
- Variables environnement prioritaires
- Validation au démarrage

### 12-Factor
- Config séparée du code
- Backing services via URLs
- Logs comme flux d'événements
- Admin processes séparés

### CLI
- Help complet et clair
- Erreurs informatives
- Exit codes standard (0=OK, 1=erreur)
- Pas de panic en production

### Sécurité
- Pas de secrets loggués
- Pas de credentials hardcodés
- Variables sensibles protégées
- Permissions fichiers config

---

## 📚 Références

- `.github/prompts/review.md`
- `.github/prompts/common.md`
- [12-Factor App](https://12factor.net/)
- [CLI Best Practices](https://clig.dev/)