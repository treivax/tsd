# Documentation TSD

Bienvenue dans la documentation du projet **TSD** (Type System Development) - Un moteur de règles basé sur l'algorithme RETE.

## 🚀 Démarrage Rapide

- **Nouveau ?** Commencez par le [Guide de Démarrage Rapide](QUICK_START.md)
- **Installation** : Consultez le [Guide d'Installation](INSTALLATION.md)
- **Premier exemple** : Suivez le [Tutorial](TUTORIAL.md)

## 📚 Documentation par Catégorie

### 🎓 Guides Utilisateur

| Document | Description | Audience |
|----------|-------------|----------|
| [Quick Start](QUICK_START.md) | Démarrage rapide en 5 minutes | Débutant |
| [Tutorial](TUTORIAL.md) | Tutorial complet avec exemples | Débutant |
| [User Guide](USER_GUIDE.md) | Guide utilisateur complet | Intermédiaire |
| [Grammar Guide](GRAMMAR_GUIDE.md) | Grammaire du langage TSD | Intermédiaire |

### ⚙️ Configuration

| Document | Description | Audience |
|----------|-------------|----------|
| [Configuration Globale](configuration/README.md) | **Guide complet de configuration** | Tous |
| [RETE Configuration](RETE_CONFIGURATION.md) | Configuration réseau RETE | Avancé |
| [Logging Guide](LOGGING_GUIDE.md) | Configuration du logging | Tous |
| [Authentication](AUTHENTICATION.md) | Configuration authentification | Admin |

### 🏗️ Architecture

| Document | Description | Audience |
|----------|-------------|----------|
| [Architecture](ARCHITECTURE.md) | Vue d'ensemble architecture | Développeur |
| [Working Memory](WORKING_MEMORY.md) | Gestion de la mémoire de travail | Développeur |
| [Bindings Design](architecture/BINDINGS_DESIGN.md) | Design des bindings | Avancé |
| [Bindings Performance](architecture/BINDINGS_PERFORMANCE.md) | Analyse performance | Avancé |

### 🔌 API & Intégration

| Document | Description | Audience |
|----------|-------------|----------|
| [API Reference](API_REFERENCE.md) | Référence API complète | Développeur |
| [Public API](api/PUBLIC_API.md) | API publique | Développeur |

### 🤝 Contribution

| Document | Description | Audience |
|----------|-------------|----------|
| [Contributing](CONTRIBUTING.md) | Guide de contribution | Contributeur |
| [../CHANGELOG.md](../CHANGELOG.md) | Historique des changements | Tous |

---

## 📖 Structure de la Documentation

```
docs/
├── README.md                          # Ce fichier - Index global
│
├── guides/                            # Guides utilisateur
│   ├── (en construction)
│
├── configuration/                     # Configuration système
│   └── README.md                      # ★ Guide configuration complet
│
├── api/                               # Documentation API
│   └── PUBLIC_API.md                  # API publique
│
├── architecture/                      # Architecture & Design
│   ├── BINDINGS_DESIGN.md
│   ├── BINDINGS_PERFORMANCE.md
│   ├── BINDINGS_ANALYSIS.md
│   ├── BINDINGS_STATUS_REPORT.md
│   └── CODE_REVIEW_BINDINGS.md
│
├── QUICK_START.md                     # Démarrage rapide
├── INSTALLATION.md                    # Installation
├── TUTORIAL.md                        # Tutorial complet
├── USER_GUIDE.md                      # Guide utilisateur
├── GRAMMAR_GUIDE.md                   # Grammaire TSD
├── ARCHITECTURE.md                    # Architecture système
├── API_REFERENCE.md                   # Référence API
├── RETE_CONFIGURATION.md              # Configuration RETE
├── LOGGING_GUIDE.md                   # Configuration logging
├── AUTHENTICATION.md                  # Authentification
├── WORKING_MEMORY.md                  # Working Memory
├── CONTRIBUTING.md                    # Guide contribution
└── INMEMORY_ONLY_MIGRATION.md         # Migration stockage
```

---

## 🎯 Documentation par Cas d'Usage

### Je veux...

#### ...démarrer rapidement
1. [Quick Start](QUICK_START.md) - 5 minutes
2. [Installation](INSTALLATION.md) - Installer TSD
3. [Tutorial](TUTORIAL.md) - Premier exemple

#### ...comprendre le langage TSD
1. [Grammar Guide](GRAMMAR_GUIDE.md) - Syntaxe complète
2. [User Guide](USER_GUIDE.md) - Utilisation avancée
3. [API Reference](API_REFERENCE.md) - API programmatique

#### ...configurer TSD pour mon cas d'usage
1. ⭐ [Configuration Globale](configuration/README.md) - **Démarrer ici**
2. [RETE Configuration](RETE_CONFIGURATION.md) - Configuration moteur
3. [Logging Guide](LOGGING_GUIDE.md) - Logs et debugging

#### ...déployer en production
1. [Configuration Globale](configuration/README.md) - Profil Production
2. [Authentication](AUTHENTICATION.md) - Sécuriser l'API
3. [RETE Configuration](RETE_CONFIGURATION.md) - Optimiser performance

#### ...comprendre l'architecture
1. [Architecture](ARCHITECTURE.md) - Vue d'ensemble
2. [Working Memory](WORKING_MEMORY.md) - Gestion mémoire
3. [Bindings Design](architecture/BINDINGS_DESIGN.md) - Design interne

#### ...contribuer au projet
1. [Contributing](CONTRIBUTING.md) - Guide contribution
2. [Architecture](ARCHITECTURE.md) - Comprendre le système
3. [Code Review](architecture/CODE_REVIEW_BINDINGS.md) - Standards

---

## 🔧 Configuration Rapide

### Profils Prédéfinis

```go
// Développement - Logs détaillés, pas de cache
config := rete.DefaultChainPerformanceConfig()
logger := rete.NewLogger(rete.LogLevelDebug, os.Stdout)

// Production - Performance maximale
config := rete.HighPerformanceConfig()
config.PrometheusEnabled = true

// Embarqué - Mémoire minimale
config := rete.LowMemoryConfig()

// Tests - Déterministe
config := rete.DisabledCachesConfig()
```

📖 **Détails** : [Configuration Globale](configuration/README.md)

---

## 📊 Composants Configurables

| Composant | Configuration | Documentation |
|-----------|---------------|---------------|
| **Réseau RETE** | `ChainPerformanceConfig` | [RETE Config](RETE_CONFIGURATION.md) |
| **Transactions** | `TransactionOptions` | [Config Globale](configuration/README.md#transactionoptions) |
| **Beta Sharing** | `BetaSharingConfig` | [RETE Config](RETE_CONFIGURATION.md#betasharingconfig) |
| **Logger** | `Logger`, `LogLevel` | [Logging Guide](LOGGING_GUIDE.md) |
| **Serveur** | `ServerConfig` | [Config Globale](configuration/README.md#server) |
| **Client** | `ClientConfig` | [Config Globale](configuration/README.md#client) |
| **Auth** | `AuthConfig` | [Authentication](AUTHENTICATION.md) |
| **Storage** | `StorageConfig` | [Config Globale](configuration/README.md#storage) |

---

## 🎓 Parcours d'Apprentissage

### Parcours Débutant (2-4 heures)

1. ✅ [Quick Start](QUICK_START.md) - 15 min
2. ✅ [Installation](INSTALLATION.md) - 15 min
3. ✅ [Tutorial](TUTORIAL.md) - 1h
4. ✅ [Grammar Guide](GRAMMAR_GUIDE.md) - 1h
5. ✅ [Configuration Globale](configuration/README.md) - 1h

### Parcours Développeur (1-2 jours)

1. ✅ Parcours Débutant
2. ✅ [User Guide](USER_GUIDE.md) - 2h
3. ✅ [API Reference](API_REFERENCE.md) - 2h
4. ✅ [Architecture](ARCHITECTURE.md) - 2h
5. ✅ [RETE Configuration](RETE_CONFIGURATION.md) - 2h
6. ✅ [Working Memory](WORKING_MEMORY.md) - 1h

### Parcours Avancé (1 semaine)

1. ✅ Parcours Développeur
2. ✅ [Bindings Design](architecture/BINDINGS_DESIGN.md)
3. ✅ [Bindings Performance](architecture/BINDINGS_PERFORMANCE.md)
4. ✅ [Code Review](architecture/CODE_REVIEW_BINDINGS.md)
5. ✅ Code source : `/rete`, `/constraint`

---

## 🔍 Recherche Rapide

### Par Sujet

- **Authentification** : [AUTHENTICATION.md](AUTHENTICATION.md)
- **API** : [API_REFERENCE.md](API_REFERENCE.md)
- **Caches** : [RETE_CONFIGURATION.md](RETE_CONFIGURATION.md#cache-de-hash)
- **Configuration** : [configuration/README.md](configuration/README.md) ⭐
- **Déploiement** : [configuration/README.md#production](configuration/README.md#production)
- **Grammaire** : [GRAMMAR_GUIDE.md](GRAMMAR_GUIDE.md)
- **Installation** : [INSTALLATION.md](INSTALLATION.md)
- **Logging** : [LOGGING_GUIDE.md](LOGGING_GUIDE.md)
- **Performance** : [RETE_CONFIGURATION.md](RETE_CONFIGURATION.md#profils-prédéfinis)
- **RETE** : [RETE_CONFIGURATION.md](RETE_CONFIGURATION.md)
- **Transactions** : [configuration/README.md#transactionoptions](configuration/README.md#transactionoptions)

### Par Type d'Utilisation

- **CLI** : [USER_GUIDE.md](USER_GUIDE.md)
- **API Programmatique** : [API_REFERENCE.md](API_REFERENCE.md)
- **Serveur HTTP/HTTPS** : [configuration/README.md#server](configuration/README.md#server)
- **Docker** : [configuration/README.md#exemple-2--production-avec-docker](configuration/README.md#exemple-2--production-avec-docker)

---

## 📝 Conventions de Documentation

### Langue

- **GoDoc** : Anglais (convention Go)
- **Documentation utilisateur** : Français
- **Commentaires code** : Français
- **README modules** : Français

### Format

- **Markdown** : GitHub Flavored Markdown
- **Code** : Blocs avec syntaxe highlighting
- **Exemples** : Testables et fonctionnels

### Standards

Voir [.github/prompts/document.md](../.github/prompts/document.md) pour les standards complets.

---

## 🆘 Besoin d'Aide ?

### Questions Fréquentes

**Q: Comment configurer TSD pour la production ?**  
A: Consultez [Configuration Globale - Profil Production](configuration/README.md#production)

**Q: Comment activer HTTPS ?**  
A: Consultez [Configuration Serveur HTTPS](configuration/README.md#exemple-https)

**Q: Comment optimiser les performances ?**  
A: Consultez [RETE Configuration - High Performance](RETE_CONFIGURATION.md#configuration-haute-performance)

**Q: Quelle est la différence entre les profils de config ?**  
A: Consultez [Profils de Déploiement](configuration/README.md#profils-de-déploiement)

### Support

1. **Documentation** : Cherchez dans cet index
2. **Issues** : [GitHub Issues](https://github.com/yourusername/tsd/issues)
3. **Contribution** : [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 🚀 Prochaines Étapes

Après avoir lu la documentation :

1. ✅ Installer TSD : [INSTALLATION.md](INSTALLATION.md)
2. ✅ Faire le tutorial : [TUTORIAL.md](TUTORIAL.md)
3. ✅ Configurer pour votre cas : [configuration/README.md](configuration/README.md)
4. ✅ Lire le guide utilisateur : [USER_GUIDE.md](USER_GUIDE.md)
5. ✅ Contribuer : [CONTRIBUTING.md](CONTRIBUTING.md)

---

## 📄 License

TSD est distribué sous licence MIT. Voir [LICENSE](../LICENSE) pour plus de détails.

---

**Version** : 1.0.0  
**Dernière mise à jour** : 2025-01-XX  
**Mainteneur** : TSD Contributors

💡 **Astuce** : Marquez cette page pour y revenir facilement !