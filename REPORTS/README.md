# REPORTS - Rapports de l'Assistant IA

Ce répertoire contient **tous les rapports générés par l'assistant IA** suite à ses actions sur le projet.

## 📋 Règle Obligatoire

**TOUS les rapports produits par l'assistant IA doivent être stockés dans ce répertoire.**

### Types de Rapports

- 🧹 **Rapports de nettoyage** - Détails des opérations de nettoyage
- 📊 **Rapports d'analyse** - Analyses de code, métriques, statistiques
- 🔍 **Rapports d'audit** - Audits de sécurité, performance, qualité
- 🐛 **Rapports de débogage** - Investigations de bugs, résolutions
- ⚡ **Rapports d'optimisation** - Améliorations de performance
- 🔄 **Rapports de refactoring** - Restructurations de code
- 🧪 **Rapports de tests** - Résultats de tests, couverture
- 📝 **Autres rapports** - Tout rapport généré par l'IA

## ⚠️ Important

Ce répertoire est **exclu de Git** (via `.gitignore`). Les rapports sont générés localement et ne sont pas versionnés.

### Pourquoi ?

- Les rapports sont des artefacts temporaires de travail
- Ils ne constituent pas de la documentation utilisateur
- Ils peuvent être volumineux et changer fréquemment
- Ils sont spécifiques à une session de travail locale

## 📁 Structure

```
REPORTS/
├── README.md (ce fichier)
├── CLEANUP_SUMMARY.md
├── DOCUMENTATION_CLEANUP_REPORT.md
└── [autres rapports futurs...]
```

## 📝 Rapports Actuels

### CLEANUP_SUMMARY.md
Résumé concis du nettoyage de la documentation (Janvier 2025)

### DOCUMENTATION_CLEANUP_REPORT.md
Rapport détaillé du nettoyage complet de la documentation du projet TSD

## 🔍 Consultation

Pour consulter les rapports :
```bash
cd REPORTS
ls -la
cat CLEANUP_SUMMARY.md
```

---

**Note** : Si vous avez besoin de conserver un rapport de manière permanente, déplacez-le dans la documentation officielle (`docs/`) et retirez-le de ce répertoire.