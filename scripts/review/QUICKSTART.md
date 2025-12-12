# ⚡ Quick Start - Review Automatisée

Guide rapide pour lancer une review automatisée du module `constraint/`.

---

## 🚀 Setup en 3 Minutes

### 1. Installer Dépendances

```bash
# Ubuntu/Debian
sudo apt-get update && sudo apt-get install -y jq curl

# macOS
brew install jq curl
```

### 2. Configurer API Claude

```bash
# Obtenir clé API : https://console.anthropic.com/settings/keys

# Exporter la clé (remplacer YOUR_KEY)
export ANTHROPIC_API_KEY='sk-ant-api03-YOUR_KEY_HERE'

# Vérifier
echo $ANTHROPIC_API_KEY
```

### 3. Lancer Review

```bash
cd /path/to/tsd

# Exécution automatique sans confirmation
AUTO_CONFIRM=1 ./scripts/review/run_review.sh
```

**C'est tout !** Le script va :
- ✅ Exécuter 6 sessions de review
- ✅ Générer rapports dans `REPORTS/`
- ✅ Créer synthèse globale
- ⏱️ Durée : ~15-30 minutes

---

## 📊 Résultats

### Rapports Générés

```
REPORTS/
├── REVIEW_CONSTRAINT_SESSION_1_STATE_API.md
├── REVIEW_CONSTRAINT_SESSION_2_VALIDATION.md
├── REVIEW_CONSTRAINT_SESSION_3_PKG_VALIDATOR.md
├── REVIEW_CONSTRAINT_SESSION_4_TYPES_DOMAIN.md
├── REVIEW_CONSTRAINT_SESSION_5_FACTS_ACTIONS.md
├── REVIEW_CONSTRAINT_SESSION_6_CONFIG_CLI.md
└── REVIEW_CONSTRAINT_SUMMARY.md  ← Commencer ici
```

### Lecture Recommandée

```bash
# 1. Lire synthèse
cat REPORTS/REVIEW_CONSTRAINT_SUMMARY.md

# 2. Identifier problèmes critiques
grep -n "🔴" REPORTS/REVIEW_CONSTRAINT_SESSION_*.md

# 3. Lire rapports détaillés par priorité
cat REPORTS/REVIEW_CONSTRAINT_SESSION_1_STATE_API.md
cat REPORTS/REVIEW_CONSTRAINT_SESSION_2_VALIDATION.md
```

---

## 🔧 Options Utiles

### Mode Interactif (avec confirmation)

```bash
./scripts/review/run_review.sh
```

### Mode Debug (sauvegarder prompts)

```bash
DEBUG=1 AUTO_CONFIRM=1 ./scripts/review/run_review.sh
# Prompts dans : /tmp/review_prompt_*.txt
```

### En Arrière-Plan (avec log)

```bash
nohup ./scripts/review/run_review.sh > review.log 2>&1 &

# Suivre progression
tail -f review.log
```

---

## 💰 Coût Estimé

**API Claude Sonnet-4** :
- ~300K tokens input + ~100K tokens output
- **Coût** : ~$2.50 pour review complète
- 💡 **Tip** : $5 gratuits pour nouveaux comptes

---

## ❌ Problèmes Courants

### "ANTHROPIC_API_KEY non définie"

```bash
# Solution
export ANTHROPIC_API_KEY='votre-clé-ici'
```

### "jq: command not found"

```bash
# Ubuntu
sudo apt-get install jq

# macOS
brew install jq
```

### "Rate limit exceeded"

- Attendre 1-2 minutes
- Vérifier quota : https://console.anthropic.com/settings/usage
- Relancer : le script reprend où il s'est arrêté

---

## 📚 Documentation Complète

Voir `README.md` pour :
- Configuration avancée
- Troubleshooting détaillé
- Personnalisation
- Coûts détaillés
- Workflow recommandé

---

## ✅ Checklist Avant Lancement

- [ ] `jq` et `curl` installés
- [ ] API key Claude configurée (`echo $ANTHROPIC_API_KEY`)
- [ ] Dans répertoire projet (`cd /path/to/tsd`)
- [ ] Script exécutable (`chmod +x scripts/review/run_review.sh`)
- [ ] ~$3 de crédit API disponible (ou free tier actif)

---

## 🎯 Après Review

### Prioriser Actions

```bash
# Problèmes critiques (à traiter en premier)
grep "🔴 Critiques" REPORTS/REVIEW_CONSTRAINT_SESSION_*.md

# Problèmes majeurs (important)
grep "🟡 Majeurs" REPORTS/REVIEW_CONSTRAINT_SESSION_*.md

# Problèmes mineurs (amélioration)
grep "🟢 Mineurs" REPORTS/REVIEW_CONSTRAINT_SESSION_*.md
```

### Créer Plan Action

1. Lire synthèse complète
2. Grouper problèmes similaires
3. Estimer effort corrections
4. Créer issues/tickets
5. Planifier refactoring

---

**Besoin d'aide ?** Voir `README.md` ou documentation API Claude.