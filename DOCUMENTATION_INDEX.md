# 📚 Index de Documentation TSD

**Version** : 2.0.0  
**Dernière mise à jour** : 2026-01-02  
**Statut** : ✅ Documentation Nettoyée et Consolidée

---

## 🎯 Guide Rapide

### Je suis...

#### 👤 Nouvel Utilisateur
1. [README Principal](README.md) - Vue d'ensemble du projet
2. [Installation](docs/installation.md) - Démarrage rapide
3. [Guides Utilisateur](docs/guides.md) - Apprendre TSD
4. [Exemples](examples/) - Programmes d'exemple

#### 💻 Développeur
1. [API Documentation](docs/api.md) - Interface Go
2. [Architecture](docs/architecture.md) - Comprendre RETE
3. [Configuration](docs/configuration.md) - Optimisation
4. [CONTRIBUTING.md](CONTRIBUTING.md) - Contribuer au projet

#### 🔧 Mainteneur
1. [MAINTENANCE_QUICKREF.md](MAINTENANCE_QUICKREF.md) - Référence rapide
2. [.github/prompts/maintain.md](.github/prompts/maintain.md) - Standards maintenance
3. [CHANGELOG.md](CHANGELOG.md) - Historique versions
4. [REPORTS/](REPORTS/) - Rapports techniques

---

## 📁 Structure du Projet

### Racine (9 fichiers essentiels)

```
tsd/
├── README.md                             # 📖 Documentation principale
├── CHANGELOG.md                          # 📝 Historique versions
├── CHANGELOG_v1.1.0.md                   # Archive v1.1.0
├── CHANGELOG_v1.2.0.md                   # Archive v1.2.0
├── CONTRIBUTING.md                       # 🤝 Guide contribution
├── SECURITY.md                           # 🔒 Politique sécurité
├── MAINTENANCE_QUICKREF.md               # 🔧 Référence maintenance
├── TODO_BUILTIN_ACTIONS_INTEGRATION.md   # 📋 TODO actif
└── TODO_VULNERABILITIES.md               # ⚠️ TODO CRITIQUE
```

### Documentation (docs/)

```
docs/
├── README.md                    # Index principal documentation
│
├── 📚 Guides Utilisateur
│   ├── guides.md                # Guide complet débutant → avancé
│   ├── installation.md          # Installation et démarrage
│   ├── user-guide/              # Guides thématiques
│   │   ├── fact-assignments.md
│   │   ├── fact-comparisons.md
│   │   └── type-system.md
│   └── tutorials/               # Tutoriels pas-à-pas
│       └── primary-keys-tutorial.md
│
├── 🔧 Référence Technique
│   ├── reference.md             # Référence complète syntaxe
│   ├── api.md                   # API Go publique
│   ├── architecture.md          # Architecture RETE
│   ├── configuration.md         # Configuration et profils
│   ├── primary-keys.md          # Clés primaires (#field)
│   ├── internal-ids.md          # Système _id_ interne
│   └── no-condition-rules.md    # Règles sans condition
│
├── ⚡ Fonctionnalités
│   ├── actions/                 # [NOUVEAU] Actions consolidées
│   │   ├── README.md            # CRUD + Xuple
│   │   ├── XUPLE_ACTION_IMPLEMENTATION.md
│   │   ├── XUPLE_DEMONSTRATION.md
│   │   └── XUPLE_REPONSE_UTILISATEUR.md
│   └── syntax-changes.md        # [NOUVEAU] Évolution syntaxe
│
├── 🏗️ Architecture Détaillée
│   ├── architecture/
│   │   ├── id-generation.md    # Algorithme génération IDs
│   │   └── diagrams/            # Diagrammes architecture
│   ├── api/
│   │   └── id-generator.md     # API ID Generator
│   └── implementation/          # Détails implémentation
│
├── 🔄 Migration et Historique
│   └── migration/
│       └── from-v1.x.md         # Guide migration v1.x → v2.0
│
└── 📦 Archives
    └── archive/                 # Documentation pré-v2.0
        ├── constraint/
        ├── rete/
        └── pre-v2.0/
```

### Archives (ARCHIVES/)

```
ARCHIVES/
├── cleanup-20260102/            # [NOUVEAU] Nettoyage 2026-01-02
│   ├── README.md                # Documentation archivage
│   └── *.md                     # 83 fichiers archivés
├── DOC_CONSOLIDATION_2025.md
├── architecture/
├── migration/
├── restructuration/
└── sessions/
```

### Rapports (REPORTS/)

```
REPORTS/
├── DEEP_CLEAN_REPORT_2026-01-02.md    # [NOUVEAU] Rapport nettoyage
├── DEEP_CLEAN_REPORT_2025-12-16.md
├── DEAD_CODE_REMOVAL_2025-12-16.md
└── README.md
```

---

## 🔍 Trouver ce que je cherche

### Par Sujet

| Sujet | Document |
|-------|----------|
| **Installation** | [docs/installation.md](docs/installation.md) |
| **Syntaxe TSD** | [docs/reference.md](docs/reference.md) |
| **Actions** | [docs/actions/README.md](docs/actions/README.md) |
| **Clés Primaires** | [docs/primary-keys.md](docs/primary-keys.md) |
| **IDs Internes** | [docs/internal-ids.md](docs/internal-ids.md) |
| **Comparaisons** | [docs/user-guide/fact-comparisons.md](docs/user-guide/fact-comparisons.md) |
| **Affectations** | [docs/user-guide/fact-assignments.md](docs/user-guide/fact-assignments.md) |
| **Types** | [docs/user-guide/type-system.md](docs/user-guide/type-system.md) |
| **API Go** | [docs/api.md](docs/api.md) |
| **Architecture** | [docs/architecture.md](docs/architecture.md) |
| **Configuration** | [docs/configuration.md](docs/configuration.md) |
| **Migration v1.x** | [docs/migration/from-v1.x.md](docs/migration/from-v1.x.md) |
| **Changements Syntaxe** | [docs/syntax-changes.md](docs/syntax-changes.md) |

### Par Niveau

#### Débutant
1. [README.md](README.md) - Vue d'ensemble
2. [docs/installation.md](docs/installation.md) - Installation
3. [docs/guides.md](docs/guides.md) - Guide débutant
4. [docs/user-guide/](docs/user-guide/) - Guides thématiques

#### Intermédiaire
1. [docs/primary-keys.md](docs/primary-keys.md) - Clés primaires
2. [docs/actions/README.md](docs/actions/README.md) - Actions
3. [docs/configuration.md](docs/configuration.md) - Configuration
4. [examples/](examples/) - Exemples complets

#### Avancé
1. [docs/architecture.md](docs/architecture.md) - Architecture RETE
2. [docs/api.md](docs/api.md) - API Go
3. [docs/internal-ids.md](docs/internal-ids.md) - IDs internes
4. [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution

### Par Type

#### Documentation Utilisateur (🇫🇷 Français)
- Tous les fichiers dans `docs/`
- README principal
- CONTRIBUTING.md

#### Documentation Code (🇬🇧 Anglais)
- GoDoc dans le code source
- Commentaires inline

#### Documentation Technique
- [docs/architecture/](docs/architecture/)
- [docs/api/](docs/api/)
- [REPORTS/](REPORTS/)

---

## 🆕 Changements Récents

### Nettoyage Documentation (2026-01-02)

**✅ Terminé** - Voir [REPORTS/DEEP_CLEAN_REPORT_2026-01-02.md](REPORTS/DEEP_CLEAN_REPORT_2026-01-02.md)

#### Améliorations
- ✅ Réduction 87% des fichiers racine (70+ → 9)
- ✅ 83 fichiers archivés (préservés)
- ✅ Documentation consolidée (actions, syntaxe)
- ✅ Structure claire et maintenable

#### Nouveaux Documents
- `docs/syntax-changes.md` - Changements syntaxe consolidés
- `docs/actions/README.md` - Documentation actions consolidée
- `ARCHIVES/cleanup-20260102/README.md` - Documentation archivage

#### Fichiers Archivés
Tous les rapports temporaires et TODO obsolètes ont été déplacés vers :
- `ARCHIVES/cleanup-20260102/` (83 fichiers)

---

## 📋 TODO Actifs

### CRITIQUE ⚠️
- **TODO_VULNERABILITIES.md** - Mise à jour Go 1.24.11+ (sécurité)

### Prioritaire 📌
- **TODO_BUILTIN_ACTIONS_INTEGRATION.md** - Intégration builtin executor

---

## 🤝 Contribution

### Comment Contribuer
1. Lire [CONTRIBUTING.md](CONTRIBUTING.md)
2. Consulter [.github/prompts/common.md](.github/prompts/common.md) - Standards
3. Consulter [.github/prompts/develop.md](.github/prompts/develop.md) - Développement
4. Créer une PR

### Standards Documentation
- **Langue** : Français (utilisateur), Anglais (code)
- **Format** : Markdown
- **Structure** : Organisée par thème
- **Maintenance** : Suivre [maintain.md](.github/prompts/maintain.md)

---

## 🔗 Liens Importants

### Projet
- **GitHub** : [https://github.com/chrlesur/tsd](https://github.com/chrlesur/tsd)
- **Issues** : [GitHub Issues](https://github.com/chrlesur/tsd/issues)

### Documentation
- **Index Principal** : [docs/README.md](docs/README.md)
- **Installation** : [docs/installation.md](docs/installation.md)
- **Guides** : [docs/guides.md](docs/guides.md)
- **API** : [docs/api.md](docs/api.md)

### Maintenance
- **Référence Rapide** : [MAINTENANCE_QUICKREF.md](MAINTENANCE_QUICKREF.md)
- **Standards** : [.github/prompts/maintain.md](.github/prompts/maintain.md)
- **Rapports** : [REPORTS/](REPORTS/)

---

## 📊 Statistiques

### Documentation
- **Fichiers Markdown racine** : 9 (essentiels)
- **Fichiers docs/** : ~88 (organisés)
- **Fichiers archivés** : 83 (cleanup-20260102)
- **TODO actifs** : 2 (pertinents)

### Projet
- **Version** : 2.0.0
- **Go minimum** : 1.21+
- **Statut** : ✅ Production Ready

---

## 🆘 Aide

### Je ne trouve pas...

#### Un fichier ancien
→ Vérifier dans [ARCHIVES/cleanup-20260102/](ARCHIVES/cleanup-20260102/)

#### Comment faire quelque chose
→ Commencer par [docs/guides.md](docs/guides.md)

#### Documentation API
→ Voir [docs/api.md](docs/api.md)

#### Problème technique
→ Consulter [GitHub Issues](https://github.com/chrlesur/tsd/issues)

---

## 📝 Notes

### Fichiers Archivés
Les fichiers dans `ARCHIVES/cleanup-20260102/` sont **préservés** mais **non maintenus**.
Pour information historique uniquement.

### Documentation Archive
La documentation pré-v2.0 est dans `docs/archive/`.
Consulter `docs/migration/from-v1.x.md` pour migration.

### Maintenance
Pour maintenir la documentation, suivre les standards dans :
- [.github/prompts/maintain.md](.github/prompts/maintain.md)
- [MAINTENANCE_QUICKREF.md](MAINTENANCE_QUICKREF.md)

---

**Bon développement avec TSD ! 🚀**