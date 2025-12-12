# 🔍 Automatisation Review Code - Module Constraint

Ce répertoire contient les scripts et fichiers pour automatiser la review complète du module `constraint/`.

---

## 📋 Contenu

### Fichiers Sessions
- `SESSION_1_STATE_API.md` - State Management & API (CRITIQUE)
- `SESSION_2_VALIDATION.md` - Validation Layer (HAUTE)
- `SESSION_3_PKG_VALIDATOR.md` - Package Validator (HAUTE)
- `SESSION_4_TYPES_DOMAIN.md` - Types & Domain (MOYENNE)
- `SESSION_5_FACTS_ACTIONS.md` - Facts, Actions & Logic (BASSE)
- `SESSION_6_CONFIG_CLI.md` - Config & CLI (BASSE)

### Script Principal
- `run_review.sh` - Script d'automatisation utilisant l'API Claude

---

## 🚀 Installation & Configuration

### 1. Prérequis

**Système** :
```bash
# Ubuntu/Debian
sudo apt-get install jq curl

# macOS
brew install jq curl
```

**API Key Claude (Anthropic)** :
1. Créer compte sur https://console.anthropic.com/
2. Obtenir API key dans Settings > API Keys
3. Configurer dans environnement

### 2. Configuration API

#### Option A - Variable d'Environnement (Temporaire)

```bash
export ANTHROPIC_API_KEY='sk-ant-api03-...'
```

#### Option B - Fichier ~/.bashrc (Permanent)

```bash
echo 'export ANTHROPIC_API_KEY="sk-ant-api03-..."' >> ~/.bashrc
source ~/.bashrc
```

#### Option C - Fichier .env Local

```bash
# Créer .env dans le répertoire projet
echo 'ANTHROPIC_API_KEY=sk-ant-api03-...' > .env

# Charger avant exécution
source .env
./scripts/review/run_review.sh
```

⚠️ **IMPORTANT** : Ne JAMAIS commiter la clé API dans git !

### 3. GitHub Copilot (Alternative)

**Note** : Le script actuel utilise l'API Claude directement. Pour utiliser GitHub Copilot, vous auriez besoin d'une approche différente car :

- GitHub Copilot n'a pas d'API CLI publique pour ce type de tâche
- GitHub Copilot Chat dans VS Code/Zed nécessite interaction manuelle
- Pas d'endpoint API pour batch processing

**Recommandations** :
1. ✅ **Utiliser API Claude** (solution actuelle) - Plus adapté pour automatisation
2. ❌ GitHub Copilot - Nécessite interface graphique, pas d'API batch

---

## 📝 Utilisation

### Mode Automatique Complet

```bash
cd /path/to/tsd
./scripts/review/run_review.sh
```

Le script va :
1. ✅ Vérifier configuration (API key, dépendances)
2. ✅ Exécuter les 6 sessions séquentiellement
3. ✅ Générer un rapport par session dans `REPORTS/`
4. ✅ Créer synthèse globale
5. ✅ Gérer les pauses entre sessions (rate limiting)

### Mode Sans Confirmation

```bash
AUTO_CONFIRM=1 ./scripts/review/run_review.sh
```

### Mode Debug

```bash
DEBUG=1 ./scripts/review/run_review.sh
```

Sauvegarde les prompts dans `/tmp/review_prompt_SESSION_X.txt` pour inspection.

### Exécuter une Session Spécifique

```bash
# Éditer SESSIONS array dans run_review.sh
# Commenter les sessions non désirées
SESSIONS=(
    "SESSION_1_STATE_API"
    # "SESSION_2_VALIDATION"
    # ...
)
```

---

## 📊 Sorties

### Rapports Individuels

Générés dans `REPORTS/` :
- `REVIEW_CONSTRAINT_SESSION_1_STATE_API.md`
- `REVIEW_CONSTRAINT_SESSION_2_VALIDATION.md`
- `REVIEW_CONSTRAINT_SESSION_3_PKG_VALIDATOR.md`
- `REVIEW_CONSTRAINT_SESSION_4_TYPES_DOMAIN.md`
- `REVIEW_CONSTRAINT_SESSION_5_FACTS_ACTIONS.md`
- `REVIEW_CONSTRAINT_SESSION_6_CONFIG_CLI.md`

### Synthèse Globale

`REPORTS/REVIEW_CONSTRAINT_SUMMARY.md` - Agrégation de tous les rapports

---

## ⚙️ Configuration Avancée

### Limites API

**Claude API** :
- **Rate Limit** : ~50 requêtes/minute (dépend du plan)
- **Max Tokens** : 16,000 tokens par réponse (configurable)
- **Context Window** : 200K tokens input

Le script inclut pause de 5s entre sessions pour éviter rate limiting.

### Ajuster le Modèle

Dans `run_review.sh`, ligne ~115 :

```bash
--arg model "claude-sonnet-4-20250514" \
```

Options :
- `claude-sonnet-4-20250514` - Dernier modèle (recommandé)
- `claude-3-5-sonnet-20241022` - Alternative
- `claude-opus-4-20250514` - Plus puissant mais plus lent/coûteux

### Ajuster Max Tokens

Dans `run_review.sh`, ligne ~116 :

```bash
max_tokens: 16000,
```

Augmenter si rapports tronqués (max 4096 pour Sonnet).

---

## 🔧 Troubleshooting

### Erreur : "ANTHROPIC_API_KEY non définie"

```bash
# Vérifier variable
echo $ANTHROPIC_API_KEY

# Si vide, configurer
export ANTHROPIC_API_KEY='votre-clé'
```

### Erreur : "jq: command not found"

```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq
```

### Erreur : "Rate limit exceeded"

- Attendre quelques minutes
- Vérifier quota sur https://console.anthropic.com/
- Augmenter pause entre sessions (ligne ~191)

### Rapports Tronqués

- Augmenter `max_tokens` dans le script
- Diviser fichiers volumineux en sous-sessions
- Utiliser modèle avec plus de capacité (Opus)

### API Timeout

```bash
# Ajouter timeout à curl (ligne ~123)
curl --max-time 300 -s -X POST ...
```

---

## 💰 Coûts Estimés

### API Claude (Anthropic)

**Modèle Sonnet-4** :
- Input : $3 / million tokens
- Output : $15 / million tokens

**Estimation pour 6 sessions** :
- Input : ~300K tokens × $3/M = $0.90
- Output : ~100K tokens × $15/M = $1.50
- **Total** : ~$2.50 pour review complète module

**Note** : Premiers $5 souvent gratuits (nouveaux comptes).

---

## 🔐 Sécurité

### Protection API Key

```bash
# Vérifier que .env n'est pas dans git
git check-ignore .env

# Ajouter à .gitignore si nécessaire
echo '.env' >> .gitignore
echo '*.key' >> .gitignore
```

### Permissions Fichiers

```bash
# Restreindre permissions .env
chmod 600 .env

# Restreindre permissions script
chmod 700 run_review.sh
```

### Rotation Clés

- Régénérer clés API tous les 3-6 mois
- Révoquer immédiatement si exposée
- Ne jamais logguer la clé

---

## 📚 Workflow Recommandé

### 1. Préparation

```bash
# Vérifier que le code est à jour
cd /path/to/tsd
git pull

# Vérifier configuration
./scripts/review/run_review.sh --check  # (à implémenter si besoin)
```

### 2. Exécution

```bash
# Lancer review automatique
AUTO_CONFIRM=1 ./scripts/review/run_review.sh > review.log 2>&1 &

# Suivre progression
tail -f review.log
```

### 3. Analyse

```bash
# Lire synthèse
cat REPORTS/REVIEW_CONSTRAINT_SUMMARY.md

# Prioriser actions
grep "🔴 Critiques" REPORTS/REVIEW_CONSTRAINT_SESSION_*.md
```

### 4. Actions

```bash
# Créer issues GitHub pour chaque problème critique
# Ou créer tickets dans système de suivi

# Planifier corrections selon priorités
```

---

## 🔄 Maintenance Script

### Ajouter Nouvelle Session

1. Créer `SESSION_7_NOUVEAU.md` avec structure standard
2. Ajouter à array `SESSIONS` dans `run_review.sh`
3. Tester : `DEBUG=1 ./run_review.sh`

### Modifier Template Rapport

Éditer section `## 📝 Format Rapport` dans chaque `SESSION_X.md`

### Changer API

Remplacer fonction `call_claude_api()` par autre provider (OpenAI, etc.)

---

## 🆘 Support

### Documentation Officielle

- [Claude API Docs](https://docs.anthropic.com/claude/reference/getting-started-with-the-api)
- [jq Manual](https://stedolan.github.io/jq/manual/)
- [Bash Scripting Guide](https://www.gnu.org/software/bash/manual/)

### Résolution Problèmes

1. Activer mode debug : `DEBUG=1`
2. Vérifier logs : `/tmp/review_prompt_*.txt`
3. Tester API manuellement : `curl https://api.anthropic.com/v1/messages ...`
4. Vérifier quota : https://console.anthropic.com/

---

## 📝 Notes

- Script conçu pour Linux/macOS (bash)
- Testé avec Claude Sonnet-4
- Temps estimation : ~15-30 min pour 6 sessions
- Rapports en français (peut être changé dans prompts)
- Code généré (`parser.go`) exclu automatiquement

---

**Dernière mise à jour** : 2025-12-10  
**Version** : 1.0  
**Auteur** : Assistant IA