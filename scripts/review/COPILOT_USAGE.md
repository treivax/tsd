# 🤖 Guide Utilisation review.sh avec GitHub Copilot CLI

Guide complet pour exécuter les reviews automatisées avec GitHub Copilot CLI.

---

## 📋 Prérequis

### 1. Installation Node.js et npm

**Ubuntu 25.10** :
```bash
# Via NodeSource (recommandé)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Vérifier
node --version  # v20.x.x
npm --version   # v10.x.x
```

### 2. Installation GitHub Copilot CLI

```bash
# Installer globalement
npm install -g @githubnext/github-copilot-cli

# Vérifier installation
copilot --version
```

### 3. Authentification GitHub Copilot

```bash
# Lancer authentification
copilot auth login

# Suivre instructions navigateur
# Nécessite abonnement GitHub Copilot actif
```

**Vérifier abonnement** : https://github.com/settings/copilot

---

## 🚀 Utilisation

### Mode Interactif (Recommandé pour Première Fois)

```bash
cd /path/to/tsd
./scripts/review/review.sh
```

**Le script va** :
1. ✅ Détecter fichiers `SESSION*.md` (ordre alphabétique)
2. ✅ Afficher liste des sessions
3. ⏸️ Demander confirmation
4. 🚀 Exécuter chaque session séquentiellement
5. ⏸️ Pause 10s entre sessions
6. 📊 Afficher résumé final

### Mode Automatique (Sans Confirmation)

```bash
./scripts/review/review.sh -y
# ou
./scripts/review/review.sh --yes
# ou
AUTO_CONFIRM=1 ./scripts/review/review.sh
```

### Mode Continue (Continuer sur Erreur)

```bash
./scripts/review/review.sh -y -c
# ou
./scripts/review/review.sh --yes --continue
# ou
AUTO_CONFIRM=1 AUTO_CONTINUE=1 ./scripts/review/review.sh
```

### Personnaliser Pause Entre Sessions

```bash
# Pause 30 secondes
./scripts/review/review.sh -p 30

# Pause 5 secondes (rapide)
./scripts/review/review.sh -y -p 5

# Via variable
PAUSE_SECONDS=60 ./scripts/review/review.sh -y
```

### Combinaisons Utiles

```bash
# Automatique complet avec pause courte
./scripts/review/review.sh -y -c -p 5

# Background avec log
nohup ./scripts/review/review.sh -y -c > review.log 2>&1 &
tail -f review.log
```

---

## 📊 Ce Que Fait Chaque Session

Le script exécute pour **chaque** `SESSION_X.md` :

```bash
copilot -p "Execute, as the linux user resinsec, the prompt \
  .github/prompts/review.md (de l'analyse jusqu'au refactoring \
  du code que tu dois mener en appliquant l'ensemble des \
  préconisations et solutions identifiées) en l'appliquant sur \
  le périmètre et les contraintes définis dans \
  scripts/review/SESSION_X.md ainsi que les règles et bonnes \
  pratiques définies dans .github/prompts/common.md. \
  Effectue les modifications sans conservation de l'existant \
  même si elles impliquent une modification du code qui utilise \
  cet existant. Dans le cas où le nouveau code ne serait pas \
  compatible avec l'existant, si tu ne peux corriger le code \
  appelant, décris clairement en TODO les actions qui seront \
  nécessaires pour rendre fonctionnel le code qui utilisera \
  les modifications faites." \
  --allow-all-tools
```

### Ordre d'Exécution

Les fichiers sont traités par **ordre lexicographique** :

1. `SESSION_1_STATE_API.md` → State Management & API
2. `SESSION_2_VALIDATION.md` → Validation Layer
3. `SESSION_3_PKG_VALIDATOR.md` → Package Validator
4. `SESSION_4_TYPES_DOMAIN.md` → Types & Domain
5. `SESSION_5_FACTS_ACTIONS.md` → Facts, Actions & Logic
6. `SESSION_6_CONFIG_CLI.md` → Config & CLI

### Actions Par Session

Copilot CLI va :
1. **Analyser** le code selon checklist `review.md`
2. **Identifier** problèmes (critiques, majeurs, mineurs)
3. **Refactorer** le code directement
4. **Appliquer** toutes préconisations
5. **Modifier** fichiers source
6. **Créer TODO** si code appelant incompatible
7. **Générer** rapport dans `REPORTS/`

---

## ⚠️ Comportement Important

### Modifications Directes du Code

**Le script modifie le code source directement** :
- ✅ Pas de conservation de l'existant
- ✅ Refactoring appliqué immédiatement
- ✅ Fichiers modifiés en place

**Avant de lancer** :
```bash
# Commiter tout changement en cours
git status
git add .
git commit -m "Avant review automatique"

# Ou créer branche dédiée
git checkout -b review-automatique-constraint
```

### Gestion Incompatibilités

Si modifications cassent code appelant :
- Copilot tentera de corriger le code appelant
- Si impossible : ajout **TODO** avec actions nécessaires
- Rechercher après : `grep -r "TODO" constraint/`

---

## 📈 Suivi Progression

### Logs en Direct

```bash
# Terminal 1 : Lancer review
./scripts/review/review.sh -y -c

# Terminal 2 : Suivre logs (si applicable)
watch -n 2 'ls -lht REPORTS/REVIEW_* | head -10'
```

### Interruption

```bash
# Ctrl+C pour arrêter
# Le script s'arrête proprement après session en cours

# Relancer reprend à la session suivante
# (sessions déjà traitées visibles dans REPORTS/)
```

---

## ✅ Après Exécution

### Vérifier Modifications

```bash
# Voir fichiers modifiés
git status

# Diff complet
git diff

# Par fichier
git diff constraint/program_state.go
```

### Rechercher TODO

```bash
# Trouver tous les TODO ajoutés
grep -r "TODO" constraint/ --color=always

# Avec contexte
grep -r -B2 -A2 "TODO" constraint/
```

### Tests de Non-Régression

```bash
# Tests unitaires
make test-unit

# Tests complets
make test-complete

# Vérifier compilation
make build
```

### Valider Qualité

```bash
# Linting
make lint

# Formatage
make format

# Validation complète
make validate
```

---

## 🔧 Troubleshooting

### "copilot: command not found"

```bash
# Installer Copilot CLI
npm install -g @githubnext/github-copilot-cli

# Vérifier PATH npm global
npm config get prefix
echo $PATH
```

### "Authentication required"

```bash
# Se connecter
copilot auth login

# Vérifier abonnement GitHub Copilot actif
# https://github.com/settings/copilot
```

### "No subscription found"

- Vérifier abonnement GitHub Copilot actif
- Compte individuel : $10/mois
- Via organisation : vérifier accès
- Free trial disponible : https://github.com/features/copilot

### Session Échoue

```bash
# Relancer session spécifique manuellement
cd /path/to/tsd
copilot -p "Execute review.md on SESSION_X.md..." --allow-all-tools

# Ou ignorer et continuer
./scripts/review/review.sh -y -c
```

### Modifications Trop Importantes

```bash
# Restaurer avant review
git checkout .

# Ou cherry-pick modifications utiles
git add -p

# Ou traiter sessions individuellement
# Commenter sessions dans script (à implémenter)
```

---

## 💡 Bonnes Pratiques

### Avant Review

1. **Commit propre** :
   ```bash
   git status
   git add .
   git commit -m "État avant review"
   ```

2. **Branche dédiée** :
   ```bash
   git checkout -b review-constraint-$(date +%Y%m%d)
   ```

3. **Tests passent** :
   ```bash
   make test-complete
   ```

### Pendant Review

1. **Suivre progression** : Logs dans terminal
2. **Ne pas interrompre** : Laisser session terminer
3. **Surveiller ressources** : CPU/Mémoire

### Après Review

1. **Vérifier modifications** :
   ```bash
   git diff --stat
   git diff
   ```

2. **Tester immédiatement** :
   ```bash
   make test-complete
   make validate
   ```

3. **Analyser TODO** :
   ```bash
   grep -r "TODO" constraint/ > todo_list.txt
   ```

4. **Commit granulaire** :
   ```bash
   # Par type de changement
   git add constraint/program_state*.go
   git commit -m "refactor(state): amélioration gestion état"
   
   git add constraint/*_validation*.go
   git commit -m "fix(validation): renforcement validation"
   ```

---

## 📊 Estimation

### Durée

- **Par session** : 5-15 minutes (dépend complexité)
- **Total (6 sessions)** : 30-90 minutes
- **Avec pauses (10s)** : +1 minute

### Ressources

- **CPU** : Modéré (CLI + Copilot API)
- **Réseau** : Connexion stable requise
- **Mémoire** : ~500MB (Node.js + Copilot)

---

## 🆘 Support

### Documentation Officielle

- **Copilot CLI** : https://www.npmjs.com/package/@githubnext/github-copilot-cli
- **GitHub Copilot** : https://github.com/features/copilot
- **Node.js** : https://nodejs.org/docs

### Commandes Utiles

```bash
# Version
copilot --version

# Aide
copilot --help

# Status authentification
copilot auth status

# Se déconnecter
copilot auth logout
```

---

## 🔄 Alternatives

### Si Copilot CLI Ne Fonctionne Pas

1. **API Claude** : Utiliser `run_review.sh` (déjà créé)
2. **Manuel dans Zed** : Charger sessions une par une
3. **VSCode Copilot** : Copier sessions dans chat

---

## 📝 Notes

- Script conçu pour Linux (bash)
- Testé sur Ubuntu 25.10
- Nécessite abonnement GitHub Copilot actif
- Modifications appliquées directement (git recommandé)
- Ordre sessions optimisé (critique → basse priorité)

---

**Version** : 1.0  
**Dernière mise à jour** : 2025-12-10  
**Compatible** : GitHub Copilot CLI v1.x